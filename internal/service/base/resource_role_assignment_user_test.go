package base_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func testAccCheckRoleAssignmentUserDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_role_assignment_user" {
			continue
		}

		_, rEnv, err := apiClient.EnvironmentsApi.ReadOneEnvironment(ctx, rs.Primary.Attributes["environment_id"]).Execute()

		if err != nil {

			if rEnv == nil {
				return fmt.Errorf("Response object does not exist and no error detected")
			}

			if rEnv.StatusCode == 404 {
				continue
			}

			return err
		}

		body, r, err := apiClient.UserRoleAssignmentsApi.ReadOneUserRoleAssignment(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["user_id"], rs.Primary.ID).Execute()

		if err != nil {

			if r == nil {
				return fmt.Errorf("Response object does not exist and no error detected")
			}

			if r.StatusCode == 404 {
				continue
			}

			tflog.Error(ctx, fmt.Sprintf("Error: %v", body))
			return err
		}

		return fmt.Errorf("PingOne User Role Assignment %s still exists", rs.Primary.ID)
	}

	return nil
}

func TestAccRoleAssignmentUser_Population(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_role_assignment_user.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckRoleAssignmentUserDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRoleAssignmentUserConfig_Population(resourceName, name, "Identity Data Admin"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "user_id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "role_id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "scope_population_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "scope_organization_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "scope_environment_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "read_only", "false"),
				),
			},
			{
				Config:      testAccRoleAssignmentUserConfig_Population(resourceName, name, "Environment Admin"),
				ExpectError: regexp.MustCompile(`Incompatible role and scope combination. Role: [a-z0-9\-]* \/ Scope: POPULATION`),
			},
			{
				Config:      testAccRoleAssignmentUserConfig_Population(resourceName, name, "Organization Admin"),
				ExpectError: regexp.MustCompile(`Incompatible role and scope combination. Role: [a-z0-9\-]* \/ Scope: POPULATION`),
			},
		},
	})
}

func TestAccRoleAssignmentUser_Organisation(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_role_assignment_user.%s", resourceName)

	name := resourceName
	organisationID := os.Getenv("PINGONE_ORGANIZATION_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckRoleAssignmentUserDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccRoleAssignmentUserConfig_Organisation(resourceName, name, "Identity Data Admin", organisationID),
				ExpectError: regexp.MustCompile(`Incompatible role and scope combination. Role: [a-z0-9\-]* \/ Scope: ORGANIZATION`),
			},
			{
				Config: testAccRoleAssignmentUserConfig_Organisation(resourceName, name, "Environment Admin", organisationID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "user_id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "role_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "scope_population_id", ""),
					resource.TestMatchResourceAttr(resourceFullName, "scope_organization_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "scope_environment_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "read_only", "false"),
				),
			},
			{
				Config: testAccRoleAssignmentUserConfig_Organisation(resourceName, name, "Organization Admin", organisationID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "user_id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "role_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "scope_population_id", ""),
					resource.TestMatchResourceAttr(resourceFullName, "scope_organization_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "scope_environment_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "read_only", "false"),
				),
			},
		},
	})
}

func TestAccRoleAssignmentUser_Environment(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_role_assignment_user.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironmentAndOrganisation(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckRoleAssignmentUserDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRoleAssignmentUserConfig_Environment(resourceName, name, "Identity Data Admin"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "user_id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "role_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "scope_population_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "scope_organization_id", ""),
					resource.TestMatchResourceAttr(resourceFullName, "scope_environment_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "read_only", "false"),
				),
			},
			{
				Config: testAccRoleAssignmentUserConfig_Environment(resourceName, name, "Environment Admin"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "user_id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "role_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "scope_population_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "scope_organization_id", ""),
					resource.TestMatchResourceAttr(resourceFullName, "scope_environment_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "read_only", "false"),
				),
			},
			{
				Config:      testAccRoleAssignmentUserConfig_Environment(resourceName, name, "Organization Admin"),
				ExpectError: regexp.MustCompile(`Incompatible role and scope combination. Role: [a-z0-9\-]* \/ Scope: ENVIRONMENT`),
			},
		},
	})
}

func testAccRoleAssignmentUserConfig_Population(resourceName, name, roleName string) string {
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

data "pingone_role" "%[2]s" {
  name = "%[4]s"
}

resource "pingone_role_assignment_user" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  user_id        = pingone_user.%[2]s.id
  role_id        = data.pingone_role.%[2]s.id

  scope_population_id = pingone_population.%[2]s.id
}`, acctest.GenericSandboxEnvironment(), resourceName, name, roleName)
}

func testAccRoleAssignmentUserConfig_Organisation(resourceName, name, roleName, organisationID string) string {
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

data "pingone_role" "%[2]s" {
  name = "%[4]s"
}

resource "pingone_role_assignment_user" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  user_id        = pingone_user.%[2]s.id
  role_id        = data.pingone_role.%[2]s.id

  scope_organization_id = "%[5]s"
}`, acctest.GenericSandboxEnvironment(), resourceName, name, roleName, organisationID)
}

func testAccRoleAssignmentUserConfig_Environment(resourceName, name, roleName string) string {
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

data "pingone_role" "%[2]s" {
  name = "%[4]s"
}

resource "pingone_role_assignment_user" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  user_id        = pingone_user.%[2]s.id
  role_id        = data.pingone_role.%[2]s.id

  scope_environment_id = data.pingone_environment.general_test.id
}`, acctest.GenericSandboxEnvironment(), resourceName, name, roleName)
}
