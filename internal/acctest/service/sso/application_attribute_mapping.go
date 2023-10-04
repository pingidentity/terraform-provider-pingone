package sso

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
)

func ApplicationAttributeMapping_CheckDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_application_attribute_mapping" {
			continue
		}

		shouldContinue, err := acctest.CheckParentEnvironmentDestroy(ctx, p1Client.API.ManagementAPIClient, rs.Primary.Attributes["environment_id"])
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		_, r, err := apiClient.ApplicationAttributeMappingApi.ReadOneApplicationAttributeMapping(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["application_id"], rs.Primary.ID).Execute()

		shouldContinue, err = acctest.CheckForResourceDestroy(r, err)
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		return fmt.Errorf("PingOne Application Attribute Mapping %s still exists", rs.Primary.ID)
	}

	return nil
}

func ApplicationAttributeMapping_GetIDs(resourceName string, environmentID, applicationID, resourceID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource Not found: %s", resourceName)
		}

		*resourceID = rs.Primary.ID
		*applicationID = rs.Primary.Attributes["application_id"]
		*environmentID = rs.Primary.Attributes["environment_id"]

		return nil
	}
}

func ApplicationAttributeMapping_RemovalDrift_PreConfig(ctx context.Context, apiClient *management.APIClient, t *testing.T, environmentID, applicationID, applicationAttributeMappingID string) {
	if environmentID == "" || applicationID == "" || applicationAttributeMappingID == "" {
		t.Fatalf("One of environment ID, application ID or application attribute mapping ID cannot be determined. Environment ID: %s, Application ID: %s, Application Attribute Mapping ID: %s", environmentID, applicationID, applicationAttributeMappingID)
	}

	_, err := apiClient.ApplicationAttributeMappingApi.DeleteApplicationAttributeMapping(ctx, environmentID, applicationID, applicationAttributeMappingID).Execute()
	if err != nil {
		t.Fatalf("Failed to delete Application attribute mapping: %v", err)
	}
}
