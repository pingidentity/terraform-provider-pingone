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

func TestAccCheckDigitalWalletApplicationDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.CredentialsAPIClient

	mgmtApiClient := p1Client.API.ManagementAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_digital_wallet_application" {
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

		body, r, err := apiClient.DigitalWalletAppsApi.ReadOneDigitalWalletApp(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["id"]).Execute()

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

		return fmt.Errorf("PingOne Digital Wallet Application ID %s still exists", rs.Primary.ID)
	}

	return nil
}

func TestAccGetDigitalWalletApplicationIDs(resourceName string, environmentID, resourceID, applicationID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource Not found: %s", resourceName)
		}

		*resourceID = rs.Primary.ID
		*environmentID = rs.Primary.Attributes["environment_id"]
		*applicationID = rs.Primary.Attributes["application_id"]

		return nil
	}
}

func DigitalWalletApplication_RemovalDrift_PreConfig(ctx context.Context, apiClient *credentials.APIClient, t *testing.T, environmentID, digitalWalletAppID string) {
	if environmentID == "" || digitalWalletAppID == "" {
		t.Fatalf("One of environment ID or digital wallet application ID cannot be determined. Environment ID: %s, Digital wallet app ID: %s", environmentID, digitalWalletAppID)
	}

	_, err := apiClient.DigitalWalletAppsApi.DeleteDigitalWalletApp(ctx, environmentID, digitalWalletAppID).Execute()
	if err != nil {
		t.Fatalf("Failed to delete Digital wallet app: %v", err)
	}
}
