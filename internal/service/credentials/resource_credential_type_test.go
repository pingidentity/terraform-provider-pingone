package credentials_test

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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

	mgmtApiClient := p1Client.API.ManagementAPIClient

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

func testAccGetMFAPolicyIDs(resourceName string, environmentID, resourceID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource Not found: %s", resourceName)
		}

		*resourceID = rs.Primary.ID
		*environmentID = rs.Primary.Attributes["environment_id"]

		return nil
	}
}

func TestAccMFAPolicy_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_policy.%s", resourceName)

	name := resourceName

	var resourceID, environmentID string

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckMFAPolicyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccMFAPolicyConfig_FullSMS(resourceName, name),
				Check:  testAccGetMFAPolicyIDs(resourceFullName, &environmentID, &resourceID),
			},
			// Replan after removal preconfig
			{
				PreConfig: func() {
					var ctx = context.Background()
					p1Client, err := acctest.TestClient(ctx)

					if err != nil {
						t.Fatalf("Failed to get API client: %v", err)
					}

					apiClient := p1Client.API.MFAAPIClient

					if environmentID == "" || resourceID == "" {
						t.Fatalf("One of environment ID or resource ID cannot be determined. Environment ID: %s, Resource ID: %s", environmentID, resourceID)
					}

					_, err = apiClient.DeviceAuthenticationPolicyApi.DeleteDeviceAuthenticationPolicy(ctx, environmentID, resourceID).Execute()
					if err != nil {
						t.Fatalf("Failed to delete MFA Policy: %v", err)
					}
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccCredentialType_NewEnv(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_credential_type.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := acctest.ResourceNameGen()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCredentialTypeDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCredentialTypeConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "title", name),
				),
			},
		},
	})
}

