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

func UserApplicationRoleAssignment_CheckDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_user_application_role_assignment" {
			continue
		}

		shouldContinue, err := acctest.CheckParentEnvironmentDestroy(ctx, p1Client.API.ManagementAPIClient, rs.Primary.Attributes["environment_id"])
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		response, r, err := apiClient.UserApplicationRoleAssignmentsApi.ReadUserApplicationRoleAssignments(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["user_id"]).Execute()

		shouldContinue, err = acctest.CheckForResourceDestroy(r, err)
		if err != nil {
			return err
		}

		if response == nil {
			return fmt.Errorf("PingOne User Application Role Assignment list cannot be found")
		}

		for _, roleAssignment := range response.Embedded.GetRoles() {
			if v, ok := roleAssignment.UserApplicationRoleAssignment.GetIdOk(); ok && v != nil && *v == rs.Primary.Attributes["application_role_id"] {
				shouldContinue = false
			}
		}

		if shouldContinue {
			continue
		}

		return fmt.Errorf("PingOne User Application Role Assignment %s still exists", rs.Primary.ID)
	}

	return nil
}

func UserApplicationRoleAssignment_GetIDs(resourceName string, environmentID, userID, resourceID *string) resource.TestCheckFunc {
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

func UserApplicationRoleAssignment_RemovalDrift_PreConfig(ctx context.Context, apiClient *management.APIClient, t *testing.T, environmentID, userID, applicationRoleID string) {
	if environmentID == "" || userID == "" || applicationRoleID == "" {
		t.Fatalf("One of environment ID, user ID or application role ID cannot be determined. Environment ID: %s, Application ID: %s, Application Role ID: %s", environmentID, userID, applicationRoleID)
	}

	_, err := apiClient.UserApplicationRoleAssignmentsApi.DeleteUserApplicationRoleAssignment(ctx, environmentID, userID, applicationRoleID).Execute()
	if err != nil {
		t.Fatalf("Failed to delete User application role assignment: %v", err)
	}
}
