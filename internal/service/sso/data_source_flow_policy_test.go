package sso_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func testAccCheckFlowPolicyDestroy(s *terraform.State) error {
	return nil
}

func TestAccFlowPolicyDataSource_ByIDFull(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_flow_policy.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironmentFeatureFlag(t, acctest.ENUMFEATUREFLAG_DAVINCI) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckFlowPolicyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFlowPolicyDataSourceConfig_ByIDFull(resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1DVResourceIDRegexp),
					resource.TestMatchResourceAttr(dataSourceFullName, "flow_policy_id", verify.P1DVResourceIDRegexp),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(dataSourceFullName, "name", regexp.MustCompile(`^Test Flow Policy( 2)?$`)),
					resource.TestCheckResourceAttr(dataSourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "davinci_application.#", "1"),
					resource.TestMatchResourceAttr(dataSourceFullName, "davinci_application.0.id", verify.P1DVResourceIDRegexp),
					resource.TestMatchResourceAttr(dataSourceFullName, "davinci_application.0.name", regexp.MustCompile(`^Test Application( 2)?$`)),
					resource.TestCheckResourceAttr(dataSourceFullName, "trigger.#", "1"),
					resource.TestCheckResourceAttr(dataSourceFullName, "trigger.0.type", "AUTHENTICATION"),
				),
			},
		},
	})
}

func TestAccFlowPolicyDataSource_NotFound(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironmentFeatureFlag(t, acctest.ENUMFEATUREFLAG_DAVINCI) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckFlowPolicyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccFlowPolicyDataSourceConfig_NotFoundByID(resourceName),
				ExpectError: regexp.MustCompile("Error when calling `ReadOneFlowPolicy`: Unable to find an active Flow Policy with ID: '[a-f0-9]{32}' in Environment '[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}'"),
			},
		},
	})
}

func testAccFlowPolicyDataSourceConfig_ByIDFull(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_flow_policies" "%[2]s" {
  environment_id = data.pingone_environment.davinci_test.id

  scim_filter = "(trigger.type eq \"AUTHENTICATION\")"
}

data "pingone_flow_policy" "%[2]s" {
  environment_id = data.pingone_environment.davinci_test.id

  flow_policy_id = data.pingone_flow_policies.%[2]s.ids[0]
}`, acctest.DaVinciFlowPolicySandboxEnvironment(), resourceName)
}

func testAccFlowPolicyDataSourceConfig_NotFoundByID(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_flow_policy" "%[2]s" {
  environment_id = data.pingone_environment.davinci_test.id

  flow_policy_id = "07ae09dea68df5530269c242487fbaf8" // dummy ID
}
`, acctest.DaVinciFlowPolicySandboxEnvironment(), resourceName)
}
