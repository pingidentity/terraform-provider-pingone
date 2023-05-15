package credentials_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
)

func TestAccCredentialTypeDataSource_ByIDFull(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_credential_type.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.pingone_credential_type.%s", resourceName)

	name := acctest.ResourceNameGen()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCredentialTypeDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCredentialTypeDataSource_ByIDFull(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceFullName, "id"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "environment_id"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "credential_type_id"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "title", resourceFullName, "title"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "description", resourceFullName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "card_type", resourceFullName, "card_type"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "card_design_template", resourceFullName, "card_design_template"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "metadata.%", resourceFullName, "metadata.%"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "metadata.fields.%", resourceFullName, "metadata.fields.%"),
				),
			},
			{
				Config:  testAccCredentialTypeDataSource_ByIDFull(resourceName, name),
				Destroy: true,
			},
		},
	})
}

func TestAccCredentialTypeDataSource_NotFound(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCredentialTypeDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccCredentialTypeDataSource_NotFoundByID(resourceName),
				ExpectError: regexp.MustCompile("Error: Cannot find credential type"),
			},
		},
	})
}

func TestAccCredentialTypeDataSourceInvalidConfig(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCredentialTypeDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccCredentialTypeDataSource_NoEnvironmentID(resourceName),
				ExpectError: regexp.MustCompile("Error: Missing required argument"),
			},
			{
				Config:      testAccCredentialTypeDataSource_NoCredentialTypeID(resourceName),
				ExpectError: regexp.MustCompile("Error: Credential Type ID not provided"),
			},
		},
	})
}

func testAccCredentialTypeDataSource_ByIDFull(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_credential_type" "%[2]s" {
  environment_id = data.pingone_environment.credentials_test.id
  title = "%[3]s"
  description = "%[3]s Example Description"
  card_type = "%[3]s"
  card_design_template = "<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 740 480\"><rect fill=\"none\" width=\"736\" height=\"476\" stroke=\"#CACED3\" stroke-width=\"3\" rx=\"10\" ry=\"10\" x=\"2\" y=\"2\"></rect><rect fill=\"$${cardColor}\" height=\"476\" rx=\"10\" ry=\"10\" width=\"736\" x=\"2\" y=\"2\" opacity=\"$${bgOpacityPercent}\"></rect><line y2=\"160\" x2=\"695\" y1=\"160\" x1=\"42.5\" stroke=\"$${textColor}\"></line><text fill=\"$${textColor}\" font-weight=\"450\" font-size=\"30\" x=\"160\" y=\"90\">$${cardTitle}</text><text fill=\"$${textColor}\" font-size=\"25\" font-weight=\"300\" x=\"160\" y=\"130\">$${cardSubtitle}</text></svg>"

  metadata = {
    name = "%[3]s"
	columns = 1
    description = "%[3]s Example Description"
    version = 5
    bg_opacity_percent = 100
    card_color = "#000000"
    text_color = "#eff0f1"

    fields = [
      {
        type = "Directory Attribute"
        title = "givenName"
        attribute = "name.given"
        is_visible = false
      },
      {
        type = "Directory Attribute"
        title = "surname"
        attribute = "name.family"
        is_visible = false
      },
      {
        type = "Directory Attribute"
        title = "jobTitle"
        attribute = "title"
        is_visible = false
      },
      {
        type = "Directory Attribute"
        title = "displayName"
        attribute = "name.formatted"
        is_visible = false
      },
      {
        type = "Directory Attribute"
        title = "mail"
        attribute = "email"
        is_visible = false
      },
      {
        type = "Directory Attribute"
        title = "preferredLanguage"
        attribute = "preferredLanguage"
        is_visible = false        
      },
      {
        type = "Directory Attribute"
        title = "id"
        attribute = "id"
        is_visible = false
      }
    ]
  }
}

data "pingone_credential_type" "%[2]s" {
	environment_id = data.pingone_environment.credentials_test.id
	credential_type_id = resource.pingone_credential_type.%[2]s.id

  }`, acctest.CredentialsSandboxEnvironment(), resourceName, name)
}

func testAccCredentialTypeDataSource_NotFoundByID(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_credential_type" "%[2]s" {
	environment_id = data.pingone_environment.credentials_test.id
	credential_type_id = "9c052a8a-14be-44e4-8f07-2662569994ce" // dummy ID that conforms to UUID v4

  }`, acctest.CredentialsSandboxEnvironment(), resourceName)
}

func testAccCredentialTypeDataSource_NoEnvironmentID(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_credential_type" "%[2]s" {
	credential_type_id = "9c052a8a-14be-44e4-8f07-2662569994ce" // dummy ID that conforms to UUID v4

  }`, acctest.CredentialsSandboxEnvironment(), resourceName)
}

func testAccCredentialTypeDataSource_NoCredentialTypeID(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_credential_type" "%[2]s" {
	environment_id = data.pingone_environment.credentials_test.id

  }`, acctest.CredentialsSandboxEnvironment(), resourceName)
}
