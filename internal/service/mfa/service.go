// Copyright © 2026 Ping Identity Corporation

package mfa

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
		NewApplicationPushCredentialResource,
		NewFIDO2PolicyResource,
		NewMFADevicePolicyResource,
		NewMFADevicePolicyDefaultResource,
		NewMFASettingsResource,
	}
	resources = append(resources, BetaResources()...)

	return resources
}

func DataSources() []func() datasource.DataSource {
	dataSources := []func() datasource.DataSource{
		NewMFADevicePoliciesDataSource,
	}
	dataSources = append(dataSources, BetaDataSources()...)

	return dataSources
}
