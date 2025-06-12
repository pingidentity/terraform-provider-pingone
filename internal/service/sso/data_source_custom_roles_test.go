// Copyright © 2025 Ping Identity Corporation

package sso_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/sso"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccCustomRolesDataSource_GetAll(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullName := fmt.Sprintf("data.pingone_custom_roles.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.CustomRole_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCustomRolesDataSourceConfig_GetAll(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "ids.#", "3"),
					resource.TestMatchResourceAttr(dataSourceFullName, "ids.0", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "ids.1", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "ids.2", verify.P1ResourceIDRegexpFullString),
				),
			},
		},
	})
}

func testAccCustomRolesDataSourceConfig_GetAll(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_role" "%[2]s_environment_admin" {
  name = "Environment Admin"
}

data "pingone_role" "%[2]s_organization_admin" {
  name = "Organization Admin"
}

resource "pingone_custom_role" "%[2]s-dependent-role" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s Datasource Dependent Role"
  applicable_to = [
    "ENVIRONMENT",
    "POPULATION"
  ]
  can_be_assigned_by = [
    {
      id = pingone_custom_role.%[2]s-parent.id
    }
  ]
  description = "My custom dependent role for datasource test"
  permissions = [
    {
      id = "permissions:read:userRoleAssignments"
    },
    {
      id = "permissions:read:groupRoleAssignments"
    },
  ]
}

resource "pingone_custom_role" "%[2]s-parent" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  applicable_to = [
    "ENVIRONMENT",
    "POPULATION"
  ]
  can_be_assigned_by = [
    {
      id = data.pingone_role.%[2]s_environment_admin.id
    },
    {
      id = data.pingone_role.%[2]s_organization_admin.id
    }
  ]
  description = "My custom role for datasource test"
  permissions = [
    {
      id = "permissions:read:gatewayRoleAssignments"
    },
    {
      id = "permissions:update:userRoleAssignments"
    }
  ]
}

resource "pingone_custom_role" "%[2]s-simple" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s simple"
  applicable_to = [
    "ENVIRONMENT"
  ]
  can_be_assigned_by = [
    {
      id = data.pingone_role.%[2]s_environment_admin.id
    }
  ]
  description = "My simple custom role for datasource test"
  permissions = [
    {
      id = "permissions:update:userRoleAssignments"
    }
  ]
}

data "pingone_custom_roles" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  depends_on = [
    pingone_custom_role.%[2]s-parent,
    pingone_custom_role.%[2]s-dependent-role,
    pingone_custom_role.%[2]s-simple,
  ]
}
`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
