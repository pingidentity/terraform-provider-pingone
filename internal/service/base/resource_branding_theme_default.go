// Copyright © 2025 Ping Identity Corporation

package base

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/legacysdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type BrandingThemeDefaultResource serviceClientType

type brandingThemeDefaultResourceModel struct {
	Id              pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId   pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	BrandingThemeId pingonetypes.ResourceIDValue `tfsdk:"branding_theme_id"`
	Default         types.Bool                   `tfsdk:"default"`
}

// Framework interfaces
var (
	_ resource.Resource                = &BrandingThemeDefaultResource{}
	_ resource.ResourceWithConfigure   = &BrandingThemeDefaultResource{}
	_ resource.ResourceWithImportState = &BrandingThemeDefaultResource{}
)

// New Object
func NewBrandingThemeDefaultResource() resource.Resource {
	return &BrandingThemeDefaultResource{}
}

// Metadata
func (r *BrandingThemeDefaultResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_branding_theme_default"
}

// Schema.
func (r *BrandingThemeDefaultResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to manage the default PingOne branding theme for an environment.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to set branding settings for."),
			),

			"branding_theme_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the branding theme to activate as the environment default."),
			),

			"default": schema.BoolAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("Confirms whether this theme is the environment's default branding configuration.").Description,
				Computed:    true,

				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseNonNullStateForUnknown(),
				},
			},
		},
	}
}

func (r *BrandingThemeDefaultResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	resourceConfig, ok := req.ProviderData.(legacysdk.ResourceType)
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

func (r *BrandingThemeDefaultResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state brandingThemeDefaultResourceModel

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
	brandingThemeDefault := plan.expand()

	// Run the API call
	var response *management.BrandingThemeDefault
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.BrandingThemesApi.UpdateBrandingThemeDefault(ctx, plan.EnvironmentId.ValueString(), plan.BrandingThemeId.ValueString()).BrandingThemeDefault(*brandingThemeDefault).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateBrandingThemeDefault",
		legacysdk.DefaultCustomError,
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

func (r *BrandingThemeDefaultResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *brandingThemeDefaultResourceModel

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
	var response *management.BrandingThemeDefault
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.BrandingThemesApi.ReadBrandingThemeDefault(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadBrandingThemeDefault",
		legacysdk.CustomErrorResourceNotFoundWarning,
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

func (r *BrandingThemeDefaultResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *BrandingThemeDefaultResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *brandingThemeDefaultResourceModel

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

	bootstrapDefaultThemeId, d := r.fetchBootstapDefaultThemeId(ctx, data.EnvironmentId.ValueString())
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	if bootstrapDefaultThemeId == nil {
		configuration := management.NewBrandingThemeConfiguration(
			management.ENUMBRANDINGTHEMEBACKGROUNDTYPE_DEFAULT,
			"#686f77",
			"#007CBA",
			"#ffffff",
			"#ffffff",
			"#686f77",
			"#007CBA",
			management.ENUMBRANDINGLOGOTYPE_NONE,
		)

		configuration.SetBackgroundColor("#ededed")
		configuration.SetName("Ping Default")

		defaultTheme := management.NewBrandingTheme(
			*configuration,
			true,
			management.ENUMBRANDINGTHEMETEMPLATE_DEFAULT,
		)

		// Run the API call
		var response *management.BrandingTheme
		resp.Diagnostics.Append(legacysdk.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				fO, fR, fErr := r.Client.ManagementAPIClient.BrandingThemesApi.CreateBrandingTheme(ctx, data.EnvironmentId.ValueString()).BrandingTheme(*defaultTheme).Execute()
				return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
			},
			"CreateBrandingTheme",
			legacysdk.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			&response,
		)...)

		bootstrapDefaultThemeId = response.Id
	}

	brandingThemeDefault := management.NewBrandingThemeDefault(true)

	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.BrandingThemesApi.UpdateBrandingThemeDefault(ctx, data.EnvironmentId.ValueString(), *bootstrapDefaultThemeId).BrandingThemeDefault(*brandingThemeDefault).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateBrandingThemeDefault",
		legacysdk.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		nil,
	)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *BrandingThemeDefaultResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

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

	defaultThemeId, d := r.fetchDefaultThemeId(ctx, attributes["environment_id"])
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	if defaultThemeId == nil {
		resp.Diagnostics.AddError(
			"Default theme not found",
			"Unable to find the default theme for the environment.",
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_id"), attributes["environment_id"])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("branding_theme_id"), defaultThemeId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), defaultThemeId)...)
}

func (r *BrandingThemeDefaultResource) fetchBootstapDefaultThemeId(ctx context.Context, environmentID string) (*string, diag.Diagnostics) {
	return r.fetchThemeId(ctx, environmentID, true)
}

func (r *BrandingThemeDefaultResource) fetchDefaultThemeId(ctx context.Context, environmentID string) (*string, diag.Diagnostics) {
	return r.fetchThemeId(ctx, environmentID, false)
}

func (r *BrandingThemeDefaultResource) fetchThemeId(ctx context.Context, environmentID string, bootstrapDefault bool) (*string, diag.Diagnostics) {
	var diags diag.Diagnostics

	var response *string
	diags.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			pagedIterator := r.Client.ManagementAPIClient.BrandingThemesApi.ReadBrandingThemes(ctx, environmentID).Execute()

			var initialHttpResponse *http.Response

			for pageCursor, err := range pagedIterator {
				if err != nil {
					return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, environmentID, nil, pageCursor.HTTPResponse, err)
				}

				if initialHttpResponse == nil {
					initialHttpResponse = pageCursor.HTTPResponse
				}

				if brandingThemes, ok := pageCursor.EntityArray.Embedded.GetThemesOk(); ok {

					for _, brandingTheme := range brandingThemes {
						if bootstrapDefault && *brandingTheme.GetConfiguration().Name == "Ping Default" {
							defaultThemeId := brandingTheme.GetId()
							return &defaultThemeId, pageCursor.HTTPResponse, nil
						}

						if !bootstrapDefault && brandingTheme.GetDefault() {
							defaultThemeId := brandingTheme.GetId()
							return &defaultThemeId, pageCursor.HTTPResponse, nil
						}
					}
				}
			}

			return nil, initialHttpResponse, nil
		},
		"ReadBrandingThemes",
		legacysdk.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		&response,
	)...)
	if diags.HasError() {
		return nil, diags
	}

	return response, diags
}

func (p *brandingThemeDefaultResourceModel) expand() *management.BrandingThemeDefault {
	return management.NewBrandingThemeDefault(true)
}

func (p *brandingThemeDefaultResourceModel) toState(apiObject *management.BrandingThemeDefault) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = p.BrandingThemeId
	p.Default = framework.BoolOkToTF(apiObject.GetDefaultOk())

	return diags
}
