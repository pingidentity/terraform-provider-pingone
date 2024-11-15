package int32validator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// int32IsGreaterThanPathValue validates if the input string is base64encoded.
type int32IsGreaterThanPathValue struct {
	PathExpression path.Expression
}

// Description describes the validation in plain text formatting.
func (v int32IsGreaterThanPathValue) Description(ctx context.Context) string {
	return fmt.Sprintf("Ensure the integer is greater than the values configured at paths that match the following path expression: %v", v.PathExpression)
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v int32IsGreaterThanPathValue) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// Validate runs the main validation logic of the validator, reading configuration data out of `req` and updating `resp` with diagnostics.
func (v int32IsGreaterThanPathValue) ValidateInt32(ctx context.Context, req validator.Int32Request, resp *validator.Int32Response) {
	// If the value is unknown or null, there is nothing to validate.
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	matchedPaths, d := req.Config.PathMatches(ctx, v.PathExpression)
	resp.Diagnostics.Append(d...)

	// Collect all errors
	if resp.Diagnostics.HasError() {
		return
	}

	for _, mp := range matchedPaths {
		// If the user specifies the same attribute this validator is applied to,
		// also as part of the input, skip it
		if mp.Equal(req.Path) {
			continue
		}

		var mpVal types.Int32
		diags := req.Config.GetAttribute(ctx, mp, &mpVal)
		resp.Diagnostics.Append(diags...)

		// Collect all errors
		if diags.HasError() {
			continue
		}

		// Delay validation until all involved attribute have a known value
		if mpVal.IsUnknown() {
			return
		}

		if mpVal.ValueInt32() >= req.ConfigValue.ValueInt32() {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Provided value is not valid",
				fmt.Sprintf("The value provided %d is not valid.  Ensure that the provided value is greater than the value configured at path %v.", req.ConfigValue.ValueInt32(), mp),
			)
			continue
		}
	}

	if resp.Diagnostics.HasError() {
		return
	}
}

// IsGreaterThanPathValue checks if an int32 is greater than the value configured in the provided path expression.
func IsGreaterThanPathValue(pathExpression path.Expression) validator.Int32 {
	return &int32IsGreaterThanPathValue{
		PathExpression: pathExpression,
	}
}
