package sso_test

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

func testAccCheckRoleAssignmentApplicationDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_application_role_assignment" {
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

		body, r, err := apiClient.ApplicationRoleAssignmentsApi.ReadOneApplicationRoleAssignment(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["application_id"], rs.Primary.ID).Execute()

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

		return fmt.Errorf("PingOne Application Role Assignment %s still exists", rs.Primary.ID)
	}

	return nil
}

func testAccGetRoleAssignmentApplicationIDs(resourceName string, environmentID, applicationID, resourceID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource Not found: %s", resourceName)
		}

		*resourceID = rs.Primary.ID
		*applicationID = rs.Primary.Attributes["application_id"]
		*environmentID = rs.Primary.Attributes["environment_id"]

		return nil
	}
}

func TestAccRoleAssignmentApplication_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application_role_assignment.%s", resourceName)

	name := resourceName

	var resourceID, applicationID, environmentID string

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckRoleAssignmentApplicationDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccRoleAssignmentApplicationConfig_Population(resourceName, name, "Identity Data Admin"),
				Check:  testAccGetRoleAssignmentApplicationIDs(resourceFullName, &environmentID, &applicationID, &resourceID),
			},
			// Replan after removal preconfig
			{
				PreConfig: func() {
					var ctx = context.Background()
					p1Client, err := acctest.TestClient(ctx)

					if err != nil {
						t.Fatalf("Failed to get API client: %v", err)
					}

					apiClient := p1Client.API.ManagementAPIClient

					if environmentID == "" || applicationID == "" || resourceID == "" {
						t.Fatalf("One of environment ID, applcation ID or resource ID cannot be determined. Environment ID: %s, Application ID: %s, Resource ID: %s", environmentID, applicationID, resourceID)
					}

					_, err = apiClient.ApplicationRoleAssignmentsApi.DeleteApplicationRoleAssignment(ctx, environmentID, applicationID, resourceID).Execute()
					if err != nil {
						t.Fatalf("Failed to delete Application role assignment: %v", err)
					}
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
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckRoleAssignmentApplicationDestroy,
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
					resource.TestCheckResourceAttr(resourceFullName, "scope_organization_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "scope_environment_id", ""),
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
				ExpectError: regexp.MustCompile(`Incompatible role and scope combination. Role: [a-z0-9\-]* \/ Scope: POPULATION`),
			},
			{
				Config:      testAccRoleAssignmentApplicationConfig_Population(resourceName, name, "Organization Admin"),
				ExpectError: regexp.MustCompile(`Incompatible role and scope combination. Role: [a-z0-9\-]* \/ Scope: POPULATION`),
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
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckRoleAssignmentApplicationDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccRoleAssignmentApplicationConfig_Organisation(resourceName, name, "Identity Data Admin", organisationID),
				ExpectError: regexp.MustCompile(`Incompatible role and scope combination. Role: [a-z0-9\-]* \/ Scope: ORGANIZATION`),
			},
			{
				Config: testAccRoleAssignmentApplicationConfig_Organisation(resourceName, name, "Environment Admin", organisationID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "application_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "role_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "scope_population_id", ""),
					resource.TestMatchResourceAttr(resourceFullName, "scope_organization_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "scope_environment_id", ""),
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
					resource.TestCheckResourceAttr(resourceFullName, "scope_population_id", ""),
					resource.TestMatchResourceAttr(resourceFullName, "scope_organization_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "scope_environment_id", ""),
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
		PreCheck:                 func() { acctest.PreCheckEnvironmentAndOrganisation(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckRoleAssignmentApplicationDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRoleAssignmentApplicationConfig_Environment(resourceName, name, "Identity Data Admin"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "application_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "role_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "scope_population_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "scope_organization_id", ""),
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
					resource.TestCheckResourceAttr(resourceFullName, "scope_population_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "scope_organization_id", ""),
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
				ExpectError: regexp.MustCompile(`Incompatible role and scope combination. Role: [a-z0-9\-]* \/ Scope: ENVIRONMENT`),
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
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckRoleAssignmentApplicationDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccRoleAssignmentApplicationConfig_SystemApplication(environmentName, licenseID, resourceName, "Environment Admin"),
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
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckRoleAssignmentApplicationDestroy,
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
				ExpectError:  regexp.MustCompile(`Invalid import ID specified \(".*"\).  The ID should be in the format "environment_id/application_id/role_assignment_id" and must match regex: .*`),
			},
			{
				ResourceName:  resourceFullName,
				ImportStateId: "/",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`Invalid import ID specified \(".*"\).  The ID should be in the format "environment_id/application_id/role_assignment_id" and must match regex: .*`),
			},
			{
				ResourceName:  resourceFullName,
				ImportStateId: "badformat/badformat/badformat",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`Invalid import ID specified \(".*"\).  The ID should be in the format "environment_id/application_id/role_assignment_id" and must match regex: .*`),
			},
		},
	})
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

func testAccRoleAssignmentApplicationConfig_SystemApplication(environmentName, licenseID, resourceName, roleName string) string {
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
