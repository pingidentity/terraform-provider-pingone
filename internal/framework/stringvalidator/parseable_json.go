package stringvalidator

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.String = StringParseableJSONValidator{}

// StringParseableJSONValidator validates that string is parseable JSON.
type StringParseableJSONValidator struct{}

// Description describes the validation in plain text formatting.
func (v StringParseableJSONValidator) Description(_ context.Context) string {
	return "Ensure that the provided string is valid JSON."
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v StringParseableJSONValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// Validate performs the validation.
func (v StringParseableJSONValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()

	jsonBytes, err := json.Marshal(value)
	if err != nil {
		resp.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			req.Path,
			v.Description(ctx),
			value,
		))

		return
	}

	var condition map[string]interface{}
	err = json.Unmarshal([]byte(jsonBytes), &condition)
	if err != nil {
		resp.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			req.Path,
			v.Description(ctx),
			string(jsonBytes),
		))

		return
	}
}

// IsParseableJSON checks that a set of path.Expression,
// including the attribute the validator is applied to,
// must have a true value.
//
// Relative path.Expression will be resolved using the attribute being
// validated.
func IsParseableJSON() validator.String {
	return StringParseableJSONValidator{}
}
