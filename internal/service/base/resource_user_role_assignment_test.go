// Copyright © 2025 Ping Identity Corporation

package base_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/base"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/sso"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccRoleAssignmentUser_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_user_role_assignment.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	var roleAssignmentID, userID, environmentID string

	var p1Client *client.Client
	var ctx = context.Background()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)

			p1Client = acctest.PreCheckTestClient(ctx, t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.RoleAssignmentUser_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Test removal of the resource
			{
				Config: testAccRoleAssignmentUserConfig_Population(resourceName, name, "Identity Data Admin"),
				Check:  base.RoleAssignmentUser_GetIDs(resourceFullName, &environmentID, &userID, &roleAssignmentID),
			},
			{
				PreConfig: func() {
					base.RoleAssignmentUser_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID, userID, roleAssignmentID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the user
			{
				Config: testAccRoleAssignmentUserConfig_Population(resourceName, name, "Identity Data Admin"),
				Check:  base.RoleAssignmentUser_GetIDs(resourceFullName, &environmentID, &userID, &roleAssignmentID),
			},
			{
				PreConfig: func() {
					sso.User_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID, userID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the environment
			{
				Config: testAccRoleAssignmentUserConfig_NewEnv(environmentName, licenseID, resourceName, name, "Identity Data Admin"),
				Check:  base.RoleAssignmentUser_GetIDs(resourceFullName, &environmentID, &userID, &roleAssignmentID),
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

func TestAccRoleAssignmentUser_Application(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_user_role_assignment.%s", resourceName)

	name := resourceName

	successCheck := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "user_id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "role_id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "scope_application_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckNoResourceAttr(resourceFullName, "scope_population_id"),
		resource.TestCheckNoResourceAttr(resourceFullName, "scope_organization_id"),
		resource.TestCheckNoResourceAttr(resourceFullName, "scope_environment_id"),
		resource.TestCheckResourceAttr(resourceFullName, "read_only", "false"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.RoleAssignmentUser_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRoleAssignmentUserConfig_Application(resourceName, name, "Application Owner"),
				Check:  successCheck,
			},
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["user_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      testAccRoleAssignmentUserConfig_Application(resourceName, name, "Identity Data Admin"),
				ExpectError: regexp.MustCompile(`Incompatible role and scope combination`),
			},
			{
				Config:      testAccRoleAssignmentUserConfig_Application(resourceName, name, "Environment Admin"),
				ExpectError: regexp.MustCompile(`Incompatible role and scope combination`),
			},
			{
				Config:      testAccRoleAssignmentUserConfig_Application(resourceName, name, "Organization Admin"),
				ExpectError: regexp.MustCompile(`Incompatible role and scope combination`),
			},
		},
	})
}

func TestAccRoleAssignmentUser_Population(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_user_role_assignment.%s", resourceName)

	name := resourceName

	successCheck := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "user_id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "role_id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "scope_population_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckNoResourceAttr(resourceFullName, "scope_organization_id"),
		resource.TestCheckNoResourceAttr(resourceFullName, "scope_environment_id"),
		resource.TestCheckResourceAttr(resourceFullName, "read_only", "false"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.RoleAssignmentUser_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRoleAssignmentUserConfig_Population(resourceName, name, "Identity Data Admin"),
				Check:  successCheck,
			},
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["user_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      testAccRoleAssignmentUserConfig_Population(resourceName, name, "Application Owner"),
				ExpectError: regexp.MustCompile(`Incompatible role and scope combination`),
			},
			{
				Config:      testAccRoleAssignmentUserConfig_Population(resourceName, name, "Environment Admin"),
				ExpectError: regexp.MustCompile(`Incompatible role and scope combination`),
			},
			{
				Config:      testAccRoleAssignmentUserConfig_Population(resourceName, name, "Organization Admin"),
				ExpectError: regexp.MustCompile(`Incompatible role and scope combination`),
			},
		},
	})
}

