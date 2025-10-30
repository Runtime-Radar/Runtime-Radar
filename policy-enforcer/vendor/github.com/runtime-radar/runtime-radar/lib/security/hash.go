package security

import (
	"crypto/sha512"
	"encoding/hex"
)

// HashSHA512 produces hash of arbitrary []byte input.
func HashSHA512(in []byte) []byte {
	h := sha512.New()
	h.Write(in)

	return h.Sum(nil)
}

// HashSHA512AsHex produces hash of arbitrary []byte input encoded as hex string.
func HashSHA512AsHex(in []byte) string {
	b := HashSHA512(in)

	return hex.EncodeToString(b)
}

// HashSaltedSHA512 produces hash of arbitrary []byte input with addition of salt.
func HashSaltedSHA512(in, salt []byte) []byte {
	h := sha512.New()
	h.Write(in)
	h.Write(salt)

	return h.Sum(nil)
}

// HashSaltedSHA512AsHex produces hash of arbitrary []byte input with addition of salt encoded as hex string.
func HashSaltedSHA512AsHex(in, salt []byte) string {
	b := HashSaltedSHA512(in, salt)

	return hex.EncodeToString(b)
}
