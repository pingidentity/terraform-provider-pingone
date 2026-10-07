// Copyright © 2026 Ping Identity Corporation

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
	acctestlegacysdk "github.com/pingidentity/terraform-provider-pingone/internal/acctest/legacysdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/base"
	baselegacysdk "github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/base/legacysdk"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccNotificationSettingsEmail_SMTP_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_notification_settings_email.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

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
		CheckDestroy:             base.NotificationSettingsEmail_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Test removal of the environment
			{
				Config: testAccNotificationSettingsEmail_SMTPConfig_Full(environmentName, licenseID, resourceName),
				Check:  base.NotificationSettingsEmail_GetIDs(resourceFullName, &environmentID),
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

func TestAccNotificationSettingsEmail_SMTP_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_notification_settings_email.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	fullStep := resource.TestStep{
		Config: testAccNotificationSettingsEmail_SMTPConfig_Full(environmentName, licenseID, resourceName),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "host", "pingidentity.com"),
			resource.TestCheckResourceAttr(resourceFullName, "port", "25"),
			resource.TestCheckResourceAttr(resourceFullName, "protocol", "SMTPS"),
			resource.TestCheckResourceAttr(resourceFullName, "username", "smtpuser"),
			resource.TestCheckResourceAttr(resourceFullName, "password", "smtpuserpassword"),
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
		CheckDestroy:             base.NotificationSettingsEmail_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full from scratch
			fullStep,
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("resource not found: %s", resourceFullName)
						}

						return rs.Primary.ID, nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password",
					"auth_token",
				},
			},
		},
	})
}

func TestAccNotificationSettingsEmail_SMTP_EmailSources(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_notification_settings_email.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	fromFull := resource.TestStep{
		Config: testAccNotificationSettingsEmail_SMTPConfig_FromFull(environmentName, licenseID, resourceName),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "from.email_address", "noreply@pingidentity.com"),
			resource.TestCheckResourceAttr(resourceFullName, "from.name", "Stubbed From Address"),
			resource.TestCheckNoResourceAttr(resourceFullName, "reply_to.email_address"),
			resource.TestCheckNoResourceAttr(resourceFullName, "reply_to.name"),
		),
	}

	fromMinimal := resource.TestStep{
		Config: testAccNotificationSettingsEmail_SMTPConfig_FromMinimal(environmentName, licenseID, resourceName),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "from.email_address", "noreply@pingidentity.com"),
			resource.TestCheckNoResourceAttr(resourceFullName, "from.name"),
			resource.TestCheckNoResourceAttr(resourceFullName, "reply_to.email_address"),
			resource.TestCheckNoResourceAttr(resourceFullName, "reply_to.name"),
		),
	}

	replyToFull := resource.TestStep{
		Config: testAccNotificationSettingsEmail_SMTPConfig_ReplyToFull(environmentName, licenseID, resourceName),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "from.email_address", "noreply@pingidentity.com"),
			resource.TestCheckResourceAttr(resourceFullName, "from.name", "Stubbed From Address"),
			resource.TestCheckResourceAttr(resourceFullName, "reply_to.email_address", "reply@pingidentity.com"),
			resource.TestCheckResourceAttr(resourceFullName, "reply_to.name", "Stubbed Reply To Address"),
		),
	}

	replyToMinimal := resource.TestStep{
		Config: testAccNotificationSettingsEmail_SMTPConfig_ReplyToMinimal(environmentName, licenseID, resourceName),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "from.email_address", "noreply@pingidentity.com"),
			resource.TestCheckNoResourceAttr(resourceFullName, "from.name"),
			resource.TestCheckResourceAttr(resourceFullName, "reply_to.email_address", "reply@pingidentity.com"),
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
		CheckDestroy:             base.NotificationSettingsEmail_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			replyToMinimal,
			fromFull,
			// Variant 1 New
			fromFull,
			{
				Config:  testAccNotificationSettingsEmail_SMTPConfig_FromFull(environmentName, licenseID, resourceName),
				Destroy: true,
			},
			// Variant 2 New
			fromMinimal,
			{
				Config:  testAccNotificationSettingsEmail_SMTPConfig_FromMinimal(environmentName, licenseID, resourceName),
				Destroy: true,
			},
			// Variant 3 New
			replyToFull,
			{
				Config:  testAccNotificationSettingsEmail_SMTPConfig_ReplyToFull(environmentName, licenseID, resourceName),
				Destroy: true,
			},
			// Variant 3 New
			replyToMinimal,
			{
				Config:  testAccNotificationSettingsEmail_SMTPConfig_ReplyToMinimal(environmentName, licenseID, resourceName),
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

func TestAccNotificationSettingsEmail_SMTP_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_notification_settings_email.%s", resourceName)

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
		CheckDestroy:             base.NotificationSettingsEmail_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccNotificationSettingsEmail_SMTPConfig_Full(environmentName, licenseID, resourceName),
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
				ImportStateId: "badformat/badformat",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`Unexpected Import Identifier`),
			},
		},
	})
}

func testAccNotificationSettingsEmail_SMTPConfig_Full(environmentName, licenseID, resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_notification_settings_email" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  host     = "pingidentity.com"
  port     = 25
  username = "smtpuser"
  password = "smtpuserpassword"

  from = {
    email_address = "noreply@pingidentity.com"
    name          = "Stubbed From Address"
  }

  reply_to = {
    email_address = "reply@pingidentity.com"
    name          = "Stubbed Reply To Address"
  }
}`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName)
}

func testAccNotificationSettingsEmail_SMTPConfig_FromFull(environmentName, licenseID, resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_notification_settings_email" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  host     = "pingidentity.com"
  port     = 25
  username = "smtpuser"
  password = "smtpuserpassword"

  from = {
    email_address = "noreply@pingidentity.com"
    name          = "Stubbed From Address"
  }
}`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName)
}

func testAccNotificationSettingsEmail_SMTPConfig_FromMinimal(environmentName, licenseID, resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_notification_settings_email" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  host     = "pingidentity.com"
  port     = 25
  username = "smtpuser"
  password = "smtpuserpassword"

  from = {
    email_address = "noreply@pingidentity.com"
  }
}`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName)
}

func testAccNotificationSettingsEmail_SMTPConfig_ReplyToFull(environmentName, licenseID, resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_notification_settings_email" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  host     = "pingidentity.com"
  port     = 25
  username = "smtpuser"
  password = "smtpuserpassword"

  from = {
    email_address = "noreply@pingidentity.com"
    name          = "Stubbed From Address"
  }

  reply_to = {
    email_address = "reply@pingidentity.com"
    name          = "Stubbed Reply To Address"
  }
}`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName)
}

