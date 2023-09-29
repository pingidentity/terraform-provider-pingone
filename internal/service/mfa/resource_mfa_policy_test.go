package mfa_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccCheckMFAPolicyDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.MFAAPIClient

	apiClientManagement := p1Client.API.ManagementAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_mfa_policy" {
			continue
		}

		_, rEnv, err := apiClientManagement.EnvironmentsApi.ReadOneEnvironment(ctx, rs.Primary.Attributes["environment_id"]).Execute()

		if err != nil {

			if rEnv == nil {
				return fmt.Errorf("Response object does not exist and no error detected")
			}

			if rEnv.StatusCode == 404 {
				continue
			}

			return err
		}

		body, r, err := apiClient.DeviceAuthenticationPolicyApi.ReadOneDeviceAuthenticationPolicy(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.ID).Execute()

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

		return fmt.Errorf("PingOne MFA Policy Instance %s still exists", rs.Primary.ID)
	}

	return nil
}

func TestAccGetMFAPolicyIDs(resourceName string, environmentID, resourceID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource Not found: %s", resourceName)
		}

		*resourceID = rs.Primary.ID
		*environmentID = rs.Primary.Attributes["environment_id"]

		return nil
	}
}

func TestAccMFAPolicy_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_policy.%s", resourceName)

	name := resourceName

	var resourceID, environmentID string

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.TestAccCheckMFAPolicyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccMFAPolicyConfig_FullSMS(resourceName, name),
				Check:  mfa.TestAccGetMFAPolicyIDs(resourceFullName, &environmentID, &resourceID),
			},
			// Replan after removal preconfig
			{
				PreConfig: func() {
					var ctx = context.Background()
					p1Client, err := acctest.TestClient(ctx)

					if err != nil {
						t.Fatalf("Failed to get API client: %v", err)
					}

					apiClient := p1Client.API.MFAAPIClient

					if environmentID == "" || resourceID == "" {
						t.Fatalf("One of environment ID or resource ID cannot be determined. Environment ID: %s, Resource ID: %s", environmentID, resourceID)
					}

					_, err = apiClient.DeviceAuthenticationPolicyApi.DeleteDeviceAuthenticationPolicy(ctx, environmentID, resourceID).Execute()
					if err != nil {
						t.Fatalf("Failed to delete MFA Policy: %v", err)
					}
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccMFAPolicy_NewEnv(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_policy.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.TestAccCheckMFAPolicyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFAPolicyConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
				),
			},
		},
	})
}

func TestAccMFAPolicy_SMS_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_policy.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.TestAccCheckMFAPolicyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFAPolicyConfig_FullSMS(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.pairing_disabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.otp_lifetime_duration", "75"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.otp_lifetime_timeunit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.otp_failure_count", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.otp_failure_cooldown_duration", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.otp_failure_cooldown_timeunit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "security_key.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "platform.#", "0"),
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

func TestAccMFAPolicy_SMS_Minimal(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_policy.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.TestAccCheckMFAPolicyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFAPolicyConfig_MinimalSMS(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.pairing_disabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.otp_lifetime_duration", "30"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.otp_lifetime_timeunit", "MINUTES"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.otp_failure_count", "3"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.otp_failure_cooldown_duration", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.otp_failure_cooldown_timeunit", "MINUTES"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "security_key.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "platform.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "NONE"),
				),
			},
		},
	})
}

func TestAccMFAPolicy_SMS_Change(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_policy.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.TestAccCheckMFAPolicyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFAPolicyConfig_FullSMS(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.pairing_disabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.otp_lifetime_duration", "75"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.otp_lifetime_timeunit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.otp_failure_count", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.otp_failure_cooldown_duration", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.otp_failure_cooldown_timeunit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "security_key.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "platform.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "SMS_THEN_EMAIL"),
				),
			},
			{
				Config: testAccMFAPolicyConfig_MinimalSMS(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.pairing_disabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.otp_lifetime_duration", "30"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.otp_lifetime_timeunit", "MINUTES"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.otp_failure_count", "3"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.otp_failure_cooldown_duration", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.otp_failure_cooldown_timeunit", "MINUTES"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "security_key.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "platform.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "NONE"),
				),
			},
			{
				Config: testAccMFAPolicyConfig_FullSMS(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.pairing_disabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.otp_lifetime_duration", "75"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.otp_lifetime_timeunit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.otp_failure_count", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.otp_failure_cooldown_duration", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.otp_failure_cooldown_timeunit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "security_key.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "platform.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "SMS_THEN_EMAIL"),
				),
			},
		},
	})
}

func TestAccMFAPolicy_Voice_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_policy.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.TestAccCheckMFAPolicyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFAPolicyConfig_FullVoice(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.pairing_disabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.otp_lifetime_duration", "75"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.otp_lifetime_timeunit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.otp_failure_count", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.otp_failure_cooldown_duration", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.otp_failure_cooldown_timeunit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "security_key.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "platform.#", "0"),
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

func TestAccMFAPolicy_Voice_Minimal(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_policy.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.TestAccCheckMFAPolicyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFAPolicyConfig_MinimalVoice(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.pairing_disabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.otp_lifetime_duration", "30"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.otp_lifetime_timeunit", "MINUTES"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.otp_failure_count", "3"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.otp_failure_cooldown_duration", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.otp_failure_cooldown_timeunit", "MINUTES"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "security_key.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "platform.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "NONE"),
				),
			},
		},
	})
}

func TestAccMFAPolicy_Voice_Change(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_policy.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.TestAccCheckMFAPolicyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFAPolicyConfig_FullVoice(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.pairing_disabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.otp_lifetime_duration", "75"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.otp_lifetime_timeunit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.otp_failure_count", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.otp_failure_cooldown_duration", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.otp_failure_cooldown_timeunit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "security_key.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "platform.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "SMS_THEN_EMAIL"),
				),
			},
			{
				Config: testAccMFAPolicyConfig_MinimalVoice(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.pairing_disabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.otp_lifetime_duration", "30"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.otp_lifetime_timeunit", "MINUTES"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.otp_failure_count", "3"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.otp_failure_cooldown_duration", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.otp_failure_cooldown_timeunit", "MINUTES"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "security_key.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "platform.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "NONE"),
				),
			},
			{
				Config: testAccMFAPolicyConfig_FullVoice(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.pairing_disabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.otp_lifetime_duration", "75"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.otp_lifetime_timeunit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.otp_failure_count", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.otp_failure_cooldown_duration", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.otp_failure_cooldown_timeunit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "security_key.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "platform.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "SMS_THEN_EMAIL"),
				),
			},
		},
	})
}

