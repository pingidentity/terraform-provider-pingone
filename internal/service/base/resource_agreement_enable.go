package base

import (
	"context"
	"fmt"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type AgreementEnableResource serviceClientType

type AgreementEnableResourceModel struct {
	Id            pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	AgreementId   pingonetypes.ResourceIDValue `tfsdk:"agreement_id"`
	Enabled       types.Bool                   `tfsdk:"enabled"`
}

// Framework interfaces
var (
	_ resource.Resource                = &AgreementEnableResource{}
	_ resource.ResourceWithConfigure   = &AgreementEnableResource{}
	_ resource.ResourceWithImportState = &AgreementEnableResource{}
)

// New Object
func NewAgreementEnableResource() resource.Resource {
	return &AgreementEnableResource{}
}

// Metadata
func (r *AgreementEnableResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_agreement_enable"
}

// Schema.
func (r *AgreementEnableResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	const attrMinLength = 1

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage the enabled status of an agreement in a PingOne environment.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment configured with an agreement to enable/disable."),
			),

			"agreement_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the agreement to set the enabled status for."),
			),

			"enabled": schema.BoolAttribute{
				Description: "A boolean that specifies the current enabled state of the agreement. The agreement must support the default language to be enabled. It cannot be disabled if it is referenced by a sign-on policy action. When an agreement is disabled, it is not used anywhere that it is configured across PingOne.",
				Required:    true,
			},
		},
	}
}

func (r *AgreementEnableResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *AgreementEnableResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state AgreementEnableResourceModel

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

	var agreementResponse *management.Agreement
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.AgreementsResourcesApi.ReadOneAgreement(ctx, plan.EnvironmentId.ValueString(), plan.AgreementId.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadOneAgreement",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		&agreementResponse,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the model for the API
	agreementEnable := plan.expand(agreementResponse)

	// Run the API call
	var response *management.Agreement
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.AgreementsResourcesApi.UpdateAgreement(ctx, plan.EnvironmentId.ValueString(), plan.AgreementId.ValueString()).Agreement(*agreementEnable).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateAgreement",
		agreementEnableUpdateCustomErrorHandler,
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

func (r *AgreementEnableResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *AgreementEnableResourceModel

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
	var response *management.Agreement
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.AgreementsResourcesApi.ReadOneAgreement(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadOneAgreement",
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

func (r *AgreementEnableResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state AgreementEnableResourceModel

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

	var agreementResponse *management.Agreement
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.AgreementsResourcesApi.ReadOneAgreement(ctx, plan.EnvironmentId.ValueString(), plan.AgreementId.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadOneAgreement",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		&agreementResponse,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the model for the API
	agreementEnable := plan.expand(agreementResponse)

	// Run the API call
	var response *management.Agreement
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.AgreementsResourcesApi.UpdateAgreement(ctx, plan.EnvironmentId.ValueString(), plan.Id.ValueString()).Agreement(*agreementEnable).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateAgreement",
		agreementEnableUpdateCustomErrorHandler,
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

func (r *AgreementEnableResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *AgreementEnableResourceModel

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

	var agreementResponse *management.Agreement
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.AgreementsResourcesApi.ReadOneAgreement(ctx, data.EnvironmentId.ValueString(), data.AgreementId.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadOneAgreement",
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
			fO, fR, fErr := r.Client.ManagementAPIClient.AgreementsResourcesApi.UpdateAgreement(ctx, data.EnvironmentId.ValueString(), data.AgreementId.ValueString()).Agreement(*agreementDisable).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateAgreement",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
		nil,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *AgreementEnableResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:     "agreement_id",
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
}

func (p *AgreementEnableResourceModel) expand(existingObject *management.Agreement) *management.Agreement {

	data := management.NewAgreement(p.Enabled.ValueBool(), existingObject.GetName())

	if v, ok := existingObject.GetDescriptionOk(); ok {
		data.SetDescription(*v)
	}

	if v, ok := existingObject.GetReconsentPeriodDaysOk(); ok {
		data.SetReconsentPeriodDays(*v)
	}

	return data
}

func (p *AgreementEnableResourceModel) toState(apiObject *management.Agreement) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.PingOneResourceIDToTF(apiObject.GetId())
	p.EnvironmentId = framework.PingOneResourceIDToTF(*apiObject.GetEnvironment().Id)
	p.AgreementId = framework.PingOneResourceIDToTF(apiObject.GetId())
	p.Enabled = framework.BoolOkToTF(apiObject.GetEnabledOk())

	return diags
}

func agreementEnableUpdateCustomErrorHandler(_ *http.Response, p1Error *model.P1Error) diag.Diagnostics {
	var diags diag.Diagnostics

	if p1Error != nil {
		if v, ok := p1Error.GetDetailsOk(); ok && v != nil && len(v) > 0 {
			if v[0].GetCode() == "CONSTRAINT_VIOLATION" {
				if match, _ := regexp.MatchString("The agreement can not be enabled without supporting the default language configured for the environment.", v[0].GetMessage()); match {
					diags.AddError(
						v[0].GetMessage(),
						"The agreement must have an enabled agreement localization for the default language of the environment.  Ensure that a `pingone_agreement_localization`, `pingone_agreement_localization_revision` and `pingone_agreement_localization_enable` resource exist for the default language, or the environment's default language is re-configured using the `pingone_language_update` resource.",
					)

					return diags
				}
			}
		}
	}

	return diags
}
