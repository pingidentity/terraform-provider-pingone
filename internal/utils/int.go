// Copyright © 2025 Ping Identity Corporation

// Package utils provides utility functions for common operations in the PingOne Terraform provider.
package utils

// IntSliceToAnySlice converts a slice of int values to a slice of interface{} values.
// This function is useful for preparing integer data for Terraform framework attributes that expect []any types.
// Each int in the input slice is converted to an interface{} value in the output slice.
func IntSliceToAnySlice(v []int) []any {
	var result []interface{}
	for _, s := range v {
		result = append(result, s)
	}
	return result
}

// Int32SliceToAnySlice converts a slice of int32 values to a slice of interface{} values.
// This function is useful for preparing 32-bit integer data for Terraform framework attributes that expect []any types.
// Each int32 in the input slice is converted to an interface{} value in the output slice.
func Int32SliceToAnySlice(v []int32) []any {
	var result []interface{}
	for _, s := range v {
		result = append(result, s)
	}
	return result
}

// Int64SliceToAnySlice converts a slice of int64 values to a slice of interface{} values.
// This function is useful for preparing 64-bit integer data for Terraform framework attributes that expect []any types.
// Each int64 in the input slice is converted to an interface{} value in the output slice.
func Int64SliceToAnySlice(v []int64) []any {
	var result []interface{}
	for _, s := range v {
		result = append(result, s)
	}
	return result
}
