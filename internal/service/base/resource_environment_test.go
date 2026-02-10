// Copyright © 2026 Ping Identity Corporation

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
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	acctestlegacysdk "github.com/pingidentity/terraform-provider-pingone/internal/acctest/legacysdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/base"
	baselegacysdk "github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/base/legacysdk"
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
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoBeta(t)
			p1Client = acctestlegacysdk.PreCheckTestClient(ctx, t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             baselegacysdk.Environment_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccEnvironmentConfig_Minimal(resourceName, name, licenseID),
				Check:  base.Environment_GetIDs(resourceFullName, &environmentID),
			},
			{
				PreConfig: func() {
					baselegacysdk.Environment_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID)
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
	region := os.Getenv("PINGONE_REGION_CODE")
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
			resource.TestCheckResourceAttr(resourceFullName, "services.#", "2"),
			resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "services.*", map[string]string{
				"type": "SSO",
			}),
			resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "services.*", map[string]string{
				"type":             "PingFederate",
				"console_url":      "https://my-console-url",
				"bookmarks.#":      "2",
				"bookmarks.0.name": "Bookmark 1",
				"bookmarks.0.url":  "https://my-bookmark-1",
				"bookmarks.1.name": "Bookmark 2",
				"bookmarks.1.url":  "https://my-bookmark-2",
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
			resource.TestCheckResourceAttr(resourceFullName, "services.#", "2"),
			resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "services.*", map[string]string{
				"type": "SSO",
			}),
			resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "services.*", map[string]string{
				"type":             "PingFederate",
				"console_url":      "https://my-console-url",
				"bookmarks.#":      "2",
				"bookmarks.0.name": "Bookmark 1",
				"bookmarks.0.url":  "https://my-bookmark-1",
				"bookmarks.1.name": "Bookmark 2",
				"bookmarks.1.url":  "https://my-bookmark-2",
			}),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             baselegacysdk.Environment_CheckDestroy,
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
							return "", fmt.Errorf("resource not found: %s", resourceFullName)
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
	region := os.Getenv("PINGONE_REGION_CODE")

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
			resource.TestCheckResourceAttr(resourceFullName, "services.#", "1"),
			resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "services.*", map[string]string{
				"type": "SSO",
			}),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             baselegacysdk.Environment_CheckDestroy,
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
	region := "NA"

	if os.Getenv("PINGONE_REGION_CODE") == "NA" {
		region = "EU"
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             baselegacysdk.Environment_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccEnvironmentConfig_MinimalWithRegion(resourceName, name, region, licenseID),
				ExpectError: regexp.MustCompile(fmt.Sprintf(`Allowed regions: \[%[1]s(?: [A-Z]{2})?|(?: [A-Z]{2})?%[1]s\]\.`, model.FindRegionByAPICode(management.EnumRegionCode(os.Getenv("PINGONE_REGION_CODE"))).APICode)),
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
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             nil,
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
				Config:      testAccEnvironmentConfig_MinimalWithType(resourceName, name, "SANDBOX", licenseID),
				ExpectError: regexp.MustCompile(`Data protection notice - The environment type cannot be changed from PRODUCTION to SANDBOX`),
			},
		},
	})
}

func TestAccEnvironment_ServiceSwitching(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGenEnvironment()
	resourceFullName := fmt.Sprintf("pingone_environment.%s", resourceName)

	name := resourceName
	region := os.Getenv("PINGONE_REGION_CODE")
	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             baselegacysdk.Environment_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironmentConfig_Minimal(resourceName, name, licenseID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "services.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "services.*", map[string]string{
						"type": "SSO",
					}),
				),
			},
			{
				Config: testAccEnvironmentConfig_Full(resourceName, name, region, licenseID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "services.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "services.*", map[string]string{
						"type": "SSO",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "services.*", map[string]string{
						"type":             "PingFederate",
						"console_url":      "https://my-console-url",
						"bookmarks.#":      "2",
						"bookmarks.0.name": "Bookmark 1",
						"bookmarks.0.url":  "https://my-bookmark-1",
						"bookmarks.1.name": "Bookmark 2",
						"bookmarks.1.url":  "https://my-bookmark-2",
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
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             baselegacysdk.Environment_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironmentConfig_DynamicServices(resourceName, name, licenseID, services1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "services.#", "12"), // check all the custom services provision, except the WORKFORCE services
				),
			},
			{
				Config: testAccEnvironmentConfig_DynamicServices(resourceName, name, licenseID, services2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "services.#", "4"), // check they can be modified downward
				),
			},
			{
				Config: testAccEnvironmentConfig_DynamicServices(resourceName, name, licenseID, services3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "services.#", "9"), // check they can be modified upward
				),
			},
			{
				Config: testAccEnvironmentConfig_Minimal(resourceName, name, licenseID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "services.#", "1"), // check they can be defaulted
				),
			},
			{
				Config: testAccEnvironmentConfig_DynamicServices(resourceName, name, licenseID, services3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "services.#", "9"), // check they can be un-defaulted
				),
			},
		},
	})
}

