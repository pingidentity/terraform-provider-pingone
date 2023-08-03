package base_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
)

func testAccGetMFAPolicyIDs(resourceName string, environmentID, resourceID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource Not found: %s", resourceName)
		}

		*resourceID = rs.Primary.ID
		*environmentID = rs.Primary.Attributes["environment_id"]

		return nil
	}
}

func TestAccMFAPolicy_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_policy.%s", resourceName)

	name := resourceName

	var resourceID, environmentID string

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckMFAPolicyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccMFAPolicyConfig_FullSMS(resourceName, name),
				Check:  testAccGetMFAPolicyIDs(resourceFullName, &environmentID, &resourceID),
			},
			// Replan after removal preconfig
			{
				PreConfig: func() {
					var ctx = context.Background()
					p1Client, err := acctest.TestClient(ctx)

					if err != nil {
						t.Fatalf("Failed to get API client: %v", err)
					}

					apiClient := p1Client.API.MFAAPIClient

					if environmentID == "" || resourceID == "" {
						t.Fatalf("One of environment ID or resource ID cannot be determined. Environment ID: %s, Resource ID: %s", environmentID, resourceID)
					}

					_, err = apiClient.DeviceAuthenticationPolicyApi.DeleteDeviceAuthenticationPolicy(ctx, environmentID, resourceID).Execute()
					if err != nil {
						t.Fatalf("Failed to delete MFA Policy: %v", err)
					}
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccCustomDomainSSL_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	certificateFile := os.Getenv("PINGONE_DOMAIN_CERTIFICATE_PEM")
	intermediateFile := os.Getenv("PINGONE_DOMAIN_INTERMEDIATE_CERTIFICATE_PEM")
	privateKeyFile := os.Getenv("PINGONE_DOMAIN_KEY_PEM")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironmentAndCustomDomainSSL(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCustomDomainDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccCustomDomainSSLConfig_Full(environmentName, licenseID, resourceName, certificateFile, intermediateFile, privateKeyFile),
				ExpectError: regexp.MustCompile(`Cannot add SSL certificate settings to the custom domain - Custom domain status must be 'SSL_CERTIFICATE_REQUIRED' or 'ACTIVE' in order to import a certificate`),
			},
		},
	})
}

func testAccCustomDomainSSLConfig_Full(environmentName, licenseID, resourceName, certificateFile, intermediateFile, privateKeyFile string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_custom_domain" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  domain_name = "terraformdev.ping-eng.com"
}

resource "pingone_custom_domain_ssl" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  custom_domain_id = pingone_custom_domain.%[3]s.id

  certificate_pem_file               = <<EOT
%[4]s
EOT
  intermediate_certificates_pem_file = <<EOT
%[5]s
EOT
  private_key_pem_file               = <<EOT
%[6]s
EOT

}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, certificateFile, intermediateFile, privateKeyFile)
}
