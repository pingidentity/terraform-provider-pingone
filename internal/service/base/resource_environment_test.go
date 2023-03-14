package base_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func testAccCheckEnvironmentDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_environment" {
			continue
		}

		_, r, err := apiClient.EnvironmentsApi.ReadOneEnvironment(ctx, rs.Primary.ID).Execute()

		if r.StatusCode == 404 {
			continue
		}

		if err != nil {
			return err
		}

		return fmt.Errorf("PingOne Environment Instance %s still exists", rs.Primary.ID)
	}

	return nil
}

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
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEnvironmentDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironmentConfig_Full(resourceName, name, description, environmentType, region, licenseID, solution, populationName, populationDescription, serviceOneType, serviceTwoType, serviceTwoURL, serviceTwoBookmarkNameOne, serviceTwoBookmarkURLOne, serviceTwoBookmarkNameTwo, serviceTwoBookmarkURLTwo),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", description),
					resource.TestCheckResourceAttr(resourceFullName, "type", environmentType),
					resource.TestCheckResourceAttr(resourceFullName, "region", region),
					resource.TestCheckResourceAttr(resourceFullName, "license_id", licenseID),
					resource.TestMatchResourceAttr(resourceFullName, "organization_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "solution", solution),
					resource.TestMatchResourceAttr(resourceFullName, "default_population_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "default_population.0.name", populationName),
					resource.TestCheckResourceAttr(resourceFullName, "default_population.0.description", populationDescription),
					resource.TestCheckResourceAttr(resourceFullName, "service.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "service.*", map[string]string{
						"type":        serviceOneType,
						"console_url": "",
						"bookmark.#":  "0",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "service.*", map[string]string{
						"type":            serviceTwoType,
						"console_url":     serviceTwoURL,
						"bookmark.#":      "2",
						"bookmark.0.name": serviceTwoBookmarkNameOne,
						"bookmark.0.url":  serviceTwoBookmarkURLOne,
						"bookmark.1.name": serviceTwoBookmarkNameTwo,
						"bookmark.1.url":  serviceTwoBookmarkURLTwo,
					}),
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
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEnvironmentDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironmentConfig_Minimal(resourceName, name, environmentType, licenseID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckNoResourceAttr(resourceFullName, "description"),
					resource.TestCheckResourceAttr(resourceFullName, "type", environmentType),
					resource.TestCheckResourceAttr(resourceFullName, "region", region),
					resource.TestCheckNoResourceAttr(resourceFullName, "solution"),
					resource.TestCheckResourceAttr(resourceFullName, "license_id", licenseID),
					resource.TestMatchResourceAttr(resourceFullName, "organization_id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "default_population_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "default_population.0.name", "Default"),
					resource.TestCheckNoResourceAttr(resourceFullName, "default_population.0.description"),
					resource.TestCheckResourceAttr(resourceFullName, "service.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "service.*", map[string]string{
						"type":        "SSO",
						"console_url": "",
						"bookmark.#":  "0",
					}),
				),
			},
		},
	})
}

func TestAccEnvironment_NonCompatibleRegion(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGenEnvironment()

	name := resourceName
	environmentType := "SANDBOX"
	licenseID := os.Getenv("PINGONE_LICENSE_ID")
	region := "NorthAmerica"

	if os.Getenv("PINGONE_REGION") == "NorthAmerica" {
		region = "Europe"
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEnvironmentDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccEnvironmentConfig_MinimalWithRegion(resourceName, name, environmentType, region, licenseID),
				ExpectError: regexp.MustCompile(fmt.Sprintf(`Incompatible environment region for the organization tenant.  Expecting regions \[%s\], region provided: %s`, model.FindRegionByName(os.Getenv("PINGONE_REGION")).Region, model.FindRegionByName(region).Region)),
			},
		},
	})
}

func TestAccEnvironment_DeleteProductionEnvironmentProtection(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGenEnvironment()
	//resourceFullName := fmt.Sprintf("pingone_environment.%s", resourceName)

	name := resourceName
	environmentType := "SANDBOX"
	licenseID := os.Getenv("PINGONE_LICENSE_ID")
	region := os.Getenv("PINGONE_REGION")

	// os.Setenv("PINGONE_FORCE_DELETE_PRODUCTION_TYPE", "false")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { t.Skipf("Test to be defined") },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEnvironmentDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccEnvironmentConfig_MinimalWithRegion(resourceName, name, environmentType, region, licenseID),
				ExpectError: regexp.MustCompile(`Not defined`),
			},
		},
	})
}

