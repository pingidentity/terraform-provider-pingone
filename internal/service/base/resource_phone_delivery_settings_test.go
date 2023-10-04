package base_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/base"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccPhoneDeliverySettings_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_phone_delivery_settings.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	var phoneDeliverySettingsID, environmentID string

	var p1Client *client.Client
	var ctx = context.Background()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)

			p1Client = acctest.PreCheckTestClient(ctx, t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.PhoneDeliverySettings_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccPhoneDeliverySettingsConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check:  base.PhoneDeliverySettings_GetIDs(resourceFullName, &environmentID, &phoneDeliverySettingsID),
			},
			// Replan after removal preconfig
			{
				PreConfig: func() {
					base.PhoneDeliverySettings_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID, phoneDeliverySettingsID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the environment
			{
				Config: testAccPhoneDeliverySettingsConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check:  base.PhoneDeliverySettings_GetIDs(resourceFullName, &environmentID, &phoneDeliverySettingsID),
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

func TestAccPhoneDeliverySettings_NewEnv(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_phone_delivery_settings.%s", resourceName)

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
		CheckDestroy:             base.PhoneDeliverySettings_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPhoneDeliverySettingsConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
				),
			},
		},
	})
}

func TestAccPhoneDeliverySettings_Custom_Twilio(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_phone_delivery_settings.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	skipTwilio, err := strconv.ParseBool(os.Getenv("PINGONE_TWILIO_TEST_SKIP"))
	if err != nil {
		skipTwilio = false
	}

	twilioSID := os.Getenv("PINGONE_TWILIO_SID")
	twilioAuthToken := os.Getenv("PINGONE_TWILIO_AUTH_TOKEN")
	number := os.Getenv("PINGONE_TWILIO_NUMBER")

	check := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "provider_type", "CUSTOM_TWILIO"),
		resource.TestCheckNoResourceAttr(resourceFullName, "provider_custom"),
		resource.TestCheckResourceAttr(resourceFullName, "provider_custom_twilio.sid", twilioSID),
		resource.TestCheckResourceAttr(resourceFullName, "provider_custom_twilio.auth_token", twilioAuthToken),
		resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "provider_custom_twilio.selected_numbers.*", map[string]string{
			"number":   number,
			"selected": "true",
			"type":     "PHONE_NUMBER",
		}),
		resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "provider_custom_twilio.service_numbers.*", map[string]string{
			"available":      "true",
			"capabilities.#": "2",
			"capabilities.0": "SMS",
			"capabilities.1": "VOICE",
			"number":         number,
			"selected":       "true",
			"type":           "PHONE_NUMBER",
		}),
		resource.TestCheckNoResourceAttr(resourceFullName, "provider_custom_syniverse"),
		resource.TestMatchResourceAttr(resourceFullName, "created_at", verify.RFC3339Regexp),
		resource.TestMatchResourceAttr(resourceFullName, "updated_at", verify.RFC3339Regexp),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckTwilio(t, skipTwilio)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.PhoneDeliverySettings_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPhoneDeliverySettingsConfig_Custom_Twilio(environmentName, licenseID, resourceName, twilioSID, twilioAuthToken, number),
				Check:  check,
			},
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
					"provider_custom_twilio.auth_token",
					"provider_custom_twilio.selected_numbers.#",
					"provider_custom_twilio.selected_numbers.0.%",
					"provider_custom_twilio.selected_numbers.0.number",
					"provider_custom_twilio.selected_numbers.0.selected",
					"provider_custom_twilio.selected_numbers.0.type",
				},
			},
			{
				Config:  testAccPhoneDeliverySettingsConfig_Custom_Twilio(environmentName, licenseID, resourceName, twilioSID, twilioAuthToken, number),
				Destroy: true,
			},
			// Errors
			{
				Config:      testAccPhoneDeliverySettingsConfig_Custom_Twilio(environmentName, licenseID, resourceName, "unknownsid", twilioAuthToken, number),
				ExpectError: regexp.MustCompile(`Authentication error`),
			},
			{
				Config:      testAccPhoneDeliverySettingsConfig_Custom_Twilio(environmentName, licenseID, resourceName, twilioSID, "unknownauthtoken", number),
				ExpectError: regexp.MustCompile(`Authentication error`),
			},
		},
	})
}

