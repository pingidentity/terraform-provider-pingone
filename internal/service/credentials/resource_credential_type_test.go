package credentials_test

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/patrickcping/pingone-go-sdk-v2/credentials"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func testAccCheckCredentialTypeDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.CredentialsAPIClient
	ctx = context.WithValue(ctx, credentials.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})

	mgmtApiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_credential_type" {
			continue
		}

		_, rEnv, err := mgmtApiClient.EnvironmentsApi.ReadOneEnvironment(ctx, rs.Primary.Attributes["environment_id"]).Execute()

		if err != nil {

			if rEnv == nil {
				return fmt.Errorf("Response object does not exist and no error detected")
			}

			if rEnv.StatusCode == 404 {
				continue
			}

			return err
		}

		body, r, err := apiClient.CredentialTypesApi.ReadOneCredentialType(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["id"]).Execute()

		if err != nil {

			if r == nil {
				return fmt.Errorf("Response object does not exist and no error detected")
			}

			if r.StatusCode == 404 {
				continue
			}

			tflog.Error(ctx, fmt.Sprintf("Error: %v", body))
			return err
		} else {

			if body.DeletedAt != nil {

				// Note: Credential Types are "soft delted" and may be returned via the ReadOneCredentialType call.
				// If the DeletedAt attribute exists, it is considered deleted, handle similar to a 404.
				return err
			}

		}

		return fmt.Errorf("PingOne Credential Type ID %s still exists", rs.Primary.ID)
	}

	return nil
}
func TestAccCredentialType_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_credential_type.%s", resourceName)

	name := acctest.ResourceNameGen()

	data, _ := ioutil.ReadFile("../../acctest/test_assets/image/credential_background_base64.png")
	backgroundImage := base64.StdEncoding.EncodeToString(data) //string(data)

	data, _ = ioutil.ReadFile("../../acctest/test_assets/image/credential_logo_base64.png")
	logoImage := base64.StdEncoding.EncodeToString(data) //string(data)

	fullStep := resource.TestStep{
		Config: testAccCredentialTypeConfig_Full(resourceName, name, backgroundImage, logoImage),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
			resource.TestCheckResourceAttr(resourceFullName, "title", name),
			resource.TestCheckResourceAttr(resourceFullName, "description", fmt.Sprintf("%s Example Description", name)),
			resource.TestCheckResourceAttr(resourceFullName, "card_type", "VerifiedEmployee"),
			resource.TestCheckResourceAttrSet(resourceFullName, "card_design_template"),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.name", name),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.description", fmt.Sprintf("%s Example Description", name)),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.version", "5"),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.columns", "1"),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.background_image", backgroundImage),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.logo_image", logoImage),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.bg_opacity_percent", "100"),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.card_color", "#000000"),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.text_color", "#eff0f1"),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.fields.#", "7"),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.fields.3.id", "Directory Attribute -> displayName"),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.fields.3.type", "Directory Attribute"),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.fields.3.title", "displayName"),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.fields.3.attribute", "name.formatted"),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.fields.3.is_visible", "false"),
		),
	}

	updatedName := acctest.ResourceNameGen()

	minimalStep := resource.TestStep{
		Config: testAccCredentialTypeConfig_Minimal(resourceName, updatedName),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
			resource.TestCheckResourceAttr(resourceFullName, "title", updatedName),
			resource.TestCheckResourceAttr(resourceFullName, "description", fmt.Sprintf("%s Example Description", updatedName)),
			resource.TestCheckResourceAttr(resourceFullName, "card_type", "DemonstrationCard"),
			resource.TestCheckResourceAttrSet(resourceFullName, "card_design_template"),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.name", updatedName),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.version", "5"),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.columns", "3"),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.bg_opacity_percent", "100"),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.card_color", "#000000"),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.text_color", "#eff0f1"),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.fields.#", "3"),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.fields.2.id", "Alphanumeric Text -> Company"),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.fields.2.type", "Alphanumeric Text"),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.fields.2.title", "Company"),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.fields.2.value", "Example"),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.fields.2.is_visible", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.fields.0.id", "Issued Timestamp -> timestamp"),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.fields.0.type", "Issued Timestamp"),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.fields.0.title", "timestamp"),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.fields.0.is_visible", "false"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCredentialTypeDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// full
			fullStep,
			{
				Config:  testAccCredentialTypeConfig_Full(resourceName, name, backgroundImage, logoImage),
				Destroy: true,
			},
			minimalStep,
			{
				Config:  testAccCredentialTypeConfig_Minimal(resourceName, updatedName),
				Destroy: true,
			},
			fullStep,
			minimalStep,
			fullStep,
		},
	})
}

