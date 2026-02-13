// Copyright © 2026 Ping Identity Corporation

package sso_test

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
			//
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
					resource.TestCheckResourceAttrSet(dataSourceFullName, "hidden_from_app_portal"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "icon.id"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "icon.href"),
					resource.TestMatchResourceAttr(dataSourceFullName, "client_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "pkce_enforcement"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "token_endpoint_auth_method"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "apply_default_theme"),
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
					resource.TestCheckResourceAttrSet(dataSourceFullName, "hidden_from_app_portal"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "icon.id"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "icon.href"),
					resource.TestMatchResourceAttr(dataSourceFullName, "client_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "pkce_enforcement"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "token_endpoint_auth_method"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "apply_default_theme"),
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
					resource.TestCheckResourceAttrSet(dataSourceFullName, "hidden_from_app_portal"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "icon.id"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "icon.href"),
					resource.TestMatchResourceAttr(dataSourceFullName, "client_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "pkce_enforcement"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "token_endpoint_auth_method"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "apply_default_theme"),
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
					resource.TestCheckResourceAttrSet(dataSourceFullName, "hidden_from_app_portal"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "icon.id"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "icon.href"),
					resource.TestMatchResourceAttr(dataSourceFullName, "client_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "pkce_enforcement"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "token_endpoint_auth_method"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "apply_default_theme"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "enable_default_theme_footer"),
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
					resource.TestCheckResourceAttrSet(dataSourceFullName, "hidden_from_app_portal"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "icon.id"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "icon.href"),
					resource.TestMatchResourceAttr(dataSourceFullName, "client_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "pkce_enforcement"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "token_endpoint_auth_method"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "apply_default_theme"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "enable_default_theme_footer"),
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
					resource.TestCheckResourceAttrSet(dataSourceFullName, "hidden_from_app_portal"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "icon.id"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "icon.href"),
					resource.TestMatchResourceAttr(dataSourceFullName, "client_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "pkce_enforcement"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "token_endpoint_auth_method"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "apply_default_theme"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "enable_default_theme_footer"),
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

// Have to look up the ID first via the same data source
data "pingone_system_application" "portal_lookup" {
  environment_id = data.pingone_environment.general_test.id
  name           = "PingOne Application Portal"
}

data "pingone_system_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  application_id = data.pingone_system_application.portal_lookup.id
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

// Have to look up the ID first via the same data source
data "pingone_system_application" "selfservice_lookup" {
  environment_id = data.pingone_environment.general_test.id
  name           = "PingOne Self-Service - MyAccount"
}

data "pingone_system_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  application_id = data.pingone_system_application.selfservice_lookup.id
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
