package sso

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/service"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type IdentityProviderResource serviceClientType

type IdentityProviderResourceModel struct {
	Id                       types.String `tfsdk:"id"`
	EnvironmentId            types.String `tfsdk:"environment_id"`
	Name                     types.String `tfsdk:"name"`
	Description              types.String `tfsdk:"description"`
	Enabled                  types.Bool   `tfsdk:"enabled"`
	RegistrationPopulationId types.String `tfsdk:"registration_population_id"`
	LoginButtonIcon          types.Object `tfsdk:"login_button_icon"`
	Icon                     types.Object `tfsdk:"icon"`
	Facebook                 types.List   `tfsdk:"facebook"`
	Google                   types.List   `tfsdk:"google"`
	LinkedIn                 types.List   `tfsdk:"linkedin"`
	Yahoo                    types.List   `tfsdk:"yahoo"`
	Amazon                   types.List   `tfsdk:"amazon"`
	Twitter                  types.List   `tfsdk:"twitter"`
	Apple                    types.List   `tfsdk:"apple"`
	Paypal                   types.List   `tfsdk:"paypal"`
	Microsoft                types.List   `tfsdk:"microsoft"`
	Github                   types.List   `tfsdk:"github"`
	OpenIDConnect            types.List   `tfsdk:"openid_connect"`
	Saml                     types.List   `tfsdk:"saml"`
}

type IdentityProviderClientIdClientSecretResourceModel struct {
	ClientId     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
}

type IdentityProviderLoginButtonIcon service.ImageResourceModel

type IdentityProviderIcon service.ImageResourceModel

type IdentityProviderFacebookResourceModel struct {
	AppId     types.String `tfsdk:"app_id"`
	AppSecret types.String `tfsdk:"app_secret"`
}

type IdentityProviderGoogleResourceModel IdentityProviderClientIdClientSecretResourceModel

type IdentityProviderLinkedInResourceModel IdentityProviderClientIdClientSecretResourceModel

type IdentityProviderYahooResourceModel IdentityProviderClientIdClientSecretResourceModel

type IdentityProviderAmazonResourceModel IdentityProviderClientIdClientSecretResourceModel

type IdentityProviderTwitterResourceModel IdentityProviderClientIdClientSecretResourceModel

type IdentityProviderAppleResourceModel struct {
	TeamId                 types.String `tfsdk:"team_id"`
	KeyId                  types.String `tfsdk:"key_id"`
	ClientId               types.String `tfsdk:"client_id"`
	ClientSecretSigningKey types.String `tfsdk:"client_secret_signing_key"`
}

type IdentityProviderPaypalResourceModel struct {
	ClientId          types.String `tfsdk:"client_id"`
	ClientSecret      types.String `tfsdk:"client_secret"`
	ClientEnvironment types.String `tfsdk:"client_environment"`
}

type IdentityProviderMicrosoftResourceModel IdentityProviderClientIdClientSecretResourceModel

type IdentityProviderGithubResourceModel IdentityProviderClientIdClientSecretResourceModel

type IdentityProviderOIDCResourceModel struct {
	AuthorizationEndpoint   types.String `tfsdk:"authorization_endpoint"`
	ClientId                types.String `tfsdk:"client_id"`
	ClientSecret            types.String `tfsdk:"client_secret"`
	DiscoveryEndpoint       types.String `tfsdk:"discovery_endpoint"`
	Issuer                  types.String `tfsdk:"issuer"`
	JwksEndpoint            types.String `tfsdk:"jwks_endpoint"`
	Scopes                  types.Set    `tfsdk:"scopes"`
	TokenEndpoint           types.String `tfsdk:"token_endpoint"`
	TokenEndpointAuthMethod types.String `tfsdk:"token_endpoint_auth_method"`
	UserinfoEndpoint        types.String `tfsdk:"userinfo_endpoint"`
}

type IdentityProviderSAMLResourceModel struct {
	AuthenticationRequestSigned   types.Bool   `tfsdk:"authentication_request_signed"`
	IdpEntityId                   types.String `tfsdk:"idp_entity_id"`
	SpEntityId                    types.String `tfsdk:"sp_entity_id"`
	IdpVerificationCertificateIds types.Set    `tfsdk:"idp_verification_certificate_ids"`
	SpSigningKeyId                types.String `tfsdk:"sp_signing_key_id"`
	SsoBinding                    types.String `tfsdk:"sso_binding"`
	SsoEndpoint                   types.String `tfsdk:"sso_endpoint"`
	SloBinding                    types.String `tfsdk:"slo_binding"`
	SloEndpoint                   types.String `tfsdk:"slo_endpoint"`
	SloResponseEndpoint           types.String `tfsdk:"slo_response_endpoint"`
	SloWindow                     types.Int64  `tfsdk:"slo_window"`
}

var (
	identityProviderFacebookTFObjectTypes = map[string]attr.Type{
		"app_id":     types.StringType,
		"app_secret": types.StringType,
	}

	identityProviderClientIDClientSecretTFObjectTypes = map[string]attr.Type{
		"client_id":     types.StringType,
		"client_secret": types.StringType,
	}

	identityProviderAppleTFObjectTypes = map[string]attr.Type{
		"team_id":                   types.StringType,
		"key_id":                    types.StringType,
		"client_id":                 types.StringType,
		"client_secret_signing_key": types.StringType,
	}

	identityProviderPaypalTFObjectTypes = map[string]attr.Type{
		"client_id":          types.StringType,
		"client_secret":      types.StringType,
		"client_environment": types.StringType,
	}

	identityProviderOIDCTFObjectTypes = map[string]attr.Type{
		"authorization_endpoint":     types.StringType,
		"client_id":                  types.StringType,
		"client_secret":              types.StringType,
		"discovery_endpoint":         types.StringType,
		"issuer":                     types.StringType,
		"jwks_endpoint":              types.StringType,
		"scopes":                     types.SetType{ElemType: types.StringType},
		"token_endpoint":             types.StringType,
		"token_endpoint_auth_method": types.StringType,
		"userinfo_endpoint":          types.StringType,
	}

	identityProviderSAMLTFObjectTypes = map[string]attr.Type{
		"authentication_request_signed":    types.BoolType,
		"idp_entity_id":                    types.StringType,
		"sp_entity_id":                     types.StringType,
		"idp_verification_certificate_ids": types.SetType{ElemType: types.StringType},
		"sp_signing_key_id":                types.StringType,
		"sso_binding":                      types.StringType,
		"sso_endpoint":                     types.StringType,
		"slo_binding":                      types.StringType,
		"slo_endpoint":                     types.StringType,
		"slo_response_endpoint":            types.StringType,
		"slo_window":                       types.Int64Type,
	}
)

// Framework interfaces
var (
	_ resource.Resource                = &IdentityProviderResource{}
	_ resource.ResourceWithConfigure   = &IdentityProviderResource{}
	_ resource.ResourceWithImportState = &IdentityProviderResource{}
)

// New Object
func NewIdentityProviderResource() resource.Resource {
	return &IdentityProviderResource{}
}

// Metadata
func (r *IdentityProviderResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_identity_provider"
}

