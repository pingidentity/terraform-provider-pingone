package sso_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccUsersDataSource_BySCIMFilter(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullName := fmt.Sprintf("data.pingone_users.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckUserDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccUsersDataSourceConfig_BySCIMFilter(resourceName, fmt.Sprintf(`(username eq \"%s-1\") OR (username eq \"%s-2\")`, name, name), name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(dataSourceFullName, "ids.#", "2"),
					resource.TestMatchResourceAttr(dataSourceFullName, "ids.0", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(dataSourceFullName, "ids.1", verify.P1ResourceIDRegexp),
				),
			},
		},
	})
}

func TestAccUsersDataSource_ByDataFilter(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullName := fmt.Sprintf("data.pingone_users.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckUserDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccUsersDataSourceConfig_ByDataFilter1(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(dataSourceFullName, "ids.#", "2"),
					resource.TestMatchResourceAttr(dataSourceFullName, "ids.0", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(dataSourceFullName, "ids.1", verify.P1ResourceIDRegexp),
				),
			},
			{
				Config: testAccUsersDataSourceConfig_ByDataFilter2(resourceName, name),
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

func TestAccUsersDataSource_NotFound(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullName := fmt.Sprintf("data.pingone_users.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckUserDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccUsersDataSourceConfig_NotFound(resourceName, fmt.Sprintf(`(username eq \"%s-1\") OR (username eq \"%s-2\")`, name, name), name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(dataSourceFullName, "ids.#", "0"),
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

func testAccUsersDataSourceConfig_NotFound(resourceName, filter, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_population" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

data "pingone_users" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  scim_filter = "%[4]s"
}`, acctest.GenericSandboxEnvironment(), resourceName, name, filter)
}
