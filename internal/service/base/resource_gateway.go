package base

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func ResourceGateway() *schema.Resource {
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
			// "type": {
			// 	Description:  fmt.Sprintf("The type of gateway resource. Options are `%s`, `%s` and `%s`.", string(management.ENUMGATEWAYTYPE_PING_FEDERATE), string(management.ENUMGATEWAYTYPE_API_GATEWAY_INTEGRATION), string(management.ENUMGATEWAYTYPE_LDAP)),
			// 	Type:         schema.TypeString,
			// 	Required:     true,
			// 	ForceNew:     true,
			// 	ValidateFunc: validation.StringInSlice([]string{string(management.ENUMGATEWAYTYPE_PING_FEDERATE), string(management.ENUMGATEWAYTYPE_API_GATEWAY_INTEGRATION), string(management.ENUMGATEWAYTYPE_LDAP)}, false),
			// },
			"enabled": {
				Description: "Indicates whether the gateway is enabled.",
				Type:        schema.TypeBool,
				Required:    true,
			},
			"pingfederate": {
				Description:  fmt.Sprintf("Sets the **%s** gateway type.", string(management.ENUMGATEWAYTYPE_PING_FEDERATE)),
				Type:         schema.TypeList,
				MaxItems:     1,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"pingfederate", "api_gateway", "ldap"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{},
				},
			},
			"api_gateway": {
				Description:  fmt.Sprintf("Sets the **%s** gateway type.", string(management.ENUMGATEWAYTYPE_API_GATEWAY_INTEGRATION)),
				Type:         schema.TypeList,
				MaxItems:     1,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"pingfederate", "api_gateway", "ldap"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{},
				},
			},
			"ldap": {
				Description:  fmt.Sprintf("Sets the **%s** gateway type.", string(management.ENUMGATEWAYTYPE_LDAP)),
				Type:         schema.TypeList,
				MaxItems:     1,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"pingfederate", "api_gateway", "ldap"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bind_dn": {
							Description:      "The distinguished name information to bind to the LDAP database (for example, `uid=pingone,dc=bxretail,dc=org`).",
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
						},
						"bind_password": {
							Description:      "The Bind password for the LDAP database.",
							Type:             schema.TypeString,
							Required:         true,
							Sensitive:        true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
						},
						"connection_security": {
							Description:      fmt.Sprintf("The connection security type. Options are `%s`, `%s`, and `%s`.", string(management.ENUMGATEWAYLDAPSECURITY_NONE), string(management.ENUMGATEWAYLDAPSECURITY_TLS), string(management.ENUMGATEWAYLDAPSECURITY_START_TLS)),
							Type:             schema.TypeString,
							Optional:         true,
							Default:          management.ENUMGATEWAYLDAPSECURITY_NONE,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{string(management.ENUMGATEWAYLDAPSECURITY_NONE), string(management.ENUMGATEWAYLDAPSECURITY_TLS), string(management.ENUMGATEWAYLDAPSECURITY_START_TLS)}, false)),
						},
						"kerberos_service_account_password": {
							Description: "The password for the Kerberos service account.",
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
						},
						"kerberos_service_account_upn": {
							Description:  "The Kerberos service account user principal name (for example, `username@bxretail.org`).",
							Type:         schema.TypeString,
							Optional:     true,
							RequiredWith: []string{"ldap.0.kerberos_service_account_password", "ldap.0.kerberos_retain_previous_credentials_mins"},
						},
						"kerberos_retain_previous_credentials_mins": {
							Description: "The number of minutes for which the previous credentials are persisted.",
							Type:        schema.TypeInt,
							Optional:    true,
						},
						"servers": {
							Description: "A list of LDAP server host name and port number combinations (for example, [`ds1.bxretail.org:636`, `ds2.bxretail.org:636`]).",
							Type:        schema.TypeSet,
							Optional:    true,
							Elem: &schema.Schema{
								Type:             schema.TypeString,
								ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
							},
						},
						"validate_tls_certificates": {
							Description: "Indicates whether or not to trust all SSL certificates (defaults to `true`). If this value is `false`, TLS certificates are not validated. When the value is set to `true`, only certificates that are signed by the default JVM CAs, or the CA certs that the customer has uploaded to the certificate service are trusted.",
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
						},
						"vendor": {
							Description:      fmt.Sprintf("The LDAP vendor. Options are `%s`, `%s`, `%s`, `%s`, `%s`, `%s`, `%s`, and `%s`.", string(management.ENUMGATEWAYVENDOR_PING_DIRECTORY), string(management.ENUMGATEWAYVENDOR_MICROSOFT_ACTIVE_DIRECTORY), string(management.ENUMGATEWAYVENDOR_ORACLE_DIRECTORY_SERVER_ENTERPRISE_EDITION), string(management.ENUMGATEWAYVENDOR_ORACLE_UNIFIED_DIRECTORY), string(management.ENUMGATEWAYVENDOR_CA_DIRECTORY), string(management.ENUMGATEWAYVENDOR_OPEN_DJ_DIRECTORY), string(management.ENUMGATEWAYVENDOR_IBM__TIVOLI_SECURITY_DIRECTORY_SERVER), string(management.ENUMGATEWAYVENDOR_LDAP_V3_COMPLIANT_DIRECTORY_SERVER)),
							Type:             schema.TypeString,
							Required:         true,
							ForceNew:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{string(management.ENUMGATEWAYVENDOR_PING_DIRECTORY), string(management.ENUMGATEWAYVENDOR_MICROSOFT_ACTIVE_DIRECTORY), string(management.ENUMGATEWAYVENDOR_ORACLE_DIRECTORY_SERVER_ENTERPRISE_EDITION), string(management.ENUMGATEWAYVENDOR_ORACLE_UNIFIED_DIRECTORY), string(management.ENUMGATEWAYVENDOR_CA_DIRECTORY), string(management.ENUMGATEWAYVENDOR_OPEN_DJ_DIRECTORY), string(management.ENUMGATEWAYVENDOR_IBM__TIVOLI_SECURITY_DIRECTORY_SERVER), string(management.ENUMGATEWAYVENDOR_LDAP_V3_COMPLIANT_DIRECTORY_SERVER)}, false)),
						},
						"user_type": {
							Description: "A collection of properties that define how users should be provisioned in PingOne. The `user_type` block specifies which user properties in PingOne correspond to the user properties in an external LDAP directory. You can use an LDAP browser to view the user properties in the external LDAP directory.",
							Type:        schema.TypeSet,
							Required:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Description:      "Identifies the user type. This correlates to the `password.external.gateway.userType.id` User property.",
										Type:             schema.TypeString,
										Required:         true,
										ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
										ForceNew:         true,
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
										Optional:    true,
									},
									"user_link_attributes": {
										Description: "A list of strings that represent LDAP attribute names that uniquely identify the user, and link to users in PingOne.",
										Type:        schema.TypeList,
										Optional:    true,
										Elem: &schema.Schema{
											Type:             schema.TypeString,
											ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
										},
									},
									"user_migration_lookup_filter_pattern": {
										Description: "The LDAP user search filter to use to match users against the entered user identifier at login. For example, `(((uid=${identifier})(mail=${identifier}))`. Alternatively, this can be a search against the user directory.",
										Type:        schema.TypeString,
										Optional:    true,
									},
									"user_migration_population_id": {
										Description:      "The ID of the population to use to create user entries during lookup.",
										Type:             schema.TypeString,
										Optional:         true,
										ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
									},
									"user_migration_attribute_mapping": {
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
									"push_password_changes_to_ldap": {
										Description: "A boolean that determines whether password updates in PingOne should be pushed to the user's record in LDAP.  If false, the user cannot change the password and have it updated in the remote LDAP directory. In this case, operations for forgotten passwords or resetting of passwords are not available to a user referencing this gateway.",
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
									},
								},
							},
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
		sdk.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	if diags.HasError() {
		return diags
	}

	respObject := resp.(*management.CreateGateway201Response)

	if gateway := respObject.Gateway; gateway != nil && gateway.GetId() != "" {
		d.SetId(gateway.GetId())
	} else if gateway := respObject.GatewayLDAP; gateway != nil && gateway.GetId() != "" {
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

		if v, ok := gateway.GetDescriptionOk(); ok {
			d.Set("description", v)
		} else {
			d.Set("description", nil)
		}

		if gateway.GetType() == management.ENUMGATEWAYTYPE_PING_FEDERATE {
			d.Set("pingfederate", make([]map[string]interface{}, 1))
		} else if gateway.GetType() == management.ENUMGATEWAYTYPE_API_GATEWAY_INTEGRATION {
			d.Set("api_gateway", make([]map[string]interface{}, 1))
		} else {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Unknown gateway type: %s", gateway.GetType()),
				Detail:   fmt.Sprintf("Only %s, %s and %s gateway types are supported in this provider.  Please raise an issue with the provider maintainers.", string(management.ENUMGATEWAYTYPE_PING_FEDERATE), string(management.ENUMGATEWAYTYPE_API_GATEWAY_INTEGRATION), string(management.ENUMGATEWAYTYPE_LDAP)),
			})

			return diags
		}
	} else if gateway := respObject.GatewayLDAP; gateway != nil && gateway.GetId() != "" {

		d.Set("name", gateway.GetName())
		d.Set("enabled", gateway.GetEnabled())

		if v, ok := gateway.GetDescriptionOk(); ok {
			d.Set("description", v)
		} else {
			d.Set("description", nil)
		}

		d.Set("ldap", flattenLDAPOptions(gateway))

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
		sdk.DefaultCustomError,
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

func expandGatewayRequest(d *schema.ResourceData) (*management.CreateGatewayRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	gatewayRequest := &management.CreateGatewayRequest{}

	if _, ok := d.GetOk("ldap"); ok {

		userTypes := expandLDAPUserTypes(d.Get("ldap.0.user_type").(*schema.Set))

		gateway := *management.NewGatewayLDAP(d.Get("name").(string), management.ENUMGATEWAYTYPE_LDAP, d.Get("enabled").(bool), d.Get("ldap.0.bind_dn").(string), d.Get("ldap.0.bind_password").(string), userTypes, management.EnumGatewayVendor(d.Get("ldap.0.vendor").(string)))

		if v, ok := d.GetOk("ldap.0.connection_security"); ok {
			gateway.SetConnectionSecurity(management.EnumGatewayLDAPSecurity(v.(string)))
		}

		if v, ok := d.GetOk("ldap.0.kerberos_service_account_upn"); ok {
			kerberos := management.NewGatewayLDAPAllOfKerberos(v.(string))

			if v1, ok := d.GetOk("ldap.0.kerberos_service_account_password"); ok {
				kerberos.SetServiceAccountPassword(v1.(string))
			}

			if v1, ok := d.GetOk("ldap.0.kerberos_retain_previous_credentials_mins"); ok {
				kerberos.SetMinutesToRetainPreviousCredentials(int32(v1.(int)))
			}

			gateway.SetKerberos(*kerberos)
		}

		if v, ok := d.GetOk("ldap.0.servers"); ok {
			obj := make([]string, 0)
			for _, str := range v.(*schema.Set).List() {
				obj = append(obj, str.(string))
			}
			gateway.SetServersHostAndPort(obj)
		}

		if v, ok := d.GetOk("ldap.0.validate_tls_certificates"); ok {
			gateway.SetValidateTlsCertificates(v.(bool))
		} else {
			gateway.SetValidateTlsCertificates(false)
		}

		gatewayRequest.GatewayLDAP = &gateway

	} else {

		var gatewayType management.EnumGatewayType
		if _, ok := d.GetOk("pingfederate"); ok {
			gatewayType = management.ENUMGATEWAYTYPE_PING_FEDERATE
		}

		if _, ok := d.GetOk("api_gateway"); ok {
			gatewayType = management.ENUMGATEWAYTYPE_API_GATEWAY_INTEGRATION
		}

		if string(gatewayType) == "" {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Cannot determine the gateway type",
				Detail:   "Ensure that either `pingfederate`, `api_gateway` or `ldap` block is set.",
			})

			return nil, diags
		}

		gateway := *management.NewGateway(d.Get("name").(string), gatewayType, d.Get("enabled").(bool)) // Gateway |  (optional)

		if v, ok := d.GetOk("description"); ok {
			gateway.SetDescription(v.(string))
		}

		gatewayRequest.Gateway = &gateway
	}

	return gatewayRequest, diags
}