// Schema.
func (r *IdentityProviderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	const attrMinLength = 1
	const appleKeyIdLength = 10
	const appleTeamIdLength = 10
	const samlSloWindowMin = 1
	const samlSloWindowMax = 24

	providerAttributeList := []string{"facebook", "google", "linkedin", "yahoo", "amazon", "twitter", "apple", "paypal", "microsoft", "github", "openid_connect", "saml"}

	enabledDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether the identity provider is enabled in the environment.",
	).DefaultValue(false)

	registrationPopulationIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the population ID to create users in, when using just-in-time provisioning. Setting this attribute gives management of linked users to the IdP and also triggers just-in-time provisioning of new users to the population ID provided.",
	).AppendMarkdownString("Must be a valid PingOne resource ID.")

	loginButtonIconIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The ID for the identity provider icon to use as the login button.  This can be retrieved from the `id` parameter of the `pingone_image` resource.  Must be a valid PingOne resource ID.",
	)

	loginButtonIconHrefDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The URL or fully qualified path to the identity provider icon to use as the login button.  This can be retrieved from the `uploaded_image[0].href` parameter of the `pingone_image` resource.",
	)

	iconIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The ID for the identity provider icon to use as the login button.  This can be retrieved from the `id` parameter of the `pingone_image` resource.  Must be a valid PingOne resource ID.",
	)

	iconHrefDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The URL or fully qualified path to the identity provider icon to use as the login button.  This can be retrieved from the `uploaded_image[0].href` parameter of the `pingone_image` resource.",
	)

	paypalClientEnvironmentDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the PayPal environment.",
	).AllowedValues("sandbox", "live")

	oidcTokenEndpointAuthMethodDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the OIDC identity provider's token endpoint authentication method.",
	).AllowedValuesEnum(management.AllowedEnumIdentityProviderOIDCTokenAuthMethodEnumValues).DefaultValue(string(management.ENUMIDENTITYPROVIDEROIDCTOKENAUTHMETHOD_CLIENT_SECRET_BASIC))

	samlAuthenticationRequestSignedDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether the SAML authentication request will be signed when sending to the identity provider. Set this to `true` if the external IDP is included in an authentication policy to be used by applications that are accessed using a mix of default URLS and custom Domains URLs.",
	).DefaultValue(false)

	samlIdpEntityIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the entity ID URI that is checked against the `issuerId` tag in the incoming response.",
	)

	samlSSOBindingDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the binding for the authentication request.",
	).AllowedValuesEnum(management.AllowedEnumIdentityProviderSAMLSSOBindingEnumValues)

	samlSLOBindingDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the binding protocol to be used for the logout response.",
	).AllowedValuesEnum(management.AllowedEnumIdentityProviderSAMLSLOBindingEnumValues).DefaultValue(string(management.ENUMIDENTITYPROVIDERSAMLSLOBINDING_POST))

	samlSLOResponseEndpointDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the endpoint URL to submit the logout response.  If a value is not provided, the `slo_endpoint` property value is used to submit SLO response.  This value must be a URL that uses http or https.",
	)

	samlSLOWindowDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("An integer that defines how long (hours) PingOne can exchange logout messages with the application, specifically a logout request from the application, since the initial request. The minimum value is `%d` hour and the maximum is `%d` hours.", samlSloWindowMin, samlSloWindowMax),
	)

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage Identity Providers in a PingOne environment.",

		Attributes: map[string]schema.Attribute{
			"id": framework.Attr_ID(),

			"environment_id": framework.Attr_LinkID(
				framework.SchemaAttributeDescriptionFromMarkdown("The ID of the environment to create the identity provider in."),
			),

			"name": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the name of the identity provider.").Description,
				Required:    true,

				Validators: []validator.String{
					stringvalidator.LengthAtLeast(attrMinLength),
				},
			},

			"description": schema.StringAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the description of the identity provider.").Description,
				Optional:    true,
			},

			"enabled": schema.BoolAttribute{
				Description:         enabledDescription.Description,
				MarkdownDescription: enabledDescription.MarkdownDescription,
				Optional:            true,
				Computed:            true,

				Default: booldefault.StaticBool(false),
			},

			"registration_population_id": schema.StringAttribute{
				Description:         registrationPopulationIdDescription.Description,
				MarkdownDescription: registrationPopulationIdDescription.MarkdownDescription,
				Optional:            true,

				Validators: []validator.String{
					verify.P1ResourceIDValidator(),
				},
			},

			"login_button_icon": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that specifies the HREF and ID for the identity provider icon to use in the login button.").Description,
				Optional:    true,

				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description:         loginButtonIconIdDescription.Description,
						MarkdownDescription: loginButtonIconIdDescription.MarkdownDescription,
						Required:            true,

						Validators: []validator.String{
							verify.P1ResourceIDValidator(),
						},
					},

					"href": schema.StringAttribute{
						Description:         loginButtonIconHrefDescription.Description,
						MarkdownDescription: loginButtonIconHrefDescription.MarkdownDescription,
						Required:            true,

						Validators: []validator.String{
							stringvalidator.RegexMatches(verify.IsURLWithHTTPS, "Value must be a valid URL with `https://` prefix."),
						},
					},
				},
			},

			"icon": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that specifies the HREF and ID for the identity provider icon.").Description,
				Optional:    true,

				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description:         iconIdDescription.Description,
						MarkdownDescription: iconIdDescription.MarkdownDescription,
						Required:            true,

						Validators: []validator.String{
							verify.P1ResourceIDValidator(),
						},
					},

					"href": schema.StringAttribute{
						Description:         iconHrefDescription.Description,
						MarkdownDescription: iconHrefDescription.MarkdownDescription,
						Required:            true,

						Validators: []validator.String{
							stringvalidator.RegexMatches(verify.IsURLWithHTTPS, "Value must be a valid URL with `https://` prefix."),
						},
					},
				},
			},
		},

		Blocks: map[string]schema.Block{

			// The providers
			"facebook": identityProviderSchemaBlock(
				framework.SchemaAttributeDescriptionFromMarkdown("A single block that specifies options for connectivity to the Facebook social identity provider."),

				map[string]schema.Attribute{
					"app_id": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the application ID from Facebook.").Description,
						Required:    true,

						Validators: []validator.String{
							stringvalidator.LengthAtLeast(attrMinLength),
						},
					},

					"app_secret": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the application secret from Facebook.").Description,
						Required:    true,
						Sensitive:   true,

						Validators: []validator.String{
							stringvalidator.LengthAtLeast(attrMinLength),
						},
					},
				},

				providerAttributeList,
			),

			"google": identityProviderSchemaBlock(
				framework.SchemaAttributeDescriptionFromMarkdown("A single block that specifies options for connectivity to the Google social identity provider."),

				identityProviderClientIdClientSecretAttributes("Google"),

				providerAttributeList,
			),

			"linkedin": identityProviderSchemaBlock(
				framework.SchemaAttributeDescriptionFromMarkdown("A single block that specifies options for connectivity to the LinkedIn social identity provider."),

				identityProviderClientIdClientSecretAttributes("LinkedIn"),

				providerAttributeList,
			),

			"yahoo": identityProviderSchemaBlock(
				framework.SchemaAttributeDescriptionFromMarkdown("A single block that specifies options for connectivity to the Yahoo social identity provider."),

				identityProviderClientIdClientSecretAttributes("Yahoo"),

				providerAttributeList,
			),

			"amazon": identityProviderSchemaBlock(
				framework.SchemaAttributeDescriptionFromMarkdown("A single block that specifies options for connectivity to the Amazon social identity provider."),

				identityProviderClientIdClientSecretAttributes("Amazon"),

				providerAttributeList,
			),

			"twitter": identityProviderSchemaBlock(
				framework.SchemaAttributeDescriptionFromMarkdown("A single block that specifies options for connectivity to the Twitter social identity provider."),

				identityProviderClientIdClientSecretAttributes("Twitter"),

				providerAttributeList,
			),

			"apple": identityProviderSchemaBlock(
				framework.SchemaAttributeDescriptionFromMarkdown("A single block that specifies options for connectivity to the Apple social identity provider."),

				map[string]schema.Attribute{
					"client_id": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the application ID from Apple. This is the identifier obtained after registering a services ID in the Apple developer portal.").Description,
						Required:    true,

						Validators: []validator.String{
							stringvalidator.LengthAtLeast(attrMinLength),
						},
					},

					"client_secret_signing_key": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the private key that is used to generate a client secret.").Description,
						Required:    true,
						Sensitive:   true,

						Validators: []validator.String{
							stringvalidator.LengthAtLeast(attrMinLength),
						},
					},

					"key_id": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A 10-character string that Apple uses to identify an authentication key.").Description,
						Required:    true,

						Validators: []validator.String{
							stringvalidator.LengthBetween(appleKeyIdLength, appleKeyIdLength),
						},
					},

					"team_id": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A 10-character string that Apple uses to identify teams.").Description,
						Required:    true,

						Validators: []validator.String{
							stringvalidator.LengthBetween(appleTeamIdLength, appleTeamIdLength),
						},
					},
				},

				providerAttributeList,
			),

			"paypal": identityProviderSchemaBlock(
				framework.SchemaAttributeDescriptionFromMarkdown("A single block that specifies options for connectivity to the Paypal social identity provider."),

				map[string]schema.Attribute{
					"client_id": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the application ID from Paypal.").Description,
						Required:    true,

						Validators: []validator.String{
							stringvalidator.LengthAtLeast(attrMinLength),
						},
					},

					"client_secret": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the application secret from PayPal.").Description,
						Required:    true,
						Sensitive:   true,

						Validators: []validator.String{
							stringvalidator.LengthAtLeast(attrMinLength),
						},
					},

					"client_environment": schema.StringAttribute{
						Description:         paypalClientEnvironmentDescription.Description,
						MarkdownDescription: paypalClientEnvironmentDescription.MarkdownDescription,
						Required:            true,

						Validators: []validator.String{
							stringvalidator.OneOf("sandbox", "live"),
						},
					},
				},

				providerAttributeList,
			),

			"microsoft": identityProviderSchemaBlock(
				framework.SchemaAttributeDescriptionFromMarkdown("A single block that specifies options for connectivity to the Microsoft social identity provider."),

				identityProviderClientIdClientSecretAttributes("Microsoft"),

				providerAttributeList,
			),

			"github": identityProviderSchemaBlock(
				framework.SchemaAttributeDescriptionFromMarkdown("A single block that specifies options for connectivity to the Github social identity provider."),

				identityProviderClientIdClientSecretAttributes("Github"),

				providerAttributeList,
			),

			"openid_connect": identityProviderSchemaBlock(
				framework.SchemaAttributeDescriptionFromMarkdown("A single block that specifies options for connectivity to an OpenID Connect compliant identity provider."),

				map[string]schema.Attribute{
					"authorization_endpoint": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the the OIDC identity provider's authorization endpoint. This value must be a URL that uses https.").Description,
						Required:    true,

						Validators: []validator.String{
							stringvalidator.RegexMatches(verify.IsURLWithHTTPS, "Value must be a valid URL with `https://` prefix."),
						},
					},

					"client_id": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the application client ID from the OIDC identity provider.").Description,
						Required:    true,

						Validators: []validator.String{
							stringvalidator.LengthAtLeast(attrMinLength),
						},
					},

					"client_secret": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the application client secret from the OIDC identity provider.").Description,
						Required:    true,
						Sensitive:   true,

						Validators: []validator.String{
							stringvalidator.LengthAtLeast(attrMinLength),
						},
					},

					"discovery_endpoint": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the OIDC identity provider's discovery endpoint. This value must be a URL that uses https.").Description,
						Optional:    true,

						Validators: []validator.String{
							stringvalidator.RegexMatches(verify.IsURLWithHTTPS, "Value must be a valid URL with `https://` prefix."),
						},
					},

					"issuer": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the issuer to which the authentication is sent for the OIDC identity provider. This value must be a URL that uses https.").Description,
						Required:    true,

						Validators: []validator.String{
							stringvalidator.RegexMatches(verify.IsURLWithHTTPS, "Value must be a valid URL with `https://` prefix."),
						},
					},

					"jwks_endpoint": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the OIDC identity provider's jwks endpoint. This value must be a URL that uses https.").Description,
						Required:    true,

						Validators: []validator.String{
							stringvalidator.RegexMatches(verify.IsURLWithHTTPS, "Value must be a valid URL with `https://` prefix."),
						},
					},

					"scopes": schema.SetAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("An array that specifies the scopes to include in the authentication request to the OIDC identity provider.").Description,
						Required:    true,

						ElementType: types.StringType,

						Validators: []validator.Set{
							setvalidator.SizeAtLeast(attrMinLength),
							setvalidator.ValueStringsAre(
								stringvalidator.LengthAtLeast(attrMinLength),
							),
						},
					},

					"token_endpoint": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the OIDC identity provider's token endpoint. This value must be a URL that uses https.").Description,
						Required:    true,

						Validators: []validator.String{
							stringvalidator.RegexMatches(verify.IsURLWithHTTPS, "Value must be a valid URL with `https://` prefix."),
						},
					},

					"token_endpoint_auth_method": schema.StringAttribute{
						Description:         oidcTokenEndpointAuthMethodDescription.Description,
						MarkdownDescription: oidcTokenEndpointAuthMethodDescription.MarkdownDescription,
						Optional:            true,
						Computed:            true,

						Default: stringdefault.StaticString(string(management.ENUMIDENTITYPROVIDEROIDCTOKENAUTHMETHOD_CLIENT_SECRET_BASIC)),

						Validators: []validator.String{
							stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumIdentityProviderOIDCTokenAuthMethodEnumValues)...),
						},
					},

					"userinfo_endpoint": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the OIDC identity provider's userInfo endpoint. This value must be a URL that uses https.").Description,
						Optional:    true,

						Validators: []validator.String{
							stringvalidator.RegexMatches(verify.IsURLWithHTTPS, "Value must be a valid URL with `https://` prefix."),
						},
					},
				},

				providerAttributeList,
			),

			"saml": identityProviderSchemaBlock(
				framework.SchemaAttributeDescriptionFromMarkdown("A single block that specifies options for connectivity to a SAML 2.0 compliant identity provider."),

				map[string]schema.Attribute{
					"authentication_request_signed": schema.BoolAttribute{
						Description:         samlAuthenticationRequestSignedDescription.Description,
						MarkdownDescription: samlAuthenticationRequestSignedDescription.MarkdownDescription,
						Optional:            true,
						Computed:            true,

						Default: booldefault.StaticBool(false),
					},

					"idp_entity_id": schema.StringAttribute{
						Description:         samlIdpEntityIdDescription.Description,
						MarkdownDescription: samlIdpEntityIdDescription.MarkdownDescription,
						Required:            true,

						Validators: []validator.String{
							stringvalidator.LengthAtLeast(attrMinLength),
						},
					},

					"sp_entity_id": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the service provider's entity ID, used to look up the application.").Description,
						Required:    true,

						Validators: []validator.String{
							stringvalidator.LengthAtLeast(attrMinLength),
						},
					},

					"idp_verification_certificate_ids": schema.SetAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("An unordered list that specifies the identity provider's certificate IDs used to verify the signature on the signed assertion from the identity provider. Signing is done with a private key and verified with a public key.  Items must be valid PingOne resource IDs.").Description,
						Required:    true,

						ElementType: types.StringType,

						Validators: []validator.Set{
							setvalidator.SizeAtLeast(attrMinLength),
							setvalidator.ValueStringsAre(
								verify.P1ResourceIDValidator(),
							),
						},
					},

					"sp_signing_key_id": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the service provider's signing key ID.  Must be a valid PingOne resource ID.").Description,
						Optional:    true,

						Validators: []validator.String{
							verify.P1ResourceIDValidator(),
						},
					},

					"sso_binding": schema.StringAttribute{
						Description:         samlSSOBindingDescription.Description,
						MarkdownDescription: samlSSOBindingDescription.MarkdownDescription,
						Required:            true,

						Validators: []validator.String{
							stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumIdentityProviderSAMLSSOBindingEnumValues)...),
						},
					},

					"sso_endpoint": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the SSO endpoint for the authentication request.  This value must be a URL that uses http or https.").Description,
						Required:    true,

						Validators: []validator.String{
							stringvalidator.RegexMatches(verify.IsURLWithHTTPorHTTPS, "Value must be a valid URL with `http://` or `https://` prefix."),
						},
					},

					"slo_binding": schema.StringAttribute{
						Description:         samlSLOBindingDescription.Description,
						MarkdownDescription: samlSLOBindingDescription.MarkdownDescription,
						Optional:            true,
						Computed:            true,

						Default: stringdefault.StaticString(string(management.ENUMIDENTITYPROVIDERSAMLSLOBINDING_POST)),

						Validators: []validator.String{
							stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumIdentityProviderSAMLSSOBindingEnumValues)...),
						},
					},

					"slo_endpoint": schema.StringAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the logout endpoint URL. This is an optional property. However, if a logout endpoint URL is not defined, logout actions result in an error.  This value must be a URL that uses http or https.").Description,
						Optional:    true,

						Validators: []validator.String{
							stringvalidator.RegexMatches(verify.IsURLWithHTTPorHTTPS, "Value must be a valid URL with `http://` or `https://` prefix."),
						},
					},

					"slo_response_endpoint": schema.StringAttribute{
						Description:         samlSLOResponseEndpointDescription.Description,
						MarkdownDescription: samlSLOResponseEndpointDescription.MarkdownDescription,
						Optional:            true,

						Validators: []validator.String{
							stringvalidator.RegexMatches(verify.IsURLWithHTTPorHTTPS, "Value must be a valid URL with `http://` or `https://` prefix."),
						},
					},

					"slo_window": schema.Int64Attribute{
						Description:         samlSLOWindowDescription.Description,
						MarkdownDescription: samlSLOWindowDescription.MarkdownDescription,
						Optional:            true,

						Validators: []validator.Int64{
							int64validator.Between(samlSloWindowMin, samlSloWindowMax),
						},
					},
				},

				providerAttributeList,
			),
		},
	}
}

