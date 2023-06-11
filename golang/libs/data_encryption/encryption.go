package data_encryption

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"errors"
)

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5Trimming(encrypt []byte) []byte {
	padding := encrypt[len(encrypt)-1]
	return encrypt[:len(encrypt)-int(padding)]
}

func EncryptSecretKey(keyGenerator KeyGenerator) (string, error) {
	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, &keyGenerator.PublicKey, keyGenerator.GetSyncKeyInfo())
	return base64.StdEncoding.EncodeToString(cipherText), err
}

func DecryptSecretKey(data string, keyGenerator *KeyGenerator) (err error) {
	cipherData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return
	}
	decryptedKey, err := rsa.DecryptPKCS1v15(rand.Reader, &keyGenerator.privateKey, cipherData)
	if err != nil {
		return
	}
	keyGenerator.secretKeyIv = decryptedKey[:15]
	keyGenerator.SecretKey = decryptedKey[16:]
	return
}

func Encrypt(data string, keyGenerator KeyGenerator) (cipherText string, err error) {
	block, err := aes.NewCipher(keyGenerator.SecretKey)
	if err != nil {
		return
	}
	if data == "" {
		return cipherText, errors.New("encrypt: Empty plain text")
	}
	ecb := cipher.NewCBCEncrypter(block, keyGenerator.secretKeyIv)
	content := PKCS5Padding([]byte(data), block.BlockSize())
	crypted := make([]byte, len(content))
	ecb.CryptBlocks(crypted, content)
	cipherText = base64.StdEncoding.EncodeToString(crypted)
	return
}

func Decrypt(data string, keyGenerator KeyGenerator) (plainText string, err error) {
	block, err := aes.NewCipher(keyGenerator.SecretKey)
	if err != nil {
		return
	}
	if data == "" {
		return plainText, errors.New("decrypt: Empty encrypted text")
	}
	crypt, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return
	}
	ecb := cipher.NewCBCDecrypter(block, keyGenerator.secretKeyIv)
	decrypted := make([]byte, len(crypt))
	ecb.CryptBlocks(decrypted, crypt)

	return string(PKCS5Trimming(decrypted)), err
}
