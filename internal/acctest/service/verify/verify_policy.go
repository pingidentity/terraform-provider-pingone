package verify

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/patrickcping/pingone-go-sdk-v2/verify"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
)

func VerifyPolicy_CheckDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.VerifyAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_verify_policy" {
			continue
		}

		shouldContinue, err := acctest.CheckParentEnvironmentDestroy(ctx, p1Client.API.ManagementAPIClient, rs.Primary.Attributes["environment_id"])
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		_, r, err := apiClient.VerifyPoliciesApi.ReadOneVerifyPolicy(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["id"]).Execute()

		shouldContinue, err = acctest.CheckForResourceDestroy(r, err)
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		return fmt.Errorf("PingOne Verify Policy %s still exists", rs.Primary.ID)
	}

	return nil
}

func VerifyPolicy_GetIDs(resourceName string, environmentID, resourceID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource Not found: %s", resourceName)
		}

		*resourceID = rs.Primary.ID
		*environmentID = rs.Primary.Attributes["environment_id"]

		return nil
	}
}

func VerifyPolicy_RemovalDrift_PreConfig(ctx context.Context, apiClient *verify.APIClient, t *testing.T, environmentID, verifyPolicyID string) {
	if environmentID == "" || verifyPolicyID == "" {
		t.Fatalf("One of environment ID or verify policy ID cannot be determined. Environment ID: %s, Verify Policy ID: %s", environmentID, verifyPolicyID)
	}

	_, err := apiClient.VerifyPoliciesApi.DeleteVerifyPolicy(ctx, environmentID, verifyPolicyID).Execute()
	if err != nil {
		t.Fatalf("Failed to delete Verify policy: %v", err)
	}
}
