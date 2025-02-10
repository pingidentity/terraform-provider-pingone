// Copyright © 2025 Ping Identity Corporation

package base

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
)

func BrandingTheme_CheckDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_branding_theme" {
			continue
		}

		shouldContinue, err := acctest.CheckParentEnvironmentDestroy(ctx, apiClient, rs.Primary.Attributes["environment_id"])
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		_, r, err := apiClient.BrandingThemesApi.ReadOneBrandingTheme(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.ID).Execute()

		shouldContinue, err = acctest.CheckForResourceDestroy(r, err)
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		return fmt.Errorf("PingOne Branding Theme Instance %s still exists", rs.Primary.ID)
	}

	return nil
}

func BrandingTheme_GetIDs(resourceName string, environmentID, resourceID *string) resource.TestCheckFunc {
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

func BrandingTheme_RemovalDrift_PreConfig(ctx context.Context, apiClient *management.APIClient, t *testing.T, environmentID, brandingThemeID string) {
	if environmentID == "" || brandingThemeID == "" {
		t.Fatalf("One of environment ID or branding theme ID cannot be determined. Environment ID: %s, Branding Theme ID: %s", environmentID, brandingThemeID)
	}

	_, err := apiClient.BrandingThemesApi.DeleteBrandingTheme(ctx, environmentID, brandingThemeID).Execute()
	if err != nil {
		t.Fatalf("Failed to delete branding theme: %v", err)
	}
}
