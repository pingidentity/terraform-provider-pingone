package mfa

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/patrickcping/pingone-go-sdk-v2/mfa"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
)

type serviceClientType struct {
	Client *mfa.APIClient
}

func Resources() []func() resource.Resource {
	return []func() resource.Resource{
		NewApplicationPushCredentialResource,
		NewFIDO2PolicyResource,
		NewMFAPoliciesResource,
	}
}

func DataSources() []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewMFAPoliciesDataSource,
	}
}

func PrepareClient(ctx context.Context, resourceConfig framework.ResourceType) (*mfa.APIClient, error) {

	if resourceConfig.Client.API == nil || resourceConfig.Client.API.MFAAPIClient == nil {
		return nil, fmt.Errorf("Expected the PingOne \"mfa\" client, got nil.  Please report this issue to the provider maintainers.")
	}

	tflog.Info(ctx, "PingOne provider \"mfa\" client init successful")

	return resourceConfig.Client.API.MFAAPIClient, nil

}
