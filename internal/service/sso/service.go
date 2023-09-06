package sso

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
)

type serviceClientType struct {
	Client *management.APIClient
}

func Resources() []func() resource.Resource {
	return []func() resource.Resource{
		NewApplicationAttributeMappingResource,
		NewApplicationFlowPolicyAssignmentResource,
		NewApplicationResourceGrantResource,
		NewGroupResource,
		NewIdentityProviderAttributeResource,
		NewPopulationResource,
		NewResourceAttributeResource,
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
		NewSchemaDataSource,
		NewUserDataSource,
		NewUsersDataSource,
	}
}

func PrepareClient(ctx context.Context, resourceConfig framework.ResourceType) (*management.APIClient, error) {

	if resourceConfig.Client.API == nil || resourceConfig.Client.API.ManagementAPIClient == nil {
		return nil, fmt.Errorf("Expected the PingOne \"management\" client, got nil.  Please report this issue to the provider maintainers.")
	}

	tflog.Info(ctx, "PingOne provider \"management\" client init successful")

	return resourceConfig.Client.API.ManagementAPIClient, nil

}
