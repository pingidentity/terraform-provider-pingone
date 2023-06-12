package sso

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
)

// Types
type SchemaAttributeResource struct {
	client *management.APIClient
	region model.RegionMapping
}

type SchemaAttributeResourceModel struct {
	Id            types.String `tfsdk:"id"`
	EnvironmentId types.String `tfsdk:"environment_id"`
	SchemaId      types.String `tfsdk:"schema_id"`
	Name          types.String `tfsdk:"name"`
	DisplayName   types.String `tfsdk:"display_name"`
	Description   types.String `tfsdk:"description"`
	Enabled       types.Bool   `tfsdk:"enabled"`
	Type          types.String `tfsdk:"type"`
	SchemaType    types.String `tfsdk:"schema_type"`
	Multivalued   types.Bool   `tfsdk:"multivalued"`
	Unique        types.Bool   `tfsdk:"unique"`
	Required      types.Bool   `tfsdk:"required"`
	LdapAttribute types.String `tfsdk:"ldap_attribute"`
}

// Framework interfaces
var (
	_ resource.Resource                = &SchemaAttributeResource{}
	_ resource.ResourceWithConfigure   = &SchemaAttributeResource{}
	_ resource.ResourceWithImportState = &SchemaAttributeResource{}
)

// New Object
func NewSchemaAttributeResource() resource.Resource {
	return &SchemaAttributeResource{}
}

// Metadata
func (r *SchemaAttributeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_schema_attribute"
}

// Schema.
func (r *SchemaAttributeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	const attrMinLength = 1

	enabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Indicates whether or not the attribute is enabled.",
	).DefaultValue("true")

	typeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The type of the attribute.",
	).AllowedValuesEnum(management.AllowedEnumSchemaAttributeTypeEnumValues).AppendMarkdownString(
		fmt.Sprintf("`%s` and `%s` attributes cannot be created, but standard attributes of those types may be updated. `%s` attributes are limited by size (total size must not exceed 16KB)", string(management.ENUMSCHEMAATTRIBUTETYPE_COMPLEX), string(management.ENUMSCHEMAATTRIBUTETYPE_BOOLEAN), string(management.ENUMSCHEMAATTRIBUTETYPE_JSON)),
	).RequiresReplace().DefaultValue(string(management.ENUMSCHEMAATTRIBUTETYPE_STRING))

	uniqueDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Indicates whether or not the attribute must have a unique value within the PingOne environment.",
	).RequiresReplace().DefaultValue("false")

	multivaluedDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Indicates whether the attribute has multiple values or a single one. Maximum number of values stored is 1,000.",
	).RequiresReplace().DefaultValue("false")

	schemaTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The schema type of the attribute.",
	).AllowedValuesEnum(management.AllowedEnumSchemaAttributeSchemaTypeEnumValues).AppendMarkdownString(
		fmt.Sprintf("`%s` and `%s` attributes are supplied by default. `%s` attributes cannot be updated or deleted. `%s` attributes cannot be deleted, but their mutable properties can be updated. `%s` attributes can be deleted, and their mutable properties can be updated. New attributes are created with a schema type of `%s`.", management.ENUMSCHEMAATTRIBUTESCHEMATYPE_CORE, management.ENUMSCHEMAATTRIBUTESCHEMATYPE_STANDARD, management.ENUMSCHEMAATTRIBUTESCHEMATYPE_CORE, management.ENUMSCHEMAATTRIBUTESCHEMATYPE_STANDARD, management.ENUMSCHEMAATTRIBUTESCHEMATYPE_CUSTOM, management.ENUMSCHEMAATTRIBUTESCHEMATYPE_CUSTOM),
	)

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage PingOne schema attributes.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to create the schema attribute in."),
			),

			"schema_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the schema to apply the schema attribute to."),
			),

			"name": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The system name of the schema attribute.").Description,
				Required:    true,

				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			"display_name": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The display name of the attribute such as 'T-shirt size'. If provided, it must not be an empty string. Valid characters consist of any Unicode letter, mark (for example, accent or umlaut), numeric character, forward slash, dot, apostrophe, underscore, space, or hyphen.").Description,
				Optional:    true,
			},

			"description": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A description of the attribute. If provided, it must not be an empty string. Valid characters consists of any Unicode letter, mark (for example, accent or umlaut), numeric character, punctuation character, or space.").Description,
				Optional:    true,
			},

			"enabled": schema.BoolAttribute{
				Description:         enabledDescription.Description,
				MarkdownDescription: enabledDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,

				Default: booldefault.StaticBool(true),
			},

			"type": schema.StringAttribute{
				Description:         typeDescription.Description,
				MarkdownDescription: typeDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},

				Default: stringdefault.StaticString(string(management.ENUMSCHEMAATTRIBUTETYPE_STRING)),

				Validators: []validator.String{
					stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumSchemaAttributeTypeEnumValues)...),
				},
			},

			"unique": schema.BoolAttribute{
				Description:         uniqueDescription.Description,
				MarkdownDescription: uniqueDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,

				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},

				Default: booldefault.StaticBool(false),
			},

			"multivalued": schema.BoolAttribute{
				Description:         multivaluedDescription.Description,
				MarkdownDescription: multivaluedDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,

				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},

				Default: booldefault.StaticBool(false),
			},

			"required": schema.BoolAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("Indicates whether or not the attribute is required.").Description,
				Computed:    true,
			},

			"ldap_attribute": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The unique identifier for the LDAP attribute.").Description,
				Computed:    true,
			},

			"schema_type": schema.StringAttribute{
				Description:         schemaTypeDescription.Description,
				MarkdownDescription: schemaTypeDescription.MarkdownDescription,
				Computed:            true,
			},
		},
	}
}

