package authorize_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func testAccCheckDecisionEndpointDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.AuthorizeAPIClient

	apiClientManagement := p1Client.API.ManagementAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_authorize_decision_endpoint" {
			continue
		}

		_, rEnv, err := apiClientManagement.EnvironmentsApi.ReadOneEnvironment(ctx, rs.Primary.Attributes["environment_id"]).Execute()

		if err != nil {

			if rEnv == nil {
				return fmt.Errorf("Response object does not exist and no error detected")
			}

			if rEnv.StatusCode == 404 {
				continue
			}

			return err
		}

		body, r, err := apiClient.PolicyDecisionManagementApi.ReadOneDecisionEndpoint(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.ID).Execute()

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

		return fmt.Errorf("PingOne Decision Endpoint Instance %s still exists", rs.Primary.ID)
	}

	return nil
}

func testAccGetDecisionEndpointIDs(resourceName string, environmentID, resourceID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource Not found: %s", resourceName)
		}

		*resourceID = rs.Primary.ID
		*environmentID = rs.Primary.Attributes["environment_id"]

		return nil
	}
}

func TestAccDecisionEndpoint_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_decision_endpoint.%s", resourceName)

	name := resourceName

	var resourceID, environmentID string

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDecisionEndpointDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccDecisionEndpointConfig_Minimal(resourceName, name),
				Check:  testAccGetDecisionEndpointIDs(resourceFullName, &environmentID, &resourceID),
			},
			// Replan after removal preconfig
			{
				PreConfig: func() {
					var ctx = context.Background()
					p1Client, err := acctest.TestClient(ctx)

					if err != nil {
						t.Fatalf("Failed to get API client: %v", err)
					}

					apiClient := p1Client.API.AuthorizeAPIClient

					if environmentID == "" || resourceID == "" {
						t.Fatalf("One of environment ID or resource ID cannot be determined. Environment ID: %s, Resource ID: %s", environmentID, resourceID)
					}

					_, err = apiClient.PolicyDecisionManagementApi.DeleteDecisionEndpoint(ctx, environmentID, resourceID).Execute()
					if err != nil {
						t.Fatalf("Failed to delete Decision Endpoint: %v", err)
					}
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccDecisionEndpoint_NewEnv(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_decision_endpoint.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDecisionEndpointDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDecisionEndpointConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
				),
			},
		},
	})
}

func TestAccDecisionEndpoint_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_decision_endpoint.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDecisionEndpointDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDecisionEndpointConfig_Full(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", "Test description"),
					resource.TestCheckResourceAttr(resourceFullName, "owned", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "record_recent_requests", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "alternate_id", fmt.Sprintf("%s-1", name)),
					// resource.TestMatchResourceAttr(resourceFullName, "authorization_version_id", verify.P1ResourceIDRegexp),
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

						return fmt.Sprintf("%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccDecisionEndpoint_Minimal(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_decision_endpoint.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDecisionEndpointDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDecisionEndpointConfig_Minimal(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", ""),
					resource.TestCheckResourceAttr(resourceFullName, "owned", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "record_recent_requests", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "alternate_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "authorization_version_id", ""),
				),
			},
		},
	})
}

func TestAccDecisionEndpoint_Change(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_decision_endpoint.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDecisionEndpointDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDecisionEndpointConfig_Full(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", "Test description"),
					resource.TestCheckResourceAttr(resourceFullName, "owned", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "record_recent_requests", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "alternate_id", fmt.Sprintf("%s-1", name)),
					// resource.TestMatchResourceAttr(resourceFullName, "authorization_version_id", verify.P1ResourceIDRegexp),
				),
			},
			{
				Config: testAccDecisionEndpointConfig_Minimal(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", ""),
					resource.TestCheckResourceAttr(resourceFullName, "owned", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "record_recent_requests", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "alternate_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "authorization_version_id", ""),
				),
			},
			{
				Config: testAccDecisionEndpointConfig_Full(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", "Test description"),
					resource.TestCheckResourceAttr(resourceFullName, "owned", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "record_recent_requests", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "alternate_id", fmt.Sprintf("%s-1", name)),
					// resource.TestMatchResourceAttr(resourceFullName, "authorization_version_id", verify.P1ResourceIDRegexp),
				),
			},
		},
	})
}

func TestAccDecisionEndpoint_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_decision_endpoint.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDecisionEndpointDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccDecisionEndpointConfig_Minimal(resourceName, name),
			},
			// Errors
			{
				ResourceName: resourceFullName,
				ImportState:  true,
				ExpectError:  regexp.MustCompile(`Invalid import ID specified \(".*"\).  The ID should be in the format "environment_id/decision_endpoint_id" and must match regex: .*`),
			},
			{
				ResourceName:  resourceFullName,
				ImportStateId: "/",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`Invalid import ID specified \(".*"\).  The ID should be in the format "environment_id/decision_endpoint_id" and must match regex: .*`),
			},
			{
				ResourceName:  resourceFullName,
				ImportStateId: "badformat/badformat",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`Invalid import ID specified \(".*"\).  The ID should be in the format "environment_id/decision_endpoint_id" and must match regex: .*`),
			},
		},
	})
}

func testAccDecisionEndpointConfig_NewEnv(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_authorize_decision_endpoint" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s"

  record_recent_requests = false
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccDecisionEndpointConfig_Full(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_authorize_decision_endpoint" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Test description"

  record_recent_requests = true
  alternate_id           = "%[3]s-1"
  // authorization_version_id = "78e14034-ecc9-417e-9c1d-df8cafdcd04b"
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccDecisionEndpointConfig_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_authorize_decision_endpoint" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  record_recent_requests = false
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
