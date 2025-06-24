// Copyright © 2025 Ping Identity Corporation

package sso_test

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/base"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/sso"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccApplication_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	var applicationID, environmentID string

	var p1Client *client.Client
	var ctx = context.Background()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)

			p1Client = acctest.PreCheckTestClient(ctx, t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Application_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccApplicationConfig_OIDC_MinimalWeb(resourceName, name),
				Check:  sso.Application_GetIDs(resourceFullName, &environmentID, &applicationID),
			},
			{
				PreConfig: func() {
					sso.Application_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID, applicationID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the environment
			{
				Config: testAccApplicationConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check:  sso.Application_GetIDs(resourceFullName, &environmentID, &applicationID),
			},
			{
				PreConfig: func() {
					base.Environment_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccApplication_NewEnv(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Application_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
				),
			},
		},
	})
}

func TestAccApplication_OIDCFullWeb(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application.%s", resourceName)

	name := resourceName

	data, _ := os.ReadFile("../../acctest/test_assets/image/image-logo.gif")
	image := base64.StdEncoding.EncodeToString(data)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Application_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationConfig_OIDC_FullWeb(resourceName, name, image),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", "My test OIDC app"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "login_page_url", "https://www.pingidentity.com"),
					resource.TestMatchResourceAttr(resourceFullName, "icon.id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "icon.href", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_role_type", "ADMIN_USERS_ONLY"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.groups.#", "1"),
					resource.TestMatchResourceAttr(resourceFullName, "access_control_group_options.groups.0", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.type", "ANY_GROUP"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.type", "WEB_APP"),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.client_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.home_page_url", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.initiate_login_uri", "https://www.pingidentity.com/initiate"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.target_link_uri", "https://www.pingidentity.com/target"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.grant_types.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.grant_types.*", "AUTHORIZATION_CODE"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.grant_types.*", "REFRESH_TOKEN"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.response_types.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.response_types.*", "CODE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.token_endpoint_auth_method", "CLIENT_SECRET_BASIC"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.par_requirement", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.par_timeout", "60"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.pkce_enforcement", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.idp_signoff", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.redirect_uris.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.redirect_uris.*", "https://pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.redirect_uris.*", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.allow_wildcard_in_redirect_uris", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.post_logout_redirect_uris.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.post_logout_redirect_uris.*", "https://pingidentity.com/logout"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.post_logout_redirect_uris.*", "https://www.pingidentity.com/logout"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.refresh_token_duration", "3000000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.refresh_token_rolling_duration", "30000000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.refresh_token_rolling_grace_period_duration", "80000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.additional_refresh_token_replay_protection_enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.support_unsigned_request_object", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.require_signed_request_object", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.cors_settings.behavior", "ALLOW_SPECIFIC_ORIGINS"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.cors_settings.origins.#", "8"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "http://localhost"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "https://localhost"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "http://auth.pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "https://auth.pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "http://*.pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "https://*.pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "http://192.168.1.1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "https://192.168.1.1"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.mobile_app"),
					resource.TestCheckNoResourceAttr(resourceFullName, "saml_options"),
					resource.TestCheckNoResourceAttr(resourceFullName, "external_link_options"),
					resource.TestCheckNoResourceAttr(resourceFullName, "wsfed_options"),
					resource.TestCheckResourceAttr(resourceFullName, "hidden_from_app_portal", "true"),
				),
			},
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccApplication_OIDCMinimalWeb(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Application_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationConfig_OIDC_MinimalWeb(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckNoResourceAttr(resourceFullName, "description"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckNoResourceAttr(resourceFullName, "login_page_url"),
					resource.TestCheckNoResourceAttr(resourceFullName, "icon"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.type", "WEB_APP"),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.client_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.home_page_url"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.initiate_login_uri"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.target_link_uri"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.grant_types.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.grant_types.*", "AUTHORIZATION_CODE"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.grant_types.*", "REFRESH_TOKEN"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.response_types.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.response_types.*", "CODE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.token_endpoint_auth_method", "CLIENT_SECRET_BASIC"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.par_requirement", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.par_timeout", "60"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.pkce_enforcement", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.redirect_uris.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.redirect_uris.*", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.allow_wildcard_in_redirect_uris", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.post_logout_redirect_uris.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.refresh_token_duration", "2592000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.refresh_token_rolling_duration", "15552000"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.refresh_token_rolling_grace_period_duration"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.additional_refresh_token_replay_protection_enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.support_unsigned_request_object", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.require_signed_request_object", "false"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.cors_settings"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.mobile_app"),
					resource.TestCheckNoResourceAttr(resourceFullName, "saml_options"),
					resource.TestCheckNoResourceAttr(resourceFullName, "external_link_options"),
					resource.TestCheckNoResourceAttr(resourceFullName, "wsfed_options"),
					resource.TestCheckResourceAttr(resourceFullName, "hidden_from_app_portal", "false"),
				),
			},
		},
	})
}

func TestAccApplication_OIDCWebUpdate(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application.%s", resourceName)

	name := resourceName

	data, _ := os.ReadFile("../../acctest/test_assets/image/image-logo.gif")
	image := base64.StdEncoding.EncodeToString(data)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Application_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationConfig_OIDC_MinimalWeb(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckNoResourceAttr(resourceFullName, "description"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckNoResourceAttr(resourceFullName, "login_page_url"),
					resource.TestCheckNoResourceAttr(resourceFullName, "icon"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.type", "WEB_APP"),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.client_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.home_page_url"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.initiate_login_uri"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.target_link_uri"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.grant_types.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.grant_types.*", "AUTHORIZATION_CODE"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.grant_types.*", "REFRESH_TOKEN"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.response_types.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.response_types.*", "CODE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.token_endpoint_auth_method", "CLIENT_SECRET_BASIC"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.par_requirement", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.par_timeout", "60"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.pkce_enforcement", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.redirect_uris.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.redirect_uris.*", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.allow_wildcard_in_redirect_uris", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.post_logout_redirect_uris.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.refresh_token_duration", "2592000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.refresh_token_rolling_duration", "15552000"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.refresh_token_rolling_grace_period_duration"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.additional_refresh_token_replay_protection_enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.support_unsigned_request_object", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.require_signed_request_object", "false"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.cors_settings"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.mobile_app"),
					resource.TestCheckNoResourceAttr(resourceFullName, "saml_options"),
					resource.TestCheckNoResourceAttr(resourceFullName, "external_link_options"),
					resource.TestCheckNoResourceAttr(resourceFullName, "wsfed_options"),
					resource.TestCheckResourceAttr(resourceFullName, "hidden_from_app_portal", "false"),
				),
			},
			{
				Config: testAccApplicationConfig_OIDC_FullWeb(resourceName, name, image),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", "My test OIDC app"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "login_page_url", "https://www.pingidentity.com"),
					resource.TestMatchResourceAttr(resourceFullName, "icon.id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "icon.href", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_role_type", "ADMIN_USERS_ONLY"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.groups.#", "1"),
					resource.TestMatchResourceAttr(resourceFullName, "access_control_group_options.groups.0", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.type", "ANY_GROUP"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.type", "WEB_APP"),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.client_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.home_page_url", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.initiate_login_uri", "https://www.pingidentity.com/initiate"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.target_link_uri", "https://www.pingidentity.com/target"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.grant_types.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.grant_types.*", "AUTHORIZATION_CODE"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.grant_types.*", "REFRESH_TOKEN"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.response_types.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.response_types.*", "CODE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.token_endpoint_auth_method", "CLIENT_SECRET_BASIC"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.par_requirement", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.par_timeout", "60"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.pkce_enforcement", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.idp_signoff", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.redirect_uris.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.redirect_uris.*", "https://pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.redirect_uris.*", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.allow_wildcard_in_redirect_uris", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.post_logout_redirect_uris.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.post_logout_redirect_uris.*", "https://pingidentity.com/logout"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.post_logout_redirect_uris.*", "https://www.pingidentity.com/logout"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.refresh_token_duration", "3000000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.refresh_token_rolling_duration", "30000000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.refresh_token_rolling_grace_period_duration", "80000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.additional_refresh_token_replay_protection_enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.support_unsigned_request_object", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.require_signed_request_object", "false"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.mobile_app"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.cors_settings.behavior", "ALLOW_SPECIFIC_ORIGINS"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.cors_settings.origins.#", "8"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "http://localhost"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "https://localhost"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "http://auth.pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "https://auth.pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "http://*.pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "https://*.pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "http://192.168.1.1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "https://192.168.1.1"),
					resource.TestCheckNoResourceAttr(resourceFullName, "saml_options"),
					resource.TestCheckNoResourceAttr(resourceFullName, "external_link_options"),
					resource.TestCheckNoResourceAttr(resourceFullName, "wsfed_options"),
					resource.TestCheckResourceAttr(resourceFullName, "hidden_from_app_portal", "true"),
				),
			},
			{
				Config: testAccApplicationConfig_OIDC_MinimalWeb(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckNoResourceAttr(resourceFullName, "description"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckNoResourceAttr(resourceFullName, "login_page_url"),
					resource.TestCheckNoResourceAttr(resourceFullName, "icon"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.type", "WEB_APP"),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.client_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.home_page_url"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.initiate_login_uri"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.target_link_uri"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.grant_types.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.grant_types.*", "AUTHORIZATION_CODE"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.grant_types.*", "REFRESH_TOKEN"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.response_types.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.response_types.*", "CODE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.token_endpoint_auth_method", "CLIENT_SECRET_BASIC"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.par_requirement", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.par_timeout", "60"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.pkce_enforcement", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.redirect_uris.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.redirect_uris.*", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.allow_wildcard_in_redirect_uris", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.post_logout_redirect_uris.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.refresh_token_duration", "2592000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.refresh_token_rolling_duration", "15552000"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.refresh_token_rolling_grace_period_duration"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.additional_refresh_token_replay_protection_enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.support_unsigned_request_object", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.require_signed_request_object", "false"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.cors_settings"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.mobile_app"),
					resource.TestCheckNoResourceAttr(resourceFullName, "saml_options"),
					resource.TestCheckNoResourceAttr(resourceFullName, "external_link_options"),
					resource.TestCheckNoResourceAttr(resourceFullName, "wsfed_options"),
					resource.TestCheckResourceAttr(resourceFullName, "hidden_from_app_portal", "false"),
				),
			},
		},
	})
}

func TestAccApplication_OIDCFullNative(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application.%s", resourceName)

	name := resourceName

	data, _ := os.ReadFile("../../acctest/test_assets/image/image-logo.gif")
	image := base64.StdEncoding.EncodeToString(data)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Application_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationConfig_OIDC_FullNative(resourceName, name, image),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", "My test OIDC app"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "login_page_url", "https://www.pingidentity.com"),
					resource.TestMatchResourceAttr(resourceFullName, "icon.id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "icon.href", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_role_type", "ADMIN_USERS_ONLY"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.groups.#", "1"),
					resource.TestMatchResourceAttr(resourceFullName, "access_control_group_options.groups.0", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.type", "ANY_GROUP"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.type", "NATIVE_APP"),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.client_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.home_page_url"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.initiate_login_uri"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.target_link_uri"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.grant_types.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.grant_types.*", "CLIENT_CREDENTIALS"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.response_types.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.token_endpoint_auth_method", "CLIENT_SECRET_BASIC"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.par_requirement", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.par_timeout", "60"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.pkce_enforcement", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.redirect_uris.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.allow_wildcard_in_redirect_uris", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.post_logout_redirect_uris.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.refresh_token_duration", "2592000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.refresh_token_rolling_duration", "15552000"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.refresh_token_rolling_grace_period_duration"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.additional_refresh_token_replay_protection_enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.support_unsigned_request_object", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.require_signed_request_object", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.cors_settings.behavior", "ALLOW_NO_ORIGINS"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.cors_settings.origins.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.bundle_id", fmt.Sprintf("com.%s.bundle", resourceName)),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.package_name", fmt.Sprintf("com.%s.package", resourceName)),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.huawei_app_id", resourceName),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.huawei_package_name", fmt.Sprintf("com.%s.huaweipackage", resourceName)),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.passcode_refresh_seconds", "45"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.universal_app_link", "https://applink.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.excluded_platforms.0", "IOS"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.cache_duration.amount", "30"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.cache_duration.units", "HOURS"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.google_play.verification_type", "INTERNAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.google_play.verification_key", "verificationkeydoesnotexist"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.google_play.decryption_key", "decryptionkeydoesnotexist"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.google_play.service_account_credentials_json"),
					resource.TestCheckNoResourceAttr(resourceFullName, "saml_options"),
					resource.TestCheckNoResourceAttr(resourceFullName, "external_link_options"),
					resource.TestCheckNoResourceAttr(resourceFullName, "wsfed_options"),
					resource.TestCheckResourceAttr(resourceFullName, "hidden_from_app_portal", "true"),
				),
			},
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"oidc_options.mobile_app.integrity_detection.google_play.decryption_key",
					"oidc_options.mobile_app.integrity_detection.google_play.verification_key",
					"oidc_options.mobile_app.integrity_detection.google_play.service_account_credentials_json",
				},
			},
		},
	})
}

func TestAccApplication_OIDCMinimalNative(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Application_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationConfig_OIDC_MinimalNative(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckNoResourceAttr(resourceFullName, "description"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckNoResourceAttr(resourceFullName, "login_page_url"),
					resource.TestCheckNoResourceAttr(resourceFullName, "icon"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.type", "NATIVE_APP"),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.client_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.home_page_url"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.initiate_login_uri"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.target_link_uri"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.grant_types.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.grant_types.*", "CLIENT_CREDENTIALS"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.response_types.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.token_endpoint_auth_method", "CLIENT_SECRET_BASIC"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.par_requirement", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.par_timeout", "60"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.pkce_enforcement", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.redirect_uris.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.allow_wildcard_in_redirect_uris", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.post_logout_redirect_uris.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.refresh_token_duration", "2592000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.refresh_token_rolling_duration", "15552000"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.refresh_token_rolling_grace_period_duration"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.additional_refresh_token_replay_protection_enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.support_unsigned_request_object", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.require_signed_request_object", "false"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.cors_settings"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.mobile_app.bundle_id"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.mobile_app.package_name"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.mobile_app.huawei_app_id"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.mobile_app.huawei_package_name"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.mobile_app.universal_app_link"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.excluded_platforms.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.passcode_refresh_seconds", "30"),
					resource.TestCheckNoResourceAttr(resourceFullName, "saml_options"),
					resource.TestCheckNoResourceAttr(resourceFullName, "external_link_options"),
					resource.TestCheckNoResourceAttr(resourceFullName, "wsfed_options"),
					resource.TestCheckResourceAttr(resourceFullName, "hidden_from_app_portal", "false"),
				),
			},
		},
	})
}

