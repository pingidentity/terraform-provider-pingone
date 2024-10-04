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
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/patrickcping/pingone-go-sdk-v2/authorize"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type EditorAttributeResource serviceClientType

type editorAttributeResourceModel struct {
	Id               pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId    pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	DefaultValue     types.String                 `tfsdk:"default_value"`
	Description      types.String                 `tfsdk:"description"`
	FullName         types.String                 `tfsdk:"full_name"`
	ManagedEntity    types.Object                 `tfsdk:"managed_entity"`
	Name             types.String                 `tfsdk:"name"`
	Parent           types.Object                 `tfsdk:"parent"`
	Processor        types.Object                 `tfsdk:"processor"`
	RepetitionSource types.Object                 `tfsdk:"repetition_source"`
	Resolvers        types.List                   `tfsdk:"resolvers"`
	ValueType        types.Object                 `tfsdk:"value_type"`
	Version          types.String                 `tfsdk:"version"`
	ValueSchema      types.String                 `tfsdk:"value_schema"`
}

type editorAttributeManagedEntityResourceModel struct {
	Owner        types.Object `tfsdk:"owner"`
	Reference    types.Object `tfsdk:"reference"`
	Restrictions types.Object `tfsdk:"restrictions"`
}

type editorAttributeManagedEntityOwnerResourceModel struct {
	Service types.Object `tfsdk:"service"`
}

type editorAttributeManagedEntityOwnerServiceResourceModel struct {
	Name types.String `tfsdk:"name"`
}

type editorAttributeManagedEntityReferenceResourceModel struct {
	Id         types.String `tfsdk:"id"`
	Type       types.String `tfsdk:"type"`
	Name       types.String `tfsdk:"name"`
	UiDeepLink types.String `tfsdk:"ui_deep_link"`
}

type editorAttributeManagedEntityRestrictionsResourceModel struct {
	ReadOnly         types.Bool `tfsdk:"read_only"`
	DisallowChildren types.Bool `tfsdk:"disallow_children"`
}

type editorAttributeResolversConditionResourceModel struct {
	Type types.String `tfsdk:"type"`
}

var (
	editorAttributeManagedEntityTFObjectTypes = map[string]attr.Type{
		"owner":        types.ObjectType{AttrTypes: editorAttributeManagedEntityOwnerTFObjectTypes},
		"reference":    types.ObjectType{AttrTypes: editorAttributeManagedEntityReferenceTFObjectTypes},
		"restrictions": types.ObjectType{AttrTypes: editorAttributeManagedEntityRestrictionsTFObjectTypes},
	}

	editorAttributeManagedEntityOwnerTFObjectTypes = map[string]attr.Type{
		"service": types.ObjectType{AttrTypes: editorAttributeManagedEntityOwnerServiceTFObjectTypes},
	}

	editorAttributeManagedEntityOwnerServiceTFObjectTypes = map[string]attr.Type{
		"name": types.StringType,
	}

	editorAttributeManagedEntityReferenceTFObjectTypes = map[string]attr.Type{
		"id":           types.StringType,
		"type":         types.StringType,
		"name":         types.StringType,
		"ui_deep_link": types.StringType,
	}

	editorAttributeManagedEntityRestrictionsTFObjectTypes = map[string]attr.Type{
		"read_only":         types.BoolType,
		"disallow_children": types.BoolType,
	}

	editorAttributeResolversConditionTFObjectTypes = map[string]attr.Type{
		"type": types.StringType,
	}

	editorAttributeValueTypeTFObjectTypes = map[string]attr.Type{
		"type": types.StringType,
	}
)

// Framework interfaces
var (
	_ resource.Resource                = &EditorAttributeResource{}
	_ resource.ResourceWithConfigure   = &EditorAttributeResource{}
	_ resource.ResourceWithImportState = &EditorAttributeResource{}
)

// New Object
func NewEditorAttributeResource() resource.Resource {
	return &EditorAttributeResource{}
}

