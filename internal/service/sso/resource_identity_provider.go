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

	providerAttributeList := []string{"facebook", "google", "linkedin", "yahoo", "amazon", "twitter", "apple", "paypal", "microsoft", "github", "generic_oidc", "saml_options"}

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
				ExactlyOneOf: providerAttributeList,
				Elem:         clientIdClientSecretSchema("Google"),
			},
			"linkedin": {
				Description:  "Options for Identity provider connectivity to LinkedIn.",
				Type:         schema.TypeList,
				MaxItems:     1,
				Optional:     true,
				ExactlyOneOf: providerAttributeList,
				Elem:         clientIdClientSecretSchema("LinkedIn"),
			},
			"yahoo": {
				Description:  "Options for Identity provider connectivity to Yahoo.",
				Type:         schema.TypeList,
				MaxItems:     1,
				Optional:     true,
				ExactlyOneOf: providerAttributeList,
				Elem:         clientIdClientSecretSchema("Yahoo"),
			},
			"amazon": {
				Description:  "Options for Identity provider connectivity to Amazon.",
				Type:         schema.TypeList,
				MaxItems:     1,
				Optional:     true,
				ExactlyOneOf: providerAttributeList,
				Elem:         clientIdClientSecretSchema("Amazon"),
			},
			"twitter": {
				Description:  "Options for Identity provider connectivity to Twitter.",
				Type:         schema.TypeList,
				MaxItems:     1,
				Optional:     true,
				ExactlyOneOf: providerAttributeList,
				Elem:         clientIdClientSecretSchema("Twitter"),
			},
			"apple": {
				Description:  "Options for Identity provider connectivity to Apple.",
				Type:         schema.TypeList,
				MaxItems:     1,
				Optional:     true,
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
				ExactlyOneOf: providerAttributeList,
				Elem:         clientIdClientSecretSchema("Microsoft"),
			},
			"github": {
				Description:  "Options for Identity provider connectivity to Github.",
				Type:         schema.TypeList,
				MaxItems:     1,
				Optional:     true,
				ExactlyOneOf: providerAttributeList,
				Elem:         clientIdClientSecretSchema("Github"),
			},
			"generic_oidc": {
				Description:  "Options for Identity provider connectivity to a generic OpenID Connect service.",
				Type:         schema.TypeList,
				MaxItems:     1,
				Optional:     true,
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
			"generic_saml": {
				Description:  "Options for Identity provider connectivity to a generic SAML service.",
				Type:         schema.TypeList,
				MaxItems:     1,
				Optional:     true,
				ExactlyOneOf: providerAttributeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"authentication_request_signed": {
							Description: "A boolean that specifies whether the SAML authentication request will be signed when sending to the identity provider. Set this to true if the external IDP is included in an authentication policy to be used by applications that are accessed using a mix of default URLS and custom Domains URLs.",
							Type:        schema.TypeBool,
							Optional:    true,
						},
						"idp_entity_id": {
							Description: "A string that specifies the entity ID URI that is checked against the issuerId tag in the incoming response.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"sp_entity_id": {
							Description: "A string that specifies the service provider's entity ID, used to look up the application.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"idp_verification_certificate_ids": {
							Description: "A list that specifies the identity provider's certificate IDs used to verify the signature on the signed assertion from the identity provider. Signing is done with a private key and verified with a public key.",
							Type:        schema.TypeSet,
							Optional:    true,
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
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{string(management.ENUMIDENTITYPROVIDERSAMLSSOBINDING_POST), string(management.ENUMIDENTITYPROVIDERSAMLSSOBINDING_REDIRECT)}, false)),
						},
						"sso_endpoint": {
							Description:      "A string that specifies the SSO endpoint for the authentication request.",
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IsURLWithHTTPorHTTPS),
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
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
			},
		},
	}
}

func resourceIdentityProviderCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	idpRequest, diags := expandIdentityProvider(d)
	if diags.HasError() {
		return diags
	}

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.IdentityProviderManagementIdentityProvidersApi.CreateIdentityProvider(ctx, d.Get("environment_id").(string)).IdentityProvider(*idpRequest).Execute()
		},
		"CreateIdentityProvider",
		sdk.DefaultCustomError,
		sdk.DefaultCreateReadRetryable,
	)
	if diags.HasError() {
		return diags
	}

	respObject := resp.(*management.IdentityProvider)

	if respObject.IdentityProviderOIDC != nil && respObject.IdentityProviderOIDC.GetId() != "" {
		d.SetId(respObject.IdentityProviderOIDC.GetId())
	} else if respObject.IdentityProviderSAML != nil && respObject.IdentityProviderSAML.GetId() != "" {
		d.SetId(respObject.IdentityProviderSAML.GetId())
	} else {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Cannot determine application ID from API response for application: %s", d.Get("name")),
			Detail:   fmt.Sprintf("Full response object: %v\n", resp),
		})

		return diags
	}

	return resourceIdentityProviderRead(ctx, d, meta)
}

func resourceIdentityProviderRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	resp, diags := sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.IdentityProviderManagementIdentityProvidersApi.ReadOneIdentityProvider(ctx, d.Get("environment_id").(string), d.Id()).Execute()
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

	return diags
}

func resourceIdentityProviderUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	idpRequest, diags := expandIdentityProvider(d)
	if diags.HasError() {
		return diags
	}

	_, diags = sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			return apiClient.IdentityProviderManagementIdentityProvidersApi.UpdateIdentityProvider(ctx, d.Get("environment_id").(string), d.Id()).IdentityProvider(*idpRequest).Execute()
		},
		"UpdateIdentityProvider",
		sdk.DefaultCustomError,
		sdk.DefaultRetryable,
	)
	if diags.HasError() {
		return diags
	}

	return resourceIdentityProviderRead(ctx, d, meta)
}

func resourceIdentityProviderDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p1Client := meta.(*client.Client)
	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})
	var diags diag.Diagnostics

	_, diags = sdk.ParseResponse(
		ctx,

		func() (interface{}, *http.Response, error) {
			r, err := apiClient.IdentityProviderManagementIdentityProvidersApi.DeleteIdentityProvider(ctx, d.Get("environment_id").(string), d.Id()).Execute()
			return nil, r, err
		},
		"DeleteIdentityProvider",
		sdk.CustomErrorResourceNotFoundWarning,
		sdk.DefaultRetryable,
	)
	if diags.HasError() {
		return diags
	}

	return diags
}

func resourceIdentityProviderImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	attributes := strings.SplitN(d.Id(), "/", 2)

	if len(attributes) != 2 {
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

	enabled := management.ENUMENABLEDSTATUS_DISABLED

	if d.Get("enabled").(bool) {
		enabled = management.ENUMENABLEDSTATUS_ENABLED
	}

	common := *management.NewIdentityProviderCommon(enabled, d.Get("name").(string), management.ENUMIDENTITYPROVIDEREXT_OPENID_CONNECT)

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
		requestObject.IdentityProviderGoogle, diags = expandIdPGoogle(v.([]interface{}), common)
		processedCount += 1
	}

	if v, ok := d.GetOk("linkedin"); ok {
		requestObject.IdentityProviderLinkedIn, diags = expandIdPLinkedIn(v.([]interface{}), common)
		processedCount += 1
	}

	if v, ok := d.GetOk("yahoo"); ok {
		requestObject.IdentityProviderYahoo, diags = expandIdPYahoo(v.([]interface{}), common)
		processedCount += 1
	}

	if v, ok := d.GetOk("amazon"); ok {
		requestObject.IdentityProviderAmazon, diags = expandIdPAmazon(v.([]interface{}), common)
		processedCount += 1
	}

	if v, ok := d.GetOk("twitter"); ok {
		requestObject.IdentityProviderTwitter, diags = expandIdPTwitter(v.([]interface{}), common)
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
		requestObject.IdentityProviderMicrosoft, diags = expandIdPMicrosoft(v.([]interface{}), common)
		processedCount += 1
	}

	if v, ok := d.GetOk("github"); ok {
		requestObject.IdentityProviderGithub, diags = expandIdPGithub(v.([]interface{}), common)
		processedCount += 1
	}

	if v, ok := d.GetOk("generic_oidc"); ok {
		requestObject.IdentityProviderOIDC, diags = expandIdPOIDC(v.([]interface{}), common)
		processedCount += 1
	}

	if v, ok := d.GetOk("generic_saml"); ok {
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

	if v != nil && len(v) > 0 && v[0] != nil {

		idp := v[0].(map[string]interface{})

		appId := idp["app_id"].(string)
		appSecret := idp["app_secret"].(string)

		idpObj := management.NewIdentityProviderFacebook(common.GetEnabled(), common.GetName(), management.ENUMIDENTITYPROVIDEREXT_FACEBOOK, appId, appSecret)

		idpObj.SetDescription(common.GetDescription())
		idpObj.SetRegistration(common.GetRegistration())
		idpObj.SetLoginButtonIcon(common.GetLoginButtonIcon())
		idpObj.SetIcon(common.GetIcon())

		return idpObj, diags
	}
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Error,
		Summary:  "Block `facebook` must be defined when using the facebook identity provider type",
	})

	return nil, diags

}

