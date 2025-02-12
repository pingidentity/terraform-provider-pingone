package authorize

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/patrickcping/pingone-go-sdk-v2/authorize"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type TrustFrameworkAttributeResource serviceClientType

type trustFrameworkAttributeResourceModel struct {
	Id               pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId    pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	DefaultValue     types.String                 `tfsdk:"default_value"`
	Description      types.String                 `tfsdk:"description"`
	Type             types.String                 `tfsdk:"type"`
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

type trustFrameworkAttributeResolversConditionResourceModel struct {
	Type types.String `tfsdk:"type"`
}

// Framework interfaces
var (
	_ resource.Resource                = &TrustFrameworkAttributeResource{}
	_ resource.ResourceWithConfigure   = &TrustFrameworkAttributeResource{}
	_ resource.ResourceWithImportState = &TrustFrameworkAttributeResource{}
)

// New Object
func NewTrustFrameworkAttributeResource() resource.Resource {
	return &TrustFrameworkAttributeResource{}
}

// Metadata
func (r *TrustFrameworkAttributeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_authorize_trust_framework_attribute"
}

func (r *TrustFrameworkAttributeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	// schema descriptions and validation settings
	const attrMinLength = 1

	typeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that describes the resource type.",
	).AllowedValuesEnum(authorize.AllowedEnumAuthorizeEditorDataDefinitionsAttributeDefinitionDTOTypeEnumValues)

	managedEntityDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"An object that specifies configuration settings for a system-assigned set of restrictions and metadata related to the resource.",
	)

	valueSchemaDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("A string that specifies the JSON schema defition, where the output type is `%s`.", authorize.ENUMAUTHORIZEEDITORDATAVALUETYPEDTO_JSON),
	)

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage an authorization attribute for the PingOne Authorize Trust Framework in a PingOne environment.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(), // DONE

			"environment_id": framework.Attr_LinkID( // DONE
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to configure the Authorize editor attribute in."),
			),

			"default_value": schema.StringAttribute{ // DONE
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the value to use if no resolvers are defined or if an error occurred with the resolvers or processors.").Description,
				Optional:    true,
			},

			"description": schema.StringAttribute{ // DONE
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies a description to apply to the authorization attribute resource.").Description,
				Optional:    true,
			},

			"full_name": schema.StringAttribute{ // DONE
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies a unique name generated by the system for each attribute resource. It is the concatenation of names in the attribute resource hierarchy.").Description,
				Computed:    true,
			},

			"type": schema.StringAttribute{ // DONE
				Description:         typeDescription.Description,
				MarkdownDescription: typeDescription.MarkdownDescription,
				Computed:            true,

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"managed_entity": schema.SingleNestedAttribute{ // TODO: DOC ERROR - Object Not in docs
				Description:         managedEntityDescription.Description,
				MarkdownDescription: managedEntityDescription.MarkdownDescription,
				Computed:            true,

				Attributes: managedEntityObjectSchemaAttributes(),
			},

			"name": schema.StringAttribute{ // DONE
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies a user-friendly authorization attribute name.  The value must be unique.").Description,
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			"parent": parentObjectSchema("attribute"),

			"processor": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("An object that specifies configuration settings for the authorization attribute resource's processor.").Description,
				Optional:    true,

				Attributes: dataProcessorObjectSchemaAttributes(),
			},

			"repetition_source": repetitionSourceObjectSchema("attribute"),

			"resolvers": schema.ListNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A list of objects that specifies configuration settings for the authorization attribute's resolvers.").Description,
				Optional:    true,

				// TODO: Workaround
				Computed: true,
				Default:  listdefault.StaticValue(types.ListValueMust(types.ObjectType{AttrTypes: editorDataResolverTFObjectTypes}, []attr.Value{})),

				NestedObject: schema.NestedAttributeObject{
					Attributes: dataResolverObjectSchemaAttributes(),
				},
			},

			"value_schema": schema.StringAttribute{
				Description:         valueSchemaDescription.Description,
				MarkdownDescription: valueSchemaDescription.MarkdownDescription,
				Optional:            true,
			},

			"value_type": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("An object that specifies configuration settings for the output value type of the authorization attribute.").Description,
				Required:    true,

				Attributes: valueTypeObjectSchemaAttributes(),
			},

			"version": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that describes a random ID generated by the system for concurrency control purposes.").Description,
				Computed:    true,
			},
		},
	}
}

func (r *TrustFrameworkAttributeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TrustFrameworkAttributeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state trustFrameworkAttributeResourceModel

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
	trustFrameworkAttribute, d := plan.expand(ctx, nil)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *authorize.AuthorizeEditorDataDefinitionsAttributeDefinitionDTO
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.AuthorizeAPIClient.AuthorizeEditorAttributesApi.CreateAttribute(ctx, plan.EnvironmentId.ValueString()).AuthorizeEditorDataDefinitionsAttributeDefinitionDTO(*trustFrameworkAttribute).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"CreateAttribute",
		framework.DefaultCustomError,
		retryAuthorizeEditorCreateUpdate,
		&response,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(ctx, response)...)
	if !resp.Diagnostics.HasError() {
		resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
	}
}

