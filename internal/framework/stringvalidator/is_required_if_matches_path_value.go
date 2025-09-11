// Copyright © 2025 Ping Identity Corporation

package stringvalidator

import (
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/schemavalidator"
)

// IsRequiredIfMatchesPathValue creates a validator that requires the current attribute when another attribute matches a specific value.
// It returns a validator that ensures the current string attribute is not null or empty
// when any of the specified path attributes equal the target value.
//
// The targetValue parameter specifies the value that the path attributes are checked against.
// The expressions parameter defines the paths to attributes that should be checked for the target value.
// If any expression matches the target value, the current attribute becomes required.
func IsRequiredIfMatchesPathValue(targetValue basetypes.StringValue, expressions ...path.Expression) validator.String {
	return schemavalidator.IsRequiredIfMatchesPathValueValidator{
		TargetValue: targetValue,
		Expressions: expressions,
	}
}
