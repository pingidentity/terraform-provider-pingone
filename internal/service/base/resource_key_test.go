package base_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func testAccCheckKeyDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_key" {
			continue
		}

		_, rEnv, err := apiClient.EnvironmentsApi.ReadOneEnvironment(ctx, rs.Primary.Attributes["environment_id"]).Execute()

		if err != nil {

			if rEnv == nil {
				return fmt.Errorf("Response object does not exist and no error detected")
			}

			if rEnv.StatusCode == 404 {
				continue
			}

			return err
		}

		body, r, err := apiClient.CertificateManagementApi.GetKey(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.ID).Accept(management.ENUMGETKEYACCEPTHEADER_JSON).Execute()

		if err != nil {

			if r == nil {
				return fmt.Errorf("Response object does not exist and no error detected")
			}

			if r.StatusCode == 404 {
				continue
			}

			tflog.Error(ctx, fmt.Sprintf("Error: %v", body))
			return err
		}

		return fmt.Errorf("PingOne Key Instance %s still exists", rs.Primary.ID)
	}

	return nil
}

func testAccGetKeyIDs(resourceName string, environmentID, resourceID *string) resource.TestCheckFunc {
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

func TestAccKey_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_key.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	var resourceID, environmentID string

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckKeyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccKeyConfig_Minimal(environmentName, licenseID, resourceName, name),
				Check:  testAccGetKeyIDs(resourceFullName, &environmentID, &resourceID),
			},
			// Replan after removal preconfig
			{
				PreConfig: func() {
					var ctx = context.Background()
					p1Client, err := acctest.TestClient(ctx)

					if err != nil {
						t.Fatalf("Failed to get API client: %v", err)
					}

					apiClient := p1Client.API.ManagementAPIClient

					if environmentID == "" || resourceID == "" {
						t.Fatalf("One of environment ID or resource ID cannot be determined. Environment ID: %s, Resource ID: %s", environmentID, resourceID)
					}

					_, err = apiClient.CertificateManagementApi.DeleteKey(ctx, environmentID, resourceID).Execute()
					if err != nil {
						t.Fatalf("Failed to delete Key: %v", err)
					}
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccKey_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_key.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckKeyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccKeyConfig_Full(environmentName, licenseID, resourceName, name, "ENCRYPTION", true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "algorithm", "RSA"),
					resource.TestCheckResourceAttr(resourceFullName, "key_length", "3072"),
					resource.TestCheckResourceAttr(resourceFullName, "signature_algorithm", "SHA512withRSA"),
					resource.TestCheckResourceAttr(resourceFullName, "subject_dn", fmt.Sprintf("CN=%s, OU=Ping Identity, O=Ping Identity, L=, ST=, C=US", name)),
					resource.TestCheckResourceAttr(resourceFullName, "usage_type", "ENCRYPTION"),
					resource.TestCheckResourceAttr(resourceFullName, "validity_period", "3650"),
					resource.TestCheckResourceAttr(resourceFullName, "issuer_dn", "CN=My CA, OU=Ping Identity, O=Ping Identity, L=, ST=, C=US"),
					resource.TestCheckResourceAttr(resourceFullName, "default", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "serial_number", "1662023413215"),
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
			},
			{
				Config: testAccKeyConfig_Full(environmentName, licenseID, resourceName, name, "SIGNING", false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "algorithm", "RSA"),
					resource.TestCheckResourceAttr(resourceFullName, "key_length", "3072"),
					resource.TestCheckResourceAttr(resourceFullName, "signature_algorithm", "SHA512withRSA"),
					resource.TestCheckResourceAttr(resourceFullName, "subject_dn", fmt.Sprintf("CN=%s, OU=Ping Identity, O=Ping Identity, L=, ST=, C=US", name)),
					resource.TestCheckResourceAttr(resourceFullName, "usage_type", "SIGNING"),
					resource.TestCheckResourceAttr(resourceFullName, "validity_period", "3650"),
					resource.TestCheckResourceAttr(resourceFullName, "issuer_dn", "CN=My CA, OU=Ping Identity, O=Ping Identity, L=, ST=, C=US"),
					resource.TestCheckResourceAttr(resourceFullName, "default", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "serial_number", "1662023413215"),
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
			},
			{
				Config: testAccKeyConfig_Full(environmentName, licenseID, resourceName, name, "SSL/TLS", false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "usage_type", "SSL/TLS"),
					resource.TestCheckResourceAttr(resourceFullName, "default", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "issuer_dn", "CN=My CA, OU=Ping Identity, O=Ping Identity, L=, ST=, C=US"),
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
			},
			{
				Config: testAccKeyConfig_Full(environmentName, licenseID, resourceName, name, "ISSUANCE", false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "usage_type", "ISSUANCE"),
					resource.TestCheckResourceAttr(resourceFullName, "default", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "issuer_dn", "CN=My CA, OU=Ping Identity, O=Ping Identity, L=, ST=, C=US"),
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
			},
		},
	})
}

