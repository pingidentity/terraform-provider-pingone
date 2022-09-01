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
	pingone "github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
)

func testAccCheckRoleAssignmentGatewayDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, pingone.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_gateway_role_assignment" {
			continue
		}

		_, rEnv, err := apiClient.EnvironmentsApi.ReadOneEnvironment(ctx, rs.Primary.Attributes["environment_id"]).Execute()

		if rEnv.StatusCode == 404 {
			continue
		}

		if err != nil {
			return err
		}

		body, r, err := apiClient.GatewayRoleAssignmentsApi.ReadOneGatewayRoleAssignment(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["gateway_id"], rs.Primary.ID).Execute()

		if r.StatusCode == 404 {
			continue
		}

		if err != nil {
			tflog.Error(ctx, fmt.Sprintf("Error: %v", body))
			return err
		}

		return fmt.Errorf("PingOne Gateway Role Assignment %s still exists", rs.Primary.ID)
	}

	return nil
}

func TestAccRoleAssignmentGateway_Population(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_gateway_role_assignment.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")
	region := os.Getenv("PINGONE_REGION")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckRoleAssignmentGatewayDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRoleAssignmentGatewayConfig_Population(environmentName, resourceName, "Identity Data Admin", licenseID, region),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, "id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "environment_id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "gateway_id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "role_id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "scope_population_id"),
					resource.TestCheckResourceAttr(resourceFullName, "scope_organization_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "scope_environment_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "read_only", "false"),
				),
			},
			{
				Config:      testAccRoleAssignmentGatewayConfig_Population(environmentName, resourceName, "Environment Admin", licenseID, region),
				ExpectError: regexp.MustCompile(`Incompatible role and scope combination. Role: [a-z0-9\-]* \/ Scope: POPULATION`),
			},
			{
				Config:      testAccRoleAssignmentGatewayConfig_Population(environmentName, resourceName, "Organization Admin", licenseID, region),
				ExpectError: regexp.MustCompile(`Incompatible role and scope combination. Role: [a-z0-9\-]* \/ Scope: POPULATION`),
			},
		},
	})
}

func TestAccRoleAssignmentGateway_Organisation(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_gateway_role_assignment.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")
	region := os.Getenv("PINGONE_REGION")
	organisationID := os.Getenv("PINGONE_ORGANIZATION_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckRoleAssignmentGatewayDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccRoleAssignmentGatewayConfig_Organisation(environmentName, resourceName, "Identity Data Admin", licenseID, region, organisationID),
				ExpectError: regexp.MustCompile(`Incompatible role and scope combination. Role: [a-z0-9\-]* \/ Scope: ORGANIZATION`),
			},
			{
				Config: testAccRoleAssignmentGatewayConfig_Organisation(environmentName, resourceName, "Environment Admin", licenseID, region, organisationID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, "id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "environment_id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "gateway_id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "role_id"),
					resource.TestCheckResourceAttr(resourceFullName, "scope_population_id", ""),
					resource.TestCheckResourceAttrSet(resourceFullName, "scope_organization_id"),
					resource.TestCheckResourceAttr(resourceFullName, "scope_environment_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "read_only", "false"),
				),
			},
			{
				Config: testAccRoleAssignmentGatewayConfig_Organisation(environmentName, resourceName, "Organization Admin", licenseID, region, organisationID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, "id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "environment_id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "gateway_id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "role_id"),
					resource.TestCheckResourceAttr(resourceFullName, "scope_population_id", ""),
					resource.TestCheckResourceAttrSet(resourceFullName, "scope_organization_id"),
					resource.TestCheckResourceAttr(resourceFullName, "scope_environment_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "read_only", "false"),
				),
			},
		},
	})
}

func TestAccRoleAssignmentGateway_Environment(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_gateway_role_assignment.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")
	region := os.Getenv("PINGONE_REGION")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironmentAndOrganisation(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckRoleAssignmentGatewayDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRoleAssignmentGatewayConfig_Environment(environmentName, resourceName, "Identity Data Admin", licenseID, region),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, "id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "environment_id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "gateway_id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "role_id"),
					resource.TestCheckResourceAttr(resourceFullName, "scope_population_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "scope_organization_id", ""),
					resource.TestCheckResourceAttrSet(resourceFullName, "scope_environment_id"),
					resource.TestCheckResourceAttr(resourceFullName, "read_only", "false"),
				),
			},
			{
				Config: testAccRoleAssignmentGatewayConfig_Environment(environmentName, resourceName, "Environment Admin", licenseID, region),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, "id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "environment_id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "gateway_id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "role_id"),
					resource.TestCheckResourceAttr(resourceFullName, "scope_population_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "scope_organization_id", ""),
					resource.TestCheckResourceAttrSet(resourceFullName, "scope_environment_id"),
					resource.TestCheckResourceAttr(resourceFullName, "read_only", "false"),
				),
			},
			{
				Config:      testAccRoleAssignmentGatewayConfig_Environment(environmentName, resourceName, "Organization Admin", licenseID, region),
				ExpectError: regexp.MustCompile(`Incompatible role and scope combination. Role: [a-z0-9\-]* \/ Scope: ENVIRONMENT`),
			},
		},
	})
}

