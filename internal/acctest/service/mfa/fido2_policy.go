// Copyright © 2025 Ping Identity Corporation

package mfa

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/patrickcping/pingone-go-sdk-v2/mfa"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
)

func FIDO2Policy_CheckDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.MFAAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_mfa_fido2_policy" {
			continue
		}

		shouldContinue, err := acctest.CheckParentEnvironmentDestroy(ctx, p1Client.API.ManagementAPIClient, rs.Primary.Attributes["environment_id"])
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		_, r, err := apiClient.FIDO2PolicyApi.ReadOneFIDO2Policy(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.ID).Execute()

		shouldContinue, err = acctest.CheckForResourceDestroy(r, err)
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		return fmt.Errorf("PingOne MFA FIDO2 Policy Instance %s still exists", rs.Primary.ID)
	}

	return nil
}

func FIDO2Policy_GetIDs(resourceName string, environmentID, resourceID *string) resource.TestCheckFunc {
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

func FIDO2Policy_RemovalDrift_PreConfig(ctx context.Context, apiClient *mfa.APIClient, t *testing.T, environmentID, fido2PolicyID string) {
	if environmentID == "" || fido2PolicyID == "" {
		t.Fatalf("One of environment ID or FIDO2 Policy ID cannot be determined. Environment ID: %s, FIDO2 Policy ID: %s", environmentID, fido2PolicyID)
	}

	_, err := apiClient.FIDO2PolicyApi.DeleteFIDO2Policy(ctx, environmentID, fido2PolicyID).Execute()
	if err != nil {
		t.Fatalf("Failed to delete FIDO2 Policy: %v", err)
	}
}
