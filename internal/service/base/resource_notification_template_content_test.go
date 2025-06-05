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

func TestAccNotificationTemplateContent_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_notification_template_content.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	name := "strong_authentication"
	locale := "en-GB"

	var notificationTemplateContentID, templateName, environmentID string

	var p1Client *client.Client
	var ctx = context.Background()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)

			p1Client = acctest.PreCheckTestClient(ctx, t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.NotificationTemplateContent_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccNotificationTemplateContentConfig_NewLocale_Minimal(environmentName, licenseID, resourceName, name, locale),
				Check:  base.NotificationTemplateContent_GetIDs(resourceFullName, &environmentID, &templateName, &notificationTemplateContentID),
			},
			{
				PreConfig: func() {
					base.NotificationTemplateContent_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID, templateName, notificationTemplateContentID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the environment
			{
				Config: testAccNotificationTemplateContentConfig_NewLocale_Minimal(environmentName, licenseID, resourceName, name, locale),
				Check:  base.NotificationTemplateContent_GetIDs(resourceFullName, &environmentID, &templateName, &notificationTemplateContentID),
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

func TestAccNotificationTemplateContent_OverrideDefaultLocale(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_notification_template_content.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	name := "strong_authentication"
	locale := "en"

	check := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "template_name", name),
		resource.TestCheckResourceAttr(resourceFullName, "locale", locale),
		resource.TestCheckResourceAttr(resourceFullName, "default", "false"),
		resource.TestCheckNoResourceAttr(resourceFullName, "variant"),
		resource.TestCheckNoResourceAttr(resourceFullName, "email"),
		resource.TestCheckResourceAttrSet(resourceFullName, "push.body"),
		resource.TestCheckResourceAttrSet(resourceFullName, "push.title"),
		resource.TestCheckNoResourceAttr(resourceFullName, "sms"),
		resource.TestCheckNoResourceAttr(resourceFullName, "voice"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.NotificationTemplateContent_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNotificationTemplateContentConfig_DefaultVariant_Push_Minimal(environmentName, licenseID, resourceName, name, locale),
				Check:  check,
			},
			// We destroy and retest due to config bootstrapping
			{
				Config:  testAccNotificationTemplateContentConfig_DefaultVariant_Push_Minimal(environmentName, licenseID, resourceName, name, locale),
				Destroy: true,
			},
			{
				Config: testAccNotificationTemplateContentConfig_DefaultVariant_Push_Minimal(environmentName, licenseID, resourceName, name, locale),
				Check:  check,
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

						return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["template_name"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNotificationTemplateContent_NewLocale(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_notification_template_content.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	name := "strong_authentication"
	locale := "en-GB"

	check := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "template_name", name),
		resource.TestCheckResourceAttr(resourceFullName, "locale", locale),
		resource.TestCheckResourceAttr(resourceFullName, "default", "false"),
		resource.TestCheckNoResourceAttr(resourceFullName, "variant"),
		resource.TestCheckNoResourceAttr(resourceFullName, "email"),
		resource.TestCheckResourceAttr(resourceFullName, "push.body", "Min - Please approve this transaction."),
		resource.TestCheckResourceAttr(resourceFullName, "push.title", "Min - BX Retail Transaction Request"),
		resource.TestCheckNoResourceAttr(resourceFullName, "sms"),
		resource.TestCheckNoResourceAttr(resourceFullName, "voice"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.NotificationTemplateContent_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNotificationTemplateContentConfig_NewLocale_Minimal(environmentName, licenseID, resourceName, name, locale),
				Check:  check,
			},
			// We destroy and retest due to config bootstrapping
			{
				Config:  testAccNotificationTemplateContentConfig_NewLocale_Minimal(environmentName, licenseID, resourceName, name, locale),
				Destroy: true,
			},
			{
				Config: testAccNotificationTemplateContentConfig_NewLocale_Minimal(environmentName, licenseID, resourceName, name, locale),
				Check:  check,
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

						return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["template_name"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Languages not added but the local used should result in error
			{
				Config:  testAccNotificationTemplateContentConfig_NewLocale_Minimal(environmentName, licenseID, resourceName, name, locale),
				Destroy: true,
			},
			{
				Config:      testAccNotificationTemplateContentConfig_DefaultVariant_Push_Minimal(environmentName, licenseID, resourceName, name, locale),
				ExpectError: regexp.MustCompile("The locale is not valid for the environment."),
			},
		},
	})
}

func TestAccNotificationTemplateContent_NewVariant(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_notification_template_content.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	name := "verification_code_template"
	locale := "en"
	variant := "My New Variant"

	check := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "template_name", name),
		resource.TestCheckResourceAttr(resourceFullName, "locale", locale),
		resource.TestCheckResourceAttr(resourceFullName, "default", "false"),
		resource.TestCheckResourceAttr(resourceFullName, "variant", variant),
		resource.TestCheckResourceAttrSet(resourceFullName, "email.body"),
		resource.TestCheckResourceAttrSet(resourceFullName, "email.subject"),
		resource.TestCheckNoResourceAttr(resourceFullName, "push"),
		resource.TestCheckNoResourceAttr(resourceFullName, "sms"),
		resource.TestCheckNoResourceAttr(resourceFullName, "voice"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.NotificationTemplateContent_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNotificationTemplateContentConfig_NewVariant_Minimal(environmentName, licenseID, resourceName, name, locale, variant),
				Check:  check,
			},
			// We destroy and retest due to config bootstrapping
			{
				Config:  testAccNotificationTemplateContentConfig_NewVariant_Minimal(environmentName, licenseID, resourceName, name, locale, variant),
				Destroy: true,
			},
			{
				Config: testAccNotificationTemplateContentConfig_NewVariant_Minimal(environmentName, licenseID, resourceName, name, locale, variant),
				Check:  check,
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

						return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["template_name"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNotificationTemplateContent_ChangeVariant(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_notification_template_content.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	name := "verification_code_template"
	locale := "en"

	variant1 := "My New Variant"
	check1 := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "variant", variant1),
	)

	variant2 := "My Second Variant"
	check2 := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "variant", variant2),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.NotificationTemplateContent_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Change defined variants
			{
				Config: testAccNotificationTemplateContentConfig_NewVariant_Minimal(environmentName, licenseID, resourceName, name, locale, variant1),
				Check:  check1,
			},
			{
				Config: testAccNotificationTemplateContentConfig_NewVariant_Minimal(environmentName, licenseID, resourceName, name, locale, variant2),
				Check:  check2,
			},
			{
				Config:  testAccNotificationTemplateContentConfig_NewVariant_Minimal(environmentName, licenseID, resourceName, name, locale, variant2),
				Destroy: true,
			},
			// From no variant, to variant
			{
				Config: testAccNotificationTemplateContentConfig_NoVariant_Minimal(environmentName, licenseID, resourceName, name, locale),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr(resourceFullName, "variant"),
				),
			},
			{
				Config: testAccNotificationTemplateContentConfig_NewVariant_Minimal(environmentName, licenseID, resourceName, name, locale, variant2),
				Check:  check2,
			},
		},
	})
}

