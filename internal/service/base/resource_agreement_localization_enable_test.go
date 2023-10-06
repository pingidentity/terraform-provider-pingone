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

var (
	agreementLocalizationEnableExternalProviders = map[string]resource.ExternalProvider{
		"time": {
			Source:            "hashicorp/time",
			VersionConstraint: "0.9.1",
		},
	}
)

func TestAccAgreementLocalizationEnable_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_agreement_localization_enable.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	var agreementLocalizationID, agreementID, environmentID string

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
		CheckDestroy:             base.AgreementLocalizationEnable_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		ExternalProviders:        agreementEnableExternalProviders,
		Steps: []resource.TestStep{
			// Test removal of the agreement localization skipped - cannot delete the localization when it has an effective revision
			// Test removal of the agreement
			{
				Config: testAccAgreementLocalizationEnableConfig_NewEnv(environmentName, licenseID, resourceName, "fr"),
				Check:  base.AgreementLocalizationEnable_GetIDs(resourceFullName, &environmentID, &agreementID, &agreementLocalizationID),
			},
			{
				PreConfig: func() {
					base.Agreement_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID, agreementID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config:  testAccAgreementEnableConfig_NewEnv(environmentName, licenseID, resourceName, "fr"),
				Destroy: true,
			},
			// Test removal of the environment
			{
				SkipFunc: func() (bool, error) { return true, fmt.Errorf("TBC") },
				Config:   testAccAgreementLocalizationEnableConfig_NewEnv(environmentName, licenseID, resourceName, "fr"),
				Check:    base.AgreementLocalizationEnable_GetIDs(resourceFullName, &environmentID, &agreementID, &agreementLocalizationID),
			},
			{
				SkipFunc: func() (bool, error) { return true, fmt.Errorf("TBC") },
				PreConfig: func() {
					base.Environment_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAgreementLocalizationEnable_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_agreement_localization_enable.%s", resourceName)

	enabledCheck := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "agreement_id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "agreement_localization_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
	)

	disabledCheck := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "agreement_id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "agreement_localization_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "enabled", "false"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.AgreementLocalizationEnable_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		ExternalProviders:        agreementEnableExternalProviders,
		Steps: []resource.TestStep{
			// Enabled
			{
				Config: testAccAgreementLocalizationEnableConfig_Enable(resourceName),
				Check:  enabledCheck,
			},
			{
				Config:  testAccAgreementLocalizationEnableConfig_Enable(resourceName),
				Destroy: true,
			},
			// Disabled
			{
				Config: testAccAgreementLocalizationEnableConfig_Disable(resourceName),
				Check:  disabledCheck,
			},
			{
				Config:  testAccAgreementLocalizationEnableConfig_Disable(resourceName),
				Destroy: true,
			},
			// Change
			{
				Config: testAccAgreementLocalizationEnableConfig_Enable(resourceName),
				Check:  enabledCheck,
			},
			{
				Config: testAccAgreementLocalizationEnableConfig_Disable(resourceName),
				Check:  disabledCheck,
			},
			{
				Config: testAccAgreementLocalizationEnableConfig_Enable(resourceName),
				Check:  enabledCheck,
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

						return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["agreement_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAgreementLocalizationEnable_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_agreement_localization_enable.%s", resourceName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.AgreementLocalizationEnable_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		ExternalProviders:        agreementEnableExternalProviders,
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccAgreementLocalizationEnableConfig_Enable(resourceName),
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
				ImportStateId: "badformat/badformat/badformat",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`Unexpected Import Identifier`),
			},
		},
	})
}

func testAccAgreementLocalizationEnableConfig_NewEnv(environmentName, licenseID, resourceName, locale string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_language" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  locale = "%[4]s"
}

resource "pingone_language_update" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  language_id = data.pingone_language.%[3]s.id
  default     = false
  enabled     = true
}

resource "pingone_agreement" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  name        = "%[3]s"
  description = "An agreement for general Terms and Conditions"
}

resource "pingone_agreement_localization" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  agreement_id   = pingone_agreement.%[3]s.id
  language_id    = pingone_language_update.%[3]s.id

  display_name = "%[3]s"
}

resource "time_static" "%[3]s" {}

resource "pingone_agreement_localization_revision" "%[3]s" {
  environment_id            = pingone_environment.%[2]s.id
  agreement_id              = pingone_agreement.%[3]s.id
  agreement_localization_id = pingone_agreement_localization.%[3]s.id

  content_type      = "text/html"
  effective_at      = time_static.%[3]s.id
  require_reconsent = true
  text              = <<EOT
			<h1>Conditions de service</h1>
			
			Veuillez accepter les termes et conditions.
			
			<h2>Utilisation des données</h2>
			
			Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
			
			<h2>Soutien</h2>
			
			Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
			EOT
}

resource "pingone_agreement_localization_enable" "%[3]s" {
  environment_id            = pingone_environment.%[2]s.id
  agreement_id              = pingone_agreement.%[3]s.id
  agreement_localization_id = pingone_agreement_localization.%[3]s.id

  enabled = true

  depends_on = [
    pingone_agreement_localization_revision.%[3]s
  ]
}
`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, locale)
}

func testAccAgreementLocalizationEnableConfig_Enable(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_agreement" "%[2]s" {
  environment_id = data.pingone_environment.agreement_test.id

  name = "AgreementLocalizationEnable"
}

data "pingone_agreement_localization" "%[2]s" {
  environment_id = data.pingone_environment.agreement_test.id
  agreement_id   = data.pingone_agreement.%[2]s.id

  locale = "en"
}

resource "pingone_agreement_localization_enable" "%[2]s" {
  environment_id            = data.pingone_environment.agreement_test.id
  agreement_id              = data.pingone_agreement.%[2]s.id
  agreement_localization_id = data.pingone_agreement_localization.%[2]s.id

  enabled = "true"
}

`, acctest.AgreementSandboxEnvironment(), resourceName)
}

func testAccAgreementLocalizationEnableConfig_Disable(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_agreement" "%[2]s" {
  environment_id = data.pingone_environment.agreement_test.id

  name = "AgreementLocalizationEnable"
}

data "pingone_agreement_localization" "%[2]s" {
  environment_id = data.pingone_environment.agreement_test.id
  agreement_id   = data.pingone_agreement.%[2]s.id

  locale = "en"
}

resource "pingone_agreement_localization_enable" "%[2]s" {
  environment_id            = data.pingone_environment.agreement_test.id
  agreement_id              = data.pingone_agreement.%[2]s.id
  agreement_localization_id = data.pingone_agreement_localization.%[2]s.id

  enabled = "false"
}
`, acctest.AgreementSandboxEnvironment(), resourceName)
}
