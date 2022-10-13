package sso_test

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
)

func testAccCheckApplicationDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_application" {
			continue
		}

		_, rEnv, err := apiClient.EnvironmentsApi.ReadOneEnvironment(ctx, rs.Primary.Attributes["environment_id"]).Execute()

		if err != nil {

			if rEnv == nil {
				return fmt.Errorf("Response object does not exist and no error detected")
			}

			if rEnv.StatusCode == 404 {
				continue
			}

			return err
		}

		body, r, err := apiClient.ApplicationsApi.ReadOneApplication(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.ID).Execute()

		if err != nil {

			if r == nil {
				return fmt.Errorf("Response object does not exist and no error detected")
			}

			if r.StatusCode == 404 {
				continue
			}

			tflog.Error(ctx, fmt.Sprintf("Error: %v", body))
			return err
		}

		return fmt.Errorf("PingOne Application Instance %s still exists", rs.Primary.ID)
	}

	return nil
}

func TestAccApplication_NewEnv(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckApplicationDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
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

	data, _ := ioutil.ReadFile("../../acctest/test_assets/image/image-logo.gif")
	image := base64.StdEncoding.EncodeToString(data)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckApplicationDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationConfig_OIDC_FullWeb(resourceName, name, image),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", "My test OIDC app"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "login_page_url", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "icon.#", "1"),
					resource.TestMatchResourceAttr(resourceFullName, "icon.0.id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "icon.0.href", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_role_type", "ADMIN_USERS_ONLY"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "access_control_group_options.*", map[string]string{
						"type":     "ANY_GROUP",
						"groups.#": "1",
					}),
					resource.TestMatchTypeSetElemNestedAttrs(resourceFullName, "access_control_group_options.*", map[string]*regexp.Regexp{
						"groups.0": regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
					}),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.type", "WEB_APP"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.home_page_url", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.grant_types.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.grant_types.*", "AUTHORIZATION_CODE"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.grant_types.*", "REFRESH_TOKEN"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.response_types.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.response_types.*", "CODE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.token_endpoint_authn_method", "CLIENT_SECRET_BASIC"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.pkce_enforcement", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.redirect_uris.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.redirect_uris.*", "https://pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.redirect_uris.*", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.post_logout_redirect_uris.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.post_logout_redirect_uris.*", "https://pingidentity.com/logout"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.post_logout_redirect_uris.*", "https://www.pingidentity.com/logout"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.refresh_token_duration", "3000000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.refresh_token_rolling_duration", "30000000"),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.0.client_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.0.client_secret", regexp.MustCompile(`[a-zA-Z0-9-~_]{10,}`)),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.support_unsigned_request_object", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.bundle_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.package_name", ""),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "external_link_options.#", "0"),
				),
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
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckApplicationDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationConfig_OIDC_MinimalWeb(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", ""),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "login_page_url", ""),
					resource.TestCheckResourceAttr(resourceFullName, "icon.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.type", "WEB_APP"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.home_page_url", ""),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.grant_types.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.grant_types.*", "AUTHORIZATION_CODE"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.grant_types.*", "REFRESH_TOKEN"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.response_types.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.response_types.*", "CODE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.token_endpoint_authn_method", "CLIENT_SECRET_BASIC"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.pkce_enforcement", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.redirect_uris.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.redirect_uris.*", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.post_logout_redirect_uris.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.refresh_token_duration", "2592000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.refresh_token_rolling_duration", "15552000"),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.0.client_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.0.client_secret", regexp.MustCompile(`[a-zA-Z0-9-~_]{10,}`)),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.support_unsigned_request_object", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.bundle_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.package_name", ""),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "external_link_options.#", "0"),
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

	data, _ := ioutil.ReadFile("../../acctest/test_assets/image/image-logo.gif")
	image := base64.StdEncoding.EncodeToString(data)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckApplicationDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationConfig_OIDC_MinimalWeb(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", ""),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "login_page_url", ""),
					resource.TestCheckResourceAttr(resourceFullName, "icon.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.type", "WEB_APP"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.home_page_url", ""),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.grant_types.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.grant_types.*", "AUTHORIZATION_CODE"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.grant_types.*", "REFRESH_TOKEN"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.response_types.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.response_types.*", "CODE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.token_endpoint_authn_method", "CLIENT_SECRET_BASIC"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.pkce_enforcement", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.redirect_uris.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.redirect_uris.*", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.post_logout_redirect_uris.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.refresh_token_duration", "2592000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.refresh_token_rolling_duration", "15552000"),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.0.client_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.0.client_secret", regexp.MustCompile(`[a-zA-Z0-9-~_]{10,}`)),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.support_unsigned_request_object", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.bundle_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.package_name", ""),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "external_link_options.#", "0"),
				),
			},
			{
				Config: testAccApplicationConfig_OIDC_FullWeb(resourceName, name, image),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", "My test OIDC app"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "login_page_url", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "icon.#", "1"),
					resource.TestMatchResourceAttr(resourceFullName, "icon.0.id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "icon.0.href", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_role_type", "ADMIN_USERS_ONLY"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "access_control_group_options.*", map[string]string{
						"type":     "ANY_GROUP",
						"groups.#": "1",
					}),
					resource.TestMatchTypeSetElemNestedAttrs(resourceFullName, "access_control_group_options.*", map[string]*regexp.Regexp{
						"groups.0": regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
					}),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.type", "WEB_APP"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.home_page_url", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.grant_types.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.grant_types.*", "AUTHORIZATION_CODE"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.grant_types.*", "REFRESH_TOKEN"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.response_types.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.response_types.*", "CODE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.token_endpoint_authn_method", "CLIENT_SECRET_BASIC"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.pkce_enforcement", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.redirect_uris.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.redirect_uris.*", "https://pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.redirect_uris.*", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.post_logout_redirect_uris.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.post_logout_redirect_uris.*", "https://pingidentity.com/logout"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.post_logout_redirect_uris.*", "https://www.pingidentity.com/logout"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.refresh_token_duration", "3000000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.refresh_token_rolling_duration", "30000000"),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.0.client_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.0.client_secret", regexp.MustCompile(`[a-zA-Z0-9-~_]{10,}`)),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.support_unsigned_request_object", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.bundle_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.package_name", ""),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "external_link_options.#", "0"),
				),
			},
			{
				Config: testAccApplicationConfig_OIDC_MinimalWeb(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", ""),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "login_page_url", ""),
					resource.TestCheckResourceAttr(resourceFullName, "icon.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.type", "WEB_APP"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.home_page_url", ""),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.grant_types.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.grant_types.*", "AUTHORIZATION_CODE"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.grant_types.*", "REFRESH_TOKEN"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.response_types.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.response_types.*", "CODE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.token_endpoint_authn_method", "CLIENT_SECRET_BASIC"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.pkce_enforcement", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.redirect_uris.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.redirect_uris.*", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.post_logout_redirect_uris.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.refresh_token_duration", "2592000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.refresh_token_rolling_duration", "15552000"),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.0.client_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.0.client_secret", regexp.MustCompile(`[a-zA-Z0-9-~_]{10,}`)),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.support_unsigned_request_object", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.bundle_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.package_name", ""),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "external_link_options.#", "0"),
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

	data, _ := ioutil.ReadFile("../../acctest/test_assets/image/image-logo.gif")
	image := base64.StdEncoding.EncodeToString(data)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckApplicationDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationConfig_OIDC_FullNative(resourceName, name, image),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", "My test OIDC app"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "login_page_url", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "icon.#", "1"),
					resource.TestMatchResourceAttr(resourceFullName, "icon.0.id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "icon.0.href", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_role_type", "ADMIN_USERS_ONLY"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "access_control_group_options.*", map[string]string{
						"type":     "ANY_GROUP",
						"groups.#": "1",
					}),
					resource.TestMatchTypeSetElemNestedAttrs(resourceFullName, "access_control_group_options.*", map[string]*regexp.Regexp{
						"groups.0": regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
					}),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.type", "NATIVE_APP"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.home_page_url", ""),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.grant_types.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.grant_types.*", "CLIENT_CREDENTIALS"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.response_types.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.token_endpoint_authn_method", "CLIENT_SECRET_BASIC"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.pkce_enforcement", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.redirect_uris.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.post_logout_redirect_uris.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.refresh_token_duration", "2592000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.refresh_token_rolling_duration", "15552000"),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.0.client_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.0.client_secret", regexp.MustCompile(`[a-zA-Z0-9-~_]{10,}`)),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.support_unsigned_request_object", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.bundle_id", fmt.Sprintf("com.%s.bundle", resourceName)),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.package_name", fmt.Sprintf("com.%s.package", resourceName)),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.passcode_refresh_seconds", "45"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.universal_app_link", "https://applink.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.integrity_detection.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.integrity_detection.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.integrity_detection.0.excluded_platforms.0", "IOS"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.integrity_detection.0.cache_duration.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.integrity_detection.0.cache_duration.0.amount", "30"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.integrity_detection.0.cache_duration.0.units", "HOURS"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.bundle_id", fmt.Sprintf("com.%s.bundle", resourceName)),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.package_name", fmt.Sprintf("com.%s.package", resourceName)),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "external_link_options.#", "0"),
				),
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
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckApplicationDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationConfig_OIDC_MinimalNative(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", ""),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "login_page_url", ""),
					resource.TestCheckResourceAttr(resourceFullName, "icon.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.type", "NATIVE_APP"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.home_page_url", ""),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.grant_types.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.grant_types.*", "CLIENT_CREDENTIALS"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.response_types.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.token_endpoint_authn_method", "CLIENT_SECRET_BASIC"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.pkce_enforcement", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.redirect_uris.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.post_logout_redirect_uris.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.refresh_token_duration", "2592000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.refresh_token_rolling_duration", "15552000"),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.0.client_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.0.client_secret", regexp.MustCompile(`[a-zA-Z0-9-~_]{10,}`)),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.support_unsigned_request_object", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.bundle_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.package_name", ""),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.universal_app_link", ""),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.integrity_detection.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.integrity_detection.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.integrity_detection.0.excluded_platforms.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.passcode_refresh_seconds", "30"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.bundle_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.package_name", ""),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "external_link_options.#", "0"),
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

	data, _ := ioutil.ReadFile("../../acctest/test_assets/image/image-logo.gif")
	image := base64.StdEncoding.EncodeToString(data)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckApplicationDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationConfig_OIDC_FullNative(resourceName, name, image),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", "My test OIDC app"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "login_page_url", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "icon.#", "1"),
					resource.TestMatchResourceAttr(resourceFullName, "icon.0.id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "icon.0.href", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_role_type", "ADMIN_USERS_ONLY"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "access_control_group_options.*", map[string]string{
						"type":     "ANY_GROUP",
						"groups.#": "1",
					}),
					resource.TestMatchTypeSetElemNestedAttrs(resourceFullName, "access_control_group_options.*", map[string]*regexp.Regexp{
						"groups.0": regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
					}),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.type", "NATIVE_APP"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.home_page_url", ""),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.grant_types.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.grant_types.*", "CLIENT_CREDENTIALS"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.response_types.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.token_endpoint_authn_method", "CLIENT_SECRET_BASIC"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.pkce_enforcement", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.redirect_uris.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.post_logout_redirect_uris.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.refresh_token_duration", "2592000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.refresh_token_rolling_duration", "15552000"),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.0.client_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.0.client_secret", regexp.MustCompile(`[a-zA-Z0-9-~_]{10,}`)),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.support_unsigned_request_object", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.bundle_id", fmt.Sprintf("com.%s.bundle", resourceName)),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.package_name", fmt.Sprintf("com.%s.package", resourceName)),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.passcode_refresh_seconds", "45"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.universal_app_link", "https://applink.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.integrity_detection.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.integrity_detection.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.integrity_detection.0.excluded_platforms.0", "IOS"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.integrity_detection.0.cache_duration.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.integrity_detection.0.cache_duration.0.amount", "30"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.integrity_detection.0.cache_duration.0.units", "HOURS"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.bundle_id", fmt.Sprintf("com.%s.bundle", resourceName)),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.package_name", fmt.Sprintf("com.%s.package", resourceName)),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "external_link_options.#", "0"),
				),
			},
			{
				Config: testAccApplicationConfig_OIDC_MinimalNative(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", ""),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "login_page_url", ""),
					resource.TestCheckResourceAttr(resourceFullName, "icon.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.type", "NATIVE_APP"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.home_page_url", ""),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.grant_types.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.grant_types.*", "CLIENT_CREDENTIALS"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.response_types.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.token_endpoint_authn_method", "CLIENT_SECRET_BASIC"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.pkce_enforcement", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.redirect_uris.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.post_logout_redirect_uris.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.refresh_token_duration", "2592000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.refresh_token_rolling_duration", "15552000"),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.0.client_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.0.client_secret", regexp.MustCompile(`[a-zA-Z0-9-~_]{10,}`)),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.support_unsigned_request_object", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.bundle_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.package_name", ""),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.universal_app_link", ""),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.integrity_detection.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.integrity_detection.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.integrity_detection.0.excluded_platforms.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.passcode_refresh_seconds", "30"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.bundle_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.package_name", ""),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "external_link_options.#", "0"),
				),
			},
			{
				Config: testAccApplicationConfig_OIDC_FullNative(resourceName, name, image),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", "My test OIDC app"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "login_page_url", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "icon.#", "1"),
					resource.TestMatchResourceAttr(resourceFullName, "icon.0.id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "icon.0.href", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_role_type", "ADMIN_USERS_ONLY"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "access_control_group_options.*", map[string]string{
						"type":     "ANY_GROUP",
						"groups.#": "1",
					}),
					resource.TestMatchTypeSetElemNestedAttrs(resourceFullName, "access_control_group_options.*", map[string]*regexp.Regexp{
						"groups.0": regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
					}),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.type", "NATIVE_APP"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.home_page_url", ""),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.grant_types.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.grant_types.*", "CLIENT_CREDENTIALS"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.response_types.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.token_endpoint_authn_method", "CLIENT_SECRET_BASIC"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.pkce_enforcement", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.redirect_uris.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.post_logout_redirect_uris.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.refresh_token_duration", "2592000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.refresh_token_rolling_duration", "15552000"),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.0.client_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.0.client_secret", regexp.MustCompile(`[a-zA-Z0-9-~_]{10,}`)),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.support_unsigned_request_object", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.bundle_id", fmt.Sprintf("com.%s.bundle", resourceName)),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.package_name", fmt.Sprintf("com.%s.package", resourceName)),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.passcode_refresh_seconds", "45"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.universal_app_link", "https://applink.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.integrity_detection.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.integrity_detection.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.integrity_detection.0.excluded_platforms.0", "IOS"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.integrity_detection.0.cache_duration.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.integrity_detection.0.cache_duration.0.amount", "30"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.integrity_detection.0.cache_duration.0.units", "HOURS"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.bundle_id", fmt.Sprintf("com.%s.bundle", resourceName)),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.package_name", fmt.Sprintf("com.%s.package", resourceName)),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "external_link_options.#", "0"),
				),
			},
		},
	})
}

