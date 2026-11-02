// Copyright © 2026 Ping Identity Corporation

package stringvalidator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

var _ validator.String = isURLValidator{}

// isURLValidator validates that a string is a URL with one of the allowed schemes.
type isURLValidator struct {
	Schemes []string
}

// Description describes the validation in plain text formatting.
func (v isURLValidator) Description(_ context.Context) string {
	return fmt.Sprintf("Value must be a valid URL with %s prefix.", verify.QuotedSchemes(v.Schemes))
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v isURLValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateString performs the validation.
func (v isURLValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	if err := verify.IsURLWithSchemes(req.ConfigValue.ValueString(), v.Schemes...); err != nil {
		resp.Diagnostics.Append(validatordiag.InvalidAttributeValueMatchDiagnostic(
			req.Path,
			v.Description(ctx),
			req.ConfigValue.ValueString(),
		))
	}
}

// IsURLWithHTTPSValidator returns a validator that checks a string is a valid URL with an `https://` scheme.
func IsURLWithHTTPSValidator() validator.String {
	return isURLValidator{Schemes: []string{"https"}}
}

// IsURLWithHTTPorHTTPSValidator returns a validator that checks a string is a valid URL with an `http://` or `https://` scheme.
func IsURLWithHTTPorHTTPSValidator() validator.String {
	return isURLValidator{Schemes: []string{"http", "https"}}
}