func TestAccPhoneDeliverySettings_Custom_Syniverse(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_phone_delivery_settings.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	skipSyniverse, err := strconv.ParseBool(os.Getenv("PINGONE_SYNIVERSE_TEST_SKIP"))
	if err != nil {
		skipSyniverse = false
	}

	syniverseAuthToken := os.Getenv("PINGONE_SYNIVERSE_AUTH_TOKEN")
	number := os.Getenv("PINGONE_SYNIVERSE_NUMBER")

	check := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "provider_type", "CUSTOM_SYNIVERSE"),
		resource.TestCheckNoResourceAttr(resourceFullName, "provider_custom"),
		resource.TestCheckNoResourceAttr(resourceFullName, "provider_custom_twilio"),
		resource.TestCheckResourceAttr(resourceFullName, "provider_custom_syniverse.auth_token", syniverseAuthToken),
		resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "provider_custom_syniverse.selected_numbers.*", map[string]string{
			"number":   number,
			"selected": "true",
			"type":     "PHONE_NUMBER",
		}),
		resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "provider_custom_syniverse.service_numbers.*", map[string]string{
			"available":      "true",
			"capabilities.#": "2",
			"capabilities.0": "SMS",
			"capabilities.1": "VOICE",
			"number":         number,
			"selected":       "true",
			"type":           "PHONE_NUMBER",
		}),
		resource.TestMatchResourceAttr(resourceFullName, "created_at", verify.RFC3339Regexp),
		resource.TestMatchResourceAttr(resourceFullName, "updated_at", verify.RFC3339Regexp),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckSyniverse(t, skipSyniverse)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.PhoneDeliverySettings_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPhoneDeliverySettingsConfig_Custom_Syniverse(environmentName, licenseID, resourceName, syniverseAuthToken, number),
				Check:  check,
			},
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
					"provider_custom_syniverse.auth_token",
					"provider_custom_syniverse.selected_numbers.#",
					"provider_custom_syniverse.selected_numbers.0.%",
					"provider_custom_syniverse.selected_numbers.0.number",
					"provider_custom_syniverse.selected_numbers.0.selected",
					"provider_custom_syniverse.selected_numbers.0.type",
				},
			},
			{
				Config:  testAccPhoneDeliverySettingsConfig_Custom_Syniverse(environmentName, licenseID, resourceName, syniverseAuthToken, number),
				Destroy: true,
			},
			// Errors
			{
				Config:      testAccPhoneDeliverySettingsConfig_Custom_Syniverse(environmentName, licenseID, resourceName, "unknownauthtoken", number),
				ExpectError: regexp.MustCompile(`Authentication Error`),
			},
		},
	})
}

