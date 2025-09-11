// Copyright © 2025 Ping Identity Corporation

package stringvalidator

import (
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/schemavalidator"
)

// ShouldBeDefinedValueIfPathMatchesValue creates a validator that requires a specific value when another attribute matches a target value.
// It returns a validator that ensures the current string attribute equals the specified value
// when any of the specified path attributes match the target value.
//
// The attributeValue parameter specifies the required value for the current attribute.
// The targetPathValue parameter specifies the value that the path attributes must match.
// The expressions parameter defines the paths to attributes that should be checked for the target value.
func ShouldBeDefinedValueIfPathMatchesValue(attributeValue basetypes.StringValue, targetPathValue basetypes.StringValue, expressions ...path.Expression) validator.String {
	return schemavalidator.ShouldBeDefinedValueIfPathMatchesValueValidator{
		AttributeValue:  attributeValue,
		TargetPathValue: targetPathValue,
		Expressions:     expressions,
	}
}
