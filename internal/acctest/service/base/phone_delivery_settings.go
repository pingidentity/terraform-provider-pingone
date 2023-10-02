package base

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
)

func TestAccCheckPhoneDeliverySettingsDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_phone_delivery_settings" {
			continue
		}

		shouldContinue, err := acctest.CheckParentEnvironmentDestroy(ctx, apiClient, rs.Primary.Attributes["environment_id"])
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		_, r, err := apiClient.PhoneDeliverySettingsApi.ReadOnePhoneDeliverySettings(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.ID).Execute()

		shouldContinue, err = acctest.CheckForResourceDestroy(r, err)
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		return fmt.Errorf("PingOne Phone Delivery Settings %s still exists", rs.Primary.ID)
	}

	return nil
}

func TestAccGetPhoneDeliverySettingsIDs(resourceName string, environmentID, resourceID *string) resource.TestCheckFunc {
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

func PhoneDeliverySettings_RemovalDrift_PreConfig(ctx context.Context, apiClient *management.APIClient, t *testing.T, environmentID, phoneDeliverySettingsID string) {
	if environmentID == "" || phoneDeliverySettingsID == "" {
		t.Fatalf("One of environment ID or phone delivery settings ID cannot be determined. Environment ID: %s, Phone Delivery Settings ID: %s", environmentID, phoneDeliverySettingsID)
	}

	_, err := apiClient.PhoneDeliverySettingsApi.DeletePhoneDeliverySettings(ctx, environmentID, phoneDeliverySettingsID).Execute()
	if err != nil {
		t.Fatalf("Failed to delete phone delivery settings: %v", err)
	}
}
