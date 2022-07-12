package base_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
)

func TestAccEnvironment_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGenEnvironment()
	resourceFullName := fmt.Sprintf("pingone_environment.%s", resourceName)

	name := resourceName
	description := "Test description"
	environmentType := "SANDBOX"
	region := os.Getenv("PINGONE_REGION")
	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	solution := "CUSTOMER"

	populationName := acctest.ResourceNameGenDefaultPopulation()
	populationDescription := "Test population"

	serviceOneType := "SSO"
	serviceTwoType := "PingFederate"
	serviceTwoURL := "https://my-console-url"
	serviceTwoBookmarkNameOne := "Bookmark 1"
	serviceTwoBookmarkURLOne := "https://my-bookmark-1"
	serviceTwoBookmarkNameTwo := "Bookmark 2"
	serviceTwoBookmarkURLTwo := "https://my-bookmark-2"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      acctest.TestAccCheckEnvironmentDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironmentConfig_Full(resourceName, name, description, environmentType, region, licenseID, solution, populationName, populationDescription, serviceOneType, serviceTwoType, serviceTwoURL, serviceTwoBookmarkNameOne, serviceTwoBookmarkURLOne, serviceTwoBookmarkNameTwo, serviceTwoBookmarkURLTwo),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", description),
					resource.TestCheckResourceAttr(resourceFullName, "type", environmentType),
					resource.TestCheckResourceAttr(resourceFullName, "region", region),
					resource.TestCheckResourceAttr(resourceFullName, "license_id", licenseID),
					// resource.TestCheckResourceAttr(resourceFullName, "solution", solution),
					resource.TestCheckResourceAttr(resourceFullName, "default_population.0.name", populationName),
					resource.TestCheckResourceAttr(resourceFullName, "default_population.0.description", populationDescription),
					resource.TestCheckResourceAttr(resourceFullName, "service.#", "2"),
					resource.TestCheckResourceAttr(resourceFullName, "service.0.type", serviceOneType),
					resource.TestCheckResourceAttr(resourceFullName, "service.1.type", serviceTwoType),
					resource.TestCheckResourceAttr(resourceFullName, "service.1.console_url", serviceTwoURL),
					resource.TestCheckResourceAttr(resourceFullName, "service.1.bookmark.#", "2"),
					resource.TestCheckResourceAttr(resourceFullName, "service.1.bookmark.0.name", serviceTwoBookmarkNameOne),
					resource.TestCheckResourceAttr(resourceFullName, "service.1.bookmark.0.url", serviceTwoBookmarkURLOne),
					resource.TestCheckResourceAttr(resourceFullName, "service.1.bookmark.1.name", serviceTwoBookmarkNameTwo),
					resource.TestCheckResourceAttr(resourceFullName, "service.1.bookmark.1.url", serviceTwoBookmarkURLTwo),
				),
			},
		},
	})
}

func TestAccEnvironment_Minimal(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGenEnvironment()
	resourceFullName := fmt.Sprintf("pingone_environment.%s", resourceName)

	name := resourceName
	environmentType := "SANDBOX"
	licenseID := os.Getenv("PINGONE_LICENSE_ID")
	region := os.Getenv("PINGONE_REGION")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      acctest.TestAccCheckEnvironmentDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironmentConfig_Minimal(resourceName, name, environmentType, region, licenseID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "type", environmentType),
					resource.TestCheckResourceAttr(resourceFullName, "region", region),
					resource.TestCheckResourceAttr(resourceFullName, "license_id", licenseID),
					resource.TestCheckResourceAttr(resourceFullName, "default_population.0.name", "Default"),
					resource.TestCheckResourceAttr(resourceFullName, "service.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "service.0.type", "SSO"),
				),
			},
		},
	})
}

