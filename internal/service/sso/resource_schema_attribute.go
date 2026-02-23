// Copyright © 2026 Ping Identity Corporation

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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	boolplanmodifierinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/boolplanmodifier"
	boolvalidatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/boolvalidator"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/legacysdk"
	objectvalidatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/objectvalidator"
	setplanmodifierinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/setplanmodifier"
	setvalidatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/setvalidator"
	stringplanmodifierinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/stringplanmodifier"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type SchemaAttributeResource serviceClientType

type SchemaAttributeResourceModelV1 struct {
	Id               pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId    pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	Description      types.String                 `tfsdk:"description"`
	DisplayName      types.String                 `tfsdk:"display_name"`
	Enabled          types.Bool                   `tfsdk:"enabled"`
	EnumeratedValues types.Set                    `tfsdk:"enumerated_values"`
	LdapAttribute    types.String                 `tfsdk:"ldap_attribute"`
	Multivalued      types.Bool                   `tfsdk:"multivalued"`
	Name             types.String                 `tfsdk:"name"`
	RegexValidation  types.Object                 `tfsdk:"regex_validation"`
	SubAttributes    types.Set                    `tfsdk:"sub_attributes"`
	Required         types.Bool                   `tfsdk:"required"`
	SchemaId         pingonetypes.ResourceIDValue `tfsdk:"schema_id"`
	SchemaType       types.String                 `tfsdk:"schema_type"`
	Type             types.String                 `tfsdk:"type"`
	Unique           types.Bool                   `tfsdk:"unique"`
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

type SchemaAttributeSubAttributesResourceModel struct {
	Description types.String `tfsdk:"description"`
	DisplayName types.String `tfsdk:"display_name"`
	Enabled     types.Bool   `tfsdk:"enabled"`
	Multivalued types.Bool   `tfsdk:"multivalued"`
	Name        types.String `tfsdk:"name"`
	Required    types.Bool   `tfsdk:"required"`
	SchemaType  types.String `tfsdk:"schema_type"`
	Type        types.String `tfsdk:"type"`
	Unique      types.Bool   `tfsdk:"unique"`
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

	schemaAttributeSubAttributesTFObjectTypes = map[string]attr.Type{
		"description":  types.StringType,
		"display_name": types.StringType,
		"enabled":      types.BoolType,
		"multivalued":  types.BoolType,
		"name":         types.StringType,
		"required":     types.BoolType,
		"schema_type":  types.StringType,
		"type":         types.StringType,
		"unique":       types.BoolType,
	}
)

const (
	schemaName = "User"

	immutableDataLossProtectionDescription = "This field is immutable and cannot be changed once defined.  To protect against accidental data loss, this resource must be replaced manually (for example, by using Terraform's plan `-replace` command option https://developer.hashicorp.com/terraform/cli/commands/plan#replace-address).  Any data that is stored against this resource must be manually exported before the resource is removed and re-imported once the resource has been replaced."

	immutableDataLossProtectionMarkdownDescription = "This field is immutable and cannot be changed once defined.  To protect against accidental data loss, this resource must be replaced manually (for example, by using Terraform's [plan `-replace` command option](https://developer.hashicorp.com/terraform/cli/commands/plan#replace-address)).  Any data that is stored against this resource must be manually exported before the resource is removed and re-imported once the resource has been replaced."
)

// Framework interfaces
var (
	_ resource.Resource                 = &SchemaAttributeResource{}
	_ resource.ResourceWithConfigure    = &SchemaAttributeResource{}
	_ resource.ResourceWithImportState  = &SchemaAttributeResource{}
	_ resource.ResourceWithModifyPlan   = &SchemaAttributeResource{}
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

	schemaIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The ID of the schema the schema attribute is applied to.",
	)

	enabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Indicates whether or not the attribute is enabled. Can be updated for `STANDARD` attributes.",
	).DefaultValue("true")

	typeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The type of the attribute.",
	).AllowedValuesEnum(management.AllowedEnumSchemaAttributeTypeEnumValues).AppendMarkdownString(
		fmt.Sprintf("`%s` and `%s` attributes cannot be created, but standard attributes of those types may be updated. `%s` attributes are limited by size (total size must not exceed 16KB)", string(management.ENUMSCHEMAATTRIBUTETYPE_COMPLEX), string(management.ENUMSCHEMAATTRIBUTETYPE_BOOLEAN), string(management.ENUMSCHEMAATTRIBUTETYPE_JSON)),
	).UnmodifiableDataLossProtection().DefaultValue(string(management.ENUMSCHEMAATTRIBUTETYPE_STRING))

	uniqueDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Indicates whether or not the attribute must have a unique value within the PingOne environment. Can only be set where the attribute type is `STRING`. Can be updated for `STANDARD` attributes.",
	).DefaultValue("false")

	multivaluedDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Indicates whether the attribute has multiple values or a single one. Maximum number of values stored is 1,000.",
	).UnmodifiableDataLossProtection().DefaultValue("false")

	schemaTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The schema type of the attribute.",
	).AllowedValuesEnum(management.AllowedEnumSchemaAttributeSchemaTypeEnumValues).AppendMarkdownString(
		fmt.Sprintf("`%s` and `%s` attributes are supplied by default. `%s` attributes cannot be updated or deleted. `%s` attributes cannot be deleted, but their mutable properties can be updated. `%s` attributes can be deleted, and their mutable properties can be updated. New attributes are created with a schema type of `%s`.", management.ENUMSCHEMAATTRIBUTESCHEMATYPE_CORE, management.ENUMSCHEMAATTRIBUTESCHEMATYPE_STANDARD, management.ENUMSCHEMAATTRIBUTESCHEMATYPE_CORE, management.ENUMSCHEMAATTRIBUTESCHEMATYPE_STANDARD, management.ENUMSCHEMAATTRIBUTESCHEMATYPE_CUSTOM, management.ENUMSCHEMAATTRIBUTESCHEMATYPE_CUSTOM),
	)

	enumeratedValuesDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A set of one or more enumerated values for the attribute. If provided, it must not be an empty set.  Can only be set where the attribute type is `STRING` and cannot be set alongside `regex_validation`.  If the attribute has been created without enumerated values and this parameter is added later, this will trigger a replacement plan of the attribute resource.  If the attribute has been created with enumerated values that are subsequently removed, this will update without needing to replace the attribute resource.",
	)

	regexValidationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object representation of the optional regular expression representation of this attribute.  Can only be set where the attribute type is `STRING` and cannot be set alongside `enumerated_values`. Can be updated for `STANDARD` attributes.",
	)

	subAttributesDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The set of sub-attributes of this attribute. Only `COMPLEX` attribute types can have sub-attributes, and only one-level of nesting is allowed. The leaf attribute definition must have a type of `STRING` or `JSON`. A `COMPLEX` attribute definition must have at least one child attribute definition.",
	)

	resp.Schema = schema.Schema{

		Version: 1,

		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage PingOne schema attributes. Attributes with a `schema_type` of `STANDARD` are supported, but must be imported into Terraform state before they can be managed. Attributes with a `schema_type` of `CORE` are not supported and should be read using the `pingone_schema_attribute` data source.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to create the schema attribute in."),
			),

			"schema_id": schema.StringAttribute{
				Description:         schemaIdDescription.Description,
				MarkdownDescription: schemaIdDescription.MarkdownDescription,
				Computed:            true,

				CustomType: pingonetypes.ResourceIDType{},

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseNonNullStateForUnknown(),
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
				Computed:    true,
			},

			"description": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A description of the attribute. If provided, it must not be an empty string. Valid characters consists of any Unicode letter, mark (for example, accent or umlaut), numeric character, punctuation character, or space.").Description,
				Optional:    true,
				Computed:    true,
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
					stringplanmodifierinternal.UnmodifiableDataLossProtectionIf(
						unmodifiableDataLossProtectionIfStringConfigValueSet,
						immutableDataLossProtectionDescription,
						immutableDataLossProtectionMarkdownDescription,
					),
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
					boolplanmodifierinternal.UnmodifiableDataLossProtectionIf(
						unmodifiableDataLossProtectionIfBoolConfigValueSet,
						immutableDataLossProtectionDescription,
						immutableDataLossProtectionMarkdownDescription,
					),
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
					setplanmodifierinternal.UnmodifiableDataLossProtectionIf(
						setplanmodifierinternal.UnmodifiableDataLossProtectionIfPreviouslyNull(),
						"The attribute has been previously created without enumerated values validation.  To add enumerated values validation, the attribute must be replaced.",
						"The attribute has been previously created without enumerated values validation.  To add enumerated values validation, the attribute must be replaced.",
					),
					setplanmodifierinternal.UnmodifiableDataLossProtectionIf(
						unmodifiableDataLossProtectionIfElementRemoved,
						"Enumerated values cannot be deleted but can be archived.",
						"Enumerated values cannot be deleted but can be archived.",
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

			"sub_attributes": schema.SetNestedAttribute{
				Description:         subAttributesDescription.Description,
				MarkdownDescription: subAttributesDescription.MarkdownDescription,
				Computed:            true,

				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("The system name of the sub-attribute.").Description,
							Computed:    true,
						},
						"display_name": schema.StringAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("The display name of the sub-attribute.").Description,
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("A description of the sub-attribute.").Description,
							Computed:    true,
						},
						"enabled": schema.BoolAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("Indicates whether or not the sub-attribute is enabled.").Description,
							Computed:    true,
						},
						"required": schema.BoolAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("Indicates whether or not the sub-attribute is required.").Description,
							Computed:    true,
						},
						"schema_type": schema.StringAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("The schema type of the sub-attribute.").Description,
							Computed:    true,
						},
						"type": schema.StringAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("The type of the sub-attribute.").Description,
							Computed:    true,
						},
						"unique": schema.BoolAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("Indicates whether or not the sub-attribute must have a unique value within the PingOne environment.").Description,
							Computed:    true,
						},
						"multivalued": schema.BoolAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("Indicates whether the sub-attribute has multiple values or a single one.").Description,
							Computed:    true,
						},
					},
				},
			},

			"required": schema.BoolAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("Indicates whether or not the attribute is required.").Description,
				Computed:    true,

				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseNonNullStateForUnknown(),
				},
			},

			"ldap_attribute": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The unique identifier for the LDAP attribute.").Description,
				Computed:    true,

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseNonNullStateForUnknown(),
				},
			},

			"schema_type": schema.StringAttribute{
				Description:         schemaTypeDescription.Description,
				MarkdownDescription: schemaTypeDescription.MarkdownDescription,
				Computed:            true,

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseNonNullStateForUnknown(),
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

func (r *SchemaAttributeResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() || req.State.Raw.IsNull() {
		return
	}

	var state, config SchemaAttributeResourceModelV1
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.SchemaType.ValueString() == string(management.ENUMSCHEMAATTRIBUTESCHEMATYPE_STANDARD) {

		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("type"), state.Type)...)
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("multivalued"), state.Multivalued)...)
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("enumerated_values"), state.EnumeratedValues)...)
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("ldap_attribute"), state.LdapAttribute)...)

		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("sub_attributes"), state.SubAttributes)...)

		if config.DisplayName.IsNull() || config.DisplayName.IsUnknown() {
			resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("display_name"), state.DisplayName)...)
		}

		if config.Description.IsNull() || config.Description.IsUnknown() {
			resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("description"), state.Description)...)
		}

		return
	}

	if state.SchemaType.ValueString() == string(management.ENUMSCHEMAATTRIBUTESCHEMATYPE_CUSTOM) {
		if config.DisplayName.IsNull() {
			resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("display_name"), types.StringNull())...)
		}

		if config.Description.IsNull() {
			resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("description"), types.StringNull())...)
		}
	}
}

