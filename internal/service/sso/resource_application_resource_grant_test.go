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
					resource.TestMatchResourceAttr(resourceFullName, "application_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "resource_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckNoResourceAttr(resourceFullName, "custom_resource_id"),
					resource.TestCheckResourceAttr(resourceFullName, "resource_type", "OPENID_CONNECT"),
					resource.TestCheckResourceAttr(resourceFullName, "scopes.#", "4"),
					resource.TestMatchResourceAttr(resourceFullName, "scopes.0", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "scopes.1", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "scopes.2", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "scopes.3", verify.P1ResourceIDRegexpFullString),
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
			// // Test error catch on update
			// {
			// 	Config: testAccApplicationResourceGrantConfig_OpenIDResource_InvalidOpenIDScope(resourceName, name),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			// 	),
			// 	//ExpectError: regexp.MustCompile(`Invalid scope`),
			// },
			// {
			// 	Config:  testAccApplicationResourceGrantConfig_OpenIDResource(resourceName, name),
			// 	Destroy: true,
			// },
			// // Test error catch on from new
			// {
			// 	Config: testAccApplicationResourceGrantConfig_OpenIDResource_InvalidOpenIDScope(resourceName, name),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			// 	),
			// 	//ExpectError: regexp.MustCompile(`Invalid scope`),
			// },
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
					resource.TestMatchResourceAttr(resourceFullName, "application_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "resource_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "custom_resource_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "resource_type", "CUSTOM"),
					resource.TestCheckResourceAttr(resourceFullName, "scopes.#", "5"),
					resource.TestMatchResourceAttr(resourceFullName, "scopes.0", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "scopes.1", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "scopes.2", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "scopes.3", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "scopes.4", verify.P1ResourceIDRegexpFullString),
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
				Config: testAccApplicationResourceGrantConfig_CustomResource(resourceName, name),
				Taint: []string{
					fmt.Sprintf("pingone_resource_scope.%[1]s-1", name),
					fmt.Sprintf("pingone_resource_scope.%[1]s-2", name),
					fmt.Sprintf("pingone_resource_scope.%[1]s-3", name),
					fmt.Sprintf("pingone_resource_scope.%[1]s-4", name),
					fmt.Sprintf("pingone_resource_scope.%[1]s-5", name),
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
				),
			},
		},
	})
}

