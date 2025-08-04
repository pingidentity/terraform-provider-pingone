// Copyright © 2025 Ping Identity Corporation

package base_test

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
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/base"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccBrandingTheme_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_branding_theme.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	var brandingThemeID, environmentID string

	var p1Client *client.Client
	var ctx = context.Background()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)

			p1Client = acctest.PreCheckTestClient(ctx, t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.BrandingTheme_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccBrandingThemeConfig_Minimal(resourceName, name),
				Check:  base.BrandingTheme_GetIDs(resourceFullName, &environmentID, &brandingThemeID),
			},
			{
				PreConfig: func() {
					base.BrandingTheme_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID, brandingThemeID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the environment
			{
				Config: testAccBrandingThemeConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check:  base.BrandingTheme_GetIDs(resourceFullName, &environmentID, &brandingThemeID),
			},
			{
				PreConfig: func() {
					base.Environment_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccBrandingTheme_NewEnv(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_branding_theme.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.BrandingTheme_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBrandingThemeConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
				),
			},
		},
	})
}

func TestAccBrandingTheme_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_branding_theme.%s", resourceName)

	name := resourceName

	logoData, _ := os.ReadFile("../../acctest/test_assets/image/image-logo.gif")
	logo := base64.StdEncoding.EncodeToString(logoData)

	backgroundData, _ := os.ReadFile("../../acctest/test_assets/image/image-background.jpg")
	background := base64.StdEncoding.EncodeToString(backgroundData)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.BrandingTheme_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBrandingThemeConfig_Full(resourceName, name, logo, background),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "template", "split"),
					resource.TestCheckResourceAttr(resourceFullName, "default", "false"),
					resource.TestMatchResourceAttr(resourceFullName, "logo.id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "logo.href", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca)|(sg))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
					resource.TestMatchResourceAttr(resourceFullName, "background_image.id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "background_image.href", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca)|(sg))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
					resource.TestCheckNoResourceAttr(resourceFullName, "background_color"),
					resource.TestCheckResourceAttr(resourceFullName, "use_default_background", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "body_text_color", "#8620FF"),
					resource.TestCheckResourceAttr(resourceFullName, "button_color", "#0CFFFB"),
					resource.TestCheckResourceAttr(resourceFullName, "button_text_color", "#FF6C6C"),
					resource.TestCheckResourceAttr(resourceFullName, "card_color", "#0FFF39"),
					resource.TestCheckResourceAttr(resourceFullName, "footer_text", "What do you call a can opener that doesn't work? A can't opener."),
					resource.TestCheckResourceAttr(resourceFullName, "heading_text_color", "#FF0005"),
					resource.TestCheckResourceAttr(resourceFullName, "link_text_color", "#8A7F06"),
				),
			},
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBrandingTheme_Minimal(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_branding_theme.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.BrandingTheme_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBrandingThemeConfig_Minimal(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "template", "split"),
					resource.TestCheckResourceAttr(resourceFullName, "default", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "logo.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "background_image.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "background_color", "#FF00F0"),
					resource.TestCheckResourceAttr(resourceFullName, "use_default_background", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "body_text_color", "#8620FF"),
					resource.TestCheckResourceAttr(resourceFullName, "button_color", "#0CFFFB"),
					resource.TestCheckResourceAttr(resourceFullName, "button_text_color", "#FF6C6C"),
					resource.TestCheckResourceAttr(resourceFullName, "card_color", "#0FFF39"),
					resource.TestCheckNoResourceAttr(resourceFullName, "footer_text"),
					resource.TestCheckResourceAttr(resourceFullName, "heading_text_color", "#FF0005"),
					resource.TestCheckResourceAttr(resourceFullName, "link_text_color", "#8A7F06"),
				),
			},
		},
	})
}

