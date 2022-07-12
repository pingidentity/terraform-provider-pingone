package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
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
					Description:      "The PingOne region to use.  Options are `EU`, `US`, `ASIA`, `CA`.  Default value can be set with the `PINGONE_REGION` environment variable.",
					ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"EU", "US", "ASIA", "CA"}, false)),
				},
				"force_delete_production_type": {
					Type:        schema.TypeBool,
					Optional:    true,
					Default:     false,
					Description: "Choose whether to force-delete any configuration that has a `PRODUCTION` type parameter.  By default, `PRODUCTION` type configuration will not destroy to protect stored data.",
				},
			},

			DataSourcesMap: map[string]*schema.Resource{
				"pingone_environment": base.DatasourceEnvironment(),

				"pingone_schema": sso.DatasourceSchema(),
			},

			ResourcesMap: map[string]*schema.Resource{
				"pingone_environment": base.ResourceEnvironment(),

				"pingone_group":            sso.ResourceGroup(),
				"pingone_population":       sso.ResourcePopulation(),
				"pingone_user":             sso.ResourceUser(),
				"pingone_schema_attribute": sso.ResourceSchemaAttribute(),
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
