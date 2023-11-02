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
	dataSourceFullName := fmt.Sprintf("data.pingone_application.%s", resourceName)

	name := resourceName

	data, _ := os.ReadFile("../../acctest/test_assets/image/image-logo.gif")
	image := base64.StdEncoding.EncodeToString(data)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Application_CheckDestroy,
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
					resource.TestCheckResourceAttr(dataSourceFullName, "oidc_options.0.home_page_url", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(dataSourceFullName, "oidc_options.0.initiate_login_uri", "https://www.pingidentity.com/initiate"),
					resource.TestCheckResourceAttr(dataSourceFullName, "oidc_options.0.target_link_uri", "https://www.pingidentity.com/target"),
					resource.TestCheckResourceAttr(dataSourceFullName, "oidc_options.0.grant_types.#", "2"),
					resource.TestCheckTypeSetElemAttr(dataSourceFullName, "oidc_options.0.grant_types.*", "AUTHORIZATION_CODE"),
					resource.TestCheckTypeSetElemAttr(dataSourceFullName, "oidc_options.0.grant_types.*", "REFRESH_TOKEN"),
					resource.TestCheckResourceAttr(dataSourceFullName, "oidc_options.0.response_types.#", "1"),
					resource.TestCheckTypeSetElemAttr(dataSourceFullName, "oidc_options.0.response_types.*", "CODE"),
					resource.TestCheckResourceAttr(dataSourceFullName, "oidc_options.0.token_endpoint_authn_method", "CLIENT_SECRET_BASIC"),
					resource.TestCheckResourceAttr(dataSourceFullName, "oidc_options.0.pkce_enforcement", "OPTIONAL"),
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
					resource.TestCheckResourceAttr(dataSourceFullName, "saml_options.#", "0"),
					resource.TestCheckResourceAttr(dataSourceFullName, "external_link_options.#", "0"),
					resource.TestCheckResourceAttr(dataSourceFullName, "hidden_from_app_portal", "true"),
				),
			},
		},
	})
}

