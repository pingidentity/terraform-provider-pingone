package credentials

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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
							Description:         "",
							MarkdownDescription: "",
							Optional:            true,
						},
						"card_color": schema.StringAttribute{
							Description:         "",
							MarkdownDescription: "",
							Optional:            true,
						},
						"description": schema.StringAttribute{
							Description:         "",
							MarkdownDescription: "",
							Optional:            true,
						},
						"text_color": schema.StringAttribute{
							Description:         "",
							MarkdownDescription: "",
							Optional:            true,
						},
						"version": schema.Int64Attribute{
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
					},
					Blocks: map[string]schema.Block{
						"fields": schema.ListNestedBlock{
							Description:         "",
							MarkdownDescription: "",

							NestedObject: schema.NestedBlockObject{

								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description:         "",
										MarkdownDescription: "",
										Optional:            true,
									},
									"type": schema.StringAttribute{
										Description:         "",
										MarkdownDescription: "",
										Optional:            true,
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

	/*management.ReadOneApplicationRequest(ctx, )
	if CredentialType.GetApplication().Id{
			// make sure it exists

		}

	    t.GetOk("oidc_options"); ok {
		    var application *management.ApplicationOIDC
		    application, diags = expandApplicationOIDC(d)
		    if diags.HasError() {
		        return diags
		    }
		    applicationRequest.ApplicationOIDC = application
		} */

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
				field.SetId(i.Attribute.ValueString())
				field.SetTitle(i.Title.ValueString())
				field.SetIsVisible(i.IsVisible.ValueBool())
				field.SetType(credentials.ENUMCREDENTIALTYPEMETADATAFIELDSTYPE_DIRECTORY_ATTRIBUTE)
				field.SetAttribute(i.Attribute.ValueString())
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

	// tempCardTemplate := "<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 740 480\"><g><rect fill=\"none\" width=\"736\" height=\"476\" stroke=\"#808993\" stroke-width=\"3\" rx=\"10\" ry=\"10\" x=\"2\" y=\"2\"></rect><rect fill=\"${cardColor}\" height=\"476\" rx=\"10\" ry=\"10\" width=\"736\" x=\"2\" y=\"2\"></rect><image href=\"${backgroundImage}\" opacity=\"${bgOpacityPercent}\" height=\"476\" rx=\"10\" ry=\"10\" width=\"736\" x=\"2\" y=\"2\"></image><text fill=\"${textColor}\" font-weight=\"450\" font-size=\"30\" x=\"160\" y=\"80\">${cardTitle}</text><text fill=\"${textColor}\" font-size=\"25\" font-weight=\"300\" x=\"160\" y=\"120\">${cardTitle}</text><image href=\"${logoImage}\" x=\"42\" y=\"43\" height=\"90px\" width=\"90px\"></image><line xmlns=\"http://www.w3.org/2000/svg\" y2=\"160\" x2=\"695\" y1=\"160\" x1=\"42.5\" stroke=\"#808993\"></line><image href=\"data:image/jpeg;base64,${fields[0].value}\" x=\"42.5\" y=\"180\" rx=\"80px\" ry=\"80px\" height=\"130px\" width=\"130px\"></image><text fill=\"${textColor}\" font-weight=\"500\" font-size=\"20\" x=\"190\" y=\"230\">${fields[1].title}: ${fields[1].value}</text><text fill=\"${textColor}\" font-weight=\"500\" font-size=\"20\" x=\"190\" y=\"272\">${fields[2].title}: ${fields[2].value}</text><text fill=\"${textColor}\" font-weight=\"500\" font-size=\"20\" x=\"190\" y=\"314\">${fields[3].title}: ${fields[3].value}</text><text fill=\"${textColor}\" font-weight=\"500\" font-size=\"20\" x=\"190\" y=\"356\">${fields[3].title}: ${fields[3].value}</text></g></svg>"
	// tempCardTemplate := "<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 740 480\"><rect fill=\"none\" width=\"736\" height=\"476\" stroke=\"#CACED3\" stroke-width=\"3\" rx=\"10\" ry=\"10\" x=\"2\" y=\"2\"></rect><rect fill=\"${cardColor}\" height=\"476\" rx=\"10\" ry=\"10\" width=\"736\" x=\"2\" y=\"2\" opacity=\"${bgOpacityPercent}\"></rect><image href=\"${backgroundImage}\" opacity=\"${bgOpacityPercent}\" height=\"476\" rx=\"10\" ry=\"10\" width=\"736\" x=\"2\" y=\"2\"></image><line y2=\"160\" x2=\"695\" y1=\"160\" x1=\"42.5\" stroke=\"${textColor}\"></line><text fill=\"${textColor}\" font-weight=\"450\" font-size=\"30\" x=\"45\" y=\"90\">${cardTitle}</text><text fill=\"${textColor}\" font-size=\"25\" font-weight=\"300\" x=\"45\" y=\"130\">${cardTitle}</text><text fill=\"${textColor}\" font-weight=\"500\" font-size=\"20\" x=\"50\" y=\"228\">${fields[0].title}: ${fields[0].value}</text></svg>"
	tempCardTemplate := "<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 740 480\"><rect fill=\"none\" width=\"736\" height=\"476\" stroke=\"#CACED3\" stroke-width=\"3\" rx=\"10\" ry=\"10\" x=\"2\" y=\"2\"></rect><rect fill=\"${cardColor}\" height=\"476\" rx=\"10\" ry=\"10\" width=\"736\" x=\"2\" y=\"2\" opacity=\"${bgOpacityPercent}\"></rect><image href=\"${backgroundImage}\" opacity=\"${bgOpacityPercent}\" height=\"476\" rx=\"10\" ry=\"10\" width=\"736\" x=\"2\" y=\"2\"></image><image href=\"${logoImage}\" x=\"42\" y=\"43\" height=\"90px\" width=\"90px\"></image><line y2=\"160\" x2=\"695\" y1=\"160\" x1=\"42.5\" stroke=\"${textColor}\"></line><text fill=\"${textColor}\" font-weight=\"450\" font-size=\"30\" x=\"160\" y=\"90\">${cardTitle}</text><text fill=\"${textColor}\" font-size=\"25\" font-weight=\"300\" x=\"160\" y=\"130\"></text></svg>"
	//data := credentials.NewCredentialType(p.CardDesignTemplate.ValueString(), *cardMetadata, p.Title.ValueString())
	data := credentials.NewCredentialType(tempCardTemplate, *cardMetadata, p.Title.ValueString())

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

	return diags
}
