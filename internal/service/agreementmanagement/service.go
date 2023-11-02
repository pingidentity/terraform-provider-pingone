package agreementmanagement

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/patrickcping/pingone-go-sdk-v2/agreementmanagement"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
)

func Resources() []func() resource.Resource {
	return []func() resource.Resource{}
}

func DataSources() []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func PrepareClient(ctx context.Context, resourceConfig framework.ResourceType) (*agreementmanagement.APIClient, error) {

	if resourceConfig.Client.API == nil || resourceConfig.Client.API.AgreementManagementAPIClient == nil {
		return nil, fmt.Errorf("Expected the PingOne \"agreement management\" client, got nil.  Please report this issue to the provider maintainers.")
	}

	tflog.Info(ctx, "PingOne provider \"agreement management\" client init successful")

	return resourceConfig.Client.API.AgreementManagementAPIClient, nil

}
