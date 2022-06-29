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

	solution := "CUSTOMER"

	populationName := resourceNameGen()
	populationDescription := "Test population"

	serviceOneType := "SSO"
	serviceTwoType := "PING_FEDERATE"
	serviceTwoURL := "https://my-console-url"
	serviceTwoBookmarkNameOne := "Bookmark 1"
	serviceTwoBookmarkURLOne := "https://my-bookmark-1"
	serviceTwoBookmarkNameTwo := "Bookmark 2"
	serviceTwoBookmarkURLTwo := "https://my-bookmark-2"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckEnvironment(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPingOneEnvironmentFull(resourceName, name, description, environmentType, region, licenseID, solution, populationName, populationDescription, serviceOneType, serviceTwoType, serviceTwoURL, serviceTwoBookmarkNameOne, serviceTwoBookmarkURLOne, serviceTwoBookmarkNameTwo, serviceTwoBookmarkURLTwo),
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

func TestAccPingOneEnvironment_Minimal(t *testing.T) {

	resourceName := resourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_environment.%s", resourceName)

	name := resourceName
	environmentType := "SANDBOX"
	licenseID := os.Getenv("PINGONE_LICENSE_ID")
	region := os.Getenv("PINGONE_REGION")

	resource.Test(t, resource.TestCase{
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
					resource.TestCheckResourceAttr(resourceFullName, "default_population.0.name", "Default"),
					resource.TestCheckResourceAttr(resourceFullName, "service.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "service.0.type", "SSO"),
				),
			},
		},
	})
}

func testAccPingOneEnvironmentFull(resourceName, name, description, environmentType, region, licenseID, solution, populationName, populationDescription, serviceOneType, serviceTwoType, serviceTwoURL, serviceTwoBookmarkNameOne, serviceTwoBookmarkURLOne, serviceTwoBookmarkNameTwo, serviceTwoBookmarkURLTwo string) string {
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

func testAccPingOneEnvironmentMinimal(resourceName, name, environmentType, region, licenseID string) string {
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

func TestBuildBOMProductsCreateRequest(t *testing.T) {
	t.Fatalf("Not implemented")
}

func TestFlattenBOMProducts(t *testing.T) {
	t.Fatalf("Not implemented")
}

func TestFlattenBOMProductsBookmarkList(t *testing.T) {
	t.Fatalf("Not implemented")
}