func TestAccMFAPolicy_Email_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_policy.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.TestAccCheckMFAPolicyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFAPolicyConfig_FullEmail(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.pairing_disabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.otp_lifetime_duration", "75"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.otp_lifetime_timeunit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.otp_failure_count", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.otp_failure_cooldown_duration", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.otp_failure_cooldown_timeunit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "security_key.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "platform.#", "0"),
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

func TestAccMFAPolicy_Email_Minimal(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_policy.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.TestAccCheckMFAPolicyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFAPolicyConfig_MinimalEmail(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.pairing_disabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.otp_lifetime_duration", "30"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.otp_lifetime_timeunit", "MINUTES"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.otp_failure_count", "3"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.otp_failure_cooldown_duration", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.otp_failure_cooldown_timeunit", "MINUTES"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "security_key.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "platform.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "NONE"),
				),
			},
		},
	})
}

func TestAccMFAPolicy_Email_Change(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_policy.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.TestAccCheckMFAPolicyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFAPolicyConfig_FullEmail(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.pairing_disabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.otp_lifetime_duration", "75"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.otp_lifetime_timeunit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.otp_failure_count", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.otp_failure_cooldown_duration", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.otp_failure_cooldown_timeunit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "security_key.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "platform.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "SMS_THEN_EMAIL"),
				),
			},
			{
				Config: testAccMFAPolicyConfig_MinimalEmail(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.pairing_disabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.otp_lifetime_duration", "30"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.otp_lifetime_timeunit", "MINUTES"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.otp_failure_count", "3"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.otp_failure_cooldown_duration", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.otp_failure_cooldown_timeunit", "MINUTES"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "security_key.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "platform.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "NONE"),
				),
			},
			{
				Config: testAccMFAPolicyConfig_FullEmail(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.pairing_disabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.otp_lifetime_duration", "75"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.otp_lifetime_timeunit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.otp_failure_count", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.otp_failure_cooldown_duration", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.otp_failure_cooldown_timeunit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "security_key.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "platform.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "SMS_THEN_EMAIL"),
				),
			},
		},
	})
}

func TestAccMFAPolicy_Mobile_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_policy.%s", resourceName)

	name := resourceName

	firebaseCredentials := os.Getenv("PINGONE_GOOGLE_FIREBASE_CREDENTIALS")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckGoogleFirebaseCredentials(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.TestAccCheckMFAPolicyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFAPolicyConfig_FullMobile(resourceName, name, firebaseCredentials),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.otp_failure_count", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.otp_failure_cooldown_duration", "125"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.otp_failure_cooldown_timeunit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.application.#", "3"),
					resource.TestMatchTypeSetElemNestedAttrs(resourceFullName, "mobile.0.application.*", map[string]*regexp.Regexp{
						"id":                                      verify.P1ResourceIDRegexpFullString,
						"pairing_disabled":                        regexp.MustCompile(`^true$`),
						"push_enabled":                            regexp.MustCompile(`^true$`),
						"push_timeout_duration":                   regexp.MustCompile(`^100$`),
						"push_timeout_timeunit":                   regexp.MustCompile(`^SECONDS$`),
						"pairing_key_lifetime_duration":           regexp.MustCompile(`^3$`),
						"pairing_key_lifetime_timeunit":           regexp.MustCompile(`^HOURS$`),
						"push_limit_count":                        regexp.MustCompile(`^10$`),
						"push_limit_lock_duration":                regexp.MustCompile(`^260$`),
						"push_limit_lock_duration_timeunit":       regexp.MustCompile(`^SECONDS$`),
						"push_limit_time_period_duration":         regexp.MustCompile(`^300$`),
						"push_limit_time_period_timeunit":         regexp.MustCompile(`^SECONDS$`),
						"otp_enabled":                             regexp.MustCompile(`^true$`),
						"device_authorization_enabled":            regexp.MustCompile(`^true$`),
						"device_authorization_extra_verification": regexp.MustCompile(`^restrictive$`),
						"auto_enrollment_enabled":                 regexp.MustCompile(`^true$`),
						"integrity_detection":                     regexp.MustCompile(`^$`),
					}),
					resource.TestMatchTypeSetElemNestedAttrs(resourceFullName, "mobile.0.application.*", map[string]*regexp.Regexp{
						"id":                                      verify.P1ResourceIDRegexpFullString,
						"pairing_disabled":                        regexp.MustCompile(`^false$`),
						"push_enabled":                            regexp.MustCompile(`^false$`),
						"push_timeout_duration":                   regexp.MustCompile(`^40$`),
						"push_timeout_timeunit":                   regexp.MustCompile(`^SECONDS$`),
						"pairing_key_lifetime_duration":           regexp.MustCompile(`^10$`),
						"pairing_key_lifetime_timeunit":           regexp.MustCompile(`^MINUTES$`),
						"push_limit_count":                        regexp.MustCompile(`^5$`),
						"push_limit_lock_duration":                regexp.MustCompile(`^30$`),
						"push_limit_lock_duration_timeunit":       regexp.MustCompile(`^MINUTES$`),
						"push_limit_time_period_duration":         regexp.MustCompile(`^10$`),
						"push_limit_time_period_timeunit":         regexp.MustCompile(`^MINUTES$`),
						"otp_enabled":                             regexp.MustCompile(`^true$`),
						"device_authorization_extra_verification": regexp.MustCompile(`^$`),
						"integrity_detection":                     regexp.MustCompile(`^permissive$`),
					}),
					resource.TestMatchTypeSetElemNestedAttrs(resourceFullName, "mobile.0.application.*", map[string]*regexp.Regexp{
						"id":                                      verify.P1ResourceIDRegexpFullString,
						"pairing_disabled":                        regexp.MustCompile(`^false$`),
						"push_enabled":                            regexp.MustCompile(`^true$`),
						"push_timeout_duration":                   regexp.MustCompile(`^40$`),
						"push_timeout_timeunit":                   regexp.MustCompile(`^SECONDS$`),
						"pairing_key_lifetime_duration":           regexp.MustCompile(`^55$`),
						"pairing_key_lifetime_timeunit":           regexp.MustCompile(`^MINUTES$`),
						"push_limit_count":                        regexp.MustCompile(`^5$`),
						"push_limit_lock_duration":                regexp.MustCompile(`^25$`),
						"push_limit_lock_duration_timeunit":       regexp.MustCompile(`^MINUTES$`),
						"push_limit_time_period_duration":         regexp.MustCompile(`^5$`),
						"push_limit_time_period_timeunit":         regexp.MustCompile(`^MINUTES$`),
						"otp_enabled":                             regexp.MustCompile(`^true$`),
						"device_authorization_enabled":            regexp.MustCompile(`^false$`),
						"device_authorization_extra_verification": regexp.MustCompile(`^$`),
						"auto_enrollment_enabled":                 regexp.MustCompile(`^true$`),
						"integrity_detection":                     regexp.MustCompile(`^permissive$`),
					}),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "security_key.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "platform.#", "0"),
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

func TestAccMFAPolicy_Mobile_IntegrityDetectionErrors(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.TestAccCheckMFAPolicyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFAPolicyConfig_MobileIntegrityDetectionError_1(resourceName, name),
				// Integrity detection (`mobile.application.integrity_detection`) has no effect when the MFAPolicy resource has integrity detection disabled
				ExpectError: regexp.MustCompile("Integrity detection \\(`mobile\\.application\\.integrity_detection`\\) has no effect when the Application resource has integrity detection disabled"),
			},
			{
				Config: testAccMFAPolicyConfig_MobileIntegrityDetectionError_2(resourceName, name),
				// Integrity detection (`mobile.application.integrity_detection`) must be set when the MFAPolicy resource has integrity detection enabled
				ExpectError: regexp.MustCompile("Integrity detection \\(`mobile\\.application\\.integrity_detection`\\) must be set when the Application resource has integrity detection enabled"),
			},
		},
	})
}

