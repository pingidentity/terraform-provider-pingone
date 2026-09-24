// Copyright © 2026 Ping Identity Corporation

package base

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/legacysdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
)

// Types
type PhoneDeliverySettingsDataSource serviceClientType

type PhoneDeliverySettingsDataSourceModel struct {
	Id                      pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId           pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	PhoneDeliverySettingsId pingonetypes.ResourceIDValue `tfsdk:"phone_delivery_settings_id"`
	Name                    types.String                 `tfsdk:"name"`
	ProviderType            types.String                 `tfsdk:"provider_type"`
	ProviderCustom          types.Object                 `tfsdk:"provider_custom"`
	ProviderCustomTwilio    types.Object                 `tfsdk:"provider_custom_twilio"`
	ProviderCustomSyniverse types.Object                 `tfsdk:"provider_custom_syniverse"`
	CreatedAt               timetypes.RFC3339            `tfsdk:"created_at"`
	UpdatedAt               timetypes.RFC3339            `tfsdk:"updated_at"`
}

// Framework interfaces
var (
	_ datasource.DataSource = &PhoneDeliverySettingsDataSource{}
)

// New Object
func NewPhoneDeliverySettingsDataSource() datasource.DataSource {
	return &PhoneDeliverySettingsDataSource{}
}

// Metadata
func (r *PhoneDeliverySettingsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_phone_delivery_settings"
}

