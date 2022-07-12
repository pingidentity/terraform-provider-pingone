package base_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
)

func TestAccRoleDataSource_ByNameFull(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_role.%s", resourceName)
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
				Config: testAccRoleDataSourceConfig_ByNameFull(environmentName, resourceName, "Organization Admin", region, licenseID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceFullName, "id"),
					resource.TestCheckResourceAttr(dataSourceFullName, "name", "Organization Admin"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "description"),
				),
			},
			{
				Config: testAccRoleDataSourceConfig_ByNameFull(environmentName, resourceName, "Environment Admin", region, licenseID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceFullName, "id"),
					resource.TestCheckResourceAttr(dataSourceFullName, "name", "Environment Admin"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "description"),
				),
			},
			{
				Config: testAccRoleDataSourceConfig_ByNameFull(environmentName, resourceName, "Identity Data Admin", region, licenseID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceFullName, "id"),
					resource.TestCheckResourceAttr(dataSourceFullName, "name", "Identity Data Admin"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "description"),
				),
			},
			{
				Config: testAccRoleDataSourceConfig_ByNameFull(environmentName, resourceName, "Client Application Developer", region, licenseID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceFullName, "id"),
					resource.TestCheckResourceAttr(dataSourceFullName, "name", "Client Application Developer"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "description"),
				),
			},
			{
				Config: testAccRoleDataSourceConfig_ByNameFull(environmentName, resourceName, "Identity Data Read Only", region, licenseID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceFullName, "id"),
					resource.TestCheckResourceAttr(dataSourceFullName, "name", "Identity Data Read Only"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "description"),
				),
			},
			{
				Config: testAccRoleDataSourceConfig_ByNameFull(environmentName, resourceName, "Configuration Read Only", region, licenseID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceFullName, "id"),
					resource.TestCheckResourceAttr(dataSourceFullName, "name", "Configuration Read Only"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "description"),
				),
			},
		},
	})
}

func testAccRoleDataSourceConfig_ByNameFull(environmentName, resourceName, name, region, licenseID string) string {
	return fmt.Sprintf(`
		resource "pingone_environment" "%[1]s" {
			name = "%[1]s"
			type = "SANDBOX"
			license_id = "%[5]s"
			region = "%[4]s"
			default_population {}
			service {}
		}

		data "pingone_role" "%[2]s" {
			name = "%[3]s"
		}`, environmentName, resourceName, name, region, licenseID)
}
