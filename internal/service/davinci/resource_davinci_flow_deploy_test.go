// Copyright © 2025 Ping Identity Corporation

//go:build beta

package davinci_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/pingidentity/pingone-go-client/pingone"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	acctestlegacysdk "github.com/pingidentity/terraform-provider-pingone/internal/acctest/legacysdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/base"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/testhcl"
)

func TestAccDavinciFlowDeploy_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_davinci_flow_deploy.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")
	var environmentId string
	var id string

	var p1Client *pingone.APIClient
	var ctx = context.Background()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckBeta(t)

			p1Client = acctest.PreCheckTestClient(ctx, t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             davinciFlow_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Test removal of the flow
			{
				Config: davinciFlowDeploy_FirstDeployHCL(t, resourceName, false),
				Check:  davinciFlowDeploy_GetIDs(resourceFullName, &environmentId, &id),
			},
			{
				PreConfig: func() {
					davinciFlow_Delete(ctx, p1Client, t, environmentId, id)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the environment
			{
				Config: davinciFlowDeploy_NewEnvHCL(environmentName, licenseID, resourceName),
				Check:  davinciFlowDeploy_GetIDs(resourceFullName, &environmentId, &id),
			},
			{
				PreConfig: func() {
					base.Environment_RemovalDrift_PreConfig(ctx, p1Client, t, environmentId)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccDavinciFlowDeploy_Clean(t *testing.T) {
	testAccDavinciFlowDeploy(t, false)
}

func TestAccDavinciFlowDeploy_WithBootstrap(t *testing.T) {
	testAccDavinciFlowDeploy(t, true)
}

func testAccDavinciFlowDeploy(t *testing.T, withBootstrap bool) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	var lastDeployTime string

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             davinciFlow_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				// Create the flow
				Config: davinciFlowDeploy_FlowOnlyHCL(t, resourceName, withBootstrap),
			},
			{
				// Initial deploy on create
				Config: davinciFlowDeploy_FirstDeployHCL(t, resourceName, withBootstrap),
				Check: resource.ComposeTestCheckFunc(
					davinciFlowDeploy_checkExpectedDeployTimestamp(true, &lastDeployTime),
					davinciFlowDeploy_CheckComputedValues(resourceName),
					davinciFlowDeploy_GetDeployedTimestamp(&lastDeployTime),
				),
			},
			{
				// Expect no additional deploy
				Config: davinciFlowDeploy_FirstNoDeployHCL(t, resourceName, withBootstrap),
				Check: resource.ComposeTestCheckFunc(
					davinciFlowDeploy_checkExpectedDeployTimestamp(false, &lastDeployTime),
					davinciFlowDeploy_CheckComputedValues(resourceName),
					davinciFlowDeploy_GetDeployedTimestamp(&lastDeployTime),
				),
			},
			{
				// Expect deploy
				Config: davinciFlowDeploy_SecondDeployHCL(t, resourceName, withBootstrap),
				Check: resource.ComposeTestCheckFunc(
					davinciFlowDeploy_checkExpectedDeployTimestamp(true, &lastDeployTime),
					davinciFlowDeploy_CheckComputedValues(resourceName),
					davinciFlowDeploy_GetDeployedTimestamp(&lastDeployTime),
				),
			},
			{
				// Expect no additional deploy
				Config: davinciFlowDeploy_SecondNoDeployHCL(t, resourceName, withBootstrap),
				Check: resource.ComposeTestCheckFunc(
					davinciFlowDeploy_checkExpectedDeployTimestamp(false, &lastDeployTime),
					davinciFlowDeploy_CheckComputedValues(resourceName),
					davinciFlowDeploy_GetDeployedTimestamp(&lastDeployTime),
				),
			},
			// Import is not supported in this resource
		},
	})
}

func TestAccDavinciFlowDeploy_BrokenFlow_Clean(t *testing.T) {
	testAccDavinciFlowDeploy_BrokenFlow(t, false)
}

func TestAccDavinciFlowDeploy_BrokenFlow_WithBootstrap(t *testing.T) {
	testAccDavinciFlowDeploy_BrokenFlow(t, true)
}

