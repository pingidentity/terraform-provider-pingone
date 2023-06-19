package boolvalidator

import (
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func MustNotBeTrueIfPathSetToValue(pathAttributeValue basetypes.StringValue, expression path.Expression) validator.Bool {
	return boolMustNotBeValueIfPathSetToValue{
		BoolValue:          types.BoolValue(true),
		PathAttributeValue: pathAttributeValue,
		PathExpression:     expression,
	}
}
