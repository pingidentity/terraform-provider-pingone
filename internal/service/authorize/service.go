// Copyright © 2025 Ping Identity Corporation

// Package authorize provides Terraform resources and data sources for managing PingOne Authorize service configurations.
// This package includes resources for API services, API service deployments, application resources, roles, and permissions management.
package authorize

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone"
)

// serviceClientType holds the PingOne client configuration for the authorize service.
// Client provides access to the PingOne API client instance used for authorize service operations.
type serviceClientType struct {
	// Client is the PingOne SDK client used to interact with PingOne Authorize APIs
	Client *pingone.Client
}

// Resources returns a slice of functions that create Terraform resource instances for the authorize service.
// Each function in the returned slice creates a specific resource type managed by the PingOne Authorize service.
// This includes API services, API service deployments, application resources, roles, and permissions.
func Resources() []func() resource.Resource {
	return []func() resource.Resource{
		NewAPIServiceDeploymentResource,
		NewAPIServiceOperationResource,
		NewAPIServiceResource,
		NewApplicationResourcePermissionResource,
		NewApplicationRolePermissionResource,
		NewApplicationRoleResource,
	}
}

// DataSources returns a slice of functions that create Terraform data source instances for the authorize service.
// Each function in the returned slice creates a specific data source type that can read authorize service configurations.
// Currently, this service does not provide any data sources.
func DataSources() []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}
