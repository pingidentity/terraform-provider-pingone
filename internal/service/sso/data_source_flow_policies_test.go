package sso_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
)

func TestAccFlowPoliciesDataSource_BySCIMFilter(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullName := fmt.Sprintf("data.pingone_flow_policies.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckFlowPolicyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFlowPoliciesDataSourceConfig_BySCIMFilter(resourceName, "(trigger.type eq \"AUTHENTICATION\")", name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(dataSourceFullName, "ids.#", "2"),
					resource.TestMatchResourceAttr(dataSourceFullName, "ids.0", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "ids.1", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
				),
			},
		},
	})
}

func TestAccFlowPoliciesDataSource_ByDataFilter(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullName := fmt.Sprintf("data.pingone_flow_policies.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckFlowPolicyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFlowPoliciesDataSourceConfig_ByDataFilter1(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(dataSourceFullName, "ids.#", "2"),
					resource.TestMatchResourceAttr(dataSourceFullName, "ids.0", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "ids.1", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
				),
			},
			// {
			// 	Config: testAccFlowPoliciesDataSourceConfig_ByDataFilter2(resourceName, name),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		resource.TestMatchResourceAttr(dataSourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
			// 		resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
			// 		resource.TestCheckResourceAttr(dataSourceFullName, "ids.#", "3"),
			// 		resource.TestMatchResourceAttr(dataSourceFullName, "ids.0", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
			// 		resource.TestMatchResourceAttr(dataSourceFullName, "ids.1", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
			// 		resource.TestMatchResourceAttr(dataSourceFullName, "ids.2", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
			// 	),
			// },
		},
	})
}

func TestAccFlowPoliciesDataSource_NotFound(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullName := fmt.Sprintf("data.pingone_flow_policies.%s", resourceName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckFlowPolicyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFlowPoliciesDataSourceConfig_NotFound(resourceName, "(trigger.type eq \"NOTAUTHENTICATION\")"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
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
  environment_id = data.pingone_environment.general_test.id

  scim_filter = "%[4]s"
}
`, acctest.GenericSandboxEnvironment(), resourceName, name, filter)
}

func testAccFlowPoliciesDataSourceConfig_ByDataFilter1(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_flow_policies" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  data_filter {
    name   = "trigger.type"
    values = ["AUTHENTICATION"]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

// func testAccFlowPoliciesDataSourceConfig_ByDataFilter2(resourceName, name string) string {
// 	return fmt.Sprintf(`
// 	%[1]s

// resource "pingone_population" "%[2]s-1" {
//   environment_id = data.pingone_environment.general_test.id
//   name           = "%[3]s-1"
// }

// resource "pingone_population" "%[2]s-2" {
//   environment_id = data.pingone_environment.general_test.id
//   name           = "%[3]s-2"
// }

// resource "pingone_population" "%[2]s-3" {
//   environment_id = data.pingone_environment.general_test.id
//   name           = "%[3]s-3"
// }

// data "pingone_flow_policies" "%[2]s" {
//   environment_id = data.pingone_environment.general_test.id

//   data_filter {
//     name   = "name"
//     values = ["%[3]s-1", "%[3]s-2", "%[3]s-3", ]
//   }

//   depends_on = [
//     pingone_population.%[2]s-1,
//     pingone_population.%[2]s-2,
//     pingone_population.%[2]s-3,
//   ]
// }`, acctest.GenericSandboxEnvironment(), resourceName, name)
// }

func testAccFlowPoliciesDataSourceConfig_NotFound(resourceName, filter string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_flow_policies" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  scim_filter = "%[3]s"
}
`, acctest.GenericSandboxEnvironment(), resourceName, filter)
}
