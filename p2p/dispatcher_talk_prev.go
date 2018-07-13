package p2p

import (
	"bytes"
	"encoding/hex"

	"github.com/gizo-network/gizo/core"
	"github.com/gizo-network/gizo/helpers"
	"github.com/gizo-network/gizo/job"
	"github.com/gorilla/websocket"
	"github.com/kpango/glg"
	funk "github.com/thoas/go-funk"
	melody "gopkg.in/olahol/melody.v1"
)

//BroadcastWorkers sends message to all workers
func (d Dispatcher) BroadcastWorkers(m []byte) {
	for s := range d.GetWorkers() {
		s.Write(m)
	}
}

//BroadcastPeers sends message to all peers
func (d Dispatcher) BroadcastPeers(m []byte) {
	for peer := range d.GetPeers() {
		switch n := peer.(type) {
		case *melody.Session:
			n.Write(m)
			break
		case *websocket.Conn:
			n.WriteMessage(websocket.BinaryMessage, m)
			break
		}
	}
}

//MulticastPeers sends message to specified peers
func (d Dispatcher) MulticastPeers(m []byte, peers []string) {
	for peer, info := range d.GetPeers() {
		if funk.ContainsString(peers, hex.EncodeToString(info.GetPub())) {
			switch n := peer.(type) {
			case *melody.Session:
				n.Write(m)
				break
			case *websocket.Conn:
				n.WriteMessage(websocket.BinaryMessage, m)
				break
			}
		}
	}
}

