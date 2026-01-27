// Copyright © 2025 Ping Identity Corporation

package sso

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/legacysdk"
)

func User_CheckDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := legacysdk.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_user" {
			continue
		}

		shouldContinue, err := legacysdk.CheckParentEnvironmentDestroy(ctx, p1Client.API.ManagementAPIClient, rs.Primary.Attributes["environment_id"])
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		_, r, err := apiClient.UsersApi.ReadUser(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.ID).Execute()

		shouldContinue, err = acctest.CheckForResourceDestroy(r, err)
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		return fmt.Errorf("PingOne User %s still exists", rs.Primary.ID)
	}

	return nil
}

func User_GetIDs(resourceName string, environmentID, resourceID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
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

func User_RemovalDrift_PreConfig(ctx context.Context, apiClient *management.APIClient, t *testing.T, environmentID, userID string) {
	if environmentID == "" || userID == "" {
		t.Fatalf("One of environment ID or user ID cannot be determined. Environment ID: %s, User ID: %s", environmentID, userID)
	}

	_, err := apiClient.UsersApi.DeleteUser(ctx, environmentID, userID).Execute()
	if err != nil {
		t.Fatalf("Failed to delete user: %v", err)
	}
}

func User_CreateUser_PreConfig(ctx context.Context, apiClient *management.APIClient, t *testing.T, environmentID, name string, resourceID, populationID *string) {
	if environmentID == "" || name == "" {
		t.Fatalf("One of environment ID or user name cannot be determined. Environment ID: %s, User name: %s", environmentID, name)
	}

	userData := management.NewUser(
		fmt.Sprintf("%s@ping-eng.com", name),
		name,
	)

	if populationID != nil {
		population := *management.NewUserPopulation(*populationID)
		userData.SetPopulation(population)
	}

	fO, _, fErr := apiClient.UsersApi.CreateUser(context.Background(), environmentID).ContentType("application/vnd.pingidentity.user.import+json").User(*userData).Execute()
	if fErr != nil {
		t.Fatalf("Failed to create user: %v", fErr)
	}

	*resourceID = fO.GetId()
}
