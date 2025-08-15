// Copyright © 2025 Ping Identity Corporation

package sso

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

func ApplicationResource_CheckDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := legacysdk.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_application_resource" {
			continue
		}

		shouldContinue, err := legacysdk.CheckParentEnvironmentDestroy(ctx, p1Client.API.ManagementAPIClient, rs.Primary.Attributes["environment_id"])
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		_, r, err := apiClient.ApplicationResourcesApi.ReadOneApplicationResource(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["resource_id"], rs.Primary.ID).Execute()

		shouldContinue, err = acctest.CheckForResourceDestroy(r, err)
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		return fmt.Errorf("PingOne Application Resource %s still exists", rs.Primary.ID)
	}

	return nil
}

func ApplicationResource_GetIDs(resourceName string, environmentID, customResourceID, resourceID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource Not found: %s", resourceName)
		}

		if resourceID != nil {
			*resourceID = rs.Primary.ID
		}

		if customResourceID != nil {
			*customResourceID = rs.Primary.Attributes["resource_id"]
		}

		if environmentID != nil {
			*environmentID = rs.Primary.Attributes["environment_id"]
		}

		return nil
	}
}

func ApplicationResource_RemovalDrift_PreConfig(ctx context.Context, apiClient *management.APIClient, t *testing.T, environmentID, customResourceID, resourceID string) {
	if environmentID == "" || customResourceID == "" || resourceID == "" {
		t.Fatalf("One of environment ID, custom resource ID or application resource ID cannot be determined. Environment ID: %s, Resource ID: %s, Application resource ID: %s", environmentID, customResourceID, resourceID)
	}

	_, err := apiClient.ApplicationResourcesApi.DeleteApplicationResource(ctx, environmentID, customResourceID, resourceID).Execute()
	if err != nil {
		t.Fatalf("Failed to delete application resource: %v", err)
	}
}
