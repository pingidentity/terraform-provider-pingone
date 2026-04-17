// Copyright © 2026 Ping Identity Corporation

package listvalidator

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/schemavalidator"
)

// IsRequiredIfMatchesPathValue validates if the provided attribute value equals
// the value at the provided path expression(s).  If matched, the current argument is required.
//
// If a list of expressions is provided, all expressions are checked until a match is found,
// or the list of expressions is exhausted.
func IsRequiredIfMatchesPathValue(targetValue attr.Value, expressions ...path.Expression) validator.List {
	return schemavalidator.IsRequiredIfMatchesPathValueValidator{
		TargetValue: targetValue,
		Expressions: expressions,
	}
}
