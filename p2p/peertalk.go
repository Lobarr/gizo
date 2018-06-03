package p2p

import (
	"github.com/gizo-network/gizo/helpers"
)

const (
	//HELLO initiates adjacency between nodes
	HELLO = "HELLO"
	//INVALIDMESSAGE when a node receives an unknown message
	INVALIDMESSAGE = "INVALIDMESSAGE" // invalid message
	//CONNFULL when the max number of workers is reached
	CONNFULL = "CONNFULL"
	//JOB when a job is sent or received
	JOB = "JOB"
	//INVALIDSIGNATURE when a message has an invalid signature
	INVALIDSIGNATURE = "INVALIDSIGNATURE"
	//RESULT job result received from worker
	RESULT = "RESULT"
	//SHUT when a node shuts off
	SHUT = "SHUT"
	//SHUTACK acknoledgement that node is shutting off
	SHUTACK = "SHUTACK"
	//BLOCK when a block is sent or received
	BLOCK = "BLOCK"
	//BLOCKREQ request for block
	BLOCKREQ = "BLOCKREQ"
	//BLOCKRES response for block
	BLOCKRES = "BLOCKRES"
	//PEERCONNECT broadcast that a peer connected to immediate neighbours
	PEERCONNECT = "PEERCONNECT"
	//PEERDISCONNECT broadcast that a peer disconnected to immediate neighbours
	PEERDISCONNECT = "PEERDISCONNECT"
	//CANCEL sent for cancelling exec
	CANCEL = "CANCEL"
)

//HelloMessage initiates adjacency connection between nodes
func HelloMessage(payload []byte) ([]byte, error) {
	return helpers.Serialize(NewPeerMessage(HELLO, payload, nil))
}

//InvalidMessage when a node receives an unknown message
func InvalidMessage() ([]byte, error) {
	return helpers.Serialize(NewPeerMessage(INVALIDMESSAGE, nil, nil))
}

//ConnFullMessage when the max number of workers is reached
func ConnFullMessage() ([]byte, error) {
	return helpers.Serialize(NewPeerMessage(CONNFULL, nil, nil))
}

//JobMessage when a job is sent or received
func JobMessage(payload, priv []byte) ([]byte, error) {
	return helpers.Serialize(NewPeerMessage(JOB, payload, priv))
}

//InvalidSignature when a message has an invalid signature
func InvalidSignature() ([]byte, error) {
	return helpers.Serialize(NewPeerMessage(INVALIDMESSAGE, nil, nil))
}

//ResultMessage job result received from worker
func ResultMessage(payload, priv []byte) ([]byte, error) {
	return helpers.Serialize(NewPeerMessage(RESULT, payload, priv))
}

//ShutMessage when a node shuts off
func ShutMessage(priv []byte) ([]byte, error) {
	return helpers.Serialize(NewPeerMessage(SHUT, nil, priv))
}

//ShutAckMessage acknoledgement that node is shutting off
func ShutAckMessage(priv []byte) ([]byte, error) {
	return helpers.Serialize(NewPeerMessage(SHUTACK, nil, priv))
}

//BlockMessage when a block is sent or received
func BlockMessage(payload, priv []byte) ([]byte, error) {
	return helpers.Serialize(NewPeerMessage(BLOCK, payload, priv))
}

//BlockReqMessage request for block
func BlockReqMessage(payload, priv []byte) ([]byte, error) {
	return helpers.Serialize(NewPeerMessage(BLOCKREQ, payload, priv))
}

//BlockResMessage response for block
func BlockResMessage(payload, priv []byte) ([]byte, error) {
	return helpers.Serialize(NewPeerMessage(BLOCKRES, payload, priv))
}

//PeerConnectMessage broadcast that a peer connected to immediate neighbours
func PeerConnectMessage(payload, priv []byte) ([]byte, error) {
	return helpers.Serialize(NewPeerMessage(PEERCONNECT, payload, priv))
}

//PeerDisconnectMessage broadcast that a peer disconnected to immediate neighbours
func PeerDisconnectMessage(payload, priv []byte) ([]byte, error) {
	return helpers.Serialize(NewPeerMessage(PEERDISCONNECT, payload, priv))
}

//CancelMessage sent for cancelling exec
func CancelMessage(priv []byte) ([]byte, error) {
	return helpers.Serialize(NewPeerMessage(CANCEL, nil, priv))
}
