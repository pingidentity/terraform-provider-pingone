package mfa_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func testAccCheckMFASettingsDestroy(s *terraform.State) error {
	return nil
}

func testAccGetMFASettingsIDs(resourceName string, environmentID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource Not found: %s", resourceName)
		}

		*environmentID = rs.Primary.Attributes["environment_id"]

		return nil
	}
}

func TestAccMFASettings_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_settings.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	var environmentID string

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckMFASettingsDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccMFASettingsConfig_Minimal(environmentName, licenseID, resourceName),
				Check:  testAccGetMFASettingsIDs(resourceFullName, &environmentID),
			},
			// Replan after removal preconfig
			{
				PreConfig: func() {
					var ctx = context.Background()
					p1Client, err := acctest.TestClient(ctx)

					if err != nil {
						t.Fatalf("Failed to get API client: %v", err)
					}

					apiClient := p1Client.API.MFAAPIClient

					if environmentID == "" {
						t.Fatalf("One of environment ID or resource ID cannot be determined. Environment ID: %s", environmentID)
					}

					_, _, err = apiClient.MFASettingsApi.ResetMFASettings(ctx, environmentID).Execute()
					if err != nil {
						t.Fatalf("Failed to reset MFA settings: %v", err)
					}
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccMFASettings_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_settings.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckMFASettingsDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFASettingsConfig_Full(environmentName, licenseID, resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "pairing.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "pairing.0.max_allowed_devices", "7"),
					resource.TestCheckResourceAttr(resourceFullName, "pairing.0.pairing_key_format", "NUMERIC"),
					resource.TestCheckResourceAttr(resourceFullName, "lockout.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "lockout.0.failure_count", "13"),
					resource.TestCheckResourceAttr(resourceFullName, "lockout.0.duration_seconds", "8"),
					resource.TestCheckResourceAttr(resourceFullName, "phone_extensions_enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "authentication.#", "1"),
					// resource.TestCheckResourceAttr(resourceFullName, "authentication.0.device_selection", "PROMPT_TO_SELECT"),
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

						return rs.Primary.ID, nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
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
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckMFASettingsDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFASettingsConfig_Minimal(environmentName, licenseID, resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "pairing.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "pairing.0.max_allowed_devices", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "pairing.0.pairing_key_format", "NUMERIC"),
					resource.TestCheckResourceAttr(resourceFullName, "lockout.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "phone_extensions_enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "authentication.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "authentication.0.device_selection", "DEFAULT_TO_FIRST"),
				),
			},
			{
				Config: testAccMFASettingsConfig_LockoutMinimal(environmentName, licenseID, resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "pairing.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "pairing.0.max_allowed_devices", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "pairing.0.pairing_key_format", "NUMERIC"),
					resource.TestCheckResourceAttr(resourceFullName, "lockout.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "lockout.0.failure_count", "13"),
					resource.TestCheckResourceAttr(resourceFullName, "lockout.0.duration_seconds", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "phone_extensions_enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "authentication.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "authentication.0.device_selection", "DEFAULT_TO_FIRST"),
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
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckMFASettingsDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFASettingsConfig_Full(environmentName, licenseID, resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "pairing.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "pairing.0.max_allowed_devices", "7"),
					resource.TestCheckResourceAttr(resourceFullName, "pairing.0.pairing_key_format", "NUMERIC"),
					resource.TestCheckResourceAttr(resourceFullName, "lockout.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "lockout.0.failure_count", "13"),
					resource.TestCheckResourceAttr(resourceFullName, "lockout.0.duration_seconds", "8"),
					resource.TestCheckResourceAttr(resourceFullName, "phone_extensions_enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "authentication.#", "1"),
					// resource.TestCheckResourceAttr(resourceFullName, "authentication.0.device_selection", "PROMPT_TO_SELECT"),
				),
			},
			{
				Config: testAccMFASettingsConfig_Minimal(environmentName, licenseID, resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "pairing.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "pairing.0.max_allowed_devices", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "pairing.0.pairing_key_format", "NUMERIC"),
					resource.TestCheckResourceAttr(resourceFullName, "lockout.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "phone_extensions_enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "authentication.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "authentication.0.device_selection", "DEFAULT_TO_FIRST"),
				),
			},
			{
				Config: testAccMFASettingsConfig_Full(environmentName, licenseID, resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "pairing.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "pairing.0.max_allowed_devices", "7"),
					resource.TestCheckResourceAttr(resourceFullName, "pairing.0.pairing_key_format", "NUMERIC"),
					resource.TestCheckResourceAttr(resourceFullName, "lockout.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "lockout.0.failure_count", "13"),
					resource.TestCheckResourceAttr(resourceFullName, "lockout.0.duration_seconds", "8"),
					resource.TestCheckResourceAttr(resourceFullName, "phone_extensions_enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "authentication.#", "1"),
					// resource.TestCheckResourceAttr(resourceFullName, "authentication.0.device_selection", "PROMPT_TO_SELECT"),
				),
			},
		},
	})
}
func TestAccMFASettings_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_settings.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckMFASettingsDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccMFASettingsConfig_Minimal(environmentName, licenseID, resourceName),
			},
			// Errors
			{
				ResourceName:  resourceFullName,
				ImportStateId: "/",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`Invalid import ID specified \(".*"\).  The ID should be in the format "environment_id" and must match regex: .*`),
			},
			{
				ResourceName:  resourceFullName,
				ImportStateId: "badformat/badformat",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`Invalid import ID specified \(".*"\).  The ID should be in the format "environment_id" and must match regex: .*`),
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

  phone_extensions_enabled = true

  //   authentication {
  //     device_selection = "PROMPT_TO_SELECT"
  //   }

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
    failure_count = 13
  }

}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName)
}
