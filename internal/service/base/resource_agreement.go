package base

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/boolplanmodifierinternal"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
)

// Types
type AgreementResource struct {
	client *management.APIClient
	region model.RegionMapping
}

type AgreementResourceModel struct {
	Id                  types.String  `tfsdk:"id"`
	EnvironmentId       types.String  `tfsdk:"environment_id"`
	Name                types.String  `tfsdk:"name"`
	Description         types.String  `tfsdk:"description"`
	Enabled             types.Bool    `tfsdk:"enabled"`
	ReconsentPeriodDays types.Float64 `tfsdk:"reconsent_period_days"`
	LocalizedText       types.Set     `tfsdk:"localized_text"`
}

type LocalizedTextModel struct {
	LanguageId           types.String `tfsdk:"language_id"`
	DisplayName          types.String `tfsdk:"display_name"`
	Locale               types.String `tfsdk:"locale"`
	Enabled              types.Bool   `tfsdk:"enabled"`
	UXTextCheckboxAccept types.String `tfsdk:"text_checkbox_accept"`
	UXTextButtonContinue types.String `tfsdk:"text_button_continue"`
	UXTextButtonDecline  types.String `tfsdk:"text_button_decline"`
	LatestRevision       types.List   `tfsdk:"latest_revision"`
}

type LatestRevisionModel struct {
	RevisionId       types.String `tfsdk:"revision_id"`
	ContentType      types.String `tfsdk:"content_type"`
	EffectiveAt      types.String `tfsdk:"effective_at"`
	NotValidAfter    types.String `tfsdk:"not_valid_after"`
	RequireReconsent types.Bool   `tfsdk:"require_reconsent"`
	Text             types.String `tfsdk:"text"`
}

var (
	localizedTextTFObjectTypes = map[string]attr.Type{
		"language_id":          types.StringType,
		"display_name":         types.StringType,
		"locale":               types.StringType,
		"enabled":              types.BoolType,
		"text_checkbox_accept": types.StringType,
		"text_button_continue": types.StringType,
		"text_button_decline":  types.StringType,
		"latest_revision":      types.ListType{ElemType: types.ObjectType{AttrTypes: revisionTFObjectTypes}},
	}

	revisionTFObjectTypes = map[string]attr.Type{
		"revision_id":       types.StringType,
		"content_type":      types.StringType,
		"effective_at":      types.StringType,
		"not_valid_after":   types.StringType,
		"require_reconsent": types.BoolType,
		"text":              types.StringType,
	}
)

// Framework interfaces
var (
	_ resource.Resource                = &AgreementResource{}
	_ resource.ResourceWithConfigure   = &AgreementResource{}
	_ resource.ResourceWithImportState = &AgreementResource{}
)

// New Object
func NewAgreementResource() resource.Resource {
	return &AgreementResource{}
}

// Metadata
func (r *AgreementResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_agreement"
}