// Metadata
func (r *EditorAttributeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_authorize_editor_attribute"
}

func (r *EditorAttributeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	// schema descriptions and validation settings
	const attrMinLength = 1

	managedEntityDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A system-assigned set of restrictions and metadata related to the resource.",
	)

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage authorization attributes for the PingOne Authorize Trust Framework in a PingOne environment.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(), // DONE

			"environment_id": framework.Attr_LinkID( // DONE
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to configure the Authorize editor attribute in."),
			),

			"default_value": schema.StringAttribute{ // DONE
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The value to use if no resolvers are defined or if an error occurred with the resolvers or processors.").Description,
				Optional:    true,
			},

			"description": schema.StringAttribute{ // DONE
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The attribute resource's description.").Description,
				Optional:    true,
			},

			"full_name": schema.StringAttribute{ // DONE
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A unique name generated by the system for each attribute resource. It is the concatenation of names in the attribute resource hierarchy.").Description,
				Optional:    true,
			},

			"managed_entity": schema.SingleNestedAttribute{ // TODO: DOC ERROR - Object Not in docs
				Description:         managedEntityDescription.Description,
				MarkdownDescription: managedEntityDescription.MarkdownDescription,
				Optional:            true,

				Attributes: map[string]schema.Attribute{
					"owner": schema.SingleNestedAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
						Computed:    true,

						Attributes: map[string]schema.Attribute{
							"service": schema.SingleNestedAttribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
								Computed:    true,

								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
										Computed:    true,
									},
								},
							},
						},
					},

					"reference": schema.SingleNestedAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
						Computed:    true,

						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
								Computed:    true,
							},

							"type": schema.StringAttribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
								Computed:    true,
							},

							"name": schema.StringAttribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
								Computed:    true,
							},

							"ui_deep_link": schema.StringAttribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
								Computed:    true,
							},
						},
					},

					"restrictions": schema.SingleNestedAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
						Computed:    true,

						Attributes: map[string]schema.Attribute{
							"read_only": schema.BoolAttribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
								Computed:    true,
							},

							"disallow_children": schema.BoolAttribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
								Computed:    true,
							},
						},
					},
				},
			},

			"name": schema.StringAttribute{ // DONE
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A user-friendly attribute name.").Description,
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			"parent": parentObjectSchema("attribute"),

			"processor": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("An object that specifies configuration settings for the attribute resource's processor.").Description,
				Optional:    true,

				Attributes: dataProcessorObjectSchemaAttributes(),
			},

			"repetition_source": repetitionSourceObjectSchema("attribute"),

			"resolvers": schema.ListNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The attribute resource's resolvers.").Description,
				Optional:    true,

				NestedObject: schema.NestedAttributeObject{
					Attributes: dataResolverObjectSchemaAttributes(),
				},
			},

			"value_type": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The value type object for the attribute.").Description,
				Required:    true,

				Attributes: valueTypeObjectSchemaAttributes(),
			},

			"version": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A random ID generated by the system for concurrency control purposes.").Description,
				Computed:    true,
			},
		},
	}
}

