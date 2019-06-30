package p2p

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/boltdb/bolt"
	"github.com/gammazero/nexus/client"
	"github.com/gammazero/nexus/wamp"
	"github.com/gizo-network/gizo/core"
	"github.com/gizo-network/gizo/crypt"
	"github.com/gizo-network/gizo/helpers"
	"github.com/gizo-network/gizo/job/queue/qItem"
	"github.com/kpango/glg"
)

//Worker worker node
type Worker struct {
	Pub        []byte //public key of the node
	Dispatcher string
	shortlist  []string // array of dispatchers received from centrum
	priv       []byte   //private key of the node
	uptime     int64    //time since node has been up
	conn       *client.Client
	interrupt  chan os.Signal
	shutdown   chan struct{}
	busy       bool
	state      string
	item       qItem.Item
	logger     *glg.Glg
}

//JobTopic channel that dispatcher emits jobs to worker
func (w Worker) JobTopic() string {
	return fmt.Sprintf("worker.%v.job", w.GetPubString())
}

//ResultTopic channel that worker emits result to dispatcher
func (w Worker) ResultTopic() string {
	return fmt.Sprintf("worker.%v.result", w.GetPubString())
}

//CancelTopic channel that dispatcher emits cancel req to worker
func (w Worker) CancelTopic() string {
	return fmt.Sprintf("worker.%v.cancel", w.GetPubString())
}

//GetItem returns worker item
func (w Worker) GetItem() qItem.Item {
	return w.item
}

//SetItem sets item
func (w *Worker) SetItem(i qItem.Item) {
	w.item = i
}

//GetShortlist returs shortlists
func (w Worker) GetShortlist() []string {
	return w.shortlist
}

//SetShortlist sets shortlists
func (w *Worker) SetShortlist(s []string) {
	w.shortlist = s
}

//GetBusy returns if worker is busy
func (w Worker) GetBusy() bool {
	return w.busy
}

//SetBusy sets busy status
func (w *Worker) SetBusy(b bool) {
	w.busy = b
}

//GetState returns worker state
func (w Worker) GetState() string {
	return w.state
}

//SetState sets worker state
func (w *Worker) SetState(s string) {
	w.state = s
}

//GetPubByte returns public key as bytes
func (w Worker) GetPubByte() []byte {
	return w.Pub
}

//GetPubString returns public key as string
func (w Worker) GetPubString() string {
	return hex.EncodeToString(w.Pub)
}

//GetPrivByte returns private key as byte
func (w Worker) GetPrivByte() []byte {
	return w.priv
}

//GetPrivString returns private key as string
func (w Worker) GetPrivString() string {
	return hex.EncodeToString(w.priv)
}

//GetUptme returns node uptime
func (w Worker) GetUptme() int64 {
	return w.uptime
}

//GetDispatcher returns connected dispatcher
func (w Worker) GetDispatcher() string {
	return w.Dispatcher
}

//SetDispatcher sets connected dispatcher
func (w *Worker) SetDispatcher(d string) {
	w.Dispatcher = d
}

//GetUptimeString returns uptime as string
func (w Worker) GetUptimeString() string {
	return time.Unix(w.uptime, 0).Sub(time.Now()).String()
}

//Start runs a worker
func (w Worker) Start() {
	w.GetDispatchers()
	w.Connect()
	go w.WatchInterrupt()
	w.JobSubscription()
	w.CancelJobSubscription()
	select {
	case <-w.conn.Done():
		w.Connect()
	}
}

//CancelJobSubscription handles cancellation requests
func (w Worker) CancelJobSubscription() {
	handler := func(args wamp.List, kwargs, details wamp.Dict) {
		w.item.GetExec().Cancel()
	}
	w.conn.Subscribe(w.CancelTopic(), handler, nil)
}

//JobSubscription handles job execution
func (w Worker) JobSubscription() {
	handler := func(args wamp.List, kwargs, details wamp.Dict) {
		glg.Log("P2P: received job")
		if w.GetState() != LIVE {
			w.SetState(LIVE)
		}
		w.SetBusy(true)
		itemStr, _ := wamp.AsString(args[0])
		fmt.Println(itemStr) //TODO: remove logging
		var item *qItem.Item
		err := json.Unmarshal([]byte(itemStr), &item)
		if err != nil {
			return
		}
		w.SetItem(*item)
		exec := w.item.Job.Execute(w.item.GetExec(), w.GetDispatcher())
		w.item.SetExec(exec)
		execBytes, _ := json.Marshal(w.item.GetExec())
		w.conn.Publish(w.ResultTopic(), nil, wamp.List{string(execBytes)}, nil)

	}
	w.conn.Subscribe(w.JobTopic(), handler, nil)
}

