// Copyright © 2025 Ping Identity Corporation

package mfa_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	acctestlegacysdk "github.com/pingidentity/terraform-provider-pingone/internal/acctest/legacysdk"
	baselegacysdk "github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/base/legacysdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/mfa"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccMFADevicePolicyDefault_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_device_policy_default.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	var environmentID string

	var p1Client *client.Client
	var ctx = context.Background()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoBeta(t)
			p1Client = acctestlegacysdk.PreCheckTestClient(ctx, t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.MFADevicePolicyDefault_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Test removal of the environment
			{
				Config: testAccMFADevicePolicyDefaultConfig_Full(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return fmt.Errorf("Resource Not found: %s", resourceFullName)
						}
						environmentID = rs.Primary.Attributes["environment_id"]
						return nil
					},
				),
			},
			{
				PreConfig: func() {
					baselegacysdk.Environment_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccMFADevicePolicyDefault_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_device_policy_default.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.MFADevicePolicyDefault_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFADevicePolicyDefaultConfig_Full(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "authentication.device_selection", "DEFAULT_TO_FIRST"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "SMS_THEN_EMAIL"),
					resource.TestCheckResourceAttr(resourceFullName, "ignore_user_lock", "true"),
					resource.TestMatchResourceAttr(resourceFullName, "notifications_policy.id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "updated_at", verify.RFC3339Regexp),
					resource.TestCheckResourceAttr(resourceFullName, "remember_me.web.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "remember_me.web.life_time.duration", "60"),
					resource.TestCheckResourceAttr(resourceFullName, "remember_me.web.life_time.time_unit", "MINUTES"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.pairing_disabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.otp.lifetime.duration", "75"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.otp.lifetime.time_unit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.otp.failure.count", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.otp.failure.cool_down.duration", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.otp.failure.cool_down.time_unit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.otp.otp_length", "7"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.prompt_for_nickname_on_pairing", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "email.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "email.pairing_disabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.applications.%", "1"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, fmt.Sprintf("pingone_application.%s", resourceName), "mobile.applications.%s.auto_enrollment.enabled", "true"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, fmt.Sprintf("pingone_application.%s", resourceName), "mobile.applications.%s.device_authorization.enabled", "true"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, fmt.Sprintf("pingone_application.%s", resourceName), "mobile.applications.%s.device_authorization.extra_verification", "permissive"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, fmt.Sprintf("pingone_application.%s", resourceName), "mobile.applications.%s.integrity_detection", "permissive"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, fmt.Sprintf("pingone_application.%s", resourceName), "mobile.applications.%s.push.enabled", "false"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, fmt.Sprintf("pingone_application.%s", resourceName), "mobile.applications.%s.push_limit.count", "5"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, fmt.Sprintf("pingone_application.%s", resourceName), "mobile.applications.%s.pairing_key_lifetime.duration", "10"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, fmt.Sprintf("pingone_application.%s", resourceName), "mobile.applications.%s.push_timeout.duration", "45"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.pairing_disabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.prompt_for_nickname_on_pairing", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.pairing_disabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.prompt_for_nickname_on_pairing", "false"),
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

						return rs.Primary.Attributes["environment_id"], nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccMFADevicePolicyDefault_Minimal(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_device_policy_default.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.MFADevicePolicyDefault_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFADevicePolicyDefaultConfig_Minimal(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "authentication.device_selection", "DEFAULT_TO_FIRST"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "NONE"),
					resource.TestCheckResourceAttr(resourceFullName, "ignore_user_lock", "false"),
					resource.TestCheckNoResourceAttr(resourceFullName, "notifications_policy.id"),
					resource.TestMatchResourceAttr(resourceFullName, "updated_at", verify.RFC3339Regexp),
					resource.TestCheckResourceAttr(resourceFullName, "remember_me.web.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "remember_me.web.life_time.duration", "30"),
					resource.TestCheckResourceAttr(resourceFullName, "remember_me.web.life_time.time_unit", "MINUTES"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.enabled", "false"),
				),
			},
		},
	})
}

func TestAccMFADevicePolicyDefault_Change(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_device_policy_default.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.MFADevicePolicyDefault_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFADevicePolicyDefaultConfig_Full(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "authentication.device_selection", "DEFAULT_TO_FIRST"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "SMS_THEN_EMAIL"),
					resource.TestCheckResourceAttr(resourceFullName, "ignore_user_lock", "true"),
					resource.TestMatchResourceAttr(resourceFullName, "notifications_policy.id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "remember_me.web.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "remember_me.web.life_time.duration", "60"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.enabled", "true"),
				),
			},
			{
				Config: testAccMFADevicePolicyDefaultConfig_Minimal(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "authentication.device_selection", "DEFAULT_TO_FIRST"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "NONE"),
					resource.TestCheckResourceAttr(resourceFullName, "ignore_user_lock", "false"),
					resource.TestCheckNoResourceAttr(resourceFullName, "notifications_policy.id"),
					resource.TestCheckResourceAttr(resourceFullName, "remember_me.web.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "remember_me.web.life_time.duration", "30"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.enabled", "false"),
				),
			},
			{
				Config: testAccMFADevicePolicyDefaultConfig_Full(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "authentication.device_selection", "DEFAULT_TO_FIRST"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "SMS_THEN_EMAIL"),
					resource.TestCheckResourceAttr(resourceFullName, "ignore_user_lock", "true"),
					resource.TestCheckResourceAttrSet(resourceFullName, "notifications_policy.id"),
					resource.TestCheckResourceAttr(resourceFullName, "remember_me.web.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.enabled", "true"),
				),
			},
		},
	})
}

func TestAccMFADevicePolicyDefault_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_device_policy_default.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.MFADevicePolicyDefault_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccMFADevicePolicyDefaultConfig_Minimal(environmentName, licenseID, resourceName, name),
			},
			// Errors
			{
				ResourceName:  resourceFullName,
				ImportStateId: "/",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`Unexpected Import Identifier`),
			},
			{
				ResourceName:  resourceFullName,
				ImportStateId: "badformat",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`Unexpected Import Identifier`),
			},
		},
	})
}

func TestAccMFADevicePolicyDefault_PingID_Full(t *testing.T) {
	// t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_device_policy_default.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckRegionSupportsWorkforce(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.MFADevicePolicyDefault_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFADevicePolicyDefaultConfig_PingID_Full(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "policy_type", "PING_ONE_ID"),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "authentication.device_selection", "PROMPT_TO_SELECT"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "SMS_THEN_EMAIL"),
					resource.TestCheckResourceAttr(resourceFullName, "ignore_user_lock", "true"),
					resource.TestMatchResourceAttr(resourceFullName, "updated_at", verify.RFC3339Regexp),
					resource.TestCheckResourceAttr(resourceFullName, "remember_me.web.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "email.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.applications.%", "1"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, fmt.Sprintf("data.pingone_application.%s", resourceName), "mobile.applications.%s.type", "pingIdAppConfig"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, fmt.Sprintf("data.pingone_application.%s", resourceName), "mobile.applications.%s.biometrics_enabled", "true"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, fmt.Sprintf("data.pingone_application.%s", resourceName), "mobile.applications.%s.integrity_detection", "permissive"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, fmt.Sprintf("data.pingone_application.%s", resourceName), "mobile.applications.%s.otp.enabled", "true"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, fmt.Sprintf("data.pingone_application.%s", resourceName), "mobile.applications.%s.pairing_disabled", "false"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, fmt.Sprintf("data.pingone_application.%s", resourceName), "mobile.applications.%s.pairing_key_lifetime.duration", "15"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, fmt.Sprintf("data.pingone_application.%s", resourceName), "mobile.applications.%s.pairing_key_lifetime.time_unit", "MINUTES"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, fmt.Sprintf("data.pingone_application.%s", resourceName), "mobile.applications.%s.push.enabled", "true"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, fmt.Sprintf("data.pingone_application.%s", resourceName), "mobile.applications.%s.push.number_matching.enabled", "true"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, fmt.Sprintf("data.pingone_application.%s", resourceName), "mobile.applications.%s.push_limit.count", "10"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, fmt.Sprintf("data.pingone_application.%s", resourceName), "mobile.applications.%s.push_limit.lock_duration.duration", "45"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, fmt.Sprintf("data.pingone_application.%s", resourceName), "mobile.applications.%s.push_limit.lock_duration.time_unit", "MINUTES"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, fmt.Sprintf("data.pingone_application.%s", resourceName), "mobile.applications.%s.push_limit.time_period.duration", "15"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, fmt.Sprintf("data.pingone_application.%s", resourceName), "mobile.applications.%s.push_limit.time_period.time_unit", "MINUTES"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, fmt.Sprintf("data.pingone_application.%s", resourceName), "mobile.applications.%s.new_request_duration_configuration.device_timeout.duration", "30"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, fmt.Sprintf("data.pingone_application.%s", resourceName), "mobile.applications.%s.new_request_duration_configuration.device_timeout.time_unit", "SECONDS"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, fmt.Sprintf("data.pingone_application.%s", resourceName), "mobile.applications.%s.new_request_duration_configuration.total_timeout.duration", "60"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, fmt.Sprintf("data.pingone_application.%s", resourceName), "mobile.applications.%s.new_request_duration_configuration.total_timeout.time_unit", "SECONDS"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, fmt.Sprintf("data.pingone_application.%s", resourceName), "mobile.applications.%s.ip_pairing_configuration.any_ip_address", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "desktop.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "desktop.otp.failure.count", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "yubikey.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oath_token.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oath_token.otp.failure.count", "3"),
				),
			},
		},
	})
}