func (r *EditorAttributeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *EditorAttributeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state editorAttributeResourceModel

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
	editorAttribute, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *authorize.AuthorizeEditorDataDefinitionsAttributeDefinitionDTO
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.AuthorizeAPIClient.AuthorizeEditorAttributesApi.CreateAttribute(ctx, plan.EnvironmentId.ValueString()).AuthorizeEditorDataDefinitionsAttributeDefinitionDTO(*editorAttribute).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"CreateAttribute",
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

func (r *EditorAttributeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *editorAttributeResourceModel

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
	var response *authorize.AuthorizeEditorDataDefinitionsAttributeDefinitionDTO
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.AuthorizeAPIClient.AuthorizeEditorAttributesApi.GetAttribute(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"GetAttribute",
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

func (r *EditorAttributeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state editorAttributeResourceModel

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
	editorAttribute, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *authorize.AuthorizeEditorDataDefinitionsAttributeDefinitionDTO
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.AuthorizeAPIClient.AuthorizeEditorAttributesApi.UpdateAttribute(ctx, plan.EnvironmentId.ValueString(), plan.Id.ValueString()).AuthorizeEditorDataDefinitionsAttributeDefinitionDTO(*editorAttribute).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateAttribute",
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

func (r *EditorAttributeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *editorAttributeResourceModel

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
			fR, fErr := r.Client.AuthorizeAPIClient.AuthorizeEditorAttributesApi.DeleteAttribute(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, fR, fErr)
		},
		"DeleteAttribute",
		framework.CustomErrorResourceNotFoundWarning,
		nil,
		nil,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *EditorAttributeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:     "authorize_editor_attribute_id",
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

func (p *editorAttributeResourceModel) expand(ctx context.Context) (*authorize.AuthorizeEditorDataDefinitionsAttributeDefinitionDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	valueType, d := expandEditorValueType(ctx, p.ValueType)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	// Main object
	data := authorize.NewAuthorizeEditorDataDefinitionsAttributeDefinitionDTO(
		p.Name.ValueString(),
		*valueType,
	)

	if !p.DefaultValue.IsNull() && !p.DefaultValue.IsUnknown() {
		data.SetDefaultValue(p.DefaultValue.ValueString())
	}

	if !p.Description.IsNull() && !p.Description.IsUnknown() {
		data.SetDescription(p.Description.ValueString())
	}

	if !p.FullName.IsNull() && !p.FullName.IsUnknown() {
		data.SetFullName(p.FullName.ValueString())
	}

	if !p.ManagedEntity.IsNull() && !p.ManagedEntity.IsUnknown() {
		var plan *editorAttributeManagedEntityResourceModel
		diags.Append(p.ManagedEntity.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		managedEntity, d := plan.expand(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetManagedEntity(*managedEntity)
	}

	if !p.Parent.IsNull() && !p.Parent.IsUnknown() {
		parent, d := expandEditorParent(ctx, p.Parent)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetParent(*parent)
	}

	if !p.Processor.IsNull() && !p.Processor.IsUnknown() {
		processor, d := expandEditorDataProcessor(ctx, p.Processor)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetProcessor(*processor)
	}

	if !p.RepetitionSource.IsNull() && !p.RepetitionSource.IsUnknown() {
		repetitionSource, d := expandEditorRepetitionSource(ctx, p.RepetitionSource)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetRepetitionSource(*repetitionSource)
	}

	if !p.Resolvers.IsNull() && !p.Resolvers.IsUnknown() {
		var plan []editorDataResolverResourceModel
		diags.Append(p.Resolvers.ElementsAs(ctx, &plan, false)...)
		if diags.HasError() {
			return nil, diags
		}

		resolvers := make([]authorize.AuthorizeEditorDataResolverDTO, 0, len(plan))

		for _, v := range plan {
			resolver, d := v.expand(ctx)
			diags.Append(d...)
			if diags.HasError() {
				return nil, diags
			}

			resolvers = append(resolvers, *resolver)
		}

		data.SetResolvers(resolvers)
	}

	if !p.ValueSchema.IsNull() && !p.ValueSchema.IsUnknown() {
		data.SetValueSchema(p.ValueSchema.ValueString())
	}

	if !p.Version.IsNull() && !p.Version.IsUnknown() {
		data.SetVersion(p.Version.ValueString())
	}

	return data, diags
}

func (p *editorAttributeManagedEntityResourceModel) expand(ctx context.Context) (*authorize.AuthorizeEditorDataManagedEntityDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	var ownerPlan *editorAttributeManagedEntityOwnerResourceModel
	diags.Append(p.Owner.As(ctx, &ownerPlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return nil, diags
	}

	owner, d := ownerPlan.expand(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	data := authorize.NewAuthorizeEditorDataManagedEntityDTO(*owner)

	if !p.Reference.IsNull() && !p.Reference.IsUnknown() {
		var plan *editorAttributeManagedEntityReferenceResourceModel
		diags.Append(p.Reference.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		reference := plan.expand()

		data.SetReference(*reference)
	}

	if !p.Restrictions.IsNull() && !p.Restrictions.IsUnknown() {
		var plan *editorAttributeManagedEntityRestrictionsResourceModel
		diags.Append(p.Restrictions.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		restrictions := plan.expand()

		data.SetRestrictions(*restrictions)
	}

	return data, diags
}

func (p *editorAttributeManagedEntityOwnerResourceModel) expand(ctx context.Context) (*authorize.AuthorizeEditorDataManagedEntityOwnerDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	var servicePlan *editorAttributeManagedEntityOwnerServiceResourceModel
	diags.Append(p.Service.As(ctx, &servicePlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return nil, diags
	}

	service := servicePlan.expand()

	data := authorize.NewAuthorizeEditorDataManagedEntityOwnerDTO(*service)

	return data, diags
}

func (p *editorAttributeManagedEntityOwnerServiceResourceModel) expand() *authorize.AuthorizeEditorDataServiceObjectDTO {

	data := authorize.NewAuthorizeEditorDataServiceObjectDTO(
		p.Name.ValueString(),
	)

	return data
}

func (p *editorAttributeManagedEntityReferenceResourceModel) expand() *authorize.AuthorizeEditorDataManagedEntityManagedEntityReferenceDTO {

	data := authorize.NewAuthorizeEditorDataManagedEntityManagedEntityReferenceDTO()

	if !p.Id.IsNull() && !p.Id.IsUnknown() {
		data.SetId(p.Id.ValueString())
	}

	if !p.Type.IsNull() && !p.Type.IsUnknown() {
		data.SetType(p.Type.ValueString())
	}

	if !p.Name.IsNull() && !p.Name.IsUnknown() {
		data.SetName(p.Name.ValueString())
	}

	if !p.UiDeepLink.IsNull() && !p.UiDeepLink.IsUnknown() {
		data.SetUiDeepLink(p.UiDeepLink.ValueString())
	}

	return data
}

func (p *editorAttributeManagedEntityRestrictionsResourceModel) expand() *authorize.AuthorizeEditorDataManagedEntityRestrictionsDTO {

	data := authorize.NewAuthorizeEditorDataManagedEntityRestrictionsDTO()

	if !p.ReadOnly.IsNull() && !p.ReadOnly.IsUnknown() {
		data.SetReadOnly(p.ReadOnly.ValueBool())
	}

	if !p.DisallowChildren.IsNull() && !p.DisallowChildren.IsUnknown() {
		data.SetDisallowChildren(p.DisallowChildren.ValueBool())
	}

	return data
}

func (p *editorAttributeResourceModel) toState(ctx context.Context, apiObject *authorize.AuthorizeEditorDataDefinitionsAttributeDefinitionDTO) diag.Diagnostics {
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
	p.DefaultValue = framework.StringOkToTF(apiObject.GetDefaultValueOk())
	p.Description = framework.StringOkToTF(apiObject.GetDescriptionOk())
	p.FullName = framework.StringOkToTF(apiObject.GetFullNameOk())

	p.ManagedEntity, d = editorAttributeManagedEntityOkToTF(apiObject.GetManagedEntityOk())
	diags.Append(d...)

	p.Name = framework.StringOkToTF(apiObject.GetNameOk())

	p.Parent, d = editorParentOkToTF(apiObject.GetParentOk())
	diags.Append(d...)

	processor, ok := apiObject.GetProcessorOk()

	p.Processor, d = editorDataProcessorOkToTF(ctx, processor, ok)
	diags.Append(d...)

	p.RepetitionSource, d = editorRepetitionSourceOkToTF(apiObject.GetRepetitionSourceOk())
	diags.Append(d...)

	resolvers, ok := apiObject.GetResolversOk()
	p.Resolvers, d = editorResolversOkToListTF(ctx, resolvers, ok)
	diags.Append(d...)

	p.ValueSchema = framework.StringOkToTF(apiObject.GetValueSchemaOk())

	p.ValueType, d = editorValueTypeOkToTF(apiObject.GetValueTypeOk())
	diags.Append(d...)

	p.Version = framework.StringOkToTF(apiObject.GetVersionOk())

	return diags
}

func editorAttributeManagedEntityOkToTF(apiObject *authorize.AuthorizeEditorDataManagedEntityDTO, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags, d diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(editorAttributeManagedEntityTFObjectTypes), diags
	}

	owner, d := editorAttributeManagedEntityOwnerOkToTF(apiObject.GetOwnerOk())
	diags.Append(d...)

	reference, d := editorAttributeManagedEntityReferenceOkToTF(apiObject.GetReferenceOk())
	diags.Append(d...)

	restrictions, d := editorAttributeManagedEntityRestrictionsOkToTF(apiObject.GetRestrictionsOk())
	diags.Append(d...)

	objValue, d := types.ObjectValue(editorAttributeManagedEntityTFObjectTypes, map[string]attr.Value{
		"owner":        owner,
		"reference":    reference,
		"restrictions": restrictions,
	})
	diags.Append(d...)

	return objValue, diags
}

func editorAttributeManagedEntityOwnerOkToTF(apiObject *authorize.AuthorizeEditorDataManagedEntityOwnerDTO, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(editorAttributeManagedEntityOwnerTFObjectTypes), diags
	}

	service, d := editorAttributeManagedEntityOwnerServiceOkToTF(apiObject.GetServiceOk())
	diags.Append(d...)

	objValue, d := types.ObjectValue(editorAttributeManagedEntityOwnerTFObjectTypes, map[string]attr.Value{
		"service": service,
	})
	diags.Append(d...)

	return objValue, diags
}

func editorAttributeManagedEntityOwnerServiceOkToTF(apiObject *authorize.AuthorizeEditorDataServiceObjectDTO, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(editorAttributeManagedEntityOwnerServiceTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(editorAttributeManagedEntityOwnerServiceTFObjectTypes, map[string]attr.Value{
		"name": framework.StringOkToTF(apiObject.GetNameOk()),
	})
	diags.Append(d...)

	return objValue, diags
}

func editorAttributeManagedEntityReferenceOkToTF(apiObject *authorize.AuthorizeEditorDataManagedEntityManagedEntityReferenceDTO, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(editorAttributeManagedEntityReferenceTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(editorAttributeManagedEntityReferenceTFObjectTypes, map[string]attr.Value{
		"id":           framework.StringOkToTF(apiObject.GetIdOk()),
		"type":         framework.StringOkToTF(apiObject.GetTypeOk()),
		"name":         framework.StringOkToTF(apiObject.GetNameOk()),
		"ui_deep_link": framework.StringOkToTF(apiObject.GetUiDeepLinkOk()),
	})
	diags.Append(d...)

	return objValue, diags
}

func editorAttributeManagedEntityRestrictionsOkToTF(apiObject *authorize.AuthorizeEditorDataManagedEntityRestrictionsDTO, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(editorAttributeManagedEntityRestrictionsTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(editorAttributeManagedEntityRestrictionsTFObjectTypes, map[string]attr.Value{
		"read_only":         framework.BoolOkToTF(apiObject.GetReadOnlyOk()),
		"disallow_children": framework.BoolOkToTF(apiObject.GetDisallowChildrenOk()),
	})
	diags.Append(d...)

	return objValue, diags
}

func editorAttributeResolversConditionOkToTF(apiObject *authorize.AuthorizeEditorDataConditionDTO, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(editorAttributeResolversConditionTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(editorAttributeResolversConditionTFObjectTypes, map[string]attr.Value{
		"type": framework.EnumOkToTF(apiObject.GetTypeOk()),
	})
	diags.Append(d...)

	return objValue, diags
}