func testAccNotificationSettingsEmail_SMTPConfig_ReplyToMinimal(environmentName, licenseID, resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_notification_settings_email" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  host     = "pingidentity.com"
  port     = 25
  username = "smtpuser"
  password = "smtpuserpassword"

  from = {
    email_address = "noreply@pingidentity.com"
  }

  reply_to = {
    email_address = "reply@pingidentity.com"
  }
}`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName)
}

func TestAccNotificationSettingsEmail_CustomProvider_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_notification_settings_email.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

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
		CheckDestroy:             base.NotificationSettingsEmail_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Test removal of the environment
			{
				Config: testAccNotificationSettingsEmail_CustomProviderConfig_GET_Full(environmentName, licenseID, resourceName),
				Check:  base.NotificationSettingsEmail_GetIDs(resourceFullName, &environmentID),
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

func TestAccNotificationSettingsEmail_CustomProvider_GET(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_notification_settings_email.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()
	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	// Test switching to SMTP
	smtpFullSwitchStep := resource.TestStep{
		Config: testAccNotificationSettingsEmail_SMTPConfig_Full(environmentName, licenseID, resourceName),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "host", "pingidentity.com"),
			resource.TestCheckResourceAttr(resourceFullName, "port", "25"),
			resource.TestCheckResourceAttr(resourceFullName, "protocol", "SMTPS"),
			resource.TestCheckResourceAttr(resourceFullName, "username", "smtpuser"),
			resource.TestCheckResourceAttr(resourceFullName, "password", "smtpuserpassword"),
			resource.TestCheckResourceAttr(resourceFullName, "from.email_address", "noreply@pingidentity.com"),
			resource.TestCheckResourceAttr(resourceFullName, "from.name", "Stubbed From Address"),
			resource.TestCheckResourceAttr(resourceFullName, "reply_to.email_address", "reply@pingidentity.com"),
			resource.TestCheckResourceAttr(resourceFullName, "reply_to.name", "Stubbed Reply To Address"),
		),
	}

	minimalStep := resource.TestStep{
		Config: testAccNotificationSettingsEmail_CustomProviderConfig_GET_Minimal(environmentName, licenseID, resourceName),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "custom_provider_name", "CustomProviderName"),
			resource.TestCheckResourceAttr(resourceFullName, "protocol", "HTTP"),
			resource.TestCheckResourceAttr(resourceFullName, "requests.#", "1"),
			resource.TestCheckResourceAttr(resourceFullName, "requests.0.method", "GET"),
			resource.TestCheckResourceAttr(resourceFullName, "requests.0.headers.Content-Type", "application/x-www-form-urlencoded"),
			resource.TestCheckResourceAttr(resourceFullName, "requests.0.headers.subject", "${subject}"),
			resource.TestCheckResourceAttr(resourceFullName, "requests.0.headers.reply-to", "${reply_to}"),
			resource.TestCheckResourceAttr(resourceFullName, "requests.0.headers.from", "${from}"),
			resource.TestCheckNoResourceAttr(resourceFullName, "requests.0.body"),
			resource.TestCheckResourceAttr(resourceFullName, "requests.0.url", "https://api.pingidentity.com/send-email.apx?${to}-${message}"),
		),
	}

	minimalStepDestroy := resource.TestStep{
		Config:  minimalStep.Config,
		Destroy: true,
	}

	fullStep := resource.TestStep{
		Config: testAccNotificationSettingsEmail_CustomProviderConfig_GET_Full(environmentName, licenseID, resourceName),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "custom_provider_name", "UpdatedCustomProviderName"),
			resource.TestCheckResourceAttr(resourceFullName, "from.name", "Updated Test Sender"),
			resource.TestCheckResourceAttr(resourceFullName, "reply_to.name", "Updated Test Reply To"),
			resource.TestCheckResourceAttr(resourceFullName, "protocol", "HTTP"),
			resource.TestCheckResourceAttr(resourceFullName, "requests.#", "1"),
			resource.TestCheckResourceAttr(resourceFullName, "requests.0.method", "GET"),
			resource.TestCheckResourceAttr(resourceFullName, "requests.0.headers.Content-Type", "application/x-www-form-urlencoded"),
			resource.TestCheckResourceAttr(resourceFullName, "requests.0.headers.updated-subject", "updated-${subject}"),
			resource.TestCheckResourceAttr(resourceFullName, "requests.0.headers.updated-reply-to", "updated-${reply_to}"),
			resource.TestCheckResourceAttr(resourceFullName, "requests.0.headers.updated-from", "updated-${from}"),
			resource.TestCheckNoResourceAttr(resourceFullName, "requests.0.body"),
			resource.TestCheckResourceAttr(resourceFullName, "requests.0.url", "https://api.pingidentity.com/updated-send-email.aspx?${to}-${message}"),
		),
	}

	fullStepDestroy := resource.TestStep{
		Config:  fullStep.Config,
		Destroy: true,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.NotificationSettingsEmail_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			minimalStep,
			smtpFullSwitchStep,
			minimalStep,
			minimalStepDestroy,
			fullStep,
			smtpFullSwitchStep,
			minimalStep,
			fullStep,
			fullStepDestroy,
			fullStep,
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("resource not found: %s", resourceFullName)
						}

						return rs.Primary.ID, nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password",
					"auth_token",
				},
			},
		},
	})
}

func testAccNotificationSettingsEmail_CustomProviderConfig_GET_Minimal(environmentName, licenseID, resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_notification_settings_email" "%[3]s" {
  environment_id       = pingone_environment.%[2]s.id
  custom_provider_name = "CustomProviderName"

  username = "customuser"
  password = "customuserpassword"

  protocol = "HTTP"

  from = {
    email_address = "no-reply@pingidentity.com"
  }

  reply_to = {
    email_address = "updated-reply@pingidentity.com"
  }

  requests = [
    {
      method = "GET"
      headers = {
        "Content-Type" = "application/x-www-form-urlencoded"
        "subject"      = "$${subject}"
        "reply-to"     = "$${reply_to}"
        "from"         = "$${from}"
      }
      url = "https://api.pingidentity.com/send-email.apx?$${to}-$${message}"
    }
  ]
}`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName)
}

