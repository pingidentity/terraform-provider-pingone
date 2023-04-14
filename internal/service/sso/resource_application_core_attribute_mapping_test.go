package sso_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func testAccCheckApplicationCoreAttributeMappingDestroy(s *terraform.State) error {
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
		if rs.Type != "pingone_application_core_attribute_mapping" {
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

		body, r, err := apiClient.ApplicationAttributeMappingApi.ReadOneApplicationAttributeMapping(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["application_id"], rs.Primary.ID).Execute()

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

		return fmt.Errorf("PingOne Application Attribute Mapping %s still exists", rs.Primary.ID)
	}

	return nil
}

func TestAccApplicationCoreAttributeMapping_OIDC(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application_core_attribute_mapping.%s", resourceName)

	name := resourceName

	fullStep := resource.TestStep{
		Config: testAccApplicationCoreAttributeMappingConfig_OIDC(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
			resource.TestMatchResourceAttr(resourceFullName, "application_id", verify.P1ResourceIDRegexp),
			resource.TestCheckResourceAttr(resourceFullName, "name", "sub"),
			resource.TestCheckResourceAttr(resourceFullName, "required", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "value", "${user.email}"),
			resource.TestCheckResourceAttr(resourceFullName, "mapping_type", "CORE"),
		),
	}

	coreAttrNameAppTypeStep := resource.TestStep{
		Config:      testAccApplicationCoreAttributeMappingConfig_OIDC_SAML_Name(resourceName, name),
		ExpectError: regexp.MustCompile(`Invalid parameter value - Not a core attribute`),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckApplicationCoreAttributeMappingDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			fullStep,
			{
				Config:  testAccApplicationCoreAttributeMappingConfig_OIDC(resourceName, name),
				Destroy: true,
			},
			// Allow core attribute name on other app types
			coreAttrNameAppTypeStep,
		},
	})
}

func TestAccApplicationCoreAttributeMapping_SAML(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application_core_attribute_mapping.%s", resourceName)

	name := resourceName

	fullStep := resource.TestStep{
		Config: testAccApplicationCoreAttributeMappingConfig_SAML_Full(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
			resource.TestMatchResourceAttr(resourceFullName, "application_id", verify.P1ResourceIDRegexp),
			resource.TestCheckResourceAttr(resourceFullName, "name", "saml_subject"),
			resource.TestCheckResourceAttr(resourceFullName, "required", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "value", "${user.email}"),
			resource.TestCheckResourceAttr(resourceFullName, "mapping_type", "CORE"),
			resource.TestCheckResourceAttr(resourceFullName, "saml_subject_nameformat", "urn:oasis:names:tc:SAML:2.0:attrname-format:uri"),
		),
	}

	minimalStep := resource.TestStep{
		Config: testAccApplicationCoreAttributeMappingConfig_SAML_Minimal(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
			resource.TestMatchResourceAttr(resourceFullName, "application_id", verify.P1ResourceIDRegexp),
			resource.TestCheckResourceAttr(resourceFullName, "name", "saml_subject"),
			resource.TestCheckResourceAttr(resourceFullName, "required", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "value", "${user.email}"),
			resource.TestCheckResourceAttr(resourceFullName, "mapping_type", "CORE"),
			resource.TestCheckNoResourceAttr(resourceFullName, "saml_subject_nameformat"),
		),
	}

	coreAttrNameAppTypeStep := resource.TestStep{
		Config:      testAccApplicationCoreAttributeMappingConfig_SAML_OIDC_Name(resourceName, name),
		ExpectError: regexp.MustCompile(`Invalid parameter value - Not a core attribute`),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckApplicationCoreAttributeMappingDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			fullStep,
			{
				Config:  testAccApplicationCoreAttributeMappingConfig_SAML_Full(resourceName, name),
				Destroy: true,
			},
			// Minimal
			minimalStep,
			{
				Config:  testAccApplicationCoreAttributeMappingConfig_SAML_Minimal(resourceName, name),
				Destroy: true,
			},
			// Change
			minimalStep,
			fullStep,
			minimalStep,
			{
				Config:  testAccApplicationCoreAttributeMappingConfig_SAML_Minimal(resourceName, name),
				Destroy: true,
			},
			// Allow core attribute name on other application types
			coreAttrNameAppTypeStep,
		},
	})
}

func TestAccApplicationCoreAttributeMapping_BadApplication(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckApplicationCoreAttributeMappingDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccApplicationCoreAttributeMappingConfig_BadApplication(resourceName, name),
				ExpectError: regexp.MustCompile("Invalid parameter value - Unmappable application type"),
			},
		},
	})
}

func TestAccApplicationCoreAttributeMapping_Expression(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application_core_attribute_mapping.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckApplicationCoreAttributeMappingDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationCoreAttributeMappingConfig_Expression(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "application_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "name", "sub"),
					resource.TestCheckResourceAttr(resourceFullName, "required", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "value", "${user.name.given + ', ' + user.name.family}"),
					resource.TestCheckResourceAttr(resourceFullName, "mapping_type", "CORE"),
				),
			},
		},
	})
}