// Schema
func (r *PhoneDeliverySettingsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	dataSourceExactlyOneOfRelativePaths := []string{
		"phone_delivery_settings_id",
		"name",
	}

	phoneDeliverySettingsIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the ID of the phone delivery settings to retrieve configuration for.  Must be a valid PingOne resource ID.",
	).ExactlyOneOf(dataSourceExactlyOneOfRelativePaths)

	nameDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the name of the phone delivery settings to retrieve configuration for.  Relevant only for settings of provider type `CUSTOM_PROVIDER`.",
	).ExactlyOneOf(dataSourceExactlyOneOfRelativePaths)

	providerTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the type of the phone delivery service.",
	).AllowedValuesEnum(management.AllowedEnumNotificationsSettingsPhoneDeliverySettingsProviderEnumValues)

	// Custom provider
	providerCustomDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single nested attribute with attributes that describe custom phone delivery settings.",
	)

	providerCustomNameDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the name of the custom provider used to identify in the PingOne platform.",
	)

	providerCustomAuthenticationMethodDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The custom provider account's authentication method.",
	).AllowedValuesComplex(map[string]string{
		string(management.ENUMNOTIFICATIONSSETTINGSPHONEDELIVERYSETTINGSCUSTOMAUTHMETHOD_BASIC):         "`username` and `password` parameters are required to be set",
		string(management.ENUMNOTIFICATIONSSETTINGSPHONEDELIVERYSETTINGSCUSTOMAUTHMETHOD_BEARER):        "`token` parameter is required to be set",
		string(management.ENUMNOTIFICATIONSSETTINGSPHONEDELIVERYSETTINGSCUSTOMAUTHMETHOD_OAUTH2):        "`auth_url` parameter is required to be set.  `grant_type` defaults to `CLIENT_CREDENTIALS`, where `client_id` and `client_secret` parameters are required to be set, or `JWT_BEARER`, where `assertion` parameter is required to be set",
		string(management.ENUMNOTIFICATIONSSETTINGSPHONEDELIVERYSETTINGSCUSTOMAUTHMETHOD_CUSTOM_HEADER): "`header_name` and `header_value` parameters are required to be set",
	})

	providerCustomAuthenticationUsernameDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("A string that specifies the username for the custom provider account. Required when `method` is `%s`", management.ENUMNOTIFICATIONSSETTINGSPHONEDELIVERYSETTINGSCUSTOMAUTHMETHOD_BASIC),
	)

	providerCustomAuthenticationPasswordDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("A string that specifies the password for the custom provider account. Required when `method` is `%s`.  This is a sensitive parameter and is not returned by the service.", management.ENUMNOTIFICATIONSSETTINGSPHONEDELIVERYSETTINGSCUSTOMAUTHMETHOD_BASIC),
	)

	providerCustomAuthenticationAuthTokenDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("A string that specifies the authentication token to use for the custom provider account. Required when `method` is `%s`.  This is a sensitive parameter and is not returned by the service.", management.ENUMNOTIFICATIONSSETTINGSPHONEDELIVERYSETTINGSCUSTOMAUTHMETHOD_BEARER),
	)

	providerCustomAuthenticationAuthUrlDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("A string that specifies the URL of the authorization server that issues the access token for the custom provider account. Required when `method` is `%s`", management.ENUMNOTIFICATIONSSETTINGSPHONEDELIVERYSETTINGSCUSTOMAUTHMETHOD_OAUTH2),
	)

	providerCustomAuthenticationGrantTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The grant type used to request the access token from the authorization server.  Relevant when `method` is `OAUTH2`.",
	).AllowedValuesEnum(management.AllowedEnumNotificationsSettingsPhoneDeliverySettingsCustomAuthGrantTypeEnumValues).DefaultValue(string(management.ENUMNOTIFICATIONSSETTINGSPHONEDELIVERYSETTINGSCUSTOMAUTHGRANTTYPE_CLIENT_CREDENTIALS))

	providerCustomAuthenticationAssertionDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("A string that specifies the JWT assertion used to request the access token from the authorization server.  Must be a valid JWT. Required when `grant_type` is `%s`.  This is a sensitive parameter and is not returned by the service.", management.ENUMNOTIFICATIONSSETTINGSPHONEDELIVERYSETTINGSCUSTOMAUTHGRANTTYPE_JWT_BEARER),
	)

	providerCustomAuthenticationClientIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("A string that specifies the client ID used to request the access token from the authorization server. Required when `grant_type` is `%s`", management.ENUMNOTIFICATIONSSETTINGSPHONEDELIVERYSETTINGSCUSTOMAUTHGRANTTYPE_CLIENT_CREDENTIALS),
	)

	providerCustomAuthenticationClientSecretDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("A string that specifies the client secret used to request the access token from the authorization server. Required when `grant_type` is `%s`.  This is a sensitive parameter and is not returned by the service.", management.ENUMNOTIFICATIONSSETTINGSPHONEDELIVERYSETTINGSCUSTOMAUTHGRANTTYPE_CLIENT_CREDENTIALS),
	)

	providerCustomAuthenticationScopesDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A set of strings that specifies the scopes to request in the access token from the authorization server, for example, `sms:send`, `voice:send`.",
	)

	providerCustomAuthenticationHeaderNameDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("A string that specifies the name of the custom header used to authenticate requests to the custom provider. Required when `method` is `%s`", management.ENUMNOTIFICATIONSSETTINGSPHONEDELIVERYSETTINGSCUSTOMAUTHMETHOD_CUSTOM_HEADER),
	)

	providerCustomAuthenticationHeaderValueDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("A string that specifies the value of the custom header used to authenticate requests to the custom provider. Required when `method` is `%s`.  This is a sensitive parameter and is not returned by the service.", management.ENUMNOTIFICATIONSSETTINGSPHONEDELIVERYSETTINGSCUSTOMAUTHMETHOD_CUSTOM_HEADER),
	)

	providerCustomAuthenticationClientAuthenticationMethodDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("A string that specifies the method used to send the OAuth 2.0 client credentials to the authorization server.  This is a read-only property computed by the service, relevant when `method` is `%s`.", management.ENUMNOTIFICATIONSSETTINGSPHONEDELIVERYSETTINGSCUSTOMAUTHMETHOD_OAUTH2),
	).AllowedValuesEnum(management.AllowedEnumNotificationsSettingsPhoneDeliverySettingsCustomAuthClientAuthenticationMethodEnumValues)

	providerCustomNumbersCapabilitiesDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A collection of the types of phone delivery service capabilities.",
	).AllowedValuesEnum(management.AllowedEnumNotificationsSettingsPhoneDeliverySettingsCustomNumbersCapabilityEnumValues)

	providerCustomNumbersSupportedCountriesDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Specifies the `number`'s supported countries for notification recipients, depending on the phone number type.  If an SMS template has an alphanumeric `sender` ID and also has short code, the `sender` ID will be used for destination countries that support both alphanumeric senders and short codes. For Unites States and Canada that don't support alphanumeric sender IDs, a short code will be used if both an alphanumeric sender and a short code are specified.",
	)

	providerCustomNumbersTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the type of phone number.",
	).AllowedValuesEnum(management.AllowedEnumNotificationsSettingsPhoneDeliverySettingsCustomNumbersTypeEnumValues)

	providerCustomRequestsAfterTagDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"For voice OTP notifications only.  A string that specifies a closing tag which is commonly used by custom providers for defining a pause between each number in the OTP number string.  Example value: `</Say> <Pause length=\"1\"/>`",
	)

	providerCustomRequestsBeforeTagDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"For voice OTP notifications only.  A string that specifies an opening tag which is commonly used by custom providers for defining a pause between each number in the OTP number string.  Possible value: `<Say>`.",
	)

	providerCustomRequestsBodyDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Optional when the `method` is `POST`.  A string that specifies the notification's request body. The body should include the `${to}` and `${message}` mandatory variables.",
	)

	providerCustomRequestsDeliveryMethodDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the notification's delivery method.",
	).AllowedValuesEnum(management.AllowedEnumNotificationsSettingsPhoneDeliverySettingsCustomDeliveryMethodEnumValues)

	providerCustomRequestsHeadersDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A map of strings that specifies the notification's request headers, matching the format of the request body.",
	)

	providerCustomRequestsMethodDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the type of HTTP request method.",
	).AllowedValuesEnum(management.AllowedEnumNotificationsSettingsPhoneDeliverySettingsCustomRequestMethodEnumValues)

	providerCustomRequestsPhoneNumberFormatDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the phone number format.",
	).AllowedValuesComplex(map[string]string{
		string(management.ENUMNOTIFICATIONSSETTINGSPHONEDELIVERYSETTINGSCUSTOMNUMBERFORMAT_FULL):        "The phone number format with a leading `+` sign, in the E.164 standard format.  For example: `+14155552671`",
		string(management.ENUMNOTIFICATIONSSETTINGSPHONEDELIVERYSETTINGSCUSTOMNUMBERFORMAT_NUMBER_ONLY): "The phone number format without a leading `+` sign, in the E.164 standard format.  For example: `14155552671`",
	})

	providerCustomRequestsUrlDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The provider's remote gateway or customer gateway URL.",
	)

	// Twilio provider
	providerCustomTwilioDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single nested attribute with attributes that describe phone delivery settings for a custom Twilio account.",
	)

	providerCustomTwilioAuthTokenDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The secret key of the Twilio account.  This is a sensitive parameter and is not returned by the service.",
	)

	providerCustomTwilioSidDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The public ID of the Twilio account.",
	)

	// Syniverse provider
	providerCustomSyniverseDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single nested attribute with attributes that describe phone delivery settings for a custom syniverse account.",
	)

	providerCustomSyniverseAuthTokenDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The secret key of the Syniverse account.  This is a sensitive parameter and is not returned by the service.",
	)

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Datasource to retrieve phone delivery settings in a PingOne environment by ID or by name.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment that is configured with the phone delivery settings."),
			),

			"phone_delivery_settings_id": schema.StringAttribute{
				Description:         phoneDeliverySettingsIdDescription.Description,
				MarkdownDescription: phoneDeliverySettingsIdDescription.MarkdownDescription,
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
						path.MatchRelative().AtParent().AtName("phone_delivery_settings_id"),
					),
				},
			},

			"provider_type": schema.StringAttribute{
				Description:         providerTypeDescription.Description,
				MarkdownDescription: providerTypeDescription.MarkdownDescription,
				Computed:            true,
			},

			"provider_custom": schema.SingleNestedAttribute{
				Description:         providerCustomDescription.Description,
				MarkdownDescription: providerCustomDescription.MarkdownDescription,
				Computed:            true,

				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Description: providerCustomNameDescription.Description,
						Computed:    true,
					},

					"authentication": schema.SingleNestedAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that provides authentication settings for authenticating to the custom service API.").Description,
						Computed:    true,

						Attributes: map[string]schema.Attribute{
							"method": schema.StringAttribute{
								Description:         providerCustomAuthenticationMethodDescription.Description,
								MarkdownDescription: providerCustomAuthenticationMethodDescription.MarkdownDescription,
								Computed:            true,
							},

							"username": schema.StringAttribute{
								Description: providerCustomAuthenticationUsernameDescription.Description,
								Computed:    true,
							},

							"password": schema.StringAttribute{
								Description:         providerCustomAuthenticationPasswordDescription.Description,
								MarkdownDescription: providerCustomAuthenticationPasswordDescription.MarkdownDescription,
								Computed:            true,
								Sensitive:           true,
							},

							"auth_token": schema.StringAttribute{
								Description:         providerCustomAuthenticationAuthTokenDescription.Description,
								MarkdownDescription: providerCustomAuthenticationAuthTokenDescription.MarkdownDescription,
								Computed:            true,
								Sensitive:           true,
							},

							"auth_url": schema.StringAttribute{
								Description:         providerCustomAuthenticationAuthUrlDescription.Description,
								MarkdownDescription: providerCustomAuthenticationAuthUrlDescription.MarkdownDescription,
								Computed:            true,
							},

							"grant_type": schema.StringAttribute{
								Description:         providerCustomAuthenticationGrantTypeDescription.Description,
								MarkdownDescription: providerCustomAuthenticationGrantTypeDescription.MarkdownDescription,
								Computed:            true,
							},

							"assertion": schema.StringAttribute{
								Description:         providerCustomAuthenticationAssertionDescription.Description,
								MarkdownDescription: providerCustomAuthenticationAssertionDescription.MarkdownDescription,
								Computed:            true,
								Sensitive:           true,
							},

							"client_id": schema.StringAttribute{
								Description:         providerCustomAuthenticationClientIdDescription.Description,
								MarkdownDescription: providerCustomAuthenticationClientIdDescription.MarkdownDescription,
								Computed:            true,
							},

							"client_secret": schema.StringAttribute{
								Description:         providerCustomAuthenticationClientSecretDescription.Description,
								MarkdownDescription: providerCustomAuthenticationClientSecretDescription.MarkdownDescription,
								Computed:            true,
								Sensitive:           true,
							},

							"scopes": schema.SetAttribute{
								Description:         providerCustomAuthenticationScopesDescription.Description,
								MarkdownDescription: providerCustomAuthenticationScopesDescription.MarkdownDescription,
								Computed:            true,

								ElementType: types.StringType,
							},

							"header_name": schema.StringAttribute{
								Description:         providerCustomAuthenticationHeaderNameDescription.Description,
								MarkdownDescription: providerCustomAuthenticationHeaderNameDescription.MarkdownDescription,
								Computed:            true,
							},

							"header_value": schema.StringAttribute{
								Description:         providerCustomAuthenticationHeaderValueDescription.Description,
								MarkdownDescription: providerCustomAuthenticationHeaderValueDescription.MarkdownDescription,
								Computed:            true,
								Sensitive:           true,
							},

							"client_authentication_method": schema.StringAttribute{
								Description:         providerCustomAuthenticationClientAuthenticationMethodDescription.Description,
								MarkdownDescription: providerCustomAuthenticationClientAuthenticationMethodDescription.MarkdownDescription,
								Computed:            true,
							},
						},
					},

					"numbers": schema.SetNestedAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("One or more objects that describe the numbers to use for phone delivery.").Description,
						Computed:    true,

						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"available": schema.BoolAttribute{
									Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether the number is currently available in the provider account.").Description,
									Computed:    true,
								},

								"capabilities": schema.SetAttribute{
									Description:         providerCustomNumbersCapabilitiesDescription.Description,
									MarkdownDescription: providerCustomNumbersCapabilitiesDescription.MarkdownDescription,
									Computed:            true,

									ElementType: types.StringType,
								},

								"number": schema.StringAttribute{
									Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the phone number, toll-free number or short code.").Description,
									Computed:    true,
								},

								"selected": schema.BoolAttribute{
									Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether the number is currently available in the provider account.").Description,
									Computed:    true,
								},

								"supported_countries": schema.SetAttribute{
									Description:         providerCustomNumbersSupportedCountriesDescription.Description,
									MarkdownDescription: providerCustomNumbersSupportedCountriesDescription.MarkdownDescription,
									Computed:            true,

									ElementType: types.StringType,
								},

								"type": schema.StringAttribute{
									Description:         providerCustomNumbersTypeDescription.Description,
									MarkdownDescription: providerCustomNumbersTypeDescription.MarkdownDescription,
									Computed:            true,
								},
							},
						},
					},

					"requests": schema.SetNestedAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("One or more objects that describe the outbound custom notification requests.").Description,
						Computed:    true,

						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"after_tag": schema.StringAttribute{
									Description: providerCustomRequestsAfterTagDescription.Description,
									Computed:    true,
								},

								"before_tag": schema.StringAttribute{
									Description: providerCustomRequestsBeforeTagDescription.Description,
									Computed:    true,
								},

								"body": schema.StringAttribute{
									Description: providerCustomRequestsBodyDescription.Description,
									Computed:    true,
								},

								"delivery_method": schema.StringAttribute{
									Description:         providerCustomRequestsDeliveryMethodDescription.Description,
									MarkdownDescription: providerCustomRequestsDeliveryMethodDescription.MarkdownDescription,
									Computed:            true,
								},

								"headers": schema.MapAttribute{
									Description:         providerCustomRequestsHeadersDescription.Description,
									MarkdownDescription: providerCustomRequestsHeadersDescription.MarkdownDescription,
									Computed:            true,

									ElementType: types.StringType,
								},

								"method": schema.StringAttribute{
									Description:         providerCustomRequestsMethodDescription.Description,
									MarkdownDescription: providerCustomRequestsMethodDescription.MarkdownDescription,
									Computed:            true,
								},

								"phone_number_format": schema.StringAttribute{
									Description:         providerCustomRequestsPhoneNumberFormatDescription.Description,
									MarkdownDescription: providerCustomRequestsPhoneNumberFormatDescription.MarkdownDescription,
									Computed:            true,
								},

								"url": schema.StringAttribute{
									Description:         providerCustomRequestsUrlDescription.Description,
									MarkdownDescription: providerCustomRequestsUrlDescription.MarkdownDescription,
									Computed:            true,
								},
							},
						},
					},
				},
			},

			"provider_custom_twilio": schema.SingleNestedAttribute{
				Description:         providerCustomTwilioDescription.Description,
				MarkdownDescription: providerCustomTwilioDescription.MarkdownDescription,
				Computed:            true,

				Attributes: map[string]schema.Attribute{
					"auth_token": schema.StringAttribute{
						Description:         providerCustomTwilioAuthTokenDescription.Description,
						MarkdownDescription: providerCustomTwilioAuthTokenDescription.MarkdownDescription,
						Computed:            true,
						Sensitive:           true,
					},

					"sid": schema.StringAttribute{
						Description:         providerCustomTwilioSidDescription.Description,
						MarkdownDescription: providerCustomTwilioSidDescription.MarkdownDescription,
						Computed:            true,
					},

					"selected_numbers": schema.SetNestedAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("One or more objects that describe the numbers selected for phone delivery.").Description,
						Computed:    true,

						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"number": schema.StringAttribute{
									Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the phone number, toll-free number or short code that has been configured in Twilio.").Description,
									Computed:    true,
								},

								"selected": schema.BoolAttribute{
									Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether the number is currently available in the provider account.").Description,
									Computed:    true,
								},

								"supported_countries": schema.SetAttribute{
									Description:         providerCustomNumbersSupportedCountriesDescription.Description,
									MarkdownDescription: providerCustomNumbersSupportedCountriesDescription.MarkdownDescription,
									Computed:            true,

									ElementType: types.StringType,
								},

								"type": schema.StringAttribute{
									Description:         providerCustomNumbersTypeDescription.Description,
									MarkdownDescription: providerCustomNumbersTypeDescription.MarkdownDescription,
									Computed:            true,
								},
							},
						},
					},

					"service_numbers": schema.SetNestedAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("One or more objects that describe the numbers that are defined in the Twilio service.").Description,
						Computed:    true,

						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"available": schema.BoolAttribute{
									Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether the number is currently available in the provider account.").Description,
									Computed:    true,
								},

								"capabilities": schema.SetAttribute{
									Description:         providerCustomNumbersCapabilitiesDescription.Description,
									MarkdownDescription: providerCustomNumbersCapabilitiesDescription.MarkdownDescription,
									Computed:            true,

									ElementType: types.StringType,
								},

								"number": schema.StringAttribute{
									Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the phone number, toll-free number or short code.").Description,
									Computed:    true,
								},

								"selected": schema.BoolAttribute{
									Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether the number is currently available in the provider account.").Description,
									Computed:    true,
								},

								"supported_countries": schema.SetAttribute{
									Description:         providerCustomNumbersSupportedCountriesDescription.Description,
									MarkdownDescription: providerCustomNumbersSupportedCountriesDescription.MarkdownDescription,
									Computed:            true,

									ElementType: types.StringType,
								},

								"type": schema.StringAttribute{
									Description:         providerCustomNumbersTypeDescription.Description,
									MarkdownDescription: providerCustomNumbersTypeDescription.MarkdownDescription,
									Computed:            true,
								},
							},
						},
					},
				},
			},

			"provider_custom_syniverse": schema.SingleNestedAttribute{
				Description:         providerCustomSyniverseDescription.Description,
				MarkdownDescription: providerCustomSyniverseDescription.MarkdownDescription,
				Computed:            true,

				Attributes: map[string]schema.Attribute{
					"auth_token": schema.StringAttribute{
						Description:         providerCustomSyniverseAuthTokenDescription.Description,
						MarkdownDescription: providerCustomSyniverseAuthTokenDescription.MarkdownDescription,
						Computed:            true,
						Sensitive:           true,
					},

					"selected_numbers": schema.SetNestedAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("One or more objects that describe the numbers selected for phone delivery.").Description,
						Computed:    true,

						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"number": schema.StringAttribute{
									Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the phone number, toll-free number or short code that has been configured in Syniverse.").Description,
									Computed:    true,
								},

								"selected": schema.BoolAttribute{
									Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether the number is currently available in the provider account.").Description,
									Computed:    true,
								},

								"supported_countries": schema.SetAttribute{
									Description:         providerCustomNumbersSupportedCountriesDescription.Description,
									MarkdownDescription: providerCustomNumbersSupportedCountriesDescription.MarkdownDescription,
									Computed:            true,

									ElementType: types.StringType,
								},

								"type": schema.StringAttribute{
									Description:         providerCustomNumbersTypeDescription.Description,
									MarkdownDescription: providerCustomNumbersTypeDescription.MarkdownDescription,
									Computed:            true,
								},
							},
						},
					},

					"service_numbers": schema.SetNestedAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("One or more objects that describe the numbers that are defined in the Syniverse service.").Description,
						Computed:    true,

						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"available": schema.BoolAttribute{
									Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether the number is currently available in the provider account.").Description,
									Computed:    true,
								},

								"capabilities": schema.SetAttribute{
									Description:         providerCustomNumbersCapabilitiesDescription.Description,
									MarkdownDescription: providerCustomNumbersCapabilitiesDescription.MarkdownDescription,
									Computed:            true,

									ElementType: types.StringType,
								},

								"number": schema.StringAttribute{
									Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the phone number, toll-free number or short code.").Description,
									Computed:    true,
								},

								"selected": schema.BoolAttribute{
									Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether the number is currently available in the provider account.").Description,
									Computed:    true,
								},

								"supported_countries": schema.SetAttribute{
									Description:         providerCustomNumbersSupportedCountriesDescription.Description,
									MarkdownDescription: providerCustomNumbersSupportedCountriesDescription.MarkdownDescription,
									Computed:            true,

									ElementType: types.StringType,
								},

								"type": schema.StringAttribute{
									Description:         providerCustomNumbersTypeDescription.Description,
									MarkdownDescription: providerCustomNumbersTypeDescription.MarkdownDescription,
									Computed:            true,
								},
							},
						},
					},
				},
			},

			"created_at": schema.StringAttribute{
				Description: "A string that specifies the time the resource was created.",
				Computed:    true,

				CustomType: timetypes.RFC3339Type{},
			},

			"updated_at": schema.StringAttribute{
				Description: "A string that specifies the time the resource was last updated.",
				Computed:    true,

				CustomType: timetypes.RFC3339Type{},
			},
		},
	}
}

