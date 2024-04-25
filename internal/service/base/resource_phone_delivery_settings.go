package base

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	setvalidatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/setvalidator"
	stringvalidatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/stringvalidator"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type PhoneDeliverySettingsResource serviceClientType

type PhoneDeliverySettingsResourceModel struct {
	Id                      pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId           pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	ProviderType            types.String                 `tfsdk:"provider_type"`
	ProviderCustom          types.Object                 `tfsdk:"provider_custom"`
	ProviderCustomTwilio    types.Object                 `tfsdk:"provider_custom_twilio"`
	ProviderCustomSyniverse types.Object                 `tfsdk:"provider_custom_syniverse"`
	CreatedAt               timetypes.RFC3339            `tfsdk:"created_at"`
	UpdatedAt               timetypes.RFC3339            `tfsdk:"updated_at"`
}

type PhoneDeliverySettingsProviderCustomResourceModel struct {
	Authentication types.Object `tfsdk:"authentication"`
	Name           types.String `tfsdk:"name"`
	Numbers        types.Set    `tfsdk:"numbers"`
	Requests       types.Set    `tfsdk:"requests"`
}

type PhoneDeliverySettingsProviderCustomAuthenticationResourceModel struct {
	Method    types.String `tfsdk:"method"`
	Password  types.String `tfsdk:"password"`
	AuthToken types.String `tfsdk:"auth_token"`
	Username  types.String `tfsdk:"username"`
}

type PhoneDeliverySettingsProviderCustomNumbersResourceModel struct {
	SupportedCountries types.Set    `tfsdk:"supported_countries"`
	Type               types.String `tfsdk:"type"`
	Selected           types.Bool   `tfsdk:"selected"`
	Available          types.Bool   `tfsdk:"available"`
	Number             types.String `tfsdk:"number"`
	Capabilities       types.Set    `tfsdk:"capabilities"`
}

type PhoneDeliverySettingsProviderCustomSelectedNumbersResourceModel struct {
	SupportedCountries types.Set    `tfsdk:"supported_countries"`
	Type               types.String `tfsdk:"type"`
	Selected           types.Bool   `tfsdk:"selected"`
	Number             types.String `tfsdk:"number"`
}

type PhoneDeliverySettingsProviderCustomRequestsResourceModel struct {
	DeliveryMethod    types.String `tfsdk:"delivery_method"`
	Url               types.String `tfsdk:"url"`
	Method            types.String `tfsdk:"method"`
	Body              types.String `tfsdk:"body"`
	Headers           types.Map    `tfsdk:"headers"`
	BeforeTag         types.String `tfsdk:"before_tag"`
	AfterTag          types.String `tfsdk:"after_tag"`
	PhoneNumberFormat types.String `tfsdk:"phone_number_format"`
}

type PhoneDeliverySettingsProviderCustomTwilioResourceModel struct {
	Sid             types.String `tfsdk:"sid"`
	AuthToken       types.String `tfsdk:"auth_token"`
	SelectedNumbers types.Set    `tfsdk:"selected_numbers"`
	ServiceNumbers  types.Set    `tfsdk:"service_numbers"`
}

type PhoneDeliverySettingsProviderCustomSyniverseResourceModel struct {
	AuthToken       types.String `tfsdk:"auth_token"`
	SelectedNumbers types.Set    `tfsdk:"selected_numbers"`
	ServiceNumbers  types.Set    `tfsdk:"service_numbers"`
}

var (
	customTFObjectTypes = map[string]attr.Type{
		"authentication": types.ObjectType{
			AttrTypes: customAuthenticationTFObjectTypes,
		},
		"name": types.StringType,
		"numbers": types.SetType{ElemType: types.ObjectType{
			AttrTypes: customNumbersTFObjectTypes,
		}},
		"requests": types.SetType{ElemType: types.ObjectType{
			AttrTypes: customRequestsTFObjectTypes,
		}},
	}

	customAuthenticationTFObjectTypes = map[string]attr.Type{
		"method":     types.StringType,
		"password":   types.StringType,
		"auth_token": types.StringType,
		"username":   types.StringType,
	}

	customNumbersTFObjectTypes = map[string]attr.Type{
		"available":           types.BoolType,
		"capabilities":        types.SetType{ElemType: types.StringType},
		"number":              types.StringType,
		"selected":            types.BoolType,
		"supported_countries": types.SetType{ElemType: types.StringType},
		"type":                types.StringType,
	}

	customRequestsTFObjectTypes = map[string]attr.Type{
		"after_tag":           types.StringType,
		"before_tag":          types.StringType,
		"body":                types.StringType,
		"delivery_method":     types.StringType,
		"headers":             types.MapType{ElemType: types.StringType},
		"method":              types.StringType,
		"phone_number_format": types.StringType,
		"url":                 types.StringType,
	}

	twilioTFObjectTypes = map[string]attr.Type{
		"auth_token": types.StringType,
		"sid":        types.StringType,
		"selected_numbers": types.SetType{ElemType: types.ObjectType{
			AttrTypes: customSelectedNumbersTFObjectTypes,
		}},
		"service_numbers": types.SetType{ElemType: types.ObjectType{
			AttrTypes: customNumbersTFObjectTypes,
		}},
	}

	syniverseTFObjectTypes = map[string]attr.Type{
		"auth_token": types.StringType,
		"selected_numbers": types.SetType{ElemType: types.ObjectType{
			AttrTypes: customSelectedNumbersTFObjectTypes,
		}},
		"service_numbers": types.SetType{ElemType: types.ObjectType{
			AttrTypes: customNumbersTFObjectTypes,
		}},
	}

	customSelectedNumbersTFObjectTypes = map[string]attr.Type{
		"number":              types.StringType,
		"selected":            types.BoolType,
		"supported_countries": types.SetType{ElemType: types.StringType},
		"type":                types.StringType,
	}
)

// Framework interfaces
var (
	_ resource.Resource                = &PhoneDeliverySettingsResource{}
	_ resource.ResourceWithConfigure   = &PhoneDeliverySettingsResource{}
	_ resource.ResourceWithImportState = &PhoneDeliverySettingsResource{}
)

// New Object
func NewPhoneDeliverySettingsResource() resource.Resource {
	return &PhoneDeliverySettingsResource{}
}

// Metadata
func (r *PhoneDeliverySettingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_phone_delivery_settings"
}

