// Copyright © 2025 Ping Identity Corporation

package sso_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/sso"
)

func TestAccPasswordPolicyDataSource_ByNameFull(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_password_policy.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	name := acctest.ResourceNameGen()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.PasswordPolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPasswordPolicyDataSourceConfig_ByNameFull(resourceName, name, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceFullName, "id"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "environment_id"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "password_policy_id"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "name", resourceFullName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "description", resourceFullName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "default", resourceFullName, "default"),
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
					//resource.TestCheckResourceAttrPair(dataSourceFullName, "population_count", resourceFullName, "population_count"),
				),
			},
			{
				Config: testAccPasswordPolicyDataSourceConfig_ByNameFull(resourceName, name, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceFullName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "name", resourceFullName, "name"),
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
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.PasswordPolicy_CheckDestroy,
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
					resource.TestCheckResourceAttrPair(dataSourceFullName, "default", resourceFullName, "default"),
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
					//resource.TestCheckResourceAttrPair(dataSourceFullName, "population_count", resourceFullName, "population_count"),
				),
			},
		},
	})
}

func TestAccPasswordPolicyDataSource_NotFound(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.PasswordPolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccPasswordPolicyDataSourceConfig_NotFoundByName(resourceName),
				ExpectError: regexp.MustCompile("Cannot find the password policy from name"),
			},
			{
				Config:      testAccPasswordPolicyDataSourceConfig_NotFoundByID(resourceName),
				ExpectError: regexp.MustCompile("Error when calling `ReadOnePasswordPolicy`: The request could not be completed. The requested resource was not found."),
			},
		},
	})
}

func testAccPasswordPolicyDataSourceConfig_ByNameFull(resourceName, name string, insensitivityCheck bool) string {

	// If insensitivityCheck is true, alter the case of the name
	nameComparator := name
	if insensitivityCheck {
		nameComparator = acctest.AlterStringCasing(nameComparator)
	}

	return fmt.Sprintf(`
		%[1]s

resource "pingone_password_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  description = "My new password policy"

  excludes_commonly_used_passwords = true
  excludes_profile_data            = true
  not_similar_to_current           = true

  history = {
    count          = 6
    retention_days = 365
  }

  length = {
    min = 8
    max = 255
  }

  password_age_max = 182
  password_age_min = 1

  lockout = {
    duration_seconds = 900
    failure_count    = 5
  }

  min_characters = {
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

  name = "%[4]s"

  depends_on = [
    pingone_password_policy.%[2]s
  ]
}`, acctest.GenericSandboxEnvironment(), resourceName, name, nameComparator)
}

func testAccPasswordPolicyDataSourceConfig_ByIDFull(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_password_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  description = "My new password policy"

  excludes_commonly_used_passwords = true
  excludes_profile_data            = true
  not_similar_to_current           = true

  history = {
    count          = 6
    retention_days = 365
  }

  length = {
    min = 8
    max = 255
  }

  password_age_max = 182
  password_age_min = 1

  lockout = {
    duration_seconds = 900
    failure_count    = 5
  }

  min_characters = {
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
