// Copyright © 2025 Ping Identity Corporation

package base_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccCertificateDataSource_ByNameFull(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_certificate.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	pem_cert := os.Getenv("PINGONE_KEY_PEM_CERT")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckPEMCert(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCertificateDataSourceConfig_ByNameFull(environmentName, licenseID, resourceName, pem_cert),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "name", "terraform"),
					resource.TestCheckResourceAttr(dataSourceFullName, "algorithm", "RSA"),
					resource.TestCheckResourceAttr(dataSourceFullName, "key_length", "4096"),
					resource.TestCheckResourceAttr(dataSourceFullName, "signature_algorithm", "SHA256withRSA"),
					resource.TestCheckResourceAttr(dataSourceFullName, "subject_dn", "C=GB,O=Ping Identity,OU=Non-Production Testing,CN=terraform"),
					resource.TestCheckResourceAttr(dataSourceFullName, "usage_type", "SIGNING"),
					resource.TestCheckResourceAttr(dataSourceFullName, "validity_period", "3650"),
					resource.TestCheckResourceAttr(dataSourceFullName, "issuer_dn", "C=GB,O=Ping Identity,OU=Non-Production Testing,CN=terraform"),
					resource.TestCheckResourceAttr(dataSourceFullName, "default", "false"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "serial_number"),
					resource.TestMatchResourceAttr(dataSourceFullName, "expires_at", regexp.MustCompile(`^(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d\.\d+)|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d)|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d)$`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "starts_at", regexp.MustCompile(`^(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d\.\d+)|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d)|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d)$`)),
					resource.TestCheckResourceAttr(dataSourceFullName, "status", "VALID"),
				),
			},
		},
	})
}

func TestAccCertificateDataSource_ByIDFull(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_certificate.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	pem_cert := os.Getenv("PINGONE_KEY_PEM_CERT")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckPEMCert(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCertificateDataSourceConfig_ByIDFull(environmentName, licenseID, resourceName, pem_cert),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "name", "terraform"),
					resource.TestCheckResourceAttr(dataSourceFullName, "algorithm", "RSA"),
					resource.TestCheckResourceAttr(dataSourceFullName, "key_length", "4096"),
					resource.TestCheckResourceAttr(dataSourceFullName, "signature_algorithm", "SHA256withRSA"),
					resource.TestCheckResourceAttr(dataSourceFullName, "subject_dn", "C=GB,O=Ping Identity,OU=Non-Production Testing,CN=terraform"),
					resource.TestCheckResourceAttr(dataSourceFullName, "usage_type", "SIGNING"),
					resource.TestCheckResourceAttr(dataSourceFullName, "validity_period", "3650"),
					resource.TestCheckResourceAttr(dataSourceFullName, "issuer_dn", "C=GB,O=Ping Identity,OU=Non-Production Testing,CN=terraform"),
					resource.TestCheckResourceAttr(dataSourceFullName, "default", "false"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "serial_number"),
					resource.TestMatchResourceAttr(dataSourceFullName, "expires_at", regexp.MustCompile(`^(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d\.\d+)|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d)|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d)$`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "starts_at", regexp.MustCompile(`^(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d\.\d+)|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d)|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d)$`)),
					resource.TestCheckResourceAttr(dataSourceFullName, "status", "VALID"),
				),
			},
		},
	})
}

func TestAccCertificateDataSource_NotFound(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
			acctest.PreCheckPEMCert(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccCertificateDataSourceConfig_NotFoundByName(resourceName),
				ExpectError: regexp.MustCompile("Cannot find certificate doesnotexist"),
			},
			{
				Config:      testAccCertificateDataSourceConfig_NotFoundByID(resourceName),
				ExpectError: regexp.MustCompile("Error when calling `GetCertificate`: Certificate not found for id: 9c052a8a-14be-44e4-8f07-2662569994ce and environmentId: [0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}"),
			},
		},
	})
}

func testAccCertificateDataSourceConfig_ByNameFull(environmentName, licenseID, resourceName, pem string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_certificate" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  pem_file = <<EOT
%[4]s
EOT

  usage_type = "SIGNING"
}

data "pingone_certificate" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  name = "terraform"

  depends_on = [
    pingone_certificate.%[3]s
  ]
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, pem)
}

func testAccCertificateDataSourceConfig_ByIDFull(environmentName, licenseID, resourceName, pem string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_certificate" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  pem_file = <<EOT
%[4]s
EOT

  usage_type = "SIGNING"
}

data "pingone_certificate" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  certificate_id = pingone_certificate.%[3]s.id
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, pem)
}

func testAccCertificateDataSourceConfig_NotFoundByName(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

data "pingone_certificate" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "doesnotexist"
}`, acctest.GenericSandboxEnvironment(), resourceName)
}

func testAccCertificateDataSourceConfig_NotFoundByID(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

data "pingone_certificate" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  certificate_id = "9c052a8a-14be-44e4-8f07-2662569994ce" // dummy ID that conforms to UUID v4
}`, acctest.GenericSandboxEnvironment(), resourceName)
}
