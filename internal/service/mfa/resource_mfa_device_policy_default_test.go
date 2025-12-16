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
					resource.TestCheckResourceAttr(resourceFullName, "totp.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "fido2.enabled", "true"),
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

func testAccMFADevicePolicyDefaultConfig_Full(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_mfa_device_policy_default" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  policy_type    = "pingone_mfa"

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
  policy_type    = "pingone_mfa"

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
					resource.TestCheckResourceAttr(resourceFullName, "policy_type", "pingid"),
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
					resource.TestCheckResourceAttr(resourceFullName, "mobile.applications.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.applications.0.biometrics_enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "mobile.applications.0.type", "pingIdAppConfig"),
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
					resource.TestCheckResourceAttr(resourceFullName, "policy_type", "pingid"),
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

func testAccMFADevicePolicyDefaultConfig_PingID_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_mfa_device_policy_default" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  policy_type    = "pingid"

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
  policy_type    = "pingone_mfa"

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
  policy_type    = "pingone_mfa"

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
  policy_type    = "pingid"

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
  policy_type    = "pingone_mfa"

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

func testAccMFADevicePolicyDefaultConfig_InvalidRememberMeDurationMinutes(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_mfa_device_policy_default" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  policy_type    = "pingone_mfa"

  name = "%[4]s"

  remember_me = {
    web = {
      enabled = true
      life_time = {
        duration  = 0
        time_unit = "MINUTES"
      }
    }
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

func testAccMFADevicePolicyDefaultConfig_InvalidRememberMeDurationHours(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_mfa_device_policy_default" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  policy_type    = "pingone_mfa"

  name = "%[4]s"

  remember_me = {
    web = {
      enabled = true
      life_time = {
        duration  = 2500
        time_unit = "HOURS"
      }
    }
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

func testAccMFADevicePolicyDefaultConfig_InvalidRememberMeDurationDays(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_mfa_device_policy_default" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  policy_type    = "pingone_mfa"

  name = "%[4]s"

  remember_me = {
    web = {
      enabled = true
      life_time = {
        duration  = 100
        time_unit = "DAYS"
      }
    }
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
  policy_type    = "pingid"

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

    applications = [
      {
        id = data.pingone_application.%[2]s.id

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
    ]

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
func TestAccMFADevicePolicyDefault_Validation(t *testing.T) {
	t.Parallel()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	testCases := map[string]func(t *testing.T){
		"PingOneMFA_Validation": func(t *testing.T) {
			t.Parallel()

			resourceName := acctest.ResourceNameGen()
			environmentName := acctest.ResourceNameGenEnvironment()
			name := resourceName

			resource.Test(t, resource.TestCase{
				PreCheck: func() {
					acctest.PreCheckClient(t)
				},
				ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
				CheckDestroy:             mfa.MFADevicePolicyDefault_CheckDestroy,
				ErrorCheck:               acctest.ErrorCheck(t),
				Steps: []resource.TestStep{
					// Invalid notifications_policy ID format
					{
						Config:      testAccMFADevicePolicyDefaultConfig_InvalidNotificationsPolicyID(environmentName, licenseID, resourceName, name),
						ExpectError: regexp.MustCompile(`The PingOne resource ID is malformed`),
					},
					// Invalid remember_me duration for MINUTES
					{
						Config:      testAccMFADevicePolicyDefaultConfig_InvalidRememberMeDurationMinutes(environmentName, licenseID, resourceName, name),
						ExpectError: regexp.MustCompile(`Attribute remember_me.web.life_time.duration value must be between 1 and`),
					},
					// Invalid remember_me duration for HOURS
					{
						Config:      testAccMFADevicePolicyDefaultConfig_InvalidRememberMeDurationHours(environmentName, licenseID, resourceName, name),
						ExpectError: regexp.MustCompile(`Attribute remember_me.web.life_time.duration value must be between 1 and`),
					},
					// Invalid remember_me duration for DAYS
					{
						Config:      testAccMFADevicePolicyDefaultConfig_InvalidRememberMeDurationDays(environmentName, licenseID, resourceName, name),
						ExpectError: regexp.MustCompile(`Attribute remember_me.web.life_time.duration value must be between 1 and`),
					},
				},
			})
		},
		"PingID_ValidationErrors": func(t *testing.T) {
			t.Parallel()

			resourceName := acctest.ResourceNameGen()
			name := resourceName

			resource.Test(t, resource.TestCase{
				PreCheck: func() {
					acctest.PreCheckClient(t)
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
				},
			})
		},
		"PingID_Structure": func(t *testing.T) {
			t.Parallel()

			resourceName := acctest.ResourceNameGen()
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
					// Missing desktop block
					{
						Config:      testAccMFADevicePolicyDefaultConfig_PingID_MissingDesktop(resourceName, name),
						ExpectError: regexp.MustCompile(`The argument desktop is required because\s+policy_type is configured as:\s+"pingid"`),
					},
					// Missing yubikey block
					{
						Config:      testAccMFADevicePolicyDefaultConfig_PingID_MissingYubikey(resourceName, name),
						ExpectError: regexp.MustCompile(`The argument yubikey is required because\s+policy_type is configured as:\s+"pingid"`),
					},
					// Mobile must be enabled for PingID policies
					{
						Config:      testAccMFADevicePolicyDefaultConfig_PingID_MobileDisabled(resourceName, name),
						ExpectError: regexp.MustCompile(`Attribute mobile.enabled must be true when attribute policy_type value is`),
					},
				},
			})
		},
		"PingID_MobileApp": func(t *testing.T) {
			t.Parallel()

			resourceName := acctest.ResourceNameGen()
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
					// Auto enrollment conflict
					{
						Config:      testAccMFADevicePolicyDefaultConfig_PingID_MobileApp_AutoEnrollment(resourceName, name),
						ExpectError: regexp.MustCompile(`The argument cannot be defined if the value "pingid" is present`),
					},
					// Device authorization conflict
					{
						Config:      testAccMFADevicePolicyDefaultConfig_PingID_MobileApp_DeviceAuthorization(resourceName, name),
						ExpectError: regexp.MustCompile(`The argument cannot be defined if the value "pingid" is present`),
					},
					// Missing new_request_duration_configuration
					{
						Config:      testAccMFADevicePolicyDefaultConfig_PingID_MobileApp_MissingNewRequestDuration(resourceName, name),
						ExpectError: regexp.MustCompile(`The argument mobile.applications\[0\].new_request_duration_configuration is\s+required because\s+policy_type is configured as:\s+"pingid"`),
					},
					// Missing ip_pairing_configuration
					{
						Config:      testAccMFADevicePolicyDefaultConfig_PingID_MobileApp_MissingIPPairing(resourceName, name),
						ExpectError: regexp.MustCompile(`The argument mobile.applications\[0\].ip_pairing_configuration is\s+required\s+because\s+policy_type is configured as:\s+"pingid"`),
					},
				},
			})
		},
		"PingOneMFA_MobileApp": func(t *testing.T) {
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
						ExpectError: regexp.MustCompile(`The argument mobile.applications\[0\].auto_enrollment is required because\s+policy_type is configured as:\s+"pingone_mfa"`),
					},
					// Missing device_authorization
					{
						Config:      testAccMFADevicePolicyDefaultConfig_PingOneMFA_MobileApp_MissingDeviceAuthorization(environmentName, licenseID, resourceName, name),
						ExpectError: regexp.MustCompile(`The argument mobile.applications\[0\].device_authorization is required because\s+policy_type is configured as:\s+"pingone_mfa"`),
					},
					// Missing integrity_detection
					{
						Config:      testAccMFADevicePolicyDefaultConfig_PingOneMFA_MobileApp_MissingIntegrityDetection(environmentName, licenseID, resourceName, name),
						ExpectError: regexp.MustCompile(`The argument mobile.applications\[0\].integrity_detection is required because\s+policy_type is configured as:\s+"pingone_mfa"`),
					},
					// Biometrics enabled conflict
					{
						Config:      testAccMFADevicePolicyDefaultConfig_PingOneMFA_MobileApp_BiometricsEnabled(environmentName, licenseID, resourceName, name),
						ExpectError: regexp.MustCompile(`The argument cannot be defined if the value "pingone_mfa" is present`),
					},
					// New request duration configuration conflict
					{
						Config:      testAccMFADevicePolicyDefaultConfig_PingOneMFA_MobileApp_NewRequestDuration(environmentName, licenseID, resourceName, name),
						ExpectError: regexp.MustCompile(`The argument cannot be defined if the value "pingone_mfa" is present`),
					},
					// IP pairing configuration conflict
					{
						Config:      testAccMFADevicePolicyDefaultConfig_PingOneMFA_MobileApp_IPPairing(environmentName, licenseID, resourceName, name),
						ExpectError: regexp.MustCompile(`The argument cannot be defined if the value "pingone_mfa" is present`),
					},
				},
			})
		},
		"PingID_IPPairing": func(t *testing.T) {
			t.Parallel()

			resourceName := acctest.ResourceNameGen()
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
					// Invalid CIDR
					{
						Config:      testAccMFADevicePolicyDefaultConfig_PingID_IPPairing_InvalidCIDR(resourceName, name),
						ExpectError: regexp.MustCompile(`Expected value to be in CIDR notation`),
					},
					// Missing only_these_ip_addresses when any_ip_address is false
					{
						Config:      testAccMFADevicePolicyDefaultConfig_PingID_IPPairing_MissingIPs(resourceName, name),
						ExpectError: regexp.MustCompile(`The argument\s+mobile.applications\[0\].ip_pairing_configuration.only_these_ip_addresses is\s+required because\s+mobile.applications\[0\].ip_pairing_configuration.any_ip_address is configured\s+as: false`),
					},
				},
			})
		},
		"Desktop": func(t *testing.T) {
			t.Parallel()

			resourceName := acctest.ResourceNameGen()
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
					// OTP failure count too high
					{
						Config:      testAccMFADevicePolicyDefaultConfig_PingID_Desktop_OTPCountHigh(resourceName, name),
						ExpectError: regexp.MustCompile(`Attribute desktop.otp.failure.count value must be between 1 and 7`),
					},
					// Pairing key lifetime too long (HOURS)
					{
						Config:      testAccMFADevicePolicyDefaultConfig_PingID_Desktop_PairingKeyLifetimeHigh(resourceName, name),
						ExpectError: regexp.MustCompile(`Attribute desktop.pairing_key_lifetime.duration value must be between 1 and\s+48`),
					},
				},
			})
		},
		"RememberMe": func(t *testing.T) {
			t.Parallel()

			resourceName := acctest.ResourceNameGen()
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
					// Duration out of range for MINUTES
					{
						Config:      testAccMFADevicePolicyDefaultConfig_RememberMe_MinutesHigh(resourceName, name),
						ExpectError: regexp.MustCompile(`Attribute remember_me.web.life_time.duration value must be between 1 and[\s\n]+129600`),
					},
					// Duration out of range for HOURS
					{
						Config:      testAccMFADevicePolicyDefaultConfig_RememberMe_HoursHigh(resourceName, name),
						ExpectError: regexp.MustCompile(`Attribute remember_me.web.life_time.duration value must be between 1 and[\s\n]+2160`),
					},
					// Duration out of range for DAYS
					{
						Config:      testAccMFADevicePolicyDefaultConfig_RememberMe_DaysHigh(resourceName, name),
						ExpectError: regexp.MustCompile(`Attribute remember_me.web.life_time.duration value must be between 1 and[\s\n]+90`),
					},
				},
			})
		},
		"MobilePushLimit": func(t *testing.T) {
			t.Parallel()

			resourceName := acctest.ResourceNameGen()
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
					// Count out of range
					{
						Config:      testAccMFADevicePolicyDefaultConfig_MobilePushLimit_CountHigh(resourceName, name),
						ExpectError: regexp.MustCompile(`Attribute mobile.applications\[0\].push_limit.count value must be between 1 and[\s\n]+50`),
					},
				},
			})
		},
		"MobileNewRequestDuration": func(t *testing.T) {
			t.Parallel()

			resourceName := acctest.ResourceNameGen()
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
					// Device timeout duration out of range
					{
						Config:      testAccMFADevicePolicyDefaultConfig_MobileNewRequestDuration_DeviceTimeoutHigh(resourceName, name),
						ExpectError: regexp.MustCompile(`Attribute[\s\n]+mobile.applications\[0\].new_request_duration_configuration.device_timeout.duration[\s\n]+value must be between 15 and 75`),
					},
					// Total timeout duration out of range
					{
						Config:      testAccMFADevicePolicyDefaultConfig_MobileNewRequestDuration_TotalTimeoutHigh(resourceName, name),
						ExpectError: regexp.MustCompile(`Attribute[\s\n]+mobile.applications\[0\].new_request_duration_configuration.total_timeout.duration[\s\n]+value must be between 30 and 90`),
					},
				},
			})
		},
		"Authentication": func(t *testing.T) {
			t.Parallel()

			resourceName := acctest.ResourceNameGen()
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
					{
						Config:      testAccMFADevicePolicyDefaultConfig_Authentication(resourceName, name, "INVALID_VALUE"),
						ExpectError: regexp.MustCompile(`Attribute authentication.device_selection value must be one of:`),
					},
				},
			})
		},
		"NewDeviceNotification": func(t *testing.T) {
			t.Parallel()

			resourceName := acctest.ResourceNameGen()
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
					{
						Config:      testAccMFADevicePolicyDefaultConfig_NewDeviceNotification(resourceName, name, "INVALID_VALUE"),
						ExpectError: regexp.MustCompile(`Attribute new_device_notification value must be one of:`),
					},
				},
			})
		},
		"MobileOtpFailureCount": func(t *testing.T) {
			t.Parallel()

			resourceName := acctest.ResourceNameGen()
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
					{
						Config:      testAccMFADevicePolicyDefaultConfig_MobileOtpFailureCount(resourceName, name, 0),
						ExpectError: regexp.MustCompile(`Attribute mobile.otp.failure.count value must be between 1 and 7`),
					},
					{
						Config:      testAccMFADevicePolicyDefaultConfig_MobileOtpFailureCount(resourceName, name, 8),
						ExpectError: regexp.MustCompile(`Attribute mobile.otp.failure.count value must be between 1 and 7`),
					},
				},
			})
		},
		"MobileNewRequestDuration_Low": func(t *testing.T) {
			t.Parallel()

			resourceName := acctest.ResourceNameGen()
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
					{
						Config:      testAccMFADevicePolicyDefaultConfig_MobileNewRequestDuration_DeviceTimeout(resourceName, name, 14),
						ExpectError: regexp.MustCompile(`Attribute[\s\n]+mobile.applications\[0\].new_request_duration_configuration.device_timeout.duration[\s\n]+value must be between 15 and 75`),
					},
					{
						Config:      testAccMFADevicePolicyDefaultConfig_MobileNewRequestDuration_TotalTimeout(resourceName, name, 29),
						ExpectError: regexp.MustCompile(`Attribute[\s\n]+mobile.applications\[0\].new_request_duration_configuration.total_timeout.duration[\s\n]+value must be between 30 and 90`),
					},
				},
			})
		},
		"MobileIntegrityDetection": func(t *testing.T) {
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
					{
						Config:      testAccMFADevicePolicyDefaultConfig_MobileIntegrityDetection(environmentName, licenseID, resourceName, name, "INVALID_VALUE"),
						ExpectError: regexp.MustCompile(`Attribute mobile.applications\[0\].integrity_detection value must be one of:`),
					},
				},
			})
		},
	}

	for name, testFunc := range testCases {
		t.Run(name, testFunc)
	}
}

