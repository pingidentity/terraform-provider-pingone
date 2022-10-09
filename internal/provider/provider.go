package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
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
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("PINGONE_CLIENT_ID", nil),
					Description: "Client ID for the worker app client.  Default value can be set with the `PINGONE_CLIENT_ID` environment variable.",
				},
				"client_secret": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("PINGONE_CLIENT_SECRET", nil),
					Description: "Client secret for the worker app client.  Default value can be set with the `PINGONE_CLIENT_SECRET` environment variable.",
				},
				"environment_id": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("PINGONE_ENVIRONMENT_ID", nil),
					Description: "Environment ID for the worker app client.  Default value can be set with the `PINGONE_ENVIRONMENT_ID` environment variable.",
				},
				"region": {
					Type:             schema.TypeString,
					Required:         true,
					DefaultFunc:      schema.EnvDefaultFunc("PINGONE_REGION", nil),
					Description:      "The PingOne region to use.  Options are `AsiaPacific` `Canada` `Europe` and `NorthAmerica`.  Default value can be set with the `PINGONE_REGION` environment variable.",
					ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(model.RegionsAvailableList(), false)),
				},
				"force_delete_production_type": {
					Type:        schema.TypeBool,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("PINGONE_FORCE_DELETE_PRODUCTION_TYPE", false),
					Description: "Choose whether to force-delete any configuration that has a `PRODUCTION` type parameter.  The platform default is that `PRODUCTION` type configuration will not destroy without intervention to protect stored data.  By default this parameter is set to `false` and can be overridden with the `PINGONE_FORCE_DELETE_PRODUCTION_TYPE` environment variable.",
				},
			},

			DataSourcesMap: map[string]*schema.Resource{
				"pingone_certificate":                    base.DatasourceCertificate(),
				"pingone_certificate_export":             base.DatasourceCertificateExport(),
				"pingone_certificate_signing_request":    base.DatasourceCertificateSigningRequest(),
				"pingone_environment":                    base.DatasourceEnvironment(),
				"pingone_license":                        base.DatasourceLicense(),
				"pingone_licenses":                       base.DatasourceLicenses(),
				"pingone_language":                       base.DatasourceLanguage(),
				"pingone_role":                           base.DatasourceRole(),
				"pingone_trusted_email_domain_dkim":      base.DatasourceTrustedEmailDomainDKIM(),
				"pingone_trusted_email_domain_ownership": base.DatasourceTrustedEmailDomainOwnership(),
				"pingone_trusted_email_domain_spf":       base.DatasourceTrustedEmailDomainSPF(),

				"pingone_password_policy": sso.DatasourcePasswordPolicy(),
				"pingone_resource":        sso.DatasourceResource(),
				"pingone_resource_scope":  sso.DatasourceResourceScope(),
				"pingone_schema":          sso.DatasourceSchema(),
				"pingone_user":            sso.DatasourceUser(),
				"pingone_users":           sso.DatasourceUsers(),
			},

			ResourcesMap: map[string]*schema.Resource{
				"pingone_authorize_decision_endpoint": authorize.ResourceDecisionEndpoint(),

				"pingone_certificate":                  base.ResourceCertificate(),
				"pingone_certificate_signing_response": base.ResourceCertificateSigningResponse(),
				"pingone_custom_domain":                base.ResourceCustomDomain(),
				"pingone_custom_domain_verify":         base.ResourceCustomDomainVerify(),
				"pingone_custom_domain_ssl":            base.ResourceCustomDomainSSL(),
				"pingone_environment":                  base.ResourceEnvironment(),
				"pingone_gateway":                      base.ResourceGateway(),
				"pingone_gateway_credential":           base.ResourceGatewayCredential(),
				"pingone_gateway_role_assignment":      base.ResourceGatewayRoleAssignment(),
				"pingone_key":                          base.ResourceKey(),
				"pingone_language":                     base.ResourceLanguage(),
				"pingone_language_update":              base.ResourceLanguageUpdate(),
				"pingone_role_assignment_user":         base.ResourceRoleAssignmentUser(),
				"pingone_trusted_email_domain":         base.ResourceTrustedEmailDomain(),
				"pingone_webhook":                      base.ResourceWebhook(),

				"pingone_application":                           sso.ResourceApplication(),
				"pingone_application_attribute_mapping":         sso.ResourceApplicationAttributeMapping(),
				"pingone_application_sign_on_policy_assignment": sso.ResourceApplicationSignOnPolicyAssignment(),
				"pingone_application_resource_grant":            sso.ResourceApplicationResourceGrant(),
				"pingone_application_role_assignment":           sso.ResourceApplicationRoleAssignment(),
				"pingone_identity_provider":                     sso.ResourceIdentityProvider(),
				"pingone_identity_provider_attribute":           sso.ResourceIdentityProviderAttribute(),
				"pingone_group":                                 sso.ResourceGroup(),
				"pingone_group_nesting":                         sso.ResourceGroupNesting(),
				"pingone_password_policy":                       sso.ResourcePasswordPolicy(),
				"pingone_population":                            sso.ResourcePopulation(),
				"pingone_resource":                              sso.ResourceResource(),
				"pingone_resource_scope":                        sso.ResourceResourceScope(),
				"pingone_schema_attribute":                      sso.ResourceSchemaAttribute(),
				"pingone_sign_on_policy":                        sso.ResourceSignOnPolicy(),
				"pingone_sign_on_policy_action":                 sso.ResourceSignOnPolicyAction(),
				"pingone_user":                                  sso.ResourceUser(),

				"pingone_mfa_application_push_credential": mfa.ResourceApplicationPushCredential(),
				"pingone_mfa_policy":                      mfa.ResourceMFAPolicy(),
				"pingone_mfa_settings":                    mfa.ResourceMFASettings(),
			},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {

		config := &client.Config{
			ClientID:      d.Get("client_id").(string),
			ClientSecret:  d.Get("client_secret").(string),
			EnvironmentID: d.Get("environment_id").(string),
			Region:        d.Get("region").(string),
			ForceDelete:   d.Get("force_delete_production_type").(bool),
		}

		client, err := config.APIClient(ctx)

		if err != nil {
			return nil, diag.FromErr(err)
		}

		return client, nil
	}
}
