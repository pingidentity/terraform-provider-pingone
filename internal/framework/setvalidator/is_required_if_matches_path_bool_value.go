// Copyright © 2025 Ping Identity Corporation

package setvalidator

import (
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/schemavalidator"
)

// IsRequiredIfMatchesPathBoolValue validates if the provided boolean value equals
// the value at the provided path expression(s).  If matched, the current argument is required.
//
// If a list of expressions is provided, all expressions are checked until a match is found,
// or the list of expressions is exhausted.
func IsRequiredIfMatchesPathBoolValue(targetValue basetypes.BoolValue, expressions ...path.Expression) validator.Set {
	return schemavalidator.IsRequiredIfMatchesPathBoolValueValidator{
		TargetValue: targetValue,
		Expressions: expressions,
	}
}