func TestAccApplication_OIDCNativeUpdate(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application.%s", resourceName)

	name := resourceName

	data, _ := os.ReadFile("../../acctest/test_assets/image/image-logo.gif")
	image := base64.StdEncoding.EncodeToString(data)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Application_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationConfig_OIDC_FullNative(resourceName, name, image),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", "My test OIDC app"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "login_page_url", "https://www.pingidentity.com"),
					resource.TestMatchResourceAttr(resourceFullName, "icon.id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "icon.href", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_role_type", "ADMIN_USERS_ONLY"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.groups.#", "1"),
					resource.TestMatchResourceAttr(resourceFullName, "access_control_group_options.groups.0", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.type", "ANY_GROUP"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.type", "NATIVE_APP"),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.client_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.home_page_url"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.initiate_login_uri"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.target_link_uri"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.grant_types.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.grant_types.*", "CLIENT_CREDENTIALS"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.response_types.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.token_endpoint_auth_method", "CLIENT_SECRET_BASIC"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.par_requirement", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.par_timeout", "60"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.pkce_enforcement", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.redirect_uris.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.allow_wildcard_in_redirect_uris", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.post_logout_redirect_uris.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.refresh_token_duration", "2592000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.refresh_token_rolling_duration", "15552000"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.refresh_token_rolling_grace_period_duration"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.additional_refresh_token_replay_protection_enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.support_unsigned_request_object", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.require_signed_request_object", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.cors_settings.behavior", "ALLOW_NO_ORIGINS"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.cors_settings.origins.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.bundle_id", fmt.Sprintf("com.%s.bundle", resourceName)),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.package_name", fmt.Sprintf("com.%s.package", resourceName)),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.huawei_app_id", resourceName),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.huawei_package_name", fmt.Sprintf("com.%s.huaweipackage", resourceName)),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.passcode_refresh_seconds", "45"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.universal_app_link", "https://applink.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.excluded_platforms.0", "IOS"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.cache_duration.amount", "30"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.cache_duration.units", "HOURS"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.google_play.verification_type", "INTERNAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.google_play.verification_key", "verificationkeydoesnotexist"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.google_play.decryption_key", "decryptionkeydoesnotexist"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.google_play.service_account_credentials_json"),
					resource.TestCheckNoResourceAttr(resourceFullName, "saml_options"),
					resource.TestCheckNoResourceAttr(resourceFullName, "external_link_options"),
					resource.TestCheckNoResourceAttr(resourceFullName, "wsfed_options"),
					resource.TestCheckResourceAttr(resourceFullName, "hidden_from_app_portal", "true"),
				),
			},
			{
				Config: testAccApplicationConfig_OIDC_MinimalNative(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckNoResourceAttr(resourceFullName, "description"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckNoResourceAttr(resourceFullName, "login_page_url"),
					resource.TestCheckNoResourceAttr(resourceFullName, "icon"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.type", "NATIVE_APP"),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.client_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.home_page_url"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.initiate_login_uri"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.target_link_uri"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.grant_types.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.grant_types.*", "CLIENT_CREDENTIALS"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.response_types.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.token_endpoint_auth_method", "CLIENT_SECRET_BASIC"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.par_requirement", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.par_timeout", "60"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.pkce_enforcement", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.redirect_uris.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.allow_wildcard_in_redirect_uris", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.post_logout_redirect_uris.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.refresh_token_duration", "2592000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.refresh_token_rolling_duration", "15552000"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.refresh_token_rolling_grace_period_duration"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.additional_refresh_token_replay_protection_enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.support_unsigned_request_object", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.require_signed_request_object", "false"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.cors_settings"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.mobile_app.bundle_id"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.mobile_app.package_name"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.mobile_app.huawei_app_id"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.mobile_app.huawei_package_name"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.mobile_app.universal_app_link"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.excluded_platforms.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.passcode_refresh_seconds", "30"),
					resource.TestCheckNoResourceAttr(resourceFullName, "saml_options"),
					resource.TestCheckNoResourceAttr(resourceFullName, "external_link_options"),
					resource.TestCheckNoResourceAttr(resourceFullName, "wsfed_options"),
					resource.TestCheckResourceAttr(resourceFullName, "hidden_from_app_portal", "false"),
				),
			},
			{
				Config: testAccApplicationConfig_OIDC_FullNative(resourceName, name, image),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", "My test OIDC app"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "login_page_url", "https://www.pingidentity.com"),
					resource.TestMatchResourceAttr(resourceFullName, "icon.id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "icon.href", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_role_type", "ADMIN_USERS_ONLY"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.groups.#", "1"),
					resource.TestMatchResourceAttr(resourceFullName, "access_control_group_options.groups.0", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.type", "ANY_GROUP"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.type", "NATIVE_APP"),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.client_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.home_page_url"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.initiate_login_uri"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.target_link_uri"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.grant_types.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.grant_types.*", "CLIENT_CREDENTIALS"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.response_types.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.token_endpoint_auth_method", "CLIENT_SECRET_BASIC"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.par_requirement", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.par_timeout", "60"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.pkce_enforcement", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.redirect_uris.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.allow_wildcard_in_redirect_uris", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.post_logout_redirect_uris.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.refresh_token_duration", "2592000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.refresh_token_rolling_duration", "15552000"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.refresh_token_rolling_grace_period_duration"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.additional_refresh_token_replay_protection_enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.support_unsigned_request_object", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.require_signed_request_object", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.cors_settings.behavior", "ALLOW_NO_ORIGINS"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.cors_settings.origins.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.bundle_id", fmt.Sprintf("com.%s.bundle", resourceName)),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.package_name", fmt.Sprintf("com.%s.package", resourceName)),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.huawei_app_id", resourceName),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.huawei_package_name", fmt.Sprintf("com.%s.huaweipackage", resourceName)),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.passcode_refresh_seconds", "45"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.universal_app_link", "https://applink.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.excluded_platforms.0", "IOS"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.cache_duration.amount", "30"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.cache_duration.units", "HOURS"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.google_play.verification_type", "INTERNAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.google_play.verification_key", "verificationkeydoesnotexist"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.google_play.decryption_key", "decryptionkeydoesnotexist"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.google_play.service_account_credentials_json"),
					resource.TestCheckNoResourceAttr(resourceFullName, "saml_options"),
					resource.TestCheckNoResourceAttr(resourceFullName, "external_link_options"),
					resource.TestCheckNoResourceAttr(resourceFullName, "wsfed_options"),
					resource.TestCheckResourceAttr(resourceFullName, "hidden_from_app_portal", "true"),
				),
			},
		},
	})
}

func TestAccApplication_NativeKerberos(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application.%s", resourceName)

	name := resourceName

	withKerberosTestStep := resource.TestStep{
		Config: testAccApplicationConfig_OIDC_NativeKerberos(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "oidc_options.certificate_based_authentication.key_id", verify.P1ResourceIDRegexpFullString),
		),
	}

	withoutKerberosTestStep := resource.TestStep{
		Config: testAccApplicationConfig_OIDC_MinimalNative(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "oidc_options.certificate_based_authentication.#", "0"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
			acctest.PreCheckRegionSupportsWorkforce(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Application_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Invalid configs
			{
				Config:      testAccApplicationConfig_OIDC_NativeKerberosIncorrectKeyType(resourceName, name),
				ExpectError: regexp.MustCompile("Error when calling `CreateApplication`: Key with ID '[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}' in Environment '[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}' is not for ISSUANCE. Usage type should be ISSUANCE."),
			},
			{
				Config:      testAccApplicationConfig_OIDC_NativeKerberosIncorrectApplicationType(resourceName, name),
				ExpectError: regexp.MustCompile("Invalid configuration"),
			},
			// With
			withKerberosTestStep,
			{
				Config:  testAccApplicationConfig_OIDC_NativeKerberos(resourceName, name),
				Destroy: true,
			},
			// Without
			withoutKerberosTestStep,
			{
				Config:  testAccApplicationConfig_OIDC_MinimalNative(resourceName, name),
				Destroy: true,
			},
			// Change
			withKerberosTestStep,
			withoutKerberosTestStep,
			withKerberosTestStep,
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:  testAccApplicationConfig_OIDC_MinimalNative(resourceName, name),
				Destroy: true,
			},
		},
	})
}

func TestAccApplication_NativeMobile(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application.%s", resourceName)

	name := resourceName

	withMobileTestStepFull := resource.TestStep{
		Config: testAccApplicationConfig_OIDC_NativeMobile_Full(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.bundle_id", fmt.Sprintf("com.%s.bundle", resourceName)),
			resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.package_name", fmt.Sprintf("com.%s.package", resourceName)),
			resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.huawei_app_id", resourceName),
			resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.huawei_package_name", fmt.Sprintf("com.%s.huaweipackage", resourceName)),
			resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.passcode_refresh_seconds", "45"),
			resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.universal_app_link", "https://applink.com"),
			resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.enabled", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.excluded_platforms.0", "IOS"),
			resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.cache_duration.amount", "30"),
			resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.cache_duration.units", "HOURS"),
			resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.google_play.verification_type", "INTERNAL"),
			resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.google_play.verification_key", "verificationkeydoesnotexist"),
			resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.google_play.decryption_key", "decryptionkeydoesnotexist"),
			resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.google_play.service_account_credentials_json"),
		),
	}

	withMobileTestStepMinimal := resource.TestStep{
		Config: testAccApplicationConfig_OIDC_NativeMobile_Minimal(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.mobile_app.bundle_id"),
			resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.mobile_app.package_name"),
			resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.mobile_app.huawei_app_id"),
			resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.mobile_app.huawei_package_name"),
			resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.passcode_refresh_seconds", "30"),
			resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.mobile_app.universal_app_link"),
			resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.enabled", "false"),
		),
	}

	withoutMobileTestStep := resource.TestStep{
		Config: testAccApplicationConfig_OIDC_MinimalNative(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.mobile_app.bundle_id"),
			resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.mobile_app.package_name"),
			resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.mobile_app.huawei_app_id"),
			resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.mobile_app.huawei_package_name"),
			resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.mobile_app.universal_app_link"),
			resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.enabled", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.excluded_platforms.#", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.passcode_refresh_seconds", "30"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Application_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// With
			withMobileTestStepFull,
			{
				Config:  testAccApplicationConfig_OIDC_NativeMobile_Full(resourceName, name),
				Destroy: true,
			},
			// Without
			withMobileTestStepMinimal,
			{
				Config:  testAccApplicationConfig_OIDC_NativeMobile_Minimal(resourceName, name),
				Destroy: true,
			},
			// Change
			withMobileTestStepFull,
			withMobileTestStepMinimal,
			withMobileTestStepFull,
			withoutMobileTestStep,
			withMobileTestStepFull,
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"oidc_options.mobile_app.integrity_detection.google_play.decryption_key",
					"oidc_options.mobile_app.integrity_detection.google_play.verification_key",
				},
			},
			{
				Config:  testAccApplicationConfig_OIDC_NativeMobile_Full(resourceName, name),
				Destroy: true,
			},
		},
	})
}

func TestAccApplication_NativeMobile_IntegrityDetection(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application.%s", resourceName)

	name := resourceName

	googleJsonKey := os.Getenv("PINGONE_GOOGLE_JSON_KEY")

	testStepFull := resource.TestStep{
		Config: testAccApplicationConfig_OIDC_NativeMobile_IntegrityDetection_Full(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.enabled", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.excluded_platforms.0", "IOS"),
			resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.cache_duration.amount", "45"),
			resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.cache_duration.units", "MINUTES"),
			resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.google_play.verification_type", "INTERNAL"),
			resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.google_play.verification_key", "verificationkeydoesnotexist"),
			resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.google_play.decryption_key", "decryptionkeydoesnotexist"),
			resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.google_play.service_account_credentials_json"),
		),
	}

	testStepMinimal := resource.TestStep{
		Config: testAccApplicationConfig_OIDC_NativeMobile_IntegrityDetection_Minimal(resourceName, name, googleJsonKey),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.enabled", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.excluded_platforms.#", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.cache_duration.amount", "30"),
			resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.cache_duration.units", "MINUTES"),
			resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.google_play.verification_type", "GOOGLE"),
			resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.google_play.verification_key"),
			resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.google_play.decryption_key"),
			resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.google_play.service_account_credentials_json", googleJsonKey),
		),
	}

	excludeGoogle := resource.TestStep{
		Config: testAccApplicationConfig_OIDC_NativeMobile_IntegrityDetection_ExcludeGoogle(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.enabled", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.excluded_platforms.0", "GOOGLE"),
			resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.cache_duration.amount", "30"),
			resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.cache_duration.units", "MINUTES"),
			resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.google_play.#", "0"),
		),
	}

	excludeIOS := resource.TestStep{
		Config: testAccApplicationConfig_OIDC_NativeMobile_IntegrityDetection_ExcludeIOS(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.enabled", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.excluded_platforms.0", "IOS"),
			resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.cache_duration.amount", "30"),
			resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.cache_duration.units", "MINUTES"),
			resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.google_play.verification_type", "INTERNAL"),
			resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.google_play.verification_key", "verificationkeydoesnotexist"),
			resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.google_play.decryption_key", "decryptionkeydoesnotexist"),
			resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.google_play.service_account_credentials_json"),
		),
	}

	testStepWithout := resource.TestStep{
		Config: testAccApplicationConfig_OIDC_NativeMobile_Minimal(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.enabled", "false"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
			acctest.PreCheckGoogleJSONKey(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Application_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// With
			testStepFull,
			{
				Config:  testAccApplicationConfig_OIDC_NativeMobile_IntegrityDetection_Full(resourceName, name),
				Destroy: true,
			},
			// Without
			testStepMinimal,
			{
				Config:  testAccApplicationConfig_OIDC_NativeMobile_IntegrityDetection_Minimal(resourceName, name, googleJsonKey),
				Destroy: true,
			},
			// Without
			excludeGoogle,
			{
				Config:  testAccApplicationConfig_OIDC_NativeMobile_IntegrityDetection_ExcludeGoogle(resourceName, name),
				Destroy: true,
			},
			// Without
			excludeIOS,
			{
				Config:  testAccApplicationConfig_OIDC_NativeMobile_IntegrityDetection_ExcludeIOS(resourceName, name),
				Destroy: true,
			},
			// Change
			testStepFull,
			testStepMinimal,
			testStepFull,
			testStepWithout,
			testStepFull,
			excludeGoogle,
			excludeIOS,
			{
				Config:  testAccApplicationConfig_OIDC_NativeMobile_IntegrityDetection_Full(resourceName, name),
				Destroy: true,
			},
		},
	})
}

