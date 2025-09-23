// Copyright © 2025 Ping Identity Corporation

package acctest

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"slices"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
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

var ProtoV6ProviderFactories map[string]func() (tfprotov6.ProviderServer, error) = protoV6ProviderFactoriesInit(context.Background(), "pingone")

func protoV6ProviderFactoriesInit(ctx context.Context, providerNames ...string) map[string]func() (tfprotov6.ProviderServer, error) {
	factories := make(map[string]func() (tfprotov6.ProviderServer, error), len(providerNames))

	for _, name := range providerNames {

		factories[name] = func() (tfprotov6.ProviderServer, error) {
			providerServerFactory, err := provider.ProviderServerFactoryV6(ctx, getProviderTestingVersion())

			if err != nil {
				return nil, err
			}

			return providerServerFactory(), nil
		}
	}

	return factories
}

func getProviderTestingVersion() string {
	returnVar := "dev"
	if v := os.Getenv("PINGONE_TESTING_PROVIDER_VERSION"); v != "" {
		returnVar = v
	}
	return returnVar
}

type TestData struct {
	Invalid string
	Valid   string
}

type MinMaxChecks struct {
	Minimal resource.TestCheckFunc
	Full    resource.TestCheckFunc
}

func PreCheckClient(t *testing.T) {
	if v := os.Getenv("PINGONE_CLIENT_ID"); v == "" {
		t.Fatal("PINGONE_CLIENT_ID is missing and must be set")
	}

	if v := os.Getenv("PINGONE_CLIENT_SECRET"); v == "" {
		t.Fatal("PINGONE_CLIENT_SECRET is missing and must be set")
	}

	if v := os.Getenv("PINGONE_ENVIRONMENT_ID"); v == "" {
		t.Fatal("PINGONE_ENVIRONMENT_ID is missing and must be set")
	}

	if v := os.Getenv("PINGONE_REGION_CODE"); v == "" {
		t.Fatal("PINGONE_REGION_CODE is missing and must be set")
	}
}

func PreCheckNoBeta(t *testing.T) {
	if v := os.Getenv("TESTACC_BETA"); v == "true" {
		t.Skip("Skipping test because TESTACC_BETA is set to true")
	}
}

func PreCheckBeta(t *testing.T) {
	if v := os.Getenv("TESTACC_BETA"); v != "true" {
		t.Skip("Skipping test because TESTACC_BETA is not set to true")
	}
}

func PreCheckNoTestAccFlaky(t *testing.T) {
	if v := os.Getenv("TESTACC_FLAKY"); v == "true" {
		t.Skip("Skipping test because TESTACC_FLAKY is set to true")
	}
}

func PreCheckTestAccFlaky(t *testing.T) {
	if v := os.Getenv("TESTACC_FLAKY"); v != "true" {
		t.Skip("Skipping test because TESTACC_FLAKY is not set to true")
	}
}

func PreCheckOrganisationName(t *testing.T) {
	if v := os.Getenv("PINGONE_ORGANIZATION_NAME"); v == "" {
		t.Fatal("PINGONE_ORGANIZATION_NAME is missing and must be set")
	}
}

func PreCheckOrganisationID(t *testing.T) {
	if v := os.Getenv("PINGONE_ORGANIZATION_ID"); v == "" {
		t.Fatal("PINGONE_ORGANIZATION_ID is missing and must be set")
	}
}

func PreCheckNewEnvironment(t *testing.T) {
	if v := os.Getenv("PINGONE_LICENSE_ID"); v == "" {
		t.Fatal("PINGONE_LICENSE_ID is missing and must be set")
	}
}

func PreCheckDomainVerification(t *testing.T) {

	skipEmailDomainVerified, err := strconv.ParseBool(os.Getenv("PINGONE_EMAIL_DOMAIN_TEST_SKIP"))
	if err != nil {
		skipEmailDomainVerified = false
	}

	if skipEmailDomainVerified {
		t.Skipf("Email domain verified integration tests are skipped")
	}

	if v := os.Getenv("PINGONE_VERIFIED_EMAIL_DOMAIN"); v == "" {
		t.Fatal("PINGONE_VERIFIED_EMAIL_DOMAIN is missing and must be set")
	}
}

func PreCheckSupportsRegion(t *testing.T, supportedRegionCodes []string) {
	if v := os.Getenv("PINGONE_REGION_CODE"); !slices.Contains(supportedRegionCodes, v) {
		t.Skipf("Test not supported in the %s region", v)
	}
}

func PreCheckRegionSupportsWorkforce(t *testing.T) {
	if v := os.Getenv("PINGONE_REGION_CODE"); v == "CA" || v == "SG" {
		t.Skipf("Workforce environment not supported in the Canada or Singapore regions")
	}
}

