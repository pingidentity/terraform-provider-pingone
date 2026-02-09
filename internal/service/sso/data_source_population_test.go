// Copyright © 2026 Ping Identity Corporation

package sso_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	acctestlegacysdk "github.com/pingidentity/terraform-provider-pingone/internal/acctest/legacysdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/sso"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccPopulationDataSource_ByNameFull(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_population.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)
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
				Config: testAccPopulationDataSourceConfig_ByNameFull(environmentName, licenseID, resourceName, name, false),
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
					resource.TestCheckResourceAttr(dataSourceFullName, "preferred_language", "pl"),
					resource.TestMatchResourceAttr(dataSourceFullName, "theme.id", verify.P1ResourceIDRegexpFullString),
				),
			},
			{
				Config: testAccPopulationDataSourceConfig_ByNameFull(environmentName, licenseID, resourceName, name, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "name", name),
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
				Config: testAccPopulationDataSourceConfig_ByIDFull(environmentName, licenseID, resourceName, name),
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
					resource.TestCheckResourceAttr(dataSourceFullName, "preferred_language", "pt"),
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
			acctest.PreCheckNoBeta(t)
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

func testAccPopulationDataSourceConfig_ByNameFull(environmentName, licenseID, resourceName, name string, insensitivityCheck bool) string {

	// If insensitivityCheck is true, alter the case of the name
	nameComparator := name
	if insensitivityCheck {
		nameComparator = acctest.AlterStringCasing(nameComparator)
	}

	return fmt.Sprintf(`
	%[1]s

resource "pingone_password_policy" "%[2]s" {
  environment_id = pingone_environment.%[4]s.id
  name           = "%[3]s"
}

data "pingone_language" "%[2]s" {
  environment_id = pingone_environment.%[4]s.id

  locale = "pl"
}

resource "pingone_language_update" "%[2]s" {
  environment_id = pingone_environment.%[4]s.id

  language_id = data.pingone_language.%[2]s.id
  enabled     = true
}

resource "pingone_branding_theme" "%[2]s" {
  environment_id = pingone_environment.%[4]s.id

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
  environment_id          = pingone_environment.%[4]s.id
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
  environment_id = pingone_environment.%[4]s.id

  name = "%[5]s"

  depends_on = [pingone_population.%[2]s-name]
}
`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), resourceName, name, environmentName, nameComparator)
}

func testAccPopulationDataSourceConfig_ByIDFull(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_password_policy" "%[2]s" {
  environment_id = pingone_environment.%[4]s.id
  name           = "%[3]s"
}

data "pingone_language" "%[2]s" {
  environment_id = pingone_environment.%[4]s.id

  locale = "pt"
}

resource "pingone_language_update" "%[2]s" {
  environment_id = pingone_environment.%[4]s.id

  language_id = data.pingone_language.%[2]s.id
  enabled     = true
}

resource "pingone_branding_theme" "%[2]s" {
  environment_id = pingone_environment.%[4]s.id

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
  environment_id = pingone_environment.%[4]s.id
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
  environment_id = pingone_environment.%[4]s.id

  population_id = pingone_population.%[2]s-name.id
}`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), resourceName, name, environmentName)
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
