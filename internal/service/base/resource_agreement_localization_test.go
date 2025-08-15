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
	acctestlegacysdk "github.com/pingidentity/terraform-provider-pingone/internal/acctest/legacysdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/base"
	baselegacysdk "github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/base/legacysdk"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccAgreementLocalization_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_agreement_localization.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	var agreementLocalizationID, agreementID, environmentID string

	var p1Client *client.Client
	var ctx = context.Background()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)

			p1Client = acctestlegacysdk.PreCheckTestClient(ctx, t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.AgreementLocalization_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Test removal of the resource
			{
				Config: testAccAgreementLocalizationConfig_Minimal(environmentName, licenseID, resourceName, name),
				Check:  base.AgreementLocalization_GetIDs(resourceFullName, &environmentID, &agreementID, &agreementLocalizationID),
			},
			{
				PreConfig: func() {
					base.AgreementLocalization_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID, agreementID, agreementLocalizationID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the agreement
			{
				Config: testAccAgreementLocalizationConfig_Minimal(environmentName, licenseID, resourceName, name),
				Check:  base.AgreementLocalization_GetIDs(resourceFullName, &environmentID, &agreementID, &agreementLocalizationID),
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
				Config: testAccAgreementLocalizationConfig_Minimal(environmentName, licenseID, resourceName, name),
				Check:  base.AgreementLocalization_GetIDs(resourceFullName, &environmentID, &agreementID, &agreementLocalizationID),
			},
			{
				PreConfig: func() {
					baselegacysdk.Environment_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAgreementLocalization_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_agreement_localization.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	nameFull := fmt.Sprintf("%s-full", resourceName)
	nameMin := fmt.Sprintf("%s-min", resourceName)

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	fullCheck := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "agreement_id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "language_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "display_name", nameFull),
		resource.TestCheckResourceAttr(resourceFullName, "locale", "en-GB"),
		resource.TestCheckResourceAttr(resourceFullName, "enabled", "false"),
		resource.TestCheckResourceAttr(resourceFullName, "text_checkbox_accept", "Yeah"),
		resource.TestCheckResourceAttr(resourceFullName, "text_button_continue", "Move on"),
		resource.TestCheckResourceAttr(resourceFullName, "text_button_decline", "Nah"),
	)

	minimalCheck := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "agreement_id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "language_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "display_name", nameMin),
		resource.TestCheckResourceAttr(resourceFullName, "locale", "en-GB"),
		resource.TestCheckResourceAttr(resourceFullName, "enabled", "false"),
		resource.TestCheckNoResourceAttr(resourceFullName, "text_checkbox_accept"),
		resource.TestCheckNoResourceAttr(resourceFullName, "text_button_continue"),
		resource.TestCheckNoResourceAttr(resourceFullName, "text_button_decline"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.AgreementLocalization_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			{
				Config: testAccAgreementLocalizationConfig_Full(environmentName, licenseID, resourceName, nameFull),
				Check:  fullCheck,
			},
			{
				Config:  testAccAgreementLocalizationConfig_Full(environmentName, licenseID, resourceName, nameFull),
				Destroy: true,
			},
			// Minimal
			{
				Config: testAccAgreementLocalizationConfig_Minimal(environmentName, licenseID, resourceName, nameMin),
				Check:  minimalCheck,
			},
			{
				Config:  testAccAgreementLocalizationConfig_Minimal(environmentName, licenseID, resourceName, nameMin),
				Destroy: true,
			},
			// Change
			{
				Config: testAccAgreementLocalizationConfig_Full(environmentName, licenseID, resourceName, nameFull),
				Check:  fullCheck,
			},
			{
				Config: testAccAgreementLocalizationConfig_Minimal(environmentName, licenseID, resourceName, nameMin),
				Check:  minimalCheck,
			},
			{
				Config: testAccAgreementLocalizationConfig_Full(environmentName, licenseID, resourceName, nameFull),
				Check:  fullCheck,
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

func TestAccAgreementLocalization_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_agreement_localization.%s", resourceName)

	name := resourceName

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.AgreementLocalization_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccAgreementLocalizationConfig_Minimal(environmentName, licenseID, resourceName, name),
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

func testAccAgreementLocalizationConfig_Full(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_language" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  locale = "en-GB"
}

resource "pingone_language_update" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  language_id = pingone_language.%[3]s.id
  enabled     = true
}

resource "pingone_agreement" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  name = "%[4]s"
}

resource "pingone_agreement_localization" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  agreement_id   = pingone_agreement.%[3]s.id
  language_id    = pingone_language_update.%[3]s.id

  display_name = "%[4]s"

  text_checkbox_accept = "Yeah"
  text_button_continue = "Move on"
  text_button_decline  = "Nah"
}
`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccAgreementLocalizationConfig_Minimal(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_language" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  locale = "en-GB"
}

resource "pingone_language_update" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  language_id = pingone_language.%[3]s.id
  enabled     = true
}

resource "pingone_agreement" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  name = "%[3]s"
}

resource "pingone_agreement_localization" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  agreement_id   = pingone_agreement.%[3]s.id
  language_id    = pingone_language_update.%[3]s.id

  display_name = "%[4]s"
}
`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}
