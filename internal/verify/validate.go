// Copyright © 2025 Ping Identity Corporation

// Package verify provides validation utilities and constants for the PingOne Terraform provider.
// This package contains validators, regular expressions, and validation functions for ensuring
// data integrity and format compliance across PingOne resources.
package verify

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Regular expression patterns for validation across the provider
var (
	// P1ResourceIDRegexp matches PingOne resource UUID format (without anchors)
	P1ResourceIDRegexp = regexp.MustCompile(`[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`)
	// P1ResourceIDRegexpFullString matches complete PingOne resource UUID format
	P1ResourceIDRegexpFullString = regexp.MustCompile(fmt.Sprintf(`^%s$`, P1ResourceIDRegexp.String()))
	// P1DVResourceIDRegexp matches PingOne DaVinci resource ID format (without anchors)
	P1DVResourceIDRegexp = regexp.MustCompile(`[a-f0-9]{32}`)
	// P1DVResourceIDRegexpFullString matches complete PingOne DaVinci resource ID format
	P1DVResourceIDRegexpFullString = regexp.MustCompile(fmt.Sprintf(`^%s$`, P1DVResourceIDRegexp.String()))
	// RFC3339Regexp matches RFC3339 datetime format
	RFC3339Regexp = regexp.MustCompile(`^((?:(\d{4}-\d{2}-\d{2})T(\d{2}:\d{2}:\d{2}(?:\.\d+)?))(Z|[\+-]\d{2}:\d{2})?)$`)
	// IPv4Regexp matches IPv4 addresses with optional CIDR notation (without anchors)
	IPv4Regexp = regexp.MustCompile(`(?:\d{1,3}\.){3}\d{1,3}(?:/\d{1,2})?`)
	// IPv4RegexpFull matches complete IPv4 addresses with optional CIDR notation
	IPv4RegexpFull = regexp.MustCompile(fmt.Sprintf(`^%s$`, IPv4Regexp.String()))
	// IPv6Regexp matches IPv6 addresses with optional CIDR notation (without anchors)
	IPv6Regexp = regexp.MustCompile(`(?:[A-F0-9]{1,4}:){7}[A-F0-9]{1,4}(?:/\d{1,3})?`)
	// IPv6RegexpFull matches complete IPv6 addresses with optional CIDR notation
	IPv6RegexpFull = regexp.MustCompile(fmt.Sprintf(`^%s$`, IPv6Regexp.String()))
	// IPv4IPv6Regexp matches either IPv4 or IPv6 addresses
	IPv4IPv6Regexp = regexp.MustCompile(fmt.Sprintf(`%s|%s`, IPv4RegexpFull.String(), IPv6RegexpFull.String()))
	// HexColorCode matches hexadecimal color codes in 3 or 6 digit format
	HexColorCode = regexp.MustCompile(`^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$`)
	// IsDomain matches domain name format
	IsDomain = regexp.MustCompile(`^(?:[\w-]+\.)+[a-z]{2,}$`)
	// IsURLWithHTTPorHTTPS matches URLs with HTTP or HTTPS protocol
	IsURLWithHTTPorHTTPS = regexp.MustCompile(fmt.Sprintf(`^http[s]{0,1}:\/\/%s`, urlRegexStringWithoutProtocol))
	// IsURLWithHTTPS matches URLs with HTTPS protocol only
	IsURLWithHTTPS = regexp.MustCompile(fmt.Sprintf(`^https:\/\/%s`, urlRegexStringWithoutProtocol))
	// IsTwoCharCountryCode matches ISO 3166-1 alpha-2 country codes
	IsTwoCharCountryCode = regexp.MustCompile(`^(A(D|E|F|G|I|L|M|N|O|R|S|T|Q|U|W|X|Z)|B(A|B|D|E|F|G|H|I|J|L|M|N|O|R|S|T|V|W|Y|Z)|C(A|C|D|F|G|H|I|K|L|M|N|O|R|U|V|X|Y|Z)|D(E|J|K|M|O|Z)|E(C|E|G|H|R|S|T)|F(I|J|K|M|O|R)|G(A|B|D|E|F|G|H|I|L|M|N|P|Q|R|S|T|U|W|Y)|H(K|M|N|R|T|U)|I(D|E|Q|L|M|N|O|R|S|T)|J(E|M|O|P)|K(E|G|H|I|M|N|P|R|W|Y|Z)|L(A|B|C|I|K|R|S|T|U|V|Y)|M(A|C|D|E|F|G|H|K|L|M|N|O|Q|P|R|S|T|U|V|W|X|Y|Z)|N(A|C|E|F|G|I|L|O|P|R|U|Z)|OM|P(A|E|F|G|H|K|L|M|N|R|S|T|W|Y)|QA|R(E|O|S|U|W)|S(A|B|C|D|E|G|H|I|J|K|L|M|N|O|R|T|V|Y|Z)|T(C|D|F|G|H|J|K|L|M|N|O|R|T|V|W|Z)|U(A|G|M|S|Y|Z)|V(A|C|E|G|I|N|U)|W(F|S)|Y(E|T)|Z(A|M|W))$`)
	// IsHostname matches hostname format including paths and query parameters
	IsHostname = regexp.MustCompile(`^(?:[\w-]+\.)+[a-z]{2,}(?:\/[\w-]+)*(?:\/[\w.-]+)?\/?(?:\?.*)?$`)
)