func TestAccEnvironment_ServicesTags(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGenEnvironment()
	resourceFullName := fmt.Sprintf("pingone_environment.%s", resourceName)

	name := resourceName
	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             baselegacysdk.Environment_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironmentConfig_DVTags(resourceName, name, licenseID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "services.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "services.*", map[string]string{
						"type":   "DaVinci",
						"tags.#": "1",
						"tags.0": "DAVINCI_MINIMAL",
					}),
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
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             baselegacysdk.Environment_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccEnvironmentConfig_Workforce(resourceName, name, licenseID),
				ExpectError: regexp.MustCompile(`Cannot create workforce environments`),
			},
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

func TestAccEnvironment_DeploymentID(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckRegionSupportsWorkforce(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             nil,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: acctest.WorkforceV2SandboxEnvironment(),
				Check: resource.ComposeTestCheckFunc(
					testCheckEnvironmentServiceHasDeploymentID("data.pingone_environment.workforce_test", "PingID-v2"),
				),
			},
		},
	})
}

func TestAccEnvironment_ServiceInvalidAddPingID(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGenEnvironment()

	name := resourceName
	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             baselegacysdk.Environment_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironmentConfig_Minimal(resourceName, name, licenseID),
			},
			{
				Config:      testAccEnvironmentConfig_AddPingID(resourceName, name, licenseID),
				ExpectError: regexp.MustCompile(`The .* service cannot be added via Terraform configuration`),
			},
		},
	})
}

func TestAccEnvironment_ServiceInvalidRemovePingID(t *testing.T) {
	resourceName := acctest.ResourceNameGenEnvironment()
	resourceFullName := fmt.Sprintf("pingone_environment.%s", resourceName)

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	var environmentID string

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckRegionSupportsWorkforce(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             nil,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Step 1: Find the environment ID
			{
				Config: acctest.WorkforceV2SandboxEnvironment(),
				Check:  base.Environment_GetIDs("data.pingone_environment.workforce_test", &environmentID),
			},
			// Step 2: Import with PingID present in config
			{
				Config:       testAccEnvironmentConfig_WorkforceWithPingID(resourceName, acctest.WorkforceV2SandboxEnvironmentName, licenseID),
				ResourceName: resourceFullName,
				ImportStateIdFunc: func(_ *terraform.State) (string, error) {
					return environmentID, nil
				},
				ImportState:        true,
				ImportStateVerify:  false,
				ImportStatePersist: true,
			},
			// Step 3: Remove PingID in config -> Expect Error
			{
				Config:      testAccEnvironmentConfig_Workforce(resourceName, acctest.WorkforceV2SandboxEnvironmentName, licenseID),
				ExpectError: regexp.MustCompile(`Cannot remove .* service`),
			},
			// Step 4: Prevent destroy of Workforce environment
			{
				Config: testAccEnvironmentConfig_RemoveWorkforcePreventDestroy(resourceName, acctest.WorkforceV2SandboxEnvironmentName, licenseID),
			},
		},
	})
}

