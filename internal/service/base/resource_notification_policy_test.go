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
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
)

func testAccCheckNotificationPolicyDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})

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

func TestAccNotificationPolicy_NewEnv(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_notification_policy.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckNotificationPolicyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNotificationPolicy_NewEnv(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
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
			resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestCheckResourceAttr(resourceFullName, "default", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "quota.#", "1"),
		),
	}

	fullStep2 := resource.TestStep{
		Config: testAccNotificationPolicyConfig_Full(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestCheckResourceAttr(resourceFullName, "default", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "quota.#", "1"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckNotificationPolicyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			fullStep1,
			fullStep2,
			fullStep1,
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
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
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
    type      = "USER"
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
    type      = "USER"
    used   = 55
    unused = 45
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
