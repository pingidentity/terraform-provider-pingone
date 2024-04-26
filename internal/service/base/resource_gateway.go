package base

import (
	"context"
	"fmt"
	"net/http"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type GatewayResource serviceClientType

type gatewayResourceModel struct {
	Id            pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	Name          types.String                 `tfsdk:"name"`
	Description   types.String                 `tfsdk:"description"`
	Type          types.String                 `tfsdk:"type"`
	Enabled       types.Bool                   `tfsdk:"enabled"`

	// LDAP
	BindDN                  types.String `tfsdk:"bind_dn"`
	BindPassword            types.String `tfsdk:"bind_password"`
	ConnectionSecurity      types.String `tfsdk:"connection_security"`
	FollowReferrals         types.Bool   `tfsdk:"follow_referrals"`
	Kerberos                types.Object `tfsdk:"kerberos"`
	Servers                 types.Set    `tfsdk:"servers"`
	ValidateTLSCertificates types.Bool   `tfsdk:"validate_tls_certificates"`
	Vendor                  types.String `tfsdk:"vendor"`
	UserTypes               types.Set    `tfsdk:"user_types"`

	// Radius
	RadiusClients             types.Set                    `tfsdk:"radius_clients"`
	RadiusDavinciPolicyId     pingonetypes.ResourceIDValue `tfsdk:"radius_davinci_policy_id"`
	RadiusDefaultSharedSecret types.String                 `tfsdk:"radius_default_shared_secret"`
	RadiusNetworkPolicyServer types.Object                 `tfsdk:"radius_network_policy_server"`
}

type gatewayKerberosResourceModel struct {
	ServiceAccountPassword        types.String `tfsdk:"service_account_password"`
	ServiceAccountUPN             types.String `tfsdk:"service_account_upn"`
	RetainPreviousCredentialsMins types.Int64  `tfsdk:"retain_previous_credentials_mins"`
}

type gatewayUserTypeResourceModel struct {
	AllowPasswordChanges                 types.Bool                   `tfsdk:"allow_password_changes"`
	Id                                   pingonetypes.ResourceIDValue `tfsdk:"id"`
	Name                                 types.String                 `tfsdk:"name"`
	NewUserLookup                        types.Object                 `tfsdk:"new_user_lookup"`
	PasswordAuthority                    types.String                 `tfsdk:"password_authority"`
	PushPasswordChangesToLDAP            types.Bool                   `tfsdk:"push_password_changes_to_ldap"`
	SearchBaseDN                         types.String                 `tfsdk:"search_base_dn"`
	UpdateUserOnSuccessfulAuthentication types.Bool                   `tfsdk:"update_user_on_successful_authentication"`
	UserLinkAttributes                   types.List                   `tfsdk:"user_link_attributes"`
}

type gatewayUserTypeNewUserLookupResourceModel struct {
	AttributeMappings types.Set                    `tfsdk:"attribute_mappings"`
	LDAPFilterPattern types.String                 `tfsdk:"ldap_filter_pattern"`
	PopulationId      pingonetypes.ResourceIDValue `tfsdk:"population_id"`
}

type gatewayUserTypeMigrationAttributeMappingResourceModel struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

type gatewayRadiusClientsResourceModel struct {
	IP           types.String `tfsdk:"ip"`
	SharedSecret types.String `tfsdk:"shared_secret"`
}

type gatewayRadiusNetworkPolicyServerResourceModel struct {
	IP   types.String `tfsdk:"ip"`
	Port types.Int64  `tfsdk:"port"`
}

var (
	gatewayKerberosTFObjectTypes = map[string]attr.Type{
		"retain_previous_credentials_mins": types.Int64Type,
		"service_account_password":         types.StringType,
		"service_account_upn":              types.StringType,
	}

	gatewayUserTypesTFObjectTypes = map[string]attr.Type{
		"allow_password_changes":        types.BoolType,
		"id":                            pingonetypes.ResourceIDType{},
		"name":                          types.StringType,
		"new_user_lookup":               types.ObjectType{AttrTypes: gatewayUserTypesNewUserLookupTFObjectTypes},
		"password_authority":            types.StringType,
		"push_password_changes_to_ldap": types.BoolType,
		"search_base_dn":                types.StringType,
		"update_user_on_successful_authentication": types.BoolType,
		"user_link_attributes":                     types.ListType{ElemType: types.StringType},
	}

	gatewayUserTypesNewUserLookupTFObjectTypes = map[string]attr.Type{
		"attribute_mappings": types.SetType{
			ElemType: types.ObjectType{
				AttrTypes: gatewayUserTypesNewUserLookupAttributeMappingTFObjectTypes,
			},
		},
		"ldap_filter_pattern": types.StringType,
		"population_id":       pingonetypes.ResourceIDType{},
	}

	gatewayUserTypesNewUserLookupAttributeMappingTFObjectTypes = map[string]attr.Type{
		"name":  types.StringType,
		"value": types.StringType,
	}

	gatewayRadiusClientsTFObjectTypes = map[string]attr.Type{
		"ip":            types.StringType,
		"shared_secret": types.StringType,
	}

	gatewayRadiusNetworkPolicyServerTFObjectTypes = map[string]attr.Type{
		"ip":   types.StringType,
		"port": types.Int64Type,
	}
)

// Framework interfaces
var (
	_ resource.Resource                = &GatewayResource{}
	_ resource.ResourceWithConfigure   = &GatewayResource{}
	_ resource.ResourceWithImportState = &GatewayResource{}
	_ resource.ResourceWithModifyPlan  = &GatewayResource{}
)

// New Object
func NewGatewayResource() resource.Resource {
	return &GatewayResource{}
}

// Metadata
func (r *GatewayResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_gateway"
}

func (r *GatewayResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	typeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the type of gateway.",
	).AllowedValuesEnum(management.AllowedEnumGatewayTypeEnumValues).RequiresReplace()

	connectionSecurityDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"For LDAP gateways only: A string that specifies the connection security type.",
	).AllowedValuesEnum(management.AllowedEnumGatewayTypeLDAPSecurityEnumValues).DefaultValue(management.ENUMGATEWAYTYPELDAPSECURITY_NONE)

	followReferralsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that, when set to true, PingOne sends LDAP queries per referrals it receives from the LDAP servers.",
	).DefaultValue(false)

	kerberosServiceAccountUpnDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the Kerberos service account user principal name (for example, `username@bxretail.org`).",
	)

	serversDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"For LDAP gateways only: A set of LDAP server host name and port number combinations (for example, [`ds1.bxretail.org:636`, `ds2.bxretail.org:636`]).",
	)

	validateTlsCertificatesDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"For LDAP gateways only: A boolean that specifies whether or not to trust all SSL certificates, including self-signed (defaults to `true`). If this value is `false`, TLS certificates are not validated. When the value is set to `true`, only certificates that are signed by the default JVM CAs, or the CA certs that the customer has uploaded to the certificate service are trusted.",
	).DefaultValue(true)

	vendorDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"For LDAP gateways only: A string that specifies the LDAP vendor.",
	).AllowedValuesEnum(management.AllowedEnumGatewayVendorEnumValues).RequiresReplace()

	userTypesDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"For LDAP gateways only: A set of objects that define how users should be provisioned in PingOne. The `user_types` set of objects specifies which user properties in PingOne correspond to the user properties in an external LDAP directory. You can use an LDAP browser to view the user properties in the external LDAP directory.",
	)

	userTypesAllowPasswordChangesDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that, if set to `false`, the user cannot change the password in the remote LDAP directory. In this case, operations for forgotten passwords or resetting of passwords are not available to a user referencing this gateway.",
	).DefaultValue(false)

	userTypesUpdateUserOnSuccessfulAuthenticationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that, if set to `true`, when users sign on through an LDAP Gateway client, user attributes are updated based on responses from the LDAP server.",
	).DefaultValue(false)

	userTypesIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Identifies the user type. This correlates to the `password.external.gateway.userType.id` User property.",
	)

	userTypesPasswordAuthorityDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the password authority for the user type.",
	).AllowedValuesEnum(management.AllowedEnumGatewayPasswordAuthorityEnumValues).AppendMarkdownString(fmt.Sprintf("If set to `%s`, PingOne authenticates with the external directory initially, then PingOne authenticates all subsequent sign-ons.", string(management.ENUMGATEWAYPASSWORDAUTHORITY_PING_ONE)))

	userTypesUserMigrationLookupFilterPatternDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The LDAP user search filter to use to match users against the entered user identifier at login. For example, `(((uid=${identifier})(mail=${identifier}))`. Alternatively, this can be a search against the user directory.",
	)

	userTypesUserMigrationLookupAttributeMappingDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A set of objects supplying a mapping of PingOne attributes to external LDAP attributes. One of the entries must be a mapping for `username`. This is required for the PingOne user schema.",
	)

	userTypesUserMigrationLookupAttributeMappingNameDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the name of a user attribute in PingOne. See [Users properties](https://apidocs.pingidentity.com/pingone/platform/v1/api/#users) for the complete list of available PingOne user attributes.",
	)

	userTypesUserMigrationLookupAttributeMappingValueDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the reference to the corresponding external LDAP attribute.  Values are in the format `${ldapAttributes.mail}`, while Terraform HCL requires an additional `$` prefix character. For example, `$${ldapAttributes.mail}`.",
	)

	userTypesPushPasswordChangesToLdapDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether password updates in PingOne should be pushed to the user's record in LDAP.  If false, the user cannot change the password and have it updated in the remote LDAP directory. In this case, operations for forgotten passwords or resetting of passwords are not available to a user referencing this gateway.",
	).DefaultValue(false)

	radiusClientSharedSecretDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the shared secret for the RADIUS client. If this value is not provided, the shared secret specified with `radius_default_shared_secret` is used. If you are not providing a shared secret for the client, this parameter is optional.",
	)

	ldapRequiredSchemaPaths := []path.Expression{}
	radiusRequiredSchemaPaths := []path.Expression{}

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage gateway configuration in a PingOne environment.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to manage the gateway in."),
			),

			"name": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the name of the gateway resource.").Description,
				Required:    true,
			},

			"description": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies a description to apply to the gateway resource.").Description,
				Optional:    true,
			},

			"type": schema.StringAttribute{
				Description:         typeDescription.Description,
				MarkdownDescription: typeDescription.MarkdownDescription,
				Required:            true,

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},

				Validators: []validator.String{
					stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumGatewayTypeEnumValues)...),
				},
			},

			"enabled": schema.BoolAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether the gateway is enabled in the environment.").Description,
				Required:    true,
			},

			// LDAP
			"bind_dn": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("For LDAP gateways only: A string that specifies the distinguished name information to bind to the LDAP database (for example, `uid=pingone,dc=bxretail,dc=org`).").Description,
				Optional:    true,

				Validators: []validator.String{
					stringvalidator.AlsoRequires(ldapRequiredSchemaPaths...),
				},
			},

			"bind_password": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("For LDAP gateways only: A string that specifies the bind password for the LDAP database.").Description,
				Optional:    true,
				Sensitive:   true,

				Validators: []validator.String{
					stringvalidator.AlsoRequires(ldapRequiredSchemaPaths...),
				},
			},

			"connection_security": schema.StringAttribute{
				Description:         connectionSecurityDescription.Description,
				MarkdownDescription: connectionSecurityDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,

				//Default: stringdefault.StaticString(string(management.ENUMGATEWAYTYPELDAPSECURITY_NONE)),

				Validators: []validator.String{
					stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumGatewayTypeLDAPSecurityEnumValues)...),
					stringvalidator.AlsoRequires(ldapRequiredSchemaPaths...),
				},
			},

			"follow_referrals": schema.BoolAttribute{
				Description:         followReferralsDescription.Description,
				MarkdownDescription: followReferralsDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,

				//Default: booldefault.StaticBool(false),
			},

			"kerberos": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("For LDAP gateways only: A single object that specifies Kerberos connection details.").Description,
				Optional:    true,

				Validators: []validator.Object{
					objectvalidator.AlsoRequires(ldapRequiredSchemaPaths...),
				},

				Attributes: map[string]schema.Attribute{
					"service_account_password": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the password for the Kerberos service account.").Description,
						Optional:    true,
						Sensitive:   true,
					},

					"service_account_upn": schema.StringAttribute{
						Description:         kerberosServiceAccountUpnDescription.Description,
						MarkdownDescription: kerberosServiceAccountUpnDescription.MarkdownDescription,
						Required:            true,
					},

					"retain_previous_credentials_mins": schema.Int64Attribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that specifies the number of minutes for which the previous credentials are persisted.").Description,
						Optional:    true,
					},
				},
			},

			"servers": schema.SetAttribute{
				Description:         serversDescription.Description,
				MarkdownDescription: serversDescription.MarkdownDescription,
				Optional:            true,

				ElementType: types.StringType,

				Validators: []validator.Set{
					setvalidator.AlsoRequires(ldapRequiredSchemaPaths...),
				},
			},

			"validate_tls_certificates": schema.BoolAttribute{
				Description:         validateTlsCertificatesDescription.Description,
				MarkdownDescription: validateTlsCertificatesDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,

				//Default: booldefault.StaticBool(true),

				Validators: []validator.Bool{
					boolvalidator.AlsoRequires(ldapRequiredSchemaPaths...),
				},
			},

			"vendor": schema.StringAttribute{
				Description:         vendorDescription.Description,
				MarkdownDescription: vendorDescription.MarkdownDescription,
				Optional:            true,

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},

				Validators: []validator.String{
					stringvalidator.AlsoRequires(ldapRequiredSchemaPaths...),
					stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumGatewayVendorEnumValues)...),
				},
			},

			"user_types": schema.SetNestedAttribute{
				Description:         userTypesDescription.Description,
				MarkdownDescription: userTypesDescription.MarkdownDescription,
				Optional:            true,

				Validators: []validator.Set{
					setvalidator.AlsoRequires(ldapRequiredSchemaPaths...),
				},

				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"allow_password_changes": schema.BoolAttribute{
							Description:         userTypesAllowPasswordChangesDescription.Description,
							MarkdownDescription: userTypesAllowPasswordChangesDescription.MarkdownDescription,
							Optional:            true,
							Computed:            true,

							Default: booldefault.StaticBool(false),
						},

						"update_user_on_successful_authentication": schema.BoolAttribute{
							Description:         userTypesUpdateUserOnSuccessfulAuthenticationDescription.Description,
							MarkdownDescription: userTypesUpdateUserOnSuccessfulAuthenticationDescription.MarkdownDescription,
							Optional:            true,
							Computed:            true,

							Default: booldefault.StaticBool(false),
						},

						"id": schema.StringAttribute{
							Description:         userTypesIdDescription.Description,
							MarkdownDescription: userTypesIdDescription.MarkdownDescription,
							Computed:            true,

							CustomType: pingonetypes.ResourceIDType{},
						},

						"name": schema.StringAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the name of the user type.").Description,
							Required:    true,
						},

						"password_authority": schema.StringAttribute{
							Description:         userTypesPasswordAuthorityDescription.Description,
							MarkdownDescription: userTypesPasswordAuthorityDescription.MarkdownDescription,
							Required:            true,

							Validators: []validator.String{
								stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumGatewayPasswordAuthorityEnumValues)...),
							},
						},

						"search_base_dn": schema.StringAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the LDAP base domain name (DN) for this user type.").Description,
							Required:    true,
						},

						"user_link_attributes": schema.ListAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("A list of strings that represent LDAP attribute names that uniquely identify the user, and link to users in PingOne.").Description,
							Required:    true,

							ElementType: types.StringType,
						},

						"new_user_lookup": schema.SingleNestedAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that describes the configurations for initially authenticating new users who will be migrated to PingOne. Note: If there are multiple users having the same user name, only the first user processed is provisioned.").Description,
							Optional:    true,

							Attributes: map[string]schema.Attribute{
								"ldap_filter_pattern": schema.StringAttribute{
									Description:         userTypesUserMigrationLookupFilterPatternDescription.Description,
									MarkdownDescription: userTypesUserMigrationLookupFilterPatternDescription.MarkdownDescription,
									Required:            true,
								},

								"population_id": schema.StringAttribute{
									Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the ID of the population to use to create user entries during lookup.  Must be a valid PingOne resource ID.").Description,
									Required:    true,

									CustomType: pingonetypes.ResourceIDType{},
								},

								"attribute_mappings": schema.SetNestedAttribute{
									Description:         userTypesUserMigrationLookupAttributeMappingDescription.Description,
									MarkdownDescription: userTypesUserMigrationLookupAttributeMappingDescription.MarkdownDescription,
									Required:            true,

									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"name": schema.StringAttribute{
												Description:         userTypesUserMigrationLookupAttributeMappingNameDescription.Description,
												MarkdownDescription: userTypesUserMigrationLookupAttributeMappingNameDescription.MarkdownDescription,
												Required:            true,
											},

											"value": schema.StringAttribute{
												Description:         userTypesUserMigrationLookupAttributeMappingValueDescription.Description,
												MarkdownDescription: userTypesUserMigrationLookupAttributeMappingValueDescription.MarkdownDescription,
												Required:            true,
											},
										},
									},
								},
							},
						},

						"push_password_changes_to_ldap": schema.BoolAttribute{
							Description:         userTypesPushPasswordChangesToLdapDescription.Description,
							MarkdownDescription: userTypesPushPasswordChangesToLdapDescription.MarkdownDescription,
							Optional:            true,
							Computed:            true,

							Default: booldefault.StaticBool(false),
						},
					},
				},
			},

			"radius_davinci_policy_id": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("For RADIUS gateways only: A string that specifies the ID of the DaVinci flow policy to use.  Must be a valid PingOne resource ID.").Description,
				Optional:    true,

				CustomType: pingonetypes.ResourceIDType{},

				Validators: []validator.String{
					stringvalidator.AlsoRequires(radiusRequiredSchemaPaths...),
				},
			},

			"radius_default_shared_secret": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("For RADIUS gateways only: A strign that specifies the value to use for the shared secret if the shared secret is not provided for one or more of the RADIUS clients specified.").Description,
				Optional:    true,
				Sensitive:   true,

				Validators: []validator.String{
					stringvalidator.AlsoRequires(radiusRequiredSchemaPaths...),
				},
			},

			"radius_clients": schema.SetNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("For RADIUS gateways only: A set of objects describing RADIUS client connections.").Description,
				Optional:    true,

				Validators: []validator.Set{
					setvalidator.AlsoRequires(radiusRequiredSchemaPaths...),
					setvalidator.SizeAtLeast(1),
				},

				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"ip": schema.StringAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the IP address of the RADIUS client.").Description,
							Required:    true,

							Validators: []validator.String{
								stringvalidator.RegexMatches(verify.IPv4Regexp, "The IP address must be a valid IPv4 address."),
							},
						},

						"shared_secret": schema.StringAttribute{
							Description:         radiusClientSharedSecretDescription.Description,
							MarkdownDescription: radiusClientSharedSecretDescription.MarkdownDescription,
							Optional:            true,
							Sensitive:           true,
						},
					},
				},
			},

			"radius_network_policy_server": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("For RADIUS gateways only: A single object that allows configuration of the RADIUS gateway to authenticate using the MS-CHAP v2 protocol.").Description,
				Optional:    true,

				Validators: []validator.Object{
					objectvalidator.AlsoRequires(radiusRequiredSchemaPaths...),
				},

				Attributes: map[string]schema.Attribute{
					"ip": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the IP address of the Network Policy Server (NPS).").Description,
						Required:    true,

						Validators: []validator.String{
							stringvalidator.RegexMatches(verify.IPv4Regexp, "The IP address must be a valid IPv4 address."),
						},
					},

					"port": schema.Int64Attribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that specifies the port number of the NPS.").Description,
						Required:    true,
					},
				},
			},
		},
	}
}