func TestAccMFADevicePolicyDefault_PingID_Minimal(t *testing.T) {
	// t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_device_policy_default.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckRegionSupportsWorkforce(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.MFADevicePolicyDefault_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFADevicePolicyDefaultConfig_PingID_Minimal(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "policy_type", "PING_ONE_ID"),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "desktop.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "yubikey.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oath_token.enabled", "false"),
				),
			},
		},
	})
}

func TestAccMFADevicePolicyDefault_PingID_Change(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_device_policy_default.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
			acctest.PreCheckRegionSupportsWorkforce(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.MFADevicePolicyDefault_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFADevicePolicyDefaultConfig_PingID_Minimal_WithNotificationPolicy(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "policy_type", "PING_ONE_ID"),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "desktop.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "yubikey.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oath_token.enabled", "false"),
				),
			},
			{
				Config: testAccMFADevicePolicyDefaultConfig_PingID_Full(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "policy_type", "PING_ONE_ID"),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.enabled", "true"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, fmt.Sprintf("data.pingone_application.%s", resourceName), "mobile.applications.%s.type", "pingIdAppConfig"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "email.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "desktop.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "yubikey.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oath_token.enabled", "true"),
				),
			},
		},
	})
}

func TestAccMFADevicePolicyDefault_Validation(t *testing.T) {
	t.Parallel()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	testCases := map[string]func(t *testing.T){
		"General_Validation": func(t *testing.T) {
			t.Parallel()

			resourceName := acctest.ResourceNameGen()
			environmentName := acctest.ResourceNameGenEnvironment()
			name := resourceName

			resource.Test(t, resource.TestCase{
				PreCheck: func() {
					acctest.PreCheckClient(t)
					acctest.PreCheckNewEnvironment(t)
				},
				ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
				CheckDestroy:             mfa.MFADevicePolicyDefault_CheckDestroy,
				ErrorCheck:               acctest.ErrorCheck(t),
				Steps: []resource.TestStep{
					// Notifications Policy - Invalid ID format
					{
						Config:      testAccMFADevicePolicyDefaultConfig_InvalidNotificationsPolicyID(environmentName, licenseID, resourceName, name),
						ExpectError: regexp.MustCompile(`The PingOne resource ID is malformed`),
					},
					// Authentication - Invalid device_selection
					{
						Config:      testAccMFADevicePolicyDefaultConfig_Authentication(resourceName, name, "INVALID_VALUE"),
						ExpectError: regexp.MustCompile(`Attribute authentication.device_selection value must be one of:`),
					},
					// New Device Notification - Invalid value
					{
						Config:      testAccMFADevicePolicyDefaultConfig_NewDeviceNotification(resourceName, name, "INVALID_VALUE"),
						ExpectError: regexp.MustCompile(`Attribute new_device_notification value must be one of:`),
					},
					// Remember Me - Duration out of range for MINUTES
					{
						Config:      testAccMFADevicePolicyDefaultConfig_RememberMe_MinutesHigh(resourceName, name),
						ExpectError: regexp.MustCompile(`Attribute remember_me.web.life_time.duration value must be between 1 and[\s\n]+129600`),
					},
					// Remember Me - Duration out of range for HOURS
					{
						Config:      testAccMFADevicePolicyDefaultConfig_RememberMe_HoursHigh(resourceName, name),
						ExpectError: regexp.MustCompile(`Attribute remember_me.web.life_time.duration value must be between 1 and[\s\n]+2160`),
					},
					// Remember Me - Duration out of range for DAYS
					{
						Config:      testAccMFADevicePolicyDefaultConfig_RememberMe_DaysHigh(resourceName, name),
						ExpectError: regexp.MustCompile(`Attribute remember_me.web.life_time.duration value must be between 1 and[\s\n]+90`),
					},
				},
			})
		},
		"PingOneMFA_Mobile_Validation": func(t *testing.T) {
			t.Parallel()

			resourceName := acctest.ResourceNameGen()
			environmentName := acctest.ResourceNameGenEnvironment()
			name := resourceName

			resource.Test(t, resource.TestCase{
				PreCheck: func() {
					acctest.PreCheckClient(t)
					acctest.PreCheckNewEnvironment(t)
				},
				ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
				CheckDestroy:             mfa.MFADevicePolicyDefault_CheckDestroy,
				ErrorCheck:               acctest.ErrorCheck(t),
				Steps: []resource.TestStep{
					// Missing auto_enrollment
					{
						Config:      testAccMFADevicePolicyDefaultConfig_PingOneMFA_MobileApp_MissingAutoEnrollment(environmentName, licenseID, resourceName, name),
						ExpectError: regexp.MustCompile(`The argument\s+mobile.applications\[.+\].auto_enrollment\s+is required because\s+policy_type is configured as:\s+"PING_ONE_MFA"`),
					},
					// Missing device_authorization
					{
						Config:      testAccMFADevicePolicyDefaultConfig_PingOneMFA_MobileApp_MissingDeviceAuthorization(environmentName, licenseID, resourceName, name),
						ExpectError: regexp.MustCompile(`The argument\s+mobile.applications\[.+\].device_authorization\s+is required because\s+policy_type is configured as:\s+"PING_ONE_MFA"`),
					},
					// Missing integrity_detection
					{
						Config:      testAccMFADevicePolicyDefaultConfig_PingOneMFA_MobileApp_MissingIntegrityDetection(environmentName, licenseID, resourceName, name),
						ExpectError: regexp.MustCompile(`The argument\s+mobile.applications\[.+\].integrity_detection\s+is required because\s+policy_type is configured as:\s+"PING_ONE_MFA"`),
					},
					// Invalid integrity_detection
					{
						Config:      testAccMFADevicePolicyDefaultConfig_MobileIntegrityDetection(environmentName, licenseID, resourceName, name, "INVALID_VALUE"),
						ExpectError: regexp.MustCompile(`Attribute\s+mobile.applications\[.+\].integrity_detection\s+value must be one of:`),
					},
					// Biometrics enabled conflict
					{
						Config:      testAccMFADevicePolicyDefaultConfig_PingOneMFA_MobileApp_BiometricsEnabled(environmentName, licenseID, resourceName, name),
						ExpectError: regexp.MustCompile(`The argument cannot be defined if the value "PING_ONE_MFA" is present`),
					},
					// New request duration configuration conflict
					{
						Config:      testAccMFADevicePolicyDefaultConfig_PingOneMFA_MobileApp_NewRequestDuration(environmentName, licenseID, resourceName, name),
						ExpectError: regexp.MustCompile(`The argument cannot be defined if the value "PING_ONE_MFA" is present`),
					},
					// IP pairing configuration conflict
					{
						Config:      testAccMFADevicePolicyDefaultConfig_PingOneMFA_MobileApp_IPPairing(environmentName, licenseID, resourceName, name),
						ExpectError: regexp.MustCompile(`The argument cannot be defined if the value "PING_ONE_MFA" is present`),
					},
				},
			})
		},
		"PingID_Structure_Validation": func(t *testing.T) {
			t.Parallel()

			resourceName := acctest.ResourceNameGen()
			name := resourceName

			resource.Test(t, resource.TestCase{
				PreCheck: func() {
					acctest.PreCheckClient(t)
					acctest.PreCheckNewEnvironment(t)
					acctest.PreCheckRegionSupportsWorkforce(t)
				},
				ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
				CheckDestroy:             mfa.MFADevicePolicyDefault_CheckDestroy,
				ErrorCheck:               acctest.ErrorCheck(t),
				Steps: []resource.TestStep{
					// Desktop should conflict with PingOneMFA policy type
					{
						Config:      testAccMFADevicePolicyDefaultConfig_PingID_DesktopWithPingOneMFA(resourceName, name),
						ExpectError: regexp.MustCompile(`Invalid argument combination`),
					},
					// Yubikey should conflict with PingOneMFA policy type
					{
						Config:      testAccMFADevicePolicyDefaultConfig_PingID_YubikeyWithPingOneMFA(resourceName, name),
						ExpectError: regexp.MustCompile(`Invalid argument combination`),
					},
					// Missing desktop block
					{
						Config:      testAccMFADevicePolicyDefaultConfig_PingID_MissingDesktop(resourceName, name),
						ExpectError: regexp.MustCompile(`The argument\s+desktop\s+is required because\s+policy_type is configured as:\s+"PING_ONE_ID"`),
					},
					// Missing yubikey block
					{
						Config:      testAccMFADevicePolicyDefaultConfig_PingID_MissingYubikey(resourceName, name),
						ExpectError: regexp.MustCompile(`The argument\s+yubikey\s+is required because\s+policy_type is configured as:\s+"PING_ONE_ID"`),
					},
					// Mobile must be enabled for PingID policies
					{
						Config:      testAccMFADevicePolicyDefaultConfig_PingID_MobileDisabled(resourceName, name),
						ExpectError: regexp.MustCompile(`Attribute mobile.enabled must be true when attribute policy_type value is`),
					},
					// Auto enrollment conflict
					{
						Config:      testAccMFADevicePolicyDefaultConfig_PingID_MobileApp_AutoEnrollment(resourceName, name),
						ExpectError: regexp.MustCompile(`The argument cannot be defined if the value "PING_ONE_ID" is present`),
					},
					// Device authorization conflict
					{
						Config:      testAccMFADevicePolicyDefaultConfig_PingID_MobileApp_DeviceAuthorization(resourceName, name),
						ExpectError: regexp.MustCompile(`The argument cannot be defined if the value "PING_ONE_ID" is present`),
					},
					// Missing new_request_duration_configuration
					{
						Config:      testAccMFADevicePolicyDefaultConfig_PingID_MobileApp_MissingNewRequestDuration(resourceName, name),
						ExpectError: regexp.MustCompile(`The argument\s+mobile.applications\[.+\].new_request_duration_configuration\s+is\s+required because\s+policy_type is configured as:\s+"PING_ONE_ID"`),
					},
					// Missing ip_pairing_configuration
					{
						Config:      testAccMFADevicePolicyDefaultConfig_PingID_MobileApp_MissingIPPairing(resourceName, name),
						ExpectError: regexp.MustCompile(`The argument\s+mobile.applications\[.+\].ip_pairing_configuration\s+is\s+required\s+because\s+policy_type is configured as:\s+"PING_ONE_ID"`),
					},
				},
			})
		},
		"Common_Field_Validation": func(t *testing.T) {
			t.Parallel()

			resourceName := acctest.ResourceNameGen()
			environmentName := acctest.ResourceNameGenEnvironment()
			name := resourceName

			resource.Test(t, resource.TestCase{
				PreCheck: func() {
					acctest.PreCheckClient(t)
					acctest.PreCheckNewEnvironment(t)
					acctest.PreCheckRegionSupportsWorkforce(t)
				},
				ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
				CheckDestroy:             mfa.MFADevicePolicyDefault_CheckDestroy,
				ErrorCheck:               acctest.ErrorCheck(t),
				Steps: []resource.TestStep{
					// Email - OTP failure count too high
					{
						Config:      testAccMFADevicePolicyDefaultConfig_EmailOtpFailureCount(resourceName, name, 8),
						ExpectError: regexp.MustCompile(`Attribute email.otp.failure.count value must be between 1 and 7`),
					},
					// Email - OTP failure cool down duration too high
					{
						Config:      testAccMFADevicePolicyDefaultConfig_EmailOtpFailureCoolDownDuration(resourceName, name, 31),
						ExpectError: regexp.MustCompile(`Attribute email.otp.failure.cool_down.duration value must be between 0 and[\s\n]+30`),
					},
					// Email - OTP lifetime duration too high
					{
						Config:      testAccMFADevicePolicyDefaultConfig_EmailOtpLifetimeDuration(resourceName, name, 121),
						ExpectError: regexp.MustCompile(`Attribute email.otp.lifetime.duration value must be between 1 and 120`),
					},
					// SMS - OTP failure count too high
					{
						Config:      testAccMFADevicePolicyDefaultConfig_SmsOtpFailureCount(resourceName, name, 8),
						ExpectError: regexp.MustCompile(`Attribute sms.otp.failure.count value must be between 1 and 7`),
					},
					// SMS - OTP failure cool down duration too high
					{
						Config:      testAccMFADevicePolicyDefaultConfig_SmsOtpFailureCoolDownDuration(resourceName, name, 31),
						ExpectError: regexp.MustCompile(`Attribute sms.otp.failure.cool_down.duration value must be between 0 and[\s\n]+30`),
					},
					// SMS - OTP lifetime duration too high
					{
						Config:      testAccMFADevicePolicyDefaultConfig_SmsOtpLifetimeDuration(resourceName, name, 121),
						ExpectError: regexp.MustCompile(`Attribute sms.otp.lifetime.duration value must be between 1 and 120`),
					},
					// Voice - OTP failure count too high
					{
						Config:      testAccMFADevicePolicyDefaultConfig_VoiceOtpFailureCount(resourceName, name, 8),
						ExpectError: regexp.MustCompile(`Attribute voice.otp.failure.count value must be between 1 and 7`),
					},
					// Voice - OTP failure cool down duration too high
					{
						Config:      testAccMFADevicePolicyDefaultConfig_VoiceOtpFailureCoolDownDuration(resourceName, name, 31),
						ExpectError: regexp.MustCompile(`Attribute voice.otp.failure.cool_down.duration value must be between 0 and[\s\n]+30`),
					},
					// Voice - OTP lifetime duration too high
					{
						Config:      testAccMFADevicePolicyDefaultConfig_VoiceOtpLifetimeDuration(resourceName, name, 121),
						ExpectError: regexp.MustCompile(`Attribute voice.otp.lifetime.duration value must be between 1 and 120`),
					},
					// Mobile - Push limit lock duration too high
					{
						Config:      testAccMFADevicePolicyDefaultConfig_MobilePushLimitLockDuration(resourceName, name, 121),
						ExpectError: regexp.MustCompile(`Attribute\s+mobile.applications\[.+\].push_limit.lock_duration.duration\s+value must[\s\n]+be between[\s\n]+1 and[\s\n]+120`),
					},
					// Mobile - Push limit time period too high
					{
						Config:      testAccMFADevicePolicyDefaultConfig_MobilePushLimitTimePeriod(resourceName, name, 121),
						ExpectError: regexp.MustCompile(`Attribute\s+mobile.applications\[.+\].push_limit.time_period.duration\s+value must[\s\n]+be between[\s\n]+1 and[\s\n]+120`),
					},
					// Mobile - Push timeout duration too high
					{
						Config:      testAccMFADevicePolicyDefaultConfig_MobilePushTimeoutDuration(environmentName, licenseID, resourceName, name, 121),
						ExpectError: regexp.MustCompile(`Attribute\s+mobile.applications\[.+\].push_timeout.duration\s+value must[\s\n]+be between[\s\n]+1 and[\s\n]+120`),
					},
					// Mobile - OTP failure cool down duration too high
					{
						Config:      testAccMFADevicePolicyDefaultConfig_MobileOtpFailureCoolDownDuration(resourceName, name, 31),
						ExpectError: regexp.MustCompile(`Attribute mobile.otp.failure.cool_down.duration value must be between 2 and[\s\n]+30`),
					},
					// TOTP - OTP failure cool down duration too high
					{
						Config:      testAccMFADevicePolicyDefaultConfig_TotpOtpFailureCoolDownDuration(resourceName, name, 31),
						ExpectError: regexp.MustCompile(`Attribute totp.otp.failure.cool_down.duration value must be between 1 and[\s\n]+30`),
					},
				},
			})
		},
		"PingID_Field_Validation": func(t *testing.T) {
			t.Parallel()

			resourceName := acctest.ResourceNameGen()
			name := resourceName

			resource.Test(t, resource.TestCase{
				PreCheck: func() {
					acctest.PreCheckClient(t)
					acctest.PreCheckNewEnvironment(t)
					acctest.PreCheckRegionSupportsWorkforce(t)
				},
				ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
				CheckDestroy:             mfa.MFADevicePolicyDefault_CheckDestroy,
				ErrorCheck:               acctest.ErrorCheck(t),
				Steps: []resource.TestStep{
					// IP Pairing - Invalid CIDR
					{
						Config:      testAccMFADevicePolicyDefaultConfig_PingID_IPPairing_InvalidCIDR(resourceName, name),
						ExpectError: regexp.MustCompile(`Expected value to be in CIDR notation`),
					},
					// IP Pairing - Missing only_these_ip_addresses when any_ip_address is false
					{
						Config:      testAccMFADevicePolicyDefaultConfig_PingID_IPPairing_MissingIPs(resourceName, name),
						ExpectError: regexp.MustCompile(`The argument\s+mobile.applications\[.+\].ip_pairing_configuration.only_these_ip_addresses\s+is\s+required\s+because\s+mobile.applications\[.+\].ip_pairing_configuration.any_ip_address\s+is\s+configured\s+as: false`),
					},
					// Desktop - OTP failure count too high
					{
						Config:      testAccMFADevicePolicyDefaultConfig_PingID_Desktop_OTPCountHigh(resourceName, name),
						ExpectError: regexp.MustCompile(`Attribute desktop.otp.failure.count value must be between 1 and 7`),
					},
					// Desktop - Pairing key lifetime too long (HOURS)
					{
						Config:      testAccMFADevicePolicyDefaultConfig_PingID_Desktop_PairingKeyLifetimeHigh(resourceName, name),
						ExpectError: regexp.MustCompile(`Attribute desktop.pairing_key_lifetime.duration value must be between 1 and\s+48`),
					},
					// Mobile - Push limit count out of range
					{
						Config:      testAccMFADevicePolicyDefaultConfig_MobilePushLimit_CountHigh(resourceName, name),
						ExpectError: regexp.MustCompile(`Attribute\s+mobile.applications\[.+\].push_limit.count\s+value must be between 1 and[\s\n]+50`),
					},
					// Mobile - Device timeout duration out of range
					{
						Config:      testAccMFADevicePolicyDefaultConfig_MobileNewRequestDuration_DeviceTimeoutHigh(resourceName, name),
						ExpectError: regexp.MustCompile(`Attribute[\s\n]+mobile.applications\[.+\].new_request_duration_configuration.device_timeout.duration[\s\n]+value must be between 15 and 75`),
					},
					// Mobile - Total timeout duration out of range
					{
						Config:      testAccMFADevicePolicyDefaultConfig_MobileNewRequestDuration_TotalTimeoutHigh(resourceName, name),
						ExpectError: regexp.MustCompile(`Attribute[\s\n]+mobile.applications\[.+\].new_request_duration_configuration.total_timeout.duration[\s\n]+value must be between 30 and 90`),
					},
					// Mobile - Device timeout duration too low
					{
						Config:      testAccMFADevicePolicyDefaultConfig_MobileNewRequestDuration_DeviceTimeout(resourceName, name, 14),
						ExpectError: regexp.MustCompile(`Attribute[\s\n]+mobile.applications\[.+\].new_request_duration_configuration.device_timeout.duration[\s\n]+value must be between 15 and 75`),
					},
					// Mobile - Total timeout duration too low
					{
						Config:      testAccMFADevicePolicyDefaultConfig_MobileNewRequestDuration_TotalTimeout(resourceName, name, 29),
						ExpectError: regexp.MustCompile(`Attribute[\s\n]+mobile.applications\[.+\].new_request_duration_configuration.total_timeout.duration[\s\n]+value must be between 30 and 90`),
					},
					// Mobile - OTP failure count out of range (0)
					{
						Config:      testAccMFADevicePolicyDefaultConfig_MobileOtpFailureCount(resourceName, name, 0),
						ExpectError: regexp.MustCompile(`Attribute mobile.otp.failure.count value must be between 1 and 7`),
					},
					// Mobile - OTP failure count out of range (8)
					{
						Config:      testAccMFADevicePolicyDefaultConfig_MobileOtpFailureCount(resourceName, name, 8),
						ExpectError: regexp.MustCompile(`Attribute mobile.otp.failure.count value must be between 1 and 7`),
					},
				},
			})
		},
	}

	for name, testFunc := range testCases {
		t.Run(name, testFunc)
	}
}

