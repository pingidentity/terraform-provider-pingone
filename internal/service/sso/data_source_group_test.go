// Copyright © 2026 Ping Identity Corporation

package sso_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/sso"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccGroupDataSource_ByNameFull(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_group.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Group_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGroupDataSourceConfig_ByNameFull(resourceName, name, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "environment_id", resourceFullName, "environment_id"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "group_id", resourceFullName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "name", resourceFullName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "description", resourceFullName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "population_id", resourceFullName, "population_id"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "user_filter", resourceFullName, "user_filter"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "external_id", resourceFullName, "external_id"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "custom_data", resourceFullName, "custom_data"),
				),
			},
			{
				Config: testAccGroupDataSourceConfig_ByNameFull(resourceName, name, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "environment_id", resourceFullName, "environment_id"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "group_id", resourceFullName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "name", resourceFullName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "description", resourceFullName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "population_id", resourceFullName, "population_id"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "user_filter", resourceFullName, "user_filter"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "external_id", resourceFullName, "external_id"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "custom_data", resourceFullName, "custom_data"),
				),
			},
		},
	})
}

func TestAccGroupDataSource_ByIDFull(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_group.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Group_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGroupDataSourceConfig_ByIDFull(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "environment_id", resourceFullName, "environment_id"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "group_id", resourceFullName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "name", resourceFullName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "description", resourceFullName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "population_id", resourceFullName, "population_id"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "user_filter", resourceFullName, "user_filter"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "external_id", resourceFullName, "external_id"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "custom_data", resourceFullName, "custom_data"),
				),
			},
		},
	})
}

func TestAccGroupDataSource_NotFound(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Group_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccGroupDataSourceConfig_NotFoundByName(resourceName),
				ExpectError: regexp.MustCompile(`Group not found`),
			},
			{
				Config:      testAccGroupDataSourceConfig_NotFoundByID(resourceName),
				ExpectError: regexp.MustCompile("Group not found"),
			},
		},
	})
}

func testAccGroupDataSourceConfig_ByNameFull(resourceName, name string, insensitivityCheck bool) string {

	// If insensitivityCheck is true, alter the case of the name
	nameComparator := name
	if insensitivityCheck {
		nameComparator = acctest.AlterStringCasing(nameComparator)
	}

	return fmt.Sprintf(`
	%[1]s

resource "pingone_population" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_group" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name          = "%[3]s"
  description   = "Test description"
  population_id = pingone_population.%[2]s.id
  user_filter   = "email ew \"@test.com\""
  external_id   = "external_1234"

  custom_data = jsonencode({ "hello" = "world" })
}

data "pingone_group" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[4]s"

  depends_on = [pingone_group.%[2]s]
}
`, acctest.GenericSandboxEnvironment(), resourceName, name, nameComparator)
}

func testAccGroupDataSourceConfig_ByIDFull(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_population" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_group" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name          = "%[3]s"
  description   = "Test description"
  population_id = pingone_population.%[2]s.id
  user_filter   = "email ew \"@test.com\""
  external_id   = "external_1234"

  custom_data = jsonencode({ "hello" = "world" })
}

data "pingone_group" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  group_id = pingone_group.%[2]s.id
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccGroupDataSourceConfig_NotFoundByName(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_group" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "doesnotexist"
}
`, acctest.GenericSandboxEnvironment(), resourceName)
}

func testAccGroupDataSourceConfig_NotFoundByID(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_group" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  group_id = "9c052a8a-14be-44e4-8f07-2662569994ce" // dummy ID that conforms to UUID v4
}
`, acctest.GenericSandboxEnvironment(), resourceName)
}
