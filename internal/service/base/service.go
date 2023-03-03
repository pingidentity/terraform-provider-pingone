package base

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/patrickcping/pingone-go-sdk-v2/agreementmanagement"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
)

func Resources() []func() resource.Resource {
	return []func() resource.Resource{
		NewAgreementResource,
		NewAgreementEnableResource,
		NewAgreementLocalizationResource,
		NewAgreementLocalizationEnableResource,
		NewAgreementLocalizationRevisionResource,
		NewTrustedEmailAddressResource,
		NewNotificationPolicyResource,
		NewNotificationSettingsEmailResource,
	}
}

func DataSources() []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewAgreementDataSource,
		NewAgreementLocalizationDataSource,
		NewTrustedEmailDomainDataSource,
	}
}

func prepareClient(ctx context.Context, resourceConfig framework.ResourceType) (*management.APIClient, error) {

	if resourceConfig.Client.API == nil || resourceConfig.Client.API.ManagementAPIClient == nil {
		return nil, fmt.Errorf("Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
	}

	tflog.Info(ctx, "PingOne provider client init successful")

	return resourceConfig.Client.API.ManagementAPIClient, nil

}

func prepareAgreementMgmtClient(ctx context.Context, resourceConfig framework.ResourceType) (*agreementmanagement.APIClient, error) {

	if resourceConfig.Client.API == nil || resourceConfig.Client.API.AgreementManagementAPIClient == nil {
		return nil, fmt.Errorf("Expected the PingOne client, got nil.  Please report this issue to the provider maintainers.")
	}

	tflog.Info(ctx, "PingOne provider client init successful")

	return resourceConfig.Client.API.AgreementManagementAPIClient, nil

}
