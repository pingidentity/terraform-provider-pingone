package base

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
)

// Types
type BrandingThemeDefaultResource struct {
	client *management.APIClient
	region model.RegionMapping
}

type brandingThemeDefaultResourceModel struct {
	Id              types.String `tfsdk:"id"`
	EnvironmentId   types.String `tfsdk:"environment_id"`
	BrandingThemeId types.String `tfsdk:"branding_theme_id"`
	Default         types.Bool   `tfsdk:"default"`
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
					boolplanmodifier.UseStateForUnknown(),
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

func (r *BrandingThemeDefaultResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state brandingThemeDefaultResourceModel

	if r.client == nil {
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
	response, d := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.BrandingThemesApi.UpdateBrandingThemeDefault(ctx, plan.EnvironmentId.ValueString(), plan.BrandingThemeId.ValueString()).BrandingThemeDefault(*brandingThemeDefault).Execute()
		},
		"UpdateBrandingThemeDefault",
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
	resp.Diagnostics.Append(state.toState(response.(*management.BrandingThemeDefault))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *BrandingThemeDefaultResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *brandingThemeDefaultResourceModel

	if r.client == nil {
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
	response, diags := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.BrandingThemesApi.ReadBrandingThemeDefault(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
		},
		"ReadBrandingThemeDefault",
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
	resp.Diagnostics.Append(data.toState(response.(*management.BrandingThemeDefault))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *BrandingThemeDefaultResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *BrandingThemeDefaultResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *brandingThemeDefaultResourceModel

	if r.client == nil {
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

	bootstrapDefaultThemeId, d := r.fetchBootstapDefaultThemeId(ctx, r.client, data.EnvironmentId.ValueString())
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
		response, d := framework.ParseResponse(
			ctx,

			func() (interface{}, *http.Response, error) {
				return r.client.BrandingThemesApi.CreateBrandingTheme(ctx, data.EnvironmentId.ValueString()).BrandingTheme(*defaultTheme).Execute()
			},
			"CreateBrandingTheme",
			framework.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
		)
		resp.Diagnostics.Append(d...)

		bootstrapDefaultThemeId = response.(*management.BrandingTheme).Id
	}

	brandingThemeDefault := management.NewBrandingThemeDefault(true)

	_, d = framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.BrandingThemesApi.UpdateBrandingThemeDefault(ctx, data.EnvironmentId.ValueString(), *bootstrapDefaultThemeId).BrandingThemeDefault(*brandingThemeDefault).Execute()
		},
		"UpdateBrandingThemeDefault",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	resp.Diagnostics.Append(d...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *BrandingThemeDefaultResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splitLength := 1
	attributes := strings.SplitN(req.ID, "/", splitLength)

	if len(attributes) != splitLength {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("invalid id (\"%s\") specified, should be in format \"environment_id\"", req.ID),
		)
		return
	}

	defaultThemeId, d := r.fetchBootstapDefaultThemeId(ctx, r.client, attributes[0])
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_id"), attributes[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), defaultThemeId)...)
}

func (r *BrandingThemeDefaultResource) fetchBootstapDefaultThemeId(ctx context.Context, apiClient *management.APIClient, environmentID string) (*string, diag.Diagnostics) {
	var diags diag.Diagnostics

	response, diags := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.BrandingThemesApi.ReadBrandingThemes(ctx, environmentID).Execute()
		},
		"ReadBrandingThemes",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	diags.Append(diags...)
	if diags.HasError() {
		return nil, diags
	}

	if brandingThemes, ok := response.(*management.EntityArray).Embedded.GetThemesOk(); ok {

		for _, brandingTheme := range brandingThemes {
			if *brandingTheme.GetConfiguration().Name == "Ping Default" {
				defaultThemeId := brandingTheme.GetId()
				return &defaultThemeId, diags
			}
		}
	}

	return nil, diags
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
