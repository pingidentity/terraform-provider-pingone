// Copyright © 2026 Ping Identity Corporation

package davinci

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/pingidentity/pingone-go-client/pingone"
)

type serviceClientType struct {
	Client *pingone.APIClient
}

func Resources() []func() resource.Resource {
	resources := []func() resource.Resource{}
	resources = append(resources, BetaResources()...)

	return resources
}

func DataSources() []func() datasource.DataSource {
	dataSources := []func() datasource.DataSource{}
	dataSources = append(dataSources, BetaDataSources()...)

	return dataSources
}
