package credentials

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/patrickcping/pingone-go-sdk-v2/credentials"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
)

type serviceClientType struct {
	Client *credentials.APIClient
}

func Resources() []func() resource.Resource {
	return []func() resource.Resource{
		NewCredentialIssuerProfileResource,
		NewCredentialTypeResource,
		NewDigitalWalletApplicationResource,
		NewCredentialIssuanceRuleResource,
	}
}

func DataSources() []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewCredentialIssuerProfileDataSource,
		NewDigitalWalletApplicationDataSource,
		NewDigitalWalletApplicationsDataSource,
		NewCredentialTypeDataSource,
		NewCredentialTypesDataSource,
		NewCredentialIssuanceRuleDataSource,
	}
}

func PrepareClient(ctx context.Context, resourceConfig framework.ResourceType) (*credentials.APIClient, error) {

	if resourceConfig.Client.API == nil || resourceConfig.Client.API.CredentialsAPIClient == nil {
		return nil, fmt.Errorf("Expected the PingOne \"credentials\" client, got nil.  Please report this issue to the provider maintainers.")
	}

	tflog.Info(ctx, "PingOne provider \"credentials\" client init successful")

	return resourceConfig.Client.API.CredentialsAPIClient, nil

}
