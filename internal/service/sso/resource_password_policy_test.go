package sso_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func testAccCheckPasswordPolicyDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_password_policy" {
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

		body, r, err := apiClient.PasswordPoliciesApi.ReadOnePasswordPolicy(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.ID).Execute()

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

		return fmt.Errorf("PingOne Password Policy %s still exists", rs.Primary.ID)
	}

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

func TestAccPasswordPolicy_NewEnv(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_password_policy.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckPasswordPolicyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPasswordPolicyConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
				),
			},
		},
	})
}

func TestAccPasswordPolicy_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_password_policy.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckPasswordPolicyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPasswordPolicyConfig_Full(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", "Test description"),
					resource.TestCheckResourceAttr(resourceFullName, "environment_default", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "bypass_policy", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "exclude_commonly_used_passwords", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "exclude_profile_data", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "password_history.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "password_history.0.prior_password_count", "10"),
					resource.TestCheckResourceAttr(resourceFullName, "password_history.0.retention_days", "150"),
					resource.TestCheckResourceAttr(resourceFullName, "password_length.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "password_length.0.min", "8"),
					resource.TestCheckResourceAttr(resourceFullName, "password_length.0.max", "255"),
					resource.TestCheckResourceAttr(resourceFullName, "account_lockout.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "account_lockout.0.duration_seconds", "30"),
					resource.TestCheckResourceAttr(resourceFullName, "account_lockout.0.fail_count", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "min_characters.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "min_characters.0.alphabetical_uppercase", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "min_characters.0.alphabetical_lowercase", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "min_characters.0.numeric", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "min_characters.0.special_characters", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "password_age.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "password_age.0.max", "35"),
					resource.TestCheckResourceAttr(resourceFullName, "password_age.0.min", "2"),
					resource.TestCheckResourceAttr(resourceFullName, "max_repeated_characters", "2"),
					resource.TestCheckResourceAttr(resourceFullName, "min_complexity", "7"),
					resource.TestCheckResourceAttr(resourceFullName, "min_unique_characters", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "not_similar_to_current", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "population_count", "0"),
				),
			},
		},
	})
}

func TestAccPasswordPolicy_Minimal(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_password_policy.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckPasswordPolicyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPasswordPolicyConfig_Minimal(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", ""),
					resource.TestCheckResourceAttr(resourceFullName, "environment_default", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "bypass_policy", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "exclude_commonly_used_passwords", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "exclude_profile_data", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "password_history.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "password_length.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "account_lockout.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "min_characters.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "password_age.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "max_repeated_characters", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "min_complexity", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "min_unique_characters", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "not_similar_to_current", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "population_count", "0"),
				),
			},
		},
	})
}

func testAccPasswordPolicyConfig_NewEnv(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_password_policy" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s"
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccPasswordPolicyConfig_Full(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_password_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  description = "Test description"

  exclude_commonly_used_passwords = true
  exclude_profile_data            = true
  not_similar_to_current          = true

  password_history {
    prior_password_count = 10
    retention_days       = 150
  }

  password_length {
    min = 8
    max = 255
  }

  password_age {
    max = 35
    min = 2
  }

  account_lockout {
    duration_seconds = 30
    fail_count       = 5
  }

  min_characters {
    alphabetical_uppercase = 1
    alphabetical_lowercase = 1
    numeric                = 1
    special_characters     = 1
  }

  max_repeated_characters = 2
  min_complexity          = 7
  min_unique_characters   = 5
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccPasswordPolicyConfig_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_password_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
