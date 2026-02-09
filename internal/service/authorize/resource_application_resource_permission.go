// Copyright © 2026 Ping Identity Corporation

package authorize

import (
	"context"
	"fmt"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/authorize"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/legacysdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type ApplicationResourcePermissionResource serviceClientType

type ApplicationResourcePermissionResourceModel struct {
	Id                    pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId         pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	ApplicationResourceId pingonetypes.ResourceIDValue `tfsdk:"application_resource_id"`
	Action                types.String                 `tfsdk:"action"`
	Description           types.String                 `tfsdk:"description"`
	Resource              types.Object                 `tfsdk:"resource"`
}

var (
	applicationResourcePermissionResourceTFObjectTypes = map[string]attr.Type{
		"id":   pingonetypes.ResourceIDType{},
		"name": types.StringType,
	}
)

// Framework interfaces
var (
	_ resource.Resource                = &ApplicationResourcePermissionResource{}
	_ resource.ResourceWithConfigure   = &ApplicationResourcePermissionResource{}
	_ resource.ResourceWithImportState = &ApplicationResourcePermissionResource{}
)

// New Object
func NewApplicationResourcePermissionResource() resource.Resource {
	return &ApplicationResourcePermissionResource{}
}

// Metadata
func (r *ApplicationResourcePermissionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application_resource_permission"
}

// Schema.
func (r *ApplicationResourcePermissionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	const attrMinLength = 1

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage application resource permissions in a PingOne environment.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to create the resource attribute in."),
			),

			"application_resource_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the application resource to create and manage permissions for."),
			),

			"action": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the action associated with this permission.  The action must contain only Unicode letters, marks, numbers, spaces, forward slashes, dots, apostrophes, underscores, or hyphens.").Description,
				Required:    true,

				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(`^[\p{L}\p{M}\p{N} /.'_-]*$`), "The action must contain only Unicode letters, marks, numbers, spaces, forward slashes, dots, apostrophes, underscores, or hyphens"),
				},
			},

			"description": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the permission's description.").Description,
				Optional:    true,
			},

			"resource": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that describes the associated application resource.").Description,
				Computed:    true,

				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the ID for the associated application resource.").Description,
						Computed:    true,

						CustomType: pingonetypes.ResourceIDType{},
					},

					"name": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that represents the name of the associated application resource.").Description,
						Computed:    true,
					},
				},
			},
		},
	}
}

func (r *ApplicationResourcePermissionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ApplicationResourcePermissionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state ApplicationResourcePermissionResourceModel

	if r.Client == nil || r.Client.AuthorizeAPIClient == nil {
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
	applicationResourcePermission := plan.expand()

	// Run the API call
	var response *authorize.ApplicationResourcePermission
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.AuthorizeAPIClient.ApplicationResourcePermissionsApi.CreateApplicationPermission(ctx, plan.EnvironmentId.ValueString(), plan.ApplicationResourceId.ValueString()).ApplicationResourcePermission(*applicationResourcePermission).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"CreateApplicationPermission",
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

func (r *ApplicationResourcePermissionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *ApplicationResourcePermissionResourceModel

	if r.Client == nil || r.Client.AuthorizeAPIClient == nil {
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
	var response *authorize.ApplicationResourcePermission
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.AuthorizeAPIClient.ApplicationResourcePermissionsApi.ReadOneApplicationPermission(ctx, data.EnvironmentId.ValueString(), data.ApplicationResourceId.ValueString(), data.Id.ValueString()).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadOneApplicationPermission",
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

func (r *ApplicationResourcePermissionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state ApplicationResourcePermissionResourceModel

	if r.Client == nil || r.Client.AuthorizeAPIClient == nil {
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
	resourceAttribute := plan.expand()

	// Run the API call
	var response *authorize.ApplicationResourcePermission
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.AuthorizeAPIClient.ApplicationResourcePermissionsApi.UpdateApplicationPermission(ctx, plan.EnvironmentId.ValueString(), plan.ApplicationResourceId.ValueString(), plan.Id.ValueString()).ApplicationResourcePermission(*resourceAttribute).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateApplicationPermission",
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

func (r *ApplicationResourcePermissionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *ApplicationResourcePermissionResourceModel

	if r.Client == nil || r.Client.AuthorizeAPIClient == nil {
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
			fR, fErr := r.Client.AuthorizeAPIClient.ApplicationResourcePermissionsApi.DeleteApplicationPermission(ctx, data.EnvironmentId.ValueString(), data.ApplicationResourceId.ValueString(), data.Id.ValueString()).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, fR, fErr)
		},
		"DeleteApplicationPermission",
		legacysdk.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
		nil,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ApplicationResourcePermissionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:  "application_resource_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:     "application_resource_permission_id",
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

func (p *ApplicationResourcePermissionResourceModel) expand() *authorize.ApplicationResourcePermission {

	data := authorize.NewApplicationResourcePermission(
		p.Action.ValueString(),
	)

	if !p.Description.IsNull() {
		data.SetDescription(p.Description.ValueString())
	}

	return data
}

func (p *ApplicationResourcePermissionResourceModel) toState(apiObject *authorize.ApplicationResourcePermission) diag.Diagnostics {
	var diags, d diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.PingOneResourceIDOkToTF(apiObject.GetIdOk())
	p.Action = framework.StringOkToTF(apiObject.GetActionOk())
	p.Description = framework.StringOkToTF(apiObject.GetDescriptionOk())
	p.Resource, d = toStateApplicationResourcePermissionResourceOk(apiObject.GetResourceOk())
	diags.Append(d...)

	return diags
}

func toStateApplicationResourcePermissionResourceOk(apiObject *authorize.ApplicationResourcePermissionResource, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(applicationResourcePermissionResourceTFObjectTypes), diags
	}

	o := map[string]attr.Value{
		"id":   framework.PingOneResourceIDOkToTF(apiObject.GetIdOk()),
		"name": framework.StringOkToTF(apiObject.GetNameOk()),
	}

	returnVar, d := types.ObjectValue(applicationResourcePermissionResourceTFObjectTypes, o)
	diags.Append(d...)

	return returnVar, diags
}
