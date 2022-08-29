package sso_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
)

func testAccCheckIdentityProviderDestroy(s *terraform.State) error {
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
		if rs.Type != "pingone_identity_provider" {
			continue
		}

		_, rEnv, err := apiClient.EnvironmentsApi.ReadOneEnvironment(ctx, rs.Primary.Attributes["environment_id"]).Execute()

		if rEnv.StatusCode == 404 {
			continue
		}

		if err != nil {
			return err
		}

		body, r, err := apiClient.IdentityProvidersApi.ReadOneIdentityProvider(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.ID).Execute()

		if r.StatusCode == 404 {
			continue
		}

		if err != nil {
			tflog.Error(ctx, fmt.Sprintf("Error: %v", body))
			return err
		}

		return fmt.Errorf("PingOne Identity Provider Instance %s still exists", rs.Primary.ID)
	}

	return nil
}

func TestAccIdentityProvider_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_identity_provider.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckIdentityProviderDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProviderConfig_Full(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", "My test identity provider"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestMatchResourceAttr(resourceFullName, "registration_population_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "login_button_icon.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "icon.#", "0"),
				),
			},
		},
	})
}

func TestAccIdentityProvider_Minimal(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_identity_provider.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckIdentityProviderDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProviderConfig_Minimal(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
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

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckIdentityProviderDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProviderConfig_Full(environmentName, licenseID, resourceName, fmt.Sprintf("%s 1", name)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "name", fmt.Sprintf("%s 1", name)),
					resource.TestCheckResourceAttr(resourceFullName, "description", "My test identity provider"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestMatchResourceAttr(resourceFullName, "registration_population_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "login_button_icon.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "icon.#", "0"),
				),
			},
			{
				Config: testAccIdentityProviderConfig_Minimal(environmentName, licenseID, resourceName, fmt.Sprintf("%s 2", name)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "name", fmt.Sprintf("%s 2", name)),
					resource.TestCheckResourceAttr(resourceFullName, "description", ""),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "registration_population_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "login_button_icon.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "icon.#", "0"),
				),
			},
			{
				Config: testAccIdentityProviderConfig_Full(environmentName, licenseID, resourceName, fmt.Sprintf("%s 1", name)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "name", fmt.Sprintf("%s 1", name)),
					resource.TestCheckResourceAttr(resourceFullName, "description", "My test identity provider"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestMatchResourceAttr(resourceFullName, "registration_population_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "login_button_icon.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "icon.#", "0"),
				),
			},
		},
	})
}

func TestAccIdentityProvider_Facebook(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_identity_provider.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckIdentityProviderDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProviderConfig_Facebook1(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "facebook.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "facebook.0.app_id", "dummyappid1"),
					resource.TestCheckResourceAttr(resourceFullName, "facebook.0.app_secret", "dummyappsecret1"),
					resource.TestCheckResourceAttr(resourceFullName, "google.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "github.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "generic_oidc.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "generic_saml.#", "0"),
				),
			},
			{
				Config: testAccIdentityProviderConfig_Facebook2(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "facebook.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "facebook.0.app_id", "dummyappid2"),
					resource.TestCheckResourceAttr(resourceFullName, "facebook.0.app_secret", "dummyappsecret2"),
					resource.TestCheckResourceAttr(resourceFullName, "google.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "github.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "generic_oidc.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "generic_saml.#", "0"),
				),
			},
		},
	})
}

func TestAccIdentityProvider_Google(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_identity_provider.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckIdentityProviderDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProviderConfig_Google1(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "facebook.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "google.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "google.0.client_id", "dummyclientid1"),
					resource.TestCheckResourceAttr(resourceFullName, "google.0.client_secret", "dummyclientsecret1"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "github.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "generic_oidc.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "generic_saml.#", "0"),
				),
			},
			{
				Config: testAccIdentityProviderConfig_Google2(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "facebook.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "google.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "google.0.client_id", "dummyclientid2"),
					resource.TestCheckResourceAttr(resourceFullName, "google.0.client_secret", "dummyclientsecret2"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "github.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "generic_oidc.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "generic_saml.#", "0"),
				),
			},
		},
	})
}

