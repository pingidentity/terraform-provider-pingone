package sso

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/patrickcping/pingone-go-sdk-v2/mfa"
)

func ResourceSecret_CheckDestroy(s *terraform.State) error {
	return nil
}

func ResourceSecret_GetIDs(resourceName string, environmentID, resourceID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource Not found: %s", resourceName)
		}

		*resourceID = rs.Primary.Attributes["resource_id"]
		*environmentID = rs.Primary.Attributes["environment_id"]

		return nil
	}
}

func ResourceSecret_RemovalDrift_PreConfig(ctx context.Context, apiClient *mfa.APIClient, t *testing.T, environmentID, resourceID string) {
}