func TestAccApplication_OIDCFullCustom(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application.%s", resourceName)

	name := resourceName

	data, _ := os.ReadFile("../../acctest/test_assets/image/image-logo.gif")
	image := base64.StdEncoding.EncodeToString(data)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Application_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationConfig_OIDC_FullCustom(resourceName, name, image),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", "My test OIDC app"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "login_page_url", "https://www.pingidentity.com"),
					resource.TestMatchResourceAttr(resourceFullName, "icon.id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "icon.href", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_role_type", "ADMIN_USERS_ONLY"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.groups.#", "1"),
					resource.TestMatchResourceAttr(resourceFullName, "access_control_group_options.groups.0", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.type", "ANY_GROUP"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.type", "CUSTOM_APP"),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.client_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.home_page_url", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.initiate_login_uri", "https://www.pingidentity.com/initiate"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.target_link_uri", "https://www.pingidentity.com/target"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.grant_types.#", "4"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.grant_types.*", "AUTHORIZATION_CODE"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.grant_types.*", "CLIENT_CREDENTIALS"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.grant_types.*", "IMPLICIT"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.grant_types.*", "REFRESH_TOKEN"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.device_path_id"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.device_custom_verification_uri"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.device_timeout", "600"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.device_polling_interval", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.idp_signoff", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.response_types.#", "3"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.response_types.*", "CODE"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.response_types.*", "TOKEN"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.response_types.*", "ID_TOKEN"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.token_endpoint_auth_method", "CLIENT_SECRET_BASIC"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.par_requirement", "REQUIRED"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.par_timeout", "180"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.pkce_enforcement", "REQUIRED"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.redirect_uris.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.redirect_uris.*", "https://pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.redirect_uris.*", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.allow_wildcard_in_redirect_uris", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.post_logout_redirect_uris.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.post_logout_redirect_uris.*", "https://pingidentity.com/logout"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.post_logout_redirect_uris.*", "https://www.pingidentity.com/logout"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.refresh_token_duration", "3000000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.refresh_token_rolling_duration", "30000000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.refresh_token_rolling_grace_period_duration", "80000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.additional_refresh_token_replay_protection_enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.support_unsigned_request_object", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.require_signed_request_object", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.cors_settings.behavior", "ALLOW_SPECIFIC_ORIGINS"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.cors_settings.origins.#", "8"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "http://localhost"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "https://localhost"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "http://auth.pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "https://auth.pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "http://*.pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "https://*.pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "http://192.168.1.1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "https://192.168.1.1"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.mobile_app"),
					resource.TestCheckNoResourceAttr(resourceFullName, "saml_options"),
					resource.TestCheckNoResourceAttr(resourceFullName, "external_link_options"),
					resource.TestCheckNoResourceAttr(resourceFullName, "wsfed_options"),
					resource.TestCheckResourceAttr(resourceFullName, "hidden_from_app_portal", "true"),
				),
			},
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccApplication_OIDCMinimalCustom(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Application_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationConfig_OIDC_MinimalCustom(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckNoResourceAttr(resourceFullName, "description"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckNoResourceAttr(resourceFullName, "login_page_url"),
					resource.TestCheckNoResourceAttr(resourceFullName, "icon"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.type", "CUSTOM_APP"),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.client_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.home_page_url"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.initiate_login_uri"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.target_link_uri"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.grant_types.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.grant_types.*", "CLIENT_CREDENTIALS"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.device_path_id"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.device_custom_verification_uri"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.device_timeout", "600"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.device_polling_interval", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.idp_signoff", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.response_types.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.token_endpoint_auth_method", "CLIENT_SECRET_BASIC"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.par_requirement", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.par_timeout", "60"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.pkce_enforcement", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.redirect_uris.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.allow_wildcard_in_redirect_uris", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.post_logout_redirect_uris.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.refresh_token_duration", "2592000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.refresh_token_rolling_duration", "15552000"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.refresh_token_rolling_grace_period_duration"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.additional_refresh_token_replay_protection_enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.support_unsigned_request_object", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.require_signed_request_object", "false"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.cors_settings"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.mobile_app"),
					resource.TestCheckNoResourceAttr(resourceFullName, "saml_options"),
					resource.TestCheckNoResourceAttr(resourceFullName, "external_link_options"),
					resource.TestCheckNoResourceAttr(resourceFullName, "wsfed_options"),
					resource.TestCheckResourceAttr(resourceFullName, "hidden_from_app_portal", "false"),
				),
			},
		},
	})
}

func TestAccApplication_OIDCCustomUpdate(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application.%s", resourceName)

	name := resourceName

	data, _ := os.ReadFile("../../acctest/test_assets/image/image-logo.gif")
	image := base64.StdEncoding.EncodeToString(data)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Application_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationConfig_OIDC_FullCustom(resourceName, name, image),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", "My test OIDC app"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "login_page_url", "https://www.pingidentity.com"),
					resource.TestMatchResourceAttr(resourceFullName, "icon.id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "icon.href", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_role_type", "ADMIN_USERS_ONLY"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.groups.#", "1"),
					resource.TestMatchResourceAttr(resourceFullName, "access_control_group_options.groups.0", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.type", "ANY_GROUP"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.type", "CUSTOM_APP"),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.client_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.home_page_url", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.initiate_login_uri", "https://www.pingidentity.com/initiate"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.target_link_uri", "https://www.pingidentity.com/target"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.grant_types.#", "4"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.grant_types.*", "AUTHORIZATION_CODE"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.grant_types.*", "CLIENT_CREDENTIALS"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.grant_types.*", "IMPLICIT"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.grant_types.*", "REFRESH_TOKEN"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.device_path_id"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.device_custom_verification_uri"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.device_timeout", "600"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.device_polling_interval", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.idp_signoff", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.response_types.#", "3"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.response_types.*", "CODE"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.response_types.*", "TOKEN"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.response_types.*", "ID_TOKEN"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.token_endpoint_auth_method", "CLIENT_SECRET_BASIC"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.par_requirement", "REQUIRED"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.par_timeout", "180"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.pkce_enforcement", "REQUIRED"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.redirect_uris.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.redirect_uris.*", "https://pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.redirect_uris.*", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.allow_wildcard_in_redirect_uris", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.post_logout_redirect_uris.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.post_logout_redirect_uris.*", "https://pingidentity.com/logout"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.post_logout_redirect_uris.*", "https://www.pingidentity.com/logout"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.refresh_token_duration", "3000000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.refresh_token_rolling_duration", "30000000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.refresh_token_rolling_grace_period_duration", "80000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.additional_refresh_token_replay_protection_enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.support_unsigned_request_object", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.require_signed_request_object", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.cors_settings.behavior", "ALLOW_SPECIFIC_ORIGINS"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.cors_settings.origins.#", "8"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "http://localhost"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "https://localhost"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "http://auth.pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "https://auth.pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "http://*.pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "https://*.pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "http://192.168.1.1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "https://192.168.1.1"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.mobile_app"),
					resource.TestCheckNoResourceAttr(resourceFullName, "saml_options"),
					resource.TestCheckNoResourceAttr(resourceFullName, "external_link_options"),
					resource.TestCheckNoResourceAttr(resourceFullName, "wsfed_options"),
					resource.TestCheckResourceAttr(resourceFullName, "hidden_from_app_portal", "true"),
				),
			},
			{
				Config: testAccApplicationConfig_OIDC_MinimalCustom(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckNoResourceAttr(resourceFullName, "description"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckNoResourceAttr(resourceFullName, "login_page_url"),
					resource.TestCheckNoResourceAttr(resourceFullName, "icon"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.type", "CUSTOM_APP"),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.client_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.home_page_url"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.initiate_login_uri"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.target_link_uri"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.grant_types.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.grant_types.*", "CLIENT_CREDENTIALS"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.device_path_id"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.device_custom_verification_uri"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.device_timeout", "600"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.device_polling_interval", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.idp_signoff", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.response_types.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.token_endpoint_auth_method", "CLIENT_SECRET_BASIC"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.par_requirement", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.par_timeout", "60"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.pkce_enforcement", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.redirect_uris.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.allow_wildcard_in_redirect_uris", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.post_logout_redirect_uris.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.refresh_token_duration", "2592000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.refresh_token_rolling_duration", "15552000"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.refresh_token_rolling_grace_period_duration"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.additional_refresh_token_replay_protection_enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.support_unsigned_request_object", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.require_signed_request_object", "false"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.cors_settings"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.mobile_app"),
					resource.TestCheckNoResourceAttr(resourceFullName, "saml_options"),
					resource.TestCheckNoResourceAttr(resourceFullName, "external_link_options"),
					resource.TestCheckNoResourceAttr(resourceFullName, "wsfed_options"),
					resource.TestCheckResourceAttr(resourceFullName, "hidden_from_app_portal", "false"),
				),
			},
			{
				Config: testAccApplicationConfig_OIDC_FullCustom(resourceName, name, image),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", "My test OIDC app"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "login_page_url", "https://www.pingidentity.com"),
					resource.TestMatchResourceAttr(resourceFullName, "icon.id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "icon.href", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_role_type", "ADMIN_USERS_ONLY"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.groups.#", "1"),
					resource.TestMatchResourceAttr(resourceFullName, "access_control_group_options.groups.0", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.type", "ANY_GROUP"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.type", "CUSTOM_APP"),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.client_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.home_page_url", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.initiate_login_uri", "https://www.pingidentity.com/initiate"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.target_link_uri", "https://www.pingidentity.com/target"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.grant_types.#", "4"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.grant_types.*", "AUTHORIZATION_CODE"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.grant_types.*", "CLIENT_CREDENTIALS"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.grant_types.*", "IMPLICIT"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.grant_types.*", "REFRESH_TOKEN"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.device_path_id"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.device_custom_verification_uri"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.device_timeout", "600"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.device_polling_interval", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.idp_signoff", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.response_types.#", "3"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.response_types.*", "CODE"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.response_types.*", "TOKEN"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.response_types.*", "ID_TOKEN"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.token_endpoint_auth_method", "CLIENT_SECRET_BASIC"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.par_requirement", "REQUIRED"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.par_timeout", "180"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.pkce_enforcement", "REQUIRED"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.redirect_uris.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.redirect_uris.*", "https://pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.redirect_uris.*", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.allow_wildcard_in_redirect_uris", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.post_logout_redirect_uris.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.post_logout_redirect_uris.*", "https://pingidentity.com/logout"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.post_logout_redirect_uris.*", "https://www.pingidentity.com/logout"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.refresh_token_duration", "3000000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.refresh_token_rolling_duration", "30000000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.refresh_token_rolling_grace_period_duration", "80000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.additional_refresh_token_replay_protection_enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.support_unsigned_request_object", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.require_signed_request_object", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.cors_settings.behavior", "ALLOW_SPECIFIC_ORIGINS"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.cors_settings.origins.#", "8"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "http://localhost"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "https://localhost"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "http://auth.pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "https://auth.pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "http://*.pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "https://*.pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "http://192.168.1.1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "https://192.168.1.1"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.mobile_app"),
					resource.TestCheckNoResourceAttr(resourceFullName, "saml_options"),
					resource.TestCheckNoResourceAttr(resourceFullName, "external_link_options"),
					resource.TestCheckNoResourceAttr(resourceFullName, "wsfed_options"),
					resource.TestCheckResourceAttr(resourceFullName, "hidden_from_app_portal", "true"),
				),
			},
		},
	})
}

func TestAccApplication_OIDCCustom_Device(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application.%s", resourceName)

	name := resourceName

	testStepFull := resource.TestStep{
		Config: testAccApplicationConfig_OIDC_Custom_Device_Full(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "oidc_options.grant_types.#", "2"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.grant_types.*", "DEVICE_CODE"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.grant_types.*", "REFRESH_TOKEN"),
			resource.TestCheckResourceAttr(resourceFullName, "oidc_options.device_path_id", "mobileAppId-1"),
			resource.TestCheckResourceAttr(resourceFullName, "oidc_options.device_custom_verification_uri", "https://pingidentity.com/verification1"),
			resource.TestCheckResourceAttr(resourceFullName, "oidc_options.device_timeout", "500"),
			resource.TestCheckResourceAttr(resourceFullName, "oidc_options.device_polling_interval", "10"),
		),
	}

	testStepMinimal := resource.TestStep{
		Config: testAccApplicationConfig_OIDC_Custom_Device_Minimal(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "oidc_options.grant_types.#", "1"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.grant_types.*", "DEVICE_CODE"),
			resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.device_path_id"),
			resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.device_custom_verification_uri"),
			resource.TestCheckResourceAttr(resourceFullName, "oidc_options.device_timeout", "600"),
			resource.TestCheckResourceAttr(resourceFullName, "oidc_options.device_polling_interval", "5"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Application_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// With
			testStepFull,
			{
				Config:  testAccApplicationConfig_OIDC_Custom_Device_Full(resourceName, name),
				Destroy: true,
			},
			// Without
			testStepMinimal,
			{
				Config:  testAccApplicationConfig_OIDC_Custom_Device_Minimal(resourceName, name),
				Destroy: true,
			},
			// Change
			testStepFull,
			testStepMinimal,
			testStepFull,
			{
				Config:  testAccApplicationConfig_OIDC_Custom_Device_Full(resourceName, name),
				Destroy: true,
			},
		},
	})
}

func TestAccApplication_OIDCFullService(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application.%s", resourceName)

	name := resourceName

	data, _ := os.ReadFile("../../acctest/test_assets/image/image-logo.gif")
	image := base64.StdEncoding.EncodeToString(data)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Application_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationConfig_OIDC_FullService(resourceName, name, image),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", "My test OIDC app"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "login_page_url", "https://www.pingidentity.com"),
					resource.TestMatchResourceAttr(resourceFullName, "icon.id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "icon.href", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_role_type", "ADMIN_USERS_ONLY"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.groups.#", "1"),
					resource.TestMatchResourceAttr(resourceFullName, "access_control_group_options.groups.0", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.type", "ANY_GROUP"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.type", "SERVICE"),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.client_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.home_page_url", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.initiate_login_uri", "https://www.pingidentity.com/initiate"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.target_link_uri", "https://www.pingidentity.com/target"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.grant_types.#", "4"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.grant_types.*", "AUTHORIZATION_CODE"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.grant_types.*", "CLIENT_CREDENTIALS"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.grant_types.*", "IMPLICIT"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.grant_types.*", "REFRESH_TOKEN"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.response_types.#", "3"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.response_types.*", "CODE"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.response_types.*", "TOKEN"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.response_types.*", "ID_TOKEN"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.token_endpoint_auth_method", "CLIENT_SECRET_BASIC"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.par_requirement", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.par_timeout", "60"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.pkce_enforcement", "REQUIRED"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.redirect_uris.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.redirect_uris.*", "https://pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.redirect_uris.*", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.allow_wildcard_in_redirect_uris", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.post_logout_redirect_uris.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.post_logout_redirect_uris.*", "https://pingidentity.com/logout"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.post_logout_redirect_uris.*", "https://www.pingidentity.com/logout"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.refresh_token_duration", "3000000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.refresh_token_rolling_duration", "30000000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.refresh_token_rolling_grace_period_duration", "80000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.additional_refresh_token_replay_protection_enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.support_unsigned_request_object", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.require_signed_request_object", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.cors_settings.behavior", "ALLOW_SPECIFIC_ORIGINS"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.cors_settings.origins.#", "8"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "http://localhost"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "https://localhost"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "http://auth.pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "https://auth.pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "http://*.pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "https://*.pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "http://192.168.1.1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "https://192.168.1.1"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.mobile_app"),
					resource.TestCheckNoResourceAttr(resourceFullName, "saml_options"),
					resource.TestCheckNoResourceAttr(resourceFullName, "external_link_options"),
					resource.TestCheckNoResourceAttr(resourceFullName, "wsfed_options"),
					resource.TestCheckResourceAttr(resourceFullName, "hidden_from_app_portal", "true"),
				),
			},
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccApplication_OIDCMinimalService(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Application_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationConfig_OIDC_MinimalService(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckNoResourceAttr(resourceFullName, "description"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckNoResourceAttr(resourceFullName, "login_page_url"),
					resource.TestCheckNoResourceAttr(resourceFullName, "icon"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.type", "SERVICE"),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.client_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.home_page_url"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.initiate_login_uri"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.target_link_uri"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.grant_types.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.grant_types.*", "AUTHORIZATION_CODE"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.grant_types.*", "REFRESH_TOKEN"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.response_types.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.response_types.*", "CODE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.token_endpoint_auth_method", "CLIENT_SECRET_BASIC"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.par_requirement", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.par_timeout", "60"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.pkce_enforcement", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.redirect_uris.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.redirect_uris.*", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.allow_wildcard_in_redirect_uris", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.post_logout_redirect_uris.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.refresh_token_duration", "2592000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.refresh_token_rolling_duration", "15552000"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.refresh_token_rolling_grace_period_duration"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.additional_refresh_token_replay_protection_enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.support_unsigned_request_object", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.require_signed_request_object", "false"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.cors_settings"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.mobile_app"),
					resource.TestCheckNoResourceAttr(resourceFullName, "saml_options"),
					resource.TestCheckNoResourceAttr(resourceFullName, "external_link_options"),
					resource.TestCheckNoResourceAttr(resourceFullName, "wsfed_options"),
					resource.TestCheckResourceAttr(resourceFullName, "hidden_from_app_portal", "false"),
				),
			},
		},
	})
}

