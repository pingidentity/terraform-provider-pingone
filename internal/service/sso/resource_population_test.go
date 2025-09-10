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
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/base"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/sso"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccPopulation_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_population.%s", resourceName)

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
			acctest.PreCheckNoBeta(t)
			p1Client = acctest.PreCheckTestClient(ctx, t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Population_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccPopulationConfig_Minimal(resourceName, name),
				Check:  sso.Population_GetIDs(resourceFullName, &environmentID, &populationID),
			},
			{
				PreConfig: func() {
					sso.Population_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID, populationID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the environment
			{
				Config: testAccPopulationConfig_NewEnv(environmentName, licenseID, resourceName, name),
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

func TestAccPopulation_NewEnv(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_population.%s", resourceName)

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
		CheckDestroy:             sso.Population_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPopulationConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "preferred_language", "en"),
					resource.TestMatchResourceAttr(resourceFullName, "theme.id", verify.P1ResourceIDRegexpFullString),
				),
			},
		},
	})
}

func TestAccPopulation_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_population.%s", resourceName)
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
		CheckDestroy:             sso.Population_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPopulationConfig_Full(environmentName, licenseID, resourceName, name),
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

						return fmt.Sprintf("%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccPopulation_Minimal(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_population.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Population_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPopulationConfig_Minimal(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckNoResourceAttr(resourceFullName, "description"),
					resource.TestCheckNoResourceAttr(resourceFullName, "password_policy.id"),
					resource.TestCheckResourceAttr(resourceFullName, "preferred_language", "en"),
					resource.TestMatchResourceAttr(resourceFullName, "theme.id", verify.P1ResourceIDRegexpFullString),
				),
			},
		},
	})
}

func TestAccPopulation_PasswordPolicy(t *testing.T) {
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
		CheckDestroy:             sso.Population_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPopulationConfig_PasswordPolicyNested(resourceName, name),
			},
			{
				Config: testAccPopulationConfig_PasswordPolicyString(resourceName, name),
			},
			{
				Config:      testAccPopulationConfig_PasswordPolicyConflict(resourceName, name),
				ExpectError: regexp.MustCompile(`Invalid Attribute Combination`),
			},
		},
	})
}

func TestAccPopulation_DataProtection(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_population.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	var p1Client *client.Client
	var ctx = context.Background()

	var populationID, environmentID, userID string

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoBeta(t)
			p1Client = acctest.PreCheckTestClient(ctx, t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Population_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPopulationConfig_DataProtection_Sandbox(environmentName, licenseID, resourceName, name),
				Check:  sso.Population_GetIDs(resourceFullName, &environmentID, &populationID),
			},
			{
				Config: testAccPopulationConfig_DataProtection_Sandbox(environmentName, licenseID, resourceName, name),
				PreConfig: func() {
					sso.User_CreateUser_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID, name, &userID, &populationID)
				},
				Destroy:     true,
				ExpectError: regexp.MustCompile("Error when calling `DeletePopulation`: The request could not be completed"),
			},
			{
				Config: testAccPopulationConfig_DataProtection_Sandbox(environmentName, licenseID, resourceName, name),
				PreConfig: func() {
					sso.User_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID, userID)
				},
				Destroy: true,
			},
			{
				Config: testAccPopulationConfig_DataProtection_Production(environmentName, licenseID, resourceName, name),
				Check:  sso.Population_GetIDs(resourceFullName, &environmentID, &populationID),
			},
			{
				Config: testAccPopulationConfig_DataProtection_Production(environmentName, licenseID, resourceName, name),
				PreConfig: func() {
					sso.User_CreateUser_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID, name, &userID, &populationID)
				},
				Destroy:     true,
				ExpectError: regexp.MustCompile("Error when calling `DeletePopulation`: The request could not be completed"),
			},
			{
				Config: testAccPopulationConfig_DataProtection_Sandbox(environmentName, licenseID, resourceName, name),
				PreConfig: func() {
					sso.User_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID, userID)
				},
				ExpectError: regexp.MustCompile("Data protection notice - The environment type cannot be changed from PRODUCTION to SANDBOX"),
			},
		},
	})
}

func TestAccPopulation_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_population.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Population_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccPopulationConfig_Minimal(resourceName, name),
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

func testAccPopulationConfig_NewEnv(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_population" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s"
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccPopulationConfig_Full(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_password_policy" "%[2]s" {
  environment_id = pingone_environment.%[4]s.id
  name           = "%[3]s"
}

data "pingone_language" "%[3]s" {
  environment_id = pingone_environment.%[4]s.id

  locale = "es"
}

resource "pingone_language_update" "%[3]s" {
  environment_id = pingone_environment.%[4]s.id

  language_id = data.pingone_language.%[3]s.id
  enabled     = true
}

resource "pingone_branding_theme" "%[3]s" {
  environment_id = pingone_environment.%[4]s.id

  name     = "%[3]s"
  template = "split"

  background_color   = "#FF00F0"
  button_text_color  = "#FF6C6C"
  heading_text_color = "#FF0005"
  card_color         = "#0FFF39"
  body_text_color    = "#8620FF"
  link_text_color    = "#8A7F06"
  button_color       = "#0CFFFB"
}

resource "pingone_population" "%[2]s" {
  environment_id = pingone_environment.%[4]s.id
  name           = "%[3]s"
  description    = "Test description"
  password_policy = {
    id = pingone_password_policy.%[2]s.id
  }
  preferred_language      = pingone_language_update.%[3]s.locale
  alternative_identifiers = ["identifier1", "identifier2"]
  theme = {
    id = pingone_branding_theme.%[3]s.id
  }
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), resourceName, name, environmentName)
}

func testAccPopulationConfig_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_population" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccPopulationConfig_PasswordPolicyNested(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s
resource "pingone_password_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
}

resource "pingone_population" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  password_policy = {
    id = pingone_password_policy.%[2]s.id
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccPopulationConfig_PasswordPolicyString(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s
resource "pingone_password_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
}

resource "pingone_population" "%[2]s" {
  environment_id     = data.pingone_environment.general_test.id
  name               = "%[3]s"
  password_policy_id = pingone_password_policy.%[2]s.id
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccPopulationConfig_PasswordPolicyConflict(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s
resource "pingone_password_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
}

resource "pingone_population" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  password_policy = {
    id = pingone_password_policy.%[2]s.id
  }
  password_policy_id = pingone_password_policy.%[2]s.id
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccPopulationConfig_DataProtection_Sandbox(environmentName, licenseID, resourceName, name string) string {
	return testAccPopulationConfig_DataProtection(environmentName, licenseID, resourceName, name, management.ENUMENVIRONMENTTYPE_SANDBOX)
}

func testAccPopulationConfig_DataProtection_Production(environmentName, licenseID, resourceName, name string) string {
	return testAccPopulationConfig_DataProtection(environmentName, licenseID, resourceName, name, management.ENUMENVIRONMENTTYPE_PRODUCTION)
}

func testAccPopulationConfig_DataProtection(environmentName, licenseID, resourceName, name string, environmentType management.EnumEnvironmentType) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_population" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s"
}`, acctest.MinimalEnvironmentNoPopulation(environmentName, licenseID, environmentType), environmentName, resourceName, name)
}
