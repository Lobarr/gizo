package p2p

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"math/big"
)

//PeerMessage message sent between nodes
type PeerMessage struct {
	Message   string
	Payload   []byte
	Signature [][]byte
}

//NewPeerMessage initializes peer message
func NewPeerMessage(message string, payload []byte, priv []byte) PeerMessage {
	pm := PeerMessage{Message: message, Payload: payload}
	if priv != nil {
		pm.sign(priv)
	}
	return pm
}

//GetMessage return message
func (m PeerMessage) GetMessage() string {
	return m.Message
}

//GetPayload returns payload
func (m PeerMessage) GetPayload() []byte {
	return m.Payload
}

//GetSignature returns signature
func (m PeerMessage) GetSignature() [][]byte {
	return m.Signature
}

//SetMessage sets message
func (m *PeerMessage) SetMessage(message string) {
	m.Message = message
}

//SetPayload sets payload
func (m *PeerMessage) SetPayload(payload []byte) {
	m.Payload = payload
}

func (m *PeerMessage) setSignature(sig [][]byte) {
	m.Signature = sig
}

func (m *PeerMessage) sign(priv []byte) error {
	hash := sha256.Sum256(bytes.Join(
		[][]byte{
			[]byte(m.GetMessage()),
			m.GetPayload(),
		},
		[]byte{},
	))
	privateKey, _ := x509.ParseECPrivateKey(priv)
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		return err
	}
	var temp [][]byte
	temp = append(temp, r.Bytes(), s.Bytes())
	m.setSignature(temp)
	return nil
}

//VerifySignature verifies message signature
func (m *PeerMessage) VerifySignature(pub string) (bool, error) {
	pubBytes, err := hex.DecodeString(pub)
	if err != nil {
		return false, err
	}
	var r big.Int
	var s big.Int
	r.SetBytes(m.GetSignature()[0])
	s.SetBytes(m.GetSignature()[1])

	publicKey, _ := x509.ParsePKIXPublicKey(pubBytes)
	hash := sha256.Sum256(bytes.Join(
		[][]byte{
			[]byte(m.GetMessage()),
			m.GetPayload(),
		},
		[]byte{},
	))
	switch pubConv := publicKey.(type) {
	case *ecdsa.PublicKey:
		return ecdsa.Verify(pubConv, hash[:], &r, &s), nil
	default:
		return false, nil
	}
}
