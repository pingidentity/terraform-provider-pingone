// Copyright © 2025 Ping Identity Corporation

package int32validator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// int32IsDivisibleBy validates if the input string is base64encoded.
type int32IsDivisibleBy struct {
	Denominator int32
}

// Description describes the validation in plain text formatting.
func (v int32IsDivisibleBy) Description(ctx context.Context) string {
	return fmt.Sprintf("Ensure the integer is exactly divisible by %d.", v.Denominator)
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v int32IsDivisibleBy) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// Validate runs the main validation logic of the validator, reading configuration data out of `req` and updating `resp` with diagnostics.
func (v int32IsDivisibleBy) ValidateInt32(ctx context.Context, req validator.Int32Request, resp *validator.Int32Response) {
	// If the value is unknown or null, there is nothing to validate.
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	if req.ConfigValue.ValueInt32()%v.Denominator != 0 {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Provided value is not valid",
			fmt.Sprintf("The value provided %d is not valid.  Ensure that the provided value is exactly divisible by %d.", req.ConfigValue.ValueInt32(), v.Denominator),
		)
		return
	}
}

// IsDivisibleBy checks if an int32 is exactly divisible by the provided denominator.
func IsDivisibleBy(denominator int32) validator.Int32 {
	return &int32IsDivisibleBy{
		Denominator: denominator,
	}
}
