package sso_test

import (
	"context"
	"fmt"
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

		if rEnv.StatusCode == 404 {
			continue
		}

		if err != nil {
			return err
		}

		body, r, err := apiClient.ApplicationsApi.ReadOneApplication(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.ID).Execute()

		if r.StatusCode == 404 {
			continue
		}

		if err != nil {
			tflog.Error(ctx, fmt.Sprintf("Error: %v", body))
			return err
		}

		return fmt.Errorf("PingOne Application Instance %s still exists", rs.Primary.ID)
	}

	return nil
}

func TestAccApplication_OIDCFullWeb(t *testing.T) {
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
				Config: testAccApplicationConfig_OIDCFullWeb(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", "My test OIDC app"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "login_page_url", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "icon.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_role_type", "ADMIN_USERS_ONLY"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.0.type", "ANY_GROUP"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.0.groups.#", "1"),
					resource.TestMatchResourceAttr(resourceFullName, "access_control_group_options.0.groups.0", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.type", "WEB_APP"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.home_page_url", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.grant_types.#", "2"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.grant_types.0", "AUTHORIZATION_CODE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.grant_types.1", "REFRESH_TOKEN"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.response_types.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.response_types.0", "CODE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.token_endpoint_authn_method", "CLIENT_SECRET_BASIC"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.pkce_enforcement", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.redirect_uris.#", "2"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.redirect_uris.0", "https://pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.redirect_uris.1", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.post_logout_redirect_uris.#", "2"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.post_logout_redirect_uris.0", "https://pingidentity.com/logout"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.post_logout_redirect_uris.1", "https://www.pingidentity.com/logout"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.refresh_token_duration", "3000000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.refresh_token_rolling_duration", "30000000"),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.0.client_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.0.client_secret", regexp.MustCompile(`[a-zA-Z0-9-~_]{10,}`)),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.support_unsigned_request_object", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.bundle_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.package_name", ""),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.#", "0"),
				),
			},
		},
	})
}

func TestAccApplication_OIDCMinimalWeb(t *testing.T) {
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
				Config: testAccApplicationConfig_OIDCMinimalWeb(environmentName, licenseID, resourceName, name),
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
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.grant_types.0", "AUTHORIZATION_CODE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.grant_types.1", "REFRESH_TOKEN"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.response_types.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.response_types.0", "CODE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.token_endpoint_authn_method", "CLIENT_SECRET_BASIC"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.pkce_enforcement", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.redirect_uris.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.redirect_uris.0", "https://www.pingidentity.com"),
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
				),
			},
		},
	})
}

func TestAccApplication_OIDCWebUpdate(t *testing.T) {
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
				Config: testAccApplicationConfig_OIDCMinimalWeb(environmentName, licenseID, resourceName, name),
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
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.grant_types.0", "AUTHORIZATION_CODE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.grant_types.1", "REFRESH_TOKEN"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.response_types.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.response_types.0", "CODE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.token_endpoint_authn_method", "CLIENT_SECRET_BASIC"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.pkce_enforcement", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.redirect_uris.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.redirect_uris.0", "https://www.pingidentity.com"),
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
				),
			},
			{
				Config: testAccApplicationConfig_OIDCFullWeb(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", "My test OIDC app"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "login_page_url", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "icon.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_role_type", "ADMIN_USERS_ONLY"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.0.type", "ANY_GROUP"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.0.groups.#", "1"),
					resource.TestMatchResourceAttr(resourceFullName, "access_control_group_options.0.groups.0", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.type", "WEB_APP"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.home_page_url", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.grant_types.#", "2"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.grant_types.0", "AUTHORIZATION_CODE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.grant_types.1", "REFRESH_TOKEN"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.response_types.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.response_types.0", "CODE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.token_endpoint_authn_method", "CLIENT_SECRET_BASIC"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.pkce_enforcement", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.redirect_uris.#", "2"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.redirect_uris.0", "https://pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.redirect_uris.1", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.post_logout_redirect_uris.#", "2"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.post_logout_redirect_uris.0", "https://pingidentity.com/logout"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.post_logout_redirect_uris.1", "https://www.pingidentity.com/logout"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.refresh_token_duration", "3000000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.refresh_token_rolling_duration", "30000000"),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.0.client_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.0.client_secret", regexp.MustCompile(`[a-zA-Z0-9-~_]{10,}`)),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.support_unsigned_request_object", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.bundle_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.package_name", ""),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.#", "0"),
				),
			},
			{
				Config: testAccApplicationConfig_OIDCMinimalWeb(environmentName, licenseID, resourceName, name),
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
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.grant_types.0", "AUTHORIZATION_CODE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.grant_types.1", "REFRESH_TOKEN"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.response_types.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.response_types.0", "CODE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.token_endpoint_authn_method", "CLIENT_SECRET_BASIC"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.pkce_enforcement", "OPTIONAL"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.redirect_uris.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.redirect_uris.0", "https://www.pingidentity.com"),
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
				),
			},
		},
	})
}

