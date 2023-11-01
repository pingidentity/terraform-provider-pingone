package base_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/base"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccEnvironmentsDataSource_BySCIMFilter(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullName := fmt.Sprintf("data.pingone_environments.%s", resourceName)

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Environment_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironmentsDataSourceConfig_BySCIMFilter(resourceName, name, licenseID, fmt.Sprintf(`(name sw \"%s\")`, name)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "ids.#", "3"),
					resource.TestMatchResourceAttr(dataSourceFullName, "ids.0", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "ids.1", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "ids.2", verify.P1ResourceIDRegexpFullString),
				),
			},
		},
	})
}

func TestAccEnvironmentsDataSource_NotFound(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullName := fmt.Sprintf("data.pingone_environments.%s", resourceName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Environment_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironmentsDataSourceConfig_NotFound(resourceName, fmt.Sprintf(`(organization.id eq \"%s\")`, uuid.New().String())),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "ids.#", "0"),
				),
			},
		},
	})
}

func testAccEnvironmentsDataSourceConfig_BySCIMFilter(resourceName, name, licenseID, filter string) string {
	return fmt.Sprintf(`


resource "pingone_environment" "%[1]s-1" {
  name       = "%[2]s-1"
  type       = "SANDBOX"
  license_id = "%[3]s"

  service {
  }
}

resource "pingone_environment" "%[1]s-2" {
  name       = "%[2]s-2"
  type       = "SANDBOX"
  license_id = "%[3]s"

  service {
  }
}

resource "pingone_environment" "%[1]s-3" {
  name       = "%[2]s-3"
  type       = "SANDBOX"
  license_id = "%[3]s"

  service {
  }
}

data "pingone_environments" "%[1]s" {

  scim_filter = "%[4]s"

  depends_on = [
    pingone_environment.%[1]s-1,
    pingone_environment.%[1]s-2,
    pingone_environment.%[1]s-3,
  ]
}
`, resourceName, name, licenseID, filter)
}

func testAccEnvironmentsDataSourceConfig_NotFound(resourceName, filter string) string {
	return fmt.Sprintf(`
data "pingone_environments" "%[1]s" {
  scim_filter = "%[2]s"
}
`, resourceName, filter)
}
