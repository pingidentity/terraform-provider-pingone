package risk

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/patrickcping/pingone-go-sdk-v2/risk"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
)

type serviceClientType struct {
	Client *risk.APIClient
}

func Resources() []func() resource.Resource {
	return []func() resource.Resource{
		NewRiskPolicyResource,
		NewRiskPredictorResource,
	}
}

func DataSources() []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func PrepareClient(ctx context.Context, resourceConfig framework.ResourceType) (*risk.APIClient, error) {

	if resourceConfig.Client.API == nil || resourceConfig.Client.API.RiskAPIClient == nil {
		return nil, fmt.Errorf("Expected the PingOne \"risk\" client, got nil.  Please report this issue to the provider maintainers.")
	}

	tflog.Info(ctx, "PingOne provider \"risk\" client init successful")

	return resourceConfig.Client.API.RiskAPIClient, nil

}
