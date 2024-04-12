package sso

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type ApplicationSecretResource serviceClientType

type ApplicationSecretResourceModel struct {
	EnvironmentId           types.String `tfsdk:"environment_id"`
	ApplicationId           types.String `tfsdk:"application_id"`
	Previous                types.Object `tfsdk:"previous"`
	Secret                  types.String `tfsdk:"secret"`
	RegenerateTriggerValues types.Map    `tfsdk:"regenerate_trigger_values"`
}

type ApplicationSecretPreviousResourceModel struct {
	Secret    types.String      `tfsdk:"secret"`
	ExpiresAt timetypes.RFC3339 `tfsdk:"expires_at"`
	LastUsed  types.String      `tfsdk:"last_used"`
}

// Framework interfaces
var (
	_ resource.Resource                = &ApplicationSecretResource{}
	_ resource.ResourceWithConfigure   = &ApplicationSecretResource{}
	_ resource.ResourceWithImportState = &ApplicationSecretResource{}
	_ resource.ResourceWithModifyPlan  = &ApplicationSecretResource{}
)

// New Object
func NewApplicationSecretResource() resource.Resource {
	return &ApplicationSecretResource{}
}

// Metadata
func (r *ApplicationSecretResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application_secret"
}

// Schema
func (r *ApplicationSecretResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to generate and read an application secret for an administrator defined application configured in PingOne.",

		Attributes: map[string]schema.Attribute{
			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to generate an application secret in."),
			),

			"application_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the application to generate the application secret for. The value for `application_id` may come from the `id` attribute of the `pingone_application` resource or data source."),
			),

			"previous": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("An object that specifies the previous secret, when it expires, and when it was last used.").Description,
				Optional:    true,

				Attributes: map[string]schema.Attribute{
					"secret": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the previous application secret. This property is returned in the response if the previous secret is not expired.").Description,
						Computed:    true,
						Sensitive:   true,
					},

					"expires_at": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A timestamp that specifies how long this secret is saved (and can be used) before it expires. Supported time range is 1 minute to 30 days.").Description,
						Optional:    true,

						CustomType: timetypes.RFC3339Type{},
					},

					"last_used": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A timestamp that specifies when the previous secret was last used.").Description,
						Computed:    true,
					},
				},
			},

			"secret": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The application secret used to authenticate to the authorization server.").Description,
				Computed:    true,
				Sensitive:   true,
			},

			"regenerate_trigger_values": schema.MapAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A meta-argument map of values that, if any values are changed, will force regeneration of the application secret.  Adding values to and removing values from the map will not trigger a secret regeneration.  This parameter can be used to control time-based rotation using Terraform.").Description,
				Optional:    true,

				ElementType: types.StringType,
			},
		},
	}
}

// ModifyPlan
func (r *ApplicationSecretResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {

	// Destruction plan
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan, state types.Map
	var planValues, stateValues map[string]attr.Value

	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("regenerate_trigger_values"), &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	planValues = plan.Elements()

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("regenerate_trigger_values"), &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	stateValues = state.Elements()

	for k, v := range planValues {
		if stateValue, ok := stateValues[k]; ok && (v == types.StringUnknown() || !stateValue.Equal(v)) {
			resp.RequiresReplace = path.Paths{path.Root("regenerate_trigger_values")}
			break
		}
	}

}

func (r *ApplicationSecretResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ApplicationSecretResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state ApplicationSecretResourceModel

	if r.Client.ManagementAPIClient == nil {
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
	applicationSecret, d := plan.expand()
	resp.Diagnostics = append(resp.Diagnostics, d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.ApplicationSecretApi.UpdateApplicationSecret(ctx, plan.EnvironmentId.ValueString(), plan.ApplicationId.ValueString()).ApplicationSecret(*applicationSecret).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateApplicationSecret",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		nil,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var response *management.ApplicationSecret
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.ApplicationSecretApi.ReadApplicationSecret(ctx, plan.EnvironmentId.ValueString(), plan.ApplicationId.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadApplicationSecret",
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

func (r *ApplicationSecretResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *ApplicationSecretResourceModel

	if r.Client.ManagementAPIClient == nil {
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
	var response *management.ApplicationSecret
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.ApplicationSecretApi.ReadApplicationSecret(ctx, data.EnvironmentId.ValueString(), data.ApplicationId.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadApplicationSecret",
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

func (r *ApplicationSecretResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state ApplicationSecretResourceModel

	if r.Client.ManagementAPIClient == nil {
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
	applicationSecret, d := plan.expand()
	resp.Diagnostics = append(resp.Diagnostics, d...)
	if resp.Diagnostics.HasError() {
		return
	}

	var response *management.ApplicationSecret
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.ApplicationSecretApi.UpdateApplicationSecret(ctx, plan.EnvironmentId.ValueString(), plan.ApplicationId.ValueString()).ApplicationSecret(*applicationSecret).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateApplicationSecret",
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

func (r *ApplicationSecretResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *ApplicationSecretResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:  "application_id",
			Regexp: verify.P1ResourceIDRegexp,
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

func (p *ApplicationSecretResourceModel) expand() (*management.ApplicationSecret, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := management.NewApplicationSecret()

	if !p.Previous.IsNull() && !p.Previous.IsUnknown() {
		var plan ApplicationSecretPreviousResourceModel
		d := p.Previous.As(context.Background(), &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		if !plan.ExpiresAt.IsNull() && !plan.ExpiresAt.IsUnknown() {

			expiresAt, d := plan.ExpiresAt.ValueRFC3339Time()
			diags.Append(d...)
			if diags.HasError() {
				return nil, diags
			}

			data.SetPrevious(*management.NewApplicationSecretPrevious(expiresAt))
		}
	}

	return data, diags
}

func (p *ApplicationSecretResourceModel) toState(apiObject *management.ApplicationSecret) diag.Diagnostics {
	var diags, d diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.EnvironmentId = framework.StringToTF(*apiObject.GetEnvironment().Id)
	p.Secret = framework.StringOkToTF(apiObject.GetSecretOk())
	p.Previous, d = applicationSecretPreviousOkToTF(apiObject.GetPreviousOk())
	diags.Append(d...)

	return diags
}
