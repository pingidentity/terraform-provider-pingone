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
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func testAccCheckResourceAttributeDestroy(s *terraform.State) error {
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
		if rs.Type != "pingone_resource_attribute" {
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

		body, r, err := apiClient.ResourceAttributesApi.ReadOneResourceAttribute(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["resource_id"], rs.Primary.ID).Execute()

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

		return fmt.Errorf("PingOne Resource Mapping %s still exists", rs.Primary.ID)
	}

	return nil
}

func TestAccResourceAttribute_OIDC_Custom(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_resource_attribute.%s", resourceName)

	name := resourceName

	fullStep := resource.TestStep{
		Config: testAccResourceAttributeConfig_OIDC_Custom_Full(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
			resource.TestMatchResourceAttr(resourceFullName, "resource_id", verify.P1ResourceIDRegexp),
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestCheckResourceAttr(resourceFullName, "value", "${user.email}"),
			resource.TestCheckResourceAttr(resourceFullName, "id_token_enabled", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "userinfo_enabled", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "type", "CUSTOM"),
		),
	}

	minimalStep := resource.TestStep{
		Config: testAccResourceAttributeConfig_OIDC_Custom_Minimal(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
			resource.TestMatchResourceAttr(resourceFullName, "resource_id", verify.P1ResourceIDRegexp),
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestCheckResourceAttr(resourceFullName, "value", "${user.name.given}"),
			resource.TestCheckResourceAttr(resourceFullName, "id_token_enabled", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "userinfo_enabled", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "type", "CUSTOM"),
		),
	}

	expressionStep := resource.TestStep{
		Config: testAccResourceAttributeConfig_OIDC_Expression(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestCheckResourceAttr(resourceFullName, "value", "${user.name.given + ', ' + user.name.family}"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckResourceAttributeDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			fullStep,
			{
				Config:  testAccResourceAttributeConfig_OIDC_Custom_Full(resourceName, name),
				Destroy: true,
			},
			// Minimal
			minimalStep,
			{
				Config:  testAccResourceAttributeConfig_OIDC_Custom_Minimal(resourceName, name),
				Destroy: true,
			},
			// Change
			fullStep,
			minimalStep,
			fullStep,
			{
				Config:  testAccResourceAttributeConfig_OIDC_Custom_Full(resourceName, name),
				Destroy: true,
			},
			// Expression
			expressionStep,
		},
	})
}

func TestAccResourceAttribute_OIDC_ReservedAttributeName(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckResourceAttributeDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccResourceAttributeConfig_OIDC_ReservedAttributeName(resourceName, name),
				ExpectError: regexp.MustCompile("Invalid attribute name `[a-zA-Z]*` for the configured OpenID Connect resource."),
			},
		},
	})
}

func TestAccResourceAttribute_OIDC_Predefined(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_resource_attribute.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	fullStep := resource.TestStep{
		Config: testAccResourceAttributeConfig_OIDC_Predefined_Full(environmentName, licenseID, resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
			resource.TestMatchResourceAttr(resourceFullName, "resource_id", verify.P1ResourceIDRegexp),
			resource.TestCheckResourceAttr(resourceFullName, "name", "email"),
			resource.TestCheckResourceAttr(resourceFullName, "value", fmt.Sprintf("${user.%s}", name)),
			resource.TestCheckResourceAttr(resourceFullName, "id_token_enabled", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "userinfo_enabled", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "type", "PREDEFINED"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckResourceAttributeDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			fullStep,
		},
	})
}

func TestAccResourceAttribute_Custom(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_resource_attribute.%s", resourceName)

	name := resourceName

	fullStep := resource.TestStep{
		Config: testAccResourceAttributeConfig_Custom_Full(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
			resource.TestMatchResourceAttr(resourceFullName, "resource_id", verify.P1ResourceIDRegexp),
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestCheckResourceAttr(resourceFullName, "value", "${user.email}"),
			resource.TestCheckResourceAttr(resourceFullName, "type", "CUSTOM"),
		),
	}

	expressionStep := resource.TestStep{
		Config: testAccResourceAttributeConfig_Custom_Expression(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestCheckResourceAttr(resourceFullName, "value", "${user.name.given + ', ' + user.name.family}"),
		),
	}

	invalidParameterStep := resource.TestStep{
		Config:      testAccResourceAttributeConfig_Custom_BadParameters(resourceName, name),
		ExpectError: regexp.MustCompile(`Invalid parameter value - Parameter doesn't apply to resource type`),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckResourceAttributeDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			fullStep,
			{
				Config:  testAccResourceAttributeConfig_Custom_Full(resourceName, name),
				Destroy: true,
			},
			// Expression
			expressionStep,
			// Change
			fullStep,
			expressionStep,
			fullStep,
			{
				Config:  testAccResourceAttributeConfig_Custom_Full(resourceName, name),
				Destroy: true,
			},
			// Bad parameters
			invalidParameterStep,
		},
	})
}

