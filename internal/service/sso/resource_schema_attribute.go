package sso

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	boolvalidatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/boolvalidator"
	objectvalidatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/objectvalidator"
	setplanmodifierinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/setplanmodifier"
	setvalidatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/setvalidator"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type SchemaAttributeResource serviceClientType

type SchemaAttributeResourceModelV1 struct {
	Id               types.String `tfsdk:"id"`
	EnvironmentId    types.String `tfsdk:"environment_id"`
	Description      types.String `tfsdk:"description"`
	DisplayName      types.String `tfsdk:"display_name"`
	Enabled          types.Bool   `tfsdk:"enabled"`
	EnumeratedValues types.Set    `tfsdk:"enumerated_values"`
	LdapAttribute    types.String `tfsdk:"ldap_attribute"`
	Multivalued      types.Bool   `tfsdk:"multivalued"`
	Name             types.String `tfsdk:"name"`
	RegexValidation  types.Object `tfsdk:"regex_validation"`
	Required         types.Bool   `tfsdk:"required"`
	SchemaId         types.String `tfsdk:"schema_id"`
	SchemaName       types.String `tfsdk:"schema_name"`
	SchemaType       types.String `tfsdk:"schema_type"`
	Type             types.String `tfsdk:"type"`
	Unique           types.Bool   `tfsdk:"unique"`
}

type SchemaAttributeEnumeratedValuesResourceModel struct {
	Archived    types.Bool   `tfsdk:"archived"`
	Description types.String `tfsdk:"description"`
	Value       types.String `tfsdk:"value"`
}

type SchemaAttributeRegexValidationModel struct {
	Pattern                     types.String `tfsdk:"pattern"`
	Requirements                types.String `tfsdk:"requirements"`
	ValuesPatternShouldMatch    types.Set    `tfsdk:"values_pattern_should_match"`
	ValuesPatternShouldNotMatch types.Set    `tfsdk:"values_pattern_should_not_match"`
}

var (
	schemaAttributeEnumeratedValuesTFObjectTypes = map[string]attr.Type{
		"archived":    types.BoolType,
		"description": types.StringType,
		"value":       types.StringType,
	}

	schemaAttributeRegexValidationTFObjectTypes = map[string]attr.Type{
		"pattern":                         types.StringType,
		"requirements":                    types.StringType,
		"values_pattern_should_match":     types.SetType{ElemType: types.StringType},
		"values_pattern_should_not_match": types.SetType{ElemType: types.StringType},
	}
)

