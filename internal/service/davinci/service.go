// Copyright © 2025 Ping Identity Corporation

package davinci

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	oldsdk "github.com/patrickcping/pingone-go-sdk-v2/pingone"
	"github.com/pingidentity/pingone-go-client/pingone"
)

type serviceClientType struct {
	Client           *pingone.APIClient
	ManagementClient *oldsdk.Client
}

func Resources() []func() resource.Resource {
	return []func() resource.Resource{
		NewDavinciVariableResource,
	}
}

func DataSources() []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}
