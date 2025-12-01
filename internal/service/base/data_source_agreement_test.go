// Copyright © 2025 Ping Identity Corporation

package base_test

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/base"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccAgreementDataSource_ByNameFull(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_agreement.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Agreement_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAgreementDataSourceConfig_ByNameFull(resourceName, name, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "agreement_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "name", name),
					resource.TestCheckResourceAttr(dataSourceFullName, "enabled", "false"),
					resource.TestCheckResourceAttr(dataSourceFullName, "description", "Before the crowbar was invented, Crows would just drink at home."),
					resource.TestCheckResourceAttr(dataSourceFullName, "reconsent_period_days", "31"),
					resource.TestCheckResourceAttr(dataSourceFullName, "total_user_consent_count", "0"),
					resource.TestCheckResourceAttr(dataSourceFullName, "expired_user_consent_count", "0"),
					resource.TestMatchResourceAttr(dataSourceFullName, "consent_counts_updated_at", regexp.MustCompile(`^(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d\.\d+)|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d)|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d)$`)),
				),
			},
			// Case insensitivity check
			{
				Config: testAccAgreementDataSourceConfig_ByNameFull(resourceName, name, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "name", name),
				),
			},
		},
	})
}

func TestAccAgreementDataSource_ByIDFull(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_agreement.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Agreement_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAgreementDataSourceConfig_ByIDFull(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "agreement_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "name", name),
					resource.TestCheckResourceAttr(dataSourceFullName, "enabled", "false"),
					resource.TestCheckResourceAttr(dataSourceFullName, "description", "Before the crowbar was invented, Crows would just drink at home."),
					resource.TestCheckResourceAttr(dataSourceFullName, "reconsent_period_days", "31"),
					resource.TestCheckResourceAttr(dataSourceFullName, "total_user_consent_count", "0"),
					resource.TestCheckResourceAttr(dataSourceFullName, "expired_user_consent_count", "0"),
					resource.TestMatchResourceAttr(dataSourceFullName, "consent_counts_updated_at", regexp.MustCompile(`^(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d\.\d+)|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d)|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d)$`)),
				),
			},
		},
	})
}

func TestAccAgreementDataSource_NotFound(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Agreement_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccAgreementDataSourceConfig_NotFoundByName(resourceName),
				ExpectError: regexp.MustCompile(`Cannot find agreement from name`),
			},
			{
				Config:      testAccAgreementDataSourceConfig_NotFoundByID(resourceName),
				ExpectError: regexp.MustCompile("Error when calling `ReadOneAgreement`: The request could not be completed. The requested resource was not found."),
			},
		},
	})
}

func testAccAgreementDataSourceConfig_ByNameFull(resourceName, name string, insensitivityCheck bool) string {
	date := time.Now().In(time.UTC).Add(time.Hour * time.Duration(1))

	// If insensitivityCheck is true, alter the case of the name
	nameComparator := name
	if insensitivityCheck {
		nameComparator = acctest.AlterStringCasing(nameComparator)
	}

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
}

resource "pingone_agreement_localization_revision" "%[2]s" {
  environment_id            = data.pingone_environment.general_test.id
  agreement_id              = pingone_agreement.%[2]s.id
  agreement_localization_id = pingone_agreement_localization.%[2]s.id

  content_type      = "text/html"
  effective_at      = "%[5]s"
  require_reconsent = true
  text              = <<EOT
  <h1>Test</h1>
EOT

}

data "pingone_agreement" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[4]s"

  depends_on = [
    pingone_agreement_localization_revision.%[2]s
  ]
}
`, acctest.GenericSandboxEnvironment(), resourceName, name, nameComparator, date.Format(time.RFC3339))
}

func testAccAgreementDataSourceConfig_ByIDFull(resourceName, name string) string {
	date := time.Now().In(time.UTC).Add(time.Hour * time.Duration(1))

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
}

resource "pingone_agreement_localization_revision" "%[2]s" {
  environment_id            = data.pingone_environment.general_test.id
  agreement_id              = pingone_agreement.%[2]s.id
  agreement_localization_id = pingone_agreement_localization.%[2]s.id

  content_type      = "text/html"
  effective_at      = "%[4]s"
  require_reconsent = true
  text              = <<EOT
	  <h1>Test</h1>
EOT

}

data "pingone_agreement" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  agreement_id = pingone_agreement.%[2]s.id

  depends_on = [
    pingone_agreement_localization_revision.%[2]s
  ]
}`, acctest.GenericSandboxEnvironment(), resourceName, name, date.Format(time.RFC3339))
}

func testAccAgreementDataSourceConfig_NotFoundByName(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_agreement" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "doesnotexist"
}
`, acctest.GenericSandboxEnvironment(), resourceName)
}

func testAccAgreementDataSourceConfig_NotFoundByID(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_agreement" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  agreement_id = "9c052a8a-14be-44e4-8f07-2662569994ce" // dummy ID that conforms to UUID v4
}
`, acctest.GenericSandboxEnvironment(), resourceName)
}