func TestAccMFADevicePolicyDefault_BOMValidation(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")
	region := os.Getenv("PINGONE_REGION_CODE")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.MFADevicePolicyDefault_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccMFADevicePolicyDefaultConfig_BOMValidation(environmentName, licenseID, region, resourceName, name),
				ExpectError: regexp.MustCompile("Unsupported Policy Type"),
			},
		},
	})
}

func TestAccMFADevicePolicyDefault_BOMValidation_CrossTypes(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	environmentNamePingID := acctest.ResourceNameGenEnvironment() + "-pingid"
	environmentNameMFA := acctest.ResourceNameGenEnvironment() + "-mfa"

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")
	region := os.Getenv("PINGONE_REGION_CODE")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.MFADevicePolicyDefault_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// 1. Environment with PingID (no MFA) -> Try to create PING_ONE_MFA policy
			{
				Config:      testAccMFADevicePolicyDefaultConfig_BOMValidation_PingIDEnv(environmentNamePingID, licenseID, region, resourceName, name),
				ExpectError: regexp.MustCompile("Unsupported Policy Type"),
			},
			// 2. Environment with MFA (no PingID) -> Try to create PING_ONE_ID policy
			{
				Config:      testAccMFADevicePolicyDefaultConfig_BOMValidation_MFAEnv(environmentNameMFA, licenseID, region, resourceName, name),
				ExpectError: regexp.MustCompile("Unsupported Policy Type"),
			},
		},
	})
}

