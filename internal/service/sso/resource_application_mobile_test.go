// Copyright © 2026 Ping Identity Corporation

package sso_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	acctestlegacysdk "github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/sso"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// TestAccApplication_NativeMobile_NewFields covers the mobile fields that are
// settable on applications created via the API: apple_team_id,
// integrity_detection.apple_app_attest_settings, and push.number_matching.
//
// The PingID mobile app configuration settings (ping_id_settings,
// allow_authentication_device_locked, allow_un_pair_and_change_device_from_app,
// mobile_app_security_pin_require, mobile_device_biometrics_setting,
// notify_optional/require_mobile_version_updates, otp_push_notification_enabled,
// show_authentication_information) are only accepted by the PingOne API on the
// environment's pre-provisioned PingID Mobile application; creating an
// application with those properties set fails with a constraint violation, so
// they cannot be covered by acceptance tests on a created resource.
func TestAccApplication_NativeMobile_NewFields(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application.%s", resourceName)

	name := resourceName

	fullCheck := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "oidc_options.type", "NATIVE_APP"),
		resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.apple_team_id", "ABC123DEFG"),
		resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.enabled", "true"),
		resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.excluded_platforms.0", "GOOGLE"),
		resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.cache_duration.amount", "60"),
		resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.integrity_detection.apple_app_attest_settings.enabled", "true"),
		resource.TestCheckResourceAttr(resourceFullName, "oidc_options.mobile_app.push.number_matching.type", "ENTER_MANUALLY"),
	)

	minimalCheck := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "oidc_options.type", "NATIVE_APP"),
		resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.mobile_app.apple_team_id"),
		resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.mobile_app.push"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             acctestlegacysdk.Application_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			{
				Config: testAccApplicationConfig_OIDC_NativeMobile_NewFields_Full(resourceName),
				Check:  fullCheck,
			},
			// Minimal (fields removed)
			{
				Config: testAccApplicationConfig_OIDC_NativeMobile_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			// Back to full
			{
				Config: testAccApplicationConfig_OIDC_NativeMobile_NewFields_Full(resourceName),
				Check:  fullCheck,
			},
			// Test importing the resource
			{
				Config:       testAccApplicationConfig_OIDC_NativeMobile_NewFields_Full(resourceName),
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("resource not found: %s", resourceFullName)
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

func testAccApplicationConfig_OIDC_NativeMobile_NewFields_Full(resourceName string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[2]s"
  enabled        = true

  oidc_options = {
    type                       = "NATIVE_APP"
    grant_types                = ["CLIENT_CREDENTIALS"]
    token_endpoint_auth_method = "CLIENT_SECRET_BASIC"

    mobile_app = {
      apple_team_id = "ABC123DEFG"

      integrity_detection = {
        enabled            = true
        excluded_platforms = ["GOOGLE"]

        cache_duration = {
          amount = 60
          units  = "MINUTES"
        }

        apple_app_attest_settings = {
          enabled = true
        }
      }

      push = {
        number_matching = {
          type = "ENTER_MANUALLY"
        }
      }
    }
  }
}`, acctest.GenericSandboxEnvironment(), resourceName)
}
