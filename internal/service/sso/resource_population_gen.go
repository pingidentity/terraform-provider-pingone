// Copyright © 2025 Ping Identity Corporation
// Code generated by ping-terraform-plugin-framework-generator

package sso

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/client"
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
	PasswordPolicy   types.Object                 `tfsdk:"password_policy"`
	UserCount        types.Int32                  `tfsdk:"user_count"`
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
			},
			"password_policy": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Required:    true,
						CustomType:  pingonetypes.ResourceIDType{},
						Description: "The ID of the password policy that is used for this population. If absent, the environment's default is used. Must be a valid PingOne resource ID.",
					},
				},
				Optional:    true,
				Description: "The object reference to the password policy resource. This is an optional property.",
				Validators: []validator.Object{
					objectvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("password_policy_id")),
				},
			},
			"password_policy_id": schema.StringAttribute{
				Description:        framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the ID of a password policy to assign to the population.  Must be a valid PingOne resource ID.").Description,
				DeprecationMessage: "This attribute is deprecated and will be removed in a future release. Please use the `password_policy.id` attribute instead.",
				Optional:           true,
				CustomType:         pingonetypes.ResourceIDType{},
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("password_policy")),
				},
			},
			"user_count": schema.Int32Attribute{
				Computed:    true,
				Description: "The number of users that belong to the population",
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
	if !model.PasswordPolicy.IsNull() {
		passwordPolicyValue := &management.PopulationPasswordPolicy{}
		passwordPolicyAttrs := model.PasswordPolicy.Attributes()
		passwordPolicyValue.Id = passwordPolicyAttrs["id"].(pingonetypes.ResourceIDValue).ValueString()
		result.PasswordPolicy = passwordPolicyValue
	} else if !model.PasswordPolicyId.IsNull() {
		// password_policy_id
		result.PasswordPolicy = &management.PopulationPasswordPolicy{
			Id: model.PasswordPolicyId.ValueString(),
		}
	}

	return result, nil
}

func (state *populationResourceModel) readClientResponse(response *management.Population) diag.Diagnostics {
	var respDiags, diags diag.Diagnostics
	// description
	state.Description = types.StringPointerValue(response.Description)
	// id
	idValue := framework.PingOneResourceIDToTF(response.GetId())
	state.Id = idValue
	// name
	state.Name = types.StringValue(response.Name)
	// password_policy_id
	if !state.PasswordPolicyId.IsNull() {
		var passwordPolicyIdValue pingonetypes.ResourceIDValue
		if response.PasswordPolicy == nil {
			passwordPolicyIdValue = pingonetypes.NewResourceIDNull()
		} else {
			passwordPolicyIdValue = framework.PingOneResourceIDToTF(response.PasswordPolicy.Id)
		}
		state.PasswordPolicyId = passwordPolicyIdValue
	} else {
		// password_policy
		passwordPolicyAttrTypes := map[string]attr.Type{
			"id": pingonetypes.ResourceIDType{},
		}
		var passwordPolicyValue types.Object
		if response.PasswordPolicy == nil {
			passwordPolicyValue = types.ObjectNull(passwordPolicyAttrTypes)
		} else {
			passwordPolicyValue, diags = types.ObjectValue(passwordPolicyAttrTypes, map[string]attr.Value{
				"id": framework.PingOneResourceIDToTF(response.PasswordPolicy.Id),
			})
			respDiags.Append(diags...)
		}
		state.PasswordPolicy = passwordPolicyValue
	}
	// user_count
	state.UserCount = types.Int32PointerValue(response.UserCount)
	return respDiags
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