func TestAccApplication_OIDCServiceUpdate(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application.%s", resourceName)

	name := resourceName

	data, _ := os.ReadFile("../../acctest/test_assets/image/image-logo.gif")
	image := base64.StdEncoding.EncodeToString(data)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Application_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationConfig_OIDC_FullService(resourceName, name, image),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", "My test OIDC app"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "login_page_url", "https://www.pingidentity.com"),
					resource.TestMatchResourceAttr(resourceFullName, "icon.id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "icon.href", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_role_type", "ADMIN_USERS_ONLY"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.groups.#", "1"),
					resource.TestMatchResourceAttr(resourceFullName, "access_control_group_options.groups.0", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.type", "ANY_GROUP"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.type", "SERVICE"),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.client_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.home_page_url", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.initiate_login_uri", "https://www.pingidentity.com/initiate"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.target_link_uri", "https://www.pingidentity.com/target"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.grant_types.#", "4"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.grant_types.*", "AUTHORIZATION_CODE"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.grant_types.*", "CLIENT_CREDENTIALS"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.grant_types.*", "IMPLICIT"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.grant_types.*", "REFRESH_TOKEN"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.response_types.#", "3"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.response_types.*", "CODE"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.response_types.*", "TOKEN"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.response_types.*", "ID_TOKEN"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.token_endpoint_auth_method", "CLIENT_SECRET_BASIC"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.par_requirement", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.par_timeout", "60"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.pkce_enforcement", "REQUIRED"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.redirect_uris.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.redirect_uris.*", "https://pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.redirect_uris.*", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.allow_wildcard_in_redirect_uris", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.post_logout_redirect_uris.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.post_logout_redirect_uris.*", "https://pingidentity.com/logout"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.post_logout_redirect_uris.*", "https://www.pingidentity.com/logout"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.refresh_token_duration", "3000000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.refresh_token_rolling_duration", "30000000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.refresh_token_rolling_grace_period_duration", "80000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.additional_refresh_token_replay_protection_enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.support_unsigned_request_object", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.require_signed_request_object", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.cors_settings.behavior", "ALLOW_SPECIFIC_ORIGINS"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.cors_settings.origins.#", "8"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "http://localhost"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "https://localhost"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "http://auth.pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "https://auth.pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "http://*.pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "https://*.pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "http://192.168.1.1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "https://192.168.1.1"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.mobile_app"),
					resource.TestCheckNoResourceAttr(resourceFullName, "saml_options"),
					resource.TestCheckNoResourceAttr(resourceFullName, "external_link_options"),
					resource.TestCheckNoResourceAttr(resourceFullName, "wsfed_options"),
					resource.TestCheckResourceAttr(resourceFullName, "hidden_from_app_portal", "true"),
				),
			},
			{
				Config: testAccApplicationConfig_OIDC_MinimalService(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckNoResourceAttr(resourceFullName, "description"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckNoResourceAttr(resourceFullName, "login_page_url"),
					resource.TestCheckNoResourceAttr(resourceFullName, "icon"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.type", "SERVICE"),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.client_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.home_page_url"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.initiate_login_uri"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.target_link_uri"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.grant_types.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.grant_types.*", "AUTHORIZATION_CODE"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.grant_types.*", "REFRESH_TOKEN"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.response_types.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.response_types.*", "CODE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.token_endpoint_auth_method", "CLIENT_SECRET_BASIC"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.par_requirement", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.par_timeout", "60"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.pkce_enforcement", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.redirect_uris.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.redirect_uris.*", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.allow_wildcard_in_redirect_uris", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.post_logout_redirect_uris.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.refresh_token_duration", "2592000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.refresh_token_rolling_duration", "15552000"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.refresh_token_rolling_grace_period_duration"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.additional_refresh_token_replay_protection_enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.support_unsigned_request_object", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.require_signed_request_object", "false"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.cors_settings"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.mobile_app"),
					resource.TestCheckNoResourceAttr(resourceFullName, "saml_options"),
					resource.TestCheckNoResourceAttr(resourceFullName, "external_link_options"),
					resource.TestCheckNoResourceAttr(resourceFullName, "wsfed_options"),
					resource.TestCheckResourceAttr(resourceFullName, "hidden_from_app_portal", "false"),
				),
			},
			{
				Config: testAccApplicationConfig_OIDC_FullService(resourceName, name, image),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", "My test OIDC app"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "login_page_url", "https://www.pingidentity.com"),
					resource.TestMatchResourceAttr(resourceFullName, "icon.id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "icon.href", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_role_type", "ADMIN_USERS_ONLY"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.groups.#", "1"),
					resource.TestMatchResourceAttr(resourceFullName, "access_control_group_options.groups.0", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.type", "ANY_GROUP"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.type", "SERVICE"),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.client_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.home_page_url", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.initiate_login_uri", "https://www.pingidentity.com/initiate"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.target_link_uri", "https://www.pingidentity.com/target"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.grant_types.#", "4"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.grant_types.*", "AUTHORIZATION_CODE"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.grant_types.*", "CLIENT_CREDENTIALS"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.grant_types.*", "IMPLICIT"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.grant_types.*", "REFRESH_TOKEN"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.response_types.#", "3"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.response_types.*", "CODE"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.response_types.*", "TOKEN"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.response_types.*", "ID_TOKEN"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.token_endpoint_auth_method", "CLIENT_SECRET_BASIC"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.par_requirement", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.par_timeout", "60"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.pkce_enforcement", "REQUIRED"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.redirect_uris.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.redirect_uris.*", "https://pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.redirect_uris.*", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.allow_wildcard_in_redirect_uris", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.post_logout_redirect_uris.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.post_logout_redirect_uris.*", "https://pingidentity.com/logout"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.post_logout_redirect_uris.*", "https://www.pingidentity.com/logout"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.refresh_token_duration", "3000000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.refresh_token_rolling_duration", "30000000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.refresh_token_rolling_grace_period_duration", "80000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.additional_refresh_token_replay_protection_enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.support_unsigned_request_object", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.require_signed_request_object", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.cors_settings.behavior", "ALLOW_SPECIFIC_ORIGINS"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.cors_settings.origins.#", "8"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "http://localhost"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "https://localhost"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "http://auth.pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "https://auth.pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "http://*.pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "https://*.pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "http://192.168.1.1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "https://192.168.1.1"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.mobile_app"),
					resource.TestCheckNoResourceAttr(resourceFullName, "saml_options"),
					resource.TestCheckNoResourceAttr(resourceFullName, "external_link_options"),
					resource.TestCheckNoResourceAttr(resourceFullName, "wsfed_options"),
					resource.TestCheckResourceAttr(resourceFullName, "hidden_from_app_portal", "true"),
				),
			},
		},
	})
}

func TestAccApplication_OIDCFullSPA(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application.%s", resourceName)

	name := resourceName

	data, _ := os.ReadFile("../../acctest/test_assets/image/image-logo.gif")
	image := base64.StdEncoding.EncodeToString(data)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Application_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationConfig_OIDC_FullSPA(resourceName, name, image),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", "My test OIDC app"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "login_page_url", "https://www.pingidentity.com"),
					resource.TestMatchResourceAttr(resourceFullName, "icon.id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "icon.href", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_role_type", "ADMIN_USERS_ONLY"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.groups.#", "1"),
					resource.TestMatchResourceAttr(resourceFullName, "access_control_group_options.groups.0", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.type", "ANY_GROUP"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.type", "SINGLE_PAGE_APP"),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.client_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.home_page_url", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.initiate_login_uri", "https://www.pingidentity.com/initiate"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.target_link_uri", "https://www.pingidentity.com/target"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.grant_types.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.grant_types.*", "AUTHORIZATION_CODE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.response_types.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.response_types.*", "CODE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.token_endpoint_auth_method", "NONE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.par_requirement", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.par_timeout", "60"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.pkce_enforcement", "S256_REQUIRED"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.idp_signoff", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.redirect_uris.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.redirect_uris.*", "https://pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.redirect_uris.*", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.allow_wildcard_in_redirect_uris", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.post_logout_redirect_uris.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.post_logout_redirect_uris.*", "https://pingidentity.com/logout"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.post_logout_redirect_uris.*", "https://www.pingidentity.com/logout"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.refresh_token_duration", "2592000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.refresh_token_rolling_duration", "15552000"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.refresh_token_rolling_grace_period_duration"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.additional_refresh_token_replay_protection_enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.support_unsigned_request_object", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.require_signed_request_object", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.cors_settings.behavior", "ALLOW_SPECIFIC_ORIGINS"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.cors_settings.origins.#", "8"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "http://localhost"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "https://localhost"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "http://auth.pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "https://auth.pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "http://*.pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "https://*.pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "http://192.168.1.1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "https://192.168.1.1"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.mobile_app"),
					resource.TestCheckNoResourceAttr(resourceFullName, "saml_options"),
					resource.TestCheckNoResourceAttr(resourceFullName, "external_link_options"),
					resource.TestCheckNoResourceAttr(resourceFullName, "wsfed_options"),
					resource.TestCheckResourceAttr(resourceFullName, "hidden_from_app_portal", "true"),
				),
			},
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccApplication_OIDCMinimalSPA(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Application_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationConfig_OIDC_MinimalSPA(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckNoResourceAttr(resourceFullName, "description"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckNoResourceAttr(resourceFullName, "login_page_url"),
					resource.TestCheckNoResourceAttr(resourceFullName, "icon"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.type", "SINGLE_PAGE_APP"),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.client_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.home_page_url"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.initiate_login_uri"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.target_link_uri"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.grant_types.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.grant_types.*", "AUTHORIZATION_CODE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.response_types.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.response_types.*", "CODE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.token_endpoint_auth_method", "NONE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.par_requirement", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.par_timeout", "60"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.pkce_enforcement", "S256_REQUIRED"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.redirect_uris.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.redirect_uris.*", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.allow_wildcard_in_redirect_uris", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.post_logout_redirect_uris.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.refresh_token_duration", "2592000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.refresh_token_rolling_duration", "15552000"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.refresh_token_rolling_grace_period_duration"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.additional_refresh_token_replay_protection_enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.support_unsigned_request_object", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.require_signed_request_object", "false"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.cors_settings"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.mobile_app"),
					resource.TestCheckNoResourceAttr(resourceFullName, "saml_options"),
					resource.TestCheckNoResourceAttr(resourceFullName, "external_link_options"),
					resource.TestCheckNoResourceAttr(resourceFullName, "wsfed_options"),
					resource.TestCheckResourceAttr(resourceFullName, "hidden_from_app_portal", "false"),
				),
			},
		},
	})
}

func TestAccApplication_OIDCSPAUpdate(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application.%s", resourceName)

	name := resourceName

	data, _ := os.ReadFile("../../acctest/test_assets/image/image-logo.gif")
	image := base64.StdEncoding.EncodeToString(data)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Application_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationConfig_OIDC_MinimalSPA(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckNoResourceAttr(resourceFullName, "description"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckNoResourceAttr(resourceFullName, "login_page_url"),
					resource.TestCheckNoResourceAttr(resourceFullName, "icon"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.type", "SINGLE_PAGE_APP"),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.client_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.home_page_url"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.initiate_login_uri"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.target_link_uri"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.grant_types.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.grant_types.*", "AUTHORIZATION_CODE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.response_types.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.response_types.*", "CODE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.token_endpoint_auth_method", "NONE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.par_requirement", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.par_timeout", "60"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.pkce_enforcement", "S256_REQUIRED"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.redirect_uris.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.redirect_uris.*", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.allow_wildcard_in_redirect_uris", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.post_logout_redirect_uris.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.refresh_token_duration", "2592000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.refresh_token_rolling_duration", "15552000"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.refresh_token_rolling_grace_period_duration"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.additional_refresh_token_replay_protection_enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.support_unsigned_request_object", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.require_signed_request_object", "false"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.cors_settings"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.mobile_app"),
					resource.TestCheckNoResourceAttr(resourceFullName, "saml_options"),
					resource.TestCheckNoResourceAttr(resourceFullName, "external_link_options"),
					resource.TestCheckNoResourceAttr(resourceFullName, "wsfed_options"),
					resource.TestCheckResourceAttr(resourceFullName, "hidden_from_app_portal", "false"),
				),
			},
			{
				Config: testAccApplicationConfig_OIDC_FullSPA(resourceName, name, image),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", "My test OIDC app"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "login_page_url", "https://www.pingidentity.com"),
					resource.TestMatchResourceAttr(resourceFullName, "icon.id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "icon.href", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_role_type", "ADMIN_USERS_ONLY"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.groups.#", "1"),
					resource.TestMatchResourceAttr(resourceFullName, "access_control_group_options.groups.0", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.type", "ANY_GROUP"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.type", "SINGLE_PAGE_APP"),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.client_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.home_page_url", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.initiate_login_uri", "https://www.pingidentity.com/initiate"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.target_link_uri", "https://www.pingidentity.com/target"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.grant_types.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.grant_types.*", "AUTHORIZATION_CODE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.response_types.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.response_types.*", "CODE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.token_endpoint_auth_method", "NONE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.par_requirement", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.par_timeout", "60"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.pkce_enforcement", "S256_REQUIRED"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.idp_signoff", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.redirect_uris.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.redirect_uris.*", "https://pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.redirect_uris.*", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.allow_wildcard_in_redirect_uris", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.post_logout_redirect_uris.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.post_logout_redirect_uris.*", "https://pingidentity.com/logout"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.post_logout_redirect_uris.*", "https://www.pingidentity.com/logout"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.refresh_token_duration", "2592000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.refresh_token_rolling_duration", "15552000"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.refresh_token_rolling_grace_period_duration"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.additional_refresh_token_replay_protection_enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.support_unsigned_request_object", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.require_signed_request_object", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.cors_settings.behavior", "ALLOW_SPECIFIC_ORIGINS"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.cors_settings.origins.#", "8"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "http://localhost"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "https://localhost"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "http://auth.pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "https://auth.pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "http://*.pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "https://*.pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "http://192.168.1.1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "https://192.168.1.1"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.mobile_app"),
					resource.TestCheckNoResourceAttr(resourceFullName, "saml_options"),
					resource.TestCheckNoResourceAttr(resourceFullName, "external_link_options"),
					resource.TestCheckNoResourceAttr(resourceFullName, "wsfed_options"),
					resource.TestCheckResourceAttr(resourceFullName, "hidden_from_app_portal", "true"),
				),
			},
			{
				Config: testAccApplicationConfig_OIDC_MinimalSPA(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckNoResourceAttr(resourceFullName, "description"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckNoResourceAttr(resourceFullName, "login_page_url"),
					resource.TestCheckNoResourceAttr(resourceFullName, "icon"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.type", "SINGLE_PAGE_APP"),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.client_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.home_page_url"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.initiate_login_uri"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.target_link_uri"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.grant_types.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.grant_types.*", "AUTHORIZATION_CODE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.response_types.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.response_types.*", "CODE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.token_endpoint_auth_method", "NONE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.par_requirement", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.par_timeout", "60"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.pkce_enforcement", "S256_REQUIRED"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.redirect_uris.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.redirect_uris.*", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.allow_wildcard_in_redirect_uris", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.post_logout_redirect_uris.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.refresh_token_duration", "2592000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.refresh_token_rolling_duration", "15552000"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.refresh_token_rolling_grace_period_duration"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.additional_refresh_token_replay_protection_enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.support_unsigned_request_object", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.require_signed_request_object", "false"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.cors_settings"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.mobile_app"),
					resource.TestCheckNoResourceAttr(resourceFullName, "saml_options"),
					resource.TestCheckNoResourceAttr(resourceFullName, "external_link_options"),
					resource.TestCheckNoResourceAttr(resourceFullName, "wsfed_options"),
					resource.TestCheckResourceAttr(resourceFullName, "hidden_from_app_portal", "false"),
				),
			},
		},
	})
}

