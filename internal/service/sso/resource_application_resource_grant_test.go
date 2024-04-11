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

func TestAccApplicationResourceGrant_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application_resource_grant.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	var applicationResourceGrantID, applicationID, environmentID string

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
		CheckDestroy:             sso.ApplicationResourceGrant_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Test removal of the resource
			{
				Config: testAccApplicationResourceGrantConfig_CustomResource(resourceName, name),
				Check:  sso.ApplicationResourceGrant_GetIDs(resourceFullName, &environmentID, &applicationID, &applicationResourceGrantID),
			},
			{
				PreConfig: func() {
					sso.ApplicationResourceGrant_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID, applicationID, applicationResourceGrantID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the application
			{
				Config: testAccApplicationResourceGrantConfig_CustomResource(resourceName, name),
				Check:  sso.ApplicationResourceGrant_GetIDs(resourceFullName, &environmentID, &applicationID, &applicationResourceGrantID),
			},
			{
				PreConfig: func() {
					sso.Application_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID, applicationID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the environment
			{
				Config: testAccApplicationResourceGrantConfig_SelfService(environmentName, licenseID, resourceName),
				Check:  sso.ApplicationResourceGrant_GetIDs(resourceFullName, &environmentID, &applicationID, &applicationResourceGrantID),
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

func TestAccApplicationResourceGrant_OpenIDResource(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application_resource_grant.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.ApplicationResourceGrant_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationResourceGrantConfig_OpenIDResource(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "resource_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "resource_name", "openid"),
					resource.TestCheckResourceAttr(resourceFullName, "scopes.#", "2"),
					resource.TestMatchResourceAttr(resourceFullName, "scopes.0", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "scopes.1", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "scope_names.#", "2"),
					resource.TestCheckResourceAttr(resourceFullName, "scope_names.0", "email"),
					resource.TestCheckResourceAttr(resourceFullName, "scope_names.1", "profile"),
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

						return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["application_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Test error catch on update
			{
				Config:      testAccApplicationResourceGrantConfig_OpenIDResource_InvalidOpenIDScope(resourceName, name),
				ExpectError: regexp.MustCompile(`Invalid scope`),
			},
			{
				Config:  testAccApplicationResourceGrantConfig_OpenIDResource(resourceName, name),
				Destroy: true,
			},
			// Test error catch on from new
			{
				Config:      testAccApplicationResourceGrantConfig_OpenIDResource_InvalidOpenIDScope(resourceName, name),
				ExpectError: regexp.MustCompile(`Invalid scope`),
			},
		},
	})
}

func TestAccApplicationResourceGrant_CustomResource(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application_resource_grant.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.ApplicationResourceGrant_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationResourceGrantConfig_CustomResource(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "resource_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "resource_name", name),
					resource.TestCheckResourceAttr(resourceFullName, "scopes.#", "3"),
					resource.TestMatchResourceAttr(resourceFullName, "scopes.0", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "scopes.1", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "scopes.2", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "scope_names.#", "3"),
					resource.TestCheckResourceAttr(resourceFullName, "scope_names.0", fmt.Sprintf("%s-1", name)),
					resource.TestCheckResourceAttr(resourceFullName, "scope_names.1", fmt.Sprintf("%s-2", name)),
					resource.TestCheckResourceAttr(resourceFullName, "scope_names.2", fmt.Sprintf("%s-3", name)),
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

						return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["application_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccApplicationResourceGrant_SystemApplication(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application_resource_grant.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.ApplicationResourceGrant_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationResourceGrantConfig_SelfService(environmentName, licenseID, resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "resource_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "resource_name", "PingOne API"),
					resource.TestCheckResourceAttr(resourceFullName, "scopes.#", "8"),
					resource.TestMatchResourceAttr(resourceFullName, "scopes.0", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "scopes.1", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "scopes.2", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "scopes.3", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "scopes.4", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "scopes.5", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "scopes.6", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "scopes.7", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "scope_names.#", "8"),
					resource.TestCheckResourceAttr(resourceFullName, "scope_names.0", "p1:create:device"),
					resource.TestCheckResourceAttr(resourceFullName, "scope_names.1", "p1:create:pairingKey"),
					resource.TestCheckResourceAttr(resourceFullName, "scope_names.2", "p1:delete:device"),
					resource.TestCheckResourceAttr(resourceFullName, "scope_names.3", "p1:read:device"),
					resource.TestCheckResourceAttr(resourceFullName, "scope_names.4", "p1:read:pairingKey"),
					resource.TestCheckResourceAttr(resourceFullName, "scope_names.5", "p1:read:user"),
					resource.TestCheckResourceAttr(resourceFullName, "scope_names.6", "p1:update:device"),
					resource.TestCheckResourceAttr(resourceFullName, "scope_names.7", "p1:update:user"),
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

						return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["application_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccApplicationResourceGrantConfig_Portal(environmentName, licenseID, resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "resource_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "resource_name", "PingOne API"),
					resource.TestCheckResourceAttr(resourceFullName, "scopes.#", "8"),
					resource.TestMatchResourceAttr(resourceFullName, "scopes.0", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "scopes.1", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "scopes.2", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "scopes.3", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "scopes.4", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "scopes.5", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "scopes.6", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "scopes.7", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "scope_names.#", "8"),
					resource.TestCheckResourceAttr(resourceFullName, "scope_names.0", "p1:create:device"),
					resource.TestCheckResourceAttr(resourceFullName, "scope_names.1", "p1:create:pairingKey"),
					resource.TestCheckResourceAttr(resourceFullName, "scope_names.2", "p1:delete:device"),
					resource.TestCheckResourceAttr(resourceFullName, "scope_names.3", "p1:read:device"),
					resource.TestCheckResourceAttr(resourceFullName, "scope_names.4", "p1:read:pairingKey"),
					resource.TestCheckResourceAttr(resourceFullName, "scope_names.5", "p1:read:user"),
					resource.TestCheckResourceAttr(resourceFullName, "scope_names.6", "p1:update:device"),
					resource.TestCheckResourceAttr(resourceFullName, "scope_names.7", "p1:update:user"),
				),
			},
			// Test console error catch - TODO
		},
	})
}

