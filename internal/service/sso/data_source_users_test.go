package sso_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
)

func TestAccUsersDataSource_BySCIMFilter(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullName := fmt.Sprintf("data.pingone_users.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckUserDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccUsersDataSourceConfig_BySCIMFilter(resourceName, fmt.Sprintf(`(username eq \"%s-1\") OR (username eq \"%s-2\")`, name, name), name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(dataSourceFullName, "ids.#", "2"),
					resource.TestMatchResourceAttr(dataSourceFullName, "ids.0", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "ids.1", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
				),
			},
		},
	})
}

func TestAccUsersDataSource_ByDataFilter(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullName := fmt.Sprintf("data.pingone_users.%s", resourceName)

	organizationID := os.Getenv("PINGONE_ORGANIZATION_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckUserDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccUsersDataSourceConfig_ByDataFilter1(resourceName, organizationID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(dataSourceFullName, "ids.#", "2"),
					resource.TestMatchResourceAttr(dataSourceFullName, "ids.0", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "ids.1", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
				),
			},
			{
				Config: testAccUsersDataSourceConfig_ByDataFilter2(resourceName, organizationID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(dataSourceFullName, "ids.#", "3"),
					resource.TestMatchResourceAttr(dataSourceFullName, "ids.0", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "ids.1", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "ids.2", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
				),
			},
		},
	})
}

func testAccUsersDataSourceConfig_BySCIMFilter(resourceName, filter, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_population" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_user" "%[2]s-1" {
  environment_id = data.pingone_environment.general_test.id

  username      = "%[3]s-1"
  email         = "noreply@pingidentity.com"
  population_id = pingone_population.%[2]s.id
}

resource "pingone_user" "%[2]s-2" {
  environment_id = data.pingone_environment.general_test.id

  username      = "%[3]s-2"
  email         = "noreply@pingidentity.com"
  population_id = pingone_population.%[2]s.id
}

resource "pingone_user" "%[2]s-3" {
  environment_id = data.pingone_environment.general_test.id

  username      = "%[3]s-3"
  email         = "noreply@pingidentity.com"
  population_id = pingone_population.%[2]s.id
}

data "pingone_users" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  scim_filter = "%[4]s"

  depends_on = [
    pingone_user.%[2]s-1,
    pingone_user.%[2]s-2,
    pingone_user.%[2]s-3,
    pingone_population.%[2]s,
  ]
}`, acctest.GenericSandboxEnvironment(), resourceName, name, filter)
}

func testAccUsersDataSourceConfig_ByDataFilter1(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_population" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_user" "%[2]s-1" {
  environment_id = data.pingone_environment.general_test.id

  username      = "%[3]s-1"
  email         = "noreply@pingidentity.com"
  population_id = pingone_population.%[2]s.id
}

resource "pingone_user" "%[2]s-2" {
  environment_id = data.pingone_environment.general_test.id

  username      = "%[3]s-2"
  email         = "noreply@pingidentity.com"
  population_id = pingone_population.%[2]s.id
}

resource "pingone_user" "%[2]s-3" {
  environment_id = data.pingone_environment.general_test.id

  username      = "%[3]s-3"
  email         = "noreply@pingidentity.com"
  population_id = pingone_population.%[2]s.id
}

data "pingone_users" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  data_filter {
    name   = "username"
    values = ["%[3]s-1", "%[3]s-2"]
  }

  data_filter {
    name   = "population.id"
    values = [pingone_population.%[2]s.id]
  }

  depends_on = [
    pingone_user.%[2]s-1,
    pingone_user.%[2]s-2,
    pingone_user.%[2]s-3,
    pingone_population.%[2]s,
  ]
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccUsersDataSourceConfig_ByDataFilter2(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_population" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_user" "%[2]s-1" {
  environment_id = data.pingone_environment.general_test.id

  username      = "%[3]s-1"
  email         = "noreply@pingidentity.com"
  population_id = pingone_population.%[2]s.id
}

resource "pingone_user" "%[2]s-2" {
  environment_id = data.pingone_environment.general_test.id

  username      = "%[3]s-2"
  email         = "noreply@pingidentity.com"
  population_id = pingone_population.%[2]s.id
}

resource "pingone_user" "%[2]s-3" {
  environment_id = data.pingone_environment.general_test.id

  username      = "%[3]s-3"
  email         = "noreply@pingidentity.com"
  population_id = pingone_population.%[2]s.id
}

data "pingone_users" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  data_filter {
    name   = "population.id"
    values = [pingone_population.%[2]s.id]
  }

  depends_on = [
    pingone_user.%[2]s-1,
    pingone_user.%[2]s-2,
    pingone_user.%[2]s-3,
    pingone_population.%[2]s,
  ]
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
