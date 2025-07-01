// Copyright © 2025 Ping Identity Corporation

package base_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/base"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccEnvironmentDataSource_ByNameFull(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGenEnvironment()
	resourceFullName := fmt.Sprintf("pingone_environment.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	name := resourceName
	description := "Test description"
	environmentType := "SANDBOX"
	region := os.Getenv("PINGONE_REGION_CODE")
	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	solution := "CUSTOMER"

	serviceOneType := "SSO"
	serviceTwoType := "PingFederate"
	serviceTwoURL := "https://my-console-url"
	serviceTwoBookmarkNameOne := "Bookmark 1"
	serviceTwoBookmarkURLOne := "https://my-bookmark-1"
	serviceTwoBookmarkNameTwo := "Bookmark 2"
	serviceTwoBookmarkURLTwo := "https://my-bookmark-2"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Environment_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironmentDataSourceConfig_ByNameFull(resourceName, name, description, environmentType, region, licenseID, solution, serviceOneType, serviceTwoType, serviceTwoURL, serviceTwoBookmarkNameOne, serviceTwoBookmarkURLOne, serviceTwoBookmarkNameTwo, serviceTwoBookmarkURLTwo),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "name", resourceFullName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "description", resourceFullName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "type", resourceFullName, "type"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "region", resourceFullName, "region"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "solution", resourceFullName, "solution"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "license_id", resourceFullName, "license_id"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "organization_id", resourceFullName, "organization_id"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "service.%", resourceFullName, "service.%"),
				),
			},
		},
	})
}

func TestAccEnvironmentDataSource_ByNameMinimal(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGenEnvironment()
	resourceFullName := fmt.Sprintf("pingone_environment.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	name := resourceName
	environmentType := "SANDBOX"
	licenseID := os.Getenv("PINGONE_LICENSE_ID")
	region := os.Getenv("PINGONE_REGION_CODE")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Environment_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironmentDataSourceConfig_ByNameMinimal(resourceName, name, environmentType, region, licenseID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "name", resourceFullName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "description", resourceFullName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "type", resourceFullName, "type"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "region", resourceFullName, "region"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "solution", resourceFullName, "solution"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "license_id", resourceFullName, "license_id"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "organization_id", resourceFullName, "organization_id"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "service.%", resourceFullName, "service.%"),
				),
			},
		},
	})
}

func TestAccEnvironmentDataSource_ByIDFull(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGenEnvironment()
	resourceFullName := fmt.Sprintf("pingone_environment.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	name := resourceName
	description := "Test description"
	environmentType := "SANDBOX"
	region := os.Getenv("PINGONE_REGION_CODE")
	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	solution := "CUSTOMER"

	serviceOneType := "SSO"
	serviceTwoType := "PingFederate"
	serviceTwoURL := "https://my-console-url"
	serviceTwoBookmarkNameOne := "Bookmark 1"
	serviceTwoBookmarkURLOne := "https://my-bookmark-1"
	serviceTwoBookmarkNameTwo := "Bookmark 2"
	serviceTwoBookmarkURLTwo := "https://my-bookmark-2"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Environment_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironmentDataSourceConfig_ByIDFull(resourceName, name, description, environmentType, region, licenseID, solution, serviceOneType, serviceTwoType, serviceTwoURL, serviceTwoBookmarkNameOne, serviceTwoBookmarkURLOne, serviceTwoBookmarkNameTwo, serviceTwoBookmarkURLTwo),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "name", resourceFullName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "description", resourceFullName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "type", resourceFullName, "type"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "region", resourceFullName, "region"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "solution", resourceFullName, "solution"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "license_id", resourceFullName, "license_id"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "organization_id", resourceFullName, "organization_id"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "service.%", resourceFullName, "service.%"),
				),
			},
		},
	})
}

func TestAccEnvironmentDataSource_ByIDMinimal(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGenEnvironment()
	resourceFullName := fmt.Sprintf("pingone_environment.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	name := resourceName
	environmentType := "SANDBOX"
	licenseID := os.Getenv("PINGONE_LICENSE_ID")
	region := os.Getenv("PINGONE_REGION_CODE")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Environment_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironmentDataSourceConfig_ByIDMinimal(resourceName, name, environmentType, region, licenseID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "name", resourceFullName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "description", resourceFullName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "type", resourceFullName, "type"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "region", resourceFullName, "region"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "solution", resourceFullName, "solution"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "license_id", resourceFullName, "license_id"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "organization_id", resourceFullName, "organization_id"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "service.%", resourceFullName, "service.%"),
				),
			},
		},
	})
}

func TestAccEnvironmentDataSource_Workforce(t *testing.T) {
	t.Parallel()

	dataSourceFullName := "data.pingone_environment.workforce_test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
			acctest.PreCheckRegionSupportsWorkforce(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             nil, // No destroy check as purely data sources
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironmentDataSourceConfig_WorkForceV1ByName(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "name", acctest.WorkforceV1SandboxEnvironmentName),
					resource.TestCheckResourceAttr(dataSourceFullName, "solution", "WORKFORCE"),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourceFullName, "services.*", map[string]string{
						"type": "PingID",
					}),
				),
			},
			{
				Config: testAccEnvironmentDataSourceConfig_WorkForceV2ByName(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "name", acctest.WorkforceV2SandboxEnvironmentName),
					resource.TestCheckResourceAttr(dataSourceFullName, "solution", "WORKFORCE"),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourceFullName, "services.*", map[string]string{
						"type": "PingID-v2",
					}),
				),
			},
		},
	})
}