func testAccRoleAssignmentGatewayConfig_Population(environmentName, resourceName, roleName, licenseID, region string) string {
	return fmt.Sprintf(`
		resource "pingone_environment" "%[1]s" {
			name = "%[1]s"
			type = "SANDBOX"
			license_id = "%[4]s"
			region = "%[5]s"
			default_population {}
			service {}
		}

		resource "pingone_gateway" "%[2]s" {
			environment_id  = "${pingone_environment.%[1]s.id}"
			name 			= "%[2]s"
			enabled 		= true
		  
			oidc_options {
				type                        = "WORKER"
				grant_types                 = ["CLIENT_CREDENTIALS"]
				token_endpoint_authn_method = "CLIENT_SECRET_BASIC"
			}
		}

		data "pingone_role" "%[2]s" {
			name = "%[3]s"
		}

		resource "pingone_gateway_role_assignment" "%[2]s" {
			environment_id = "${pingone_environment.%[1]s.id}"
			gateway_id = "${pingone_gateway.%[2]s.id}"
			role_id = "${data.pingone_role.%[2]s.id}"

			scope_population_id = "${pingone_environment.%[1]s.default_population_id}"
		}`, environmentName, resourceName, roleName, licenseID, region)
}

func testAccRoleAssignmentGatewayConfig_Organisation(environmentName, resourceName, roleName, licenseID, region, organisationID string) string {
	return fmt.Sprintf(`
		resource "pingone_environment" "%[1]s" {
			name = "%[1]s"
			type = "SANDBOX"
			license_id = "%[4]s"
			region = "%[5]s"
			default_population {}
			service {}
		}

		resource "pingone_gateway" "%[2]s" {
			environment_id  = "${pingone_environment.%[1]s.id}"
			name 			= "%[2]s"
			enabled 		= true
		  
			oidc_options {
				type                        = "WORKER"
				grant_types                 = ["CLIENT_CREDENTIALS"]
				token_endpoint_authn_method = "CLIENT_SECRET_BASIC"
			}
		}

		data "pingone_role" "%[2]s" {
			name = "%[3]s"
		}

		resource "pingone_gateway_role_assignment" "%[2]s" {
			environment_id = "${pingone_environment.%[1]s.id}"
			gateway_id = "${pingone_gateway.%[2]s.id}"
			role_id = "${data.pingone_role.%[2]s.id}"

			scope_organization_id = "%[6]s"
		}`, environmentName, resourceName, roleName, licenseID, region, organisationID)
}

func testAccRoleAssignmentGatewayConfig_Environment(environmentName, resourceName, roleName, licenseID, region string) string {
	return fmt.Sprintf(`
		resource "pingone_environment" "%[1]s" {
			name = "%[1]s"
			type = "SANDBOX"
			license_id = "%[4]s"
			region = "%[5]s"
			default_population {}
			service {}
		}

		resource "pingone_gateway" "%[2]s" {
			environment_id  = "${pingone_environment.%[1]s.id}"
			name 			= "%[2]s"
			enabled 		= true
		  
			oidc_options {
				type                        = "WORKER"
				grant_types                 = ["CLIENT_CREDENTIALS"]
				token_endpoint_authn_method = "CLIENT_SECRET_BASIC"
			}
		}

		data "pingone_role" "%[2]s" {
			name = "%[3]s"
		}

		resource "pingone_gateway_role_assignment" "%[2]s" {
			environment_id = "${pingone_environment.%[1]s.id}"
			gateway_id = "${pingone_gateway.%[2]s.id}"
			role_id = "${data.pingone_role.%[2]s.id}"

			scope_environment_id = "${pingone_environment.%[1]s.id}"
		}`, environmentName, resourceName, roleName, licenseID, region)
}
