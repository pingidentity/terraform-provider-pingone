// Copyright © 2025 Ping Identity Corporation

package base_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/base"
	validation "github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccNotificationPolicyDataSource_All(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullName := fmt.Sprintf("data.pingone_notification_policy.%s", resourceName)

	name := acctest.ResourceNameGen()

	findByID := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(dataSourceFullName, "id", validation.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", validation.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(dataSourceFullName, "notification_policy_id", validation.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(dataSourceFullName, "name", name),
		resource.TestCheckResourceAttr(dataSourceFullName, "default", "false"),

		// Quota
		resource.TestCheckResourceAttr(dataSourceFullName, "quota.#", "1"),
		resource.TestCheckResourceAttr(dataSourceFullName, "quota.0.type", "ENVIRONMENT"),
		resource.TestCheckResourceAttr(dataSourceFullName, "quota.0.delivery_methods.#", "2"),
		resource.TestCheckResourceAttr(dataSourceFullName, "quota.0.total", "100"),
		resource.TestCheckNoResourceAttr(dataSourceFullName, "quota.0.used"),
		resource.TestCheckNoResourceAttr(dataSourceFullName, "quota.0.unused"),

		// Country Limit
		resource.TestCheckResourceAttr(dataSourceFullName, "country_limit.type", "ALLOWED"),
		resource.TestCheckResourceAttr(dataSourceFullName, "country_limit.delivery_methods.#", "2"),
		resource.TestCheckResourceAttr(dataSourceFullName, "country_limit.countries.#", "2"),

		// Cooldown Configuration - Email
		resource.TestCheckResourceAttr(dataSourceFullName, "cooldown_configuration.email.enabled", "true"),
		resource.TestCheckResourceAttr(dataSourceFullName, "cooldown_configuration.email.group_by", "RECIPIENT"),
		resource.TestCheckResourceAttr(dataSourceFullName, "cooldown_configuration.email.resend_limit", "5"),
		resource.TestCheckResourceAttr(dataSourceFullName, "cooldown_configuration.email.periods.#", "3"),
		resource.TestCheckResourceAttr(dataSourceFullName, "cooldown_configuration.email.periods.0.duration", "10"),
		resource.TestCheckResourceAttr(dataSourceFullName, "cooldown_configuration.email.periods.0.time_unit", "MINUTES"),
		resource.TestCheckResourceAttr(dataSourceFullName, "cooldown_configuration.email.periods.1.duration", "30"),
		resource.TestCheckResourceAttr(dataSourceFullName, "cooldown_configuration.email.periods.1.time_unit", "SECONDS"),
		resource.TestCheckResourceAttr(dataSourceFullName, "cooldown_configuration.email.periods.2.duration", "1"),
		resource.TestCheckResourceAttr(dataSourceFullName, "cooldown_configuration.email.periods.2.time_unit", "MINUTES"),

		// Cooldown Configuration - SMS
		resource.TestCheckResourceAttr(dataSourceFullName, "cooldown_configuration.sms.enabled", "true"),
		resource.TestCheckResourceAttr(dataSourceFullName, "cooldown_configuration.sms.group_by", "RECIPIENT"),
		resource.TestCheckResourceAttr(dataSourceFullName, "cooldown_configuration.sms.resend_limit", "3"),
		resource.TestCheckResourceAttr(dataSourceFullName, "cooldown_configuration.sms.periods.#", "3"),
		resource.TestCheckResourceAttr(dataSourceFullName, "cooldown_configuration.sms.periods.0.duration", "5"),
		resource.TestCheckResourceAttr(dataSourceFullName, "cooldown_configuration.sms.periods.0.time_unit", "MINUTES"),
		resource.TestCheckResourceAttr(dataSourceFullName, "cooldown_configuration.sms.periods.1.duration", "30"),
		resource.TestCheckResourceAttr(dataSourceFullName, "cooldown_configuration.sms.periods.1.time_unit", "SECONDS"),
		resource.TestCheckResourceAttr(dataSourceFullName, "cooldown_configuration.sms.periods.2.duration", "2"),
		resource.TestCheckResourceAttr(dataSourceFullName, "cooldown_configuration.sms.periods.2.time_unit", "MINUTES"),

		// Cooldown Configuration - Voice
		resource.TestCheckResourceAttr(dataSourceFullName, "cooldown_configuration.voice.enabled", "false"),
		resource.TestCheckNoResourceAttr(dataSourceFullName, "cooldown_configuration.voice.group_by"),
		resource.TestCheckNoResourceAttr(dataSourceFullName, "cooldown_configuration.voice.resend_limit"),
		resource.TestCheckNoResourceAttr(dataSourceFullName, "cooldown_configuration.voice.periods"),

		// Cooldown Configuration - WhatsApp
		resource.TestCheckResourceAttr(dataSourceFullName, "cooldown_configuration.whats_app.enabled", "false"),
		resource.TestCheckNoResourceAttr(dataSourceFullName, "cooldown_configuration.whats_app.group_by"),
		resource.TestCheckNoResourceAttr(dataSourceFullName, "cooldown_configuration.whats_app.resend_limit"),
		resource.TestCheckNoResourceAttr(dataSourceFullName, "cooldown_configuration.whats_app.periods"),

		// Provider Configuration
		resource.TestCheckResourceAttr(dataSourceFullName, "provider_configuration.conditions.#", "2"),
		// First condition - with countries specified
		resource.TestCheckResourceAttr(dataSourceFullName, "provider_configuration.conditions.0.delivery_methods.#", "2"),
		resource.TestCheckResourceAttr(dataSourceFullName, "provider_configuration.conditions.0.countries.#", "2"),
		resource.TestCheckResourceAttr(dataSourceFullName, "provider_configuration.conditions.0.fallback_chain.#", "1"),
		resource.TestMatchResourceAttr(dataSourceFullName, "provider_configuration.conditions.0.fallback_chain.0.id", validation.P1ResourceIDRegexpFullString),
		// Second condition - default (no countries)
		resource.TestCheckResourceAttr(dataSourceFullName, "provider_configuration.conditions.1.delivery_methods.#", "2"),
		resource.TestCheckNoResourceAttr(dataSourceFullName, "provider_configuration.conditions.1.countries"),
		resource.TestCheckResourceAttr(dataSourceFullName, "provider_configuration.conditions.1.fallback_chain.#", "1"),
		resource.TestMatchResourceAttr(dataSourceFullName, "provider_configuration.conditions.1.fallback_chain.0.id", validation.P1ResourceIDRegexpFullString),
	)

	findByName := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(dataSourceFullName, "id", validation.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", validation.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(dataSourceFullName, "notification_policy_id", validation.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(dataSourceFullName, "name", name),
		resource.TestCheckResourceAttr(dataSourceFullName, "default", "false"),

		// Quota
		resource.TestCheckResourceAttr(dataSourceFullName, "quota.#", "1"),
		resource.TestCheckResourceAttr(dataSourceFullName, "quota.0.type", "ENVIRONMENT"),
		resource.TestCheckResourceAttr(dataSourceFullName, "quota.0.delivery_methods.#", "2"),
		resource.TestCheckResourceAttr(dataSourceFullName, "quota.0.total", "100"),
		resource.TestCheckNoResourceAttr(dataSourceFullName, "quota.0.used"),
		resource.TestCheckNoResourceAttr(dataSourceFullName, "quota.0.unused"),

		// Country Limit - not set in this test
		resource.TestCheckResourceAttr(dataSourceFullName, "country_limit.type", "NONE"),
		resource.TestCheckNoResourceAttr(dataSourceFullName, "country_limit.delivery_methods"),
		resource.TestCheckNoResourceAttr(dataSourceFullName, "country_limit.countries"),

		// Cooldown Configuration - Email (defaults)
		resource.TestCheckResourceAttr(dataSourceFullName, "cooldown_configuration.email.enabled", "false"),
		resource.TestCheckNoResourceAttr(dataSourceFullName, "cooldown_configuration.email.group_by"),
		resource.TestCheckNoResourceAttr(dataSourceFullName, "cooldown_configuration.email.resend_limit"),
		resource.TestCheckNoResourceAttr(dataSourceFullName, "cooldown_configuration.email.periods"),

		// Cooldown Configuration - SMS (defaults)
		resource.TestCheckResourceAttr(dataSourceFullName, "cooldown_configuration.sms.enabled", "false"),
		resource.TestCheckNoResourceAttr(dataSourceFullName, "cooldown_configuration.sms.group_by"),
		resource.TestCheckNoResourceAttr(dataSourceFullName, "cooldown_configuration.sms.resend_limit"),
		resource.TestCheckNoResourceAttr(dataSourceFullName, "cooldown_configuration.sms.periods"),

		// Cooldown Configuration - Voice (defaults)
		resource.TestCheckResourceAttr(dataSourceFullName, "cooldown_configuration.voice.enabled", "false"),
		resource.TestCheckNoResourceAttr(dataSourceFullName, "cooldown_configuration.voice.group_by"),
		resource.TestCheckNoResourceAttr(dataSourceFullName, "cooldown_configuration.voice.resend_limit"),
		resource.TestCheckNoResourceAttr(dataSourceFullName, "cooldown_configuration.voice.periods"),

		// Cooldown Configuration - WhatsApp (defaults)
		resource.TestCheckResourceAttr(dataSourceFullName, "cooldown_configuration.whats_app.enabled", "false"),
		resource.TestCheckNoResourceAttr(dataSourceFullName, "cooldown_configuration.whats_app.group_by"),
		resource.TestCheckNoResourceAttr(dataSourceFullName, "cooldown_configuration.whats_app.resend_limit"),
		resource.TestCheckNoResourceAttr(dataSourceFullName, "cooldown_configuration.whats_app.periods"),

		// Provider Configuration - not set in this test
		resource.TestCheckNoResourceAttr(dataSourceFullName, "provider_configuration.conditions"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.NotificationPolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNotificationPolicyDataSource_FindByID(resourceName, name),
				Check:  findByID,
			},
			{
				Config:  testAccNotificationPolicyDataSource_FindByID(resourceName, name),
				Destroy: true,
			},
			{
				Config: testAccNotificationPolicyDataSource_FindByName(resourceName, name),
				Check:  findByName,
			},
		},
	})
}

