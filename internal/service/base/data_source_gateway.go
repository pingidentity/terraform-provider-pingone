package base

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	validation "github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type GatewayDataSource serviceClientType

type gatewayDataSourceModel struct {
	Id                                    types.String `tfsdk:"id"`
	EnvironmentId                         types.String `tfsdk:"environment_id"`
	GatewayId                             types.String `tfsdk:"gateway_id"`
	Name                                  types.String `tfsdk:"name"`
	Description                           types.String `tfsdk:"description"`
	Enabled                               types.Bool   `tfsdk:"enabled"`
	Type                                  types.String `tfsdk:"type"`
	BindDN                                types.String `tfsdk:"bind_dn"`
	ConnectionSecurity                    types.String `tfsdk:"connection_security"`
	KerberosServiceAccountUPN             types.String `tfsdk:"kerberos_service_account_upn"`
	KerberosRetainPreviousCredentialsMins types.Int64  `tfsdk:"kerberos_retain_previous_credentials_mins"`
	Servers                               types.Set    `tfsdk:"servers"`
	ValidateTLSCertificates               types.Bool   `tfsdk:"validate_tls_certificates"`
	Vendor                                types.String `tfsdk:"vendor"`
	UserType                              types.Set    `tfsdk:"user_type"`
	RadiusDavinciPolicyId                 types.String `tfsdk:"radius_davinci_policy_id"`
	RadiusDefaultSharedSecret             types.String `tfsdk:"radius_default_shared_secret"`
	RadiusClient                          types.Set    `tfsdk:"radius_client"`
}

var (
	gatewayRadiusClientTFObjectTypes = map[string]attr.Type{
		"ip":            types.StringType,
		"shared_secret": types.StringType,
	}

	gatewayLdapUserTypeTFObjectTypes = map[string]attr.Type{
		"id":                            types.StringType,
		"name":                          types.StringType,
		"password_authority":            types.StringType,
		"push_password_changes_to_ldap": types.BoolType,
		"search_base_dn":                types.StringType,
		"user_link_attributes":          types.ListType{ElemType: types.StringType},
		"user_migration":                types.ListType{ElemType: types.ObjectType{AttrTypes: gatewayLdapUserMigrationTFObjectTypes}},
	}

	gatewayLdapUserMigrationTFObjectTypes = map[string]attr.Type{
		"lookup_filter_pattern": types.StringType,
		"population_id":         types.StringType,
		"attribute_mapping":     types.SetType{ElemType: types.ObjectType{AttrTypes: gatewayLdapUserLookupAttributeMappingTFObjectTypes}},
	}

	gatewayLdapUserLookupAttributeMappingTFObjectTypes = map[string]attr.Type{
		"name":  types.StringType,
		"value": types.StringType,
	}
)

// Framework interfaces
var (
	_ datasource.DataSource = &GatewayDataSource{}
)

// New Object
func NewGatewayDataSource() datasource.DataSource {
	return &GatewayDataSource{}
}

// Metadata
func (r *GatewayDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_gateway"
}

