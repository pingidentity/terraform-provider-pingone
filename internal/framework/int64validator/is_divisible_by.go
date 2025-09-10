// Copyright © 2025 Ping Identity Corporation

// Package int64validator provides custom int64 validators for the Terraform Plugin Framework.
// This package contains validators that check integer attribute constraints, ranges,
// and mathematical relationships for the PingOne provider's specific requirements.
package int64validator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// int64IsDivisibleBy validates that an int64 value is exactly divisible by a specified denominator.
// It checks that the value modulo the denominator equals zero, ensuring clean mathematical division.
type int64IsDivisibleBy struct {
	// Denominator is the value that the int64 must be exactly divisible by
	Denominator int64
}

// Description describes the validation in plain text formatting.
// It returns a human-readable description of the validation rule for error messages and documentation.
func (v int64IsDivisibleBy) Description(ctx context.Context) string {
	return fmt.Sprintf("Ensure the integer is exactly divisible by %d.", v.Denominator)
}

// MarkdownDescription describes the validation in Markdown formatting.
// It returns the same description as Description() but formatted for Markdown documentation.
func (v int64IsDivisibleBy) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateInt64 performs the validation logic for divisibility checking.
// It checks that the int64 value is exactly divisible by the denominator with no remainder.
// Null and unknown values are considered valid and skip validation.
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

// IsDivisibleBy creates a validator that checks if an int64 value is exactly divisible by the specified denominator.
// It returns a validator that ensures the integer value has no remainder when divided by the denominator.
// This is useful for enforcing constraints like multiples of specific values or alignment requirements.
//
// The denominator parameter specifies the value that the int64 must be exactly divisible by.
// Values that do not divide evenly will result in a validation error.
func IsDivisibleBy(denominator int64) validator.Int64 {
	return &int64IsDivisibleBy{
		Denominator: denominator,
	}
}
