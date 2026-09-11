// Copyright © 2026 Ping Identity Corporation

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
	acctestlegacysdk "github.com/pingidentity/terraform-provider-pingone/internal/acctest/legacysdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/base"
	baselegacysdk "github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/base/legacysdk"
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
			acctest.PreCheckNoBeta(t)
			p1Client = acctestlegacysdk.PreCheckTestClient(ctx, t)
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
					baselegacysdk.Environment_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID)
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
			acctest.PreCheckNoBeta(t)
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
			acctest.PreCheckNoBeta(t)
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
					resource.TestCheckResourceAttr(resourceFullName, "header", "<h1>Welcome to PingOne</h1>"),
					resource.TestCheckResourceAttr(resourceFullName, "application_background_color", "#F8F8F8"),
					resource.TestCheckResourceAttr(resourceFullName, "header_background_color", "#4A4A4A"),
					resource.TestCheckResourceAttr(resourceFullName, "focus_rectangle_color", "#D033FF"),
					resource.TestCheckResourceAttr(resourceFullName, "title_text_color", "#4A4A4A"),
					resource.TestCheckResourceAttr(resourceFullName, "sub_title_text_color", "#4A4A4A"),
					resource.TestCheckResourceAttr(resourceFullName, "body_text_size", "15px"),
					resource.TestCheckResourceAttr(resourceFullName, "body_text_weight", "400"),
					resource.TestCheckResourceAttr(resourceFullName, "button_corner_radius", "2px"),
					resource.TestCheckResourceAttr(resourceFullName, "button_border_width", "1px"),
					resource.TestCheckResourceAttr(resourceFullName, "global_font", "\"Helvetica Neue\", Helvetica, sans-serif"),
					resource.TestCheckResourceAttr(resourceFullName, "logo_height", "56px"),
					resource.TestCheckResourceAttr(resourceFullName, "card_horizontal_alignment", "left"),
					resource.TestCheckResourceAttr(resourceFullName, "card_vertical_alignment", "center"),
					resource.TestCheckResourceAttr(resourceFullName, "input_label_position", "float"),
					resource.TestCheckResourceAttr(resourceFullName, "button_border_color", "#112233"),
					resource.TestCheckResourceAttr(resourceFullName, "button_hover_state_border_color", "#223344"),
					resource.TestCheckResourceAttr(resourceFullName, "button_hover_state_fill_color", "#334455"),
					resource.TestCheckResourceAttr(resourceFullName, "button_hover_state_text_color", "#556677"),
					resource.TestCheckResourceAttr(resourceFullName, "card_border_color", "#CACED3"),
					resource.TestCheckResourceAttr(resourceFullName, "card_border_width", "2px"),
					resource.TestCheckResourceAttr(resourceFullName, "card_corner_radius", "4px"),
					resource.TestCheckResourceAttr(resourceFullName, "card_shadow", "0px 2px 0px 0px rgba(0,0,0,0.1)"),
					resource.TestCheckResourceAttr(resourceFullName, "input_border_width", "1px"),
					resource.TestCheckResourceAttr(resourceFullName, "input_box_border_color", "#667788"),
					resource.TestCheckResourceAttr(resourceFullName, "input_corner_radius", "5px"),
					resource.TestCheckResourceAttr(resourceFullName, "input_label_text_color", "#798087"),
					resource.TestCheckResourceAttr(resourceFullName, "input_label_text_size", "14px"),
					resource.TestCheckResourceAttr(resourceFullName, "input_label_text_weight", "500"),
					resource.TestCheckResourceAttr(resourceFullName, "input_value_text_color", "#798087"),
					resource.TestCheckResourceAttr(resourceFullName, "input_value_text_size", "13px"),
					resource.TestCheckResourceAttr(resourceFullName, "input_value_text_weight", "600"),
					resource.TestCheckResourceAttr(resourceFullName, "link_text_hover_color", "#007CBA"),
					resource.TestCheckResourceAttr(resourceFullName, "link_text_size", "15px"),
					resource.TestCheckResourceAttr(resourceFullName, "link_text_weight", "400"),
					resource.TestCheckResourceAttr(resourceFullName, "sub_title_text_size", "15px"),
					resource.TestCheckResourceAttr(resourceFullName, "sub_title_text_weight", "400"),
					resource.TestCheckResourceAttr(resourceFullName, "title_text_size", "20px"),
					resource.TestCheckResourceAttr(resourceFullName, "title_text_weight", "600"),
					resource.TestCheckResourceAttr(resourceFullName, "button_text_size", "16px"),
					resource.TestCheckResourceAttr(resourceFullName, "button_text_weight", "400"),
				),
			},
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
		},
	})
}

