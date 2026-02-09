// Copyright © 2026 Ping Identity Corporation

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

func ApplicationPushCredential_CheckDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := legacysdk.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.MFAAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_mfa_application_push_credential" {
			continue
		}

		shouldContinue, err := legacysdk.CheckParentEnvironmentDestroy(ctx, p1Client.API.ManagementAPIClient, rs.Primary.Attributes["environment_id"])
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		_, r, err := apiClient.ApplicationsApplicationMFAPushCredentialsApi.ReadOneMFAPushCredential(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["application_id"], rs.Primary.ID).Execute()

		shouldContinue, err = acctest.CheckForResourceDestroy(r, err)
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		return fmt.Errorf("PingOne Application MFA Push Credential %s still exists", rs.Primary.ID)
	}

	return nil
}

func ApplicationPushCredential_GetIDs(resourceName string, environmentID, applicationID, resourceID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		if resourceID != nil {
			*resourceID = rs.Primary.ID
		}

		if applicationID != nil {
			*applicationID = rs.Primary.Attributes["application_id"]
		}

		if environmentID != nil {
			*environmentID = rs.Primary.Attributes["environment_id"]
		}

		return nil
	}
}

func ApplicationPushCredential_RemovalDrift_PreConfig(ctx context.Context, apiClient *mfa.APIClient, t *testing.T, environmentID, applicationID, applicationPushCredentialID string) {
	if environmentID == "" || applicationID == "" || applicationPushCredentialID == "" {
		t.Fatalf("One of environment ID, application ID or application push credential ID cannot be determined. Environment ID: %s, Application ID: %s, Application Push Credential ID: %s", environmentID, applicationID, applicationPushCredentialID)
	}

	_, err := apiClient.ApplicationsApplicationMFAPushCredentialsApi.DeleteMFAPushCredential(ctx, environmentID, applicationID, applicationPushCredentialID).Execute()
	if err != nil {
		t.Fatalf("Failed to delete MFA Application push credential: %v", err)
	}
}