func TestAccApplication_OIDCFullWorker(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application.%s", resourceName)

	name := resourceName

	data, _ := os.ReadFile("../../acctest/test_assets/image/image-logo.png")
	image := base64.StdEncoding.EncodeToString(data)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Application_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationConfig_OIDC_FullWorker(resourceName, name, image),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", "My test OIDC app"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "login_page_url", "https://www.pingidentity.com"),
					resource.TestMatchResourceAttr(resourceFullName, "icon.id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "icon.href", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_role_type", "ADMIN_USERS_ONLY"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.groups.#", "2"),
					resource.TestMatchResourceAttr(resourceFullName, "access_control_group_options.groups.0", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "access_control_group_options.groups.1", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.type", "ANY_GROUP"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.type", "WORKER"),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.client_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.home_page_url"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.initiate_login_uri"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.target_link_uri"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.grant_types.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.grant_types.*", "CLIENT_CREDENTIALS"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.response_types.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.token_endpoint_auth_method", "CLIENT_SECRET_BASIC"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.par_requirement", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.par_timeout", "60"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.pkce_enforcement", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.redirect_uris.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.allow_wildcard_in_redirect_uris", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.post_logout_redirect_uris.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.refresh_token_duration", "2592000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.refresh_token_rolling_duration", "15552000"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.refresh_token_rolling_grace_period_duration"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.additional_refresh_token_replay_protection_enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.support_unsigned_request_object", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.require_signed_request_object", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.cors_settings.behavior", "ALLOW_SPECIFIC_ORIGINS"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.cors_settings.origins.#", "8"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "http://localhost"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "https://localhost"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "http://auth.pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "https://auth.pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "http://*.pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "https://*.pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "http://192.168.1.1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "https://192.168.1.1"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.mobile_app"),
					resource.TestCheckNoResourceAttr(resourceFullName, "saml_options"),
					resource.TestCheckNoResourceAttr(resourceFullName, "external_link_options"),
					resource.TestCheckNoResourceAttr(resourceFullName, "wsfed_options"),
					resource.TestCheckResourceAttr(resourceFullName, "hidden_from_app_portal", "true"),
				),
			},
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccApplication_OIDCMinimalWorker(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Application_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationConfig_OIDC_MinimalWorker(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckNoResourceAttr(resourceFullName, "description"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckNoResourceAttr(resourceFullName, "login_page_url"),
					resource.TestCheckNoResourceAttr(resourceFullName, "icon"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_role_type", "ADMIN_USERS_ONLY"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.type", "WORKER"),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.client_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.home_page_url"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.initiate_login_uri"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.target_link_uri"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.grant_types.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.grant_types.*", "CLIENT_CREDENTIALS"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.response_types.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.token_endpoint_auth_method", "CLIENT_SECRET_BASIC"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.par_requirement", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.par_timeout", "60"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.pkce_enforcement", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.redirect_uris.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.allow_wildcard_in_redirect_uris", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.post_logout_redirect_uris.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.refresh_token_duration", "2592000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.refresh_token_rolling_duration", "15552000"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.refresh_token_rolling_grace_period_duration"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.additional_refresh_token_replay_protection_enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.support_unsigned_request_object", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.require_signed_request_object", "false"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.cors_settings"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.mobile_app"),
					resource.TestCheckNoResourceAttr(resourceFullName, "saml_options"),
					resource.TestCheckNoResourceAttr(resourceFullName, "external_link_options"),
					resource.TestCheckNoResourceAttr(resourceFullName, "wsfed_options"),
					resource.TestCheckResourceAttr(resourceFullName, "hidden_from_app_portal", "false"),
				),
			},
		},
	})
}

func TestAccApplication_OIDCWorkerUpdate(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application.%s", resourceName)

	name := resourceName

	data, _ := os.ReadFile("../../acctest/test_assets/image/image-logo.png")
	image := base64.StdEncoding.EncodeToString(data)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Application_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationConfig_OIDC_MinimalWorker(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckNoResourceAttr(resourceFullName, "description"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckNoResourceAttr(resourceFullName, "login_page_url"),
					resource.TestCheckNoResourceAttr(resourceFullName, "icon"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_role_type", "ADMIN_USERS_ONLY"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.type", "WORKER"),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.client_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.home_page_url"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.initiate_login_uri"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.target_link_uri"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.grant_types.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.grant_types.*", "CLIENT_CREDENTIALS"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.response_types.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.token_endpoint_auth_method", "CLIENT_SECRET_BASIC"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.par_requirement", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.par_timeout", "60"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.pkce_enforcement", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.redirect_uris.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.allow_wildcard_in_redirect_uris", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.post_logout_redirect_uris.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.refresh_token_duration", "2592000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.refresh_token_rolling_duration", "15552000"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.refresh_token_rolling_grace_period_duration"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.additional_refresh_token_replay_protection_enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.support_unsigned_request_object", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.require_signed_request_object", "false"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.cors_settings"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.mobile_app"),
					resource.TestCheckNoResourceAttr(resourceFullName, "saml_options"),
					resource.TestCheckNoResourceAttr(resourceFullName, "external_link_options"),
					resource.TestCheckNoResourceAttr(resourceFullName, "wsfed_options"),
					resource.TestCheckResourceAttr(resourceFullName, "hidden_from_app_portal", "false"),
				),
			},
			{
				Config: testAccApplicationConfig_OIDC_FullWorker(resourceName, name, image),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", "My test OIDC app"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "login_page_url", "https://www.pingidentity.com"),
					resource.TestMatchResourceAttr(resourceFullName, "icon.id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "icon.href", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_role_type", "ADMIN_USERS_ONLY"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.groups.#", "2"),
					resource.TestMatchResourceAttr(resourceFullName, "access_control_group_options.groups.0", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "access_control_group_options.groups.1", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.type", "ANY_GROUP"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.type", "WORKER"),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.client_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.home_page_url"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.initiate_login_uri"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.target_link_uri"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.grant_types.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.grant_types.*", "CLIENT_CREDENTIALS"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.response_types.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.token_endpoint_auth_method", "CLIENT_SECRET_BASIC"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.par_requirement", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.par_timeout", "60"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.pkce_enforcement", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.redirect_uris.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.allow_wildcard_in_redirect_uris", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.post_logout_redirect_uris.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.refresh_token_duration", "2592000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.refresh_token_rolling_duration", "15552000"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.refresh_token_rolling_grace_period_duration"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.additional_refresh_token_replay_protection_enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.support_unsigned_request_object", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.require_signed_request_object", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.cors_settings.behavior", "ALLOW_SPECIFIC_ORIGINS"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.cors_settings.origins.#", "8"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "http://localhost"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "https://localhost"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "http://auth.pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "https://auth.pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "http://*.pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "https://*.pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "http://192.168.1.1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.cors_settings.origins.*", "https://192.168.1.1"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.mobile_app"),
					resource.TestCheckNoResourceAttr(resourceFullName, "saml_options"),
					resource.TestCheckNoResourceAttr(resourceFullName, "external_link_options"),
					resource.TestCheckNoResourceAttr(resourceFullName, "wsfed_options"),
					resource.TestCheckResourceAttr(resourceFullName, "hidden_from_app_portal", "true"),
				),
			},
			{
				Config: testAccApplicationConfig_OIDC_MinimalWorker(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckNoResourceAttr(resourceFullName, "description"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckNoResourceAttr(resourceFullName, "login_page_url"),
					resource.TestCheckNoResourceAttr(resourceFullName, "icon"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_role_type", "ADMIN_USERS_ONLY"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.type", "WORKER"),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.client_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.home_page_url"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.initiate_login_uri"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.target_link_uri"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.grant_types.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.grant_types.*", "CLIENT_CREDENTIALS"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.response_types.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.token_endpoint_auth_method", "CLIENT_SECRET_BASIC"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.par_requirement", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.par_timeout", "60"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.pkce_enforcement", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.redirect_uris.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.allow_wildcard_in_redirect_uris", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.post_logout_redirect_uris.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.refresh_token_duration", "2592000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.refresh_token_rolling_duration", "15552000"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.refresh_token_rolling_grace_period_duration"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.additional_refresh_token_replay_protection_enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.support_unsigned_request_object", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.require_signed_request_object", "false"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.cors_settings"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.mobile_app"),
					resource.TestCheckNoResourceAttr(resourceFullName, "saml_options"),
					resource.TestCheckNoResourceAttr(resourceFullName, "external_link_options"),
					resource.TestCheckNoResourceAttr(resourceFullName, "wsfed_options"),
					resource.TestCheckResourceAttr(resourceFullName, "hidden_from_app_portal", "false"),
				),
			},
		},
	})
}

// OIDC Use Cases
func TestAccApplication_OIDC_WildcardInRedirectURI(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Application_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccApplicationConfig_OIDC_WildcardInRedirect(resourceName, name, false),
				ExpectError: regexp.MustCompile("Invalid configuration"),
			},
			{
				Config: testAccApplicationConfig_OIDC_WildcardInRedirect(resourceName, name, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.redirect_uris.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.redirect_uris.*", "https://www.pingidentity.com/*"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.allow_wildcard_in_redirect_uris", "true"),
				),
			},
		},
	})
}

type OIDCLocalhostTest struct {
	Hostname string
	Valid    string
}

func TestAccApplication_OIDC_LocalhostAddresses(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Application_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Localhost
			{
				Config: testAccApplicationConfig_OIDC_LocalhostAddresses(resourceName, name, "localhost"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "login_page_url", "http://localhost/login"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.home_page_url", "http://localhost/home"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.initiate_login_uri", "http://localhost/init"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.redirect_uris.*", "http://localhost/callback"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.post_logout_redirect_uris.*", "http://localhost/logout"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.target_link_uri", "http://localhost/link"),
				),
			},
			{
				Config:  testAccApplicationConfig_OIDC_LocalhostAddresses(resourceName, name, "localhost"),
				Destroy: true,
			},
			// 127.0.0.1
			{
				Config: testAccApplicationConfig_OIDC_LocalhostAddresses(resourceName, name, "127.0.0.1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "login_page_url", "http://127.0.0.1/login"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.home_page_url", "http://127.0.0.1/home"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.initiate_login_uri", "http://127.0.0.1/init"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.redirect_uris.*", "http://127.0.0.1/callback"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.post_logout_redirect_uris.*", "http://127.0.0.1/logout"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.target_link_uri", "http://127.0.0.1/link"),
				),
			},
			{
				Config:  testAccApplicationConfig_OIDC_LocalhostAddresses(resourceName, name, "127.0.0.1"),
				Destroy: true,
			},
		},
	})
}

func TestAccApplication_OIDC_NativeAppAddresses(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Application_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Localhost
			{
				Config: testAccApplicationConfig_OIDC_NativeAppAddresses(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.redirect_uris.*", "com.myapp.app://callback"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.post_logout_redirect_uris.*", "com.myapp.app://logout"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.target_link_uri", "com.myapp.app://target"),
				),
			},
		},
	})
}

func TestAccApplication_OIDC_JwtTokenAuth(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application.%s", resourceName)

	name := resourceName

	clientSecretBasic := resource.TestStep{
		Config: testAccApplicationConfig_OIDC_MinimalWeb(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "oidc_options.token_endpoint_auth_method", "CLIENT_SECRET_BASIC"),
			resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.jwks"),
			resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.jwks_url"),
		),
	}

	privateKeyJwtJWKS := resource.TestStep{
		Config: testAccApplicationConfig_OIDC_MinimalWeb_PrivateKeyJWT_JWKS(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "oidc_options.token_endpoint_auth_method", "PRIVATE_KEY_JWT"),
			resource.TestCheckResourceAttr(resourceFullName, "oidc_options.jwks", "{\n\t\"keys\": [\n\t  {\n\t\t\"kty\": \"RSA\",\n\t\t\"e\": \"AQAB\",\n\t\t\"use\": \"sig\",\n\t\t\"kid\": \"12345\",\n\t\t\"alg\": \"RS256\",\n\t\t\"n\": \"0vx7agoebGcQSuuPiLJXZptN9nndrQmbXEPCR0VH7jhV1JvKFvVsenY4rz5BnCNRS7U2mFF9K2BWXTZiaF4f3hjd4J0AOnHZV9KbV7L5Cp-1PEXF12R3HCkcqsn9b2I0xEs5OYLQaGbhV8v1FwTD2jZX0tAiw4i+SN9GJkPV2ZCOnF-8RPZCVDG9LZGFq4c9-YNPvRwT7B9-EN0kDYEKsOmGiJ0PVPAVTPZ9EVZVdLg7SwamytKcP4fz_BLTjCojz2W9KIL5UZGenQR5S7KAZxJ0T0DO8Q4kqNVdF7OOrBizX6-qQ9ZC1l6HJ6Sq9ye0oW2jTiMlNQxgc5vNRgHFAmb4DNA\"\n\t  }\n\t]\n}\n"),
			resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.jwks_url"),
		),
	}

	privateKeyJwtJWKSURL := resource.TestStep{
		Config: testAccApplicationConfig_OIDC_MinimalWeb_PrivateKeyJWT_JWKSURL(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "oidc_options.token_endpoint_auth_method", "PRIVATE_KEY_JWT"),
			resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.jwks"),
			resource.TestCheckResourceAttr(resourceFullName, "oidc_options.jwks_url", "https://pingidentity.com/jwks"),
		),
	}

	clientSecretJWT := resource.TestStep{
		Config: testAccApplicationConfig_OIDC_MinimalWeb_ClientSecretJWT(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "oidc_options.token_endpoint_auth_method", "CLIENT_SECRET_JWT"),
			resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.jwks"),
			resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.jwks_url"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Application_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Private key jwt jwks
			privateKeyJwtJWKS,
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:  testAccApplicationConfig_OIDC_MinimalWeb_PrivateKeyJWT_JWKS(resourceName, name),
				Destroy: true,
			},
			// Private key jwt jwks_url
			privateKeyJwtJWKSURL,
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:  testAccApplicationConfig_OIDC_MinimalWeb_PrivateKeyJWT_JWKSURL(resourceName, name),
				Destroy: true,
			},
			// Client secret jwt
			clientSecretJWT,
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:  testAccApplicationConfig_OIDC_MinimalWeb_ClientSecretJWT(resourceName, name),
				Destroy: true,
			},
			// Update
			clientSecretBasic,
			privateKeyJwtJWKS,
			privateKeyJwtJWKSURL,
			clientSecretJWT,
			clientSecretBasic,
		},
	})
}