func TestAccIdentityProvider_LinkedIn(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_identity_provider.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckIdentityProviderDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProviderConfig_LinkedIn1(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "facebook.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "google.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.0.client_id", "dummyclientid1"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.0.client_secret", "dummyclientsecret1"),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "github.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "generic_oidc.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "generic_saml.#", "0"),
				),
			},
			{
				Config: testAccIdentityProviderConfig_LinkedIn2(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "facebook.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "google.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.0.client_id", "dummyclientid2"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.0.client_secret", "dummyclientsecret2"),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "github.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "generic_oidc.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "generic_saml.#", "0"),
				),
			},
		},
	})
}

func TestAccIdentityProvider_Yahoo(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_identity_provider.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckIdentityProviderDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProviderConfig_Yahoo1(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "facebook.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "google.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.0.client_id", "dummyclientid1"),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.0.client_secret", "dummyclientsecret1"),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "github.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "generic_oidc.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "generic_saml.#", "0"),
				),
			},
			{
				Config: testAccIdentityProviderConfig_Yahoo2(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "facebook.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "google.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.0.client_id", "dummyclientid2"),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.0.client_secret", "dummyclientsecret2"),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "github.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "generic_oidc.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "generic_saml.#", "0"),
				),
			},
		},
	})
}

func TestAccIdentityProvider_Amazon(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_identity_provider.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckIdentityProviderDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProviderConfig_Amazon1(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "facebook.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "google.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.0.client_id", "dummyclientid1"),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.0.client_secret", "dummyclientsecret1"),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "github.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "generic_oidc.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "generic_saml.#", "0"),
				),
			},
			{
				Config: testAccIdentityProviderConfig_Amazon2(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "facebook.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "google.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.0.client_id", "dummyclientid2"),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.0.client_secret", "dummyclientsecret2"),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "github.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "generic_oidc.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "generic_saml.#", "0"),
				),
			},
		},
	})
}

func TestAccIdentityProvider_Twitter(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_identity_provider.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckIdentityProviderDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProviderConfig_Twitter1(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "facebook.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "google.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.0.client_id", "dummyclientid1"),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.0.client_secret", "dummyclientsecret1"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "github.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "generic_oidc.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "generic_saml.#", "0"),
				),
			},
			{
				Config: testAccIdentityProviderConfig_Twitter2(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "facebook.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "google.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "linkedin.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "yahoo.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "amazon.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.0.client_id", "dummyclientid2"),
					resource.TestCheckResourceAttr(resourceFullName, "twitter.0.client_secret", "dummyclientsecret2"),
					resource.TestCheckResourceAttr(resourceFullName, "apple.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "paypal.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "github.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "generic_oidc.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "generic_saml.#", "0"),
				),
			},
		},
	})
}

func TestAccIdentityProvider_Apple(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_identity_provider.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckIdentityProviderDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProviderConfig_Apple1(environmentName, licenseID, resourceName, name),
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
					resource.TestCheckResourceAttr(resourceFullName, "paypal.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "github.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "generic_oidc.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "generic_saml.#", "0"),
				),
			},
			{
				Config: testAccIdentityProviderConfig_Apple2(environmentName, licenseID, resourceName, name),
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
					resource.TestCheckResourceAttr(resourceFullName, "paypal.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "github.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "generic_oidc.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "generic_saml.#", "0"),
				),
			},
		},
	})
}

func TestAccIdentityProvider_Paypal(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_identity_provider.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckIdentityProviderDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProviderConfig_Paypal1(environmentName, licenseID, resourceName, name),
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
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "github.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "generic_oidc.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "generic_saml.#", "0"),
				),
			},
			{
				Config: testAccIdentityProviderConfig_Paypal2(environmentName, licenseID, resourceName, name),
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
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "github.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "generic_oidc.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "generic_saml.#", "0"),
				),
			},
		},
	})
}

func TestAccIdentityProvider_Microsoft(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_identity_provider.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckIdentityProviderDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProviderConfig_Microsoft1(environmentName, licenseID, resourceName, name),
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
					resource.TestCheckResourceAttr(resourceFullName, "github.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "generic_oidc.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "generic_saml.#", "0"),
				),
			},
			{
				Config: testAccIdentityProviderConfig_Microsoft2(environmentName, licenseID, resourceName, name),
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
					resource.TestCheckResourceAttr(resourceFullName, "github.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "generic_oidc.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "generic_saml.#", "0"),
				),
			},
		},
	})
}