func expandLDAPUserTypes(c *schema.Set) []management.GatewayLDAPAllOfUserTypes {

	userTypes := make([]management.GatewayLDAPAllOfUserTypes, 0)

	for _, v := range c.List() {
		obj := v.(map[string]interface{})

		newUserLookup := expandLDAPUserLookup(obj)

		userType := *management.NewGatewayLDAPAllOfUserTypes(
			obj["id"].(string),
			obj["name"].(string),
			newUserLookup,
			management.EnumGatewayPasswordAuthority(obj["password_authority"].(string)),
		)

		if v, ok := obj["push_password_changes_to_ldap"].(bool); ok {
			userType.SetAllowPasswordChanges(v)
		}

		if v, ok := obj["search_base_dn"].(string); ok && v != "" {
			userType.SetSearchBaseDn(v)
		}

		if v, ok := obj["user_link_attributes"].([]interface{}); ok && len(v) > 0 && v[0] != "" {
			obj := make([]string, 0)
			for _, str := range v {
				obj = append(obj, str.(string))
			}
			userType.SetOrderedCorrelationAttributes(obj)
		}

		userTypes = append(userTypes, userType)
	}

	return userTypes

}

func expandLDAPUserLookup(c map[string]interface{}) management.GatewayLDAPAllOfNewUserLookup {

	attributeMappings := expandLDAPUserLookupAttributeMappings(c["user_migration_attribute_mapping"].(*schema.Set).List())

	userLookup := *management.NewGatewayLDAPAllOfNewUserLookup(attributeMappings)

	if v, ok := c["user_migration_lookup_filter_pattern"].(string); ok && v != "" {
		userLookup.SetLdapFilterPattern(v)
	}

	if v, ok := c["user_migration_population_id"].(string); ok && v != "" {
		userLookup.SetPopulation(*management.NewGatewayLDAPAllOfNewUserLookupPopulation(v))
	}

	return userLookup

}