func TestAccApplicationDataSource_OIDCAppByName(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullName := fmt.Sprintf("data.pingone_application.%s", resourceName)

	name := resourceName

	data, _ := os.ReadFile("../../acctest/test_assets/image/image-logo.gif")
	image := base64.StdEncoding.EncodeToString(data)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Application_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationDataSource_OIDCAppByName(resourceName, name, image),
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
					resource.TestCheckResourceAttr(dataSourceFullName, "oidc_options.0.type", "NATIVE_APP"),
					resource.TestCheckNoResourceAttr(dataSourceFullName, "oidc_options.0.home_page_url"),
					resource.TestCheckNoResourceAttr(dataSourceFullName, "oidc_options.0.initiate_login_uri"),
					resource.TestCheckNoResourceAttr(dataSourceFullName, "oidc_options.0.target_link_uri"),
					resource.TestCheckResourceAttr(dataSourceFullName, "oidc_options.0.grant_types.#", "1"),
					resource.TestCheckTypeSetElemAttr(dataSourceFullName, "oidc_options.0.grant_types.*", "CLIENT_CREDENTIALS"),
					resource.TestCheckResourceAttr(dataSourceFullName, "oidc_options.0.response_types.#", "0"),
					resource.TestCheckResourceAttr(dataSourceFullName, "oidc_options.0.token_endpoint_authn_method", "CLIENT_SECRET_BASIC"),
					resource.TestCheckResourceAttr(dataSourceFullName, "oidc_options.0.pkce_enforcement", "OPTIONAL"),
					resource.TestCheckResourceAttr(dataSourceFullName, "oidc_options.0.redirect_uris.#", "0"),
					resource.TestCheckResourceAttr(dataSourceFullName, "oidc_options.0.allow_wildcards_in_redirect_uris", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "oidc_options.0.post_logout_redirect_uris.#", "0"),
					resource.TestCheckResourceAttr(dataSourceFullName, "oidc_options.0.refresh_token_duration", "2592000"),
					resource.TestCheckResourceAttr(dataSourceFullName, "oidc_options.0.refresh_token_rolling_duration", "15552000"),
					resource.TestCheckResourceAttr(dataSourceFullName, "oidc_options.0.refresh_token_rolling_grace_period_duration", "0"),
					resource.TestCheckResourceAttr(dataSourceFullName, "oidc_options.0.additional_refresh_token_replay_protection_enabled", "true"),
					resource.TestMatchResourceAttr(dataSourceFullName, "oidc_options.0.client_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "oidc_options.0.client_secret", regexp.MustCompile(`[a-zA-Z0-9-~_]{10,}`)),
					resource.TestCheckResourceAttr(dataSourceFullName, "oidc_options.0.support_unsigned_request_object", "false"),
					resource.TestCheckResourceAttr(dataSourceFullName, "oidc_options.0.require_signed_request_object", "false"),
					resource.TestCheckResourceAttr(dataSourceFullName, "oidc_options.0.mobile_app.#", "1"),
					resource.TestCheckResourceAttr(dataSourceFullName, "oidc_options.0.mobile_app.0.bundle_id", fmt.Sprintf("com.%s.bundle", resourceName)),
					resource.TestCheckResourceAttr(dataSourceFullName, "oidc_options.0.mobile_app.0.package_name", fmt.Sprintf("com.%s.package", resourceName)),
					resource.TestCheckResourceAttr(dataSourceFullName, "oidc_options.0.mobile_app.0.huawei_app_id", resourceName),
					resource.TestCheckResourceAttr(dataSourceFullName, "oidc_options.0.mobile_app.0.huawei_package_name", fmt.Sprintf("com.%s.huaweipackage", resourceName)),
					resource.TestCheckResourceAttr(dataSourceFullName, "oidc_options.0.mobile_app.0.passcode_refresh_seconds", "45"),
					resource.TestCheckResourceAttr(dataSourceFullName, "oidc_options.0.mobile_app.0.universal_app_link", "https://applink.com"),
					resource.TestCheckResourceAttr(dataSourceFullName, "oidc_options.0.mobile_app.0.integrity_detection.#", "1"),
					resource.TestCheckResourceAttr(dataSourceFullName, "oidc_options.0.mobile_app.0.integrity_detection.0.enabled", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "oidc_options.0.mobile_app.0.integrity_detection.0.excluded_platforms.0", "IOS"),
					resource.TestCheckResourceAttr(dataSourceFullName, "oidc_options.0.mobile_app.0.integrity_detection.0.cache_duration.#", "1"),
					resource.TestCheckResourceAttr(dataSourceFullName, "oidc_options.0.mobile_app.0.integrity_detection.0.cache_duration.0.amount", "30"),
					resource.TestCheckResourceAttr(dataSourceFullName, "oidc_options.0.mobile_app.0.integrity_detection.0.cache_duration.0.units", "HOURS"),
					resource.TestCheckResourceAttr(dataSourceFullName, "oidc_options.0.mobile_app.0.integrity_detection.0.google_play.#", "1"),
					resource.TestCheckResourceAttr(dataSourceFullName, "oidc_options.0.mobile_app.0.integrity_detection.0.google_play.0.verification_type", "INTERNAL"),
					resource.TestCheckResourceAttr(dataSourceFullName, "oidc_options.0.mobile_app.0.integrity_detection.0.google_play.0.verification_key", "DUMMY_SUPPRESS_VALUE"),
					resource.TestCheckResourceAttr(dataSourceFullName, "oidc_options.0.mobile_app.0.integrity_detection.0.google_play.0.decryption_key", "DUMMY_SUPPRESS_VALUE"),
					resource.TestCheckNoResourceAttr(dataSourceFullName, "oidc_options.0.mobile_app.0.integrity_detection.0.google_play.0.service_account_credentials_json"),
					resource.TestCheckResourceAttr(dataSourceFullName, "saml_options.#", "0"),
					resource.TestCheckResourceAttr(dataSourceFullName, "external_link_options.#", "0"),
					resource.TestCheckResourceAttr(dataSourceFullName, "hidden_from_app_portal", "true"),
				),
			},
		},
	})
}

