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

func TestAccRoleAssignmentApplication_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application_role_assignment.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	var roleAssignmentID, applicationID, environmentID string

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
		CheckDestroy:             sso.RoleAssignmentApplication_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Test removal of the resource
			{
				Config: testAccRoleAssignmentApplicationConfig_Population(resourceName, name, "Identity Data Admin"),
				Check:  sso.RoleAssignmentApplication_GetIDs(resourceFullName, &environmentID, &applicationID, &roleAssignmentID),
			},
			{
				PreConfig: func() {
					sso.RoleAssignmentApplication_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID, applicationID, roleAssignmentID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the application
			{
				Config: testAccRoleAssignmentApplicationConfig_Population(resourceName, name, "Identity Data Admin"),
				Check:  sso.RoleAssignmentApplication_GetIDs(resourceFullName, &environmentID, &applicationID, &roleAssignmentID),
			},
			{
				PreConfig: func() {
					sso.Application_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID, applicationID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the environment
			{
				Config: testAccRoleAssignmentApplicationConfig_NewEnv(environmentName, licenseID, resourceName, name, "Identity Data Admin"),
				Check:  sso.RoleAssignmentApplication_GetIDs(resourceFullName, &environmentID, &applicationID, &roleAssignmentID),
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

func TestAccRoleAssignmentApplication_Population(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application_role_assignment.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.RoleAssignmentApplication_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRoleAssignmentApplicationConfig_Population(resourceName, name, "Identity Data Admin"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "application_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "role_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "scope_population_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckNoResourceAttr(resourceFullName, "scope_organization_id"),
					resource.TestCheckNoResourceAttr(resourceFullName, "scope_environment_id"),
					resource.TestCheckResourceAttr(resourceFullName, "read_only", "false"),
				),
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

						return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["application_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      testAccRoleAssignmentApplicationConfig_Population(resourceName, name, "Environment Admin"),
				ExpectError: regexp.MustCompile(`Incompatible role and scope combination`),
			},
			{
				Config:      testAccRoleAssignmentApplicationConfig_Population(resourceName, name, "Organization Admin"),
				ExpectError: regexp.MustCompile(`Incompatible role and scope combination`),
			},
		},
	})
}

func TestAccRoleAssignmentApplication_Organisation(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application_role_assignment.%s", resourceName)

	name := resourceName
	organisationID := os.Getenv("PINGONE_ORGANIZATION_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.RoleAssignmentApplication_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccRoleAssignmentApplicationConfig_Organisation(resourceName, name, "Identity Data Admin", organisationID),
				ExpectError: regexp.MustCompile(`Incompatible role and scope combination`),
			},
			{
				Config: testAccRoleAssignmentApplicationConfig_Organisation(resourceName, name, "Environment Admin", organisationID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "application_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "role_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckNoResourceAttr(resourceFullName, "scope_population_id"),
					resource.TestMatchResourceAttr(resourceFullName, "scope_organization_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckNoResourceAttr(resourceFullName, "scope_environment_id"),
					resource.TestCheckResourceAttr(resourceFullName, "read_only", "false"),
				),
			},
			{
				Config: testAccRoleAssignmentApplicationConfig_Organisation(resourceName, name, "Organization Admin", organisationID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "application_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "role_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckNoResourceAttr(resourceFullName, "scope_population_id"),
					resource.TestMatchResourceAttr(resourceFullName, "scope_organization_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckNoResourceAttr(resourceFullName, "scope_environment_id"),
					resource.TestCheckResourceAttr(resourceFullName, "read_only", "false"),
				),
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

						return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["application_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccRoleAssignmentApplication_Environment(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application_role_assignment.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckOrganisationID(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.RoleAssignmentApplication_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRoleAssignmentApplicationConfig_Environment(resourceName, name, "Identity Data Admin"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "application_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "role_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckNoResourceAttr(resourceFullName, "scope_population_id"),
					resource.TestCheckNoResourceAttr(resourceFullName, "scope_organization_id"),
					resource.TestMatchResourceAttr(resourceFullName, "scope_environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "read_only", "false"),
				),
			},
			{
				Config: testAccRoleAssignmentApplicationConfig_Environment(resourceName, name, "Environment Admin"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "application_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "role_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckNoResourceAttr(resourceFullName, "scope_population_id"),
					resource.TestCheckNoResourceAttr(resourceFullName, "scope_organization_id"),
					resource.TestMatchResourceAttr(resourceFullName, "scope_environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "read_only", "false"),
				),
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

						return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["application_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      testAccRoleAssignmentApplicationConfig_Environment(resourceName, name, "Organization Admin"),
				ExpectError: regexp.MustCompile(`Incompatible role and scope combination`),
			},
		},
	})
}

func TestAccRoleAssignmentApplication_SystemApplication(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.RoleAssignmentApplication_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccRoleAssignmentApplicationConfig_UnmappableSystemApplication(environmentName, licenseID, resourceName, "Environment Admin"),
				ExpectError: regexp.MustCompile(`Invalid parameter value - Unmappable application type`),
			},
		},
	})
}