func TestAccEnvironment_DeleteProductionEnvironment(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGenEnvironment()
	// resourceFullName := fmt.Sprintf("pingone_environment.%s", resourceName)

	name := resourceName
	environmentType := "SANDBOX"
	licenseID := os.Getenv("PINGONE_LICENSE_ID")
	region := os.Getenv("PINGONE_REGION")

	// os.Setenv("PINGONE_FORCE_DELETE_PRODUCTION_TYPE", "true")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { t.Skipf("Test to be defined") },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEnvironmentDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccEnvironmentConfig_MinimalWithRegion(resourceName, name, environmentType, region, licenseID),
				ExpectError: regexp.MustCompile(`Not defined`),
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
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEnvironmentDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironmentConfig_Full(resourceName, name, description, environmentType, region, licenseID, solution, populationName, populationDescription, serviceOneType, serviceTwoType, serviceTwoURL, serviceTwoBookmarkNameOne, serviceTwoBookmarkURLOne, serviceTwoBookmarkNameTwo, serviceTwoBookmarkURLTwo),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestMatchResourceAttr(resourceFullName, "default_population_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "default_population.0.name", populationName),
					resource.TestCheckResourceAttr(resourceFullName, "default_population.0.description", populationDescription),
					resource.TestCheckResourceAttr(resourceFullName, "service.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "service.*", map[string]string{
						"type":        serviceOneType,
						"console_url": "",
						"bookmark.#":  "0",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "service.*", map[string]string{
						"type":            serviceTwoType,
						"console_url":     serviceTwoURL,
						"bookmark.#":      "2",
						"bookmark.0.name": serviceTwoBookmarkNameOne,
						"bookmark.0.url":  serviceTwoBookmarkURLOne,
						"bookmark.1.name": serviceTwoBookmarkNameTwo,
						"bookmark.1.url":  serviceTwoBookmarkURLTwo,
					}),
				),
			},
		},
	})
}

func TestAccEnvironment_EnvironmentTypeSwitching(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGenEnvironment()
	resourceFullName := fmt.Sprintf("pingone_environment.%s", resourceName)

	name := resourceName
	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEnvironmentDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironmentConfig_Minimal(resourceName, name, "SANDBOX", licenseID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "type", "SANDBOX"),
				),
			},
			{
				Config: testAccEnvironmentConfig_Minimal(resourceName, name, "PRODUCTION", licenseID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "type", "PRODUCTION"),
				),
			},
			{
				Config: testAccEnvironmentConfig_Minimal(resourceName, name, "SANDBOX", licenseID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "type", "SANDBOX"),
				),
			},
		},
	})
}

func TestAccEnvironment_ServiceAndPopulationSwitching(t *testing.T) {
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
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEnvironmentDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironmentConfig_Minimal(resourceName, name, "SANDBOX", licenseID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "default_population_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "default_population.0.name", "Default"),
					resource.TestCheckResourceAttr(resourceFullName, "default_population.0.description", ""),
					resource.TestCheckResourceAttr(resourceFullName, "service.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "service.*", map[string]string{
						"type":        "SSO",
						"console_url": "",
						"bookmark.#":  "0",
					}),
				),
			},
			{
				Config: testAccEnvironmentConfig_Full(resourceName, name, description, environmentType, region, licenseID, solution, populationName, populationDescription, serviceOneType, serviceTwoType, serviceTwoURL, serviceTwoBookmarkNameOne, serviceTwoBookmarkURLOne, serviceTwoBookmarkNameTwo, serviceTwoBookmarkURLTwo),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "default_population_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "default_population.0.name", populationName),
					resource.TestCheckResourceAttr(resourceFullName, "default_population.0.description", populationDescription),
					resource.TestCheckResourceAttr(resourceFullName, "service.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "service.*", map[string]string{
						"type":        serviceOneType,
						"console_url": "",
						"bookmark.#":  "0",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "service.*", map[string]string{
						"type":            serviceTwoType,
						"console_url":     serviceTwoURL,
						"bookmark.#":      "2",
						"bookmark.0.name": serviceTwoBookmarkNameOne,
						"bookmark.0.url":  serviceTwoBookmarkURLOne,
						"bookmark.1.name": serviceTwoBookmarkNameTwo,
						"bookmark.1.url":  serviceTwoBookmarkURLTwo,
					}),
				),
			},
		},
	})
}

