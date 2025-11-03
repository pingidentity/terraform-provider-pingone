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
	acctestlegacysdk "github.com/pingidentity/terraform-provider-pingone/internal/acctest/legacysdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/base"
	baselegacysdk "github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/base/legacysdk"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccNotificationPolicy_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_notification_policy.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	var notificationPolicyID, environmentID string

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
		CheckDestroy:             base.NotificationPolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccNotificationPolicyConfig_Minimal(resourceName, name),
				Check:  base.NotificationPolicy_GetIDs(resourceFullName, &environmentID, &notificationPolicyID),
			},
			{
				PreConfig: func() {
					base.NotificationPolicy_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID, notificationPolicyID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the environment
			{
				Config: testAccNotificationPolicy_NewEnv(environmentName, licenseID, resourceName, name),
				Check:  base.NotificationPolicy_GetIDs(resourceFullName, &environmentID, &notificationPolicyID),
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

func TestAccNotificationPolicy_NewEnv(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_notification_policy.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.NotificationPolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNotificationPolicy_NewEnv(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
				),
			},
		},
	})
}

func TestAccNotificationPolicy_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_notification_policy.%s", resourceName)

	name := resourceName

	fullStep1 := resource.TestStep{
		Config: testAccNotificationPolicyConfig_Full(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestCheckResourceAttr(resourceFullName, "default", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "quota.#", "2"),
			resource.TestCheckResourceAttr(resourceFullName, "country_limit.type", "DENIED"),
			// Cooldown Configuration - email (enabled with all fields)
			resource.TestCheckResourceAttr(resourceFullName, "cooldown_configuration.email.enabled", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "cooldown_configuration.email.group_by", "USER"),
			resource.TestCheckResourceAttr(resourceFullName, "cooldown_configuration.email.resend_limit", "5"),
			resource.TestCheckResourceAttr(resourceFullName, "cooldown_configuration.email.periods.#", "3"),
			// Cooldown Configuration - sms (enabled with all fields)
			resource.TestCheckResourceAttr(resourceFullName, "cooldown_configuration.sms.enabled", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "cooldown_configuration.sms.group_by", "PHONE"),
			resource.TestCheckResourceAttr(resourceFullName, "cooldown_configuration.sms.resend_limit", "3"),
			resource.TestCheckResourceAttr(resourceFullName, "cooldown_configuration.sms.periods.#", "3"),
			// Cooldown Configuration - voice (disabled, no other fields)
			resource.TestCheckResourceAttr(resourceFullName, "cooldown_configuration.voice.enabled", "false"),
			// Cooldown Configuration - whats_app (disabled, no other fields)
			resource.TestCheckResourceAttr(resourceFullName, "cooldown_configuration.whats_app.enabled", "false"),
			// Provider Configuration
			resource.TestCheckResourceAttr(resourceFullName, "provider_configuration.conditions.#", "2"),
			resource.TestMatchResourceAttr(resourceFullName, "provider_configuration.conditions.0.fallback_chain.0.id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "provider_configuration.conditions.1.fallback_chain.0.id", verify.P1ResourceIDRegexpFullString),
		),
	}

	minimalStep := resource.TestStep{
		Config: testAccNotificationPolicyConfig_Minimal(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestCheckResourceAttr(resourceFullName, "default", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "quota.#", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "country_limit.type", "NONE"),
			resource.TestCheckResourceAttr(resourceFullName, "cooldown_configuration.email.enabled", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "cooldown_configuration.sms.enabled", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "cooldown_configuration.voice.enabled", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "cooldown_configuration.whats_app.enabled", "false"),
			resource.TestCheckNoResourceAttr(resourceFullName, "provider_configuration"),
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
		CheckDestroy:             base.NotificationPolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			fullStep1,
			{
				Config:  testAccNotificationPolicyConfig_Full(resourceName, name),
				Destroy: true,
			},
			// Minimal
			minimalStep,
			{
				Config:  testAccNotificationPolicyConfig_Minimal(resourceName, name),
				Destroy: true,
			},
			// Change
			fullStep1,
			minimalStep,
			fullStep1,
			// Test importing the resource
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
			},
		},
	})
}

