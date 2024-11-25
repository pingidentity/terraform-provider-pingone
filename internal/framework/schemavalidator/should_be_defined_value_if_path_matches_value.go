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
	_ validator.Bool    = ShouldBeDefinedValueIfPathMatchesValueValidator{}
	_ validator.Float32 = ShouldBeDefinedValueIfPathMatchesValueValidator{}
	_ validator.Float64 = ShouldBeDefinedValueIfPathMatchesValueValidator{}
	_ validator.Int32   = ShouldBeDefinedValueIfPathMatchesValueValidator{}
	_ validator.Int64   = ShouldBeDefinedValueIfPathMatchesValueValidator{}
	_ validator.List    = ShouldBeDefinedValueIfPathMatchesValueValidator{}
	_ validator.Map     = ShouldBeDefinedValueIfPathMatchesValueValidator{}
	_ validator.Number  = ShouldBeDefinedValueIfPathMatchesValueValidator{}
	_ validator.Object  = ShouldBeDefinedValueIfPathMatchesValueValidator{}
	_ validator.Set     = ShouldBeDefinedValueIfPathMatchesValueValidator{}
	_ validator.String  = ShouldBeDefinedValueIfPathMatchesValueValidator{}
)

// ShouldBeDefinedValueIfPathMatchesValueValidator validates that the attribute to which this validator is configured should be a specific value if the given path matches the given value.
//
// If a list of expressions is provided, all expressions are checked until a match is found,
// or the list of expressions is exhausted.
type ShouldBeDefinedValueIfPathMatchesValueValidator struct {
	AttributeValue  basetypes.StringValue
	TargetPathValue basetypes.StringValue
	Expressions     path.Expressions
}

type ShouldBeDefinedValueIfPathMatchesValueValidatorRequest struct {
	Config         tfsdk.Config
	ConfigValue    attr.Value
	Path           path.Path
	PathExpression path.Expression
}

type ShouldBeDefinedValueIfPathMatchesValueValidatorResponse struct {
	Diagnostics diag.Diagnostics
}

// Description describes the validation in plain text formatting.
func (v ShouldBeDefinedValueIfPathMatchesValueValidator) Description(_ context.Context) string {
	return fmt.Sprintf("The argument must be value \"%s\" if value \"%s\" is present at the defined path: %v", v.AttributeValue.ValueString(), v.TargetPathValue.ValueString(), v.Expressions)
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v ShouldBeDefinedValueIfPathMatchesValueValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// Validate runs the main validation logic of the validator, reading configuration data out of `req` and updating `resp` with diagnostics.
func (v ShouldBeDefinedValueIfPathMatchesValueValidator) Validate(ctx context.Context, req ShouldBeDefinedValueIfPathMatchesValueValidatorRequest, resp *ShouldBeDefinedValueIfPathMatchesValueValidatorResponse) {

	// If it's unknown, then we can't evaluate it's definition.
	if req.ConfigValue.IsUnknown() {
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

			// If the matched path value is unknown, we cannot compare
			// values, so continue to other matched paths.
			if matchedPathValue.IsUnknown() {
				continue
			}

			// Found a matched path.  Compare the matched path to the provided path.
			// If a matched value, return an error.
			if v.TargetPathValue.Equal(matchedPathValue) && !req.ConfigValue.Equal(v.AttributeValue) {

				resp.Diagnostics.AddAttributeError(
					req.Path,
					"Invalid argument combination",
					v.Description(ctx),
				)
			}
		}
	}
}

func (validator ShouldBeDefinedValueIfPathMatchesValueValidator) ValidateBool(ctx context.Context, req validator.BoolRequest, resp *validator.BoolResponse) {
	validateReq := ShouldBeDefinedValueIfPathMatchesValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &ShouldBeDefinedValueIfPathMatchesValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (validator ShouldBeDefinedValueIfPathMatchesValueValidator) ValidateFloat32(ctx context.Context, req validator.Float32Request, resp *validator.Float32Response) {
	validateReq := ShouldBeDefinedValueIfPathMatchesValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &ShouldBeDefinedValueIfPathMatchesValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (validator ShouldBeDefinedValueIfPathMatchesValueValidator) ValidateFloat64(ctx context.Context, req validator.Float64Request, resp *validator.Float64Response) {
	validateReq := ShouldBeDefinedValueIfPathMatchesValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &ShouldBeDefinedValueIfPathMatchesValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (validator ShouldBeDefinedValueIfPathMatchesValueValidator) ValidateInt32(ctx context.Context, req validator.Int32Request, resp *validator.Int32Response) {
	validateReq := ShouldBeDefinedValueIfPathMatchesValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &ShouldBeDefinedValueIfPathMatchesValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (validator ShouldBeDefinedValueIfPathMatchesValueValidator) ValidateInt64(ctx context.Context, req validator.Int64Request, resp *validator.Int64Response) {
	validateReq := ShouldBeDefinedValueIfPathMatchesValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &ShouldBeDefinedValueIfPathMatchesValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (validator ShouldBeDefinedValueIfPathMatchesValueValidator) ValidateList(ctx context.Context, req validator.ListRequest, resp *validator.ListResponse) {
	validateReq := ShouldBeDefinedValueIfPathMatchesValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &ShouldBeDefinedValueIfPathMatchesValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (validator ShouldBeDefinedValueIfPathMatchesValueValidator) ValidateMap(ctx context.Context, req validator.MapRequest, resp *validator.MapResponse) {
	validateReq := ShouldBeDefinedValueIfPathMatchesValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &ShouldBeDefinedValueIfPathMatchesValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (validator ShouldBeDefinedValueIfPathMatchesValueValidator) ValidateNumber(ctx context.Context, req validator.NumberRequest, resp *validator.NumberResponse) {
	validateReq := ShouldBeDefinedValueIfPathMatchesValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &ShouldBeDefinedValueIfPathMatchesValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (validator ShouldBeDefinedValueIfPathMatchesValueValidator) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	validateReq := ShouldBeDefinedValueIfPathMatchesValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &ShouldBeDefinedValueIfPathMatchesValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (validator ShouldBeDefinedValueIfPathMatchesValueValidator) ValidateSet(ctx context.Context, req validator.SetRequest, resp *validator.SetResponse) {
	validateReq := ShouldBeDefinedValueIfPathMatchesValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &ShouldBeDefinedValueIfPathMatchesValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (validator ShouldBeDefinedValueIfPathMatchesValueValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	validateReq := ShouldBeDefinedValueIfPathMatchesValueValidatorRequest{
		Config:         req.Config,
		ConfigValue:    req.ConfigValue,
		Path:           req.Path,
		PathExpression: req.PathExpression,
	}
	validateResp := &ShouldBeDefinedValueIfPathMatchesValueValidatorResponse{}

	validator.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}
