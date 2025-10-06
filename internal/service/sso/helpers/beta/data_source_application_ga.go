// Copyright © 2025 Ping Identity Corporation

// This file relates to a beta feature described in CDI-492 and should be modified or removed on completion of CDI-631

//go:build !beta

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
	}
}
