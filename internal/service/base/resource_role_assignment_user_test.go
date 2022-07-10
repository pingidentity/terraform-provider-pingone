package base_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	pingone "github.com/patrickcping/pingone-go/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
)

func testAccCheckRoleAssignmentUserDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API
	ctx = context.WithValue(ctx, pingone.ContextServerVariables, map[string]string{
		"suffix": p1Client.RegionSuffix,
	})

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_role_assignment_user" {
			continue
		}

		body, r, err := apiClient.UsersUserRoleAssignmentsApi.ReadOneUserRoleAssignment(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["user_id"], rs.Primary.ID).Execute()

		if r.StatusCode == 404 {
			continue
		}

		if err != nil {
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

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")
	region := os.Getenv("PINGONE_REGION")

	roleName := "Identity Data Admin"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckRoleAssignmentUserDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRoleAssignmentUserConfig_Population(environmentName, resourceName, roleName, licenseID, region),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, "id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "environment_id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "user_id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "role_id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "scope_population_id"),
					resource.TestCheckNoResourceAttr(resourceFullName, "scope_organization_id"),
					resource.TestCheckNoResourceAttr(resourceFullName, "scope_environment_id"),
					resource.TestCheckResourceAttr(resourceFullName, "read_only", "false"),
				),
			},
		},
	})
}

func TestAccRoleAssignmentUser_Organisation(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_role_assignment_user.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")
	region := os.Getenv("PINGONE_REGION")
	organisationID := os.Getenv("PINGONE_ORGANISATION_ID")

	roleName := "Identity Data Admin"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckRoleAssignmentUserDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRoleAssignmentUserConfig_Organisation(environmentName, resourceName, roleName, licenseID, region, organisationID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, "id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "environment_id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "user_id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "role_id"),
					resource.TestCheckNoResourceAttr(resourceFullName, "scope_population_id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "scope_organization_id"),
					resource.TestCheckNoResourceAttr(resourceFullName, "scope_environment_id"),
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

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")
	region := os.Getenv("PINGONE_REGION")

	roleName := "Identity Data Admin"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironmentAndOrganisation(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckRoleAssignmentUserDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRoleAssignmentUserConfig_Environment(environmentName, resourceName, roleName, licenseID, region),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, "id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "environment_id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "user_id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "role_id"),
					resource.TestCheckNoResourceAttr(resourceFullName, "scope_population_id"),
					resource.TestCheckNoResourceAttr(resourceFullName, "scope_organization_id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "scope_environment_id"),
					resource.TestCheckResourceAttr(resourceFullName, "read_only", "false"),
				),
			},
		},
	})
}

func testAccRoleAssignmentUserConfig_Population(environmentName, resourceName, roleName, licenseID, region string) string {
	return fmt.Sprintf(`
		resource "pingone_environment" "%[1]s" {
			name = "%[1]s"
			type = "SANDBOX"
			license_id = "%[4]s"
			region = "%[5]s"
			default_population {}
			service {}
		}

		resource "pingone_user" "%[2]s" {
			environment_id = "${pingone_environment.%[1]s.id}"
			population_id = "${pingone_environment.%[1]s.default_population_id}"

			username = "%[2]s"
		}

		data "pingone_role" "%[2]s" {
			name = "%[3]s"
		}

		resource "pingone_role_assignment_user" "%[2]s" {
			environment_id = "${pingone_environment.%[1]s.id}"
			user_id = "${pingone_user.%[2]s.id}"
			role_id = "${data.pingone_role.%[2]s.id}"

			scope_population_id = "${pingone_environment.%[1]s.default_population_id}"
		}`, environmentName, resourceName, roleName, licenseID, region)
}

func testAccRoleAssignmentUserConfig_Organisation(environmentName, resourceName, roleName, licenseID, region, organisationID string) string {
	return fmt.Sprintf(`
		resource "pingone_environment" "%[1]s" {
			name = "%[1]s"
			type = "SANDBOX"
			license_id = "%[4]s"
			region = "%[5]s"
			default_population {}
			service {}
		}

		resource "pingone_user" "%[2]s" {
			environment_id = "${pingone_environment.%[1]s.id}"
			population_id = "${pingone_environment.%[1]s.default_population_id}"

			username = "%[2]s"
		}

		data "pingone_role" "%[2]s" {
			name = "%[3]s"
		}

		resource "pingone_role_assignment_user" "%[2]s" {
			environment_id = "${pingone_environment.%[1]s.id}"
			user_id = "${pingone_user.%[2]s.id}"
			role_id = "${data.pingone_role.%[2]s.id}"

			scope_organization_id = "%[6]s"
		}`, environmentName, resourceName, roleName, licenseID, region, organisationID)
}

func testAccRoleAssignmentUserConfig_Environment(environmentName, resourceName, roleName, licenseID, region string) string {
	return fmt.Sprintf(`
		resource "pingone_environment" "%[1]s" {
			name = "%[1]s"
			type = "SANDBOX"
			license_id = "%[4]s"
			region = "%[5]s"
			default_population {}
			service {}
		}

		resource "pingone_user" "%[2]s" {
			environment_id = "${pingone_environment.%[1]s.id}"
			population_id = "${pingone_environment.%[1]s.default_population_id}"

			username = "%[2]s"
		}

		data "pingone_role" "%[2]s" {
			name = "%[3]s"
		}

		resource "pingone_role_assignment_user" "%[2]s" {
			environment_id = "${pingone_environment.%[1]s.id}"
			user_id = "${pingone_user.%[2]s.id}"
			role_id = "${data.pingone_role.%[2]s.id}"

			scope_environment_id = "${pingone_environment.%[1]s.id}"
		}`, environmentName, resourceName, roleName, licenseID, region)
}
