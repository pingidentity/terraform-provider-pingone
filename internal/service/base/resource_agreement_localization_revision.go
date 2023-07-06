package base

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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type AgreementLocalizationRevisionResource struct {
	managementClient *management.APIClient
	region           model.RegionMapping
}

type AgreementLocalizationRevisionResourceModel struct {
	Id                      types.String `tfsdk:"id"`
	EnvironmentId           types.String `tfsdk:"environment_id"`
	AgreementId             types.String `tfsdk:"agreement_id"`
	AgreementLocalizationId types.String `tfsdk:"agreement_localization_id"`
	ContentType             types.String `tfsdk:"content_type"`
	EffectiveAt             types.String `tfsdk:"effective_at"`
	NotValidAfter           types.String `tfsdk:"not_valid_after"`
	RequireReconsent        types.Bool   `tfsdk:"require_reconsent"`
	Text                    types.String `tfsdk:"text"`
}

// Framework interfaces
var (
	_ resource.Resource                = &AgreementLocalizationRevisionResource{}
	_ resource.ResourceWithConfigure   = &AgreementLocalizationRevisionResource{}
	_ resource.ResourceWithImportState = &AgreementLocalizationRevisionResource{}
)

// New Object
func NewAgreementLocalizationRevisionResource() resource.Resource {
	return &AgreementLocalizationRevisionResource{}
}

// Metadata
func (r *AgreementLocalizationRevisionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_agreement_localization_revision"
}

// Schema.
func (r *AgreementLocalizationRevisionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	const attrMinLength = 1

	contentTypeFmt := "The content type to apply to the revision text configured in the `text` parameter. Options are `text/html` and `text/plain`, as defined by %s."
	contentTypeDescription := framework.SchemaAttributeDescription{
		MarkdownDescription: fmt.Sprintf(contentTypeFmt, "[rfc-6838](https://datatracker.ietf.org/doc/html/rfc6838#section-4.2.1) and [Media Types/text](https://www.iana.org/assignments/media-types/media-types.xhtml#text)"),
		Description:         fmt.Sprintf(strings.ReplaceAll(contentTypeFmt, "`", "\""), "rfc-6838 (https://datatracker.ietf.org/doc/html/rfc6838#section-4.2.1) and  Media Types/text (https://www.iana.org/assignments/media-types/media-types.xhtml#text)"),
	}

	// "Text or HTML for the revision. HTML support includes \"tags\" (italicize, bold, links, headers, paragraph, line breaks), \"link (a) tags\" (allow href, style, target attributes), \"block tags (p, b, h)\" (allow style and align attributes).",
	// "Text or HTML for the revision. HTML support includes **tags** (italicize, bold, links, headers, paragraph, line breaks), **link (a) tags** (allow href, style, target attributes), **block tags (p, b, h)** (allow style and align attributes).",

	textDescriptionFmt := "Text or HTML for the revision. HTML support includes **tags** (italicize, bold, links, headers, paragraph, line breaks), **link (a) tags** (allow href, style, target attributes), **block tags (p, b, h)** (allow style and align attributes)."
	textDescription := framework.SchemaAttributeDescription{
		MarkdownDescription: textDescriptionFmt,
		Description:         strings.ReplaceAll(textDescriptionFmt, "**", "\""),
	}

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage agreement localization revisions in a PingOne environment.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to associate the agreement localization revision with."),
			),

			"agreement_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the agreement to associate the agreement localization revision with."),
			),

			"agreement_localization_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the agreement localization to associate the revision with."),
			),

			"content_type": schema.StringAttribute{
				Description:         contentTypeDescription.Description,
				MarkdownDescription: contentTypeDescription.MarkdownDescription,
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(string(management.ENUMAGREEMENTREVISIONCONTENTTYPE_HTML), string(management.ENUMAGREEMENTREVISIONCONTENTTYPE_PLAIN)),
				},
			},

			"effective_at": schema.StringAttribute{
				Description: "The start date that the revision is presented to users.  The effective date must be unique for each language agreement, and the property value can be the present date or a future date only.  Must be a valid RFC3339 date/time string.  If left undefined, will default to the current date and time (the revision will be effective immediately).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(verify.RFC3339Regexp, "Attribute must be a valid RFC3339 date/time string."),
				},
			},

			"not_valid_after": schema.StringAttribute{
				Description: "Specifies whether the revision is still valid in the context of all revisions for a language. This property is calculated dynamically at read time, taking into consideration the agreement language, the language enabled property, and the agreement enabled property. When a new revision is added, this attribute's property values for all other previous revisions might be impacted. For example, if a new revision becomes effective and it forces reconsent, then all older revisions are no longer valid.",
				Computed:    true,
			},

			"require_reconsent": schema.BoolAttribute{
				Description: "Whether the user is required to provide a renewed consent to the language revision after it becomes effective.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
				Required: true,
			},

			"text": schema.StringAttribute{
				Description:         textDescription.Description,
				MarkdownDescription: textDescription.MarkdownDescription,
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},
		},
	}
}

func (r *AgreementLocalizationRevisionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.managementClient = preparedClient
	r.region = resourceConfig.Client.API.Region
}

func (r *AgreementLocalizationRevisionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state AgreementLocalizationRevisionResourceModel

	if r.managementClient == nil {
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
	localizationRevision, d := plan.expand()
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *management.AgreementLanguageRevision
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return r.managementClient.AgreementRevisionsResourcesApi.CreateAgreementLanguageRevision(ctx, plan.EnvironmentId.ValueString(), plan.AgreementId.ValueString(), plan.AgreementLocalizationId.ValueString()).AgreementLanguageRevision(*localizationRevision).Execute()
		},
		"CreateAgreementLanguageRevision",
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
	resp.Diagnostics.Append(state.toState(response, localizationRevision.Text)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *AgreementLocalizationRevisionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *AgreementLocalizationRevisionResourceModel

	if r.managementClient == nil {
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
	var response *management.AgreementLanguageRevision
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return r.managementClient.AgreementRevisionsResourcesApi.ReadOneAgreementLanguageRevision(ctx, data.EnvironmentId.ValueString(), data.AgreementId.ValueString(), data.AgreementLocalizationId.ValueString(), data.Id.ValueString()).Execute()
		},
		"ReadOneAgreementLanguageRevision",
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
	resp.Diagnostics.Append(data.toState(response, data.Text.ValueString())...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AgreementLocalizationRevisionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *AgreementLocalizationRevisionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *AgreementLocalizationRevisionResourceModel

	if r.managementClient == nil {
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
			r, err := r.managementClient.AgreementRevisionsResourcesApi.DeleteAgreementLanguageRevision(ctx, data.EnvironmentId.ValueString(), data.AgreementId.ValueString(), data.AgreementLocalizationId.ValueString(), data.Id.ValueString()).Execute()
			return nil, r, err
		},
		"DeleteAgreementLanguageRevision",
		agreementLocalizationRevisionDeleteErrorHandler,
		sdk.DefaultCreateReadRetryable,
		nil,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *AgreementLocalizationRevisionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splitLength := 4
	attributes := strings.SplitN(req.ID, "/", splitLength)

	if len(attributes) != splitLength {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("invalid id (\"%s\") specified, should be in format \"environment_id/agreement_id/agreement_localization_id/agreement_localization_revision_id\"", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_id"), attributes[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("agreement_id"), attributes[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("agreement_localization_id"), attributes[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), attributes[1])...)
}

func (p *AgreementLocalizationRevisionResourceModel) expand() (*management.AgreementLanguageRevision, diag.Diagnostics) {
	var diags diag.Diagnostics

	var t time.Time

	if !p.EffectiveAt.IsNull() && !p.EffectiveAt.IsUnknown() {
		var e error
		t, e = time.Parse(time.RFC3339, p.EffectiveAt.ValueString())
		if e != nil {
			diags.AddError(
				"Invalid data format",
				"Cannot convert effectve_at to a date/time.  Please check the format is a valid RFC3339 date time format.")
			return nil, diags
		}
	} else {
		t = time.Now().Local()
	}

	data := management.NewAgreementLanguageRevision(
		management.EnumAgreementRevisionContentType(p.ContentType.ValueString()),
		t,
		p.RequireReconsent.ValueBool(),
		p.Text.ValueString(),
	)

	return data, diags
}

func (p *AgreementLocalizationRevisionResourceModel) toState(apiObject *management.AgreementLanguageRevision, revisionText string) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.StringToTF(apiObject.GetId())
	p.AgreementId = framework.StringToTF(*apiObject.GetAgreement().Id)
	p.AgreementLocalizationId = framework.StringToTF(*apiObject.GetLanguage().Id)
	p.ContentType = framework.EnumOkToTF(apiObject.GetContentTypeOk())
	p.EffectiveAt = framework.TimeOkToTF(apiObject.GetEffectiveAtOk())
	p.NotValidAfter = framework.TimeOkToTF(apiObject.GetNotValidAfterOk())
	p.RequireReconsent = framework.BoolOkToTF(apiObject.GetRequireReconsentOk())
	p.Text = framework.StringToTF(revisionText)

	return diags
}

func agreementLocalizationRevisionDeleteErrorHandler(error model.P1Error) diag.Diagnostics {
	var diags diag.Diagnostics

	// Deleted outside of TF
	if error.GetCode() == "NOT_FOUND" {
		diags.AddWarning(
			"Resource not found on delete.",
			error.GetMessage(),
		)

		return diags
	}

	// Last action in the policy
	if v, ok := error.GetDetailsOk(); ok && v != nil && len(v) > 0 {
		if v[0].GetCode() == "CONSTRAINT_VIOLATION" {
			if match, _ := regexp.MatchString("A currently effective revision cannot be deleted.", v[0].GetMessage()); match {
				diags.AddWarning(
					"Cannot delete the agreement localization revision, a currently effective revision cannot be deleted.",
					"The revision is left in place but no longer managed by the provider.",
				)

				return diags
			}
		}
	}

	return nil
}
