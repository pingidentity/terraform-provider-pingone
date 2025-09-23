// Copyright © 2025 Ping Identity Corporation

package verify

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/verify"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/legacysdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
)

// Types
type VoicePhraseDataSource serviceClientType

type voicePhraseDataSourceModel struct {
	Id            pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	VoicePhraseId pingonetypes.ResourceIDValue `tfsdk:"voice_phrase_id"`
	DisplayName   types.String                 `tfsdk:"display_name"`
	CreatedAt     timetypes.RFC3339            `tfsdk:"created_at"`
	UpdatedAt     timetypes.RFC3339            `tfsdk:"updated_at"`
}

// Framework interfaces
var (
	_ datasource.DataSource = &VoicePhraseDataSource{}
)

// New Object
func NewVoicePhraseDataSource() datasource.DataSource {
	return &VoicePhraseDataSource{}
}

// Metadata
func (r *VoicePhraseDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_verify_voice_phrase"
}

func (r *VoicePhraseDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	// schema descriptions and validation settings
	const attrMinLength = 1

	dataSourceExactlyOneOfRelativePaths := []string{
		"voice_phrase_id",
		"display_name",
	}

	voicePhraseIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Identifier (UUID) associated with the voice phrase.",
	).ExactlyOneOf(dataSourceExactlyOneOfRelativePaths)

	displayNameDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Name of the voice phrase container displayed in PingOne Admin UI or other administrative interface managing the container.",
	).ExactlyOneOf(dataSourceExactlyOneOfRelativePaths)

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Data source to find a PingOne Verify Voice Phrase by its Voice Phrase Id or Name.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("PingOne environment identifier (UUID) in which the verify voice phrase exists."),
			),

			"voice_phrase_id": schema.StringAttribute{
				Description:         voicePhraseIdDescription.Description,
				MarkdownDescription: voicePhraseIdDescription.MarkdownDescription,
				Optional:            true,

				CustomType: pingonetypes.ResourceIDType{},

				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName("display_name"),
					),
				},
			},

			"display_name": schema.StringAttribute{
				Description:         displayNameDescription.Description,
				MarkdownDescription: displayNameDescription.MarkdownDescription,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName("voice_phrase_id"),
					),
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			"created_at": schema.StringAttribute{
				Description: "Date and time the verify phrase was created.",
				Computed:    true,

				CustomType: timetypes.RFC3339Type{},
			},

			"updated_at": schema.StringAttribute{
				Description: "Date and time the verify phrase was updated. Can be null.",
				Computed:    true,

				CustomType: timetypes.RFC3339Type{},
			},
		},
	}
}

func (r *VoicePhraseDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *VoicePhraseDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *voicePhraseDataSourceModel

	if r.Client == nil || r.Client.VerifyAPIClient == nil {
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

	var voicePhrase *verify.VoicePhrase

	if !data.VoicePhraseId.IsNull() {

		// Run the API call
		resp.Diagnostics.Append(legacysdk.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				fO, fR, fErr := r.Client.VerifyAPIClient.VoicePhrasesApi.ReadOneVoicePhrase(ctx, data.EnvironmentId.ValueString(), data.VoicePhraseId.ValueString()).Execute()
				return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
			},
			"ReadOneVoicePhrase",
			legacysdk.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			&voicePhrase,
		)...)
		if resp.Diagnostics.HasError() {
			return
		}

	} else if !data.DisplayName.IsNull() {
		// Run the API call
		resp.Diagnostics.Append(legacysdk.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				pagedIterator := r.Client.VerifyAPIClient.VoicePhrasesApi.ReadAllVoicePhrases(ctx, data.EnvironmentId.ValueString()).Execute()

				var initialHttpResponse *http.Response

				for pageCursor, err := range pagedIterator {
					if err != nil {
						return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, pageCursor.HTTPResponse, err)
					}

					if initialHttpResponse == nil {
						initialHttpResponse = pageCursor.HTTPResponse
					}

					if voicePhrases, ok := pageCursor.EntityArray.Embedded.GetVoicePhrasesOk(); ok {

						for _, voicePhraseItem := range voicePhrases {

							if strings.EqualFold(voicePhraseItem.GetDisplayName(), data.DisplayName.ValueString()) {
								return &voicePhraseItem, pageCursor.HTTPResponse, nil
							}
						}
					}
				}

				return nil, initialHttpResponse, nil
			},
			"ReadAllVoicePhrases",
			legacysdk.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			&voicePhrase,
		)...)
		if resp.Diagnostics.HasError() {
			return
		}

		if voicePhrase == nil {
			resp.Diagnostics.AddError(
				"Cannot find voice phrase from display name",
				fmt.Sprintf("The voice phrase display name %s for environment %s cannot be found", data.DisplayName.String(), data.EnvironmentId.String()),
			)
			return
		}
	} else {
		resp.Diagnostics.AddError(
			"Missing parameter",
			"Cannot find the requested PingOne Voice Phrase: voice_phrase_id or display name argument must be set.",
		)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(voicePhrase)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (p *voicePhraseDataSourceModel) toState(apiObject *verify.VoicePhrase) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.PingOneResourceIDOkToTF(apiObject.GetIdOk())
	p.EnvironmentId = framework.PingOneResourceIDToTF(*apiObject.GetEnvironment().Id)
	p.VoicePhraseId = framework.PingOneResourceIDOkToTF(apiObject.GetIdOk())
	p.DisplayName = framework.StringOkToTF(apiObject.GetDisplayNameOk())
	p.CreatedAt = framework.TimeOkToTF(apiObject.GetCreatedAtOk())
	p.UpdatedAt = framework.TimeOkToTF(apiObject.GetUpdatedAtOk())

	return diags
}
