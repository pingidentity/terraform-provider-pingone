// Copyright © 2025 Ping Identity Corporation

package risk

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/patrickcping/pingone-go-sdk-v2/risk"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
)

func RiskPolicy_CheckDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.RiskAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_risk_policy" {
			continue
		}

		shouldContinue, err := acctest.CheckParentEnvironmentDestroy(ctx, p1Client.API.ManagementAPIClient, rs.Primary.Attributes["environment_id"])
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		_, r, err := apiClient.RiskPoliciesApi.ReadOneRiskPolicySet(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.ID).Execute()

		shouldContinue, err = acctest.CheckForResourceDestroy(r, err)
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		return fmt.Errorf("PingOne risk policy %s still exists", rs.Primary.ID)
	}

	return nil
}

func RiskPolicy_GetIDs(resourceName string, environmentID, resourceID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource Not found: %s", resourceName)
		}

		if resourceID != nil {
			*resourceID = rs.Primary.ID
		}

		if environmentID != nil {
			*environmentID = rs.Primary.Attributes["environment_id"]
		}

		return nil
	}
}

func RiskPolicy_RemovalDrift_PreConfig(ctx context.Context, apiClient *risk.APIClient, t *testing.T, environmentID, riskPolicyID string) {
	if environmentID == "" || riskPolicyID == "" {
		t.Fatalf("One of environment ID or risk policy ID cannot be determined. Environment ID: %s, Risk Policy ID: %s", environmentID, riskPolicyID)
	}

	_, err := apiClient.RiskPoliciesApi.DeleteRiskPolicySet(ctx, environmentID, riskPolicyID).Execute()
	if err != nil {
		t.Fatalf("Failed to delete Risk Policy: %v", err)
	}
}
