package authorize

import (
	"context"
	"fmt"
	"log"
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
	ValueSchema      types.String                 `tfsdk:"value_schema"`
	ValueType        types.Object                 `tfsdk:"value_type"`
	Version          types.String                 `tfsdk:"version"`
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

type editorAttributeReferenceDataResourceModel struct {
	Id pingonetypes.ResourceIDValue `tfsdk:"id"`
}

type editorAttributeParentResourceModel editorAttributeReferenceDataResourceModel

type editorAttributeProcessorResourceModel struct {
	Name types.String `tfsdk:"name"`
	Type types.String `tfsdk:"type"`
}

type editorAttributeRepetitionSourceResourceModel editorAttributeReferenceDataResourceModel

type editorAttributeResolversResourceModel struct {
	Condition types.Object `tfsdk:"condition"`
	Name      types.String `tfsdk:"name"`
	Processor types.Object `tfsdk:"processor"`
	Type      types.String `tfsdk:"type"`
}

type editorAttributeResolversConditionResourceModel struct {
	Type types.String `tfsdk:"type"`
}

type editorAttributeResolversProcessorResourceModel struct {
	Name types.String `tfsdk:"name"`
	Type types.String `tfsdk:"type"`
}

type editorAttributeValueTypeResourceModel struct {
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

	editorAttributeReferenceObjectTFObjectTypes = map[string]attr.Type{
		"id": pingonetypes.ResourceIDType{},
	}

	editorAttributeProcessorTFObjectTypes = map[string]attr.Type{
		"name": types.StringType,
		"type": types.StringType,
	}

	editorAttributeResolversTFObjectTypes = map[string]attr.Type{
		"condition": types.ObjectType{AttrTypes: editorAttributeResolversConditionTFObjectTypes},
		"name":      types.StringType,
		"processor": types.ObjectType{AttrTypes: editorAttributeResolversProcessorTFObjectTypes},
		"type":      types.StringType,
	}

	editorAttributeResolversConditionTFObjectTypes = map[string]attr.Type{
		"type": types.StringType,
	}

	editorAttributeResolversProcessorTFObjectTypes = map[string]attr.Type{
		"name": types.StringType,
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

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage Authorize editor attributes in a PingOne environment.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to configure the Authorize editor attribute in."),
			),

			"default_value": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
				Optional:    true,
			},

			"description": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
				Optional:    true,
			},

			"full_name": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
				Optional:    true,
			},

			"managed_entity": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
				Optional:    true,

				Attributes: map[string]schema.Attribute{
					"owner": schema.SingleNestedAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
						Required:    true,

						Attributes: map[string]schema.Attribute{
							"service": schema.SingleNestedAttribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
								Required:    true,

								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
										Required:    true,
									},
								},
							},
						},
					},

					"reference": schema.SingleNestedAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
						Optional:    true,

						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
								Optional:    true,
							},

							"type": schema.StringAttribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
								Optional:    true,
							},

							"name": schema.StringAttribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
								Optional:    true,
							},

							"ui_deep_link": schema.StringAttribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
								Optional:    true,
							},
						},
					},

					"restrictions": schema.SingleNestedAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
						Optional:    true,

						Attributes: map[string]schema.Attribute{
							"read_only": schema.BoolAttribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
								Optional:    true,
							},

							"disallow_children": schema.BoolAttribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
								Optional:    true,
							},
						},
					},
				},
			},

			"name": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			"parent": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
				Optional:    true,

				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
						Required:    true,

						CustomType: pingonetypes.ResourceIDType{},
					},
				},
			},

			"processor": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
				Optional:    true,

				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
						Required:    true,
					},

					"type": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
						Required:    true,
					},
				},
			},

			"repetition_source": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
				Optional:    true,

				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
						Required:    true,

						CustomType: pingonetypes.ResourceIDType{},
					},
				},
			},

			"resolvers": schema.ListNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
				Optional:    true,

				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"condition": schema.ListNestedAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
							Optional:    true,

							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"type": schema.StringAttribute{
										Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
										Required:    true,
									},
								},
							},
						},

						"name": schema.StringAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
							Optional:    true,
						},

						"processor": schema.ListNestedAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
							Optional:    true,

							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
										Required:    true,
									},

									"type": schema.StringAttribute{
										Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
										Required:    true,
									},
								},
							},
						},

						"type": schema.StringAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
							Required:    true,
						},
					},
				},
			},

			"value_schema": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
				Optional:    true,
			},

			"value_type": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
				Optional:    true,

				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
						Required:    true,

						//TODO: Enum validation
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
	resp.Diagnostics.Append(state.toState(response)...)
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
	resp.Diagnostics.Append(data.toState(response)...)
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
	resp.Diagnostics.Append(state.toState(response)...)
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

	var valueTypePlan *editorAttributeValueTypeResourceModel
	diags.Append(p.ValueType.As(ctx, &valueTypePlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return nil, diags
	}

	valueType := authorize.NewAuthorizeEditorDataValueTypeDTO(valueTypePlan.Type.ValueString())

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
		var plan *editorAttributeParentResourceModel
		diags.Append(p.Parent.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		parent, d := plan.expand(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetParent(*parent)
	}

	if !p.Processor.IsNull() && !p.Processor.IsUnknown() {
		var plan *editorAttributeProcessorResourceModel
		diags.Append(p.Processor.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		processor, d := plan.expand(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetProcessor(*processor)
	}

	if !p.RepetitionSource.IsNull() && !p.RepetitionSource.IsUnknown() {
		var plan *editorAttributeRepetitionSourceResourceModel
		diags.Append(p.RepetitionSource.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		repetitionSource, d := plan.expand(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		data.SetRepetitionSource(*repetitionSource)
	}

	if !p.Resolvers.IsNull() && !p.Resolvers.IsUnknown() {
		var plan []editorAttributeResolversResourceModel
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

		reference := plan.expand(ctx)

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

		restrictions := plan.expand(ctx)

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

	service := servicePlan.expand(ctx)

	data := authorize.NewAuthorizeEditorDataManagedEntityOwnerDTO(*service)

	return data, diags
}

func (p *editorAttributeManagedEntityOwnerServiceResourceModel) expand(ctx context.Context) *authorize.AuthorizeEditorDataServiceObjectDTO {

	data := authorize.NewAuthorizeEditorDataServiceObjectDTO(
		p.Name.ValueString(),
	)

	return data
}

func (p *editorAttributeManagedEntityReferenceResourceModel) expand(ctx context.Context) *authorize.AuthorizeEditorDataManagedEntityManagedEntityReferenceDTO {

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

func (p *editorAttributeManagedEntityRestrictionsResourceModel) expand(ctx context.Context) *authorize.AuthorizeEditorDataManagedEntityRestrictionsDTO {

	data := authorize.NewAuthorizeEditorDataManagedEntityRestrictionsDTO()

	if !p.ReadOnly.IsNull() && !p.ReadOnly.IsUnknown() {
		data.SetReadOnly(p.ReadOnly.ValueBool())
	}

	if !p.DisallowChildren.IsNull() && !p.DisallowChildren.IsUnknown() {
		data.SetDisallowChildren(p.DisallowChildren.ValueBool())
	}

	return data
}

func (p *editorAttributeParentResourceModel) expand(ctx context.Context) (*authorize.AuthorizeEditorDataReferenceObjectDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	log.Fatalf("Not implemented")

	return nil, diags
}

func (p *editorAttributeProcessorResourceModel) expand(ctx context.Context) (*authorize.AuthorizeEditorDataProcessorDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	log.Fatalf("Not implemented")

	return nil, diags
}

func (p *editorAttributeRepetitionSourceResourceModel) expand(ctx context.Context) (*authorize.AuthorizeEditorDataReferenceObjectDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	log.Fatalf("Not implemented")

	return nil, diags
}

func (p *editorAttributeResolversResourceModel) expand(ctx context.Context) (*authorize.AuthorizeEditorDataResolverDTO, diag.Diagnostics) {
	var diags diag.Diagnostics

	log.Fatalf("Not implemented")

	return nil, diags
}

func (p *editorAttributeResourceModel) toState(apiObject *authorize.AuthorizeEditorDataDefinitionsAttributeDefinitionDTO) diag.Diagnostics {
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
	diags = append(diags, d...)

	p.Name = framework.StringOkToTF(apiObject.GetNameOk())

	p.Parent, d = editorAttributeParentOkToTF(apiObject.GetParentOk())
	diags = append(diags, d...)

	p.Processor, d = editorAttributeProcessorOkToTF(apiObject.GetProcessorOk())
	diags = append(diags, d...)

	p.RepetitionSource, d = editorAttributeRepetitionSourceOkToTF(apiObject.GetRepetitionSourceOk())
	diags = append(diags, d...)

	p.Resolvers, d = editorAttributeResolversOkToTF(apiObject.GetResolversOk())
	diags = append(diags, d...)

	p.ValueSchema = framework.StringOkToTF(apiObject.GetValueSchemaOk())

	p.ValueType, d = editorAttributeValueTypeOkToTF(apiObject.GetValueTypeOk())
	diags = append(diags, d...)

	p.Version = framework.StringOkToTF(apiObject.GetVersionOk())

	return diags
}

func editorAttributeManagedEntityOkToTF(apiObject *authorize.AuthorizeEditorDataManagedEntityDTO, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags, d diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(editorAttributeManagedEntityTFObjectTypes), diags
	}

	owner, d := editorAttributeManagedEntityOwnerOkToTF(apiObject.GetOwnerOk())
	diags = append(diags, d...)

	reference, d := editorAttributeManagedEntityReferenceOkToTF(apiObject.GetReferenceOk())
	diags = append(diags, d...)

	restrictions, d := editorAttributeManagedEntityRestrictionsOkToTF(apiObject.GetRestrictionsOk())
	diags = append(diags, d...)

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
	diags = append(diags, d...)

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

func editorAttributeDataReferenceObjectOkToTF(apiObject *authorize.AuthorizeEditorDataReferenceObjectDTO, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(editorAttributeReferenceObjectTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(editorAttributeReferenceObjectTFObjectTypes, map[string]attr.Value{
		"id": framework.PingOneResourceIDOkToTF(apiObject.GetIdOk()),
	})
	diags.Append(d...)

	return objValue, diags
}

func editorAttributeParentOkToTF(apiObject *authorize.AuthorizeEditorDataReferenceObjectDTO, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	return editorAttributeDataReferenceObjectOkToTF(apiObject, ok)
}

func editorAttributeProcessorOkToTF(apiObject *authorize.AuthorizeEditorDataProcessorDTO, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(editorAttributeProcessorTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(editorAttributeProcessorTFObjectTypes, map[string]attr.Value{
		"name": framework.StringOkToTF(apiObject.GetNameOk()),
		"type": framework.StringOkToTF(apiObject.GetTypeOk()),
	})
	diags.Append(d...)

	return objValue, diags
}

func editorAttributeRepetitionSourceOkToTF(apiObject *authorize.AuthorizeEditorDataReferenceObjectDTO, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	return editorAttributeDataReferenceObjectOkToTF(apiObject, ok)
}

func editorAttributeResolversOkToTF(apiObject []authorize.AuthorizeEditorDataResolverDTO, ok bool) (basetypes.ListValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	tfObjType := types.ObjectType{AttrTypes: editorAttributeResolversTFObjectTypes}

	if !ok || apiObject == nil {
		return types.ListNull(tfObjType), diags
	}

	flattenedList := []attr.Value{}
	for _, v := range apiObject {

		condition, d := editorAttributeResolversConditionOkToTF(v.GetConditionOk())
		diags = append(diags, d...)

		processor, d := editorAttributeResolversProcessorOkToTF(v.GetProcessorOk())
		diags = append(diags, d...)

		flattenedObj, d := types.ObjectValue(editorAttributeResolversTFObjectTypes, map[string]attr.Value{
			"condition": condition,
			"name":      framework.StringOkToTF(v.GetNameOk()),
			"processor": processor,
			"type":      framework.StringOkToTF(v.GetTypeOk()),
		})
		diags.Append(d...)

		flattenedList = append(flattenedList, flattenedObj)
	}

	returnVar, d := types.ListValue(tfObjType, flattenedList)
	diags.Append(d...)

	return returnVar, diags
}

func editorAttributeResolversConditionOkToTF(apiObject *authorize.AuthorizeEditorDataConditionDTO, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(editorAttributeResolversConditionTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(editorAttributeResolversConditionTFObjectTypes, map[string]attr.Value{
		"type": framework.StringOkToTF(apiObject.GetTypeOk()),
	})
	diags.Append(d...)

	return objValue, diags
}

func editorAttributeResolversProcessorOkToTF(apiObject *authorize.AuthorizeEditorDataProcessorDTO, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(editorAttributeResolversProcessorTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(editorAttributeResolversProcessorTFObjectTypes, map[string]attr.Value{
		"name": framework.StringOkToTF(apiObject.GetNameOk()),
		"type": framework.StringOkToTF(apiObject.GetTypeOk()),
	})
	diags.Append(d...)

	return objValue, diags
}

func editorAttributeValueTypeOkToTF(apiObject *authorize.AuthorizeEditorDataValueTypeDTO, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(editorAttributeValueTypeTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(editorAttributeValueTypeTFObjectTypes, map[string]attr.Value{
		"type": framework.StringOkToTF(apiObject.GetTypeOk()),
	})
	diags.Append(d...)

	return objValue, diags
}