func TestAccApplicationResourceGrant_Change(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application_resource_grant.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.ApplicationResourceGrant_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationResourceGrantConfig_OpenIDResource(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "resource_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "resource_name", "openid"),
					resource.TestCheckResourceAttr(resourceFullName, "scopes.#", "2"),
					resource.TestMatchResourceAttr(resourceFullName, "scopes.0", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "scopes.1", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "scope_names.#", "2"),
					resource.TestCheckResourceAttr(resourceFullName, "scope_names.0", "email"),
					resource.TestCheckResourceAttr(resourceFullName, "scope_names.1", "profile"),
				),
			},
			{
				Config: testAccApplicationResourceGrantConfig_CustomResource(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "resource_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "resource_name", name),
					resource.TestCheckResourceAttr(resourceFullName, "scopes.#", "3"),
					resource.TestMatchResourceAttr(resourceFullName, "scopes.0", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "scopes.1", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "scopes.2", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "scope_names.#", "3"),
					resource.TestCheckResourceAttr(resourceFullName, "scope_names.0", fmt.Sprintf("%s-1", name)),
					resource.TestCheckResourceAttr(resourceFullName, "scope_names.1", fmt.Sprintf("%s-2", name)),
					resource.TestCheckResourceAttr(resourceFullName, "scope_names.2", fmt.Sprintf("%s-3", name)),
				),
			},
			{
				Config: testAccApplicationResourceGrantConfig_OpenIDResource(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "resource_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "resource_name", "openid"),
					resource.TestCheckResourceAttr(resourceFullName, "scopes.#", "2"),
					resource.TestMatchResourceAttr(resourceFullName, "scopes.0", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "scopes.1", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "scope_names.#", "2"),
					resource.TestCheckResourceAttr(resourceFullName, "scope_names.0", "email"),
					resource.TestCheckResourceAttr(resourceFullName, "scope_names.1", "profile"),
				),
			},
		},
	})
}

