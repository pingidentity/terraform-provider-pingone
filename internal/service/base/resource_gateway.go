package base

import (
	"context"
	"fmt"
	"hash/crc32"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
	"golang.org/x/exp/slices"
)

func ResourceGateway() *schema.Resource {

	ldapSchemaAttrList := []string{"bind_dn", "bind_password", "servers", "vendor"}
	radiusSchemaAttrList := []string{"radius_davinci_policy_id", "radius_client"}

	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage PingOne gateways.",

		CreateContext: resourceGatewayCreate,
		ReadContext:   resourceGatewayRead,
		UpdateContext: resourceGatewayUpdate,
		DeleteContext: resourceGatewayDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceGatewayImport,
		},

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Description:      "The ID of the environment to create the gateway in.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
				ForceNew:         true,
			},
			"name": {
				Description:      "The name of the gateway resource.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
			},
			"description": {
				Description: "A description to apply to the gateway resource.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"type": {
				Description:  fmt.Sprintf("The type of gateway resource. Options are `%s`, `%s`, `%s`, `%s` and `%s`.", string(management.ENUMGATEWAYTYPE_PING_FEDERATE), string(management.ENUMGATEWAYTYPE_API_GATEWAY_INTEGRATION), string(management.ENUMGATEWAYTYPE_LDAP), string(management.ENUMGATEWAYTYPE_RADIUS), string(management.ENUMGATEWAYTYPE_PING_INTELLIGENCE)),
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{string(management.ENUMGATEWAYTYPE_PING_FEDERATE), string(management.ENUMGATEWAYTYPE_API_GATEWAY_INTEGRATION), string(management.ENUMGATEWAYTYPE_LDAP), string(management.ENUMGATEWAYTYPE_RADIUS), string(management.ENUMGATEWAYTYPE_PING_INTELLIGENCE)}, false),
			},
			"enabled": {
				Description: "Indicates whether the gateway is enabled.",
				Type:        schema.TypeBool,
				Required:    true,
			},

			// LDAP
			"bind_dn": {
				Description:      "For LDAP gateways only: The distinguished name information to bind to the LDAP database (for example, `uid=pingone,dc=bxretail,dc=org`).",
				Type:             schema.TypeString,
				Optional:         true,
				RequiredWith:     ldapSchemaAttrList,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
			},
			"bind_password": {
				Description:      "For LDAP gateways only: The Bind password for the LDAP database.",
				Type:             schema.TypeString,
				Optional:         true,
				RequiredWith:     ldapSchemaAttrList,
				Sensitive:        true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
			},
			"connection_security": {
				Description:      fmt.Sprintf("For LDAP gateways only: The connection security type. Options are `%s`, `%s`, and `%s`.", string(management.ENUMGATEWAYTYPELDAPSECURITY_NONE), string(management.ENUMGATEWAYTYPELDAPSECURITY_TLS), string(management.ENUMGATEWAYTYPELDAPSECURITY_START_TLS)),
				Type:             schema.TypeString,
				Optional:         true,
				Default:          management.ENUMGATEWAYTYPELDAPSECURITY_NONE,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{string(management.ENUMGATEWAYTYPELDAPSECURITY_NONE), string(management.ENUMGATEWAYTYPELDAPSECURITY_TLS), string(management.ENUMGATEWAYTYPELDAPSECURITY_START_TLS)}, false)),
			},
			"kerberos_service_account_password": {
				Description: "For LDAP gateways only: The password for the Kerberos service account.",
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
			},
			"kerberos_service_account_upn": {
				Description: "For LDAP gateways only: The Kerberos service account user principal name (for example, `username@bxretail.org`).",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"kerberos_retain_previous_credentials_mins": {
				Description: "For LDAP gateways only: The number of minutes for which the previous credentials are persisted.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"servers": {
				Description:  "For LDAP gateways only: A list of LDAP server host name and port number combinations (for example, [`ds1.bxretail.org:636`, `ds2.bxretail.org:636`]).",
				Type:         schema.TypeSet,
				Optional:     true,
				RequiredWith: ldapSchemaAttrList,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
				},
			},
			"validate_tls_certificates": {
				Description: "For LDAP gateways only: Indicates whether or not to trust all SSL certificates (defaults to `true`). If this value is `false`, TLS certificates are not validated. When the value is set to `true`, only certificates that are signed by the default JVM CAs, or the CA certs that the customer has uploaded to the certificate service are trusted.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"vendor": {
				Description:      fmt.Sprintf("For LDAP gateways only: The LDAP vendor. Options are `%s`, `%s`, `%s`, `%s`, `%s`, `%s`, `%s`, and `%s`.", string(management.ENUMGATEWAYVENDOR_PING_DIRECTORY), string(management.ENUMGATEWAYVENDOR_MICROSOFT_ACTIVE_DIRECTORY), string(management.ENUMGATEWAYVENDOR_ORACLE_DIRECTORY_SERVER_ENTERPRISE_EDITION), string(management.ENUMGATEWAYVENDOR_ORACLE_UNIFIED_DIRECTORY), string(management.ENUMGATEWAYVENDOR_CA_DIRECTORY), string(management.ENUMGATEWAYVENDOR_OPEN_DJ_DIRECTORY), string(management.ENUMGATEWAYVENDOR_IBM__TIVOLI_SECURITY_DIRECTORY_SERVER), string(management.ENUMGATEWAYVENDOR_LDAP_V3_COMPLIANT_DIRECTORY_SERVER)),
				Type:             schema.TypeString,
				Optional:         true,
				RequiredWith:     ldapSchemaAttrList,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{string(management.ENUMGATEWAYVENDOR_PING_DIRECTORY), string(management.ENUMGATEWAYVENDOR_MICROSOFT_ACTIVE_DIRECTORY), string(management.ENUMGATEWAYVENDOR_ORACLE_DIRECTORY_SERVER_ENTERPRISE_EDITION), string(management.ENUMGATEWAYVENDOR_ORACLE_UNIFIED_DIRECTORY), string(management.ENUMGATEWAYVENDOR_CA_DIRECTORY), string(management.ENUMGATEWAYVENDOR_OPEN_DJ_DIRECTORY), string(management.ENUMGATEWAYVENDOR_IBM__TIVOLI_SECURITY_DIRECTORY_SERVER), string(management.ENUMGATEWAYVENDOR_LDAP_V3_COMPLIANT_DIRECTORY_SERVER)}, false)),
			},
			"user_type": {
				Description: "For LDAP gateways only: A collection of properties that define how users should be provisioned in PingOne. The `user_type` block specifies which user properties in PingOne correspond to the user properties in an external LDAP directory. You can use an LDAP browser to view the user properties in the external LDAP directory.",
				Type:        schema.TypeSet,
				Optional:    true,
				Set:         userItemsHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "Identifies the user type. This correlates to the `password.external.gateway.userType.id` User property.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"name": {
							Description:      "The name of the user type.",
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
						},
						"password_authority": {
							Description:      fmt.Sprintf("This can be either `%s` or `%s`. If set to `%s`, PingOne authenticates with the external directory initially, then PingOne authenticates all subsequent sign-ons.", string(management.ENUMGATEWAYPASSWORDAUTHORITY_PING_ONE), string(management.ENUMGATEWAYPASSWORDAUTHORITY_LDAP), string(management.ENUMGATEWAYPASSWORDAUTHORITY_PING_ONE)),
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{string(management.ENUMGATEWAYPASSWORDAUTHORITY_PING_ONE), string(management.ENUMGATEWAYPASSWORDAUTHORITY_LDAP)}, false)),
						},
						"search_base_dn": {
							Description: "The LDAP base domain name (DN) for this user type.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"user_link_attributes": {
							Description: "A list of strings that represent LDAP attribute names that uniquely identify the user, and link to users in PingOne.",
							Type:        schema.TypeList,
							Required:    true,
							Elem: &schema.Schema{
								Type:             schema.TypeString,
								ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
							},
						},
						"user_migration": {
							Description: "The configurations for initially authenticating new users who will be migrated to PingOne. Note: If there are multiple users having the same user name, only the first user processed is provisioned.",
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"lookup_filter_pattern": {
										Description: "The LDAP user search filter to use to match users against the entered user identifier at login. For example, `(((uid=${identifier})(mail=${identifier}))`. Alternatively, this can be a search against the user directory.",
										Type:        schema.TypeString,
										Required:    true,
									},
									"population_id": {
										Description:      "The ID of the population to use to create user entries during lookup.",
										Type:             schema.TypeString,
										Required:         true,
										ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
									},
									"attribute_mapping": {
										Description: "A collection of properties that define how users should be provisioned in PingOne. The `user_type` block specifies which user properties in PingOne correspond to the user properties in an external LDAP directory. You can use an LDAP browser to view the user properties in the external LDAP directory.",
										Type:        schema.TypeSet,
										Required:    true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Description:      "The name of a user attribute in PingOne. See [Users properties](https://apidocs.pingidentity.com/pingone/platform/v1/api/#users) for the complete list of available PingOne user attributes.",
													Type:             schema.TypeString,
													Required:         true,
													ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
												},
												"value": {
													Description:      "A reference to the corresponding external LDAP attribute.  Values are in the format `${ldapAttributes.mail}`, while Terraform HCL requires an additional `$` prefix character. For example, `$${ldapAttributes.mail}`",
													Type:             schema.TypeString,
													Required:         true,
													ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
												},
											},
										},
									},
								},
							},
						},
						"push_password_changes_to_ldap": {
							Description: "A boolean that determines whether password updates in PingOne should be pushed to the user's record in LDAP.  If false, the user cannot change the password and have it updated in the remote LDAP directory. In this case, operations for forgotten passwords or resetting of passwords are not available to a user referencing this gateway.",
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
						},
					},
				},
			},

			// RADIUS
			"radius_davinci_policy_id": {
				Description:      "For RADIUS gateways only: The ID of the DaVinci flow policy to use.",
				Type:             schema.TypeString,
				Optional:         true,
				RequiredWith:     radiusSchemaAttrList,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"radius_default_shared_secret": {
				Description:      "For RADIUS gateways only: Value to use for the shared secret if the shared secret is not provided for one or more of the RADIUS clients specified.",
				Type:             schema.TypeString,
				Optional:         true,
				Sensitive:        true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
			},
			"radius_client": {
				Description:  "For RADIUS gateways only: A collection of RADIUS clients.",
				Type:         schema.TypeSet,
				Optional:     true,
				RequiredWith: radiusSchemaAttrList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Description:      "The IP of the RADIUS client.",
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IsIPv4Address),
						},
						"shared_secret": {
							Description:      "The shared secret for the RADIUS client. If this value is not provided, the shared secret specified with `default_shared_secret` is used. If you are not providing a shared secret for the client, this parameter is optional.",
							Type:             schema.TypeString,
							Optional:         true,
							Sensitive:        true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
						},
					},
				},
			},
		},
	}
}

func resourceGatewayCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	gatewayRequest, diags := expandGatewayRequest(d)
	if diags.HasError() {
		return diags
	}

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.GatewaysApi.CreateGateway(ctx, d.Get("environment_id").(string)).CreateGatewayRequest(*gatewayRequest).Execute()
		},
		"CreateGateway",
		gatewayWriteErrors,
		sdk.DefaultCreateReadRetryable,
	)
	if diags.HasError() {
		return diags
	}

	respObject := resp.(*management.CreateGateway201Response)

	if gateway := respObject.Gateway; gateway != nil && gateway.GetId() != "" {
		d.SetId(gateway.GetId())
	} else if gateway := respObject.GatewayTypeLDAP; gateway != nil && gateway.GetId() != "" {
		d.SetId(gateway.GetId())
	} else if gateway := respObject.GatewayTypeRADIUS; gateway != nil && gateway.GetId() != "" {
		d.SetId(gateway.GetId())
	}

	return resourceGatewayRead(ctx, d, meta)
}

func resourceGatewayRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.GatewaysApi.ReadOneGateway(ctx, d.Get("environment_id").(string), d.Id()).Execute()
		},
		"ReadOneGateway",
		sdk.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
	)
	if diags.HasError() {
		return diags
	}

	if resp == nil {
		d.SetId("")
		return nil
	}

	respObject := resp.(*management.CreateGateway201Response)

	if gateway := respObject.Gateway; gateway != nil && gateway.GetId() != "" {

		d.Set("name", gateway.GetName())
		d.Set("enabled", gateway.GetEnabled())
		d.Set("type", gateway.GetType())

		if v, ok := gateway.GetDescriptionOk(); ok {
			d.Set("description", v)
		} else {
			d.Set("description", nil)
		}
	} else if gateway := respObject.GatewayTypeLDAP; gateway != nil && gateway.GetId() != "" {

		d.Set("name", gateway.GetName())
		d.Set("enabled", gateway.GetEnabled())
		d.Set("type", gateway.GetType())

		if v, ok := gateway.GetDescriptionOk(); ok {
			d.Set("description", v)
		} else {
			d.Set("description", nil)
		}

		d.Set("bind_dn", gateway.GetBindDN())

		d.Set("vendor", string(gateway.GetVendor()))

		d.Set("connection_security", string(gateway.GetConnectionSecurity()))

		if v, ok := gateway.GetKerberosOk(); ok {
			d.Set("kerberos_service_account_upn", v.GetServiceAccountUserPrincipalName())

			if v1, ok := v.GetMinutesToRetainPreviousCredentialsOk(); ok {
				d.Set("kerberos_retain_previous_credentials_mins", v1)
			} else {
				d.Set("kerberos_retain_previous_credentials_mins", nil)
			}

		} else {
			d.Set("kerberos_service_account_upn", nil)
			d.Set("kerberos_service_account_password", nil)
			d.Set("kerberos_retain_previous_credentials_mins", nil)
		}

		if v, ok := gateway.GetServersHostAndPortOk(); ok {
			d.Set("servers", v)
		} else {
			d.Set("servers", nil)
		}

		if v, ok := gateway.GetValidateTlsCertificatesOk(); ok {
			d.Set("validate_tls_certificates", v)
		} else {
			d.Set("validate_tls_certificates", nil)
		}

		d.Set("user_type", flattenUserType(gateway.GetUserTypes()))

	} else if gateway := respObject.GatewayTypeRADIUS; gateway != nil && gateway.GetId() != "" {

		d.Set("name", gateway.GetName())
		d.Set("enabled", gateway.GetEnabled())
		d.Set("type", gateway.GetType())

		if v, ok := gateway.GetDescriptionOk(); ok {
			d.Set("description", v)
		} else {
			d.Set("description", nil)
		}

		d.Set("radius_davinci_policy_id", gateway.GetDavinci().Policy.Id)

		if v, ok := gateway.GetDefaultSharedSecretOk(); ok {
			d.Set("radius_default_shared_secret", v)
		} else {
			d.Set("radius_default_shared_secret", nil)
		}

		radiusClientsFlattened := make([]interface{}, 0)
		for _, v := range gateway.GetRadiusClients() {

			radiusClient := map[string]interface{}{
				"ip": v.GetIp(),
			}

			if v1, ok := v.GetSharedSecretOk(); ok {
				radiusClient["shared_secret"] = *v1
			}

			radiusClientsFlattened = append(radiusClientsFlattened, radiusClient)
		}

		d.Set("radius_client", radiusClientsFlattened)
	}

	return diags
}

func resourceGatewayUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	gatewayRequest, diags := expandGatewayRequest(d)
	if diags.HasError() {
		return diags
	}

	_, diags = sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.GatewaysApi.UpdateGateway(ctx, d.Get("environment_id").(string), d.Id()).CreateGatewayRequest(*gatewayRequest).Execute()
		},
		"UpdateGateway",
		gatewayWriteErrors,
		sdk.DefaultRetryable,
	)
	if diags.HasError() {
		return diags
	}

	return resourceGatewayRead(ctx, d, meta)
}

func resourceGatewayDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	_, diags = sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			r, err := apiClient.GatewaysApi.DeleteGateway(ctx, d.Get("environment_id").(string), d.Id()).Execute()
			return nil, r, err
		},
		"DeleteGateway",
		sdk.CustomErrorResourceNotFoundWarning,
		sdk.DefaultRetryable,
	)
	if diags.HasError() {
		return diags
	}

	return diags
}

func resourceGatewayImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	splitLength := 2
	attributes := strings.SplitN(d.Id(), "/", splitLength)

	if len(attributes) != splitLength {
		return nil, fmt.Errorf("invalid id (\"%s\") specified, should be in format \"environmentID/gatewayID\"", d.Id())
	}

	environmentID, gatewayID := attributes[0], attributes[1]

	d.Set("environment_id", environmentID)
	d.SetId(gatewayID)

	resourceGatewayRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}

var (
	gatewayWriteErrors = func(error model.P1Error) diag.Diagnostics {
		var diags diag.Diagnostics

		// Invalid shared secret combination
		if details, ok := error.GetDetailsOk(); ok && details != nil && len(details) > 0 {
			if code, ok := details[0].GetCodeOk(); ok && *code == "INVALID_VALUE" {
				diags = diag.FromErr(fmt.Errorf(details[0].GetMessage()))

				return diags
			}
		}

		return nil
	}
)

func expandGatewayRequest(d *schema.ResourceData) (*management.CreateGatewayRequest, diag.Diagnostics) {
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

func userItemsHash(v interface{}) int {

	c := int(crc32.ChecksumIEEE([]byte(v.(map[string]interface{})["name"].(string))))
	if c >= 0 {
		return c
	}
	if -c >= 0 {
		return -c
	}
	// v == MinInt
	return 0
}

func flattenUserType(c []management.GatewayTypeLDAPAllOfUserTypes) *schema.Set {

	items := schema.NewSet(userItemsHash, nil)

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

		items.Add(item)

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
