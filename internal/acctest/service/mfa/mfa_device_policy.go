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
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/legacysdk"
)

func MFADevicePolicy_CheckDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := legacysdk.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.MFAAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_mfa_device_policy" {
			continue
		}

		shouldContinue, err := legacysdk.CheckParentEnvironmentDestroy(ctx, p1Client.API.ManagementAPIClient, rs.Primary.Attributes["environment_id"])
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		_, r, err := apiClient.DeviceAuthenticationPolicyApi.ReadOneDeviceAuthenticationPolicy(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.ID).Execute()

		shouldContinue, err = acctest.CheckForResourceDestroy(r, err)
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		return fmt.Errorf("PingOne MFA Policy Instance %s still exists", rs.Primary.ID)
	}

	return nil
}

func MFADevicePolicy_GetIDs(resourceName string, environmentID, resourceID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
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

func MFADevicePolicy_RemovalDrift_PreConfig(ctx context.Context, apiClient *mfa.APIClient, t *testing.T, environmentID, mfaDevicePolicyID string) {
	if environmentID == "" || mfaDevicePolicyID == "" {
		t.Fatalf("One of environment ID or MFA device policy ID cannot be determined. Environment ID: %s, MFA Device Policy ID: %s", environmentID, mfaDevicePolicyID)
	}

	_, err := apiClient.DeviceAuthenticationPolicyApi.DeleteDeviceAuthenticationPolicy(ctx, environmentID, mfaDevicePolicyID).Execute()
	if err != nil {
		t.Fatalf("Failed to delete MFA Policy: %v", err)
	}
}

func TestCheckMFADevicePolicyApplicationMapResourceAttr(name, applicationResource, keyPattern, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[applicationResource]
		if !ok {
			return fmt.Errorf("resource not found: %s", applicationResource)
		}

		return resource.TestCheckResourceAttr(name, fmt.Sprintf(keyPattern, rs.Primary.ID), value)(s)
	}
}

func TestCheckMFADevicePolicyApplicationMapNoResourceAttr(name, applicationResource, keyPattern string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[applicationResource]
		if !ok {
			return fmt.Errorf("resource not found: %s", applicationResource)
		}

		return resource.TestCheckNoResourceAttr(name, fmt.Sprintf(keyPattern, rs.Primary.ID))(s)
	}
}
