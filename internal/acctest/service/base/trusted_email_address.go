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

func TrustedEmailAddress_CheckDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_trusted_email_address" {
			continue
		}

		shouldContinue, err := acctest.CheckParentEnvironmentDestroy(ctx, apiClient, rs.Primary.Attributes["environment_id"])
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		_, r, err := apiClient.TrustedEmailAddressesApi.ReadOneTrustedEmailAddress(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["email_domain_id"], rs.Primary.ID).Execute()

		shouldContinue, err = acctest.CheckForResourceDestroy(r, err)
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		return fmt.Errorf("PingOne trusted email address %s still exists", rs.Primary.ID)
	}

	return nil
}

func TrustedEmailAddress_GetIDs(resourceName string, environmentID, emailDomainID, resourceID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource Not found: %s", resourceName)
		}

		*resourceID = rs.Primary.ID
		*emailDomainID = rs.Primary.Attributes["email_domain_id"]
		*environmentID = rs.Primary.Attributes["environment_id"]

		return nil
	}
}

func TrustedEmailAddress_RemovalDrift_PreConfig(ctx context.Context, apiClient *management.APIClient, t *testing.T, environmentID, emailDomainID, trustedEmailAddressID string) {
	if environmentID == "" || emailDomainID == "" || trustedEmailAddressID == "" {
		t.Fatalf("One of environment ID, email domain ID or trusted email address ID cannot be determined. Environment ID: %s, Email Domain ID: %s, Trusted Email Address ID: %s", environmentID, emailDomainID, trustedEmailAddressID)
	}

	_, err := apiClient.TrustedEmailAddressesApi.DeleteTrustedEmailAddress(ctx, environmentID, emailDomainID, trustedEmailAddressID).Execute()
	if err != nil {
		t.Fatalf("Failed to delete trusted email address: %v", err)
	}
}
