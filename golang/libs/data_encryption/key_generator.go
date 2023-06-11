package data_encryption

import (
	"crypto/rand"
	"crypto/rsa"

	"github.com/kataras/golog"
)

type KeyGenerator struct {
	asyncKeySize int
	syncKeySize  int
	secretKey    []byte
	secretKeyIv  []byte
	publicKey    rsa.PublicKey
	privateKey   rsa.PrivateKey
}

func NewGenerator() KeyGenerator {
	return KeyGenerator{
		asyncKeySize: 2048,
		syncKeySize:  16,
		secretKey:    nil,
		secretKeyIv:  nil,
	}
}

func (keyGenerator *KeyGenerator) generateKeys(sync, async bool) {
	if sync {
		keyGenerator.secretKeyIv = make([]byte, 16)
		if _, err := rand.Read(keyGenerator.secretKeyIv); err != nil {
			golog.Fatal("Unable to generate secretKeyIv\n", err)
		}

		keyGenerator.secretKey = make([]byte, keyGenerator.syncKeySize)
		if _, err := rand.Read(keyGenerator.secretKey); err != nil {
			golog.Fatal("Unable to generate secretKey\n", err)
		}
	}
	if async {
		privateKey, err := rsa.GenerateKey(rand.Reader, keyGenerator.asyncKeySize)
		if err != nil {
			golog.Fatal("Unable to generate RSA keys\n", err)
		}
		publicKey := &privateKey.PublicKey

		keyGenerator.privateKey = *privateKey
		keyGenerator.publicKey = *publicKey
	}
}

func (keyGenerator KeyGenerator) GetSyncKeyInfo() (info []byte) {
	info = append(info, keyGenerator.secretKeyIv...)
	info = append(info, keyGenerator.secretKey...)
	return
}
