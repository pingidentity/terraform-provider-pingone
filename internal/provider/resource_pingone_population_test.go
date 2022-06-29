package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccPingOnePopulation_Full(t *testing.T) {

	resourceName := resourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_population.%s", resourceName)

	name := resourceName
	description := "Test description"

	licenseID := os.Getenv("PINGONE_LICENSE_ID")
	region := os.Getenv("PINGONE_REGION")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckEnvironment(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPingOnePopulationFull(resourceName, name, description, licenseID, region),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", description),
				),
			},
		},
	})
}

func TestAccPingOnePopulation_Minimal(t *testing.T) {

	resourceName := resourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_population.%s", resourceName)

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")
	region := os.Getenv("PINGONE_REGION")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckEnvironment(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPingOnePopulationMinimal(resourceName, name, licenseID, region),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
				),
			},
		},
	})
}

func testAccPingOnePopulationFull(resourceName, name, description, licenseID, region string) string {
	return fmt.Sprintf(`
		resource "pingone_environment" "%[1]s" {
			name = "%[2]s"
			type = "SANDBOX"
			license_id = "%[4]s"
			region = "%[5]s"
			default_population {}
			service {}
		}

		resource "pingone_population" "%[1]s" {
			environment_id = "${pingone_environment.%[1]s.id}"
			name = "%[2]s"
			description = "%[3]s"
		}`, resourceName, name, description, licenseID, region)
}

func testAccPingOnePopulationMinimal(resourceName, name, licenseID, region string) string {
	return fmt.Sprintf(`
		resource "pingone_environment" "%[1]s" {
			name = "%[2]s"
			type = "SANDBOX"
			license_id = "%[3]s"
			region = "%[4]s"
			default_population {}
			service {}
		}

		resource "pingone_population" "%[1]s" {
			environment_id = "${pingone_environment.%[1]s.id}"
			name = "%[2]s"
		}`, resourceName, name, licenseID, region)
}
