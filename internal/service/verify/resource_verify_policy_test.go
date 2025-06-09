// Copyright © 2025 Ping Identity Corporation

package verify_test

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
	baselegacysdk "github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/base/legacysdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/verify"
	"github.com/pingidentity/terraform-provider-pingone/internal/client"
	validation "github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccVerifyPolicy_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_verify_policy.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	var verifyPolicyID, environmentID string

	var p1Client *client.Client
	var ctx = context.Background()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)

			p1Client = acctestlegacysdk.PreCheckTestClient(ctx, t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             verify.VerifyPolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccVerifyPolicy_Full(resourceName, name),
				Check:  verify.VerifyPolicy_GetIDs(resourceFullName, &environmentID, &verifyPolicyID),
			},
			{
				PreConfig: func() {
					verify.VerifyPolicy_RemovalDrift_PreConfig(ctx, p1Client.API.VerifyAPIClient, t, environmentID, verifyPolicyID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the environment
			{
				Config: testAccVerifyPolicyConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check:  verify.VerifyPolicy_GetIDs(resourceFullName, &environmentID, &verifyPolicyID),
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

func TestAccVerifyPolicy_NewEnv(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_verify_policy.%s", resourceName)

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
		CheckDestroy:             verify.VerifyPolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVerifyPolicyConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", validation.P1ResourceIDRegexpFullString),
				),
			},
		},
	})
}

