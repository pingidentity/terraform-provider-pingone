// Copyright © 2025 Ping Identity Corporation

package base_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
)

func TestAccRoleDataSource_ByNameFull(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_role.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRoleDataSourceConfig_ByNameFull(resourceName, "Organization Admin"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceFullName, "id"),
					resource.TestCheckResourceAttr(dataSourceFullName, "name", "Organization Admin"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "description"),
					resource.TestMatchResourceAttr(dataSourceFullName, "applicable_to.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "permissions.#", regexp.MustCompile(`^[1-9]\d*$`)),
				),
			},
			{
				Config: testAccRoleDataSourceConfig_ByNameFull(resourceName, "Environment Admin"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceFullName, "id"),
					resource.TestCheckResourceAttr(dataSourceFullName, "name", "Environment Admin"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "description"),
					resource.TestMatchResourceAttr(dataSourceFullName, "applicable_to.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "permissions.#", regexp.MustCompile(`^[1-9]\d*$`)),
				),
			},
			{
				Config: testAccRoleDataSourceConfig_ByNameFull(resourceName, "Identity Data Admin"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceFullName, "id"),
					resource.TestCheckResourceAttr(dataSourceFullName, "name", "Identity Data Admin"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "description"),
					resource.TestMatchResourceAttr(dataSourceFullName, "applicable_to.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "permissions.#", regexp.MustCompile(`^[1-9]\d*$`)),
				),
			},
			{
				Config: testAccRoleDataSourceConfig_ByNameFull(resourceName, "Client Application Developer"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceFullName, "id"),
					resource.TestCheckResourceAttr(dataSourceFullName, "name", "Client Application Developer"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "description"),
					resource.TestMatchResourceAttr(dataSourceFullName, "applicable_to.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "permissions.#", regexp.MustCompile(`^[1-9]\d*$`)),
				),
			},
			{
				Config: testAccRoleDataSourceConfig_ByNameFull(resourceName, "Identity Data Read Only"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceFullName, "id"),
					resource.TestCheckResourceAttr(dataSourceFullName, "name", "Identity Data Read Only"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "description"),
					resource.TestMatchResourceAttr(dataSourceFullName, "applicable_to.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "permissions.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestMatchTypeSetElemNestedAttrs(dataSourceFullName, "permissions.*", map[string]*regexp.Regexp{
						"id":          regexp.MustCompile(`.+`),
						"classifier":  regexp.MustCompile(`.+`),
						"description": regexp.MustCompile(`.+`),
					}),
				),
			},
			{
				Config: testAccRoleDataSourceConfig_ByNameFull(resourceName, "Configuration Read Only"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceFullName, "id"),
					resource.TestCheckResourceAttr(dataSourceFullName, "name", "Configuration Read Only"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "description"),
					resource.TestMatchResourceAttr(dataSourceFullName, "applicable_to.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "permissions.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestMatchTypeSetElemNestedAttrs(dataSourceFullName, "permissions.*", map[string]*regexp.Regexp{
						"id":          regexp.MustCompile(`.+`),
						"classifier":  regexp.MustCompile(`.+`),
						"description": regexp.MustCompile(`.+`),
					}),
				),
			},
			{
				Config: testAccRoleDataSourceConfig_ByNameFull(resourceName, "DaVinci Admin"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceFullName, "id"),
					resource.TestCheckResourceAttr(dataSourceFullName, "name", "DaVinci Admin"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "description"),
					resource.TestMatchResourceAttr(dataSourceFullName, "applicable_to.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "permissions.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestMatchTypeSetElemNestedAttrs(dataSourceFullName, "permissions.*", map[string]*regexp.Regexp{
						"id":          regexp.MustCompile(`.+`),
						"classifier":  regexp.MustCompile(`.+`),
						"description": regexp.MustCompile(`.+`),
					}),
				),
			},
			{
				Config: testAccRoleDataSourceConfig_ByNameFull(resourceName, "DaVinci Admin Read Only"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceFullName, "id"),
					resource.TestCheckResourceAttr(dataSourceFullName, "name", "DaVinci Admin Read Only"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "description"),
					resource.TestMatchResourceAttr(dataSourceFullName, "applicable_to.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "permissions.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestMatchTypeSetElemNestedAttrs(dataSourceFullName, "permissions.*", map[string]*regexp.Regexp{
						"id":          regexp.MustCompile(`.+`),
						"classifier":  regexp.MustCompile(`.+`),
						"description": regexp.MustCompile(`.+`),
					}),
				),
			},
		},
	})
}

func TestAccRoleDataSource_ByIdFull(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_role.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRoleDataSourceConfig_ByIdFull(resourceName, "Organization Admin"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceFullName, "id"),
					resource.TestCheckResourceAttr(dataSourceFullName, "name", "Organization Admin"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "description"),
					resource.TestMatchResourceAttr(dataSourceFullName, "applicable_to.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "permissions.#", regexp.MustCompile(`^[1-9]\d*$`)),
				),
			},
		},
	})
}

func TestAccRoleDataSource_NotFound(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccRoleDataSourceConfig_NotFoundByName(resourceName),
				ExpectError: regexp.MustCompile("Invalid Attribute Value Match"),
			},
			// {
			// 	Config:      testAccRoleDataSourceConfig_NotFoundByID(resourceName),
			// 	ExpectError: regexp.MustCompile("Error when calling `GetRole`: Role not found for id: 9c052a8a-14be-44e4-8f07-2662569994ce and environmentId: [0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}"),
			// },
		},
	})
}

func testAccRoleDataSourceConfig_ByNameFull(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

data "pingone_role" "%[2]s" {
  name = "%[3]s"
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRoleDataSourceConfig_ByIdFull(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

data "pingone_role" "%[2]s-lookup" {
  name = "%[3]s"
}

data "pingone_role" "%[2]s" {
  role_id = data.pingone_role.%[2]s-lookup.id
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRoleDataSourceConfig_NotFoundByName(resourceName string) string {
	return fmt.Sprintf(`
		%[1]s

data "pingone_role" "%[2]s" {
  name = "doesnotexist"
}`, acctest.GenericSandboxEnvironment(), resourceName)
}
