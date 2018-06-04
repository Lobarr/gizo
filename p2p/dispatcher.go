package p2p

import (
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

	"github.com/Lobarr/lane"
	upnp "github.com/NebulousLabs/go-upnp"
	"github.com/hprose/hprose-golang/rpc"
	rpc_ws "github.com/hprose/hprose-golang/rpc/websocket"

	"github.com/gizo-network/gizo/core/difficulty"
	"github.com/gizo-network/gizo/core/merkletree"

	"github.com/gizo-network/gizo/job"

	"github.com/gizo-network/gizo/helpers"
	"github.com/gizo-network/gizo/job/queue"
	funk "github.com/thoas/go-funk"
	melody "gopkg.in/olahol/melody.v1"

	"github.com/boltdb/bolt"

	"github.com/kpango/glg"

	"github.com/gizo-network/gizo/benchmark"
	"github.com/gizo-network/gizo/cache"
	"github.com/gizo-network/gizo/core"
	"github.com/gizo-network/gizo/crypt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var (
	//ErrJobsFull occurs when jobs map is full
	ErrJobsFull = errors.New("Jobs map full")
)

//TODO: verify job hasn't been modified

//Dispatcher dispatcher node
type Dispatcher struct {
	IP        string
	Port      uint   //port
	Pub       []byte //public key of the node
	priv      []byte //private key of the node
	uptime    int64  //time since node has been up
	jobPQ     *queue.JobPriorityQueue
	workers   map[*melody.Session]*WorkerInfo //worker nodes in dispatcher's area
	peers     map[interface{}]*DispatcherInfo // neighbours of the node
	workerPQ  *WorkerPriorityQueue
	bench     benchmark.Engine //benchmark of node
	wWS       *melody.Melody   //workers ws server
	dWS       *melody.Melody   //dispatchers ws server
	rpcHTTP   *rpc.HTTPService // rpc servce
	rpcWS     *rpc_ws.WebSocketService
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
}

//GetJobs returns jobs held in memory to be written to the bc
func (d Dispatcher) GetJobs() map[string]job.Job {
	return d.jobs
}

//watches the queue of jobs to be written to the bc
func (d Dispatcher) watchWriteQ() {
	//TODO: write to bc if it's taking too long
	for {
		if d.GetWriteQ().Empty() == false {
			jobs := d.GetWriteQ().Dequeue()
			d.WriteJobs(jobs.(map[string]job.Job))
		}
	}
}