func (r *SchemaAttributeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state SchemaAttributeResourceModelV1

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

	// Get the schema ID
	schema, d := fetchSchemaFromName(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), schemaName)
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
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.SchemasApi.CreateAttribute(ctx, plan.EnvironmentId.ValueString(), schema.GetId()).SchemaAttribute(*schemaAttribute).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"CreateAttribute",
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

func (r *SchemaAttributeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *SchemaAttributeResourceModelV1

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

	// Run the API call
	var response *management.SchemaAttribute
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.SchemasApi.ReadOneAttribute(ctx, data.EnvironmentId.ValueString(), data.SchemaId.ValueString(), data.Id.ValueString()).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadOneAttribute",
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

func (r *SchemaAttributeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state SchemaAttributeResourceModelV1

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

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.SchemaType.ValueString() == string(management.ENUMSCHEMAATTRIBUTESCHEMATYPE_STANDARD) {
		immutableStandardChanges := standardImmutableChanges(plan, state)
		if len(immutableStandardChanges) > 0 {
			resp.Diagnostics.AddError(
				"Invalid update for STANDARD schema attribute",
				fmt.Sprintf("STANDARD schema attributes can only update 'enabled' (all types) and additionally 'unique' and 'regex_validation' for STRING attributes. Immutable attributes changed: %v", immutableStandardChanges),
			)
			return
		}
		normalizeStandardImmutableState(&plan, state)
	}

	// Get the schema ID
	schema, d := fetchSchemaFromName(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), schemaName)
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
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.SchemasApi.UpdateAttributePut(ctx, plan.EnvironmentId.ValueString(), schema.GetId(), plan.Id.ValueString()).SchemaAttribute(*schemaAttribute).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateAttributePut",
		legacysdk.DefaultCustomError,
		nil,
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

func standardImmutableChanges(plan, state SchemaAttributeResourceModelV1) []string {
	immutableStandardChanges := make([]string, 0)

	if !plan.Name.Equal(state.Name) {
		immutableStandardChanges = append(immutableStandardChanges, "name")
	}

	if !plan.Type.Equal(state.Type) {
		immutableStandardChanges = append(immutableStandardChanges, "type")
	}

	if !plan.Multivalued.Equal(state.Multivalued) {
		immutableStandardChanges = append(immutableStandardChanges, "multivalued")
	}

	if !plan.EnumeratedValues.Equal(state.EnumeratedValues) {
		immutableStandardChanges = append(immutableStandardChanges, "enumerated_values")
	}

	if !plan.DisplayName.IsNull() && !plan.DisplayName.IsUnknown() && !plan.DisplayName.Equal(state.DisplayName) {
		immutableStandardChanges = append(immutableStandardChanges, "display_name")
	}

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() && !plan.Description.Equal(state.Description) {
		immutableStandardChanges = append(immutableStandardChanges, "description")
	}

	if !plan.SchemaType.IsNull() && !plan.SchemaType.IsUnknown() && !plan.SchemaType.Equal(state.SchemaType) {
		immutableStandardChanges = append(immutableStandardChanges, "schema_type")
	}

	if state.Type.ValueString() != string(management.ENUMSCHEMAATTRIBUTETYPE_STRING) {
		if !plan.Unique.Equal(state.Unique) {
			immutableStandardChanges = append(immutableStandardChanges, "unique")
		}

		if !plan.RegexValidation.Equal(state.RegexValidation) {
			immutableStandardChanges = append(immutableStandardChanges, "regex_validation")
		}
	}

	return immutableStandardChanges
}

func normalizeStandardImmutableState(plan *SchemaAttributeResourceModelV1, state SchemaAttributeResourceModelV1) {
	plan.Name = state.Name
	plan.Type = state.Type
	plan.SubAttributes = state.SubAttributes
	plan.Multivalued = state.Multivalued
	plan.EnumeratedValues = state.EnumeratedValues
	plan.DisplayName = state.DisplayName
	plan.Description = state.Description
	plan.SchemaType = state.SchemaType
}

func (r *SchemaAttributeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *SchemaAttributeResourceModelV1

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

	if data.SchemaType.ValueString() == string(management.ENUMSCHEMAATTRIBUTESCHEMATYPE_STANDARD) {
		// Reset to defaults
		updateModel := *data
		updateModel.Enabled = types.BoolValue(true)
		updateModel.Unique = types.BoolValue(false)
		updateModel.RegexValidation = types.ObjectNull(schemaAttributeRegexValidationTFObjectTypes)

		updateAttribute, d := updateModel.expand(ctx, "UPDATE")
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}

		resp.Diagnostics.Append(legacysdk.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				fO, fR, fErr := r.Client.ManagementAPIClient.SchemasApi.UpdateAttributePut(ctx, data.EnvironmentId.ValueString(), data.SchemaId.ValueString(), data.Id.ValueString()).SchemaAttribute(*updateAttribute).Execute()
				return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
			},
			"UpdateAttributePut",
			legacysdk.CustomErrorResourceNotFoundWarning,
			nil,
			nil,
		)...)

		return
	}

	// Run the API call
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fR, fErr := r.Client.ManagementAPIClient.SchemasApi.DeleteAttribute(ctx, data.EnvironmentId.ValueString(), data.SchemaId.ValueString(), data.Id.ValueString()).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, fR, fErr)
		},
		"DeleteAttribute",
		legacysdk.CustomErrorResourceNotFoundWarning,
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

	if r.Client == nil || r.Client.ManagementAPIClient == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.",
		)
		return
	}

	attribute, _, err := r.Client.ManagementAPIClient.SchemasApi.ReadOneAttribute(ctx, attributes["environment_id"], attributes["schema_id"], attributes["schema_attribute_id"]).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Cannot import schema attribute",
			fmt.Sprintf("Failed to read schema attribute during import: %v", err),
		)
		return
	}

	if attribute.GetSchemaType() == management.ENUMSCHEMAATTRIBUTESCHEMATYPE_CORE {
		resp.Diagnostics.AddError(
			"Invalid import for CORE schema attribute",
			"CORE schema attributes are immutable and cannot be managed with the pingone_schema_attribute resource. Use the pingone_schema_attribute data source instead.",
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

	if action == "CREATE" {
		data.SetSchemaType(management.ENUMSCHEMAATTRIBUTESCHEMATYPE_CUSTOM)
	} else if !p.SchemaType.IsNull() && !p.SchemaType.IsUnknown() {
		data.SetSchemaType(management.EnumSchemaAttributeSchemaType(p.SchemaType.ValueString()))
	}

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

	data.SetRequired(p.Required.ValueBool())

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
			var valuesPatternPlan []types.String
			diags.Append(plan.ValuesPatternShouldMatch.ElementsAs(ctx, &valuesPatternPlan, false)...)
			if diags.HasError() {
				return nil, diags
			}

			values, d := framework.TFTypeStringSliceToStringSlice(valuesPatternPlan, path.Root("regex_validation").AtName("values_pattern_should_match"))
			diags.Append(d...)
			if diags.HasError() {
				return nil, diags
			}

			regexValidation.SetValuesPatternShouldMatch(values)
		}

		if !plan.ValuesPatternShouldNotMatch.IsNull() && !plan.ValuesPatternShouldNotMatch.IsUnknown() {
			var valuesPatternPlan []types.String
			diags.Append(plan.ValuesPatternShouldNotMatch.ElementsAs(ctx, &valuesPatternPlan, false)...)
			if diags.HasError() {
				return nil, diags
			}

			values, d := framework.TFTypeStringSliceToStringSlice(valuesPatternPlan, path.Root("regex_validation").AtName("values_pattern_should_not_match"))
			diags.Append(d...)
			if diags.HasError() {
				return nil, diags
			}

			regexValidation.SetValuesPatternShouldNotMatch(values)
		}

		data.SetRegexValidation(*regexValidation)
	}

	if !p.SubAttributes.IsNull() && !p.SubAttributes.IsUnknown() {
		var plan []SchemaAttributeSubAttributesResourceModel
		diags.Append(p.SubAttributes.ElementsAs(ctx, &plan, false)...)
		if diags.HasError() {
			return nil, diags
		}

		subAttributes := make([]management.SchemaAttribute, 0, len(plan))
		for _, v := range plan {
			subAttribute := management.NewSchemaAttribute(v.Enabled.ValueBool(), v.Name.ValueString(), management.EnumSchemaAttributeType(v.Type.ValueString()))

			if !v.DisplayName.IsNull() && !v.DisplayName.IsUnknown() {
				subAttribute.SetDisplayName(v.DisplayName.ValueString())
			}

			if !v.Description.IsNull() && !v.Description.IsUnknown() {
				subAttribute.SetDescription(v.Description.ValueString())
			}

			if !v.Required.IsNull() && !v.Required.IsUnknown() {
				subAttribute.SetRequired(v.Required.ValueBool())
			}

			if !v.SchemaType.IsNull() && !v.SchemaType.IsUnknown() {
				subAttribute.SetSchemaType(management.EnumSchemaAttributeSchemaType(v.SchemaType.ValueString()))
			}

			if !v.Unique.IsNull() && !v.Unique.IsUnknown() {
				subAttribute.SetUnique(v.Unique.ValueBool())
			}

			if !v.Multivalued.IsNull() && !v.Multivalued.IsUnknown() {
				subAttribute.SetMultiValued(v.Multivalued.ValueBool())
			}

			subAttributes = append(subAttributes, *subAttribute)
		}

		data.SetSubAttributes(subAttributes)
	}

	return &data, diags
}

