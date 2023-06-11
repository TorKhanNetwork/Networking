package client

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"net"
)

func GetMacAddress() (mac string, err error) {
	ifas, err := net.Interfaces()
	if err != nil {
		return
	}
	for _, ifa := range ifas {
		a := ifa.HardwareAddr.String()
		if a != "" {
			mac = a
			return
		}
	}
	return
}

func ParseRsaPublicKeyFromPemStr(pubPEM string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte("-----BEGIN PUBLIC KEY-----\n" + pubPEM + "\n-----END PUBLIC KEY-----"))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch pub := pub.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		break // fall through
	}
	return nil, errors.New("key type is not RSA")
}
