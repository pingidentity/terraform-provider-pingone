package base

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
)

// Types
type LicenseDataSource serviceClientType

type licenseDataSourceModel struct {
	Id                        pingonetypes.ResourceIDValue `tfsdk:"id"`
	OrganizationId            pingonetypes.ResourceIDValue `tfsdk:"organization_id"`
	LicenseId                 pingonetypes.ResourceIDValue `tfsdk:"license_id"`
	Name                      types.String                 `tfsdk:"name"`
	Package                   types.String                 `tfsdk:"package"`
	Status                    types.String                 `tfsdk:"status"`
	ReplacesLicenseId         pingonetypes.ResourceIDValue `tfsdk:"replaces_license_id"`
	ReplacedByLicenseId       pingonetypes.ResourceIDValue `tfsdk:"replaced_by_license_id"`
	BeginsAt                  timetypes.RFC3339            `tfsdk:"begins_at"`
	ExpiresAt                 timetypes.RFC3339            `tfsdk:"expires_at"`
	TerminatesAt              timetypes.RFC3339            `tfsdk:"terminates_at"`
	AssignedEnvironmentsCount types.Int64                  `tfsdk:"assigned_environments_count"`
	AdvancedServices          types.Object                 `tfsdk:"advanced_services"`
	Authorize                 types.Object                 `tfsdk:"authorize"`
	Credentials               types.Object                 `tfsdk:"credentials"`
	Environments              types.Object                 `tfsdk:"environments"`
	Fraud                     types.Object                 `tfsdk:"fraud"`
	Gateways                  types.Object                 `tfsdk:"gateways"`
	Intelligence              types.Object                 `tfsdk:"intelligence"`
	Mfa                       types.Object                 `tfsdk:"mfa"`
	Orchestrate               types.Object                 `tfsdk:"orchestrate"`
	Users                     types.Object                 `tfsdk:"users"`
	Verify                    types.Object                 `tfsdk:"verify"`
}

var (
	licenseAdvancedServicesTFObjectTypes = map[string]attr.Type{
		"pingid": types.ObjectType{AttrTypes: licenseAdvancedServicesPingidTFObjectTypes},
	}

	licenseAdvancedServicesPingidTFObjectTypes = map[string]attr.Type{
		"included": types.BoolType,
		"type":     types.StringType,
	}

	licenseAuthorizeTFObjectTypes = map[string]attr.Type{
		"allow_api_access_management": types.BoolType,
		"allow_dynamic_authorization": types.BoolType,
	}

	licenseCredentialsTFObjectTypes = map[string]attr.Type{
		"allow_credentials": types.BoolType,
	}

	licenseEnvironmentsTFObjectTypes = map[string]attr.Type{
		"allow_add_resources": types.BoolType,
		"allow_connections":   types.BoolType,
		"allow_custom_domain": types.BoolType,
		"allow_custom_schema": types.BoolType,
		"allow_production":    types.BoolType,
		"max":                 types.Int64Type,
		"regions":             types.SetType{ElemType: types.StringType},
	}

	licenseFraudTFObjectTypes = map[string]attr.Type{
		"allow_bot_malicious_device_detection": types.BoolType,
		"allow_account_protection":             types.BoolType,
	}

	licenseGatewaysTFObjectTypes = map[string]attr.Type{
		"allow_ldap_gateway":     types.BoolType,
		"allow_kerberos_gateway": types.BoolType,
		"allow_radius_gateway":   types.BoolType,
	}

	licenseIntelligenceTFObjectTypes = map[string]attr.Type{
		"allow_advanced_predictors":         types.BoolType,
		"allow_geo_velocity":                types.BoolType,
		"allow_anonymous_network_detection": types.BoolType,
		"allow_reputation":                  types.BoolType,
		"allow_data_consent":                types.BoolType,
		"allow_risk":                        types.BoolType,
	}

	licenseMfaTFObjectTypes = map[string]attr.Type{
		"allow_push_notification":              types.BoolType,
		"allow_notification_outside_whitelist": types.BoolType,
		"allow_fido2_devices":                  types.BoolType,
		"allow_voice_otp":                      types.BoolType,
		"allow_email_otp":                      types.BoolType,
		"allow_sms_otp":                        types.BoolType,
		"allow_totp":                           types.BoolType,
	}

	licenseOrchestrateTFObjectTypes = map[string]attr.Type{
		"allow_orchestration": types.BoolType,
	}

	licenseUsersTFObjectTypes = map[string]attr.Type{
		"allow_password_management_notifications": types.BoolType,
		"allow_identity_providers":                types.BoolType,
		"allow_my_account":                        types.BoolType,
		"allow_password_only_authentication":      types.BoolType,
		"allow_password_policy":                   types.BoolType,
		"allow_provisioning":                      types.BoolType,
		"allow_inbound_provisioning":              types.BoolType,
		"allow_role_assignment":                   types.BoolType,
		"allow_verification_flow":                 types.BoolType,
		"allow_update_self":                       types.BoolType,
		"entitled_to_support":                     types.BoolType,
		"max":                                     types.Int64Type,
		"max_hard_limit":                          types.Int64Type,
		"annual_active_included":                  types.Int64Type,
		"monthly_active_included":                 types.Int64Type,
	}

	licenseVerifyTFObjectTypes = map[string]attr.Type{
		"allow_push_notifications":   types.BoolType,
		"allow_document_match":       types.BoolType,
		"allow_face_match":           types.BoolType,
		"allow_manual_id_inspection": types.BoolType,
	}
)

