package sso

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
)

func TestAccCheckGroupNestingDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_group_nesting" {
			continue
		}

		body, r, err := apiClient.GroupsApi.ReadOneGroupNesting(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["group_id"], rs.Primary.ID).Execute()

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

		return fmt.Errorf("PingOne Group Nesting Instance %s still exists", rs.Primary.ID)
	}

	return nil
}

func TestAccGetGroupNestingIDs(resourceName string, environmentID, groupID, resourceID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource Not found: %s", resourceName)
		}

		*resourceID = rs.Primary.ID
		*groupID = rs.Primary.Attributes["group_id"]
		*environmentID = rs.Primary.Attributes["environment_id"]

		return nil
	}
}

func GroupNesting_RemovalDrift_PreConfig(ctx context.Context, apiClient *management.APIClient, t *testing.T, environmentID, groupID, groupNestingID string) {
	if environmentID == "" || groupID == "" || groupNestingID == "" {
		t.Fatalf("One of environment ID, group ID or group nesting ID cannot be determined. Environment ID: %s, Group ID: %s, Group Nesting ID: %s", environmentID, groupID, groupNestingID)
	}

	_, err := apiClient.GroupsApi.DeleteGroupNesting(ctx, environmentID, groupID, groupNestingID).Execute()
	if err != nil {
		t.Fatalf("Failed to delete group nesting: %v", err)
	}
}
