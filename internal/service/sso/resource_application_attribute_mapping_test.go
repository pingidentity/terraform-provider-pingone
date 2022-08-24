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

func testAccCheckApplicationAttributeMappingDestroy(s *terraform.State) error {
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
		if rs.Type != "pingone_application_attribute_mapping" {
			continue
		}

		_, rEnv, err := apiClient.EnvironmentsApi.ReadOneEnvironment(ctx, rs.Primary.Attributes["environment_id"]).Execute()

		if rEnv.StatusCode == 404 {
			continue
		}

		if err != nil {
			return err
		}

		body, r, err := apiClient.ApplicationsApplicationAttributeMappingApi.ReadOneApplicationAttributeMapping(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["application_id"], rs.Primary.ID).Execute()

		if r.StatusCode == 404 {
			continue
		}

		if err != nil {
			tflog.Error(ctx, fmt.Sprintf("Error: %v", body))
			return err
		}

		return fmt.Errorf("PingOne Application Attribute Mapping %s still exists", rs.Primary.ID)
	}

	return nil
}

func TestAccApplicationAttributeMapping_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application_attribute_mapping.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckApplicationAttributeMappingDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationAttributeMappingConfig_Full(environmentName, resourceName, licenseID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, "id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "environment_id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "application_id"),
					resource.TestCheckResourceAttr(resourceFullName, "name", "email"),
					resource.TestCheckResourceAttr(resourceFullName, "required", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "value", "${user.email}"),
					resource.TestCheckResourceAttr(resourceFullName, "mapping_type", "CUSTOM"),
				),
			},
		},
	})
}

func TestAccApplicationAttributeMapping_Minimal(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application_attribute_mapping.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckApplicationAttributeMappingDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationAttributeMappingConfig_Minimal(environmentName, resourceName, licenseID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, "id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "environment_id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "application_id"),
					resource.TestCheckResourceAttr(resourceFullName, "name", "email"),
					resource.TestCheckResourceAttr(resourceFullName, "required", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "value", "${user.email}"),
					resource.TestCheckResourceAttr(resourceFullName, "mapping_type", "CUSTOM"),
				),
			},
		},
	})
}

func TestAccApplicationAttributeMapping_Expression(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application_attribute_mapping.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckApplicationAttributeMappingDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationAttributeMappingConfig_Expression(environmentName, resourceName, licenseID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, "id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "environment_id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "application_id"),
					resource.TestCheckResourceAttr(resourceFullName, "name", "full_name"),
					resource.TestCheckResourceAttr(resourceFullName, "required", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "value", "${user.name.given + ', ' + user.name.family}"),
					resource.TestCheckResourceAttr(resourceFullName, "mapping_type", "CUSTOM"),
				),
			},
		},
	})
}

func TestAccApplicationAttributeMapping_Change(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application_attribute_mapping.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckApplicationAttributeMappingDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationAttributeMappingConfig_Minimal(environmentName, resourceName, licenseID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "name", "email"),
					resource.TestCheckResourceAttr(resourceFullName, "required", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "value", "${user.email}"),
				),
			},
			{
				Config: testAccApplicationAttributeMappingConfig_Expression(environmentName, resourceName, licenseID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "name", "full_name"),
					resource.TestCheckResourceAttr(resourceFullName, "required", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "value", "${user.name.given + ', ' + user.name.family}"),
				),
			},
			{
				Config: testAccApplicationAttributeMappingConfig_Full(environmentName, resourceName, licenseID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "name", "email"),
					resource.TestCheckResourceAttr(resourceFullName, "required", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "value", "${user.email}"),
				),
			},
			{
				Config:      testAccApplicationAttributeMappingConfig_ReservedAttributeName(environmentName, resourceName, licenseID),
				ExpectError: regexp.MustCompile("Attribute name '[a-zA-Z]*' is not valid for the '[A-Z_]*' application"),
			},
			{
				Config: testAccApplicationAttributeMappingConfig_Minimal(environmentName, resourceName, licenseID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "name", "email"),
					resource.TestCheckResourceAttr(resourceFullName, "required", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "value", "${user.email}"),
				),
			},
		},
	})
}

func TestAccApplicationAttributeMapping_ReservedAttributeName(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckApplicationAttributeMappingDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccApplicationAttributeMappingConfig_ReservedAttributeName(environmentName, resourceName, licenseID),
				ExpectError: regexp.MustCompile("Attribute name '[a-zA-Z]*' is not valid for the '[A-Z_]*' application"),
			},
		},
	})
}

