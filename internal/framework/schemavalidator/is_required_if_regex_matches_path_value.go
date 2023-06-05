package schemavalidator

// Influenced from github.com/hashicorp/terraform-plugin-framework-validators/internal/schemavalidator/*

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

var (
	_ validator.Bool    = IsRequiredIfRegexMatchesPathValueValidator{}
	_ validator.Float64 = IsRequiredIfRegexMatchesPathValueValidator{}
	_ validator.Int64   = IsRequiredIfRegexMatchesPathValueValidator{}
	_ validator.List    = IsRequiredIfRegexMatchesPathValueValidator{}
	_ validator.Map     = IsRequiredIfRegexMatchesPathValueValidator{}
	_ validator.Number  = IsRequiredIfRegexMatchesPathValueValidator{}
	_ validator.Object  = IsRequiredIfRegexMatchesPathValueValidator{}
	_ validator.Set     = IsRequiredIfRegexMatchesPathValueValidator{}
	_ validator.String  = IsRequiredIfRegexMatchesPathValueValidator{}
)

// IsRequiredIfRegexMatchesPathValueValidator validates if the provided regex matches
// the value at the provided path expression(s). If matched, the current argument is required.
//
// If a list of expressions is provided, all expressions are checked until a match is found,
// or the list of expressions is exhausted.
type IsRequiredIfRegexMatchesPathValueValidator struct {
	Regexp      *regexp.Regexp
	Message     string
	Expressions path.Expressions
}

type IsRequiredIfRegexMatchesPathValueValidatorRequest struct {
	Config         tfsdk.Config
	ConfigValue    attr.Value
	Path           path.Path
	PathExpression path.Expression
}

type IsRequiredIfRegexMatchesPathValueValidatorResponse struct {
	Diagnostics diag.Diagnostics
}

// Description describes the validation in plain text formatting.
func (v IsRequiredIfRegexMatchesPathValueValidator) Description(_ context.Context) string {
	if v.Message != "" {
		return v.Message
	}
	return fmt.Sprintf("The argument is required if the regular expression %s is present at the defined path: %v", v.Regexp, v.Expressions)
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v IsRequiredIfRegexMatchesPathValueValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// Validate runs the main validation logic of the validator, reading configuration data out of `req` and updating `resp` with diagnostics.
func (v IsRequiredIfRegexMatchesPathValueValidator) Validate(ctx context.Context, req IsRequiredIfRegexMatchesPathValueValidatorRequest, resp *IsRequiredIfRegexMatchesPathValueValidatorResponse) {
	// Combine the given path expressions with the current attribute path
	// expression. This call automatically handles relative and absolute
	// expressions.
	expressions := req.PathExpression.MergeExpressions(v.Expressions...)

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
			if v.Regexp.MatchString(matchedPathValue.String()) && (req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown()) {
				resp.Diagnostics.Append(validatordiag.InvalidAttributeValueMatchDiagnostic(
					req.Path,
					v.Description(ctx),
					fmt.Sprintf("%s is undefined.", req.PathExpression.String()),
				))
			}
		}
	}
}

func (validator IsRequiredIfRegexMatchesPathValueValidator) ValidateBool(ctx context.Context, req validator.BoolRequest, resp *validator.BoolResponse) {
	validateReq := IsRequiredIfRegexMatchesPathValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &IsRequiredIfRegexMatchesPathValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (validator IsRequiredIfRegexMatchesPathValueValidator) ValidateFloat64(ctx context.Context, req validator.Float64Request, resp *validator.Float64Response) {
	validateReq := IsRequiredIfRegexMatchesPathValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &IsRequiredIfRegexMatchesPathValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (validator IsRequiredIfRegexMatchesPathValueValidator) ValidateInt64(ctx context.Context, req validator.Int64Request, resp *validator.Int64Response) {
	validateReq := IsRequiredIfRegexMatchesPathValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &IsRequiredIfRegexMatchesPathValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (validator IsRequiredIfRegexMatchesPathValueValidator) ValidateList(ctx context.Context, req validator.ListRequest, resp *validator.ListResponse) {
	validateReq := IsRequiredIfRegexMatchesPathValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &IsRequiredIfRegexMatchesPathValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (validator IsRequiredIfRegexMatchesPathValueValidator) ValidateMap(ctx context.Context, req validator.MapRequest, resp *validator.MapResponse) {
	validateReq := IsRequiredIfRegexMatchesPathValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &IsRequiredIfRegexMatchesPathValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (validator IsRequiredIfRegexMatchesPathValueValidator) ValidateNumber(ctx context.Context, req validator.NumberRequest, resp *validator.NumberResponse) {
	validateReq := IsRequiredIfRegexMatchesPathValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &IsRequiredIfRegexMatchesPathValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (validator IsRequiredIfRegexMatchesPathValueValidator) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	validateReq := IsRequiredIfRegexMatchesPathValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &IsRequiredIfRegexMatchesPathValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (validator IsRequiredIfRegexMatchesPathValueValidator) ValidateSet(ctx context.Context, req validator.SetRequest, resp *validator.SetResponse) {
	validateReq := IsRequiredIfRegexMatchesPathValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &IsRequiredIfRegexMatchesPathValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (validator IsRequiredIfRegexMatchesPathValueValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	validateReq := IsRequiredIfRegexMatchesPathValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &IsRequiredIfRegexMatchesPathValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}
