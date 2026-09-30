// Copyright © 2026 Ping Identity Corporation

package boolvalidator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.Bool = boolAlsoRequiresIfTrueValidator{}

// boolAlsoRequiresIfTrueValidator validates that the path expressions have a
// non-null value when the attribute is set to true.
type boolAlsoRequiresIfTrueValidator struct {
	PathExpressions path.Expressions
}

// Description describes the validation in plain text formatting.
func (v boolAlsoRequiresIfTrueValidator) Description(_ context.Context) string {
	return fmt.Sprintf("Ensure that if an attribute is set to true, these are also set: %v", v.PathExpressions)
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v boolAlsoRequiresIfTrueValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// Validate performs the validation.
func (v boolAlsoRequiresIfTrueValidator) ValidateBool(ctx context.Context, req validator.BoolRequest, resp *validator.BoolResponse) {

	// If attribute configuration is null, there is nothing else to validate
	// If attribute configuration is unknown, delay the validation until it is known.
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	// If attribute configuration is not true, there is nothing else to validate
	if !req.ConfigValue.ValueBool() {
		return
	}

	for _, expression := range v.PathExpressions {
		matchedPaths, diags := req.Config.PathMatches(ctx, expression)

		resp.Diagnostics.Append(diags...)

		// Collect all errors
		if diags.HasError() {
			continue
		}

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

			if mpVal.IsNull() {
				resp.Diagnostics.Append(validatordiag.InvalidAttributeCombinationDiagnostic(
					req.Path,
					fmt.Sprintf("Attribute %q must be specified when %q is set to true", mp, req.Path),
				))
			}
		}
	}
}

// AlsoRequiresIfTrue checks that a set of path.Expression has a non-null value
// if the current attribute is set to true.
//
// Relative path.Expression will be resolved using the attribute or block
// being validated.
func AlsoRequiresIfTrue(expressions ...path.Expression) validator.Bool {
	return boolAlsoRequiresIfTrueValidator{
		PathExpressions: expressions,
	}
}