func testCheckEnvironmentServiceHasDeploymentID(resourceName, serviceType string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		for key, value := range rs.Primary.Attributes {
			// The key format is services.HASH.type
			if strings.HasPrefix(key, "services.") && strings.HasSuffix(key, ".type") && value == serviceType {
				baseKey := strings.TrimSuffix(key, ".type")
				deploymentIdKey := fmt.Sprintf("%s.deployment.id", baseKey)

				if v, ok := rs.Primary.Attributes[deploymentIdKey]; ok && v != "" {
					return nil
				}

				return fmt.Errorf("Deployment ID not found for service type %s", serviceType)
			}
		}

		return fmt.Errorf("Service type %s not found in environment", serviceType)
	}
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

  services = [
    {
      type = "SSO"
    },
    {
      type        = "PingFederate"
      console_url = "https://my-console-url"
      bookmarks = [
        {
          name = "Bookmark 1"
          url  = "https://my-bookmark-1"
        },
        {
          name = "Bookmark 2"
          url  = "https://my-bookmark-2"
        }
      ]
    }
  ]
}`, resourceName, name, region, licenseID)
}

func testAccEnvironmentConfig_DynamicServices(resourceName, name, licenseID string, services []string) string {
	return fmt.Sprintf(`


variable "services_%[1]s" {
  type    = list(string)
  default = ["%[4]s"]
}

resource "pingone_environment" "%[1]s" {
  name       = "%[2]s"
  license_id = "%[3]s"

  services = [
    for serviceType in var.services_%[1]s : {
      type = serviceType
    }
  ]
}`, resourceName, name, licenseID, strings.Join(services, "\",\""))
}

func testAccEnvironmentConfig_DVTags(resourceName, name, licenseID string) string {
	return fmt.Sprintf(`
resource "pingone_environment" "%[1]s" {
  name       = "%[2]s"
  license_id = "%[3]s"

  services = [
    {
      type = "SSO"
    },
    {
      type = "DaVinci"
      tags = ["DAVINCI_MINIMAL"]
    }
  ]
}`, resourceName, name, licenseID)
}

func testAccEnvironmentConfig_Minimal(resourceName, name, licenseID string) string {
	return fmt.Sprintf(`
resource "pingone_environment" "%[1]s" {
  name       = "%[2]s"
  license_id = "%[3]s"

  services = [
    {
      type = "SSO"
    }
  ]
}`, resourceName, name, licenseID)
}

func testAccEnvironmentConfig_AddPingID(resourceName, name, licenseID string) string {
	return fmt.Sprintf(`
resource "pingone_environment" "%[1]s" {
  name       = "%[2]s"
  license_id = "%[3]s"

  services = [
    {
      type = "SSO"
    },
    {
      type = "PingID-v2"
    }
  ]
}`, resourceName, name, licenseID)
}

func testAccEnvironmentConfig_WorkforceWithPingID(resourceName, name, licenseID string) string {
	return fmt.Sprintf(`
resource "pingone_environment" "%[1]s" {
  name       = "%[2]s"
  license_id = "%[3]s"
  solution   = "WORKFORCE"

  services = [
    {
      type = "SSO"
    },
    {
      type = "PingID-v2"
    }
  ]
}`, resourceName, name, licenseID)
}

func testAccEnvironmentConfig_RemoveWorkforcePreventDestroy(resourceName, name, licenseID string) string {
	return fmt.Sprintf(`
removed {
	from = pingone_environment.%[1]s

	lifecycle {
		destroy = false
	}
}`, resourceName, name, licenseID)
}

func testAccEnvironmentConfig_Workforce(resourceName, name, licenseID string) string {
	return fmt.Sprintf(`
resource "pingone_environment" "%[1]s" {
  name       = "%[2]s"
  license_id = "%[3]s"
  solution   = "WORKFORCE"

  services = [
    {
      type = "SSO"
    }
  ]

  lifecycle {
    prevent_destroy = true
  }
}`, resourceName, name, licenseID)
}

func testAccEnvironmentConfig_MinimalWithType(resourceName, name, environmentType, licenseID string) string {
	return fmt.Sprintf(`
resource "pingone_environment" "%[1]s" {
  name       = "%[2]s"
  type       = "%[3]s"
  license_id = "%[4]s"

  services = [
    {
      type = "SSO"
    }
  ]
}`, resourceName, name, environmentType, licenseID)
}

func testAccEnvironmentConfig_MinimalWithRegion(resourceName, name, region, licenseID string) string {
	return fmt.Sprintf(`
resource "pingone_environment" "%[1]s" {
  name       = "%[2]s"
  region     = "%[3]s"
  license_id = "%[4]s"

  services = [
    {
      type = "SSO"
    }
  ]
}`, resourceName, name, region, licenseID)
}