func TestAccEnvironmentDataSource_NotFound(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGenEnvironment()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Environment_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccEnvironmentDataSourceConfig_NotFoundByName(resourceName),
				ExpectError: regexp.MustCompile("Cannot find environment from name"),
			},
			{
				Config:      testAccEnvironmentDataSourceConfig_NotFoundByID(resourceName),
				ExpectError: regexp.MustCompile("Error when calling `ReadOneEnvironment`: Unable to find environment with ID: '9c052a8a-14be-44e4-8f07-2662569994ce'"),
			},
		},
	})
}

func testAccEnvironmentDataSourceConfig_ByNameFull(resourceName, name, description, environmentType, region, licenseID, solution, serviceOneType, serviceTwoType, serviceTwoURL, serviceTwoBookmarkNameOne, serviceTwoBookmarkURLOne, serviceTwoBookmarkNameTwo, serviceTwoBookmarkURLTwo string) string {
	return fmt.Sprintf(`
resource "pingone_environment" "%[1]s" {
  name        = "%[2]s"
  description = "%[3]s"
  type        = "%[4]s"
  region      = "%[5]s"
  license_id  = "%[6]s"

  services = [
    {
      type = "%[8]s"
    },
    {
      type        = "%[9]s"
      console_url = "%[10]s"
      bookmarks = [
        {
          name = "%[11]s"
          url  = "%[12]s"
        },
        {
          name = "%[13]s"
          url  = "%[14]s"
        }
      ]
    }
  ]
}

data "pingone_environment" "%[1]s" {
  name = "%[2]s"

  depends_on = [
    pingone_environment.%[1]s
  ]
}`, resourceName, name, description, environmentType, region, licenseID, solution, serviceOneType, serviceTwoType, serviceTwoURL, serviceTwoBookmarkNameOne, serviceTwoBookmarkURLOne, serviceTwoBookmarkNameTwo, serviceTwoBookmarkURLTwo)
}

func testAccEnvironmentDataSourceConfig_ByNameMinimal(resourceName, name, environmentType, region, licenseID string) string {
	return fmt.Sprintf(`
resource "pingone_environment" "%[1]s" {
  name       = "%[2]s"
  type       = "%[3]s"
  region     = "%[4]s"
  license_id = "%[5]s"

  services = [
    {
      type = "SSO"
    }
  ]
}
data "pingone_environment" "%[1]s" {
  name = "%[2]s"

  depends_on = [
    pingone_environment.%[1]s
  ]
}
`, resourceName, name, environmentType, region, licenseID)
}

func testAccEnvironmentDataSourceConfig_ByIDFull(resourceName, name, description, environmentType, region, licenseID, solution, serviceOneType, serviceTwoType, serviceTwoURL, serviceTwoBookmarkNameOne, serviceTwoBookmarkURLOne, serviceTwoBookmarkNameTwo, serviceTwoBookmarkURLTwo string) string {
	return fmt.Sprintf(`
resource "pingone_environment" "%[1]s" {
  name        = "%[2]s"
  description = "%[3]s"
  type        = "%[4]s"
  region      = "%[5]s"
  license_id  = "%[6]s"

  services = [
    {
      type = "%[8]s"
    },
    {
      type        = "%[9]s"
      console_url = "%[10]s"
      bookmarks = [
        {
          name = "%[11]s"
          url  = "%[12]s"
        },
        {
          name = "%[13]s"
          url  = "%[14]s"
        }
      ]
    }
  ]
}
data "pingone_environment" "%[1]s" {
  environment_id = pingone_environment.%[1]s.id
}`, resourceName, name, description, environmentType, region, licenseID, solution, serviceOneType, serviceTwoType, serviceTwoURL, serviceTwoBookmarkNameOne, serviceTwoBookmarkURLOne, serviceTwoBookmarkNameTwo, serviceTwoBookmarkURLTwo)
}

func testAccEnvironmentDataSourceConfig_ByIDMinimal(resourceName, name, environmentType, region, licenseID string) string {
	return fmt.Sprintf(`
resource "pingone_environment" "%[1]s" {
  name       = "%[2]s"
  type       = "%[3]s"
  region     = "%[4]s"
  license_id = "%[5]s"

  services = [
    {
      type = "SSO"
    }
  ]
}
data "pingone_environment" "%[1]s" {
  environment_id = pingone_environment.%[1]s.id
}
`, resourceName, name, environmentType, region, licenseID)
}

func testAccEnvironmentDataSourceConfig_WorkForceV1ByName() string {
	return acctest.WorkforceV1SandboxEnvironment()
}

func testAccEnvironmentDataSourceConfig_WorkForceV2ByName() string {
	return acctest.WorkforceV2SandboxEnvironment()
}

func testAccEnvironmentDataSourceConfig_NotFoundByName(resourceName string) string {
	return fmt.Sprintf(`

data "pingone_environment" "%[1]s" {
  name = "doesnotexist"
}`, resourceName)
}

func testAccEnvironmentDataSourceConfig_NotFoundByID(resourceName string) string {
	return fmt.Sprintf(`

data "pingone_environment" "%[1]s" {
  environment_id = "9c052a8a-14be-44e4-8f07-2662569994ce" // dummy ID that conforms to UUID v4
}`, resourceName)
}
