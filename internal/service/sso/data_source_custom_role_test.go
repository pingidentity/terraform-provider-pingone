package sso_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/sso"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

const (
	// Default "Custom Roles Admin" administrator role id
	testAccCustomRoleDataSource_CustomRolesAdminRoleID = "6f770b08-793f-4393-b2aa-b1d1587a0324"
	// Default "Environment Admin" administrator role id
	testAccCustomRoleDataSource_EnvironmentAdminRoleID = "29ddce68-cd7f-4b2a-b6fc-f7a19553b496"
)

func TestAccCustomRoleDataSource_ByNameFull(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_custom_role.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.CustomRole_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCustomRoleDataSourceConfig_ByNameFull(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceFullName, "applicable_to.#", "2"),
					resource.TestCheckTypeSetElemAttr(dataSourceFullName, "applicable_to.*", "ENVIRONMENT"),
					resource.TestCheckTypeSetElemAttr(dataSourceFullName, "applicable_to.*", "POPULATION"),
					resource.TestCheckResourceAttr(dataSourceFullName, "can_assign.#", "1"),
					resource.TestMatchResourceAttr(dataSourceFullName, "can_assign.0.id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourceFullName, "can_be_assigned_by.*",
						map[string]string{
							"id": testAccCustomRoleDataSource_CustomRolesAdminRoleID,
						},
					),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourceFullName, "can_be_assigned_by.*",
						map[string]string{
							"id": testAccCustomRoleDataSource_EnvironmentAdminRoleID,
						},
					),
					resource.TestCheckResourceAttr(dataSourceFullName, "description", "My custom role for datasource test"),
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourceFullName, "permissions.*",
						map[string]string{
							"id": "permissions:read:gatewayRoleAssignments",
						},
					),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourceFullName, "permissions.*",
						map[string]string{
							"id": "permissions:update:userRoleAssignments",
						},
					),
					resource.TestMatchResourceAttr(dataSourceFullName, "role_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "type", "CUSTOM"),
				),
			},
		},
	})
}

func TestAccCustomRoleDataSource_ByIDFull(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_custom_role.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.CustomRole_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCustomRoleDataSourceConfig_ByIDFull(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceFullName, "applicable_to.#", "2"),
					resource.TestCheckTypeSetElemAttr(dataSourceFullName, "applicable_to.*", "ENVIRONMENT"),
					resource.TestCheckTypeSetElemAttr(dataSourceFullName, "applicable_to.*", "POPULATION"),
					resource.TestCheckResourceAttr(dataSourceFullName, "can_assign.#", "1"),
					resource.TestMatchResourceAttr(dataSourceFullName, "can_assign.0.id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourceFullName, "can_be_assigned_by.*",
						map[string]string{
							"id": testAccCustomRoleDataSource_CustomRolesAdminRoleID,
						},
					),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourceFullName, "can_be_assigned_by.*",
						map[string]string{
							"id": testAccCustomRoleDataSource_EnvironmentAdminRoleID,
						},
					),
					resource.TestCheckResourceAttr(dataSourceFullName, "description", "My custom role for datasource test"),
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "name", name),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourceFullName, "permissions.*",
						map[string]string{
							"id": "permissions:read:gatewayRoleAssignments",
						},
					),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourceFullName, "permissions.*",
						map[string]string{
							"id": "permissions:update:userRoleAssignments",
						},
					),
					resource.TestCheckResourceAttr(dataSourceFullName, "type", "CUSTOM"),
				),
			},
		},
	})
}

func TestAccCustomRoleDataSource_NotFound(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.CustomRole_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccCustomRoleDataSourceConfig_NotFoundByName(resourceName),
				ExpectError: regexp.MustCompile(`Error: Cannot find the custom role from name`),
			},
			{
				Config:      testAccCustomRoleDataSourceConfig_NotFoundByID(resourceName),
				ExpectError: regexp.MustCompile("There is no role with id"),
			},
		},
	})
}

func testAccCustomRoleDataSourceConfig_Full(resourceName, name string) string {
	return fmt.Sprintf(`
resource "pingone_custom_role" "%[1]s-dependent-role" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[2]s Datasource Dependent Role"
  applicable_to = [
    "ENVIRONMENT",
    "POPULATION"
  ]
  can_be_assigned_by = [
    {
      id = pingone_custom_role.%[1]s-parent.id
    }
  ]
  description = "My custom dependent role for datasource test"
  permissions = [
    {
      id = "permissions:read:userRoleAssignments"
    },
    {
      id = "permissions:read:groupRoleAssignments"
    },
  ]
}

resource "pingone_custom_role" "%[1]s-parent" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[2]s"
  applicable_to = [
    "ENVIRONMENT",
    "POPULATION"
  ]
  can_be_assigned_by = [
    {
      id = "%[3]s"
    },
    {
      id = "%[4]s"
    }
  ]
  description = "My custom role for datasource test"
  permissions = [
    {
      id = "permissions:read:gatewayRoleAssignments"
    },
    {
      id = "permissions:update:userRoleAssignments"
    }
  ]
}
	`, resourceName, name,
		testAccCustomRoleDataSource_CustomRolesAdminRoleID, testAccCustomRoleDataSource_EnvironmentAdminRoleID)
}

func testAccCustomRoleDataSourceConfig_ByNameFull(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

	%[3]s

data "pingone_custom_role" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  depends_on     = [pingone_custom_role.%[2]s-dependent-role]

  name = pingone_custom_role.%[2]s-parent.name
}
`, acctest.GenericSandboxEnvironment(), resourceName, testAccCustomRoleDataSourceConfig_Full(resourceName, name))
}

func testAccCustomRoleDataSourceConfig_ByIDFull(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

    %[3]s

data "pingone_custom_role" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  depends_on     = [pingone_custom_role.%[2]s-dependent-role]

  role_id = pingone_custom_role.%[2]s-parent.id
}`, acctest.GenericSandboxEnvironment(), resourceName, testAccCustomRoleDataSourceConfig_Full(resourceName, name))
}

func testAccCustomRoleDataSourceConfig_NotFoundByName(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_custom_role" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "doesnotexist"
}
`, acctest.GenericSandboxEnvironment(), resourceName)
}

func testAccCustomRoleDataSourceConfig_NotFoundByID(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_custom_role" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  role_id = "9c052a8a-14be-44e4-8f07-2662569994ce" // dummy ID that conforms to UUID v4
}
`, acctest.GenericSandboxEnvironment(), resourceName)
}
