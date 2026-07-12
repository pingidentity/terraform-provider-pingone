// Copyright © 2026 Ping Identity Corporation

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

func TestAccMFADevicePolicy_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_device_policy.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	var mfaDevicePolicyID, environmentID string

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
		CheckDestroy:             mfa.MFADevicePolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccMFADevicePolicyConfig_FullSMS(resourceName, name),
				Check:  mfa.MFADevicePolicy_GetIDs(resourceFullName, &environmentID, &mfaDevicePolicyID),
			},
			{
				PreConfig: func() {
					mfa.MFADevicePolicy_RemovalDrift_PreConfig(ctx, p1Client.API.MFAAPIClient, t, environmentID, mfaDevicePolicyID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the environment
			{
				Config: testAccMFADevicePolicyConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check:  mfa.MFADevicePolicy_GetIDs(resourceFullName, &environmentID, &mfaDevicePolicyID),
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

func TestAccMFADevicePolicy_NewEnv(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_device_policy.%s", resourceName)

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
		CheckDestroy:             mfa.MFADevicePolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFADevicePolicyConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
				),
			},
		},
	})
}

func TestAccMFADevicePolicy_SMS_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_device_policy.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.MFADevicePolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFADevicePolicyConfig_FullSMS(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "authentication.device_selection", "DEFAULT_TO_FIRST"),
					resource.TestCheckResourceAttr(resourceFullName, "ignore_user_lock", "true"),
					resource.TestMatchResourceAttr(resourceFullName, "notifications_policy.id", verify.P1ResourceIDRegexpFullString),
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
					resource.TestCheckResourceAttr(resourceFullName, "voice.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "whats_app.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "SMS_THEN_EMAIL"),
					// Backward compatibility (CDI-1259): a config that omits `policy_type` and the
					// new PingID-only attributes must still plan/apply cleanly, with `policy_type`
					// computed to `PING_ONE_MFA` and the PingID-only attributes null.
					resource.TestCheckResourceAttr(resourceFullName, "policy_type", "PING_ONE_MFA"),
					resource.TestCheckNoResourceAttr(resourceFullName, "desktop"),
					resource.TestCheckNoResourceAttr(resourceFullName, "yubikey"),
					resource.TestCheckResourceAttr(resourceFullName, "oath_token.enabled", "false"),
				),
			},
			// Test importing the resource
			{
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

func TestAccMFADevicePolicy_SMS_Minimal(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_device_policy.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.MFADevicePolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFADevicePolicyConfig_MinimalSMS(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "authentication.device_selection", "DEFAULT_TO_FIRST"),
					resource.TestCheckResourceAttr(resourceFullName, "ignore_user_lock", "false"),
					resource.TestCheckNoResourceAttr(resourceFullName, "notifications_policy.id"),
					resource.TestCheckResourceAttr(resourceFullName, "remember_me.web.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "remember_me.web.life_time.duration", "30"),
					resource.TestCheckResourceAttr(resourceFullName, "remember_me.web.life_time.time_unit", "MINUTES"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.pairing_disabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.otp.lifetime.duration", "30"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.otp.lifetime.time_unit", "MINUTES"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.otp.failure.count", "3"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.otp.failure.cool_down.duration", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.otp.failure.cool_down.time_unit", "MINUTES"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.otp.otp_length", "6"),
					resource.TestCheckNoResourceAttr(resourceFullName, "sms.prompt_for_nickname_on_pairing"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "whats_app.enabled", "false"),
					resource.TestCheckNoResourceAttr(resourceFullName, "whats_app.pairing_disabled"),
					resource.TestCheckResourceAttr(resourceFullName, "whats_app.otp.lifetime.duration", "30"),
					resource.TestCheckResourceAttr(resourceFullName, "whats_app.otp.lifetime.time_unit", "MINUTES"),
					resource.TestCheckResourceAttr(resourceFullName, "whats_app.otp.failure.count", "3"),
					resource.TestCheckResourceAttr(resourceFullName, "whats_app.otp.failure.cool_down.duration", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "whats_app.otp.failure.cool_down.time_unit", "MINUTES"),
					resource.TestCheckResourceAttr(resourceFullName, "whats_app.otp.otp_length", "6"),
					resource.TestCheckNoResourceAttr(resourceFullName, "whats_app.prompt_for_nickname_on_pairing"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "NONE"),
					// Backward compatibility (CDI-1259): a config that omits `policy_type` and the
					// new PingID-only attributes must still plan/apply cleanly, with `policy_type`
					// computed to `PING_ONE_MFA` and the PingID-only attributes null.
					resource.TestCheckResourceAttr(resourceFullName, "policy_type", "PING_ONE_MFA"),
					resource.TestCheckNoResourceAttr(resourceFullName, "desktop"),
					resource.TestCheckNoResourceAttr(resourceFullName, "yubikey"),
					resource.TestCheckResourceAttr(resourceFullName, "oath_token.enabled", "false"),
				),
			},
			// Re-apply the same, pre-existing-shaped config a second time to prove there is no
			// perpetual diff once `policy_type` is computed and stored in state (AC4 / Task 4's
			// no-perpetual-diff guarantee, exercised end-to-end here against the real API).
			{
				Config: testAccMFADevicePolicyConfig_MinimalSMS(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "policy_type", "PING_ONE_MFA"),
					resource.TestCheckNoResourceAttr(resourceFullName, "desktop"),
					resource.TestCheckNoResourceAttr(resourceFullName, "yubikey"),
					resource.TestCheckResourceAttr(resourceFullName, "oath_token.enabled", "false"),
				),
				ExpectNonEmptyPlan: false,
			},
			// Test importing the resource, confirming a clean round-trip for the computed
			// `policy_type` and the null PingID-only attributes.
			{
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

func TestAccMFADevicePolicy_SMS_Change(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_device_policy.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.MFADevicePolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFADevicePolicyConfig_FullSMS(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "authentication.device_selection", "DEFAULT_TO_FIRST"),
					resource.TestCheckResourceAttr(resourceFullName, "ignore_user_lock", "true"),
					resource.TestMatchResourceAttr(resourceFullName, "notifications_policy.id", verify.P1ResourceIDRegexpFullString),
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
					resource.TestCheckResourceAttr(resourceFullName, "voice.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "whats_app.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "SMS_THEN_EMAIL"),
				),
			},
			{
				Config: testAccMFADevicePolicyConfig_MinimalSMS(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "authentication.device_selection", "DEFAULT_TO_FIRST"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.pairing_disabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.otp.lifetime.duration", "30"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.otp.lifetime.time_unit", "MINUTES"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.otp.failure.count", "3"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.otp.failure.cool_down.duration", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.otp.failure.cool_down.time_unit", "MINUTES"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.otp.otp_length", "6"),
					resource.TestCheckNoResourceAttr(resourceFullName, "sms.prompt_for_nickname_on_pairing"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "whats_app.enabled", "false"),
					resource.TestCheckNoResourceAttr(resourceFullName, "whats_app.pairing_disabled"),
					resource.TestCheckResourceAttr(resourceFullName, "whats_app.otp.lifetime.duration", "30"),
					resource.TestCheckResourceAttr(resourceFullName, "whats_app.otp.lifetime.time_unit", "MINUTES"),
					resource.TestCheckResourceAttr(resourceFullName, "whats_app.otp.failure.count", "3"),
					resource.TestCheckResourceAttr(resourceFullName, "whats_app.otp.failure.cool_down.duration", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "whats_app.otp.failure.cool_down.time_unit", "MINUTES"),
					resource.TestCheckResourceAttr(resourceFullName, "whats_app.otp.otp_length", "6"),
					resource.TestCheckNoResourceAttr(resourceFullName, "whats_app.prompt_for_nickname_on_pairing"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "NONE"),
				),
			},
			{
				Config: testAccMFADevicePolicyConfig_FullSMS(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "authentication.device_selection", "DEFAULT_TO_FIRST"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.pairing_disabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.otp.lifetime.duration", "75"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.otp.lifetime.time_unit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.otp.failure.count", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.otp.failure.cool_down.duration", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.otp.failure.cool_down.time_unit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.otp.otp_length", "7"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.prompt_for_nickname_on_pairing", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "whats_app.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "SMS_THEN_EMAIL"),
				),
			},
		},
	})
}

func TestAccMFADevicePolicy_Voice_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_device_policy.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.MFADevicePolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFADevicePolicyConfig_FullVoice(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "authentication.device_selection", "PROMPT_TO_SELECT"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.pairing_disabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.otp.lifetime.duration", "75"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.otp.lifetime.time_unit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.otp.failure.count", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.otp.failure.cool_down.duration", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.otp.failure.cool_down.time_unit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.otp.otp_length", "7"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.prompt_for_nickname_on_pairing", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "email.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "SMS_THEN_EMAIL"),
				),
			},
			// Test importing the resource
			{
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

func TestAccMFADevicePolicy_Voice_Minimal(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_device_policy.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.MFADevicePolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFADevicePolicyConfig_MinimalVoice(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "authentication.device_selection", "DEFAULT_TO_FIRST"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.pairing_disabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.otp.lifetime.duration", "30"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.otp.lifetime.time_unit", "MINUTES"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.otp.failure.count", "3"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.otp.failure.cool_down.duration", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.otp.failure.cool_down.time_unit", "MINUTES"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.otp.otp_length", "6"),
					resource.TestCheckNoResourceAttr(resourceFullName, "voice.prompt_for_nickname_on_pairing"),
					resource.TestCheckResourceAttr(resourceFullName, "email.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "NONE"),
				),
			},
		},
	})
}

func TestAccMFADevicePolicy_Voice_Change(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_device_policy.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.MFADevicePolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFADevicePolicyConfig_FullVoice(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "authentication.device_selection", "PROMPT_TO_SELECT"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.pairing_disabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.otp.lifetime.duration", "75"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.otp.lifetime.time_unit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.otp.failure.count", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.otp.failure.cool_down.duration", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.otp.failure.cool_down.time_unit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.otp.otp_length", "7"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.prompt_for_nickname_on_pairing", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "email.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "SMS_THEN_EMAIL"),
				),
			},
			{
				Config: testAccMFADevicePolicyConfig_MinimalVoice(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "authentication.device_selection", "DEFAULT_TO_FIRST"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.pairing_disabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.otp.lifetime.duration", "30"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.otp.lifetime.time_unit", "MINUTES"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.otp.failure.count", "3"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.otp.failure.cool_down.duration", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.otp.failure.cool_down.time_unit", "MINUTES"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.otp.otp_length", "6"),
					resource.TestCheckNoResourceAttr(resourceFullName, "voice.prompt_for_nickname_on_pairing"),
					resource.TestCheckResourceAttr(resourceFullName, "email.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "NONE"),
				),
			},
			{
				Config: testAccMFADevicePolicyConfig_FullVoice(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "authentication.device_selection", "PROMPT_TO_SELECT"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.pairing_disabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.otp.lifetime.duration", "75"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.otp.lifetime.time_unit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.otp.failure.count", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.otp.failure.cool_down.duration", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.otp.failure.cool_down.time_unit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.otp.otp_length", "7"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.prompt_for_nickname_on_pairing", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "email.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "SMS_THEN_EMAIL"),
				),
			},
		},
	})
}

