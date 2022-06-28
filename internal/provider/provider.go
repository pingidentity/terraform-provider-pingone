package provider

import (
	"context"
	"log"

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

type apiClient struct {
	// Add whatever fields, client or connection info, etc. here
	// you would need to setup to communicate with the upstream
	// API.
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {

		log.Printf("[INFO] PingOne Client configuring")

		config := &Config{
			ClientID:      d.Get("client_id").(string),
			ClientSecret:  d.Get("client_secret").(string),
			EnvironmentID: d.Get("environment_id").(string),
			Region:        d.Get("region").(string),
		}

		client, err := config.APIClient(ctx)

		if err != nil {
			return nil, diag.FromErr(err)
		}

		return client, nil
	}
}