func (p *SchemaAttributeResourceModelV1) toState(apiObject *management.SchemaAttribute) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	var d diag.Diagnostics

	p.Id = framework.PingOneResourceIDOkToTF(apiObject.GetIdOk())
	p.EnvironmentId = framework.PingOneResourceIDOkToTF(apiObject.Environment.GetIdOk())
	p.Description = framework.StringOkToTF(apiObject.GetDescriptionOk())
	p.DisplayName = framework.StringOkToTF(apiObject.GetDisplayNameOk())
	p.Enabled = framework.BoolOkToTF(apiObject.GetEnabledOk())

	p.EnumeratedValues, d = schemaAttributeEnumeratedValuesOkToTF(apiObject.GetEnumeratedValuesOk())
	diags.Append(d...)

	p.LdapAttribute = framework.StringOkToTF(apiObject.GetLdapAttributeOk())
	p.Multivalued = framework.BoolOkToTF(apiObject.GetMultiValuedOk())
	if p.Multivalued.IsNull() {
		p.Multivalued = types.BoolValue(false)
	}
	p.Name = framework.StringOkToTF(apiObject.GetNameOk())

	p.RegexValidation, d = schemaAttributeRegexValidationOkToTF(apiObject.GetRegexValidationOk())
	diags.Append(d...)

	p.SubAttributes, d = schemaAttributeSubAttributesOkToTF(apiObject.GetSubAttributesOk())
	diags.Append(d...)

	p.Required = framework.BoolOkToTF(apiObject.GetRequiredOk())
	p.SchemaId = framework.PingOneResourceIDOkToTF(apiObject.Schema.GetIdOk())
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

