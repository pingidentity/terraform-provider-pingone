package acctest

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/provider"
	"github.com/pingidentity/terraform-provider-pingone/internal/provider/sdkv2"
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

var ProtoV5ProviderFactories map[string]func() (tfprotov5.ProviderServer, error) = protoV5ProviderFactoriesInit(context.Background(), "pingone")

func init() {
	Provider = sdkv2.New("dev")()

	// Always allocate a new provider instance each invocation, otherwise gRPC
	// ProviderConfigure() can overwrite configuration during concurrent testing.
	ProviderFactories = map[string]func() (*schema.Provider, error){
		"pingone": func() (*schema.Provider, error) {
			provider := sdkv2.New("acctest")()

			if provider == nil {
				return nil, fmt.Errorf("Cannot initiate provider factory")
			}
			return provider, nil
		},
	}
}

func protoV5ProviderFactoriesInit(ctx context.Context, providerNames ...string) map[string]func() (tfprotov5.ProviderServer, error) {
	factories := make(map[string]func() (tfprotov5.ProviderServer, error), len(providerNames))

	for _, name := range providerNames {

		factories[name] = func() (tfprotov5.ProviderServer, error) {
			providerServerFactory, _, err := provider.ProviderServerFactoryV5(ctx, "acctest")

			if err != nil {
				return nil, err
			}

			return providerServerFactory(), nil
		}
	}

	return factories
}

type TestData struct {
	Invalid string
	Valid   string
}

type MinMaxChecks struct {
	Minimal resource.TestCheckFunc
	Full    resource.TestCheckFunc
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

func PreCheckEnvironmentDomainVerified(t *testing.T) {

	PreCheckEnvironment(t)
	if v := os.Getenv("PINGONE_VERIFIED_EMAIL_DOMAIN"); v == "" {
		t.Fatal("PINGONE_VERIFIED_EMAIL_DOMAIN is missing and must be set")
	}
}

func PreCheckWorkforceEnvironment(t *testing.T) {

	PreCheckEnvironment(t)
	if v := os.Getenv("PINGONE_REGION"); v == "Canada" {
		t.Skipf("Workforce environment not supported in the Canada region")
	}
}

func PreCheckEnvironmentAndPKCS12(t *testing.T) {

	PreCheckEnvironment(t)
	if v := os.Getenv("PINGONE_KEY_PKCS12"); v == "" {
		t.Fatal("PINGONE_KEY_PKCS12 is missing and must be set")
	}
}

func PreCheckEnvironmentAndPKCS12WithCSR(t *testing.T) {

	PreCheckEnvironmentAndPKCS12(t)
	if v := os.Getenv("PINGONE_KEY_PKCS10_CSR"); v == "" {
		t.Fatal("PINGONE_KEY_PKCS10_CSR is missing and must be set")
	}

	if v := os.Getenv("PINGONE_KEY_PEM_CSR"); v == "" {
		t.Fatal("PINGONE_KEY_PEM_CSR is missing and must be set")
	}
}

func PreCheckEnvironmentAndPKCS12WithCSRResponse(t *testing.T) {

	PreCheckEnvironmentAndPKCS12(t)
	if v := os.Getenv("PINGONE_KEY_PEM_CSR_RESPONSE"); v == "" {
		t.Fatal("PINGONE_KEY_PEM_CSR_RESPONSE is missing and must be set")
	}
}

func PreCheckEnvironmentAndPKCS12WithCerts(t *testing.T) {

	PreCheckEnvironmentAndPKCS12(t)
	if v := os.Getenv("PINGONE_KEY_PKCS7_CERT"); v == "" {
		t.Fatal("PINGONE_KEY_PKCS7_CERT is missing and must be set")
	}

	if v := os.Getenv("PINGONE_KEY_PEM_CERT"); v == "" {
		t.Fatal("PINGONE_KEY_PEM_CERT is missing and must be set")
	}
}

func PreCheckEnvironmentAndPKCS7(t *testing.T) {

	PreCheckEnvironment(t)
	if v := os.Getenv("PINGONE_KEY_PKCS7_CERT"); v == "" {
		t.Fatal("PINGONE_KEY_PKCS7_CERT is missing and must be set")
	}
}

func PreCheckEnvironmentAndPEM(t *testing.T) {

	PreCheckEnvironment(t)
	if v := os.Getenv("PINGONE_KEY_PEM_CERT"); v == "" {
		t.Fatal("PINGONE_KEY_PEM_CERT is missing and must be set")
	}
}

func PreCheckEnvironmentAndCustomDomainSSL(t *testing.T) {

	PreCheckEnvironment(t)
	if v := os.Getenv("PINGONE_DOMAIN_CERTIFICATE_PEM"); v == "" {
		t.Fatal("PINGONE_DOMAIN_CERTIFICATE_PEM is missing and must be set")
	}

	if v := os.Getenv("PINGONE_DOMAIN_INTERMEDIATE_CERTIFICATE_PEM"); v == "" {
		t.Fatal("PINGONE_DOMAIN_INTERMEDIATE_CERTIFICATE_PEM is missing and must be set")
	}

	if v := os.Getenv("PINGONE_DOMAIN_KEY_PEM"); v == "" {
		t.Fatal("PINGONE_DOMAIN_KEY_PEM is missing and must be set")
	}
}

func PreCheckEnvironmentAndOrganisation(t *testing.T) {

	PreCheckEnvironment(t)
	if v := os.Getenv("PINGONE_ORGANIZATION_ID"); v == "" {
		t.Fatal("PINGONE_ORGANIZATION_ID is missing and must be set")
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

func ResourceNameGenEnvironment() string {
	return fmt.Sprintf("tf-testacc-dynamic-%s", ResourceNameGen())
}

func ResourceNameGenDefaultPopulation() string {
	return fmt.Sprintf("default-%s", ResourceNameGen())
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

func TestAccCheckEnvironmentDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_environment" {
			continue
		}

		_, r, err := apiClient.EnvironmentsApi.ReadOneEnvironment(ctx, rs.Primary.ID).Execute()

		if r.StatusCode == 404 {
			continue
		}

		if err != nil {
			return err
		}

		return fmt.Errorf("PingOne Environment Instance %s still exists", rs.Primary.ID)
	}

	return nil
}

func MinimalSandboxEnvironment(resourceName, licenseID string) string {
	return fmt.Sprintf(`
		resource "pingone_environment" "%[1]s" {
			name = "%[1]s"
			type = "SANDBOX"
			license_id = "%[2]s"
			default_population {
			}
			service {
				type = "SSO"
			}
			service {
				type = "MFA"
			}
		}`, resourceName, licenseID)
}

func GenericSandboxEnvironment() string {
	return `
		data "pingone_environment" "general_test" {
			name = "tf-testacc-dynamic-general-test"
		}`
}

func WorkforceSandboxEnvironment() string {
	return `
		data "pingone_environment" "workforce_test" {
			name = "tf-testacc-static-workforce-test"
		}`
}

func DomainVerifiedSandboxEnvironment() string {
	return `
		data "pingone_environment" "domainverified_test" {
			name = "tf-testacc-static-domainverified-test"
		}`
}

func AgreementSandboxEnvironment() string {
	return `
		data "pingone_environment" "agreement_test" {
			name = "tf-testacc-static-agreements-test"
		}`
}
