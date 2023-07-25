package sso

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

func ResourceIdentityProvider() *schema.Resource {

	providerAttributeList := []string{"facebook", "google", "linkedin", "yahoo", "amazon", "twitter", "apple", "paypal", "microsoft", "github", "openid_connect", "saml"}

	return &schema.Resource{

		// This description is used by the documentation generator and the language server.
		Description: "Resource to create and manage PingOne Identity Providers.",

		CreateContext: resourceIdentityProviderCreate,
		ReadContext:   resourceIdentityProviderRead,
		UpdateContext: resourceIdentityProviderUpdate,
		DeleteContext: resourceIdentityProviderDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceIdentityProviderImport,
		},

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Description:      "The ID of the environment to create the identity provider in.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
				ForceNew:         true,
			},
			"name": {
				Description:      "A string that specifies the name of the identity provider.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
			},
			"description": {
				Description: "A string that specifies the description of the identity provider.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"enabled": {
				Description: "A boolean that specifies whether the identity provider is enabled in the environment.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"registration_population_id": {
				Description:      "Setting this attribute gives management of linked users to the IdP and also triggers just-in-time provisioning of new users to the population ID provided.",
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
			},
			"login_button_icon": {
				Description: "The HREF and the ID for the identity provider icon to use as the login button.",
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description:      "The ID for the identity provider icon to use as the login button.",
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
						},
						"href": {
							Description:      "The HREF for the identity provider icon to use as the login button.",
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IsURLWithHTTPS),
						},
					},
				},
			},
			"icon": {
				Description: "The HREF and the ID for the identity provider icon.",
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description:      "The ID for the identity provider icon.",
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
						},
						"href": {
							Description:      "The HREF for the identity provider icon.",
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IsURLWithHTTPS),
						},
					},
				},
			},

			// The providers
			"facebook": {
				Description:  "Options for Identity provider connectivity to Facebook.",
				Type:         schema.TypeList,
				MaxItems:     1,
				Optional:     true,
				ExactlyOneOf: providerAttributeList,
				ForceNew:     true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"app_id": {
							Description:      "A string that specifies the application ID from Facebook.",
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
						},
						"app_secret": {
							Description:      "A string that specifies the application secret from Facebook.",
							Type:             schema.TypeString,
							Required:         true,
							Sensitive:        true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
						},
					},
				},
			},
			"google": {
				Description:  "Options for Identity provider connectivity to Google.",
				Type:         schema.TypeList,
				MaxItems:     1,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: providerAttributeList,
				Elem:         clientIdClientSecretSchema("Google"),
			},
			"linkedin": {
				Description:  "Options for Identity provider connectivity to LinkedIn.",
				Type:         schema.TypeList,
				MaxItems:     1,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: providerAttributeList,
				Elem:         clientIdClientSecretSchema("LinkedIn"),
			},
			"yahoo": {
				Description:  "Options for Identity provider connectivity to Yahoo.",
				Type:         schema.TypeList,
				MaxItems:     1,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: providerAttributeList,
				Elem:         clientIdClientSecretSchema("Yahoo"),
			},
			"amazon": {
				Description:  "Options for Identity provider connectivity to Amazon.",
				Type:         schema.TypeList,
				MaxItems:     1,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: providerAttributeList,
				Elem:         clientIdClientSecretSchema("Amazon"),
			},
			"twitter": {
				Description:  "Options for Identity provider connectivity to Twitter.",
				Type:         schema.TypeList,
				MaxItems:     1,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: providerAttributeList,
				Elem:         clientIdClientSecretSchema("Twitter"),
			},
			"apple": {
				Description:  "Options for Identity provider connectivity to Apple.",
				Type:         schema.TypeList,
				MaxItems:     1,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: providerAttributeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"client_id": {
							Description:      "A string that specifies the application ID from Apple. This is the identifier obtained after registering a services ID in the Apple developer portal.",
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
						},
						"client_secret_signing_key": {
							Description:      "A string that specifies the private key that is used to generate a client secret.",
							Type:             schema.TypeString,
							Required:         true,
							Sensitive:        true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
						},
						"key_id": {
							Description:      "A 10-character string that Apple uses to identify an authentication key.",
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(10, 10)),
						},
						"team_id": {
							Description:      "A 10-character string that Apple uses to identify teams.",
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(10, 10)),
						},
					},
				},
			},
			"paypal": {
				Description:  "Options for Identity provider connectivity to Paypal.",
				Type:         schema.TypeList,
				MaxItems:     1,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: providerAttributeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"client_id": {
							Description:      "A string that specifies the application ID from PayPal.",
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
						},
						"client_secret": {
							Description:      "A string that specifies the application secret from PayPal.",
							Type:             schema.TypeString,
							Required:         true,
							Sensitive:        true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
						},
						"client_environment": {
							Description:      "A string that specifies the PayPal environment. Options are `sandbox`, and `live`.",
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"sandbox", "live"}, false)),
						},
					},
				},
			},
			"microsoft": {
				Description:  "Options for Identity provider connectivity to Microsoft.",
				Type:         schema.TypeList,
				MaxItems:     1,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: providerAttributeList,
				Elem:         clientIdClientSecretSchema("Microsoft"),
			},
			"github": {
				Description:  "Options for Identity provider connectivity to Github.",
				Type:         schema.TypeList,
				MaxItems:     1,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: providerAttributeList,
				Elem:         clientIdClientSecretSchema("Github"),
			},
			"openid_connect": {
				Description:  "Options for Identity provider connectivity to a generic OpenID Connect service.",
				Type:         schema.TypeList,
				MaxItems:     1,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: providerAttributeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"authorization_endpoint": {
							Description:      "A string that specifies the the OIDC identity provider's authorization endpoint. This value must be a URL that uses https.",
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IsURLWithHTTPS),
						},
						"client_id": {
							Description:      "A string that specifies the application ID from the OIDC identity provider.",
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
						},
						"client_secret": {
							Description:      "A string that specifies the application secret from the OIDC identity provider.",
							Type:             schema.TypeString,
							Required:         true,
							Sensitive:        true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
						},
						"discovery_endpoint": {
							Description:      "A string that specifies the OIDC identity provider's discovery endpoint. This value must be a URL that uses https.",
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IsURLWithHTTPS),
						},
						"issuer": {
							Description:      "A string that specifies the issuer to which the authentication is sent for the OIDC identity provider. This value must be a URL that uses https.",
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IsURLWithHTTPS),
						},
						"jwks_endpoint": {
							Description:      "A string that specifies the OIDC identity provider's jwks endpoint. This value must be a URL that uses https.",
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
						},
						"scopes": {
							Description: "An array that specifies the scopes to include in the authentication request to the OIDC identity provider.",
							Type:        schema.TypeSet,
							Required:    true,
							Elem: &schema.Schema{
								Type:             schema.TypeString,
								ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
							},
						},
						"token_endpoint": {
							Description:      "A string that specifies the OIDC identity provider's token endpoint.",
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IsURLWithHTTPS),
						},
						"token_endpoint_auth_method": {
							Description:      fmt.Sprintf("A string that specifies the OIDC identity provider's token endpoint authentication method. Options are `%s` (default), `%s`, and `%s`.", string(management.ENUMIDENTITYPROVIDEROIDCTOKENAUTHMETHOD_CLIENT_SECRET_BASIC), string(management.ENUMIDENTITYPROVIDEROIDCTOKENAUTHMETHOD_CLIENT_SECRET_POST), string(management.ENUMIDENTITYPROVIDEROIDCTOKENAUTHMETHOD_NONE)),
							Type:             schema.TypeString,
							Optional:         true,
							Default:          string(management.ENUMIDENTITYPROVIDEROIDCTOKENAUTHMETHOD_CLIENT_SECRET_BASIC),
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{string(management.ENUMIDENTITYPROVIDEROIDCTOKENAUTHMETHOD_CLIENT_SECRET_BASIC), string(management.ENUMIDENTITYPROVIDEROIDCTOKENAUTHMETHOD_CLIENT_SECRET_POST), string(management.ENUMIDENTITYPROVIDEROIDCTOKENAUTHMETHOD_NONE)}, false)),
						},
						"userinfo_endpoint": {
							Description:      "A string that specifies the OIDC identity provider's userInfo endpoint.",
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IsURLWithHTTPS),
						},
					},
				},
			},
			"saml": {
				Description:  "Options for Identity provider connectivity to a generic SAML service.",
				Type:         schema.TypeList,
				MaxItems:     1,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: providerAttributeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"authentication_request_signed": {
							Description: "A boolean that specifies whether the SAML authentication request will be signed when sending to the identity provider. Set this to true if the external IDP is included in an authentication policy to be used by applications that are accessed using a mix of default URLS and custom Domains URLs.",
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
						},
						"idp_entity_id": {
							Description: "A string that specifies the entity ID URI that is checked against the issuerId tag in the incoming response.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"sp_entity_id": {
							Description: "A string that specifies the service provider's entity ID, used to look up the application.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"idp_verification_certificate_ids": {
							Description: "A list that specifies the identity provider's certificate IDs used to verify the signature on the signed assertion from the identity provider. Signing is done with a private key and verified with a public key.",
							Type:        schema.TypeSet,
							Required:    true,
							Elem: &schema.Schema{
								Type:             schema.TypeString,
								ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
							},
						},
						"sp_signing_key_id": {
							Description:      "A string that specifies the service provider's signing key ID.",
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(verify.ValidP1ResourceID),
						},
						"sso_binding": {
							Description:      fmt.Sprintf("A string that specifies the binding for the authentication request. Options are `%s` and `%s`.", string(management.ENUMIDENTITYPROVIDERSAMLSSOBINDING_POST), string(management.ENUMIDENTITYPROVIDERSAMLSSOBINDING_REDIRECT)),
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{string(management.ENUMIDENTITYPROVIDERSAMLSSOBINDING_POST), string(management.ENUMIDENTITYPROVIDERSAMLSSOBINDING_REDIRECT)}, false)),
						},
						"sso_endpoint": {
							Description:      "A string that specifies the SSO endpoint for the authentication request.",
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IsURLWithHTTPorHTTPS),
						},
						"slo_binding": {
							Description:      fmt.Sprintf("A string that specifies the binding protocol to be used for the logout response. Options are `%s` and `%s`.  Existing configurations with no data default to `%s`.", string(management.ENUMIDENTITYPROVIDERSAMLSLOBINDING_REDIRECT), string(management.ENUMIDENTITYPROVIDERSAMLSLOBINDING_POST), string(management.ENUMIDENTITYPROVIDERSAMLSLOBINDING_POST)),
							Type:             schema.TypeString,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{string(management.ENUMIDENTITYPROVIDERSAMLSLOBINDING_REDIRECT), string(management.ENUMIDENTITYPROVIDERSAMLSLOBINDING_POST)}, false)),
							Optional:         true,
							Default:          string(management.ENUMIDENTITYPROVIDERSAMLSLOBINDING_POST),
						},
						"slo_endpoint": {
							Description:      "A string that specifies the logout endpoint URL. This is an optional property. However, if a logout endpoint URL is not defined, logout actions result in an error.",
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IsURLWithHTTPorHTTPS),
						},
						"slo_response_endpoint": {
							Description:      "A string that specifies the endpoint URL to submit the logout response. If a value is not provided, the `slo_endpoint` property value is used to submit SLO response.",
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IsURLWithHTTPorHTTPS),
						},
						"slo_window": {
							Description:      "An integer that defines how long (hours) PingOne can exchange logout messages with the application, specifically a logout request from the application, since the initial request. The minimum value is `1` hour and the maximum is `24` hours.",
							Type:             schema.TypeInt,
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(1, 24)),
						},
					},
				},
			},
		},
	}
}

