package base

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func DatasourceLicense() *schema.Resource {
	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Datasource to read detailed PingOne license data, selected by the license ID.",

		ReadContext: datasourcePingOneLicenseRead,

		Schema: map[string]*schema.Schema{
			"organization_id": {
				Description:      "A string that specifies the organization resource’s unique identifier associated with the license.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"license_id": {
				Description:      "A string that specifies the license resource’s unique identifier.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"name": {
				Description: "A string that specifies a descriptive name for the license.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"package": {
				Description: "A string that specifies the license template on which this license is based. Options are `TRIAL`, `STANDARD`, `PREMIUM`, `MFA`, `RISK`, `MFARISK`, and `GLOBAL`.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"status": {
				Description: "A string that specifies the status of the license. Options are `ACTIVE`, `EXPIRED`, and `FUTURE`.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"replaces_license_id": {
				Description: "A string that specifies the license ID of the license that is replaced by this license.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"replaced_by_license_id": {
				Description: "A string that specifies the license ID of the license that replaces this license.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"begins_at": {
				Description: "The date and time this license begins.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"expires_at": {
				Description: "The date and time this license expires. `TRIAL` licenses stop access to PingOne services at expiration. All other licenses trigger an event to send a notification when the license expires but do not block services.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"terminates_at": {
				Description: "An attribute that designates the exact date and time when this license terminates access to PingOne services.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"assigned_environments_count": {
				Description: "An integer that specifies the total number of environments associated with this license.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"advanced_services": {
				Description: "A block that describes features related to **advanced services**.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"pingid": {
							Description: "A block that describes features related to **PingID** advanced service.",
							Type:        schema.TypeList,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"included": {
										Description: "A boolean that specifies whether the PingID advanced service is enabled in the organization.",
										Type:        schema.TypeBool,
										Computed:    true,
									},
									"type": {
										Description: "A string that specifies the type of PingID advanced service.",
										Type:        schema.TypeString,
										Computed:    true,
									},
								},
							},
						},
					},
				},
			},
			"authorize": {
				Description: "A block that describes features related to the **authorize** services.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allow_api_access_management": {
							Description: "A boolean that specifies whether to enable the PingOne Authorize API access management feature.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"allow_dynamic_authorization": {
							Description: "A boolean that specifies whether to enable the PingOne Authorize dynamic authorization feature.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
					},
				},
			},
			"credentials": {
				Description: "A block that describes features related to the **credentials** services.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allow_credentials": {
							Description: "A boolean that specifies whether to enable the PingOne Credentials feature.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
					},
				},
			},
			"environments": {
				Description: "A block that describes features related to the **environments** in the organization.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allow_add_resources": {
							Description: "A boolean that specifies whether the license supports creation of resources in the specified environment.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"allow_connections": {
							Description: "A boolean that specifies whether the license supports creation of application connections in the specified environment.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"allow_custom_domain": {
							Description: "A boolean that specifies whether the license supports creation of a custom domain in the specified environment.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"allow_custom_schema": {
							Description: "A boolean that specifies whether the license supports using custom schema attributes in the specified environment.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"allow_production": {
							Description: "A boolean that specifies whether production environments are allowed.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"max": {
							Description: "An integer that specifies the maximum number of environments allowed.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"regions": {
							Description: "A string that specifies the allowed regions associated with environments. Options are `NA`, `EU`, `CA` and `AP`.",
							Type:        schema.TypeSet,
							Computed:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"fraud": {
				Description: "A block that describes features related to the **fraud** services.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allow_bot_malicious_device_detection": {
							Description: "A boolean that specifies whether to enable the Malicious device detection features of PingOne Fraud.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"allow_account_protection": {
							Description: "A boolean that specifies whether to enable the account protection features of PingOne Fraud.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
					},
				},
			},
			"gateways": {
				Description: "A block that describes features related to the **gateway** services.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allow_ldap_gateway": {
							Description: "A boolean that specifies whether to enable the LDAP Gateway features of PingOne.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"allow_kerberos_gateway": {
							Description: "A boolean that specifies whether to enable the Kerberos Gateway features of PingOne.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"allow_radius_gateway": {
							Description: "A boolean that specifies whether to enable the RADIUS Gateway features of PingOne.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
					},
				},
			},
			"intelligence": {
				Description: "A block that describes features related to the **intelligence** services.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allow_advanced_predictors": {
							Description: "A boolean that specifies whether your license permits you to configure advanced risk features.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"allow_geo_velocity": {
							Description: "A boolean that specifies whether to use the intelligence geo-velocity feature. For `TRIAL` (unpaid) licenses, the default value is `true`. For `ADMIN`, `GLOBAL`, `RISK`, and `MFARISK`, the default value is `true`.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"allow_anonymous_network_detection": {
							Description: "A boolean that specifies whether to use the intelligence anonymous network detection feature. For `TRIAL` (unpaid) licenses, the default value is `true`. For `ADMIN`, `GLOBAL`, `RISK`, and `MFARISK`, the default value is `true`.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"allow_reputation": {
							Description: "A boolean that specifies whether to use the intelligence IP reputation feature. For `TRIAL` (unpaid) licenses, the default value is `true`. For `ADMIN`, `GLOBAL`, `RISK`, and `MFARISK`, the default value is `true`.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"allow_data_consent": {
							Description: "A boolean that specifies whether the customer has opted in to allow user and event behavior analytics (UEBA) data collection.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"allow_risk": {
							Description: "A boolean that specifies whether your license permits you to configure risk features such as sign-on policies that include rules to detect anomalous changes to your locations (such as impossible travel). This capability is supported for `TRIAL`, `RISK`, and `MFARISK` license packages. Note: The sharing of user data to enable our machine-learning engine, which is integral to PingOne Risk, is captured in the license property `intelligence.allow_data_consent`, but it is not set to `true` by default in any license package. This license capability always requires active consent by the customer before it can be enabled, and if consent is given, then it allows the full scope of intelligence features included in PingOne Risk (and PingOne Risk plus MFA).",
							Type:        schema.TypeBool,
							Computed:    true,
						},
					},
				},
			},
			"mfa": {
				Description: "A block that describes features related to the **mfa** service.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allow_push_notification": {
							Description: "A boolean that specifies whether push notifications are allowed. For `TRIAL` (unpaid) licenses, the default value is `true`. For other license package types, adoption of the feature determines the default value.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"allow_notification_outside_whitelist": {
							Description: "A boolean that specifies whether the license supports sending notifications outside of the environment's whitelist.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"allow_fido2_devices": {
							Description: "A boolean that specifies whether FIDO2 devices are allowed. For `TRIAL` (unpaid) licenses, the default value is `true`. For other license package types, adoption of the feature determines the default value.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"allow_voice_otp": {
							Description: "A boolean that specifies whether Voice OTP devices are allowed.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"allow_email_otp": {
							Description: "A boolean that specifies whether Email OTP devices are allowed.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"allow_sms_otp": {
							Description: "A boolean that specifies whether SMS OTP devices are allowed.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"allow_totp": {
							Description: "A boolean that specifies whether TOTP devices are allowed.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
					},
				},
			},
			"orchestrate": {
				Description: "A block that describes features related to the **identity orchestration** services.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allow_orchestration": {
							Description: "A boolean that specifies whether the core orchestration services are allowed.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
					},
				},
			},
			"users": {
				Description: "A block that describes features related to the **users** in the organization.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allow_password_management_notifications": {
							Description: "A boolean that specifies whether the license supports sending password management notifications.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"allow_identity_providers": {
							Description: "A boolean that specifies whether the license supports using external identity providers in the specified environment.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"allow_my_account": {
							Description: "A boolean that specifies whether the license supports using My Account capabilities in the specified environment.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"allow_password_only_authentication": {
							Description: "A boolean that specifies whether the license supports using password only login capabilities in the specified environment.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"allow_password_policy": {
							Description: "A boolean that specifies whether the license supports using password policies in the specified environment.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"allow_provisioning": {
							Description: "A boolean that specifies whether the license supports using provisioning capabilities in the specified environment.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"allow_inbound_provisioning": {
							Description: "A boolean that specifies whether the license supports using inbound provisioning capabilities in the specified environment.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"allow_role_assignment": {
							Description: "A boolean that specifies whether the license supports role assignments in the specified environment.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"allow_verification_flow": {
							Description: "A boolean that specifies whether the license supports using verification flows in the specified environment.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"allow_update_self": {
							Description: "A boolean that specifies whether the license supports allowing users to update their own profile.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"entitled_to_support": {
							Description: "A boolean that specifies whether the license allows PingOne support.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"max": {
							Description: "An integer that specifies the maximum number of users allowed per environment.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"max_hard_limit": {
							Description: "An integer that specifies the maximum number of users (hard limit) allowed per environment.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"annual_active_included": {
							Description: "An integer that specifies a soft limit on the number of active identities across all environments on the license per year. This property is not visible if a value is not provided at the time the license is created.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"monthly_active_included": {
							Description: "An integer that specifies a soft limit on the number of active identities across all environments on the license per month. This property is not visible if a value is not provided at the time the license is created.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
					},
				},
			},
			"verify": {
				Description: "A block that describes features related to the **verify** services.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allow_push_notifications": {
							Description: "A boolean that specifies whether to enable the PingOne Verify push notifications feature.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"allow_document_match": {
							Description: "A boolean that specifies whether to enable the PingOne Verify document matching feature.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"allow_face_match": {
							Description: "A boolean that specifies whether to enable the PingOne Verify face matching feature.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"allow_manual_id_inspection": {
							Description: "A boolean that specifies whether to enable the PingOne Verify manual ID inspection feature.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func datasourcePingOneLicenseRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	var resp management.License

	licenseResp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.LicensesApi.ReadOneLicense(ctx, d.Get("organization_id").(string), d.Get("license_id").(string)).Execute()
		},
		"ReadOneLicense",
		sdk.DefaultCustomError,
		sdk.DefaultRetryable,
	)
	if diags.HasError() {
		return diags
	}

	resp = *licenseResp.(*management.License)

	d.SetId(resp.GetId())
	d.Set("name", resp.GetName())
	d.Set("package", resp.GetPackage())
	d.Set("status", resp.GetStatus())
	d.Set("replaces_license_id", resp.GetReplacesLicense().Id)
	d.Set("replaced_by_license_id", resp.GetReplacedByLicense().Id)
	d.Set("begins_at", resp.GetBeginsAt())
	d.Set("expires_at", resp.GetExpiresAt())
	d.Set("terminates_at", resp.GetTerminatesAt())
	d.Set("assigned_environments_count", int(resp.GetAssignedEnvironmentsCount()))

	d.Set("advanced_services", flattenLicenseAdvancedServices(resp.GetAdvancedServices()))
	d.Set("authorize", flattenLicenseAuthorize(resp.GetAuthorize()))
	d.Set("credentials", flattenLicenseCredentials(resp.GetCredentials()))
	d.Set("environments", flattenLicenseEnvironments(resp.GetEnvironments()))
	d.Set("fraud", flattenLicenseFraud(resp.GetFraud()))
	d.Set("gateways", flattenLicenseGateways(resp.GetGateways()))
	d.Set("intelligence", flattenLicenseIntelligence(resp.GetIntelligence()))
	d.Set("mfa", flattenLicenseMFA(resp.GetMfa()))
	d.Set("orchestrate", flattenLicenseOrchestrate(resp.GetOrchestrate()))
	d.Set("users", flattenLicenseUsers(resp.GetUsers()))
	d.Set("verify", flattenLicenseVerify(resp.GetVerify()))

	return diags
}

