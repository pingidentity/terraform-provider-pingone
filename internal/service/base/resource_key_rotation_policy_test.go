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
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func testAccCheckKeyRotationPolicyDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_key_rotation_policy" {
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

		body, r, err := apiClient.KeyRotationPoliciesApi.GetKeyRotationPolicy(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.ID).Execute()

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

		return fmt.Errorf("PingOne Key Rotation Policy Instance %s still exists", rs.Primary.ID)
	}

	return nil
}

func TestAccKeyRotationPolicy_NewEnv(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_key_rotation_policy.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckKeyRotationPolicyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccKeyRotationPolicyConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
				),
			},
		},
	})
}

func TestAccKeyRotationPolicy_All(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_key_rotation_policy.%s", resourceName)

	nameFull := fmt.Sprintf("%sfull", resourceName)
	nameMin := fmt.Sprintf("%smin", resourceName)

	fullTest := resource.TestStep{
		Config: testAccKeyRotationPolicyConfig_Full(resourceName, nameFull),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
			resource.TestCheckResourceAttr(resourceFullName, "name", nameFull),
			resource.TestCheckResourceAttr(resourceFullName, "algorithm", "RSA"),
			resource.TestMatchResourceAttr(resourceFullName, "current_key_id", verify.P1ResourceIDRegexp),
			resource.TestMatchResourceAttr(resourceFullName, "next_key_id", verify.P1ResourceIDRegexp),
			resource.TestCheckResourceAttr(resourceFullName, "subject_dn", fmt.Sprintf("CN=%s, OU=Ping Identity, O=Ping Identity, L=, ST=, C=US", nameFull)),
			resource.TestCheckResourceAttr(resourceFullName, "key_length", "3072"),
			resource.TestCheckResourceAttr(resourceFullName, "rotation_period", "31"),
			resource.TestCheckResourceAttr(resourceFullName, "signature_algorithm", "SHA256withRSA"),
			resource.TestCheckResourceAttr(resourceFullName, "usage_type", "SIGNING"),
			resource.TestCheckResourceAttr(resourceFullName, "validity_period", "340"),
			resource.TestMatchResourceAttr(resourceFullName, "rotated_at", regexp.MustCompile(`^(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d\.\d+)|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d)|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d)$`)),
		),
	}

	minimalTest := resource.TestStep{
		Config: testAccKeyRotationPolicyConfig_Minimal(resourceName, nameMin),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
			resource.TestCheckResourceAttr(resourceFullName, "name", nameMin),
			resource.TestCheckResourceAttr(resourceFullName, "algorithm", "RSA"),
			resource.TestMatchResourceAttr(resourceFullName, "current_key_id", verify.P1ResourceIDRegexp),
			resource.TestMatchResourceAttr(resourceFullName, "next_key_id", verify.P1ResourceIDRegexp),
			resource.TestCheckResourceAttr(resourceFullName, "subject_dn", fmt.Sprintf("CN=%s, OU=Ping Identity, O=Ping Identity, L=, ST=, C=US", nameMin)),
			resource.TestCheckResourceAttr(resourceFullName, "key_length", "2048"),
			resource.TestCheckResourceAttr(resourceFullName, "rotation_period", "90"),
			resource.TestCheckResourceAttr(resourceFullName, "signature_algorithm", "SHA256withRSA"),
			resource.TestCheckResourceAttr(resourceFullName, "usage_type", "SIGNING"),
			resource.TestCheckResourceAttr(resourceFullName, "validity_period", "365"),
			resource.TestMatchResourceAttr(resourceFullName, "rotated_at", regexp.MustCompile(`^(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d\.\d+)|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d)|(\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d)$`)),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironmentAndPEM(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckKeyRotationPolicyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			fullTest,
			{
				Config:  testAccKeyRotationPolicyConfig_Full(resourceName, nameFull),
				Destroy: true,
			},
			// Minimal
			minimalTest,
			{
				Config:  testAccKeyRotationPolicyConfig_Minimal(resourceName, nameMin),
				Destroy: true,
			},
			// Update
		},
	})
}

func testAccKeyRotationPolicyConfig_NewEnv(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_key_rotation_policy" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  name = "%[4]s"

  algorithm           = "RSA"
  subject_dn          = "CN=%[4]s, OU=Ping Identity, O=Ping Identity, L=, ST=, C=US"
  key_length          = 3072
  signature_algorithm = "SHA256withRSA"
  usage_type          = "SIGNING"
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccKeyRotationPolicyConfig_Full(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_key_rotation_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  algorithm           = "RSA"
  subject_dn          = "CN=%[3]s, OU=Ping Identity, O=Ping Identity, L=, ST=, C=US"
  key_length          = 3072
  signature_algorithm = "SHA256withRSA"
  usage_type          = "SIGNING"
  rotation_period     = 31
  validity_period     = 340
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccKeyRotationPolicyConfig_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_key_rotation_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  algorithm           = "RSA"
  subject_dn          = "CN=%[3]s, OU=Ping Identity, O=Ping Identity, L=, ST=, C=US"
  key_length          = 2048
  signature_algorithm = "SHA256withRSA"
  usage_type          = "SIGNING"
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
