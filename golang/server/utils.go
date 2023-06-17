package server

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"strings"
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