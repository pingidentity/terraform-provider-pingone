// Copyright © 2025 Ping Identity Corporation

// Package acctest provides common functions and utilities for acceptance testing in the PingOne Terraform provider.
package acctest

import (
	"context"
	"fmt"
	"net/http"
	"os"
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

var (
	// ProviderFactories is a static map containing only the main provider instance.
	// Use other ProviderFactories functions, such as FactoriesAlternate,
	// for tests requiring special provider configurations.
	ProviderFactories map[string]func() (*schema.Provider, error)

	// Provider is the "main" provider instance.
	// This Provider can be used in testing code for API calls without requiring
	// the use of saving and referencing specific ProviderFactories instances.
	// PreCheck(t) must be called before using this provider instance.
	Provider *schema.Provider

	// ProtoV6ProviderFactories are used to instantiate a provider during acceptance testing.
	// The factory function will be invoked for every Terraform CLI command executed
	// to create a provider server to which the CLI can reattach.
	ProtoV6ProviderFactories map[string]func() (tfprotov6.ProviderServer, error) = protoV6ProviderFactoriesInit(context.Background(), "pingone")
)

// protoV6ProviderFactoriesInit initializes provider factories for the given provider names.
// It returns a map where each key is a provider name string and each value is a factory function
// that creates tfprotov6.ProviderServer instances for acceptance testing with Protocol 6.
// The ctx parameter provides the context for provider initialization and configuration.
// The providerNames parameter is a variadic string slice specifying which providers to create factories for.
// This function creates isolated provider instances for each test execution to ensure test independence.
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

// getProviderTestingVersion returns the provider version string to use for testing.
// It returns a string value representing the provider version, defaulting to "dev" for development builds.
// The returned version can be overridden by setting the PINGONE_TESTING_PROVIDER_VERSION environment variable
// to specify a custom version string for testing purposes.
func getProviderTestingVersion() string {
	returnVar := "dev"
	if v := os.Getenv("PINGONE_TESTING_PROVIDER_VERSION"); v != "" {
		returnVar = v
	}
	return returnVar
}

// TestData contains test data values for validation testing.
type TestData struct {
	Invalid string
	Valid   string
}

// MinMaxChecks contains test check functions for minimal and full configurations.
type MinMaxChecks struct {
	Minimal resource.TestCheckFunc
	Full    resource.TestCheckFunc
}

// EnumFeatureFlag represents feature flags for testing specific functionality.
type EnumFeatureFlag string

const (
	// ENUMFEATUREFLAG_DAVINCI represents the DaVinci feature flag for testing.
	ENUMFEATUREFLAG_DAVINCI EnumFeatureFlag = "DAVINCI"
)

// PreCheckClient performs pre-test validation checks for required client environment variables.
// This function validates that all necessary PingOne client configuration environment variables are set
// before running acceptance tests. It terminates the test with a fatal error if any required variables are missing.
// The t parameter is the testing instance used to report fatal errors.
// Required environment variables: PINGONE_CLIENT_ID, PINGONE_CLIENT_SECRET, PINGONE_ENVIRONMENT_ID, PINGONE_REGION_CODE.
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

// PreCheckNoFeatureFlag performs pre-test validation when no specific feature flag is required.
// This function ensures that tests run when no feature flag requirements are needed.
// It delegates to PreCheckFeatureFlag with an empty string to indicate no flag requirement.
// The t parameter is the testing instance used for test management and potential skipping.
func PreCheckNoFeatureFlag(t *testing.T) {
	PreCheckFeatureFlag(t, "")
}

// PreCheckFeatureFlag performs pre-test validation for specific feature flag requirements.
// This function skips tests when the required feature flag is not enabled via the FEATURE_FLAG environment variable.
// The t parameter is the testing instance used for skipping tests when conditions are not met.
// The flag parameter specifies the EnumFeatureFlag value that must be set in the FEATURE_FLAG environment variable.
// If the environment variable does not match the required flag value, the test is skipped with an informative message.
func PreCheckFeatureFlag(t *testing.T, flag EnumFeatureFlag) {
	if v := os.Getenv("FEATURE_FLAG"); v != string(flag) {
		t.Skipf("Skipping feature flag test.  Flag required: \"%s\"", string(flag))
	}
}

// PreCheckNoTestAccFlaky skips tests when flaky test execution is explicitly enabled.
// This function prevents execution of tests that are not marked as flaky when the TESTACC_FLAKY
// environment variable is set to "true", allowing for selective execution of stable tests only.
// The t parameter is the testing instance used for skipping tests based on the flaky test configuration.
func PreCheckNoTestAccFlaky(t *testing.T) {
	if v := os.Getenv("TESTACC_FLAKY"); v == "true" {
		t.Skip("Skipping test because TESTACC_FLAKY is set to true")
	}
}

// PreCheckTestAccFlaky skips tests when flaky test execution is not explicitly enabled.
// This function ensures that tests marked as flaky only run when the TESTACC_FLAKY environment
// variable is specifically set to "true", preventing unreliable tests from running in normal test suites.
// The t parameter is the testing instance used for skipping tests based on the flaky test configuration.
func PreCheckTestAccFlaky(t *testing.T) {
	if v := os.Getenv("TESTACC_FLAKY"); v != "true" {
		t.Skip("Skipping test because TESTACC_FLAKY is not set to true")
	}
}

// PreCheckOrganisationName performs pre-test validation for the required organization name environment variable.
// This function ensures that the PINGONE_ORGANIZATION_NAME environment variable is set before running tests
// that require organization name identification. It terminates the test with a fatal error if the variable is missing.
// The t parameter is the testing instance used to report fatal errors when the organization name is not configured.
func PreCheckOrganisationName(t *testing.T) {
	if v := os.Getenv("PINGONE_ORGANIZATION_NAME"); v == "" {
		t.Fatal("PINGONE_ORGANIZATION_NAME is missing and must be set")
	}
}

// PreCheckOrganisationID performs pre-test validation for the required organization ID environment variable.
// This function ensures that the PINGONE_ORGANIZATION_ID environment variable is set before running tests
// that require organization ID identification. It terminates the test with a fatal error if the variable is missing.
// The t parameter is the testing instance used to report fatal errors when the organization ID is not configured.
func PreCheckOrganisationID(t *testing.T) {
	if v := os.Getenv("PINGONE_ORGANIZATION_ID"); v == "" {
		t.Fatal("PINGONE_ORGANIZATION_ID is missing and must be set")
	}
}

// PreCheckNewEnvironment performs pre-test validation for creating new PingOne environments.
// This function ensures that the PINGONE_LICENSE_ID environment variable is set before running tests
// that need to create new environments, as a valid license ID is required for environment provisioning.
// The t parameter is the testing instance used to report fatal errors when the license ID is not configured.
func PreCheckNewEnvironment(t *testing.T) {
	if v := os.Getenv("PINGONE_LICENSE_ID"); v == "" {
		t.Fatal("PINGONE_LICENSE_ID is missing and must be set")
	}
}

// PreCheckDomainVerification performs pre-test validation for domain verification testing requirements.
// This function checks if domain verification tests should be skipped based on the PINGONE_EMAIL_DOMAIN_TEST_SKIP
// environment variable and ensures the PINGONE_VERIFIED_EMAIL_DOMAIN variable is set when tests should run.
// The t parameter is the testing instance used for skipping tests or reporting fatal errors.
// The function skips tests if PINGONE_EMAIL_DOMAIN_TEST_SKIP is set to true, otherwise requires PINGONE_VERIFIED_EMAIL_DOMAIN to be configured.
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

// PreCheckRegionSupportsWorkforce skips tests for PingOne regions that don't support workforce environments.
// This function prevents execution of workforce-related tests in regions where this functionality is not available,
// specifically skipping tests when PINGONE_REGION_CODE is set to "CA" (Canada) or "SG" (Singapore).
// The t parameter is the testing instance used for skipping tests with an informative message about regional limitations.
func PreCheckRegionSupportsWorkforce(t *testing.T) {
	if v := os.Getenv("PINGONE_REGION_CODE"); v == "CA" || v == "SG" {
		t.Skipf("Workforce environment not supported in the Canada or Singapore regions")
	}
}

// PreCheckPKCS12Key performs pre-test validation for PKCS12 key file environment variables.
// This function ensures that both PINGONE_KEY_PKCS12 and PINGONE_KEY_PKCS12_PASSWORD environment variables
// are set before running tests that require PKCS12 certificate handling. It terminates the test with a fatal error if either variable is missing.
// The t parameter is the testing instance used to report fatal errors when PKCS12 configuration is incomplete.
func PreCheckPKCS12Key(t *testing.T) {
	if v := os.Getenv("PINGONE_KEY_PKCS12"); v == "" {
		t.Fatal("PINGONE_KEY_PKCS12 is missing and must be set")
	}

	if v := os.Getenv("PINGONE_KEY_PKCS12_PASSWORD"); v == "" {
		t.Fatal("PINGONE_KEY_PKCS12_PASSWORD is missing and must be set")
	}
}

// PreCheckAPNSPKCS8Key performs pre-test validation for Apple Push Notification Service PKCS8 key environment variables.
// This function ensures that the PINGONE_KEY_PKCS8 environment variable is set before running tests
// that require APNS certificate handling with PKCS8 key format. It terminates the test with a fatal error if the variable is missing.
// The t parameter is the testing instance used to report fatal errors when APNS PKCS8 configuration is incomplete.
func PreCheckAPNSPKCS8Key(t *testing.T) {
	if v := os.Getenv("PINGONE_KEY_PKCS8"); v == "" {
		t.Fatal("PINGONE_KEY_PKCS8 is missing and must be set")
	}
}

// PreCheckPKCS12UnencryptedKey performs pre-test validation for unencrypted PKCS12 key environment variables.
// This function ensures that the PINGONE_KEY_PKCS12_UNENCRYPTED environment variable is set before running tests
// that require unencrypted PKCS12 certificate handling. It terminates the test with a fatal error if the variable is missing.
// The t parameter is the testing instance used to report fatal errors when unencrypted PKCS12 configuration is incomplete.
func PreCheckPKCS12UnencryptedKey(t *testing.T) {
	if v := os.Getenv("PINGONE_KEY_PKCS12_UNENCRYPTED"); v == "" {
		t.Fatal("PINGONE_KEY_PKCS12_UNENCRYPTED is missing and must be set")
	}
}

// PreCheckPKCS12WithCSR performs pre-test validation for PKCS12 certificate signing request environment variables.
// This function ensures that both PINGONE_KEY_PKCS10_CSR and PINGONE_KEY_PEM_CSR environment variables are set
// before running tests that require certificate signing request handling with PKCS12 format. It terminates the test with a fatal error if either variable is missing.
// The t parameter is the testing instance used to report fatal errors when CSR configuration is incomplete.
func PreCheckPKCS12WithCSR(t *testing.T) {
	if v := os.Getenv("PINGONE_KEY_PKCS10_CSR"); v == "" {
		t.Fatal("PINGONE_KEY_PKCS10_CSR is missing and must be set")
	}

	if v := os.Getenv("PINGONE_KEY_PEM_CSR"); v == "" {
		t.Fatal("PINGONE_KEY_PEM_CSR is missing and must be set")
	}
}

// PreCheckPKCS12CSRResponse performs pre-test validation for PKCS12 certificate signing request response environment variables.
// This function ensures that the PINGONE_KEY_PEM_CSR_RESPONSE environment variable is set before running tests
// that require processing of certificate signing request responses. It terminates the test with a fatal error if the variable is missing.
// The t parameter is the testing instance used to report fatal errors when CSR response configuration is incomplete.
func PreCheckPKCS12CSRResponse(t *testing.T) {
	if v := os.Getenv("PINGONE_KEY_PEM_CSR_RESPONSE"); v == "" {
		t.Fatal("PINGONE_KEY_PEM_CSR_RESPONSE is missing and must be set")
	}
}

// PreCheckPKCS7Cert performs pre-test validation for PKCS7 certificate environment variables.
// This function ensures that the PINGONE_KEY_PKCS7_CERT environment variable is set before running tests
// that require PKCS7 certificate handling. It terminates the test with a fatal error if the variable is missing.
// The t parameter is the testing instance used to report fatal errors when PKCS7 certificate configuration is incomplete.
func PreCheckPKCS7Cert(t *testing.T) {
	if v := os.Getenv("PINGONE_KEY_PKCS7_CERT"); v == "" {
		t.Fatal("PINGONE_KEY_PKCS7_CERT is missing and must be set")
	}
}

// PreCheckPEMCert performs pre-test validation for PEM certificate environment variables.
// This function ensures that the PINGONE_KEY_PEM_CERT environment variable is set before running tests
// that require PEM certificate handling. It terminates the test with a fatal error if the variable is missing.
// The t parameter is the testing instance used to report fatal errors when PEM certificate configuration is incomplete.
func PreCheckPEMCert(t *testing.T) {
	if v := os.Getenv("PINGONE_KEY_PEM_CERT"); v == "" {
		t.Fatal("PINGONE_KEY_PEM_CERT is missing and must be set")
	}
}

// PreCheckGoogleJSONKey performs pre-test validation for Google JSON key environment variables.
// This function ensures that the PINGONE_GOOGLE_JSON_KEY environment variable is set before running tests
// that require Google service account JSON key integration. It terminates the test with a fatal error if the variable is missing.
// The t parameter is the testing instance used to report fatal errors when Google JSON key configuration is incomplete.
func PreCheckGoogleJSONKey(t *testing.T) {
	if v := os.Getenv("PINGONE_GOOGLE_JSON_KEY"); v == "" {
		t.Fatal("PINGONE_GOOGLE_JSON_KEY is missing and must be set")
	}
}

// PreCheckGoogleFirebaseCredentials performs pre-test validation for Google Firebase credentials environment variables.
// This function ensures that the PINGONE_GOOGLE_FIREBASE_CREDENTIALS environment variable is set before running tests
// that require Google Firebase integration. It terminates the test with a fatal error if the variable is missing.
// The t parameter is the testing instance used to report fatal errors when Firebase credentials configuration is incomplete.
func PreCheckGoogleFirebaseCredentials(t *testing.T) {
	if v := os.Getenv("PINGONE_GOOGLE_FIREBASE_CREDENTIALS"); v == "" {
		t.Fatal("PINGONE_GOOGLE_FIREBASE_CREDENTIALS is missing and must be set")
	}
}

// PreCheckCustomDomainSSL performs pre-test validation for custom domain SSL certificate environment variables.
// This function ensures that all required SSL certificate environment variables are set before running tests
// that require custom domain SSL configuration: PINGONE_DOMAIN_CERTIFICATE_PEM, PINGONE_DOMAIN_INTERMEDIATE_CERTIFICATE_PEM, and PINGONE_DOMAIN_KEY_PEM.
// The t parameter is the testing instance used to report fatal errors when SSL certificate configuration is incomplete.
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

// PreCheckTwilio performs pre-test validation for Twilio integration environment variables.
// This function checks if Twilio integration tests should be skipped based on the PINGONE_TWILIO_TEST_SKIP
// environment variable and ensures required Twilio configuration variables are set when tests should run.
// The t parameter is the testing instance used for skipping tests or reporting fatal errors.
// Required variables when not skipped: PINGONE_TWILIO_SID, PINGONE_TWILIO_AUTH_TOKEN, PINGONE_TWILIO_NUMBER.
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

// PreCheckSyniverse performs pre-test validation for Syniverse integration environment variables.
// This function checks if Syniverse integration tests should be skipped based on the PINGONE_SYNIVERSE_TEST_SKIP
// environment variable and ensures required Syniverse configuration variables are set when tests should run.
// The t parameter is the testing instance used for skipping tests or reporting fatal errors.
// Required variables when not skipped: PINGONE_SYNIVERSE_AUTH_TOKEN, PINGONE_SYNIVERSE_NUMBER.
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

// ErrorCheck returns a function that can be used as an error check callback in acceptance tests.
// This function creates a resource.ErrorCheckFunc that simply passes through any errors unchanged,
// providing a standardized way to handle errors during test execution without additional processing.
// The t parameter is the testing instance, though it is not used in the current implementation.
// The returned function maintains the original error or returns nil if no error occurred.
func ErrorCheck(t *testing.T) resource.ErrorCheckFunc {
	return func(err error) error {
		if err == nil {
			return nil
		}
		return err
	}
}

// ResourceNameGen generates a random resource name string for use in acceptance testing.
// This function creates a 10-character random string using only alphabetic characters,
// suitable for naming test resources to ensure uniqueness across test runs.
// The returned string contains only letters (a-z, A-Z) and is appropriate for resource naming conventions.
func ResourceNameGen() string {
	strlen := 10
	return acctest.RandStringFromCharSet(strlen, acctest.CharSetAlpha)
}

// ResourceNameGenEnvironment generates a random environment name string for use in acceptance testing.
// This function creates a unique environment name by combining the prefix "tf-testacc-dynamic-"
// with a randomly generated string, ensuring test environments have predictable naming patterns while remaining unique.
// The returned string follows the format "tf-testacc-dynamic-{random}" where {random} is a 10-character alphabetic string.
func ResourceNameGenEnvironment() string {
	return fmt.Sprintf("tf-testacc-dynamic-%s", ResourceNameGen())
}

// TestClient creates and configures a PingOne API client for use in acceptance testing.
// This function initializes a client using environment variables for authentication and configuration,
// setting up the necessary region, credentials, and global options for test operations.
// The ctx parameter provides the context for client initialization and API operations.
// Returns a configured client.Client instance or an error if client creation fails.
// The client is configured with test-appropriate settings including force delete options for populations containing users.
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

// PreCheckTestClient creates a configured test client and performs pre-test validation.
// This function combines client creation with error handling, ensuring that a valid API client
// is available before test execution begins. It terminates the test with a fatal error if client creation fails.
// The ctx parameter provides the context for client initialization.
// The t parameter is the testing instance used to report fatal errors.
// Returns a configured client.Client instance ready for use in test operations.
func PreCheckTestClient(ctx context.Context, t *testing.T) *client.Client {
	p1Client, err := TestClient(ctx)

	if err != nil {
		t.Fatalf("Failed to get API client: %v", err)
	}

	return p1Client
}

// MinimalSandboxEnvironment returns a Terraform configuration string for a minimal sandbox environment with default population.
// This function creates a complete environment configuration that includes both the environment resource
// and a default population resource, suitable for tests requiring a fully configured sandbox environment.
// The resourceName parameter specifies the name to use for the Terraform resources.
// The licenseID parameter specifies the PingOne license ID required for environment creation.
// Returns a formatted Terraform configuration string that can be used in acceptance tests.
func MinimalSandboxEnvironment(resourceName, licenseID string) string {
	return fmt.Sprintf(`
	%[1]s
		
	resource "pingone_population_default" "%[2]s" {
		environment_id = pingone_environment.%[2]s.id

		name = "%[2]s"
	}
`, MinimalSandboxEnvironmentNoPopulation(resourceName, licenseID), resourceName)
}

// MinimalSandboxEnvironmentNoPopulation returns a Terraform configuration string for a minimal sandbox environment without default population.
// This function creates an environment configuration without any population resources,
// suitable for tests that need to manage populations separately or don't require them.
// The resourceName parameter specifies the name to use for the Terraform environment resource.
// The licenseID parameter specifies the PingOne license ID required for environment creation.
// Returns a formatted Terraform configuration string for a sandbox environment.
func MinimalSandboxEnvironmentNoPopulation(resourceName, licenseID string) string {
	return MinimalEnvironmentNoPopulation(resourceName, licenseID, management.ENUMENVIRONMENTTYPE_SANDBOX)
}

// MinimalEnvironmentNoPopulation returns a Terraform configuration string for a minimal environment of the specified type without default population.
// This function creates a flexible environment configuration that can be used for any environment type,
// with all standard PingOne services enabled but no population resources included.
// The resourceName parameter specifies the name to use for the Terraform environment resource.
// The licenseID parameter specifies the PingOne license ID required for environment creation.
// The environmentType parameter specifies the type of environment to create (e.g., SANDBOX, PRODUCTION).
// Returns a formatted Terraform configuration string with SSO, MFA, Risk, Credentials, and Verify services enabled.
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

// GenericSandboxEnvironment returns a Terraform configuration string for referencing a generic sandbox environment.
// This function provides a data source configuration that references a pre-existing test environment
// named "tf-testacc-dynamic-general-test", suitable for tests that need a stable environment reference.
// Returns a formatted Terraform data source configuration string for the general test environment.
func GenericSandboxEnvironment() string {
	return `
		data "pingone_environment" "general_test" {
			name = "tf-testacc-dynamic-general-test"
		}`
}

const (
	// WorkforceV1SandboxEnvironmentName is the name of the static workforce v1 test environment.
	WorkforceV1SandboxEnvironmentName = "tf-testacc-static-workforce-test"
	// WorkforceV2SandboxEnvironmentName is the name of the static workforce v2 test environment.
	WorkforceV2SandboxEnvironmentName = "tf-testacc-static-workforce-v2-test"
)

// WorkforceV1SandboxEnvironment returns a Terraform configuration string for referencing a static workforce v1 environment.
// This function provides a data source configuration that references a pre-existing workforce test environment
// configured with PingID v1 integration, suitable for tests requiring workforce functionality with legacy PingID.
// Returns a formatted Terraform data source configuration string for the workforce v1 test environment.
func WorkforceV1SandboxEnvironment() string {
	return fmt.Sprintf(`
		data "pingone_environment" "workforce_test" {
			name = "%s"
		}`, WorkforceV1SandboxEnvironmentName)
}

// WorkforceV2SandboxEnvironment returns a Terraform configuration string for referencing a static workforce v2 environment.
// This function provides a data source configuration that references a pre-existing workforce test environment
// configured with PingID v2 integration, suitable for tests requiring workforce functionality with modern PingID.
// Returns a formatted Terraform data source configuration string for the workforce v2 test environment.
func WorkforceV2SandboxEnvironment() string {
	return fmt.Sprintf(`
		data "pingone_environment" "workforce_test" {
			name = "%s"
		}`, WorkforceV2SandboxEnvironmentName)
}

// DomainVerifiedSandboxEnvironment returns a Terraform configuration string for referencing a domain-verified sandbox environment.
// This function provides a data source configuration that references a pre-existing test environment
// with verified domain configuration, suitable for tests requiring domain verification functionality.
// Returns a formatted Terraform data source configuration string for the domain-verified test environment.
func DomainVerifiedSandboxEnvironment() string {
	return `
		data "pingone_environment" "domainverified_test" {
			name = "tf-testacc-static-domainverified-test"
		}`
}

// AgreementSandboxEnvironment returns a Terraform configuration string for referencing an agreement-enabled sandbox environment.
// This function provides a data source configuration that references a pre-existing test environment
// configured for agreement and terms of service testing functionality.
// Returns a formatted Terraform data source configuration string for the agreement test environment.
func AgreementSandboxEnvironment() string {
	return `
		data "pingone_environment" "agreement_test" {
			name = "tf-testacc-static-agreements-test"
		}`
}

// DaVinciFlowPolicySandboxEnvironment returns a Terraform configuration string for referencing a DaVinci-enabled sandbox environment.
// This function provides a data source configuration that references a pre-existing test environment
// configured with DaVinci flow policy integration, suitable for tests requiring DaVinci orchestration functionality.
// Returns a formatted Terraform data source configuration string for the DaVinci test environment.
func DaVinciFlowPolicySandboxEnvironment() string {
	return `
		data "pingone_environment" "davinci_test" {
			name = "tf-testacc-static-davinci-test"
		}`
}

// CheckParentEnvironmentDestroy checks whether a parent environment has been properly destroyed or converted to production.
// This function verifies environment destruction by checking if the environment is deleted or has been
// converted to a production environment type, which is considered equivalent to destruction for testing purposes.
// The ctx parameter provides the context for API operations.
// The apiClient parameter is the PingOne Management API client for making environment queries.
// The environmentID parameter specifies the environment to check for destruction.
// Returns true if the environment is destroyed/converted, false otherwise, and any error encountered during the check.
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

// CheckParentUserDestroy checks whether a parent user resource has been properly destroyed.
// This function verifies user deletion by attempting to read the user resource and checking
// for the appropriate HTTP response indicating the user no longer exists.
// The ctx parameter provides the context for API operations.
// The apiClient parameter is the PingOne Management API client for making user queries.
// The environmentID parameter specifies the environment containing the user.
// The userID parameter specifies the user to check for destruction.
// Returns true if the user is destroyed, false otherwise, and any error encountered during the check.
func CheckParentUserDestroy(ctx context.Context, apiClient *management.APIClient, environmentID, userID string) (bool, error) {
	_, r, err := apiClient.UsersApi.ReadUser(ctx, environmentID, userID).Execute()

	return CheckForResourceDestroy(r, err)
}

// CheckForResourceDestroy checks whether a resource has been properly destroyed using the default HTTP status code.
// This function provides a standardized way to verify resource destruction by checking for HTTP 404 responses,
// which indicates the resource no longer exists on the server.
// The r parameter is the HTTP response from the resource query attempt.
// The err parameter is any error returned from the resource query attempt.
// Returns true if the resource is destroyed (404 response), false otherwise, and any error encountered during the check.
func CheckForResourceDestroy(r *http.Response, err error) (bool, error) {
	defaultDestroyHttpCode := 404
	return CheckForResourceDestroyCustomHTTPCode(r, err, defaultDestroyHttpCode)
}

// CheckForResourceDestroyCustomHTTPCode checks whether a resource has been properly destroyed using a custom HTTP status code.
// This function provides a flexible way to verify resource destruction by checking for a specified HTTP response code
// that indicates the resource no longer exists, accommodating APIs that use non-standard status codes for resource absence.
// The r parameter is the HTTP response from the resource query attempt.
// The err parameter is any error returned from the resource query attempt.
// The customHttpCode parameter specifies the HTTP status code to interpret as successful destruction.
// Returns true if the resource is destroyed (custom response code), false otherwise, and any error encountered during the check.
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
