// Copyright © 2026 Ping Identity Corporation

package credentials

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/credentials"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/legacysdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
)

// Types
type CredentialTypeVersionDataSource serviceClientType

type CredentialTypeVersionDataSourceModel struct {
	Id                      pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId           pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	CredentialTypeId        pingonetypes.ResourceIDValue `tfsdk:"credential_type_id"`
	CredentialTypeVersionId pingonetypes.ResourceIDValue `tfsdk:"credential_type_version_id"`
	Number                  types.Int32                  `tfsdk:"number"`
	CreatedAt               timetypes.RFC3339            `tfsdk:"created_at"`
}

// Framework interfaces
var (
	_ datasource.DataSource = &CredentialTypeVersionDataSource{}
)

// New Object
func NewCredentialTypeVersionDataSource() datasource.DataSource {
	return &CredentialTypeVersionDataSource{}
}

// Metadata
func (r *CredentialTypeVersionDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_credential_type_version"
}

func (r *CredentialTypeVersionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Datasource to retrieve a single version of a PingOne Credentials credential type by its version ID.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment in which the credential type exists."),
			),

			"credential_type_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the credential type to retrieve the version for."),
			),

			"credential_type_version_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the credential type version to retrieve."),
			),

			"number": schema.Int32Attribute{
				Description: "Version number of the credential type version.",
				Computed:    true,
			},

			"created_at": schema.StringAttribute{
				Description: "Date and time the credential type version was created.",
				Computed:    true,

				CustomType: timetypes.RFC3339Type{},
			},
		},
	}
}

func (r *CredentialTypeVersionDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *CredentialTypeVersionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *CredentialTypeVersionDataSourceModel

	if r.Client == nil || r.Client.CredentialsAPIClient == nil {
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
	var response *credentials.CredentialTypeVersion
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.CredentialsAPIClient.CredentialTypesApi.ReadOneCredentialTypeVersion(ctx, data.EnvironmentId.ValueString(), data.CredentialTypeId.ValueString(), data.CredentialTypeVersionId.ValueString()).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadOneCredentialTypeVersion",
		legacysdk.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		&response,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(data.EnvironmentId.ValueString(), response)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (p *CredentialTypeVersionDataSourceModel) toState(environmentID string, apiObject *credentials.CredentialTypeVersion) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil || environmentID == "" {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.PingOneResourceIDToTF(environmentID)
	p.CredentialTypeVersionId = framework.PingOneResourceIDOkToTF(apiObject.GetIdOk())
	p.Number = framework.Int32OkToTF(apiObject.GetVersionOk())
	p.CreatedAt = framework.TimeOkToTF(apiObject.GetCreatedAtOk())

	return diags
}
