// Copyright © 2025 Ping Identity Corporation

package sso_test

import (
	"encoding/base64"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/sso"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccApplicationDataSource_OIDCAppByID(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	name := resourceName

	data, _ := os.ReadFile("../../acctest/test_assets/image/image-logo.gif")
	image := base64.StdEncoding.EncodeToString(data)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Application_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationDataSource_OIDCAppByID(resourceName, name, image),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "environment_id", resourceFullName, "environment_id"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "name", resourceFullName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "description", resourceFullName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "enabled", resourceFullName, "enabled"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "tags", resourceFullName, "tags"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "login_page_url", resourceFullName, "login_page_url"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "icon", resourceFullName, "icon"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "access_control_role_type", resourceFullName, "access_control_role_type"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "access_control_group_options", resourceFullName, "access_control_group_options"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "oidc_options", resourceFullName, "oidc_options"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "saml_options", resourceFullName, "saml_options"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "external_link_options", resourceFullName, "external_link_options"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "hidden_from_app_portal", resourceFullName, "hidden_from_app_portal"),
				),
			},
		},
	})
}

func TestAccApplicationDataSource_OIDCAppByName(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	name := resourceName

	data, _ := os.ReadFile("../../acctest/test_assets/image/image-logo.gif")
	image := base64.StdEncoding.EncodeToString(data)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Application_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationDataSource_OIDCAppByName(resourceName, name, image),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "environment_id", resourceFullName, "environment_id"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "name", resourceFullName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "description", resourceFullName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "enabled", resourceFullName, "enabled"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "tags", resourceFullName, "tags"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "login_page_url", resourceFullName, "login_page_url"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "icon", resourceFullName, "icon"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "access_control_role_type", resourceFullName, "access_control_role_type"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "access_control_group_options", resourceFullName, "access_control_group_options"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "oidc_options", resourceFullName, "oidc_options"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "saml_options", resourceFullName, "saml_options"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "external_link_options", resourceFullName, "external_link_options"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "hidden_from_app_portal", resourceFullName, "hidden_from_app_portal"),
				),
			},
		},
	})
}

func TestAccApplicationDataSource_ExternalLinkAppByID(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	name := resourceName

	data, _ := os.ReadFile("../../acctest/test_assets/image/image-logo.gif")
	image := base64.StdEncoding.EncodeToString(data)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Application_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationDataSource_ExternalLinkAppByID(resourceName, name, image),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "environment_id", resourceFullName, "environment_id"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "name", resourceFullName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "description", resourceFullName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "enabled", resourceFullName, "enabled"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "tags", resourceFullName, "tags"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "login_page_url", resourceFullName, "login_page_url"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "icon", resourceFullName, "icon"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "access_control_role_type", resourceFullName, "access_control_role_type"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "access_control_group_options", resourceFullName, "access_control_group_options"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "oidc_options", resourceFullName, "oidc_options"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "saml_options", resourceFullName, "saml_options"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "external_link_options", resourceFullName, "external_link_options"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "hidden_from_app_portal", resourceFullName, "hidden_from_app_portal"),
				),
			},
		},
	})
}

func TestAccApplicationDataSource_ExternalLinkAppByName(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Application_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationDataSource_ExternalLinkAppByName(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "environment_id", resourceFullName, "environment_id"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "name", resourceFullName, "name"),
					//resource.TestCheckResourceAttrPair(dataSourceFullName, "description", resourceFullName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "enabled", resourceFullName, "enabled"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "tags", resourceFullName, "tags"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "login_page_url", resourceFullName, "login_page_url"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "icon", resourceFullName, "icon"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "access_control_role_type", resourceFullName, "access_control_role_type"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "access_control_group_options", resourceFullName, "access_control_group_options"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "oidc_options", resourceFullName, "oidc_options"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "saml_options", resourceFullName, "saml_options"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "external_link_options", resourceFullName, "external_link_options"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "hidden_from_app_portal", resourceFullName, "hidden_from_app_portal"),
				),
			},
		},
	})
}

