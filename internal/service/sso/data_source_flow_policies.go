package sso

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/filter"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/davincitypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
)

// Types
type FlowPoliciesDataSource serviceClientType

type FlowPoliciesDataSourceModel struct {
	EnvironmentId pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	Id            pingonetypes.ResourceIDValue `tfsdk:"id"`
	ScimFilter    types.String                 `tfsdk:"scim_filter"`
	DataFilters   types.List                   `tfsdk:"data_filters"`
	Ids           types.List                   `tfsdk:"ids"`
}

// Framework interfaces
var (
	_ datasource.DataSource = &FlowPoliciesDataSource{}
)

// New Object
func NewFlowPoliciesDataSource() datasource.DataSource {
	return &FlowPoliciesDataSource{}
}

// Metadata
func (r *FlowPoliciesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_flow_policies"
}

// Schema
func (r *FlowPoliciesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	filterableAttributes := []string{"trigger.type"}

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Datasource to retrieve multiple PingOne DaVinci Flow Policies.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(framework.SchemaAttributeDescriptionFromMarkdown(
				"The ID of the environment to filter DaVinci flow policies from.",
			)),

			"scim_filter": framework.Attr_SCIMFilter(framework.SchemaAttributeDescriptionFromMarkdown(
				"A SCIM filter to apply to the DaVinci flow policy selection.  A SCIM filter offers the greatest flexibility in filtering DaVinci flow policies.",
			),
				filterableAttributes,
				[]string{"scim_filter", "data_filters"},
			),

			"data_filters": framework.Attr_DataFilter(framework.SchemaAttributeDescriptionFromMarkdown(
				"Individual data filters to apply to the DaVinci flow policy selection.",
			),
				filterableAttributes,
				[]string{"scim_filter", "data_filters"},
			),

			"ids": framework.Attr_DataSourceReturnIDsByElement(framework.SchemaAttributeDescriptionFromMarkdown(
				"The list of resulting IDs of DaVinci flow policies that have been successfully retrieved and filtered.",
			),
				davincitypes.ResourceIDType{}),
		},
	}
}

func (r *FlowPoliciesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *FlowPoliciesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *FlowPoliciesDataSourceModel

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

	var filterFunction func() management.EntityArrayPagedIterator

	if !data.ScimFilter.IsNull() {

		filterFunction = r.Client.ManagementAPIClient.FlowPoliciesApi.ReadAllFlowPolicies(ctx, data.EnvironmentId.ValueString()).Filter(data.ScimFilter.ValueString()).Execute

	} else if !data.DataFilters.IsNull() {

		var dataFilterIn []framework.DataFilterModel
		resp.Diagnostics.Append(data.DataFilters.ElementsAs(ctx, &dataFilterIn, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		filterSet := make([]interface{}, 0)

		for _, v := range dataFilterIn {

			values := framework.TFListToStringSlice(ctx, v.Values)
			tflog.Debug(ctx, "Filter set loop", map[string]interface{}{
				"name":          v.Name.ValueString(),
				"len(elements)": fmt.Sprintf("%d", len(v.Values.Elements())),
				"len(values)":   fmt.Sprintf("%d", len(values)),
			})
			filterSet = append(filterSet, map[string]interface{}{
				"name":   v.Name.ValueString(),
				"values": values,
			})
		}

		scimFilter := filter.BuildScimFilter(filterSet, map[string]string{})

		tflog.Debug(ctx, "SCIM Filter", map[string]interface{}{
			"scimFilter": scimFilter,
		})

		filterFunction = r.Client.ManagementAPIClient.FlowPoliciesApi.ReadAllFlowPolicies(ctx, data.EnvironmentId.ValueString()).Filter(scimFilter).Execute

	} else {
		resp.Diagnostics.AddError(
			"Missing parameter",
			"Cannot find the requested DaVinci flow policies. scim_filter or data_filter must be set.",
		)
		return
	}

	var flowPolicyIDs []string
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			pagedIterator := filterFunction()

			var initialHttpResponse *http.Response

			foundIDs := make([]string, 0)

			for pageCursor, err := range pagedIterator {
				if err != nil {
					return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, pageCursor.HTTPResponse, err)
				}

				if initialHttpResponse == nil {
					initialHttpResponse = pageCursor.HTTPResponse
				}

				if pageCursor.EntityArray.Embedded != nil && pageCursor.EntityArray.Embedded.FlowPolicies != nil {
					for _, flowPolicy := range pageCursor.EntityArray.Embedded.GetFlowPolicies() {
						foundIDs = append(foundIDs, flowPolicy.GetId())
					}
				}
			}

			return foundIDs, initialHttpResponse, nil
		},
		"ReadAllFlowPolicies",
		framework.DefaultCustomError,
		nil,
		&flowPolicyIDs,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(data.EnvironmentId.ValueString(), flowPolicyIDs)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (p *FlowPoliciesDataSourceModel) toState(environmentID string, flowPolicyIDs []string) diag.Diagnostics {
	var diags diag.Diagnostics

	if flowPolicyIDs == nil || environmentID == "" {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	var d diag.Diagnostics

	p.Id = framework.PingOneResourceIDToTF(environmentID)
	p.Ids = framework.DaVinciResourceIDListToTF(flowPolicyIDs)
	diags.Append(d...)

	return diags
}
