// Copyright © 2025 Ping Identity Corporation

// Package verify provides Terraform resources and data sources for managing PingOne Verify service configurations.
// This package includes resources for verify policies, voice phrases, and voice phrase content management.
package verify

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone"
)

// serviceClientType holds the PingOne client configuration for the verify service.
// Client provides access to the PingOne API client instance used for verify service operations.
type serviceClientType struct {
	// Client is the PingOne SDK client used to interact with PingOne Verify APIs
	Client *pingone.Client
}

// Resources returns a slice of functions that create Terraform resource instances for the verify service.
// Each function in the returned slice creates a specific resource type managed by the PingOne Verify service.
// This includes verify policies, voice phrases, and voice phrase content resources.
func Resources() []func() resource.Resource {
	return []func() resource.Resource{
		NewVerifyPolicyResource,
		NewVoicePhraseResource,
		NewVoicePhraseContentResource,
	}
}

// DataSources returns a slice of functions that create Terraform data source instances for the verify service.
// Each function in the returned slice creates a specific data source type that can read verify service configurations.
// This includes data sources for verify policies, voice phrases, and voice phrase content.
func DataSources() []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewVerifyPolicyDataSource,
		NewVerifyPoliciesDataSource,
		NewVoicePhraseDataSource,
		NewVoicePhraseContentDataSource,
		NewVoicePhraseContentsDataSource,
	}
}
