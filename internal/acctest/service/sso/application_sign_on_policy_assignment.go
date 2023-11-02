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

func ApplicationSignOnPolicyAssignment_CheckDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_application_sign_on_policy_assignment" {
			continue
		}

		shouldContinue, err := acctest.CheckParentEnvironmentDestroy(ctx, apiClient, rs.Primary.Attributes["environment_id"])
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		_, r, err := apiClient.ApplicationSignOnPolicyAssignmentsApi.ReadOneSignOnPolicyAssignment(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["application_id"], rs.Primary.ID).Execute()

		shouldContinue, err = acctest.CheckForResourceDestroy(r, err)
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		return fmt.Errorf("PingOne Application Sign On Policy assignment %s still exists", rs.Primary.ID)
	}

	return nil
}

func ApplicationSignOnPolicyAssignment_GetIDs(resourceName string, environmentID, applicationID, signOnPolicyID, signOnPolicyAssignmentID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource Not found: %s", resourceName)
		}

		*signOnPolicyAssignmentID = rs.Primary.ID
		*applicationID = rs.Primary.Attributes["application_id"]
		*environmentID = rs.Primary.Attributes["environment_id"]
		*signOnPolicyID = rs.Primary.Attributes["sign_on_policy_id"]

		return nil
	}
}

func ApplicationSignOnPolicyAssignment_RemovalDrift_PreConfig(ctx context.Context, apiClient *management.APIClient, t *testing.T, environmentID, applicationID, signOnPolicyAssignmentID string) {
	if environmentID == "" || applicationID == "" || signOnPolicyAssignmentID == "" {
		t.Fatalf("One of environment ID, application ID or resource ID cannot be determined. Environment ID: %s, Application ID: %s, Sign On Policy Assignment ID: %s", environmentID, applicationID, signOnPolicyAssignmentID)
	}

	_, err := apiClient.ApplicationSignOnPolicyAssignmentsApi.DeleteSignOnPolicyAssignment(ctx, environmentID, applicationID, signOnPolicyAssignmentID).Execute()
	if err != nil {
		t.Fatalf("Failed to delete application sign-on policy assignment: %v", err)
	}
}
