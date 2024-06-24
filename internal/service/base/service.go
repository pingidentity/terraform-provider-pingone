package base

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
		NewAgreementEnableResource,
		NewAgreementLocalizationEnableResource,
		NewAgreementLocalizationResource,
		NewAgreementLocalizationRevisionResource,
		NewAgreementResource,
		NewBrandingSettingsResource,
		NewBrandingThemeDefaultResource,
		NewBrandingThemeResource,
		NewCustomDomainResource,
		NewCustomDomainSSLResource,
		NewCustomDomainVerifyResource,
		NewEnvironmentResource,
		NewFormResource,
		NewFormsRecaptchaV2Resource,
		NewGatewayResource,
		NewGatewayRoleAssignmentResource,
		NewIdentityPropagationPlanResource,
		NewImageResource,
		NewKeyResource,
		NewKeyRotationPolicyResource,
		NewNotificationPolicyResource,
		NewNotificationSettingsEmailResource,
		NewNotificationSettingsResource,
		NewNotificationTemplateContentResource,
		NewPhoneDeliverySettingsResource,
		NewSystemApplicationResource,
		NewTrustedEmailAddressResource,
		NewTrustedEmailDomainResource,
		NewUserRoleAssignmentResource,
		NewWebhookResource,
	}
}

func DataSources() []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewAgreementDataSource,
		NewAgreementLocalizationDataSource,
		NewEnvironmentDataSource,
		NewEnvironmentsDataSource,
		NewGatewayDataSource,
		NewLicensesDataSource,
		NewOrganizationDataSource,
		NewPhoneDeliverySettingsListDataSource,
		NewRoleDataSource,
		NewRolesDataSource,
		NewTrustedEmailDomainDataSource,
		NewTrustedEmailDomainDKIMDataSource,
		NewTrustedEmailDomainOwnershipDataSource,
		NewUserRoleAssignmentsDataSource,
	}
}
