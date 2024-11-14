package base

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
)

// Types
type PhoneDeliverySettingsListDataSource serviceClientType

type PhoneDeliverySettingsListDataSourceModel struct {
	Id            pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	Ids           types.List                   `tfsdk:"ids"`
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

			"environment_id": framework.Attr_LinkID(framework.SchemaAttributeDescriptionFromMarkdown(
				"The ID of the environment to filter phone delivery settings senders from.",
			)),

			"ids": framework.Attr_DataSourceReturnIDs(framework.SchemaAttributeDescriptionFromMarkdown(
				"The list of resulting IDs of phone delivery settings senders that have been successfully retrieved.",
			)),
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

	r.Client = resourceConfig.Client.API
	if r.Client == nil {
		resp.Diagnostics.AddError(
			"Client not initialised",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.",
		)
		return
	}
}

func (r *PhoneDeliverySettingsListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *PhoneDeliverySettingsListDataSourceModel

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

	var response []string
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			pagedIterator := r.Client.ManagementAPIClient.PhoneDeliverySettingsApi.ReadAllPhoneDeliverySettings(ctx, data.EnvironmentId.ValueString()).Execute()

			phoneDeliverySettingsIds := make([]string, 0)

			var initialHttpResponse *http.Response

			for pageCursor, err := range pagedIterator {
				if err != nil {
					return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, pageCursor.HTTPResponse, err)
				}

				if initialHttpResponse == nil {
					initialHttpResponse = pageCursor.HTTPResponse
				}

				if pageCursor.EntityArray.Embedded != nil && pageCursor.EntityArray.Embedded.PhoneDeliverySettings != nil {
					for _, item := range pageCursor.EntityArray.Embedded.GetPhoneDeliverySettings() {
						if v := item.NotificationsSettingsPhoneDeliverySettingsCustom; v != nil {
							phoneDeliverySettingsIds = append(phoneDeliverySettingsIds, v.GetId())
						}

						if v := item.NotificationsSettingsPhoneDeliverySettingsTwilioSyniverse; v != nil {
							phoneDeliverySettingsIds = append(phoneDeliverySettingsIds, v.GetId())
						}
					}
				}

			}

			return phoneDeliverySettingsIds, initialHttpResponse, nil
		},
		"ReadAllPhoneDeliverySettings",
		framework.DefaultCustomError,
		nil,
		&response,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(response)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (p *PhoneDeliverySettingsListDataSourceModel) toState(apiObject []string) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	var d diag.Diagnostics

	p.Id = p.EnvironmentId

	p.Ids, d = framework.StringSliceToTF(apiObject)
	diags.Append(d...)

	return diags
}
