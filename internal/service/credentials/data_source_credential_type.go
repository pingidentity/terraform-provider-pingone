package credentials

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/credentials"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
)

// Types
type CredentialTypeDataSource serviceClientType

type CredentialTypeDataSourceModel struct {
	Id                 types.String `tfsdk:"id"`
	EnvironmentId      types.String `tfsdk:"environment_id"`
	CredentialTypeId   types.String `tfsdk:"credential_type_id"`
	Title              types.String `tfsdk:"title"`
	Description        types.String `tfsdk:"description"`
	CardType           types.String `tfsdk:"card_type"`
	CardDesignTemplate types.String `tfsdk:"card_design_template"`
	Metadata           types.Object `tfsdk:"metadata"`
}

type MetadataDataSourceModel struct {
	BackgroundImage  types.String `tfsdk:"background_image"`
	BgOpacityPercent types.Int64  `tfsdk:"bg_opacity_percent"`
	CardColor        types.String `tfsdk:"card_color"`
	Columns          types.Int64  `tfsdk:"columns"`
	Description      types.String `tfsdk:"description"`
	TextColor        types.String `tfsdk:"text_color"`
	Version          types.Int64  `tfsdk:"version"`
	LogoImage        types.String `tfsdk:"logo_image"`
	Name             types.String `tfsdk:"name"`
	Fields           types.List   `tfsdk:"fields"`
}

type FieldsDataSourceModel struct {
	Id        types.String `tfsdk:"id"`
	Type      types.String `tfsdk:"type"`
	Title     types.String `tfsdk:"title"`
	IsVisible types.Bool   `tfsdk:"is_visible"`
	Attribute types.String `tfsdk:"attribute"`
	Value     types.String `tfsdk:"value"`
}

