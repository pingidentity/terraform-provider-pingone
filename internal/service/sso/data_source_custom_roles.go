package sso

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
)

// Types
type CustomRolesDataSource serviceClientType

type CustomRolesDataSourceModel struct {
	EnvironmentId pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	Id            pingonetypes.ResourceIDValue `tfsdk:"id"`
	Ids           types.List                   `tfsdk:"ids"`
}

// Framework interfaces
var (
	_ datasource.DataSource = &CustomRolesDataSource{}
)

// New Object
func NewCustomRolesDataSource() datasource.DataSource {
	return &CustomRolesDataSource{}
}

// Metadata
func (r *CustomRolesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_custom_roles"
}

// Schema
func (r *CustomRolesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Datasource to retrieve multiple PingOne custom administrator roles.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(framework.SchemaAttributeDescriptionFromMarkdown(
				"The ID of the environment to read custom roles from.",
			)),

			"ids": framework.Attr_DataSourceReturnIDs(framework.SchemaAttributeDescriptionFromMarkdown(
				"The list of resulting IDs of custom roles that have been successfully retrieved.",
			)),
		},
	}
}

func (r *CustomRolesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *CustomRolesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *CustomRolesDataSourceModel

	if r.Client == nil || r.Client.ManagementAPIClient == nil {
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

	// Get all custom roles
	var customRoleIDs []string
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			pagedIterator := r.Client.ManagementAPIClient.CustomAdminRolesApi.ReadAllCustomAdminRoles(ctx, data.EnvironmentId.ValueString()).Execute()

			var initialHttpResponse *http.Response

			for pageCursor, err := range pagedIterator {
				if err != nil {
					return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, pageCursor.HTTPResponse, err)
				}

				if initialHttpResponse == nil {
					initialHttpResponse = pageCursor.HTTPResponse
				}

				if roles, ok := pageCursor.EntityArray.Embedded.GetRolesOk(); ok {

					for _, roleObj := range roles {
						roleInstance := roleObj.CustomAdminRole
						if roleInstance == nil {
							continue
						}
						roleId, ok := roleInstance.GetIdOk()
						if ok {
							customRoleIDs = append(customRoleIDs, *roleId)
						}
					}
				}
			}

			return nil, initialHttpResponse, nil
		},
		"ReadAllCustomAdminRoles",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		nil,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(data.EnvironmentId.ValueString(), customRoleIDs)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (p *CustomRolesDataSourceModel) toState(environmentID string, customRoleIDs []string) diag.Diagnostics {
	var diags diag.Diagnostics

	if environmentID == "" {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	var d diag.Diagnostics

	p.Id = framework.PingOneResourceIDToTF(environmentID)
	p.Ids, d = framework.StringSliceToTF(customRoleIDs)
	diags.Append(d...)

	return diags
}
