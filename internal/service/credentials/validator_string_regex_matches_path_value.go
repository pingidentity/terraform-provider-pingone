package credentials

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Ensure our implementation satisfies the validator.Int64 interface.
//var _ validator.Int64 = &int64IsGreaterThanValidator{}

// int64IsGreaterThanValidator is the underlying type implementing Int64IsGreaterThan.
type stringRegexMatchesPathValueValidator struct {
	regexp      *regexp.Regexp
	message     string
	expressions path.Expressions
}

// Description returns a plaintext string describing the validator.
func (v stringRegexMatchesPathValueValidator) Description(_ context.Context) string {
	if v.message != "" {
		return v.message
	}
	return fmt.Sprintf("value must match regular expression '%s'", v.regexp)
}

// MarkdownDescription returns a Markdown formatted string describing the validator.
func (v stringRegexMatchesPathValueValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// Validate performs the validation logic for the validator.
func (v stringRegexMatchesPathValueValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	// If the value is unknown or null, there is nothing to validate.
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

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

			if !v.regexp.MatchString(matchedPathValue.String()) {
				resp.Diagnostics.Append(validatordiag.InvalidAttributeValueMatchDiagnostic(
					req.Path,
					v.Description(ctx),
					matchedPathValue.String(),
				))
			}
		}
	}
}

// Int64IsGreaterThan checks that any Int64 values in the paths described by the
// path.Expression are less than the current attribute value.
func RegexMatchesPathValue(regexp *regexp.Regexp, message string, expressions ...path.Expression) validator.String {
	return &stringRegexMatchesPathValueValidator{
		regexp:      regexp,
		message:     message,
		expressions: expressions,
	}
}
