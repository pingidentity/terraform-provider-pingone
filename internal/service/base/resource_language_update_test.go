// Copyright © 2025 Ping Identity Corporation

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
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccLanguageUpdate_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_language_update.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	var languageID, environmentID string

	var p1Client *client.Client
	var ctx = context.Background()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
			acctest.PreCheckTestAccFlaky(t)

			acctest.PreCheckNoBeta(t)
			p1Client = acctest.PreCheckTestClient(ctx, t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.LanguageUpdate_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Test removal of the environment
			{
				Config: testAccLanguageUpdateConfig_Full(environmentName, licenseID, resourceName, "de-DE", true),
				Check:  base.LanguageUpdate_GetIDs(resourceFullName, &environmentID, &languageID),
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

func TestAccLanguageUpdate_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_language_update.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.LanguageUpdate_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLanguageUpdateConfig_Full(environmentName, licenseID, resourceName, "de-DE", true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", "German (Germany)"),
					resource.TestCheckResourceAttr(resourceFullName, "locale", "de-DE"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "default", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "customer_added", "true"),
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

func TestAccLanguageUpdate_Minimal(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_language_update.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.LanguageUpdate_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLanguageUpdateConfig_Full(environmentName, licenseID, resourceName, "fr-FR", true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
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
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.LanguageUpdate_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLanguageUpdateConfig_SystemDefined(environmentName, licenseID, resourceName, "fr", true, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
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
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
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
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", "French"),
					resource.TestCheckResourceAttr(resourceFullName, "locale", "fr"),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "default", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "customer_added", "false"),
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

func TestAccLanguageUpdate_Change(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_language_update.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.LanguageUpdate_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLanguageUpdateConfig_Minimal(environmentName, licenseID, resourceName, "fr-FR", true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
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
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
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
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
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
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
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

func TestAccLanguageUpdate_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_language_update.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.LanguageUpdate_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccLanguageUpdateConfig_Full(environmentName, licenseID, resourceName, "de-DE", true),
			},
			// Errors
			{
				ResourceName: resourceFullName,
				ImportState:  true,
				ExpectError:  regexp.MustCompile(`Invalid import ID specified \(".*"\).  The ID should be in the format "environment_id/language_id" and must match regex: .*`),
			},
			{
				ResourceName:  resourceFullName,
				ImportStateId: "/",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`Invalid import ID specified \(".*"\).  The ID should be in the format "environment_id/language_id" and must match regex: .*`),
			},
			{
				ResourceName:  resourceFullName,
				ImportStateId: "badformat/badformat",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`Invalid import ID specified \(".*"\).  The ID should be in the format "environment_id/language_id" and must match regex: .*`),
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
