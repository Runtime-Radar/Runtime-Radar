package cipher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
)

var allowedKeyLengths = []int{16, 24, 32}

func EncryptStringAsHex(in, hKey string) string {
	key, err := hex.DecodeString(hKey)
	if err != nil {
		panic(err)
	}

	out := Encrypt([]byte(in), key)

	return hex.EncodeToString(out)
}

func DecryptHexAsString(h, hKey string) string {
	key, err := hex.DecodeString(hKey)
	if err != nil {
		panic(err)
	}

	b, err := hex.DecodeString(h)
	if err != nil {
		panic(err)
	}

	out := Decrypt(b, key)

	return string(out)
}

func Encrypt(in, key []byte) []byte {
	plaintext := in

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return ciphertext
}

func Decrypt(in, key []byte) []byte {
	ciphertext := in

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)

	return ciphertext
}

type Crypt struct {
	key []byte
}

func (c *Crypt) EncryptStringAsHex(in string) string {
	ciphertext := Encrypt([]byte(in), c.key)
	return hex.EncodeToString(ciphertext)
}

func (c *Crypt) DecryptHexAsString(in string) string {
	b, err := hex.DecodeString(in)
	if err != nil {
		panic(err)
	}

	out := Decrypt(b, c.key)
	return string(out)
}

type Crypter interface {
	EncryptStringAsHex(in string) string
	DecryptHexAsString(in string) string
}

func NewCrypt(h string) (*Crypt, error) {
	key, err := ParseKey(h)
	if err != nil {
		return nil, err
	}

	return &Crypt{key: key}, nil
}

func ParseKey(h string) ([]byte, error) {
	if h == "" {
		return nil, errors.New("key is empty")
	}
	key, err := hex.DecodeString(h)
	if err != nil {
		return nil, err
	}

	// check key length
	for _, keyLen := range allowedKeyLengths {
		if len(key) == keyLen {
			return key, nil
		}
	}

	return nil, errors.New("invalid key length")
}
