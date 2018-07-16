package p2p

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/gammazero/nexus/wamp"

	nx_router "github.com/gammazero/nexus/router"

	"github.com/Lobarr/lane"
	upnp "github.com/NebulousLabs/go-upnp"
	"github.com/hprose/hprose-golang/rpc"

	"github.com/gizo-network/gizo/core/difficulty"
	"github.com/gizo-network/gizo/core/merkletree"

	"github.com/gizo-network/gizo/job"

	"github.com/gizo-network/gizo/helpers"
	"github.com/gizo-network/gizo/job/queue"
	funk "github.com/thoas/go-funk"

	"github.com/boltdb/bolt"

	"github.com/kpango/glg"

	"github.com/gammazero/nexus/client"
	"github.com/gizo-network/gizo/benchmark"
	"github.com/gizo-network/gizo/cache"
	"github.com/gizo-network/gizo/core"
	"github.com/gizo-network/gizo/crypt"
	"github.com/gorilla/mux"
)

var (
	//ErrJobsFull occurs when jobs map is full
	ErrJobsFull = errors.New("Jobs map full")
)

//TODO: verify job hasn't been modified

//Dispatcher dispatcher node
type Dispatcher struct {
	ip        string
	port      uint   //port
	pub       []byte //public key of the node
	priv      []byte //private key of the node
	uptime    int64  //time since node has been up
	jobPQ     *queue.JobPriorityQueue
	workers   map[string]*WorkerInfo //worker nodes in dispatcher's area
	workerPQ  *WorkerPriorityQueue
	peers     map[string]*client.Client
	bench     benchmark.Engine //benchmark of node
	wamp      nx_router.Router
	rpc       *rpc.HTTPService // rpc servce
	router    *mux.Router
	jc        *cache.JobCache  //job cache
	bc        *core.BlockChain //blockchain
	db        *bolt.DB         //holds topology table
	mu        *sync.Mutex
	jobs      map[string]job.Job // holds done jobs and new jobs submitted to the network before being placed in the bc
	interrupt chan os.Signal     //used watch interrupts
	writeQ    *lane.Queue        // queue of job (execs) to be written to the db
	centrum   *Centrum
	discover  *upnp.IGD //used for upnp external ip discovery
	new       bool      // if true, sends a new dispatcher to centrum else sends a wake with token
	dClient   *client.Client
	wClient   *client.Client
}

//GetJobs returns jobs held in memory to be written to the bc
func (d Dispatcher) GetJobs() map[string]job.Job {
	return d.jobs
}

//AddPeer adds peer connection
func (d *Dispatcher) AddPeer(p string, c *client.Client) {
	d.peers[p] = c
}

//GetPeers returns peers
func (d Dispatcher) GetPeers() map[string]*client.Client {
	return d.peers
}

//watches the queue of jobs to be written to the bc
func (d Dispatcher) watchWriteQ() {
	//TODO: write to bc if it's taking too long
	for {
		if !d.GetWriteQ().Empty() {
			jobs := d.GetWriteQ().Dequeue()
			d.WriteJobsAndPublish(jobs.(map[string]job.Job))
		}
	}
}

//WriteJobsAndPublish writes jobs to the bc
func (d Dispatcher) WriteJobsAndPublish(jobs map[string]job.Job) {
	nodes := []*merkletree.MerkleNode{}
	for _, job := range jobs {
		node, err := merkletree.NewNode(job, nil, nil)
		if err != nil {
			glg.Fatal(err)
		}
		nodes = append(nodes, node)
	}
	tree, err := merkletree.NewMerkleTree(nodes)
	if err != nil {
		glg.Fatal(err)
	}
	latestBlock, err := d.GetBC().GetLatestBlock()
	if err != nil {
		glg.Fatal(err)
	}
	nextHeight, err := d.GetBC().GetNextHeight()
	if err != nil {
		glg.Fatal(err)
	}
	block, err := core.NewBlock(*tree, latestBlock.GetHeader().GetHash(), nextHeight, uint8(difficulty.Difficulty(d.GetBenchmarks(), *d.GetBC())), d.GetPubString())
	if err != nil {
		glg.Fatal(err)
	}
	err = d.GetBC().AddBlock(block)
	if err != nil {
		glg.Fatal(err)
	}
	bBytes, _ := helpers.Serialize(block)
	err = d.dClient.Publish(BLOCK, nil, wamp.List{string(bBytes)}, nil)
	if err != nil {
		glg.Fatal(err)
	}
}

