// Copyright © 2025 Ping Identity Corporation

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
	_ validator.Bool    = IsRequiredIfMatchesPathValueValidator{}
	_ validator.Float32 = IsRequiredIfMatchesPathValueValidator{}
	_ validator.Float64 = IsRequiredIfMatchesPathValueValidator{}
	_ validator.Int32   = IsRequiredIfMatchesPathValueValidator{}
	_ validator.Int64   = IsRequiredIfMatchesPathValueValidator{}
	_ validator.List    = IsRequiredIfMatchesPathValueValidator{}
	_ validator.Map     = IsRequiredIfMatchesPathValueValidator{}
	_ validator.Number  = IsRequiredIfMatchesPathValueValidator{}
	_ validator.Object  = IsRequiredIfMatchesPathValueValidator{}
	_ validator.Set     = IsRequiredIfMatchesPathValueValidator{}
	_ validator.String  = IsRequiredIfMatchesPathValueValidator{}
)

// IsRequiredIfMatchesPathValueValidator validates if the provided string value equals
// the value at the provided path expression(s).  If matched, the current arguemnt is required.
//
// If a list of expressions is provided, all expressions are checked until a match is found,
// or the list of expressions is exhausted.
type IsRequiredIfMatchesPathValueValidator struct {
	TargetValue basetypes.StringValue
	Expressions path.Expressions
}

type IsRequiredIfMatchesPathValueValidatorRequest struct {
	Config         tfsdk.Config
	ConfigValue    attr.Value
	Path           path.Path
	PathExpression path.Expression
}

type IsRequiredIfMatchesPathValueValidatorResponse struct {
	Diagnostics diag.Diagnostics
}

// Description describes the validation in plain text formatting.
func (v IsRequiredIfMatchesPathValueValidator) Description(_ context.Context) string {
	return fmt.Sprintf("The argument is required if the value %s is present at the defined path: %v", v.TargetValue.ValueString(), v.Expressions)
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v IsRequiredIfMatchesPathValueValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// Validate runs the main validation logic of the validator, reading configuration data out of `req` and updating `resp` with diagnostics.
func (v IsRequiredIfMatchesPathValueValidator) Validate(ctx context.Context, req IsRequiredIfMatchesPathValueValidatorRequest, resp *IsRequiredIfMatchesPathValueValidatorResponse) {
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
			// If a matched value, and the current argument has not been set, return an error.
			if v.TargetValue.Equal(matchedPathValue) && (req.ConfigValue.IsNull()) {

				resp.Diagnostics.AddAttributeError(
					matchedPath,
					"Missing required argument",
					fmt.Sprintf("The argument %s is required because %s is configured as: %s.", req.Path, matchedPath, v.TargetValue),
				)
			}
		}
	}
}

func (validator IsRequiredIfMatchesPathValueValidator) ValidateBool(ctx context.Context, req validator.BoolRequest, resp *validator.BoolResponse) {
	validateReq := IsRequiredIfMatchesPathValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &IsRequiredIfMatchesPathValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (validator IsRequiredIfMatchesPathValueValidator) ValidateFloat32(ctx context.Context, req validator.Float32Request, resp *validator.Float32Response) {
	validateReq := IsRequiredIfMatchesPathValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &IsRequiredIfMatchesPathValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (validator IsRequiredIfMatchesPathValueValidator) ValidateFloat64(ctx context.Context, req validator.Float64Request, resp *validator.Float64Response) {
	validateReq := IsRequiredIfMatchesPathValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &IsRequiredIfMatchesPathValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (validator IsRequiredIfMatchesPathValueValidator) ValidateInt32(ctx context.Context, req validator.Int32Request, resp *validator.Int32Response) {
	validateReq := IsRequiredIfMatchesPathValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &IsRequiredIfMatchesPathValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (validator IsRequiredIfMatchesPathValueValidator) ValidateInt64(ctx context.Context, req validator.Int64Request, resp *validator.Int64Response) {
	validateReq := IsRequiredIfMatchesPathValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &IsRequiredIfMatchesPathValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (validator IsRequiredIfMatchesPathValueValidator) ValidateList(ctx context.Context, req validator.ListRequest, resp *validator.ListResponse) {
	validateReq := IsRequiredIfMatchesPathValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &IsRequiredIfMatchesPathValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (validator IsRequiredIfMatchesPathValueValidator) ValidateMap(ctx context.Context, req validator.MapRequest, resp *validator.MapResponse) {
	validateReq := IsRequiredIfMatchesPathValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &IsRequiredIfMatchesPathValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (validator IsRequiredIfMatchesPathValueValidator) ValidateNumber(ctx context.Context, req validator.NumberRequest, resp *validator.NumberResponse) {
	validateReq := IsRequiredIfMatchesPathValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &IsRequiredIfMatchesPathValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (validator IsRequiredIfMatchesPathValueValidator) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	validateReq := IsRequiredIfMatchesPathValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &IsRequiredIfMatchesPathValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (validator IsRequiredIfMatchesPathValueValidator) ValidateSet(ctx context.Context, req validator.SetRequest, resp *validator.SetResponse) {
	validateReq := IsRequiredIfMatchesPathValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &IsRequiredIfMatchesPathValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (validator IsRequiredIfMatchesPathValueValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	validateReq := IsRequiredIfMatchesPathValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &IsRequiredIfMatchesPathValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}