func TestAccNotificationTemplateContent_InvalidData(t *testing.T) {
	t.Parallel()

	environmentName := acctest.ResourceNameGenEnvironment()
	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resourceName := acctest.ResourceNameGen()

	name := acctest.TestData{
		Invalid: "strong_authentication_doesnotexist",
		Valid:   "strong_authentication",
	}
	locale := acctest.TestData{
		Invalid: "en-ZZ",
		Valid:   "en-GB",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.NotificationTemplateContent_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccNotificationTemplateContentConfig_DefaultVariant_Push_Minimal(environmentName, licenseID, resourceName, name.Invalid, locale.Valid),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Match`),
			},
			{
				Config:      testAccNotificationTemplateContentConfig_DefaultVariant_Push_Minimal(environmentName, licenseID, resourceName, name.Valid, locale.Invalid),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Match`),
			},
			{
				Config:      testAccNotificationTemplateContentConfig_DuplicateLocale(environmentName, licenseID, resourceName, name.Valid, "en"),
				ExpectError: regexp.MustCompile(`Customized content for the template, locale and variant combination already exists.`),
			},
		},
	})
}

func TestAccNotificationTemplateContent_Email(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_notification_template_content.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	locale := "en"

	template := acctest.TestData{
		//Invalid: "", // Invalid not tested, no templates without email
		Valid: "email_verification_admin",
	}

	check := acctest.MinMaxChecks{
		Minimal: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "template_name", template.Valid),
			resource.TestCheckResourceAttr(resourceFullName, "locale", locale),
			resource.TestCheckResourceAttr(resourceFullName, "email.subject", "Min - Verify your email address"),
			resource.TestCheckResourceAttr(resourceFullName, "email.body", "Min - <p>\n\tUse the code provided to verify your email address with PingOne. If you think you received this email in error, contact your supervisor.\n\t</p>\n\t<p>\n\tTo verify your email address:\n\t<ol type=\"1\">\n\t  <li>Sign-on to the self-service portal. For instructions, see <a href=\"https://docs.pingidentity.com/bundle/pingone/page/snd1631892368614.html\">Managing your PingOne user profile</a>.</li>\n\t  <li>Enter your verification code in the <b>Contact</b> section: ${code}</li>\n\t  <li>Click <b>Verify</b>.</li>\n\t</ol>\n\t</p>\n"),
			resource.TestCheckNoResourceAttr(resourceFullName, "email.from"),
			resource.TestCheckNoResourceAttr(resourceFullName, "email.reply_to"),
			resource.TestCheckResourceAttr(resourceFullName, "email.character_set", "UTF-8"),
			resource.TestCheckResourceAttr(resourceFullName, "email.content_type", "text/html"),
			resource.TestCheckNoResourceAttr(resourceFullName, "push"),
			resource.TestCheckNoResourceAttr(resourceFullName, "sms"),
			resource.TestCheckNoResourceAttr(resourceFullName, "voice"),
		),
		Full: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "template_name", template.Valid),
			resource.TestCheckResourceAttr(resourceFullName, "locale", locale),
			resource.TestCheckResourceAttr(resourceFullName, "email.subject", "Full - Verify your email address"),
			resource.TestCheckResourceAttr(resourceFullName, "email.body", "Full - <p>\n\tUse the code provided to verify your email address with PingOne. If you think you received this email in error, contact your supervisor.\n\t</p>\n\t<p>\n\tTo verify your email address:\n\t<ol type=\"1\">\n\t  <li>Sign-on to the self-service portal. For instructions, see <a href=\"https://docs.pingidentity.com/bundle/pingone/page/snd1631892368614.html\">Managing your PingOne user profile</a>.</li>\n\t  <li>Enter your verification code in the <b>Contact</b> section: ${code}</li>\n\t  <li>Click <b>Verify</b>.</li>\n\t</ol>\n\t</p>\n"),
			// resource.TestCheckResourceAttr(resourceFullName, "email.from.#", "1"),
			// resource.TestCheckResourceAttr(resourceFullName, "email.from.name", "BX Retail"),
			// resource.TestCheckResourceAttr(resourceFullName, "email.from.address", "noreply@bxretail.org"),
			// resource.TestCheckResourceAttr(resourceFullName, "email.reply_to.#", "1"),
			// resource.TestCheckResourceAttr(resourceFullName, "email.reply_to.name", "BX Retail Reply"),
			// resource.TestCheckResourceAttr(resourceFullName, "email.reply_to.address", "reply@bxretail.org"),
			resource.TestCheckResourceAttr(resourceFullName, "email.character_set", "iso-8859-5"),
			resource.TestCheckResourceAttr(resourceFullName, "email.content_type", "text/plain"),
			resource.TestCheckNoResourceAttr(resourceFullName, "push"),
			resource.TestCheckNoResourceAttr(resourceFullName, "sms"),
			resource.TestCheckNoResourceAttr(resourceFullName, "voice"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.NotificationTemplateContent_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Minimal from new
			{
				Config: testAccNotificationTemplateContentConfig_DefaultVariant_Email_Minimal(environmentName, licenseID, resourceName, template.Valid, locale),
				Check:  check.Minimal,
			},
			{
				Config:  testAccNotificationTemplateContentConfig_DefaultVariant_Email_Minimal(environmentName, licenseID, resourceName, template.Valid, locale),
				Destroy: true,
			},
			// Full from new
			{
				Config: testAccNotificationTemplateContentConfig_DefaultVariant_Email_Full(environmentName, licenseID, resourceName, template.Valid, locale),
				Check:  check.Full,
			},
			{
				Config:  testAccNotificationTemplateContentConfig_DefaultVariant_Email_Full(environmentName, licenseID, resourceName, template.Valid, locale),
				Destroy: true,
			},
			// Update
			{
				Config: testAccNotificationTemplateContentConfig_DefaultVariant_Email_Full(environmentName, licenseID, resourceName, template.Valid, locale),
				Check:  check.Full,
			},
			{
				Config: testAccNotificationTemplateContentConfig_DefaultVariant_Email_Minimal(environmentName, licenseID, resourceName, template.Valid, locale),
				Check:  check.Minimal,
			},
			{
				Config: testAccNotificationTemplateContentConfig_DefaultVariant_Email_Full(environmentName, licenseID, resourceName, template.Valid, locale),
				Check:  check.Full,
			},
		},
	})
}