func (r *GatewayDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	// schema descriptions and validation settings
	const attrMinLength = 1

	gatewayIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The identifier (UUID) of the gateway.",
	).ExactlyOneOf([]string{"gateway_id", "name"}).AppendMarkdownString("Must be a valid PingOne resource ID.")

	nameDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The name of the gateway.",
	).ExactlyOneOf([]string{"gateway_id", "name"})

	typeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Specifies the type of gateway resource.",
	).AllowedValuesEnum(management.AllowedEnumGatewayTypeEnumValues)

	bindDNDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"For LDAP gateways only: The distinguished name information to bind to the LDAP database (for example, `uid=pingone,dc=bxretail,dc=org`).",
	)

	connectionSecurity := framework.SchemaAttributeDescriptionFromMarkdown(
		"For LDAP gateways only: The connection security type.",
	).AllowedValuesEnum(management.AllowedEnumGatewayTypeLDAPSecurityEnumValues)

	kerberosServiceAccountUPNDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"For LDAP gateways only: The Kerberos service account user principal name (for example, `username@bxretail.org`).",
	)

	serversDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"For LDAP gateways only: A list of LDAP server host name and port number combinations (for example, [`ds1.bxretail.org:636`, `ds2.bxretail.org:636`]).",
	)

	validateTLSCertificatesDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"For LDAP gateways only: Indicates whether or not to trust all SSL certificates (defaults to `true`). If this value is `false`, TLS certificates are not validated. When the value is set to `true`, only certificates that are signed by the default JVM CAs, or the CA certs that the customer has uploaded to the certificate service are trusted.",
	)

	vendorDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"For LDAP gateways only: The LDAP vendor",
	).AllowedValuesEnum(management.AllowedEnumGatewayVendorEnumValues)

	userTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"For LDAP gateways only: A collection of properties that define how users should be provisioned in PingOne. The `user_type` block specifies which user properties in PingOne correspond to the user properties in an external LDAP directory. You can use an LDAP browser to view the user properties in the external LDAP directory.",
	)

	userTypeIdsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Identifies the user type. This correlates to the `password.external.gateway.user_type.id` User property.",
	)

	userTypePasswordAuthorityDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("This can be either `%s` or `%s`. If set to `%s`, PingOne authenticates with the external directory initially, then PingOne authenticates all subsequent sign-ons.", string(management.ENUMGATEWAYPASSWORDAUTHORITY_PING_ONE), string(management.ENUMGATEWAYPASSWORDAUTHORITY_LDAP), string(management.ENUMGATEWAYPASSWORDAUTHORITY_PING_ONE)),
	)

	userMigrationLookupFilterPatternDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The LDAP user search filter to use to match users against the entered user identifier at login. For example, `(((uid=${identifier})(mail=${identifier}))`. Alternatively, this can be a search against the user directory.",
	)

	userMigrationAttributeMappingDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A collection of properties that define how users should be provisioned in PingOne. The `user_type` block specifies which user properties in PingOne correspond to the user properties in an external LDAP directory. You can use an LDAP browser to view the user properties in the external LDAP directory.",
	)

	userMigrationAttributeMappingNameDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The name of a user attribute in PingOne. See [Users properties](https://apidocs.pingidentity.com/pingone/platform/v1/api/#users) for the complete list of available PingOne user attributes.",
	)

	userMigrationAttributeMappingValueDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A reference to the corresponding external LDAP attribute.  Values are in the format `${ldapAttributes.mail}`, while Terraform HCL requires an additional `$` prefix character. For example, `$${ldapAttributes.mail}`",
	)

	radiusClientSharedSecretDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The shared secret for the RADIUS client. If this value is not provided, the shared secret specified with `default_shared_secret` is used.",
	)

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Data source to retrieve a PingOne gateway in an environment from ID or by name.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),
			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("PingOne environment identifier (UUID) in which the gateway exists."),
			),
			"gateway_id": schema.StringAttribute{
				Description:         gatewayIdDescription.Description,
				MarkdownDescription: gatewayIdDescription.MarkdownDescription,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName("name"),
					),
					validation.P1ResourceIDValidator(),
				},
			},
			"name": schema.StringAttribute{
				Description:         nameDescription.Description,
				MarkdownDescription: nameDescription.MarkdownDescription,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName("gateway_id"),
					),
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},
			"description": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the description of the gateway.").Description,
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether the gateway is enabled in the environment.").Description,
				Computed:    true,
			},
			"type": schema.StringAttribute{
				Description:         typeDescription.Description,
				MarkdownDescription: typeDescription.MarkdownDescription,
				Computed:            true,
			},

			// LDAP
			"bind_dn": schema.StringAttribute{
				Description:         bindDNDescription.Description,
				MarkdownDescription: bindDNDescription.MarkdownDescription,
				Computed:            true,
			},
			"connection_security": schema.StringAttribute{
				Description:         connectionSecurity.Description,
				MarkdownDescription: connectionSecurity.MarkdownDescription,
				Computed:            true,
			},
			"kerberos_service_account_upn": schema.StringAttribute{
				Description:         kerberosServiceAccountUPNDescription.Description,
				MarkdownDescription: kerberosServiceAccountUPNDescription.MarkdownDescription,
				Computed:            true,
			},
			"kerberos_retain_previous_credentials_mins": schema.Int64Attribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("For LDAP gateways only: The number of minutes for which the previous credentials are persisted.").Description,
				Computed:    true,
			},
			"servers": schema.SetAttribute{
				Description:         serversDescription.Description,
				MarkdownDescription: serversDescription.MarkdownDescription,
				ElementType:         types.StringType,
				Computed:            true,
			},
			"validate_tls_certificates": schema.BoolAttribute{
				Description:         validateTLSCertificatesDescription.Description,
				MarkdownDescription: validateTLSCertificatesDescription.MarkdownDescription,
				Computed:            true,
			},
			"vendor": schema.StringAttribute{
				Description:         vendorDescription.Description,
				MarkdownDescription: vendorDescription.MarkdownDescription,
				Computed:            true,
			},
			"user_type": schema.SetNestedAttribute{
				Description:         userTypeDescription.Description,
				MarkdownDescription: userTypeDescription.MarkdownDescription,
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description:         userTypeIdsDescription.Description,
							MarkdownDescription: userTypeIdsDescription.MarkdownDescription,
							Computed:            true,
						},
						"name": schema.StringAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("The name of the user type.").Description,
							Computed:    true,
						},
						"password_authority": schema.StringAttribute{
							Description:         userTypePasswordAuthorityDescription.Description,
							MarkdownDescription: userTypePasswordAuthorityDescription.MarkdownDescription,
							Computed:            true,
						},
						"push_password_changes_to_ldap": schema.BoolAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("Determines whether password updates in PingOne should be pushed to the user's record in LDAP.  If false, the user cannot change the password and have it updated in the remote LDAP directory. In this case, operations for forgotten passwords or resetting of passwords are not available to a user referencing this gateway.").Description,
							Computed:    true,
						},
						"search_base_dn": schema.StringAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("The LDAP base domain name (DN) for this user type.").Description,
							Computed:    true,
						},
						"user_link_attributes": schema.ListAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("Represents LDAP attribute names that uniquely identify the user, and link to users in PingOne.").Description,
							ElementType: types.StringType,
							Computed:    true,
						},
						"user_migration": schema.ListNestedAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("The configurations for initially authenticating new users who will be migrated to PingOne. Note: If there are multiple users having the same user name, only the first user processed is provisioned.").Description,
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"lookup_filter_pattern": schema.StringAttribute{
										Description:         userMigrationLookupFilterPatternDescription.Description,
										MarkdownDescription: userMigrationLookupFilterPatternDescription.MarkdownDescription,
										Computed:            true,
									},
									"population_id": schema.StringAttribute{
										Description: framework.SchemaAttributeDescriptionFromMarkdown("The ID of the population to use to create user entries during lookup.").Description,
										Computed:    true,
									},
									"attribute_mapping": schema.SetNestedAttribute{
										Description:         userMigrationAttributeMappingDescription.Description,
										MarkdownDescription: userMigrationAttributeMappingDescription.MarkdownDescription,
										Computed:            true,
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"name": schema.StringAttribute{
													Description:         userMigrationAttributeMappingNameDescription.Description,
													MarkdownDescription: userMigrationAttributeMappingNameDescription.MarkdownDescription,
													Computed:            true,
												},
												"value": schema.StringAttribute{
													Description:         userMigrationAttributeMappingValueDescription.Description,
													MarkdownDescription: userMigrationAttributeMappingValueDescription.MarkdownDescription,
													Computed:            true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},

			// RADIUS
			"radius_davinci_policy_id": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("For RADIUS gateways only: The ID of the DaVinci flow policy to use.").Description,
				Computed:    true,
			},
			"radius_default_shared_secret": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("For RADIUS gateways only: Value to use for the shared secret if the shared secret is not provided for one or more of the RADIUS clients specified.").Description,
				Computed:    true,
				Sensitive:   true,
			},
			"radius_client": schema.SetNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("For RADIUS gateways only: A collection of RADIUS clients.").Description,
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"ip": schema.StringAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("The IP of the RADIUS client.").Description,
							Computed:    true,
						},
						"shared_secret": schema.StringAttribute{
							Description:         radiusClientSharedSecretDescription.Description,
							MarkdownDescription: radiusClientSharedSecretDescription.MarkdownDescription,
							Computed:            true,
							Sensitive:           true,
						},
					},
				},
			},
		},
	}
}

