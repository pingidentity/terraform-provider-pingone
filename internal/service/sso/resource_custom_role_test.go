// Copyright © 2025 Ping Identity Corporation

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

func TestAccCustomRole_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_custom_role.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	var customRoleID, environmentID string

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
		CheckDestroy:             sso.CustomRole_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccCustomRoleConfig_Minimal(resourceName, name),
				Check:  sso.CustomRole_GetIDs(resourceFullName, &environmentID, &customRoleID),
			},
			{
				PreConfig: func() {
					sso.CustomRole_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID, customRoleID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the environment
			{
				Config: testAccCustomRoleConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check:  sso.CustomRole_GetIDs(resourceFullName, &environmentID, &customRoleID),
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

func TestAccCustomRole_NewEnv(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_custom_role.%s", resourceName)

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
		CheckDestroy:             sso.CustomRole_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCustomRoleConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
				),
			},
		},
	})
}

func TestAccCustomRole_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_custom_role.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.CustomRole_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCustomRoleConfig_Full(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "can_assign.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "type", "CUSTOM"),
				),
			},
			{
				RefreshState: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					// After a refresh, can_assign should see the second custom role
					resource.TestCheckResourceAttr(resourceFullName, "can_assign.#", "1"),
					resource.TestMatchResourceAttr(resourceFullName, "can_assign.0.id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "type", "CUSTOM"),
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

						return fmt.Sprintf("%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCustomRole_Minimal(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_custom_role.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.CustomRole_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCustomRoleConfig_Minimal(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
				),
			},
		},
	})
}

func TestAccCustomRole_MinimalMaximal(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_custom_role.%s", resourceName)

	name := resourceName

	fullCheckNoCanAssign := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "can_assign.#", "0"),
		resource.TestCheckResourceAttr(resourceFullName, "type", "CUSTOM"),
	)
	fullCheckWithCanAssign := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
		// After a refresh, can_assign should see the second custom role
		resource.TestCheckResourceAttr(resourceFullName, "can_assign.#", "1"),
		resource.TestMatchResourceAttr(resourceFullName, "can_assign.0.id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "type", "CUSTOM"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
			acctest.PreCheckTestAccFlaky(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.CustomRole_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			{
				Config: testAccCustomRoleConfig_Full(resourceName, name),
				Check:  fullCheckNoCanAssign,
			},
			{
				RefreshState: true,
				Check:        fullCheckWithCanAssign,
			},
			// Minimal
			{
				Config: testAccCustomRoleConfig_Minimal(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "can_assign.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "type", "CUSTOM"),
				),
			},
			// Back to full
			{
				Config: testAccCustomRoleConfig_Full(resourceName, name),
				Check:  fullCheckNoCanAssign,
			},
			{
				RefreshState: true,
				Check:        fullCheckWithCanAssign,
			},
		},
	})
}

func TestAccCustomRole_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_custom_role.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.CustomRole_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccCustomRoleConfig_Minimal(resourceName, name),
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

func testAccCustomRoleConfig_NewEnv(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

data "pingone_role" "%[3]s_environment_admin" {
  name = "Environment Admin"
}

resource "pingone_custom_role" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s"
  applicable_to = [
    "ENVIRONMENT"
  ]
  can_be_assigned_by = [
    {
      id = data.pingone_role.%[3]s_environment_admin.id
    }
  ]
  permissions = [
    {
      id = "permissions:read:userRoleAssignments"
    }
  ]
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccCustomRoleConfig_Full(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

data "pingone_role" "%[2]s_environment_admin" {
  name = "Environment Admin"
}

data "pingone_role" "%[2]s_organization_admin" {
  name = "Organization Admin"
}

resource "pingone_custom_role" "%[2]s-dependent_role" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s Dependent Role"
  applicable_to = [
    "ENVIRONMENT",
    "POPULATION",
    "APPLICATION"
  ]
  can_be_assigned_by = [
    {
      id = pingone_custom_role.%[2]s.id
    }
  ]
  description = "My custom dependent role"
  permissions = [
    {
      id = "permissions:read:userRoleAssignments"
    },
    {
      id = "permissions:read:groupRoleAssignments"
    },
  ]
}

resource "pingone_custom_role" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  applicable_to = [
    "ENVIRONMENT",
    "POPULATION",
    "APPLICATION"
  ]
  can_be_assigned_by = [
    {
      id = data.pingone_role.%[2]s_environment_admin.id
    },
    {
      id = data.pingone_role.%[2]s_organization_admin.id
    }
  ]
  description = "My custom role"
  permissions = [
    {
      id = "permissions:read:gatewayRoleAssignments"
    },
    {
      id = "permissions:update:userRoleAssignments"
    }
  ]
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccCustomRoleConfig_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_custom_role" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  applicable_to = [
    "ENVIRONMENT"
  ]
  can_be_assigned_by = []
  permissions = [
    {
      id = "permissions:read:userRoleAssignments"
    }
  ]
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