// ModifyPlan
func (r *GatewayResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {

	// Destruction plan
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan gatewayResourceModel
	// Read Terraform plan and state data into the model
	resp.Diagnostics.Append(resp.Plan.Get(ctx, &plan)...)

	if plan.Type.Equal(types.StringValue(string(management.ENUMGATEWAYTYPE_LDAP))) {
		if plan.ConnectionSecurity.IsNull() {
			resp.Plan.SetAttribute(ctx, path.Root("connection_security"), types.StringValue(string(management.ENUMGATEWAYTYPELDAPSECURITY_NONE)))
		}

		if plan.FollowReferrals.IsNull() {
			resp.Plan.SetAttribute(ctx, path.Root("follow_referrals"), types.BoolValue(false))
		}

		if plan.ValidateTLSCertificates.IsNull() {
			resp.Plan.SetAttribute(ctx, path.Root("validate_tls_certificates"), types.BoolValue(true))
		}
	}

}

func (r *GatewayResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *GatewayResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state gatewayResourceModel

	if r.Client.MFAAPIClient == nil {
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
	createGatewayRequest, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *management.CreateGateway201Response
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.GatewaysApi.CreateGateway(ctx, plan.EnvironmentId.ValueString()).CreateGatewayRequest(*createGatewayRequest).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"CreateGateway",
		gatewayWriteErrors,
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

func (r *GatewayResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *gatewayResourceModel

	if r.Client.MFAAPIClient == nil {
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
	var response *management.CreateGateway201Response
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.GatewaysApi.ReadOneGateway(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadOneGateway",
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

func (r *GatewayResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state gatewayResourceModel

	if r.Client.MFAAPIClient == nil {
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
	createGatewayRequest, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *management.CreateGateway201Response
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.GatewaysApi.UpdateGateway(ctx, plan.EnvironmentId.ValueString(), plan.Id.ValueString()).CreateGatewayRequest(*createGatewayRequest).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateGateway",
		gatewayWriteErrors,
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

func (r *GatewayResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *gatewayResourceModel

	if r.Client.MFAAPIClient == nil {
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
			fR, fErr := r.Client.ManagementAPIClient.GatewaysApi.DeleteGateway(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, fR, fErr)
		},
		"DeleteGateway",
		framework.CustomErrorResourceNotFoundWarning,
		nil,
		nil,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *GatewayResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:     "gateway_id",
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

func (p *gatewayResourceModel) expand(ctx context.Context) (*management.CreateGatewayRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Main object
	data := &management.CreateGatewayRequest{}

	gatewayType := management.EnumGatewayType(p.Type.ValueString())

	if slices.Contains([]management.EnumGatewayType{
		"PING_FEDERATE",
		"PING_INTELLIGENCE",
		"API_GATEWAY_INTEGRATION",
	}, gatewayType) {

		gateway := *management.NewGateway(p.Name.ValueString(), gatewayType, p.Enabled.ValueBool())

		if !p.Description.IsNull() && !p.Description.IsUnknown() {
			gateway.SetDescription(p.Description.ValueString())
		}

		data.Gateway = &gateway

	} else if gatewayType == management.ENUMGATEWAYTYPE_LDAP {

		servers := make([]string, 0)

		if !p.Servers.IsNull() && !p.Servers.IsUnknown() {
			var serversPlan []string
			diags.Append(p.Servers.ElementsAs(ctx, &serversPlan, false)...)
			if diags.HasError() {
				return nil, diags
			}

			for _, server := range serversPlan {
				servers = append(servers, server)
			}
		}

		gateway := *management.NewGatewayTypeLDAP(
			p.Name.ValueString(),
			gatewayType,
			p.Enabled.ValueBool(),
			p.BindDN.ValueString(),
			p.BindPassword.ValueString(),
			servers,
			management.EnumGatewayVendor(p.Vendor.ValueString()),
		)

		if !p.ConnectionSecurity.IsNull() && !p.ConnectionSecurity.IsUnknown() {
			gateway.SetConnectionSecurity(management.EnumGatewayTypeLDAPSecurity(p.ConnectionSecurity.ValueString()))
		}

		if !p.Kerberos.IsNull() && !p.Kerberos.IsUnknown() {

			var kerberosPlan gatewayKerberosResourceModel
			diags.Append(p.Kerberos.As(ctx, &kerberosPlan, basetypes.ObjectAsOptions{
				UnhandledNullAsEmpty:    false,
				UnhandledUnknownAsEmpty: false,
			})...)
			if diags.HasError() {
				return nil, diags
			}

			kerberos := management.NewGatewayTypeLDAPAllOfKerberos(kerberosPlan.ServiceAccountUPN.ValueString())

			if !kerberosPlan.ServiceAccountPassword.IsNull() && !kerberosPlan.ServiceAccountPassword.IsUnknown() {
				kerberos.SetServiceAccountPassword(kerberosPlan.ServiceAccountPassword.ValueString())
			}

			if !kerberosPlan.RetainPreviousCredentialsMins.IsNull() && !kerberosPlan.RetainPreviousCredentialsMins.IsUnknown() {
				kerberos.SetMinutesToRetainPreviousCredentials(int32(kerberosPlan.RetainPreviousCredentialsMins.ValueInt64()))
			}

			gateway.SetKerberos(*kerberos)
		}

		if !p.ValidateTLSCertificates.IsNull() && !p.ValidateTLSCertificates.IsUnknown() {
			gateway.SetValidateTlsCertificates(p.ValidateTLSCertificates.ValueBool())
		}

		if !p.UserTypes.IsNull() && !p.UserTypes.IsUnknown() {
			var userTypesPlan []gatewayUserTypeResourceModel
			diags.Append(p.UserTypes.ElementsAs(ctx, &userTypesPlan, false)...)
			if diags.HasError() {
				return nil, diags
			}

			userTypes := make([]management.GatewayTypeLDAPAllOfUserTypes, 0)

			for _, userTypePlan := range userTypesPlan {
				userType, d := userTypePlan.expandLDAPUserType(ctx)
				diags.Append(d...)
				if diags.HasError() {
					return nil, diags
				}

				userTypes = append(userTypes, *userType)
			}

			gateway.SetUserTypes(userTypes)
		}

		data.GatewayTypeLDAP = &gateway

	} else if gatewayType == management.ENUMGATEWAYTYPE_RADIUS {

		radiusClients := make([]management.GatewayTypeRADIUSAllOfRadiusClients, 0)

		if !p.RadiusClients.IsNull() && !p.RadiusClients.IsUnknown() {
			var radiusClientsPlan []gatewayRadiusClientsResourceModel
			diags.Append(p.RadiusClients.ElementsAs(ctx, &radiusClientsPlan, false)...)
			if diags.HasError() {
				return nil, diags
			}

			for _, client := range radiusClientsPlan {

				radiusClientObj := *management.NewGatewayTypeRADIUSAllOfRadiusClients(client.IP.ValueString())

				if !client.SharedSecret.IsNull() && !client.SharedSecret.IsUnknown() {
					radiusClientObj.SetSharedSecret(client.SharedSecret.ValueString())
				}

				radiusClients = append(radiusClients, radiusClientObj)
			}
		}

		gateway := *management.NewGatewayTypeRADIUS(
			p.Name.ValueString(),
			gatewayType,
			p.Enabled.ValueBool(),
			*management.NewGatewayTypeRADIUSAllOfDavinci(*management.NewGatewayTypeRADIUSAllOfDavinciPolicy(p.RadiusDavinciPolicyId.ValueString())),
			radiusClients,
		)

		if !p.RadiusDefaultSharedSecret.IsNull() && !p.RadiusDefaultSharedSecret.IsUnknown() {
			gateway.SetDefaultSharedSecret(p.RadiusDefaultSharedSecret.ValueString())
		}

		data.GatewayTypeRADIUS = &gateway

	} else {

		diags.AddAttributeError(
			path.Root("type"),
			"Unsupported gateway type",
			"The gateway type value of %s is not supported in the provider.  Ensure that the configuration of the resource is set correctly.",
		)

		return nil, diags
	}

	return data, diags
}

func (p *gatewayUserTypeResourceModel) expandLDAPUserType(ctx context.Context) (*management.GatewayTypeLDAPAllOfUserTypes, diag.Diagnostics) {
	var diags diag.Diagnostics

	var userLinkAttributesPlan []string
	diags.Append(p.UserLinkAttributes.ElementsAs(ctx, &userLinkAttributesPlan, false)...)
	if diags.HasError() {
		return nil, diags
	}

	data := management.NewGatewayTypeLDAPAllOfUserTypes(
		p.Name.ValueString(),
		userLinkAttributesPlan,
		management.EnumGatewayPasswordAuthority(p.PasswordAuthority.ValueString()),
		p.SearchBaseDN.ValueString(),
	)

	if !p.PushPasswordChangesToLDAP.IsNull() && !p.PushPasswordChangesToLDAP.IsUnknown() {
		data.SetAllowPasswordChanges(p.PushPasswordChangesToLDAP.ValueBool())
	}

	if !p.NewUserLookup.IsNull() && !p.NewUserLookup.IsUnknown() {

		var newUserLookupPlan gatewayUserTypeNewUserLookupResourceModel
		diags.Append(p.NewUserLookup.As(ctx, &newUserLookupPlan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		newUserLookup, d := newUserLookupPlan.expandLDAPUserTypeNewUserLookup(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}
		data.SetNewUserLookup(*newUserLookup)
	}

	return data, diags
}

func (p *gatewayUserTypeNewUserLookupResourceModel) expandLDAPUserTypeNewUserLookup(ctx context.Context) (*management.GatewayTypeLDAPAllOfNewUserLookup, diag.Diagnostics) {
	var diags diag.Diagnostics

	var attributeMappingsPlan []gatewayUserTypeMigrationAttributeMappingResourceModel
	diags.Append(p.AttributeMappings.ElementsAs(ctx, &attributeMappingsPlan, false)...)
	if diags.HasError() {
		return nil, diags
	}

	attributeMappings := make([]management.GatewayTypeLDAPAllOfNewUserLookupAttributeMappings, 0)

	for _, attributeMapping := range attributeMappingsPlan {
		attributeMappings = append(attributeMappings, *management.NewGatewayTypeLDAPAllOfNewUserLookupAttributeMappings(
			attributeMapping.Name.ValueString(),
			attributeMapping.Value.ValueString(),
		))
	}

	data := management.NewGatewayTypeLDAPAllOfNewUserLookup(
		attributeMappings,
		p.LDAPFilterPattern.ValueString(),
		*management.NewGatewayTypeLDAPAllOfNewUserLookupPopulation(p.PopulationId.ValueString()),
	)

	return data, diags
}

func (p *gatewayResourceModel) toState(apiObject *management.CreateGateway201Response) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)
		return diags
	}

	var d diag.Diagnostics
	fieldObject := apiObject.GetActualInstance()

	switch t := fieldObject.(type) {
	case *management.Gateway:
		p.Id = framework.PingOneResourceIDToTF(t.GetId())
		p.EnvironmentId = framework.PingOneResourceIDToTF(*t.GetEnvironment().Id)
		p.Name = framework.StringOkToTF(t.GetNameOk())
		p.Description = framework.StringOkToTF(t.GetDescriptionOk())
		p.Type = framework.EnumOkToTF(t.GetTypeOk())
		p.Enabled = framework.BoolOkToTF(t.GetEnabledOk())

		// LDAP
		p.BindDN = types.StringNull()
		p.BindPassword = types.StringNull()
		p.ConnectionSecurity = types.StringNull()
		p.FollowReferrals = types.BoolNull()
		p.Kerberos = types.ObjectNull(gatewayKerberosTFObjectTypes)
		p.Servers = types.SetNull(types.StringType)
		p.ValidateTLSCertificates = types.BoolNull()
		p.Vendor = types.StringNull()
		p.UserTypes = types.SetNull(types.ObjectType{AttrTypes: gatewayUserTypesTFObjectTypes})

		// Radius
		p.RadiusDavinciPolicyId = pingonetypes.NewResourceIDNull()
		p.RadiusDefaultSharedSecret = types.StringNull()
		p.RadiusClients = types.SetNull(types.ObjectType{AttrTypes: gatewayRadiusClientsTFObjectTypes})
		p.RadiusNetworkPolicyServer = types.ObjectNull(gatewayRadiusNetworkPolicyServerTFObjectTypes)

	case *management.GatewayTypeLDAP:
		p.Id = framework.PingOneResourceIDToTF(t.GetId())
		p.EnvironmentId = framework.PingOneResourceIDToTF(*t.GetEnvironment().Id)
		p.Name = framework.StringOkToTF(t.GetNameOk())
		p.Description = framework.StringOkToTF(t.GetDescriptionOk())
		p.Type = framework.EnumOkToTF(t.GetTypeOk())
		p.Enabled = framework.BoolOkToTF(t.GetEnabledOk())

		// LDAP
		p.BindDN = framework.StringOkToTF(t.GetBindDNOk())
		p.BindPassword = framework.StringOkToTF(t.GetBindPasswordOk())
		p.ConnectionSecurity = framework.EnumOkToTF(t.GetConnectionSecurityOk())
		p.FollowReferrals = framework.BoolOkToTF(t.GetFollowReferralsOk())
		p.Kerberos, d = toStateKerberosOk(t.GetKerberosOk())
		diags.Append(d...)

		p.Servers = framework.StringSetOkToTF(t.GetServersHostAndPortOk())

		p.ValidateTLSCertificates = framework.BoolOkToTF(t.GetValidateTlsCertificatesOk())
		p.Vendor = framework.EnumOkToTF(t.GetVendorOk())

		p.UserTypes, d = toStateUserTypesOk(t.GetUserTypesOk())
		diags.Append(d...)

		// Radius
		p.RadiusDavinciPolicyId = pingonetypes.NewResourceIDNull()
		p.RadiusDefaultSharedSecret = types.StringNull()
		p.RadiusClients = types.SetNull(types.ObjectType{AttrTypes: gatewayRadiusClientsTFObjectTypes})
		p.RadiusNetworkPolicyServer = types.ObjectNull(gatewayRadiusNetworkPolicyServerTFObjectTypes)

	case *management.GatewayTypeRADIUS:
		p.Id = framework.PingOneResourceIDToTF(t.GetId())
		p.EnvironmentId = framework.PingOneResourceIDToTF(*t.GetEnvironment().Id)
		p.Name = framework.StringOkToTF(t.GetNameOk())
		p.Description = framework.StringOkToTF(t.GetDescriptionOk())
		p.Type = framework.EnumOkToTF(t.GetTypeOk())
		p.Enabled = framework.BoolOkToTF(t.GetEnabledOk())

		// LDAP
		p.BindDN = types.StringNull()
		p.BindPassword = types.StringNull()
		p.ConnectionSecurity = types.StringNull()
		p.FollowReferrals = types.BoolNull()
		p.Kerberos = types.ObjectNull(gatewayKerberosTFObjectTypes)
		p.Servers = types.SetNull(types.StringType)
		p.ValidateTLSCertificates = types.BoolNull()
		p.Vendor = types.StringNull()
		p.UserTypes = types.SetNull(types.ObjectType{AttrTypes: gatewayUserTypesTFObjectTypes})

		// Radius
		if dv, ok := t.GetDavinciOk(); ok {
			if policy, ok := dv.GetPolicyOk(); ok {
				p.RadiusDavinciPolicyId = framework.PingOneResourceIDOkToTF(policy.GetIdOk())
			}
		}
		p.RadiusDefaultSharedSecret = framework.StringOkToTF(t.GetDefaultSharedSecretOk())
		p.RadiusClients, d = toStateRadiusClientOk(t.GetRadiusClientsOk())
		diags.Append(d...)
		p.RadiusNetworkPolicyServer, d = toStateRadiusNetworkPolicyServerOk(t.GetNetworkPolicyServerOk())
		diags.Append(d...)
	}

	return diags
}

func toStateRadiusClientOk(apiObject []management.GatewayTypeRADIUSAllOfRadiusClients, ok bool) (types.Set, diag.Diagnostics) {
	var diags diag.Diagnostics

	tfObjType := types.ObjectType{AttrTypes: gatewayRadiusClientsTFObjectTypes}

	if !ok || apiObject == nil {
		return types.SetNull(tfObjType), diags
	}

	objectList := []attr.Value{}
	for _, client := range apiObject {

		o := map[string]attr.Value{
			"ip":            framework.StringOkToTF(client.GetIpOk()),
			"shared_secret": framework.StringOkToTF(client.GetSharedSecretOk()),
		}

		objValue, d := types.ObjectValue(gatewayRadiusClientsTFObjectTypes, o)
		diags.Append(d...)

		objectList = append(objectList, objValue)
	}

	returnVar, d := types.SetValue(tfObjType, objectList)
	diags.Append(d...)

	return returnVar, diags
}

func toStateRadiusNetworkPolicyServerOk(apiObject *management.GatewayTypeRADIUSAllOfNetworkPolicyServer, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(gatewayRadiusNetworkPolicyServerTFObjectTypes), diags
	}

	o := map[string]attr.Value{
		"ip":   framework.StringOkToTF(apiObject.GetIpOk()),
		"port": framework.Int32OkToTF(apiObject.GetPortOk()),
	}

	returnVar, d := types.ObjectValue(gatewayRadiusNetworkPolicyServerTFObjectTypes, o)
	diags.Append(d...)

	return returnVar, diags
}

func toStateKerberosOk(apiObject *management.GatewayTypeLDAPAllOfKerberos, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(gatewayKerberosTFObjectTypes), diags
	}

	o := map[string]attr.Value{
		"service_account_password":         framework.StringOkToTF(apiObject.GetServiceAccountPasswordOk()),
		"service_account_upn":              framework.StringOkToTF(apiObject.GetServiceAccountUserPrincipalNameOk()),
		"retain_previous_credentials_mins": framework.Int32OkToTF(apiObject.GetMinutesToRetainPreviousCredentialsOk()),
	}

	returnVar, d := types.ObjectValue(gatewayKerberosTFObjectTypes, o)
	diags.Append(d...)

	return returnVar, diags
}

func toStateUserTypesOk(apiObject []management.GatewayTypeLDAPAllOfUserTypes, ok bool) (types.Set, diag.Diagnostics) {
	var diags, d diag.Diagnostics

	tfObjType := types.ObjectType{AttrTypes: gatewayUserTypesTFObjectTypes}

	if !ok || apiObject == nil {
		return types.SetNull(tfObjType), diags
	}

	objectList := []attr.Value{}
	for _, userType := range apiObject {

		o := map[string]attr.Value{
			"id":                            framework.PingOneResourceIDOkToTF(userType.GetIdOk()),
			"name":                          framework.StringOkToTF(userType.GetNameOk()),
			"password_authority":            framework.EnumOkToTF(userType.GetPasswordAuthorityOk()),
			"search_base_dn":                framework.StringOkToTF(userType.GetSearchBaseDnOk()),
			"user_link_attributes":          framework.StringListOkToTF(userType.GetOrderedCorrelationAttributesOk()),
			"push_password_changes_to_ldap": framework.BoolOkToTF(userType.GetAllowPasswordChangesOk()),
			"allow_password_changes":        framework.BoolOkToTF(userType.GetAllowPasswordChangesOk()),
			"update_user_on_successful_authentication": framework.BoolOkToTF(userType.GetUpdateUserOnSuccessfulAuthenticationOk()),
		}

		o["new_user_lookup"], d = toStateUserTypesNewUserLookupOk(userType.GetNewUserLookupOk())
		diags.Append(d...)

		objValue, d := types.ObjectValue(gatewayUserTypesTFObjectTypes, o)
		diags.Append(d...)

		objectList = append(objectList, objValue)
	}

	returnVar, d := types.SetValue(tfObjType, objectList)
	diags.Append(d...)

	return returnVar, diags
}

func toStateUserTypesNewUserLookupOk(apiObject *management.GatewayTypeLDAPAllOfNewUserLookup, ok bool) (types.Object, diag.Diagnostics) {
	var diags, d diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(gatewayUserTypesNewUserLookupTFObjectTypes), diags
	}

	o := map[string]attr.Value{
		"ldap_filter_pattern": framework.StringOkToTF(apiObject.GetLdapFilterPatternOk()),
	}

	o["attribute_mappings"], d = toStateUserTypesNewUserLookupAttributeMappingsOk(apiObject.GetAttributeMappingsOk())
	diags.Append(d...)

	if v, ok := apiObject.GetPopulationOk(); ok {
		o["population_id"] = framework.PingOneResourceIDOkToTF(v.GetIdOk())
	}

	returnVar, d := types.ObjectValue(gatewayUserTypesNewUserLookupTFObjectTypes, o)
	diags.Append(d...)

	return returnVar, diags
}

func toStateUserTypesNewUserLookupAttributeMappingsOk(apiObject []management.GatewayTypeLDAPAllOfNewUserLookupAttributeMappings, ok bool) (types.Set, diag.Diagnostics) {
	var diags, d diag.Diagnostics

	tfObjType := types.ObjectType{AttrTypes: gatewayUserTypesNewUserLookupAttributeMappingTFObjectTypes}

	if !ok || apiObject == nil {
		return types.SetNull(tfObjType), diags
	}

	objectList := []attr.Value{}
	for _, userType := range apiObject {

		o := map[string]attr.Value{
			"name":  framework.StringOkToTF(userType.GetNameOk()),
			"value": framework.StringOkToTF(userType.GetValueOk()),
		}

		objValue, d := types.ObjectValue(gatewayUserTypesNewUserLookupAttributeMappingTFObjectTypes, o)
		diags.Append(d...)

		objectList = append(objectList, objValue)
	}

	returnVar, d := types.SetValue(tfObjType, objectList)
	diags.Append(d...)

	return returnVar, diags
}

var (
	gatewayWriteErrors = func(error model.P1Error) diag.Diagnostics {
		var diags diag.Diagnostics

		// Invalid shared secret combination
		if details, ok := error.GetDetailsOk(); ok && details != nil && len(details) > 0 {
			if code, ok := details[0].GetCodeOk(); ok && *code == "INVALID_VALUE" {
				diags.AddError(
					"Invalid Value",
					details[0].GetMessage(),
				)

				return diags
			}
		}

		return framework.DefaultCustomError(error)
	}
)