func flattenLicenseAdvancedServices(v management.LicenseAdvancedServices) []map[string]interface{} {
	item := map[string]interface{}{
		"pingid": flattenLicenseAdvancedServicesPingID(v.GetPingId()),
	}

	return append(make([]map[string]interface{}, 0), item)
}

func flattenLicenseAdvancedServicesPingID(v management.LicenseAdvancedServicesPingId) []map[string]interface{} {
	item := map[string]interface{}{
		"included": v.GetIncluded(),
		"type":     v.GetType(),
	}

	return append(make([]map[string]interface{}, 0), item)
}

func flattenLicenseAuthorize(v management.LicenseAuthorize) []map[string]interface{} {
	item := map[string]interface{}{
		"allow_api_access_management": v.GetAllowApiAccessManagement(),
		"allow_dynamic_authorization": v.GetAllowDynamicAuthorization(),
	}

	return append(make([]map[string]interface{}, 0), item)
}

func flattenLicenseCredentials(v management.LicenseCredentials) []map[string]interface{} {
	item := map[string]interface{}{
		"allow_credentials": v.GetAllowCredentials(),
	}

	return append(make([]map[string]interface{}, 0), item)
}

func flattenLicenseEnvironments(v management.LicenseEnvironments) []map[string]interface{} {
	item := map[string]interface{}{
		"allow_add_resources": v.GetAllowConnections(),
		"allow_connections":   v.GetAllowConnections(),
		"allow_custom_domain": v.GetAllowCustomDomain(),
		"allow_custom_schema": v.GetAllowCustomSchema(),
		"allow_production":    v.GetAllowProduction(),
		"max":                 int(v.GetMax()),
		"regions":             v.GetRegions(),
	}

	return append(make([]map[string]interface{}, 0), item)
}

