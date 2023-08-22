package sso_test

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

func testAccCheckIdentityProviderDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_identity_provider" {
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

		body, r, err := apiClient.IdentityProvidersApi.ReadOneIdentityProvider(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.ID).Execute()

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

		return fmt.Errorf("PingOne Identity Provider Instance %s still exists", rs.Primary.ID)
	}

	return nil
}

func testAccGetIdentityProviderIDs(resourceName string, environmentID, resourceID *string) resource.TestCheckFunc {
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

func TestAccIdentityProvider_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_identity_provider.%s", resourceName)

	name := resourceName

	var resourceID, environmentID string

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckIdentityProviderDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccIdentityProviderConfig_Minimal(resourceName, name),
				Check:  testAccGetIdentityProviderIDs(resourceFullName, &environmentID, &resourceID),
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

					_, err = apiClient.IdentityProvidersApi.DeleteIdentityProvider(ctx, environmentID, resourceID).Execute()
					if err != nil {
						t.Fatalf("Failed to delete identity provider: %v", err)
					}
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccIdentityProvider_NewEnv(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_identity_provider.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckIdentityProviderDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProviderConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
				),
			},
		},
	})
}

func TestAccIdentityProvider_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_identity_provider.%s", resourceName)

	name := resourceName

	data, _ := ioutil.ReadFile("../../acctest/test_assets/image/image-logo.gif")
	image := base64.StdEncoding.EncodeToString(data)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckIdentityProviderDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProviderConfig_Full(resourceName, name, image),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", "My test identity provider"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestMatchResourceAttr(resourceFullName, "registration_population_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "login_button_icon.#", "1"),
					resource.TestMatchResourceAttr(resourceFullName, "login_button_icon.0.id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "login_button_icon.0.href", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
					resource.TestCheckResourceAttr(resourceFullName, "icon.#", "1"),
					resource.TestMatchResourceAttr(resourceFullName, "icon.0.id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "icon.0.href", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
				),
			},
		},
	})
}

func TestAccIdentityProvider_Minimal(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_identity_provider.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckIdentityProviderDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProviderConfig_Minimal(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", ""),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "registration_population_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "login_button_icon.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "icon.#", "0"),
				),
			},
		},
	})
}

