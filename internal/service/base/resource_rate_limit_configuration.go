// Copyright © 2025 Ping Identity Corporation

package base

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/legacysdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type RateLimitConfigurationResource serviceClientType

type RateLimitConfigurationResourceModel struct {
	Id            pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	Type          types.String                 `tfsdk:"type"`
	Value         types.String                 `tfsdk:"value"`
	CreatedAt     timetypes.RFC3339            `tfsdk:"created_at"`
	UpdatedAt     timetypes.RFC3339            `tfsdk:"updated_at"`
}

// Framework interfaces
var (
	_ resource.Resource                = &RateLimitConfigurationResource{}
	_ resource.ResourceWithConfigure   = &RateLimitConfigurationResource{}
	_ resource.ResourceWithImportState = &RateLimitConfigurationResource{}
)

// New Object
func NewRateLimitConfigurationResource() resource.Resource {
	return &RateLimitConfigurationResource{}
}

// Metadata
func (r *RateLimitConfigurationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_rate_limit_configuration"
}

// Schema
func (r *RateLimitConfigurationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	providerDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Resource to create and manage rate limit configurations in PingOne. Rate limit configurations allow you to exclude specific IP addresses or CIDR ranges from rate limiting.",
	)

	typeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The type of rate limit configuration. Currently, the only type supported is `WHITELIST`, indicating that the IP address in `value` is to be excluded from rate limiting.",
	).DefaultValue(string(management.ENUMRATELIMITCONFIGURATIONTYPE_WHITELIST))

	valueDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The IP address (IPv4 or IPv6), or a CIDR range, for the IP address or addresses to be excluded from rate limiting.",
	)

	createdAtDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the time the resource was created in RFC3339 format.",
	)

	updatedAtDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the time the resource was last updated in RFC3339 format.",
	)

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: providerDescription.MarkdownDescription,
		Description:         providerDescription.Description,

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),
			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to create the rate limit configuration in."),
			),
			"type": schema.StringAttribute{
				MarkdownDescription: typeDescription.MarkdownDescription,
				Description:         typeDescription.Description,
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(string(management.ENUMRATELIMITCONFIGURATIONTYPE_WHITELIST)),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(string(management.ENUMRATELIMITCONFIGURATIONTYPE_WHITELIST)),
				},
			},

			"value": schema.StringAttribute{
				MarkdownDescription: valueDescription.MarkdownDescription,
				Description:         valueDescription.Description,
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(verify.IPv4IPv6Regexp, "Values must be valid IPv4 or IPv6 CIDR format."),
				},
			},

			"created_at": schema.StringAttribute{
				MarkdownDescription: createdAtDescription.MarkdownDescription,
				Description:         createdAtDescription.Description,
				Computed:            true,
				CustomType:          timetypes.RFC3339Type{},
			},

			"updated_at": schema.StringAttribute{
				MarkdownDescription: updatedAtDescription.MarkdownDescription,
				Description:         updatedAtDescription.Description,
				Computed:            true,
				CustomType:          timetypes.RFC3339Type{},
			},
		},
	}
}

func (r *RateLimitConfigurationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *RateLimitConfigurationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state RateLimitConfigurationResourceModel

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
	rateLimitConfiguration := plan.expand()

	// Run the API call
	var response *management.RateLimitConfiguration
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.RateLimitingApi.CreateRateLimitConfiguration(ctx, plan.EnvironmentId.ValueString()).RateLimitConfiguration(*rateLimitConfiguration).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"CreateRateLimitConfiguration",
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

func (r *RateLimitConfigurationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *RateLimitConfigurationResourceModel

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
	var response *management.RateLimitConfiguration
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.RateLimitingApi.ReadOneRateLimitConfiguration(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadOneRateLimitConfiguration",
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

func (r *RateLimitConfigurationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Rate limit configurations are immutable - all fields require replacement
	// This method should never be called as all fields have RequiresReplace plan modifiers
}

func (r *RateLimitConfigurationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *RateLimitConfigurationResourceModel

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
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fR, fErr := r.Client.ManagementAPIClient.RateLimitingApi.DeleteRateLimitConfiguration(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, fR, fErr)
		},
		"DeleteRateLimitConfiguration",
		legacysdk.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
		nil,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *RateLimitConfigurationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:     "rate_limit_configuration_id",
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

	for _, idComponent := range idComponents {
		pathKey := idComponent.Label

		if idComponent.PrimaryID {
			pathKey = "id"
		}

		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(pathKey), attributes[idComponent.Label])...)
	}
}

func (p *RateLimitConfigurationResourceModel) expand() *management.RateLimitConfiguration {
	data := management.NewRateLimitConfiguration(
		management.EnumRateLimitConfigurationType(p.Type.ValueString()),
		p.Value.ValueString(),
	)

	return data
}

func (p *RateLimitConfigurationResourceModel) toState(apiObject *management.RateLimitConfiguration) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.PingOneResourceIDToTF(apiObject.GetId())
	p.Type = framework.EnumOkToTF(apiObject.GetTypeOk())
	p.Value = framework.StringOkToTF(apiObject.GetValueOk())
	p.CreatedAt = framework.TimeOkToTF(apiObject.GetCreatedAtOk())
	p.UpdatedAt = framework.TimeOkToTF(apiObject.GetUpdatedAtOk())

	return diags
}
