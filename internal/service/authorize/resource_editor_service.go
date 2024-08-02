package authorize

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
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/authorize"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type EditorServiceResource serviceClientType

type editorServiceResourceModel struct {
	Id              pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId   pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	Name            types.String                 `tfsdk:"name"`
	FullName        types.String                 `tfsdk:"full_name"`
	Description     types.String                 `tfsdk:"description"`
	Parent          types.Object                 `tfsdk:"parent"`
	Type            types.String                 `tfsdk:"type"`
	CacheSettings   types.Object                 `tfsdk:"cache_settings"`
	ServiceType     types.String                 `tfsdk:"service_type"`
	Version         types.String                 `tfsdk:"version"`
	Processor       types.Object                 `tfsdk:"processor"`
	ValueType       types.Object                 `tfsdk:"value_type"`
	ServiceSettings types.Object                 `tfsdk:"service_settings"`
}

type editorServiceParentResourceModel editorAttributeReferenceDataResourceModel

type editorServiceCacheSettingResourceModel struct {
	TtlSeconds types.Int32 `tfsdk:"ttl_seconds"`
}

type editorServiceProcessorResourceModel struct {
	Name types.String `tfsdk:"name"`
	Type types.String `tfsdk:"type"`
}

type editorServiceValueTypeResourceModel struct {
	Type types.String `tfsdk:"type"`
}

type editorServiceServiceSettingsResourceModel struct {
	MaximumConcurrentRequests types.Int32   `tfsdk:"maximum_concurrent_requests"`
	MaximumRequestsPerSecond  types.Float64 `tfsdk:"maximum_requests_per_second"`

	// HTTP
	TimeoutMilliseconds types.Int32  `tfsdk:"timeout_milliseconds"`
	Url                 types.String `tfsdk:"url"`
	Verb                types.String `tfsdk:"verb"`
	Body                types.String `tfsdk:"body"`
	ContentType         types.String `tfsdk:"content_type"`
	Headers             types.Set    `tfsdk:"headers"`
	Authentication      types.Object `tfsdk:"authentication"`
	TlsSettings         types.Object `tfsdk:"tls_settings"`

	// Connector
	Channel       types.String `tfsdk:"channel"`
	Code          types.String `tfsdk:"code"`
	Capability    types.String `tfsdk:"capability"`
	SchemaVersion types.Int32  `tfsdk:"schema_version"`
	InputMappings types.List   `tfsdk:"input_mappings"`
}

type editorServiceServiceSettingsHeaderResourceModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.Object `tfsdk:"value"`
}

type editorServiceServiceSettingsHeaderValueResourceModel struct {
	Type types.String `tfsdk:"type"`
}

type editorServiceServiceSettingsAuthenticationResourceModel struct {
	Type types.String `tfsdk:"type"`
}

type editorServiceServiceSettingsTlsSettingsResourceModel struct {
	TlsValidationType types.String `tfsdk:"tls_validation_type"`
}

type editorServiceServiceSettingsInputMappingResourceModel struct {
	Property types.String `tfsdk:"property"`
	Type     types.String `tfsdk:"type"`
}

var (
	editorServiceStatementTFObjectTypes = map[string]attr.Type{}

	editorServiceConditionTFObjectTypes = map[string]attr.Type{
		"type": types.StringType,
	}

	editorServiceEffectSettingsTFObjectTypes = map[string]attr.Type{
		"type": types.StringType,
	}
)

// Framework interfaces
var (
	_ resource.Resource                = &EditorServiceResource{}
	_ resource.ResourceWithConfigure   = &EditorServiceResource{}
	_ resource.ResourceWithImportState = &EditorServiceResource{}
)

// New Object
func NewEditorServiceResource() resource.Resource {
	return &EditorServiceResource{}
}

// Metadata
func (r *EditorServiceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_authorize_editor_service"
}