func TestAccApplication_OIDCFullNative(t *testing.T) {
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
				Config: testAccApplicationConfig_OIDCFullNative(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", "My test OIDC app"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "login_page_url", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "icon.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_role_type", "ADMIN_USERS_ONLY"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.0.type", "ANY_GROUP"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.0.groups.#", "1"),
					resource.TestMatchResourceAttr(resourceFullName, "access_control_group_options.0.groups.0", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.type", "NATIVE_APP"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.home_page_url", ""),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.grant_types.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.grant_types.0", "CLIENT_CREDENTIALS"),
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
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.bundle_id", "com.pingidentity.bundle"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.package_name", "com.pingidentity.package"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.passcode_refresh_seconds", "45"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.integrity_detection.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.integrity_detection.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.integrity_detection.0.cache_duration.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.integrity_detection.0.cache_duration.0.amount", "30"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.integrity_detection.0.cache_duration.0.units", "HOURS"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.bundle_id", "com.pingidentity.bundle"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.package_name", "com.pingidentity.package"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.#", "0"),
				),
			},
		},
	})
}

func TestAccApplication_OIDCMinimalNative(t *testing.T) {
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
				Config: testAccApplicationConfig_OIDCMinimalNative(environmentName, licenseID, resourceName, name),
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
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.grant_types.0", "CLIENT_CREDENTIALS"),
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
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.integrity_detection.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.integrity_detection.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.passcode_refresh_seconds", "30"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.bundle_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.package_name", ""),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.#", "0"),
				),
			},
		},
	})
}

func TestAccApplication_OIDCNativeUpdate(t *testing.T) {
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
				Config: testAccApplicationConfig_OIDCFullNative(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", "My test OIDC app"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "login_page_url", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "icon.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_role_type", "ADMIN_USERS_ONLY"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.0.type", "ANY_GROUP"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.0.groups.#", "1"),
					resource.TestMatchResourceAttr(resourceFullName, "access_control_group_options.0.groups.0", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.type", "NATIVE_APP"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.home_page_url", ""),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.grant_types.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.grant_types.0", "CLIENT_CREDENTIALS"),
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
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.bundle_id", "com.pingidentity.bundle"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.package_name", "com.pingidentity.package"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.passcode_refresh_seconds", "45"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.integrity_detection.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.integrity_detection.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.integrity_detection.0.cache_duration.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.integrity_detection.0.cache_duration.0.amount", "30"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.integrity_detection.0.cache_duration.0.units", "HOURS"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.bundle_id", "com.pingidentity.bundle"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.package_name", "com.pingidentity.package"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.#", "0"),
				),
			},
			{
				Config: testAccApplicationConfig_OIDCMinimalNative(environmentName, licenseID, resourceName, name),
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
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.grant_types.0", "CLIENT_CREDENTIALS"),
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
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.integrity_detection.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.integrity_detection.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.passcode_refresh_seconds", "30"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.bundle_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.package_name", ""),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.#", "0"),
				),
			},
			{
				Config: testAccApplicationConfig_OIDCFullNative(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", "My test OIDC app"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "login_page_url", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "icon.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_role_type", "ADMIN_USERS_ONLY"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.0.type", "ANY_GROUP"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.0.groups.#", "1"),
					resource.TestMatchResourceAttr(resourceFullName, "access_control_group_options.0.groups.0", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.type", "NATIVE_APP"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.home_page_url", ""),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.grant_types.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.grant_types.0", "CLIENT_CREDENTIALS"),
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
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.bundle_id", "com.pingidentity.bundle"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.package_name", "com.pingidentity.package"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.passcode_refresh_seconds", "45"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.integrity_detection.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.integrity_detection.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.integrity_detection.0.cache_duration.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.integrity_detection.0.cache_duration.0.amount", "30"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.0.integrity_detection.0.cache_duration.0.units", "HOURS"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.bundle_id", "com.pingidentity.bundle"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.package_name", "com.pingidentity.package"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.#", "0"),
				),
			},
		},
	})
}

