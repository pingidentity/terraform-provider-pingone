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
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/base"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/sso"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccRoleAssignmentGroup_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_group_role_assignment.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	var roleAssignmentID, groupID, environmentID string

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
		CheckDestroy:             sso.RoleAssignmentGroup_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Test removal of the resource
			{
				Config: testAccRoleAssignmentGroupConfig_Population(resourceName, name, "Identity Data Admin"),
				Check:  sso.RoleAssignmentGroup_GetIDs(resourceFullName, &environmentID, &groupID, &roleAssignmentID),
			},
			{
				PreConfig: func() {
					sso.RoleAssignmentGroup_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID, groupID, roleAssignmentID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the application
			{
				Config: testAccRoleAssignmentGroupConfig_Population(resourceName, name, "Identity Data Admin"),
				Check:  sso.RoleAssignmentGroup_GetIDs(resourceFullName, &environmentID, &groupID, &roleAssignmentID),
			},
			{
				PreConfig: func() {
					sso.Group_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID, groupID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the environment
			{
				Config: testAccRoleAssignmentGroupConfig_NewEnv(environmentName, licenseID, resourceName, name, "Identity Data Admin"),
				Check:  sso.RoleAssignmentGroup_GetIDs(resourceFullName, &environmentID, &groupID, &roleAssignmentID),
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

func TestAccRoleAssignmentGroup_Application(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_group_role_assignment.%s", resourceName)

	name := resourceName

	successCheck := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "group_id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "role_id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "scope_application_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckNoResourceAttr(resourceFullName, "scope_population_id"),
		resource.TestCheckNoResourceAttr(resourceFullName, "scope_organization_id"),
		resource.TestCheckNoResourceAttr(resourceFullName, "scope_environment_id"),
		resource.TestCheckResourceAttr(resourceFullName, "read_only", "false"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.RoleAssignmentGroup_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRoleAssignmentGroupConfig_Application(resourceName, name, "Application Owner"),
				Check:  successCheck,
			},
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["group_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      testAccRoleAssignmentGroupConfig_Application(resourceName, name, "Identity Data Admin"),
				ExpectError: regexp.MustCompile(`Incompatible role and scope combination`),
			},
			{
				Config:      testAccRoleAssignmentGroupConfig_Application(resourceName, name, "Environment Admin"),
				ExpectError: regexp.MustCompile(`Incompatible role and scope combination`),
			},
			{
				Config:      testAccRoleAssignmentGroupConfig_Application(resourceName, name, "Organization Admin"),
				ExpectError: regexp.MustCompile(`Incompatible role and scope combination`),
			},
		},
	})
}

func TestAccRoleAssignmentGroup_Population(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_group_role_assignment.%s", resourceName)

	name := resourceName

	successCheck := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "group_id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "role_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckNoResourceAttr(resourceFullName, "scope_application_id"),
		resource.TestMatchResourceAttr(resourceFullName, "scope_population_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckNoResourceAttr(resourceFullName, "scope_organization_id"),
		resource.TestCheckNoResourceAttr(resourceFullName, "scope_environment_id"),
		resource.TestCheckResourceAttr(resourceFullName, "read_only", "false"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.RoleAssignmentGroup_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRoleAssignmentGroupConfig_Population(resourceName, name, "Identity Data Admin"),
				Check:  successCheck,
			},
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["group_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      testAccRoleAssignmentGroupConfig_Population(resourceName, name, "Application Owner"),
				ExpectError: regexp.MustCompile(`Incompatible role and scope combination`),
			},
			{
				Config:      testAccRoleAssignmentGroupConfig_Population(resourceName, name, "Environment Admin"),
				ExpectError: regexp.MustCompile(`Incompatible role and scope combination`),
			},
			{
				Config:      testAccRoleAssignmentGroupConfig_Population(resourceName, name, "Organization Admin"),
				ExpectError: regexp.MustCompile(`Incompatible role and scope combination`),
			},
		},
	})
}

