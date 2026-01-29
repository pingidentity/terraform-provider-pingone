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

func AgreementLocalization_CheckDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := legacysdk.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_agreement_localization" {
			continue
		}

		shouldContinue, err := legacysdk.CheckParentEnvironmentDestroy(ctx, apiClient, rs.Primary.Attributes["environment_id"])
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		_, r, err := apiClient.AgreementLanguagesResourcesApi.ReadOneAgreementLanguage(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["agreement_id"], rs.Primary.ID).Execute()

		shouldContinue, err = acctest.CheckForResourceDestroy(r, err)
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		return fmt.Errorf("PingOne agreement localization %s still exists", rs.Primary.ID)
	}

	return nil
}

func AgreementLocalization_GetIDs(resourceName string, environmentID, agreementID, resourceID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		if resourceID != nil {
			*resourceID = rs.Primary.ID
		}

		if agreementID != nil {
			*agreementID = rs.Primary.Attributes["agreement_id"]
		}

		if environmentID != nil {
			*environmentID = rs.Primary.Attributes["environment_id"]
		}

		return nil
	}
}

func AgreementLocalization_RemovalDrift_PreConfig(ctx context.Context, apiClient *management.APIClient, t *testing.T, environmentID, agreementID, agreementLocalizationID string) {
	if environmentID == "" || agreementID == "" || agreementLocalizationID == "" {
		t.Fatalf("One of environment ID, agreement ID or agreement localization ID cannot be determined. Environment ID: %s, Agreement ID: %s, Agreement Localization ID: %s", environmentID, agreementID, agreementLocalizationID)
	}

	agreementLanguage, r, err := apiClient.AgreementLanguagesResourcesApi.ReadOneAgreementLanguage(ctx, environmentID, agreementID, agreementLocalizationID).Execute()
	if err != nil || r.StatusCode != 200 {
		t.Fatalf("Failed to read agreement language to delete: status code: %d, error: %v", r.StatusCode, err)
	}

	if agreementLanguage == nil {
		t.Fatalf("Agreement language not found.  Cannot disable before delete")
	}

	agreementLanguage.SetEnabled(false)

	_, r, err = apiClient.AgreementLanguagesResourcesApi.UpdateAgreementLanguage(ctx, environmentID, agreementID, agreementLocalizationID).AgreementLanguage(*agreementLanguage).Execute()
	if err != nil || r.StatusCode != 200 {
		t.Fatalf("Failed to disable agreement language prior to delete: status code: %d, error: %v", r.StatusCode, err)
	}

	_, err = apiClient.AgreementLanguagesResourcesApi.DeleteAgreementLanguage(ctx, environmentID, agreementID, agreementLocalizationID).Execute()
	if err != nil {
		t.Fatalf("Failed to delete agreement localisation: %v", err)
	}
}