func TestAccApplicationDataSource_SAMLAppByID(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	data, _ := os.ReadFile("../../acctest/test_assets/image/image-logo.gif")
	image := base64.StdEncoding.EncodeToString(data)

	pem_cert := os.Getenv("PINGONE_KEY_PEM_CERT")
	pkcs7_cert := os.Getenv("PINGONE_KEY_PKCS7_CERT")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
			acctest.PreCheckNoBeta(t)
			acctest.PreCheckPKCS7Cert(t)
			acctest.PreCheckPEMCert(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Application_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationDataSource_SAMLAppByID(environmentName, licenseID, resourceName, name, image, pkcs7_cert, pem_cert),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "environment_id", resourceFullName, "environment_id"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "name", resourceFullName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "description", resourceFullName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "enabled", resourceFullName, "enabled"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "tags", resourceFullName, "tags"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "login_page_url", resourceFullName, "login_page_url"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "icon", resourceFullName, "icon"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "access_control_role_type", resourceFullName, "access_control_role_type"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "access_control_group_options", resourceFullName, "access_control_group_options"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "oidc_options", resourceFullName, "oidc_options"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "saml_options", resourceFullName, "saml_options"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "external_link_options", resourceFullName, "external_link_options"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "hidden_from_app_portal", resourceFullName, "hidden_from_app_portal"),
				),
			},
		},
	})
}

func TestAccApplicationDataSource_SAMLAppByName(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Application_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationDataSource_SAMLAppByName(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "environment_id", resourceFullName, "environment_id"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "name", resourceFullName, "name"),
					//resource.TestCheckResourceAttrPair(dataSourceFullName, "description", resourceFullName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "enabled", resourceFullName, "enabled"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "tags", resourceFullName, "tags"),
					//resource.TestCheckResourceAttrPair(dataSourceFullName, "login_page_url", resourceFullName, "login_page_url"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "icon", resourceFullName, "icon"),
					//resource.TestCheckResourceAttrPair(dataSourceFullName, "access_control_role_type", resourceFullName, "access_control_role_type"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "access_control_group_options", resourceFullName, "access_control_group_options"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "oidc_options", resourceFullName, "oidc_options"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "saml_options", resourceFullName, "saml_options"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "external_link_options", resourceFullName, "external_link_options"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "hidden_from_app_portal", resourceFullName, "hidden_from_app_portal"),
				),
			},
		},
	})
}

func TestAccApplicationDataSource_WSFedAppByID(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	data, _ := os.ReadFile("../../acctest/test_assets/image/image-logo.gif")
	image := base64.StdEncoding.EncodeToString(data)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Application_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationDataSource_WSFedAppByID(environmentName, licenseID, resourceName, name, image),
				Check:  testAccApplicationConfig_WSFed_FullCheck(dataSourceFullName, name),
			},
		},
	})
}

func TestAccApplicationDataSource_WSFedAppByName(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Application_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationDataSource_WSFedAppByName(resourceName, name),
				Check:  testAccApplicationConfig_WSFed_MinimalCheck(dataSourceFullName, name),
			},
		},
	})
}

func TestAccApplicationDataSource_FailureChecks(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Application_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccApplicationDataSource_FindByIDFail(resourceName, name),
				ExpectError: regexp.MustCompile("Error: Error when calling `ReadOneApplication`: Unable to find Application with ID"),
			},
			{
				Config:      testAccApplicationDataSource_FindByNameFail(resourceName, name),
				ExpectError: regexp.MustCompile("Error: Cannot find the application from name"),
			},
		},
	})
}

func testAccApplicationDataSource_OIDCAppByID(resourceName, name, image string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_group" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
}

