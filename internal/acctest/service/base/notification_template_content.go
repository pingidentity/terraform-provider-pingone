// Copyright © 2026 Ping Identity Corporation

package base

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/pingidentity/pingone-go-sdk-legacy/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/legacysdk"
)

// NotificationTemplateContent_GetDefaultContentID finds predefined content for the locale and delivery method.
func NotificationTemplateContent_GetDefaultContentID(ctx context.Context, apiClient *management.APIClient, t *testing.T, environmentID, templateName, locale string, deliveryMethod management.EnumTemplateContentDeliveryMethod) string {
	t.Helper()

	pagedIterator := apiClient.NotificationsTemplatesApi.ReadAllTemplateContents(ctx, environmentID, management.EnumTemplateName(templateName)).Execute()

	for pageCursor, err := range pagedIterator {
		if err != nil {
			t.Fatalf("Failed to read notification template contents: %v", err)
		}

		if contents, ok := pageCursor.EntityArray.Embedded.GetContentsOk(); ok {

			for _, content := range contents {

				common, ok := content.GetActualInstance().(interface {
					GetId() string
					GetDefault() bool
					GetLocale() string
					GetDeliveryMethod() management.EnumTemplateContentDeliveryMethod
				})

				if ok && common.GetDefault() && common.GetLocale() == locale && common.GetDeliveryMethod() == deliveryMethod {
					return common.GetId()
				}
			}
		}
	}

	t.Fatalf("Cannot find default %s content for template %s and locale %s in environment %s", deliveryMethod, templateName, locale, environmentID)

	return ""
}

func NotificationTemplateContent_CheckDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := legacysdk.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_notification_template_content" {
			continue
		}

		shouldContinue, err := legacysdk.CheckParentEnvironmentDestroy(ctx, apiClient, rs.Primary.Attributes["environment_id"])
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		_, r, err := apiClient.NotificationsTemplatesApi.ReadOneContent(ctx, rs.Primary.Attributes["environment_id"], management.EnumTemplateName(rs.Primary.Attributes["template_name"]), rs.Primary.ID).Execute()

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

func NotificationTemplateContent_GetIDs(resourceName string, environmentID, templateName, resourceID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		if resourceID != nil {
			*resourceID = rs.Primary.ID
		}

		if templateName != nil {
			*templateName = rs.Primary.Attributes["template_name"]
		}

		if environmentID != nil {
			*environmentID = rs.Primary.Attributes["environment_id"]
		}

		return nil
	}
}

func NotificationTemplateContent_RemovalDrift_PreConfig(ctx context.Context, apiClient *management.APIClient, t *testing.T, environmentID, templateName, resourceID string) {
	if environmentID == "" || templateName == "" || resourceID == "" {
		t.Fatalf("One of environment ID, template name or resource ID cannot be determined. Environment ID: %s, Template Name: %s, Resource ID: %s", environmentID, templateName, resourceID)
	}

	_, err := apiClient.NotificationsTemplatesApi.DeleteContent(ctx, environmentID, management.EnumTemplateName(templateName), resourceID).Execute()
	if err != nil {
		t.Fatalf("Failed to delete notification template contents: %v", err)
	}
}
