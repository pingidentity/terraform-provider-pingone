// Copyright © 2025 Ping Identity Corporation

package sso

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/legacysdk"
)

func ApplicationFlowPolicyAssignment_CheckDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := legacysdk.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_application_flow_policy_assignment" {
			continue
		}

		shouldContinue, err := legacysdk.CheckParentEnvironmentDestroy(ctx, p1Client.API.ManagementAPIClient, rs.Primary.Attributes["environment_id"])
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		_, r, err := apiClient.ApplicationFlowPolicyAssignmentsApi.ReadOneFlowPolicyAssignment(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["application_id"], rs.Primary.ID).Execute()

		shouldContinue, err = acctest.CheckForResourceDestroy(r, err)
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		return fmt.Errorf("PingOne Application Flow Policy assignment %s still exists", rs.Primary.ID)
	}

	return nil
}

func ApplicationFlowPolicyAssignment_GetIDs(resourceName string, environmentID, applicationID, resourceID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		if resourceID != nil {
			*resourceID = rs.Primary.ID
		}

		if applicationID != nil {
			*applicationID = rs.Primary.Attributes["application_id"]
		}

		if environmentID != nil {
			*environmentID = rs.Primary.Attributes["environment_id"]
		}

		return nil
	}
}

func ApplicationFlowPolicyAssignment_RemovalDrift_PreConfig(ctx context.Context, apiClient *management.APIClient, t *testing.T, environmentID, applicationID, applicationFlowPolicyAssignmentID string) {
	if environmentID == "" || applicationID == "" || applicationFlowPolicyAssignmentID == "" {
		t.Fatalf("One of environment ID, application ID or application flow policy assignment ID cannot be determined. Environment ID: %s, Application ID: %s, Application flow policy assignment ID: %s", environmentID, applicationID, applicationFlowPolicyAssignmentID)
	}

	_, err := apiClient.ApplicationFlowPolicyAssignmentsApi.DeleteFlowPolicyAssignment(ctx, environmentID, applicationID, applicationFlowPolicyAssignmentID).Execute()
	if err != nil {
		t.Fatalf("Failed to delete Application flow policy assignment: %v", err)
	}
}
