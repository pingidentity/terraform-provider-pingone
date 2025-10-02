// Copyright © 2025 Ping Identity Corporation

package base

import (
	"context"
	"fmt"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/legacysdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type AgreementLocalizationResource serviceClientType

type AgreementLocalizationResourceModel struct {
	Id                   pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId        pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	AgreementId          pingonetypes.ResourceIDValue `tfsdk:"agreement_id"`
	LanguageId           pingonetypes.ResourceIDValue `tfsdk:"language_id"`
	DisplayName          types.String                 `tfsdk:"display_name"`
	Locale               types.String                 `tfsdk:"locale"`
	Enabled              types.Bool                   `tfsdk:"enabled"`
	UXTextCheckboxAccept types.String                 `tfsdk:"text_checkbox_accept"`
	UXTextButtonContinue types.String                 `tfsdk:"text_button_continue"`
	UXTextButtonDecline  types.String                 `tfsdk:"text_button_decline"`
}

// Framework interfaces
var (
	_ resource.Resource                = &AgreementLocalizationResource{}
	_ resource.ResourceWithConfigure   = &AgreementLocalizationResource{}
	_ resource.ResourceWithImportState = &AgreementLocalizationResource{}
)

// New Object
func NewAgreementLocalizationResource() resource.Resource {
	return &AgreementLocalizationResource{}
}

// Metadata
func (r *AgreementLocalizationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_agreement_localization"
}

// Schema.
func (r *AgreementLocalizationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	const attrMinLength = 1

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage agreement localizations in a PingOne environment.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to associate the agreement localization with."),
			),

			"agreement_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the agreement to associate the agreement localization with."),
			),

			"language_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the language in the PingOne environment that the localization applies to."),
			),

			"display_name": schema.StringAttribute{
				Description: "A string used as the title of the agreement for the language presented to the user.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			"locale": schema.StringAttribute{
				Description: "A string used as the locale code of the agreement localization to retrieve. Either `agreement_localization_id`, `display_name` or `locale` can be used to retrieve the agreement localization, but cannot be set together.",
				Computed:    true,
			},

			"enabled": schema.BoolAttribute{
				Description: "A boolean that specifies whether the localization (and it's revision text) is enabled in the agreement.",
				Computed:    true,
			},

			"text_checkbox_accept": schema.StringAttribute{
				Description: "A string that specifies the text next to the \"accept\" checkbox in the end user interface. Accepted character are unicode letters, combining marks, numeric characters, whitespace, and punctuation characters.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(`^[\p{L}\p{M}\p{N}\p{Zs}\p{P}]+$`), "Accepted character are unicode letters, combining marks, numeric characters, whitespace, and punctuation characters."),
				},
			},

			"text_button_continue": schema.StringAttribute{
				Description: "A string that specifies the text next to the \"continue\" button in the end user interface. Accepted character are unicode letters, combining marks, numeric characters, whitespace, and punctuation characters.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(`^[\p{L}\p{M}\p{N}\p{Zs}\p{P}]+$`), "Accepted character are unicode letters, combining marks, numeric characters, whitespace, and punctuation characters."),
				},
			},

			"text_button_decline": schema.StringAttribute{
				Description: "A string that specifies the text next to the \"decline\" button in the end user interface. Accepted character are unicode letters, combining marks, numeric characters, whitespace, and punctuation characters.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(`^[\p{L}\p{M}\p{N}\p{Zs}\p{P}]+$`), "Accepted character are unicode letters, combining marks, numeric characters, whitespace, and punctuation characters."),
				},
			},
		},
	}
}

