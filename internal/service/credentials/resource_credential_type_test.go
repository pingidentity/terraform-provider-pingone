// Copyright © 2026 Ping Identity Corporation

package credentials_test

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	acctestlegacysdk "github.com/pingidentity/terraform-provider-pingone/internal/acctest/legacysdk"
	baselegacysdk "github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/base/legacysdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/credentials"
	"github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccCredentialType_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_credential_type.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	var credentialTypeID, environmentID string

	var p1Client *client.Client
	var ctx = context.Background()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoBeta(t)
			p1Client = acctestlegacysdk.PreCheckTestClient(ctx, t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             credentials.CredentialType_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccCredentialTypeConfig_Minimal(resourceName, name),
				Check:  credentials.CredentialType_GetIDs(resourceFullName, &environmentID, &credentialTypeID),
			},
			{
				PreConfig: func() {
					credentials.CredentialType_RemovalDrift_PreConfig(ctx, p1Client.API.CredentialsAPIClient, t, environmentID, credentialTypeID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the environment
			{
				Config: testAccCredentialTypeConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check:  credentials.CredentialType_GetIDs(resourceFullName, &environmentID, &credentialTypeID),
			},
			{
				PreConfig: func() {
					baselegacysdk.Environment_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID)
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

	data, _ := os.ReadFile("../../acctest/test_assets/image/credential_background.png")
	backgroundImage := base64.StdEncoding.EncodeToString(data)

	data, _ = os.ReadFile("../../acctest/test_assets/image/credential_logo.png")
	logoImage := base64.StdEncoding.EncodeToString(data)

	cardDesignTemplate := `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 740 480">
<rect fill="none" width="736" height="476" stroke="#CACED3" stroke-width="3" rx="10" ry="10" x="2" y="2"></rect>
<rect fill="${cardColor}" height="476" rx="10" ry="10" width="736" x="2" y="2" opacity="${bgOpacityPercent}"></rect>
<image href="${backgroundImage}" opacity="${bgOpacityPercent}" height="476" rx="10" ry="10" width="736" x="2" y="2"></image>
<image href="${logoImage}" x="42" y="43" height="90px" width="90px"></image>
<line y2="160" x2="695" y1="160" x1="42.5" stroke="${textColor}"></line>
<text fill="${textColor}" font-weight="450" font-size="30" x="160" y="90">${cardTitle}</text>
<text fill="${textColor}" font-size="25" font-weight="300" x="160" y="130">${cardSubtitle}</text>
</svg>
`

	fullStep := resource.TestStep{
		Config: testAccCredentialTypeConfig_Full(resourceName, name, backgroundImage, logoImage),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "issuer_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "title", name),
			resource.TestCheckResourceAttr(resourceFullName, "description", fmt.Sprintf("%s Example Description", name)),
			resource.TestCheckResourceAttr(resourceFullName, "card_type", "VerifiedEmployee"),
			resource.TestCheckResourceAttr(resourceFullName, "card_design_template", cardDesignTemplate),
			resource.TestCheckResourceAttr(resourceFullName, "management_mode", "AUTOMATED"),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.name", name),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.description", fmt.Sprintf("%s Example Description", name)),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.version", "5"), // ensures calculated default is 5
			resource.TestCheckResourceAttr(resourceFullName, "metadata.columns", "1"),
			resource.TestMatchResourceAttr(resourceFullName, "metadata.background_image", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca)|(sg))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
			resource.TestMatchResourceAttr(resourceFullName, "metadata.logo_image", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca)|(sg))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.bg_opacity_percent", "100"),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.card_color", "#000000"),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.text_color", "#eff0f1"),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.fields.#", "8"),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.fields.3.id", "Directory Attribute -> displayName"),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.fields.3.type", "Directory Attribute"),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.fields.3.title", "displayName"),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.fields.3.attribute", "name.formatted"),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.fields.3.is_visible", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.fields.3.required", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.fields.6.id", "Directory Attribute -> id"),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.fields.6.type", "Directory Attribute"),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.fields.6.title", "id"),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.fields.6.attribute", "id"),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.fields.6.is_visible", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.fields.6.required", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.fields.7.file_support", "REFERENCE_FILE"),
			resource.TestCheckResourceAttr(resourceFullName, "revoke_on_delete", "true"),
			resource.TestMatchResourceAttr(resourceFullName, "created_at", verify.RFC3339Regexp),
			resource.TestMatchResourceAttr(resourceFullName, "updated_at", verify.RFC3339Regexp),
		),
	}

	updatedName := acctest.ResourceNameGen()
	updatedCardDesignTemplate := `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 740 480">
<rect fill="none" width="736" height="476" stroke="#CACED3" stroke-width="3" rx="10" ry="10" x="2" y="2"></rect>
<rect fill="" height="476" rx="10" ry="10" width="736" x="2" y="2" opacity=""></rect>
<line y2="160" x2="695" y1="160" x1="42.5" stroke=""></line>
<text fill="" font-weight="450" font-size="30" x="160" y="90">${cardTitle}</text>
<text font-size="25" font-weight="300" x="160" y="130">${cardSubtitle}</text>
</svg>
`

	minimalStep := resource.TestStep{
		Config: testAccCredentialTypeConfig_Minimal(resourceName, updatedName),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "issuer_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "title", updatedName),
			resource.TestCheckResourceAttr(resourceFullName, "description", fmt.Sprintf("%s Example Description", updatedName)),
			resource.TestCheckResourceAttr(resourceFullName, "card_type", "DemonstrationCard"),
			resource.TestCheckResourceAttr(resourceFullName, "card_design_template", updatedCardDesignTemplate),
			resource.TestCheckResourceAttr(resourceFullName, "management_mode", "AUTOMATED"),
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
			resource.TestCheckResourceAttr(resourceFullName, "revoke_on_delete", "false"),
			resource.TestMatchResourceAttr(resourceFullName, "created_at", verify.RFC3339Regexp),
			resource.TestMatchResourceAttr(resourceFullName, "updated_at", verify.RFC3339Regexp),
		),
	}

	updateManagementModeStep := resource.TestStep{
		Config: testAccCredentialTypeConfig_ManagedCredential(resourceName, updatedName),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "issuer_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "title", updatedName),
			resource.TestCheckResourceAttr(resourceFullName, "description", fmt.Sprintf("%s Example Description", updatedName)),
			resource.TestCheckResourceAttr(resourceFullName, "card_type", "DemonstrationCard"),
			resource.TestCheckResourceAttr(resourceFullName, "card_design_template", updatedCardDesignTemplate),
			resource.TestCheckResourceAttr(resourceFullName, "management_mode", "MANAGED"),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.name", updatedName),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.version", "5"), // ensures calculated default is 5
			resource.TestCheckResourceAttr(resourceFullName, "metadata.fields.#", "3"),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.fields.0.id", "Issued Timestamp -> timestamp"),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.fields.0.type", "Issued Timestamp"),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.fields.0.title", "timestamp"),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.fields.0.is_visible", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.fields.1.id", "Alphanumeric Text -> selfie"),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.fields.1.type", "Alphanumeric Text"),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.fields.1.title", "selfie"),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.fields.1.is_visible", "false"),
			resource.TestCheckNoResourceAttr(resourceFullName, "metadata.fields.1.value"),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.fields.2.id", "Alphanumeric Text -> other"),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.fields.2.type", "Alphanumeric Text"),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.fields.2.title", "other"),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.fields.2.is_visible", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "metadata.fields.2.value", "sample"),
			resource.TestCheckNoResourceAttr(resourceFullName, "metadata.columns"),
			resource.TestCheckNoResourceAttr(resourceFullName, "metadata.bg_opacity_percent"),
			resource.TestCheckNoResourceAttr(resourceFullName, "metadata.card_color"),
			resource.TestCheckNoResourceAttr(resourceFullName, "metadata.text_color"),
			resource.TestCheckResourceAttr(resourceFullName, "revoke_on_delete", "false"),
			resource.TestMatchResourceAttr(resourceFullName, "created_at", verify.RFC3339Regexp),
			resource.TestMatchResourceAttr(resourceFullName, "updated_at", verify.RFC3339Regexp),
		),
	}

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
			// full - new
			fullStep,
			{
				Config:  testAccCredentialTypeConfig_Full(resourceName, name, backgroundImage, logoImage),
				Destroy: true,
			},
			// minimal - new
			minimalStep,
			{
				Config:  testAccCredentialTypeConfig_Minimal(resourceName, updatedName),
				Destroy: true,
			},
			// update
			fullStep,
			updateManagementModeStep,
			fullStep,
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("resource not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
			// clear
			{
				Config:  testAccCredentialTypeConfig_ManagedCredential(resourceName, updatedName),
				Destroy: true,
			},
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
				Config:      testAccCredentialTypeConfig_InvalidTitle(resourceName, ""),
				ExpectError: regexp.MustCompile("Error: Invalid Attribute Value Length"),
				Destroy:     true,
			},
			{
				Config:      testAccCredentialTypeConfig_InvalidCardColorHexValue(resourceName, name),
				ExpectError: regexp.MustCompile("Attribute metadata.card_color expected value to contain a valid 6-digit\nhexadecimal color code"),
				Destroy:     true,
			},
			{
				Config:      testAccCredentialTypeConfig_InvalidTextColorHexValue(resourceName, name),
				ExpectError: regexp.MustCompile("Attribute metadata.text_color expected value to contain a valid 6-digit\nhexadecimal color code"),
				Destroy:     true,
			},
			{
				Config:      testAccCredentialTypeConfig_InvalidBackgroundOpacityValue(resourceName, name),
				ExpectError: regexp.MustCompile("Attribute metadata.bg_opacity_percent value must be between 0 and 100"),
				Destroy:     true,
			},
			{
				Config:      testAccCredentialTypeConfig_EmptyFieldsArray(resourceName, name),
				ExpectError: regexp.MustCompile("Attribute metadata.fields list must contain at least 1 elements"),
				Destroy:     true,
			},
			{
				Config:      testAccCredentialTypeConfig_InvalidFileSupportValue(resourceName, name),
				ExpectError: regexp.MustCompile("Error: Invalid Attribute Value Match"),
				Destroy:     true,
			},
			{
				Config:      testAccCredentialTypeConfig_InvalidManagementModeValueCombination(resourceName, name),
				ExpectError: regexp.MustCompile("Error: Invalid credential type configuration"),
				Destroy:     true,
			},
		},
	})
}