func TestAccApplication_OIDCFullCustom(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application.%s", resourceName)

	name := resourceName

	data, _ := ioutil.ReadFile("../../acctest/test_assets/image/image-logo.gif")
	image := base64.StdEncoding.EncodeToString(data)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckApplicationDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationConfig_OIDC_FullCustom(resourceName, name, image),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", "My test OIDC app"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "login_page_url", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "icon.#", "1"),
					resource.TestMatchResourceAttr(resourceFullName, "icon.0.id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "icon.0.href", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_role_type", "ADMIN_USERS_ONLY"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "access_control_group_options.*", map[string]string{
						"type":     "ANY_GROUP",
						"groups.#": "1",
					}),
					resource.TestMatchTypeSetElemNestedAttrs(resourceFullName, "access_control_group_options.*", map[string]*regexp.Regexp{
						"groups.0": regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
					}),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.type", "CUSTOM_APP"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.home_page_url", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.grant_types.#", "4"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.grant_types.*", "AUTHORIZATION_CODE"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.grant_types.*", "CLIENT_CREDENTIALS"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.grant_types.*", "IMPLICIT"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.grant_types.*", "REFRESH_TOKEN"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.response_types.#", "3"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.response_types.*", "CODE"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.response_types.*", "TOKEN"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.response_types.*", "ID_TOKEN"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.token_endpoint_authn_method", "CLIENT_SECRET_BASIC"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.pkce_enforcement", "REQUIRED"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.redirect_uris.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.redirect_uris.*", "https://pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.redirect_uris.*", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.post_logout_redirect_uris.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.post_logout_redirect_uris.*", "https://pingidentity.com/logout"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.post_logout_redirect_uris.*", "https://www.pingidentity.com/logout"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.refresh_token_duration", "3000000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.refresh_token_rolling_duration", "30000000"),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.0.client_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.0.client_secret", regexp.MustCompile(`[a-zA-Z0-9-~_]{10,}`)),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.support_unsigned_request_object", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.bundle_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.package_name", ""),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "external_link_options.#", "0"),
				),
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
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckApplicationDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationConfig_OIDC_MinimalCustom(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", ""),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "login_page_url", ""),
					resource.TestCheckResourceAttr(resourceFullName, "icon.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.type", "CUSTOM_APP"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.home_page_url", ""),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.grant_types.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.grant_types.*", "CLIENT_CREDENTIALS"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.response_types.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.token_endpoint_authn_method", "CLIENT_SECRET_BASIC"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.pkce_enforcement", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.redirect_uris.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.post_logout_redirect_uris.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.refresh_token_duration", "2592000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.refresh_token_rolling_duration", "15552000"),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.0.client_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.0.client_secret", regexp.MustCompile(`[a-zA-Z0-9-~_]{10,}`)),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.support_unsigned_request_object", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.bundle_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.package_name", ""),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "external_link_options.#", "0"),
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

	data, _ := ioutil.ReadFile("../../acctest/test_assets/image/image-logo.gif")
	image := base64.StdEncoding.EncodeToString(data)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckApplicationDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationConfig_OIDC_FullCustom(resourceName, name, image),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", "My test OIDC app"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "login_page_url", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "icon.#", "1"),
					resource.TestMatchResourceAttr(resourceFullName, "icon.0.id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "icon.0.href", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_role_type", "ADMIN_USERS_ONLY"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "access_control_group_options.*", map[string]string{
						"type":     "ANY_GROUP",
						"groups.#": "1",
					}),
					resource.TestMatchTypeSetElemNestedAttrs(resourceFullName, "access_control_group_options.*", map[string]*regexp.Regexp{
						"groups.0": regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
					}),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.type", "CUSTOM_APP"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.home_page_url", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.grant_types.#", "4"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.grant_types.*", "AUTHORIZATION_CODE"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.grant_types.*", "CLIENT_CREDENTIALS"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.grant_types.*", "IMPLICIT"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.grant_types.*", "REFRESH_TOKEN"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.response_types.#", "3"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.response_types.*", "CODE"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.response_types.*", "TOKEN"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.response_types.*", "ID_TOKEN"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.token_endpoint_authn_method", "CLIENT_SECRET_BASIC"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.pkce_enforcement", "REQUIRED"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.redirect_uris.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.redirect_uris.*", "https://pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.redirect_uris.*", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.post_logout_redirect_uris.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.post_logout_redirect_uris.*", "https://pingidentity.com/logout"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.post_logout_redirect_uris.*", "https://www.pingidentity.com/logout"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.refresh_token_duration", "3000000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.refresh_token_rolling_duration", "30000000"),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.0.client_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.0.client_secret", regexp.MustCompile(`[a-zA-Z0-9-~_]{10,}`)),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.support_unsigned_request_object", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.bundle_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.package_name", ""),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "external_link_options.#", "0"),
				),
			},
			{
				Config: testAccApplicationConfig_OIDC_MinimalCustom(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", ""),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "login_page_url", ""),
					resource.TestCheckResourceAttr(resourceFullName, "icon.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.type", "CUSTOM_APP"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.home_page_url", ""),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.grant_types.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.grant_types.*", "CLIENT_CREDENTIALS"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.response_types.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.token_endpoint_authn_method", "CLIENT_SECRET_BASIC"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.pkce_enforcement", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.redirect_uris.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.post_logout_redirect_uris.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.refresh_token_duration", "2592000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.refresh_token_rolling_duration", "15552000"),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.0.client_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.0.client_secret", regexp.MustCompile(`[a-zA-Z0-9-~_]{10,}`)),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.support_unsigned_request_object", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.bundle_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.package_name", ""),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "external_link_options.#", "0"),
				),
			},
			{
				Config: testAccApplicationConfig_OIDC_FullCustom(resourceName, name, image),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", "My test OIDC app"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "login_page_url", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "icon.#", "1"),
					resource.TestMatchResourceAttr(resourceFullName, "icon.0.id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "icon.0.href", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_role_type", "ADMIN_USERS_ONLY"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "access_control_group_options.*", map[string]string{
						"type":     "ANY_GROUP",
						"groups.#": "1",
					}),
					resource.TestMatchTypeSetElemNestedAttrs(resourceFullName, "access_control_group_options.*", map[string]*regexp.Regexp{
						"groups.0": regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
					}),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.type", "CUSTOM_APP"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.home_page_url", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.grant_types.#", "4"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.grant_types.*", "AUTHORIZATION_CODE"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.grant_types.*", "CLIENT_CREDENTIALS"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.grant_types.*", "IMPLICIT"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.grant_types.*", "REFRESH_TOKEN"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.response_types.#", "3"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.response_types.*", "CODE"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.response_types.*", "TOKEN"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.response_types.*", "ID_TOKEN"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.token_endpoint_authn_method", "CLIENT_SECRET_BASIC"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.pkce_enforcement", "REQUIRED"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.redirect_uris.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.redirect_uris.*", "https://pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.redirect_uris.*", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.post_logout_redirect_uris.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.post_logout_redirect_uris.*", "https://pingidentity.com/logout"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.post_logout_redirect_uris.*", "https://www.pingidentity.com/logout"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.refresh_token_duration", "3000000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.refresh_token_rolling_duration", "30000000"),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.0.client_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.0.client_secret", regexp.MustCompile(`[a-zA-Z0-9-~_]{10,}`)),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.support_unsigned_request_object", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.bundle_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.package_name", ""),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "external_link_options.#", "0"),
				),
			},
		},
	})
}