// SAML
func TestAccApplication_SAMLFull(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application.%s", resourceName)

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
			acctest.PreCheckPKCS7Cert(t)
			acctest.PreCheckPEMCert(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Application_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationConfig_SAML_Full(environmentName, licenseID, resourceName, name, image, pkcs7_cert, pem_cert),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", "My test SAML app"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "login_page_url", "https://www.pingidentity.com"),
					resource.TestMatchResourceAttr(resourceFullName, "icon.id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "icon.href", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_role_type", "ADMIN_USERS_ONLY"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.groups.#", "2"),
					resource.TestMatchResourceAttr(resourceFullName, "access_control_group_options.groups.0", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "access_control_group_options.groups.1", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.type", "ANY_GROUP"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options"),
					resource.TestCheckNoResourceAttr(resourceFullName, "wsfed_options"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.home_page_url", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.type", "WEB_APP"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.acs_urls.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "saml_options.acs_urls.*", "https://pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "saml_options.acs_urls.*", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.assertion_duration", "3600"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.assertion_signed_enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.default_target_url", "https://www.pingidentity.com/relaystate"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.idp_signing_key.algorithm", "SHA384withECDSA"),
					resource.TestMatchResourceAttr(resourceFullName, "saml_options.idp_signing_key.key_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.enable_requested_authn_context", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.nameid_format", "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.response_is_signed", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.session_not_on_or_after_duration", "64"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.slo_binding", "HTTP_REDIRECT"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.slo_endpoint", "https://www.pingidentity.com/sloendpoint"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.slo_response_endpoint", "https://www.pingidentity.com/sloresponseendpoint"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.slo_window", "3"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.sp_encryption.algorithm", "AES_256"),
					resource.TestMatchResourceAttr(resourceFullName, "saml_options.sp_encryption.certificate.id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.sp_entity_id", fmt.Sprintf("sp:entity:%s", resourceName)),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.sp_verification.authn_request_signed", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.sp_verification.certificate_ids.#", "2"),
					resource.TestMatchResourceAttr(resourceFullName, "saml_options.sp_verification.certificate_ids.0", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "saml_options.sp_verification.certificate_ids.1", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.cors_settings.behavior", "ALLOW_SPECIFIC_ORIGINS"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.virtual_server_id_settings.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.virtual_server_id_settings.virtual_server_ids.#", "4"),
					// Default goes to the last one in the list
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.virtual_server_id_settings.virtual_server_ids.3.vs_id", "virtualserver1"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.virtual_server_id_settings.virtual_server_ids.3.default", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.virtual_server_id_settings.virtual_server_ids.0.vs_id", "virtualserver2"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.virtual_server_id_settings.virtual_server_ids.0.default", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.virtual_server_id_settings.virtual_server_ids.1.vs_id", "virtualserver3"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.virtual_server_id_settings.virtual_server_ids.1.default", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.virtual_server_id_settings.virtual_server_ids.2.vs_id", "virtualserver4"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.virtual_server_id_settings.virtual_server_ids.2.default", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.cors_settings.origins.#", "8"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "saml_options.cors_settings.origins.*", "http://localhost"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "saml_options.cors_settings.origins.*", "https://localhost"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "saml_options.cors_settings.origins.*", "http://auth.pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "saml_options.cors_settings.origins.*", "https://auth.pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "saml_options.cors_settings.origins.*", "http://*.pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "saml_options.cors_settings.origins.*", "https://*.pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "saml_options.cors_settings.origins.*", "http://192.168.1.1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "saml_options.cors_settings.origins.*", "https://192.168.1.1"),
					resource.TestCheckNoResourceAttr(resourceFullName, "external_link_options"),
					resource.TestCheckResourceAttr(resourceFullName, "hidden_from_app_portal", "true"),
				),
			},
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccApplication_SAMLMinimal(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application.%s", resourceName)

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
			acctest.PreCheckPKCS7Cert(t)
			acctest.PreCheckPEMCert(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Application_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationConfig_SAML_Minimal(environmentName, licenseID, resourceName, name, image, pkcs7_cert, pem_cert),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckNoResourceAttr(resourceFullName, "description"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckNoResourceAttr(resourceFullName, "login_page_url"),
					resource.TestCheckNoResourceAttr(resourceFullName, "icon"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options"),
					resource.TestCheckNoResourceAttr(resourceFullName, "wsfed_options"),
					resource.TestCheckNoResourceAttr(resourceFullName, "saml_options.home_page_url"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.type", "WEB_APP"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.acs_urls.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "saml_options.acs_urls.*", "https://pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.assertion_duration", "3600"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.assertion_signed_enabled", "true"),
					resource.TestCheckNoResourceAttr(resourceFullName, "saml_options.default_target_url"),
					resource.TestMatchResourceAttr(resourceFullName, "saml_options.idp_signing_key.key_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckNoResourceAttr(resourceFullName, "saml_options.enable_requested_authn_context"),
					resource.TestCheckNoResourceAttr(resourceFullName, "saml_options.nameid_format"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.response_is_signed", "false"),
					resource.TestCheckNoResourceAttr(resourceFullName, "saml_options.session_not_on_or_after_duration"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.slo_binding", "HTTP_POST"),
					resource.TestCheckNoResourceAttr(resourceFullName, "saml_options.slo_endpoint"),
					resource.TestCheckNoResourceAttr(resourceFullName, "saml_options.slo_response_endpoint"),
					resource.TestCheckNoResourceAttr(resourceFullName, "saml_options.slo_window"),
					resource.TestCheckNoResourceAttr(resourceFullName, "saml_options.sp_encryption"),
					resource.TestCheckNoResourceAttr(resourceFullName, "saml_options.virtual_server_id_settings"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.sp_entity_id", fmt.Sprintf("sp:entity:%s", resourceName)),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.sp_verification.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.cors_settings.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "hidden_from_app_portal", "false"),
				),
			},
		},
	})
}

func TestAccApplication_SAMLVirtualServerIdSettingsOrdering(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application.%s", resourceName)

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
			acctest.PreCheckPKCS7Cert(t)
			acctest.PreCheckPEMCert(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Application_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Initial configuration
			{
				Config: testAccApplicationConfig_SAML_VirtualServerIdSettingsOrdering(environmentName, licenseID, resourceName, name, image, pkcs7_cert, pem_cert),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.virtual_server_id_settings.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.virtual_server_id_settings.virtual_server_ids.#", "4"),
					// Ensure the default virtual server ID is returned last
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.virtual_server_id_settings.virtual_server_ids.3.vs_id", "virtualserver1"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.virtual_server_id_settings.virtual_server_ids.3.default", "true"),
				),
			},
			// Change the order of virtual server IDs
			{
				Config: testAccApplicationConfig_SAML_VirtualServerIdSettingsReordered(environmentName, licenseID, resourceName, name, image, pkcs7_cert, pem_cert),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.virtual_server_id_settings.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.virtual_server_id_settings.virtual_server_ids.#", "4"),
					// Ensure the new default virtual server ID is returned last
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.virtual_server_id_settings.virtual_server_ids.3.vs_id", "virtualserver3"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.virtual_server_id_settings.virtual_server_ids.3.default", "true"),
				),
			},
		},
	})
}

func TestAccApplication_ExternalLinkFull(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application.%s", resourceName)

	name := resourceName

	data, _ := os.ReadFile("../../acctest/test_assets/image/image-logo.gif")
	image := base64.StdEncoding.EncodeToString(data)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Application_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationConfig_ExternalLinkFull(resourceName, name, image),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", "My test external link app"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestMatchResourceAttr(resourceFullName, "icon.id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "icon.href", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.groups.#", "2"),
					resource.TestMatchResourceAttr(resourceFullName, "access_control_group_options.groups.0", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "access_control_group_options.groups.1", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.type", "ANY_GROUP"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options"),
					resource.TestCheckNoResourceAttr(resourceFullName, "saml_options"),
					resource.TestCheckNoResourceAttr(resourceFullName, "wsfed_options"),
					resource.TestCheckResourceAttr(resourceFullName, "external_link_options.home_page_url", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "hidden_from_app_portal", "true"),
				),
			},
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccApplication_ExternalLinkMinimal(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Application_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationConfig_ExternalLinkMinimal(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckNoResourceAttr(resourceFullName, "description"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckNoResourceAttr(resourceFullName, "icon"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options"),
					resource.TestCheckNoResourceAttr(resourceFullName, "saml_options"),
					resource.TestCheckNoResourceAttr(resourceFullName, "wsfed_options"),
					resource.TestCheckResourceAttr(resourceFullName, "external_link_options.home_page_url", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "hidden_from_app_portal", "false"),
				),
			},
		},
	})
}

// WS-Fed
func TestAccApplication_WSFedFull(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application.%s", resourceName)

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
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Application_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationConfig_WSFed_FullNewEnv(environmentName, licenseID, resourceName, name, image),
				Check:  testAccApplicationConfig_WSFed_FullCheck(resourceFullName, name),
			},
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccApplication_WSFedMinimal(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Application_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationConfig_WSFed_MinimalNewEnv(environmentName, licenseID, resourceName, name),
				Check:  testAccApplicationConfig_WSFed_MinimalCheck(resourceFullName, name),
			},
		},
	})
}

func TestAccApplication_WSFedMinimalMaximal(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application.%s", resourceName)

	name := resourceName

	data, _ := os.ReadFile("../../acctest/test_assets/image/image-logo.gif")
	image := base64.StdEncoding.EncodeToString(data)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Application_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationConfig_WSFed_Minimal(resourceName, name),
				Check:  testAccApplicationConfig_WSFed_MinimalCheck(resourceFullName, name),
			},
			{
				// Delete the minimal model
				Config:  testAccApplicationConfig_WSFed_Minimal(resourceName, name),
				Destroy: true,
			},
			{
				// Re-create with a complete model
				Config: testAccApplicationConfig_WSFed_Full(resourceName, name, image),
				Check:  testAccApplicationConfig_WSFed_FullCheck(resourceFullName, name),
			},
			{
				// Back to minimal model
				Config: testAccApplicationConfig_WSFed_Minimal(resourceName, name),
				Check:  testAccApplicationConfig_WSFed_MinimalCheck(resourceFullName, name),
			},
			{
				// Back to complete model
				Config: testAccApplicationConfig_WSFed_Full(resourceName, name, image),
				Check:  testAccApplicationConfig_WSFed_FullCheck(resourceFullName, name),
			},
			{
				// Test importing the resource
				Config:       testAccApplicationConfig_WSFed_Full(resourceName, name, image),
				ResourceName: fmt.Sprintf("pingone_application.%s", resourceName),
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["id"]), nil
					}
				}(),
				ImportStateVerifyIdentifierAttribute: "id",
				ImportState:                          true,
				ImportStateVerify:                    true,
			},
		},
	})
}

func testAccApplicationConfig_WSFed_MinimalCheck(resourceFullName, name string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "name", name),
		resource.TestCheckNoResourceAttr(resourceFullName, "description"),
		resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
		resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
		resource.TestCheckNoResourceAttr(resourceFullName, "login_page_url"),
		resource.TestCheckNoResourceAttr(resourceFullName, "icon"),
		resource.TestCheckNoResourceAttr(resourceFullName, "external_link_options"),
		resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options"),
		resource.TestCheckNoResourceAttr(resourceFullName, "saml_options"),
		resource.TestCheckResourceAttr(resourceFullName, "hidden_from_app_portal", "false"),
		resource.TestCheckResourceAttr(resourceFullName, "wsfed_options.audience_restriction", "urn:federation:MicrosoftOnline"),
		resource.TestCheckResourceAttr(resourceFullName, "wsfed_options.domain_name", "my.domain.name.example.com"),
		resource.TestCheckResourceAttr(resourceFullName, "wsfed_options.idp_signing_key.algorithm", "SHA384withECDSA"),
		resource.TestMatchResourceAttr(resourceFullName, "wsfed_options.idp_signing_key.key_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "wsfed_options.reply_url", "https://example.com"),
		resource.TestCheckResourceAttr(resourceFullName, "wsfed_options.type", "WEB_APP"),
	)
}