func testAccMFADevicePolicyDefaultConfig_Full(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s"
  enabled        = true

  oidc_options = {
    type                       = "NATIVE_APP"
    grant_types                = ["AUTHORIZATION_CODE"]
    response_types             = ["CODE"]
    token_endpoint_auth_method = "NONE"
    redirect_uris              = ["https://example.com"]
  }
}

resource "pingone_mfa_device_policy_default" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  policy_type    = "PING_ONE_MFA"

  name = "%[4]s"

  authentication = {
    device_selection = "DEFAULT_TO_FIRST"
  }

  new_device_notification = "SMS_THEN_EMAIL"
  ignore_user_lock        = true

  notifications_policy = {
    id = pingone_notification_policy.%[3]s.id
  }

  remember_me = {
    web = {
      enabled = true
      life_time = {
        duration  = 60
        time_unit = "MINUTES"
      }
    }
  }

  sms = {
    enabled                        = true
    pairing_disabled               = true
    prompt_for_nickname_on_pairing = true
    otp = {
      failure = {
        count = 5
        cool_down = {
          duration  = 5
          time_unit = "SECONDS"
        }
      }
      lifetime = {
        duration  = 75
        time_unit = "SECONDS"
      }
      otp_length = 7
    }
  }

  email = {
    enabled                        = true
    pairing_disabled               = false
    prompt_for_nickname_on_pairing = false
    otp = {
      failure = {
        count = 3
        cool_down = {
          duration  = 2
          time_unit = "MINUTES"
        }
      }
      lifetime = {
        duration  = 30
        time_unit = "MINUTES"
      }
      otp_length = 6
    }
  }

  voice = {
    enabled                        = true
    pairing_disabled               = false
    prompt_for_nickname_on_pairing = false
    otp = {
      failure = {
        count = 3
        cool_down = {
          duration  = 2
          time_unit = "MINUTES"
        }
      }
      lifetime = {
        duration  = 30
        time_unit = "MINUTES"
      }
      otp_length = 6
    }
  }

  mobile = {
    enabled                        = true
    prompt_for_nickname_on_pairing = false
    otp = {
      failure = {
        count = 3
        cool_down = {
          duration  = 2
          time_unit = "MINUTES"
        }
      }
    }

    applications = {
      (pingone_application.%[3]s.id) = {
        auto_enrollment = {
          enabled = true
        }

        device_authorization = {
          enabled            = true
          extra_verification = "permissive"
        }

        integrity_detection = "permissive"

        otp = {
          enabled = true
        }

        pairing_disabled = false

        pairing_key_lifetime = {
          duration  = 10
          time_unit = "MINUTES"
        }

        push = {
          enabled = false
        }

        push_limit = {
          count = 5
          lock_duration = {
            duration  = 30
            time_unit = "MINUTES"
          }
          time_period = {
            duration  = 10
            time_unit = "MINUTES"
          }
        }

        push_timeout = {
          duration  = 45
          time_unit = "SECONDS"
        }
      }
    }
  }

  totp = {
    enabled                        = true
    pairing_disabled               = false
    prompt_for_nickname_on_pairing = false
    otp = {
      failure = {
        count = 3
        cool_down = {
          duration  = 2
          time_unit = "MINUTES"
        }
      }
    }
  }

  fido2 = {
    enabled                        = true
    pairing_disabled               = false
    prompt_for_nickname_on_pairing = false
  }
}

resource "pingone_notification_policy" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s"

  lifecycle {
    create_before_destroy = true
  }
}
`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccMFADevicePolicyDefaultConfig_Minimal(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_mfa_device_policy_default" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  policy_type    = "PING_ONE_MFA"

  name = "%[4]s"

  sms = {
    enabled = false
  }

  voice = {
    enabled = false
  }

  email = {
    enabled = false
  }

  mobile = {
    enabled = false
  }

  totp = {
    enabled = false
  }

  fido2 = {
    enabled = false
  }
}`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccMFADevicePolicyDefaultConfig_PingID_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_mfa_device_policy_default" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  policy_type    = "PING_ONE_ID"

  name = "%[3]s"

  sms = {
    enabled = false
  }

  voice = {
    enabled = false
  }

  email = {
    enabled = false
  }

  mobile = {
    enabled = true
  }

  totp = {
    enabled = false
  }

  desktop = {
    enabled = false
  }

  yubikey = {
    enabled = false
  }

  oath_token = {
    enabled = false
  }
}`, acctest.WorkforceV2SandboxEnvironment(), resourceName, name)
}

func testAccMFADevicePolicyDefaultConfig_PingID_Minimal_WithNotificationPolicy(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_notification_policy" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  name           = "%[3]s"
}

resource "pingone_mfa_device_policy_default" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  policy_type    = "PING_ONE_ID"

  name = "%[3]s"

  sms = {
    enabled = false
  }

  voice = {
    enabled = false
  }

  email = {
    enabled = false
  }

  mobile = {
    enabled = true
  }

  totp = {
    enabled = false
  }

  desktop = {
    enabled = false
  }

  yubikey = {
    enabled = false
  }

  oath_token = {
    enabled = false
  }
}`, acctest.WorkforceV2SandboxEnvironment(), resourceName, name)
}

func testAccMFADevicePolicyDefaultConfig_PingID_DesktopWithPingOneMFA(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_mfa_device_policy_default" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  policy_type    = "PING_ONE_MFA"

  name = "%[3]s"

  sms = {
    enabled = false
  }

  voice = {
    enabled = false
  }

  email = {
    enabled = false
  }

  mobile = {
    enabled = false
  }

  totp = {
    enabled = false
  }

  fido2 = {
    enabled = false
  }

  desktop = {
    enabled = true
  }
}`, acctest.WorkforceV2SandboxEnvironment(), resourceName, name)
}

func testAccMFADevicePolicyDefaultConfig_PingID_YubikeyWithPingOneMFA(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_mfa_device_policy_default" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  policy_type    = "PING_ONE_MFA"

  name = "%[3]s"

  sms = {
    enabled = false
  }

  voice = {
    enabled = false
  }

  email = {
    enabled = false
  }

  mobile = {
    enabled = false
  }

  totp = {
    enabled = false
  }

  fido2 = {
    enabled = false
  }

  yubikey = {
    enabled = true
  }
}`, acctest.WorkforceV2SandboxEnvironment(), resourceName, name)
}

func testAccMFADevicePolicyDefaultConfig_PingID_MobileDisabled(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_mfa_device_policy_default" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  policy_type    = "PING_ONE_ID"

  name = "%[3]s"

  sms = {
    enabled = false
  }

  voice = {
    enabled = false
  }

  email = {
    enabled = false
  }

  mobile = {
    enabled = false
  }

  totp = {
    enabled = false
  }
}`, acctest.WorkforceV2SandboxEnvironment(), resourceName, name)
}

func testAccMFADevicePolicyDefaultConfig_InvalidNotificationsPolicyID(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_mfa_device_policy_default" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  policy_type    = "PING_ONE_MFA"

  name = "%[4]s"

  notifications_policy = {
    id = "invalid-uuid-format"
  }

  sms = {
    enabled = false
  }

  voice = {
    enabled = false
  }

  email = {
    enabled = false
  }

  mobile = {
    enabled = false
  }

  totp = {
    enabled = false
  }

  fido2 = {
    enabled = false
  }
}`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccMFADevicePolicyDefaultConfig_PingID_Full(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_notification_policy" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  name           = "%[3]s"
}

