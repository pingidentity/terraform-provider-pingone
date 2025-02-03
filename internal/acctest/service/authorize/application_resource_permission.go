// Copyright © 2025 Ping Identity Corporation

package authorize

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/patrickcping/pingone-go-sdk-v2/authorize"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
)

func ApplicationResourcePermission_CheckDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.AuthorizeAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_application_resource_permission" {
			continue
		}

		shouldContinue, err := acctest.CheckParentEnvironmentDestroy(ctx, p1Client.API.ManagementAPIClient, rs.Primary.Attributes["environment_id"])
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		_, r, err := apiClient.ApplicationResourcePermissionsApi.ReadOneApplicationPermission(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["application_resource_id"], rs.Primary.ID).Execute()

		shouldContinue, err = acctest.CheckForResourceDestroy(r, err)
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		return fmt.Errorf("PingOne Application Resource Permission Instance %s still exists", rs.Primary.ID)
	}

	return nil
}

func ApplicationResourcePermission_GetIDs(resourceName string, environmentID, applicationResourceID, resourceID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource Not found: %s", resourceName)
		}

		if resourceID != nil {
			*resourceID = rs.Primary.ID
		}

		if applicationResourceID != nil {
			*applicationResourceID = rs.Primary.Attributes["application_resource_id"]
		}

		if environmentID != nil {
			*environmentID = rs.Primary.Attributes["environment_id"]
		}

		return nil
	}
}

func ApplicationResourcePermission_RemovalDrift_PreConfig(ctx context.Context, apiClient *authorize.APIClient, t *testing.T, environmentID, applicationResourceID, permissionID string) {
	if environmentID == "" || applicationResourceID == "" || permissionID == "" {
		t.Fatalf("One of environment ID or decision endpoint ID cannot be determined. Environment ID: %s, Application Resource ID: %s, Permission ID: %s", environmentID, applicationResourceID, permissionID)
	}

	_, err := apiClient.ApplicationResourcePermissionsApi.DeleteApplicationPermission(ctx, environmentID, applicationResourceID, permissionID).Execute()
	if err != nil {
		t.Fatalf("Failed to delete application resource permission: %v", err)
	}
}
