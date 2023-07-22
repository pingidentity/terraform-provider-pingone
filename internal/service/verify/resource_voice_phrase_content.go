package verify

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/patrickcping/pingone-go-sdk-v2/verify"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	validation "github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type VoicePhraseContentResource struct {
	client *verify.APIClient
	region model.RegionMapping
}

type voicePhraseContentResourceModel struct {
	Id            types.String `tfsdk:"id"`
	EnvironmentId types.String `tfsdk:"environment_id"`
	VoicePhraseId types.String `tfsdk:"voice_phrase_id"`
	Locale        types.String `tfsdk:"locale"`
	Content       types.String `tfsdk:"content"`
	CreatedAt     types.String `tfsdk:"created_at"`
	UpdatedAt     types.String `tfsdk:"updated_at"`
}

// Framework interfaces
var (
	_ resource.Resource                = &VoicePhraseContentResource{}
	_ resource.ResourceWithConfigure   = &VoicePhraseContentResource{}
	_ resource.ResourceWithImportState = &VoicePhraseContentResource{}
)

// New Object
func NewVoicePhraseContentResource() resource.Resource {
	return &VoicePhraseContentResource{}
}

// Metadata
func (r *VoicePhraseContentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_voice_phrase_content"
}

func (r *VoicePhraseContentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	// schema descriptions and validation settings
	const attrMinLength = 1

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
		Description: "Resource to configure the voice enrollment or verification requirements when configuring a `verify_policy` for voice verification.\n\n" +
			"A `voice_phrase_id` is obtained by configuring the `voice_phrase` container with a name. The actual phrases to speak are defined in the `voice_phrase_contents` configuration, where the content has a locale and the phrase to speak, written in the language required by the locale.",

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

			"locale": schema.StringAttribute{
				Description: "Language localization requirement for the voice phrase contents.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf(validation.FullIsoList()...),
				},
			},

			"content": schema.StringAttribute{
				Description:         contentDescription.Description,
				MarkdownDescription: contentDescription.MarkdownDescription,
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
				},
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

func (r *VoicePhraseContentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *VoicePhraseContentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state voicePhraseContentResourceModel

	if r.client == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the model for the API
	VoicePhraseContent, d := plan.expand()
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *verify.VoicePhraseContents
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return r.client.VoicePhraseContentsApi.CreateVoicePhraseContent(ctx, plan.EnvironmentId.ValueString(), plan.VoicePhraseId.ValueString()).VoicePhraseContents(*VoicePhraseContent).Execute()
		},
		"CreateVoicePhraseContent",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		&response,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(response)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *VoicePhraseContentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *voicePhraseContentResourceModel

	if r.client == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *verify.VoicePhraseContents
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return r.client.VoicePhraseContentsApi.ReadOneVoicePhraseContent(ctx, data.EnvironmentId.ValueString(), data.VoicePhraseId.ValueString(), data.Id.ValueString()).Execute()
		},
		"ReadOneVoicePhraseContent",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
		&response,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Remove from state if resource is not found
	if response == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(response)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VoicePhraseContentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state voicePhraseContentResourceModel

	if r.client == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the model for the API
	VoicePhraseContent, d := plan.expand()
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *verify.VoicePhraseContents
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return r.client.VoicePhraseContentsApi.UpdateVoicePhraseContent(ctx, plan.EnvironmentId.ValueString(), plan.VoicePhraseId.ValueString(), plan.Id.ValueString()).VoicePhraseContents(*VoicePhraseContent).Execute()
		},
		"UpdateVoicePhraseContent",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		&response,
	)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(response)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *VoicePhraseContentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *voicePhraseContentResourceModel

	if r.client == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			r, err := r.client.VoicePhraseContentsApi.DeleteVoicePhraseContent(ctx, data.EnvironmentId.ValueString(), data.VoicePhraseId.ValueString(), data.Id.ValueString()).Execute()
			return nil, r, err
		},
		"DeleteVoicePhrase",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
		nil,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *VoicePhraseContentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splitLength := 3
	attributes := strings.SplitN(req.ID, "/", splitLength)

	if len(attributes) != splitLength {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("invalid id (\"%s\") specified, should be in format \"environment_id/voice_phrase_id/voice_phrase_content_id\"", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_id"), attributes[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("voice_phrase_id"), attributes[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), attributes[1])...)
}

func (p *voicePhraseContentResourceModel) expand() (*verify.VoicePhraseContents, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := verify.NewVoicePhraseContentsWithDefaults()

	data.SetId(p.Id.ValueString())

	if !p.VoicePhraseId.IsNull() && !p.VoicePhraseId.IsUnknown() {
		data.SetVoicePhrase(*verify.NewVoicePhraseContentsVoicePhrase(p.VoicePhraseId.ValueString()))
	}

	if !p.Content.IsNull() && !p.Content.IsUnknown() {
		data.SetContent(p.Content.ValueString())
	}

	if !p.Locale.IsNull() && !p.Locale.IsUnknown() {
		data.SetLocale(p.Locale.ValueString())
	}

	if !p.CreatedAt.IsNull() && !p.CreatedAt.IsUnknown() {
		createdAt, err := time.Parse(time.RFC3339, p.CreatedAt.ValueString())
		if err != nil {
			diags.AddError(
				"Unexpected Value",
				fmt.Sprintf("Unexpected createdAt value: %s. Please report this to the provider maintainers.", err.Error()),
			)
		}
		data.SetCreatedAt(createdAt)
	}

	if !p.UpdatedAt.IsNull() && !p.UpdatedAt.IsUnknown() {
		updatedAt, err := time.Parse(time.RFC3339, p.UpdatedAt.ValueString())
		if err != nil {
			diags.AddError(
				"Unexpected Value",
				fmt.Sprintf("Unexpected updatedAt value: %s. Please report this to the provider maintainers.", err.Error()),
			)
		}
		data.SetUpdatedAt(updatedAt)

		if data == nil {
			diags.AddError(
				"Unexpected Value",
				"Verify Policy object was unexpectedly null on expansion. Please report this to the provider maintainers.",
			)
		}
	}

	return data, diags
}

func (p *voicePhraseContentResourceModel) toState(apiObject *verify.VoicePhraseContents) diag.Diagnostics {
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
	p.VoicePhraseId = framework.StringToTF(apiObject.GetVoicePhrase().Id)
	p.Locale = framework.StringOkToTF(apiObject.GetLocaleOk())
	p.Content = framework.StringOkToTF(apiObject.GetContentOk())
	p.CreatedAt = framework.TimeOkToTF(apiObject.GetCreatedAtOk())
	p.UpdatedAt = framework.TimeOkToTF(apiObject.GetUpdatedAtOk())

	return diags
}
