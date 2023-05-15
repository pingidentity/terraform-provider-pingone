package utils

import (
	"encoding/json"
	"reflect"
)

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
