// Copyright © 2026 Ping Identity Corporation

package base_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccSystemApplicationDataSource_PingOnePortalByID(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullName := fmt.Sprintf("data.pingone_system_application.%s", resourceName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSystemApplicationDataSource_PingOnePortalByID(resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "type", "PING_ONE_PORTAL"),
					resource.TestCheckResourceAttr(dataSourceFullName, "protocol", "OPENID_CONNECT"),
					resource.TestCheckResourceAttr(dataSourceFullName, "name", "PingOne Application Portal"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "description"),
					resource.TestCheckResourceAttr(dataSourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "access_control_role_type", "ADMIN_USERS_ONLY"),
					resource.TestCheckResourceAttr(dataSourceFullName, "access_control_group_options.type", "ALL_GROUPS"),
					resource.TestCheckResourceAttr(dataSourceFullName, "access_control_group_options.groups.#", "2"),
					resource.TestCheckResourceAttr(dataSourceFullName, "hidden_from_app_portal", "false"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "icon.id"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "icon.href"),
					resource.TestMatchResourceAttr(dataSourceFullName, "client_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "pkce_enforcement", "OPTIONAL"),
					resource.TestCheckResourceAttr(dataSourceFullName, "token_endpoint_auth_method", "NONE"),
					resource.TestCheckResourceAttr(dataSourceFullName, "apply_default_theme", "true"),
				),
			},
		},
	})
}

func TestAccSystemApplicationDataSource_PingOnePortalByName(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullName := fmt.Sprintf("data.pingone_system_application.%s", resourceName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSystemApplicationDataSource_PingOnePortalByName(resourceName, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "type", "PING_ONE_PORTAL"),
					resource.TestCheckResourceAttr(dataSourceFullName, "protocol", "OPENID_CONNECT"),
					resource.TestCheckResourceAttr(dataSourceFullName, "name", "PingOne Application Portal"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "description"),
					resource.TestCheckResourceAttr(dataSourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "hidden_from_app_portal", "false"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "icon.id"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "icon.href"),
					resource.TestMatchResourceAttr(dataSourceFullName, "client_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "pkce_enforcement", "OPTIONAL"),
					resource.TestCheckResourceAttr(dataSourceFullName, "token_endpoint_auth_method", "NONE"),
				),
			},
			{
				Config: testAccSystemApplicationDataSource_PingOnePortalByName(resourceName, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "type", "PING_ONE_PORTAL"),
					resource.TestCheckResourceAttr(dataSourceFullName, "protocol", "OPENID_CONNECT"),
					resource.TestCheckResourceAttr(dataSourceFullName, "name", "PingOne Application Portal"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "description"),
					resource.TestCheckResourceAttr(dataSourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "hidden_from_app_portal", "false"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "icon.id"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "icon.href"),
					resource.TestMatchResourceAttr(dataSourceFullName, "client_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "pkce_enforcement", "OPTIONAL"),
					resource.TestCheckResourceAttr(dataSourceFullName, "token_endpoint_auth_method", "NONE"),
				),
			},
		},
	})
}

func TestAccSystemApplicationDataSource_PingOneSelfServiceByID(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullName := fmt.Sprintf("data.pingone_system_application.%s", resourceName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSystemApplicationDataSource_PingOneSelfServiceByID(resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "type", "PING_ONE_SELF_SERVICE"),
					resource.TestCheckResourceAttr(dataSourceFullName, "protocol", "OPENID_CONNECT"),
					resource.TestCheckResourceAttr(dataSourceFullName, "name", "PingOne Self-Service - MyAccount"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "description"),
					resource.TestCheckResourceAttr(dataSourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "access_control_role_type", "ADMIN_USERS_ONLY"),
					resource.TestCheckResourceAttr(dataSourceFullName, "access_control_group_options.type", "ALL_GROUPS"),
					resource.TestCheckResourceAttr(dataSourceFullName, "access_control_group_options.groups.#", "2"),
					resource.TestCheckResourceAttr(dataSourceFullName, "hidden_from_app_portal", "false"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "icon.id"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "icon.href"),
					resource.TestMatchResourceAttr(dataSourceFullName, "client_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "pkce_enforcement", "OPTIONAL"),
					resource.TestCheckResourceAttr(dataSourceFullName, "token_endpoint_auth_method", "NONE"),
					resource.TestCheckResourceAttr(dataSourceFullName, "apply_default_theme", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "enable_default_theme_footer", "true"),
				),
			},
		},
	})
}

