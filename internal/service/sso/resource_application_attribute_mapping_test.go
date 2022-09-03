package sso_test

import (
	"context"
	"fmt"
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

		body, r, err := apiClient.ApplicationAttributeMappingApi.ReadOneApplicationAttributeMapping(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["application_id"], rs.Primary.ID).Execute()

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

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckApplicationAttributeMappingDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationAttributeMappingConfig_Full(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
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

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckApplicationAttributeMappingDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationAttributeMappingConfig_Minimal(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
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

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckApplicationAttributeMappingDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationAttributeMappingConfig_Expression(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
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

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckApplicationAttributeMappingDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationAttributeMappingConfig_Minimal(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "name", "email"),
					resource.TestCheckResourceAttr(resourceFullName, "required", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "value", "${user.email}"),
				),
			},
			{
				Config: testAccApplicationAttributeMappingConfig_Expression(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "name", "full_name"),
					resource.TestCheckResourceAttr(resourceFullName, "required", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "value", "${user.name.given + ', ' + user.name.family}"),
				),
			},
			{
				Config: testAccApplicationAttributeMappingConfig_Full(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "name", "email"),
					resource.TestCheckResourceAttr(resourceFullName, "required", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "value", "${user.email}"),
				),
			},
			{
				Config:      testAccApplicationAttributeMappingConfig_ReservedAttributeName(resourceName, name),
				ExpectError: regexp.MustCompile("Attribute name '[a-zA-Z]*' is not valid for the '[A-Z_]*' application"),
			},
			{
				Config: testAccApplicationAttributeMappingConfig_Minimal(resourceName, name),
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

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckApplicationAttributeMappingDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccApplicationAttributeMappingConfig_ReservedAttributeName(resourceName, name),
				ExpectError: regexp.MustCompile("Attribute name '[a-zA-Z]*' is not valid for the '[A-Z_]*' application"),
			},
		},
	})
}

func testAccApplicationAttributeMappingConfig_Full(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

		resource "pingone_application" "%[2]s" {
			environment_id  = "${data.pingone_environment.general_test.id}"
			name 			= "%[3]s"
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
			environment_id = "${data.pingone_environment.general_test.id}"
			application_id = "${pingone_application.%[2]s.id}"
			
			name 		= "email"
			required 	= true
			value		= "$${user.email}"
		}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccApplicationAttributeMappingConfig_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

		resource "pingone_application" "%[2]s" {
			environment_id  = "${data.pingone_environment.general_test.id}"
			name 			= "%[3]s"
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
			environment_id = "${data.pingone_environment.general_test.id}"
			application_id = "${pingone_application.%[2]s.id}"
			
			name 		= "email"
			value		= "$${user.email}"
		}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccApplicationAttributeMappingConfig_Expression(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

		resource "pingone_application" "%[2]s" {
			environment_id  = "${data.pingone_environment.general_test.id}"
			name 			= "%[3]s"
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
			environment_id = "${data.pingone_environment.general_test.id}"
			application_id = "${pingone_application.%[2]s.id}"
			
			name 		= "full_name"
			value		= "$${user.name.given + ', ' + user.name.family}"
		}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccApplicationAttributeMappingConfig_ReservedAttributeName(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

		resource "pingone_application" "%[2]s" {
			environment_id  = "${data.pingone_environment.general_test.id}"
			name 			= "%[3]s"
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
			environment_id = "${data.pingone_environment.general_test.id}"
			application_id = "${pingone_application.%[2]s.id}"
			
			name 		= "aud"
			value		= "$${'test'}"
		}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
