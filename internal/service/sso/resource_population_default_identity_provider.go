// Copyright © 2025 Ping Identity Corporation

package sso

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
type PopulationDefaultIdpResource serviceClientType

type PopulationDefaultIdpResourceModel struct {
	EnvironmentId      pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	PopulationId       pingonetypes.ResourceIDValue `tfsdk:"population_id"`
	IdentityProviderId pingonetypes.ResourceIDValue `tfsdk:"identity_provider_id"`
	Type               types.String                 `tfsdk:"type"`
}

// Framework interfaces
var (
	_ resource.Resource                = &PopulationDefaultIdpResource{}
	_ resource.ResourceWithConfigure   = &PopulationDefaultIdpResource{}
	_ resource.ResourceWithImportState = &PopulationDefaultIdpResource{}
)

// New Object
func NewPopulationDefaultIdpResource() resource.Resource {
	return &PopulationDefaultIdpResource{}
}

// Metadata
func (r *PopulationDefaultIdpResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_population_default_identity_provider"
}

// Schema.
func (r *PopulationDefaultIdpResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	const attrMinLength = 1

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage the default Identity Provider for a given population in a PingOne environment.",

		Attributes: map[string]schema.Attribute{
			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment that contains the population to assign a default Identity provider to."),
			),

			"population_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the population to assign the default Identity Provider to."),
			),

			"identity_provider_id": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The ID of the Identity Provider to assign as the default for the given population.  To specify PingOne as the default identity provider, leave this field undefined, or assign a null value.  When defined, must be a valid PingOne resource ID.").Description,
				Optional:    true,

				CustomType: pingonetypes.ResourceIDType{},
			},

			"type": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The type of the Identity Provider.").Description,
				Computed:    true,
			},
		},
	}
}

func (r *PopulationDefaultIdpResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *PopulationDefaultIdpResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state PopulationDefaultIdpResourceModel

	if r.Client.ManagementAPIClient == nil {
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

	// Build the model for the API
	populationDefaultIdp := plan.expand()

	// Run the API call
	var response *management.PopulationDefaultIdp
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.PopulationsApi.UpdatePopulationDefaultIdp(ctx, plan.EnvironmentId.ValueString(), plan.PopulationId.ValueString()).PopulationDefaultIdp(*populationDefaultIdp).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdatePopulationDefaultIdp",
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

func (r *PopulationDefaultIdpResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *PopulationDefaultIdpResourceModel

	if r.Client.ManagementAPIClient == nil {
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
	var response *management.PopulationDefaultIdp
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.PopulationsApi.ReadOnePopulationDefaultIdp(ctx, data.EnvironmentId.ValueString(), data.PopulationId.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadOnePopulationDefaultIdp",
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

func (r *PopulationDefaultIdpResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state PopulationDefaultIdpResourceModel

	if r.Client.ManagementAPIClient == nil {
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

	// Build the model for the API
	populationDefaultIdp := plan.expand()

	// Run the API call
	var response *management.PopulationDefaultIdp
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.PopulationsApi.UpdatePopulationDefaultIdp(ctx, plan.EnvironmentId.ValueString(), plan.PopulationId.ValueString()).PopulationDefaultIdp(*populationDefaultIdp).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdatePopulationDefaultIdp",
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

func (r *PopulationDefaultIdpResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *PopulationDefaultIdpResourceModel

	if r.Client.ManagementAPIClient == nil {
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

	populationDefaultIdp := management.NewPopulationDefaultIdp()

	// Run the API call
	var response *management.PopulationDefaultIdp
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.PopulationsApi.UpdatePopulationDefaultIdp(ctx, data.EnvironmentId.ValueString(), data.PopulationId.ValueString()).PopulationDefaultIdp(*populationDefaultIdp).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdatePopulationDefaultIdp",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		&response,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *PopulationDefaultIdpResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:  "population_id",
			Regexp: verify.P1ResourceIDRegexp,
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

func (p *PopulationDefaultIdpResourceModel) expand() *management.PopulationDefaultIdp {

	data := management.NewPopulationDefaultIdp()

	if !p.IdentityProviderId.IsNull() && !p.IdentityProviderId.IsUnknown() {
		data.SetId(p.IdentityProviderId.ValueString())
	}

	return data
}

func (p *PopulationDefaultIdpResourceModel) toState(apiObject *management.PopulationDefaultIdp) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.IdentityProviderId = framework.PingOneResourceIDToTF(apiObject.GetId())
	p.Type = framework.EnumOkToTF(apiObject.GetTypeOk())

	return diags
}