func clientIdClientSecretSchema(providerName string) *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"client_id": {
				Description:      fmt.Sprintf("A string that specifies the application ID from %s.", providerName),
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
			},
			"client_secret": {
				Description:      fmt.Sprintf("A string that specifies the application secret from %s.", providerName),
				Type:             schema.TypeString,
				Required:         true,
				Sensitive:        true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
			},
		},
	}
}

func resourceIdentityProviderCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient

	var diags diag.Diagnostics

	idpRequest, diags := expandIdentityProvider(d)
	if diags.HasError() {
		return diags
	}

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return apiClient.IdentityProvidersApi.CreateIdentityProvider(ctx, d.Get("environment_id").(string)).IdentityProvider(*idpRequest).Execute()
		},
		"CreateIdentityProvider",
		sdk.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	if diags.HasError() {
		return diags
	}

	respObject := resp.(*management.IdentityProvider)

	if respObject.IdentityProviderApple != nil && respObject.IdentityProviderApple.GetId() != "" {
		d.SetId(respObject.IdentityProviderApple.GetId())
	} else if respObject.IdentityProviderClientIDClientSecret != nil && respObject.IdentityProviderClientIDClientSecret.GetId() != "" {
		d.SetId(respObject.IdentityProviderClientIDClientSecret.GetId())
	} else if respObject.IdentityProviderFacebook != nil && respObject.IdentityProviderFacebook.GetId() != "" {
		d.SetId(respObject.IdentityProviderFacebook.GetId())
	} else if respObject.IdentityProviderOIDC != nil && respObject.IdentityProviderOIDC.GetId() != "" {
		d.SetId(respObject.IdentityProviderOIDC.GetId())
	} else if respObject.IdentityProviderPaypal != nil && respObject.IdentityProviderPaypal.GetId() != "" {
		d.SetId(respObject.IdentityProviderPaypal.GetId())
	} else if respObject.IdentityProviderSAML != nil && respObject.IdentityProviderSAML.GetId() != "" {
		d.SetId(respObject.IdentityProviderSAML.GetId())
	} else {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Cannot determine ID from API response for identity provider: %s", d.Get("name")),
			Detail:   fmt.Sprintf("Full response object: %v\n", resp),
		})

		return diags
	}

	return resourceIdentityProviderRead(ctx, d, meta)
}

func resourceIdentityProviderRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient

	var diags diag.Diagnostics

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return apiClient.IdentityProvidersApi.ReadOneIdentityProvider(ctx, d.Get("environment_id").(string), d.Id()).Execute()
		},
		"ReadOneIdentityProvider",
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

	respObject := resp.(*management.IdentityProvider)

	// flatten
	values := map[string]interface{}{
		"name":                       nil,
		"description":                nil,
		"enabled":                    nil,
		"registration_population_id": nil,
		"login_button_icon":          nil,
		"icon":                       nil,
		"facebook":                   nil,
		"google":                     nil,
		"linkedin":                   nil,
		"yahoo":                      nil,
		"amazon":                     nil,
		"twitter":                    nil,
		"apple":                      nil,
		"paypal":                     nil,
		"microsoft":                  nil,
		"github":                     nil,
		"openid_connect":             nil,
		"saml":                       nil,
	}

	switch respObject.GetActualInstance().(type) {
	case *management.IdentityProviderClientIDClientSecret:
		idpObject := respObject.IdentityProviderClientIDClientSecret

		var schemaAttribute string
		switch idpObject.GetType() {
		case management.ENUMIDENTITYPROVIDEREXT_AMAZON:
			schemaAttribute = "amazon"
		case management.ENUMIDENTITYPROVIDEREXT_GITHUB:
			schemaAttribute = "github"
		case management.ENUMIDENTITYPROVIDEREXT_GOOGLE:
			schemaAttribute = "google"
		case management.ENUMIDENTITYPROVIDEREXT_LINKEDIN:
			schemaAttribute = "linkedin"
		case management.ENUMIDENTITYPROVIDEREXT_MICROSOFT:
			schemaAttribute = "microsoft"
		case management.ENUMIDENTITYPROVIDEREXT_TWITTER:
			schemaAttribute = "twitter"
		case management.ENUMIDENTITYPROVIDEREXT_YAHOO:
			schemaAttribute = "yahoo"
		}

		values["name"] = idpObject.GetName()

		values["enabled"] = idpObject.GetEnabled()

		if v, ok := idpObject.GetDescriptionOk(); ok {
			values["description"] = v
		}

		if v, ok := idpObject.GetRegistrationOk(); ok {
			values["registration_population_id"] = v.GetPopulation().Id
		}

		if v, ok := idpObject.GetLoginButtonIconOk(); ok {
			values["login_button_icon"] = flattenLoginButtonIcon(v)
		}

		if v, ok := idpObject.GetIconOk(); ok {
			values["icon"] = flattenIdPIcon(v)
		}

		values[schemaAttribute] = flattenClientIdClientSecret(idpObject.GetClientId(), idpObject.GetClientSecret())

	case *management.IdentityProviderApple:
		idpObject := respObject.IdentityProviderApple
		schemaAttribute := "apple"

		values["name"] = idpObject.GetName()

		values["enabled"] = idpObject.GetEnabled()

		if v, ok := idpObject.GetDescriptionOk(); ok {
			values["description"] = v
		}

		if v, ok := idpObject.GetRegistrationOk(); ok {
			values["registration_population_id"] = v.GetPopulation().Id
		}

		if v, ok := idpObject.GetLoginButtonIconOk(); ok {
			values["login_button_icon"] = flattenLoginButtonIcon(v)
		}

		if v, ok := idpObject.GetIconOk(); ok {
			values["icon"] = flattenIdPIcon(v)
		}

		values[schemaAttribute] = flattenApple(idpObject.GetClientId(), idpObject.GetClientSecretSigningKey(), idpObject.GetKeyId(), idpObject.GetTeamId())

	case *management.IdentityProviderFacebook:
		idpObject := respObject.IdentityProviderFacebook
		schemaAttribute := "facebook"

		values["name"] = idpObject.GetName()

		values["enabled"] = idpObject.GetEnabled()

		if v, ok := idpObject.GetDescriptionOk(); ok {
			values["description"] = v
		}

		if v, ok := idpObject.GetRegistrationOk(); ok {
			values["registration_population_id"] = v.GetPopulation().Id
		}

		if v, ok := idpObject.GetLoginButtonIconOk(); ok {
			values["login_button_icon"] = flattenLoginButtonIcon(v)
		}

		if v, ok := idpObject.GetIconOk(); ok {
			values["icon"] = flattenIdPIcon(v)
		}

		values[schemaAttribute] = flattenFacebook(idpObject.GetAppId(), idpObject.GetAppSecret())

	case *management.IdentityProviderOIDC:
		idpObject := respObject.IdentityProviderOIDC
		schemaAttribute := "openid_connect"

		values["name"] = idpObject.GetName()

		values["enabled"] = idpObject.GetEnabled()

		if v, ok := idpObject.GetDescriptionOk(); ok {
			values["description"] = v
		}

		if v, ok := idpObject.GetRegistrationOk(); ok {
			values["registration_population_id"] = v.GetPopulation().Id
		}

		if v, ok := idpObject.GetLoginButtonIconOk(); ok {
			values["login_button_icon"] = flattenLoginButtonIcon(v)
		}

		if v, ok := idpObject.GetIconOk(); ok {
			values["icon"] = flattenIdPIcon(v)
		}

		values[schemaAttribute] = flattenOIDC(idpObject)

	case *management.IdentityProviderPaypal:
		idpObject := respObject.IdentityProviderPaypal
		schemaAttribute := "paypal"

		values["name"] = idpObject.GetName()

		values["enabled"] = idpObject.GetEnabled()

		if v, ok := idpObject.GetDescriptionOk(); ok {
			values["description"] = v
		}

		if v, ok := idpObject.GetRegistrationOk(); ok {
			values["registration_population_id"] = v.GetPopulation().Id
		}

		if v, ok := idpObject.GetLoginButtonIconOk(); ok {
			values["login_button_icon"] = flattenLoginButtonIcon(v)
		}

		if v, ok := idpObject.GetIconOk(); ok {
			values["icon"] = flattenIdPIcon(v)
		}

		values[schemaAttribute] = flattenPaypal(idpObject.GetClientId(), idpObject.GetClientSecret(), idpObject.GetClientEnvironment())

	case *management.IdentityProviderSAML:
		idpObject := respObject.IdentityProviderSAML
		schemaAttribute := "saml"

		values["name"] = idpObject.GetName()

		values["enabled"] = idpObject.GetEnabled()

		if v, ok := idpObject.GetDescriptionOk(); ok {
			values["description"] = v
		}

		if v, ok := idpObject.GetRegistrationOk(); ok {
			values["registration_population_id"] = v.GetPopulation().Id
		}

		if v, ok := idpObject.GetLoginButtonIconOk(); ok {
			values["login_button_icon"] = flattenLoginButtonIcon(v)
		}

		if v, ok := idpObject.GetIconOk(); ok {
			values["icon"] = flattenIdPIcon(v)
		}

		values[schemaAttribute] = flattenSAML(idpObject)
	}

	d.Set("name", values["name"])
	d.Set("description", values["description"])
	d.Set("enabled", values["enabled"])

	d.Set("registration_population_id", values["registration_population_id"])

	d.Set("login_button_icon", values["login_button_icon"])
	d.Set("icon", values["icon"])

	d.Set("facebook", values["facebook"])
	d.Set("google", values["google"])
	d.Set("linkedin", values["linkedin"])
	d.Set("yahoo", values["yahoo"])
	d.Set("amazon", values["amazon"])
	d.Set("twitter", values["twitter"])
	d.Set("apple", values["apple"])
	d.Set("paypal", values["paypal"])
	d.Set("microsoft", values["microsoft"])
	d.Set("github", values["github"])
	d.Set("openid_connect", values["openid_connect"])
	d.Set("saml", values["saml"])

	return diags
}

func resourceIdentityProviderUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient

	var diags diag.Diagnostics

	idpRequest, diags := expandIdentityProvider(d)
	if diags.HasError() {
		return diags
	}

	_, diags = sdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return apiClient.IdentityProvidersApi.UpdateIdentityProvider(ctx, d.Get("environment_id").(string), d.Id()).IdentityProvider(*idpRequest).Execute()
		},
		"UpdateIdentityProvider",
		sdk.DefaultCustomError,
		nil,
	)
	if diags.HasError() {
		return diags
	}

	return resourceIdentityProviderRead(ctx, d, meta)
}

func resourceIdentityProviderDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient

	var diags diag.Diagnostics

	_, diags = sdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			r, err := apiClient.IdentityProvidersApi.DeleteIdentityProvider(ctx, d.Get("environment_id").(string), d.Id()).Execute()
			return nil, r, err
		},
		"DeleteIdentityProvider",
		sdk.CustomErrorResourceNotFoundWarning,
		nil,
	)
	if diags.HasError() {
		return diags
	}

	return diags
}

func resourceIdentityProviderImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	splitLength := 2
	attributes := strings.SplitN(d.Id(), "/", splitLength)

	if len(attributes) != splitLength {
		return nil, fmt.Errorf("invalid id (\"%s\") specified, should be in format \"environmentID/identityProviderID\"", d.Id())
	}

	environmentID, identityProviderID := attributes[0], attributes[1]

	d.Set("environment_id", environmentID)
	d.SetId(identityProviderID)

	resourceIdentityProviderRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}

