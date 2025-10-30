package security

import (
	"crypto/rand"
	mrand "math/rand"
)

const (
	AlphaNum = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
)

// RandAlphaNum generates random alphanumeric string.
func RandAlphaNum(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = AlphaNum[mrand.Intn(len(AlphaNum))] // nolint:gosec
	}

	return string(b)
}

// Rand generates n random bytes.
func Rand(n int) []byte {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil { // nolint:gosec
		panic(err) // something should be really wrong if we got here
	}

	return b
}
