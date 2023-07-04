package stringvalidator

import (
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/schemavalidator"
)

func ShouldBeDefinedValueIfPathMatchesValue(attributeValue basetypes.StringValue, targetPathValue basetypes.StringValue, expressions ...path.Expression) validator.String {
	return schemavalidator.ShouldBeDefinedValueIfPathMatchesValueValidator{
		AttributeValue:  attributeValue,
		TargetPathValue: targetPathValue,
		Expressions:     expressions,
	}
}
