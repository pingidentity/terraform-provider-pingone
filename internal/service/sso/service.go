package sso

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone"
)

type serviceClientType struct {
	Client *pingone.Client
}

func Resources() []func() resource.Resource {
	return []func() resource.Resource{
		NewApplicationAttributeMappingResource,
		NewApplicationFlowPolicyAssignmentResource,
		NewApplicationResourceGrantResource,
		NewGroupResource,
		NewGroupNestingResource,
		NewIdentityProviderAttributeResource,
		NewPopulationResource,
		NewPopulationDefaultResource,
		NewResourceAttributeResource,
		NewResourceScopeResource,
		NewResourceScopeOpenIDResource,
		NewResourceScopePingOneAPIResource,
		NewSchemaAttributeResource,
		NewUserResource,
	}
}

func DataSources() []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewFlowPoliciesDataSource,
		NewFlowPolicyDataSource,
		NewPopulationDataSource,
		NewPopulationsDataSource,
		NewResourceDataSource,
		NewResourceScopeDataSource,
		NewSchemaDataSource,
		NewUserDataSource,
		NewUsersDataSource,
	}
}