func TestAccPhoneDeliverySettings_Custom(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_phone_delivery_settings.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	fullCheck := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "provider_type", "CUSTOM_PROVIDER"),
		resource.TestCheckResourceAttr(resourceFullName, "provider_custom.name", name),

		resource.TestCheckResourceAttr(resourceFullName, "provider_custom.authentication.method", "BASIC"),
		resource.TestCheckResourceAttr(resourceFullName, "provider_custom.authentication.username", "testusername"),
		resource.TestCheckResourceAttr(resourceFullName, "provider_custom.authentication.password", "testpassword"),
		resource.TestCheckNoResourceAttr(resourceFullName, "provider_custom.authentication.auth_token"),

		resource.TestCheckResourceAttr(resourceFullName, "provider_custom.numbers.#", "3"),
		resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "provider_custom.numbers.*", map[string]string{
			"available":             "true",
			"capabilities.#":        "2",
			"capabilities.0":        "SMS",
			"capabilities.1":        "VOICE",
			"number":                "+441234567890",
			"selected":              "true",
			"supported_countries.#": "4",
			"supported_countries.0": "DE",
			"supported_countries.1": "FR",
			"supported_countries.2": "GB",
			"supported_countries.3": "US",
			"type":                  "TOLL_FREE",
		}),
		resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "provider_custom.numbers.*", map[string]string{
			"available":             "false",
			"capabilities.#":        "1",
			"capabilities.0":        "VOICE",
			"number":                "+441234567891",
			"supported_countries.#": "1",
			"supported_countries.0": "US",
			"selected":              "false",
			"type":                  "SHORT_CODE",
		}),
		resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "provider_custom.numbers.*", map[string]string{
			"available":      "false",
			"capabilities.#": "1",
			"capabilities.0": "SMS",
			"number":         "+441234567892",
			"selected":       "false",
			"type":           "PHONE_NUMBER",
		}),

		resource.TestCheckResourceAttr(resourceFullName, "provider_custom.requests.#", "4"),
		resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "provider_custom.requests.*", map[string]string{
			"body":                 "{\"from\":\"${from}\",\"message\":\"${message}\",\"to\":\"${to}\"}",
			"delivery_method":      "SMS",
			"headers.%":            "2",
			"headers.content-type": "application/json",
			"headers.testheader":   "testvalue1",
			"method":               "POST",
			"phone_number_format":  "FULL",
			"url":                  "https://pingdevops.com/fake-send-to-test",
		}),
		resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "provider_custom.requests.*", map[string]string{
			"delivery_method":     "SMS",
			"method":              "GET",
			"headers.%":           "1",
			"headers.testheader":  "testvalue2",
			"phone_number_format": "NUMBER_ONLY",
			"url":                 "https://pingdevops.com/fake-send-to-test?to=${to}&from=${from}&message=${message}",
		}),
		resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "provider_custom.requests.*", map[string]string{
			"after_tag":            "</Say> <Pause length=\"1\"/>",
			"before_tag":           "<Say>",
			"body":                 "{\"from\":\"${from}\",\"message\":\"${message}\",\"to\":\"${to}\"}",
			"delivery_method":      "VOICE",
			"headers.%":            "2",
			"headers.content-type": "application/json",
			"headers.testheader":   "testvalue3",
			"method":               "POST",
			"phone_number_format":  "FULL",
			"url":                  "https://pingdevops.com/fake-send-to-test",
		}),
		resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "provider_custom.requests.*", map[string]string{
			"after_tag":           "</Say> <Pause length=\"1\"/>",
			"before_tag":          "<Say>",
			"delivery_method":     "VOICE",
			"headers.%":           "1",
			"headers.testheader":  "testvalue4",
			"method":              "GET",
			"phone_number_format": "NUMBER_ONLY",
			"url":                 "https://pingdevops.com/fake-send-to-test?to=${to}&from=${from}&message=${message}",
		}),

		resource.TestCheckNoResourceAttr(resourceFullName, "provider_custom_twilio"),
		resource.TestCheckNoResourceAttr(resourceFullName, "provider_custom_syniverse"),
		resource.TestMatchResourceAttr(resourceFullName, "created_at", verify.RFC3339Regexp),
		resource.TestMatchResourceAttr(resourceFullName, "updated_at", verify.RFC3339Regexp),
	)

	minimalCheck := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "provider_type", "CUSTOM_PROVIDER"),
		resource.TestCheckResourceAttr(resourceFullName, "provider_custom.name", name),

		resource.TestCheckResourceAttr(resourceFullName, "provider_custom.authentication.method", "BEARER"),
		resource.TestCheckNoResourceAttr(resourceFullName, "provider_custom.authentication.username"),
		resource.TestCheckNoResourceAttr(resourceFullName, "provider_custom.authentication.password"),
		resource.TestCheckResourceAttr(resourceFullName, "provider_custom.authentication.auth_token", "testtoken"),

		resource.TestCheckNoResourceAttr(resourceFullName, "provider_custom.numbers"),

		resource.TestCheckResourceAttr(resourceFullName, "provider_custom.requests.#", "1"),
		resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "provider_custom.requests.*", map[string]string{
			"delivery_method": "SMS",
			"method":          "GET",
			"url":             "https://pingdevops.com/fake-send-to-test?to=${to}&message=${message}",
		}),

		resource.TestCheckNoResourceAttr(resourceFullName, "provider_custom_twilio"),
		resource.TestCheckNoResourceAttr(resourceFullName, "provider_custom_syniverse"),
		resource.TestMatchResourceAttr(resourceFullName, "created_at", verify.RFC3339Regexp),
		resource.TestMatchResourceAttr(resourceFullName, "updated_at", verify.RFC3339Regexp),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.PhoneDeliverySettings_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			{
				Config: testAccPhoneDeliverySettingsConfig_Custom_Full(environmentName, licenseID, resourceName, name),
				Check:  fullCheck,
			},
			{
				Config:  testAccPhoneDeliverySettingsConfig_Custom_Full(environmentName, licenseID, resourceName, name),
				Destroy: true,
			},
			// Minimal
			{
				Config: testAccPhoneDeliverySettingsConfig_Custom_Minimal(environmentName, licenseID, resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config:  testAccPhoneDeliverySettingsConfig_Custom_Minimal(environmentName, licenseID, resourceName, name),
				Destroy: true,
			},
			// update
			{
				Config: testAccPhoneDeliverySettingsConfig_Custom_Full(environmentName, licenseID, resourceName, name),
				Check:  fullCheck,
			},
			{
				Config: testAccPhoneDeliverySettingsConfig_Custom_Minimal(environmentName, licenseID, resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config: testAccPhoneDeliverySettingsConfig_Custom_Full(environmentName, licenseID, resourceName, name),
				Check:  fullCheck,
			},
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
					"provider_custom.authentication.password",
				},
			},
		},
	})
}

