// Copyright © 2025 Ping Identity Corporation

package stringvalidator

import (
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/schemavalidator"
)

// ConflictsIfMatchesPathValue creates a validator that checks for conflicts when a path attribute matches a specific value.
// It returns a validator that prevents the current string attribute from being set when any of the specified
// path attributes equal the target value. This is useful for enforcing mutual exclusion between attributes.
//
// The targetValue parameter specifies the value that the path attributes are checked against.
// The expressions parameter defines the paths to attributes that should be checked for conflicts.
func ConflictsIfMatchesPathValue(targetValue basetypes.StringValue, expressions ...path.Expression) validator.String {
	return schemavalidator.ConflictsIfMatchesPathValueValidator{
		TargetValue: targetValue,
		Expressions: expressions,
	}
}