//assigns next job in the queue to the next available worker
func (d *Dispatcher) deployJobs() {
	for {
		if !d.GetWorkerPQ().getPQ().Empty() {
			if !d.GetJobPQ().GetPQ().Empty() {
				d.mu.Lock()
				w := d.GetWorkerPQ().Pop()
				if !d.GetWorker(w).GetShut() {
					j := d.GetJobPQ().Pop()
					if j.GetExec().GetStatus() != job.CANCELLED {
						j.GetExec().SetBy(d.GetWorker(w).GetPub())
						d.GetWorker(w).Assign(&j)
						glg.Info("P2P: dispatched job")
						worker := d.GetWorker(w)
						jBytes, _ := helpers.Serialize(j)
						err := d.wClient.Publish(worker.JobTopic(), nil, wamp.List{string(jBytes)}, nil)
						if err != nil {
							glg.Fatal(err)
						}
					} else {
						j.ResultsChan() <- j
					}
				} else {
					delete(d.GetWorkers(), w)
				}
				d.mu.Unlock()
			}
		}
	}
}

//BlockSubscribe subscribe and publish to block topic
func (d *Dispatcher) BlockSubscribe(peer *client.Client) {
	handler := func(args wamp.List, kwargs wamp.Dict, details wamp.Dict) {
		var block *core.Block
		blk, _ := wamp.AsString(args[0])
		err := helpers.Deserialize([]byte(blk), &block)
		if err != nil {
			return //TODO: handle error
		}
		check, _ := d.GetBC().GetBlockInfo(block.GetHeader().GetHash())
		if check == nil {
			d.GetBC().AddBlock(block)
			err := d.dClient.Publish(BLOCK, nil, wamp.List{blk}, nil)
			if err != nil {
				glg.Fatal(err)
			}
		}
	}
	peer.Subscribe(BLOCK, handler, nil)
}

//WorkerDisconnect handles worker disconnect
func (d Dispatcher) WorkerDisconnect() {
	//TODO: handle unexpected disconnect
	handler := func(args wamp.List, kwargs wamp.Dict, details wamp.Dict) {
		d.mu.Lock()
		glg.Log("Dispatcher: worker disconnected")
		worker := args[0].(string)
		d.GetWorker(worker).SetShut(true)
		if d.GetWorker(worker).GetJob() != nil {
			d.GetJobPQ().PushItem(*d.GetWorker(worker).GetJob(), job.BOOST)
		}
		d.mu.Unlock()
	}
	d.wClient.Subscribe(WORKERDISCONNECT, handler, nil)
}

//BlockReq handles peer request for block
func (d Dispatcher) BlockReq(ctx context.Context, args wamp.List, kwargs, details wamp.Dict) *client.InvokeResult {
	blockinfo, _ := d.GetBC().GetBlockInfo(args[0].(string))
	block := blockinfo.GetBlock()
	bBytes, _ := helpers.Serialize(block)
	return &client.InvokeResult{Args: wamp.List{string(bBytes)}}
}