func TestAccNotificationTemplateContent_Push(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_notification_template_content.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	locale := "en"

	template := acctest.TestData{
		Valid:   "strong_authentication",
		Invalid: "email_verification_admin",
	}

	check := acctest.MinMaxChecks{
		Minimal: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "template_name", template.Valid),
			resource.TestCheckResourceAttr(resourceFullName, "locale", locale),
			resource.TestCheckNoResourceAttr(resourceFullName, "email"),
			resource.TestCheckResourceAttr(resourceFullName, "push.body", "Min - Please approve this transaction."),
			resource.TestCheckResourceAttr(resourceFullName, "push.title", "Min - BX Retail Transaction Request"),
			resource.TestCheckResourceAttr(resourceFullName, "push.category", "BANNER_BUTTONS"),
			resource.TestCheckNoResourceAttr(resourceFullName, "sms"),
			resource.TestCheckNoResourceAttr(resourceFullName, "voice"),
		),
		Full: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "template_name", template.Valid),
			resource.TestCheckResourceAttr(resourceFullName, "locale", locale),
			resource.TestCheckNoResourceAttr(resourceFullName, "email"),
			resource.TestCheckResourceAttr(resourceFullName, "push.body", "Full - Please approve this transaction."),
			resource.TestCheckResourceAttr(resourceFullName, "push.title", "Full - BX Retail Transaction Request"),
			resource.TestCheckResourceAttr(resourceFullName, "push.category", "WITHOUT_BANNER_BUTTONS"),
			resource.TestCheckNoResourceAttr(resourceFullName, "sms"),
			resource.TestCheckNoResourceAttr(resourceFullName, "voice"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.NotificationTemplateContent_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Minimal from new
			{
				Config: testAccNotificationTemplateContentConfig_DefaultVariant_Push_Minimal(environmentName, licenseID, resourceName, template.Valid, locale),
				Check:  check.Minimal,
			},
			{
				Config:  testAccNotificationTemplateContentConfig_DefaultVariant_Push_Minimal(environmentName, licenseID, resourceName, template.Valid, locale),
				Destroy: true,
			},
			// Full from new
			{
				Config: testAccNotificationTemplateContentConfig_DefaultVariant_Push_Full(environmentName, licenseID, resourceName, template.Valid, locale),
				Check:  check.Full,
			},
			{
				Config:  testAccNotificationTemplateContentConfig_DefaultVariant_Push_Full(environmentName, licenseID, resourceName, template.Valid, locale),
				Destroy: true,
			},
			// Update
			{
				Config: testAccNotificationTemplateContentConfig_DefaultVariant_Push_Full(environmentName, licenseID, resourceName, template.Valid, locale),
				Check:  check.Full,
			},
			{
				Config: testAccNotificationTemplateContentConfig_DefaultVariant_Push_Minimal(environmentName, licenseID, resourceName, template.Valid, locale),
				Check:  check.Minimal,
			},
			{
				Config: testAccNotificationTemplateContentConfig_DefaultVariant_Push_Full(environmentName, licenseID, resourceName, template.Valid, locale),
				Check:  check.Full,
			},
			{
				Config:  testAccNotificationTemplateContentConfig_DefaultVariant_Push_Full(environmentName, licenseID, resourceName, template.Valid, locale),
				Destroy: true,
			},
			// Bad method for template
			{
				Config:      testAccNotificationTemplateContentConfig_DefaultVariant_Push_Full(environmentName, licenseID, resourceName, template.Invalid, locale),
				ExpectError: regexp.MustCompile(`The configured delivery method does not apply to the selected template.`),
			},
		},
	})
}

