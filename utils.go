package fmk

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"math"
)

var (
	CipherTextTooShort error = errors.New("FMK: Ciphertext too short")
)

const (
        FRONT_SLASH = 47 // literally ascii/utf8 for '/'
)

// This is constant with respect to N, too much crypto packages messing with my head
func equal(a []byte, b []byte) bool {
	is_equal := true
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			is_equal = false
		}
	}
	return is_equal
}

func encrypt(plaintext []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, CipherTextTooShort
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

func armor(start []byte) string {
	return base64.StdEncoding.EncodeToString(start)
}

func dearmor(end string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(end)
}

func pow(a, b int) byte {
	return byte(math.Pow(float64(a), float64(b)))
}

func ParseURL(url string) []string {
        urlBytes := []byte(url)
        wordBuilder := make([]byte, 0)
        result := make([]string, 0)
        for _, char := range urlBytes {
                if char == FRONT_SLASH {
                        result = append(result, string(wordBuilder))
                        wordBuilder = make([]byte, 0)
                } else {
                        wordBuilder = append(wordBuilder, char)
                }
        }
        result = append(result, string(wordBuilder))
        return result
}
