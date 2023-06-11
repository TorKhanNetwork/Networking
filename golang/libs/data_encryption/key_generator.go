package data_encryption

import (
	"crypto/rand"
	"crypto/rsa"

	"github.com/kataras/golog"
)

type KeyGenerator struct {
	asyncKeySize int
	syncKeySize  int
	SecretKey    []byte
	secretKeyIv  []byte
	PublicKey    rsa.PublicKey
	privateKey   rsa.PrivateKey
}

func NewGenerator() KeyGenerator {
	return KeyGenerator{
		asyncKeySize: 2048,
		syncKeySize:  16,
		SecretKey:    nil,
		secretKeyIv:  nil,
	}
}

func (keyGenerator *KeyGenerator) GenerateKeys(sync, async bool) {
	if sync {
		keyGenerator.secretKeyIv = make([]byte, 16)
		if _, err := rand.Read(keyGenerator.secretKeyIv); err != nil {
			golog.Fatal("Unable to generate secretKeyIv\n", err)
		}

		keyGenerator.SecretKey = make([]byte, keyGenerator.syncKeySize)
		if _, err := rand.Read(keyGenerator.SecretKey); err != nil {
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
		keyGenerator.PublicKey = *publicKey
	}
}

func (keyGenerator KeyGenerator) GetSyncKeyInfo() (info []byte) {
	info = append(info, keyGenerator.secretKeyIv...)
	info = append(info, keyGenerator.SecretKey...)
	return
}
