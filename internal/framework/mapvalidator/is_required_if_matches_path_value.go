// Copyright © 2025 Ping Identity Corporation

package mapvalidator

import (
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/schemavalidator"
)

// IsRequiredIfMatchesPathValue validates if the provided string value equals
// the value at the provided path expression(s).  If matched, the current arguemnt is required.
//
// If a list of expressions is provided, all expressions are checked until a match is found,
// or the list of expressions is exhausted.
func IsRequiredIfMatchesPathValue(targetValue basetypes.StringValue, expressions ...path.Expression) validator.Map {
	return schemavalidator.IsRequiredIfMatchesPathValueValidator{
		TargetValue: targetValue,
		Expressions: expressions,
	}
}
