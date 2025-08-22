// Copyright © 2025 Ping Identity Corporation

package base

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

func Agreement_CheckDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := legacysdk.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_agreement" {
			continue
		}

		shouldContinue, err := legacysdk.CheckParentEnvironmentDestroy(ctx, apiClient, rs.Primary.Attributes["environment_id"])
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		_, r, err := apiClient.AgreementsResourcesApi.ReadOneAgreement(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.ID).Execute()

		shouldContinue, err = acctest.CheckForResourceDestroy(r, err)
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		return fmt.Errorf("PingOne agreement %s still exists", rs.Primary.ID)
	}

	return nil
}

func Agreement_GetIDs(resourceName string, environmentID, resourceID *string) resource.TestCheckFunc {
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

func Agreement_RemovalDrift_PreConfig(ctx context.Context, apiClient *management.APIClient, t *testing.T, environmentID, agreementID string) {
	if environmentID == "" || agreementID == "" {
		t.Fatalf("One of environment ID or agreement ID cannot be determined. Environment ID: %s, Agreement ID: %s", environmentID, agreementID)
	}

	agreement, r, err := apiClient.AgreementsResourcesApi.ReadOneAgreement(ctx, environmentID, agreementID).Execute()
	if err != nil || r.StatusCode != 200 {
		t.Fatalf("Failed to read agreement to delete: status code: %d, error: %v", r.StatusCode, err)
	}

	if agreement == nil {
		t.Fatalf("Agreement not found.  Cannot disable before delete")
	}

	agreement.SetEnabled(false)

	_, r, err = apiClient.AgreementsResourcesApi.UpdateAgreement(ctx, environmentID, agreementID).Agreement(*agreement).Execute()
	if err != nil || r.StatusCode != 200 {
		t.Fatalf("Failed to disable agreement prior to delete: status code: %d, error: %v", r.StatusCode, err)
	}

	_, err = apiClient.AgreementsResourcesApi.DeleteAgreement(ctx, environmentID, agreementID).Execute()
	if err != nil {
		t.Fatalf("Failed to delete agreement: %v", err)
	}
}