func (r *AgreementLocalizationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *AgreementLocalizationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state AgreementLocalizationResourceModel

	if r.Client == nil || r.Client.ManagementAPIClient == nil {
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

	var language *management.Language
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.LanguagesApi.ReadOneLanguage(ctx, plan.EnvironmentId.ValueString(), plan.LanguageId.ValueString()).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadOneLanguage",
		legacysdk.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		&language,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if v, ok := language.GetEnabledOk(); !ok || !*v {
		resp.Diagnostics.AddError(
			"Invalid language parameter",
			fmt.Sprintf("The language with ID %s needs to be enabled in the environment before it can be assigned to a localized agreement", plan.LanguageId.ValueString()),
		)
		return
	}

	locale := language.GetLocale()

	// Build the model for the API
	localization := plan.expand(locale)

	// Run the API call
	var response *management.AgreementLanguage
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.AgreementLanguagesResourcesApi.CreateAgreementLanguage(ctx, plan.EnvironmentId.ValueString(), plan.AgreementId.ValueString()).AgreementLanguage(*localization).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"CreateAgreementLanguage",
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
	resp.Diagnostics.Append(state.toState(response, language)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *AgreementLocalizationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *AgreementLocalizationResourceModel

	if r.Client == nil || r.Client.ManagementAPIClient == nil {
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
	var response *management.AgreementLanguage
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.AgreementLanguagesResourcesApi.ReadOneAgreementLanguage(ctx, data.EnvironmentId.ValueString(), data.AgreementId.ValueString(), data.Id.ValueString()).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadOneAgreementLanguage",
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

	var language *management.Language
	if data.LanguageId.IsNull() || data.LanguageId.IsUnknown() {
		var d diag.Diagnostics
		language, d = findLanguageByLocale_Framework(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), response.GetLocale())
		resp.Diagnostics.Append(d...)
	} else {
		resp.Diagnostics.Append(legacysdk.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				fO, fR, fErr := r.Client.ManagementAPIClient.LanguagesApi.ReadOneLanguage(ctx, data.EnvironmentId.ValueString(), data.LanguageId.ValueString()).Execute()
				return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
			},
			"ReadOneLanguage",
			legacysdk.CustomErrorResourceNotFoundWarning,
			sdk.DefaultCreateReadRetryable,
			&language,
		)...)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(response, language)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AgreementLocalizationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state AgreementLocalizationResourceModel

	if r.Client == nil || r.Client.ManagementAPIClient == nil {
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

	var language *management.Language
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.LanguagesApi.ReadOneLanguage(ctx, plan.EnvironmentId.ValueString(), plan.LanguageId.ValueString()).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadOneLanguage",
		legacysdk.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		&language,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if v, ok := language.GetEnabledOk(); !ok || !*v {
		resp.Diagnostics.AddError(
			"Invalid language parameter",
			fmt.Sprintf("The language with ID %s needs to be enabled in the environment before it can be assigned to a localized agreement", plan.LanguageId.ValueString()),
		)
		return
	}

	locale := language.GetLocale()

	// Build the model for the API
	localization := plan.expand(locale)

	// Run the API call
	var response *management.AgreementLanguage
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.AgreementLanguagesResourcesApi.UpdateAgreementLanguage(ctx, plan.EnvironmentId.ValueString(), plan.AgreementId.ValueString(), plan.Id.ValueString()).AgreementLanguage(*localization).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateAgreementLanguage",
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
	resp.Diagnostics.Append(state.toState(response, language)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *AgreementLocalizationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *AgreementLocalizationResourceModel

	if r.Client == nil || r.Client.ManagementAPIClient == nil {
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
			fR, fErr := r.Client.ManagementAPIClient.AgreementLanguagesResourcesApi.DeleteAgreementLanguage(ctx, data.EnvironmentId.ValueString(), data.AgreementId.ValueString(), data.Id.ValueString()).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, fR, fErr)
		},
		"DeleteAgreementLanguage",
		agreementLocalizationDeleteErrorHandler,
		sdk.DefaultCreateReadRetryable,
		nil,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *AgreementLocalizationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:  "agreement_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:     "agreement_localization_id",
			Regexp:    verify.P1ResourceIDRegexp,
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

func (p *AgreementLocalizationResourceModel) expand(locale string) *management.AgreementLanguage {
	data := management.NewAgreementLanguage(
		p.DisplayName.ValueString(),
		p.Enabled.ValueBool(),
		locale,
	)

	userExperience := management.NewAgreementLanguageUserExperience()
	setUx := false

	if !p.UXTextCheckboxAccept.IsNull() && !p.UXTextCheckboxAccept.IsUnknown() {
		userExperience.SetAcceptCheckboxText(p.UXTextCheckboxAccept.ValueString())
		if !setUx {
			setUx = true
		}
	}

	if !p.UXTextButtonContinue.IsNull() && !p.UXTextButtonContinue.IsUnknown() {
		userExperience.SetContinueButtonText(p.UXTextButtonContinue.ValueString())
		if !setUx {
			setUx = true
		}
	}

	if !p.UXTextButtonDecline.IsNull() && !p.UXTextButtonDecline.IsUnknown() {
		userExperience.SetDeclineButtonText(p.UXTextButtonDecline.ValueString())
		if !setUx {
			setUx = true
		}
	}

	if setUx {
		data.SetUserExperience(*userExperience)
	}

	return data
}

func (p *AgreementLocalizationResourceModel) toState(apiObject *management.AgreementLanguage, languageApiObject *management.Language) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.PingOneResourceIDToTF(apiObject.GetId())
	p.AgreementId = framework.PingOneResourceIDToTF(*apiObject.GetAgreement().Id)
	p.LanguageId = framework.PingOneResourceIDOkToTF(languageApiObject.GetIdOk())
	p.DisplayName = framework.StringOkToTF(apiObject.GetDisplayNameOk())
	p.Locale = framework.StringOkToTF(apiObject.GetLocaleOk())
	p.Enabled = framework.BoolOkToTF(apiObject.GetEnabledOk())

	if v, ok := apiObject.GetUserExperienceOk(); ok {
		p.UXTextCheckboxAccept = framework.StringOkToTF(v.GetAcceptCheckboxTextOk())
		p.UXTextButtonContinue = framework.StringOkToTF(v.GetContinueButtonTextOk())
		p.UXTextButtonDecline = framework.StringOkToTF(v.GetDeclineButtonTextOk())
	}

	return diags
}

func agreementLocalizationDeleteErrorHandler(r *http.Response, p1Error *model.P1Error) diag.Diagnostics {
	var diags diag.Diagnostics

	diags.Append(legacysdk.CustomErrorResourceNotFoundWarning(r, p1Error)...)

	if v, ok := p1Error.GetDetailsOk(); ok && v != nil && len(v) > 0 {
		if v[0].GetCode() == "CONSTRAINT_VIOLATION" {
			if match, _ := regexp.MatchString("Agreement language with effective revision can not be deleted.", v[0].GetMessage()); match {
				diags.AddWarning(
					"Cannot delete the agreement localization, a localization with a currently effective revision cannot be deleted.",
					"The agreement localization is left in place but no longer managed by the provider.",
				)

				return diags
			}
		}
	}

	return diags
}
