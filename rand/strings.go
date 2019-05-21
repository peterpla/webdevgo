package rand

import (
	"crypto/rand"
	"encoding/base64"
)

const RememberTokenBytes = 32

// RememberToken returns a fixed-length remember token
func RememberToken() (string, error) {
	return String(RememberTokenBytes)
}

// String returns a string of length n, containing random base64 encoded data,
// or on error, an empty string ""
func String(nBytes int) (string, error) {
	b, err := Bytes(nBytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// Bytes returns []byte of length n containing random data and nil,
// or will return an error if there was one. This uses the
// crypto/rand package so it is safe to use with things
// like remember tokens
func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}
