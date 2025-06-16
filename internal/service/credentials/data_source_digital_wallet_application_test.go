// Copyright © 2025 Ping Identity Corporation

package credentials_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/credentials"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccDigitalWalletApplicationDataSource_ByIDFull(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullName := fmt.Sprintf("data.pingone_digital_wallet_application.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             credentials.DigitalWalletApplication_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDigitalWalletApplicationDataSource_ByIDFull(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "digital_wallet_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "application_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "name", name),
					resource.TestCheckResourceAttr(dataSourceFullName, "app_open_url", fmt.Sprintf("https://www.example.com/%s", name)),
				),
			},
		},
	})
}

func TestAccDigitalWalletApplicationDataSource_ByApplicationIDFull(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullName := fmt.Sprintf("data.pingone_digital_wallet_application.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             credentials.DigitalWalletApplication_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDigitalWalletApplicationDataSource_ByApplicationIDFull(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "digital_wallet_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "application_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "name", name),
					resource.TestCheckResourceAttr(dataSourceFullName, "app_open_url", fmt.Sprintf("https://www.example.com/%s", name)),
				),
			},
		},
	})
}

func TestAccDigitalWalletApplicationDataSource_ByNameFull(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullName := fmt.Sprintf("data.pingone_digital_wallet_application.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             credentials.DigitalWalletApplication_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDigitalWalletApplicationDataSource_ByNameFull(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "digital_wallet_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "application_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "name", name),
					resource.TestCheckResourceAttr(dataSourceFullName, "app_open_url", fmt.Sprintf("https://www.example.com/%s", name)),
				),
			},
		},
	})
}

func TestAccDigitalWalletApplicationDataSource_NotFound(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             credentials.DigitalWalletApplication_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccDigitalWalletApplicationDataSource_NotFoundByID(resourceName),
				ExpectError: regexp.MustCompile("Error: Error when calling `ReadOneDigitalWalletApplication`: The request could not be completed. The requested resource was not found."),
			},
			{
				Config:      testAccDigitalWalletApplicationDataSource_NotFoundByApplicationID(resourceName),
				ExpectError: regexp.MustCompile("Error: Cannot find digital wallet application from application_id"),
			},
			{
				Config:      testAccDigitalWalletApplicationDataSource_NotFoundByName(resourceName),
				ExpectError: regexp.MustCompile("Error: Cannot find digital wallet application from name"),
			},
		},
	})
}

func testAccDigitalWalletApplicationDataSource_ByIDFull(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s
resource "pingone_application" "%[2]s-appname" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
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

resource "pingone_digital_wallet_application" "%[2]s-walletappname" {
  environment_id = data.pingone_environment.general_test.id
  application_id = resource.pingone_application.%[2]s-appname.id
  name           = "%[3]s"
  app_open_url   = "https://www.example.com/%[3]s"

  depends_on = [resource.pingone_application.%[2]s-appname]
}

data "pingone_digital_wallet_application" "%[2]s" {
  environment_id    = data.pingone_environment.general_test.id
  digital_wallet_id = resource.pingone_digital_wallet_application.%[2]s-walletappname.id

  depends_on = [resource.pingone_digital_wallet_application.%[2]s-walletappname]

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccDigitalWalletApplicationDataSource_ByApplicationIDFull(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s
resource "pingone_application" "%[2]s-appname" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true

  oidc_options = {
    type                       = "NATIVE_APP"
    grant_types                = ["CLIENT_CREDENTIALS"]
    token_endpoint_auth_method = "CLIENT_SECRET_BASIC"

    mobile_app = {
      bundle_id                = "com.pingidentity.ios_wallet_byappid"
      package_name             = "com.pingidentity.android_wallet_byappid"
      passcode_refresh_seconds = 30
    }
  }
}

resource "pingone_digital_wallet_application" "%[2]s-walletappname" {
  environment_id = data.pingone_environment.general_test.id
  application_id = resource.pingone_application.%[2]s-appname.id
  name           = "%[3]s"
  app_open_url   = "https://www.example.com/%[3]s"

  depends_on = [resource.pingone_application.%[2]s-appname]
}

data "pingone_digital_wallet_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  application_id = resource.pingone_application.%[2]s-appname.id

  depends_on = [resource.pingone_digital_wallet_application.%[2]s-walletappname]

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccDigitalWalletApplicationDataSource_ByNameFull(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s
resource "pingone_application" "%[2]s-appname" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true

  oidc_options = {
    type                       = "NATIVE_APP"
    grant_types                = ["CLIENT_CREDENTIALS"]
    token_endpoint_auth_method = "CLIENT_SECRET_BASIC"

    mobile_app = {
      bundle_id                = "com.pingidentity.ios_%[2]s"
      package_name             = "com.pingidentity.android_%[2]s"
      passcode_refresh_seconds = 30
    }
  }
}

resource "pingone_digital_wallet_application" "%[2]s-walletappname" {
  environment_id = data.pingone_environment.general_test.id
  application_id = resource.pingone_application.%[2]s-appname.id
  name           = "%[3]s"
  app_open_url   = "https://www.example.com/%[3]s"

  depends_on = [resource.pingone_application.%[2]s-appname]
}

data "pingone_digital_wallet_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  depends_on = [resource.pingone_digital_wallet_application.%[2]s-walletappname]

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccDigitalWalletApplicationDataSource_NotFoundByID(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_digital_wallet_application" "%[2]s" {
  environment_id    = data.pingone_environment.general_test.id
  digital_wallet_id = "9c052a8a-14be-44e4-8f07-2662569994ce" // dummy ID that conforms to UUID v4

}`, acctest.GenericSandboxEnvironment(), resourceName)
}

func testAccDigitalWalletApplicationDataSource_NotFoundByApplicationID(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_digital_wallet_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  application_id = "9c052a8a-14be-44e4-8f07-2662569994ce" // dummy ID that conforms to UUID v4

}`, acctest.GenericSandboxEnvironment(), resourceName)
}

func testAccDigitalWalletApplicationDataSource_NotFoundByName(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_digital_wallet_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "dummy value"

}`, acctest.GenericSandboxEnvironment(), resourceName)
}
