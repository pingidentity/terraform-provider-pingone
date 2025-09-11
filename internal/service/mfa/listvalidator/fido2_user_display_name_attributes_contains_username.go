// Copyright © 2025 Ping Identity Corporation

// Package listvalidator provides custom list validators for MFA service configurations.
// This package contains validators that ensure proper configuration of FIDO2 policies and related MFA settings.
package listvalidator

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// listFIDO2UserDisplayNameAttributeContainsUsernameValidator validates that a list of user display name attributes contains the `username` attribute.
// This validator implements the validator.List interface to ensure FIDO2 policies include the required username attribute.
var _ validator.List = listFIDO2UserDisplayNameAttributeContainsUsernameValidator{}

// listFIDO2UserDisplayNameAttributeContainsUsernameValidator validates that a list of user display name attributes contains the `username` attribute.
// This validator ensures that FIDO2 user display name configuration includes the mandatory username attribute.
type listFIDO2UserDisplayNameAttributeContainsUsernameValidator struct{}

// FIDO2PolicyUserDisplayNameAttributesAttributesResourceModel represents the structure for FIDO2 policy user display name attributes in Terraform configurations.
// This model defines the name and sub-attributes for display name configuration in FIDO2 policies.
type FIDO2PolicyUserDisplayNameAttributesAttributesResourceModel struct {
	// Name is the name of the user display name attribute
	Name types.String `tfsdk:"name"`
	// SubAttributes contains any sub-attributes associated with the display name attribute
	SubAttributes types.List `tfsdk:"sub_attributes"`
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

// FIDO2UserDisplayNameAttributeContainsUsername returns a validator that ensures the provided user display name attributes list contains the `username` attribute.
// This function creates and returns a new instance of the FIDO2 user display name validator.
// The validator is used to enforce that FIDO2 policies include the mandatory username attribute in their display name configuration.
func FIDO2UserDisplayNameAttributeContainsUsername() validator.List {
	return &listFIDO2UserDisplayNameAttributeContainsUsernameValidator{}
}
