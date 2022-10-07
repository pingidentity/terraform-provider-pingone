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
	pingone "github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
)

func testAccCheckMFAPolicyDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.MFAAPIClient
	apiClientManagement := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, pingone.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})

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

func TestAccMFAPolicy_NewEnv(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_policy.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckMFAPolicyDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFAPolicyConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
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
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckMFAPolicyDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFAPolicyConfig_FullSMS(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.otp_lifetime_duration", "75"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.otp_lifetime_timeunit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.otp_failure_count", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.otp_failure_cooldown_duration", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.otp_failure_cooldown_timeunit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "security_key.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "platform.0.enabled", "false"),
				),
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
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckMFAPolicyDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFAPolicyConfig_MinimalSMS(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.otp_lifetime_duration", "30"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.otp_lifetime_timeunit", "MINUTES"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.otp_failure_count", "3"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.otp_failure_cooldown_duration", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.otp_failure_cooldown_timeunit", "MINUTES"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "security_key.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "platform.0.enabled", "false"),
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
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckMFAPolicyDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFAPolicyConfig_FullSMS(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.otp_lifetime_duration", "75"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.otp_lifetime_timeunit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.otp_failure_count", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.otp_failure_cooldown_duration", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.otp_failure_cooldown_timeunit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "security_key.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "platform.0.enabled", "false"),
				),
			},
			{
				Config: testAccMFAPolicyConfig_MinimalSMS(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.otp_lifetime_duration", "30"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.otp_lifetime_timeunit", "MINUTES"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.otp_failure_count", "3"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.otp_failure_cooldown_duration", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.otp_failure_cooldown_timeunit", "MINUTES"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "security_key.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "platform.0.enabled", "false"),
				),
			},
			{
				Config: testAccMFAPolicyConfig_FullSMS(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.otp_lifetime_duration", "75"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.otp_lifetime_timeunit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.otp_failure_count", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.otp_failure_cooldown_duration", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.otp_failure_cooldown_timeunit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "security_key.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "platform.0.enabled", "false"),
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
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckMFAPolicyDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFAPolicyConfig_FullVoice(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.otp_lifetime_duration", "75"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.otp_lifetime_timeunit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.otp_failure_count", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.otp_failure_cooldown_duration", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.otp_failure_cooldown_timeunit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "security_key.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "platform.0.enabled", "false"),
				),
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
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckMFAPolicyDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFAPolicyConfig_MinimalVoice(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.otp_lifetime_duration", "30"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.otp_lifetime_timeunit", "MINUTES"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.otp_failure_count", "3"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.otp_failure_cooldown_duration", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.otp_failure_cooldown_timeunit", "MINUTES"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "security_key.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "platform.0.enabled", "false"),
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
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckMFAPolicyDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFAPolicyConfig_FullVoice(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.otp_lifetime_duration", "75"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.otp_lifetime_timeunit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.otp_failure_count", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.otp_failure_cooldown_duration", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.otp_failure_cooldown_timeunit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "security_key.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "platform.0.enabled", "false"),
				),
			},
			{
				Config: testAccMFAPolicyConfig_MinimalVoice(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.otp_lifetime_duration", "30"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.otp_lifetime_timeunit", "MINUTES"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.otp_failure_count", "3"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.otp_failure_cooldown_duration", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.otp_failure_cooldown_timeunit", "MINUTES"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "security_key.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "platform.0.enabled", "false"),
				),
			},
			{
				Config: testAccMFAPolicyConfig_FullVoice(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.otp_lifetime_duration", "75"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.otp_lifetime_timeunit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.otp_failure_count", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.otp_failure_cooldown_duration", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.otp_failure_cooldown_timeunit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "security_key.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "platform.0.enabled", "false"),
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
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckMFAPolicyDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFAPolicyConfig_FullEmail(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.otp_lifetime_duration", "75"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.otp_lifetime_timeunit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.otp_failure_count", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.otp_failure_cooldown_duration", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.otp_failure_cooldown_timeunit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "security_key.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "platform.0.enabled", "false"),
				),
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
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckMFAPolicyDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFAPolicyConfig_MinimalEmail(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.otp_lifetime_duration", "30"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.otp_lifetime_timeunit", "MINUTES"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.otp_failure_count", "3"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.otp_failure_cooldown_duration", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.otp_failure_cooldown_timeunit", "MINUTES"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "security_key.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "platform.0.enabled", "false"),
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
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckMFAPolicyDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFAPolicyConfig_FullEmail(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.otp_lifetime_duration", "75"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.otp_lifetime_timeunit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.otp_failure_count", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.otp_failure_cooldown_duration", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.otp_failure_cooldown_timeunit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "security_key.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "platform.0.enabled", "false"),
				),
			},
			{
				Config: testAccMFAPolicyConfig_MinimalEmail(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.otp_lifetime_duration", "30"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.otp_lifetime_timeunit", "MINUTES"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.otp_failure_count", "3"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.otp_failure_cooldown_duration", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.otp_failure_cooldown_timeunit", "MINUTES"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "security_key.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "platform.0.enabled", "false"),
				),
			},
			{
				Config: testAccMFAPolicyConfig_FullEmail(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.otp_lifetime_duration", "75"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.otp_lifetime_timeunit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.otp_failure_count", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.otp_failure_cooldown_duration", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.otp_failure_cooldown_timeunit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "security_key.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "platform.0.enabled", "false"),
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

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckMFAPolicyDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFAPolicyConfig_FullMobile(resourceName, name),
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
						"id":           regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
						"push_enabled": regexp.MustCompile(`^false$`),
						"otp_enabled":  regexp.MustCompile(`^true$`),
						//"device_authorization_enabled": regexp.MustCompile(`^false$`),
						"device_authorization_extra_verification": regexp.MustCompile(`^$`),
						"auto_enrollment_enabled":                 regexp.MustCompile(`^true$`),
						"integrity_detection":                     regexp.MustCompile(`^restrictive$`),
					}),
					resource.TestMatchTypeSetElemNestedAttrs(resourceFullName, "mobile.0.application.*", map[string]*regexp.Regexp{
						"id":           regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
						"push_enabled": regexp.MustCompile(`^false$`),
						"otp_enabled":  regexp.MustCompile(`^true$`),
						"device_authorization_extra_verification": regexp.MustCompile(`^$`),
						"integrity_detection":                     regexp.MustCompile(`^restrictive$`),
					}),
					resource.TestMatchTypeSetElemNestedAttrs(resourceFullName, "mobile.0.application.*", map[string]*regexp.Regexp{
						"id":                           regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
						"push_enabled":                 regexp.MustCompile(`^false$`),
						"otp_enabled":                  regexp.MustCompile(`^true$`),
						"device_authorization_enabled": regexp.MustCompile(`^false$`),
						"device_authorization_extra_verification": regexp.MustCompile(`^$`),
						"auto_enrollment_enabled":                 regexp.MustCompile(`^true$`),
						"integrity_detection":                     regexp.MustCompile(`^permissive$`),
					}),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "security_key.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "platform.0.enabled", "false"),
				),
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
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckMFAPolicyDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
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
					resource.TestCheckResourceAttr(resourceFullName, "security_key.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "platform.0.enabled", "false"),
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

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckMFAPolicyDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFAPolicyConfig_FullMobile(resourceName, name),
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
						"id":           regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
						"push_enabled": regexp.MustCompile(`^false$`),
						"otp_enabled":  regexp.MustCompile(`^true$`),
						//"device_authorization_enabled": regexp.MustCompile(`^false$`),
						"device_authorization_extra_verification": regexp.MustCompile(`^$`),
						"auto_enrollment_enabled":                 regexp.MustCompile(`^true$`),
						"integrity_detection":                     regexp.MustCompile(`^restrictive$`),
					}),
					resource.TestMatchTypeSetElemNestedAttrs(resourceFullName, "mobile.0.application.*", map[string]*regexp.Regexp{
						"id":           regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
						"push_enabled": regexp.MustCompile(`^false$`),
						"otp_enabled":  regexp.MustCompile(`^true$`),
						"device_authorization_extra_verification": regexp.MustCompile(`^$`),
						"integrity_detection":                     regexp.MustCompile(`^restrictive$`),
					}),
					resource.TestMatchTypeSetElemNestedAttrs(resourceFullName, "mobile.0.application.*", map[string]*regexp.Regexp{
						"id":                           regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
						"push_enabled":                 regexp.MustCompile(`^false$`),
						"otp_enabled":                  regexp.MustCompile(`^true$`),
						"device_authorization_enabled": regexp.MustCompile(`^false$`),
						"device_authorization_extra_verification": regexp.MustCompile(`^$`),
						"auto_enrollment_enabled":                 regexp.MustCompile(`^true$`),
						"integrity_detection":                     regexp.MustCompile(`^permissive$`),
					}),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "security_key.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "platform.0.enabled", "false"),
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
					resource.TestCheckResourceAttr(resourceFullName, "security_key.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "platform.0.enabled", "false"),
				),
			},
			{
				Config: testAccMFAPolicyConfig_FullMobile(resourceName, name),
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
						"id":           regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
						"push_enabled": regexp.MustCompile(`^false$`),
						"otp_enabled":  regexp.MustCompile(`^true$`),
						//"device_authorization_enabled": regexp.MustCompile(`^false$`),
						"device_authorization_extra_verification": regexp.MustCompile(`^$`),
						"auto_enrollment_enabled":                 regexp.MustCompile(`^true$`),
						"integrity_detection":                     regexp.MustCompile(`^restrictive$`),
					}),
					resource.TestMatchTypeSetElemNestedAttrs(resourceFullName, "mobile.0.application.*", map[string]*regexp.Regexp{
						"id":           regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
						"push_enabled": regexp.MustCompile(`^false$`),
						"otp_enabled":  regexp.MustCompile(`^true$`),
						"device_authorization_extra_verification": regexp.MustCompile(`^$`),
						"integrity_detection":                     regexp.MustCompile(`^restrictive$`),
					}),
					resource.TestMatchTypeSetElemNestedAttrs(resourceFullName, "mobile.0.application.*", map[string]*regexp.Regexp{
						"id":                           regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
						"push_enabled":                 regexp.MustCompile(`^false$`),
						"otp_enabled":                  regexp.MustCompile(`^true$`),
						"device_authorization_enabled": regexp.MustCompile(`^false$`),
						"device_authorization_extra_verification": regexp.MustCompile(`^$`),
						"auto_enrollment_enabled":                 regexp.MustCompile(`^true$`),
						"integrity_detection":                     regexp.MustCompile(`^permissive$`),
					}),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "security_key.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "platform.0.enabled", "false"),
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
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckMFAPolicyDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFAPolicyConfig_FullTotp(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.otp_failure_count", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.otp_failure_cooldown_duration", "125"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.otp_failure_cooldown_timeunit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "security_key.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "platform.0.enabled", "false"),
				),
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
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckMFAPolicyDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFAPolicyConfig_MinimalTotp(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.otp_failure_count", "3"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.otp_failure_cooldown_duration", "2"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.otp_failure_cooldown_timeunit", "MINUTES"),
					resource.TestCheckResourceAttr(resourceFullName, "security_key.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "platform.0.enabled", "false"),
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
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckMFAPolicyDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFAPolicyConfig_FullTotp(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "sms.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "voice.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "email.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.otp_failure_count", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.otp_failure_cooldown_duration", "125"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.otp_failure_cooldown_timeunit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "security_key.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "platform.0.enabled", "false"),
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
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.otp_failure_count", "3"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.otp_failure_cooldown_duration", "2"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.otp_failure_cooldown_timeunit", "MINUTES"),
					resource.TestCheckResourceAttr(resourceFullName, "security_key.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "platform.0.enabled", "false"),
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
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.otp_failure_count", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.otp_failure_cooldown_duration", "125"),
					resource.TestCheckResourceAttr(resourceFullName, "totp.0.otp_failure_cooldown_timeunit", "SECONDS"),
					resource.TestCheckResourceAttr(resourceFullName, "security_key.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "platform.0.enabled", "false"),
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
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckMFAPolicyDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
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
					// resource.TestMatchResourceAttr(resourceFullName, "security_key.0.fido_policy_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
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
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckMFAPolicyDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
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
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckMFAPolicyDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
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
					// resource.TestMatchResourceAttr(resourceFullName, "security_key.0.fido_policy_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
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
					// resource.TestMatchResourceAttr(resourceFullName, "security_key.0.fido_policy_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
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
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckMFAPolicyDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
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
					// resource.TestMatchResourceAttr(resourceFullName, "platform.0.fido_policy_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
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
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckMFAPolicyDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
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
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckMFAPolicyDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
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
					// resource.TestMatchResourceAttr(resourceFullName, "platform.0.fido_policy_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
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
					// resource.TestMatchResourceAttr(resourceFullName, "platform.0.fido_policy_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
				),
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

	security_key {
		enabled = true
	}

	platform {
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
			
			  sms {
				  enabled = true

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
		  
			  security_key {
				  enabled = false
			  }
		  
			  platform {
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
		  
			  security_key {
				  enabled = false
			  }
		  
			  platform {
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
			
			  sms {
				  enabled = false
			  }
		  
			  voice {
				enabled = true

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
		  
			  security_key {
				  enabled = false
			  }
		  
			  platform {
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
		  
			  security_key {
				  enabled = false
			  }
		  
			  platform {
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
			
			  sms {
				  enabled = false
			  }
		  
			  voice {
				  enabled = false
			  }
		  
			  email {
				enabled = true

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
		  
			  security_key {
				  enabled = false
			  }
		  
			  platform {
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
		  
			  security_key {
				  enabled = false
			  }
		  
			  platform {
				  enabled = false
			  }
		  
		  }`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccMFAPolicyConfig_FullMobile(resourceName, name string) string {
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
				  enabled = true
				  cache_duration {
					amount = 30
					units  = "HOURS"
				  }
				}
			  }
		  
			  bundle_id    = "com.%[2]s1.bundle"
			  package_name = "com.%[2]s1.package"
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
				}
			  }
		  
			  bundle_id    = "com.%[2]s2.bundle"
			  package_name = "com.%[2]s2.package"
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
				}
			  }
		  
			  bundle_id    = "com.%[2]s3.bundle"
			  package_name = "com.%[2]s3.package"
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
					id = pingone_application.%[2]s-1.id

					push_enabled = false //true
					otp_enabled = true

					// device_authorization_enabled = true
					// device_authorization_extra_verification = "restrictive"

					auto_enrollment_enabled = true

					integrity_detection = "restrictive"
				}

				application {
					id = pingone_application.%[2]s-2.id

					push_enabled = false
					otp_enabled = true
				}

				application {
					id = pingone_application.%[2]s-3.id

					push_enabled = false //true
					otp_enabled = true

					device_authorization_enabled = false

					auto_enrollment_enabled = true

					integrity_detection = "permissive"
				}
			  }
		  
			  totp {
				  enabled = false
			  }
		  
			  security_key {
				  enabled = false
			  }
		  
			  platform {
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
		  
			  security_key {
				  enabled = false
			  }
		  
			  platform {
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

				otp_failure_count = 5

				otp_failure_cooldown_duration = 125
				otp_failure_cooldown_timeunit = "SECONDS"
			  }
		  
			  security_key {
				  enabled = false
			  }
		  
			  platform {
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
		  
			  security_key {
				  enabled = false
			  }
		  
			  platform {
				  enabled = false
			  }
		  
		  }`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccMFAPolicyConfig_FullSecurityKey(resourceName, name string) string {
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
				  enabled = false
			  }
		  
			  security_key {
				enabled = true

				// fido_policy_id = 
			  }
		  
			  platform {
				  enabled = false
			  }
		  
		  }`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccMFAPolicyConfig_MinimalSecurityKey(resourceName, name string) string {
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

				// fido_policy_id = 
			  }
		  
		  }`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccMFAPolicyConfig_MinimalPlatform(resourceName, name string) string {
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
