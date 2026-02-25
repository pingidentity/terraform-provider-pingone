// Copyright © 2026 Ping Identity Corporation

package boolvalidator

import (
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func MustBeFalseIfPathSetToValue(pathAttributeValue basetypes.StringValue, expression path.Expression) validator.Bool {
	return boolMustBeValueIfPathSetToValue{
		BoolValue:          types.BoolValue(false),
		PathAttributeValue: pathAttributeValue,
		PathExpression:     expression,
	}
}
