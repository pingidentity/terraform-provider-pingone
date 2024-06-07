package sso

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type ApplicationResourceResource serviceClientType

type ApplicationResourceResourceModel struct {
	Id            pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	ResourceId    pingonetypes.ResourceIDValue `tfsdk:"resource_id"`
	ResourceName  types.String                 `tfsdk:"resource_name"`
	Name          types.String                 `tfsdk:"name"`
	Description   types.String                 `tfsdk:"description"`
	Parent        types.Object                 `tfsdk:"parent"`
}

var (
	applicationResourceParentTFObjectTypes = map[string]attr.Type{
		"type": types.StringType,
		"id":   pingonetypes.ResourceIDType{},
	}
)

// Framework interfaces
var (
	_ resource.Resource                = &ApplicationResourceResource{}
	_ resource.ResourceWithConfigure   = &ApplicationResourceResource{}
	_ resource.ResourceWithImportState = &ApplicationResourceResource{}
)

// New Object
func NewApplicationResourceResource() resource.Resource {
	return &ApplicationResourceResource{}
}

// Metadata
func (r *ApplicationResourceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application_resource"
}

// Schema.
func (r *ApplicationResourceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	const attrMinLength = 1

	resourceIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The ID of the resource that is assigned as an application resource.",
	)

	resourceNameDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The name of the resource to should be assigned as an application resource.",
	).RequiresReplace()

	parentTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The application resource's parent type.",
	).AllowedValuesEnum(management.AllowedEnumResourceApplicationResourceTypeEnumValues)

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage application resources in PingOne.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to create the resource attribute in."),
			),

			"resource_id": schema.StringAttribute{
				Description:         resourceIdDescription.Description,
				MarkdownDescription: resourceIdDescription.MarkdownDescription,
				Computed:            true,

				CustomType: pingonetypes.ResourceIDType{},

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"resource_name": schema.StringAttribute{
				Description:         resourceNameDescription.Description,
				MarkdownDescription: resourceNameDescription.MarkdownDescription,
				Required:            true,

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			"name": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the application resource name.  The value must be unique for the resource.").Description,
				Required:    true,

				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			"description": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the application resource's description.").Description,
				Optional:    true,
			},

			"parent": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that describes the application resource's parent.").Description,
				Computed:    true,

				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the application resource's parent ID.").Description,
						Computed:    true,

						CustomType: pingonetypes.ResourceIDType{},
					},

					"type": schema.StringAttribute{
						Description:         parentTypeDescription.Description,
						MarkdownDescription: parentTypeDescription.MarkdownDescription,
						Computed:            true,
					},
				},
			},
		},
	}
}