func testAccMFADevicePolicyDefaultConfig_MobileIntegrityDetection(environmentName, licenseID, resourceName, name, integrityDetection string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_mfa_device_policy_default" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  policy_type    = "pingone_mfa"
  name           = "%[4]s"

  mobile = {
    enabled = true
    applications = [{
      id = "11111111-1111-1111-1111-111111111111"
      otp = { enabled = true }
      auto_enrollment = { enabled = true }
      device_authorization = { enabled = true }
      integrity_detection = "%[5]s"
    }]
  }
  sms = { enabled = false }
  voice = { enabled = false }
  email = { enabled = false }
  totp = { enabled = false }
}`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name, integrityDetection)
}

func testAccMFADevicePolicyDefaultConfig_PingID_MissingDesktop(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_mfa_device_policy_default" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  policy_type    = "pingid"
  name           = "%[3]s"

  mobile = { enabled = true }
  yubikey = { enabled = false }
  sms = { enabled = false }
  voice = { enabled = false }
  email = { enabled = false }
  totp = { enabled = false }
}`, acctest.WorkforceV2SandboxEnvironment(), resourceName, name)
}

func testAccMFADevicePolicyDefaultConfig_PingID_MissingYubikey(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_mfa_device_policy_default" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  policy_type    = "pingid"
  name           = "%[3]s"

  mobile = { enabled = true }
  desktop = { enabled = false }
  sms = { enabled = false }
  voice = { enabled = false }
  email = { enabled = false }
  totp = { enabled = false }
}`, acctest.WorkforceV2SandboxEnvironment(), resourceName, name)
}

func testAccMFADevicePolicyDefaultConfig_PingID_MobileApp_AutoEnrollment(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_mfa_device_policy_default" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  policy_type    = "pingid"
  name           = "%[3]s"

  mobile = {
    enabled = true
    applications = [{
      id = "11111111-1111-1111-1111-111111111111"
      otp = { enabled = true }
      auto_enrollment = { enabled = true }
      new_request_duration_configuration = {
        device_timeout = { duration = 25 }
        total_timeout = { duration = 40 }
      }
      ip_pairing_configuration = { any_ip_address = true }
    }]
  }
  desktop = { enabled = false }
  yubikey = { enabled = false }
  sms = { enabled = false }
  voice = { enabled = false }
  email = { enabled = false }
  totp = { enabled = false }
}`, acctest.WorkforceV2SandboxEnvironment(), resourceName, name)
}

