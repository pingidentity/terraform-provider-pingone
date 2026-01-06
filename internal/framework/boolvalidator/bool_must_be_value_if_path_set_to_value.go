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

// boolMustBeValueIfPathSetToValue validates that set contains at least min elements.
type boolMustBeValueIfPathSetToValue struct {
	BoolValue          basetypes.BoolValue
	PathAttributeValue basetypes.StringValue
	PathExpression     path.Expression
}

// Description describes the validation in plain text formatting.
func (v boolMustBeValueIfPathSetToValue) Description(_ context.Context) string {
	return fmt.Sprintf("Ensure that the boolean is set to %v when the value %s is present in the following path expression: %v", v.BoolValue, v.PathAttributeValue, v.PathExpression)
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v boolMustBeValueIfPathSetToValue) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// Validate performs the validation.
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
			fmt.Sprintf("Attribute %s must be %t when attribute %v value is %s", req.Path, v.BoolValue.ValueBool(), v.PathExpression, v.PathAttributeValue.ValueString()),
		))
	}
}

// MustBeValueIfPathSetToValue checks that the boolean is set to the required value if a string value is present in the provided path.Expression.
func MustBeValueIfPathSetToValue(boolValue basetypes.BoolValue, pathAttributeValue basetypes.StringValue, expression path.Expression) validator.Bool {
	return boolMustBeValueIfPathSetToValue{
		BoolValue:          boolValue,
		PathAttributeValue: pathAttributeValue,
		PathExpression:     expression,
	}
}
