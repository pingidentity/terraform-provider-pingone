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

func RoleAssignmentGroup_CheckDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_group_role_assignment" {
			continue
		}

		shouldContinue, err := acctest.CheckParentEnvironmentDestroy(ctx, p1Client.API.ManagementAPIClient, rs.Primary.Attributes["environment_id"])
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		_, r, err := apiClient.GroupRoleAssignmentsApi.ReadOneGroupRoleAssignment(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["group_id"], rs.Primary.ID).Execute()

		shouldContinue, err = acctest.CheckForResourceDestroy(r, err)
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		return fmt.Errorf("PingOne Group Role Assignment %s still exists", rs.Primary.ID)
	}

	return nil
}

func RoleAssignmentGroup_GetIDs(resourceName string, environmentID, groupID, resourceID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource Not found: %s", resourceName)
		}

		if resourceID != nil {
			*resourceID = rs.Primary.ID
		}

		if groupID != nil {
			*groupID = rs.Primary.Attributes["group_id"]
		}

		if environmentID != nil {
			*environmentID = rs.Primary.Attributes["environment_id"]
		}

		return nil
	}
}

func RoleAssignmentGroup_RemovalDrift_PreConfig(ctx context.Context, apiClient *management.APIClient, t *testing.T, environmentID, groupID, groupRoleAssignmentID string) {
	if environmentID == "" || groupID == "" || groupRoleAssignmentID == "" {
		t.Fatalf("One of environment ID, group ID or group role assignment ID cannot be determined. Environment ID: %s, Group ID: %s, Group Role Assignment ID: %s", environmentID, groupID, groupRoleAssignmentID)
	}

	_, err := apiClient.GroupRoleAssignmentsApi.DeleteGroupRoleAssignment(ctx, environmentID, groupID, groupRoleAssignmentID).Execute()
	if err != nil {
		t.Fatalf("Failed to delete Group role assignment: %v", err)
	}
}
