package base_test

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
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
)

func testAccCheckImageDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_image" {
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

		body, r, err := apiClient.ImagesApi.ReadImage(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.ID).Execute()

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

		return fmt.Errorf("PingOne Image Instance %s still exists", rs.Primary.ID)
	}

	return nil
}

func TestAccImage_PNG(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_image.%s", resourceName)

	data, _ := ioutil.ReadFile("../../acctest/test_assets/image/image-logo.png")
	image := base64.StdEncoding.EncodeToString(data)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironmentAndPKCS7(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckImageDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccImageConfig_Image(resourceName, image),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "uploaded_image.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "uploaded_image.0.width", "901"),
					resource.TestCheckResourceAttr(resourceFullName, "uploaded_image.0.height", "901"),
					resource.TestCheckResourceAttr(resourceFullName, "uploaded_image.0.type", "png"),
					resource.TestMatchResourceAttr(resourceFullName, "uploaded_image.0.href", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
				),
			},
		},
	})
}

func TestAccImage_JPG(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_image.%s", resourceName)

	data, _ := ioutil.ReadFile("../../acctest/test_assets/image/image-logo.jpg")
	image := base64.StdEncoding.EncodeToString(data)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironmentAndPKCS7(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckImageDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccImageConfig_Image(resourceName, image),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "uploaded_image.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "uploaded_image.0.width", "901"),
					resource.TestCheckResourceAttr(resourceFullName, "uploaded_image.0.height", "901"),
					resource.TestCheckResourceAttr(resourceFullName, "uploaded_image.0.type", "png"),
					resource.TestMatchResourceAttr(resourceFullName, "uploaded_image.0.href", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
				),
			},
		},
	})
}

func TestAccImage_GIF(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_image.%s", resourceName)

	data, _ := ioutil.ReadFile("../../acctest/test_assets/image/image-logo.gif")
	image := base64.StdEncoding.EncodeToString(data)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironmentAndPKCS7(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckImageDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccImageConfig_Image(resourceName, image),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "uploaded_image.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "uploaded_image.0.width", "901"),
					resource.TestCheckResourceAttr(resourceFullName, "uploaded_image.0.height", "901"),
					resource.TestCheckResourceAttr(resourceFullName, "uploaded_image.0.type", "png"),
					resource.TestMatchResourceAttr(resourceFullName, "uploaded_image.0.href", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
				),
			},
		},
	})
}

func testAccImageConfig_Image(resourceName, image string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_image" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  image_file_base64 = "%[3]s"

}`, acctest.GenericSandboxEnvironment(), resourceName, image)
}
