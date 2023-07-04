package credentials_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func testAccCheckDigitalWalletApplicationDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.CredentialsAPIClient

	mgmtApiClient := p1Client.API.ManagementAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_digital_wallet_application" {
			continue
		}

		_, rEnv, err := mgmtApiClient.EnvironmentsApi.ReadOneEnvironment(ctx, rs.Primary.Attributes["environment_id"]).Execute()

		if err != nil {

			if rEnv == nil {
				return fmt.Errorf("Response object does not exist and no error detected")
			}

			if rEnv.StatusCode == 404 {
				continue
			}

			return err
		}

		body, r, err := apiClient.DigitalWalletAppsApi.ReadOneDigitalWalletApp(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["id"]).Execute()

		if err != nil {

			if r == nil {
				return fmt.Errorf("Response object does not exist and no error detected")
			}

			if r.StatusCode == 404 {
				continue
			}

			tflog.Error(ctx, fmt.Sprintf("Error: %v", body))
			return err
		}

		return fmt.Errorf("PingOne Digital Wallet Application ID e %s still exists", rs.Primary.ID)
	}

	return nil
}

func TestAccDigitalWalletApplication_NewEnv(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_digital_wallet_application.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := acctest.ResourceNameGen()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCredentialTypeDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDigitalWalletApplication_NewEnv(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
				),
			},
		},
	})
}

func TestAccDigitalWalletApplication_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_digital_wallet_application.%s", resourceName)

	appOpenUrl := "https://www.example.com/appopen"
	name := acctest.ResourceNameGen()

	fullStep := resource.TestStep{
		Config: testAccDigitalWalletApplication_Full(resourceName, name, appOpenUrl),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
			resource.TestMatchResourceAttr(resourceFullName, "application_id", verify.P1ResourceIDRegexp),
			resource.TestCheckResourceAttr(resourceFullName, "app_open_url", appOpenUrl),
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
		),
	}

	updatedAppOpenUrl := "https://www.example.com/v2/appopen"
	updatedName := acctest.ResourceNameGen()

	updateStep := resource.TestStep{
		Config: testAccDigitalWalletApplication_Full(resourceName, updatedName, updatedAppOpenUrl),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
			resource.TestMatchResourceAttr(resourceFullName, "application_id", verify.P1ResourceIDRegexp),
			resource.TestCheckResourceAttr(resourceFullName, "app_open_url", updatedAppOpenUrl),
			resource.TestCheckResourceAttr(resourceFullName, "name", updatedName),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDigitalWalletApplicationDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			fullStep,
			{
				Config:  testAccDigitalWalletApplication_Full(resourceName, name, appOpenUrl),
				Destroy: true,
			},
			updateStep,
			{
				Config:  testAccDigitalWalletApplication_Full(resourceName, updatedName, updatedAppOpenUrl),
				Destroy: true,
			},
			// changes
			fullStep,
			updateStep,
			{
				Config:  testAccDigitalWalletApplication_Full(resourceName, name, appOpenUrl),
				Destroy: true,
			},
			{
				Config:  testAccDigitalWalletApplication_Full(resourceName, updatedName, updatedAppOpenUrl),
				Destroy: true,
			},
		},
	})
}

func TestAccDigitalWalletApplication_InvalidNativeApplication(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	appOpenUrl := "https://www.example.com/appopen"
	name := acctest.ResourceNameGen()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDigitalWalletApplicationDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccDigitalWalletApplication_NativeAppMissing(resourceName, name, appOpenUrl),
				ExpectError: regexp.MustCompile("(Error: Error when calling `ReadOneApplication`: Unable to find Application with ID: '9c052a8a-14be-44e4-8f07-2662569994ce' in Environment).*"),
			},
			{
				Config:      testAccDigitalWalletApplication_InvalidAppType(resourceName, name, appOpenUrl),
				ExpectError: regexp.MustCompile("Error: Application referenced in `application.id` is OIDC, but is not the required `Native` OIDC application type"),
			},
			{
				Config:      testAccDigitalWalletApplication_NativeAppMobileNotConfigured(resourceName, name, appOpenUrl),
				ExpectError: regexp.MustCompile("Error: Application referenced in `application.id` does not contain mobile application configuration"),
			},
			{
				Config:      testAccDigitalWalletApplication_NativeAppInvalidMobileConfiguration(resourceName, name, appOpenUrl),
				ExpectError: regexp.MustCompile("Error: Application referenced in `application.id` does not contain mobile application configuration"),
			},
		},
	})
}

func TestAccDigitalWalletApplication_InvalidAppOpenUrl(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	appOpenUrl := "www.example.com"
	name := acctest.ResourceNameGen()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDigitalWalletApplicationDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccDigitalWalletApplication_InvalidAppOpenUrl(resourceName, name, appOpenUrl),
				ExpectError: regexp.MustCompile("Error: Error when calling `CreateDigitalWalletApplication`: Validation Error : \\[appOpenUrl must be a valid URL\\]"),
			},
		},
	})
}

