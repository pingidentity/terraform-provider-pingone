package listvalidator

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ validator.List = listFIDO2UserDisplayNameAttributeContainsUsernameValidator{}

// listFIDO2UserDisplayNameAttributeContainsUsernameValidator validates that a list of user display name attributes contains the `username` attribute.
type listFIDO2UserDisplayNameAttributeContainsUsernameValidator struct{}

type FIDO2PolicyUserDisplayNameAttributesAttributesResourceModel struct {
	Name          types.String `tfsdk:"name"`
	SubAttributes types.List   `tfsdk:"sub_attributes"`
}

// Description describes the validation in plain text formatting.
func (v listFIDO2UserDisplayNameAttributeContainsUsernameValidator) Description(_ context.Context) string {
	return "Ensure that the provided user display name attributes list contains the `username` attribute"
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v listFIDO2UserDisplayNameAttributeContainsUsernameValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// Validate performs the validation.
func (v listFIDO2UserDisplayNameAttributeContainsUsernameValidator) ValidateList(ctx context.Context, req validator.ListRequest, resp *validator.ListResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	var userDisplayNameAttributesAttributesPlan []FIDO2PolicyUserDisplayNameAttributesAttributesResourceModel
	resp.Diagnostics.Append(req.ConfigValue.ElementsAs(ctx, &userDisplayNameAttributesAttributesPlan, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	found := false
	for _, attributePlan := range userDisplayNameAttributesAttributesPlan {
		if attributePlan.Name.Equal(types.StringValue("username")) {
			found = true
			break
		}
	}

	if !found {
		resp.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			req.Path,
			v.Description(ctx),
			"none",
		))
	}

}

func FIDO2UserDisplayNameAttributeContainsUsername() validator.List {
	return &listFIDO2UserDisplayNameAttributeContainsUsernameValidator{}
}
