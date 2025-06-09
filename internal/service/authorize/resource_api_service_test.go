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

func TestAccAPIService_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_api_service.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	var apiServiceID, environmentID string

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
		CheckDestroy:             authorize.APIService_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccAPIServiceConfig_PingOneSSO_Minimal(resourceName, name),
				Check:  authorize.APIService_GetIDs(resourceFullName, &environmentID, &apiServiceID),
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
				Config: testAccAPIServiceConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check:  authorize.APIService_GetIDs(resourceFullName, &environmentID, &apiServiceID),
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

func TestAccAPIService_NewEnv(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_api_service.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             authorize.APIService_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAPIServiceConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
				),
			},
		},
	})
}

func TestAccAPIService_PingOneSSO_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_api_service.%s", resourceName)

	name := resourceName

	fullStep := resource.TestStep{
		Config: testAccAPIServiceConfig_PingOneSSO_Full(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "access_control.custom.enabled", "true"),
			resource.TestMatchResourceAttr(resourceFullName, "authorization_server.resource_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "authorization_server.type", "PINGONE_SSO"),
			resource.TestCheckResourceAttr(resourceFullName, "base_urls.#", "3"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "base_urls.*", fmt.Sprintf("https://api.bxretail.org/%s", name)),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "base_urls.*", fmt.Sprintf("https://api.bxretail.org/%s/1", name)),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "base_urls.*", fmt.Sprintf("https://api.bxretail.org/%s/2", name)),
			resource.TestCheckResourceAttr(resourceFullName, "directory.type", "PINGONE_SSO"),
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestMatchResourceAttr(resourceFullName, "policy_id", verify.P1ResourceIDRegexpFullString),
		),
	}

	minimalStep := resource.TestStep{
		Config: testAccAPIServiceConfig_PingOneSSO_Minimal(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "access_control.custom.enabled", "true"),
			resource.TestMatchResourceAttr(resourceFullName, "authorization_server.resource_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "authorization_server.type", "PINGONE_SSO"),
			resource.TestCheckResourceAttr(resourceFullName, "base_urls.#", "2"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "base_urls.*", fmt.Sprintf("https://api.bxretail.org/%s", name)),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "base_urls.*", fmt.Sprintf("https://api.bxretail.org/%s/1", name)),
			resource.TestCheckResourceAttr(resourceFullName, "directory.type", "PINGONE_SSO"),
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestMatchResourceAttr(resourceFullName, "policy_id", verify.P1ResourceIDRegexpFullString),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             authorize.APIService_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full from scratch
			fullStep,
			{
				Config:  testAccAPIServiceConfig_PingOneSSO_Full(resourceName, name),
				Destroy: true,
			},
			// Minimal from scratch
			minimalStep,
			{
				Config:  testAccAPIServiceConfig_PingOneSSO_Minimal(resourceName, name),
				Destroy: true,
			},
			// Update
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

						return fmt.Sprintf("%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAPIService_External_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_api_service.%s", resourceName)

	name := resourceName

	fullStep := resource.TestStep{
		Config: testAccAPIServiceConfig_External_Full(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "access_control.custom.enabled", "false"),
			resource.TestCheckNoResourceAttr(resourceFullName, "authorization_server.resource_id"),
			resource.TestCheckResourceAttr(resourceFullName, "authorization_server.type", "EXTERNAL"),
			resource.TestCheckResourceAttr(resourceFullName, "base_urls.#", "3"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "base_urls.*", fmt.Sprintf("https://api.bxretail.org/%s", name)),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "base_urls.*", fmt.Sprintf("https://api.bxretail.org/%s/1", name)),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "base_urls.*", fmt.Sprintf("https://api.bxretail.org/%s/2", name)),
			resource.TestCheckResourceAttr(resourceFullName, "directory.type", "EXTERNAL"),
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestCheckNoResourceAttr(resourceFullName, "policy_id"),
		),
	}

	minimalStep := resource.TestStep{
		Config: testAccAPIServiceConfig_External_Minimal(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "access_control.custom.enabled", "false"),
			resource.TestCheckNoResourceAttr(resourceFullName, "authorization_server.resource_id"),
			resource.TestCheckResourceAttr(resourceFullName, "authorization_server.type", "EXTERNAL"),
			resource.TestCheckResourceAttr(resourceFullName, "base_urls.#", "2"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "base_urls.*", fmt.Sprintf("https://api.bxretail.org/%s", name)),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "base_urls.*", fmt.Sprintf("https://api.bxretail.org/%s/1", name)),
			resource.TestCheckResourceAttr(resourceFullName, "directory.type", "EXTERNAL"),
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
		CheckDestroy:             authorize.APIService_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full from scratch
			fullStep,
			{
				Config:  testAccAPIServiceConfig_External_Full(resourceName, name),
				Destroy: true,
			},
			// Minimal from scratch
			minimalStep,
			{
				Config:  testAccAPIServiceConfig_External_Minimal(resourceName, name),
				Destroy: true,
			},
			// Update
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

						return fmt.Sprintf("%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAPIService_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_api_service.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             authorize.APIService_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccAPIServiceConfig_PingOneSSO_Minimal(resourceName, name),
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

func testAccAPIServiceConfig_NewEnv(environmentName, licenseID, resourceName, name string) string {
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
}`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccAPIServiceConfig_PingOneSSO_Full(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_resource" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  audience                      = "%[3]s"
  access_token_validity_seconds = 3600
}

resource "pingone_authorize_api_service" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  access_control = {
    custom = {
      enabled = true
    }
  }

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
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccAPIServiceConfig_PingOneSSO_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_resource" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  audience                      = "%[3]s"
  access_token_validity_seconds = 3600
}

resource "pingone_authorize_api_service" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  access_control = {
    custom = {
      enabled = true
    }
  }

  base_urls = [
    "https://api.bxretail.org/%[3]s",
    "https://api.bxretail.org/%[3]s/1"
  ]

  authorization_server = {
    resource_id = pingone_resource.%[2]s.id
    type        = "PINGONE_SSO"
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccAPIServiceConfig_External_Full(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_authorize_api_service" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  access_control = {
    custom = {
      enabled = false
    }
  }

  base_urls = [
    "https://api.bxretail.org/%[3]s",
    "https://api.bxretail.org/%[3]s/2",
    "https://api.bxretail.org/%[3]s/1"
  ]

  authorization_server = {
    type = "EXTERNAL"
  }

  directory = {
    type = "EXTERNAL"
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccAPIServiceConfig_External_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_authorize_api_service" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  base_urls = [
    "https://api.bxretail.org/%[3]s",
    "https://api.bxretail.org/%[3]s/1"
  ]

  authorization_server = {
    type = "EXTERNAL"
  }

  directory = {
    type = "EXTERNAL"
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
