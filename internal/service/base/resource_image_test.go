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

func TestAccImage_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_image.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	data, _ := os.ReadFile("../../acctest/test_assets/image/image-logo.png")
	image := base64.StdEncoding.EncodeToString(data)

	var imageID, environmentID string

	var p1Client *client.Client
	var ctx = context.Background()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)

			acctest.PreCheckNoBeta(t)
			p1Client = acctest.PreCheckTestClient(ctx, t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Image_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccImageConfig_Image(resourceName, image),
				Check:  base.Image_GetIDs(resourceFullName, &environmentID, &imageID),
			},
			{
				PreConfig: func() {
					base.Image_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID, imageID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the environment
			{
				Config: testAccImageConfig_NewEnv(environmentName, licenseID, resourceName, image),
				Check:  base.Image_GetIDs(resourceFullName, &environmentID, &imageID),
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

func TestAccImage_PNG(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_image.%s", resourceName)

	data, _ := os.ReadFile("../../acctest/test_assets/image/image-logo.png")
	image := base64.StdEncoding.EncodeToString(data)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Image_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccImageConfig_Image(resourceName, image),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "uploaded_image.width", "901"),
					resource.TestCheckResourceAttr(resourceFullName, "uploaded_image.height", "901"),
					resource.TestCheckResourceAttr(resourceFullName, "uploaded_image.type", "png"),
					resource.TestMatchResourceAttr(resourceFullName, "uploaded_image.href", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca)|(sg))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
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
				ImportStateVerifyIgnore: []string{
					"image_file_base64",
				},
			},
		},
	})
}

func TestAccImage_JPG(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_image.%s", resourceName)

	data, _ := os.ReadFile("../../acctest/test_assets/image/image-logo.jpg")
	image := base64.StdEncoding.EncodeToString(data)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Image_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccImageConfig_Image(resourceName, image),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "uploaded_image.width", "901"),
					resource.TestCheckResourceAttr(resourceFullName, "uploaded_image.height", "901"),
					resource.TestCheckResourceAttr(resourceFullName, "uploaded_image.type", "png"),
					resource.TestMatchResourceAttr(resourceFullName, "uploaded_image.href", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca)|(sg))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
				),
			},
		},
	})
}

func TestAccImage_GIF(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_image.%s", resourceName)

	data, _ := os.ReadFile("../../acctest/test_assets/image/image-logo.gif")
	image := base64.StdEncoding.EncodeToString(data)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Image_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccImageConfig_Image(resourceName, image),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "uploaded_image.width", "901"),
					resource.TestCheckResourceAttr(resourceFullName, "uploaded_image.height", "901"),
					resource.TestCheckResourceAttr(resourceFullName, "uploaded_image.type", "png"),
					resource.TestMatchResourceAttr(resourceFullName, "uploaded_image.href", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca)|(sg))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
				),
			},
		},
	})
}

func TestAccImage_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_image.%s", resourceName)

	data, _ := os.ReadFile("../../acctest/test_assets/image/image-logo.png")
	image := base64.StdEncoding.EncodeToString(data)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Image_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccImageConfig_Image(resourceName, image),
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

func testAccImageConfig_NewEnv(environmentName, licenseID, resourceName, image string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_image" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  image_file_base64 = "%[4]s"

}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, image)
}

func testAccImageConfig_Image(resourceName, image string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_image" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  image_file_base64 = "%[3]s"

}`, acctest.GenericSandboxEnvironment(), resourceName, image)
}