func TestAccNotificationTemplateContent_SMS(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_notification_template_content.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	locale := "en"

	template := acctest.TestData{
		Valid:   "strong_authentication",
		Invalid: "email_verification_admin",
	}

	check := acctest.MinMaxChecks{
		Minimal: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "template_name", template.Valid),
			resource.TestCheckResourceAttr(resourceFullName, "locale", locale),
			resource.TestCheckNoResourceAttr(resourceFullName, "email"),
			resource.TestCheckNoResourceAttr(resourceFullName, "push"),
			resource.TestCheckResourceAttr(resourceFullName, "sms.content", "Min - Please approve this transaction with passcode ${otp}."),
			resource.TestCheckNoResourceAttr(resourceFullName, "sms.sender"),
			resource.TestCheckNoResourceAttr(resourceFullName, "voice"),
		),
		Full: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "template_name", template.Valid),
			resource.TestCheckResourceAttr(resourceFullName, "locale", locale),
			resource.TestCheckNoResourceAttr(resourceFullName, "email"),
			resource.TestCheckNoResourceAttr(resourceFullName, "push"),
			resource.TestCheckResourceAttr(resourceFullName, "sms.content", "Full - Please approve this transaction with passcode ${otp}."),
			resource.TestCheckResourceAttr(resourceFullName, "sms.sender", "BX Retail"),
			resource.TestCheckNoResourceAttr(resourceFullName, "voice"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.NotificationTemplateContent_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Minimal from new
			{
				Config: testAccNotificationTemplateContentConfig_DefaultVariant_SMS_Minimal(environmentName, licenseID, resourceName, template.Valid, locale),
				Check:  check.Minimal,
			},
			{
				Config:  testAccNotificationTemplateContentConfig_DefaultVariant_SMS_Minimal(environmentName, licenseID, resourceName, template.Valid, locale),
				Destroy: true,
			},
			// Full from new
			{
				Config: testAccNotificationTemplateContentConfig_DefaultVariant_SMS_Full(environmentName, licenseID, resourceName, template.Valid, locale),
				Check:  check.Full,
			},
			{
				Config:  testAccNotificationTemplateContentConfig_DefaultVariant_SMS_Full(environmentName, licenseID, resourceName, template.Valid, locale),
				Destroy: true,
			},
			// Update
			{
				Config: testAccNotificationTemplateContentConfig_DefaultVariant_SMS_Full(environmentName, licenseID, resourceName, template.Valid, locale),
				Check:  check.Full,
			},
			{
				Config: testAccNotificationTemplateContentConfig_DefaultVariant_SMS_Minimal(environmentName, licenseID, resourceName, template.Valid, locale),
				Check:  check.Minimal,
			},
			{
				Config: testAccNotificationTemplateContentConfig_DefaultVariant_SMS_Full(environmentName, licenseID, resourceName, template.Valid, locale),
				Check:  check.Full,
			},
			{
				Config:  testAccNotificationTemplateContentConfig_DefaultVariant_SMS_Full(environmentName, licenseID, resourceName, template.Valid, locale),
				Destroy: true,
			},
			// Bad method for template
			{
				Config:      testAccNotificationTemplateContentConfig_DefaultVariant_SMS_Full(environmentName, licenseID, resourceName, template.Invalid, locale),
				ExpectError: regexp.MustCompile(`The configured delivery method does not apply to the selected template.`),
			},
		},
	})
}

