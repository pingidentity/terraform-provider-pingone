// Copyright © 2025 Ping Identity Corporation

package sso_test

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
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/sso"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccIdentityProvider_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_identity_provider.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	var resourceID, environmentID string

	var p1Client *client.Client
	var ctx = context.Background()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)

			p1Client = acctest.PreCheckTestClient(ctx, t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.IdentityProvider_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccIdentityProviderConfig_Minimal(resourceName, name),
				Check:  sso.IdentityProvider_GetIDs(resourceFullName, &environmentID, &resourceID),
			},
			{
				PreConfig: func() {
					sso.IdentityProvider_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID, resourceID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the environment
			{
				Config: testAccIdentityProviderConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check:  sso.IdentityProvider_GetIDs(resourceFullName, &environmentID, &resourceID),
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

func TestAccIdentityProvider_NewEnv(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_identity_provider.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.IdentityProvider_CheckDestroy,
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

func TestAccIdentityProvider_Change(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_identity_provider.%s", resourceName)

	name := resourceName

	data, _ := os.ReadFile("../../acctest/test_assets/image/image-logo.gif")
	image := base64.StdEncoding.EncodeToString(data)

	fullStep := resource.TestStep{
		Config: testAccIdentityProviderConfig_Full(resourceName, fmt.Sprintf("%s 1", name), image),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "name", fmt.Sprintf("%s 1", name)),
			resource.TestCheckResourceAttr(resourceFullName, "description", "My test identity provider"),
			resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
			resource.TestMatchResourceAttr(resourceFullName, "registration_population_id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "login_button_icon.id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "login_button_icon.href", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
			resource.TestMatchResourceAttr(resourceFullName, "icon.id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "icon.href", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
		),
	}

	minimalStep := resource.TestStep{
		Config: testAccIdentityProviderConfig_Minimal(resourceName, fmt.Sprintf("%s 2", name)),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "name", fmt.Sprintf("%s 2", name)),
			resource.TestCheckNoResourceAttr(resourceFullName, "description"),
			resource.TestCheckResourceAttr(resourceFullName, "enabled", "false"),
			resource.TestCheckNoResourceAttr(resourceFullName, "registration_population_id"),
			resource.TestCheckResourceAttr(resourceFullName, "login_button_icon.%", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "icon.%", "0"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.IdentityProvider_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full step
			fullStep,
			{
				Config:  testAccIdentityProviderConfig_Full(resourceName, fmt.Sprintf("%s 1", name), image),
				Destroy: true,
			},
			// Minimal step
			minimalStep,
			{
				Config:  testAccIdentityProviderConfig_Minimal(resourceName, fmt.Sprintf("%s 2", name)),
				Destroy: true,
			},
			// Change
			fullStep,
			minimalStep,
			fullStep,
		},
	})
}

func TestAccIdentityProvider_Facebook(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_identity_provider.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.IdentityProvider_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProviderConfig_Facebook1(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "facebook.app_id", "dummyappid1"),
					resource.TestCheckResourceAttr(resourceFullName, "facebook.app_secret", "dummyappsecret1"),
					// resource.TestMatchResourceAttr(resourceFullName, "facebook.callback_url", regexp.MustCompile(`^https:\/\/auth\.pingone\.(?:eu|com|asia|ca)\/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}\/rp\/callback\/facebook$`)),
					resource.TestCheckResourceAttr(resourceFullName, "google.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin_oidc.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "github.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.%", "0"),
				),
			},
			{
				Config: testAccIdentityProviderConfig_Facebook2(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "facebook.app_id", "dummyappid2"),
					resource.TestCheckResourceAttr(resourceFullName, "facebook.app_secret", "dummyappsecret2"),
					// resource.TestMatchResourceAttr(resourceFullName, "facebook.callback_url", regexp.MustCompile(`^https:\/\/auth\.pingone\.(?:eu|com|asia|ca)\/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}\/rp\/callback\/facebook$`)),
					resource.TestCheckResourceAttr(resourceFullName, "google.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin_oidc.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "github.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.%", "0"),
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
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.IdentityProvider_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProviderConfig_Google1(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "facebook.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "google.client_id", "dummyclientid1"),
					resource.TestCheckResourceAttr(resourceFullName, "google.client_secret", "dummyclientsecret1"),
					// resource.TestMatchResourceAttr(resourceFullName, "google.callback_url", regexp.MustCompile(`^https:\/\/auth\.pingone\.(?:eu|com|asia|ca)\/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}\/rp\/callback\/google$`)),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin_oidc.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "github.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.%", "0"),
				),
			},
			{
				Config: testAccIdentityProviderConfig_Google2(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "facebook.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "google.client_id", "dummyclientid2"),
					resource.TestCheckResourceAttr(resourceFullName, "google.client_secret", "dummyclientsecret2"),
					// resource.TestMatchResourceAttr(resourceFullName, "google.callback_url", regexp.MustCompile(`^https:\/\/auth\.pingone\.(?:eu|com|asia|ca)\/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}\/rp\/callback\/google$`)),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin_oidc.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "github.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.%", "0"),
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
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.IdentityProvider_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProviderConfig_LinkedIn1(resourceName, name, "linkedin"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "facebook.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "google.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.client_id", "dummyclientid1"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.client_secret", "dummyclientsecret1"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin_oidc.%", "0"),
					// resource.TestMatchResourceAttr(resourceFullName, "linkedin.callback_url", regexp.MustCompile(`^https:\/\/auth\.pingone\.(?:eu|com|asia|ca)\/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}\/rp\/callback\/linkedin$`)),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "github.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.%", "0"),
				),
			},
			{
				Config: testAccIdentityProviderConfig_LinkedIn2(resourceName, name, "linkedin"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "facebook.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "google.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.client_id", "dummyclientid2"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.client_secret", "dummyclientsecret2"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin_oidc.%", "0"),
					// resource.TestMatchResourceAttr(resourceFullName, "linkedin.callback_url", regexp.MustCompile(`^https:\/\/auth\.pingone\.(?:eu|com|asia|ca)\/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}\/rp\/callback\/linkedin$`)),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "github.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.%", "0"),
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

func TestAccIdentityProvider_LinkedInOIDC(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_identity_provider.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.IdentityProvider_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProviderConfig_LinkedIn1(resourceName, name, "linkedin_oidc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "facebook.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "google.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin_oidc.client_id", "dummyclientid1"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin_oidc.client_secret", "dummyclientsecret1"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.%", "0"),
					// resource.TestMatchResourceAttr(resourceFullName, "linkedin.callback_url", regexp.MustCompile(`^https:\/\/auth\.pingone\.(?:eu|com|asia|ca)\/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}\/rp\/callback\/linkedin$`)),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "github.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.%", "0"),
				),
			},
			{
				Config: testAccIdentityProviderConfig_LinkedIn2(resourceName, name, "linkedin_oidc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "facebook.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "google.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin_oidc.client_id", "dummyclientid2"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin_oidc.client_secret", "dummyclientsecret2"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.%", "0"),
					// resource.TestMatchResourceAttr(resourceFullName, "linkedin.callback_url", regexp.MustCompile(`^https:\/\/auth\.pingone\.(?:eu|com|asia|ca)\/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}\/rp\/callback\/linkedin$`)),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "github.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.%", "0"),
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
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.IdentityProvider_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProviderConfig_Yahoo1(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "facebook.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "google.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin_oidc.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.client_id", "dummyclientid1"),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.client_secret", "dummyclientsecret1"),
					// resource.TestMatchResourceAttr(resourceFullName, "yahoo.callback_url", regexp.MustCompile(`^https:\/\/auth\.pingone\.(?:eu|com|asia|ca)\/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}\/rp\/callback\/yahoo$`)),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "github.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.%", "0"),
				),
			},
			{
				Config: testAccIdentityProviderConfig_Yahoo2(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "facebook.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "google.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin_oidc.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.client_id", "dummyclientid2"),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.client_secret", "dummyclientsecret2"),
					// resource.TestMatchResourceAttr(resourceFullName, "yahoo.callback_url", regexp.MustCompile(`^https:\/\/auth\.pingone\.(?:eu|com|asia|ca)\/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}\/rp\/callback\/yahoo$`)),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "github.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.%", "0"),
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
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.IdentityProvider_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProviderConfig_Amazon1(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "facebook.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "google.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin_oidc.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.client_id", "dummyclientid1"),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.client_secret", "dummyclientsecret1"),
					// resource.TestMatchResourceAttr(resourceFullName, "amazon.callback_url", regexp.MustCompile(`^https:\/\/auth\.pingone\.(?:eu|com|asia|ca)\/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}\/rp\/callback\/amazon$`)),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "github.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.%", "0"),
				),
			},
			{
				Config: testAccIdentityProviderConfig_Amazon2(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "facebook.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "google.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin_oidc.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.client_id", "dummyclientid2"),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.client_secret", "dummyclientsecret2"),
					// resource.TestMatchResourceAttr(resourceFullName, "amazon.callback_url", regexp.MustCompile(`^https:\/\/auth\.pingone\.(?:eu|com|asia|ca)\/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}\/rp\/callback\/amazon$`)),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "github.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.%", "0"),
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
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.IdentityProvider_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProviderConfig_Twitter1(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "facebook.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "google.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin_oidc.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.client_id", "dummyclientid1"),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.client_secret", "dummyclientsecret1"),
					// resource.TestMatchResourceAttr(resourceFullName, "twitter.callback_url", regexp.MustCompile(`^https:\/\/auth\.pingone\.(?:eu|com|asia|ca)\/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}\/rp\/callback\/twitter$`)),
					resource.TestCheckResourceAttr(resourceFullName, "apple.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "github.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.%", "0"),
				),
			},
			{
				Config: testAccIdentityProviderConfig_Twitter2(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "facebook.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "google.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin_oidc.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.client_id", "dummyclientid2"),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.client_secret", "dummyclientsecret2"),
					// resource.TestMatchResourceAttr(resourceFullName, "twitter.callback_url", regexp.MustCompile(`^https:\/\/auth\.pingone\.(?:eu|com|asia|ca)\/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}\/rp\/callback\/twitter$`)),
					resource.TestCheckResourceAttr(resourceFullName, "apple.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "github.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.%", "0"),
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
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.IdentityProvider_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProviderConfig_Apple1(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "facebook.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "google.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin_oidc.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.client_id", "dummyclientid1"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.client_secret_signing_key", "-----BEGIN PRIVATE KEY-----dummyclientsecretsigningkey1-----END PRIVATE KEY-----"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.key_id", "dummykeyi1"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.team_id", "dummyteam1"),
					// resource.TestMatchResourceAttr(resourceFullName, "apple.callback_url", regexp.MustCompile(`^https:\/\/auth\.pingone\.(?:eu|com|asia|ca)\/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}\/rp\/callback\/apple$`)),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "github.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.%", "0"),
				),
			},
			{
				Config: testAccIdentityProviderConfig_Apple2(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "facebook.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "google.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin_oidc.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.client_id", "dummyclientid2"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.client_secret_signing_key", "-----BEGIN PRIVATE KEY-----dummyclientsecretsigningkey2-----END PRIVATE KEY-----"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.key_id", "dummykeyi2"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.team_id", "dummyteam2"),
					// resource.TestMatchResourceAttr(resourceFullName, "apple.callback_url", regexp.MustCompile(`^https:\/\/auth\.pingone\.(?:eu|com|asia|ca)\/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}\/rp\/callback\/apple$`)),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "github.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.%", "0"),
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
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.IdentityProvider_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProviderConfig_Paypal1(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "facebook.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "google.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin_oidc.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.client_id", "dummyclientid1"),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.client_secret", "dummyclientsecret1"),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.client_environment", "sandbox"),
					// resource.TestMatchResourceAttr(resourceFullName, "paypal.callback_url", regexp.MustCompile(`^https:\/\/auth\.pingone\.(?:eu|com|asia|ca)\/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}\/rp\/callback\/paypal$`)),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "github.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.%", "0"),
				),
			},
			{
				Config: testAccIdentityProviderConfig_Paypal2(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "facebook.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "google.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin_oidc.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.client_id", "dummyclientid2"),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.client_secret", "dummyclientsecret2"),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.client_environment", "live"),
					// resource.TestMatchResourceAttr(resourceFullName, "paypal.callback_url", regexp.MustCompile(`^https:\/\/auth\.pingone\.(?:eu|com|asia|ca)\/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}\/rp\/callback\/paypal$`)),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "github.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.%", "0"),
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
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.IdentityProvider_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProviderConfig_Microsoft1(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "facebook.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "google.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin_oidc.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.client_id", "dummyclientid1"),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.client_secret", "dummyclientsecret1"),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.tenant_id", "dummytenantid1"),
					// resource.TestMatchResourceAttr(resourceFullName, "microsoft.callback_url", regexp.MustCompile(`^https:\/\/auth\.pingone\.(?:eu|com|asia|ca)\/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}\/rp\/callback\/microsoft$`)),
					resource.TestCheckResourceAttr(resourceFullName, "github.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.%", "0"),
				),
			},
			{
				Config: testAccIdentityProviderConfig_Microsoft2(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "facebook.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "google.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin_oidc.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.client_id", "dummyclientid2"),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.client_secret", "dummyclientsecret2"),
					// resource.TestMatchResourceAttr(resourceFullName, "microsoft.callback_url", regexp.MustCompile(`^https:\/\/auth\.pingone\.(?:eu|com|asia|ca)\/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}\/rp\/callback\/microsoft$`)),
					resource.TestCheckResourceAttr(resourceFullName, "github.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.%", "0"),
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
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.IdentityProvider_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProviderConfig_Github1(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "facebook.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "google.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin_oidc.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "github.client_id", "dummyclientid1"),
					resource.TestCheckResourceAttr(resourceFullName, "github.client_secret", "dummyclientsecret1"),
					// resource.TestMatchResourceAttr(resourceFullName, "github.callback_url", regexp.MustCompile(`^https:\/\/auth\.pingone\.(?:eu|com|asia|ca)\/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}\/rp\/callback\/github$`)),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.%", "0"),
				),
			},
			{
				Config: testAccIdentityProviderConfig_Github2(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "facebook.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "google.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin_oidc.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "github.client_id", "dummyclientid2"),
					resource.TestCheckResourceAttr(resourceFullName, "github.client_secret", "dummyclientsecret2"),
					// resource.TestMatchResourceAttr(resourceFullName, "github.callback_url", regexp.MustCompile(`^https:\/\/auth\.pingone\.(?:eu|com|asia|ca)\/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}\/rp\/callback\/github$`)),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.%", "0"),
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

	data, _ := os.ReadFile("../../acctest/test_assets/image/image-logo.gif")
	image := base64.StdEncoding.EncodeToString(data)

	fullStep := resource.TestStep{
		Config: testAccIdentityProviderConfig_OIDCFull(resourceName, name, image),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "login_button_icon.id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "login_button_icon.href", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
			resource.TestMatchResourceAttr(resourceFullName, "icon.id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "icon.href", regexp.MustCompile(`^https:\/\/uploads\.pingone\.((eu)|(com)|(asia)|(ca))\/environments\/[a-zA-Z0-9-]*\/images\/[a-zA-Z0-9-]*_[a-zA-Z0-9-]*_original\.png$`)),
			resource.TestCheckResourceAttr(resourceFullName, "facebook.%", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "google.%", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "linkedin.%", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "linkedin_oidc.%", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "yahoo.%", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "amazon.%", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "twitter.%", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "apple.%", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "paypal.%", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "microsoft.%", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "github.%", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "openid_connect.authorization_endpoint", "https://www.pingidentity.com/authz"),
			resource.TestCheckResourceAttr(resourceFullName, "openid_connect.client_id", "dummyclientid1"),
			resource.TestCheckResourceAttr(resourceFullName, "openid_connect.client_secret", "dummyclientsecret1"),
			resource.TestCheckResourceAttr(resourceFullName, "openid_connect.discovery_endpoint", "https://www.pingidentity.com/discovery"),
			resource.TestCheckResourceAttr(resourceFullName, "openid_connect.issuer", "https://www.pingidentity.com/issuer"),
			resource.TestCheckResourceAttr(resourceFullName, "openid_connect.pkce_method", "S256"),
			resource.TestCheckResourceAttr(resourceFullName, "openid_connect.jwks_endpoint", "https://www.pingidentity.com/jwks"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "openid_connect.scopes.*", "openid"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "openid_connect.scopes.*", "scope1"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "openid_connect.scopes.*", "scope2"),
			resource.TestCheckResourceAttr(resourceFullName, "openid_connect.token_endpoint", "https://www.pingidentity.com/token"),
			resource.TestCheckResourceAttr(resourceFullName, "openid_connect.token_endpoint_auth_method", "CLIENT_SECRET_POST"),
			resource.TestCheckResourceAttr(resourceFullName, "openid_connect.userinfo_endpoint", "https://www.pingidentity.com/userinfo"),
			// resource.TestMatchResourceAttr(resourceFullName, "openid_connect.callback_url", regexp.MustCompile(`^https:\/\/auth\.pingone\.(?:eu|com|asia|ca)\/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}\/rp\/callback\/openid_connect$`)),
			resource.TestCheckResourceAttr(resourceFullName, "saml.%", "0"),
		),
	}

	minimalStep := resource.TestStep{
		Config: testAccIdentityProviderConfig_OIDCMinimal(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "facebook.%", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "google.%", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "linkedin.%", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "linkedin_oidc.%", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "yahoo.%", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "amazon.%", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "twitter.%", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "apple.%", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "paypal.%", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "microsoft.%", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "github.%", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "openid_connect.authorization_endpoint", "https://www.pingidentity.com/authz2"),
			resource.TestCheckResourceAttr(resourceFullName, "openid_connect.client_id", "dummyclientid2"),
			resource.TestCheckResourceAttr(resourceFullName, "openid_connect.client_secret", "dummyclientsecret2"),
			resource.TestCheckNoResourceAttr(resourceFullName, "openid_connect.discovery_endpoint"),
			resource.TestCheckResourceAttr(resourceFullName, "openid_connect.issuer", "https://www.pingidentity.com/issuer2"),
			resource.TestCheckResourceAttr(resourceFullName, "openid_connect.pkce_method", "NONE"),
			resource.TestCheckResourceAttr(resourceFullName, "openid_connect.jwks_endpoint", "https://www.pingidentity.com/jwks2"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "openid_connect.scopes.*", "openid"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "openid_connect.scopes.*", "scope3"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "openid_connect.scopes.*", "scope4"),
			resource.TestCheckResourceAttr(resourceFullName, "openid_connect.token_endpoint", "https://www.pingidentity.com/token2"),
			resource.TestCheckResourceAttr(resourceFullName, "openid_connect.token_endpoint_auth_method", "CLIENT_SECRET_BASIC"),
			resource.TestCheckNoResourceAttr(resourceFullName, "openid_connect.userinfo_endpoint"),
			// resource.TestMatchResourceAttr(resourceFullName, "openid_connect.callback_url", regexp.MustCompile(`^https:\/\/auth\.pingone\.(?:eu|com|asia|ca)\/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}\/rp\/callback\/openid_connect$`)),
			resource.TestCheckResourceAttr(resourceFullName, "saml.%", "0"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.IdentityProvider_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			fullStep,
			{
				Config:  testAccIdentityProviderConfig_OIDCFull(resourceName, name, image),
				Destroy: true,
			},
			// Minimal
			minimalStep,
			{
				Config:  testAccIdentityProviderConfig_OIDCMinimal(resourceName, name),
				Destroy: true,
			},
			// Change
			fullStep,
			minimalStep,
			fullStep,
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

	pem_cert := os.Getenv("PINGONE_KEY_PEM_CERT")

	fullStep := resource.TestStep{
		Config: testAccIdentityProviderConfig_SAMLFull(resourceName, name, pem_cert),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "facebook.%", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "google.%", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "linkedin.%", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "linkedin_oidc.%", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "yahoo.%", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "amazon.%", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "twitter.%", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "apple.%", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "paypal.%", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "microsoft.%", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "github.%", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "openid_connect.%", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "saml.authentication_request_signed", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "saml.idp_entity_id", fmt.Sprintf("idp:%s", name)),
			resource.TestCheckResourceAttr(resourceFullName, "saml.sp_entity_id", fmt.Sprintf("sp:%s", name)),
			resource.TestCheckResourceAttr(resourceFullName, "saml.idp_verification.certificates.#", "1"),
			resource.TestMatchResourceAttr(resourceFullName, "saml.idp_verification.certificates.0.id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "saml.sp_signing.key.id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "saml.sp_signing.algorithm", "SHA512withRSA"),
			resource.TestCheckResourceAttr(resourceFullName, "saml.sso_binding", "HTTP_POST"),
			resource.TestCheckResourceAttr(resourceFullName, "saml.sso_endpoint", "https://www.pingidentity.com/sso"),
			resource.TestCheckResourceAttr(resourceFullName, "saml.slo_binding", "HTTP_REDIRECT"),
			resource.TestCheckResourceAttr(resourceFullName, "saml.slo_endpoint", "https://dummy-slo-endpoint.pingidentity.com"),
			resource.TestCheckResourceAttr(resourceFullName, "saml.slo_response_endpoint", "https://dummy-slo-response-endpoint.pingidentity.com"),
			resource.TestCheckResourceAttr(resourceFullName, "saml.slo_window", "1"),
		),
	}

	minimalStep := resource.TestStep{
		Config: testAccIdentityProviderConfig_SAMLMinimal(resourceName, name, pem_cert),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "facebook.%", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "google.%", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "linkedin.%", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "linkedin_oidc.%", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "yahoo.%", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "amazon.%", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "twitter.%", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "apple.%", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "paypal.%", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "microsoft.%", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "github.%", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "openid_connect.%", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "saml.authentication_request_signed", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "saml.idp_entity_id", fmt.Sprintf("idp:%s-1", name)),
			resource.TestCheckResourceAttr(resourceFullName, "saml.sp_entity_id", fmt.Sprintf("sp:%s-1", name)),
			resource.TestCheckResourceAttr(resourceFullName, "saml.idp_verification.certificates.#", "1"),
			resource.TestMatchResourceAttr(resourceFullName, "saml.idp_verification.certificates.0.id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckNoResourceAttr(resourceFullName, "saml.sp_signing"),
			resource.TestCheckResourceAttr(resourceFullName, "saml.sso_binding", "HTTP_REDIRECT"),
			resource.TestCheckResourceAttr(resourceFullName, "saml.sso_endpoint", "https://www.pingidentity.com/sso1"),
			resource.TestCheckResourceAttr(resourceFullName, "saml.slo_binding", "HTTP_POST"),
			resource.TestCheckNoResourceAttr(resourceFullName, "saml.slo_endpoint"),
			resource.TestCheckNoResourceAttr(resourceFullName, "saml.slo_response_endpoint"),
			resource.TestCheckNoResourceAttr(resourceFullName, "saml.slo_window"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
			acctest.PreCheckPEMCert(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.IdentityProvider_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			fullStep,
			{
				Config:  testAccIdentityProviderConfig_SAMLFull(resourceName, name, pem_cert),
				Destroy: true,
			},
			// Minimal
			minimalStep,
			{
				Config:  testAccIdentityProviderConfig_SAMLMinimal(resourceName, name, pem_cert),
				Destroy: true,
			},
			// Change
			fullStep,
			minimalStep,
			fullStep,
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
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.IdentityProvider_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProviderConfig_Apple1(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "facebook.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "google.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin_oidc.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.client_id", "dummyclientid1"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.client_secret_signing_key", "-----BEGIN PRIVATE KEY-----dummyclientsecretsigningkey1-----END PRIVATE KEY-----"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.key_id", "dummykeyi1"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.team_id", "dummyteam1"),
					// resource.TestMatchResourceAttr(resourceFullName, "apple.callback_url", regexp.MustCompile(`^https:\/\/auth\.pingone\.(?:eu|com|asia|ca)\/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}\/rp\/callback\/apple$`)),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "github.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.%", "0"),
				),
			},
			{
				Config: testAccIdentityProviderConfig_Github2(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "facebook.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "google.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin_oidc.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "github.client_id", "dummyclientid2"),
					resource.TestCheckResourceAttr(resourceFullName, "github.client_secret", "dummyclientsecret2"),
					// resource.TestMatchResourceAttr(resourceFullName, "github.callback_url", regexp.MustCompile(`^https:\/\/auth\.pingone\.(?:eu|com|asia|ca)\/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}\/rp\/callback\/github$`)),
					resource.TestCheckResourceAttr(resourceFullName, "openid_connect.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "saml.%", "0"),
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
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.IdentityProvider_CheckDestroy,
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

func testAccIdentityProviderConfig_NewEnv(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s
resource "pingone_identity_provider" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s"

  google = {
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

  icon = {
    id   = pingone_image.%[2]s.id
    href = pingone_image.%[2]s.uploaded_image.href
  }

  login_button_icon = {
    id   = pingone_image.%[2]s.id
    href = pingone_image.%[2]s.uploaded_image.href
  }

  google = {
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

  google = {
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

  facebook = {
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

  facebook = {
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

  google = {
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

  google = {
    client_id     = "dummyclientid2"
    client_secret = "dummyclientsecret2"
  }
}
		`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccIdentityProviderConfig_LinkedIn1(resourceName, name, linkedInType string) string {
	return fmt.Sprintf(`
		%[1]s
resource "pingone_identity_provider" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  %[4]s = {
    client_id     = "dummyclientid1"
    client_secret = "dummyclientsecret1"
  }
}
		`, acctest.GenericSandboxEnvironment(), resourceName, name, linkedInType)
}

func testAccIdentityProviderConfig_LinkedIn2(resourceName, name, linkedInType string) string {
	return fmt.Sprintf(`
		%[1]s
resource "pingone_identity_provider" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  %[4]s = {
    client_id     = "dummyclientid2"
    client_secret = "dummyclientsecret2"
  }
}
		`, acctest.GenericSandboxEnvironment(), resourceName, name, linkedInType)
}

func testAccIdentityProviderConfig_Yahoo1(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s
resource "pingone_identity_provider" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  yahoo = {
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

  yahoo = {
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

  amazon = {
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

  amazon = {
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

  twitter = {
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

  twitter = {
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

  apple = {
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

  apple = {
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

  paypal = {
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

  paypal = {
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

  microsoft = {
    client_id     = "dummyclientid1"
    client_secret = "dummyclientsecret1"
    tenant_id     = "dummytenantid1"
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

  microsoft = {
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

  github = {
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

  github = {
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

  icon = {
    id   = pingone_image.%[2]s.id
    href = pingone_image.%[2]s.uploaded_image.href
  }

  login_button_icon = {
    id   = pingone_image.%[2]s.id
    href = pingone_image.%[2]s.uploaded_image.href
  }

  openid_connect = {
    authorization_endpoint     = "https://www.pingidentity.com/authz"
    client_id                  = "dummyclientid1"
    client_secret              = "dummyclientsecret1"
    discovery_endpoint         = "https://www.pingidentity.com/discovery"
    issuer                     = "https://www.pingidentity.com/issuer"
    pkce_method                = "S256"
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

  openid_connect = {
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

func testAccIdentityProviderConfig_SAMLFull(resourceName, name, pem string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_key" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name                = "%[3]s"
  algorithm           = "EC"
  key_length          = 256
  signature_algorithm = "SHA384withECDSA"
  subject_dn          = "CN=%[3]s, OU=Ping Identity, O=Ping Identity, L=, ST=, C=US"
  usage_type          = "SIGNING"
  validity_period     = 365
}

resource "pingone_certificate" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  pem_file = <<EOT
%[4]s
EOT

  usage_type = "SIGNING"
}

resource "pingone_identity_provider" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  saml = {
    authentication_request_signed = true
    idp_entity_id                 = "idp:%[3]s"
    sp_entity_id                  = "sp:%[3]s"
    idp_verification = {
      certificates = [
        {
          id = pingone_certificate.%[2]s.id
        }
      ]
    }
    sp_signing = {
      key = {
        id = pingone_key.%[2]s.id
      }
      algorithm = "SHA512withRSA"
    }
    sso_binding           = "HTTP_POST"
    sso_endpoint          = "https://www.pingidentity.com/sso"
    slo_binding           = "HTTP_REDIRECT"
    slo_endpoint          = "https://dummy-slo-endpoint.pingidentity.com"
    slo_response_endpoint = "https://dummy-slo-response-endpoint.pingidentity.com"
    slo_window            = 1
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name, pem)
}

func testAccIdentityProviderConfig_SAMLMinimal(resourceName, name, pem string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_key" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name                = "%[3]s"
  algorithm           = "EC"
  key_length          = 256
  signature_algorithm = "SHA384withECDSA"
  subject_dn          = "CN=%[3]s, OU=Ping Identity, O=Ping Identity, L=, ST=, C=US"
  usage_type          = "SIGNING"
  validity_period     = 365
}

resource "pingone_certificate" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  pem_file = <<EOT
%[4]s
EOT

  usage_type = "SIGNING"
}

resource "pingone_identity_provider" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s1"

  saml = {
    idp_entity_id = "idp:%[3]s-1"
    sp_entity_id  = "sp:%[3]s-1"
    idp_verification = {
      certificates = [
        {
          id = pingone_certificate.%[2]s.id
        }
      ]
    }
    sso_binding  = "HTTP_REDIRECT"
    sso_endpoint = "https://www.pingidentity.com/sso1"
  }
}
		`, acctest.GenericSandboxEnvironment(), resourceName, name, pem)
}
