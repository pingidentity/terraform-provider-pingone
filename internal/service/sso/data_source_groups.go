package sso

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/filter"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type GroupsDataSource serviceClientType

type GroupsDataSourceModel struct {
	EnvironmentId types.String `tfsdk:"environment_id"`
	Id            types.String `tfsdk:"id"`
	ScimFilter    types.String `tfsdk:"scim_filter"`
	DataFilters   types.List   `tfsdk:"data_filters"`
	Ids           types.List   `tfsdk:"ids"`
}

// Framework interfaces
var (
	_ datasource.DataSource = &GroupsDataSource{}
)

// New Object
func NewGroupsDataSource() datasource.DataSource {
	return &GroupsDataSource{}
}

// Metadata
func (r *GroupsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_groups"
}

// Schema
func (r *GroupsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	filterableAttributes := []string{"id", "name", "population.id", "externalId"}

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Datasource to filter and retrieve multiple PingOne groups in an environment.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the ID of the environment to filter groups from.  Must be a valid PingOne resource ID.").Description,
				Required:    true,
				Validators: []validator.String{
					verify.P1ResourceIDValidator(),
				},
			},

			"scim_filter": framework.Attr_SCIMFilter(framework.SchemaAttributeDescriptionFromMarkdown(
				"A SCIM filter to apply to the group selection.  A SCIM filter offers the greatest flexibility in filtering groups.",
			),
				filterableAttributes,
				[]string{"scim_filter", "data_filters"},
			),

			"data_filters": framework.Attr_DataFilter(framework.SchemaAttributeDescriptionFromMarkdown(
				"Individual data filters to apply to the group selection.",
			),
				filterableAttributes,
				[]string{"scim_filter", "data_filters"},
			),

			"ids": framework.Attr_DataSourceReturnIDs(framework.SchemaAttributeDescriptionFromMarkdown(
				"The list of resulting IDs of groups that have been successfully retrieved and filtered.",
			)),
		},
	}
}

func (r *GroupsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *GroupsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *GroupsDataSourceModel

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

	var filterFunction sdk.SDKInterfaceFunc

	if !data.ScimFilter.IsNull() {

		filterFunction = func() (any, *http.Response, error) {
			return r.Client.ManagementAPIClient.GroupsApi.ReadAllGroups(ctx, data.EnvironmentId.ValueString()).Filter(data.ScimFilter.ValueString()).Execute()
		}

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

		filterFunction = func() (any, *http.Response, error) {
			return r.Client.ManagementAPIClient.GroupsApi.ReadAllGroups(ctx, data.EnvironmentId.ValueString()).Filter(scimFilter).Execute()
		}

	} else {
		resp.Diagnostics.AddError(
			"Missing parameter",
			"Cannot find the requested group. group_id or name must be set.",
		)
		return
	}

	var entityArray *management.EntityArray
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		filterFunction,
		"ReadAllGroups",
		framework.DefaultCustomError,
		nil,
		&entityArray,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(data.EnvironmentId.ValueString(), entityArray.Embedded.GetGroups())...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (p *GroupsDataSourceModel) toState(environmentID string, groups []management.Group) diag.Diagnostics {
	var diags diag.Diagnostics

	if groups == nil || environmentID == "" {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	list := make([]string, 0)
	for _, item := range groups {
		list = append(list, item.GetId())
	}

	var d diag.Diagnostics

	p.Id = framework.StringToTF(environmentID)
	p.Ids, d = framework.StringSliceToTF(list)
	diags.Append(d...)

	return diags
}
