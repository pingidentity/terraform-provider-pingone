// Copyright © 2025 Ping Identity Corporation

// Package utils provides utility functions for common operations in the PingOne Terraform provider.
package utils

import (
	"encoding/json"
	"reflect"
)

// DeepEqualJSON compares two JSON byte arrays for deep equality of their content.
// It returns true if both byte arrays represent the same JSON structure and values, false otherwise.
// The comparison is performed by unmarshaling both JSON byte arrays into interface{} values
// and using reflect.DeepEqual to compare the resulting structures.
// This function is useful for detecting changes in JSON configuration data stored in Terraform state.
func DeepEqualJSON(a, b []byte) bool {
	var aj, bj interface{}
	if err := json.Unmarshal(a, &aj); err != nil {
		return false
	}
	if err := json.Unmarshal(b, &bj); err != nil {
		return false
	}
	return reflect.DeepEqual(aj, bj)
}
