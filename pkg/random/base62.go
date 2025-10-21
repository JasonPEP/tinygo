package random

import (
	"crypto/rand"
	"math/big"
)

const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var base = big.NewInt(int64(len(alphabet)))

// Code generates a random base62 string with the given length.
// It uses crypto/rand for better randomness.
func Code(length int) (string, error) {
	if length <= 0 {
		return "", nil
	}
	b := make([]byte, length)
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, base)
		if err != nil {
			return "", err
		}
		b[i] = alphabet[n.Int64()]
	}
	return string(b), nil
}