func TestAccNotificationTemplateContent_Voice(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_notification_template_content.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	locale := "en"

	template := acctest.TestData{
		Valid:   "strong_authentication",
		Invalid: "email_verification_admin",
	}

	check := acctest.MinMaxChecks{
		Minimal: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "template_name", template.Valid),
			resource.TestCheckResourceAttr(resourceFullName, "locale", locale),
			resource.TestCheckNoResourceAttr(resourceFullName, "email"),
			resource.TestCheckNoResourceAttr(resourceFullName, "push"),
			resource.TestCheckNoResourceAttr(resourceFullName, "sms"),
			resource.TestCheckResourceAttr(resourceFullName, "voice.content", "Min - Hello <pause1sec> your authentication code is <sayCharValue>${otp}</sayCharValue><pause1sec><pause1sec><repeatMessage val=2>I repeat <pause1sec>your code is <sayCharValue>${otp}</sayCharValue></repeatMessage>"),
			resource.TestCheckResourceAttr(resourceFullName, "voice.type", "Alice"),
		),
		Full: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "template_name", template.Valid),
			resource.TestCheckResourceAttr(resourceFullName, "locale", locale),
			resource.TestCheckNoResourceAttr(resourceFullName, "email"),
			resource.TestCheckNoResourceAttr(resourceFullName, "push"),
			resource.TestCheckNoResourceAttr(resourceFullName, "sms"),
			resource.TestCheckResourceAttr(resourceFullName, "voice.content", "Full - Hello <pause1sec> your authentication code is <sayCharValue>${otp}</sayCharValue><pause1sec><pause1sec><repeatMessage val=2>I repeat <pause1sec>your code is <sayCharValue>${otp}</sayCharValue></repeatMessage>"),
			resource.TestCheckResourceAttr(resourceFullName, "voice.type", "Man"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.NotificationTemplateContent_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Minimal from new
			{
				Config: testAccNotificationTemplateContentConfig_DefaultVariant_Voice_Minimal(environmentName, licenseID, resourceName, template.Valid, locale),
				Check:  check.Minimal,
			},
			{
				Config:  testAccNotificationTemplateContentConfig_DefaultVariant_Voice_Minimal(environmentName, licenseID, resourceName, template.Valid, locale),
				Destroy: true,
			},
			// Full from new
			{
				Config: testAccNotificationTemplateContentConfig_DefaultVariant_Voice_Full(environmentName, licenseID, resourceName, template.Valid, locale),
				Check:  check.Full,
			},
			{
				Config:  testAccNotificationTemplateContentConfig_DefaultVariant_Voice_Full(environmentName, licenseID, resourceName, template.Valid, locale),
				Destroy: true,
			},
			// Update
			{
				Config: testAccNotificationTemplateContentConfig_DefaultVariant_Voice_Full(environmentName, licenseID, resourceName, template.Valid, locale),
				Check:  check.Full,
			},
			{
				Config: testAccNotificationTemplateContentConfig_DefaultVariant_Voice_Minimal(environmentName, licenseID, resourceName, template.Valid, locale),
				Check:  check.Minimal,
			},
			{
				Config: testAccNotificationTemplateContentConfig_DefaultVariant_Voice_Full(environmentName, licenseID, resourceName, template.Valid, locale),
				Check:  check.Full,
			},
			{
				Config:  testAccNotificationTemplateContentConfig_DefaultVariant_Voice_Full(environmentName, licenseID, resourceName, template.Valid, locale),
				Destroy: true,
			},
			// Bad method for template
			{
				Config:      testAccNotificationTemplateContentConfig_DefaultVariant_Voice_Full(environmentName, licenseID, resourceName, template.Invalid, locale),
				ExpectError: regexp.MustCompile(`The configured delivery method does not apply to the selected template.`),
			},
		},
	})
}

