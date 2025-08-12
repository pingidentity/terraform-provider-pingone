// Copyright © 2025 Ping Identity Corporation

package acctest

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	clientconfig "github.com/pingidentity/pingone-go-client/config"
	"github.com/pingidentity/pingone-go-client/oauth2"
	"github.com/pingidentity/pingone-go-client/pingone"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
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
			providerServerFactory, err := provider.ProviderServerFactoryV6(ctx, GetProviderTestingVersion())

			if err != nil {
				return nil, err
			}

			return providerServerFactory(), nil
		}
	}

	return factories
}

func GetProviderTestingVersion() string {
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

type EnumFeatureFlag string

const (
	ENUMFEATUREFLAG_DAVINCI EnumFeatureFlag = "DAVINCI"
)

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

func PreCheckNoFeatureFlag(t *testing.T) {
	PreCheckFeatureFlag(t, "")
}

func PreCheckFeatureFlag(t *testing.T, flag EnumFeatureFlag) {
	if v := os.Getenv("FEATURE_FLAG"); v != string(flag) {
		t.Skipf("Skipping feature flag test.  Flag required: \"%s\"", string(flag))
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

func TestClient(ctx context.Context) (*pingone.APIClient, error) {
	regionSuffix, ok := framework.RegionSuffixFromCode(strings.ToLower(os.Getenv("PINGONE_REGION_CODE")))
	if !ok {
		return nil, fmt.Errorf("invalid PINGONE_REGION_CODE: %s", os.Getenv("PINGONE_REGION_CODE"))
	}
	config := clientconfig.NewConfiguration().
		WithGrantType(oauth2.GrantTypeClientCredentials).
		WithTopLevelDomain(regionSuffix)

	pingOneConfig := pingone.NewConfiguration(config)
	pingOneConfig.UserAgent = framework.UserAgent("", GetProviderTestingVersion())

	return pingone.NewAPIClient(pingOneConfig)

}

func PreCheckTestClient(ctx context.Context, t *testing.T) *pingone.APIClient {
	p1Client, err := TestClient(ctx)

	if err != nil {
		t.Fatalf("Failed to get API client: %v", err)
	}

	return p1Client
}

func DaVinciSandboxEnvironment(withBootstrapConfig bool) string {
	if withBootstrapConfig {
		generalName := "general_test"
		return DaVinciBootstrappedSandboxEnvironment(&generalName)
	} else {
		return GenericSandboxEnvironment()
	}
}

func GenericSandboxEnvironment() string {
	return `
		data "pingone_environment" "general_test" {
			name = "tf-testacc-dynamic-general-test"
		}`
}

func DaVinciBootstrappedSandboxEnvironment(dataSourceName *string) string {
	var name string
	if dataSourceName != nil {
		name = *dataSourceName
	} else {
		name = "davinci_bootstrapped_test"
	}
	return fmt.Sprintf(`
		data "pingone_environment" "%s" {
			name = "tf-testacc-dynamic-davinci-bootstrapped-test"
		}`, name)
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

func DaVinciFlowPolicySandboxEnvironment() string {
	return `
		data "pingone_environment" "davinci_test" {
			name = "tf-testacc-static-davinci-test"
		}`
}

func CheckParentEnvironmentDestroy(ctx context.Context, apiClient *pingone.APIClient, environmentID string) (bool, error) {
	environmentIdUuid, err := uuid.Parse(environmentID)
	if err != nil {
		return false, fmt.Errorf("unable to parse environment id '%s' as uuid: %v", environmentID, err)
	}

	//TODO remove placeholder expand once pingone-go-client is updated to remove this requirement
	environment, r, err := apiClient.EnvironmentsApi.GetEnvironmentById(ctx, environmentIdUuid).Expand("placeholder").Execute()

	destroyed, err := CheckForResourceDestroy(r, err)
	if err != nil {
		return destroyed, err
	}

	if destroyed {
		return destroyed, nil
	} else {
		if environment != nil && environment.Type == pingone.ENVIRONMENTTYPEVALUE_PRODUCTION {
			return true, nil
		} else {
			return false, nil
		}
	}
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