func expandIdentityProvider(d *schema.ResourceData) (*management.IdentityProvider, diag.Diagnostics) {
	var diags diag.Diagnostics

	common := *management.NewIdentityProviderCommon(d.Get("enabled").(bool), d.Get("name").(string), management.ENUMIDENTITYPROVIDEREXT_OPENID_CONNECT)

	if v, ok := d.GetOk("description"); ok {
		common.SetDescription(v.(string))
	}

	if v, ok := d.GetOk("registration_population_id"); ok {
		registrationPopulation := *management.NewIdentityProviderCommonRegistrationPopulation()
		registrationPopulation.SetId(v.(string))
		registration := *management.NewIdentityProviderCommonRegistration()
		registration.SetPopulation(registrationPopulation)
		common.SetRegistration(registration)
	}

	if v, ok := d.GetOk("login_button_icon"); ok {
		if j, okJ := v.([]interface{}); okJ && j != nil && len(j) > 0 {
			attrs := j[0].(map[string]interface{})
			icon := *management.NewIdentityProviderCommonLoginButtonIcon()
			icon.SetId(attrs["id"].(string))
			icon.SetHref(attrs["href"].(string))
			common.SetLoginButtonIcon(icon)
		}
	}

	if v, ok := d.GetOk("icon"); ok {
		if j, okJ := v.([]interface{}); okJ && j != nil && len(j) > 0 {
			attrs := j[0].(map[string]interface{})
			icon := *management.NewIdentityProviderCommonIcon()
			icon.SetId(attrs["id"].(string))
			icon.SetHref(attrs["href"].(string))
			common.SetIcon(icon)
		}
	}

	requestObject := &management.IdentityProvider{}

	processedCount := 0

	if v, ok := d.GetOk("facebook"); ok {
		requestObject.IdentityProviderFacebook, diags = expandIdPFacebook(v.([]interface{}), common)
		processedCount += 1
	}

	if v, ok := d.GetOk("google"); ok {
		requestObject.IdentityProviderClientIDClientSecret, diags = expandIdPGoogle(v.([]interface{}), common)
		processedCount += 1
	}

	if v, ok := d.GetOk("linkedin"); ok {
		requestObject.IdentityProviderClientIDClientSecret, diags = expandIdPLinkedIn(v.([]interface{}), common)
		processedCount += 1
	}

	if v, ok := d.GetOk("yahoo"); ok {
		requestObject.IdentityProviderClientIDClientSecret, diags = expandIdPYahoo(v.([]interface{}), common)
		processedCount += 1
	}

	if v, ok := d.GetOk("amazon"); ok {
		requestObject.IdentityProviderClientIDClientSecret, diags = expandIdPAmazon(v.([]interface{}), common)
		processedCount += 1
	}

	if v, ok := d.GetOk("twitter"); ok {
		requestObject.IdentityProviderClientIDClientSecret, diags = expandIdPTwitter(v.([]interface{}), common)
		processedCount += 1
	}

	if v, ok := d.GetOk("apple"); ok {
		requestObject.IdentityProviderApple, diags = expandIdPApple(v.([]interface{}), common)
		processedCount += 1
	}

	if v, ok := d.GetOk("paypal"); ok {
		requestObject.IdentityProviderPaypal, diags = expandIdPPaypal(v.([]interface{}), common)
		processedCount += 1
	}

	if v, ok := d.GetOk("microsoft"); ok {
		requestObject.IdentityProviderClientIDClientSecret, diags = expandIdPMicrosoft(v.([]interface{}), common)
		processedCount += 1
	}

	if v, ok := d.GetOk("github"); ok {
		requestObject.IdentityProviderClientIDClientSecret, diags = expandIdPGithub(v.([]interface{}), common)
		processedCount += 1
	}

	if v, ok := d.GetOk("openid_connect"); ok {
		requestObject.IdentityProviderOIDC, diags = expandIdPOIDC(v.([]interface{}), common)
		processedCount += 1
	}

	if v, ok := d.GetOk("saml"); ok {
		requestObject.IdentityProviderSAML, diags = expandIdPSAML(v.([]interface{}), common)
		processedCount += 1
	}

	if processedCount > 1 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "More than one identity provider type configured.  This is not supported.",
		})
		return nil, diags
	} else if processedCount == 0 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "No identity provider types configured.  This is not supported.",
		})
		return nil, diags
	}

	return requestObject, diags

}