func (r *ApplicationResourceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ApplicationResourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state ApplicationResourceResourceModel

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

	var resourceResponse *management.Resource
	var d diag.Diagnostics
	if !plan.ResourceId.IsNull() && !plan.ResourceId.IsUnknown() {
		resourceResponse, d = fetchResourceFromID(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), plan.ResourceId.ValueString(), false)
	} else {
		resourceResponse, d = fetchResourceFromName(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), plan.ResourceName.ValueString(), false)
	}
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the model for the API
	applicationResource := plan.expand()

	resp.Diagnostics.Append(plan.validate(resourceResponse.GetType())...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var resourceAttributeResponse *management.ResourceApplicationResource
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.ApplicationResourcesApi.CreateApplicationResource(ctx, plan.EnvironmentId.ValueString(), resourceResponse.GetId()).ResourceApplicationResource(*applicationResource).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"CreateApplicationResource",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		&resourceAttributeResponse,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(resourceAttributeResponse, resourceResponse)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *ApplicationResourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *ApplicationResourceResourceModel

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

	var resourceResponse *management.Resource
	var d diag.Diagnostics
	if !data.ResourceId.IsNull() && !data.ResourceId.IsUnknown() {
		resourceResponse, d = fetchResourceFromID(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), data.ResourceId.ValueString(), true)
	} else {
		resourceResponse, d = fetchResourceFromName(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), data.ResourceName.ValueString(), true)
	}

	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Remove from state if resource is not found
	if resourceResponse == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	// Run the API call
	var resourceAttributeResponse *management.ResourceApplicationResource
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.ApplicationResourcesApi.ReadOneApplicationResource(ctx, data.EnvironmentId.ValueString(), data.ResourceId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadOneApplicationResource",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
		&resourceAttributeResponse,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Remove from state if resource is not found
	if resourceAttributeResponse == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(resourceAttributeResponse, resourceResponse)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ApplicationResourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state ApplicationResourceResourceModel

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

	var resourceResponse *management.Resource
	var d diag.Diagnostics
	if !plan.ResourceId.IsNull() && !plan.ResourceId.IsUnknown() {
		resourceResponse, d = fetchResourceFromID(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), plan.ResourceId.ValueString(), false)
	} else {
		resourceResponse, d = fetchResourceFromName(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), plan.ResourceName.ValueString(), false)
	}

	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(plan.validate(resourceResponse.GetType())...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the model for the API
	resourceAttribute := plan.expand()

	// Run the API call
	var resourceAttributeResponse *management.ResourceApplicationResource
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.ApplicationResourcesApi.UpdateApplicationResource(ctx, plan.EnvironmentId.ValueString(), resourceResponse.GetId(), plan.Id.ValueString()).ResourceApplicationResource(*resourceAttribute).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateApplicationResource",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		&resourceAttributeResponse,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(resourceAttributeResponse, resourceResponse)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *ApplicationResourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *ApplicationResourceResourceModel

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

	var resource *management.Resource
	var d diag.Diagnostics
	if !data.ResourceId.IsNull() && !data.ResourceId.IsUnknown() {
		resource, d = fetchResourceFromID(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), data.ResourceId.ValueString(), true)
	} else {
		resource, d = fetchResourceFromName(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), data.ResourceName.ValueString(), true)
	}

	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	if resource == nil {
		return
	}

	// Run the API call
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fR, fErr := r.Client.ManagementAPIClient.ApplicationResourcesApi.DeleteApplicationResource(ctx, data.EnvironmentId.ValueString(), resource.GetId(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, fR, fErr)
		},
		"DeleteApplicationResource",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
		nil,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ApplicationResourceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:  "resource_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:     "application_resource_id",
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

func (p *ApplicationResourceResourceModel) expand() *management.ResourceApplicationResource {

	data := management.NewResourceApplicationResource(
		p.Name.ValueString(),
	)

	if !p.Description.IsNull() {
		data.SetDescription(p.Description.ValueString())
	}

	return data
}

func (p *ApplicationResourceResourceModel) validate(resourceType management.EnumResourceType) diag.Diagnostics {
	var diags diag.Diagnostics

	// Check that we're using a custom resource
	if resourceType != management.ENUMRESOURCETYPE_CUSTOM {
		diags.AddError(
			"Invalid parameter value - Invalid resource type",
			"Resources that are of type PingOne API or OpenID Connect cannot be application resource.  Only custom resources can be application resources.  Please ensure that the resource configured in the `resource_name` parameter is a custom resource.",
		)
	}

	return diags
}

func (p *ApplicationResourceResourceModel) toState(apiObject *management.ResourceApplicationResource, resourceApiObject *management.Resource) diag.Diagnostics {
	var diags, d diag.Diagnostics

	if apiObject == nil || resourceApiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.PingOneResourceIDOkToTF(apiObject.GetIdOk())
	p.ResourceId = framework.PingOneResourceIDOkToTF(resourceApiObject.GetIdOk())
	p.ResourceName = framework.StringOkToTF(resourceApiObject.GetNameOk())
	p.Name = framework.StringOkToTF(apiObject.GetNameOk())
	p.Description = framework.StringOkToTF(apiObject.GetDescriptionOk())
	p.Parent, d = toStateApplicationResourceParentOk(apiObject.GetParentOk())
	diags.Append(d...)

	return diags
}

func toStateApplicationResourceParentOk(apiObject *management.ResourceApplicationResourceParent, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(applicationResourceParentTFObjectTypes), diags
	}

	o := map[string]attr.Value{
		"type": framework.EnumOkToTF(apiObject.GetTypeOk()),
		"id":   framework.PingOneResourceIDOkToTF(apiObject.GetIdOk()),
	}

	returnVar, d := types.ObjectValue(applicationResourceParentTFObjectTypes, o)
	diags.Append(d...)

	return returnVar, diags
}
