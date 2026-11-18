// Copyright © 2026 Ping Identity Corporation

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
	acctestlegacysdk "github.com/pingidentity/terraform-provider-pingone/internal/acctest/legacysdk"
	baselegacysdk "github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/base/legacysdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/sso"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccResource_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_resource.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	var resourceID, environmentID string

	var p1Client *client.Client
	var ctx = context.Background()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoBeta(t)
			p1Client = acctestlegacysdk.PreCheckTestClient(ctx, t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Resource_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccResourceConfig_Minimal(resourceName, name),
				Check:  sso.Resource_GetIDs(resourceFullName, &environmentID, &resourceID),
			},
			{
				PreConfig: func() {
					sso.Resource_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID, resourceID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the environment
			{
				Config: testAccResourceConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check:  sso.Resource_GetIDs(resourceFullName, &environmentID, &resourceID),
			},
			{
				PreConfig: func() {
					baselegacysdk.Environment_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccResource_NewEnv(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_resource.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Resource_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
				),
			},
		},
	})
}

func TestAccResource_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_resource.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Resource_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceConfig_Full(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", "Test Resource"),
					resource.TestCheckResourceAttr(resourceFullName, "type", "CUSTOM"),
					resource.TestCheckResourceAttr(resourceFullName, "audience", fmt.Sprintf("%s-1", name)),
					resource.TestCheckResourceAttr(resourceFullName, "access_token_validity_seconds", "7200"),
					resource.TestCheckResourceAttr(resourceFullName, "introspect_endpoint_auth_method", "CLIENT_SECRET_POST"),
					resource.TestCheckResourceAttr(resourceFullName, "application_permissions_settings.claim_enabled", "true"),
				),
			},
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("resource not found: %s", resourceFullName)
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

func TestAccResource_Minimal(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_resource.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Resource_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceConfig_Minimal(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckNoResourceAttr(resourceFullName, "description"),
					resource.TestCheckResourceAttr(resourceFullName, "type", "CUSTOM"),
					resource.TestCheckResourceAttr(resourceFullName, "audience", name),
					resource.TestCheckResourceAttr(resourceFullName, "access_token_validity_seconds", "3600"),
					resource.TestCheckResourceAttr(resourceFullName, "introspect_endpoint_auth_method", "CLIENT_SECRET_BASIC"),
					resource.TestCheckResourceAttr(resourceFullName, "application_permissions_settings.claim_enabled", "false"),
				),
			},
		},
	})
}

func TestAccResource_Change(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_resource.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Resource_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceConfig_Minimal(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckNoResourceAttr(resourceFullName, "description"),
					resource.TestCheckResourceAttr(resourceFullName, "type", "CUSTOM"),
					resource.TestCheckResourceAttr(resourceFullName, "audience", name),
					resource.TestCheckResourceAttr(resourceFullName, "access_token_validity_seconds", "3600"),
					resource.TestCheckResourceAttr(resourceFullName, "introspect_endpoint_auth_method", "CLIENT_SECRET_BASIC"),
					resource.TestCheckResourceAttr(resourceFullName, "application_permissions_settings.claim_enabled", "false"),
				),
			},
			{
				Config: testAccResourceConfig_Full(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", "Test Resource"),
					resource.TestCheckResourceAttr(resourceFullName, "type", "CUSTOM"),
					resource.TestCheckResourceAttr(resourceFullName, "audience", fmt.Sprintf("%s-1", name)),
					resource.TestCheckResourceAttr(resourceFullName, "access_token_validity_seconds", "7200"),
					resource.TestCheckResourceAttr(resourceFullName, "introspect_endpoint_auth_method", "CLIENT_SECRET_POST"),
					resource.TestCheckResourceAttr(resourceFullName, "application_permissions_settings.claim_enabled", "true"),
				),
			},
			{
				Config: testAccResourceConfig_Minimal(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckNoResourceAttr(resourceFullName, "description"),
					resource.TestCheckResourceAttr(resourceFullName, "type", "CUSTOM"),
					resource.TestCheckResourceAttr(resourceFullName, "audience", name),
					resource.TestCheckResourceAttr(resourceFullName, "access_token_validity_seconds", "3600"),
					resource.TestCheckResourceAttr(resourceFullName, "introspect_endpoint_auth_method", "CLIENT_SECRET_BASIC"),
					resource.TestCheckResourceAttr(resourceFullName, "application_permissions_settings.claim_enabled", "false"),
				),
			},
		},
	})
}

func TestAccResource_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_resource.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Resource_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccResourceConfig_Minimal(resourceName, name),
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

