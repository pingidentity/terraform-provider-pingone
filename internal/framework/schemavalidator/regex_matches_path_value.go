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
	_ validator.Bool    = RegexMatchesPathValueValidator{}
	_ validator.Float64 = RegexMatchesPathValueValidator{}
	_ validator.Int64   = RegexMatchesPathValueValidator{}
	_ validator.List    = RegexMatchesPathValueValidator{}
	_ validator.Map     = RegexMatchesPathValueValidator{}
	_ validator.Number  = RegexMatchesPathValueValidator{}
	_ validator.Object  = RegexMatchesPathValueValidator{}
	_ validator.Set     = RegexMatchesPathValueValidator{}
	_ validator.String  = RegexMatchesPathValueValidator{}
)

// RegexMatchesPathValueValidator validates if the provided regex matches
// the value at the provided path expression(s).  If a list of expressions is provided,
// all expressions are checked until a match is found, or the list of expressions is exhausted.
type RegexMatchesPathValueValidator struct {
	Regexp      *regexp.Regexp
	Message     string
	Expressions path.Expressions
}

type RegexMatchesPathValueValidatorRequest struct {
	Config         tfsdk.Config
	ConfigValue    attr.Value
	Path           path.Path
	PathExpression path.Expression
}

type RegexMatchesPathValueValidatorResponse struct {
	Diagnostics diag.Diagnostics
}

// Description describes the validation in plain text formatting.
func (v RegexMatchesPathValueValidator) Description(_ context.Context) string {
	if v.Message != "" {
		return v.Message
	}
	return fmt.Sprintf("The value at path %v must match regular expression '%s'", v.Expressions, v.Regexp)
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v RegexMatchesPathValueValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// Validate runs the main validation logic of the validator, reading configuration data out of `req` and updating `resp` with diagnostics.
func (v RegexMatchesPathValueValidator) Validate(ctx context.Context, req RegexMatchesPathValueValidatorRequest, resp *RegexMatchesPathValueValidatorResponse) {
	// If the value is null, there is nothing to validate. This validator
	// is only concerned if the source value has been set. The value's
	// content is not used in the validation decision.
	if req.ConfigValue.IsNull() {
		return
	}

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
			// If there is not a regex match, return the provided error message.
			if !v.Regexp.MatchString(matchedPathValue.String()) {
				resp.Diagnostics.Append(validatordiag.InvalidAttributeValueMatchDiagnostic(
					req.Path,
					v.Description(ctx),
					matchedPathValue.String(),
				))
			}
		}
	}
}

func (validator RegexMatchesPathValueValidator) ValidateBool(ctx context.Context, req validator.BoolRequest, resp *validator.BoolResponse) {
	validateReq := RegexMatchesPathValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &RegexMatchesPathValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (validator RegexMatchesPathValueValidator) ValidateFloat64(ctx context.Context, req validator.Float64Request, resp *validator.Float64Response) {
	validateReq := RegexMatchesPathValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &RegexMatchesPathValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (validator RegexMatchesPathValueValidator) ValidateInt64(ctx context.Context, req validator.Int64Request, resp *validator.Int64Response) {
	validateReq := RegexMatchesPathValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &RegexMatchesPathValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (validator RegexMatchesPathValueValidator) ValidateList(ctx context.Context, req validator.ListRequest, resp *validator.ListResponse) {
	validateReq := RegexMatchesPathValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &RegexMatchesPathValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (validator RegexMatchesPathValueValidator) ValidateMap(ctx context.Context, req validator.MapRequest, resp *validator.MapResponse) {
	validateReq := RegexMatchesPathValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &RegexMatchesPathValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (validator RegexMatchesPathValueValidator) ValidateNumber(ctx context.Context, req validator.NumberRequest, resp *validator.NumberResponse) {
	validateReq := RegexMatchesPathValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &RegexMatchesPathValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (validator RegexMatchesPathValueValidator) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	validateReq := RegexMatchesPathValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &RegexMatchesPathValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (validator RegexMatchesPathValueValidator) ValidateSet(ctx context.Context, req validator.SetRequest, resp *validator.SetResponse) {
	validateReq := RegexMatchesPathValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &RegexMatchesPathValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (validator RegexMatchesPathValueValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	validateReq := RegexMatchesPathValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &RegexMatchesPathValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}