func expandLDAPUserLookupAttributeMappings(c []interface{}) []management.GatewayLDAPAllOfNewUserLookupAttributeMappings {
	mappings := make([]management.GatewayLDAPAllOfNewUserLookupAttributeMappings, 0)

	for _, v := range c {

		obj := v.(map[string]interface{})

		mappings = append(mappings, *management.NewGatewayLDAPAllOfNewUserLookupAttributeMappings(obj["name"].(string), obj["value"].(string)))

	}

	return mappings
}

func flattenLDAPOptions(gateway *management.GatewayLDAP) interface{} {

	// Required
	item := map[string]interface{}{
		"bind_dn":       gateway.GetBindDN(),
		"bind_password": gateway.GetBindPassword(),
		"vendor":        string(gateway.GetVendor()),
		"user_type":     flattenUserType(gateway.GetUserTypes()),
	}

	// Optional
	if v, ok := gateway.GetConnectionSecurityOk(); ok {
		item["connection_security"] = string(*v)
	} else {
		item["connection_security"] = nil
	}

	if v, ok := gateway.GetServersHostAndPortOk(); ok {
		item["servers"] = v
	} else {
		item["servers"] = nil
	}

	if v, ok := gateway.GetValidateTlsCertificatesOk(); ok {
		item["validate_tls_certificates"] = v
	} else {
		item["validate_tls_certificates"] = nil
	}

	if v, ok := gateway.GetKerberosOk(); ok {
		item["kerberos_service_account_upn"] = v.GetServiceAccountUserPrincipalName()

		if v1, ok := v.GetServiceAccountPasswordOk(); ok {
			item["kerberos_service_account_password"] = v1
		} else {
			item["kerberos_service_account_password"] = nil
		}

		if v1, ok := v.GetMinutesToRetainPreviousCredentialsOk(); ok {
			item["kerberos_retain_previous_credentials_mins"] = v1
		} else {
			item["kerberos_retain_previous_credentials_mins"] = nil
		}

	} else {
		item["kerberos_service_account_upn"] = nil
		item["kerberos_service_account_password"] = nil
		item["kerberos_retain_previous_credentials_mins"] = nil
	}

	items := make([]interface{}, 0)
	return append(items, item)
}

