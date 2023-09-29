package base

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
)

func TestAccCheckBrandingThemeDefaultDestroy(s *terraform.State) error {
	return nil
}

func TestAccGetAgreementIDs(resourceName string, environmentID, resourceID *string) resource.TestCheckFunc {
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

func Agreement_RemovalDrift_PreConfig(ctx context.Context, apiClient *management.APIClient, t *testing.T, environmentID, agreementID string) {
	if environmentID == "" || agreementID == "" {
		t.Fatalf("One of environment ID or agreement ID cannot be determined. Environment ID: %s, Agreement ID: %s", environmentID, agreementID)
	}

	_, err := apiClient.AgreementsResourcesApi.DeleteAgreement(ctx, environmentID, agreementID).Execute()
	if err != nil {
		t.Fatalf("Failed to delete agreement: %v", err)
	}
}
