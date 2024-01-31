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

// stringB64ContentTypeValidator validates that string is parseable JSON.
type stringB64ContentTypeValidator struct {
	VerifyContentTypes []string
}

// Description describes the validation in plain text formatting.
func (v stringB64ContentTypeValidator) Description(_ context.Context) string {
	return fmt.Sprintf("Ensure that the provided string is a base64 encoded content type of either %s", strings.Join(v.VerifyContentTypes, ", "))
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v stringB64ContentTypeValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// Validate performs the validation.
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

// IsB64ContentType checks that a set of path.Expression,
// including the attribute the validator is applied to,
// is a valid base64 encoded string that is one of the specified content types.
//
// Relative path.Expression will be resolved using the attribute being
// validated.
func IsB64ContentType(verifyContentTypes ...string) validator.String {
	return stringB64ContentTypeValidator{
		VerifyContentTypes: verifyContentTypes,
	}
}