func TestAccCredentialType_MetaData(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	//resourceFullName := fmt.Sprintf("pingone_credential_type.%s", resourceName)

	name := acctest.ResourceNameGen()

	data, _ := ioutil.ReadFile("../../acctest/test_assets/image/image-background.jpg") // >90kb
	backgroundImage := base64.StdEncoding.EncodeToString(data)                         //string(data)

	data, _ = ioutil.ReadFile("../../acctest/test_assets/image/image-logo.gif") // >50kb
	logoImage := base64.StdEncoding.EncodeToString(data)                        //string(data)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCredentialTypeDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccCredentialTypeConfig_InvalidTitle(resourceName, ""),
				ExpectError: regexp.MustCompile("Error: Invalid Attribute Value Length"),
			},
			{
				Config:      testAccCredentialTypeConfig_InvalidCardColorHexValue(resourceName, name),
				ExpectError: regexp.MustCompile("Error: Invalid Attribute Value Match"),
			},
			{
				Config:      testAccCredentialTypeConfig_InvalidTextColorHexValue(resourceName, name),
				ExpectError: regexp.MustCompile("Error: Invalid Attribute Value Match"),
			},
			{
				Config:      testAccCredentialTypeConfig_InvalidBackgroundOpacityValue(resourceName, name),
				ExpectError: regexp.MustCompile("Error: Invalid Attribute Value"),
			},
			{
				Config:      testAccCredentialTypeConfig_InvalidVersion(resourceName, name),
				ExpectError: regexp.MustCompile("Error: Invalid Attribute Value"),
			},
			{
				Config:      testAccCredentialTypeConfig_NoSVGTagCardDesignTemplate(resourceName, name),
				ExpectError: regexp.MustCompile("Error: Invalid Attribute Value"),
			},
			{
				Config:      testAccCredentialTypeConfig_EmptyFieldsArray(resourceName, name),
				ExpectError: regexp.MustCompile("Error: Error when calling `CreateCredentialType`: Validation Error : \\[metadata.fields must not be empty\\]"),
			},
			{
				Config:      testAccCredentialTypeConfig_ImageSizeExceeded(resourceName, name, backgroundImage, logoImage),
				ExpectError: regexp.MustCompile("Error: Invalid Attribute Value Length"),
			},
		},
	})
}

func testAccCredentialTypeConfig_Full(resourceName, name, backgroundImage, logoImage string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_credential_type" "%[2]s" {
  environment_id = data.pingone_environment.credentials_test.id
  title = "%[3]s"
  description = "%[3]s Example Description"
  card_type = "VerifiedEmployee"
  card_design_template = "<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 740 480\"><rect fill=\"none\" width=\"736\" height=\"476\" stroke=\"#CACED3\" stroke-width=\"3\" rx=\"10\" ry=\"10\" x=\"2\" y=\"2\"></rect><rect fill=\"$${cardColor}\" height=\"476\" rx=\"10\" ry=\"10\" width=\"736\" x=\"2\" y=\"2\" opacity=\"$${bgOpacityPercent}\"></rect><image href=\"$${backgroundImage}\" opacity=\"$${bgOpacityPercent}\" height=\"476\" rx=\"10\" ry=\"10\" width=\"736\" x=\"2\" y=\"2\"></image><image href=\"$${logoImage}\" x=\"42\" y=\"43\" height=\"90px\" width=\"90px\"></image><line y2=\"160\" x2=\"695\" y1=\"160\" x1=\"42.5\" stroke=\"$${textColor}\"></line><text fill=\"$${textColor}\" font-weight=\"450\" font-size=\"30\" x=\"160\" y=\"90\">$${cardTitle}</text><text fill=\"$${textColor}\" font-size=\"25\" font-weight=\"300\" x=\"160\" y=\"130\">$${cardSubtitle}</text></svg>"

  metadata = {
    name = "%[3]s"
    description = "%[3]s Example Description"
    version = 5
	columns = 1
	background_image = "%[4]s"
	logo_image = "%[5]s"
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
}`, acctest.CredentialsSandboxEnvironment(), resourceName, name, backgroundImage, logoImage)
}

func testAccCredentialTypeConfig_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_credential_type" "%[2]s" {
  environment_id = data.pingone_environment.credentials_test.id
  title = "%[3]s"
  description = "%[3]s Example Description"
  card_type = "DemonstrationCard"
  card_design_template = "<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 740 480\"><rect fill=\"none\" width=\"736\" height=\"476\" stroke=\"#CACED3\" stroke-width=\"3\" rx=\"10\" ry=\"10\" x=\"2\" y=\"2\"></rect><rect fill=\"$${cardColor}\" height=\"476\" rx=\"10\" ry=\"10\" width=\"736\" x=\"2\" y=\"2\" opacity=\"$${bgOpacityPercent}\"></rect><line y2=\"160\" x2=\"695\" y1=\"160\" x1=\"42.5\" stroke=\"$${textColor}\"></line><text fill=\"$${textColor}\" font-weight=\"450\" font-size=\"30\" x=\"160\" y=\"90\">$${cardTitle}</text></svg>"

  metadata = {
    name = "%[3]s"
    version = 5
	columns = 3
    bg_opacity_percent = 100
    card_color = "#000000"
    text_color = "#eff0f1"

    fields = [
      {
        type = "Issued Timestamp"
        title = "timestamp"
        is_visible = false
      },
      {
        type = "Directory Attribute"
        title = "surname"
        attribute = "name.family"
        is_visible = false
      },
      {
        type = "Alphanumeric Text"
        title = "Company"
        value = "Example"
        is_visible = false
      }
    ]
  }
}`, acctest.CredentialsSandboxEnvironment(), resourceName, name)
}

