package crypt

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
)

//GenKeys returns private and public keys
func GenKeys() (private string, public string) {
	privKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	privKeyBytes, _ := x509.MarshalECPrivateKey(privKey)
	pubKeyBytes, _ := x509.MarshalPKIXPublicKey(&privKey.PublicKey)
	return hex.EncodeToString(privKeyBytes), hex.EncodeToString(pubKeyBytes)
}

//GenKeysBytes returns private and public keys as bytes
func GenKeysBytes() (private []byte, public []byte) {
	privKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	privKeyBytes, _ := x509.MarshalECPrivateKey(privKey)
	pubKeyBytes, _ := x509.MarshalPKIXPublicKey(&privKey.PublicKey)
	return privKeyBytes, pubKeyBytes
}
