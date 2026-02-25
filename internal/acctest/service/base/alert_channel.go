// Copyright © 2026 Ping Identity Corporation

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

func AlertChannel_CheckDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := legacysdk.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient

mainloop:
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_alert_channel" {
			continue
		}

		shouldContinue, err := legacysdk.CheckParentEnvironmentDestroy(ctx, apiClient, rs.Primary.Attributes["environment_id"])
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		pagedIterator := apiClient.AlertingApi.ReadAllAlertChannels(ctx, rs.Primary.Attributes["environment_id"]).Execute()

		found := false

	pagedIteratorLoop:
		for pageCursor, err := range pagedIterator {
			shouldContinue, err = acctest.CheckForResourceDestroy(pageCursor.HTTPResponse, err)
			if err != nil {
				return err
			}

			// Environment not found
			if shouldContinue {
				continue mainloop
			}

			// Find the resource in the list
			if embedded, ok := pageCursor.EntityArray.GetEmbeddedOk(); ok {
				if alertChannels, ok := embedded.GetAlertChannelsOk(); ok {
					for _, alertChannel := range alertChannels {
						if alertChannel.GetId() == rs.Primary.ID {
							found = true
							break pagedIteratorLoop
						}
					}
				}
			}
		}

		if !found {
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
			return fmt.Errorf("resource not found: %s", resourceName)
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
