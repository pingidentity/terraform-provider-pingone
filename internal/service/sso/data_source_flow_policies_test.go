// Copyright © 2025 Ping Identity Corporation

package sso_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	acctestlegacysdk "github.com/pingidentity/terraform-provider-pingone/internal/acctest/legacysdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/sso"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccFlowPoliciesDataSource_BySCIMFilter(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullName := fmt.Sprintf("data.pingone_flow_policies.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			t.Skip("Skipping until DaVinci capability merged")
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.FlowPolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFlowPoliciesDataSourceConfig_BySCIMFilter(environmentName, licenseID, resourceName, `(trigger.type eq \"AUTHENTICATION\")`),
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

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			t.Skip("Skipping until DaVinci capability merged")
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.FlowPolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFlowPoliciesDataSourceConfig_ByDataFilters(environmentName, licenseID, resourceName),
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

func testAccFlowPoliciesDataSourceConfig_BySCIMFilter(environmentName, licenseID, resourceName, filter string) string {
	return fmt.Sprintf(`
%[1]s

data "pingone_flow_policies" "%[2]s" {
  environment_id = pingone_environment.%[2]s.id

  scim_filter = "%[3]s"

  depends_on = [
    davinci_application_flow_policy.%[2]s-1,
    davinci_application_flow_policy.%[2]s-2,
  ]
}
`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), resourceName, filter)
}

func testAccFlowPoliciesDataSourceConfig_ByDataFilters(environmentName, licenseID, resourceName string) string {
	return fmt.Sprintf(`
%[1]s

data "pingone_flow_policies" "%[2]s" {
  environment_id = pingone_environment.%[2]s.id

  data_filters = [
    {
      name   = "trigger.type"
      values = ["AUTHENTICATION"]
    }
  ]

  depends_on = [
    davinci_application_flow_policy.%[2]s-1,
    davinci_application_flow_policy.%[2]s-2,
  ]
}`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), resourceName)
}

func testAccFlowPoliciesDataSourceConfig_NotFound(resourceName, filter string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_flow_policies" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  scim_filter = "%[3]s"
}
`, acctest.GenericSandboxEnvironment(), resourceName, filter)
}
