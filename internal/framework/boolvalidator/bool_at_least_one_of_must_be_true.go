package boolvalidator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var _ validator.Bool = boolAtLeastOneOfMustBeTrueValidator{}

// boolAtLeastOneOfMustBeTrueValidator validates that set contains at least min elements.
type boolAtLeastOneOfMustBeTrueValidator struct {
	AttributeDefault   basetypes.BoolValue
	ExpressionDefaults basetypes.BoolValue
	PathExpressions    path.Expressions
}

type boolAtLeastOneOfMustBeTrueValidatorExpression struct {
	expressionDefault basetypes.BoolValue
}

// Description describes the validation in plain text formatting.
func (v boolAtLeastOneOfMustBeTrueValidator) Description(_ context.Context) string {
	return fmt.Sprintf("Ensure that at least one attribute is true from the following: %v", v.PathExpressions)
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

	// If it's null and the default is marked as true, then return ok
	if req.ConfigValue.IsNull() {
		if v.AttributeDefault.Equal(types.BoolValue(true)) {
			return
		} else {
			resp.Diagnostics.Append(validatordiag.InvalidAttributeCombinationDiagnostic(
				req.Path,
				v.Description(ctx),
			))
		}
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

			// If another attribute is true, or is null and we have a default, we're good
			if mpVal.Equal(types.BoolValue(true)) || mpVal.IsNull() && v.ExpressionDefaults.Equal(types.BoolValue(true)) {
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
func AtLeastOneOfMustBeTrue(attributeDefault, expressionDefaults basetypes.BoolValue, expressions ...path.Expression) validator.Bool {
	return boolAtLeastOneOfMustBeTrueValidator{
		AttributeDefault:   attributeDefault,
		ExpressionDefaults: expressionDefaults,
		PathExpressions:    expressions,
	}
}