func testAccApplicationConfig_WSFed_FullCheck(resourceFullName, name string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "name", name),
		resource.TestCheckResourceAttr(resourceFullName, "description", "My test WS-Fed app"),
		resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
		resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
		resource.TestCheckResourceAttr(resourceFullName, "login_page_url", "https://www.pingidentity.com"),
		resource.TestMatchResourceAttr(resourceFullName, "icon.id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "icon.href", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
		resource.TestCheckResourceAttr(resourceFullName, "access_control_role_type", "ADMIN_USERS_ONLY"),
		resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.groups.#", "2"),
		resource.TestMatchResourceAttr(resourceFullName, "access_control_group_options.groups.0", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "access_control_group_options.groups.1", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.type", "ANY_GROUP"),
		resource.TestCheckResourceAttr(resourceFullName, "hidden_from_app_portal", "true"),
		resource.TestCheckNoResourceAttr(resourceFullName, "external_link_options"),
		resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options"),
		resource.TestCheckNoResourceAttr(resourceFullName, "saml_options"),
		resource.TestCheckResourceAttr(resourceFullName, "wsfed_options.audience_restriction", "urn:federation:Example"),
		resource.TestCheckResourceAttr(resourceFullName, "wsfed_options.cors_settings.behavior", "ALLOW_SPECIFIC_ORIGINS"),
		resource.TestCheckResourceAttr(resourceFullName, "wsfed_options.cors_settings.origins.#", "8"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "wsfed_options.cors_settings.origins.*", "http://localhost"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "wsfed_options.cors_settings.origins.*", "https://localhost"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "wsfed_options.cors_settings.origins.*", "http://auth.pingidentity.com"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "wsfed_options.cors_settings.origins.*", "https://auth.pingidentity.com"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "wsfed_options.cors_settings.origins.*", "http://*.pingidentity.com"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "wsfed_options.cors_settings.origins.*", "https://*.pingidentity.com"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "wsfed_options.cors_settings.origins.*", "http://192.168.1.1"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "wsfed_options.cors_settings.origins.*", "https://192.168.1.1"),
		resource.TestCheckResourceAttr(resourceFullName, "wsfed_options.domain_name", "my.updated.domain.name.example.com"),
		resource.TestCheckResourceAttr(resourceFullName, "wsfed_options.idp_signing_key.algorithm", "SHA384withECDSA"),
		resource.TestMatchResourceAttr(resourceFullName, "wsfed_options.idp_signing_key.key_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "wsfed_options.kerberos.gateways.#", "1"),
		resource.TestMatchResourceAttr(resourceFullName, "wsfed_options.kerberos.gateways.0.id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "wsfed_options.kerberos.gateways.0.type", "LDAP"),
		resource.TestMatchResourceAttr(resourceFullName, "wsfed_options.kerberos.gateways.0.user_type.id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "wsfed_options.reply_url", "https://example.com"),
		resource.TestCheckResourceAttr(resourceFullName, "wsfed_options.slo_endpoint", "https://example.com/slo"),
		resource.TestCheckResourceAttr(resourceFullName, "wsfed_options.subject_name_identifier_format", "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress"),
		resource.TestCheckResourceAttr(resourceFullName, "wsfed_options.type", "WEB_APP"),
	)
}

func TestAccApplication_Enabled(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Application_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationConfig_Enabled(resourceName, name, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "false"),
				),
			},
			{
				Config: testAccApplicationConfig_Enabled(resourceName, name, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
				),
			},
			{
				Config: testAccApplicationConfig_Enabled(resourceName, name, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "false"),
				),
			},
		},
	})
}

func TestAccApplication_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Application_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccApplicationConfig_OIDC_MinimalWeb(resourceName, name),
			},
			// Errors
			{
				ResourceName: resourceFullName,
				ImportState:  true,
				ExpectError:  regexp.MustCompile(`Unexpected Import Identifier`),
			},
			{
				ResourceName:  resourceFullName,
				ImportStateId: "/",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`Unexpected Import Identifier`),
			},
			{
				ResourceName:  resourceFullName,
				ImportStateId: "badformat/badformat",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`Unexpected Import Identifier`),
			},
		},
	})
}

func testAccApplicationConfig_NewEnv(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s"
  enabled        = true

  oidc_options = {
    type                       = "WEB_APP"
    grant_types                = ["AUTHORIZATION_CODE", "REFRESH_TOKEN"]
    response_types             = ["CODE"]
    token_endpoint_auth_method = "CLIENT_SECRET_BASIC"
    redirect_uris              = ["https://www.pingidentity.com"]
  }
}
`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccApplicationConfig_OIDC_FullWeb(resourceName, name, image string) string {
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

    home_page_url      = "https://www.pingidentity.com"
    initiate_login_uri = "https://www.pingidentity.com/initiate"
    target_link_uri    = "https://www.pingidentity.com/target"
    pkce_enforcement   = "OPTIONAL"

    idp_signoff = true

    par_requirement = "OPTIONAL"
    par_timeout     = 60

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
`, acctest.GenericSandboxEnvironment(), resourceName, name, image)
}

func testAccApplicationConfig_OIDC_MinimalWeb(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true

  oidc_options = {
    type                       = "WEB_APP"
    grant_types                = ["REFRESH_TOKEN", "AUTHORIZATION_CODE"]
    response_types             = ["CODE"]
    token_endpoint_auth_method = "CLIENT_SECRET_BASIC"
    redirect_uris              = ["https://www.pingidentity.com"]
  }
}
`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
func testAccApplicationConfig_OIDC_MinimalWeb_PrivateKeyJWT_JWKS(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true

  oidc_options = {
    type                       = "WEB_APP"
    grant_types                = ["REFRESH_TOKEN", "AUTHORIZATION_CODE"]
    response_types             = ["CODE"]
    token_endpoint_auth_method = "PRIVATE_KEY_JWT"
    redirect_uris              = ["https://www.pingidentity.com"]

    jwks = <<EOF
{
	"keys": [
	  {
		"kty": "RSA",
		"e": "AQAB",
		"use": "sig",
		"kid": "12345",
		"alg": "RS256",
		"n": "0vx7agoebGcQSuuPiLJXZptN9nndrQmbXEPCR0VH7jhV1JvKFvVsenY4rz5BnCNRS7U2mFF9K2BWXTZiaF4f3hjd4J0AOnHZV9KbV7L5Cp-1PEXF12R3HCkcqsn9b2I0xEs5OYLQaGbhV8v1FwTD2jZX0tAiw4i+SN9GJkPV2ZCOnF-8RPZCVDG9LZGFq4c9-YNPvRwT7B9-EN0kDYEKsOmGiJ0PVPAVTPZ9EVZVdLg7SwamytKcP4fz_BLTjCojz2W9KIL5UZGenQR5S7KAZxJ0T0DO8Q4kqNVdF7OOrBizX6-qQ9ZC1l6HJ6Sq9ye0oW2jTiMlNQxgc5vNRgHFAmb4DNA"
	  }
	]
}
EOF
  }
}
`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
func testAccApplicationConfig_OIDC_MinimalWeb_PrivateKeyJWT_JWKSURL(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true

  oidc_options = {
    type                       = "WEB_APP"
    grant_types                = ["REFRESH_TOKEN", "AUTHORIZATION_CODE"]
    response_types             = ["CODE"]
    token_endpoint_auth_method = "PRIVATE_KEY_JWT"
    redirect_uris              = ["https://www.pingidentity.com"]

    jwks_url = "https://pingidentity.com/jwks"
  }
}
`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
func testAccApplicationConfig_OIDC_MinimalWeb_ClientSecretJWT(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true

  oidc_options = {
    type                       = "WEB_APP"
    grant_types                = ["REFRESH_TOKEN", "AUTHORIZATION_CODE"]
    response_types             = ["CODE"]
    token_endpoint_auth_method = "CLIENT_SECRET_JWT"
    redirect_uris              = ["https://www.pingidentity.com"]
  }
}
`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccApplicationConfig_OIDC_FullNative(resourceName, name, image string) string {
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

    cors_settings = {
      behavior = "ALLOW_NO_ORIGINS"
    }

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
`, acctest.GenericSandboxEnvironment(), resourceName, name, image)
}

func testAccApplicationConfig_OIDC_MinimalNative(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true

  oidc_options = {
    type                       = "NATIVE_APP"
    grant_types                = ["CLIENT_CREDENTIALS"]
    token_endpoint_auth_method = "CLIENT_SECRET_BASIC"
  }
}
`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccApplicationConfig_OIDC_NativeKerberos(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_key" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id

  name                = "%[3]s"
  algorithm           = "EC"
  key_length          = 256
  signature_algorithm = "SHA384withECDSA"
  subject_dn          = "CN=%[3]s, OU=Ping Identity, O=Ping Identity, L=, ST=, C=US"
  usage_type          = "ISSUANCE"
  validity_period     = 365
}

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  name           = "%[3]s"
  enabled        = true

  oidc_options = {
    type                       = "NATIVE_APP"
    grant_types                = ["CLIENT_CREDENTIALS"]
    token_endpoint_auth_method = "CLIENT_SECRET_BASIC"

    certificate_based_authentication = {
      key_id = pingone_key.%[2]s.id
    }
  }
}
`, acctest.WorkforceSandboxEnvironment(), resourceName, name)
}

func testAccApplicationConfig_OIDC_NativeKerberosIncorrectKeyType(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_key" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id

  name                = "%[3]s"
  algorithm           = "EC"
  key_length          = 256
  signature_algorithm = "SHA384withECDSA"
  subject_dn          = "CN=%[3]s, OU=Ping Identity, O=Ping Identity, L=, ST=, C=US"
  usage_type          = "SIGNING"
  validity_period     = 365
}

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  name           = "%[3]s"
  enabled        = true

  oidc_options = {
    type                       = "NATIVE_APP"
    grant_types                = ["CLIENT_CREDENTIALS"]
    token_endpoint_auth_method = "CLIENT_SECRET_BASIC"

    certificate_based_authentication = {
      key_id = pingone_key.%[2]s.id
    }
  }
}
`, acctest.WorkforceSandboxEnvironment(), resourceName, name)
}

func testAccApplicationConfig_OIDC_NativeKerberosIncorrectApplicationType(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_key" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id

  name                = "%[3]s"
  algorithm           = "EC"
  key_length          = 256
  signature_algorithm = "SHA384withECDSA"
  subject_dn          = "CN=%[3]s, OU=Ping Identity, O=Ping Identity, L=, ST=, C=US"
  usage_type          = "ISSUANCE"
  validity_period     = 365
}

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  name           = "%[3]s"
  enabled        = true

  oidc_options = {
    type                       = "WORKER"
    grant_types                = ["CLIENT_CREDENTIALS"]
    token_endpoint_auth_method = "CLIENT_SECRET_BASIC"

    certificate_based_authentication = {
      key_id = pingone_key.%[2]s.id
    }
  }
}
`, acctest.WorkforceSandboxEnvironment(), resourceName, name)
}

func testAccApplicationConfig_OIDC_NativeMobile_Full(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true

  oidc_options = {
    type                       = "NATIVE_APP"
    grant_types                = ["CLIENT_CREDENTIALS"]
    token_endpoint_auth_method = "CLIENT_SECRET_BASIC"

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
`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccApplicationConfig_OIDC_NativeMobile_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true

  oidc_options = {
    type                       = "NATIVE_APP"
    grant_types                = ["CLIENT_CREDENTIALS"]
    token_endpoint_auth_method = "CLIENT_SECRET_BASIC"

    mobile_app = {}
  }
}
`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccApplicationConfig_OIDC_NativeMobile_IntegrityDetection_Full(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true

  oidc_options = {
    type                       = "NATIVE_APP"
    grant_types                = ["CLIENT_CREDENTIALS"]
    token_endpoint_auth_method = "CLIENT_SECRET_BASIC"

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
          amount = 45
          units  = "MINUTES"
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
`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccApplicationConfig_OIDC_NativeMobile_IntegrityDetection_Minimal(resourceName, name, googleJsonKey string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true

  oidc_options = {
    type                       = "NATIVE_APP"
    grant_types                = ["CLIENT_CREDENTIALS"]
    token_endpoint_auth_method = "CLIENT_SECRET_BASIC"

    mobile_app = {
      bundle_id           = "com.%[2]s.bundle"
      package_name        = "com.%[2]s.package"
      huawei_app_id       = "%[2]s"
      huawei_package_name = "com.%[2]s.huaweipackage"

      passcode_refresh_seconds = 45

      universal_app_link = "https://applink.com"

      integrity_detection = {
        enabled = true

        cache_duration = {
          amount = 30
        }

        google_play = {
          verification_type                = "GOOGLE"
          service_account_credentials_json = jsonencode(%[4]s)
        }
      }
    }
  }
}
`, acctest.GenericSandboxEnvironment(), resourceName, name, googleJsonKey)
}

func testAccApplicationConfig_OIDC_NativeMobile_IntegrityDetection_ExcludeGoogle(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true

  oidc_options = {
    type                       = "NATIVE_APP"
    grant_types                = ["CLIENT_CREDENTIALS"]
    token_endpoint_auth_method = "CLIENT_SECRET_BASIC"

    mobile_app = {
      bundle_id           = "com.%[2]s.bundle"
      package_name        = "com.%[2]s.package"
      huawei_app_id       = "%[2]s"
      huawei_package_name = "com.%[2]s.huaweipackage"

      passcode_refresh_seconds = 45

      universal_app_link = "https://applink.com"

      integrity_detection = {
        enabled = true

        cache_duration = {
          amount = 30
        }

        excluded_platforms = ["GOOGLE"]
      }
    }
  }
}
`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccApplicationConfig_OIDC_NativeMobile_IntegrityDetection_ExcludeIOS(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true

  oidc_options = {
    type                       = "NATIVE_APP"
    grant_types                = ["CLIENT_CREDENTIALS"]
    token_endpoint_auth_method = "CLIENT_SECRET_BASIC"

    mobile_app = {
      bundle_id           = "com.%[2]s.bundle"
      package_name        = "com.%[2]s.package"
      huawei_app_id       = "%[2]s"
      huawei_package_name = "com.%[2]s.huaweipackage"

      passcode_refresh_seconds = 45

      universal_app_link = "https://applink.com"

      integrity_detection = {
        enabled = true

        cache_duration = {
          amount = 30
        }

        excluded_platforms = ["IOS"]

        google_play = {
          verification_type = "INTERNAL"
          decryption_key    = "decryptionkeydoesnotexist"
          verification_key  = "verificationkeydoesnotexist"
        }
      }
    }
  }
}
`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccApplicationConfig_OIDC_FullCustom(resourceName, name, image string) string {
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
    type = "CUSTOM_APP"
    grant_types = [
      "AUTHORIZATION_CODE",
      "CLIENT_CREDENTIALS",
      "IMPLICIT",
      "REFRESH_TOKEN"
    ]
    response_types = [
      "CODE",
      "TOKEN",
      "ID_TOKEN"
    ]
    token_endpoint_auth_method      = "CLIENT_SECRET_BASIC"
    redirect_uris                   = ["https://www.pingidentity.com", "https://pingidentity.com"]
    allow_wildcard_in_redirect_uris = true
    post_logout_redirect_uris       = ["https://www.pingidentity.com/logout", "https://pingidentity.com/logout"]

    refresh_token_duration                             = 3000000
    refresh_token_rolling_duration                     = 30000000
    refresh_token_rolling_grace_period_duration        = 80000
    additional_refresh_token_replay_protection_enabled = false

    home_page_url      = "https://www.pingidentity.com"
    initiate_login_uri = "https://www.pingidentity.com/initiate"
    target_link_uri    = "https://www.pingidentity.com/target"
    pkce_enforcement   = "REQUIRED"

    idp_signoff = true

    par_requirement = "REQUIRED"
    par_timeout     = 180

    require_signed_request_object = true

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
`, acctest.GenericSandboxEnvironment(), resourceName, name, image)
}

