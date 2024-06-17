package sso

import (
	"context"
	"fmt"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/mfa"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type UserResource serviceClientType

type UserResourceModel struct {
	Id            pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	Username      types.String                 `tfsdk:"username"`
	Email         types.String                 `tfsdk:"email"`
	//EmailVerified     types.Bool   `tfsdk:"email_verified"`
	Enabled           types.Bool                   `tfsdk:"enabled"`
	PopulationId      pingonetypes.ResourceIDValue `tfsdk:"population_id"`
	Account           types.Object                 `tfsdk:"account"`
	Address           types.Object                 `tfsdk:"address"`
	ExternalId        types.String                 `tfsdk:"external_id"`
	IdentityProvider  types.Object                 `tfsdk:"identity_provider"`
	Lifecycle         types.Object                 `tfsdk:"user_lifecycle"`
	Locale            types.String                 `tfsdk:"locale"`
	MFAEnabled        types.Bool                   `tfsdk:"mfa_enabled"`
	MobilePhone       types.String                 `tfsdk:"mobile_phone"`
	Name              types.Object                 `tfsdk:"name"`
	Nickname          types.String                 `tfsdk:"nickname"`
	Password          types.Object                 `tfsdk:"password"`
	Photo             types.Object                 `tfsdk:"photo"`
	PreferredLanguage types.String                 `tfsdk:"preferred_language"`
	PrimaryPhone      types.String                 `tfsdk:"primary_phone"`
	Timezone          types.String                 `tfsdk:"timezone"`
	Title             types.String                 `tfsdk:"title"`
	Type              types.String                 `tfsdk:"type"`
	VerifyStatus      types.String                 `tfsdk:"verify_status"`
}

type UserAccountResourceModel struct {
	CanAuthenticate types.Bool        `tfsdk:"can_authenticate"`
	LockedAt        timetypes.RFC3339 `tfsdk:"locked_at"`
	Status          types.String      `tfsdk:"status"`
}

type UserAddressResourceModel struct {
	CountryCode   types.String `tfsdk:"country_code"`
	Locality      types.String `tfsdk:"locality"`
	PostalCode    types.String `tfsdk:"postal_code"`
	Region        types.String `tfsdk:"region"`
	StreetAddress types.String `tfsdk:"street_address"`
}

type UserIdentityProviderResourceModel struct {
	Id   pingonetypes.ResourceIDValue `tfsdk:"id"`
	Type types.String                 `tfsdk:"type"`
}

type UserLifecycleResourceModel struct {
	Status                   types.String `tfsdk:"status"`
	SuppressVerificationCode types.Bool   `tfsdk:"suppress_verification_code"`
}

type UserNameResourceModel struct {
	Family          types.String `tfsdk:"family"`
	Formatted       types.String `tfsdk:"formatted"`
	Given           types.String `tfsdk:"given"`
	HonorificPrefix types.String `tfsdk:"honorific_prefix"`
	HonorificSuffix types.String `tfsdk:"honorific_suffix"`
	Middle          types.String `tfsdk:"middle"`
}

type UserPasswordResourceModel struct {
	ForceChange  types.Bool   `tfsdk:"force_change"`
	InitialValue types.String `tfsdk:"initial_value"`
	External     types.Object `tfsdk:"external"`
}

type UserPasswordExternalResourceModel struct {
	Gateway types.Object `tfsdk:"gateway"`
}

type UserPasswordExternalGatewayResourceModel struct {
	Id                    pingonetypes.ResourceIDValue `tfsdk:"id"`
	Type                  types.String                 `tfsdk:"type"`
	UserTypeId            pingonetypes.ResourceIDValue `tfsdk:"user_type_id"`
	CorrelationAttributes types.Map                    `tfsdk:"correlation_attributes"`
}

type UserPhotoResourceModel struct {
	Href types.String `tfsdk:"href"`
}

var (
	userAccountTFObjectTypes = map[string]attr.Type{
		"can_authenticate": types.BoolType,
		"locked_at":        timetypes.RFC3339Type{},
		"status":           types.StringType,
	}

	userAddressTFObjectTypes = map[string]attr.Type{
		"country_code":   types.StringType,
		"locality":       types.StringType,
		"postal_code":    types.StringType,
		"region":         types.StringType,
		"street_address": types.StringType,
	}

	userIdentityProviderTFObjectTypes = map[string]attr.Type{
		"id":   pingonetypes.ResourceIDType{},
		"type": types.StringType,
	}

	userLifecycleTFObjectTypes = map[string]attr.Type{
		"status":                     types.StringType,
		"suppress_verification_code": types.BoolType,
	}

	userNameTFObjectTypes = map[string]attr.Type{
		"family":           types.StringType,
		"formatted":        types.StringType,
		"given":            types.StringType,
		"honorific_prefix": types.StringType,
		"honorific_suffix": types.StringType,
		"middle":           types.StringType,
	}

	userPasswordTFObjectTypes = map[string]attr.Type{
		"force_change":  types.BoolType,
		"initial_value": types.StringType,
		"external": types.ObjectType{
			AttrTypes: userPasswordExternalTFObjectTypes,
		},
	}

	userPasswordExternalTFObjectTypes = map[string]attr.Type{
		"gateway": types.ObjectType{
			AttrTypes: userPasswordExternalGatewayTFObjectTypes,
		},
	}

	userPasswordExternalGatewayTFObjectTypes = map[string]attr.Type{
		"id":           pingonetypes.ResourceIDType{},
		"type":         types.StringType,
		"user_type_id": pingonetypes.ResourceIDType{},
		"correlation_attributes": types.MapType{
			ElemType: types.StringType,
		},
	}

	userPhotoTFObjectTypes = map[string]attr.Type{
		"href": types.StringType,
	}
)

// Framework interfaces
var (
	_ resource.Resource                = &UserResource{}
	_ resource.ResourceWithConfigure   = &UserResource{}
	_ resource.ResourceWithImportState = &UserResource{}
	_ resource.ResourceWithModifyPlan  = &UserResource{}
)

// New Object
func NewUserResource() resource.Resource {
	return &UserResource{}
}

// Metadata
func (r *UserResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

// Schema.
func (r *UserResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	const attrMinLength = 1

	const attrUsernameMaxLength = 128
	const attrAddressLocalityMaxLength = 256
	const attrAddressPostalCodeMaxLength = 40
	const attrAddressRegionMaxLength = 256
	const attrAddressStreetAddressMaxLength = 256
	const attrExternalIdMaxLength = 1024
	const attrLocaleMaxLength = 256
	const attrPhoneMaxLength = 32
	const attrNameFamilyMaxLength = 256
	const attrNameFormattedMaxLength = 256
	const attrNameGivenMaxLength = 256
	const attrNameMiddleMaxLength = 256
	const attrNicknameMaxLength = 256
	const attrTitleMaxLength = 256
	const attrTypeMaxLength = 256

	usernameDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the user name, which must be provided and must be unique within an environment. The `username` must either be a well-formed email address or a string. The string can contain any letters, numbers, combining characters, math and currency symbols, dingbats and drawing characters, and invisible whitespace",
	)

	emailDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the user's email address, which must be provided and valid. For more information about email address formatting, see section 3.4 of [RFC 2822, Internet Message Format](http://www.faqs.org/rfcs/rfc2822.html).",
	)

	enabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether the user is enabled. This attribute is set to `true` by default when the user is created.",
	)

	accountCanAuthenticateDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether the user can authenticate. If the value is set to `false`, the account is locked or the user is disabled, and unless specified otherwise in administrative configuration, the user will be unable to authenticate.",
	)

	accountStatusDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the the account locked state.",
	).AllowedValuesEnum(management.AllowedEnumUserStatusEnumValues)

	addressCountryCodeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the country name component. When specified, the value must be in [ISO 3166-1](https://www.iso.org/iso-3166-country-codes.html) \"alpha-2\" code format. For example, the country codes for the United States and Sweden are `US` and `SE`, respectively. Valid characters consist of two upper-case letters.",
	)

	identityProviderIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that identifies the external identity provider used to authenticate the user. If not provided, PingOne is the identity provider. This attribute is required if the identity provider is authoritative for just-in-time user provisioning.",
	).RequiresReplace()

	identityProviderTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the type of identity provider used to authenticate the user.",
	).AllowedValuesEnum(management.AllowedEnumIdentityProviderEnumValues).AppendMarkdownString(
		"The default value of `PING_ONE` is set when a value for `id` is not provided in this object.",
	)

	userLifecycleDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that specifies the user's identity lifecycle information.",
	).RequiresReplace()

	userLifecycleStatusDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the status of the account lifecycle.",
	).AllowedValuesEnum(management.AllowedEnumUserLifecycleStatusEnumValues).AppendMarkdownString(
		" This property value is only allowed to be set when importing a user to set the initial account status. If the initial status is set to `VERIFICATION_REQUIRED` and an email address is provided, a verification email is sent.",
	).RequiresReplace()

	userLifecycleSuppressVerificationCodeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether to suppress the verification code when the user is imported and the `status` is set to `VERIFICATION_REQUIRED`. If this property is set to `true`, no verification email is sent to the user. If this property is omitted or set to `false`, a verification email is sent automatically to the user.",
	).RequiresReplace()

	localeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the user's default location. This may be explicitly set to null when updating a user to unset it. This is used for purposes of localizing such items as currency, date time format, or numerical representations. If provided, it must be a valid language tag as defined in [RFC 5646](https://www.rfc-editor.org/rfc/rfc5646.html). The following are example tags: `fr`, `en-US`, `es-419`, `az-Arab`, `man-Nkoo-GN`. The string can contain any letters, numbers, combining characters, math and currency symbols, dingbats and drawing characters, and invisible whitespace. It can have a length of no more than 256 characters.",
	)

	mfaEnabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether multi-factor authentication is enabled.",
	).DefaultValue(false)

	mobilePhoneDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the user's native phone number. This might also match the `primary_phone` attribute. This may be explicitly set to null when updating a user to unset it. Valid phone numbers must have at least one number and a maximum character length of 32.",
	)

	nameFamilyDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the family name of the user, or Last in most Western languages (for example, `Jensen` given the full name `Ms. Barbara J Jensen, III`). This may be explicitly set to null when updating a name to unset it. Valid characters consist of any Unicode letter, mark (for example, accent, umlaut), space, dot, apostrophe, or hyphen. It can have a length of no more than 256 characters.",
	)

	nameFormattedDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the fully formatted name of the user (for example `Ms. Barbara J Jensen, III`). This can be explicitly set to null when updating a name to unset it. Valid characters consist of any Unicode letter, mark (for example, accent, umlaut), space, dot, apostrophe, or hyphen. It can have a length of no more than 256 characters.",
	)

	nameGivenDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the given name of the user, or First in most Western languages (for example, `Barbara` given the full name `Ms. Barbara J Jensen, III`). This may be explicitly set to null when updating a name to unset it. The string can contain any letters, numbers, combining characters, math and currency symbols, dingbats and drawing characters, and invisible whitespace. It can have a length of no more than 256 characters.",
	)

	nameHonorificPrefixDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the honorific prefix(es) of the user, or title in most Western languages (for example, `Ms.` given the full name `Ms. Barbara Jane Jensen, III`). This can be explicitly set to null when updating a name to unset it.",
	)

	nameHonorificSuffixDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the honorific suffix(es) of the user, or suffix in most Western languages (for example, `III` given the full name `Ms. Barbara Jane Jensen, III`). This can be explicitly set to null when updating a name to unset it.",
	)

	nameMiddleDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the middle name(s) of the user (for exmple, `Jane` given the full name `Ms. Barbara Jane Jensen, III`). This can be explicitly set to null when updating a name to unset it. The string can contain any letters, numbers, combining characters, math and currency symbols, dingbats and drawing characters, and invisible whitespace. It can have a length of no more than 256 characters.",
	)

	nicknameDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the user's nickname. This can be explicitly set to null when updating a user to unset it. The string can contain any letters, numbers, combining characters, math and currency symbols, dingbats and drawing characters, and invisible whitespace. It can have a length of no more than 256 characters.",
	)

	passwordForceChangeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether the user is forced to change the password on the next log in.",
	).DefaultValue("false")

	passwordExternalGatewayTypeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that indicates one of the supported gateway types.",
	).AllowedValuesEnum(management.AllowedEnumGatewayTypeEnumValues)

	photoHrefDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The URI that is a uniform resource locator (as defined in [Section 1.1.3 of RFC 3986](https://www.rfc-editor.org/rfc/rfc3986#section-1.3)) that points to a resource location representing the user's image. This can be removed from a user by setting the photo attribute to null. If provided, the resource must be a file (for example, a GIF, JPEG, or PNG image file) rather than a web page containing an image. It must be a valid URL that starts with the HTTP or HTTPS scheme.",
	)

	preferredLanguageDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the user's preferred written or spoken languages. This may be explicitly set to null when updating a user to unset it. If provided, the format of the value must be a valid language range and the same as the HTTP `Accept-Language` header field (not including `Accept-Language:` prefix) and is specified in [Section 5.3.5 of RFC 7231](https://datatracker.ietf.org/doc/html/rfc7231#section-5.3.5). For example: `en-US`, `en-gb;q=0.8`, `en;q=0.7`.",
	)

	primaryPhoneDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the user's primary phone number. This might also match the `mobile_phone` attribute. This may be explicitly set to null when updating a user to unset it. Valid phone numbers must have at least one number and a maximum character length of 32.",
	)

	timezoneDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the user's time zone. This can be explicitly set to null when updating a user to unset it. If provided, it must conform with the IANA Time Zone database format [RFC 6557](https://www.rfc-editor.org/rfc/rfc6557.html), also known as the \"Olson\" time zone database format [Olson-TZ](https://www.iana.org/time-zones) for example, `America/Los_Angeles`.",
	)

	titleDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the user's title, such as `Vice President`. This can be explicitly set to null when updating a user to unset it. The string can contain any letters, numbers, combining characters, math and currency symbols, dingbats and drawing characters, and invisible whitespace. It can have a length of no more than 256 characters.",
	)

	typeDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the user's type, which is optional. This can be explicitly set to null when updating a user to unset it. This attribute is organization-specific and has no special meaning within the PingOne platform. It is a free-text field that could have values of (for example) `Contractor`, `Employee`, `Intern`, `Temp`, `External`, and `Unknown`. The string can contain any letters, numbers, combining characters, math and currency symbols, dingbats and drawing characters, and invisible whitespace. It can have a length of no more than 256 characters.",
	)

	verifyStatusDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that indicates whether ID verification can be done for the user.",
	).AllowedValuesEnum(management.AllowedEnumUserVerifyStatusEnumValues).AppendMarkdownString(
		"If the user verification status is `DISABLED`, a new verification status cannot be created for that user until the status is changed to `ENABLED`.",
	).RequiresReplace()

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage a PingOne user in an environment.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to manage the user in."),
			),

			"username": schema.StringAttribute{
				Description:         usernameDescription.Description,
				MarkdownDescription: usernameDescription.MarkdownDescription,
				Required:            true,

				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
					stringvalidator.LengthAtMost(attrUsernameMaxLength),
					stringvalidator.RegexMatches(regexp.MustCompile(`^[\p{L}\p{M}\p{Zs}\p{S}\p{N}\p{P}]*$`), `must match the regex ^[\p{L}\p{M}\p{Zs}\p{S}\p{N}\p{P}]*$`),
				},
			},

			"email": schema.StringAttribute{
				Description:         emailDescription.Description,
				MarkdownDescription: emailDescription.MarkdownDescription,
				Required:            true,

				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			// "email_verified": schema.BoolAttribute{
			// 	Description: framework.SchemaAttributeDescriptionFromMarkdown(
			// 		"A boolean that specifies whether the user's email is verified. An email address can be verified during account verification. If the email address used to request the verification code is the same as the user,s email at verification time (and the verification code is valid), then the email is verified. The value of this property can be set on user import.",
			// 	).Description,
			// 	Computed: true,

			// 	PlanModifiers: []planmodifier.Bool{
			// 		boolplanmodifier.UseStateForUnknown(),
			// 	},
			// },

			"enabled": schema.BoolAttribute{
				Description:         enabledDescription.Description,
				MarkdownDescription: enabledDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,

				Default: booldefault.StaticBool(true),
			},

			"population_id": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown(
					"The identifier of the population resource associated with the user.  Must be a valid PingOne resource ID.",
				).Description,
				Required: true,

				CustomType: pingonetypes.ResourceIDType{},
			},

			"account": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown(
					"A single object that specifies the user's account information.",
				).Description,
				Optional: true,
				Computed: true,

				Attributes: map[string]schema.Attribute{
					"can_authenticate": schema.BoolAttribute{
						Description:         accountCanAuthenticateDescription.Description,
						MarkdownDescription: accountCanAuthenticateDescription.MarkdownDescription,
						Optional:            true,
						Computed:            true,
					},

					"locked_at": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown(
							"The time the specified user account was locked. This property might be absent if the account is unlocked or if the account was locked out automatically by failed password attempts.",
						).Description,
						Computed: true,

						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},

						CustomType: timetypes.RFC3339Type{},
					},

					"status": schema.StringAttribute{
						Description:         accountStatusDescription.Description,
						MarkdownDescription: accountStatusDescription.MarkdownDescription,
						Required:            true,

						Validators: []validator.String{
							stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumUserStatusEnumValues)...),
						},
					},
				},
			},

			"address": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown(
					"A single object that specifies the user's address information.",
				).Description,
				Optional: true,

				Attributes: map[string]schema.Attribute{
					"country_code": schema.StringAttribute{
						Description:         addressCountryCodeDescription.Description,
						MarkdownDescription: addressCountryCodeDescription.MarkdownDescription,
						Optional:            true,

						Validators: []validator.String{
							stringvalidator.RegexMatches(verify.IsTwoCharCountryCode, `must be a valid two character country code`),
						},
					},

					"locality": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown(
							"A string that specifies the city or locality component of the address. The string can contain any letters, numbers, combining characters, math and currency symbols, dingbats and drawing characters, and invisible whitespace. It can have a length of no more than 256 characters.",
						).Description,
						Optional: true,

						Validators: []validator.String{
							stringvalidator.LengthAtLeast(attrMinLength),
							stringvalidator.LengthAtMost(attrAddressLocalityMaxLength),
							stringvalidator.RegexMatches(regexp.MustCompile(`^[\p{L}\p{M}\p{Zs}\p{S}\p{N}\p{P}]*$`), `must match the regex ^[\p{L}\p{M}\p{Zs}\p{S}\p{N}\p{P}]*$`),
						},
					},

					"postal_code": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown(
							"A string that specifies the ZIP code or postal code component of the address. The string can contain any letters, numbers, combining characters, math and currency symbols, dingbats and drawing characters, and invisible whitespace. It can have a length of no more than 40 characters.",
						).Description,
						Optional: true,

						Validators: []validator.String{
							stringvalidator.LengthAtLeast(attrMinLength),
							stringvalidator.LengthAtMost(attrAddressPostalCodeMaxLength),
							stringvalidator.RegexMatches(regexp.MustCompile(`^[\p{L}\p{M}\p{Zs}\p{S}\p{N}\p{P}]*$`), `must match the regex ^[\p{L}\p{M}\p{Zs}\p{S}\p{N}\p{P}]*$`),
						},
					},

					"region": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown(
							"A string that specifies the state, province, or region component of the address. The string can contain any letters, numbers, combining characters, math and currency symbols, dingbats and drawing characters, and invisible whitespace. It can have a length of no more than 256 characters.",
						).Description,
						Optional: true,

						Validators: []validator.String{
							stringvalidator.LengthAtLeast(attrMinLength),
							stringvalidator.LengthAtMost(attrAddressRegionMaxLength),
							stringvalidator.RegexMatches(regexp.MustCompile(`^[\p{L}\p{M}\p{Zs}\p{S}\p{N}\p{P}]*$`), `must match the regex ^[\p{L}\p{M}\p{Zs}\p{S}\p{N}\p{P}]*$`),
						},
					},

					"street_address": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown(
							"A string that specifies the full street address component, which may include house number, street name, P.O. box, and multi-line extended street address information. This attribute may contain newlines. It can have a length of no more than 256 characters.",
						).Description,
						Optional: true,

						Validators: []validator.String{
							stringvalidator.LengthAtLeast(attrMinLength),
							stringvalidator.LengthAtMost(attrAddressStreetAddressMaxLength),
							stringvalidator.RegexMatches(regexp.MustCompile(`^[\p{L}\p{M}\p{N}\p{Zs}\p{P}\n\r]*$`), `must match the regex ^[\p{L}\p{M}\p{N}\p{Zs}\p{P}\n\r]*$`),
						},
					},
				},
			},

			"external_id": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown(
					"A string that specifies an identifier for the user resource as defined by the provisioning client. This may be explicitly set to null when updating a user to unset it. The external id attribute simplifies the correlation of the user in PingOne with the user's account in another system of record. The platform does not use this attribute directly in any way, but it is used by Ping Identity's Data Sync product. It can have a length of no more than 1024 characters.",
				).Description,
				Optional: true,

				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
					stringvalidator.LengthAtMost(attrExternalIdMaxLength),
				},
			},

			"identity_provider": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown(
					"A single object that specifies the user's identity provider information.",
				).Description,
				Optional: true,
				Computed: true,

				Default: objectdefault.StaticValue(func() basetypes.ObjectValue {
					o := map[string]attr.Value{
						"id":   pingonetypes.NewResourceIDNull(),
						"type": types.StringValue(string(management.ENUMIDENTITYPROVIDER_PING_ONE)),
					}

					objValue, d := types.ObjectValue(userIdentityProviderTFObjectTypes, o)
					resp.Diagnostics.Append(d...)

					return objValue
				}()),

				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description:         identityProviderIdDescription.Description,
						MarkdownDescription: identityProviderIdDescription.MarkdownDescription,
						Optional:            true,

						CustomType: pingonetypes.ResourceIDType{},

						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},

					"type": schema.StringAttribute{
						Description:         identityProviderTypeDescription.Description,
						MarkdownDescription: identityProviderTypeDescription.MarkdownDescription,
						Computed:            true,

						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},

			"user_lifecycle": schema.SingleNestedAttribute{
				Description:         userLifecycleDescription.Description,
				MarkdownDescription: userLifecycleDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,

				Default: objectdefault.StaticValue(func() basetypes.ObjectValue {
					o := map[string]attr.Value{
						"status":                     framework.StringToTF(string(management.ENUMUSERLIFECYCLESTATUS_ACCOUNT_OK)),
						"suppress_verification_code": types.BoolNull(),
					}

					objValue, d := types.ObjectValue(userLifecycleTFObjectTypes, o)
					resp.Diagnostics.Append(d...)

					return objValue
				}()),

				Attributes: map[string]schema.Attribute{
					"status": schema.StringAttribute{
						Description:         userLifecycleStatusDescription.Description,
						MarkdownDescription: userLifecycleStatusDescription.MarkdownDescription,
						Optional:            true,

						Validators: []validator.String{
							stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumUserLifecycleStatusEnumValues)...),
						},

						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},

					"suppress_verification_code": schema.BoolAttribute{
						Description:         userLifecycleSuppressVerificationCodeDescription.Description,
						MarkdownDescription: userLifecycleSuppressVerificationCodeDescription.MarkdownDescription,
						Optional:            true,

						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.RequiresReplace(),
						},
					},
				},

				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
			},

			"locale": schema.StringAttribute{
				Description:         localeDescription.Description,
				MarkdownDescription: localeDescription.MarkdownDescription,
				Optional:            true,

				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
					stringvalidator.LengthAtMost(attrLocaleMaxLength),
					stringvalidator.RegexMatches(regexp.MustCompile(`^[\p{L}\p{M}\p{Zs}\p{S}\p{N}\p{P}]*$`), `must match the regex ^[\p{L}\p{M}\p{Zs}\p{S}\p{N}\p{P}]*$`),
				},
			},

			"mfa_enabled": schema.BoolAttribute{
				Description:         mfaEnabledDescription.Description,
				MarkdownDescription: mfaEnabledDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,

				Default: booldefault.StaticBool(false),
			},

			"mobile_phone": schema.StringAttribute{
				Description:         mobilePhoneDescription.Description,
				MarkdownDescription: mobilePhoneDescription.MarkdownDescription,
				Optional:            true,

				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
					stringvalidator.LengthAtMost(attrPhoneMaxLength),
					stringvalidator.RegexMatches(regexp.MustCompile(`[0-9]+`), `must have at least one number`),
				},
			},

			"name": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown(
					"A single object that specifies the user's name information.",
				).Description,
				Optional: true,

				Attributes: map[string]schema.Attribute{
					"family": schema.StringAttribute{
						Description:         nameFamilyDescription.Description,
						MarkdownDescription: nameFamilyDescription.MarkdownDescription,
						Optional:            true,

						Validators: []validator.String{
							stringvalidator.LengthAtLeast(attrMinLength),
							stringvalidator.LengthAtMost(attrNameFamilyMaxLength),
							stringvalidator.RegexMatches(regexp.MustCompile(`^[\p{L}\p{M}\p{N}' .-]*$`), `must match regex ^[\p{L}\p{M}\p{N}' .-]*$`),
						},
					},

					"formatted": schema.StringAttribute{
						Description:         nameFormattedDescription.Description,
						MarkdownDescription: nameFormattedDescription.MarkdownDescription,
						Optional:            true,

						Validators: []validator.String{
							stringvalidator.LengthAtLeast(attrMinLength),
							stringvalidator.LengthAtMost(attrNameFormattedMaxLength),
							stringvalidator.RegexMatches(regexp.MustCompile(`^[\p{L}\p{M}\p{N}' .-]*$`), `must match regex ^[\p{L}\p{M}\p{N}' .-]*$`),
						},
					},

					"given": schema.StringAttribute{
						Description:         nameGivenDescription.Description,
						MarkdownDescription: nameGivenDescription.MarkdownDescription,
						Optional:            true,

						Validators: []validator.String{
							stringvalidator.LengthAtLeast(attrMinLength),
							stringvalidator.LengthAtMost(attrNameGivenMaxLength),
							stringvalidator.RegexMatches(regexp.MustCompile(`^[\p{L}\p{M}\p{Zs}\p{S}\p{N}\p{P}]*$`), `must match regex ^[\p{L}\p{M}\p{Zs}\p{S}\p{N}\p{P}]*$`),
						},
					},

					"honorific_prefix": schema.StringAttribute{
						Description:         nameHonorificPrefixDescription.Description,
						MarkdownDescription: nameHonorificPrefixDescription.MarkdownDescription,
						Optional:            true,
					},

					"honorific_suffix": schema.StringAttribute{
						Description:         nameHonorificSuffixDescription.Description,
						MarkdownDescription: nameHonorificSuffixDescription.MarkdownDescription,
						Optional:            true,
					},

					"middle": schema.StringAttribute{
						Description:         nameMiddleDescription.Description,
						MarkdownDescription: nameMiddleDescription.MarkdownDescription,
						Optional:            true,

						Validators: []validator.String{
							stringvalidator.LengthAtLeast(attrMinLength),
							stringvalidator.LengthAtMost(attrNameMiddleMaxLength),
							stringvalidator.RegexMatches(regexp.MustCompile(`^[\p{L}\p{M}\p{Zs}\p{S}\p{N}\p{P}]*$`), `must match regex ^[\p{L}\p{M}\p{Zs}\p{S}\p{N}\p{P}]*$`),
						},
					},
				},
			},

			"nickname": schema.StringAttribute{
				Description:         nicknameDescription.Description,
				MarkdownDescription: nicknameDescription.MarkdownDescription,
				Optional:            true,

				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
					stringvalidator.LengthAtMost(attrNicknameMaxLength),
					stringvalidator.RegexMatches(regexp.MustCompile(`^[\p{L}\p{M}\p{Zs}\p{S}\p{N}\p{P}]*$`), `must match regex ^[\p{L}\p{M}\p{Zs}\p{S}\p{N}\p{P}]*$`),
				},
			},

			"password": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown(
					"A single object that specifies the user's password information.",
				).Description,
				Optional: true,

				Attributes: map[string]schema.Attribute{
					"force_change": schema.BoolAttribute{
						Description:         passwordForceChangeDescription.Description,
						MarkdownDescription: passwordForceChangeDescription.MarkdownDescription,
						Optional:            true,
						Computed:            true,

						Default: booldefault.StaticBool(false),
					},

					"initial_value": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown(
							"A string that specifies the user's initial password value. The string is either in cleartext or pre-encoded format.  User passwords cannot be extracted from the platfom.  This value, if defined or changed on the PingOne service by an identity administrator or the user account's owner, will not be refreshed in the Terraform state.",
						).Description,
						Optional:  true,
						Sensitive: true,
					},

					"external": schema.SingleNestedAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown(
							"A single object that maps the information relevant to the user's password, and its association to external directories.",
						).Description,
						Optional: true,

						Attributes: map[string]schema.Attribute{
							"gateway": schema.SingleNestedAttribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown(
									"A single object that contains the external gateway properties. When this is value is specified, the user's password is managed in an external directory.",
								).Description,
								Required: true,

								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: framework.SchemaAttributeDescriptionFromMarkdown(
											"A string that specifies the UUID of the linked gateway that references the remote directory.  Must be a valid PingOne resource ID.",
										).Description,
										Optional: true,

										CustomType: pingonetypes.ResourceIDType{},
									},

									"type": schema.StringAttribute{
										Description:         passwordExternalGatewayTypeDescription.Description,
										MarkdownDescription: passwordExternalGatewayTypeDescription.MarkdownDescription,
										Optional:            true,

										Validators: []validator.String{
											stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumGatewayTypeEnumValues)...),
										},
									},

									"user_type_id": schema.StringAttribute{
										Description: framework.SchemaAttributeDescriptionFromMarkdown(
											"A string that specifies the UUID of a user type in the list of user types for the LDAP gateway.  Must be a valid PingOne resource ID.",
										).Description,
										Optional: true,

										CustomType: pingonetypes.ResourceIDType{},
									},

									"correlation_attributes": schema.MapAttribute{
										Description: framework.SchemaAttributeDescriptionFromMarkdown(
											"A string map that maps the external LDAP directory attributes to PingOne attributes. PingOne uses these values to read the attributes from the external LDAP directory and map them to the corresponding PingOne attributes.",
										).Description,
										Optional: true,

										ElementType: types.StringType,
									},
								},
							},
						},
					},
				},
			},

			"photo": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown(
					"A single object that describes the user's photo information.",
				).Description,
				Optional: true,

				Attributes: map[string]schema.Attribute{
					"href": schema.StringAttribute{
						Description:         photoHrefDescription.Description,
						MarkdownDescription: photoHrefDescription.MarkdownDescription,
						Required:            true,

						Validators: []validator.String{
							stringvalidator.RegexMatches(verify.IsURLWithHTTPorHTTPS, "must be a valid URL with HTTP or HTTPS"),
						},
					},
				},
			},

			"preferred_language": schema.StringAttribute{
				Description:         preferredLanguageDescription.Description,
				MarkdownDescription: preferredLanguageDescription.MarkdownDescription,
				Optional:            true,
			},

			"primary_phone": schema.StringAttribute{
				Description:         primaryPhoneDescription.Description,
				MarkdownDescription: primaryPhoneDescription.MarkdownDescription,
				Optional:            true,

				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
					stringvalidator.LengthAtMost(attrPhoneMaxLength),
					stringvalidator.RegexMatches(regexp.MustCompile(`[0-9]+`), `must have at least one number`),
				},
			},

			"timezone": schema.StringAttribute{
				Description:         timezoneDescription.Description,
				MarkdownDescription: timezoneDescription.MarkdownDescription,
				Optional:            true,

				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(`^\w+\/\w+$`), `must match regex ^\w+\/\w+$`),
				},
			},

			"title": schema.StringAttribute{
				Description:         titleDescription.Description,
				MarkdownDescription: titleDescription.MarkdownDescription,
				Optional:            true,

				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
					stringvalidator.LengthAtMost(attrTitleMaxLength),
					stringvalidator.RegexMatches(regexp.MustCompile(`^[\p{L}\p{M}\p{Zs}\p{S}\p{N}\p{P}]*$`), `must match regex ^[\p{L}\p{M}\p{Zs}\p{S}\p{N}\p{P}]*$`),
				},
			},

			"type": schema.StringAttribute{
				Description:         typeDescription.Description,
				MarkdownDescription: typeDescription.MarkdownDescription,
				Optional:            true,

				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
					stringvalidator.LengthAtMost(attrTypeMaxLength),
					stringvalidator.RegexMatches(regexp.MustCompile(`^[\p{L}\p{M}\p{Zs}\p{S}\p{N}\p{P}]*$`), `must match regex ^[\p{L}\p{M}\p{Zs}\p{S}\p{N}\p{P}]*$`),
				},
			},

			"verify_status": schema.StringAttribute{
				Description:         verifyStatusDescription.Description,
				MarkdownDescription: verifyStatusDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,

				Default: stringdefault.StaticString(string(management.ENUMUSERVERIFYSTATUS_NOT_INITIATED)),

				Validators: []validator.String{
					stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumUserVerifyStatusEnumValues)...),
				},

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *UserResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Destruction plan
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan UserResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// if plan.EmailVerified.IsUnknown() {
	// 	resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("email_verified"), types.BoolNull())...)
	// }

	if plan.Account.IsNull() || plan.Account.IsUnknown() {

		o := map[string]attr.Value{
			"can_authenticate": types.BoolValue(true),
			"locked_at":        timetypes.NewRFC3339Null(),
			"status":           types.StringValue(string(management.ENUMUSERSTATUS_OK)),
		}

		if !plan.Enabled.IsNull() && !plan.Enabled.IsUnknown() {
			o["can_authenticate"] = plan.Enabled
		} else {
			o["can_authenticate"] = types.BoolValue(true)
		}

		objValue, d := types.ObjectValue(userAccountTFObjectTypes, o)
		resp.Diagnostics.Append(d...)

		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("account"), objValue)...)
	} else {
		var accountPlan *UserAccountResourceModel
		resp.Diagnostics.Append(resp.Plan.GetAttribute(ctx, path.Root("account"), &accountPlan)...)

		if !accountPlan.Status.Equal(types.StringValue("LOCKED")) {
			resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("account").AtName("locked_at"), timetypes.NewRFC3339Null())...)
		}
	}
}

