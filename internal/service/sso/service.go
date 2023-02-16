package sso

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func Resources() []func() resource.Resource {
	return []func() resource.Resource{}
}

func DataSources() []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}
