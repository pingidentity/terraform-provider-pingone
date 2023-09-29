package base

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
)

func TestAccCheckNotificationSettingsDestroy(s *terraform.State) error {
	return nil
}

func TestAccGetNotificationSettingsIDs(resourceName string, resourceID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource Not found: %s", resourceName)
		}

		*resourceID = rs.Primary.ID

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