func TestAccMFADevicePolicy_Email_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_device_policy.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.MFADevicePolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFADevicePolicyConfig_FullEmail(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "authentication.device_selection", "ALWAYS_DISPLAY_DEVICES"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "email.pairing_disabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "email.otp.lifetime.duration", "75"),
					resource.TestCheckResourceAttr(resourceFullName, "email.otp.lifetime.time_unit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "email.otp.failure.count", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "email.otp.failure.cool_down.duration", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "email.otp.failure.cool_down.time_unit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "email.otp.otp_length", "7"),
					resource.TestCheckResourceAttr(resourceFullName, "email.prompt_for_nickname_on_pairing", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "SMS_THEN_EMAIL"),
				),
			},
			// Test importing the resource
			{
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

func TestAccMFADevicePolicy_Email_Minimal(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_device_policy.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.MFADevicePolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFADevicePolicyConfig_MinimalEmail(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "authentication.device_selection", "DEFAULT_TO_FIRST"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "email.pairing_disabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.otp.lifetime.duration", "30"),
					resource.TestCheckResourceAttr(resourceFullName, "email.otp.lifetime.time_unit", "MINUTES"),
					resource.TestCheckResourceAttr(resourceFullName, "email.otp.failure.count", "3"),
					resource.TestCheckResourceAttr(resourceFullName, "email.otp.failure.cool_down.duration", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "email.otp.failure.cool_down.time_unit", "MINUTES"),
					resource.TestCheckResourceAttr(resourceFullName, "email.otp.otp_length", "6"),
					resource.TestCheckNoResourceAttr(resourceFullName, "email.prompt_for_nickname_on_pairing"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "NONE"),
				),
			},
		},
	})
}

func TestAccMFADevicePolicy_Email_Change(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_device_policy.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.MFADevicePolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFADevicePolicyConfig_FullEmail(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "authentication.device_selection", "ALWAYS_DISPLAY_DEVICES"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "email.pairing_disabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "email.otp.lifetime.duration", "75"),
					resource.TestCheckResourceAttr(resourceFullName, "email.otp.lifetime.time_unit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "email.otp.failure.count", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "email.otp.failure.cool_down.duration", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "email.otp.failure.cool_down.time_unit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "email.otp.otp_length", "7"),
					resource.TestCheckResourceAttr(resourceFullName, "email.prompt_for_nickname_on_pairing", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "SMS_THEN_EMAIL"),
				),
			},
			{
				Config: testAccMFADevicePolicyConfig_MinimalEmail(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "authentication.device_selection", "DEFAULT_TO_FIRST"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "email.pairing_disabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.otp.lifetime.duration", "30"),
					resource.TestCheckResourceAttr(resourceFullName, "email.otp.lifetime.time_unit", "MINUTES"),
					resource.TestCheckResourceAttr(resourceFullName, "email.otp.failure.count", "3"),
					resource.TestCheckResourceAttr(resourceFullName, "email.otp.failure.cool_down.duration", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "email.otp.failure.cool_down.time_unit", "MINUTES"),
					resource.TestCheckResourceAttr(resourceFullName, "email.otp.otp_length", "6"),
					resource.TestCheckNoResourceAttr(resourceFullName, "email.prompt_for_nickname_on_pairing"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "NONE"),
				),
			},
			{
				Config: testAccMFADevicePolicyConfig_FullEmail(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "authentication.device_selection", "ALWAYS_DISPLAY_DEVICES"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "email.pairing_disabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "email.otp.lifetime.duration", "75"),
					resource.TestCheckResourceAttr(resourceFullName, "email.otp.lifetime.time_unit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "email.otp.failure.count", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "email.otp.failure.cool_down.duration", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "email.otp.failure.cool_down.time_unit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "email.otp.otp_length", "7"),
					resource.TestCheckResourceAttr(resourceFullName, "email.prompt_for_nickname_on_pairing", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "SMS_THEN_EMAIL"),
				),
			},
		},
	})
}

func TestAccMFADevicePolicy_Mobile_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_device_policy.%s", resourceName)

	name := resourceName

	firebaseCredentials := os.Getenv("PINGONE_GOOGLE_FIREBASE_CREDENTIALS")

	application1FullName := fmt.Sprintf("pingone_application.%s-1", name)
	application2FullName := fmt.Sprintf("pingone_application.%s-2", name)
	application3FullName := fmt.Sprintf("pingone_application.%s-3", name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
			acctest.PreCheckGoogleFirebaseCredentials(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.MFADevicePolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFADevicePolicyConfig_FullMobile(resourceName, name, firebaseCredentials),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.otp.failure.count", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.otp.failure.cool_down.duration", "125"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.otp.failure.cool_down.time_unit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.applications.%", "3"),

					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application1FullName, "mobile.applications.%s.pairing_disabled", "true"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application1FullName, "mobile.applications.%s.push.enabled", "true"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application1FullName, "mobile.applications.%s.push.number_matching.enabled", "true"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application1FullName, "mobile.applications.%s.push_timeout.duration", "100"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application1FullName, "mobile.applications.%s.push_timeout.time_unit", "SECONDS"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application1FullName, "mobile.applications.%s.pairing_key_lifetime.duration", "3"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application1FullName, "mobile.applications.%s.pairing_key_lifetime.time_unit", "HOURS"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application1FullName, "mobile.applications.%s.push_limit.count", "10"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application1FullName, "mobile.applications.%s.push_limit.lock_duration.duration", "260"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application1FullName, "mobile.applications.%s.push_limit.lock_duration.time_unit", "SECONDS"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application1FullName, "mobile.applications.%s.push_limit.time_period.duration", "300"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application1FullName, "mobile.applications.%s.push_limit.time_period.time_unit", "SECONDS"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application1FullName, "mobile.applications.%s.otp.enabled", "true"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application1FullName, "mobile.applications.%s.device_authorization.enabled", "true"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application1FullName, "mobile.applications.%s.device_authorization.extra_verification", "restrictive"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application1FullName, "mobile.applications.%s.auto_enrollment.enabled", "true"),
					mfa.TestCheckMFADevicePolicyApplicationMapNoResourceAttr(resourceFullName, application1FullName, "mobile.applications.%s.integrity_detection"),

					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application2FullName, "mobile.applications.%s.pairing_disabled", "false"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application2FullName, "mobile.applications.%s.push.enabled", "false"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application2FullName, "mobile.applications.%s.push.number_matching.enabled", "false"),
					mfa.TestCheckMFADevicePolicyApplicationMapNoResourceAttr(resourceFullName, application2FullName, "mobile.applications.%s.push_timeout.duration"),
					mfa.TestCheckMFADevicePolicyApplicationMapNoResourceAttr(resourceFullName, application2FullName, "mobile.applications.%s.push_timeout.time_unit"),
					mfa.TestCheckMFADevicePolicyApplicationMapNoResourceAttr(resourceFullName, application2FullName, "mobile.applications.%s.pairing_key_lifetime.duration"),
					mfa.TestCheckMFADevicePolicyApplicationMapNoResourceAttr(resourceFullName, application2FullName, "mobile.applications.%s.pairing_key_lifetime.time_unit"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application2FullName, "mobile.applications.%s.push_limit.count", "5"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application2FullName, "mobile.applications.%s.push_limit.lock_duration.duration", "30"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application2FullName, "mobile.applications.%s.push_limit.lock_duration.time_unit", "MINUTES"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application2FullName, "mobile.applications.%s.push_limit.time_period.duration", "10"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application2FullName, "mobile.applications.%s.push_limit.time_period.time_unit", "MINUTES"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application2FullName, "mobile.applications.%s.otp.enabled", "true"),
					mfa.TestCheckMFADevicePolicyApplicationMapNoResourceAttr(resourceFullName, application2FullName, "mobile.applications.%s.device_authorization.extra_verification"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application2FullName, "mobile.applications.%s.integrity_detection", "permissive"),

					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application3FullName, "mobile.applications.%s.pairing_disabled", "false"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application3FullName, "mobile.applications.%s.push.enabled", "true"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application3FullName, "mobile.applications.%s.push.number_matching.enabled", "false"),
					mfa.TestCheckMFADevicePolicyApplicationMapNoResourceAttr(resourceFullName, application3FullName, "mobile.applications.%s.push_timeout.duration"),
					mfa.TestCheckMFADevicePolicyApplicationMapNoResourceAttr(resourceFullName, application3FullName, "mobile.applications.%s.push_timeout.time_unit"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application3FullName, "mobile.applications.%s.pairing_key_lifetime.duration", "55"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application3FullName, "mobile.applications.%s.pairing_key_lifetime.time_unit", "MINUTES"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application3FullName, "mobile.applications.%s.push_limit.count", "5"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application3FullName, "mobile.applications.%s.push_limit.lock_duration.duration", "25"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application3FullName, "mobile.applications.%s.push_limit.lock_duration.time_unit", "MINUTES"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application3FullName, "mobile.applications.%s.push_limit.time_period.duration", "5"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application3FullName, "mobile.applications.%s.push_limit.time_period.time_unit", "MINUTES"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application3FullName, "mobile.applications.%s.otp.enabled", "true"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application3FullName, "mobile.applications.%s.device_authorization.enabled", "false"),
					mfa.TestCheckMFADevicePolicyApplicationMapNoResourceAttr(resourceFullName, application3FullName, "mobile.applications.%s.device_authorization.extra_verification"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application3FullName, "mobile.applications.%s.auto_enrollment.enabled", "true"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application3FullName, "mobile.applications.%s.integrity_detection", "permissive"),

					resource.TestCheckResourceAttr(resourceFullName, "mobile.prompt_for_nickname_on_pairing", "true"),

					resource.TestCheckResourceAttr(resourceFullName, "totp.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "SMS_THEN_EMAIL"),
				),
			},
			// Test importing the resource
			{
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

func TestAccMFADevicePolicy_Mobile_IntegrityDetectionErrors(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.MFADevicePolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccMFADevicePolicyConfig_MobileIntegrityDetectionError_1(resourceName, name),
				ExpectError: regexp.MustCompile("Mobile application integrity detection setting has no effect"),
			},
			{
				Config:      testAccMFADevicePolicyConfig_MobileIntegrityDetectionError_2(resourceName, name),
				ExpectError: regexp.MustCompile("Invalid mobile application integrity detection setting"),
			},
		},
	})
}

