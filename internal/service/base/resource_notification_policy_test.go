package base_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func testAccCheckNotificationPolicyDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_notification_policy" {
			continue
		}

		_, rEnv, err := apiClient.EnvironmentsApi.ReadOneEnvironment(ctx, rs.Primary.Attributes["environment_id"]).Execute()

		if err != nil {

			if rEnv == nil {
				return fmt.Errorf("Response object does not exist and no error detected")
			}

			if rEnv.StatusCode == 404 {
				continue
			}

			return err
		}

		body, r, err := apiClient.NotificationsPoliciesApi.ReadOneNotificationsPolicy(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.ID).Execute()

		if err != nil {

			if r == nil {
				return fmt.Errorf("Response object does not exist and no error detected")
			}

			if r.StatusCode == 404 {
				continue
			}

			tflog.Error(ctx, fmt.Sprintf("Error: %v", body))
			return err
		}

		return fmt.Errorf("PingOne Notification Policy %s still exists", rs.Primary.ID)
	}

	return nil
}

func testAccGetNotificationPolicyIDs(resourceName string, environmentID, resourceID *string) resource.TestCheckFunc {
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

func TestAccNotificationPolicy_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_notification_policy.%s", resourceName)

	name := resourceName

	var resourceID, environmentID string

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckNotificationPolicyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccNotificationPolicyConfig_Minimal(resourceName, name),
				Check:  testAccGetNotificationPolicyIDs(resourceFullName, &environmentID, &resourceID),
			},
			// Replan after removal preconfig
			{
				PreConfig: func() {
					var ctx = context.Background()
					p1Client, err := acctest.TestClient(ctx)

					if err != nil {
						t.Fatalf("Failed to get API client: %v", err)
					}

					apiClient := p1Client.API.ManagementAPIClient

					if environmentID == "" || resourceID == "" {
						t.Fatalf("One of environment ID or resource ID cannot be determined. Environment ID: %s, Resource ID: %s", environmentID, resourceID)
					}

					_, err = apiClient.NotificationsPoliciesApi.DeleteNotificationsPolicy(ctx, environmentID, resourceID).Execute()
					if err != nil {
						t.Fatalf("Failed to delete Notification Policy: %v", err)
					}
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
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckNotificationPolicyDestroy,
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
			resource.TestCheckResourceAttr(resourceFullName, "quota.#", "1"),
			resource.TestCheckResourceAttr(resourceFullName, "country_limit.type", "DENIED"),
		),
	}

	fullStep2 := resource.TestStep{
		Config: testAccNotificationPolicyConfig_Full(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestCheckResourceAttr(resourceFullName, "default", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "quota.#", "1"),
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
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckNotificationPolicyDestroy,
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
			// Full change
			fullStep1,
			fullStep2,
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
			resource.TestCheckResourceAttr(resourceFullName, "quota.#", "1"),
			resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "quota.*", map[string]string{
				"type": "ENVIRONMENT",
				// "delivery_methods.#": "2",
				// "delivery_methods.0": "SMS",
				// "delivery_methods.1": "Voice",
				"total":  "10000",
				"unused": "",
				"used":   "",
			}),
		),
	}

	quotaUser := resource.TestStep{
		Config: testAccNotificationPolicyConfig_QuotaUser(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "quota.#", "1"),
			resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "quota.*", map[string]string{
				"type": "USER",
				// "delivery_methods.#": "2",
				// "delivery_methods.0": "SMS",
				// "delivery_methods.1": "Voice",
				"total":  "",
				"unused": "45",
				"used":   "40",
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
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckNotificationPolicyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
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
			// Invalid
			{
				Config:      testAccNotificationPolicyConfig_QuotaUser_Invalid(resourceName, name),
				ExpectError: regexp.MustCompile("Invalid parameter"),
			},
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
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckNotificationPolicyDestroy,
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
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckNotificationPolicyDestroy,
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
				ExpectError:  regexp.MustCompile(`Invalid import ID specified \(".*"\).  The ID should be in the format "environment_id/notification_policy_id".`),
			},
			{
				ResourceName:  resourceFullName,
				ImportStateId: "/",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`Invalid import ID specified \(".*"\).  The ID should be in the format "environment_id/notification_policy_id".`),
			},
			{
				ResourceName:  resourceFullName,
				ImportStateId: "badformat/badformat",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`Invalid import ID specified \(".*"\).  The ID should be in the format "environment_id/notification_policy_id".`),
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
    type  = "ENVIRONMENT"
    total = 10000
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
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccNotificationPolicyConfig_QuotaUser(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_notification_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  quota {
    type   = "USER"
    used   = 40
    unused = 45
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
