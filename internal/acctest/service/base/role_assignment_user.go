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

func TestAccCheckRoleAssignmentUserDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_role_assignment_user" {
			continue
		}

		shouldContinue, err := acctest.CheckParentEnvironmentDestroy(ctx, apiClient, rs.Primary.Attributes["environment_id"])
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		_, r, err := apiClient.UserRoleAssignmentsApi.ReadOneUserRoleAssignment(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["user_id"], rs.Primary.ID).Execute()

		shouldContinue, err = acctest.CheckForResourceDestroy(r, err)
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		return fmt.Errorf("PingOne User Role Assignment %s still exists", rs.Primary.ID)
	}

	return nil
}

func TestAccGetRoleAssignmentUserIDs(resourceName string, environmentID, userID, resourceID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource Not found: %s", resourceName)
		}

		*resourceID = rs.Primary.ID
		*userID = rs.Primary.Attributes["user_id"]
		*environmentID = rs.Primary.Attributes["environment_id"]

		return nil
	}
}

func RoleAssignmentUser_RemovalDrift_PreConfig(ctx context.Context, apiClient *management.APIClient, t *testing.T, environmentID, userID, roleAssignmentID string) {
	if environmentID == "" || userID == "" || roleAssignmentID == "" {
		t.Fatalf("One of environment ID, user ID or resource ID cannot be determined. Environment ID: %s, User ID: %s, Role assignment ID: %s", environmentID, userID, roleAssignmentID)
	}

	_, err := apiClient.UserRoleAssignmentsApi.DeleteUserRoleAssignment(ctx, environmentID, userID, roleAssignmentID).Execute()
	if err != nil {
		t.Fatalf("Failed to delete user role assignment: %v", err)
	}
}
