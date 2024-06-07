package sso

import (
	"context"
	"fmt"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type PopulationResource struct {
	serviceClientType
	options client.GlobalOptions
}

type PopulationResourceModel struct {
	Id               pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId    pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	Name             types.String                 `tfsdk:"name"`
	Description      types.String                 `tfsdk:"description"`
	PasswordPolicyId pingonetypes.ResourceIDValue `tfsdk:"password_policy_id"`
}

// Framework interfaces
var (
	_ resource.Resource                = &PopulationResource{}
	_ resource.ResourceWithConfigure   = &PopulationResource{}
	_ resource.ResourceWithImportState = &PopulationResource{}
)

// New Object
func NewPopulationResource() resource.Resource {
	return &PopulationResource{}
}

// Metadata
func (r *PopulationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_population"
}

// Schema.
func (r *PopulationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	const attrMinLength = 1

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage a PingOne population in an environment.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to create the population in."),
			),

			"name": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the name of the population.").Description,
				Required:    true,

				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			"description": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies a description to apply to the population.").Description,
				Optional:    true,
			},

			"password_policy_id": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the ID of a password policy to assign to the population.  Must be a valid PingOne resource ID.").Description,
				Optional:    true,

				CustomType: pingonetypes.ResourceIDType{},
			},
		},
	}
}

func (r *PopulationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	if resourceConfig.Client.GlobalOptions != nil {
		r.options = *resourceConfig.Client.GlobalOptions
	}

}

func (r *PopulationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state PopulationResourceModel

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
	response, d := PingOnePopulationCreate(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), *population)

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

func (r *PopulationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *PopulationResourceModel

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
		"ReadOnePopulation",
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

func (r *PopulationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state PopulationResourceModel

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
		"UpdatePopulation",
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

func (r *PopulationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *PopulationResourceModel

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

	hasUsersAssigned, d := r.hasUsersAssigned(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString())
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	if hasUsersAssigned {
		d := r.checkControlsAndDeletePopulationUsers(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString())
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Run the API call
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fR, fErr := r.Client.ManagementAPIClient.PopulationsApi.DeletePopulation(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, fR, fErr)
		},
		"DeletePopulation",
		framework.CustomErrorResourceNotFoundWarning,
		nil,
		nil,
	)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *PopulationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

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

func (p *PopulationResourceModel) expand() *management.Population {

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

func (p *PopulationResourceModel) toState(apiObject *management.Population) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.PingOneResourceIDOkToTF(apiObject.GetIdOk())
	p.EnvironmentId = framework.PingOneResourceIDOkToTF(apiObject.Environment.GetIdOk())
	p.Name = framework.StringOkToTF(apiObject.GetNameOk())
	p.Description = framework.StringOkToTF(apiObject.GetDescriptionOk())

	if v, ok := apiObject.GetPasswordPolicyOk(); ok {
		p.PasswordPolicyId = framework.PingOneResourceIDOkToTF(v.GetIdOk())
	} else {
		p.PasswordPolicyId = pingonetypes.NewResourceIDNull()
	}

	return diags
}

func (r *PopulationResource) hasUsersAssigned(ctx context.Context, environmentID, populationID string) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	users, d := r.readUsers(ctx, environmentID, populationID)
	diags.Append(d...)
	if diags.HasError() {
		return false, diags
	}

	if len(users) > 0 {
		return true, diags
	}

	return false, diags
}

func (r *PopulationResource) readUsers(ctx context.Context, environmentID, populationID string) ([]management.User, diag.Diagnostics) {
	var diags diag.Diagnostics

	if m, err := regexp.MatchString(verify.P1ResourceIDRegexpFullString.String(), populationID); err == nil && m {

		scimFilter := fmt.Sprintf(`population.id eq "%s"`, populationID)

		// Run the API call
		var entityArray *management.EntityArray
		diags.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				fO, fR, fErr := r.Client.ManagementAPIClient.UsersApi.ReadAllUsers(ctx, environmentID).Filter(scimFilter).Execute()
				return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, environmentID, fO, fR, fErr)
			},
			"ReadAllUsers",
			framework.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			&entityArray,
		)...)

		if diags.HasError() {
			return nil, diags
		}

		return entityArray.Embedded.GetUsers(), nil
	}

	if r.options.Population.ContainsUsersForceDelete {
		diags.AddError(
			"Data protection notice",
			fmt.Sprintf("For data protection reasons, it could not be determined whether users exist in the population %[2]s in environment %[1]s. Any users in this population will not be deleted.", environmentID, populationID),
		)
	}

	return nil, diags
}

func (r *PopulationResource) checkEnvironmentControls(ctx context.Context, environmentID, populationID string) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	if r.options.Population.ContainsUsersForceDelete {
		// Check if the environment is a sandbox type.  We'll only delete users in sandbox environments
		var environmentResponse *management.Environment
		diags.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				fO, fR, fErr := r.Client.ManagementAPIClient.EnvironmentsApi.ReadOneEnvironment(ctx, environmentID).Execute()
				return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, environmentID, fO, fR, fErr)
			},
			"ReadOneEnvironment-DeletePopulation",
			framework.DefaultCustomError,
			nil,
			&environmentResponse,
		)...)
		if diags.HasError() {
			return false, diags
		}

		if v, ok := environmentResponse.GetTypeOk(); ok && *v == management.ENUMENVIRONMENTTYPE_SANDBOX {
			return true, diags
		} else {
			diags.AddWarning(
				"Data protection notice",
				fmt.Sprintf("For data protection reasons, the provider configuration `global_options.population.contains_users_force_delete` has no effect on environment ID %[1]s as it has a type set to `PRODUCTION`.  Users in this population will not be deleted.\n"+
					"If you wish to force delete population %[2]s in environment %[1]s, please review and remove user data manually.", environmentID, populationID),
			)
		}
	}

	return false, diags
}

func (r *PopulationResource) checkControlsAndDeletePopulationUsers(ctx context.Context, environmentID, populationID string) diag.Diagnostics {
	var diags diag.Diagnostics

	environmentControlsOk, d := r.checkEnvironmentControls(ctx, environmentID, populationID)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	if environmentControlsOk {

		loopCounter := 1
		for loopCounter > 0 {

			users, d := r.readUsers(ctx, environmentID, populationID)
			diags.Append(d...)
			if diags.HasError() {
				return diags
			}

			// DELETE USERS
			if len(users) == 0 {
				break
			} else {
				for _, user := range users {
					var entityArray *management.EntityArray
					diags.Append(framework.ParseResponse(
						ctx,

						func() (any, *http.Response, error) {
							fR, fErr := r.Client.ManagementAPIClient.UsersApi.DeleteUser(ctx, environmentID, user.GetId()).Execute()
							return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, environmentID, nil, fR, fErr)
						},
						"DeleteUser-DeletePopulation",
						framework.DefaultCustomError,
						nil,
						&entityArray,
					)...)

					if diags.HasError() {
						return diags
					}
				}
			}
		}
	}

	return diags
}

func PingOnePopulationCreate(ctx context.Context, apiClient *management.APIClient, environmentID string, population management.Population) (*management.Population, diag.Diagnostics) {
	var diags diag.Diagnostics

	var returnVar *management.Population
	diags.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := apiClient.PopulationsApi.CreatePopulation(ctx, environmentID).Population(population).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, apiClient, environmentID, fO, fR, fErr)
		},
		"CreatePopulation",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		&returnVar,
	)...)

	if diags.HasError() {
		return nil, diags
	}

	return returnVar, diags
}
