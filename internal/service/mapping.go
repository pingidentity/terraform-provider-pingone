package service

import (
	"fmt"

	"golang.org/x/exp/slices"
)

type ServiceMapping struct {
	PlatformCode  string
	ProviderCode  string
	SolutionType  string
	ConflictsWith []string
}

func servicesMapping() []ServiceMapping {

	return []ServiceMapping{
		{
			PlatformCode: "PING_ONE_BASE",
			ProviderCode: "SSO",
		},
		{
			PlatformCode: "PING_ONE_PROVISIONING",
			ProviderCode: "SSO_PROVISIONING",
		},
		{
			PlatformCode: "PING_ONE_MFA",
			ProviderCode: "MFA",
		},
		{
			PlatformCode: "PING_ONE_RISK",
			ProviderCode: "Risk",
		},
		{
			PlatformCode: "PING_ONE_VERIFY",
			ProviderCode: "Verify",
		},
		{
			PlatformCode: "PING_ONE_CREDENTIALS",
			ProviderCode: "Credentials",
		},
		{
			PlatformCode: "PING_INTELLIGENCE",
			ProviderCode: "APIIntelligence",
		},
		{
			PlatformCode: "PING_ONE_AUTHORIZE",
			ProviderCode: "Authorize",
		},
		{
			PlatformCode: "PING_ONE_FRAUD",
			ProviderCode: "Fraud",
		},
		{
			PlatformCode: "PING_ID",
			ProviderCode: "PingID",
		},
		{
			PlatformCode: "PING_FEDERATE",
			ProviderCode: "PingFederate",
		},
		{
			PlatformCode: "PING_ACCESS",
			ProviderCode: "PingAccess",
		},
		{
			PlatformCode: "PING_DIRECTORY",
			ProviderCode: "PingDirectory",
		},
		{
			PlatformCode: "PING_AUTHORIZE",
			ProviderCode: "PingAuthorize",
		},
		{
			PlatformCode: "PING_CENTRAL",
			ProviderCode: "PingCentral",
		},
	}

}

func ServiceFromProviderCode(providerCode string) (ServiceMapping, error) {

	idx := slices.IndexFunc(servicesMapping(), func(c ServiceMapping) bool { return c.ProviderCode == providerCode })

	if idx < 0 {
		return ServiceMapping{}, fmt.Errorf("Cannot find service by provider code %s", providerCode)
	}

	return servicesMapping()[idx], nil
}

func ServiceFromPlatformCode(platformCode string) (ServiceMapping, error) {

	idx := slices.IndexFunc(servicesMapping(), func(c ServiceMapping) bool { return c.PlatformCode == platformCode })

	if idx < 0 {
		return ServiceMapping{}, fmt.Errorf("Cannot find service by provider code %s", platformCode)
	}

	return servicesMapping()[idx], nil
}
