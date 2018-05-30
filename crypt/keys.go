package crypt

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"

	"github.com/gizo-network/gizo/helpers"
)

//GenKeys returns private and public keys
func GenKeys() (private string, public string) {
	logger := helpers.Logger()
	privKey, err := ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
	if err != nil {
		logger.Fatal(err)
	}
	privKeyBytes, err := x509.MarshalECPrivateKey(privKey)
	if err != nil {
		logger.Fatal(err)
	}
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(&privKey.PublicKey)
	if err != nil {
		logger.Fatal(err)
	}
	return hex.EncodeToString(privKeyBytes), hex.EncodeToString(pubKeyBytes)
}