// Schema.
func (r *AgreementResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	const attrMinLength = 1

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage agreements in a PingOne environment.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_EnvironmentID(framework.SchemaDescription{
				Description: "The ID of the environment to associate the agreement with."},
			),

			"name": schema.StringAttribute{
				Description: "A string that specifies the name of the agreement to configure.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			"description": schema.StringAttribute{
				Description: "A string that specifies the description of the agreement.",
				Optional:    true,
			},

			"enabled": schema.BoolAttribute{
				Description: "A boolean that specifies the current enabled state of the agreement. The agreement must support the default language to be enabled. It cannot be disabled if it is referenced by a sign-on policy action. When an agreement is disabled, it is not used anywhere that it is configured across PingOne.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifierinternal.BoolDefaultValue(
						types.BoolValue(true),
						"The default value for this attribute is \"true\".",
						"The default value for this attribute is `true`.",
					),
				},
			},

			"reconsent_period_days": schema.Float64Attribute{
				Description: "A number that specifies the number of days until a consent to this agreement expires.",
				Optional:    true,
			},
		},

		Blocks: map[string]schema.Block{
			"localized_text": schema.SetNestedBlock{
				Description: "One or more blocks of localized agreement text",

				NestedObject: schema.NestedBlockObject{

					Attributes: map[string]schema.Attribute{
						"language_id": schema.StringAttribute{
							Description: "A string that specifies the UUID that identifies the language ID.",
							Computed:    true,
						},

						"display_name": schema.StringAttribute{
							Description: "A string used as the title of the agreement for the language presented to the user.",
							Required:    true,
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(attrMinLength),
							},
						},

						"locale": schema.StringAttribute{
							Description: "", // TODO
							Required:    true,
							Validators: []validator.String{ // TODO
								stringvalidator.LengthAtLeast(attrMinLength),
							},
						},

						"enabled": schema.BoolAttribute{
							Description: "A boolean that specifies whether a localized text is enabled in the agreement.",
							Optional:    true,
							PlanModifiers: []planmodifier.Bool{
								boolplanmodifierinternal.BoolDefaultValue(
									types.BoolValue(true),
									"The default value for this attribute is \"true\".",
									"The default value for this attribute is `true`.",
								),
							},
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

					Blocks: map[string]schema.Block{
						"latest_revision": schema.ListNestedBlock{
							Description: "A single block of agreement text for the latest revision.  Changes to attributes in this block will create a new agreement revision in PingOne.  Once a new revision has been created, previous revisions are no longer modifiable.",

							NestedObject: schema.NestedBlockObject{

								Attributes: map[string]schema.Attribute{
									"revision_id": schema.StringAttribute{
										Description: "A string that specifies the UUID that identifies the revision ID.",
										Computed:    true,
									},

									"content_type": schema.StringAttribute{
										Description:         "The content type of text. Options are text/html and text/plain, as defined by rfc-6838 (https://datatracker.ietf.org/doc/html/rfc6838#section-4.2.1) and Media Types/text (https://www.iana.org/assignments/media-types/media-types.xhtml#text).",
										MarkdownDescription: "The content type of `text`. Options are `text/html` and `text/plain`, as defined by [rfc-6838](https://datatracker.ietf.org/doc/html/rfc6838#section-4.2.1) and [Media Types/text](https://www.iana.org/assignments/media-types/media-types.xhtml#text).",
										Required:            true,
										Validators: []validator.String{
											stringvalidator.LengthAtLeast(attrMinLength),
										},
									},

									"effective_at": schema.StringAttribute{
										Description: "The start date that the revision is presented to users.  The effective date must be unique for each language agreement, and the property value can be the present date or a future date only.",
										Required:    true,
										Validators: []validator.String{ // TODO (RFC3339)
											stringvalidator.LengthAtLeast(attrMinLength),
										},
									},

									"not_valid_after": schema.StringAttribute{
										Description: "Specifies whether the revision is still valid in the context of all revisions for a language. This property is calculated dynamically at read time, taking into consideration the agreement language, the language enabled property, and the agreement enabled property. When a new revision is added, this attribute's property values for all other previous revisions might be impacted. For example, if a new revision becomes effective and it forces reconsent, then all older revisions are no longer valid.",
										Computed:    true,
									},

									"require_reconsent": schema.BoolAttribute{
										Description: "Whether the user is required to provide consent to the language revision after it becomes effective.",
										Required:    true,
									},

									"text": schema.StringAttribute{
										Description:         "Text or HTML for the revision. HTML support includes \"tags\" (italicize, bold, links, headers, paragraph, line breaks), \"link (a) tags\" (allow href, style, target attributes), \"block tags (p, b, h)\" (allow style and align attributes).",
										MarkdownDescription: "Text or HTML for the revision. HTML support includes **tags** (italicize, bold, links, headers, paragraph, line breaks), **link (a) tags** (allow href, style, target attributes), **block tags (p, b, h)** (allow style and align attributes).",
										Required:            true,
										Validators: []validator.String{
											stringvalidator.LengthAtLeast(attrMinLength),
										},
									},
								},
							},
							Validators: []validator.List{
								listvalidator.IsRequired(),
								listvalidator.SizeAtLeast(1),
								listvalidator.SizeAtMost(1),
							},
						},
					},
				},
				Validators: []validator.Set{
					setvalidator.IsRequired(),
					setvalidator.SizeAtLeast(1),
				},
			},
		},
	}
}

