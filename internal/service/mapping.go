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
			ProviderCode: "RISK",
		},
		{
			PlatformCode: "PING_ONE_VERIFY",
			ProviderCode: "VERIFY",
		},
		{
			PlatformCode: "PING_ONE_CREDENTIALS",
			ProviderCode: "CREDENTIALS",
		},
		{
			PlatformCode: "PING_INTELLIGENCE",
			ProviderCode: "API_INTELLIGENCE",
		},
		{
			PlatformCode: "PING_ONE_AUTHORIZE",
			ProviderCode: "AUTHORIZE",
		},
		{
			PlatformCode: "PING_ONE_FRAUD",
			ProviderCode: "FRAUD",
		},
		{
			PlatformCode: "PING_ID",
			ProviderCode: "PING_ID",
		},
		{
			PlatformCode: "PING_FEDERATE",
			ProviderCode: "PING_FEDERATE",
		},
		{
			PlatformCode: "PING_ACCESS",
			ProviderCode: "PING_ACCESS",
		},
		{
			PlatformCode: "PING_DIRECTORY",
			ProviderCode: "PING_DIRECTORY",
		},
		{
			PlatformCode: "PING_AUTHORIZE",
			ProviderCode: "PING_AUTHORIZE",
		},
		{
			PlatformCode: "PING_CENTRAL",
			ProviderCode: "PING_CENTRAL",
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
