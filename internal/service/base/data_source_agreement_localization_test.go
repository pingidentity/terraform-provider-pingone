package base_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccAgreementLocalizationDataSource_ByNameFull(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_agreement_localization.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckAgreementLocalizationDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAgreementLocalizationDataSourceConfig_ByNameFull(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(dataSourceFullName, "agreement_id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(dataSourceFullName, "agreement_localization_id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(dataSourceFullName, "language_id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(dataSourceFullName, "display_name", name),
					resource.TestCheckResourceAttr(dataSourceFullName, "locale", "en"),
					resource.TestCheckResourceAttr(dataSourceFullName, "enabled", "false"),
					resource.TestCheckResourceAttr(dataSourceFullName, "text_checkbox_accept", "Yeah"),
					resource.TestCheckResourceAttr(dataSourceFullName, "text_button_continue", "Move on"),
					resource.TestCheckResourceAttr(dataSourceFullName, "text_button_decline", "Nah"),
					resource.TestMatchResourceAttr(dataSourceFullName, "current_revision_id", verify.P1ResourceIDRegexp),
				),
			},
		},
	})
}

func TestAccAgreementLocalizationDataSource_ByLocaleFull(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_agreement_localization.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckAgreementLocalizationDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAgreementLocalizationDataSourceConfig_ByLocaleFull(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(dataSourceFullName, "agreement_id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(dataSourceFullName, "agreement_localization_id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(dataSourceFullName, "language_id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(dataSourceFullName, "display_name", name),
					resource.TestCheckResourceAttr(dataSourceFullName, "locale", "en"),
					resource.TestCheckResourceAttr(dataSourceFullName, "enabled", "false"),
					resource.TestCheckResourceAttr(dataSourceFullName, "text_checkbox_accept", "Yeah"),
					resource.TestCheckResourceAttr(dataSourceFullName, "text_button_continue", "Move on"),
					resource.TestCheckResourceAttr(dataSourceFullName, "text_button_decline", "Nah"),
					resource.TestMatchResourceAttr(dataSourceFullName, "current_revision_id", verify.P1ResourceIDRegexp),
				),
			},
		},
	})
}

func TestAccAgreementLocalizationDataSource_ByIDFull(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_agreement_localization.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckAgreementLocalizationDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAgreementLocalizationDataSourceConfig_ByIDFull(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(dataSourceFullName, "agreement_id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(dataSourceFullName, "agreement_localization_id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(dataSourceFullName, "language_id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(dataSourceFullName, "display_name", name),
					resource.TestCheckResourceAttr(dataSourceFullName, "locale", "en"),
					resource.TestCheckResourceAttr(dataSourceFullName, "enabled", "false"),
					resource.TestCheckResourceAttr(dataSourceFullName, "text_checkbox_accept", "Yeah"),
					resource.TestCheckResourceAttr(dataSourceFullName, "text_button_continue", "Move on"),
					resource.TestCheckResourceAttr(dataSourceFullName, "text_button_decline", "Nah"),
					resource.TestMatchResourceAttr(dataSourceFullName, "current_revision_id", verify.P1ResourceIDRegexp),
				),
			},
		},
	})
}

func TestAccAgreementLocalizationDataSource_NotFound(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckAgreementLocalizationDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccAgreementLocalizationDataSourceConfig_NotFoundByName(resourceName),
				ExpectError: regexp.MustCompile(`Cannot find agreement localization from name or locale`),
			},
			{
				Config:      testAccAgreementLocalizationDataSourceConfig_NotFoundByLocale(resourceName),
				ExpectError: regexp.MustCompile(`Cannot find agreement localization from name or locale`),
			},
			{
				Config:      testAccAgreementLocalizationDataSourceConfig_NotFoundByID(resourceName),
				ExpectError: regexp.MustCompile("Error when calling `ReadOneAgreementLanguage`: The request could not be completed. The requested resource was not found."),
			},
		},
	})
}

func testAccAgreementLocalizationDataSourceConfig_ByNameFull(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_agreement" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name                  = "%[3]s"
  description           = "Before the crowbar was invented, Crows would just drink at home."
  reconsent_period_days = 31

}

data "pingone_language" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  locale = "en"
}

