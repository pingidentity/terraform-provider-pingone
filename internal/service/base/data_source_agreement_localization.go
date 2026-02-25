// Copyright © 2026 Ping Identity Corporation

package base

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
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/legacysdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type AgreementLocalizationDataSource serviceClientType

type AgreementLocalizationDataSourceModel struct {
	Id                      pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId           pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	AgreementId             pingonetypes.ResourceIDValue `tfsdk:"agreement_id"`
	AgreementLocalizationId pingonetypes.ResourceIDValue `tfsdk:"agreement_localization_id"`
	LanguageId              pingonetypes.ResourceIDValue `tfsdk:"language_id"`
	DisplayName             types.String                 `tfsdk:"display_name"`
	Locale                  types.String                 `tfsdk:"locale"`
	Enabled                 types.Bool                   `tfsdk:"enabled"`
	UXTextCheckboxAccept    types.String                 `tfsdk:"text_checkbox_accept"`
	UXTextButtonContinue    types.String                 `tfsdk:"text_button_continue"`
	UXTextButtonDecline     types.String                 `tfsdk:"text_button_decline"`
	CurrentRevisionId       pingonetypes.ResourceIDValue `tfsdk:"current_revision_id"`
}

// Framework interfaces
var (
	_ datasource.DataSource = &AgreementLocalizationDataSource{}
)

// New Object
func NewAgreementLocalizationDataSource() datasource.DataSource {
	return &AgreementLocalizationDataSource{}
}

// Metadata
func (r *AgreementLocalizationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_agreement_localization"
}

// Schema
func (r *AgreementLocalizationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	nameLength := 1

	agreementLocalizationIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The ID of the agreement localization language to retrieve.",
	).ExactlyOneOf([]string{"agreement_localization_id", "display_name", "locale"})

	displayNameDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string used as the title of the agreement localization to retrieve.",
	).ExactlyOneOf([]string{"agreement_localization_id", "display_name", "locale"})

	localeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string used as the locale code of the agreement localization to retrieve.",
	).AllowedValuesEnum(verify.FullIsoList()).ExactlyOneOf([]string{"agreement_localization_id", "display_name", "locale"})

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Datasource to retrieve details of an agreement localization in a PingOne environment.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment that is configured with the agreement localization."),
			),

			"agreement_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The UUID that identifies the agreement ID that the localization is applied to."),
			),

			"agreement_localization_id": schema.StringAttribute{
				Description:         agreementLocalizationIdDescription.Description,
				MarkdownDescription: agreementLocalizationIdDescription.MarkdownDescription,
				Optional:            true,

				CustomType: pingonetypes.ResourceIDType{},

				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName("display_name"),
						path.MatchRelative().AtParent().AtName("locale"),
					),
				},
			},

			"display_name": schema.StringAttribute{
				Description:         displayNameDescription.Description,
				MarkdownDescription: displayNameDescription.MarkdownDescription,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName("agreement_localization_id"),
						path.MatchRelative().AtParent().AtName("locale"),
					),
					stringvalidator.LengthAtLeast(nameLength),
				},
			},

			"locale": schema.StringAttribute{
				Description:         localeDescription.Description,
				MarkdownDescription: localeDescription.MarkdownDescription,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName("agreement_localization_id"),
						path.MatchRelative().AtParent().AtName("display_name"),
					),
					stringvalidator.OneOf(verify.FullIsoList()...),
				},
			},

			"language_id": schema.StringAttribute{
				Description: "The ID of the language used for the agreement localization.",
				Computed:    true,

				CustomType: pingonetypes.ResourceIDType{},
			},

			"enabled": schema.BoolAttribute{
				Description: "A boolean that specifies whether the localization (and it's revision text) is enabled in the agreement.",
				Computed:    true,
			},

			"text_checkbox_accept": schema.StringAttribute{
				Description: "A string that specifies the text next to the \"accept\" checkbox in the end user interface.",
				Computed:    true,
			},

			"text_button_continue": schema.StringAttribute{
				Description: "A string that specifies the text next to the \"continue\" button in the end user interface.",
				Computed:    true,
			},

			"text_button_decline": schema.StringAttribute{
				Description: "A string that specifies the text next to the \"decline\" button in the end user interface.",
				Computed:    true,
			},

			"current_revision_id": schema.StringAttribute{
				Description: "A string that specifies the UUID of the current revision associated with this agreement localization resource. The current revision is the one shown to users for new consents in the language.",
				Computed:    true,

				CustomType: pingonetypes.ResourceIDType{},
			},
		},
	}
}