func flattenLicenseFraud(v management.LicenseFraud) []map[string]interface{} {
	item := map[string]interface{}{
		"allow_bot_malicious_device_detection": v.GetAllowBotMaliciousDeviceDetection(),
		"allow_account_protection":             v.GetAllowAccountProtection(),
	}

	return append(make([]map[string]interface{}, 0), item)
}

func flattenLicenseGateways(v management.LicenseGateways) []map[string]interface{} {
	item := map[string]interface{}{
		"allow_ldap_gateway":     v.GetAllowLdapGateway(),
		"allow_kerberos_gateway": v.GetAllowKerberosGateway(),
		"allow_radius_gateway":   v.GetAllowRadiusGateway(),
	}

	return append(make([]map[string]interface{}, 0), item)
}

func flattenLicenseIntelligence(v management.LicenseIntelligence) []map[string]interface{} {
	item := map[string]interface{}{
		"allow_advanced_predictors":         v.GetAllowAdvancedPredictors(),
		"allow_geo_velocity":                v.GetAllowGeoVelocity(),
		"allow_anonymous_network_detection": v.GetAllowAnonymousNetworkDetection(),
		"allow_reputation":                  v.GetAllowReputation(),
		"allow_data_consent":                v.GetAllowDataConsent(),
		"allow_risk":                        v.GetAllowRisk(),
	}

	return append(make([]map[string]interface{}, 0), item)
}

