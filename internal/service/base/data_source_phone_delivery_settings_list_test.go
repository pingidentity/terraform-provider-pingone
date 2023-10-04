package base_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/base"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccPhoneDeliverySettingsListDataSource_ByAll(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullName := fmt.Sprintf("data.pingone_phone_delivery_settings_list.%s", resourceName)

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	environmentName := acctest.ResourceNameGenEnvironment()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.TestAccCheckEnvironmentDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPhoneDeliverySettingsListDataSourceConfig_ByAll(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "ids.#", "1"),
					resource.TestMatchResourceAttr(dataSourceFullName, "ids.0", verify.P1ResourceIDRegexpFullString),
				),
			},
		},
	})
}

func TestAccPhoneDeliverySettingsListDataSource_NotFound(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullName := fmt.Sprintf("data.pingone_phone_delivery_settings_list.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.TestAccCheckEnvironmentDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPhoneDeliverySettingsListDataSourceConfig_NotFound(environmentName, licenseID, resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "ids.#", "0"),
				),
			},
		},
	})
}

func testAccPhoneDeliverySettingsListDataSourceConfig_ByAll(environmentName, licenseID, resourceName, name string) string {
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

data "pingone_phone_delivery_settings_list" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  depends_on = [
    pingone_phone_delivery_settings.%[3]s,
  ]
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccPhoneDeliverySettingsListDataSourceConfig_NotFound(environmentName, licenseID, resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_phone_delivery_settings_list" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName)
}
