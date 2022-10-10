package mfa_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
)

func testAccCheckMFASettingsDestroy(s *terraform.State) error {
	return nil
}

func TestAccMFASettings_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_settings.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckMFASettingsDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFASettingsConfig_Full(environmentName, licenseID, resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "pairing.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "pairing.0.max_allowed_devices", "7"),
					resource.TestCheckResourceAttr(resourceFullName, "pairing.0.pairing_key_format", "NUMERIC"),
					resource.TestCheckResourceAttr(resourceFullName, "lockout.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "lockout.0.failure_count", "13"),
					resource.TestCheckResourceAttr(resourceFullName, "lockout.0.duration_seconds", "8"),
					resource.TestCheckResourceAttr(resourceFullName, "authentication.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "authentication.0.device_selection", "PROMPT_TO_SELECT"),
				),
			},
		},
	})
}

func TestAccMFASettings_Minimal(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_settings.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckMFASettingsDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFASettingsConfig_Minimal(environmentName, licenseID, resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "pairing.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "pairing.0.max_allowed_devices", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "pairing.0.pairing_key_format", "NUMERIC"),
					resource.TestCheckResourceAttr(resourceFullName, "lockout.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "authentication.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "authentication.0.device_selection", "PROMPT_TO_SELECT"),
				),
			},
			{
				Config: testAccMFASettingsConfig_LockoutMinimal(environmentName, licenseID, resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "pairing.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "pairing.0.max_allowed_devices", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "pairing.0.pairing_key_format", "NUMERIC"),
					resource.TestCheckResourceAttr(resourceFullName, "lockout.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "lockout.0.failure_count", "13"),
					resource.TestCheckResourceAttr(resourceFullName, "lockout.0.duration_seconds", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "authentication.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "authentication.0.device_selection", "PROMPT_TO_SELECT"),
				),
			},
		},
	})
}

func TestAccMFASettings_Change(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_settings.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckMFASettingsDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFASettingsConfig_Full(environmentName, licenseID, resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "pairing.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "pairing.0.max_allowed_devices", "7"),
					resource.TestCheckResourceAttr(resourceFullName, "pairing.0.pairing_key_format", "NUMERIC"),
					resource.TestCheckResourceAttr(resourceFullName, "lockout.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "lockout.0.failure_count", "13"),
					resource.TestCheckResourceAttr(resourceFullName, "lockout.0.duration_seconds", "8"),
					resource.TestCheckResourceAttr(resourceFullName, "authentication.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "authentication.0.device_selection", "PROMPT_TO_SELECT"),
				),
			},
			{
				Config: testAccMFASettingsConfig_Minimal(environmentName, licenseID, resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "pairing.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "pairing.0.max_allowed_devices", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "pairing.0.pairing_key_format", "NUMERIC"),
					resource.TestCheckResourceAttr(resourceFullName, "lockout.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "authentication.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "authentication.0.device_selection", "PROMPT_TO_SELECT"),
				),
			},
			{
				Config: testAccMFASettingsConfig_Full(environmentName, licenseID, resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "pairing.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "pairing.0.max_allowed_devices", "7"),
					resource.TestCheckResourceAttr(resourceFullName, "pairing.0.pairing_key_format", "NUMERIC"),
					resource.TestCheckResourceAttr(resourceFullName, "lockout.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "lockout.0.failure_count", "13"),
					resource.TestCheckResourceAttr(resourceFullName, "lockout.0.duration_seconds", "8"),
					resource.TestCheckResourceAttr(resourceFullName, "authentication.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "authentication.0.device_selection", "PROMPT_TO_SELECT"),
				),
			},
		},
	})
}

func testAccMFASettingsConfig_Full(environmentName, licenseID, resourceName string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_mfa_settings" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  pairing {
    max_allowed_devices = 7
    pairing_key_format  = "NUMERIC"
  }

  lockout {
    failure_count    = 13
    duration_seconds = 8
  }

  authentication {
    device_selection = "PROMPT_TO_SELECT"
  }

}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName)
}

func testAccMFASettingsConfig_Minimal(environmentName, licenseID, resourceName string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_mfa_settings" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  pairing {
    pairing_key_format = "NUMERIC"
  }

  authentication {
    device_selection = "PROMPT_TO_SELECT"
  }

}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName)
}

func testAccMFASettingsConfig_LockoutMinimal(environmentName, licenseID, resourceName string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_mfa_settings" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  pairing {
    pairing_key_format = "NUMERIC"
  }

  lockout {
    failure_count    = 13
  }

  authentication {
    device_selection = "PROMPT_TO_SELECT"
  }

}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName)
}