func TestAccNotificationPolicyDataSource_FailureChecks(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.NotificationPolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccNotificationPolicyDataSource_FindByIDFail(resourceName),
				ExpectError: regexp.MustCompile("Notification policy not found"),
			},
			{
				Config:      testAccNotificationPolicyDataSource_FindByNameFail(resourceName),
				ExpectError: regexp.MustCompile("Notification policy not found"),
			},
		},
	})
}

func testAccNotificationPolicyDataSource_FindByID(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_phone_delivery_settings" "%[2]s_provider_1" {
  environment_id = data.pingone_environment.general_test.id

  provider_custom = {
    name = "Test Provider 1"
    authentication = {
      method   = "BASIC"
      username = "test-user-1"
      password = "test-password-1"
    }
    requests = [{
      delivery_method = "SMS"
      url             = "https://example1.com/sms"
      method          = "POST"
    }]
  }
}

resource "pingone_phone_delivery_settings" "%[2]s_provider_2" {
  environment_id = data.pingone_environment.general_test.id

  provider_custom = {
    name = "Test Provider 2"
    authentication = {
      method   = "BASIC"
      username = "test-user-2"
      password = "test-password-2"
    }
    requests = [{
      delivery_method = "SMS"
      url             = "https://example2.com/sms"
      method          = "POST"
    }]
  }
}