data "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  name           = "PingID Mobile"
}

resource "pingone_mfa_device_policy_default" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  policy_type    = "PING_ONE_ID"

  name = "%[3]s"

  authentication = {
    device_selection = "PROMPT_TO_SELECT"
  }

  new_device_notification = "SMS_THEN_EMAIL"
  ignore_user_lock        = true

  notifications_policy = {
    id = pingone_notification_policy.%[2]s.id
  }

  remember_me = {
    web = {
      enabled = true
      life_time = {
        duration  = 60
        time_unit = "MINUTES"
      }
    }
  }

  sms = {
    enabled                        = true
    pairing_disabled               = false
    prompt_for_nickname_on_pairing = true
    otp = {
      failure = {
        count = 5
        cool_down = {
          duration  = 5
          time_unit = "SECONDS"
        }
      }
      lifetime = {
        duration  = 75
        time_unit = "SECONDS"
      }
      otp_length = 7
    }
  }

  voice = {
    enabled                        = true
    pairing_disabled               = false
    prompt_for_nickname_on_pairing = false
    otp = {
      failure = {
        count = 4
        cool_down = {
          duration  = 3
          time_unit = "MINUTES"
        }
      }
      lifetime = {
        duration  = 30
        time_unit = "MINUTES"
      }
      otp_length = 8
    }
  }

  email = {
    enabled                        = true
    pairing_disabled               = false
    prompt_for_nickname_on_pairing = false
    otp = {
      failure = {
        count = 3
        cool_down = {
          duration  = 2
          time_unit = "MINUTES"
        }
      }
      lifetime = {
        duration  = 30
        time_unit = "MINUTES"
      }
      otp_length = 6
    }
  }

  mobile = {
    enabled                        = true
    prompt_for_nickname_on_pairing = false

    applications = {
      (data.pingone_application.%[2]s.id) = {
        type = "pingIdAppConfig"

        biometrics_enabled = true

        integrity_detection = "permissive"

        otp = {
          enabled = true
        }

        pairing_disabled = false

        pairing_key_lifetime = {
          duration  = 15
          time_unit = "MINUTES"
        }

        push = {
          enabled = true
          number_matching = {
            enabled = true
          }
        }

        push_limit = {
          count = 10
          lock_duration = {
            duration  = 45
            time_unit = "MINUTES"
          }
          time_period = {
            duration  = 15
            time_unit = "MINUTES"
          }
        }

        new_request_duration_configuration = {
          device_timeout = {
            duration  = 30
            time_unit = "SECONDS"
          }
          total_timeout = {
            duration  = 60
            time_unit = "SECONDS"
          }
        }

        ip_pairing_configuration = {
          any_ip_address = true
        }
      }
    }

    otp = {
      failure = {
        count = 3
        cool_down = {
          duration  = 2
          time_unit = "MINUTES"
        }
      }
    }
  }

  totp = {
    enabled                        = true
    pairing_disabled               = false
    prompt_for_nickname_on_pairing = false
    otp = {
      failure = {
        count = 3
        cool_down = {
          duration  = 2
          time_unit = "MINUTES"
        }
      }
    }
  }

  desktop = {
    enabled = true
    otp = {
      failure = {
        count = 5
        cool_down = {
          duration  = 3
          time_unit = "MINUTES"
        }
      }
    }
    pairing_key_lifetime = {
      duration  = 48
      time_unit = "HOURS"
    }
  }

  yubikey = {
    enabled = true
    otp = {
      failure = {
        count = 3
        cool_down = {
          duration  = 2
          time_unit = "MINUTES"
        }
      }
    }
    pairing_disabled = false
  }

  oath_token = {
    enabled                        = true
    pairing_disabled               = false
    prompt_for_nickname_on_pairing = false
    otp = {
      failure = {
        count = 3
        cool_down = {
          duration  = 2
          time_unit = "MINUTES"
        }
      }
    }
  }
}
`, acctest.WorkforceV2SandboxEnvironment(), resourceName, name)
}

func testAccMFADevicePolicyDefaultConfig_MobileIntegrityDetection(environmentName, licenseID, resourceName, name, integrityDetection string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_mfa_device_policy_default" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  policy_type    = "PING_ONE_MFA"
  name           = "%[4]s"

  mobile = {
    enabled = true
    applications = {
      "11111111-1111-1111-1111-111111111111" = {
        otp                  = { enabled = true }
        auto_enrollment      = { enabled = true }
        device_authorization = { enabled = true }
        integrity_detection  = "%[5]s"
      }
    }
  }
  sms   = { enabled = false }
  voice = { enabled = false }
  email = { enabled = false }
  totp  = { enabled = false }
}`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name, integrityDetection)
}

func testAccMFADevicePolicyDefaultConfig_PingID_MissingDesktop(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_mfa_device_policy_default" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  policy_type    = "PING_ONE_ID"
  name           = "%[3]s"

  mobile  = { enabled = true }
  yubikey = { enabled = false }
  sms     = { enabled = false }
  voice   = { enabled = false }
  email   = { enabled = false }
  totp    = { enabled = false }
}`, acctest.WorkforceV2SandboxEnvironment(), resourceName, name)
}

func testAccMFADevicePolicyDefaultConfig_PingID_MissingYubikey(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_mfa_device_policy_default" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  policy_type    = "PING_ONE_ID"
  name           = "%[3]s"

  mobile  = { enabled = true }
  desktop = { enabled = false }
  sms     = { enabled = false }
  voice   = { enabled = false }
  email   = { enabled = false }
  totp    = { enabled = false }
}`, acctest.WorkforceV2SandboxEnvironment(), resourceName, name)
}

func testAccMFADevicePolicyDefaultConfig_PingID_MobileApp_AutoEnrollment(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_mfa_device_policy_default" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  policy_type    = "PING_ONE_ID"
  name           = "%[3]s"

  mobile = {
    enabled = true
    applications = {
      "11111111-1111-1111-1111-111111111111" = {
        otp             = { enabled = true }
        auto_enrollment = { enabled = true }
        new_request_duration_configuration = {
          device_timeout = { duration = 25 }
          total_timeout  = { duration = 40 }
        }
        ip_pairing_configuration = { any_ip_address = true }
      }
    }
  }
  desktop = { enabled = false }
  yubikey = { enabled = false }
  sms     = { enabled = false }
  voice   = { enabled = false }
  email   = { enabled = false }
  totp    = { enabled = false }
}`, acctest.WorkforceV2SandboxEnvironment(), resourceName, name)
}

func testAccMFADevicePolicyDefaultConfig_PingID_MobileApp_DeviceAuthorization(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_mfa_device_policy_default" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  policy_type    = "PING_ONE_ID"
  name           = "%[3]s"

  mobile = {
    enabled = true
    applications = {
      "11111111-1111-1111-1111-111111111111" = {
        otp                  = { enabled = true }
        device_authorization = { enabled = true }
        new_request_duration_configuration = {
          device_timeout = { duration = 25 }
          total_timeout  = { duration = 40 }
        }
        ip_pairing_configuration = { any_ip_address = true }
      }
    }
  }
  desktop = { enabled = false }
  yubikey = { enabled = false }
  sms     = { enabled = false }
  voice   = { enabled = false }
  email   = { enabled = false }
  totp    = { enabled = false }
}`, acctest.WorkforceV2SandboxEnvironment(), resourceName, name)
}

func testAccMFADevicePolicyDefaultConfig_PingID_MobileApp_MissingNewRequestDuration(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_mfa_device_policy_default" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  policy_type    = "PING_ONE_ID"
  name           = "%[3]s"

  mobile = {
    enabled = true
    applications = {
      "11111111-1111-1111-1111-111111111111" = {
        otp                      = { enabled = true }
        ip_pairing_configuration = { any_ip_address = true }
      }
    }
  }
  desktop = { enabled = false }
  yubikey = { enabled = false }
  sms     = { enabled = false }
  voice   = { enabled = false }
  email   = { enabled = false }
  totp    = { enabled = false }
}`, acctest.WorkforceV2SandboxEnvironment(), resourceName, name)
}

func testAccMFADevicePolicyDefaultConfig_PingID_MobileApp_MissingIPPairing(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_mfa_device_policy_default" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  policy_type    = "PING_ONE_ID"
  name           = "%[3]s"

  mobile = {
    enabled = true
    applications = {
      "11111111-1111-1111-1111-111111111111" = {
        otp = { enabled = true }
        new_request_duration_configuration = {
          device_timeout = { duration = 25 }
          total_timeout  = { duration = 40 }
        }
      }
    }
  }
  desktop = { enabled = false }
  yubikey = { enabled = false }
  sms     = { enabled = false }
  voice   = { enabled = false }
  email   = { enabled = false }
  totp    = { enabled = false }
}`, acctest.WorkforceV2SandboxEnvironment(), resourceName, name)
}

