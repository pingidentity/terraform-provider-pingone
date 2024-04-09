package base

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
)

func IdentityPropagationPlan_CheckDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_identity_propagation_plan" {
			continue
		}

		shouldContinue, err := acctest.CheckParentEnvironmentDestroy(ctx, apiClient, rs.Primary.Attributes["environment_id"])
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		_, r, err := apiClient.IdentityPropagationPlansApi.ReadOnePlan(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.ID).Execute()

		shouldContinue, err = acctest.CheckForResourceDestroy(r, err)
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		return fmt.Errorf("PingOne Identity Propagation Plan Instance %s still exists", rs.Primary.ID)
	}

	return nil
}

func IdentityPropagationPlan_GetIDs(resourceName string, environmentID, resourceID *string) resource.TestCheckFunc {
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

func IdentityPropagationPlan_RemovalDrift_PreConfig(ctx context.Context, apiClient *management.APIClient, t *testing.T, environmentID, identityPropagationPlanID string) {
	if environmentID == "" || identityPropagationPlanID == "" {
		t.Fatalf("One of environment ID or identity propagation plan ID cannot be determined. Environment ID: %s, Identity Propagation Plan ID: %s", environmentID, identityPropagationPlanID)
	}

	_, err := apiClient.IdentityPropagationPlansApi.DeletePlan(ctx, environmentID, identityPropagationPlanID).Execute()
	if err != nil {
		t.Fatalf("Failed to delete Identity Propagation Plan: %v", err)
	}
}
