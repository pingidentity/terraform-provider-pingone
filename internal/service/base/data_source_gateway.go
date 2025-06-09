// Copyright © 2025 Ping Identity Corporation

package base

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/davincitypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/legacysdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
)

// Types
type GatewayDataSource serviceClientType

type gatewayDataSourceModel struct {
	Id            pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	GatewayId     pingonetypes.ResourceIDValue `tfsdk:"gateway_id"`
	Name          types.String                 `tfsdk:"name"`
	Description   types.String                 `tfsdk:"description"`
	Type          types.String                 `tfsdk:"type"`
	Enabled       types.Bool                   `tfsdk:"enabled"`

	// LDAP
	BindDN                  types.String `tfsdk:"bind_dn"`
	ConnectionSecurity      types.String `tfsdk:"connection_security"`
	FollowReferrals         types.Bool   `tfsdk:"follow_referrals"`
	Kerberos                types.Object `tfsdk:"kerberos"`
	Servers                 types.Set    `tfsdk:"servers"`
	ValidateTLSCertificates types.Bool   `tfsdk:"validate_tls_certificates"`
	Vendor                  types.String `tfsdk:"vendor"`
	UserTypes               types.Map    `tfsdk:"user_types"`

	// Radius
	RadiusClients             types.Set                    `tfsdk:"radius_clients"`
	RadiusDavinciPolicyId     davincitypes.ResourceIDValue `tfsdk:"radius_davinci_policy_id"`
	RadiusDefaultSharedSecret types.String                 `tfsdk:"radius_default_shared_secret"`
	RadiusNetworkPolicyServer types.Object                 `tfsdk:"radius_network_policy_server"`
}

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

	followReferralsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"For LDAP gateways only: A boolean that, when set to true, PingOne sends LDAP queries per referrals it receives from the LDAP servers.",
	)

	kerberosServiceAccountUpnDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The Kerberos service account user principal name (for example, `username@bxretail.org`).",
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

	userTypesDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"For LDAP gateways only: A collection of properties that define how users should be provisioned in PingOne. The `user_type` block specifies which user properties in PingOne correspond to the user properties in an external LDAP directory. You can use an LDAP browser to view the user properties in the external LDAP directory.",
	)

	userTypesAllowPasswordChangesDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that, if set to `false`, the user cannot change the password in the remote LDAP directory. In this case, operations for forgotten passwords or resetting of passwords are not available to a user referencing this gateway.",
	)

	userTypesUpdateUserOnSuccessfulAuthenticationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that, if set to `true`, when users sign on through an LDAP Gateway client, user attributes are updated based on responses from the LDAP server.",
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

				CustomType: pingonetypes.ResourceIDType{},

				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName("name"),
					),
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
			"type": schema.StringAttribute{
				Description:         typeDescription.Description,
				MarkdownDescription: typeDescription.MarkdownDescription,
				Computed:            true,
			},
			"enabled": schema.BoolAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether the gateway is enabled in the environment.").Description,
				Computed:    true,
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
			"follow_referrals": schema.BoolAttribute{
				Description:         followReferralsDescription.Description,
				MarkdownDescription: followReferralsDescription.MarkdownDescription,
				Computed:            true,
			},
			"kerberos": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("For LDAP gateways only: A single object that specifies Kerberos connection details.").Description,
				Computed:    true,

				Attributes: map[string]schema.Attribute{
					"service_account_password": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the password for the Kerberos service account.").Description,
						Computed:    true,
						Sensitive:   true,
					},

					"service_account_upn": schema.StringAttribute{
						Description:         kerberosServiceAccountUpnDescription.Description,
						MarkdownDescription: kerberosServiceAccountUpnDescription.MarkdownDescription,
						Computed:            true,
					},

					"retain_previous_credentials_mins": schema.Int32Attribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that specifies the number of minutes for which the previous credentials are persisted.").Description,
						Computed:    true,
					},
				},
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
			"user_types": schema.MapNestedAttribute{
				Description:         userTypesDescription.Description,
				MarkdownDescription: userTypesDescription.MarkdownDescription,
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"allow_password_changes": schema.BoolAttribute{
							Description:         userTypesAllowPasswordChangesDescription.Description,
							MarkdownDescription: userTypesAllowPasswordChangesDescription.MarkdownDescription,
							Computed:            true,
						},
						"update_user_on_successful_authentication": schema.BoolAttribute{
							Description:         userTypesUpdateUserOnSuccessfulAuthenticationDescription.Description,
							MarkdownDescription: userTypesUpdateUserOnSuccessfulAuthenticationDescription.MarkdownDescription,
							Computed:            true,
						},
						"id": schema.StringAttribute{
							Description:         userTypeIdsDescription.Description,
							MarkdownDescription: userTypeIdsDescription.MarkdownDescription,
							Computed:            true,

							CustomType: pingonetypes.ResourceIDType{},
						},
						"password_authority": schema.StringAttribute{
							Description:         userTypePasswordAuthorityDescription.Description,
							MarkdownDescription: userTypePasswordAuthorityDescription.MarkdownDescription,
							Computed:            true,
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
						"new_user_lookup": schema.SingleNestedAttribute{
							Description: framework.SchemaAttributeDescriptionFromMarkdown("The configurations for initially authenticating new users who will be migrated to PingOne. Note: If there are multiple users having the same user name, only the first user processed is provisioned.").Description,
							Computed:    true,

							Attributes: map[string]schema.Attribute{
								"ldap_filter_pattern": schema.StringAttribute{
									Description:         userMigrationLookupFilterPatternDescription.Description,
									MarkdownDescription: userMigrationLookupFilterPatternDescription.MarkdownDescription,
									Computed:            true,
								},
								"population_id": schema.StringAttribute{
									Description: framework.SchemaAttributeDescriptionFromMarkdown("The ID of the population to use to create user entries during lookup.").Description,
									Computed:    true,

									CustomType: pingonetypes.ResourceIDType{},
								},
								"attribute_mappings": schema.SetNestedAttribute{
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

			// RADIUS
			"radius_davinci_policy_id": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("For RADIUS gateways only: The ID of the PingOne DaVinci flow policy to use.").Description,
				Computed:    true,

				CustomType: davincitypes.ResourceIDType{},
			},
			"radius_default_shared_secret": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("For RADIUS gateways only: Value to use for the shared secret if the shared secret is not provided for one or more of the RADIUS clients specified.").Description,
				Computed:    true,
				Sensitive:   true,
			},
			"radius_clients": schema.SetNestedAttribute{
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
			"radius_network_policy_server": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("For RADIUS gateways only: A single object that allows configuration of the RADIUS gateway to authenticate using the MS-CHAP v2 protocol.").Description,
				Computed:    true,

				Attributes: map[string]schema.Attribute{
					"ip": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the IP address of the Network Policy Server (NPS).").Description,
						Computed:    true,
					},

					"port": schema.Int32Attribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that specifies the port number of the NPS.").Description,
						Computed:    true,
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

	if r.Client == nil || r.Client.ManagementAPIClient == nil {
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
		resp.Diagnostics.Append(legacysdk.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				fO, fR, fErr := r.Client.ManagementAPIClient.GatewaysApi.ReadOneGateway(ctx, data.EnvironmentId.ValueString(), data.GatewayId.ValueString()).Execute()
				return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
			},
			"ReadOneGateway",
			legacysdk.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			&response,
		)...)
		if resp.Diagnostics.HasError() {
			return
		}

		gatewayInstance = response.GetActualInstance()

	} else if !data.Name.IsNull() {
		// Run the API call
		var response *management.CreateGateway201Response
		resp.Diagnostics.Append(legacysdk.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				pagedIterator := r.Client.ManagementAPIClient.GatewaysApi.ReadAllGateways(ctx, data.EnvironmentId.ValueString()).Execute()

				var initialHttpResponse *http.Response

				for pageCursor, err := range pagedIterator {
					if err != nil {
						return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, pageCursor.HTTPResponse, err)
					}

					if initialHttpResponse == nil {
						initialHttpResponse = pageCursor.HTTPResponse
					}

					if gateways, ok := pageCursor.EntityArray.Embedded.GetGatewaysOk(); ok {

						for _, gatewayObject := range gateways {
							if gateway := gatewayObject.Gateway; gateway != nil && gateway.GetId() != "" && gateway.GetName() == data.Name.ValueString() {
								return &management.CreateGateway201Response{
									Gateway: gateway,
								}, pageCursor.HTTPResponse, nil

							} else if gateway := gatewayObject.GatewayTypeLDAP; gateway != nil && gateway.GetId() != "" && gateway.GetName() == data.Name.ValueString() {
								return &management.CreateGateway201Response{
									GatewayTypeLDAP: gateway,
								}, pageCursor.HTTPResponse, nil

							} else if gateway := gatewayObject.GatewayTypeRADIUS; gateway != nil && gateway.GetId() != "" && gateway.GetName() == data.Name.ValueString() {
								return &management.CreateGateway201Response{
									GatewayTypeRADIUS: gateway,
								}, pageCursor.HTTPResponse, nil

							}
						}

					}
				}

				return nil, initialHttpResponse, nil
			},
			"ReadAllGateways",
			legacysdk.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			&response,
		)...)
		if resp.Diagnostics.HasError() {
			return
		}

		if response == nil {
			resp.Diagnostics.AddError(
				"Cannot find the gateway from name",
				fmt.Sprintf("The gateway name %s for environment %s cannot be found", data.Name.String(), data.EnvironmentId.String()),
			)
			return
		}

		gatewayInstance = response.GetActualInstance()
	} else {
		resp.Diagnostics.AddError(
			"Missing parameter",
			"Cannot find the requested PingOne Gateway: gateway_id or name argument must be set.",
		)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(ctx, gatewayInstance)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (p *gatewayDataSourceModel) toState(ctx context.Context, apiObject interface{}) diag.Diagnostics {
	var diags, d diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	switch t := apiObject.(type) {
	case *management.Gateway:
		p.Id = framework.PingOneResourceIDToTF(t.GetId())
		p.GatewayId = framework.PingOneResourceIDToTF(t.GetId())
		p.EnvironmentId = framework.PingOneResourceIDToTF(*t.GetEnvironment().Id)
		p.Name = framework.StringOkToTF(t.GetNameOk())
		p.Description = framework.StringOkToTF(t.GetDescriptionOk())
		p.Type = framework.EnumOkToTF(t.GetTypeOk())
		p.Enabled = framework.BoolOkToTF(t.GetEnabledOk())

		// LDAP
		p.BindDN = types.StringNull()
		p.ConnectionSecurity = types.StringNull()
		p.FollowReferrals = types.BoolNull()
		p.Kerberos = types.ObjectNull(gatewayKerberosTFObjectTypes)
		p.Servers = types.SetNull(types.StringType)
		p.ValidateTLSCertificates = types.BoolNull()
		p.Vendor = types.StringNull()
		p.UserTypes = types.MapNull(types.ObjectType{AttrTypes: gatewayUserTypesTFObjectTypes})

		// Radius
		p.RadiusDavinciPolicyId = davincitypes.NewResourceIDNull()
		p.RadiusDefaultSharedSecret = types.StringNull()
		p.RadiusClients = types.SetNull(types.ObjectType{AttrTypes: gatewayRadiusClientsTFObjectTypes})
		p.RadiusNetworkPolicyServer = types.ObjectNull(gatewayRadiusNetworkPolicyServerTFObjectTypes)

	case *management.GatewayTypeLDAP:
		p.Id = framework.PingOneResourceIDToTF(t.GetId())
		p.GatewayId = framework.PingOneResourceIDToTF(t.GetId())
		p.EnvironmentId = framework.PingOneResourceIDToTF(*t.GetEnvironment().Id)
		p.Name = framework.StringOkToTF(t.GetNameOk())
		p.Description = framework.StringOkToTF(t.GetDescriptionOk())
		p.Type = framework.EnumOkToTF(t.GetTypeOk())
		p.Enabled = framework.BoolOkToTF(t.GetEnabledOk())

		// LDAP
		p.BindDN = framework.StringOkToTF(t.GetBindDNOk())
		p.ConnectionSecurity = framework.EnumOkToTF(t.GetConnectionSecurityOk())
		p.FollowReferrals = framework.BoolOkToTF(t.GetFollowReferralsOk())

		var kerberosPlan *gatewayKerberosResourceModelV1

		diags.Append(p.Kerberos.As(ctx, &kerberosPlan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return diags
		}

		kerberosObj, ok := t.GetKerberosOk()
		p.Kerberos, d = toStateKerberosOk(kerberosObj, ok, kerberosPlan)
		diags.Append(d...)

		p.Servers = framework.StringSetOkToTF(t.GetServersHostAndPortOk())

		p.ValidateTLSCertificates = framework.BoolOkToTF(t.GetValidateTlsCertificatesOk())
		p.Vendor = framework.EnumOkToTF(t.GetVendorOk())

		p.UserTypes, d = toStateUserTypesOk(t.GetUserTypesOk())
		diags.Append(d...)

		// Radius
		p.RadiusDavinciPolicyId = davincitypes.NewResourceIDNull()
		p.RadiusDefaultSharedSecret = types.StringNull()
		p.RadiusClients = types.SetNull(types.ObjectType{AttrTypes: gatewayRadiusClientsTFObjectTypes})
		p.RadiusNetworkPolicyServer = types.ObjectNull(gatewayRadiusNetworkPolicyServerTFObjectTypes)

	case *management.GatewayTypeRADIUS:
		p.Id = framework.PingOneResourceIDToTF(t.GetId())
		p.GatewayId = framework.PingOneResourceIDToTF(t.GetId())
		p.EnvironmentId = framework.PingOneResourceIDToTF(*t.GetEnvironment().Id)
		p.Name = framework.StringOkToTF(t.GetNameOk())
		p.Description = framework.StringOkToTF(t.GetDescriptionOk())
		p.Type = framework.EnumOkToTF(t.GetTypeOk())
		p.Enabled = framework.BoolOkToTF(t.GetEnabledOk())

		// LDAP
		p.BindDN = types.StringNull()
		p.ConnectionSecurity = types.StringNull()
		p.FollowReferrals = types.BoolNull()
		p.Kerberos = types.ObjectNull(gatewayKerberosTFObjectTypes)
		p.Servers = types.SetNull(types.StringType)
		p.ValidateTLSCertificates = types.BoolNull()
		p.Vendor = types.StringNull()
		p.UserTypes = types.MapNull(types.ObjectType{AttrTypes: gatewayUserTypesTFObjectTypes})

		// Radius
		if dv, ok := t.GetDavinciOk(); ok {
			if policy, ok := dv.GetPolicyOk(); ok {
				p.RadiusDavinciPolicyId = framework.DaVinciResourceIDOkToTF(policy.GetIdOk())
			}
		}
		p.RadiusDefaultSharedSecret = framework.StringOkToTF(t.GetDefaultSharedSecretOk())
		p.RadiusClients, d = toStateRadiusClientOk(t.GetRadiusClientsOk())
		diags.Append(d...)
		p.RadiusNetworkPolicyServer, d = toStateRadiusNetworkPolicyServerOk(t.GetNetworkPolicyServerOk())
		diags.Append(d...)

	default:
		diags.AddError(
			"Undefined gateway type",
			"Cannot identify the gateway type from the data object.  Please report this to the provider maintainers.",
		)

	}

	return diags
}
