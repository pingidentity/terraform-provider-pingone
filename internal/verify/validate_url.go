// Copyright © 2026 Ping Identity Corporation

package verify

import (
	"fmt"
	"net/url"
	"strings"
)

// IsURLWithHTTPS checks a value is a valid URL with an `https://` scheme.
// Suitable for use with resource.TestCheckResourceAttrWith in acceptance tests.
func IsURLWithHTTPS(value string) error {
	return IsURLWithSchemes(value, "https")
}

// IsURLWithHTTPorHTTPS checks a value is a valid URL with an `http://` or `https://` scheme.
// Suitable for use with resource.TestCheckResourceAttrWith in acceptance tests.
func IsURLWithHTTPorHTTPS(value string) error {
	return IsURLWithSchemes(value, "http", "https")
}

// QuotedSchemes renders a list of URL schemes as quoted `scheme://` items joined by " or ",
// for use in validator diagnostic descriptions.
func QuotedSchemes(schemes []string) string {
	quoted := make([]string, len(schemes))
	for i, scheme := range schemes {
		quoted[i] = fmt.Sprintf("`%s://`", scheme)
	}
	return strings.Join(quoted, " or ")
}

// IsURLWithSchemes checks a value is a valid URL with one of the given schemes.
// A value is considered valid when it parses as a URL (net/url.Parse), includes a host,
// and has a case-insensitive scheme match. Percent-encoded characters are allowed.
func IsURLWithSchemes(value string, schemes ...string) error {
	u, err := url.Parse(value)
	if err != nil {
		return fmt.Errorf("value is not a valid URL: %w", err)
	}

	if u.Hostname() == "" {
		return fmt.Errorf("URL does not include a host, got: %s", value)
	}

	for _, scheme := range schemes {
		if strings.EqualFold(u.Scheme, scheme) {
			return nil
		}
	}

	return fmt.Errorf("URL scheme must be %s, got: %s", QuotedSchemes(schemes), value)
}
