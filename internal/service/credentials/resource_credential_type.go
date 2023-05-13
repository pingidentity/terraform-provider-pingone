package credentials

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/patrickcping/pingone-go-sdk-v2/credentials"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
)

// Types
type CredentialTypeResource struct {
	client *credentials.APIClient
	region model.RegionMapping
}

type CredentialTypeResourceModel struct {
	Id                 types.String `tfsdk:"id"`
	EnvironmentId      types.String `tfsdk:"environment_id"`
	Title              types.String `tfsdk:"title"`
	Description        types.String `tfsdk:"description"`
	CardType           types.String `tfsdk:"card_type"`
	CardDesignTemplate types.String `tfsdk:"card_design_template"`
	Metadata           types.Object `tfsdk:"metadata"`
}

type MetadataModel struct {
	BackgroundImage  types.String `tfsdk:"background_image"`
	BgOpacityPercent types.Int64  `tfsdk:"bg_opacity_percent"`
	CardColor        types.String `tfsdk:"card_color"`
	Description      types.String `tfsdk:"description"`
	TextColor        types.String `tfsdk:"text_color"`
	Version          types.Int64  `tfsdk:"version"` // Watch Item - Best practice is to allow service to set, but if version of 5 or higher is not provided, creds do not appear in UI!
	LogoImage        types.String `tfsdk:"logo_image"`
	Name             types.String `tfsdk:"name"`
	Fields           types.List   `tfsdk:"fields"`
}

type FieldsModel struct {
	Id        types.String `tfsdk:"id"`
	Type      types.String `tfsdk:"type"`
	Title     types.String `tfsdk:"title"`
	IsVisible types.Bool   `tfsdk:"is_visible"`
	Attribute types.String `tfsdk:"attribute"`
	Value     types.String `tfsdk:"value"`
}

var (
	metadataServiceTFObjectTypes = map[string]attr.Type{
		"background_image":   types.StringType,
		"bg_opacity_percent": types.Int64Type,
		"card_color":         types.StringType,
		"description":        types.StringType,
		"text_color":         types.StringType,
		"version":            types.Int64Type,
		"logo_image":         types.StringType,
		"name":               types.StringType,
		"fields":             types.ListType{ElemType: types.ObjectType{AttrTypes: innerFieldsServiceTFObjectTypes}},
	}

	innerFieldsServiceTFObjectTypes = map[string]attr.Type{
		"id":         types.StringType,
		"type":       types.StringType,
		"title":      types.StringType,
		"is_visible": types.BoolType,
		"attribute":  types.StringType,
		"value":      types.StringType,
	}
)

// Framework interfaces
var (
	_ resource.Resource                = &CredentialTypeResource{}
	_ resource.ResourceWithConfigure   = &CredentialTypeResource{}
	_ resource.ResourceWithImportState = &CredentialTypeResource{}
)

// New Object
func NewCredentialTypeResource() resource.Resource {
	return &CredentialTypeResource{}
}

// Metadata
func (r *CredentialTypeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_credential_type"
}

