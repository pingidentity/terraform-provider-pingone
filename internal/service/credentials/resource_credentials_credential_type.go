package credentials

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
	Metadata           types.List   `tfsdk:"metadata"`
}

type MetadataModel struct {
	BackgroundImage  types.String `tfsdk:"background_image"`
	BgOpacityPercent types.Int64  `tfsdk:"bg_opacity_percent"`
	CardColor        types.String `tfsdk:"card_color"`
	Description      types.String `tfsdk:"description"`
	TextColor        types.String `tfsdk:"text_color"`
	Version          types.Int64  `tfsdk:"version"` // Watch Item - Best practice is to allow service to set, but if version of 5 or higher is not provided, creds do not appear in UI
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
	resp.TypeName = req.ProviderTypeName + "_credentials_credential_type"
}

func (r *CredentialTypeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage PingOne Credentials credential types.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(framework.SchemaDescription{
				Description: "The ID of the environment to create the credential type in."},
			),

			"title": schema.StringAttribute{
				Description: "card title",
				Required:    true,
			},

			"description": schema.StringAttribute{
				MarkdownDescription: "test",
				Optional:            true,
			},

			"card_type": schema.StringAttribute{
				MarkdownDescription: "test",
				Optional:            true,
			},

			"card_design_template": schema.StringAttribute{
				MarkdownDescription: "test",
				Required:            true,
			},
		},
		Blocks: map[string]schema.Block{
			"metadata": schema.ListNestedBlock{
				Description:         "",
				MarkdownDescription: "",
				NestedObject: schema.NestedBlockObject{

					Attributes: map[string]schema.Attribute{
						"background_image": schema.StringAttribute{
							Description:         "",
							MarkdownDescription: "",
							Optional:            true,
						},
						"bg_opacity_percent": schema.Int64Attribute{
							Description: "A numnber containing the percent opacity of the background image in the credential. High percentage opacity may make displayed text difficult to read.",
							Optional:    true,
							Validators: []validator.Int64{
								int64validator.Between(0, 100),
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
						},
						"logo_image": schema.StringAttribute{
							Description:         "",
							MarkdownDescription: "",
							Optional:            true,
						},
						"name": schema.StringAttribute{
							Description:         "",
							MarkdownDescription: "",
							Optional:            true,
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
								int64validator.AtLeast(5),
							},
						},
					},
					Blocks: map[string]schema.Block{
						"fields": schema.ListNestedBlock{
							Description:         "",
							MarkdownDescription: "",

							NestedObject: schema.NestedBlockObject{

								Attributes: map[string]schema.Attribute{
									// Placeholder. The id value is constructed per specific API requirements. Future may allow user-provided id.
									"id": schema.StringAttribute{
										Optional: true,
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
									},
									"attribute": schema.StringAttribute{
										Description:         "",
										MarkdownDescription: "",
										Optional:            true,
									},
									"value": schema.StringAttribute{
										Description:         "",
										MarkdownDescription: "",
										Optional:            true,
									},
									"is_visible": schema.BoolAttribute{
										Description:         "",
										MarkdownDescription: "",
										Optional:            true,
									},
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
	splitLength := 3
	attributes := strings.SplitN(req.ID, "/", splitLength)

	if len(attributes) != splitLength {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("invalid id (\"%s\") specified, should be in format \"environment_id/credentials_credential_type_id/\"", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_id"), attributes[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), attributes[2])...)
}

func (p *CredentialTypeResourceModel) expand(ctx context.Context) (*credentials.CredentialType, diag.Diagnostics) {
	var diags diag.Diagnostics

	var metadata []MetadataModel
	diags.Append(p.Metadata.ElementsAs(ctx, &metadata, false)...)
	if diags.HasError() {
		return nil, diags
	}

	// TODO: need expand metadata & fields
	cardMetadata := credentials.NewCredentialTypeMetaDataWithDefaults()
	for _, v := range metadata {
		if !v.Fields.IsNull() && !v.Fields.IsUnknown() {
			var innerFields []FieldsModel
			fields := make([]credentials.CredentialTypeMetaDataFieldsInner, 0)
			diags.Append(v.Fields.ElementsAs(ctx, &innerFields, false)...)
			if diags.HasError() {
				return nil, diags
			}
			for _, i := range innerFields {
				field := *credentials.NewCredentialTypeMetaDataFieldsInnerWithDefaults()

				attrType := credentials.EnumCredentialTypeMetaDataFieldsType(i.Type.ValueString())
				attrId := i.Type.ValueString() + " -> " + i.Title.ValueString() // construct id per API requirements

				if attrType == credentials.ENUMCREDENTIALTYPEMETADATAFIELDSTYPE_ALPHANUMERIC_TEXT {
					field.SetValue(i.Value.ValueString()) // required if static text attribute - todo: need to test & error if not provided
				}

				if attrType == credentials.ENUMCREDENTIALTYPEMETADATAFIELDSTYPE_DIRECTORY_ATTRIBUTE {
					field.SetAttribute(i.Attribute.ValueString()) // required if directory attribute - todo: need to test & error if not provided

					// todo: check if the attribute exists, if it doesn't error? or warn?
				}

				field.SetId(attrId)
				field.SetType(attrType)
				field.SetTitle(i.Title.ValueString())
				field.SetIsVisible(i.IsVisible.ValueBool())

				fields = append(fields, field)
			}
			// complete the meta data object
			cardMetadata.SetFields(fields)
		}

		if !v.Name.IsNull() && !v.Name.IsUnknown() {
			cardMetadata.SetName(v.Name.ValueString())
		}

		if !v.BackgroundImage.IsNull() && !v.BackgroundImage.IsUnknown() {
			cardMetadata.SetBackgroundImage(v.BackgroundImage.ValueString())
		}

		if !v.BgOpacityPercent.IsNull() && !v.BgOpacityPercent.IsUnknown() {
			cardMetadata.SetBgOpacityPercent(int32(v.BgOpacityPercent.ValueInt64()))
		}

		if !v.CardColor.IsNull() && !v.CardColor.IsUnknown() {
			cardMetadata.SetCardColor(v.CardColor.ValueString())
		}

		if !v.LogoImage.IsNull() && !v.LogoImage.IsUnknown() {
			cardMetadata.SetLogoImage(v.LogoImage.ValueString())
		}

		if !v.TextColor.IsNull() && !v.TextColor.IsUnknown() {
			cardMetadata.SetTextColor(v.TextColor.ValueString())
		}

		if !v.Version.IsNull() && !v.Version.IsUnknown() {
			cardMetadata.SetVersion(int32(v.Version.ValueInt64()))
		}
	}

	data := credentials.NewCredentialType(p.CardDesignTemplate.ValueString(), *cardMetadata, p.Title.ValueString())

	data.SetDescription(p.Description.ValueString())
	data.SetCardType(p.CardType.ValueString())

	return data, diags
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

	p.Id = framework.StringToTF(apiObject.GetId())
	p.EnvironmentId = framework.StringToTF(apiObject.GetEnvironment().Id)
	p.Title = framework.StringToTF(apiObject.GetTitle())
	p.Description = framework.StringToTF((apiObject.GetDescription()))
	p.CardType = framework.StringToTF(apiObject.GetCardType())
	p.CardDesignTemplate = framework.StringToTF(apiObject.GetCardDesignTemplate())

	// TODO metadata & fields handling...

	return diags
}
