// Copyright © 2025 Ping Identity Corporation

package stringvalidator

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/schemavalidator"
)

// RegexMatchesPathValue creates a validator that checks if the current attribute value matches a regex when path values match.
// It returns a validator that ensures the current string attribute matches the provided regular expression
// when any of the specified path attributes are present and have values.
//
// The regexp parameter defines the regular expression pattern that the current attribute value must match.
// The message parameter provides a custom error message when the validation fails.
// The expressions parameter defines the paths to attributes that trigger this regex validation.
// If any expression has a value, the current attribute must match the regex pattern.
func RegexMatchesPathValue(regexp *regexp.Regexp, message string, expressions ...path.Expression) validator.String {
	return schemavalidator.RegexMatchesPathValueValidator{
		Regexp:      regexp,
		Message:     message,
		Expressions: expressions,
	}
}