func TestAccCredentialType_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_credential_type.%s", resourceName)

	name := acctest.ResourceNameGen()

	data, _ := os.ReadFile("../../acctest/test_assets/image/credential_background_base64.png")
	backgroundImage := "data:image/png;base64," + base64.StdEncoding.EncodeToString(data)

	data, _ = os.ReadFile("../../acctest/test_assets/image/credential_logo_base64.png")
	logoImage := "data:image/png;base64," + base64.StdEncoding.EncodeToString(data)

	// Note: If template is defined directly in HCL, the ${variableName} variables are escaped with $$, such as $${variableName}.
	// The HCL in the test case has the escaped variable name.  The test value does not ensuring it is saved to, and returned from, state properly.
	// Future: Move the design template to a test asset file.
	cardDesignTemplate := "<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 740 480\"><rect fill=\"none\" width=\"736\" height=\"476\" stroke=\"#CACED3\" stroke-width=\"3\" rx=\"10\" ry=\"10\" x=\"2\" y=\"2\"></rect><rect fill=\"${cardColor}\" height=\"476\" rx=\"10\" ry=\"10\" width=\"736\" x=\"2\" y=\"2\" opacity=\"${bgOpacityPercent}\"></rect><image href=\"${backgroundImage}\" opacity=\"${bgOpacityPercent}\" height=\"476\" rx=\"10\" ry=\"10\" width=\"736\" x=\"2\" y=\"2\"></image><image href=\"${logoImage}\" x=\"42\" y=\"43\" height=\"90px\" width=\"90px\"></image><line y2=\"160\" x2=\"695\" y1=\"160\" x1=\"42.5\" stroke=\"${textColor}\"></line><text fill=\"${textColor}\" font-weight=\"450\" font-size=\"30\" x=\"160\" y=\"90\">${cardTitle}</text><text fill=\"${textColor}\" font-size=\"25\" font-weight=\"300\" x=\"160\" y=\"130\">${cardSubtitle}</text></svg>"

	fullStep := resource.TestStep{
		Config: testAccCredentialTypeConfig_Full(resourceName, name, backgroundImage, logoImage),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
			resource.TestCheckResourceAttr(resourceFullName, "title", name),
			resource.TestCheckResourceAttr(resourceFullName, "description", fmt.Sprintf("%s Example Description", name)),
			resource.TestCheckResourceAttr(resourceFullName, "card_type", "VerifiedEmployee"),
			resource.TestCheckResourceAttr(resourceFullName, "card_design_template", cardDesignTemplate),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.name", name),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.description", fmt.Sprintf("%s Example Description", name)),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.version", "5"), // ensures calculated default is 5
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

	// Note: If template is defined directly in HCL, the ${variableName} variables are escaped with $$, such as $${variableName}.
	// The HCL in the test case has the escaped variable name.  The test value does not ensuring it is saved to, and returned from, state properly.
	// Future: Move the design template to a test asset file.
	updatedCardDesignTemplate := "<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 740 480\"><rect fill=\"none\" width=\"736\" height=\"476\" stroke=\"#CACED3\" stroke-width=\"3\" rx=\"10\" ry=\"10\" x=\"2\" y=\"2\"></rect><rect fill=\"\" height=\"476\" rx=\"10\" ry=\"10\" width=\"736\" x=\"2\" y=\"2\" opacity=\"\"></rect><line y2=\"160\" x2=\"695\" y1=\"160\" x1=\"42.5\" stroke=\"\"></line><text fill=\"\" font-weight=\"450\" font-size=\"30\" x=\"160\" y=\"90\">${cardTitle}</text><text font-size=\"25\" font-weight=\"300\" x=\"160\" y=\"130\">${cardSubtitle}</text></svg>"

	minimalStep := resource.TestStep{
		Config: testAccCredentialTypeConfig_Minimal(resourceName, updatedName),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
			resource.TestCheckResourceAttr(resourceFullName, "title", updatedName),
			resource.TestCheckResourceAttr(resourceFullName, "description", fmt.Sprintf("%s Example Description", updatedName)),
			resource.TestCheckResourceAttr(resourceFullName, "card_type", "DemonstrationCard"),
			resource.TestCheckResourceAttr(resourceFullName, "card_design_template", updatedCardDesignTemplate),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.name", updatedName),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.version", "5"), // ensures calculated default is 5
			resource.TestCheckResourceAttr(resourceFullName, "metadata.fields.#", "1"),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.fields.0.id", "Issued Timestamp -> timestamp"),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.fields.0.type", "Issued Timestamp"),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.fields.0.title", "timestamp"),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.fields.0.is_visible", "false"),

			resource.TestCheckNoResourceAttr(resourceFullName, "metadata.columns"),
			resource.TestCheckNoResourceAttr(resourceFullName, "metadata.bg_opacity_percent"),
			resource.TestCheckNoResourceAttr(resourceFullName, "metadata.card_color"),
			resource.TestCheckNoResourceAttr(resourceFullName, "metadata.text_color"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             nil, //testAccCheckCredentialTypeDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// full - new
			fullStep,
			{
				Config: testAccCredentialTypeConfig_Full(resourceName, name, backgroundImage, logoImage),
				//Destroy: true,
			},
			// minimal - new
			minimalStep,
			{
				Config:  testAccCredentialTypeConfig_Minimal(resourceName, updatedName),
				Destroy: true,
			},
			// update
			fullStep,
			minimalStep,
			fullStep,
			// clear
			{
				Config:  testAccCredentialTypeConfig_Minimal(resourceName, updatedName),
				Destroy: true,
			},
			{
				Config:  testAccCredentialTypeConfig_Full(resourceName, name, backgroundImage, logoImage),
				Destroy: true,
			},
		},
	})
}

func TestAccCredentialType_MetaData(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	name := acctest.ResourceNameGen()

	data, _ := os.ReadFile("../../acctest/test_assets/image/image-background.jpg") // >90kb
	backgroundImage := "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(data)

	data, _ = os.ReadFile("../../acctest/test_assets/image/image-logo.gif") // >50kb
	logoImage := "data:image/gif;base64," + base64.StdEncoding.EncodeToString(data)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCredentialTypeDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccCredentialTypeConfig_InvalidTitle(resourceName, ""),
				ExpectError: regexp.MustCompile("Error: Invalid Attribute Value Length"),
				Destroy:     true,
			},
			{
				Config:      testAccCredentialTypeConfig_InvalidCardColorHexValue(resourceName, name),
				ExpectError: regexp.MustCompile("Error: Invalid Attribute Value Match"),
				Destroy:     true,
			},
			{
				Config:      testAccCredentialTypeConfig_InvalidTextColorHexValue(resourceName, name),
				ExpectError: regexp.MustCompile("Error: Invalid Attribute Value Match"),
				Destroy:     true,
			},
			{
				Config:      testAccCredentialTypeConfig_InvalidBackgroundOpacityValue(resourceName, name),
				ExpectError: regexp.MustCompile("Error: Invalid Attribute Value"),
				Destroy:     true,
			},
			{
				Config:      testAccCredentialTypeConfig_InvalidVersion(resourceName, name),
				ExpectError: regexp.MustCompile("Error: Invalid Configuration for Read-Only Attribute"),
				Destroy:     true,
			},
			{
				Config:      testAccCredentialTypeConfig_EmptyFieldsArray(resourceName, name),
				ExpectError: regexp.MustCompile("Error: Invalid Attribute Value"),
				Destroy:     true,
			},
			{
				Config:      testAccCredentialTypeConfig_ImageSizeExceeded(resourceName, name, backgroundImage, logoImage),
				ExpectError: regexp.MustCompile("Error: Invalid Attribute Value Length"),
				Destroy:     true,
			},
		},
	})
}