func schemaAttributeSubAttributesOkToTF(apiObject []management.SchemaAttribute, ok bool) (types.Set, diag.Diagnostics) {
	var diags diag.Diagnostics
	tfObjType := types.ObjectType{AttrTypes: schemaAttributeSubAttributesTFObjectTypes}

	if !ok || len(apiObject) == 0 {
		return types.SetNull(tfObjType), diags
	}

	flattenedList := []attr.Value{}
	for _, v := range apiObject {
		objMap := map[string]attr.Value{
			"description":  framework.StringOkToTF(v.GetDescriptionOk()),
			"display_name": framework.StringOkToTF(v.GetDisplayNameOk()),
			"enabled":      framework.BoolOkToTF(v.GetEnabledOk()),
			"multivalued":  framework.BoolOkToTF(v.GetMultiValuedOk()),
			"name":         framework.StringOkToTF(v.GetNameOk()),
			"required":     framework.BoolOkToTF(v.GetRequiredOk()),
			"schema_type":  framework.EnumOkToTF(v.GetSchemaTypeOk()),
			"type":         framework.EnumOkToTF(v.GetTypeOk()),
			"unique":       framework.BoolOkToTF(v.GetUniqueOk()),
		}

		flattenedObj, d := types.ObjectValue(schemaAttributeSubAttributesTFObjectTypes, objMap)
		diags.Append(d...)

		flattenedList = append(flattenedList, flattenedObj)
	}

	returnVar, d := types.SetValue(tfObjType, flattenedList)
	diags.Append(d...)

	return returnVar, diags
}

