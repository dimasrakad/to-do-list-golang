package utils

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)

	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func GenerateRandomToken(byteLength int) (string, error) {
	bytes, err := GenerateRandomBytes(byteLength)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(bytes), nil
}
