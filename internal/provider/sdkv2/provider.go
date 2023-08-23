package sdkv2

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/service/authorize"
	"github.com/pingidentity/terraform-provider-pingone/internal/service/base"
	"github.com/pingidentity/terraform-provider-pingone/internal/service/mfa"
	"github.com/pingidentity/terraform-provider-pingone/internal/service/sso"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
		desc := s.Description
		if s.Default != nil {
			desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
		}
		return strings.TrimSpace(desc)
	}
}

// New provider function
func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{

			Schema: map[string]*schema.Schema{
				"client_id": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Client ID for the worker app client.  Default value can be set with the `PINGONE_CLIENT_ID` environment variable.  Must provide only one of `api_access_token` (when obtaining the worker token outside of the provider) and `client_id` (when the provider should fetch the worker token during operations).  Must be configured with `client_secret` and `environment_id`.",
				},
				"client_secret": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Client secret for the worker app client.  Default value can be set with the `PINGONE_CLIENT_SECRET` environment variable.  Must be configured with `client_id` and `environment_id`.",
				},
				"environment_id": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Environment ID for the worker app client.  Default value can be set with the `PINGONE_ENVIRONMENT_ID` environment variable.  Must be configured with `client_id` and `client_secret`.",
				},
				"api_access_token": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The access token used for provider resource management against the PingOne management API.  Default value can be set with the `PINGONE_API_ACCESS_TOKEN` environment variable.  Must provide only one of `api_access_token` (when obtaining the worker token outside of the provider) and `client_id` (when the provider should fetch the worker token during operations).",
				},
				"region": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The PingOne region to use.  Options are `AsiaPacific` `Canada` `Europe` and `NorthAmerica`.  Default value can be set with the `PINGONE_REGION` environment variable.",
				},
				"force_delete_production_type": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "Choose whether to force-delete any configuration that has a `PRODUCTION` type parameter.  The platform default is that `PRODUCTION` type configuration will not destroy without intervention to protect stored data.  By default this parameter is set to `false` and can be overridden with the `PINGONE_FORCE_DELETE_PRODUCTION_TYPE` environment variable.",
				},
				"http_proxy": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Full URL for the http/https proxy service, for example `http://127.0.0.1:8090`.  Default value can be set with the `HTTP_PROXY` or `HTTPS_PROXY` environment variables.",
				},
				"service_endpoints": {
					Type:        schema.TypeList,
					Optional:    true,
					MaxItems:    1,
					Description: "A single block containing configuration items to override the service API endpoints of PingOne.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"auth_hostname": {
								Type:        schema.TypeString,
								Required:    true,
								Description: "Hostname for the PingOne authentication service API.  Default value can be set with the `PINGONE_AUTH_SERVICE_HOSTNAME` environment variable.",
							},
							"api_hostname": {
								Type:        schema.TypeString,
								Required:    true,
								Description: "Hostname for the PingOne management service API.  Default value can be set with the `PINGONE_API_SERVICE_HOSTNAME` environment variable.",
							},
						},
					},
				},
			},

			DataSourcesMap: map[string]*schema.Resource{
				"pingone_certificate":                    base.DatasourceCertificate(),
				"pingone_certificate_export":             base.DatasourceCertificateExport(),
				"pingone_certificate_signing_request":    base.DatasourceCertificateSigningRequest(),
				"pingone_language":                       base.DatasourceLanguage(),
				"pingone_license":                        base.DatasourceLicense(),
				"pingone_licenses":                       base.DatasourceLicenses(),
				"pingone_role":                           base.DatasourceRole(),
				"pingone_trusted_email_domain_dkim":      base.DatasourceTrustedEmailDomainDKIM(),
				"pingone_trusted_email_domain_ownership": base.DatasourceTrustedEmailDomainOwnership(),
				"pingone_trusted_email_domain_spf":       base.DatasourceTrustedEmailDomainSPF(),

				"pingone_password_policy":    sso.DatasourcePasswordPolicy(),
				"pingone_resource_attribute": sso.DatasourceResourceAttribute(),
				"pingone_resource_scope":     sso.DatasourceResourceScope(),
			},

			ResourcesMap: map[string]*schema.Resource{
				"pingone_authorize_decision_endpoint": authorize.ResourceDecisionEndpoint(),

				"pingone_certificate":                   base.ResourceCertificate(),
				"pingone_certificate_signing_response":  base.ResourceCertificateSigningResponse(),
				"pingone_gateway":                       base.ResourceGateway(),
				"pingone_gateway_credential":            base.ResourceGatewayCredential(),
				"pingone_gateway_role_assignment":       base.ResourceGatewayRoleAssignment(),
				"pingone_image":                         base.ResourceImage(),
				"pingone_key":                           base.ResourceKey(),
				"pingone_language":                      base.ResourceLanguage(),
				"pingone_language_update":               base.ResourceLanguageUpdate(),
				"pingone_notification_template_content": base.ResourceNotificationTemplateContent(),
				"pingone_role_assignment_user":          base.ResourceRoleAssignmentUser(),

				"pingone_application":                           sso.ResourceApplication(),
				"pingone_application_sign_on_policy_assignment": sso.ResourceApplicationSignOnPolicyAssignment(),
				"pingone_application_role_assignment":           sso.ResourceApplicationRoleAssignment(),
				"pingone_group":                                 sso.ResourceGroup(),
				"pingone_group_nesting":                         sso.ResourceGroupNesting(),
				"pingone_identity_provider":                     sso.ResourceIdentityProvider(),
				"pingone_password_policy":                       sso.ResourcePasswordPolicy(),
				"pingone_resource":                              sso.ResourceResource(),
				"pingone_resource_scope":                        sso.ResourceResourceScope(),
				"pingone_resource_scope_openid":                 sso.ResourceResourceScopeOpenID(),
				"pingone_resource_scope_pingone_api":            sso.ResourceResourceScopePingOneAPI(),
				"pingone_sign_on_policy":                        sso.ResourceSignOnPolicy(),
				"pingone_sign_on_policy_action":                 sso.ResourceSignOnPolicyAction(),

				"pingone_mfa_fido_policy": mfa.ResourceFIDOPolicy(),
				"pingone_mfa_policy":      mfa.ResourceMFAPolicy(),
				"pingone_mfa_settings":    mfa.ResourceMFASettings(),
			},
		}

		p.ConfigureContextFunc = configure(version)

		return p
	}
}

