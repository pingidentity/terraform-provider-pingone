// Copyright © 2026 Ping Identity Corporation

package verify

import (
	"fmt"
	"regexp"
	"strings"

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
var urlRegexStringWithoutProtocol = `(?:(?:[\w-]+\.)+[a-z]{2,}|localhost)(?:\:[\d]{1,})*(?:\/[\w.-]+)*(?:\/[\w\:.-]+)?\/?(?:\?.*)?$`
var IsURLWithHTTPorHTTPS = regexp.MustCompile(fmt.Sprintf(`^http[s]{0,1}:\/\/%s`, urlRegexStringWithoutProtocol))
var IsURLWithHTTPS = regexp.MustCompile(fmt.Sprintf(`^https:\/\/%s`, urlRegexStringWithoutProtocol))
var IsTwoCharCountryCode = regexp.MustCompile(`^[A-Z]{2}$`)
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

// Locale
func LocaleValidator() *regexp.Regexp {
	isoList := FullIsoList()
	escaped := make([]string, len(isoList))
	for i, v := range isoList {
		escaped[i] = regexp.QuoteMeta(v)
	}
	pattern := "(" + strings.Join(escaped, "|") + ")"
	return regexp.MustCompile(pattern)
}
