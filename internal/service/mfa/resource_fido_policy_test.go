package mfa_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/mfa"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccFIDOPolicy_NewEnv(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_fido_policy.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		//PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		PreCheck:                 func() { t.Skipf("Resource deprecated for new environments") },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.FIDOPolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFIDOPolicyConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
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
		//PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		PreCheck:                 func() { t.Skipf("Resource deprecated for new environments") },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.FIDOPolicy_CheckDestroy,
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
		//PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		PreCheck:                 func() { t.Skipf("Resource deprecated for new environments") },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.FIDOPolicy_CheckDestroy,
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
		//PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		PreCheck:                 func() { t.Skipf("Resource deprecated for new environments") },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.FIDOPolicy_CheckDestroy,
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