func testAccNotificationSettingsEmail_CustomProviderConfig_GET_Full(environmentName, licenseID, resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_notification_settings_email" "%[3]s" {
  environment_id       = pingone_environment.%[2]s.id
  custom_provider_name = "UpdatedCustomProviderName"

  auth_token = "customauthtoken"

  protocol = "HTTP"

  from = {
    name          = "Updated Test Sender"
    email_address = "updated-no-reply@pingidentity.com"
  }

  reply_to = {
    name          = "Updated Test Reply To"
    email_address = "updated-reply@pingidentity.com"
  }

  requests = [
    {
      method = "GET"
      headers = {
        "Content-Type"     = "application/x-www-form-urlencoded"
        "updated-subject"  = "updated-$${subject}"
        "updated-reply-to" = "updated-$${reply_to}"
        "updated-from"     = "updated-$${from}"
      }
      url = "https://api.pingidentity.com/updated-send-email.aspx?$${to}-$${message}"
    }
  ]
}`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName)
}

func TestAccNotificationSettingsEmail_CustomProvider_POST_NoKeyValuesBody(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_notification_settings_email.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()
	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	// Test switching to SMTP
	smtpFullSwitchStep := resource.TestStep{
		Config: testAccNotificationSettingsEmail_SMTPConfig_Full(environmentName, licenseID, resourceName),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "host", "pingidentity.com"),
			resource.TestCheckResourceAttr(resourceFullName, "port", "25"),
			resource.TestCheckResourceAttr(resourceFullName, "protocol", "SMTPS"),
			resource.TestCheckResourceAttr(resourceFullName, "username", "smtpuser"),
			resource.TestCheckResourceAttr(resourceFullName, "password", "smtpuserpassword"),
			resource.TestCheckResourceAttr(resourceFullName, "from.email_address", "noreply@pingidentity.com"),
			resource.TestCheckResourceAttr(resourceFullName, "from.name", "Stubbed From Address"),
			resource.TestCheckResourceAttr(resourceFullName, "reply_to.email_address", "reply@pingidentity.com"),
			resource.TestCheckResourceAttr(resourceFullName, "reply_to.name", "Stubbed Reply To Address"),
		),
	}

	minimalStep := resource.TestStep{
		Config: testAccNotificationSettingsEmail_CustomProviderConfig_POST_NoKeyValuesBody_Minimal(environmentName, licenseID, resourceName),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "custom_provider_name", "CustomProviderName"),
			resource.TestCheckResourceAttr(resourceFullName, "protocol", "HTTP"),
			resource.TestCheckResourceAttr(resourceFullName, "requests.#", "1"),
			resource.TestCheckResourceAttr(resourceFullName, "requests.0.method", "POST"),
			resource.TestCheckResourceAttr(resourceFullName, "requests.0.url", "https://api.pingidentity.com/send-email"),
			resource.TestCheckNoResourceAttr(resourceFullName, "from.name"),
			resource.TestCheckNoResourceAttr(resourceFullName, "reply_to.name"),
			resource.TestCheckResourceAttr(resourceFullName, "requests.0.headers.Content-Type", "application/x-www-form-urlencoded"),
			resource.TestCheckResourceAttr(resourceFullName, "requests.0.headers.subject", "${subject}"),
			resource.TestCheckResourceAttr(resourceFullName, "requests.0.headers.reply-to", "${reply_to}"),
			resource.TestCheckResourceAttr(resourceFullName, "requests.0.headers.from", "${from}"),
			resource.TestCheckNoResourceAttr(resourceFullName, "requests.0.body"),
		),
	}

	minimalStepDestroy := resource.TestStep{
		Config:  minimalStep.Config,
		Destroy: true,
	}

	fullStep := resource.TestStep{
		Config: testAccNotificationSettingsEmail_CustomProviderConfig_POST_NoKeyValuesBody_Full(environmentName, licenseID, resourceName),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "custom_provider_name", "UpdatedCustomProviderName"),
			resource.TestCheckResourceAttr(resourceFullName, "protocol", "HTTP"),
			resource.TestCheckResourceAttr(resourceFullName, "requests.#", "1"),
			resource.TestCheckResourceAttr(resourceFullName, "requests.0.method", "POST"),
			resource.TestCheckResourceAttr(resourceFullName, "requests.0.url", "https://api.pingidentity.com/updated-send-email"),
			resource.TestCheckResourceAttr(resourceFullName, "from.name", "Updated Test Sender"),
			resource.TestCheckResourceAttr(resourceFullName, "reply_to.name", "Updated Test Reply To"),
			resource.TestCheckResourceAttr(resourceFullName, "requests.0.headers.Content-Type", "application/x-www-form-urlencoded"),
			resource.TestCheckResourceAttr(resourceFullName, "requests.0.headers.updated-subject", "updated-${subject}"),
			resource.TestCheckResourceAttr(resourceFullName, "requests.0.headers.updated-reply-to", "updated-${reply_to}"),
			resource.TestCheckResourceAttr(resourceFullName, "requests.0.headers.updated-from", "updated-${from}"),
		),
	}

	fullStepDestroy := resource.TestStep{
		Config:  fullStep.Config,
		Destroy: true,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.NotificationSettingsEmail_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			minimalStep,
			smtpFullSwitchStep,
			minimalStep,
			minimalStepDestroy,
			fullStep,
			smtpFullSwitchStep,
			minimalStep,
			fullStep,
			fullStepDestroy,
			fullStep,
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("resource not found: %s", resourceFullName)
						}

						return rs.Primary.ID, nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password",
					"auth_token",
				},
			},
		},
	})
}

func testAccNotificationSettingsEmail_CustomProviderConfig_POST_NoKeyValuesBody_Minimal(environmentName, licenseID, resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_notification_settings_email" "%[3]s" {
  environment_id       = pingone_environment.%[2]s.id
  custom_provider_name = "CustomProviderName"

  username = "customuser"
  password = "customuserpassword"

  protocol = "HTTP"

  from = {
    email_address = "no-reply@pingidentity.com"
  }

  reply_to = {
    email_address = "reply@pingidentity.com"
  }

  requests = [
    {
      method = "POST"
      headers = {
        "Content-Type" = "application/x-www-form-urlencoded"
        "subject"      = "$${subject}"
        "reply-to"     = "$${reply_to}"
        "from"         = "$${from}"
      }
      url = "https://api.pingidentity.com/send-email"
    }
  ]
}`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName)
}

