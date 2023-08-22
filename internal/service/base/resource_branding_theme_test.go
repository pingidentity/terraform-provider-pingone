package base_test

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func testAccCheckBrandingThemeDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_branding_theme" {
			continue
		}

		_, rEnv, err := apiClient.EnvironmentsApi.ReadOneEnvironment(ctx, rs.Primary.Attributes["environment_id"]).Execute()

		if err != nil {

			if rEnv == nil {
				return fmt.Errorf("Response object does not exist and no error detected")
			}

			if rEnv.StatusCode == 404 {
				continue
			}

			return err
		}

		body, r, err := apiClient.BrandingThemesApi.ReadOneBrandingTheme(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.ID).Execute()

		if err != nil {

			if r == nil {
				return fmt.Errorf("Response object does not exist and no error detected")
			}

			if r.StatusCode == 404 {
				continue
			}

			tflog.Error(ctx, fmt.Sprintf("Error: %v", body))
			return err
		}

		return fmt.Errorf("PingOne Branding Theme Instance %s still exists", rs.Primary.ID)
	}

	return nil
}

func testAccGetBrandingThemeIDs(resourceName string, environmentID, resourceID *string) resource.TestCheckFunc {
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

func TestAccBrandingTheme_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_branding_theme.%s", resourceName)

	name := resourceName

	var resourceID, environmentID string

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckBrandingThemeDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccBrandingThemeConfig_Minimal(resourceName, name),
				Check:  testAccGetBrandingThemeIDs(resourceFullName, &environmentID, &resourceID),
			},
			// Replan after removal preconfig
			{
				PreConfig: func() {
					var ctx = context.Background()
					p1Client, err := acctest.TestClient(ctx)

					if err != nil {
						t.Fatalf("Failed to get API client: %v", err)
					}

					apiClient := p1Client.API.ManagementAPIClient

					if environmentID == "" || resourceID == "" {
						t.Fatalf("One of environment ID or resource ID cannot be determined. Environment ID: %s, Resource ID: %s", environmentID, resourceID)
					}

					_, err = apiClient.BrandingThemesApi.DeleteBrandingTheme(ctx, environmentID, resourceID).Execute()
					if err != nil {
						t.Fatalf("Failed to delete branding theme: %v", err)
					}
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
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckBrandingThemeDestroy,
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

	logoData, _ := ioutil.ReadFile("../../acctest/test_assets/image/image-logo.gif")
	logo := base64.StdEncoding.EncodeToString(logoData)

	backgroundData, _ := ioutil.ReadFile("../../acctest/test_assets/image/image-background.jpg")
	background := base64.StdEncoding.EncodeToString(backgroundData)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckBrandingThemeDestroy,
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
					resource.TestCheckResourceAttr(resourceFullName, "logo.#", "1"),
					resource.TestMatchResourceAttr(resourceFullName, "logo.0.id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "logo.0.href", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
					resource.TestCheckResourceAttr(resourceFullName, "background_image.#", "1"),
					resource.TestMatchResourceAttr(resourceFullName, "background_image.0.id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "background_image.0.href", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
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
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckBrandingThemeDestroy,
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
					resource.TestCheckResourceAttr(resourceFullName, "logo.#", "0"),
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
		},
	})
}

func TestAccBrandingTheme_Change(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_branding_theme.%s", resourceName)

	name := resourceName

	logoData, _ := ioutil.ReadFile("../../acctest/test_assets/image/image-logo.gif")
	logo := base64.StdEncoding.EncodeToString(logoData)

	backgroundData, _ := ioutil.ReadFile("../../acctest/test_assets/image/image-background.jpg")
	background := base64.StdEncoding.EncodeToString(backgroundData)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckBrandingThemeDestroy,
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
					resource.TestCheckResourceAttr(resourceFullName, "logo.#", "1"),
					resource.TestMatchResourceAttr(resourceFullName, "logo.0.id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "logo.0.href", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
					resource.TestCheckResourceAttr(resourceFullName, "background_image.#", "1"),
					resource.TestMatchResourceAttr(resourceFullName, "background_image.0.id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "background_image.0.href", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
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
					resource.TestCheckResourceAttr(resourceFullName, "logo.#", "0"),
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
					resource.TestCheckResourceAttr(resourceFullName, "logo.#", "1"),
					resource.TestMatchResourceAttr(resourceFullName, "logo.0.id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "logo.0.href", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
					resource.TestCheckResourceAttr(resourceFullName, "background_image.#", "1"),
					resource.TestMatchResourceAttr(resourceFullName, "background_image.0.id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "background_image.0.href", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
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
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckBrandingThemeDestroy,
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
				ExpectError:  regexp.MustCompile(`Invalid import ID specified \(".*"\).  The ID should be in the format "environment_id/branding_theme_id".`),
			},
			{
				ResourceName:  resourceFullName,
				ImportStateId: "/",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`Invalid import ID specified \(".*"\).  The ID should be in the format "environment_id/branding_theme_id".`),
			},
			{
				ResourceName:  resourceFullName,
				ImportStateId: "badformat/badformat",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`Invalid import ID specified \(".*"\).  The ID should be in the format "environment_id/branding_theme_id".`),
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

  logo {
    id   = pingone_image.%[2]s-logo.id
    href = pingone_image.%[2]s-logo.uploaded_image[0].href
  }

  background_image {
    id   = pingone_image.%[2]s-background.id
    href = pingone_image.%[2]s-background.uploaded_image[0].href
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