func TestAccSystemApplicationDataSource_PingOneSelfServiceByName(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullName := fmt.Sprintf("data.pingone_system_application.%s", resourceName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSystemApplicationDataSource_PingOneSelfServiceByName(resourceName, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "type", "PING_ONE_SELF_SERVICE"),
					resource.TestCheckResourceAttr(dataSourceFullName, "protocol", "OPENID_CONNECT"),
					resource.TestCheckResourceAttr(dataSourceFullName, "name", "PingOne Self-Service - MyAccount"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "description"),
					resource.TestCheckResourceAttr(dataSourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "hidden_from_app_portal", "false"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "icon.id"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "icon.href"),
					resource.TestMatchResourceAttr(dataSourceFullName, "client_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "pkce_enforcement", "OPTIONAL"),
					resource.TestCheckResourceAttr(dataSourceFullName, "token_endpoint_auth_method", "NONE"),
				),
			},
			{
				Config: testAccSystemApplicationDataSource_PingOneSelfServiceByName(resourceName, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "type", "PING_ONE_SELF_SERVICE"),
					resource.TestCheckResourceAttr(dataSourceFullName, "protocol", "OPENID_CONNECT"),
					resource.TestCheckResourceAttr(dataSourceFullName, "name", "PingOne Self-Service - MyAccount"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "description"),
					resource.TestCheckResourceAttr(dataSourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "hidden_from_app_portal", "false"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "icon.id"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "icon.href"),
					resource.TestMatchResourceAttr(dataSourceFullName, "client_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "pkce_enforcement", "OPTIONAL"),
					resource.TestCheckResourceAttr(dataSourceFullName, "token_endpoint_auth_method", "NONE"),
				),
			},
		},
	})
}

func TestAccSystemApplicationDataSource_NotFound(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccSystemApplicationDataSource_NotFoundByID(resourceName),
				ExpectError: regexp.MustCompile("Error: Error when calling `ReadOneApplication`: Unable to find Application with ID"),
			},
			{
				Config:      testAccSystemApplicationDataSource_NotFoundByName(resourceName),
				ExpectError: regexp.MustCompile("Cannot find the system application from name"),
			},
		},
	})
}

func TestAccSystemApplicationDataSource_NotSystemApplication(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccSystemApplicationDataSource_NotSystemApplication(resourceName, name),
				ExpectError: regexp.MustCompile("Application is not a system application"),
			},
		},
	})
}

func testAccSystemApplicationDataSource_PingOnePortalByID(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_group" "%[2]s-1" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[2]s-1"
}

resource "pingone_group" "%[2]s-2" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[2]s-2"
}

resource "pingone_system_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  type    = "PING_ONE_PORTAL"
  enabled = true

  access_control_role_type = "ADMIN_USERS_ONLY"
  access_control_group_options = {
    groups = [
      pingone_group.%[2]s-2.id,
      pingone_group.%[2]s-1.id,
    ]

    type = "ALL_GROUPS"
  }

  apply_default_theme = true

}

data "pingone_system_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  application_id = pingone_system_application.%[2]s.id
}`, acctest.GenericSandboxEnvironment(), resourceName)
}

func testAccSystemApplicationDataSource_PingOnePortalByName(resourceName string, insensitivityCheck bool) string {
	// If insensitivityCheck is true, alter the case of the name
	nameComparator := "PingOne Application Portal"
	if insensitivityCheck {
		nameComparator = acctest.AlterStringCasing(nameComparator)
	}

	return fmt.Sprintf(`
%[1]s

data "pingone_system_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
}`, acctest.GenericSandboxEnvironment(), resourceName, nameComparator)
}

func testAccSystemApplicationDataSource_PingOneSelfServiceByID(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_group" "%[2]s-1" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[2]s-1"
}

resource "pingone_group" "%[2]s-2" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[2]s-2"
}

resource "pingone_system_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  type    = "PING_ONE_SELF_SERVICE"
  enabled = true

  access_control_role_type = "ADMIN_USERS_ONLY"
  access_control_group_options = {
    groups = [
      pingone_group.%[2]s-2.id,
      pingone_group.%[2]s-1.id,
    ]

    type = "ALL_GROUPS"
  }

  apply_default_theme         = true
  enable_default_theme_footer = true

}

data "pingone_system_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  application_id = pingone_system_application.%[2]s.id
}`, acctest.GenericSandboxEnvironment(), resourceName)
}

func testAccSystemApplicationDataSource_PingOneSelfServiceByName(resourceName string, insensitivityCheck bool) string {
	// If insensitivityCheck is true, alter the case of the name
	nameComparator := "PingOne Self-Service - MyAccount"
	if insensitivityCheck {
		nameComparator = acctest.AlterStringCasing(nameComparator)
	}

	return fmt.Sprintf(`
%[1]s

data "pingone_system_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
}`, acctest.GenericSandboxEnvironment(), resourceName, nameComparator)
}

func testAccSystemApplicationDataSource_NotFoundByID(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

data "pingone_system_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  application_id = "9c052a8a-14be-44e4-8f07-2662569994ce" // dummy ID that conforms to UUID v4
}`, acctest.GenericSandboxEnvironment(), resourceName)
}

func testAccSystemApplicationDataSource_NotFoundByName(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

data "pingone_system_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "NonExistentSystemApplication"
}`, acctest.GenericSandboxEnvironment(), resourceName)
}

func testAccSystemApplicationDataSource_NotSystemApplication(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true

  oidc_options = {
    type                       = "SINGLE_PAGE_APP"
    grant_types                = ["AUTHORIZATION_CODE"]
    response_types             = ["CODE"]
    token_endpoint_auth_method = "NONE"
    redirect_uris              = ["https://www.example.com"]
  }
}

data "pingone_system_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  application_id = pingone_application.%[2]s.id
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