func testAccMFADevicePolicyDefaultConfig_PingOneMFA_MobileApp_MissingAutoEnrollment(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_mfa_device_policy_default" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  policy_type    = "PING_ONE_MFA"
  name           = "%[4]s"

  mobile = {
    enabled = true
    applications = {
      "11111111-1111-1111-1111-111111111111" = {
        otp                  = { enabled = true }
        device_authorization = { enabled = true }
        integrity_detection  = "permissive"
      }
    }
  }
  sms   = { enabled = false }
  voice = { enabled = false }
  email = { enabled = false }
  totp  = { enabled = false }
}`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccMFADevicePolicyDefaultConfig_PingOneMFA_MobileApp_MissingDeviceAuthorization(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_mfa_device_policy_default" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  policy_type    = "PING_ONE_MFA"
  name           = "%[4]s"

  mobile = {
    enabled = true
    applications = {
      "11111111-1111-1111-1111-111111111111" = {
        otp                 = { enabled = true }
        auto_enrollment     = { enabled = true }
        integrity_detection = "permissive"
      }
    }
  }
  sms   = { enabled = false }
  voice = { enabled = false }
  email = { enabled = false }
  totp  = { enabled = false }
}`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccMFADevicePolicyDefaultConfig_PingOneMFA_MobileApp_MissingIntegrityDetection(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_mfa_device_policy_default" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  policy_type    = "PING_ONE_MFA"
  name           = "%[4]s"

  mobile = {
    enabled = true
    applications = {
      "11111111-1111-1111-1111-111111111111" = {
        otp                  = { enabled = true }
        auto_enrollment      = { enabled = true }
        device_authorization = { enabled = true }
      }
    }
  }
  sms   = { enabled = false }
  voice = { enabled = false }
  email = { enabled = false }
  totp  = { enabled = false }
}`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccMFADevicePolicyDefaultConfig_PingOneMFA_MobileApp_BiometricsEnabled(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_mfa_device_policy_default" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  policy_type    = "PING_ONE_MFA"
  name           = "%[4]s"

  mobile = {
    enabled = true
    applications = {
      "11111111-1111-1111-1111-111111111111" = {
        otp                  = { enabled = true }
        auto_enrollment      = { enabled = true }
        device_authorization = { enabled = true }
        integrity_detection  = "permissive"
        biometrics_enabled   = true
      }
    }
  }
  sms   = { enabled = false }
  voice = { enabled = false }
  email = { enabled = false }
  totp  = { enabled = false }
}`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccMFADevicePolicyDefaultConfig_PingOneMFA_MobileApp_NewRequestDuration(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_mfa_device_policy_default" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  policy_type    = "PING_ONE_MFA"
  name           = "%[4]s"

  mobile = {
    enabled = true
    applications = {
      "11111111-1111-1111-1111-111111111111" = {
        otp                  = { enabled = true }
        auto_enrollment      = { enabled = true }
        device_authorization = { enabled = true }
        integrity_detection  = "permissive"
        new_request_duration_configuration = {
          device_timeout = { duration = 25 }
          total_timeout  = { duration = 40 }
        }
      }
    }
  }
  sms   = { enabled = false }
  voice = { enabled = false }
  email = { enabled = false }
  totp  = { enabled = false }
}`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccMFADevicePolicyDefaultConfig_PingOneMFA_MobileApp_IPPairing(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_mfa_device_policy_default" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  policy_type    = "PING_ONE_MFA"
  name           = "%[4]s"

  mobile = {
    enabled = true
    applications = {
      "11111111-1111-1111-1111-111111111111" = {
        otp                      = { enabled = true }
        auto_enrollment          = { enabled = true }
        device_authorization     = { enabled = true }
        integrity_detection      = "permissive"
        ip_pairing_configuration = { any_ip_address = true }
      }
    }
  }
  sms   = { enabled = false }
  voice = { enabled = false }
  email = { enabled = false }
  totp  = { enabled = false }
}`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccMFADevicePolicyDefaultConfig_PingID_IPPairing_InvalidCIDR(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_mfa_device_policy_default" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  policy_type    = "PING_ONE_ID"
  name           = "%[3]s"

  mobile = {
    enabled = true
    applications = {
      "11111111-1111-1111-1111-111111111111" = {
        otp = { enabled = true }
        new_request_duration_configuration = {
          device_timeout = { duration = 25 }
          total_timeout  = { duration = 40 }
        }
        ip_pairing_configuration = {
          any_ip_address          = false
          only_these_ip_addresses = ["192.168.1.1"]
        }
      }
    }
  }
  desktop = { enabled = false }
  yubikey = { enabled = false }
  sms     = { enabled = false }
  voice   = { enabled = false }
  email   = { enabled = false }
  totp    = { enabled = false }
}`, acctest.WorkforceV2SandboxEnvironment(), resourceName, name)
}

func testAccMFADevicePolicyDefaultConfig_PingID_IPPairing_MissingIPs(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_mfa_device_policy_default" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  policy_type    = "PING_ONE_ID"
  name           = "%[3]s"

  mobile = {
    enabled = true
    applications = {
      "11111111-1111-1111-1111-111111111111" = {
        otp = { enabled = true }
        new_request_duration_configuration = {
          device_timeout = { duration = 25 }
          total_timeout  = { duration = 40 }
        }
        ip_pairing_configuration = {
          any_ip_address = false
        }
      }
    }
  }
  desktop = { enabled = false }
  yubikey = { enabled = false }
  sms     = { enabled = false }
  voice   = { enabled = false }
  email   = { enabled = false }
  totp    = { enabled = false }
}`, acctest.WorkforceV2SandboxEnvironment(), resourceName, name)
}

func testAccMFADevicePolicyDefaultConfig_PingID_Desktop_OTPCountHigh(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_mfa_device_policy_default" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  policy_type    = "PING_ONE_ID"
  name           = "%[3]s"

  mobile = { enabled = true }
  desktop = {
    enabled = true
    otp = {
      failure = {
        count     = 8
        cool_down = { duration = 2, time_unit = "MINUTES" }
      }
    }
  }
  yubikey = { enabled = false }
  sms     = { enabled = false }
  voice   = { enabled = false }
  email   = { enabled = false }
  totp    = { enabled = false }
}`, acctest.WorkforceV2SandboxEnvironment(), resourceName, name)
}

func testAccMFADevicePolicyDefaultConfig_PingID_Desktop_PairingKeyLifetimeHigh(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_mfa_device_policy_default" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  policy_type    = "PING_ONE_ID"
  name           = "%[3]s"

  mobile = { enabled = true }
  desktop = {
    enabled = true
    pairing_key_lifetime = {
      duration  = 50
      time_unit = "HOURS"
    }
  }
  yubikey = { enabled = false }
  sms     = { enabled = false }
  voice   = { enabled = false }
  email   = { enabled = false }
  totp    = { enabled = false }
}`, acctest.WorkforceV2SandboxEnvironment(), resourceName, name)
}

func testAccMFADevicePolicyDefaultConfig_RememberMe_MinutesHigh(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_mfa_device_policy_default" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  policy_type    = "PING_ONE_ID"
  name           = "%[3]s"

  remember_me = {
    web = {
      enabled = true
      life_time = {
        duration  = 129601
        time_unit = "MINUTES"
      }
    }
  }

  mobile  = { enabled = true }
  desktop = { enabled = false }
  yubikey = { enabled = false }
  sms     = { enabled = false }
  voice   = { enabled = false }
  email   = { enabled = false }
  totp    = { enabled = false }
}`, acctest.WorkforceV2SandboxEnvironment(), resourceName, name)
}

func testAccMFADevicePolicyDefaultConfig_RememberMe_HoursHigh(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_mfa_device_policy_default" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  policy_type    = "PING_ONE_ID"
  name           = "%[3]s"

  remember_me = {
    web = {
      enabled = true
      life_time = {
        duration  = 2161
        time_unit = "HOURS"
      }
    }
  }

  mobile  = { enabled = true }
  desktop = { enabled = false }
  yubikey = { enabled = false }
  sms     = { enabled = false }
  voice   = { enabled = false }
  email   = { enabled = false }
  totp    = { enabled = false }
}`, acctest.WorkforceV2SandboxEnvironment(), resourceName, name)
}

func testAccMFADevicePolicyDefaultConfig_RememberMe_DaysHigh(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_mfa_device_policy_default" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  policy_type    = "PING_ONE_ID"
  name           = "%[3]s"

  remember_me = {
    web = {
      enabled = true
      life_time = {
        duration  = 91
        time_unit = "DAYS"
      }
    }
  }

  mobile  = { enabled = true }
  desktop = { enabled = false }
  yubikey = { enabled = false }
  sms     = { enabled = false }
  voice   = { enabled = false }
  email   = { enabled = false }
  totp    = { enabled = false }
}`, acctest.WorkforceV2SandboxEnvironment(), resourceName, name)
}

func testAccMFADevicePolicyDefaultConfig_MobilePushLimit_CountHigh(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_mfa_device_policy_default" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  policy_type    = "PING_ONE_ID"
  name           = "%[3]s"

  mobile = {
    enabled = true
    applications = {
      "11111111-1111-1111-1111-111111111111" = {
        otp = { enabled = true }
        push_limit = {
          count         = 51
          lock_duration = { duration = 30, time_unit = "MINUTES" }
          time_period   = { duration = 10, time_unit = "MINUTES" }
        }
        new_request_duration_configuration = {
          device_timeout = { duration = 25 }
          total_timeout  = { duration = 40 }
        }
        ip_pairing_configuration = { any_ip_address = true }
      }
    }
  }
  desktop = { enabled = false }
  yubikey = { enabled = false }
  sms     = { enabled = false }
  voice   = { enabled = false }
  email   = { enabled = false }
  totp    = { enabled = false }
}`, acctest.WorkforceV2SandboxEnvironment(), resourceName, name)
}

