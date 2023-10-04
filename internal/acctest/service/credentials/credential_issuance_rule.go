package credentials

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/patrickcping/pingone-go-sdk-v2/credentials"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
)

func TestAccCheckCredentialIssuanceRuleDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.CredentialsAPIClient

	mgmtApiClient := p1Client.API.ManagementAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_credential_issuance_rule" {
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

		body, r, err := apiClient.CredentialIssuanceRulesApi.ReadOneCredentialIssuanceRule(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["credential_type_id"], rs.Primary.Attributes["id"]).Execute()

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

		return fmt.Errorf("PingOne Credential Issuance Rule ID %s still exists", rs.Primary.ID)
	}

	return nil
}

func TestAccGetCredentialIssuanceRuleIDs(resourceName string, environmentID, credentialTypeID, digitalWalletApplicationID, resourceID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource Not found: %s", resourceName)
		}

		*resourceID = rs.Primary.ID
		*credentialTypeID = rs.Primary.Attributes["credential_type_id"]
		*digitalWalletApplicationID = rs.Primary.Attributes["digital_wallet_application_id"]
		*environmentID = rs.Primary.Attributes["environment_id"]

		return nil
	}
}

func CredentialIssuanceRule_RemovalDrift_PreConfig(ctx context.Context, apiClient *credentials.APIClient, t *testing.T, environmentID, credentialTypeID, credentialIssuanceRuleID string) {
	if environmentID == "" || credentialIssuanceRuleID == "" {
		t.Fatalf("One of environment ID, credential type ID or credential issuance rule ID cannot be determined. Environment ID: %s, Credential Type ID: %s, Credential Issuance rule ID: %s", environmentID, credentialTypeID, credentialIssuanceRuleID)
	}

	_, err := apiClient.CredentialIssuanceRulesApi.DeleteCredentialIssuanceRule(ctx, environmentID, credentialTypeID, credentialIssuanceRuleID).Execute()
	if err != nil {
		t.Fatalf("Failed to delete Credential issuance rule: %v", err)
	}
}