func TestAccMFADevicePolicy_Mobile_BadMFADevicePolicyErrors(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.MFADevicePolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccMFADevicePolicyConfig_MobileBadApplicationError_1(resourceName, name),
				ExpectError: regexp.MustCompile("Application not found"),
			},
			{
				Config:      testAccMFADevicePolicyConfig_MobileBadApplicationError_2(resourceName, name),
				ExpectError: regexp.MustCompile("Invalid application type"),
			},
			{
				Config:      testAccMFADevicePolicyConfig_MobileBadApplicationError_3(resourceName, name),
				ExpectError: regexp.MustCompile("Invalid application type"),
			},
			{
				Config:      testAccMFADevicePolicyConfig_MobileBadApplicationError_4(resourceName, name),
				ExpectError: regexp.MustCompile("Missing application configuration"),
			},
		},
	})
}

func TestAccMFADevicePolicy_Mobile_Minimal(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_device_policy.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.MFADevicePolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFADevicePolicyConfig_MinimalMobile(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.otp.failure.count", "3"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.otp.failure.cool_down.duration", "2"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.otp.failure.cool_down.time_unit", "MINUTES"),
					resource.TestCheckNoResourceAttr(resourceFullName, "mobile.prompt_for_nickname_on_pairing"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "NONE"),
				),
			},
		},
	})
}

func TestAccMFADevicePolicy_Mobile_Change(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_device_policy.%s", resourceName)

	name := resourceName

	firebaseCredentials := os.Getenv("PINGONE_GOOGLE_FIREBASE_CREDENTIALS")

	application1FullName := fmt.Sprintf("pingone_application.%s-1", name)
	application2FullName := fmt.Sprintf("pingone_application.%s-2", name)
	application3FullName := fmt.Sprintf("pingone_application.%s-3", name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
			acctest.PreCheckGoogleFirebaseCredentials(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.MFADevicePolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFADevicePolicyConfig_FullMobile(resourceName, name, firebaseCredentials),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.otp.failure.count", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.otp.failure.cool_down.duration", "125"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.otp.failure.cool_down.time_unit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.applications.%", "3"),

					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application1FullName, "mobile.applications.%s.pairing_disabled", "true"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application1FullName, "mobile.applications.%s.push.enabled", "true"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application1FullName, "mobile.applications.%s.push.number_matching.enabled", "true"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application1FullName, "mobile.applications.%s.push_timeout.duration", "100"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application1FullName, "mobile.applications.%s.push_timeout.time_unit", "SECONDS"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application1FullName, "mobile.applications.%s.pairing_key_lifetime.duration", "3"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application1FullName, "mobile.applications.%s.pairing_key_lifetime.time_unit", "HOURS"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application1FullName, "mobile.applications.%s.push_limit.count", "10"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application1FullName, "mobile.applications.%s.push_limit.lock_duration.duration", "260"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application1FullName, "mobile.applications.%s.push_limit.lock_duration.time_unit", "SECONDS"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application1FullName, "mobile.applications.%s.push_limit.time_period.duration", "300"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application1FullName, "mobile.applications.%s.push_limit.time_period.time_unit", "SECONDS"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application1FullName, "mobile.applications.%s.otp.enabled", "true"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application1FullName, "mobile.applications.%s.device_authorization.enabled", "true"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application1FullName, "mobile.applications.%s.device_authorization.extra_verification", "restrictive"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application1FullName, "mobile.applications.%s.auto_enrollment.enabled", "true"),
					mfa.TestCheckMFADevicePolicyApplicationMapNoResourceAttr(resourceFullName, application1FullName, "mobile.applications.%s.integrity_detection"),

					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application2FullName, "mobile.applications.%s.pairing_disabled", "false"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application2FullName, "mobile.applications.%s.push.enabled", "false"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application2FullName, "mobile.applications.%s.push.number_matching.enabled", "false"),
					mfa.TestCheckMFADevicePolicyApplicationMapNoResourceAttr(resourceFullName, application2FullName, "mobile.applications.%s.push_timeout.duration"),
					mfa.TestCheckMFADevicePolicyApplicationMapNoResourceAttr(resourceFullName, application2FullName, "mobile.applications.%s.push_timeout.time_unit"),
					mfa.TestCheckMFADevicePolicyApplicationMapNoResourceAttr(resourceFullName, application2FullName, "mobile.applications.%s.pairing_key_lifetime.duration"),
					mfa.TestCheckMFADevicePolicyApplicationMapNoResourceAttr(resourceFullName, application2FullName, "mobile.applications.%s.pairing_key_lifetime.time_unit"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application2FullName, "mobile.applications.%s.push_limit.count", "5"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application2FullName, "mobile.applications.%s.push_limit.lock_duration.duration", "30"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application2FullName, "mobile.applications.%s.push_limit.lock_duration.time_unit", "MINUTES"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application2FullName, "mobile.applications.%s.push_limit.time_period.duration", "10"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application2FullName, "mobile.applications.%s.push_limit.time_period.time_unit", "MINUTES"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application2FullName, "mobile.applications.%s.otp.enabled", "true"),
					mfa.TestCheckMFADevicePolicyApplicationMapNoResourceAttr(resourceFullName, application2FullName, "mobile.applications.%s.device_authorization.extra_verification"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application2FullName, "mobile.applications.%s.integrity_detection", "permissive"),

					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application3FullName, "mobile.applications.%s.pairing_disabled", "false"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application3FullName, "mobile.applications.%s.push.enabled", "true"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application3FullName, "mobile.applications.%s.push.number_matching.enabled", "false"),
					mfa.TestCheckMFADevicePolicyApplicationMapNoResourceAttr(resourceFullName, application3FullName, "mobile.applications.%s.push_timeout.duration"),
					mfa.TestCheckMFADevicePolicyApplicationMapNoResourceAttr(resourceFullName, application3FullName, "mobile.applications.%s.push_timeout.time_unit"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application3FullName, "mobile.applications.%s.pairing_key_lifetime.duration", "55"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application3FullName, "mobile.applications.%s.pairing_key_lifetime.time_unit", "MINUTES"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application3FullName, "mobile.applications.%s.push_limit.count", "5"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application3FullName, "mobile.applications.%s.push_limit.lock_duration.duration", "25"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application3FullName, "mobile.applications.%s.push_limit.lock_duration.time_unit", "MINUTES"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application3FullName, "mobile.applications.%s.push_limit.time_period.duration", "5"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application3FullName, "mobile.applications.%s.push_limit.time_period.time_unit", "MINUTES"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application3FullName, "mobile.applications.%s.otp.enabled", "true"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application3FullName, "mobile.applications.%s.device_authorization.enabled", "false"),
					mfa.TestCheckMFADevicePolicyApplicationMapNoResourceAttr(resourceFullName, application3FullName, "mobile.applications.%s.device_authorization.extra_verification"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application3FullName, "mobile.applications.%s.auto_enrollment.enabled", "true"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application3FullName, "mobile.applications.%s.integrity_detection", "permissive"),

					resource.TestCheckResourceAttr(resourceFullName, "mobile.prompt_for_nickname_on_pairing", "true"),

					resource.TestCheckResourceAttr(resourceFullName, "totp.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "SMS_THEN_EMAIL"),
				),
			},
			{
				Config: testAccMFADevicePolicyConfig_MinimalMobile(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.otp.failure.count", "3"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.otp.failure.cool_down.duration", "2"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.otp.failure.cool_down.time_unit", "MINUTES"),
					resource.TestCheckNoResourceAttr(resourceFullName, "mobile.prompt_for_nickname_on_pairing"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "NONE"),
				),
			},
			{
				Config: testAccMFADevicePolicyConfig_FullMobile(resourceName, name, firebaseCredentials),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.otp.failure.count", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.otp.failure.cool_down.duration", "125"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.otp.failure.cool_down.time_unit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.applications.%", "3"),

					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application1FullName, "mobile.applications.%s.pairing_disabled", "true"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application1FullName, "mobile.applications.%s.push.enabled", "true"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application1FullName, "mobile.applications.%s.push.number_matching.enabled", "true"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application1FullName, "mobile.applications.%s.push_timeout.duration", "100"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application1FullName, "mobile.applications.%s.push_timeout.time_unit", "SECONDS"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application1FullName, "mobile.applications.%s.pairing_key_lifetime.duration", "3"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application1FullName, "mobile.applications.%s.pairing_key_lifetime.time_unit", "HOURS"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application1FullName, "mobile.applications.%s.push_limit.count", "10"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application1FullName, "mobile.applications.%s.push_limit.lock_duration.duration", "260"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application1FullName, "mobile.applications.%s.push_limit.lock_duration.time_unit", "SECONDS"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application1FullName, "mobile.applications.%s.push_limit.time_period.duration", "300"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application1FullName, "mobile.applications.%s.push_limit.time_period.time_unit", "SECONDS"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application1FullName, "mobile.applications.%s.otp.enabled", "true"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application1FullName, "mobile.applications.%s.device_authorization.enabled", "true"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application1FullName, "mobile.applications.%s.device_authorization.extra_verification", "restrictive"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application1FullName, "mobile.applications.%s.auto_enrollment.enabled", "true"),
					mfa.TestCheckMFADevicePolicyApplicationMapNoResourceAttr(resourceFullName, application1FullName, "mobile.applications.%s.integrity_detection"),

					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application2FullName, "mobile.applications.%s.pairing_disabled", "false"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application2FullName, "mobile.applications.%s.push.enabled", "false"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application2FullName, "mobile.applications.%s.push.number_matching.enabled", "false"),
					mfa.TestCheckMFADevicePolicyApplicationMapNoResourceAttr(resourceFullName, application2FullName, "mobile.applications.%s.push_timeout.duration"),
					mfa.TestCheckMFADevicePolicyApplicationMapNoResourceAttr(resourceFullName, application2FullName, "mobile.applications.%s.push_timeout.time_unit"),
					mfa.TestCheckMFADevicePolicyApplicationMapNoResourceAttr(resourceFullName, application2FullName, "mobile.applications.%s.pairing_key_lifetime.duration"),
					mfa.TestCheckMFADevicePolicyApplicationMapNoResourceAttr(resourceFullName, application2FullName, "mobile.applications.%s.pairing_key_lifetime.time_unit"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application2FullName, "mobile.applications.%s.push_limit.count", "5"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application2FullName, "mobile.applications.%s.push_limit.lock_duration.duration", "30"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application2FullName, "mobile.applications.%s.push_limit.lock_duration.time_unit", "MINUTES"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application2FullName, "mobile.applications.%s.push_limit.time_period.duration", "10"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application2FullName, "mobile.applications.%s.push_limit.time_period.time_unit", "MINUTES"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application2FullName, "mobile.applications.%s.otp.enabled", "true"),
					mfa.TestCheckMFADevicePolicyApplicationMapNoResourceAttr(resourceFullName, application2FullName, "mobile.applications.%s.device_authorization.extra_verification"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application2FullName, "mobile.applications.%s.integrity_detection", "permissive"),

					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application3FullName, "mobile.applications.%s.pairing_disabled", "false"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application3FullName, "mobile.applications.%s.push.enabled", "true"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application3FullName, "mobile.applications.%s.push.number_matching.enabled", "false"),
					mfa.TestCheckMFADevicePolicyApplicationMapNoResourceAttr(resourceFullName, application3FullName, "mobile.applications.%s.push_timeout.duration"),
					mfa.TestCheckMFADevicePolicyApplicationMapNoResourceAttr(resourceFullName, application3FullName, "mobile.applications.%s.push_timeout.time_unit"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application3FullName, "mobile.applications.%s.pairing_key_lifetime.duration", "55"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application3FullName, "mobile.applications.%s.pairing_key_lifetime.time_unit", "MINUTES"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application3FullName, "mobile.applications.%s.push_limit.count", "5"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application3FullName, "mobile.applications.%s.push_limit.lock_duration.duration", "25"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application3FullName, "mobile.applications.%s.push_limit.lock_duration.time_unit", "MINUTES"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application3FullName, "mobile.applications.%s.push_limit.time_period.duration", "5"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application3FullName, "mobile.applications.%s.push_limit.time_period.time_unit", "MINUTES"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application3FullName, "mobile.applications.%s.otp.enabled", "true"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application3FullName, "mobile.applications.%s.device_authorization.enabled", "false"),
					mfa.TestCheckMFADevicePolicyApplicationMapNoResourceAttr(resourceFullName, application3FullName, "mobile.applications.%s.device_authorization.extra_verification"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application3FullName, "mobile.applications.%s.auto_enrollment.enabled", "true"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, application3FullName, "mobile.applications.%s.integrity_detection", "permissive"),

					resource.TestCheckResourceAttr(resourceFullName, "mobile.prompt_for_nickname_on_pairing", "true"),

					resource.TestCheckResourceAttr(resourceFullName, "totp.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "SMS_THEN_EMAIL"),
				),
			},
		},
	})
}