func TestAccResourceAttribute_Custom_CoreAttribute(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_resource_attribute.%s", resourceName)

	name := resourceName

	fullStep := resource.TestStep{
		Config: testAccResourceAttributeConfig_Custom_Core_Full(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
			resource.TestMatchResourceAttr(resourceFullName, "resource_id", verify.P1ResourceIDRegexp),
			resource.TestCheckResourceAttr(resourceFullName, "name", "sub"),
			resource.TestCheckResourceAttr(resourceFullName, "value", "${user.email}"),
			resource.TestCheckResourceAttr(resourceFullName, "type", "CORE"),
		),
	}

	expressionStep := resource.TestStep{
		Config: testAccResourceAttributeConfig_Custom_Core_Expression(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "name", "sub"),
			resource.TestCheckResourceAttr(resourceFullName, "value", "${user.name.given + ', ' + user.name.family}"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckResourceAttributeDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			fullStep,
			{
				Config:  testAccResourceAttributeConfig_Custom_Core_Full(resourceName, name),
				Destroy: true,
			},
			// Expression
			expressionStep,
			// Change
			fullStep,
			expressionStep,
			fullStep,
			{
				Config:  testAccResourceAttributeConfig_Custom_Core_Full(resourceName, name),
				Destroy: true,
			},
		},
	})
}

func TestAccResourceAttribute_BadResource(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckResourceAttributeDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccResourceAttributeConfig_BadResource(resourceName, name),
				ExpectError: regexp.MustCompile("Invalid parameter value - Invalid resource type"),
			},
		},
	})
}

func testAccResourceAttributeConfig_OIDC_Custom_Full(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

data "pingone_resource" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "openid"
}

resource "pingone_resource_attribute" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  resource_id    = data.pingone_resource.%[2]s.id

  name  = "%[3]s"
  value = "$${user.email}"

  id_token_enabled = true
  userinfo_enabled = false
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccResourceAttributeConfig_OIDC_Custom_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

data "pingone_resource" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "openid"
}

resource "pingone_resource_attribute" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  resource_id    = data.pingone_resource.%[2]s.id

  name  = "%[3]s"
  value = "$${user.name.given}"
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccResourceAttributeConfig_OIDC_Expression(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

data "pingone_resource" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "openid"
}

resource "pingone_resource_attribute" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  resource_id    = data.pingone_resource.%[2]s.id

  name  = "%[3]s"
  value = "$${user.name.given + ', ' + user.name.family}"
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccResourceAttributeConfig_OIDC_ReservedAttributeName(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

data "pingone_resource" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "openid"
}

resource "pingone_resource_attribute" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  resource_id    = data.pingone_resource.%[2]s.id

  name  = "aud"
  value = "$${'test'}"
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccResourceAttributeConfig_OIDC_Predefined_Full(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

data "pingone_schema" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  name = "User"
}

resource "pingone_schema_attribute" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  schema_id      = data.pingone_schema.%[3]s.id

  name         = "%[4]s"
  display_name = "My Attribute"
  description  = "My new attribute"

  type        = "STRING"
  unique      = false
  multivalued = false
}

data "pingone_resource" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  name = "openid"
}

resource "pingone_resource_attribute" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  resource_id    = data.pingone_resource.%[3]s.id

  name  = "email"
  value = "$${user.%[4]s}"
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccResourceAttributeConfig_Custom_Full(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_resource" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_resource_attribute" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  resource_id    = pingone_resource.%[2]s.id

  name  = "%[3]s"
  value = "$${user.email}"
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccResourceAttributeConfig_Custom_Expression(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_resource" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_resource_attribute" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  resource_id    = pingone_resource.%[2]s.id

  name  = "%[3]s"
  value = "$${user.name.given + ', ' + user.name.family}"
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccResourceAttributeConfig_Custom_Core_Full(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_resource" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_resource_attribute" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  resource_id    = pingone_resource.%[2]s.id

  name  = "sub"
  value = "$${user.email}"
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccResourceAttributeConfig_Custom_Core_Expression(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_resource" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_resource_attribute" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  resource_id    = pingone_resource.%[2]s.id

  name  = "sub"
  value = "$${user.name.given + ', ' + user.name.family}"
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccResourceAttributeConfig_Custom_BadParameters(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_resource" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_resource_attribute" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  resource_id    = pingone_resource.%[2]s.id

  name  = "%[3]s"
  value = "$${user.email}"

  id_token_enabled = true
  userinfo_enabled = false
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccResourceAttributeConfig_BadResource(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

data "pingone_resource" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "PingOne API"
}

resource "pingone_resource_attribute" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  resource_id    = data.pingone_resource.%[2]s.id

  name  = "%[3]s"
  value = "$${user.email}"
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
