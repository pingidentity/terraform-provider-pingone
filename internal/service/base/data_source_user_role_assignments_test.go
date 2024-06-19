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

func TestAccUserRoleAssignmentsDataSource_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_user_role_assignments.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	name := resourceName

	organisationID := os.Getenv("PINGONE_ORGANIZATION_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.RoleAssignmentUser_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccUserRoleAssignmentsDataSourceConfig_Full(resourceName, name, organisationID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "user_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "role_assignments.#", "3"),
					resource.TestMatchTypeSetElemNestedAttrs(dataSourceFullName, "role_assignments.*", map[string]*regexp.Regexp{
						"id":         verify.P1ResourceIDRegexpFullString,
						"scope.id":   verify.P1ResourceIDRegexpFullString,
						"scope.type": regexp.MustCompile(`^ENVIRONMENT$`),
						"role_id":    verify.P1ResourceIDRegexpFullString,
						"read_only":  regexp.MustCompile(`^false$`),
					}),
					resource.TestMatchTypeSetElemNestedAttrs(dataSourceFullName, "role_assignments.*", map[string]*regexp.Regexp{
						"id":         verify.P1ResourceIDRegexpFullString,
						"scope.id":   verify.P1ResourceIDRegexpFullString,
						"scope.type": regexp.MustCompile(`^POPULATION$`),
						"role_id":    verify.P1ResourceIDRegexpFullString,
						"read_only":  regexp.MustCompile(`^false$`),
					}),
					resource.TestMatchTypeSetElemNestedAttrs(dataSourceFullName, "role_assignments.*", map[string]*regexp.Regexp{
						"id":         verify.P1ResourceIDRegexpFullString,
						"scope.id":   verify.P1ResourceIDRegexpFullString,
						"scope.type": regexp.MustCompile(`^ORGANIZATION$`),
						"role_id":    verify.P1ResourceIDRegexpFullString,
						"read_only":  regexp.MustCompile(`^false$`),
					}),
				),
			},
		},
	})
}

func TestAccUserRoleAssignmentsDataSource_NotFound(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_user_role_assignments.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.RoleAssignmentUser_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccUserRoleAssignmentsDataSourceConfig_NoRoles(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "user_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "role_assignments.#", "0"),
				),
			},
		},
	})
}

func testAccUserRoleAssignmentsDataSourceConfig_Full(resourceName, name, organisationID string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_population" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_user" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  population_id  = pingone_population.%[2]s.id

  username = "%[3]s"
  email    = "foouser@pingidentity.com"
}

data "pingone_role" "environment_admin" {
  name = "Environment Admin"
}

data "pingone_role" "client_application_developer" {
  name = "Client Application Developer"
}

data "pingone_role" "identity_data_admin" {
  name = "Identity Data Admin"
}

resource "pingone_user_role_assignment" "%[2]s-1" {
  environment_id = data.pingone_environment.general_test.id
  user_id        = pingone_user.%[2]s.id
  role_id        = data.pingone_role.environment_admin.id

  scope_organization_id = "%[4]s"
}

resource "pingone_user_role_assignment" "%[2]s-2" {
  environment_id = data.pingone_environment.general_test.id
  user_id        = pingone_user.%[2]s.id
  role_id        = data.pingone_role.client_application_developer.id

  scope_environment_id = data.pingone_environment.general_test.id
}

resource "pingone_user_role_assignment" "%[2]s-3" {
  environment_id = data.pingone_environment.general_test.id
  user_id        = pingone_user.%[2]s.id
  role_id        = data.pingone_role.identity_data_admin.id

  scope_population_id = pingone_population.%[2]s.id
}

data "pingone_user_role_assignments" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  user_id = pingone_user.%[2]s.id

  depends_on = [
    pingone_user_role_assignment.%[2]s-1,
    pingone_user_role_assignment.%[2]s-2,
    pingone_user_role_assignment.%[2]s-3,
  ]
}
`, acctest.GenericSandboxEnvironment(), resourceName, name, organisationID)
}

func testAccUserRoleAssignmentsDataSourceConfig_NoRoles(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_population" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_user" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  population_id  = pingone_population.%[2]s.id

  username = "%[3]s"
  email    = "foouser@pingidentity.com"
}

data "pingone_user_role_assignments" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  user_id = pingone_user.%[2]s.id
}
`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