func testAccNotificationSettingsEmail_CustomProviderConfig_POST_NoKeyValuesBody_Full(environmentName, licenseID, resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_notification_settings_email" "%[3]s" {
  environment_id       = pingone_environment.%[2]s.id
  custom_provider_name = "UpdatedCustomProviderName"

  auth_token = "customauthtoken"

  protocol = "HTTP"

  from = {
    name          = "Updated Test Sender"
    email_address = "updated-no-reply@pingidentity.com"
  }

  reply_to = {
    name          = "Updated Test Reply To"
    email_address = "updated-reply@pingidentity.com"
  }

  requests = [
    {
      method = "POST"
      headers = {
        "Content-Type"     = "application/x-www-form-urlencoded"
        "updated-subject"  = "updated-$${subject}"
        "updated-reply-to" = "updated-$${reply_to}"
        "updated-from"     = "updated-$${from}"
      }
      url = "https://api.pingidentity.com/updated-send-email"
    }
  ]
}`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName)
}

func TestAccNotificationSettingsEmail_CustomProvider_POST_KeyValuesBody(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_notification_settings_email.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()
	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	// Test switching to SMTP
	smtpFullSwitchStep := resource.TestStep{
		Config: testAccNotificationSettingsEmail_SMTPConfig_Full(environmentName, licenseID, resourceName),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "host", "pingidentity.com"),
			resource.TestCheckResourceAttr(resourceFullName, "port", "25"),
			resource.TestCheckResourceAttr(resourceFullName, "protocol", "SMTPS"),
			resource.TestCheckResourceAttr(resourceFullName, "username", "smtpuser"),
			resource.TestCheckResourceAttr(resourceFullName, "password", "smtpuserpassword"),
			resource.TestCheckResourceAttr(resourceFullName, "from.email_address", "noreply@pingidentity.com"),
			resource.TestCheckResourceAttr(resourceFullName, "from.name", "Stubbed From Address"),
			resource.TestCheckResourceAttr(resourceFullName, "reply_to.email_address", "reply@pingidentity.com"),
			resource.TestCheckResourceAttr(resourceFullName, "reply_to.name", "Stubbed Reply To Address"),
		),
	}

	minimalStep := resource.TestStep{
		Config: testAccNotificationSettingsEmail_CustomProviderConfig_POST_KeyValuesBody_Minimal(environmentName, licenseID, resourceName),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "custom_provider_name", "CustomProviderName"),
			resource.TestCheckResourceAttr(resourceFullName, "protocol", "HTTP"),
			resource.TestCheckResourceAttr(resourceFullName, "requests.#", "1"),
			resource.TestCheckResourceAttr(resourceFullName, "requests.0.method", "POST"),
			resource.TestCheckResourceAttr(resourceFullName, "requests.0.url", "https://api.pingidentity.com/send-email"),
			resource.TestCheckNoResourceAttr(resourceFullName, "from.name"),
			resource.TestCheckNoResourceAttr(resourceFullName, "reply_to.name"),
			resource.TestCheckResourceAttr(resourceFullName, "requests.0.headers.Content-Type", "application/x-www-form-urlencoded"),
			resource.TestCheckResourceAttr(resourceFullName, "requests.0.headers.subject", "${subject}"),
			resource.TestCheckResourceAttr(resourceFullName, "requests.0.headers.reply-to", "${reply_to}"),
			resource.TestCheckResourceAttr(resourceFullName, "requests.0.headers.from", "${from}"),
			resource.TestCheckResourceAttr(resourceFullName, "requests.0.body", "to=${to}&message=${message}"),
		),
	}

	minimalStepDestroy := resource.TestStep{
		Config:  minimalStep.Config,
		Destroy: true,
	}

	fullStep := resource.TestStep{
		Config: testAccNotificationSettingsEmail_CustomProviderConfig_POST_KeyValuesBody_Full(environmentName, licenseID, resourceName),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "custom_provider_name", "UpdatedCustomProviderName"),
			resource.TestCheckResourceAttr(resourceFullName, "protocol", "HTTP"),
			resource.TestCheckResourceAttr(resourceFullName, "requests.#", "1"),
			resource.TestCheckResourceAttr(resourceFullName, "requests.0.method", "POST"),
			resource.TestCheckResourceAttr(resourceFullName, "requests.0.url", "https://api.pingidentity.com/updated-send-email"),
			resource.TestCheckResourceAttr(resourceFullName, "from.name", "Updated Test Sender"),
			resource.TestCheckResourceAttr(resourceFullName, "reply_to.name", "Updated Test Reply To"),
			resource.TestCheckResourceAttr(resourceFullName, "requests.0.headers.Content-Type", "application/x-www-form-urlencoded"),
			resource.TestCheckResourceAttr(resourceFullName, "requests.0.headers.updated-subject", "updated-${subject}"),
			resource.TestCheckResourceAttr(resourceFullName, "requests.0.headers.updated-reply-to", "updated-${reply_to}"),
			resource.TestCheckResourceAttr(resourceFullName, "requests.0.headers.updated-from", "updated-${from}"),
			resource.TestCheckResourceAttr(resourceFullName, "requests.0.body", "updated-to=${to}&message=${message}"),
		),
	}

	fullStepDestroy := resource.TestStep{
		Config:  fullStep.Config,
		Destroy: true,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.NotificationSettingsEmail_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			minimalStep,
			smtpFullSwitchStep,
			minimalStep,
			minimalStepDestroy,
			fullStep,
			smtpFullSwitchStep,
			minimalStep,
			fullStep,
			fullStepDestroy,
			fullStep,
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("resource not found: %s", resourceFullName)
						}

						return rs.Primary.ID, nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password",
					"auth_token",
				},
			},
		},
	})
}

func testAccNotificationSettingsEmail_CustomProviderConfig_POST_KeyValuesBody_Minimal(environmentName, licenseID, resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_notification_settings_email" "%[3]s" {
  environment_id       = pingone_environment.%[2]s.id
  custom_provider_name = "CustomProviderName"

  username = "customuser"
  password = "customuserpassword"

  protocol = "HTTP"

  from = {
    email_address = "no-reply@pingidentity.com"
  }

  reply_to = {
    email_address = "reply@pingidentity.com"
  }

  requests = [
    {
      method = "POST"
      headers = {
        "Content-Type" = "application/x-www-form-urlencoded"
        "subject"      = "$${subject}"
        "reply-to"     = "$${reply_to}"
        "from"         = "$${from}"
      }
      body = "to=$${to}&message=$${message}"
      url  = "https://api.pingidentity.com/send-email"
    }
  ]
}`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName)
}

