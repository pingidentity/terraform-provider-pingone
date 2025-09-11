// Copyright © 2025 Ping Identity Corporation

// Package sso provides Terraform resources and data sources for managing PingOne SSO (Single Sign-On) service configurations.
// This package includes resources for applications, identity providers, users, groups, populations, resources, schemas, and sign-on policies.
package sso

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone"
)

// serviceClientType holds the PingOne client configuration for the SSO service.
// Client provides access to the PingOne API client instance used for SSO service operations.
type serviceClientType struct {
	// Client is the PingOne SDK client used to interact with PingOne SSO APIs
	Client *pingone.Client
}

// Resources returns a slice of functions that create Terraform resource instances for the SSO service.
// Each function in the returned slice creates a specific resource type managed by the PingOne SSO service.
// This includes applications, identity providers, users, groups, populations, resources, schemas, and sign-on policies.
func Resources() []func() resource.Resource {
	return []func() resource.Resource{
		NewAdministratorSecurityResource,
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

// DataSources returns a slice of functions that create Terraform data source instances for the SSO service.
// Each function in the returned slice creates a specific data source type that can read SSO service configurations.
// This includes data sources for applications, identity providers, users, groups, populations, resources, schemas, and policies.
func DataSources() []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewAdministratorSecurityDataSource,
		NewApplicationDataSource,
		NewApplicationFlowPolicyAssignmentsDataSource,
		NewApplicationSecretDataSource,
		NewApplicationSignOnPolicyAssignmentsDataSource,
		NewCustomRoleDataSource,
		NewCustomRolesDataSource,
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
