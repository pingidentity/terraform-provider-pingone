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

func TestAccFlowPoliciesDataSource_BySCIMFilter(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullName := fmt.Sprintf("data.pingone_flow_policies.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.FlowPolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFlowPoliciesDataSourceConfig_BySCIMFilter(resourceName, `(trigger.type eq \"AUTHENTICATION\")`, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "ids.#", "2"),
					resource.TestMatchResourceAttr(dataSourceFullName, "ids.0", verify.P1DVResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "ids.1", verify.P1DVResourceIDRegexpFullString),
				),
			},
		},
	})
}

func TestAccFlowPoliciesDataSource_ByDataFilters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullName := fmt.Sprintf("data.pingone_flow_policies.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.FlowPolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFlowPoliciesDataSourceConfig_ByDataFilters(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "ids.#", "2"),
					resource.TestMatchResourceAttr(dataSourceFullName, "ids.0", verify.P1DVResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "ids.1", verify.P1DVResourceIDRegexpFullString),
				),
			},
		},
	})
}

func TestAccFlowPoliciesDataSource_NotFound(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullName := fmt.Sprintf("data.pingone_flow_policies.%s", resourceName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.FlowPolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFlowPoliciesDataSourceConfig_NotFound(resourceName, `(trigger.type eq \"NOTAUTHENTICATION\")`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "ids.#", "0"),
				),
			},
		},
	})
}

func testAccFlowPoliciesDataSourceConfig_BySCIMFilter(resourceName, filter, name string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_flow_policies" "%[2]s" {
  environment_id = data.pingone_environment.davinci_test.id

  scim_filter = "%[4]s"
}
`, acctest.DaVinciFlowPolicySandboxEnvironment(), resourceName, name, filter)
}

func testAccFlowPoliciesDataSourceConfig_ByDataFilters(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_flow_policies" "%[2]s" {
  environment_id = data.pingone_environment.davinci_test.id

  data_filters = [
    {
      name   = "trigger.type"
      values = ["AUTHENTICATION"]
    }
  ]
}`, acctest.DaVinciFlowPolicySandboxEnvironment(), resourceName, name)
}

func testAccFlowPoliciesDataSourceConfig_NotFound(resourceName, filter string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_flow_policies" "%[2]s" {
  environment_id = data.pingone_environment.davinci_test.id

  scim_filter = "%[3]s"
}
`, acctest.DaVinciFlowPolicySandboxEnvironment(), resourceName, filter)
}