//WorkerConnect handles workers request to join area
func (d *Dispatcher) WorkerConnect(ctx context.Context, args wamp.List, kwargs, details wamp.Dict) *client.InvokeResult {
	if len(d.GetWorkers()) < MaxWorkers {
		worker, _ := wamp.AsString(args[0])
		d.AddWorker(worker)
		d.centrum.ConnectWorker()
		d.GetWorkerPQ().Push(worker, 0)
		//handle results
		handler := func(args wamp.List, kwargs, details wamp.Dict) {
			d.mu.Lock()
			var exec *job.Exec
			execStr, _ := wamp.AsString(args[0])
			err := helpers.Deserialize([]byte(execStr), &exec)
			if err != nil {
				return //TODO: handle error's better
			}
			d.GetWorker(worker).GetJob().SetExec(exec)
			d.GetWorker(worker).GetJob().ResultsChan() <- *d.GetWorker(worker).GetJob()
			j := d.GetWorker(worker).GetJob().GetJob()
			d.GetWorker(worker).SetJob(nil)
			j.AddExec(*exec)
			d.AddJob(j)
			//TODO: send to requester's message broker
			if !d.GetWorker(worker).GetShut() {
				d.GetWorkerPQ().Push(worker, 0)
			}
			d.mu.Unlock()
		}
		d.wClient.Subscribe(d.GetWorker(worker).ResultTopic(), handler, nil)
		return &client.InvokeResult{Args: wamp.List{CONNECTED}}
	}
	return &client.InvokeResult{Args: wamp.List{CONNFULL}}
}

//AddJob keeps job in memory before being written to the bc
func (d *Dispatcher) AddJob(j job.Job) {
	if len(d.GetJobs()) < merkletree.MaxTreeJobs {
		job, ok := d.GetJobs()[j.GetID()]
		if ok {
			job.AddExec(j.GetLatestExec())
			d.jobs[j.GetID()] = job
			return
		}
		d.jobs[j.GetID()] = j
	} else {
		d.GetWriteQ().Enqueue(d.GetJobs())
		d.EmptyJobs()
		d.jobs[j.GetID()] = j
	}
}

//EmptyJobs empties the jobs held in memory
func (d *Dispatcher) EmptyJobs() {
	d.jobs = make(map[string]job.Job)
}

//GetWorkerPQ returns the workers priotity queue
func (d Dispatcher) GetWorkerPQ() *WorkerPriorityQueue {
	return d.workerPQ
}

//GetAssignedWorker returns worker assigned to execute job
func (d Dispatcher) GetAssignedWorker(hash string) string {
	for key, val := range d.GetWorkers() {
		if val.GetJob().GetJob().GetHash() == hash {
			return key
		}
	}
	return ""
}

//GetIP returns ip address of node
func (d Dispatcher) GetIP() string {
	return d.ip
}

//SetIP sets ip address of node
func (d *Dispatcher) SetIP(ip string) {
	d.ip = ip
}

//GetPubByte returns public key of node as bytes
func (d Dispatcher) GetPubByte() []byte {
	return d.pub
}

//GetPubString returns public key of node as bytes
func (d Dispatcher) GetPubString() string {
	return hex.EncodeToString(d.pub)
}

//GetPrivByte returns private key of node as bytes
func (d Dispatcher) GetPrivByte() []byte {
	return d.priv
}

//GetPrivString returns private key of node as string
func (d Dispatcher) GetPrivString() string {
	return hex.EncodeToString(d.priv)
}

//GetWriteQ returns queue of jobs to be written to the bc
func (d Dispatcher) GetWriteQ() *lane.Queue {
	return d.writeQ
}

//GetJobPQ returns priority of job execs
func (d Dispatcher) GetJobPQ() *queue.JobPriorityQueue {
	return d.jobPQ
}

//GetWorkers returns workers in the standard area
func (d Dispatcher) GetWorkers() map[string]*WorkerInfo {
	return d.workers
}

//GetWorker returns specified worker
func (d Dispatcher) GetWorker(s string) *WorkerInfo {
	return d.GetWorkers()[s]
}

//AddWorker sets worker
func (d *Dispatcher) AddWorker(pub string) {
	d.GetWorkers()[pub] = NewWorkerInfo(pub)
}

//GetBC returns blockchain object
func (d Dispatcher) GetBC() *core.BlockChain {
	return d.bc
}

//GetPort returns node port
func (d Dispatcher) GetPort() int {
	return int(d.port)
}

//GetUptime returns uptime of node
func (d Dispatcher) GetUptime() int64 {
	return d.uptime
}

//GetUptimeString returns uptime of node as string
func (d Dispatcher) GetUptimeString() string {
	return time.Unix(d.uptime, 0).Sub(time.Now()).String()
}

