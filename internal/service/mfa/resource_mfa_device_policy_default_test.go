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

func TestAccMFADevicePolicyDefault_PingOneMFA_Validation(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

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
				ExpectError: regexp.MustCompile(`Attribute remember_me.web.life_time.duration value must be at most 2160`),
			},
			// Invalid remember_me duration for DAYS
			{
				Config:      testAccMFADevicePolicyDefaultConfig_InvalidRememberMeDurationDays(environmentName, licenseID, resourceName, name),
				ExpectError: regexp.MustCompile(`Attribute remember_me.web.life_time.duration value must be at most 90`),
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

// func TestAccMFADevicePolicyDefault_PingID_Full(t *testing.T) {
// 	t.Parallel()

// 	resourceName := acctest.ResourceNameGen()
// 	resourceFullName := fmt.Sprintf("pingone_mfa_device_policy_default.%s", resourceName)

// 	environmentName := acctest.ResourceNameGenEnvironment()

// 	name := resourceName

// 	licenseID := os.Getenv("PINGONE_LICENSE_ID")

// 	resource.Test(t, resource.TestCase{
// 		PreCheck: func() {
// 			acctest.PreCheckNoTestAccFlaky(t)
// 			acctest.PreCheckClient(t)
// 			acctest.PreCheckNewEnvironment(t)
// 			acctest.PreCheckNoBeta(t)
// 		},
// 		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
// 		CheckDestroy:             mfa.MFADevicePolicyDefault_CheckDestroy,
// 		ErrorCheck:               acctest.ErrorCheck(t),
// 		Steps: []resource.TestStep{
// 			// Step 1: Create environment and let it initialize
// 			{
// 				Config: testAccMFADevicePolicyDefaultConfig_PingID_EnvironmentOnly(environmentName, licenseID),
// 				Check: resource.ComposeTestCheckFunc(
// 					resource.TestMatchResourceAttr(fmt.Sprintf("pingone_environment.%s", environmentName), "id", verify.P1ResourceIDRegexpFullString),
// 				),
// 			},
// 			// Step 2: Now add the device policy after environment has initialized
// 			{
// 				Config: testAccMFADevicePolicyDefaultConfig_PingID_Full(environmentName, licenseID, resourceName, name),
// 				Check: resource.ComposeTestCheckFunc(
// 					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
// 					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
// 					resource.TestCheckResourceAttr(resourceFullName, "policy_type", "pingid"),
// 					resource.TestCheckResourceAttr(resourceFullName, "name", name),
// 					resource.TestCheckResourceAttr(resourceFullName, "mobile.enabled", "true"),
// 					resource.TestCheckResourceAttr(resourceFullName, "desktop.enabled", "true"),
// 					resource.TestCheckResourceAttr(resourceFullName, "desktop.otp.failure.count", "5"),
// 					resource.TestCheckResourceAttr(resourceFullName, "yubikey.enabled", "true"),
// 					resource.TestCheckResourceAttr(resourceFullName, "oath_token.enabled", "true"),
// 					resource.TestCheckResourceAttr(resourceFullName, "oath_token.otp.failure.count", "3"),
// 				),
// 			},
// 		},
// 	})
// }

// func TestAccMFADevicePolicyDefault_PingID_Minimal(t *testing.T) {
// 	t.Parallel()

// 	resourceName := acctest.ResourceNameGen()
// 	resourceFullName := fmt.Sprintf("pingone_mfa_device_policy_default.%s", resourceName)

// 	environmentName := acctest.ResourceNameGenEnvironment()

// 	name := resourceName

// 	licenseID := os.Getenv("PINGONE_LICENSE_ID")

// 	resource.Test(t, resource.TestCase{
// 		PreCheck: func() {
// 			acctest.PreCheckNoTestAccFlaky(t)
// 			acctest.PreCheckClient(t)
// 			acctest.PreCheckNewEnvironment(t)
// 			acctest.PreCheckNoBeta(t)
// 		},
// 		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
// 		CheckDestroy:             mfa.MFADevicePolicyDefault_CheckDestroy,
// 		ErrorCheck:               acctest.ErrorCheck(t),
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccMFADevicePolicyDefaultConfig_PingID_Minimal(environmentName, licenseID, resourceName, name),
// 				Check: resource.ComposeTestCheckFunc(
// 					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
// 					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
// 					resource.TestCheckResourceAttr(resourceFullName, "policy_type", "pingid"),
// 					resource.TestCheckResourceAttr(resourceFullName, "name", name),
// 					resource.TestCheckResourceAttr(resourceFullName, "mobile.enabled", "true"),
// 					resource.TestCheckResourceAttr(resourceFullName, "sms.enabled", "false"),
// 					resource.TestCheckResourceAttr(resourceFullName, "voice.enabled", "false"),
// 					resource.TestCheckResourceAttr(resourceFullName, "email.enabled", "false"),
// 					resource.TestCheckResourceAttr(resourceFullName, "totp.enabled", "false"),
// 				),
// 			},
// 		},
// 	})
// }

// func TestAccMFADevicePolicyDefault_PingID_ValidationErrors(t *testing.T) {
// 	t.Parallel()

// 	resourceName := acctest.ResourceNameGen()

// 	environmentName := acctest.ResourceNameGenEnvironment()

// 	name := resourceName

// 	licenseID := os.Getenv("PINGONE_LICENSE_ID")

// 	resource.Test(t, resource.TestCase{
// 		PreCheck: func() {
// 			acctest.PreCheckNoTestAccFlaky(t)
// 			acctest.PreCheckClient(t)
// 			acctest.PreCheckNewEnvironment(t)
// 			acctest.PreCheckNoBeta(t)
// 		},
// 		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
// 		CheckDestroy:             mfa.MFADevicePolicyDefault_CheckDestroy,
// 		ErrorCheck:               acctest.ErrorCheck(t),
// 		Steps: []resource.TestStep{
// 			// Desktop should conflict with PingOneMFA policy type
// 			{
// 				Config:      testAccMFADevicePolicyDefaultConfig_PingID_DesktopWithPingOneMFA(environmentName, licenseID, resourceName, name),
// 				ExpectError: regexp.MustCompile(`Invalid argument combination`),
// 			},
// 			// Yubikey should conflict with PingOneMFA policy type
// 			{
// 				Config:      testAccMFADevicePolicyDefaultConfig_PingID_YubikeyWithPingOneMFA(environmentName, licenseID, resourceName, name),
// 				ExpectError: regexp.MustCompile(`Invalid argument combination`),
// 			},
// 			// Mobile must be enabled for PingID policies
// 			{
// 				Config:      testAccMFADevicePolicyDefaultConfig_PingID_MobileDisabled(environmentName, licenseID, resourceName, name),
// 				ExpectError: regexp.MustCompile(`Attribute mobile.enabled must be true when attribute policy_type value is`),
// 			},
// 		},
// 	})
// }

// func testAccMFADevicePolicyDefaultConfig_PingID_EnvironmentOnly(environmentName, licenseID string) string {
// 	return acctestlegacysdk.MinimalPingIDSandboxEnvironment(environmentName, licenseID)
// }

// func testAccMFADevicePolicyDefaultConfig_PingID_Full(environmentName, licenseID, resourceName, name string) string {
// 	return fmt.Sprintf(`
// 		%[1]s

// resource "pingone_mfa_device_policy_default" "%[3]s" {
//   environment_id = pingone_environment.%[2]s.id
//   policy_type    = "pingid"

//   depends_on = [pingone_population_default.%[2]s]

//   name = "%[3]s"

//   authentication = {
//     device_selection = "PROMPT_TO_SELECT"
//   }

//   new_device_notification = "SMS_THEN_EMAIL"

//   sms = {
//     enabled = true
//   }

//   voice = {
//     enabled = true
//   }

//   email = {
//     enabled = true
//   }

//   mobile = {
//     enabled = true
//     otp = {
//       failure = {
//         count = 3
//         cool_down = {
//           duration  = 2
//           time_unit = "MINUTES"
//         }
//       }
//     }
//   }

//   totp = {
//     enabled = true
//   }

//   desktop = {
//     enabled = true
//     otp = {
//       failure = {
//         count = 5
//         cool_down = {
//           duration  = 3
//           time_unit = "MINUTES"
//         }
//       }
//     }
//   }

//   yubikey = {
//     enabled = true
//   }

//   oath_token = {
//     enabled = true
//     otp = {
//       failure = {
//         count = 3
//         cool_down = {
//           duration  = 2
//           time_unit = "MINUTES"
//         }
//       }
//     }
//   }
// }`, acctestlegacysdk.MinimalPingIDSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
// }

// func testAccMFADevicePolicyDefaultConfig_PingID_Minimal(environmentName, licenseID, resourceName, name string) string {
// 	return fmt.Sprintf(`
// 		%[1]s

// resource "pingone_mfa_device_policy_default" "%[3]s" {
//   environment_id = pingone_environment.%[2]s.id
//   policy_type    = "pingid"

//   name = "%[3]s"

//   sms = {
//     enabled = false
//   }

//   voice = {
//     enabled = false
//   }

//   email = {
//     enabled = true
//   }

//   mobile = {
//     enabled = true
//   }

//   totp = {
//     enabled = true
//   }
// }`, acctestlegacysdk.MinimalPingIDSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
// }

// func testAccMFADevicePolicyDefaultConfig_PingID_DesktopWithPingOneMFA(environmentName, licenseID, resourceName, name string) string {
// 	return fmt.Sprintf(`
// 		%[1]s

// resource "pingone_mfa_device_policy_default" "%[3]s" {
//   environment_id = pingone_environment.%[2]s.id
//   policy_type    = "pingone_mfa"

//   name = "%[3]s"

//   sms = {
//     enabled = false
//   }

//   voice = {
//     enabled = false
//   }

//   email = {
//     enabled = false
//   }

//   mobile = {
//     enabled = false
//   }

//   totp = {
//     enabled = false
//   }

//   fido2 = {
//     enabled = false
//   }

//   desktop = {
//     enabled = true
//   }
// }`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
// }

// func testAccMFADevicePolicyDefaultConfig_PingID_YubikeyWithPingOneMFA(environmentName, licenseID, resourceName, name string) string {
// 	return fmt.Sprintf(`
// 		%[1]s

// resource "pingone_mfa_device_policy_default" "%[3]s" {
//   environment_id = pingone_environment.%[2]s.id
//   policy_type    = "pingone_mfa"

//   name = "%[3]s"

//   sms = {
//     enabled = false
//   }

//   voice = {
//     enabled = false
//   }

//   email = {
//     enabled = false
//   }

//   mobile = {
//     enabled = false
//   }

//   totp = {
//     enabled = false
//   }

//   fido2 = {
//     enabled = false
//   }

//   yubikey = {
//     enabled = true
//   }
// }`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
// }

// func testAccMFADevicePolicyDefaultConfig_PingID_MobileDisabled(environmentName, licenseID, resourceName, name string) string {
// 	return fmt.Sprintf(`
// 		%[1]s

// resource "pingone_mfa_device_policy_default" "%[3]s" {
//   environment_id = pingone_environment.%[2]s.id
//   policy_type    = "pingid"

//   name = "%[3]s"

//   sms = {
//     enabled = false
//   }

//   voice = {
//     enabled = false
//   }

//   email = {
//     enabled = false
//   }

//   mobile = {
//     enabled = false
//   }

//   totp = {
//     enabled = false
//   }
// }`, acctestlegacysdk.MinimalPingIDSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
// }

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
