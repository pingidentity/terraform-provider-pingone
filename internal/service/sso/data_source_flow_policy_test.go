package sso_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
)

func testAccCheckFlowPolicyDestroy(s *terraform.State) error {
	return nil
}

func TestAccFlowPolicyDataSource_ByNameFull(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_flow_policy.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	name := "Test Flow 1"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckFlowPolicyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFlowPolicyDataSourceConfig_ByNameFull(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "flow_policy_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(dataSourceFullName, "name", name),
					resource.TestCheckResourceAttr(dataSourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "davinci_application.#", "1"),
					resource.TestMatchResourceAttr(dataSourceFullName, "davinci_application.0.id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(dataSourceFullName, "davinci_application.0.name", "Test Application"),
					resource.TestCheckResourceAttr(dataSourceFullName, "trigger.#", "1"),
					resource.TestCheckResourceAttr(dataSourceFullName, "trigger.0.type", "AUTHENTICATION"),
				),
			},
		},
	})
}

func TestAccFlowPolicyDataSource_ByIDFull(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_flow_policy.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	name := "Test Flow 1"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckFlowPolicyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFlowPolicyDataSourceConfig_ByIDFull(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "flow_policy_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(dataSourceFullName, "name", name),
					resource.TestCheckResourceAttr(dataSourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "davinci_application.#", "1"),
					resource.TestMatchResourceAttr(dataSourceFullName, "davinci_application.0.id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(dataSourceFullName, "davinci_application.0.name", "Test Application"),
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
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckFlowPolicyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccFlowPolicyDataSourceConfig_NotFoundByName(resourceName),
				ExpectError: regexp.MustCompile(`Cannot find flow_policy from name`),
			},
			{
				Config:      testAccFlowPolicyDataSourceConfig_NotFoundByID(resourceName),
				ExpectError: regexp.MustCompile("Error when calling `ReadOneFlowPolicy`: The request could not be completed. The requested resource was not found."),
			},
		},
	})
}

func testAccFlowPolicyDataSourceConfig_ByNameFull(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_flow_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}
`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccFlowPolicyDataSourceConfig_ByIDFull(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_flow_policy" "%[2]s-name" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

data "pingone_flow_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  flow_policy_id = pingone_flow_policy.%[2]s-name.id
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccFlowPolicyDataSourceConfig_NotFoundByName(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_flow_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "doesnotexist"
}
`, acctest.GenericSandboxEnvironment(), resourceName)
}

func testAccFlowPolicyDataSourceConfig_NotFoundByID(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_flow_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  flow_policy_id = "9c052a8a-14be-44e4-8f07-2662569994ce" // dummy ID that conforms to UUID v4
}
`, acctest.GenericSandboxEnvironment(), resourceName)
}
