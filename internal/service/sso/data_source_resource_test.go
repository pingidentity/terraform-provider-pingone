package sso_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccResourceDataSource_ByNameFull(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_resource.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckResourceDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceDataSourceConfig_ByNameFull(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(dataSourceFullName, "resource_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "name", resourceFullName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "description", resourceFullName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "type", resourceFullName, "type"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "audience", resourceFullName, "audience"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "access_token_validity_seconds", resourceFullName, "access_token_validity_seconds"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "introspect_endpoint_auth_method", resourceFullName, "introspect_endpoint_auth_method"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "client_secret", resourceFullName, "client_secret"),
				),
			},
		},
	})
}

func TestAccResourceDataSource_ByNameSystem(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_resource.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceDataSourceConfig_ByNameSystem(resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(dataSourceFullName, "resource_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(dataSourceFullName, "name", "openid"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "description"),
					resource.TestCheckResourceAttr(dataSourceFullName, "type", "OPENID_CONNECT"),
					resource.TestCheckNoResourceAttr(dataSourceFullName, "audience"),
					resource.TestCheckNoResourceAttr(dataSourceFullName, "access_token_validity_seconds"),
					resource.TestCheckNoResourceAttr(dataSourceFullName, "introspect_endpoint_auth_method"),
					resource.TestCheckNoResourceAttr(dataSourceFullName, "client_secret"),
				),
			},
		},
	})
}

func TestAccResourceDataSource_ByIDFull(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_resource.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckResourceDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceDataSourceConfig_ByIDFull(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(dataSourceFullName, "resource_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "name", resourceFullName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "description", resourceFullName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "type", resourceFullName, "type"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "audience", resourceFullName, "audience"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "access_token_validity_seconds", resourceFullName, "access_token_validity_seconds"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "introspect_endpoint_auth_method", resourceFullName, "introspect_endpoint_auth_method"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "client_secret", resourceFullName, "client_secret"),
				),
			},
		},
	})
}

func TestAccResourceDataSource_NotFound(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckResourceDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccResourceDataSourceConfig_NotFoundByName(resourceName),
				ExpectError: regexp.MustCompile("Cannot find resource from name"),
			},
			{
				Config:      testAccResourceDataSourceConfig_NotFoundByID(resourceName),
				ExpectError: regexp.MustCompile("Error when calling `ReadOneResource`: The request could not be completed. The requested resource was not found."),
			},
		},
	})
}

func testAccResourceDataSourceConfig_ByNameFull(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_resource" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name        = "%[3]s"
  description = "Test Resource"

  audience                      = "%[3]s"
  access_token_validity_seconds = 7200
}

data "pingone_resource" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  depends_on = [
    pingone_resource.%[2]s
  ]
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccResourceDataSourceConfig_ByIDFull(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_resource" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name        = "%[3]s"
  description = "Test Resource"

  audience                      = "%[3]s"
  access_token_validity_seconds = 7200
}

data "pingone_resource" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  resource_id = pingone_resource.%[2]s.id
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccResourceDataSourceConfig_ByNameSystem(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_resource" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "openid"
}`, acctest.GenericSandboxEnvironment(), resourceName)
}

func testAccResourceDataSourceConfig_NotFoundByName(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_resource" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "doesnotexist"
}`, acctest.GenericSandboxEnvironment(), resourceName)
}

func testAccResourceDataSourceConfig_NotFoundByID(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_resource" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  resource_id = "9c052a8a-14be-44e4-8f07-2662569994ce" // dummy ID that conforms to UUID v4
}`, acctest.GenericSandboxEnvironment(), resourceName)
}