func TestAccApplicationDataSource_ExternalLinkAppByID(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullName := fmt.Sprintf("data.pingone_application.%s", resourceName)

	name := resourceName

	data, _ := os.ReadFile("../../acctest/test_assets/image/image-logo.gif")
	image := base64.StdEncoding.EncodeToString(data)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Application_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationDataSource_ExternalLinkAppByID(resourceName, name, image),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "name", name),
					resource.TestCheckResourceAttr(dataSourceFullName, "description", "My test external link app"),
					resource.TestCheckResourceAttr(dataSourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "icon.#", "1"),
					resource.TestMatchResourceAttr(dataSourceFullName, "icon.0.id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "icon.0.href", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
					resource.TestCheckResourceAttr(dataSourceFullName, "access_control_group_options.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourceFullName, "access_control_group_options.*", map[string]string{
						"type":     "ANY_GROUP",
						"groups.#": "2",
					}),
					resource.TestMatchTypeSetElemNestedAttrs(dataSourceFullName, "access_control_group_options.*", map[string]*regexp.Regexp{
						"groups.0": verify.P1ResourceIDRegexpFullString,
						"groups.1": verify.P1ResourceIDRegexpFullString,
					}),
					resource.TestCheckResourceAttr(dataSourceFullName, "oidc_options.#", "0"),
					resource.TestCheckResourceAttr(dataSourceFullName, "saml_options.#", "0"),
					resource.TestCheckResourceAttr(dataSourceFullName, "external_link_options.#", "1"),
					resource.TestCheckResourceAttr(dataSourceFullName, "external_link_options.0.home_page_url", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(dataSourceFullName, "hidden_from_app_portal", "true"),
				),
			},
		},
	})
}

func TestAccApplicationDataSource_ExternalLinkAppByName(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullName := fmt.Sprintf("data.pingone_application.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Application_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationDataSource_ExternalLinkAppByName(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "name", name),
					resource.TestCheckNoResourceAttr(dataSourceFullName, "description"),
					resource.TestCheckResourceAttr(dataSourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "icon.#", "0"),
					resource.TestCheckResourceAttr(dataSourceFullName, "oidc_options.#", "0"),
					resource.TestCheckResourceAttr(dataSourceFullName, "saml_options.#", "0"),
					resource.TestCheckResourceAttr(dataSourceFullName, "external_link_options.#", "1"),
					resource.TestCheckResourceAttr(dataSourceFullName, "external_link_options.0.home_page_url", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(dataSourceFullName, "hidden_from_app_portal", "false"),
				),
			},
		},
	})
}