func testAccDavinciFlowDeploy_BrokenFlow(t *testing.T, withBootstrap bool) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             davinciFlow_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: davinciFlowDeploy_BrokenFlowHCL(t, resourceName, withBootstrap),
				// TRIAGE-31546: Right now attempting to deploy this broken flow returns a 500 error from the API
				ExpectError: regexp.MustCompile(`There was an unexpected error with the service`),
			},
		},
	})
}

func davinciFlowDeploy_GetIDs(resourceName string, environmentId, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource Not found: %s", resourceName)
		}
		if environmentId != nil {
			*environmentId = rs.Primary.Attributes["environment_id"]
		}
		if id != nil {
			*id = rs.Primary.Attributes["id"]
		}

		return nil
	}
}

func davinciFlowDeploy_GetDeployedTimestamp(lastDeploy *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var ctx = context.Background()

		p1Client, err := acctest.TestClient(ctx)

		if err != nil {
			return err
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "pingone_davinci_flow" {
				continue
			}

			flowResponse, _, err := p1Client.DaVinciFlowsApi.GetFlowById(ctx, uuid.MustParse(rs.Primary.Attributes["environment_id"]), rs.Primary.Attributes["id"]).Execute()
			if err != nil {
				return err
			}

			if flowResponse != nil && flowResponse.DeployedAt != nil {
				*lastDeploy = flowResponse.DeployedAt.Format(time.RFC3339)
			} else {
				return fmt.Errorf("unable to determine last deployed time for flow %s", rs.Primary.ID)
			}
		}

		return nil
	}
}

func davinciFlowDeploy_CheckComputedValues(resourceName string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(fmt.Sprintf("pingone_davinci_flow_deploy.%s", resourceName), "id"),
	)
}

func davinciFlowDeploy_checkExpectedDeployTimestamp(expectRedeploy bool, lastDeployTime *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var currentDeployTimestamp string
		err := davinciFlowDeploy_GetDeployedTimestamp(&currentDeployTimestamp)(s)
		if err != nil {
			return err
		}
		if lastDeployTime == nil {
			return fmt.Errorf("lastDeployTime pointer is nil")
		}
		if currentDeployTimestamp != *lastDeployTime && !expectRedeploy {
			return fmt.Errorf("Current flow was redeployed unexpectedly. Previous deploy time: %s, Current deploy time: %s", *lastDeployTime, currentDeployTimestamp)
		} else if currentDeployTimestamp == *lastDeployTime && expectRedeploy {
			return fmt.Errorf("Expected the current flow to have been redeployed, but it was not. Previous deploy time: %s, Current deploy time: %s", *lastDeployTime, currentDeployTimestamp)
		}
		return nil
	}
}

func davinciFlowDeploy_FlowOnlyHCL(t *testing.T, resourceName string, withBootstrap bool) string {
	hcl, err := testhcl.ReadTestHcl("pingone_davinci_flow/ootb_device_management.tf")
	if err != nil {
		t.Fatalf("failed to read HCL in davinciFlow_DeviceManagementMainFlowHCL: %v", err)
	}
	return fmt.Sprintf(hcl, acctest.DaVinciSandboxEnvironment(withBootstrap), resourceName)
}

func davinciFlowDeploy_FirstDeployHCL(t *testing.T, resourceName string, withBootstrap bool) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_davinci_flow_deploy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  flow_id        = pingone_davinci_flow.%[2]s.id
}
`, davinciFlowDeploy_FlowOnlyHCL(t, resourceName, withBootstrap), resourceName)
}

// Ensure that adding triggers doesn't cause a redeploy
func davinciFlowDeploy_FirstNoDeployHCL(t *testing.T, resourceName string, withBootstrap bool) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_davinci_flow_deploy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  flow_id        = pingone_davinci_flow.%[2]s.id
  deploy_trigger_values = {
    "trigger" = "initial"
  }
}
`, davinciFlowDeploy_FlowOnlyHCL(t, resourceName, withBootstrap), resourceName)
}

func davinciFlowDeploy_SecondDeployHCL(t *testing.T, resourceName string, withBootstrap bool) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_davinci_flow_deploy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  flow_id        = pingone_davinci_flow.%[2]s.id
  deploy_trigger_values = {
    "trigger"    = "updated"
    "newtrigger" = "new"
  }
}
`, davinciFlowDeploy_FlowOnlyHCL(t, resourceName, withBootstrap), resourceName)
}

// Ensure that removing triggers doesn't cause a redeploy
func davinciFlowDeploy_SecondNoDeployHCL(t *testing.T, resourceName string, withBootstrap bool) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_davinci_flow_deploy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  flow_id        = pingone_davinci_flow.%[2]s.id
  deploy_trigger_values = {
    "trigger" = "updated"
  }
}
`, davinciFlowDeploy_FlowOnlyHCL(t, resourceName, withBootstrap), resourceName)
}

