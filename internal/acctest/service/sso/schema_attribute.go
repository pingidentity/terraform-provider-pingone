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

func TestAccCheckSchemaAttributeDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_schema_attribute" {
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

		body, r, err := apiClient.SchemasApi.ReadOneAttribute(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["schema_id"], rs.Primary.ID).Execute()

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

		return fmt.Errorf("PingOne Schema Attribute Instance %s still exists", rs.Primary.ID)
	}

	return nil
}

func TestAccGetSchemaAttributeIDs(resourceName string, environmentID, schemaID, resourceID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource Not found: %s", resourceName)
		}

		*resourceID = rs.Primary.ID
		*schemaID = rs.Primary.Attributes["schema_id"]
		*environmentID = rs.Primary.Attributes["environment_id"]

		return nil
	}
}

func SchemaAttribute_RemovalDrift_PreConfig(ctx context.Context, apiClient *management.APIClient, t *testing.T, environmentID, schemaID, schemaAttributeID string) {
	if environmentID == "" || schemaID == "" || schemaAttributeID == "" {
		t.Fatalf("One of environment ID, schema ID or schema attribute ID cannot be determined. Environment ID: %s, Schema ID: %s, Schema Attribute ID: %s", environmentID, schemaID, schemaAttributeID)
	}

	_, err := apiClient.SchemasApi.DeleteAttribute(ctx, environmentID, schemaID, schemaAttributeID).Execute()
	if err != nil {
		t.Fatalf("Failed to delete schema attribute: %v", err)
	}
}
