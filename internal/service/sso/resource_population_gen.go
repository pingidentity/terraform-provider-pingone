// Copyright © 2025 Ping Identity Corporation
// Code generated by ping-terraform-plugin-framework-generator

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
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

var (
	_ resource.Resource                = &populationResource{}
	_ resource.ResourceWithConfigure   = &populationResource{}
	_ resource.ResourceWithImportState = &populationResource{}
)

func NewPopulationResource() resource.Resource {
	return &populationResource{}
}

type populationResource struct {
	serviceClientType
	options client.GlobalOptions
}

func (r *populationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_population"
}

func (r *populationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

type populationResourceModel struct {
	Description      types.String                 `tfsdk:"description"`
	EnvironmentId    pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	Id               pingonetypes.ResourceIDValue `tfsdk:"id"`
	Name             types.String                 `tfsdk:"name"`
	PasswordPolicyId pingonetypes.ResourceIDValue `tfsdk:"password_policy_id"`
}

func (r *populationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Resource to create and manage a PingOne population in an environment.",
		Attributes: map[string]schema.Attribute{
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "A string that specifies the description of the population.",
			},
			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to create and manage the population in."),
			),
			"id": framework.Attr_ID(),
			"name": schema.StringAttribute{
				Required:    true,
				Description: "A string that specifies the population name, which must be provided and must be unique within an environment.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"password_policy_id": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the ID of a password policy to assign to the population.  Must be a valid PingOne resource ID.").Description,
				Optional:    true,
				CustomType:  pingonetypes.ResourceIDType{},
			},
		},
	}
}

func (model *populationResourceModel) buildClientStruct() (*management.Population, diag.Diagnostics) {
	result := &management.Population{}
	// description
	result.Description = model.Description.ValueStringPointer()
	// name
	result.Name = model.Name.ValueString()
	// password_policy
	if !model.PasswordPolicyId.IsNull() {
		passwordPolicyValue := &management.PopulationPasswordPolicy{
			Id: model.PasswordPolicyId.ValueString(),
		}
		result.PasswordPolicy = passwordPolicyValue
	}

	return result, nil
}

func (state *populationResourceModel) readClientResponse(response *management.Population) diag.Diagnostics {
	// description
	state.Description = types.StringPointerValue(response.Description)
	// id
	idValue := framework.PingOneResourceIDToTF(response.GetId())
	state.Id = idValue
	// name
	state.Name = types.StringValue(response.Name)
	// password_policy_id
	var passwordPolicyIdValue pingonetypes.ResourceIDValue
	if response.PasswordPolicy == nil {
		passwordPolicyIdValue = pingonetypes.NewResourceIDNull()
	} else {
		passwordPolicyIdValue = framework.PingOneResourceIDToTF(response.PasswordPolicy.Id)
	}
	state.PasswordPolicyId = passwordPolicyIdValue
	return nil
}

func (r *populationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data populationResourceModel

	if r.Client == nil || r.Client.ManagementAPIClient == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create API call logic
	clientData, diags := data.buildClientStruct()
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var responseData *management.Population
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.PopulationsApi.CreatePopulation(ctx, data.EnvironmentId.ValueString()).Population(*clientData).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"CreatePopulation",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		&responseData,
	)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *populationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data populationResourceModel

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

	// Read API call logic
	var responseData *management.Population
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.PopulationsApi.ReadOnePopulation(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadOnePopulation",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
		&responseData,
	)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Remove from state if resource is not found
	if responseData == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *populationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data populationResourceModel

	if r.Client == nil || r.Client.ManagementAPIClient == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update API call logic
	clientData, diags := data.buildClientStruct()
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var responseData *management.Population
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.PopulationsApi.UpdatePopulation(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Population(*clientData).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdatePopulation",
		framework.DefaultCustomError,
		nil,
		&responseData,
	)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *populationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data populationResourceModel

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

	// Delete API call logic
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fR, fErr := r.Client.ManagementAPIClient.PopulationsApi.DeletePopulation(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, fR, fErr)
		},
		"DeletePopulation",
		populationDeleteCustomErrorHandler,
		nil,
		nil,
	)...)
}

func (r *populationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

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

func populationDeleteCustomErrorHandler(r *http.Response, p1Error *model.P1Error) diag.Diagnostics {
	var diags diag.Diagnostics

	if p1Error != nil {
		// Env must contain at least one population
		if details, ok := p1Error.GetDetailsOk(); ok && details != nil && len(details) > 0 {
			if code, ok := details[0].GetCodeOk(); ok && *code == "CONSTRAINT_VIOLATION" {
				if message, ok := details[0].GetMessageOk(); ok {
					m, err := regexp.MatchString(`must contain at least one population`, *message)
					if err == nil && m {
						diags.AddWarning(
							"Constraint violation",
							fmt.Sprintf("A constraint violation error was encountered: %s\n\nThe population has been removed from Terraform state, but has been left in place in the environment.", p1Error.GetMessage()),
						)

						return diags
					}
				}
			}
		}
	}

	return diags
}

func (r *populationResource) hasUsersAssigned(ctx context.Context, environmentID, populationID string) (bool, diag.Diagnostics) {
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

func (r *populationResource) readUsers(ctx context.Context, environmentID, populationID string) ([]management.User, diag.Diagnostics) {
	var diags diag.Diagnostics

	m, err := regexp.MatchString(verify.P1ResourceIDRegexpFullString.String(), populationID)
	if err != nil {
		diags.AddError(
			"Population ID validation",
			fmt.Sprintf("An error occurred while validating the population ID: %s", err.Error()),
		)
		return nil, diags
	}

	if m {

		scimFilter := fmt.Sprintf(`population.id eq "%s"`, populationID)

		// Run the API call
		var users []management.User
		diags.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				pagedIterator := r.Client.ManagementAPIClient.UsersApi.ReadAllUsers(ctx, environmentID).Filter(scimFilter).Execute()

				var initialHttpResponse *http.Response

				foundUsers := make([]management.User, 0)

				for pageCursor, err := range pagedIterator {
					if err != nil {
						return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, environmentID, nil, pageCursor.HTTPResponse, err)
					}

					if initialHttpResponse == nil {
						initialHttpResponse = pageCursor.HTTPResponse
					}

					if pageCursor.EntityArray.Embedded != nil && pageCursor.EntityArray.Embedded.Users != nil {
						foundUsers = append(foundUsers, pageCursor.EntityArray.Embedded.GetUsers()...)
					}
				}

				return foundUsers, initialHttpResponse, nil
			},
			"ReadAllUsers",
			framework.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			&users,
		)...)

		if diags.HasError() {
			return nil, diags
		}

		return users, nil
	}

	if r.options.Population.ContainsUsersForceDelete {
		diags.AddError(
			"Data protection notice",
			fmt.Sprintf("For data protection reasons, it could not be determined whether users exist in the population %[2]s in environment %[1]s. Any users in this population will not be deleted.", environmentID, populationID),
		)
	}

	return nil, diags
}

func (r *populationResource) checkEnvironmentControls(ctx context.Context, environmentID, populationID string) (bool, diag.Diagnostics) {
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

func (r *populationResource) checkControlsAndDeletePopulationUsers(ctx context.Context, environmentID, populationID string) diag.Diagnostics {
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