// TestAccBrandingTheme_FullLocalized exercises header_localized / footer_localized in
// isolation: the PingOne API cannot remove these objects once configured (see
// scratch/pingone-branding-theme-api-echo-issues.md), so the config is applied once and
// never transitioned away from.
func TestAccBrandingTheme_FullLocalized(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_branding_theme.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.BrandingTheme_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBrandingThemeConfig_FullLocalized(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "header_localized.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "header_localized.default_language", "en"),
					resource.TestCheckResourceAttr(resourceFullName, "header_localized.content_type", "HTML"),
					resource.TestCheckResourceAttr(resourceFullName, "header_localized.content.%", "2"),
					resource.TestCheckResourceAttr(resourceFullName, "header_localized.content.en.input_text", "<h1>Welcome</h1>"),
					resource.TestCheckResourceAttr(resourceFullName, "header_localized.content.fr.input_text", "<h1>Bienvenue</h1>"),
					resource.TestCheckResourceAttr(resourceFullName, "footer_localized.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "footer_localized.default_language", "en"),
					resource.TestCheckResourceAttr(resourceFullName, "footer_localized.content_type", "HTML"),
					resource.TestCheckResourceAttr(resourceFullName, "footer_localized.content.%", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "footer_localized.content.en.input_text", "<p>Copyright</p>"),
				),
			},
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
			acctest.PreCheckNoBeta(t)
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
					resource.TestCheckNoResourceAttr(resourceFullName, "header"),
					resource.TestCheckNoResourceAttr(resourceFullName, "header_localized"),
					resource.TestCheckNoResourceAttr(resourceFullName, "footer_localized"),
					resource.TestCheckNoResourceAttr(resourceFullName, "title_text_color"),
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
			acctest.PreCheckNoBeta(t)
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
					resource.TestCheckNoResourceAttr(resourceFullName, "header"),
					resource.TestCheckResourceAttr(resourceFullName, "heading_text_color", "#FF0005"),
					resource.TestCheckResourceAttr(resourceFullName, "link_text_color", "#8A7F06"),
				),
			},
			// NOTE: no localized-to-unlocalized transition step here — the PingOne API cannot
			// remove headerLocalized/footerLocalized once set (see scratch/pingone-branding-theme-api-echo-issues.md),
			// so a config that removes the block would churn indefinitely.
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

