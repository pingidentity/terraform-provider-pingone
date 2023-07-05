package sso_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccUserDataSource_ByNameFull(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_user.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccUserDataSourceConfig_ByNameFull(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(dataSourceFullName, "user_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(dataSourceFullName, "username", name),
					resource.TestCheckResourceAttr(dataSourceFullName, "email", fmt.Sprintf("%s@pingidentity.com", name)),
					resource.TestCheckResourceAttr(dataSourceFullName, "status", "ENABLED"),
					resource.TestMatchResourceAttr(dataSourceFullName, "population_id", verify.P1ResourceIDRegexp),
				),
			},
		},
	})
}

func TestAccUserDataSource_ByEmailFull(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_user.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccUserDataSourceConfig_ByEmailFull(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(dataSourceFullName, "user_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(dataSourceFullName, "username", name),
					resource.TestCheckResourceAttr(dataSourceFullName, "email", fmt.Sprintf("%s@pingidentity.com", name)),
					resource.TestCheckResourceAttr(dataSourceFullName, "status", "ENABLED"),
					resource.TestMatchResourceAttr(dataSourceFullName, "population_id", verify.P1ResourceIDRegexp),
				),
			},
		},
	})
}

func TestAccUserDataSource_ByIDFull(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_user.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccUserDataSourceConfig_ByIDFull(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(dataSourceFullName, "user_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(dataSourceFullName, "username", name),
					resource.TestCheckResourceAttr(dataSourceFullName, "email", fmt.Sprintf("%s@pingidentity.com", name)),
					resource.TestCheckResourceAttr(dataSourceFullName, "status", "ENABLED"),
					resource.TestMatchResourceAttr(dataSourceFullName, "population_id", verify.P1ResourceIDRegexp),
				),
			},
		},
	})
}

func TestAccUserDataSource_NotFound(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccUserDataSourceConfig_NotFoundByName(resourceName),
				ExpectError: regexp.MustCompile("Cannot find user"),
			},
			{
				Config:      testAccUserDataSourceConfig_NotFoundByID(resourceName),
				ExpectError: regexp.MustCompile("Cannot find user"),
			},
		},
	})
}

func testAccUserDataSourceConfig_ByNameFull(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_population" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_user" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  username      = "%[3]s"
  email         = "%[3]s@pingidentity.com"
  population_id = pingone_population.%[2]s.id
}

data "pingone_user" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  username = "%[3]s"

  depends_on = [
    pingone_user.%[2]s,
  ]
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccUserDataSourceConfig_ByEmailFull(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_population" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_user" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  username      = "%[3]s"
  email         = "%[3]s@pingidentity.com"
  population_id = pingone_population.%[2]s.id
}

data "pingone_user" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  email = "%[3]s@pingidentity.com"

  depends_on = [
    pingone_user.%[2]s,
  ]
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccUserDataSourceConfig_ByIDFull(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_population" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_user" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  username      = "%[3]s"
  email         = "%[3]s@pingidentity.com"
  population_id = pingone_population.%[2]s.id
}

data "pingone_user" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  user_id = pingone_user.%[2]s.id
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccUserDataSourceConfig_NotFoundByName(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_user" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  username = "doesnotexist"
}`, acctest.GenericSandboxEnvironment(), resourceName)
}

func testAccUserDataSourceConfig_NotFoundByID(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_user" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  user_id = "9c052a8a-14be-44e4-8f07-2662569994ce" // dummy ID that conforms to UUID v4
}`, acctest.GenericSandboxEnvironment(), resourceName)
}
