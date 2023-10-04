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

	var ctx = context.Background()
	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		t.Fatalf("Failed to get API client: %v", err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
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
				Check:  base.NotificationTemplateContent_GetIDs(resourceFullName, &environmentID, &templateName, &notificationTemplateContentID),
			},
			// Replan after removal preconfig
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
		resource.TestCheckResourceAttr(resourceFullName, "variant", ""),
		resource.TestCheckResourceAttr(resourceFullName, "email.#", "0"),
		resource.TestCheckResourceAttr(resourceFullName, "push.#", "1"),
		resource.TestCheckResourceAttr(resourceFullName, "sms.#", "0"),
		resource.TestCheckResourceAttr(resourceFullName, "voice.#", "0"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
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
		resource.TestCheckResourceAttr(resourceFullName, "variant", ""),
		resource.TestCheckResourceAttr(resourceFullName, "email.#", "0"),
		resource.TestCheckResourceAttr(resourceFullName, "push.#", "1"),
		resource.TestCheckResourceAttr(resourceFullName, "sms.#", "0"),
		resource.TestCheckResourceAttr(resourceFullName, "voice.#", "0"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
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

	name := "strong_authentication"
	locale := "en"
	variant := "My New Variant"

	check := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "template_name", name),
		resource.TestCheckResourceAttr(resourceFullName, "locale", locale),
		resource.TestCheckResourceAttr(resourceFullName, "default", "false"),
		resource.TestCheckResourceAttr(resourceFullName, "variant", variant),
		resource.TestCheckResourceAttr(resourceFullName, "email.#", "0"),
		resource.TestCheckResourceAttr(resourceFullName, "push.#", "1"),
		resource.TestCheckResourceAttr(resourceFullName, "sms.#", "0"),
		resource.TestCheckResourceAttr(resourceFullName, "voice.#", "0"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
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

	name := "strong_authentication"
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
				Config: testAccNotificationTemplateContentConfig_DefaultVariant_Push_Minimal(environmentName, licenseID, resourceName, name, locale),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "variant", ""),
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
				ExpectError: regexp.MustCompile(`expected template_name to be one of \["email_verification_admin" "email_verification_user" "general" "transaction" "verification_code_template" "recovery_code_template" "device_pairing" "strong_authentication" "email_phone_verification" "id_verification" "credential_issued" "credential_updated" "digital_wallet_pairing" "credential_revoked"\], got strong_authentication_doesnotexist`),
			},
			{
				Config:      testAccNotificationTemplateContentConfig_DefaultVariant_Push_Minimal(environmentName, licenseID, resourceName, name.Valid, locale.Invalid),
				ExpectError: regexp.MustCompile(`expected locale to be one of \["af" "af-ZA" "ar" "ar-AE" "ar-BH" "ar-DZ" "ar-EG" "ar-IQ" "ar-JO" "ar-KW" "ar-LB" "ar-LY" "ar-MA" "ar-OM" "ar-QA" "ar-SA" "ar-SY" "ar-TN" "ar-YE" "az" "az-AZ" "be" "be-BY" "bg" "bg-BG" "bs-BA" "ca" "ca-ES" "cs" "cs-CZ" "cy" "cy-GB" "da" "da-DK" "de" "de-AT" "de-CH" "de-DE" "de-LI" "de-LU" "dv" "dv-MV" "el" "el-GR" "en" "en-AU" "en-BZ" "en-CA" "en-CB" "en-GB" "en-IE" "en-JM" "en-NZ" "en-PH" "en-TT" "en-US" "en-ZA" "en-ZW" "eo" "es" "es-AR" "es-BO" "es-CL" "es-CO" "es-CR" "es-DO" "es-EC" "es-ES" "es-GT" "es-HN" "es-MX" "es-NI" "es-PA" "es-PE" "es-PR" "es-PY" "es-SV" "es-UY" "es-VE" "et" "et-EE" "eu" "eu-ES" "fa" "fa-IR" "fi" "fi-FI" "fo" "fo-FO" "fr" "fr-BE" "fr-CA" "fr-CH" "fr-FR" "fr-LU" "fr-MC" "gl" "gl-ES" "gu" "gu-IN" "he" "he-IL" "hi" "hi-IN" "hr" "hr-BA" "hr-HR" "hu" "hu-HU" "hy" "hy-AM" "id" "id-ID" "is" "is-IS" "it" "it-CH" "it-IT" "ja" "ja-JP" "ka" "ka-GE" "kk" "kk-KZ" "kn" "kn-IN" "ko" "ko-KR" "kok" "kok-IN" "ky" "ky-KG" "lt" "lt-LT" "lv" "lv-LV" "mi" "mi-NZ" "mk" "mk-MK" "mn" "mn-MN" "mr" "mr-IN" "ms" "ms-BN" "ms-MY" "mt" "mt-MT" "nb" "nb-NO" "nl" "nl-BE" "nl-NL" "nn-NO" "ns" "ns-ZA" "pa" "pa-IN" "pl" "pl-PL" "ps" "ps-AR" "pt" "pt-BR" "pt-PT" "qu" "qu-BO" "qu-EC" "qu-PE" "ro" "ro-RO" "ru" "ru-RU" "sa" "sa-IN" "se" "se-FI" "se-FI" "se-FI" "se-NO" "se-SE" "se-SE" "se-SE" "sk" "sk-SK" "sl" "sl-SI" "sq" "sq-AL" "sr-BA" "sr-SP" "sv" "sv-FI" "sv-SE" "sw" "sw-KE" "syr" "syr-SY" "ta" "ta-IN" "te" "te-IN" "th" "th-TH" "tl" "tl-PH" "tn" "tn-ZA" "tr" "tr-TR" "tt" "tt-RU" "ts" "uk" "uk-UA" "ur" "ur-PK" "uz" "uz-UZ" "uz-UZ" "vi" "vi-VN" "xh" "xh-ZA" "zh" "zh-CN" "zh-HK" "zh-MO" "zh-SG" "zh-TW" "zu" "zu-ZA"\], got en-ZZ`),
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
		Invalid: "", // Invalid not tested, no templates without email
		Valid:   "email_verification_admin",
	}

	check := acctest.MinMaxChecks{
		Minimal: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "template_name", template.Valid),
			resource.TestCheckResourceAttr(resourceFullName, "locale", locale),
			resource.TestCheckResourceAttr(resourceFullName, "email.#", "1"),
			resource.TestCheckResourceAttr(resourceFullName, "email.0.subject", "Min - Verify your email address"),
			resource.TestCheckResourceAttr(resourceFullName, "email.0.body", "Min - <p>\n\tUse the code provided to verify your email address with PingOne. If you think you received this email in error, contact your supervisor.\n\t</p>\n\t<p>\n\tTo verify your email address:\n\t<ol type=\"1\">\n\t  <li>Sign-on to the self-service portal. For instructions, see <a href=\"https://docs.pingidentity.com/bundle/pingone/page/snd1631892368614.html\">Managing your PingOne user profile</a>.</li>\n\t  <li>Enter your verification code in the <b>Contact</b> section: ${code}</li>\n\t  <li>Click <b>Verify</b>.</li>\n\t</ol>\n\t</p>\n"),
			resource.TestCheckResourceAttr(resourceFullName, "email.0.from.#", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "email.0.reply_to.#", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "email.0.character_set", "UTF-8"),
			resource.TestCheckResourceAttr(resourceFullName, "email.0.content_type", "text/html"),
			resource.TestCheckResourceAttr(resourceFullName, "push.#", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "sms.#", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "voice.#", "0"),
		),
		Full: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "template_name", template.Valid),
			resource.TestCheckResourceAttr(resourceFullName, "locale", locale),
			resource.TestCheckResourceAttr(resourceFullName, "email.#", "1"),
			resource.TestCheckResourceAttr(resourceFullName, "email.0.subject", "Full - Verify your email address"),
			resource.TestCheckResourceAttr(resourceFullName, "email.0.body", "Full - <p>\n\tUse the code provided to verify your email address with PingOne. If you think you received this email in error, contact your supervisor.\n\t</p>\n\t<p>\n\tTo verify your email address:\n\t<ol type=\"1\">\n\t  <li>Sign-on to the self-service portal. For instructions, see <a href=\"https://docs.pingidentity.com/bundle/pingone/page/snd1631892368614.html\">Managing your PingOne user profile</a>.</li>\n\t  <li>Enter your verification code in the <b>Contact</b> section: ${code}</li>\n\t  <li>Click <b>Verify</b>.</li>\n\t</ol>\n\t</p>\n"),
			// resource.TestCheckResourceAttr(resourceFullName, "email.0.from.#", "1"),
			// resource.TestCheckResourceAttr(resourceFullName, "email.0.from.0.name", "BX Retail"),
			// resource.TestCheckResourceAttr(resourceFullName, "email.0.from.0.address", "noreply@bxretail.org"),
			// resource.TestCheckResourceAttr(resourceFullName, "email.0.reply_to.#", "1"),
			// resource.TestCheckResourceAttr(resourceFullName, "email.0.reply_to.0.name", "BX Retail Reply"),
			// resource.TestCheckResourceAttr(resourceFullName, "email.0.reply_to.0.address", "reply@bxretail.org"),
			resource.TestCheckResourceAttr(resourceFullName, "email.0.character_set", "iso-8859-5"),
			resource.TestCheckResourceAttr(resourceFullName, "email.0.content_type", "text/plain"),
			resource.TestCheckResourceAttr(resourceFullName, "push.#", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "sms.#", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "voice.#", "0"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
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
			resource.TestCheckResourceAttr(resourceFullName, "email.#", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "push.#", "1"),
			resource.TestCheckResourceAttr(resourceFullName, "push.0.body", "Min - Please approve this transaction."),
			resource.TestCheckResourceAttr(resourceFullName, "push.0.title", "Min - BX Retail Transaction Request"),
			resource.TestCheckResourceAttr(resourceFullName, "push.0.category", "BANNER_BUTTONS"),
			resource.TestCheckResourceAttr(resourceFullName, "sms.#", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "voice.#", "0"),
		),
		Full: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "template_name", template.Valid),
			resource.TestCheckResourceAttr(resourceFullName, "locale", locale),
			resource.TestCheckResourceAttr(resourceFullName, "email.#", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "push.#", "1"),
			resource.TestCheckResourceAttr(resourceFullName, "push.0.body", "Full - Please approve this transaction."),
			resource.TestCheckResourceAttr(resourceFullName, "push.0.title", "Full - BX Retail Transaction Request"),
			resource.TestCheckResourceAttr(resourceFullName, "push.0.category", "WITHOUT_BANNER_BUTTONS"),
			resource.TestCheckResourceAttr(resourceFullName, "sms.#", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "voice.#", "0"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
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
			resource.TestCheckResourceAttr(resourceFullName, "email.#", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "push.#", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "sms.#", "1"),
			resource.TestCheckResourceAttr(resourceFullName, "sms.0.content", "Min - Please approve this transaction with passcode ${otp}."),
			resource.TestCheckResourceAttr(resourceFullName, "sms.0.sender", ""),
			resource.TestCheckResourceAttr(resourceFullName, "voice.#", "0"),
		),
		Full: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "template_name", template.Valid),
			resource.TestCheckResourceAttr(resourceFullName, "locale", locale),
			resource.TestCheckResourceAttr(resourceFullName, "email.#", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "push.#", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "sms.#", "1"),
			resource.TestCheckResourceAttr(resourceFullName, "sms.0.content", "Full - Please approve this transaction with passcode ${otp}."),
			resource.TestCheckResourceAttr(resourceFullName, "sms.0.sender", "BX Retail"),
			resource.TestCheckResourceAttr(resourceFullName, "voice.#", "0"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
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
			resource.TestCheckResourceAttr(resourceFullName, "email.#", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "push.#", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "sms.#", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "voice.#", "1"),
			resource.TestCheckResourceAttr(resourceFullName, "voice.0.content", "Min - Hello <pause1sec> your authentication code is <sayCharValue>${otp}</sayCharValue><pause1sec><pause1sec><repeatMessage val=2>I repeat <pause1sec>your code is <sayCharValue>${otp}</sayCharValue></repeatMessage>"),
			resource.TestCheckResourceAttr(resourceFullName, "voice.0.type", "Alice"),
		),
		Full: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "template_name", template.Valid),
			resource.TestCheckResourceAttr(resourceFullName, "locale", locale),
			resource.TestCheckResourceAttr(resourceFullName, "email.#", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "push.#", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "sms.#", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "voice.#", "1"),
			resource.TestCheckResourceAttr(resourceFullName, "voice.0.content", "Full - Hello <pause1sec> your authentication code is <sayCharValue>${otp}</sayCharValue><pause1sec><pause1sec><repeatMessage val=2>I repeat <pause1sec>your code is <sayCharValue>${otp}</sayCharValue></repeatMessage>"),
			resource.TestCheckResourceAttr(resourceFullName, "voice.0.type", "Man"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
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
				ExpectError:  regexp.MustCompile(`Invalid import ID specified \(".*"\).  The ID should be in the format "environment_id/template_name/notification_template_content_id" and must match regex: .*`),
			},
			{
				ResourceName:  resourceFullName,
				ImportStateId: "/",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`Invalid import ID specified \(".*"\).  The ID should be in the format "environment_id/template_name/notification_template_content_id" and must match regex: .*`),
			},
			{
				ResourceName:  resourceFullName,
				ImportStateId: "badformat/badformat/badformat",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`Invalid import ID specified \(".*"\).  The ID should be in the format "environment_id/template_name/notification_template_content_id" and must match regex: .*`),
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

  push {
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

  push {
    body  = "Min - Please approve this transaction."
    title = "Min - BX Retail Transaction Request"
  }
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name, locale, variant)
}

func testAccNotificationTemplateContentConfig_DuplicateLocale(environmentName, licenseID, resourceName, name, locale string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_notification_template_content" "%[3]s-1" {
  environment_id = pingone_environment.%[2]s.id
  template_name  = "%[4]s"
  locale         = "%[5]s"
  variant        = "test-duplicate-locale"

  push {
    body  = "1 - Please approve this transaction."
    title = "1 - BX Retail Transaction Request"
  }
}

resource "pingone_notification_template_content" "%[3]s-2" {
  environment_id = pingone_environment.%[2]s.id
  template_name  = "%[4]s"
  locale         = "%[5]s"
  variant        = "test-duplicate-locale"

  push {
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

  email {
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

  email {
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

    // from {
    //   name    = "BX Retail"
    //   address = "noreply@bxretail.org"
    // }

    // reply_to {
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

  push {
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

  push {
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

  sms {
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

  sms {
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

  voice {
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

  voice {
    content = "Full - Hello <pause1sec> your authentication code is <sayCharValue>$${otp}</sayCharValue><pause1sec><pause1sec><repeatMessage val=2>I repeat <pause1sec>your code is <sayCharValue>$${otp}</sayCharValue></repeatMessage>"
    type    = "Man"
  }
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name, locale)
}
