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
)

func FormsRecaptchaV2_CheckDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_forms_recaptcha_v2" {
			continue
		}

		shouldContinue, err := acctest.CheckParentEnvironmentDestroy(ctx, apiClient, rs.Primary.Attributes["environment_id"])
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		_, r, err := apiClient.RecaptchaConfigurationApi.ReadRecaptchaConfiguration(ctx, rs.Primary.Attributes["environment_id"]).Execute()

		destroyHttpCode := 204
		shouldContinue, err = acctest.CheckForResourceDestroyCustomHTTPCode(r, err, destroyHttpCode)
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		return fmt.Errorf("PingOne Forms Recaptcha v2 config Instance %s still exists", rs.Primary.Attributes["environment_id"])
	}

	return nil
}

func FormsRecaptchaV2_GetIDs(resourceName string, environmentID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource Not found: %s", resourceName)
		}

		if environmentID != nil {
			*environmentID = rs.Primary.Attributes["environment_id"]
		}

		return nil
	}
}

func FormsRecaptchaV2_RemovalDrift_PreConfig(ctx context.Context, apiClient *management.APIClient, t *testing.T, environmentID string) {
	if environmentID == "" {
		t.Fatalf("One of environment ID or form ID cannot be determined. Environment ID: %s", environmentID)
	}

	_, err := apiClient.RecaptchaConfigurationApi.DeleteRecaptchaConfiguration(ctx, environmentID).Execute()
	if err != nil {
		t.Fatalf("Failed to delete FormsRecaptchaV2: %v", err)
	}
}