func testAccCredentialTypeConfig_ImageSizeExceeded(resourceName, name, backgroundImage, logoImage string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_credential_type" "%[2]s" {
  environment_id = data.pingone_environment.credentials_test.id
  title = "%[3]s"
  description = "%[3]s Example Description"
  card_type = "VerifiedEmployee"
  card_design_template = "<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 740 480\"><rect fill=\"none\" width=\"736\" height=\"476\" stroke=\"#CACED3\" stroke-width=\"3\" rx=\"10\" ry=\"10\" x=\"2\" y=\"2\"></rect><rect fill=\"$${cardColor}\" height=\"476\" rx=\"10\" ry=\"10\" width=\"736\" x=\"2\" y=\"2\" opacity=\"$${bgOpacityPercent}\"></rect><image href=\"$${backgroundImage}\" opacity=\"$${bgOpacityPercent}\" height=\"476\" rx=\"10\" ry=\"10\" width=\"736\" x=\"2\" y=\"2\"></image><image href=\"$${logoImage}\" x=\"42\" y=\"43\" height=\"90px\" width=\"90px\"></image><line y2=\"160\" x2=\"695\" y1=\"160\" x1=\"42.5\" stroke=\"$${textColor}\"></line><text fill=\"$${textColor}\" font-weight=\"450\" font-size=\"30\" x=\"160\" y=\"90\">$${cardTitle}</text><text fill=\"$${textColor}\" font-size=\"25\" font-weight=\"300\" x=\"160\" y=\"130\">$${cardSubtitle}</text></svg>"

  metadata = {
    name = "%[3]s"
    description = "%[3]s Example Description"
    version = 5
	background_image = "%[4]s"
	logo_image = "%[5]s"
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
}`, acctest.CredentialsSandboxEnvironment(), resourceName, name, backgroundImage, logoImage)
}

func testAccCredentialTypeConfig_InvalidTitle(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_credential_type" "%[2]s" {
  environment_id = data.pingone_environment.credentials_test.id
  title = "%[3]s"
  description = "%[3]s Example Description"
  card_type = "DemonstrationCard"
  card_design_template = "<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 740 480\"><rect fill=\"none\" width=\"736\" height=\"476\" stroke=\"#CACED3\" stroke-width=\"3\" rx=\"10\" ry=\"10\" x=\"2\" y=\"2\"></rect><rect fill=\"$${cardColor}\" height=\"476\" rx=\"10\" ry=\"10\" width=\"736\" x=\"2\" y=\"2\" opacity=\"$${bgOpacityPercent}\"></rect><line y2=\"160\" x2=\"695\" y1=\"160\" x1=\"42.5\" stroke=\"$${textColor}\"></line><text fill=\"$${textColor}\" font-weight=\"450\" font-size=\"30\" x=\"160\" y=\"90\">$${cardTitle}</text></svg>"

  metadata = {
    name = "%[3]s"
    version = 5
    bg_opacity_percent = 100
    card_color = "#000000"
    text_color = "#eff0f1"

    fields = [
      {
        type = "Issued Timestamp"
        title = "timestamp"
        is_visible = false
      },
      {
        type = "Directory Attribute"
        title = "surname"
        attribute = "name.family"
        is_visible = false
      },
      {
        type = "Alphanumeric Text"
        title = "Company"
        value = "Example"
        is_visible = false
      }
    ]
  }
}`, acctest.CredentialsSandboxEnvironment(), resourceName, name)
}

