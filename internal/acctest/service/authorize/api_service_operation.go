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
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/legacysdk"
)

func APIServiceOperation_CheckDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := legacysdk.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.AuthorizeAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_authorize_api_service_operation" {
			continue
		}

		shouldContinue, err := legacysdk.CheckParentEnvironmentDestroy(ctx, p1Client.API.ManagementAPIClient, rs.Primary.Attributes["environment_id"])
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		_, r, err := apiClient.APIServerOperationsApi.ReadOneAPIServerOperation(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["api_service_id"], rs.Primary.ID).Execute()

		shouldContinue, err = acctest.CheckForResourceDestroy(r, err)
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		return fmt.Errorf("PingOne API Service Operation Instance %s still exists", rs.Primary.ID)
	}

	return nil
}

func APIServiceOperation_GetIDs(resourceName string, environmentID, apiServiceID, resourceID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource Not found: %s", resourceName)
		}

		if resourceID != nil {
			*resourceID = rs.Primary.ID
		}

		if environmentID != nil {
			*environmentID = rs.Primary.Attributes["environment_id"]
		}

		if apiServiceID != nil {
			*apiServiceID = rs.Primary.Attributes["api_service_id"]
		}

		return nil
	}
}

func APIServiceOperation_RemovalDrift_PreConfig(ctx context.Context, apiClient *authorize.APIClient, t *testing.T, environmentID, apiServiceID, apiServiceOperationID string) {
	if environmentID == "" || apiServiceID == "" || apiServiceOperationID == "" {
		t.Fatalf("One of environment ID, API service ID or API service operation ID cannot be determined. Environment ID: %s, API Service ID: %s, API Service Operation ID: %s", environmentID, apiServiceID, apiServiceOperationID)
	}

	_, err := apiClient.APIServerOperationsApi.DeleteAPIServerOperation(ctx, environmentID, apiServiceID, apiServiceOperationID).Execute()
	if err != nil {
		t.Fatalf("Failed to delete API service operation: %v", err)
	}
}
