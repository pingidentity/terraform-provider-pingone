package base

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
)

func TestAccCheckRoleAssignmentGatewayDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_gateway_role_assignment" {
			continue
		}

		shouldContinue, err := acctest.CheckParentEnvironmentDestroy(ctx, apiClient, rs.Primary.Attributes["environment_id"])
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		_, r, err := apiClient.GatewayRoleAssignmentsApi.ReadOneGatewayRoleAssignment(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["gateway_id"], rs.Primary.ID).Execute()

		shouldContinue, err = acctest.CheckForResourceDestroy(r, err)
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		return fmt.Errorf("PingOne Gateway Role Assignment %s still exists", rs.Primary.ID)
	}

	return nil
}

func TestAccGetRoleAssignmentGatewayIDs(resourceName string, environmentID, gatewayID, resourceID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource Not found: %s", resourceName)
		}

		*resourceID = rs.Primary.ID
		*gatewayID = rs.Primary.Attributes["gateway_id"]
		*environmentID = rs.Primary.Attributes["environment_id"]

		return nil
	}
}

func RoleAssignmentGateway_RemovalDrift_PreConfig(ctx context.Context, apiClient *management.APIClient, t *testing.T, environmentID, gatewayID, roleAssignmentID string) {
	if environmentID == "" || gatewayID == "" || roleAssignmentID == "" {
		t.Fatalf("One of environment ID, gateway ID or role assignment ID cannot be determined. Environment ID: %s, Gateway ID: %s, Role Assignment ID: %s", environmentID, gatewayID, roleAssignmentID)
	}

	_, err := apiClient.GatewayRoleAssignmentsApi.DeleteGatewayRoleAssignment(ctx, environmentID, gatewayID, roleAssignmentID).Execute()
	if err != nil {
		t.Fatalf("Failed to delete Gateway role assignment: %v", err)
	}
}
