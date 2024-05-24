package authorize_test

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
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccApplicationRolePermission_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_application_role_permission.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	var applicationRolePermissionID, applicationRoleID, environmentID string

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
		CheckDestroy:             authorize.ApplicationRolePermission_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccApplicationRolePermissionConfig_Minimal(resourceName, name),
				Check:  authorize.ApplicationRolePermission_GetIDs(resourceFullName, &environmentID, &applicationRoleID, &applicationRolePermissionID),
			},
			{
				PreConfig: func() {
					authorize.ApplicationRolePermission_RemovalDrift_PreConfig(ctx, p1Client.API.AuthorizeAPIClient, t, environmentID, applicationRoleID, applicationRolePermissionID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the application role
			{
				Config: testAccApplicationRolePermissionConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check:  authorize.ApplicationRolePermission_GetIDs(resourceFullName, &environmentID, &applicationRoleID, &applicationRolePermissionID),
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
				Config: testAccApplicationRolePermissionConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check:  authorize.ApplicationRolePermission_GetIDs(resourceFullName, &environmentID, &applicationRoleID, &applicationRolePermissionID),
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

func TestAccApplicationRolePermission_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_application_role_permission.%s", resourceName)

	name := resourceName

	fullCheck := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "name", name),
		resource.TestCheckResourceAttr(resourceFullName, "description", "Test application role"),
	)

	minimalCheck := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "name", name),
		resource.TestCheckNoResourceAttr(resourceFullName, "description"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             authorize.ApplicationRolePermission_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			{
				Config: testAccApplicationRolePermissionConfig_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config:  testAccApplicationRolePermissionConfig_Full(resourceName, name),
				Destroy: true,
			},
			// Minimal
			{
				Config: testAccApplicationRolePermissionConfig_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config:  testAccApplicationRolePermissionConfig_Minimal(resourceName, name),
				Destroy: true,
			},
			// Change
			{
				Config: testAccApplicationRolePermissionConfig_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config: testAccApplicationRolePermissionConfig_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config: testAccApplicationRolePermissionConfig_Full(resourceName, name),
				Check:  fullCheck,
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

						return fmt.Sprintf("%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccApplicationRolePermission_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_application_role_permission.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             authorize.ApplicationRolePermission_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccApplicationRolePermissionConfig_Minimal(resourceName, name),
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
				ImportStateId: "badformat/badformat",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`Unexpected Import Identifier`),
			},
		},
	})
}

func testAccApplicationRolePermissionConfig_Full(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_authorize_application_role_permission" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Test application role"
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccApplicationRolePermissionConfig_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_authorize_application_role_permission" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
