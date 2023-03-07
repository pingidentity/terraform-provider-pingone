package base_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
)

func TestAccEnvironmentsDataSource_BySCIMFilter(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullName := fmt.Sprintf("data.pingone_environments.%s", resourceName)

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEnvironmentDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironmentsDataSourceConfig_BySCIMFilter(resourceName, name, licenseID, fmt.Sprintf(`(name sw \"%s\")`, name)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(dataSourceFullName, "ids.#", "3"),
					resource.TestMatchResourceAttr(dataSourceFullName, "ids.0", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "ids.1", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "ids.2", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
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
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEnvironmentDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironmentsDataSourceConfig_NotFound(resourceName, fmt.Sprintf(`(organization.id eq \"%s\")`, uuid.New().String())),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
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
  default_population {
  }
  service {
  }
}

resource "pingone_environment" "%[1]s-2" {
  name       = "%[2]s-2"
  type       = "SANDBOX"
  license_id = "%[3]s"
  default_population {
  }
  service {
  }
}

resource "pingone_environment" "%[1]s-3" {
  name       = "%[2]s-3"
  type       = "SANDBOX"
  license_id = "%[3]s"
  default_population {
  }
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
