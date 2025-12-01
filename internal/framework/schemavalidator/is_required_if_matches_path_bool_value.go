// Copyright © 2025 Ping Identity Corporation

package schemavalidator

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
	_ validator.Bool    = IsRequiredIfMatchesPathBoolValueValidator{}
	_ validator.Float32 = IsRequiredIfMatchesPathBoolValueValidator{}
	_ validator.Float64 = IsRequiredIfMatchesPathBoolValueValidator{}
	_ validator.Int32   = IsRequiredIfMatchesPathBoolValueValidator{}
	_ validator.Int64   = IsRequiredIfMatchesPathBoolValueValidator{}
	_ validator.List    = IsRequiredIfMatchesPathBoolValueValidator{}
	_ validator.Map     = IsRequiredIfMatchesPathBoolValueValidator{}
	_ validator.Number  = IsRequiredIfMatchesPathBoolValueValidator{}
	_ validator.Object  = IsRequiredIfMatchesPathBoolValueValidator{}
	_ validator.Set     = IsRequiredIfMatchesPathBoolValueValidator{}
	_ validator.String  = IsRequiredIfMatchesPathBoolValueValidator{}
)

// IsRequiredIfMatchesPathBoolValueValidator validates if the provided boolean value equals
// the value at the provided path expression(s).  If matched, the current argument is required.
//
// If a list of expressions is provided, all expressions are checked until a match is found,
// or the list of expressions is exhausted.
type IsRequiredIfMatchesPathBoolValueValidator struct {
	TargetValue basetypes.BoolValue
	Expressions path.Expressions
}

type IsRequiredIfMatchesPathBoolValueValidatorRequest struct {
	Config         tfsdk.Config
	ConfigValue    attr.Value
	Path           path.Path
	PathExpression path.Expression
}

type IsRequiredIfMatchesPathBoolValueValidatorResponse struct {
	Diagnostics diag.Diagnostics
}

func (v IsRequiredIfMatchesPathBoolValueValidator) Description(_ context.Context) string {
	return fmt.Sprintf("If the value at the provided path matches %s, the current attribute is required.", v.TargetValue)
}

func (v IsRequiredIfMatchesPathBoolValueValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v IsRequiredIfMatchesPathBoolValueValidator) Validate(ctx context.Context, req IsRequiredIfMatchesPathBoolValueValidatorRequest, resp *IsRequiredIfMatchesPathBoolValueValidatorResponse) {
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

			// Convert to BoolValue for comparison
			boolValue, ok := matchedPathValue.(basetypes.BoolValue)
			if !ok {
				continue
			}

			// Found a matched path. Compare the matched path to the provided path.
			// If a matched value, and the current argument has not been set, return an error.
			if v.TargetValue.Equal(boolValue) && (req.ConfigValue.IsNull()) {

				resp.Diagnostics.AddAttributeError(
					matchedPath,
					"Missing required argument",
					fmt.Sprintf("The argument %s is required because %s is configured as: %t.", req.Path, matchedPath, v.TargetValue.ValueBool()),
				)
			}
		}
	}
}

func (validator IsRequiredIfMatchesPathBoolValueValidator) ValidateBool(ctx context.Context, req validator.BoolRequest, resp *validator.BoolResponse) {
	validateReq := IsRequiredIfMatchesPathBoolValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &IsRequiredIfMatchesPathBoolValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (validator IsRequiredIfMatchesPathBoolValueValidator) ValidateFloat32(ctx context.Context, req validator.Float32Request, resp *validator.Float32Response) {
	validateReq := IsRequiredIfMatchesPathBoolValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &IsRequiredIfMatchesPathBoolValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (validator IsRequiredIfMatchesPathBoolValueValidator) ValidateFloat64(ctx context.Context, req validator.Float64Request, resp *validator.Float64Response) {
	validateReq := IsRequiredIfMatchesPathBoolValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &IsRequiredIfMatchesPathBoolValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (validator IsRequiredIfMatchesPathBoolValueValidator) ValidateInt32(ctx context.Context, req validator.Int32Request, resp *validator.Int32Response) {
	validateReq := IsRequiredIfMatchesPathBoolValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &IsRequiredIfMatchesPathBoolValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (validator IsRequiredIfMatchesPathBoolValueValidator) ValidateInt64(ctx context.Context, req validator.Int64Request, resp *validator.Int64Response) {
	validateReq := IsRequiredIfMatchesPathBoolValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &IsRequiredIfMatchesPathBoolValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (validator IsRequiredIfMatchesPathBoolValueValidator) ValidateList(ctx context.Context, req validator.ListRequest, resp *validator.ListResponse) {
	validateReq := IsRequiredIfMatchesPathBoolValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &IsRequiredIfMatchesPathBoolValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (validator IsRequiredIfMatchesPathBoolValueValidator) ValidateMap(ctx context.Context, req validator.MapRequest, resp *validator.MapResponse) {
	validateReq := IsRequiredIfMatchesPathBoolValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &IsRequiredIfMatchesPathBoolValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (validator IsRequiredIfMatchesPathBoolValueValidator) ValidateNumber(ctx context.Context, req validator.NumberRequest, resp *validator.NumberResponse) {
	validateReq := IsRequiredIfMatchesPathBoolValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &IsRequiredIfMatchesPathBoolValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (validator IsRequiredIfMatchesPathBoolValueValidator) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	validateReq := IsRequiredIfMatchesPathBoolValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &IsRequiredIfMatchesPathBoolValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (validator IsRequiredIfMatchesPathBoolValueValidator) ValidateSet(ctx context.Context, req validator.SetRequest, resp *validator.SetResponse) {
	validateReq := IsRequiredIfMatchesPathBoolValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &IsRequiredIfMatchesPathBoolValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (validator IsRequiredIfMatchesPathBoolValueValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	validateReq := IsRequiredIfMatchesPathBoolValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &IsRequiredIfMatchesPathBoolValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}
