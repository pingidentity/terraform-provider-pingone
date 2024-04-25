package base

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/mfa"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type GatewayResource serviceClientType

type GatewayResourceModel struct {
	Id            pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	Name          types.String                 `tfsdk:"name"`
	Description   types.String                 `tfsdk:"description"`
	Type          types.String                 `tfsdk:"type"`
	Enabled       types.Bool                   `tfsdk:"enabled"`

	// LDAP
	BindDN                                types.String `tfsdk:"bind_dn"`
	BindPassword                          types.String `tfsdk:"bind_password"`
	ConnectionSecurity                    types.String `tfsdk:"connection_security"`
	KerberosServiceAccountPassword        types.String `tfsdk:"kerberos_service_account_password"`
	KerberosServiceAccountUPN             types.String `tfsdk:"kerberos_service_account_upn"`
	KerberosRetainPreviousCredentialsMins types.String `tfsdk:"kerberos_retain_previous_credentials_mins"`
	Servers                               types.Set    `tfsdk:"servers"`
	ValidateTLSCertificates               types.Bool   `tfsdk:"validate_tls_certificates"`
	Vendor                                types.String `tfsdk:"vendor"`
	UserTypes                             types.Set    `tfsdk:"user_types"`

	// Radius
	RadiusDavinciPolicyId     pingonetypes.ResourceIDValue `tfsdk:"radius_davinci_policy_id"`
	RadiusDefaultSharedSecret types.String                 `tfsdk:"radius_default_shared_secret"`
	RadiusClient              types.Set                    `tfsdk:"radius_client"`
}

type GatewayUserTypeResourceModel struct {
	Id                        pingonetypes.ResourceIDValue `tfsdk:"id"`
	Name                      types.String                 `tfsdk:"name"`
	PasswordAuthority         types.String                 `tfsdk:"password_authority"`
	SearchBaseDN              types.String                 `tfsdk:"search_base_dn"`
	UserLinkAttributes        types.String                 `tfsdk:"user_link_attributes"`
	UserMigration             types.Object                 `tfsdk:"user_migration"`
	PushPasswordChangesToLDAP types.Bool                   `tfsdk:"push_password_changes_to_ldap"`
}

type GatewayUserTypeMigrationResourceModel struct {
	LookupFilterPattern types.String                 `tfsdk:"lookup_filter_pattern"`
	PopulationId        pingonetypes.ResourceIDValue `tfsdk:"population_id"`
	AttributeMapping    types.Set                    `tfsdk:"attribute_mapping"`
}

type GatewayUserTypeMigrationAttributeMappingResourceModel struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

type GatewayRadiusClientResourceModel struct {
	IP           types.String `tfsdk:"ip"`
	SharedSecret types.String `tfsdk:"shared_secret"`
}

var (
	GatewayUserTypesTFObjectTypes = map[string]attr.Type{
		"id":                            pingonetypes.ResourceIDType{},
		"name":                          types.StringType,
		"password_authority":            types.StringType,
		"search_base_dn":                types.StringType,
		"user_link_attributes":          types.StringType,
		"user_migration":                types.ObjectType{},
		"push_password_changes_to_ldap": types.BoolType,
	}

	GatewayUserTypesMigrationTFObjectTypes = map[string]attr.Type{
		"lookup_filter_pattern": types.StringType,
		"population_id":         pingonetypes.ResourceIDType{},
		"attribute_mapping":     types.SetType{},
	}

	GatewayUserTypesMigrationAttributeMappingTFObjectTypes = map[string]attr.Type{
		"name":  types.StringType,
		"value": types.StringType,
	}

	GatewayRadiusClientTFObjectTypes = map[string]attr.Type{
		"ip":            types.StringType,
		"shared_secret": types.StringType,
	}
)

