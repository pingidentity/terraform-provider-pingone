package utils

import "math/rand"

func RandStringFromCharSet(strlen int, charSet string) string {
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = charSet[RandIntRange(0, len(charSet))]
	}
	return string(result)
}

func RandIntRange(min int, max int) int {
	return rand.Intn(max-min) + min
}
