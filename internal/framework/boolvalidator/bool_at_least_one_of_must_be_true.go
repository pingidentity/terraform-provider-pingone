// Copyright © 2025 Ping Identity Corporation

// Package boolvalidator provides custom boolean validators for the Terraform Plugin Framework.
// This package contains validators that check boolean attribute constraints and combinations
// for the PingOne provider's specific validation requirements.
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

// boolAtLeastOneOfMustBeTrueValidator validates that at least one boolean attribute in a group must be true.
// It ensures that among a set of boolean attributes (including the one being validated),
// at least one must have a true value, considering both explicit values and defaults.
type boolAtLeastOneOfMustBeTrueValidator struct {
	// AttributeDefault is the default value for the attribute being validated
	AttributeDefault basetypes.BoolValue
	// ExpressionDefaults is the default value for attributes in the path expressions
	ExpressionDefaults basetypes.BoolValue
	// PathExpressions contains the paths to other attributes that are part of this validation group
	PathExpressions path.Expressions
}

// Description describes the validation in plain text formatting.
// It returns a human-readable description of the validation rule for error messages and documentation.
func (v boolAtLeastOneOfMustBeTrueValidator) Description(_ context.Context) string {
	return fmt.Sprintf("Ensure that at least one attribute is true from the following: %v", v.PathExpressions)
}

// MarkdownDescription describes the validation in Markdown formatting.
// It returns the same description as Description() but formatted for Markdown documentation.
func (v boolAtLeastOneOfMustBeTrueValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateBool performs the validation logic for the at-least-one-true constraint.
// It checks that among the current attribute and the specified path expressions,
// at least one boolean attribute has a true value (either explicitly set or via defaults).
// The validation considers null values with defaults as potentially true values.
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

// AtLeastOneOfMustBeTrue creates a validator that ensures at least one boolean in a group is true.
// It returns a validator that checks the current attribute and specified path expressions
// to ensure at least one has a true value, considering both explicit values and defaults.
//
// The attributeDefault parameter specifies the default value for the current attribute being validated.
// The expressionDefaults parameter specifies the default value for attributes in the path expressions.
// The expressions parameter contains the paths to other attributes that are part of this validation group.
//
// Relative path expressions will be resolved using the attribute being validated as the base.
func AtLeastOneOfMustBeTrue(attributeDefault, expressionDefaults basetypes.BoolValue, expressions ...path.Expression) validator.Bool {
	return boolAtLeastOneOfMustBeTrueValidator{
		AttributeDefault:   attributeDefault,
		ExpressionDefaults: expressionDefaults,
		PathExpressions:    expressions,
	}
}
