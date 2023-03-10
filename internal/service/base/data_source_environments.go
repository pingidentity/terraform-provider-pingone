package base

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
)

// Types
type EnvironmentsDataSource struct {
	client *management.APIClient
	region model.RegionMapping
}

type EnvironmentsDataSourceModel struct {
	Id         types.String `tfsdk:"id"`
	ScimFilter types.String `tfsdk:"scim_filter"`
	Ids        types.List   `tfsdk:"ids"`
}

// Framework interfaces
var (
	_ datasource.DataSource = &EnvironmentsDataSource{}
)

// New Object
func NewEnvironmentsDataSource() datasource.DataSource {
	return &EnvironmentsDataSource{}
}

// Metadata
func (r *EnvironmentsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_environments"
}

// Schema
func (r *EnvironmentsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	scimFilterDescriptionFmt := "A SCIM filter to apply to the environment selection.  A SCIM filter offers the greatest flexibility in filtering environments.  SCIM operators can be used in the following ways: `sw` (starts with) supports the `name` attribute; `eq` (equal to) supports the `id`, `organization.id`, `license.id` attributes; `and` (logical AND) can be used to aggregate conditions.  For example, `(name sw \"TEST-\") AND (license.id eq \"${var.license_id}\")`"
	scimFilterDescription := framework.SchemaDescription{
		MarkdownDescription: scimFilterDescriptionFmt,
		Description:         strings.ReplaceAll(scimFilterDescriptionFmt, "`", "\""),
	}

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Datasource to retrieve multiple PingOne environments.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"scim_filter": schema.StringAttribute{
				Description:         scimFilterDescription.Description,
				MarkdownDescription: scimFilterDescription.MarkdownDescription,
				Required:            true,
			},

			"ids": framework.Attr_DataSourceReturnIDs(framework.SchemaDescription{
				Description: "The list of resulting IDs of environments that have been successfully retrieved and filtered.",
			}),
		},
	}
}

func (r *EnvironmentsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *EnvironmentsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *EnvironmentsDataSourceModel

	if r.client == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": r.region.URLSuffix,
	})

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	filterFunction := func() (interface{}, *http.Response, error) {
		return r.client.EnvironmentsApi.ReadAllEnvironments(ctx).Filter(data.ScimFilter.ValueString()).Execute()
	}

	response, diags := framework.ParseResponse(
		ctx,

		filterFunction,
		"ReadAllEnvironments",
		framework.DefaultCustomError,
		sdk.DefaultRetryable,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	entityArray := response.(*management.EntityArray)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(entityArray.Embedded.GetEnvironments())...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (p *EnvironmentsDataSourceModel) toState(environments []management.Environment) diag.Diagnostics {
	var diags diag.Diagnostics

	if environments == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	list := make([]string, 0)
	for _, item := range environments {
		list = append(list, item.GetId())
	}

	var d diag.Diagnostics

	if p.Id.IsNull() {
		p.Id = framework.StringToTF(uuid.New().String())
	}

	p.Ids, d = framework.StringSliceToTF(list)
	diags.Append(d...)

	return diags
}