func configure(version string) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		var diags diag.Diagnostics
		var config client.Config

		// Set the defaults
		debugLogMessage := "[v5] Provider parameter %s missing, defaulting to environment variable"
		if v, ok := d.Get("client_id").(string); ok && v != "" {
			config.ClientID = v
		}

		if v, ok := d.Get("client_secret").(string); ok && v != "" {
			config.ClientSecret = v
		}

		if v, ok := d.Get("environment_id").(string); ok && v != "" {
			config.EnvironmentID = v
		}

		if v, ok := d.Get("api_access_token").(string); ok && v != "" {
			config.AccessToken = v
		}

		if v, ok := d.Get("region").(string); ok && v != "" {
			config.Region = v
		}

		if v, ok := d.Get("force_delete_production_type").(bool); ok {
			config.ForceDelete = v
		} else {
			forceDelete, err := strconv.ParseBool(os.Getenv("PINGONE_FORCE_DELETE_PRODUCTION_TYPE"))
			if err != nil {
				forceDelete = false
			}
			tflog.Debug(ctx, fmt.Sprintf(debugLogMessage, "force_delete_production_type"), map[string]interface{}{
				"env_var":       "PINGONE_FORCE_DELETE_PRODUCTION_TYPE",
				"env_var_value": forceDelete,
			})
			config.ForceDelete = forceDelete
		}

		if v, ok := d.Get("http_proxy").(string); ok && v != "" {
			config.ProxyURL = &v
		}

		if v, ok := d.Get("service_endpoints").([]interface{}); ok && len(v) > 0 && v[0] != nil {
			if v, ok := d.Get("auth_hostname").(string); ok && v != "" {
				config.AuthHostnameOverride = &v
			}

			if v, ok := d.Get("api_hostname").(string); ok && v != "" {
				config.APIHostnameOverride = &v
			}
		}

		client, err := config.APIClient(ctx, version)

		if err != nil {
			return nil, diag.FromErr(err)
		}

		return client, diags
	}
}
