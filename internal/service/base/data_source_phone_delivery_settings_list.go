package base

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
type PhoneDeliverySettingsListDataSource serviceClientType

type PhoneDeliverySettingsListDataSourceModel struct {
	Id            types.String `tfsdk:"id"`
	EnvironmentId types.String `tfsdk:"environment_id"`
	Ids           types.List   `tfsdk:"ids"`
}

// Framework interfaces
var (
	_ datasource.DataSource = &PhoneDeliverySettingsListDataSource{}
)

// New Object
func NewPhoneDeliverySettingsListDataSource() datasource.DataSource {
	return &PhoneDeliverySettingsListDataSource{}
}

// Metadata
func (r *PhoneDeliverySettingsListDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_phone_delivery_settings_list"
}

// Schema
func (r *PhoneDeliverySettingsListDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Datasource to retrieve multiple phone delivery settings in a PingOne environments.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(framework.SchemaAttributeDescription{
				Description: "The ID of the environment to filter phone delivery settings senders from.",
			}),

			"ids": framework.Attr_DataSourceReturnIDs(framework.SchemaAttributeDescription{
				Description: "The list of resulting IDs of phone delivery settings senders that have been successfully retrieved.",
			}),
		},
	}
}

func (r *PhoneDeliverySettingsListDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.Client = preparedClient
}

func (r *PhoneDeliverySettingsListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *PhoneDeliverySettingsListDataSourceModel

	if r.Client == nil {
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
			return r.Client.PhoneDeliverySettingsApi.ReadAllPhoneDeliverySettings(ctx, data.EnvironmentId.ValueString()).Execute()
		},
		"ReadAllPhoneDeliverySettings",
		framework.DefaultCustomError,
		nil,
		&entityArray,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(entityArray.Embedded.GetPhoneDeliverySettings())...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (p *PhoneDeliverySettingsListDataSourceModel) toState(apiObject []management.NotificationsSettingsPhoneDeliverySettings) diag.Diagnostics {
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
		if v := item.NotificationsSettingsPhoneDeliverySettingsCustom; v != nil {
			list = append(list, v.GetId())
		}

		if v := item.NotificationsSettingsPhoneDeliverySettingsTwilioSyniverse; v != nil {
			list = append(list, v.GetId())
		}
	}

	var d diag.Diagnostics

	p.Id = p.EnvironmentId

	p.Ids, d = framework.StringSliceToTF(list)
	diags.Append(d...)

	return diags
}