// Framework interfaces
var (
	_ resource.Resource                 = &SchemaAttributeResource{}
	_ resource.ResourceWithConfigure    = &SchemaAttributeResource{}
	_ resource.ResourceWithImportState  = &SchemaAttributeResource{}
	_ resource.ResourceWithUpgradeState = &SchemaAttributeResource{}
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
	const schemaName = "User"

	schemaIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The ID of the schema the schema attribute is applied to.",
	)

	schemaNameDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The name of the schema to apply the schema attribute to.",
	).AllowedValues(schemaName).DefaultValue(schemaName)

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

	enumeratedValuesDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A set of one or more enumerated values for the attribute. If provided, it must not be an empty set.  Can only be set where the attribute type is `STRING` and cannot be set alongside `regex_validation`.  If the attribute has been created without enumerated values and this parameter is added later, this will trigger a replacement plan of the attribute resource.  If the attribute has been created with enumerated values that are subsequently removed, this will update without needing to replace the attribute resource.",
	)

	regexValidationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object representation of the optional regular expression representation of this attribute.  Can only be set where the attribute type is `STRING` and cannot be set alongside `enumerated_values`.",
	)

	resp.Schema = schema.Schema{

		Version: 1,

		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage PingOne schema attributes.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to create the schema attribute in."),
			),

			"schema_id": schema.StringAttribute{
				Description:         schemaIdDescription.Description,
				MarkdownDescription: schemaIdDescription.MarkdownDescription,
				Computed:            true,

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"schema_name": schema.StringAttribute{
				Description:         schemaNameDescription.Description,
				MarkdownDescription: schemaNameDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,

				Default: stringdefault.StaticString(schemaName),

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},

				Validators: []validator.String{
					stringvalidator.OneOf(schemaName),
				},
			},

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

				Validators: []validator.Bool{
					boolvalidatorinternal.MustNotBeTrueIfPathSetToValue(
						types.StringValue(string(management.ENUMSCHEMAATTRIBUTETYPE_JSON)),
						path.MatchRoot("type"),
					),
					boolvalidatorinternal.MustNotBeTrueIfPathSetToValue(
						types.StringValue(string(management.ENUMSCHEMAATTRIBUTETYPE_COMPLEX)),
						path.MatchRoot("type"),
					),
					boolvalidatorinternal.MustNotBeTrueIfPathSetToValue(
						types.StringValue(string(management.ENUMSCHEMAATTRIBUTETYPE_BOOLEAN)),
						path.MatchRoot("type"),
					),
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

			"enumerated_values": schema.SetNestedAttribute{
				Description:         enumeratedValuesDescription.Description,
				MarkdownDescription: enumeratedValuesDescription.MarkdownDescription,
				Optional:            true,

				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"value": schema.StringAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the value of the enumerated value item. If provided, it must not be an empty string.").Description,
							Required:    true,
						},

						"archived": schema.BoolAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether the enumerated value is archived. Archived values cannot be added to a user, but existing archived values are preserved. This allows clients that read the schema to know all possible values of an attribute.").Description,
							Optional:    true,
							Computed:    true,

							Default: booldefault.StaticBool(false),
						},

						"description": schema.StringAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the description of the enumerated value.").Description,
							Optional:    true,
						},
					},
				},

				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplaceIf(
						setplanmodifierinternal.RequiresReplaceIfPreviouslyNull(),
						"The attribute has been previously created without enumerated values validation.  To add enumerated values validation, the attribute must be replaced.",
						"The attribute has been previously created without enumerated values validation.  To add enumerated values validation, the attribute must be replaced.",
					),
				},

				Validators: []validator.Set{
					setvalidator.SizeAtLeast(attrMinLength),
					setvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("regex_validation")),
					setvalidatorinternal.ConflictsIfMatchesPathValue(types.StringValue(string(management.ENUMSCHEMAATTRIBUTETYPE_JSON)), path.MatchRelative().AtParent().AtName("type")),
					setvalidatorinternal.ConflictsIfMatchesPathValue(types.StringValue(string(management.ENUMSCHEMAATTRIBUTETYPE_BOOLEAN)), path.MatchRelative().AtParent().AtName("type")),
					setvalidatorinternal.ConflictsIfMatchesPathValue(types.StringValue(string(management.ENUMSCHEMAATTRIBUTETYPE_COMPLEX)), path.MatchRelative().AtParent().AtName("type")),
				},
			},

			"regex_validation": schema.SingleNestedAttribute{
				Description:         regexValidationDescription.Description,
				MarkdownDescription: regexValidationDescription.MarkdownDescription,
				Optional:            true,

				Attributes: map[string]schema.Attribute{
					"pattern": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the regular expression to which the attribute must conform.").Description,
						Required:    true,
					},

					"requirements": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies a developer friendly description of the regular expression requirements.").Description,
						Required:    true,
					},

					"values_pattern_should_match": schema.SetAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A set of one or more strings matching the regular expression.").Description,
						Optional:    true,

						ElementType: types.StringType,
					},

					"values_pattern_should_not_match": schema.SetAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A set of one or more strings that do not match the regular expression.").Description,
						Optional:    true,

						ElementType: types.StringType,
					},
				},

				Validators: []validator.Object{
					objectvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("enumerated_values")),
					objectvalidatorinternal.ConflictsIfMatchesPathValue(types.StringValue(string(management.ENUMSCHEMAATTRIBUTETYPE_JSON)), path.MatchRelative().AtParent().AtName("type")),
					objectvalidatorinternal.ConflictsIfMatchesPathValue(types.StringValue(string(management.ENUMSCHEMAATTRIBUTETYPE_BOOLEAN)), path.MatchRelative().AtParent().AtName("type")),
					objectvalidatorinternal.ConflictsIfMatchesPathValue(types.StringValue(string(management.ENUMSCHEMAATTRIBUTETYPE_COMPLEX)), path.MatchRelative().AtParent().AtName("type")),
				},
			},

			"required": schema.BoolAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("Indicates whether or not the attribute is required.").Description,
				Computed:    true,

				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},

			"ldap_attribute": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The unique identifier for the LDAP attribute.").Description,
				Computed:    true,

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"schema_type": schema.StringAttribute{
				Description:         schemaTypeDescription.Description,
				MarkdownDescription: schemaTypeDescription.MarkdownDescription,
				Computed:            true,

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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

	r.Client = resourceConfig.Client.API
	if r.Client == nil {
		resp.Diagnostics.AddError(
			"Client not initialised",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.",
		)
		return
	}
}