func TestAccMFADevicePolicy_Totp_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_device_policy.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.MFADevicePolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFADevicePolicyConfig_FullTotp(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.pairing_disabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.passcode_grace_period", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.otp.failure.count", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.otp.failure.cool_down.duration", "125"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.otp.failure.cool_down.time_unit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.prompt_for_nickname_on_pairing", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.uri_parameters.%", "3"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.uri_parameters.param1", "value 1"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.uri_parameters.param2", "value 2"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.uri_parameters.param3", "value 3"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "SMS_THEN_EMAIL"),
				),
			},
			// Test importing the resource
			{
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

func TestAccMFADevicePolicy_Totp_Minimal(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_device_policy.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.MFADevicePolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFADevicePolicyConfig_MinimalTotp(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.pairing_disabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.passcode_grace_period", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.otp.failure.count", "3"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.otp.failure.cool_down.duration", "2"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.otp.failure.cool_down.time_unit", "MINUTES"),
					resource.TestCheckNoResourceAttr(resourceFullName, "totp.prompt_for_nickname_on_pairing"),
					resource.TestCheckNoResourceAttr(resourceFullName, "totp.uri_parameters"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "NONE"),
				),
			},
		},
	})
}

func TestAccMFADevicePolicy_Totp_Change(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_device_policy.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.MFADevicePolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFADevicePolicyConfig_FullTotp(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.pairing_disabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.passcode_grace_period", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.otp.failure.count", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.otp.failure.cool_down.duration", "125"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.otp.failure.cool_down.time_unit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.prompt_for_nickname_on_pairing", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.uri_parameters.%", "3"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.uri_parameters.param1", "value 1"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.uri_parameters.param2", "value 2"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.uri_parameters.param3", "value 3"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "SMS_THEN_EMAIL"),
				),
			},
			{
				Config: testAccMFADevicePolicyConfig_MinimalTotp(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.pairing_disabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.passcode_grace_period", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.otp.failure.count", "3"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.otp.failure.cool_down.duration", "2"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.otp.failure.cool_down.time_unit", "MINUTES"),
					resource.TestCheckNoResourceAttr(resourceFullName, "totp.prompt_for_nickname_on_pairing"),
					resource.TestCheckNoResourceAttr(resourceFullName, "totp.uri_parameters"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "NONE"),
				),
			},
			{
				Config: testAccMFADevicePolicyConfig_FullTotp(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.pairing_disabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.passcode_grace_period", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.otp.failure.count", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.otp.failure.cool_down.duration", "125"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.otp.failure.cool_down.time_unit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.prompt_for_nickname_on_pairing", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.uri_parameters.%", "3"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.uri_parameters.param1", "value 1"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.uri_parameters.param2", "value 2"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.uri_parameters.param3", "value 3"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "SMS_THEN_EMAIL"),
				),
			},
		},
	})
}

func TestAccMFADevicePolicy_FIDO2_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_device_policy.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.MFADevicePolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFADevicePolicyConfig_FullFIDO2(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.failure.count", "4"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.failure.cool_down.duration", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.failure.cool_down.time_unit", "MINUTES"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.pairing_disabled", "true"),
					resource.TestMatchResourceAttr(resourceFullName, "fido2.fido2_policy_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.prompt_for_nickname_on_pairing", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "SMS_THEN_EMAIL"),
				),
			},
			// Test importing the resource
			{
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

func TestAccMFADevicePolicy_FIDO2_Minimal(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_device_policy.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.MFADevicePolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFADevicePolicyConfig_MinimalFIDO2(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.failure.count", "3"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.failure.cool_down.duration", "2"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.failure.cool_down.time_unit", "MINUTES"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.pairing_disabled", "false"),
					resource.TestCheckNoResourceAttr(resourceFullName, "fido2.fido2_policy_id"),
					resource.TestCheckNoResourceAttr(resourceFullName, "fido2.prompt_for_nickname_on_pairing"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "NONE"),
				),
			},
		},
	})
}

func TestAccMFADevicePolicy_FIDO2_Change(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_device_policy.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.MFADevicePolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFADevicePolicyConfig_FullFIDO2(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.failure.count", "4"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.failure.cool_down.duration", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.failure.cool_down.time_unit", "MINUTES"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.pairing_disabled", "true"),
					resource.TestMatchResourceAttr(resourceFullName, "fido2.fido2_policy_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.prompt_for_nickname_on_pairing", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "SMS_THEN_EMAIL"),
				),
			},
			{
				Config: testAccMFADevicePolicyConfig_MinimalFIDO2(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.failure.count", "3"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.failure.cool_down.duration", "2"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.failure.cool_down.time_unit", "MINUTES"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.pairing_disabled", "false"),
					resource.TestCheckNoResourceAttr(resourceFullName, "fido2.fido2_policy_id"),
					resource.TestCheckNoResourceAttr(resourceFullName, "fido2.prompt_for_nickname_on_pairing"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "NONE"),
				),
			},
			{
				Config: testAccMFADevicePolicyConfig_FullFIDO2(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.failure.count", "4"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.failure.cool_down.duration", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.failure.cool_down.time_unit", "MINUTES"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.pairing_disabled", "true"),
					resource.TestMatchResourceAttr(resourceFullName, "fido2.fido2_policy_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.prompt_for_nickname_on_pairing", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "SMS_THEN_EMAIL"),
				),
			},
		},
	})
}

func TestAccMFADevicePolicy_FIDO2_Disabled(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_device_policy.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.MFADevicePolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFADevicePolicyConfig_FullFIDO2(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.enabled", "true"),
					resource.TestCheckResourceAttrSet(resourceFullName, "fido2.failure.count"),
					resource.TestCheckResourceAttrSet(resourceFullName, "fido2.failure.cool_down.duration"),
					resource.TestCheckResourceAttrSet(resourceFullName, "fido2.failure.cool_down.time_unit"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.pairing_disabled", "true"),
					resource.TestMatchResourceAttr(resourceFullName, "fido2.fido2_policy_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.prompt_for_nickname_on_pairing", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "SMS_THEN_EMAIL"),
				),
			},
			{
				Config: testAccMFADevicePolicyConfig_DisableFIDO2(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.enabled", "false"),
					resource.TestCheckResourceAttrSet(resourceFullName, "fido2.failure.count"),
					resource.TestCheckResourceAttrSet(resourceFullName, "fido2.failure.cool_down.duration"),
					resource.TestCheckResourceAttrSet(resourceFullName, "fido2.failure.cool_down.time_unit"),
					resource.TestCheckNoResourceAttr(resourceFullName, "fido2.pairing_disabled"),
					resource.TestCheckNoResourceAttr(resourceFullName, "fido2.fido2_policy_id"),
					resource.TestCheckNoResourceAttr(resourceFullName, "fido2.prompt_for_nickname_on_pairing"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "NONE"),
				),
			},
			{
				Config: testAccMFADevicePolicyConfig_FullFIDO2(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.enabled", "true"),
					resource.TestCheckResourceAttrSet(resourceFullName, "fido2.failure.count"),
					resource.TestCheckResourceAttrSet(resourceFullName, "fido2.failure.cool_down.duration"),
					resource.TestCheckResourceAttrSet(resourceFullName, "fido2.failure.cool_down.time_unit"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.pairing_disabled", "true"),
					resource.TestMatchResourceAttr(resourceFullName, "fido2.fido2_policy_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.prompt_for_nickname_on_pairing", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "SMS_THEN_EMAIL"),
				),
			},
		},
	})
}

func TestAccMFADevicePolicy_Desktop_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_device_policy.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckRegionSupportsWorkforce(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.MFADevicePolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFADevicePolicyConfig_FullDesktop(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "policy_type", "PING_ONE_ID"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "desktop.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "desktop.otp.failure.count", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "desktop.otp.failure.cool_down.duration", "125"),
					resource.TestCheckResourceAttr(resourceFullName, "desktop.otp.failure.cool_down.time_unit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "desktop.pairing_disabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "desktop.pairing_key_lifetime.duration", "48"),
					resource.TestCheckResourceAttr(resourceFullName, "desktop.pairing_key_lifetime.time_unit", "HOURS"),
					resource.TestCheckResourceAttr(resourceFullName, "desktop.prompt_for_nickname_on_pairing", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "yubikey.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oath_token.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "SMS_THEN_EMAIL"),
				),
			},
			// Test importing the resource
			{
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

func TestAccMFADevicePolicy_Desktop_Minimal(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_device_policy.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckRegionSupportsWorkforce(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.MFADevicePolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFADevicePolicyConfig_MinimalDesktop(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "policy_type", "PING_ONE_ID"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "desktop.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "desktop.otp.failure.count", "3"),
					resource.TestCheckResourceAttr(resourceFullName, "desktop.otp.failure.cool_down.duration", "2"),
					resource.TestCheckResourceAttr(resourceFullName, "desktop.otp.failure.cool_down.time_unit", "MINUTES"),
					resource.TestCheckResourceAttr(resourceFullName, "desktop.pairing_disabled", "false"),
					resource.TestCheckNoResourceAttr(resourceFullName, "desktop.pairing_key_lifetime"),
					resource.TestCheckNoResourceAttr(resourceFullName, "desktop.prompt_for_nickname_on_pairing"),
					resource.TestCheckResourceAttr(resourceFullName, "yubikey.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oath_token.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "NONE"),
				),
			},
		},
	})
}