func testAccNotificationSettingsEmail_CustomProviderConfig_POST_KeyValuesBody_Full(environmentName, licenseID, resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_notification_settings_email" "%[3]s" {
  environment_id       = pingone_environment.%[2]s.id
  custom_provider_name = "UpdatedCustomProviderName"

  auth_token = "customauthtoken"

  protocol = "HTTP"

  from = {
    name          = "Updated Test Sender"
    email_address = "updated-no-reply@pingidentity.com"
  }

  reply_to = {
    name          = "Updated Test Reply To"
    email_address = "updated-reply@pingidentity.com"
  }

  requests = [
    {
      method = "POST"
      headers = {
        "Content-Type"     = "application/x-www-form-urlencoded"
        "updated-subject"  = "updated-$${subject}"
        "updated-reply-to" = "updated-$${reply_to}"
        "updated-from"     = "updated-$${from}"
      }
      body = "updated-to=$${to}&message=$${message}"
      url  = "https://api.pingidentity.com/updated-send-email"
    }
  ]
}`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName)
}

func TestAccNotificationSettingsEmail_CustomProvider_POST_RawBody(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_notification_settings_email.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()
	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	// Test switching to SMTP
	smtpFullSwitchStep := resource.TestStep{
		Config: testAccNotificationSettingsEmail_SMTPConfig_Full(environmentName, licenseID, resourceName),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "host", "pingidentity.com"),
			resource.TestCheckResourceAttr(resourceFullName, "port", "25"),
			resource.TestCheckResourceAttr(resourceFullName, "protocol", "SMTPS"),
			resource.TestCheckResourceAttr(resourceFullName, "username", "smtpuser"),
			resource.TestCheckResourceAttr(resourceFullName, "password", "smtpuserpassword"),
			resource.TestCheckResourceAttr(resourceFullName, "from.email_address", "noreply@pingidentity.com"),
			resource.TestCheckResourceAttr(resourceFullName, "from.name", "Stubbed From Address"),
			resource.TestCheckResourceAttr(resourceFullName, "reply_to.email_address", "reply@pingidentity.com"),
			resource.TestCheckResourceAttr(resourceFullName, "reply_to.name", "Stubbed Reply To Address"),
		),
	}

	// Test switching to Custom Provider
	customProviderSwitchStep := resource.TestStep{
		Config: testAccNotificationSettingsEmail_CustomProviderConfig_POST_RawBody_Full(environmentName, licenseID, resourceName),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "custom_provider_name", "UpdatedCustomProviderName"),
			resource.TestCheckResourceAttr(resourceFullName, "auth_token", "customauthtoken"),
			resource.TestCheckResourceAttr(resourceFullName, "from.name", "Updated Test Sender"),
			resource.TestCheckResourceAttr(resourceFullName, "from.email_address", "updated-no-reply@pingidentity.com"),
			resource.TestCheckResourceAttr(resourceFullName, "reply_to.name", "Updated Test Reply To"),
			resource.TestCheckResourceAttr(resourceFullName, "reply_to.email_address", "updated-reply@pingidentity.com"),
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
		CheckDestroy:             base.NotificationSettingsEmail_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			smtpFullSwitchStep,
			customProviderSwitchStep,
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("resource not found: %s", resourceFullName)
						}

						return rs.Primary.ID, nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password",
					"auth_token",
				},
			},
		},
	})
}

func testAccNotificationSettingsEmail_CustomProviderConfig_POST_RawBody_Full(environmentName, licenseID, resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_notification_settings_email" "%[3]s" {
  environment_id       = pingone_environment.%[2]s.id
  custom_provider_name = "UpdatedCustomProviderName"

  auth_token = "customauthtoken"

  protocol = "HTTP"

  from = {
    name          = "Updated Test Sender"
    email_address = "updated-no-reply@pingidentity.com"
  }

  reply_to = {
    name          = "Updated Test Reply To"
    email_address = "updated-reply@pingidentity.com"
  }

  requests = [
    {
      method = "POST"
      headers = {
        "Content-Type"     = "application/json"
        "updated-subject"  = "updated-$${subject}"
        "updated-reply-to" = "updated-$${reply_to}"
        "updated-from"     = "updated-$${from}"
      }
      body = <<EOF
{
  $${message},
  "to": [
    "$${to}"
  ],
  "toOverride": [
    "test@testaddress.com"
  ]
}
EOF
      url  = "https://api.pingidentity.com/updated-send-email"
    }
  ]
}`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName)
}

