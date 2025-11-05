// Copyright © 2025 Ping Identity Corporation

package sso_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/sso"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccResourceScopesDataSource_ByResourceID(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullName := fmt.Sprintf("data.pingone_resource_scopes.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.ResourceScope_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceScopesDataSourceConfig_ByResourceID(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "resource_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "ids.#", "3"),
					resource.TestMatchResourceAttr(dataSourceFullName, "ids.0", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "ids.1", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "ids.2", verify.P1ResourceIDRegexpFullString),
				),
			},
		},
	})
}

func TestAccResourceScopesDataSource_OpenID(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullName := fmt.Sprintf("data.pingone_resource_scopes.%s", resourceName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.ResourceScope_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceScopesDataSourceConfig_OpenID(resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "resource_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttrWith(dataSourceFullName, "ids.#", func(value string) error {
						if value == "0" {
							return fmt.Errorf("expected at least 1 scope ID, got 0")
						}
						return nil
					}),
				),
			},
		},
	})
}

func TestAccResourceScopesDataSource_PingOneAPI(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullName := fmt.Sprintf("data.pingone_resource_scopes.%s", resourceName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.ResourceScope_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceScopesDataSourceConfig_PingOneAPI(resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "resource_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttrWith(dataSourceFullName, "ids.#", func(value string) error {
						if value == "0" {
							return fmt.Errorf("expected at least 1 scope ID, got 0")
						}
						return nil
					}),
				),
			},
		},
	})
}

func TestAccResourceScopesDataSource_NotFound(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.ResourceScope_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccResourceScopesDataSourceConfig_NotFoundByID(resourceName),
				ExpectError: regexp.MustCompile("The requested resource was not found"),
			},
		},
	})
}

func testAccResourceScopesDataSourceConfig_ByResourceID(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_resource" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
}

resource "pingone_resource_scope" "%[2]s-1" {
  environment_id = data.pingone_environment.general_test.id
  resource_id    = pingone_resource.%[2]s.id
  name           = "%[3]s-1"
  description    = "Test scope 1"
}

resource "pingone_resource_scope" "%[2]s-2" {
  environment_id = data.pingone_environment.general_test.id
  resource_id    = pingone_resource.%[2]s.id
  name           = "%[3]s-2"
  description    = "Test scope 2"
}

resource "pingone_resource_scope" "%[2]s-3" {
  environment_id = data.pingone_environment.general_test.id
  resource_id    = pingone_resource.%[2]s.id
  name           = "%[3]s-3"
  description    = "Test scope 3"
}

data "pingone_resource_scopes" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  resource_id    = pingone_resource.%[2]s.id

  depends_on = [
    pingone_resource_scope.%[2]s-1,
    pingone_resource_scope.%[2]s-2,
    pingone_resource_scope.%[2]s-3,
  ]
}
`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccResourceScopesDataSourceConfig_OpenID(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_resource" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "openid"
}

data "pingone_resource_scopes" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  resource_id    = data.pingone_resource.%[2]s.id
}
`, acctest.GenericSandboxEnvironment(), resourceName)
}

func testAccResourceScopesDataSourceConfig_PingOneAPI(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_resource" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "PingOne API"
}

data "pingone_resource_scopes" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  resource_id    = data.pingone_resource.%[2]s.id
}
`, acctest.GenericSandboxEnvironment(), resourceName)
}

func testAccResourceScopesDataSourceConfig_NotFoundByID(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_resource_scopes" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  resource_id    = "9c052a8a-14be-44e4-8f07-2662569994ce" // dummy generic ID
}
`, acctest.GenericSandboxEnvironment(), resourceName)
}
