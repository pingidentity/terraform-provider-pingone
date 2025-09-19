// Copyright © 2025 Ping Identity Corporation

// Package credentials provides Terraform resources and data sources for managing PingOne Credentials service configurations.
// This package includes resources for credential issuer profiles, credential types, digital wallet applications, and credential issuance rules.
package credentials

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone"
)

// serviceClientType holds the PingOne client configuration for the credentials service.
// Client provides access to the PingOne API client instance used for credentials service operations.
type serviceClientType struct {
	// Client is the PingOne SDK client used to interact with PingOne Credentials APIs
	Client *pingone.Client
}

// Resources returns a slice of functions that create Terraform resource instances for the credentials service.
// Each function in the returned slice creates a specific resource type managed by the PingOne Credentials service.
// This includes credential issuer profiles, credential types, digital wallet applications, and credential issuance rules.
func Resources() []func() resource.Resource {
	return []func() resource.Resource{
		NewCredentialIssuerProfileResource,
		NewCredentialTypeResource,
		NewDigitalWalletApplicationResource,
		NewCredentialIssuanceRuleResource,
	}
}

// DataSources returns a slice of functions that create Terraform data source instances for the credentials service.
// Each function in the returned slice creates a specific data source type that can read credentials service configurations.
// This includes data sources for credential issuer profiles, digital wallet applications, credential types, and credential issuance rules.
func DataSources() []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewCredentialIssuerProfileDataSource,
		NewDigitalWalletApplicationDataSource,
		NewDigitalWalletApplicationsDataSource,
		NewCredentialTypeDataSource,
		NewCredentialTypesDataSource,
		NewCredentialIssuanceRuleDataSource,
	}
}