func (r *UserResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *UserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state UserResourceModel

	if r.Client == nil || r.Client.ManagementAPIClient == nil {
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
	user, userEnabled, userMFAEnabled, d := plan.expand(ctx, false)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	// Create the user
	var createUserResponse *management.User
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.UsersApi.CreateUser(ctx, plan.EnvironmentId.ValueString()).ContentType("application/vnd.pingidentity.user.import+json").User(*user).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"CreateUser",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		&createUserResponse,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the user enabled attribute
	var updateUserEnabledResponse *management.UserEnabled
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.EnableUsersApi.UpdateUserEnabled(ctx, plan.EnvironmentId.ValueString(), createUserResponse.GetId()).UserEnabled(*userEnabled).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateUserEnabled",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		&updateUserEnabledResponse,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the MFA enabled attribute
	var updateUserMfaEnabledResponse *mfa.UserMFAEnabled
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.MFAAPIClient.EnableUsersMFAApi.UpdateUserMFAEnabled(ctx, plan.EnvironmentId.ValueString(), createUserResponse.GetId()).UserMFAEnabled(*userMFAEnabled).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateUserMFAEnabled",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		&updateUserMfaEnabledResponse,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	//User account lock
	if updateUserEnabledResponse.GetEnabled() {
		if account, ok := user.GetAccountOk(); ok {
			if status, ok := account.GetStatusOk(); ok && *status == management.ENUMUSERSTATUS_LOCKED {
				accountStatus := management.ENUMUSERACCOUNTCONTENTTYPEHEADER_LOCKJSON

				resp.Diagnostics.Append(framework.ParseResponse(
					ctx,

					func() (any, *http.Response, error) {
						fO, fR, fErr := r.Client.ManagementAPIClient.UserAccountsApi.UserAccount(ctx, plan.EnvironmentId.ValueString(), createUserResponse.GetId()).ContentType(accountStatus).Execute()
						return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
					},
					"UserAccount",
					framework.DefaultCustomError,
					userAccountEnabledRetryable,
					nil,
				)...)
				if resp.Diagnostics.HasError() {
					return
				}
			}
		}
	}

	// Read the user object again, as other attributes may have changed following the update API calls
	var finalUserResponse *management.User
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.UsersApi.ReadUser(ctx, plan.EnvironmentId.ValueString(), createUserResponse.GetId()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadUser",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		&finalUserResponse,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(ctx, finalUserResponse)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *UserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *UserResourceModel

	if r.Client == nil || r.Client.ManagementAPIClient == nil {
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
	var response *management.User
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.UsersApi.ReadUser(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadUser",
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

	var responseEnabled *management.UserEnabled
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.EnableUsersApi.ReadUserEnabled(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadUserEnabled",
		framework.CustomErrorResourceNotFoundWarning,
		sdk.DefaultCreateReadRetryable,
		&responseEnabled,
	)...)

	// Remove from state if resource is not found
	if responseEnabled == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(data.toState(ctx, response)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state UserResourceModel

	if r.Client == nil || r.Client.ManagementAPIClient == nil {
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
	user, userEnabled, userMFAEnabled, d := plan.expand(ctx, true)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API calls
	if state.PopulationId.ValueString() != plan.PopulationId.ValueString() {

		userPopulation := management.NewUserPopulation(plan.PopulationId.ValueString())

		resp.Diagnostics.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				fO, fR, fErr := r.Client.ManagementAPIClient.UserPopulationsApi.UpdateUserPopulation(ctx, plan.EnvironmentId.ValueString(), plan.Id.ValueString()).UserPopulation(*userPopulation).Execute()
				return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
			},
			"UpdateUserPopulation",
			framework.DefaultCustomError,
			userAccountEnabledRetryable,
			nil,
		)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.UsersApi.UpdateUserPut(ctx, plan.EnvironmentId.ValueString(), plan.Id.ValueString()).User(*user).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateUserPut",
		framework.DefaultCustomError,
		nil,
		nil,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var updateUserEnabledResponse *management.UserEnabled
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.EnableUsersApi.UpdateUserEnabled(ctx, plan.EnvironmentId.ValueString(), plan.Id.ValueString()).UserEnabled(*userEnabled).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateUserEnabled",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		&updateUserEnabledResponse,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the MFA enabled attribute
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.MFAAPIClient.EnableUsersMFAApi.UpdateUserMFAEnabled(ctx, plan.EnvironmentId.ValueString(), plan.Id.ValueString()).UserMFAEnabled(*userMFAEnabled).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateUserMFAEnabled",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		nil,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	//User account lock
	if updateUserEnabledResponse.GetEnabled() {
		accountStatus := management.ENUMUSERACCOUNTCONTENTTYPEHEADER_UNLOCKJSON
		if account, ok := user.GetAccountOk(); ok {
			if status, ok := account.GetStatusOk(); ok && *status == management.ENUMUSERSTATUS_LOCKED {
				accountStatus = management.ENUMUSERACCOUNTCONTENTTYPEHEADER_LOCKJSON
			}
		}

		resp.Diagnostics.Append(framework.ParseResponse(
			ctx,

			func() (any, *http.Response, error) {
				fO, fR, fErr := r.Client.ManagementAPIClient.UserAccountsApi.UserAccount(ctx, plan.EnvironmentId.ValueString(), plan.Id.ValueString()).ContentType(accountStatus).Execute()
				return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
			},
			"UserAccount",
			framework.DefaultCustomError,
			userAccountEnabledRetryable,
			nil,
		)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	var finalUserResponse *management.User
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.UsersApi.ReadUser(ctx, plan.EnvironmentId.ValueString(), plan.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadUser",
		framework.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
		&finalUserResponse,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the state to save
	state = plan

	// Save updated data into Terraform state
	resp.Diagnostics.Append(state.toState(ctx, finalUserResponse)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *UserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *UserResourceModel

	if r.Client == nil || r.Client.ManagementAPIClient == nil {
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
			fR, fErr := r.Client.ManagementAPIClient.UsersApi.DeleteUser(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, fR, fErr)
		},
		"DeleteUser",
		framework.CustomErrorResourceNotFoundWarning,
		nil,
		nil,
	)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *UserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:     "user_id",
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

var (
	userAccountEnabledRetryable = func(ctx context.Context, r *http.Response, p1error *model.P1Error) bool {

		// Catch observed race condition
		if r.StatusCode == http.StatusUnsupportedMediaType {
			tflog.Warn(ctx, "Unexpected 415 Error detected. Available for retry..")
			return true
		}

		return sdk.DefaultCreateReadRetryable(ctx, r, p1error)
	}
)

func (p *UserResourceModel) expand(ctx context.Context, isUpdate bool) (*management.User, *management.UserEnabled, *mfa.UserMFAEnabled, diag.Diagnostics) {
	var diags diag.Diagnostics

	userData := management.NewUser(p.Email.ValueString(), p.Username.ValueString())

	population := *management.NewUserPopulation(p.PopulationId.ValueString())
	userData.SetPopulation(population)

	userEnabledData := management.NewUserEnabled()
	if !p.Enabled.IsNull() && !p.Enabled.IsUnknown() {
		userEnabledData.SetEnabled(p.Enabled.ValueBool())
	}

	if !p.Account.IsNull() && !p.Account.IsUnknown() {
		var plan UserAccountResourceModel
		diags.Append(p.Account.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, nil, nil, diags
		}

		v := management.NewUserAccount(
			plan.CanAuthenticate.ValueBool(),
			management.EnumUserStatus(plan.Status.ValueString()),
		)

		userData.SetAccount(*v)
	}

	if !p.Address.IsNull() && !p.Address.IsUnknown() {
		var plan UserAddressResourceModel
		diags.Append(p.Address.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, nil, nil, diags
		}

		v := management.NewUserAddress()

		if !plan.CountryCode.IsNull() && !plan.CountryCode.IsUnknown() {
			v.SetCountryCode(plan.CountryCode.ValueString())
		}

		if !plan.Locality.IsNull() && !plan.Locality.IsUnknown() {
			v.SetLocality(plan.Locality.ValueString())
		}

		if !plan.PostalCode.IsNull() && !plan.PostalCode.IsUnknown() {
			v.SetPostalCode(plan.PostalCode.ValueString())
		}

		if !plan.Region.IsNull() && !plan.Region.IsUnknown() {
			v.SetRegion(plan.Region.ValueString())
		}

		if !plan.StreetAddress.IsNull() && !plan.StreetAddress.IsUnknown() {
			v.SetStreetAddress(plan.StreetAddress.ValueString())
		}

		userData.SetAddress(*v)
	}

	if !p.ExternalId.IsNull() && !p.ExternalId.IsUnknown() {
		userData.SetExternalId(p.ExternalId.ValueString())
	}

	if !isUpdate && !p.IdentityProvider.IsNull() && !p.IdentityProvider.IsUnknown() {

		var plan UserIdentityProviderResourceModel
		diags.Append(p.IdentityProvider.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, nil, nil, diags
		}

		if !plan.Id.IsNull() && !plan.Id.IsUnknown() {
			v := management.NewUserIdentityProvider()
			v.SetId(plan.Id.ValueString())

			userData.SetIdentityProvider(*v)
		}
	}

	if !p.Lifecycle.IsNull() && !p.Lifecycle.IsUnknown() {
		var plan UserLifecycleResourceModel
		diags.Append(p.Lifecycle.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, nil, nil, diags
		}

		v := management.NewUserLifecycle()

		if !plan.Status.IsNull() && !plan.Status.IsUnknown() {
			v.SetStatus(management.EnumUserLifecycleStatus(plan.Status.ValueString()))
		}

		if !plan.SuppressVerificationCode.IsNull() && !plan.SuppressVerificationCode.IsUnknown() {
			v.SetSuppressVerificationCode(plan.SuppressVerificationCode.ValueBool())
		}

		userData.SetLifecycle(*v)
	}

	if !p.Locale.IsNull() && !p.Locale.IsUnknown() {
		userData.SetLocale(p.Locale.ValueString())
	}

	if !p.MobilePhone.IsNull() && !p.MobilePhone.IsUnknown() {
		userData.SetMobilePhone(p.MobilePhone.ValueString())
	}

	if !p.Name.IsNull() && !p.Name.IsUnknown() {
		var plan UserNameResourceModel
		diags.Append(p.Name.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, nil, nil, diags
		}

		v := management.NewUserName()

		if !plan.Family.IsNull() && !plan.Family.IsUnknown() {
			v.SetFamily(plan.Family.ValueString())
		}

		if !plan.Formatted.IsNull() && !plan.Formatted.IsUnknown() {
			v.SetFormatted(plan.Formatted.ValueString())
		}

		if !plan.Given.IsNull() && !plan.Given.IsUnknown() {
			v.SetGiven(plan.Given.ValueString())
		}

		if !plan.HonorificPrefix.IsNull() && !plan.HonorificPrefix.IsUnknown() {
			v.SetHonorificPrefix(plan.HonorificPrefix.ValueString())
		}

		if !plan.HonorificSuffix.IsNull() && !plan.HonorificSuffix.IsUnknown() {
			v.SetHonorificSuffix(plan.HonorificSuffix.ValueString())
		}

		if !plan.Middle.IsNull() && !plan.Middle.IsUnknown() {
			v.SetMiddle(plan.Middle.ValueString())
		}

		userData.SetName(*v)
	}

	if !p.Nickname.IsNull() && !p.Nickname.IsUnknown() {
		userData.SetNickname(p.Nickname.ValueString())
	}

	if !p.Password.IsNull() && !p.Password.IsUnknown() {
		var plan UserPasswordResourceModel
		diags.Append(p.Password.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, nil, nil, diags
		}

		v := management.NewUserPassword()

		if !plan.ForceChange.IsNull() && !plan.ForceChange.IsUnknown() {
			v.SetForceChange(plan.ForceChange.ValueBool())
		}

		if !plan.InitialValue.IsNull() && !plan.InitialValue.IsUnknown() {
			v.SetValue(plan.InitialValue.ValueString())
		}

		if !plan.External.IsNull() && !plan.External.IsUnknown() {
			external := management.NewUserPasswordExternal()
			v.SetExternal(*external)
		}

		userData.SetPassword(*v)
	}

	if !p.Photo.IsNull() && !p.Photo.IsUnknown() {
		var plan UserPhotoResourceModel
		diags.Append(p.Photo.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, nil, nil, diags
		}

		v := management.NewUserPhoto(plan.Href.ValueString())

		userData.SetPhoto(*v)
	}

	if !p.PreferredLanguage.IsNull() && !p.PreferredLanguage.IsUnknown() {
		userData.SetPreferredLanguage(p.PreferredLanguage.ValueString())
	}

	if !p.PrimaryPhone.IsNull() && !p.PrimaryPhone.IsUnknown() {
		userData.SetPrimaryPhone(p.PrimaryPhone.ValueString())
	}

	if !p.Timezone.IsNull() && !p.Timezone.IsUnknown() {
		userData.SetTimezone(p.Timezone.ValueString())
	}

	if !p.Title.IsNull() && !p.Title.IsUnknown() {
		userData.SetTitle(p.Title.ValueString())
	}

	if !p.Type.IsNull() && !p.Type.IsUnknown() {
		userData.SetType(p.Type.ValueString())
	}

	if !p.VerifyStatus.IsNull() && !p.VerifyStatus.IsUnknown() {
		userData.SetVerifyStatus(management.EnumUserVerifyStatus(p.VerifyStatus.ValueString()))
	}

	userMFAEnabledData := mfa.NewUserMFAEnabled(p.MFAEnabled.ValueBool())

	return userData, userEnabledData, userMFAEnabledData, diags
}

func (p *UserResourceModel) toState(ctx context.Context, apiObject *management.User) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.PingOneResourceIDOkToTF(apiObject.GetIdOk())
	p.EnvironmentId = framework.PingOneResourceIDOkToTF(apiObject.Environment.GetIdOk())
	p.Username = framework.StringOkToTF(apiObject.GetUsernameOk())
	p.Email = framework.StringOkToTF(apiObject.GetEmailOk())
	//p.EmailVerified = framework.BoolOkToTF(apiObject.GetEmailVerifiedOk())
	p.Enabled = framework.BoolOkToTF(apiObject.GetEnabledOk())

	var d diag.Diagnostics

	p.PopulationId = framework.PingOneResourceIDOkToTF(apiObject.Population.GetIdOk())
	p.Account, d = p.userAccountOkToTF(apiObject.GetAccountOk())
	diags = append(diags, d...)

	p.Address, d = p.userAddressOkToTF(apiObject.GetAddressOk())
	diags = append(diags, d...)

	p.ExternalId = framework.StringOkToTF(apiObject.GetExternalIdOk())
	p.IdentityProvider, d = p.userIdentityProviderOkToTF(apiObject.GetIdentityProviderOk())
	diags = append(diags, d...)

	lifecycle, lifecycleOk := apiObject.GetLifecycleOk()
	p.Lifecycle, d = p.userLifecycleOkToTF(ctx, lifecycle, lifecycleOk)
	diags = append(diags, d...)

	p.Locale = framework.StringOkToTF(apiObject.GetLocaleOk())
	p.MFAEnabled = framework.BoolOkToTF(apiObject.GetMfaEnabledOk())
	p.MobilePhone = framework.StringOkToTF(apiObject.GetMobilePhoneOk())
	p.Name, d = p.userNameOkToTF(apiObject.GetNameOk())
	diags = append(diags, d...)

	p.Nickname = framework.StringOkToTF(apiObject.GetNicknameOk())

	password, passwordOk := apiObject.GetPasswordOk()
	p.Password, d = p.userPasswordOkToTF(ctx, password, passwordOk)
	diags = append(diags, d...)

	p.Photo, d = p.photoOkToTF(apiObject.GetPhotoOk())
	diags = append(diags, d...)

	p.PreferredLanguage = framework.StringOkToTF(apiObject.GetPreferredLanguageOk())
	p.PrimaryPhone = framework.StringOkToTF(apiObject.GetPrimaryPhoneOk())
	p.Timezone = framework.StringOkToTF(apiObject.GetTimezoneOk())
	p.Title = framework.StringOkToTF(apiObject.GetTitleOk())
	p.Type = framework.StringOkToTF(apiObject.GetTypeOk())
	p.VerifyStatus = framework.EnumOkToTF(apiObject.GetVerifyStatusOk())

	return diags
}

func (p *UserResourceModel) userAccountOkToTF(apiObject *management.UserAccount, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(userAccountTFObjectTypes), diags
	}

	objMap := map[string]attr.Value{
		"can_authenticate": framework.BoolOkToTF(apiObject.GetCanAuthenticateOk()),
		"locked_at":        framework.TimeOkToTF(apiObject.GetLockedAtOk()),
		"status":           framework.EnumOkToTF(apiObject.GetStatusOk()),
	}

	objValue, d := types.ObjectValue(userAccountTFObjectTypes, objMap)
	diags.Append(d...)

	return objValue, diags
}

func (p *UserResourceModel) userAddressOkToTF(apiObject *management.UserAddress, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(userAddressTFObjectTypes), diags
	}

	objMap := map[string]attr.Value{
		"country_code":   framework.StringOkToTF(apiObject.GetCountryCodeOk()),
		"locality":       framework.StringOkToTF(apiObject.GetLocalityOk()),
		"postal_code":    framework.StringOkToTF(apiObject.GetPostalCodeOk()),
		"region":         framework.StringOkToTF(apiObject.GetRegionOk()),
		"street_address": framework.StringOkToTF(apiObject.GetStreetAddressOk()),
	}

	objValue, d := types.ObjectValue(userAddressTFObjectTypes, objMap)
	diags.Append(d...)

	return objValue, diags
}

func (p *UserResourceModel) userIdentityProviderOkToTF(apiObject *management.UserIdentityProvider, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(userIdentityProviderTFObjectTypes), diags
	}

	objMap := map[string]attr.Value{
		"id":   framework.PingOneResourceIDOkToTF(apiObject.GetIdOk()),
		"type": framework.EnumOkToTF(apiObject.GetTypeOk()),
	}

	objValue, d := types.ObjectValue(userIdentityProviderTFObjectTypes, objMap)
	diags.Append(d...)

	return objValue, diags
}

func (p *UserResourceModel) userLifecycleOkToTF(ctx context.Context, apiObject *management.UserLifecycle, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(userLifecycleTFObjectTypes), diags
	}

	suppressVerificationCode := types.BoolNull()
	if !p.Lifecycle.IsNull() && !p.Lifecycle.IsUnknown() {
		var plan *UserLifecycleResourceModel
		diags.Append(p.Lifecycle.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return types.ObjectNull(userLifecycleTFObjectTypes), diags
		}

		suppressVerificationCode = plan.SuppressVerificationCode
	}

	objMap := map[string]attr.Value{
		"status":                     framework.EnumOkToTF(apiObject.GetStatusOk()),
		"suppress_verification_code": suppressVerificationCode,
	}

	objValue, d := types.ObjectValue(userLifecycleTFObjectTypes, objMap)
	diags.Append(d...)

	return objValue, diags
}

func (p *UserResourceModel) userNameOkToTF(apiObject *management.UserName, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(userNameTFObjectTypes), diags
	}

	objMap := map[string]attr.Value{
		"family":           framework.StringOkToTF(apiObject.GetFamilyOk()),
		"formatted":        framework.StringOkToTF(apiObject.GetFormattedOk()),
		"given":            framework.StringOkToTF(apiObject.GetGivenOk()),
		"honorific_prefix": framework.StringOkToTF(apiObject.GetHonorificPrefixOk()),
		"honorific_suffix": framework.StringOkToTF(apiObject.GetHonorificSuffixOk()),
		"middle":           framework.StringOkToTF(apiObject.GetMiddleOk()),
	}

	objValue, d := types.ObjectValue(userNameTFObjectTypes, objMap)
	diags.Append(d...)

	return objValue, diags
}

func (p *UserResourceModel) userPasswordOkToTF(ctx context.Context, apiObject *management.UserPassword, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// The API object might be nil even though the plan is not.  We need to fill the state from the plan if it is
	if (!ok || apiObject == nil) && p.Password.IsNull() {
		return types.ObjectNull(userPasswordTFObjectTypes), diags
	}

	var plan *UserPasswordResourceModel
	diags.Append(p.Password.As(ctx, &plan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return types.ObjectNull(userPasswordTFObjectTypes), diags
	}

	externalObject := types.ObjectNull(userPasswordExternalTFObjectTypes)

	objMap := map[string]attr.Value{
		"force_change":  plan.ForceChange,
		"initial_value": plan.InitialValue,
		"external":      externalObject,
	}

	objValue, d := types.ObjectValue(userPasswordTFObjectTypes, objMap)
	diags.Append(d...)

	return objValue, diags
}

func (p *UserResourceModel) photoOkToTF(apiObject *management.UserPhoto, ok bool) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(userPhotoTFObjectTypes), diags
	}

	objMap := map[string]attr.Value{
		"href": framework.StringOkToTF(apiObject.GetHrefOk()),
	}

	objValue, d := types.ObjectValue(userPhotoTFObjectTypes, objMap)
	diags.Append(d...)

	return objValue, diags
}
