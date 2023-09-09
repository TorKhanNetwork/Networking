package data_encryption

import (
	"crypto/rand"
	"crypto/rsa"

	"github.com/kataras/golog"
)

type KeysGenerator struct {
	asyncKeySize int
	syncKeySize  int
	SecretKey    []byte
	secretKeyIv  []byte
	PublicKey    rsa.PublicKey
	PrivateKey   rsa.PrivateKey
}

func NewGenerator() KeysGenerator {
	return KeysGenerator{
		asyncKeySize: 2048,
		syncKeySize:  16,
		SecretKey:    nil,
		secretKeyIv:  nil,
	}
}

func (keysGenerator *KeysGenerator) GenerateKeys(sync, async bool) {
	if sync {
		keysGenerator.secretKeyIv = make([]byte, 16)
		if _, err := rand.Read(keysGenerator.secretKeyIv); err != nil {
			golog.Fatal("Unable to generate secretKeyIv\n", err)
		}

		keysGenerator.SecretKey = make([]byte, keysGenerator.syncKeySize)
		if _, err := rand.Read(keysGenerator.SecretKey); err != nil {
			golog.Fatal("Unable to generate secretKey\n", err)
		}
	}
	if async {
		privateKey, err := rsa.GenerateKey(rand.Reader, keysGenerator.asyncKeySize)
		if err != nil {
			golog.Fatal("Unable to generate RSA keys\n", err)
		}
		publicKey := privateKey.PublicKey

		keysGenerator.PrivateKey = *privateKey
		keysGenerator.PublicKey = publicKey
	}
}

func (keysGenerator KeysGenerator) GetSyncKeyInfo() (info []byte) {
	info = append(info, keysGenerator.secretKeyIv...)
	info = append(info, keysGenerator.SecretKey...)
	return
}