// Framework interfaces
var (
	_ datasource.DataSource = &LicenseDataSource{}
)

// New Object
func NewLicenseDataSource() datasource.DataSource {
	return &LicenseDataSource{}
}

// Metadata
func (r *LicenseDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_license"
}

func (r *LicenseDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	packageDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the license template on which this license is based. Options are `TRIAL`, `STANDARD`, `PREMIUM`, `MFA`, `RISK`, `MFARISK`, and `GLOBAL`.",
	)

	statusDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the status of the license.",
	).AllowedValuesEnum(management.AllowedEnumLicenseStatusEnumValues)

	expiresAtDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The RFC3339 date and time this license expires. `TRIAL` licenses stop access to PingOne services at expiration. All other licenses trigger an event to send a notification when the license expires but do not block services.",
	)

	advancedServicesDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that describes features related to **advanced services**.",
	)

	advancedServicesPingidDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that describes features related to **PingID** advanced service.",
	)

	authorizeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that describes features related to the **authorize** services.",
	)

	credentialsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that describes features related to the **credentials** services.",
	)

	environmentsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that describes features related to the **environments** in the organization.",
	)

	environmentsRegionsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the allowed regions associated with environments.",
	).AllowedValuesEnum(management.AllowedEnumRegionCodeEnumValues)

	fraudDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that describes features related to the **fraud** services.",
	)

	gatewaysDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that describes features related to the **gateway** services.",
	)

	intelligenceDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that describes features related to the **intelligence** services.",
	)

	intelligenceAllowGeoVelocityDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether to use the intelligence geo-velocity feature. For `TRIAL` (unpaid) licenses, the default value is `true`. For `ADMIN`, `GLOBAL`, `RISK`, and `MFARISK`, the default value is `true`.",
	)

	intelligenceAllowAnonymousNetworkDetectionDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether to use the intelligence anonymous network detection feature. For `TRIAL` (unpaid) licenses, the default value is `true`. For `ADMIN`, `GLOBAL`, `RISK`, and `MFARISK`, the default value is `true`.",
	)

	intelligenceAllowReputationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether to use the intelligence IP reputation feature. For `TRIAL` (unpaid) licenses, the default value is `true`. For `ADMIN`, `GLOBAL`, `RISK`, and `MFARISK`, the default value is `true`.",
	)

	intelligenceAllowRiskDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether your license permits you to configure risk features such as sign-on policies that include rules to detect anomalous changes to your locations (such as impossible travel). This capability is supported for `TRIAL`, `RISK`, and `MFARISK` license packages. Note: The sharing of user data to enable our machine-learning engine, which is integral to PingOne Risk, is captured in the license property `intelligence.allow_data_consent`, but it is not set to `true` by default in any license package. This license capability always requires active consent by the customer before it can be enabled, and if consent is given, then it allows the full scope of intelligence features included in PingOne Risk (and PingOne Risk plus MFA).",
	)

	mfaDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that describes features related to the **mfa** service.",
	)

	mfaAllowPushNotificationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether push notifications are allowed. For `TRIAL` (unpaid) licenses, the default value is `true`. For other license package types, adoption of the feature determines the default value.",
	)

	mfaAllowFido2DevicesNotificationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether FIDO2 devices are allowed. For `TRIAL` (unpaid) licenses, the default value is `true`. For other license package types, adoption of the feature determines the default value.",
	)

	orchestrateDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that describes features related to the **identity orchestration** services.",
	)

	usersDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that describes features related to the **users** in the organization.",
	)

	verifyDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that describes features related to the **verify** services.",
	)

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Data source to read detailed PingOne license data, selected by the license ID.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"organization_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the organization resource’s unique identifier associated with the license."),
			),

			"license_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the license resource’s unique identifier."),
			),

			"name": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies a descriptive name for the license.").Description,
				Computed:    true,
			},

			"package": schema.StringAttribute{
				Description:         packageDescription.Description,
				MarkdownDescription: packageDescription.MarkdownDescription,
				Computed:            true,
			},

			"status": schema.StringAttribute{
				Description:         statusDescription.Description,
				MarkdownDescription: statusDescription.MarkdownDescription,
				Computed:            true,
			},

			"replaces_license_id": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the license ID of the license that is replaced by this license.").Description,
				Computed:    true,

				CustomType: pingonetypes.ResourceIDType{},
			},

			"replaced_by_license_id": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the license ID of the license that replaces this license.").Description,
				Computed:    true,

				CustomType: pingonetypes.ResourceIDType{},
			},

			"begins_at": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The RFC3339 date and time this license begins.").Description,
				Computed:    true,

				CustomType: timetypes.RFC3339Type{},
			},

			"expires_at": schema.StringAttribute{
				Description:         expiresAtDescription.Description,
				MarkdownDescription: expiresAtDescription.MarkdownDescription,
				Computed:            true,

				CustomType: timetypes.RFC3339Type{},
			},

			"terminates_at": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("The RFC3339 date and time when this license terminates access to PingOne services.").Description,
				Computed:    true,

				CustomType: timetypes.RFC3339Type{},
			},

			"assigned_environments_count": schema.Int64Attribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that specifies the total number of environments associated with this license.").Description,
				Computed:    true,
			},

			"advanced_services": schema.SingleNestedAttribute{
				Description:         advancedServicesDescription.Description,
				MarkdownDescription: advancedServicesDescription.MarkdownDescription,
				Computed:            true,

				Attributes: map[string]schema.Attribute{
					"pingid": schema.SingleNestedAttribute{
						Description:         advancedServicesPingidDescription.Description,
						MarkdownDescription: advancedServicesPingidDescription.MarkdownDescription,
						Computed:            true,

						Attributes: map[string]schema.Attribute{
							"included": schema.BoolAttribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether the PingID advanced service is enabled in the organization.").Description,
								Computed:    true,
							},

							"type": schema.StringAttribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the type of PingID advanced service.").Description,
								Computed:    true,
							},
						},
					},
				},
			},

			"authorize": schema.SingleNestedAttribute{
				Description:         authorizeDescription.Description,
				MarkdownDescription: authorizeDescription.MarkdownDescription,
				Computed:            true,

				Attributes: map[string]schema.Attribute{
					"allow_api_access_management": schema.BoolAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether to enable the PingOne Authorize API access management feature.").Description,
						Computed:    true,
					},

					"allow_dynamic_authorization": schema.BoolAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether to enable the PingOne Authorize dynamic authorization feature.").Description,
						Computed:    true,
					},
				},
			},

			"credentials": schema.SingleNestedAttribute{
				Description:         credentialsDescription.Description,
				MarkdownDescription: credentialsDescription.MarkdownDescription,
				Computed:            true,

				Attributes: map[string]schema.Attribute{
					"allow_credentials": schema.BoolAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether to enable the PingOne Credentials feature.").Description,
						Computed:    true,
					},
				},
			},

			"environments": schema.SingleNestedAttribute{
				Description:         environmentsDescription.Description,
				MarkdownDescription: environmentsDescription.MarkdownDescription,
				Computed:            true,

				Attributes: map[string]schema.Attribute{
					"allow_add_resources": schema.BoolAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether the license supports creation of resources in the specified environment.").Description,
						Computed:    true,
					},

					"allow_connections": schema.BoolAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether the license supports creation of application connections in the specified environment.").Description,
						Computed:    true,
					},

					"allow_custom_domain": schema.BoolAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether the license supports creation of a custom domain in the specified environment.").Description,
						Computed:    true,
					},

					"allow_custom_schema": schema.BoolAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether the license supports using custom schema attributes in the specified environment.").Description,
						Computed:    true,
					},

					"allow_production": schema.BoolAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether production environments are allowed.").Description,
						Computed:    true,
					},

					"max": schema.Int64Attribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that specifies the maximum number of environments allowed.").Description,
						Computed:    true,
					},

					"regions": schema.SetAttribute{
						Description:         environmentsRegionsDescription.Description,
						MarkdownDescription: environmentsRegionsDescription.MarkdownDescription,
						Computed:            true,

						ElementType: types.StringType,
					},
				},
			},

			"fraud": schema.SingleNestedAttribute{
				Description:         fraudDescription.Description,
				MarkdownDescription: fraudDescription.MarkdownDescription,
				Computed:            true,

				Attributes: map[string]schema.Attribute{
					"allow_bot_malicious_device_detection": schema.BoolAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether to enable the Malicious device detection features of PingOne Fraud.").Description,
						Computed:    true,
					},

					"allow_account_protection": schema.BoolAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether to enable the account protection features of PingOne Fraud.").Description,
						Computed:    true,
					},
				},
			},

			"gateways": schema.SingleNestedAttribute{
				Description:         gatewaysDescription.Description,
				MarkdownDescription: gatewaysDescription.MarkdownDescription,
				Computed:            true,

				Attributes: map[string]schema.Attribute{
					"allow_ldap_gateway": schema.BoolAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether to enable the LDAP Gateway features of PingOne.").Description,
						Computed:    true,
					},

					"allow_kerberos_gateway": schema.BoolAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether to enable the Kerberos Gateway features of PingOne.").Description,
						Computed:    true,
					},

					"allow_radius_gateway": schema.BoolAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether to enable the RADIUS Gateway features of PingOne.").Description,
						Computed:    true,
					},
				},
			},

			"intelligence": schema.SingleNestedAttribute{
				Description:         intelligenceDescription.Description,
				MarkdownDescription: intelligenceDescription.MarkdownDescription,
				Computed:            true,

				Attributes: map[string]schema.Attribute{
					"allow_advanced_predictors": schema.BoolAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether your license permits you to configure advanced risk features.").Description,
						Computed:    true,
					},

					"allow_geo_velocity": schema.BoolAttribute{
						Description:         intelligenceAllowGeoVelocityDescription.Description,
						MarkdownDescription: intelligenceAllowGeoVelocityDescription.MarkdownDescription,
						Computed:            true,
					},

					"allow_anonymous_network_detection": schema.BoolAttribute{
						Description:         intelligenceAllowAnonymousNetworkDetectionDescription.Description,
						MarkdownDescription: intelligenceAllowAnonymousNetworkDetectionDescription.MarkdownDescription,
						Computed:            true,
					},

					"allow_reputation": schema.BoolAttribute{
						Description:         intelligenceAllowReputationDescription.Description,
						MarkdownDescription: intelligenceAllowReputationDescription.MarkdownDescription,
						Computed:            true,
					},

					"allow_data_consent": schema.BoolAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether the customer has opted in to allow user and event behavior analytics (UEBA) data collection.").Description,
						Computed:    true,
					},

					"allow_risk": schema.BoolAttribute{
						Description:         intelligenceAllowRiskDescription.Description,
						MarkdownDescription: intelligenceAllowRiskDescription.MarkdownDescription,
						Computed:            true,
					},
				},
			},

			"mfa": schema.SingleNestedAttribute{
				Description:         mfaDescription.Description,
				MarkdownDescription: mfaDescription.MarkdownDescription,
				Computed:            true,

				Attributes: map[string]schema.Attribute{
					"allow_push_notification": schema.BoolAttribute{
						Description:         mfaAllowPushNotificationDescription.Description,
						MarkdownDescription: mfaAllowPushNotificationDescription.MarkdownDescription,
						Computed:            true,
					},

					"allow_notification_outside_whitelist": schema.BoolAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether the license supports sending notifications outside of the environment's whitelist.").Description,
						Computed:    true,
					},

					"allow_fido2_devices": schema.BoolAttribute{
						Description:         mfaAllowFido2DevicesNotificationDescription.Description,
						MarkdownDescription: mfaAllowFido2DevicesNotificationDescription.MarkdownDescription,
						Computed:            true,
					},

					"allow_voice_otp": schema.BoolAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether Voice OTP devices are allowed.").Description,
						Computed:    true,
					},

					"allow_email_otp": schema.BoolAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether Email OTP devices are allowed.").Description,
						Computed:    true,
					},

					"allow_sms_otp": schema.BoolAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether SMS OTP devices are allowed.").Description,
						Computed:    true,
					},

					"allow_totp": schema.BoolAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether TOTP devices are allowed.").Description,
						Computed:    true,
					},
				},
			},

			"orchestrate": schema.SingleNestedAttribute{
				Description:         orchestrateDescription.Description,
				MarkdownDescription: orchestrateDescription.MarkdownDescription,
				Computed:            true,

				Attributes: map[string]schema.Attribute{
					"allow_orchestration": schema.BoolAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether the core orchestration services are allowed.").Description,
						Computed:    true,
					},
				},
			},

			"users": schema.SingleNestedAttribute{
				Description:         usersDescription.Description,
				MarkdownDescription: usersDescription.MarkdownDescription,
				Computed:            true,

				Attributes: map[string]schema.Attribute{
					"allow_password_management_notifications": schema.BoolAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether the license supports sending password management notifications.").Description,
						Computed:    true,
					},

					"allow_identity_providers": schema.BoolAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether the license supports using external identity providers in the specified environment.").Description,
						Computed:    true,
					},

					"allow_my_account": schema.BoolAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether the license supports using My Account capabilities in the specified environment.").Description,
						Computed:    true,
					},

					"allow_password_only_authentication": schema.BoolAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether the license supports using password only login capabilities in the specified environment.").Description,
						Computed:    true,
					},

					"allow_password_policy": schema.BoolAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether the license supports using password policies in the specified environment.").Description,
						Computed:    true,
					},

					"allow_provisioning": schema.BoolAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether the license supports using provisioning capabilities in the specified environment.").Description,
						Computed:    true,
					},

					"allow_inbound_provisioning": schema.BoolAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether the license supports using inbound provisioning capabilities in the specified environment.").Description,
						Computed:    true,
					},

					"allow_role_assignment": schema.BoolAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether the license supports role assignments in the specified environment.").Description,
						Computed:    true,
					},

					"allow_verification_flow": schema.BoolAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether the license supports using verification flows in the specified environment.").Description,
						Computed:    true,
					},

					"allow_update_self": schema.BoolAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether the license supports allowing users to update their own profile.").Description,
						Computed:    true,
					},

					"entitled_to_support": schema.BoolAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether the license allows PingOne support.").Description,
						Computed:    true,
					},

					"max": schema.Int64Attribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that specifies the maximum number of users allowed per environment.").Description,
						Computed:    true,
					},

					"max_hard_limit": schema.Int64Attribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that specifies the maximum number of users (hard limit) allowed per environment.").Description,
						Computed:    true,
					},

					"annual_active_included": schema.Int64Attribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that specifies a soft limit on the number of active identities across all environments on the license per year. This property is not visible if a value is not provided at the time the license is created.").Description,
						Computed:    true,
					},

					"monthly_active_included": schema.Int64Attribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("An integer that specifies a soft limit on the number of active identities across all environments on the license per month. This property is not visible if a value is not provided at the time the license is created.").Description,
						Computed:    true,
					},
				},
			},

			"verify": schema.SingleNestedAttribute{
				Description:         verifyDescription.Description,
				MarkdownDescription: verifyDescription.MarkdownDescription,
				Computed:            true,

				Attributes: map[string]schema.Attribute{
					"allow_push_notifications": schema.BoolAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether to enable the PingOne Verify push notifications feature.").Description,
						Computed:    true,
					},

					"allow_document_match": schema.BoolAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether to enable the PingOne Verify document matching feature.").Description,
						Computed:    true,
					},

					"allow_face_match": schema.BoolAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether to enable the PingOne Verify face matching feature.").Description,
						Computed:    true,
					},

					"allow_manual_id_inspection": schema.BoolAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether to enable the PingOne Verify manual ID inspection feature.").Description,
						Computed:    true,
					},
				},
			},
		},
	}
}

