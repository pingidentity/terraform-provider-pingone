package objectvalidator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// objectIsRequiredIfMatchesPathValueValidator validates if the provided string value equals
// the value at the provided path expression(s).  If matched, the current arguemnt is required.
//
// If a list of expressions is provided, all expressions are checked until a match is found,
// or the list of expressions is exhausted.
type objectIsRequiredIfMatchesPathValueValidator struct {
	targetValue basetypes.StringValue
	expressions path.Expressions
}

// Description describes the validation in plain text formatting.
func (v objectIsRequiredIfMatchesPathValueValidator) Description(_ context.Context) string {
	return fmt.Sprintf("The argument is required if the value %s is present at the defined path: %v", v.targetValue.ValueString(), v.expressions)
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v objectIsRequiredIfMatchesPathValueValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// Validate runs the main validation logic of the validator, reading configuration data out of `req` and updating `resp` with diagnostics.
func (v objectIsRequiredIfMatchesPathValueValidator) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	// Combine the given path expressions with the current attribute path
	// expression. This call automatically handles relative and absolute
	// expressions.
	expressions := req.PathExpression.MergeExpressions(v.expressions...)

	// For each expression, find matching paths.
	for _, expression := range expressions {
		// Find paths matching the expression in the configuration data.
		matchedPaths, diags := req.Config.PathMatches(ctx, expression)

		resp.Diagnostics.Append(diags...)
		if diags.HasError() {
			continue
		}

		// For each matched path, get the data and compare.
		for _, matchedPath := range matchedPaths {
			// Fetch the generic attr.Value at the given path. This ensures any
			// potential parent value of a different type, which can be a null
			// or unknown value, can be safely checked without raising a type
			// conversion error.
			var matchedPathValue attr.Value

			diags := req.Config.GetAttribute(ctx, matchedPath, &matchedPathValue)
			resp.Diagnostics.Append(diags...)
			if diags.HasError() {
				continue
			}

			// If the matched path value is null or unknown, we cannot compare
			// values, so continue to other matched paths.
			if matchedPathValue.IsNull() || matchedPathValue.IsUnknown() {
				continue
			}

			// Found a matched path.  Compare the matched path to the provided path.
			// If a matched value, and the current argument has not been set, return an error.
			if v.targetValue.Equal(matchedPathValue) && (req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown()) {

				resp.Diagnostics.AddAttributeError(
					matchedPath,
					"Missing required argument",
					fmt.Sprintf("The argument %s is required because %s is configured as: %s.", req.Path, matchedPath, v.targetValue),
				)
			}
		}
	}
}

// IsRequiredIfMatchesPathValue validates if the provided string value equals
// the value at the provided path expression(s).  If matched, the current arguemnt is required.
//
// If a list of expressions is provided, all expressions are checked until a match is found,
// or the list of expressions is exhausted.
func IsRequiredIfMatchesPathValue(targetValue basetypes.StringValue, expressions ...path.Expression) validator.Object {
	return &objectIsRequiredIfMatchesPathValueValidator{
		targetValue: targetValue,
		expressions: expressions,
	}
}
