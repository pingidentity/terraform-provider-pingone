package sso_test

import (
	"context"
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

func TestAccIdentityProviderAttribute_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_identity_provider_attribute.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	var identityProviderAttributeID, identityProviderID, environmentID string

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
		CheckDestroy:             sso.IdentityProviderAttribute_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Test removal of the resource
			{
				Config: testAccIdentityProviderAttributeConfig_Minimal(resourceName, name),
				Check:  sso.IdentityProviderAttribute_GetIDs(resourceFullName, &environmentID, &identityProviderID, &identityProviderAttributeID),
			},
			// Replan after removal preconfig
			{
				PreConfig: func() {
					sso.IdentityProviderAttribute_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID, identityProviderID, identityProviderAttributeID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the IDP
			{
				Config: testAccIdentityProviderAttributeConfig_Minimal(resourceName, name),
				Check:  sso.IdentityProviderAttribute_GetIDs(resourceFullName, &environmentID, &identityProviderID, &identityProviderAttributeID),
			},
			// Replan after removal preconfig
			{
				PreConfig: func() {
					sso.IdentityProvider_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID, identityProviderID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the environment
			{
				Config: testAccIdentityProviderAttributeConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check:  sso.IdentityProviderAttribute_GetIDs(resourceFullName, &environmentID, &identityProviderID, &identityProviderAttributeID),
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

func TestAccIdentityProviderAttribute_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_identity_provider_attribute.%s", resourceName)

	name := resourceName

	fullStep := resource.TestStep{
		Config: testAccIdentityProviderAttributeConfig_Full(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "identity_provider_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "name", "email"),
			resource.TestCheckResourceAttr(resourceFullName, "update", "ALWAYS"),
			resource.TestCheckResourceAttr(resourceFullName, "value", "${providerAttributes.emailAddress.value}"),
			resource.TestCheckResourceAttr(resourceFullName, "mapping_type", "CUSTOM"),
		),
	}

	minimalStep := resource.TestStep{
		Config: testAccIdentityProviderAttributeConfig_Minimal(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "identity_provider_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "name", "email"),
			resource.TestCheckResourceAttr(resourceFullName, "update", "EMPTY_ONLY"),
			resource.TestCheckResourceAttr(resourceFullName, "value", "${providerAttributes.name.givenName}"),
			resource.TestCheckResourceAttr(resourceFullName, "mapping_type", "CUSTOM"),
		),
	}

	expressionStep := resource.TestStep{
		Config: testAccIdentityProviderAttributeConfig_Expression(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "name", "name.given"),
			resource.TestCheckResourceAttr(resourceFullName, "update", "ALWAYS"),
			resource.TestCheckResourceAttr(resourceFullName, "value", "${providerAttributes.name.givenName + ', ' + providerAttributes.name.givenName}"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.IdentityProviderAttribute_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			fullStep,
			{
				Config:  testAccIdentityProviderAttributeConfig_Full(resourceName, name),
				Destroy: true,
			},
			// Minimal
			minimalStep,
			{
				Config:  testAccIdentityProviderAttributeConfig_Minimal(resourceName, name),
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

						return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["identity_provider_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:  testAccIdentityProviderAttributeConfig_Full(resourceName, name),
				Destroy: true,
			},
			// Expression
			expressionStep,
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["identity_provider_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccIdentityProviderAttribute_ReservedAttributeName(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.IdentityProviderAttribute_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccIdentityProviderAttributeConfig_ReservedAttributeName(resourceName, name),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Match`),
			},
		},
	})
}

func TestAccIdentityProviderAttribute_Core(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_identity_provider_attribute.%s", resourceName)

	name := resourceName

	fullStep := resource.TestStep{
		Config: testAccIdentityProviderAttributeConfig_Core_Full(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "identity_provider_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "name", "username"),
			resource.TestCheckResourceAttr(resourceFullName, "update", "EMPTY_ONLY"),
			resource.TestCheckResourceAttr(resourceFullName, "value", "${providerAttributes.emailAddress.value}"),
			resource.TestCheckResourceAttr(resourceFullName, "mapping_type", "CORE"),
		),
	}

	updateStep := resource.TestStep{
		Config: testAccIdentityProviderAttributeConfig_Core_Update(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "identity_provider_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "name", "username"),
			resource.TestCheckResourceAttr(resourceFullName, "update", "EMPTY_ONLY"),
			resource.TestCheckResourceAttr(resourceFullName, "value", "${providerAttributes.name.displayName}"),
			resource.TestCheckResourceAttr(resourceFullName, "mapping_type", "CORE"),
		),
	}

	expressionStep := resource.TestStep{
		Config: testAccIdentityProviderAttributeConfig_Core_Expression(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "name", "username"),
			resource.TestCheckResourceAttr(resourceFullName, "value", "${providerAttributes.name.displayName + ', ' + providerAttributes.name.displayName}"),
		),
	}

	badParameterStep := resource.TestStep{
		Config:      testAccIdentityProviderAttributeConfig_Core_BadParameter(resourceName, name),
		ExpectError: regexp.MustCompile(`Invalid parameter value - Parameter doesn't apply to attribute type`),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.IdentityProviderAttribute_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			fullStep,
			{
				Config:  testAccIdentityProviderAttributeConfig_Core_Full(resourceName, name),
				Destroy: true,
			},
			// Change
			fullStep,
			updateStep,
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

						return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["identity_provider_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:  testAccIdentityProviderAttributeConfig_Core_Full(resourceName, name),
				Destroy: true,
			},
			// Expression
			expressionStep,
			{
				Config:  testAccIdentityProviderAttributeConfig_Core_Expression(resourceName, name),
				Destroy: true,
			},
			// Bad parameters
			badParameterStep,
		},
	})
}

func TestAccIdentityProviderAttribute_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_identity_provider_attribute.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.IdentityProviderAttribute_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccIdentityProviderAttributeConfig_Minimal(resourceName, name),
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
				ImportStateId: "badformat/badformat/badformat",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`Unexpected Import Identifier`),
			},
		},
	})
}