func TestAccNotificationSettingsEmail_CustomProvider_OAuth2(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_notification_settings_email.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()
	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	// Test switching to SMTP
	smtpFullSwitchStep := resource.TestStep{
		Config: testAccNotificationSettingsEmail_SMTPConfig_Full(environmentName, licenseID, resourceName),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "host", "pingidentity.com"),
			resource.TestCheckResourceAttr(resourceFullName, "port", "25"),
			resource.TestCheckResourceAttr(resourceFullName, "protocol", "SMTPS"),
			resource.TestCheckResourceAttr(resourceFullName, "username", "smtpuser"),
			resource.TestCheckResourceAttr(resourceFullName, "password", "smtpuserpassword"),
			resource.TestCheckResourceAttr(resourceFullName, "from.email_address", "noreply@pingidentity.com"),
			resource.TestCheckResourceAttr(resourceFullName, "from.name", "Stubbed From Address"),
			resource.TestCheckResourceAttr(resourceFullName, "reply_to.email_address", "reply@pingidentity.com"),
			resource.TestCheckResourceAttr(resourceFullName, "reply_to.name", "Stubbed Reply To Address"),
		),
	}

	minimalStep := resource.TestStep{
		Config: testAccNotificationSettingsEmail_CustomProviderConfig_OAuth2_Minimal(environmentName, licenseID, resourceName),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "custom_provider_name", "CustomProviderName"),
			resource.TestCheckResourceAttr(resourceFullName, "protocol", "HTTP"),
			resource.TestCheckResourceAttr(resourceFullName, "auth_url", "https://auth.pingone.com/oauth2/token"),
			resource.TestCheckResourceAttr(resourceFullName, "client_id", "testclientid"),
			resource.TestCheckResourceAttr(resourceFullName, "client_secret", "testclientsecret"),
			resource.TestCheckResourceAttr(resourceFullName, "grant_type", "CLIENT_CREDENTIALS"),
			resource.TestCheckResourceAttr(resourceFullName, "client_authentication_method", "BASIC_AUTH_HEADER"),
			resource.TestCheckNoResourceAttr(resourceFullName, "assertion"),
			resource.TestCheckNoResourceAttr(resourceFullName, "scopes.#"),
			resource.TestCheckResourceAttr(resourceFullName, "requests.#", "1"),
			resource.TestCheckResourceAttr(resourceFullName, "requests.0.method", "POST"),
			resource.TestCheckResourceAttr(resourceFullName, "requests.0.url", "https://api.pingidentity.com/send-email"),
			resource.TestCheckNoResourceAttr(resourceFullName, "from.name"),
			resource.TestCheckNoResourceAttr(resourceFullName, "reply_to.name"),
			resource.TestCheckResourceAttr(resourceFullName, "requests.0.headers.Content-Type", "application/x-www-form-urlencoded"),
			resource.TestCheckResourceAttr(resourceFullName, "requests.0.headers.subject", "${subject}"),
			resource.TestCheckResourceAttr(resourceFullName, "requests.0.headers.reply-to", "${reply_to}"),
			resource.TestCheckResourceAttr(resourceFullName, "requests.0.headers.from", "${from}"),
			resource.TestCheckResourceAttr(resourceFullName, "requests.0.body", "to=${to}&message=${message}"),
		),
	}

	minimalStepDestroy := resource.TestStep{
		Config:  minimalStep.Config,
		Destroy: true,
	}

	fullStep := resource.TestStep{
		Config: testAccNotificationSettingsEmail_CustomProviderConfig_OAuth2_Full(environmentName, licenseID, resourceName),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "custom_provider_name", "UpdatedCustomProviderName"),
			resource.TestCheckResourceAttr(resourceFullName, "protocol", "HTTP"),
			resource.TestCheckResourceAttr(resourceFullName, "auth_url", "https://auth.pingone.com/updated-oauth2/token"),
			resource.TestCheckResourceAttr(resourceFullName, "client_id", "updatedtestclientid"),
			resource.TestCheckResourceAttr(resourceFullName, "client_secret", "updatedtestclientsecret"),
			resource.TestCheckResourceAttr(resourceFullName, "grant_type", "CLIENT_CREDENTIALS"),
			resource.TestCheckResourceAttr(resourceFullName, "client_authentication_method", "BASIC_AUTH_HEADER"),
			resource.TestCheckResourceAttr(resourceFullName, "scopes.#", "2"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "scopes.*", "mail.send"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "scopes.*", "mail.read"),
			resource.TestCheckResourceAttr(resourceFullName, "from.name", "Updated Test Sender"),
			resource.TestCheckResourceAttr(resourceFullName, "reply_to.name", "Updated Test Reply To"),
			resource.TestCheckResourceAttr(resourceFullName, "requests.#", "1"),
			resource.TestCheckResourceAttr(resourceFullName, "requests.0.method", "POST"),
			resource.TestCheckResourceAttr(resourceFullName, "requests.0.url", "https://api.pingidentity.com/updated-send-email"),
			resource.TestCheckResourceAttr(resourceFullName, "requests.0.headers.Content-Type", "application/x-www-form-urlencoded"),
			resource.TestCheckResourceAttr(resourceFullName, "requests.0.headers.updated-subject", "updated-${subject}"),
			resource.TestCheckResourceAttr(resourceFullName, "requests.0.headers.updated-reply-to", "updated-${reply_to}"),
			resource.TestCheckResourceAttr(resourceFullName, "requests.0.headers.updated-from", "updated-${from}"),
			resource.TestCheckResourceAttr(resourceFullName, "requests.0.body", "updated-to=${to}&message=${message}"),
		),
	}

	fullStepDestroy := resource.TestStep{
		Config:  fullStep.Config,
		Destroy: true,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.NotificationSettingsEmail_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			minimalStep,
			smtpFullSwitchStep,
			minimalStep,
			minimalStepDestroy,
			fullStep,
			smtpFullSwitchStep,
			minimalStep,
			fullStep,
			fullStepDestroy,
			fullStep,
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("resource not found: %s", resourceFullName)
						}

						return rs.Primary.ID, nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password",
					"auth_token",
					"client_secret",
				},
			},
		},
	})
}

func testAccNotificationSettingsEmail_CustomProviderConfig_OAuth2_Minimal(environmentName, licenseID, resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_notification_settings_email" "%[3]s" {
  environment_id       = pingone_environment.%[2]s.id
  custom_provider_name = "CustomProviderName"

  auth_url      = "https://auth.pingone.com/oauth2/token"
  client_id     = "testclientid"
  client_secret = "testclientsecret"

  protocol = "HTTP"

  from = {
    email_address = "no-reply@pingidentity.com"
  }

  reply_to = {
    email_address = "reply@pingidentity.com"
  }

  requests = [
    {
      method = "POST"
      headers = {
        "Content-Type" = "application/x-www-form-urlencoded"
        "subject"      = "$${subject}"
        "reply-to"     = "$${reply_to}"
        "from"         = "$${from}"
      }
      body = "to=$${to}&message=$${message}"
      url  = "https://api.pingidentity.com/send-email"
    }
  ]
}`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName)
}

func testAccNotificationSettingsEmail_CustomProviderConfig_OAuth2_Full(environmentName, licenseID, resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_notification_settings_email" "%[3]s" {
  environment_id       = pingone_environment.%[2]s.id
  custom_provider_name = "UpdatedCustomProviderName"

  auth_url      = "https://auth.pingone.com/updated-oauth2/token"
  client_id     = "updatedtestclientid"
  client_secret = "updatedtestclientsecret"
  grant_type    = "CLIENT_CREDENTIALS"
  scopes        = ["mail.send", "mail.read"]

  protocol = "HTTP"

  from = {
    name          = "Updated Test Sender"
    email_address = "updated-no-reply@pingidentity.com"
  }

  reply_to = {
    name          = "Updated Test Reply To"
    email_address = "updated-reply@pingidentity.com"
  }

  requests = [
    {
      method = "POST"
      headers = {
        "Content-Type"     = "application/x-www-form-urlencoded"
        "updated-subject"  = "updated-$${subject}"
        "updated-reply-to" = "updated-$${reply_to}"
        "updated-from"     = "updated-$${from}"
      }
      body = "updated-to=$${to}&message=$${message}"
      url  = "https://api.pingidentity.com/updated-send-email"
    }
  ]
}`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName)
}

