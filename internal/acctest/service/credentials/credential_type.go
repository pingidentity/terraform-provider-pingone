// Copyright © 2026 Ping Identity Corporation

package credentials

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/patrickcping/pingone-go-sdk-v2/credentials"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/legacysdk"
)

func CredentialType_CheckDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := legacysdk.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.CredentialsAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_credential_type" {
			continue
		}

		shouldContinue, err := legacysdk.CheckParentEnvironmentDestroy(ctx, p1Client.API.ManagementAPIClient, rs.Primary.Attributes["environment_id"])
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		body, r, err := apiClient.CredentialTypesApi.ReadOneCredentialType(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["id"]).Execute()

		shouldContinue, err = acctest.CheckForResourceDestroy(r, err)
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		} else {

			if body.DeletedAt != nil {

				// Note: Credential Types are "soft delted" and may be returned via the ReadOneCredentialType call.
				// If the DeletedAt attribute exists, it is considered deleted, handle similar to a 404.
				return err
			}

		}

		return fmt.Errorf("PingOne Credential Type ID %s still exists", rs.Primary.ID)
	}

	return nil
}

func CredentialType_GetIDs(resourceName string, environmentID, resourceID *string) resource.TestCheckFunc {
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

func CredentialType_RemovalDrift_PreConfig(ctx context.Context, apiClient *credentials.APIClient, t *testing.T, environmentID, credentialTypeID string) {
	if environmentID == "" || credentialTypeID == "" {
		t.Fatalf("One of environment ID or credential type ID cannot be determined. Environment ID: %s, Credential Type ID: %s", environmentID, credentialTypeID)
	}

	_, err := apiClient.CredentialTypesApi.DeleteCredentialType(ctx, environmentID, credentialTypeID).Execute()
	if err != nil {
		t.Fatalf("Failed to delete Credential type: %v", err)
	}
}