func testAccIdentityProviderAttributeConfig_NewEnv(environmentName, licenseID, resourceName, name string) string {
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

resource "pingone_identity_provider_attribute" "%[3]s" {
  environment_id       = pingone_environment.%[2]s.id
  identity_provider_id = pingone_identity_provider.%[3]s.id

  name   = "email"
  update = "ALWAYS"
  value  = "$${providerAttributes.emailAddress.value}"
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccIdentityProviderAttributeConfig_Full(resourceName, name string) string {
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

resource "pingone_identity_provider_attribute" "%[2]s" {
  environment_id       = data.pingone_environment.general_test.id
  identity_provider_id = pingone_identity_provider.%[2]s.id

  name   = "email"
  update = "ALWAYS"
  value  = "$${providerAttributes.emailAddress.value}"
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccIdentityProviderAttributeConfig_Minimal(resourceName, name string) string {
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

resource "pingone_identity_provider_attribute" "%[2]s" {
  environment_id       = data.pingone_environment.general_test.id
  identity_provider_id = pingone_identity_provider.%[2]s.id

  name  = "email"
  value = "$${providerAttributes.name.givenName}"
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccIdentityProviderAttributeConfig_Expression(resourceName, name string) string {
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

resource "pingone_identity_provider_attribute" "%[2]s" {
  environment_id       = data.pingone_environment.general_test.id
  identity_provider_id = pingone_identity_provider.%[2]s.id

  name   = "name.given"
  update = "ALWAYS"
  value  = "$${providerAttributes.name.givenName + ', ' + providerAttributes.name.givenName}"
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccIdentityProviderAttributeConfig_ReservedAttributeName(resourceName, name string) string {
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

resource "pingone_identity_provider_attribute" "%[2]s" {
  environment_id       = data.pingone_environment.general_test.id
  identity_provider_id = pingone_identity_provider.%[2]s.id

  name   = "account"
  update = "ALWAYS"
  value  = "$${'test'}"
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccIdentityProviderAttributeConfig_Core_Full(resourceName, name string) string {
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

resource "pingone_identity_provider_attribute" "%[2]s" {
  environment_id       = data.pingone_environment.general_test.id
  identity_provider_id = pingone_identity_provider.%[2]s.id

  name  = "username"
  value = "$${providerAttributes.emailAddress.value}"
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccIdentityProviderAttributeConfig_Core_Update(resourceName, name string) string {
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

resource "pingone_identity_provider_attribute" "%[2]s" {
  environment_id       = data.pingone_environment.general_test.id
  identity_provider_id = pingone_identity_provider.%[2]s.id

  name  = "username"
  value = "$${providerAttributes.name.displayName}"
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccIdentityProviderAttributeConfig_Core_Expression(resourceName, name string) string {
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

resource "pingone_identity_provider_attribute" "%[2]s" {
  environment_id       = data.pingone_environment.general_test.id
  identity_provider_id = pingone_identity_provider.%[2]s.id

  name  = "username"
  value = "$${providerAttributes.name.displayName + ', ' + providerAttributes.name.displayName}"
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccIdentityProviderAttributeConfig_Core_BadParameter(resourceName, name string) string {
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

resource "pingone_identity_provider_attribute" "%[2]s" {
  environment_id       = data.pingone_environment.general_test.id
  identity_provider_id = pingone_identity_provider.%[2]s.id

  name   = "username"
  update = "ALWAYS"
  value  = "$${providerAttributes.name.displayName + ', ' + providerAttributes.name.displayName}"
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
