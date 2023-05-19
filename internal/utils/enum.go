package utils

import "encoding/json"

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