func testAccCredentialTypeConfig_InvalidCardColorHexValue(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_credential_type" "%[2]s" {
  environment_id = data.pingone_environment.credentials_test.id
  title = "%[3]s"
  description = "%[3]s Example Description"
  card_type = "DemonstrationCard"
  card_design_template = "<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 740 480\"><rect fill=\"none\" width=\"736\" height=\"476\" stroke=\"#CACED3\" stroke-width=\"3\" rx=\"10\" ry=\"10\" x=\"2\" y=\"2\"></rect><rect fill=\"$${cardColor}\" height=\"476\" rx=\"10\" ry=\"10\" width=\"736\" x=\"2\" y=\"2\" opacity=\"$${bgOpacityPercent}\"></rect><line y2=\"160\" x2=\"695\" y1=\"160\" x1=\"42.5\" stroke=\"$${textColor}\"></line><text fill=\"$${textColor}\" font-weight=\"450\" font-size=\"30\" x=\"160\" y=\"90\">$${cardTitle}</text></svg>"

  metadata = {
    name = "%[3]s"
    version = 5
    bg_opacity_percent = 100
    card_color = "InvalidColor"
    text_color = "#eff0f1"

    fields = [
      {
        type = "Issued Timestamp"
        title = "timestamp"
        is_visible = false
      },
    ]
  }
}`, acctest.CredentialsSandboxEnvironment(), resourceName, name)
}

func testAccCredentialTypeConfig_InvalidTextColorHexValue(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s
resource "pingone_credential_type" "%[2]s" {
	environment_id = data.pingone_environment.credentials_test.id
	title = "%[3]s"
	description = "%[3]s Example Description"
	card_type = "DemonstrationCard"
	card_design_template = "<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 740 480\"><rect fill=\"none\" width=\"736\" height=\"476\" stroke=\"#CACED3\" stroke-width=\"3\" rx=\"10\" ry=\"10\" x=\"2\" y=\"2\"></rect><rect fill=\"$${cardColor}\" height=\"476\" rx=\"10\" ry=\"10\" width=\"736\" x=\"2\" y=\"2\" opacity=\"$${bgOpacityPercent}\"></rect><line y2=\"160\" x2=\"695\" y1=\"160\" x1=\"42.5\" stroke=\"$${textColor}\"></line><text fill=\"$${textColor}\" font-weight=\"450\" font-size=\"30\" x=\"160\" y=\"90\">$${cardTitle}</text></svg>"
  
	metadata = {
	  name = "%[3]s"
	  version = 5
	  bg_opacity_percent = 100
	  card_color = "#000000"
	  text_color = "InvalidColor"
  
	  fields = [
		{
		  type = "Issued Timestamp"
		  title = "timestamp"
		  is_visible = false
		},
	  ]
	}
  }`, acctest.CredentialsSandboxEnvironment(), resourceName, name)
}

