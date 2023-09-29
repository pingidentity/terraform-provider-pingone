package base_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/base"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccOrganizationDataSource_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullName := fmt.Sprintf("data.pingone_organization.%s", resourceName)

	organizationID := os.Getenv("PINGONE_ORGANIZATION_ID")
	organizationName := os.Getenv("PINGONE_ORGANIZATION_NAME")

	testCheck := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(dataSourceFullName, "organization_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(dataSourceFullName, "name", organizationName),
		resource.TestMatchResourceAttr(dataSourceFullName, "description", regexp.MustCompile(`^[a-zA-Z0-9 -_\\.]*$`)),
		resource.TestCheckResourceAttr(dataSourceFullName, "type", "INTERNAL"),
		resource.TestCheckResourceAttr(dataSourceFullName, "billing_connection_ids.#", "1"),
		resource.TestMatchResourceAttr(dataSourceFullName, "billing_connection_ids.0", regexp.MustCompile(`^[a-zA-Z0-9]*$`)),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckOrganisationID(t)
			acctest.PreCheckOrganisationName(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.TestAccCheckOrganizationDestroy,
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
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
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
