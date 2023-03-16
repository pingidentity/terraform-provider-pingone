package verify

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var P1ResourceIDRegexp = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
var P1DVResourceIDRegexp = regexp.MustCompile(`^[a-f0-9]{32}`)
var RFC3339Regexp = regexp.MustCompile(`^((?:(\d{4}-\d{2}-\d{2})T(\d{2}:\d{2}:\d{2}(?:\.\d+)?))(Z|[\+-]\d{2}:\d{2})?)$`)

// SDKv2
func ValidP1ResourceID(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if value == "" {
		return ws, errors
	}
	if !P1ResourceIDRegexp.MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q PingOne resource ID is malformed(%q): %q",
			k, P1ResourceIDRegexp, value))
	}

	return
}

// Framework
func P1ResourceIDValidator() validator.String {
	return stringvalidator.RegexMatches(P1ResourceIDRegexp, fmt.Sprintf("The PingOne resource ID is malformed.  Must match regex %q", P1ResourceIDRegexp))
}

func P1DVResourceIDValidator() validator.String {
	return stringvalidator.RegexMatches(P1DVResourceIDRegexp, fmt.Sprintf("The PingOne DaVinci resource ID is malformed.  Must match regex %q", P1DVResourceIDRegexp))
}
