package p2p

import (
	"github.com/gizo-network/gizo/helpers"
)

const (
	HELLO            = "HELLO"
	INVALIDMESSAGE   = "INVALIDMESSAGE" // invalid message
	CONNFULL         = "CONNFULL"       // max workers reached
	JOB              = "JOB"
	INVALIDSIGNATURE = "JOB"
	RESULT           = "RESULT"
	SHUT             = "SHUT"
	SHUTACK          = "SHUTACK"
	BLOCK            = "BLOCK"
	BLOCKREQ         = "BLOCKREQ"
	BLOCKRES         = "BLOCKRES"
	PEERCONNECT      = "PEERCONNECT"
	PEERDISCONNECT   = "PEERDISCONNECT"
	CANCEL           = "CANCEL"
)

//HelloMessage initiates adjacency connection between nodes
func HelloMessage(payload []byte) ([]byte, error) {
	return helpers.Serialize(NewPeerMessage(HELLO, payload, nil))
}

//InvalidMessage when a node receives an unknown message
func InvalidMessage() ([]byte, error) {
	return helpers.Serialize(NewPeerMessage(INVALIDMESSAGE, nil, nil))
}

func ConnFullMessage() ([]byte, error) {
	return helpers.Serialize(NewPeerMessage(CONNFULL, nil, nil))
}

func JobMessage(payload, priv []byte) ([]byte, error) {
	return helpers.Serialize(NewPeerMessage(JOB, payload, priv))
}

func InvalidSignature() ([]byte, error) {
	return helpers.Serialize(NewPeerMessage(INVALIDMESSAGE, nil, nil))
}

func ResultMessage(payload, priv []byte) ([]byte, error) {
	return helpers.Serialize(NewPeerMessage(RESULT, payload, priv))
}

func ShutMessage(priv []byte) ([]byte, error) {
	return helpers.Serialize(NewPeerMessage(SHUT, nil, priv))
}

func ShutAckMessage(priv []byte) ([]byte, error) {
	return helpers.Serialize(NewPeerMessage(SHUTACK, nil, priv))
}

func BlockMessage(payload, priv []byte) ([]byte, error) {
	return helpers.Serialize(NewPeerMessage(BLOCK, payload, priv))
}

func BlockReqMessage(payload, priv []byte) ([]byte, error) {
	return helpers.Serialize(NewPeerMessage(BLOCKREQ, payload, priv))
}

func BlockResMessage(payload, priv []byte) ([]byte, error) {
	return helpers.Serialize(NewPeerMessage(BLOCKRES, payload, priv))
}

func PeerConnectMessage(payload, priv []byte) ([]byte, error) {
	return helpers.Serialize(NewPeerMessage(PEERCONNECT, payload, priv))
}

func PeerDisconnectMessage(payload, priv []byte) ([]byte, error) {
	return helpers.Serialize(NewPeerMessage(PEERDISCONNECT, payload, priv))
}

func CancelMessage(priv []byte) ([]byte, error) {
	return helpers.Serialize(NewPeerMessage(CANCEL, nil, priv))
}