func (r *LicenseDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *LicenseDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *licenseDataSourceModel

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

	// Run the API call
	var response *management.License
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return r.Client.ManagementAPIClient.LicensesApi.ReadOneLicense(ctx, data.OrganizationId.ValueString(), data.LicenseId.ValueString()).Execute()
		},
		"ReadOneLicense",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		&response,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(response)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (p *licenseDataSourceModel) toState(apiObject *management.License) diag.Diagnostics {
	var diags, d diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.PingOneResourceIDToTF(apiObject.GetId())
	p.OrganizationId = framework.PingOneResourceIDToTF(*apiObject.GetOrganization().Id)
	p.LicenseId = framework.PingOneResourceIDToTF(apiObject.GetId())
	p.Name = framework.StringOkToTF(apiObject.GetNameOk())
	p.Package = framework.StringOkToTF(apiObject.GetPackageOk())
	p.Status = framework.EnumOkToTF(apiObject.GetStatusOk())

	if v, ok := apiObject.GetReplacesLicenseOk(); ok {
		p.ReplacesLicenseId = framework.PingOneResourceIDOkToTF(v.GetIdOk())
	} else {
		p.ReplacesLicenseId = pingonetypes.NewResourceIDNull()
	}

	if v, ok := apiObject.GetReplacedByLicenseOk(); ok {
		p.ReplacedByLicenseId = framework.PingOneResourceIDOkToTF(v.GetIdOk())
	} else {
		p.ReplacedByLicenseId = pingonetypes.NewResourceIDNull()
	}

	p.BeginsAt = framework.TimeOkToTF(apiObject.GetBeginsAtOk())
	p.ExpiresAt = framework.TimeOkToTF(apiObject.GetExpiresAtOk())
	p.TerminatesAt = framework.TimeOkToTF(apiObject.GetTerminatesAtOk())
	p.AssignedEnvironmentsCount = framework.Int32OkToInt64TF(apiObject.GetAssignedEnvironmentsCountOk())

	p.AdvancedServices, d = licenseAdvancedServicesOkToTF(apiObject.GetAdvancedServicesOk())
	diags.Append(d...)

	p.Authorize, d = licenseAuthorizeOkToTF(apiObject.GetAuthorizeOk())
	diags.Append(d...)

	p.Credentials, d = licenseCredentialsOkToTF(apiObject.GetCredentialsOk())
	diags.Append(d...)

	p.Environments, d = licenseEnvironmentsOkToTF(apiObject.GetEnvironmentsOk())
	diags.Append(d...)

	p.Fraud, d = licenseFraudOkToTF(apiObject.GetFraudOk())
	diags.Append(d...)

	p.Gateways, d = licenseGatewaysOkToTF(apiObject.GetGatewaysOk())
	diags.Append(d...)

	p.Intelligence, d = licenseIntelligenceOkToTF(apiObject.GetIntelligenceOk())
	diags.Append(d...)

	p.Mfa, d = licenseMfaOkToTF(apiObject.GetMfaOk())
	diags.Append(d...)

	p.Orchestrate, d = licenseOrchestrateOkToTF(apiObject.GetOrchestrateOk())
	diags.Append(d...)

	p.Users, d = licenseUsersOkToTF(apiObject.GetUsersOk())
	diags.Append(d...)

	p.Verify, d = licenseVerifyOkToTF(apiObject.GetVerifyOk())
	diags.Append(d...)

	return diags
}

