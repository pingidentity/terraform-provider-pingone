package base_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/base"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccEnvironment_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_environment.%s", resourceName)

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	var environmentID string

	var p1Client *client.Client
	var ctx = context.Background()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)

			p1Client = acctest.PreCheckTestClient(ctx, t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Environment_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccEnvironmentConfig_Minimal(resourceName, name, licenseID),
				Check:  base.Environment_GetIDs(resourceFullName, &environmentID),
			},
			{
				PreConfig: func() {
					base.Environment_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccEnvironment_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGenEnvironment()
	resourceFullName := fmt.Sprintf("pingone_environment.%s", resourceName)

	name := resourceName
	region := os.Getenv("PINGONE_REGION")
	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	fullStepVariant1 := resource.TestStep{

		Config: testAccEnvironmentConfig_Full(resourceName, fmt.Sprintf("%s-1", name), region, licenseID),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "name", fmt.Sprintf("%s-1", name)),
			resource.TestCheckResourceAttr(resourceFullName, "description", "Test description"),
			resource.TestCheckResourceAttr(resourceFullName, "type", "SANDBOX"),
			resource.TestCheckResourceAttr(resourceFullName, "region", region),
			resource.TestCheckResourceAttr(resourceFullName, "license_id", licenseID),
			resource.TestMatchResourceAttr(resourceFullName, "organization_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "solution", "CUSTOMER"),
			resource.TestCheckNoResourceAttr(resourceFullName, "default_population_id"),
			resource.TestCheckNoResourceAttr(resourceFullName, "default_population.0.name"),
			resource.TestCheckNoResourceAttr(resourceFullName, "default_population.0.description"),
			resource.TestCheckResourceAttr(resourceFullName, "service.#", "2"),
			resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "service.*", map[string]string{
				"type":        "SSO",
				"console_url": "",
				"bookmark.#":  "0",
			}),
			resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "service.*", map[string]string{
				"type":            "PingFederate",
				"console_url":     "https://my-console-url",
				"bookmark.#":      "2",
				"bookmark.0.name": "Bookmark 1",
				"bookmark.0.url":  "https://my-bookmark-1",
				"bookmark.1.name": "Bookmark 2",
				"bookmark.1.url":  "https://my-bookmark-2",
			}),
		),
	}

	fullStepVariant2 := resource.TestStep{

		Config: testAccEnvironmentConfig_Full(resourceName, fmt.Sprintf("%s-2", name), region, licenseID),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "name", fmt.Sprintf("%s-2", name)),
			resource.TestCheckResourceAttr(resourceFullName, "description", "Test description"),
			resource.TestCheckResourceAttr(resourceFullName, "type", "SANDBOX"),
			resource.TestCheckResourceAttr(resourceFullName, "region", region),
			resource.TestCheckResourceAttr(resourceFullName, "license_id", licenseID),
			resource.TestMatchResourceAttr(resourceFullName, "organization_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "solution", "CUSTOMER"),
			resource.TestCheckNoResourceAttr(resourceFullName, "default_population_id"),
			resource.TestCheckNoResourceAttr(resourceFullName, "default_population.0.name"),
			resource.TestCheckNoResourceAttr(resourceFullName, "default_population.0.description"),
			resource.TestCheckResourceAttr(resourceFullName, "service.#", "2"),
			resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "service.*", map[string]string{
				"type":        "SSO",
				"console_url": "",
				"bookmark.#":  "0",
			}),
			resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "service.*", map[string]string{
				"type":            "PingFederate",
				"console_url":     "https://my-console-url",
				"bookmark.#":      "2",
				"bookmark.0.name": "Bookmark 1",
				"bookmark.0.url":  "https://my-bookmark-1",
				"bookmark.1.name": "Bookmark 2",
				"bookmark.1.url":  "https://my-bookmark-2",
			}),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Environment_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			fullStepVariant1,
			fullStepVariant2,
			fullStepVariant1,
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return rs.Primary.ID, nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
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

	minimalStep := resource.TestStep{
		Config: testAccEnvironmentConfig_Minimal(resourceName, name, licenseID),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestCheckNoResourceAttr(resourceFullName, "description"),
			resource.TestCheckResourceAttr(resourceFullName, "type", environmentType),
			resource.TestCheckResourceAttr(resourceFullName, "region", region),
			resource.TestCheckNoResourceAttr(resourceFullName, "solution"),
			resource.TestCheckResourceAttr(resourceFullName, "license_id", licenseID),
			resource.TestMatchResourceAttr(resourceFullName, "organization_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckNoResourceAttr(resourceFullName, "default_population_id"),
			resource.TestCheckNoResourceAttr(resourceFullName, "default_population.0.name"),
			resource.TestCheckNoResourceAttr(resourceFullName, "default_population.0.description"),
			resource.TestCheckResourceAttr(resourceFullName, "service.#", "1"),
			resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "service.*", map[string]string{
				"type":        "SSO",
				"console_url": "",
				"bookmark.#":  "0",
			}),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Environment_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			minimalStep,
		},
	})
}

