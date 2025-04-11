// Copyright © 2025 Ping Identity Corporation

package credentials

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/patrickcping/pingone-go-sdk-v2/credentials"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
)

func CredentialIssuanceRule_CheckDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.CredentialsAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_credential_issuance_rule" {
			continue
		}

		shouldContinue, err := acctest.CheckParentEnvironmentDestroy(ctx, p1Client.API.ManagementAPIClient, rs.Primary.Attributes["environment_id"])
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		_, r, err := apiClient.CredentialIssuanceRulesApi.ReadOneCredentialIssuanceRule(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["credential_type_id"], rs.Primary.Attributes["id"]).Execute()

		shouldContinue, err = acctest.CheckForResourceDestroy(r, err)
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		return fmt.Errorf("PingOne Credential Issuance Rule ID %s still exists", rs.Primary.ID)
	}

	return nil
}

func CredentialIssuanceRule_GetIDs(resourceName string, environmentID, credentialTypeID, applicationID, digitalWalletApplicationID, resourceID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource Not found: %s", resourceName)
		}

		if resourceID != nil {
			*resourceID = rs.Primary.ID
		}

		if credentialTypeID != nil {
			*credentialTypeID = rs.Primary.Attributes["credential_type_id"]
		}

		if applicationID != nil {
			app, ok := s.RootModule().Resources[strings.Replace(resourceName, "_credential_issuance_rule.", "_application.", 1)]
			if !ok {
				return fmt.Errorf("Resource Not found: %s", resourceName)
			}
			*applicationID = app.Primary.ID
		}

		if digitalWalletApplicationID != nil {
			*digitalWalletApplicationID = rs.Primary.Attributes["digital_wallet_application_id"]
		}

		if environmentID != nil {
			*environmentID = rs.Primary.Attributes["environment_id"]
		}

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