func TestAccMFADevicePolicy_Desktop_Change(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_device_policy.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckRegionSupportsWorkforce(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.MFADevicePolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFADevicePolicyConfig_FullDesktop(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "policy_type", "PING_ONE_ID"),
					resource.TestCheckResourceAttr(resourceFullName, "desktop.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "desktop.otp.failure.count", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "desktop.otp.failure.cool_down.duration", "125"),
					resource.TestCheckResourceAttr(resourceFullName, "desktop.otp.failure.cool_down.time_unit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "desktop.pairing_disabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "desktop.pairing_key_lifetime.duration", "48"),
					resource.TestCheckResourceAttr(resourceFullName, "desktop.pairing_key_lifetime.time_unit", "HOURS"),
					resource.TestCheckResourceAttr(resourceFullName, "desktop.prompt_for_nickname_on_pairing", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "yubikey.enabled", "false"),
				),
			},
			{
				Config: testAccMFADevicePolicyConfig_MinimalDesktop(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "policy_type", "PING_ONE_ID"),
					resource.TestCheckResourceAttr(resourceFullName, "desktop.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "desktop.otp.failure.count", "3"),
					resource.TestCheckResourceAttr(resourceFullName, "desktop.otp.failure.cool_down.duration", "2"),
					resource.TestCheckResourceAttr(resourceFullName, "desktop.otp.failure.cool_down.time_unit", "MINUTES"),
					resource.TestCheckResourceAttr(resourceFullName, "desktop.pairing_disabled", "false"),
					resource.TestCheckNoResourceAttr(resourceFullName, "desktop.pairing_key_lifetime"),
					resource.TestCheckNoResourceAttr(resourceFullName, "desktop.prompt_for_nickname_on_pairing"),
					resource.TestCheckResourceAttr(resourceFullName, "yubikey.enabled", "false"),
				),
			},
			{
				Config: testAccMFADevicePolicyConfig_FullDesktop(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "policy_type", "PING_ONE_ID"),
					resource.TestCheckResourceAttr(resourceFullName, "desktop.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "desktop.otp.failure.count", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "desktop.otp.failure.cool_down.duration", "125"),
					resource.TestCheckResourceAttr(resourceFullName, "desktop.otp.failure.cool_down.time_unit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "desktop.pairing_disabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "desktop.pairing_key_lifetime.duration", "48"),
					resource.TestCheckResourceAttr(resourceFullName, "desktop.pairing_key_lifetime.time_unit", "HOURS"),
					resource.TestCheckResourceAttr(resourceFullName, "desktop.prompt_for_nickname_on_pairing", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "yubikey.enabled", "false"),
				),
			},
		},
	})
}

// TestAccMFADevicePolicy_Desktop_NoPolicyTypeError guards against a regression of a
// data-loss bug (CDI-1259 QA finding): PingID-only attributes (`desktop`, `yubikey`,
// and the PingID-only mobile-application sub-fields) were silently accepted and
// silently dropped - with no plan-time error - when `policy_type` was omitted from
// configuration and left to resolve via its schema default of `PING_ONE_MFA`. This
// happened because the schema's `ConflictsIfMatchesPathValue`/`IsRequiredIfMatchesPathValue`
// validators read `policy_type`'s raw config value, which is null when omitted, and
// skip their check on null. `ValidateConfig` on the resource now closes that gap by
// resolving the effective (post-default) `policy_type` value and re-running the
// conflict check. This test asserts a config that sets `desktop` without an explicit
// `policy_type` now fails at plan time with a clear error, rather than silently
// planning clean and dropping `desktop` on apply.
func TestAccMFADevicePolicy_Desktop_NoPolicyTypeError(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.MFADevicePolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccMFADevicePolicyConfig_DesktopNoPolicyTypeError(resourceName, name),
				ExpectError: regexp.MustCompile("Invalid argument combination"),
				PlanOnly:    true,
			},
		},
	})
}

func TestAccMFADevicePolicy_Yubikey_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_device_policy.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckRegionSupportsWorkforce(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.MFADevicePolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFADevicePolicyConfig_FullYubikey(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "policy_type", "PING_ONE_ID"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "desktop.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "yubikey.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "yubikey.otp.failure.count", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "yubikey.otp.failure.cool_down.duration", "125"),
					resource.TestCheckResourceAttr(resourceFullName, "yubikey.otp.failure.cool_down.time_unit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "yubikey.pairing_disabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "yubikey.pairing_key_lifetime.duration", "48"),
					resource.TestCheckResourceAttr(resourceFullName, "yubikey.pairing_key_lifetime.time_unit", "HOURS"),
					resource.TestCheckResourceAttr(resourceFullName, "yubikey.prompt_for_nickname_on_pairing", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oath_token.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "SMS_THEN_EMAIL"),
				),
			},
			// Test importing the resource
			{
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

func TestAccMFADevicePolicy_Yubikey_Minimal(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_device_policy.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckRegionSupportsWorkforce(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.MFADevicePolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFADevicePolicyConfig_MinimalYubikey(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "policy_type", "PING_ONE_ID"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "desktop.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "yubikey.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "yubikey.otp.failure.count", "3"),
					resource.TestCheckResourceAttr(resourceFullName, "yubikey.otp.failure.cool_down.duration", "2"),
					resource.TestCheckResourceAttr(resourceFullName, "yubikey.otp.failure.cool_down.time_unit", "MINUTES"),
					resource.TestCheckResourceAttr(resourceFullName, "yubikey.pairing_disabled", "false"),
					resource.TestCheckNoResourceAttr(resourceFullName, "yubikey.pairing_key_lifetime"),
					resource.TestCheckNoResourceAttr(resourceFullName, "yubikey.prompt_for_nickname_on_pairing"),
					resource.TestCheckResourceAttr(resourceFullName, "oath_token.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "NONE"),
				),
			},
		},
	})
}

func TestAccMFADevicePolicy_Yubikey_Change(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_device_policy.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckRegionSupportsWorkforce(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.MFADevicePolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFADevicePolicyConfig_FullYubikey(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "policy_type", "PING_ONE_ID"),
					resource.TestCheckResourceAttr(resourceFullName, "yubikey.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "yubikey.otp.failure.count", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "yubikey.otp.failure.cool_down.duration", "125"),
					resource.TestCheckResourceAttr(resourceFullName, "yubikey.otp.failure.cool_down.time_unit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "yubikey.pairing_disabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "yubikey.pairing_key_lifetime.duration", "48"),
					resource.TestCheckResourceAttr(resourceFullName, "yubikey.pairing_key_lifetime.time_unit", "HOURS"),
					resource.TestCheckResourceAttr(resourceFullName, "yubikey.prompt_for_nickname_on_pairing", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "desktop.enabled", "false"),
				),
			},
			{
				Config: testAccMFADevicePolicyConfig_MinimalYubikey(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "policy_type", "PING_ONE_ID"),
					resource.TestCheckResourceAttr(resourceFullName, "yubikey.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "yubikey.otp.failure.count", "3"),
					resource.TestCheckResourceAttr(resourceFullName, "yubikey.otp.failure.cool_down.duration", "2"),
					resource.TestCheckResourceAttr(resourceFullName, "yubikey.otp.failure.cool_down.time_unit", "MINUTES"),
					resource.TestCheckResourceAttr(resourceFullName, "yubikey.pairing_disabled", "false"),
					resource.TestCheckNoResourceAttr(resourceFullName, "yubikey.pairing_key_lifetime"),
					resource.TestCheckNoResourceAttr(resourceFullName, "yubikey.prompt_for_nickname_on_pairing"),
					resource.TestCheckResourceAttr(resourceFullName, "desktop.enabled", "false"),
				),
			},
			{
				Config: testAccMFADevicePolicyConfig_FullYubikey(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "policy_type", "PING_ONE_ID"),
					resource.TestCheckResourceAttr(resourceFullName, "yubikey.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "yubikey.otp.failure.count", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "yubikey.otp.failure.cool_down.duration", "125"),
					resource.TestCheckResourceAttr(resourceFullName, "yubikey.otp.failure.cool_down.time_unit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "yubikey.pairing_disabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "yubikey.pairing_key_lifetime.duration", "48"),
					resource.TestCheckResourceAttr(resourceFullName, "yubikey.pairing_key_lifetime.time_unit", "HOURS"),
					resource.TestCheckResourceAttr(resourceFullName, "yubikey.prompt_for_nickname_on_pairing", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "desktop.enabled", "false"),
				),
			},
		},
	})
}

func TestAccMFADevicePolicy_OathToken_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_device_policy.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.MFADevicePolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFADevicePolicyConfig_FullOathToken(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "policy_type", "PING_ONE_MFA"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.enabled", "false"),
					resource.TestCheckNoResourceAttr(resourceFullName, "desktop"),
					resource.TestCheckNoResourceAttr(resourceFullName, "yubikey"),
					resource.TestCheckResourceAttr(resourceFullName, "oath_token.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oath_token.otp.failure.count", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "oath_token.otp.failure.cool_down.duration", "125"),
					resource.TestCheckResourceAttr(resourceFullName, "oath_token.otp.failure.cool_down.time_unit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "oath_token.pairing_disabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oath_token.pairing_key_lifetime.duration", "48"),
					resource.TestCheckResourceAttr(resourceFullName, "oath_token.pairing_key_lifetime.time_unit", "HOURS"),
					resource.TestCheckResourceAttr(resourceFullName, "oath_token.prompt_for_nickname_on_pairing", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "SMS_THEN_EMAIL"),
				),
			},
			// Test importing the resource
			{
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

func TestAccMFADevicePolicy_OathToken_Minimal(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_device_policy.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.MFADevicePolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				// oath_token is omitted from the config entirely to exercise its
				// schema-level default (AD4: enabled defaults to false).
				Config: testAccMFADevicePolicyConfig_MinimalOathToken(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "policy_type", "PING_ONE_MFA"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.enabled", "false"),
					resource.TestCheckNoResourceAttr(resourceFullName, "desktop"),
					resource.TestCheckNoResourceAttr(resourceFullName, "yubikey"),
					resource.TestCheckResourceAttr(resourceFullName, "oath_token.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oath_token.otp.failure.count", "3"),
					resource.TestCheckResourceAttr(resourceFullName, "oath_token.otp.failure.cool_down.duration", "2"),
					resource.TestCheckResourceAttr(resourceFullName, "oath_token.otp.failure.cool_down.time_unit", "MINUTES"),
					resource.TestCheckResourceAttr(resourceFullName, "oath_token.pairing_disabled", "false"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oath_token.pairing_key_lifetime"),
					resource.TestCheckResourceAttr(resourceFullName, "oath_token.prompt_for_nickname_on_pairing", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "NONE"),
				),
			},
		},
	})
}

