package sso_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
)

func TestAccSchemaDataSource_ByNameFull(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_schema.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      acctest.TestAccCheckEnvironmentDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSchemaDataSourceConfig_ByNameFull(environmentName, licenseID, resourceName, "User"),
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

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      acctest.TestAccCheckEnvironmentDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSchemaDataSourceConfig_ByIDFull(environmentName, licenseID, resourceName, "User"),
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

func testAccSchemaDataSourceConfig_ByNameFull(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s
		
		data "pingone_schema" "%[3]s" {
			environment_id = "${pingone_environment.%[2]s.id}"

			name = "%[4]s"
		}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccSchemaDataSourceConfig_ByIDFull(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

	data "pingone_schema" "%[2]s" {
		environment_id = "${pingone_environment.%[2]s.id}"

		name = "%[4]s"
	}

	data "pingone_schema" "%[3]s" {
		environment_id = "${pingone_environment.%[2]s.id}"

		schema_id = "${data.pingone_schema.%[2]s.id}"
	}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}
