package base_test

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

func testAccGetMFAPolicyIDs(resourceName string, environmentID, resourceID *string) resource.TestCheckFunc {
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

func TestAccMFAPolicy_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_policy.%s", resourceName)

	name := resourceName

	var resourceID, environmentID string

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckMFAPolicyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccMFAPolicyConfig_FullSMS(resourceName, name),
				Check:  testAccGetMFAPolicyIDs(resourceFullName, &environmentID, &resourceID),
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

					if environmentID == "" || resourceID == "" {
						t.Fatalf("One of environment ID or resource ID cannot be determined. Environment ID: %s, Resource ID: %s", environmentID, resourceID)
					}

					_, err = apiClient.DeviceAuthenticationPolicyApi.DeleteDeviceAuthenticationPolicy(ctx, environmentID, resourceID).Execute()
					if err != nil {
						t.Fatalf("Failed to delete MFA Policy: %v", err)
					}
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccLanguageUpdate_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_language_update.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckLanguageDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLanguageUpdateConfig_Full(environmentName, licenseID, resourceName, "de-DE", true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "name", "German (Germany)"),
					resource.TestCheckResourceAttr(resourceFullName, "locale", "de-DE"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "default", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "customer_added", "true"),
				),
			},
		},
	})
}

func TestAccLanguageUpdate_Minimal(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_language_update.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckLanguageDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLanguageUpdateConfig_Full(environmentName, licenseID, resourceName, "fr-FR", true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "name", "French (France)"),
					resource.TestCheckResourceAttr(resourceFullName, "locale", "fr-FR"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "default", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "customer_added", "true"),
				),
			},
		},
	})
}

func TestAccLanguageUpdate_SystemDefined(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_language_update.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckLanguageDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLanguageUpdateConfig_SystemDefined(environmentName, licenseID, resourceName, "fr", true, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "name", "French"),
					resource.TestCheckResourceAttr(resourceFullName, "locale", "fr"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "default", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "customer_added", "false"),
				),
			},
			{
				Config: testAccLanguageUpdateConfig_SystemDefined(environmentName, licenseID, resourceName, "fr", true, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "name", "French"),
					resource.TestCheckResourceAttr(resourceFullName, "locale", "fr"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "default", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "customer_added", "false"),
				),
			},
			{
				Config:      testAccLanguageUpdateConfig_SystemDefined(environmentName, licenseID, resourceName, "fr", false, true),
				ExpectError: regexp.MustCompile("Invalid combination of `enabled` and `default`, the language must be enabled to be made default. Got enabled=false, default=true."),
			},
			{
				Config: testAccLanguageUpdateConfig_SystemDefined(environmentName, licenseID, resourceName, "fr", false, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "name", "French"),
					resource.TestCheckResourceAttr(resourceFullName, "locale", "fr"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "default", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "customer_added", "false"),
				),
			},
		},
	})
}

func TestAccLanguageUpdate_Change(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_language_update.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckLanguageDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLanguageUpdateConfig_Minimal(environmentName, licenseID, resourceName, "fr-FR", true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "name", "French (France)"),
					resource.TestCheckResourceAttr(resourceFullName, "locale", "fr-FR"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "default", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "customer_added", "true"),
				),
			},
			{
				Config: testAccLanguageUpdateConfig_Minimal(environmentName, licenseID, resourceName, "fr-FR", false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "name", "French (France)"),
					resource.TestCheckResourceAttr(resourceFullName, "locale", "fr-FR"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "default", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "customer_added", "true"),
				),
			},
			{
				Config: testAccLanguageUpdateConfig_Full(environmentName, licenseID, resourceName, "fr-FR", true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "name", "French (France)"),
					resource.TestCheckResourceAttr(resourceFullName, "locale", "fr-FR"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "default", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "customer_added", "true"),
				),
			},
			{
				Config: testAccLanguageUpdateConfig_Minimal(environmentName, licenseID, resourceName, "fr-FR", false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "name", "French (France)"),
					resource.TestCheckResourceAttr(resourceFullName, "locale", "fr-FR"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "default", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "customer_added", "true"),
				),
			},
		},
	})
}

func testAccLanguageUpdateConfig_Full(environmentName, licenseID, resourceName, locale string, enabled bool) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_language" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  locale = "%[4]s"
}

resource "pingone_language_update" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  language_id = pingone_language.%[3]s.id
  enabled     = %[5]t
  default     = true
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, locale, enabled)
}

func testAccLanguageUpdateConfig_Minimal(environmentName, licenseID, resourceName, locale string, enabled bool) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_language" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  locale = "%[4]s"
}

resource "pingone_language_update" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  language_id = pingone_language.%[3]s.id
  enabled     = %[5]t
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, locale, enabled)
}

func testAccLanguageUpdateConfig_SystemDefined(environmentName, licenseID, resourceName, locale string, enabled, defaultValue bool) string {
	return fmt.Sprintf(`
		%[1]s

data "pingone_language" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  locale = "%[4]s"
}

resource "pingone_language_update" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  language_id = data.pingone_language.%[3]s.id
  enabled     = %[5]t
  default     = %[6]t
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, locale, enabled, defaultValue)
}
