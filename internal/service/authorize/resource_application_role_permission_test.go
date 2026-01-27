// Copyright © 2025 Ping Identity Corporation

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
	acctestlegacysdk "github.com/pingidentity/terraform-provider-pingone/internal/acctest/legacysdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/authorize"
	baselegacysdk "github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/base/legacysdk"
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
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoBeta(t)
			p1Client = acctestlegacysdk.PreCheckTestClient(ctx, t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             authorize.ApplicationRolePermission_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccApplicationRolePermissionConfig_Full(resourceName, name),
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
				Config: testAccApplicationRolePermissionConfig_Full(resourceName, name),
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
					baselegacysdk.Environment_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID)
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
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "application_role_id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "application_resource_permission_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "action", name),
		resource.TestCheckResourceAttr(resourceFullName, "description", "Test permission"),
		resource.TestCheckResourceAttr(resourceFullName, "key", fmt.Sprintf("%s:%s", name, name)),
		resource.TestMatchResourceAttr(resourceFullName, "resource.id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "resource.name", name),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
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
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("resource not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["application_role_id"], rs.Primary.Attributes["application_resource_permission_id"]), nil
					}
				}(),
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "application_resource_permission_id",
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
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             authorize.ApplicationRolePermission_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccApplicationRolePermissionConfig_Full(resourceName, name),
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

func testAccApplicationRolePermissionConfig_NewEnv(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_resource" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  name = "%[4]s"
}

resource "pingone_application_resource" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  resource_name  = pingone_resource.%[3]s.name

  name = "%[4]s"
}

resource "pingone_application_resource_permission" "%[3]s" {
  environment_id          = pingone_environment.%[2]s.id
  application_resource_id = pingone_application_resource.%[3]s.id

  action = "%[4]s"
}


resource "pingone_authorize_application_role" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  name = "%[4]s"
}

resource "pingone_authorize_application_role_permission" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  application_role_id                = pingone_authorize_application_role.%[3]s.id
  application_resource_permission_id = pingone_application_resource_permission.%[3]s.id
}`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccApplicationRolePermissionConfig_Full(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_resource" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_application_resource" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  resource_name  = pingone_resource.%[2]s.name

  name = "%[3]s"
}

resource "pingone_application_resource_permission" "%[2]s" {
  environment_id          = data.pingone_environment.general_test.id
  application_resource_id = pingone_application_resource.%[2]s.id

  action      = "%[3]s"
  description = "Test permission"
}


resource "pingone_authorize_application_role" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_authorize_application_role_permission" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  application_role_id                = pingone_authorize_application_role.%[2]s.id
  application_resource_permission_id = pingone_application_resource_permission.%[2]s.id
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
