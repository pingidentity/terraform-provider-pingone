package verify

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var P1ResourceIDRegexp = regexp.MustCompile(`[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`)
var P1ResourceIDRegexpFullString = regexp.MustCompile(fmt.Sprintf(`^%s$`, P1ResourceIDRegexp.String()))
var P1DVResourceIDRegexp = regexp.MustCompile(`[a-f0-9]{32}`)
var P1DVResourceIDRegexpFullString = regexp.MustCompile(fmt.Sprintf(`^%s$`, P1DVResourceIDRegexp.String()))
var RFC3339Regexp = regexp.MustCompile(`^((?:(\d{4}-\d{2}-\d{2})T(\d{2}:\d{2}:\d{2}(?:\.\d+)?))(Z|[\+-]\d{2}:\d{2})?)$`)
var IPv4Regexp = regexp.MustCompile(`(?:\d{1,3}\.){3}\d{1,3}(?:/\d{1,2})?`)
var IPv4RegexpFull = regexp.MustCompile(fmt.Sprintf(`^%s$`, IPv4Regexp.String()))
var IPv6Regexp = regexp.MustCompile(`(?:[A-F0-9]{1,4}:){7}[A-F0-9]{1,4}(?:/\d{1,3})?`)
var IPv6RegexpFull = regexp.MustCompile(fmt.Sprintf(`^%s$`, IPv6Regexp.String()))
var IPv4IPv6Regexp = regexp.MustCompile(fmt.Sprintf(`%s|%s`, IPv4RegexpFull.String(), IPv6RegexpFull.String()))
var HexColorCode = regexp.MustCompile(`^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$`)

var IsDomain = regexp.MustCompile(`^(?:[\w-]+\.)+[a-z]{2,}$`)
var urlRegexStringWithoutProtocol = `(?:[\w-]+\.)+[a-z]{2,}(?:\/[\w.-]+)*(?:\/[\w\:.-]+)?\/?(?:\?.*)?$`
var IsURLWithHTTPorHTTPS = regexp.MustCompile(fmt.Sprintf(`^http[s]{0,1}:\/\/%s`, urlRegexStringWithoutProtocol))
var IsURLWithHTTPS = regexp.MustCompile(fmt.Sprintf(`^https:\/\/%s`, urlRegexStringWithoutProtocol))
var IsTwoCharCountryCode = regexp.MustCompile(`^(A(D|E|F|G|I|L|M|N|O|R|S|T|Q|U|W|X|Z)|B(A|B|D|E|F|G|H|I|J|L|M|N|O|R|S|T|V|W|Y|Z)|C(A|C|D|F|G|H|I|K|L|M|N|O|R|U|V|X|Y|Z)|D(E|J|K|M|O|Z)|E(C|E|G|H|R|S|T)|F(I|J|K|M|O|R)|G(A|B|D|E|F|G|H|I|L|M|N|P|Q|R|S|T|U|W|Y)|H(K|M|N|R|T|U)|I(D|E|Q|L|M|N|O|R|S|T)|J(E|M|O|P)|K(E|G|H|I|M|N|P|R|W|Y|Z)|L(A|B|C|I|K|R|S|T|U|V|Y)|M(A|C|D|E|F|G|H|K|L|M|N|O|Q|P|R|S|T|U|V|W|X|Y|Z)|N(A|C|E|F|G|I|L|O|P|R|U|Z)|OM|P(A|E|F|G|H|K|L|M|N|R|S|T|W|Y)|QA|R(E|O|S|U|W)|S(A|B|C|D|E|G|H|I|J|K|L|M|N|O|R|T|V|Y|Z)|T(C|D|F|G|H|J|K|L|M|N|O|R|T|V|W|Z)|U(A|G|M|S|Y|Z)|V(A|C|E|G|I|N|U)|W(F|S)|Y(E|T)|Z(A|M|W))$`)
var IsHostname = regexp.MustCompile(`^(?:[\w-]+\.)+[a-z]{2,}(?:\/[\w-]+)*(?:\/[\w.-]+)?\/?(?:\?.*)?$`)

// SDKv2
func ValidP1ResourceID(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if value == "" {
		return ws, errors
	}
	if !P1ResourceIDRegexpFullString.MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q PingOne resource ID is malformed(%q): %q",
			k, P1ResourceIDRegexpFullString, value))
	}

	return
}

// Framework
func P1ResourceIDValidator() validator.String {
	return stringvalidator.RegexMatches(P1ResourceIDRegexpFullString, fmt.Sprintf("The PingOne resource ID is malformed.  Must match regex %q", P1ResourceIDRegexpFullString))
}

func P1DVResourceIDValidator() validator.String {
	return stringvalidator.RegexMatches(P1DVResourceIDRegexpFullString, fmt.Sprintf("The PingOne DaVinci resource ID is malformed.  Must match regex %q", P1DVResourceIDRegexpFullString))
}