func TestAccMFAPolicy_Mobile_BadMFAPolicyErrors(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.TestAccCheckMFAPolicyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccMFAPolicyConfig_MobileBadApplicationError_1(resourceName, name),
				ExpectError: regexp.MustCompile("Application referenced in `mobile.application.id` does not exist"),
			},
			{
				Config:      testAccMFAPolicyConfig_MobileBadApplicationError_2(resourceName, name),
				ExpectError: regexp.MustCompile("Application referenced in `mobile.application.id` is not of type OIDC"),
			},
			{
				Config:      testAccMFAPolicyConfig_MobileBadApplicationError_3(resourceName, name),
				ExpectError: regexp.MustCompile("Application referenced in `mobile.application.id` is OIDC, but is not the required `Native` OIDC application type"),
			},
			{
				Config:      testAccMFAPolicyConfig_MobileBadApplicationError_4(resourceName, name),
				ExpectError: regexp.MustCompile("Application referenced in `mobile.application.id` does not contain mobile application configuration"),
			},
		},
	})
}

func TestAccMFAPolicy_Mobile_Minimal(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_policy.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.TestAccCheckMFAPolicyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFAPolicyConfig_MinimalMobile(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.otp_failure_count", "3"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.otp_failure_cooldown_duration", "2"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.otp_failure_cooldown_timeunit", "MINUTES"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "security_key.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "platform.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "NONE"),
				),
			},
		},
	})
}

