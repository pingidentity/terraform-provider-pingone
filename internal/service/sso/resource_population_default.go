package sso

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type PopulationDefaultResource serviceClientType

type PopulationDefaultResourceModel struct {
	Id               types.String `tfsdk:"id"`
	EnvironmentId    types.String `tfsdk:"environment_id"`
	Name             types.String `tfsdk:"name"`
	Description      types.String `tfsdk:"description"`
	PasswordPolicyId types.String `tfsdk:"password_policy_id"`
}

// Framework interfaces
var (
	_ resource.Resource                = &PopulationDefaultResource{}
	_ resource.ResourceWithConfigure   = &PopulationDefaultResource{}
	_ resource.ResourceWithImportState = &PopulationDefaultResource{}
)

// New Object
func NewPopulationDefaultResource() resource.Resource {
	return &PopulationDefaultResource{}
}

// Metadata
func (r *PopulationDefaultResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_population_default"
}

// Schema.
func (r *PopulationDefaultResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	const attrMinLength = 1

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage the default PingOne population.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to manage the default population in."),
			),

			"name": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The name of the default population.").Description,
				Required:    true,

				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			"description": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A description to apply to the default population.").Description,
				Optional:    true,
			},

			"password_policy_id": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The ID of a password policy to assign to the default population.").Description,
				Optional:    true,

				Validators: []validator.String{
					verify.P1ResourceIDValidator(),
				},
			},
		},
	}
}

func (r *PopulationDefaultResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *PopulationDefaultResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state PopulationDefaultResourceModel

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
	population := plan.expand()

	// Run the API call
	response, d := PingOnePopulationDefaultCreate(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), *population)

	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(response)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *PopulationDefaultResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *PopulationDefaultResourceModel

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
	var response *management.Population
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.PopulationsApi.ReadOnePopulation(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadOnePopulation-Default",
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

func (r *PopulationDefaultResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state PopulationDefaultResourceModel

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
	population := plan.expand()

	// Run the API call
	var response *management.Population
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.PopulationsApi.UpdatePopulation(ctx, plan.EnvironmentId.ValueString(), plan.Id.ValueString()).Population(*population).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdatePopulation-Default",
		framework.DefaultCustomError,
		nil,
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

func (r *PopulationDefaultResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *PopulationDefaultResourceModel

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
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fR, fErr := r.Client.ManagementAPIClient.PopulationsApi.DeletePopulation(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, fR, fErr)
		},
		"DeletePopulation-Default",
		framework.CustomErrorResourceNotFoundWarning,
		nil,
		nil,
	)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *PopulationDefaultResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:     "population_id",
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

func (p *PopulationDefaultResourceModel) expand() *management.Population {

	data := management.NewPopulation(p.Name.ValueString())

	if !p.Description.IsNull() && !p.Description.IsUnknown() {
		data.SetDescription(p.Description.ValueString())
	}

	if !p.PasswordPolicyId.IsNull() && !p.PasswordPolicyId.IsUnknown() {
		data.SetPasswordPolicy(
			*management.NewPopulationPasswordPolicy(p.PasswordPolicyId.ValueString()),
		)
	}

	return data
}

func (p *PopulationDefaultResourceModel) toState(apiObject *management.Population) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.StringOkToTF(apiObject.GetIdOk())
	p.EnvironmentId = framework.StringOkToTF(apiObject.Environment.GetIdOk())
	p.Name = framework.StringOkToTF(apiObject.GetNameOk())
	p.Description = framework.StringOkToTF(apiObject.GetDescriptionOk())

	if v, ok := apiObject.GetPasswordPolicyOk(); ok {
		p.PasswordPolicyId = framework.StringOkToTF(v.GetIdOk())
	} else {
		p.PasswordPolicyId = types.StringNull()
	}

	return diags
}

func PingOnePopulationDefaultCreate(ctx context.Context, apiClient *management.APIClient, environmentID string, population management.Population) (*management.Population, diag.Diagnostics) {
	var diags diag.Diagnostics

	var returnVar *management.Population
	diags.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := apiClient.PopulationsApi.CreatePopulation(ctx, environmentID).Population(population).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, apiClient, environmentID, fO, fR, fErr)
		},
		"CreatePopulation-Default",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		&returnVar,
	)...)

	if diags.HasError() {
		return nil, diags
	}

	return returnVar, diags
}
