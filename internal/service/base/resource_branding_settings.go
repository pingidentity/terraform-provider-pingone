// Copyright © 2025 Ping Identity Corporation

package base

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/service"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type BrandingSettingsResource serviceClientType

type brandingSettingsResourceModelV1 struct {
	Id            pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	CompanyName   types.String                 `tfsdk:"company_name"`
	LogoImage     types.Object                 `tfsdk:"logo_image"`
}

// Framework interfaces
var (
	_ resource.Resource                 = &BrandingSettingsResource{}
	_ resource.ResourceWithConfigure    = &BrandingSettingsResource{}
	_ resource.ResourceWithImportState  = &BrandingSettingsResource{}
	_ resource.ResourceWithUpgradeState = &BrandingSettingsResource{}
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
		"The URL or fully qualified path to the logo file used for branding.  This can be retrieved from the `uploaded_image.href` parameter of the `pingone_image` resource.",
	)

	resp.Schema = schema.Schema{

		Version: 1,

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

			"logo_image": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that specifies the HREF and ID for the company logo.").Description,
				Optional:    true,
				Computed:    true,

				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description:         logoIdDescription.Description,
						MarkdownDescription: logoIdDescription.MarkdownDescription,
						Required:            true,

						CustomType: pingonetypes.ResourceIDType{},
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

	r.Client = resourceConfig.Client.API
	if r.Client == nil {
		resp.Diagnostics.AddError(
			"Client not initialised",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.",
		)
		return
	}
}

func (r *BrandingSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state brandingSettingsResourceModelV1

	if r.Client == nil || r.Client.ManagementAPIClient == nil {
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
			fO, fR, fErr := r.Client.ManagementAPIClient.BrandingSettingsApi.UpdateBrandingSettings(ctx, plan.EnvironmentId.ValueString()).BrandingSettings(*brandingSettings).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
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
	var data *brandingSettingsResourceModelV1

	if r.Client == nil || r.Client.ManagementAPIClient == nil {
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
			fO, fR, fErr := r.Client.ManagementAPIClient.BrandingSettingsApi.ReadBrandingSettings(ctx, data.EnvironmentId.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadBrandingSettings",
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

func (r *BrandingSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state brandingSettingsResourceModelV1

	if r.Client == nil || r.Client.ManagementAPIClient == nil {
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
			fO, fR, fErr := r.Client.ManagementAPIClient.BrandingSettingsApi.UpdateBrandingSettings(ctx, plan.EnvironmentId.ValueString()).BrandingSettings(*brandingSettings).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
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
	var data *brandingSettingsResourceModelV1

	if r.Client == nil || r.Client.ManagementAPIClient == nil {
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
			fO, fR, fErr := r.Client.ManagementAPIClient.BrandingSettingsApi.UpdateBrandingSettings(ctx, data.EnvironmentId.ValueString()).BrandingSettings(*brandingSettings).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"Update::UpdateBrandingSettings",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
		nil,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *BrandingSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idComponents := []framework.ImportComponent{
		{
			Label:     "environment_id",
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

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_id"), framework.PingOneResourceIDToTF(attributes["environment_id"]))...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), framework.PingOneResourceIDToTF(attributes["environment_id"]))...)
}

func (p *brandingSettingsResourceModelV1) expand(ctx context.Context) (*management.BrandingSettings, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := management.NewBrandingSettings()

	if !p.CompanyName.IsNull() && !p.CompanyName.IsUnknown() {
		data.SetCompanyName(p.CompanyName.ValueString())
	} else {
		data.SetCompanyName("")
	}

	if !p.LogoImage.IsNull() && !p.LogoImage.IsUnknown() {

		var plan service.ImageResourceModel
		diags.Append(p.LogoImage.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetLogo(*management.NewBrandingSettingsLogo(plan.Href.ValueString(), plan.Id.ValueString()))
	}

	return data, diags
}

func (p *brandingSettingsResourceModelV1) toState(apiObject *management.BrandingSettings) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.PingOneResourceIDToTF(apiObject.GetId())
	p.CompanyName = framework.EnumOkToTF(apiObject.GetCompanyNameOk())

	logoImage, d := service.ImageOkToTF(apiObject.GetLogoOk())
	diags.Append(d...)
	p.LogoImage = logoImage

	return diags
}
