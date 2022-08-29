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
	pingone "github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
)

func testAccCheckIdentityProviderAttributeDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, pingone.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_identity_provider_attribute" {
			continue
		}

		_, rEnv, err := apiClient.EnvironmentsApi.ReadOneEnvironment(ctx, rs.Primary.Attributes["environment_id"]).Execute()

		if rEnv.StatusCode == 404 {
			continue
		}

		if err != nil {
			return err
		}

		body, r, err := apiClient.IdentityProviderAttributesApi.ReadOneIdentityProviderAttribute(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["identity_provider_id"], rs.Primary.ID).Execute()

		if r.StatusCode == 404 {
			continue
		}

		if err != nil {
			tflog.Error(ctx, fmt.Sprintf("Error: %v", body))
			return err
		}

		return fmt.Errorf("PingOne Identity Provider attribute %s still exists", rs.Primary.ID)
	}

	return nil
}

func TestAccIdentityProviderAttribute_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_identity_provider_attribute.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckIdentityProviderAttributeDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProviderAttributeConfig_Full(environmentName, resourceName, licenseID, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, "id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "environment_id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "identity_provider_id"),
					resource.TestCheckResourceAttr(resourceFullName, "name", "email"),
					resource.TestCheckResourceAttr(resourceFullName, "update", "EMPTY_ONLY"),
					resource.TestCheckResourceAttr(resourceFullName, "value", "${providerAttributes.emailAddress.value}"),
					resource.TestCheckResourceAttr(resourceFullName, "mapping_type", "CUSTOM"),
				),
			},
		},
	})
}

func TestAccIdentityProviderAttribute_Minimal(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_identity_provider_attribute.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckIdentityProviderAttributeDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProviderAttributeConfig_Minimal(environmentName, resourceName, licenseID, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, "id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "environment_id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "identity_provider_id"),
					resource.TestCheckResourceAttr(resourceFullName, "name", "email"),
					resource.TestCheckResourceAttr(resourceFullName, "update", "ALWAYS"),
					resource.TestCheckResourceAttr(resourceFullName, "value", "${providerAttributes.emailAddress.value}"),
					resource.TestCheckResourceAttr(resourceFullName, "mapping_type", "CUSTOM"),
				),
			},
		},
	})
}

func TestAccIdentityProviderAttribute_Expression(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_identity_provider_attribute.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckIdentityProviderAttributeDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProviderAttributeConfig_Expression(environmentName, resourceName, licenseID, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, "id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "environment_id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "identity_provider_id"),
					resource.TestCheckResourceAttr(resourceFullName, "name", "name.given"),
					resource.TestCheckResourceAttr(resourceFullName, "update", "ALWAYS"),
					resource.TestCheckResourceAttr(resourceFullName, "value", "${providerAttributes.name.givenName}"),
					resource.TestCheckResourceAttr(resourceFullName, "mapping_type", "CUSTOM"),
				),
			},
		},
	})
}

func TestAccIdentityProviderAttribute_Change(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_identity_provider_attribute.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckIdentityProviderAttributeDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProviderAttributeConfig_Minimal(environmentName, resourceName, licenseID, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "name", "email"),
					resource.TestCheckResourceAttr(resourceFullName, "update", "ALWAYS"),
					resource.TestCheckResourceAttr(resourceFullName, "value", "${providerAttributes.emailAddress.value}"),
				),
			},
			{
				Config: testAccIdentityProviderAttributeConfig_Expression(environmentName, resourceName, licenseID, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "name", "name.given"),
					resource.TestCheckResourceAttr(resourceFullName, "update", "ALWAYS"),
					resource.TestCheckResourceAttr(resourceFullName, "value", "${providerAttributes.name.givenName}"),
				),
			},
			{
				Config: testAccIdentityProviderAttributeConfig_Full(environmentName, resourceName, licenseID, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "name", "email"),
					resource.TestCheckResourceAttr(resourceFullName, "update", "EMPTY_ONLY"),
					resource.TestCheckResourceAttr(resourceFullName, "value", "${providerAttributes.emailAddress.value}"),
				),
			},
			{
				Config:      testAccIdentityProviderAttributeConfig_ReservedAttributeName(environmentName, resourceName, licenseID, name),
				ExpectError: regexp.MustCompile(`expected name to not be any of \[[a-zA-Z ]*\], got [a-zA-Z]*`),
			},
			{
				Config: testAccIdentityProviderAttributeConfig_Minimal(environmentName, resourceName, licenseID, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "name", "email"),
					resource.TestCheckResourceAttr(resourceFullName, "update", "ALWAYS"),
					resource.TestCheckResourceAttr(resourceFullName, "value", "${providerAttributes.emailAddress.value}"),
				),
			},
		},
	})
}

func TestAccIdentityProviderAttribute_ReservedAttributeName(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckIdentityProviderAttributeDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccIdentityProviderAttributeConfig_ReservedAttributeName(environmentName, resourceName, licenseID, name),
				ExpectError: regexp.MustCompile(`expected name to not be any of \[[a-zA-Z ]*\], got [a-zA-Z]*`),
			},
		},
	})
}

func testAccIdentityProviderAttributeConfig_Full(environmentName, resourceName, licenseID, name string) string {
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

		resource "pingone_identity_provider_attribute" "%[3]s" {
			environment_id = "${pingone_environment.%[2]s.id}"
			identity_provider_id = "${pingone_identity_provider.%[3]s.id}"
			
			name 		= "email"
			update 		= "EMPTY_ONLY"
			value		= "$${providerAttributes.emailAddress.value}"
		}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccIdentityProviderAttributeConfig_Minimal(environmentName, resourceName, licenseID, name string) string {
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

		resource "pingone_identity_provider_attribute" "%[3]s" {
			environment_id = "${pingone_environment.%[2]s.id}"
			identity_provider_id = "${pingone_identity_provider.%[3]s.id}"
			
			name 		= "email"
			update 		= "ALWAYS"
			value		= "$${providerAttributes.emailAddress.value}"
		}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccIdentityProviderAttributeConfig_Expression(environmentName, resourceName, licenseID, name string) string {
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

		resource "pingone_identity_provider_attribute" "%[3]s" {
			environment_id = "${pingone_environment.%[2]s.id}"
			identity_provider_id = "${pingone_identity_provider.%[3]s.id}"
			
			name 		= "name.given"
			update 		= "ALWAYS"
			value		= "$${providerAttributes.name.givenName}"
		}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccIdentityProviderAttributeConfig_ReservedAttributeName(environmentName, resourceName, licenseID, name string) string {
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

		resource "pingone_identity_provider_attribute" "%[3]s" {
			environment_id = "${pingone_environment.%[2]s.id}"
			identity_provider_id = "${pingone_identity_provider.%[3]s.id}"
			
			name 		= "account"
			update 		= "ALWAYS"
			value		= "$${'test'}"
		}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}
