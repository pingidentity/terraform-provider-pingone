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

func TestAccCheckNotificationTemplateContentDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_notification_template_content" {
			continue
		}

		shouldContinue, err := acctest.CheckParentEnvironmentDestroy(ctx, apiClient, rs.Primary.Attributes["environment_id"])
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		_, r, err := apiClient.NotificationsTemplatesApi.ReadOneContent(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["template_name"], rs.Primary.ID).Execute()

		shouldContinue, err = acctest.CheckForResourceDestroy(r, err)
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		return fmt.Errorf("PingOne User Role Assignment %s still exists", rs.Primary.ID)
	}

	return nil
}

func TestAccGetNotificationTemplateContentIDs(resourceName string, environmentID, templateName, resourceID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource Not found: %s", resourceName)
		}

		*resourceID = rs.Primary.ID
		*templateName = rs.Primary.Attributes["template_name"]
		*environmentID = rs.Primary.Attributes["environment_id"]

		return nil
	}
}

func NotificationTemplateContent_RemovalDrift_PreConfig(ctx context.Context, apiClient *management.APIClient, t *testing.T, environmentID, templateName, resourceID string) {
	if environmentID == "" || templateName == "" || resourceID == "" {
		t.Fatalf("One of environment ID, template name or resource ID cannot be determined. Environment ID: %s, Template Name: %s, Resource ID: %s", environmentID, templateName, resourceID)
	}

	_, err := apiClient.NotificationsTemplatesApi.DeleteContent(ctx, environmentID, templateName, resourceID).Execute()
	if err != nil {
		t.Fatalf("Failed to delete notification template contents: %v", err)
	}
}