func testAccApplicationCoreAttributeMappingConfig_OIDC(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true

  oidc_options {
    type                        = "SINGLE_PAGE_APP"
    grant_types                 = ["AUTHORIZATION_CODE"]
    response_types              = ["CODE"]
    pkce_enforcement            = "S256_REQUIRED"
    token_endpoint_authn_method = "NONE"
    redirect_uris               = ["https://www.pingidentity.com"]
  }
}

resource "pingone_application_core_attribute_mapping" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  application_id = pingone_application.%[2]s.id

  name  = "sub"
  value = "$${user.email}"
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccApplicationCoreAttributeMappingConfig_OIDC_SAML_Name(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true

  oidc_options {
    type                        = "SINGLE_PAGE_APP"
    grant_types                 = ["AUTHORIZATION_CODE"]
    response_types              = ["CODE"]
    pkce_enforcement            = "S256_REQUIRED"
    token_endpoint_authn_method = "NONE"
    redirect_uris               = ["https://www.pingidentity.com"]
  }
}

resource "pingone_application_core_attribute_mapping" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  application_id = pingone_application.%[2]s.id

  name  = "saml_subject"
  value = "$${user.email}"
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccApplicationCoreAttributeMappingConfig_SAML_Full(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_key" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name                = "%[3]s"
  algorithm           = "RSA"
  key_length          = 4096
  signature_algorithm = "SHA512withRSA"
  subject_dn          = "CN=%[3]s, OU=BX Retail, O=BX Retail, L=, ST=, C=US"
  usage_type          = "SIGNING"
  validity_period     = 365
}

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true

  saml_options {
    acs_urls           = ["https://my-saas-app.com"]
    assertion_duration = 3600
    sp_entity_id       = "sp:entity:%[3]s"

    idp_signing_key {
      key_id    = pingone_key.%[2]s.id
      algorithm = pingone_key.%[2]s.signature_algorithm
    }
  }
}

resource "pingone_application_core_attribute_mapping" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  application_id = pingone_application.%[2]s.id

  name  = "saml_subject"
  value = "$${user.email}"

  saml_subject_nameformat = "urn:oasis:names:tc:SAML:2.0:attrname-format:uri"

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccApplicationCoreAttributeMappingConfig_SAML_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_key" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name                = "%[3]s"
  algorithm           = "RSA"
  key_length          = 4096
  signature_algorithm = "SHA512withRSA"
  subject_dn          = "CN=%[3]s, OU=BX Retail, O=BX Retail, L=, ST=, C=US"
  usage_type          = "SIGNING"
  validity_period     = 365
}

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true

  saml_options {
    acs_urls           = ["https://my-saas-app.com"]
    assertion_duration = 3600
    sp_entity_id       = "sp:entity:%[3]s"

    idp_signing_key {
      key_id    = pingone_key.%[2]s.id
      algorithm = pingone_key.%[2]s.signature_algorithm
    }
  }
}

resource "pingone_application_core_attribute_mapping" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  application_id = pingone_application.%[2]s.id

  name  = "saml_subject"
  value = "$${user.email}"
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccApplicationCoreAttributeMappingConfig_SAML_OIDC_Name(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_key" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name                = "%[3]s"
  algorithm           = "RSA"
  key_length          = 4096
  signature_algorithm = "SHA512withRSA"
  subject_dn          = "CN=%[3]s, OU=BX Retail, O=BX Retail, L=, ST=, C=US"
  usage_type          = "SIGNING"
  validity_period     = 365
}

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true

  saml_options {
    acs_urls           = ["https://my-saas-app.com"]
    assertion_duration = 3600
    sp_entity_id       = "sp:entity:%[3]s"

    idp_signing_key {
      key_id    = pingone_key.%[2]s.id
      algorithm = pingone_key.%[2]s.signature_algorithm
    }
  }
}

resource "pingone_application_core_attribute_mapping" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  application_id = pingone_application.%[2]s.id

  name  = "sub"
  value = "$${user.email}"
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccApplicationCoreAttributeMappingConfig_BadApplication(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true

  external_link_options {
    home_page_url = "https://demo.bxretail.org/"
  }
}

resource "pingone_application_core_attribute_mapping" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  application_id = pingone_application.%[2]s.id

  name  = "saml_subject"
  value = "$${user.name.given + ', ' + user.name.family}"
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccApplicationCoreAttributeMappingConfig_Expression(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true

  oidc_options {
    type                        = "SINGLE_PAGE_APP"
    grant_types                 = ["AUTHORIZATION_CODE"]
    response_types              = ["CODE"]
    pkce_enforcement            = "S256_REQUIRED"
    token_endpoint_authn_method = "NONE"
    redirect_uris               = ["https://www.pingidentity.com"]
  }
}

resource "pingone_application_core_attribute_mapping" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  application_id = pingone_application.%[2]s.id

  name  = "sub"
  value = "$${user.name.given + ', ' + user.name.family}"
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
