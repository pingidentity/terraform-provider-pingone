package sso

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
	"github.com/pingidentity/terraform-provider-pingone/internal/sdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/service"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Types
type IdentityProviderResource serviceClientType

type IdentityProviderResourceModel struct {
	Id                       pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId            pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	Name                     types.String                 `tfsdk:"name"`
	Description              types.String                 `tfsdk:"description"`
	Enabled                  types.Bool                   `tfsdk:"enabled"`
	RegistrationPopulationId pingonetypes.ResourceIDValue `tfsdk:"registration_population_id"`
	LoginButtonIcon          types.Object                 `tfsdk:"login_button_icon"`
	Icon                     types.Object                 `tfsdk:"icon"`
	Facebook                 types.Object                 `tfsdk:"facebook"`
	Google                   types.Object                 `tfsdk:"google"`
	LinkedIn                 types.Object                 `tfsdk:"linkedin"`
	Yahoo                    types.Object                 `tfsdk:"yahoo"`
	Amazon                   types.Object                 `tfsdk:"amazon"`
	Twitter                  types.Object                 `tfsdk:"twitter"`
	Apple                    types.Object                 `tfsdk:"apple"`
	Paypal                   types.Object                 `tfsdk:"paypal"`
	Microsoft                types.Object                 `tfsdk:"microsoft"`
	Github                   types.Object                 `tfsdk:"github"`
	OpenIDConnect            types.Object                 `tfsdk:"openid_connect"`
	Saml                     types.Object                 `tfsdk:"saml"`
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
	PkceMethod              types.String `tfsdk:"pkce_method"`
	JwksEndpoint            types.String `tfsdk:"jwks_endpoint"`
	Scopes                  types.Set    `tfsdk:"scopes"`
	TokenEndpoint           types.String `tfsdk:"token_endpoint"`
	TokenEndpointAuthMethod types.String `tfsdk:"token_endpoint_auth_method"`
	UserinfoEndpoint        types.String `tfsdk:"userinfo_endpoint"`
}

type IdentityProviderSAMLResourceModel struct {
	AuthenticationRequestSigned types.Bool   `tfsdk:"authentication_request_signed"`
	IdpEntityId                 types.String `tfsdk:"idp_entity_id"`
	SpEntityId                  types.String `tfsdk:"sp_entity_id"`
	IdpVerification             types.Object `tfsdk:"idp_verification"`
	SpSigning                   types.Object `tfsdk:"sp_signing"`
	SsoBinding                  types.String `tfsdk:"sso_binding"`
	SsoEndpoint                 types.String `tfsdk:"sso_endpoint"`
	SloBinding                  types.String `tfsdk:"slo_binding"`
	SloEndpoint                 types.String `tfsdk:"slo_endpoint"`
	SloResponseEndpoint         types.String `tfsdk:"slo_response_endpoint"`
	SloWindow                   types.Int64  `tfsdk:"slo_window"`
}

type IdentityProviderSAMLResourceIdPVerificationModel struct {
	Certificates types.Set `tfsdk:"certificates"`
}

type IdentityProviderSAMLResourceIdPVerificationCertificatesModel struct {
	Id pingonetypes.ResourceIDValue `tfsdk:"id"`
}

type IdentityProviderSAMLResourceSpSigningModel struct {
	Key       types.Object `tfsdk:"key"`
	Algorithm types.String `tfsdk:"algorithm"`
}

