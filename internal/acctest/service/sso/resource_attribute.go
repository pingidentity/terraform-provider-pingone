// Copyright © 2026 Ping Identity Corporation

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

func ResourceAttribute_CheckDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := legacysdk.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_resource_attribute" {
			continue
		}

		shouldContinue, err := legacysdk.CheckParentEnvironmentDestroy(ctx, p1Client.API.ManagementAPIClient, rs.Primary.Attributes["environment_id"])
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		_, r, err := apiClient.ResourceAttributesApi.ReadOneResourceAttribute(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["resource_id"], rs.Primary.ID).Execute()

		shouldContinue, err = acctest.CheckForResourceDestroy(r, err)
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		return fmt.Errorf("PingOne Resource Mapping %s still exists", rs.Primary.ID)
	}

	return nil
}

func ResourceAttribute_GetIDs(resourceName string, environmentID, oidcResourceID, resourceID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		if resourceID != nil {
			*resourceID = rs.Primary.ID
		}

		if oidcResourceID != nil {
			*oidcResourceID = rs.Primary.Attributes["resource_id"]
		}

		if environmentID != nil {
			*environmentID = rs.Primary.Attributes["environment_id"]
		}

		return nil
	}
}

func ResourceAttribute_RemovalDrift_PreConfig(ctx context.Context, apiClient *management.APIClient, t *testing.T, environmentID, customResourceID, resourceID string) {
	if environmentID == "" || customResourceID == "" || resourceID == "" {
		t.Fatalf("One of environment ID, custom resource ID or resource attribute ID cannot be determined. Environment ID: %s, Resource ID: %s, Resource Attribute ID: %s", environmentID, customResourceID, resourceID)
	}

	_, err := apiClient.ResourceAttributesApi.DeleteResourceAttribute(ctx, environmentID, customResourceID, resourceID).Execute()
	if err != nil {
		t.Fatalf("Failed to delete resource attribute: %v", err)
	}
}