func TestAccPhoneDeliverySettings_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_phone_delivery_settings.%s", resourceName)

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
		CheckDestroy:             base.PhoneDeliverySettings_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccPhoneDeliverySettingsConfig_NewEnv(environmentName, licenseID, resourceName, name),
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

func testAccPhoneDeliverySettingsConfig_NewEnv(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_phone_delivery_settings" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  provider_custom = {
    name = "%[4]s"

    authentication = {
      method     = "BEARER"
      auth_token = "testtoken"
    }

    requests = [
      {
        delivery_method     = "SMS"
        method              = "GET"
        phone_number_format = "FULL"
        url                 = "https://pingdevops.com/fake-send-to-test?to=$${to}&message=$${message}"
      }
    ]
  }
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccPhoneDeliverySettingsConfig_Custom_Twilio(environmentName, licenseID, resourceName, twilioSID, twilioAuthToken, number string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_phone_delivery_settings" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  provider_custom_twilio = {
    sid        = "%[4]s"
    auth_token = "%[5]s"

    selected_numbers = [
      {
        number = "%[6]s"
        type   = "PHONE_NUMBER"
      }
    ]
  }
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, twilioSID, twilioAuthToken, number)
}

func testAccPhoneDeliverySettingsConfig_Custom_Syniverse(environmentName, licenseID, resourceName, syniverseAuthToken, number string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_phone_delivery_settings" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  provider_custom_syniverse = {
    auth_token = "%[4]s"

    numbers = [
      {
        number   = "%[5]s"
        selected = true
      }
    ]
  }
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, syniverseAuthToken, number)
}

func testAccPhoneDeliverySettingsConfig_Custom_Full(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_phone_delivery_settings" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  provider_custom = {
    name = "%[4]s"

    authentication = {
      method   = "BASIC"
      username = "testusername"
      password = "testpassword"
    }

    numbers = [
      {
        available           = "true"
        capabilities        = ["VOICE", "SMS"]
        number              = "+441234567890"
        selected            = "true"
        supported_countries = ["US", "FR", "GB", "DE"]
        type                = "TOLL_FREE"
      },

      {
        capabilities        = ["VOICE"]
        number              = "+441234567891"
        supported_countries = ["US"]
        type                = "SHORT_CODE"
      },

      {
        available    = "false"
        capabilities = ["SMS"]
        number       = "+441234567892"
        selected     = "false"
        type         = "PHONE_NUMBER"
      }
    ]

    requests = [
      {
        body = jsonencode({
          "to"      = "$${to}",
          "from"    = "$${from}",
          "message" = "$${message}"
        })
        delivery_method = "SMS"
        headers = {
          "content-type" = "application/json",
          testheader     = "testvalue1",
        }
        method              = "POST"
        phone_number_format = "FULL"
        url                 = "https://pingdevops.com/fake-send-to-test"
      },

      {
        delivery_method = "SMS"
        headers = {
          testheader = "testvalue2",
        }
        method              = "GET"
        phone_number_format = "NUMBER_ONLY"
        url                 = "https://pingdevops.com/fake-send-to-test?to=$${to}&from=$${from}&message=$${message}"
      },

      {
        after_tag  = "</Say> <Pause length=\"1\"/>"
        before_tag = "<Say>"
        body = jsonencode({
          "to"      = "$${to}",
          "from"    = "$${from}",
          "message" = "$${message}"
        })
        delivery_method = "VOICE"
        headers = {
          "content-type" = "application/json",
          testheader     = "testvalue3",
        }
        method              = "POST"
        phone_number_format = "FULL"
        url                 = "https://pingdevops.com/fake-send-to-test"
      },

      {
        after_tag       = "</Say> <Pause length=\"1\"/>"
        before_tag      = "<Say>"
        delivery_method = "VOICE"
        headers = {
          testheader = "testvalue4",
        }
        method              = "GET"
        phone_number_format = "NUMBER_ONLY"
        url                 = "https://pingdevops.com/fake-send-to-test?to=$${to}&from=$${from}&message=$${message}"
      }
    ]
  }
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccPhoneDeliverySettingsConfig_Custom_Minimal(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_phone_delivery_settings" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  provider_custom = {
    name = "%[4]s"

    authentication = {
      method     = "BEARER"
      auth_token = "testtoken"
    }

    requests = [
      {
        delivery_method     = "SMS"
        method              = "GET"
        phone_number_format = "FULL"
        url                 = "https://pingdevops.com/fake-send-to-test?to=$${to}&message=$${message}"
      }
    ]
  }
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}