resource "pingone_image" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  image_file_base64 = "%[4]s"
}

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "My test OIDC app"

  login_page_url = "https://www.pingidentity.com"

  icon = {
    id   = pingone_image.%[2]s.id
    href = pingone_image.%[2]s.uploaded_image.href
  }

  access_control_role_type = "ADMIN_USERS_ONLY"

  access_control_group_options = {
    type = "ANY_GROUP"

    groups = [
      pingone_group.%[2]s.id
    ]
  }

  hidden_from_app_portal = true


  enabled = true

  oidc_options = {
    type                            = "WEB_APP"
    grant_types                     = ["REFRESH_TOKEN", "AUTHORIZATION_CODE"]
    response_types                  = ["CODE"]
    token_endpoint_auth_method      = "CLIENT_SECRET_BASIC"
    redirect_uris                   = ["https://www.pingidentity.com", "https://pingidentity.com"]
    allow_wildcard_in_redirect_uris = true
    post_logout_redirect_uris       = ["https://www.pingidentity.com/logout", "https://pingidentity.com/logout"]

    refresh_token_duration                             = 3000000
    refresh_token_rolling_duration                     = 30000000
    refresh_token_rolling_grace_period_duration        = 80000
    additional_refresh_token_replay_protection_enabled = false
    idp_signoff                                        = true

    home_page_url      = "https://www.pingidentity.com"
    initiate_login_uri = "https://www.pingidentity.com/initiate"
    target_link_uri    = "https://www.pingidentity.com/target"
    pkce_enforcement   = "OPTIONAL"

    support_unsigned_request_object = true

    cors_settings = {
      behavior = "ALLOW_SPECIFIC_ORIGINS"
      origins = [
        "http://localhost",
        "https://localhost",
        "http://auth.pingidentity.com",
        "https://auth.pingidentity.com",
        "http://*.pingidentity.com",
        "https://*.pingidentity.com",
        "http://192.168.1.1",
        "https://192.168.1.1",
      ]
    }
  }
}
data "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  application_id = pingone_application.%[2]s.id
}`, acctest.GenericSandboxEnvironment(), resourceName, name, image)
}

func testAccApplicationDataSource_OIDCAppByName(resourceName, name, image string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_group" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
}

resource "pingone_image" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  image_file_base64 = "%[4]s"
}

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "My test OIDC app"

  login_page_url = "https://www.pingidentity.com"

  icon = {
    id   = pingone_image.%[2]s.id
    href = pingone_image.%[2]s.uploaded_image.href
  }

  access_control_role_type = "ADMIN_USERS_ONLY"

  access_control_group_options = {
    type = "ANY_GROUP"

    groups = [
      pingone_group.%[2]s.id
    ]
  }

  hidden_from_app_portal = true

  enabled = true

  oidc_options = {
    type                       = "NATIVE_APP"
    grant_types                = ["CLIENT_CREDENTIALS"]
    token_endpoint_auth_method = "CLIENT_SECRET_BASIC"

    allow_wildcard_in_redirect_uris = true

    mobile_app = {
      bundle_id           = "com.%[2]s.bundle"
      package_name        = "com.%[2]s.package"
      huawei_app_id       = "%[2]s"
      huawei_package_name = "com.%[2]s.huaweipackage"

      passcode_refresh_seconds = 45

      universal_app_link = "https://applink.com"

      integrity_detection = {
        enabled = true

        excluded_platforms = ["IOS"]

        cache_duration = {
          amount = 30
          units  = "HOURS"
        }

        google_play = {
          verification_type = "INTERNAL"
          decryption_key    = "decryptionkeydoesnotexist"
          verification_key  = "verificationkeydoesnotexist"
        }
      }
    }
  }
}
data "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  depends_on = [pingone_application.%[2]s]
}`, acctest.GenericSandboxEnvironment(), resourceName, name, image)
}

func testAccApplicationDataSource_ExternalLinkAppByID(resourceName, name, image string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_group" "%[2]s-1" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-1"
}

resource "pingone_group" "%[2]s-2" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-2"
}

resource "pingone_image" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  image_file_base64 = "%[4]s"
}

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "My test external link app"

  icon = {
    id   = pingone_image.%[2]s.id
    href = pingone_image.%[2]s.uploaded_image.href
  }

  access_control_group_options = {
    type = "ANY_GROUP"

    groups = [
      pingone_group.%[2]s-2.id,
      pingone_group.%[2]s-1.id
    ]
  }

  hidden_from_app_portal = true

  enabled = true

  external_link_options = {
    home_page_url = "https://www.pingidentity.com"
  }
}
data "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  application_id = pingone_application.%[2]s.id
}`, acctest.GenericSandboxEnvironment(), resourceName, name, image)
}

func testAccApplicationDataSource_ExternalLinkAppByName(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true

  external_link_options = {
    home_page_url = "https://www.pingidentity.com"
  }
}
data "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  depends_on = [pingone_application.%[2]s]
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccApplicationDataSource_FindByIDFail(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  application_id = "9c052a8a-14be-44e4-8f07-2662569994ce" // dummy ID that conforms to UUID v4


}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccApplicationDataSource_FindByNameFail(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccApplicationDataSource_SAMLAppByID(environmentName, licenseID, resourceName, name, image, pkcs7_cert, pem_cert string) string {
	return fmt.Sprintf(`
		%[1]s


resource "pingone_group" "%[3]s-1" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s-1"
}

resource "pingone_group" "%[3]s-2" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s-2"
}

resource "pingone_key" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  name                = "%[4]s"
  algorithm           = "EC"
  key_length          = 256
  signature_algorithm = "SHA384withECDSA"
  subject_dn          = "CN=%[4]s, OU=Ping Identity, O=Ping Identity, L=, ST=, C=US"
  usage_type          = "SIGNING"
  validity_period     = 365
}

