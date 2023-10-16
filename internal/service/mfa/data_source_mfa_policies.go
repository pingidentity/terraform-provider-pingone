package mfa

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/mfa"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
)

// Types
type MFAPoliciesDataSource serviceClientType

type MFAPoliciesDataSourceModel struct {
	Id            types.String `tfsdk:"id"`
	EnvironmentId types.String `tfsdk:"environment_id"`
	Ids           types.List   `tfsdk:"ids"`
}

// Framework interfaces
var (
	_ datasource.DataSource = &MFAPoliciesDataSource{}
)

// New Object
func NewMFAPoliciesDataSource() datasource.DataSource {
	return &MFAPoliciesDataSource{}
}

// Metadata
func (r *MFAPoliciesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mfa_policies"
}

// Schema
func (r *MFAPoliciesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Datasource to retrieve the IDs of multiple PingOne MFA Device policies.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to select MFA device policies from."),
			),

			"ids": framework.Attr_DataSourceReturnIDs(framework.SchemaAttributeDescription{
				Description: "The list of resulting IDs of MFA Device Policies that have been successfully retrieved and filtered.",
			}),
		},
	}
}

func (r *MFAPoliciesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *MFAPoliciesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *MFAPoliciesDataSourceModel

	if r.Client.MFAAPIClient == nil {
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

	filterFunction := func() (any, *http.Response, error) {
		return r.Client.MFAAPIClient.DeviceAuthenticationPolicyApi.ReadDeviceAuthenticationPolicies(ctx, data.EnvironmentId.ValueString()).Execute()
	}

	var entityArray *mfa.EntityArray
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		filterFunction,
		"ReadDeviceAuthenticationPolicies",
		framework.DefaultCustomError,
		nil,
		&entityArray,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(entityArray.Embedded.GetDeviceAuthenticationPolicies())...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (p *MFAPoliciesDataSourceModel) toState(apiObject []mfa.DeviceAuthenticationPolicy) diag.Diagnostics {
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

	if p.Id.IsNull() {
		p.Id = framework.StringToTF(uuid.New().String())
	}

	p.EnvironmentId = framework.StringToTF(p.EnvironmentId.ValueString())

	p.Ids, d = framework.StringSliceToTF(list)
	diags.Append(d...)

	return diags
}