func TestAccEnvironment_Services(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGenEnvironment()
	resourceFullName := fmt.Sprintf("pingone_environment.%s", resourceName)

	name := resourceName
	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	services1 := []string{`SSO`, `MFA`, `Risk`, `Verify`, `Credentials`, `APIIntelligence`, `Authorize`, `PingFederate`, `PingAccess`, `PingDirectory`, `PingAuthorize`, `PingCentral`}
	services2 := []string{`SSO`, `MFA`, `Risk`, `Verify`}
	services3 := []string{`SSO`, `MFA`, `Risk`, `Verify`, `PingFederate`, `PingAccess`, `PingDirectory`, `PingAuthorize`, `PingCentral`}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEnvironmentDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironmentConfig_DynamicServices(resourceName, name, licenseID, services1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "service.#", "12"), // check all the custom services provision, except the WORKFORCE services
				),
			},
			{
				Config: testAccEnvironmentConfig_DynamicServices(resourceName, name, licenseID, services2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "service.#", "4"), // check they can be modified downward
				),
			},
			{
				Config: testAccEnvironmentConfig_DynamicServices(resourceName, name, licenseID, services3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "service.#", "9"), // check they can be modified upward
				),
			},
		},
	})
}

func testAccEnvironmentConfig_Full(resourceName, name, description, environmentType, region, licenseID, solution, populationName, populationDescription, serviceOneType, serviceTwoType, serviceTwoURL, serviceTwoBookmarkNameOne, serviceTwoBookmarkURLOne, serviceTwoBookmarkNameTwo, serviceTwoBookmarkURLTwo string) string {
	return fmt.Sprintf(`
resource "pingone_environment" "%[1]s" {
  name        = "%[2]s"
  description = "%[3]s"
  type        = "%[4]s"
  region      = "%[5]s"
  license_id  = "%[6]s"
  solution    = "%[7]s"
  default_population {
    name        = "%[8]s"
    description = "%[9]s"
  }
  service {
    type = "%[10]s"
  }
  service {
    type        = "%[11]s"
    console_url = "%[12]s"
    bookmark {
      name = "%[13]s"
      url  = "%[14]s"
    }
    bookmark {
      name = "%[15]s"
      url  = "%[16]s"
    }
  }
}`, resourceName, name, description, environmentType, region, licenseID, solution, populationName, populationDescription, serviceOneType, serviceTwoType, serviceTwoURL, serviceTwoBookmarkNameOne, serviceTwoBookmarkURLOne, serviceTwoBookmarkNameTwo, serviceTwoBookmarkURLTwo)
}

func testAccEnvironmentConfig_DynamicServices(resourceName, name, licenseID string, services []string) string {

	composedServices := composeServices(services)

	return fmt.Sprintf(`
resource "pingone_environment" "%[1]s" {
  name       = "%[2]s"
  license_id = "%[3]s"
  default_population {
  }
			%[4]s
}`, resourceName, name, licenseID, composedServices)
}

func testAccEnvironmentConfig_Minimal(resourceName, name, environmentType, licenseID string) string {
	return fmt.Sprintf(`
resource "pingone_environment" "%[1]s" {
  name       = "%[2]s"
  type       = "%[3]s"
  license_id = "%[4]s"
  default_population {
  }
  service {
  }
}`, resourceName, name, environmentType, licenseID)
}

func testAccEnvironmentConfig_MinimalWithRegion(resourceName, name, environmentType, region, licenseID string) string {
	return fmt.Sprintf(`
resource "pingone_environment" "%[1]s" {
  name       = "%[2]s"
  type       = "%[3]s"
  region     = "%[4]s"
  license_id = "%[5]s"
  default_population {
  }
  service {
  }
}`, resourceName, name, environmentType, region, licenseID)
}

func composeServices(services []string) string {

	var composedServices = fmt.Sprintf("service {\n  type = \"%s\"\n}", strings.Join(services, "\"\n}\nservice {\n  type = \""))

	return composedServices

}
