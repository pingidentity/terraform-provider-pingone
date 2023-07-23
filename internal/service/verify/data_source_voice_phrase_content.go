package verify

import (
	"context"
	"fmt"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/patrickcping/pingone-go-sdk-v2/verify"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	validation "github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type VoicePhraseContentDataSource struct {
	client *verify.APIClient
	region model.RegionMapping
}

type voicePhraseContentDataSourceModel struct {
	Id                   types.String `tfsdk:"id"`
	EnvironmentId        types.String `tfsdk:"environment_id"`
	VoicePhraseContentId types.String `tfsdk:"voice_phrase_content_id"`
	VoicePhraseId        types.String `tfsdk:"voice_phrase_id"`
	Locale               types.String `tfsdk:"locale"`
	Content              types.String `tfsdk:"content"`
	CreatedAt            types.String `tfsdk:"created_at"`
	UpdatedAt            types.String `tfsdk:"updated_at"`
}

// Framework interfaces
var (
	_ datasource.DataSource = &VoicePhraseContentDataSource{}
)

// New Object
func NewVoicePhraseContentDataSource() datasource.DataSource {
	return &VoicePhraseContentDataSource{}
}

// Metadata
func (r *VoicePhraseContentDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_voice_phrase_content"
}

func (r *VoicePhraseContentDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	// schema descriptions and validation settings

	// P1 Platform does not set a traditional UUID as the default phrase ID value
	const defaultVoicePhraseId = "exceptional_experiences"

	phraseIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"For a customer-defined phrase, the identifier (UUID) of the `voice_phrase` associated with the `voice_phrase_content` configuration. For pre-defined phrases, a string value.",
	)

	contentDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The phrase a user must speak as part of the voice enrollment or verification. The phrase must be written in the language and character set required by the language specified in the `locale` property.",
	)

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Data source to find PingOne Voice Phrase Contents by its Voice Phrase Content Id for a speicific PingOne Voice Phrase Id.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("PingOne environment identifier (UUID) in which the verify voice phrase exists."),
			),

			"voice_phrase_content_id": schema.StringAttribute{
				Description: "Identifier (UUID) associated with the voice phrase content.",
				Required:    true,
				Validators: []validator.String{
					validation.P1ResourceIDValidator(),
				},
			},

			"voice_phrase_id": schema.StringAttribute{
				Description:         phraseIdDescription.Description,
				MarkdownDescription: phraseIdDescription.MarkdownDescription,
				Required:            true,
				Validators: []validator.String{
					stringvalidator.Any(
						validation.P1ResourceIDValidator(),
						stringvalidator.RegexMatches(regexp.MustCompile(defaultVoicePhraseId), "Unexpected error with the pre-defined, default value. Please report this issue to the provider maintainers."),
					),
				},
			},

			"locale": schema.StringAttribute{
				Description: "Language localization requirement for the voice phrase contents.",
				Computed:    true,
			},

			"content": schema.StringAttribute{
				Description:         contentDescription.Description,
				MarkdownDescription: contentDescription.MarkdownDescription,
				Computed:            true,
			},

			"created_at": schema.StringAttribute{
				Description: "Date and time the verify phrase content was created.",
				Computed:    true,
			},

			"updated_at": schema.StringAttribute{
				Description: "Date and time the verify phrase content was updated. Can be null.",
				Computed:    true,
			},
		},
	}
}

func (r *VoicePhraseContentDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *VoicePhraseContentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *voicePhraseContentDataSourceModel

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

	// Run the API call
	var response *verify.VoicePhraseContents
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return r.client.VoicePhraseContentsApi.ReadOneVoicePhraseContent(ctx, data.EnvironmentId.ValueString(), data.VoicePhraseId.ValueString(), data.VoicePhraseContentId.ValueString()).Execute()
		},
		"ReadOneVoicePhraseContent",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		&response,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(response)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (p *voicePhraseContentDataSourceModel) toState(apiObject *verify.VoicePhraseContents) diag.Diagnostics {
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
	p.VoicePhraseContentId = framework.StringOkToTF(apiObject.GetIdOk())
	p.VoicePhraseId = framework.StringToTF(apiObject.GetVoicePhrase().Id)
	p.Locale = framework.StringOkToTF(apiObject.GetLocaleOk())
	p.Content = framework.StringOkToTF(apiObject.GetContentOk())
	p.CreatedAt = framework.TimeOkToTF(apiObject.GetCreatedAtOk())
	p.UpdatedAt = framework.TimeOkToTF(apiObject.GetUpdatedAtOk())

	return diags
}
