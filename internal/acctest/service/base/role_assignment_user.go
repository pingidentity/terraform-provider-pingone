package base

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tflog"
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

		_, rEnv, err := apiClient.EnvironmentsApi.ReadOneEnvironment(ctx, rs.Primary.Attributes["environment_id"]).Execute()

		if err != nil {

			if rEnv == nil {
				return fmt.Errorf("Response object does not exist and no error detected")
			}

			if rEnv.StatusCode == 404 {
				continue
			}

			return err
		}

		body, r, err := apiClient.UserRoleAssignmentsApi.ReadOneUserRoleAssignment(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["user_id"], rs.Primary.ID).Execute()

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