func (r *GatewayDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *GatewayDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *gatewayDataSourceModel

	if r.Client.ManagementAPIClient == nil {
		resp.Diagnostics.AddError(
			"Client not initialized",
			"Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
		return
	}

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var gatewayInstance interface{}

	// Gateway API does not support SCIM filtering
	if !data.GatewayId.IsNull() {
		// Run the API call
		var response *management.CreateGateway201Response
		resp.Diagnostics.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				fO, fR, fErr := r.Client.ManagementAPIClient.GatewaysApi.ReadOneGateway(ctx, data.EnvironmentId.ValueString(), data.GatewayId.ValueString()).Execute()
				return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
			},
			"ReadOneGateway",
			framework.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			&response,
		)...)
		if resp.Diagnostics.HasError() {
			return
		}

		gatewayInstance = response.GetActualInstance()

	} else if !data.Name.IsNull() {
		// Run the API call
		var entityArray *management.EntityArray
		resp.Diagnostics.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				fO, fR, fErr := r.Client.ManagementAPIClient.GatewaysApi.ReadAllGateways(ctx, data.EnvironmentId.ValueString()).Execute()
				return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
			},
			"ReadAllGateways",
			framework.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			&entityArray,
		)...)
		if resp.Diagnostics.HasError() {
			return
		}

		if gateways, ok := entityArray.Embedded.GetGatewaysOk(); ok {
			found := false

			for _, gatewayObject := range gateways {
				if gateway := gatewayObject.Gateway; gateway != nil && gateway.GetId() != "" && gateway.GetName() == data.Name.ValueString() {
					gatewayInstance = gateway
					found = true
					break

				} else if gateway := gatewayObject.GatewayTypeLDAP; gateway != nil && gateway.GetId() != "" && gateway.GetName() == data.Name.ValueString() {
					gatewayInstance = gateway
					found = true
					break

				} else if gateway := gatewayObject.GatewayTypeRADIUS; gateway != nil && gateway.GetId() != "" && gateway.GetName() == data.Name.ValueString() {
					gatewayInstance = gateway
					found = true
					break

				}
			}

			if !found {
				resp.Diagnostics.AddError(
					"Cannot find the gateway from name",
					fmt.Sprintf("The gateway name %s for environment %s cannot be found", data.Name.String(), data.EnvironmentId.String()),
				)
				return
			}

		}
	} else {
		resp.Diagnostics.AddError(
			"Missing parameter",
			"Cannot find the requested PingOne Gateway: gateway_id or name argument must be set.",
		)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(gatewayInstance)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (p *gatewayDataSourceModel) toState(apiObject interface{}) diag.Diagnostics {
	var diags, d diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	switch v := apiObject.(type) {
	case *management.Gateway:
		p.Id = framework.StringOkToTF(v.GetIdOk())
		p.EnvironmentId = framework.StringToTF(*v.GetEnvironment().Id)
		p.GatewayId = framework.StringOkToTF(v.GetIdOk())
		p.Name = framework.StringOkToTF(v.GetNameOk())
		p.Description = framework.StringOkToTF(v.GetDescriptionOk())
		p.Enabled = framework.BoolOkToTF(v.GetEnabledOk())
		p.Type = framework.EnumOkToTF(v.GetTypeOk())

	case *management.GatewayTypeLDAP:
		p.Id = framework.StringOkToTF(v.GetIdOk())
		p.EnvironmentId = framework.StringToTF(*v.GetEnvironment().Id)
		p.GatewayId = framework.StringOkToTF(v.GetIdOk())
		p.Name = framework.StringOkToTF(v.GetNameOk())
		p.Description = framework.StringOkToTF(v.GetDescriptionOk())
		p.Enabled = framework.BoolOkToTF(v.GetEnabledOk())
		p.Type = framework.EnumOkToTF(v.GetTypeOk())
		p.BindDN = framework.StringOkToTF(v.GetBindDNOk())
		p.Vendor = framework.EnumOkToTF(v.GetVendorOk())
		p.ConnectionSecurity = framework.EnumOkToTF(v.GetConnectionSecurityOk())
		p.Servers = framework.StringSetOkToTF(v.GetServersHostAndPortOk())
		p.ValidateTLSCertificates = framework.BoolOkToTF(v.GetValidateTlsCertificatesOk())

		if v1, ok := v.GetKerberosOk(); ok {
			p.KerberosServiceAccountUPN = framework.StringOkToTF(v1.GetServiceAccountUserPrincipalNameOk())
			p.KerberosRetainPreviousCredentialsMins = framework.Int32OkToTF(v1.GetMinutesToRetainPreviousCredentialsOk())
		} else {
			p.KerberosServiceAccountUPN = types.StringNull()
			p.KerberosRetainPreviousCredentialsMins = types.Int64Null()
		}

		p.UserType, d = p.gatewayLdapUserTypeOkToTF(v.GetUserTypesOk())
		diags.Append(d...)

	case *management.GatewayTypeRADIUS:
		p.Id = framework.StringOkToTF(v.GetIdOk())
		p.EnvironmentId = framework.StringToTF(*v.GetEnvironment().Id)
		p.GatewayId = framework.StringOkToTF(v.GetIdOk())
		p.Name = framework.StringOkToTF(v.GetNameOk())
		p.Description = framework.StringOkToTF(v.GetDescriptionOk())
		p.Enabled = framework.BoolOkToTF(v.GetEnabledOk())
		p.Type = framework.EnumOkToTF(v.GetTypeOk())

		if v1, ok := v.GetDavinciOk(); ok {
			p.RadiusDavinciPolicyId = framework.StringToTF(v1.GetPolicy().Id)
		} else {
			p.RadiusDavinciPolicyId = types.StringNull()
		}

		p.RadiusDefaultSharedSecret = framework.StringOkToTF(v.GetDefaultSharedSecretOk())

		p.RadiusClient, d = p.gatewayRadiusClientOkToTF(v.GetRadiusClientsOk())
		diags.Append(d...)

	default:
		diags.AddError(
			"Undefined gateway type",
			"Cannot identify the gateway type from the data object.  Please report this to the provider maintainers.",
		)

	}

	return diags
}
func (p *gatewayDataSourceModel) gatewayLdapUserTypeOkToTF(apiObject []management.GatewayTypeLDAPAllOfUserTypes, ok bool) (basetypes.SetValue, diag.Diagnostics) {
	var diags, d diag.Diagnostics
	tfObjType := types.ObjectType{AttrTypes: gatewayLdapUserTypeTFObjectTypes}

	if !ok || apiObject == nil {
		return types.SetNull(tfObjType), diags
	}

	userTypes := []attr.Value{}
	for _, v := range apiObject {

		// build user migration object
		userMigration := types.ListNull(types.ObjectType{AttrTypes: gatewayLdapUserMigrationTFObjectTypes})
		if v1, ok := v.GetNewUserLookupOk(); ok {
			attributeMapObj, d := p.gatewayLdapUserLookupAttributeMappingsOkToTF(v1.GetAttributeMappingsOk())
			diags.Append(d...)

			o := map[string]attr.Value{
				"lookup_filter_pattern": framework.StringOkToTF(v1.GetLdapFilterPatternOk()),
				"population_id":         framework.StringToTF(v1.GetPopulation().Id),
				"attribute_mapping":     attributeMapObj,
			}

			flattenedObj, d := types.ObjectValue(gatewayLdapUserMigrationTFObjectTypes, o)
			diags.Append(d...)

			objValue, d := types.ListValue(types.ObjectType{AttrTypes: gatewayLdapUserMigrationTFObjectTypes}, append([]attr.Value{}, flattenedObj))
			diags.Append(d...)

			userMigration = objValue
		}

		userTypesObj, d := types.ObjectValue(gatewayLdapUserTypeTFObjectTypes, map[string]attr.Value{
			"id":                            framework.StringOkToTF(v.GetIdOk()),
			"name":                          framework.StringOkToTF(v.GetNameOk()),
			"password_authority":            framework.EnumOkToTF(v.GetPasswordAuthorityOk()),
			"search_base_dn":                framework.StringOkToTF(v.GetSearchBaseDnOk()),
			"user_link_attributes":          framework.StringListOkToTF(v.GetOrderedCorrelationAttributesOk()),
			"push_password_changes_to_ldap": framework.BoolOkToTF(v.GetAllowPasswordChangesOk()),
			"user_migration":                userMigration,
		})
		diags.Append(d...)

		userTypes = append(userTypes, userTypesObj)
	}

	returnVar, d := types.SetValue(tfObjType, userTypes)
	diags.Append(d...)

	return returnVar, diags
}

func (p *gatewayDataSourceModel) gatewayRadiusClientOkToTF(apiObject []management.GatewayTypeRADIUSAllOfRadiusClients, ok bool) (basetypes.SetValue, diag.Diagnostics) {
	var diags diag.Diagnostics
	tfObjType := types.ObjectType{AttrTypes: gatewayRadiusClientTFObjectTypes}

	if !ok || apiObject == nil {
		return types.SetNull(tfObjType), diags
	}

	radiusClients := []attr.Value{}
	for _, v := range apiObject {

		radiusClient := map[string]attr.Value{
			"ip":            framework.StringOkToTF(v.GetIpOk()),
			"shared_secret": framework.StringOkToTF(v.GetSharedSecretOk()),
		}
		radiusClientsObj, d := types.ObjectValue(gatewayRadiusClientTFObjectTypes, radiusClient)
		diags.Append(d...)

		radiusClients = append(radiusClients, radiusClientsObj)
	}

	returnVar, d := types.SetValue(tfObjType, radiusClients)
	diags.Append(d...)

	return returnVar, diags
}

func (p *gatewayDataSourceModel) gatewayLdapUserLookupAttributeMappingsOkToTF(apiObject []management.GatewayTypeLDAPAllOfNewUserLookupAttributeMappings, ok bool) (basetypes.SetValue, diag.Diagnostics) {
	var diags diag.Diagnostics
	tfObjType := types.ObjectType{AttrTypes: gatewayLdapUserLookupAttributeMappingTFObjectTypes}

	if !ok || apiObject == nil {
		return types.SetNull(tfObjType), diags
	}

	attributeMappings := []attr.Value{}
	for _, v := range apiObject {

		attributeMap := map[string]attr.Value{
			"name":  framework.StringOkToTF(v.GetNameOk()),
			"value": framework.StringOkToTF(v.GetValueOk()),
		}
		attributeMappingObj, d := types.ObjectValue(gatewayLdapUserLookupAttributeMappingTFObjectTypes, attributeMap)
		diags.Append(d...)

		attributeMappings = append(attributeMappings, attributeMappingObj)
	}

	returnVar, d := types.SetValue(tfObjType, attributeMappings)
	diags.Append(d...)

	return returnVar, diags
}