func TestAccApplication_OIDCFullService(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application.%s", resourceName)

	name := resourceName

	data, _ := ioutil.ReadFile("../../acctest/test_assets/image/image-logo.gif")
	image := base64.StdEncoding.EncodeToString(data)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckApplicationDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationConfig_OIDC_FullService(resourceName, name, image),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", "My test OIDC app"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "login_page_url", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "icon.#", "1"),
					resource.TestMatchResourceAttr(resourceFullName, "icon.0.id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "icon.0.href", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_role_type", "ADMIN_USERS_ONLY"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "access_control_group_options.*", map[string]string{
						"type":     "ANY_GROUP",
						"groups.#": "1",
					}),
					resource.TestMatchTypeSetElemNestedAttrs(resourceFullName, "access_control_group_options.*", map[string]*regexp.Regexp{
						"groups.0": regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
					}),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.type", "SERVICE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.home_page_url", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.grant_types.#", "4"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.grant_types.*", "AUTHORIZATION_CODE"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.grant_types.*", "CLIENT_CREDENTIALS"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.grant_types.*", "IMPLICIT"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.grant_types.*", "REFRESH_TOKEN"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.response_types.#", "3"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.response_types.*", "CODE"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.response_types.*", "TOKEN"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.response_types.*", "ID_TOKEN"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.token_endpoint_authn_method", "CLIENT_SECRET_BASIC"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.pkce_enforcement", "REQUIRED"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.redirect_uris.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.redirect_uris.*", "https://pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.redirect_uris.*", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.post_logout_redirect_uris.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.post_logout_redirect_uris.*", "https://pingidentity.com/logout"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.post_logout_redirect_uris.*", "https://www.pingidentity.com/logout"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.refresh_token_duration", "3000000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.refresh_token_rolling_duration", "30000000"),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.0.client_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.0.client_secret", regexp.MustCompile(`[a-zA-Z0-9-~_]{10,}`)),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.support_unsigned_request_object", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.bundle_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.package_name", ""),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "external_link_options.#", "0"),
				),
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
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckApplicationDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationConfig_OIDC_MinimalService(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", ""),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "login_page_url", ""),
					resource.TestCheckResourceAttr(resourceFullName, "icon.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.type", "SERVICE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.home_page_url", ""),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.grant_types.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.grant_types.*", "AUTHORIZATION_CODE"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.grant_types.*", "REFRESH_TOKEN"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.response_types.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.response_types.*", "CODE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.token_endpoint_authn_method", "CLIENT_SECRET_BASIC"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.pkce_enforcement", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.redirect_uris.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.redirect_uris.*", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.post_logout_redirect_uris.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.refresh_token_duration", "2592000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.refresh_token_rolling_duration", "15552000"),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.0.client_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.0.client_secret", regexp.MustCompile(`[a-zA-Z0-9-~_]{10,}`)),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.support_unsigned_request_object", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.bundle_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.package_name", ""),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "external_link_options.#", "0"),
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

	data, _ := ioutil.ReadFile("../../acctest/test_assets/image/image-logo.gif")
	image := base64.StdEncoding.EncodeToString(data)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckApplicationDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationConfig_OIDC_FullService(resourceName, name, image),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", "My test OIDC app"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "login_page_url", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "icon.#", "1"),
					resource.TestMatchResourceAttr(resourceFullName, "icon.0.id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "icon.0.href", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_role_type", "ADMIN_USERS_ONLY"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "access_control_group_options.*", map[string]string{
						"type":     "ANY_GROUP",
						"groups.#": "1",
					}),
					resource.TestMatchTypeSetElemNestedAttrs(resourceFullName, "access_control_group_options.*", map[string]*regexp.Regexp{
						"groups.0": regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
					}),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.type", "SERVICE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.home_page_url", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.grant_types.#", "4"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.grant_types.*", "AUTHORIZATION_CODE"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.grant_types.*", "CLIENT_CREDENTIALS"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.grant_types.*", "IMPLICIT"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.grant_types.*", "REFRESH_TOKEN"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.response_types.#", "3"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.response_types.*", "CODE"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.response_types.*", "TOKEN"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.response_types.*", "ID_TOKEN"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.token_endpoint_authn_method", "CLIENT_SECRET_BASIC"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.pkce_enforcement", "REQUIRED"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.redirect_uris.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.redirect_uris.*", "https://pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.redirect_uris.*", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.post_logout_redirect_uris.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.post_logout_redirect_uris.*", "https://pingidentity.com/logout"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.post_logout_redirect_uris.*", "https://www.pingidentity.com/logout"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.refresh_token_duration", "3000000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.refresh_token_rolling_duration", "30000000"),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.0.client_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.0.client_secret", regexp.MustCompile(`[a-zA-Z0-9-~_]{10,}`)),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.support_unsigned_request_object", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.bundle_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.package_name", ""),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "external_link_options.#", "0"),
				),
			},
			{
				Config: testAccApplicationConfig_OIDC_MinimalService(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", ""),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "login_page_url", ""),
					resource.TestCheckResourceAttr(resourceFullName, "icon.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.type", "SERVICE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.home_page_url", ""),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.grant_types.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.grant_types.*", "AUTHORIZATION_CODE"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.grant_types.*", "REFRESH_TOKEN"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.response_types.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.response_types.*", "CODE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.token_endpoint_authn_method", "CLIENT_SECRET_BASIC"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.pkce_enforcement", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.redirect_uris.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.redirect_uris.*", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.post_logout_redirect_uris.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.refresh_token_duration", "2592000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.refresh_token_rolling_duration", "15552000"),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.0.client_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.0.client_secret", regexp.MustCompile(`[a-zA-Z0-9-~_]{10,}`)),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.support_unsigned_request_object", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.bundle_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.package_name", ""),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "external_link_options.#", "0"),
				),
			},
			{
				Config: testAccApplicationConfig_OIDC_FullService(resourceName, name, image),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", "My test OIDC app"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "login_page_url", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "icon.#", "1"),
					resource.TestMatchResourceAttr(resourceFullName, "icon.0.id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "icon.0.href", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_role_type", "ADMIN_USERS_ONLY"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "access_control_group_options.*", map[string]string{
						"type":     "ANY_GROUP",
						"groups.#": "1",
					}),
					resource.TestMatchTypeSetElemNestedAttrs(resourceFullName, "access_control_group_options.*", map[string]*regexp.Regexp{
						"groups.0": regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
					}),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.type", "SERVICE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.home_page_url", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.grant_types.#", "4"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.grant_types.*", "AUTHORIZATION_CODE"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.grant_types.*", "CLIENT_CREDENTIALS"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.grant_types.*", "IMPLICIT"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.grant_types.*", "REFRESH_TOKEN"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.response_types.#", "3"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.response_types.*", "CODE"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.response_types.*", "TOKEN"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.response_types.*", "ID_TOKEN"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.token_endpoint_authn_method", "CLIENT_SECRET_BASIC"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.pkce_enforcement", "REQUIRED"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.redirect_uris.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.redirect_uris.*", "https://pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.redirect_uris.*", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.post_logout_redirect_uris.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.post_logout_redirect_uris.*", "https://pingidentity.com/logout"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.post_logout_redirect_uris.*", "https://www.pingidentity.com/logout"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.refresh_token_duration", "3000000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.refresh_token_rolling_duration", "30000000"),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.0.client_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.0.client_secret", regexp.MustCompile(`[a-zA-Z0-9-~_]{10,}`)),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.support_unsigned_request_object", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.bundle_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.package_name", ""),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "external_link_options.#", "0"),
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

	data, _ := ioutil.ReadFile("../../acctest/test_assets/image/image-logo.gif")
	image := base64.StdEncoding.EncodeToString(data)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckApplicationDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationConfig_OIDC_FullSPA(resourceName, name, image),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", "My test OIDC app"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "login_page_url", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "icon.#", "1"),
					resource.TestMatchResourceAttr(resourceFullName, "icon.0.id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "icon.0.href", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_role_type", "ADMIN_USERS_ONLY"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "access_control_group_options.*", map[string]string{
						"type":     "ANY_GROUP",
						"groups.#": "1",
					}),
					resource.TestMatchTypeSetElemNestedAttrs(resourceFullName, "access_control_group_options.*", map[string]*regexp.Regexp{
						"groups.0": regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
					}),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.type", "SINGLE_PAGE_APP"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.home_page_url", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.grant_types.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.grant_types.*", "AUTHORIZATION_CODE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.response_types.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.response_types.*", "CODE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.token_endpoint_authn_method", "NONE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.pkce_enforcement", "S256_REQUIRED"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.redirect_uris.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.redirect_uris.*", "https://pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.redirect_uris.*", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.post_logout_redirect_uris.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.post_logout_redirect_uris.*", "https://pingidentity.com/logout"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.post_logout_redirect_uris.*", "https://www.pingidentity.com/logout"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.refresh_token_duration", "2592000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.refresh_token_rolling_duration", "15552000"),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.0.client_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.0.client_secret", regexp.MustCompile(`[a-zA-Z0-9-~_]{10,}`)),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.support_unsigned_request_object", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.bundle_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.package_name", ""),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "external_link_options.#", "0"),
				),
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
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckApplicationDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationConfig_OIDC_MinimalSPA(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", ""),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "login_page_url", ""),
					resource.TestCheckResourceAttr(resourceFullName, "icon.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.type", "SINGLE_PAGE_APP"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.home_page_url", ""),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.grant_types.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.grant_types.*", "AUTHORIZATION_CODE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.response_types.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.response_types.*", "CODE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.token_endpoint_authn_method", "NONE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.pkce_enforcement", "S256_REQUIRED"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.redirect_uris.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.redirect_uris.*", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.post_logout_redirect_uris.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.refresh_token_duration", "2592000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.refresh_token_rolling_duration", "15552000"),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.0.client_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.0.client_secret", regexp.MustCompile(`[a-zA-Z0-9-~_]{10,}`)),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.support_unsigned_request_object", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.bundle_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.package_name", ""),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "external_link_options.#", "0"),
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

	data, _ := ioutil.ReadFile("../../acctest/test_assets/image/image-logo.gif")
	image := base64.StdEncoding.EncodeToString(data)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckApplicationDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationConfig_OIDC_MinimalSPA(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", ""),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "login_page_url", ""),
					resource.TestCheckResourceAttr(resourceFullName, "icon.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.type", "SINGLE_PAGE_APP"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.home_page_url", ""),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.grant_types.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.grant_types.*", "AUTHORIZATION_CODE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.response_types.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.response_types.*", "CODE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.token_endpoint_authn_method", "NONE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.pkce_enforcement", "S256_REQUIRED"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.redirect_uris.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.redirect_uris.*", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.post_logout_redirect_uris.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.refresh_token_duration", "2592000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.refresh_token_rolling_duration", "15552000"),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.0.client_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.0.client_secret", regexp.MustCompile(`[a-zA-Z0-9-~_]{10,}`)),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.support_unsigned_request_object", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.bundle_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.package_name", ""),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "external_link_options.#", "0"),
				),
			},
			{
				Config: testAccApplicationConfig_OIDC_FullSPA(resourceName, name, image),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", "My test OIDC app"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "login_page_url", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "icon.#", "1"),
					resource.TestMatchResourceAttr(resourceFullName, "icon.0.id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "icon.0.href", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_role_type", "ADMIN_USERS_ONLY"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "access_control_group_options.*", map[string]string{
						"type":     "ANY_GROUP",
						"groups.#": "1",
					}),
					resource.TestMatchTypeSetElemNestedAttrs(resourceFullName, "access_control_group_options.*", map[string]*regexp.Regexp{
						"groups.0": regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
					}),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.type", "SINGLE_PAGE_APP"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.home_page_url", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.grant_types.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.grant_types.*", "AUTHORIZATION_CODE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.response_types.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.response_types.*", "CODE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.token_endpoint_authn_method", "NONE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.pkce_enforcement", "S256_REQUIRED"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.redirect_uris.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.redirect_uris.*", "https://pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.redirect_uris.*", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.post_logout_redirect_uris.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.post_logout_redirect_uris.*", "https://pingidentity.com/logout"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.post_logout_redirect_uris.*", "https://www.pingidentity.com/logout"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.refresh_token_duration", "2592000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.refresh_token_rolling_duration", "15552000"),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.0.client_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.0.client_secret", regexp.MustCompile(`[a-zA-Z0-9-~_]{10,}`)),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.support_unsigned_request_object", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.bundle_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.package_name", ""),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "external_link_options.#", "0"),
				),
			},
			{
				Config: testAccApplicationConfig_OIDC_MinimalSPA(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", ""),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "login_page_url", ""),
					resource.TestCheckResourceAttr(resourceFullName, "icon.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.type", "SINGLE_PAGE_APP"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.home_page_url", ""),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.grant_types.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.grant_types.*", "AUTHORIZATION_CODE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.response_types.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.response_types.*", "CODE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.token_endpoint_authn_method", "NONE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.pkce_enforcement", "S256_REQUIRED"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.redirect_uris.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.redirect_uris.*", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.post_logout_redirect_uris.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.refresh_token_duration", "2592000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.refresh_token_rolling_duration", "15552000"),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.0.client_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.0.client_secret", regexp.MustCompile(`[a-zA-Z0-9-~_]{10,}`)),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.support_unsigned_request_object", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.bundle_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.package_name", ""),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "external_link_options.#", "0"),
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

	data, _ := ioutil.ReadFile("../../acctest/test_assets/image/image-logo.png")
	image := base64.StdEncoding.EncodeToString(data)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckApplicationDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationConfig_OIDC_FullWorker(resourceName, name, image),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", "My test OIDC app"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "login_page_url", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "icon.#", "1"),
					resource.TestMatchResourceAttr(resourceFullName, "icon.0.id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "icon.0.href", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_role_type", "ADMIN_USERS_ONLY"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "access_control_group_options.*", map[string]string{
						"type":     "ANY_GROUP",
						"groups.#": "2",
					}),
					resource.TestMatchTypeSetElemNestedAttrs(resourceFullName, "access_control_group_options.*", map[string]*regexp.Regexp{
						"groups.0": regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
						"groups.1": regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
					}),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.type", "WORKER"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.home_page_url", ""),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.grant_types.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.grant_types.*", "CLIENT_CREDENTIALS"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.response_types.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.token_endpoint_authn_method", "CLIENT_SECRET_BASIC"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.pkce_enforcement", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.redirect_uris.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.post_logout_redirect_uris.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.refresh_token_duration", "2592000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.refresh_token_rolling_duration", "15552000"),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.0.client_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.0.client_secret", regexp.MustCompile(`[a-zA-Z0-9-~_]{10,}`)),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.support_unsigned_request_object", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.bundle_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.package_name", ""),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "external_link_options.#", "0"),
				),
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
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckApplicationDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationConfig_OIDC_MinimalWorker(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", ""),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "login_page_url", ""),
					resource.TestCheckResourceAttr(resourceFullName, "icon.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_role_type", "ADMIN_USERS_ONLY"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.type", "WORKER"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.home_page_url", ""),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.grant_types.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.grant_types.*", "CLIENT_CREDENTIALS"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.response_types.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.token_endpoint_authn_method", "CLIENT_SECRET_BASIC"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.pkce_enforcement", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.redirect_uris.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.post_logout_redirect_uris.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.refresh_token_duration", "2592000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.refresh_token_rolling_duration", "15552000"),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.0.client_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.0.client_secret", regexp.MustCompile(`[a-zA-Z0-9-~_]{10,}`)),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.support_unsigned_request_object", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.bundle_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.package_name", ""),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "external_link_options.#", "0"),
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

	data, _ := ioutil.ReadFile("../../acctest/test_assets/image/image-logo.png")
	image := base64.StdEncoding.EncodeToString(data)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckApplicationDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationConfig_OIDC_MinimalWorker(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", ""),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "login_page_url", ""),
					resource.TestCheckResourceAttr(resourceFullName, "icon.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_role_type", "ADMIN_USERS_ONLY"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.type", "WORKER"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.home_page_url", ""),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.grant_types.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.grant_types.*", "CLIENT_CREDENTIALS"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.response_types.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.token_endpoint_authn_method", "CLIENT_SECRET_BASIC"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.pkce_enforcement", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.redirect_uris.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.post_logout_redirect_uris.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.refresh_token_duration", "2592000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.refresh_token_rolling_duration", "15552000"),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.0.client_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.0.client_secret", regexp.MustCompile(`[a-zA-Z0-9-~_]{10,}`)),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.support_unsigned_request_object", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.bundle_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.package_name", ""),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "external_link_options.#", "0"),
				),
			},
			{
				Config: testAccApplicationConfig_OIDC_FullWorker(resourceName, name, image),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", "My test OIDC app"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "login_page_url", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "icon.#", "1"),
					resource.TestMatchResourceAttr(resourceFullName, "icon.0.id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "icon.0.href", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_role_type", "ADMIN_USERS_ONLY"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "access_control_group_options.*", map[string]string{
						"type":     "ANY_GROUP",
						"groups.#": "2",
					}),
					resource.TestMatchTypeSetElemNestedAttrs(resourceFullName, "access_control_group_options.*", map[string]*regexp.Regexp{
						"groups.0": regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
						"groups.1": regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
					}),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.type", "WORKER"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.home_page_url", ""),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.grant_types.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.grant_types.*", "CLIENT_CREDENTIALS"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.response_types.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.token_endpoint_authn_method", "CLIENT_SECRET_BASIC"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.pkce_enforcement", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.redirect_uris.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.post_logout_redirect_uris.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.refresh_token_duration", "2592000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.refresh_token_rolling_duration", "15552000"),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.0.client_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.0.client_secret", regexp.MustCompile(`[a-zA-Z0-9-~_]{10,}`)),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.support_unsigned_request_object", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.bundle_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.package_name", ""),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "external_link_options.#", "0"),
				),
			},
			{
				Config: testAccApplicationConfig_OIDC_MinimalWorker(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", ""),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "login_page_url", ""),
					resource.TestCheckResourceAttr(resourceFullName, "icon.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_role_type", "ADMIN_USERS_ONLY"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.type", "WORKER"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.home_page_url", ""),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.grant_types.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "oidc_options.0.grant_types.*", "CLIENT_CREDENTIALS"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.response_types.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.token_endpoint_authn_method", "CLIENT_SECRET_BASIC"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.pkce_enforcement", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.redirect_uris.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.post_logout_redirect_uris.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.refresh_token_duration", "2592000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.refresh_token_rolling_duration", "15552000"),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.0.client_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.0.client_secret", regexp.MustCompile(`[a-zA-Z0-9-~_]{10,}`)),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.support_unsigned_request_object", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.bundle_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.package_name", ""),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "external_link_options.#", "0"),
				),
			},
		},
	})
}