func TestAccIdentityProvider_Github(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_identity_provider.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckIdentityProviderDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProviderConfig_Github1(environmentName, licenseID, resourceName, name),
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
					resource.TestCheckResourceAttr(resourceFullName, "generic_oidc.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "generic_saml.#", "0"),
				),
			},
			{
				Config: testAccIdentityProviderConfig_Github2(environmentName, licenseID, resourceName, name),
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
					resource.TestCheckResourceAttr(resourceFullName, "generic_oidc.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "generic_saml.#", "0"),
				),
			},
		},
	})
}

func TestAccIdentityProvider_OIDC(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_identity_provider.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckIdentityProviderDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProviderConfig_OIDCFull(environmentName, licenseID, resourceName, name),
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
					resource.TestCheckResourceAttr(resourceFullName, "generic_oidc.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "generic_oidc.0.authorization_endpoint", "https://www.pingidentity.com/authz"),
					resource.TestCheckResourceAttr(resourceFullName, "generic_oidc.0.client_id", "dummyclientid1"),
					resource.TestCheckResourceAttr(resourceFullName, "generic_oidc.0.client_secret", "dummyclientsecret1"),
					resource.TestCheckResourceAttr(resourceFullName, "generic_oidc.0.discovery_endpoint", "https://www.pingidentity.com/discovery"),
					resource.TestCheckResourceAttr(resourceFullName, "generic_oidc.0.issuer", "https://www.pingidentity.com/issuer"),
					resource.TestCheckResourceAttr(resourceFullName, "generic_oidc.0.jwks_endpoint", "https://www.pingidentity.com/jwks"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "generic_oidc.0.scopes.*", "openid"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "generic_oidc.0.scopes.*", "scope1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "generic_oidc.0.scopes.*", "scope2"),
					resource.TestCheckResourceAttr(resourceFullName, "generic_oidc.0.token_endpoint", "https://www.pingidentity.com/token"),
					resource.TestCheckResourceAttr(resourceFullName, "generic_oidc.0.token_endpoint_auth_method", "CLIENT_SECRET_POST"),
					resource.TestCheckResourceAttr(resourceFullName, "generic_oidc.0.userinfo_endpoint", "https://www.pingidentity.com/userinfo"),
					resource.TestCheckResourceAttr(resourceFullName, "generic_saml.#", "0"),
				),
			},
			{
				Config: testAccIdentityProviderConfig_OIDCMinimal(environmentName, licenseID, resourceName, name),
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
					resource.TestCheckResourceAttr(resourceFullName, "generic_oidc.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "generic_oidc.0.authorization_endpoint", "https://www.pingidentity.com/authz2"),
					resource.TestCheckResourceAttr(resourceFullName, "generic_oidc.0.client_id", "dummyclientid2"),
					resource.TestCheckResourceAttr(resourceFullName, "generic_oidc.0.client_secret", "dummyclientsecret2"),
					resource.TestCheckResourceAttr(resourceFullName, "generic_oidc.0.discovery_endpoint", ""),
					resource.TestCheckResourceAttr(resourceFullName, "generic_oidc.0.issuer", "https://www.pingidentity.com/issuer2"),
					resource.TestCheckResourceAttr(resourceFullName, "generic_oidc.0.jwks_endpoint", "https://www.pingidentity.com/jwks2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "generic_oidc.0.scopes.*", "openid"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "generic_oidc.0.scopes.*", "scope3"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "generic_oidc.0.scopes.*", "scope4"),
					resource.TestCheckResourceAttr(resourceFullName, "generic_oidc.0.token_endpoint", "https://www.pingidentity.com/token2"),
					resource.TestCheckResourceAttr(resourceFullName, "generic_oidc.0.token_endpoint_auth_method", "CLIENT_SECRET_BASIC"),
					resource.TestCheckResourceAttr(resourceFullName, "generic_oidc.0.userinfo_endpoint", ""),
					resource.TestCheckResourceAttr(resourceFullName, "generic_saml.#", "0"),
				),
			},
			{
				Config: testAccIdentityProviderConfig_OIDCFull(environmentName, licenseID, resourceName, name),
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
					resource.TestCheckResourceAttr(resourceFullName, "generic_oidc.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "generic_oidc.0.authorization_endpoint", "https://www.pingidentity.com/authz"),
					resource.TestCheckResourceAttr(resourceFullName, "generic_oidc.0.client_id", "dummyclientid1"),
					resource.TestCheckResourceAttr(resourceFullName, "generic_oidc.0.client_secret", "dummyclientsecret1"),
					resource.TestCheckResourceAttr(resourceFullName, "generic_oidc.0.discovery_endpoint", "https://www.pingidentity.com/discovery"),
					resource.TestCheckResourceAttr(resourceFullName, "generic_oidc.0.issuer", "https://www.pingidentity.com/issuer"),
					resource.TestCheckResourceAttr(resourceFullName, "generic_oidc.0.jwks_endpoint", "https://www.pingidentity.com/jwks"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "generic_oidc.0.scopes.*", "openid"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "generic_oidc.0.scopes.*", "scope1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "generic_oidc.0.scopes.*", "scope2"),
					resource.TestCheckResourceAttr(resourceFullName, "generic_oidc.0.token_endpoint", "https://www.pingidentity.com/token"),
					resource.TestCheckResourceAttr(resourceFullName, "generic_oidc.0.token_endpoint_auth_method", "CLIENT_SECRET_POST"),
					resource.TestCheckResourceAttr(resourceFullName, "generic_oidc.0.userinfo_endpoint", "https://www.pingidentity.com/userinfo"),
					resource.TestCheckResourceAttr(resourceFullName, "generic_saml.#", "0"),
				),
			},
		},
	})
}

