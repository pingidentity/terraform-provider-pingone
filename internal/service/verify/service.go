package verify

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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

func PrepareClient(ctx context.Context, resourceConfig framework.ResourceType) (*verify.APIClient, error) {

	if resourceConfig.Client.API == nil || resourceConfig.Client.API.VerifyAPIClient == nil {
		return nil, fmt.Errorf("Expected the PingOne \"verify\" client, got nil.  Please report this issue to the provider maintainers.")
	}

	tflog.Info(ctx, "PingOne provider \"verify\" client init successful")

	return resourceConfig.Client.API.VerifyAPIClient, nil

}
