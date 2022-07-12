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

	region := os.Getenv("PINGONE_REGION")
	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      acctest.TestAccCheckEnvironmentDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSchemaDataSourceConfig_ByNameFull(environmentName, resourceName, "User", region, licenseID),
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

	region := os.Getenv("PINGONE_REGION")
	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      acctest.TestAccCheckEnvironmentDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSchemaDataSourceConfig_ByIDFull(environmentName, resourceName, "User", region, licenseID),
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

func testAccSchemaDataSourceConfig_ByNameFull(environmentName, resourceName, name, region, licenseID string) string {
	return fmt.Sprintf(`
		resource "pingone_environment" "%[1]s" {
			name = "%[1]s"
			type = "SANDBOX"
			license_id = "%[5]s"
			region = "%[4]s"
			default_population {}
			service {}
		}
		data "pingone_schema" "%[2]s" {
			environment_id = "${pingone_environment.%[1]s.id}"

			name = "%[3]s"
		}`, environmentName, resourceName, name, region, licenseID)
}

func testAccSchemaDataSourceConfig_ByIDFull(environmentName, resourceName, name, region, licenseID string) string {
	return fmt.Sprintf(`
	resource "pingone_environment" "%[1]s" {
		name = "%[1]s"
		type = "SANDBOX"
		license_id = "%[5]s"
		region = "%[4]s"
		default_population {}
		service {}
	}
	data "pingone_schema" "%[1]s" {
		environment_id = "${pingone_environment.%[1]s.id}"

		name = "%[3]s"
	}
	data "pingone_schema" "%[2]s" {
		environment_id = "${pingone_environment.%[1]s.id}"

		schema_id = "${data.pingone_schema.%[1]s.id}"
	}`, environmentName, resourceName, name, region, licenseID)
}