//WriteJobs writes jobs to the bc
func (d Dispatcher) WriteJobs(jobs map[string]job.Job) {
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
	blockBytes, err := helpers.Serialize(block)
	if err != nil {
		glg.Fatal(err)
	}
	bm, err := BlockMessage(blockBytes, d.GetPrivByte())
	if err != nil {
		glg.Fatal(err)
	}
	d.BroadcastPeers(bm)
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
func (d Dispatcher) GetAssignedWorker(hash string) *melody.Session {
	for key, val := range d.GetWorkers() {
		if val.GetJob().GetJob().GetHash() == hash {
			return key
		}
	}
	return nil
}

//GetIP returns ip address of node
func (d Dispatcher) GetIP() string {
	return d.IP
}

//SetIP sets ip address of node
func (d *Dispatcher) SetIP(ip string) {
	d.IP = ip
}

//GetPubByte returns public key of node as bytes
func (d Dispatcher) GetPubByte() []byte {
	return d.Pub
}

//GetPubString returns public key of node as bytes
func (d Dispatcher) GetPubString() string {
	return hex.EncodeToString(d.Pub)
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
func (d Dispatcher) GetWorkers() map[*melody.Session]*WorkerInfo {
	return d.workers
}

//GetWorker returns specified worker
func (d Dispatcher) GetWorker(s *melody.Session) *WorkerInfo {
	return d.GetWorkers()[s]
}

//SetWorker sets worker
func (d *Dispatcher) SetWorker(s *melody.Session, w *WorkerInfo) {
	d.GetWorkers()[s] = w
}

//GetPeers returns peers of a peer
func (d Dispatcher) GetPeers() map[interface{}]*DispatcherInfo {
	return d.peers
}

//GetPeersPubs returns public keys of peers
func (d Dispatcher) GetPeersPubs() []string {
	var temp []string
	for _, info := range d.GetPeers() {
		temp = append(temp, hex.EncodeToString(info.GetPub()))
	}
	return temp
}

//GetPeer returns dispatcher info of specific peer
func (d Dispatcher) GetPeer(n interface{}) *DispatcherInfo {
	return d.GetPeers()[n]
}

//AddPeer adds peers
func (d *Dispatcher) AddPeer(s interface{}, n *DispatcherInfo) {
	d.GetPeers()[s] = n
}

//GetBC returns blockchain object
func (d Dispatcher) GetBC() *core.BlockChain {
	return d.bc
}

//GetPort returns node port
func (d Dispatcher) GetPort() int {
	return int(d.Port)
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

//GetWWS returns workers ws server
func (d Dispatcher) GetWWS() *melody.Melody {
	return d.wWS
}

func (d *Dispatcher) setWWS(m *melody.Melody) {
	d.wWS = m
}

//GetDWS returns dispatchers ws server
func (d Dispatcher) GetDWS() *melody.Melody {
	return d.dWS
}

func (d *Dispatcher) setDWS(m *melody.Melody) {
	d.dWS = m
}

//GetRPCHTTP returns hprose rpc http server
func (d Dispatcher) GetRPCHTTP() *rpc.HTTPService {
	return d.rpcHTTP
}

func (d Dispatcher) setRPCHTTP(s *rpc.HTTPService) {
	d.rpcHTTP = s
}

//GetRPCWS returns hprose rpc ws server
func (d Dispatcher) GetRPCWS() *rpc_ws.WebSocketService {
	return d.rpcWS
}

func (d Dispatcher) setRPCWS(s *rpc_ws.WebSocketService) {
	d.rpcWS = s
}

//WorkerExists checks if worker is in node's standard area
func (d Dispatcher) WorkerExists(s *melody.Session) bool {
	_, ok := d.workers[s]
	return ok
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
			sm, err := ShutMessage(d.GetPrivByte())
			if err != nil {
				glg.Fatal(err)
			}
			d.BroadcastWorkers(sm)
			time.Sleep(time.Second * 3) // give neighbors and workers 3 seconds to disconnect
			os.Exit(0)
		case syscall.SIGQUIT:
			os.Exit(1)
		}
	}
}

//Start spins up services
func (d Dispatcher) Start() {
	verify, err := d.GetBC().Verify()
	if err != nil {
		glg.Fatal(err)
	}
	if !verify {
		glg.Fatal("Dispatcher: blockchain not verified")
	}
	go d.deployJobs()
	go d.watchWriteQ()
	go d.watchInterrupt()
	d.GetDispatchersAndSync()
	d.wWS.Upgrader.ReadBufferSize = 100000
	d.wWS.Upgrader.WriteBufferSize = 100000
	d.wWS.Config.MessageBufferSize = 100000
	d.wWS.Config.MaxMessageSize = 100000
	d.wWS.Upgrader.EnableCompression = true
	d.dWS.Upgrader.ReadBufferSize = 100000
	d.dWS.Upgrader.WriteBufferSize = 100000
	d.dWS.Config.MessageBufferSize = 100000
	d.dWS.Config.MaxMessageSize = 100000
	d.dWS.Upgrader.EnableCompression = true
	d.router.HandleFunc("/d", func(w http.ResponseWriter, r *http.Request) {
		d.dWS.HandleRequest(w, r)
	})
	d.router.HandleFunc("/w", func(w http.ResponseWriter, r *http.Request) {
		d.wWS.HandleRequest(w, r)
	})
	d.wPeerTalk()
	d.dPeerTalk()
	d.RPCHTTP()
	d.RPCWS()
	d.router.Handle("/rpc", d.GetRPCHTTP()).Methods("POST")
	d.router.Handle("/rpc-ws", d.GetRPCWS()).Methods("POST")
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
	syncVersion := new(Version)
	syncPeer := new(websocket.Conn)
	dispatchers, ok := res["dispatchers"]
	if !ok {
		glg.Warn(ErrNoDispatchers)
		// d.GetBC().InitGenesisBlock(d.GetPubString())
		return
	}
	for _, dispatcher := range dispatchers.([]string) {
		addr, err := ParseAddr(dispatcher)
		if err == nil && addr["pub"] != d.GetPubString() {
			var v Version
			wsURL := fmt.Sprintf("ws://%v:%v/d", addr["ip"], addr["port"])
			versionURL := fmt.Sprintf("http://%v:%v/rpc", addr["ip"], addr["port"])
			dailer := websocket.Dialer{
				Proxy:           http.ProxyFromEnvironment,
				ReadBufferSize:  10000,
				WriteBufferSize: 10000,
			}
			conn, _, err := dailer.Dial(wsURL, nil)
			if err != nil {
				continue
			}
			conn.EnableWriteCompression(true)
			pubBytes, err := hex.DecodeString(addr["pub"])
			if err != nil {
				glg.Fatal(err)
			}
			d.AddPeer(conn, NewDispatcherInfo(pubBytes, []string{}))
			go d.handleNodeConnect(conn)
			_, err = s.New().Get(versionURL).ReceiveSuccess(&v)
			if err != nil {
				glg.Fatal(err)
			}
			if syncVersion.GetHeight() < v.GetHeight() {
				syncVersion = &v
				syncPeer = conn
			}
		}
	}
	if syncVersion.GetHeight() != 0 {
		glg.Warn("Dispatcher: node sync in progress")
		blocks, err := d.GetBC().GetBlockHashesHex()
		if err != nil {
			glg.Fatal(err)
		}
		for _, hash := range syncVersion.GetBlocks() {
			if !funk.ContainsString(blocks, hash) {
				hashBytes, err := hex.DecodeString(hash)
				if err != nil {
					glg.Fatal(err)
				}
				brm, err := BlockReqMessage(hashBytes, d.GetPrivByte())
				if err != nil {
					glg.Fatal(err)
				}
				syncPeer.WriteMessage(websocket.BinaryMessage, brm)
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
	var priv, pub []byte
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

	var dbFile string
	if os.Getenv("ENV") == "dev" {
		dbFile = path.Join(core.IndexPathDev, fmt.Sprintf(NodeDB, "dispatcher"))
	} else {
		dbFile = path.Join(core.IndexPathProd, fmt.Sprintf(NodeDB, "dispatcher"))
	}

	if helpers.FileExists(dbFile) {
		glg.Warn("Dispatcher: using existing keypair and benchmark")
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
		// bc.InitGenesisBlock(hex.EncodeToString(pub))
		jc := cache.NewJobCache(bc)
		return &Dispatcher{
			IP:        ip,
			Pub:       pub,
			priv:      priv,
			Port:      uint(port),
			uptime:    time.Now().Unix(),
			bench:     bench,
			jobPQ:     queue.NewJobPriorityQueue(),
			workers:   make(map[*melody.Session]*WorkerInfo),
			workerPQ:  NewWorkerPriorityQueue(),
			peers:     make(map[interface{}]*DispatcherInfo),
			jc:        jc,
			bc:        bc,
			db:        db,
			router:    mux.NewRouter(),
			wWS:       melody.New(),
			dWS:       melody.New(),
			rpcHTTP:   rpc.NewHTTPService(),
			rpcWS:     rpc_ws.NewWebSocketService(),
			mu:        new(sync.Mutex),
			interrupt: interrupt,
			writeQ:    lane.NewQueue(),
			centrum:   centrum,
			discover:  discover,
			new:       false,
			jobs:      make(map[string]job.Job),
		}
	}

	_priv, _pub := crypt.GenKeys()
	priv, err = hex.DecodeString(_priv)
	if err != nil {
		glg.Fatal(err)
	}
	pub, err = hex.DecodeString(_pub)
	if err != nil {
		glg.Fatal(err)
	}
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
	jc := cache.NewJobCache(bc)
	return &Dispatcher{
		IP:        ip,
		Pub:       pub,
		priv:      priv,
		Port:      uint(port),
		uptime:    time.Now().Unix(),
		bench:     bench,
		jobPQ:     queue.NewJobPriorityQueue(),
		workers:   make(map[*melody.Session]*WorkerInfo),
		workerPQ:  NewWorkerPriorityQueue(),
		peers:     make(map[interface{}]*DispatcherInfo),
		jc:        jc,
		bc:        bc,
		db:        db,
		router:    mux.NewRouter(),
		wWS:       melody.New(),
		dWS:       melody.New(),
		rpcHTTP:   rpc.NewHTTPService(),
		rpcWS:     rpc_ws.NewWebSocketService(),
		mu:        new(sync.Mutex),
		interrupt: interrupt,
		writeQ:    lane.NewQueue(),
		centrum:   centrum,
		discover:  discover,
		new:       true,
		jobs:      make(map[string]job.Job),
	}
}
