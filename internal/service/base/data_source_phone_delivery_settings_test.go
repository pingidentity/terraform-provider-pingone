// Copyright © 2026 Ping Identity Corporation

package base_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	acctestlegacysdk "github.com/pingidentity/terraform-provider-pingone/internal/acctest/legacysdk"
	baselegacysdk "github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/base/legacysdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccPhoneDeliverySettingsDataSource_ById(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_phone_delivery_settings.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.pingone_phone_delivery_settings.%s", resourceName)

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	environmentName := acctest.ResourceNameGenEnvironment()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             baselegacysdk.Environment_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPhoneDeliverySettingsDataSourceConfig_ById(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "phone_delivery_settings_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "phone_delivery_settings_id", resourceFullName, "id"),
					resource.TestCheckResourceAttr(dataSourceFullName, "provider_type", "CUSTOM_PROVIDER"),
					resource.TestCheckResourceAttr(dataSourceFullName, "provider_custom.name", name),
					resource.TestCheckResourceAttr(dataSourceFullName, "provider_custom.authentication.method", "BEARER"),
					// `auth_token` is a sensitive, write-only parameter: the service
					// doesn't return it and the data source has no plan to carry it
					// over from, so it reads back as null.
					resource.TestCheckNoResourceAttr(dataSourceFullName, "provider_custom.authentication.auth_token"),
					resource.TestCheckNoResourceAttr(dataSourceFullName, "provider_custom.numbers"),
					resource.TestCheckNoResourceAttr(dataSourceFullName, "provider_custom_twilio"),
					resource.TestCheckNoResourceAttr(dataSourceFullName, "provider_custom_syniverse"),
					resource.TestMatchResourceAttr(dataSourceFullName, "created_at", verify.RFC3339Regexp),
					resource.TestMatchResourceAttr(dataSourceFullName, "updated_at", verify.RFC3339Regexp),
				),
			},
		},
	})
}

func TestAccPhoneDeliverySettingsDataSource_ByName(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_phone_delivery_settings.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.pingone_phone_delivery_settings.%s", resourceName)

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	environmentName := acctest.ResourceNameGenEnvironment()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             baselegacysdk.Environment_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPhoneDeliverySettingsDataSourceConfig_ByName(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "phone_delivery_settings_id", resourceFullName, "id"),
					resource.TestCheckResourceAttr(dataSourceFullName, "name", name),
					resource.TestCheckResourceAttr(dataSourceFullName, "provider_type", "CUSTOM_PROVIDER"),
					resource.TestCheckResourceAttr(dataSourceFullName, "provider_custom.name", name),
					resource.TestCheckResourceAttr(dataSourceFullName, "provider_custom.authentication.method", "BEARER"),
					resource.TestCheckNoResourceAttr(dataSourceFullName, "provider_custom_twilio"),
					resource.TestCheckNoResourceAttr(dataSourceFullName, "provider_custom_syniverse"),
				),
			},
		},
	})
}

func TestAccPhoneDeliverySettingsDataSource_NotFound(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             baselegacysdk.Environment_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccPhoneDeliverySettingsDataSourceConfig_NotFound(environmentName, licenseID, resourceName),
				ExpectError: regexp.MustCompile(`Cannot find phone delivery settings from name`),
			},
		},
	})
}

func TestAccPhoneDeliverySettingsDataSource_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             baselegacysdk.Environment_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccPhoneDeliverySettingsDataSourceConfig_BadParameters(environmentName, licenseID, resourceName),
				ExpectError: regexp.MustCompile(`Invalid Attribute Combination`),
			},
		},
	})
}

func testAccPhoneDeliverySettingsDataSourceConfig_ById(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_phone_delivery_settings" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  provider_custom = {
    name = "%[4]s"

    authentication = {
      method     = "BEARER"
      auth_token = "testtoken"
    }

    requests = [
      {
        delivery_method     = "SMS"
        method              = "GET"
        phone_number_format = "FULL"
        url                 = "https://pingdevops.com/fake-send-to-test?to=$${to}&message=$${message}"
      }
    ]
  }
}

data "pingone_phone_delivery_settings" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  phone_delivery_settings_id = pingone_phone_delivery_settings.%[3]s.id
}`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccPhoneDeliverySettingsDataSourceConfig_ByName(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_phone_delivery_settings" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  provider_custom = {
    name = "%[4]s"

    authentication = {
      method     = "BEARER"
      auth_token = "testtoken"
    }

    requests = [
      {
        delivery_method     = "SMS"
        method              = "GET"
        phone_number_format = "FULL"
        url                 = "https://pingdevops.com/fake-send-to-test?to=$${to}&message=$${message}"
      }
    ]
  }
}

data "pingone_phone_delivery_settings" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  name = pingone_phone_delivery_settings.%[3]s.provider_custom.name
}`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccPhoneDeliverySettingsDataSourceConfig_NotFound(environmentName, licenseID, resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_phone_delivery_settings" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  name = "tf-testacc-unknown-name"
}`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName)
}

func testAccPhoneDeliverySettingsDataSourceConfig_BadParameters(environmentName, licenseID, resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_phone_delivery_settings" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
}`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName)
}