func TestAccResource_Jwks(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_resource.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Resource_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure with jwks
			{
				Config: testAccResourceConfig_Jwks(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "introspect_endpoint_auth_method", "PRIVATE_KEY_JWT"),
					resource.TestCheckResourceAttr(resourceFullName, "jwks", testAccResourceJwksString),
				),
			},
			// Switch to jwks_url
			{
				Config: testAccResourceConfig_JwksUrl(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "introspect_endpoint_auth_method", "PRIVATE_KEY_JWT"),
					resource.TestCheckNoResourceAttr(resourceFullName, "jwks"),
					resource.TestCheckResourceAttr(resourceFullName, "jwks_url", "https://example.com/.well-known/jwks.json"),
				),
			},
			// Revert to default auth method without jwks attributes
			{
				Config: testAccResourceConfig_Minimal(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "introspect_endpoint_auth_method", "CLIENT_SECRET_BASIC"),
					resource.TestCheckNoResourceAttr(resourceFullName, "jwks"),
					resource.TestCheckNoResourceAttr(resourceFullName, "jwks_url"),
				),
			},
			// Errors: PRIVATE_KEY_JWT with neither jwks nor jwks_url
			{
				Config:      testAccResourceConfig_JwksNeither(resourceName, name),
				ExpectError: regexp.MustCompile(`jwksUrl or jwks is required when introspectEndpointAuthMethod=PRIVATE_KEY_JWT`),
			},
			// Errors: both jwks and jwks_url set
			{
				Config:      testAccResourceConfig_JwksBoth(resourceName, name),
				ExpectError: regexp.MustCompile(`Invalid Attribute Combination`),
			},
			// Revert to a valid minimal config so the post-test destroy does not
			// re-render the conflicting-attributes config from the error step above
			{
				Config: testAccResourceConfig_Minimal(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "introspect_endpoint_auth_method", "CLIENT_SECRET_BASIC"),
					resource.TestCheckNoResourceAttr(resourceFullName, "jwks"),
					resource.TestCheckNoResourceAttr(resourceFullName, "jwks_url"),
				),
			},
		},
	})
}

// testAccResourceJwksString is the expected state value of the `jwks`
// attribute when configured with testAccResourceJwksConfig (jsonencode of the
// fixture JWKS below).
const testAccResourceJwksString = `{"keys":[{"alg":"RS256","e":"AQAB","kid":"cdi1443-test","kty":"RSA","n":"0vx7agoebGcQSuuPiLJXZptN9nndrQmbXEps2aiAFbWhM78LhWx4cbbfAAtb0CmzzzyYV8dBnrKamTmULCcBLine1V2d70uH7jXGNKLBaaNjZjxBUxQXKFCpbmHBcjW1z5IbjgkG9ESj8STbLIimSbvW7wPE_Td0dzRRPUHMYu7YmC3RLWkv12M7O4MTFGM9vLhUdAD0msh0uDdrEBFQ0dO7YbXCHJ1FFG6gJXf9cmU3g8PLAETWdanT6wT36B0-HmF4FCYUktLXBGikmLQyqxCYBNfPTGqbUc0ETzAToZg2yx3TFW1hQxTVBwOq0ephT6who7IW38SyLv0WaBSG27Z9c9TH0mp2OqSjIiFfr7XcvKX8Q"}]}`

const testAccResourceJwksConfig = `jsonencode({
    keys = [
      {
        alg = "RS256"
        e   = "AQAB"
        kid = "cdi1443-test"
        kty = "RSA"
        n   = "0vx7agoebGcQSuuPiLJXZptN9nndrQmbXEps2aiAFbWhM78LhWx4cbbfAAtb0CmzzzyYV8dBnrKamTmULCcBLine1V2d70uH7jXGNKLBaaNjZjxBUxQXKFCpbmHBcjW1z5IbjgkG9ESj8STbLIimSbvW7wPE_Td0dzRRPUHMYu7YmC3RLWkv12M7O4MTFGM9vLhUdAD0msh0uDdrEBFQ0dO7YbXCHJ1FFG6gJXf9cmU3g8PLAETWdanT6wT36B0-HmF4FCYUktLXBGikmLQyqxCYBNfPTGqbUc0ETzAToZg2yx3TFW1hQxTVBwOq0ephT6who7IW38SyLv0WaBSG27Z9c9TH0mp2OqSjIiFfr7XcvKX8Q"
      },
    ]
  })`

func testAccResourceConfig_Jwks(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_resource" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  introspect_endpoint_auth_method = "PRIVATE_KEY_JWT"

  jwks = %s
}`, acctest.GenericSandboxEnvironment(), resourceName, name, testAccResourceJwksConfig)
}

func testAccResourceConfig_JwksUrl(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_resource" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  introspect_endpoint_auth_method = "PRIVATE_KEY_JWT"

  jwks_url = "https://example.com/.well-known/jwks.json"
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccResourceConfig_JwksNeither(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_resource" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  introspect_endpoint_auth_method = "PRIVATE_KEY_JWT"
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccResourceConfig_JwksBoth(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_resource" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  introspect_endpoint_auth_method = "PRIVATE_KEY_JWT"

  jwks = %s

  jwks_url = "https://example.com/.well-known/jwks.json"
}`, acctest.GenericSandboxEnvironment(), resourceName, name, testAccResourceJwksConfig)
}

func testAccResourceConfig_NewEnv(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_resource" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  name = "%[4]s"
}`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccResourceConfig_Full(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_resource" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name        = "%[3]s"
  description = "Test Resource"

  audience                      = "%[3]s-1"
  access_token_validity_seconds = 7200

  introspect_endpoint_auth_method = "CLIENT_SECRET_POST"

  application_permissions_settings = {
    claim_enabled = true
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccResourceConfig_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_resource" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
