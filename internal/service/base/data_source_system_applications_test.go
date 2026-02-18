// Copyright © 2026 Ping Identity Corporation

package base_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccSystemApplicationsDataSource_GetAll(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullName := fmt.Sprintf("data.pingone_system_applications.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSystemApplicationsDataSourceConfig_GetAll(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "ids.#", "2"),
					resource.TestMatchResourceAttr(dataSourceFullName, "ids.0", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "ids.1", verify.P1ResourceIDRegexpFullString),
					testAccSystemApplicationsDataSource_CheckIDs(resourceName),
				),
			},
		},
	})
}

func testAccSystemApplicationsDataSourceConfig_GetAll(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_system_application" "%[2]s_portal" {
  environment_id = data.pingone_environment.general_test.id
  type           = "PING_ONE_PORTAL"
}

data "pingone_system_application" "%[2]s_self_service" {
  environment_id = data.pingone_environment.general_test.id
  type           = "PING_ONE_SELF_SERVICE"
}

data "pingone_system_applications" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  depends_on = [
    data.pingone_system_application.%[2]s_portal,
    data.pingone_system_application.%[2]s_self_service,
  ]
}
`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

// Check that the correct two IDs are returned in the list data source
func testAccSystemApplicationsDataSource_CheckIDs(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		portalDataSourceName := fmt.Sprintf("data.pingone_system_application.%s_portal", resourceName)
		selfServiceDataSourceName := fmt.Sprintf("data.pingone_system_application.%s_self_service", resourceName)
		listDataSourceName := fmt.Sprintf("data.pingone_system_applications.%s", resourceName)

		portalDs, ok := s.RootModule().Resources[portalDataSourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", portalDataSourceName)
		}
		selfServiceDs, ok := s.RootModule().Resources[selfServiceDataSourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", selfServiceDataSourceName)
		}
		listDs, ok := s.RootModule().Resources[listDataSourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", listDataSourceName)
		}

		portalID := portalDs.Primary.ID
		selfServiceID := selfServiceDs.Primary.ID

		// Get the ids from the list data source to compare
		id1, ok := listDs.Primary.Attributes["ids.0"]
		if !ok {
			return fmt.Errorf("attribute not found: ids.0")
		}
		id2, ok := listDs.Primary.Attributes["ids.1"]
		if !ok {
			return fmt.Errorf("attribute not found: ids.1")
		}

		foundPortalID := false
		foundSelfServiceID := false
		if id1 == portalID || id2 == portalID {
			foundPortalID = true
		}
		if id1 == selfServiceID || id2 == selfServiceID {
			foundSelfServiceID = true
		}

		if !foundPortalID {
			return fmt.Errorf("portal application ID from single application data source not found in list data source IDs")
		}
		if !foundSelfServiceID {
			return fmt.Errorf("self-service application ID from single application data source not found in list data source IDs")
		}

		return nil
	}
}