func (r *TrustFrameworkAttributeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *trustFrameworkAttributeResourceModel

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
		nil,
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
	if !resp.Diagnostics.HasError() {
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	}
}

func (r *TrustFrameworkAttributeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state trustFrameworkAttributeResourceModel

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

	// Run the API call
	var getResponse *authorize.AuthorizeEditorDataDefinitionsAttributeDefinitionDTO
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.AuthorizeAPIClient.AuthorizeEditorAttributesApi.GetAttribute(ctx, plan.EnvironmentId.ValueString(), plan.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"GetAttribute-Update",
		framework.DefaultCustomError,
		retryAuthorizeEditorCreateUpdate,
		&getResponse,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	version := getResponse.GetVersion()

	// Build the model for the API
	trustFrameworkAttribute, d := plan.expand(ctx, &version)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *authorize.AuthorizeEditorDataDefinitionsAttributeDefinitionDTO
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.AuthorizeAPIClient.AuthorizeEditorAttributesApi.UpdateAttribute(ctx, plan.EnvironmentId.ValueString(), plan.Id.ValueString()).AuthorizeEditorDataDefinitionsAttributeDefinitionDTO(*trustFrameworkAttribute).Execute()
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
	if !resp.Diagnostics.HasError() {
		resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
	}
}

func (r *TrustFrameworkAttributeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *trustFrameworkAttributeResourceModel

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

	deleteStateConf := &retry.StateChangeConf{
		Pending: []string{
			"200",
		},
		Target: []string{
			"404",
			"ERROR",
		},
		Refresh: func() (interface{}, string, error) {
			// Run the API call
			resp.Diagnostics.Append(framework.ParseResponse(
				ctx,

				func() (any, *http.Response, error) {
					fR, fErr := r.Client.AuthorizeAPIClient.AuthorizeEditorAttributesApi.DeleteAttribute(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
					return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, fR, fErr)
				},
				"DeleteAttribute",
				framework.CustomErrorResourceNotFoundWarning,
				retryAuthorizeEditorDelete,
				nil,
			)...)
			if resp.Diagnostics.HasError() {
				return nil, "ERROR", fmt.Errorf("Error deleting authorize attribute (%s)", data.Id.ValueString())
			}

			fO, fR, fErr := r.Client.AuthorizeAPIClient.AuthorizeEditorAttributesApi.GetAttribute(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			getResp, r, err := framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)

			if err != nil || r == nil {
				return getResp, "ERROR", err
			}

			base := 10
			return getResp, strconv.FormatInt(int64(r.StatusCode), base), nil
		},
		Timeout:                   20 * time.Minute,
		Delay:                     1 * time.Second,
		MinTimeout:                500 * time.Millisecond,
		ContinuousTargetOccurence: 2,
	}
	_, err := deleteStateConf.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Authorize Attribute Delete Timeout",
			fmt.Sprintf("Error waiting for authorize attribute (%s) to be deleted: %s", data.Id.ValueString(), err),
		)

		return
	}
}

func (r *TrustFrameworkAttributeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:     "authorize_trust_framework_attribute_id",
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

func (p *trustFrameworkAttributeResourceModel) expand(ctx context.Context, updateVersionId *string) (*authorize.AuthorizeEditorDataDefinitionsAttributeDefinitionDTO, diag.Diagnostics) {
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

	if updateVersionId != nil {
		data.SetVersion(*updateVersionId)

		if !p.Id.IsNull() && !p.Id.IsUnknown() {
			data.SetId(p.Id.ValueString())
		}
	}

	return data, diags
}

func (p *trustFrameworkAttributeResourceModel) toState(ctx context.Context, apiObject *authorize.AuthorizeEditorDataDefinitionsAttributeDefinitionDTO) diag.Diagnostics {
	var diags, d diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)
		return diags
	}

	p.Id = framework.PingOneResourceIDOkToTF(apiObject.GetIdOk())
	// p.EnvironmentId = framework.PingOneResourceIDToTF(*apiObject.GetEnvironment().Id)
	p.DefaultValue = framework.StringOkToTF(apiObject.GetDefaultValueOk())
	p.Description = framework.StringOkToTF(apiObject.GetDescriptionOk())
	p.Type = framework.EnumOkToTF(apiObject.GetTypeOk())
	p.FullName = framework.StringOkToTF(apiObject.GetFullNameOk())

	p.ManagedEntity, d = editorManagedEntityOkToTF(apiObject.GetManagedEntityOk())
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