func testAccMFADevicePolicyDefaultConfig_PingID_MobileApp_DeviceAuthorization(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_mfa_device_policy_default" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  policy_type    = "pingid"
  name           = "%[3]s"

  mobile = {
    enabled = true
    applications = [{
      id = "11111111-1111-1111-1111-111111111111"
      otp = { enabled = true }
      device_authorization = { enabled = true }
      new_request_duration_configuration = {
        device_timeout = { duration = 25 }
        total_timeout = { duration = 40 }
      }
      ip_pairing_configuration = { any_ip_address = true }
    }]
  }
  desktop = { enabled = false }
  yubikey = { enabled = false }
  sms = { enabled = false }
  voice = { enabled = false }
  email = { enabled = false }
  totp = { enabled = false }
}`, acctest.WorkforceV2SandboxEnvironment(), resourceName, name)
}

func testAccMFADevicePolicyDefaultConfig_PingID_MobileApp_MissingNewRequestDuration(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_mfa_device_policy_default" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  policy_type    = "pingid"
  name           = "%[3]s"

  mobile = {
    enabled = true
    applications = [{
      id = "11111111-1111-1111-1111-111111111111"
      otp = { enabled = true }
      ip_pairing_configuration = { any_ip_address = true }
    }]
  }
  desktop = { enabled = false }
  yubikey = { enabled = false }
  sms = { enabled = false }
  voice = { enabled = false }
  email = { enabled = false }
  totp = { enabled = false }
}`, acctest.WorkforceV2SandboxEnvironment(), resourceName, name)
}

