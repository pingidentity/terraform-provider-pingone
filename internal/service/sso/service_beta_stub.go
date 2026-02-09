// Copyright © 2026 Ping Identity Corporation

//go:build !beta

package sso

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Stub method to return an empty list when beta flag is not enabled.
func BetaResources() []func() resource.Resource {
	// Do not add resources here. Beta resources should be added in service_beta.go
	return []func() resource.Resource{}
}

// Stub method to return an empty list when beta flag is not enabled
func BetaDataSources() []func() datasource.DataSource {
	// Do not add data sources here. Beta data sources should be added in service_beta.go
	return []func() datasource.DataSource{}
}
