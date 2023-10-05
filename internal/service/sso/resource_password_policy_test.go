package sso_test

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
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/sso"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccPasswordPolicy_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_password_policy.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	var passwordPolicyID, environmentID string

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
		CheckDestroy:             sso.PasswordPolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccPasswordPolicyConfig_Minimal(resourceName, name),
				Check:  sso.PasswordPolicy_GetIDs(resourceFullName, &environmentID, &passwordPolicyID),
			},
			{
				PreConfig: func() {
					sso.PasswordPolicy_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID, passwordPolicyID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the environment
			{
				SkipFunc: func() (bool, error) { return true, fmt.Errorf("TBC") },
				Config:   testAccPasswordPolicyConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check:    sso.PasswordPolicy_GetIDs(resourceFullName, &environmentID, &passwordPolicyID),
			},
			{
				SkipFunc: func() (bool, error) { return true, fmt.Errorf("TBC") },
				PreConfig: func() {
					base.Environment_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID)
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
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.PasswordPolicy_CheckDestroy,
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
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.PasswordPolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPasswordPolicyConfig_Full(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
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

func TestAccPasswordPolicy_Minimal(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_password_policy.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.PasswordPolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPasswordPolicyConfig_Minimal(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
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

func TestAccPasswordPolicy_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_password_policy.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.PasswordPolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccPasswordPolicyConfig_Minimal(resourceName, name),
			},
			// Errors
			{
				ResourceName: resourceFullName,
				ImportState:  true,
				ExpectError:  regexp.MustCompile(`Invalid import ID specified \(".*"\).  The ID should be in the format "environment_id/password_policy_id" and must match regex: .*`),
			},
			{
				ResourceName:  resourceFullName,
				ImportStateId: "/",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`Invalid import ID specified \(".*"\).  The ID should be in the format "environment_id/password_policy_id" and must match regex: .*`),
			},
			{
				ResourceName:  resourceFullName,
				ImportStateId: "badformat/badformat",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`Invalid import ID specified \(".*"\).  The ID should be in the format "environment_id/password_policy_id" and must match regex: .*`),
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
