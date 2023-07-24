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
type VoicePhraseContentsDataSource struct {
	client *verify.APIClient
	region model.RegionMapping
}

type voicePhraseContentsDataSourceModel struct {
	Id            types.String `tfsdk:"id"`
	EnvironmentId types.String `tfsdk:"environment_id"`
	VoicePhraseId types.String `tfsdk:"voice_phrase_id"`
	Ids           types.List   `tfsdk:"ids"`
}

// Framework interfaces
var (
	_ datasource.DataSource = &VoicePhraseContentsDataSource{}
)

// New Object
func NewVoicePhraseContentsDataSource() datasource.DataSource {
	return &VoicePhraseContentsDataSource{}
}

// Metadata
func (r *VoicePhraseContentsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_voice_phrase_contents"
}

func (r *VoicePhraseContentsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	// schema descriptions and validation settings

	// P1 Platform does not set a traditional UUID as the default phrase ID value
	const defaultVoicePhraseId = "exceptional_experiences"

	phraseIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"For a customer-defined phrase, the identifier (UUID) of the `voice_phrase` associated with the `voice_phrase_content` configuration. For pre-defined phrases, a string value.",
	)

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Data source to retrieve a list of PingOne Voice Phrase Contents. Filtering the list by SCIM or data filter currently is not supported.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("PingOne environment identifier (UUID) in which the verify voice phrase exists."),
			),

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

			"ids": framework.Attr_DataSourceReturnIDs(framework.SchemaAttributeDescription{
				Description: "The list of resulting IDs of voice phrase contents that have been successfully retrieved.",
			}),
		},
	}
}

func (r *VoicePhraseContentsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *VoicePhraseContentsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *voicePhraseContentsDataSourceModel

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
	var entityArray *verify.EntityArray
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return r.client.VoicePhraseContentsApi.ReadAllVoicePhraseContents(ctx, data.EnvironmentId.ValueString(), data.VoicePhraseId.ValueString()).Execute()
		},
		"ReadAllVoicePhraseContents",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		&entityArray,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(data.EnvironmentId.ValueString(), data.VoicePhraseId.ValueString(), entityArray.Embedded.GetContents())...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (p *voicePhraseContentsDataSourceModel) toState(environmentID string, voicePhraseID string, voicePhraseContents []verify.VoicePhraseContents) diag.Diagnostics {
	var diags diag.Diagnostics

	if voicePhraseContents == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	list := make([]string, 0)
	for _, item := range voicePhraseContents {
		list = append(list, item.GetId())
	}

	var d diag.Diagnostics

	p.Id = framework.StringToTF(environmentID)
	p.VoicePhraseId = framework.StringToTF(voicePhraseID)
	p.Ids, d = framework.StringSliceToTF(list)
	diags.Append(d...)

	return diags
}