func TestAccMFAPolicy_Mobile_Change(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_policy.%s", resourceName)

	name := resourceName

	firebaseCredentials := os.Getenv("PINGONE_GOOGLE_FIREBASE_CREDENTIALS")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckGoogleFirebaseCredentials(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.TestAccCheckMFAPolicyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFAPolicyConfig_FullMobile(resourceName, name, firebaseCredentials),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.otp_failure_count", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.otp_failure_cooldown_duration", "125"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.otp_failure_cooldown_timeunit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.application.#", "3"),
					resource.TestMatchTypeSetElemNestedAttrs(resourceFullName, "mobile.0.application.*", map[string]*regexp.Regexp{
						"id":                                      verify.P1ResourceIDRegexpFullString,
						"pairing_disabled":                        regexp.MustCompile(`^true$`),
						"push_enabled":                            regexp.MustCompile(`^true$`),
						"push_timeout_duration":                   regexp.MustCompile(`^100$`),
						"push_timeout_timeunit":                   regexp.MustCompile(`^SECONDS$`),
						"pairing_key_lifetime_duration":           regexp.MustCompile(`^3$`),
						"pairing_key_lifetime_timeunit":           regexp.MustCompile(`^HOURS$`),
						"push_limit_count":                        regexp.MustCompile(`^10$`),
						"push_limit_lock_duration":                regexp.MustCompile(`^260$`),
						"push_limit_lock_duration_timeunit":       regexp.MustCompile(`^SECONDS$`),
						"push_limit_time_period_duration":         regexp.MustCompile(`^300$`),
						"push_limit_time_period_timeunit":         regexp.MustCompile(`^SECONDS$`),
						"otp_enabled":                             regexp.MustCompile(`^true$`),
						"device_authorization_enabled":            regexp.MustCompile(`^true$`),
						"device_authorization_extra_verification": regexp.MustCompile(`^restrictive$`),
						"auto_enrollment_enabled":                 regexp.MustCompile(`^true$`),
						"integrity_detection":                     regexp.MustCompile(`^$`),
					}),
					resource.TestMatchTypeSetElemNestedAttrs(resourceFullName, "mobile.0.application.*", map[string]*regexp.Regexp{
						"id":                                      verify.P1ResourceIDRegexpFullString,
						"pairing_disabled":                        regexp.MustCompile(`^false$`),
						"push_enabled":                            regexp.MustCompile(`^false$`),
						"push_timeout_duration":                   regexp.MustCompile(`^40$`),
						"push_timeout_timeunit":                   regexp.MustCompile(`^SECONDS$`),
						"pairing_key_lifetime_duration":           regexp.MustCompile(`^10$`),
						"pairing_key_lifetime_timeunit":           regexp.MustCompile(`^MINUTES$`),
						"push_limit_count":                        regexp.MustCompile(`^5$`),
						"push_limit_lock_duration":                regexp.MustCompile(`^30$`),
						"push_limit_lock_duration_timeunit":       regexp.MustCompile(`^MINUTES$`),
						"push_limit_time_period_duration":         regexp.MustCompile(`^10$`),
						"push_limit_time_period_timeunit":         regexp.MustCompile(`^MINUTES$`),
						"otp_enabled":                             regexp.MustCompile(`^true$`),
						"device_authorization_extra_verification": regexp.MustCompile(`^$`),
						"integrity_detection":                     regexp.MustCompile(`^permissive$`),
					}),
					resource.TestMatchTypeSetElemNestedAttrs(resourceFullName, "mobile.0.application.*", map[string]*regexp.Regexp{
						"id":                                      verify.P1ResourceIDRegexpFullString,
						"pairing_disabled":                        regexp.MustCompile(`^false$`),
						"push_enabled":                            regexp.MustCompile(`^true$`),
						"push_timeout_duration":                   regexp.MustCompile(`^40$`),
						"push_timeout_timeunit":                   regexp.MustCompile(`^SECONDS$`),
						"pairing_key_lifetime_duration":           regexp.MustCompile(`^55$`),
						"pairing_key_lifetime_timeunit":           regexp.MustCompile(`^MINUTES$`),
						"push_limit_count":                        regexp.MustCompile(`^5$`),
						"push_limit_lock_duration":                regexp.MustCompile(`^25$`),
						"push_limit_lock_duration_timeunit":       regexp.MustCompile(`^MINUTES$`),
						"push_limit_time_period_duration":         regexp.MustCompile(`^5$`),
						"push_limit_time_period_timeunit":         regexp.MustCompile(`^MINUTES$`),
						"otp_enabled":                             regexp.MustCompile(`^true$`),
						"device_authorization_enabled":            regexp.MustCompile(`^false$`),
						"device_authorization_extra_verification": regexp.MustCompile(`^$`),
						"auto_enrollment_enabled":                 regexp.MustCompile(`^true$`),
						"integrity_detection":                     regexp.MustCompile(`^permissive$`),
					}),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "security_key.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "platform.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "SMS_THEN_EMAIL"),
				),
			},
			{
				Config: testAccMFAPolicyConfig_MinimalMobile(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.otp_failure_count", "3"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.otp_failure_cooldown_duration", "2"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.otp_failure_cooldown_timeunit", "MINUTES"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "security_key.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "platform.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "NONE"),
				),
			},
			{
				Config: testAccMFAPolicyConfig_FullMobile(resourceName, name, firebaseCredentials),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.otp_failure_count", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.otp_failure_cooldown_duration", "125"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.otp_failure_cooldown_timeunit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.application.#", "3"),
					resource.TestMatchTypeSetElemNestedAttrs(resourceFullName, "mobile.0.application.*", map[string]*regexp.Regexp{
						"id":                                      verify.P1ResourceIDRegexpFullString,
						"pairing_disabled":                        regexp.MustCompile(`^true$`),
						"push_enabled":                            regexp.MustCompile(`^true$`),
						"push_timeout_duration":                   regexp.MustCompile(`^100$`),
						"push_timeout_timeunit":                   regexp.MustCompile(`^SECONDS$`),
						"pairing_key_lifetime_duration":           regexp.MustCompile(`^3$`),
						"pairing_key_lifetime_timeunit":           regexp.MustCompile(`^HOURS$`),
						"push_limit_count":                        regexp.MustCompile(`^10$`),
						"push_limit_lock_duration":                regexp.MustCompile(`^260$`),
						"push_limit_lock_duration_timeunit":       regexp.MustCompile(`^SECONDS$`),
						"push_limit_time_period_duration":         regexp.MustCompile(`^300$`),
						"push_limit_time_period_timeunit":         regexp.MustCompile(`^SECONDS$`),
						"otp_enabled":                             regexp.MustCompile(`^true$`),
						"device_authorization_enabled":            regexp.MustCompile(`^true$`),
						"device_authorization_extra_verification": regexp.MustCompile(`^restrictive$`),
						"auto_enrollment_enabled":                 regexp.MustCompile(`^true$`),
						"integrity_detection":                     regexp.MustCompile(`^$`),
					}),
					resource.TestMatchTypeSetElemNestedAttrs(resourceFullName, "mobile.0.application.*", map[string]*regexp.Regexp{
						"id":                                      verify.P1ResourceIDRegexpFullString,
						"pairing_disabled":                        regexp.MustCompile(`^false$`),
						"push_enabled":                            regexp.MustCompile(`^false$`),
						"push_timeout_duration":                   regexp.MustCompile(`^40$`),
						"push_timeout_timeunit":                   regexp.MustCompile(`^SECONDS$`),
						"pairing_key_lifetime_duration":           regexp.MustCompile(`^10$`),
						"pairing_key_lifetime_timeunit":           regexp.MustCompile(`^MINUTES$`),
						"push_limit_count":                        regexp.MustCompile(`^5$`),
						"push_limit_lock_duration":                regexp.MustCompile(`^30$`),
						"push_limit_lock_duration_timeunit":       regexp.MustCompile(`^MINUTES$`),
						"push_limit_time_period_duration":         regexp.MustCompile(`^10$`),
						"push_limit_time_period_timeunit":         regexp.MustCompile(`^MINUTES$`),
						"otp_enabled":                             regexp.MustCompile(`^true$`),
						"device_authorization_extra_verification": regexp.MustCompile(`^$`),
						"integrity_detection":                     regexp.MustCompile(`^permissive$`),
					}),
					resource.TestMatchTypeSetElemNestedAttrs(resourceFullName, "mobile.0.application.*", map[string]*regexp.Regexp{
						"id":                                      verify.P1ResourceIDRegexpFullString,
						"pairing_disabled":                        regexp.MustCompile(`^false$`),
						"push_enabled":                            regexp.MustCompile(`^true$`),
						"push_timeout_duration":                   regexp.MustCompile(`^40$`),
						"push_timeout_timeunit":                   regexp.MustCompile(`^SECONDS$`),
						"pairing_key_lifetime_duration":           regexp.MustCompile(`^55$`),
						"pairing_key_lifetime_timeunit":           regexp.MustCompile(`^MINUTES$`),
						"push_limit_count":                        regexp.MustCompile(`^5$`),
						"push_limit_lock_duration":                regexp.MustCompile(`^25$`),
						"push_limit_lock_duration_timeunit":       regexp.MustCompile(`^MINUTES$`),
						"push_limit_time_period_duration":         regexp.MustCompile(`^5$`),
						"push_limit_time_period_timeunit":         regexp.MustCompile(`^MINUTES$`),
						"otp_enabled":                             regexp.MustCompile(`^true$`),
						"device_authorization_enabled":            regexp.MustCompile(`^false$`),
						"device_authorization_extra_verification": regexp.MustCompile(`^$`),
						"auto_enrollment_enabled":                 regexp.MustCompile(`^true$`),
						"integrity_detection":                     regexp.MustCompile(`^permissive$`),
					}),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "security_key.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "platform.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "SMS_THEN_EMAIL"),
				),
			},
		},
	})
}

func TestAccMFAPolicy_Totp_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_policy.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.TestAccCheckMFAPolicyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFAPolicyConfig_FullTotp(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.pairing_disabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.otp_failure_count", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.otp_failure_cooldown_duration", "125"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.otp_failure_cooldown_timeunit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "security_key.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "platform.#", "0"),
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

func TestAccMFAPolicy_Totp_Minimal(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_policy.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.TestAccCheckMFAPolicyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFAPolicyConfig_MinimalTotp(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.pairing_disabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.otp_failure_count", "3"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.otp_failure_cooldown_duration", "2"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.otp_failure_cooldown_timeunit", "MINUTES"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "security_key.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "platform.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "NONE"),
				),
			},
		},
	})
}