func expandIdPGoogle(v []interface{}, common management.IdentityProviderCommon) (*management.IdentityProviderGoogle, diag.Diagnostics) {
	var diags diag.Diagnostics

	if v != nil && len(v) > 0 && v[0] != nil {

		idp := v[0].(map[string]interface{})

		clientId := idp["client_id"].(string)
		clientSecret := idp["client_secret"].(string)

		idpObj := management.NewIdentityProviderGoogle(common.GetEnabled(), common.GetName(), management.ENUMIDENTITYPROVIDEREXT_GOOGLE, clientId, clientSecret)

		idpObj.SetDescription(common.GetDescription())
		idpObj.SetRegistration(common.GetRegistration())
		idpObj.SetLoginButtonIcon(common.GetLoginButtonIcon())
		idpObj.SetIcon(common.GetIcon())

		return idpObj, diags
	}
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Error,
		Summary:  "Block `google` must be defined when using the google identity provider type",
	})

	return nil, diags

}

func expandIdPLinkedIn(v []interface{}, common management.IdentityProviderCommon) (*management.IdentityProviderLinkedIn, diag.Diagnostics) {
	var diags diag.Diagnostics

	if v != nil && len(v) > 0 && v[0] != nil {

		idp := v[0].(map[string]interface{})

		clientId := idp["client_id"].(string)
		clientSecret := idp["client_secret"].(string)

		idpObj := management.NewIdentityProviderLinkedIn(common.GetEnabled(), common.GetName(), management.ENUMIDENTITYPROVIDEREXT_LINKEDIN, clientId, clientSecret)

		idpObj.SetDescription(common.GetDescription())
		idpObj.SetRegistration(common.GetRegistration())
		idpObj.SetLoginButtonIcon(common.GetLoginButtonIcon())
		idpObj.SetIcon(common.GetIcon())

		return idpObj, diags
	}
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Error,
		Summary:  "Block `linkedin` must be defined when using the LinkedIn identity provider type",
	})

	return nil, diags

}

func expandIdPYahoo(v []interface{}, common management.IdentityProviderCommon) (*management.IdentityProviderYahoo, diag.Diagnostics) {
	var diags diag.Diagnostics

	if v != nil && len(v) > 0 && v[0] != nil {

		idp := v[0].(map[string]interface{})

		clientId := idp["client_id"].(string)
		clientSecret := idp["client_secret"].(string)

		idpObj := management.NewIdentityProviderYahoo(common.GetEnabled(), common.GetName(), management.ENUMIDENTITYPROVIDEREXT_YAHOO, clientId, clientSecret)

		idpObj.SetDescription(common.GetDescription())
		idpObj.SetRegistration(common.GetRegistration())
		idpObj.SetLoginButtonIcon(common.GetLoginButtonIcon())
		idpObj.SetIcon(common.GetIcon())

		return idpObj, diags
	}
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Error,
		Summary:  "Block `yahoo` must be defined when using the Yahoo identity provider type",
	})

	return nil, diags

}

