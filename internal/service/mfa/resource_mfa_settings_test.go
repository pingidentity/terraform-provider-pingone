package mfa_test

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
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/mfa"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccMFASettings_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_settings.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	var environmentID string

	var p1Client *client.Client
	var ctx = context.Background()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)

			p1Client = acctest.PreCheckTestClient(ctx, t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.MFASettings_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccMFASettingsConfig_Minimal(environmentName, licenseID, resourceName),
				Check:  mfa.MFASettings_GetIDs(resourceFullName, &environmentID),
			},
			{
				PreConfig: func() {
					mfa.MFASettings_RemovalDrift_PreConfig(ctx, p1Client.API.MFAAPIClient, t, environmentID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the environment
			{
				Config: testAccMFASettingsConfig_Minimal(environmentName, licenseID, resourceName),
				Check:  mfa.MFASettings_GetIDs(resourceFullName, &environmentID),
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

func TestAccMFASettings_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_settings.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.MFASettings_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFASettingsConfig_Full(environmentName, licenseID, resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "pairing.max_allowed_devices", "7"),
					resource.TestCheckResourceAttr(resourceFullName, "pairing.pairing_key_format", "NUMERIC"),
					resource.TestCheckResourceAttr(resourceFullName, "lockout.failure_count", "13"),
					resource.TestCheckResourceAttr(resourceFullName, "lockout.duration_seconds", "8"),
					resource.TestCheckResourceAttr(resourceFullName, "phone_extensions.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "users.mfa_enabled", "true"),
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

						return rs.Primary.Attributes["environment_id"], nil
					}
				}(),
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "environment_id",
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
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.MFASettings_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFASettingsConfig_Minimal(environmentName, licenseID, resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "pairing.max_allowed_devices", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "pairing.pairing_key_format", "NUMERIC"),
					resource.TestCheckNoResourceAttr(resourceFullName, "lockout.failure_count"),
					resource.TestCheckNoResourceAttr(resourceFullName, "lockout.duration_seconds"),
					resource.TestCheckResourceAttr(resourceFullName, "phone_extensions.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "users.mfa_enabled", "false"),
				),
			},
			{
				Config: testAccMFASettingsConfig_LockoutMinimal(environmentName, licenseID, resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "pairing.max_allowed_devices", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "pairing.pairing_key_format", "NUMERIC"),
					resource.TestCheckResourceAttr(resourceFullName, "lockout.failure_count", "13"),
					resource.TestCheckNoResourceAttr(resourceFullName, "lockout.duration_seconds"),
					resource.TestCheckResourceAttr(resourceFullName, "phone_extensions.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "users.mfa_enabled", "false"),
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
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.MFASettings_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFASettingsConfig_Full(environmentName, licenseID, resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "pairing.max_allowed_devices", "7"),
					resource.TestCheckResourceAttr(resourceFullName, "pairing.pairing_key_format", "NUMERIC"),
					resource.TestCheckResourceAttr(resourceFullName, "lockout.failure_count", "13"),
					resource.TestCheckResourceAttr(resourceFullName, "lockout.duration_seconds", "8"),
					resource.TestCheckResourceAttr(resourceFullName, "phone_extensions.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "users.mfa_enabled", "true"),
				),
			},
			{
				Config: testAccMFASettingsConfig_Minimal(environmentName, licenseID, resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "pairing.max_allowed_devices", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "pairing.pairing_key_format", "NUMERIC"),
					resource.TestCheckNoResourceAttr(resourceFullName, "lockout.failure_count"),
					resource.TestCheckNoResourceAttr(resourceFullName, "lockout.duration_seconds"),
					resource.TestCheckResourceAttr(resourceFullName, "phone_extensions.enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "users.mfa_enabled", "false"),
				),
			},
			{
				Config: testAccMFASettingsConfig_Full(environmentName, licenseID, resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "pairing.max_allowed_devices", "7"),
					resource.TestCheckResourceAttr(resourceFullName, "pairing.pairing_key_format", "NUMERIC"),
					resource.TestCheckResourceAttr(resourceFullName, "lockout.failure_count", "13"),
					resource.TestCheckResourceAttr(resourceFullName, "lockout.duration_seconds", "8"),
					resource.TestCheckResourceAttr(resourceFullName, "phone_extensions.enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "users.mfa_enabled", "true"),
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
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.MFASettings_CheckDestroy,
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

func testAccMFASettingsConfig_Full(environmentName, licenseID, resourceName string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_mfa_settings" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  lockout = {
    failure_count    = 13
    duration_seconds = 8
  }

  pairing = {
    max_allowed_devices = 7
    pairing_key_format  = "NUMERIC"
  }

  phone_extensions = {
    enabled = true
  }

  users = {
    mfa_enabled = true
  }

}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName)
}

func testAccMFASettingsConfig_Minimal(environmentName, licenseID, resourceName string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_mfa_settings" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  pairing = {
    pairing_key_format = "NUMERIC"
  }
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName)
}

func testAccMFASettingsConfig_LockoutMinimal(environmentName, licenseID, resourceName string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_mfa_settings" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  pairing = {
    pairing_key_format = "NUMERIC"
  }

  lockout = {
    failure_count = 13
  }
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName)
}
