// Copyright © 2025 Ping Identity Corporation

package verify_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	acctestlegacysdk "github.com/pingidentity/terraform-provider-pingone/internal/acctest/legacysdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/verify"
	validation "github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccVerifyPolicyDataSource_All(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullName := fmt.Sprintf("data.pingone_verify_policy.%s", resourceName)

	name := acctest.ResourceNameGen()
	updatedName := acctest.ResourceNameGen()

	environmentName := acctest.ResourceNameGenEnvironment()
	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	// P1Verify creates a default policy automatically when the service is enabled.
	// We will use this policy for the default policy lookup tests.
	defaultPolicyName := "Default Verify Policy"
	defaultPolicyDescription := "Default Verify Policy based on Environment Capabilities"

	findByID := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(dataSourceFullName, "id", validation.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", validation.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(dataSourceFullName, "name", name),
		resource.TestCheckResourceAttr(dataSourceFullName, "description", name),
		resource.TestCheckResourceAttr(dataSourceFullName, "default", "false"),

		resource.TestCheckResourceAttr(dataSourceFullName, "government_id.verify", "REQUIRED"),
		resource.TestCheckResourceAttr(dataSourceFullName, "government_id.inspection_type", "AUTOMATIC"),
		resource.TestCheckResourceAttr(dataSourceFullName, "government_id.fail_expired_id", "true"),
		resource.TestCheckResourceAttr(dataSourceFullName, "government_id.provider_auto", "VERIFF"),
		resource.TestCheckResourceAttr(dataSourceFullName, "government_id.provider_manual", "MITEK"),
		resource.TestCheckResourceAttr(dataSourceFullName, "government_id.retry_attempts", "1"),

		resource.TestCheckResourceAttr(dataSourceFullName, "facial_comparison.verify", "REQUIRED"),
		resource.TestCheckResourceAttr(dataSourceFullName, "facial_comparison.threshold", "HIGH"),

		resource.TestCheckResourceAttr(dataSourceFullName, "liveness.verify", "REQUIRED"),
		resource.TestCheckResourceAttr(dataSourceFullName, "liveness.threshold", "HIGH"),
		resource.TestCheckResourceAttr(dataSourceFullName, "liveness.retry_attempts", "3"),

		resource.TestCheckResourceAttr(dataSourceFullName, "email.verify", "REQUIRED"),
		resource.TestCheckResourceAttr(dataSourceFullName, "email.create_mfa_device", "true"),
		resource.TestCheckResourceAttr(dataSourceFullName, "email.otp.attempts.count", "4"),
		resource.TestCheckResourceAttr(dataSourceFullName, "email.otp.lifetime.duration", "16"),
		resource.TestCheckResourceAttr(dataSourceFullName, "email.otp.lifetime.time_unit", "MINUTES"),
		resource.TestCheckResourceAttr(dataSourceFullName, "email.otp.deliveries.count", "5"),
		resource.TestCheckResourceAttr(dataSourceFullName, "email.otp.deliveries.cooldown.duration", "33"),
		resource.TestCheckResourceAttr(dataSourceFullName, "email.otp.deliveries.cooldown.time_unit", "SECONDS"),
		resource.TestCheckResourceAttr(dataSourceFullName, "email.otp.notification.template_name", "email_phone_verification"),
		resource.TestCheckResourceAttr(dataSourceFullName, "email.otp.notification.variant_name", "english_b"),

		resource.TestCheckResourceAttr(dataSourceFullName, "phone.verify", "REQUIRED"),
		resource.TestCheckResourceAttr(dataSourceFullName, "phone.create_mfa_device", "true"),
		resource.TestCheckResourceAttr(dataSourceFullName, "phone.otp.attempts.count", "2"),
		resource.TestCheckResourceAttr(dataSourceFullName, "phone.otp.lifetime.duration", "7"),
		resource.TestCheckResourceAttr(dataSourceFullName, "phone.otp.lifetime.time_unit", "MINUTES"),
		resource.TestCheckResourceAttr(dataSourceFullName, "phone.otp.deliveries.count", "3"),
		resource.TestCheckResourceAttr(dataSourceFullName, "phone.otp.deliveries.cooldown.duration", "16"),
		resource.TestCheckResourceAttr(dataSourceFullName, "phone.otp.deliveries.cooldown.time_unit", "SECONDS"),
		resource.TestCheckResourceAttr(dataSourceFullName, "phone.otp.notification.template_name", "email_phone_verification"),
		resource.TestCheckResourceAttr(dataSourceFullName, "phone.otp.notification.variant_name", "variant23_b"),

		resource.TestCheckResourceAttr(dataSourceFullName, "voice.verify", "DISABLED"),
		resource.TestCheckResourceAttr(dataSourceFullName, "voice.enrollment", "false"),
		resource.TestCheckResourceAttr(dataSourceFullName, "voice.comparison_threshold", "MEDIUM"),
		resource.TestCheckResourceAttr(dataSourceFullName, "voice.liveness_threshold", "MEDIUM"),
		resource.TestCheckResourceAttr(dataSourceFullName, "voice.text_dependent.samples", "3"),
		resource.TestCheckResourceAttr(dataSourceFullName, "voice.text_dependent.voice_phrase_id", "exceptional_experiences"),
		resource.TestCheckResourceAttr(dataSourceFullName, "voice.reference_data.retain_original_recordings", "false"),
		resource.TestCheckResourceAttr(dataSourceFullName, "voice.reference_data.update_on_reenrollment", "true"),
		resource.TestCheckResourceAttr(dataSourceFullName, "voice.reference_data.update_on_verification", "true"),

		resource.TestCheckResourceAttr(dataSourceFullName, "transaction.timeout.duration", "27"),
		resource.TestCheckResourceAttr(dataSourceFullName, "transaction.timeout.time_unit", "MINUTES"),
		resource.TestCheckResourceAttr(dataSourceFullName, "transaction.data_collection.timeout.duration", "12"),
		resource.TestCheckResourceAttr(dataSourceFullName, "transaction.data_collection.timeout.time_unit", "MINUTES"),
		resource.TestCheckResourceAttr(dataSourceFullName, "transaction.data_collection_only", "false"),

		resource.TestMatchResourceAttr(dataSourceFullName, "created_at", validation.RFC3339Regexp),
		resource.TestMatchResourceAttr(dataSourceFullName, "updated_at", validation.RFC3339Regexp),
	)

	findByName := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(dataSourceFullName, "id", validation.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", validation.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(dataSourceFullName, "name", updatedName),
		resource.TestCheckResourceAttr(dataSourceFullName, "description", updatedName),
		resource.TestCheckResourceAttr(dataSourceFullName, "default", "false"),

		resource.TestCheckResourceAttr(dataSourceFullName, "government_id.verify", "DISABLED"),
		resource.TestCheckNoResourceAttr(dataSourceFullName, "government_id.inspection_type"),

		resource.TestCheckResourceAttr(dataSourceFullName, "facial_comparison.verify", "DISABLED"),
		resource.TestCheckResourceAttr(dataSourceFullName, "facial_comparison.threshold", "MEDIUM"),

		resource.TestCheckResourceAttr(dataSourceFullName, "liveness.verify", "REQUIRED"),
		resource.TestCheckResourceAttr(dataSourceFullName, "liveness.threshold", "LOW"),

		resource.TestCheckResourceAttr(dataSourceFullName, "email.verify", "DISABLED"),
		resource.TestCheckResourceAttr(dataSourceFullName, "email.create_mfa_device", "false"),
		resource.TestCheckResourceAttr(dataSourceFullName, "email.otp.attempts.count", "5"),
		resource.TestCheckResourceAttr(dataSourceFullName, "email.otp.lifetime.duration", "10"),
		resource.TestCheckResourceAttr(dataSourceFullName, "email.otp.lifetime.time_unit", "MINUTES"),
		resource.TestCheckResourceAttr(dataSourceFullName, "email.otp.deliveries.count", "3"),
		resource.TestCheckResourceAttr(dataSourceFullName, "email.otp.deliveries.cooldown.duration", "30"),
		resource.TestCheckResourceAttr(dataSourceFullName, "email.otp.deliveries.cooldown.time_unit", "SECONDS"),
		resource.TestCheckResourceAttr(dataSourceFullName, "email.otp.notification.template_name", "email_phone_verification"),

		resource.TestCheckResourceAttr(dataSourceFullName, "phone.verify", "DISABLED"),
		resource.TestCheckResourceAttr(dataSourceFullName, "phone.create_mfa_device", "false"),
		resource.TestCheckResourceAttr(dataSourceFullName, "phone.otp.attempts.count", "5"),
		resource.TestCheckResourceAttr(dataSourceFullName, "phone.otp.lifetime.duration", "5"),
		resource.TestCheckResourceAttr(dataSourceFullName, "phone.otp.lifetime.time_unit", "MINUTES"),
		resource.TestCheckResourceAttr(dataSourceFullName, "phone.otp.deliveries.count", "3"),
		resource.TestCheckResourceAttr(dataSourceFullName, "phone.otp.deliveries.cooldown.duration", "30"),
		resource.TestCheckResourceAttr(dataSourceFullName, "phone.otp.deliveries.cooldown.time_unit", "SECONDS"),
		resource.TestCheckResourceAttr(dataSourceFullName, "phone.otp.notification.template_name", "email_phone_verification"),

		resource.TestCheckResourceAttr(dataSourceFullName, "voice.verify", "REQUIRED"),
		resource.TestCheckResourceAttr(dataSourceFullName, "voice.enrollment", "true"),
		resource.TestCheckResourceAttr(dataSourceFullName, "voice.comparison_threshold", "HIGH"),
		resource.TestCheckResourceAttr(dataSourceFullName, "voice.liveness_threshold", "HIGH"),
		resource.TestCheckResourceAttr(dataSourceFullName, "voice.text_dependent.samples", "5"),
		resource.TestMatchResourceAttr(dataSourceFullName, "voice.text_dependent.voice_phrase_id", validation.P1ResourceIDRegexp),
		resource.TestCheckResourceAttr(dataSourceFullName, "voice.reference_data.retain_original_recordings", "false"),
		resource.TestCheckResourceAttr(dataSourceFullName, "voice.reference_data.update_on_reenrollment", "false"),
		resource.TestCheckResourceAttr(dataSourceFullName, "voice.reference_data.update_on_verification", "false"),

		resource.TestCheckResourceAttr(dataSourceFullName, "transaction.timeout.duration", "30"),
		resource.TestCheckResourceAttr(dataSourceFullName, "transaction.timeout.time_unit", "MINUTES"),
		resource.TestCheckResourceAttr(dataSourceFullName, "transaction.data_collection.timeout.duration", "15"),
		resource.TestCheckResourceAttr(dataSourceFullName, "transaction.data_collection.timeout.time_unit", "MINUTES"),
		resource.TestCheckResourceAttr(dataSourceFullName, "transaction.data_collection_only", "false"),

		resource.TestMatchResourceAttr(dataSourceFullName, "created_at", validation.RFC3339Regexp),
		resource.TestMatchResourceAttr(dataSourceFullName, "updated_at", validation.RFC3339Regexp),
	)

	findDefault := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(dataSourceFullName, "id", validation.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", validation.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(dataSourceFullName, "name", defaultPolicyName),
		resource.TestCheckResourceAttr(dataSourceFullName, "description", defaultPolicyDescription),
		resource.TestCheckResourceAttr(dataSourceFullName, "default", "true"),
		resource.TestMatchResourceAttr(dataSourceFullName, "created_at", validation.RFC3339Regexp),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             verify.VerifyPolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVerifyPolicy_FindByID(environmentName, licenseID, resourceName, name),
				Check:  findByID,
			},
			{
				Config:  testAccVerifyPolicy_FindByID(environmentName, licenseID, resourceName, name),
				Destroy: true,
			},
			{
				Config: testAccVerifyPolicy_FindByName(environmentName, licenseID, resourceName, updatedName, false),
				Check:  findByName,
			},
			{
				Config: testAccVerifyPolicy_FindByName(environmentName, licenseID, resourceName, updatedName, true),
				Check:  findByName,
			},
			{
				Config:  testAccVerifyPolicy_FindByName(environmentName, licenseID, resourceName, updatedName, false),
				Destroy: true,
			},
			{
				Config: testAccVerifyPolicy_FindDefaultPolicy(environmentName, licenseID, resourceName, updatedName),
				Check:  findDefault,
			},
		},
	})
}