// Schema.
func (r *PhoneDeliverySettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	const attrMinLength = 1

	providerTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the type of the phone delivery service.",
	).AllowedValuesEnum(management.AllowedEnumNotificationsSettingsPhoneDeliverySettingsProviderEnumValues)

	// Custom provider
	providerCustomDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single nested attribute with attributes that describe custom phone delivery settings.",
	).ExactlyOneOf([]string{
		"provider_custom",
		"provider_custom_twilio",
		"provider_custom_syniverse",
	})

	providerCustomAuthenticationMethodDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The custom provider account's authentication method.",
	).AllowedValuesComplex(map[string]string{
		string(management.ENUMNOTIFICATIONSSETTINGSPHONEDELIVERYSETTINGSCUSTOMAUTHMETHOD_BASIC):  "`username` and `password` parameters are required to be set",
		string(management.ENUMNOTIFICATIONSSETTINGSPHONEDELIVERYSETTINGSCUSTOMAUTHMETHOD_BEARER): "`token` parameter is required to be set",
	})

	providerCustomAuthenticationUsernameDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("A string that specifies the username for the custom provider account. Required when `method` is `%s`", management.ENUMNOTIFICATIONSSETTINGSPHONEDELIVERYSETTINGSCUSTOMAUTHMETHOD_BASIC),
	)

	providerCustomAuthenticationPasswordDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("A string that specifies the password for the custom provider account. Required when `method` is `%s`", management.ENUMNOTIFICATIONSSETTINGSPHONEDELIVERYSETTINGSCUSTOMAUTHMETHOD_BASIC),
	)

	providerCustomAuthenticationAuthTokenDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("A string that specifies the authentication token to use for the custom provider account. Required when `method` is `%s`", management.ENUMNOTIFICATIONSSETTINGSPHONEDELIVERYSETTINGSCUSTOMAUTHMETHOD_BEARER),
	)

	providerCustomNumbersCapabilitiesDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A collection of the types of phone delivery service capabilities.",
	).AllowedValuesEnum(management.AllowedEnumNotificationsSettingsPhoneDeliverySettingsCustomNumbersCapabilityEnumValues)

	providerCustomNumbersSupportedCountriesDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"Specifies the `number`'s supported countries for notification recipients, depending on the phone number type.  If an SMS template has an alphanumeric `sender` ID and also has short code, the `sender` ID will be used for destination countries that support both alphanumeric senders and short codes. For Unites States and Canada that don't support alphanumeric sender IDs, a short code will be used if both an alphanumeric sender and a short code are specified.\n" +
			"    - `SHORT_CODE`: A collection containing a single 2-character ISO country code, for example, `US`, `GB`, `CA`.\n" +
			"    If the custom provider is of `type` `CUSTOM_PROVIDER`, this attribute must not be empty or null.\n" +
			"    For other custom provider types, if this attribute is null (empty is not supported), the specified short code `number` can only be used to dispatch notifications to United States recipient numbers.\n" +
			"    - `TOLL_FREE`: A collection of valid 2-character country ISO codes, for example, `US`, `GB`, `CA`.\n" +
			"    If the custom provider is of `type` `CUSTOM_PROVIDER`, this attribute must not be empty or null.\n" +
			"    For other custom provider types, if this attribute is null (empty is not supported), the specified toll-free `number` can only be used to dispatch notifications to United States recipient numbers.\n" +
			"    - `PHONE_NUMBER`: this attribute cannot be specified.",
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
		"Optional when the `method` is `POST`.  A string that specifies the notification's request body. The body should include the `${to}` and `${message}` mandatory variables. For some vendors, the optional `${from}` variable may also be required. For example `messageType=ARN&message=${message}&phoneNumber=${to}&sender=${from}`.  In addition, you can use [dynamic variables](https://apidocs.pingidentity.com/pingone/platform/v1/api/#notifications-templates-dynamic-variables) and the following optional variables:\n" +
			"    - `${voice}` - the type of voice configured for notifications\n" +
			"    - `${locale}` - locale\n" +
			"    - `${otp}` - OTP\n" +
			"    - `${user.username}` - user's username\n" +
			"    - `${user.name.given}` - user's given name\n" +
			"    - `${user.name.family}` - user's family name",
	)

	providerCustomRequestsDeliveryMethodDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the notification's delivery method.",
	).AllowedValuesEnum(management.AllowedEnumNotificationsSettingsPhoneDeliverySettingsCustomDeliveryMethodEnumValues)

	providerCustomRequestsHeadersDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A map of strings that specifies the notification's request headers, matching the format of the request body. The header should include only one of the following if the `method` is set to `POST`:\n" +
			"    - `content-type` = `application/x-www-form-urlencoded` (where the `body` should be form encoded)\n" +
			"    - `content-type` = `application/json` (where the `body` should be JSON encoded)",
	)

	providerCustomRequestsMethodDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the type of HTTP request method.",
	).AllowedValuesEnum(management.AllowedEnumNotificationsSettingsPhoneDeliverySettingsCustomRequestMethodEnumValues)

	providerCustomRequestsPhoneNumberFormatDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the phone number format.",
	).AllowedValuesComplex(map[string]string{
		string(management.ENUMNOTIFICATIONSSETTINGSPHONEDELIVERYSETTINGSCUSTOMNUMBERFORMAT_FULL):        "The phone number format with a leading `+` sign, in the E.164 standard format.  For example: `+14155552671`",
		string(management.ENUMNOTIFICATIONSSETTINGSPHONEDELIVERYSETTINGSCUSTOMNUMBERFORMAT_NUMBER_ONLY): "The phone number format without a leading `+` sign, in the E.164 standard format.  For example: `14155552671`",
	}).DefaultValue(string(management.ENUMNOTIFICATIONSSETTINGSPHONEDELIVERYSETTINGSCUSTOMNUMBERFORMAT_FULL))

	providerCustomRequestsUrlDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The provider's remote gateway or customer gateway URL.  For requests using the `POST` method, use the provider's remote gateway URL.  For requests using the `GET` method, use the provider's remote gateway URL, including the `${to}` and `${message}` mandatory variables, and the optional `${from}` variable, for example: `https://api.transmitsms.com/send-sms.json?to=${to}&from=${from}&message=${message}`",
	)

	// Twilio provider
	providerCustomTwilioDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single nested attribute with attributes that describe phone delivery settings for a custom Twilio account.",
	).ExactlyOneOf([]string{
		"provider_custom",
		"provider_custom_twilio",
		"provider_custom_syniverse",
	})

	providerCustomTwilioAuthTokenDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The secret key of the Twilio account.",
	).RequiresReplace()

	providerCustomTwilioSidDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The public ID of the Twilio account.",
	).RequiresReplace()

	// Syniverse provider
	providerCustomSyniverseDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single nested attribute with attributes that describe phone delivery settings for a custom syniverse account.",
	).ExactlyOneOf([]string{
		"provider_custom",
		"provider_custom_twilio",
		"provider_custom_syniverse",
	})

	providerCustomSyniverseAuthTokenDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The secret key of the Syniverse account.",
	).RequiresReplace()

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage SMS/voice delivery settings in a PingOne environment.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to configure SMS/voice settings for."),
			),

			"provider_type": schema.StringAttribute{
				Description:         providerTypeDescription.Description,
				MarkdownDescription: providerTypeDescription.MarkdownDescription,
				Computed:            true,

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"provider_custom": schema.SingleNestedAttribute{
				Description:         providerCustomDescription.Description,
				MarkdownDescription: providerCustomDescription.MarkdownDescription,
				Optional:            true,

				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("The string that specifies the name of the custom provider used to identify in the PingOne platform.").Description,
						Required:    true,
					},

					"authentication": schema.SingleNestedAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that provides authentication settings for authenticating to the custom service API.").Description,
						Required:    true,

						Attributes: map[string]schema.Attribute{
							"method": schema.StringAttribute{
								Description:         providerCustomAuthenticationMethodDescription.Description,
								MarkdownDescription: providerCustomAuthenticationMethodDescription.MarkdownDescription,
								Required:            true,

								Validators: []validator.String{
									stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumNotificationsSettingsPhoneDeliverySettingsCustomAuthMethodEnumValues)...),
								},
							},

							"username": schema.StringAttribute{
								Description:         providerCustomAuthenticationUsernameDescription.Description,
								MarkdownDescription: providerCustomAuthenticationUsernameDescription.MarkdownDescription,
								Optional:            true,

								Validators: []validator.String{
									stringvalidatorinternal.IsRequiredIfMatchesPathValue(
										types.StringValue(string(management.ENUMNOTIFICATIONSSETTINGSPHONEDELIVERYSETTINGSCUSTOMAUTHMETHOD_BASIC)),
										path.MatchRelative().AtParent().AtName("method"),
									),
								},
							},

							"password": schema.StringAttribute{
								Description:         providerCustomAuthenticationPasswordDescription.Description,
								MarkdownDescription: providerCustomAuthenticationPasswordDescription.MarkdownDescription,
								Optional:            true,
								Sensitive:           true,

								Validators: []validator.String{
									stringvalidatorinternal.IsRequiredIfMatchesPathValue(
										types.StringValue(string(management.ENUMNOTIFICATIONSSETTINGSPHONEDELIVERYSETTINGSCUSTOMAUTHMETHOD_BASIC)),
										path.MatchRelative().AtParent().AtName("method"),
									),
								},
							},

							"auth_token": schema.StringAttribute{
								Description:         providerCustomAuthenticationAuthTokenDescription.Description,
								MarkdownDescription: providerCustomAuthenticationAuthTokenDescription.MarkdownDescription,
								Optional:            true,
								Sensitive:           true,

								Validators: []validator.String{
									stringvalidatorinternal.IsRequiredIfMatchesPathValue(
										types.StringValue(string(management.ENUMNOTIFICATIONSSETTINGSPHONEDELIVERYSETTINGSCUSTOMAUTHMETHOD_BEARER)),
										path.MatchRelative().AtParent().AtName("method"),
									),
								},
							},
						},
					},

					"numbers": schema.SetNestedAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("One or more objects that describe the numbers to use for phone delivery.").Description,
						Optional:    true,

						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"available": schema.BoolAttribute{
									Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether the number is currently available in the provider account.").Description,
									Optional:    true,
									Computed:    true,

									Default: booldefault.StaticBool(false),
								},

								"capabilities": schema.SetAttribute{
									Description:         providerCustomNumbersCapabilitiesDescription.Description,
									MarkdownDescription: providerCustomNumbersCapabilitiesDescription.MarkdownDescription,
									Required:            true,

									ElementType: types.StringType,

									Validators: []validator.Set{
										setvalidator.ValueStringsAre(
											stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumNotificationsSettingsPhoneDeliverySettingsCustomNumbersCapabilityEnumValues)...),
										),
										setvalidator.SizeAtLeast(attrMinLength),
									},
								},

								"number": schema.StringAttribute{
									Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the phone number, toll-free number or short code.").Description,
									Required:    true,
								},

								"selected": schema.BoolAttribute{
									Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether the number is currently available in the provider account.").Description,
									Optional:    true,
									Computed:    true,

									Default: booldefault.StaticBool(false),
								},

								"supported_countries": schema.SetAttribute{
									Description:         providerCustomNumbersSupportedCountriesDescription.Description,
									MarkdownDescription: providerCustomNumbersSupportedCountriesDescription.MarkdownDescription,
									Optional:            true,

									ElementType: types.StringType,

									Validators: []validator.Set{
										setvalidator.Any(
											// Can be set if `type` is `SHORT_CODE` or `TOLL_FREE`, must also be at least one in size and be a 2 letter country code
											setvalidator.All(
												setvalidator.All(
													setvalidator.Any(
														setvalidatorinternal.IsRequiredIfMatchesPathValue(
															types.StringValue(string(management.ENUMNOTIFICATIONSSETTINGSPHONEDELIVERYSETTINGSCUSTOMNUMBERSTYPE_SHORT_CODE)),
															path.MatchRelative().AtParent().AtName("type"),
														),
														setvalidator.SizeAtMost(attrMinLength),
													),
													setvalidatorinternal.IsRequiredIfMatchesPathValue(
														types.StringValue(string(management.ENUMNOTIFICATIONSSETTINGSPHONEDELIVERYSETTINGSCUSTOMNUMBERSTYPE_TOLL_FREE)),
														path.MatchRelative().AtParent().AtName("type"),
													),
												),
												setvalidator.ValueStringsAre(
													stringvalidator.RegexMatches(verify.IsTwoCharCountryCode, "must be a valid two character country code"),
												),
												setvalidator.SizeAtLeast(attrMinLength),
											),

											// Cannot be set if `type` is `PHONE_NUMBER`
											setvalidatorinternal.ConflictsIfMatchesPathValue(
												types.StringValue(string(management.ENUMNOTIFICATIONSSETTINGSPHONEDELIVERYSETTINGSCUSTOMNUMBERSTYPE_PHONE_NUMBER)),
												path.MatchRelative().AtParent().AtName("type"),
											),
										),
									},
								},

								"type": schema.StringAttribute{
									Description:         providerCustomNumbersTypeDescription.Description,
									MarkdownDescription: providerCustomNumbersTypeDescription.MarkdownDescription,
									Required:            true,

									Validators: []validator.String{
										stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumNotificationsSettingsPhoneDeliverySettingsCustomNumbersTypeEnumValues)...),
									},
								},
							},
						},
					},

					"requests": schema.SetNestedAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("One or more objects that describe the outbound custom notification requests.").Description,
						Required:    true,

						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"after_tag": schema.StringAttribute{
									Description:         providerCustomRequestsAfterTagDescription.Description,
									MarkdownDescription: providerCustomRequestsAfterTagDescription.MarkdownDescription,
									Optional:            true,

									Validators: []validator.String{
										stringvalidatorinternal.ConflictsIfMatchesPathValue(
											types.StringValue(string(management.ENUMNOTIFICATIONSSETTINGSPHONEDELIVERYSETTINGSCUSTOMDELIVERYMETHOD_SMS)),
											path.MatchRelative().AtParent().AtName("delivery_method"),
										),
									},
								},

								"before_tag": schema.StringAttribute{
									Description:         providerCustomRequestsBeforeTagDescription.Description,
									MarkdownDescription: providerCustomRequestsBeforeTagDescription.MarkdownDescription,
									Optional:            true,

									Validators: []validator.String{
										stringvalidatorinternal.ConflictsIfMatchesPathValue(
											types.StringValue(string(management.ENUMNOTIFICATIONSSETTINGSPHONEDELIVERYSETTINGSCUSTOMDELIVERYMETHOD_SMS)),
											path.MatchRelative().AtParent().AtName("delivery_method"),
										),
									},
								},

								"body": schema.StringAttribute{
									Description:         providerCustomRequestsBodyDescription.Description,
									MarkdownDescription: providerCustomRequestsBodyDescription.MarkdownDescription,
									Optional:            true,

									Validators: []validator.String{
										stringvalidator.Any(
											stringvalidator.All(
												stringvalidator.RegexMatches(regexp.MustCompile(`\$\{to\}.*\$\{message\}`), "Body must have `${to}` and `${message}` mandatory variables"),
												stringvalidatorinternal.IsRequiredIfMatchesPathValue(
													types.StringValue(string(management.ENUMNOTIFICATIONSSETTINGSPHONEDELIVERYSETTINGSCUSTOMREQUESTMETHOD_POST)),
													path.MatchRelative().AtParent().AtName("method"),
												),
											),
											stringvalidatorinternal.IsRequiredIfMatchesPathValue(
												types.StringValue(string(management.ENUMNOTIFICATIONSSETTINGSPHONEDELIVERYSETTINGSCUSTOMREQUESTMETHOD_GET)),
												path.MatchRelative().AtParent().AtName("method"),
											),
										),
									},
								},

								"delivery_method": schema.StringAttribute{
									Description:         providerCustomRequestsDeliveryMethodDescription.Description,
									MarkdownDescription: providerCustomRequestsDeliveryMethodDescription.MarkdownDescription,
									Required:            true,

									Validators: []validator.String{
										stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumNotificationsSettingsPhoneDeliverySettingsCustomDeliveryMethodEnumValues)...),
									},
								},

								"headers": schema.MapAttribute{
									Description:         providerCustomRequestsHeadersDescription.Description,
									MarkdownDescription: providerCustomRequestsHeadersDescription.MarkdownDescription,
									Optional:            true,

									ElementType: types.StringType,
								},

								"method": schema.StringAttribute{
									Description:         providerCustomRequestsMethodDescription.Description,
									MarkdownDescription: providerCustomRequestsMethodDescription.MarkdownDescription,
									Required:            true,

									Validators: []validator.String{
										stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumNotificationsSettingsPhoneDeliverySettingsCustomRequestMethodEnumValues)...),
									},
								},

								"phone_number_format": schema.StringAttribute{
									Description:         providerCustomRequestsPhoneNumberFormatDescription.Description,
									MarkdownDescription: providerCustomRequestsPhoneNumberFormatDescription.MarkdownDescription,
									Optional:            true,
									Computed:            true,

									Default: stringdefault.StaticString(string(management.ENUMNOTIFICATIONSSETTINGSPHONEDELIVERYSETTINGSCUSTOMNUMBERFORMAT_FULL)),

									Validators: []validator.String{
										stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumNotificationsSettingsPhoneDeliverySettingsCustomNumberFormatEnumValues)...),
									},
								},

								"url": schema.StringAttribute{
									Description:         providerCustomRequestsUrlDescription.Description,
									MarkdownDescription: providerCustomRequestsUrlDescription.MarkdownDescription,
									Required:            true,

									Validators: []validator.String{
										stringvalidator.Any(
											stringvalidator.All(
												stringvalidator.RegexMatches(regexp.MustCompile(`\$\{to\}.*\$\{message\}`), "URL must have `${to}` and `${message}` mandatory variables"),
												stringvalidatorinternal.IsRequiredIfMatchesPathValue(
													types.StringValue(string(management.ENUMNOTIFICATIONSSETTINGSPHONEDELIVERYSETTINGSCUSTOMREQUESTMETHOD_GET)),
													path.MatchRelative().AtParent().AtName("method"),
												),
											),
											stringvalidatorinternal.IsRequiredIfMatchesPathValue(
												types.StringValue(string(management.ENUMNOTIFICATIONSSETTINGSPHONEDELIVERYSETTINGSCUSTOMREQUESTMETHOD_POST)),
												path.MatchRelative().AtParent().AtName("method"),
											),
										),
										stringvalidator.RegexMatches(verify.IsURLWithHTTPS, "URL must be a valid HTTPS URL"),
									},
								},
							},
						},
					},
				},

				Validators: []validator.Object{
					objectvalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName("provider_custom"),
						path.MatchRelative().AtParent().AtName("provider_custom_twilio"),
						path.MatchRelative().AtParent().AtName("provider_custom_syniverse"),
					),
				},
			},

			"provider_custom_twilio": schema.SingleNestedAttribute{
				Description:         providerCustomTwilioDescription.Description,
				MarkdownDescription: providerCustomTwilioDescription.MarkdownDescription,
				Optional:            true,

				Attributes: map[string]schema.Attribute{
					"auth_token": schema.StringAttribute{
						Description:         providerCustomTwilioAuthTokenDescription.Description,
						MarkdownDescription: providerCustomTwilioAuthTokenDescription.MarkdownDescription,
						Required:            true,
						Sensitive:           true,

						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},

					"sid": schema.StringAttribute{
						Description:         providerCustomTwilioSidDescription.Description,
						MarkdownDescription: providerCustomTwilioSidDescription.MarkdownDescription,
						Required:            true,

						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},

					"selected_numbers": schema.SetNestedAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("One or more objects that describe the numbers to use for phone delivery.").Description,
						Required:    true,

						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"number": schema.StringAttribute{
									Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the phone number, toll-free number or short code that has been configured in Twilio.").Description,
									Required:    true,
								},

								"selected": schema.BoolAttribute{
									Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether the number is currently available in the provider account.").Description,
									Computed:    true,

									Default: booldefault.StaticBool(true),
								},

								"supported_countries": schema.SetAttribute{
									Description:         providerCustomNumbersSupportedCountriesDescription.Description,
									MarkdownDescription: providerCustomNumbersSupportedCountriesDescription.MarkdownDescription,
									Optional:            true,

									ElementType: types.StringType,

									Validators: []validator.Set{
										setvalidator.Any(
											// Can be set if `type` is `SHORT_CODE` or `TOLL_FREE`, must also be at least one in size and be a 2 letter country code
											setvalidator.All(
												setvalidator.All(
													setvalidator.Any(
														setvalidatorinternal.IsRequiredIfMatchesPathValue(
															types.StringValue(string(management.ENUMNOTIFICATIONSSETTINGSPHONEDELIVERYSETTINGSCUSTOMNUMBERSTYPE_SHORT_CODE)),
															path.MatchRelative().AtParent().AtName("type"),
														),
														setvalidator.SizeAtMost(attrMinLength),
													),
													setvalidatorinternal.IsRequiredIfMatchesPathValue(
														types.StringValue(string(management.ENUMNOTIFICATIONSSETTINGSPHONEDELIVERYSETTINGSCUSTOMNUMBERSTYPE_TOLL_FREE)),
														path.MatchRelative().AtParent().AtName("type"),
													),
												),
												setvalidator.ValueStringsAre(
													stringvalidator.RegexMatches(verify.IsTwoCharCountryCode, "must be a valid two character country code"),
												),
												setvalidator.SizeAtLeast(attrMinLength),
											),

											// Cannot be set if `type` is `PHONE_NUMBER`
											setvalidatorinternal.ConflictsIfMatchesPathValue(
												types.StringValue(string(management.ENUMNOTIFICATIONSSETTINGSPHONEDELIVERYSETTINGSCUSTOMNUMBERSTYPE_PHONE_NUMBER)),
												path.MatchRelative().AtParent().AtName("type"),
											),
										),
									},
								},

								"type": schema.StringAttribute{
									Description:         providerCustomNumbersTypeDescription.Description,
									MarkdownDescription: providerCustomNumbersTypeDescription.MarkdownDescription,
									Required:            true,

									Validators: []validator.String{
										stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumNotificationsSettingsPhoneDeliverySettingsCustomNumbersTypeEnumValues)...),
									},
								},
							},
						},
					},

					"service_numbers": schema.SetNestedAttribute{
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

									Default: booldefault.StaticBool(true),
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

				Validators: []validator.Object{
					objectvalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName("provider_custom"),
						path.MatchRelative().AtParent().AtName("provider_custom_twilio"),
						path.MatchRelative().AtParent().AtName("provider_custom_syniverse"),
					),
				},
			},

			"provider_custom_syniverse": schema.SingleNestedAttribute{
				Description:         providerCustomSyniverseDescription.Description,
				MarkdownDescription: providerCustomSyniverseDescription.MarkdownDescription,
				Optional:            true,

				Attributes: map[string]schema.Attribute{
					"auth_token": schema.StringAttribute{
						Description:         providerCustomSyniverseAuthTokenDescription.Description,
						MarkdownDescription: providerCustomSyniverseAuthTokenDescription.MarkdownDescription,
						Required:            true,
						Sensitive:           true,

						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},

					"selected_numbers": schema.SetNestedAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("One or more objects that describe the numbers to use for phone delivery.").Description,
						Required:    true,

						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"number": schema.StringAttribute{
									Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the phone number, toll-free number or short code that has been configured in Twilio.").Description,
									Required:    true,
								},

								"selected": schema.BoolAttribute{
									Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that specifies whether the number is currently available in the provider account.").Description,
									Computed:    true,
								},

								"supported_countries": schema.SetAttribute{
									Description:         providerCustomNumbersSupportedCountriesDescription.Description,
									MarkdownDescription: providerCustomNumbersSupportedCountriesDescription.MarkdownDescription,
									Optional:            true,

									ElementType: types.StringType,

									Validators: []validator.Set{
										setvalidator.Any(
											// Can be set if `type` is `SHORT_CODE` or `TOLL_FREE`, must also be at least one in size and be a 2 letter country code
											setvalidator.All(
												setvalidator.All(
													setvalidator.Any(
														setvalidatorinternal.IsRequiredIfMatchesPathValue(
															types.StringValue(string(management.ENUMNOTIFICATIONSSETTINGSPHONEDELIVERYSETTINGSCUSTOMNUMBERSTYPE_SHORT_CODE)),
															path.MatchRelative().AtParent().AtName("type"),
														),
														setvalidator.SizeAtMost(attrMinLength),
													),
													setvalidatorinternal.IsRequiredIfMatchesPathValue(
														types.StringValue(string(management.ENUMNOTIFICATIONSSETTINGSPHONEDELIVERYSETTINGSCUSTOMNUMBERSTYPE_TOLL_FREE)),
														path.MatchRelative().AtParent().AtName("type"),
													),
												),
												setvalidator.ValueStringsAre(
													stringvalidator.RegexMatches(verify.IsTwoCharCountryCode, "must be a valid two character country code"),
												),
												setvalidator.SizeAtLeast(attrMinLength),
											),

											// Cannot be set if `type` is `PHONE_NUMBER`
											setvalidatorinternal.ConflictsIfMatchesPathValue(
												types.StringValue(string(management.ENUMNOTIFICATIONSSETTINGSPHONEDELIVERYSETTINGSCUSTOMNUMBERSTYPE_PHONE_NUMBER)),
												path.MatchRelative().AtParent().AtName("type"),
											),
										),
									},
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

				Validators: []validator.Object{
					objectvalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName("provider_custom"),
						path.MatchRelative().AtParent().AtName("provider_custom_twilio"),
						path.MatchRelative().AtParent().AtName("provider_custom_syniverse"),
					),
				},
			},

			"created_at": schema.StringAttribute{
				Description: "A string that specifies the time the resource was created.",
				Computed:    true,

				CustomType: timetypes.RFC3339Type{},

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"updated_at": schema.StringAttribute{
				Description: "A string that specifies the time the resource was last updated.",
				Computed:    true,

				CustomType: timetypes.RFC3339Type{},
			},
		},
	}
}

