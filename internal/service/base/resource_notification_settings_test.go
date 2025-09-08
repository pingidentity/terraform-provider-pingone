// Copyright © 2025 Ping Identity Corporation

package base_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/base"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccNotificationSettings_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_notification_settings.%s", resourceName)

	name := resourceName

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	var notificationSettingsID string

	var p1Client *client.Client
	var ctx = context.Background()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoBeta(t)
			p1Client = acctest.PreCheckTestClient(ctx, t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.NotificationSettings_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccNotificationSettingsConfig_Minimal(environmentName, licenseID, resourceName, name),
				Check:  base.NotificationSettings_GetIDs(resourceFullName, &notificationSettingsID),
			},
			{
				PreConfig: func() {
					base.NotificationSettings_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, notificationSettingsID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccNotificationSettings_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_notification_settings.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	name := resourceName

	fullStep := resource.TestStep{
		Config: testAccNotificationSettingsConfig_Full(environmentName, licenseID, resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "delivery_mode", "ALL"),
			resource.TestCheckResourceAttr(resourceFullName, "provider_fallback_chain.#", "2"),
			resource.TestMatchResourceAttr(resourceFullName, "provider_fallback_chain.0", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "provider_fallback_chain.1", "PINGONE_TWILIO"),
			resource.TestCheckResourceAttr(resourceFullName, "allowed_list.#", "3"),
			resource.TestMatchResourceAttr(resourceFullName, "allowed_list.0.user_id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "allowed_list.1.user_id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "allowed_list.2.user_id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "updated_at", verify.RFC3339Regexp),
			resource.TestCheckResourceAttr(resourceFullName, "from.email_address", "noreply@pingidentity.com"),
			resource.TestCheckNoResourceAttr(resourceFullName, "from.name"),
			resource.TestCheckResourceAttr(resourceFullName, "reply_to.email_address", "noreply@pingidentity.com"),
			resource.TestCheckNoResourceAttr(resourceFullName, "reply_to.name"),
		),
	}

	minimalStep := resource.TestStep{
		Config: testAccNotificationSettingsConfig_Minimal(environmentName, licenseID, resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "delivery_mode", "ALL"),
			resource.TestCheckResourceAttr(resourceFullName, "provider_fallback_chain.#", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "allowed_list.#", "0"),
			resource.TestMatchResourceAttr(resourceFullName, "updated_at", verify.RFC3339Regexp),
			resource.TestCheckResourceAttr(resourceFullName, "from.email_address", "noreply@pingidentity.com"),
			resource.TestCheckResourceAttr(resourceFullName, "from.name", "PingOne"),
			resource.TestCheckNoResourceAttr(resourceFullName, "reply_to.email_address"),
			resource.TestCheckNoResourceAttr(resourceFullName, "reply_to.name"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.NotificationSettings_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full from scratch
			fullStep,
			{
				Config:  testAccNotificationSettingsConfig_Full(environmentName, licenseID, resourceName, name),
				Destroy: true,
			},
			// Minimal from scratch
			minimalStep,
			{
				Config:  testAccNotificationSettingsConfig_Minimal(environmentName, licenseID, resourceName, name),
				Destroy: true,
			},
			// Update
			fullStep,
			minimalStep,
			fullStep,
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

func TestAccNotificationSettings_EmailSources(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_notification_settings.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	fromFull := resource.TestStep{
		Config: testAccNotificationSettingsConfig_FromFull(environmentName, licenseID, resourceName),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "from.email_address", "noreply@pingidentity.com"),
			resource.TestCheckResourceAttr(resourceFullName, "from.name", "Stubbed From Address"),
			resource.TestCheckNoResourceAttr(resourceFullName, "reply_to.email_address"),
			resource.TestCheckNoResourceAttr(resourceFullName, "reply_to.name"),
		),
	}

	fromMinimal := resource.TestStep{
		Config: testAccNotificationSettingsConfig_FromMinimal(environmentName, licenseID, resourceName),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "from.email_address", "noreply@pingidentity.com"),
			resource.TestCheckNoResourceAttr(resourceFullName, "from.name"),
			resource.TestCheckNoResourceAttr(resourceFullName, "reply_to.email_address"),
			resource.TestCheckNoResourceAttr(resourceFullName, "reply_to.name"),
		),
	}

	replyToFull := resource.TestStep{
		Config: testAccNotificationSettingsConfig_ReplyToFull(environmentName, licenseID, resourceName),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "from.email_address", "noreply@pingidentity.com"),
			resource.TestCheckResourceAttr(resourceFullName, "from.name", "PingOne"),
			resource.TestCheckResourceAttr(resourceFullName, "reply_to.email_address", "noreply@pingidentity.com"),
			resource.TestCheckResourceAttr(resourceFullName, "reply_to.name", "Stubbed Reply To Address"),
		),
	}

	replyToMinimal := resource.TestStep{
		Config: testAccNotificationSettingsConfig_ReplyToMinimal(environmentName, licenseID, resourceName),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "from.email_address", "noreply@pingidentity.com"),
			resource.TestCheckResourceAttr(resourceFullName, "from.name", "PingOne"),
			resource.TestCheckResourceAttr(resourceFullName, "reply_to.email_address", "noreply@pingidentity.com"),
			resource.TestCheckNoResourceAttr(resourceFullName, "reply_to.name"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.NotificationSettings_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			replyToMinimal,
			fromFull,
			// Variant 1 New
			fromFull,
			{
				Config:  testAccNotificationSettingsConfig_FromFull(environmentName, licenseID, resourceName),
				Destroy: true,
			},
			// Variant 2 New
			fromMinimal,
			{
				Config:  testAccNotificationSettingsConfig_FromMinimal(environmentName, licenseID, resourceName),
				Destroy: true,
			},
			// Variant 3 New
			replyToFull,
			{
				Config:  testAccNotificationSettingsConfig_ReplyToFull(environmentName, licenseID, resourceName),
				Destroy: true,
			},
			// Variant 3 New
			replyToMinimal,
			{
				Config:  testAccNotificationSettingsConfig_ReplyToMinimal(environmentName, licenseID, resourceName),
				Destroy: true,
			},
			// Update
			fromFull,
			fromMinimal,
			replyToFull,
			replyToMinimal,
			fromFull,
		},
	})
}

