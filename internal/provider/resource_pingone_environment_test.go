package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccPingOneEnvironment_Full(t *testing.T) {

	resourceName := resourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_environment.%s", resourceName)

	name := resourceName
	description := "Test description"
	environmentType := "SANDBOX"
	region := os.Getenv("PINGONE_REGION")
	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckEnvironment(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPingOneEnvironmentFull(resourceName, name, description, environmentType, region, licenseID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", description),
					resource.TestCheckResourceAttr(resourceFullName, "type", environmentType),
					resource.TestCheckResourceAttr(resourceFullName, "region", region),
					resource.TestCheckResourceAttr(resourceFullName, "license_id", licenseID),
				),
			},
		},
	})
}

func TestAccPingOneEnvironment_Minimal(t *testing.T) {

	resourceName := resourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_environment.%s", resourceName)

	name := resourceName
	environmentType := "SANDBOX"
	licenseID := os.Getenv("PINGONE_LICENSE_ID")
	region := os.Getenv("PINGONE_REGION")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckEnvironment(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPingOneEnvironmentMinimal(resourceName, name, environmentType, region, licenseID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "type", environmentType),
					resource.TestCheckResourceAttr(resourceFullName, "region", region),
					resource.TestCheckResourceAttr(resourceFullName, "license_id", licenseID),
				),
			},
		},
	})
}

func testAccPingOneEnvironmentFull(resourceName, name, description, environmentType, region, licenseID string) string {
	return fmt.Sprintf(`
		resource "pingone_environment" "%[1]s" {
			name = "%[2]s"
			description = "%[3]s"
			type = "%[4]s"
			region = "%[5]s"
			license_id = "%[6]s"
		}`, resourceName, name, description, environmentType, region, licenseID)
}

func testAccPingOneEnvironmentMinimal(resourceName, name, environmentType, region, licenseID string) string {
	return fmt.Sprintf(`
		resource "pingone_environment" "%[1]s" {
			name = "%[2]s"
			type = "%[3]s"
			region = "%[4]s"
			license_id = "%[5]s"
		}`, resourceName, name, environmentType, region, licenseID)
}