const (
	testAccNotificationSettingsEmailJWTAssertion        = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJ0ZXN0Y2xpZW50aWQiLCJzdWIiOiJ0ZXN0Y2xpZW50aWQiLCJhdWQiOiJodHRwczovL2F1dGgucGluZ29uZS5jb20vb2F1dGgyL3Rva2VuIiwiZXhwIjo0MTAyNDQ0ODAwfQ.dGVzdC1zaWduYXR1cmU"
	testAccNotificationSettingsEmailJWTAssertionUpdated = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJ1cGRhdGVkdGVzdGNsaWVudGlkIiwic3ViIjoidXBkYXRlZHRlc3RjbGllbnRpZCIsImF1ZCI6Imh0dHBzOi8vYXV0aC5waW5nb25lLmNvbS91cGRhdGVkLW9hdXRoMi90b2tlbiIsImV4cCI6NDEwMjQ0NDgwMH0.dGVzdC1zaWduYXR1cmU"
)

func TestAccNotificationSettingsEmail_CustomProvider_OAuth2_JWTBearer(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_notification_settings_email.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()
	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	minimalStep := resource.TestStep{
		Config: testAccNotificationSettingsEmail_CustomProviderConfig_Authentication(environmentName, licenseID, resourceName, fmt.Sprintf(`
  auth_url   = "https://auth.pingone.com/oauth2/token"
  grant_type = "JWT_BEARER"
  assertion  = "%s"
`, testAccNotificationSettingsEmailJWTAssertion)),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "custom_provider_name", "CustomProviderName"),
			resource.TestCheckResourceAttr(resourceFullName, "protocol", "HTTP"),
			resource.TestCheckResourceAttr(resourceFullName, "auth_url", "https://auth.pingone.com/oauth2/token"),
			resource.TestCheckResourceAttr(resourceFullName, "grant_type", "JWT_BEARER"),
			resource.TestCheckResourceAttr(resourceFullName, "assertion", testAccNotificationSettingsEmailJWTAssertion),
			resource.TestCheckResourceAttr(resourceFullName, "client_authentication_method", "BODY"),
			resource.TestCheckNoResourceAttr(resourceFullName, "client_id"),
			resource.TestCheckNoResourceAttr(resourceFullName, "client_secret"),
			resource.TestCheckNoResourceAttr(resourceFullName, "scopes.#"),
		),
	}

	minimalStepDestroy := resource.TestStep{
		Config:  minimalStep.Config,
		Destroy: true,
	}

	fullStep := resource.TestStep{
		Config: testAccNotificationSettingsEmail_CustomProviderConfig_Authentication(environmentName, licenseID, resourceName, fmt.Sprintf(`
  auth_url   = "https://auth.pingone.com/updated-oauth2/token"
  grant_type = "JWT_BEARER"
  assertion  = "%s"
  scopes     = ["mail.send", "mail.read"]
`, testAccNotificationSettingsEmailJWTAssertionUpdated)),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "auth_url", "https://auth.pingone.com/updated-oauth2/token"),
			resource.TestCheckResourceAttr(resourceFullName, "grant_type", "JWT_BEARER"),
			resource.TestCheckResourceAttr(resourceFullName, "assertion", testAccNotificationSettingsEmailJWTAssertionUpdated),
			resource.TestCheckResourceAttr(resourceFullName, "client_authentication_method", "BODY"),
			resource.TestCheckResourceAttr(resourceFullName, "scopes.#", "2"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "scopes.*", "mail.send"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "scopes.*", "mail.read"),
			resource.TestCheckNoResourceAttr(resourceFullName, "client_id"),
			resource.TestCheckNoResourceAttr(resourceFullName, "client_secret"),
		),
	}

	// Test switching to the client credentials grant type
	clientCredentialsSwitchStep := resource.TestStep{
		Config: testAccNotificationSettingsEmail_CustomProviderConfig_OAuth2_Minimal(environmentName, licenseID, resourceName),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "auth_url", "https://auth.pingone.com/oauth2/token"),
			resource.TestCheckResourceAttr(resourceFullName, "grant_type", "CLIENT_CREDENTIALS"),
			resource.TestCheckResourceAttr(resourceFullName, "client_id", "testclientid"),
			resource.TestCheckResourceAttr(resourceFullName, "client_secret", "testclientsecret"),
			resource.TestCheckResourceAttr(resourceFullName, "client_authentication_method", "BASIC_AUTH_HEADER"),
			resource.TestCheckNoResourceAttr(resourceFullName, "assertion"),
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
		CheckDestroy:             base.NotificationSettingsEmail_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			minimalStep,
			fullStep,
			clientCredentialsSwitchStep,
			minimalStep,
			minimalStepDestroy,
			fullStep,
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("resource not found: %s", resourceFullName)
						}

						return rs.Primary.ID, nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"assertion",
				},
			},
		},
	})
}

