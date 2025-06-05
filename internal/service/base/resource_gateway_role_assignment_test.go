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
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccRoleAssignmentGateway_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_gateway_role_assignment.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	var roleAssignmentID, gatewayID, environmentID string

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
		CheckDestroy:             base.RoleAssignmentGateway_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Test removal of the resource
			{
				Config: testAccRoleAssignmentGatewayConfig_Environment(environmentName, licenseID, resourceName, name, "Identity Data Admin"),
				Check:  base.RoleAssignmentGateway_GetIDs(resourceFullName, &environmentID, &gatewayID, &roleAssignmentID),
			},
			{
				PreConfig: func() {
					base.RoleAssignmentGateway_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID, gatewayID, roleAssignmentID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the gateway
			{
				Config: testAccRoleAssignmentGatewayConfig_Environment(environmentName, licenseID, resourceName, name, "Identity Data Admin"),
				Check:  base.RoleAssignmentGateway_GetIDs(resourceFullName, &environmentID, &gatewayID, &roleAssignmentID),
			},
			{
				PreConfig: func() {
					base.Gateway_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID, gatewayID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the environment
			{
				Config: testAccRoleAssignmentGatewayConfig_Environment(environmentName, licenseID, resourceName, name, "Identity Data Admin"),
				Check:  base.RoleAssignmentGateway_GetIDs(resourceFullName, &environmentID, &gatewayID, &roleAssignmentID),
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

func TestAccRoleAssignmentGateway_Application(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_gateway_role_assignment.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	successCheck := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "gateway_id", verify.P1ResourceIDRegexpFullString),
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
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.RoleAssignmentGateway_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRoleAssignmentGatewayConfig_Application(environmentName, licenseID, resourceName, name, "Application Owner"),
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

						return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["gateway_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      testAccRoleAssignmentGatewayConfig_Application(environmentName, licenseID, resourceName, name, "Identity Data Admin"),
				ExpectError: regexp.MustCompile(`Incompatible role and scope combination`),
			},
			{
				Config:      testAccRoleAssignmentGatewayConfig_Application(environmentName, licenseID, resourceName, name, "Environment Admin"),
				ExpectError: regexp.MustCompile(`Incompatible role and scope combination`),
			},
			{
				Config:      testAccRoleAssignmentGatewayConfig_Application(environmentName, licenseID, resourceName, name, "Organization Admin"),
				ExpectError: regexp.MustCompile(`Incompatible role and scope combination`),
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

	successCheck := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "gateway_id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "role_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckNoResourceAttr(resourceFullName, "scope_application_id"),
		resource.TestMatchResourceAttr(resourceFullName, "scope_population_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckNoResourceAttr(resourceFullName, "scope_organization_id"),
		resource.TestCheckNoResourceAttr(resourceFullName, "scope_environment_id"),
		resource.TestCheckResourceAttr(resourceFullName, "read_only", "false"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.RoleAssignmentGateway_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRoleAssignmentGatewayConfig_Population(environmentName, licenseID, resourceName, name, "Identity Data Admin"),
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

						return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["gateway_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      testAccRoleAssignmentGatewayConfig_Population(environmentName, licenseID, resourceName, name, "Application Owner"),
				ExpectError: regexp.MustCompile(`Incompatible role and scope combination`),
			},
			{
				Config:      testAccRoleAssignmentGatewayConfig_Population(environmentName, licenseID, resourceName, name, "Environment Admin"),
				ExpectError: regexp.MustCompile(`Incompatible role and scope combination`),
			},
			{
				Config:      testAccRoleAssignmentGatewayConfig_Population(environmentName, licenseID, resourceName, name, "Organization Admin"),
				ExpectError: regexp.MustCompile(`Incompatible role and scope combination`),
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

	successCheck := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "gateway_id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "role_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckNoResourceAttr(resourceFullName, "scope_application_id"),
		resource.TestCheckNoResourceAttr(resourceFullName, "scope_population_id"),
		resource.TestCheckNoResourceAttr(resourceFullName, "scope_organization_id"),
		resource.TestMatchResourceAttr(resourceFullName, "scope_environment_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "read_only", "false"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckOrganisationID(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.RoleAssignmentGateway_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRoleAssignmentGatewayConfig_Environment(environmentName, licenseID, resourceName, name, "Identity Data Admin"),
				Check:  successCheck,
			},
			{
				Config: testAccRoleAssignmentGatewayConfig_Environment(environmentName, licenseID, resourceName, name, "Environment Admin"),
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

						return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["gateway_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccRoleAssignmentGatewayConfig_Environment(environmentName, licenseID, resourceName, name, "Application Owner"),
				Check:  successCheck,
			},
			{
				Config:      testAccRoleAssignmentGatewayConfig_Environment(environmentName, licenseID, resourceName, name, "Organization Admin"),
				ExpectError: regexp.MustCompile(`Incompatible role and scope combination`),
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
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.RoleAssignmentGateway_CheckDestroy,
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

func testAccRoleAssignmentGatewayConfig_Application(environmentName, licenseID, resourceName, name, roleName string) string {
	return fmt.Sprintf(`
		%[1]s

		%[5]s

resource "pingone_gateway" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true

  type = "PING_FEDERATE"
}

resource "pingone_application" "%[2]s" {
  environment_id = pingone_environment.%[6]s.id
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

data "pingone_role" "%[2]s" {
  name = "%[4]s"
}

resource "pingone_gateway_role_assignment" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  gateway_id     = pingone_gateway.%[2]s.id
  role_id        = data.pingone_role.%[2]s.id

  scope_application_id = pingone_application.%[2]s.id
}`, acctest.GenericSandboxEnvironment(), resourceName, name, roleName, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName)
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

resource "pingone_population" "%[2]s" {
  environment_id = pingone_environment.%[6]s.id
  name           = "%[3]s"
}

data "pingone_role" "%[2]s" {
  name = "%[4]s"
}

resource "pingone_gateway_role_assignment" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  gateway_id     = pingone_gateway.%[2]s.id
  role_id        = data.pingone_role.%[2]s.id

  scope_population_id = pingone_population.%[2]s.id
}`, acctest.GenericSandboxEnvironment(), resourceName, name, roleName, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName)
}

func testAccRoleAssignmentGatewayConfig_Environment(environmentName, licenseID, resourceName, name, roleName string) string {
	return fmt.Sprintf(`
		%[1]s

		%[5]s

resource "pingone_gateway" "%[2]s" {
  environment_id = pingone_environment.%[6]s.id
  name           = "%[3]s"
  enabled        = true

  type = "PING_FEDERATE"
}

data "pingone_role" "%[2]s" {
  name = "%[4]s"
}

resource "pingone_gateway_role_assignment" "%[2]s" {
  environment_id = pingone_environment.%[6]s.id
  gateway_id     = pingone_gateway.%[2]s.id
  role_id        = data.pingone_role.%[2]s.id

  scope_environment_id = data.pingone_environment.general_test.id
}`, acctest.GenericSandboxEnvironment(), resourceName, name, roleName, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName)
}