func TestAccKey_Minimal(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_key.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckKeyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccKeyConfig_Minimal(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "algorithm", "EC"),
					resource.TestCheckResourceAttr(resourceFullName, "key_length", "256"),
					resource.TestCheckResourceAttr(resourceFullName, "signature_algorithm", "SHA384withECDSA"),
					resource.TestCheckResourceAttr(resourceFullName, "subject_dn", fmt.Sprintf("CN=%s, OU=Ping Identity, O=Ping Identity, L=, ST=, C=US", name)),
					resource.TestCheckResourceAttr(resourceFullName, "usage_type", "SIGNING"),
					resource.TestCheckResourceAttr(resourceFullName, "validity_period", "365"),
					resource.TestCheckResourceAttr(resourceFullName, "issuer_dn", fmt.Sprintf("CN=%s, OU=Ping Identity, O=Ping Identity, L=, ST=, C=US", name)),
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

func TestAccKey_PKCS12(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_key.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	pkcs12 := os.Getenv("PINGONE_KEY_PKCS12")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironmentAndPKCS12(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckKeyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccKeyConfig_PKCS12(environmentName, licenseID, resourceName, pkcs12),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", "terraform"),
					resource.TestCheckResourceAttr(resourceFullName, "algorithm", "RSA"),
					resource.TestCheckResourceAttr(resourceFullName, "key_length", "4096"),
					resource.TestCheckResourceAttr(resourceFullName, "signature_algorithm", "SHA256withRSA"),
					resource.TestCheckResourceAttr(resourceFullName, "subject_dn", "CN=terraform, OU=Non-Production Testing, O=Ping Identity, C=GB"),
					resource.TestCheckResourceAttr(resourceFullName, "usage_type", "SIGNING"),
					resource.TestCheckResourceAttr(resourceFullName, "validity_period", "3650"),
					resource.TestCheckResourceAttr(resourceFullName, "issuer_dn", "CN=terraform, OU=Non-Production Testing, O=Ping Identity, C=GB"),
					resource.TestCheckResourceAttr(resourceFullName, "default", "false"),
					resource.TestCheckResourceAttrSet(resourceFullName, "serial_number"),
					resource.TestMatchResourceAttr(resourceFullName, "expires_at", regexp.MustCompile(`^(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d\.\d+)|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d)|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d)$`)),
					resource.TestMatchResourceAttr(resourceFullName, "starts_at", regexp.MustCompile(`^(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d\.\d+)|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d)|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d)$`)),
					resource.TestCheckResourceAttr(resourceFullName, "status", "VALID"),
				),
			},
			{
				Config: testAccKeyConfig_PKCS12(environmentName, licenseID, resourceName, pkcs12),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", "terraform"),
					resource.TestCheckResourceAttr(resourceFullName, "algorithm", "RSA"),
					resource.TestCheckResourceAttr(resourceFullName, "key_length", "4096"),
					resource.TestCheckResourceAttr(resourceFullName, "signature_algorithm", "SHA256withRSA"),
					resource.TestCheckResourceAttr(resourceFullName, "subject_dn", "CN=terraform, OU=Non-Production Testing, O=Ping Identity, C=GB"),
					resource.TestCheckResourceAttr(resourceFullName, "usage_type", "SIGNING"),
					resource.TestCheckResourceAttr(resourceFullName, "validity_period", "3650"),
					resource.TestCheckResourceAttr(resourceFullName, "issuer_dn", "CN=terraform, OU=Non-Production Testing, O=Ping Identity, C=GB"),
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

func TestAccKey_Change(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_key.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	pkcs12 := os.Getenv("PINGONE_KEY_PKCS12")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironmentAndPKCS12(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckKeyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccKeyConfig_Full(environmentName, licenseID, resourceName, name, "ENCRYPTION", true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "algorithm", "RSA"),
					resource.TestCheckResourceAttr(resourceFullName, "key_length", "3072"),
					resource.TestCheckResourceAttr(resourceFullName, "signature_algorithm", "SHA512withRSA"),
					resource.TestCheckResourceAttr(resourceFullName, "subject_dn", fmt.Sprintf("CN=%s, OU=Ping Identity, O=Ping Identity, L=, ST=, C=US", name)),
					resource.TestCheckResourceAttr(resourceFullName, "usage_type", "ENCRYPTION"),
					resource.TestCheckResourceAttr(resourceFullName, "validity_period", "3650"),
					resource.TestCheckResourceAttr(resourceFullName, "issuer_dn", "CN=My CA, OU=Ping Identity, O=Ping Identity, L=, ST=, C=US"),
					resource.TestCheckResourceAttr(resourceFullName, "default", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "serial_number", "1662023413215"),
					resource.TestMatchResourceAttr(resourceFullName, "expires_at", regexp.MustCompile(`^(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d\.\d+)|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d)|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d)$`)),
					resource.TestMatchResourceAttr(resourceFullName, "starts_at", regexp.MustCompile(`^(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d\.\d+)|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d)|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d)$`)),
					resource.TestCheckResourceAttr(resourceFullName, "status", "VALID"),
				),
			},
			{
				Config: testAccKeyConfig_Full(environmentName, licenseID, resourceName, name, "ENCRYPTION", false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "default", "false"),
				),
			},
			{
				Config: testAccKeyConfig_Minimal(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "algorithm", "EC"),
					resource.TestCheckResourceAttr(resourceFullName, "key_length", "256"),
					resource.TestCheckResourceAttr(resourceFullName, "signature_algorithm", "SHA384withECDSA"),
					resource.TestCheckResourceAttr(resourceFullName, "subject_dn", fmt.Sprintf("CN=%s, OU=Ping Identity, O=Ping Identity, L=, ST=, C=US", name)),
					resource.TestCheckResourceAttr(resourceFullName, "usage_type", "SIGNING"),
					resource.TestCheckResourceAttr(resourceFullName, "validity_period", "365"),
					resource.TestCheckResourceAttr(resourceFullName, "issuer_dn", fmt.Sprintf("CN=%s, OU=Ping Identity, O=Ping Identity, L=, ST=, C=US", name)),
					resource.TestCheckResourceAttr(resourceFullName, "default", "false"),
					resource.TestCheckResourceAttrSet(resourceFullName, "serial_number"),
					resource.TestMatchResourceAttr(resourceFullName, "expires_at", regexp.MustCompile(`^(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d\.\d+)|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d)|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d)$`)),
					resource.TestMatchResourceAttr(resourceFullName, "starts_at", regexp.MustCompile(`^(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d\.\d+)|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d)|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d)$`)),
					resource.TestCheckResourceAttr(resourceFullName, "status", "VALID"),
				),
			},
			{
				Config: testAccKeyConfig_PKCS12(environmentName, licenseID, resourceName, pkcs12),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", "terraform"),
					resource.TestCheckResourceAttr(resourceFullName, "algorithm", "RSA"),
					resource.TestCheckResourceAttr(resourceFullName, "key_length", "4096"),
					resource.TestCheckResourceAttr(resourceFullName, "signature_algorithm", "SHA256withRSA"),
					resource.TestCheckResourceAttr(resourceFullName, "subject_dn", "CN=terraform, OU=Non-Production Testing, O=Ping Identity, C=GB"),
					resource.TestCheckResourceAttr(resourceFullName, "usage_type", "SIGNING"),
					resource.TestCheckResourceAttr(resourceFullName, "validity_period", "3650"),
					resource.TestCheckResourceAttr(resourceFullName, "issuer_dn", "CN=terraform, OU=Non-Production Testing, O=Ping Identity, C=GB"),
					resource.TestCheckResourceAttr(resourceFullName, "default", "false"),
					resource.TestCheckResourceAttrSet(resourceFullName, "serial_number"),
					resource.TestMatchResourceAttr(resourceFullName, "expires_at", regexp.MustCompile(`^(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d\.\d+)|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d)|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d)$`)),
					resource.TestMatchResourceAttr(resourceFullName, "starts_at", regexp.MustCompile(`^(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d\.\d+)|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d)|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d)$`)),
					resource.TestCheckResourceAttr(resourceFullName, "status", "VALID"),
				),
			},
			{
				Config: testAccKeyConfig_Full(environmentName, licenseID, resourceName, name, "ENCRYPTION", false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "algorithm", "RSA"),
					resource.TestCheckResourceAttr(resourceFullName, "key_length", "3072"),
					resource.TestCheckResourceAttr(resourceFullName, "signature_algorithm", "SHA512withRSA"),
					resource.TestCheckResourceAttr(resourceFullName, "subject_dn", fmt.Sprintf("CN=%s, OU=Ping Identity, O=Ping Identity, L=, ST=, C=US", name)),
					resource.TestCheckResourceAttr(resourceFullName, "usage_type", "ENCRYPTION"),
					resource.TestCheckResourceAttr(resourceFullName, "validity_period", "3650"),
					resource.TestCheckResourceAttr(resourceFullName, "issuer_dn", "CN=My CA, OU=Ping Identity, O=Ping Identity, L=, ST=, C=US"),
					resource.TestCheckResourceAttr(resourceFullName, "default", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "serial_number", "1662023413215"),
					resource.TestMatchResourceAttr(resourceFullName, "expires_at", regexp.MustCompile(`^(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d\.\d+)|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d)|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d)$`)),
					resource.TestMatchResourceAttr(resourceFullName, "starts_at", regexp.MustCompile(`^(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d\.\d+)|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d)|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d)$`)),
					resource.TestCheckResourceAttr(resourceFullName, "status", "VALID"),
				),
			},
		},
	})
}