func (r *CredentialTypeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	// schema descriptions and validation settings
	const attrMinLength = 1
	const attrMinVersion = 5
	const attrMinPercent = 0
	const attrMaxPercent = 100
	const imageMaxSize = 70000 // todo: change back to 50k after i update my test image, which oddly works at 65k

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage PingOne Credentials credential types.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(framework.SchemaDescription{
				Description: "The ID of the environment to create the credential type in."},
			),

			"title": schema.StringAttribute{
				Description: "A string that specifies the title of the credential. Verification sites are expected to be able to request the issued credential from the compatible wallet app using the credential title.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			"description": schema.StringAttribute{
				Description: "",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			"card_type": schema.StringAttribute{
				Description: "",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			"card_design_template": schema.StringAttribute{
				Description: "A string that specifies an SVG formatted image containing placeholders for the credential fields that need to be displayed in the image.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(`^<svg.*>[\s\S]*<\/svg>\s*$`), "expected value to contain a valid PingOne Credentials SVG card template."),
				},
			},

			"metadata": schema.SingleNestedAttribute{
				Description:         "",
				MarkdownDescription: "",
				Required:            true,

				Attributes: map[string]schema.Attribute{
					"background_image": schema.StringAttribute{
						Description:         "",
						MarkdownDescription: "",
						Optional:            true,
						Validators: []validator.String{
							stringvalidator.LengthAtMost(imageMaxSize),
							//stringvalidator.RegexMatches(regexp.MustCompile(`^data:image\/(\w+);base64,`), "base64encoded image must include Content-type and Content-encoding prefix, such as data:image/jpeg;base64, data:image/svg;base64, or data:image/png;base64."),
							isbase64EncodedValidator{},
						},
					},

					"bg_opacity_percent": schema.Int64Attribute{
						Description: "A numnber containing the percent opacity of the background image in the credential. High percentage opacity may make displayed text difficult to read.",
						Optional:    true,
						Validators: []validator.Int64{
							int64validator.Between(attrMinPercent, attrMaxPercent),
						},
					},

					"card_color": schema.StringAttribute{
						Description: "A string containing a 6-digit hexadecimal color code specifying the color of the credential.",
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.RegexMatches(
								regexp.MustCompile(`^#([A-Fa-f0-9]{6})$`),
								"expected value to contain a valid 6-digit hexadecimal color code, prefixed with a hash (#) symbol."),
						},
					},

					"description": schema.StringAttribute{
						Description:         "",
						MarkdownDescription: "",
						Optional:            true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(attrMinLength),
						},
					},

					"logo_image": schema.StringAttribute{
						Description:         "",
						MarkdownDescription: "",
						Optional:            true,
						Validators: []validator.String{
							isbase64EncodedValidator{},
							stringvalidator.LengthAtMost(imageMaxSize),
						},
					},

					"name": schema.StringAttribute{
						Description:         "",
						MarkdownDescription: "",
						Optional:            true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(attrMinLength),
						},
					},

					"text_color": schema.StringAttribute{
						Description: "A string containing a 6-digit hexadecimal color code specifying the color of the credential text.",
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.RegexMatches(
								regexp.MustCompile(`^#([A-Fa-f0-9]{6})$`),
								"expected value to contain a valid 6-digit hexadecimal color code, prefixed with a hash (#) symbol."),
						},
					},

					"version": schema.Int64Attribute{
						Description:         "",
						MarkdownDescription: "",
						Required:            true, // not requried in schema, but credentials will not display in P1 admin console if not provided
						Validators: []validator.Int64{
							int64validator.AtLeast(attrMinVersion),
						},
					},

					"fields": schema.ListNestedAttribute{
						Description:         "",
						MarkdownDescription: "",
						Required:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Description:         "",
									MarkdownDescription: "",
									Computed:            true,
								},
								"type": schema.StringAttribute{
									Description:         "",
									MarkdownDescription: "",
									Required:            true,
									Validators: []validator.String{
										stringvalidator.OneOf(
											string(credentials.ENUMCREDENTIALTYPEMETADATAFIELDSTYPE_ALPHANUMERIC_TEXT),
											string(credentials.ENUMCREDENTIALTYPEMETADATAFIELDSTYPE_DIRECTORY_ATTRIBUTE),
											string(credentials.ENUMCREDENTIALTYPEMETADATAFIELDSTYPE_ISSUED_TIMESTAMP)),
									},
								},
								"title": schema.StringAttribute{
									Description:         "",
									MarkdownDescription: "",
									Optional:            true,
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(attrMinLength),
									},
								},
								"attribute": schema.StringAttribute{
									Description:         "",
									MarkdownDescription: "",
									Optional:            true,
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(attrMinLength),
										stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("value")),
									},
								},
								"value": schema.StringAttribute{
									Description:         "",
									MarkdownDescription: "",
									Optional:            true,
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(attrMinLength),
										stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("attribute")),
									},
								},
								"is_visible": schema.BoolAttribute{
									Description:         "",
									MarkdownDescription: "",
									Optional:            true,
									Validators:          []validator.Bool{},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r *CredentialTypeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CredentialTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state CredentialTypeResourceModel

	if r.client == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	ctx = context.WithValue(ctx, credentials.ContextServerVariables, map[string]string{
		"suffix": r.region.URLSuffix,
	})

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the model for the API
	credentialType, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	response, d := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.CredentialTypesApi.CreateCredentialType(ctx, plan.EnvironmentId.ValueString()).CredentialType(*credentialType).Execute()
		},
		"CreateCredentialType",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(response.(*credentials.CredentialType))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *CredentialTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *CredentialTypeResourceModel

	if r.client == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	ctx = context.WithValue(ctx, credentials.ContextServerVariables, map[string]string{
		"suffix": r.region.URLSuffix,
	})

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	response, diags := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.CredentialTypesApi.ReadOneCredentialType(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
		},
		"ReadOneCredentialType",
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

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(response.(*credentials.CredentialType))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CredentialTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state CredentialTypeResourceModel

	if r.client == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	ctx = context.WithValue(ctx, credentials.ContextServerVariables, map[string]string{
		"suffix": r.region.URLSuffix,
	})

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the model for the API
	credentialType, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	response, diags := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.CredentialTypesApi.UpdateCredentialType(ctx, plan.EnvironmentId.ValueString(), plan.Id.ValueString()).CredentialType(*credentialType).Execute()
		},
		"UpdateCredentialType",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(response.(*credentials.CredentialType))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *CredentialTypeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *CredentialTypeResourceModel

	if r.client == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	ctx = context.WithValue(ctx, credentials.ContextServerVariables, map[string]string{
		"suffix": r.region.URLSuffix,
	})

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	_, diags := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			r, err := r.client.CredentialTypesApi.DeleteCredentialType(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return nil, r, err
		},
		"DeleteCredentialType",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *CredentialTypeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splitLength := 2
	attributes := strings.SplitN(req.ID, "/", splitLength)

	if len(attributes) != splitLength {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("invalid id (\"%s\") specified, should be in format \"environment_id/credential_type_id/\"", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_id"), attributes[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), attributes[1])...)
}

func (p *CredentialTypeResourceModel) expand(ctx context.Context) (*credentials.CredentialType, diag.Diagnostics) {
	var diags diag.Diagnostics

	credentialTypeMetaData := credentials.NewCredentialTypeMetaData()

	// expand credential type metadata and metadata fields
	if !p.Metadata.IsNull() && !p.Metadata.IsUnknown() {
		var metadata MetadataModel
		d := p.Metadata.As(ctx, &metadata, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		credentialTypeMetaData, d = metadata.expandMetaDataModel(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}
	}

	data := credentials.NewCredentialType(p.CardDesignTemplate.ValueString(), *credentialTypeMetaData, p.Title.ValueString())
	data.SetDescription(p.Description.ValueString())
	data.SetCardType(p.CardType.ValueString())

	return data, diags
}

func (p *MetadataModel) expandMetaDataModel(ctx context.Context) (*credentials.CredentialTypeMetaData, diag.Diagnostics) {
	var diags diag.Diagnostics

	cardMetadata := credentials.NewCredentialTypeMetaDataWithDefaults()

	if !p.Name.IsNull() && !p.Name.IsUnknown() {
		cardMetadata.SetName(p.Name.ValueString())
	}

	if !p.BackgroundImage.IsNull() && !p.BackgroundImage.IsUnknown() {
		cardMetadata.SetBackgroundImage(p.BackgroundImage.ValueString())
	}

	if !p.BgOpacityPercent.IsNull() && !p.BgOpacityPercent.IsUnknown() {
		cardMetadata.SetBgOpacityPercent(int32(p.BgOpacityPercent.ValueInt64()))
	}

	if !p.CardColor.IsNull() && !p.CardColor.IsUnknown() {
		cardMetadata.SetCardColor(p.CardColor.ValueString())
	}

	if !p.Description.IsNull() && !p.Description.IsUnknown() {
		cardMetadata.SetDescription(p.Description.ValueString())
	}

	if !p.LogoImage.IsNull() && !p.LogoImage.IsUnknown() {
		cardMetadata.SetLogoImage(p.LogoImage.ValueString())
	}

	if !p.TextColor.IsNull() && !p.TextColor.IsUnknown() {
		cardMetadata.SetTextColor(p.TextColor.ValueString())
	}

	if !p.Version.IsNull() && !p.Version.IsUnknown() {
		cardMetadata.SetVersion(int32(p.Version.ValueInt64()))
	}

	// expand fields
	if !p.Fields.IsNull() && !p.Fields.IsUnknown() {
		var innerFields []FieldsModel
		fields := make([]credentials.CredentialTypeMetaDataFieldsInner, 0)
		diags.Append(p.Fields.ElementsAs(ctx, &innerFields, false)...)
		if diags.HasError() {
			return nil, diags
		}
		for _, v := range innerFields {
			field, d := v.expandFields(ctx)
			diags.Append(d...)
			if diags.HasError() {
				return nil, diags
			}

			fields = append(fields, *field)
		}
		// complete the meta data object
		cardMetadata.SetFields(fields)
	}

	return cardMetadata, diags
}

func (p *FieldsModel) expandFields(ctx context.Context) (*credentials.CredentialTypeMetaDataFieldsInner, diag.Diagnostics) {
	var diags diag.Diagnostics

	innerFields := credentials.NewCredentialTypeMetaDataFieldsInnerWithDefaults()

	attrType := credentials.EnumCredentialTypeMetaDataFieldsType(p.Type.ValueString())
	attrId := p.Type.ValueString() + " -> " + p.Title.ValueString() // construct id per API requirements

	if attrType == credentials.ENUMCREDENTIALTYPEMETADATAFIELDSTYPE_ALPHANUMERIC_TEXT {
		innerFields.SetValue(p.Value.ValueString()) // required if static text attribute - todo: need to test & error if not provided at schema?
	}

	if attrType == credentials.ENUMCREDENTIALTYPEMETADATAFIELDSTYPE_DIRECTORY_ATTRIBUTE {
		innerFields.SetAttribute(p.Attribute.ValueString()) // required if directory attribute - todo: need to test & error if not provided at schema?

		// todo: check if the attribute exists, if it doesn't error? or warn?
		// current APIs makes this... not simple to do
	}

	innerFields.SetId(attrId)
	innerFields.SetType(attrType)
	innerFields.SetTitle(p.Title.ValueString())
	innerFields.SetIsVisible(p.IsVisible.ValueBool())

	return innerFields, diags
}

func (p *CredentialTypeResourceModel) toState(apiObject *credentials.CredentialType) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)
		return diags
	}

	// credential attributes
	p.Id = framework.StringOkToTF(apiObject.GetIdOk())
	p.EnvironmentId = framework.StringToTF(apiObject.GetEnvironment().Id)
	p.Title = framework.StringOkToTF(apiObject.GetTitleOk())
	p.Description = framework.StringOkToTF(apiObject.GetDescriptionOk())
	p.CardType = framework.StringOkToTF(apiObject.GetCardTypeOk())
	p.CardDesignTemplate = framework.StringOkToTF(apiObject.GetCardDesignTemplateOk())

	// credential metadata
	metadata, d := toStateMetadata(apiObject.GetMetadataOk())
	diags.Append(d...)
	p.Metadata = metadata

	return diags
}

