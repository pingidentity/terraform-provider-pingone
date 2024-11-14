package base

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
)

// Types
type EnvironmentsDataSource serviceClientType

type EnvironmentsDataSourceModel struct {
	Id         pingonetypes.ResourceIDValue `tfsdk:"id"`
	ScimFilter types.String                 `tfsdk:"scim_filter"`
	Ids        types.List                   `tfsdk:"ids"`
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

	scimFilterDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A SCIM filter to apply to the environment selection.  A SCIM filter offers the greatest flexibility in filtering environments.  SCIM operators can be used in the following ways: `sw` (starts with) supports the `name` attribute; `eq` (equal to) supports the `id`, `organization.id`, `license.id` attributes; `and` (logical AND) can be used to aggregate conditions.  For example, `(name sw \"TEST-\") AND (license.id eq \"${var.license_id}\")`",
	)

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

			"ids": framework.Attr_DataSourceReturnIDs(framework.SchemaAttributeDescriptionFromMarkdown(
				"The list of resulting IDs of environments that have been successfully retrieved and filtered.",
			)),
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

	r.Client = resourceConfig.Client.API
	if r.Client == nil {
		resp.Diagnostics.AddError(
			"Client not initialised",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.",
		)
		return
	}
}

func (r *EnvironmentsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *EnvironmentsDataSourceModel

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

	var environments []management.Environment
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			pagedIterator := r.Client.ManagementAPIClient.EnvironmentsApi.ReadAllEnvironments(ctx).Filter(data.ScimFilter.ValueString()).Execute()

			returnEnvironments := make([]management.Environment, 0)

			var initialHttpResponse *http.Response

			for pageCursor, err := range pagedIterator {
				if err != nil {
					return nil, pageCursor.HTTPResponse, err
				}

				if initialHttpResponse == nil {
					initialHttpResponse = pageCursor.HTTPResponse
				}

				if environments, ok := pageCursor.EntityArray.Embedded.GetEnvironmentsOk(); ok {
					returnEnvironments = append(returnEnvironments, environments...)
				}
			}

			return returnEnvironments, initialHttpResponse, nil
		},
		"ReadAllEnvironments",
		framework.DefaultCustomError,
		nil,
		&environments,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(environments)...)
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
		p.Id = framework.PingOneResourceIDToTF(uuid.New().String())
	}

	p.Ids, d = framework.StringSliceToTF(list)
	diags.Append(d...)

	return diags
}
