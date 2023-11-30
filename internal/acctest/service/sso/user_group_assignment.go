package sso

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
)

func UserGroupAssignment_CheckDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_user_group_assignment" {
			continue
		}

		shouldContinue, err := acctest.CheckParentEnvironmentDestroy(ctx, p1Client.API.ManagementAPIClient, rs.Primary.Attributes["environment_id"])
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		_, r, err := apiClient.GroupMembershipApi.ReadOneGroupMembershipForUser(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["user_id"], rs.Primary.Attributes["group_id"]).Execute()

		shouldContinue, err = acctest.CheckForResourceDestroy(r, err)
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		return fmt.Errorf("PingOne User Group Assignment %s still exists", rs.Primary.ID)
	}

	return nil
}

func UserGroupAssignment_GetIDs(resourceName string, environmentID, userID, groupID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource Not found: %s", resourceName)
		}

		*environmentID = rs.Primary.Attributes["environment_id"]
		*userID = rs.Primary.Attributes["user_id"]
		*groupID = rs.Primary.Attributes["group_id"]

		return nil
	}
}

func UserGroupAssignment_RemovalDrift_PreConfig(ctx context.Context, apiClient *management.APIClient, t *testing.T, environmentID, userID, groupID string) {
	if environmentID == "" || userID == "" || groupID == "" {
		t.Fatalf("One of environment ID, user ID or group ID cannot be determined. Environment ID: %s, User ID: %s, Group ID: %s", environmentID, userID, groupID)
	}

	_, err := apiClient.GroupMembershipApi.RemoveUserFromGroup(ctx, environmentID, userID, groupID).Execute()
	if err != nil {
		t.Fatalf("Failed to delete user group membership: %v", err)
	}
}