func TestAccIdentityProvider_Change(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_identity_provider.%s", resourceName)

	name := resourceName

	data, _ := ioutil.ReadFile("../../acctest/test_assets/image/image-logo.gif")
	image := base64.StdEncoding.EncodeToString(data)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckIdentityProviderDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProviderConfig_Full(resourceName, fmt.Sprintf("%s 1", name), image),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", fmt.Sprintf("%s 1", name)),
					resource.TestCheckResourceAttr(resourceFullName, "description", "My test identity provider"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestMatchResourceAttr(resourceFullName, "registration_population_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "login_button_icon.#", "1"),
					resource.TestMatchResourceAttr(resourceFullName, "login_button_icon.0.id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "login_button_icon.0.href", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
					resource.TestCheckResourceAttr(resourceFullName, "icon.#", "1"),
					resource.TestMatchResourceAttr(resourceFullName, "icon.0.id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "icon.0.href", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
				),
			},
			{
				Config: testAccIdentityProviderConfig_Minimal(resourceName, fmt.Sprintf("%s 2", name)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", fmt.Sprintf("%s 2", name)),
					resource.TestCheckResourceAttr(resourceFullName, "description", ""),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "registration_population_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "login_button_icon.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "icon.#", "0"),
				),
			},
			{
				Config: testAccIdentityProviderConfig_Full(resourceName, fmt.Sprintf("%s 1", name), image),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", fmt.Sprintf("%s 1", name)),
					resource.TestCheckResourceAttr(resourceFullName, "description", "My test identity provider"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestMatchResourceAttr(resourceFullName, "registration_population_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "login_button_icon.#", "1"),
					resource.TestMatchResourceAttr(resourceFullName, "login_button_icon.0.id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "login_button_icon.0.href", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
					resource.TestCheckResourceAttr(resourceFullName, "icon.#", "1"),
					resource.TestMatchResourceAttr(resourceFullName, "icon.0.id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "icon.0.href", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
				),
			},
		},
	})
}

func TestAccIdentityProvider_Facebook(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_identity_provider.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckIdentityProviderDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProviderConfig_Facebook1(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "facebook.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "facebook.0.app_id", "dummyappid1"),
					resource.TestCheckResourceAttr(resourceFullName, "facebook.0.app_secret", "dummyappsecret1"),
					// resource.TestMatchResourceAttr(resourceFullName, "facebook.0.callback_url", regexp.MustCompile(`^https:\/\/auth\.pingone\.(?:eu|com|asia|ca)\/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}\/rp\/callback\/facebook$`)),
					resource.TestCheckResourceAttr(resourceFullName, "google.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "github.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.#", "0"),
				),
			},
			{
				Config: testAccIdentityProviderConfig_Facebook2(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "facebook.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "facebook.0.app_id", "dummyappid2"),
					resource.TestCheckResourceAttr(resourceFullName, "facebook.0.app_secret", "dummyappsecret2"),
					// resource.TestMatchResourceAttr(resourceFullName, "facebook.0.callback_url", regexp.MustCompile(`^https:\/\/auth\.pingone\.(?:eu|com|asia|ca)\/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}\/rp\/callback\/facebook$`)),
					resource.TestCheckResourceAttr(resourceFullName, "google.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "github.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.#", "0"),
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

func TestAccIdentityProvider_Google(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_identity_provider.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckIdentityProviderDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProviderConfig_Google1(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "facebook.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "google.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "google.0.client_id", "dummyclientid1"),
					resource.TestCheckResourceAttr(resourceFullName, "google.0.client_secret", "dummyclientsecret1"),
					// resource.TestMatchResourceAttr(resourceFullName, "google.0.callback_url", regexp.MustCompile(`^https:\/\/auth\.pingone\.(?:eu|com|asia|ca)\/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}\/rp\/callback\/google$`)),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "github.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.#", "0"),
				),
			},
			{
				Config: testAccIdentityProviderConfig_Google2(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "facebook.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "google.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "google.0.client_id", "dummyclientid2"),
					resource.TestCheckResourceAttr(resourceFullName, "google.0.client_secret", "dummyclientsecret2"),
					// resource.TestMatchResourceAttr(resourceFullName, "google.0.callback_url", regexp.MustCompile(`^https:\/\/auth\.pingone\.(?:eu|com|asia|ca)\/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}\/rp\/callback\/google$`)),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "github.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.#", "0"),
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

func TestAccIdentityProvider_LinkedIn(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_identity_provider.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckIdentityProviderDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProviderConfig_LinkedIn1(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "facebook.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "google.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.0.client_id", "dummyclientid1"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.0.client_secret", "dummyclientsecret1"),
					// resource.TestMatchResourceAttr(resourceFullName, "linkedin.0.callback_url", regexp.MustCompile(`^https:\/\/auth\.pingone\.(?:eu|com|asia|ca)\/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}\/rp\/callback\/linkedin$`)),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "github.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.#", "0"),
				),
			},
			{
				Config: testAccIdentityProviderConfig_LinkedIn2(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "facebook.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "google.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.0.client_id", "dummyclientid2"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.0.client_secret", "dummyclientsecret2"),
					// resource.TestMatchResourceAttr(resourceFullName, "linkedin.0.callback_url", regexp.MustCompile(`^https:\/\/auth\.pingone\.(?:eu|com|asia|ca)\/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}\/rp\/callback\/linkedin$`)),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "github.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.#", "0"),
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

func TestAccIdentityProvider_Yahoo(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_identity_provider.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckIdentityProviderDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProviderConfig_Yahoo1(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "facebook.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "google.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.0.client_id", "dummyclientid1"),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.0.client_secret", "dummyclientsecret1"),
					// resource.TestMatchResourceAttr(resourceFullName, "yahoo.0.callback_url", regexp.MustCompile(`^https:\/\/auth\.pingone\.(?:eu|com|asia|ca)\/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}\/rp\/callback\/yahoo$`)),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "github.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.#", "0"),
				),
			},
			{
				Config: testAccIdentityProviderConfig_Yahoo2(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "facebook.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "google.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.0.client_id", "dummyclientid2"),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.0.client_secret", "dummyclientsecret2"),
					// resource.TestMatchResourceAttr(resourceFullName, "yahoo.0.callback_url", regexp.MustCompile(`^https:\/\/auth\.pingone\.(?:eu|com|asia|ca)\/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}\/rp\/callback\/yahoo$`)),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "github.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.#", "0"),
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

func TestAccIdentityProvider_Amazon(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_identity_provider.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckIdentityProviderDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProviderConfig_Amazon1(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "facebook.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "google.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.0.client_id", "dummyclientid1"),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.0.client_secret", "dummyclientsecret1"),
					// resource.TestMatchResourceAttr(resourceFullName, "amazon.0.callback_url", regexp.MustCompile(`^https:\/\/auth\.pingone\.(?:eu|com|asia|ca)\/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}\/rp\/callback\/amazon$`)),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "github.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.#", "0"),
				),
			},
			{
				Config: testAccIdentityProviderConfig_Amazon2(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "facebook.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "google.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.0.client_id", "dummyclientid2"),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.0.client_secret", "dummyclientsecret2"),
					// resource.TestMatchResourceAttr(resourceFullName, "amazon.0.callback_url", regexp.MustCompile(`^https:\/\/auth\.pingone\.(?:eu|com|asia|ca)\/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}\/rp\/callback\/amazon$`)),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "github.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.#", "0"),
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

func TestAccIdentityProvider_Twitter(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_identity_provider.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckIdentityProviderDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProviderConfig_Twitter1(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "facebook.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "google.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.0.client_id", "dummyclientid1"),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.0.client_secret", "dummyclientsecret1"),
					// resource.TestMatchResourceAttr(resourceFullName, "twitter.0.callback_url", regexp.MustCompile(`^https:\/\/auth\.pingone\.(?:eu|com|asia|ca)\/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}\/rp\/callback\/twitter$`)),
					resource.TestCheckResourceAttr(resourceFullName, "apple.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "github.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.#", "0"),
				),
			},
			{
				Config: testAccIdentityProviderConfig_Twitter2(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "facebook.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "google.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.0.client_id", "dummyclientid2"),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.0.client_secret", "dummyclientsecret2"),
					// resource.TestMatchResourceAttr(resourceFullName, "twitter.0.callback_url", regexp.MustCompile(`^https:\/\/auth\.pingone\.(?:eu|com|asia|ca)\/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}\/rp\/callback\/twitter$`)),
					resource.TestCheckResourceAttr(resourceFullName, "apple.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "github.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.#", "0"),
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

func TestAccIdentityProvider_Apple(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_identity_provider.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckIdentityProviderDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProviderConfig_Apple1(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "facebook.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "google.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.0.client_id", "dummyclientid1"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.0.client_secret_signing_key", "-----BEGIN PRIVATE KEY-----dummyclientsecretsigningkey1-----END PRIVATE KEY-----"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.0.key_id", "dummykeyi1"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.0.team_id", "dummyteam1"),
					// resource.TestMatchResourceAttr(resourceFullName, "apple.0.callback_url", regexp.MustCompile(`^https:\/\/auth\.pingone\.(?:eu|com|asia|ca)\/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}\/rp\/callback\/apple$`)),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "github.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.#", "0"),
				),
			},
			{
				Config: testAccIdentityProviderConfig_Apple2(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "facebook.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "google.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.0.client_id", "dummyclientid2"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.0.client_secret_signing_key", "-----BEGIN PRIVATE KEY-----dummyclientsecretsigningkey2-----END PRIVATE KEY-----"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.0.key_id", "dummykeyi2"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.0.team_id", "dummyteam2"),
					// resource.TestMatchResourceAttr(resourceFullName, "apple.0.callback_url", regexp.MustCompile(`^https:\/\/auth\.pingone\.(?:eu|com|asia|ca)\/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}\/rp\/callback\/apple$`)),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "github.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.#", "0"),
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

func TestAccIdentityProvider_Paypal(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_identity_provider.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckIdentityProviderDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProviderConfig_Paypal1(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "facebook.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "google.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.0.client_id", "dummyclientid1"),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.0.client_secret", "dummyclientsecret1"),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.0.client_environment", "sandbox"),
					// resource.TestMatchResourceAttr(resourceFullName, "paypal.0.callback_url", regexp.MustCompile(`^https:\/\/auth\.pingone\.(?:eu|com|asia|ca)\/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}\/rp\/callback\/paypal$`)),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "github.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.#", "0"),
				),
			},
			{
				Config: testAccIdentityProviderConfig_Paypal2(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "facebook.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "google.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.0.client_id", "dummyclientid2"),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.0.client_secret", "dummyclientsecret2"),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.0.client_environment", "live"),
					// resource.TestMatchResourceAttr(resourceFullName, "paypal.0.callback_url", regexp.MustCompile(`^https:\/\/auth\.pingone\.(?:eu|com|asia|ca)\/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}\/rp\/callback\/paypal$`)),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "github.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.#", "0"),
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

func TestAccIdentityProvider_Microsoft(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_identity_provider.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckIdentityProviderDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProviderConfig_Microsoft1(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "facebook.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "google.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.0.client_id", "dummyclientid1"),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.0.client_secret", "dummyclientsecret1"),
					// resource.TestMatchResourceAttr(resourceFullName, "microsoft.0.callback_url", regexp.MustCompile(`^https:\/\/auth\.pingone\.(?:eu|com|asia|ca)\/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}\/rp\/callback\/microsoft$`)),
					resource.TestCheckResourceAttr(resourceFullName, "github.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.#", "0"),
				),
			},
			{
				Config: testAccIdentityProviderConfig_Microsoft2(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "facebook.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "google.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.0.client_id", "dummyclientid2"),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.0.client_secret", "dummyclientsecret2"),
					// resource.TestMatchResourceAttr(resourceFullName, "microsoft.0.callback_url", regexp.MustCompile(`^https:\/\/auth\.pingone\.(?:eu|com|asia|ca)\/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}\/rp\/callback\/microsoft$`)),
					resource.TestCheckResourceAttr(resourceFullName, "github.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.#", "0"),
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

func TestAccIdentityProvider_Github(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_identity_provider.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckIdentityProviderDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProviderConfig_Github1(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "facebook.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "google.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "github.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "github.0.client_id", "dummyclientid1"),
					resource.TestCheckResourceAttr(resourceFullName, "github.0.client_secret", "dummyclientsecret1"),
					// resource.TestMatchResourceAttr(resourceFullName, "github.0.callback_url", regexp.MustCompile(`^https:\/\/auth\.pingone\.(?:eu|com|asia|ca)\/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}\/rp\/callback\/github$`)),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.#", "0"),
				),
			},
			{
				Config: testAccIdentityProviderConfig_Github2(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "facebook.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "google.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "github.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "github.0.client_id", "dummyclientid2"),
					resource.TestCheckResourceAttr(resourceFullName, "github.0.client_secret", "dummyclientsecret2"),
					// resource.TestMatchResourceAttr(resourceFullName, "github.0.callback_url", regexp.MustCompile(`^https:\/\/auth\.pingone\.(?:eu|com|asia|ca)\/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}\/rp\/callback\/github$`)),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.#", "0"),
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

func TestAccIdentityProvider_OIDC(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_identity_provider.%s", resourceName)

	name := resourceName

	data, _ := ioutil.ReadFile("../../acctest/test_assets/image/image-logo.gif")
	image := base64.StdEncoding.EncodeToString(data)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckIdentityProviderDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProviderConfig_OIDCFull(resourceName, name, image),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "login_button_icon.#", "1"),
					resource.TestMatchResourceAttr(resourceFullName, "login_button_icon.0.id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "login_button_icon.0.href", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
					resource.TestCheckResourceAttr(resourceFullName, "icon.#", "1"),
					resource.TestMatchResourceAttr(resourceFullName, "icon.0.id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "icon.0.href", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
					resource.TestCheckResourceAttr(resourceFullName, "facebook.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "google.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "github.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.0.authorization_endpoint", "https://www.pingidentity.com/authz"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.0.client_id", "dummyclientid1"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.0.client_secret", "dummyclientsecret1"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.0.discovery_endpoint", "https://www.pingidentity.com/discovery"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.0.issuer", "https://www.pingidentity.com/issuer"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.0.jwks_endpoint", "https://www.pingidentity.com/jwks"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "openid_connect.0.scopes.*", "openid"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "openid_connect.0.scopes.*", "scope1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "openid_connect.0.scopes.*", "scope2"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.0.token_endpoint", "https://www.pingidentity.com/token"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.0.token_endpoint_auth_method", "CLIENT_SECRET_POST"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.0.userinfo_endpoint", "https://www.pingidentity.com/userinfo"),
					// resource.TestMatchResourceAttr(resourceFullName, "openid_connect.0.callback_url", regexp.MustCompile(`^https:\/\/auth\.pingone\.(?:eu|com|asia|ca)\/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}\/rp\/callback\/openid_connect$`)),
					resource.TestCheckResourceAttr(resourceFullName, "saml.#", "0"),
				),
			},
			{
				Config: testAccIdentityProviderConfig_OIDCMinimal(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "facebook.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "google.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "github.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.0.authorization_endpoint", "https://www.pingidentity.com/authz2"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.0.client_id", "dummyclientid2"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.0.client_secret", "dummyclientsecret2"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.0.discovery_endpoint", ""),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.0.issuer", "https://www.pingidentity.com/issuer2"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.0.jwks_endpoint", "https://www.pingidentity.com/jwks2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "openid_connect.0.scopes.*", "openid"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "openid_connect.0.scopes.*", "scope3"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "openid_connect.0.scopes.*", "scope4"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.0.token_endpoint", "https://www.pingidentity.com/token2"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.0.token_endpoint_auth_method", "CLIENT_SECRET_BASIC"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.0.userinfo_endpoint", ""),
					// resource.TestMatchResourceAttr(resourceFullName, "openid_connect.0.callback_url", regexp.MustCompile(`^https:\/\/auth\.pingone\.(?:eu|com|asia|ca)\/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}\/rp\/callback\/openid_connect$`)),
					resource.TestCheckResourceAttr(resourceFullName, "saml.#", "0"),
				),
			},
			{
				Config: testAccIdentityProviderConfig_OIDCFull(resourceName, name, image),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "login_button_icon.#", "1"),
					resource.TestMatchResourceAttr(resourceFullName, "login_button_icon.0.id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "login_button_icon.0.href", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
					resource.TestCheckResourceAttr(resourceFullName, "icon.#", "1"),
					resource.TestMatchResourceAttr(resourceFullName, "icon.0.id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "icon.0.href", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
					resource.TestCheckResourceAttr(resourceFullName, "facebook.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "google.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "github.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.0.authorization_endpoint", "https://www.pingidentity.com/authz"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.0.client_id", "dummyclientid1"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.0.client_secret", "dummyclientsecret1"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.0.discovery_endpoint", "https://www.pingidentity.com/discovery"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.0.issuer", "https://www.pingidentity.com/issuer"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.0.jwks_endpoint", "https://www.pingidentity.com/jwks"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "openid_connect.0.scopes.*", "openid"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "openid_connect.0.scopes.*", "scope1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "openid_connect.0.scopes.*", "scope2"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.0.token_endpoint", "https://www.pingidentity.com/token"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.0.token_endpoint_auth_method", "CLIENT_SECRET_POST"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.0.userinfo_endpoint", "https://www.pingidentity.com/userinfo"),
					// resource.TestMatchResourceAttr(resourceFullName, "openid_connect.0.callback_url", regexp.MustCompile(`^https:\/\/auth\.pingone\.(?:eu|com|asia|ca)\/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}\/rp\/callback\/openid_connect$`)),
					resource.TestCheckResourceAttr(resourceFullName, "saml.#", "0"),
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

func TestAccIdentityProvider_SAML(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_identity_provider.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { t.Skipf("Test to be re-defined") }, // Needs redefinition
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckIdentityProviderDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProviderConfig_SAMLFull(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "facebook.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "google.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "github.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.0.authentication_request_signed", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.0.idp_entity_id", "idp:entity"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.0.sp_entity_id", "sp:entity"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.0.idp_verification_certificate_ids", "https://www.pingidentity.com/discovery"),
					resource.TestMatchResourceAttr(resourceFullName, "saml.0.sp_signing_key_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "saml.0.sso_binding", "HTTP_POST"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.0.sso_endpoint", "https://www.pingidentity.com/sso"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.0.slo_binding", "HTTP_REDIRECT"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.0.slo_endpoint", "https://dummy-slo-endpoint.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.0.slo_response_endpoint", "https://dummy-slo-response-endpoint.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.0.slo_window", "1"),
				),
			},
			{
				Config: testAccIdentityProviderConfig_SAMLMinimal(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "facebook.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "google.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "github.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.0.authentication_request_signed", ""),
					resource.TestCheckResourceAttr(resourceFullName, "saml.0.idp_entity_id", "idp:entity"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.0.sp_entity_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "saml.0.idp_verification_certificate_ids.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.0.sp_signing_key_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "saml.0.sso_binding", "HTTP_REDIRECT"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.0.sso_endpoint", "https://www.pingidentity.com/sso"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.0.slo_binding", ""),
					resource.TestCheckResourceAttr(resourceFullName, "saml.0.slo_endpoint", ""),
					resource.TestCheckResourceAttr(resourceFullName, "saml.0.slo_response_endpoint", ""),
					resource.TestCheckResourceAttr(resourceFullName, "saml.0.slo_window", "0"),
				),
			},
			{
				Config: testAccIdentityProviderConfig_SAMLFull(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "facebook.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "google.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "github.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.0.authentication_request_signed", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.0.idp_entity_id", "idp:entity"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.0.sp_entity_id", "sp:entity"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.0.idp_verification_certificate_ids", "https://www.pingidentity.com/discovery"),
					resource.TestMatchResourceAttr(resourceFullName, "saml.0.sp_signing_key_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "saml.0.sso_binding", "HTTP_POST"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.0.sso_endpoint", "https://www.pingidentity.com/sso"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.0.slo_binding", "HTTP_REDIRECT"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.0.slo_endpoint", "https://dummy-slo-endpoint.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.0.slo_response_endpoint", "https://dummy-slo-response-endpoint.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.0.slo_window", "1"),
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

func TestAccIdentityProvider_ChangeProvider(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_identity_provider.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckIdentityProviderDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProviderConfig_Apple1(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "facebook.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "google.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.0.client_id", "dummyclientid1"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.0.client_secret_signing_key", "-----BEGIN PRIVATE KEY-----dummyclientsecretsigningkey1-----END PRIVATE KEY-----"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.0.key_id", "dummykeyi1"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.0.team_id", "dummyteam1"),
					// resource.TestMatchResourceAttr(resourceFullName, "apple.0.callback_url", regexp.MustCompile(`^https:\/\/auth\.pingone\.(?:eu|com|asia|ca)\/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}\/rp\/callback\/apple$`)),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "github.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.#", "0"),
				),
			},
			{
				Config: testAccIdentityProviderConfig_Github2(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "facebook.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "google.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "github.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "github.0.client_id", "dummyclientid2"),
					resource.TestCheckResourceAttr(resourceFullName, "github.0.client_secret", "dummyclientsecret2"),
					// resource.TestMatchResourceAttr(resourceFullName, "github.0.callback_url", regexp.MustCompile(`^https:\/\/auth\.pingone\.(?:eu|com|asia|ca)\/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}\/rp\/callback\/github$`)),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.#", "0"),
				),
			},
		},
	})
}

func TestAccIdentityProvider_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_identity_provider.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckIdentityProviderDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccIdentityProviderConfig_Minimal(resourceName, name),
			},
			// Errors
			{
				ResourceName: resourceFullName,
				ImportState:  true,
				ExpectError:  regexp.MustCompile(`Invalid import ID specified \(".*"\).  The ID should be in the format "environment_id/identity_provider_id".`),
			},
			{
				ResourceName:  resourceFullName,
				ImportStateId: "/",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`Invalid import ID specified \(".*"\).  The ID should be in the format "environment_id/identity_provider_id".`),
			},
			{
				ResourceName:  resourceFullName,
				ImportStateId: "badformat/badformat",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`Invalid import ID specified \(".*"\).  The ID should be in the format "environment_id/identity_provider_id".`),
			},
		},
	})
}

func testAccIdentityProviderConfig_NewEnv(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s
resource "pingone_identity_provider" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s"

  google {
    client_id     = "testclientid"
    client_secret = "testclientsecret"
  }
}
		`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccIdentityProviderConfig_Full(resourceName, name, image string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_population" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_image" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  image_file_base64 = "%[4]s"
}

resource "pingone_identity_provider" "%[2]s" {
  environment_id             = data.pingone_environment.general_test.id
  name                       = "%[3]s"
  description                = "My test identity provider"
  enabled                    = true
  registration_population_id = pingone_population.%[2]s.id

  icon {
    id   = pingone_image.%[2]s.id
    href = pingone_image.%[2]s.uploaded_image[0].href
  }

  login_button_icon {
    id   = pingone_image.%[2]s.id
    href = pingone_image.%[2]s.uploaded_image[0].href
  }

  google {
    client_id     = "testclientid"
    client_secret = "testclientsecret"
  }
}
		`, acctest.GenericSandboxEnvironment(), resourceName, name, image)
}

func testAccIdentityProviderConfig_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s
resource "pingone_identity_provider" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  google {
    client_id     = "testclientid"
    client_secret = "testclientsecret"
  }
}
		`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccIdentityProviderConfig_Facebook1(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s
resource "pingone_identity_provider" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  facebook {
    app_id     = "dummyappid1"
    app_secret = "dummyappsecret1"
  }
}
		`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccIdentityProviderConfig_Facebook2(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s
resource "pingone_identity_provider" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  facebook {
    app_id     = "dummyappid2"
    app_secret = "dummyappsecret2"
  }
}
		`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccIdentityProviderConfig_Google1(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s
resource "pingone_identity_provider" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  google {
    client_id     = "dummyclientid1"
    client_secret = "dummyclientsecret1"
  }
}
		`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccIdentityProviderConfig_Google2(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s
resource "pingone_identity_provider" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  google {
    client_id     = "dummyclientid2"
    client_secret = "dummyclientsecret2"
  }
}
		`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccIdentityProviderConfig_LinkedIn1(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s
resource "pingone_identity_provider" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  linkedin {
    client_id     = "dummyclientid1"
    client_secret = "dummyclientsecret1"
  }
}
		`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccIdentityProviderConfig_LinkedIn2(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s
resource "pingone_identity_provider" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  linkedin {
    client_id     = "dummyclientid2"
    client_secret = "dummyclientsecret2"
  }
}
		`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccIdentityProviderConfig_Yahoo1(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s
resource "pingone_identity_provider" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  yahoo {
    client_id     = "dummyclientid1"
    client_secret = "dummyclientsecret1"
  }
}
		`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccIdentityProviderConfig_Yahoo2(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s
resource "pingone_identity_provider" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  yahoo {
    client_id     = "dummyclientid2"
    client_secret = "dummyclientsecret2"
  }
}
		`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccIdentityProviderConfig_Amazon1(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s
resource "pingone_identity_provider" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  amazon {
    client_id     = "dummyclientid1"
    client_secret = "dummyclientsecret1"
  }
}
		`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccIdentityProviderConfig_Amazon2(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s
resource "pingone_identity_provider" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  amazon {
    client_id     = "dummyclientid2"
    client_secret = "dummyclientsecret2"
  }
}
		`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccIdentityProviderConfig_Twitter1(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s
resource "pingone_identity_provider" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  twitter {
    client_id     = "dummyclientid1"
    client_secret = "dummyclientsecret1"
  }
}
		`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccIdentityProviderConfig_Twitter2(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s
resource "pingone_identity_provider" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  twitter {
    client_id     = "dummyclientid2"
    client_secret = "dummyclientsecret2"
  }
}
		`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccIdentityProviderConfig_Apple1(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s
resource "pingone_identity_provider" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  apple {
    client_id                 = "dummyclientid1"
    client_secret_signing_key = "-----BEGIN PRIVATE KEY-----dummyclientsecretsigningkey1-----END PRIVATE KEY-----"
    key_id                    = "dummykeyi1"
    team_id                   = "dummyteam1"
  }
}
		`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccIdentityProviderConfig_Apple2(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s
resource "pingone_identity_provider" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  apple {
    client_id                 = "dummyclientid2"
    client_secret_signing_key = "-----BEGIN PRIVATE KEY-----dummyclientsecretsigningkey2-----END PRIVATE KEY-----"
    key_id                    = "dummykeyi2"
    team_id                   = "dummyteam2"
  }
}
		`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccIdentityProviderConfig_Paypal1(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s
resource "pingone_identity_provider" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  paypal {
    client_id          = "dummyclientid1"
    client_secret      = "dummyclientsecret1"
    client_environment = "sandbox"
  }
}
		`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccIdentityProviderConfig_Paypal2(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s
resource "pingone_identity_provider" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  paypal {
    client_id          = "dummyclientid2"
    client_secret      = "dummyclientsecret2"
    client_environment = "live"
  }
}
		`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccIdentityProviderConfig_Microsoft1(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s
resource "pingone_identity_provider" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  microsoft {
    client_id     = "dummyclientid1"
    client_secret = "dummyclientsecret1"
  }
}
		`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccIdentityProviderConfig_Microsoft2(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s
resource "pingone_identity_provider" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  microsoft {
    client_id     = "dummyclientid2"
    client_secret = "dummyclientsecret2"
  }
}
		`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccIdentityProviderConfig_Github1(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s
resource "pingone_identity_provider" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  github {
    client_id     = "dummyclientid1"
    client_secret = "dummyclientsecret1"
  }
}
		`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccIdentityProviderConfig_Github2(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s
resource "pingone_identity_provider" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  github {
    client_id     = "dummyclientid2"
    client_secret = "dummyclientsecret2"
  }
}
		`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccIdentityProviderConfig_OIDCFull(resourceName, name, image string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_image" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  image_file_base64 = "%[4]s"
}

resource "pingone_identity_provider" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  icon {
    id   = pingone_image.%[2]s.id
    href = pingone_image.%[2]s.uploaded_image[0].href
  }

  login_button_icon {
    id   = pingone_image.%[2]s.id
    href = pingone_image.%[2]s.uploaded_image[0].href
  }

  openid_connect {
    authorization_endpoint     = "https://www.pingidentity.com/authz"
    client_id                  = "dummyclientid1"
    client_secret              = "dummyclientsecret1"
    discovery_endpoint         = "https://www.pingidentity.com/discovery"
    issuer                     = "https://www.pingidentity.com/issuer"
    jwks_endpoint              = "https://www.pingidentity.com/jwks"
    scopes                     = ["openid", "scope1", "scope2"]
    token_endpoint             = "https://www.pingidentity.com/token"
    token_endpoint_auth_method = "CLIENT_SECRET_POST"
    userinfo_endpoint          = "https://www.pingidentity.com/userinfo"
  }
}
		`, acctest.GenericSandboxEnvironment(), resourceName, name, image)
}