func TestAccApplicationResourceGrant_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application_resource_grant.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.ApplicationResourceGrant_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccApplicationResourceGrantConfig_CustomResource(resourceName, name),
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

func testAccApplicationResourceGrantConfig_OpenIDResource(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true

  oidc_options = {
    type                        = "SINGLE_PAGE_APP"
    grant_types                 = ["AUTHORIZATION_CODE"]
    response_types              = ["CODE"]
    pkce_enforcement            = "S256_REQUIRED"
    token_endpoint_authn_method = "NONE"
    redirect_uris               = ["https://www.pingidentity.com"]
  }
}

resource "pingone_application_resource_grant" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  application_id = pingone_application.%[2]s.id

  resource_name = "openid"
  scope_names = [
    "email",
    "profile",
  ]
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccApplicationResourceGrantConfig_OpenIDResource_InvalidOpenIDScope(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true

  oidc_options = {
    type                        = "SINGLE_PAGE_APP"
    grant_types                 = ["AUTHORIZATION_CODE"]
    response_types              = ["CODE"]
    pkce_enforcement            = "S256_REQUIRED"
    token_endpoint_authn_method = "NONE"
    redirect_uris               = ["https://www.pingidentity.com"]
  }
}

resource "pingone_application_resource_grant" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  application_id = pingone_application.%[2]s.id

  resource_name = "openid"
  scope_names = [
    "email",
    "profile",
    "openid",
  ]
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccApplicationResourceGrantConfig_CustomResource(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true

  oidc_options = {
    type                        = "SINGLE_PAGE_APP"
    grant_types                 = ["AUTHORIZATION_CODE"]
    response_types              = ["CODE"]
    pkce_enforcement            = "S256_REQUIRED"
    token_endpoint_authn_method = "NONE"
    redirect_uris               = ["https://www.pingidentity.com"]
  }
}

resource "pingone_resource" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_resource_scope" "%[2]s-1" {
  environment_id = data.pingone_environment.general_test.id
  resource_id    = pingone_resource.%[2]s.id

  name = "%[3]s-1"
}

resource "pingone_resource_scope" "%[2]s-2" {
  environment_id = data.pingone_environment.general_test.id
  resource_id    = pingone_resource.%[2]s.id

  name = "%[3]s-2"
}

resource "pingone_resource_scope" "%[2]s-3" {
  environment_id = data.pingone_environment.general_test.id
  resource_id    = pingone_resource.%[2]s.id

  name = "%[3]s-3"
}

resource "pingone_application_resource_grant" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  application_id = pingone_application.%[2]s.id

  resource_name = pingone_resource.%[2]s.name
  scope_names = [
    pingone_resource_scope.%[2]s-1.name,
    pingone_resource_scope.%[2]s-2.name,
    pingone_resource_scope.%[2]s-3.name
  ]
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccApplicationResourceGrantConfig_SelfService(environmentName, licenseID, resourceName string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_system_application" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  type    = "PING_ONE_SELF_SERVICE"
  enabled = true

  apply_default_theme         = true
  enable_default_theme_footer = true
}

resource "pingone_application_resource_grant" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  application_id = pingone_system_application.%[3]s.id

  resource_name = "PingOne API"

  scope_names = [
    "p1:create:device",
    "p1:create:pairingKey",
    "p1:delete:device",
    "p1:read:device",
    "p1:read:pairingKey",
    "p1:read:user",
    "p1:update:device",
    "p1:update:user",
  ]
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName)
}

func testAccApplicationResourceGrantConfig_Portal(environmentName, licenseID, resourceName string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_system_application" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  type    = "PING_ONE_PORTAL"
  enabled = true
}

resource "pingone_application_resource_grant" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  application_id = pingone_system_application.%[3]s.id

  resource_name = "PingOne API"

  scope_names = [
    "p1:create:device",
    "p1:create:pairingKey",
    "p1:delete:device",
    "p1:read:device",
    "p1:read:pairingKey",
    "p1:read:user",
    "p1:update:device",
    "p1:update:user",
  ]
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName)
}
