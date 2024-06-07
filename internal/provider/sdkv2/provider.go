package sdkv2

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/service/authorize"
	"github.com/pingidentity/terraform-provider-pingone/internal/service/base"
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
				"region_code": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The PingOne region to use, which selects the appropriate service endpoints.  Options are `AP` (for Asia-Pacific `.asia` tenants), `AU` (for Asia-Pacific `.com.au` tenants), `CA` (for Canada `.ca` tenants), `EU` (for Europe `.eu` tenants) and `NA` (for North America `.com` tenants).  Default value can be set with the `PINGONE_REGION_CODE` environment variable.",
				},
				"http_proxy": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Full URL for the http/https proxy service, for example `http://127.0.0.1:8090`.  Default value can be set with the `HTTP_PROXY` or `HTTPS_PROXY` environment variables.",
				},
				"append_user_agent": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "A custom string value to append to the end of the `User-Agent` header when making API requests to the PingOne service. Default value can be set with the `PINGONE_TF_APPEND_USER_AGENT` environment variable.",
				},
				"global_options": {
					Type:        schema.TypeList,
					Optional:    true,
					MaxItems:    1,
					Description: "A single block containing configuration items to override API behaviours in PingOne.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"population": {
								Type:        schema.TypeList,
								Optional:    true,
								MaxItems:    1,
								Description: "A single block containing configuration items to override population resource settings in PingOne.",
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"contains_users_force_delete": {
											Type:        schema.TypeBool,
											Optional:    true,
											Description: "Choose whether to force-delete populations that contain users not managed by Terraform. Useful for development and testing use cases, and only applies if the environment that contains the population is of type `SANDBOX`. The platform default is that populations cannot be removed if they contain user data. By default this parameter is set to `false`.",
										},
									},
								},
							},
						},
					},
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
				"pingone_certificate":                 base.DatasourceCertificate(),
				"pingone_certificate_export":          base.DatasourceCertificateExport(),
				"pingone_certificate_signing_request": base.DatasourceCertificateSigningRequest(),
				"pingone_language":                    base.DatasourceLanguage(),
				"pingone_license":                     base.DatasourceLicense(),
				"pingone_trusted_email_domain_spf":    base.DatasourceTrustedEmailDomainSPF(),

				"pingone_resource_attribute": sso.DatasourceResourceAttribute(),
			},

			ResourcesMap: map[string]*schema.Resource{
				"pingone_authorize_decision_endpoint": authorize.ResourceDecisionEndpoint(),

				"pingone_certificate":                   base.ResourceCertificate(),
				"pingone_certificate_signing_response":  base.ResourceCertificateSigningResponse(),
				"pingone_gateway_credential":            base.ResourceGatewayCredential(),
				"pingone_language":                      base.ResourceLanguage(),
				"pingone_language_update":               base.ResourceLanguageUpdate(),
				"pingone_notification_template_content": base.ResourceNotificationTemplateContent(),

				"pingone_application_sign_on_policy_assignment": sso.ResourceApplicationSignOnPolicyAssignment(),
				"pingone_sign_on_policy_action":                 sso.ResourceSignOnPolicyAction(),
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

		if v, ok := d.Get("region_code").(string); ok && v != "" {
			regionCode := management.EnumRegionCode(v)
			config.RegionCode = &regionCode
		}

		config.GlobalOptions = &client.GlobalOptions{
			Population: &client.PopulationOptions{
				ContainsUsersForceDelete: false,
			},
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

		if v, ok := d.Get("global_options").([]interface{}); ok && len(v) > 0 && v[0] != nil {
			if v, ok := d.Get("population").([]interface{}); ok && len(v) > 0 && v[0] != nil {
				if v, ok := d.Get("contains_users_force_delete").(bool); ok {
					config.GlobalOptions.Population.ContainsUsersForceDelete = v
				}
			}
		}

		if v, ok := d.Get("append_user_agent").(string); ok && v != "" {
			config.UserAgentAppend = &v
		} else if v := strings.TrimSpace(os.Getenv("PINGONE_TF_APPEND_USER_AGENT")); v != "" {
			config.UserAgentAppend = &v
		}

		client, err := config.APIClient(ctx, version)

		if err != nil {
			return nil, diag.FromErr(err)
		}

		return client, diags
	}
}
