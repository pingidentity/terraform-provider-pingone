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

var _ validator.Bool = boolMustBeValueIfPathSetToValue{}

// boolMustBeValueIfPathSetToValue validates that a boolean attribute has a specific value when another attribute matches a condition.
// It ensures that when a path attribute equals a specified value, the current boolean attribute
// must be set to the required boolean value.
type boolMustBeValueIfPathSetToValue struct {
	// BoolValue is the required boolean value for the current attribute when the condition is met
	BoolValue basetypes.BoolValue
	// PathAttributeValue is the value that the path attribute must match to trigger the validation
	PathAttributeValue basetypes.StringValue
	// PathExpression defines the path to the attribute whose value is being checked
	PathExpression path.Expression
}

// Description describes the validation in plain text formatting.
// It returns a human-readable description of the validation rule for error messages and documentation.
func (v boolMustBeValueIfPathSetToValue) Description(_ context.Context) string {
	return fmt.Sprintf("Ensure that the boolean is set to %v when the value %s is present in the following path expression: %v", v.BoolValue, v.PathAttributeValue, v.PathExpression)
}

// MarkdownDescription describes the validation in Markdown formatting.
// It returns the same description as Description() but formatted for Markdown documentation.
func (v boolMustBeValueIfPathSetToValue) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateBool performs the validation logic for conditional boolean value requirements.
// It checks that when the path attribute equals the specified value, the current boolean attribute
// is set to the required boolean value. Unknown values are deferred until they become known.
func (v boolMustBeValueIfPathSetToValue) ValidateBool(ctx context.Context, req validator.BoolRequest, resp *validator.BoolResponse) {

	// If unknown, we can't check until it is known
	if req.ConfigValue.IsUnknown() {
		return
	}

	// If attribute configuration is already the value we want, we're done
	if req.ConfigValue.Equal(v.BoolValue) {
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

// MustBeValueIfPathSetToValue creates a validator that enforces a boolean value based on another attribute's value.
// It returns a validator that ensures the current boolean attribute equals the required value
// when the specified path attribute matches the provided value.
//
// The boolValue parameter specifies the required boolean value for the current attribute.
// The pathAttributeValue parameter specifies the value that the path attribute must match.
// The expression parameter defines the path to the attribute whose value is being checked.
func MustBeValueIfPathSetToValue(boolValue basetypes.BoolValue, pathAttributeValue basetypes.StringValue, expression path.Expression) validator.Bool {
	return boolMustBeValueIfPathSetToValue{
		BoolValue:          boolValue,
		PathAttributeValue: pathAttributeValue,
		PathExpression:     expression,
	}
}
