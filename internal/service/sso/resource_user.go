package sso

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type UserResource struct {
	client *management.APIClient
	region model.RegionMapping
}

type UserResourceModel struct {
	Id            types.String `tfsdk:"id"`
	EnvironmentId types.String `tfsdk:"environment_id"`
	Username      types.String `tfsdk:"username"`
	Email         types.String `tfsdk:"email"`
	Status        types.String `tfsdk:"status"`
	PopulationId  types.String `tfsdk:"population_id"`
}

// Framework interfaces
var (
	_ resource.Resource                = &UserResource{}
	_ resource.ResourceWithConfigure   = &UserResource{}
	_ resource.ResourceWithImportState = &UserResource{}
)

// New Object
func NewUserResource() resource.Resource {
	return &UserResource{}
}

// Metadata
func (r *UserResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

// Schema.
func (r *UserResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	const attrMinLength = 1

	statusDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The enabled status of the user.",
	).AllowedValues([]string{"ENABLED", "DISABLED"}).DefaultValue("ENABLED")

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage PingOne users.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to create the user in."),
			),

			"username": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The username of the user.").Description,
				Required:    true,

				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			"email": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The email address of the user.").Description,
				Required:    true,

				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			"status": schema.StringAttribute{
				Description:         statusDescription.Description,
				MarkdownDescription: statusDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,

				Default: stringdefault.StaticString("ENABLED"),

				Validators: []validator.String{
					stringvalidator.OneOf("ENABLED", "DISABLED"),
				},
			},

			"population_id": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The population ID to add the user to.").Description,
				Required:    true,

				Validators: []validator.String{
					verify.P1ResourceIDValidator(),
				},
			},
		},
	}
}

func (r *UserResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	preparedClient, err := prepareClient(ctx, resourceConfig)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			err.Error(),
		)

		return
	}

	r.client = preparedClient
	r.region = resourceConfig.Client.API.Region
}

func (r *UserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state UserResourceModel

	if r.client == nil {
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
	user, userEnabled := plan.expand()

	// Run the API call
	response, d := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.UsersApi.CreateUser(ctx, plan.EnvironmentId.ValueString()).User(*user).Execute()
		},
		"CreateUser",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)

	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	userResponse := response.(*management.User)

	responseEnabled, d := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.EnableUsersApi.UpdateUserEnabled(ctx, plan.EnvironmentId.ValueString(), userResponse.GetId()).UserEnabled(*userEnabled).Execute()
		},
		"UpdateUserEnabled",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)

	resp.Diagnostics.Append(d...)

	// Create the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(userResponse, responseEnabled.(*management.UserEnabled))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *UserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *UserResourceModel

	if r.client == nil {
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
	response, d := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.UsersApi.ReadUser(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
		},
		"ReadUser",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
	)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Remove from state if resource is not found
	if response == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	responseEnabled, d := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.EnableUsersApi.ReadUserEnabled(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
		},
		"ReadUserEnabled",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
	)
	resp.Diagnostics.Append(d...)

	// Remove from state if resource is not found
	if responseEnabled == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(response.(*management.User), responseEnabled.(*management.UserEnabled))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state UserResourceModel

	if r.client == nil {
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
	user, userEnabled := plan.expand()

	// Run the API call
	response, d := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.UsersApi.UpdateUserPut(ctx, plan.EnvironmentId.ValueString(), plan.Id.ValueString()).User(*user).Execute()
		},
		"UpdateUserPut",
		framework.DefaultCustomError,
		nil,
	)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	responseEnabled, d := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.EnableUsersApi.UpdateUserEnabled(ctx, plan.EnvironmentId.ValueString(), plan.Id.ValueString()).UserEnabled(*userEnabled).Execute()
		},
		"UpdateUserEnabled",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)

	resp.Diagnostics.Append(d...)

	// Create the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(response.(*management.User), responseEnabled.(*management.UserEnabled))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *UserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *UserResourceModel

	if r.client == nil {
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
	_, d := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			r, err := r.client.UsersApi.DeleteUser(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return nil, r, err
		},
		"DeleteUser",
		framework.CustomErrorResourceNotFoundWarning,
		nil,
	)
	resp.Diagnostics.Append(d...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *UserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splitLength := 2
	attributes := strings.SplitN(req.ID, "/", splitLength)

	if len(attributes) != splitLength {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("invalid id (\"%s\") specified, should be in format \"environment_id/user_id\"", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_id"), attributes[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), attributes[1])...)
}

func (p *UserResourceModel) expand() (*management.User, *management.UserEnabled) {

	userData := management.NewUser(p.Email.ValueString(), p.Username.ValueString())

	population := *management.NewUserPopulation(p.PopulationId.ValueString())
	userData.SetPopulation(population)

	userEnabledData := management.NewUserEnabled()
	if p.Status.ValueString() == "ENABLED" {
		userEnabledData.SetEnabled(true)
	} else {
		userEnabledData.SetEnabled(false)
	}

	return userData, userEnabledData
}

func (p *UserResourceModel) toState(apiObject *management.User, apiObjectEnabled *management.UserEnabled) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil || apiObjectEnabled == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.StringOkToTF(apiObject.GetIdOk())
	p.EnvironmentId = framework.StringOkToTF(apiObject.Environment.GetIdOk())
	p.Username = framework.StringOkToTF(apiObject.GetUsernameOk())
	p.Email = framework.StringOkToTF(apiObject.GetEmailOk())

	if v, ok := apiObjectEnabled.GetEnabledOk(); ok && *v {
		p.Status = framework.StringToTF("ENABLED")
	} else {
		p.Status = framework.StringToTF("DISABLED")
	}

	p.PopulationId = framework.StringOkToTF(apiObject.Population.GetIdOk())

	return diags
}
