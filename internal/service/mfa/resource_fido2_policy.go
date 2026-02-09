// Copyright © 2026 Ping Identity Corporation

package mfa

import (
	"context"
	"fmt"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/patrickcping/pingone-go-sdk-v2/mfa"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	int32validatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/int32validator"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/legacysdk"
	setvalidatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/setvalidator"
	stringvalidatorinternal "github.com/pingidentity/terraform-provider-pingone/internal/framework/stringvalidator"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	listvalidatormfa "github.com/pingidentity/terraform-provider-pingone/internal/service/mfa/listvalidator"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type FIDO2PolicyResource serviceClientType

type FIDO2PolicyResourceModel struct {
	Id                            pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId                 pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	Name                          types.String                 `tfsdk:"name"`
	Description                   types.String                 `tfsdk:"description"`
	Default                       types.Bool                   `tfsdk:"default"`
	AttestationRequirements       types.String                 `tfsdk:"attestation_requirements"`
	AuthenticatorAttachment       types.String                 `tfsdk:"authenticator_attachment"`
	BackupEligibility             types.Object                 `tfsdk:"backup_eligibility"`
	DeviceDisplayName             types.String                 `tfsdk:"device_display_name"`
	DiscoverableCredentials       types.String                 `tfsdk:"discoverable_credentials"`
	MdsAuthenticatorsRequirements types.Object                 `tfsdk:"mds_authenticators_requirements"`
	RelyingPartyId                types.String                 `tfsdk:"relying_party_id"`
	UserDisplayNameAttributes     types.Object                 `tfsdk:"user_display_name_attributes"`
	UserPresenceTimeout           types.Object                 `tfsdk:"user_presence_timeout"`
	UserVerification              types.Object                 `tfsdk:"user_verification"`
}

type FIDO2PolicyBackupEligibilityResourceModel struct {
	Allow                       types.Bool `tfsdk:"allow"`
	EnforceDuringAuthentication types.Bool `tfsdk:"enforce_during_authentication"`
}

type FIDO2PolicyMdsAuthenticatorsRequirementsResourceModel struct {
	AllowedAuthenticatorIDs     types.Set    `tfsdk:"allowed_authenticator_ids"`
	EnforceDuringAuthentication types.Bool   `tfsdk:"enforce_during_authentication"`
	Option                      types.String `tfsdk:"option"`
}

type FIDO2PolicyUserDisplayNameAttributesResourceModel struct {
	Attributes types.List `tfsdk:"attributes"`
}

type FIDO2PolicyUserDisplayNameAttributesAttributesResourceModel struct {
	Name          types.String `tfsdk:"name"`
	SubAttributes types.List   `tfsdk:"sub_attributes"`
}

type FIDO2PolicyUserDisplayNameAttributesAttributesSubAttributesResourceModel struct {
	Name types.String `tfsdk:"name"`
}

type FIDO2PolicyUserPresenceTimeoutResourceModel struct {
	Duration types.Int32  `tfsdk:"duration"`
	TimeUnit types.String `tfsdk:"time_unit"`
}

type FIDO2PolicyUserVerificationResourceModel struct {
	EnforceDuringAuthentication types.Bool   `tfsdk:"enforce_during_authentication"`
	Option                      types.String `tfsdk:"option"`
}

var (
	fido2PolicyBackupEligibilityTFObjectTypes = map[string]attr.Type{
		"allow":                         types.BoolType,
		"enforce_during_authentication": types.BoolType,
	}

	fido2PolicyMdsAuthenticatorRequirementsTFObjectTypes = map[string]attr.Type{
		"allowed_authenticator_ids":     types.SetType{ElemType: types.StringType},
		"enforce_during_authentication": types.BoolType,
		"option":                        types.StringType,
	}

	fido2PolicyUserDisplayNameAttributesTFObjectTypes = map[string]attr.Type{
		"attributes": types.ListType{ElemType: types.ObjectType{AttrTypes: fido2PolicyUserDisplayNameAttributesAttributesTFObjectTypes}},
	}

	fido2PolicyUserDisplayNameAttributesAttributesTFObjectTypes = map[string]attr.Type{
		"sub_attributes": types.ListType{ElemType: types.ObjectType{AttrTypes: fido2PolicyUserDisplayNameAttributesAttributesSubAttributesTFObjectTypes}},
		"name":           types.StringType,
	}

	fido2PolicyUserDisplayNameAttributesAttributesSubAttributesTFObjectTypes = map[string]attr.Type{
		"name": types.StringType,
	}

	fido2PolicyUserPresenceTimeoutTFObjectTypes = map[string]attr.Type{
		"duration":  types.Int32Type,
		"time_unit": types.StringType,
	}

	fido2PolicyUserVerificationTFObjectTypes = map[string]attr.Type{
		"enforce_during_authentication": types.BoolType,
		"option":                        types.StringType,
	}
)

// Framework interfaces
var (
	_ resource.Resource                = &FIDO2PolicyResource{}
	_ resource.ResourceWithConfigure   = &FIDO2PolicyResource{}
	_ resource.ResourceWithImportState = &FIDO2PolicyResource{}
)

// New Object
func NewFIDO2PolicyResource() resource.Resource {
	return &FIDO2PolicyResource{}
}

// Metadata
func (r *FIDO2PolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mfa_fido2_policy"
}

func (r *FIDO2PolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	// schema descriptions and validation settings
	const attrMinLength = 1
	const attrNameMaxLength = 256
	const attrDeviceDisplayNameMaxLength = 100
	const attrMinLifetimeDurationMinutes = 1
	const attrMaxLifetimeDurationMinutes = 10
	const attrMinLifetimeDurationSeconds = 60
	const attrMaxLifetimeDurationSeconds = 600
	const attrDefaultUserPresenceTimeoutDuration = 2

	attestationRequirementsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the level of attestation to apply.",
	).AllowedValuesComplex(map[string]string{
		string(mfa.ENUMFIDO2POLICYATTESTATIONREQUIREMENTS_DIRECT): "perform attestation",
		string(mfa.ENUMFIDO2POLICYATTESTATIONREQUIREMENTS_NONE):   "don't perform attestation",
	}).AppendMarkdownString(fmt.Sprintf("If `%s` is specified, the `mds_authentication_requirements.option` parameter should also be set to `%s`.", mfa.ENUMFIDO2POLICYATTESTATIONREQUIREMENTS_NONE, mfa.ENUMFIDO2POLICYMDSAUTHENTICATOROPTION_NONE))

	authenticatorAttachmentDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the types of authenticators that are allowed.",
	).AllowedValuesComplex(map[string]string{
		string(mfa.ENUMFIDO2POLICYAUTHENTICATORATTACHMENT_PLATFORM):       "only allow the use of FIDO device authenticators that contain an internal authenticator (such as a face or fingerprint scanner)",
		string(mfa.ENUMFIDO2POLICYAUTHENTICATORATTACHMENT_CROSS_PLATFORM): "allow use of cross-platform authenticators, which are external to the accessing device (such as a security key)",
		string(mfa.ENUMFIDO2POLICYAUTHENTICATORATTACHMENT_BOTH):           "allow both categories of authenticators",
	})

	backupEligibilityAllowDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether to allow users to register and authenticate with a device that uses cloud-synced credentials.",
	)

	backupEligibilityEnforceDuringAuthnDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether the backup eligibility of the device should be checked again at each authentication attempt.  Set to `true` if you want the backup eligibility of the device to be checked again at each authentication attempt and not just once during registration. Set to `false` to have it checked only at registration.",
	)

	deviceDisplayNameDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The name to display for the device in registration and authentication windows. Can be up to 100 characters. If you want to use translatable text (configured for each language under **Languages** in the Admin Console), you can use any of the keys listed on the `FIDO Policy` page of the `Self-Service` module and the `Sign On Policy` module. The value of the parameter should include only the part of the key name that comes after the module name, for example, `fidoPolicy.deviceDisplayName01` or `fidoPolicy.deviceDisplayName07`. See each language under the **Languages** section of the admin console UI for the full list of keys. For more information on translatable keys, see [Modifying translatable keys](https://docs.pingidentity.com/access/sources/dita/topic?category=p1&resourceid=pingone_modifying_translatable_keys) in the PingOne documentation.",
	)

	discoverableCredentialsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the behaviour when registered users are authenticating without providing credentials",
	).AllowedValuesComplex(map[string]string{
		string(mfa.ENUMFIDO2POLICYDISCOVERABLECREDENTIALS_DISCOURAGED): "discoverable credentials are not used, even when supported by the FIDO device. In cases where use of discoverable credentials is required by the FIDO device itself, this setting does not override the device setting",
		string(mfa.ENUMFIDO2POLICYDISCOVERABLECREDENTIALS_REQUIRED):    "require the use of discoverable credentials. This option is required for usernameless authentication",
		string(mfa.ENUMFIDO2POLICYDISCOVERABLECREDENTIALS_PREFERRED):   "use discoverable credentials where possible",
	})

	mdsAuthenticatorRequirementsAllowedAuthenticatorIDsDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("A set of strings that is used if `option` is set to `%s`, to specify the mdsIdentitfer IDs of authenticators that are allowed in the policy.", mfa.ENUMFIDO2POLICYMDSAUTHENTICATOROPTION_SPECIFIC),
	)

	mdsAuthenticatorRequirementsEnforceDuringAuthnDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether devices characteristics related to verification are checked again on each authentication attempt.  Set to `true` if you want the device characteristics related to attestation to be checked again at each authentication attempt and not just once during registration. Set to `false` to have them checked only at registration.",
	)

	mdsAuthenticatorRequirementsOptionDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the types of device that are allowed on the basis of the attestation provided.",
	).AllowedValuesComplex(map[string]string{
		string(mfa.ENUMFIDO2POLICYMDSAUTHENTICATOROPTION_NONE):       "do not request attestation, allow all FIDO devices",
		string(mfa.ENUMFIDO2POLICYMDSAUTHENTICATOROPTION_AUDIT_ONLY): "attestation is requested and the information is used for logging purposes, but the information is not used for filtering authenticators",
		string(mfa.ENUMFIDO2POLICYMDSAUTHENTICATOROPTION_GLOBAL):     "allow use of all FIDO authenticators listed in the Global Authenticators table",
		string(mfa.ENUMFIDO2POLICYMDSAUTHENTICATOROPTION_CERTIFIED):  "allow only FIDO Certified authenticators",
		string(mfa.ENUMFIDO2POLICYMDSAUTHENTICATOROPTION_SPECIFIC):   "allow only the authenticators specified with the `allowed_authenticator_ids` parameter",
	})

	relyingPartyIDDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The ID of the relying party. The value should be a domain name, such as `bxretail.org` (in lower-case characters).",
	)

	userDisplayNameAttributesAttributesDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A list of objects that describe attributes associated with the users's account that can be displayed during registration and authentication.\n" +
			"    - The content of the list should reflect the preferred order.\n" +
			"    - If the first attribute is empty for the user, PingOne will continue through the list until a non-empty attribute is found.\n" +
			"    - You can specify any user attribute (including custom attributes) that meet the following criteria: attribute type must be String, validation cannot be set to enumerated values.\n" +
			"    - The array must contain the user attribute `username` to ensure that there is at least one non-empty attribute.\n" +
			"    - You can have a maximum of six user attributes in the list.",
	)

	userDisplayNameAttributesAttributesNameDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The name of the attribute in PingOne, for example `username` or `email`.  The attribute can be any user attribute, including a custom attribute, that is a string data type and does not have enumerated values configured.  If you want to use the `name` attribute for the user (or any attribute that is a complex data type), you must also specify the `sub_attributes` parameter, which can be either the `given` and `family` user attributes or the `formatted` user attribute.",
	)

	userDisplayNameAttributesAttributesSubAttributesDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A lsit of objects that describe the sub attributes to use when `name` is configured to use an attribute that is a complex data type.",
	)

	userDisplayNameAttributesAttributesSubAttributesNameDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The name of a complex attribute's sub attribute in PingOne, for example `given` or `formatted` where the parent object has a name value of `name`.",
	)

	userPresenceTimeoutDescription := framework.SchemaAttributeDescriptionFromMarkdown("A single nested object that specifies the user presence timeout settings, used to control the amount of time a user has to perform a user presence gesture with their FIDO device. If not provided, defaults to 2 minutes.")

	userPresenceTimeoutDurationDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The amount of time (minutes or seconds) a user presence gesture will be accepted for the authentication request. Minimum is one minute (60 seconds); maxiumum is ten minutes (600 seconds).",
	).DefaultValue(attrDefaultUserPresenceTimeoutDuration)

	userPresenceTimeoutTimeUnitDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The units for specifying the amount of time a user presence gesture will be accepted for the authentication request.",
	).AllowedValuesEnum(mfa.AllowedEnumTimeUnitEnumValues).DefaultValue(string(mfa.ENUMTIMEUNIT_MINUTES))

	userVerificationEnforceDuringAuthnDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether device characteristics related to user verification are to be checked again at each authentication attempt. Set to `true` if you want the device characteristics related to user verification to be checked again at each authentication attempt and not just once during registration. Set to `false` to have them checked only at registration.",
	)

	userVerificationOptionDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the type of user verification to perform.",
	).AllowedValuesEnum(mfa.AllowedEnumFIDO2PolicyUserVerificationOptionEnumValues).AllowedValuesComplex(map[string]string{
		string(mfa.ENUMFIDO2POLICYUSERVERIFICATIONOPTION_REQUIRED):    "only FIDO devices supporting user verification can be used",
		string(mfa.ENUMFIDO2POLICYUSERVERIFICATIONOPTION_DISCOURAGED): "user verification is not required, even when supported by the FIDO device. In cases where user verification is required by the FIDO device itself, this setting does not override the device setting",
		string(mfa.ENUMFIDO2POLICYUSERVERIFICATIONOPTION_PREFERRED):   "user verification is required if the user's FIDO device supports it, but is not required if the user's device does not support it",
	}).AppendMarkdownString("For usernameless flows, only FIDO devices supporting user verification can be used, regardless of the value configured in this parameter.")

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage FIDO2 policies in a PingOne environment.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to configure the FIDO2 policy in."),
			),

			"name": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the unique, friendly name for this FIDO2 policy.").Description,
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
					stringvalidator.LengthAtMost(attrNameMaxLength),
				},
			},

			"description": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the description of the FIDO2 policy.").Description,
				Optional:    true,
			},

			"default": schema.BoolAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A boolean that describes whether this policy should serve as the default FIDO policy.").Description,
				Computed:    true,

				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseNonNullStateForUnknown(),
				},
			},

			"attestation_requirements": schema.StringAttribute{
				Description:         attestationRequirementsDescription.Description,
				MarkdownDescription: attestationRequirementsDescription.MarkdownDescription,
				Required:            true,

				Validators: []validator.String{
					stringvalidator.OneOf(utils.EnumSliceToStringSlice(mfa.AllowedEnumFIDO2PolicyAttestationRequirementsEnumValues)...),
					stringvalidatorinternal.ShouldBeDefinedValueIfPathMatchesValue(
						types.StringValue(string(mfa.ENUMFIDO2POLICYATTESTATIONREQUIREMENTS_NONE)),
						types.StringValue(string(mfa.ENUMFIDO2POLICYMDSAUTHENTICATOROPTION_NONE)),
						path.MatchRoot("mds_authenticators_requirements").AtName("option"),
					),
				},
			},

			"authenticator_attachment": schema.StringAttribute{
				Description:         authenticatorAttachmentDescription.Description,
				MarkdownDescription: authenticatorAttachmentDescription.MarkdownDescription,
				Required:            true,

				Validators: []validator.String{
					stringvalidator.OneOf(utils.EnumSliceToStringSlice(mfa.AllowedEnumFIDO2PolicyAuthenticatorAttachmentEnumValues)...),
				},
			},

			"backup_eligibility": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A single nested object that contains settings used to control whether users should be allowed to register and authenticate with a device that uses cloud-synced credentials, such as a passkey.").Description,
				Required:    true,

				Attributes: map[string]schema.Attribute{
					"allow": schema.BoolAttribute{
						Description:         backupEligibilityAllowDescription.Description,
						MarkdownDescription: backupEligibilityAllowDescription.MarkdownDescription,
						Required:            true,
					},

					"enforce_during_authentication": schema.BoolAttribute{
						Description:         backupEligibilityEnforceDuringAuthnDescription.Description,
						MarkdownDescription: backupEligibilityEnforceDuringAuthnDescription.MarkdownDescription,
						Required:            true,
					},
				},
			},

			"device_display_name": schema.StringAttribute{
				Description:         deviceDisplayNameDescription.Description,
				MarkdownDescription: deviceDisplayNameDescription.MarkdownDescription,
				Required:            true,

				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
					stringvalidator.LengthAtMost(attrDeviceDisplayNameMaxLength),
				},
			},

			"discoverable_credentials": schema.StringAttribute{
				Description:         discoverableCredentialsDescription.Description,
				MarkdownDescription: discoverableCredentialsDescription.MarkdownDescription,
				Required:            true,

				Validators: []validator.String{
					stringvalidator.OneOf(utils.EnumSliceToStringSlice(mfa.AllowedEnumFIDO2PolicyDiscoverableCredentialsEnumValues)...),
				},
			},

			"mds_authenticators_requirements": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A single nested object that specifies MDS authenticator requirements, used to specify whether attestation is requested from the authenticator, and whether this information is used to restrict authenticator usage.").Description,
				Required:    true,

				Attributes: map[string]schema.Attribute{
					"allowed_authenticator_ids": schema.SetAttribute{
						Description:         mdsAuthenticatorRequirementsAllowedAuthenticatorIDsDescription.Description,
						MarkdownDescription: mdsAuthenticatorRequirementsAllowedAuthenticatorIDsDescription.MarkdownDescription,
						Optional:            true,

						ElementType: types.StringType,

						Validators: []validator.Set{
							setvalidator.SizeAtLeast(1),
							setvalidatorinternal.IsRequiredIfMatchesPathValue(
								types.StringValue(string(mfa.ENUMFIDO2POLICYMDSAUTHENTICATOROPTION_SPECIFIC)),
								path.MatchRelative().AtParent().AtName("option"),
							),
							setvalidatorinternal.ConflictsIfMatchesPathValue(
								types.StringValue(string(mfa.ENUMFIDO2POLICYMDSAUTHENTICATOROPTION_NONE)),
								path.MatchRelative().AtParent().AtName("option"),
							),
							setvalidatorinternal.ConflictsIfMatchesPathValue(
								types.StringValue(string(mfa.ENUMFIDO2POLICYMDSAUTHENTICATOROPTION_CERTIFIED)),
								path.MatchRelative().AtParent().AtName("option"),
							),
							setvalidatorinternal.ConflictsIfMatchesPathValue(
								types.StringValue(string(mfa.ENUMFIDO2POLICYMDSAUTHENTICATOROPTION_GLOBAL)),
								path.MatchRelative().AtParent().AtName("option"),
							),
							setvalidatorinternal.ConflictsIfMatchesPathValue(
								types.StringValue(string(mfa.ENUMFIDO2POLICYMDSAUTHENTICATOROPTION_AUDIT_ONLY)),
								path.MatchRelative().AtParent().AtName("option"),
							),
						},
					},

					"enforce_during_authentication": schema.BoolAttribute{
						Description:         mdsAuthenticatorRequirementsEnforceDuringAuthnDescription.Description,
						MarkdownDescription: mdsAuthenticatorRequirementsEnforceDuringAuthnDescription.MarkdownDescription,
						Required:            true,
					},

					"option": schema.StringAttribute{
						Description:         mdsAuthenticatorRequirementsOptionDescription.Description,
						MarkdownDescription: mdsAuthenticatorRequirementsOptionDescription.MarkdownDescription,
						Required:            true,

						Validators: []validator.String{
							stringvalidator.OneOf(utils.EnumSliceToStringSlice(mfa.AllowedEnumFIDO2PolicyMDSAuthenticatorOptionEnumValues)...),
							stringvalidatorinternal.ShouldBeDefinedValueIfPathMatchesValue(
								types.StringValue(string(mfa.ENUMFIDO2POLICYMDSAUTHENTICATOROPTION_NONE)),
								types.StringValue(string(mfa.ENUMFIDO2POLICYATTESTATIONREQUIREMENTS_NONE)),
								path.MatchRoot("attestation_requirements"),
							),
						},
					},
				},
			},

			"relying_party_id": schema.StringAttribute{
				Description:         relyingPartyIDDescription.Description,
				MarkdownDescription: relyingPartyIDDescription.MarkdownDescription,
				Required:            true,

				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
					stringvalidator.RegexMatches(regexp.MustCompile(`^(?:[\w-]+\.)+(?:[a-z]{2,}|xn--[a-z0-9]+)$`), "must be a valid domain name"),
				},
			},

			"user_display_name_attributes": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A single nested object that specifies the string associated with the users's account that is displayed during registration and authentication.").Description,
				Required:    true,

				Attributes: map[string]schema.Attribute{
					"attributes": schema.ListNestedAttribute{
						Description:         userDisplayNameAttributesAttributesDescription.Description,
						MarkdownDescription: userDisplayNameAttributesAttributesDescription.MarkdownDescription,
						Required:            true,

						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Description:         userDisplayNameAttributesAttributesNameDescription.Description,
									MarkdownDescription: userDisplayNameAttributesAttributesNameDescription.MarkdownDescription,
									Required:            true,

									Validators: []validator.String{
										stringvalidator.LengthAtLeast(attrMinLength),
									},
								},

								"sub_attributes": schema.ListNestedAttribute{
									Description:         userDisplayNameAttributesAttributesSubAttributesDescription.Description,
									MarkdownDescription: userDisplayNameAttributesAttributesSubAttributesDescription.MarkdownDescription,
									Optional:            true,

									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"name": schema.StringAttribute{
												Description:         userDisplayNameAttributesAttributesSubAttributesNameDescription.Description,
												MarkdownDescription: userDisplayNameAttributesAttributesSubAttributesNameDescription.MarkdownDescription,
												Required:            true,

												Validators: []validator.String{
													stringvalidator.LengthAtLeast(attrMinLength),
												},
											},
										},
									},

									Validators: []validator.List{
										listvalidator.SizeAtLeast(1),
									},
								},
							},
						},

						Validators: []validator.List{
							listvalidator.SizeAtLeast(1),
							listvalidatormfa.FIDO2UserDisplayNameAttributeContainsUsername(),
						},
					},
				},
			},

			"user_presence_timeout": schema.SingleNestedAttribute{
				Description: userPresenceTimeoutDescription.Description,
				Optional:    true,
				Computed:    true,
				Default: objectdefault.StaticValue(types.ObjectValueMust(
					fido2PolicyUserPresenceTimeoutTFObjectTypes,
					map[string]attr.Value{
						"duration":  types.Int32Value(attrDefaultUserPresenceTimeoutDuration),
						"time_unit": types.StringValue(string(mfa.ENUMTIMEUNIT_MINUTES)),
					},
				)),

				Attributes: map[string]schema.Attribute{
					"duration": schema.Int32Attribute{
						Description:         userPresenceTimeoutDurationDescription.Description,
						MarkdownDescription: userPresenceTimeoutDurationDescription.MarkdownDescription,
						Optional:            true,
						Computed:            true,
						Default:             int32default.StaticInt32(attrDefaultUserPresenceTimeoutDuration),

						Validators: []validator.Int32{
							int32validator.Any(
								int32validator.All(
									int32validator.Between(attrMinLifetimeDurationMinutes, attrMaxLifetimeDurationMinutes),
									int32validatorinternal.RegexMatchesPathValue(
										regexp.MustCompile(`MINUTES`),
										fmt.Sprintf("If `time_unit` is `MINUTES`, the allowed duration range is %d - %d.", attrMinLifetimeDurationMinutes, attrMaxLifetimeDurationMinutes),
										path.MatchRelative().AtParent().AtName("time_unit"),
									),
								),
								int32validator.All(
									int32validator.Between(attrMinLifetimeDurationSeconds, attrMaxLifetimeDurationSeconds),
									int32validatorinternal.RegexMatchesPathValue(
										regexp.MustCompile(`SECONDS`),
										fmt.Sprintf("If `time_unit` is `SECONDS`, the allowed duration range is %d - %d.", attrMinLifetimeDurationSeconds, attrMaxLifetimeDurationSeconds),
										path.MatchRelative().AtParent().AtName("time_unit"),
									),
								),
							),
						},
					},
					"time_unit": schema.StringAttribute{
						Description:         userPresenceTimeoutTimeUnitDescription.Description,
						MarkdownDescription: userPresenceTimeoutTimeUnitDescription.MarkdownDescription,
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString(string(mfa.ENUMTIMEUNIT_MINUTES)),

						Validators: []validator.String{
							stringvalidator.OneOf(utils.EnumSliceToStringSlice(mfa.AllowedEnumTimeUnitEnumValues)...),
						},
					},
				},
			},

			"user_verification": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A single nested object that specifies user verification settings, used to control whether the user must perform a gesture (such as a public key credential, fingerprint scan, or a PIN code) when registering or authenticating with their FIDO device.").Description,
				Required:    true,

				Attributes: map[string]schema.Attribute{
					"enforce_during_authentication": schema.BoolAttribute{
						Description:         userVerificationEnforceDuringAuthnDescription.Description,
						MarkdownDescription: userVerificationEnforceDuringAuthnDescription.MarkdownDescription,
						Required:            true,
					},

					"option": schema.StringAttribute{
						Description:         userVerificationOptionDescription.Description,
						MarkdownDescription: userVerificationOptionDescription.MarkdownDescription,
						Required:            true,

						Validators: []validator.String{
							stringvalidator.OneOf(utils.EnumSliceToStringSlice(mfa.AllowedEnumFIDO2PolicyUserVerificationOptionEnumValues)...),
						},
					},
				},
			},
		},
	}
}

