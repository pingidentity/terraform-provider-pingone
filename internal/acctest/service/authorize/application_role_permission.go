// Copyright © 2026 Ping Identity Corporation

package authorize

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/patrickcping/pingone-go-sdk-v2/authorize"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/legacysdk"
)

func ApplicationRolePermission_CheckDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := legacysdk.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.AuthorizeAPIClient

mainloop:
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_authorize_application_role_permission" {
			continue
		}

		shouldContinue, err := legacysdk.CheckParentEnvironmentDestroy(ctx, p1Client.API.ManagementAPIClient, rs.Primary.Attributes["environment_id"])
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		pagedIterator := apiClient.ApplicationRolePermissionsApi.ReadApplicationRolePermissions(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["application_role_id"]).Execute()

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

			if pageCursor.EntityArray.Embedded != nil && pageCursor.EntityArray.Embedded.Permissions != nil {
				for _, permission := range pageCursor.EntityArray.Embedded.Permissions {
					if v := permission.ApplicationRolePermission; v != nil && v.GetId() == rs.Primary.ID {
						found = true
						break pagedIteratorLoop
					}
				}
			}

		}

		if !found {
			continue
		}

		return fmt.Errorf("PingOne Application Role Permission Instance %s still exists", rs.Primary.ID)
	}

	return nil
}

func ApplicationRolePermission_GetIDs(resourceName string, environmentID, applicationRoleID, resourceID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		if applicationRoleID != nil {
			*applicationRoleID = rs.Primary.Attributes["application_role_id"]
		}

		if environmentID != nil {
			*environmentID = rs.Primary.Attributes["environment_id"]
		}

		if resourceID != nil {
			*resourceID = rs.Primary.Attributes["application_resource_permission_id"]
		}

		return nil
	}
}

func ApplicationRolePermission_RemovalDrift_PreConfig(ctx context.Context, apiClient *authorize.APIClient, t *testing.T, environmentID, applicationRoleID, applicationRolePermissionID string) {
	if environmentID == "" || applicationRoleID == "" || applicationRolePermissionID == "" {
		t.Fatalf("One of environment ID, application role ID or application role permission ID cannot be determined. Environment ID: %s, Application Role ID: %s, Application Role Permission ID: %s", environmentID, applicationRoleID, applicationRolePermissionID)
	}

	_, err := apiClient.ApplicationRolePermissionsApi.DeleteApplicationRolePermission(ctx, environmentID, applicationRoleID, applicationRolePermissionID).Execute()
	if err != nil {
		t.Fatalf("Failed to delete application role permission: %v", err)
	}
}