func (r *PhoneDeliverySettingsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *PhoneDeliverySettingsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *PhoneDeliverySettingsDataSourceModel

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

	var response *management.NotificationsSettingsPhoneDeliverySettings

	if !data.PhoneDeliverySettingsId.IsNull() {

		// Run the API call
		resp.Diagnostics.Append(legacysdk.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				fO, fR, fErr := r.Client.ManagementAPIClient.PhoneDeliverySettingsApi.ReadOnePhoneDeliverySettings(ctx, data.EnvironmentId.ValueString(), data.PhoneDeliverySettingsId.ValueString()).Execute()
				return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
			},
			"ReadOnePhoneDeliverySettings",
			legacysdk.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			&response,
		)...)
		if resp.Diagnostics.HasError() {
			return
		}

		if response == nil {
			resp.Diagnostics.AddError(
				"Cannot find phone delivery settings from ID",
				fmt.Sprintf("The phone delivery settings %s for environment %s cannot be found", data.PhoneDeliverySettingsId.String(), data.EnvironmentId.String()),
			)
			return
		}

	} else if !data.Name.IsNull() {

		// Run the API call
		resp.Diagnostics.Append(legacysdk.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				pagedIterator := r.Client.ManagementAPIClient.PhoneDeliverySettingsApi.ReadAllPhoneDeliverySettings(ctx, data.EnvironmentId.ValueString()).Execute()

				var initialHttpResponse *http.Response

				for pageCursor, err := range pagedIterator {
					if err != nil {
						return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, pageCursor.HTTPResponse, err)
					}

					if initialHttpResponse == nil {
						initialHttpResponse = pageCursor.HTTPResponse
					}

					if phoneDeliverySettings, ok := pageCursor.EntityArray.Embedded.GetPhoneDeliverySettingsOk(); ok {
						for _, phoneDeliverySettingsItem := range phoneDeliverySettings {

							// Only custom provider settings have a name to match
							if v := phoneDeliverySettingsItem.NotificationsSettingsPhoneDeliverySettingsCustom; v != nil {
								if strings.EqualFold(v.GetName(), data.Name.ValueString()) {
									return &phoneDeliverySettingsItem, pageCursor.HTTPResponse, nil
								}
							}
						}
					}
				}

				return nil, initialHttpResponse, nil
			},
			"ReadAllPhoneDeliverySettings",
			legacysdk.DefaultCustomError,
			sdk.DefaultCreateReadRetryable,
			&response,
		)...)
		if resp.Diagnostics.HasError() {
			return
		}

		if response == nil {
			resp.Diagnostics.AddError(
				"Cannot find phone delivery settings from name",
				fmt.Sprintf("The phone delivery settings name %s for environment %s cannot be found", data.Name.String(), data.EnvironmentId.String()),
			)
			return
		}

	} else {
		resp.Diagnostics.AddError(
			"Missing parameter",
			"Cannot find the requested phone delivery settings. phone_delivery_settings_id or name must be set.",
		)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(ctx, response)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (p *PhoneDeliverySettingsDataSourceModel) toState(ctx context.Context, apiObject *management.NotificationsSettingsPhoneDeliverySettings) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	apiObjectCommon := management.NotificationsSettingsPhoneDeliverySettingsCommon{}

	if v := apiObject.NotificationsSettingsPhoneDeliverySettingsCustom; v != nil {
		apiObjectCommon = management.NotificationsSettingsPhoneDeliverySettingsCommon{
			Id:          v.Id,
			Environment: v.Environment,
			Provider:    v.Provider,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
		}
	}

	if v := apiObject.NotificationsSettingsPhoneDeliverySettingsTwilioSyniverse; v != nil {
		apiObjectCommon = management.NotificationsSettingsPhoneDeliverySettingsCommon{
			Id:          v.Id,
			Environment: v.Environment,
			Provider:    v.Provider,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
		}
	}

	p.Id = framework.PingOneResourceIDOkToTF(apiObjectCommon.GetIdOk())
	p.EnvironmentId = framework.PingOneResourceIDToTF(*apiObjectCommon.GetEnvironment().Id)
	p.PhoneDeliverySettingsId = framework.PingOneResourceIDOkToTF(apiObjectCommon.GetIdOk())
	p.ProviderType = framework.EnumOkToTF(apiObjectCommon.GetProviderOk())
	p.CreatedAt = framework.TimeOkToTF(apiObjectCommon.GetCreatedAtOk())
	p.UpdatedAt = framework.TimeOkToTF(apiObjectCommon.GetUpdatedAtOk())

	var d diag.Diagnostics

	if p.ProviderType.Equal(types.StringValue(string(management.ENUMNOTIFICATIONSSETTINGSPHONEDELIVERYSETTINGSPROVIDER_PROVIDER))) {
		p.ProviderCustom, d = toStatePhoneDeliverySettingsProviderCustom(ctx, nil, apiObject.NotificationsSettingsPhoneDeliverySettingsCustom)
		diags.Append(d...)
	} else {
		p.ProviderCustom = types.ObjectNull(customTFObjectTypes)
	}

	if p.ProviderType.Equal(types.StringValue(string(management.ENUMNOTIFICATIONSSETTINGSPHONEDELIVERYSETTINGSPROVIDER_TWILIO))) {
		p.ProviderCustomTwilio, d = toStatePhoneDeliverySettingsProviderCustomTwilio(ctx, nil, apiObject.NotificationsSettingsPhoneDeliverySettingsTwilioSyniverse)
		diags.Append(d...)
	} else {
		p.ProviderCustomTwilio = types.ObjectNull(twilioTFObjectTypes)
	}

	if p.ProviderType.Equal(types.StringValue(string(management.ENUMNOTIFICATIONSSETTINGSPHONEDELIVERYSETTINGSPROVIDER_SYNIVERSE))) {
		p.ProviderCustomSyniverse, d = toStatePhoneDeliverySettingsProviderCustomSyniverse(ctx, nil, apiObject.NotificationsSettingsPhoneDeliverySettingsTwilioSyniverse)
		diags.Append(d...)
	} else {
		p.ProviderCustomSyniverse = types.ObjectNull(syniverseTFObjectTypes)
	}

	return diags
}