func licenseAdvancedServicesOkToTF(apiObject *management.LicenseAdvancedServices, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(licenseAdvancedServicesTFObjectTypes), diags
	}

	pingid, d := licenseAdvancedServicesPingIdOkToTF(apiObject.GetPingIdOk())
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(licenseAdvancedServicesTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(licenseAdvancedServicesTFObjectTypes, map[string]attr.Value{
		"pingid": pingid,
	})
	diags.Append(d...)

	return objValue, diags
}

func licenseAdvancedServicesPingIdOkToTF(apiObject *management.LicenseAdvancedServicesPingId, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(licenseAdvancedServicesPingidTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(licenseAdvancedServicesPingidTFObjectTypes, map[string]attr.Value{
		"included": framework.BoolOkToTF(apiObject.GetIncludedOk()),
		"type":     framework.StringOkToTF(apiObject.GetTypeOk()),
	})
	diags.Append(d...)

	return objValue, diags
}

func licenseAuthorizeOkToTF(apiObject *management.LicenseAuthorize, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(licenseAuthorizeTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(licenseAuthorizeTFObjectTypes, map[string]attr.Value{
		"allow_api_access_management": framework.BoolOkToTF(apiObject.GetAllowApiAccessManagementOk()),
		"allow_dynamic_authorization": framework.BoolOkToTF(apiObject.GetAllowDynamicAuthorizationOk()),
	})
	diags.Append(d...)

	return objValue, diags
}

func licenseCredentialsOkToTF(apiObject *management.LicenseCredentials, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(licenseCredentialsTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(licenseCredentialsTFObjectTypes, map[string]attr.Value{
		"allow_credentials": framework.BoolOkToTF(apiObject.GetAllowCredentialsOk()),
	})
	diags.Append(d...)

	return objValue, diags
}

func licenseEnvironmentsOkToTF(apiObject *management.LicenseEnvironments, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(licenseEnvironmentsTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(licenseEnvironmentsTFObjectTypes, map[string]attr.Value{
		"allow_add_resources": framework.BoolOkToTF(apiObject.GetAllowAddResourcesOk()),
		"allow_connections":   framework.BoolOkToTF(apiObject.GetAllowConnectionsOk()),
		"allow_custom_domain": framework.BoolOkToTF(apiObject.GetAllowCustomDomainOk()),
		"allow_custom_schema": framework.BoolOkToTF(apiObject.GetAllowCustomSchemaOk()),
		"allow_production":    framework.BoolOkToTF(apiObject.GetAllowProductionOk()),
		"max":                 framework.Int32OkToInt64TF(apiObject.GetMaxOk()),
		"regions":             framework.EnumSetOkToTF(apiObject.GetRegionsOk()),
	})
	diags.Append(d...)

	return objValue, diags
}

func licenseFraudOkToTF(apiObject *management.LicenseFraud, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(licenseFraudTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(licenseFraudTFObjectTypes, map[string]attr.Value{
		"allow_bot_malicious_device_detection": framework.BoolOkToTF(apiObject.GetAllowBotMaliciousDeviceDetectionOk()),
		"allow_account_protection":             framework.BoolOkToTF(apiObject.GetAllowAccountProtectionOk()),
	})
	diags.Append(d...)

	return objValue, diags
}

func licenseGatewaysOkToTF(apiObject *management.LicenseGateways, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(licenseGatewaysTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(licenseGatewaysTFObjectTypes, map[string]attr.Value{
		"allow_ldap_gateway":     framework.BoolOkToTF(apiObject.GetAllowLdapGatewayOk()),
		"allow_kerberos_gateway": framework.BoolOkToTF(apiObject.GetAllowKerberosGatewayOk()),
		"allow_radius_gateway":   framework.BoolOkToTF(apiObject.GetAllowRadiusGatewayOk()),
	})
	diags.Append(d...)

	return objValue, diags
}

func licenseIntelligenceOkToTF(apiObject *management.LicenseIntelligence, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(licenseIntelligenceTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(licenseIntelligenceTFObjectTypes, map[string]attr.Value{
		"allow_advanced_predictors":         framework.BoolOkToTF(apiObject.GetAllowAdvancedPredictorsOk()),
		"allow_geo_velocity":                framework.BoolOkToTF(apiObject.GetAllowGeoVelocityOk()),
		"allow_anonymous_network_detection": framework.BoolOkToTF(apiObject.GetAllowAnonymousNetworkDetectionOk()),
		"allow_reputation":                  framework.BoolOkToTF(apiObject.GetAllowReputationOk()),
		"allow_data_consent":                framework.BoolOkToTF(apiObject.GetAllowDataConsentOk()),
		"allow_risk":                        framework.BoolOkToTF(apiObject.GetAllowRiskOk()),
	})
	diags.Append(d...)

	return objValue, diags
}

func licenseMfaOkToTF(apiObject *management.LicenseMfa, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(licenseMfaTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(licenseMfaTFObjectTypes, map[string]attr.Value{
		"allow_push_notification":              framework.BoolOkToTF(apiObject.GetAllowPushNotificationOk()),
		"allow_notification_outside_whitelist": framework.BoolOkToTF(apiObject.GetAllowNotificationOutsideWhitelistOk()),
		"allow_fido2_devices":                  framework.BoolOkToTF(apiObject.GetAllowFido2DevicesOk()),
		"allow_voice_otp":                      framework.BoolOkToTF(apiObject.GetAllowVoiceOtpOk()),
		"allow_email_otp":                      framework.BoolOkToTF(apiObject.GetAllowEmailOtpOk()),
		"allow_sms_otp":                        framework.BoolOkToTF(apiObject.GetAllowSmsOtpOk()),
		"allow_totp":                           framework.BoolOkToTF(apiObject.GetAllowTotpOk()),
	})
	diags.Append(d...)

	return objValue, diags
}

func licenseOrchestrateOkToTF(apiObject *management.LicenseOrchestrate, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(licenseOrchestrateTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(licenseOrchestrateTFObjectTypes, map[string]attr.Value{
		"allow_orchestration": framework.BoolOkToTF(apiObject.GetAllowOrchestrationOk()),
	})
	diags.Append(d...)

	return objValue, diags
}

func licenseUsersOkToTF(apiObject *management.LicenseUsers, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(licenseUsersTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(licenseUsersTFObjectTypes, map[string]attr.Value{
		"allow_password_management_notifications": framework.BoolOkToTF(apiObject.GetAllowPasswordManagementNotificationsOk()),
		"allow_identity_providers":                framework.BoolOkToTF(apiObject.GetAllowIdentityProvidersOk()),
		"allow_my_account":                        framework.BoolOkToTF(apiObject.GetAllowMyAccountOk()),
		"allow_password_only_authentication":      framework.BoolOkToTF(apiObject.GetAllowPasswordOnlyAuthenticationOk()),
		"allow_password_policy":                   framework.BoolOkToTF(apiObject.GetAllowPasswordPolicyOk()),
		"allow_provisioning":                      framework.BoolOkToTF(apiObject.GetAllowProvisioningOk()),
		"allow_inbound_provisioning":              framework.BoolOkToTF(apiObject.GetAllowInboundProvisioningOk()),
		"allow_role_assignment":                   framework.BoolOkToTF(apiObject.GetAllowRoleAssignmentOk()),
		"allow_verification_flow":                 framework.BoolOkToTF(apiObject.GetAllowVerificationFlowOk()),
		"allow_update_self":                       framework.BoolOkToTF(apiObject.GetAllowUpdateSelfOk()),
		"entitled_to_support":                     framework.BoolOkToTF(apiObject.GetEntitledToSupportOk()),
		"max":                                     framework.Int32OkToInt64TF(apiObject.GetMaxOk()),
		"max_hard_limit":                          framework.Int32OkToInt64TF(apiObject.GetHardLimitMaxOk()),
		"annual_active_included":                  framework.Int32OkToInt64TF(apiObject.GetAnnualActiveIncludedOk()),
		"monthly_active_included":                 framework.Int32OkToInt64TF(apiObject.GetMonthlyActiveIncludedOk()),
	})
	diags.Append(d...)

	return objValue, diags
}

func licenseVerifyOkToTF(apiObject *management.LicenseVerify, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(licenseVerifyTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(licenseVerifyTFObjectTypes, map[string]attr.Value{
		"allow_push_notifications":   framework.BoolOkToTF(apiObject.GetAllowPushNotificationsOk()),
		"allow_document_match":       framework.BoolOkToTF(apiObject.GetAllowDocumentMatchOk()),
		"allow_face_match":           framework.BoolOkToTF(apiObject.GetAllowFaceMatchOk()),
		"allow_manual_id_inspection": framework.BoolOkToTF(apiObject.GetAllowManualIdInspectionOk()),
	})
	diags.Append(d...)

	return objValue, diags
}
