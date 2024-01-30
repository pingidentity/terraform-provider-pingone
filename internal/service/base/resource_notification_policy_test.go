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
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)

			p1Client = acctest.PreCheckTestClient(ctx, t)
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
					base.Environment_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID)
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
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
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
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
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
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
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
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
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

func TestAccNotificationPolicy_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_notification_policy.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
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

  quota {
    type  = "ENVIRONMENT"
    total = 10000
  }
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccNotificationPolicyConfig_Full(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_notification_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  quota {
    type             = "ENVIRONMENT"
    delivery_methods = ["SMS", "Voice"]
    total            = 10000
  }

  quota {
    type             = "ENVIRONMENT"
    delivery_methods = ["Email"]
    total            = 10000
  }

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

func testAccNotificationPolicyConfig_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

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

  quota {
    type  = "ENVIRONMENT"
    total = 10000
  }

  quota {
    type             = "ENVIRONMENT"
    delivery_methods = ["Email"]
    total            = 500
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccNotificationPolicyConfig_QuotaUser(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_notification_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  quota {
    type             = "USER"
    delivery_methods = ["SMS"]
    used             = 40
    unused           = 45
  }

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

  quota {
    type   = "USER"
    used   = 55
    unused = 45
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccNotificationPolicyConfig_QuotaUser_InvalidDeliveryMethod(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_notification_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  quota {
    type             = "USER"
    delivery_methods = ["SMS", "Email"]
    total            = 100
  }
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
