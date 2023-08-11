package verify

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/verify"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
)

func Resources() []func() resource.Resource {
	return []func() resource.Resource{
		NewVerifyPolicyResource,
		NewVoicePhraseResource,
		NewVoicePhraseContentResource,
	}
}

func DataSources() []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewVerifyPolicyDataSource,
		NewVerifyPoliciesDataSource,
		NewVoicePhraseDataSource,
		NewVoicePhraseContentDataSource,
		NewVoicePhraseContentsDataSource,
	}
}

func prepareClient(ctx context.Context, resourceConfig framework.ResourceType) (*verify.APIClient, error) {

	if resourceConfig.Client.API == nil || resourceConfig.Client.API.VerifyAPIClient == nil {
		return nil, fmt.Errorf("Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
	}

	tflog.Info(ctx, "PingOne provider client init successful")

	return resourceConfig.Client.API.VerifyAPIClient, nil

}

func prepareMgmtClient(ctx context.Context, resourceConfig framework.ResourceType) (*management.APIClient, error) {

	if resourceConfig.Client.API == nil || resourceConfig.Client.API.ManagementAPIClient == nil {
		return nil, fmt.Errorf("Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
	}

	tflog.Info(ctx, "PingOne provider client init successful")

	return resourceConfig.Client.API.ManagementAPIClient, nil

}