func (r *AgreementLocalizationDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *AgreementLocalizationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *AgreementLocalizationDataSourceModel

	if r.Client == nil || r.Client.ManagementAPIClient == nil {
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

	var agreementLocalization *management.AgreementLanguage

	if !data.DisplayName.IsNull() || !data.Locale.IsNull() {

		// Run the API call
		resp.Diagnostics.Append(legacysdk.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				pagedIterator := r.Client.ManagementAPIClient.AgreementLanguagesResourcesApi.ReadAllAgreementLanguages(ctx, data.EnvironmentId.ValueString(), data.AgreementId.ValueString()).Execute()

				var initialHttpResponse *http.Response

				for pageCursor, err := range pagedIterator {
					if err != nil {
						return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, pageCursor.HTTPResponse, err)
					}

					if initialHttpResponse == nil {
						initialHttpResponse = pageCursor.HTTPResponse
					}

					if agreementLocalizations, ok := pageCursor.EntityArray.Embedded.GetLanguagesOk(); ok {

						for _, localizationItem := range agreementLocalizations {

							if (!data.DisplayName.IsNull() && localizationItem.AgreementLanguage.GetDisplayName() == data.DisplayName.ValueString()) ||
								(!data.Locale.IsNull() && localizationItem.AgreementLanguage.GetLocale() == data.Locale.ValueString()) {
								return localizationItem.AgreementLanguage, pageCursor.HTTPResponse, nil
							}

						}
					}
				}

				return nil, initialHttpResponse, nil
			},
			"ReadAllAgreementLanguages",
			legacysdk.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			&agreementLocalization,
		)...)
		if resp.Diagnostics.HasError() {
			return
		}

		if agreementLocalization == nil {
			var identifier string
			if !data.DisplayName.IsNull() {
				identifier = data.DisplayName.String()
			} else if !data.Locale.IsNull() {
				identifier = data.Locale.String()
			}

			resp.Diagnostics.AddError(
				"Cannot find agreement localization from name or locale",
				fmt.Sprintf("The agreement localization %s for environment %s cannot be found", identifier, data.EnvironmentId.String()),
			)
			return
		}

	} else if !data.AgreementLocalizationId.IsNull() {

		// Run the API call
		resp.Diagnostics.Append(legacysdk.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				fO, fR, fErr := r.Client.ManagementAPIClient.AgreementLanguagesResourcesApi.ReadOneAgreementLanguage(ctx, data.EnvironmentId.ValueString(), data.AgreementId.ValueString(), data.AgreementLocalizationId.ValueString()).Execute()
				return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
			},
			"ReadOneAgreementLanguage",
			legacysdk.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			&agreementLocalization,
		)...)
		if resp.Diagnostics.HasError() {
			return
		}

	} else {
		resp.Diagnostics.AddError(
			"Missing parameter",
			"Cannot find the requested agreement localization. agreement_localization_id or display_name must be set.",
		)
		return
	}

	languageResponse, _ := findLanguageByLocale(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), agreementLocalization.GetLocale())

	if languageResponse == nil {
		resp.Diagnostics.AddError(
			"Cannot find language from locale",
			fmt.Sprintf("Cannot find the requested language from the locale %s of the agreement localization %s in environment %s.  Please report this error to the provider maintainers.", agreementLocalization.GetLocale(), agreementLocalization.GetId(), data.EnvironmentId.String()),
		)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(agreementLocalization, languageResponse.GetId())...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (p *AgreementLocalizationDataSourceModel) toState(apiObject *management.AgreementLanguage, languageId string) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.PingOneResourceIDToTF(apiObject.GetId())
	p.AgreementLocalizationId = framework.PingOneResourceIDToTF(apiObject.GetId())
	p.LanguageId = framework.PingOneResourceIDToTF(languageId)
	p.DisplayName = framework.StringOkToTF(apiObject.GetDisplayNameOk())
	p.Locale = framework.StringOkToTF(apiObject.GetLocaleOk())
	p.Enabled = framework.BoolOkToTF(apiObject.GetEnabledOk())

	if v, ok := apiObject.GetUserExperienceOk(); ok {
		p.UXTextCheckboxAccept = framework.StringOkToTF(v.GetAcceptCheckboxTextOk())
		p.UXTextButtonContinue = framework.StringOkToTF(v.GetContinueButtonTextOk())
		p.UXTextButtonDecline = framework.StringOkToTF(v.GetDeclineButtonTextOk())
	}

	if v, ok := apiObject.GetCurrentRevisionOk(); ok {
		p.CurrentRevisionId = framework.PingOneResourceIDOkToTF(v.GetIdOk())
	} else {
		p.CurrentRevisionId = pingonetypes.NewResourceIDNull()
	}

	return diags
}