func expandIdPAmazon(v []interface{}, common management.IdentityProviderCommon) (*management.IdentityProviderAmazon, diag.Diagnostics) {
	var diags diag.Diagnostics

	if v != nil && len(v) > 0 && v[0] != nil {

		idp := v[0].(map[string]interface{})

		clientId := idp["client_id"].(string)
		clientSecret := idp["client_secret"].(string)

		idpObj := management.NewIdentityProviderAmazon(common.GetEnabled(), common.GetName(), management.ENUMIDENTITYPROVIDEREXT_AMAZON, clientId, clientSecret)

		idpObj.SetDescription(common.GetDescription())
		idpObj.SetRegistration(common.GetRegistration())
		idpObj.SetLoginButtonIcon(common.GetLoginButtonIcon())
		idpObj.SetIcon(common.GetIcon())

		return idpObj, diags
	}
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Error,
		Summary:  "Block `amazon` must be defined when using the Amazon identity provider type",
	})

	return nil, diags

}

func expandIdPTwitter(v []interface{}, common management.IdentityProviderCommon) (*management.IdentityProviderTwitter, diag.Diagnostics) {
	var diags diag.Diagnostics

	if v != nil && len(v) > 0 && v[0] != nil {

		idp := v[0].(map[string]interface{})

		clientId := idp["client_id"].(string)
		clientSecret := idp["client_secret"].(string)

		idpObj := management.NewIdentityProviderTwitter(common.GetEnabled(), common.GetName(), management.ENUMIDENTITYPROVIDEREXT_TWITTER, clientId, clientSecret)

		idpObj.SetDescription(common.GetDescription())
		idpObj.SetRegistration(common.GetRegistration())
		idpObj.SetLoginButtonIcon(common.GetLoginButtonIcon())
		idpObj.SetIcon(common.GetIcon())

		return idpObj, diags
	}
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Error,
		Summary:  "Block `twitter` must be defined when using the Twitter identity provider type",
	})

	return nil, diags

}

func expandIdPApple(v []interface{}, common management.IdentityProviderCommon) (*management.IdentityProviderApple, diag.Diagnostics) {
	var diags diag.Diagnostics

	if v != nil && len(v) > 0 && v[0] != nil {

		idp := v[0].(map[string]interface{})

		clientId := idp["client_id"].(string)
		clientSecretSigningKey := idp["client_secret_signing_key"].(string)
		keyId := idp["key_id"].(string)
		teamId := idp["team_id"].(string)

		idpObj := management.NewIdentityProviderApple(common.GetEnabled(), common.GetName(), management.ENUMIDENTITYPROVIDEREXT_APPLE, clientId, clientSecretSigningKey, keyId, teamId)

		idpObj.SetDescription(common.GetDescription())
		idpObj.SetRegistration(common.GetRegistration())
		idpObj.SetLoginButtonIcon(common.GetLoginButtonIcon())
		idpObj.SetIcon(common.GetIcon())

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

	if v != nil && len(v) > 0 && v[0] != nil {

		idp := v[0].(map[string]interface{})

		clientId := idp["client_id"].(string)
		clientSecret := idp["client_secret"].(string)
		clientEnvironment := idp["client_environment"].(string)

		idpObj := management.NewIdentityProviderPaypal(common.GetEnabled(), common.GetName(), management.ENUMIDENTITYPROVIDEREXT_PAYPAL, clientId, clientSecret, clientEnvironment)

		idpObj.SetDescription(common.GetDescription())
		idpObj.SetRegistration(common.GetRegistration())
		idpObj.SetLoginButtonIcon(common.GetLoginButtonIcon())
		idpObj.SetIcon(common.GetIcon())

		return idpObj, diags
	}
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Error,
		Summary:  "Block `paypal` must be defined when using the Paypal identity provider type",
	})

	return nil, diags

}

func expandIdPMicrosoft(v []interface{}, common management.IdentityProviderCommon) (*management.IdentityProviderMicrosoft, diag.Diagnostics) {
	var diags diag.Diagnostics

	if v != nil && len(v) > 0 && v[0] != nil {

		idp := v[0].(map[string]interface{})

		clientId := idp["client_id"].(string)
		clientSecret := idp["client_secret"].(string)

		idpObj := management.NewIdentityProviderMicrosoft(common.GetEnabled(), common.GetName(), management.ENUMIDENTITYPROVIDEREXT_MICROSOFT, clientId, clientSecret)

		idpObj.SetDescription(common.GetDescription())
		idpObj.SetRegistration(common.GetRegistration())
		idpObj.SetLoginButtonIcon(common.GetLoginButtonIcon())
		idpObj.SetIcon(common.GetIcon())

		return idpObj, diags
	}
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Error,
		Summary:  "Block `microsoft` must be defined when using the Microsoft identity provider type",
	})

	return nil, diags

}

