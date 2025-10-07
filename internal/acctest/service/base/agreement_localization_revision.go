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

func AgreementLocalizationRevision_CheckDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := legacysdk.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_agreement_localization_revision" {
			continue
		}

		shouldContinue, err := legacysdk.CheckParentEnvironmentDestroy(ctx, apiClient, rs.Primary.Attributes["environment_id"])
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		_, r, err := apiClient.AgreementRevisionsResourcesApi.ReadOneAgreementLanguageRevision(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["agreement_id"], rs.Primary.Attributes["agreement_localization_id"], rs.Primary.ID).Execute()

		shouldContinue, err = acctest.CheckForResourceDestroy(r, err)
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		return fmt.Errorf("PingOne agreement localization revision %s still exists", rs.Primary.ID)
	}

	return nil
}

func AgreementLocalizationRevision_GetIDs(resourceName string, environmentID, agreementID, agreementLocalizationID, resourceID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource Not found: %s", resourceName)
		}

		if resourceID != nil {
			*resourceID = rs.Primary.ID
		}

		if agreementLocalizationID != nil {
			*agreementLocalizationID = rs.Primary.Attributes["agreement_localization_id"]
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

func AgreementLocalizationRevision_RemovalDrift_PreConfig(ctx context.Context, apiClient *management.APIClient, t *testing.T, environmentID, agreementID, agreementLocalizationID, agreementLocalizationRevisionID string) {
	if environmentID == "" || agreementID == "" || agreementLocalizationID == "" || agreementLocalizationRevisionID == "" {
		t.Fatalf("One of environment ID, agreement ID, agreement localization ID or agreement localization revision ID cannot be determined. Environment ID: %s, Agreement ID: %s, Agreement Localization ID: %s, Agreement Localization Revision ID: %s", environmentID, agreementID, agreementLocalizationID, agreementLocalizationRevisionID)
	}

	_, err := apiClient.AgreementRevisionsResourcesApi.DeleteAgreementLanguageRevision(ctx, environmentID, agreementID, agreementLocalizationID, agreementLocalizationRevisionID).Execute()
	if err != nil {
		t.Fatalf("Failed to delete agreement localisation revision: %v", err)
	}
}