func TestAccCredentialType_CardDesignTemplate(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	name := acctest.ResourceNameGen()

	data, _ := os.ReadFile("../../acctest/test_assets/image/credential_background_base64.png")
	backgroundImage := "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(data)

	data, _ = os.ReadFile("../../acctest/test_assets/image/image-logo.gif") // >50kb
	logoImage := "data:image/gif;base64," + base64.StdEncoding.EncodeToString(data)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCredentialTypeDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccCredentialTypeConfig_CardDesignTemplate_NoSVG(resourceName, name),
				ExpectError: regexp.MustCompile("Error: Invalid Attribute Value Match"),
				Destroy:     true,
			},
			{
				Config:      testAccCredentialTypeConfig_CardDesignTemplate_NoBackgroundImage(resourceName, name, backgroundImage),
				ExpectError: regexp.MustCompile("Error: Invalid Attribute Value Match"),
				Destroy:     true,
			},
			{
				Config:      testAccCredentialTypeConfig_CardDesignTemplate_NoCardColor(resourceName, name),
				ExpectError: regexp.MustCompile("Error: Invalid Attribute Value Match"),
				Destroy:     true,
			},
			{
				Config:      testAccCredentialTypeConfig_CardDesignTemplate_NoLogoImage(resourceName, name, logoImage),
				ExpectError: regexp.MustCompile("Error: Invalid Attribute Value Match"),
				Destroy:     true,
			},
			{
				Config:      testAccCredentialTypeConfig_CardDesignTemplate_NoSubtitle(resourceName, name),
				ExpectError: regexp.MustCompile("Error: Invalid Attribute Value Match"),
				Destroy:     true,
			},
			{
				Config:      testAccCredentialTypeConfig_CardDesignTemplate_NoTextColor(resourceName, name),
				ExpectError: regexp.MustCompile("Error: Invalid Attribute Value Match"),
				Destroy:     true,
			},
		},
	})
}