func testAccApplicationConfig_OIDC_MinimalCustom(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true

  oidc_options = {
    type                       = "CUSTOM_APP"
    grant_types                = ["CLIENT_CREDENTIALS"]
    token_endpoint_auth_method = "CLIENT_SECRET_BASIC"
  }
}
`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccApplicationConfig_OIDC_Custom_Device_Full(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true

  oidc_options = {
    type                       = "CUSTOM_APP"
    grant_types                = ["DEVICE_CODE", "REFRESH_TOKEN"]
    token_endpoint_auth_method = "NONE"

    device_path_id                 = "mobileAppId-1"
    device_custom_verification_uri = "https://pingidentity.com/verification1"
    device_timeout                 = 500
    device_polling_interval        = 10
  }
}
`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccApplicationConfig_OIDC_Custom_Device_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true

  oidc_options = {
    type                       = "CUSTOM_APP"
    grant_types                = ["DEVICE_CODE"]
    token_endpoint_auth_method = "CLIENT_SECRET_BASIC"
  }
}
`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccApplicationConfig_OIDC_FullService(resourceName, name, image string) string {
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
    type = "SERVICE"
    grant_types = [
      "AUTHORIZATION_CODE",
      "CLIENT_CREDENTIALS",
      "IMPLICIT",
      "REFRESH_TOKEN"
    ]
    response_types = [
      "CODE",
      "TOKEN",
      "ID_TOKEN"
    ]
    token_endpoint_auth_method      = "CLIENT_SECRET_BASIC"
    redirect_uris                   = ["https://www.pingidentity.com", "https://pingidentity.com"]
    allow_wildcard_in_redirect_uris = true
    post_logout_redirect_uris       = ["https://www.pingidentity.com/logout", "https://pingidentity.com/logout"]

    refresh_token_duration                             = 3000000
    refresh_token_rolling_duration                     = 30000000
    refresh_token_rolling_grace_period_duration        = 80000
    additional_refresh_token_replay_protection_enabled = false

    home_page_url      = "https://www.pingidentity.com"
    initiate_login_uri = "https://www.pingidentity.com/initiate"
    target_link_uri    = "https://www.pingidentity.com/target"
    pkce_enforcement   = "REQUIRED"

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
`, acctest.GenericSandboxEnvironment(), resourceName, name, image)
}

func testAccApplicationConfig_OIDC_MinimalService(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true

  oidc_options = {
    type                       = "SERVICE"
    grant_types                = ["REFRESH_TOKEN", "AUTHORIZATION_CODE"]
    response_types             = ["CODE"]
    token_endpoint_auth_method = "CLIENT_SECRET_BASIC"
    redirect_uris              = ["https://www.pingidentity.com"]
  }
}
`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccApplicationConfig_OIDC_FullSPA(resourceName, name, image string) string {
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
    type                            = "SINGLE_PAGE_APP"
    grant_types                     = ["AUTHORIZATION_CODE"]
    response_types                  = ["CODE"]
    pkce_enforcement                = "S256_REQUIRED"
    token_endpoint_auth_method      = "NONE"
    redirect_uris                   = ["https://www.pingidentity.com", "https://pingidentity.com"]
    allow_wildcard_in_redirect_uris = true
    post_logout_redirect_uris       = ["https://www.pingidentity.com/logout", "https://pingidentity.com/logout"]
    home_page_url                   = "https://www.pingidentity.com"
    initiate_login_uri              = "https://www.pingidentity.com/initiate"
    target_link_uri                 = "https://www.pingidentity.com/target"

    support_unsigned_request_object = true

    idp_signoff = true

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
`, acctest.GenericSandboxEnvironment(), resourceName, name, image)
}

func testAccApplicationConfig_OIDC_MinimalSPA(resourceName, name string) string {
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
    pkce_enforcement           = "S256_REQUIRED"
    token_endpoint_auth_method = "NONE"
    redirect_uris              = ["https://www.pingidentity.com"]
  }
}
`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccApplicationConfig_OIDC_FullWorker(resourceName, name, image string) string {
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
      pingone_group.%[2]s-2.id,
      pingone_group.%[2]s-1.id
    ]
  }

  hidden_from_app_portal = true

  enabled = true

  oidc_options = {
    type                       = "WORKER"
    grant_types                = ["CLIENT_CREDENTIALS"]
    token_endpoint_auth_method = "CLIENT_SECRET_BASIC"

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
`, acctest.GenericSandboxEnvironment(), resourceName, name, image)
}

func testAccApplicationConfig_OIDC_MinimalWorker(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true

  oidc_options = {
    type                       = "WORKER"
    grant_types                = ["CLIENT_CREDENTIALS"]
    token_endpoint_auth_method = "CLIENT_SECRET_BASIC"
  }
}
`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccApplicationConfig_OIDC_WildcardInRedirect(resourceName, name string, wildcardInRedirect bool) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true

  oidc_options = {
    type                            = "SINGLE_PAGE_APP"
    grant_types                     = ["AUTHORIZATION_CODE"]
    response_types                  = ["CODE"]
    pkce_enforcement                = "S256_REQUIRED"
    token_endpoint_auth_method      = "NONE"
    redirect_uris                   = ["https://www.pingidentity.com/*"]
    allow_wildcard_in_redirect_uris = %[4]t
  }
}
`, acctest.GenericSandboxEnvironment(), resourceName, name, wildcardInRedirect)
}

func testAccApplicationConfig_OIDC_LocalhostAddresses(resourceName, name, hostname string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true

  login_page_url = "http://%[4]s/login" # https with the exception of localhost

  oidc_options = {
    home_page_url              = "http://%[4]s/home" # https with the exception of localhost
    initiate_login_uri         = "http://%[4]s/init" # https with the exception of localhost
    type                       = "SINGLE_PAGE_APP"
    grant_types                = ["AUTHORIZATION_CODE"]
    response_types             = ["CODE"]
    pkce_enforcement           = "S256_REQUIRED"
    token_endpoint_auth_method = "NONE"
    redirect_uris              = ["http://%[4]s/callback"] # https with the exception of localhost
    post_logout_redirect_uris  = ["http://%[4]s/logout"]   # either http or https
    target_link_uri            = "http://%[4]s/link"       # either http or https
  }
}
`, acctest.GenericSandboxEnvironment(), resourceName, name, hostname)
}

func testAccApplicationConfig_OIDC_NativeAppAddresses(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true

  oidc_options = {
    type                       = "NATIVE_APP"
    grant_types                = ["AUTHORIZATION_CODE"]
    response_types             = ["CODE"]
    pkce_enforcement           = "S256_REQUIRED"
    token_endpoint_auth_method = "NONE"
    redirect_uris              = ["com.myapp.app://callback"]
    post_logout_redirect_uris  = ["com.myapp.app://logout"]
    initiate_login_uri         = "https://pingidentity.com/target"
    target_link_uri            = "com.myapp.app://target"
  }
}
`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccApplicationConfig_SAML_Full(environmentName, licenseID, resourceName, name, image, pkcs7_cert, pem_cert string) string {
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
    idp_signing_key = {
      key_id    = pingone_key.%[3]s.id
      algorithm = pingone_key.%[3]s.signature_algorithm
    }
    enable_requested_authn_context   = true
    nameid_format                    = "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress"
    response_is_signed               = true
    session_not_on_or_after_duration = 64
    slo_binding                      = "HTTP_REDIRECT"
    slo_endpoint                     = "https://www.pingidentity.com/sloendpoint"
    slo_response_endpoint            = "https://www.pingidentity.com/sloresponseendpoint"
    slo_window                       = 3

    default_target_url = "https://www.pingidentity.com/relaystate"

    sp_verification = {
      authn_request_signed = true
      certificate_ids = [
        pingone_certificate.%[3]s-2.id,
        pingone_certificate.%[3]s-1.id,
      ]
    }

    virtual_server_id_settings = {
      enabled = true
      virtual_server_ids = [
        {
          vs_id   = "virtualserver1"
          default = true
        },
        {
          vs_id = "virtualserver2"
        },
        {
          vs_id = "virtualserver3"
        },
        {
          vs_id = "virtualserver4"
        },
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
  }
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name, image, pkcs7_cert, pem_cert)
}

func testAccApplicationConfig_SAML_Minimal(environmentName, licenseID, resourceName, name, image, pkcs7_cert, pem_cert string) string {
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
  enabled        = true

  saml_options = {
    acs_urls           = ["https://pingidentity.com"]
    assertion_duration = 3600
    sp_entity_id       = "sp:entity:%[3]s"

    idp_signing_key = {
      key_id    = pingone_key.%[3]s.id
      algorithm = pingone_key.%[3]s.signature_algorithm
    }
  }
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name, image, pkcs7_cert, pem_cert)
}

func testAccApplicationConfig_SAML_VirtualServerIdSettingsOrdering(environmentName, licenseID, resourceName, name, image, pkcs7_cert, pem_cert string) string {
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
  description    = "Test SAML app with virtual server ID settings ordering"
  enabled        = true

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
    idp_signing_key = {
      key_id    = pingone_key.%[3]s.id
      algorithm = pingone_key.%[3]s.signature_algorithm
    }
    enable_requested_authn_context   = true
    nameid_format                    = "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress"
    response_is_signed               = true
    session_not_on_or_after_duration = 64
    slo_binding                      = "HTTP_REDIRECT"
    slo_endpoint                     = "https://www.pingidentity.com/sloendpoint"
    slo_response_endpoint            = "https://www.pingidentity.com/sloresponseendpoint"
    slo_window                       = 3

    default_target_url = "https://www.pingidentity.com/relaystate"

    sp_verification = {
      authn_request_signed = true
      certificate_ids = [
        pingone_certificate.%[3]s-2.id,
        pingone_certificate.%[3]s-1.id,
      ]
    }

    virtual_server_id_settings = {
      enabled = false
      virtual_server_ids = [
        {
          vs_id   = "virtualserver2"
          default = false
        },
        {
          vs_id   = "virtualserver3"
          default = false
        },
        {
          vs_id   = "virtualserver4"
          default = false
        },
        {
          vs_id   = "virtualserver1"
          default = true
        }
      ]
    }
  }
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name, image, pkcs7_cert, pem_cert)
}

func testAccApplicationConfig_SAML_VirtualServerIdSettingsReordered(environmentName, licenseID, resourceName, name, image, pkcs7_cert, pem_cert string) string {
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
  description    = "Test SAML app with reordered virtual server ID settings"
  enabled        = true

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
    idp_signing_key = {
      key_id    = pingone_key.%[3]s.id
      algorithm = pingone_key.%[3]s.signature_algorithm
    }
    enable_requested_authn_context   = true
    nameid_format                    = "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress"
    response_is_signed               = true
    session_not_on_or_after_duration = 64
    slo_binding                      = "HTTP_REDIRECT"
    slo_endpoint                     = "https://www.pingidentity.com/sloendpoint"
    slo_response_endpoint            = "https://www.pingidentity.com/sloresponseendpoint"
    slo_window                       = 3

    default_target_url = "https://www.pingidentity.com/relaystate"

    sp_verification = {
      authn_request_signed = true
      certificate_ids = [
        pingone_certificate.%[3]s-2.id,
        pingone_certificate.%[3]s-1.id,
      ]
    }

    virtual_server_id_settings = {
      enabled = true
      virtual_server_ids = [
        {
          vs_id   = "virtualserver1"
          default = false
        },
        {
          vs_id   = "virtualserver2"
          default = false
        },
        {
          vs_id   = "virtualserver4"
          default = false
        },
        {
          vs_id   = "virtualserver3"
          default = true
        }
      ]
    }
  }
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name, image, pkcs7_cert, pem_cert)
}

func testAccApplicationConfig_ExternalLinkFull(resourceName, name, image string) string {
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
}`, acctest.GenericSandboxEnvironment(), resourceName, name, image)
}

func testAccApplicationConfig_ExternalLinkMinimal(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true

  external_link_options = {
    home_page_url = "https://www.pingidentity.com"
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccApplicationConfig_Enabled(resourceName, name string, enabled bool) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = %[4]t

  oidc_options = {
    type                       = "SINGLE_PAGE_APP"
    grant_types                = ["AUTHORIZATION_CODE"]
    response_types             = ["CODE"]
    pkce_enforcement           = "S256_REQUIRED"
    token_endpoint_auth_method = "NONE"
    redirect_uris              = ["https://www.pingidentity.com"]
  }
}
`, acctest.GenericSandboxEnvironment(), resourceName, name, enabled)
}

func testAccApplicationConfig_WSFed_Full(resourceName, name, image string) string {
	return testAccApplicationConfig_WSFed_FullWithEnv(acctest.GenericSandboxEnvironment(), "data.pingone_environment.general_test.id", resourceName, name, image)
}

func testAccApplicationConfig_WSFed_FullNewEnv(environmentName, licenseID, resourceName, name, image string) string {
	return testAccApplicationConfig_WSFed_FullWithEnv(acctest.MinimalSandboxEnvironment(environmentName, licenseID), fmt.Sprintf("pingone_environment.%s.id", environmentName), resourceName, name, image)
}

func testAccApplicationConfig_WSFed_FullWithEnv(environmentHcl, environmentIdResourceLink, resourceName, name, image string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_group" "%[3]s-1" {
  environment_id = %[2]s
  name           = "%[4]s-1"
}

resource "pingone_group" "%[3]s-2" {
  environment_id = %[2]s
  name           = "%[4]s-2"
}

resource "pingone_key" "%[3]s" {
  environment_id = %[2]s

  name                = "%[4]s"
  algorithm           = "EC"
  key_length          = 256
  signature_algorithm = "SHA384withECDSA"
  subject_dn          = "CN=%[4]s, OU=Ping Identity, O=Ping Identity, L=, ST=, C=US"
  usage_type          = "SIGNING"
  validity_period     = 365
}

resource "pingone_image" "%[3]s" {
  environment_id = %[2]s

  image_file_base64 = "%[5]s"
}

resource "pingone_population" "%[3]s" {
  environment_id = %[2]s

  name = "%[4]s"
}

resource "pingone_gateway" "%[3]s" {
  environment_id = %[2]s
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
  environment_id = %[2]s
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
}`, environmentHcl, environmentIdResourceLink, resourceName, name, image)
}

func testAccApplicationConfig_WSFed_Minimal(resourceName, name string) string {
	return testAccApplicationConfig_WSFed_MinimalWithEnv(acctest.GenericSandboxEnvironment(), "data.pingone_environment.general_test.id", resourceName, name)
}

func testAccApplicationConfig_WSFed_MinimalNewEnv(environmentName, licenseID, resourceName, name string) string {
	return testAccApplicationConfig_WSFed_MinimalWithEnv(acctest.MinimalSandboxEnvironment(environmentName, licenseID), fmt.Sprintf("pingone_environment.%s.id", environmentName), resourceName, name)
}

func testAccApplicationConfig_WSFed_MinimalWithEnv(environmentHcl, environmentIdResourceLink, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_key" "%[3]s" {
  environment_id = %[2]s

  name                = "%[4]s"
  algorithm           = "EC"
  key_length          = 256
  signature_algorithm = "SHA384withECDSA"
  subject_dn          = "CN=%[4]s, OU=Ping Identity, O=Ping Identity, L=, ST=, C=US"
  usage_type          = "SIGNING"
  validity_period     = 365
}

resource "pingone_application" "%[3]s" {
  environment_id = %[2]s
  name           = "%[4]s"
  enabled        = true

  wsfed_options = {
    domain_name = "my.domain.name.example.com"
    idp_signing_key = {
      key_id    = pingone_key.%[3]s.id
      algorithm = pingone_key.%[3]s.signature_algorithm
    }
    reply_url = "https://example.com"
    type      = "WEB_APP"
  }
}`, environmentHcl, environmentIdResourceLink, resourceName, name)
}
