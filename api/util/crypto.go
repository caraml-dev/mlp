package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"io"
)

func CreateHash(key string) string {
	hasher := sha256.Sum256([]byte(key))
	return hex.EncodeToString(hasher[:])
}

func Encrypt(plainText, passphrase string) (string, error) {
	key, err := hex.DecodeString(passphrase)
	if err != nil {
		return "", err
	}
	block, _ := aes.NewCipher(key)
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	cipherText := gcm.Seal(nonce, nonce, []byte(plainText), nil)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func Decrypt(cipherText, passphrase string) (string, error) {
	key, err := hex.DecodeString(passphrase)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonceSize := gcm.NonceSize()
	ciphertextByte, _ := base64.StdEncoding.DecodeString(cipherText)
	nonce, ciphertext := ciphertextByte[:nonceSize], ciphertextByte[nonceSize:]
	plainText, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}