func testAccApplicationAttributeMappingConfig_Full(environmentName, resourceName, licenseID string) string {
	return fmt.Sprintf(`
		resource "pingone_environment" "%[1]s" {
			name = "%[1]s"
			type = "SANDBOX"
			license_id = "%[3]s"
			default_population {}
			service {}
		}

		resource "pingone_application" "%[2]s" {
			environment_id  = "${pingone_environment.%[1]s.id}"
			name 			= "%[2]s"
			enabled 		= true
		  
			oidc_options {
				type                        = "SINGLE_PAGE_APP"
				grant_types                 = ["AUTHORIZATION_CODE"]
				response_types              = ["CODE"]
				pkce_enforcement            = "S256_REQUIRED"
				token_endpoint_authn_method = "NONE"
				redirect_uris               = ["https://www.pingidentity.com"]
			}
		}

		resource "pingone_application_attribute_mapping" "%[2]s" {
			environment_id = "${pingone_environment.%[1]s.id}"
			application_id = "${pingone_application.%[2]s.id}"
			
			name 		= "email"
			required 	= true
			value		= "$${user.email}"
		}`, environmentName, resourceName, licenseID)
}

func testAccApplicationAttributeMappingConfig_Minimal(environmentName, resourceName, licenseID string) string {
	return fmt.Sprintf(`
		resource "pingone_environment" "%[1]s" {
			name = "%[1]s"
			type = "SANDBOX"
			license_id = "%[3]s"
			default_population {}
			service {}
		}

		resource "pingone_application" "%[2]s" {
			environment_id  = "${pingone_environment.%[1]s.id}"
			name 			= "%[2]s"
			enabled 		= true
		  
			oidc_options {
				type                        = "SINGLE_PAGE_APP"
				grant_types                 = ["AUTHORIZATION_CODE"]
				response_types              = ["CODE"]
				pkce_enforcement            = "S256_REQUIRED"
				token_endpoint_authn_method = "NONE"
				redirect_uris               = ["https://www.pingidentity.com"]
			}
		}

		resource "pingone_application_attribute_mapping" "%[2]s" {
			environment_id = "${pingone_environment.%[1]s.id}"
			application_id = "${pingone_application.%[2]s.id}"
			
			name 		= "email"
			value		= "$${user.email}"
		}`, environmentName, resourceName, licenseID)
}

func testAccApplicationAttributeMappingConfig_Expression(environmentName, resourceName, licenseID string) string {
	return fmt.Sprintf(`
		resource "pingone_environment" "%[1]s" {
			name = "%[1]s"
			type = "SANDBOX"
			license_id = "%[3]s"
			default_population {}
			service {}
		}

		resource "pingone_application" "%[2]s" {
			environment_id  = "${pingone_environment.%[1]s.id}"
			name 			= "%[2]s"
			enabled 		= true
		  
			oidc_options {
				type                        = "SINGLE_PAGE_APP"
				grant_types                 = ["AUTHORIZATION_CODE"]
				response_types              = ["CODE"]
				pkce_enforcement            = "S256_REQUIRED"
				token_endpoint_authn_method = "NONE"
				redirect_uris               = ["https://www.pingidentity.com"]
			}
		}

		resource "pingone_application_attribute_mapping" "%[2]s" {
			environment_id = "${pingone_environment.%[1]s.id}"
			application_id = "${pingone_application.%[2]s.id}"
			
			name 		= "full_name"
			value		= "$${user.name.given + ', ' + user.name.family}"
		}`, environmentName, resourceName, licenseID)
}

func testAccApplicationAttributeMappingConfig_ReservedAttributeName(environmentName, resourceName, licenseID string) string {
	return fmt.Sprintf(`
		resource "pingone_environment" "%[1]s" {
			name = "%[1]s"
			type = "SANDBOX"
			license_id = "%[3]s"
			default_population {}
			service {}
		}

		resource "pingone_application" "%[2]s" {
			environment_id  = "${pingone_environment.%[1]s.id}"
			name 			= "%[2]s"
			enabled 		= true
		  
			oidc_options {
				type                        = "SINGLE_PAGE_APP"
				grant_types                 = ["AUTHORIZATION_CODE"]
				response_types              = ["CODE"]
				pkce_enforcement            = "S256_REQUIRED"
				token_endpoint_authn_method = "NONE"
				redirect_uris               = ["https://www.pingidentity.com"]
			}
		}

		resource "pingone_application_attribute_mapping" "%[2]s" {
			environment_id = "${pingone_environment.%[1]s.id}"
			application_id = "${pingone_application.%[2]s.id}"
			
			name 		= "aud"
			value		= "$${'test'}"
		}`, environmentName, resourceName, licenseID)
}
