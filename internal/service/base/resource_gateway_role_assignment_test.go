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
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/base"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccCheckRoleAssignmentGatewayDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_gateway_role_assignment" {
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

		body, r, err := apiClient.GatewayRoleAssignmentsApi.ReadOneGatewayRoleAssignment(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["gateway_id"], rs.Primary.ID).Execute()

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

		return fmt.Errorf("PingOne Gateway Role Assignment %s still exists", rs.Primary.ID)
	}

	return nil
}

func TestAccGetRoleAssignmentGatewayIDs(resourceName string, environmentID, gatewayID, resourceID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource Not found: %s", resourceName)
		}

		*resourceID = rs.Primary.ID
		*gatewayID = rs.Primary.Attributes["gateway_id"]
		*environmentID = rs.Primary.Attributes["environment_id"]

		return nil
	}
}

func TestAccRoleAssignmentGateway_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_gateway_role_assignment.%s", resourceName)

	name := resourceName

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	var resourceID, gatewayID, environmentID string

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.TestAccCheckRoleAssignmentGatewayDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Test removal of the resource
			{
				Config: testAccRoleAssignmentGatewayConfig_Environment(environmentName, licenseID, resourceName, name, "Identity Data Admin"),
				Check:  base.TestAccGetRoleAssignmentGatewayIDs(resourceFullName, &environmentID, &gatewayID, &resourceID),
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

					if environmentID == "" || gatewayID == "" || resourceID == "" {
						t.Fatalf("One of environment ID or resource ID cannot be determined. Environment ID: %s, Gateway ID: %s, Resource ID: %s", environmentID, gatewayID, resourceID)
					}

					_, err = apiClient.GatewayRoleAssignmentsApi.DeleteGatewayRoleAssignment(ctx, environmentID, gatewayID, resourceID).Execute()
					if err != nil {
						t.Fatalf("Failed to delete Gateway role assignment: %v", err)
					}
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the gateway
			{
				Config: testAccRoleAssignmentGatewayConfig_Environment(environmentName, licenseID, resourceName, name, "Identity Data Admin"),
				Check:  base.TestAccGetRoleAssignmentGatewayIDs(resourceFullName, &environmentID, &gatewayID, &resourceID),
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

					if environmentID == "" || gatewayID == "" || resourceID == "" {
						t.Fatalf("One of environment ID or resource ID cannot be determined. Environment ID: %s, Gateway ID: %s, Resource ID: %s", environmentID, gatewayID, resourceID)
					}

					_, err = apiClient.GatewaysApi.DeleteGateway(ctx, environmentID, gatewayID).Execute()
					if err != nil {
						t.Fatalf("Failed to delete Gateway: %v", err)
					}
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccRoleAssignmentGateway_Population(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_gateway_role_assignment.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.TestAccCheckRoleAssignmentGatewayDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRoleAssignmentGatewayConfig_Population(environmentName, licenseID, resourceName, name, "Identity Data Admin"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "gateway_id", verify.P1ResourceIDRegexpFullString),
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

						return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["gateway_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      testAccRoleAssignmentGatewayConfig_Population(environmentName, licenseID, resourceName, name, "Environment Admin"),
				ExpectError: regexp.MustCompile(`Incompatible role and scope combination. Role: [a-z0-9\-]* \/ Scope: POPULATION`),
			},
			{
				Config:      testAccRoleAssignmentGatewayConfig_Population(environmentName, licenseID, resourceName, name, "Organization Admin"),
				ExpectError: regexp.MustCompile(`Incompatible role and scope combination. Role: [a-z0-9\-]* \/ Scope: POPULATION`),
			},
		},
	})
}

func TestAccRoleAssignmentGateway_Environment(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_gateway_role_assignment.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckOrganisationID(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.TestAccCheckRoleAssignmentGatewayDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRoleAssignmentGatewayConfig_Environment(environmentName, licenseID, resourceName, name, "Identity Data Admin"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "gateway_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "role_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "scope_population_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "scope_organization_id", ""),
					resource.TestMatchResourceAttr(resourceFullName, "scope_environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "read_only", "false"),
				),
			},
			{
				Config: testAccRoleAssignmentGatewayConfig_Environment(environmentName, licenseID, resourceName, name, "Environment Admin"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "gateway_id", verify.P1ResourceIDRegexpFullString),
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

						return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["gateway_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      testAccRoleAssignmentGatewayConfig_Environment(environmentName, licenseID, resourceName, name, "Organization Admin"),
				ExpectError: regexp.MustCompile(`Incompatible role and scope combination. Role: [a-z0-9\-]* \/ Scope: ENVIRONMENT`),
			},
		},
	})
}

func TestAccRoleAssignmentGateway_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_gateway_role_assignment.%s", resourceName)

	name := resourceName

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.TestAccCheckRoleAssignmentGatewayDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccRoleAssignmentGatewayConfig_Environment(environmentName, licenseID, resourceName, name, "Identity Data Admin"),
			},
			// Errors
			{
				ResourceName: resourceFullName,
				ImportState:  true,
				ExpectError:  regexp.MustCompile(`Invalid import ID specified \(".*"\).  The ID should be in the format "environment_id/gateway_id/role_assignment_id" and must match regex: .*`),
			},
			{
				ResourceName:  resourceFullName,
				ImportStateId: "/",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`Invalid import ID specified \(".*"\).  The ID should be in the format "environment_id/gateway_id/role_assignment_id" and must match regex: .*`),
			},
			{
				ResourceName:  resourceFullName,
				ImportStateId: "badformat/badformat/badformat",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`Invalid import ID specified \(".*"\).  The ID should be in the format "environment_id/gateway_id/role_assignment_id" and must match regex: .*`),
			},
		},
	})
}

func testAccRoleAssignmentGatewayConfig_Population(environmentName, licenseID, resourceName, name, roleName string) string {
	return fmt.Sprintf(`
		%[1]s

		%[5]s

resource "pingone_gateway" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true

  type = "PING_FEDERATE"
}

data "pingone_role" "%[2]s" {
  name = "%[4]s"
}

resource "pingone_gateway_role_assignment" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  gateway_id     = pingone_gateway.%[2]s.id
  role_id        = data.pingone_role.%[2]s.id

  scope_population_id = pingone_environment.%[6]s.default_population_id
}`, acctest.GenericSandboxEnvironment(), resourceName, name, roleName, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName)
}

func testAccRoleAssignmentGatewayConfig_Environment(environmentName, licenseID, resourceName, name, roleName string) string {
	return fmt.Sprintf(`
		%[1]s

		%[5]s

resource "pingone_gateway" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true

  type = "PING_FEDERATE"
}

data "pingone_role" "%[2]s" {
  name = "%[4]s"
}

resource "pingone_gateway_role_assignment" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  gateway_id     = pingone_gateway.%[2]s.id
  role_id        = data.pingone_role.%[2]s.id

  scope_environment_id = pingone_environment.%[6]s.id
}`, acctest.GenericSandboxEnvironment(), resourceName, name, roleName, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName)
}
