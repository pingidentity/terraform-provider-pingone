package mfa

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/patrickcping/pingone-go-sdk-v2/mfa"
)

func MFASettings_CheckDestroy(s *terraform.State) error {
	return nil
}

func MFASettings_GetIDs(resourceName string, environmentID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource Not found: %s", resourceName)
		}

		if environmentID != nil {
			*environmentID = rs.Primary.Attributes["environment_id"]
		}

		return nil
	}
}

func MFASettings_RemovalDrift_PreConfig(ctx context.Context, apiClient *mfa.APIClient, t *testing.T, environmentID string) {
	if environmentID == "" {
		t.Fatalf("The environment ID cannot be determined. Environment ID: %s", environmentID)
	}

	_, _, err := apiClient.MFASettingsApi.ResetMFASettings(ctx, environmentID).Execute()
	if err != nil {
		t.Fatalf("Failed to reset MFA settings: %v", err)
	}
}