func (r *EditorServiceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	// schema descriptions and validation settings
	const attrMinLength = 1

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage Authorize editor services in a PingOne environment.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to configure the Authorize editor service in."),
			),

			"name": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
				Required:    true,

				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			"description": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
				Optional:    true,
			},

			"type": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
				Optional:    true,
			},

			"enabled": schema.BoolAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
				Optional:    true,
			},

			"statements": schema.ListNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
				Optional:    true,

				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{},
				},
			},

			"condition": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
				Optional:    true,

				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
						Required:    true,
					},
				},
			},

			"effect_settings": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
				Required:    true,

				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
						Required:    true,
					},
				},
			},

			"version": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
				Optional:    true,
			},
		},
	}
}

func (r *EditorServiceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *EditorServiceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state editorServiceResourceModel

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
	editorService, d := plan.expandCreate(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *authorize.CreateService201Response
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.AuthorizeAPIClient.AuthorizeEditorServicesApi.CreateService(ctx, plan.EnvironmentId.ValueString()).CreateServiceRequest(*editorService).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"CreateService",
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

func (r *EditorServiceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *editorServiceResourceModel

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
	var response *authorize.CreateService201Response
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.AuthorizeAPIClient.AuthorizeEditorServicesApi.GetService(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"GetService",
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

func (r *EditorServiceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state editorServiceResourceModel

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
	editorService, d := plan.expandUpdate(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *authorize.CreateService201Response
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.AuthorizeAPIClient.AuthorizeEditorServicesApi.UpdateService(ctx, plan.EnvironmentId.ValueString(), plan.Id.ValueString()).UpdateServiceRequest(*editorService).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateService",
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
	resp.Diagnostics.Append(state.toState(response)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *EditorServiceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *editorServiceResourceModel

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
			fR, fErr := r.Client.AuthorizeAPIClient.AuthorizeEditorServicesApi.DeleteService(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, fR, fErr)
		},
		"DeleteService",
		framework.CustomErrorResourceNotFoundWarning,
		nil,
		nil,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *EditorServiceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:     "authorize_editor_service_id",
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

func (p *editorServiceResourceModel) expandCreate(ctx context.Context) (*authorize.CreateServiceRequest, diag.Diagnostics) {
	var diags, d diag.Diagnostics

	connectorService, d := p.expandConnectorService(ctx)
	diags.Append(d...)

	httpService, d := p.expandHttpService(ctx)
	diags.Append(d...)

	noneService, d := p.expandNoneService(ctx)
	diags.Append(d...)

	if diags.HasError() {
		return nil, diags
	}

	data := authorize.CreateServiceRequest{
		AuthorizeEditorDataServicesConnectorServiceDefinitionDTO: connectorService,
		AuthorizeEditorDataServicesHttpServiceDefinitionDTO:      httpService,
		AuthorizeEditorDataServicesNoneServiceDefinitionDTO:      noneService,
	}

	return &data, diags
}

func (p *editorServiceResourceModel) expandUpdate(ctx context.Context) (*authorize.UpdateServiceRequest, diag.Diagnostics) {
	var diags, d diag.Diagnostics

	connectorService, d := p.expandConnectorService(ctx)
	diags.Append(d...)

	httpService, d := p.expandHttpService(ctx)
	diags.Append(d...)

	noneService, d := p.expandNoneService(ctx)
	diags.Append(d...)

	if diags.HasError() {
		return nil, diags
	}

	data := authorize.UpdateServiceRequest{
		AuthorizeEditorDataServicesConnectorServiceDefinitionDTO: connectorService,
		AuthorizeEditorDataServicesHttpServiceDefinitionDTO:      httpService,
		AuthorizeEditorDataServicesNoneServiceDefinitionDTO:      noneService,
	}

	return &data, diags
}

func (p *editorServiceResourceModel) toState(apiObject *authorize.CreateService201Response) diag.Diagnostics {
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
	p.Description = framework.StringOkToTF(apiObject.GetDescriptionOk())
	p.Enabled = framework.BoolOkToTF(apiObject.GetEnabledOk())

	p.Statements, d = editorServiceStatementsOkToTF(apiObject.GetStatementsOk())
	diags.Append(d...)

	p.Condition, d = editorServiceConditionOkToTF(apiObject.GetConditionOk())
	diags.Append(d...)

	p.EffectSettings, d = editorServiceEffectSettingsOkToTF(apiObject.GetEffectSettingsOk())
	diags.Append(d...)

	return diags
}
