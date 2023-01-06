package sso_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
)

func TestAccSchemaDataSource_ByNameFull(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_schema.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSchemaDataSourceConfig_ByNameFull(resourceName, "User"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "schema_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(dataSourceFullName, "name", "User"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "description"),
				),
			},
		},
	})
}

func TestAccSchemaDataSource_ByIDFull(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_schema.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSchemaDataSourceConfig_ByIDFull(resourceName, "User"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "schema_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(dataSourceFullName, "name", "User"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "description"),
				),
			},
		},
	})
}

func TestAccSchemaDataSource_NotFound(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccSchemaDataSourceConfig_NotFoundByName(resourceName),
				ExpectError: regexp.MustCompile("Cannot find schema doesnotexist"),
			},
			{
				Config:      testAccSchemaDataSourceConfig_NotFoundByID(resourceName),
				ExpectError: regexp.MustCompile("Error when calling `ReadOneSchema`: The request could not be completed. The requested resource was not found."),
			},
		},
	})
}

func testAccSchemaDataSourceConfig_ByNameFull(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

data "pingone_schema" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccSchemaDataSourceConfig_ByIDFull(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_schema" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

data "pingone_schema" "%[2]s-2" {
  environment_id = data.pingone_environment.general_test.id

  schema_id = data.pingone_schema.%[2]s.id
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccSchemaDataSourceConfig_NotFoundByName(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_schema" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "doesnotexist"
}
`, acctest.GenericSandboxEnvironment(), resourceName)
}

func testAccSchemaDataSourceConfig_NotFoundByID(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_schema" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  schema_id = "9c052a8a-14be-44e4-8f07-2662569994ce" // dummy ID that conforms to UUID v4
}
`, acctest.GenericSandboxEnvironment(), resourceName)
}
