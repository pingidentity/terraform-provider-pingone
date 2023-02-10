package verify

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var p1ResourceID = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

// SDKv2
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

// Framework
func P1ResourceIDValidator() validator.String {
	return stringvalidator.RegexMatches(p1ResourceID, fmt.Sprintf("The PingOne resource ID is malformed.  Must match regex %q", p1ResourceID))
}
