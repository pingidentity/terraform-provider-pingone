package credentials

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure our implementation satisfies the validator.Int64 interface.
//var _ validator.Int64 = &int64IsGreaterThanValidator{}

// int64IsGreaterThanValidator is the underlying type implementing Int64IsGreaterThan.
type stringIsRequiredIfPathValueValidator struct {
	targetValue basetypes.StringValue
	expressions path.Expressions
}

// Description returns a plaintext string describing the validator.
func (v stringIsRequiredIfPathValueValidator) Description(_ context.Context) string {
	return fmt.Sprintf("If configured, must be greater than %s attributes", v.expressions)
}

// MarkdownDescription returns a Markdown formatted string describing the validator.
func (v stringIsRequiredIfPathValueValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// Validate performs the validation logic for the validator.
func (v stringIsRequiredIfPathValueValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
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
			// If a matched path, and the current property has not been set, return an error.
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

// Int64IsGreaterThan checks that any Int64 values in the paths described by the
// path.Expression are less than the current attribute value.
func IsRequiredIfPathValue(targetValue basetypes.StringValue, expressions ...path.Expression) validator.String {
	return &stringIsRequiredIfPathValueValidator{
		targetValue: targetValue,
		expressions: expressions,
	}
}
