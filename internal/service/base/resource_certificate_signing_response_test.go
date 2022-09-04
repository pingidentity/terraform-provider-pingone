package base_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
)

func TestAccCertificateSigningResponse_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_certificate_signing_response.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	pkcs12 := os.Getenv("PINGONE_KEY_PKCS12")
	pemResponse := os.Getenv("PINGONE_KEY_PEM_CSR_RESPONSE")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironmentAndPKCS12WithCSRResponse(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      acctest.TestAccCheckEnvironmentDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCertificateSigningResponseConfig_Full(environmentName, licenseID, resourceName, pkcs12, pemResponse),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "name", "terraform (Test CA)"),
					resource.TestCheckResourceAttr(resourceFullName, "algorithm", "RSA"),
					resource.TestCheckResourceAttr(resourceFullName, "key_length", "4096"),
					resource.TestCheckResourceAttr(resourceFullName, "signature_algorithm", "SHA256withRSA"),
					resource.TestCheckResourceAttr(resourceFullName, "subject_dn", "CN=terraform, OU=Non-Production Testing, O=Ping Identity, C=GB"),
					resource.TestCheckResourceAttr(resourceFullName, "usage_type", "SIGNING"),
					resource.TestCheckResourceAttr(resourceFullName, "validity_period", "3650"),
					resource.TestCheckResourceAttr(resourceFullName, "issuer_dn", "CN=Test CA, OU=Non-Production Testing, O=Ping Identity, C=GB"),
					resource.TestCheckResourceAttr(resourceFullName, "default", "false"),
					resource.TestCheckResourceAttrSet(resourceFullName, "serial_number"),
					resource.TestMatchResourceAttr(resourceFullName, "expires_at", regexp.MustCompile(`^(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d\.\d+)|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d)|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d)$`)),
					resource.TestMatchResourceAttr(resourceFullName, "starts_at", regexp.MustCompile(`^(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d\.\d+)|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d)|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d)$`)),
					resource.TestCheckResourceAttr(resourceFullName, "status", "VALID"),
				),
			},
		},
	})
}

func testAccCertificateSigningResponseConfig_Full(environmentName, licenseID, resourceName, pkcs12, pemResponse string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_key" "%[3]s" {
	environment_id = pingone_environment.%[2]s.id
			  
	pkcs12_file_base64 = <<EOT
%[4]s
EOT
			  
	usage_type = "SIGNING"
}
	
resource "pingone_certificate_signing_response" "%[3]s" {
	environment_id = pingone_environment.%[2]s.id
	
	key_id = pingone_key.%[3]s.id
	pem_ca_response_file = <<EOT
%[5]s
EOT
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, pkcs12, pemResponse)
}