func TestAccRoleAssignmentUser_Organisation(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_user_role_assignment.%s", resourceName)

	name := resourceName
	organisationID := os.Getenv("PINGONE_ORGANIZATION_ID")

	successCheck := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "user_id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "role_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckNoResourceAttr(resourceFullName, "scope_population_id"),
		resource.TestMatchResourceAttr(resourceFullName, "scope_organization_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckNoResourceAttr(resourceFullName, "scope_environment_id"),
		resource.TestCheckResourceAttr(resourceFullName, "read_only", "false"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.RoleAssignmentUser_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccRoleAssignmentUserConfig_Organisation(resourceName, name, "Identity Data Admin", organisationID),
				ExpectError: regexp.MustCompile(`Incompatible role and scope combination`),
			},
			{
				Config:      testAccRoleAssignmentUserConfig_Organisation(resourceName, name, "Application Owner", organisationID),
				ExpectError: regexp.MustCompile(`Incompatible role and scope combination`),
			},
			{
				Config: testAccRoleAssignmentUserConfig_Organisation(resourceName, name, "Environment Admin", organisationID),
				Check:  successCheck,
			},
			{
				Config: testAccRoleAssignmentUserConfig_Organisation(resourceName, name, "Organization Admin", organisationID),
				Check:  successCheck,
			},
			{
				Config: testAccRoleAssignmentUserConfig_Organisation(resourceName, name, "DaVinci Admin", organisationID),
				Check:  successCheck,
			},
			{
				Config: testAccRoleAssignmentUserConfig_Organisation(resourceName, name, "DaVinci Admin Read Only", organisationID),
				Check:  successCheck,
			},
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["user_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccRoleAssignmentUser_Environment(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_user_role_assignment.%s", resourceName)

	name := resourceName

	successCheck := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "user_id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "role_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckNoResourceAttr(resourceFullName, "scope_population_id"),
		resource.TestCheckNoResourceAttr(resourceFullName, "scope_organization_id"),
		resource.TestMatchResourceAttr(resourceFullName, "scope_environment_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "read_only", "false"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
			acctest.PreCheckOrganisationID(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.RoleAssignmentUser_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRoleAssignmentUserConfig_Environment(resourceName, name, "Identity Data Admin"),
				Check:  successCheck,
			},
			{
				Config: testAccRoleAssignmentUserConfig_Environment(resourceName, name, "Environment Admin"),
				Check:  successCheck,
			},
			{
				Config: testAccRoleAssignmentUserConfig_Environment(resourceName, name, "DaVinci Admin"),
				Check:  successCheck,
			},
			{
				Config: testAccRoleAssignmentUserConfig_Environment(resourceName, name, "DaVinci Admin Read Only"),
				Check:  successCheck,
			},
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["user_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccRoleAssignmentUserConfig_Environment(resourceName, name, "Application Owner"),
				Check:  successCheck,
			},
			{
				Config:      testAccRoleAssignmentUserConfig_Environment(resourceName, name, "Organization Admin"),
				ExpectError: regexp.MustCompile(`Incompatible role and scope combination`),
			},
		},
	})
}

func TestAccRoleAssignmentUser_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_user_role_assignment.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.RoleAssignmentUser_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccRoleAssignmentUserConfig_Population(resourceName, name, "Identity Data Admin"),
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

func testAccRoleAssignmentUserConfig_NewEnv(environmentName, licenseID, resourceName, name, roleName string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_population" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  name = "%[4]s"
}

resource "pingone_user" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  population_id  = pingone_population.%[3]s.id

  username = "%[4]s"
  email    = "foouser@pingidentity.com"
}

data "pingone_role" "%[3]s" {
  name = "%[5]s"
}

resource "pingone_user_role_assignment" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  user_id        = pingone_user.%[3]s.id
  role_id        = data.pingone_role.%[3]s.id

  scope_population_id = pingone_population.%[3]s.id
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name, roleName)
}

func testAccRoleAssignmentUserConfig_Application(resourceName, name, roleName string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true

  oidc_options = {
    type                       = "WEB_APP"
    grant_types                = ["REFRESH_TOKEN", "AUTHORIZATION_CODE"]
    response_types             = ["CODE"]
    token_endpoint_auth_method = "CLIENT_SECRET_BASIC"
    redirect_uris              = ["https://www.pingidentity.com"]
  }
}

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

resource "pingone_user_role_assignment" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  user_id        = pingone_user.%[2]s.id
  role_id        = data.pingone_role.%[2]s.id

  scope_application_id = pingone_application.%[2]s.id
}`, acctest.GenericSandboxEnvironment(), resourceName, name, roleName)
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

resource "pingone_user_role_assignment" "%[2]s" {
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

resource "pingone_user_role_assignment" "%[2]s" {
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

resource "pingone_user_role_assignment" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  user_id        = pingone_user.%[2]s.id
  role_id        = data.pingone_role.%[2]s.id

  scope_environment_id = data.pingone_environment.general_test.id
}`, acctest.GenericSandboxEnvironment(), resourceName, name, roleName)
}