func testAccMFADevicePolicyDefaultConfig_PingID_MobileApp_MissingIPPairing(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_mfa_device_policy_default" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  policy_type    = "pingid"
  name           = "%[3]s"

  mobile = {
    enabled = true
    applications = [{
      id = "11111111-1111-1111-1111-111111111111"
      otp = { enabled = true }
      new_request_duration_configuration = {
        device_timeout = { duration = 25 }
        total_timeout = { duration = 40 }
      }
    }]
  }
  desktop = { enabled = false }
  yubikey = { enabled = false }
  sms = { enabled = false }
  voice = { enabled = false }
  email = { enabled = false }
  totp = { enabled = false }
}`, acctest.WorkforceV2SandboxEnvironment(), resourceName, name)
}

func testAccMFADevicePolicyDefaultConfig_PingOneMFA_MobileApp_MissingAutoEnrollment(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_mfa_device_policy_default" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  policy_type    = "pingone_mfa"
  name           = "%[4]s"

  mobile = {
    enabled = true
    applications = [{
      id = "11111111-1111-1111-1111-111111111111"
      otp = { enabled = true }
      device_authorization = { enabled = true }
      integrity_detection = "permissive"
    }]
  }
  sms = { enabled = false }
  voice = { enabled = false }
  email = { enabled = false }
  totp = { enabled = false }
}`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccMFADevicePolicyDefaultConfig_PingOneMFA_MobileApp_MissingDeviceAuthorization(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_mfa_device_policy_default" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  policy_type    = "pingone_mfa"
  name           = "%[4]s"

  mobile = {
    enabled = true
    applications = [{
      id = "11111111-1111-1111-1111-111111111111"
      otp = { enabled = true }
      auto_enrollment = { enabled = true }
      integrity_detection = "permissive"
    }]
  }
  sms = { enabled = false }
  voice = { enabled = false }
  email = { enabled = false }
  totp = { enabled = false }
}`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccMFADevicePolicyDefaultConfig_PingOneMFA_MobileApp_MissingIntegrityDetection(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_mfa_device_policy_default" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  policy_type    = "pingone_mfa"
  name           = "%[4]s"

  mobile = {
    enabled = true
    applications = [{
      id = "11111111-1111-1111-1111-111111111111"
      otp = { enabled = true }
      auto_enrollment = { enabled = true }
      device_authorization = { enabled = true }
    }]
  }
  sms = { enabled = false }
  voice = { enabled = false }
  email = { enabled = false }
  totp = { enabled = false }
}`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccMFADevicePolicyDefaultConfig_PingOneMFA_MobileApp_BiometricsEnabled(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_mfa_device_policy_default" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  policy_type    = "pingone_mfa"
  name           = "%[4]s"

  mobile = {
    enabled = true
    applications = [{
      id = "11111111-1111-1111-1111-111111111111"
      otp = { enabled = true }
      auto_enrollment = { enabled = true }
      device_authorization = { enabled = true }
      integrity_detection = "permissive"
      biometrics_enabled = true
    }]
  }
  sms = { enabled = false }
  voice = { enabled = false }
  email = { enabled = false }
  totp = { enabled = false }
}`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccMFADevicePolicyDefaultConfig_PingOneMFA_MobileApp_NewRequestDuration(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_mfa_device_policy_default" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  policy_type    = "pingone_mfa"
  name           = "%[4]s"

  mobile = {
    enabled = true
    applications = [{
      id = "11111111-1111-1111-1111-111111111111"
      otp = { enabled = true }
      auto_enrollment = { enabled = true }
      device_authorization = { enabled = true }
      integrity_detection = "permissive"
      new_request_duration_configuration = {
        device_timeout = { duration = 25 }
        total_timeout = { duration = 40 }
      }
    }]
  }
  sms = { enabled = false }
  voice = { enabled = false }
  email = { enabled = false }
  totp = { enabled = false }
}`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccMFADevicePolicyDefaultConfig_PingOneMFA_MobileApp_IPPairing(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_mfa_device_policy_default" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  policy_type    = "pingone_mfa"
  name           = "%[4]s"

  mobile = {
    enabled = true
    applications = [{
      id = "11111111-1111-1111-1111-111111111111"
      otp = { enabled = true }
      auto_enrollment = { enabled = true }
      device_authorization = { enabled = true }
      integrity_detection = "permissive"
      ip_pairing_configuration = { any_ip_address = true }
    }]
  }
  sms = { enabled = false }
  voice = { enabled = false }
  email = { enabled = false }
  totp = { enabled = false }
}`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccMFADevicePolicyDefaultConfig_PingID_IPPairing_InvalidCIDR(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_mfa_device_policy_default" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  policy_type    = "pingid"
  name           = "%[3]s"

  mobile = {
    enabled = true
    applications = [{
      id = "11111111-1111-1111-1111-111111111111"
      otp = { enabled = true }
      new_request_duration_configuration = {
        device_timeout = { duration = 25 }
        total_timeout = { duration = 40 }
      }
      ip_pairing_configuration = {
        any_ip_address = false
        only_these_ip_addresses = ["192.168.1.1"]
      }
    }]
  }
  desktop = { enabled = false }
  yubikey = { enabled = false }
  sms = { enabled = false }
  voice = { enabled = false }
  email = { enabled = false }
  totp = { enabled = false }
}`, acctest.WorkforceV2SandboxEnvironment(), resourceName, name)
}

