package base_test

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
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func testAccCheckPhoneDeliverySettingsDestroy(s *terraform.State) error {
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
		if rs.Type != "pingone_phone_delivery_settings" {
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

		body, r, err := apiClient.PhoneDeliverySettingsApi.ReadOnePhoneDeliverySettings(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.ID).Execute()

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

		return fmt.Errorf("PingOne Phone Delivery Settings %s still exists", rs.Primary.ID)
	}

	return nil
}

func TestAccPhoneDeliverySettings_NewEnv(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_phone_delivery_settings.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckPhoneDeliverySettingsDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPhoneDeliverySettingsConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
				),
			},
		},
	})
}

func TestAccPhoneDeliverySettings_Custom_Twilio(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_phone_delivery_settings.%s", resourceName)

	twilioSID := os.Getenv("PINGONE_TWILIO_SID")
	twilioAuthToken := os.Getenv("PINGONE_TWILIO_AUTH_TOKEN")

	check := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
		resource.TestCheckResourceAttr(resourceFullName, "provider_type", "CUSTOM_TWILIO"),
		resource.TestCheckNoResourceAttr(resourceFullName, "provider_custom"),
		resource.TestCheckResourceAttr(resourceFullName, "provider_custom_twilio.sid", twilioSID),
		resource.TestCheckResourceAttr(resourceFullName, "provider_custom_twilio.auth_token", twilioAuthToken),
		resource.TestCheckNoResourceAttr(resourceFullName, "provider_custom_syniverse"),
		resource.TestMatchResourceAttr(resourceFullName, "created_at", verify.RFC3339Regexp),
		resource.TestMatchResourceAttr(resourceFullName, "updated_at", verify.RFC3339Regexp),
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckPhoneDeliverySettingsDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPhoneDeliverySettingsConfig_Custom_Twilio(resourceName, twilioSID, twilioAuthToken),
				Check:  check,
			},
			{
				Config:  testAccPhoneDeliverySettingsConfig_Custom_Twilio(resourceName, twilioSID, twilioAuthToken),
				Destroy: true,
			},
			// Errors
			{
				Config:      testAccPhoneDeliverySettingsConfig_Custom_Twilio(resourceName, "unknownsid", twilioAuthToken),
				ExpectError: regexp.MustCompile(`uhhm, that didn't work`),
			},
			{
				Config:      testAccPhoneDeliverySettingsConfig_Custom_Twilio(resourceName, twilioSID, "unknownauthtoken"),
				ExpectError: regexp.MustCompile(`uhhm, that didn't work`),
			},
		},
	})
}

func TestAccPhoneDeliverySettings_Custom_Syniverse(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_phone_delivery_settings.%s", resourceName)

	syniverseAuthToken := os.Getenv("PINGONE_SYNIVERSE_AUTH_TOKEN")

	check := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
		resource.TestCheckResourceAttr(resourceFullName, "provider_type", "CUSTOM_SYNIVERSE"),
		resource.TestCheckNoResourceAttr(resourceFullName, "provider_custom"),
		resource.TestCheckNoResourceAttr(resourceFullName, "provider_custom_twilio"),
		resource.TestCheckResourceAttr(resourceFullName, "provider_custom_syniverse.auth_token", syniverseAuthToken),
		resource.TestMatchResourceAttr(resourceFullName, "created_at", verify.RFC3339Regexp),
		resource.TestMatchResourceAttr(resourceFullName, "updated_at", verify.RFC3339Regexp),
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckPhoneDeliverySettingsDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPhoneDeliverySettingsConfig_Custom_Syniverse(resourceName, syniverseAuthToken),
				Check:  check,
			},
			{
				Config:  testAccPhoneDeliverySettingsConfig_Custom_Syniverse(resourceName, syniverseAuthToken),
				Destroy: true,
			},
			// Errors
			{
				Config:      testAccPhoneDeliverySettingsConfig_Custom_Syniverse(resourceName, "unknownauthtoken"),
				ExpectError: regexp.MustCompile(`uhhm, that didn't work`),
			},
		},
	})
}