func TestAccApplication_SAMLFull(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application.%s", resourceName)

	name := resourceName

	data, _ := ioutil.ReadFile("../../acctest/test_assets/image/image-logo.gif")
	image := base64.StdEncoding.EncodeToString(data)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckApplicationDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationConfig_SAML_Full(resourceName, name, image),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", "My test SAML app"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "login_page_url", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "icon.#", "1"),
					resource.TestMatchResourceAttr(resourceFullName, "icon.0.id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "icon.0.href", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_role_type", "ADMIN_USERS_ONLY"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "access_control_group_options.*", map[string]string{
						"type":     "ANY_GROUP",
						"groups.#": "2",
					}),
					resource.TestMatchTypeSetElemNestedAttrs(resourceFullName, "access_control_group_options.*", map[string]*regexp.Regexp{
						"groups.0": regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
						"groups.1": regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
					}),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.0.type", "WEB_APP"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.0.acs_urls.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "saml_options.0.acs_urls.*", "https://pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "saml_options.0.acs_urls.*", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.0.assertion_duration", "3600"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.0.assertion_signed_enabled", "false"),
					resource.TestMatchResourceAttr(resourceFullName, "saml_options.0.idp_signing_key_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.0.nameid_format", "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.0.response_is_signed", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.0.slo_binding", "HTTP_REDIRECT"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.0.slo_endpoint", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.0.slo_response_endpoint", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.0.sp_entity_id", fmt.Sprintf("sp:entity:%s", resourceName)),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.0.sp_verification_certificate_ids.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "external_link_options.#", "0"),
				),
			},
		},
	})
}

