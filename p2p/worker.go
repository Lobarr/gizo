package p2p

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/boltdb/bolt"
	"github.com/gizo-network/gizo/helpers"

	"github.com/gizo-network/gizo/core"
	"github.com/gizo-network/gizo/crypt"
	"github.com/gizo-network/gizo/job/queue/qItem"
	"github.com/gorilla/websocket"
	"github.com/kpango/glg"
)

//Worker worker node
type Worker struct {
	Pub        []byte //public key of the node
	Dispatcher string
	shortlist  []string // array of dispatchers received from centrum
	priv       []byte   //private key of the node
	uptime     int64    //time since node has been up
	conn       *websocket.Conn
	interrupt  chan os.Signal
	shutdown   chan struct{}
	busy       bool
	state      string
	item       qItem.Item
	logger     *glg.Glg
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

//Start starts running
func (w *Worker) Start() {
	//TODO: implemented cancel and getstatus
	w.GetDispatchers()
	w.Connect()
	go w.WatchInterrupt()
	hm, err := HelloMessage(w.GetPubByte())
	if err != nil {
		//FIXME:
	}
	w.conn.WriteMessage(websocket.BinaryMessage, hm)
	for {
		_, message, err := w.conn.ReadMessage()
		if err != nil {
			w.Connect()
		}
		var m PeerMessage
		err = helpers.Deserialize(message, &m)
		if err != nil {
			//FIXME:
		}
		switch m.GetMessage() {
		case CONNFULL:
			w.Disconnect()
			w.Connect()
			break
		case HELLO:
			if w.GetDispatcher() != hex.EncodeToString(m.GetPayload()) {
				w.Disconnect()
				w.Connect()
			}
			w.SetState(INIT)
			glg.Info("P2P: connected to dispatcher")
			break
		case JOB:
			glg.Info("P2P: job received")
			if w.GetState() != LIVE {
				w.SetState(LIVE)
			}
			w.SetBusy(true)
			verify, err := m.VerifySignature(w.GetDispatcher())
			if err != nil {
				//FIXME:
			}
			if verify {
				var item qItem.Item
				err = helpers.Deserialize(m.GetPayload(), &item)
				w.SetItem(item)
				exec := w.item.Job.Execute(w.item.GetExec(), w.GetDispatcher())
				w.item.SetExec(exec)
				resultBytes, err := helpers.Serialize(w.item.GetExec())
				if err != nil {
					//FIXME: handle error
				}
				rm, err := ResultMessage(resultBytes, w.GetPrivByte())
				if err != nil {
					//FIXME:
				}
				w.conn.WriteMessage(websocket.BinaryMessage, rm)
			} else {
				is, err := InvalidSignature()
				if err != nil {
					//FIXME:
				}
				w.conn.WriteMessage(websocket.BinaryMessage, is)
				w.Disconnect()
			}
			w.SetBusy(false)
			break
		case CANCEL:
			glg.Info("P2P: job cancelled")
			verify, err := m.VerifySignature(w.GetDispatcher())
			if err != nil {
				//FIXME:
			}
			if verify {
				w.item.GetExec().Cancel()
			} else {
				is, err := InvalidSignature()
				if err != nil {
					//FIXME:
				}
				w.conn.WriteMessage(websocket.BinaryMessage, is)
			}
			break
		case SHUT:
			w.Connect()
			break
		case SHUTACK:
			for {
				if w.GetBusy() {
					continue
				} else {
					break
				}
			} // wait until worker not busy
			w.Disconnect()
			w.SetState(DOWN)
			glg.Info("Worker: graceful shutdown")
			os.Exit(0)
		default:
			w.Disconnect() //look for new dispatcher
			break
		}
	}
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
			url := fmt.Sprintf("ws://%v:%v/w", addr["ip"], addr["port"])
			if err = w.Dial(url); err == nil {
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
	dailer := websocket.Dialer{
		Proxy:           http.ProxyFromEnvironment,
		ReadBufferSize:  10000,
		WriteBufferSize: 10000,
	}
	conn, _, err := dailer.Dial(url, nil)
	if err != nil {
		return err
	}
	conn.EnableWriteCompression(true)
	w.conn = conn
	return nil
}

//WatchInterrupt watchs for interrupt
func (w Worker) WatchInterrupt() {
	select {
	case i := <-w.interrupt:
		glg.Warn("Worker: interrupt detected")
		switch i {
		case syscall.SIGINT, syscall.SIGTERM:
			sm, err := ShutMessage(w.GetPrivByte())
			if err != nil {
				//FIXME:
			}
			w.conn.WriteMessage(websocket.BinaryMessage, sm)
			break
		case syscall.SIGQUIT:
			os.Exit(1)
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
			glg.Fatal(err)
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