func TestAccNotificationTemplateContent_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_notification_template_content.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	name := "strong_authentication"
	locale := "en-GB"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.NotificationTemplateContent_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccNotificationTemplateContentConfig_NewLocale_Minimal(environmentName, licenseID, resourceName, name, locale),
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
				ImportStateId: "badformat/badformat/badformat",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`Unexpected Import Identifier`),
			},
		},
	})
}

func testAccNotificationTemplateContentConfig_NewLocale_Minimal(environmentName, licenseID, resourceName, name, locale string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_language" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  locale = "%[5]s"
}

resource "pingone_notification_template_content" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  template_name  = "%[4]s"
  locale         = "%[5]s"

  push = {
    body  = "Min - Please approve this transaction."
    title = "Min - BX Retail Transaction Request"
  }

  depends_on = [
    pingone_language.%[3]s
  ]
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name, locale)
}

func testAccNotificationTemplateContentConfig_NewVariant_Minimal(environmentName, licenseID, resourceName, name, locale, variant string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_notification_template_content" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  template_name  = "%[4]s"
  locale         = "%[5]s"
  variant        = "%[6]s"

  email = {
    body    = <<EOT
Test $${code.value}
EOT
    subject = "Test"
  }
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name, locale, variant)
}

func testAccNotificationTemplateContentConfig_NoVariant_Minimal(environmentName, licenseID, resourceName, name, locale string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_notification_template_content" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  template_name  = "%[4]s"
  locale         = "%[5]s"

  email = {
    body    = <<EOT
Test $${code.value}
EOT
    subject = "Test"
  }
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name, locale)
}

