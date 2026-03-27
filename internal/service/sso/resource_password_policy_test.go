// Copyright © 2026 Ping Identity Corporation

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
	acctestlegacysdk "github.com/pingidentity/terraform-provider-pingone/internal/acctest/legacysdk"
	baselegacysdk "github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/base/legacysdk"
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
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoBeta(t)
			p1Client = acctestlegacysdk.PreCheckTestClient(ctx, t)
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
				Config: testAccPasswordPolicyConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check:  sso.PasswordPolicy_GetIDs(resourceFullName, &environmentID, &passwordPolicyID),
			},
			{
				PreConfig: func() {
					baselegacysdk.Environment_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID)
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
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoBeta(t)
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
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
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
					resource.TestCheckResourceAttr(resourceFullName, "default", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "excludes_commonly_used_passwords", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "excludes_profile_data", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "history.count", "10"),
					resource.TestCheckResourceAttr(resourceFullName, "history.retention_days", "150"),
					resource.TestCheckResourceAttr(resourceFullName, "length.min", "12"),
					resource.TestCheckResourceAttr(resourceFullName, "length.max", "255"),
					resource.TestCheckResourceAttr(resourceFullName, "lockout.duration_seconds", "30"),
					resource.TestCheckResourceAttr(resourceFullName, "lockout.failure_count", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "min_characters.alphabetical_uppercase", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "min_characters.alphabetical_lowercase", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "min_characters.numeric", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "min_characters.special_characters", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "alphabet_sequence_rule.max_length", "3"),
					resource.TestCheckResourceAttr(resourceFullName, "number_sequence_rule.max_length", "3"),
					resource.TestCheckResourceAttr(resourceFullName, "qwerty_sequence_rule.max_length", "3"),
					resource.TestCheckResourceAttr(resourceFullName, "shifted_number_row_sequence_rule.max_length", "3"),
					resource.TestCheckResourceAttr(resourceFullName, "password_age_max", "35"),
					resource.TestCheckResourceAttr(resourceFullName, "password_age_min", "2"),
					resource.TestCheckResourceAttr(resourceFullName, "max_repeated_characters", "2"),
					resource.TestCheckResourceAttr(resourceFullName, "min_complexity", "7"),
					resource.TestCheckResourceAttr(resourceFullName, "min_unique_characters", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "not_similar_to_current", "true"),
					//resource.TestCheckResourceAttr(resourceFullName, "population_count", "1"),
				),
			},
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("resource not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"population_count", // this is ignored because it is 0 (not returned) on recording initial creation state, but is returned on import read, leading to a difference between the state after create and the state after re-import
				},
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
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
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
					resource.TestCheckNoResourceAttr(resourceFullName, "description"),
					resource.TestCheckResourceAttr(resourceFullName, "default", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "excludes_commonly_used_passwords", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "excludes_profile_data", "false"),
					resource.TestCheckNoResourceAttr(resourceFullName, "history"),
					resource.TestCheckNoResourceAttr(resourceFullName, "length"),
					resource.TestCheckNoResourceAttr(resourceFullName, "lockout"),
					resource.TestCheckNoResourceAttr(resourceFullName, "min_characters"),
					resource.TestCheckNoResourceAttr(resourceFullName, "alphabet_sequence_rule"),
					resource.TestCheckNoResourceAttr(resourceFullName, "number_sequence_rule"),
					resource.TestCheckNoResourceAttr(resourceFullName, "qwerty_sequence_rule"),
					resource.TestCheckNoResourceAttr(resourceFullName, "shifted_number_row_sequence_rule"),
					resource.TestCheckNoResourceAttr(resourceFullName, "password_age_max"),
					resource.TestCheckNoResourceAttr(resourceFullName, "password_age_min"),
					resource.TestCheckNoResourceAttr(resourceFullName, "max_repeated_characters"),
					resource.TestCheckNoResourceAttr(resourceFullName, "min_complexity"),
					resource.TestCheckNoResourceAttr(resourceFullName, "min_unique_characters"),
					resource.TestCheckResourceAttr(resourceFullName, "not_similar_to_current", "false"),
					resource.TestCheckNoResourceAttr(resourceFullName, "population_count"),
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
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
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

