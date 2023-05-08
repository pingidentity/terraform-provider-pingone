package credentials

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/patrickcping/pingone-go-sdk-v2/credentials"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type CredentialTypeDataSource struct {
	client *credentials.APIClient
	region model.RegionMapping
}

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
	Description      types.String `tfsdk:"description"`
	TextColor        types.String `tfsdk:"text_color"`
	Version          types.Int64  `tfsdk:"version"` // Watch Item - Best practice is to allow service to set, but if version of 5 or higher is not provided, creds do not appear in UI!
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
	resp.TypeName = req.ProviderTypeName + "_credentials_credential_type"
}

func (r *CredentialTypeDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Data to retrieve a PingOne Credentials credential type by its Credential Type Id. The credential_type_id is the only parameter to ensure a single credential is retrieved.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(framework.SchemaDescription{
				Description: "The ID of the environment to create the credential type in."},
			),

			"credential_type_id": schema.StringAttribute{
				Description: "The ID of the credential type.",
				Optional:    true,
				Validators: []validator.String{
					verify.P1ResourceIDValidator(),
				},
			},

			"title": schema.StringAttribute{
				Description: "A string that specifies the title of the credential. Verification sites are expected to be able to request the issued credential from the compatible wallet app using the credential title.",
				Computed:    true,
			},

			"description": schema.StringAttribute{
				Description: "",
				Computed:    true,
			},

			"card_type": schema.StringAttribute{
				Description: "",
				Computed:    true,
			},

			"card_design_template": schema.StringAttribute{
				Description: "A string that specifies an SVG formatted image containing placeholders for the credential fields that need to be displayed in the image.",
				Computed:    true,
			},

			"metadata": schema.SingleNestedAttribute{
				Description:         "",
				MarkdownDescription: "",
				Computed:            true,

				Attributes: map[string]schema.Attribute{
					"background_image": schema.StringAttribute{
						Description:         "",
						MarkdownDescription: "",
						Computed:            true,
					},

					"bg_opacity_percent": schema.Int64Attribute{
						Description: "A numnber containing the percent opacity of the background image in the credential. High percentage opacity may make displayed text difficult to read.",
						Computed:    true,
					},

					"card_color": schema.StringAttribute{
						Description: "A string containing a 6-digit hexadecimal color code specifying the color of the credential.",
						Computed:    true,
					},

					"description": schema.StringAttribute{
						Description:         "",
						MarkdownDescription: "",
						Computed:            true,
					},

					"logo_image": schema.StringAttribute{
						Description:         "",
						MarkdownDescription: "",
						Computed:            true,
					},

					"name": schema.StringAttribute{
						Description:         "",
						MarkdownDescription: "",
						Computed:            true,
					},

					"text_color": schema.StringAttribute{
						Description: "A string containing a 6-digit hexadecimal color code specifying the color of the credential text.",
						Computed:    true},

					"version": schema.Int64Attribute{
						Description:         "",
						MarkdownDescription: "",
						Computed:            true,
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
									Computed:            true,
								},
								"title": schema.StringAttribute{
									Description:         "",
									MarkdownDescription: "",
									Computed:            true,
								},
								"attribute": schema.StringAttribute{
									Description:         "",
									MarkdownDescription: "",
									Computed:            true,
								},
								"value": schema.StringAttribute{
									Description:         "",
									MarkdownDescription: "",
									Computed:            true,
								},
								"is_visible": schema.BoolAttribute{
									Description:         "",
									MarkdownDescription: "",
									Computed:            true,
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

	r.client = preparedClient
	r.region = resourceConfig.Client.API.Region
}

func (r *CredentialTypeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *CredentialTypeDataSourceModel

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
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	response, diags := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.CredentialTypesApi.ReadOneCredentialType(ctx, data.EnvironmentId.ValueString(), data.CredentialTypeId.ValueString()).Execute()
		},
		"ReadOneCredentialType",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Error if not found
	if response == nil {
		resp.Diagnostics.AddError(
			"Cannot find credential type",
			fmt.Sprintf("The credential type ID %s for environment %s cannot be found.", data.CredentialTypeId.String(), data.EnvironmentId.String()),
		)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(response.(*credentials.CredentialType))...)
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
	p.EnvironmentId = framework.StringToTF(apiObject.GetEnvironment().Id)
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
			"type":       enumCredentialTypeMetaDataFieldsDataSourceOkToTF(v.GetTypeOk()),
		}
		innerflattenedObj, d := types.ObjectValue(innerFieldsDataSourceServiceTFObjectTypes, fieldsMap)
		diags.Append(d...)

		innerflattenedList = append(innerflattenedList, innerflattenedObj)
	}
	fields, d := types.ListValue(tfInnerObjType, innerflattenedList)
	diags.Append(d...)

	return fields, diags
}

func enumCredentialTypeMetaDataFieldsDataSourceOkToTF(v *credentials.EnumCredentialTypeMetaDataFieldsType, ok bool) basetypes.StringValue {
	if !ok || v == nil {
		return types.StringNull()
	} else {
		return types.StringValue(string(*v))
	}
}