func (r *FIDO2PolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *FIDO2PolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state FIDO2PolicyResourceModel

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
	fido2Policy, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *mfa.FIDO2Policy
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.MFAAPIClient.FIDO2PolicyApi.CreateFIDO2Policy(ctx, plan.EnvironmentId.ValueString()).FIDO2Policy(*fido2Policy).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"CreateFIDO2Policy",
		legacysdk.DefaultCustomError,
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

func (r *FIDO2PolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *FIDO2PolicyResourceModel

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
	var response *mfa.FIDO2Policy
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.MFAAPIClient.FIDO2PolicyApi.ReadOneFIDO2Policy(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadOneFIDO2Policy",
		legacysdk.CustomErrorResourceNotFoundWarning,
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

func (r *FIDO2PolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state FIDO2PolicyResourceModel

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
	fido2Policy, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *mfa.FIDO2Policy
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.MFAAPIClient.FIDO2PolicyApi.UpdateFIDO2Policy(ctx, plan.EnvironmentId.ValueString(), plan.Id.ValueString()).FIDO2Policy(*fido2Policy).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateFIDO2Policy",
		legacysdk.DefaultCustomError,
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

func (r *FIDO2PolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *FIDO2PolicyResourceModel

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
	resp.Diagnostics.Append(legacysdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fR, fErr := r.Client.MFAAPIClient.FIDO2PolicyApi.DeleteFIDO2Policy(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return legacysdk.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, fR, fErr)
		},
		"DeleteFIDO2Policy",
		mfaFido2PolicyDeleteCustomError,
		nil,
		nil,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

var mfaFido2PolicyDeleteCustomError = func(r *http.Response, p1Error *model.P1Error) diag.Diagnostics {
	var diags diag.Diagnostics

	if p1Error != nil {
		// Undeletable default FIDO2 policy
		if v, ok := p1Error.GetDetailsOk(); ok && v != nil && len(v) > 0 {
			if v[0].GetCode() == "CONSTRAINT_VIOLATION" {
				if match, _ := regexp.MatchString("cannot delete the default policy", v[0].GetMessage()); match {

					diags.AddWarning("Cannot delete the default MFA FIDO2 policy", "Due to API restrictions, the provider cannot delete the default FIDO2 policy for an environment.  The policy has been removed from Terraform state but has been left in place in the PingOne service.")

					return diags
				}
			}
		}
	}

	diags.Append(legacysdk.CustomErrorResourceNotFoundWarning(r, p1Error)...)
	return diags
}

func (r *FIDO2PolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:     "fido2_policy_id",
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

func (p *FIDO2PolicyResourceModel) expand(ctx context.Context) (*mfa.FIDO2Policy, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Backup eligibility
	var backupEligibilityPlan FIDO2PolicyBackupEligibilityResourceModel
	diags.Append(p.BackupEligibility.As(ctx, &backupEligibilityPlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return nil, diags
	}
	backupEligibility := mfa.NewFIDO2PolicyBackupEligibility(
		backupEligibilityPlan.Allow.ValueBool(),
		backupEligibilityPlan.EnforceDuringAuthentication.ValueBool(),
	)

	// MDS Authenticator Requirements
	var mdsAuthenticatorRequirementsPlan FIDO2PolicyMdsAuthenticatorsRequirementsResourceModel
	diags.Append(p.MdsAuthenticatorsRequirements.As(ctx, &mdsAuthenticatorRequirementsPlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return nil, diags
	}

	mdsAuthenticatorsRequirements := mfa.NewFIDO2PolicyMdsAuthenticatorsRequirements(
		mdsAuthenticatorRequirementsPlan.EnforceDuringAuthentication.ValueBool(),
		mfa.EnumFIDO2PolicyMDSAuthenticatorOption(mdsAuthenticatorRequirementsPlan.Option.ValueString()),
	)

	if !mdsAuthenticatorRequirementsPlan.AllowedAuthenticatorIDs.IsNull() && !mdsAuthenticatorRequirementsPlan.AllowedAuthenticatorIDs.IsUnknown() {
		allowedAuthenticators := make([]mfa.FIDO2PolicyMdsAuthenticatorsRequirementsAllowedAuthenticatorsInner, 0)

		var allowedAuthenticatorIDsPlan []types.String
		diags.Append(mdsAuthenticatorRequirementsPlan.AllowedAuthenticatorIDs.ElementsAs(ctx, &allowedAuthenticatorIDsPlan, false)...)
		if diags.HasError() {
			return nil, diags
		}

		allowedAuthenticatorIDs, d := framework.TFTypeStringSliceToStringSlice(allowedAuthenticatorIDsPlan, path.Root("mds_authenticators_requirements").AtName("allowed_authenticator_ids"))
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		for _, allowedAuthenticatorIDPlan := range allowedAuthenticatorIDs {

			allowedAuthenticator := *mfa.NewFIDO2PolicyMdsAuthenticatorsRequirementsAllowedAuthenticatorsInner(
				allowedAuthenticatorIDPlan,
			)

			allowedAuthenticators = append(allowedAuthenticators, allowedAuthenticator)
		}

		mdsAuthenticatorsRequirements.SetAllowedAuthenticators(allowedAuthenticators)
	}

	// User display name attributes
	var userDisplayNameAttributesPlan FIDO2PolicyUserDisplayNameAttributesResourceModel
	diags.Append(p.UserDisplayNameAttributes.As(ctx, &userDisplayNameAttributesPlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return nil, diags
	}

	attributes := make([]mfa.FIDO2PolicyUserDisplayNameAttributesAttributesInner, 0)
	if !userDisplayNameAttributesPlan.Attributes.IsNull() && !userDisplayNameAttributesPlan.Attributes.IsUnknown() {

		var userDisplayNameAttributesAttributesPlan []FIDO2PolicyUserDisplayNameAttributesAttributesResourceModel
		diags.Append(userDisplayNameAttributesPlan.Attributes.ElementsAs(ctx, &userDisplayNameAttributesAttributesPlan, false)...)
		if diags.HasError() {
			return nil, diags
		}

		for _, attributePlan := range userDisplayNameAttributesAttributesPlan {

			attribute := *mfa.NewFIDO2PolicyUserDisplayNameAttributesAttributesInner(
				attributePlan.Name.ValueString(),
			)

			if !attributePlan.SubAttributes.IsNull() && !attributePlan.SubAttributes.IsUnknown() {
				var userDisplayNameAttributesAttributesSubAttributesPlan []FIDO2PolicyUserDisplayNameAttributesAttributesSubAttributesResourceModel
				diags.Append(attributePlan.SubAttributes.ElementsAs(ctx, &userDisplayNameAttributesAttributesSubAttributesPlan, false)...)
				if diags.HasError() {
					return nil, diags
				}

				subAttributes := make([]mfa.FIDO2PolicyUserDisplayNameAttributesAttributesInnerSubAttributesInner, 0)

				for _, subAttributePlan := range userDisplayNameAttributesAttributesSubAttributesPlan {
					subAttributes = append(subAttributes, *mfa.NewFIDO2PolicyUserDisplayNameAttributesAttributesInnerSubAttributesInner(
						subAttributePlan.Name.ValueString(),
					))
				}

				attribute.SetSubAttributes(subAttributes)
			}

			attributes = append(attributes, attribute)
		}
	}

	userDisplayNameAttributes := mfa.NewFIDO2PolicyUserDisplayNameAttributes(
		attributes,
	)

	// User verification
	var userVerificationPlan FIDO2PolicyUserVerificationResourceModel
	diags.Append(p.UserVerification.As(ctx, &userVerificationPlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return nil, diags
	}
	userVerification := mfa.NewFIDO2PolicyUserVerification(
		userVerificationPlan.EnforceDuringAuthentication.ValueBool(),
		mfa.EnumFIDO2PolicyUserVerificationOption(userVerificationPlan.Option.ValueString()),
	)

	// Main object
	data := mfa.NewFIDO2Policy(
		mfa.EnumFIDO2PolicyAttestationRequirements(p.AttestationRequirements.ValueString()),
		mfa.EnumFIDO2PolicyAuthenticatorAttachment(p.AuthenticatorAttachment.ValueString()),
		*backupEligibility,
		p.DeviceDisplayName.ValueString(),
		mfa.EnumFIDO2PolicyDiscoverableCredentials(p.DiscoverableCredentials.ValueString()),
		*mdsAuthenticatorsRequirements,
		p.Name.ValueString(),
		p.RelyingPartyId.ValueString(),
		*userDisplayNameAttributes,
		*userVerification,
	)

	// User presence timeout
	if !p.UserPresenceTimeout.IsNull() && !p.UserPresenceTimeout.IsUnknown() {
		var userPresenceTimeoutPlan FIDO2PolicyUserPresenceTimeoutResourceModel
		diags.Append(p.UserPresenceTimeout.As(ctx, &userPresenceTimeoutPlan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if diags.HasError() {
			return nil, diags
		}

		userPresenceTimeout := mfa.NewFIDO2PolicyUserPresenceTimeout()
		userPresenceTimeout.SetDuration(userPresenceTimeoutPlan.Duration.ValueInt32())
		userPresenceTimeout.SetTimeUnit(mfa.EnumTimeUnit(userPresenceTimeoutPlan.TimeUnit.ValueString()))
		data.SetUserPresenceTimeout(*userPresenceTimeout)
	}

	if !p.Description.IsNull() && !p.Description.IsUnknown() {
		data.SetDescription(p.Description.ValueString())
	}

	if !p.Default.IsNull() && !p.Default.IsUnknown() {
		data.SetDefault(p.Default.ValueBool())
	} else {
		data.SetDefault(false)
	}

	return data, diags
}

func (p *FIDO2PolicyResourceModel) toState(apiObject *mfa.FIDO2Policy) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)
		return diags
	}

	var d diag.Diagnostics

	p.Id = framework.PingOneResourceIDOkToTF(apiObject.GetIdOk())
	p.EnvironmentId = framework.PingOneResourceIDToTF(*apiObject.GetEnvironment().Id)
	p.Name = framework.StringOkToTF(apiObject.GetNameOk())
	p.Description = framework.StringOkToTF(apiObject.GetDescriptionOk())
	p.Default = framework.BoolOkToTF(apiObject.GetDefaultOk())
	p.AttestationRequirements = framework.EnumOkToTF(apiObject.GetAttestationRequirementsOk())
	p.AuthenticatorAttachment = framework.EnumOkToTF(apiObject.GetAuthenticatorAttachmentOk())

	p.BackupEligibility, d = toStateBackupEligibility(apiObject.GetBackupEligibilityOk())
	diags.Append(d...)

	p.DeviceDisplayName = framework.StringOkToTF(apiObject.GetDeviceDisplayNameOk())
	p.DiscoverableCredentials = framework.EnumOkToTF(apiObject.GetDiscoverableCredentialsOk())

	p.MdsAuthenticatorsRequirements, d = toStateMdsAuthenticatorsRequirements(apiObject.GetMdsAuthenticatorsRequirementsOk())
	diags.Append(d...)

	p.RelyingPartyId = framework.StringOkToTF(apiObject.GetRelyingPartyIdOk())

	p.UserDisplayNameAttributes, d = toStateUserDisplayNameAttributes(apiObject.GetUserDisplayNameAttributesOk())
	diags.Append(d...)

	p.UserPresenceTimeout, d = toStateUserPresenceTimeout(apiObject.GetUserPresenceTimeoutOk())
	diags.Append(d...)

	p.UserVerification, d = toStateUserVerification(apiObject.GetUserVerificationOk())
	diags.Append(d...)

	return diags
}

func toStateBackupEligibility(apiObject *mfa.FIDO2PolicyBackupEligibility, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(fido2PolicyBackupEligibilityTFObjectTypes), nil
	}

	o := map[string]attr.Value{
		"allow":                         framework.BoolOkToTF(apiObject.GetAllowOk()),
		"enforce_during_authentication": framework.BoolOkToTF(apiObject.GetEnforceDuringAuthenticationOk()),
	}

	objValue, d := types.ObjectValue(fido2PolicyBackupEligibilityTFObjectTypes, o)
	diags.Append(d...)

	return objValue, diags
}

func toStateMdsAuthenticatorsRequirements(apiObject *mfa.FIDO2PolicyMdsAuthenticatorsRequirements, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(fido2PolicyMdsAuthenticatorRequirementsTFObjectTypes), nil
	}

	o := map[string]attr.Value{
		"allowed_authenticator_ids":     types.SetNull(types.StringType),
		"enforce_during_authentication": framework.BoolOkToTF(apiObject.GetEnforceDuringAuthenticationOk()),
		"option":                        framework.EnumOkToTF(apiObject.GetOptionOk()),
	}

	if v, ok := apiObject.GetAllowedAuthenticatorsOk(); ok {
		allowedAuthenticatorsList := make([]string, 0)
		for _, item := range v {
			if id, ok := item.GetIdOk(); ok {
				allowedAuthenticatorsList = append(allowedAuthenticatorsList, *id)
			}
		}

		o["allowed_authenticator_ids"] = framework.StringSetToTF(allowedAuthenticatorsList)
	}

	objValue, d := types.ObjectValue(fido2PolicyMdsAuthenticatorRequirementsTFObjectTypes, o)
	diags.Append(d...)

	return objValue, diags
}

func toStateUserDisplayNameAttributes(apiObject *mfa.FIDO2PolicyUserDisplayNameAttributes, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(fido2PolicyUserDisplayNameAttributesTFObjectTypes), nil
	}

	attributesList, d := toStateUserDisplayNameAttributesAttributes(apiObject.GetAttributesOk())
	diags.Append(d...)

	o := map[string]attr.Value{
		"attributes": attributesList,
	}

	objValue, d := types.ObjectValue(fido2PolicyUserDisplayNameAttributesTFObjectTypes, o)
	diags.Append(d...)

	return objValue, diags
}

func toStateUserDisplayNameAttributesAttributes(apiObject []mfa.FIDO2PolicyUserDisplayNameAttributesAttributesInner, ok bool) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	tfObjType := types.ObjectType{AttrTypes: fido2PolicyUserDisplayNameAttributesAttributesTFObjectTypes}

	if !ok || len(apiObject) == 0 {
		return types.ListNull(tfObjType), diags
	}

	flattenedList := []attr.Value{}
	for _, v := range apiObject {

		subAttributesList, d := toStateUserDisplayNameAttributesAttributesSubAttributes(v.GetSubAttributesOk())
		diags.Append(d...)

		objMap := map[string]attr.Value{
			"sub_attributes": subAttributesList,
			"name":           framework.StringOkToTF(v.GetNameOk()),
		}

		flattenedObj, d := types.ObjectValue(fido2PolicyUserDisplayNameAttributesAttributesTFObjectTypes, objMap)
		diags.Append(d...)

		flattenedList = append(flattenedList, flattenedObj)
	}

	returnVar, d := types.ListValue(tfObjType, flattenedList)
	diags.Append(d...)

	return returnVar, diags
}

func toStateUserDisplayNameAttributesAttributesSubAttributes(apiObject []mfa.FIDO2PolicyUserDisplayNameAttributesAttributesInnerSubAttributesInner, ok bool) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	tfObjType := types.ObjectType{AttrTypes: fido2PolicyUserDisplayNameAttributesAttributesSubAttributesTFObjectTypes}

	if !ok || len(apiObject) == 0 {
		return types.ListNull(tfObjType), diags
	}

	flattenedList := []attr.Value{}
	for _, v := range apiObject {

		objMap := map[string]attr.Value{
			"name": framework.StringOkToTF(v.GetNameOk()),
		}

		flattenedObj, d := types.ObjectValue(fido2PolicyUserDisplayNameAttributesAttributesSubAttributesTFObjectTypes, objMap)
		diags.Append(d...)

		flattenedList = append(flattenedList, flattenedObj)
	}

	returnVar, d := types.ListValue(tfObjType, flattenedList)
	diags.Append(d...)

	return returnVar, diags
}

func toStateUserPresenceTimeout(apiObject *mfa.FIDO2PolicyUserPresenceTimeout, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(fido2PolicyUserPresenceTimeoutTFObjectTypes), nil
	}

	o := map[string]attr.Value{
		"duration":  framework.Int32OkToTF(apiObject.GetDurationOk()),
		"time_unit": framework.EnumOkToTF(apiObject.GetTimeUnitOk()),
	}

	objValue, d := types.ObjectValue(fido2PolicyUserPresenceTimeoutTFObjectTypes, o)
	diags.Append(d...)

	return objValue, diags
}

func toStateUserVerification(apiObject *mfa.FIDO2PolicyUserVerification, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(fido2PolicyUserVerificationTFObjectTypes), nil
	}

	o := map[string]attr.Value{
		"enforce_during_authentication": framework.BoolOkToTF(apiObject.GetEnforceDuringAuthenticationOk()),
		"option":                        framework.EnumOkToTF(apiObject.GetOptionOk()),
	}

	objValue, d := types.ObjectValue(fido2PolicyUserVerificationTFObjectTypes, o)
	diags.Append(d...)

	return objValue, diags
}