func TestAccNotificationPolicy_Quotas(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_notification_policy.%s", resourceName)

	name := resourceName

	quotaEnvironment := resource.TestStep{
		Config: testAccNotificationPolicyConfig_QuotaEnvironment(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "quota.#", "2"),
			resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "quota.*", map[string]string{
				"type":               "ENVIRONMENT",
				"delivery_methods.#": "2",
				"delivery_methods.0": "SMS",
				"delivery_methods.1": "Voice",
				"total":              "10000",
				"unused":             "",
				"used":               "",
			}),
			resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "quota.*", map[string]string{
				"type":               "ENVIRONMENT",
				"delivery_methods.#": "1",
				"delivery_methods.0": "Email",
				"total":              "500",
				"unused":             "",
				"used":               "",
			}),
		),
	}

	quotaUser := resource.TestStep{
		Config: testAccNotificationPolicyConfig_QuotaUser(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "quota.#", "1"),
			resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "quota.*", map[string]string{
				"type":               "USER",
				"delivery_methods.#": "1",
				"delivery_methods.0": "SMS",
				"total":              "",
				"unused":             "45",
				"used":               "40",
			}),
		),
	}

	quotaUnlimited := resource.TestStep{
		Config: testAccNotificationPolicyConfig_QuotaUnlimited(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "quota.#", "0"),
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
		CheckDestroy:             base.NotificationPolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Invalid
			{
				Config:      testAccNotificationPolicyConfig_QuotaUser_Invalid(resourceName, name),
				ExpectError: regexp.MustCompile("Invalid parameter"),
			},
			{
				Config:      testAccNotificationPolicyConfig_QuotaUser_InvalidDeliveryMethod(resourceName, name),
				ExpectError: regexp.MustCompile("Invalid Attribute Value Match"),
			},
			// Variant 1 New
			quotaEnvironment,
			{
				Config:  testAccNotificationPolicyConfig_QuotaEnvironment(resourceName, name),
				Destroy: true,
			},
			// Variant 2 New
			quotaUser,
			{
				Config:  testAccNotificationPolicyConfig_QuotaUser(resourceName, name),
				Destroy: true,
			},
			// Variant 3 New
			quotaUnlimited,
			{
				Config:  testAccNotificationPolicyConfig_QuotaUnlimited(resourceName, name),
				Destroy: true,
			},
			// Update
			quotaEnvironment,
			quotaUser,
			quotaUnlimited,
			quotaEnvironment,
		},
	})
}