func TestAccCredentialType_CardDesignTemplate(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	name := acctest.ResourceNameGen()

	data, _ := os.ReadFile("../../acctest/test_assets/image/credential_background.png")
	backgroundImage := base64.StdEncoding.EncodeToString(data)

	data, _ = os.ReadFile("../../acctest/test_assets/image/image-logo.gif") // >50kb
	logoImage := base64.StdEncoding.EncodeToString(data)

	// escape certain messages to test
	noCardColorErrorMsg := regexp.QuoteMeta("Attribute metadata.card_color The metadata.card_color argument is defined but\nthe card_design_template does not have a ${cardColor} element.")
	noSubTitleErrorMsg := regexp.QuoteMeta("Attribute description The description argument is defined but the\ncard_design_template does not have a ${cardSubtitle} element.")
	noTextColorErrorMsg := regexp.QuoteMeta("Attribute metadata.text_color The metadata.text_color argument is defined but\nthe card_design_template does not have a ${textColor} element.")
	noBackgroundImageErrorMsg := regexp.QuoteMeta("Attribute metadata.background_image The metadata.background_image argument is\ndefined but the card_design_template does not have a ${backgroundImage}\nelement.")
	noLogoImageErrorMsg := regexp.QuoteMeta("Attribute metadata.logo_image The metadata.logo_image argument is defined but\nthe card_design_template does not have a ${logoImage} element.")

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
				Config:      testAccCredentialTypeConfig_CardDesignTemplate_NoSVG(resourceName, name),
				ExpectError: regexp.MustCompile("Attribute card_design_template expected value to contain a valid PingOne\nCredentials SVG card template."),
				Destroy:     true,
			},
			{
				Config:      testAccCredentialTypeConfig_CardDesignTemplate_NoBackgroundImage(resourceName, name, backgroundImage),
				ExpectError: regexp.MustCompile(noBackgroundImageErrorMsg),
				Destroy:     true,
			},
			{
				Config:      testAccCredentialTypeConfig_CardDesignTemplate_NoCardColor(resourceName, name),
				ExpectError: regexp.MustCompile(noCardColorErrorMsg),
				Destroy:     true,
			},
			{
				Config:      testAccCredentialTypeConfig_CardDesignTemplate_NoLogoImage(resourceName, name, logoImage),
				ExpectError: regexp.MustCompile(noLogoImageErrorMsg),
				Destroy:     true,
			},
			{
				Config:      testAccCredentialTypeConfig_CardDesignTemplate_NoSubtitle(resourceName, name),
				ExpectError: regexp.MustCompile(noSubTitleErrorMsg),
				Destroy:     true,
			},
			{
				Config:      testAccCredentialTypeConfig_CardDesignTemplate_NoTextColor(resourceName, name),
				ExpectError: regexp.MustCompile(noTextColorErrorMsg),
				Destroy:     true,
			},
		},
	})
}