func TestAccMFADevicePolicy_OathToken_Change(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_device_policy.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.MFADevicePolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFADevicePolicyConfig_FullOathToken(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "policy_type", "PING_ONE_MFA"),
					resource.TestCheckResourceAttr(resourceFullName, "oath_token.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oath_token.otp.failure.count", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "oath_token.otp.failure.cool_down.duration", "125"),
					resource.TestCheckResourceAttr(resourceFullName, "oath_token.otp.failure.cool_down.time_unit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "oath_token.pairing_disabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oath_token.pairing_key_lifetime.duration", "48"),
					resource.TestCheckResourceAttr(resourceFullName, "oath_token.pairing_key_lifetime.time_unit", "HOURS"),
					resource.TestCheckResourceAttr(resourceFullName, "oath_token.prompt_for_nickname_on_pairing", "true"),
				),
			},
			{
				Config: testAccMFADevicePolicyConfig_MinimalOathToken(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "policy_type", "PING_ONE_MFA"),
					resource.TestCheckResourceAttr(resourceFullName, "oath_token.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "oath_token.otp.failure.count", "3"),
					resource.TestCheckResourceAttr(resourceFullName, "oath_token.otp.failure.cool_down.duration", "2"),
					resource.TestCheckResourceAttr(resourceFullName, "oath_token.otp.failure.cool_down.time_unit", "MINUTES"),
					resource.TestCheckResourceAttr(resourceFullName, "oath_token.pairing_disabled", "false"),
					resource.TestCheckNoResourceAttr(resourceFullName, "oath_token.pairing_key_lifetime"),
					resource.TestCheckResourceAttr(resourceFullName, "oath_token.prompt_for_nickname_on_pairing", "false"),
				),
			},
			{
				Config: testAccMFADevicePolicyConfig_FullOathToken(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "policy_type", "PING_ONE_MFA"),
					resource.TestCheckResourceAttr(resourceFullName, "oath_token.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oath_token.otp.failure.count", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "oath_token.otp.failure.cool_down.duration", "125"),
					resource.TestCheckResourceAttr(resourceFullName, "oath_token.otp.failure.cool_down.time_unit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "oath_token.pairing_disabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "oath_token.pairing_key_lifetime.duration", "48"),
					resource.TestCheckResourceAttr(resourceFullName, "oath_token.pairing_key_lifetime.time_unit", "HOURS"),
					resource.TestCheckResourceAttr(resourceFullName, "oath_token.prompt_for_nickname_on_pairing", "true"),
				),
			},
		},
	})
}

func TestAccMFADevicePolicy_Mobile_PingID_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_device_policy.%s", resourceName)

	name := resourceName

	applicationFullName := fmt.Sprintf("pingone_application.%s", name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckRegionSupportsWorkforce(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.MFADevicePolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFADevicePolicyConfig_FullMobilePingID(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "policy_type", "PING_ONE_ID"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.applications.%", "1"),

					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, applicationFullName, "mobile.applications.%s.biometrics_enabled", "true"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, applicationFullName, "mobile.applications.%s.otp.enabled", "true"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, applicationFullName, "mobile.applications.%s.new_request_duration_configuration.device_timeout.duration", "30"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, applicationFullName, "mobile.applications.%s.new_request_duration_configuration.device_timeout.time_unit", "SECONDS"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, applicationFullName, "mobile.applications.%s.new_request_duration_configuration.total_timeout.duration", "60"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, applicationFullName, "mobile.applications.%s.new_request_duration_configuration.total_timeout.time_unit", "SECONDS"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, applicationFullName, "mobile.applications.%s.ip_pairing_configuration.any_ip_address", "false"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, applicationFullName, "mobile.applications.%s.ip_pairing_configuration.only_these_ip_addresses.#", "2"),

					resource.TestCheckResourceAttr(resourceFullName, "totp.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "desktop.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "yubikey.enabled", "false"),
				),
			},
			// Test importing the resource
			{
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

func TestAccMFADevicePolicy_Mobile_PingID_AnyIPAddress(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_device_policy.%s", resourceName)

	name := resourceName

	applicationFullName := fmt.Sprintf("pingone_application.%s", name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckRegionSupportsWorkforce(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.MFADevicePolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFADevicePolicyConfig_MinimalMobilePingID(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "policy_type", "PING_ONE_ID"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.applications.%", "1"),

					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, applicationFullName, "mobile.applications.%s.new_request_duration_configuration.device_timeout.duration", "25"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, applicationFullName, "mobile.applications.%s.new_request_duration_configuration.device_timeout.time_unit", "SECONDS"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, applicationFullName, "mobile.applications.%s.new_request_duration_configuration.total_timeout.duration", "40"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, applicationFullName, "mobile.applications.%s.new_request_duration_configuration.total_timeout.time_unit", "SECONDS"),
					mfa.TestCheckMFADevicePolicyApplicationMapResourceAttr(resourceFullName, applicationFullName, "mobile.applications.%s.ip_pairing_configuration.any_ip_address", "true"),
					mfa.TestCheckMFADevicePolicyApplicationMapNoResourceAttr(resourceFullName, applicationFullName, "mobile.applications.%s.ip_pairing_configuration.only_these_ip_addresses"),
				),
			},
		},
	})
}

func TestAccMFADevicePolicy_DataModel(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_device_policy.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.MFADevicePolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Minimal from new
			{
				Config: testAccMFADevicePolicyConfig_MinimalFIDO2(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "authentication.device_selection", "DEFAULT_TO_FIRST"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "NONE"),
					resource.TestCheckResourceAttr(resourceFullName, "default", "false"),
				),
			},
			{
				Config:  testAccMFADevicePolicyConfig_MinimalFIDO2(resourceName, name),
				Destroy: true,
			},
			// Full from new
			{
				Config: testAccMFADevicePolicyConfig_FullFIDO2(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "authentication.device_selection", "PROMPT_TO_SELECT"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "SMS_THEN_EMAIL"),
					resource.TestCheckResourceAttr(resourceFullName, "default", "false"),
				),
			},
			{
				Config:  testAccMFADevicePolicyConfig_FullFIDO2(resourceName, name),
				Destroy: true,
			},
			// Update
			{
				Config: testAccMFADevicePolicyConfig_FullFIDO2(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "authentication.device_selection", "PROMPT_TO_SELECT"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "SMS_THEN_EMAIL"),
					resource.TestCheckResourceAttr(resourceFullName, "default", "false"),
				),
			},
			{
				Config: testAccMFADevicePolicyConfig_MinimalFIDO2(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "authentication.device_selection", "DEFAULT_TO_FIRST"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "NONE"),
					resource.TestCheckResourceAttr(resourceFullName, "default", "false"),
				),
			},
			{
				Config: testAccMFADevicePolicyConfig_FullFIDO2(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "authentication.device_selection", "PROMPT_TO_SELECT"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "SMS_THEN_EMAIL"),
					resource.TestCheckResourceAttr(resourceFullName, "default", "false"),
				),
			},
		},
	})
}

func TestAccMFADevicePolicy_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_device_policy.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.MFADevicePolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccMFADevicePolicyConfig_FullSMS(resourceName, name),
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

func TestAccMFADevicePolicy_DeleteDependentSOPFinalAction(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_device_policy.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.MFADevicePolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFADevicePolicyConfig_DeleteDependentSOPFinalAction(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
				),
			},
			{
				Config:  testAccMFADevicePolicyConfig_DeleteDependentSOPFinalAction(resourceName, name),
				Destroy: true,
			},
		},
	})
}