func TestAccMFAPolicy_Totp_Change(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_policy.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.TestAccCheckMFAPolicyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFAPolicyConfig_FullTotp(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.pairing_disabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.otp_failure_count", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.otp_failure_cooldown_duration", "125"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.otp_failure_cooldown_timeunit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "security_key.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "platform.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "SMS_THEN_EMAIL"),
				),
			},
			{
				Config: testAccMFAPolicyConfig_MinimalTotp(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.pairing_disabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.otp_failure_count", "3"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.otp_failure_cooldown_duration", "2"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.otp_failure_cooldown_timeunit", "MINUTES"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "security_key.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "platform.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "NONE"),
				),
			},
			{
				Config: testAccMFAPolicyConfig_FullTotp(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.pairing_disabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.otp_failure_count", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.otp_failure_cooldown_duration", "125"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.otp_failure_cooldown_timeunit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "security_key.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "platform.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "SMS_THEN_EMAIL"),
				),
			},
		},
	})
}

func TestAccMFAPolicy_FIDO2_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_policy.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.TestAccCheckMFAPolicyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFAPolicyConfig_FullFIDO2(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.0.pairing_disabled", "true"),
					resource.TestMatchResourceAttr(resourceFullName, "fido2.0.fido2_policy_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "security_key.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "platform.#", "0"),
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

func TestAccMFAPolicy_FIDO2_Minimal(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_policy.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.TestAccCheckMFAPolicyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFAPolicyConfig_MinimalFIDO2(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.0.pairing_disabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.0.fido2_policy_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "security_key.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "platform.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "NONE"),
				),
			},
		},
	})
}

func TestAccMFAPolicy_FIDO2_Change(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_policy.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.TestAccCheckMFAPolicyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFAPolicyConfig_FullFIDO2(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.0.pairing_disabled", "true"),
					resource.TestMatchResourceAttr(resourceFullName, "fido2.0.fido2_policy_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "security_key.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "platform.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "SMS_THEN_EMAIL"),
				),
			},
			{
				Config: testAccMFAPolicyConfig_MinimalFIDO2(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.0.pairing_disabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.0.fido2_policy_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "security_key.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "platform.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "NONE"),
				),
			},
			{
				Config: testAccMFAPolicyConfig_FullFIDO2(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.0.pairing_disabled", "true"),
					resource.TestMatchResourceAttr(resourceFullName, "fido2.0.fido2_policy_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "security_key.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "platform.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "SMS_THEN_EMAIL"),
				),
			},
		},
	})
}

func TestAccMFAPolicy_SecurityKey_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_policy.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		//PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		PreCheck:                 func() { t.Skipf("Device method deprecated for new environments") },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.TestAccCheckMFAPolicyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFAPolicyConfig_FullSecurityKey(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "security_key.0.enabled", "true"),
					resource.TestMatchResourceAttr(resourceFullName, "security_key.0.fido_policy_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "platform.0.enabled", "false"),
				),
			},
		},
	})
}

func TestAccMFAPolicy_SecurityKey_Minimal(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_policy.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		//PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		PreCheck:                 func() { t.Skipf("Device method deprecated for new environments") },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.TestAccCheckMFAPolicyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFAPolicyConfig_MinimalSecurityKey(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "security_key.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "security_key.0.fido_policy_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "platform.0.enabled", "false"),
				),
			},
		},
	})
}

func TestAccMFAPolicy_SecurityKey_Change(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_policy.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		//PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		PreCheck:                 func() { t.Skipf("Device method deprecated for new environments") },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.TestAccCheckMFAPolicyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFAPolicyConfig_FullSecurityKey(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "security_key.0.enabled", "true"),
					resource.TestMatchResourceAttr(resourceFullName, "security_key.0.fido_policy_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "platform.0.enabled", "false"),
				),
			},
			{
				Config: testAccMFAPolicyConfig_MinimalSecurityKey(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "security_key.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "security_key.0.fido_policy_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "platform.0.enabled", "false"),
				),
			},
			{
				Config: testAccMFAPolicyConfig_FullSecurityKey(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "security_key.0.enabled", "true"),
					resource.TestMatchResourceAttr(resourceFullName, "security_key.0.fido_policy_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "platform.0.enabled", "false"),
				),
			},
		},
	})
}

func TestAccMFAPolicy_Platform_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_policy.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		//PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		PreCheck:                 func() { t.Skipf("Device method deprecated for new environments") },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.TestAccCheckMFAPolicyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFAPolicyConfig_FullPlatform(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "security_key.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "platform.0.enabled", "true"),
					resource.TestMatchResourceAttr(resourceFullName, "platform.0.fido_policy_id", verify.P1ResourceIDRegexpFullString),
				),
			},
		},
	})
}

func TestAccMFAPolicy_Platform_Minimal(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_policy.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		//PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		PreCheck:                 func() { t.Skipf("Device method deprecated for new environments") },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.TestAccCheckMFAPolicyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFAPolicyConfig_MinimalPlatform(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "security_key.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "platform.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "platform.0.fido_policy_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "EMAIL_THEN_SMS"),
				),
			},
		},
	})
}

func TestAccMFAPolicy_Platform_Change(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_policy.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		//PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		PreCheck:                 func() { t.Skipf("Device method deprecated for new environments") },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.TestAccCheckMFAPolicyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFAPolicyConfig_FullPlatform(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "security_key.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "platform.0.enabled", "true"),
					resource.TestMatchResourceAttr(resourceFullName, "platform.0.fido_policy_id", verify.P1ResourceIDRegexpFullString),
				),
			},
			{
				Config: testAccMFAPolicyConfig_MinimalPlatform(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "security_key.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "platform.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "platform.0.fido_policy_id", ""),
				),
			},
			{
				Config: testAccMFAPolicyConfig_FullPlatform(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "security_key.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "platform.0.enabled", "true"),
					resource.TestMatchResourceAttr(resourceFullName, "platform.0.fido_policy_id", verify.P1ResourceIDRegexpFullString),
				),
			},
		},
	})
}

func TestAccMFAPolicy_DataModel(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_policy.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.TestAccCheckMFAPolicyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Minimal from new
			{
				Config: testAccMFAPolicyConfig_MinimalFIDO2(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "device_selection", "DEFAULT_TO_FIRST"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "NONE"),
				),
			},
			{
				Config:  testAccMFAPolicyConfig_MinimalFIDO2(resourceName, name),
				Destroy: true,
			},
			// Full from new
			{
				Config: testAccMFAPolicyConfig_FullFIDO2(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "device_selection", "PROMPT_TO_SELECT"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "SMS_THEN_EMAIL"),
				),
			},
			{
				Config:  testAccMFAPolicyConfig_FullFIDO2(resourceName, name),
				Destroy: true,
			},
			// Update
			{
				Config: testAccMFAPolicyConfig_FullFIDO2(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "device_selection", "PROMPT_TO_SELECT"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "SMS_THEN_EMAIL"),
				),
			},
			{
				Config: testAccMFAPolicyConfig_MinimalFIDO2(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "device_selection", "DEFAULT_TO_FIRST"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "NONE"),
				),
			},
			{
				Config: testAccMFAPolicyConfig_FullFIDO2(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "device_selection", "PROMPT_TO_SELECT"),
					resource.TestCheckResourceAttr(resourceFullName, "new_device_notification", "SMS_THEN_EMAIL"),
				),
			},
		},
	})
}

