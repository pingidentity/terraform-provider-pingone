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

func TestAccApplicationFlowPolicyAssignmentsDataSource_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullName := fmt.Sprintf("data.pingone_application_flow_policy_assignments.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckFeatureFlag(t, acctest.ENUMFEATUREFLAG_DAVINCI)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.FlowPolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationFlowPolicyAssignmentsDataSourceConfig_Full(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "application_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "ids.#", "2"),
					resource.TestMatchResourceAttr(dataSourceFullName, "ids.0", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "ids.1", verify.P1ResourceIDRegexpFullString),
				),
			},
		},
	})
}

func TestAccApplicationFlowPolicyAssignmentsDataSource_NotFound(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullName := fmt.Sprintf("data.pingone_application_flow_policy_assignments.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckFeatureFlag(t, acctest.ENUMFEATUREFLAG_DAVINCI)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.FlowPolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationFlowPolicyAssignmentsDataSourceConfig_NotFound(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "application_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "ids.#", "0"),
				),
			},
		},
	})
}

func testAccApplicationFlowPolicyAssignmentsDataSourceConfig_Full(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.davinci_test.id
  name           = "%[3]s"
  enabled        = true

  oidc_options = {
    type                       = "SINGLE_PAGE_APP"
    grant_types                = ["AUTHORIZATION_CODE"]
    response_types             = ["CODE"]
    pkce_enforcement           = "S256_REQUIRED"
    token_endpoint_auth_method = "NONE"
    redirect_uris              = ["https://www.pingidentity.com"]
  }
}

data "pingone_flow_policies" "%[2]s" {
  environment_id = data.pingone_environment.davinci_test.id

  scim_filter = "(trigger.type eq \"AUTHENTICATION\")"
}

resource "pingone_application_flow_policy_assignment" "%[2]s" {
  environment_id = data.pingone_environment.davinci_test.id
  application_id = pingone_application.%[2]s.id

  flow_policy_id = data.pingone_flow_policies.%[2]s.ids[0]

  priority = 2
}

resource "pingone_application_flow_policy_assignment" "%[2]s-2" {
  environment_id = data.pingone_environment.davinci_test.id
  application_id = pingone_application.%[2]s.id

  flow_policy_id = data.pingone_flow_policies.%[2]s.ids[1]

  priority = 1
}

data "pingone_application_flow_policy_assignments" "%[2]s" {
  environment_id = data.pingone_environment.davinci_test.id

  application_id = pingone_application.%[2]s.id

  depends_on = [
    pingone_application_flow_policy_assignment.%[2]s,
    pingone_application_flow_policy_assignment.%[2]s-2,
  ]
}
`, acctest.DaVinciFlowPolicySandboxEnvironment(), resourceName, name)
}

func testAccApplicationFlowPolicyAssignmentsDataSourceConfig_NotFound(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.davinci_test.id
  name           = "%[3]s"
  enabled        = true

  oidc_options = {
    type                       = "SINGLE_PAGE_APP"
    grant_types                = ["AUTHORIZATION_CODE"]
    response_types             = ["CODE"]
    pkce_enforcement           = "S256_REQUIRED"
    token_endpoint_auth_method = "NONE"
    redirect_uris              = ["https://www.pingidentity.com"]
  }
}

data "pingone_application_flow_policy_assignments" "%[2]s" {
  environment_id = data.pingone_environment.davinci_test.id

  application_id = pingone_application.%[2]s.id
}
`, acctest.DaVinciFlowPolicySandboxEnvironment(), resourceName, name)
}