func testAccNotificationTemplateContentConfig_DuplicateLocale(environmentName, licenseID, resourceName, name, locale string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_notification_template_content" "%[3]s-1" {
  environment_id = pingone_environment.%[2]s.id
  template_name  = "%[4]s"
  locale         = "%[5]s"
  variant        = "test-duplicate-locale"

  push = {
    body  = "1 - Please approve this transaction."
    title = "1 - BX Retail Transaction Request"
  }
}

resource "pingone_notification_template_content" "%[3]s-2" {
  environment_id = pingone_environment.%[2]s.id
  template_name  = "%[4]s"
  locale         = "%[5]s"
  variant        = "test-duplicate-locale"

  push = {
    body  = "2 - Please approve this transaction."
    title = "2 - BX Retail Transaction Request"
  }

  depends_on = [
    pingone_notification_template_content.%[3]s-1
  ]
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name, locale)

}

func testAccNotificationTemplateContentConfig_DefaultVariant_Email_Minimal(environmentName, licenseID, resourceName, name, locale string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_notification_template_content" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  template_name  = "%[4]s"
  locale         = "%[5]s"

  email = {
    body    = <<EOT
Min - <p>
	Use the code provided to verify your email address with PingOne. If you think you received this email in error, contact your supervisor.
	</p>
	<p>
	To verify your email address:
	<ol type="1">
	  <li>Sign-on to the self-service portal. For instructions, see <a href="https://docs.pingidentity.com/bundle/pingone/page/snd1631892368614.html">Managing your PingOne user profile</a>.</li>
	  <li>Enter your verification code in the <b>Contact</b> section: $${code}</li>
	  <li>Click <b>Verify</b>.</li>
	</ol>
	</p>
EOT
    subject = "Min - Verify your email address"
  }
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name, locale)
}

