package base_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func testAccCheckOrganizationDestroy(s *terraform.State) error {
	return nil
}

func TestAccOrganizationDataSource_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullName := fmt.Sprintf("data.pingone_organization.%s", resourceName)

	organizationID := os.Getenv("PINGONE_ORGANIZATION_ID")
	organizationName := os.Getenv("PINGONE_ORGANIZATION_NAME")

	region := os.Getenv("PINGONE_REGION")

	domainTld := "not-set"
	switch region {
	case "Europe":
		domainTld = "eu"
	case "NorthAmerica":
		domainTld = "com"
	case "Canada":
		domainTld = "ca"
	case "AsiaPacific":
		domainTld = "ap"
	}

	testCheck := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexp),
		resource.TestMatchResourceAttr(dataSourceFullName, "organization_id", verify.P1ResourceIDRegexp),
		resource.TestCheckResourceAttr(dataSourceFullName, "name", organizationName),
		resource.TestMatchResourceAttr(dataSourceFullName, "description", regexp.MustCompile(`^[a-zA-Z0-9 -_\\.]*$`)),
		resource.TestCheckResourceAttr(dataSourceFullName, "type", "INTERNAL"),
		resource.TestCheckResourceAttr(dataSourceFullName, "billing_connection_ids.#", "1"),
		resource.TestMatchResourceAttr(dataSourceFullName, "billing_connection_ids.0", regexp.MustCompile(`^[a-zA-Z0-9]*$`)),
		resource.TestCheckResourceAttr(dataSourceFullName, "base_url_api", fmt.Sprintf("api.pingone.%s", domainTld)),
		resource.TestCheckResourceAttr(dataSourceFullName, "base_url_auth", fmt.Sprintf("auth.pingone.%s", domainTld)),
		resource.TestCheckResourceAttr(dataSourceFullName, "base_url_orchestrate", fmt.Sprintf("orchestrate-api.pingone.%s", domainTld)),
		resource.TestCheckResourceAttr(dataSourceFullName, "base_url_agreement_management", fmt.Sprintf("agreement-mgmt.pingone.%s", domainTld)),
		resource.TestCheckResourceAttr(dataSourceFullName, "base_url_console", fmt.Sprintf("console.pingone.%s", domainTld)),
		resource.TestCheckResourceAttr(dataSourceFullName, "base_url_apps", fmt.Sprintf("apps.pingone.%s", domainTld)),
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckOrganisation(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckOrganizationDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationDataSourceConfig_ByIDFull(resourceName, organizationID),
				Check:  testCheck,
			},
			{
				Config: testAccOrganizationDataSourceConfig_ByNameFull(resourceName, organizationName),
				Check:  testCheck,
			},
		},
	})
}

func TestAccOrganizationDataSource_NotFound(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccOrganizationDataSourceConfig_NotFoundByName(resourceName),
				ExpectError: regexp.MustCompile("Cannot find organization from name"),
			},
			{
				Config:      testAccOrganizationDataSourceConfig_NotFoundByID(resourceName),
				ExpectError: regexp.MustCompile("Error when calling `ReadOneOrganization`: The request could not be completed. You do not have access to this resource."),
			},
		},
	})
}

func testAccOrganizationDataSourceConfig_ByIDFull(resourceName, organizationID string) string {
	return fmt.Sprintf(`
data "pingone_organization" "%[1]s" {
  organization_id = "%[2]s"
}`, resourceName, organizationID)
}

func testAccOrganizationDataSourceConfig_ByNameFull(resourceName, organizationName string) string {
	return fmt.Sprintf(`
data "pingone_organization" "%[1]s" {
  name = "%[2]s"
}`, resourceName, organizationName)
}

func testAccOrganizationDataSourceConfig_NotFoundByID(resourceName string) string {
	return fmt.Sprintf(`
data "pingone_organization" "%[1]s" {
  organization_id = "9c052a8a-14be-44e4-8f07-2662569994ce" // dummy ID that conforms to UUID v4
}`, resourceName)
}

func testAccOrganizationDataSourceConfig_NotFoundByName(resourceName string) string {
	return fmt.Sprintf(`
data "pingone_organization" "%[1]s" {
  name = "doesnotexist"
}`, resourceName)
}
