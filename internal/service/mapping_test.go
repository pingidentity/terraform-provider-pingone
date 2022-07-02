package service

import (
	"testing"
)

func TestServiceFromProviderCode(t *testing.T) {

	codeTests := map[string]string{
		"SSO":            "PING_ONE_BASE",
		"MFA":            "PING_ONE_MFA",
		"PingFederate":   "PING_FEDERATE",
		"DOES_NOT_EXIST": "",
	}

	for providerCode, platformCode := range codeTests {

		v, err := ServiceFromProviderCode(providerCode)
		if err != nil && platformCode != "" {
			t.Fatalf("serviceFromProviderCode errored with %v but expected %s", err, platformCode)
		} else {

			if v.PlatformCode != platformCode {
				t.Fatalf("serviceFromProviderCode resulted in %v, expected %s", v, platformCode)
			}
		}
	}
}

func TestServiceFromPlatformCode(t *testing.T) {

	codeTests := map[string]string{
		"PING_ONE_BASE":  "SSO",
		"PING_ONE_MFA":   "MFA",
		"PING_FEDERATE":  "PingFederate",
		"DOES_NOT_EXIST": "",
	}

	for platformCode, providerCode := range codeTests {

		v, err := ServiceFromPlatformCode(platformCode)

		if err != nil && providerCode != "" {
			t.Fatalf("serviceFromPlatformCode errored with %v but expected %s", err, providerCode)
		} else {

			if v.ProviderCode != providerCode {
				t.Fatalf("serviceFromPlatformCode resulted in %v, expected %s", v, providerCode)
			}
		}
	}
}
