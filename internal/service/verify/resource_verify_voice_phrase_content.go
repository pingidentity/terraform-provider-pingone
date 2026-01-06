// Copyright © 2025 Ping Identity Corporation

package verify

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/verify"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/legacysdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	validation "github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type VoicePhraseContentResource serviceClientType

type voicePhraseContentResourceModel struct {
	Id            pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	VoicePhraseId pingonetypes.ResourceIDValue `tfsdk:"voice_phrase_id"`
	Locale        types.String                 `tfsdk:"locale"`
	Content       types.String                 `tfsdk:"content"`
	CreatedAt     timetypes.RFC3339            `tfsdk:"created_at"`
	UpdatedAt     timetypes.RFC3339            `tfsdk:"updated_at"`
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
	resp.TypeName = req.ProviderTypeName + "_verify_voice_phrase_content"
}

func (r *VoicePhraseContentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	// schema descriptions and validation settings
	const attrMinLength = 1

	phraseIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The identifier (UUID) of the `voice_phrase` associated with the `voice_phrase_content` configuration.",
	).RequiresReplace()

	localeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Language localization requirement for the voice phrase contents.",
	).AllowedValuesEnum(validation.FullIsoList())

	contentDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The phrase a user must speak as part of the voice enrollment or verification. The phrase must be written in the language and character set required by the language specified in the `locale` property.",
	)

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to configure the phrases to speak during voice verification enrollment or validation.",

		DeprecationMessage: "Deprecation notice: This resource is deprecated and will be removed in a future release. Please use alternative verification methods.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("PingOne environment identifier (UUID) in which the verify voice phrase exists."),
			),

			"voice_phrase_id": schema.StringAttribute{
				Description:         phraseIdDescription.Description,
				MarkdownDescription: phraseIdDescription.MarkdownDescription,
				Required:            true,

				CustomType: pingonetypes.ResourceIDType{},

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			"locale": schema.StringAttribute{
				Description:         localeDescription.Description,
				MarkdownDescription: localeDescription.MarkdownDescription,
				Required:            true,
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

				CustomType: timetypes.RFC3339Type{},

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"updated_at": schema.StringAttribute{
				Description: "Date and time the verify phrase content was updated. Can be null.",
				Computed:    true,

				CustomType: timetypes.RFC3339Type{},
			},
		},
	}
}

func (r *VoicePhraseContentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *VoicePhraseContentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state voicePhraseContentResourceModel

	if r.Client == nil || r.Client.VerifyAPIClient == nil {
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
	voicePhraseContent, d := plan.expand()
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *verify.VoicePhraseContents
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.VerifyAPIClient.VoicePhraseContentsApi.CreateVoicePhraseContent(ctx, plan.EnvironmentId.ValueString(), plan.VoicePhraseId.ValueString()).VoicePhraseContents(*voicePhraseContent).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"CreateVoicePhraseContent",
		legacysdk.DefaultCustomError,
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

	if r.Client == nil || r.Client.VerifyAPIClient == nil {
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
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.VerifyAPIClient.VoicePhraseContentsApi.ReadOneVoicePhraseContent(ctx, data.EnvironmentId.ValueString(), data.VoicePhraseId.ValueString(), data.Id.ValueString()).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadOneVoicePhraseContent",
		legacysdk.CustomErrorResourceNotFoundWarning,
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

	if r.Client == nil || r.Client.VerifyAPIClient == nil {
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
	voicePhraseContent, d := plan.expand()
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *verify.VoicePhraseContents
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.VerifyAPIClient.VoicePhraseContentsApi.UpdateVoicePhraseContent(ctx, plan.EnvironmentId.ValueString(), plan.VoicePhraseId.ValueString(), plan.Id.ValueString()).VoicePhraseContents(*voicePhraseContent).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateVoicePhraseContent",
		legacysdk.DefaultCustomError,
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

	if r.Client == nil || r.Client.VerifyAPIClient == nil {
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
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fR, fErr := r.Client.VerifyAPIClient.VoicePhraseContentsApi.DeleteVoicePhraseContent(ctx, data.EnvironmentId.ValueString(), data.VoicePhraseId.ValueString(), data.Id.ValueString()).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, fR, fErr)
		},
		"DeleteVoicePhrase",
		legacysdk.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
		nil,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *VoicePhraseContentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
			Regexp: validation.P1ResourceIDRegexp,
		},
		{
			Label:  "voice_phrase_id",
			Regexp: validation.P1ResourceIDRegexp,
		},
		{
			Label:     "voice_phrase_content_id",
			Regexp:    validation.P1ResourceIDRegexp,
			PrimaryID: true,
		},
	}

	attributes, err := framework.ParseImportID(req.ID, idComponents...)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			err.Error(),
		)
		return
	}

	for _, idComponent := range idComponents {
		pathKey := idComponent.Label

		if idComponent.PrimaryID {
			pathKey = "id"
		}

		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(pathKey), attributes[idComponent.Label])...)
	}
}

func (p *voicePhraseContentResourceModel) expand() (*verify.VoicePhraseContents, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := verify.NewVoicePhraseContentsWithDefaults()
	if data == nil {
		diags.AddError(
			"Unexpected Error on Expand",
			"Cannot expand new voice phrase content with defaults.  Please report this to the provider maintainers.",
		)
	}

	if !p.VoicePhraseId.IsNull() && !p.VoicePhraseId.IsUnknown() {
		data.SetVoicePhrase(*verify.NewVoicePhraseContentsVoicePhrase(p.VoicePhraseId.ValueString()))
	}

	if !p.Content.IsNull() && !p.Content.IsUnknown() {
		data.SetContent(p.Content.ValueString())
	}

	if !p.Locale.IsNull() && !p.Locale.IsUnknown() {
		data.SetLocale(p.Locale.ValueString())
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

	p.Id = framework.PingOneResourceIDOkToTF(apiObject.GetIdOk())
	p.EnvironmentId = framework.PingOneResourceIDToTF(*apiObject.GetEnvironment().Id)
	p.VoicePhraseId = framework.PingOneResourceIDToTF(apiObject.GetVoicePhrase().Id)
	p.Locale = framework.StringOkToTF(apiObject.GetLocaleOk())
	p.Content = framework.StringOkToTF(apiObject.GetContentOk())
	p.CreatedAt = framework.TimeOkToTF(apiObject.GetCreatedAtOk())
	p.UpdatedAt = framework.TimeOkToTF(apiObject.GetUpdatedAtOk())

	return diags
}
