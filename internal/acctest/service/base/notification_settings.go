// Copyright © 2026 Ping Identity Corporation

package base

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
)

func NotificationSettings_CheckDestroy(s *terraform.State) error {
	return nil
}

func NotificationSettings_GetIDs(resourceName string, resourceID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		if resourceID != nil {
			*resourceID = rs.Primary.ID
		}

		return nil
	}
}

func NotificationSettings_RemovalDrift_PreConfig(ctx context.Context, apiClient *management.APIClient, t *testing.T, notificationSettingsID string) {
	if notificationSettingsID == "" {
		t.Fatalf("Notification Settings ID cannot be determined. Notification Settings ID: %s", notificationSettingsID)
	}

	_, _, err := apiClient.NotificationsSettingsApi.DeleteNotificationsSettings(ctx, notificationSettingsID).Execute()
	if err != nil {
		t.Fatalf("Failed to delete notification settings: %v", err)
	}
}
