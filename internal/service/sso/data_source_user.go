package sso

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/filter"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type UserDataSource struct {
	client *management.APIClient
	region model.RegionMapping
}

type UserDataSourceModel struct {
	Id            types.String `tfsdk:"id"`
	UserId        types.String `tfsdk:"user_id"`
	EnvironmentId types.String `tfsdk:"environment_id"`
	Username      types.String `tfsdk:"username"`
	Email         types.String `tfsdk:"email"`
	Status        types.String `tfsdk:"status"`
	PopulationId  types.String `tfsdk:"population_id"`
}

// Framework interfaces
var (
	_ datasource.DataSource = &UserDataSource{}
)

// New Object
func NewUserDataSource() datasource.DataSource {
	return &UserDataSource{}
}

// Metadata
func (r *UserDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

// Schema
func (r *UserDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	const attrMinLength = 1

	statusDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The enabled status of the user.",
	).AllowedValues("ENABLED", "DISABLED")

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Datasource to read PingOne user data.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to create the user in."),
			),

			"user_id": schema.StringAttribute{
				Description: "The ID of the user.",
				Optional:    true,
				Computed:    true,

				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(
						path.MatchRoot("user_id"),
						path.MatchRoot("username"),
						path.MatchRoot("email"),
					),
					verify.P1ResourceIDValidator(),
				},
			},

			"username": schema.StringAttribute{
				Description: "The username of the user.",
				Optional:    true,
				Computed:    true,

				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(
						path.MatchRoot("user_id"),
						path.MatchRoot("username"),
						path.MatchRoot("email"),
					),
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			"email": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The email address of the user.").Description,
				Optional:    true,
				Computed:    true,

				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(
						path.MatchRoot("user_id"),
						path.MatchRoot("username"),
						path.MatchRoot("email"),
					),
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			"status": schema.StringAttribute{
				Description:         statusDescription.Description,
				MarkdownDescription: statusDescription.MarkdownDescription,
				Computed:            true,
			},

			"population_id": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The population ID the user is assigned to.").Description,
				Computed:    true,
			},
		},
	}
}

func (r *UserDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	preparedClient, err := PrepareClient(ctx, resourceConfig)
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

func (r *UserDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *UserDataSourceModel

	if r.client == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var user management.User
	var scimFilter string

	if !data.Username.IsNull() {

		scimFilter = filter.BuildScimFilter(
			append(make([]interface{}, 0), map[string]interface{}{
				"name":   "username",
				"values": []string{data.Username.ValueString()},
			}), map[string]string{})

	} else if !data.UserId.IsNull() {

		scimFilter = filter.BuildScimFilter(
			append(make([]interface{}, 0), map[string]interface{}{
				"name":   "id",
				"values": []string{data.UserId.ValueString()},
			}), map[string]string{})

	} else if !data.Email.IsNull() {

		scimFilter = filter.BuildScimFilter(
			append(make([]interface{}, 0), map[string]interface{}{
				"name":   "email",
				"values": []string{data.Email.ValueString()},
			}), map[string]string{})

	} else {
		resp.Diagnostics.AddError(
			"Missing parameter",
			"Cannot find the requested user. user_id, username or email must be set.",
		)
		return
	}

	var response *management.EntityArray
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return r.client.UsersApi.ReadAllUsers(ctx, data.EnvironmentId.ValueString()).Filter(scimFilter).Execute()
		},
		"ReadAllUsers",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		&response,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var responseEnabled *management.UserEnabled
	if users, ok := response.Embedded.GetUsersOk(); ok && len(users) > 0 && users[0].Id != nil {

		user = users[0]

		resp.Diagnostics.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				return r.client.EnableUsersApi.ReadUserEnabled(ctx, data.EnvironmentId.ValueString(), user.GetId()).Execute()
			},
			"ReadUserEnabled",
			framework.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			&responseEnabled,
		)...)

	} else {
		resp.Diagnostics.AddError(
			"Cannot find user",
			"Cannot find the requested user from the provided values. Please check the user_id, username or email parameters.",
		)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(&user, responseEnabled)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (p *UserDataSourceModel) toState(apiObject *management.User, apiObjectEnabled *management.UserEnabled) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil || apiObjectEnabled == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.StringOkToTF(apiObject.GetIdOk())
	p.UserId = framework.StringOkToTF(apiObject.GetIdOk())
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