func TestAccApplicationResourceGrant_CustomResource_SimultaneousGrantRemoval(t *testing.T) {
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
				Config: testAccApplicationResourceGrantConfig_CustomResource_SimultaneousGrantRemoval(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
				),
			},
			{
				Config: testAccApplicationResourceGrantConfig_CustomResource_SimultaneousGrantRemoval(resourceName, name),
				Taint: []string{
					fmt.Sprintf("pingone_resource_scope.%[1]s-1", name),
					fmt.Sprintf("pingone_resource_scope.%[1]s-2", name),
					fmt.Sprintf("pingone_resource_scope.%[1]s-3", name),
					fmt.Sprintf("pingone_resource_scope.%[1]s-4", name),
					fmt.Sprintf("pingone_resource_scope.%[1]s-5", name),
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
				),
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
					resource.TestMatchResourceAttr(resourceFullName, "application_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "resource_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckNoResourceAttr(resourceFullName, "custom_resource_id"),
					resource.TestCheckResourceAttr(resourceFullName, "resource_type", "PINGONE_API"),
					resource.TestCheckResourceAttr(resourceFullName, "scopes.#", "4"),
					resource.TestMatchResourceAttr(resourceFullName, "scopes.0", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "scopes.1", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "scopes.2", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "scopes.3", verify.P1ResourceIDRegexpFullString),
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
					resource.TestMatchResourceAttr(resourceFullName, "application_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "resource_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckNoResourceAttr(resourceFullName, "custom_resource_id"),
					resource.TestCheckResourceAttr(resourceFullName, "resource_type", "PINGONE_API"),
					resource.TestCheckResourceAttr(resourceFullName, "scopes.#", "3"),
					resource.TestMatchResourceAttr(resourceFullName, "scopes.0", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "scopes.1", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "scopes.2", verify.P1ResourceIDRegexpFullString),
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
					resource.TestMatchResourceAttr(resourceFullName, "application_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "resource_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckNoResourceAttr(resourceFullName, "custom_resource_id"),
					resource.TestCheckResourceAttr(resourceFullName, "resource_type", "OPENID_CONNECT"),
					resource.TestCheckResourceAttr(resourceFullName, "scopes.#", "4"),
					resource.TestMatchResourceAttr(resourceFullName, "scopes.0", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "scopes.1", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "scopes.2", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "scopes.3", verify.P1ResourceIDRegexpFullString),
				),
			},
			{
				Config: testAccApplicationResourceGrantConfig_CustomResource(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "application_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "resource_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "custom_resource_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "resource_type", "CUSTOM"),
					resource.TestCheckResourceAttr(resourceFullName, "scopes.#", "5"),
					resource.TestMatchResourceAttr(resourceFullName, "scopes.0", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "scopes.1", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "scopes.2", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "scopes.3", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "scopes.4", verify.P1ResourceIDRegexpFullString),
				),
			},
			{
				Config: testAccApplicationResourceGrantConfig_OpenIDResource(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "application_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "resource_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckNoResourceAttr(resourceFullName, "custom_resource_id"),
					resource.TestCheckResourceAttr(resourceFullName, "resource_type", "OPENID_CONNECT"),
					resource.TestCheckResourceAttr(resourceFullName, "scopes.#", "4"),
					resource.TestMatchResourceAttr(resourceFullName, "scopes.0", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "scopes.1", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "scopes.2", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "scopes.3", verify.P1ResourceIDRegexpFullString),
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

data "pingone_resource_scope" "%[2]s_openid" {
  environment_id = data.pingone_environment.general_test.id

  resource_type = "OPENID_CONNECT"
  name          = "openid"
}

data "pingone_resource_scope" "%[2]s_email" {
  environment_id = data.pingone_environment.general_test.id

  resource_type = "OPENID_CONNECT"
  name          = "email"
}

data "pingone_resource_scope" "%[2]s_profile" {
  environment_id = data.pingone_environment.general_test.id

  resource_type = "OPENID_CONNECT"
  name          = "profile"
}

resource "pingone_resource_attribute" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  resource_type = "OPENID_CONNECT"
  name          = "%[3]s"
  value         = "$${user.name.given}"
}

resource "pingone_resource_scope_openid" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[2]s"

  mapped_claims = [
    pingone_resource_attribute.%[2]s.id
  ]
}

resource "pingone_application_resource_grant" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  application_id = pingone_application.%[2]s.id

  resource_type = "OPENID_CONNECT"
  scopes = [
    data.pingone_resource_scope.%[2]s_openid.id,
    data.pingone_resource_scope.%[2]s_email.id,
    data.pingone_resource_scope.%[2]s_profile.id,
    pingone_resource_scope_openid.%[2]s.id,
  ]
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

// func testAccApplicationResourceGrantConfig_OpenIDResource_InvalidOpenIDScope(resourceName, name string) string {
// 	return fmt.Sprintf(`
// 		%[1]s

// resource "pingone_application" "%[2]s" {
//   environment_id = data.pingone_environment.general_test.id
//   name           = "%[3]s"
//   enabled        = true

//   oidc_options = {
//     type                        = "SINGLE_PAGE_APP"
//     grant_types                 = ["AUTHORIZATION_CODE"]
//     response_types              = ["CODE"]
//     pkce_enforcement            = "S256_REQUIRED"
//     token_endpoint_authn_method = "NONE"
//     redirect_uris               = ["https://www.pingidentity.com"]
//   }
// }

// data "pingone_resource_scope" "%[2]s_email" {
//   environment_id = data.pingone_environment.general_test.id

//   resource_type = "OPENID_CONNECT"
//   name          = "email"
// }

// data "pingone_resource_scope" "%[2]s_profile" {
//   environment_id = data.pingone_environment.general_test.id

//   resource_type = "OPENID_CONNECT"
//   name          = "profile"
// }

// resource "pingone_resource_attribute" "%[2]s" {
//   environment_id = data.pingone_environment.general_test.id

//   resource_type = "OPENID_CONNECT"
//   name          = "%[3]s"
//   value         = "$${user.name.given}"
// }

// resource "pingone_resource_scope_openid" "%[2]s" {
//   environment_id = data.pingone_environment.general_test.id

//   name = "%[2]s"

//   mapped_claims = [
//     pingone_resource_attribute.%[2]s.id
//   ]
// }

// resource "pingone_application_resource_grant" "%[2]s" {
//   environment_id = data.pingone_environment.general_test.id
//   application_id = pingone_application.%[2]s.id

//   resource_type = "OPENID_CONNECT"
//   scopes = [
//     data.pingone_resource_scope.%[2]s_email.id,
//     data.pingone_resource_scope.%[2]s_profile.id,
//     pingone_resource_scope_openid.%[2]s.id,
//   ]
// }`, acctest.GenericSandboxEnvironment(), resourceName, name)
// }

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

resource "pingone_resource_scope" "%[2]s-4" {
  environment_id = data.pingone_environment.general_test.id
  resource_id    = pingone_resource.%[2]s.id

  name = "%[3]s-4"
}

resource "pingone_resource_scope" "%[2]s-5" {
  environment_id = data.pingone_environment.general_test.id
  resource_id    = pingone_resource.%[2]s.id

  name = "%[3]s-5"
}

resource "pingone_application_resource_grant" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  application_id = pingone_application.%[2]s.id

  resource_type      = "CUSTOM"
  custom_resource_id = pingone_resource.%[2]s.id
  scopes = [
    pingone_resource_scope.%[2]s-2.id,
    pingone_resource_scope.%[2]s-1.id,
    pingone_resource_scope.%[2]s-3.id,
    pingone_resource_scope.%[2]s-4.id,
    pingone_resource_scope.%[2]s-5.id,
  ]
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccApplicationResourceGrantConfig_CustomResource_SimultaneousGrantRemoval(resourceName, name string) string {
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

resource "pingone_resource_scope" "%[2]s-4" {
  environment_id = data.pingone_environment.general_test.id
  resource_id    = pingone_resource.%[2]s.id

  name = "%[3]s-4"
}

resource "pingone_resource_scope" "%[2]s-5" {
  environment_id = data.pingone_environment.general_test.id
  resource_id    = pingone_resource.%[2]s.id

  name = "%[3]s-5"
}

resource "pingone_application_resource_grant" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  application_id = pingone_application.%[2]s.id

  resource_type = "CUSTOM"
  custom_resource_id  = pingone_resource.%[2]s.id

  scopes = [
    pingone_resource_scope.%[2]s-1.id,
    pingone_resource_scope.%[2]s-2.id,
    pingone_resource_scope.%[2]s-3.id,
    pingone_resource_scope.%[2]s-4.id,
    pingone_resource_scope.%[2]s-5.id,
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

data "pingone_resource_scope" "%[2]s_read_user" {
  environment_id = pingone_environment.%[2]s.id

  resource_type = "PINGONE_API"
  name          = "p1:read:user"
}

data "pingone_resource_scope" "%[2]s_update_user" {
  environment_id = pingone_environment.%[2]s.id

  resource_type = "PINGONE_API"
  name          = "p1:update:user"
}

data "pingone_resource_scope" "%[2]s_create_device" {
  environment_id = pingone_environment.%[2]s.id

  resource_type = "PINGONE_API"
  name          = "p1:create:device"
}

data "pingone_resource_scope" "%[2]s_create_pairing_key" {
  environment_id = pingone_environment.%[2]s.id

  resource_type = "PINGONE_API"
  name          = "p1:create:pairingKey"
}

resource "pingone_application_resource_grant" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  application_id = pingone_system_application.%[3]s.id

  resource_type = "PINGONE_API"

  scopes = [
    data.pingone_resource_scope.%[2]s_read_user.id,
    data.pingone_resource_scope.%[2]s_update_user.id,
    data.pingone_resource_scope.%[2]s_create_device.id,
    data.pingone_resource_scope.%[2]s_create_pairing_key.id,
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

data "pingone_resource_scope" "pingone_api_read_user" {
  environment_id = pingone_environment.%[2]s.id

  resource_type = "PINGONE_API"
  name          = "p1:read:user"
}

data "pingone_resource_scope" "pingone_api_update_user" {
  environment_id = pingone_environment.%[2]s.id

  resource_type = "PINGONE_API"
  name          = "p1:update:user"
}

data "pingone_resource_scope" "pingone_api_create_device" {
  environment_id = pingone_environment.%[2]s.id

  resource_type = "PINGONE_API"
  name          = "p1:create:device"
}

data "pingone_resource_scope" "pingone_api_create_pairing_key" {
  environment_id = pingone_environment.%[2]s.id

  resource_type = "PINGONE_API"
  name          = "p1:create:pairingKey"
}

resource "pingone_application_resource_grant" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  application_id = pingone_system_application.%[3]s.id

  resource_type = "PINGONE_API"

  scopes = [
    data.pingone_resource_scope.pingone_api_read_user.id,
    data.pingone_resource_scope.pingone_api_update_user.id,
    data.pingone_resource_scope.pingone_api_create_device.id,
  ]
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName)
}
