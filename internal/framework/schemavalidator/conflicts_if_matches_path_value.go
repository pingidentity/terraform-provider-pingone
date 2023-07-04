package schemavalidator

// Influenced from github.com/hashicorp/terraform-plugin-framework-validators/internal/schemavalidator/*

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var (
	_ validator.Bool    = ConflictsIfMatchesPathValueValidator{}
	_ validator.Float64 = ConflictsIfMatchesPathValueValidator{}
	_ validator.Int64   = ConflictsIfMatchesPathValueValidator{}
	_ validator.List    = ConflictsIfMatchesPathValueValidator{}
	_ validator.Map     = ConflictsIfMatchesPathValueValidator{}
	_ validator.Number  = ConflictsIfMatchesPathValueValidator{}
	_ validator.Object  = ConflictsIfMatchesPathValueValidator{}
	_ validator.Set     = ConflictsIfMatchesPathValueValidator{}
	_ validator.String  = ConflictsIfMatchesPathValueValidator{}
)

// ConflictsIfMatchesPathValueValidator validates that if the provided string value equals
// the value at the provided path expression(s), the attribute with the validator will be in conflict and should produce a plan error.
//
// If a list of expressions is provided, all expressions are checked until a match is found,
// or the list of expressions is exhausted.
type ConflictsIfMatchesPathValueValidator struct {
	TargetValue basetypes.StringValue
	Expressions path.Expressions
}

type ConflictsIfMatchesPathValueValidatorRequest struct {
	Config         tfsdk.Config
	ConfigValue    attr.Value
	Path           path.Path
	PathExpression path.Expression
}

type ConflictsIfMatchesPathValueValidatorResponse struct {
	Diagnostics diag.Diagnostics
}

// Description describes the validation in plain text formatting.
func (v ConflictsIfMatchesPathValueValidator) Description(_ context.Context) string {
	return fmt.Sprintf("The argument cannot be defined if the value \"%s\" is present at the defined path: %v", v.TargetValue.ValueString(), v.Expressions)
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v ConflictsIfMatchesPathValueValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// Validate runs the main validation logic of the validator, reading configuration data out of `req` and updating `resp` with diagnostics.
func (v ConflictsIfMatchesPathValueValidator) Validate(ctx context.Context, req ConflictsIfMatchesPathValueValidatorRequest, resp *ConflictsIfMatchesPathValueValidatorResponse) {

	// If not set then nothing to do.  If it's unknown, then we can't evaluate it's definition.
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
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
			// If a matched value, return an error.
			if v.TargetValue.Equal(matchedPathValue) {

				resp.Diagnostics.AddAttributeError(
					req.Path,
					"Invalid argument combination",
					v.Description(ctx),
				)
			}
		}
	}
}

func (validator ConflictsIfMatchesPathValueValidator) ValidateBool(ctx context.Context, req validator.BoolRequest, resp *validator.BoolResponse) {
	validateReq := ConflictsIfMatchesPathValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &ConflictsIfMatchesPathValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (validator ConflictsIfMatchesPathValueValidator) ValidateFloat64(ctx context.Context, req validator.Float64Request, resp *validator.Float64Response) {
	validateReq := ConflictsIfMatchesPathValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &ConflictsIfMatchesPathValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (validator ConflictsIfMatchesPathValueValidator) ValidateInt64(ctx context.Context, req validator.Int64Request, resp *validator.Int64Response) {
	validateReq := ConflictsIfMatchesPathValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &ConflictsIfMatchesPathValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (validator ConflictsIfMatchesPathValueValidator) ValidateList(ctx context.Context, req validator.ListRequest, resp *validator.ListResponse) {
	validateReq := ConflictsIfMatchesPathValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &ConflictsIfMatchesPathValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (validator ConflictsIfMatchesPathValueValidator) ValidateMap(ctx context.Context, req validator.MapRequest, resp *validator.MapResponse) {
	validateReq := ConflictsIfMatchesPathValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &ConflictsIfMatchesPathValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (validator ConflictsIfMatchesPathValueValidator) ValidateNumber(ctx context.Context, req validator.NumberRequest, resp *validator.NumberResponse) {
	validateReq := ConflictsIfMatchesPathValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &ConflictsIfMatchesPathValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (validator ConflictsIfMatchesPathValueValidator) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	validateReq := ConflictsIfMatchesPathValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &ConflictsIfMatchesPathValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (validator ConflictsIfMatchesPathValueValidator) ValidateSet(ctx context.Context, req validator.SetRequest, resp *validator.SetResponse) {
	validateReq := ConflictsIfMatchesPathValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &ConflictsIfMatchesPathValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (validator ConflictsIfMatchesPathValueValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	validateReq := ConflictsIfMatchesPathValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &ConflictsIfMatchesPathValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}