//GetBench returns node benchmark
func (d Dispatcher) GetBench() benchmark.Engine {
	return d.bench
}

//GetBenchmarks return node benchmarks per difficulty
func (d Dispatcher) GetBenchmarks() []benchmark.Benchmark {
	return d.bench.GetData()
}

//GetJC returns job cache object
func (d Dispatcher) GetJC() *cache.JobCache {
	return d.jc
}

//GetRPC returns hprose rpc http server
func (d Dispatcher) GetRPC() *rpc.HTTPService {
	return d.rpc
}

func (d Dispatcher) setRPC(s *rpc.HTTPService) {
	d.rpc = s
}

func (d Dispatcher) watchInterrupt() {
	select {
	case i := <-d.interrupt:
		glg.Warn("Dispatcher: interrupt detected")
		switch i {
		case syscall.SIGINT, syscall.SIGTERM:
			//TODO: get jobs from workers, add all jobs held in memory to block, broadcast block to the network
			res, err := d.centrum.Sleep()
			if err != nil {
				glg.Fatal(err)
			} else if res["status"].(string) != "success" {
				glg.Fatal("Centrum: " + res["status"].(string))
			}
			d.dClient.Close()
			d.wClient.Close()
			d.wamp.Close()
			os.Exit(0)
		case syscall.SIGQUIT:
			os.Exit(1)
		}
	}
}

//Start spins up services
func (d Dispatcher) Start() {
	go d.deployJobs()
	go d.watchWriteQ()
	go d.watchInterrupt()
	d.GetDispatchersAndSync()
	verify, err := d.GetBC().Verify()
	if err != nil {
		glg.Fatal(err)
	}
	if !verify {
		glg.Fatal("Dispatcher: blockchain not verified")
	}
	d.WorkerDisconnect()
	if err = d.dClient.Register(BLOCKREQ, d.BlockReq, nil); err != nil {
		glg.Fatal(err)
	}
	if err = d.wClient.Register(WORKERCONNECT, d.WorkerConnect, nil); err != nil {
		glg.Fatal(err)
	}
	d.RPC()
	wampRouter := nx_router.NewWebsocketServer(d.wamp)
	wampRouter.Upgrader.EnableCompression = true
	wampRouter.Upgrader.ReadBufferSize = 1000000
	wampRouter.Upgrader.WriteBufferSize = 1000000
	d.router.Handle("/rpc", d.GetRPC()).Methods("POST")
	d.router.Handle("/wamp", wampRouter)
	status := make(map[string]string)
	status["status"] = "running"
	status["pub"] = d.GetPubString()
	statusBytes, err := json.Marshal(status)
	if err != nil {
		glg.Fatal(err)
	}
	d.router.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Write(statusBytes)
	})
	d.router.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		hashes, err := d.GetBC().GetBlockHashesHex()
		if err != nil {
			glg.Fatal(err)
		}
		height, err := d.GetBC().GetLatestHeight()
		if err != nil {
			glg.Fatal(err)
		}
		versionBytes, err := helpers.Serialize(NewVersion(GizoVersion, height, hashes))
		if err != nil {
			glg.Fatal(err)
		}
		w.Write(versionBytes)
	})

	err = d.discover.Forward(uint16(d.GetPort()), "gizo dispatcher node")
	if err != nil {
		log.Fatal(err)
	}
	if d.new {
		go d.Register()
	} else {
		res, err := d.centrum.Wake(d.GetPubString(), d.GetIP(), d.GetPort())
		if err != nil || res["status"].(string) != "success" {
			d.Register()
		}
	}
	d.jc = cache.NewJobCache(d.GetBC())
	fmt.Println(http.ListenAndServe(":"+strconv.FormatInt(int64(d.GetPort()), 10), d.router))
}

//SaveToken saves token got from centrum to db
func (d Dispatcher) SaveToken() {
	err := d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(NodeBucket))
		if err := b.Put([]byte("token"), []byte(d.centrum.GetToken())); err != nil {
			glg.Fatal(err)
		}
		return nil
	})
	if err != nil {
		glg.Fatal(err)
	}
}