func (r *PhoneDeliverySettingsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *PhoneDeliverySettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state PhoneDeliverySettingsResourceModel

	if r.Client.ManagementAPIClient == nil {
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
	phoneDeliverySettings, d := plan.expand(ctx, nil)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var createResponse *management.NotificationsSettingsPhoneDeliverySettings
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.PhoneDeliverySettingsApi.CreatePhoneDeliverySettings(ctx, plan.EnvironmentId.ValueString()).NotificationsSettingsPhoneDeliverySettings(*phoneDeliverySettings).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"CreatePhoneDeliverySettings",
		phoneDeliverySettingsCreateUpdateCustomErrorHandler,
		sdk.DefaultCreateReadRetryable,
		&createResponse,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If twilio or syniverse, then a numbers set will be returned from the services API calls.  We need to merge with configured numbers set.
	var response *management.NotificationsSettingsPhoneDeliverySettings

	if !plan.ProviderCustomTwilio.IsNull() && !plan.ProviderCustomTwilio.IsUnknown() ||
		!plan.ProviderCustomSyniverse.IsNull() && !plan.ProviderCustomSyniverse.IsUnknown() {

		phoneDeliverySettingsId, numbers, d := parsePhoneDeliverySettingsNumbers(createResponse)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}

		phoneDeliverySettingsUpdate, d := plan.expand(ctx, numbers)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}

		resp.Diagnostics.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				fO, fR, fErr := r.Client.ManagementAPIClient.PhoneDeliverySettingsApi.UpdatePhoneDeliverySettings(ctx, plan.EnvironmentId.ValueString(), phoneDeliverySettingsId).NotificationsSettingsPhoneDeliverySettings(*phoneDeliverySettingsUpdate).Execute()
				return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
			},
			"UpdatePhoneDeliverySettings",
			phoneDeliverySettingsCreateUpdateCustomErrorHandler,
			sdk.DefaultCreateReadRetryable,
			&response,
		)...)
		if resp.Diagnostics.HasError() {
			return
		}

	} else {
		response = createResponse
	}

	// Create the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(ctx, response)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *PhoneDeliverySettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *PhoneDeliverySettingsResourceModel

	if r.Client.ManagementAPIClient == nil {
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
	var response *management.NotificationsSettingsPhoneDeliverySettings
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.PhoneDeliverySettingsApi.ReadOnePhoneDeliverySettings(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadOnePhoneDeliverySettings",
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
	resp.Diagnostics.Append(data.toState(ctx, response)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PhoneDeliverySettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state PhoneDeliverySettingsResourceModel

	if r.Client.ManagementAPIClient == nil {
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
	phoneDeliverySettings, d := plan.expand(ctx, nil)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *management.NotificationsSettingsPhoneDeliverySettings
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.PhoneDeliverySettingsApi.UpdatePhoneDeliverySettings(ctx, plan.EnvironmentId.ValueString(), plan.Id.ValueString()).NotificationsSettingsPhoneDeliverySettings(*phoneDeliverySettings).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdatePhoneDeliverySettings",
		phoneDeliverySettingsCreateUpdateCustomErrorHandler,
		sdk.DefaultCreateReadRetryable,
		&response,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(ctx, response)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *PhoneDeliverySettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *PhoneDeliverySettingsResourceModel

	if r.Client.ManagementAPIClient == nil {
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
			fR, fErr := r.Client.ManagementAPIClient.PhoneDeliverySettingsApi.DeletePhoneDeliverySettings(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, fR, fErr)
		},
		"DeletePhoneDeliverySettings",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
		nil,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *PhoneDeliverySettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:     "phone_delivery_settings_id",
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

func phoneDeliverySettingsCreateUpdateCustomErrorHandler(error model.P1Error) diag.Diagnostics {
	var diags diag.Diagnostics

	// Invalid composition
	if details, ok := error.GetDetailsOk(); ok && details != nil && len(details) > 0 {
		if code, ok := details[0].GetCodeOk(); ok && *code == "INVALID_VALUE" {
			diags.AddError(
				"Authentication error",
				fmt.Sprintf("%s. Please check the credentials used to connect to Twilio/Syniverse and retry.", details[0].GetMessage()),
			)

			return diags
		}
	}

	return nil
}

func (p *PhoneDeliverySettingsResourceModel) expand(ctx context.Context, serviceNumbers []management.NotificationsSettingsPhoneDeliverySettingsCustomNumbers) (*management.NotificationsSettingsPhoneDeliverySettings, diag.Diagnostics) {
	var diags diag.Diagnostics

	data := management.NotificationsSettingsPhoneDeliverySettings{
		NotificationsSettingsPhoneDeliverySettingsCustom:          nil,
		NotificationsSettingsPhoneDeliverySettingsTwilioSyniverse: nil,
	}

	if !p.ProviderCustom.IsNull() && !p.ProviderCustom.IsUnknown() {

		if len(serviceNumbers) > 0 {
			diags.AddWarning(
				"Invalid combination of selected/service numbers",
				"The existence of service numbers is not expected for custom provider types.  Please raise an issue with the provider maintainers.",
			)
		}

		var providerPlan PhoneDeliverySettingsProviderCustomResourceModel
		diags.Append(p.ProviderCustom.As(ctx, &providerPlan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		// Expand authentication
		var authenticationPlan PhoneDeliverySettingsProviderCustomAuthenticationResourceModel
		diags.Append(providerPlan.Authentication.As(ctx, &authenticationPlan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		authentication := management.NewNotificationsSettingsPhoneDeliverySettingsCustomAllOfAuthentication(
			management.EnumNotificationsSettingsPhoneDeliverySettingsCustomAuthMethod(authenticationPlan.Method.ValueString()),
		)

		if !authenticationPlan.Password.IsNull() && !authenticationPlan.Password.IsUnknown() {
			authentication.SetPassword(authenticationPlan.Password.ValueString())
		}

		if !authenticationPlan.Username.IsNull() && !authenticationPlan.Username.IsUnknown() {
			authentication.SetUsername(authenticationPlan.Username.ValueString())
		}

		if !authenticationPlan.AuthToken.IsNull() && !authenticationPlan.AuthToken.IsUnknown() {
			authentication.SetAuthToken(authenticationPlan.AuthToken.ValueString())
		}

		// Expand requests
		requests := make([]management.NotificationsSettingsPhoneDeliverySettingsCustomRequest, 0)

		if !providerPlan.Requests.IsNull() && !providerPlan.Requests.IsUnknown() {
			var requestsPlan []PhoneDeliverySettingsProviderCustomRequestsResourceModel
			diags.Append(providerPlan.Requests.ElementsAs(ctx, &requestsPlan, false)...)
			if diags.HasError() {
				return nil, diags
			}

			for _, requestPlan := range requestsPlan {

				request := management.NewNotificationsSettingsPhoneDeliverySettingsCustomRequest(
					management.EnumNotificationsSettingsPhoneDeliverySettingsCustomDeliveryMethod(requestPlan.DeliveryMethod.ValueString()),
					requestPlan.Url.ValueString(),
					management.EnumNotificationsSettingsPhoneDeliverySettingsCustomRequestMethod(requestPlan.Method.ValueString()),
					management.EnumNotificationsSettingsPhoneDeliverySettingsCustomNumberFormat(requestPlan.PhoneNumberFormat.ValueString()),
				)

				if !requestPlan.Body.IsNull() && !requestPlan.Body.IsUnknown() {
					request.SetBody(requestPlan.Body.ValueString())
				}

				if !requestPlan.Headers.IsNull() && !requestPlan.Headers.IsUnknown() {
					var headersPlan map[string]string
					diags.Append(requestPlan.Headers.ElementsAs(ctx, &headersPlan, false)...)
					if diags.HasError() {
						return nil, diags
					}

					request.SetHeaders(headersPlan)
				}

				if !requestPlan.BeforeTag.IsNull() && !requestPlan.BeforeTag.IsUnknown() {
					request.SetBeforeTag(requestPlan.BeforeTag.ValueString())
				}

				if !requestPlan.AfterTag.IsNull() && !requestPlan.AfterTag.IsUnknown() {
					request.SetAfterTag(requestPlan.AfterTag.ValueString())
				}

				requests = append(requests, *request)
			}
		}

		providerData := management.NewNotificationsSettingsPhoneDeliverySettingsCustom(
			management.ENUMNOTIFICATIONSSETTINGSPHONEDELIVERYSETTINGSPROVIDER_PROVIDER,
			providerPlan.Name.ValueString(),
			requests,
			*authentication,
		)

		if !providerPlan.Numbers.IsNull() && !providerPlan.Numbers.IsUnknown() {
			var numbersPlan []PhoneDeliverySettingsProviderCustomNumbersResourceModel
			diags.Append(providerPlan.Numbers.ElementsAs(ctx, &numbersPlan, false)...)
			if diags.HasError() {
				return nil, diags
			}

			numbers, d := phoneDeliverySettingsExpandNumbers(ctx, numbersPlan)
			diags.Append(d...)
			if diags.HasError() {
				return nil, diags
			}

			providerData.SetNumbers(numbers)
		}

		data.NotificationsSettingsPhoneDeliverySettingsCustom = providerData
	}

	if !p.ProviderCustomTwilio.IsNull() && !p.ProviderCustomTwilio.IsUnknown() {
		var providerPlan PhoneDeliverySettingsProviderCustomTwilioResourceModel
		diags.Append(p.ProviderCustomTwilio.As(ctx, &providerPlan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		providerData := management.NewNotificationsSettingsPhoneDeliverySettingsTwilioSyniverse(
			management.ENUMNOTIFICATIONSSETTINGSPHONEDELIVERYSETTINGSPROVIDER_TWILIO,
			providerPlan.Sid.ValueString(),
			providerPlan.AuthToken.ValueString(),
		)

		if !providerPlan.SelectedNumbers.IsNull() && !providerPlan.SelectedNumbers.IsUnknown() && len(serviceNumbers) > 0 {
			var selectedNumbersPlan []PhoneDeliverySettingsProviderCustomSelectedNumbersResourceModel
			diags.Append(providerPlan.SelectedNumbers.ElementsAs(ctx, &selectedNumbersPlan, false)...)
			if diags.HasError() {
				return nil, diags
			}

			numbers := make([]management.NotificationsSettingsPhoneDeliverySettingsCustomNumbers, 0)

			for _, serviceNumber := range serviceNumbers {

				overriddenServiceNumber := serviceNumber

				for _, selectedNumberPlan := range selectedNumbersPlan {
					if serviceNumber.GetNumber() == selectedNumberPlan.Number.ValueString() {

						overriddenServiceNumber.SetSelected(true)

						if !selectedNumberPlan.SupportedCountries.IsNull() && !selectedNumberPlan.SupportedCountries.IsUnknown() {
							var supportedCountries []string
							diags.Append(selectedNumberPlan.SupportedCountries.ElementsAs(ctx, &supportedCountries, false)...)
							if diags.HasError() {
								return nil, diags
							}

							overriddenServiceNumber.SetSupportedCountries(supportedCountries)
						}
					}
				}

				numbers = append(numbers, overriddenServiceNumber)
			}

			providerData.SetNumbers(numbers)
		}

		data.NotificationsSettingsPhoneDeliverySettingsTwilioSyniverse = providerData
	}

	if !p.ProviderCustomSyniverse.IsNull() && !p.ProviderCustomSyniverse.IsUnknown() {
		var providerPlan PhoneDeliverySettingsProviderCustomSyniverseResourceModel
		diags.Append(p.ProviderCustomSyniverse.As(ctx, &providerPlan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		providerData := management.NewNotificationsSettingsPhoneDeliverySettingsTwilioSyniverse(
			management.ENUMNOTIFICATIONSSETTINGSPHONEDELIVERYSETTINGSPROVIDER_SYNIVERSE,
			"",
			providerPlan.AuthToken.ValueString(),
		)

		data.NotificationsSettingsPhoneDeliverySettingsTwilioSyniverse = providerData
	}

	return &data, diags
}

func parsePhoneDeliverySettingsNumbers(apiObject *management.NotificationsSettingsPhoneDeliverySettings) (string, []management.NotificationsSettingsPhoneDeliverySettingsCustomNumbers, diag.Diagnostics) {
	var diags diag.Diagnostics

	apiObjectInstance := apiObject.NotificationsSettingsPhoneDeliverySettingsTwilioSyniverse

	if apiObjectInstance == nil {
		diags.AddError(
			"Invalid phone delivery settings API response",
			"Twilio or Syniverse phone delivery settings must present in the API response.  Please raise this as an issue with the provider maintainers.",
		)
		return "", nil, diags
	}

	return apiObjectInstance.GetId(), apiObjectInstance.GetNumbers(), diags
}

func phoneDeliverySettingsExpandNumbers(ctx context.Context, numbersPlan []PhoneDeliverySettingsProviderCustomNumbersResourceModel) ([]management.NotificationsSettingsPhoneDeliverySettingsCustomNumbers, diag.Diagnostics) {
	var diags diag.Diagnostics

	numbers := make([]management.NotificationsSettingsPhoneDeliverySettingsCustomNumbers, 0)

	for _, numberPlan := range numbersPlan {

		number := management.NewNotificationsSettingsPhoneDeliverySettingsCustomNumbers(
			numberPlan.Number.ValueString(),
			management.EnumNotificationsSettingsPhoneDeliverySettingsCustomNumbersType(numberPlan.Type.ValueString()),
		)

		if !numberPlan.Capabilities.IsNull() && !numberPlan.Capabilities.IsUnknown() {
			var capabilitiesSlice []string
			diags.Append(numberPlan.Capabilities.ElementsAs(ctx, &capabilitiesSlice, false)...)
			if diags.HasError() {
				return nil, diags
			}

			capabilities := make([]management.EnumNotificationsSettingsPhoneDeliverySettingsCustomNumbersCapability, 0)

			for _, capability := range capabilitiesSlice {
				capabilities = append(capabilities, management.EnumNotificationsSettingsPhoneDeliverySettingsCustomNumbersCapability(capability))
			}
			number.SetCapabilities(capabilities)
		}

		if !numberPlan.Selected.IsNull() && !numberPlan.Selected.IsUnknown() {
			number.SetSelected(numberPlan.Selected.ValueBool())
		}

		if !numberPlan.Available.IsNull() && !numberPlan.Available.IsUnknown() {
			number.SetAvailable(numberPlan.Available.ValueBool())
		}

		if !numberPlan.SupportedCountries.IsNull() && !numberPlan.SupportedCountries.IsUnknown() {
			var supportedCountries []string
			diags.Append(numberPlan.SupportedCountries.ElementsAs(ctx, &supportedCountries, false)...)
			if diags.HasError() {
				return nil, diags
			}

			number.SetSupportedCountries(supportedCountries)
		}

		numbers = append(numbers, *number)
	}

	return numbers, diags
}

func (p *PhoneDeliverySettingsResourceModel) toState(ctx context.Context, apiObject *management.NotificationsSettingsPhoneDeliverySettings) diag.Diagnostics {
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
	p.ProviderType = framework.EnumOkToTF(apiObjectCommon.GetProviderOk())
	p.CreatedAt = framework.TimeOkToTF(apiObjectCommon.GetCreatedAtOk())
	p.UpdatedAt = framework.TimeOkToTF(apiObjectCommon.GetUpdatedAtOk())

	var d diag.Diagnostics

	if p.ProviderType.Equal(types.StringValue(string(management.ENUMNOTIFICATIONSSETTINGSPHONEDELIVERYSETTINGSPROVIDER_PROVIDER))) {
		var providerPlan *PhoneDeliverySettingsProviderCustomResourceModel
		diags.Append(p.ProviderCustom.As(ctx, &providerPlan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return diags
		}

		p.ProviderCustom, d = p.toStatePhoneDeliverySettingsProviderCustom(ctx, providerPlan, apiObject.NotificationsSettingsPhoneDeliverySettingsCustom)
		diags.Append(d...)
	} else {
		p.ProviderCustom = types.ObjectNull(customTFObjectTypes)
	}

	if p.ProviderType.Equal(types.StringValue(string(management.ENUMNOTIFICATIONSSETTINGSPHONEDELIVERYSETTINGSPROVIDER_TWILIO))) {
		var providerPlan *PhoneDeliverySettingsProviderCustomTwilioResourceModel
		diags.Append(p.ProviderCustomTwilio.As(ctx, &providerPlan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return diags
		}

		p.ProviderCustomTwilio, d = p.toStatePhoneDeliverySettingsProviderCustomTwilio(ctx, providerPlan, apiObject.NotificationsSettingsPhoneDeliverySettingsTwilioSyniverse)
		diags.Append(d...)
	} else {
		p.ProviderCustomTwilio = types.ObjectNull(twilioTFObjectTypes)
	}

	if p.ProviderType.Equal(types.StringValue(string(management.ENUMNOTIFICATIONSSETTINGSPHONEDELIVERYSETTINGSPROVIDER_SYNIVERSE))) {
		var providerPlan *PhoneDeliverySettingsProviderCustomSyniverseResourceModel
		diags.Append(p.ProviderCustomSyniverse.As(ctx, &providerPlan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return diags
		}

		p.ProviderCustomSyniverse, d = p.toStatePhoneDeliverySettingsProviderCustomSyniverse(ctx, providerPlan, apiObject.NotificationsSettingsPhoneDeliverySettingsTwilioSyniverse)
		diags.Append(d...)
	} else {
		p.ProviderCustomSyniverse = types.ObjectNull(syniverseTFObjectTypes)
	}

	return diags
}

func (p *PhoneDeliverySettingsResourceModel) toStatePhoneDeliverySettingsProviderCustom(ctx context.Context, planData *PhoneDeliverySettingsProviderCustomResourceModel, apiObject *management.NotificationsSettingsPhoneDeliverySettingsCustom) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if apiObject == nil || apiObject.GetId() == "" {
		return types.ObjectNull(customTFObjectTypes), diags
	}

	objMap := map[string]attr.Value{
		"name": framework.StringOkToTF(apiObject.GetNameOk()),
	}

	var d diag.Diagnostics

	var authenticationPlan *PhoneDeliverySettingsProviderCustomAuthenticationResourceModel

	if planData != nil {
		diags.Append(planData.Authentication.As(ctx, &authenticationPlan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return types.ObjectNull(customTFObjectTypes), diags
		}
	}

	authentication, ok := apiObject.GetAuthenticationOk()
	objMap["authentication"], d = phoneDeliverySettingsCustomAuthenticationOkToTF(authenticationPlan, authentication, ok)
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(customTFObjectTypes), diags
	}

	objMap["numbers"], d = phoneDeliverySettingsCustomNumbersOkToTF(apiObject.GetNumbersOk())
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(customTFObjectTypes), diags
	}

	objMap["requests"], d = phoneDeliverySettingsCustomRequestsOkToTF(apiObject.GetRequestsOk())
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(customTFObjectTypes), diags
	}

	objValue, d := types.ObjectValue(customTFObjectTypes, objMap)
	diags.Append(d...)

	return objValue, diags
}

func phoneDeliverySettingsCustomAuthenticationOkToTF(planData *PhoneDeliverySettingsProviderCustomAuthenticationResourceModel, apiObject *management.NotificationsSettingsPhoneDeliverySettingsCustomAllOfAuthentication, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(customAuthenticationTFObjectTypes), diags
	}

	objMap := map[string]attr.Value{
		"method":     framework.EnumOkToTF(apiObject.GetMethodOk()),
		"password":   types.StringNull(),
		"auth_token": types.StringNull(),
		"username":   framework.StringOkToTF(apiObject.GetUsernameOk()),
	}

	if planData != nil {
		objMap["password"] = planData.Password
		objMap["auth_token"] = planData.AuthToken
	}

	returnVar, d := types.ObjectValue(customAuthenticationTFObjectTypes, objMap)
	diags.Append(d...)

	return returnVar, diags
}

func phoneDeliverySettingsCustomNumbersOkToTF(apiObject []management.NotificationsSettingsPhoneDeliverySettingsCustomNumbers, ok bool) (basetypes.SetValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	tfObjType := types.ObjectType{AttrTypes: customNumbersTFObjectTypes}

	if !ok || len(apiObject) == 0 {
		return types.SetNull(tfObjType), diags
	}

	flattenedList := []attr.Value{}
	for _, v := range apiObject {

		objMap := map[string]attr.Value{
			"supported_countries": framework.StringSetOkToTF(v.GetSupportedCountriesOk()),
			"type":                framework.EnumOkToTF(v.GetTypeOk()),
			"selected":            framework.BoolOkToTF(v.GetSelectedOk()),
			"available":           framework.BoolOkToTF(v.GetAvailableOk()),
			"number":              framework.StringOkToTF(v.GetNumberOk()),
			"capabilities":        framework.EnumSetOkToTF(v.GetCapabilitiesOk()),
		}

		flattenedObj, d := types.ObjectValue(customNumbersTFObjectTypes, objMap)
		diags.Append(d...)

		flattenedList = append(flattenedList, flattenedObj)
	}

	returnVar, d := types.SetValue(tfObjType, flattenedList)
	diags.Append(d...)

	return returnVar, diags
}

func phoneDeliverySettingsCustomRequestsOkToTF(apiObject []management.NotificationsSettingsPhoneDeliverySettingsCustomRequest, ok bool) (basetypes.SetValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	tfObjType := types.ObjectType{AttrTypes: customRequestsTFObjectTypes}

	if !ok || len(apiObject) == 0 {
		return types.SetNull(tfObjType), diags
	}

	flattenedList := []attr.Value{}
	for _, v := range apiObject {

		objMap := map[string]attr.Value{
			"delivery_method":     framework.EnumOkToTF(v.GetDeliveryMethodOk()),
			"url":                 framework.StringOkToTF(v.GetUrlOk()),
			"method":              framework.EnumOkToTF(v.GetMethodOk()),
			"body":                framework.StringOkToTF(v.GetBodyOk()),
			"headers":             framework.StringMapOkToTF(v.GetHeadersOk()),
			"before_tag":          framework.StringOkToTF(v.GetBeforeTagOk()),
			"after_tag":           framework.StringOkToTF(v.GetAfterTagOk()),
			"phone_number_format": framework.EnumOkToTF(v.GetPhoneNumberFormatOk()),
		}

		flattenedObj, d := types.ObjectValue(customRequestsTFObjectTypes, objMap)
		diags.Append(d...)

		flattenedList = append(flattenedList, flattenedObj)
	}

	returnVar, d := types.SetValue(tfObjType, flattenedList)
	diags.Append(d...)

	return returnVar, diags
}

func (p *PhoneDeliverySettingsResourceModel) toStatePhoneDeliverySettingsProviderCustomTwilio(ctx context.Context, planData *PhoneDeliverySettingsProviderCustomTwilioResourceModel, apiObject *management.NotificationsSettingsPhoneDeliverySettingsTwilioSyniverse) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if apiObject == nil || apiObject.GetId() == "" {
		return types.ObjectNull(twilioTFObjectTypes), diags
	}

	objMap := map[string]attr.Value{
		"sid":        framework.StringOkToTF(apiObject.GetSidOk()),
		"auth_token": types.StringNull(),
	}

	if planData != nil {
		objMap["auth_token"] = planData.AuthToken
	}

	var d diag.Diagnostics

	numbers, ok := apiObject.GetNumbersOk()
	objMap["service_numbers"], d = phoneDeliverySettingsTwilioSyniverseNumbersOkToTF(numbers, ok)
	diags.Append(d...)

	selectedNumbers := types.SetNull(types.ObjectType{AttrTypes: customSelectedNumbersTFObjectTypes})

	if planData != nil {
		selectedNumbers = planData.SelectedNumbers
	}

	objMap["selected_numbers"], d = phoneDeliverySettingsTwilioSyniverseSelectedNumbersOkToTF(ctx, selectedNumbers, numbers, ok)
	diags.Append(d...)

	objValue, d := types.ObjectValue(twilioTFObjectTypes, objMap)
	diags.Append(d...)

	return objValue, diags
}

func (p *PhoneDeliverySettingsResourceModel) toStatePhoneDeliverySettingsProviderCustomSyniverse(ctx context.Context, planData *PhoneDeliverySettingsProviderCustomSyniverseResourceModel, apiObject *management.NotificationsSettingsPhoneDeliverySettingsTwilioSyniverse) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if apiObject == nil || apiObject.GetId() == "" {
		return types.ObjectNull(syniverseTFObjectTypes), diags
	}

	objMap := map[string]attr.Value{
		"auth_token": types.StringNull(),
	}

	if planData != nil {
		objMap["auth_token"] = planData.AuthToken
	}

	var d diag.Diagnostics

	numbers, ok := apiObject.GetNumbersOk()
	objMap["service_numbers"], d = phoneDeliverySettingsTwilioSyniverseNumbersOkToTF(numbers, ok)
	diags.Append(d...)

	selectedNumbers := types.SetNull(types.ObjectType{AttrTypes: customSelectedNumbersTFObjectTypes})

	if planData != nil {
		selectedNumbers = planData.SelectedNumbers
	}

	objMap["selected_numbers"], d = phoneDeliverySettingsTwilioSyniverseSelectedNumbersOkToTF(ctx, selectedNumbers, numbers, ok)
	diags.Append(d...)

	objValue, d := types.ObjectValue(syniverseTFObjectTypes, objMap)
	diags.Append(d...)

	return objValue, diags
}

func phoneDeliverySettingsTwilioSyniverseNumbersOkToTF(apiObject []management.NotificationsSettingsPhoneDeliverySettingsCustomNumbers, ok bool) (basetypes.SetValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	tfObjType := types.ObjectType{AttrTypes: customNumbersTFObjectTypes}

	if !ok || len(apiObject) == 0 {
		return types.SetNull(tfObjType), diags
	}

	flattenedNumbersList := []attr.Value{}
	for _, v := range apiObject {

		if vNumber, ok := v.GetNumberOk(); ok {

			objMap := map[string]attr.Value{
				"supported_countries": framework.StringSetOkToTF(v.GetSupportedCountriesOk()),
				"type":                framework.EnumOkToTF(v.GetTypeOk()),
				"selected":            framework.BoolOkToTF(v.GetSelectedOk()),
				"available":           framework.BoolOkToTF(v.GetAvailableOk()),
				"number":              framework.StringToTF(*vNumber),
				"capabilities":        framework.EnumSetOkToTF(v.GetCapabilitiesOk()),
			}

			flattenedNumberObj, d := types.ObjectValue(customNumbersTFObjectTypes, objMap)
			diags.Append(d...)

			flattenedNumbersList = append(flattenedNumbersList, flattenedNumberObj)
		}
	}

	returnVarNumbers, d := types.SetValue(tfObjType, flattenedNumbersList)
	diags.Append(d...)

	return returnVarNumbers, diags
}

func phoneDeliverySettingsTwilioSyniverseSelectedNumbersOkToTF(ctx context.Context, plan basetypes.SetValue, apiObject []management.NotificationsSettingsPhoneDeliverySettingsCustomNumbers, ok bool) (basetypes.SetValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	tfObjType := types.ObjectType{AttrTypes: customSelectedNumbersTFObjectTypes}

	if !ok || len(apiObject) == 0 || plan.IsNull() || plan.IsUnknown() {
		return types.SetNull(tfObjType), diags
	}

	// Get the list of numbers
	selectedNumbers := make([]string, 0)
	if !plan.IsNull() && !plan.IsUnknown() {
		var numbersPlan []PhoneDeliverySettingsProviderCustomSelectedNumbersResourceModel
		diags.Append(plan.ElementsAs(ctx, &numbersPlan, false)...)
		if diags.HasError() {
			return types.SetNull(tfObjType), diags
		}

		for _, v := range numbersPlan {
			selectedNumbers = append(selectedNumbers, v.Number.ValueString())
		}
	}

	flattenedList := []attr.Value{}
	for _, v := range apiObject {

		if vNumber, ok := v.GetNumberOk(); ok {
			if slices.Contains(selectedNumbers, *vNumber) {

				selectedObjMap := map[string]attr.Value{
					"supported_countries": framework.StringSetOkToTF(v.GetSupportedCountriesOk()),
					"type":                framework.EnumOkToTF(v.GetTypeOk()),
					"selected":            framework.BoolOkToTF(v.GetSelectedOk()),
					"number":              framework.StringToTF(*vNumber),
				}

				flattenedObj, d := types.ObjectValue(customSelectedNumbersTFObjectTypes, selectedObjMap)
				diags.Append(d...)

				flattenedList = append(flattenedList, flattenedObj)
			}
		}
	}

	returnVar, d := types.SetValue(tfObjType, flattenedList)
	diags.Append(d...)

	return returnVar, diags
}
