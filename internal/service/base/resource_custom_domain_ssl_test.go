package base_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
)

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
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
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
