package boolvalidator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ validator.Bool = boolAtLeastOneOfMustBeTrueValidator{}

// boolAtLeastOneOfMustBeTrueValidator validates that set contains at least min elements.
type boolAtLeastOneOfMustBeTrueValidator struct {
	PathExpressions path.Expressions
}

// Description describes the validation in plain text formatting.
func (v boolAtLeastOneOfMustBeTrueValidator) Description(_ context.Context) string {
	return fmt.Sprintf("Ensure that at least one attribute is true from the following: ", v.PathExpressions)
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v boolAtLeastOneOfMustBeTrueValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// Validate performs the validation.
func (v boolAtLeastOneOfMustBeTrueValidator) ValidateBool(ctx context.Context, req validator.BoolRequest, resp *validator.BoolResponse) {
	// If attribute configuration is true, it cannot conflict with others
	if req.ConfigValue.ValueBool() {
		return
	}

	expressions := req.PathExpression.MergeExpressions(v.PathExpressions...)

	for _, expression := range expressions {
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

			// If another attribute is true, we're good
			if mpVal.Equal(types.BoolValue(true)) {
				return
			}
		}

	}

	resp.Diagnostics.Append(validatordiag.InvalidAttributeCombinationDiagnostic(
		req.Path,
		v.Description(ctx),
	))
}

// AtLeastOneOfMustBeTrue checks that a set of path.Expression,
// including the attribute the validator is applied to,
// must have a true value.
//
// Relative path.Expression will be resolved using the attribute being
// validated.
func AtLeastOneOfMustBeTrue(expressions ...path.Expression) validator.Bool {
	return boolAtLeastOneOfMustBeTrueValidator{
		PathExpressions: expressions,
	}
}