func testAccMFADevicePolicyConfig_NewEnv(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_mfa_device_policy" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s"

  sms = {
    enabled = true
  }

  voice = {
    enabled = true
  }

  email = {
    enabled = true
  }

  mobile = {
    enabled = true
  }

  totp = {
    enabled = true
  }

  fido2 = {
    enabled = true
  }

}`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccMFADevicePolicyConfig_FullSMS(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_notification_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  lifecycle {
    create_before_destroy = true
  }
}

resource "pingone_mfa_device_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

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

  authentication = {
    device_selection = "DEFAULT_TO_FIRST"
  }

  sms = {
    enabled          = true
    pairing_disabled = true

    prompt_for_nickname_on_pairing = true

    otp = {
      lifetime = {
        duration  = 75
        time_unit = "SECONDS"
      }

      failure = {
        count = 5

        cool_down = {
          duration  = 5
          time_unit = "SECONDS"
        }
      }

      otp_length = 7
    }
  }

  voice = {
    enabled = false
  }

  email = {
    enabled = false
  }

  whats_app = {
    enabled = false

    pairing_disabled = true

    prompt_for_nickname_on_pairing = true

    otp = {
      lifetime = {
        duration  = 30
        time_unit = "MINUTES"
      }

      failure = {
        count = 3

        cool_down = {
          duration  = 0
          time_unit = "MINUTES"
        }
      }

      otp_length = 7
    }
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

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccMFADevicePolicyConfig_MinimalSMS(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_mfa_device_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  sms = {
    enabled = true
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

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccMFADevicePolicyConfig_FullVoice(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_mfa_device_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  new_device_notification = "SMS_THEN_EMAIL"

  authentication = {
    device_selection = "PROMPT_TO_SELECT"
  }

  sms = {
    enabled = false
  }

  voice = {
    enabled          = true
    pairing_disabled = true

    prompt_for_nickname_on_pairing = true

    otp = {
      lifetime = {
        duration  = 75
        time_unit = "SECONDS"
      }

      failure = {
        count = 5

        cool_down = {
          duration  = 5
          time_unit = "SECONDS"
        }
      }

      otp_length = 7
    }
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

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccMFADevicePolicyConfig_MinimalVoice(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_mfa_device_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  sms = {
    enabled = false
  }

  voice = {
    enabled = true
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

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccMFADevicePolicyConfig_FullEmail(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_mfa_device_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  new_device_notification = "SMS_THEN_EMAIL"

  authentication = {
    device_selection = "ALWAYS_DISPLAY_DEVICES"
  }

  sms = {
    enabled = false
  }

  voice = {
    enabled = false
  }

  email = {
    enabled          = true
    pairing_disabled = true

    prompt_for_nickname_on_pairing = true

    otp = {
      lifetime = {
        duration  = 75
        time_unit = "SECONDS"
      }

      failure = {
        count = 5

        cool_down = {
          duration  = 5
          time_unit = "SECONDS"
        }
      }

      otp_length = 7
    }
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

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccMFADevicePolicyConfig_MinimalEmail(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_mfa_device_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  sms = {
    enabled = false
  }

  voice = {
    enabled = false
  }

  email = {
    enabled = true
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

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccMFADevicePolicyConfig_FullMobile(resourceName, name, key string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[2]s-1" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-1"
  description    = "My test OIDC app for MFA Policy"

  login_page_url = "https://www.pingidentity.com"

  enabled = true

  oidc_options = {
    type                       = "NATIVE_APP"
    grant_types                = ["CLIENT_CREDENTIALS"]
    token_endpoint_auth_method = "CLIENT_SECRET_BASIC"

    mobile_app = {
      bundle_id    = "com.%[2]s1.bundle"
      package_name = "com.%[2]s1.package"

      passcode_refresh_seconds = 45

      integrity_detection = {
        enabled = false
      }
    }
  }
}

resource "pingone_mfa_application_push_credential" "%[2]s-1" {
  environment_id = data.pingone_environment.general_test.id
  application_id = pingone_application.%[2]s-1.id

  fcm = {
    google_service_account_credentials = jsonencode(%[4]s)
  }
}

resource "pingone_application" "%[2]s-2" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-2"
  description    = "My test OIDC app for MFA Policy"

  login_page_url = "https://www.pingidentity.com"

  enabled = true

  oidc_options = {
    type                       = "NATIVE_APP"
    grant_types                = ["CLIENT_CREDENTIALS"]
    token_endpoint_auth_method = "CLIENT_SECRET_BASIC"

    mobile_app = {
      bundle_id    = "com.%[2]s2.bundle"
      package_name = "com.%[2]s2.package"

      passcode_refresh_seconds = 45

      integrity_detection = {
        enabled = true
        cache_duration = {
          amount = 30
          units  = "HOURS"
        }
        google_play = {
          verification_type = "INTERNAL"
          decryption_key    = "dummykeydoesnotexist"
          verification_key  = "dummykeydoesnotexist"
        }
      }
    }
  }
}

resource "pingone_application" "%[2]s-3" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-3"
  description    = "My test OIDC app for MFA Policy"

  login_page_url = "https://www.pingidentity.com"

  enabled = true

  oidc_options = {
    type                       = "NATIVE_APP"
    grant_types                = ["CLIENT_CREDENTIALS"]
    token_endpoint_auth_method = "CLIENT_SECRET_BASIC"

    mobile_app = {
      bundle_id    = "com.%[2]s3.bundle"
      package_name = "com.%[2]s3.package"

      passcode_refresh_seconds = 45

      integrity_detection = {
        enabled = true
        cache_duration = {
          amount = 30
          units  = "HOURS"
        }
        google_play = {
          verification_type = "INTERNAL"
          decryption_key    = "dummykeydoesnotexist"
          verification_key  = "dummykeydoesnotexist"
        }
      }
    }
  }
}

resource "pingone_mfa_application_push_credential" "%[2]s-3" {
  environment_id = data.pingone_environment.general_test.id
  application_id = pingone_application.%[2]s-3.id

  fcm = {
    google_service_account_credentials = jsonencode(%[4]s)
  }
}

resource "pingone_application" "%[2]s-4" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-4"
  description    = "My test OIDC app for MFA Policy"

  login_page_url = "https://www.pingidentity.com"

  enabled = true

  oidc_options = {
    type                       = "NATIVE_APP"
    grant_types                = ["CLIENT_CREDENTIALS"]
    token_endpoint_auth_method = "CLIENT_SECRET_BASIC"

    mobile_app = {
      bundle_id    = "com.%[2]s4.bundle"
      package_name = "com.%[2]s4.package"

      passcode_refresh_seconds = 45

      integrity_detection = {
        enabled = true
        cache_duration = {
          amount = 30
          units  = "HOURS"
        }
        google_play = {
          verification_type = "INTERNAL"
          decryption_key    = "dummykeydoesnotexist"
          verification_key  = "dummykeydoesnotexist"
        }
      }
    }
  }
}

resource "pingone_mfa_application_push_credential" "%[2]s-4" {
  environment_id = data.pingone_environment.general_test.id
  application_id = pingone_application.%[2]s-4.id

  fcm = {
    google_service_account_credentials = jsonencode(%[4]s)
  }
}

resource "pingone_mfa_device_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  new_device_notification = "SMS_THEN_EMAIL"

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

    prompt_for_nickname_on_pairing = true

    otp = {
      failure = {
        count = 5

        cool_down = {
          duration  = 125
          time_unit = "SECONDS"
        }
      }
    }

    applications = {
      (pingone_application.%[2]s-1.id) = {

        pairing_disabled = true

        push = {
          enabled = true
          number_matching = {
            enabled = true
          }
        }

        push_timeout = {
          duration = 100
        }

        pairing_key_lifetime = {
          duration  = 3
          time_unit = "HOURS"
        }

        push_limit = {
          count = 10

          lock_duration = {
            duration  = 260
            time_unit = "SECONDS"
          }

          time_period = {
            duration  = 300
            time_unit = "SECONDS"
          }
        }

        otp = {
          enabled = true
        }

        device_authorization = {
          enabled            = true
          extra_verification = "restrictive"
        }

        auto_enrollment = {
          enabled = true
        }
      },
      (pingone_application.%[2]s-2.id) = {

        pairing_disabled = false

        push = {
          enabled = false
        }

        otp = {
          enabled = true
        }

        integrity_detection = "permissive"
      },
      (pingone_application.%[2]s-3.id) = {

        push = {
          enabled = true
        }

        otp = {
          enabled = true
        }

        pairing_key_lifetime = {
          duration  = 55
          time_unit = "MINUTES"
        }

        push_limit = {
          lock_duration = {
            duration  = 25
            time_unit = "MINUTES"
          }
          time_period = {
            duration  = 5
            time_unit = "MINUTES"
          }
        }

        device_authorization = {
          enabled = false
        }

        auto_enrollment = {
          enabled = true
        }

        integrity_detection = "permissive"
      }
    }
  }

  totp = {
    enabled = false
  }

  fido2 = {
    enabled = false
  }

  depends_on = [
    pingone_mfa_application_push_credential.%[2]s-1,
    pingone_mfa_application_push_credential.%[2]s-3
  ]

}`, acctest.GenericSandboxEnvironment(), resourceName, name, key)
}

func testAccMFADevicePolicyConfig_MobileIntegrityDetectionError_1(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "My test OIDC app for MFA Policy"

  login_page_url = "https://www.pingidentity.com"

  enabled = true

  oidc_options = {
    type                       = "NATIVE_APP"
    grant_types                = ["CLIENT_CREDENTIALS"]
    token_endpoint_auth_method = "CLIENT_SECRET_BASIC"

    mobile_app = {
      bundle_id    = "com.%[2]s.bundle"
      package_name = "com.%[2]s.package"

      passcode_refresh_seconds = 45

      integrity_detection = {
        enabled = false
      }
    }
  }
}

