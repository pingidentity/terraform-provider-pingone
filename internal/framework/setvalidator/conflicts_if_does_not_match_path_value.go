// Copyright © 2026 Ping Identity Corporation

package setvalidator

import (
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/schemavalidator"
)

func ConflictsIfDoesNotMatchPathValue(targetValue basetypes.StringValue, expressions ...path.Expression) validator.Set {
	return schemavalidator.ConflictsIfDoesNotMatchPathValueValidator{
		TargetValue: targetValue,
		Expressions: expressions,
	}
}