func TestAccNotificationPolicy_CountryLimit(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_notification_policy.%s", resourceName)

	name := resourceName

	countryLimitNone := resource.TestStep{
		Config: testAccNotificationPolicyConfig_CountryLimitNone(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "country_limit.type", "NONE"),
			resource.TestCheckResourceAttr(resourceFullName, "country_limit.delivery_methods.#", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "country_limit.countries.#", "0"),
		),
	}

	countryLimitAllowed := resource.TestStep{
		Config: testAccNotificationPolicyConfig_CountryLimitAllowed(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "country_limit.type", "ALLOWED"),
			resource.TestCheckResourceAttr(resourceFullName, "country_limit.delivery_methods.#", "2"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "country_limit.delivery_methods.*", "Voice"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "country_limit.delivery_methods.*", "SMS"),
			resource.TestCheckResourceAttr(resourceFullName, "country_limit.countries.#", "5"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "country_limit.countries.*", "NP"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "country_limit.countries.*", "HM"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "country_limit.countries.*", "GQ"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "country_limit.countries.*", "GE"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "country_limit.countries.*", "FR"),
		),
	}

	countryLimitDenied := resource.TestStep{
		Config: testAccNotificationPolicyConfig_CountryLimitDenied(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "country_limit.type", "DENIED"),
			resource.TestCheckResourceAttr(resourceFullName, "country_limit.delivery_methods.#", "1"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "country_limit.delivery_methods.*", "Voice"),
			resource.TestCheckResourceAttr(resourceFullName, "country_limit.countries.#", "3"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "country_limit.countries.*", "GB"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "country_limit.countries.*", "NZ"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "country_limit.countries.*", "NO"),
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
		CheckDestroy:             base.NotificationPolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Variant 1 New
			countryLimitNone,
			{
				Config:  testAccNotificationPolicyConfig_CountryLimitNone(resourceName, name),
				Destroy: true,
			},
			// Variant 2 New
			countryLimitAllowed,
			{
				Config:  testAccNotificationPolicyConfig_CountryLimitAllowed(resourceName, name),
				Destroy: true,
			},
			// Variant 3 New
			countryLimitDenied,
			{
				Config:  testAccNotificationPolicyConfig_CountryLimitDenied(resourceName, name),
				Destroy: true,
			},
			// Update
			countryLimitNone,
			countryLimitAllowed,
			countryLimitDenied,
			countryLimitNone,
			{
				Config:  testAccNotificationPolicyConfig_CountryLimitNone(resourceName, name),
				Destroy: true,
			},
			// Invalid
			{
				Config:      testAccNotificationPolicyConfig_CountryLimit_InvalidCombination(resourceName, name),
				ExpectError: regexp.MustCompile(`Invalid argument combination`),
			},
			{
				Config:      testAccNotificationPolicyConfig_CountryLimit_BadCountryCode(resourceName, name),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Match`),
			},
			{
				Config:  testAccNotificationPolicyConfig_CountryLimitAllowed(resourceName, name),
				Destroy: true,
			},
		},
	})
}

func TestAccNotificationPolicy_CooldownConfiguration(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_notification_policy.%s", resourceName)

	name := resourceName

	cooldownDisabled := resource.TestStep{
		Config: testAccNotificationPolicyConfig_CooldownDisabled(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			// Email disabled
			resource.TestCheckResourceAttr(resourceFullName, "cooldown_configuration.email.enabled", "false"),
			resource.TestCheckNoResourceAttr(resourceFullName, "cooldown_configuration.email.group_by"),
			resource.TestCheckNoResourceAttr(resourceFullName, "cooldown_configuration.email.resend_limit"),
			resource.TestCheckResourceAttr(resourceFullName, "cooldown_configuration.email.periods.#", "0"),
			// SMS disabled
			resource.TestCheckResourceAttr(resourceFullName, "cooldown_configuration.sms.enabled", "false"),
			resource.TestCheckNoResourceAttr(resourceFullName, "cooldown_configuration.sms.group_by"),
			resource.TestCheckNoResourceAttr(resourceFullName, "cooldown_configuration.sms.resend_limit"),
			resource.TestCheckResourceAttr(resourceFullName, "cooldown_configuration.sms.periods.#", "0"),
			// Voice disabled
			resource.TestCheckResourceAttr(resourceFullName, "cooldown_configuration.voice.enabled", "false"),
			resource.TestCheckNoResourceAttr(resourceFullName, "cooldown_configuration.voice.group_by"),
			resource.TestCheckNoResourceAttr(resourceFullName, "cooldown_configuration.voice.resend_limit"),
			resource.TestCheckResourceAttr(resourceFullName, "cooldown_configuration.voice.periods.#", "0"),
			// WhatsApp disabled
			resource.TestCheckResourceAttr(resourceFullName, "cooldown_configuration.whats_app.enabled", "false"),
			resource.TestCheckNoResourceAttr(resourceFullName, "cooldown_configuration.whats_app.group_by"),
			resource.TestCheckNoResourceAttr(resourceFullName, "cooldown_configuration.whats_app.resend_limit"),
			resource.TestCheckResourceAttr(resourceFullName, "cooldown_configuration.whats_app.periods.#", "0"),
		),
	}

	cooldownEnabled := resource.TestStep{
		Config: testAccNotificationPolicyConfig_CooldownEnabled(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			// Email cooldown
			resource.TestCheckResourceAttr(resourceFullName, "cooldown_configuration.email.enabled", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "cooldown_configuration.email.group_by", "USER"),
			resource.TestCheckResourceAttr(resourceFullName, "cooldown_configuration.email.resend_limit", "5"),
			resource.TestCheckResourceAttr(resourceFullName, "cooldown_configuration.email.periods.#", "3"),
			resource.TestCheckResourceAttr(resourceFullName, "cooldown_configuration.email.periods.0.duration", "30"),
			resource.TestCheckResourceAttr(resourceFullName, "cooldown_configuration.email.periods.0.time_unit", "SECONDS"),
			resource.TestCheckResourceAttr(resourceFullName, "cooldown_configuration.email.periods.1.duration", "60"),
			resource.TestCheckResourceAttr(resourceFullName, "cooldown_configuration.email.periods.1.time_unit", "SECONDS"),
			resource.TestCheckResourceAttr(resourceFullName, "cooldown_configuration.email.periods.2.duration", "2"),
			resource.TestCheckResourceAttr(resourceFullName, "cooldown_configuration.email.periods.2.time_unit", "MINUTES"),
			// SMS cooldown
			resource.TestCheckResourceAttr(resourceFullName, "cooldown_configuration.sms.enabled", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "cooldown_configuration.sms.group_by", "PHONE"),
			resource.TestCheckResourceAttr(resourceFullName, "cooldown_configuration.sms.resend_limit", "3"),
			resource.TestCheckResourceAttr(resourceFullName, "cooldown_configuration.sms.periods.#", "3"),
			resource.TestCheckResourceAttr(resourceFullName, "cooldown_configuration.sms.periods.0.duration", "45"),
			resource.TestCheckResourceAttr(resourceFullName, "cooldown_configuration.sms.periods.0.time_unit", "SECONDS"),
			resource.TestCheckResourceAttr(resourceFullName, "cooldown_configuration.sms.periods.1.duration", "90"),
			resource.TestCheckResourceAttr(resourceFullName, "cooldown_configuration.sms.periods.1.time_unit", "SECONDS"),
			resource.TestCheckResourceAttr(resourceFullName, "cooldown_configuration.sms.periods.2.duration", "3"),
			resource.TestCheckResourceAttr(resourceFullName, "cooldown_configuration.sms.periods.2.time_unit", "MINUTES"),
			// Voice disabled
			resource.TestCheckResourceAttr(resourceFullName, "cooldown_configuration.voice.enabled", "false"),
			resource.TestCheckNoResourceAttr(resourceFullName, "cooldown_configuration.voice.group_by"),
			resource.TestCheckNoResourceAttr(resourceFullName, "cooldown_configuration.voice.resend_limit"),
			resource.TestCheckResourceAttr(resourceFullName, "cooldown_configuration.voice.periods.#", "0"),
			// WhatsApp disabled
			resource.TestCheckResourceAttr(resourceFullName, "cooldown_configuration.whats_app.enabled", "false"),
			resource.TestCheckNoResourceAttr(resourceFullName, "cooldown_configuration.whats_app.group_by"),
			resource.TestCheckNoResourceAttr(resourceFullName, "cooldown_configuration.whats_app.resend_limit"),
			resource.TestCheckResourceAttr(resourceFullName, "cooldown_configuration.whats_app.periods.#", "0"),
		),
	}

	cooldownMixed := resource.TestStep{
		Config: testAccNotificationPolicyConfig_CooldownMixed(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			// Email enabled with RECIPIENT grouping
			resource.TestCheckResourceAttr(resourceFullName, "cooldown_configuration.email.enabled", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "cooldown_configuration.email.group_by", "RECIPIENT"),
			resource.TestCheckResourceAttr(resourceFullName, "cooldown_configuration.email.resend_limit", "10"),
			resource.TestCheckResourceAttr(resourceFullName, "cooldown_configuration.email.periods.#", "3"),
			// SMS disabled
			resource.TestCheckResourceAttr(resourceFullName, "cooldown_configuration.sms.enabled", "false"),
			resource.TestCheckNoResourceAttr(resourceFullName, "cooldown_configuration.sms.group_by"),
			resource.TestCheckNoResourceAttr(resourceFullName, "cooldown_configuration.sms.resend_limit"),
			resource.TestCheckResourceAttr(resourceFullName, "cooldown_configuration.sms.periods.#", "0"),
			// Voice enabled
			resource.TestCheckResourceAttr(resourceFullName, "cooldown_configuration.voice.enabled", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "cooldown_configuration.voice.group_by", "PHONE"),
			resource.TestCheckResourceAttr(resourceFullName, "cooldown_configuration.voice.resend_limit", "2"),
			resource.TestCheckResourceAttr(resourceFullName, "cooldown_configuration.voice.periods.#", "3"),
			// WhatsApp disabled
			resource.TestCheckResourceAttr(resourceFullName, "cooldown_configuration.whats_app.enabled", "false"),
			resource.TestCheckNoResourceAttr(resourceFullName, "cooldown_configuration.whats_app.group_by"),
			resource.TestCheckNoResourceAttr(resourceFullName, "cooldown_configuration.whats_app.resend_limit"),
			resource.TestCheckResourceAttr(resourceFullName, "cooldown_configuration.whats_app.periods.#", "0"),
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
		CheckDestroy:             base.NotificationPolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Variant 1 New
			cooldownDisabled,
			{
				Config:  testAccNotificationPolicyConfig_CooldownDisabled(resourceName, name),
				Destroy: true,
			},
			// Variant 2 New
			cooldownEnabled,
			{
				Config:  testAccNotificationPolicyConfig_CooldownEnabled(resourceName, name),
				Destroy: true,
			},
			// Variant 3 New
			cooldownMixed,
			{
				Config:  testAccNotificationPolicyConfig_CooldownMixed(resourceName, name),
				Destroy: true,
			},
			// Update
			cooldownDisabled,
			cooldownEnabled,
			cooldownMixed,
			cooldownDisabled,
			{
				Config:  testAccNotificationPolicyConfig_CooldownDisabled(resourceName, name),
				Destroy: true,
			},
			// Invalid - periods must have exactly 3 elements
			{
				Config:      testAccNotificationPolicyConfig_CooldownInvalidPeriodsTooFew(resourceName, name),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value`),
			},
			{
				Config:      testAccNotificationPolicyConfig_CooldownInvalidPeriodsTooMany(resourceName, name),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value`),
			},
			// Invalid - resend_limit required when enabled
			{
				Config:      testAccNotificationPolicyConfig_CooldownInvalidMissingResendLimit(resourceName, name),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			// Invalid - resend_limit out of range
			{
				Config:      testAccNotificationPolicyConfig_CooldownInvalidResendLimitTooLow(resourceName, name),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value`),
			},
			{
				Config:      testAccNotificationPolicyConfig_CooldownInvalidResendLimitTooHigh(resourceName, name),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value`),
			},
			// Invalid - duration out of range for SECONDS
			{
				Config:      testAccNotificationPolicyConfig_CooldownInvalidDurationSecondsTooLow(resourceName, name),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value`),
			},
			{
				Config:      testAccNotificationPolicyConfig_CooldownInvalidDurationSecondsTooHigh(resourceName, name),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value`),
			},
			// Invalid - duration out of range for MINUTES
			{
				Config:      testAccNotificationPolicyConfig_CooldownInvalidDurationMinutesTooLow(resourceName, name),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value`),
			},
			{
				Config:      testAccNotificationPolicyConfig_CooldownInvalidDurationMinutesTooHigh(resourceName, name),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value`),
			},
		},
	})
}

func TestAccNotificationPolicy_ProviderConfiguration(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_notification_policy.%s", resourceName)

	name := resourceName

	providerConfigBasic := resource.TestStep{
		Config: testAccNotificationPolicyConfig_ProviderConfigBasic(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			// Provider Configuration - Two conditions
			resource.TestCheckResourceAttr(resourceFullName, "provider_configuration.conditions.#", "2"),
			// First condition - specific countries
			resource.TestCheckResourceAttr(resourceFullName, "provider_configuration.conditions.0.delivery_methods.#", "2"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "provider_configuration.conditions.0.delivery_methods.*", "SMS"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "provider_configuration.conditions.0.delivery_methods.*", "VOICE"),
			resource.TestCheckResourceAttr(resourceFullName, "provider_configuration.conditions.0.countries.#", "2"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "provider_configuration.conditions.0.countries.*", "US"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "provider_configuration.conditions.0.countries.*", "CA"),
			resource.TestCheckResourceAttr(resourceFullName, "provider_configuration.conditions.0.fallback_chain.#", "1"),
			resource.TestMatchResourceAttr(resourceFullName, "provider_configuration.conditions.0.fallback_chain.0.id", verify.P1ResourceIDRegexpFullString),
			// Second condition - default fallback
			resource.TestCheckResourceAttr(resourceFullName, "provider_configuration.conditions.1.delivery_methods.#", "2"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "provider_configuration.conditions.1.delivery_methods.*", "SMS"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "provider_configuration.conditions.1.delivery_methods.*", "VOICE"),
			resource.TestCheckResourceAttr(resourceFullName, "provider_configuration.conditions.1.countries.#", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "provider_configuration.conditions.1.fallback_chain.#", "1"),
			resource.TestMatchResourceAttr(resourceFullName, "provider_configuration.conditions.1.fallback_chain.0.id", verify.P1ResourceIDRegexpFullString),
		),
	}

	providerConfigMultiple := resource.TestStep{
		Config: testAccNotificationPolicyConfig_ProviderConfigMultiple(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			// Not configured
			resource.TestCheckResourceAttr(resourceFullName, "quota.#", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "country_limit.type", "NONE"),
			resource.TestCheckResourceAttr(resourceFullName, "cooldown_configuration.email.enabled", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "cooldown_configuration.sms.enabled", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "cooldown_configuration.voice.enabled", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "cooldown_configuration.whats_app.enabled", "false"),
			// Provider Configuration - Three conditions
			resource.TestCheckResourceAttr(resourceFullName, "provider_configuration.conditions.#", "3"),
			// First condition - US only with multiple providers
			resource.TestCheckResourceAttr(resourceFullName, "provider_configuration.conditions.0.countries.#", "1"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "provider_configuration.conditions.0.countries.*", "US"),
			resource.TestCheckResourceAttr(resourceFullName, "provider_configuration.conditions.0.fallback_chain.#", "2"),
			// Second condition - EU countries
			resource.TestCheckResourceAttr(resourceFullName, "provider_configuration.conditions.1.countries.#", "3"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "provider_configuration.conditions.1.countries.*", "GB"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "provider_configuration.conditions.1.countries.*", "DE"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "provider_configuration.conditions.1.countries.*", "FR"),
			resource.TestCheckResourceAttr(resourceFullName, "provider_configuration.conditions.1.fallback_chain.#", "1"),
			// Third condition - default
			resource.TestCheckResourceAttr(resourceFullName, "provider_configuration.conditions.2.countries.#", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "provider_configuration.conditions.2.fallback_chain.#", "1"),
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
		CheckDestroy:             base.NotificationPolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Variant 1 New
			providerConfigBasic,
			{
				Config:  testAccNotificationPolicyConfig_ProviderConfigBasic(resourceName, name),
				Destroy: true,
			},
			// Variant 2 New
			providerConfigMultiple,
			{
				Config:  testAccNotificationPolicyConfig_ProviderConfigMultiple(resourceName, name),
				Destroy: true,
			},
			// Update
			providerConfigBasic,
			providerConfigMultiple,
			{
				Config:  testAccNotificationPolicyConfig_ProviderConfigMultiple(resourceName, name),
				Destroy: true,
			},
			// Invalid - conditions must have at least one element
			{
				Config:      testAccNotificationPolicyConfig_ProviderConfigInvalidNoConditions(resourceName, name),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value`),
			},
			// Invalid - fallback_chain must have at least one element
			{
				Config:      testAccNotificationPolicyConfig_ProviderConfigInvalidNoFallbackChain(resourceName, name),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value`),
			},
		},
	})
}

func TestAccNotificationPolicy_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_notification_policy.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.NotificationPolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccNotificationPolicyConfig_Minimal(resourceName, name),
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

func testAccNotificationPolicy_NewEnv(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_notification_policy" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  name = "%[4]s"

  quota = [
    {
      type  = "ENVIRONMENT"
      total = 10000
    }
  ]
}`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccNotificationPolicyConfig_Full(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_phone_delivery_settings" "%[2]s_provider_1" {
  environment_id = data.pingone_environment.general_test.id

  provider_custom = {
    name = "%[2]s Provider 1"
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
    name = "%[2]s Provider 2"
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

  name = "%[3]s"

  quota = [
    {
      type             = "ENVIRONMENT"
      delivery_methods = ["SMS", "Voice"]
      total            = 10000
    },
    {
      type             = "ENVIRONMENT"
      delivery_methods = ["Email"]
      total            = 10000
    }
  ]

  country_limit = {
    type             = "DENIED"
    delivery_methods = ["Voice"]
    countries = [
      "NO",
      "GB",
      "NZ",
    ]
  }

  cooldown_configuration = {
    email = {
      enabled      = true
      group_by     = "USER"
      resend_limit = 5
      periods = [
        {
          duration  = 30
          time_unit = "SECONDS"
        },
        {
          duration  = 60
          time_unit = "SECONDS"
        },
        {
          duration  = 2
          time_unit = "MINUTES"
        }
      ]
    }

    sms = {
      enabled      = true
      group_by     = "PHONE"
      resend_limit = 3
      periods = [
        {
          duration  = 45
          time_unit = "SECONDS"
        },
        {
          duration  = 90
          time_unit = "SECONDS"
        },
        {
          duration  = 3
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

  depends_on = [
    pingone_phone_delivery_settings.%[2]s_provider_1,
    pingone_phone_delivery_settings.%[2]s_provider_2
  ]
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccNotificationPolicyConfig_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_phone_delivery_settings" "%[2]s_provider_1" {
  environment_id = data.pingone_environment.general_test.id

  provider_custom = {
    name = "%[2]s Provider 1"
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
    name = "%[2]s Provider 2"
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

  name = "%[3]s"
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccNotificationPolicyConfig_QuotaEnvironment(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_notification_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  quota = [
    {
      type  = "ENVIRONMENT"
      total = 10000
    },
    {
      type             = "ENVIRONMENT"
      delivery_methods = ["Email"]
      total            = 500
    }
  ]
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccNotificationPolicyConfig_QuotaUser(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_notification_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  quota = [
    {
      type             = "USER"
      delivery_methods = ["SMS"]
      used             = 40
      unused           = 45
    }
  ]

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccNotificationPolicyConfig_QuotaUnlimited(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_notification_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccNotificationPolicyConfig_QuotaUser_Invalid(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_notification_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  quota = [
    {
      type   = "USER"
      used   = 55
      unused = 45
    }
  ]
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccNotificationPolicyConfig_QuotaUser_InvalidDeliveryMethod(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_notification_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  quota = [
    {
      type             = "USER"
      delivery_methods = ["SMS", "Email"]
      total            = 100
    }
  ]
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccNotificationPolicyConfig_CountryLimitNone(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_notification_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  country_limit = {
    type = "NONE"
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccNotificationPolicyConfig_CountryLimitAllowed(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_notification_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  country_limit = {
    type = "ALLOWED"
    countries = [
      "GQ",
      "NP",
      "GE",
      "FR",
      "HM",
    ]
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccNotificationPolicyConfig_CountryLimitDenied(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_notification_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  country_limit = {
    type             = "DENIED"
    delivery_methods = ["Voice"]
    countries = [
      "NO",
      "GB",
      "NZ",
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccNotificationPolicyConfig_CountryLimit_InvalidCombination(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_notification_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  country_limit = {
    type             = "NONE"
    delivery_methods = ["Voice"]
    countries = [
      "NO",
      "GB",
      "NZ",
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccNotificationPolicyConfig_CountryLimit_BadCountryCode(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_notification_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  country_limit = {
    type             = "ALLOWED"
    delivery_methods = ["Voice"]
    countries = [
      "NO",
      "GBE",
      "NZ",
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccNotificationPolicyConfig_CooldownDisabled(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_notification_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  cooldown_configuration = {
    email = {
      enabled = false
    }

    sms = {
      enabled = false
    }

    voice = {
      enabled = false
    }

    whats_app = {
      enabled = false
    }
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccNotificationPolicyConfig_CooldownEnabled(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_notification_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  cooldown_configuration = {
    email = {
      enabled      = true
      group_by     = "USER"
      resend_limit = 5
      periods = [
        {
          duration  = 30
          time_unit = "SECONDS"
        },
        {
          duration  = 60
          time_unit = "SECONDS"
        },
        {
          duration  = 2
          time_unit = "MINUTES"
        }
      ]
    }

    sms = {
      enabled      = true
      group_by     = "PHONE"
      resend_limit = 3
      periods = [
        {
          duration  = 45
          time_unit = "SECONDS"
        },
        {
          duration  = 90
          time_unit = "SECONDS"
        },
        {
          duration  = 3
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
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccNotificationPolicyConfig_CooldownMixed(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_notification_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  cooldown_configuration = {
    email = {
      enabled      = true
      group_by     = "RECIPIENT"
      resend_limit = 10
      periods = [
        {
          duration  = 10
          time_unit = "SECONDS"
        },
        {
          duration  = 20
          time_unit = "SECONDS"
        },
        {
          duration  = 1
          time_unit = "MINUTES"
        }
      ]
    }

    sms = {
      enabled = false
    }

    voice = {
      enabled      = true
      group_by     = "PHONE"
      resend_limit = 2
      periods = [
        {
          duration  = 15
          time_unit = "SECONDS"
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

    whats_app = {
      enabled = false
    }
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccNotificationPolicyConfig_ProviderConfigBasic(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_phone_delivery_settings" "%[2]s_provider_1" {
  environment_id = data.pingone_environment.general_test.id

  provider_custom = {
    name = "%[2]s Provider 1"
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
    name = "%[2]s Provider 2"
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

  name = "%[3]s"

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

  depends_on = [
    pingone_phone_delivery_settings.%[2]s_provider_1,
    pingone_phone_delivery_settings.%[2]s_provider_2
  ]
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccNotificationPolicyConfig_ProviderConfigMultiple(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_phone_delivery_settings" "%[2]s_provider_1" {
  environment_id = data.pingone_environment.general_test.id

  provider_custom = {
    name = "%[2]s Provider 1"
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
    name = "%[2]s Provider 2"
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

resource "pingone_phone_delivery_settings" "%[2]s_provider_3" {
  environment_id = data.pingone_environment.general_test.id

  provider_custom = {
    name = "%[2]s Provider 3"
    authentication = {
      method   = "BASIC"
      username = "test-user-3"
      password = "test-password-3"
    }
    requests = [{
      delivery_method = "SMS"
      url             = "https://example3.com/sms"
      method          = "POST"
    }]
  }
}

resource "pingone_phone_delivery_settings" "%[2]s_provider_4" {
  environment_id = data.pingone_environment.general_test.id

  provider_custom = {
    name = "%[2]s Provider 4"
    authentication = {
      method   = "BASIC"
      username = "test-user-4"
      password = "test-password-4"
    }
    requests = [{
      delivery_method = "VOICE"
      url             = "https://example4.com/voice"
      method          = "POST"
    }]
  }
}

resource "pingone_notification_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  provider_configuration = {
    conditions = [
      {
        delivery_methods = ["SMS"]
        countries        = ["US"]
        fallback_chain = [
          {
            id = pingone_phone_delivery_settings.%[2]s_provider_1.id
          },
          {
            id = pingone_phone_delivery_settings.%[2]s_provider_2.id
          }
        ]
      },
      {
        delivery_methods = ["SMS", "VOICE"]
        countries        = ["GB", "DE", "FR"]
        fallback_chain = [{
          id = pingone_phone_delivery_settings.%[2]s_provider_3.id
        }]
      },
      {
        delivery_methods = ["SMS", "VOICE"]
        fallback_chain = [{
          id = pingone_phone_delivery_settings.%[2]s_provider_4.id
        }]
      }
    ]
  }

  depends_on = [
    pingone_phone_delivery_settings.%[2]s_provider_1,
    pingone_phone_delivery_settings.%[2]s_provider_2,
    pingone_phone_delivery_settings.%[2]s_provider_3,
    pingone_phone_delivery_settings.%[2]s_provider_4
  ]
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccNotificationPolicyConfig_CooldownInvalidPeriodsTooFew(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_notification_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  cooldown_configuration = {
    email = {
      enabled      = true
      group_by     = "USER"
      resend_limit = 5
      periods = [
        {
          duration  = 30
          time_unit = "SECONDS"
        },
        {
          duration  = 60
          time_unit = "SECONDS"
        }
      ]
    }

    sms = {
      enabled = false
    }

    voice = {
      enabled = false
    }

    whats_app = {
      enabled = false
    }
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccNotificationPolicyConfig_CooldownInvalidPeriodsTooMany(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_notification_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  cooldown_configuration = {
    email = {
      enabled      = true
      group_by     = "USER"
      resend_limit = 5
      periods = [
        {
          duration  = 30
          time_unit = "SECONDS"
        },
        {
          duration  = 60
          time_unit = "SECONDS"
        },
        {
          duration  = 2
          time_unit = "MINUTES"
        },
        {
          duration  = 3
          time_unit = "MINUTES"
        }
      ]
    }

    sms = {
      enabled = false
    }

    voice = {
      enabled = false
    }

    whats_app = {
      enabled = false
    }
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccNotificationPolicyConfig_CooldownInvalidMissingResendLimit(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_notification_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  cooldown_configuration = {
    email = {
      enabled  = true
      group_by = "USER"
      periods = [
        {
          duration  = 30
          time_unit = "SECONDS"
        },
        {
          duration  = 60
          time_unit = "SECONDS"
        },
        {
          duration  = 2
          time_unit = "MINUTES"
        }
      ]
    }

    sms = {
      enabled = false
    }

    voice = {
      enabled = false
    }

    whats_app = {
      enabled = false
    }
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccNotificationPolicyConfig_CooldownInvalidResendLimitTooLow(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_notification_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  cooldown_configuration = {
    email = {
      enabled      = true
      group_by     = "USER"
      resend_limit = 0
      periods = [
        {
          duration  = 30
          time_unit = "SECONDS"
        },
        {
          duration  = 60
          time_unit = "SECONDS"
        },
        {
          duration  = 2
          time_unit = "MINUTES"
        }
      ]
    }

    sms = {
      enabled = false
    }

    voice = {
      enabled = false
    }

    whats_app = {
      enabled = false
    }
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccNotificationPolicyConfig_CooldownInvalidResendLimitTooHigh(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_notification_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  cooldown_configuration = {
    email = {
      enabled      = true
      group_by     = "USER"
      resend_limit = 11
      periods = [
        {
          duration  = 30
          time_unit = "SECONDS"
        },
        {
          duration  = 60
          time_unit = "SECONDS"
        },
        {
          duration  = 2
          time_unit = "MINUTES"
        }
      ]
    }

    sms = {
      enabled = false
    }

    voice = {
      enabled = false
    }

    whats_app = {
      enabled = false
    }
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccNotificationPolicyConfig_CooldownInvalidDurationSecondsTooLow(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_notification_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  cooldown_configuration = {
    email = {
      enabled      = true
      group_by     = "USER"
      resend_limit = 5
      periods = [
        {
          duration  = 9
          time_unit = "SECONDS"
        },
        {
          duration  = 60
          time_unit = "SECONDS"
        },
        {
          duration  = 2
          time_unit = "MINUTES"
        }
      ]
    }

    sms = {
      enabled = false
    }

    voice = {
      enabled = false
    }

    whats_app = {
      enabled = false
    }
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccNotificationPolicyConfig_CooldownInvalidDurationSecondsTooHigh(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_notification_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  cooldown_configuration = {
    email = {
      enabled      = true
      group_by     = "USER"
      resend_limit = 5
      periods = [
        {
          duration  = 601
          time_unit = "SECONDS"
        },
        {
          duration  = 60
          time_unit = "SECONDS"
        },
        {
          duration  = 2
          time_unit = "MINUTES"
        }
      ]
    }

    sms = {
      enabled = false
    }

    voice = {
      enabled = false
    }

    whats_app = {
      enabled = false
    }
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccNotificationPolicyConfig_CooldownInvalidDurationMinutesTooLow(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_notification_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  cooldown_configuration = {
    email = {
      enabled      = true
      group_by     = "USER"
      resend_limit = 5
      periods = [
        {
          duration  = 30
          time_unit = "SECONDS"
        },
        {
          duration  = 60
          time_unit = "SECONDS"
        },
        {
          duration  = 0
          time_unit = "MINUTES"
        }
      ]
    }

    sms = {
      enabled = false
    }

    voice = {
      enabled = false
    }

    whats_app = {
      enabled = false
    }
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccNotificationPolicyConfig_CooldownInvalidDurationMinutesTooHigh(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_notification_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  cooldown_configuration = {
    email = {
      enabled      = true
      group_by     = "USER"
      resend_limit = 5
      periods = [
        {
          duration  = 30
          time_unit = "SECONDS"
        },
        {
          duration  = 60
          time_unit = "SECONDS"
        },
        {
          duration  = 11
          time_unit = "MINUTES"
        }
      ]
    }

    sms = {
      enabled = false
    }

    voice = {
      enabled = false
    }

    whats_app = {
      enabled = false
    }
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccNotificationPolicyConfig_ProviderConfigInvalidNoConditions(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_notification_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  provider_configuration = {
    conditions = []
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccNotificationPolicyConfig_ProviderConfigInvalidNoFallbackChain(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_notification_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  provider_configuration = {
    conditions = [
      {
        delivery_methods = ["SMS", "VOICE"]
        fallback_chain   = []
      }
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