// func TestAccIdentityProvider_SAML(t *testing.T) {
// 	t.Parallel()

// 	resourceName := acctest.ResourceNameGen()
// 	resourceFullName := fmt.Sprintf("pingone_identity_provider.%s", resourceName)

// 	environmentName := acctest.ResourceNameGenEnvironment()

// 	name := resourceName

// 	licenseID := os.Getenv("PINGONE_LICENSE_ID")

// 	resource.Test(t, resource.TestCase{
// 		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
// 		ProviderFactories: acctest.ProviderFactories,
// 		CheckDestroy:      testAccCheckIdentityProviderDestroy,
// 		ErrorCheck:        acctest.ErrorCheck(t),
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccIdentityProviderConfig_SAMLFull(environmentName, licenseID, resourceName, name),
// 				Check: resource.ComposeTestCheckFunc(
// 					resource.TestCheckResourceAttr(resourceFullName, "facebook.#", "0"),
// 					resource.TestCheckResourceAttr(resourceFullName, "google.#", "0"),
// 					resource.TestCheckResourceAttr(resourceFullName, "linkedin.#", "0"),
// 					resource.TestCheckResourceAttr(resourceFullName, "yahoo.#", "0"),
// 					resource.TestCheckResourceAttr(resourceFullName, "amazon.#", "0"),
// 					resource.TestCheckResourceAttr(resourceFullName, "twitter.#", "0"),
// 					resource.TestCheckResourceAttr(resourceFullName, "apple.#", "0"),
// 					resource.TestCheckResourceAttr(resourceFullName, "paypal.#", "0"),
// 					resource.TestCheckResourceAttr(resourceFullName, "microsoft.#", "0"),
// 					resource.TestCheckResourceAttr(resourceFullName, "github.#", "0"),
// 					resource.TestCheckResourceAttr(resourceFullName, "generic_oidc.#", "0"),
// 					resource.TestCheckResourceAttr(resourceFullName, "generic_saml.#", "1"),
// 					resource.TestCheckResourceAttr(resourceFullName, "generic_saml.0.authentication_request_signed", "true"),
// 					resource.TestCheckResourceAttr(resourceFullName, "generic_saml.0.idp_entity_id", "idp:entity"),
// 					resource.TestCheckResourceAttr(resourceFullName, "generic_saml.0.sp_entity_id", "sp:entity"),
// 					resource.TestCheckResourceAttr(resourceFullName, "generic_saml.0.idp_verification_certificate_ids", "https://www.pingidentity.com/discovery"),
// 					resource.TestMatchResourceAttr(resourceFullName, "generic_saml.0.sp_signing_key_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
// 					resource.TestCheckResourceAttr(resourceFullName, "generic_saml.0.sso_binding", "HTTP_POST"),
// 					resource.TestCheckResourceAttr(resourceFullName, "generic_saml.0.sso_endpoint", "https://www.pingidentity.com/sso"),
// 				),
// 			},
// 			{
// 				Config: testAccIdentityProviderConfig_SAMLMinimal(environmentName, licenseID, resourceName, name),
// 				Check: resource.ComposeTestCheckFunc(
// 					resource.TestCheckResourceAttr(resourceFullName, "facebook.#", "0"),
// 					resource.TestCheckResourceAttr(resourceFullName, "google.#", "0"),
// 					resource.TestCheckResourceAttr(resourceFullName, "linkedin.#", "0"),
// 					resource.TestCheckResourceAttr(resourceFullName, "yahoo.#", "0"),
// 					resource.TestCheckResourceAttr(resourceFullName, "amazon.#", "0"),
// 					resource.TestCheckResourceAttr(resourceFullName, "twitter.#", "0"),
// 					resource.TestCheckResourceAttr(resourceFullName, "apple.#", "0"),
// 					resource.TestCheckResourceAttr(resourceFullName, "paypal.#", "0"),
// 					resource.TestCheckResourceAttr(resourceFullName, "microsoft.#", "0"),
// 					resource.TestCheckResourceAttr(resourceFullName, "github.#", "0"),
// 					resource.TestCheckResourceAttr(resourceFullName, "generic_oidc.#", "0"),
// 					resource.TestCheckResourceAttr(resourceFullName, "generic_saml.#", "1"),
// 					resource.TestCheckResourceAttr(resourceFullName, "generic_saml.0.authentication_request_signed", ""),
// 					resource.TestCheckResourceAttr(resourceFullName, "generic_saml.0.idp_entity_id", "idp:entity"),
// 					resource.TestCheckResourceAttr(resourceFullName, "generic_saml.0.sp_entity_id", ""),
// 					resource.TestCheckResourceAttr(resourceFullName, "generic_saml.0.idp_verification_certificate_ids.#", "0"),
// 					resource.TestCheckResourceAttr(resourceFullName, "generic_saml.0.sp_signing_key_id", ""),
// 					resource.TestCheckResourceAttr(resourceFullName, "generic_saml.0.sso_binding", "HTTP_REDIRECT"),
// 					resource.TestCheckResourceAttr(resourceFullName, "generic_saml.0.sso_endpoint", "https://www.pingidentity.com/sso"),
// 				),
// 			},
// 			{
// 				Config: testAccIdentityProviderConfig_SAMLFull(environmentName, licenseID, resourceName, name),
// 				Check: resource.ComposeTestCheckFunc(
// 					resource.TestCheckResourceAttr(resourceFullName, "facebook.#", "0"),
// 					resource.TestCheckResourceAttr(resourceFullName, "google.#", "0"),
// 					resource.TestCheckResourceAttr(resourceFullName, "linkedin.#", "0"),
// 					resource.TestCheckResourceAttr(resourceFullName, "yahoo.#", "0"),
// 					resource.TestCheckResourceAttr(resourceFullName, "amazon.#", "0"),
// 					resource.TestCheckResourceAttr(resourceFullName, "twitter.#", "0"),
// 					resource.TestCheckResourceAttr(resourceFullName, "apple.#", "0"),
// 					resource.TestCheckResourceAttr(resourceFullName, "paypal.#", "0"),
// 					resource.TestCheckResourceAttr(resourceFullName, "microsoft.#", "0"),
// 					resource.TestCheckResourceAttr(resourceFullName, "github.#", "0"),
// 					resource.TestCheckResourceAttr(resourceFullName, "generic_oidc.#", "0"),
// 					resource.TestCheckResourceAttr(resourceFullName, "generic_saml.#", "1"),
// 					resource.TestCheckResourceAttr(resourceFullName, "generic_saml.0.authentication_request_signed", "true"),
// 					resource.TestCheckResourceAttr(resourceFullName, "generic_saml.0.idp_entity_id", "idp:entity"),
// 					resource.TestCheckResourceAttr(resourceFullName, "generic_saml.0.sp_entity_id", "sp:entity"),
// 					resource.TestCheckResourceAttr(resourceFullName, "generic_saml.0.idp_verification_certificate_ids", "https://www.pingidentity.com/discovery"),
// 					resource.TestMatchResourceAttr(resourceFullName, "generic_saml.0.sp_signing_key_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
// 					resource.TestCheckResourceAttr(resourceFullName, "generic_saml.0.sso_binding", "HTTP_POST"),
// 					resource.TestCheckResourceAttr(resourceFullName, "generic_saml.0.sso_endpoint", "https://www.pingidentity.com/sso"),
// 				),
// 			},
// 		},
// 	})
// }