func (r *AgreementResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *AgreementResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state AgreementResourceModel

	if r.client == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": r.region.URLSuffix,
	})

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	/* ********
	Validation
	******** */
	var localizedTextPlan []LocalizedTextModel
	resp.Diagnostics.Append(plan.LocalizedText.ElementsAs(ctx, &localizedTextPlan, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if len(localizedTextPlan) < 1 {
		resp.Diagnostics.AddError(
			"Missing parameter",
			"At least one `localized_text` block must be configured.",
		)
		return
	}

	for _, localizedTextItem := range localizedTextPlan {
		language, _ := findLanguageByLocale(ctx, r.client, plan.EnvironmentId.ValueString(), localizedTextItem.Locale.ValueString())

		if language == nil || !language.GetEnabled() {
			resp.Diagnostics.AddError(
				"Invalid parameter",
				fmt.Sprintf("The language for locale %s needs to be configured and enabled for the environment.  Hint: use the `pingone_language` and `pingone_language_update` resources.", localizedTextItem.Locale.ValueString()),
			)
			return
		}
	}

	/* ********
	The base agreement
	******** */

	// Build the model for the API
	createAgreement := plan.expand(true)

	// Run the API call
	createResponse, d := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.AgreementsResourcesApi.CreateAgreement(ctx, plan.EnvironmentId.ValueString()).Agreement(*createAgreement).Execute()
		},
		"CreateAgreement",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	agreementId := createResponse.(*management.Agreement).GetId()

	/* ********
	Agreement - Localized text
	******** */

	var localizedTextResponses = make([]struct {
		agreementLanguage *management.AgreementLanguage
		revision          *management.AgreementLanguageRevision
	}, 0)
	for _, localizedTextItem := range localizedTextPlan {

		// Build the model for the API
		agreementLanguage := localizedTextItem.expand()

		// Run the API call
		localizedTextResponse, d := framework.ParseResponse(
			ctx,

			func() (interface{}, *http.Response, error) {
				return r.client.AgreementLanguagesResourcesApi.CreateAgreementLanguage(ctx, plan.EnvironmentId.ValueString(), agreementId).AgreementLanguage(*agreementLanguage).Execute()
			},
			"CreateAgreementLanguage-Create",
			framework.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
		)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {

			diags := deleteAgreement(ctx, r.client, plan.EnvironmentId.ValueString(), agreementId)
			resp.Diagnostics.Append(diags...)

			return
		}

		agreementLanguageId := localizedTextResponse.(*management.AgreementLanguage).GetId()

		/* ********
		Agreement - Localized text - revisions
		******** */
		var revisionPlan []LatestRevisionModel
		resp.Diagnostics.Append(localizedTextItem.LatestRevision.ElementsAs(ctx, &revisionPlan, false)...)
		if resp.Diagnostics.HasError() {

			diags := deleteAgreement(ctx, r.client, plan.EnvironmentId.ValueString(), agreementId)
			resp.Diagnostics.Append(diags...)

			return
		}

		if len(revisionPlan) < 1 || len(revisionPlan) > 1 {
			resp.Diagnostics.AddError(
				"Invalid parameter",
				"Exactly one `latest_revision` block must be configured.",
			)

			diags := deleteAgreement(ctx, r.client, plan.EnvironmentId.ValueString(), agreementId)
			resp.Diagnostics.Append(diags...)

			return
		}

		// Build the model for the API
		revision, d := revisionPlan[0].expand()
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {

			diags := deleteAgreement(ctx, r.client, plan.EnvironmentId.ValueString(), agreementId)
			resp.Diagnostics.Append(diags...)

			return
		}

		// Run the API call
		revisionResponse, d := framework.ParseResponse(
			ctx,

			func() (interface{}, *http.Response, error) {
				return r.client.AgreementRevisionsResourcesApi.CreateAgreementLanguageRevision(ctx, plan.EnvironmentId.ValueString(), agreementId, agreementLanguageId).AgreementLanguageRevision(*revision).Execute()
			},
			"CreateAgreementLanguageRevision-Create",
			framework.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
		)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {

			diags := deleteAgreement(ctx, r.client, plan.EnvironmentId.ValueString(), agreementId)
			resp.Diagnostics.Append(diags...)

			return
		}

		localizedTextResponses = append(localizedTextResponses, struct {
			agreementLanguage *management.AgreementLanguage
			revision          *management.AgreementLanguageRevision
		}{
			agreementLanguage: localizedTextResponse.(*management.AgreementLanguage),
			revision:          revisionResponse.(*management.AgreementLanguageRevision),
		})

	}

	/* ********
	Agreement - Final agreement - we do this to round off the attributes that are required to be set after the rest of the agreement is populated
	******** */
	agreement := plan.expand(false)

	response, d := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.AgreementsResourcesApi.UpdateAgreement(ctx, plan.EnvironmentId.ValueString(), agreementId).Agreement(*agreement).Execute()
		},
		"UpdateAgreement-Create",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		diags := deleteAgreement(ctx, r.client, plan.EnvironmentId.ValueString(), agreementId)
		resp.Diagnostics.Append(diags...)

		return
	}

	// Create the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(ctx, response.(*management.Agreement), localizedTextResponses)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *AgreementResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *AgreementResourceModel

	if r.client == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": r.region.URLSuffix,
	})

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	/* ********
	The base agreement
	******** */
	// Run the API call
	response, diags := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.AgreementsResourcesApi.ReadOneAgreement(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
		},
		"ReadOneAgreement-Read",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Remove from state if resource is not found
	if response == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	/* ********
	Agreement - Localized text
	******** */

	// Run the API call
	localizedTextsResponse, diags := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.AgreementLanguagesResourcesApi.ReadAllAgreementLanguages(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
		},
		"ReadAllAgreementLanguages-Read",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	entityArray := localizedTextsResponse.(*management.EntityArray)

	var localizedTextResponses = make([]struct {
		agreementLanguage *management.AgreementLanguage
		revision          *management.AgreementLanguageRevision
	}, 0)
	if languages, ok := entityArray.Embedded.GetLanguagesOk(); ok {

		for _, language := range languages {

			languageId := language.AgreementLanguage.GetId()
			latestRevisionId := language.AgreementLanguage.GetCurrentRevision().Id

			/* ********
			Agreement - Localized text - revisions
			******** */
			revisionResponse, diags := framework.ParseResponse(
				ctx,

				func() (interface{}, *http.Response, error) {
					return r.client.AgreementRevisionsResourcesApi.ReadOneAgreementLanguageRevision(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString(), languageId, *latestRevisionId).Execute()
				},
				"ReadOneAgreementLanguageRevision-Read",
				framework.CustomErrorResourceNotFoundWarning,
				sdk.DefaultCreateReadRetryable,
			)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}

			localizedTextResponses = append(localizedTextResponses, struct {
				agreementLanguage *management.AgreementLanguage
				revision          *management.AgreementLanguageRevision
			}{
				agreementLanguage: language.AgreementLanguage,
				revision:          revisionResponse.(*management.AgreementLanguageRevision),
			})
		}

	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(ctx, response.(*management.Agreement), localizedTextResponses)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AgreementResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *AgreementResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *AgreementResourceModel

	if r.client == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": r.region.URLSuffix,
	})

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	diags := deleteAgreement(ctx, r.client, data.EnvironmentId.ValueString(), data.Id.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func deleteAgreement(ctx context.Context, apiClient *management.APIClient, environmentId, agreementId string) diag.Diagnostics {
	var diags diag.Diagnostics

	_, diags = framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			r, err := apiClient.AgreementsResourcesApi.DeleteAgreement(ctx, environmentId, agreementId).Execute()
			return nil, r, err
		},
		"DeleteAgreement",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
	)

	return diags
}