func (r *SchemaAttributeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	preparedClient, err := prepareClient(ctx, resourceConfig)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			err.Error(),
		)

		return
	}

	r.client = preparedClient
	r.region = resourceConfig.Client.API.Region
}

func (r *SchemaAttributeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state SchemaAttributeResourceModel

	if r.client == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": r.region.URLSuffix,
	})

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the model for the API
	schemaAttribute, d := plan.expand("CREATE")
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	response, d := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.SchemasApi.CreateAttribute(ctx, plan.EnvironmentId.ValueString(), plan.SchemaId.ValueString()).SchemaAttribute(*schemaAttribute).Execute()
		},
		"CreateAttribute",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)

	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(response.(*management.SchemaAttribute))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *SchemaAttributeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *SchemaAttributeResourceModel

	if r.client == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": r.region.URLSuffix,
	})

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	response, d := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.SchemasApi.ReadOneAttribute(ctx, data.EnvironmentId.ValueString(), data.SchemaId.ValueString(), data.Id.ValueString()).Execute()
		},
		"ReadOneAttribute",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
	)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Remove from state if resource is not found
	if response == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(response.(*management.SchemaAttribute))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SchemaAttributeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state SchemaAttributeResourceModel

	if r.client == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": r.region.URLSuffix,
	})

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the model for the API
	schemaAttribute, d := plan.expand("UPDATE")
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	response, d := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return r.client.SchemasApi.UpdateAttributePatch(ctx, plan.EnvironmentId.ValueString(), plan.SchemaId.ValueString(), plan.Id.ValueString()).SchemaAttribute(*schemaAttribute).Execute()
		},
		"UpdateAttributePatch",
		framework.DefaultCustomError,
		nil,
	)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(response.(*management.SchemaAttribute))...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *SchemaAttributeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *SchemaAttributeResourceModel

	if r.client == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": r.region.URLSuffix,
	})

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	_, d := framework.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			r, err := r.client.SchemasApi.DeleteAttribute(ctx, data.EnvironmentId.ValueString(), data.SchemaId.ValueString(), data.Id.ValueString()).Execute()
			return nil, r, err
		},
		"DeleteAttribute",
		framework.CustomErrorResourceNotFoundWarning,
		nil,
	)
	resp.Diagnostics.Append(d...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *SchemaAttributeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splitLength := 2
	attributes := strings.SplitN(req.ID, "/", splitLength)

	if len(attributes) != splitLength {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("invalid id (\"%s\") specified, should be in format \"environment_id/schema_attribute_id\"", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_id"), attributes[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), attributes[1])...)
}

func (p *SchemaAttributeResourceModel) expand(action string) (*management.SchemaAttribute, diag.Diagnostics) {
	var diags diag.Diagnostics

	attrType := p.Type.ValueString()

	if (attrType == "BOOLEAN" || attrType == "COMPLEX") && action == "CREATE" {
		diags.AddError(
			"Invalid attribute type",
			fmt.Sprintf("Cannot create attributes of type BOOLEAN or COMPLEX.  Custom attributes must be either STRING or JSON.  Attribute type found: %s", attrType),
		)
		return nil, diags
	}

	data := *management.NewSchemaAttribute(p.Enabled.ValueBool(), p.Name.ValueString(), management.EnumSchemaAttributeType(attrType))

	if !p.DisplayName.IsNull() && !p.DisplayName.IsUnknown() {
		data.SetDisplayName(p.DisplayName.ValueString())
	}

	if !p.Description.IsNull() && !p.Description.IsUnknown() {
		data.SetDescription(p.Description.ValueString())
	}

	attrUnique := p.Unique.ValueBool()

	if attrUnique && attrType != "STRING" {
		diags.AddError(
			"Invalid attribute type",
			fmt.Sprintf("Cannot set attribute unique parameter when the attribute type is not STRING.  Attribute type found: %s", attrType),
		)
		return nil, diags
	}

	data.SetUnique(attrUnique)

	data.SetMultiValued(p.Multivalued.ValueBool())

	data.SetRequired(p.Unique.ValueBool())

	return &data, diags
}

func (p *SchemaAttributeResourceModel) toState(apiObject *management.SchemaAttribute) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.StringOkToTF(apiObject.GetIdOk())
	p.EnvironmentId = framework.StringOkToTF(apiObject.Environment.GetIdOk())
	p.SchemaId = framework.StringOkToTF(apiObject.Schema.GetIdOk())
	p.Name = framework.StringOkToTF(apiObject.GetNameOk())
	p.DisplayName = framework.StringOkToTF(apiObject.GetDisplayNameOk())
	p.Description = framework.StringOkToTF(apiObject.GetDescriptionOk())
	p.Enabled = framework.BoolOkToTF(apiObject.GetEnabledOk())
	p.Type = framework.EnumOkToTF(apiObject.GetTypeOk())
	p.Unique = framework.BoolOkToTF(apiObject.GetUniqueOk())
	p.Multivalued = framework.BoolOkToTF(apiObject.GetMultiValuedOk())
	p.Required = framework.BoolOkToTF(apiObject.GetRequiredOk())
	p.LdapAttribute = framework.StringOkToTF(apiObject.GetLdapAttributeOk())
	p.SchemaType = framework.EnumOkToTF(apiObject.GetSchemaTypeOk())

	return diags
}
