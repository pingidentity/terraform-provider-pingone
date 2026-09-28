// Copyright © 2026 Ping Identity Corporation

package credentials

import (
	"context"
	"fmt"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/credentials"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/legacysdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
)

// Types
type CredentialTypeVersionsDataSource serviceClientType

type CredentialTypeVersionsDataSourceModel struct {
	Id               pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId    pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	CredentialTypeId pingonetypes.ResourceIDValue `tfsdk:"credential_type_id"`
	Filter           types.String                 `tfsdk:"filter"`
	Versions         types.List                   `tfsdk:"versions"`
}

type CredentialTypeVersionDataSourceInnerModel struct {
	Id        pingonetypes.ResourceIDValue `tfsdk:"id"`
	Number    types.Int32                  `tfsdk:"number"`
	CreatedAt timetypes.RFC3339            `tfsdk:"created_at"`
}

// Framework interfaces
var (
	_ datasource.DataSource = &CredentialTypeVersionsDataSource{}
)

// New Object
func NewCredentialTypeVersionsDataSource() datasource.DataSource {
	return &CredentialTypeVersionsDataSource{}
}

// Metadata
func (r *CredentialTypeVersionsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_credential_type_versions"
}

var credentialTypeVersionsFilterRegexp = regexp.MustCompile(`^\s*version\s+eq\s+\d+\s*$`)

func (r *CredentialTypeVersionsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Datasource to retrieve a list of the versions of a PingOne Credentials credential type.  Credential type versions record the history of changes made to a credential type; a new version is created each time the credential type is updated, and the list is returned newest first.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment in which the credential type exists."),
			),

			"credential_type_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the credential type to retrieve versions for."),
			),

			"filter": schema.StringAttribute{
				Description: "An optional SCIM filter expression to apply to the credential type version selection, in the form `version eq <number>` (for example, `version eq 2`).  The `version` attribute is the only supported filter attribute, and the `eq` operator is the only supported operator.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(credentialTypeVersionsFilterRegexp, "must be a filter expression in the form `version eq <number>`, where `<number>` is an integer (for example, `version eq 2`)."),
				},
			},

			"versions": schema.ListNestedAttribute{
				Description: "The list of credential type versions that have been successfully retrieved, newest first.",
				Computed:    true,

				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Identifier (UUID) of the credential type version.",
							Computed:    true,

							CustomType: pingonetypes.ResourceIDType{},
						},

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
				},
			},
		},
	}
}

func (r *CredentialTypeVersionsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *CredentialTypeVersionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *CredentialTypeVersionsDataSourceModel

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
	var credentialTypeVersions []credentials.CredentialTypeVersion
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			apiRequest := r.Client.CredentialsAPIClient.CredentialTypesApi.
				ReadAllCredentialTypeVersions(ctx, data.EnvironmentId.ValueString(), data.CredentialTypeId.ValueString())

			if !data.Filter.IsNull() {
				apiRequest = apiRequest.Filter(data.Filter.ValueString())
			}

			pagedIterator := apiRequest.Execute()

			var initialHttpResponse *http.Response

			foundVersions := make([]credentials.CredentialTypeVersion, 0)

			for pageCursor, err := range pagedIterator {
				if err != nil {
					return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, pageCursor.HTTPResponse, err)
				}

				if initialHttpResponse == nil {
					initialHttpResponse = pageCursor.HTTPResponse
				}

				if pageCursor.EntityArray.Embedded != nil && pageCursor.EntityArray.Embedded.Versions != nil {
					foundVersions = append(foundVersions, pageCursor.EntityArray.Embedded.GetVersions()...)
				}
			}

			return foundVersions, initialHttpResponse, nil
		},
		"ReadAllCredentialTypeVersions",
		legacysdk.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		&credentialTypeVersions,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(data.EnvironmentId.ValueString(), credentialTypeVersions)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func toStateCredentialTypeVersionRecord(credentialTypeVersion *credentials.CredentialTypeVersion, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || credentialTypeVersion == nil {
		return types.ObjectNull(credentialTypeVersionRecordTFObjectTypes), diags
	}

	versionMap := map[string]attr.Value{
		"id":         framework.PingOneResourceIDOkToTF(credentialTypeVersion.GetIdOk()),
		"number":     framework.Int32OkToTF(credentialTypeVersion.GetVersionOk()),
		"created_at": framework.TimeOkToTF(credentialTypeVersion.GetCreatedAtOk()),
	}

	flattenedObj, d := types.ObjectValue(credentialTypeVersionRecordTFObjectTypes, versionMap)
	diags.Append(d...)

	return flattenedObj, diags
}

func (p *CredentialTypeVersionsDataSourceModel) toState(environmentID string, credentialTypeVersions []credentials.CredentialTypeVersion) diag.Diagnostics {
	var diags diag.Diagnostics

	if credentialTypeVersions == nil || environmentID == "" {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	var d diag.Diagnostics

	p.Id = framework.PingOneResourceIDToTF(environmentID)

	flattenedVersions := make([]attr.Value, 0, len(credentialTypeVersions))

	for _, credentialTypeVersion := range credentialTypeVersions {
		versionObj, d := toStateCredentialTypeVersionRecord(&credentialTypeVersion, true)
		diags.Append(d...)
		flattenedVersions = append(flattenedVersions, versionObj)
	}

	p.Versions, d = types.ListValue(types.ObjectType{AttrTypes: credentialTypeVersionRecordTFObjectTypes}, flattenedVersions)
	diags.Append(d...)

	return diags
}

var credentialTypeVersionRecordTFObjectTypes = map[string]attr.Type{
	"id":         pingonetypes.ResourceIDType{},
	"number":     types.Int32Type,
	"created_at": timetypes.RFC3339Type{},
}