func davinciFlowDeploy_NewEnvHCL(environmentName, licenseID, resourceName string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_davinci_flow" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[3]s"
  description    = "This is a demo flow"
  color          = "#00FF00"

  graph_data = {
    elements = {
      nodes = {
        "8bnj41592a" = {
          data = {
            id              = "8bnj41592a"
            node_type       = "CONNECTION"
            connector_id    = "pingOneSSOConnector"
            label           = "PingOne"
            status          = "configured"
            capability_name = "userLookup"
            type            = "action"
            properties = jsonencode({
              "additionalUserProperties" : {
                "value" : []
              },
              "username" : {
                "value" : "[\n  {\n    \"children\": [\n      {\n        \"text\": \"5282e30d-6e05-499c-ae68-0069fba776f1\"\n      }\n    ]\n  }\n]"
              },
              "population" : {
                "value" : "c9f3fb3f-11e9-4eb0-b4ba-9fb7789a8418"
              },
              "userIdentifierForFindUser" : {
                "value" : "[\n  {\n    \"children\": [\n      {\n        \"text\": \"5282e30d-6e05-499c-ae68-0069fba776f1\"\n      }\n    ]\n  }\n]"
              }
            })
          }
          position = {
            x = 420
            y = 360
          }
          group      = "nodes"
          removed    = false
          selected   = false
          selectable = true
          locked     = false
          grabbable  = true
          pannable   = false
          classes    = ""
        }
      }
    }

    data = "{}"

    box_selection_enabled = true
    user_zooming_enabled  = true
    zooming_enabled       = true
    zoom                  = 1
    min_zoom              = 0.01
    max_zoom              = 10000
    panning_enabled       = true
    user_panning_enabled  = true

    pan = {
      x = 0
      y = 0
    }

    renderer = jsonencode({
      "name" : "null"
    })
  }

  output_schema = {
    output = {
      type                 = "object"
      properties           = jsonencode({})
      additionalProperties = true
    }
  }

  trigger = {
    type = "AUTHENTICATION"
  }
}

resource "pingone_davinci_flow_deploy" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  flow_id        = pingone_davinci_flow.%[3]s.id
  deploy_trigger_values = {
    "trigger" = "initial"
  }
}
`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName)
}

func davinciFlowDeploy_BrokenFlowHCL(t *testing.T, resourceName string, withBootstrap bool) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_davinci_connector_instance" "%[2]s-errors" {
  environment_id = data.pingone_environment.general_test.id
  connector = {
    id = "errorConnector"
  }
  name = "%[2]s-errors"
}

resource "pingone_davinci_flow" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name  = "brokenFlow"
  color = "#FFC8C1"

  graph_data = {
    elements = {
      nodes = {
        "zlgncqvws2" = {
          data = {
            id            = "zlgncqvws2"
            node_type     = "CONNECTION"
            connection_id = pingone_davinci_connector_instance.%[2]s-errors.id
            connector_id  = "errorConnector"
            name          = "Error Message"
            label         = "Error Message"
            status        = "unconfigured"
          }
          position = {
            x = 400
            y = 409.0909090042114
          }
          group      = "nodes"
          removed    = false
          selected   = false
          selectable = true
          locked     = false
          grabbable  = true
          pannable   = false
          classes    = ""
        }
      }
    }
    data                 = "{}"
    zooming_enabled      = true
    user_zooming_enabled = true
    zoom                 = 1
    min_zoom             = 1e-50
    max_zoom             = 1e+50
    panning_enabled      = true
    user_panning_enabled = true
    pan = {
      x = 0
      y = 0
    }
    box_selection_enabled = true
    renderer = jsonencode({
      "name" : "null"
    })
  }
}

resource "pingone_davinci_flow_deploy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  flow_id        = pingone_davinci_flow.%[2]s.id
}
`, acctest.DaVinciSandboxEnvironment(withBootstrap), resourceName)
}