func expandIdPFacebook(v []interface{}, common management.IdentityProviderCommon) (*management.IdentityProviderFacebook, diag.Diagnostics) {
	var diags diag.Diagnostics

	if len(v) > 0 && v[0] != nil {

		idp := v[0].(map[string]interface{})

		appId := idp["app_id"].(string)
		appSecret := idp["app_secret"].(string)

		idpObj := management.NewIdentityProviderFacebook(common.GetEnabled(), common.GetName(), management.ENUMIDENTITYPROVIDEREXT_FACEBOOK, appId, appSecret)

		if v, ok := common.GetDescriptionOk(); ok {
			idpObj.SetDescription(*v)
		}

		if v, ok := common.GetRegistrationOk(); ok {
			idpObj.SetRegistration(*v)
		}

		if v, ok := common.GetLoginButtonIconOk(); ok {
			idpObj.SetLoginButtonIcon(*v)
		}

		if v, ok := common.GetIconOk(); ok {
			idpObj.SetIcon(*v)
		}

		return idpObj, diags
	}
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Error,
		Summary:  "Block `facebook` must be defined when using the facebook identity provider type",
	})

	return nil, diags

}

func expandIdPClientIdClientSecret(v []interface{}, common management.IdentityProviderCommon, idpType management.EnumIdentityProviderExt) (*management.IdentityProviderClientIDClientSecret, diag.Diagnostics) {
	var diags diag.Diagnostics

	if len(v) > 0 && v[0] != nil {

		idp := v[0].(map[string]interface{})

		clientId := idp["client_id"].(string)
		clientSecret := idp["client_secret"].(string)

		idpObj := management.NewIdentityProviderClientIDClientSecret(common.GetEnabled(), common.GetName(), idpType, clientId, clientSecret)

		if v, ok := common.GetDescriptionOk(); ok {
			idpObj.SetDescription(*v)
		}

		if v, ok := common.GetRegistrationOk(); ok {
			idpObj.SetRegistration(*v)
		}

		if v, ok := common.GetLoginButtonIconOk(); ok {
			idpObj.SetLoginButtonIcon(*v)
		}

		if v, ok := common.GetIconOk(); ok {
			idpObj.SetIcon(*v)
		}

		return idpObj, diags
	}
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Error,
		Summary:  "Identity provider block not defined correctly.  Please raise an issue.",
	})

	return nil, diags

}

func expandIdPGoogle(v []interface{}, common management.IdentityProviderCommon) (*management.IdentityProviderClientIDClientSecret, diag.Diagnostics) {
	return expandIdPClientIdClientSecret(v, common, management.ENUMIDENTITYPROVIDEREXT_GOOGLE)
}

func expandIdPLinkedIn(v []interface{}, common management.IdentityProviderCommon) (*management.IdentityProviderClientIDClientSecret, diag.Diagnostics) {
	return expandIdPClientIdClientSecret(v, common, management.ENUMIDENTITYPROVIDEREXT_LINKEDIN)
}

func expandIdPYahoo(v []interface{}, common management.IdentityProviderCommon) (*management.IdentityProviderClientIDClientSecret, diag.Diagnostics) {
	return expandIdPClientIdClientSecret(v, common, management.ENUMIDENTITYPROVIDEREXT_YAHOO)
}

func expandIdPAmazon(v []interface{}, common management.IdentityProviderCommon) (*management.IdentityProviderClientIDClientSecret, diag.Diagnostics) {
	return expandIdPClientIdClientSecret(v, common, management.ENUMIDENTITYPROVIDEREXT_AMAZON)
}

func expandIdPTwitter(v []interface{}, common management.IdentityProviderCommon) (*management.IdentityProviderClientIDClientSecret, diag.Diagnostics) {
	return expandIdPClientIdClientSecret(v, common, management.ENUMIDENTITYPROVIDEREXT_TWITTER)
}

func expandIdPApple(v []interface{}, common management.IdentityProviderCommon) (*management.IdentityProviderApple, diag.Diagnostics) {
	var diags diag.Diagnostics

	if len(v) > 0 && v[0] != nil {

		idp := v[0].(map[string]interface{})

		clientId := idp["client_id"].(string)
		clientSecretSigningKey := idp["client_secret_signing_key"].(string)
		keyId := idp["key_id"].(string)
		teamId := idp["team_id"].(string)

		idpObj := management.NewIdentityProviderApple(common.GetEnabled(), common.GetName(), management.ENUMIDENTITYPROVIDEREXT_APPLE, clientId, clientSecretSigningKey, keyId, teamId)

		if v, ok := common.GetDescriptionOk(); ok {
			idpObj.SetDescription(*v)
		}

		if v, ok := common.GetRegistrationOk(); ok {
			idpObj.SetRegistration(*v)
		}

		if v, ok := common.GetLoginButtonIconOk(); ok {
			idpObj.SetLoginButtonIcon(*v)
		}

		if v, ok := common.GetIconOk(); ok {
			idpObj.SetIcon(*v)
		}

		return idpObj, diags
	}
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Error,
		Summary:  "Block `apple` must be defined when using the Apple identity provider type",
	})

	return nil, diags

}

func expandIdPPaypal(v []interface{}, common management.IdentityProviderCommon) (*management.IdentityProviderPaypal, diag.Diagnostics) {
	var diags diag.Diagnostics

	if len(v) > 0 && v[0] != nil {

		idp := v[0].(map[string]interface{})

		clientId := idp["client_id"].(string)
		clientSecret := idp["client_secret"].(string)
		clientEnvironment := idp["client_environment"].(string)

		idpObj := management.NewIdentityProviderPaypal(common.GetEnabled(), common.GetName(), management.ENUMIDENTITYPROVIDEREXT_PAYPAL, clientId, clientSecret, clientEnvironment)

		if v, ok := common.GetDescriptionOk(); ok {
			idpObj.SetDescription(*v)
		}

		if v, ok := common.GetRegistrationOk(); ok {
			idpObj.SetRegistration(*v)
		}

		if v, ok := common.GetLoginButtonIconOk(); ok {
			idpObj.SetLoginButtonIcon(*v)
		}

		if v, ok := common.GetIconOk(); ok {
			idpObj.SetIcon(*v)
		}

		return idpObj, diags
	}
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Error,
		Summary:  "Block `paypal` must be defined when using the Paypal identity provider type",
	})

	return nil, diags

}

