// Copyright © 2025 Ping Identity Corporation

package int64validator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// int64IsDivisibleBy validates if the input string is base64encoded.
type int64IsDivisibleBy struct {
	Denominator int64
}

// Description describes the validation in plain text formatting.
func (v int64IsDivisibleBy) Description(ctx context.Context) string {
	return fmt.Sprintf("Ensure the integer is exactly divisible by %d.", v.Denominator)
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v int64IsDivisibleBy) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// Validate runs the main validation logic of the validator, reading configuration data out of `req` and updating `resp` with diagnostics.
func (v int64IsDivisibleBy) ValidateInt64(ctx context.Context, req validator.Int64Request, resp *validator.Int64Response) {
	// If the value is unknown or null, there is nothing to validate.
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	if req.ConfigValue.ValueInt64()%v.Denominator != 0 {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Provided value is not valid",
			fmt.Sprintf("The value provided %d is not valid.  Ensure that the provided value is exactly divisible by %d.", req.ConfigValue.ValueInt64(), v.Denominator),
		)
		return
	}
}

// IsDivisibleBy checks if an int64 is exactly divisible by the provided denominator.
func IsDivisibleBy(denominator int64) validator.Int64 {
	return &int64IsDivisibleBy{
		Denominator: denominator,
	}
}
