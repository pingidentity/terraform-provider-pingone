// Copyright © 2025 Ping Identity Corporation

package boolvalidator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var _ validator.Bool = boolMustNotBeValueIfPathSetToValue{}

// boolMustNotBeValueIfPathSetToValue validates that a boolean attribute does not have a specific value when another attribute matches a condition.
// It ensures that when a path attribute equals a specified value, the current boolean attribute
// must not be set to the prohibited boolean value.
type boolMustNotBeValueIfPathSetToValue struct {
	// BoolValue is the prohibited boolean value for the current attribute when the condition is met
	BoolValue basetypes.BoolValue
	// PathAttributeValue is the value that the path attribute must match to trigger the validation
	PathAttributeValue basetypes.StringValue
	// PathExpression defines the path to the attribute whose value is being checked
	PathExpression path.Expression
}

// Description describes the validation in plain text formatting.
// It returns a human-readable description of the validation rule for error messages and documentation.
func (v boolMustNotBeValueIfPathSetToValue) Description(_ context.Context) string {
	return fmt.Sprintf("Ensure that the boolean is not set to %v when the value %s is present in the following path expression: %v", v.BoolValue, v.PathAttributeValue, v.PathExpression)
}

// MarkdownDescription describes the validation in Markdown formatting.
// It returns the same description as Description() but formatted for Markdown documentation.
func (v boolMustNotBeValueIfPathSetToValue) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateBool performs the validation logic for conditional boolean value prohibition.
// It checks that when the path attribute equals the specified value, the current boolean attribute
// is not set to the prohibited boolean value. Unknown values are deferred until they become known.
func (v boolMustNotBeValueIfPathSetToValue) ValidateBool(ctx context.Context, req validator.BoolRequest, resp *validator.BoolResponse) {

	// If the value is not the one we want to validate paths for, this validator is not applicable
	if !req.ConfigValue.Equal(v.BoolValue) {
		return
	}

	matchedPaths, diags := req.Config.PathMatches(ctx, v.PathExpression)

	resp.Diagnostics.Append(diags...)

	// Collect all errors
	if diags.HasError() {
		return
	}

	matchedValue := false

	for _, mp := range matchedPaths {
		// If the user specifies the same attribute this validator is applied to,
		// also as part of the input, skip it
		if mp.Equal(req.Path) {
			continue
		}

		var mpVal attr.Value
		diags := req.Config.GetAttribute(ctx, mp, &mpVal)
		resp.Diagnostics.Append(diags...)

		// Collect all errors
		if diags.HasError() {
			continue
		}

		// Delay validation until all involved attribute have a known value
		if mpVal.IsUnknown() {
			return
		}

		// If the path value is equal to the required value, we're done
		if mpVal.Equal(v.PathAttributeValue) {
			matchedValue = true
			break
		}
	}

	if matchedValue {
		resp.Diagnostics.Append(validatordiag.InvalidAttributeCombinationDiagnostic(
			req.Path,
			v.Description(ctx),
		))
	}
}

// MustNotBeValueIfPathSetToValue creates a validator that prohibits a boolean value based on another attribute's value.
// It returns a validator that ensures the current boolean attribute does not equal the prohibited value
// when the specified path attribute matches the provided value.
//
// The boolValue parameter specifies the prohibited boolean value for the current attribute.
// The pathAttributeValue parameter specifies the value that the path attribute must match.
// The expression parameter defines the path to the attribute whose value is being checked.
func MustNotBeValueIfPathSetToValue(boolValue basetypes.BoolValue, pathAttributeValue basetypes.StringValue, expression path.Expression) validator.Bool {
	return boolMustNotBeValueIfPathSetToValue{
		BoolValue:          boolValue,
		PathAttributeValue: pathAttributeValue,
		PathExpression:     expression,
	}
}