func TestAccBrandingTheme_ValidationChecks(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	name := acctest.ResourceNameGen()

	backgroundData, _ := os.ReadFile("../../acctest/test_assets/image/image-background.jpg")
	background := base64.StdEncoding.EncodeToString(backgroundData)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.BrandingTheme_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccBrandingThemeConfig_BackgroundColorAndBackgroundImageConflict(resourceName, name, background),
				ExpectError: regexp.MustCompile(`Error: Invalid Attribute Combination`),
				Destroy:     true,
			},
			{
				Config:      testAccBrandingThemeConfig_InvalidLocalizedContentType(resourceName, name),
				ExpectError: regexp.MustCompile(`value must be one of`),
				Destroy:     true,
			},
			{
				Config:      testAccBrandingThemeConfig_InvalidTitleTextColor(resourceName, name),
				ExpectError: regexp.MustCompile(`Value must be a valid hex color code.`),
				Destroy:     true,
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
			acctest.PreCheckNoBeta(t)
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

}`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
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

  use_default_background = false

  button_text_color  = "#FF6C6C"
  heading_text_color = "#FF0005"
  card_color         = "#0FFF39"
  body_text_color    = "#8620FF"
  link_text_color    = "#8A7F06"
  button_color       = "#0CFFFB"

  footer_text = "What do you call a can opener that doesn't work? A can't opener."

  application_background_color = "#F8F8F8"
  header_background_color      = "#4A4A4A"
  focus_rectangle_color        = "#D033FF"
  title_text_color             = "#4A4A4A"
  sub_title_text_color         = "#4A4A4A"
  body_text_size               = "15px"
  body_text_weight             = "400"
  button_corner_radius         = "2px"
  button_border_width          = "1px"
  global_font                  = "\"Helvetica Neue\", Helvetica, sans-serif"
  logo_height                  = "56px"
  card_horizontal_alignment    = "left"
  card_vertical_alignment      = "center"
  input_label_position         = "float"
  header                       = "<h1>Welcome to PingOne</h1>"

  button_border_color             = "#112233"
  button_hover_state_border_color = "#223344"
  button_hover_state_fill_color   = "#334455"
  button_hover_state_text_color   = "#556677"
  card_border_color               = "#CACED3"
  card_border_width               = "2px"
  card_corner_radius              = "4px"
  card_shadow                     = "0px 2px 0px 0px rgba(0,0,0,0.1)"
  input_border_width              = "1px"
  input_box_border_color          = "#667788"
  input_corner_radius             = "5px"
  input_label_text_color          = "#798087"
  input_label_text_size           = "14px"
  input_label_text_weight         = "500"
  input_value_text_color          = "#798087"
  input_value_text_size           = "13px"
  input_value_text_weight         = "600"
  link_text_hover_color           = "#007CBA"
  link_text_size                  = "15px"
  link_text_weight                = "400"
  sub_title_text_size             = "15px"
  sub_title_text_weight           = "400"
  title_text_size                 = "20px"
  title_text_weight               = "600"
  button_text_size                = "16px"
  button_text_weight              = "400"

  // NOTE: header_localized / footer_localized are exercised in TestAccBrandingTheme_FullLocalized.
  // The PingOne API cannot remove these objects once set (see scratch/pingone-branding-theme-api-echo-issues.md),
  // so the full→minimal transitions below must not cross the localized boundary.

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

func testAccBrandingThemeConfig_FullLocalized(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_branding_theme" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name     = "%[3]s"
  template = "split"

  button_text_color  = "#FF6C6C"
  heading_text_color = "#FF0005"
  card_color         = "#0FFF39"
  body_text_color    = "#8620FF"
  link_text_color    = "#8A7F06"
  button_color       = "#0CFFFB"

  header_localized = {
    enabled          = true
    default_language = "en"
    content_type     = "HTML"
    content = {
      en = { input_text = "<h1>Welcome</h1>" }
      fr = { input_text = "<h1>Bienvenue</h1>" }
    }
  }

  footer_localized = {
    enabled          = true
    default_language = "en"
    content_type     = "HTML"
    content = {
      en = { input_text = "<p>Copyright</p>" }
    }
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccBrandingThemeConfig_InvalidLocalizedContentType(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_branding_theme" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name     = "%[3]s"
  template = "split"

  button_text_color  = "#FF6C6C"
  heading_text_color = "#FF0005"
  card_color         = "#0FFF39"
  body_text_color    = "#8620FF"
  link_text_color    = "#8A7F06"
  button_color       = "#0CFFFB"

  header_localized = {
    enabled          = true
    default_language = "en"
    content_type     = "Markdown"
    content = {
      en = { input_text = "<h1>Welcome</h1>" }
    }
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccBrandingThemeConfig_InvalidTitleTextColor(resourceName, name string) string {
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
  title_text_color   = "red"

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccBrandingThemeConfig_BackgroundColorAndBackgroundImageConflict(resourceName, name, background string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_image" "%[2]s-background" {
  environment_id = data.pingone_environment.general_test.id

  image_file_base64 = "%[4]s"
}

resource "pingone_branding_theme" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name     = "%[3]s"
  template = "split"

  background_color = "#FF00F0"
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

}`, acctest.GenericSandboxEnvironment(), resourceName, name, background)
}
