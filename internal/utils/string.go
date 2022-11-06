package utils

import (
	"crypto/rand"
	"math/big"
)

func RandStringFromCharSet(strlen int, charSet string) (string, error) {
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		num, err := RandIntRange(len(charSet))
		if err != nil {
			return "", err
		}

		result[i] = charSet[num.Int64()]
	}
	return string(result), nil
}

func RandIntRange(max int) (*big.Int, error) {
	return rand.Int(rand.Reader, big.NewInt(int64(max)))
}
