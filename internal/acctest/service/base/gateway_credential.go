package base

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

func TestAccCheckGatewayCredentialDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_gateway_credential" {
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

		body, r, err := apiClient.GatewayCredentialsApi.ReadOneGatewayCredential(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["gateway_id"], rs.Primary.ID).Execute()

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

		return fmt.Errorf("PingOne Gateway Credential Instance %s still exists", rs.Primary.ID)
	}

	return nil
}

func TestAccGetGatewayCredentialIDs(resourceName string, environmentID, gatewayID, resourceID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource Not found: %s", resourceName)
		}

		*resourceID = rs.Primary.ID
		*gatewayID = rs.Primary.Attributes["gateway_id"]
		*environmentID = rs.Primary.Attributes["environment_id"]

		return nil
	}
}

func GatewayCredential_RemovalDrift_PreConfig(ctx context.Context, apiClient *management.APIClient, t *testing.T, environmentID, gatewayID, gatewayCredentialID string) {
	if environmentID == "" || gatewayID == "" || gatewayCredentialID == "" {
		t.Fatalf("One of environment ID, gateway ID or gateway credential ID cannot be determined. Environment ID: %s, Gateway ID: %s, Gateway Credential ID: %s", environmentID, gatewayID, gatewayCredentialID)
	}

	_, err := apiClient.GatewayCredentialsApi.DeleteGatewayCredential(ctx, environmentID, gatewayID, gatewayCredentialID).Execute()
	if err != nil {
		t.Fatalf("Failed to delete gateway credential: %v", err)
	}
}