func TestAccApplication_OIDCFullSPA(t *testing.T) {
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
				Config: testAccApplicationConfig_OIDCFullSPA(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", "My test OIDC app"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "login_page_url", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "icon.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_role_type", "ADMIN_USERS_ONLY"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.0.type", "ANY_GROUP"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.0.groups.#", "1"),
					resource.TestMatchResourceAttr(resourceFullName, "access_control_group_options.0.groups.0", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.type", "SINGLE_PAGE_APP"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.home_page_url", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.grant_types.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.grant_types.0", "AUTHORIZATION_CODE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.response_types.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.response_types.0", "CODE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.token_endpoint_authn_method", "NONE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.pkce_enforcement", "S256_REQUIRED"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.redirect_uris.#", "2"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.redirect_uris.0", "https://pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.redirect_uris.1", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.post_logout_redirect_uris.#", "2"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.post_logout_redirect_uris.0", "https://pingidentity.com/logout"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.post_logout_redirect_uris.1", "https://www.pingidentity.com/logout"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.refresh_token_duration", "2592000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.refresh_token_rolling_duration", "15552000"),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.0.client_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.0.client_secret", regexp.MustCompile(`[a-zA-Z0-9-~_]{10,}`)),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.support_unsigned_request_object", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.bundle_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.package_name", ""),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.#", "0"),
				),
			},
		},
	})
}

func TestAccApplication_OIDCMinimalSPA(t *testing.T) {
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
				Config: testAccApplicationConfig_OIDCMinimalSPA(environmentName, licenseID, resourceName, name),
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
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.grant_types.0", "AUTHORIZATION_CODE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.response_types.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.response_types.0", "CODE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.token_endpoint_authn_method", "NONE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.pkce_enforcement", "S256_REQUIRED"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.redirect_uris.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.redirect_uris.0", "https://www.pingidentity.com"),
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
				),
			},
		},
	})
}

