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

func TestAccResourceScopeDataSource_ByNameFull(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_resource_scope.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.ResourceScope_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceScopeDataSourceConfig_ByNameFull(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "resource_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "custom_resource_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "resource_type", "CUSTOM"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "name", resourceFullName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "description", resourceFullName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "schema_attributes", resourceFullName, "schema_attributes"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "mapped_claims", resourceFullName, "mapped_claims"),
				),
			},
		},
	})
}

func TestAccResourceScopeDataSource_ByNameSystem(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_resource_scope.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.ResourceScope_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceScopeDataSourceConfig_ByNameSystem(resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "resource_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckNoResourceAttr(dataSourceFullName, "custom_resource_id"),
					resource.TestCheckResourceAttr(dataSourceFullName, "resource_type", "OPENID_CONNECT"),
					resource.TestMatchResourceAttr(dataSourceFullName, "resource_scope_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "name", "email"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "description"),
					resource.TestCheckResourceAttr(dataSourceFullName, "schema_attributes.#", "0"),
					resource.TestCheckResourceAttr(dataSourceFullName, "mapped_claims.#", "0"),
				),
			},
		},
	})
}

func TestAccResourceScopeDataSource_ByIDFull(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_resource_scope.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.ResourceScope_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceScopeDataSourceConfig_ByIDFull(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "resource_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "custom_resource_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "resource_type", "CUSTOM"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "name", resourceFullName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "description", resourceFullName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "schema_attributes", resourceFullName, "schema_attributes"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "mapped_claims", resourceFullName, "mapped_claims"),
				),
			},
		},
	})
}

func TestAccResourceScopeDataSource_ByIDSchemaAttributes(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_resource_scope_pingone_api.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.pingone_resource_scope.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.ResourceScope_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceScopeDataSourceConfig_ByIDSchemaAttributes(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceFullName, "schema_attributes", resourceFullName, "schema_attributes"),
					resource.TestCheckResourceAttr(dataSourceFullName, "mapped_claims.#", "0"),
				),
			},
		},
	})
}

func TestAccResourceScopeDataSource_ByIDMappedClaims(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_resource_scope_openid.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.pingone_resource_scope.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.ResourceScope_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceScopeDataSourceConfig_ByIDMappedClaims(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceFullName, "schema_attributes.#", "0"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "mapped_claims", resourceFullName, "mapped_claims"),
				),
			},
		},
	})
}

func TestAccResourceScopeDataSource_NotFound(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.ResourceScope_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccResourceScopeDataSourceConfig_NotFoundByName(resourceName),
				ExpectError: regexp.MustCompile("Cannot find resource scope"),
			},
			{
				Config:      testAccResourceScopeDataSourceConfig_NotFoundByID(resourceName),
				ExpectError: regexp.MustCompile("Error when calling `ReadOneResourceScope`: The request could not be completed. The requested resource was not found."),
			},
		},
	})
}

func testAccResourceScopeDataSourceConfig_ByNameFull(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_resource" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_resource_scope" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  resource_id    = pingone_resource.%[2]s.id

  name = "%[3]s"
}

data "pingone_resource_scope" "%[3]s" {
  environment_id     = data.pingone_environment.general_test.id
  resource_type      = "CUSTOM"
  custom_resource_id = pingone_resource.%[2]s.id

  name = "%[3]s"

  depends_on = [
    pingone_resource_scope.%[2]s
  ]
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccResourceScopeDataSourceConfig_ByIDFull(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_resource" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_resource_scope" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  resource_id    = pingone_resource.%[2]s.id

  name = "%[3]s"
}

data "pingone_resource_scope" "%[2]s" {
  environment_id     = data.pingone_environment.general_test.id
  resource_type      = "CUSTOM"
  custom_resource_id = pingone_resource.%[2]s.id

  resource_scope_id = pingone_resource_scope.%[2]s.id
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccResourceScopeDataSourceConfig_ByNameSystem(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_resource_scope" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  resource_type  = "OPENID_CONNECT"

  name = "email"
}`, acctest.GenericSandboxEnvironment(), resourceName)
}

func testAccResourceScopeDataSourceConfig_ByIDSchemaAttributes(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_resource_scope_pingone_api" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "p1:read:user:%[3]s"

  schema_attributes = [
    "name.given",
    "name.family",
  ]
}

data "pingone_resource_scope" "%[2]s" {
  environment_id    = data.pingone_environment.general_test.id
  resource_type     = "PINGONE_API"
  resource_scope_id = pingone_resource_scope_pingone_api.%[2]s.id

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccResourceScopeDataSourceConfig_ByIDMappedClaims(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_resource_attribute" "%[2]s-1" {
  environment_id = data.pingone_environment.general_test.id
  resource_type  = "OPENID_CONNECT"

  name  = "%[3]s-1"
  value = "$${user.name.given}"
}

resource "pingone_resource_attribute" "%[2]s-2" {
  environment_id = data.pingone_environment.general_test.id
  resource_type  = "OPENID_CONNECT"

  name  = "%[3]s-2"
  value = "$${user.name.family}"
}

resource "pingone_resource_attribute" "%[2]s-3" {
  environment_id = data.pingone_environment.general_test.id
  resource_type  = "OPENID_CONNECT"

  name  = "%[3]s-3"
  value = "$${user.email}"
}

resource "pingone_resource_scope_openid" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name        = "%[3]s"
  description = "My resource scope"

  mapped_claims = [
    pingone_resource_attribute.%[2]s-2.id,
    pingone_resource_attribute.%[2]s-3.id,
    pingone_resource_attribute.%[2]s-1.id
  ]
}

data "pingone_resource_scope" "%[2]s" {
  environment_id    = data.pingone_environment.general_test.id
  resource_type     = "OPENID_CONNECT"
  resource_scope_id = pingone_resource_scope_openid.%[2]s.id

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccResourceScopeDataSourceConfig_NotFoundByName(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_resource" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[2]s"
}

data "pingone_resource_scope" "%[2]s" {
  environment_id     = data.pingone_environment.general_test.id
  resource_type      = "CUSTOM"
  custom_resource_id = pingone_resource.%[2]s.id

  name = "doesnotexist"
}`, acctest.GenericSandboxEnvironment(), resourceName)
}

func testAccResourceScopeDataSourceConfig_NotFoundByID(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_resource" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[2]s"
}

data "pingone_resource_scope" "%[2]s" {
  environment_id     = data.pingone_environment.general_test.id
  resource_type      = "CUSTOM"
  custom_resource_id = pingone_resource.%[2]s.id

  resource_scope_id = "9c052a8a-14be-44e4-8f07-2662569994ce" // dummy ID that conforms to UUID v4
}`, acctest.GenericSandboxEnvironment(), resourceName)
}
