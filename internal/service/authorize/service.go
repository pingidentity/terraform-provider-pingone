// Copyright © 2026 Ping Identity Corporation

package authorize

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/pingidentity/pingone-go-sdk-legacy/pingone"
)

type serviceClientType struct {
	Client *pingone.Client
}

func Resources() []func() resource.Resource {
	resources := []func() resource.Resource{
		NewAPIServiceDeploymentResource,
		NewAPIServiceOperationResource,
		NewAPIServiceResource,
		NewApplicationResourcePermissionResource,
		NewApplicationRolePermissionResource,
		NewApplicationRoleResource,
	}
	resources = append(resources, BetaResources()...)

	return resources
}

func DataSources() []func() datasource.DataSource {
	dataSources := []func() datasource.DataSource{}
	dataSources = append(dataSources, BetaDataSources()...)

	return dataSources
}
