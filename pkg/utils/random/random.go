package random

import (
	"encoding/base64"
	"github.com/google/uuid"
)

func UniqueString() (string, error) {
	u := uuid.New()
	b, err := u.MarshalBinary()
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}