func TestAccVerifyPolicyDataSource_FailureChecks(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	name := acctest.ResourceNameGen()

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
		CheckDestroy:             verify.VerifyPolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccVerifyPolicy_FindByIDFail(environmentName, licenseID, resourceName, name),
				ExpectError: regexp.MustCompile("Error: Error when calling `ReadOneVerifyPolicy`: verifyPolicy could not be found"),
			},
			{
				Config:      testAccVerifyPolicy_FindByNameFail(environmentName, licenseID, resourceName, name),
				ExpectError: regexp.MustCompile("Error: Cannot find verify policy from name"),
			},
		},
	})
}

func testAccVerifyPolicy_FindByID(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s
resource "pingone_verify_policy" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s"
  description    = "%[4]s"

  government_id = {
    verify          = "REQUIRED"
    inspection_type = "AUTOMATIC"
    fail_expired_id = true
    provider_auto   = "VERIFF"
    provider_manual = "MITEK"
    retry_attempts  = "1"
  }

  facial_comparison = {
    verify    = "REQUIRED"
    threshold = "HIGH"
  }

  liveness = {
    verify         = "REQUIRED"
    threshold      = "HIGH"
    retry_attempts = "3"
  }

  email = {
    verify = "REQUIRED"
    create_mfa_device : true
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
        variant_name = "english_b"
      }
    }
  }

  phone = {
    verify = "REQUIRED"
    create_mfa_device : true
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
        variant_name = "variant23_b"
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

  depends_on = [pingone_environment.%[2]s]
}