func TestAccApplication_OIDCSPAUpdate(t *testing.T) {
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
				Config: testAccApplicationConfig_OIDCMinimalSPA(environmentName, licenseID, resourceName, name),
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
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.grant_types.0", "AUTHORIZATION_CODE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.response_types.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.response_types.0", "CODE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.token_endpoint_authn_method", "NONE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.pkce_enforcement", "S256_REQUIRED"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.redirect_uris.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.redirect_uris.0", "https://www.pingidentity.com"),
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
				),
			},
			{
				Config: testAccApplicationConfig_OIDCFullSPA(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", "My test OIDC app"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "login_page_url", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "icon.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_role_type", "ADMIN_USERS_ONLY"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.0.type", "ANY_GROUP"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.0.groups.#", "1"),
					resource.TestMatchResourceAttr(resourceFullName, "access_control_group_options.0.groups.0", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.type", "SINGLE_PAGE_APP"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.home_page_url", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.grant_types.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.grant_types.0", "AUTHORIZATION_CODE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.response_types.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.response_types.0", "CODE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.token_endpoint_authn_method", "NONE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.pkce_enforcement", "S256_REQUIRED"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.redirect_uris.#", "2"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.redirect_uris.0", "https://pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.redirect_uris.1", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.post_logout_redirect_uris.#", "2"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.post_logout_redirect_uris.0", "https://pingidentity.com/logout"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.post_logout_redirect_uris.1", "https://www.pingidentity.com/logout"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.refresh_token_duration", "2592000"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.refresh_token_rolling_duration", "15552000"),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.0.client_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.0.client_secret", regexp.MustCompile(`[a-zA-Z0-9-~_]{10,}`)),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.support_unsigned_request_object", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.mobile_app.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.bundle_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.package_name", ""),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.#", "0"),
				),
			},
			{
				Config: testAccApplicationConfig_OIDCMinimalSPA(environmentName, licenseID, resourceName, name),
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
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.grant_types.0", "AUTHORIZATION_CODE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.response_types.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.response_types.0", "CODE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.token_endpoint_authn_method", "NONE"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.pkce_enforcement", "S256_REQUIRED"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.redirect_uris.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.redirect_uris.0", "https://www.pingidentity.com"),
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
				),
			},
		},
	})
}

func TestAccApplication_OIDCFullWorker(t *testing.T) {
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
				Config: testAccApplicationConfig_OIDCFullWorker(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", "My test OIDC app"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "login_page_url", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "icon.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_role_type", "ADMIN_USERS_ONLY"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.0.type", "ANY_GROUP"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.0.groups.#", "1"),
					resource.TestMatchResourceAttr(resourceFullName, "access_control_group_options.0.groups.0", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.type", "WORKER"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.home_page_url", ""),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.grant_types.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.grant_types.0", "CLIENT_CREDENTIALS"),
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
				),
			},
		},
	})
}

func TestAccApplication_OIDCMinimalWorker(t *testing.T) {
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
				Config: testAccApplicationConfig_OIDCMinimalWorker(environmentName, licenseID, resourceName, name),
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
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.grant_types.0", "CLIENT_CREDENTIALS"),
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
				),
			},
		},
	})
}

func TestAccApplication_OIDCWorkerUpdate(t *testing.T) {
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
				Config: testAccApplicationConfig_OIDCMinimalWorker(environmentName, licenseID, resourceName, name),
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
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.grant_types.0", "CLIENT_CREDENTIALS"),
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
				),
			},
			{
				Config: testAccApplicationConfig_OIDCFullWorker(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", "My test OIDC app"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "login_page_url", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "icon.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_role_type", "ADMIN_USERS_ONLY"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.0.type", "ANY_GROUP"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.0.groups.#", "1"),
					resource.TestMatchResourceAttr(resourceFullName, "access_control_group_options.0.groups.0", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.type", "WORKER"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.home_page_url", ""),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.grant_types.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.grant_types.0", "CLIENT_CREDENTIALS"),
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
				),
			},
			{
				Config: testAccApplicationConfig_OIDCMinimalWorker(environmentName, licenseID, resourceName, name),
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
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.0.grant_types.0", "CLIENT_CREDENTIALS"),
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
				),
			},
		},
	})
}

func TestAccApplication_SAMLFull(t *testing.T) {
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
				Config: testAccApplicationConfig_SAMLFull(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", "My test SAML app"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "tags.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "login_page_url", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "icon.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_role_type", "ADMIN_USERS_ONLY"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.0.type", "ANY_GROUP"),
					resource.TestCheckResourceAttr(resourceFullName, "access_control_group_options.0.groups.#", "1"),
					resource.TestMatchResourceAttr(resourceFullName, "access_control_group_options.0.groups.0", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.0.type", "WEB_APP"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.0.acs_urls.#", "2"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.0.acs_urls.0", "https://pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.0.acs_urls.1", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.0.assertion_duration", "3600"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.0.assertion_signed_enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.0.idp_signing_key_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.0.nameid_format", "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.0.response_is_signed", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.0.slo_binding", "HTTP_REDIRECT"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.0.slo_endpoint", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.0.slo_response_endpoint", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.0.sp_entity_id", "sp:entity:localhost"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.0.sp_verification_certificate_ids.#", "0"),
				),
			},
		},
	})
}