resource "pingone_mfa_device_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

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

    otp = {
      failure = {
        count = 5

        cool_down = {
          duration  = 125
          time_unit = "SECONDS"
        }
      }
    }

    applications = {
      (pingone_application.%[2]s.id) = {

        push = {
          enabled = true
        }

        otp = {
          enabled = true
        }

        device_authorization = {
          enabled            = true
          extra_verification = "restrictive"
        }

        auto_enrollment = {
          enabled = true
        }

        integrity_detection = "permissive"
      }
    }
  }

  totp = {
    enabled = false
  }

  fido2 = {
    enabled = false
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccMFADevicePolicyConfig_MobileIntegrityDetectionError_2(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "My test OIDC app for MFA Policy"

  login_page_url = "https://www.pingidentity.com"

  enabled = true

  oidc_options = {
    type                       = "NATIVE_APP"
    grant_types                = ["CLIENT_CREDENTIALS"]
    token_endpoint_auth_method = "CLIENT_SECRET_BASIC"

    mobile_app = {
      bundle_id    = "com.%[2]s.bundle"
      package_name = "com.%[2]s.package"

      passcode_refresh_seconds = 45

      integrity_detection = {
        enabled = true
        cache_duration = {
          amount = 30
          units  = "HOURS"
        }
        google_play = {
          verification_type = "INTERNAL"
          decryption_key    = "dummykeydoesnotexist"
          verification_key  = "dummykeydoesnotexist"
        }
      }
    }
  }
}

resource "pingone_mfa_device_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

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

    otp = {
      failure = {
        count = 5

        cool_down = {
          duration  = 125
          time_unit = "SECONDS"
        }
      }
    }

    applications = {
      (pingone_application.%[2]s.id) = {

        push = {
          enabled = true
        }

        otp = {
          enabled = true
        }

        device_authorization = {
          enabled            = true
          extra_verification = "restrictive"
        }

        auto_enrollment = {
          enabled = true
        }
      }
    }
  }

  totp = {
    enabled = false
  }

  fido2 = {
    enabled = false
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccMFADevicePolicyConfig_MobileBadApplicationError_1(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_mfa_device_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

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

    otp = {
      failure = {
        count = 5

        cool_down = {
          duration  = 125
          time_unit = "SECONDS"
        }
      }
    }

    applications = {
      "9bf6c075-78ba-4cd6-a5b1-96ec144d66ef" = { // Fake ID

        push = {
          enabled = true
        }

        otp = {
          enabled = true
        }

        device_authorization = {
          enabled            = true
          extra_verification = "restrictive"
        }

        auto_enrollment = {
          enabled = true
        }
      }
    }
  }

  totp = {
    enabled = false
  }

  fido2 = {
    enabled = false
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccMFADevicePolicyConfig_MobileBadApplicationError_2(resourceName, name string) string {
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
  description    = "My test SAML app for MFA Policy"
  login_page_url = "https://www.pingidentity.com"

  enabled = true

  saml_options = {
    acs_urls           = ["https://my-saas-app.com"]
    assertion_duration = 3600
    sp_entity_id       = "sp:%[2]s:localhost"

    idp_signing_key = {
      key_id    = pingone_key.%[3]s.id
      algorithm = pingone_key.%[3]s.signature_algorithm
    }
  }
}

resource "pingone_mfa_device_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

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

    otp = {
      failure = {
        count = 5

        cool_down = {
          duration  = 125
          time_unit = "SECONDS"
        }
      }
    }

    applications = {
      (pingone_application.%[2]s.id) = {

        push = {
          enabled = true
        }

        otp = {
          enabled = true
        }

        device_authorization = {
          enabled            = true
          extra_verification = "restrictive"
        }

        auto_enrollment = {
          enabled = true
        }
      }
    }
  }

  totp = {
    enabled = false
  }

  fido2 = {
    enabled = false
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccMFADevicePolicyConfig_MobileBadApplicationError_3(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "My test OIDC app for MFA Policy"

  login_page_url = "https://www.pingidentity.com"

  enabled = true

  oidc_options = {
    type                       = "WORKER"
    grant_types                = ["CLIENT_CREDENTIALS"]
    token_endpoint_auth_method = "CLIENT_SECRET_BASIC"
  }
}

resource "pingone_mfa_device_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

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

    otp = {
      failure = {
        count = 5

        cool_down = {
          duration  = 125
          time_unit = "SECONDS"
        }
      }
    }

    applications = {
      (pingone_application.%[2]s.id) = {

        push = {
          enabled = true
        }

        otp = {
          enabled = true
        }

        device_authorization = {
          enabled            = true
          extra_verification = "restrictive"
        }

        auto_enrollment = {
          enabled = true
        }
      }
    }
  }

  totp = {
    enabled = false
  }

  fido2 = {
    enabled = false
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccMFADevicePolicyConfig_MobileBadApplicationError_4(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "My test OIDC app for MFA Policy"

  login_page_url = "https://www.pingidentity.com"

  enabled = true

  oidc_options = {
    type                       = "NATIVE_APP"
    grant_types                = ["CLIENT_CREDENTIALS"]
    token_endpoint_auth_method = "CLIENT_SECRET_BASIC"
  }
}

resource "pingone_mfa_device_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

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

    otp = {
      failure = {
        count = 5

        cool_down = {
          duration  = 125
          time_unit = "SECONDS"
        }
      }
    }

    applications = {
      (pingone_application.%[2]s.id) = {

        push = {
          enabled = false
        }

        otp = {
          enabled = true
        }

        device_authorization = {
          enabled = false
        }

        auto_enrollment = {
          enabled = true
        }
      }
    }
  }

  totp = {
    enabled = false
  }

  fido2 = {
    enabled = false
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccMFADevicePolicyConfig_MinimalMobile(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_mfa_device_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

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

  fido2 = {
    enabled = false
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccMFADevicePolicyConfig_FullTotp(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_mfa_device_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  new_device_notification = "SMS_THEN_EMAIL"

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
    enabled               = true
    pairing_disabled      = true
    passcode_grace_period = 5

    prompt_for_nickname_on_pairing = true

    uri_parameters = {
      "param1" = "value 1"
      "param3" = "value 3"
      "param2" = "value 2"
    }

    otp = {
      failure = {
        count = 5

        cool_down = {
          duration  = 125
          time_unit = "SECONDS"
        }
      }
    }
  }

  fido2 = {
    enabled = false
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccMFADevicePolicyConfig_MinimalTotp(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_mfa_device_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

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
    enabled = true
  }

  fido2 = {
    enabled = false
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccMFADevicePolicyConfig_FullFIDO2(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_mfa_fido2_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  attestation_requirements = "NONE"
  authenticator_attachment = "PLATFORM"

  backup_eligibility = {
    allow                         = false
    enforce_during_authentication = true
  }

  device_display_name = "fidoPolicy.deviceDisplayName02"

  discoverable_credentials = "DISCOURAGED"

  mds_authenticators_requirements = {
    enforce_during_authentication = false
    option                        = "NONE"
  }

  relying_party_id = "ping-devops.com"

  user_display_name_attributes = {
    attributes = [
      {
        name = "username"
      }
    ]
  }

  user_verification = {
    enforce_during_authentication = false
    option                        = "DISCOURAGED"
  }
}

resource "pingone_mfa_device_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  authentication = {
    device_selection = "PROMPT_TO_SELECT"
  }

  new_device_notification = "SMS_THEN_EMAIL"

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
    enabled          = true
    pairing_disabled = true

    failure = {
      count = 4
      cool_down = {
        duration  = 5
        time_unit = "MINUTES"
      }
    }

    prompt_for_nickname_on_pairing = true

    fido2_policy_id = pingone_mfa_fido2_policy.%[2]s.id
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccMFADevicePolicyConfig_MinimalFIDO2(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_mfa_fido2_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  attestation_requirements = "NONE"
  authenticator_attachment = "PLATFORM"

  backup_eligibility = {
    allow                         = false
    enforce_during_authentication = true
  }

  device_display_name = "fidoPolicy.deviceDisplayName02"

  discoverable_credentials = "DISCOURAGED"

  mds_authenticators_requirements = {
    enforce_during_authentication = false
    option                        = "NONE"
  }

  relying_party_id = "ping-devops.com"

  user_display_name_attributes = {
    attributes = [
      {
        name = "username"
      }
    ]
  }

  user_verification = {
    enforce_during_authentication = false
    option                        = "DISCOURAGED"
  }
}

resource "pingone_mfa_device_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

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
    enabled = true
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccMFADevicePolicyConfig_FullDesktop(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_mfa_device_policy" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  policy_type    = "PING_ONE_ID"

  name = "%[3]s"

  new_device_notification = "SMS_THEN_EMAIL"

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

    otp = {
      failure = {
        count = 5

        cool_down = {
          duration  = 125
          time_unit = "SECONDS"
        }
      }
    }

    pairing_disabled = true

    pairing_key_lifetime = {
      duration  = 48
      time_unit = "HOURS"
    }

    prompt_for_nickname_on_pairing = true
  }

  yubikey = {
    enabled = false
  }

}`, acctest.WorkforceV2SandboxEnvironment(), resourceName, name)
}

// testAccMFADevicePolicyConfig_DesktopNoPolicyTypeError deliberately omits `policy_type`
// while setting `desktop` - the exact CDI-1259 QA-reported bypass scenario. It uses the
// generic (non-Workforce) sandbox environment because the expected failure is a plan-time
// config validation error raised before any PingID-specific entitlement would be needed.
func testAccMFADevicePolicyConfig_DesktopNoPolicyTypeError(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_mfa_device_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

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

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccMFADevicePolicyConfig_MinimalDesktop(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_mfa_device_policy" "%[2]s" {
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

  fido2 = {
    enabled = false
  }

  desktop = {
    enabled = true
  }

  yubikey = {
    enabled = false
  }

}`, acctest.WorkforceV2SandboxEnvironment(), resourceName, name)
}

func testAccMFADevicePolicyConfig_FullYubikey(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_mfa_device_policy" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  policy_type    = "PING_ONE_ID"

  name = "%[3]s"

  new_device_notification = "SMS_THEN_EMAIL"

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
    enabled = false
  }

  yubikey = {
    enabled = true

    otp = {
      failure = {
        count = 5

        cool_down = {
          duration  = 125
          time_unit = "SECONDS"
        }
      }
    }

    pairing_disabled = true

    pairing_key_lifetime = {
      duration  = 48
      time_unit = "HOURS"
    }

    prompt_for_nickname_on_pairing = true
  }

}`, acctest.WorkforceV2SandboxEnvironment(), resourceName, name)
}

func testAccMFADevicePolicyConfig_MinimalYubikey(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_mfa_device_policy" "%[2]s" {
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

  fido2 = {
    enabled = false
  }

  desktop = {
    enabled = false
  }

  yubikey = {
    enabled = true
  }

}`, acctest.WorkforceV2SandboxEnvironment(), resourceName, name)
}

func testAccMFADevicePolicyConfig_FullOathToken(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_mfa_device_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  new_device_notification = "SMS_THEN_EMAIL"

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

  oath_token = {
    enabled = true

    otp = {
      failure = {
        count = 5

        cool_down = {
          duration  = 125
          time_unit = "SECONDS"
        }
      }
    }

    pairing_disabled = true

    pairing_key_lifetime = {
      duration  = 48
      time_unit = "HOURS"
    }

    prompt_for_nickname_on_pairing = true
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccMFADevicePolicyConfig_MinimalOathToken(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_mfa_device_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

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

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccMFADevicePolicyConfig_FullMobilePingID(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  name           = "%[3]s"
  description    = "My test OIDC app for MFA Policy"

  login_page_url = "https://www.pingidentity.com"

  enabled = true

  oidc_options = {
    type                       = "NATIVE_APP"
    grant_types                = ["CLIENT_CREDENTIALS"]
    token_endpoint_auth_method = "CLIENT_SECRET_BASIC"

    mobile_app = {
      bundle_id    = "com.%[2]s.bundle"
      package_name = "com.%[2]s.package"

      passcode_refresh_seconds = 45

      integrity_detection = {
        enabled = false
      }
    }
  }
}

resource "pingone_mfa_device_policy" "%[2]s" {
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

    applications = {
      (pingone_application.%[2]s.id) = {

        biometrics_enabled = true

        otp = {
          enabled = true
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
          any_ip_address          = false
          only_these_ip_addresses = ["192.168.1.1/32", "10.0.0.0/8"]
        }
      }
    }
  }

  totp = {
    enabled = false
  }

  fido2 = {
    enabled = false
  }

  desktop = {
    enabled = false
  }

  yubikey = {
    enabled = false
  }

}`, acctest.WorkforceV2SandboxEnvironment(), resourceName, name)
}

func testAccMFADevicePolicyConfig_MinimalMobilePingID(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  name           = "%[3]s"
  description    = "My test OIDC app for MFA Policy"

  login_page_url = "https://www.pingidentity.com"

  enabled = true

  oidc_options = {
    type                       = "NATIVE_APP"
    grant_types                = ["CLIENT_CREDENTIALS"]
    token_endpoint_auth_method = "CLIENT_SECRET_BASIC"

    mobile_app = {
      bundle_id    = "com.%[2]s.bundle"
      package_name = "com.%[2]s.package"

      passcode_refresh_seconds = 45

      integrity_detection = {
        enabled = false
      }
    }
  }
}

resource "pingone_mfa_device_policy" "%[2]s" {
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

    applications = {
      (pingone_application.%[2]s.id) = {

        otp = {
          enabled = true
        }

        new_request_duration_configuration = {
          device_timeout = {
            duration = 25
          }
          total_timeout = {
            duration = 40
          }
        }

        ip_pairing_configuration = {
          any_ip_address = true
        }
      }
    }
  }

  totp = {
    enabled = false
  }

  fido2 = {
    enabled = false
  }

  desktop = {
    enabled = false
  }

  yubikey = {
    enabled = false
  }

}`, acctest.WorkforceV2SandboxEnvironment(), resourceName, name)
}

func testAccMFADevicePolicyConfig_DisableFIDO2(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_mfa_fido2_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  attestation_requirements = "NONE"
  authenticator_attachment = "PLATFORM"

  backup_eligibility = {
    allow                         = false
    enforce_during_authentication = true
  }

  device_display_name = "fidoPolicy.deviceDisplayName02"

  discoverable_credentials = "DISCOURAGED"

  mds_authenticators_requirements = {
    enforce_during_authentication = false
    option                        = "NONE"
  }

  relying_party_id = "ping-devops.com"

  user_display_name_attributes = {
    attributes = [
      {
        name = "username"
      }
    ]
  }

  user_verification = {
    enforce_during_authentication = false
    option                        = "DISCOURAGED"
  }
}

resource "pingone_mfa_device_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

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
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccMFADevicePolicyConfig_DeleteDependentSOPFinalAction(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_sign_on_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
}

resource "pingone_mfa_device_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

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

}

resource "pingone_sign_on_policy_action" "%[2]s" {
  environment_id    = data.pingone_environment.general_test.id
  sign_on_policy_id = pingone_sign_on_policy.%[2]s.id

  priority = "1"

  mfa {
    device_sign_on_policy_id = pingone_mfa_device_policy.%[2]s.id
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
