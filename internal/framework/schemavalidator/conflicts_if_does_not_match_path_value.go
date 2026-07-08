// Copyright © 2026 Ping Identity Corporation

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
	_ validator.Bool    = ConflictsIfDoesNotMatchPathValueValidator{}
	_ validator.Float32 = ConflictsIfDoesNotMatchPathValueValidator{}
	_ validator.Float64 = ConflictsIfDoesNotMatchPathValueValidator{}
	_ validator.Int32   = ConflictsIfDoesNotMatchPathValueValidator{}
	_ validator.Int64   = ConflictsIfDoesNotMatchPathValueValidator{}
	_ validator.List    = ConflictsIfDoesNotMatchPathValueValidator{}
	_ validator.Map     = ConflictsIfDoesNotMatchPathValueValidator{}
	_ validator.Number  = ConflictsIfDoesNotMatchPathValueValidator{}
	_ validator.Object  = ConflictsIfDoesNotMatchPathValueValidator{}
	_ validator.Set     = ConflictsIfDoesNotMatchPathValueValidator{}
	_ validator.String  = ConflictsIfDoesNotMatchPathValueValidator{}
)

// ConflictsIfDoesNotMatchPathValueValidator validates that, if the attribute with the validator is
// defined, the value at the provided path expression(s) must equal the target value. If a non-null,
// non-target value is found at any of the provided paths, the attribute is in conflict and a plan
// error is produced.
//
// This is the inverse of ConflictsIfMatchesPathValueValidator: use it to permit an attribute only
// when a sibling attribute holds a specific value (for example, allow `verify_policy_id` only when
// `action` is `VERIFY`).
//
// If a list of expressions is provided, all expressions are checked until a conflicting value is
// found, or the list of expressions is exhausted.
type ConflictsIfDoesNotMatchPathValueValidator struct {
	TargetValue basetypes.StringValue
	Expressions path.Expressions
}

type ConflictsIfDoesNotMatchPathValueValidatorRequest struct {
	Config         tfsdk.Config
	ConfigValue    attr.Value
	Path           path.Path
	PathExpression path.Expression
}

type ConflictsIfDoesNotMatchPathValueValidatorResponse struct {
	Diagnostics diag.Diagnostics
}

// Description describes the validation in plain text formatting.
func (v ConflictsIfDoesNotMatchPathValueValidator) Description(_ context.Context) string {
	return fmt.Sprintf("The argument cannot be defined unless the value \"%s\" is present at the defined path: %v", v.TargetValue.ValueString(), v.Expressions)
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v ConflictsIfDoesNotMatchPathValueValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// Validate runs the main validation logic of the validator, reading configuration data out of `req` and updating `resp` with diagnostics.
func (v ConflictsIfDoesNotMatchPathValueValidator) Validate(ctx context.Context, req ConflictsIfDoesNotMatchPathValueValidatorRequest, resp *ConflictsIfDoesNotMatchPathValueValidatorResponse) {

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

			// Found a matched path with a known value.  If it does not equal the
			// target value, the attribute is in conflict and an error is returned.
			if !v.TargetValue.Equal(matchedPathValue) {

				resp.Diagnostics.AddAttributeError(
					req.Path,
					"Invalid argument combination",
					v.Description(ctx),
				)
			}
		}
	}
}

func (validator ConflictsIfDoesNotMatchPathValueValidator) ValidateBool(ctx context.Context, req validator.BoolRequest, resp *validator.BoolResponse) {
	validateReq := ConflictsIfDoesNotMatchPathValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &ConflictsIfDoesNotMatchPathValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (validator ConflictsIfDoesNotMatchPathValueValidator) ValidateFloat32(ctx context.Context, req validator.Float32Request, resp *validator.Float32Response) {
	validateReq := ConflictsIfDoesNotMatchPathValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &ConflictsIfDoesNotMatchPathValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (validator ConflictsIfDoesNotMatchPathValueValidator) ValidateFloat64(ctx context.Context, req validator.Float64Request, resp *validator.Float64Response) {
	validateReq := ConflictsIfDoesNotMatchPathValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &ConflictsIfDoesNotMatchPathValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (validator ConflictsIfDoesNotMatchPathValueValidator) ValidateInt32(ctx context.Context, req validator.Int32Request, resp *validator.Int32Response) {
	validateReq := ConflictsIfDoesNotMatchPathValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &ConflictsIfDoesNotMatchPathValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (validator ConflictsIfDoesNotMatchPathValueValidator) ValidateInt64(ctx context.Context, req validator.Int64Request, resp *validator.Int64Response) {
	validateReq := ConflictsIfDoesNotMatchPathValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &ConflictsIfDoesNotMatchPathValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (validator ConflictsIfDoesNotMatchPathValueValidator) ValidateList(ctx context.Context, req validator.ListRequest, resp *validator.ListResponse) {
	validateReq := ConflictsIfDoesNotMatchPathValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &ConflictsIfDoesNotMatchPathValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (validator ConflictsIfDoesNotMatchPathValueValidator) ValidateMap(ctx context.Context, req validator.MapRequest, resp *validator.MapResponse) {
	validateReq := ConflictsIfDoesNotMatchPathValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &ConflictsIfDoesNotMatchPathValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (validator ConflictsIfDoesNotMatchPathValueValidator) ValidateNumber(ctx context.Context, req validator.NumberRequest, resp *validator.NumberResponse) {
	validateReq := ConflictsIfDoesNotMatchPathValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &ConflictsIfDoesNotMatchPathValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (validator ConflictsIfDoesNotMatchPathValueValidator) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	validateReq := ConflictsIfDoesNotMatchPathValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &ConflictsIfDoesNotMatchPathValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (validator ConflictsIfDoesNotMatchPathValueValidator) ValidateSet(ctx context.Context, req validator.SetRequest, resp *validator.SetResponse) {
	validateReq := ConflictsIfDoesNotMatchPathValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &ConflictsIfDoesNotMatchPathValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (validator ConflictsIfDoesNotMatchPathValueValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	validateReq := ConflictsIfDoesNotMatchPathValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &ConflictsIfDoesNotMatchPathValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}
