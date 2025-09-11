// Copyright © 2025 Ping Identity Corporation

// Package mfa provides Terraform resources and data sources for managing PingOne MFA (Multi-Factor Authentication) service configurations.
// This package includes resources for MFA device policies, FIDO2 policies, application push credentials, and MFA settings.
package mfa

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone"
)

// serviceClientType holds the PingOne client configuration for the MFA service.
// Client provides access to the PingOne API client instance used for MFA service operations.
type serviceClientType struct {
	// Client is the PingOne SDK client used to interact with PingOne MFA APIs
	Client *pingone.Client
}

// Resources returns a slice of functions that create Terraform resource instances for the MFA service.
// Each function in the returned slice creates a specific resource type managed by the PingOne MFA service.
// This includes MFA device policies, FIDO2 policies, application push credentials, and MFA settings.
func Resources() []func() resource.Resource {
	return []func() resource.Resource{
		NewApplicationPushCredentialResource,
		NewFIDO2PolicyResource,
		NewMFADevicePolicyResource,
		NewMFASettingsResource,
	}
}

// DataSources returns a slice of functions that create Terraform data source instances for the MFA service.
// Each function in the returned slice creates a specific data source type that can read MFA service configurations.
// This includes data sources for MFA device policies and related configurations.
func DataSources() []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewMFADevicePoliciesDataSource,
	}
}
