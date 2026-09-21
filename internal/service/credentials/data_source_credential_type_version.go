// Copyright © 2026 Ping Identity Corporation

package credentials

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
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
	Snapshot                types.Object                 `tfsdk:"snapshot"`
}

type CredentialTypeVersionSnapshotDataSourceModel struct {
	Title              types.String      `tfsdk:"title"`
	Description        types.String      `tfsdk:"description"`
	CardType           types.String      `tfsdk:"card_type"`
	CardDesignTemplate types.String      `tfsdk:"card_design_template"`
	ManagementMode     types.String      `tfsdk:"management_mode"`
	Metadata           types.Object      `tfsdk:"metadata"`
	RevokeOnDelete     types.Bool        `tfsdk:"revoke_on_delete"`
	CreatedAt          timetypes.RFC3339 `tfsdk:"created_at"`
	UpdatedAt          timetypes.RFC3339 `tfsdk:"updated_at"`
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

			"snapshot": schema.SingleNestedAttribute{
				Description: "The full credential type content as it was at this version.",
				Computed:    true,

				Attributes: map[string]schema.Attribute{
					"title": schema.StringAttribute{
						Description: "Title of the credential.",
						Computed:    true,
					},

					"description": schema.StringAttribute{
						Description: "A description of the credential type.",
						Computed:    true,
					},

					"card_type": schema.StringAttribute{
						Description: "A descriptor of the credential type. Can be non-identity types such as proof of employment or proof of insurance.",
						Computed:    true,
					},

					"card_design_template": schema.StringAttribute{
						Description: "An SVG formatted image containing placeholders for the credentials fields that need to be displayed in the image.",
						Computed:    true,
					},

					"management_mode": schema.StringAttribute{
						Description: "Specifies the management mode of the credential type.",
						Computed:    true,
					},

					"revoke_on_delete": schema.BoolAttribute{
						Description: "Specifies whether a user's issued verifiable credentials are automatically revoked when the credential type is deleted.",
						Computed:    true,
					},

					"created_at": schema.StringAttribute{
						Description: "Date and time the credential type version was created.",
						Computed:    true,

						CustomType: timetypes.RFC3339Type{},
					},

					"updated_at": schema.StringAttribute{
						Description: "Date and time the object was updated. Can be null.",
						Computed:    true,

						CustomType: timetypes.RFC3339Type{},
					},

					"metadata": schema.SingleNestedAttribute{
						Description: "Contains the names, data types, and other metadata related to the credential.",
						Computed:    true,

						Attributes: map[string]schema.Attribute{
							"background_image": schema.StringAttribute{
								Description: "URL or fully qualified path to the image file used for the credential background.",
								Computed:    true,
							},

							"bg_opacity_percent": schema.Int32Attribute{
								Description: "Percent opacity of the background image in the credential.",
								Computed:    true,
							},

							"card_color": schema.StringAttribute{
								Description: "Color to show on the credential.",
								Computed:    true,
							},

							"columns": schema.Int32Attribute{
								Description: "Number of columns to organize the fields displayed on the credential.",
								Computed:    true,
							},

							"description": schema.StringAttribute{
								Description: "Description of the credential.",
								Computed:    true,
							},

							"logo_image": schema.StringAttribute{
								Description: "URL or fully qualified path to the image file used for the credential logo.",
								Computed:    true,
							},

							"name": schema.StringAttribute{
								Description: "Name of the credential.",
								Computed:    true,
							},

							"text_color": schema.StringAttribute{
								Description: "Color of the text to show on the credential.",
								Computed:    true,
							},

							"version": schema.Int32Attribute{
								Description: "Version of this credential metadata.",
								Computed:    true,
							},

							"fields": schema.ListNestedAttribute{
								Description: "Array of objects representing the credential fields.",
								Computed:    true,

								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"id": schema.StringAttribute{
											Description: "Identifier of the field object.",
											Computed:    true,

											CustomType: pingonetypes.ResourceIDType{},
										},

										"type": schema.StringAttribute{
											Description: "Type of data in the field.",
											Computed:    true,
										},

										"title": schema.StringAttribute{
											Description: "Descriptive text when showing the field.",
											Computed:    true,
										},

										"file_support": schema.StringAttribute{
											Description: "Specifies how an image is stored in the credential field.",
											Computed:    true,
										},

										"is_visible": schema.BoolAttribute{
											Description: "Specifies whether the field should be visible to viewers of the credential.",
											Computed:    true,
										},

										"attribute": schema.StringAttribute{
											Description: "Name of the PingOne Directory attribute. Present if field.type is Directory Attribute.",
											Computed:    true,
										},

										"value": schema.StringAttribute{
											Description: "The text to appear on the credential for a field.type of Alphanumeric Text.",
											Computed:    true,
										},

										"required": schema.BoolAttribute{
											Description: "Specifies whether the field is required for the credential.",
											Computed:    true,
										},
									},
								},
							},
						},
					},
				},
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

	snapshot, d := toStateCredentialTypeVersionSnapshot(apiObject.GetSnapshotOk())
	diags.Append(d...)
	p.Snapshot = snapshot

	return diags
}

func toStateCredentialTypeVersionSnapshot(snapshot *credentials.CredentialType, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || snapshot == nil {
		return types.ObjectNull(snapshotTFObjectTypes), diags
	}

	// core snapshot attributes
	snapshotMap := map[string]attr.Value{
		"title":                framework.StringOkToTF(snapshot.GetTitleOk()),
		"description":          framework.StringOkToTF(snapshot.GetDescriptionOk()),
		"card_type":            framework.StringOkToTF(snapshot.GetCardTypeOk()),
		"card_design_template": framework.StringOkToTF(snapshot.GetCardDesignTemplateOk()),
		"revoke_on_delete":     types.BoolNull(),
		"created_at":           framework.TimeOkToTF(snapshot.GetCreatedAtOk()),
		"updated_at":           framework.TimeOkToTF(snapshot.GetUpdatedAtOk()),
	}

	if v, ok := snapshot.GetManagementOk(); ok {
		snapshotMap["management_mode"] = framework.EnumOkToTF(v.GetModeOk())
	} else {
		snapshotMap["management_mode"] = types.StringNull()
	}

	if v, ok := snapshot.GetOnDeleteOk(); ok {
		snapshotMap["revoke_on_delete"] = framework.BoolOkToTF(v.GetRevokeIssuedCredentialsOk())
	}

	// snapshot metadata object
	metadata, d := toStateMetadataDataSource(snapshot.GetMetadataOk())
	diags.Append(d...)

	snapshotMap["metadata"] = metadata

	flattenedObj, d := types.ObjectValue(snapshotTFObjectTypes, snapshotMap)
	diags.Append(d...)

	return flattenedObj, diags
}

var snapshotTFObjectTypes = map[string]attr.Type{
	"title":                types.StringType,
	"description":          types.StringType,
	"card_type":            types.StringType,
	"card_design_template": types.StringType,
	"management_mode":      types.StringType,
	"metadata":             types.ObjectType{AttrTypes: metadataDataSourceServiceTFObjectTypes},
	"revoke_on_delete":     types.BoolType,
	"created_at":           timetypes.RFC3339Type{},
	"updated_at":           timetypes.RFC3339Type{},
}
