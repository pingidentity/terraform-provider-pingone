package base_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/base"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccSystemApplication_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_system_application.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	var applicationID, environmentID string

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
		CheckDestroy:             base.SystemApplication_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Test removal of the environment
			{
				Config: testAccSystemApplicationConfig_NewEnv(environmentName, licenseID, resourceName),
				Check:  base.SystemApplication_GetIDs(resourceFullName, &environmentID, &applicationID),
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

func TestAccSystemApplication_NewEnv(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_system_application.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.SystemApplication_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSystemApplicationConfig_NewEnv(environmentName, licenseID, resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
				),
			},
		},
	})
}

func TestAccSystemApplication_Portal(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_system_application.%s", resourceName)

	applicationType := string(management.ENUMAPPLICATIONTYPE_PING_ONE_PORTAL)

	fullCheck := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "type", applicationType),
		resource.TestCheckResourceAttr(resourceFullName, "name", "PingOne Application Portal"),
		resource.TestCheckResourceAttr(resourceFullName, "enabled", "false"),
		resource.TestCheckResourceAttr(resourceFullName, "access_control_role_type", "ADMIN_USERS_ONLY"),
		resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.groups.#", "2"),
		resource.TestMatchResourceAttr(resourceFullName, "access_control_group_options.groups.0", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "access_control_group_options.groups.1", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.type", "ALL_GROUPS"),
		resource.TestCheckResourceAttr(resourceFullName, "apply_default_theme", "true"),
		resource.TestCheckNoResourceAttr(resourceFullName, "enable_default_theme_footer"),
	)

	minimalCheck := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "type", applicationType),
		resource.TestCheckResourceAttr(resourceFullName, "name", "PingOne Application Portal"),
		resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
		resource.TestCheckNoResourceAttr(resourceFullName, "access_control_role_type"),
		resource.TestCheckResourceAttr(resourceFullName, "apply_default_theme", "false"),
		resource.TestCheckNoResourceAttr(resourceFullName, "enable_default_theme_footer"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.SystemApplication_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			{
				Config: testAccSystemApplicationConfig_Portal_Full(resourceName, false),
				Check:  fullCheck,
			},
			{
				Config:  testAccSystemApplicationConfig_Portal_Full(resourceName, false),
				Destroy: true,
			},
			// Minimal
			{
				Config: testAccSystemApplicationConfig_Minimal(resourceName, applicationType),
				Check:  minimalCheck,
			},
			{
				Config:  testAccSystemApplicationConfig_Minimal(resourceName, applicationType),
				Destroy: true,
			},
			// Change
			{
				Config: testAccSystemApplicationConfig_Portal_Full(resourceName, false),
				Check:  fullCheck,
			},
			{
				Config: testAccSystemApplicationConfig_Minimal(resourceName, applicationType),
				Check:  minimalCheck,
			},
			{
				Config: testAccSystemApplicationConfig_Portal_Full(resourceName, false),
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

func TestAccSystemApplication_SelfService(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_system_application.%s", resourceName)

	applicationType := string(management.ENUMAPPLICATIONTYPE_PING_ONE_SELF_SERVICE)

	fullCheck := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "type", applicationType),
		resource.TestCheckResourceAttr(resourceFullName, "name", "PingOne Self-Service - MyAccount"),
		resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
		resource.TestCheckResourceAttr(resourceFullName, "access_control_role_type", "ADMIN_USERS_ONLY"),
		resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.groups.#", "2"),
		resource.TestMatchResourceAttr(resourceFullName, "access_control_group_options.groups.0", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "access_control_group_options.groups.1", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.type", "ALL_GROUPS"),
		resource.TestCheckResourceAttr(resourceFullName, "apply_default_theme", "true"),
		resource.TestCheckResourceAttr(resourceFullName, "enable_default_theme_footer", "true"),
	)

	minimalCheck := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "type", applicationType),
		resource.TestCheckResourceAttr(resourceFullName, "name", "PingOne Self-Service - MyAccount"),
		resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
		resource.TestCheckNoResourceAttr(resourceFullName, "access_control_role_type"),
		resource.TestCheckResourceAttr(resourceFullName, "apply_default_theme", "false"),
		resource.TestCheckNoResourceAttr(resourceFullName, "enable_default_theme_footer"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.SystemApplication_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			{
				Config: testAccSystemApplicationConfig_SelfService_Full(resourceName, true),
				Check:  fullCheck,
			},
			{
				Config:  testAccSystemApplicationConfig_SelfService_Full(resourceName, true),
				Destroy: true,
			},
			// Minimal
			{
				Config: testAccSystemApplicationConfig_Minimal(resourceName, applicationType),
				Check:  minimalCheck,
			},
			{
				Config:  testAccSystemApplicationConfig_Minimal(resourceName, applicationType),
				Destroy: true,
			},
			// Change
			{
				Config: testAccSystemApplicationConfig_SelfService_Full(resourceName, true),
				Check:  fullCheck,
			},
			{
				Config: testAccSystemApplicationConfig_Minimal(resourceName, applicationType),
				Check:  minimalCheck,
			},
			{
				Config: testAccSystemApplicationConfig_SelfService_Full(resourceName, true),
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

func TestAccSystemApplication_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_system_application.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.SystemApplication_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccSystemApplicationConfig_NewEnv(environmentName, licenseID, resourceName),
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

func testAccSystemApplicationConfig_NewEnv(environmentName, licenseID, resourceName string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_system_application" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  type    = "PING_ONE_PORTAL"
  enabled = false
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName)
}

func testAccSystemApplicationConfig_Portal_Full(resourceName string, enabled bool) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_group" "%[2]s-1" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[2]s-1"
}

resource "pingone_group" "%[2]s-2" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[2]s-2"
}

resource "pingone_system_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  type    = "PING_ONE_PORTAL"
  enabled = %[3]t

  access_control_role_type = "ADMIN_USERS_ONLY"
  access_control_group_options = {
    groups = [
      pingone_group.%[2]s-2.id,
      pingone_group.%[2]s-1.id,
    ]

    type = "ALL_GROUPS"
  }

  apply_default_theme = true

}`, acctest.GenericSandboxEnvironment(), resourceName, enabled)
}

func testAccSystemApplicationConfig_SelfService_Full(resourceName string, enabled bool) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_group" "%[2]s-1" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[2]s-1"
}

resource "pingone_group" "%[2]s-2" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[2]s-2"
}

resource "pingone_system_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  type    = "PING_ONE_SELF_SERVICE"
  enabled = %[3]t

  access_control_role_type = "ADMIN_USERS_ONLY"
  access_control_group_options = {
    groups = [
      pingone_group.%[2]s-2.id,
      pingone_group.%[2]s-1.id,
    ]

    type = "ALL_GROUPS"
  }

  apply_default_theme         = true
  enable_default_theme_footer = true

}`, acctest.GenericSandboxEnvironment(), resourceName, enabled)
}

func testAccSystemApplicationConfig_Minimal(resourceName, applicationType string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_system_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  type    = "%[3]s"
  enabled = true
}`, acctest.GenericSandboxEnvironment(), resourceName, applicationType)
}
