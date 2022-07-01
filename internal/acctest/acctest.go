package acctest

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/provider"
)

// ProviderFactories is a static map containing only the main provider instance
//
// Use other ProviderFactories functions, such as FactoriesAlternate,
// for tests requiring special provider configurations.
var ProviderFactories map[string]func() (*schema.Provider, error)

// Provider is the "main" provider instance
//
// This Provider can be used in testing code for API calls without requiring
// the use of saving and referencing specific ProviderFactories instances.
//
// PreCheck(t) must be called before using this provider instance.
var Provider *schema.Provider

// providerFactories are used to instantiate a provider during acceptance testing.
// The factory function will be invoked for every Terraform CLI command executed
// to create a provider server to which the CLI can reattach.

func init() {
	Provider = provider.New("dev")()

	// Always allocate a new provider instance each invocation, otherwise gRPC
	// ProviderConfigure() can overwrite configuration during concurrent testing.
	ProviderFactories = map[string]func() (*schema.Provider, error){
		"pingone": func() (*schema.Provider, error) {
			return provider.New("dev")(), nil
		},
	}
}

func PreCheck(t *testing.T) {

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

func PreCheckEnvironment(t *testing.T) {

	PreCheck(t)
	if v := os.Getenv("PINGONE_LICENSE_ID"); v == "" {
		t.Fatal("PINGONE_LICENSE_ID is missing and must be set")
	}
}

func ErrorCheck(t *testing.T) resource.ErrorCheckFunc {
	return func(err error) error {
		if err == nil {
			return nil
		}
		return err
	}
}

func ResourceNameGen() string {
	strlen := 10
	return acctest.RandStringFromCharSet(strlen, acctest.CharSetAlpha)
}

func TestClient(ctx context.Context) (*client.Client, error) {

	config := &client.Config{
		ClientID:      os.Getenv("PINGONE_CLIENT_ID"),
		ClientSecret:  os.Getenv("PINGONE_CLIENT_SECRET"),
		EnvironmentID: os.Getenv("PINGONE_ENVIRONMENT_ID"),
		Region:        os.Getenv("PINGONE_REGION"),
		ForceDelete:   false,
	}

	return config.APIClient(ctx)

}