func PreCheckPKCS12Key(t *testing.T) {
	if v := os.Getenv("PINGONE_KEY_PKCS12"); v == "" {
		t.Fatal("PINGONE_KEY_PKCS12 is missing and must be set")
	}

	if v := os.Getenv("PINGONE_KEY_PKCS12_PASSWORD"); v == "" {
		t.Fatal("PINGONE_KEY_PKCS12_PASSWORD is missing and must be set")
	}
}

func PreCheckAPNSPKCS8Key(t *testing.T) {
	if v := os.Getenv("PINGONE_KEY_PKCS8"); v == "" {
		t.Fatal("PINGONE_KEY_PKCS8 is missing and must be set")
	}
}

func PreCheckPKCS12UnencryptedKey(t *testing.T) {
	if v := os.Getenv("PINGONE_KEY_PKCS12_UNENCRYPTED"); v == "" {
		t.Fatal("PINGONE_KEY_PKCS12_UNENCRYPTED is missing and must be set")
	}
}

func PreCheckPKCS12WithCSR(t *testing.T) {
	if v := os.Getenv("PINGONE_KEY_PKCS10_CSR"); v == "" {
		t.Fatal("PINGONE_KEY_PKCS10_CSR is missing and must be set")
	}

	if v := os.Getenv("PINGONE_KEY_PEM_CSR"); v == "" {
		t.Fatal("PINGONE_KEY_PEM_CSR is missing and must be set")
	}
}

func PreCheckPKCS12CSRResponse(t *testing.T) {
	if v := os.Getenv("PINGONE_KEY_PEM_CSR_RESPONSE"); v == "" {
		t.Fatal("PINGONE_KEY_PEM_CSR_RESPONSE is missing and must be set")
	}
}

func PreCheckPKCS7Cert(t *testing.T) {
	if v := os.Getenv("PINGONE_KEY_PKCS7_CERT"); v == "" {
		t.Fatal("PINGONE_KEY_PKCS7_CERT is missing and must be set")
	}
}

func PreCheckPEMCert(t *testing.T) {
	if v := os.Getenv("PINGONE_KEY_PEM_CERT"); v == "" {
		t.Fatal("PINGONE_KEY_PEM_CERT is missing and must be set")
	}
}

func PreCheckGoogleJSONKey(t *testing.T) {
	if v := os.Getenv("PINGONE_GOOGLE_JSON_KEY"); v == "" {
		t.Fatal("PINGONE_GOOGLE_JSON_KEY is missing and must be set")
	}
}

func PreCheckGoogleFirebaseCredentials(t *testing.T) {
	if v := os.Getenv("PINGONE_GOOGLE_FIREBASE_CREDENTIALS"); v == "" {
		t.Fatal("PINGONE_GOOGLE_FIREBASE_CREDENTIALS is missing and must be set")
	}
}