func testAccCredentialTypeConfig_NewEnv(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_credential_type" "%[3]s" {
  environment_id       = pingone_environment.%[2]s.id
  title                = "%[4]s"
  description          = "%[4]s"
  card_type            = "DemonstrationCard"
  card_design_template = "<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 740 480\"><rect fill=\"none\" width=\"736\" height=\"476\" stroke=\"#CACED3\" stroke-width=\"3\" rx=\"10\" ry=\"10\" x=\"2\" y=\"2\"></rect><rect fill=\"\" height=\"476\" rx=\"10\" ry=\"10\" width=\"736\" x=\"2\" y=\"2\" opacity=\"\"></rect><line y2=\"160\" x2=\"695\" y1=\"160\" x1=\"42.5\" stroke=\"\"></line><text fill=\"\" font-weight=\"450\" font-size=\"30\" x=\"160\" y=\"90\">$${cardTitle}</text><text font-size=\"25\" font-weight=\"300\" x=\"160\" y=\"130\">$${cardSubtitle}</text></svg>"
  metadata = {
    name = "%[4]s"

    fields = [
      {
        type       = "Issued Timestamp"
        title      = "timestamp"
        is_visible = false
      }
    ]
  }

  depends_on = [pingone_environment.%[2]s]
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccCredentialTypeConfig_Full(resourceName, name, backgroundImage, logoImage string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_credential_type" "%[2]s" {
  environment_id       = data.pingone_environment.general_test.id
  title                = "%[3]s"
  description          = "%[3]s Example Description"
  card_type            = "VerifiedEmployee"
  card_design_template = "<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 740 480\"><rect fill=\"none\" width=\"736\" height=\"476\" stroke=\"#CACED3\" stroke-width=\"3\" rx=\"10\" ry=\"10\" x=\"2\" y=\"2\"></rect><rect fill=\"$${cardColor}\" height=\"476\" rx=\"10\" ry=\"10\" width=\"736\" x=\"2\" y=\"2\" opacity=\"$${bgOpacityPercent}\"></rect><image href=\"$${backgroundImage}\" opacity=\"$${bgOpacityPercent}\" height=\"476\" rx=\"10\" ry=\"10\" width=\"736\" x=\"2\" y=\"2\"></image><image href=\"$${logoImage}\" x=\"42\" y=\"43\" height=\"90px\" width=\"90px\"></image><line y2=\"160\" x2=\"695\" y1=\"160\" x1=\"42.5\" stroke=\"$${textColor}\"></line><text fill=\"$${textColor}\" font-weight=\"450\" font-size=\"30\" x=\"160\" y=\"90\">$${cardTitle}</text><text fill=\"$${textColor}\" font-size=\"25\" font-weight=\"300\" x=\"160\" y=\"130\">$${cardSubtitle}</text></svg>"

  metadata = {
    name               = "%[3]s"
    description        = "%[3]s Example Description"
    columns            = 1
    background_image   = "%[4]s"
    logo_image         = "%[5]s"
    bg_opacity_percent = 100
    card_color         = "#000000"
    text_color         = "#eff0f1"

    fields = [
      {
        type       = "Directory Attribute"
        title      = "givenName"
        attribute  = "name.given"
        is_visible = false
      },
      {
        type       = "Directory Attribute"
        title      = "surname"
        attribute  = "name.family"
        is_visible = false
      },
      {
        type       = "Directory Attribute"
        title      = "jobTitle"
        attribute  = "title"
        is_visible = false
      },
      {
        type       = "Directory Attribute"
        title      = "displayName"
        attribute  = "name.formatted"
        is_visible = false
      },
      {
        type       = "Directory Attribute"
        title      = "mail"
        attribute  = "email"
        is_visible = false
      },
      {
        type       = "Directory Attribute"
        title      = "preferredLanguage"
        attribute  = "preferredLanguage"
        is_visible = false
      },
      {
        type       = "Directory Attribute"
        title      = "id"
        attribute  = "id"
        is_visible = false
      }
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name, backgroundImage, logoImage)
}

func testAccCredentialTypeConfig_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_credential_type" "%[2]s" {
  environment_id       = data.pingone_environment.general_test.id
  title                = "%[3]s"
  description          = "%[3]s Example Description"
  card_type            = "DemonstrationCard"
  card_design_template = "<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 740 480\"><rect fill=\"none\" width=\"736\" height=\"476\" stroke=\"#CACED3\" stroke-width=\"3\" rx=\"10\" ry=\"10\" x=\"2\" y=\"2\"></rect><rect fill=\"\" height=\"476\" rx=\"10\" ry=\"10\" width=\"736\" x=\"2\" y=\"2\" opacity=\"\"></rect><line y2=\"160\" x2=\"695\" y1=\"160\" x1=\"42.5\" stroke=\"\"></line><text fill=\"\" font-weight=\"450\" font-size=\"30\" x=\"160\" y=\"90\">$${cardTitle}</text><text font-size=\"25\" font-weight=\"300\" x=\"160\" y=\"130\">$${cardSubtitle}</text></svg>"

  metadata = {
    name = "%[3]s"

    fields = [
      {
        type       = "Issued Timestamp"
        title      = "timestamp"
        is_visible = false
      }
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccCredentialTypeConfig_ImageSizeExceeded(resourceName, name, backgroundImage, logoImage string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_credential_type" "%[2]s" {
  environment_id       = data.pingone_environment.general_test.id
  title                = "%[3]s"
  description          = "%[3]s Example Description"
  card_type            = "VerifiedEmployee"
  card_design_template = "<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 740 480\"><rect fill=\"none\" width=\"736\" height=\"476\" stroke=\"#CACED3\" stroke-width=\"3\" rx=\"10\" ry=\"10\" x=\"2\" y=\"2\"></rect><rect fill=\"$${cardColor}\" height=\"476\" rx=\"10\" ry=\"10\" width=\"736\" x=\"2\" y=\"2\" opacity=\"$${bgOpacityPercent}\"></rect><image href=\"$${backgroundImage}\" opacity=\"$${bgOpacityPercent}\" height=\"476\" rx=\"10\" ry=\"10\" width=\"736\" x=\"2\" y=\"2\"></image><image href=\"$${logoImage}\" x=\"42\" y=\"43\" height=\"90px\" width=\"90px\"></image><line y2=\"160\" x2=\"695\" y1=\"160\" x1=\"42.5\" stroke=\"$${textColor}\"></line><text fill=\"$${textColor}\" font-weight=\"450\" font-size=\"30\" x=\"160\" y=\"90\">$${cardTitle}</text><text fill=\"$${textColor}\" font-size=\"25\" font-weight=\"300\" x=\"160\" y=\"130\">$${cardSubtitle}</text></svg>"

  metadata = {
    name               = "%[3]s"
    description        = "%[3]s Example Description"
    background_image   = "%[4]s"
    logo_image         = "%[5]s"
    bg_opacity_percent = 100
    card_color         = "#000000"
    text_color         = "#eff0f1"

    fields = [
      {
        type       = "Directory Attribute"
        title      = "givenName"
        attribute  = "name.given"
        is_visible = false
      },
      {
        type       = "Directory Attribute"
        title      = "surname"
        attribute  = "name.family"
        is_visible = false
      },
      {
        type       = "Directory Attribute"
        title      = "jobTitle"
        attribute  = "title"
        is_visible = false
      },
      {
        type       = "Directory Attribute"
        title      = "displayName"
        attribute  = "name.formatted"
        is_visible = false
      },
      {
        type       = "Directory Attribute"
        title      = "mail"
        attribute  = "email"
        is_visible = false
      },
      {
        type       = "Directory Attribute"
        title      = "preferredLanguage"
        attribute  = "preferredLanguage"
        is_visible = false
      },
      {
        type       = "Directory Attribute"
        title      = "id"
        attribute  = "id"
        is_visible = false
      }
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name, backgroundImage, logoImage)
}

func testAccCredentialTypeConfig_InvalidTitle(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_credential_type" "%[2]s" {
  environment_id       = data.pingone_environment.general_test.id
  title                = "%[3]s"
  description          = "%[3]s Example Description"
  card_type            = "DemonstrationCard"
  card_design_template = "<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 740 480\"><rect fill=\"none\" width=\"736\" height=\"476\" stroke=\"#CACED3\" stroke-width=\"3\" rx=\"10\" ry=\"10\" x=\"2\" y=\"2\"></rect><rect fill=\"$${cardColor}\" height=\"476\" rx=\"10\" ry=\"10\" width=\"736\" x=\"2\" y=\"2\" opacity=\"$${bgOpacityPercent}\"></rect><line y2=\"160\" x2=\"695\" y1=\"160\" x1=\"42.5\" stroke=\"$${textColor}\"></line><text fill=\"$${textColor}\" font-weight=\"450\" font-size=\"30\" x=\"160\" y=\"90\">$${cardTitle}</text></svg>"

  metadata = {
    name               = "%[3]s"
    bg_opacity_percent = 100
    card_color         = "#000000"
    text_color         = "#eff0f1"

    fields = [
      {
        type       = "Issued Timestamp"
        title      = "timestamp"
        is_visible = false
      },
      {
        type       = "Directory Attribute"
        title      = "surname"
        attribute  = "name.family"
        is_visible = false
      },
      {
        type       = "Alphanumeric Text"
        title      = "Company"
        value      = "Example"
        is_visible = false
      }
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccCredentialTypeConfig_InvalidCardColorHexValue(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_credential_type" "%[2]s" {
  environment_id       = data.pingone_environment.general_test.id
  title                = "%[3]s"
  description          = "%[3]s Example Description"
  card_type            = "DemonstrationCard"
  card_design_template = "<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 740 480\"><rect fill=\"none\" width=\"736\" height=\"476\" stroke=\"#CACED3\" stroke-width=\"3\" rx=\"10\" ry=\"10\" x=\"2\" y=\"2\"></rect><rect fill=\"$${cardColor}\" height=\"476\" rx=\"10\" ry=\"10\" width=\"736\" x=\"2\" y=\"2\" opacity=\"$${bgOpacityPercent}\"></rect><line y2=\"160\" x2=\"695\" y1=\"160\" x1=\"42.5\" stroke=\"$${textColor}\"></line><text fill=\"$${textColor}\" font-weight=\"450\" font-size=\"30\" x=\"160\" y=\"90\">$${cardTitle}</text></svg>"

  metadata = {
    name               = "%[3]s"
    bg_opacity_percent = 100
    card_color         = "InvalidColor"
    text_color         = "#eff0f1"

    fields = [
      {
        type       = "Issued Timestamp"
        title      = "timestamp"
        is_visible = false
      },
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccCredentialTypeConfig_InvalidTextColorHexValue(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s
resource "pingone_credential_type" "%[2]s" {
  environment_id       = data.pingone_environment.general_test.id
  title                = "%[3]s"
  description          = "%[3]s Example Description"
  card_type            = "DemonstrationCard"
  card_design_template = "<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 740 480\"><rect fill=\"none\" width=\"736\" height=\"476\" stroke=\"#CACED3\" stroke-width=\"3\" rx=\"10\" ry=\"10\" x=\"2\" y=\"2\"></rect><rect fill=\"$${cardColor}\" height=\"476\" rx=\"10\" ry=\"10\" width=\"736\" x=\"2\" y=\"2\" opacity=\"$${bgOpacityPercent}\"></rect><line y2=\"160\" x2=\"695\" y1=\"160\" x1=\"42.5\" stroke=\"$${textColor}\"></line><text fill=\"$${textColor}\" font-weight=\"450\" font-size=\"30\" x=\"160\" y=\"90\">$${cardTitle}</text></svg>"

  metadata = {
    name               = "%[3]s"
    bg_opacity_percent = 100
    card_color         = "#000000"
    text_color         = "InvalidColor"

    fields = [
      {
        type       = "Issued Timestamp"
        title      = "timestamp"
        is_visible = false
      },
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccCredentialTypeConfig_InvalidBackgroundOpacityValue(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s
resource "pingone_credential_type" "%[2]s" {
  environment_id       = data.pingone_environment.general_test.id
  title                = "%[3]s"
  description          = "%[3]s Example Description"
  card_type            = "DemonstrationCard"
  card_design_template = "<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 740 480\"><rect fill=\"none\" width=\"736\" height=\"476\" stroke=\"#CACED3\" stroke-width=\"3\" rx=\"10\" ry=\"10\" x=\"2\" y=\"2\"></rect><rect fill=\"$${cardColor}\" height=\"476\" rx=\"10\" ry=\"10\" width=\"736\" x=\"2\" y=\"2\" opacity=\"$${bgOpacityPercent}\"></rect><line y2=\"160\" x2=\"695\" y1=\"160\" x1=\"42.5\" stroke=\"$${textColor}\"></line><text fill=\"$${textColor}\" font-weight=\"450\" font-size=\"30\" x=\"160\" y=\"90\">$${cardTitle}</text></svg>"

  metadata = {
    name               = "%[3]s"
    bg_opacity_percent = 101
    card_color         = "#000000"
    text_color         = "#000000"

    fields = [
      {
        type       = "Issued Timestamp"
        title      = "timestamp"
        is_visible = false
      },
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccCredentialTypeConfig_InvalidVersion(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s
resource "pingone_credential_type" "%[2]s" {
  environment_id       = data.pingone_environment.general_test.id
  title                = "%[3]s"
  description          = "%[3]s Example Description"
  card_type            = "DemonstrationCard"
  card_design_template = "<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 740 480\"><rect fill=\"none\" width=\"736\" height=\"476\" stroke=\"#CACED3\" stroke-width=\"3\" rx=\"10\" ry=\"10\" x=\"2\" y=\"2\"></rect><rect fill=\"$${cardColor}\" height=\"476\" rx=\"10\" ry=\"10\" width=\"736\" x=\"2\" y=\"2\" opacity=\"$${bgOpacityPercent}\"></rect><line y2=\"160\" x2=\"695\" y1=\"160\" x1=\"42.5\" stroke=\"$${textColor}\"></line><text fill=\"$${textColor}\" font-weight=\"450\" font-size=\"30\" x=\"160\" y=\"90\">$${cardTitle}</text></svg>"

  metadata = {
    name               = "%[3]s"
    version            = 4
    bg_opacity_percent = 100
    card_color         = "#000000"
    text_color         = "#000000"

    fields = [
      {
        type       = "Issued Timestamp"
        title      = "timestamp"
        is_visible = false
      },
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccCredentialTypeConfig_EmptyFieldsArray(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s
resource "pingone_credential_type" "%[2]s" {
  environment_id       = data.pingone_environment.general_test.id
  title                = "%[3]s"
  description          = "%[3]s Example Description"
  card_type            = "DemonstrationCard"
  card_design_template = "<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 740 480\"><rect fill=\"none\" width=\"736\" height=\"476\" stroke=\"#CACED3\" stroke-width=\"3\" rx=\"10\" ry=\"10\" x=\"2\" y=\"2\"></rect><rect fill=\"$${cardColor}\" height=\"476\" rx=\"10\" ry=\"10\" width=\"736\" x=\"2\" y=\"2\" opacity=\"$${bgOpacityPercent}\"></rect><line y2=\"160\" x2=\"695\" y1=\"160\" x1=\"42.5\" stroke=\"$${textColor}\"></line><text fill=\"$${textColor}\" font-weight=\"450\" font-size=\"30\" x=\"160\" y=\"90\">$${cardTitle}</text></svg>"

  metadata = {
    name               = "%[3]s"
    bg_opacity_percent = 100
    card_color         = "#000000"
    text_color         = "#000000"

    fields = [
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccCredentialTypeConfig_CardDesignTemplate_NoSVG(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s
resource "pingone_credential_type" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  title          = "%[3]s"
  description    = "%[3]s Example Description"
  card_type      = "DemonstrationCard"

  # missing svg tags
  card_design_template = "<rect fill=\"none\" width=\"736\" height=\"476\" stroke=\"#CACED3\" stroke-width=\"3\" rx=\"10\" ry=\"10\" x=\"2\" y=\"2\"></rect><rect fill=\"$${cardColor}\" height=\"476\" rx=\"10\" ry=\"10\" width=\"736\" x=\"2\" y=\"2\" opacity=\"$${bgOpacityPercent}\"></rect><line y2=\"160\" x2=\"695\" y1=\"160\" x1=\"42.5\" stroke=\"$${textColor}\"></line><text fill=\"$${textColor}\" font-weight=\"450\" font-size=\"30\" x=\"160\" y=\"90\">$${cardTitle}</text>"

  metadata = {
    name               = "%[3]s"
    bg_opacity_percent = 100
    card_color         = "#000000"
    text_color         = "#000000"

    fields = [
      {
        type       = "Issued Timestamp"
        title      = "timestamp"
        is_visible = false
      },
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccCredentialTypeConfig_CardDesignTemplate_NoBackgroundImage(resourceName, name, backgroundImage string) string {
	return fmt.Sprintf(`
	%[1]s
resource "pingone_credential_type" "%[2]s" {
  environment_id       = data.pingone_environment.general_test.id
  title                = "%[3]s"
  description          = "%[3]s Example Description"
  card_type            = "DemonstrationCard"
  card_design_template = "<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 740 480\"><rect fill=\"none\" width=\"736\" height=\"476\" stroke=\"#CACED3\" stroke-width=\"3\" rx=\"10\" ry=\"10\" x=\"2\" y=\"2\"></rect><rect fill=\"$${cardColor}\" height=\"476\" rx=\"10\" ry=\"10\" width=\"736\" x=\"2\" y=\"2\" opacity=\"$${bgOpacityPercent}\"></rect><line y2=\"160\" x2=\"695\" y1=\"160\" x1=\"42.5\" stroke=\"$${textColor}\"></line><text fill=\"$${textColor}\" font-weight=\"450\" font-size=\"30\" x=\"160\" y=\"90\">$${cardTitle}</text></svg>"

  metadata = {
    name               = "%[3]s"
    bg_opacity_percent = 100
    card_color         = "#000000"
    text_color         = "#000000"
    background_image   = "%[4]s" # {backgroundImage} is missing from card_design_template

    fields = [
      {
        type       = "Issued Timestamp"
        title      = "timestamp"
        is_visible = false
      },
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name, backgroundImage)
}

func testAccCredentialTypeConfig_CardDesignTemplate_NoCardColor(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s
resource "pingone_credential_type" "%[2]s" {
  environment_id       = data.pingone_environment.general_test.id
  title                = "%[3]s"
  description          = "%[3]s Example Description"
  card_type            = "DemonstrationCard"
  card_design_template = "<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 740 480\"><rect fill=\"none\" width=\"736\" height=\"476\" stroke=\"#CACED3\" stroke-width=\"3\" rx=\"10\" ry=\"10\" x=\"2\" y=\"2\"></rect><rect height=\"476\" rx=\"10\" ry=\"10\" width=\"736\" x=\"2\" y=\"2\" opacity=\"$${bgOpacityPercent}\"></rect><line y2=\"160\" x2=\"695\" y1=\"160\" x1=\"42.5\" stroke=\"$${textColor}\"></line><text fill=\"$${textColor}\" font-weight=\"450\" font-size=\"30\" x=\"160\" y=\"90\">$${cardTitle}</text></svg>"

  metadata = {
    name               = "%[3]s"
    bg_opacity_percent = 100
    card_color         = "#000000" # {cardColor} is missing from card_design_template
    text_color         = "#000000"

    fields = [
      {
        type       = "Issued Timestamp"
        title      = "timestamp"
        is_visible = false
      },
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccCredentialTypeConfig_CardDesignTemplate_NoLogoImage(resourceName, name, logoImage string) string {
	return fmt.Sprintf(`
	%[1]s
resource "pingone_credential_type" "%[2]s" {
  environment_id       = data.pingone_environment.general_test.id
  title                = "%[3]s"
  description          = "%[3]s Example Description"
  card_type            = "DemonstrationCard"
  card_design_template = "<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 740 480\"><rect fill=\"none\" width=\"736\" height=\"476\" stroke=\"#CACED3\" stroke-width=\"3\" rx=\"10\" ry=\"10\" x=\"2\" y=\"2\"></rect><rect fill=\"$${cardColor}\" height=\"476\" rx=\"10\" ry=\"10\" width=\"736\" x=\"2\" y=\"2\" opacity=\"$${bgOpacityPercent}\"></rect><image href=\"\" opacity=\"$${bgOpacityPercent}\" height=\"476\" rx=\"10\" ry=\"10\" width=\"736\" x=\"2\" y=\"2\"></image><image href=\"\" x=\"42\" y=\"43\" height=\"90px\" width=\"90px\"></image><line y2=\"160\" x2=\"695\" y1=\"160\" x1=\"42.5\" stroke=\"$${textColor}\"></line><text fill=\"$${textColor}\" font-weight=\"450\" font-size=\"30\" x=\"160\" y=\"90\">$${cardTitle}</text><text fill=\"$${textColor}\" font-size=\"25\" font-weight=\"300\" x=\"160\" y=\"130\">$${cardSubtitle}</text></svg>"

  metadata = {
    name               = "%[3]s"
    description        = "%[3]s Example Description"
    bg_opacity_percent = 100
    card_color         = "#000000"
    text_color         = "#000000"
    logo_image         = "%[4]s" # {logoImage} is missing from card_design_template

    fields = [
      {
        type       = "Issued Timestamp"
        title      = "timestamp"
        is_visible = false
      },
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name, logoImage)
}

func testAccCredentialTypeConfig_CardDesignTemplate_NoSubtitle(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s
resource "pingone_credential_type" "%[2]s" {
  environment_id       = data.pingone_environment.general_test.id
  title                = "%[3]s"
  description          = "%[3]s Example Description"
  card_type            = "DemonstrationCard"
  card_design_template = "<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 740 480\"><rect fill=\"none\" width=\"736\" height=\"476\" stroke=\"#CACED3\" stroke-width=\"3\" rx=\"10\" ry=\"10\" x=\"2\" y=\"2\"></rect><rect fill=\"$${cardColor}\" height=\"476\" rx=\"10\" ry=\"10\" width=\"736\" x=\"2\" y=\"2\" opacity=\"$${bgOpacityPercent}\"></rect><image href=\"\" opacity=\"$${bgOpacityPercent}\" height=\"476\" rx=\"10\" ry=\"10\" width=\"736\" x=\"2\" y=\"2\"></image><image href=\"\" x=\"42\" y=\"43\" height=\"90px\" width=\"90px\"></image><line y2=\"160\" x2=\"695\" y1=\"160\" x1=\"42.5\" stroke=\"$${textColor}\"></line><text fill=\"$${textColor}\" font-weight=\"450\" font-size=\"30\" x=\"160\" y=\"90\">$${cardTitle}</text><text fill=\"$${textColor}\" font-size=\"25\" font-weight=\"300\" x=\"160\" y=\"130\"></text></svg>"

  metadata = {
    name               = "%[3]s"
    description        = "%[3]s Example Description" # {subTitle} missing from template - description maps to card subtitle value
    bg_opacity_percent = 100
    card_color         = "#000000"
    text_color         = "#000000"

    fields = [
      {
        type       = "Issued Timestamp"
        title      = "timestamp"
        is_visible = false
      },
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccCredentialTypeConfig_CardDesignTemplate_NoTextColor(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s
resource "pingone_credential_type" "%[2]s" {
  environment_id       = data.pingone_environment.general_test.id
  title                = "%[3]s"
  description          = "%[3]s Example Description"
  card_type            = "DemonstrationCard"
  card_design_template = "<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 740 480\"><rect fill=\"none\" width=\"736\" height=\"476\" stroke=\"#CACED3\" stroke-width=\"3\" rx=\"10\" ry=\"10\" x=\"2\" y=\"2\"></rect><rect fill=\"$${cardColor}\" height=\"476\" rx=\"10\" ry=\"10\" width=\"736\" x=\"2\" y=\"2\" opacity=\"$${bgOpacityPercent}\"></rect><line y2=\"160\" x2=\"695\" y1=\"160\" x1=\"42.5\" stroke=\"\"></line><text fill=\"\" font-weight=\"450\" font-size=\"30\" x=\"160\" y=\"90\">$${cardTitle}</text></svg>"

  metadata = {
    name               = "%[3]s"
    bg_opacity_percent = 100
    card_color         = "#000000"
    text_color         = "#000000" # {textColor} is missing from card_design_template

    fields = [
      {
        type       = "Issued Timestamp"
        title      = "timestamp"
        is_visible = false
      },
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