resource "pingone_image" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  image_file_base64 = "%[5]s"
}

resource "pingone_certificate" "%[3]s-1" {
  environment_id = pingone_environment.%[2]s.id

  pkcs7_file_base64 = <<EOT
%[6]s
EOT

  usage_type = "SIGNING"
}

resource "pingone_certificate" "%[3]s-2" {
  environment_id = pingone_environment.%[2]s.id

  pem_file = <<EOT
%[7]s
EOT

  usage_type = "SIGNING"
}

resource "pingone_certificate" "%[3]s-enc" {
  environment_id = pingone_environment.%[2]s.id
  pem_file       = <<EOT
%[7]s
EOT
  usage_type     = "ENCRYPTION"
}

resource "pingone_application" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s"
  description    = "My test SAML app"
  login_page_url = "https://www.pingidentity.com"

  icon = {
    id   = pingone_image.%[3]s.id
    href = pingone_image.%[3]s.uploaded_image.href
  }

  access_control_role_type = "ADMIN_USERS_ONLY"

  access_control_group_options = {
    type = "ANY_GROUP"

    groups = [
      pingone_group.%[3]s-2.id,
      pingone_group.%[3]s-1.id
    ]
  }

  hidden_from_app_portal = true

  enabled = true

  saml_options = {
    type               = "WEB_APP"
    home_page_url      = "https://www.pingidentity.com"
    acs_urls           = ["https://www.pingidentity.com", "https://pingidentity.com"]
    assertion_duration = 3600
    sp_entity_id       = "sp:entity:%[3]s"

    sp_encryption = {
      algorithm = "AES_256"
      certificate = {
        id = pingone_certificate.%[3]s-enc.id
      }
    }

    assertion_signed_enabled = false

    enable_requested_authn_context   = true
    nameid_format                    = "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress"
    response_is_signed               = true
    session_not_on_or_after_duration = 64
    slo_binding                      = "HTTP_REDIRECT"
    slo_endpoint                     = "https://www.pingidentity.com/sloendpoint"
    slo_response_endpoint            = "https://www.pingidentity.com/sloresponseendpoint"
    slo_window                       = 3

    default_target_url = "https://www.pingidentity.com/relaystate"

    idp_signing_key = {
      algorithm = pingone_key.%[3]s.signature_algorithm
      key_id    = pingone_key.%[3]s.id
    }

    sp_verification = {
      authn_request_signed = true
      certificate_ids = [
        pingone_certificate.%[3]s-2.id,
        pingone_certificate.%[3]s-1.id,
      ]
    }

    cors_settings = {
      behavior = "ALLOW_SPECIFIC_ORIGINS"
      origins = [
        "http://localhost",
        "https://localhost",
        "http://auth.pingidentity.com",
        "https://auth.pingidentity.com",
        "http://*.pingidentity.com",
        "https://*.pingidentity.com",
        "http://192.168.1.1",
        "https://192.168.1.1",
      ]
    }

    virtual_server_id_settings = {
      enabled = true
      virtual_server_ids = [
        {
          vs_id = "virtualserver1"
        },
        {
          vs_id   = "virtualserver2"
          default = false
        },
        {
          vs_id = "virtualserver4"
        },
        {
          vs_id   = "virtualserver3"
          default = true
        }
      ]
    }
  }
}
data "pingone_application" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  application_id = pingone_application.%[3]s.id
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name, image, pkcs7_cert, pem_cert)
}

func testAccApplicationDataSource_SAMLAppByName(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_key" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name                = "%[3]s"
  algorithm           = "EC"
  key_length          = 256
  signature_algorithm = "SHA384withECDSA"
  subject_dn          = "CN=%[3]s, OU=Ping Identity, O=Ping Identity, L=, ST=, C=US"
  usage_type          = "SIGNING"
  validity_period     = 365
}

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true

  saml_options = {
    acs_urls           = ["https://pingidentity.com"]
    assertion_duration = 3600
    sp_entity_id       = "sp:entity:%[2]s"
    idp_signing_key = {
      key_id    = pingone_key.%[2]s.id
      algorithm = pingone_key.%[2]s.signature_algorithm
    }
  }
}
data "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  depends_on = [pingone_application.%[2]s]
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccApplicationDataSource_WSFedAppByID(environmentName, licenseID, resourceName, name, image string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_group" "%[3]s-1" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s-1"
}

resource "pingone_group" "%[3]s-2" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s-2"
}

