// Copyright © 2025 Ping Identity Corporation

// Package verify provides validation utilities and constants for the PingOne Terraform provider.
package verify

import (
	"fmt"
	"slices"
	"strings"
)

// illegalOIDCattributeNames contains OIDC attribute names that are illegal and cannot be used for custom attributes
var illegalOIDCattributeNames = []string{"acr", "amr", "aud", "auth_time", "client_id", "env", "exp", "iat", "iss", "jti", "org", "p1.*", "scope", "sid", "sub"}

// overrideOIDCattributeNames contains OIDC attribute names that can be overridden with custom configurations
var overrideOIDCattributeNames = []string{"address.country", "address.formatted", "address.locality", "address.postal_code", "address.region", "address.street_address", "birthdate", "email", "email_verified", "family_name", "gender", "given_name", "locale", "middle_name", "name", "nickname", "phone_number", "phone_number_verified", "picture", "preferred_username", "profile", "updated_at", "website", "zoneinfo"}

// IllegalOIDCattributeNamesList returns a slice of OIDC attribute names that are illegal for custom use.
// These attribute names are reserved by the OIDC specification and PingOne and cannot be used
// for custom attribute mapping or schema extensions.
func IllegalOIDCattributeNamesList() []string {
	return illegalOIDCattributeNames
}

// IllegalOIDCAttributeNameString returns a formatted string containing all illegal OIDC attribute names.
// The attribute names are sorted alphabetically and formatted with backticks for documentation purposes.
// This function is useful for generating validation error messages that need to display
// the complete list of prohibited OIDC attribute names.
func IllegalOIDCAttributeNameString() string {
	slices.Sort(illegalOIDCattributeNames)

	v := make([]string, len(illegalOIDCattributeNames))
	for i, c := range illegalOIDCattributeNames {
		v[i] = fmt.Sprintf("`%s`", c)
	}
	return strings.Join(v, ", ")
}

// OverrideOIDCAttributeNameList returns a slice of OIDC attribute names that can be overridden.
// These attribute names represent standard OIDC claims that can be customized with
// alternative attribute mappings or configurations in PingOne applications.
func OverrideOIDCAttributeNameList() []string {
	return overrideOIDCattributeNames
}

// OverrideOIDCAttributeNameString returns a formatted string containing all overridable OIDC attribute names.
// The attribute names are sorted alphabetically and formatted with backticks for documentation purposes.
// This function is useful for generating documentation or informational messages that need to display
// the complete list of OIDC attribute names that can be customized.
func OverrideOIDCAttributeNameString() string {
	slices.Sort(overrideOIDCattributeNames)

	v := make([]string, len(overrideOIDCattributeNames))
	for i, c := range overrideOIDCattributeNames {
		v[i] = fmt.Sprintf("`%s`", c)
	}
	return strings.Join(v, ", ")
}
