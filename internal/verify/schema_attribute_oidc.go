package verify

import (
	"fmt"
	"slices"
	"strings"
)

var illegalOIDCattributeNames = []string{"acr", "amr", "aud", "auth_time", "client_id", "env", "exp", "iat", "iss", "jti", "org", "p1.*", "scope", "sid", "sub"}
var overrideOIDCattributeNames = []string{"address.country", "address.formatted", "address.locality", "address.postal_code", "address.region", "address.street_address", "birthdate", "email", "email_verified", "family_name", "gender", "given_name", "locale", "middle_name", "name", "nickname", "phone_number", "phone_number_verified", "picture", "preferred_username", "profile", "updated_at", "website", "zoneinfo"}

func IllegalOIDCattributeNamesList() []string {
	return illegalOIDCattributeNames
}

func IllegalOIDCAttributeNameString() string {
	slices.Sort(illegalOIDCattributeNames)

	v := make([]string, len(illegalOIDCattributeNames))
	for i, c := range illegalOIDCattributeNames {
		v[i] = fmt.Sprintf("`%s`", c)
	}
	return strings.Join(v, ", ")
}

func OverrideOIDCAttributeNameList() []string {
	return overrideOIDCattributeNames
}

func OverrideOIDCAttributeNameString() string {
	slices.Sort(overrideOIDCattributeNames)

	v := make([]string, len(overrideOIDCattributeNames))
	for i, c := range overrideOIDCattributeNames {
		v[i] = fmt.Sprintf("`%s`", c)
	}
	return strings.Join(v, ", ")
}