func testAccMFADevicePolicyDefaultConfig_PingID_IPPairing_MissingIPs(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_mfa_device_policy_default" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  policy_type    = "pingid"
  name           = "%[3]s"

  mobile = {
    enabled = true
    applications = [{
      id = "11111111-1111-1111-1111-111111111111"
      otp = { enabled = true }
      new_request_duration_configuration = {
        device_timeout = { duration = 25 }
        total_timeout = { duration = 40 }
      }
      ip_pairing_configuration = {
        any_ip_address = false
      }
    }]
  }
  desktop = { enabled = false }
  yubikey = { enabled = false }
  sms = { enabled = false }
  voice = { enabled = false }
  email = { enabled = false }
  totp = { enabled = false }
}`, acctest.WorkforceV2SandboxEnvironment(), resourceName, name)
}

func testAccMFADevicePolicyDefaultConfig_PingID_Desktop_OTPCountHigh(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_mfa_device_policy_default" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  policy_type    = "pingid"
  name           = "%[3]s"

  mobile = { enabled = true }
  desktop = {
    enabled = true
    otp = {
      failure = {
        count = 8
        cool_down = { duration = 2, time_unit = "MINUTES" }
      }
    }
  }
  yubikey = { enabled = false }
  sms = { enabled = false }
  voice = { enabled = false }
  email = { enabled = false }
  totp = { enabled = false }
}`, acctest.WorkforceV2SandboxEnvironment(), resourceName, name)
}

