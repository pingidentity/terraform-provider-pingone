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

func AlertChannel_CheckDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_alert_channel" {
			continue
		}

		shouldContinue, err := acctest.CheckParentEnvironmentDestroy(ctx, apiClient, rs.Primary.Attributes["environment_id"])
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		listResponse, r, err := apiClient.AlertingApi.ReadAllAlertChannels(ctx, rs.Primary.Attributes["environment_id"]).Execute()

		shouldContinue, err = acctest.CheckForResourceDestroy(r, err)
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		// Find the resource in the list
		var response *management.AlertChannel
		if embedded, ok := listResponse.GetEmbeddedOk(); ok {
			if alertChannels, ok := embedded.GetAlertChannelsOk(); ok {
				for _, alertChannel := range alertChannels {
					if alertChannel.GetId() == rs.Primary.ID {
						response = &alertChannel
						break
					}
				}
			}
		}

		if response == nil {
			continue
		}

		return fmt.Errorf("PingOne Alert Channel Instance %s still exists", rs.Primary.ID)
	}

	return nil
}

func AlertChannel_GetIDs(resourceName string, environmentID, resourceID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource Not found: %s", resourceName)
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

func AlertChannel_RemovalDrift_PreConfig(ctx context.Context, apiClient *management.APIClient, t *testing.T, environmentID, alertChannelID string) {
	if environmentID == "" || alertChannelID == "" {
		t.Fatalf("One of environment ID or alert channel ID cannot be determined. Environment ID: %s, Alert Channel ID: %s", environmentID, alertChannelID)
	}

	_, err := apiClient.AlertingApi.DeleteAlertChannel(ctx, environmentID, alertChannelID).Execute()
	if err != nil {
		t.Fatalf("Failed to delete Alert Channel: %v", err)
	}
}
