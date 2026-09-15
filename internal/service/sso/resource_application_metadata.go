// Copyright © 2026 Ping Identity Corporation

package sso

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/legacysdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type ApplicationMetadataResource serviceClientType

type ApplicationMetadataResourceModel struct {
	EnvironmentId pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	ApplicationId pingonetypes.ResourceIDValue `tfsdk:"application_id"`
	Metadata      jsontypes.Normalized         `tfsdk:"metadata"`
}

// Framework interfaces
var (
	_ resource.Resource                = &ApplicationMetadataResource{}
	_ resource.ResourceWithConfigure   = &ApplicationMetadataResource{}
	_ resource.ResourceWithImportState = &ApplicationMetadataResource{}
)

// New Object
func NewApplicationMetadataResource() resource.Resource {
	return &ApplicationMetadataResource{}
}

// Metadata
func (r *ApplicationMetadataResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application_metadata"
}

// Schema
func (r *ApplicationMetadataResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage custom metadata for an application in PingOne.  Metadata is a user-defined JSON object that is associated with the application, and can be removed by destroying this resource.",

		Attributes: map[string]schema.Attribute{
			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to manage the application metadata in."),
			),

			"application_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the application to manage the metadata for. The value for `application_id` may come from the `id` attribute of the `pingone_application` resource or data source."),
			),

			"metadata": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A JSON string that specifies user-defined custom metadata for the application.  The top level of the JSON must be an object (a map of key-value pairs).  If metadata is not included in the request, the existing application metadata will be removed; to remove application metadata, destroy this resource.").Description,
				Required:    true,

				CustomType: jsontypes.NormalizedType{},

				Validators: []validator.String{
					jsonObjectValidator{},
				},
			},
		},
	}
}

func (r *ApplicationMetadataResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ApplicationMetadataResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state ApplicationMetadataResourceModel

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
	applicationMetadata, d := plan.expand()
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *management.ApplicationMetadata
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.ApplicationMetadataApi.UpdateApplicationMetadata(ctx, plan.EnvironmentId.ValueString(), plan.ApplicationId.ValueString()).ApplicationMetadata(*applicationMetadata).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateApplicationMetadata",
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

func (r *ApplicationMetadataResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *ApplicationMetadataResourceModel

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

	// Guard against corrupt state missing the identifier attributes, which would
	// otherwise issue a GET against empty path segments.
	if data.EnvironmentId.IsNull() || data.ApplicationId.IsNull() {
		resp.Diagnostics.AddError(
			"Missing resource identifiers",
			"The environment ID or application ID is missing from the Terraform state.  The state for this resource appears to be corrupt; remove the resource from state and import it again.",
		)
		return
	}

	// Run the API call
	var response *management.ApplicationMetadata
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.ApplicationMetadataApi.ReadApplicationMetadata(ctx, data.EnvironmentId.ValueString(), data.ApplicationId.ValueString()).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadApplicationMetadata",
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

func (r *ApplicationMetadataResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state ApplicationMetadataResourceModel

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
	applicationMetadata, d := plan.expand()
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *management.ApplicationMetadata
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.ApplicationMetadataApi.UpdateApplicationMetadata(ctx, plan.EnvironmentId.ValueString(), plan.ApplicationId.ValueString()).ApplicationMetadata(*applicationMetadata).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateApplicationMetadata",
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

func (r *ApplicationMetadataResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *ApplicationMetadataResourceModel

	if r.Client == nil || r.Client.ManagementAPIClient == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report it to the provider maintainers.")
		return
	}

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call.  There is no dedicated DELETE endpoint for application metadata;
	// a PUT with no `metadata` property removes the existing metadata from the application.
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.ApplicationMetadataApi.UpdateApplicationMetadata(ctx, data.EnvironmentId.ValueString(), data.ApplicationId.ValueString()).ApplicationMetadata(*management.NewApplicationMetadata()).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateApplicationMetadata",
		legacysdk.CustomErrorResourceNotFoundWarning,
		nil,
		nil,
	)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ApplicationMetadataResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

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

func (p *ApplicationMetadataResourceModel) expand() (*management.ApplicationMetadata, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := management.NewApplicationMetadata()

	if !p.Metadata.IsNull() && !p.Metadata.IsUnknown() {
		var metadata map[string]interface{}
		diags.Append(p.Metadata.Unmarshal(&metadata)...)
		if !diags.HasError() {
			data.SetMetadata(metadata)
		}
	}

	return data, diags
}

func (p *ApplicationMetadataResourceModel) toState(apiObject *management.ApplicationMetadata) diag.Diagnostics {
	var diags, d diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	// The API does not consistently echo the `environmentId` and `applicationId` properties
	// on the PUT response, so only overwrite the plan/state values when the API returns them.
	if v, ok := apiObject.GetEnvironmentIdOk(); ok && v != nil {
		p.EnvironmentId = framework.PingOneResourceIDToTF(*v)
	}

	if v, ok := apiObject.GetApplicationIdOk(); ok && v != nil {
		p.ApplicationId = framework.PingOneResourceIDToTF(*v)
	}

	p.Metadata, d = framework.JSONNormalizedOkToTF(apiObject.GetMetadataOk())
	diags.Append(d...)

	return diags
}

// jsonObjectValidator validates that the string value is a JSON object at the top level.
type jsonObjectValidator struct{}

func (v jsonObjectValidator) Description(_ context.Context) string {
	return "Ensure the string contains a JSON object value."
}

func (v jsonObjectValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v jsonObjectValidator) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	var jsonObject map[string]interface{}
	if err := json.Unmarshal([]byte(req.ConfigValue.ValueString()), &jsonObject); err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid JSON Object Value",
			fmt.Sprintf("A string value was provided that is not a JSON object (a map of key-value pairs).\n\nGiven Value: %s\nError: %s", req.ConfigValue.ValueString(), err.Error()),
		)
	}
}
