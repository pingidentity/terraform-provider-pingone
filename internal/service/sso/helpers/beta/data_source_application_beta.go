// Copyright © 2025 Ping Identity Corporation

// This file relates to a beta feature described in CDI-492

//go:build beta

package beta

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
)

func DataSourceSchemaItems() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"client_id": schema.StringAttribute{
			Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the application ID used to authenticate to the authorization server.").Description,
			Computed:    true,
		},
		"initial_client_secret": schema.StringAttribute{
			Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the initial import client secret used to authenticate to the authorization server. If left undefined, the service will generate a value. Note this field's value will not change if the secret is rotated. After initial import, the `pingone_application_secret` resource or data source should be used to refer to the current active value of the application's secret.").Beta("To modify the value of this field, the environment must be enabled with the feature flag to allow importing applications with administrator defined client ID and client secret values.").Description,
			Computed:    true,
			Sensitive:   true,
		},
	}
}
