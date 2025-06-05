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
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/authorize"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/base"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccApplicationResourcePermission_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application_resource_permission.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	var resourcePermissionID, applicationResourceID, environmentID string

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
		CheckDestroy:             authorize.ApplicationResourcePermission_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Test removal of the permission resource
			{
				Config: testAccApplicationResourcePermissionConfig_Custom_Full(resourceName, name),
				Check:  authorize.ApplicationResourcePermission_GetIDs(resourceFullName, &environmentID, &applicationResourceID, &resourcePermissionID),
			},
			{
				PreConfig: func() {
					authorize.ApplicationResourcePermission_RemovalDrift_PreConfig(ctx, p1Client.API.AuthorizeAPIClient, t, environmentID, applicationResourceID, resourcePermissionID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the application resource
			{
				Config: testAccApplicationResourcePermissionConfig_Custom_Full(resourceName, name),
				Check:  authorize.ApplicationResourcePermission_GetIDs(resourceFullName, &environmentID, &applicationResourceID, &resourcePermissionID),
			},
			{
				PreConfig: func() {
					authorize.ApplicationResource_RemovalDrift_PreConfig(ctx, p1Client.API, t, environmentID, applicationResourceID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the resource
			{
				Config: testAccApplicationResourcePermissionConfig_Custom_Full(resourceName, name),
				Check:  authorize.ApplicationResourcePermission_GetIDs(resourceFullName, &environmentID, &applicationResourceID, &resourcePermissionID),
			},
			{
				PreConfig: func() {
					authorize.ApplicationResource_Resource_RemovalDrift_PreConfig(ctx, p1Client.API, t, environmentID, applicationResourceID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the environment
			{
				Config: testAccApplicationResourcePermissionConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check:  authorize.ApplicationResourcePermission_GetIDs(resourceFullName, &environmentID, &applicationResourceID, &resourcePermissionID),
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

func TestAccApplicationResourcePermission_NewEnv(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application_resource_permission.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             authorize.ApplicationResourcePermission_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationResourcePermissionConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
				),
			},
		},
	})
}

func TestAccApplicationResourcePermission_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application_resource_permission.%s", resourceName)

	name := resourceName

	fullStep := resource.TestStep{
		Config: testAccApplicationResourcePermissionConfig_Custom_Full(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "application_resource_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "action", name),
			resource.TestCheckResourceAttr(resourceFullName, "description", "My custom application resource permission"),
			resource.TestMatchResourceAttr(resourceFullName, "resource.id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "resource.name", name),
		),
	}

	minimalStep := resource.TestStep{
		Config: testAccApplicationResourcePermissionConfig_Custom_Minimal(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "application_resource_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "action", name),
			resource.TestCheckNoResourceAttr(resourceFullName, "description"),
			resource.TestMatchResourceAttr(resourceFullName, "resource.id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "resource.name", name),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             authorize.ApplicationResourcePermission_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			fullStep,
			{
				Config:  testAccApplicationResourcePermissionConfig_Custom_Full(resourceName, name),
				Destroy: true,
			},
			// Minimal
			minimalStep,
			{
				Config:  testAccApplicationResourcePermissionConfig_Custom_Minimal(resourceName, name),
				Destroy: true,
			},
			// Change
			fullStep,
			minimalStep,
			fullStep,
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["application_resource_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:  testAccApplicationResourcePermissionConfig_Custom_Full(resourceName, name),
				Destroy: true,
			},
		},
	})
}

func TestAccApplicationResourcePermission_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application_resource_permission.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             authorize.ApplicationResourcePermission_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccApplicationResourcePermissionConfig_Custom_Minimal(resourceName, name),
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

func testAccApplicationResourcePermissionConfig_NewEnv(environmentName, licenseID, resourceName, name string) string {
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
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccApplicationResourcePermissionConfig_Custom_Full(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_resource" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_application_resource" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  resource_name  = pingone_resource.%[2]s.name

  name        = "%[3]s"
  description = "My custom application resource"
}

resource "pingone_application_resource_permission" "%[2]s" {
  environment_id          = data.pingone_environment.general_test.id
  application_resource_id = pingone_application_resource.%[2]s.id

  action      = "%[3]s"
  description = "My custom application resource permission"
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccApplicationResourcePermissionConfig_Custom_Minimal(resourceName, name string) string {
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

  action = "%[3]s"
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