func TestAccEnvironment_NonCompatibleRegion(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")
	region := "NorthAmerica"

	if os.Getenv("PINGONE_REGION") == "NorthAmerica" {
		region = "Europe"
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Environment_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccEnvironmentConfig_MinimalWithRegion(resourceName, name, region, licenseID),
				ExpectError: regexp.MustCompile(fmt.Sprintf(`Incompatible environment region for the organization tenant.  Allowed regions: \[%s\].`, model.FindRegionByName(os.Getenv("PINGONE_REGION")).Region)),
			},
		},
	})
}

func TestAccEnvironment_DeleteProductionEnvironmentProtection(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGenEnvironment()
	//resourceFullName := fmt.Sprintf("pingone_environment.%s", resourceName)

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")
	region := os.Getenv("PINGONE_REGION")

	// os.Setenv("PINGONE_FORCE_DELETE_PRODUCTION_TYPE", "false")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { t.Skipf("Test to be defined") },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Environment_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccEnvironmentConfig_MinimalWithRegion(resourceName, name, region, licenseID),
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

	licenseID := os.Getenv("PINGONE_LICENSE_ID")
	region := os.Getenv("PINGONE_REGION")

	// os.Setenv("PINGONE_FORCE_DELETE_PRODUCTION_TYPE", "true")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { t.Skipf("Test to be defined") },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Environment_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccEnvironmentConfig_MinimalWithRegion(resourceName, name, region, licenseID),
				ExpectError: regexp.MustCompile(`Not defined`),
			},
		},
	})
}

///////////////////
// Deprecated start

func TestAccEnvironment_MinimalWithPopulation(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGenEnvironment()
	resourceFullName := fmt.Sprintf("pingone_environment.%s", resourceName)

	name := resourceName
	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Environment_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironmentConfig_MinimalWithPopulation(resourceName, name, licenseID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "default_population_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "default_population.0.name", name),
					resource.TestCheckResourceAttr(resourceFullName, "default_population.0.description", fmt.Sprintf("%s description", name)),
				),
			},
		},
	})
}

// Deprecated end
///////////////////

func TestAccEnvironment_EnvironmentTypeSwitching(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGenEnvironment()
	resourceFullName := fmt.Sprintf("pingone_environment.%s", resourceName)

	name := resourceName
	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Environment_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironmentConfig_MinimalWithType(resourceName, name, "SANDBOX", licenseID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "type", "SANDBOX"),
				),
			},
			{
				Config: testAccEnvironmentConfig_MinimalWithType(resourceName, name, "PRODUCTION", licenseID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "type", "PRODUCTION"),
				),
			},
			{
				Config: testAccEnvironmentConfig_MinimalWithType(resourceName, name, "SANDBOX", licenseID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "type", "SANDBOX"),
				),
			},
		},
	})
}

