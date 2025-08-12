// Copyright © 2025 Ping Identity Corporation

package base

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/pingidentity/pingone-go-client/pingone"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
)

func Environment_CheckDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_environment" {
			continue
		}

		environmentIdUuid, err := uuid.Parse(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("unable to parse environment id '%s' as uuid: %v", environmentIdUuid, err)
		}

		_, r, err := p1Client.EnvironmentsApi.GetEnvironmentById(ctx, environmentIdUuid).Execute()

		if r.StatusCode == 404 {
			continue
		}

		if err != nil {
			return err
		}

		return fmt.Errorf("PingOne Environment Instance %s still exists", rs.Primary.ID)
	}

	return nil
}

func Environment_GetIDs(resourceName string, resourceID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource Not found: %s", resourceName)
		}

		if resourceID != nil {
			*resourceID = rs.Primary.ID
		}

		return nil
	}
}

func Environment_RemovalDrift_PreConfig(ctx context.Context, apiClient *pingone.APIClient, t *testing.T, environmentID string) {
	environmentIdUuid, err := uuid.Parse(environmentID)
	if err != nil {
		t.Fatalf("unable to parse environment id '%s' as uuid: %v", environmentIdUuid, err)
	}

	if environmentID == "" {
		t.Fatalf("Environment ID cannot be determined. Environment ID: %s", environmentID)
	}

	_, err = apiClient.EnvironmentsApi.DeleteEnvironmentById(ctx, environmentIdUuid).Execute()
	if err != nil {
		t.Fatalf("Failed to delete environment: %v", err)
	}
}
