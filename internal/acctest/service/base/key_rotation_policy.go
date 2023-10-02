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

func TestAccCheckKeyRotationPolicyDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_key_rotation_policy" {
			continue
		}

		shouldContinue, err := acctest.CheckParentEnvironmentDestroy(ctx, apiClient, rs.Primary.Attributes["environment_id"])
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		_, r, err := apiClient.KeyRotationPoliciesApi.GetKeyRotationPolicy(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.ID).Execute()

		shouldContinue, err = acctest.CheckForResourceDestroy(r, err)
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		return fmt.Errorf("PingOne Key Rotation Policy Instance %s still exists", rs.Primary.ID)
	}

	return nil
}

func TestAccGetKeyRotationPolicyIDs(resourceName string, environmentID, resourceID *string) resource.TestCheckFunc {
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

func KeyRotationPolicy_RemovalDrift_PreConfig(ctx context.Context, apiClient *management.APIClient, t *testing.T, environmentID, keyRotationPolicyID string) {
	if environmentID == "" || keyRotationPolicyID == "" {
		t.Fatalf("One of environment ID or key rotation policy ID cannot be determined. Environment ID: %s, Key Rotation Policy ID: %s", environmentID, keyRotationPolicyID)
	}

	_, err := apiClient.KeyRotationPoliciesApi.DeleteKeyRotationPolicy(ctx, environmentID, keyRotationPolicyID).Execute()
	if err != nil {
		t.Fatalf("Failed to delete key rotation policy: %v", err)
	}
}
