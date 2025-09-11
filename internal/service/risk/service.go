// Copyright © 2025 Ping Identity Corporation

// Package risk provides Terraform resources and data sources for managing PingOne Risk service configurations.
// This package includes resources for risk policies and risk predictors used in risk-based authentication and fraud detection.
package risk

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone"
)

// serviceClientType holds the PingOne client configuration for the risk service.
// Client provides access to the PingOne API client instance used for risk service operations.
type serviceClientType struct {
	// Client is the PingOne SDK client used to interact with PingOne Risk APIs
	Client *pingone.Client
}

// Resources returns a slice of functions that create Terraform resource instances for the risk service.
// Each function in the returned slice creates a specific resource type managed by the PingOne Risk service.
// This includes risk policies and risk predictors for risk-based authentication and fraud detection.
func Resources() []func() resource.Resource {
	return []func() resource.Resource{
		NewRiskPolicyResource,
		NewRiskPredictorResource,
	}
}

// DataSources returns a slice of functions that create Terraform data source instances for the risk service.
// Each function in the returned slice creates a specific data source type that can read risk service configurations.
// Currently, this service does not provide any data sources.
func DataSources() []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}
