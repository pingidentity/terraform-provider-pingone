// Copyright © 2025 Ping Identity Corporation

// Package stringvalidator provides custom string validators for the Terraform Plugin Framework.
// This package contains validators that check string attribute constraints, formatting,
// and content validation for the PingOne provider's specific requirements.
package stringvalidator

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.String = stringB64ContentTypeValidator{}

// stringB64ContentTypeValidator validates that a string contains base64-encoded content of specific types.
// It decodes the base64 string and checks that the content matches one of the expected MIME types
// using HTTP content type detection.
type stringB64ContentTypeValidator struct {
	// VerifyContentTypes is the list of acceptable MIME content types for the decoded content
	VerifyContentTypes []string
}

// Description describes the validation in plain text formatting.
// It returns a human-readable description of the validation rule for error messages and documentation.
func (v stringB64ContentTypeValidator) Description(_ context.Context) string {
	return fmt.Sprintf("Ensure that the provided string is a base64 encoded content type of either %s", strings.Join(v.VerifyContentTypes, ", "))
}

// MarkdownDescription describes the validation in Markdown formatting.
// It returns the same description as Description() but formatted for Markdown documentation.
func (v stringB64ContentTypeValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateString performs the validation logic for base64 content type checking.
// It decodes the base64 string value and uses HTTP content type detection to verify
// that the decoded content matches one of the expected MIME types.
// Null and unknown values are considered valid and skip validation.
func (v stringB64ContentTypeValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	archive, err := base64.StdEncoding.DecodeString(req.ConfigValue.ValueString())
	if err != nil {
		resp.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			req.Path,
			v.Description(ctx),
			req.ConfigValue.ValueString(),
		))

		return
	}

	contentType := http.DetectContentType(archive)

	found := false

	for _, verifyContentType := range v.VerifyContentTypes {
		if strings.EqualFold(verifyContentType, contentType) {
			found = true
			break
		}
	}

	if !found {
		resp.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			req.Path,
			v.Description(ctx),
			req.ConfigValue.ValueString(),
		))

		return
	}

}

// IsB64ContentType creates a validator that checks base64-encoded content against expected MIME types.
// It returns a validator that decodes the base64 string and verifies the content type
// matches one of the specified acceptable content types.
//
// The verifyContentTypes parameter contains the list of acceptable MIME content types
// that the decoded base64 content must match (case-insensitive comparison).
//
// This validator is useful for attributes that accept base64-encoded files where the
// file type must be restricted (e.g., images, certificates, or specific document types).
func IsB64ContentType(verifyContentTypes ...string) validator.String {
	return stringB64ContentTypeValidator{
		VerifyContentTypes: verifyContentTypes,
	}
}