// Framework interfaces
var (
	_ resource.Resource                = &GatewayResource{}
	_ resource.ResourceWithConfigure   = &GatewayResource{}
	_ resource.ResourceWithImportState = &GatewayResource{}
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

	kerberosServiceAccountUpnDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"For LDAP gateways only: A string that specifies the Kerberos service account user principal name (for example, `username@bxretail.org`).",
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

	userTypesPasswordAuthorityDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the password authority for the user type.",
	).AllowedValuesEnum(management.AllowedEnumGatewayPasswordAuthorityEnumValues).AppendMarkdownString(fmt.Sprintf("If set to `%s`, PingOne authenticates with the external directory initially, then PingOne authenticates all subsequent sign-ons.", string(management.ENUMGATEWAYPASSWORDAUTHORITY_PING_ONE)))

	userTypesUserMigrationLookupFilterPatternDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The LDAP user search filter to use to match users against the entered user identifier at login. For example, `(((uid=${identifier})(mail=${identifier}))`. Alternatively, this can be a search against the user directory.",
	)

	userTypesUserMigrationLookupAttributeMappingDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"",
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

				Default: stringdefault.StaticString(string(management.ENUMGATEWAYTYPELDAPSECURITY_NONE)),

				Validators: []validator.String{
					stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumGatewayTypeLDAPSecurityEnumValues)...),
					stringvalidator.AlsoRequires(ldapRequiredSchemaPaths...),
				},
			},

			"kerberos_service_account_password": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("For LDAP gateways only: A string that specifies the password for the Kerberos service account.").Description,
				Optional:    true,
				Sensitive:   true,

				Validators: []validator.String{
					stringvalidator.AlsoRequires(ldapRequiredSchemaPaths...),
				},
			},

			"kerberos_service_account_upn": schema.StringAttribute{
				Description:         kerberosServiceAccountUpnDescription.Description,
				MarkdownDescription: kerberosServiceAccountUpnDescription.MarkdownDescription,
				Optional:            true,

				Validators: []validator.String{
					stringvalidator.AlsoRequires(ldapRequiredSchemaPaths...),
				},
			},

			"kerberos_retain_previous_credentials_mins": schema.Int64Attribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("For LDAP gateways only: An integer that specifies the number of minutes for which the previous credentials are persisted.").Description,
				Optional:    true,

				Validators: []validator.Int64{
					int64validator.AlsoRequires(ldapRequiredSchemaPaths...),
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

				Default: booldefault.StaticBool(true),

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
						"id": schema.StringAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
							Computed:    true,
						},

						"name": schema.StringAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("").Description,
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

						"user_link_attributes": schema.SetAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("A list of strings that represent LDAP attribute names that uniquely identify the user, and link to users in PingOne.").Description,
							Required:    true,

							ElementType: types.StringType,
						},

						"user_migration": schema.SingleNestedAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that describes the configurations for initially authenticating new users who will be migrated to PingOne. Note: If there are multiple users having the same user name, only the first user processed is provisioned.").Description,
							Optional:    true,

							Attributes: map[string]schema.Attribute{
								"lookup_filter_pattern": schema.StringAttribute{
									Description:         userTypesUserMigrationLookupFilterPatternDescription.Description,
									MarkdownDescription: userTypesUserMigrationLookupFilterPatternDescription.MarkdownDescription,
									Required:            true,
								},

								"population_id": schema.StringAttribute{
									Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the ID of the population to use to create user entries during lookup.  Must be a valid PingOne resource ID.").Description,
									Required:    true,

									CustomType: pingonetypes.ResourceIDType{},
								},

								"attribute_mapping": schema.SetNestedAttribute{
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

			"radius_client": schema.SetNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("For RADIUS gateways only: A set of objects describing RADIUS client connections.").Description,
				Optional:    true,

				Validators: []validator.Set{
					setvalidator.AlsoRequires(radiusRequiredSchemaPaths...),
				},

				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"ip": schema.StringAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the IP address of the RADIUS client.").Description,
							Required:    true,

							Validators: []validator.String{ipv4},
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
		},
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
	var plan, state GatewayResourceModel

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

func (r *GatewayResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *GatewayResourceModel

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
	var plan, state GatewayResourceModel

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

func (r *GatewayResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *GatewayResourceModel

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

	var diags diag.Diagnostics

	gatewayRequest := &management.CreateGatewayRequest{}

	gatewayType := management.EnumGatewayType(d.Get("type").(string))

	diags = append(diags, checkIllegalParamsForGatewayType(d, gatewayType)...)
	if diags.HasError() {
		return nil, diags
	}

	if slices.Contains([]management.EnumGatewayType{
		"PING_FEDERATE",
		"PING_INTELLIGENCE",
		"API_GATEWAY_INTEGRATION",
	}, gatewayType) {

		gateway := *management.NewGateway(d.Get("name").(string), gatewayType, d.Get("enabled").(bool)) // Gateway |  (optional)

		if v, ok := d.GetOk("description"); ok {
			gateway.SetDescription(v.(string))
		}

		gatewayRequest.Gateway = &gateway

	} else if gatewayType == management.ENUMGATEWAYTYPE_LDAP {

		serversHostAndPort := make([]string, 0)

		if v, ok := d.GetOk("servers"); ok {
			if c := v.(*schema.Set).List(); len(c) > 0 && c[0] != "" {

				for _, str := range v.(*schema.Set).List() {
					serversHostAndPort = append(serversHostAndPort, str.(string))
				}

			}
		}

		gateway := *management.NewGatewayTypeLDAP(
			d.Get("name").(string),
			gatewayType,
			d.Get("enabled").(bool),
			d.Get("bind_dn").(string),
			d.Get("bind_password").(string),
			serversHostAndPort,
			management.EnumGatewayVendor(d.Get("vendor").(string)),
		)

		if v, ok := d.GetOk("connection_security"); ok {
			gateway.SetConnectionSecurity(management.EnumGatewayTypeLDAPSecurity(v.(string)))
		}

		if v, ok := d.GetOk("kerberos_service_account_upn"); ok {
			kerberos := management.NewGatewayTypeLDAPAllOfKerberos(v.(string))

			if v1, ok := d.GetOk("kerberos_service_account_password"); ok {
				kerberos.SetServiceAccountPassword(v1.(string))
			}

			if v1, ok := d.GetOk("kerberos_retain_previous_credentials_mins"); ok {
				kerberos.SetMinutesToRetainPreviousCredentials(int32(v1.(int)))
			}

			gateway.SetKerberos(*kerberos)
		}

		if v, ok := d.GetOk("validate_tls_certificates"); ok {
			gateway.SetValidateTlsCertificates(v.(bool))
		} else {
			gateway.SetValidateTlsCertificates(false)
		}

		if v, ok := d.GetOk("user_type"); ok {
			gateway.SetUserTypes(expandLDAPUserTypes(v.(*schema.Set)))
		}

		gatewayRequest.GatewayTypeLDAP = &gateway

	} else if gatewayType == management.ENUMGATEWAYTYPE_RADIUS {

		radiusClients := make([]management.GatewayTypeRADIUSAllOfRadiusClients, 0)

		if v, ok := d.GetOk("radius_client"); ok {
			if c := v.(*schema.Set).List(); len(c) > 0 && c[0] != "" {

				for _, client := range c {
					clientMap := client.(map[string]interface{})
					radiusClientObj := *management.NewGatewayTypeRADIUSAllOfRadiusClients(clientMap["ip"].(string))

					if v, ok := clientMap["shared_secret"].(string); ok && v != "" {
						radiusClientObj.SetSharedSecret(v)
					}

					radiusClients = append(radiusClients, radiusClientObj)
				}

			} else {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Missing RADIUS Clients",
					Detail:   "Ensure that the `radius_client` parameter is set appropriately.",
				})

				return nil, diags
			}
		}

		gateway := *management.NewGatewayTypeRADIUS(
			d.Get("name").(string),
			gatewayType,
			d.Get("enabled").(bool),
			*management.NewGatewayTypeRADIUSAllOfDavinci(*management.NewGatewayTypeRADIUSAllOfDavinciPolicy(d.Get("radius_davinci_policy_id").(string))),
			radiusClients,
		)

		if v, ok := d.GetOk("radius_default_shared_secret"); ok {
			gateway.SetDefaultSharedSecret(v.(string))
		}

		gatewayRequest.GatewayTypeRADIUS = &gateway

	} else {

		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Cannot determine the gateway type",
			Detail:   "Ensure that the `type` parameter is set appropriately.",
		})

		return nil, diags
	}

	return gatewayRequest, diags
}

