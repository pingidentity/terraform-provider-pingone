// Copyright © 2025 Ping Identity Corporation

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

func TestAccPopulationDefault_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_population_default.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	var populationID, environmentID string

	var p1Client *client.Client
	var ctx = context.Background()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)

			acctest.PreCheckNoBeta(t)
			p1Client = acctest.PreCheckTestClient(ctx, t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.PopulationDefault_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Test removal of the environment
			{
				Config: testAccPopulationDefaultConfig_Full(environmentName, licenseID, resourceName, name),
				Check:  sso.Population_GetIDs(resourceFullName, &environmentID, &populationID),
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

func TestAccPopulationDefault_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_population_default.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.PopulationDefault_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPopulationDefaultConfig_Full(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", "Test description"),
					resource.TestMatchResourceAttr(resourceFullName, "password_policy.id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "alternative_identifiers.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "alternative_identifiers.*", "identifier1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "alternative_identifiers.*", "identifier2"),
					resource.TestCheckResourceAttr(resourceFullName, "preferred_language", "es"),
					resource.TestMatchResourceAttr(resourceFullName, "theme.id", verify.P1ResourceIDRegexpFullString),
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

						return rs.Primary.Attributes["environment_id"], nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccPopulationDefault_Minimal(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_population_default.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.PopulationDefault_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPopulationDefaultConfig_Minimal(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckNoResourceAttr(resourceFullName, "description"),
					resource.TestCheckNoResourceAttr(resourceFullName, "password_policy_id"),
					resource.TestCheckNoResourceAttr(resourceFullName, "password_policy"),
					resource.TestCheckResourceAttr(resourceFullName, "preferred_language", "en"),
					resource.TestMatchResourceAttr(resourceFullName, "theme.id", verify.P1ResourceIDRegexpFullString),
				),
			},
		},
	})
}

func TestAccPopulationDefault_PasswordPolicy(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Population_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPopulationDefaultConfig_PasswordPolicyNested(environmentName, licenseID, resourceName, name),
			},
			{
				Config: testAccPopulationDefaultConfig_PasswordPolicyString(environmentName, licenseID, resourceName, name),
			},
			{
				Config:      testAccPopulationDefaultConfig_PasswordPolicyConflict(environmentName, licenseID, resourceName, name),
				ExpectError: regexp.MustCompile(`Invalid Attribute Combination`),
			},
		},
	})
}

func TestAccPopulationDefault_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_population_default.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.PopulationDefault_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccPopulationDefaultConfig_Minimal(environmentName, licenseID, resourceName, name),
			},
			// Errors
			{
				ResourceName:  resourceFullName,
				ImportStateId: "/",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`Unexpected Import Identifier`),
			},
			{
				ResourceName:  resourceFullName,
				ImportStateId: "badformat",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`Unexpected Import Identifier`),
			},
		},
	})
}

func testAccPopulationDefaultConfig_Full(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_password_policy" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s"
}

data "pingone_language" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  locale = "es"
}

resource "pingone_language_update" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  language_id = data.pingone_language.%[3]s.id
  enabled     = true
}

resource "pingone_branding_theme" "%[3]s" {
  environment_id     = pingone_environment.%[2]s.id
  name               = "%[3]s"
  template           = "split"
  background_color   = "#FF00F0"
  button_text_color  = "#FF6C6C"
  heading_text_color = "#FF0005"
  card_color         = "#0FFF39"
  body_text_color    = "#8620FF"
  link_text_color    = "#8A7F06"
  button_color       = "#0CFFFB"
}

resource "pingone_population_default" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s"
  description    = "Test description"
  password_policy = {
    id = pingone_password_policy.%[3]s.id
  }
  alternative_identifiers = ["identifier1", "identifier2"]
  preferred_language      = pingone_language_update.%[3]s.locale
  theme = {
    id = pingone_branding_theme.%[3]s.id
  }
}`, acctest.MinimalSandboxEnvironmentNoPopulation(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccPopulationDefaultConfig_Minimal(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_population_default" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s"
}`, acctest.MinimalSandboxEnvironmentNoPopulation(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccPopulationDefaultConfig_PasswordPolicyNested(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s
resource "pingone_password_policy" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s"
}

resource "pingone_population_default" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s"
  password_policy = {
    id = pingone_password_policy.%[3]s.id
  }
}`, acctest.MinimalSandboxEnvironmentNoPopulation(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccPopulationDefaultConfig_PasswordPolicyString(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s
resource "pingone_password_policy" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s"
}

resource "pingone_population_default" "%[3]s" {
  environment_id     = pingone_environment.%[2]s.id
  name               = "%[4]s"
  password_policy_id = pingone_password_policy.%[3]s.id
}`, acctest.MinimalSandboxEnvironmentNoPopulation(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccPopulationDefaultConfig_PasswordPolicyConflict(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s
resource "pingone_password_policy" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s"
}

resource "pingone_population_default" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s"
  password_policy = {
    id = pingone_password_policy.%[3]s.id
  }
  password_policy_id = pingone_password_policy.%[3]s.id
}`, acctest.MinimalSandboxEnvironmentNoPopulation(environmentName, licenseID), environmentName, resourceName, name)
}
