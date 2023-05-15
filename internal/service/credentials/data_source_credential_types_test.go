package credentials_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccCredentialTypesDataSource_NoFilter(t *testing.T) {
	// If run in parallel, unique environments are needed to prevent collisions within the same environment.
	//t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullName := fmt.Sprintf("data.pingone_credential_types.%s", resourceName)

	name := acctest.ResourceNameGen()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCredentialTypeDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCredentialTypesDataSource_NoFilter(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(dataSourceFullName, "ids.#", "3"),
					resource.TestMatchResourceAttr(dataSourceFullName, "ids.0", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(dataSourceFullName, "ids.1", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(dataSourceFullName, "ids.2", verify.P1ResourceIDRegexp),
				),
			},
			{
				Config:  testAccCredentialTypesDataSource_NoFilter(resourceName, name),
				Destroy: true,
			},
		},
	})
}

func TestAccCredentialTypesDataSource_NotFound(t *testing.T) {
	// If run in parallel, unique environments are needed to prevent collisions within the same environment.
	//t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullName := fmt.Sprintf("data.pingone_credential_types.%s", resourceName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCredentialTypeDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCredentialTypesDataSource_NotFound(resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(dataSourceFullName, "ids.#", "0"),
				),
			},
		},
	})
}

func testAccCredentialTypesDataSource_NoFilter(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_credential_type" "%[2]s-1" {
  environment_id       = data.pingone_environment.credentials_test.id
  title                = "%[3]s"
  description          = "%[3]s Example Description"
  card_type            = "%[3]s"
  card_design_template = "<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 740 480\"><rect fill=\"none\" width=\"736\" height=\"476\" stroke=\"#CACED3\" stroke-width=\"3\" rx=\"10\" ry=\"10\" x=\"2\" y=\"2\"></rect><rect fill=\"$${cardColor}\" height=\"476\" rx=\"10\" ry=\"10\" width=\"736\" x=\"2\" y=\"2\" opacity=\"$${bgOpacityPercent}\"></rect><line y2=\"160\" x2=\"695\" y1=\"160\" x1=\"42.5\" stroke=\"$${textColor}\"></line><text fill=\"$${textColor}\" font-weight=\"450\" font-size=\"30\" x=\"160\" y=\"90\">$${cardTitle}</text><text fill=\"$${textColor}\" font-size=\"25\" font-weight=\"300\" x=\"160\" y=\"130\">$${cardSubtitle}</text></svg>"

  metadata = {
    name               = "%[3]s"
    columns            = 1
    description        = "%[3]s Example Description"
    version            = 5
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
}

resource "pingone_credential_type" "%[2]s-2" {
  environment_id       = data.pingone_environment.credentials_test.id
  title                = "%[3]s"
  description          = "%[3]s Example Description"
  card_type            = "%[3]s"
  card_design_template = "<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 740 480\"><rect fill=\"none\" width=\"736\" height=\"476\" stroke=\"#CACED3\" stroke-width=\"3\" rx=\"10\" ry=\"10\" x=\"2\" y=\"2\"></rect><rect fill=\"$${cardColor}\" height=\"476\" rx=\"10\" ry=\"10\" width=\"736\" x=\"2\" y=\"2\" opacity=\"$${bgOpacityPercent}\"></rect><line y2=\"160\" x2=\"695\" y1=\"160\" x1=\"42.5\" stroke=\"$${textColor}\"></line><text fill=\"$${textColor}\" font-weight=\"450\" font-size=\"30\" x=\"160\" y=\"90\">$${cardTitle}</text><text fill=\"$${textColor}\" font-size=\"25\" font-weight=\"300\" x=\"160\" y=\"130\">$${cardSubtitle}</text></svg>"

  metadata = {
    name               = "%[3]s"
    columns            = 1
    description        = "%[3]s Example Description"
    version            = 5
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
}

resource "pingone_credential_type" "%[2]s-3" {
  environment_id       = data.pingone_environment.credentials_test.id
  title                = "%[3]s"
  description          = "%[3]s Example Description"
  card_type            = "%[3]s"
  card_design_template = "<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 740 480\"><rect fill=\"none\" width=\"736\" height=\"476\" stroke=\"#CACED3\" stroke-width=\"3\" rx=\"10\" ry=\"10\" x=\"2\" y=\"2\"></rect><rect fill=\"$${cardColor}\" height=\"476\" rx=\"10\" ry=\"10\" width=\"736\" x=\"2\" y=\"2\" opacity=\"$${bgOpacityPercent}\"></rect><line y2=\"160\" x2=\"695\" y1=\"160\" x1=\"42.5\" stroke=\"$${textColor}\"></line><text fill=\"$${textColor}\" font-weight=\"450\" font-size=\"30\" x=\"160\" y=\"90\">$${cardTitle}</text><text fill=\"$${textColor}\" font-size=\"25\" font-weight=\"300\" x=\"160\" y=\"130\">$${cardSubtitle}</text></svg>"

  metadata = {
    name               = "%[3]s"
    columns            = 1
    description        = "%[3]s Example Description"
    version            = 5
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
}
data "pingone_credential_types" "%[2]s" {
  environment_id = data.pingone_environment.credentials_test.id

  depends_on = [pingone_credential_type.%[2]s-1, pingone_credential_type.%[2]s-2, pingone_credential_type.%[2]s-3]

}`, acctest.CredentialsSandboxEnvironment(), resourceName, name)
}

func testAccCredentialTypesDataSource_NotFound(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_credential_types" "%[2]s" {
  environment_id = data.pingone_environment.credentials_test.id

}`, acctest.CredentialsSandboxEnvironment(), resourceName)
}