func TestAccRoleAssignmentGroup_Organisation(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_group_role_assignment.%s", resourceName)

	name := resourceName
	organisationID := os.Getenv("PINGONE_ORGANIZATION_ID")

	successCheck := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "group_id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "role_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckNoResourceAttr(resourceFullName, "scope_application_id"),
		resource.TestCheckNoResourceAttr(resourceFullName, "scope_population_id"),
		resource.TestMatchResourceAttr(resourceFullName, "scope_organization_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckNoResourceAttr(resourceFullName, "scope_environment_id"),
		resource.TestCheckResourceAttr(resourceFullName, "read_only", "false"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.RoleAssignmentGroup_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccRoleAssignmentGroupConfig_Organisation(resourceName, name, "Application Owner", organisationID),
				ExpectError: regexp.MustCompile(`Incompatible role and scope combination`),
			},
			{
				Config:      testAccRoleAssignmentGroupConfig_Organisation(resourceName, name, "Identity Data Admin", organisationID),
				ExpectError: regexp.MustCompile(`Incompatible role and scope combination`),
			},
			{
				Config: testAccRoleAssignmentGroupConfig_Organisation(resourceName, name, "Environment Admin", organisationID),
				Check:  successCheck,
			},
			{
				Config: testAccRoleAssignmentGroupConfig_Organisation(resourceName, name, "Organization Admin", organisationID),
				Check:  successCheck,
			},
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["group_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccRoleAssignmentGroup_Environment(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_group_role_assignment.%s", resourceName)

	name := resourceName

	successCheck := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "group_id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "role_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckNoResourceAttr(resourceFullName, "scope_application_id"),
		resource.TestCheckNoResourceAttr(resourceFullName, "scope_population_id"),
		resource.TestCheckNoResourceAttr(resourceFullName, "scope_organization_id"),
		resource.TestMatchResourceAttr(resourceFullName, "scope_environment_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "read_only", "false"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckOrganisationID(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.RoleAssignmentGroup_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRoleAssignmentGroupConfig_Environment(resourceName, name, "Application Owner"),
				Check:  successCheck,
			},
			{
				Config: testAccRoleAssignmentGroupConfig_Environment(resourceName, name, "Identity Data Admin"),
				Check:  successCheck,
			},
			{
				Config: testAccRoleAssignmentGroupConfig_Environment(resourceName, name, "Environment Admin"),
				Check:  successCheck,
			},
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["group_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      testAccRoleAssignmentGroupConfig_Environment(resourceName, name, "Organization Admin"),
				ExpectError: regexp.MustCompile(`Incompatible role and scope combination`),
			},
		},
	})
}

func TestAccRoleAssignmentGroup_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_group_role_assignment.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.RoleAssignmentGroup_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccRoleAssignmentGroupConfig_Population(resourceName, name, "Identity Data Admin"),
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

func testAccRoleAssignmentGroupConfig_NewEnv(environmentName, licenseID, resourceName, name, roleName string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_population" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  name = "%[4]s"
}

resource "pingone_group" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s"
}

data "pingone_role" "%[3]s" {
  name = "%[5]s"
}

resource "pingone_group_role_assignment" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  group_id       = pingone_group.%[3]s.id
  role_id        = data.pingone_role.%[3]s.id

  scope_population_id = pingone_population.%[3]s.id
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name, roleName)
}

func testAccRoleAssignmentGroupConfig_Application(resourceName, name, roleName string) string {
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

resource "pingone_group" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
}

data "pingone_role" "%[2]s" {
  name = "%[4]s"
}

resource "pingone_group_role_assignment" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  group_id       = pingone_group.%[2]s.id
  role_id        = data.pingone_role.%[2]s.id

  scope_application_id = pingone_application.%[2]s.id
}`, acctest.GenericSandboxEnvironment(), resourceName, name, roleName)
}

func testAccRoleAssignmentGroupConfig_Population(resourceName, name, roleName string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_population" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_group" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
}

data "pingone_role" "%[2]s" {
  name = "%[4]s"
}

resource "pingone_group_role_assignment" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  group_id       = pingone_group.%[2]s.id
  role_id        = data.pingone_role.%[2]s.id

  scope_population_id = pingone_population.%[2]s.id
}`, acctest.GenericSandboxEnvironment(), resourceName, name, roleName)
}

func testAccRoleAssignmentGroupConfig_Organisation(resourceName, name, roleName, organisationID string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_group" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
}

data "pingone_role" "%[2]s" {
  name = "%[4]s"
}

resource "pingone_group_role_assignment" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  group_id       = pingone_group.%[2]s.id
  role_id        = data.pingone_role.%[2]s.id

  scope_organization_id = "%[5]s"
}`, acctest.GenericSandboxEnvironment(), resourceName, name, roleName, organisationID)
}

func testAccRoleAssignmentGroupConfig_Environment(resourceName, name, roleName string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_group" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
}

data "pingone_role" "%[2]s" {
  name = "%[4]s"
}

resource "pingone_group_role_assignment" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  group_id       = pingone_group.%[2]s.id
  role_id        = data.pingone_role.%[2]s.id

  scope_environment_id = data.pingone_environment.general_test.id
}`, acctest.GenericSandboxEnvironment(), resourceName, name, roleName)
}
