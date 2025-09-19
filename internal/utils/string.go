// Copyright © 2025 Ping Identity Corporation

// Package utils provides utility functions for common operations in the PingOne Terraform provider.
package utils

import (
	"crypto/rand"
	"math/big"
)

// RandStringFromCharSet generates a cryptographically secure random string of specified length from a character set.
// It returns the generated string and any error encountered during random number generation.
// The strlen parameter specifies the desired length of the generated string.
// The charSet parameter defines the characters that can be used in the generated string.
// This function is useful for generating secure random identifiers, passwords, or tokens.
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

// RandIntRange generates a cryptographically secure random integer within the specified range.
// It returns a *big.Int containing the random value and any error encountered during generation.
// The max parameter defines the exclusive upper bound for the random number (0 to max-1).
// This function is used internally by RandStringFromCharSet for character selection.
func RandIntRange(max int) (*big.Int, error) {
	return rand.Int(rand.Reader, big.NewInt(int64(max)))
}

// StringSliceToAnySlice converts a slice of strings to a slice of interface{} values.
// This function is useful for preparing string data for Terraform framework attributes that expect []any types.
// Each string in the input slice is converted to an interface{} value in the output slice.
func StringSliceToAnySlice(v []string) []any {
	var result []interface{}
	for _, s := range v {
		result = append(result, s)
	}
	return result
}