func toStateMetadata(metadata *credentials.CredentialTypeMetaData, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	// core metadata object
	metadataMap := map[string]attr.Value{
		"background_image":   framework.StringOkToTF(metadata.GetBackgroundImageOk()),
		"bg_opacity_percent": framework.Int32OkToTF(metadata.GetBgOpacityPercentOk()),
		"card_color":         framework.StringOkToTF(metadata.GetCardColorOk()),
		"description":        framework.StringOkToTF(metadata.GetDescriptionOk()),
		"text_color":         framework.StringOkToTF(metadata.GetTextColorOk()),
		"version":            framework.Int32OkToTF(metadata.GetVersionOk()),
		"logo_image":         framework.StringOkToTF(metadata.GetLogoImageOk()),
		"name":               framework.StringOkToTF(metadata.GetNameOk()),
	}

	// metadata fields objects
	fields, d := toStateFields(metadata.GetFieldsOk())
	diags.Append(d...)

	metadataMap["fields"] = fields
	flattenedObj, d := types.ObjectValue(metadataServiceTFObjectTypes, metadataMap)
	diags.Append(d...)

	return flattenedObj, diags
}

func toStateFields(innerFields []credentials.CredentialTypeMetaDataFieldsInner, ok bool) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	tfInnerObjType := types.ObjectType{AttrTypes: innerFieldsServiceTFObjectTypes}
	innerflattenedList := []attr.Value{}
	for _, v := range innerFields {

		fieldsMap := map[string]attr.Value{
			"id":         framework.StringOkToTF(v.GetIdOk()),
			"title":      framework.StringOkToTF(v.GetTitleOk()),
			"attribute":  framework.StringOkToTF(v.GetAttributeOk()),
			"value":      framework.StringOkToTF(v.GetValueOk()),
			"is_visible": framework.BoolOkToTF(v.GetIsVisibleOk()),
			"type":       enumCredentialTypeMetaDataFieldsOkToTF(v.GetTypeOk()),
		}
		innerflattenedObj, d := types.ObjectValue(innerFieldsServiceTFObjectTypes, fieldsMap)
		diags.Append(d...)

		innerflattenedList = append(innerflattenedList, innerflattenedObj)
	}
	fields, d := types.ListValue(tfInnerObjType, innerflattenedList)
	diags.Append(d...)

	return fields, diags
}

func enumCredentialTypeMetaDataFieldsOkToTF(v *credentials.EnumCredentialTypeMetaDataFieldsType, ok bool) basetypes.StringValue {
	if !ok || v == nil {
		return types.StringNull()
	} else {
		return types.StringValue(string(*v))
	}
}