func testAccMFADevicePolicyDefaultConfig_MobileNewRequestDuration_DeviceTimeoutHigh(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_mfa_device_policy_default" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  policy_type    = "PING_ONE_ID"
  name           = "%[3]s"

  mobile = {
    enabled = true
    applications = {
      "11111111-1111-1111-1111-111111111111" = {
        otp = { enabled = true }
        new_request_duration_configuration = {
          device_timeout = { duration = 76 }
          total_timeout  = { duration = 40 }
        }
        ip_pairing_configuration = { any_ip_address = true }
      }
    }
  }
  desktop = { enabled = false }
  yubikey = { enabled = false }
  sms     = { enabled = false }
  voice   = { enabled = false }
  email   = { enabled = false }
  totp    = { enabled = false }
}`, acctest.WorkforceV2SandboxEnvironment(), resourceName, name)
}

func testAccMFADevicePolicyDefaultConfig_MobileNewRequestDuration_TotalTimeoutHigh(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_mfa_device_policy_default" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  policy_type    = "PING_ONE_ID"
  name           = "%[3]s"

  mobile = {
    enabled = true
    applications = {
      "11111111-1111-1111-1111-111111111111" = {
        otp = { enabled = true }
        new_request_duration_configuration = {
          device_timeout = { duration = 25 }
          total_timeout  = { duration = 91 }
        }
        ip_pairing_configuration = { any_ip_address = true }
      }
    }
  }
  desktop = { enabled = false }
  yubikey = { enabled = false }
  sms     = { enabled = false }
  voice   = { enabled = false }
  email   = { enabled = false }
  totp    = { enabled = false }
}`, acctest.WorkforceV2SandboxEnvironment(), resourceName, name)
}

func testAccMFADevicePolicyDefaultConfig_Authentication(resourceName, name, deviceSelection string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_mfa_device_policy_default" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  policy_type    = "PING_ONE_ID"
  name           = "%[3]s"

  authentication = {
    device_selection = "%[4]s"
  }

  mobile  = { enabled = true }
  desktop = { enabled = false }
  yubikey = { enabled = false }
  sms     = { enabled = false }
  voice   = { enabled = false }
  email   = { enabled = false }
  totp    = { enabled = false }
}`, acctest.WorkforceV2SandboxEnvironment(), resourceName, name, deviceSelection)
}

func testAccMFADevicePolicyDefaultConfig_NewDeviceNotification(resourceName, name, notification string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_mfa_device_policy_default" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  policy_type    = "PING_ONE_ID"
  name           = "%[3]s"

  new_device_notification = "%[4]s"

  mobile  = { enabled = true }
  desktop = { enabled = false }
  yubikey = { enabled = false }
  sms     = { enabled = false }
  voice   = { enabled = false }
  email   = { enabled = false }
  totp    = { enabled = false }
}`, acctest.WorkforceV2SandboxEnvironment(), resourceName, name, notification)
}

func testAccMFADevicePolicyDefaultConfig_MobileOtpFailureCount(resourceName, name string, count int) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_mfa_device_policy_default" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  policy_type    = "PING_ONE_ID"
  name           = "%[3]s"

  mobile = {
    enabled = true
    otp = {
      failure = {
        count     = %[4]d
        cool_down = { duration = 2, time_unit = "MINUTES" }
      }
    }
    applications = {
      "11111111-1111-1111-1111-111111111111" = {
        otp = { enabled = true }
        new_request_duration_configuration = {
          device_timeout = { duration = 25 }
          total_timeout  = { duration = 40 }
        }
        ip_pairing_configuration = { any_ip_address = true }
      }
    }
  }
  desktop = { enabled = false }
  yubikey = { enabled = false }
  sms     = { enabled = false }
  voice   = { enabled = false }
  email   = { enabled = false }
  totp    = { enabled = false }
}`, acctest.WorkforceV2SandboxEnvironment(), resourceName, name, count)
}

func testAccMFADevicePolicyDefaultConfig_MobileNewRequestDuration_DeviceTimeout(resourceName, name string, duration int) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_mfa_device_policy_default" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  policy_type    = "PING_ONE_ID"
  name           = "%[3]s"

  mobile = {
    enabled = true
    applications = {
      "11111111-1111-1111-1111-111111111111" = {
        otp = { enabled = true }
        new_request_duration_configuration = {
          device_timeout = { duration = "%[4]d" }
          total_timeout  = { duration = 40 }
        }
        ip_pairing_configuration = { any_ip_address = true }
      }
    }
  }
  desktop = { enabled = false }
  yubikey = { enabled = false }
  sms     = { enabled = false }
  voice   = { enabled = false }
  email   = { enabled = false }
  totp    = { enabled = false }
}`, acctest.WorkforceV2SandboxEnvironment(), resourceName, name, duration)
}

func testAccMFADevicePolicyDefaultConfig_MobileNewRequestDuration_TotalTimeout(resourceName, name string, duration int) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_mfa_device_policy_default" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  policy_type    = "PING_ONE_ID"
  name           = "%[3]s"

  mobile = {
    enabled = true
    applications = {
      "11111111-1111-1111-1111-111111111111" = {
        otp = { enabled = true }
        new_request_duration_configuration = {
          device_timeout = { duration = 25 }
          total_timeout  = { duration = "%[4]d" }
        }
        ip_pairing_configuration = { any_ip_address = true }
      }
    }
  }
  desktop = { enabled = false }
  yubikey = { enabled = false }
  sms     = { enabled = false }
  voice   = { enabled = false }
  email   = { enabled = false }
  totp    = { enabled = false }
}`, acctest.WorkforceV2SandboxEnvironment(), resourceName, name, duration)
}

func testAccMFADevicePolicyDefaultConfig_BOMValidation(environmentName, licenseID, region, resourceName, name string) string {
	return fmt.Sprintf(`
resource "pingone_environment" "%[1]s" {
  name       = "%[1]s"
  type       = "SANDBOX"
  region     = "%[3]s"
  license_id = "%[2]s"
  services = [
    {
      type = "SSO"
    }
  ]
}

resource "pingone_mfa_device_policy_default" "%[4]s" {
  environment_id = pingone_environment.%[1]s.id
  policy_type    = "PING_ONE_MFA"
  name           = "%[5]s"

  mobile = { enabled = true }
  sms    = { enabled = false }
  voice  = { enabled = false }
  email  = { enabled = false }
  totp   = { enabled = false }
}`, environmentName, licenseID, region, resourceName, name)
}

func testAccMFADevicePolicyDefaultConfig_BOMValidation_PingIDEnv(environmentName, licenseID, region, resourceName, name string) string {
	return fmt.Sprintf(`
resource "pingone_environment" "%[1]s" {
  name       = "%[1]s"
  type       = "SANDBOX"
  region     = "%[3]s"
  license_id = "%[2]s"
  services = [
    {
      type = "SSO"
    },
    {
      type = "PingID"
    }
  ]
}

resource "pingone_mfa_device_policy_default" "%[4]s" {
  environment_id = pingone_environment.%[1]s.id
  policy_type    = "PING_ONE_MFA"
  name           = "%[5]s"

  mobile = { enabled = true }
  sms    = { enabled = false }
  voice  = { enabled = false }
  email  = { enabled = false }
  totp   = { enabled = false }
}`, environmentName, licenseID, region, resourceName, name)
}

func testAccMFADevicePolicyDefaultConfig_BOMValidation_MFAEnv(environmentName, licenseID, region, resourceName, name string) string {
	return fmt.Sprintf(`
resource "pingone_environment" "%[1]s" {
  name       = "%[1]s"
  type       = "SANDBOX"
  region     = "%[3]s"
  license_id = "%[2]s"
  services = [
    {
      type = "SSO"
    },
    {
      type = "MFA"
    }
  ]
}

resource "pingone_mfa_device_policy_default" "%[4]s" {
  environment_id = pingone_environment.%[1]s.id
  policy_type    = "PING_ONE_ID"
  name           = "%[5]s"

  mobile  = { enabled = true }
  desktop = { enabled = false }
  yubikey = { enabled = false }
  sms     = { enabled = false }
  voice   = { enabled = false }
  email   = { enabled = false }
  totp    = { enabled = false }
}`, environmentName, licenseID, region, resourceName, name)
}

func testAccMFADevicePolicyDefaultConfig_EmailOtpFailureCount(resourceName, name string, count int) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_mfa_device_policy_default" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  policy_type    = "PING_ONE_ID"
  name           = "%[3]s"

  email = {
    enabled = true
    otp = {
      failure = {
        count = %[4]d
        cool_down = {
          duration  = 5
          time_unit = "MINUTES"
        }
      }
    }
  }
  mobile  = { enabled = true }
  desktop = { enabled = false }
  yubikey = { enabled = false }
  sms     = { enabled = false }
  voice   = { enabled = false }
  totp    = { enabled = false }
}`, acctest.WorkforceV2SandboxEnvironment(), resourceName, name, count)
}

