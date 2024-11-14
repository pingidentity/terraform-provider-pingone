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
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
)

// Types
type MFADevicePoliciesDataSource serviceClientType

type MFADevicePoliciesDataSourceModel struct {
	Id            pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	Ids           types.List                   `tfsdk:"ids"`
}

// Framework interfaces
var (
	_ datasource.DataSource = &MFADevicePoliciesDataSource{}
)

// New Object
func NewMFADevicePoliciesDataSource() datasource.DataSource {
	return &MFADevicePoliciesDataSource{}
}

// Metadata
func (r *MFADevicePoliciesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mfa_device_policies"
}

// Schema
func (r *MFADevicePoliciesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Datasource to retrieve the IDs of multiple PingOne MFA Device policies.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to select MFA device policies from."),
			),

			"ids": framework.Attr_DataSourceReturnIDs(framework.SchemaAttributeDescriptionFromMarkdown(
				"The list of resulting IDs of MFA Device Policies that have been successfully retrieved and filtered.",
			)),
		},
	}
}

func (r *MFADevicePoliciesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *MFADevicePoliciesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *MFADevicePoliciesDataSourceModel

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
		pagedIterator := r.Client.MFAAPIClient.DeviceAuthenticationPolicyApi.ReadDeviceAuthenticationPolicies(ctx, data.EnvironmentId.ValueString()).Execute()

		devicePolicyIDs := make([]string, 0)

		var initialHttpResponse *http.Response

		for pageCursor, err := range pagedIterator {
			if err != nil {
				return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, pageCursor.HTTPResponse, err)
			}

			if initialHttpResponse == nil {
				initialHttpResponse = pageCursor.HTTPResponse
			}

			if pageCursor.EntityArray.Embedded != nil && pageCursor.EntityArray.Embedded.DeviceAuthenticationPolicies != nil {
				for _, policy := range pageCursor.EntityArray.Embedded.GetDeviceAuthenticationPolicies() {
					devicePolicyIDs = append(devicePolicyIDs, policy.GetId())
				}
			}
		}

		return nil, initialHttpResponse, nil
	}

	var devicePolicyIDs []string
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		filterFunction,
		"ReadDeviceAuthenticationPolicies",
		framework.DefaultCustomError,
		nil,
		&devicePolicyIDs,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(devicePolicyIDs)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (p *MFADevicePoliciesDataSourceModel) toState(apiObject []string) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	var d diag.Diagnostics

	if p.Id.IsNull() {
		p.Id = framework.PingOneResourceIDToTF(uuid.New().String())
	}

	p.EnvironmentId = framework.PingOneResourceIDToTF(p.EnvironmentId.ValueString())

	p.Ids, d = framework.StringSliceToTF(apiObject)
	diags.Append(d...)

	return diags
}