resource "pingone_key" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  name                = "%[4]s"
  algorithm           = "EC"
  key_length          = 256
  signature_algorithm = "SHA384withECDSA"
  subject_dn          = "CN=%[4]s, OU=Ping Identity, O=Ping Identity, L=, ST=, C=US"
  usage_type          = "SIGNING"
  validity_period     = 365
}

resource "pingone_image" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  image_file_base64 = "%[5]s"
}

resource "pingone_population" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  name = "%[4]s"
}

resource "pingone_gateway" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s"
  enabled        = true
  type           = "LDAP"
  description    = "My test gateway"

  bind_dn       = "ou=test1,dc=example,dc=com"
  bind_password = "dummyPasswordValue1"

  connection_security = "TLS"
  vendor              = "Microsoft Active Directory"

  kerberos = {
    service_account_upn              = "username@domainname"
    service_account_password         = "dummyKerberosPasswordValue"
    retain_previous_credentials_mins = 20
  }

  servers = [
    "ds1.dummyldapservice.com:636",
    "ds3.dummyldapservice.com:636",
    "ds2.dummyldapservice.com:636",
  ]

  validate_tls_certificates = false

  user_types = {
    "User Set 1" = {
      password_authority = "LDAP"
      search_base_dn     = "ou=users1,dc=example,dc=com"

      user_link_attributes = ["objectGUID", "objectSid"]
    },
    "User Set 2" = {
      password_authority = "PING_ONE"
      search_base_dn     = "ou=users,dc=example,dc=com"

      user_link_attributes = ["objectGUID", "dn", "objectSid"]

      new_user_lookup = {
        ldap_filter_pattern = "(|(uid=$${identifier})(mail=$${identifier}))"

        population_id = pingone_population.%[3]s.id

        attribute_mappings = [
          {
            name  = "username"
            value = "$${ldapAttributes.uid}"
          },
          {
            name  = "email"
            value = "$${ldapAttributes.mail}"
          },
          {
            name  = "name.family"
            value = "$${ldapAttributes.sn}"
          }
        ]
      }
    }
  }
}

resource "pingone_application" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s"
  description    = "My test WS-Fed app"
  login_page_url = "https://www.pingidentity.com"

  icon = {
    id   = pingone_image.%[3]s.id
    href = pingone_image.%[3]s.uploaded_image.href
  }

  access_control_role_type = "ADMIN_USERS_ONLY"

  access_control_group_options = {
    type = "ANY_GROUP"

    groups = [
      pingone_group.%[3]s-2.id,
      pingone_group.%[3]s-1.id
    ]
  }

  hidden_from_app_portal = true

  enabled = true

  wsfed_options = {
    audience_restriction = "urn:federation:Example"
    cors_settings = {
      behavior = "ALLOW_SPECIFIC_ORIGINS"
      origins = [
        "http://localhost",
        "https://localhost",
        "http://auth.pingidentity.com",
        "https://auth.pingidentity.com",
        "http://*.pingidentity.com",
        "https://*.pingidentity.com",
        "http://192.168.1.1",
        "https://192.168.1.1",
      ]
    }
    domain_name = "my.updated.domain.name.example.com"
    idp_signing_key = {
      key_id    = pingone_key.%[3]s.id
      algorithm = pingone_key.%[3]s.signature_algorithm
    }
    kerberos = {
      gateways = [
        {
          id   = pingone_gateway.%[3]s.id
          type = "LDAP"
          user_type = {
            id = pingone_gateway.%[3]s.user_types["User Set 2"].id
          }
        }
      ]
    }
    reply_url                      = "https://example.com"
    slo_endpoint                   = "https://example.com/slo"
    subject_name_identifier_format = "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress"
    type                           = "WEB_APP"
  }
}
data "pingone_application" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  application_id = pingone_application.%[3]s.id
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name, image)
}

func testAccApplicationDataSource_WSFedAppByName(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_key" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name                = "%[3]s"
  algorithm           = "EC"
  key_length          = 256
  signature_algorithm = "SHA384withECDSA"
  subject_dn          = "CN=%[3]s, OU=Ping Identity, O=Ping Identity, L=, ST=, C=US"
  usage_type          = "SIGNING"
  validity_period     = 365
}

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true

  wsfed_options = {
    domain_name = "my.domain.name.example.com"
    idp_signing_key = {
      key_id    = pingone_key.%[2]s.id
      algorithm = pingone_key.%[2]s.signature_algorithm
    }
    reply_url = "https://example.com"
    type      = "WEB_APP"
  }
}
data "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  depends_on = [pingone_application.%[2]s]
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