func testAccMFADevicePolicyDefaultConfig_PingID_Desktop_PairingKeyLifetimeHigh(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_mfa_device_policy_default" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  policy_type    = "pingid"
  name           = "%[3]s"

  mobile = { enabled = true }
  desktop = {
    enabled = true
    pairing_key_lifetime = {
      duration = 50
      time_unit = "HOURS"
    }
  }
  yubikey = { enabled = false }
  sms = { enabled = false }
  voice = { enabled = false }
  email = { enabled = false }
  totp = { enabled = false }
}`, acctest.WorkforceV2SandboxEnvironment(), resourceName, name)
}

func testAccMFADevicePolicyDefaultConfig_RememberMe_MinutesHigh(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_mfa_device_policy_default" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  policy_type    = "pingid"
  name           = "%[3]s"

  remember_me = {
    web = {
      enabled = true
      life_time = {
        duration = 129601
        time_unit = "MINUTES"
      }
    }
  }

  mobile = { enabled = true }
  desktop = { enabled = false }
  yubikey = { enabled = false }
  sms = { enabled = false }
  voice = { enabled = false }
  email = { enabled = false }
  totp = { enabled = false }
}`, acctest.WorkforceV2SandboxEnvironment(), resourceName, name)
}

func testAccMFADevicePolicyDefaultConfig_RememberMe_HoursHigh(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_mfa_device_policy_default" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  policy_type    = "pingid"
  name           = "%[3]s"

  remember_me = {
    web = {
      enabled = true
      life_time = {
        duration = 2161
        time_unit = "HOURS"
      }
    }
  }

  mobile = { enabled = true }
  desktop = { enabled = false }
  yubikey = { enabled = false }
  sms = { enabled = false }
  voice = { enabled = false }
  email = { enabled = false }
  totp = { enabled = false }
}`, acctest.WorkforceV2SandboxEnvironment(), resourceName, name)
}

func testAccMFADevicePolicyDefaultConfig_RememberMe_DaysHigh(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_mfa_device_policy_default" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  policy_type    = "pingid"
  name           = "%[3]s"

  remember_me = {
    web = {
      enabled = true
      life_time = {
        duration = 91
        time_unit = "DAYS"
      }
    }
  }

  mobile = { enabled = true }
  desktop = { enabled = false }
  yubikey = { enabled = false }
  sms = { enabled = false }
  voice = { enabled = false }
  email = { enabled = false }
  totp = { enabled = false }
}`, acctest.WorkforceV2SandboxEnvironment(), resourceName, name)
}

