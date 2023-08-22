package sso_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccUserDataSource_ByNameFull(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_user.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccUserDataSourceConfig_ByNameFull(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "user_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "username", resourceFullName, "username"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "email", resourceFullName, "email"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "email_verified", resourceFullName, "email_verified"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "population_id", resourceFullName, "population_id"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "enabled", resourceFullName, "enabled"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "status", resourceFullName, "status"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "account", resourceFullName, "account"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "address", resourceFullName, "address"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "external_id", resourceFullName, "external_id"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "identity_provider", resourceFullName, "identity_provider"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "user_lifecycle", resourceFullName, "user_lifecycle"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "locale", resourceFullName, "locale"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "mfa_enabled", resourceFullName, "mfa_enabled"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "mobile_phone", resourceFullName, "mobile_phone"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "name", resourceFullName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "nickname", resourceFullName, "nickname"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "password", resourceFullName, "password"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "photo", resourceFullName, "photo"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "preferred_language", resourceFullName, "preferred_language"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "primary_phone", resourceFullName, "primary_phone"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "timezone", resourceFullName, "timezone"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "title", resourceFullName, "title"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "type", resourceFullName, "type"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "verify_status", resourceFullName, "verify_status"),
				),
			},
		},
	})
}

func TestAccUserDataSource_ByEmailFull(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_user.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccUserDataSourceConfig_ByEmailFull(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "user_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "username", resourceFullName, "username"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "email", resourceFullName, "email"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "email_verified", resourceFullName, "email_verified"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "population_id", resourceFullName, "population_id"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "enabled", resourceFullName, "enabled"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "status", resourceFullName, "status"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "account", resourceFullName, "account"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "address", resourceFullName, "address"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "external_id", resourceFullName, "external_id"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "identity_provider", resourceFullName, "identity_provider"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "user_lifecycle", resourceFullName, "user_lifecycle"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "locale", resourceFullName, "locale"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "mfa_enabled", resourceFullName, "mfa_enabled"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "mobile_phone", resourceFullName, "mobile_phone"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "name", resourceFullName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "nickname", resourceFullName, "nickname"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "password", resourceFullName, "password"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "photo", resourceFullName, "photo"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "preferred_language", resourceFullName, "preferred_language"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "primary_phone", resourceFullName, "primary_phone"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "timezone", resourceFullName, "timezone"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "title", resourceFullName, "title"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "type", resourceFullName, "type"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "verify_status", resourceFullName, "verify_status"),
				),
			},
		},
	})
}

func TestAccUserDataSource_ByIDFull(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_user.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccUserDataSourceConfig_ByIDFull(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "user_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "username", resourceFullName, "username"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "email", resourceFullName, "email"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "email_verified", resourceFullName, "email_verified"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "population_id", resourceFullName, "population_id"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "enabled", resourceFullName, "enabled"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "status", resourceFullName, "status"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "account", resourceFullName, "account"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "address", resourceFullName, "address"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "external_id", resourceFullName, "external_id"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "identity_provider", resourceFullName, "identity_provider"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "user_lifecycle", resourceFullName, "user_lifecycle"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "locale", resourceFullName, "locale"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "mfa_enabled", resourceFullName, "mfa_enabled"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "mobile_phone", resourceFullName, "mobile_phone"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "name", resourceFullName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "nickname", resourceFullName, "nickname"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "password", resourceFullName, "password"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "photo", resourceFullName, "photo"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "preferred_language", resourceFullName, "preferred_language"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "primary_phone", resourceFullName, "primary_phone"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "timezone", resourceFullName, "timezone"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "title", resourceFullName, "title"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "type", resourceFullName, "type"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "verify_status", resourceFullName, "verify_status"),
				),
			},
		},
	})
}

func TestAccUserDataSource_NotFound(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccUserDataSourceConfig_NotFoundByName(resourceName),
				ExpectError: regexp.MustCompile("Cannot find user"),
			},
			{
				Config:      testAccUserDataSourceConfig_NotFoundByID(resourceName),
				ExpectError: regexp.MustCompile("Cannot find user"),
			},
		},
	})
}

func testAccUserDataSourceConfig_ByNameFull(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

data "pingone_user" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  username = "%[3]s"

  depends_on = [
    pingone_user.%[2]s,
  ]
}`, testAccUserConfig_Full(resourceName, name), resourceName, name)
}

func testAccUserDataSourceConfig_ByEmailFull(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

data "pingone_user" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  email = "%[3]s@pingidentity.com"

  depends_on = [
    pingone_user.%[2]s,
  ]
}`, testAccUserConfig_Full(resourceName, name), resourceName, name)
}

func testAccUserDataSourceConfig_ByIDFull(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_user" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  user_id = pingone_user.%[2]s.id
}`, testAccUserConfig_Full(resourceName, name), resourceName, name)
}

func testAccUserDataSourceConfig_NotFoundByName(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_user" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  username = "doesnotexist"
}`, acctest.GenericSandboxEnvironment(), resourceName)
}

func testAccUserDataSourceConfig_NotFoundByID(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_user" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  user_id = "9c052a8a-14be-44e4-8f07-2662569994ce" // dummy ID that conforms to UUID v4
}`, acctest.GenericSandboxEnvironment(), resourceName)
}
