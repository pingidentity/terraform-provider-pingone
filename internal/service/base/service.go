// Copyright © 2025 Ping Identity Corporation

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
		NewAlertChannelResource,
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
		NewLanguageTranslationResource,
		NewNotificationPolicyResource,
		NewNotificationSettingsEmailResource,
		NewNotificationSettingsResource,
		NewNotificationTemplateContentResource,
		NewPhoneDeliverySettingsResource,
		NewRoleAssignmentUserResource,
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
		NewLicenseDataSource,
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