func checkIllegalParamsForGatewayType(d *schema.ResourceData, gatewayType management.EnumGatewayType) diag.Diagnostics {
	var diags diag.Diagnostics

	attributes := make([]string, 0)

	if gatewayType != management.ENUMGATEWAYTYPE_LDAP {
		attributes = append(attributes, []string{
			"bind_dn",
			"bind_password",
			// "connection_security",
			"kerberos_service_account_password",
			"kerberos_service_account_upn",
			"kerberos_retain_previous_credentials_mins",
			"servers",
			// "validate_tls_certificates",
			"vendor",
			"user_type",
		}...)
	}

	if gatewayType != management.ENUMGATEWAYTYPE_RADIUS {
		attributes = append(attributes, []string{
			"radius_default_shared_secret",
			"radius_davinci_policy_id",
			"radius_client",
		}...)
	}

	for _, attribute := range attributes {
		if _, ok := d.GetOk(attribute); ok {

			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Unexpected parameter %s for %s gateway type.", attribute, string(gatewayType)),
				Detail:   fmt.Sprintf("The parameter %s does not apply to this gateway type.", attribute),
			})

		}

	}

	return diags

}

func expandLDAPUserTypes(c *schema.Set) []management.GatewayTypeLDAPAllOfUserTypes {

	userTypes := make([]management.GatewayTypeLDAPAllOfUserTypes, 0)

	for _, v := range c.List() {
		obj := v.(map[string]interface{})

		orderedCorrelationAttribtues := make([]string, 0)
		for _, str := range obj["user_link_attributes"].([]interface{}) {
			orderedCorrelationAttribtues = append(orderedCorrelationAttribtues, str.(string))
		}

		userType := *management.NewGatewayTypeLDAPAllOfUserTypes(
			obj["name"].(string),
			orderedCorrelationAttribtues,
			management.EnumGatewayPasswordAuthority(obj["password_authority"].(string)),
			obj["search_base_dn"].(string),
		)

		if v, ok := obj["push_password_changes_to_ldap"].(bool); ok {
			userType.SetAllowPasswordChanges(v)
		}

		if v, ok := obj["user_migration"].([]interface{}); ok && len(v) > 0 && v[0] != nil {
			userType.SetNewUserLookup(*expandLDAPUserLookup(v[0].(map[string]interface{})))
		}

		userTypes = append(userTypes, userType)
	}

	return userTypes

}