func identityProviderSchemaBlock(description framework.SchemaAttributeDescription, attributes map[string]schema.Attribute, exactlyOneOfBlockNames []string) schema.ListNestedBlock {
	description = description.ExactlyOneOf(exactlyOneOfBlockNames).RequiresReplaceBlock()

	exactlyOneOfPaths := make([]path.Expression, len(exactlyOneOfBlockNames))
	for i, blockName := range exactlyOneOfBlockNames {
		exactlyOneOfPaths[i] = path.MatchRelative().AtParent().AtName(blockName)
	}

	return schema.ListNestedBlock{
		Description:         description.Description,
		MarkdownDescription: description.MarkdownDescription,

		NestedObject: schema.NestedBlockObject{
			Attributes: attributes,
		},

		Validators: []validator.List{
			listvalidator.SizeAtMost(1),
			listvalidator.ExactlyOneOf(
				exactlyOneOfPaths...,
			),
		},

		PlanModifiers: []planmodifier.List{
			listplanmodifier.RequiresReplace(),
		},
	}
}

func identityProviderClientIdClientSecretAttributes(idpName string) map[string]schema.Attribute {
	const attrMinLength = 1

	return map[string]schema.Attribute{
		"client_id": schema.StringAttribute{
			Description: framework.SchemaAttributeDescriptionFromMarkdown(fmt.Sprintf("A string that specifies the application client ID from %s.", idpName)).Description,
			Required:    true,

			Validators: []validator.String{
				stringvalidator.LengthAtLeast(attrMinLength),
			},
		},

		"client_secret": schema.StringAttribute{
			Description: framework.SchemaAttributeDescriptionFromMarkdown(fmt.Sprintf("A string that specifies the application client secret from %s.", idpName)).Description,
			Required:    true,
			Sensitive:   true,

			Validators: []validator.String{
				stringvalidator.LengthAtLeast(attrMinLength),
			},
		},
	}
}