func expandIdPGithub(v []interface{}, common management.IdentityProviderCommon) (*management.IdentityProviderGithub, diag.Diagnostics) {
	var diags diag.Diagnostics

	if v != nil && len(v) > 0 && v[0] != nil {

		idp := v[0].(map[string]interface{})

		clientId := idp["client_id"].(string)
		clientSecret := idp["client_secret"].(string)

		idpObj := management.NewIdentityProviderGithub(common.GetEnabled(), common.GetName(), management.ENUMIDENTITYPROVIDEREXT_GITHUB, clientId, clientSecret)

		idpObj.SetDescription(common.GetDescription())
		idpObj.SetRegistration(common.GetRegistration())
		idpObj.SetLoginButtonIcon(common.GetLoginButtonIcon())
		idpObj.SetIcon(common.GetIcon())

		return idpObj, diags
	}
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Error,
		Summary:  "Block `github` must be defined when using the Github identity provider type",
	})

	return nil, diags

}

func expandIdPOIDC(v []interface{}, common management.IdentityProviderCommon) (*management.IdentityProviderOIDC, diag.Diagnostics) {
	var diags diag.Diagnostics

	if v != nil && len(v) > 0 && v[0] != nil {

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

		idpObj.SetDescription(common.GetDescription())
		idpObj.SetRegistration(common.GetRegistration())
		idpObj.SetLoginButtonIcon(common.GetLoginButtonIcon())
		idpObj.SetIcon(common.GetIcon())

		return idpObj, diags
	}
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Error,
		Summary:  "Block `generic_oidc` must be defined when using the OpenID Connect identity provider type",
	})

	return nil, diags

}

func expandIdPSAML(v []interface{}, common management.IdentityProviderCommon) (*management.IdentityProviderSAML, diag.Diagnostics) {
	var diags diag.Diagnostics

	if v != nil && len(v) > 0 && v[0] != nil {

		idp := v[0].(map[string]interface{})

		idpObj := management.NewIdentityProviderSAML(common.GetEnabled(), common.GetName(), management.ENUMIDENTITYPROVIDEREXT_SAML)

		if v, ok := idp["authentication_request_signed"].(bool); ok {
			idpObj.SetAuthnRequestSigned(v)
		}

		if v, ok := idp["idp_entity_id"].(string); ok {
			idpObj.SetIdpEntityId(v)
		}

		if v, ok := idp["sp_entity_id"].(string); ok {
			idpObj.SetSpEntityId(v)
		}

		if v, ok := idp["idp_verification_certificate_ids"].([]interface{}); ok && v != nil && len(v) > 0 && v[0] != nil {
			certificates := make([]management.ApplicationAccessControlGroupGroupsInner, 0)

			for _, certificate := range v {
				certificates = append(certificates, *management.NewApplicationAccessControlGroupGroupsInner(certificate.(string)))
			}

			verification := *management.NewIdentityProviderSAMLAllOfIdpVerification()
			verification.SetCertificates(certificates)
			idpObj.SetIdpVerification(verification)
		}

		if v, ok := idp["sp_signing_key_id"].(string); ok {
			idpObj.SetSpSigning(*management.NewIdentityProviderSAMLAllOfSpSigning(*management.NewIdentityProviderSAMLAllOfSpSigningKey(v)))
		}

		if v, ok := idp["sso_binding"].(string); ok {
			idpObj.SetSsoBinding(management.EnumIdentityProviderSAMLSSOBinding(v))
		}

		if v, ok := idp["sso_endpoint"].(string); ok {
			idpObj.SetSsoEndpoint(v)
		}

		idpObj.SetDescription(common.GetDescription())
		idpObj.SetRegistration(common.GetRegistration())
		idpObj.SetLoginButtonIcon(common.GetLoginButtonIcon())
		idpObj.SetIcon(common.GetIcon())

		return idpObj, diags
	}
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Error,
		Summary:  "Block `generic_saml` must be defined when using the SAML identity provider type",
	})

	return nil, diags

}