//Register registers dispatcher with centrum
func (d Dispatcher) Register() {
	time.Sleep(time.Second * 1)
	if err := d.centrum.NewDisptcher(d.GetPubString(), d.GetIP(), int(d.GetPort())); err != nil {
		glg.Warn("Centrum: unable to connect to network")
		glg.Fatal("Centrum: " + err.Error())
	}
	d.SaveToken()
}

//GetDispatchersAndSync get's dispatchers from centrum and syncs with the node with the highest verison
func (d *Dispatcher) GetDispatchersAndSync() {
	//TODO: speed up
	time.Sleep(time.Second * 2)
	res, err := d.centrum.GetDispatchers()
	if err != nil {
		glg.Fatal(err)
	}
	blocks, err := d.GetBC().GetBlockHashesHex()
	if err != nil {
		glg.Fatal(err)
	}
	syncVersion := new(Version)
	syncPeer := new(client.Client)
	dispatchers, ok := res["dispatchers"]
	if !ok {
		glg.Warn(ErrNoDispatchers)
		if len(blocks) == 0 {
			d.GetBC().InitGenesisBlock(d.GetPubString())
		}
		return
	}
	for _, dispatcher := range dispatchers.([]string) {
		addr, err := ParseAddr(dispatcher)
		if err == nil && addr["pub"] != d.GetPubString() {
			var v Version
			wsURL := fmt.Sprintf("ws://%v:%v/wamp", addr["ip"], addr["port"])
			versionURL := fmt.Sprintf("http://%v:%v/version", addr["ip"], addr["port"])
			conn, err := client.ConnectNet(wsURL, client.Config{
				Realm: DISPATCHERREALM,
			})

			if err != nil {
				glg.Fatal(err)
			}
			d.AddPeer(addr["pub"], conn)
			go d.BlockSubscribe(conn)
			_, err = s.New().Get(versionURL).ReceiveSuccess(&v)
			if err != nil {
				glg.Fatal(err)
			}
			if syncVersion.GetHeight() <= v.GetHeight() {
				syncVersion = &v
				syncPeer = conn
			}
		}
	}
	glg.Info("Dispatcher: node sync in progress")
	if syncVersion.GetHeight() >= 0 {
		for _, hash := range syncVersion.GetBlocks() {
			if !funk.ContainsString(blocks, hash) {
				result, err := syncPeer.Call(context.Background(), BLOCKREQ, nil, wamp.List{hash}, nil, "")
				if err != nil {
					glg.Fatal("P2P: unable to sync blockchain", err)
				}
				blk, _ := wamp.AsString(result.Arguments[0])
				var block *core.Block
				err = helpers.Deserialize([]byte(blk), &block)
				if err != nil {
					glg.Fatal("P2P: unable to sync blockchain", err)
				}
				err = block.Export()
				if err != nil {
					glg.Fatal(err)
				}
				d.GetBC().AddBlock(block)
			}
		}
	}
}

