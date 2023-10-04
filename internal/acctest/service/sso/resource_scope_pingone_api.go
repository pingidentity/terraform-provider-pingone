package sso

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
)

func ResourceScopePingOneAPI_CheckDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient

	re, err := regexp.Compile(`^p1:(read|update):user$`)
	if err != nil {
		return fmt.Errorf("Cannot compile regex check for predefined scopes.")
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_resource_scope_pingone_api" {
			continue
		}

		if m := re.MatchString(rs.Primary.Attributes["name"]); m {
			return nil
		} else {

			_, rEnv, err := apiClient.EnvironmentsApi.ReadOneEnvironment(ctx, rs.Primary.Attributes["environment_id"]).Execute()

			if err != nil {

				if rEnv == nil {
					return fmt.Errorf("Response object does not exist and no error detected")
				}

				if rEnv.StatusCode == 404 {
					continue
				}

				return err
			}

			body, r, err := apiClient.ResourcesApi.ReadOneResource(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.ID).Execute()

			if err != nil {

				if r == nil {
					return fmt.Errorf("Response object does not exist and no error detected")
				}

				if r.StatusCode == 404 {
					continue
				}

				tflog.Error(ctx, fmt.Sprintf("Error: %v", body))
				return err
			}

			return fmt.Errorf("PingOne Resource scope Instance %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

func ResourceScopePingOneAPI_GetIDs(resourceName string, environmentID, openidResourceID, resourceID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource Not found: %s", resourceName)
		}

		*resourceID = rs.Primary.ID
		*openidResourceID = rs.Primary.Attributes["resource_id"]
		*environmentID = rs.Primary.Attributes["environment_id"]

		return nil
	}
}

func ResourceScopePingOneAPI_RemovalDrift_PreConfig(ctx context.Context, apiClient *management.APIClient, t *testing.T, environmentID, openidResourceID, resourceScopeID string) {
	if environmentID == "" || openidResourceID == "" || resourceScopeID == "" {
		t.Fatalf("One of environment ID, OpenID resource ID or resource scope ID cannot be determined. Environment ID: %s, OpenID resource ID: %s, Resource Scope ID: %s", environmentID, openidResourceID, resourceScopeID)
	}

	_, err := apiClient.ResourceScopesApi.DeleteResourceScope(ctx, environmentID, openidResourceID, resourceScopeID).Execute()
	if err != nil {
		t.Fatalf("Failed to delete PingOne API resource scope: %v", err)
	}
}
