package sso_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
)

func TestAccPasswordPolicyDataSource_ByNameFull(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_password_policy.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	name := acctest.ResourceNameGen()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckPasswordPolicyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPasswordPolicyDataSourceConfig_ByNameFull(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceFullName, "id"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "environment_id"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "password_policy_id"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "name", resourceFullName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "description", resourceFullName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "environment_default", resourceFullName, "environment_default"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "bypass_policy", resourceFullName, "bypass_policy"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "exclude_commonly_used_passwords", resourceFullName, "exclude_commonly_used_passwords"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "exclude_profile_data", resourceFullName, "exclude_profile_data"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "password_history.%", resourceFullName, "password_history.%"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "password_length.%", resourceFullName, "password_length.%"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "account_lockout.%", resourceFullName, "account_lockout.%"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "min_characters.%", resourceFullName, "min_characters.%"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "password_age.%", resourceFullName, "password_age.%"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "max_repeated_characters", resourceFullName, "max_repeated_characters"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "min_complexity", resourceFullName, "min_complexity"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "min_unique_characters", resourceFullName, "min_unique_characters"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "not_similar_to_current", resourceFullName, "not_similar_to_current"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "population_count", resourceFullName, "population_count"),
				),
			},
		},
	})
}

func TestAccPasswordPolicyDataSource_ByIDFull(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_password_policy.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	name := acctest.ResourceNameGen()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckPasswordPolicyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPasswordPolicyDataSourceConfig_ByIDFull(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceFullName, "id"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "environment_id"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "password_policy_id"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "name", resourceFullName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "description", resourceFullName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "environment_default", resourceFullName, "environment_default"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "bypass_policy", resourceFullName, "bypass_policy"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "exclude_commonly_used_passwords", resourceFullName, "exclude_commonly_used_passwords"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "exclude_profile_data", resourceFullName, "exclude_profile_data"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "password_history.%", resourceFullName, "password_history.%"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "password_length.%", resourceFullName, "password_length.%"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "account_lockout.%", resourceFullName, "account_lockout.%"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "min_characters.%", resourceFullName, "min_characters.%"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "password_age.%", resourceFullName, "password_age.%"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "max_repeated_characters", resourceFullName, "max_repeated_characters"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "min_complexity", resourceFullName, "min_complexity"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "min_unique_characters", resourceFullName, "min_unique_characters"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "not_similar_to_current", resourceFullName, "not_similar_to_current"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "population_count", resourceFullName, "population_count"),
				),
			},
		},
	})
}

func TestAccPasswordPolicyDataSource_NotFound(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckPasswordPolicyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccPasswordPolicyDataSourceConfig_NotFoundByName(resourceName),
				ExpectError: regexp.MustCompile("Cannot find password policy doesnotexist"),
			},
			{
				Config:      testAccPasswordPolicyDataSourceConfig_NotFoundByID(resourceName),
				ExpectError: regexp.MustCompile("Error when calling `ReadOnePasswordPolicy`: The request could not be completed. The requested resource was not found."),
			},
		},
	})
}

func testAccPasswordPolicyDataSourceConfig_ByNameFull(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_password_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  description = "My new password policy"

  exclude_commonly_used_passwords = true
  exclude_profile_data            = true
  not_similar_to_current          = true

  password_history {
    prior_password_count = 6
    retention_days       = 365
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
}

data "pingone_password_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  depends_on = [
    pingone_password_policy.%[2]s
  ]
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccPasswordPolicyDataSourceConfig_ByIDFull(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_password_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  description = "My new password policy"

  exclude_commonly_used_passwords = true
  exclude_profile_data            = true
  not_similar_to_current          = true

  password_history {
    prior_password_count = 6
    retention_days       = 365
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
}

data "pingone_password_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  password_policy_id = pingone_password_policy.%[2]s.id
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccPasswordPolicyDataSourceConfig_NotFoundByName(resourceName string) string {
	return fmt.Sprintf(`
		%[1]s

data "pingone_password_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "doesnotexist"
}`, acctest.GenericSandboxEnvironment(), resourceName)
}

func testAccPasswordPolicyDataSourceConfig_NotFoundByID(resourceName string) string {
	return fmt.Sprintf(`
		%[1]s

data "pingone_password_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  password_policy_id = "9c052a8a-14be-44e4-8f07-2662569994ce" // dummy ID that conforms to UUID v4
}`, acctest.GenericSandboxEnvironment(), resourceName)
}
