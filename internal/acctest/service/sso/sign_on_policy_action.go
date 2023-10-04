package sso

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
)

func SignOnPolicyAction_CheckDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_sign_on_policy_action" {
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

		body, r, err := apiClient.SignOnPolicyActionsApi.ReadOneSignOnPolicyAction(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["sign_on_policy_id"], rs.Primary.ID).Execute()

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

		return fmt.Errorf("PingOne Sign on Policy Action %s still exists", rs.Primary.ID)
	}

	return nil
}

func SignOnPolicyAction_GetIDs(resourceName string, environmentID, signOnPolicyID, resourceID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource Not found: %s", resourceName)
		}

		*resourceID = rs.Primary.ID
		*signOnPolicyID = rs.Primary.Attributes["sign_on_policy_id"]
		*environmentID = rs.Primary.Attributes["environment_id"]

		return nil
	}
}

func SignOnPolicyAction_RemovalDrift_PreConfig(ctx context.Context, apiClient *management.APIClient, t *testing.T, environmentID, signOnPolicyID, signOnPolicyActionID string) {
	if environmentID == "" || signOnPolicyID == "" || signOnPolicyActionID == "" {
		t.Fatalf("One of environment ID, sign-on policy ID or sign-on policy action ID cannot be determined. Environment ID: %s, Sign-on policy ID: %s, Sign-on policy action ID: %s", environmentID, signOnPolicyID, signOnPolicyActionID)
	}

	_, err := apiClient.SignOnPolicyActionsApi.DeleteSignOnPolicyAction(ctx, environmentID, signOnPolicyID, signOnPolicyActionID).Execute()
	if err != nil {
		t.Fatalf("Failed to delete sign-on policy action: %v", err)
	}
}