func (r *SchemaAttributeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state SchemaAttributeResourceModelV1

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

	// Get the schema ID
	schema, d := fetchSchemaFromName(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), plan.SchemaName.ValueString())
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the model for the API
	schemaAttribute, d := plan.expand(ctx, "CREATE")
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *management.SchemaAttribute
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.SchemasApi.CreateAttribute(ctx, plan.EnvironmentId.ValueString(), schema.GetId()).SchemaAttribute(*schemaAttribute).Execute()
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
	resp.Diagnostics.Append(state.toState(response, schema)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *SchemaAttributeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *SchemaAttributeResourceModelV1

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

	// Run the API call
	var response *management.SchemaAttribute
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.SchemasApi.ReadOneAttribute(ctx, data.EnvironmentId.ValueString(), data.SchemaId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadOneAttribute",
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

	var schemaResponse *management.Schema
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.SchemasApi.ReadOneSchema(ctx, data.EnvironmentId.ValueString(), data.SchemaId.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadOneSchema",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
		&schemaResponse,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Remove from state if resource is not found
	if schemaResponse == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(response, schemaResponse)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SchemaAttributeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state SchemaAttributeResourceModelV1

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

	// Get the schema ID
	schema, d := fetchSchemaFromName(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), plan.SchemaName.ValueString())
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the model for the API
	schemaAttribute, d := plan.expand(ctx, "UPDATE")
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *management.SchemaAttribute
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.SchemasApi.UpdateAttributePut(ctx, plan.EnvironmentId.ValueString(), schema.GetId(), plan.Id.ValueString()).SchemaAttribute(*schemaAttribute).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateAttributePut",
		framework.DefaultCustomError,
		nil,
		&response,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(response, schema)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *SchemaAttributeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *SchemaAttributeResourceModelV1

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

	// Run the API call
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fR, fErr := r.Client.ManagementAPIClient.SchemasApi.DeleteAttribute(ctx, data.EnvironmentId.ValueString(), data.SchemaId.ValueString(), data.Id.ValueString()).Execute()
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

func (r *SchemaAttributeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:  "schema_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:     "schema_attribute_id",
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

func (p *SchemaAttributeResourceModelV1) expand(ctx context.Context, action string) (*management.SchemaAttribute, diag.Diagnostics) {
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

	data.SetSchemaType(management.ENUMSCHEMAATTRIBUTESCHEMATYPE_CUSTOM)

	if !p.DisplayName.IsNull() && !p.DisplayName.IsUnknown() {
		data.SetDisplayName(p.DisplayName.ValueString())
	}

	if !p.Description.IsNull() && !p.Description.IsUnknown() {
		data.SetDescription(p.Description.ValueString())
	}

	attrUnique := p.Unique.ValueBool()

	// This is handled in schema validation, but we optionally check here too
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

	if !p.EnumeratedValues.IsNull() && !p.EnumeratedValues.IsUnknown() {
		var plan []SchemaAttributeEnumeratedValuesResourceModel
		diags.Append(p.EnumeratedValues.ElementsAs(ctx, &plan, false)...)
		if diags.HasError() {
			return nil, diags
		}

		enumeratedValues := make([]management.SchemaAttributeEnumeratedValuesInner, 0)
		for _, v := range plan {
			enumeratedValue := management.NewSchemaAttributeEnumeratedValuesInner(v.Value.ValueString())

			if !v.Archived.IsNull() && !v.Archived.IsUnknown() {
				enumeratedValue.SetArchived(v.Archived.ValueBool())
			}

			if !v.Description.IsNull() && !v.Description.IsUnknown() {
				enumeratedValue.SetDescription(v.Description.ValueString())
			}

			enumeratedValues = append(enumeratedValues, *enumeratedValue)
		}

		data.SetEnumeratedValues(enumeratedValues)
	}

	if !p.RegexValidation.IsNull() && !p.RegexValidation.IsUnknown() {

		var plan SchemaAttributeRegexValidationModel
		diags.Append(p.RegexValidation.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		regexValidation := management.NewSchemaAttributeRegexValidation(
			plan.Pattern.ValueString(),
			plan.Requirements.ValueString(),
		)

		if !plan.ValuesPatternShouldMatch.IsNull() && !plan.ValuesPatternShouldMatch.IsUnknown() {
			var values []string
			diags.Append(plan.ValuesPatternShouldMatch.ElementsAs(ctx, &values, false)...)
			if diags.HasError() {
				return nil, diags
			}

			regexValidation.SetValuesPatternShouldMatch(values)
		}

		if !plan.ValuesPatternShouldNotMatch.IsNull() && !plan.ValuesPatternShouldNotMatch.IsUnknown() {
			var values []string
			diags.Append(plan.ValuesPatternShouldNotMatch.ElementsAs(ctx, &values, false)...)
			if diags.HasError() {
				return nil, diags
			}

			regexValidation.SetValuesPatternShouldNotMatch(values)
		}

		data.SetRegexValidation(*regexValidation)
	}

	return &data, diags
}

func (p *SchemaAttributeResourceModelV1) toState(apiObject *management.SchemaAttribute, schemaApiObject *management.Schema) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	var d diag.Diagnostics

	p.Id = framework.StringOkToTF(apiObject.GetIdOk())
	p.EnvironmentId = framework.StringOkToTF(apiObject.Environment.GetIdOk())
	p.Description = framework.StringOkToTF(apiObject.GetDescriptionOk())
	p.DisplayName = framework.StringOkToTF(apiObject.GetDisplayNameOk())
	p.Enabled = framework.BoolOkToTF(apiObject.GetEnabledOk())

	p.EnumeratedValues, d = schemaAttributeEnumeratedValuesOkToTF(apiObject.GetEnumeratedValuesOk())
	diags.Append(d...)

	p.LdapAttribute = framework.StringOkToTF(apiObject.GetLdapAttributeOk())
	p.Multivalued = framework.BoolOkToTF(apiObject.GetMultiValuedOk())
	p.Name = framework.StringOkToTF(apiObject.GetNameOk())

	p.RegexValidation, d = schemaAttributeRegexValidationOkToTF(apiObject.GetRegexValidationOk())
	diags.Append(d...)

	p.Required = framework.BoolOkToTF(apiObject.GetRequiredOk())
	p.SchemaId = framework.StringOkToTF(apiObject.Schema.GetIdOk())
	p.SchemaName = framework.StringOkToTF(schemaApiObject.GetNameOk())
	p.SchemaType = framework.EnumOkToTF(apiObject.GetSchemaTypeOk())
	p.Type = framework.EnumOkToTF(apiObject.GetTypeOk())
	p.Unique = framework.BoolOkToTF(apiObject.GetUniqueOk())

	return diags
}

func schemaAttributeEnumeratedValuesOkToTF(apiObject []management.SchemaAttributeEnumeratedValuesInner, ok bool) (types.Set, diag.Diagnostics) {
	var diags diag.Diagnostics
	tfObjType := types.ObjectType{AttrTypes: schemaAttributeEnumeratedValuesTFObjectTypes}

	if !ok || len(apiObject) == 0 {
		return types.SetNull(tfObjType), diags
	}

	flattenedList := []attr.Value{}
	for _, v := range apiObject {

		objMap := map[string]attr.Value{
			"archived":    framework.BoolOkToTF(v.GetArchivedOk()),
			"description": framework.StringOkToTF(v.GetDescriptionOk()),
			"value":       framework.StringOkToTF(v.GetValueOk()),
		}

		flattenedObj, d := types.ObjectValue(schemaAttributeEnumeratedValuesTFObjectTypes, objMap)
		diags.Append(d...)

		flattenedList = append(flattenedList, flattenedObj)
	}

	returnVar, d := types.SetValue(tfObjType, flattenedList)
	diags.Append(d...)

	return returnVar, diags
}

func schemaAttributeRegexValidationOkToTF(apiObject *management.SchemaAttributeRegexValidation, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(schemaAttributeRegexValidationTFObjectTypes), diags
	}

	objMap := map[string]attr.Value{
		"pattern":                         framework.StringOkToTF(apiObject.GetPatternOk()),
		"requirements":                    framework.StringOkToTF(apiObject.GetRequirementsOk()),
		"values_pattern_should_match":     framework.StringSetOkToTF(apiObject.GetValuesPatternShouldMatchOk()),
		"values_pattern_should_not_match": framework.StringSetOkToTF(apiObject.GetValuesPatternShouldNotMatchOk()),
	}

	flattenedObj, d := types.ObjectValue(schemaAttributeRegexValidationTFObjectTypes, objMap)
	diags.Append(d...)

	return flattenedObj, diags
}