func TestAccCredentialType_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_credential_type.%s", resourceName)

	name := resourceName

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
			// Configure
			{
				Config: testAccCredentialTypeConfig_Minimal(resourceName, name),
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

func testAccCredentialTypeConfig_NewEnv(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_credential_type" "%[3]s" {
  environment_id       = pingone_environment.%[2]s.id
  title                = "%[4]s"
  description          = "%[4]s"
  card_type            = "DemonstrationCard"
  card_design_template = <<-EOT
  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 740 480">
  <rect fill="none" width="736" height="476" stroke="#CACED3" stroke-width="3" rx="10" ry="10" x="2" y="2"></rect>
  <rect fill="" height="476" rx="10" ry="10" width="736" x="2" y="2" opacity=""></rect>
  <line y2="160" x2="695" y1="160" x1="42.5" stroke=""></line>
  <text fill="" font-weight="450" font-size="30" x="160" y="90">$${cardTitle}</text>
  <text font-size="25" font-weight="300" x="160" y="130">$${cardSubtitle}</text>
  </svg>
  EOT

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
}`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccCredentialTypeConfig_Full(resourceName, name, backgroundImage, logoImage string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_image" "%[2]s-background_image" {
  environment_id    = data.pingone_environment.general_test.id
  image_file_base64 = "%[4]s"
}

resource "pingone_image" "%[2]s-logo_image" {
  environment_id    = data.pingone_environment.general_test.id
  image_file_base64 = "%[5]s"
}

resource "pingone_credential_type" "%[2]s" {
  environment_id   = data.pingone_environment.general_test.id
  title            = "%[3]s"
  description      = "%[3]s Example Description"
  card_type        = "VerifiedEmployee"
  revoke_on_delete = true

  card_design_template = <<-EOT
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 740 480">
<rect fill="none" width="736" height="476" stroke="#CACED3" stroke-width="3" rx="10" ry="10" x="2" y="2"></rect>
<rect fill="$${cardColor}" height="476" rx="10" ry="10" width="736" x="2" y="2" opacity="$${bgOpacityPercent}"></rect>
<image href="$${backgroundImage}" opacity="$${bgOpacityPercent}" height="476" rx="10" ry="10" width="736" x="2" y="2"></image>
<image href="$${logoImage}" x="42" y="43" height="90px" width="90px"></image>
<line y2="160" x2="695" y1="160" x1="42.5" stroke="$${textColor}"></line>
<text fill="$${textColor}" font-weight="450" font-size="30" x="160" y="90">$${cardTitle}</text>
<text fill="$${textColor}" font-size="25" font-weight="300" x="160" y="130">$${cardSubtitle}</text>
</svg>
EOT

  metadata = {
    name               = "%[3]s"
    description        = "%[3]s Example Description"
    columns            = 1
    background_image   = pingone_image.%[2]s-background_image.uploaded_image.href
    logo_image         = pingone_image.%[2]s-logo_image.uploaded_image.href
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
        required   = true
      },
      {
        type         = "Directory Attribute"
        title        = "photo"
        attribute    = "photo"
        is_visible   = false
        file_support = "REFERENCE_FILE"
      }
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name, backgroundImage, logoImage)
}

func testAccCredentialTypeConfig_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_credential_type" "%[2]s" {
  environment_id   = data.pingone_environment.general_test.id
  title            = "%[3]s"
  description      = "%[3]s Example Description"
  card_type        = "DemonstrationCard"
  revoke_on_delete = false

  card_design_template = <<-EOT
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 740 480">
<rect fill="none" width="736" height="476" stroke="#CACED3" stroke-width="3" rx="10" ry="10" x="2" y="2"></rect>
<rect fill="" height="476" rx="10" ry="10" width="736" x="2" y="2" opacity=""></rect>
<line y2="160" x2="695" y1="160" x1="42.5" stroke=""></line>
<text fill="" font-weight="450" font-size="30" x="160" y="90">$${cardTitle}</text>
<text font-size="25" font-weight="300" x="160" y="130">$${cardSubtitle}</text>
</svg>
EOT

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

func testAccCredentialTypeConfig_InvalidTitle(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_credential_type" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  title          = "%[3]s"
  description    = "%[3]s Example Description"
  card_type      = "DemonstrationCard"

  card_design_template = <<-EOT
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 740 480">
<rect fill="none" width="736" height="476" stroke="#CACED3" stroke-width="3" rx="10" ry="10" x="2" y="2"></rect>
<rect fill="$${cardColor}" height="476" rx="10" ry="10" width="736" x="2" y="2" opacity="$${bgOpacityPercent}"></rect>
<line y2="160" x2="695" y1="160" x1="42.5" stroke="$${textColor}"></line>
<text fill="$${textColor}" font-weight="450" font-size="30" x="160" y="90">$${cardTitle}</text>
</svg>
EOT

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
  environment_id = data.pingone_environment.general_test.id
  title          = "%[3]s"
  description    = "%[3]s Example Description"
  card_type      = "DemonstrationCard"

  card_design_template = <<-EOT
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 740 480">
<rect fill="none" width="736" height="476" stroke="#CACED3" stroke-width="3" rx="10" ry="10" x="2" y="2"></rect>
<rect fill="$${cardColor}" height="476" rx="10" ry="10" width="736" x="2" y="2" opacity=""></rect>
<line y2="160" x2="695" y1="160" x1="42.5" stroke=""></line>
<text fill="$${textColor}" font-weight="450" font-size="30" x="160" y="90">$${cardTitle}</text>
<text font-size="25" font-weight="300" x="160" y="130">$${cardSubtitle}</text>
</svg>
EOT

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
  environment_id = data.pingone_environment.general_test.id
  title          = "%[3]s"
  description    = "%[3]s Example Description"
  card_type      = "DemonstrationCard"

  card_design_template = <<-EOT
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 740 480">
<rect fill="none" width="736" height="476" stroke="#CACED3" stroke-width="3" rx="10" ry="10" x="2" y="2"></rect>
<rect fill="$${cardColor}" height="476" rx="10" ry="10" width="736" x="2" y="2" opacity=""></rect>
<line y2="160" x2="695" y1="160" x1="42.5" stroke=""></line>
<text fill="$${textColor}" font-weight="450" font-size="30" x="160" y="90">$${cardTitle}</text>
<text font-size="25" font-weight="300" x="160" y="130">$${cardSubtitle}</text>
</svg>
EOT  

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
  environment_id = data.pingone_environment.general_test.id
  title          = "%[3]s"
  description    = "%[3]s Example Description"
  card_type      = "DemonstrationCard"

  card_design_template = <<-EOT
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 740 480">
<rect fill="none" width="736" height="476" stroke="#CACED3" stroke-width="3" rx="10" ry="10" x="2" y="2"></rect>
<rect fill="$${cardColor}" height="476" rx="10" ry="10" width="736" x="2" y="2" opacity=""></rect>
<line y2="160" x2="695" y1="160" x1="42.5" stroke=""></line>
<text fill="$${textColor}" font-weight="450" font-size="30" x="160" y="90">$${cardTitle}</text>
<text font-size="25" font-weight="300" x="160" y="130">$${cardSubtitle}</text>
</svg>
EOT

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

func testAccCredentialTypeConfig_EmptyFieldsArray(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s
resource "pingone_credential_type" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  title          = "%[3]s"
  description    = "%[3]s Example Description"
  card_type      = "DemonstrationCard"

  card_design_template = <<-EOT
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 740 480">
<rect fill="none" width="736" height="476" stroke="#CACED3" stroke-width="3" rx="10" ry="10" x="2" y="2"></rect>
<rect fill="$${cardColor}" height="476" rx="10" ry="10" width="736" x="2" y="2" opacity=""></rect>
<line y2="160" x2="695" y1="160" x1="42.5" stroke=""></line>
<text fill="$${textColor}" font-weight="450" font-size="30" x="160" y="90">$${cardTitle}</text>
<text font-size="25" font-weight="300" x="160" y="130">$${cardSubtitle}</text>
</svg>
EOT  

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

func testAccCredentialTypeConfig_InvalidFileSupportValue(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_credential_type" "%[3]s" {
  environment_id   = data.pingone_environment.general_test.id
  title            = "%[3]s"
  description      = "%[3]s Example Description"
  card_type        = "DemonstrationCard"
  revoke_on_delete = true

  card_design_template = <<-EOT
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 740 480">
<rect fill="none" width="736" height="476" stroke="#CACED3" stroke-width="3" rx="10" ry="10" x="2" y="2"></rect>
<rect fill="" height="476" rx="10" ry="10" width="736" x="2" y="2" opacity=""></rect>
<line y2="160" x2="695" y1="160" x1="42.5" stroke=""></line>
<text fill="" font-weight="450" font-size="30" x="160" y="90">$${cardTitle}</text>
<text font-size="25" font-weight="300" x="160" y="130">$${cardSubtitle}</text>
</svg>
EOT

  metadata = {
    name = "%[3]s"

    fields = [
      {
        type         = "Issued Timestamp"
        title        = "timestamp"
        is_visible   = false
        file_support = "REFERENCE_FILE"
      }
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccCredentialTypeConfig_InvalidManagementModeValueCombination(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_credential_type" "%[3]s" {
  environment_id   = data.pingone_environment.general_test.id
  title            = "%[3]s"
  description      = "%[3]s Example Description"
  card_type        = "DemonstrationCard"
  revoke_on_delete = true

  card_design_template = <<-EOT
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 740 480">
<rect fill="none" width="736" height="476" stroke="#CACED3" stroke-width="3" rx="10" ry="10" x="2" y="2"></rect>
<rect fill="" height="476" rx="10" ry="10" width="736" x="2" y="2" opacity=""></rect>
<line y2="160" x2="695" y1="160" x1="42.5" stroke=""></line>
<text fill="" font-weight="450" font-size="30" x="160" y="90">$${cardTitle}</text>
<text font-size="25" font-weight="300" x="160" y="130">$${cardSubtitle}</text>
</svg>
EOT

  metadata = {
    name = "%[3]s"

    fields = [
      {
        type       = "Alphanumeric Text"
        title      = "selfie"
        is_visible = false
      }
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
  card_design_template = <<-EOT
<rect fill="none" width="736" height="476" stroke="#CACED3" stroke-width="3" rx="10" ry="10" x="2" y="2"></rect>
<rect fill="$${cardColor}" height="476" rx="10" ry="10" width="736" x="2" y="2" opacity="$${bgOpacityPercent}"></rect>
<image href="$${backgroundImage}" opacity="$${bgOpacityPercent}" height="476" rx="10" ry="10" width="736" x="2" y="2"></image>
<image href="$${logoImage}" x="42" y="43" height="90px" width="90px"></image>
<line y2="160" x2="695" y1="160" x1="42.5" stroke="$${textColor}"></line>
<text fill="$${textColor}" font-weight="450" font-size="30" x="160" y="90">$${cardTitle}</text>
<text fill="$${textColor}" font-size="25" font-weight="300" x="160" y="130">$${cardSubtitle}</text>
EOT

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

resource "pingone_image" "%[2]s-background_image" {
  environment_id    = data.pingone_environment.general_test.id
  image_file_base64 = "%[4]s"
}

resource "pingone_credential_type" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  title          = "%[3]s"
  description    = "%[3]s Example Description"
  card_type      = "DemonstrationCard"

  card_design_template = <<-EOT
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 740 480">
<rect fill="none" width="736" height="476" stroke="#CACED3" stroke-width="3" rx="10" ry="10" x="2" y="2"></rect>
<rect fill="$${cardColor}" height="476" rx="10" ry="10" width="736" x="2" y="2" opacity="$${bgOpacityPercent}"></rect>
<line y2="160" x2="695" y1="160" x1="42.5" stroke="$${textColor}"></line>
<text fill="$${textColor}" font-weight="450" font-size="30" x="160" y="90">$${cardTitle}</text>
<text fill="$${textColor}" font-size="25" font-weight="300" x="160" y="130">$${cardSubtitle}</text>
</svg>
EOT  

  metadata = {
    name               = "%[3]s"
    bg_opacity_percent = 100
    card_color         = "#000000"
    text_color         = "#000000"
    background_image   = pingone_image.%[2]s-background_image.uploaded_image.href
    //background_image = "https://wtf.example.com"

    fields = [
      {
        type       = "Issued Timestamp"
        title      = "timestamp"
        is_visible = false
      },
    ]
  }
  depends_on = [pingone_image.%[2]s-background_image]
}`, acctest.GenericSandboxEnvironment(), resourceName, name, backgroundImage)
}