func expandLDAPUserLookup(c map[string]interface{}) *management.GatewayTypeLDAPAllOfNewUserLookup {

	attributeMappings := expandLDAPUserLookupAttributeMappings(c["attribute_mapping"].(*schema.Set).List())

	userLookup := *management.NewGatewayTypeLDAPAllOfNewUserLookup(
		attributeMappings,
		c["lookup_filter_pattern"].(string),
		*management.NewGatewayTypeLDAPAllOfNewUserLookupPopulation(c["population_id"].(string)),
	)

	return &userLookup

}

func expandLDAPUserLookupAttributeMappings(c []interface{}) []management.GatewayTypeLDAPAllOfNewUserLookupAttributeMappings {
	mappings := make([]management.GatewayTypeLDAPAllOfNewUserLookupAttributeMappings, 0)

	for _, v := range c {

		obj := v.(map[string]interface{})

		mappings = append(mappings, *management.NewGatewayTypeLDAPAllOfNewUserLookupAttributeMappings(obj["name"].(string), obj["value"].(string)))

	}

	return mappings
}

func flattenUserType(c []management.GatewayTypeLDAPAllOfUserTypes) []map[string]interface{} {

	items := make([]map[string]interface{}, 0)

	for _, v := range c {
		// Required
		item := map[string]interface{}{
			"id":                   v.GetId(),
			"name":                 v.GetName(),
			"password_authority":   string(v.GetPasswordAuthority()),
			"search_base_dn":       v.GetSearchBaseDn(),
			"user_link_attributes": v.GetOrderedCorrelationAttributes(),
		}

		// Optional

		if v1, ok := v.GetAllowPasswordChangesOk(); ok {
			item["push_password_changes_to_ldap"] = v1
		} else {
			item["push_password_changes_to_ldap"] = nil
		}

		if v1, ok := v.GetNewUserLookupOk(); ok {

			userMigrationItem := map[string]interface{}{
				"lookup_filter_pattern": v1.GetLdapFilterPattern(),
				"population_id":         v1.GetPopulation().Id,
				"attribute_mapping":     flattenLDAPUserLookupAttributeMappings(v1.GetAttributeMappings()),
			}

			userMigrationItemList := make([]map[string]interface{}, 0)
			item["user_migration"] = append(userMigrationItemList, userMigrationItem)

		} else {
			item["user_migration"] = nil
		}

		items = append(items, item)
	}

	return items
}

func flattenLDAPUserLookupAttributeMappings(c []management.GatewayTypeLDAPAllOfNewUserLookupAttributeMappings) interface{} {
	items := make([]interface{}, 0)

	for _, v := range c {
		items = append(items, map[string]interface{}{
			"name":  v.GetName(),
			"value": v.GetValue(),
		})
	}

	return items
}
