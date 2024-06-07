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

func TrustedEmailDomain_CheckDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_trusted_email_domain" {
			continue
		}

		shouldContinue, err := acctest.CheckParentEnvironmentDestroy(ctx, apiClient, rs.Primary.Attributes["environment_id"])
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		_, r, err := apiClient.TrustedEmailDomainsApi.ReadOneTrustedEmailDomain(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.ID).Execute()

		shouldContinue, err = acctest.CheckForResourceDestroy(r, err)
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		return fmt.Errorf("PingOne Trusted Email Domain Instance %s still exists", rs.Primary.ID)
	}

	return nil
}

func TrustedEmailDomain_GetIDs(resourceName string, environmentID, resourceID *string) resource.TestCheckFunc {
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

func TrustedEmailDomain_RemovalDrift_PreConfig(ctx context.Context, apiClient *management.APIClient, t *testing.T, environmentID, trustedEmailDomainID string) {
	if environmentID == "" || trustedEmailDomainID == "" {
		t.Fatalf("One of environment ID or trusted email domain ID cannot be determined. Environment ID: %s, Trusted Email Domain ID: %s", environmentID, trustedEmailDomainID)
	}

	_, err := apiClient.TrustedEmailDomainsApi.DeleteTrustedEmailDomain(ctx, environmentID, trustedEmailDomainID).Execute()
	if err != nil {
		t.Fatalf("Failed to delete trusted email domain: %v", err)
	}
}
