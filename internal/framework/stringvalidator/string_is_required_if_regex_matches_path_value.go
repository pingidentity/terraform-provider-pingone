package stringvalidator

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// stringIsRequiredIfRegexMatchesPathValueValidator validates if the provided regex matches
// the value at the provided path expression(s). If matched, the current argument is required.
//
// If a list of expressions is provided, all expressions are checked until a match is found,
// or the list of expressions is exhausted.
type stringIsRequiredIfRegexMatchesPathValueValidator struct {
	regexp      *regexp.Regexp
	message     string
	expressions path.Expressions
}

// Description describes the validation in plain text formatting.
func (v stringIsRequiredIfRegexMatchesPathValueValidator) Description(_ context.Context) string {
	if v.message != "" {
		return v.message
	}
	return fmt.Sprintf("The argument is required if the regular expression %s is present at the defined path: %v", v.regexp, v.expressions)
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v stringIsRequiredIfRegexMatchesPathValueValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// Validate runs the main validation logic of the validator, reading configuration data out of `req` and updating `resp` with diagnostics.
func (v stringIsRequiredIfRegexMatchesPathValueValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
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
			// If a regex match, and the current argument has not been set, return an error.
			if v.regexp.MatchString(matchedPathValue.String()) && (req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown()) {
				resp.Diagnostics.Append(validatordiag.InvalidAttributeValueMatchDiagnostic(
					req.Path,
					v.Description(ctx),
					fmt.Sprintf("%s is undefined.", req.PathExpression.String()),
				))
			}
		}
	}
}

// IsRequiredIfRegexMatchesPathValue validates if the provided regex matches
// the value at the provided path expression(s). If matched, the current argument is required.
//
// If a list of expressions is provided, all expressions are checked until a match is found,
// or the list of expressions is exhausted.
func IsRequiredIfRegexMatchesPathValue(regexp *regexp.Regexp, message string, expressions ...path.Expression) validator.String {
	return &stringIsRequiredIfRegexMatchesPathValueValidator{
		regexp:      regexp,
		message:     message,
		expressions: expressions,
	}
}