resource "pingone_agreement_localization" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  agreement_id   = pingone_agreement.%[2]s.id
  language_id    = data.pingone_language.%[2]s.id

  display_name = "%[3]s"

  text_checkbox_accept = "Yeah"
  text_button_continue = "Move on"
  text_button_decline  = "Nah"
}

resource "pingone_agreement_localization_revision" "%[2]s" {
  environment_id            = data.pingone_environment.general_test.id
  agreement_id              = pingone_agreement.%[2]s.id
  agreement_localization_id = pingone_agreement_localization.%[2]s.id

  content_type      = "text/html"
  require_reconsent = true
  text              = <<EOT
	<h1>Test</h1>
  EOT

}

data "pingone_agreement_localization" "%[3]s" {
  environment_id = data.pingone_environment.general_test.id
  agreement_id   = pingone_agreement.%[2]s.id

  display_name = pingone_agreement_localization.%[2]s.display_name

  depends_on = [
    pingone_agreement_localization_revision.%[2]s
  ]
}
`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccAgreementLocalizationDataSourceConfig_ByLocaleFull(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_agreement" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name                  = "%[3]s"
  description           = "Before the crowbar was invented, Crows would just drink at home."
  reconsent_period_days = 31

}

data "pingone_language" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  locale = "en"
}

resource "pingone_agreement_localization" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  agreement_id   = pingone_agreement.%[2]s.id
  language_id    = data.pingone_language.%[2]s.id

  display_name = "%[3]s"

  text_checkbox_accept = "Yeah"
  text_button_continue = "Move on"
  text_button_decline  = "Nah"
}

resource "pingone_agreement_localization_revision" "%[2]s" {
  environment_id            = data.pingone_environment.general_test.id
  agreement_id              = pingone_agreement.%[2]s.id
  agreement_localization_id = pingone_agreement_localization.%[2]s.id

  content_type      = "text/html"
  require_reconsent = true
  text              = <<EOT
	<h1>Test</h1>
  EOT

}

data "pingone_agreement_localization" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  agreement_id   = pingone_agreement.%[2]s.id

  locale = pingone_agreement_localization.%[2]s.locale

  depends_on = [
    pingone_agreement_localization_revision.%[2]s
  ]
}
`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccAgreementLocalizationDataSourceConfig_ByIDFull(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_agreement" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name                  = "%[3]s"
  description           = "Before the crowbar was invented, Crows would just drink at home."
  reconsent_period_days = 31

}

data "pingone_language" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  locale = "en"
}

resource "pingone_agreement_localization" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  agreement_id   = pingone_agreement.%[2]s.id
  language_id    = data.pingone_language.%[2]s.id

  display_name = "%[3]s"

  text_checkbox_accept = "Yeah"
  text_button_continue = "Move on"
  text_button_decline  = "Nah"
}

resource "pingone_agreement_localization_revision" "%[2]s" {
  environment_id            = data.pingone_environment.general_test.id
  agreement_id              = pingone_agreement.%[2]s.id
  agreement_localization_id = pingone_agreement_localization.%[2]s.id

  content_type      = "text/html"
  require_reconsent = true
  text              = <<EOT
	<h1>Test</h1>
  EOT

}

data "pingone_agreement_localization" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  agreement_id   = pingone_agreement.%[2]s.id

  agreement_localization_id = pingone_agreement_localization.%[2]s.id

  depends_on = [
    pingone_agreement_localization_revision.%[2]s
  ]
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccAgreementLocalizationDataSourceConfig_NotFoundByName(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_agreement" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[2]s"
}

data "pingone_agreement_localization" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  agreement_id   = pingone_agreement.%[2]s.id

  display_name = "doesnotexist"
}
`, acctest.GenericSandboxEnvironment(), resourceName)
}

func testAccAgreementLocalizationDataSourceConfig_NotFoundByLocale(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_agreement" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[2]s"
}

data "pingone_agreement_localization" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  agreement_id   = pingone_agreement.%[2]s.id

  locale = "da"
}
`, acctest.GenericSandboxEnvironment(), resourceName)
}

func testAccAgreementLocalizationDataSourceConfig_NotFoundByID(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_agreement" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[2]s"
}

data "pingone_agreement_localization" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  agreement_id   = pingone_agreement.%[2]s.id

  agreement_localization_id = "9c052a8a-14be-44e4-8f07-2662569994ce" // dummy ID that conforms to UUID v4
}
`, acctest.GenericSandboxEnvironment(), resourceName)
}