func testAccNotificationTemplateContentConfig_DefaultVariant_Email_Full(environmentName, licenseID, resourceName, name, locale string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_notification_template_content" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  template_name  = "%[4]s"
  locale         = "%[5]s"

  email = {
    body    = <<EOT
Full - <p>
	Use the code provided to verify your email address with PingOne. If you think you received this email in error, contact your supervisor.
	</p>
	<p>
	To verify your email address:
	<ol type="1">
	  <li>Sign-on to the self-service portal. For instructions, see <a href="https://docs.pingidentity.com/bundle/pingone/page/snd1631892368614.html">Managing your PingOne user profile</a>.</li>
	  <li>Enter your verification code in the <b>Contact</b> section: $${code}</li>
	  <li>Click <b>Verify</b>.</li>
	</ol>
	</p>
EOT
    subject = "Full - Verify your email address"

    // from = {
    //   name    = "BX Retail"
    //   address = "noreply@bxretail.org"
    // }

    // reply_to = {
    //   name    = "BX Retail Reply"
    //   address = "reply@bxretail.org"
    // }

    character_set = "iso-8859-5"
    content_type  = "text/plain"
  }
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name, locale)
}

func testAccNotificationTemplateContentConfig_DefaultVariant_Push_Minimal(environmentName, licenseID, resourceName, name, locale string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_notification_template_content" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  template_name  = "%[4]s"
  locale         = "%[5]s"

  push = {
    body  = "Min - Please approve this transaction."
    title = "Min - BX Retail Transaction Request"
  }
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name, locale)
}

func testAccNotificationTemplateContentConfig_DefaultVariant_Push_Full(environmentName, licenseID, resourceName, name, locale string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_notification_template_content" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  template_name  = "%[4]s"
  locale         = "%[5]s"

  push = {
    body  = "Full - Please approve this transaction."
    title = "Full - BX Retail Transaction Request"

    category = "WITHOUT_BANNER_BUTTONS"
  }
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name, locale)
}

func testAccNotificationTemplateContentConfig_DefaultVariant_SMS_Minimal(environmentName, licenseID, resourceName, name, locale string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_notification_template_content" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  template_name  = "%[4]s"
  locale         = "%[5]s"

  sms = {
    content = "Min - Please approve this transaction with passcode $${otp}."
  }
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name, locale)
}

func testAccNotificationTemplateContentConfig_DefaultVariant_SMS_Full(environmentName, licenseID, resourceName, name, locale string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_notification_template_content" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  template_name  = "%[4]s"
  locale         = "%[5]s"

  sms = {
    content = "Full - Please approve this transaction with passcode $${otp}."
    sender  = "BX Retail"
  }
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name, locale)
}

func testAccNotificationTemplateContentConfig_DefaultVariant_Voice_Minimal(environmentName, licenseID, resourceName, name, locale string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_notification_template_content" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  template_name  = "%[4]s"
  locale         = "%[5]s"

  voice = {
    content = "Min - Hello <pause1sec> your authentication code is <sayCharValue>$${otp}</sayCharValue><pause1sec><pause1sec><repeatMessage val=2>I repeat <pause1sec>your code is <sayCharValue>$${otp}</sayCharValue></repeatMessage>"
  }
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name, locale)
}

func testAccNotificationTemplateContentConfig_DefaultVariant_Voice_Full(environmentName, licenseID, resourceName, name, locale string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_notification_template_content" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  template_name  = "%[4]s"
  locale         = "%[5]s"

  voice = {
    content = "Full - Hello <pause1sec> your authentication code is <sayCharValue>$${otp}</sayCharValue><pause1sec><pause1sec><repeatMessage val=2>I repeat <pause1sec>your code is <sayCharValue>$${otp}</sayCharValue></repeatMessage>"
    type    = "Man"
  }
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name, locale)
}