func (r *IdentityProviderResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *IdentityProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state IdentityProviderResourceModel

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
	identityProvider, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *management.IdentityProvider
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.IdentityProvidersApi.CreateIdentityProvider(ctx, plan.EnvironmentId.ValueString()).IdentityProvider(*identityProvider).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"CreateIdentityProvider",
		framework.DefaultCustomError,
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

func (r *IdentityProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *IdentityProviderResourceModel

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
	var response *management.IdentityProvider
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.IdentityProvidersApi.ReadOneIdentityProvider(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"ReadOneIdentityProvider",
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
	resp.Diagnostics.Append(data.toState(response)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IdentityProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state IdentityProviderResourceModel

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
	identityProvider, d := plan.expand(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run the API call
	var response *management.IdentityProvider
	resp.Diagnostics.Append(framework.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			fO, fR, fErr := r.Client.ManagementAPIClient.IdentityProvidersApi.UpdateIdentityProvider(ctx, plan.EnvironmentId.ValueString(), plan.Id.ValueString()).IdentityProvider(*identityProvider).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, plan.EnvironmentId.ValueString(), fO, fR, fErr)
		},
		"UpdateIdentityProvider",
		framework.DefaultCustomError,
		nil,
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

func (r *IdentityProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *IdentityProviderResourceModel

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
			fR, fErr := r.Client.ManagementAPIClient.IdentityProvidersApi.DeleteIdentityProvider(ctx, data.EnvironmentId.ValueString(), data.Id.ValueString()).Execute()
			return framework.CheckEnvironmentExistsOnPermissionsError(ctx, r.Client.ManagementAPIClient, data.EnvironmentId.ValueString(), nil, fR, fErr)
		},
		"DeleteIdentityProvider",
		framework.CustomErrorResourceNotFoundWarning,
		nil,
		nil,
	)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *IdentityProviderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idComponents := []framework.ImportComponent{
		{
			Label:  "environment_id",
			Regexp: verify.P1ResourceIDRegexp,
		},
		{
			Label:     "identity_provider_id",
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

func (p *IdentityProviderResourceModel) expand(ctx context.Context) (*management.IdentityProvider, diag.Diagnostics) {
	var diags diag.Diagnostics

	common := *management.NewIdentityProviderCommon(p.Enabled.ValueBool(), p.Name.ValueString(), management.ENUMIDENTITYPROVIDEREXT_OPENID_CONNECT)

	if !p.Description.IsNull() && !p.Description.IsUnknown() {
		common.SetDescription(p.Description.ValueString())
	}

	if !p.RegistrationPopulationId.IsNull() && !p.RegistrationPopulationId.IsUnknown() {
		registrationPopulation := *management.NewIdentityProviderCommonRegistrationPopulation()
		registrationPopulation.SetId(p.RegistrationPopulationId.ValueString())
		registration := *management.NewIdentityProviderCommonRegistration()
		registration.SetPopulation(registrationPopulation)
		common.SetRegistration(registration)
	}

	if !p.LoginButtonIcon.IsNull() && !p.LoginButtonIcon.IsUnknown() {
		var plan IdentityProviderLoginButtonIcon
		d := p.LoginButtonIcon.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		icon := *management.NewIdentityProviderCommonLoginButtonIcon()
		icon.SetId(plan.Id.ValueString())
		icon.SetHref(plan.Href.ValueString())
		common.SetLoginButtonIcon(icon)

	}

	if !p.Icon.IsNull() && !p.Icon.IsUnknown() {
		var plan IdentityProviderIcon
		d := p.Icon.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		icon := *management.NewIdentityProviderCommonIcon()
		icon.SetId(plan.Id.ValueString())
		icon.SetHref(plan.Href.ValueString())
		common.SetIcon(icon)
	}

	data := &management.IdentityProvider{}

	processedCount := 0

	if !p.Facebook.IsNull() && !p.Facebook.IsUnknown() {
		var plan []IdentityProviderFacebookResourceModel
		d := p.Facebook.ElementsAs(ctx, &plan, false)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		if len(plan) == 0 {
			diags.AddError(
				"Invalid configuration",
				"The `facebook` block is declared but has no configuration.  Please report this to the provider maintainers.",
			)
		}

		planItem := plan[0]

		idpData := management.IdentityProviderFacebook{
			Enabled:         common.Enabled,
			Name:            common.Name,
			Type:            management.ENUMIDENTITYPROVIDEREXT_FACEBOOK,
			Description:     common.Description,
			Registration:    common.Registration,
			LoginButtonIcon: common.LoginButtonIcon,
			Icon:            common.Icon,
		}

		idpData.SetAppId(planItem.AppId.ValueString())
		idpData.SetAppSecret(planItem.AppSecret.ValueString())

		data.IdentityProviderFacebook = &idpData
		processedCount += 1
	}

	if !p.Google.IsNull() && !p.Google.IsUnknown() {
		var plan []IdentityProviderGoogleResourceModel
		d := p.Google.ElementsAs(ctx, &plan, false)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		if len(plan) == 0 {
			diags.AddError(
				"Invalid configuration",
				"The `google` block is declared but has no configuration.  Please report this to the provider maintainers.",
			)
		}

		planItem := plan[0]

		idpData := management.IdentityProviderClientIDClientSecret{
			Enabled:         common.Enabled,
			Name:            common.Name,
			Type:            management.ENUMIDENTITYPROVIDEREXT_GOOGLE,
			Description:     common.Description,
			Registration:    common.Registration,
			LoginButtonIcon: common.LoginButtonIcon,
			Icon:            common.Icon,
		}

		idpData.SetClientId(planItem.ClientId.ValueString())
		idpData.SetClientSecret(planItem.ClientSecret.ValueString())

		data.IdentityProviderClientIDClientSecret = &idpData
		processedCount += 1
	}

	if !p.LinkedIn.IsNull() && !p.LinkedIn.IsUnknown() {
		var plan []IdentityProviderLinkedInResourceModel
		d := p.LinkedIn.ElementsAs(ctx, &plan, false)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		if len(plan) == 0 {
			diags.AddError(
				"Invalid configuration",
				"The `linkedin` block is declared but has no configuration.  Please report this to the provider maintainers.",
			)
		}

		planItem := plan[0]

		idpData := management.IdentityProviderClientIDClientSecret{
			Enabled:         common.Enabled,
			Name:            common.Name,
			Type:            management.ENUMIDENTITYPROVIDEREXT_LINKEDIN,
			Description:     common.Description,
			Registration:    common.Registration,
			LoginButtonIcon: common.LoginButtonIcon,
			Icon:            common.Icon,
		}

		idpData.SetClientId(planItem.ClientId.ValueString())
		idpData.SetClientSecret(planItem.ClientSecret.ValueString())

		data.IdentityProviderClientIDClientSecret = &idpData
		processedCount += 1
	}

	if !p.Yahoo.IsNull() && !p.Yahoo.IsUnknown() {
		var plan []IdentityProviderYahooResourceModel
		d := p.Yahoo.ElementsAs(ctx, &plan, false)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		if len(plan) == 0 {
			diags.AddError(
				"Invalid configuration",
				"The `yahoo` block is declared but has no configuration.  Please report this to the provider maintainers.",
			)
		}

		planItem := plan[0]

		idpData := management.IdentityProviderClientIDClientSecret{
			Enabled:         common.Enabled,
			Name:            common.Name,
			Type:            management.ENUMIDENTITYPROVIDEREXT_YAHOO,
			Description:     common.Description,
			Registration:    common.Registration,
			LoginButtonIcon: common.LoginButtonIcon,
			Icon:            common.Icon,
		}

		idpData.SetClientId(planItem.ClientId.ValueString())
		idpData.SetClientSecret(planItem.ClientSecret.ValueString())

		data.IdentityProviderClientIDClientSecret = &idpData
		processedCount += 1
	}

	if !p.Amazon.IsNull() && !p.Amazon.IsUnknown() {
		var plan []IdentityProviderAmazonResourceModel
		d := p.Amazon.ElementsAs(ctx, &plan, false)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		if len(plan) == 0 {
			diags.AddError(
				"Invalid configuration",
				"The `amazon` block is declared but has no configuration.  Please report this to the provider maintainers.",
			)
		}

		planItem := plan[0]

		idpData := management.IdentityProviderClientIDClientSecret{
			Enabled:         common.Enabled,
			Name:            common.Name,
			Type:            management.ENUMIDENTITYPROVIDEREXT_AMAZON,
			Description:     common.Description,
			Registration:    common.Registration,
			LoginButtonIcon: common.LoginButtonIcon,
			Icon:            common.Icon,
		}

		idpData.SetClientId(planItem.ClientId.ValueString())
		idpData.SetClientSecret(planItem.ClientSecret.ValueString())

		data.IdentityProviderClientIDClientSecret = &idpData
		processedCount += 1
	}

	if !p.Twitter.IsNull() && !p.Twitter.IsUnknown() {
		var plan []IdentityProviderTwitterResourceModel
		d := p.Twitter.ElementsAs(ctx, &plan, false)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		if len(plan) == 0 {
			diags.AddError(
				"Invalid configuration",
				"The `twitter` block is declared but has no configuration.  Please report this to the provider maintainers.",
			)
		}

		planItem := plan[0]

		idpData := management.IdentityProviderClientIDClientSecret{
			Enabled:         common.Enabled,
			Name:            common.Name,
			Type:            management.ENUMIDENTITYPROVIDEREXT_TWITTER,
			Description:     common.Description,
			Registration:    common.Registration,
			LoginButtonIcon: common.LoginButtonIcon,
			Icon:            common.Icon,
		}

		idpData.SetClientId(planItem.ClientId.ValueString())
		idpData.SetClientSecret(planItem.ClientSecret.ValueString())

		data.IdentityProviderClientIDClientSecret = &idpData
		processedCount += 1
	}

	if !p.Apple.IsNull() && !p.Apple.IsUnknown() {
		var plan []IdentityProviderAppleResourceModel
		d := p.Apple.ElementsAs(ctx, &plan, false)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		if len(plan) == 0 {
			diags.AddError(
				"Invalid configuration",
				"The `apple` block is declared but has no configuration.  Please report this to the provider maintainers.",
			)
		}

		planItem := plan[0]

		idpData := management.IdentityProviderApple{
			Enabled:         common.Enabled,
			Name:            common.Name,
			Type:            management.ENUMIDENTITYPROVIDEREXT_APPLE,
			Description:     common.Description,
			Registration:    common.Registration,
			LoginButtonIcon: common.LoginButtonIcon,
			Icon:            common.Icon,
		}

		idpData.SetClientId(planItem.ClientId.ValueString())
		idpData.SetClientSecretSigningKey(planItem.ClientSecretSigningKey.ValueString())
		idpData.SetKeyId(planItem.KeyId.ValueString())
		idpData.SetTeamId(planItem.TeamId.ValueString())

		data.IdentityProviderApple = &idpData
		processedCount += 1
	}

	if !p.Paypal.IsNull() && !p.Paypal.IsUnknown() {
		var plan []IdentityProviderPaypalResourceModel
		d := p.Paypal.ElementsAs(ctx, &plan, false)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		if len(plan) == 0 {
			diags.AddError(
				"Invalid configuration",
				"The `paypal` block is declared but has no configuration.  Please report this to the provider maintainers.",
			)
		}

		planItem := plan[0]

		idpData := management.IdentityProviderPaypal{
			Enabled:         common.Enabled,
			Name:            common.Name,
			Type:            management.ENUMIDENTITYPROVIDEREXT_PAYPAL,
			Description:     common.Description,
			Registration:    common.Registration,
			LoginButtonIcon: common.LoginButtonIcon,
			Icon:            common.Icon,
		}

		idpData.SetClientId(planItem.ClientId.ValueString())
		idpData.SetClientSecret(planItem.ClientSecret.ValueString())
		idpData.SetClientEnvironment(planItem.ClientEnvironment.ValueString())

		data.IdentityProviderPaypal = &idpData
		processedCount += 1
	}

	if !p.Microsoft.IsNull() && !p.Microsoft.IsUnknown() {
		var plan []IdentityProviderMicrosoftResourceModel
		d := p.Microsoft.ElementsAs(ctx, &plan, false)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		if len(plan) == 0 {
			diags.AddError(
				"Invalid configuration",
				"The `microsoft` block is declared but has no configuration.  Please report this to the provider maintainers.",
			)
		}

		planItem := plan[0]

		idpData := management.IdentityProviderClientIDClientSecret{
			Enabled:         common.Enabled,
			Name:            common.Name,
			Type:            management.ENUMIDENTITYPROVIDEREXT_MICROSOFT,
			Description:     common.Description,
			Registration:    common.Registration,
			LoginButtonIcon: common.LoginButtonIcon,
			Icon:            common.Icon,
		}

		idpData.SetClientId(planItem.ClientId.ValueString())
		idpData.SetClientSecret(planItem.ClientSecret.ValueString())

		data.IdentityProviderClientIDClientSecret = &idpData
		processedCount += 1
	}

	if !p.Github.IsNull() && !p.Github.IsUnknown() {
		var plan []IdentityProviderGithubResourceModel
		d := p.Github.ElementsAs(ctx, &plan, false)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		if len(plan) == 0 {
			diags.AddError(
				"Invalid configuration",
				"The `github` block is declared but has no configuration.  Please report this to the provider maintainers.",
			)
		}

		planItem := plan[0]

		idpData := management.IdentityProviderClientIDClientSecret{
			Enabled:         common.Enabled,
			Name:            common.Name,
			Type:            management.ENUMIDENTITYPROVIDEREXT_GITHUB,
			Description:     common.Description,
			Registration:    common.Registration,
			LoginButtonIcon: common.LoginButtonIcon,
			Icon:            common.Icon,
		}

		idpData.SetClientId(planItem.ClientId.ValueString())
		idpData.SetClientSecret(planItem.ClientSecret.ValueString())

		data.IdentityProviderClientIDClientSecret = &idpData
		processedCount += 1
	}

	if !p.OpenIDConnect.IsNull() && !p.OpenIDConnect.IsUnknown() {
		var plan []IdentityProviderOIDCResourceModel
		d := p.OpenIDConnect.ElementsAs(ctx, &plan, false)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		if len(plan) == 0 {
			diags.AddError(
				"Invalid configuration",
				"The `openid_connect` block is declared but has no configuration.  Please report this to the provider maintainers.",
			)
		}

		planItem := plan[0]

		idpData := management.IdentityProviderOIDC{
			Enabled:         common.Enabled,
			Name:            common.Name,
			Type:            management.ENUMIDENTITYPROVIDEREXT_OPENID_CONNECT,
			Description:     common.Description,
			Registration:    common.Registration,
			LoginButtonIcon: common.LoginButtonIcon,
			Icon:            common.Icon,
		}

		if !planItem.AuthorizationEndpoint.IsNull() && !planItem.AuthorizationEndpoint.IsUnknown() {
			idpData.SetAuthorizationEndpoint(planItem.AuthorizationEndpoint.ValueString())
		}

		if !planItem.ClientId.IsNull() && !planItem.ClientId.IsUnknown() {
			idpData.SetClientId(planItem.ClientId.ValueString())
		}

		if !planItem.ClientSecret.IsNull() && !planItem.ClientSecret.IsUnknown() {
			idpData.SetClientSecret(planItem.ClientSecret.ValueString())
		}

		if !planItem.DiscoveryEndpoint.IsNull() && !planItem.DiscoveryEndpoint.IsUnknown() {
			idpData.SetDiscoveryEndpoint(planItem.DiscoveryEndpoint.ValueString())
		}

		if !planItem.Issuer.IsNull() && !planItem.Issuer.IsUnknown() {
			idpData.SetIssuer(planItem.Issuer.ValueString())
		}

		if !planItem.JwksEndpoint.IsNull() && !planItem.JwksEndpoint.IsUnknown() {
			idpData.SetJwksEndpoint(planItem.JwksEndpoint.ValueString())
		}

		if !planItem.Scopes.IsNull() && !planItem.Scopes.IsUnknown() {
			var scopesPlan []string
			diags.Append(planItem.Scopes.ElementsAs(ctx, &scopesPlan, false)...)
			if diags.HasError() {
				return nil, diags
			}

			idpData.SetScopes(scopesPlan)
		}

		if !planItem.TokenEndpoint.IsNull() && !planItem.TokenEndpoint.IsUnknown() {
			idpData.SetTokenEndpoint(planItem.TokenEndpoint.ValueString())
		}

		if !planItem.TokenEndpointAuthMethod.IsNull() && !planItem.TokenEndpointAuthMethod.IsUnknown() {
			idpData.SetTokenEndpointAuthMethod(management.EnumIdentityProviderOIDCTokenAuthMethod(planItem.TokenEndpointAuthMethod.ValueString()))
		}

		if !planItem.UserinfoEndpoint.IsNull() && !planItem.UserinfoEndpoint.IsUnknown() {
			idpData.SetUserInfoEndpoint(planItem.UserinfoEndpoint.ValueString())
		}

		data.IdentityProviderOIDC = &idpData
		processedCount += 1
	}

	if !p.Saml.IsNull() && !p.Saml.IsUnknown() {
		var plan []IdentityProviderSAMLResourceModel
		d := p.Saml.ElementsAs(ctx, &plan, false)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		if len(plan) == 0 {
			diags.AddError(
				"Invalid configuration",
				"The `saml` block is declared but has no configuration.  Please report this to the provider maintainers.",
			)
		}

		planItem := plan[0]

		idpData := management.IdentityProviderSAML{
			Enabled:         common.Enabled,
			Name:            common.Name,
			Type:            management.ENUMIDENTITYPROVIDEREXT_SAML,
			Description:     common.Description,
			Registration:    common.Registration,
			LoginButtonIcon: common.LoginButtonIcon,
			Icon:            common.Icon,
		}

		if !planItem.AuthenticationRequestSigned.IsNull() && !planItem.AuthenticationRequestSigned.IsUnknown() {
			idpData.SetAuthnRequestSigned(planItem.AuthenticationRequestSigned.ValueBool())
		}

		if !planItem.IdpEntityId.IsNull() && !planItem.IdpEntityId.IsUnknown() {
			idpData.SetIdpEntityId(planItem.IdpEntityId.ValueString())
		}

		if !planItem.SpEntityId.IsNull() && !planItem.SpEntityId.IsUnknown() {
			idpData.SetSpEntityId(planItem.SpEntityId.ValueString())
		}

		if !planItem.IdpVerificationCertificateIds.IsNull() && !planItem.IdpVerificationCertificateIds.IsUnknown() {
			var certificateIdsPlan []string
			diags.Append(planItem.IdpVerificationCertificateIds.ElementsAs(ctx, &certificateIdsPlan, false)...)
			if diags.HasError() {
				return nil, diags
			}

			idpVerificationCertificates := make([]management.IdentityProviderSAMLAllOfIdpVerificationCertificates, 0)
			for _, v := range certificateIdsPlan {
				idpVerificationCertificates = append(idpVerificationCertificates, *management.NewIdentityProviderSAMLAllOfIdpVerificationCertificates(v))
			}

			idpVerifcation := management.NewIdentityProviderSAMLAllOfIdpVerification(idpVerificationCertificates)
			idpData.SetIdpVerification(*idpVerifcation)
		}

		if !planItem.SpSigningKeyId.IsNull() && !planItem.SpSigningKeyId.IsUnknown() {
			spSigningKey := management.NewIdentityProviderSAMLAllOfSpSigningKey(planItem.SpSigningKeyId.ValueString())
			spSigning := management.NewIdentityProviderSAMLAllOfSpSigning(*spSigningKey)
			idpData.SetSpSigning(*spSigning)
		}

		if !planItem.SsoBinding.IsNull() && !planItem.SsoBinding.IsUnknown() {
			idpData.SetSsoBinding(management.EnumIdentityProviderSAMLSSOBinding(planItem.SsoBinding.ValueString()))
		}

		if !planItem.SsoEndpoint.IsNull() && !planItem.SsoEndpoint.IsUnknown() {
			idpData.SetSsoEndpoint(planItem.SsoEndpoint.ValueString())
		}

		if !planItem.SloBinding.IsNull() && !planItem.SloBinding.IsUnknown() {
			idpData.SetSloBinding(management.EnumIdentityProviderSAMLSLOBinding(planItem.SloBinding.ValueString()))
		}

		if !planItem.SloEndpoint.IsNull() && !planItem.SloEndpoint.IsUnknown() {
			idpData.SetSloEndpoint(planItem.SloEndpoint.ValueString())
		}

		if !planItem.SloResponseEndpoint.IsNull() && !planItem.SloResponseEndpoint.IsUnknown() {
			idpData.SetSloResponseEndpoint(planItem.SloResponseEndpoint.ValueString())
		}

		if !planItem.SloWindow.IsNull() && !planItem.SloWindow.IsUnknown() {
			idpData.SetSloWindow(int32(planItem.SloWindow.ValueInt64()))
		}

		data.IdentityProviderSAML = &idpData
		processedCount += 1
	}

	if processedCount > 1 {
		diags.AddError(
			"Invalid configuration",
			"More than one identity provider type configured.  This is not supported.  Please raise an issue with the provider maintainers.",
		)

		return nil, diags
	} else if processedCount == 0 {
		diags.AddError(
			"Invalid configuration",
			"No identity provider types configured.  This is not supported.  Please raise an issue with the provider maintainers.",
		)

		return nil, diags
	}

	return data, diags
}

func (p *IdentityProviderResourceModel) toState(apiObject *management.IdentityProvider) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiObject == nil {
		diags.AddError(
			"Data object missing",
			"Cannot convert the data object to state as the data object is nil.  Please report this to the provider maintainers.",
		)

		return diags
	}

	byteData, err := apiObject.MarshalJSON()
	if err != nil {
		diags.AddError(
			"Data object invalid",
			"Cannot convert the data object to state as the data object is not a valid type.  Please report this to the provider maintainers.",
		)

		return diags
	}

	var common management.IdentityProviderCommon
	if err := json.Unmarshal(byteData, &common); err != nil {
		diags.AddError(
			"Data object invalid",
			"Cannot convert the data object to state as the data object cannot be converted.  Please report this to the provider maintainers.",
		)

		return diags
	}

	p.Id = framework.StringOkToTF(common.GetIdOk())
	p.EnvironmentId = framework.StringOkToTF(common.Environment.GetIdOk())
	p.Name = framework.StringOkToTF(common.GetNameOk())
	p.Description = framework.StringOkToTF(common.GetDescriptionOk())
	p.Enabled = framework.BoolOkToTF(common.GetEnabledOk())

	p.RegistrationPopulationId = types.StringNull()
	if v, ok := common.GetRegistrationOk(); ok {
		if q, ok := v.GetPopulationOk(); ok {
			p.RegistrationPopulationId = framework.StringOkToTF(q.GetIdOk())
		}
	}

	var d diag.Diagnostics
	p.LoginButtonIcon, d = service.ImageOkToTF(common.GetLoginButtonIconOk())
	diags.Append(d...)

	p.Icon, d = service.ImageOkToTF(common.GetIconOk())
	diags.Append(d...)

	// The providers
	p.Facebook, d = identityProviderFacebookToTF(apiObject.IdentityProviderFacebook)
	diags.Append(d...)

	p.Google, d = identityProviderClientIDClientSecretToTF(apiObject.IdentityProviderClientIDClientSecret, management.ENUMIDENTITYPROVIDEREXT_GOOGLE)
	diags.Append(d...)

	p.LinkedIn, d = identityProviderClientIDClientSecretToTF(apiObject.IdentityProviderClientIDClientSecret, management.ENUMIDENTITYPROVIDEREXT_LINKEDIN)
	diags.Append(d...)

	p.Yahoo, d = identityProviderClientIDClientSecretToTF(apiObject.IdentityProviderClientIDClientSecret, management.ENUMIDENTITYPROVIDEREXT_YAHOO)
	diags.Append(d...)

	p.Amazon, d = identityProviderClientIDClientSecretToTF(apiObject.IdentityProviderClientIDClientSecret, management.ENUMIDENTITYPROVIDEREXT_AMAZON)
	diags.Append(d...)

	p.Twitter, d = identityProviderClientIDClientSecretToTF(apiObject.IdentityProviderClientIDClientSecret, management.ENUMIDENTITYPROVIDEREXT_TWITTER)
	diags.Append(d...)

	p.Apple, d = identityProviderAppleToTF(apiObject.IdentityProviderApple)
	diags.Append(d...)

	p.Paypal, d = identityProviderPaypalToTF(apiObject.IdentityProviderPaypal)
	diags.Append(d...)

	p.Microsoft, d = identityProviderClientIDClientSecretToTF(apiObject.IdentityProviderClientIDClientSecret, management.ENUMIDENTITYPROVIDEREXT_MICROSOFT)
	diags.Append(d...)

	p.Github, d = identityProviderClientIDClientSecretToTF(apiObject.IdentityProviderClientIDClientSecret, management.ENUMIDENTITYPROVIDEREXT_GITHUB)
	diags.Append(d...)

	p.OpenIDConnect, d = identityProviderOIDCToTF(apiObject.IdentityProviderOIDC)
	diags.Append(d...)

	p.Saml, d = identityProviderSAMLToTF(apiObject.IdentityProviderSAML)
	diags.Append(d...)

	return diags
}

func identityProviderFacebookToTF(idpApiObject *management.IdentityProviderFacebook) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	tfObjType := types.ObjectType{AttrTypes: identityProviderFacebookTFObjectTypes}

	if idpApiObject == nil || idpApiObject.GetType() != management.ENUMIDENTITYPROVIDEREXT_FACEBOOK {
		return types.ListNull(tfObjType), diags
	}

	attributesMap := map[string]attr.Value{
		"app_id":     framework.StringOkToTF(idpApiObject.GetAppIdOk()),
		"app_secret": framework.StringOkToTF(idpApiObject.GetAppSecretOk()),
	}

	flattenedObj, d := types.ObjectValue(identityProviderFacebookTFObjectTypes, attributesMap)
	diags.Append(d...)

	returnVar, d := types.ListValue(tfObjType, append([]attr.Value{}, flattenedObj))
	diags.Append(d...)

	return returnVar, diags
}

func identityProviderClientIDClientSecretToTF(idpApiObject *management.IdentityProviderClientIDClientSecret, idpType management.EnumIdentityProviderExt) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	tfObjType := types.ObjectType{AttrTypes: identityProviderClientIDClientSecretTFObjectTypes}

	if idpApiObject == nil || idpApiObject.GetType() != idpType {
		return types.ListNull(tfObjType), diags
	}

	attributesMap := map[string]attr.Value{
		"client_id":     framework.StringOkToTF(idpApiObject.GetClientIdOk()),
		"client_secret": framework.StringOkToTF(idpApiObject.GetClientSecretOk()),
	}

	flattenedObj, d := types.ObjectValue(identityProviderClientIDClientSecretTFObjectTypes, attributesMap)
	diags.Append(d...)

	returnVar, d := types.ListValue(tfObjType, append([]attr.Value{}, flattenedObj))
	diags.Append(d...)

	return returnVar, diags
}

func identityProviderAppleToTF(idpApiObject *management.IdentityProviderApple) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	tfObjType := types.ObjectType{AttrTypes: identityProviderAppleTFObjectTypes}

	if idpApiObject == nil || idpApiObject.GetType() != management.ENUMIDENTITYPROVIDEREXT_APPLE {
		return types.ListNull(tfObjType), diags
	}

	attributesMap := map[string]attr.Value{
		"team_id":                   framework.StringOkToTF(idpApiObject.GetTeamIdOk()),
		"key_id":                    framework.StringOkToTF(idpApiObject.GetKeyIdOk()),
		"client_id":                 framework.StringOkToTF(idpApiObject.GetClientIdOk()),
		"client_secret_signing_key": framework.StringOkToTF(idpApiObject.GetClientSecretSigningKeyOk()),
	}

	flattenedObj, d := types.ObjectValue(identityProviderAppleTFObjectTypes, attributesMap)
	diags.Append(d...)

	returnVar, d := types.ListValue(tfObjType, append([]attr.Value{}, flattenedObj))
	diags.Append(d...)

	return returnVar, diags
}

func identityProviderPaypalToTF(idpApiObject *management.IdentityProviderPaypal) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	tfObjType := types.ObjectType{AttrTypes: identityProviderPaypalTFObjectTypes}

	if idpApiObject == nil || idpApiObject.GetType() != management.ENUMIDENTITYPROVIDEREXT_PAYPAL {
		return types.ListNull(tfObjType), diags
	}

	attributesMap := map[string]attr.Value{
		"client_id":          framework.StringOkToTF(idpApiObject.GetClientIdOk()),
		"client_secret":      framework.StringOkToTF(idpApiObject.GetClientSecretOk()),
		"client_environment": framework.StringOkToTF(idpApiObject.GetClientEnvironmentOk()),
	}

	flattenedObj, d := types.ObjectValue(identityProviderPaypalTFObjectTypes, attributesMap)
	diags.Append(d...)

	returnVar, d := types.ListValue(tfObjType, append([]attr.Value{}, flattenedObj))
	diags.Append(d...)

	return returnVar, diags
}

