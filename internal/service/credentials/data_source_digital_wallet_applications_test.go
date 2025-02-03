// Copyright © 2025 Ping Identity Corporation

package credentials_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/credentials"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccDigitalWalletApplicationsDataSource_NoFilter(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullName := fmt.Sprintf("data.pingone_digital_wallet_applications.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := acctest.ResourceNameGen()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             credentials.DigitalWalletApplication_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDigitalWalletApplicationsDataSource_NoFilter(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "ids.#", "3"),
					resource.TestMatchResourceAttr(dataSourceFullName, "ids.0", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "ids.1", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "ids.2", verify.P1ResourceIDRegexpFullString),
				),
			},
			{
				Config:  testAccDigitalWalletApplicationsDataSource_NoFilter(environmentName, licenseID, resourceName, name),
				Destroy: true,
			},
		},
	})
}

func TestAccDigitalWalletApplicationsDataSource_NotFound(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullName := fmt.Sprintf("data.pingone_digital_wallet_applications.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := acctest.ResourceNameGen()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             credentials.DigitalWalletApplication_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDigitalWalletApplicationsDataSource_NotFound(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "ids.#", "0"),
				),
			},
		},
	})
}

func testAccDigitalWalletApplicationsDataSource_NoFilter(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s
resource "pingone_application" "%[3]s-appname1" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s-appname1"
  enabled        = true

  oidc_options = {
    type                       = "NATIVE_APP"
    grant_types                = ["CLIENT_CREDENTIALS"]
    token_endpoint_auth_method = "CLIENT_SECRET_BASIC"

    mobile_app = {
      bundle_id                = "com.pingidentity.ios_wallet_byid"
      package_name             = "com.pingidentity.android_wallet_byid"
      passcode_refresh_seconds = 30
    }
  }
}

resource "pingone_digital_wallet_application" "%[3]s-walletappname1" {
  environment_id = pingone_environment.%[2]s.id
  application_id = resource.pingone_application.%[3]s-appname1.id
  name           = "%[4]s-name1"
  app_open_url   = "https://www.example.com"
}

resource "pingone_application" "%[3]s-appname2" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s-appname2"
  enabled        = true

  oidc_options = {
    type                       = "NATIVE_APP"
    grant_types                = ["AUTHORIZATION_CODE"]
    response_types             = ["CODE"]
    pkce_enforcement           = "S256_REQUIRED"
    token_endpoint_auth_method = "NONE"
    redirect_uris = [
      "https://www.example.com/app/callback",
    ]

    mobile_app = {
      bundle_id                = "com.pingidentity.ios_wallet2"
      package_name             = "com.pingidentity.android_wallet2"
      passcode_refresh_seconds = 30
    }
  }
}

resource "pingone_digital_wallet_application" "%[3]s-walletappname2" {
  environment_id = pingone_environment.%[2]s.id
  application_id = resource.pingone_application.%[3]s-appname2.id
  name           = "%[4]s-name2"
  app_open_url   = "https://www.example.com"
}

resource "pingone_application" "%[3]s-appname3" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s-appname3"
  enabled        = true

  oidc_options = {
    type                       = "NATIVE_APP"
    grant_types                = ["CLIENT_CREDENTIALS"]
    token_endpoint_auth_method = "CLIENT_SECRET_BASIC"

    mobile_app = {
      bundle_id                = "com.pingidentity.ios_wallet3"
      package_name             = "com.pingidentity.android_wallet3"
      passcode_refresh_seconds = 30
    }
  }
}

resource "pingone_digital_wallet_application" "%[3]s-walletappname3" {
  environment_id = pingone_environment.%[2]s.id
  application_id = resource.pingone_application.%[3]s-appname3.id
  name           = "%[4]s-name3"
  app_open_url   = "https://www.example.com"
}

data "pingone_digital_wallet_applications" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  depends_on = [pingone_digital_wallet_application.%[3]s-walletappname1, pingone_digital_wallet_application.%[3]s-walletappname2, pingone_digital_wallet_application.%[3]s-walletappname3]
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccDigitalWalletApplicationsDataSource_NotFound(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_digital_wallet_applications" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}
