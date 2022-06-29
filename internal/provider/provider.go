package provider

import (
	"context"
	"fmt"
	"log"

	"golang.org/x/exp/slices"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	// schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
	// 	desc := s.Description
	// 	if s.Default != nil {
	// 		desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
	// 	}
	// 	return strings.TrimSpace(desc)
	// }
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
					Description: "Client ID for the worker app client",
				},
				"client_secret": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("PINGONE_CLIENT_SECRET", nil),
					Description: "Client secret for the worker app client",
				},
				"environment_id": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("PINGONE_ENVIRONMENT_ID", nil),
					Description: "Environment ID for the worker app client",
				},
				"region": {
					Type:         schema.TypeString,
					Required:     true,
					DefaultFunc:  schema.EnvDefaultFunc("PINGONE_REGION", nil),
					Description:  "The PingOne region to use.  Options are EU, US, ASIA, CA",
					ValidateFunc: validation.StringInSlice([]string{"EU", "US", "ASIA", "CA"}, false),
				},
				"force_delete_production_type": {
					Type:        schema.TypeBool,
					Optional:    true,
					Default:     false,
					Description: "Choose whether to force-delete any configuration that has a `PRODUCTION` type parameter.  By default, `PRODUCTION` type configuration will not destroy to protect stored data",
				},
			},

			DataSourcesMap: map[string]*schema.Resource{
				"pingone_environment": datasourcePingOneEnvironment(),
			},

			ResourcesMap: map[string]*schema.Resource{
				"pingone_environment": resourcePingOneEnvironment(),
				"pingone_population":  resourcePingOnePopulation(),
			},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {

		log.Printf("[INFO] PingOne Client configuring")

		config := &Config{
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

type ServiceMapping struct {
	PlatformCode  string
	ProviderCode  string
	SolutionType  string
	ConflictsWith []string
}

func servicesMapping() []ServiceMapping {

	return []ServiceMapping{
		{
			PlatformCode: "PING_ONE_BASE",
			ProviderCode: "SSO",
		},
		{
			PlatformCode: "PING_ONE_PROVISIONING",
			ProviderCode: "SSO_PROVISIONING",
		},
		{
			PlatformCode: "PING_ONE_MFA",
			ProviderCode: "MFA",
		},
		{
			PlatformCode: "PING_ONE_RISK",
			ProviderCode: "RISK",
		},
		{
			PlatformCode: "PING_ONE_VERIFY",
			ProviderCode: "VERIFY",
		},
		{
			PlatformCode: "PING_ONE_CREDENTIALS",
			ProviderCode: "CREDENTIALS",
		},
		{
			PlatformCode: "PING_INTELLIGENCE",
			ProviderCode: "API_INTELLIGENCE",
		},
		{
			PlatformCode: "PING_ONE_AUTHORIZE",
			ProviderCode: "AUTHORIZE",
		},
		{
			PlatformCode: "PING_ONE_FRAUD",
			ProviderCode: "FRAUD",
		},
		{
			PlatformCode: "PING_ID",
			ProviderCode: "PING_ID",
		},
		{
			PlatformCode: "PING_FEDERATE",
			ProviderCode: "PING_FEDERATE",
		},
		{
			PlatformCode: "PING_ACCESS",
			ProviderCode: "PING_ACCESS",
		},
		{
			PlatformCode: "PING_DIRECTORY",
			ProviderCode: "PING_DIRECTORY",
		},
		{
			PlatformCode: "PING_AUTHORIZE",
			ProviderCode: "PING_AUTHORIZE",
		},
		{
			PlatformCode: "PING_CENTRAL",
			ProviderCode: "PING_CENTRAL",
		},
	}

}

func serviceFromProviderCode(providerCode string) (ServiceMapping, error) {

	idx := slices.IndexFunc(servicesMapping(), func(c ServiceMapping) bool { return c.ProviderCode == providerCode })

	if idx < 0 {
		return ServiceMapping{}, fmt.Errorf("Cannot find service by provider code %s", providerCode)
	}

	return servicesMapping()[idx], nil
}

func serviceFromPlatformCode(platformCode string) (ServiceMapping, error) {

	idx := slices.IndexFunc(servicesMapping(), func(c ServiceMapping) bool { return c.PlatformCode == platformCode })

	if idx < 0 {
		return ServiceMapping{}, fmt.Errorf("Cannot find service by provider code %s", platformCode)
	}

	return servicesMapping()[idx], nil
}