func TestAccNotificationSettings_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_notification_settings.%s", resourceName)

	name := resourceName

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.NotificationSettings_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccNotificationSettingsConfig_Minimal(environmentName, licenseID, resourceName, name),
			},
			// Errors
			{
				ResourceName:  resourceFullName,
				ImportStateId: "badformat",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`Unexpected Import Identifier`),
			},
		},
	})
}

func testAccNotificationSettingsConfig_Full(environmentName, licenseID, resourceName, name string) string {
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
}

resource "pingone_population" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  name = "%[4]s"
}

resource "pingone_user" "%[3]s-1" {
  environment_id = pingone_environment.%[2]s.id

  username      = "%[4]s-1"
  email         = "noreply@pingidentity.com"
  population_id = pingone_population.%[3]s.id
}

resource "pingone_user" "%[3]s-2" {
  environment_id = pingone_environment.%[2]s.id

  username      = "%[4]s-2"
  email         = "noreply@pingidentity.com"
  population_id = pingone_population.%[3]s.id
}

resource "pingone_user" "%[3]s-3" {
  environment_id = pingone_environment.%[2]s.id

  username      = "%[4]s-3"
  email         = "noreply@pingidentity.com"
  population_id = pingone_population.%[3]s.id
}

resource "pingone_notification_settings" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  delivery_mode = "ALL"

  provider_fallback_chain = [
    pingone_phone_delivery_settings.%[3]s.id,
    "PINGONE_TWILIO",
  ]

  allowed_list = [
    {
      user_id = pingone_user.%[3]s-3.id
    },
    {
      user_id = pingone_user.%[3]s-1.id
    },
    {
      user_id = pingone_user.%[3]s-2.id
    },
  ]

  from = {
    email_address = "noreply@pingidentity.com"
  }

  reply_to = {
    email_address = "noreply@pingidentity.com"
  }
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccNotificationSettingsConfig_Minimal(environmentName, licenseID, resourceName, name string) string {
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
}

resource "pingone_population" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  name = "%[4]s"
}

resource "pingone_user" "%[3]s-1" {
  environment_id = pingone_environment.%[2]s.id

  username      = "%[4]s-1"
  email         = "noreply@pingidentity.com"
  population_id = pingone_population.%[3]s.id
}

resource "pingone_user" "%[3]s-2" {
  environment_id = pingone_environment.%[2]s.id

  username      = "%[4]s-2"
  email         = "noreply@pingidentity.com"
  population_id = pingone_population.%[3]s.id
}

resource "pingone_user" "%[3]s-3" {
  environment_id = pingone_environment.%[2]s.id

  username      = "%[4]s-3"
  email         = "noreply@pingidentity.com"
  population_id = pingone_population.%[3]s.id
}

resource "pingone_notification_settings" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccNotificationSettingsConfig_FromFull(environmentName, licenseID, resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_notification_settings" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  from = {
    email_address = "noreply@pingidentity.com"
    name          = "Stubbed From Address"
  }
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName)
}

func testAccNotificationSettingsConfig_FromMinimal(environmentName, licenseID, resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_notification_settings" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  from = {
    email_address = "noreply@pingidentity.com"
  }
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName)
}

func testAccNotificationSettingsConfig_ReplyToFull(environmentName, licenseID, resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_notification_settings" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  reply_to = {
    email_address = "noreply@pingidentity.com"
    name          = "Stubbed Reply To Address"
  }
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName)
}

func testAccNotificationSettingsConfig_ReplyToMinimal(environmentName, licenseID, resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_notification_settings" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  reply_to = {
    email_address = "noreply@pingidentity.com"
  }
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName)
}