func TestAccApplicationDataSource_SAMLAppByID(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullName := fmt.Sprintf("data.pingone_application.%s", resourceName)

	name := resourceName

	data, _ := os.ReadFile("../../acctest/test_assets/image/image-logo.gif")
	image := base64.StdEncoding.EncodeToString(data)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Application_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationDataSource_SAMLAppByID(resourceName, name, image),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "name", name),
					resource.TestCheckResourceAttr(dataSourceFullName, "description", "My test SAML app"),
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
						"groups.#": "2",
					}),
					resource.TestMatchTypeSetElemNestedAttrs(dataSourceFullName, "access_control_group_options.*", map[string]*regexp.Regexp{
						"groups.0": verify.P1ResourceIDRegexpFullString,
						"groups.1": verify.P1ResourceIDRegexpFullString,
					}),
					resource.TestCheckResourceAttr(dataSourceFullName, "oidc_options.#", "0"),
					resource.TestCheckResourceAttr(dataSourceFullName, "saml_options.#", "1"),
					resource.TestCheckResourceAttr(dataSourceFullName, "saml_options.0.home_page_url", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(dataSourceFullName, "saml_options.0.type", "WEB_APP"),
					resource.TestCheckResourceAttr(dataSourceFullName, "saml_options.0.acs_urls.#", "2"),
					resource.TestCheckTypeSetElemAttr(dataSourceFullName, "saml_options.0.acs_urls.*", "https://pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(dataSourceFullName, "saml_options.0.acs_urls.*", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(dataSourceFullName, "saml_options.0.assertion_duration", "3600"),
					resource.TestCheckResourceAttr(dataSourceFullName, "saml_options.0.assertion_signed_enabled", "false"),
					resource.TestCheckResourceAttr(dataSourceFullName, "saml_options.0.idp_signing_key.#", "1"),
					resource.TestCheckNoResourceAttr(dataSourceFullName, "saml_options.0.idp_signing_key.0.algorithm"),
					resource.TestMatchResourceAttr(dataSourceFullName, "saml_options.0.idp_signing_key.0.key_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "saml_options.0.enable_requested_authn_context", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "saml_options.0.nameid_format", "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress"),
					resource.TestCheckResourceAttr(dataSourceFullName, "saml_options.0.response_is_signed", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "saml_options.0.slo_binding", "HTTP_REDIRECT"),
					resource.TestCheckResourceAttr(dataSourceFullName, "saml_options.0.slo_endpoint", "https://www.pingidentity.com/sloendpoint"),
					resource.TestCheckResourceAttr(dataSourceFullName, "saml_options.0.slo_response_endpoint", "https://www.pingidentity.com/sloresponseendpoint"),
					resource.TestCheckResourceAttr(dataSourceFullName, "saml_options.0.slo_window", "3"),
					resource.TestCheckResourceAttr(dataSourceFullName, "saml_options.0.sp_entity_id", fmt.Sprintf("sp:entity:%s", resourceName)),
					resource.TestCheckResourceAttr(dataSourceFullName, "saml_options.0.sp_verification_certificate_ids.#", "0"),
					resource.TestCheckResourceAttr(dataSourceFullName, "external_link_options.#", "0"),
					resource.TestCheckResourceAttr(dataSourceFullName, "hidden_from_app_portal", "true"),
				),
			},
		},
	})
}

