package credentials

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/patrickcping/pingone-go-sdk-v2/credentials"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	customstringvalidator "github.com/pingidentity/terraform-provider-pingone/internal/framework/stringvalidator"
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
	Columns          types.Int64  `tfsdk:"columns"`
	Description      types.String `tfsdk:"description"`
	TextColor        types.String `tfsdk:"text_color"`
	Version          types.Int64  `tfsdk:"version"`
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
		"columns":            types.Int64Type,
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
	const attrMinColumns = 1
	const attrMaxColumns = 3
	const attrDefaultVersion = 5
	const attrMinPercent = 0
	const attrMaxPercent = 100
	const imageMaxSize = 50000

	titleDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Title of the credential. Verification sites are expected to be able to request the issued credential from the compatible wallet app using the title.  This value aligns to `${cardTitle}` in the `card_design_template`.",
	)

	credentialDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A description of the credential type. This value aligns to `${cardSubtitle}` in the `card_design_template`.",
	)

	fieldsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"In a credential, the information is stored as key-value pairs where `fields` defines those key-value pairs. Effectively, `fields.title` is the key and its value is `fields.value` or extracted from the PingOne Directory attribute named in `fields.attribute`.",
	)

	fieldsIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Identifier of the field formatted as `<fields.type> -> <fields.title>`.",
	)

	fieldsTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Type of data in the credential field. The must contain one of the following types: `Directory Attribute`, `Alphanumeric Text`, or `Issued Timestamp`.",
	)

	fieldsAttributeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Name of the PingOne Directory attribute. Present if `field.type` is `Directory Attribute`.",
	)

	fieldsValueDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The text to appear on the credential for a `field.type` of `Alphanumeric Text`.",
	)

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage the credential types used by compatible wallet applications.\n\n" +
			framework.SchemaAttributeDescriptionFromMarkdown("~> You must ensure that any fields used in the `card_design_template` are defined appropriately in `metadata.fields` or errors occur when you attempt to create a credential of that type.").MarkdownDescription,

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("PingOne environment identifier (UUID) in which the credential type exists."),
			),

			"title": schema.StringAttribute{
				Description:         titleDescription.Description,
				MarkdownDescription: titleDescription.MarkdownDescription,
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
					customstringvalidator.IsRequiredIfRegexMatchesPathValue(
						regexp.MustCompile(`\${cardTitle}`),
						framework.SchemaAttributeDescriptionFromMarkdown("The title argument is required because the ${cardTitle} element is defined in the `card_design_template`.").MarkdownDescription,
						path.MatchRoot("card_design_template"),
					),
					customstringvalidator.RegexMatchesPathValue(
						regexp.MustCompile(`\${cardTitle}`),
						"The title argument is defined but the card_design_template does not have a ${cardTitle} element.",
						path.MatchRoot("card_design_template"),
					),
				},
			},

			"description": schema.StringAttribute{
				Description:         credentialDescription.Description,
				MarkdownDescription: credentialDescription.MarkdownDescription,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
					customstringvalidator.IsRequiredIfRegexMatchesPathValue(
						regexp.MustCompile(`\${cardSubtitle}`),
						framework.SchemaAttributeDescriptionFromMarkdown("The description argument is required because the ${cardSubtitle} element is defined in the `card_design_template`.").MarkdownDescription,
						path.MatchRoot("card_design_template"),
					),
					customstringvalidator.RegexMatchesPathValue(
						regexp.MustCompile(`\${cardSubtitle}`),
						"The description argument is defined but the card_design_template does not have a ${cardSubtitle} element.",
						path.MatchRoot("card_design_template"),
					),
				},
			},

			"card_type": schema.StringAttribute{
				Description: "A descriptor of the credential type. Can be non-identity types such as proof of employment or proof of insurance.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			"card_design_template": schema.StringAttribute{
				Description: "An SVG formatted image containing placeholders for the credentials fields that need to be displayed in the image.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(`^<svg.*>[\s\S]*<\/svg>\s*$`), "expected value to contain a valid PingOne Credentials SVG card template."),
				},
			},

			"metadata": schema.SingleNestedAttribute{
				Description: "Contains the names, data types, and other metadata related to the credential.",
				Required:    true,

				Attributes: map[string]schema.Attribute{
					"background_image": schema.StringAttribute{
						Description: "A base64 encoded image of the background to show in the credential. The value must include a Content-type prefix, such as data:image/png;base64.",
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.LengthAtMost(imageMaxSize),
							// Required until P1Creds follows the standard PingOne image handling capability.
							// Attempts of other stop-gap mechanisms to detect and update Content-Type yielded inconsistent results.
							stringvalidator.RegexMatches(regexp.MustCompile(`^data:image\/(\w+);base64,`), "base64encoded image must include Content-type prefix, such as data:image/jpeg;base64, data:image/svg;base64, or data:image/png;base64."),
							customstringvalidator.IsBase64Encoded(),
							customstringvalidator.IsRequiredIfRegexMatchesPathValue(
								regexp.MustCompile(`\${backgroundImage}`),
								"The metadata.background_image argument is required because the ${backgroundImage} element is defined in the card_design_template.",
								path.MatchRoot("card_design_template"),
							),
							customstringvalidator.RegexMatchesPathValue(
								regexp.MustCompile(`\${backgroundImage}`),
								"The metadata.background_image argument is defined but the card_design_template does not have a ${backgroundImage} element.",
								path.MatchRoot("card_design_template"),
							),
						},
					},

					"bg_opacity_percent": schema.Int64Attribute{
						Description: "A numnber indicating the percent opacity of the background image in the credential. High percentage opacity may make text on the credential difficult to read.",
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
							customstringvalidator.IsRequiredIfRegexMatchesPathValue(
								regexp.MustCompile(`\${cardColor}`),
								"The metadata.card_color argument is required because the ${cardColor} element is defined in the card_design_template.",
								path.MatchRoot("card_design_template"),
							),
							customstringvalidator.RegexMatchesPathValue(
								regexp.MustCompile(`\${cardColor}`),
								"The metadata.card_color argument is defined but the card_design_template does not have a ${cardColor} element.",
								path.MatchRoot("card_design_template"),
							),
						},
					},

					"columns": schema.Int64Attribute{
						Description: "Indicates a number (between 1-3) of columns to display visible fields on the credential.",
						Optional:    true,
						Validators: []validator.Int64{
							int64validator.Between(attrMinColumns, attrMaxColumns),
						},
					},

					"description": schema.StringAttribute{
						Description: "Description of the credential.",
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(attrMinLength),
						},
					},

					"logo_image": schema.StringAttribute{
						Description: "A base64 encoded image of the logo to show in the credential. The value must include a Content-type prefix, such as data:image/png;base64.",
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.LengthAtMost(imageMaxSize),
							// Required until P1Creds follows the standard PingOne image handling capability.
							// Attempts of other stop-gap mechanisms to detect and update Content-Type yielded inconsistent results.
							stringvalidator.RegexMatches(regexp.MustCompile(`^data:image\/(\w+);base64,`), "base64encoded image must include Content-type prefix, such as data:image/jpeg;base64, data:image/svg;base64, or data:image/png;base64."),
							customstringvalidator.IsBase64Encoded(),
							customstringvalidator.IsRequiredIfRegexMatchesPathValue(
								regexp.MustCompile(`\${logoImage}`),
								"The metadata.card_color argument is required because the ${logoImage} element is defined in the card_design_template.",
								path.MatchRoot("card_design_template"),
							),
							customstringvalidator.RegexMatchesPathValue(
								regexp.MustCompile(`\${logoImage}`),
								"The metadata.logo_image argument is defined but the card_design_template does not have a ${logoImage} element.",
								path.MatchRoot("card_design_template"),
							),
						},
					},

					"name": schema.StringAttribute{
						Description: "Name of the credential.",
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(attrMinLength),
						},
					},

					"text_color": schema.StringAttribute{
						Description: "A string containing a 6-digit hexadecimal color code specifying the color of the credential text.",
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.RegexMatches(
								regexp.MustCompile(
									`^#([A-Fa-f0-9]{6})$`),
								"expected value to contain a valid 6-digit hexadecimal color code, prefixed with a hash (#) symbol."),
							customstringvalidator.IsRequiredIfRegexMatchesPathValue(
								regexp.MustCompile(`\${textColor}`),
								"The metadata.text_color argument is required because the ${textColor} element is defined in the card_design_template.",
								path.MatchRoot("card_design_template"),
							),
							customstringvalidator.RegexMatchesPathValue(
								regexp.MustCompile(`\${textColor}`),
								"The metadata.text_color argument is defined but the card_design_template does not have a ${textColor} element.",
								path.MatchRoot("card_design_template"),
							),
						},
					},
					//fix
					"version": schema.Int64Attribute{
						Description: "Number version of this credential.",
						Computed:    true,
						Default:     int64default.StaticInt64(attrDefaultVersion),
						// P1Creds has a limitation within the EarlyRelease.
						// To resolve, we will compute the "version" argument with a value of "5";
						// the same value set by the P1 admin console until resolved.
						// Below are the actual settings to use once fixed.
						// Optional: true,
						// Validators: []validator.Int64{
						// int64validator.AtLeast(attrMinVersion),
						//},
					},

					"fields": schema.ListNestedAttribute{
						Description:         fieldsDescription.Description,
						MarkdownDescription: fieldsDescription.MarkdownDescription,
						Required:            true,
						Validators: []validator.List{
							listvalidator.SizeAtLeast(attrMinLength),
						},
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Description:         fieldsIdDescription.Description,
									MarkdownDescription: fieldsIdDescription.MarkdownDescription,
									Computed:            true,
								},
								"type": schema.StringAttribute{
									Description:         fieldsTypeDescription.Description,
									MarkdownDescription: fieldsTypeDescription.MarkdownDescription,
									Required:            true,
									Validators: []validator.String{
										stringvalidator.OneOf(
											string(credentials.ENUMCREDENTIALTYPEMETADATAFIELDSTYPE_ALPHANUMERIC_TEXT),
											string(credentials.ENUMCREDENTIALTYPEMETADATAFIELDSTYPE_DIRECTORY_ATTRIBUTE),
											string(credentials.ENUMCREDENTIALTYPEMETADATAFIELDSTYPE_ISSUED_TIMESTAMP)),
									},
								},
								"title": schema.StringAttribute{
									Description: "Descriptive text when showing the field.",
									Optional:    true,
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(attrMinLength),
									},
								},
								"attribute": schema.StringAttribute{
									Description:         fieldsAttributeDescription.Description,
									MarkdownDescription: fieldsAttributeDescription.MarkdownDescription,
									Optional:            true,
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(attrMinLength),
										stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("value")),
										customstringvalidator.IsRequiredIfMatchesPathValue(basetypes.NewStringValue(string(credentials.ENUMCREDENTIALTYPEMETADATAFIELDSTYPE_DIRECTORY_ATTRIBUTE)), path.MatchRelative().AtParent().AtName("type")),
									},
								},
								"value": schema.StringAttribute{
									Description:         fieldsValueDescription.Description,
									MarkdownDescription: fieldsValueDescription.MarkdownDescription,
									Optional:            true,
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(attrMinLength),
										stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("attribute")),
										customstringvalidator.IsRequiredIfMatchesPathValue(basetypes.NewStringValue(string(credentials.ENUMCREDENTIALTYPEMETADATAFIELDSTYPE_ALPHANUMERIC_TEXT)), path.MatchRelative().AtParent().AtName("type")),
									},
								},
								"is_visible": schema.BoolAttribute{
									Description: "Specifies whether the field should be visible to viewers of the credential.",
									Optional:    true,
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
	timeoutValue := 15
	response, d := framework.ParseResponseWithCustomTimeout(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.CredentialTypesApi.CreateCredentialType(ctx, plan.EnvironmentId.ValueString()).CredentialType(*credentialType).Execute()
		},
		"CreateCredentialType",
		framework.DefaultCustomError,
		credentialTypeRetryConditions,
		time.Duration(timeoutValue)*time.Minute, // 15 mins
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

	// expand credential type metadata and metadata.fields
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

	if !p.Columns.IsNull() && !p.Columns.IsUnknown() {
		cardMetadata.SetColumns(int32(p.Columns.ValueInt64()))
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
			field, d := v.expandFields()
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

func (p *FieldsModel) expandFields() (*credentials.CredentialTypeMetaDataFieldsInner, diag.Diagnostics) {
	var diags diag.Diagnostics

	innerFields := credentials.NewCredentialTypeMetaDataFieldsInnerWithDefaults()

	attrType := credentials.EnumCredentialTypeMetaDataFieldsType(p.Type.ValueString())
	attrId := p.Type.ValueString() + " -> " + p.Title.ValueString() // construct id per P1Creds API recommendations

	if attrType == credentials.ENUMCREDENTIALTYPEMETADATAFIELDSTYPE_ALPHANUMERIC_TEXT {
		innerFields.SetValue(p.Value.ValueString())
	}

	if attrType == credentials.ENUMCREDENTIALTYPEMETADATAFIELDSTYPE_DIRECTORY_ATTRIBUTE {
		innerFields.SetAttribute(p.Attribute.ValueString())
	}

	innerFields.SetId(attrId)
	innerFields.SetType(attrType)
	innerFields.SetTitle(p.Title.ValueString())
	innerFields.SetIsVisible(p.IsVisible.ValueBool())

	if innerFields == nil {
		diags.AddWarning(
			"Unexpected Value",
			"Metadata.Fields object was unexpectedly null on expansion.  Please report this to the provider maintainers.",
		)
	}

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
	p.EnvironmentId = framework.StringToTF(*apiObject.GetEnvironment().Id)
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
		"columns":            framework.Int32OkToTF(metadata.GetColumnsOk()),
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

func credentialTypeRetryConditions(ctx context.Context, r *http.Response, p1error *model.P1Error) bool {

	var err error

	if p1error != nil {

		// Credential Issuer Profile's keys may not have propagated yet. Rare, but possible.
		if details, ok := p1error.GetDetailsOk(); ok && details != nil && len(details) > 0 {

			// detected issuer profile not fully deployed yet
			if m, err := regexp.MatchString("^issuerProfile must exist before creating credentialTypes", details[0].GetMessage()); err == nil && m {
				tflog.Warn(ctx, fmt.Sprintf("IssuerProfile (prerequisite) has not finished provisioning - %s.  Retrying...", details[0].GetMessage()))
				return true
			}
			if err != nil {
				tflog.Warn(ctx, "Cannot match error string for retry")
				return false
			}
		}

		// detected credentials service not fully deployed yet
		if m, _ := regexp.MatchString("^The actor attempting to perform the request is not authorized.", p1error.GetMessage()); err == nil && m {
			tflog.Warn(ctx, "Insufficient PingOne privileges detected")
			return true
		}
		if err != nil {
			tflog.Warn(ctx, "Cannot match error string for retry")
			return false
		}
	}

	return false
}
