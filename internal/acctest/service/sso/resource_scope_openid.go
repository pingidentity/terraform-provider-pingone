// Copyright © 2025 Ping Identity Corporation

package sso

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
)

func ResourceScopeOpenID_CheckDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient

	re, err := regexp.Compile(`^(address|email|openid|phone|profile)$`)
	if err != nil {
		return fmt.Errorf("Cannot compile regex check for predefined scopes.")
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_resource_scope_openid" {
			continue
		}

		if m := re.MatchString(rs.Primary.Attributes["name"]); m {
			return nil
		} else {

			shouldContinue, err := acctest.CheckParentEnvironmentDestroy(ctx, p1Client.API.ManagementAPIClient, rs.Primary.Attributes["environment_id"])
			if err != nil {
				return err
			}

			if shouldContinue {
				continue
			}

			_, r, err := apiClient.ResourcesApi.ReadOneResource(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.ID).Execute()

			shouldContinue, err = acctest.CheckForResourceDestroy(r, err)
			if err != nil {
				return err
			}

			if shouldContinue {
				continue
			}

			return fmt.Errorf("PingOne Resource scope Instance %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

func ResourceScopeOpenID_GetIDs(resourceName string, environmentID, openidResourceID, resourceID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource Not found: %s", resourceName)
		}

		if resourceID != nil {
			*resourceID = rs.Primary.ID
		}

		if openidResourceID != nil {
			*openidResourceID = rs.Primary.Attributes["resource_id"]
		}

		if environmentID != nil {
			*environmentID = rs.Primary.Attributes["environment_id"]
		}

		return nil
	}
}

func ResourceScopeOpenID_RemovalDrift_PreConfig(ctx context.Context, apiClient *management.APIClient, t *testing.T, environmentID, openidResourceID, resourceScopeID string) {
	if environmentID == "" || openidResourceID == "" || resourceScopeID == "" {
		t.Fatalf("One of environment ID, OIDC resource ID or resource scope ID cannot be determined. Environment ID: %s, OpenID Resource ID: %s, Resource Scope ID: %s", environmentID, openidResourceID, resourceScopeID)
	}

	_, err := apiClient.ResourceScopesApi.DeleteResourceScope(ctx, environmentID, openidResourceID, resourceScopeID).Execute()
	if err != nil {
		t.Fatalf("Failed to delete OIDC resource scope: %v", err)
	}
}