func identityProviderOIDCToTF(idpApiObject *management.IdentityProviderOIDC) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	tfObjType := types.ObjectType{AttrTypes: identityProviderOIDCTFObjectTypes}

	if idpApiObject == nil || idpApiObject.GetType() != management.ENUMIDENTITYPROVIDEREXT_OPENID_CONNECT {
		return types.ListNull(tfObjType), diags
	}

	attributesMap := map[string]attr.Value{
		"authorization_endpoint":     framework.StringOkToTF(idpApiObject.GetAuthorizationEndpointOk()),
		"client_id":                  framework.StringOkToTF(idpApiObject.GetClientIdOk()),
		"client_secret":              framework.StringOkToTF(idpApiObject.GetClientSecretOk()),
		"discovery_endpoint":         framework.StringOkToTF(idpApiObject.GetDiscoveryEndpointOk()),
		"issuer":                     framework.StringOkToTF(idpApiObject.GetIssuerOk()),
		"jwks_endpoint":              framework.StringOkToTF(idpApiObject.GetJwksEndpointOk()),
		"scopes":                     framework.StringSetOkToTF(idpApiObject.GetScopesOk()),
		"token_endpoint":             framework.StringOkToTF(idpApiObject.GetTokenEndpointOk()),
		"token_endpoint_auth_method": framework.EnumOkToTF(idpApiObject.GetTokenEndpointAuthMethodOk()),
		"userinfo_endpoint":          framework.StringOkToTF(idpApiObject.GetUserInfoEndpointOk()),
	}

	flattenedObj, d := types.ObjectValue(identityProviderOIDCTFObjectTypes, attributesMap)
	diags.Append(d...)

	returnVar, d := types.ListValue(tfObjType, append([]attr.Value{}, flattenedObj))
	diags.Append(d...)

	return returnVar, diags
}

