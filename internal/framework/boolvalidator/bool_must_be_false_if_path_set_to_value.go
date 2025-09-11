// Copyright © 2025 Ping Identity Corporation

package boolvalidator

import (
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// MustBeFalseIfPathSetToValue creates a validator that requires a boolean to be false when another attribute matches a value.
// It returns a validator that ensures the current boolean attribute is false when the specified path attribute
// equals the provided value. This is useful for enforcing conditional constraints between attributes.
//
// The pathAttributeValue parameter specifies the value that the path attribute must match.
// The expression parameter defines the path to the attribute whose value is being checked.
func MustBeFalseIfPathSetToValue(pathAttributeValue basetypes.StringValue, expression path.Expression) validator.Bool {
	return boolMustBeValueIfPathSetToValue{
		BoolValue:          types.BoolValue(false),
		PathAttributeValue: pathAttributeValue,
		PathExpression:     expression,
	}
}
