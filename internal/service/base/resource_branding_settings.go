package base

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
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
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type BrandingSettingsResource serviceClientType

type brandingSettingsResourceModel struct {
	Id            types.String `tfsdk:"id"`
	EnvironmentId types.String `tfsdk:"environment_id"`
	CompanyName   types.String `tfsdk:"company_name"`
	LogoImage     types.List   `tfsdk:"logo_image"`
}

type imageResourceModel struct {
	Id   types.String `tfsdk:"id"`
	Href types.String `tfsdk:"href"`
}

var (
	logoTFObjectTypes = map[string]attr.Type{
		"id":   types.StringType,
		"href": types.StringType,
	}
)

// Framework interfaces
var (
	_ resource.Resource                = &BrandingSettingsResource{}
	_ resource.ResourceWithConfigure   = &BrandingSettingsResource{}
	_ resource.ResourceWithImportState = &BrandingSettingsResource{}
)

// New Object
func NewBrandingSettingsResource() resource.Resource {
	return &BrandingSettingsResource{}
}

// Metadata
func (r *BrandingSettingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_branding_settings"
}

// Schema.
func (r *BrandingSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	const attrMinLength = 1

	logoIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The ID of the logo image.  This can be retrieved from the `id` parameter of the `pingone_image` resource.  Must be a valid PingOne resource ID.",
	)

	logoHrefDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The URL or fully qualified path to the logo file used for branding.  This can be retrieved from the `uploaded_image[0].href` parameter of the `pingone_image` resource.",
	)

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage the PingOne branding settings for an environment.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to set branding settings for."),
			),

			"company_name": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The company name associated with the specified environment.").Description,
				Optional:    true,
				Computed:    true,

				Default: stringdefault.StaticString(""),
			},
		},

		Blocks: map[string]schema.Block{

			"logo_image": schema.ListNestedBlock{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A single block that specifies the HREF and ID for the company logo.").Description,

				NestedObject: schema.NestedBlockObject{

					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description:         logoIdDescription.Description,
							MarkdownDescription: logoIdDescription.MarkdownDescription,
							Required:            true,

							Validators: []validator.String{
								verify.P1ResourceIDValidator(),
							},
						},

						"href": schema.StringAttribute{
							Description:         logoHrefDescription.Description,
							MarkdownDescription: logoHrefDescription.MarkdownDescription,
							Required:            true,
							Validators: []validator.String{
								stringvalidator.RegexMatches(verify.IsURLWithHTTPS, "Value must be a valid URL with `https://` prefix."),
							},
						},
					},
				},

				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},
		},
	}
}

func (r *BrandingSettingsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *BrandingSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state brandingSettingsResourceModel

	if r.Client == nil {
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
	brandingSettings, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *management.BrandingSettings
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return r.Client.BrandingSettingsApi.UpdateBrandingSettings(ctx, plan.EnvironmentId.ValueString()).BrandingSettings(*brandingSettings).Execute()
		},
		"Create::UpdateBrandingSettings",
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

func (r *BrandingSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *brandingSettingsResourceModel

	if r.Client == nil {
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
	var response *management.BrandingSettings
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return r.Client.BrandingSettingsApi.ReadBrandingSettings(ctx, data.EnvironmentId.ValueString()).Execute()
		},
		"ReadBrandingSettings",
		framework.DefaultCustomError,
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

func (r *BrandingSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state brandingSettingsResourceModel

	if r.Client == nil {
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
	brandingSettings, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *management.BrandingSettings
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return r.Client.BrandingSettingsApi.UpdateBrandingSettings(ctx, plan.EnvironmentId.ValueString()).BrandingSettings(*brandingSettings).Execute()
		},
		"Update::UpdateBrandingSettings",
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

func (r *BrandingSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *brandingSettingsResourceModel

	if r.Client == nil {
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

	brandingSettings := management.NewBrandingSettings()
	brandingSettings.SetCompanyName("")

	// Run the API call
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return r.Client.BrandingSettingsApi.UpdateBrandingSettings(ctx, data.EnvironmentId.ValueString()).BrandingSettings(*brandingSettings).Execute()
		},
		"Update::UpdateBrandingSettings",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		nil,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *BrandingSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splitLength := 1
	attributes := strings.SplitN(req.ID, "/", splitLength)

	if len(attributes) != splitLength {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("invalid id (\"%s\") specified, should be in format \"environment_id\"", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_id"), attributes[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), attributes[0])...)
}

func (p *brandingSettingsResourceModel) expand(ctx context.Context) (*management.BrandingSettings, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := management.NewBrandingSettings()

	if !p.CompanyName.IsNull() && !p.CompanyName.IsUnknown() {
		data.SetCompanyName(p.CompanyName.ValueString())
	} else {
		data.SetCompanyName("")
	}

	if !p.LogoImage.IsNull() && !p.LogoImage.IsUnknown() {

		var plan []imageResourceModel
		diags.Append(p.LogoImage.ElementsAs(ctx, &plan, false)...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetLogo(*management.NewBrandingSettingsLogo(plan[0].Href.ValueString(), plan[0].Id.ValueString()))
	}

	return data, diags
}

func (p *brandingSettingsResourceModel) toState(apiObject *management.BrandingSettings) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.StringToTF(apiObject.GetId())
	p.CompanyName = framework.EnumOkToTF(apiObject.GetCompanyNameOk())

	logoImage, d := toStateImageRef(apiObject.GetLogoOk())
	diags.Append(d...)
	p.LogoImage = logoImage

	return diags
}

func toStateImageRef(logo interface{}, ok bool) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	tfObjType := types.ObjectType{AttrTypes: logoTFObjectTypes}

	if !ok || logo == nil {
		return types.ListNull(tfObjType), diags
	}

	b, e := json.Marshal(logo)
	if e != nil {
		diags.AddError(
			"Invalid data object",
			fmt.Sprintf("Cannot remap the data object to JSON: %s.  Please report this to the provider maintainers.", e),
		)
		return types.ListNull(tfObjType), diags
	}

	var s map[string]string
	e = json.Unmarshal(b, &s)
	if e != nil {
		diags.AddError(
			"Invalid data object",
			fmt.Sprintf("Cannot remap the data object to map: %s.  Please report this to the provider maintainers.", e),
		)
		return types.ListNull(tfObjType), diags
	}

	logoMap := map[string]attr.Value{}

	if s["href"] != "" {
		logoMap["href"] = framework.StringToTF(s["href"])
	} else {
		logoMap["href"] = types.StringNull()
	}

	if s["id"] != "" {
		logoMap["id"] = framework.StringToTF(s["id"])
	} else {
		logoMap["id"] = types.StringNull()
	}

	flattenedObj, d := types.ObjectValue(logoTFObjectTypes, logoMap)
	diags.Append(d...)

	returnVar, d := types.ListValue(tfObjType, append([]attr.Value{}, flattenedObj))
	diags.Append(d...)

	return returnVar, diags
}
