package random

import (
	"crypto/rand"
	"encoding/base64"
)

func String(numberOfBytes uint) (string, error) {
	b := make([]byte, numberOfBytes)
	_, err := rand.Read(b)

	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}
