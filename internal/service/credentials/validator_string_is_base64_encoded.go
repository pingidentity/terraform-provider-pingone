package credentials

import (
	"context"
	"encoding/base64"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// stringIsBase64EncodedValidator validates if the input string is base64encoded.
type stringIsBase64EncodedValidator struct{}

// Description describes the validation in plain text formatting.
func (v stringIsBase64EncodedValidator) Description(ctx context.Context) string {
	return "Ensure the string contains a base64 encoded value."
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v stringIsBase64EncodedValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// Validate runs the main validation logic of the validator, reading configuration data out of `req` and updating `resp` with diagnostics.
func (v stringIsBase64EncodedValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	// If the value is unknown or null, there is nothing to validate.
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	s := req.ConfigValue.ValueString()

	// determine if the input value has content type
	numberOfSubstrings := 2
	re := regexp.MustCompile(`^(\w+):(\w+)\/(\w+);base64,`)
	matches := re.Split(s, numberOfSubstrings)

	if len(matches) == numberOfSubstrings {
		// Content-Type was found; obtain the value that occurs after the prefix.
		s = matches[1]
	} else {
		// Content-Type was not found; use the input string input as-is.
		s = matches[0]
	}

	// Check if the string can be base64 decoded.
	_, err := base64.StdEncoding.DecodeString(s)

	if err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Unable to Base64 Decode String Value",
			fmt.Sprintf("unable to base64 decode string at: %s", req.Path),
		)

		return
	}
}

// IsBase64Encoded checks if a string is base64 encdoed.
//
// If the string contains a Content-Type prefex, the prefix is ignored
// and the subsequent substring is evaluated.
func IsBase64Encoded() validator.String {
	return &stringIsBase64EncodedValidator{}
}
