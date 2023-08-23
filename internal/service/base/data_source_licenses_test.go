package base_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccLicensesDataSource_BySCIMFilter(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullName := fmt.Sprintf("data.pingone_licenses.%s", resourceName)

	organizationID := os.Getenv("PINGONE_ORGANIZATION_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironmentAndOrganisation(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckLicenseDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLicensesDataSourceConfig_BySCIMFilter(resourceName, organizationID, fmt.Sprintf("(status eq \\\"active\\\") and (beginsAt lt \\\"%s\\\")", time.Now().Format(time.RFC3339))),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "organization_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "ids.#", "2"),
					resource.TestMatchResourceAttr(dataSourceFullName, "ids.0", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "ids.1", verify.P1ResourceIDRegexpFullString),
				),
			},
		},
	})
}

func TestAccLicensesDataSource_ByDataFilter(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullName := fmt.Sprintf("data.pingone_licenses.%s", resourceName)

	organizationID := os.Getenv("PINGONE_ORGANIZATION_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironmentAndOrganisation(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckLicenseDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLicensesDataSourceConfig_ByDataFilter1(resourceName, organizationID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "organization_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "ids.#", "2"),
					resource.TestMatchResourceAttr(dataSourceFullName, "ids.0", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "ids.1", verify.P1ResourceIDRegexpFullString),
				),
			},
			{
				Config: testAccLicensesDataSourceConfig_ByDataFilter2(resourceName, organizationID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "organization_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "ids.#", "1"),
					resource.TestMatchResourceAttr(dataSourceFullName, "ids.0", verify.P1ResourceIDRegexpFullString),
				),
			},
			{
				Config: testAccLicensesDataSourceConfig_ByDataFilter3(resourceName, organizationID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "organization_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "ids.#", "0"),
				),
			},
		},
	})
}

func TestAccLicensesDataSource_NotFound(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullName := fmt.Sprintf("data.pingone_licenses.%s", resourceName)

	organizationID := os.Getenv("PINGONE_ORGANIZATION_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironmentAndOrganisation(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckLicenseDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLicensesDataSourceConfig_BySCIMFilter(resourceName, organizationID, fmt.Sprintf("(status eq \\\"active\\\") and (beginsAt lt \\\"%s\\\")", "2006-01-02T15:04:05Z07:00")),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "organization_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "ids.#", "0"),
				),
			},
		},
	})
}

func testAccLicensesDataSourceConfig_BySCIMFilter(resourceName, organizationID, filter string) string {
	return fmt.Sprintf(`
data "pingone_licenses" "%[1]s" {
  organization_id = "%[2]s"
  scim_filter     = "%[3]s"
}`, resourceName, organizationID, filter)
}

func testAccLicensesDataSourceConfig_ByDataFilter1(resourceName, organizationID string) string {
	return fmt.Sprintf(`
data "pingone_licenses" "%[1]s" {
  organization_id = "%[2]s"

  data_filter {
    name   = "package"
    values = ["INTERNAL", "ADMIN"]
  }

  data_filter {
    name   = "status"
    values = ["ACTIVE"]
  }
}`, resourceName, organizationID)
}

func testAccLicensesDataSourceConfig_ByDataFilter2(resourceName, organizationID string) string {
	return fmt.Sprintf(`
data "pingone_licenses" "%[1]s" {
  organization_id = "%[2]s"

  data_filter {
    name   = "name"
    values = ["INTERNAL"]
  }

  data_filter {
    name   = "status"
    values = ["ACTIVE"]
  }
}`, resourceName, organizationID)
}

func testAccLicensesDataSourceConfig_ByDataFilter3(resourceName, organizationID string) string {
	return fmt.Sprintf(`
data "pingone_licenses" "%[1]s" {
  organization_id = "%[2]s"

  data_filter {
    name   = "package"
    values = ["INTERNAL"]
  }

  data_filter {
    name   = "status"
    values = ["EXPIRED"]
  }
}`, resourceName, organizationID)
}
