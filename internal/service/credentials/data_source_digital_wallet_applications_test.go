package credentials_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccDigitalWalletApplicationsDataSource_NoFilter(t *testing.T) {
	// If run in parallel, unique environments are needed to prevent collisions within the same environment.
	//t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullName := fmt.Sprintf("data.pingone_digital_wallet_applications.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDigitalWalletApplicationDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDigitalWalletApplicationsDataSource_NoFilter(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(dataSourceFullName, "ids.#", "3"),
					resource.TestMatchResourceAttr(dataSourceFullName, "ids.0", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(dataSourceFullName, "ids.1", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(dataSourceFullName, "ids.2", verify.P1ResourceIDRegexp),
				),
			},
		},
	})
}

func TestAccDigitalWalletApplicationsDataSource_NotFound(t *testing.T) {
	// If run in parallel, unique environments are needed to prevent collisions within the same environment.
	//t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullName := fmt.Sprintf("data.pingone_digital_wallet_applications.%s", resourceName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCredentialTypeDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDigitalWalletApplicationsDataSource_NotFound(resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(dataSourceFullName, "ids.#", "0"),
				),
			},
		},
	})
}

func testAccDigitalWalletApplicationsDataSource_NoFilter(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s
resource "pingone_application" "%[2]s-appname1" {
	environment_id = data.pingone_environment.credentials_test.id
	name = "%[2]s-appname1"
	enabled        = true

	oidc_options {
	  type                        = "NATIVE_APP"
	  grant_types                 = ["CLIENT_CREDENTIALS"]
	  token_endpoint_authn_method = "CLIENT_SECRET_BASIC"
	  bundle_id        = "com.pingidentity.ios_wallet_byid"
	  package_name     = "com.pingidentity.android_wallet_byid"
  
	  mobile_app {
		 bundle_id        = "com.pingidentity.ios_wallet_byid"
		 package_name     = "com.pingidentity.android_wallet_byid"
		 passcode_refresh_seconds = 30
	  }
	}
}

resource "pingone_digital_wallet_application" "%[2]s-walletappname1" {
	environment_id = data.pingone_environment.credentials_test.id
	application_id = resource.pingone_application.%[2]s-appname1.id
	name = "%[2]s-name1"
	app_open_url = "https://www.example.com"
}

resource "pingone_application" "%[2]s-appname2" {
	environment_id = data.pingone_environment.credentials_test.id
	name = "%[2]s-appname2"
	enabled        = true

	oidc_options {
	  type                        = "NATIVE_APP"
	  grant_types                 = ["CLIENT_CREDENTIALS"]
	  token_endpoint_authn_method = "CLIENT_SECRET_BASIC"
	  bundle_id        = "com.pingidentity.ios_wallet2"
	  package_name     = "com.pingidentity.android_wallet2"
  
	  mobile_app {
		 bundle_id        = "com.pingidentity.ios_wallet2"
		 package_name     = "com.pingidentity.android_wallet2"
		 passcode_refresh_seconds = 30
	  }
	}
}

resource "pingone_digital_wallet_application" "%[2]s-walletappname2" {
	environment_id = data.pingone_environment.credentials_test.id
	application_id = resource.pingone_application.%[2]s-appname2.id
	name = "%[2]s-name2"
	app_open_url = "https://www.example.com"
}

resource "pingone_application" "%[2]s-appname3" {
	environment_id = data.pingone_environment.credentials_test.id
	name = "%[2]s-appname3"
	enabled        = true

	oidc_options {
	  type                        = "NATIVE_APP"
	  grant_types                 = ["CLIENT_CREDENTIALS"]
	  token_endpoint_authn_method = "CLIENT_SECRET_BASIC"
	  bundle_id        = "com.pingidentity.ios_wallet3"
	  package_name     = "com.pingidentity.android_wallet3"
  
	  mobile_app {
		 bundle_id        = "com.pingidentity.ios_wallet3"
		 package_name     = "com.pingidentity.android_wallet3"
		 passcode_refresh_seconds = 30
	  }
	}
}

resource "pingone_digital_wallet_application" "%[2]s-walletappname3" {
	environment_id = data.pingone_environment.credentials_test.id
	application_id = resource.pingone_application.%[2]s-appname3.id
	name = "%[2]s-name3"
	app_open_url = "https://www.example.com"
}

data "pingone_digital_wallet_applications" "%[2]s" {
	environment_id = data.pingone_environment.credentials_test.id

	depends_on = [pingone_digital_wallet_application.%[2]s-walletappname1, pingone_digital_wallet_application.%[2]s-walletappname2, pingone_digital_wallet_application.%[2]s-walletappname3]
  }`, acctest.CredentialsSandboxEnvironment(), resourceName, name)
}

func testAccDigitalWalletApplicationsDataSource_NotFound(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_digital_wallet_applications" "%[2]s" {
	environment_id = data.pingone_environment.credentials_test.id

  }`, acctest.CredentialsSandboxEnvironment(), resourceName)
}
