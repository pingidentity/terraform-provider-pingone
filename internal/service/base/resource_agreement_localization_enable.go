// Copyright © 2025 Ping Identity Corporation

package base

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type AgreementLocalizationEnableResource serviceClientType

type AgreementLocalizationEnableResourceModel struct {
	Id                      pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId           pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	AgreementId             pingonetypes.ResourceIDValue `tfsdk:"agreement_id"`
	AgreementLocalizationId pingonetypes.ResourceIDValue `tfsdk:"agreement_localization_id"`
	Enabled                 types.Bool                   `tfsdk:"enabled"`
}

// Framework interfaces
var (
	_ resource.Resource                = &AgreementLocalizationEnableResource{}
	_ resource.ResourceWithConfigure   = &AgreementLocalizationEnableResource{}
	_ resource.ResourceWithImportState = &AgreementLocalizationEnableResource{}
)

// New Object
func NewAgreementLocalizationEnableResource() resource.Resource {
	return &AgreementLocalizationEnableResource{}
}

// Metadata
func (r *AgreementLocalizationEnableResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_agreement_localization_enable"
}

// Schema.
func (r *AgreementLocalizationEnableResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	const attrMinLength = 1

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage the enabled status of an agreement localization in a PingOne environment.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment configured with an agreement localization to enable/disable."),
			),

			"agreement_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the agreement configured with an agreement localization to enable/disable."),
			),

			"agreement_localization_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the agreement localization to enable/disable."),
			),

			"enabled": schema.BoolAttribute{
				Description: "A boolean that specifies the current enabled state of the agreement localization. The agreement localization must have an active revision text to be enabled.",
				Required:    true,
			},
		},
	}
}

func (r *AgreementLocalizationEnableResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	resourceConfig, ok := req.ProviderData.(framework.ResourceType)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected the provider client, got: %T. Please report this issue to the provider maintainers.", req.ProviderData),
		)

		return
	}

	r.Client = resourceConfig.Client.API
	if r.Client == nil {
		resp.Diagnostics.AddError(
			"Client not initialised",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.",
		)
		return
	}
}

func (r *AgreementLocalizationEnableResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state AgreementLocalizationEnableResourceModel

	if r.Client == nil || r.Client.ManagementAPIClient == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var agreementResponse *management.AgreementLanguage
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.AgreementLanguagesResourcesApi.ReadOneAgreementLanguage(ctx, plan.EnvironmentId.ValueString(), plan.AgreementId.ValueString(), plan.AgreementLocalizationId.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadOneAgreementLanguage",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		&agreementResponse,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the model for the API
	agreementLocalizationEnable := plan.expand(agreementResponse)

	// Run the API call
	var response *management.AgreementLanguage
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.AgreementLanguagesResourcesApi.UpdateAgreementLanguage(ctx, plan.EnvironmentId.ValueString(), plan.AgreementId.ValueString(), plan.AgreementLocalizationId.ValueString()).AgreementLanguage(*agreementLocalizationEnable).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateAgreementLanguage",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		&response,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(response)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *AgreementLocalizationEnableResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *AgreementLocalizationEnableResourceModel

	if r.Client == nil || r.Client.ManagementAPIClient == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *management.AgreementLanguage
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.AgreementLanguagesResourcesApi.ReadOneAgreementLanguage(ctx, data.EnvironmentId.ValueString(), data.AgreementId.ValueString(), data.AgreementLocalizationId.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadOneAgreementLanguage",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
		&response,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Remove from state if resource is not found
	if response == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(response)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AgreementLocalizationEnableResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state AgreementLocalizationEnableResourceModel

	if r.Client == nil || r.Client.ManagementAPIClient == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var agreementLanguageResponse *management.AgreementLanguage
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.AgreementLanguagesResourcesApi.ReadOneAgreementLanguage(ctx, plan.EnvironmentId.ValueString(), plan.AgreementId.ValueString(), plan.AgreementLocalizationId.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadOneAgreementLanguage",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		&agreementLanguageResponse,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the model for the API
	agreementLanguageEnable := plan.expand(agreementLanguageResponse)

	// Run the API call
	var response *management.AgreementLanguage
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.AgreementLanguagesResourcesApi.UpdateAgreementLanguage(ctx, plan.EnvironmentId.ValueString(), plan.AgreementId.ValueString(), plan.AgreementLocalizationId.ValueString()).AgreementLanguage(*agreementLanguageEnable).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateAgreementLanguage",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		&response,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(response)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *AgreementLocalizationEnableResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *AgreementLocalizationEnableResourceModel

	if r.Client == nil || r.Client.ManagementAPIClient == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var agreementResponse *management.AgreementLanguage
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.AgreementLanguagesResourcesApi.ReadOneAgreementLanguage(ctx, data.EnvironmentId.ValueString(), data.AgreementId.ValueString(), data.AgreementLocalizationId.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadOneAgreementLanguage",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
		&agreementResponse,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the model for the API
	data.Enabled = types.BoolValue(false)
	agreementDisable := data.expand(agreementResponse)

	// Run the API call
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.AgreementLanguagesResourcesApi.UpdateAgreementLanguage(ctx, data.EnvironmentId.ValueString(), data.AgreementId.ValueString(), data.AgreementLocalizationId.ValueString()).AgreementLanguage(*agreementDisable).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateAgreementLanguage",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
		nil,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *AgreementLocalizationEnableResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:  "agreement_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:     "agreement_localization_id",
			Regexp:    verify.P1ResourceIDRegexp,
			PrimaryID: true,
		},
	}

	attributes, err := framework.ParseImportID(req.ID, idComponents...)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			err.Error(),
		)
		return
	}

	for _, idComponent := range idComponents {
		pathKey := idComponent.Label

		if idComponent.PrimaryID {
			pathKey = "id"
		}

		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(pathKey), attributes[idComponent.Label])...)
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("agreement_localization_id"), attributes["agreement_localization_id"])...)
}

func (p *AgreementLocalizationEnableResourceModel) expand(existingObject *management.AgreementLanguage) *management.AgreementLanguage {

	data := management.NewAgreementLanguage(existingObject.GetDisplayName(), p.Enabled.ValueBool(), existingObject.GetLocale())

	userExperience := management.NewAgreementLanguageUserExperience()

	if v, ok := existingObject.GetUserExperienceOk(); ok {

		if c, ok := v.GetAcceptCheckboxTextOk(); ok {
			userExperience.SetAcceptCheckboxText(*c)
		}

		if c, ok := v.GetContinueButtonTextOk(); ok {
			userExperience.SetContinueButtonText(*c)
		}

		if c, ok := v.GetDeclineButtonTextOk(); ok {
			userExperience.SetDeclineButtonText(*c)
		}

		data.SetUserExperience(*userExperience)
	}

	return data
}

func (p *AgreementLocalizationEnableResourceModel) toState(apiObject *management.AgreementLanguage) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.PingOneResourceIDToTF(apiObject.GetId())
	p.AgreementLocalizationId = framework.PingOneResourceIDToTF(apiObject.GetId())
	p.Enabled = framework.BoolOkToTF(apiObject.GetEnabledOk())

	return diags
}