func TestAccKey_CustomCRL(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_key.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironmentAndPKCS12(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckKeyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccKeyConfig_CustomCRL(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "usage_type", "ISSUANCE"),
					resource.TestCheckResourceAttr(resourceFullName, "custom_crl", "http://www.pingidentity.com/test.crl"),
				),
			},
			{
				Config:  testAccKeyConfig_CustomCRL(environmentName, licenseID, resourceName, name),
				Destroy: true,
			},
			{
				Config:      testAccKeyConfig_CustomCRLIncorrectKeyType(environmentName, licenseID, resourceName, name),
				ExpectError: regexp.MustCompile("Invalid attribute combination"),
			},
		},
	})
}

func TestAccKey_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_key.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckKeyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccKeyConfig_Minimal(environmentName, licenseID, resourceName, name),
			},
			// Errors
			{
				ResourceName: resourceFullName,
				ImportState:  true,
				ExpectError:  regexp.MustCompile(`Unexpected Import Identifier`),
			},
			{
				ResourceName:  resourceFullName,
				ImportStateId: "/",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`Unexpected Import Identifier`),
			},
			{
				ResourceName:  resourceFullName,
				ImportStateId: "badformat/badformat",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`Unexpected Import Identifier`),
			},
		},
	})
}

func testAccKeyConfig_Full(environmentName, licenseID, resourceName, name, usage string, defaultKey bool) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_key" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  name                = "%[4]s"
  algorithm           = "RSA"
  key_length          = 3072
  signature_algorithm = "SHA512withRSA"
  subject_dn          = "CN=%[4]s, OU=Ping Identity, O=Ping Identity, L=, ST=, C=US"
  usage_type          = "%[5]s"

  default         = %[6]t
  issuer_dn       = "CN=My CA, OU=Ping Identity, O=Ping Identity, L=, ST=, C=US"
  serial_number   = "1662023413215"
  validity_period = 3650
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name, usage, defaultKey)
}

