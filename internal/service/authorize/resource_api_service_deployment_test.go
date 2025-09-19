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

func TestAccAPIServiceDeployment_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_api_service_deployment.%s", resourceName)

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
			acctest.PreCheckTestAccFlaky(t)
			acctest.PreCheckNoBeta(t)

			p1Client = acctestlegacysdk.PreCheckTestClient(ctx, t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             authorize.APIServiceDeployment_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Test removal of the API service
			{
				Config: testAccAPIServiceDeploymentConfig_Full(resourceName, name),
				Check:  authorize.APIServiceDeployment_GetIDs(resourceFullName, &environmentID, &apiServiceID),
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
				Config: testAccAPIServiceDeploymentConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check:  authorize.APIServiceDeployment_GetIDs(resourceFullName, &environmentID, &apiServiceID),
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

func TestAccAPIServiceDeployment_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_api_service_deployment.%s", resourceName)

	name := resourceName

	fullStep := resource.TestStep{
		Config: testAccAPIServiceDeploymentConfig_Full(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "api_service_id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "authorization_version.id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "decision_endpoint.id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "deployed_at", verify.RFC3339Regexp),
			resource.TestCheckNoResourceAttr(resourceFullName, "policy.id"),
			resource.TestCheckResourceAttr(resourceFullName, "status.code", "DEPLOYMENT_SUCCESSFUL"),
			resource.TestCheckNoResourceAttr(resourceFullName, "status.error"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             authorize.APIServiceDeployment_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full from scratch
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

						return fmt.Sprintf("%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["api_service_id"]), nil
					}
				}(),
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "api_service_id",
			},
		},
	})
}

func TestAccAPIServiceDeployment_ReplaceTriggers(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_api_service_deployment.%s", resourceName)

	name := resourceName

	triggerAStep1 := resource.TestStep{
		Config: testAccAPIServiceDeploymentConfig_ReplaceTriggerA1(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "redeployment_trigger_values.triggerA", "triggerAValue1"),
			resource.TestCheckNoResourceAttr(resourceFullName, "redeployment_trigger_values.triggerB"),
		),
	}

	triggerAStep2 := resource.TestStep{
		Config: testAccAPIServiceDeploymentConfig_ReplaceTriggerA2(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "redeployment_trigger_values.triggerA", "triggerAValue2"),
			resource.TestCheckNoResourceAttr(resourceFullName, "redeployment_trigger_values.triggerB"),
		),
	}

	addTriggerBStep := resource.TestStep{
		Config: testAccAPIServiceDeploymentConfig_ReplaceTriggerB1(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "redeployment_trigger_values.triggerA", "triggerAValue2"),
			resource.TestCheckResourceAttr(resourceFullName, "redeployment_trigger_values.triggerB", "triggerBValue1"),
		),
	}

	removeTriggerBStep := resource.TestStep{
		Config: testAccAPIServiceDeploymentConfig_ReplaceTriggerB2(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "redeployment_trigger_values.triggerA", "triggerAValue2"),
			resource.TestCheckNoResourceAttr(resourceFullName, "redeployment_trigger_values.triggerB"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             authorize.APIServiceDeployment_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			triggerAStep1,
			triggerAStep2,
			addTriggerBStep,
			removeTriggerBStep,
		},
	})
}

func TestAccAPIServiceDeployment_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_api_service_deployment.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             authorize.APIServiceDeployment_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccAPIServiceDeploymentConfig_Full(resourceName, name),
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

func testAccAPIServiceDeploymentConfig_NewEnv(environmentName, licenseID, resourceName, name string) string {
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

resource "pingone_authorize_api_service_deployment" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  api_service_id = pingone_authorize_api_service.%[3]s.id
}`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccAPIServiceDeploymentConfig_Full(resourceName, name string) string {
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
}

resource "pingone_authorize_api_service_deployment" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  api_service_id = pingone_authorize_api_service.%[2]s.id
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccAPIServiceDeploymentConfig_ReplaceTriggerA1(resourceName, name string) string {
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
}

resource "pingone_authorize_api_service_deployment" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  api_service_id = pingone_authorize_api_service.%[2]s.id

  redeployment_trigger_values = {
    "triggerA" : "triggerAValue1",
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccAPIServiceDeploymentConfig_ReplaceTriggerA2(resourceName, name string) string {
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
}

resource "pingone_authorize_api_service_deployment" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  api_service_id = pingone_authorize_api_service.%[2]s.id

  redeployment_trigger_values = {
    "triggerA" : "triggerAValue2",
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccAPIServiceDeploymentConfig_ReplaceTriggerB1(resourceName, name string) string {
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
}

resource "pingone_authorize_api_service_deployment" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  api_service_id = pingone_authorize_api_service.%[2]s.id

  redeployment_trigger_values = {
    "triggerA" : "triggerAValue2",
    "triggerB" : "triggerBValue1",
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccAPIServiceDeploymentConfig_ReplaceTriggerB2(resourceName, name string) string {
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
}

resource "pingone_authorize_api_service_deployment" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  api_service_id = pingone_authorize_api_service.%[2]s.id

  redeployment_trigger_values = {
    "triggerA" : "triggerAValue2",
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
