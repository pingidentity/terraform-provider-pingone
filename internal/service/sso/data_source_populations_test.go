package sso_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/sso"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccPopulationsDataSource_BySCIMFilter(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullName := fmt.Sprintf("data.pingone_populations.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.TestAccCheckPopulationDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPopulationsDataSourceConfig_BySCIMFilter(resourceName, fmt.Sprintf(`(name eq \"%s-1\") OR (name eq \"%s-2\")`, name, name), name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "ids.#", "2"),
					resource.TestMatchResourceAttr(dataSourceFullName, "ids.0", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "ids.1", verify.P1ResourceIDRegexpFullString),
				),
			},
		},
	})
}

func TestAccPopulationsDataSource_ByDataFilter(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullName := fmt.Sprintf("data.pingone_populations.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.TestAccCheckPopulationDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPopulationsDataSourceConfig_ByDataFilter1(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "ids.#", "2"),
					resource.TestMatchResourceAttr(dataSourceFullName, "ids.0", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "ids.1", verify.P1ResourceIDRegexpFullString),
				),
			},
			{
				Config: testAccPopulationsDataSourceConfig_ByDataFilter2(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "ids.#", "3"),
					resource.TestMatchResourceAttr(dataSourceFullName, "ids.0", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "ids.1", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "ids.2", verify.P1ResourceIDRegexpFullString),
				),
			},
		},
	})
}

func TestAccPopulationsDataSource_NotFound(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullName := fmt.Sprintf("data.pingone_populations.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.TestAccCheckPopulationDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPopulationsDataSourceConfig_NotFound(resourceName, fmt.Sprintf(`(name eq \"%s-1\") OR (name eq \"%s-2\")`, name, name)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "ids.#", "0"),
				),
			},
		},
	})
}

func testAccPopulationsDataSourceConfig_BySCIMFilter(resourceName, filter, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_population" "%[2]s-1" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-1"
}

resource "pingone_population" "%[2]s-2" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-2"
}

resource "pingone_population" "%[2]s-3" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-3"
}

data "pingone_populations" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  scim_filter = "%[4]s"

  depends_on = [
    pingone_population.%[2]s-1,
    pingone_population.%[2]s-2,
    pingone_population.%[2]s-3,
  ]
}
`, acctest.GenericSandboxEnvironment(), resourceName, name, filter)
}

func testAccPopulationsDataSourceConfig_ByDataFilter1(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_population" "%[2]s-1" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-1"
}

resource "pingone_population" "%[2]s-2" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-2"
}

resource "pingone_population" "%[2]s-3" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-3"
}

data "pingone_populations" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  data_filter {
    name   = "name"
    values = ["%[3]s-1", "%[3]s-2"]
  }

  depends_on = [
    pingone_population.%[2]s-1,
    pingone_population.%[2]s-2,
    pingone_population.%[2]s-3,
  ]
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccPopulationsDataSourceConfig_ByDataFilter2(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_population" "%[2]s-1" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-1"
}

resource "pingone_population" "%[2]s-2" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-2"
}

resource "pingone_population" "%[2]s-3" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-3"
}

data "pingone_populations" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  data_filter {
    name   = "name"
    values = ["%[3]s-1", "%[3]s-2", "%[3]s-3", ]
  }

  depends_on = [
    pingone_population.%[2]s-1,
    pingone_population.%[2]s-2,
    pingone_population.%[2]s-3,
  ]
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccPopulationsDataSourceConfig_NotFound(resourceName, filter string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_populations" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  scim_filter = "%[3]s"
}
`, acctest.GenericSandboxEnvironment(), resourceName, filter)
}