func TestAccEnvironment_ServiceSwitching(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGenEnvironment()
	resourceFullName := fmt.Sprintf("pingone_environment.%s", resourceName)

	name := resourceName
	region := os.Getenv("PINGONE_REGION")
	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Environment_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironmentConfig_Minimal(resourceName, name, licenseID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "service.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "service.*", map[string]string{
						"type":        "SSO",
						"console_url": "",
						"bookmark.#":  "0",
					}),
				),
			},
			{
				Config: testAccEnvironmentConfig_Full(resourceName, name, region, licenseID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "service.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "service.*", map[string]string{
						"type":        "SSO",
						"console_url": "",
						"bookmark.#":  "0",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "service.*", map[string]string{
						"type":            "PingFederate",
						"console_url":     "https://my-console-url",
						"bookmark.#":      "2",
						"bookmark.0.name": "Bookmark 1",
						"bookmark.0.url":  "https://my-bookmark-1",
						"bookmark.1.name": "Bookmark 2",
						"bookmark.1.url":  "https://my-bookmark-2",
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
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Environment_CheckDestroy,
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
			{
				Config: testAccEnvironmentConfig_Minimal(resourceName, name, licenseID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "service.#", "1"), // check they can be defaulted
				),
			},
			{
				Config: testAccEnvironmentConfig_DynamicServices(resourceName, name, licenseID, services3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "service.#", "9"), // check they can be un-defaulted
				),
			},
		},
	})
}

func TestAccEnvironment_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_environment.%s", resourceName)

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Environment_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccEnvironmentConfig_Minimal(resourceName, name, licenseID),
			},
			// Errors
			{
				ResourceName:  resourceFullName,
				ImportStateId: "badformat",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`Unexpected Import Identifier`),
			},
		},
	})
}

func testAccEnvironmentConfig_Full(resourceName, name, region, licenseID string) string {
	return fmt.Sprintf(`
resource "pingone_environment" "%[1]s" {
  name        = "%[2]s"
  description = "Test description"
  type        = "SANDBOX"
  region      = "%[3]s"
  license_id  = "%[4]s"
  solution    = "CUSTOMER"

  service {
    type = "SSO"
  }

  service {
    type        = "PingFederate"
    console_url = "https://my-console-url"
    bookmark {
      name = "Bookmark 1"
      url  = "https://my-bookmark-1"
    }
    bookmark {
      name = "Bookmark 2"
      url  = "https://my-bookmark-2"
    }
  }
}`, resourceName, name, region, licenseID)
}

func testAccEnvironmentConfig_DynamicServices(resourceName, name, licenseID string, services []string) string {

	composedServices := composeServices(services)

	return fmt.Sprintf(`
resource "pingone_environment" "%[1]s" {
  name       = "%[2]s"
  license_id = "%[3]s"

  %[4]s
}`, resourceName, name, licenseID, composedServices)
}

func testAccEnvironmentConfig_MinimalWithPopulation(resourceName, name, licenseID string) string {
	return fmt.Sprintf(`
resource "pingone_environment" "%[1]s" {
  name       = "%[2]s"
  license_id = "%[3]s"

  default_population {
    name        = "%[2]s"
    description = "%[2]s description"
  }
}`, resourceName, name, licenseID)
}

func testAccEnvironmentConfig_Minimal(resourceName, name, licenseID string) string {
	return fmt.Sprintf(`
resource "pingone_environment" "%[1]s" {
  name       = "%[2]s"
  license_id = "%[3]s"
}`, resourceName, name, licenseID)
}

func testAccEnvironmentConfig_MinimalWithType(resourceName, name, environmentType, licenseID string) string {
	return fmt.Sprintf(`
resource "pingone_environment" "%[1]s" {
  name       = "%[2]s"
  type       = "%[3]s"
  license_id = "%[4]s"
}`, resourceName, name, environmentType, licenseID)
}

func testAccEnvironmentConfig_MinimalWithRegion(resourceName, name, region, licenseID string) string {
	return fmt.Sprintf(`
resource "pingone_environment" "%[1]s" {
  name       = "%[2]s"
  region     = "%[3]s"
  license_id = "%[4]s"
}`, resourceName, name, region, licenseID)
}

func composeServices(services []string) string {

	var composedServices = fmt.Sprintf("service {\n  type = \"%s\"\n}", strings.Join(services, "\"\n}\nservice {\n  type = \""))

	return composedServices

}
