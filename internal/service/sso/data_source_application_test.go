package sso_test

import (
	"encoding/base64"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccApplicationDataSource_OIDCAppByID(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullName := fmt.Sprintf("data.pingone_application.%s", resourceName)

	name := resourceName

	data, _ := os.ReadFile("../../acctest/test_assets/image/image-logo.gif")
	image := base64.StdEncoding.EncodeToString(data)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckApplicationDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationDataSource_OIDCAppByID(resourceName, name, image),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "name", name),
					resource.TestCheckResourceAttr(dataSourceFullName, "description", "My test OIDC app"),
					resource.TestCheckResourceAttr(dataSourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "tags.#", "0"),
					resource.TestCheckResourceAttr(dataSourceFullName, "login_page_url", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(dataSourceFullName, "icon.#", "1"),
					resource.TestMatchResourceAttr(dataSourceFullName, "icon.0.id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "icon.0.href", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
					resource.TestCheckResourceAttr(dataSourceFullName, "access_control_role_type", "ADMIN_USERS_ONLY"),
					resource.TestCheckResourceAttr(dataSourceFullName, "access_control_group_options.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourceFullName, "access_control_group_options.*", map[string]string{
						"type":     "ANY_GROUP",
						"groups.#": "1",
					}),
					resource.TestMatchTypeSetElemNestedAttrs(dataSourceFullName, "access_control_group_options.*", map[string]*regexp.Regexp{
						"groups.0": verify.P1ResourceIDRegexpFullString,
					}),
					resource.TestCheckResourceAttr(dataSourceFullName, "oidc_options.#", "1"),
					resource.TestCheckResourceAttr(dataSourceFullName, "oidc_options.0.type", "WEB_APP"),
					/*resource.TestCheckResourceAttr(dataSourceFullName, "oidc_options.0.home_page_url", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(dataSourceFullName, "oidc_options.0.initiate_login_uri", "https://www.pingidentity.com/initiate"),
					resource.TestCheckResourceAttr(dataSourceFullName, "oidc_options.0.target_link_uri", "https://www.pingidentity.com/target"),
					resource.TestCheckResourceAttr(dataSourceFullName, "oidc_options.0.grant_types.#", "2"),
					resource.TestCheckTypeSetElemAttr(dataSourceFullName, "oidc_options.0.grant_types.*", "AUTHORIZATION_CODE"),
					resource.TestCheckTypeSetElemAttr(dataSourceFullName, "oidc_options.0.grant_types.*", "REFRESH_TOKEN"),
					resource.TestCheckResourceAttr(dataSourceFullName, "oidc_options.0.response_types.#", "1"),
					resource.TestCheckTypeSetElemAttr(dataSourceFullName, "oidc_options.0.response_types.*", "CODE"),*/
					resource.TestCheckResourceAttr(dataSourceFullName, "oidc_options.0.token_endpoint_authn_method", "CLIENT_SECRET_BASIC"),
					/*resource.TestCheckResourceAttr(dataSourceFullName, "oidc_options.0.pkce_enforcement", "OPTIONAL"),
					resource.TestCheckResourceAttr(dataSourceFullName, "oidc_options.0.redirect_uris.#", "2"),
					resource.TestCheckTypeSetElemAttr(dataSourceFullName, "oidc_options.0.redirect_uris.*", "https://pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(dataSourceFullName, "oidc_options.0.redirect_uris.*", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(dataSourceFullName, "oidc_options.0.allow_wildcards_in_redirect_uris", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "oidc_options.0.post_logout_redirect_uris.#", "2"),
					resource.TestCheckTypeSetElemAttr(dataSourceFullName, "oidc_options.0.post_logout_redirect_uris.*", "https://pingidentity.com/logout"),
					resource.TestCheckTypeSetElemAttr(dataSourceFullName, "oidc_options.0.post_logout_redirect_uris.*", "https://www.pingidentity.com/logout"),
					resource.TestCheckResourceAttr(dataSourceFullName, "oidc_options.0.refresh_token_duration", "3000000"),
					resource.TestCheckResourceAttr(dataSourceFullName, "oidc_options.0.refresh_token_rolling_duration", "30000000"),
					resource.TestCheckResourceAttr(dataSourceFullName, "oidc_options.0.refresh_token_rolling_grace_period_duration", "80000"),
					resource.TestCheckResourceAttr(dataSourceFullName, "oidc_options.0.additional_refresh_token_replay_protection_enabled", "false"),
					resource.TestMatchResourceAttr(dataSourceFullName, "oidc_options.0.client_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "oidc_options.0.client_secret", regexp.MustCompile(`[a-zA-Z0-9-~_]{10,}`)),
					resource.TestCheckResourceAttr(dataSourceFullName, "oidc_options.0.support_unsigned_request_object", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "oidc_options.0.require_signed_request_object", "false"),
					resource.TestCheckResourceAttr(dataSourceFullName, "oidc_options.0.mobile_app.#", "0"),
					resource.TestCheckResourceAttr(dataSourceFullName, "oidc_options.0.bundle_id", ""),
					resource.TestCheckResourceAttr(dataSourceFullName, "oidc_options.0.package_name", ""),*/
					resource.TestCheckResourceAttr(dataSourceFullName, "saml_options.#", "0"),
					resource.TestCheckResourceAttr(dataSourceFullName, "external_link_options.#", "0"),
					resource.TestCheckResourceAttr(dataSourceFullName, "hidden_from_app_portal", "true"),
				),
			},
		},
	})
}

// additional tests to flesh out
// consider implementation to limit unnecessary duplication
// TestAccApplicationDataSource_OIDCAppByName
// TestAccApplicationDataSource_SAMLAppByID
// TestAccApplicationDataSource_SAMLAppByName
// TestAccApplicationDataSource_ExternalLinkAppByID
// TestAccApplicationDataSource_ExternalLinkAppByName
// TestAccApplicationDataSource_AppNotFoundByID
// TestAccApplicationDataSource_AppNotFoundByName

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
  tags           = []
  login_page_url = "https://www.pingidentity.com"

  icon {
    id   = pingone_image.%[2]s.id
    href = pingone_image.%[2]s.uploaded_image[0].href
  }

  access_control_role_type = "ADMIN_USERS_ONLY"

  access_control_group_options {
    type = "ANY_GROUP"

    groups = [
      pingone_group.%[2]s.id
    ]
  }

  hidden_from_app_portal = true


  enabled = true

  oidc_options {
    type                             = "WEB_APP"
    grant_types                      = ["REFRESH_TOKEN", "AUTHORIZATION_CODE"]
    response_types                   = ["CODE"]
    token_endpoint_authn_method      = "CLIENT_SECRET_BASIC"
    redirect_uris                    = ["https://www.pingidentity.com", "https://pingidentity.com"]
    allow_wildcards_in_redirect_uris = true
    post_logout_redirect_uris        = ["https://www.pingidentity.com/logout", "https://pingidentity.com/logout"]

    refresh_token_duration                             = 3000000
    refresh_token_rolling_duration                     = 30000000
    refresh_token_rolling_grace_period_duration        = 80000
    additional_refresh_token_replay_protection_enabled = false

    home_page_url      = "https://www.pingidentity.com"
    initiate_login_uri = "https://www.pingidentity.com/initiate"
    target_link_uri    = "https://www.pingidentity.com/target"
    pkce_enforcement   = "OPTIONAL"

    support_unsigned_request_object = true
  }
}
data "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  application_id = pingone_application.%[2]s.id
}`, acctest.GenericSandboxEnvironment(), resourceName, name, image)
}
