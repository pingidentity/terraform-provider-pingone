// Copyright © 2025 Ping Identity Corporation

package stringvalidator

import (
	"context"
	"encoding/base64"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// stringIsBase64EncodedValidator validates that a string contains valid base64-encoded content.
// It can handle strings with or without Content-Type prefixes and focuses on validating
// the base64 encoding of the actual content portion.
type stringIsBase64EncodedValidator struct{}

// Description describes the validation in plain text formatting.
// It returns a human-readable description of the validation rule for error messages and documentation.
func (v stringIsBase64EncodedValidator) Description(ctx context.Context) string {
	return "Ensure the string contains a base64 encoded value."
}

// MarkdownDescription describes the validation in Markdown formatting.
// It returns the same description as Description() but formatted for Markdown documentation.
func (v stringIsBase64EncodedValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateString performs the validation logic for base64 encoding verification.
// It checks that the string contains valid base64-encoded content, optionally handling
// Content-Type prefixes by extracting and validating only the content portion.
// Null and unknown values are considered valid and skip validation.
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

// IsBase64Encoded creates a validator that checks if a string contains valid base64-encoded content.
// It returns a validator that can handle strings with or without Content-Type prefixes,
// automatically extracting and validating the base64 portion of the content.
//
// If the string contains a Content-Type prefix (e.g., "data:text/plain;base64,SGVsbG8="),
// the prefix is ignored and only the subsequent base64 content is validated.
// If no prefix is found, the entire string is treated as base64 content for validation.
func IsBase64Encoded() validator.String {
	return &stringIsBase64EncodedValidator{}
}
