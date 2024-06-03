package sso_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/authorize"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/base"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/sso"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccUserApplicationRoleAssignment_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_user_application_role_assignment.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	var applicationRoleID, userID, environmentID string

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
		CheckDestroy:             sso.UserApplicationRoleAssignment_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Test removal of the resource
			{
				Config: testAccUserApplicationRoleAssignmentConfig_Full(resourceName, name),
				Check:  sso.UserApplicationRoleAssignment_GetIDs(resourceFullName, &environmentID, &userID, &applicationRoleID),
			},
			{
				PreConfig: func() {
					sso.UserApplicationRoleAssignment_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID, userID, applicationRoleID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the user
			{
				Config: testAccUserApplicationRoleAssignmentConfig_Full(resourceName, name),
				Check:  sso.UserApplicationRoleAssignment_GetIDs(resourceFullName, &environmentID, &userID, &applicationRoleID),
			},
			{
				PreConfig: func() {
					sso.User_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID, userID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the application role
			{
				Config: testAccUserApplicationRoleAssignmentConfig_Full(resourceName, name),
				Check:  sso.UserApplicationRoleAssignment_GetIDs(resourceFullName, &environmentID, &userID, &applicationRoleID),
			},
			{
				PreConfig: func() {
					authorize.ApplicationRole_RemovalDrift_PreConfig(ctx, p1Client.API.AuthorizeAPIClient, t, environmentID, applicationRoleID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the environment
			{
				Config: testAccUserApplicationRoleAssignmentConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check:  sso.UserApplicationRoleAssignment_GetIDs(resourceFullName, &environmentID, &userID, &applicationRoleID),
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

func TestAccUserApplicationRoleAssignment_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_user_application_role_assignment.%s", resourceName)

	name := resourceName

	fullStep := resource.TestStep{
		Config: testAccUserApplicationRoleAssignmentConfig_Full(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "user_id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "application_role_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestCheckResourceAttr(resourceFullName, "description", "Test application role"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.UserApplicationRoleAssignment_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			fullStep,
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["user_id"], rs.Primary.Attributes["application_role_id"]), nil
					}
				}(),
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "application_role_id",
			},
		},
	})
}

func TestAccUserApplicationRoleAssignment_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_user_application_role_assignment.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.UserApplicationRoleAssignment_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccUserApplicationRoleAssignmentConfig_Full(resourceName, name),
			},
			// Errors
			{
				ResourceName: resourceFullName,
				ImportState:  true,
				ExpectError:  regexp.MustCompile(`Unexpected Import Identifier`),
			},
			{
				ResourceName:  resourceFullName,
				ImportStateId: "/",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`Unexpected Import Identifier`),
			},
			{
				ResourceName:  resourceFullName,
				ImportStateId: "badformat/badformat/badformat",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`Unexpected Import Identifier`),
			},
		},
	})
}

func testAccUserApplicationRoleAssignmentConfig_NewEnv(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_population" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  name = "%[4]s"
}

resource "pingone_user" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  username      = "%[4]s"
  email         = "%[4]s@pingidentity.com"
  population_id = pingone_population.%[3]s.id
}

resource "pingone_authorize_application_role" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s"
  description    = "Test application role"
}

resource "pingone_user_application_role_assignment" "%[3]s" {
  environment_id      = pingone_environment.%[2]s.id
  user_id             = pingone_user.%[3]s.id
  application_role_id = pingone_authorize_application_role.%[3]s.id
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccUserApplicationRoleAssignmentConfig_Full(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_population" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_user" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  username      = "%[3]s"
  email         = "%[3]s@pingidentity.com"
  population_id = pingone_population.%[2]s.id
}

resource "pingone_authorize_application_role" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Test application role"
}

resource "pingone_user_application_role_assignment" "%[2]s" {
  environment_id      = data.pingone_environment.general_test.id
  user_id             = pingone_user.%[2]s.id
  application_role_id = pingone_authorize_application_role.%[2]s.id
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