func TestAccIdentityProvider_ChangeProvider(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_identity_provider.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckIdentityProviderDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProviderConfig_Apple1(environmentName, licenseID, resourceName, name),
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
					resource.TestCheckResourceAttr(resourceFullName, "paypal.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "microsoft.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "github.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "generic_oidc.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "generic_saml.#", "0"),
				),
			},
			{
				Config: testAccIdentityProviderConfig_Github2(environmentName, licenseID, resourceName, name),
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
					resource.TestCheckResourceAttr(resourceFullName, "generic_oidc.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "generic_saml.#", "0"),
				),
			},
		},
	})
}

func testAccIdentityProviderConfig_Full(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s
		resource "pingone_identity_provider" "%[3]s" {
			environment_id = "${pingone_environment.%[2]s.id}"
			name = "%[4]s"
			description = "My test identity provider"
			enabled = true
			registration_population_id = "${pingone_environment.%[2]s.default_population_id}"

			// icon {
			// 	id = "1"
			// 	href = "https://assets.pingone.com/ux/ui-library/4.18.0/images/logo-pingidentity.png"
			// }

			// login_button_icon {
				// 	id = "1"
				// 	href = "https://assets.pingone.com/ux/ui-library/4.18.0/images/logo-pingidentity.png"
				// }

			google {
				client_id = "testclientid"
				client_secret = "testclientsecret"
			}
		}
		`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccIdentityProviderConfig_Minimal(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s
		resource "pingone_identity_provider" "%[3]s" {
			environment_id = "${pingone_environment.%[2]s.id}"
			name = "%[4]s"
			
			google {
				client_id = "testclientid"
				client_secret = "testclientsecret"
			}
		}
		`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccIdentityProviderConfig_Facebook1(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s
		resource "pingone_identity_provider" "%[3]s" {
			environment_id = "${pingone_environment.%[2]s.id}"
			name = "%[4]s"
			
			facebook {
				app_id = "dummyappid1"
				app_secret = "dummyappsecret1"
			}
		}
		`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccIdentityProviderConfig_Facebook2(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s
		resource "pingone_identity_provider" "%[3]s" {
			environment_id = "${pingone_environment.%[2]s.id}"
			name = "%[4]s"
			
			facebook {
				app_id = "dummyappid2"
				app_secret = "dummyappsecret2"
			}
		}
		`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccIdentityProviderConfig_Google1(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s
		resource "pingone_identity_provider" "%[3]s" {
			environment_id = "${pingone_environment.%[2]s.id}"
			name = "%[4]s"
			
			google {
				client_id = "dummyclientid1"
				client_secret = "dummyclientsecret1"
			}
		}
		`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccIdentityProviderConfig_Google2(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s
		resource "pingone_identity_provider" "%[3]s" {
			environment_id = "${pingone_environment.%[2]s.id}"
			name = "%[4]s"
			
			google {
				client_id = "dummyclientid2"
				client_secret = "dummyclientsecret2"
			}
		}
		`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccIdentityProviderConfig_LinkedIn1(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s
		resource "pingone_identity_provider" "%[3]s" {
			environment_id = "${pingone_environment.%[2]s.id}"
			name = "%[4]s"
			
			linkedin {
				client_id = "dummyclientid1"
				client_secret = "dummyclientsecret1"
			}
		}
		`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccIdentityProviderConfig_LinkedIn2(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s
		resource "pingone_identity_provider" "%[3]s" {
			environment_id = "${pingone_environment.%[2]s.id}"
			name = "%[4]s"
			
			linkedin {
				client_id = "dummyclientid2"
				client_secret = "dummyclientsecret2"
			}
		}
		`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccIdentityProviderConfig_Yahoo1(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s
		resource "pingone_identity_provider" "%[3]s" {
			environment_id = "${pingone_environment.%[2]s.id}"
			name = "%[4]s"
			
			yahoo {
				client_id = "dummyclientid1"
				client_secret = "dummyclientsecret1"
			}
		}
		`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccIdentityProviderConfig_Yahoo2(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s
		resource "pingone_identity_provider" "%[3]s" {
			environment_id = "${pingone_environment.%[2]s.id}"
			name = "%[4]s"
			
			yahoo {
				client_id = "dummyclientid2"
				client_secret = "dummyclientsecret2"
			}
		}
		`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccIdentityProviderConfig_Amazon1(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s
		resource "pingone_identity_provider" "%[3]s" {
			environment_id = "${pingone_environment.%[2]s.id}"
			name = "%[4]s"
			
			amazon {
				client_id = "dummyclientid1"
				client_secret = "dummyclientsecret1"
			}
		}
		`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccIdentityProviderConfig_Amazon2(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s
		resource "pingone_identity_provider" "%[3]s" {
			environment_id = "${pingone_environment.%[2]s.id}"
			name = "%[4]s"
			
			amazon {
				client_id = "dummyclientid2"
				client_secret = "dummyclientsecret2"
			}
		}
		`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccIdentityProviderConfig_Twitter1(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s
		resource "pingone_identity_provider" "%[3]s" {
			environment_id = "${pingone_environment.%[2]s.id}"
			name = "%[4]s"
			
			twitter {
				client_id = "dummyclientid1"
				client_secret = "dummyclientsecret1"
			}
		}
		`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccIdentityProviderConfig_Twitter2(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s
		resource "pingone_identity_provider" "%[3]s" {
			environment_id = "${pingone_environment.%[2]s.id}"
			name = "%[4]s"
			
			twitter {
				client_id = "dummyclientid2"
				client_secret = "dummyclientsecret2"
			}
		}
		`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccIdentityProviderConfig_Apple1(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s
		resource "pingone_identity_provider" "%[3]s" {
			environment_id = "${pingone_environment.%[2]s.id}"
			name = "%[4]s"
			
			apple {
				client_id = "dummyclientid1"
				client_secret_signing_key = "-----BEGIN PRIVATE KEY-----dummyclientsecretsigningkey1-----END PRIVATE KEY-----"
				key_id = "dummykeyi1"
				team_id = "dummyteam1"
			}
		}
		`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccIdentityProviderConfig_Apple2(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s
		resource "pingone_identity_provider" "%[3]s" {
			environment_id = "${pingone_environment.%[2]s.id}"
			name = "%[4]s"
			
			apple {
				client_id = "dummyclientid2"
				client_secret_signing_key = "-----BEGIN PRIVATE KEY-----dummyclientsecretsigningkey2-----END PRIVATE KEY-----"
				key_id = "dummykeyi2"
				team_id = "dummyteam2"
			}
		}
		`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccIdentityProviderConfig_Paypal1(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s
		resource "pingone_identity_provider" "%[3]s" {
			environment_id = "${pingone_environment.%[2]s.id}"
			name = "%[4]s"
			
			paypal {
				client_id = "dummyclientid1"
				client_secret = "dummyclientsecret1"
				client_environment = "sandbox"
			}
		}
		`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccIdentityProviderConfig_Paypal2(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s
		resource "pingone_identity_provider" "%[3]s" {
			environment_id = "${pingone_environment.%[2]s.id}"
			name = "%[4]s"
			
			paypal {
				client_id = "dummyclientid2"
				client_secret = "dummyclientsecret2"
				client_environment = "live"
			}
		}
		`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccIdentityProviderConfig_Microsoft1(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s
		resource "pingone_identity_provider" "%[3]s" {
			environment_id = "${pingone_environment.%[2]s.id}"
			name = "%[4]s"
			
			microsoft {
				client_id = "dummyclientid1"
				client_secret = "dummyclientsecret1"
			}
		}
		`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccIdentityProviderConfig_Microsoft2(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s
		resource "pingone_identity_provider" "%[3]s" {
			environment_id = "${pingone_environment.%[2]s.id}"
			name = "%[4]s"
			
			microsoft {
				client_id = "dummyclientid2"
				client_secret = "dummyclientsecret2"
			}
		}
		`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccIdentityProviderConfig_Github1(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s
		resource "pingone_identity_provider" "%[3]s" {
			environment_id = "${pingone_environment.%[2]s.id}"
			name = "%[4]s"
			
			github {
				client_id = "dummyclientid1"
				client_secret = "dummyclientsecret1"
			}
		}
		`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccIdentityProviderConfig_Github2(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s
		resource "pingone_identity_provider" "%[3]s" {
			environment_id = "${pingone_environment.%[2]s.id}"
			name = "%[4]s"
			
			github {
				client_id = "dummyclientid2"
				client_secret = "dummyclientsecret2"
			}
		}
		`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccIdentityProviderConfig_OIDCFull(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s
		resource "pingone_identity_provider" "%[3]s" {
			environment_id = "${pingone_environment.%[2]s.id}"
			name = "%[4]s"
			
			generic_oidc {
				authorization_endpoint = "https://www.pingidentity.com/authz"
				client_id = "dummyclientid1"
				client_secret = "dummyclientsecret1"
				discovery_endpoint = "https://www.pingidentity.com/discovery"
				issuer = "https://www.pingidentity.com/issuer"
				jwks_endpoint = "https://www.pingidentity.com/jwks"
				scopes = ["openid", "scope1", "scope2"]
				token_endpoint = "https://www.pingidentity.com/token"
				token_endpoint_auth_method = "CLIENT_SECRET_POST"
				userinfo_endpoint = "https://www.pingidentity.com/userinfo"
			}
		}
		`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccIdentityProviderConfig_OIDCMinimal(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s
		resource "pingone_identity_provider" "%[3]s" {
			environment_id = "${pingone_environment.%[2]s.id}"
			name = "%[4]s"
			
			generic_oidc {
				authorization_endpoint = "https://www.pingidentity.com/authz2"
				client_id = "dummyclientid2"
				client_secret = "dummyclientsecret2"
				issuer = "https://www.pingidentity.com/issuer2"
				jwks_endpoint = "https://www.pingidentity.com/jwks2"
				scopes = ["openid", "scope3", "scope4"]
				token_endpoint = "https://www.pingidentity.com/token2"
			}
		}
		`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

// func testAccIdentityProviderConfig_SAMLFull(environmentName, licenseID, resourceName, name string) string {
// 	return fmt.Sprintf(`
// 		%[1]s
// 		resource "pingone_identity_provider" "%[3]s" {
// 			environment_id = "${pingone_environment.%[2]s.id}"
// 			name = "%[4]s"

// 			generic_saml {
// 				authentication_request_signed = true
// 				idp_entity_id = "idp:entity"
// 				sp_entity_id = "sp:entity"
// 				idp_verification_certificate_ids = []
// 				// sp_signing_key_id =
// 				sso_binding = "HTTP_POST"
// 				sso_endpoint = "https://www.pingidentity.com/sso"
// 			}
// 		}
// 		`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
// }

// func testAccIdentityProviderConfig_SAMLMinimal(environmentName, licenseID, resourceName, name string) string {
// 	return fmt.Sprintf(`
// 		%[1]s
// 		resource "pingone_identity_provider" "%[3]s" {
// 			environment_id = "${pingone_environment.%[2]s.id}"
// 			name = "%[4]s"

// 			generic_saml {
// 				idp_entity_id = "idp:entity"
// 				sso_binding = "HTTP_REDIRECT"
// 				sso_endpoint = "https://www.pingidentity.com/sso"
// 			}
// 		}
// 		`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
// }