func TestAccApplication_Enabled(t *testing.T) {
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
				Config: testAccApplicationConfig_Enabled(environmentName, licenseID, resourceName, name, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "false"),
				),
			},
			{
				Config: testAccApplicationConfig_Enabled(environmentName, licenseID, resourceName, name, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
				),
			},
			{
				Config: testAccApplicationConfig_Enabled(environmentName, licenseID, resourceName, name, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "false"),
				),
			},
		},
	})
}

// Application enabled / disabled

func TestAccApplication_SAMLMinimal(t *testing.T) {
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
				Config: testAccApplicationConfig_SAMLMinimal(environmentName, licenseID, resourceName, name),
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
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.0.acs_urls.0", "https://pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.0.assertion_duration", "3600"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.0.assertion_signed_enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.0.idp_signing_key_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.0.nameid_format", ""),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.0.response_is_signed", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.0.slo_binding", "HTTP_POST"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.0.slo_endpoint", ""),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.0.slo_response_endpoint", ""),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.0.sp_entity_id", "sp:entity:localhost"),
					resource.TestCheckResourceAttr(resourceFullName, "saml_options.0.sp_verification_certificate_ids.#", "0"),
				),
			},
		},
	})
}

func testAccApplicationConfig_OIDCFullWeb(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s
		resource "pingone_group" "%[3]s" {
			environment_id = "${pingone_environment.%[2]s.id}"
			name = "Group 1"
		}
		resource "pingone_application" "%[3]s" {
			environment_id = "${pingone_environment.%[2]s.id}"
			name = "%[4]s"
			description = "My test OIDC app"
			tags = []
			login_page_url = "https://www.pingidentity.com"

			// icon {
			// 	id = "1"
			// 	href = "https://assets.pingone.com/ux/ui-library/4.18.0/images/logo-pingidentity.png"
			// }

			access_control_role_type = "ADMIN_USERS_ONLY"

			access_control_group_options {
				type = "ANY_GROUP"

				groups = [
					"${pingone_group.%[3]s.id}"
				]
			}
			
			
			enabled = true

			oidc_options {
				type                        = "WEB_APP"
				grant_types                 = ["AUTHORIZATION_CODE", "REFRESH_TOKEN"]
				response_types              = ["CODE"]
				token_endpoint_authn_method = "CLIENT_SECRET_BASIC"
				redirect_uris               = ["https://pingidentity.com", "https://www.pingidentity.com"]
				post_logout_redirect_uris   = ["https://pingidentity.com/logout", "https://www.pingidentity.com/logout"]

				refresh_token_duration          = 3000000
				refresh_token_rolling_duration  = 30000000

				home_page_url 		= "https://www.pingidentity.com"
				pkce_enforcement 	= "OPTIONAL"


				support_unsigned_request_object = true

			}
		}
		`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccApplicationConfig_OIDCMinimalWeb(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s
		resource "pingone_application" "%[3]s" {
			environment_id = "${pingone_environment.%[2]s.id}"
			name = "%[4]s"
			enabled = true

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

func testAccApplicationConfig_OIDCFullNative(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s
		resource "pingone_group" "%[3]s" {
			environment_id = "${pingone_environment.%[2]s.id}"
			name = "Group 1"
		}
		resource "pingone_application" "%[3]s" {
			environment_id = "${pingone_environment.%[2]s.id}"
			name = "%[4]s"
			description = "My test OIDC app"
			tags = []
			login_page_url = "https://www.pingidentity.com"

			// icon {
			// 	id = "1"
			// 	href = "https://assets.pingone.com/ux/ui-library/4.18.0/images/logo-pingidentity.png"
			// }

			access_control_role_type = "ADMIN_USERS_ONLY"

			access_control_group_options {
				type = "ANY_GROUP"

				groups = [
					"${pingone_group.%[3]s.id}"
				]
			}
			
			enabled = true

			oidc_options {
				type                        = "NATIVE_APP"
				grant_types                 = ["CLIENT_CREDENTIALS"]
				token_endpoint_authn_method = "CLIENT_SECRET_BASIC"

				mobile_app {
					bundle_id 		= "com.pingidentity.bundle"
					package_name 	= "com.pingidentity.package"

					passcode_refresh_seconds = 45

					integrity_detection {
						enabled = true
						cache_duration {
							amount = 30
							units  = "HOURS"
						}
					}
				}

				bundle_id 		= "com.pingidentity.bundle"
				package_name 	= "com.pingidentity.package"
			}
		}
		`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccApplicationConfig_OIDCMinimalNative(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s
		resource "pingone_application" "%[3]s" {
			environment_id = "${pingone_environment.%[2]s.id}"
			name = "%[4]s"
			enabled = true

			oidc_options {
				type                        = "NATIVE_APP"
				grant_types                 = ["CLIENT_CREDENTIALS"]
				token_endpoint_authn_method = "CLIENT_SECRET_BASIC"
			}
		}
		`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccApplicationConfig_OIDCFullSPA(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s
		resource "pingone_group" "%[3]s" {
			environment_id = "${pingone_environment.%[2]s.id}"
			name = "Group 1"
		}
		resource "pingone_application" "%[3]s" {
			environment_id = "${pingone_environment.%[2]s.id}"
			name = "%[4]s"
			description = "My test OIDC app"
			tags = []
			login_page_url = "https://www.pingidentity.com"

			// icon {
			// 	id = "1"
			// 	href = "https://assets.pingone.com/ux/ui-library/4.18.0/images/logo-pingidentity.png"
			// }

			access_control_role_type = "ADMIN_USERS_ONLY"

			access_control_group_options {
				type = "ANY_GROUP"

				groups = [
					"${pingone_group.%[3]s.id}"
				]
			}
			
			enabled = true

			oidc_options {
				type                        = "SINGLE_PAGE_APP"
				grant_types                 = ["AUTHORIZATION_CODE"]
				response_types              = ["CODE"]
				pkce_enforcement            = "S256_REQUIRED"
				token_endpoint_authn_method = "NONE"
				redirect_uris               = ["https://pingidentity.com", "https://www.pingidentity.com"]
				post_logout_redirect_uris   = ["https://pingidentity.com/logout", "https://www.pingidentity.com/logout"]
				home_page_url				= "https://www.pingidentity.com"

				support_unsigned_request_object = true

			}
		}
		`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccApplicationConfig_OIDCMinimalSPA(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s
		resource "pingone_application" "%[3]s" {
			environment_id = "${pingone_environment.%[2]s.id}"
			name = "%[4]s"
			enabled = true

			oidc_options {
				type                        = "SINGLE_PAGE_APP"
				grant_types                 = ["AUTHORIZATION_CODE"]
				response_types              = ["CODE"]
				pkce_enforcement            = "S256_REQUIRED"
				token_endpoint_authn_method = "NONE"
				redirect_uris               = ["https://www.pingidentity.com"]
			}
		}
		`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccApplicationConfig_OIDCFullWorker(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s
		resource "pingone_group" "%[3]s" {
			environment_id = "${pingone_environment.%[2]s.id}"
			name = "Group 1"
		}
		resource "pingone_application" "%[3]s" {
			environment_id = "${pingone_environment.%[2]s.id}"
			name = "%[4]s"
			description = "My test OIDC app"
			tags = []
			login_page_url = "https://www.pingidentity.com"

			// icon {
			// 	id = "1"
			// 	href = "https://assets.pingone.com/ux/ui-library/4.18.0/images/logo-pingidentity.png"
			// }

			access_control_role_type = "ADMIN_USERS_ONLY"

			access_control_group_options {
				type = "ANY_GROUP"

				groups = [
					"${pingone_group.%[3]s.id}"
				]
			}
			
			enabled = true

			oidc_options {
				type                        = "WORKER"
				grant_types                 = ["CLIENT_CREDENTIALS"]
				token_endpoint_authn_method = "CLIENT_SECRET_BASIC"
			}
		}
		`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccApplicationConfig_OIDCMinimalWorker(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s
		resource "pingone_application" "%[3]s" {
			environment_id = "${pingone_environment.%[2]s.id}"
			name = "%[4]s"
			enabled = true

			oidc_options {
				type                        = "WORKER"
				grant_types                 = ["CLIENT_CREDENTIALS"]
				token_endpoint_authn_method = "CLIENT_SECRET_BASIC"
			}
		}
		`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccApplicationConfig_SAMLFull(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s
		resource "pingone_group" "%[3]s" {
			environment_id = "${pingone_environment.%[2]s.id}"
			name = "Group 1"
		}
		resource "pingone_application" "%[3]s" {
			environment_id = "${pingone_environment.%[2]s.id}"
			name = "%[4]s"
			description = "My test SAML app"
			tags = []
			login_page_url = "https://www.pingidentity.com"

			// icon {
			// 	id = "1"
			// 	href = "https://assets.pingone.com/ux/ui-library/4.18.0/images/logo-pingidentity.png"
			// }

			access_control_role_type = "ADMIN_USERS_ONLY"

			access_control_group_options {
				type = "ANY_GROUP"

				groups = [
					"${pingone_group.%[3]s.id}"
				]
			}

			enabled = true

			saml_options {
				type = "WEB_APP"
				acs_urls = ["https://pingidentity.com", "https://www.pingidentity.com"]
				assertion_duration = 3600
				sp_entity_id = "sp:entity:localhost"

				assertion_signed_enabled = false
				// idp_signing_key_id = 
				nameid_format = "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress"
				response_is_signed = true
				slo_binding = "HTTP_REDIRECT"
				slo_endpoint = "https://www.pingidentity.com"
				slo_response_endpoint = "https://www.pingidentity.com"

				// sp_verification_certificate_ids = []

			}
		}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccApplicationConfig_SAMLMinimal(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s
		resource "pingone_application" "%[3]s" {
			environment_id = "${pingone_environment.%[2]s.id}"
			name = "%[4]s"
			enabled = true

			saml_options {
				acs_urls = ["https://pingidentity.com"]
				assertion_duration = 3600
				sp_entity_id = "sp:entity:localhost"
			}
		}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccApplicationConfig_Enabled(environmentName, licenseID, resourceName, name string, enabled bool) string {
	return fmt.Sprintf(`
		%[1]s
		resource "pingone_application" "%[3]s" {
			environment_id = "${pingone_environment.%[2]s.id}"
			name = "%[4]s"
			enabled = %[5]t

			oidc_options {
				type                        = "SINGLE_PAGE_APP"
				grant_types                 = ["AUTHORIZATION_CODE"]
				response_types              = ["CODE"]
				pkce_enforcement            = "S256_REQUIRED"
				token_endpoint_authn_method = "NONE"
				redirect_uris               = ["https://www.pingidentity.com"]
			}
		}
		`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name, enabled)
}

// Error conditions

// func testAccApplicationConfig_NoType(environmentName, licenseID, resourceName, name string) string {
// 	return fmt.Sprintf(`
// 		%[1]s
// 		resource "pingone_application" "%[2]s" {
// 			environment_id = "${pingone_environment.%[1]s.id}"
// 			name = "%[3]s"
// 		}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), resourceName, name)
// }
