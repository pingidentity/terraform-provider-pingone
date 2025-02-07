// Copyright © 2025 Ping Identity Corporation

package utils

func IntSliceToAnySlice(v []int) []any {
	var result []interface{}
	for _, s := range v {
		result = append(result, s)
	}
	return result
}

func Int32SliceToAnySlice(v []int32) []any {
	var result []interface{}
	for _, s := range v {
		result = append(result, s)
	}
	return result
}

func Int64SliceToAnySlice(v []int64) []any {
	var result []interface{}
	for _, s := range v {
		result = append(result, s)
	}
	return result
}
