// Copyright © 2026 Ping Identity Corporation

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
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/legacysdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
)

var (
	_ datasource.DataSource              = &systemApplicationsDataSource{}
	_ datasource.DataSourceWithConfigure = &systemApplicationsDataSource{}
)

func NewSystemApplicationsDataSource() datasource.DataSource {
	return &systemApplicationsDataSource{}
}

type systemApplicationsDataSource serviceClientType

func (r *systemApplicationsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_system_applications"
}

func (r *systemApplicationsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

type systemApplicationsDataSourceModel struct {
	EnvironmentId pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	Id            pingonetypes.ResourceIDValue `tfsdk:"id"`
	Ids           types.List                   `tfsdk:"ids"`
}

func (r *systemApplicationsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Data source to retrieve built-in system applications (PingOne Self-Service or PingOne Portal) from PingOne",
		Attributes: map[string]schema.Attribute{
			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to read system applications from."),
			),
			"id": framework.Attr_ID(),
			"ids": schema.ListAttribute{
				ElementType: pingonetypes.ResourceIDType{},
				Computed:    true,
				Description: "The list of IDs of system applications that have been successfully retrieved.",
			},
		},
	}
}

func (state *systemApplicationsDataSourceModel) readClientResponse(response []string) diag.Diagnostics {
	var respDiags diag.Diagnostics
	state.Id = framework.PingOneResourceIDToTF(state.EnvironmentId.ValueString())
	state.Ids, respDiags = framework.StringSliceToTF(response)
	return respDiags
}

func (r *systemApplicationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data systemApplicationsDataSourceModel

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

	// Run the API call
	var responseData []string
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			pagedIterator := r.Client.ManagementAPIClient.ApplicationsApi.ReadAllApplications(ctx, data.EnvironmentId.ValueString()).Execute()

			var initialHttpResponse *http.Response
			foundIDs := make([]string, 0)

			for pageCursor, err := range pagedIterator {
				if err != nil {
					return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, pageCursor.HTTPResponse, err)
				}

				if initialHttpResponse == nil {
					initialHttpResponse = pageCursor.HTTPResponse
				}

				if applications, ok := pageCursor.EntityArray.Embedded.GetApplicationsOk(); ok {

					var applicationObj management.ReadOneApplication200Response
					for _, applicationObj = range applications {
						applicationInstance := applicationObj.GetActualInstance()

						switch v := applicationInstance.(type) {
						case *management.ApplicationPingOnePortal:
							if id, ok := v.GetIdOk(); ok {
								foundIDs = append(foundIDs, *id)
							}
						case *management.ApplicationPingOneSelfService:
							if id, ok := v.GetIdOk(); ok {
								foundIDs = append(foundIDs, *id)
							}
						}
					}
				}
			}

			return foundIDs, initialHttpResponse, nil
		},
		"ReadAllApplications",
		legacysdk.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		&responseData,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
