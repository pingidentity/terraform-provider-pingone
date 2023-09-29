package sso

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

func TestAccCheckApplicationSignOnPolicyAssignmentDestroy(s *terraform.State) error {
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

		body, r, err := apiClient.ApplicationSignOnPolicyAssignmentsApi.ReadOneSignOnPolicyAssignment(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["application_id"], rs.Primary.ID).Execute()

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

		return fmt.Errorf("PingOne Application Sign On Policy assignment %s still exists", rs.Primary.ID)
	}

	return nil
}

func TestAccGetApplicationSignOnPolicyAssignmentIDs(resourceName string, environmentID, applicationID, signOnPolicyID, signOnPolicyAssignmentID *string) resource.TestCheckFunc {
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
