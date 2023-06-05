package boolvalidator

import (
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// MustBeTrueIfPathSetToValue checks that the boolean is set to the required value if a string value is present in the provided path.Expression.
func MustBeFalseIfPathSetToValue(pathAttributeValue basetypes.StringValue, expression path.Expression) validator.Bool {
	return boolMustBeValueIfPathSetToValue{
		BoolValue:          types.BoolValue(false),
		PathAttributeValue: pathAttributeValue,
		PathExpression:     expression,
	}
}