func TestAccApplication_SAMLMinimal(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckApplicationDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationConfig_SAML_Minimal(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", ""),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "login_page_url", ""),
					resource.TestCheckResourceAttr(resourceFullName, "icon.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.0.type", "WEB_APP"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.0.acs_urls.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "saml_options.0.acs_urls.*", "https://pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.0.assertion_duration", "3600"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.0.assertion_signed_enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.0.idp_signing_key_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.0.nameid_format", ""),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.0.response_is_signed", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.0.slo_binding", "HTTP_POST"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.0.slo_endpoint", ""),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.0.slo_response_endpoint", ""),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.0.sp_entity_id", fmt.Sprintf("sp:entity:%s", resourceName)),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.0.sp_verification_certificate_ids.#", "0"),
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

	data, _ := ioutil.ReadFile("../../acctest/test_assets/image/image-logo.gif")
	image := base64.StdEncoding.EncodeToString(data)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckApplicationDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationConfig_ExternalLinkFull(resourceName, name, image),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", "My test external link app"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "icon.#", "1"),
					resource.TestMatchResourceAttr(resourceFullName, "icon.0.id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "icon.0.href", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "access_control_group_options.*", map[string]string{
						"type":     "ANY_GROUP",
						"groups.#": "2",
					}),
					resource.TestMatchTypeSetElemNestedAttrs(resourceFullName, "access_control_group_options.*", map[string]*regexp.Regexp{
						"groups.0": regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
						"groups.1": regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
					}),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "external_link_options.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "external_link_options.0.home_page_url", "https://www.pingidentity.com"),
				),
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
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckApplicationDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationConfig_ExternalLinkMinimal(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", ""),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "icon.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "external_link_options.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "external_link_options.0.home_page_url", "https://www.pingidentity.com"),
				),
			},
		},
	})
}

