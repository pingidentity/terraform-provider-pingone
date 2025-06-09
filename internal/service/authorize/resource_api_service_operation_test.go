// Copyright © 2025 Ping Identity Corporation

package authorize_test

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
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/authorize"
	baselegacysdk "github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/base/legacysdk"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccAPIServiceOperation_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_api_service_operation.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	var apiServiceOperationID, apiServiceID, environmentID string

	var p1Client *client.Client
	var ctx = context.Background()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)

			p1Client = acctestlegacysdk.PreCheckTestClient(ctx, t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             authorize.APIServiceOperation_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccAPIServiceOperationConfig_Minimal(resourceName, name),
				Check:  authorize.APIServiceOperation_GetIDs(resourceFullName, &environmentID, &apiServiceID, &apiServiceOperationID),
			},
			{
				PreConfig: func() {
					authorize.APIServiceOperation_RemovalDrift_PreConfig(ctx, p1Client.API.AuthorizeAPIClient, t, environmentID, apiServiceID, apiServiceOperationID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the API service
			{
				Config: testAccAPIServiceOperationConfig_Minimal(resourceName, name),
				Check:  authorize.APIServiceOperation_GetIDs(resourceFullName, &environmentID, &apiServiceID, &apiServiceOperationID),
			},
			{
				PreConfig: func() {
					authorize.APIService_RemovalDrift_PreConfig(ctx, p1Client.API.AuthorizeAPIClient, t, environmentID, apiServiceID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the environment
			{
				Config: testAccAPIServiceOperationConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check:  authorize.APIServiceOperation_GetIDs(resourceFullName, &environmentID, &apiServiceID, &apiServiceOperationID),
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

func TestAccAPIServiceOperation_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_api_service_operation.%s", resourceName)

	name := resourceName

	fullStep := resource.TestStep{
		Config: testAccAPIServiceOperationConfig_Full(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "api_service_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "access_control.group.groups.#", "2"),
			resource.TestMatchResourceAttr(resourceFullName, "access_control.group.groups.0.id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "access_control.group.groups.1.id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckNoResourceAttr(resourceFullName, "access_control.permission"),
			resource.TestCheckResourceAttr(resourceFullName, "access_control.scope.match_type", "ALL"),
			resource.TestCheckResourceAttr(resourceFullName, "access_control.scope.scopes.#", "2"),
			resource.TestMatchResourceAttr(resourceFullName, "access_control.scope.scopes.0.id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "access_control.scope.scopes.1.id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "methods.#", "3"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "methods.*", "GET"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "methods.*", "POST"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "methods.*", "PUT"),
			resource.TestCheckResourceAttr(resourceFullName, "paths.#", "3"),
			resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "paths.*", map[string]string{
				"pattern": "/test/1",
				"type":    "EXACT",
			}),
			resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "paths.*", map[string]string{
				"pattern": "/test/{variable}/*",
				"type":    "PARAMETER",
			}),
			resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "paths.*", map[string]string{
				"pattern": "/test/2",
				"type":    "EXACT",
			}),
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestCheckNoResourceAttr(resourceFullName, "policy_id"),
		),
	}

	updateStep := resource.TestStep{
		Config: testAccAPIServiceOperationConfig_Full_Update(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "api_service_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "access_control.group.groups.#", "3"),
			resource.TestMatchResourceAttr(resourceFullName, "access_control.group.groups.0.id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "access_control.group.groups.1.id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "access_control.group.groups.2.id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "access_control.permission.id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "access_control.scope.match_type", "ANY"),
			resource.TestCheckResourceAttr(resourceFullName, "access_control.scope.scopes.#", "3"),
			resource.TestMatchResourceAttr(resourceFullName, "access_control.scope.scopes.0.id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "access_control.scope.scopes.1.id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "access_control.scope.scopes.2.id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "methods.#", "4"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "methods.*", "GET"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "methods.*", "POST"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "methods.*", "PUT"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "methods.*", "DELETE"),
			resource.TestCheckResourceAttr(resourceFullName, "paths.#", "2"),
			resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "paths.*", map[string]string{
				"pattern": "/test/1",
				"type":    "EXACT",
			}),
			resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "paths.*", map[string]string{
				"pattern": "/test/{variable}/*",
				"type":    "PARAMETER",
			}),
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestCheckNoResourceAttr(resourceFullName, "policy_id"),
		),
	}

	minimalStep := resource.TestStep{
		Config: testAccAPIServiceOperationConfig_Minimal(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "api_service_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckNoResourceAttr(resourceFullName, "access_control"),
			resource.TestCheckNoResourceAttr(resourceFullName, "methods"),
			resource.TestCheckResourceAttr(resourceFullName, "paths.#", "2"),
			resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "paths.*", map[string]string{
				"pattern": "/test/1",
				"type":    "EXACT",
			}),
			resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "paths.*", map[string]string{
				"pattern": "/test/2",
				"type":    "EXACT",
			}),
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestCheckNoResourceAttr(resourceFullName, "policy_id"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             authorize.APIServiceOperation_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full from scratch
			fullStep,
			{
				Config:  testAccAPIServiceOperationConfig_Full(resourceName, name),
				Destroy: true,
			},
			// Minimal from scratch
			minimalStep,
			{
				Config:  testAccAPIServiceOperationConfig_Minimal(resourceName, name),
				Destroy: true,
			},
			// Update
			fullStep,
			minimalStep,
			fullStep,
			updateStep,
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["api_service_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAPIServiceOperation_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_api_service_operation.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             authorize.APIServiceOperation_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccAPIServiceOperationConfig_Minimal(resourceName, name),
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

func testAccAPIServiceOperationConfig_NewEnv(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_resource" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  name = "%[4]s"

  audience                      = "%[4]s"
  access_token_validity_seconds = 3600
}

resource "pingone_authorize_api_service" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  name = "%[4]s"

  base_urls = [
    "https://api.bxretail.org/%[4]s",
    "https://api.bxretail.org/%[4]s/1"
  ]

  authorization_server = {
    resource_id = pingone_resource.%[3]s.id
    type        = "PINGONE_SSO"
  }
}

resource "pingone_authorize_api_service_operation" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  api_service_id = pingone_authorize_api_service.%[3]s.id

  name = "%[4]s"

  paths = [
    {
      pattern = "/test/1"
      type    = "EXACT"
    },
    {
      pattern = "/test/2"
      type    = "EXACT"
    }
  ]
}`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccAPIServiceOperationConfig_Full(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "pingone_authorize_api_service_operation" "%[3]s" {
  environment_id = data.pingone_environment.general_test.id
  api_service_id = pingone_authorize_api_service.%[3]s.id

  name = "%[4]s"

  access_control = {
    group = {
      groups = [
        {
          id = pingone_group.%[3]s-1.id
        },
        {
          id = pingone_group.%[3]s-2.id
        },
      ]
    }

    // permission = {
    // 	id = "permissionid"
    // }

    scope = {
      match_type = "ALL"
      scopes = [
        {
          id = pingone_resource_scope.%[3]s-1.id
        },
        {
          id = pingone_resource_scope.%[3]s-2.id
        },
      ]
    }
  }

  methods = [
    "POST",
    "GET",
    "PUT",
  ]

  paths = [
    {
      pattern = "/test/1"
      type    = "EXACT"
    },
    {
      pattern = "/test/{variable}/*"
      type    = "PARAMETER"
    },
    {
      pattern = "/test/2"
      type    = "EXACT"
    }
  ]
}`, acctest.GenericSandboxEnvironment(), testAccAPIServiceOperationConfig_Full_SharedHCL(resourceName, name), resourceName, name)
}

func testAccAPIServiceOperationConfig_Full_Update(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "pingone_authorize_api_service_operation" "%[3]s" {
  environment_id = data.pingone_environment.general_test.id
  api_service_id = pingone_authorize_api_service.%[3]s.id

  name = "%[4]s"

  access_control = {
    group = {
      groups = [
        {
          id = pingone_group.%[3]s-1.id
        },
        {
          id = pingone_group.%[3]s-3.id
        },
        {
          id = pingone_group.%[3]s-2.id
        },
      ]
    }

    permission = {
      id = pingone_authorize_application_role_permission.%[3]s.application_resource_permission_id
    }

    scope = {
      match_type = "ANY"
      scopes = [
        {
          id = pingone_resource_scope.%[3]s-1.id
        },
        {
          id = pingone_resource_scope.%[3]s-3.id
        },
        {
          id = pingone_resource_scope.%[3]s-2.id
        },
      ]
    }
  }

  methods = [
    "POST",
    "PUT",
    "GET",
    "DELETE",
  ]

  paths = [
    {
      pattern = "/test/1"
      type    = "EXACT"
    },
    {
      pattern = "/test/{variable}/*"
      type    = "PARAMETER"
    },
  ]
}`, acctest.GenericSandboxEnvironment(), testAccAPIServiceOperationConfig_Full_SharedHCL(resourceName, name), resourceName, name)
}

func testAccAPIServiceOperationConfig_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "pingone_authorize_api_service_operation" "%[3]s" {
  environment_id = data.pingone_environment.general_test.id
  api_service_id = pingone_authorize_api_service.%[3]s.id

  name = "%[4]s"

  paths = [
    {
      pattern = "/test/1"
      type    = "EXACT"
    },
    {
      pattern = "/test/2"
      type    = "EXACT"
    }
  ]
}`, acctest.GenericSandboxEnvironment(), testAccAPIServiceOperationConfig_Full_SharedHCL(resourceName, name), resourceName, name)
}

func testAccAPIServiceOperationConfig_Full_SharedHCL(resourceName, name string) string {
	return fmt.Sprintf(`
resource "pingone_resource" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  audience                      = "%[3]s"
  access_token_validity_seconds = 3600
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

resource "pingone_authorize_api_service" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  base_urls = [
    "https://api.bxretail.org/%[3]s",
    "https://api.bxretail.org/%[3]s/2",
    "https://api.bxretail.org/%[3]s/1"
  ]

  authorization_server = {
    resource_id = pingone_resource.%[3]s.id
    type        = "PINGONE_SSO"
  }

  directory = {
    type = "PINGONE_SSO"
  }
}

resource "pingone_group" "%[2]s-1" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s-1"
}

resource "pingone_group" "%[2]s-2" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s-2"
}

resource "pingone_group" "%[2]s-3" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s-3"
}

resource "pingone_application_resource" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  resource_name  = pingone_resource.%[2]s.name

  name = "%[3]s"
}

resource "pingone_application_resource_permission" "%[2]s" {
  environment_id          = data.pingone_environment.general_test.id
  application_resource_id = pingone_application_resource.%[2]s.id

  action      = "%[3]s"
  description = "Test permission"
}


resource "pingone_authorize_application_role" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_authorize_application_role_permission" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  application_role_id                = pingone_authorize_application_role.%[2]s.id
  application_resource_permission_id = pingone_application_resource_permission.%[2]s.id
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