func testAccCredentialTypeConfig_InvalidBackgroundOpacityValue(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s
resource "pingone_credential_type" "%[2]s" {
	environment_id = data.pingone_environment.credentials_test.id
	title = "%[3]s"
	description = "%[3]s Example Description"
	card_type = "DemonstrationCard"
	card_design_template = "<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 740 480\"><rect fill=\"none\" width=\"736\" height=\"476\" stroke=\"#CACED3\" stroke-width=\"3\" rx=\"10\" ry=\"10\" x=\"2\" y=\"2\"></rect><rect fill=\"$${cardColor}\" height=\"476\" rx=\"10\" ry=\"10\" width=\"736\" x=\"2\" y=\"2\" opacity=\"$${bgOpacityPercent}\"></rect><line y2=\"160\" x2=\"695\" y1=\"160\" x1=\"42.5\" stroke=\"$${textColor}\"></line><text fill=\"$${textColor}\" font-weight=\"450\" font-size=\"30\" x=\"160\" y=\"90\">$${cardTitle}</text></svg>"
  
	metadata = {
	  name = "%[3]s"
	  version = 5
	  bg_opacity_percent = 101
	  card_color = "#000000"
	  text_color = "#000000"
  
	  fields = [
		{
		  type = "Issued Timestamp"
		  title = "timestamp"
		  is_visible = false
		},
	  ]
	}
  }`, acctest.CredentialsSandboxEnvironment(), resourceName, name)
}

func testAccCredentialTypeConfig_InvalidVersion(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s
resource "pingone_credential_type" "%[2]s" {
	environment_id = data.pingone_environment.credentials_test.id
	title = "%[3]s"
	description = "%[3]s Example Description"
	card_type = "DemonstrationCard"
	card_design_template = "<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 740 480\"><rect fill=\"none\" width=\"736\" height=\"476\" stroke=\"#CACED3\" stroke-width=\"3\" rx=\"10\" ry=\"10\" x=\"2\" y=\"2\"></rect><rect fill=\"$${cardColor}\" height=\"476\" rx=\"10\" ry=\"10\" width=\"736\" x=\"2\" y=\"2\" opacity=\"$${bgOpacityPercent}\"></rect><line y2=\"160\" x2=\"695\" y1=\"160\" x1=\"42.5\" stroke=\"$${textColor}\"></line><text fill=\"$${textColor}\" font-weight=\"450\" font-size=\"30\" x=\"160\" y=\"90\">$${cardTitle}</text></svg>"
  
	metadata = {
	  name = "%[3]s"
	  version = 4
	  bg_opacity_percent = 100
	  card_color = "#000000"
	  text_color = "#000000"
  
	  fields = [
		{
		  type = "Issued Timestamp"
		  title = "timestamp"
		  is_visible = false
		},
	  ]
	}
  }`, acctest.CredentialsSandboxEnvironment(), resourceName, name)
}

func testAccCredentialTypeConfig_NoSVGTagCardDesignTemplate(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s
resource "pingone_credential_type" "%[2]s" {
	environment_id = data.pingone_environment.credentials_test.id
	title = "%[3]s"
	description = "%[3]s Example Description"
	card_type = "DemonstrationCard"
	card_design_template = "<rect fill=\"none\" width=\"736\" height=\"476\" stroke=\"#CACED3\" stroke-width=\"3\" rx=\"10\" ry=\"10\" x=\"2\" y=\"2\"></rect><rect fill=\"$${cardColor}\" height=\"476\" rx=\"10\" ry=\"10\" width=\"736\" x=\"2\" y=\"2\" opacity=\"$${bgOpacityPercent}\"></rect><line y2=\"160\" x2=\"695\" y1=\"160\" x1=\"42.5\" stroke=\"$${textColor}\"></line><text fill=\"$${textColor}\" font-weight=\"450\" font-size=\"30\" x=\"160\" y=\"90\">$${cardTitle}</text>"
  
	metadata = {
	  name = "%[3]s"
	  version = 5
	  bg_opacity_percent = 100
	  card_color = "#000000"
	  text_color = "#000000"
  
	  fields = [
		{
		  type = "Issued Timestamp"
		  title = "timestamp"
		  is_visible = false
		},
	  ]
	}
  }`, acctest.CredentialsSandboxEnvironment(), resourceName, name)
}

func testAccCredentialTypeConfig_EmptyFieldsArray(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s
resource "pingone_credential_type" "%[2]s" {
	environment_id = data.pingone_environment.credentials_test.id
	title = "%[3]s"
	description = "%[3]s Example Description"
	card_type = "DemonstrationCard"
	card_design_template = "<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 740 480\"><rect fill=\"none\" width=\"736\" height=\"476\" stroke=\"#CACED3\" stroke-width=\"3\" rx=\"10\" ry=\"10\" x=\"2\" y=\"2\"></rect><rect fill=\"$${cardColor}\" height=\"476\" rx=\"10\" ry=\"10\" width=\"736\" x=\"2\" y=\"2\" opacity=\"$${bgOpacityPercent}\"></rect><line y2=\"160\" x2=\"695\" y1=\"160\" x1=\"42.5\" stroke=\"$${textColor}\"></line><text fill=\"$${textColor}\" font-weight=\"450\" font-size=\"30\" x=\"160\" y=\"90\">$${cardTitle}</text></svg>"
  
	metadata = {
	  name = "%[3]s"
	  version = 5
	  bg_opacity_percent = 100
	  card_color = "#000000"
	  text_color = "#000000"
  
	  fields = [
	  ]
	}
  }`, acctest.CredentialsSandboxEnvironment(), resourceName, name)
}
