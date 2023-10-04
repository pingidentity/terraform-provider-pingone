package credentials

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/patrickcping/pingone-go-sdk-v2/credentials"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
)

func TestAccCheckCredentialTypeDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.CredentialsAPIClient

	mgmtApiClient := p1Client.API.ManagementAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_credential_type" {
			continue
		}

		_, rEnv, err := mgmtApiClient.EnvironmentsApi.ReadOneEnvironment(ctx, rs.Primary.Attributes["environment_id"]).Execute()

		if err != nil {

			if rEnv == nil {
				return fmt.Errorf("Response object does not exist and no error detected")
			}

			if rEnv.StatusCode == 404 {
				continue
			}

			return err
		}

		body, r, err := apiClient.CredentialTypesApi.ReadOneCredentialType(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["id"]).Execute()

		if err != nil {

			if r == nil {
				return fmt.Errorf("Response object does not exist and no error detected")
			}

			if r.StatusCode == 404 {
				continue
			}

			tflog.Error(ctx, fmt.Sprintf("Error: %v", body))
			return err
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

func TestAccGetCredentialTypeIDs(resourceName string, environmentID, resourceID *string) resource.TestCheckFunc {
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

func CredentialType_RemovalDrift_PreConfig(ctx context.Context, apiClient *credentials.APIClient, t *testing.T, environmentID, credentialTypeID string) {
	if environmentID == "" || credentialTypeID == "" {
		t.Fatalf("One of environment ID or credential type ID cannot be determined. Environment ID: %s, Credential Type ID: %s", environmentID, credentialTypeID)
	}

	_, err := apiClient.CredentialTypesApi.DeleteCredentialType(ctx, environmentID, credentialTypeID).Execute()
	if err != nil {
		t.Fatalf("Failed to delete Credential type: %v", err)
	}
}
