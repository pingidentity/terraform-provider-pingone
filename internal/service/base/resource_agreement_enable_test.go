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

func TestAccAgreementEnable_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_agreement_enable.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	var agreementID, environmentID string

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
		CheckDestroy:             base.AgreementEnable_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Test removal of the agreement
			{
				Config: testAccAgreementEnableConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check:  base.AgreementEnable_GetIDs(resourceFullName, &environmentID, &agreementID),
			},
			{
				PreConfig: func() {
					base.Agreement_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID, agreementID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the environment
			{
				Config: testAccAgreementEnableConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check:  base.AgreementEnable_GetIDs(resourceFullName, &environmentID, &agreementID),
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

func TestAccAgreementEnable_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_agreement_enable.%s", resourceName)

	name := resourceName

	enabledCheck := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "agreement_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
	)

	disabledCheck := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "agreement_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "enabled", "false"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.AgreementEnable_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Enabled
			{
				Config: testAccAgreementEnableConfig_Enable(resourceName, name),
				Check:  enabledCheck,
			},
			{
				Config:  testAccAgreementEnableConfig_Enable(resourceName, name),
				Destroy: true,
			},
			// Disabled
			{
				Config: testAccAgreementEnableConfig_Disable(resourceName, name),
				Check:  disabledCheck,
			},
			{
				Config:  testAccAgreementEnableConfig_Disable(resourceName, name),
				Destroy: true,
			},
			// Change
			{
				Config: testAccAgreementEnableConfig_Enable(resourceName, name),
				Check:  enabledCheck,
			},
			{
				Config: testAccAgreementEnableConfig_Disable(resourceName, name),
				Check:  disabledCheck,
			},
			{
				Config: testAccAgreementEnableConfig_Enable(resourceName, name),
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

						return fmt.Sprintf("%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAgreementEnable_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_agreement_enable.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.AgreementEnable_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccAgreementEnableConfig_Enable(resourceName, name),
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

func testAccAgreementEnableConfig_NewEnv(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_language" "fr" {
  environment_id = pingone_environment.my_environment.id

  locale = "fr"
}

resource "pingone_language_update" "fr" {
  environment_id = pingone_environment.my_environment.id

  language_id = data.pingone_language.fr.id
  default     = true
  enabled     = true
}

resource "pingone_agreement" "my_agreement" {
  environment_id = pingone_environment.my_environment.id

  name        = "Terms and Conditions"
  description = "An agreement for general Terms and Conditions"
}

resource "pingone_agreement_localization" "my_agreement_fr" {
  environment_id = pingone_environment.my_environment.id
  agreement_id   = pingone_agreement.my_agreement.id
  language_id    = pingone_language_update.fr.id

  display_name = "Terms and Conditions - French Locale"
}

resource "time_static" "now" {}

resource "pingone_agreement_localization_revision" "my_agreement_fr_now" {
  environment_id            = pingone_environment.my_environment.id
  agreement_id              = pingone_agreement.my_agreement.id
  agreement_localization_id = pingone_agreement_localization.my_agreement_fr.id

  content_type      = "text/html"
  effective_at      = time_static.now.id
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

resource "pingone_agreement_localization_enable" "my_agreement_fr_enable" {
  environment_id            = pingone_environment.my_environment.id
  agreement_id              = pingone_agreement.my_agreement.id
  agreement_localization_id = pingone_agreement_localization.my_agreement_fr.id

  enabled = true

  depends_on = [
    pingone_agreement_localization_revision.my_agreement_fr_now
  ]
}

resource "pingone_agreement_enable" "my_agreement_enable" {
  environment_id = pingone_environment.my_environment.id
  agreement_id   = pingone_agreement.my_agreement.id

  enabled = true

  depends_on = [
    pingone_agreement_localization_enable.my_agreement_fr_enable
  ]
}
`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccAgreementEnableConfig_Enable(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_agreement" "%[2]s" {
  environment_id = data.pingone_environment.agreement_test.id

  name = "AgreementEnable"
}

resource "pingone_agreement_enable" "%[2]s" {
  environment_id = data.pingone_environment.agreement_test.id
  agreement_id   = data.pingone_agreement.%[2]s.id

  enabled = "true"
}
`, acctest.AgreementSandboxEnvironment(), resourceName, name)
}

func testAccAgreementEnableConfig_Disable(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_agreement" "%[2]s" {
  environment_id = data.pingone_environment.agreement_test.id

  name = "AgreementEnable"
}

resource "pingone_agreement_enable" "%[2]s" {
  environment_id = data.pingone_environment.agreement_test.id
  agreement_id   = data.pingone_agreement.%[2]s.id

  enabled = "false"
}
`, acctest.AgreementSandboxEnvironment(), resourceName, name)
}
