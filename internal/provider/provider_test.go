package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// providerFactories are used to instantiate a provider during acceptance testing.
// The factory function will be invoked for every Terraform CLI command executed
// to create a provider server to which the CLI can reattach.
var providerFactories = map[string]func() (*schema.Provider, error){
	"pingone": func() (*schema.Provider, error) {
		return New("dev")(), nil
	},
}

func TestProvider(t *testing.T) {
	if err := New("dev")().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {

	if v := os.Getenv("PINGONE_CLIENT_ID"); v == "" {
		t.Fatal("PINGONE_CLIENT_ID is missing and must be set")
	}

	if v := os.Getenv("PINGONE_CLIENT_SECRET"); v == "" {
		t.Fatal("PINGONE_CLIENT_SECRET is missing and must be set")
	}

	if v := os.Getenv("PINGONE_ENVIRONMENT_ID"); v == "" {
		t.Fatal("PINGONE_ENVIRONMENT_ID is missing and must be set")
	}

	if v := os.Getenv("PINGONE_REGION"); v == "" {
		t.Fatal("PINGONE_REGION is missing and must be set")
	}

}

func testAccPreCheckEnvironment(t *testing.T) {

	testAccPreCheck(t)
	if v := os.Getenv("PINGONE_LICENSE_ID"); v == "" {
		t.Fatal("PINGONE_LICENSE_ID is missing and must be set")
	}
}

func resourceNameGen() string {
	return acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
}

func TestServiceFromProviderCode(t *testing.T) {

	codeTests := map[string]string{
		"SSO":            "PING_ONE_BASE",
		"MFA":            "PING_ONE_MFA",
		"PING_FEDERATE":  "PING_FEDERATE",
		"DOES_NOT_EXIST": "",
	}

	for providerCode, platformCode := range codeTests {

		v, err := serviceFromProviderCode(providerCode)
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
		"PING_FEDERATE":  "PING_FEDERATE",
		"DOES_NOT_EXIST": "",
	}

	for platformCode, providerCode := range codeTests {

		v, err := serviceFromPlatformCode(platformCode)

		if err != nil && providerCode != "" {
			t.Fatalf("serviceFromPlatformCode errored with %v but expected %s", err, providerCode)
		} else {

			if v.ProviderCode != providerCode {
				t.Fatalf("serviceFromPlatformCode resulted in %v, expected %s", v, providerCode)
			}
		}
	}
}
