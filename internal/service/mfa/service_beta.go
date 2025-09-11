// Copyright © 2025 Ping Identity Corporation

//go:build beta

package mfa

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func BetaResources() []func() resource.Resource {
	return []func() resource.Resource{}
}

func BetaDataSources() []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}