func (r *AgreementResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splitLength := 2
	attributes := strings.SplitN(req.ID, "/", splitLength)

	if len(attributes) != splitLength {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("invalid id (\"%s\") specified, should be in format \"environment_id/agreement_id\"", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_id"), attributes[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), attributes[1])...)
}

func (p *AgreementResourceModel) expand(initialCreate bool) *management.Agreement {

	var agreementEnabled bool

	if initialCreate {
		agreementEnabled = false
	} else {
		agreementEnabled = p.Enabled.ValueBool()
	}

	data := management.NewAgreement(agreementEnabled, p.Name.ValueString())

	if !p.Description.IsNull() && !p.Description.IsUnknown() {
		data.SetDescription(p.Description.ValueString())
	}

	if !p.ReconsentPeriodDays.IsNull() && !p.ReconsentPeriodDays.IsUnknown() {
		data.SetReconsentPeriodDays(float32(p.ReconsentPeriodDays.ValueFloat64()))
	}

	return data
}

func (p *LocalizedTextModel) expand() *management.AgreementLanguage {
	data := management.NewAgreementLanguage(
		p.DisplayName.ValueString(),
		p.Enabled.ValueBool(),
		p.Locale.ValueString(),
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

func (p *LatestRevisionModel) expand() (*management.AgreementLanguageRevision, diag.Diagnostics) {
	var diags diag.Diagnostics

	t, e := time.Parse(time.RFC3339, p.EffectiveAt.ValueString())
	if e != nil {
		diags.AddError(
			"Invalid data format",
			"Cannot convert effectve_at to a date/time.  Please check the format is a valid RFC3339 date time format.")
		return nil, diags
	}

	data := management.NewAgreementLanguageRevision(
		management.EnumAgreementRevisionContentType(p.ContentType.ValueString()),
		t,
		p.RequireReconsent.ValueBool(),
		p.Text.ValueString(),
	)

	return data, diags
}

func (p *AgreementResourceModel) toState(ctx context.Context, apiObject *management.Agreement, localizedTexts []struct {
	agreementLanguage *management.AgreementLanguage
	revision          *management.AgreementLanguageRevision
}) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.StringToTF(apiObject.GetId())
	p.EnvironmentId = framework.StringToTF(*apiObject.GetEnvironment().Id)
	p.Name = framework.StringOkToTF(apiObject.GetNameOk())
	p.Description = framework.StringOkToTF(apiObject.GetDescriptionOk())
	p.Enabled = framework.BoolOkToTF(apiObject.GetEnabledOk())
	p.ReconsentPeriodDays = framework.Float32OkToTF(apiObject.GetReconsentPeriodDaysOk())

	localizedText, d := toStateLocalizedTexts(ctx, localizedTexts)
	diags.Append(d...)
	p.LocalizedText = localizedText

	return diags
}

func toStateLocalizedTexts(ctx context.Context, localizedTexts []struct {
	agreementLanguage *management.AgreementLanguage
	revision          *management.AgreementLanguageRevision
}) (types.Set, diag.Diagnostics) {
	var diags diag.Diagnostics
	tfObjType := types.ObjectType{AttrTypes: localizedTextTFObjectTypes}

	if len(localizedTexts) < 1 {
		return types.SetValueMust(tfObjType, []attr.Value{}), diags
	}

	setOfObj := []attr.Value{}

	for _, v := range localizedTexts {

		localizedTextMap := map[string]attr.Value{
			"language_id":  framework.StringOkToTF(v.agreementLanguage.GetIdOk()),
			"display_name": framework.StringOkToTF(v.agreementLanguage.GetDisplayNameOk()),
			"locale":       framework.StringOkToTF(v.agreementLanguage.GetLocaleOk()),
			"enabled":      framework.BoolOkToTF(v.agreementLanguage.GetEnabledOk()),
		}

		if ux, ok := v.agreementLanguage.GetUserExperienceOk(); ok {
			localizedTextMap["text_checkbox_accept"] = framework.StringOkToTF(ux.GetAcceptCheckboxTextOk())
			localizedTextMap["text_button_continue"] = framework.StringOkToTF(ux.GetContinueButtonTextOk())
			localizedTextMap["text_button_decline"] = framework.StringOkToTF(ux.GetDeclineButtonTextOk())
		} else {
			localizedTextMap["text_checkbox_accept"] = types.StringNull()
			localizedTextMap["text_button_continue"] = types.StringNull()
			localizedTextMap["text_button_decline"] = types.StringNull()
		}

		latestRevision, d := toStateAgreementRevision(v.revision)
		diags.Append(d...)
		if diags.HasError() {
			return types.SetValueMust(tfObjType, []attr.Value{}), diags
		}

		localizedTextMap["latest_revision"] = latestRevision

		tflog.Debug(ctx, "Revision created", map[string]interface{}{
			"localizedTextMap[\"latest_revision\"]": localizedTextMap["latest_revision"],
			"localizedTextMap":                      localizedTextMap,
		})

		flattenedObj, d := types.ObjectValue(localizedTextTFObjectTypes, localizedTextMap)
		diags.Append(d...)

		setOfObj = append(setOfObj, flattenedObj)

	}

	returnVar, d := types.SetValue(tfObjType, setOfObj)
	diags.Append(d...)

	return returnVar, diags
}

func toStateAgreementRevision(revision *management.AgreementLanguageRevision) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	tfObjType := types.ObjectType{AttrTypes: revisionTFObjectTypes}

	if revision == nil {
		return types.ListValueMust(tfObjType, []attr.Value{}), diags
	}

	revisionMap := map[string]attr.Value{
		"revision_id":       framework.StringOkToTF(revision.GetIdOk()),
		"content_type":      enumAgreementRevisionContentTypeOkToTF(revision.GetContentTypeOk()),
		"effective_at":      framework.TimeOkToTF(revision.GetEffectiveAtOk()),
		"not_valid_after":   framework.TimeOkToTF(revision.GetNotValidAfterOk()),
		"require_reconsent": framework.BoolOkToTF(revision.GetRequireReconsentOk()),
		"text":              framework.StringOkToTF(revision.GetTextOk()),
	}

	flattenedObj, d := types.ObjectValue(revisionTFObjectTypes, revisionMap)
	diags.Append(d...)

	returnVar, d := types.ListValue(tfObjType, append([]attr.Value{}, flattenedObj))
	diags.Append(d...)

	return returnVar, diags
}

func enumAgreementRevisionContentTypeOkToTF(v *management.EnumAgreementRevisionContentType, ok bool) basetypes.StringValue {
	if !ok || v == nil {
		return types.StringNull()
	} else {
		return types.StringValue(string(*v))
	}
}
