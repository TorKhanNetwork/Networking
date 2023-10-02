package server

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"strings"

	"github.com/TorKhanNetwork/Networking/data_encryption"
	"github.com/kataras/golog"
)

func ExportRsaPublicKeyToString(publicKey *rsa.PublicKey) (string, error) {
	pubkey_bytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", err
	}
	pubkey_pem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: pubkey_bytes,
		},
	)
	split := strings.Split(string(pubkey_pem), "\n")
	return strings.Join(split[1:len(split)-2], ""), nil
}

func ReadAsyncKeys(kg *data_encryption.KeysGenerator) (ok bool) {
	bytes, err := os.ReadFile("private.key")
	if err != nil {
		return
	}
	block, _ := pem.Decode(bytes)
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return
	}
	kg.PrivateKey = *key
	kg.PublicKey = key.PublicKey
	return true
}

func WriteAsyncKeys(kg data_encryption.KeysGenerator) {
	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(&kg.PrivateKey),
	}
	privateKeyFile, err := os.Create("private.key")
	if err != nil {
		golog.Errorf("Error creating private key file : %s", err)
	}
	pem.Encode(privateKeyFile, privateKeyPEM)
	privateKeyFile.Close()
}