func TestAccApplicationDataSource_SAMLAppByName(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullName := fmt.Sprintf("pingone_application.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Application_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationDataSource_SAMLAppByName(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "name", name),
					resource.TestCheckResourceAttr(dataSourceFullName, "description", ""),
					resource.TestCheckResourceAttr(dataSourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "tags.#", "0"),
					resource.TestCheckResourceAttr(dataSourceFullName, "login_page_url", ""),
					resource.TestCheckResourceAttr(dataSourceFullName, "icon.#", "0"),
					resource.TestCheckResourceAttr(dataSourceFullName, "oidc_options.#", "0"),
					resource.TestCheckResourceAttr(dataSourceFullName, "saml_options.#", "1"),
					resource.TestCheckResourceAttr(dataSourceFullName, "saml_options.0.home_page_url", ""),
					resource.TestCheckResourceAttr(dataSourceFullName, "saml_options.0.type", "WEB_APP"),
					resource.TestCheckResourceAttr(dataSourceFullName, "saml_options.0.acs_urls.#", "1"),
					resource.TestCheckTypeSetElemAttr(dataSourceFullName, "saml_options.0.acs_urls.*", "https://pingidentity.com"),
					resource.TestCheckResourceAttr(dataSourceFullName, "saml_options.0.assertion_duration", "3600"),
					resource.TestCheckResourceAttr(dataSourceFullName, "saml_options.0.assertion_signed_enabled", "true"),
					resource.TestMatchResourceAttr(dataSourceFullName, "saml_options.0.idp_signing_key_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "saml_options.0.idp_signing_key.#", "1"),
					resource.TestCheckResourceAttr(dataSourceFullName, "saml_options.0.idp_signing_key.0.algorithm", ""),
					resource.TestMatchResourceAttr(dataSourceFullName, "saml_options.0.idp_signing_key.0.key_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "saml_options.0.enable_requested_authn_context", "false"),
					resource.TestCheckResourceAttr(dataSourceFullName, "saml_options.0.nameid_format", ""),
					resource.TestCheckResourceAttr(dataSourceFullName, "saml_options.0.response_is_signed", "false"),
					resource.TestCheckResourceAttr(dataSourceFullName, "saml_options.0.slo_binding", "HTTP_POST"),
					resource.TestCheckResourceAttr(dataSourceFullName, "saml_options.0.slo_endpoint", ""),
					resource.TestCheckResourceAttr(dataSourceFullName, "saml_options.0.slo_response_endpoint", ""),
					resource.TestCheckResourceAttr(dataSourceFullName, "saml_options.0.slo_window", "0"),
					resource.TestCheckResourceAttr(dataSourceFullName, "saml_options.0.sp_entity_id", fmt.Sprintf("sp:entity:%s", resourceName)),
					resource.TestCheckResourceAttr(dataSourceFullName, "saml_options.0.sp_verification_certificate_ids.#", "0"),
					resource.TestCheckResourceAttr(dataSourceFullName, "hidden_from_app_portal", "false"),
				),
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
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
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
    type                        = "NATIVE_APP"
    grant_types                 = ["CLIENT_CREDENTIALS"]
    token_endpoint_authn_method = "CLIENT_SECRET_BASIC"

    allow_wildcards_in_redirect_uris = true

    mobile_app {
      bundle_id           = "com.%[2]s.bundle"
      package_name        = "com.%[2]s.package"
      huawei_app_id       = "%[2]s"
      huawei_package_name = "com.%[2]s.huaweipackage"

      passcode_refresh_seconds = 45

      universal_app_link = "https://applink.com"

      integrity_detection {
        enabled = true

        excluded_platforms = ["IOS"]

        cache_duration {
          amount = 30
          units  = "HOURS"
        }

        google_play {
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

  icon {
    id   = pingone_image.%[2]s.id
    href = pingone_image.%[2]s.uploaded_image[0].href
  }

  access_control_group_options {
    type = "ANY_GROUP"

    groups = [
      pingone_group.%[2]s-2.id,
      pingone_group.%[2]s-1.id
    ]
  }

  hidden_from_app_portal = true

  enabled = true

  external_link_options {
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

  external_link_options {
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

func testAccApplicationDataSource_SAMLAppByID(resourceName, name, image string) string {
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

resource "pingone_image" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  image_file_base64 = "%[4]s"
}

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "My test SAML app"
  login_page_url = "https://www.pingidentity.com"

  icon {
    id   = pingone_image.%[2]s.id
    href = pingone_image.%[2]s.uploaded_image[0].href
  }

  access_control_role_type = "ADMIN_USERS_ONLY"

  access_control_group_options {
    type = "ANY_GROUP"

    groups = [
      pingone_group.%[2]s-2.id,
      pingone_group.%[2]s-1.id
    ]
  }

  hidden_from_app_portal = true

  enabled = true

  saml_options {
    type               = "WEB_APP"
    home_page_url      = "https://www.pingidentity.com"
    acs_urls           = ["https://www.pingidentity.com", "https://pingidentity.com"]
    assertion_duration = 3600
    sp_entity_id       = "sp:entity:%[2]s"

    assertion_signed_enabled       = false
    idp_signing_key_id             = pingone_key.%[2]s.id
    enable_requested_authn_context = true
    nameid_format                  = "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress"
    response_is_signed             = true
    slo_binding                    = "HTTP_REDIRECT"
    slo_endpoint                   = "https://www.pingidentity.com/sloendpoint"
    slo_response_endpoint          = "https://www.pingidentity.com/sloresponseendpoint"
    slo_window                     = 3

    // sp_verification_certificate_ids = []

  }
}
data "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  application_id = pingone_application.%[2]s.id
}`, acctest.GenericSandboxEnvironment(), resourceName, name, image)
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

  saml_options {
    acs_urls           = ["https://pingidentity.com"]
    assertion_duration = 3600
    sp_entity_id       = "sp:entity:%[2]s"
    idp_signing_key_id = pingone_key.%[2]s.id
  }
}
data "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  depends_on = [pingone_application.%[2]s]
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
