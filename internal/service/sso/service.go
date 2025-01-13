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
		NewApplicationResource,
		NewApplicationResourceGrantResource,
		NewApplicationResourceResource,
		NewApplicationRoleAssignmentResource,
		NewApplicationSecretResource,
		NewCustomRoleResource,
		NewGroupNestingResource,
		NewGroupResource,
		NewGroupRoleAssignmentResource,
		NewIdentityProviderAttributeResource,
		NewIdentityProviderResource,
		NewPasswordPolicyResource,
		NewPopulationDefaultIdpResource,
		NewPopulationDefaultResource,
		NewPopulationResource,
		NewResourceAttributeResource,
		NewResourceResource,
		NewResourceScopeOpenIDResource,
		NewResourceScopePingOneAPIResource,
		NewResourceScopeResource,
		NewResourceSecretResource,
		NewSchemaAttributeResource,
		NewSignOnPolicyResource,
		NewUserApplicationRoleAssignmentResource,
		NewUserGroupAssignmentResource,
		NewUserResource,
	}
}

func DataSources() []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewApplicationDataSource,
		NewApplicationFlowPolicyAssignmentsDataSource,
		NewApplicationSecretDataSource,
		NewApplicationSignOnPolicyAssignmentsDataSource,
		NewCustomRoleDataSource,
		NewFlowPoliciesDataSource,
		NewFlowPolicyDataSource,
		NewGroupDataSource,
		NewGroupsDataSource,
		NewPasswordPolicyDataSource,
		NewPopulationDataSource,
		NewPopulationsDataSource,
		NewResourceDataSource,
		NewResourceScopeDataSource,
		NewResourceSecretDataSource,
		NewSchemaDataSource,
		NewUserDataSource,
		NewUsersDataSource,
	}
}
