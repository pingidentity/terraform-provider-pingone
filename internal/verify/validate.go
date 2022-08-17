package verify

import (
	"fmt"
	"regexp"
)

var p1ResourceID = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

func ValidP1ResourceID(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if value == "" {
		return ws, errors
	}
	if !p1ResourceID.MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q PingOne resource ID is malformed(%q): %q",
			k, p1ResourceID, value))
	}

	return
}