func TestAccApplication_Enabled(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckApplicationDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
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

func testAccApplicationConfig_NewEnv(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s"
  enabled        = true

  oidc_options {
    type                        = "WEB_APP"
    grant_types                 = ["AUTHORIZATION_CODE", "REFRESH_TOKEN"]
    response_types              = ["CODE"]
    token_endpoint_authn_method = "CLIENT_SECRET_BASIC"
    redirect_uris               = ["https://www.pingidentity.com"]
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


  enabled = true

  oidc_options {
    type                        = "WEB_APP"
    grant_types                 = ["REFRESH_TOKEN", "AUTHORIZATION_CODE"]
    response_types              = ["CODE"]
    token_endpoint_authn_method = "CLIENT_SECRET_BASIC"
    redirect_uris               = ["https://www.pingidentity.com", "https://pingidentity.com"]
    post_logout_redirect_uris   = ["https://www.pingidentity.com/logout", "https://pingidentity.com/logout"]

    refresh_token_duration         = 3000000
    refresh_token_rolling_duration = 30000000

    home_page_url    = "https://www.pingidentity.com"
    pkce_enforcement = "OPTIONAL"


    support_unsigned_request_object = true

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

  oidc_options {
    type                        = "WEB_APP"
    grant_types                 = ["REFRESH_TOKEN", "AUTHORIZATION_CODE"]
    response_types              = ["CODE"]
    token_endpoint_authn_method = "CLIENT_SECRET_BASIC"
    redirect_uris               = ["https://www.pingidentity.com"]
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

  enabled = true

  oidc_options {
    type                        = "NATIVE_APP"
    grant_types                 = ["CLIENT_CREDENTIALS"]
    token_endpoint_authn_method = "CLIENT_SECRET_BASIC"

    mobile_app {
      bundle_id    = "com.%[2]s.bundle"
      package_name = "com.%[2]s.package"

      passcode_refresh_seconds = 45

      universal_app_link = "https://applink.com"

      integrity_detection {
        enabled = true

        excluded_platforms = ["IOS"]

        cache_duration {
          amount = 30
          units  = "HOURS"
        }
      }
    }

    bundle_id    = "com.%[2]s.bundle"
    package_name = "com.%[2]s.package"
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

  oidc_options {
    type                        = "NATIVE_APP"
    grant_types                 = ["CLIENT_CREDENTIALS"]
    token_endpoint_authn_method = "CLIENT_SECRET_BASIC"
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

  enabled = true

  oidc_options {
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
    token_endpoint_authn_method = "CLIENT_SECRET_BASIC"
    redirect_uris               = ["https://www.pingidentity.com", "https://pingidentity.com"]
    post_logout_redirect_uris   = ["https://www.pingidentity.com/logout", "https://pingidentity.com/logout"]

    refresh_token_duration         = 3000000
    refresh_token_rolling_duration = 30000000

    home_page_url    = "https://www.pingidentity.com"
    pkce_enforcement = "REQUIRED"
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

  oidc_options {
    type                        = "CUSTOM_APP"
    grant_types                 = ["CLIENT_CREDENTIALS"]
    token_endpoint_authn_method = "CLIENT_SECRET_BASIC"
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

  enabled = true

  oidc_options {
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
    token_endpoint_authn_method = "CLIENT_SECRET_BASIC"
    redirect_uris               = ["https://www.pingidentity.com", "https://pingidentity.com"]
    post_logout_redirect_uris   = ["https://www.pingidentity.com/logout", "https://pingidentity.com/logout"]

    refresh_token_duration         = 3000000
    refresh_token_rolling_duration = 30000000

    home_page_url    = "https://www.pingidentity.com"
    pkce_enforcement = "REQUIRED"
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

  oidc_options {
    type                        = "SERVICE"
    grant_types                 = ["REFRESH_TOKEN", "AUTHORIZATION_CODE"]
    response_types              = ["CODE"]
    token_endpoint_authn_method = "CLIENT_SECRET_BASIC"
    redirect_uris               = ["https://www.pingidentity.com"]
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

  enabled = true

  oidc_options {
    type                        = "SINGLE_PAGE_APP"
    grant_types                 = ["AUTHORIZATION_CODE"]
    response_types              = ["CODE"]
    pkce_enforcement            = "S256_REQUIRED"
    token_endpoint_authn_method = "NONE"
    redirect_uris               = ["https://www.pingidentity.com", "https://pingidentity.com"]
    post_logout_redirect_uris   = ["https://www.pingidentity.com/logout", "https://pingidentity.com/logout"]
    home_page_url               = "https://www.pingidentity.com"

    support_unsigned_request_object = true

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

  oidc_options {
    type                        = "SINGLE_PAGE_APP"
    grant_types                 = ["AUTHORIZATION_CODE"]
    response_types              = ["CODE"]
    pkce_enforcement            = "S256_REQUIRED"
    token_endpoint_authn_method = "NONE"
    redirect_uris               = ["https://www.pingidentity.com"]
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
      pingone_group.%[2]s-2.id,
      pingone_group.%[2]s-1.id
    ]
  }

  enabled = true

  oidc_options {
    type                        = "WORKER"
    grant_types                 = ["CLIENT_CREDENTIALS"]
    token_endpoint_authn_method = "CLIENT_SECRET_BASIC"
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

  oidc_options {
    type                        = "WORKER"
    grant_types                 = ["CLIENT_CREDENTIALS"]
    token_endpoint_authn_method = "CLIENT_SECRET_BASIC"
  }
}
		`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccApplicationConfig_SAML_Full(resourceName, name, image string) string {
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
  signature_algorithm = "SHA224withECDSA"
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
      pingone_group.%[2]s-2.id,
      pingone_group.%[2]s-1.id
    ]
  }

  enabled = true

  saml_options {
    type               = "WEB_APP"
    acs_urls           = ["https://www.pingidentity.com", "https://pingidentity.com"]
    assertion_duration = 3600
    sp_entity_id       = "sp:entity:%[2]s"

    assertion_signed_enabled = false
    idp_signing_key_id       = pingone_key.%[2]s.id
    nameid_format            = "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress"
    response_is_signed       = true
    slo_binding              = "HTTP_REDIRECT"
    slo_endpoint             = "https://www.pingidentity.com"
    slo_response_endpoint    = "https://www.pingidentity.com"

    // sp_verification_certificate_ids = []

  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name, image)
}

func testAccApplicationConfig_SAML_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true

  saml_options {
    acs_urls           = ["https://pingidentity.com"]
    assertion_duration = 3600
    sp_entity_id       = "sp:entity:%[2]s"
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
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

  enabled = true

  external_link_options {
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

  external_link_options {
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

  oidc_options {
    type                        = "SINGLE_PAGE_APP"
    grant_types                 = ["AUTHORIZATION_CODE"]
    response_types              = ["CODE"]
    pkce_enforcement            = "S256_REQUIRED"
    token_endpoint_authn_method = "NONE"
    redirect_uris               = ["https://www.pingidentity.com"]
  }
}
		`, acctest.GenericSandboxEnvironment(), resourceName, name, enabled)
}

// Error conditions

// func testAccApplicationConfig_NoType(resourceName, name string) string {
// 	return fmt.Sprintf(`
// 		%[1]s
// 		resource "pingone_application" "%[3]s" {
// 			environment_id = data.pingone_environment.general_test.id
// 			name = "%[3]s"
// 		}`, acctest.GenericSandboxEnvironment(), resourceName, name)
// }
