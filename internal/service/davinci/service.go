// Copyright © 2025 Ping Identity Corporation

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
	return []func() resource.Resource{
		NewDavinciApplicationFlowPolicyResource,
		NewDavinciApplicationKeyResource,
		NewDavinciApplicationResource,
		NewDavinciApplicationSecretResource,
		NewDavinciConnectorInstanceResource,
		NewDavinciFlowResource,
		NewDavinciFlowDeployResource,
		NewDavinciFlowEnabledResource,
		NewDavinciVariableResource,
	}
}

func DataSources() []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewDavinciApplicationDataSource,
		NewDavinciApplicationsDataSource,
		NewDavinciConnectorDataSource,
		NewDavinciConnectorsDataSource,
		NewDavinciConnectorInstanceDataSource,
		NewDavinciConnectorInstancesDataSource,
	}
}
