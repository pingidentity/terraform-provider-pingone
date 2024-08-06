package authorize

import (
	"context"
	"encoding/json"
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
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
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

type editorServiceCacheSettingsResourceModel struct {
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
	editorServiceParentTFObjectTypes = map[string]attr.Type{
		"id": pingonetypes.ResourceIDType{},
	}

	editorServiceCacheSettingsTFObjectTypes = map[string]attr.Type{
		"ttl_seconds": types.Int32Type,
	}

	editorServiceProcessorTFObjectTypes = map[string]attr.Type{
		"name": types.StringType,
		"type": types.StringType,
	}

	editorServiceValueTypeTFObjectTypes = map[string]attr.Type{
		"type": types.StringType,
	}

	editorServiceServiceSettingsTFObjectTypes = map[string]attr.Type{
		"maximum_concurrent_requests": types.Int32Type,
		"maximum_requests_per_second": types.Float64Type,

		"timeout_milliseconds": types.Int32Type,
		"url":                  types.StringType,
		"verb":                 types.StringType,
		"body":                 types.StringType,
		"content_type":         types.StringType,
		"headers":              types.SetType{ElemType: types.ObjectType{AttrTypes: editorServiceServiceSettingsHeadersTFObjectTypes}},
		"authentication":       types.ObjectType{AttrTypes: editorServiceServiceSettingsAuthenticationTFObjectTypes},
		"tls_settings":         types.ObjectType{AttrTypes: editorServiceServiceSettingsTlsSettingsTFObjectTypes},

		"channel":        types.StringType,
		"code":           types.StringType,
		"capability":     types.StringType,
		"schema_version": types.Int32Type,
		"input_mappings": types.ListType{ElemType: types.ObjectType{AttrTypes: editorServiceServiceSettingsInputMappingsTFObjectTypes}},
	}

	editorServiceServiceSettingsHeadersTFObjectTypes = map[string]attr.Type{
		"key":   types.StringType,
		"value": types.ObjectType{AttrTypes: editorServiceServiceSettingsHeaderValueTFObjectTypes},
	}

	editorServiceServiceSettingsHeaderValueTFObjectTypes = map[string]attr.Type{
		"type": types.StringType,
	}

	editorServiceServiceSettingsAuthenticationTFObjectTypes = map[string]attr.Type{
		"type": types.StringType,
	}

	editorServiceServiceSettingsTlsSettingsTFObjectTypes = map[string]attr.Type{
		"tls_validation_type": types.StringType,
	}

	editorServiceServiceSettingsInputMappingsTFObjectTypes = map[string]attr.Type{
		"property": types.StringType,
		"type":     types.StringType,
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

	commonData, d := p.expandCommon(ctx)
	diags.Append(d...)

	connectorService, d := p.expandConnectorService(ctx, commonData)
	diags.Append(d...)

	httpService, d := p.expandHttpService(ctx, commonData)
	diags.Append(d...)

	noneService, d := p.expandNoneService(ctx, commonData)
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

	commonData, d := p.expandCommon(ctx)
	diags.Append(d...)

	connectorService, d := p.expandConnectorService(ctx, commonData)
	diags.Append(d...)

	httpService, d := p.expandHttpService(ctx, commonData)
	diags.Append(d...)

	noneService, d := p.expandNoneService(ctx, commonData)
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

func (p *editorServiceResourceModel) expandCommon(ctx context.Context) (*authorize.AuthorizeEditorDataDefinitionsServiceDefinitionDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := authorize.NewAuthorizeEditorDataDefinitionsServiceDefinitionDTO(
		p.Name.ValueString(),
		p.ServiceType.ValueString(),
	)

	if !p.FullName.IsNull() && !p.FullName.IsUnknown() {
		data.SetFullName(p.FullName.ValueString())
	}

	if !p.Description.IsNull() && !p.Description.IsUnknown() {
		data.SetDescription(p.Description.ValueString())
	}

	if !p.Version.IsNull() && !p.Version.IsUnknown() {
		data.SetVersion(p.Version.ValueString())
	}

	if !p.Parent.IsNull() && !p.Parent.IsUnknown() {
		var plan *editorServiceParentResourceModel
		diags.Append(p.Parent.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		parent := plan.expand()

		data.SetParent(*parent)
	}

	if !p.Type.IsNull() && !p.Type.IsUnknown() {
		data.SetType(p.Type.ValueString())
	}

	if !p.CacheSettings.IsNull() && !p.CacheSettings.IsUnknown() {
		var plan *editorServiceCacheSettingsResourceModel
		diags.Append(p.CacheSettings.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		cacheSettings := plan.expand()

		data.SetCacheSettings(*cacheSettings)
	}

	return data, diags
}

func (p *editorServiceParentResourceModel) expand() *authorize.AuthorizeEditorDataReferenceObjectDTO {

	data := authorize.NewAuthorizeEditorDataReferenceObjectDTO(p.Id.ValueString())

	return data
}

func (p *editorServiceCacheSettingsResourceModel) expand() *authorize.AuthorizeEditorDataCacheSettingsDTO {

	data := authorize.NewAuthorizeEditorDataCacheSettingsDTO()

	if !p.TtlSeconds.IsNull() && !p.TtlSeconds.IsUnknown() {
		data.SetTtlSeconds(p.TtlSeconds.ValueInt32())
	}

	return data
}

func (p *editorServiceResourceModel) expandConnectorService(ctx context.Context, commonData *authorize.AuthorizeEditorDataDefinitionsServiceDefinitionDTO) (*authorize.AuthorizeEditorDataServicesConnectorServiceDefinitionDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	if p.ServiceType.ValueString() == "CONNECTOR" {
		return nil, diags
	}

	// Use json.marshall and unmarshal to cast commonData to a AuthorizeEditorDataServicesConnectorServiceDefinitionDTO type
	bytes, err := json.Marshal(commonData)
	if err != nil {
		diags.AddError("Failed to marshal data", err.Error())
		return nil, diags
	}

	var data *authorize.AuthorizeEditorDataServicesConnectorServiceDefinitionDTO
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		diags.AddError("Failed to unmarshal data", err.Error())
		return nil, diags
	}

	var valueTypePlan *editorServiceValueTypeResourceModel
	diags.Append(p.ValueType.As(ctx, &valueTypePlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return nil, diags
	}

	valueType := valueTypePlan.expand(ctx)

	var serviceSettingsPlan *editorServiceServiceSettingsResourceModel
	diags.Append(p.ServiceSettings.As(ctx, &serviceSettingsPlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return nil, diags
	}

	serviceSettings, d := serviceSettingsPlan.expandConnector(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	data.SetValueType(*valueType)
	data.SetServiceSettings(*serviceSettings)

	if !p.Processor.IsNull() && !p.Processor.IsUnknown() {
		var plan *editorServiceProcessorResourceModel
		diags.Append(p.Processor.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		processor := plan.expand()

		data.SetProcessor(*processor)
	}

	return data, diags
}

func (p *editorServiceResourceModel) expandHttpService(ctx context.Context, commonData *authorize.AuthorizeEditorDataDefinitionsServiceDefinitionDTO) (*authorize.AuthorizeEditorDataServicesHttpServiceDefinitionDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	if p.ServiceType.ValueString() == "HTTP" {
		return nil, diags
	}

	// Use json.marshall and unmarshal to cast commonData to a AuthorizeEditorDataServicesHttpServiceDefinitionDTO type
	bytes, err := json.Marshal(commonData)
	if err != nil {
		diags.AddError("Failed to marshal data", err.Error())
		return nil, diags
	}

	var data *authorize.AuthorizeEditorDataServicesHttpServiceDefinitionDTO
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		diags.AddError("Failed to unmarshal data", err.Error())
		return nil, diags
	}

	var valueTypePlan *editorServiceValueTypeResourceModel
	diags.Append(p.ValueType.As(ctx, &valueTypePlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return nil, diags
	}

	valueType := valueTypePlan.expand(ctx)

	var serviceSettingsPlan *editorServiceServiceSettingsResourceModel
	diags.Append(p.ServiceSettings.As(ctx, &serviceSettingsPlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return nil, diags
	}

	serviceSettings, d := serviceSettingsPlan.expandHttp(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	data.SetValueType(*valueType)
	data.SetServiceSettings(*serviceSettings)

	if !p.Processor.IsNull() && !p.Processor.IsUnknown() {
		var plan *editorServiceProcessorResourceModel
		diags.Append(p.Processor.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		processor := plan.expand()

		data.SetProcessor(*processor)
	}

	return data, diags
}

func (p *editorServiceResourceModel) expandNoneService(ctx context.Context, commonData *authorize.AuthorizeEditorDataDefinitionsServiceDefinitionDTO) (*authorize.AuthorizeEditorDataServicesNoneServiceDefinitionDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	if p.ServiceType.ValueString() == "NONE" {
		return nil, diags
	}

	// Use json.marshall and unmarshal to cast commonData to a AuthorizeEditorDataServicesNoneServiceDefinitionDTO type
	bytes, err := json.Marshal(commonData)
	if err != nil {
		diags.AddError("Failed to marshal data", err.Error())
		return nil, diags
	}

	var data *authorize.AuthorizeEditorDataServicesNoneServiceDefinitionDTO
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		diags.AddError("Failed to unmarshal data", err.Error())
		return nil, diags
	}

	return data, diags
}

func (p *editorServiceValueTypeResourceModel) expand(ctx context.Context) *authorize.AuthorizeEditorDataValueTypeDTO {

	data := authorize.NewAuthorizeEditorDataValueTypeDTO(
		p.Type.ValueString(),
	)

	return data
}

func (p *editorServiceServiceSettingsResourceModel) expandConnector(ctx context.Context) (*authorize.AuthorizeEditorDataServiceSettingsConnectorServiceSettingsDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	inputMappings := make([]authorize.AuthorizeEditorDataInputMappingDTO, 0)

	var inputMappingsPlan []editorServiceServiceSettingsInputMappingResourceModel
	diags.Append(p.InputMappings.ElementsAs(ctx, &inputMappingsPlan, false)...)
	if diags.HasError() {
		return nil, diags
	}

	for _, inputMappingPlan := range inputMappingsPlan {
		inputMappings = append(inputMappings, *inputMappingPlan.expand())
	}

	data := authorize.NewAuthorizeEditorDataServiceSettingsConnectorServiceSettingsDTO(
		p.Channel.ValueString(),
		p.Code.ValueString(),
		p.Capability.ValueString(),
		inputMappings,
	)

	if !p.MaximumConcurrentRequests.IsNull() && !p.MaximumConcurrentRequests.IsUnknown() {
		data.SetMaximumConcurrentRequests(p.MaximumConcurrentRequests.ValueInt32())
	}

	if !p.MaximumRequestsPerSecond.IsNull() && !p.MaximumRequestsPerSecond.IsUnknown() {
		data.SetMaximumRequestsPerSecond(p.MaximumRequestsPerSecond.ValueFloat64())
	}

	if !p.SchemaVersion.IsNull() && !p.SchemaVersion.IsUnknown() {
		data.SetSchemaVersion(p.SchemaVersion.ValueInt32())
	}

	return data, diags
}

func (p *editorServiceServiceSettingsInputMappingResourceModel) expand() *authorize.AuthorizeEditorDataInputMappingDTO {

	data := authorize.NewAuthorizeEditorDataInputMappingDTO(
		p.Property.ValueString(),
		p.Type.ValueString(),
	)

	return data
}

func (p *editorServiceServiceSettingsResourceModel) expandHttp(ctx context.Context) (*authorize.AuthorizeEditorDataServiceSettingsHttpServiceSettingsDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	var authenticationPlan *editorServiceServiceSettingsAuthenticationResourceModel
	diags.Append(p.Authentication.As(ctx, &authenticationPlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return nil, diags
	}

	authentication := authenticationPlan.expand()

	var tlsSettingsPlan *editorServiceServiceSettingsTlsSettingsResourceModel
	diags.Append(p.TlsSettings.As(ctx, &tlsSettingsPlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return nil, diags
	}

	tlsSettings := tlsSettingsPlan.expand()

	data := authorize.NewAuthorizeEditorDataServiceSettingsHttpServiceSettingsDTO(
		p.Url.ValueString(),
		p.Verb.ValueString(),
		*authentication,
		*tlsSettings,
	)

	if !p.MaximumConcurrentRequests.IsNull() && !p.MaximumConcurrentRequests.IsUnknown() {
		data.SetMaximumConcurrentRequests(p.MaximumConcurrentRequests.ValueInt32())
	}

	if !p.MaximumRequestsPerSecond.IsNull() && !p.MaximumRequestsPerSecond.IsUnknown() {
		data.SetMaximumRequestsPerSecond(p.MaximumRequestsPerSecond.ValueFloat64())
	}

	if !p.TimeoutMilliseconds.IsNull() && !p.TimeoutMilliseconds.IsUnknown() {
		data.SetTimeoutMilliseconds(p.TimeoutMilliseconds.ValueInt32())
	}

	if !p.Body.IsNull() && !p.Body.IsUnknown() {
		data.SetBody(p.Body.ValueString())
	}

	if !p.ContentType.IsNull() && !p.ContentType.IsUnknown() {
		data.SetContentType(p.ContentType.ValueString())
	}

	if !p.Headers.IsNull() && !p.Headers.IsUnknown() {
		headers := make([]authorize.AuthorizeEditorDataHttpRequestHeaderDTO, 0)

		var plan []editorServiceServiceSettingsHeaderResourceModel
		diags.Append(p.Headers.ElementsAs(ctx, &plan, false)...)
		if diags.HasError() {
			return nil, diags
		}

		for _, headerPlan := range plan {
			header, d := headerPlan.expand(ctx)
			diags.Append(d...)
			if diags.HasError() {
				return nil, diags
			}

			headers = append(headers, *header)
		}

		data.SetHeaders(headers)
	}

	return data, diags
}

func (p *editorServiceServiceSettingsHeaderResourceModel) expand(ctx context.Context) (*authorize.AuthorizeEditorDataHttpRequestHeaderDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := authorize.NewAuthorizeEditorDataHttpRequestHeaderDTO(
		p.Key.ValueString(),
	)

	if !p.Value.IsNull() && !p.Value.IsUnknown() {
		var plan *editorServiceServiceSettingsHeaderValueResourceModel
		diags.Append(p.Value.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		value := plan.expand()

		data.SetValue(*value)
	}

	return data, diags
}

func (p *editorServiceServiceSettingsHeaderValueResourceModel) expand() *authorize.AuthorizeEditorDataInputDTO {

	data := authorize.NewAuthorizeEditorDataInputDTO(
		p.Type.ValueString(),
	)

	return data
}

func (p *editorServiceServiceSettingsAuthenticationResourceModel) expand() *authorize.AuthorizeEditorDataAuthenticationDTO {

	data := authorize.NewAuthorizeEditorDataAuthenticationDTO(
		p.Type.ValueString(),
	)

	return data
}

func (p *editorServiceServiceSettingsTlsSettingsResourceModel) expand() *authorize.AuthorizeEditorDataTlsSettingsDTO {

	data := authorize.NewAuthorizeEditorDataTlsSettingsDTO(
		p.TlsValidationType.ValueString(),
	)

	return data
}

func (p *editorServiceProcessorResourceModel) expand() *authorize.AuthorizeEditorDataProcessorDTO {

	data := authorize.NewAuthorizeEditorDataProcessorDTO(
		p.Name.ValueString(),
		p.Type.ValueString(),
	)

	return data
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

	apiObjectCommon := authorize.AuthorizeEditorDataDefinitionsServiceDefinitionDTO{}

	if apiObject.AuthorizeEditorDataServicesConnectorServiceDefinitionDTO != nil {
		apiObjectCommon = authorize.AuthorizeEditorDataDefinitionsServiceDefinitionDTO{
			Id:            apiObject.AuthorizeEditorDataServicesConnectorServiceDefinitionDTO.Id,
			Version:       apiObject.AuthorizeEditorDataServicesConnectorServiceDefinitionDTO.Version,
			Name:          apiObject.AuthorizeEditorDataServicesConnectorServiceDefinitionDTO.Name,
			FullName:      apiObject.AuthorizeEditorDataServicesConnectorServiceDefinitionDTO.FullName,
			Description:   apiObject.AuthorizeEditorDataServicesConnectorServiceDefinitionDTO.Description,
			Parent:        apiObject.AuthorizeEditorDataServicesConnectorServiceDefinitionDTO.Parent,
			Type:          apiObject.AuthorizeEditorDataServicesConnectorServiceDefinitionDTO.Type,
			CacheSettings: apiObject.AuthorizeEditorDataServicesConnectorServiceDefinitionDTO.CacheSettings,
			ServiceType:   apiObject.AuthorizeEditorDataServicesConnectorServiceDefinitionDTO.ServiceType,
		}
	}

	if apiObject.AuthorizeEditorDataServicesHttpServiceDefinitionDTO != nil {
		apiObjectCommon = authorize.AuthorizeEditorDataDefinitionsServiceDefinitionDTO{
			Id:            apiObject.AuthorizeEditorDataServicesHttpServiceDefinitionDTO.Id,
			Version:       apiObject.AuthorizeEditorDataServicesHttpServiceDefinitionDTO.Version,
			Name:          apiObject.AuthorizeEditorDataServicesHttpServiceDefinitionDTO.Name,
			FullName:      apiObject.AuthorizeEditorDataServicesHttpServiceDefinitionDTO.FullName,
			Description:   apiObject.AuthorizeEditorDataServicesHttpServiceDefinitionDTO.Description,
			Parent:        apiObject.AuthorizeEditorDataServicesHttpServiceDefinitionDTO.Parent,
			Type:          apiObject.AuthorizeEditorDataServicesHttpServiceDefinitionDTO.Type,
			CacheSettings: apiObject.AuthorizeEditorDataServicesHttpServiceDefinitionDTO.CacheSettings,
			ServiceType:   apiObject.AuthorizeEditorDataServicesHttpServiceDefinitionDTO.ServiceType,
		}
	}

	if apiObject.AuthorizeEditorDataServicesNoneServiceDefinitionDTO != nil {
		apiObjectCommon = authorize.AuthorizeEditorDataDefinitionsServiceDefinitionDTO{
			Id:            apiObject.AuthorizeEditorDataServicesNoneServiceDefinitionDTO.Id,
			Version:       apiObject.AuthorizeEditorDataServicesNoneServiceDefinitionDTO.Version,
			Name:          apiObject.AuthorizeEditorDataServicesNoneServiceDefinitionDTO.Name,
			FullName:      apiObject.AuthorizeEditorDataServicesNoneServiceDefinitionDTO.FullName,
			Description:   apiObject.AuthorizeEditorDataServicesNoneServiceDefinitionDTO.Description,
			Parent:        apiObject.AuthorizeEditorDataServicesNoneServiceDefinitionDTO.Parent,
			Type:          apiObject.AuthorizeEditorDataServicesNoneServiceDefinitionDTO.Type,
			CacheSettings: apiObject.AuthorizeEditorDataServicesNoneServiceDefinitionDTO.CacheSettings,
			ServiceType:   apiObject.AuthorizeEditorDataServicesNoneServiceDefinitionDTO.ServiceType,
		}
	}

	p.Id = framework.PingOneResourceIDOkToTF(apiObjectCommon.GetIdOk())
	//p.EnvironmentId = framework.PingOneResourceIDToTF(*apiObjectCommon.GetEnvironment().Id)
	p.Name = framework.StringOkToTF(apiObjectCommon.GetNameOk())
	p.FullName = framework.StringOkToTF(apiObjectCommon.GetFullNameOk())
	p.Description = framework.StringOkToTF(apiObjectCommon.GetDescriptionOk())

	p.Parent, d = editorServiceParentOkToTF(apiObjectCommon.GetParentOk())
	diags.Append(d...)

	p.Type = framework.StringOkToTF(apiObjectCommon.GetTypeOk())

	p.CacheSettings, d = editorServiceCacheSettingsOkToTF(apiObjectCommon.GetCacheSettingsOk())
	diags.Append(d...)

	p.ServiceType = framework.StringOkToTF(apiObjectCommon.GetServiceTypeOk())
	p.Version = framework.StringOkToTF(apiObjectCommon.GetVersionOk())

	p.Processor = types.ObjectNull(editorServiceProcessorTFObjectTypes)
	p.ValueType = types.ObjectNull(editorServiceValueTypeTFObjectTypes)
	p.ServiceSettings = types.ObjectNull(editorServiceServiceSettingsTFObjectTypes)

	if v := apiObject.AuthorizeEditorDataServicesConnectorServiceDefinitionDTO; v != nil {
		p.Processor, d = editorServiceProcessorOkToTF(v.GetProcessorOk())
		diags.Append(d...)

		p.ValueType, d = editorServiceValueTypeOkToTF(v.GetValueTypeOk())
		diags.Append(d...)

		p.ServiceSettings, d = editorServiceServiceSettingsConnectorOkToTF(v.GetServiceSettingsOk())
		diags.Append(d...)
	}

	if v := apiObject.AuthorizeEditorDataServicesHttpServiceDefinitionDTO; v != nil {
		p.Processor, d = editorServiceProcessorOkToTF(v.GetProcessorOk())
		diags.Append(d...)

		p.ValueType, d = editorServiceValueTypeOkToTF(v.GetValueTypeOk())
		diags.Append(d...)

		p.ServiceSettings, d = editorServiceServiceSettingsHttpOkToTF(v.GetServiceSettingsOk())
		diags.Append(d...)
	}

	// No implementation for "None" service
	// if apiObject.AuthorizeEditorDataServicesNoneServiceDefinitionDTO != nil {}

	return diags
}

func editorServiceParentOkToTF(apiObject *authorize.AuthorizeEditorDataReferenceObjectDTO, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(editorServiceParentTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(editorServiceParentTFObjectTypes, map[string]attr.Value{
		"id": framework.PingOneResourceIDOkToTF(apiObject.GetIdOk()),
	})
	diags.Append(d...)

	return objValue, diags
}

func editorServiceCacheSettingsOkToTF(apiObject *authorize.AuthorizeEditorDataCacheSettingsDTO, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(editorServiceCacheSettingsTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(editorServiceCacheSettingsTFObjectTypes, map[string]attr.Value{
		"ttl_seconds": framework.Int32OkToTF(apiObject.GetTtlSecondsOk()),
	})
	diags.Append(d...)

	return objValue, diags
}

func editorServiceProcessorOkToTF(apiObject *authorize.AuthorizeEditorDataProcessorDTO, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(editorServiceProcessorTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(editorServiceProcessorTFObjectTypes, map[string]attr.Value{
		"name": framework.StringOkToTF(apiObject.GetNameOk()),
		"type": framework.StringOkToTF(apiObject.GetTypeOk()),
	})
	diags.Append(d...)

	return objValue, diags
}

func editorServiceValueTypeOkToTF(apiObject *authorize.AuthorizeEditorDataValueTypeDTO, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(editorServiceValueTypeTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(editorServiceValueTypeTFObjectTypes, map[string]attr.Value{
		"type": framework.StringOkToTF(apiObject.GetTypeOk()),
	})
	diags.Append(d...)

	return objValue, diags
}

func editorServiceServiceSettingsConnectorOkToTF(apiObject *authorize.AuthorizeEditorDataServiceSettingsConnectorServiceSettingsDTO, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(editorServiceServiceSettingsTFObjectTypes), diags
	}

	inputMappings, d := editorServiceServiceSettingsConnectorInputMappingsOkToTF(apiObject.GetInputMappingsOk())
	diags.Append(d...)

	objValue, d := types.ObjectValue(editorServiceServiceSettingsTFObjectTypes, map[string]attr.Value{
		"maximum_concurrent_requests": framework.Int32OkToTF(apiObject.GetMaximumConcurrentRequestsOk()),
		"maximum_requests_per_second": framework.Float64OkToTF(apiObject.GetMaximumRequestsPerSecondOk()),

		"timeout_milliseconds": types.Int32Null(),
		"url":                  types.StringNull(),
		"verb":                 types.StringNull(),
		"body":                 types.StringNull(),
		"content_type":         types.StringNull(),
		"headers":              types.ObjectNull(editorServiceServiceSettingsHeadersTFObjectTypes),
		"authentication":       types.ObjectNull(editorServiceServiceSettingsAuthenticationTFObjectTypes),
		"tls_settings":         types.ObjectNull(editorServiceServiceSettingsTlsSettingsTFObjectTypes),

		"channel":        framework.StringOkToTF(apiObject.GetChannelOk()),
		"code":           framework.StringOkToTF(apiObject.GetCodeOk()),
		"capability":     framework.StringOkToTF(apiObject.GetCapabilityOk()),
		"schema_version": framework.Int32OkToTF(apiObject.GetSchemaVersionOk()),
		"input_mappings": inputMappings,
	})
	diags.Append(d...)

	return objValue, diags
}

func editorServiceServiceSettingsConnectorInputMappingsOkToTF(apiObject []authorize.AuthorizeEditorDataInputMappingDTO, ok bool) (basetypes.ListValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	tfObjType := types.ObjectType{AttrTypes: editorServiceServiceSettingsInputMappingsTFObjectTypes}

	if !ok || apiObject == nil {
		return types.ListNull(tfObjType), diags
	}

	flattenedList := []attr.Value{}
	for _, v := range apiObject {

		flattenedObj, d := types.ObjectValue(editorServiceServiceSettingsInputMappingsTFObjectTypes, map[string]attr.Value{
			"property": framework.StringOkToTF(v.GetPropertyOk()),
			"type":     framework.StringOkToTF(v.GetTypeOk()),
		})
		diags.Append(d...)

		flattenedList = append(flattenedList, flattenedObj)
	}

	returnVar, d := types.ListValue(tfObjType, flattenedList)
	diags.Append(d...)

	return returnVar, diags
}

func editorServiceServiceSettingsHttpOkToTF(apiObject *authorize.AuthorizeEditorDataServiceSettingsHttpServiceSettingsDTO, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(editorServiceServiceSettingsTFObjectTypes), diags
	}

	headers, d := editorServiceServiceSettingsHttpHeadersOkToTF(apiObject.GetHeadersOk())
	diags.Append(d...)

	authentication, d := editorServiceServiceSettingsHttpAuthenticationOkToTF(apiObject.GetAuthenticationOk())
	diags.Append(d...)

	tlsSettings, d := editorServiceServiceSettingsHttpTlsSettingsOkToTF(apiObject.GetTlsSettingsOk())
	diags.Append(d...)

	objValue, d := types.ObjectValue(editorServiceServiceSettingsTFObjectTypes, map[string]attr.Value{
		"maximum_concurrent_requests": framework.Int32OkToTF(apiObject.GetMaximumConcurrentRequestsOk()),
		"maximum_requests_per_second": framework.Float64OkToTF(apiObject.GetMaximumRequestsPerSecondOk()),

		"timeout_milliseconds": framework.Int32OkToTF(apiObject.GetTimeoutMillisecondsOk()),
		"url":                  framework.StringOkToTF(apiObject.GetUrlOk()),
		"verb":                 framework.StringOkToTF(apiObject.GetVerbOk()),
		"body":                 framework.StringOkToTF(apiObject.GetBodyOk()),
		"content_type":         framework.StringOkToTF(apiObject.GetContentTypeOk()),
		"headers":              headers,
		"authentication":       authentication,
		"tls_settings":         tlsSettings,

		"channel":        types.StringNull(),
		"code":           types.StringNull(),
		"capability":     types.StringNull(),
		"schema_version": types.Int32Null(),
		"input_mappings": types.ListNull(types.ObjectType{AttrTypes: editorServiceServiceSettingsInputMappingsTFObjectTypes}),
	})
	diags.Append(d...)

	return objValue, diags
}

func editorServiceServiceSettingsHttpHeadersOkToTF(apiObject []authorize.AuthorizeEditorDataHttpRequestHeaderDTO, ok bool) (basetypes.SetValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	tfObjType := types.ObjectType{AttrTypes: editorServiceServiceSettingsHeadersTFObjectTypes}

	if !ok || apiObject == nil {
		return types.SetNull(tfObjType), diags
	}

	flattenedList := []attr.Value{}
	for _, v := range apiObject {

		value, d := editorServiceServiceSettingsHttpHeadersValueOkToTF(v.GetValueOk())
		diags.Append(d...)

		flattenedObj, d := types.ObjectValue(editorServiceServiceSettingsHeadersTFObjectTypes, map[string]attr.Value{
			"key":   framework.StringOkToTF(v.GetKeyOk()),
			"value": value,
		})
		diags.Append(d...)

		flattenedList = append(flattenedList, flattenedObj)
	}

	returnVar, d := types.SetValue(tfObjType, flattenedList)
	diags.Append(d...)

	return returnVar, diags
}

func editorServiceServiceSettingsHttpHeadersValueOkToTF(apiObject *authorize.AuthorizeEditorDataInputDTO, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(editorServiceServiceSettingsHeaderValueTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(editorServiceServiceSettingsHeaderValueTFObjectTypes, map[string]attr.Value{
		"type": framework.StringOkToTF(apiObject.GetTypeOk()),
	})
	diags.Append(d...)

	return objValue, diags
}

func editorServiceServiceSettingsHttpAuthenticationOkToTF(apiObject *authorize.AuthorizeEditorDataAuthenticationDTO, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(editorServiceServiceSettingsAuthenticationTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(editorServiceServiceSettingsAuthenticationTFObjectTypes, map[string]attr.Value{
		"type": framework.StringOkToTF(apiObject.GetTypeOk()),
	})
	diags.Append(d...)

	return objValue, diags
}

func editorServiceServiceSettingsHttpTlsSettingsOkToTF(apiObject *authorize.AuthorizeEditorDataTlsSettingsDTO, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(editorServiceServiceSettingsTlsSettingsTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(editorServiceServiceSettingsTlsSettingsTFObjectTypes, map[string]attr.Value{
		"tls_validation_type": framework.StringOkToTF(apiObject.GetTlsValidationTypeOk()),
	})
	diags.Append(d...)

	return objValue, diags
}