// URL regex component used in URL validation patterns
var urlRegexStringWithoutProtocol = `(?:(?:[\w-]+\.)+[a-z]{2,}|localhost)(?:\:[\d]{1,})*(?:\/[\w.-]+)*(?:\/[\w\:.-]+)?\/?(?:\?.*)?$`

// ValidP1ResourceID validates PingOne resource ID format using SDKv2 validation.
// It returns any warnings and errors encountered during validation.
// The v parameter is the value to validate (should be a string).
// The k parameter is the key name for error reporting purposes.
// This function checks for empty values and validates against the PingOne UUID format.
// It is used for backward compatibility with SDKv2-based resources.
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

// P1ResourceIDValidator creates a string validator for PingOne resource IDs using the Terraform Framework.
// It returns a validator that checks string values against the PingOne UUID format.
// This validator ensures that resource IDs conform to the expected PingOne resource identifier pattern.
// It is used for Framework-based resources and provides descriptive error messages for validation failures.
func P1ResourceIDValidator() validator.String {
	return stringvalidator.RegexMatches(P1ResourceIDRegexpFullString, fmt.Sprintf("The PingOne resource ID is malformed.  Must match regex %q", P1ResourceIDRegexpFullString))
}

// P1DVResourceIDValidator creates a string validator for PingOne DaVinci resource IDs using the Terraform Framework.
// It returns a validator that checks string values against the PingOne DaVinci resource ID format.
// This validator ensures that DaVinci resource IDs conform to the expected 32-character hexadecimal pattern.
// It is used for DaVinci-specific resources and provides descriptive error messages for validation failures.
func P1DVResourceIDValidator() validator.String {
	return stringvalidator.RegexMatches(P1DVResourceIDRegexpFullString, fmt.Sprintf("The PingOne DaVinci resource ID is malformed.  Must match regex %q", P1DVResourceIDRegexpFullString))
}

// LocaleValidator creates a compiled regular expression for validating locale codes.
// It returns a regexp.Regexp that matches valid ISO language and locale codes.
// The validator is built dynamically from the complete list of supported ISO locale codes
// and ensures that locale values conform to internationalization standards.
// This function is used for validating language and locale settings across PingOne resources.
func LocaleValidator() *regexp.Regexp {
	isoList := FullIsoList()
	escaped := make([]string, len(isoList))
	for i, v := range isoList {
		escaped[i] = regexp.QuoteMeta(v)
	}
	pattern := "(" + strings.Join(escaped, "|") + ")"
	return regexp.MustCompile(pattern)
}