//Disconnect disconnects from dispatcher
func (w Worker) Disconnect() {
	w.conn.Close()
}

//Connect connects to dispatcher
func (w *Worker) Connect() {
	for i, dispatcher := range w.GetShortlist() {
		addr, err := ParseAddr(dispatcher)
		if err == nil {
			url := fmt.Sprintf("ws://%v:%v/wamp", addr["ip"], addr["port"])
			if err = w.Dial(url); err == nil {
				w.logger.Log("Worker: connected to dispatcher")
				w.SetDispatcher(addr["pub"])
				return
			}
		}
		w.SetShortlist(append(w.GetShortlist()[:i], w.GetShortlist()[i+1:]...))
	}
	w.GetDispatchers()
}

//Dial attempts ws connections to dispatcher
func (w *Worker) Dial(url string) error {
	conn, err := client.ConnectNet(url, client.Config{
		Realm: WORKERREALM,
	})
	if err != nil {
		return err
	}
	res, err := conn.Call(context.Background(), WORKERCONNECT, nil, wamp.List{w.GetPubString()}, nil, "")
	if err != nil {
		return err
	}
	if res.Arguments[0].(string) == CONNFULL {
		return errors.New(CONNFULL)
	}
	w.conn = conn
	return nil
}

//WatchInterrupt watchs for interrupt
func (w Worker) WatchInterrupt() {
	select {
	case i := <-w.interrupt:
		glg.Warn("Worker: interrupt detected")
		switch i {
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			w.conn.Publish(WORKERDISCONNECT, nil, wamp.List{w.GetPubString()}, nil)
			w.conn.Close()
			break
		}
	}
}

//GetDispatchers connects to centrum to get dispatchers
func (w *Worker) GetDispatchers() {
	c := NewCentrum()
	res, err := c.GetDispatchers()
	if err != nil {
		glg.Fatal(err)
	}
	shortlist, ok := res["dispatchers"]
	if !ok {
		glg.Warn(ErrNoDispatchers)
		os.Exit(0)
	}
	w.SetShortlist(shortlist.([]string))
}

//NewWorker initializes worker node
func NewWorker() *Worker {
	core.InitializeDataPath()
	var priv, pub []byte
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	var dbFile string
	if os.Getenv("ENV") == "dev" {
		dbFile = path.Join(core.IndexPathDev, fmt.Sprintf(NodeDB, "worker"))
	} else {
		dbFile = path.Join(core.IndexPathProd, fmt.Sprintf(NodeDB, "worker"))
	}

	if helpers.FileExists(dbFile) {
		glg.Warn("Worker: using existing keypair")
		db, err := bolt.Open(dbFile, 0600, &bolt.Options{Timeout: time.Second * 2})
		if err != nil {
			glg.Fatal("Another worker using instance")
		}
		err = db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(NodeBucket))
			priv = b.Get([]byte("priv"))
			pub = b.Get([]byte("pub"))
			return nil
		})
		if err != nil {
			glg.Fatal(err)
		}
		return &Worker{
			Pub:       pub,
			priv:      priv,
			uptime:    time.Now().Unix(),
			interrupt: interrupt,
			state:     DOWN,
			logger:    helpers.Logger(),
		}
	}
	_priv, _pub := crypt.GenKeys()
	priv, err := hex.DecodeString(_priv)
	if err != nil {
		glg.Fatal(err)
	}
	pub, err = hex.DecodeString(_pub)
	if err != nil {
		glg.Fatal(err)
	}
	db, err := bolt.Open(dbFile, 0600, &bolt.Options{Timeout: time.Second * 2})
	if err != nil {
		glg.Fatal(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte(NodeBucket))
		if err != nil {
			glg.Fatal(err)
		}

		if err = b.Put([]byte("priv"), priv); err != nil {
			glg.Fatal(err)
		}

		if err = b.Put([]byte("pub"), pub); err != nil {
			glg.Fatal(err)
		}
		return nil
	})
	if err != nil {
		glg.Fatal(err)
	}
	return &Worker{
		Pub:       pub,
		priv:      priv,
		uptime:    time.Now().Unix(),
		interrupt: interrupt,
		state:     DOWN,
		logger:    helpers.Logger(),
	}
}
