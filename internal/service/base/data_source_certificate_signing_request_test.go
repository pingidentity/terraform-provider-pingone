package base_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
)

func TestAccCertificateSigningRequestDataSource_ByIDFull(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_certificate_signing_request.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	pkcs12 := os.Getenv("PINGONE_KEY_PKCS12")
	pkcs10_csr := os.Getenv("PINGONE_KEY_PKCS10_CSR")
	pem_csr := os.Getenv("PINGONE_KEY_PEM_CSR")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironmentAndPKCS12WithCSR(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCertificateSigningRequestDataSourceConfig_ByIDFull(environmentName, licenseID, resourceName, pkcs12),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceFullName, "pkcs10_file_base64", pkcs10_csr),
					resource.TestCheckResourceAttr(dataSourceFullName, "pem_file", pem_csr),
				),
			},
		},
	})
}

func TestAccCertificateSigningRequestDataSource_NotFound(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	resource.Test(t, resource.TestCase{
		// PreCheck:                 func() { acctest.PreCheckEnvironmentAndPKCS12WithCSR(t) },
		PreCheck:                 func() { t.Skipf("https://github.com/pingidentity/terraform-provider-pingone/issues/259") },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccCertificateSigningRequestDataSourceConfig_NotFoundByID(resourceName),
				ExpectError: regexp.MustCompile("Error when calling `GetKey`: Key not found for id: 9c052a8a-14be-44e4-8f07-2662569994ce and environmentId: [0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}"),
			},
		},
	})
}

func testAccCertificateSigningRequestDataSourceConfig_ByIDFull(environmentName, licenseID, resourceName, pkcs12 string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_key" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  pkcs12_file_base64 = <<EOT
%[4]s
EOT

  usage_type = "SIGNING"
}

data "pingone_certificate_signing_request" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  key_id = pingone_key.%[3]s.id
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, pkcs12)
}

func testAccCertificateSigningRequestDataSourceConfig_NotFoundByID(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_certificate_signing_request" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  key_id = "9c052a8a-14be-44e4-8f07-2662569994ce" // dummy ID that conforms to UUID v4
}
`, acctest.GenericSandboxEnvironment(), resourceName)
}
