package mfa_test

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
	"github.com/patrickcping/pingone-go-sdk-v2/mfa"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
)

func testAccCheckFIDOPolicyDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.MFAAPIClient
	ctx = context.WithValue(ctx, mfa.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})

	apiClientManagement := p1Client.API.ManagementAPIClient
	ctxManagement := context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_mfa_fido_policy" {
			continue
		}

		_, rEnv, err := apiClientManagement.EnvironmentsApi.ReadOneEnvironment(ctxManagement, rs.Primary.Attributes["environment_id"]).Execute()

		if err != nil {

			if rEnv == nil {
				return fmt.Errorf("Response object does not exist and no error detected")
			}

			if rEnv.StatusCode == 404 {
				continue
			}

			return err
		}

		body, r, err := apiClient.FIDOPolicyApi.ReadOneFidoPolicy(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.ID).Execute()

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

		return fmt.Errorf("PingOne MFA FIDO Policy Instance %s still exists", rs.Primary.ID)
	}

	return nil
}

func TestAccFIDOPolicy_NewEnv(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_fido_policy.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckFIDOPolicyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFIDOPolicyConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
				),
			},
		},
	})
}

func TestAccFIDOPolicy_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_fido_policy.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckFIDOPolicyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFIDOPolicyConfig_Full(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", "Test FIDO policy"),
					resource.TestCheckResourceAttr(resourceFullName, "attestation_requirements", "CERTIFIED"),
					resource.TestCheckResourceAttr(resourceFullName, "allowed_authenticators.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "resident_key_requirement", "DISCOURAGED"),
					resource.TestCheckResourceAttr(resourceFullName, "enforce_during_authentication", "true"),
				),
			},
		},
	})
}

func TestAccFIDOPolicy_Minimal(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_fido_policy.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckFIDOPolicyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFIDOPolicyConfig_Minimal(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", ""),
					resource.TestCheckResourceAttr(resourceFullName, "attestation_requirements", "GLOBAL"),
					resource.TestCheckResourceAttr(resourceFullName, "allowed_authenticators.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "resident_key_requirement", "REQUIRED"),
					resource.TestCheckResourceAttr(resourceFullName, "enforce_during_authentication", "false"),
				),
			},
		},
	})
}

func TestAccFIDOPolicy_Change(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_fido_policy.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckFIDOPolicyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFIDOPolicyConfig_Full(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", "Test FIDO policy"),
					resource.TestCheckResourceAttr(resourceFullName, "attestation_requirements", "CERTIFIED"),
					resource.TestCheckResourceAttr(resourceFullName, "allowed_authenticators.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "resident_key_requirement", "DISCOURAGED"),
					resource.TestCheckResourceAttr(resourceFullName, "enforce_during_authentication", "true"),
				),
			},
			{
				Config: testAccFIDOPolicyConfig_Minimal(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", ""),
					resource.TestCheckResourceAttr(resourceFullName, "attestation_requirements", "GLOBAL"),
					resource.TestCheckResourceAttr(resourceFullName, "allowed_authenticators.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "resident_key_requirement", "REQUIRED"),
					resource.TestCheckResourceAttr(resourceFullName, "enforce_during_authentication", "false"),
				),
			},
			{
				Config: testAccFIDOPolicyConfig_Full(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", "Test FIDO policy"),
					resource.TestCheckResourceAttr(resourceFullName, "attestation_requirements", "CERTIFIED"),
					resource.TestCheckResourceAttr(resourceFullName, "allowed_authenticators.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "resident_key_requirement", "DISCOURAGED"),
					resource.TestCheckResourceAttr(resourceFullName, "enforce_during_authentication", "true"),
				),
			},
		},
	})
}

func testAccFIDOPolicyConfig_NewEnv(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_mfa_fido_policy" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s"

  attestation_requirements = "AUDIT_ONLY"
  resident_key_requirement = "DISCOURAGED"

}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccFIDOPolicyConfig_Full(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_mfa_fido_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Test FIDO policy"

  attestation_requirements = "CERTIFIED"
  resident_key_requirement = "DISCOURAGED"

  enforce_during_authentication = true

  // allowed_authenticators

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccFIDOPolicyConfig_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_mfa_fido_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  attestation_requirements = "GLOBAL"
  resident_key_requirement = "REQUIRED"

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
