package sso_test

import (
	"fmt"
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
					resource.TestCheckResourceAttrSet(dataSourceFullName, "id"),
					resource.TestCheckResourceAttr(dataSourceFullName, "name", "User"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "description"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "schema_id"),
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
					resource.TestCheckResourceAttrSet(dataSourceFullName, "id"),
					resource.TestCheckResourceAttr(dataSourceFullName, "name", "User"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "description"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "schema_id"),
				),
			},
		},
	})
}

func testAccSchemaDataSourceConfig_ByNameFull(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s
		
		data "pingone_schema" "%[2]s" {
			environment_id = "${data.pingone_environment.general_test.id}"

			name = "%[3]s"
		}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccSchemaDataSourceConfig_ByIDFull(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

	data "pingone_schema" "%[2]s" {
		environment_id = "${data.pingone_environment.general_test.id}"

		name = "%[3]s"
	}

	data "pingone_schema" "%[2]s-2" {
		environment_id = "${data.pingone_environment.general_test.id}"

		schema_id = "${data.pingone_schema.%[2]s.id}"
	}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
