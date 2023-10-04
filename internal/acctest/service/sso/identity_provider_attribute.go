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

func TestAccCheckIdentityProviderAttributeDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_identity_provider_attribute" {
			continue
		}

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

		body, r, err := apiClient.IdentityProviderAttributesApi.ReadOneIdentityProviderAttribute(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["identity_provider_id"], rs.Primary.ID).Execute()

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

		return fmt.Errorf("PingOne Identity Provider attribute %s still exists", rs.Primary.ID)
	}

	return nil
}

func TestAccGetIdentityProviderAttributeIDs(resourceName string, environmentID, identityProviderID, resourceID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource Not found: %s", resourceName)
		}

		*resourceID = rs.Primary.ID
		*identityProviderID = rs.Primary.Attributes["identity_provider_id"]
		*environmentID = rs.Primary.Attributes["environment_id"]

		return nil
	}
}

func IdentityProviderAttribute_RemovalDrift_PreConfig(ctx context.Context, apiClient *management.APIClient, t *testing.T, environmentID, identityProviderID, identityProviderAttributeID string) {
	if environmentID == "" || identityProviderID == "" || identityProviderAttributeID == "" {
		t.Fatalf("One of environment ID, identity provider ID or identity provider attribute ID cannot be determined. Environment ID: %s, Identity provider ID: %s, Identity provider attribute ID: %s", environmentID, identityProviderID, identityProviderAttributeID)
	}

	_, err := apiClient.IdentityProviderAttributesApi.DeleteIdentityProviderAttribute(ctx, environmentID, identityProviderID, identityProviderAttributeID).Execute()
	if err != nil {
		t.Fatalf("Failed to delete identity provider attribute mapping: %v", err)
	}
}
