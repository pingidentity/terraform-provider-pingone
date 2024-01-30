package sso

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
)

// Types
type ApplicationSignOnPolicyAssignmentsDataSource serviceClientType

type ApplicationSignOnPolicyAssignmentsDataSourceModel struct {
	EnvironmentId types.String `tfsdk:"environment_id"`
	ApplicationId types.String `tfsdk:"application_id"`
	Ids           types.List   `tfsdk:"ids"`
}

// Framework interfaces
var (
	_ datasource.DataSource = &ApplicationSignOnPolicyAssignmentsDataSource{}
)

// New Object
func NewApplicationSignOnPolicyAssignmentsDataSource() datasource.DataSource {
	return &ApplicationSignOnPolicyAssignmentsDataSource{}
}

// Metadata
func (r *ApplicationSignOnPolicyAssignmentsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application_sign_on_policy_assignments"
}

// Schema
func (r *ApplicationSignOnPolicyAssignmentsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Datasource to retrieve the IDs, as a collection, of PingOne Sign On Policy assignments for an application in an environment.",

		Attributes: map[string]schema.Attribute{
			"environment_id": framework.Attr_LinkID(framework.SchemaAttributeDescriptionFromMarkdown(
				"The ID of the environment to filter application sign on policy assignments from.",
			)),

			"application_id": framework.Attr_LinkID(framework.SchemaAttributeDescriptionFromMarkdown(
				"The ID of the application to filter application sign on policy assignments from.",
			)),

			"ids": framework.Attr_DataSourceReturnIDs(framework.SchemaAttributeDescriptionFromMarkdown(
				"The list of resulting IDs of application sign on policy assignments that have been successfully retrieved for an application.",
			)),
		},
	}
}

func (r *ApplicationSignOnPolicyAssignmentsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *ApplicationSignOnPolicyAssignmentsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *ApplicationSignOnPolicyAssignmentsDataSourceModel

	if r.Client.ManagementAPIClient == nil {
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

	var entityArray *management.EntityArray
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return r.Client.ManagementAPIClient.ApplicationSignOnPolicyAssignmentsApi.ReadAllSignOnPolicyAssignments(ctx, data.EnvironmentId.ValueString(), data.ApplicationId.ValueString()).Execute()
		},
		"ReadAllSignOnPolicyAssignments",
		framework.DefaultCustomError,
		nil,
		&entityArray,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(entityArray.Embedded.GetSignOnPolicyAssignments())...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (p *ApplicationSignOnPolicyAssignmentsDataSourceModel) toState(apiObject []management.SignOnPolicyAssignment) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	list := make([]string, 0)
	for _, item := range apiObject {
		list = append(list, item.GetId())
	}

	var d diag.Diagnostics

	p.Ids, d = framework.StringSliceToTF(list)
	diags.Append(d...)

	return diags
}