func unmodifiableDataLossProtectionIfElementRemoved(ctx context.Context, req planmodifier.SetRequest, resp *setplanmodifierinternal.UnmodifiableDataLossProtectionIfFuncResponse) {
	// If the configuration is unknown, this cannot be sure what to do yet.
	if req.ConfigValue.IsUnknown() {
		resp.Error = false
		return
	}

	// If the state is not null and the config value is not null, error
	if !req.StateValue.IsNull() && !req.ConfigValue.IsNull() {
		var config, state []SchemaAttributeEnumeratedValuesResourceModel
		resp.Diagnostics.Append(req.StateValue.ElementsAs(ctx, &state, false)...)
		resp.Diagnostics.Append(req.ConfigValue.ElementsAs(ctx, &config, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		if len(config) < len(state) {
			resp.Error = true
			return
		}

		for _, stateValue := range state {
			found := false
			for _, configValue := range config {
				if stateValue.Value.ValueString() == configValue.Value.ValueString() {
					found = true
					break
				}
			}

			if !found {
				resp.Error = true
				return
			}
		}
	}

	resp.Error = false
}

func unmodifiableDataLossProtectionIfStringConfigValueSet(ctx context.Context, req planmodifier.StringRequest, resp *stringplanmodifierinternal.UnmodifiableDataLossProtectionIfFuncResponse) {
	if req.ConfigValue.IsUnknown() {
		resp.Error = false
		return
	}

	if req.ConfigValue.IsNull() {
		resp.Error = false
		return
	}

	resp.Error = true
}

func unmodifiableDataLossProtectionIfBoolConfigValueSet(ctx context.Context, req planmodifier.BoolRequest, resp *boolplanmodifierinternal.UnmodifiableDataLossProtectionIfFuncResponse) {
	if req.ConfigValue.IsUnknown() {
		resp.Error = false
		return
	}

	if req.ConfigValue.IsNull() {
		resp.Error = false
		return
	}

	resp.Error = true
}
