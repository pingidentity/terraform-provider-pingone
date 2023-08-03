package base_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func testAccCheckNotificationSettingsEmailDestroy(s *terraform.State) error {
	return nil
}

func testAccGetMFAPolicyIDs(resourceName string, environmentID, resourceID *string) resource.TestCheckFunc {
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
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckMFAPolicyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccMFAPolicyConfig_FullSMS(resourceName, name),
				Check:  testAccGetMFAPolicyIDs(resourceFullName, &environmentID, &resourceID),
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

func TestAccNotificationSettingsEmail_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_notification_settings_email.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	fullStep := resource.TestStep{
		Config: testAccNotificationSettingsEmailConfig_Full(environmentName, licenseID, resourceName),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
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
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckNotificationSettingsEmailDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full from scratch
			fullStep,
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
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckNotificationSettingsEmailDestroy,
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