var (
	metadataDataSourceServiceTFObjectTypes = map[string]attr.Type{
		"background_image":   types.StringType,
		"bg_opacity_percent": types.Int64Type,
		"card_color":         types.StringType,
		"columns":            types.Int64Type,
		"description":        types.StringType,
		"text_color":         types.StringType,
		"version":            types.Int64Type,
		"logo_image":         types.StringType,
		"name":               types.StringType,
		"fields":             types.ListType{ElemType: types.ObjectType{AttrTypes: innerFieldsDataSourceServiceTFObjectTypes}},
	}

	innerFieldsDataSourceServiceTFObjectTypes = map[string]attr.Type{
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
	_ datasource.DataSource = &CredentialTypeDataSource{}
)

// New Object
func NewCredentialTypeDataSource() datasource.DataSource {
	return &CredentialTypeDataSource{}
}

// Metadata
func (r *CredentialTypeDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_credential_type"
}

func (r *CredentialTypeDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Datasource to retrieve a PingOne Credentials credential type by its Credential Type Id.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("PingOne environment identifier (UUID) in which the credential type exists."),
			),

			"credential_type_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("Identifier (UUID) associated with the credential type."),
			),

			"title": schema.StringAttribute{
				Description: "Title of the credential.",
				Computed:    true,
			},

			"description": schema.StringAttribute{
				Description: "A description of the credential type.",
				Computed:    true,
			},

			"card_type": schema.StringAttribute{
				Description: "A descriptor of the credential type. Can be non-identity types such as proof of employment or proof of insurance.",
				Computed:    true,
			},

			"card_design_template": schema.StringAttribute{
				Description: "An SVG formatted image containing placeholders for the credentials fields that need to be displayed in the image.",
				Computed:    true,
			},

			"metadata": schema.SingleNestedAttribute{
				Description: "An object that contains the names, data types, and other metadata related to the credentia",
				Computed:    true,

				Attributes: map[string]schema.Attribute{
					"background_image": schema.StringAttribute{
						Description: "A base64 encoded image of the background to show in the credential.",
						Computed:    true,
					},

					"bg_opacity_percent": schema.Int64Attribute{
						Description: "Percent opacity of the background image in the credential.",
						Computed:    true,
					},

					"card_color": schema.StringAttribute{
						Description: "Color to show on the credential.",
						Computed:    true,
					},

					"columns": schema.Int64Attribute{
						Description: "Number of columns to organize the fields displayed on the credential.",
						Computed:    true,
					},

					"description": schema.StringAttribute{
						Description: "Description of the credential.",
						Computed:    true,
					},

					"logo_image": schema.StringAttribute{
						Description: "A base64 encoded image of the logo to show in the credential.",
						Computed:    true,
					},

					"name": schema.StringAttribute{
						Description: "Name of the credential.",
						Computed:    true,
					},

					"text_color": schema.StringAttribute{
						Description: "Color of the text to show on the credential.",
						Computed:    true},

					"version": schema.Int64Attribute{
						Description: "Version of this credential.",
						Computed:    true,
					},

					"fields": schema.ListNestedAttribute{
						Description: "Array of objects representing the credential fields.",
						Computed:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Description: "Identifier of the field formatted as “<fields.type> -> <fields.title>.",
									Computed:    true,
								},
								"type": schema.StringAttribute{
									Description: "Type of data in the field.",
									Computed:    true,
								},
								"title": schema.StringAttribute{
									Description: "Descriptive text when showing the field.",
									Computed:    true,
								},
								"attribute": schema.StringAttribute{
									Description: "Name of the PingOne Directory attribute. Present if field.type is Directory Attribute.",
									Computed:    true,
								},
								"value": schema.StringAttribute{
									Description: "The text to appear on the credential for a field.type of Alphanumeric Text.",
									Computed:    true,
								},
								"is_visible": schema.BoolAttribute{
									Description: "Specifies whether the field should be visible to viewers of the credential.",
									Computed:    true,
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r *CredentialTypeDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.Client = preparedClient
}

func (r *CredentialTypeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *CredentialTypeDataSourceModel

	if r.Client == nil {
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
	var response *credentials.CredentialType
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return r.Client.CredentialTypesApi.ReadOneCredentialType(ctx, data.EnvironmentId.ValueString(), data.CredentialTypeId.ValueString()).Execute()
		},
		"ReadOneCredentialType",
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

func (p *CredentialTypeDataSourceModel) toState(apiObject *credentials.CredentialType) diag.Diagnostics {
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
	p.EnvironmentId = framework.StringToTF(*apiObject.GetEnvironment().Id)
	p.CredentialTypeId = framework.StringToTF(apiObject.GetId())
	p.Title = framework.StringOkToTF(apiObject.GetTitleOk())
	p.Description = framework.StringOkToTF(apiObject.GetDescriptionOk())
	p.CardType = framework.StringOkToTF(apiObject.GetCardTypeOk())
	p.CardDesignTemplate = framework.StringOkToTF(apiObject.GetCardDesignTemplateOk())

	// credential metadata
	metadata, d := toStateMetadataDataSource(apiObject.GetMetadataOk())
	diags.Append(d...)
	p.Metadata = metadata

	return diags
}

func toStateMetadataDataSource(metadata *credentials.CredentialTypeMetaData, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	// core metadata object
	metadataMap := map[string]attr.Value{
		"background_image":   framework.StringOkToTF(metadata.GetBackgroundImageOk()),
		"bg_opacity_percent": framework.Int32OkToTF(metadata.GetBgOpacityPercentOk()),
		"card_color":         framework.StringOkToTF(metadata.GetCardColorOk()),
		"columns":            framework.Int32OkToTF(metadata.GetColumnsOk()),
		"description":        framework.StringOkToTF(metadata.GetDescriptionOk()),
		"text_color":         framework.StringOkToTF(metadata.GetTextColorOk()),
		"version":            framework.Int32OkToTF(metadata.GetVersionOk()),
		"logo_image":         framework.StringOkToTF(metadata.GetLogoImageOk()),
		"name":               framework.StringOkToTF(metadata.GetNameOk()),
	}

	// metadata fields objects
	fields, d := toStateFieldsDataSource(metadata.GetFieldsOk())
	diags.Append(d...)

	metadataMap["fields"] = fields
	flattenedObj, d := types.ObjectValue(metadataDataSourceServiceTFObjectTypes, metadataMap)
	diags.Append(d...)

	return flattenedObj, diags
}

func toStateFieldsDataSource(innerFields []credentials.CredentialTypeMetaDataFieldsInner, ok bool) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	tfInnerObjType := types.ObjectType{AttrTypes: innerFieldsDataSourceServiceTFObjectTypes}
	innerflattenedList := []attr.Value{}
	for _, v := range innerFields {

		fieldsMap := map[string]attr.Value{
			"id":         framework.StringOkToTF(v.GetIdOk()),
			"title":      framework.StringOkToTF(v.GetTitleOk()),
			"attribute":  framework.StringOkToTF(v.GetAttributeOk()),
			"value":      framework.StringOkToTF(v.GetValueOk()),
			"is_visible": framework.BoolOkToTF(v.GetIsVisibleOk()),
			"type":       framework.EnumOkToTF(v.GetTypeOk()),
		}
		innerflattenedObj, d := types.ObjectValue(innerFieldsDataSourceServiceTFObjectTypes, fieldsMap)
		diags.Append(d...)

		innerflattenedList = append(innerflattenedList, innerflattenedObj)
	}
	fields, d := types.ListValue(tfInnerObjType, innerflattenedList)
	diags.Append(d...)

	return fields, diags
}
