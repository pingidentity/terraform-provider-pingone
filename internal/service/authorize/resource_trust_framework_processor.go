package authorize

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
	"github.com/patrickcping/pingone-go-sdk-v2/authorize"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type TrustFrameworkProcessorResource serviceClientType

type trustFrameworkProcessorResourceModel struct {
	Id            pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	Description   types.String                 `tfsdk:"description"`
	FullName      types.String                 `tfsdk:"full_name"`
	Name          types.String                 `tfsdk:"name"`
	Type          types.String                 `tfsdk:"type"`
	Parent        types.Object                 `tfsdk:"parent"`
	Processor     types.Object                 `tfsdk:"processor"`
	Version       types.String                 `tfsdk:"version"`
}

// Framework interfaces
var (
	_ resource.Resource                = &TrustFrameworkProcessorResource{}
	_ resource.ResourceWithConfigure   = &TrustFrameworkProcessorResource{}
	_ resource.ResourceWithImportState = &TrustFrameworkProcessorResource{}
)

// New Object
func NewTrustFrameworkProcessorResource() resource.Resource {
	return &TrustFrameworkProcessorResource{}
}

// Metadata
func (r *TrustFrameworkProcessorResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_authorize_trust_framework_processor"
}

func (r *TrustFrameworkProcessorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	// schema descriptions and validation settings
	const attrMinLength = 1

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage an authorization processor for the PingOne Authorize Trust Framework in a PingOne environment.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to configure the Authorize editor processor in."),
			),

			"name": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies a user-friendly name to apply to the authorization processor.  The value must be unique.").Description,
				Required:    true,

				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			"type": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that describes the type of the resource.").Description,
				Computed:    true,

				Default: stringdefault.StaticString("PROCESSOR"),
			},

			"full_name": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies a unique name generated by the system for each authorization processor resource. It is the concatenation of names in the processor resource hierarchy.").Description,
				Optional:    true,
			},

			"description": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies a description to apply to the authorization processor.").Description,
				Optional:    true,
			},

			"parent": parentObjectSchema("processor"),

			"processor": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("An object that specifies configuration settings for the authorization processor.").Description,
				Required:    true,

				Attributes: dataProcessorObjectSchemaAttributes(),
			},

			"version": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that describes a random ID generated by the system for concurrency control purposes.").Description,
				Computed:    true,
			},
		},
	}
}

func (r *TrustFrameworkProcessorResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TrustFrameworkProcessorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state trustFrameworkProcessorResourceModel

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
	trustFrameworkProcessor, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *authorize.AuthorizeEditorDataDefinitionsProcessorDefinitionDTO
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.AuthorizeAPIClient.AuthorizeEditorProcessorsApi.CreateProcessor(ctx, plan.EnvironmentId.ValueString()).AuthorizeEditorDataDefinitionsProcessorDefinitionDTO(*trustFrameworkProcessor).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"CreateProcessor",
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
	resp.Diagnostics.Append(state.toState(ctx, response)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *TrustFrameworkProcessorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *trustFrameworkProcessorResourceModel

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
	var response *authorize.AuthorizeEditorDataDefinitionsProcessorDefinitionDTO
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.AuthorizeAPIClient.AuthorizeEditorProcessorsApi.GetProcessor(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"GetProcessor",
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
	resp.Diagnostics.Append(data.toState(ctx, response)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TrustFrameworkProcessorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state trustFrameworkProcessorResourceModel

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
	trustFrameworkProcessor, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *authorize.AuthorizeEditorDataDefinitionsProcessorDefinitionDTO
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.AuthorizeAPIClient.AuthorizeEditorProcessorsApi.UpdateProcessor(ctx, plan.EnvironmentId.ValueString(), plan.Id.ValueString()).AuthorizeEditorDataDefinitionsProcessorDefinitionDTO(*trustFrameworkProcessor).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateProcessor",
		framework.DefaultCustomError,
		nil,
		&response,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(ctx, response)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *TrustFrameworkProcessorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *trustFrameworkProcessorResourceModel

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
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fR, fErr := r.Client.AuthorizeAPIClient.AuthorizeEditorProcessorsApi.DeleteProcessor(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, fR, fErr)
		},
		"DeleteProcessor",
		framework.CustomErrorResourceNotFoundWarning,
		nil,
		nil,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *TrustFrameworkProcessorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:     "authorize_trust_framework_processor_id",
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

func (p *trustFrameworkProcessorResourceModel) expand(ctx context.Context) (*authorize.AuthorizeEditorDataDefinitionsProcessorDefinitionDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	processor, d := expandEditorDataProcessor(ctx, p.Processor)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	// Main object
	data := authorize.NewAuthorizeEditorDataDefinitionsProcessorDefinitionDTO(
		p.Name.ValueString(),
		authorize.EnumAuthorizeEditorDataDefinitionsProcessorDefinitionDTOType(p.Type.ValueString()),
		*processor,
	)

	if !p.Description.IsNull() && !p.Description.IsUnknown() {
		data.SetDescription(p.Description.ValueString())
	}

	if !p.FullName.IsNull() && !p.FullName.IsUnknown() {
		data.SetFullName(p.FullName.ValueString())
	}

	if !p.Parent.IsNull() && !p.Parent.IsUnknown() {
		parent, d := expandEditorParent(ctx, p.Parent)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetParent(*parent)
	}

	if !p.Version.IsNull() && !p.Version.IsUnknown() {
		data.SetVersion(p.Version.ValueString())
	}

	return data, diags
}

func (p *trustFrameworkProcessorResourceModel) toState(ctx context.Context, apiObject *authorize.AuthorizeEditorDataDefinitionsProcessorDefinitionDTO) diag.Diagnostics {
	var diags, d diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)
		return diags
	}

	p.Id = framework.PingOneResourceIDOkToTF(apiObject.GetIdOk())
	p.EnvironmentId = framework.PingOneResourceIDToTF(*apiObject.GetEnvironment().Id)
	p.Name = framework.StringOkToTF(apiObject.GetNameOk())
	p.Type = framework.EnumOkToTF(apiObject.GetTypeOk())
	p.Description = framework.StringOkToTF(apiObject.GetDescriptionOk())
	p.FullName = framework.StringOkToTF(apiObject.GetFullNameOk())

	p.Parent, d = editorParentOkToTF(apiObject.GetParentOk())
	diags.Append(d...)

	processorVal, ok := apiObject.GetProcessorOk()
	p.Processor, d = editorDataProcessorOkToTF(ctx, processorVal, ok)
	diags.Append(d...)

	return diags
}