package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

// EncryptString encrypts plain text using AES-GCM and returns base64 encoded ciphertext.
func EncryptString(plain, key string) (string, error) {
	if key == "" {
		return plain, nil
	}
	bkey := []byte(key)
	if len(bkey) != 32 {
		return "", fmt.Errorf("encryption key must be 32 bytes")
	}
	block, err := aes.NewCipher(bkey)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	cipherText := gcm.Seal(nonce, nonce, []byte(plain), nil)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// DecryptString decrypts a base64 encoded ciphertext using AES-GCM.
func DecryptString(cipherText, key string) (string, error) {
	if key == "" {
		return cipherText, nil
	}
	data, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}
	bkey := []byte(key)
	if len(bkey) != 32 {
		return "", fmt.Errorf("encryption key must be 32 bytes")
	}
	block, err := aes.NewCipher(bkey)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	if len(data) < gcm.NonceSize() {
		return "", fmt.Errorf("ciphertext too short")
	}
	nonce := data[:gcm.NonceSize()]
	ciphertext := data[gcm.NonceSize():]
	plain, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plain), nil
}