//assigns next job in the queue to the next available worker
func (d *Dispatcher) deployJobs() {
	for {
		if d.GetWorkerPQ().getPQ().Empty() == false {
			if d.GetJobPQ().GetPQ().Empty() == false {
				d.mu.Lock()
				w := d.GetWorkerPQ().Pop()
				if !d.GetWorker(w).GetShut() {
					j := d.GetJobPQ().Pop()
					if j.GetExec().GetStatus() != job.CANCELLED {
						j.GetExec().SetBy(d.GetWorker(w).GetPub())
						d.GetWorker(w).Assign(&j)
						glg.Info("P2P: dispatched job")
						jBytes, err := helpers.Serialize(j)
						if err != nil {
							glg.Fatal(err)
						}
						jm, err := JobMessage(jBytes, d.GetPrivByte())
						if err != nil {
							glg.Fatal(err)
						}
						w.Write(jm)
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

//communication between dispatcher and worker
func (d Dispatcher) wPeerTalk() {
	//TODO: handle unexpected error
	// d.wWS.HandleDisconnect(func(s *melody.Session) {
	// 	d.mu.Lock()
	// 	glg.Log("Dispatcher: worker disconnected")
	// 	if d.GetWorker(s).GetJob() != nil {
	// 		d.GetJobPQ().PushItem(*d.GetWorker(s).GetJob(), job.BOOST)
	// 	}
	// 	d.GetWorker(s).SetShut(true)
	// 	d.mu.Unlock()
	// })
	//TODO: handle panic
	d.wWS.HandleError(func(s *melody.Session, err error) {
		d.mu.Lock()
		glg.Log("Dispatcher: worker disconnected")
		if d.GetWorker(s).GetJob() != nil {
			d.GetJobPQ().PushItem(*d.GetWorker(s).GetJob(), job.BOOST)
		}
		d.GetWorker(s).SetShut(true)
		d.mu.Unlock()
	})
	d.wWS.HandleMessageBinary(func(s *melody.Session, message []byte) {
		var m PeerMessage
		err := helpers.Deserialize(message, &m)
		if err != nil {
			glg.Fatal(err)
		}
		switch m.GetMessage() {
		case HELLO:
			d.mu.Lock()
			if len(d.GetWorkers()) < MaxWorkers {
				glg.Info("Dispatcher: worker connected")
				d.SetWorker(s, NewWorkerInfo(hex.EncodeToString(m.GetPayload())))
				hm, err := HelloMessage(d.GetPubByte())
				if err != nil {
					glg.Fatal(err)
				}
				s.Write(hm)
				d.centrum.ConnectWorker()
				d.GetWorkerPQ().Push(s, 0)
			} else {
				cfm, err := ConnFullMessage()
				if err != nil {
					glg.Fatal(err)
				}
				s.Write(cfm)
			}
			d.mu.Unlock()
			break
		case RESULT:
			d.mu.Lock()
			verify, err := m.VerifySignature(d.GetWorker(s).GetPub())
			if err != nil {
				glg.Fatal(err)
			}
			if verify {
				glg.Info("P2P: received result")
				var exec *job.Exec
				err := helpers.Deserialize(m.GetPayload(), &exec)
				if err != nil {
					glg.Fatal(err)
				}
				d.GetWorker(s).GetJob().SetExec(exec)
				d.GetWorker(s).GetJob().ResultsChan() <- *d.GetWorker(s).GetJob()
				j := d.GetWorker(s).GetJob().GetJob()
				d.GetWorker(s).SetJob(nil)
				j.AddExec(*exec)
				//TODO: send to requester's message broker
				d.AddJob(j)
			} else {
				d.GetJobPQ().PushItem(*d.GetWorker(s).GetJob(), job.HIGH)
			}
			if !d.GetWorker(s).GetShut() {
				d.GetWorkerPQ().Push(s, 0)
			}
			d.mu.Unlock()
			break
		case SHUT:
			d.mu.Lock()
			d.GetWorker(s).SetShut(true)
			sam, err := ShutAckMessage(d.GetPrivByte())
			if err != nil {
				glg.Fatal(err)
			}
			s.Write(sam)
			d.centrum.DisconnectWorker()
			d.mu.Unlock()
			break
		default:
			im, err := InvalidMessage()
			if err != nil {
				glg.Fatal(err)
			}
			s.Write(im)
			break
		}
	})
}

// communication between dispatcher and dispatcher
func (d Dispatcher) dPeerTalk() {
	d.dWS.HandleDisconnect(func(s *melody.Session) {
		d.mu.Lock()
		info := d.GetPeer(s)
		if info != nil {
			glg.Info("Dispatcher: peer disconnected")
			pdm, err := PeerDisconnectMessage(info.GetPub(), d.GetPrivByte())
			if err != nil {
				glg.Fatal(err)
			}
			d.BroadcastPeers(pdm)
			delete(d.GetPeers(), s)
		}
		d.mu.Unlock()
	})
	d.dWS.HandleMessageBinary(func(s *melody.Session, message []byte) {
		var m PeerMessage
		err := helpers.Deserialize(message, &m)
		if err != nil {
			glg.Fatal(err)
		}
		switch m.GetMessage() {
		case HELLO:
			d.mu.Lock()
			var info DispatcherInfo
			err := helpers.Deserialize(m.GetPayload(), &info)
			if err != nil {
				glg.Fatal(err)
			}
			d.AddPeer(s, &DispatcherInfo{Pub: info.GetPub(), Peers: info.GetPeers()})
			diBytes, err := helpers.Serialize(NewDispatcherInfo(d.GetPubByte(), d.GetPeersPubs()))
			if err != nil {
				glg.Fatal(err)
			}
			hm, err := HelloMessage(diBytes)
			if err != nil {
				glg.Fatal(err)
			}
			s.Write(hm)
			d.mu.Unlock()
			break
		case BLOCK:
			d.mu.Lock()
			verify, err := m.VerifySignature(hex.EncodeToString(d.GetPeer(s).GetPub()))
			if err != nil {
				glg.Fatal(err)
			}
			if verify {
				var b *core.Block
				err := helpers.Deserialize(m.GetPayload(), &b)
				if err != nil {
					glg.Fatal(err)
				}
				err = b.Export()
				if err != nil {
					glg.Fatal(err)
				}
				d.GetBC().AddBlock(b)
				var peerToRecv []string
				for _, info := range d.GetPeers() {
					//! avoids broadcast storms by not sending block back to sender and to neigbhours that are not directly connected to sender
					if !funk.ContainsString(info.GetPeers(), hex.EncodeToString(d.GetPeer(s).GetPub())) && bytes.Compare(info.GetPub(), d.GetPeer(s).GetPub()) != 0 {
						peerToRecv = append(peerToRecv, hex.EncodeToString(info.GetPub()))
					}
				}
				bm, err := BlockMessage(m.GetPayload(), d.GetPrivByte())
				if err != nil {
					glg.Fatal(err)
				}
				d.MulticastPeers(bm, peerToRecv)
			}
			d.mu.Unlock()
			break
		case BLOCKREQ:
			d.mu.Lock()
			verify, err := m.VerifySignature(hex.EncodeToString(d.GetPeer(s).GetPub()))
			if err != nil {
				glg.Fatal(err)
			}
			if verify {
				blockinfo, _ := d.GetBC().GetBlockInfo(hex.EncodeToString(m.GetPayload()))
				biBytes, err := helpers.Serialize(blockinfo.GetBlock())
				if err != nil {
					glg.Fatal(err)
				}
				brs, err := BlockResMessage(biBytes, d.GetPrivByte())
				s.Write(brs)
			}
			d.mu.Unlock()
			break
		case PEERCONNECT:
			d.mu.Lock()
			verify, err := m.VerifySignature(hex.EncodeToString(d.GetPeer(s).GetPub()))
			if err != nil {
				glg.Fatal(err)
			}
			if verify {
				d.GetPeer(s).AddPeer(hex.EncodeToString(m.GetPayload()))
			}
			d.mu.Unlock()
			break
		case PEERDISCONNECT:
			d.mu.Lock()
			verify, err := m.VerifySignature(hex.EncodeToString(d.GetPeer(s).GetPub()))
			if err != nil {
				glg.Fatal(err)
			}
			if verify {
				peers := d.GetPeer(s).GetPeers()
				for i, peer := range peers {
					if peer == hex.EncodeToString(m.GetPayload()) {
						d.GetPeer(s).SetPeers(append(peers[:i], peers[i+1:]...))
						break
					}
				}
			}
			d.mu.Unlock()
			break
		default:
			im, err := InvalidMessage()
			if err != nil {
				glg.Fatal(err)
			}
			s.Write(im)
			break
		}
	})
}

// communication between dispatcher and dispatcher
func (d Dispatcher) handleNodeConnect(conn *websocket.Conn) {
	diBytes, err := helpers.Serialize(NewDispatcherInfo(d.GetPubByte(), d.GetPeersPubs()))
	if err != nil {
		glg.Fatal(err)
	}
	hm, err := HelloMessage(diBytes)
	if err != nil {
		glg.Fatal(err)
	}
	conn.WriteMessage(websocket.BinaryMessage, hm)
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			//TODO: handle syncer disconnect - use next best version
			d.mu.Lock()
			glg.Info("Dispatcher: peer disconnected")
			info := d.GetPeer(conn)
			pdm, err := PeerDisconnectMessage(info.GetPub(), d.GetPrivByte())
			if err != nil {
				glg.Fatal(err)
			}
			d.BroadcastPeers(pdm)
			delete(d.GetPeers(), conn)
			d.mu.Unlock()
			break
		}
		var m PeerMessage
		err = helpers.Deserialize(message, &m)
		if err != nil {
			glg.Fatal(err)
		}
		switch m.GetMessage() {
		case HELLO:
			d.mu.Lock()
			var peerInfo DispatcherInfo
			err := helpers.Deserialize(m.GetPayload(), &peerInfo)
			if err != nil {
				glg.Fatal(err)
			}
			if bytes.Compare(d.GetPeer(conn).GetPub(), peerInfo.GetPub()) == 0 {
				d.GetPeer(conn).SetPeers(peerInfo.GetPeers())
			} else {
				delete(d.GetPeers(), conn)
				conn.Close()
			}
			d.mu.Unlock()
			break
		case BLOCK:
			d.mu.Lock()
			verify, err := m.VerifySignature(hex.EncodeToString(d.GetPeer(conn).GetPub()))
			if err != nil {
				glg.Fatal(err)
			}
			if verify {
				var b *core.Block
				err := helpers.Deserialize(m.GetPayload(), &b)
				if err != nil {
					glg.Fatal(err)
				}
				err = b.Export()
				if err != nil {
					glg.Fatal(err)
				}
				d.GetBC().AddBlock(b)
				var peerToRecv []string
				for _, info := range d.GetPeers() {
					//! avoids broadcast storms by not sending block back to sender and to neigbhours that are not directly connected to sender
					if !funk.ContainsString(info.GetPeers(), hex.EncodeToString(d.GetPeer(conn).GetPub())) && bytes.Compare(info.GetPub(), d.GetPeer(conn).GetPub()) != 0 {
						peerToRecv = append(peerToRecv, hex.EncodeToString(info.GetPub()))
					}
				}
				bm, err := BlockMessage(m.GetPayload(), d.GetPrivByte())
				d.MulticastPeers(bm, peerToRecv)
			}
			d.mu.Unlock()
			break
		case BLOCKRES:
			verify, err := m.VerifySignature(hex.EncodeToString(d.GetPeer(conn).GetPub()))
			if err != nil {
				glg.Fatal(err)
			}
			if verify {
				var b *core.Block
				err := helpers.Deserialize(m.GetPayload(), &b)
				if err != nil {
					glg.Fatal(err)
				}
				err = b.Export()
				if err != nil {
					glg.Fatal(err)
				}
				d.GetBC().AddBlock(b)
			}
			break
		case PEERCONNECT:
			d.mu.Lock()
			verify, err := m.VerifySignature(hex.EncodeToString(d.GetPeer(conn).GetPub()))
			if err != nil {
				glg.Fatal(err)
			}
			if verify {
				d.GetPeer(conn).AddPeer(hex.EncodeToString(m.GetPayload()))
			}
			d.mu.Unlock()
			break
		case PEERDISCONNECT:
			d.mu.Lock()
			verify, err := m.VerifySignature(hex.EncodeToString(d.GetPeer(conn).GetPub()))
			if err != nil {
				glg.Fatal(err)
			}
			if verify {
				peers := d.GetPeer(conn).GetPeers()
				for i, peer := range peers {
					if peer == hex.EncodeToString(m.GetPayload()) {
						d.GetPeer(conn).SetPeers(append(peers[:i], peers[i+1:]...))
					}
				}
			}
			d.mu.Unlock()
			break
		default:
			im, err := InvalidMessage()
			if err != nil {
				glg.Fatal(err)
			}
			conn.WriteMessage(websocket.BinaryMessage, im)
			break
		}
	}
}