func TestAccPasswordPolicy_Validation(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.PasswordPolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccPasswordPolicyConfig_WithBody(resourceName, "", ""),
				ExpectError: regexp.MustCompile(`Attribute name string length must be at least 1`),
			},
			{
				Config: testAccPasswordPolicyConfig_WithBody(resourceName, name, `
  alphabet_sequence_rule = {
    max_length = 4
  }
`),
				ExpectError: regexp.MustCompile(`Attribute alphabet_sequence_rule\.max_length value must be one of:`),
			},
			{
				Config: testAccPasswordPolicyConfig_WithBody(resourceName, name, `
  history = {
    count          = 0
    retention_days = 1
  }
`),
				ExpectError: regexp.MustCompile(`Attribute history\.count value must be at least 1`),
			},
			{
				Config: testAccPasswordPolicyConfig_WithBody(resourceName, name, `
  history = {
    count          = 1
    retention_days = 0
  }
`),
				ExpectError: regexp.MustCompile(`Attribute history\.retention_days value must be at least 1`),
			},
			{
				Config: testAccPasswordPolicyConfig_WithBody(resourceName, name, `
  length = {
    min = 8
    max = 254
  }
`),
				ExpectError: regexp.MustCompile(`Attribute length\.max value must be between 255[\s\n]+and[\s\n]+255`),
			},
			{
				Config: testAccPasswordPolicyConfig_WithBody(resourceName, name, `
  length = {
    min = 7
    max = 255
  }
`),
				ExpectError: regexp.MustCompile(`Attribute length\.min value must be between 8[\s\n]+and[\s\n]+32`),
			},
			{
				Config: testAccPasswordPolicyConfig_WithBody(resourceName, name, `
  lockout = {
    duration_seconds = 0
    failure_count    = 1
  }
`),
				ExpectError: regexp.MustCompile(`Attribute lockout\.duration_seconds value must be at least 1`),
			},
			{
				Config: testAccPasswordPolicyConfig_WithBody(resourceName, name, `
  lockout = {
    duration_seconds = 1
    failure_count    = 0
  }
`),
				ExpectError: regexp.MustCompile(`Attribute lockout\.failure_count value must be at least 1`),
			},
			{
				Config: testAccPasswordPolicyConfig_WithBody(resourceName, name, `
  min_characters = {
    alphabetical_uppercase = 2
  }
`),
				ExpectError: regexp.MustCompile(`Attribute min_characters\.alphabetical_uppercase value must be between 0[\s\n]+and[\s\n]+1`),
			},
			{
				Config: testAccPasswordPolicyConfig_WithBody(resourceName, name, `
  min_characters = {
    alphabetical_lowercase = 2
  }
`),
				ExpectError: regexp.MustCompile(`Attribute min_characters\.alphabetical_lowercase value must be between 0[\s\n]+and[\s\n]+1`),
			},
			{
				Config: testAccPasswordPolicyConfig_WithBody(resourceName, name, `
  min_characters = {
    numeric = 2
  }
`),
				ExpectError: regexp.MustCompile(`Attribute min_characters\.numeric value must be between 0[\s\n]+and[\s\n]+1`),
			},
			{
				Config: testAccPasswordPolicyConfig_WithBody(resourceName, name, `
  min_characters = {
    special_characters = 2
  }
`),
				ExpectError: regexp.MustCompile(`Attribute min_characters\.special_characters value must be between 0[\s\n]+and[\s\n]+1`),
			},
			{
				Config: testAccPasswordPolicyConfig_WithBody(resourceName, name, `
  number_sequence_rule = {
    max_length = 4
  }
`),
				ExpectError: regexp.MustCompile(`Attribute number_sequence_rule\.max_length value must be one of:`),
			},
			{
				Config: testAccPasswordPolicyConfig_WithBody(resourceName, name, `
  qwerty_sequence_rule = {
    max_length = 2
  }
`),
				ExpectError: regexp.MustCompile(`Attribute qwerty_sequence_rule\.max_length value must be between 3[\s\n]+and[\s\n]+3`),
			},
			{
				Config: testAccPasswordPolicyConfig_WithBody(resourceName, name, `
  shifted_number_row_sequence_rule = {
    max_length = 2
  }
`),
				ExpectError: regexp.MustCompile(`Attribute shifted_number_row_sequence_rule\.max_length value must be between 3[\s\n]+and[\s\n]+3`),
			},
			{
				Config:      testAccPasswordPolicyConfig_WithBody(resourceName, name, "\n  password_age_max = 0\n"),
				ExpectError: regexp.MustCompile(`Attribute password_age_max value must be at least 1`),
			},
			{
				Config:      testAccPasswordPolicyConfig_WithBody(resourceName, name, "\n  password_age_min = 0\n"),
				ExpectError: regexp.MustCompile(`Attribute password_age_min value must be at least 1`),
			},
			{
				Config:      testAccPasswordPolicyConfig_WithBody(resourceName, name, "\n  max_repeated_characters = 1\n"),
				ExpectError: regexp.MustCompile(`Attribute max_repeated_characters value must be between 2[\s\n]+and[\s\n]+2`),
			},
			{
				Config:      testAccPasswordPolicyConfig_WithBody(resourceName, name, "\n  min_complexity = 6\n"),
				ExpectError: regexp.MustCompile(`Attribute min_complexity value must be between 7[\s\n]+and[\s\n]+7`),
			},
			{
				Config:      testAccPasswordPolicyConfig_WithBody(resourceName, name, "\n  min_unique_characters = 4\n"),
				ExpectError: regexp.MustCompile(`Attribute min_unique_characters value must be between 5[\s\n]+and[\s\n]+5`),
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
}`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccPasswordPolicyConfig_Full(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_password_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  description = "Test description"

  excludes_commonly_used_passwords = true
  excludes_profile_data            = true
  not_similar_to_current           = true

  history = {
    count          = 10
    retention_days = 150
  }

  length = {
    min = 12
    max = 255
  }

  password_age_max = 35
  password_age_min = 2

  lockout = {
    duration_seconds = 30
    failure_count    = 5
  }

  min_characters = {
    alphabetical_uppercase = 0
    alphabetical_lowercase = 1
    numeric                = 1
    special_characters     = 1
  }

  alphabet_sequence_rule = {
    max_length = 3
  }

  number_sequence_rule = {
    max_length = 3
  }

  qwerty_sequence_rule = {
    max_length = 3
  }

  shifted_number_row_sequence_rule = {
    max_length = 3
  }

  max_repeated_characters = 2
  min_complexity          = 7
  min_unique_characters   = 5
}

resource "pingone_population" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name               = "%[3]s"
  password_policy_id = pingone_password_policy.%[2]s.id
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

func testAccPasswordPolicyConfig_WithBody(resourceName, name, body string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_password_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
%[4]s
}`, acctest.GenericSandboxEnvironment(), resourceName, name, body)
}
