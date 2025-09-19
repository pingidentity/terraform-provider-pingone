// Copyright © 2025 Ping Identity Corporation

package stringvalidator

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ validator.String = shouldNotContainValidator{}

// shouldNotContainValidator validates that a string does not contain any of the specified substrings.
// It checks the string content to ensure none of the prohibited values appear anywhere within it.
type shouldNotContainValidator struct {
	// values contains the list of substrings that are prohibited from appearing in the validated string
	values []types.String
}

// Description describes the validation in plain text formatting.
// It returns a human-readable description of the validation rule for error messages and documentation.
func (v shouldNotContainValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

// MarkdownDescription describes the validation in Markdown formatting.
// It returns a formatted description listing the prohibited values for documentation purposes.
func (v shouldNotContainValidator) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("value must not contain: %s", v.values)
}

// ValidateString performs the validation logic for substring prohibition checking.
// It examines the string value to ensure it does not contain any of the specified prohibited substrings.
// Null and unknown values are considered valid and skip validation.
func (v shouldNotContainValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue

	for _, otherValue := range v.values {
		if !strings.Contains(value.ValueString(), otherValue.ValueString()) {
			continue
		}

		response.Diagnostics.Append(validatordiag.InvalidAttributeValueMatchDiagnostic(
			request.Path,
			v.Description(ctx),
			value.String(),
		))

		break
	}
}

// ShouldNotContain creates a validator that ensures a string does not contain any of the specified substrings.
// It returns a validator that checks the string content and rejects values containing
// any of the prohibited substrings using case-sensitive matching.
//
// The values parameter contains the list of substrings that are prohibited from appearing
// anywhere within the validated string. The validation fails if any of these values
// are found as substrings within the attribute value.
func ShouldNotContain(values ...string) validator.String {
	frameworkValues := make([]types.String, 0, len(values))

	for _, value := range values {
		frameworkValues = append(frameworkValues, types.StringValue(value))
	}

	return shouldNotContainValidator{
		values: frameworkValues,
	}
}