func TestAccNotificationSettingsEmail_CustomProvider_CustomHeader(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_notification_settings_email.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()
	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	// Test switching to SMTP
	smtpFullSwitchStep := resource.TestStep{
		Config: testAccNotificationSettingsEmail_SMTPConfig_Full(environmentName, licenseID, resourceName),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "host", "pingidentity.com"),
			resource.TestCheckResourceAttr(resourceFullName, "username", "smtpuser"),
			resource.TestCheckNoResourceAttr(resourceFullName, "header_name"),
			resource.TestCheckNoResourceAttr(resourceFullName, "header_value"),
		),
	}

	minimalStep := resource.TestStep{
		Config: testAccNotificationSettingsEmail_CustomProviderConfig_Authentication(environmentName, licenseID, resourceName, `
  header_name  = "X-Api-Key"
  header_value = "testheadervalue"
`),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "custom_provider_name", "CustomProviderName"),
			resource.TestCheckResourceAttr(resourceFullName, "protocol", "HTTP"),
			resource.TestCheckResourceAttr(resourceFullName, "header_name", "X-Api-Key"),
			resource.TestCheckResourceAttr(resourceFullName, "header_value", "testheadervalue"),
			resource.TestCheckNoResourceAttr(resourceFullName, "auth_url"),
			resource.TestCheckNoResourceAttr(resourceFullName, "auth_token"),
			resource.TestCheckNoResourceAttr(resourceFullName, "username"),
		),
	}

	minimalStepDestroy := resource.TestStep{
		Config:  minimalStep.Config,
		Destroy: true,
	}

	fullStep := resource.TestStep{
		Config: testAccNotificationSettingsEmail_CustomProviderConfig_Authentication(environmentName, licenseID, resourceName, `
  header_name  = "Authorization-Key"
  header_value = "updatedtestheadervalue"
`),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "header_name", "Authorization-Key"),
			resource.TestCheckResourceAttr(resourceFullName, "header_value", "updatedtestheadervalue"),
		),
	}

	// Test switching from OAuth 2.0
	oauth2SwitchStep := resource.TestStep{
		Config: testAccNotificationSettingsEmail_CustomProviderConfig_OAuth2_Minimal(environmentName, licenseID, resourceName),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "auth_url", "https://auth.pingone.com/oauth2/token"),
			resource.TestCheckNoResourceAttr(resourceFullName, "header_name"),
			resource.TestCheckNoResourceAttr(resourceFullName, "header_value"),
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
		CheckDestroy:             base.NotificationSettingsEmail_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			minimalStep,
			fullStep,
			smtpFullSwitchStep,
			minimalStep,
			oauth2SwitchStep,
			minimalStep,
			minimalStepDestroy,
			fullStep,
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("resource not found: %s", resourceFullName)
						}

						return rs.Primary.ID, nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"header_value",
				},
			},
		},
	})
}

func TestAccNotificationSettingsEmail_CustomProvider_InvalidAuthentication(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	environmentName := acctest.ResourceNameGenEnvironment()
	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	invalidStep := func(authenticationHCL string, expectError *regexp.Regexp) resource.TestStep {
		return resource.TestStep{
			Config:      testAccNotificationSettingsEmail_CustomProviderConfig_Authentication(environmentName, licenseID, resourceName, authenticationHCL),
			ExpectError: expectError,
		}
	}

	missingAttribute := regexp.MustCompile(`Missing Attribute Configuration`)
	invalidCombination := regexp.MustCompile(`Invalid Attribute Combination`)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.NotificationSettingsEmail_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// JWT_BEARER requires `assertion`
			invalidStep(`
  auth_url   = "https://auth.pingone.com/oauth2/token"
  grant_type = "JWT_BEARER"
`, missingAttribute),
			// The default grant type (CLIENT_CREDENTIALS) requires `client_secret`
			invalidStep(`
  auth_url  = "https://auth.pingone.com/oauth2/token"
  client_id = "testclientid"
`, missingAttribute),
			// The default grant type (CLIENT_CREDENTIALS) requires `client_id`
			invalidStep(`
  auth_url      = "https://auth.pingone.com/oauth2/token"
  client_secret = "testclientsecret"
`, missingAttribute),
			// CLIENT_CREDENTIALS cannot be combined with `assertion`
			invalidStep(fmt.Sprintf(`
  auth_url      = "https://auth.pingone.com/oauth2/token"
  client_id     = "testclientid"
  client_secret = "testclientsecret"
  assertion     = "%s"
`, testAccNotificationSettingsEmailJWTAssertion), invalidCombination),
			// JWT_BEARER cannot be combined with `client_id`
			invalidStep(fmt.Sprintf(`
  auth_url   = "https://auth.pingone.com/oauth2/token"
  grant_type = "JWT_BEARER"
  assertion  = "%s"
  client_id  = "testclientid"
`, testAccNotificationSettingsEmailJWTAssertion), invalidCombination),
			// Invalid grant type
			invalidStep(`
  auth_url      = "https://auth.pingone.com/oauth2/token"
  grant_type    = "INVALID"
  client_id     = "testclientid"
  client_secret = "testclientsecret"
`, regexp.MustCompile(`Invalid Attribute Value Match`)),
			// OAuth 2.0 attributes require `auth_url`
			invalidStep(`
  auth_token = "testauthtoken"
  scopes     = ["mail.send"]
`, invalidCombination),
			// `header_name` requires `header_value`
			invalidStep(`
  header_name = "X-Api-Key"
`, invalidCombination),
			// `header_value` requires `header_name`
			invalidStep(`
  auth_token   = "testauthtoken"
  header_value = "testheadervalue"
`, invalidCombination),
			// Only one authentication method can be configured
			invalidStep(`
  header_name  = "X-Api-Key"
  header_value = "testheadervalue"
  auth_token   = "testauthtoken"
`, invalidCombination),
			invalidStep(`
  auth_url      = "https://auth.pingone.com/oauth2/token"
  client_id     = "testclientid"
  client_secret = "testclientsecret"
  header_name   = "X-Api-Key"
  header_value  = "testheadervalue"
`, invalidCombination),
		},
	})
}

func testAccNotificationSettingsEmail_CustomProviderConfig_Authentication(environmentName, licenseID, resourceName, authenticationHCL string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_notification_settings_email" "%[3]s" {
  environment_id       = pingone_environment.%[2]s.id
  custom_provider_name = "CustomProviderName"
  protocol             = "HTTP"
  %[4]s
  from = {
    email_address = "no-reply@pingidentity.com"
  }

  reply_to = {
    email_address = "reply@pingidentity.com"
  }

  requests = [
    {
      method = "POST"
      headers = {
        "Content-Type" = "application/x-www-form-urlencoded"
        "subject"      = "$${subject}"
        "reply-to"     = "$${reply_to}"
        "from"         = "$${from}"
      }
      body = "to=$${to}&message=$${message}"
      url  = "https://api.pingidentity.com/send-email"
    }
  ]
}`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, authenticationHCL)
}
