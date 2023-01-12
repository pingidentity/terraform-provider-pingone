package framework

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
)

// providerFactories are used to instantiate a provider during acceptance testing.
// The factory function will be invoked for every Terraform CLI command executed
// to create a provider server to which the CLI can reattach.
var providerFactories = map[string]func() (provider.Provider, error){
	"pingone": func() (provider.Provider, error) {
		provider := New("dev")()

		if provider == nil {
			return nil, fmt.Errorf("Cannot initiate provider factory")
		}
		return provider, nil
	},
}

// func TestProvider(t *testing.T) {
// 	if err := New("dev")().InternalValidate(); err != nil {
// 		t.Fatalf("err: %s", err)
// 	}
// }

func testAccPreCheck(t *testing.T) {
	vars := []string{"PINGONE_CLIENT_ID", "PINGONE_CLIENT_SECRET", "PINGONE_ENVIRONMENT_ID", "PINGONE_REGION"}
	for _, v := range vars {
		if os.Getenv(v) == "" {
			t.Fatalf("%s is missing and must be set", v)
		}
	}
}