func TestAccPhoneDeliverySettings_Custom(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_phone_delivery_settings.%s", resourceName)

	name := resourceFullName

	fullCheck := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
		resource.TestCheckResourceAttr(resourceFullName, "provider_type", "CUSTOM_PROVIDER"),
		resource.TestCheckResourceAttr(resourceFullName, "provider_custom.name", name),

		resource.TestCheckResourceAttr(resourceFullName, "provider_custom.authentication.method", "BASIC"),
		resource.TestCheckResourceAttr(resourceFullName, "provider_custom.authentication.username", "testusername"),
		resource.TestCheckResourceAttr(resourceFullName, "provider_custom.authentication.password", "testpassword"),
		resource.TestCheckNoResourceAttr(resourceFullName, "provider_custom.authentication.auth_token"),

		resource.TestCheckResourceAttr(resourceFullName, "provider_custom.numbers.#", "3"),
		resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "provider_custom.numbers.*", map[string]string{
			"available":           "true",
			"capabilities":        "", //SMS,VOICE
			"number":              "+441234567890",
			"selected":            "true",
			"supported_countries": "",
			"type":                "TOLL_FREE",
		}),
		resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "provider_custom.numbers.*", map[string]string{
			"capabilities":        "", //SMS,VOICE
			"number":              "+441234567891",
			"supported_countries": "",
			"type":                "SHORT_CODE",
		}),
		resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "provider_custom.numbers.*", map[string]string{
			"available":           "false",
			"capabilities":        "", //SMS,VOICE
			"number":              "+441234567892",
			"selected":            "false",
			"supported_countries": "",
			"type":                "PHONE_NUMBER",
		}),

		resource.TestCheckResourceAttr(resourceFullName, "provider_custom.requests.#", "4"),
		resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "provider_custom.requests.*", map[string]string{
			"body":                "{\"to\": \"${to}\", \"message\": \"${message}\"}",
			"delivery_method":     "SMS",
			"headers":             "",
			"method":              "POST",
			"phone_number_format": "FULL",
			"url":                 "https://pingdevops.com/fake-send-to-test",
		}),
		resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "provider_custom.requests.*", map[string]string{
			"delivery_method":     "SMS",
			"headers":             "",
			"method":              "GET",
			"phone_number_format": "NUMBER_ONLY",
			"url":                 "https://pingdevops.com/fake-send-to-test?to=${to}&message=${message}",
		}),
		resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "provider_custom.requests.*", map[string]string{
			"after_tag":           "</Say> <Pause length=\"1\"/>",
			"before_tag":          "<Say>",
			"body":                "{\"to\": \"${to}\", \"message\": \"${message}\"}",
			"delivery_method":     "VOICE",
			"headers":             "",
			"method":              "POST",
			"phone_number_format": "FULL",
			"url":                 "https://pingdevops.com/fake-send-to-test",
		}),
		resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "provider_custom.requests.*", map[string]string{
			"after_tag":           "</Say> <Pause length=\"1\"/>",
			"before_tag":          "<Say>",
			"delivery_method":     "VOICE",
			"headers":             "",
			"method":              "GET",
			"phone_number_format": "NUMBER_ONLY",
			"url":                 "https://pingdevops.com/fake-send-to-test?to=${to}&message=${message}",
		}),

		resource.TestCheckNoResourceAttr(resourceFullName, "provider_custom_twilio"),
		resource.TestCheckNoResourceAttr(resourceFullName, "provider_custom_syniverse"),
		resource.TestMatchResourceAttr(resourceFullName, "created_at", verify.RFC3339Regexp),
		resource.TestMatchResourceAttr(resourceFullName, "updated_at", verify.RFC3339Regexp),
	)

	minimalCheck := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
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
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckPhoneDeliverySettingsDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			{
				Config: testAccPhoneDeliverySettingsConfig_Custom_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config:  testAccPhoneDeliverySettingsConfig_Custom_Full(resourceName, name),
				Destroy: true,
			},
			// Minimal
			{
				Config: testAccPhoneDeliverySettingsConfig_Custom_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config:  testAccPhoneDeliverySettingsConfig_Custom_Minimal(resourceName, name),
				Destroy: true,
			},
			// update
			{
				Config: testAccPhoneDeliverySettingsConfig_Custom_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config: testAccPhoneDeliverySettingsConfig_Custom_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config: testAccPhoneDeliverySettingsConfig_Custom_Full(resourceName, name),
				Check:  fullCheck,
			},
		},
	})
}

func testAccPhoneDeliverySettingsConfig_NewEnv(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_phone_delivery_settings" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  provider_type = "CUSTOM_PROVIDER"

  provider_custom = {
    name = "%[4]s"

    authentication = {
      method = "BEARER"
      auth_token  = "testtoken"
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

func testAccPhoneDeliverySettingsConfig_Custom_Twilio(resourceName, twilioSID, twilioAuthToken string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_phone_delivery_settings" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  provider_type = "CUSTOM_TWILIO"

  provider_custom_twilio = {
    sid        = "%[3]s"
    auth_token = "%[4]s"
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, twilioSID, twilioAuthToken)
}

func testAccPhoneDeliverySettingsConfig_Custom_Syniverse(resourceName, syniverseAuthToken string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_phone_delivery_settings" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  provider_type = "CUSTOM_SYNIVERSE"

  provider_custom_syniverse = {
    auth_token = "%[3]s"
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, syniverseAuthToken)
}

func testAccPhoneDeliverySettingsConfig_Custom_Full(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_phone_delivery_settings" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  provider_type = "CUSTOM_PROVIDER"

  provider_custom = {
    name = "%[3]s"

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
        available           = "false"
        capabilities        = ["SMS"]
        number              = "+441234567892"
        selected            = "false"
        supported_countries = ["FR"]
        type                = "PHONE_NUMBER"
      }
    ]

    requests = [
      {
        body = jsonencode({
          "to"      = "$${to}",
          "message" = "$${message}"
        })
        delivery_method = "SMS"
        headers = {
          testheader = "testvalue1",
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
        url                 = "https://pingdevops.com/fake-send-to-test?to=$${to}&message=$${message}"
      },

      {
        after_tag  = "</Say> <Pause length=\"1\"/>"
        before_tag = "<Say>"
        body = jsonencode({
          "to"      = "$${to}",
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
          "content-type" = "application/json",
          testheader     = "testvalue4",
        }
        method              = "GET"
        phone_number_format = "NUMBER_ONLY"
        url                 = "https://pingdevops.com/fake-send-to-test?to=$${to}&message=$${message}"
      }
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccPhoneDeliverySettingsConfig_Custom_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_phone_delivery_settings" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  provider_type = "CUSTOM_PROVIDER"

  provider_custom = {
    name = "%[3]s"

    authentication = {
      method = "BEARER"
      auth_token  = "testtoken"
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
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
