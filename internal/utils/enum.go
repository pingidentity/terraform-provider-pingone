// Copyright © 2025 Ping Identity Corporation

// Package utils provides utility functions for common operations in the PingOne Terraform provider.
// This package contains helper functions for type conversions, string manipulation, JSON processing,
// and other general-purpose utilities used throughout the provider.
package utils

import "encoding/json"

// EnumToString converts a PingOne SDK enum value to its string representation.
// It returns an empty string if the conversion fails or if the enum value is invalid.
// This function is useful for converting SDK enum types to string values for Terraform state management.
func EnumToString(enum interface{}) string {
	b, e := json.Marshal(enum)
	if e != nil {
		return ""
	}

	var s string
	e = json.Unmarshal(b, &s)
	if e != nil {
		return ""
	}

	return s
}

// EnumSliceToStringSlice converts a slice of PingOne SDK enum values to a slice of strings.
// It returns an empty slice if the conversion fails or if any enum values are invalid.
// This function is useful for converting SDK enum slices to string slices for Terraform attribute handling.
func EnumSliceToStringSlice(enum interface{}) []string {
	b, e := json.Marshal(enum)
	if e != nil {
		return []string{}
	}

	var s []string
	e = json.Unmarshal(b, &s)
	if e != nil {
		return []string{}
	}

	return s
}

// EnumSliceToAnySlice converts a slice of PingOne SDK enum values to a slice of interface{} values.
// It first converts the enum slice to strings, then converts the string slice to interface{} slice.
// This function is useful for preparing enum data for Terraform framework attributes that expect []any types.
func EnumSliceToAnySlice(enum interface{}) []any {
	v := EnumSliceToStringSlice(enum)
	return StringSliceToAnySlice(v)
}
