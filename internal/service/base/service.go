// Copyright © 2025 Ping Identity Corporation

// Package base provides Terraform resources and data sources for managing core PingOne service configurations.
// This package includes resources for environments, organizations, agreements, branding, gateways, keys, notifications, and user management.
package base

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone"
)

// serviceClientType holds the PingOne client configuration for the base service.
// Client provides access to the PingOne API client instance used for base service operations.
type serviceClientType struct {
	// Client is the PingOne SDK client used to interact with core PingOne APIs
	Client *pingone.Client
}

// Resources returns a slice of functions that create Terraform resource instances for the base service.
// Each function in the returned slice creates a specific resource type managed by the core PingOne platform.
// This includes environments, agreements, branding, gateways, keys, notifications, roles, and user management resources.
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

// DataSources returns a slice of functions that create Terraform data source instances for the base service.
// Each function in the returned slice creates a specific data source type that can read core PingOne configurations.
// This includes data sources for environments, organizations, agreements, gateways, licenses, roles, and trusted email domains.
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
