// Copyright © 2025 Ping Identity Corporation

package sso_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/sso"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccPopulationDataSource_ByNameFull(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_population.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Population_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPopulationDataSourceConfig_ByNameFull(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "population_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "name", name),
					resource.TestCheckResourceAttr(dataSourceFullName, "user_count", "0"),
					resource.TestCheckResourceAttr(dataSourceFullName, "description", "Test description"),
					resource.TestMatchResourceAttr(dataSourceFullName, "password_policy_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "password_policy.id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "default", "false"),
					resource.TestCheckResourceAttr(dataSourceFullName, "alternative_identifiers.#", "2"),
					resource.TestCheckTypeSetElemAttr(dataSourceFullName, "alternative_identifiers.*", "identifier1"),
					resource.TestCheckTypeSetElemAttr(dataSourceFullName, "alternative_identifiers.*", "identifier2"),
					resource.TestCheckResourceAttr(dataSourceFullName, "preferred_language", "es"),
					resource.TestMatchResourceAttr(dataSourceFullName, "theme.id", verify.P1ResourceIDRegexpFullString),
				),
			},
		},
	})
}

func TestAccPopulationDataSource_ByIDFull(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_population.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Population_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPopulationDataSourceConfig_ByIDFull(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "population_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "name", name),
					resource.TestCheckResourceAttr(dataSourceFullName, "user_count", "0"),
					resource.TestCheckResourceAttr(dataSourceFullName, "description", "Test description"),
					resource.TestMatchResourceAttr(dataSourceFullName, "password_policy_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "password_policy.id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "default", "false"),
					resource.TestCheckResourceAttr(dataSourceFullName, "alternative_identifiers.#", "2"),
					resource.TestCheckTypeSetElemAttr(dataSourceFullName, "alternative_identifiers.*", "identifier1"),
					resource.TestCheckTypeSetElemAttr(dataSourceFullName, "alternative_identifiers.*", "identifier2"),
					resource.TestCheckResourceAttr(dataSourceFullName, "preferred_language", "es"),
					resource.TestMatchResourceAttr(dataSourceFullName, "theme.id", verify.P1ResourceIDRegexpFullString),
				),
			},
		},
	})
}

func TestAccPopulationDataSource_NotFound(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Population_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccPopulationDataSourceConfig_NotFoundByName(resourceName),
				ExpectError: regexp.MustCompile(`Population not found`),
			},
			{
				Config:      testAccPopulationDataSourceConfig_NotFoundByID(resourceName),
				ExpectError: regexp.MustCompile("Population not found"),
			},
		},
	})
}

func testAccPopulationDataSourceConfig_ByNameFull(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_password_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
}

data "pingone_language" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  locale = "es"
}

resource "pingone_language_update" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  language_id = data.pingone_language.%[2]s.id
  enabled     = true
}

resource "pingone_branding_theme" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name     = "%[2]s"
  template = "split"

  background_color   = "#FF00F0"
  button_text_color  = "#FF6C6C"
  heading_text_color = "#FF0005"
  card_color         = "#0FFF39"
  body_text_color    = "#8620FF"
  link_text_color    = "#8A7F06"
  button_color       = "#0CFFFB"
}

resource "pingone_population" "%[2]s-name" {
  environment_id          = data.pingone_environment.general_test.id
  name                    = "%[3]s"
  description             = "Test description"
  password_policy_id      = pingone_password_policy.%[2]s.id
  preferred_language      = pingone_language_update.%[2]s.locale
  alternative_identifiers = ["identifier1", "identifier2"]
  theme = {
    id = pingone_branding_theme.%[2]s.id
  }
}

data "pingone_population" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = pingone_population.%[2]s-name.name
}
`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccPopulationDataSourceConfig_ByIDFull(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_password_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
}

data "pingone_language" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  locale = "es"
}

resource "pingone_language_update" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  language_id = data.pingone_language.%[2]s.id
  enabled     = true
}

resource "pingone_branding_theme" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name     = "%[2]s"
  template = "split"

  background_color   = "#FF00F0"
  button_text_color  = "#FF6C6C"
  heading_text_color = "#FF0005"
  card_color         = "#0FFF39"
  body_text_color    = "#8620FF"
  link_text_color    = "#8A7F06"
  button_color       = "#0CFFFB"
}

resource "pingone_population" "%[2]s-name" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Test description"
  password_policy = {
    id = pingone_password_policy.%[2]s.id
  }
  preferred_language      = pingone_language_update.%[2]s.locale
  alternative_identifiers = ["identifier1", "identifier2"]
  theme = {
    id = pingone_branding_theme.%[2]s.id
  }
}

data "pingone_population" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  population_id = pingone_population.%[2]s-name.id
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccPopulationDataSourceConfig_NotFoundByName(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_population" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "doesnotexist"
}
`, acctest.GenericSandboxEnvironment(), resourceName)
}

func testAccPopulationDataSourceConfig_NotFoundByID(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_population" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  population_id = "9c052a8a-14be-44e4-8f07-2662569994ce" // dummy ID that conforms to UUID v4
}
`, acctest.GenericSandboxEnvironment(), resourceName)
}
