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
		NewApplicationRoleAssignmentResource,
		NewGroupNestingResource,
		NewGroupResource,
		NewGroupRoleAssignmentResource,
		NewIdentityProviderAttributeResource,
		NewIdentityProviderResource,
		NewPopulationDefaultResource,
		NewPopulationResource,
		NewResourceAttributeResource,
		NewResourceScopeOpenIDResource,
		NewResourceScopePingOneAPIResource,
		NewResourceScopeResource,
		NewSchemaAttributeResource,
		NewUserResource,
	}
}

func DataSources() []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewApplicationDataSource,
		NewFlowPoliciesDataSource,
		NewFlowPolicyDataSource,
		NewGroupDataSource,
		NewGroupsDataSource,
		NewPopulationDataSource,
		NewPopulationsDataSource,
		NewResourceDataSource,
		NewResourceScopeDataSource,
		NewSchemaDataSource,
		NewUserDataSource,
		NewUsersDataSource,
	}
}