func testAccMFADevicePolicyDefaultConfig_MobilePushLimit_CountHigh(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_mfa_device_policy_default" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  policy_type    = "pingid"
  name           = "%[3]s"

  mobile = {
    enabled = true
    applications = [{
      id = "11111111-1111-1111-1111-111111111111"
      otp = { enabled = true }
      push_limit = {
        count = 51
        lock_duration = { duration = 30, time_unit = "MINUTES" }
        time_period = { duration = 10, time_unit = "MINUTES" }
      }
      new_request_duration_configuration = {
        device_timeout = { duration = 25 }
        total_timeout = { duration = 40 }
      }
      ip_pairing_configuration = { any_ip_address = true }
    }]
  }
  desktop = { enabled = false }
  yubikey = { enabled = false }
  sms = { enabled = false }
  voice = { enabled = false }
  email = { enabled = false }
  totp = { enabled = false }
}`, acctest.WorkforceV2SandboxEnvironment(), resourceName, name)
}

func testAccMFADevicePolicyDefaultConfig_MobileNewRequestDuration_DeviceTimeoutHigh(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_mfa_device_policy_default" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  policy_type    = "pingid"
  name           = "%[3]s"

  mobile = {
    enabled = true
    applications = [{
      id = "11111111-1111-1111-1111-111111111111"
      otp = { enabled = true }
      new_request_duration_configuration = {
        device_timeout = { duration = 76 }
        total_timeout = { duration = 40 }
      }
      ip_pairing_configuration = { any_ip_address = true }
    }]
  }
  desktop = { enabled = false }
  yubikey = { enabled = false }
  sms = { enabled = false }
  voice = { enabled = false }
  email = { enabled = false }
  totp = { enabled = false }
}`, acctest.WorkforceV2SandboxEnvironment(), resourceName, name)
}

