package sso_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	pingone "github.com/patrickcping/pingone-go/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
)

func testAccCheckPasswordPolicyDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API
	ctx = context.WithValue(ctx, pingone.ContextServerVariables, map[string]string{
		"suffix": p1Client.RegionSuffix,
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
	region := os.Getenv("PINGONE_REGION")

	userFilter := `email ew "@test.com"`
	externalID := "external_1234"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckPasswordPolicyDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPasswordPolicyConfig_Full(environmentName, resourceName, name, description, licenseID, region, userFilter, externalID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, "id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "environment_id"),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", description),
					resource.TestCheckResourceAttrSet(resourceFullName, "population_id"),
					resource.TestCheckResourceAttr(resourceFullName, "user_filter", userFilter),
					resource.TestCheckResourceAttr(resourceFullName, "external_id", externalID),
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
	region := os.Getenv("PINGONE_REGION")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckPasswordPolicyDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPasswordPolicyConfig_Minimal(environmentName, resourceName, name, licenseID, region),
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
					resource.TestCheckResourceAttr(resourceFullName, "max_repeated_characters", ""),
					resource.TestCheckResourceAttr(resourceFullName, "min_complexity", ""),
					resource.TestCheckResourceAttr(resourceFullName, "min_unique_characters", ""),
					resource.TestCheckResourceAttr(resourceFullName, "not_similar_to_current", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "population_count", "0"),
				),
			},
		},
	})
}

func testAccPasswordPolicyConfig_Full(environmentName, resourceName, name, description, licenseID, region, userFilter, externalID string) string {
	return fmt.Sprintf(`
		resource "pingone_environment" "%[1]s" {
			name = "%[1]s"
			type = "SANDBOX"
			license_id = "%[5]s"
			region = "%[6]s"
			default_population {}
			service {}
		}

		resource "pingone_password_policy" "%[2]s" {
			environment_id = "${pingone_environment.%[1]s.id}"
			name = "%[3]s"
			
			description = "My new password policy"

			exclude_commonly_used_passwords = true
			exclude_profile_data = true
			not_similar_to_current = true

			password_history {
				prior_password_count = 6
				retention_days = 365
			}

			password_length {
				min = 8
				max = 255
			}

			password_age {
				max = 182
				min = 1
			}

			account_lockout {
				duration_seconds = 900
				fail_count = 5
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
		}`, environmentName, resourceName, name, description, licenseID, region, userFilter, externalID)
}

func testAccPasswordPolicyConfig_Minimal(environmentName, resourceName, name, licenseID, region string) string {
	return fmt.Sprintf(`
		resource "pingone_environment" "%[1]s" {
			name = "%[1]s"
			type = "SANDBOX"
			license_id = "%[4]s"
			region = "%[5]s"
			default_population {}
			service {}
		}

		resource "pingone_password_policy" "%[2]s" {
			environment_id = "${pingone_environment.%[1]s.id}"
			name = "%[3]s"
		}`, environmentName, resourceName, name, licenseID, region)
}
