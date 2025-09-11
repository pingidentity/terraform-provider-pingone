// Copyright © 2025 Ping Identity Corporation

package stringvalidator

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/schemavalidator"
)

// IsRequiredIfRegexMatchesPathValue creates a validator that requires the current attribute when another attribute matches a regex pattern.
// It returns a validator that ensures the current string attribute is not null or empty
// when any of the specified path attributes match the provided regular expression.
//
// The regexp parameter defines the regular expression pattern to match against path attribute values.
// The message parameter provides a custom error message when the validation fails.
// The expressions parameter defines the paths to attributes that should be checked against the regex pattern.
// If any expression value matches the regex, the current attribute becomes required.
func IsRequiredIfRegexMatchesPathValue(regexp *regexp.Regexp, message string, expressions ...path.Expression) validator.String {
	return schemavalidator.IsRequiredIfRegexMatchesPathValueValidator{
		Regexp:      regexp,
		Message:     message,
		Expressions: expressions,
	}
}
