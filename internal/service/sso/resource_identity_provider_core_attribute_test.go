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

func testAccCheckIdentityProviderCoreAttributeDestroy(s *terraform.State) error {
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
		if rs.Type != "pingone_identity_provider_core_attribute" {
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

		body, r, err := apiClient.IdentityProviderAttributesApi.ReadOneIdentityProviderAttribute(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["identity_provider_id"], rs.Primary.ID).Execute()

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

		return fmt.Errorf("PingOne Identity Provider attribute %s still exists", rs.Primary.ID)
	}

	return nil
}

func TestAccIdentityProviderCoreAttribute_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_identity_provider_core_attribute.%s", resourceName)

	name := resourceName

	fullStep := resource.TestStep{
		Config: testAccIdentityProviderCoreAttributeConfig_Full(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
			resource.TestMatchResourceAttr(resourceFullName, "identity_provider_id", verify.P1ResourceIDRegexp),
			resource.TestCheckResourceAttr(resourceFullName, "name", "username"),
			resource.TestCheckResourceAttr(resourceFullName, "update", "EMPTY_ONLY"),
			resource.TestCheckResourceAttr(resourceFullName, "value", "${providerAttributes.emailAddress.value}"),
			resource.TestCheckResourceAttr(resourceFullName, "mapping_type", "CORE"),
		),
	}

	updateStep := resource.TestStep{
		Config: testAccIdentityProviderCoreAttributeConfig_Update(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
			resource.TestMatchResourceAttr(resourceFullName, "identity_provider_id", verify.P1ResourceIDRegexp),
			resource.TestCheckResourceAttr(resourceFullName, "name", "username"),
			resource.TestCheckResourceAttr(resourceFullName, "update", "EMPTY_ONLY"),
			resource.TestCheckResourceAttr(resourceFullName, "value", "${providerAttributes.name.displayName}"),
			resource.TestCheckResourceAttr(resourceFullName, "mapping_type", "CORE"),
		),
	}

	expressionStep := resource.TestStep{
		Config: testAccIdentityProviderCoreAttributeConfig_Expression(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "name", "username"),
			resource.TestCheckResourceAttr(resourceFullName, "value", "${providerAttributes.name.displayName + ', ' + providerAttributes.name.displayName}"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckIdentityProviderCoreAttributeDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			fullStep,
			{
				Config:  testAccIdentityProviderCoreAttributeConfig_Full(resourceName, name),
				Destroy: true,
			},
			// Change
			fullStep,
			updateStep,
			fullStep,
			{
				Config:  testAccIdentityProviderCoreAttributeConfig_Full(resourceName, name),
				Destroy: true,
			},
			// Expression
			expressionStep,
		},
	})
}

func TestAccIdentityProviderCoreAttribute_NonCoreAttribute(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckIdentityProviderCoreAttributeDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccIdentityProviderCoreAttributeConfig_NonCoreAttribute(resourceName, name),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Match`),
			},
		},
	})
}

func testAccIdentityProviderCoreAttributeConfig_Full(resourceName, name string) string {
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

resource "pingone_identity_provider_core_attribute" "%[2]s" {
  environment_id       = data.pingone_environment.general_test.id
  identity_provider_id = pingone_identity_provider.%[2]s.id

  name  = "username"
  value = "$${providerAttributes.emailAddress.value}"
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccIdentityProviderCoreAttributeConfig_Update(resourceName, name string) string {
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

resource "pingone_identity_provider_core_attribute" "%[2]s" {
  environment_id       = data.pingone_environment.general_test.id
  identity_provider_id = pingone_identity_provider.%[2]s.id

  name  = "username"
  value = "$${providerAttributes.name.displayName}"
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccIdentityProviderCoreAttributeConfig_Expression(resourceName, name string) string {
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

resource "pingone_identity_provider_core_attribute" "%[2]s" {
  environment_id       = data.pingone_environment.general_test.id
  identity_provider_id = pingone_identity_provider.%[2]s.id

  name  = "username"
  value = "$${providerAttributes.name.displayName + ', ' + providerAttributes.name.displayName}"
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccIdentityProviderCoreAttributeConfig_NonCoreAttribute(resourceName, name string) string {
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

resource "pingone_identity_provider_core_attribute" "%[2]s" {
  environment_id       = data.pingone_environment.general_test.id
  identity_provider_id = pingone_identity_provider.%[2]s.id

  name  = "account"
  value = "$${'test'}"
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
