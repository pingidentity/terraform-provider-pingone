package base

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type FormResource serviceClientType

type formResourceModel struct {
	Id                types.String `tfsdk:"id"`
	EnvironmentId     types.String `tfsdk:"environment_id"`
	Name              types.String `tfsdk:"name"`
	Description       types.String `tfsdk:"description"`
	Category          types.String `tfsdk:"category"`
	Cols              types.Int64  `tfsdk:"cols"`
	Components        types.Object `tfsdk:"components"`
	FieldTypes        types.Set    `tfsdk:"field_types"`
	LanguageBundle    types.Map    `tfsdk:"language_bundle"`
	MarkOptional      types.Bool   `tfsdk:"mark_optional"`
	MarkRequired      types.Bool   `tfsdk:"mark_required"`
	TranslationMethod types.String `tfsdk:"translation_method"`
}

type formComponentsResourceModel struct {
	Fields types.Set `tfsdk:"fields"`
}

// Framework interfaces
var (
	_ resource.Resource                = &FormResource{}
	_ resource.ResourceWithConfigure   = &FormResource{}
	_ resource.ResourceWithImportState = &FormResource{}
)

// New Object
func NewFormResource() resource.Resource {
	return &FormResource{}
}

// Metadata
func (r *FormResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_form"
}

// Schema.
func (r *FormResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	const attrMinLength = 1
	const colsMinValue = 1
	const colsMaxValue = 4

	nameDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the form name, which must be provided and must be unique within an environment.",
	)

	descriptionDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the description of the form.",
	)

	categoryDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the type of form.",
	).AllowedValuesComplex(map[string]string{
		string(management.ENUMFORMCATEGORY_CUSTOM): "allows the form to be built with fields that do not map specifically to the PingOne directory attributes",
	},
	).DefaultValue(string(management.ENUMFORMCATEGORY_CUSTOM))

	colsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("An integer that specifies the number of columns in the form (min = `%d`; max = `%d`).", colsMinValue, colsMaxValue),
	).DefaultValue("UNKNOWN")

	componentsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that specifies the form configuration elements.",
	)

	componentFieldsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A set of objects that specifies the form fields that make up the form.",
	)

	fieldTypesDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A set of strings that specifies the field types in the form.",
	).AllowedValuesEnum(management.AllowedEnumFormFieldTypeEnumValues)

	languageBundleDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An map of strings that provides i18n keys to their translations. This object includes both the keys and their default translations. The PingOne language management service finds this object, and creates the new keys for translation for this form.",
	)

	markOptionalDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether optional fields are highlighted in the rendered form.",
	)

	markRequiredDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether required fields are highlighted in the rendered form.",
	)

	translationMethodDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies how to translate the text strings in the form.",
	).AllowedValuesEnum(management.AllowedEnumFormTranslationMethodEnumValues).DefaultValue("UNKNOWN")

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage PingOne forms for an environment.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to manage the form in."),
			),

			"name": schema.StringAttribute{
				Description:         nameDescription.Description,
				MarkdownDescription: nameDescription.MarkdownDescription,
				Required:            true,

				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			"description": schema.StringAttribute{
				Description:         descriptionDescription.Description,
				MarkdownDescription: descriptionDescription.MarkdownDescription,
				Optional:            true,
			},

			"category": schema.StringAttribute{
				Description:         categoryDescription.Description,
				MarkdownDescription: categoryDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,

				Default: stringdefault.StaticString(string(management.ENUMFORMCATEGORY_CUSTOM)),
			},

			"cols": schema.StringAttribute{
				Description:         colsDescription.Description,
				MarkdownDescription: colsDescription.MarkdownDescription,
				Optional:            true,
			},

			"components": schema.SingleNestedAttribute{
				Description:         componentsDescription.Description,
				MarkdownDescription: componentsDescription.MarkdownDescription,
				Required:            true,

				Attributes: map[string]schema.Attribute{
					"fields": schema.SetNestedAttribute{
						Description:         componentFieldsDescription.Description,
						MarkdownDescription: componentFieldsDescription.MarkdownDescription,
						Required:            true,

						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{},
						},

						Validators: []validator.Set{
							setvalidator.SizeAtLeast(attrMinLength),
						},
					},
				},
			},

			"field_types": schema.SetAttribute{
				Description:         fieldTypesDescription.Description,
				MarkdownDescription: fieldTypesDescription.MarkdownDescription,
				Computed:            true,

				ElementType: types.StringType,

				// Validators: []validator.Set{
				// 	setvalidator.ValueStringsAre(
				// 		stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumFormFieldTypeEnumValues)...),
				// 	),
				// },
			},

			"language_bundle": schema.MapAttribute{
				Description:         languageBundleDescription.Description,
				MarkdownDescription: languageBundleDescription.MarkdownDescription,
				Optional:            true,

				ElementType: types.StringType,
			},

			"mark_optional": schema.BoolAttribute{
				Description:         markOptionalDescription.Description,
				MarkdownDescription: markOptionalDescription.MarkdownDescription,
				Required:            true,
			},

			"mark_required": schema.BoolAttribute{
				Description:         markRequiredDescription.Description,
				MarkdownDescription: markRequiredDescription.MarkdownDescription,
				Required:            true,
			},

			"translation_method": schema.StringAttribute{
				Description:         translationMethodDescription.Description,
				MarkdownDescription: translationMethodDescription.MarkdownDescription,
				Optional:            true,

				Validators: []validator.String{
					stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumFormTranslationMethodEnumValues)...),
				},
			},
		},
	}
}