func testAccMFADevicePolicyDefaultConfig_MobileNewRequestDuration_TotalTimeoutHigh(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_mfa_device_policy_default" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  policy_type    = "pingid"
  name           = "%[3]s"

  mobile = {
    enabled = true
    applications = [{
      id = "11111111-1111-1111-1111-111111111111"
      otp = { enabled = true }
      new_request_duration_configuration = {
        device_timeout = { duration = 25 }
        total_timeout = { duration = 91 }
      }
      ip_pairing_configuration = { any_ip_address = true }
    }]
  }
  desktop = { enabled = false }
  yubikey = { enabled = false }
  sms = { enabled = false }
  voice = { enabled = false }
  email = { enabled = false }
  totp = { enabled = false }
}`, acctest.WorkforceV2SandboxEnvironment(), resourceName, name)
}

func testAccMFADevicePolicyDefaultConfig_Authentication(resourceName, name, deviceSelection string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_mfa_device_policy_default" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  policy_type    = "pingid"
  name           = "%[3]s"

  authentication = {
    device_selection = "%[4]s"
  }

  mobile = { enabled = true }
  desktop = { enabled = false }
  yubikey = { enabled = false }
  sms = { enabled = false }
  voice = { enabled = false }
  email = { enabled = false }
  totp = { enabled = false }
}`, acctest.WorkforceV2SandboxEnvironment(), resourceName, name, deviceSelection)
}

func testAccMFADevicePolicyDefaultConfig_NewDeviceNotification(resourceName, name, notification string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_mfa_device_policy_default" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  policy_type    = "pingid"
  name           = "%[3]s"

  new_device_notification = "%[4]s"

  mobile = { enabled = true }
  desktop = { enabled = false }
  yubikey = { enabled = false }
  sms = { enabled = false }
  voice = { enabled = false }
  email = { enabled = false }
  totp = { enabled = false }
}`, acctest.WorkforceV2SandboxEnvironment(), resourceName, name, notification)
}

func testAccMFADevicePolicyDefaultConfig_MobileOtpFailureCount(resourceName, name string, count int) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_mfa_device_policy_default" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  policy_type    = "pingid"
  name           = "%[3]s"

  mobile = {
    enabled = true
    otp = {
      failure = {
        count = %[4]d
        cool_down = { duration = 2, time_unit = "MINUTES" }
      }
    }
    applications = [{
      id = "11111111-1111-1111-1111-111111111111"
      otp = { enabled = true }
      new_request_duration_configuration = {
        device_timeout = { duration = 25 }
        total_timeout = { duration = 40 }
      }
      ip_pairing_configuration = { any_ip_address = true }
    }]
  }
  desktop = { enabled = false }
  yubikey = { enabled = false }
  sms = { enabled = false }
  voice = { enabled = false }
  email = { enabled = false }
  totp = { enabled = false }
}`, acctest.WorkforceV2SandboxEnvironment(), resourceName, name, count)
}

func testAccMFADevicePolicyDefaultConfig_MobileNewRequestDuration_DeviceTimeout(resourceName, name string, duration int) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_mfa_device_policy_default" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  policy_type    = "pingid"
  name           = "%[3]s"

  mobile = {
    enabled = true
    applications = [{
      id = "11111111-1111-1111-1111-111111111111"
      otp = { enabled = true }
      new_request_duration_configuration = {
        device_timeout = { duration = %[4]d }
        total_timeout = { duration = 40 }
      }
      ip_pairing_configuration = { any_ip_address = true }
    }]
  }
  desktop = { enabled = false }
  yubikey = { enabled = false }
  sms = { enabled = false }
  voice = { enabled = false }
  email = { enabled = false }
  totp = { enabled = false }
}`, acctest.WorkforceV2SandboxEnvironment(), resourceName, name, duration)
}

func testAccMFADevicePolicyDefaultConfig_MobileNewRequestDuration_TotalTimeout(resourceName, name string, duration int) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_mfa_device_policy_default" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id
  policy_type    = "pingid"
  name           = "%[3]s"

  mobile = {
    enabled = true
    applications = [{
      id = "11111111-1111-1111-1111-111111111111"
      otp = { enabled = true }
      new_request_duration_configuration = {
        device_timeout = { duration = 25 }
        total_timeout = { duration = %[4]d }
      }
      ip_pairing_configuration = { any_ip_address = true }
    }]
  }
  desktop = { enabled = false }
  yubikey = { enabled = false }
  sms = { enabled = false }
  voice = { enabled = false }
  email = { enabled = false }
  totp = { enabled = false }
}`, acctest.WorkforceV2SandboxEnvironment(), resourceName, name, duration)
}