func testAccIdentityProviderConfig_OIDCMinimal(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s
resource "pingone_identity_provider" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  openid_connect {
    authorization_endpoint = "https://www.pingidentity.com/authz2"
    client_id              = "dummyclientid2"
    client_secret          = "dummyclientsecret2"
    issuer                 = "https://www.pingidentity.com/issuer2"
    jwks_endpoint          = "https://www.pingidentity.com/jwks2"
    scopes                 = ["openid", "scope3", "scope4"]
    token_endpoint         = "https://www.pingidentity.com/token2"
  }
}
		`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccIdentityProviderConfig_SAMLFull(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_identity_provider" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  saml {
    authentication_request_signed    = true
    idp_entity_id                    = "idp:entity"
    sp_entity_id                     = "sp:entity"
    idp_verification_certificate_ids = []
    // sp_signing_key_id = 
    sso_binding           = "HTTP_POST"
    sso_endpoint          = "https://www.pingidentity.com/sso"
    slo_binding           = "HTTP_REDIRECT"
    slo_endpoint          = "https://dummy-slo-endpoint.pingidentity.com"
    slo_response_endpoint = "https://dummy-slo-response-endpoint.pingidentity.com"
    slo_window            = 1
  }
}
		`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccIdentityProviderConfig_SAMLMinimal(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s
resource "pingone_identity_provider" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  saml {
    idp_entity_id = "idp:entity"
    sso_binding   = "HTTP_REDIRECT"
    sso_endpoint  = "https://www.pingidentity.com/sso"
  }
}
		`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
