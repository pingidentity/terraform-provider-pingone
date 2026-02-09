// Copyright © 2026 Ping Identity Corporation

package int32validator

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/schemavalidator"
)

// RegexMatchesPathValue validates if the provided regex matches
// the value at the provided path expression(s).  If a list of expressions is provided,
// all expressions are checked until a match is found, or the list of expressions is exhausted.
func RegexMatchesPathValue(regexp *regexp.Regexp, message string, expressions ...path.Expression) validator.Int32 {
	return schemavalidator.RegexMatchesPathValueValidator{
		Regexp:      regexp,
		Message:     message,
		Expressions: expressions,
	}
}