func (r *FormResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.Client = resourceConfig.Client.API
	if r.Client == nil {
		resp.Diagnostics.AddError(
			"Client not initialised",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.",
		)
		return
	}
}

func (r *FormResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state formResourceModel

	if r.Client.ManagementAPIClient == nil {
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
	form, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *management.Form
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.FormManagementApi.CreateForm(ctx, plan.EnvironmentId.ValueString()).Form(*form).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"CreateForm",
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

func (r *FormResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *formResourceModel

	if r.Client.ManagementAPIClient == nil {
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
	var response *management.Form
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.FormManagementApi.ReadForm(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadForm",
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

func (r *FormResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state formResourceModel

	if r.Client.ManagementAPIClient == nil {
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
	form, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *management.Form
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.FormManagementApi.UpdateForm(ctx, plan.EnvironmentId.ValueString(), plan.Id.ValueString()).Form(*form).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateForm",
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

func (r *FormResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *formResourceModel

	if r.Client.ManagementAPIClient == nil {
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
			fR, fErr := r.Client.ManagementAPIClient.FormManagementApi.DeleteForm(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, fR, fErr)
		},
		"DeleteForm",
		framework.CustomErrorResourceNotFoundWarning,
		nil,
		nil,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *FormResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:     "form_id",
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
		pathForm := idComponent.Label

		if idComponent.PrimaryID {
			pathForm = "id"
		}

		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(pathForm), attributes[idComponent.Label])...)
	}
}

func (p *formResourceModel) expand(ctx context.Context) (*management.Form, diag.Diagnostics) {
	var diags diag.Diagnostics

	componentFields := make([]management.FormField, 0)
	// TODO

	data := management.NewForm(
		p.Name.ValueString(),
		management.EnumFormCategory(p.Category.ValueString()),
		*management.NewFormComponents(componentFields),
		p.MarkOptional.ValueBool(),
		p.MarkRequired.ValueBool(),
	)

	if !p.Description.IsNull() && !p.Description.IsUnknown() {
		data.SetDescription(p.Description.ValueString())
	}

	if !p.Cols.IsNull() && !p.Cols.IsUnknown() {
		data.SetCols(int32(p.Cols.ValueInt64()))
	}

	if !p.FieldTypes.IsNull() && !p.FieldTypes.IsUnknown() {
		var plan []string
		diags.Append(p.FieldTypes.ElementsAs(ctx, plan, false)...)
		if diags.HasError() {
			return nil, diags
		}

		fieldTypes := make([]management.EnumFormFieldType, 0)
		for _, v := range plan {
			fieldTypes = append(fieldTypes, management.EnumFormFieldType(v))
		}

		data.SetFieldTypes(fieldTypes)
	}

	if !p.LanguageBundle.IsNull() && !p.LanguageBundle.IsUnknown() {
		var plan map[string]string
		diags.Append(p.LanguageBundle.ElementsAs(ctx, plan, false)...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetLanguageBundle(plan)
	}

	if !p.TranslationMethod.IsNull() && !p.TranslationMethod.IsUnknown() {
		data.SetTranslationMethod(management.EnumFormTranslationMethod(p.TranslationMethod.ValueString()))
	}

	return data, diags
}

func (p *formResourceModel) toState(apiObject *management.Form) diag.Diagnostics {
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
	p.Name = framework.StringOkToTF(apiObject.GetNameOk())
	p.Description = framework.StringOkToTF(apiObject.GetDescriptionOk())
	p.Category = framework.EnumOkToTF(apiObject.GetCategoryOk())
	p.Cols = framework.Int32OkToTF(apiObject.GetColsOk())
	p.FieldTypes = framework.EnumSetOkToTF(apiObject.GetFieldTypesOk())
	p.LanguageBundle = framework.StringMapOkToTF(apiObject.GetLanguageBundleOk())
	p.MarkOptional = framework.BoolOkToTF(apiObject.GetMarkOptionalOk())
	p.MarkRequired = framework.BoolOkToTF(apiObject.GetMarkRequiredOk())
	p.TranslationMethod = framework.EnumOkToTF(apiObject.GetTranslationMethodOk())

	p.Components = formComponentsOkToTF(apiObject.GetComponentsOk())

	return diags
}

func formComponentsOkToTF(formComponents *management.FormComponents, ok bool) types.Object {
	return types.Object{}
}