func TestAccMFAPolicy_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_policy.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.TestAccCheckMFAPolicyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccMFAPolicyConfig_FullSMS(resourceName, name),
			},
			// Errors
			{
				ResourceName: resourceFullName,
				ImportState:  true,
				ExpectError:  regexp.MustCompile(`Invalid import ID specified \(".*"\).  The ID should be in the format "environment_id/mfa_device_policy_id" and must match regex: .*`),
			},
			{
				ResourceName:  resourceFullName,
				ImportStateId: "/",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`Invalid import ID specified \(".*"\).  The ID should be in the format "environment_id/mfa_device_policy_id" and must match regex: .*`),
			},
			{
				ResourceName:  resourceFullName,
				ImportStateId: "badformat/badformat",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`Invalid import ID specified \(".*"\).  The ID should be in the format "environment_id/mfa_device_policy_id" and must match regex: .*`),
			},
		},
	})
}

func testAccMFAPolicyConfig_NewEnv(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_mfa_policy" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s"

  sms {
    enabled = true
  }

  voice {
    enabled = true
  }

  email {
    enabled = true
  }

  mobile {
    enabled = true
  }

  totp {
    enabled = true
  }

  fido2 {
    enabled = true
  }

}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccMFAPolicyConfig_FullSMS(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_mfa_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  new_device_notification = "SMS_THEN_EMAIL"

  sms {
    enabled          = true
    pairing_disabled = true

    otp_lifetime_duration = 75
    otp_lifetime_timeunit = "SECONDS"

    otp_failure_count = 5

    otp_failure_cooldown_duration = 5
    otp_failure_cooldown_timeunit = "SECONDS"
  }

  voice {
    enabled = false
  }

  email {
    enabled = false
  }

  mobile {
    enabled = false
  }

  totp {
    enabled = false
  }

  fido2 {
    enabled = false
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccMFAPolicyConfig_MinimalSMS(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_mfa_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  sms {
    enabled = true
  }

  voice {
    enabled = false
  }

  email {
    enabled = false
  }

  mobile {
    enabled = false
  }

  totp {
    enabled = false
  }

  fido2 {
    enabled = false
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccMFAPolicyConfig_FullVoice(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_mfa_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  new_device_notification = "SMS_THEN_EMAIL"

  sms {
    enabled = false
  }

  voice {
    enabled          = true
    pairing_disabled = true

    otp_lifetime_duration = 75
    otp_lifetime_timeunit = "SECONDS"

    otp_failure_count = 5

    otp_failure_cooldown_duration = 5
    otp_failure_cooldown_timeunit = "SECONDS"
  }

  email {
    enabled = false
  }

  mobile {
    enabled = false
  }

  totp {
    enabled = false
  }

  fido2 {
    enabled = false
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccMFAPolicyConfig_MinimalVoice(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_mfa_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  sms {
    enabled = false
  }

  voice {
    enabled = true
  }

  email {
    enabled = false
  }

  mobile {
    enabled = false
  }

  totp {
    enabled = false
  }

  fido2 {
    enabled = false
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccMFAPolicyConfig_FullEmail(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_mfa_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  new_device_notification = "SMS_THEN_EMAIL"

  sms {
    enabled = false
  }

  voice {
    enabled = false
  }

  email {
    enabled          = true
    pairing_disabled = true

    otp_lifetime_duration = 75
    otp_lifetime_timeunit = "SECONDS"

    otp_failure_count = 5

    otp_failure_cooldown_duration = 5
    otp_failure_cooldown_timeunit = "SECONDS"
  }

  mobile {
    enabled = false
  }

  totp {
    enabled = false
  }

  fido2 {
    enabled = false
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccMFAPolicyConfig_MinimalEmail(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_mfa_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  sms {
    enabled = false
  }

  voice {
    enabled = false
  }

  email {
    enabled = true
  }

  mobile {
    enabled = false
  }

  totp {
    enabled = false
  }

  fido2 {
    enabled = false
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccMFAPolicyConfig_FullMobile(resourceName, name, key string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[2]s-1" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-1"
  description    = "My test OIDC app for MFA Policy"
  tags           = []
  login_page_url = "https://www.pingidentity.com"

  enabled = true

  oidc_options {
    type                        = "NATIVE_APP"
    grant_types                 = ["CLIENT_CREDENTIALS"]
    token_endpoint_authn_method = "CLIENT_SECRET_BASIC"

    mobile_app {
      bundle_id    = "com.%[2]s1.bundle"
      package_name = "com.%[2]s1.package"

      passcode_refresh_seconds = 45

      integrity_detection {
        enabled = false
      }
    }
  }
}

resource "pingone_mfa_application_push_credential" "%[2]s-1" {
  environment_id = data.pingone_environment.general_test.id
  application_id = pingone_application.%[2]s-1.id

  fcm {
    google_service_account_credentials = jsonencode(%[4]s)
  }
}

resource "pingone_application" "%[2]s-2" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-2"
  description    = "My test OIDC app for MFA Policy"
  tags           = []
  login_page_url = "https://www.pingidentity.com"

  enabled = true

  oidc_options {
    type                        = "NATIVE_APP"
    grant_types                 = ["CLIENT_CREDENTIALS"]
    token_endpoint_authn_method = "CLIENT_SECRET_BASIC"

    mobile_app {
      bundle_id    = "com.%[2]s2.bundle"
      package_name = "com.%[2]s2.package"

      passcode_refresh_seconds = 45

      integrity_detection {
        enabled = true
        cache_duration {
          amount = 30
          units  = "HOURS"
        }
        google_play {
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
  tags           = []
  login_page_url = "https://www.pingidentity.com"

  enabled = true

  oidc_options {
    type                        = "NATIVE_APP"
    grant_types                 = ["CLIENT_CREDENTIALS"]
    token_endpoint_authn_method = "CLIENT_SECRET_BASIC"

    mobile_app {
      bundle_id    = "com.%[2]s3.bundle"
      package_name = "com.%[2]s3.package"

      passcode_refresh_seconds = 45

      integrity_detection {
        enabled = true
        cache_duration {
          amount = 30
          units  = "HOURS"
        }
        google_play {
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

  fcm {
    google_service_account_credentials = jsonencode(%[4]s)
  }
}

resource "pingone_application" "%[2]s-4" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-4"
  description    = "My test OIDC app for MFA Policy"
  tags           = []
  login_page_url = "https://www.pingidentity.com"

  enabled = true

  oidc_options {
    type                        = "NATIVE_APP"
    grant_types                 = ["CLIENT_CREDENTIALS"]
    token_endpoint_authn_method = "CLIENT_SECRET_BASIC"

    mobile_app {
      bundle_id    = "com.%[2]s4.bundle"
      package_name = "com.%[2]s4.package"

      passcode_refresh_seconds = 45

      integrity_detection {
        enabled = true
        cache_duration {
          amount = 30
          units  = "HOURS"
        }
        google_play {
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

  fcm {
    key = "dummykey"
  }
}

resource "pingone_mfa_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  new_device_notification = "SMS_THEN_EMAIL"

  sms {
    enabled = false
  }

  voice {
    enabled = false
  }

  email {
    enabled = false
  }

  mobile {
    enabled = true

    otp_failure_count = 5

    otp_failure_cooldown_duration = 125
    otp_failure_cooldown_timeunit = "SECONDS"

    application {
      id = pingone_application.%[2]s-1.id

      pairing_disabled = true

      push_enabled          = true
      push_timeout_duration = 100

      pairing_key_lifetime_duration = 3
      pairing_key_lifetime_timeunit = "HOURS"

      push_limit_count                  = 10
      push_limit_lock_duration          = 260
      push_limit_lock_duration_timeunit = "SECONDS"
      push_limit_time_period_duration   = 300
      push_limit_time_period_timeunit   = "SECONDS"

      otp_enabled = true

      device_authorization_enabled            = true
      device_authorization_extra_verification = "restrictive"

      auto_enrollment_enabled = true
    }

    application {
      id = pingone_application.%[2]s-2.id

      pairing_disabled = false

      push_enabled = false
      otp_enabled  = true

      integrity_detection = "permissive"
    }

    application {
      id = pingone_application.%[2]s-3.id

      push_enabled = true
      otp_enabled  = true

      pairing_key_lifetime_duration = 55

      push_limit_lock_duration        = 25
      push_limit_time_period_duration = 5

      device_authorization_enabled = false

      auto_enrollment_enabled = true

      integrity_detection = "permissive"
    }
  }

  totp {
    enabled = false
  }

  fido2 {
    enabled = false
  }

  depends_on = [
    pingone_mfa_application_push_credential.%[2]s-1,
    pingone_mfa_application_push_credential.%[2]s-3
  ]

}`, acctest.GenericSandboxEnvironment(), resourceName, name, key)
}

func testAccMFAPolicyConfig_MobileIntegrityDetectionError_1(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "My test OIDC app for MFA Policy"
  tags           = []
  login_page_url = "https://www.pingidentity.com"

  enabled = true

  oidc_options {
    type                        = "NATIVE_APP"
    grant_types                 = ["CLIENT_CREDENTIALS"]
    token_endpoint_authn_method = "CLIENT_SECRET_BASIC"

    mobile_app {
      bundle_id    = "com.%[2]s.bundle"
      package_name = "com.%[2]s.package"

      passcode_refresh_seconds = 45

      integrity_detection {
        enabled = false
      }
    }
  }
}

resource "pingone_mfa_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  sms {
    enabled = false
  }

  voice {
    enabled = false
  }

  email {
    enabled = false
  }

  mobile {
    enabled = true

    otp_failure_count = 5

    otp_failure_cooldown_duration = 125
    otp_failure_cooldown_timeunit = "SECONDS"

    application {
      id = pingone_application.%[2]s.id

      push_enabled = true
      otp_enabled  = true

      device_authorization_enabled            = true
      device_authorization_extra_verification = "restrictive"

      auto_enrollment_enabled = true

      integrity_detection = "permissive"
    }
  }

  totp {
    enabled = false
  }

  fido2 {
    enabled = false
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccMFAPolicyConfig_MobileIntegrityDetectionError_2(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "My test OIDC app for MFA Policy"
  tags           = []
  login_page_url = "https://www.pingidentity.com"

  enabled = true

  oidc_options {
    type                        = "NATIVE_APP"
    grant_types                 = ["CLIENT_CREDENTIALS"]
    token_endpoint_authn_method = "CLIENT_SECRET_BASIC"

    mobile_app {
      bundle_id    = "com.%[2]s.bundle"
      package_name = "com.%[2]s.package"

      passcode_refresh_seconds = 45

      integrity_detection {
        enabled = true
        cache_duration {
          amount = 30
          units  = "HOURS"
        }
        google_play {
          verification_type = "INTERNAL"
          decryption_key    = "dummykeydoesnotexist"
          verification_key  = "dummykeydoesnotexist"
        }
      }
    }
  }
}

resource "pingone_mfa_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  sms {
    enabled = false
  }

  voice {
    enabled = false
  }

  email {
    enabled = false
  }

  mobile {
    enabled = true

    otp_failure_count = 5

    otp_failure_cooldown_duration = 125
    otp_failure_cooldown_timeunit = "SECONDS"

    application {
      id = pingone_application.%[2]s.id

      push_enabled = true
      otp_enabled  = true

      device_authorization_enabled            = true
      device_authorization_extra_verification = "restrictive"

      auto_enrollment_enabled = true
    }
  }

  totp {
    enabled = false
  }

  fido2 {
    enabled = false
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccMFAPolicyConfig_MobileBadApplicationError_1(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_mfa_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  sms {
    enabled = false
  }

  voice {
    enabled = false
  }

  email {
    enabled = false
  }

  mobile {
    enabled = true

    otp_failure_count = 5

    otp_failure_cooldown_duration = 125
    otp_failure_cooldown_timeunit = "SECONDS"

    application {
      id = "9bf6c075-78ba-4cd6-a5b1-96ec144d66ef" // Fake ID

      push_enabled = true
      otp_enabled  = true

      device_authorization_enabled            = true
      device_authorization_extra_verification = "restrictive"

      auto_enrollment_enabled = true
    }
  }

  totp {
    enabled = false
  }

  fido2 {
    enabled = false
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccMFAPolicyConfig_MobileBadApplicationError_2(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "My test SAML app for MFA Policy"
  login_page_url = "https://www.pingidentity.com"

  enabled = true

  saml_options {
    acs_urls           = ["https://my-saas-app.com"]
    assertion_duration = 3600
    sp_entity_id       = "sp:%[2]s:localhost"
  }
}

resource "pingone_mfa_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  sms {
    enabled = false
  }

  voice {
    enabled = false
  }

  email {
    enabled = false
  }

  mobile {
    enabled = true

    otp_failure_count = 5

    otp_failure_cooldown_duration = 125
    otp_failure_cooldown_timeunit = "SECONDS"

    application {
      id = pingone_application.%[2]s.id

      push_enabled = true
      otp_enabled  = true

      device_authorization_enabled            = true
      device_authorization_extra_verification = "restrictive"

      auto_enrollment_enabled = true
    }
  }

  totp {
    enabled = false
  }

  fido2 {
    enabled = false
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccMFAPolicyConfig_MobileBadApplicationError_3(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "My test OIDC app for MFA Policy"
  tags           = []
  login_page_url = "https://www.pingidentity.com"

  enabled = true

  oidc_options {
    type                        = "WORKER"
    grant_types                 = ["CLIENT_CREDENTIALS"]
    token_endpoint_authn_method = "CLIENT_SECRET_BASIC"
  }
}

resource "pingone_mfa_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  sms {
    enabled = false
  }

  voice {
    enabled = false
  }

  email {
    enabled = false
  }

  mobile {
    enabled = true

    otp_failure_count = 5

    otp_failure_cooldown_duration = 125
    otp_failure_cooldown_timeunit = "SECONDS"

    application {
      id = pingone_application.%[2]s.id

      push_enabled = true
      otp_enabled  = true

      device_authorization_enabled            = true
      device_authorization_extra_verification = "restrictive"

      auto_enrollment_enabled = true
    }
  }

  totp {
    enabled = false
  }

  fido2 {
    enabled = false
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccMFAPolicyConfig_MobileBadApplicationError_4(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "My test OIDC app for MFA Policy"
  tags           = []
  login_page_url = "https://www.pingidentity.com"

  enabled = true

  oidc_options {
    type                        = "NATIVE_APP"
    grant_types                 = ["CLIENT_CREDENTIALS"]
    token_endpoint_authn_method = "CLIENT_SECRET_BASIC"
  }
}

resource "pingone_mfa_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  sms {
    enabled = false
  }

  voice {
    enabled = false
  }

  email {
    enabled = false
  }

  mobile {
    enabled = true

    otp_failure_count = 5

    otp_failure_cooldown_duration = 125
    otp_failure_cooldown_timeunit = "SECONDS"

    application {
      id = pingone_application.%[2]s.id

      push_enabled = false
      otp_enabled  = true

      device_authorization_enabled = false

      auto_enrollment_enabled = true
    }
  }

  totp {
    enabled = false
  }

  fido2 {
    enabled = false
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccMFAPolicyConfig_MinimalMobile(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_mfa_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  sms {
    enabled = false
  }

  voice {
    enabled = false
  }

  email {
    enabled = false
  }

  mobile {
    enabled = true
  }

  totp {
    enabled = false
  }

  fido2 {
    enabled = false
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccMFAPolicyConfig_FullTotp(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_mfa_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  new_device_notification = "SMS_THEN_EMAIL"

  sms {
    enabled = false
  }

  voice {
    enabled = false
  }

  email {
    enabled = false
  }

  mobile {
    enabled = false
  }

  totp {
    enabled          = true
    pairing_disabled = true

    otp_failure_count = 5

    otp_failure_cooldown_duration = 125
    otp_failure_cooldown_timeunit = "SECONDS"
  }

  fido2 {
    enabled = false
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccMFAPolicyConfig_MinimalTotp(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_mfa_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  sms {
    enabled = false
  }

  voice {
    enabled = false
  }

  email {
    enabled = false
  }

  mobile {
    enabled = false
  }

  totp {
    enabled = true
  }

  fido2 {
    enabled = false
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccMFAPolicyConfig_FullFIDO2(resourceName, name string) string {
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

resource "pingone_mfa_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  device_selection        = "PROMPT_TO_SELECT"
  new_device_notification = "SMS_THEN_EMAIL"

  sms {
    enabled = false
  }

  voice {
    enabled = false
  }

  email {
    enabled = false
  }

  mobile {
    enabled = false
  }

  totp {
    enabled = false
  }

  fido2 {
    enabled          = true
    pairing_disabled = true

    fido2_policy_id = pingone_mfa_fido2_policy.%[2]s.id
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccMFAPolicyConfig_MinimalFIDO2(resourceName, name string) string {
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

resource "pingone_mfa_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  sms {
    enabled = false
  }

  voice {
    enabled = false
  }

  email {
    enabled = false
  }

  mobile {
    enabled = false
  }

  totp {
    enabled = false
  }

  fido2 {
    enabled = true
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccMFAPolicyConfig_FullSecurityKey(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_mfa_fido_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  attestation_requirements = "GLOBAL"
  resident_key_requirement = "REQUIRED"

}

resource "pingone_mfa_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  sms {
    enabled = false
  }

  voice {
    enabled = false
  }

  email {
    enabled = false
  }

  mobile {
    enabled = false
  }

  totp {
    enabled = false
  }

  security_key {
    enabled = true

    fido_policy_id = pingone_mfa_fido_policy.%[2]s.id
  }

  platform {
    enabled = false
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccMFAPolicyConfig_MinimalSecurityKey(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_mfa_fido_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  attestation_requirements = "GLOBAL"
  resident_key_requirement = "REQUIRED"

}

resource "pingone_mfa_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  sms {
    enabled = false
  }

  voice {
    enabled = false
  }

  email {
    enabled = false
  }

  mobile {
    enabled = false
  }

  totp {
    enabled = false
  }

  security_key {
    enabled = true
  }

  platform {
    enabled = false
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccMFAPolicyConfig_FullPlatform(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_mfa_fido_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  attestation_requirements = "GLOBAL"
  resident_key_requirement = "REQUIRED"

}

resource "pingone_mfa_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  device_selection = "PROMPT_TO_SELECT"

  sms {
    enabled = false
  }

  voice {
    enabled = false
  }

  email {
    enabled = false
  }

  mobile {
    enabled = false
  }

  totp {
    enabled = false
  }

  security_key {
    enabled = false
  }

  platform {
    enabled = true

    fido_policy_id = pingone_mfa_fido_policy.%[2]s.id
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccMFAPolicyConfig_MinimalPlatform(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_mfa_fido_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  attestation_requirements = "GLOBAL"
  resident_key_requirement = "REQUIRED"

}

resource "pingone_mfa_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  sms {
    enabled = false
  }

  voice {
    enabled = false
  }

  email {
    enabled = false
  }

  mobile {
    enabled = false
  }

  totp {
    enabled = false
  }

  security_key {
    enabled = false
  }

  platform {
    enabled = true
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