resource "pingone_notification_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  quota = [{
    type             = "ENVIRONMENT"
    delivery_methods = ["SMS", "Voice"]
    total            = 100
  }]

  country_limit = {
    type             = "ALLOWED"
    delivery_methods = ["SMS", "Voice"]
    countries        = ["US", "CA"]
  }

  cooldown_configuration = {
    email = {
      enabled      = true
      group_by     = "RECIPIENT"
      resend_limit = 5
      periods = [
        {
          duration  = 10
          time_unit = "MINUTES"
        },
        {
          duration  = 30
          time_unit = "SECONDS"
        },
        {
          duration  = 1
          time_unit = "MINUTES"
        }
      ]
    }

    sms = {
      enabled      = true
      group_by     = "RECIPIENT"
      resend_limit = 3
      periods = [
        {
          duration  = 5
          time_unit = "MINUTES"
        },
        {
          duration  = 30
          time_unit = "SECONDS"
        },
        {
          duration  = 2
          time_unit = "MINUTES"
        }
      ]
    }

    voice = {
      enabled = false
    }

    whats_app = {
      enabled = false
    }
  }

  provider_configuration = {
    conditions = [
      {
        delivery_methods = ["SMS", "VOICE"]
        countries        = ["US", "CA"]
        fallback_chain = [{
          id = pingone_phone_delivery_settings.%[2]s_provider_1.id
        }]
      },
      {
        delivery_methods = ["SMS", "VOICE"]
        fallback_chain = [{
          id = pingone_phone_delivery_settings.%[2]s_provider_2.id
        }]
      }
    ]
  }
}

data "pingone_notification_policy" "%[2]s" {
  environment_id         = data.pingone_environment.general_test.id
  notification_policy_id = pingone_notification_policy.%[2]s.id

  depends_on = [pingone_notification_policy.%[2]s]
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccNotificationPolicyDataSource_FindByName(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_notification_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  quota = [{
    type             = "ENVIRONMENT"
    delivery_methods = ["SMS", "Voice"]
    total            = 100
  }]
}

data "pingone_notification_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = pingone_notification_policy.%[2]s.name

  depends_on = [pingone_notification_policy.%[2]s]
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccNotificationPolicyDataSource_FindByIDFail(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_notification_policy" "%[2]s" {
  environment_id         = data.pingone_environment.general_test.id
  notification_policy_id = "9c052a8a-14be-44e4-8f07-2662569994ce" // dummy ID that conforms to UUID v4
}`, acctest.GenericSandboxEnvironment(), resourceName)
}

func testAccNotificationPolicyDataSource_FindByNameFail(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_notification_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "doesnotexist"
}`, acctest.GenericSandboxEnvironment(), resourceName)
}