func expandIdPMicrosoft(v []interface{}, common management.IdentityProviderCommon) (*management.IdentityProviderClientIDClientSecret, diag.Diagnostics) {
	return expandIdPClientIdClientSecret(v, common, management.ENUMIDENTITYPROVIDEREXT_MICROSOFT)
}

func expandIdPGithub(v []interface{}, common management.IdentityProviderCommon) (*management.IdentityProviderClientIDClientSecret, diag.Diagnostics) {
	return expandIdPClientIdClientSecret(v, common, management.ENUMIDENTITYPROVIDEREXT_GITHUB)
}

func expandIdPOIDC(v []interface{}, common management.IdentityProviderCommon) (*management.IdentityProviderOIDC, diag.Diagnostics) {
	var diags diag.Diagnostics

	if len(v) > 0 && v[0] != nil {

		idp := v[0].(map[string]interface{})

		authorizationEndpoint := idp["authorization_endpoint"].(string)
		clientId := idp["client_id"].(string)
		clientSecret := idp["client_secret"].(string)
		issuer := idp["issuer"].(string)
		jwksEndpoint := idp["jwks_endpoint"].(string)

		scopes := make([]string, 0)
		for _, scope := range idp["scopes"].(*schema.Set).List() {
			scopes = append(scopes, scope.(string))
		}

		tokenEndpoint := idp["token_endpoint"].(string)
		tokenEndpointAuthMethod := management.EnumIdentityProviderOIDCTokenAuthMethod(idp["token_endpoint_auth_method"].(string))

		idpObj := management.NewIdentityProviderOIDC(common.GetEnabled(), common.GetName(), management.ENUMIDENTITYPROVIDEREXT_OPENID_CONNECT, authorizationEndpoint, clientId, clientSecret, issuer, jwksEndpoint, scopes, tokenEndpoint, tokenEndpointAuthMethod)

		if v, ok := idp["discovery_endpoint"].(string); ok && v != "" {
			idpObj.SetDiscoveryEndpoint(v)
		}

		if v, ok := idp["userinfo_endpoint"].(string); ok && v != "" {
			idpObj.SetUserInfoEndpoint(v)
		}

		if v, ok := common.GetDescriptionOk(); ok {
			idpObj.SetDescription(*v)
		}

		if v, ok := common.GetRegistrationOk(); ok {
			idpObj.SetRegistration(*v)
		}

		if v, ok := common.GetLoginButtonIconOk(); ok {
			idpObj.SetLoginButtonIcon(*v)
		}

		if v, ok := common.GetIconOk(); ok {
			idpObj.SetIcon(*v)
		}

		return idpObj, diags
	}
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Error,
		Summary:  "Block `openid_connect` must be defined when using the OpenID Connect identity provider type",
	})

	return nil, diags

}

func expandIdPSAML(v []interface{}, common management.IdentityProviderCommon) (*management.IdentityProviderSAML, diag.Diagnostics) {
	var diags diag.Diagnostics

	if len(v) > 0 && v[0] != nil {

		idp := v[0].(map[string]interface{})

		var idpVerification management.IdentityProviderSAMLAllOfIdpVerification

		if v, ok := idp["idp_verification_certificate_ids"].(*schema.Set); ok {

			planCertificates := v.List()

			if len(planCertificates) > 0 && planCertificates[0] != nil {
				certificates := make([]management.IdentityProviderSAMLAllOfIdpVerificationCertificates, 0)

				for _, certificate := range planCertificates {
					certificates = append(certificates, *management.NewIdentityProviderSAMLAllOfIdpVerificationCertificates(certificate.(string)))
				}

				idpVerification = *management.NewIdentityProviderSAMLAllOfIdpVerification(certificates)
			}
		}

		idpObj := management.NewIdentityProviderSAML(common.GetEnabled(), common.GetName(), management.ENUMIDENTITYPROVIDEREXT_SAML, idp["idp_entity_id"].(string), idpVerification, idp["sp_entity_id"].(string), management.EnumIdentityProviderSAMLSSOBinding(idp["sso_binding"].(string)), idp["sso_endpoint"].(string))

		if v, ok := idp["authentication_request_signed"].(bool); ok {
			idpObj.SetAuthnRequestSigned(v)
		}

		if v, ok := idp["sp_signing_key_id"].(string); ok && v != "" {
			idpObj.SetSpSigning(*management.NewIdentityProviderSAMLAllOfSpSigning(*management.NewIdentityProviderSAMLAllOfSpSigningKey(v)))
		}

		if v1, ok := idp["slo_binding"].(string); ok && v1 != "" {
			idpObj.SetSloBinding(management.EnumIdentityProviderSAMLSLOBinding(v1))
		}

		if v1, ok := idp["slo_endpoint"].(string); ok && v1 != "" {
			idpObj.SetSloEndpoint(v1)
		}

		if v1, ok := idp["slo_response_endpoint"].(string); ok && v1 != "" {
			idpObj.SetSloResponseEndpoint(v1)
		}

		if v1, ok := idp["slo_window"].(int); ok && v1 > 0 {
			idpObj.SetSloWindow(int32(v1))
		}

		if v, ok := common.GetDescriptionOk(); ok {
			idpObj.SetDescription(*v)
		}

		if v, ok := common.GetRegistrationOk(); ok {
			idpObj.SetRegistration(*v)
		}

		if v, ok := common.GetLoginButtonIconOk(); ok {
			idpObj.SetLoginButtonIcon(*v)
		}

		if v, ok := common.GetIconOk(); ok {
			idpObj.SetIcon(*v)
		}

		return idpObj, diags
	}
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Error,
		Summary:  "Block `saml` must be defined when using the SAML identity provider type",
	})

	return nil, diags

}

