package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"

	"gitlab.com/timkado/api/daisi-rest-postgres/internal/config"
)

const (
	algorithmNonceSize = 12
	algorithmTagSize   = 16
)

func Decrypt(token string) (string, error) {
	key, err := hex.DecodeString(config.LoadConfig().SecretKey)
	if err != nil {
		return "", err
	}
	ciphertextAndNonce, err := hex.DecodeString(token)
	if err != nil {
		return "", err
	}

	if len(ciphertextAndNonce) <= algorithmNonceSize {
		return "", errors.New("ciphertext and nonce size is too short")
	}

	nonce := ciphertextAndNonce[:algorithmNonceSize]
	ciphertext := ciphertextAndNonce[algorithmNonceSize:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
