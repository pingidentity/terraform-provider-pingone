package base_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func testAccCheckBrandingThemeDefaultDestroy(s *terraform.State) error {
	return nil
}

func TestAccBrandingThemeDefault_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_branding_theme_default.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckBrandingThemeDefaultDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBrandingThemeDefaultConfig_Full(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "branding_theme_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "default", "true"),
				),
			},
		},
	})
}

func testAccBrandingThemeDefaultConfig_Full(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_branding_theme" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name     = "%[2]s"
  template = "split"

  background_color   = "#FF00F0"
  button_text_color  = "#FF6C6C"
  heading_text_color = "#FF0005"
  card_color         = "#0FFF39"
  body_text_color    = "#8620FF"
  link_text_color    = "#8A7F06"
  button_color       = "#0CFFFB"

}

resource "pingone_branding_theme_default" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  branding_theme_id = pingone_branding_theme.%[2]s.id
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
