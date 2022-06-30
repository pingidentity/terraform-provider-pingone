package sso_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
)

func TestAccGroup_Full(t *testing.T) {

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_group.%s", resourceName)

	name := resourceName
	description := "Test description"

	licenseID := os.Getenv("PINGONE_LICENSE_ID")
	region := os.Getenv("PINGONE_REGION")

	userFilter := `email ew "@test.com"`
	externalID := "external_1234"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGroupConfig_Full(resourceName, name, description, licenseID, region, userFilter, externalID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", description),
					//resource.TestCheckResourceAttr(resourceFullName, "population_id", populationID),
					resource.TestCheckResourceAttr(resourceFullName, "user_filter", userFilter),
					resource.TestCheckResourceAttr(resourceFullName, "external_id", externalID),
				),
			},
		},
	})
}

func TestAccGroup_Minimal(t *testing.T) {

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_group.%s", resourceName)

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")
	region := os.Getenv("PINGONE_REGION")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGroupConfig_Minimal(resourceName, name, licenseID, region),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
				),
			},
		},
	})
}

func testAccGroupConfig_Full(resourceName, name, description, licenseID, region, userFilter, externalID string) string {
	return fmt.Sprintf(`
		resource "pingone_environment" "%[1]s" {
			name = "%[2]s"
			type = "SANDBOX"
			license_id = "%[4]s"
			region = "%[5]s"
			default_population {}
			service {}
		}

		resource "pingone_group" "%[1]s" {
			environment_id = "${pingone_environment.%[1]s.id}"
			name = "%[2]s"
			description = "%[3]s"
			population_id = "${pingone_environment.%[1]s.default_population_id}"
			user_filter = %[6]q
			external_id = "%[7]s"
		}`, resourceName, name, description, licenseID, region, userFilter, externalID)
}

func testAccGroupConfig_Minimal(resourceName, name, licenseID, region string) string {
	return fmt.Sprintf(`
		resource "pingone_environment" "%[1]s" {
			name = "%[2]s"
			type = "SANDBOX"
			license_id = "%[3]s"
			region = "%[4]s"
			default_population {}
			service {}
		}

		resource "pingone_group" "%[1]s" {
			environment_id = "${pingone_environment.%[1]s.id}"
			name = "%[2]s"
		}`, resourceName, name, licenseID, region)
}
