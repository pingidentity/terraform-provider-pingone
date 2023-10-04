package base_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/base"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccCertificate_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_certificate.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	pem_cert := os.Getenv("PINGONE_KEY_PEM_CERT")

	var certificateID, environmentID string

	var ctx = context.Background()
	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		t.Fatalf("Failed to get API client: %v", err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.TestAccCheckCertificateDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccCertificateConfig_PEM(environmentName, licenseID, resourceName, pem_cert),
				Check:  base.TestAccGetCertificateIDs(resourceFullName, &environmentID, &certificateID),
			},
			// Replan after removal preconfig
			{
				PreConfig: func() {
					base.Certificate_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID, certificateID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the environment
			{
				Config: testAccCertificateConfig_PEM(environmentName, licenseID, resourceName, pem_cert),
				Check:  base.TestAccGetCertificateIDs(resourceFullName, &environmentID, &certificateID),
			},
			{
				PreConfig: func() {
					base.Environment_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccCertificate_PKCS7(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_certificate.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	pkcs7_cert := os.Getenv("PINGONE_KEY_PKCS7_CERT")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckPKCS7Cert(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.TestAccCheckCertificateDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCertificateConfig_PKCS7(environmentName, licenseID, resourceName, pkcs7_cert),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", "terraform"),
					resource.TestCheckResourceAttr(resourceFullName, "algorithm", "RSA"),
					resource.TestCheckResourceAttr(resourceFullName, "key_length", "4096"),
					resource.TestCheckResourceAttr(resourceFullName, "signature_algorithm", "SHA256withRSA"),
					resource.TestCheckResourceAttr(resourceFullName, "subject_dn", "C=GB,O=Ping Identity,OU=Non-Production Testing,CN=terraform"),
					resource.TestCheckResourceAttr(resourceFullName, "usage_type", "SSL/TLS"),
					resource.TestCheckResourceAttr(resourceFullName, "validity_period", "3650"),
					resource.TestCheckResourceAttr(resourceFullName, "issuer_dn", "C=GB,O=Ping Identity,OU=Non-Production Testing,CN=terraform"),
					resource.TestCheckResourceAttr(resourceFullName, "default", "false"),
					resource.TestCheckResourceAttrSet(resourceFullName, "serial_number"),
					resource.TestMatchResourceAttr(resourceFullName, "expires_at", regexp.MustCompile(`^(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d\.\d+)|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d)|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d)$`)),
					resource.TestMatchResourceAttr(resourceFullName, "starts_at", regexp.MustCompile(`^(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d\.\d+)|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d)|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d)$`)),
					resource.TestCheckResourceAttr(resourceFullName, "status", "VALID"),
				),
			},
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"pkcs7_file_base64",
					"pem_file",
				},
			},
		},
	})
}

func TestAccCertificate_PEM(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_certificate.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	pem_cert := os.Getenv("PINGONE_KEY_PEM_CERT")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckPEMCert(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.TestAccCheckCertificateDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCertificateConfig_PEM(environmentName, licenseID, resourceName, pem_cert),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", "terraform"),
					resource.TestCheckResourceAttr(resourceFullName, "algorithm", "RSA"),
					resource.TestCheckResourceAttr(resourceFullName, "key_length", "4096"),
					resource.TestCheckResourceAttr(resourceFullName, "signature_algorithm", "SHA256withRSA"),
					resource.TestCheckResourceAttr(resourceFullName, "subject_dn", "C=GB,O=Ping Identity,OU=Non-Production Testing,CN=terraform"),
					resource.TestCheckResourceAttr(resourceFullName, "usage_type", "SSL/TLS"),
					resource.TestCheckResourceAttr(resourceFullName, "validity_period", "3650"),
					resource.TestCheckResourceAttr(resourceFullName, "issuer_dn", "C=GB,O=Ping Identity,OU=Non-Production Testing,CN=terraform"),
					resource.TestCheckResourceAttr(resourceFullName, "default", "false"),
					resource.TestCheckResourceAttrSet(resourceFullName, "serial_number"),
					resource.TestMatchResourceAttr(resourceFullName, "expires_at", regexp.MustCompile(`^(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d\.\d+)|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d)|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d)$`)),
					resource.TestMatchResourceAttr(resourceFullName, "starts_at", regexp.MustCompile(`^(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d\.\d+)|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d)|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d)$`)),
					resource.TestCheckResourceAttr(resourceFullName, "status", "VALID"),
				),
			},
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"pkcs7_file_base64",
					"pem_file",
				},
			},
		},
	})
}

func TestAccCertificate_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_certificate.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	pem_cert := os.Getenv("PINGONE_KEY_PEM_CERT")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.TestAccCheckCertificateDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccCertificateConfig_PEM(environmentName, licenseID, resourceName, pem_cert),
			},
			// Errors
			{
				ResourceName: resourceFullName,
				ImportState:  true,
				ExpectError:  regexp.MustCompile(`Invalid import ID specified \(".*"\).  The ID should be in the format "environment_id/certificate_id" and must match regex: .*`),
			},
			{
				ResourceName:  resourceFullName,
				ImportStateId: "/",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`Invalid import ID specified \(".*"\).  The ID should be in the format "environment_id/certificate_id" and must match regex: .*`),
			},
			{
				ResourceName:  resourceFullName,
				ImportStateId: "badformat/badformat",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`Invalid import ID specified \(".*"\).  The ID should be in the format "environment_id/certificate_id" and must match regex: .*`),
			},
		},
	})
}

func testAccCertificateConfig_PKCS7(environmentName, licenseID, resourceName, pkcs12 string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_certificate" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  pkcs7_file_base64 = <<EOT
%[4]s
EOT

  usage_type = "SSL/TLS"
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, pkcs12)
}

func testAccCertificateConfig_PEM(environmentName, licenseID, resourceName, pem string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_certificate" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  pem_file = <<EOT
%[4]s
EOT

  usage_type = "SSL/TLS"
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, pem)
}