func testAccMFADevicePolicyDefaultConfig_EmailOtpFailureCoolDownDuration(resourceName, name string, duration int) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_mfa_device_policy_default" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  policy_type    = "PING_ONE_ID"
  name           = "%[3]s"

  email = {
    enabled = true
    otp = {
      failure = {
        count = 1
        cool_down = {
          duration  = %[4]d
          time_unit = "MINUTES"
        }
      }
    }
  }
  mobile  = { enabled = true }
  desktop = { enabled = false }
  yubikey = { enabled = false }
  sms     = { enabled = false }
  voice   = { enabled = false }
  totp    = { enabled = false }
}`, acctest.WorkforceV2SandboxEnvironment(), resourceName, name, duration)
}

func testAccMFADevicePolicyDefaultConfig_EmailOtpLifetimeDuration(resourceName, name string, duration int) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_mfa_device_policy_default" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  policy_type    = "PING_ONE_ID"
  name           = "%[3]s"

  email = {
    enabled = true
    otp = {
      lifetime = {
        duration  = %[4]d
        time_unit = "MINUTES"
      }
    }
  }
  mobile  = { enabled = true }
  desktop = { enabled = false }
  yubikey = { enabled = false }
  sms     = { enabled = false }
  voice   = { enabled = false }
  totp    = { enabled = false }
}`, acctest.WorkforceV2SandboxEnvironment(), resourceName, name, duration)
}

func testAccMFADevicePolicyDefaultConfig_SmsOtpFailureCount(resourceName, name string, count int) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_mfa_device_policy_default" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  policy_type    = "PING_ONE_ID"
  name           = "%[3]s"

  sms = {
    enabled = true
    otp = {
      failure = {
        count = %[4]d
        cool_down = {
          duration  = 5
          time_unit = "MINUTES"
        }
      }
    }
  }
  mobile  = { enabled = true }
  desktop = { enabled = false }
  yubikey = { enabled = false }
  email   = { enabled = false }
  voice   = { enabled = false }
  totp    = { enabled = false }
}`, acctest.WorkforceV2SandboxEnvironment(), resourceName, name, count)
}

func testAccMFADevicePolicyDefaultConfig_SmsOtpFailureCoolDownDuration(resourceName, name string, duration int) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_mfa_device_policy_default" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  policy_type    = "PING_ONE_ID"
  name           = "%[3]s"

  sms = {
    enabled = true
    otp = {
      failure = {
        count = 1
        cool_down = {
          duration  = %[4]d
          time_unit = "MINUTES"
        }
      }
    }
  }
  mobile  = { enabled = true }
  desktop = { enabled = false }
  yubikey = { enabled = false }
  email   = { enabled = false }
  voice   = { enabled = false }
  totp    = { enabled = false }
}`, acctest.WorkforceV2SandboxEnvironment(), resourceName, name, duration)
}

func testAccMFADevicePolicyDefaultConfig_SmsOtpLifetimeDuration(resourceName, name string, duration int) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_mfa_device_policy_default" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  policy_type    = "PING_ONE_ID"
  name           = "%[3]s"

  sms = {
    enabled = true
    otp = {
      lifetime = {
        duration  = %[4]d
        time_unit = "MINUTES"
      }
    }
  }
  mobile  = { enabled = true }
  desktop = { enabled = false }
  yubikey = { enabled = false }
  email   = { enabled = false }
  voice   = { enabled = false }
  totp    = { enabled = false }
}`, acctest.WorkforceV2SandboxEnvironment(), resourceName, name, duration)
}

func testAccMFADevicePolicyDefaultConfig_VoiceOtpFailureCount(resourceName, name string, count int) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_mfa_device_policy_default" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  policy_type    = "PING_ONE_ID"
  name           = "%[3]s"

  voice = {
    enabled = true
    otp = {
      failure = {
        count = %[4]d
        cool_down = {
          duration  = 5
          time_unit = "MINUTES"
        }
      }
    }
  }
  mobile  = { enabled = true }
  desktop = { enabled = false }
  yubikey = { enabled = false }
  email   = { enabled = false }
  sms     = { enabled = false }
  totp    = { enabled = false }
}`, acctest.WorkforceV2SandboxEnvironment(), resourceName, name, count)
}

func testAccMFADevicePolicyDefaultConfig_VoiceOtpFailureCoolDownDuration(resourceName, name string, duration int) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_mfa_device_policy_default" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  policy_type    = "PING_ONE_ID"
  name           = "%[3]s"

  voice = {
    enabled = true
    otp = {
      failure = {
        count = 1
        cool_down = {
          duration  = %[4]d
          time_unit = "MINUTES"
        }
      }
    }
  }
  mobile  = { enabled = true }
  desktop = { enabled = false }
  yubikey = { enabled = false }
  email   = { enabled = false }
  sms     = { enabled = false }
  totp    = { enabled = false }
}`, acctest.WorkforceV2SandboxEnvironment(), resourceName, name, duration)
}

func testAccMFADevicePolicyDefaultConfig_VoiceOtpLifetimeDuration(resourceName, name string, duration int) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_mfa_device_policy_default" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  policy_type    = "PING_ONE_ID"
  name           = "%[3]s"

  voice = {
    enabled = true
    otp = {
      lifetime = {
        duration  = %[4]d
        time_unit = "MINUTES"
      }
    }
  }
  mobile  = { enabled = true }
  desktop = { enabled = false }
  yubikey = { enabled = false }
  email   = { enabled = false }
  sms     = { enabled = false }
  totp    = { enabled = false }
}`, acctest.WorkforceV2SandboxEnvironment(), resourceName, name, duration)
}

func testAccMFADevicePolicyDefaultConfig_MobilePushLimitLockDuration(resourceName, name string, duration int) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_mfa_device_policy_default" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  policy_type    = "PING_ONE_ID"
  name           = "%[3]s"

  mobile = {
    enabled = true
    applications = {
      "11111111-1111-1111-1111-111111111111" = {
        type = "pingIdAppConfig"
        otp  = { enabled = true }
        push_limit = {
          lock_duration = {
            duration  = %[4]d
            time_unit = "MINUTES"
          }
        }
        new_request_duration_configuration = {
          device_timeout = { duration = 25 }
          total_timeout  = { duration = 40 }
        }
        ip_pairing_configuration = { any_ip_address = true }
      }
    }
  }
  desktop = { enabled = false }
  yubikey = { enabled = false }
  sms     = { enabled = false }
  voice   = { enabled = false }
  email   = { enabled = false }
  totp    = { enabled = false }
}`, acctest.WorkforceV2SandboxEnvironment(), resourceName, name, duration)
}

func testAccMFADevicePolicyDefaultConfig_MobilePushLimitTimePeriod(resourceName, name string, duration int) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_mfa_device_policy_default" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  policy_type    = "PING_ONE_ID"
  name           = "%[3]s"

  mobile = {
    enabled = true
    applications = {
      "11111111-1111-1111-1111-111111111111" = {
        type = "pingIdAppConfig"
        otp  = { enabled = true }
        push_limit = {
          time_period = {
            duration  = %[4]d
            time_unit = "MINUTES"
          }
        }
        new_request_duration_configuration = {
          device_timeout = { duration = 25 }
          total_timeout  = { duration = 40 }
        }
        ip_pairing_configuration = { any_ip_address = true }
      }
    }
  }
  desktop = { enabled = false }
  yubikey = { enabled = false }
  sms     = { enabled = false }
  voice   = { enabled = false }
  email   = { enabled = false }
  totp    = { enabled = false }
}`, acctest.WorkforceV2SandboxEnvironment(), resourceName, name, duration)
}

func testAccMFADevicePolicyDefaultConfig_MobilePushTimeoutDuration(environmentName, licenseID, resourceName, name string, duration int) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_mfa_device_policy_default" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  policy_type    = "PING_ONE_MFA"
  name           = "%[4]s"

  mobile = {
    enabled = true
    applications = {
      "11111111-1111-1111-1111-111111111111" = {
        otp = { enabled = true }
        push_timeout = {
          duration  = %[5]d
          time_unit = "SECONDS"
        }
        auto_enrollment      = { enabled = true }
        device_authorization = { enabled = true }
        integrity_detection  = "permissive"
      }
    }
  }
  sms   = { enabled = false }
  voice = { enabled = false }
  email = { enabled = false }
  totp  = { enabled = false }
}`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name, duration)
}

func testAccMFADevicePolicyDefaultConfig_MobileOtpFailureCoolDownDuration(resourceName, name string, duration int) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_mfa_device_policy_default" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  policy_type    = "PING_ONE_ID"
  name           = "%[3]s"

  mobile = {
    enabled = true
    applications = {
      "11111111-1111-1111-1111-111111111111" = {
        type = "pingIdAppConfig"
        otp  = { enabled = true }
        new_request_duration_configuration = {
          device_timeout = { duration = 25 }
          total_timeout  = { duration = 40 }
        }
        ip_pairing_configuration = { any_ip_address = true }
      }
    }
    otp = {
      failure = {
        count = 1
        cool_down = {
          duration  = %[4]d
          time_unit = "MINUTES"
        }
      }
    }
  }
  desktop = { enabled = false }
  yubikey = { enabled = false }
  sms     = { enabled = false }
  voice   = { enabled = false }
  email   = { enabled = false }
  totp    = { enabled = false }
}`, acctest.WorkforceV2SandboxEnvironment(), resourceName, name, duration)
}

func testAccMFADevicePolicyDefaultConfig_TotpOtpFailureCoolDownDuration(resourceName, name string, duration int) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_mfa_device_policy_default" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  policy_type    = "PING_ONE_ID"
  name           = "%[3]s"

  totp = {
    enabled = true
    otp = {
      failure = {
        count = 1
        cool_down = {
          duration  = %[4]d
          time_unit = "MINUTES"
        }
      }
    }
  }
  mobile  = { enabled = true }
  desktop = { enabled = false }
  yubikey = { enabled = false }
  sms     = { enabled = false }
  voice   = { enabled = false }
  email   = { enabled = false }
}`, acctest.WorkforceV2SandboxEnvironment(), resourceName, name, duration)
}