func PreCheckCustomDomainSSL(t *testing.T) {
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

func PreCheckTwilio(t *testing.T) {

	skipTwilio, err := strconv.ParseBool(os.Getenv("PINGONE_TWILIO_TEST_SKIP"))
	if err != nil {
		skipTwilio = false
	}

	if skipTwilio {
		t.Skipf("Twilio integration tests are skipped")
	}

	if v := os.Getenv("PINGONE_TWILIO_SID"); v == "" {
		t.Fatal("PINGONE_TWILIO_SID is missing and must be set")
	}

	if v := os.Getenv("PINGONE_TWILIO_AUTH_TOKEN"); v == "" {
		t.Fatal("PINGONE_TWILIO_AUTH_TOKEN is missing and must be set")
	}

	if v := os.Getenv("PINGONE_TWILIO_NUMBER"); v == "" {
		t.Fatal("PINGONE_TWILIO_NUMBER is missing and must be set")
	}
}

func PreCheckSyniverse(t *testing.T) {

	skipSyniverse, err := strconv.ParseBool(os.Getenv("PINGONE_SYNIVERSE_TEST_SKIP"))
	if err != nil {
		skipSyniverse = false
	}

	if skipSyniverse {
		t.Skipf("Syniverse integration tests are skipped")
	}

	if v := os.Getenv("PINGONE_SYNIVERSE_AUTH_TOKEN"); v == "" {
		t.Fatal("PINGONE_SYNIVERSE_AUTH_TOKEN is missing and must be set")
	}

	if v := os.Getenv("PINGONE_SYNIVERSE_NUMBER"); v == "" {
		t.Fatal("PINGONE_SYNIVERSE_NUMBER is missing and must be set")
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

func TestClient(ctx context.Context) (*client.Client, error) {

	regionCode := management.EnumRegionCode(os.Getenv("PINGONE_REGION_CODE"))

	config := &client.Config{
		ClientID:      os.Getenv("PINGONE_CLIENT_ID"),
		ClientSecret:  os.Getenv("PINGONE_CLIENT_SECRET"),
		EnvironmentID: os.Getenv("PINGONE_ENVIRONMENT_ID"),
		RegionCode:    &regionCode,
		GlobalOptions: &client.GlobalOptions{
			Population: &client.PopulationOptions{
				ContainsUsersForceDelete: false,
			},
		},
	}

	return config.APIClient(ctx, getProviderTestingVersion())

}

func PreCheckTestClient(ctx context.Context, t *testing.T) *client.Client {
	p1Client, err := TestClient(ctx)

	if err != nil {
		t.Fatalf("Failed to get API client: %v", err)
	}

	return p1Client
}

func MinimalSandboxEnvironment(resourceName, licenseID string) string {
	return fmt.Sprintf(`
	%[1]s
		
	resource "pingone_population_default" "%[2]s" {
		environment_id = pingone_environment.%[2]s.id

		name = "%[2]s"
	}
`, MinimalSandboxEnvironmentNoPopulation(resourceName, licenseID), resourceName)
}

func MinimalSandboxEnvironmentNoPopulation(resourceName, licenseID string) string {
	return MinimalEnvironmentNoPopulation(resourceName, licenseID, management.ENUMENVIRONMENTTYPE_SANDBOX)
}

func MinimalEnvironmentNoPopulation(resourceName, licenseID string, environmentType management.EnumEnvironmentType) string {
	return fmt.Sprintf(`
resource "pingone_environment" "%[1]s" {
	name = "%[1]s"
	license_id = "%[2]s"
	type = "%[3]s"

	services = [
		{
			type = "SSO"
		},
		{
			type = "DaVinci"
			tags = ["DAVINCI_MINIMAL"]
		},
		{
			type = "MFA"
		},
		{
			type = "Risk"
		},
		{
			type = "Credentials"
		},
		{
			type = "Verify"
		}
	]
}
`, resourceName, licenseID, string(environmentType))
}

func GenericSandboxEnvironment() string {
	return `
		data "pingone_environment" "general_test" {
			name = "tf-testacc-dynamic-general-test"
		}`
}

const (
	WorkforceV1SandboxEnvironmentName = "tf-testacc-static-workforce-test"
	WorkforceV2SandboxEnvironmentName = "tf-testacc-static-workforce-v2-test"
)

// Static environment that uses v1 PingID
func WorkforceV1SandboxEnvironment() string {
	return fmt.Sprintf(`
		data "pingone_environment" "workforce_test" {
			name = "%s"
		}`, WorkforceV1SandboxEnvironmentName)
}

// Static environment that uses v2 PingID
func WorkforceV2SandboxEnvironment() string {
	return fmt.Sprintf(`
		data "pingone_environment" "workforce_test" {
			name = "%s"
		}`, WorkforceV2SandboxEnvironmentName)
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

func CheckParentEnvironmentDestroy(ctx context.Context, apiClient *management.APIClient, environmentID string) (bool, error) {
	environment, r, err := apiClient.EnvironmentsApi.ReadOneEnvironment(ctx, environmentID).Execute()

	destroyed, err := CheckForResourceDestroy(r, err)
	if err != nil {
		return destroyed, err
	}

	if destroyed {
		return destroyed, nil
	} else {
		if environment != nil && environment.Type == management.ENUMENVIRONMENTTYPE_PRODUCTION {
			return true, nil
		} else {
			return false, nil
		}
	}
}

func CheckParentUserDestroy(ctx context.Context, apiClient *management.APIClient, environmentID, userID string) (bool, error) {
	_, r, err := apiClient.UsersApi.ReadUser(ctx, environmentID, userID).Execute()

	return CheckForResourceDestroy(r, err)
}

func CheckForResourceDestroy(r *http.Response, err error) (bool, error) {
	defaultDestroyHttpCode := 404
	return CheckForResourceDestroyCustomHTTPCode(r, err, defaultDestroyHttpCode)
}

func CheckForResourceDestroyCustomHTTPCode(r *http.Response, err error, customHttpCode int) (bool, error) {
	if err != nil {

		if r == nil {
			return false, fmt.Errorf("Response object does not exist and no error detected")
		}

		if r.StatusCode == customHttpCode {
			return true, nil
		}

		return false, err
	}

	return false, nil
}

// AlterStringCasing alternates the case of alphabetic characters in a string for testing purposes.
// It returns a string where even-indexed characters (0, 2, 4, etc.) are converted to uppercase
// and odd-indexed characters (1, 3, 5, etc.) are converted to lowercase. Non-alphabetic characters
// remain unchanged in their original positions.
// The strInput parameter must be a valid string that may contain any Unicode characters.
// This function is primarily used in acceptance tests to create case-insensitive string comparisons
// and verify that data source filters work correctly regardless of character casing.
// No external dependencies or environment variables are required for this function to operate.
func AlterStringCasing(strInput string) string {
	runes := []rune(strInput)
	for i := range runes {
		if i%2 == 0 {
			if runes[i] >= 'a' && runes[i] <= 'z' {
				runes[i] = runes[i] - ('a' - 'A')
			}
		} else {
			if runes[i] >= 'A' && runes[i] <= 'Z' {
				runes[i] = runes[i] + ('a' - 'A')
			}
		}
	}
	return string(runes)
}
