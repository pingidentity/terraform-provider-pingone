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

func TestAccApplicationSignOnPolicyAssignmentsDataSource_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullName := fmt.Sprintf("data.pingone_application_sign_on_policy_assignments.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.SignOnPolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationSignOnPolicyAssignmentsDataSourceConfig_Full(resourceName, name),
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

func TestAccApplicationSignOnPolicyAssignmentsDataSource_NotFound(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullName := fmt.Sprintf("data.pingone_application_sign_on_policy_assignments.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.SignOnPolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationSignOnPolicyAssignmentsDataSourceConfig_NotFound(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "application_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "ids.#", "0"),
				),
			},
		},
	})
}

func testAccApplicationSignOnPolicyAssignmentsDataSourceConfig_Full(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
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

resource "pingone_sign_on_policy" "%[2]s-1" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[2]s_1"
}

resource "pingone_sign_on_policy" "%[2]s-2" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[2]s_2"
}

resource "pingone_application_sign_on_policy_assignment" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  application_id = pingone_application.%[2]s.id

  sign_on_policy_id = pingone_sign_on_policy.%[2]s-1.id
  priority          = 2
}

resource "pingone_application_sign_on_policy_assignment" "%[2]s-2" {
  environment_id = data.pingone_environment.general_test.id
  application_id = pingone_application.%[2]s.id

  sign_on_policy_id = pingone_sign_on_policy.%[2]s-2.id
  priority          = 1
}

data "pingone_application_sign_on_policy_assignments" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  application_id = pingone_application.%[2]s.id

  depends_on = [
    pingone_application_sign_on_policy_assignment.%[2]s,
    pingone_application_sign_on_policy_assignment.%[2]s-2,
  ]
}
`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccApplicationSignOnPolicyAssignmentsDataSourceConfig_NotFound(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
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

data "pingone_application_sign_on_policy_assignments" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  application_id = pingone_application.%[2]s.id
}
`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
