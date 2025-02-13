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

func TrustFrameworkAttribute_CheckDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.AuthorizeAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_authorize_trust_framework_attribute" {
			continue
		}

		shouldContinue, err := acctest.CheckParentEnvironmentDestroy(ctx, p1Client.API.ManagementAPIClient, rs.Primary.Attributes["environment_id"])
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		_, r, err := apiClient.AuthorizeEditorAttributesApi.GetAttribute(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.ID).Execute()

		shouldContinue, err = acctest.CheckForResourceDestroy(r, err)
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		return fmt.Errorf("PingOne Authorize editor attribute Instance %s still exists", rs.Primary.ID)
	}

	return nil
}

func TrustFrameworkAttribute_GetIDs(resourceName string, environmentID, resourceID *string) resource.TestCheckFunc {
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

		return nil
	}
}

func TrustFrameworkAttribute_RemovalDrift_PreConfig(ctx context.Context, apiClient *authorize.APIClient, t *testing.T, environmentID, trustFrameworkAttributeID string) {
	if environmentID == "" || trustFrameworkAttributeID == "" {
		t.Fatalf("One of environment ID or authorize editor attribute ID cannot be determined. Environment ID: %s, Authorize Editor Attribute ID: %s", environmentID, trustFrameworkAttributeID)
	}

	_, err := apiClient.AuthorizeEditorAttributesApi.DeleteAttribute(ctx, environmentID, trustFrameworkAttributeID).Execute()
	if err != nil {
		t.Fatalf("Failed to delete editor attribute: %v", err)
	}
}