func testAccDigitalWalletApplication_NewEnv(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s"
  enabled        = true

  oidc_options {
    type                        = "NATIVE_APP"
    grant_types                 = ["AUTHORIZATION_CODE"]
    response_types              = ["CODE"]
    pkce_enforcement            = "S256_REQUIRED"
    token_endpoint_authn_method = "NONE"
    redirect_uris               = ["https://www.pingidentity.com"]

    mobile_app {
      bundle_id                = "com.pingidentity.ios_%[4]s"
      package_name             = "com.pingidentity.android_%[4]s"
      passcode_refresh_seconds = 30
    }
  }
}

resource "pingone_digital_wallet_application" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  application_id = resource.pingone_application.%[3]s.id
  name           = "%[4]s"
  app_open_url   = "https://www.example.com/appopen"
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccDigitalWalletApplication_Full(resourceName, name, appOpenUrl string) string {
	return fmt.Sprintf(`
	%[1]s
resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true

  oidc_options {
    type                        = "NATIVE_APP"
    grant_types                 = ["CLIENT_CREDENTIALS"]
    token_endpoint_authn_method = "CLIENT_SECRET_BASIC"

    mobile_app {
      bundle_id                = "com.pingidentity.ios_%[3]s"
      package_name             = "com.pingidentity.android_%[3]s"
      passcode_refresh_seconds = 30
    }
  }
}

resource "pingone_digital_wallet_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  application_id = resource.pingone_application.%[2]s.id
  name           = "%[3]s"
  app_open_url   = "%[4]s"
}`, acctest.GenericSandboxEnvironment(), resourceName, name, appOpenUrl)
}

func testAccDigitalWalletApplication_NativeAppMissing(resourceName, name, appOpenUrl string) string {
	return fmt.Sprintf(`
	%[1]s
resource "pingone_digital_wallet_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  application_id = "9c052a8a-14be-44e4-8f07-2662569994ce" // dummy ID that conforms to UUID v4
  name           = "%[3]s"
  app_open_url   = "%[4]s"
}`, acctest.GenericSandboxEnvironment(), resourceName, name, appOpenUrl)
}

func testAccDigitalWalletApplication_InvalidAppType(resourceName, name, appOpenUrl string) string {
	return fmt.Sprintf(`
	%[1]s
resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true

  oidc_options {
    type                        = "SINGLE_PAGE_APP"
    grant_types                 = ["AUTHORIZATION_CODE"]
    response_types              = ["CODE"]
    pkce_enforcement            = "S256_REQUIRED"
    token_endpoint_authn_method = "NONE"
    redirect_uris               = ["https://www.pingidentity.com"]
  }
}

resource "pingone_digital_wallet_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  application_id = resource.pingone_application.%[2]s.id
  name           = "%[3]s"
  app_open_url   = "%[4]s"
}`, acctest.GenericSandboxEnvironment(), resourceName, name, appOpenUrl)
}

func testAccDigitalWalletApplication_NativeAppMobileNotConfigured(resourceName, name, appOpenUrl string) string {
	return fmt.Sprintf(`
	%[1]s
resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true

  oidc_options {
    type                        = "NATIVE_APP"
    grant_types                 = ["CLIENT_CREDENTIALS"]
    token_endpoint_authn_method = "CLIENT_SECRET_BASIC"
  }
}

resource "pingone_digital_wallet_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  application_id = resource.pingone_application.%[2]s.id
  name           = "%[3]s"
  app_open_url   = "%[4]s"
}`, acctest.GenericSandboxEnvironment(), resourceName, name, appOpenUrl)
}

func testAccDigitalWalletApplication_NativeAppInvalidMobileConfiguration(resourceName, name, appOpenUrl string) string {
	return fmt.Sprintf(`
	%[1]s
resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true

  oidc_options {
    type                        = "NATIVE_APP"
    grant_types                 = ["CLIENT_CREDENTIALS"]
    token_endpoint_authn_method = "CLIENT_SECRET_BASIC"

    mobile_app {
      passcode_refresh_seconds = 30
    }
  }
}

resource "pingone_digital_wallet_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  application_id = resource.pingone_application.%[2]s.id
  name           = "%[3]s"
  app_open_url   = "%[4]s"
}`, acctest.GenericSandboxEnvironment(), resourceName, name, appOpenUrl)
}

func testAccDigitalWalletApplication_InvalidAppOpenUrl(resourceName, name, appOpenUrl string) string {
	return fmt.Sprintf(`
	%[1]s
resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true

  oidc_options {
    type                        = "NATIVE_APP"
    grant_types                 = ["CLIENT_CREDENTIALS"]
    token_endpoint_authn_method = "CLIENT_SECRET_BASIC"

    mobile_app {
      bundle_id                = "com.pingidentity.ios_%[3]s"
      package_name             = "com.pingidentity.android_%[3]s"
      passcode_refresh_seconds = 30
    }
  }
}

resource "pingone_digital_wallet_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  application_id = resource.pingone_application.%[2]s.id
  name           = "%[3]s"
  app_open_url   = "%[4]s"
}`, acctest.GenericSandboxEnvironment(), resourceName, name, appOpenUrl)
}
