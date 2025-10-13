// Copyright © 2025 Ping Identity Corporation

package base_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	acctestlegacysdk "github.com/pingidentity/terraform-provider-pingone/internal/acctest/legacysdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/base"
)

func TestAccCustomDomainSSL_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	domainPrefix := resourceName

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	certificateFile := os.Getenv("PINGONE_DOMAIN_CERTIFICATE_PEM")
	intermediateFile := os.Getenv("PINGONE_DOMAIN_INTERMEDIATE_CERTIFICATE_PEM")
	privateKeyFile := os.Getenv("PINGONE_DOMAIN_KEY_PEM")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckCustomDomain(t)
			acctest.PreCheckCustomDomainSSL(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.CustomDomain_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccCustomDomainSSLConfig_Full(environmentName, licenseID, resourceName, domainPrefix, certificateFile, intermediateFile, privateKeyFile),
				ExpectError: regexp.MustCompile(`Cannot add SSL certificate settings to the custom domain - Custom domain status must be 'SSL_CERTIFICATE_REQUIRED' or 'ACTIVE' in order to import a certificate`),
			},
		},
	})
}

func testAccCustomDomainSSLConfig_Full(environmentName, licenseID, resourceName, domainPrefix, certificateFile, intermediateFile, privateKeyFile string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_custom_domain" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  domain_name = "%[4]s.cdi-team-terraform-custom-domain-test.ping-eng.com"
}

resource "pingone_custom_domain_ssl" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  custom_domain_id = pingone_custom_domain.%[3]s.id

  certificate_pem_file               = <<EOT
%[5]s
EOT
  intermediate_certificates_pem_file = <<EOT
%[6]s
EOT
  private_key_pem_file               = <<EOT
%[7]s
EOT

}`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, domainPrefix, certificateFile, intermediateFile, privateKeyFile)
}
