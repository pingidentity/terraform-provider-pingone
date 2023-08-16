package verify

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/patrickcping/pingone-go-sdk-v2/verify"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	validation "github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type VoicePhraseDataSource struct {
	client *verify.APIClient
	region model.RegionMapping
}

type voicePhraseDataSourceModel struct {
	Id            types.String `tfsdk:"id"`
	EnvironmentId types.String `tfsdk:"environment_id"`
	VoicePhraseId types.String `tfsdk:"voice_phrase_id"`
	DisplayName   types.String `tfsdk:"display_name"`
	CreatedAt     types.String `tfsdk:"created_at"`
	UpdatedAt     types.String `tfsdk:"updated_at"`
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
		Description: "Data source to find a PingOne Voice Phrase by its Voice Phrase Id or Name.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("PingOne environment identifier (UUID) in which the verify voice phrase exists."),
			),

			"voice_phrase_id": schema.StringAttribute{
				Description:         voicePhraseIdDescription.Description,
				MarkdownDescription: voicePhraseIdDescription.MarkdownDescription,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName("display_name"),
					),
					validation.P1ResourceIDValidator(),
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
			},

			"updated_at": schema.StringAttribute{
				Description: "Date and time the verify phrase was updated. Can be null.",
				Computed:    true,
			},
		},
	}
}

func (r *VoicePhraseDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = preparedClient
	r.region = resourceConfig.Client.API.Region
}

func (r *VoicePhraseDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *voicePhraseDataSourceModel

	if r.client == nil {
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

	var voicePhrase verify.VoicePhrase

	if !data.VoicePhraseId.IsNull() {

		// Run the API call
		var response *verify.VoicePhrase
		resp.Diagnostics.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				return r.client.VoicePhrasesApi.ReadOneVoicePhrase(ctx, data.EnvironmentId.ValueString(), data.VoicePhraseId.ValueString()).Execute()
			},
			"ReadOneVoicePhrase",
			framework.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			&response,
		)...)
		if resp.Diagnostics.HasError() {
			return
		}

		voicePhrase = *response

	} else if !data.DisplayName.IsNull() {
		// Run the API call
		var entityArray *verify.EntityArray
		resp.Diagnostics.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				return r.client.VoicePhrasesApi.ReadAllVoicePhrases(ctx, data.EnvironmentId.ValueString()).Execute()
			},
			"ReadAllVoicePhrases",
			framework.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			&entityArray,
		)...)
		if resp.Diagnostics.HasError() {
			return
		}

		if voicePhrases, ok := entityArray.Embedded.GetVoicePhrasesOk(); ok {

			found := false
			for _, voicePhraseItem := range voicePhrases {

				if voicePhraseItem.GetDisplayName() == data.DisplayName.ValueString() {
					voicePhrase = voicePhraseItem
					found = true
					break
				}
			}

			if !found {
				resp.Diagnostics.AddError(
					"Cannot find voice phrase from display name",
					fmt.Sprintf("The voice phrase display name %s for environment %s cannot be found", data.DisplayName.String(), data.EnvironmentId.String()),
				)
				return
			}
		}
	} else {
		resp.Diagnostics.AddError(
			"Missing parameter",
			"Cannot find the requested PingOne Voice Phrase: voice_phrase_id or display name argument must be set.",
		)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(&voicePhrase)...)
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

	p.Id = framework.StringOkToTF(apiObject.GetIdOk())
	p.EnvironmentId = framework.StringToTF(*apiObject.GetEnvironment().Id)
	p.VoicePhraseId = framework.StringOkToTF(apiObject.GetIdOk())
	p.DisplayName = framework.StringOkToTF(apiObject.GetDisplayNameOk())
	p.CreatedAt = framework.TimeOkToTF(apiObject.GetCreatedAtOk())
	p.UpdatedAt = framework.TimeOkToTF(apiObject.GetUpdatedAtOk())

	return diags
}