func flattenLicenseMFA(v management.LicenseMfa) []map[string]interface{} {
	item := map[string]interface{}{
		"allow_push_notification":              v.GetAllowPushNotification(),
		"allow_notification_outside_whitelist": v.GetAllowNotificationOutsideWhitelist(),
		"allow_fido2_devices":                  v.GetAllowFido2Devices(),
		"allow_voice_otp":                      v.GetAllowVoiceOtp(),
		"allow_email_otp":                      v.GetAllowEmailOtp(),
		"allow_sms_otp":                        v.GetAllowSmsOtp(),
		"allow_totp":                           v.GetAllowTotp(),
	}

	return append(make([]map[string]interface{}, 0), item)
}

func flattenLicenseOrchestrate(v management.LicenseOrchestrate) []map[string]interface{} {
	item := map[string]interface{}{
		"allow_orchestration": v.GetAllowOrchestration(),
	}

	return append(make([]map[string]interface{}, 0), item)
}

func flattenLicenseUsers(v management.LicenseUsers) []map[string]interface{} {
	item := map[string]interface{}{
		"allow_password_management_notifications": v.GetAllowPasswordManagementNotifications(),
		"allow_identity_providers":                v.GetAllowIdentityProviders(),
		"allow_my_account":                        v.GetAllowMyAccount(),
		"allow_password_only_authentication":      v.GetAllowPasswordOnlyAuthentication(),
		"allow_password_policy":                   v.GetAllowPasswordPolicy(),
		"allow_provisioning":                      v.GetAllowProvisioning(),
		"allow_inbound_provisioning":              v.GetAllowInboundProvisioning(),
		"allow_role_assignment":                   v.GetAllowRoleAssignment(),
		"allow_verification_flow":                 v.GetAllowVerificationFlow(),
		"allow_update_self":                       v.GetAllowUpdateSelf(),
		"entitled_to_support":                     v.GetEntitledToSupport(),
		"max":                                     int(v.GetMax()),
		"max_hard_limit":                          int(v.GetHardLimitMax()),
		"annual_active_included":                  int(v.GetAnnualActiveIncluded()),
		"monthly_active_included":                 int(v.GetMonthlyActiveIncluded()),
	}

	return append(make([]map[string]interface{}, 0), item)
}

func flattenLicenseVerify(v management.LicenseVerify) []map[string]interface{} {
	item := map[string]interface{}{
		"allow_push_notifications":   v.GetAllowPushNotifications(),
		"allow_document_match":       v.GetAllowDocumentMatch(),
		"allow_face_match":           v.GetAllowFaceMatch(),
		"allow_manual_id_inspection": v.GetAllowManualIdInspection(),
	}

	return append(make([]map[string]interface{}, 0), item)
}