func TestAccBrandingTheme_Change(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_branding_theme.%s", resourceName)

	name := resourceName

	logoData, _ := os.ReadFile("../../acctest/test_assets/image/image-logo.gif")
	logo := base64.StdEncoding.EncodeToString(logoData)

	backgroundData, _ := os.ReadFile("../../acctest/test_assets/image/image-background.jpg")
	background := base64.StdEncoding.EncodeToString(backgroundData)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.BrandingTheme_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBrandingThemeConfig_Full(resourceName, name, logo, background),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "template", "split"),
					resource.TestCheckResourceAttr(resourceFullName, "default", "false"),
					resource.TestMatchResourceAttr(resourceFullName, "logo.id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "logo.href", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca)|(sg))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
					resource.TestMatchResourceAttr(resourceFullName, "background_image.id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "background_image.href", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca)|(sg))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
					resource.TestCheckNoResourceAttr(resourceFullName, "background_color"),
					resource.TestCheckResourceAttr(resourceFullName, "use_default_background", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "body_text_color", "#8620FF"),
					resource.TestCheckResourceAttr(resourceFullName, "button_color", "#0CFFFB"),
					resource.TestCheckResourceAttr(resourceFullName, "button_text_color", "#FF6C6C"),
					resource.TestCheckResourceAttr(resourceFullName, "card_color", "#0FFF39"),
					resource.TestCheckResourceAttr(resourceFullName, "footer_text", "What do you call a can opener that doesn't work? A can't opener."),
					resource.TestCheckResourceAttr(resourceFullName, "heading_text_color", "#FF0005"),
					resource.TestCheckResourceAttr(resourceFullName, "link_text_color", "#8A7F06"),
				),
			},
			{
				Config: testAccBrandingThemeConfig_Minimal(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "template", "split"),
					resource.TestCheckResourceAttr(resourceFullName, "default", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "logo.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "background_image.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "background_color", "#FF00F0"),
					resource.TestCheckResourceAttr(resourceFullName, "use_default_background", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "body_text_color", "#8620FF"),
					resource.TestCheckResourceAttr(resourceFullName, "button_color", "#0CFFFB"),
					resource.TestCheckResourceAttr(resourceFullName, "button_text_color", "#FF6C6C"),
					resource.TestCheckResourceAttr(resourceFullName, "card_color", "#0FFF39"),
					resource.TestCheckNoResourceAttr(resourceFullName, "footer_text"),
					resource.TestCheckResourceAttr(resourceFullName, "heading_text_color", "#FF0005"),
					resource.TestCheckResourceAttr(resourceFullName, "link_text_color", "#8A7F06"),
				),
			},
			{
				Config: testAccBrandingThemeConfig_Full(resourceName, name, logo, background),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "template", "split"),
					resource.TestCheckResourceAttr(resourceFullName, "default", "false"),
					resource.TestMatchResourceAttr(resourceFullName, "logo.id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "logo.href", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca)|(sg))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
					resource.TestMatchResourceAttr(resourceFullName, "background_image.id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "background_image.href", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca)|(sg))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
					resource.TestCheckNoResourceAttr(resourceFullName, "background_color"),
					resource.TestCheckResourceAttr(resourceFullName, "use_default_background", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "body_text_color", "#8620FF"),
					resource.TestCheckResourceAttr(resourceFullName, "button_color", "#0CFFFB"),
					resource.TestCheckResourceAttr(resourceFullName, "button_text_color", "#FF6C6C"),
					resource.TestCheckResourceAttr(resourceFullName, "card_color", "#0FFF39"),
					resource.TestCheckResourceAttr(resourceFullName, "footer_text", "What do you call a can opener that doesn't work? A can't opener."),
					resource.TestCheckResourceAttr(resourceFullName, "heading_text_color", "#FF0005"),
					resource.TestCheckResourceAttr(resourceFullName, "link_text_color", "#8A7F06"),
				),
			},
		},
	})
}

func TestAccBrandingTheme_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_branding_theme.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.BrandingTheme_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccBrandingThemeConfig_Minimal(resourceName, name),
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

func testAccBrandingThemeConfig_NewEnv(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_branding_theme" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  name     = "%[4]s"
  template = "split"

  background_color   = "#FF00F0"
  button_text_color  = "#FF6C6C"
  heading_text_color = "#FF0005"
  card_color         = "#0FFF39"
  body_text_color    = "#8620FF"
  link_text_color    = "#8A7F06"
  button_color       = "#0CFFFB"

}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccBrandingThemeConfig_Full(resourceName, name, logo, background string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_image" "%[2]s-logo" {
  environment_id = data.pingone_environment.general_test.id

  image_file_base64 = "%[4]s"
}

resource "pingone_image" "%[2]s-background" {
  environment_id = data.pingone_environment.general_test.id

  image_file_base64 = "%[5]s"
}

resource "pingone_branding_theme" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name     = "%[3]s"
  template = "split"

  logo = {
    id   = pingone_image.%[2]s-logo.id
    href = pingone_image.%[2]s-logo.uploaded_image.href
  }

  background_image = {
    id   = pingone_image.%[2]s-background.id
    href = pingone_image.%[2]s-background.uploaded_image.href
  }

  button_text_color  = "#FF6C6C"
  heading_text_color = "#FF0005"
  card_color         = "#0FFF39"
  body_text_color    = "#8620FF"
  link_text_color    = "#8A7F06"
  button_color       = "#0CFFFB"

  footer_text = "What do you call a can opener that doesn't work? A can't opener."

}`, acctest.GenericSandboxEnvironment(), resourceName, name, logo, background)
}

func testAccBrandingThemeConfig_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_branding_theme" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name     = "%[3]s"
  template = "split"

  background_color   = "#FF00F0"
  button_text_color  = "#FF6C6C"
  heading_text_color = "#FF0005"
  card_color         = "#0FFF39"
  body_text_color    = "#8620FF"
  link_text_color    = "#8A7F06"
  button_color       = "#0CFFFB"

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
