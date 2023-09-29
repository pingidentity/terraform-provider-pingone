package base

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
)

func TestAccCheckAgreementLocalizationEnableDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_agreement_localization_enable" {
			continue
		}

		_, rEnv, err := apiClient.EnvironmentsApi.ReadOneEnvironment(ctx, rs.Primary.Attributes["environment_id"]).Execute()

		if err != nil {

			if rEnv == nil {
				return fmt.Errorf("Response object does not exist and no error detected")
			}

			if rEnv.StatusCode == 404 {
				continue
			}

			return err
		}

		body, r, err := apiClient.AgreementLanguagesResourcesApi.ReadOneAgreementLanguage(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["agreement_id"], rs.Primary.ID).Execute()

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

		if !body.GetEnabled() {
			continue
		}

		return fmt.Errorf("PingOne agreement localization %s still exists", rs.Primary.ID)
	}

	return nil
}

func TestAccGetAgreementLocalizationEnableIDs(resourceName string, environmentID, agreementID, resourceID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource Not found: %s", resourceName)
		}

		*resourceID = rs.Primary.ID
		*agreementID = rs.Primary.Attributes["agreement_id"]
		*environmentID = rs.Primary.Attributes["environment_id"]

		return nil
	}
}

func AgreementLocalizationEnable_RemovalDrift_PreConfig(ctx context.Context, apiClient *management.APIClient, t *testing.T, environmentID, agreementID, agreementLocalizationID string) {
	if environmentID == "" || agreementID == "" || agreementLocalizationID == "" {
		t.Fatalf("One of environment ID, agreement ID or agreement localization ID cannot be determined. Environment ID: %s, Agreement ID: %s, Agreement Localization ID: %s", environmentID, agreementID, agreementLocalizationID)
	}

	_, err := apiClient.AgreementLanguagesResourcesApi.DeleteAgreementLanguage(ctx, environmentID, agreementID, agreementLocalizationID).Execute()
	if err != nil {
		t.Fatalf("Failed to delete agreement localisation: %v", err)
	}
}