func flattenUserType(c []management.GatewayLDAPAllOfUserTypes) interface{} {

	items := make([]interface{}, 0)

	for _, v := range c {
		// Required
		item := map[string]interface{}{
			"id":                               v.GetId(),
			"name":                             v.GetName(),
			"password_authority":               string(v.GetPasswordAuthority()),
			"user_migration_attribute_mapping": flattenLDAPUserLookupAttributeMappings(v.GetNewUserLookup().AttributeMappings),
		}

		// Optional
		if v1, ok := v.GetSearchBaseDnOk(); ok {
			item["search_base_dn"] = v1
		} else {
			item["search_base_dn"] = nil
		}

		if v1, ok := v.GetOrderedCorrelationAttributesOk(); ok {
			item["user_link_attributes"] = v1
		} else {
			item["user_link_attributes"] = nil
		}

		if v1, ok := v.GetAllowPasswordChangesOk(); ok {
			item["push_password_changes_to_ldap"] = v1
		} else {
			item["push_password_changes_to_ldap"] = nil
		}

		if v1, ok := v.GetNewUserLookupOk(); ok {

			if v2, ok := v1.GetLdapFilterPatternOk(); ok {
				item["user_migration_lookup_filter_pattern"] = v2
			} else {
				item["user_migration_lookup_filter_pattern"] = nil
			}

			if v2, ok := v1.GetPopulationOk(); ok {
				item["user_migration_population_id"] = v2.GetId()
			} else {
				item["user_migration_population_id"] = nil
			}

		} else {
			item["user_migration_lookup_filter_pattern"] = nil
			item["user_migration_population_id"] = nil
		}

		items = append(items, item)

	}

	return items
}

func flattenLDAPUserLookupAttributeMappings(c []management.GatewayLDAPAllOfNewUserLookupAttributeMappings) interface{} {
	items := make([]interface{}, 0)

	for _, v := range c {
		items = append(items, map[string]interface{}{
			"name":  v.GetName(),
			"value": v.GetValue(),
		})
	}

	return items
}
