// Copyright © 2026 Ping Identity Corporation

package credentials

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone"
)

type serviceClientType struct {
	Client *pingone.Client
}

func Resources() []func() resource.Resource {
	resources := []func() resource.Resource{
		NewCredentialIssuerProfileResource,
		NewCredentialTypeResource,
		NewDigitalWalletApplicationResource,
		NewCredentialIssuanceRuleResource,
	}
	resources = append(resources, BetaResources()...)

	return resources
}

func DataSources() []func() datasource.DataSource {
	dataSources := []func() datasource.DataSource{
		NewCredentialIssuerProfileDataSource,
		NewDigitalWalletApplicationDataSource,
		NewDigitalWalletApplicationsDataSource,
		NewCredentialTypeDataSource,
		NewCredentialTypesDataSource,
		NewCredentialIssuanceRuleDataSource,
	}
	dataSources = append(dataSources, BetaDataSources()...)

	return dataSources
}