func TestAccVerifyPolicy_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_verify_policy.%s", resourceName)

	name := acctest.ResourceNameGen()
	updatedName := acctest.ResourceNameGen()

	fullPolicy := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", validation.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", validation.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "name", name),
		resource.TestCheckResourceAttr(resourceFullName, "description", fmt.Sprintf("Description for %s", name)),
		resource.TestCheckResourceAttr(resourceFullName, "default", "false"),

		resource.TestCheckResourceAttr(resourceFullName, "government_id.verify", "REQUIRED"),
		resource.TestCheckResourceAttr(resourceFullName, "government_id.inspection_type", "AUTOMATIC"),
		resource.TestCheckResourceAttr(resourceFullName, "government_id.fail_expired_id", "true"),
		resource.TestCheckResourceAttr(resourceFullName, "government_id.provider_auto", "VERIFF"),
		resource.TestCheckResourceAttr(resourceFullName, "government_id.provider_manual", "MITEK"),
		resource.TestCheckResourceAttr(resourceFullName, "government_id.retry_attempts", "2"),

		resource.TestCheckResourceAttr(resourceFullName, "facial_comparison.verify", "REQUIRED"),
		resource.TestCheckResourceAttr(resourceFullName, "facial_comparison.threshold", "HIGH"),

		resource.TestCheckResourceAttr(resourceFullName, "liveness.verify", "REQUIRED"),
		resource.TestCheckResourceAttr(resourceFullName, "liveness.threshold", "HIGH"),
		resource.TestCheckResourceAttr(resourceFullName, "liveness.retry_attempts", "1"),

		resource.TestCheckResourceAttr(resourceFullName, "email.verify", "REQUIRED"),
		resource.TestCheckResourceAttr(resourceFullName, "email.create_mfa_device", "true"),
		resource.TestCheckResourceAttr(resourceFullName, "email.otp.attempts.count", "4"),
		resource.TestCheckResourceAttr(resourceFullName, "email.otp.lifetime.duration", "16"),
		resource.TestCheckResourceAttr(resourceFullName, "email.otp.lifetime.time_unit", "MINUTES"),
		resource.TestCheckResourceAttr(resourceFullName, "email.otp.deliveries.count", "5"),
		resource.TestCheckResourceAttr(resourceFullName, "email.otp.deliveries.cooldown.duration", "33"),
		resource.TestCheckResourceAttr(resourceFullName, "email.otp.deliveries.cooldown.time_unit", "SECONDS"),
		resource.TestCheckResourceAttr(resourceFullName, "email.otp.notification.template_name", "email_phone_verification"),
		resource.TestCheckResourceAttr(resourceFullName, "email.otp.notification.variant_name", "variantABC"),

		resource.TestCheckResourceAttr(resourceFullName, "phone.verify", "REQUIRED"),
		resource.TestCheckResourceAttr(resourceFullName, "phone.create_mfa_device", "true"),
		resource.TestCheckResourceAttr(resourceFullName, "phone.otp.attempts.count", "2"),
		resource.TestCheckResourceAttr(resourceFullName, "phone.otp.lifetime.duration", "7"),
		resource.TestCheckResourceAttr(resourceFullName, "phone.otp.lifetime.time_unit", "MINUTES"),
		resource.TestCheckResourceAttr(resourceFullName, "phone.otp.deliveries.count", "3"),
		resource.TestCheckResourceAttr(resourceFullName, "phone.otp.deliveries.cooldown.duration", "16"),
		resource.TestCheckResourceAttr(resourceFullName, "phone.otp.deliveries.cooldown.time_unit", "SECONDS"),
		resource.TestCheckResourceAttr(resourceFullName, "phone.otp.notification.template_name", "email_phone_verification"),
		resource.TestCheckResourceAttr(resourceFullName, "phone.otp.notification.variant_name", "variantABC"),

		resource.TestCheckResourceAttr(resourceFullName, "transaction.timeout.duration", "27"),
		resource.TestCheckResourceAttr(resourceFullName, "transaction.timeout.time_unit", "MINUTES"),
		resource.TestCheckResourceAttr(resourceFullName, "transaction.data_collection.timeout.duration", "12"),
		resource.TestCheckResourceAttr(resourceFullName, "transaction.data_collection.timeout.time_unit", "MINUTES"),
		resource.TestCheckResourceAttr(resourceFullName, "transaction.data_collection_only", "false"),

		resource.TestCheckResourceAttr(resourceFullName, "voice.verify", "REQUIRED"),
		resource.TestCheckResourceAttr(resourceFullName, "voice.enrollment", "true"),
		resource.TestCheckResourceAttr(resourceFullName, "voice.comparison_threshold", "HIGH"),
		resource.TestCheckResourceAttr(resourceFullName, "voice.liveness_threshold", "MEDIUM"),
		resource.TestCheckResourceAttr(resourceFullName, "voice.text_dependent.samples", "4"),
		resource.TestMatchResourceAttr(resourceFullName, "voice.text_dependent.voice_phrase_id", validation.P1ResourceIDRegexp),
		resource.TestCheckResourceAttr(resourceFullName, "voice.reference_data.retain_original_recordings", "true"),
		resource.TestCheckResourceAttr(resourceFullName, "voice.reference_data.update_on_reenrollment", "true"),
		resource.TestCheckResourceAttr(resourceFullName, "voice.reference_data.update_on_verification", "true"),

		resource.TestMatchResourceAttr(resourceFullName, "created_at", validation.RFC3339Regexp),
		resource.TestMatchResourceAttr(resourceFullName, "updated_at", validation.RFC3339Regexp),
	)

	minimalPolicy := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", validation.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", validation.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "name", updatedName),
		resource.TestCheckResourceAttr(resourceFullName, "description", fmt.Sprintf("Description for %s", updatedName)),
		resource.TestCheckResourceAttr(resourceFullName, "default", "false"),

		resource.TestCheckResourceAttr(resourceFullName, "government_id.verify", "REQUIRED"),
		resource.TestCheckResourceAttr(resourceFullName, "government_id.inspection_type", "AUTOMATIC"),
		resource.TestCheckResourceAttr(resourceFullName, "government_id.provider_auto", "MITEK"),
		resource.TestCheckResourceAttr(resourceFullName, "government_id.provider_manual", "MITEK"),

		resource.TestCheckResourceAttr(resourceFullName, "facial_comparison.verify", "DISABLED"),
		resource.TestCheckResourceAttr(resourceFullName, "facial_comparison.threshold", "MEDIUM"),

		resource.TestCheckResourceAttr(resourceFullName, "liveness.verify", "DISABLED"),
		resource.TestCheckResourceAttr(resourceFullName, "liveness.threshold", "MEDIUM"),

		resource.TestCheckResourceAttr(resourceFullName, "email.verify", "DISABLED"),
		resource.TestCheckResourceAttr(resourceFullName, "email.create_mfa_device", "false"),
		resource.TestCheckResourceAttr(resourceFullName, "email.otp.attempts.count", "5"),
		resource.TestCheckResourceAttr(resourceFullName, "email.otp.lifetime.duration", "10"),
		resource.TestCheckResourceAttr(resourceFullName, "email.otp.lifetime.time_unit", "MINUTES"),
		resource.TestCheckResourceAttr(resourceFullName, "email.otp.deliveries.count", "3"),
		resource.TestCheckResourceAttr(resourceFullName, "email.otp.deliveries.cooldown.duration", "30"),
		resource.TestCheckResourceAttr(resourceFullName, "email.otp.deliveries.cooldown.time_unit", "SECONDS"),
		resource.TestCheckResourceAttr(resourceFullName, "email.otp.notification.template_name", "email_phone_verification"),

		resource.TestCheckResourceAttr(resourceFullName, "phone.verify", "DISABLED"),
		resource.TestCheckResourceAttr(resourceFullName, "phone.create_mfa_device", "false"),
		resource.TestCheckResourceAttr(resourceFullName, "phone.otp.attempts.count", "5"),
		resource.TestCheckResourceAttr(resourceFullName, "phone.otp.lifetime.duration", "5"),
		resource.TestCheckResourceAttr(resourceFullName, "phone.otp.lifetime.time_unit", "MINUTES"),
		resource.TestCheckResourceAttr(resourceFullName, "phone.otp.deliveries.count", "3"),
		resource.TestCheckResourceAttr(resourceFullName, "phone.otp.deliveries.cooldown.duration", "30"),
		resource.TestCheckResourceAttr(resourceFullName, "phone.otp.deliveries.cooldown.time_unit", "SECONDS"),
		resource.TestCheckResourceAttr(resourceFullName, "phone.otp.notification.template_name", "email_phone_verification"),

		resource.TestCheckResourceAttr(resourceFullName, "transaction.timeout.duration", "30"),
		resource.TestCheckResourceAttr(resourceFullName, "transaction.timeout.time_unit", "MINUTES"),
		resource.TestCheckResourceAttr(resourceFullName, "transaction.data_collection.timeout.duration", "15"),
		resource.TestCheckResourceAttr(resourceFullName, "transaction.data_collection.timeout.time_unit", "MINUTES"),
		resource.TestCheckResourceAttr(resourceFullName, "transaction.data_collection_only", "false"),

		resource.TestCheckResourceAttr(resourceFullName, "voice.verify", "DISABLED"),
		resource.TestCheckResourceAttr(resourceFullName, "voice.enrollment", "false"),
		resource.TestCheckResourceAttr(resourceFullName, "voice.comparison_threshold", "MEDIUM"),
		resource.TestCheckResourceAttr(resourceFullName, "voice.liveness_threshold", "MEDIUM"),
		resource.TestCheckResourceAttr(resourceFullName, "voice.text_dependent.samples", "3"),
		resource.TestCheckResourceAttr(resourceFullName, "voice.text_dependent.voice_phrase_id", "exceptional_experiences"),
		resource.TestCheckResourceAttr(resourceFullName, "voice.reference_data.retain_original_recordings", "false"),
		resource.TestCheckResourceAttr(resourceFullName, "voice.reference_data.update_on_reenrollment", "true"),
		resource.TestCheckResourceAttr(resourceFullName, "voice.reference_data.update_on_verification", "true"),

		resource.TestMatchResourceAttr(resourceFullName, "created_at", validation.RFC3339Regexp),
		resource.TestMatchResourceAttr(resourceFullName, "updated_at", validation.RFC3339Regexp),
	)

	updateTimeUnitsPolicy := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", validation.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", validation.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "name", updatedName),
		resource.TestCheckResourceAttr(resourceFullName, "description", fmt.Sprintf("Timeunit Policy Update Description for %s", updatedName)),
		resource.TestCheckResourceAttr(resourceFullName, "default", "false"),

		resource.TestCheckResourceAttr(resourceFullName, "government_id.verify", "DISABLED"),
		resource.TestCheckNoResourceAttr(resourceFullName, "government_id.inspection_type"),

		resource.TestCheckResourceAttr(resourceFullName, "facial_comparison.verify", "REQUIRED"),
		resource.TestCheckResourceAttr(resourceFullName, "facial_comparison.threshold", "HIGH"),

		resource.TestCheckResourceAttr(resourceFullName, "liveness.verify", "OPTIONAL"),
		resource.TestCheckResourceAttr(resourceFullName, "liveness.threshold", "LOW"),

		resource.TestCheckResourceAttr(resourceFullName, "email.verify", "REQUIRED"),
		resource.TestCheckResourceAttr(resourceFullName, "email.create_mfa_device", "true"),
		resource.TestCheckResourceAttr(resourceFullName, "email.otp.attempts.count", "4"),
		resource.TestCheckResourceAttr(resourceFullName, "email.otp.lifetime.duration", "90"),
		resource.TestCheckResourceAttr(resourceFullName, "email.otp.lifetime.time_unit", "SECONDS"),
		resource.TestCheckResourceAttr(resourceFullName, "email.otp.deliveries.count", "5"),
		resource.TestCheckResourceAttr(resourceFullName, "email.otp.deliveries.cooldown.duration", "65"),
		resource.TestCheckResourceAttr(resourceFullName, "email.otp.deliveries.cooldown.time_unit", "SECONDS"),
		resource.TestCheckResourceAttr(resourceFullName, "email.otp.notification.template_name", "email_phone_verification"),

		resource.TestCheckResourceAttr(resourceFullName, "phone.verify", "REQUIRED"),
		resource.TestCheckResourceAttr(resourceFullName, "phone.create_mfa_device", "false"),
		resource.TestCheckResourceAttr(resourceFullName, "phone.otp.attempts.count", "2"),
		resource.TestCheckResourceAttr(resourceFullName, "phone.otp.lifetime.duration", "600"),
		resource.TestCheckResourceAttr(resourceFullName, "phone.otp.lifetime.time_unit", "SECONDS"),
		resource.TestCheckResourceAttr(resourceFullName, "phone.otp.deliveries.count", "1"),
		resource.TestCheckResourceAttr(resourceFullName, "phone.otp.deliveries.cooldown.duration", "5"),
		resource.TestCheckResourceAttr(resourceFullName, "phone.otp.deliveries.cooldown.time_unit", "MINUTES"),
		resource.TestCheckResourceAttr(resourceFullName, "phone.otp.notification.template_name", "email_phone_verification"),
		resource.TestCheckResourceAttr(resourceFullName, "phone.otp.notification.variant_name", "variantZYX"),

		resource.TestCheckResourceAttr(resourceFullName, "transaction.timeout.duration", "1500"),
		resource.TestCheckResourceAttr(resourceFullName, "transaction.timeout.time_unit", "SECONDS"),
		resource.TestCheckResourceAttr(resourceFullName, "transaction.data_collection.timeout.duration", "423"),
		resource.TestCheckResourceAttr(resourceFullName, "transaction.data_collection.timeout.time_unit", "SECONDS"),
		resource.TestCheckResourceAttr(resourceFullName, "transaction.data_collection_only", "true"),

		resource.TestCheckResourceAttr(resourceFullName, "voice.verify", "OPTIONAL"),
		resource.TestCheckResourceAttr(resourceFullName, "voice.enrollment", "false"),
		resource.TestCheckResourceAttr(resourceFullName, "voice.comparison_threshold", "LOW"),
		resource.TestCheckResourceAttr(resourceFullName, "voice.liveness_threshold", "LOW"),
		resource.TestCheckResourceAttr(resourceFullName, "voice.text_dependent.samples", "3"),
		resource.TestCheckResourceAttr(resourceFullName, "voice.text_dependent.voice_phrase_id", "exceptional_experiences"),
		resource.TestCheckResourceAttr(resourceFullName, "voice.reference_data.retain_original_recordings", "false"),
		resource.TestCheckResourceAttr(resourceFullName, "voice.reference_data.update_on_reenrollment", "true"),
		resource.TestCheckResourceAttr(resourceFullName, "voice.reference_data.update_on_verification", "true"),

		resource.TestMatchResourceAttr(resourceFullName, "created_at", validation.RFC3339Regexp),
		resource.TestMatchResourceAttr(resourceFullName, "updated_at", validation.RFC3339Regexp),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             verify.VerifyPolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVerifyPolicy_Full(resourceName, name),
				Check:  fullPolicy,
			},
			{
				Config:  testAccVerifyPolicy_Full(resourceName, name),
				Destroy: true,
			},
			{
				Config: testAccVerifyPolicy_Minimal(resourceName, updatedName),
				Check:  minimalPolicy,
			},
			{
				Config:  testAccVerifyPolicy_Minimal(resourceName, updatedName),
				Destroy: true,
			},
			// changes
			{
				Config: testAccVerifyPolicy_Full(resourceName, name),
				Check:  fullPolicy,
			},
			{
				Config: testAccVerifyPolicy_Minimal(resourceName, updatedName),
				Check:  minimalPolicy,
			},
			{
				Config: testAccVerifyPolicy_MinimalDisabledDevice(resourceName, updatedName),
				Check:  minimalPolicy,
			},
			{
				Config: testAccVerifyPolicy_UpdateTimeUnits(resourceName, updatedName),
				Check:  updateTimeUnitsPolicy,
			},
			{
				Config: testAccVerifyPolicy_Full(resourceName, name),
				Check:  fullPolicy,
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

						return fmt.Sprintf("%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccVerifyPolicy_ValidationChecks(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	name := acctest.ResourceNameGen()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             verify.VerifyPolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccVerifyPolicy_NoChecksDefined(resourceName, name),
				ExpectError: regexp.MustCompile(`(?s)(.*Error: Invalid Attribute Combination.*){5}`),
				Destroy:     true,
			},
			{
				Config: testAccVerifyPolicy_EmptyCheckDefinitions(resourceName, name),
				ExpectError: regexp.MustCompile(`(?s)(.*Inappropriate value for attribute \"government_id\".*)` +
					`(.*Inappropriate value for attribute \"facial_comparison\".*)` +
					`(.*Inappropriate value for attribute \"liveness\".*)` +
					`(.*Inappropriate value for attribute \"email\".*)` +
					`(.*Inappropriate value for attribute \"phone\".*)` +
					`(.*Inappropriate value for attribute \"voice\".*)`),
				Destroy: true,
			},
			{
				Config:      testAccVerifyPolicy_IncorrectTransactionDurationRange(resourceName, name),
				ExpectError: regexp.MustCompile("Error: Provided value is not valid"),
				Destroy:     true,
			},
			{
				Config:      testAccVerifyPolicy_TransactionDataCollectionDurationBeyondTimeoutDuration(resourceName, name),
				ExpectError: regexp.MustCompile("Error: Provided value is not valid"),
				Destroy:     true,
			},
			{
				Config:      testAccVerifyPolicy_GovernmentIdInspectionTypeNotAllowed(resourceName, name),
				ExpectError: regexp.MustCompile("Error: Invalid Attribute Value Match"),
				Destroy:     true,
			},
		},
	})
}