//NewDispatcher initalizes dispatcher node
func NewDispatcher(port int) *Dispatcher {
	glg.Info("Creating Dispatcher Node")
	core.InitializeDataPath()
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	var bench benchmark.Engine
	var token string
	discover, err := upnp.Discover()
	if err != nil {
		glg.Fatal(err)
	}

	ip, err := discover.ExternalIP()
	if err != nil {
		log.Fatal(err)
	}

	centrum := NewCentrum()

	wampConfig := &nx_router.Config{
		RealmConfigs: []*nx_router.RealmConfig{
			&nx_router.RealmConfig{
				URI:           wamp.URI(DISPATCHERREALM),
				AnonymousAuth: true,
			},
			&nx_router.RealmConfig{
				URI:           wamp.URI(WORKERREALM),
				AnonymousAuth: true,
			},
		},
	}

	nxr, err := nx_router.NewRouter(wampConfig, nil)
	if err != nil {
		glg.Fatal(err)
	}

	dClient, err := client.ConnectLocal(nxr, client.Config{
		Realm: DISPATCHERREALM,
	})

	if err != nil {
		glg.Fatal(err)
	}

	wClient, err := client.ConnectLocal(nxr, client.Config{
		Realm: WORKERREALM,
	})

	if err != nil {
		glg.Fatal(err)
	}

	var dbFile string
	if os.Getenv("ENV") == "dev" {
		dbFile = path.Join(core.IndexPathDev, fmt.Sprintf(NodeDB, "dispatcher"))
	} else {
		dbFile = path.Join(core.IndexPathProd, fmt.Sprintf(NodeDB, "dispatcher"))
	}

	if helpers.FileExists(dbFile) {
		var priv, pub []byte
		glg.Info("Dispatcher: using existing keypair and benchmark")
		db, err := bolt.Open(dbFile, 0600, &bolt.Options{Timeout: time.Second * 2})
		if err != nil {
			glg.Fatal(err)
		}
		err = db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(NodeBucket))
			priv = b.Get([]byte("priv"))
			pub = b.Get([]byte("pub"))
			var bench *benchmark.Engine
			err := helpers.Deserialize(b.Get([]byte("benchmark")), &bench)
			if err != nil {
				glg.Fatal(err)
			}
			//!FIXME: check if token is defined
			token = string(b.Get([]byte("token")))
			return nil
		})
		if err != nil {
			glg.Fatal(err)
		}
		centrum.SetToken(token)
		bc := core.CreateBlockChain(hex.EncodeToString(pub))
		return &Dispatcher{
			ip:        ip,
			port:      uint(port),
			pub:       pub,
			priv:      priv,
			uptime:    time.Now().Unix(),
			jobPQ:     queue.NewJobPriorityQueue(),
			workers:   make(map[string]*WorkerInfo),
			workerPQ:  NewWorkerPriorityQueue(),
			peers:     make(map[string]*client.Client),
			bench:     bench,
			wamp:      nxr,
			rpc:       rpc.NewHTTPService(),
			router:    mux.NewRouter(),
			bc:        bc,
			db:        db,
			mu:        new(sync.Mutex),
			jobs:      make(map[string]job.Job),
			interrupt: interrupt,
			writeQ:    lane.NewQueue(),
			centrum:   centrum,
			discover:  discover,
			new:       false,
			dClient:   dClient,
			wClient:   wClient,
		}
	}

	priv, pub := crypt.GenKeysBytes()

	bench = benchmark.NewEngine()
	db, err := bolt.Open(dbFile, 0600, &bolt.Options{Timeout: time.Second * 2})
	if err != nil {
		glg.Fatal(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte(NodeBucket))
		if err != nil {
			glg.Fatal(err)
		}

		benchBytes, err := helpers.Serialize(bench)
		if err != nil {
			glg.Log("here")
			glg.Fatal(err)
		}
		if err = b.Put([]byte("benchmark"), benchBytes); err != nil {
			glg.Fatal(err)
		}

		if err = b.Put([]byte("priv"), priv); err != nil {
			glg.Fatal(err)
		}

		if err = b.Put([]byte("pub"), priv); err != nil {
			glg.Fatal(err)
		}
		return nil
	})

	if err != nil {
		glg.Fatal(err)
	}
	bc := core.CreateBlockChain(hex.EncodeToString(pub))
	return &Dispatcher{
		ip:        ip,
		port:      uint(port),
		pub:       pub,
		priv:      priv,
		uptime:    time.Now().Unix(),
		jobPQ:     queue.NewJobPriorityQueue(),
		workers:   make(map[string]*WorkerInfo),
		workerPQ:  NewWorkerPriorityQueue(),
		peers:     make(map[string]*client.Client),
		bench:     bench,
		wamp:      nxr,
		rpc:       rpc.NewHTTPService(),
		router:    mux.NewRouter(),
		bc:        bc,
		db:        db,
		mu:        new(sync.Mutex),
		jobs:      make(map[string]job.Job),
		interrupt: interrupt,
		writeQ:    lane.NewQueue(),
		centrum:   centrum,
		discover:  discover,
		new:       true,
		dClient:   dClient,
		wClient:   wClient,
	}
}
