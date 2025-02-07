// Copyright © 2025 Ping Identity Corporation

package boolvalidator

import (
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/schemavalidator"
)

func ConflictsIfMatchesPathValue(targetValue basetypes.StringValue, expressions ...path.Expression) validator.Bool {
	return schemavalidator.ConflictsIfMatchesPathValueValidator{
		TargetValue: targetValue,
		Expressions: expressions,
	}
}
