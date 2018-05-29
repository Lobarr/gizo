package crypt

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"

	"github.com/kpango/glg"
)

//GenKeys returns private and public keys
func GenKeys() (private string, public string) {
	privKey, err := ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
	if err != nil {
		glg.Fatal(err)
	}
	privKeyBytes, err := x509.MarshalECPrivateKey(privKey)
	if err != nil {
		glg.Fatal(err)
	}
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(&privKey.PublicKey)
	if err != nil {
		glg.Fatal(err)
	}
	return hex.EncodeToString(privKeyBytes), hex.EncodeToString(pubKeyBytes)
}
