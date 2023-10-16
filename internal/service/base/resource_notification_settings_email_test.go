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

func TestAccNotificationSettingsEmail_RemovalDrift(t *testing.T) {
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
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)

			p1Client = acctest.PreCheckTestClient(ctx, t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.NotificationSettingsEmail_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Test removal of the environment
			{
				Config: testAccNotificationSettingsEmailConfig_Full(environmentName, licenseID, resourceName),
				Check:  base.NotificationSettingsEmail_GetIDs(resourceFullName, &environmentID),
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

func TestAccNotificationSettingsEmail_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_notification_settings_email.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	fullStep := resource.TestStep{
		Config: testAccNotificationSettingsEmailConfig_Full(environmentName, licenseID, resourceName),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "host", "smtp-example.pingidentity.com"),
			resource.TestCheckResourceAttr(resourceFullName, "port", "25"),
			resource.TestCheckResourceAttr(resourceFullName, "protocol", "SMTPS"),
			resource.TestCheckResourceAttr(resourceFullName, "username", "smtpuser"),
			resource.TestCheckResourceAttr(resourceFullName, "password", "smtpuserpassword"),
			resource.TestCheckResourceAttr(resourceFullName, "from.#", "1"),
			resource.TestCheckResourceAttr(resourceFullName, "reply_to.#", "1"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
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
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return rs.Primary.ID, nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password",
				},
			},
		},
	})
}

func TestAccNotificationSettingsEmail_EmailSources(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_notification_settings_email.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	fromFull := resource.TestStep{
		Config: testAccNotificationSettingsEmailConfig_FromFull(environmentName, licenseID, resourceName),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "from.#", "1"),
			resource.TestCheckResourceAttr(resourceFullName, "from.0.email_address", "noreply@pingidentity.com"),
			resource.TestCheckResourceAttr(resourceFullName, "from.0.name", "Stubbed From Address"),
			resource.TestCheckResourceAttr(resourceFullName, "reply_to.#", "0"),
		),
	}

	fromMinimal := resource.TestStep{
		Config: testAccNotificationSettingsEmailConfig_FromMinimal(environmentName, licenseID, resourceName),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "from.#", "1"),
			resource.TestCheckResourceAttr(resourceFullName, "from.0.email_address", "noreply@pingidentity.com"),
			resource.TestCheckNoResourceAttr(resourceFullName, "from.0.name"),
			resource.TestCheckResourceAttr(resourceFullName, "reply_to.#", "0"),
		),
	}

	replyToFull := resource.TestStep{
		Config: testAccNotificationSettingsEmailConfig_ReplyToFull(environmentName, licenseID, resourceName),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "from.#", "1"),
			resource.TestCheckResourceAttr(resourceFullName, "from.0.email_address", "noreply@pingidentity.com"),
			resource.TestCheckNoResourceAttr(resourceFullName, "from.0.name"),
			resource.TestCheckResourceAttr(resourceFullName, "reply_to.#", "1"),
			resource.TestCheckResourceAttr(resourceFullName, "reply_to.0.email_address", "reply@pingidentity.com"),
			resource.TestCheckResourceAttr(resourceFullName, "reply_to.0.name", "Stubbed Reply To Address"),
		),
	}

	replyToMinimal := resource.TestStep{
		Config: testAccNotificationSettingsEmailConfig_ReplyToMinimal(environmentName, licenseID, resourceName),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "from.#", "1"),
			resource.TestCheckResourceAttr(resourceFullName, "from.0.email_address", "noreply@pingidentity.com"),
			resource.TestCheckNoResourceAttr(resourceFullName, "from.0.name"),
			resource.TestCheckResourceAttr(resourceFullName, "reply_to.#", "1"),
			resource.TestCheckResourceAttr(resourceFullName, "reply_to.0.email_address", "reply@pingidentity.com"),
			resource.TestCheckNoResourceAttr(resourceFullName, "reply_to.0.name"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
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
				Config:  testAccNotificationSettingsEmailConfig_FromFull(environmentName, licenseID, resourceName),
				Destroy: true,
			},
			// Variant 2 New
			fromMinimal,
			{
				Config:  testAccNotificationSettingsEmailConfig_FromMinimal(environmentName, licenseID, resourceName),
				Destroy: true,
			},
			// Variant 3 New
			replyToFull,
			{
				Config:  testAccNotificationSettingsEmailConfig_ReplyToFull(environmentName, licenseID, resourceName),
				Destroy: true,
			},
			// Variant 3 New
			replyToMinimal,
			{
				Config:  testAccNotificationSettingsEmailConfig_ReplyToMinimal(environmentName, licenseID, resourceName),
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

func TestAccNotificationSettingsEmail_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_notification_settings_email.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.NotificationSettingsEmail_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccNotificationSettingsEmailConfig_Full(environmentName, licenseID, resourceName),
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

func testAccNotificationSettingsEmailConfig_Full(environmentName, licenseID, resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_notification_settings_email" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  host     = "smtp-example.pingidentity.com"
  port     = 25
  username = "smtpuser"
  password = "smtpuserpassword"

  from {
    email_address = "noreply@pingidentity.com"
  }

  reply_to {
    email_address = "reply@pingidentity.com"
  }
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName)
}

func testAccNotificationSettingsEmailConfig_FromFull(environmentName, licenseID, resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_notification_settings_email" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  host     = "smtp-example.pingidentity.com"
  port     = 25
  username = "smtpuser"
  password = "smtpuserpassword"

  from {
    email_address = "noreply@pingidentity.com"
    name          = "Stubbed From Address"
  }
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName)
}

func testAccNotificationSettingsEmailConfig_FromMinimal(environmentName, licenseID, resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_notification_settings_email" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  host     = "smtp-example.pingidentity.com"
  port     = 25
  username = "smtpuser"
  password = "smtpuserpassword"

  from {
    email_address = "noreply@pingidentity.com"
  }
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName)
}

func testAccNotificationSettingsEmailConfig_ReplyToFull(environmentName, licenseID, resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_notification_settings_email" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  host     = "smtp-example.pingidentity.com"
  port     = 25
  username = "smtpuser"
  password = "smtpuserpassword"

  from {
    email_address = "noreply@pingidentity.com"
  }

  reply_to {
    email_address = "reply@pingidentity.com"
    name          = "Stubbed Reply To Address"
  }
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName)
}

func testAccNotificationSettingsEmailConfig_ReplyToMinimal(environmentName, licenseID, resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_notification_settings_email" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  host     = "smtp-example.pingidentity.com"
  port     = 25
  username = "smtpuser"
  password = "smtpuserpassword"

  from {
    email_address = "noreply@pingidentity.com"
  }

  reply_to {
    email_address = "reply@pingidentity.com"
  }
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName)
}
