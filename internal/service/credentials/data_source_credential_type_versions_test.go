// Copyright © 2026 Ping Identity Corporation

package credentials_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	acctestlegacysdk "github.com/pingidentity/terraform-provider-pingone/internal/acctest/legacysdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/credentials"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccCredentialTypeVersionsDataSource_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullName := fmt.Sprintf("data.pingone_credential_type_versions.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := acctest.ResourceNameGen()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             credentials.CredentialType_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				// A newly created credential type may have zero or one versions, so
				// don't assert on the version list in this step.
				Config: testAccCredentialTypeVersionsDataSource_Full(environmentName, licenseID, resourceName, name, ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "credential_type_id", verify.P1ResourceIDRegexpFullString),
				),
			},
			{
				// Updating the credential type forces the service to save the
				// previous version and bump version.number, so at least two
				// versions should exist after this step. The list is returned
				// newest first, so versions.0 is the current version 2 and
				// versions.1 is the saved version 1.
				Config: testAccCredentialTypeVersionsDataSource_Full(environmentName, licenseID, resourceName, fmt.Sprintf("%s Updated", name), ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "versions.0.id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "versions.0.number", "2"),
					resource.TestMatchResourceAttr(dataSourceFullName, "versions.0.created_at", verify.RFC3339Regexp),
					resource.TestMatchResourceAttr(dataSourceFullName, "versions.1.id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "versions.1.number", "1"),
					resource.TestMatchResourceAttr(dataSourceFullName, "versions.1.created_at", verify.RFC3339Regexp),
				),
			},
			{
				// Filtering to `version eq 1` should return only the saved
				// original version.
				Config: testAccCredentialTypeVersionsDataSource_Full(environmentName, licenseID, resourceName, fmt.Sprintf("%s Updated", name), "version eq 1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceFullName, "versions.#", "1"),
					resource.TestCheckResourceAttr(dataSourceFullName, "versions.0.number", "1"),
					resource.TestMatchResourceAttr(dataSourceFullName, "versions.0.id", verify.P1ResourceIDRegexpFullString),
				),
			},
			{
				Config:  testAccCredentialTypeVersionsDataSource_Full(environmentName, licenseID, resourceName, fmt.Sprintf("%s Updated", name), ""),
				Destroy: true,
			},
		},
	})
}

func TestAccCredentialTypeVersionsDataSource_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             credentials.CredentialType_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccCredentialTypeVersionsDataSource_BadParameters_FilterFormat(resourceName),
				ExpectError: regexp.MustCompile("Error: Invalid Attribute Value"),
			},
		},
	})
}

func TestAccCredentialTypeVersionsDataSource_NotFound(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             credentials.CredentialType_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccCredentialTypeVersionsDataSource_NotFound(environmentName, licenseID, resourceName),
				ExpectError: regexp.MustCompile("Error: Error when calling `ReadAllCredentialTypeVersions`: The request could not be completed. The requested resource was not found."),
			},
		},
	})
}

func testAccCredentialTypeVersionsDataSource_Full(environmentName, licenseID, resourceName, name, filter string) string {
	// The filter attribute is emitted only when a filter is supplied, so the
	// unfiltered steps omit the attribute entirely (an empty string would
	// fail plan-time validation).
	filterAttribute := ""
	if filter != "" {
		filterAttribute = fmt.Sprintf("\n  filter = \"%s\"\n", filter)
	}

	return fmt.Sprintf(`
	%[1]s

resource "pingone_credential_type" "%[3]s" {
  environment_id       = pingone_environment.%[2]s.id
  title                = "%[4]s"
  description          = "%[4]s Example Description"
  card_type            = "%[4]s"
  card_design_template = "<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 740 480\"><rect fill=\"none\" width=\"736\" height=\"476\" stroke=\"#CACED3\" stroke-width=\"3\" rx=\"10\" ry=\"10\" x=\"2\" y=\"2\"></rect><rect fill=\"$${cardColor}\" height=\"476\" rx=\"10\" ry=\"10\" width=\"736\" x=\"2\" y=\"2\" opacity=\"$${bgOpacityPercent}\"></rect><line y2=\"160\" x2=\"695\" y1=\"160\" x1=\"42.5\" stroke=\"$${textColor}\"></line><text fill=\"$${textColor}\" font-weight=\"450\" font-size=\"30\" x=\"160\" y=\"90\">$${cardTitle}</text><text fill=\"$${textColor}\" font-size=\"25\" font-weight=\"300\" x=\"160\" y=\"130\">$${cardSubtitle}</text></svg>"

  metadata = {
    name               = "%[4]s"
    columns            = 1
    description        = "%[4]s Example Description"
    bg_opacity_percent = 100
    card_color         = "#000000"
    text_color         = "#eff0f1"

    fields = [
      {
        type       = "Directory Attribute"
        title      = "givenName"
        attribute  = "name.given"
        is_visible = false
      }
    ]
  }
  depends_on = [pingone_environment.%[2]s]
}

data "pingone_credential_type_versions" "%[3]s" {
  environment_id     = pingone_environment.%[2]s.id
  credential_type_id = pingone_credential_type.%[3]s.id
%[5]s
  depends_on = [pingone_credential_type.%[3]s]

}`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name, filterAttribute)
}

func testAccCredentialTypeVersionsDataSource_BadParameters_FilterFormat(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_credential_type_versions" "%[2]s" {
  environment_id     = data.pingone_environment.general_test.id
  credential_type_id = "9c052a8a-14be-44e4-8f07-2662569994ce" // dummy ID that conforms to UUID v4

  filter = "version ne 1"

}`, acctest.GenericSandboxEnvironment(), resourceName)
}

func testAccCredentialTypeVersionsDataSource_NotFound(environmentName, licenseID, resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_credential_type_versions" "%[3]s" {
  environment_id     = pingone_environment.%[2]s.id
  credential_type_id = "9c052a8a-14be-44e4-8f07-2662569994ce" // dummy ID that conforms to UUID v4

}`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName)
}