func testAccCredentialTypeConfig_CardDesignTemplate_NoCardColor(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s
resource "pingone_credential_type" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  title          = "%[3]s"
  description    = "%[3]s Example Description"
  card_type      = "DemonstrationCard"

  card_design_template = <<-EOT
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 740 480">
<rect fill="none" width="736" height="476" stroke="#CACED3" stroke-width="3" rx="10" ry="10" x="2" y="2"></rect>
<rect height="476" rx="10" ry="10" width="736" x="2" y="2" opacity="$${bgOpacityPercent}"></rect>
<line y2="160" x2="695" y1="160" x1="42.5" stroke="$${textColor}"></line>
<text fill="$${textColor}" font-weight="450" font-size="30" x="160" y="90">$${cardTitle}</text>
<text fill="$${textColor}" font-size="25" font-weight="300" x="160" y="130">$${cardSubtitle}</text>
</svg>
EOT  

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

resource "pingone_image" "%[2]s-logo_image" {
  environment_id    = data.pingone_environment.general_test.id
  image_file_base64 = "%[4]s"
}
resource "pingone_credential_type" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  title          = "%[3]s"
  description    = "%[3]s Example Description"
  card_type      = "DemonstrationCard"

  card_design_template = <<-EOT
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 740 480">
<rect fill="none" width="736" height="476" stroke="#CACED3" stroke-width="3" rx="10" ry="10" x="2" y="2"></rect>
<rect fill="$${cardColor}" height="476" rx="10" ry="10" width="736" x="2" y="2" opacity="$${bgOpacityPercent}"></rect>
<line y2="160" x2="695" y1="160" x1="42.5" stroke="$${textColor}"></line>
<text fill="$${textColor}" font-weight="450" font-size="30" x="160" y="90">$${cardTitle}</text>
<text fill="$${textColor}" font-size="25" font-weight="300" x="160" y="130">$${cardSubtitle}</text>
</svg>
EOT  

  metadata = {
    name               = "%[3]s"
    description        = "%[3]s Example Description"
    bg_opacity_percent = 100
    card_color         = "#000000"
    text_color         = "#000000"
    logo_image         = pingone_image.%[2]s-logo_image.uploaded_image.href # {logoImage} is missing from card_design_template

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
  environment_id = data.pingone_environment.general_test.id
  title          = "%[3]s"
  description    = "%[3]s Example Description"
  card_type      = "DemonstrationCard"

  card_design_template = <<-EOT
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 740 480">
<rect fill="none" width="736" height="476" stroke="#CACED3" stroke-width="3" rx="10" ry="10" x="2" y="2"></rect>
<rect fill="$${cardColor}" height="476" rx="10" ry="10" width="736" x="2" y="2" opacity="$${bgOpacityPercent}"></rect>
<image href="" opacity="$${bgOpacityPercent}" height="476" rx="10" ry="10" width="736" x="2" y="2"></image>
<image href="" x="42" y="43" height="90px" width="90px"></image><line y2="160" x2="695" y1="160" x1="42.5" stroke="$${textColor}"></line>
<text fill="$${textColor}" font-weight="450" font-size="30" x="160" y="90">$${cardTitle}</text>
<text fill="$${textColor}" font-size="25" font-weight="300" x="160" y="130"></text>
</svg>
EOT

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
  environment_id = data.pingone_environment.general_test.id
  title          = "%[3]s"
  description    = "%[3]s Example Description"
  card_type      = "DemonstrationCard"

  card_design_template = <<-EOT
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 740 480">
<rect fill="none" width="736" height="476" stroke="#CACED3" stroke-width="3" rx="10" ry="10" x="2" y="2"></rect>
<rect fill="$${cardColor}" height="476" rx="10" ry="10" width="736" x="2" y="2" opacity="$${bgOpacityPercent}"></rect>
<line y2="160" x2="695" y1="160" x1="42.5" stroke=""></line>
<text fill="" font-weight="450" font-size="30" x="160" y="90">$${cardTitle}</text>
<text fill="" font-size="25" font-weight="300" x="160" y="130">$${cardSubtitle}</text>
</svg>
EOT

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

func testAccCredentialTypeConfig_ManagedCredential(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_credential_type" "%[2]s" {
  environment_id   = data.pingone_environment.general_test.id
  title            = "%[3]s"
  description      = "%[3]s Example Description"
  card_type        = "DemonstrationCard"
  management_mode  = "MANAGED"
  revoke_on_delete = false

  card_design_template = <<-EOT
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 740 480">
<rect fill="none" width="736" height="476" stroke="#CACED3" stroke-width="3" rx="10" ry="10" x="2" y="2"></rect>
<rect fill="" height="476" rx="10" ry="10" width="736" x="2" y="2" opacity=""></rect>
<line y2="160" x2="695" y1="160" x1="42.5" stroke=""></line>
<text fill="" font-weight="450" font-size="30" x="160" y="90">$${cardTitle}</text>
<text font-size="25" font-weight="300" x="160" y="130">$${cardSubtitle}</text>
</svg>
EOT

  metadata = {
    name = "%[3]s"

    fields = [
      {
        type       = "Issued Timestamp"
        title      = "timestamp"
        is_visible = false
      },
      {
        type       = "Alphanumeric Text"
        title      = "selfie"
        is_visible = false
      },
      {
        type       = "Alphanumeric Text"
        title      = "other"
        is_visible = false
        value      = "sample"
      }
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
