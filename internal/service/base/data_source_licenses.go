// Copyright © 2026 Ping Identity Corporation

package base

import (
	"context"
	"fmt"
	"net/http"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/legacysdk"
)

// Types
type LicensesDataSource serviceClientType

type LicensesDataSourceModel struct {
	OrganizationId pingonetypes.ResourceIDValue `tfsdk:"organization_id"`
	Id             pingonetypes.ResourceIDValue `tfsdk:"id"`
	ScimFilter     types.String                 `tfsdk:"scim_filter"`
	DataFilters    types.List                   `tfsdk:"data_filters"`
	Ids            types.List                   `tfsdk:"ids"`
}

// Framework interfaces
var (
	_ datasource.DataSource = &LicensesDataSource{}
)

// New Object
func NewLicensesDataSource() datasource.DataSource {
	return &LicensesDataSource{}
}

// Metadata
func (r *LicensesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_licenses"
}

// Schema
func (r *LicensesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	filterableAttributes := []string{"name", "package", "status"}

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Datasource to retrieve multiple PingOne license IDs selected by a SCIM filter or a name/value list combination.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"organization_id": framework.Attr_LinkID(framework.SchemaAttributeDescriptionFromMarkdown(
				"The ID of the organization to retrieve licenses for.",
			)),

			"scim_filter": framework.Attr_SCIMFilter(framework.SchemaAttributeDescriptionFromMarkdown(
				"A SCIM filter to apply to the license selection.  A SCIM filter offers the greatest flexibility in filtering licenses.",
			).AppendMarkdownString(fmt.Sprintf("If the attribute filter is `status`, available values are `%s`, `%s`, `%s` and `%s`.", management.ENUMLICENSESTATUS_ACTIVE, management.ENUMLICENSESTATUS_EXPIRED, management.ENUMLICENSESTATUS_FUTURE, management.ENUMLICENSESTATUS_TERMINATED)),
				filterableAttributes,
				[]string{"scim_filter", "data_filters"},
			),

			"data_filters": framework.Attr_DataFilter(framework.SchemaAttributeDescriptionFromMarkdown(
				"Individual data filters to apply to the license selection.",
			).AppendMarkdownString(fmt.Sprintf("If the attribute filter is `status`, available values are `%s`, `%s`, `%s` and `%s`.", management.ENUMLICENSESTATUS_ACTIVE, management.ENUMLICENSESTATUS_EXPIRED, management.ENUMLICENSESTATUS_FUTURE, management.ENUMLICENSESTATUS_TERMINATED)),
				filterableAttributes,
				[]string{"scim_filter", "data_filters"},
			),

			"ids": framework.Attr_DataSourceReturnIDs(framework.SchemaAttributeDescriptionFromMarkdown(
				"The list of resulting IDs of licenses that have been successfully retrieved and filtered.",
			)),
		},
	}
}

func (r *LicensesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	resourceConfig, ok := req.ProviderData.(legacysdk.ResourceType)
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

func (r *LicensesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *LicensesDataSourceModel

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

	var filterFunction management.ApiReadAllLicensesRequest

	if !data.ScimFilter.IsNull() {

		filterFunction = r.Client.ManagementAPIClient.LicensesApi.ReadAllLicenses(ctx, data.OrganizationId.ValueString()).Filter(data.ScimFilter.ValueString())

	} else if !data.DataFilters.IsNull() {

		filterFunction = r.Client.ManagementAPIClient.LicensesApi.ReadAllLicenses(ctx, data.OrganizationId.ValueString())

	} else {
		resp.Diagnostics.AddError(
			"Missing parameter",
			"Cannot find the requested DaVinci flow policies. scim_filter or data_filter must be set.",
		)
		return
	}

	var licenses []management.License
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			pagedIterator := filterFunction.Execute()

			licenses := make([]management.License, 0)

			var initialHttpResponse *http.Response

			for pageCursor, err := range pagedIterator {
				if err != nil {
					return nil, pageCursor.HTTPResponse, err
				}

				if initialHttpResponse == nil {
					initialHttpResponse = pageCursor.HTTPResponse
				}

				licenses = append(licenses, pageCursor.EntityArray.Embedded.GetLicenses()...)

			}

			return licenses, initialHttpResponse, nil
		},
		"ReadAllLicenses",
		legacysdk.DefaultCustomError,
		nil,
		&licenses,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !data.DataFilters.IsNull() {
		var dataFilterPlan []framework.DataFilterModel
		resp.Diagnostics.Append(data.DataFilters.ElementsAs(ctx, &dataFilterPlan, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		var d diag.Diagnostics
		licenses, d = filterResults(ctx, dataFilterPlan, licenses)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(data.OrganizationId.ValueString(), licenses)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (p *LicensesDataSourceModel) toState(environmentID string, licenses []management.License) diag.Diagnostics {
	var diags diag.Diagnostics

	if licenses == nil || environmentID == "" {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	list := make([]string, 0)
	for _, item := range licenses {
		list = append(list, item.GetId())
	}

	var d diag.Diagnostics

	p.Id = framework.PingOneResourceIDToTF(environmentID)
	p.Ids, d = framework.StringSliceToTF(list)
	diags.Append(d...)

	return diags
}

func filterResults(ctx context.Context, filterPlan []framework.DataFilterModel, licenses []management.License) ([]management.License, diag.Diagnostics) {
	var diags diag.Diagnostics

	items := make([]management.License, 0)

	for _, license := range licenses {

		filterMap := map[string]string{
			"name":    license.GetName(),
			"package": license.GetPackage(),
			"status":  string(license.GetStatus()),
		}

		include := true

		for _, filter := range filterPlan {

			for k, v := range filterMap {
				if filter.Name.ValueString() == k {

					var filterValuesPlan []types.String
					diags.Append(filter.Values.ElementsAs(ctx, &filterValuesPlan, false)...)
					if diags.HasError() {
						return nil, diags
					}

					filterValues, d := framework.TFTypeStringSliceToStringSlice(filterValuesPlan, path.Root("data_filters"))
					diags.Append(d...)
					if diags.HasError() {
						return nil, diags
					}

					if !slices.Contains(filterValues, v) {
						include = false
					}
				}
			}

		}

		if include {
			items = append(items, license)
		}

	}

	return items, diags
}
