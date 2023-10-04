package risk

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/patrickcping/pingone-go-sdk-v2/risk"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
)

func TestAccCheckRiskPredictorDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.RiskAPIClient

	apiClientManagement := p1Client.API.ManagementAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_risk_predictor" {
			continue
		}

		_, rEnv, err := apiClientManagement.EnvironmentsApi.ReadOneEnvironment(ctx, rs.Primary.Attributes["environment_id"]).Execute()

		if err != nil {

			if rEnv == nil {
				return fmt.Errorf("Response object does not exist and no error detected")
			}

			if rEnv.StatusCode == 404 {
				continue
			}

			return err
		}

		body, r, err := apiClient.RiskAdvancedPredictorsApi.ReadOneRiskPredictor(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.ID).Execute()

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

		return fmt.Errorf("PingOne risk predictor %s still exists", rs.Primary.ID)
	}

	return nil
}

func TestAccCheckRiskPredictorDestroyUndeletable(s *terraform.State) error {
	return nil
}

func TestAccGetRiskPredictorIDs(resourceName string, environmentID, resourceID *string) resource.TestCheckFunc {
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

func RiskPredictor_RemovalDrift_PreConfig(ctx context.Context, apiClient *risk.APIClient, t *testing.T, environmentID, riskPredictorID string) {
	if environmentID == "" || riskPredictorID == "" {
		t.Fatalf("One of environment ID or risk predictor ID cannot be determined. Environment ID: %s, Risk Predictor ID: %s", environmentID, riskPredictorID)
	}

	_, err := apiClient.RiskAdvancedPredictorsApi.DeleteRiskAdvancedPredictor(ctx, environmentID, riskPredictorID).Execute()
	if err != nil {
		t.Fatalf("Failed to delete Risk Predictor: %v", err)
	}
}
