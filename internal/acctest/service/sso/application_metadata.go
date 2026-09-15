// Copyright © 2026 Ping Identity Corporation

package sso

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
)

func ApplicationMetadata_CheckDestroy(s *terraform.State) error {
	return nil
}

func ApplicationMetadata_GetIDs(resourceName string, environmentID, applicationID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		if applicationID != nil {
			*applicationID = rs.Primary.Attributes["application_id"]
		}

		if environmentID != nil {
			*environmentID = rs.Primary.Attributes["environment_id"]
		}

		return nil
	}
}

func ApplicationMetadata_RemovalDrift_PreConfig(ctx context.Context, apiClient *management.APIClient, t *testing.T, environmentID, applicationID string) {
	if environmentID == "" || applicationID == "" {
		t.Fatalf("One of environment ID or application ID cannot be determined. Environment ID: %s, Application ID: %s", environmentID, applicationID)
	}

	// A PUT with no `metadata` property removes the existing metadata from the application.
	_, _, err := apiClient.ApplicationMetadataApi.UpdateApplicationMetadata(ctx, environmentID, applicationID).ApplicationMetadata(*management.NewApplicationMetadata()).Execute()
	if err != nil {
		t.Fatalf("Failed to remove application metadata: %v", err)
	}
}

func ApplicationMetadata_Change_PreConfig(ctx context.Context, apiClient *management.APIClient, t *testing.T, environmentID, applicationID, metadata string) {
	if environmentID == "" || applicationID == "" {
		t.Fatalf("One of environment ID or application ID cannot be determined. Environment ID: %s, Application ID: %s", environmentID, applicationID)
	}

	metadataObject := management.NewApplicationMetadata()

	var metadataMap map[string]interface{}
	if err := json.Unmarshal([]byte(metadata), &metadataMap); err != nil {
		t.Fatalf("Failed to unmarshal metadata: %v", err)
	}
	metadataObject.SetMetadata(metadataMap)

	_, _, err := apiClient.ApplicationMetadataApi.UpdateApplicationMetadata(ctx, environmentID, applicationID).ApplicationMetadata(*metadataObject).Execute()
	if err != nil {
		t.Fatalf("Failed to update application metadata: %v", err)
	}
}