func TestAccEnvironment_NonPopulationServices(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGenEnvironment()
	resourceFullName := fmt.Sprintf("pingone_environment.%s", resourceName)

	name := resourceName
	description := "Test description"
	environmentType := "SANDBOX"
	region := os.Getenv("PINGONE_REGION")
	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	solution := "CUSTOMER"

	populationName := acctest.ResourceNameGenDefaultPopulation()
	populationDescription := "Test population"

	serviceOneType := "PingAccess"
	serviceTwoType := "PingFederate"
	serviceTwoURL := "https://my-console-url"
	serviceTwoBookmarkNameOne := "Bookmark 1"
	serviceTwoBookmarkURLOne := "https://my-bookmark-1"
	serviceTwoBookmarkNameTwo := "Bookmark 2"
	serviceTwoBookmarkURLTwo := "https://my-bookmark-2"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      acctest.TestAccCheckEnvironmentDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironmentConfig_Full(resourceName, name, description, environmentType, region, licenseID, solution, populationName, populationDescription, serviceOneType, serviceTwoType, serviceTwoURL, serviceTwoBookmarkNameOne, serviceTwoBookmarkURLOne, serviceTwoBookmarkNameTwo, serviceTwoBookmarkURLTwo),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", description),
					resource.TestCheckResourceAttr(resourceFullName, "type", environmentType),
					resource.TestCheckResourceAttr(resourceFullName, "region", region),
					resource.TestCheckResourceAttr(resourceFullName, "license_id", licenseID),
					// resource.TestCheckResourceAttr(resourceFullName, "solution", solution),
					resource.TestCheckResourceAttr(resourceFullName, "default_population.0.name", populationName),
					resource.TestCheckResourceAttr(resourceFullName, "default_population.0.description", populationDescription),
					resource.TestCheckResourceAttr(resourceFullName, "service.#", "2"),
					resource.TestCheckResourceAttr(resourceFullName, "service.0.type", serviceOneType),
					resource.TestCheckResourceAttr(resourceFullName, "service.1.type", serviceTwoType),
					resource.TestCheckResourceAttr(resourceFullName, "service.1.console_url", serviceTwoURL),
					resource.TestCheckResourceAttr(resourceFullName, "service.1.bookmark.#", "2"),
					resource.TestCheckResourceAttr(resourceFullName, "service.1.bookmark.0.name", serviceTwoBookmarkNameOne),
					resource.TestCheckResourceAttr(resourceFullName, "service.1.bookmark.0.url", serviceTwoBookmarkURLOne),
					resource.TestCheckResourceAttr(resourceFullName, "service.1.bookmark.1.name", serviceTwoBookmarkNameTwo),
					resource.TestCheckResourceAttr(resourceFullName, "service.1.bookmark.1.url", serviceTwoBookmarkURLTwo),
				),
			},
		},
	})
}

func testAccEnvironmentConfig_Full(resourceName, name, description, environmentType, region, licenseID, solution, populationName, populationDescription, serviceOneType, serviceTwoType, serviceTwoURL, serviceTwoBookmarkNameOne, serviceTwoBookmarkURLOne, serviceTwoBookmarkNameTwo, serviceTwoBookmarkURLTwo string) string {
	return fmt.Sprintf(`
		resource "pingone_environment" "%[1]s" {
			name = "%[2]s"
			description = "%[3]s"
			type = "%[4]s"
			region = "%[5]s"
			license_id = "%[6]s"
			default_population {
				name = "%[8]s"
				description = "%[9]s"
			}
			service {
				type = "%[10]s"
			}
			service {
				type = "%[11]s"
				console_url = "%[12]s"
				bookmark {
					name = "%[13]s"
					url = "%[14]s"
				}
				bookmark {
					name = "%[15]s"
					url = "%[16]s"
				}
			}
		}`, resourceName, name, description, environmentType, region, licenseID, solution, populationName, populationDescription, serviceOneType, serviceTwoType, serviceTwoURL, serviceTwoBookmarkNameOne, serviceTwoBookmarkURLOne, serviceTwoBookmarkNameTwo, serviceTwoBookmarkURLTwo)
}

func testAccEnvironmentConfig_Minimal(resourceName, name, environmentType, region, licenseID string) string {
	return fmt.Sprintf(`
		resource "pingone_environment" "%[1]s" {
			name = "%[2]s"
			type = "%[3]s"
			region = "%[4]s"
			license_id = "%[5]s"
			default_population {
			}
			service {
			}
		}`, resourceName, name, environmentType, region, licenseID)
}