type IdentityProviderSAMLResourceSpSigningKeyModel struct {
	Id pingonetypes.ResourceIDValue `tfsdk:"id"`
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
		"pkce_method":                types.StringType,
		"jwks_endpoint":              types.StringType,
		"scopes":                     types.SetType{ElemType: types.StringType},
		"token_endpoint":             types.StringType,
		"token_endpoint_auth_method": types.StringType,
		"userinfo_endpoint":          types.StringType,
	}

	identityProviderSAMLTFObjectTypes = map[string]attr.Type{
		"authentication_request_signed": types.BoolType,
		"idp_entity_id":                 types.StringType,
		"sp_entity_id":                  types.StringType,
		"idp_verification":              types.ObjectType{AttrTypes: identityProviderSAMLIdPVerificationTFObjectTypes},
		"sp_signing":                    types.ObjectType{AttrTypes: identityProviderSAMLSpSigningTFObjectTypes},
		"sso_binding":                   types.StringType,
		"sso_endpoint":                  types.StringType,
		"slo_binding":                   types.StringType,
		"slo_endpoint":                  types.StringType,
		"slo_response_endpoint":         types.StringType,
		"slo_window":                    types.Int64Type,
	}

	identityProviderSAMLIdPVerificationTFObjectTypes = map[string]attr.Type{
		"certificates": types.SetType{ElemType: types.ObjectType{AttrTypes: identityProviderSAMLIdPVerificationCertificateTFObjectTypes}},
	}

	identityProviderSAMLIdPVerificationCertificateTFObjectTypes = map[string]attr.Type{
		"id": pingonetypes.ResourceIDType{},
	}

	identityProviderSAMLSpSigningTFObjectTypes = map[string]attr.Type{
		"key":       types.ObjectType{AttrTypes: identityProviderSAMLSpSigningKeyTFObjectTypes},
		"algorithm": types.StringType,
	}

	identityProviderSAMLSpSigningKeyTFObjectTypes = map[string]attr.Type{
		"id": pingonetypes.ResourceIDType{},
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
		"The URL or fully qualified path to the identity provider icon to use as the login button.  This can be retrieved from the `uploaded_image.href` parameter of the `pingone_image` resource.",
	)

	iconIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The ID for the identity provider icon to use as the login button.  This can be retrieved from the `id` parameter of the `pingone_image` resource.  Must be a valid PingOne resource ID.",
	)

	iconHrefDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The URL or fully qualified path to the identity provider icon to use as the login button.  This can be retrieved from the `uploaded_image.href` parameter of the `pingone_image` resource.",
	)

	paypalClientEnvironmentDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the PayPal environment.",
	).AllowedValues("sandbox", "live")

	oidcPkceMethodDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the method for PKCE. This value auto-populates from a discovery endpoint if the OpenID Provider includes `S256` in its `code_challenge_methods_supported` claim. The plain method is not currently supported.",
	).AllowedValuesEnum(management.AllowedEnumIdentityProviderPKCEMethodEnumValues).DefaultValue(string(management.ENUMIDENTITYPROVIDERPKCEMETHOD_NONE))

	oidcTokenEndpointAuthMethodDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the OIDC identity provider's token endpoint authentication method.",
	).AllowedValuesEnum(management.AllowedEnumIdentityProviderOIDCTokenAuthMethodEnumValues).DefaultValue(string(management.ENUMIDENTITYPROVIDEROIDCTOKENAUTHMETHOD_CLIENT_SECRET_BASIC))

	samlAuthenticationRequestSignedDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A boolean that specifies whether the SAML authentication request will be signed when sending to the identity provider. Set this to `true` if the external IDP is included in an authentication policy to be used by applications that are accessed using a mix of default URLS and custom Domains URLs.",
	).DefaultValue(false)

	samlIdpEntityIdDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A string that specifies the entity ID URI that is checked against the `issuerId` tag in the incoming response.",
	)

	samlSpSigningDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"A single object that specifies settings for SAML assertion signing, including the key and the signature algorithm.  Required when `authentication_request_signed` is set to `true`.",
	)

	samlSpSigningAlgorithmDescription := framework.SchemaAttributeDescriptionFromMarkdown(
		"The signing key algorithm used by PingOne. The value will depend on which key algorithm and signature algorithm you chose when creating your signing key.",
	).AllowedValuesEnum(management.AllowedEnumIdentityProviderSAMLSigningAlgorithmEnumValues)

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

				CustomType: pingonetypes.ResourceIDType{},
			},

			"login_button_icon": schema.SingleNestedAttribute{
				Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that specifies the HREF and ID for the identity provider icon to use in the login button.").Description,
				Optional:    true,

				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description:         loginButtonIconIdDescription.Description,
						MarkdownDescription: loginButtonIconIdDescription.MarkdownDescription,
						Required:            true,

						CustomType: pingonetypes.ResourceIDType{},
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

						CustomType: pingonetypes.ResourceIDType{},
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

			// The providers
			"facebook": identityProviderSchemaAttribute(
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

			"google": identityProviderSchemaAttribute(
				framework.SchemaAttributeDescriptionFromMarkdown("A single block that specifies options for connectivity to the Google social identity provider."),

				identityProviderClientIdClientSecretAttributes("Google"),

				providerAttributeList,
			),

			"linkedin": identityProviderSchemaAttribute(
				framework.SchemaAttributeDescriptionFromMarkdown("A single block that specifies options for connectivity to the LinkedIn social identity provider."),

				identityProviderClientIdClientSecretAttributes("LinkedIn"),

				providerAttributeList,
			),

			"yahoo": identityProviderSchemaAttribute(
				framework.SchemaAttributeDescriptionFromMarkdown("A single block that specifies options for connectivity to the Yahoo social identity provider."),

				identityProviderClientIdClientSecretAttributes("Yahoo"),

				providerAttributeList,
			),

			"amazon": identityProviderSchemaAttribute(
				framework.SchemaAttributeDescriptionFromMarkdown("A single block that specifies options for connectivity to the Amazon social identity provider."),

				identityProviderClientIdClientSecretAttributes("Amazon"),

				providerAttributeList,
			),

			"twitter": identityProviderSchemaAttribute(
				framework.SchemaAttributeDescriptionFromMarkdown("A single block that specifies options for connectivity to the Twitter social identity provider."),

				identityProviderClientIdClientSecretAttributes("Twitter"),

				providerAttributeList,
			),

			"apple": identityProviderSchemaAttribute(
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

			"paypal": identityProviderSchemaAttribute(
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

			"microsoft": identityProviderSchemaAttribute(
				framework.SchemaAttributeDescriptionFromMarkdown("A single block that specifies options for connectivity to the Microsoft social identity provider."),

				identityProviderClientIdClientSecretAttributes("Microsoft"),

				providerAttributeList,
			),

			"github": identityProviderSchemaAttribute(
				framework.SchemaAttributeDescriptionFromMarkdown("A single block that specifies options for connectivity to the Github social identity provider."),

				identityProviderClientIdClientSecretAttributes("Github"),

				providerAttributeList,
			),

			"openid_connect": identityProviderSchemaAttribute(
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

					"pkce_method": schema.StringAttribute{
						Description:         oidcPkceMethodDescription.Description,
						MarkdownDescription: oidcPkceMethodDescription.MarkdownDescription,
						Optional:            true,
						Computed:            true,

						Default: stringdefault.StaticString(string(management.ENUMIDENTITYPROVIDERPKCEMETHOD_NONE)),

						Validators: []validator.String{
							stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumIdentityProviderPKCEMethodEnumValues)...),
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

			"saml": identityProviderSchemaAttribute(
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

					"idp_verification": schema.SingleNestedAttribute{
						Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that specifies settings for SAML IdP verification, including the list of IdP certificates used to verify the signature on the signed assertion of the identity provider.").Description,
						Required:    true,

						Attributes: map[string]schema.Attribute{
							"certificates": schema.SetNestedAttribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("An unordered list that specifies the identity provider's certificate IDs used to verify the signature on the signed assertion from the identity provider. Signing is done with a private key and verified with a public key.").Description,
								Required:    true,

								NestedObject: schema.NestedAttributeObject{

									Attributes: map[string]schema.Attribute{
										"id": schema.StringAttribute{
											Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the identity provider's certificate ID used to verify the signature on the signed assertion from the identity provider.  Must be a valid PingOne resource ID.").Description,
											Required:    true,

											CustomType: pingonetypes.ResourceIDType{},
										},
									},
								},

								Validators: []validator.Set{
									setvalidator.SizeAtLeast(attrMinLength),
								},
							},
						},
					},

					"sp_signing": schema.SingleNestedAttribute{
						Description:         samlSpSigningDescription.Description,
						MarkdownDescription: samlSpSigningDescription.MarkdownDescription,
						Optional:            true,

						Attributes: map[string]schema.Attribute{
							"algorithm": schema.StringAttribute{
								Description:         samlSpSigningAlgorithmDescription.Description,
								MarkdownDescription: samlSpSigningAlgorithmDescription.MarkdownDescription,
								Optional:            true,

								Validators: []validator.String{
									stringvalidator.OneOf(utils.EnumSliceToStringSlice(management.AllowedEnumIdentityProviderSAMLSigningAlgorithmEnumValues)...),
								},
							},

							"key": schema.SingleNestedAttribute{
								Description: framework.SchemaAttributeDescriptionFromMarkdown("A single object that specifies settings for the SAML Sp Signing key.").Description,
								Required:    true,

								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the service provider's signing key ID.  Must be a valid PingOne resource ID.").Description,
										Required:    true,

										CustomType: pingonetypes.ResourceIDType{},
									},
								},
							},
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

func identityProviderSchemaAttribute(description framework.SchemaAttributeDescription, attributes map[string]schema.Attribute, exactlyOneOfBlockNames []string) schema.SingleNestedAttribute {
	description = description.ExactlyOneOf(exactlyOneOfBlockNames).RequiresReplaceNestedAttributes()

	exactlyOneOfPaths := make([]path.Expression, len(exactlyOneOfBlockNames))
	for i, blockName := range exactlyOneOfBlockNames {
		exactlyOneOfPaths[i] = path.MatchRelative().AtParent().AtName(blockName)
	}

	return schema.SingleNestedAttribute{
		Description:         description.Description,
		MarkdownDescription: description.MarkdownDescription,
		Optional:            true,

		Attributes: attributes,

		Validators: []validator.Object{
			objectvalidator.ExactlyOneOf(
				exactlyOneOfPaths...,
			),
		},

		PlanModifiers: []planmodifier.Object{
			objectplanmodifier.RequiresReplace(),
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
		var plan IdentityProviderFacebookResourceModel
		d := p.Facebook.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		idpData := management.IdentityProviderFacebook{
			Enabled:         common.Enabled,
			Name:            common.Name,
			Type:            management.ENUMIDENTITYPROVIDEREXT_FACEBOOK,
			Description:     common.Description,
			Registration:    common.Registration,
			LoginButtonIcon: common.LoginButtonIcon,
			Icon:            common.Icon,
		}

		idpData.SetAppId(plan.AppId.ValueString())
		idpData.SetAppSecret(plan.AppSecret.ValueString())

		data.IdentityProviderFacebook = &idpData
		processedCount += 1
	}

	if !p.Google.IsNull() && !p.Google.IsUnknown() {
		var plan IdentityProviderGoogleResourceModel
		d := p.Google.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		idpData := management.IdentityProviderClientIDClientSecret{
			Enabled:         common.Enabled,
			Name:            common.Name,
			Type:            management.ENUMIDENTITYPROVIDEREXT_GOOGLE,
			Description:     common.Description,
			Registration:    common.Registration,
			LoginButtonIcon: common.LoginButtonIcon,
			Icon:            common.Icon,
		}

		idpData.SetClientId(plan.ClientId.ValueString())
		idpData.SetClientSecret(plan.ClientSecret.ValueString())

		data.IdentityProviderClientIDClientSecret = &idpData
		processedCount += 1
	}

	if !p.LinkedIn.IsNull() && !p.LinkedIn.IsUnknown() {
		var plan IdentityProviderLinkedInResourceModel
		d := p.LinkedIn.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		idpData := management.IdentityProviderClientIDClientSecret{
			Enabled:         common.Enabled,
			Name:            common.Name,
			Type:            management.ENUMIDENTITYPROVIDEREXT_LINKEDIN,
			Description:     common.Description,
			Registration:    common.Registration,
			LoginButtonIcon: common.LoginButtonIcon,
			Icon:            common.Icon,
		}

		idpData.SetClientId(plan.ClientId.ValueString())
		idpData.SetClientSecret(plan.ClientSecret.ValueString())

		data.IdentityProviderClientIDClientSecret = &idpData
		processedCount += 1
	}

	if !p.Yahoo.IsNull() && !p.Yahoo.IsUnknown() {
		var plan IdentityProviderYahooResourceModel
		d := p.Yahoo.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		idpData := management.IdentityProviderClientIDClientSecret{
			Enabled:         common.Enabled,
			Name:            common.Name,
			Type:            management.ENUMIDENTITYPROVIDEREXT_YAHOO,
			Description:     common.Description,
			Registration:    common.Registration,
			LoginButtonIcon: common.LoginButtonIcon,
			Icon:            common.Icon,
		}

		idpData.SetClientId(plan.ClientId.ValueString())
		idpData.SetClientSecret(plan.ClientSecret.ValueString())

		data.IdentityProviderClientIDClientSecret = &idpData
		processedCount += 1
	}

	if !p.Amazon.IsNull() && !p.Amazon.IsUnknown() {
		var plan IdentityProviderAmazonResourceModel
		d := p.Amazon.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		idpData := management.IdentityProviderClientIDClientSecret{
			Enabled:         common.Enabled,
			Name:            common.Name,
			Type:            management.ENUMIDENTITYPROVIDEREXT_AMAZON,
			Description:     common.Description,
			Registration:    common.Registration,
			LoginButtonIcon: common.LoginButtonIcon,
			Icon:            common.Icon,
		}

		idpData.SetClientId(plan.ClientId.ValueString())
		idpData.SetClientSecret(plan.ClientSecret.ValueString())

		data.IdentityProviderClientIDClientSecret = &idpData
		processedCount += 1
	}

	if !p.Twitter.IsNull() && !p.Twitter.IsUnknown() {
		var plan IdentityProviderTwitterResourceModel
		d := p.Twitter.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		idpData := management.IdentityProviderClientIDClientSecret{
			Enabled:         common.Enabled,
			Name:            common.Name,
			Type:            management.ENUMIDENTITYPROVIDEREXT_TWITTER,
			Description:     common.Description,
			Registration:    common.Registration,
			LoginButtonIcon: common.LoginButtonIcon,
			Icon:            common.Icon,
		}

		idpData.SetClientId(plan.ClientId.ValueString())
		idpData.SetClientSecret(plan.ClientSecret.ValueString())

		data.IdentityProviderClientIDClientSecret = &idpData
		processedCount += 1
	}

	if !p.Apple.IsNull() && !p.Apple.IsUnknown() {
		var plan IdentityProviderAppleResourceModel
		d := p.Apple.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		idpData := management.IdentityProviderApple{
			Enabled:         common.Enabled,
			Name:            common.Name,
			Type:            management.ENUMIDENTITYPROVIDEREXT_APPLE,
			Description:     common.Description,
			Registration:    common.Registration,
			LoginButtonIcon: common.LoginButtonIcon,
			Icon:            common.Icon,
		}

		idpData.SetClientId(plan.ClientId.ValueString())
		idpData.SetClientSecretSigningKey(plan.ClientSecretSigningKey.ValueString())
		idpData.SetKeyId(plan.KeyId.ValueString())
		idpData.SetTeamId(plan.TeamId.ValueString())

		data.IdentityProviderApple = &idpData
		processedCount += 1
	}

	if !p.Paypal.IsNull() && !p.Paypal.IsUnknown() {
		var plan IdentityProviderPaypalResourceModel
		d := p.Paypal.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		idpData := management.IdentityProviderPaypal{
			Enabled:         common.Enabled,
			Name:            common.Name,
			Type:            management.ENUMIDENTITYPROVIDEREXT_PAYPAL,
			Description:     common.Description,
			Registration:    common.Registration,
			LoginButtonIcon: common.LoginButtonIcon,
			Icon:            common.Icon,
		}

		idpData.SetClientId(plan.ClientId.ValueString())
		idpData.SetClientSecret(plan.ClientSecret.ValueString())
		idpData.SetClientEnvironment(plan.ClientEnvironment.ValueString())

		data.IdentityProviderPaypal = &idpData
		processedCount += 1
	}

	if !p.Microsoft.IsNull() && !p.Microsoft.IsUnknown() {
		var plan IdentityProviderMicrosoftResourceModel
		d := p.Microsoft.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		idpData := management.IdentityProviderClientIDClientSecret{
			Enabled:         common.Enabled,
			Name:            common.Name,
			Type:            management.ENUMIDENTITYPROVIDEREXT_MICROSOFT,
			Description:     common.Description,
			Registration:    common.Registration,
			LoginButtonIcon: common.LoginButtonIcon,
			Icon:            common.Icon,
		}

		idpData.SetClientId(plan.ClientId.ValueString())
		idpData.SetClientSecret(plan.ClientSecret.ValueString())

		data.IdentityProviderClientIDClientSecret = &idpData
		processedCount += 1
	}

	if !p.Github.IsNull() && !p.Github.IsUnknown() {
		var plan IdentityProviderGithubResourceModel
		d := p.Github.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		idpData := management.IdentityProviderClientIDClientSecret{
			Enabled:         common.Enabled,
			Name:            common.Name,
			Type:            management.ENUMIDENTITYPROVIDEREXT_GITHUB,
			Description:     common.Description,
			Registration:    common.Registration,
			LoginButtonIcon: common.LoginButtonIcon,
			Icon:            common.Icon,
		}

		idpData.SetClientId(plan.ClientId.ValueString())
		idpData.SetClientSecret(plan.ClientSecret.ValueString())

		data.IdentityProviderClientIDClientSecret = &idpData
		processedCount += 1
	}

	if !p.OpenIDConnect.IsNull() && !p.OpenIDConnect.IsUnknown() {
		var plan IdentityProviderOIDCResourceModel
		d := p.OpenIDConnect.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		idpData := management.IdentityProviderOIDC{
			Enabled:         common.Enabled,
			Name:            common.Name,
			Type:            management.ENUMIDENTITYPROVIDEREXT_OPENID_CONNECT,
			Description:     common.Description,
			Registration:    common.Registration,
			LoginButtonIcon: common.LoginButtonIcon,
			Icon:            common.Icon,
		}

		if !plan.AuthorizationEndpoint.IsNull() && !plan.AuthorizationEndpoint.IsUnknown() {
			idpData.SetAuthorizationEndpoint(plan.AuthorizationEndpoint.ValueString())
		}

		if !plan.ClientId.IsNull() && !plan.ClientId.IsUnknown() {
			idpData.SetClientId(plan.ClientId.ValueString())
		}

		if !plan.ClientSecret.IsNull() && !plan.ClientSecret.IsUnknown() {
			idpData.SetClientSecret(plan.ClientSecret.ValueString())
		}

		if !plan.DiscoveryEndpoint.IsNull() && !plan.DiscoveryEndpoint.IsUnknown() {
			idpData.SetDiscoveryEndpoint(plan.DiscoveryEndpoint.ValueString())
		}

		if !plan.Issuer.IsNull() && !plan.Issuer.IsUnknown() {
			idpData.SetIssuer(plan.Issuer.ValueString())
		}

		if !plan.PkceMethod.IsNull() && !plan.PkceMethod.IsUnknown() {
			idpData.SetPkceMethod(management.EnumIdentityProviderPKCEMethod(plan.PkceMethod.ValueString()))
		}

		if !plan.JwksEndpoint.IsNull() && !plan.JwksEndpoint.IsUnknown() {
			idpData.SetJwksEndpoint(plan.JwksEndpoint.ValueString())
		}

		if !plan.Scopes.IsNull() && !plan.Scopes.IsUnknown() {
			var scopesPlan []string
			diags.Append(plan.Scopes.ElementsAs(ctx, &scopesPlan, false)...)
			if diags.HasError() {
				return nil, diags
			}

			idpData.SetScopes(scopesPlan)
		}

		if !plan.TokenEndpoint.IsNull() && !plan.TokenEndpoint.IsUnknown() {
			idpData.SetTokenEndpoint(plan.TokenEndpoint.ValueString())
		}

		if !plan.TokenEndpointAuthMethod.IsNull() && !plan.TokenEndpointAuthMethod.IsUnknown() {
			idpData.SetTokenEndpointAuthMethod(management.EnumIdentityProviderOIDCTokenAuthMethod(plan.TokenEndpointAuthMethod.ValueString()))
		}

		if !plan.UserinfoEndpoint.IsNull() && !plan.UserinfoEndpoint.IsUnknown() {
			idpData.SetUserInfoEndpoint(plan.UserinfoEndpoint.ValueString())
		}

		data.IdentityProviderOIDC = &idpData
		processedCount += 1
	}

	if !p.Saml.IsNull() && !p.Saml.IsUnknown() {
		var plan IdentityProviderSAMLResourceModel
		d := p.Saml.As(ctx, &plan, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		idpData := management.IdentityProviderSAML{
			Enabled:         common.Enabled,
			Name:            common.Name,
			Type:            management.ENUMIDENTITYPROVIDEREXT_SAML,
			Description:     common.Description,
			Registration:    common.Registration,
			LoginButtonIcon: common.LoginButtonIcon,
			Icon:            common.Icon,
		}

		if !plan.AuthenticationRequestSigned.IsNull() && !plan.AuthenticationRequestSigned.IsUnknown() {
			idpData.SetAuthnRequestSigned(plan.AuthenticationRequestSigned.ValueBool())
		}

		if !plan.IdpEntityId.IsNull() && !plan.IdpEntityId.IsUnknown() {
			idpData.SetIdpEntityId(plan.IdpEntityId.ValueString())
		}

		if !plan.SpEntityId.IsNull() && !plan.SpEntityId.IsUnknown() {
			idpData.SetSpEntityId(plan.SpEntityId.ValueString())
		}

		if !plan.IdpVerification.IsNull() && !plan.IdpVerification.IsUnknown() {
			var idpVerificationPlan IdentityProviderSAMLResourceIdPVerificationModel
			diags.Append(plan.IdpVerification.As(ctx, &idpVerificationPlan, basetypes.ObjectAsOptions{
				UnhandledNullAsEmpty:    false,
				UnhandledUnknownAsEmpty: false,
			})...)
			if diags.HasError() {
				return nil, diags
			}

			idpVerification, d := idpVerificationPlan.expand(ctx)
			diags.Append(d...)
			if diags.HasError() {
				return nil, diags
			}

			idpData.SetIdpVerification(*idpVerification)
		}

		if !plan.SpSigning.IsNull() && !plan.SpSigning.IsUnknown() {
			var spSigningPlan IdentityProviderSAMLResourceSpSigningModel
			diags.Append(plan.SpSigning.As(ctx, &spSigningPlan, basetypes.ObjectAsOptions{
				UnhandledNullAsEmpty:    false,
				UnhandledUnknownAsEmpty: false,
			})...)
			if diags.HasError() {
				return nil, diags
			}

			spSigning, d := spSigningPlan.expand(ctx)
			diags.Append(d...)
			if diags.HasError() {
				return nil, diags
			}

			idpData.SetSpSigning(*spSigning)
		}

		if !plan.SsoBinding.IsNull() && !plan.SsoBinding.IsUnknown() {
			idpData.SetSsoBinding(management.EnumIdentityProviderSAMLSSOBinding(plan.SsoBinding.ValueString()))
		}

		if !plan.SsoEndpoint.IsNull() && !plan.SsoEndpoint.IsUnknown() {
			idpData.SetSsoEndpoint(plan.SsoEndpoint.ValueString())
		}

		if !plan.SloBinding.IsNull() && !plan.SloBinding.IsUnknown() {
			idpData.SetSloBinding(management.EnumIdentityProviderSAMLSLOBinding(plan.SloBinding.ValueString()))
		}

		if !plan.SloEndpoint.IsNull() && !plan.SloEndpoint.IsUnknown() {
			idpData.SetSloEndpoint(plan.SloEndpoint.ValueString())
		}

		if !plan.SloResponseEndpoint.IsNull() && !plan.SloResponseEndpoint.IsUnknown() {
			idpData.SetSloResponseEndpoint(plan.SloResponseEndpoint.ValueString())
		}

		if !plan.SloWindow.IsNull() && !plan.SloWindow.IsUnknown() {
			idpData.SetSloWindow(int32(plan.SloWindow.ValueInt64()))
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

func (p *IdentityProviderSAMLResourceIdPVerificationModel) expand(ctx context.Context) (*management.IdentityProviderSAMLAllOfIdpVerification, diag.Diagnostics) {
	var diags diag.Diagnostics

	var certificatesPlan []IdentityProviderSAMLResourceIdPVerificationCertificatesModel
	diags.Append(p.Certificates.ElementsAs(ctx, &certificatesPlan, false)...)
	if diags.HasError() {
		return nil, diags
	}

	certificates := make([]management.IdentityProviderSAMLAllOfIdpVerificationCertificates, 0)
	for _, certificatePlan := range certificatesPlan {
		certificate := management.NewIdentityProviderSAMLAllOfIdpVerificationCertificates(certificatePlan.Id.ValueString())
		certificates = append(certificates, *certificate)
	}

	data := management.NewIdentityProviderSAMLAllOfIdpVerification(
		certificates,
	)

	return data, diags
}

func (p *IdentityProviderSAMLResourceSpSigningModel) expand(ctx context.Context) (*management.IdentityProviderSAMLAllOfSpSigning, diag.Diagnostics) {
	var diags diag.Diagnostics

	var keyPlan IdentityProviderSAMLResourceSpSigningKeyModel
	diags.Append(p.Key.As(ctx, &keyPlan, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if diags.HasError() {
		return nil, diags
	}

	key := management.NewIdentityProviderSAMLAllOfSpSigningKey(keyPlan.Id.ValueString())

	data := management.NewIdentityProviderSAMLAllOfSpSigning(
		*key,
	)

	if !p.Algorithm.IsNull() && !p.Algorithm.IsUnknown() {
		data.SetAlgorithm(management.EnumIdentityProviderSAMLSigningAlgorithm(p.Algorithm.ValueString()))
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

	p.Id = framework.PingOneResourceIDOkToTF(common.GetIdOk())
	p.EnvironmentId = framework.PingOneResourceIDOkToTF(common.Environment.GetIdOk())
	p.Name = framework.StringOkToTF(common.GetNameOk())
	p.Description = framework.StringOkToTF(common.GetDescriptionOk())
	p.Enabled = framework.BoolOkToTF(common.GetEnabledOk())

	p.RegistrationPopulationId = pingonetypes.NewResourceIDNull()
	if v, ok := common.GetRegistrationOk(); ok {
		if q, ok := v.GetPopulationOk(); ok {
			p.RegistrationPopulationId = framework.PingOneResourceIDOkToTF(q.GetIdOk())
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

func identityProviderFacebookToTF(idpApiObject *management.IdentityProviderFacebook) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if idpApiObject == nil || idpApiObject.GetType() != management.ENUMIDENTITYPROVIDEREXT_FACEBOOK {
		return types.ObjectNull(identityProviderFacebookTFObjectTypes), diags
	}

	attributesMap := map[string]attr.Value{
		"app_id":     framework.StringOkToTF(idpApiObject.GetAppIdOk()),
		"app_secret": framework.StringOkToTF(idpApiObject.GetAppSecretOk()),
	}

	returnVar, d := types.ObjectValue(identityProviderFacebookTFObjectTypes, attributesMap)
	diags.Append(d...)

	return returnVar, diags
}

func identityProviderClientIDClientSecretToTF(idpApiObject *management.IdentityProviderClientIDClientSecret, idpType management.EnumIdentityProviderExt) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if idpApiObject == nil || idpApiObject.GetType() != idpType {
		return types.ObjectNull(identityProviderClientIDClientSecretTFObjectTypes), diags
	}

	attributesMap := map[string]attr.Value{
		"client_id":     framework.StringOkToTF(idpApiObject.GetClientIdOk()),
		"client_secret": framework.StringOkToTF(idpApiObject.GetClientSecretOk()),
	}

	returnVar, d := types.ObjectValue(identityProviderClientIDClientSecretTFObjectTypes, attributesMap)
	diags.Append(d...)

	return returnVar, diags
}

func identityProviderAppleToTF(idpApiObject *management.IdentityProviderApple) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if idpApiObject == nil || idpApiObject.GetType() != management.ENUMIDENTITYPROVIDEREXT_APPLE {
		return types.ObjectNull(identityProviderAppleTFObjectTypes), diags
	}

	attributesMap := map[string]attr.Value{
		"team_id":                   framework.StringOkToTF(idpApiObject.GetTeamIdOk()),
		"key_id":                    framework.StringOkToTF(idpApiObject.GetKeyIdOk()),
		"client_id":                 framework.StringOkToTF(idpApiObject.GetClientIdOk()),
		"client_secret_signing_key": framework.StringOkToTF(idpApiObject.GetClientSecretSigningKeyOk()),
	}

	returnVar, d := types.ObjectValue(identityProviderAppleTFObjectTypes, attributesMap)
	diags.Append(d...)

	return returnVar, diags
}

func identityProviderPaypalToTF(idpApiObject *management.IdentityProviderPaypal) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if idpApiObject == nil || idpApiObject.GetType() != management.ENUMIDENTITYPROVIDEREXT_PAYPAL {
		return types.ObjectNull(identityProviderPaypalTFObjectTypes), diags
	}

	attributesMap := map[string]attr.Value{
		"client_id":          framework.StringOkToTF(idpApiObject.GetClientIdOk()),
		"client_secret":      framework.StringOkToTF(idpApiObject.GetClientSecretOk()),
		"client_environment": framework.StringOkToTF(idpApiObject.GetClientEnvironmentOk()),
	}

	returnVar, d := types.ObjectValue(identityProviderPaypalTFObjectTypes, attributesMap)
	diags.Append(d...)

	return returnVar, diags
}

func identityProviderOIDCToTF(idpApiObject *management.IdentityProviderOIDC) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if idpApiObject == nil || idpApiObject.GetType() != management.ENUMIDENTITYPROVIDEREXT_OPENID_CONNECT {
		return types.ObjectNull(identityProviderOIDCTFObjectTypes), diags
	}

	attributesMap := map[string]attr.Value{
		"authorization_endpoint":     framework.StringOkToTF(idpApiObject.GetAuthorizationEndpointOk()),
		"client_id":                  framework.StringOkToTF(idpApiObject.GetClientIdOk()),
		"client_secret":              framework.StringOkToTF(idpApiObject.GetClientSecretOk()),
		"discovery_endpoint":         framework.StringOkToTF(idpApiObject.GetDiscoveryEndpointOk()),
		"issuer":                     framework.StringOkToTF(idpApiObject.GetIssuerOk()),
		"pkce_method":                framework.EnumOkToTF(idpApiObject.GetPkceMethodOk()),
		"jwks_endpoint":              framework.StringOkToTF(idpApiObject.GetJwksEndpointOk()),
		"scopes":                     framework.StringSetOkToTF(idpApiObject.GetScopesOk()),
		"token_endpoint":             framework.StringOkToTF(idpApiObject.GetTokenEndpointOk()),
		"token_endpoint_auth_method": framework.EnumOkToTF(idpApiObject.GetTokenEndpointAuthMethodOk()),
		"userinfo_endpoint":          framework.StringOkToTF(idpApiObject.GetUserInfoEndpointOk()),
	}

	returnVar, d := types.ObjectValue(identityProviderOIDCTFObjectTypes, attributesMap)
	diags.Append(d...)

	return returnVar, diags
}

func identityProviderSAMLToTF(idpApiObject *management.IdentityProviderSAML) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if idpApiObject == nil || idpApiObject.GetType() != management.ENUMIDENTITYPROVIDEREXT_SAML {
		return types.ObjectNull(identityProviderSAMLTFObjectTypes), diags
	}

	idpVerification, d := identityProviderSAMLIdPVerificationOkToTF(idpApiObject.GetIdpVerificationOk())
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(identityProviderSAMLTFObjectTypes), diags
	}

	spSigning, d := identityProviderSAMLSpSigningOkToTF(idpApiObject.GetSpSigningOk())
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(identityProviderSAMLTFObjectTypes), diags
	}

	attributesMap := map[string]attr.Value{
		"authentication_request_signed": framework.BoolOkToTF(idpApiObject.GetAuthnRequestSignedOk()),
		"idp_entity_id":                 framework.StringOkToTF(idpApiObject.GetIdpEntityIdOk()),
		"sp_entity_id":                  framework.StringOkToTF(idpApiObject.GetSpEntityIdOk()),
		"idp_verification":              idpVerification,
		"sp_signing":                    spSigning,
		"sso_binding":                   framework.EnumOkToTF(idpApiObject.GetSsoBindingOk()),
		"sso_endpoint":                  framework.StringOkToTF(idpApiObject.GetSsoEndpointOk()),
		"slo_binding":                   framework.EnumOkToTF(idpApiObject.GetSloBindingOk()),
		"slo_endpoint":                  framework.StringOkToTF(idpApiObject.GetSloEndpointOk()),
		"slo_response_endpoint":         framework.StringOkToTF(idpApiObject.GetSloResponseEndpointOk()),
		"slo_window":                    framework.Int32OkToTF(idpApiObject.GetSloWindowOk()),
	}

	returnVar, d := types.ObjectValue(identityProviderSAMLTFObjectTypes, attributesMap)
	diags.Append(d...)

	return returnVar, diags
}

func identityProviderSAMLIdPVerificationOkToTF(apiObject *management.IdentityProviderSAMLAllOfIdpVerification, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(identityProviderSAMLIdPVerificationTFObjectTypes), diags
	}

	certificates, d := identityProviderSAMLIdPVerificationCertificatesOkToTF(apiObject.GetCertificatesOk())
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(identityProviderSAMLIdPVerificationTFObjectTypes), diags
	}

	attributesMap := map[string]attr.Value{
		"certificates": certificates,
	}

	returnVar, d := types.ObjectValue(identityProviderSAMLIdPVerificationTFObjectTypes, attributesMap)
	diags.Append(d...)

	return returnVar, diags
}

func identityProviderSAMLIdPVerificationCertificatesOkToTF(apiObject []management.IdentityProviderSAMLAllOfIdpVerificationCertificates, ok bool) (types.Set, diag.Diagnostics) {
	var diags diag.Diagnostics

	tfObjType := types.ObjectType{AttrTypes: identityProviderSAMLIdPVerificationCertificateTFObjectTypes}

	if !ok || len(apiObject) == 0 {
		return types.SetNull(tfObjType), diags
	}

	flattenedList := []attr.Value{}
	for _, v := range apiObject {

		objMap := map[string]attr.Value{
			"id": framework.PingOneResourceIDOkToTF(v.GetIdOk()),
		}

		flattenedObj, d := types.ObjectValue(identityProviderSAMLIdPVerificationCertificateTFObjectTypes, objMap)
		diags.Append(d...)

		flattenedList = append(flattenedList, flattenedObj)
	}

	returnVar, d := types.SetValue(tfObjType, flattenedList)
	diags.Append(d...)

	return returnVar, diags
}

func identityProviderSAMLSpSigningOkToTF(apiObject *management.IdentityProviderSAMLAllOfSpSigning, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(identityProviderSAMLSpSigningTFObjectTypes), diags
	}

	key, d := identityProviderSAMLSpSigningKeyOkToTF(apiObject.GetKeyOk())
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(identityProviderSAMLSpSigningTFObjectTypes), diags
	}

	attributesMap := map[string]attr.Value{
		"key":       key,
		"algorithm": framework.EnumOkToTF(apiObject.GetAlgorithmOk()),
	}

	returnVar, d := types.ObjectValue(identityProviderSAMLSpSigningTFObjectTypes, attributesMap)
	diags.Append(d...)

	return returnVar, diags
}

func identityProviderSAMLSpSigningKeyOkToTF(apiObject *management.IdentityProviderSAMLAllOfSpSigningKey, ok bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !ok || apiObject == nil {
		return types.ObjectNull(identityProviderSAMLSpSigningKeyTFObjectTypes), diags
	}

	attributesMap := map[string]attr.Value{
		"id": framework.PingOneResourceIDOkToTF(apiObject.GetIdOk()),
	}

	returnVar, d := types.ObjectValue(identityProviderSAMLSpSigningKeyTFObjectTypes, attributesMap)
	diags.Append(d...)

	return returnVar, diags
}