func identityProviderSAMLToTF(idpApiObject *management.IdentityProviderSAML) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	tfObjType := types.ObjectType{AttrTypes: identityProviderSAMLTFObjectTypes}

	if idpApiObject == nil || idpApiObject.GetType() != management.ENUMIDENTITYPROVIDEREXT_SAML {
		return types.ListNull(tfObjType), diags
	}

	attributesMap := map[string]attr.Value{
		"authentication_request_signed": framework.BoolOkToTF(idpApiObject.GetAuthnRequestSignedOk()),
		"idp_entity_id":                 framework.StringOkToTF(idpApiObject.GetIdpEntityIdOk()),
		"sp_entity_id":                  framework.StringOkToTF(idpApiObject.GetSpEntityIdOk()),
		"sso_binding":                   framework.EnumOkToTF(idpApiObject.GetSsoBindingOk()),
		"sso_endpoint":                  framework.StringOkToTF(idpApiObject.GetSsoEndpointOk()),
		"slo_binding":                   framework.EnumOkToTF(idpApiObject.GetSloBindingOk()),
		"slo_endpoint":                  framework.StringOkToTF(idpApiObject.GetSloEndpointOk()),
		"slo_response_endpoint":         framework.StringOkToTF(idpApiObject.GetSloResponseEndpointOk()),
		"slo_window":                    framework.Int32OkToTF(idpApiObject.GetSloWindowOk()),
	}

	attributesMap["idp_verification_certificate_ids"] = types.SetNull(types.StringType)
	if v, ok := idpApiObject.GetIdpVerificationOk(); ok {
		if c, ok := v.GetCertificatesOk(); ok {
			ids := make([]string, 0)
			for _, certificate := range c {
				ids = append(ids, certificate.GetId())
			}

			attributesMap["idp_verification_certificate_ids"] = framework.StringSetToTF(ids)
		}
	}

	attributesMap["sp_signing_key_id"] = types.StringNull()
	if v, ok := idpApiObject.GetSpSigningOk(); ok {
		if c, ok := v.GetKeyOk(); ok {
			attributesMap["sp_signing_key_id"] = framework.StringOkToTF(c.GetIdOk())
		}
	}

	flattenedObj, d := types.ObjectValue(identityProviderSAMLTFObjectTypes, attributesMap)
	diags.Append(d...)

	returnVar, d := types.ListValue(tfObjType, append([]attr.Value{}, flattenedObj))
	diags.Append(d...)

	return returnVar, diags
}