data "pingone_verify_policy" "%[3]s" {
  environment_id   = pingone_environment.%[2]s.id
  verify_policy_id = pingone_verify_policy.%[3]s.id

  depends_on = [pingone_verify_policy.%[3]s]
}`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccVerifyPolicy_FindByName(environmentName, licenseID, resourceName, name string, insensitivityCheck bool) string {

	// If insensitivityCheck is true, alter the case of the name
	nameComparator := name
	if insensitivityCheck {
		nameComparator = acctest.AlterStringCasing(nameComparator)
	}

	return fmt.Sprintf(`
	%[1]s
resource "pingone_verify_voice_phrase" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  display_name   = "%[4]s"
}

resource "pingone_verify_policy" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s"
  description    = "%[4]s"

  liveness = {
    verify    = "REQUIRED"
    threshold = "LOW"
  }

  voice = {
    verify               = "REQUIRED"
    enrollment           = true
    comparison_threshold = "HIGH"
    liveness_threshold   = "HIGH"

    text_dependent = {
      samples         = "5"
      voice_phrase_id = pingone_verify_voice_phrase.%[3]s.id
    }

    reference_data = {
      retain_original_recordings = false
      update_on_reenrollment     = false
      update_on_verification     = false
    }
  }

  depends_on = [pingone_environment.%[2]s]
}

data "pingone_verify_policy" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[5]s"

  depends_on = [pingone_verify_policy.%[3]s]
}`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name, nameComparator)
}

func testAccVerifyPolicy_FindDefaultPolicy(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_verify_policy" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  default        = true

  depends_on = [pingone_environment.%[2]s]
}`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccVerifyPolicy_FindByIDFail(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_verify_policy" "%[3]s" {
  environment_id   = pingone_environment.%[2]s.id
  verify_policy_id = "9c052a8a-14be-44e4-8f07-2662569994ce" // dummy ID that conforms to UUID v4


}`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccVerifyPolicy_FindByNameFail(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_verify_policy" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s"

}`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}
