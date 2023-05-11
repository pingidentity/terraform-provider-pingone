package credentials_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccDigitalWalletApplicationDataSource_ByIDFull(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullName := fmt.Sprintf("data.pingone_digital_wallet_application.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             nil,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDigitalWalletApplicationDataSourceConfigDataSource_ByIDFull(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(dataSourceFullName, "digital_wallet_id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(dataSourceFullName, "application_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "name"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "app_open_url"),
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
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             nil,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDigitalWalletApplicationDataSourceConfigDataSource_ByApplicationIDFull(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(dataSourceFullName, "digital_wallet_id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(dataSourceFullName, "application_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "name"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "app_open_url"),
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
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             nil,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDigitalWalletApplicationDataSourceConfigDataSource_ByNameFull(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(dataSourceFullName, "digital_wallet_id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(dataSourceFullName, "application_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "name"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "app_open_url"),
				),
			},
		},
	})
}

func TestAccDigitalWalletApplicationDataSource_NotFoundByID(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             nil,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccDigitalWalletApplicationDataSourceConfigDataSource_NotFoundByID(resourceName),
				ExpectError: regexp.MustCompile("Error: Cannot find digital wallet application from id"),
			},
		},
	})
}

func TestAccDigitalWalletApplicationDataSource_NotFoundByApplicationID(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             nil,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccDigitalWalletApplicationDataSourceConfigDataSource_NotFoundByApplicationID(resourceName),
				ExpectError: regexp.MustCompile("Error: Cannot find digital wallet application from application_id"),
			},
		},
	})
}

func TestAccDigitalWalletApplicationDataSource_NotFoundByName(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             nil,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccDigitalWalletApplicationDataSourceConfigDataSource_NotFoundByName(resourceName),
				ExpectError: regexp.MustCompile("Error: Cannot find digital wallet application from name"),
			},
		},
	})
}

func testAccDigitalWalletApplicationDataSourceConfigDataSource_ByIDFull(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s
resource "pingone_application" "%[2]s-appname" {
	environment_id = data.pingone_environment.credentials_test.id
	name = "%[2]s-appname"
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

resource "pingone_digital_wallet_application" "%[2]s-walletappname" {
	environment_id = data.pingone_environment.credentials_test.id
	application_id = resource.pingone_application.%[2]s-appname.id
	name = "%[2]s-name"
	app_open_url = "https://www.example.com"
}

data "pingone_digital_wallet_application" "%[2]s" {
	environment_id = data.pingone_environment.credentials_test.id
	digital_wallet_id = resource.pingone_digital_wallet_application.%[2]s-walletappname.id

  }`, acctest.CredentialsSandboxEnvironment(), resourceName, name)
}

func testAccDigitalWalletApplicationDataSourceConfigDataSource_ByApplicationIDFull(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s
resource "pingone_application" "%[2]s-appname" {
	environment_id = data.pingone_environment.credentials_test.id
	name = "%[2]s-appname"
	enabled        = true

	oidc_options {
	  type                        = "NATIVE_APP"
	  grant_types                 = ["CLIENT_CREDENTIALS"]
	  token_endpoint_authn_method = "CLIENT_SECRET_BASIC"
	  bundle_id        = "com.pingidentity.ios_wallet_byappid"
	  package_name     = "com.pingidentity.android_wallet_byappid"
  
	  mobile_app {
		 bundle_id        = "com.pingidentity.ios_wallet_byappid"
		 package_name     = "com.pingidentity.android_wallet_byappid"
		 passcode_refresh_seconds = 30
	  }
	}
}

resource "pingone_digital_wallet_application" "%[2]s-walletappname" {
	environment_id = data.pingone_environment.credentials_test.id
	application_id = resource.pingone_application.%[2]s-appname.id
	name = "%[2]s-name"
	app_open_url = "https://www.example.com"	

	depends_on = [resource.pingone_application.%[2]s-appname]
}

data "pingone_digital_wallet_application" "%[2]s" {
	environment_id = data.pingone_environment.credentials_test.id
	application_id = resource.pingone_application.%[2]s-appname.id

	depends_on = [resource.pingone_digital_wallet_application.%[2]s-walletappname]

  }`, acctest.CredentialsSandboxEnvironment(), resourceName, name)
}

func testAccDigitalWalletApplicationDataSourceConfigDataSource_ByNameFull(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s
resource "pingone_application" "%[2]s-appname" {
	environment_id = data.pingone_environment.credentials_test.id
	name = "%[2]s-appname"
	enabled        = true

	oidc_options {
	  type                        = "NATIVE_APP"
	  grant_types                 = ["CLIENT_CREDENTIALS"]
	  token_endpoint_authn_method = "CLIENT_SECRET_BASIC"
	  bundle_id        = "com.pingidentity.ios_wallet_byname"
	  package_name     = "com.pingidentity.android_wallet_byname"
  
	  mobile_app {
		 bundle_id        = "com.pingidentity.ios_wallet_byname"
		 package_name     = "com.pingidentity.android_wallet_byname"
		 passcode_refresh_seconds = 30
	  }
	}
}

resource "pingone_digital_wallet_application" "%[2]s-walletappname" {
	environment_id = data.pingone_environment.credentials_test.id
	application_id = resource.pingone_application.%[2]s-appname.id
	name = "%[2]s-name"
	app_open_url = "https://www.example.com"	

	depends_on = [resource.pingone_application.%[2]s-appname]
}

data "pingone_digital_wallet_application" "%[2]s" {
	environment_id = data.pingone_environment.credentials_test.id
	name = "%[2]s-name"

	depends_on = [resource.pingone_digital_wallet_application.%[2]s-walletappname]

  }`, acctest.CredentialsSandboxEnvironment(), resourceName, name)
}

func testAccDigitalWalletApplicationDataSourceConfigDataSource_NotFoundByID(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_digital_wallet_application" "%[2]s" {
	environment_id = data.pingone_environment.credentials_test.id
	digital_wallet_id = "9c052a8a-14be-44e4-8f07-2662569994ce" // dummy ID that conforms to UUID v4

  }`, acctest.CredentialsSandboxEnvironment(), resourceName)
}

func testAccDigitalWalletApplicationDataSourceConfigDataSource_NotFoundByApplicationID(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_digital_wallet_application" "%[2]s" {
	environment_id = data.pingone_environment.credentials_test.id
	application_id = "9c052a8a-14be-44e4-8f07-2662569994ce" // dummy ID that conforms to UUID v4

  }`, acctest.CredentialsSandboxEnvironment(), resourceName)
}

func testAccDigitalWalletApplicationDataSourceConfigDataSource_NotFoundByName(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_digital_wallet_application" "%[2]s" {
	environment_id = data.pingone_environment.credentials_test.id
	name = "dummy value"

  }`, acctest.CredentialsSandboxEnvironment(), resourceName)
}
