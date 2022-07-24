package sso_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
)

func testAccCheckPasswordPolicyDestroy(s *terraform.State) error {
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
		if rs.Type != "pingone_password_policy" {
			continue
		}

		_, rEnv, err := apiClient.EnvironmentsApi.ReadOneEnvironment(ctx, rs.Primary.Attributes["environment_id"]).Execute()

		if rEnv.StatusCode == 404 {
			continue
		}

		if err != nil {
			return err
		}

		body, r, err := apiClient.PasswordPoliciesApi.ReadOnePasswordPolicy(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.ID).Execute()

		if r.StatusCode == 404 {
			continue
		}

		if err != nil {
			tflog.Error(ctx, fmt.Sprintf("Error: %v", body))
			return err
		}

		return fmt.Errorf("PingOne Password Policy %s still exists", rs.Primary.ID)
	}

	return nil
}

func TestAccPasswordPolicy_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_password_policy.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName
	description := "Test description"

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	priorPasswordCount := 10
	retentionDays := 150
	ageMax := 35
	ageMin := 2
	lockoutDuration := 30
	lockoutFailCount := 5

	excludeCommonPasswords := true
	excludeProfileData := true
	notSimilarToCurrent := true

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckPasswordPolicyDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPasswordPolicyConfig_Full(environmentName, licenseID, resourceName, description, priorPasswordCount, retentionDays, ageMax, ageMin, lockoutDuration, lockoutFailCount, excludeCommonPasswords, excludeProfileData, notSimilarToCurrent),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, "id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "environment_id"),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", description),
					resource.TestCheckResourceAttr(resourceFullName, "environment_default", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "bypass_policy", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "exclude_commonly_used_passwords", fmt.Sprintf("%t", excludeCommonPasswords)),
					resource.TestCheckResourceAttr(resourceFullName, "exclude_profile_data", fmt.Sprintf("%t", excludeProfileData)),
					resource.TestCheckResourceAttr(resourceFullName, "password_history.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "password_history.0.prior_password_count", fmt.Sprintf("%d", priorPasswordCount)),
					resource.TestCheckResourceAttr(resourceFullName, "password_history.0.retention_days", fmt.Sprintf("%d", retentionDays)),
					resource.TestCheckResourceAttr(resourceFullName, "password_length.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "password_length.0.min", "8"),
					resource.TestCheckResourceAttr(resourceFullName, "password_length.0.max", "255"),
					resource.TestCheckResourceAttr(resourceFullName, "account_lockout.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "account_lockout.0.duration_seconds", fmt.Sprintf("%d", lockoutDuration)),
					resource.TestCheckResourceAttr(resourceFullName, "account_lockout.0.fail_count", fmt.Sprintf("%d", lockoutFailCount)),
					resource.TestCheckResourceAttr(resourceFullName, "min_characters.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "min_characters.0.alphabetical_uppercase", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "min_characters.0.alphabetical_lowercase", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "min_characters.0.numeric", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "min_characters.0.special_characters", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "password_age.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "password_age.0.max", fmt.Sprintf("%d", ageMax)),
					resource.TestCheckResourceAttr(resourceFullName, "password_age.0.min", fmt.Sprintf("%d", ageMin)),
					resource.TestCheckResourceAttr(resourceFullName, "max_repeated_characters", "2"),
					resource.TestCheckResourceAttr(resourceFullName, "min_complexity", "7"),
					resource.TestCheckResourceAttr(resourceFullName, "min_unique_characters", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "not_similar_to_current", fmt.Sprintf("%t", notSimilarToCurrent)),
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

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckPasswordPolicyDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPasswordPolicyConfig_Minimal(environmentName, resourceName, name, licenseID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, "id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "environment_id"),
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

func testAccPasswordPolicyConfig_Full(environmentName, licenseID, resourceName, description string, priorPasswordCount, retentionDays, ageMax, ageMin, lockoutDuration, lockoutFailCount int, excludeCommonPasswords, excludeProfileData, notSimilarToCurrent bool) string {
	return fmt.Sprintf(`
		resource "pingone_environment" "%[1]s" {
			name = "%[1]s"
			type = "SANDBOX"
			license_id = "%[2]s"
			default_population {}
			service {}
		}

		resource "pingone_password_policy" "%[3]s" {
			environment_id = "${pingone_environment.%[1]s.id}"
			name = "%[3]s"
			
			description = "%[4]s"

			exclude_commonly_used_passwords = %[11]t
			exclude_profile_data = %[12]t
			not_similar_to_current = %[13]t

			password_history {
				prior_password_count = %[5]d
				retention_days = %[6]d
			}

			password_length {
				min = 8
				max = 255
			}

			password_age {
				max = %[7]d
				min = %[8]d
			}

			account_lockout {
				duration_seconds = %[9]d
				fail_count = %[10]d
			}

			min_characters {
				alphabetical_uppercase = 1
				alphabetical_lowercase = 1
				numeric = 1
				special_characters = 1
			}

			max_repeated_characters = 2
			min_complexity = 7
			min_unique_characters = 5
		}`, environmentName, licenseID, resourceName, description, priorPasswordCount, retentionDays, ageMax, ageMin, lockoutDuration, lockoutFailCount, excludeCommonPasswords, excludeProfileData, notSimilarToCurrent)
}

func testAccPasswordPolicyConfig_Minimal(environmentName, resourceName, name, licenseID string) string {
	return fmt.Sprintf(`
		resource "pingone_environment" "%[1]s" {
			name = "%[1]s"
			type = "SANDBOX"
			license_id = "%[4]s"
			default_population {}
			service {}
		}

		resource "pingone_password_policy" "%[2]s" {
			environment_id = "${pingone_environment.%[1]s.id}"
			name = "%[3]s"
		}`, environmentName, resourceName, name, licenseID)
}