func flattenLoginButtonIcon(s *management.IdentityProviderCommonLoginButtonIcon) []interface{} {

	item := map[string]interface{}{
		"id":   s.GetId(),
		"href": s.GetHref(),
	}

	items := make([]interface{}, 0)
	return append(items, item)
}

func flattenIdPIcon(s *management.IdentityProviderCommonIcon) []interface{} {

	item := map[string]interface{}{
		"id":   s.GetId(),
		"href": s.GetHref(),
	}

	items := make([]interface{}, 0)
	return append(items, item)
}

func flattenClientIdClientSecret(clientId, clientSecret string) []interface{} {

	item := map[string]interface{}{
		"client_id":     clientId,
		"client_secret": clientSecret,
	}

	items := make([]interface{}, 0)
	return append(items, item)
}

func flattenSAML(idpObject *management.IdentityProviderSAML) []interface{} {

	item := map[string]interface{}{
		"idp_entity_id": idpObject.GetIdpEntityId(),
		"sso_binding":   string(idpObject.GetSsoBinding()),
		"sso_endpoint":  idpObject.GetSsoEndpoint(),
	}

	if v, ok := idpObject.GetAuthnRequestSignedOk(); ok {
		item["authentication_request_signed"] = v
	} else {
		item["authentication_request_signed"] = nil
	}

	if v, ok := idpObject.GetSpEntityIdOk(); ok {
		item["sp_entity_id"] = *v
	} else {
		item["sp_entity_id"] = nil
	}

	if v, ok := idpObject.GetIdpVerificationOk(); ok {
		if v1, ok := v.GetCertificatesOk(); ok {
			ids := make([]string, 0)
			for _, certificate := range v1 {
				ids = append(ids, certificate.GetId())
			}
			item["idp_verification_certificate_ids"] = ids
		} else {
			item["idp_verification_certificate_ids"] = nil
		}
	} else {
		item["idp_verification_certificate_ids"] = nil
	}

	if v, ok := idpObject.GetSpSigningOk(); ok {
		item["sp_signing_key_id"] = v.GetKey().Id
	} else {
		item["sp_signing_key_id"] = nil
	}

	if v, ok := idpObject.GetSloBindingOk(); ok {
		item["slo_binding"] = v
	} else {
		item["slo_binding"] = nil
	}

	if v, ok := idpObject.GetSloEndpointOk(); ok {
		item["slo_endpoint"] = v
	} else {
		item["slo_endpoint"] = nil
	}

	if v, ok := idpObject.GetSloResponseEndpointOk(); ok {
		item["slo_response_endpoint"] = v
	} else {
		item["slo_response_endpoint"] = nil
	}

	if v, ok := idpObject.GetSloWindowOk(); ok {
		item["slo_window"] = v
	} else {
		item["slo_window"] = nil
	}

	items := make([]interface{}, 0)
	return append(items, item)
}

func flattenOIDC(idpObject *management.IdentityProviderOIDC) []interface{} {

	item := map[string]interface{}{
		"authorization_endpoint": idpObject.GetAuthorizationEndpoint(),
		"client_id":              idpObject.GetClientId(),
		"client_secret":          idpObject.GetClientSecret(),
		"issuer":                 idpObject.GetIssuer(),
		"jwks_endpoint":          idpObject.GetJwksEndpoint(),
		"scopes":                 idpObject.GetScopes(),
		"token_endpoint":         idpObject.GetTokenEndpoint(),
	}

	if v, ok := idpObject.GetDiscoveryEndpointOk(); ok {
		item["discovery_endpoint"] = v
	} else {
		item["discovery_endpoint"] = nil
	}

	if v, ok := idpObject.GetTokenEndpointAuthMethodOk(); ok {
		item["token_endpoint_auth_method"] = string(*v)
	} else {
		item["token_endpoint_auth_method"] = nil
	}

	if v, ok := idpObject.GetUserInfoEndpointOk(); ok {
		item["userinfo_endpoint"] = v
	} else {
		item["userinfo_endpoint"] = nil
	}

	items := make([]interface{}, 0)
	return append(items, item)
}

func flattenPaypal(clientId, clientSecret, clientEnvironment string) []interface{} {

	item := map[string]interface{}{
		"client_id":          clientId,
		"client_secret":      clientSecret,
		"client_environment": clientEnvironment,
	}

	items := make([]interface{}, 0)
	return append(items, item)
}

func flattenFacebook(appId, appSecret string) []interface{} {

	item := map[string]interface{}{
		"app_id":     appId,
		"app_secret": appSecret,
	}

	items := make([]interface{}, 0)
	return append(items, item)
}

func flattenApple(clientId, clientSecretSigningKey, keyId, teamId string) []interface{} {

	item := map[string]interface{}{
		"client_id":                 clientId,
		"client_secret_signing_key": clientSecretSigningKey,
		"key_id":                    keyId,
		"team_id":                   teamId,
	}

	items := make([]interface{}, 0)
	return append(items, item)
}