func testAccKeyConfig_Minimal(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_key" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  name                = "%[4]s"
  algorithm           = "EC"
  key_length          = 256
  signature_algorithm = "SHA384withECDSA"
  subject_dn          = "CN=%[4]s, OU=Ping Identity, O=Ping Identity, L=, ST=, C=US"
  usage_type          = "SIGNING"
  validity_period     = 365
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccKeyConfig_CustomCRL(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_key" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  name                = "%[4]s"
  algorithm           = "EC"
  key_length          = 256
  signature_algorithm = "SHA384withECDSA"
  subject_dn          = "CN=%[4]s, OU=Ping Identity, O=Ping Identity, L=, ST=, C=US"
  usage_type          = "ISSUANCE"
  validity_period     = 365

  custom_crl = "http://www.pingidentity.com/test.crl"
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccKeyConfig_CustomCRLIncorrectKeyType(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_key" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  name                = "%[4]s"
  algorithm           = "EC"
  key_length          = 256
  signature_algorithm = "SHA384withECDSA"
  subject_dn          = "CN=%[4]s, OU=Ping Identity, O=Ping Identity, L=, ST=, C=US"
  usage_type          = "SIGNING"
  validity_period     = 365

  custom_crl = "http://www.pingidentity.com/test.crl"
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccKeyConfig_PKCS12(environmentName, licenseID, resourceName, pkcs12 string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_key" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  pkcs12_file_base64 = <<EOT
%[4]s
EOT

  usage_type = "SIGNING"
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, pkcs12)
}