func TestAccVerifyPolicy_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_verify_policy.%s", resourceName)

	updatedName := acctest.ResourceNameGen()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             verify.VerifyPolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccVerifyPolicy_Minimal(resourceName, updatedName),
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

func testAccVerifyPolicyConfig_NewEnv(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_verify_policy" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  name        = "%[4]s"
  description = "%[4]s"

  facial_comparison = {
    verify    = "REQUIRED"
    threshold = "LOW"
  }

}`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccVerifyPolicy_Full(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_verify_voice_phrase" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  display_name   = "%[3]s"
}

resource "pingone_verify_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Description for %[3]s"

  government_id = {
    verify          = "REQUIRED"
    inspection_type = "AUTOMATIC"
    fail_expired_id = true
    provider_auto   = "VERIFF"
    provider_manual = "MITEK"
    retry_attempts  = "2"
  }

  facial_comparison = {
    verify    = "REQUIRED"
    threshold = "HIGH"
  }

  liveness = {
    verify         = "REQUIRED"
    threshold      = "HIGH"
    retry_attempts = "1"
  }

  email = {
    verify            = "REQUIRED"
    create_mfa_device = true
    otp = {
      attempts = {
        count = "4"
      }
      lifetime = {
        duration  = "16"
        time_unit = "MINUTES"
      },
      deliveries = {
        count = 5
        cooldown = {
          duration  = "33"
          time_unit = "SECONDS"
        }
      }
      notification = {
        variant_name = "variantABC"
      }
    }
  }

  phone = {
    verify            = "REQUIRED"
    create_mfa_device = true
    otp = {
      attempts = {
        count = "2"
      }
      lifetime = {
        duration  = "7"
        time_unit = "MINUTES"
      },
      deliveries = {
        count = 3
        cooldown = {
          duration  = "16"
          time_unit = "SECONDS"
        }
      }
      notification = {
        variant_name = "variantABC"
      }
    }
  }

  transaction = {
    timeout = {
      duration  = "27"
      time_unit = "MINUTES"
    }

    data_collection = {
      timeout = {
        duration  = "12"
        time_unit = "MINUTES"
      }
    }

    data_collection_only = false
  }

  voice = {
    verify               = "REQUIRED"
    enrollment           = true
    comparison_threshold = "HIGH"
    liveness_threshold   = "MEDIUM"

    text_dependent = {
      samples         = "4"
      voice_phrase_id = pingone_verify_voice_phrase.%[2]s.id
    }

    reference_data = {
      retain_original_recordings = true
      update_on_reenrollment     = true
      update_on_verification     = true
    }
  }


}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccVerifyPolicy_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s
resource "pingone_verify_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Description for %[3]s"

  government_id = {
    verify = "REQUIRED"
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccVerifyPolicy_MinimalDisabledDevice(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s
resource "pingone_verify_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Description for %[3]s"

  government_id = {
    verify = "REQUIRED"
  }

  email = {
    verify = "DISABLED"
  }

  phone = {
    verify = "DISABLED"
  }


}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccVerifyPolicy_UpdateTimeUnits(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s
resource "pingone_verify_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Timeunit Policy Update Description for %[3]s"

  government_id = {
    verify = "DISABLED"
  }

  facial_comparison = {
    verify    = "REQUIRED"
    threshold = "HIGH"
  }

  liveness = {
    verify    = "OPTIONAL"
    threshold = "LOW"
  }

  email = {
    verify            = "REQUIRED"
    create_mfa_device = true
    otp = {
      attempts = {
        count = "4"
      }
      lifetime = {
        duration  = "90"
        time_unit = "SECONDS"
      },
      deliveries = {
        count = 5
        cooldown = {
          duration  = "65"
          time_unit = "SECONDS"
        }
      }
    }
  }

  phone = {
    verify = "REQUIRED"
    otp = {
      attempts = {
        count = "2"
      }
      lifetime = {
        duration  = "600"
        time_unit = "SECONDS"
      },
      deliveries = {
        count = 1
        cooldown = {
          duration  = "5"
          time_unit = "MINUTES"
        }
      }
      notification = {
        variant_name = "variantZYX"
      }
    }
  }

  transaction = {
    timeout = {
      duration  = "1500"
      time_unit = "SECONDS"
    }

    data_collection = {
      timeout = {
        duration  = "423"
        time_unit = "SECONDS"
      }
    }

    data_collection_only = true
  }

  voice = {
    verify               = "OPTIONAL"
    enrollment           = false
    comparison_threshold = "LOW"
    liveness_threshold   = "LOW"

  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccVerifyPolicy_NoChecksDefined(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s
resource "pingone_verify_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "%[3]s"


}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccVerifyPolicy_EmptyCheckDefinitions(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s
resource "pingone_verify_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "%[3]s"

  government_id = {}

  facial_comparison = {}

  liveness = {}

  email = {}

  phone = {}

  voice = {}

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccVerifyPolicy_IncorrectTransactionDurationRange(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s
resource "pingone_verify_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "%[3]s"

  facial_comparison = {
    verify    = "REQUIRED"
    threshold = "HIGH"
  }

  transaction = {
    timeout = {
      duration  = "35"
      time_unit = "MINUTES"
    }

    data_collection = {
      timeout = {
        duration  = "2000"
        time_unit = "SECONDS"
      }
    }
  }


}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccVerifyPolicy_TransactionDataCollectionDurationBeyondTimeoutDuration(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s
resource "pingone_verify_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "%[3]s"

  facial_comparison = {
    verify    = "REQUIRED"
    threshold = "HIGH"
  }

  transaction = {
    timeout = {
      duration  = "15"
      time_unit = "MINUTES"
    }

    data_collection = {
      timeout = {
        duration  = "20"
        time_unit = "MINUTES"
      }
    }
  }


}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccVerifyPolicy_GovernmentIdInspectionTypeNotAllowed(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s
resource "pingone_verify_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "%[3]s"

  government_id = {
    verify    = "DISABLED"
    threshold = "STEP_UP"
  }

  facial_comparison = {
    verify    = "REQUIRED"
    threshold = "HIGH"
  }

  transaction = {
    timeout = {
      duration  = "35"
      time_unit = "MINUTES"
    }

    data_collection = {
      timeout = {
        duration  = "2000"
        time_unit = "SECONDS"
      }
    }
  }


}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