func TestAccRoleAssignmentApplication_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application_role_assignment.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.RoleAssignmentApplication_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccRoleAssignmentApplicationConfig_Population(resourceName, name, "Identity Data Admin"),
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

func testAccRoleAssignmentApplicationConfig_NewEnv(environmentName, licenseID, resourceName, name, roleName string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_population" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  name = "%[4]s"
}

resource "pingone_application" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s"
  enabled        = true

  oidc_options {
    type                        = "WORKER"
    grant_types                 = ["CLIENT_CREDENTIALS"]
    token_endpoint_authn_method = "CLIENT_SECRET_BASIC"
  }
}

data "pingone_role" "%[3]s" {
  name = "%[5]s"
}

resource "pingone_application_role_assignment" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  application_id = pingone_application.%[3]s.id
  role_id        = data.pingone_role.%[3]s.id

  scope_population_id = pingone_population.%[3]s.id
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name, roleName)
}

func testAccRoleAssignmentApplicationConfig_Population(resourceName, name, roleName string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_population" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true

  oidc_options {
    type                        = "WORKER"
    grant_types                 = ["CLIENT_CREDENTIALS"]
    token_endpoint_authn_method = "CLIENT_SECRET_BASIC"
  }
}

data "pingone_role" "%[2]s" {
  name = "%[4]s"
}

resource "pingone_application_role_assignment" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  application_id = pingone_application.%[2]s.id
  role_id        = data.pingone_role.%[2]s.id

  scope_population_id = pingone_population.%[2]s.id
}`, acctest.GenericSandboxEnvironment(), resourceName, name, roleName)
}

func testAccRoleAssignmentApplicationConfig_Organisation(resourceName, name, roleName, organisationID string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true

  oidc_options {
    type                        = "WORKER"
    grant_types                 = ["CLIENT_CREDENTIALS"]
    token_endpoint_authn_method = "CLIENT_SECRET_BASIC"
  }
}

data "pingone_role" "%[2]s" {
  name = "%[4]s"
}

resource "pingone_application_role_assignment" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  application_id = pingone_application.%[2]s.id
  role_id        = data.pingone_role.%[2]s.id

  scope_organization_id = "%[5]s"
}`, acctest.GenericSandboxEnvironment(), resourceName, name, roleName, organisationID)
}

func testAccRoleAssignmentApplicationConfig_Environment(resourceName, name, roleName string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true

  oidc_options {
    type                        = "WORKER"
    grant_types                 = ["CLIENT_CREDENTIALS"]
    token_endpoint_authn_method = "CLIENT_SECRET_BASIC"
  }
}

data "pingone_role" "%[2]s" {
  name = "%[4]s"
}

resource "pingone_application_role_assignment" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  application_id = pingone_application.%[2]s.id
  role_id        = data.pingone_role.%[2]s.id

  scope_environment_id = data.pingone_environment.general_test.id
}`, acctest.GenericSandboxEnvironment(), resourceName, name, roleName)
}

func testAccRoleAssignmentApplicationConfig_UnmappableSystemApplication(environmentName, licenseID, resourceName, roleName string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_system_application" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  type           = "PING_ONE_PORTAL"
  enabled        = true
}

data "pingone_role" "%[3]s" {
  name = "%[4]s"
}

resource "pingone_application_role_assignment" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  application_id = pingone_system_application.%[3]s.id
  role_id        = data.pingone_role.%[3]s.id

  scope_environment_id = pingone_environment.%[2]s.id
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, roleName)
}
