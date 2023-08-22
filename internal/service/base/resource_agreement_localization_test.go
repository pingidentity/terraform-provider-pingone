package base_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func testAccCheckAgreementLocalizationDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_agreement_localization" {
			continue
		}

		_, rEnv, err := apiClient.EnvironmentsApi.ReadOneEnvironment(ctx, rs.Primary.Attributes["environment_id"]).Execute()

		if err != nil {

			if rEnv == nil {
				return fmt.Errorf("Response object does not exist and no error detected")
			}

			if rEnv.StatusCode == 404 {
				continue
			}

			return err
		}

		body, r, err := apiClient.AgreementLanguagesResourcesApi.ReadOneAgreementLanguage(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["agreement_id"], rs.Primary.ID).Execute()

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

		return fmt.Errorf("PingOne agreement localization %s still exists", rs.Primary.ID)
	}

	return nil
}

func testAccGetAgreementLocalizationIDs(resourceName string, environmentID, agreementID, resourceID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource Not found: %s", resourceName)
		}

		*resourceID = rs.Primary.ID
		*agreementID = rs.Primary.Attributes["agreement_id"]
		*environmentID = rs.Primary.Attributes["environment_id"]

		return nil
	}
}

func TestAccAgreementLocalization_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_agreement_localization.%s", resourceName)

	name := resourceName

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	var resourceID, agreementID, environmentID string

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckAgreementLocalizationDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccAgreementLocalizationConfig_Minimal(environmentName, licenseID, resourceName, name),
				Check:  testAccGetAgreementLocalizationIDs(resourceFullName, &environmentID, &agreementID, &resourceID),
			},
			// Replan after removal preconfig
			{
				PreConfig: func() {
					var ctx = context.Background()
					p1Client, err := acctest.TestClient(ctx)

					if err != nil {
						t.Fatalf("Failed to get API client: %v", err)
					}

					apiClient := p1Client.API.ManagementAPIClient

					if environmentID == "" || agreementID == "" || resourceID == "" {
						t.Fatalf("One of environment ID, agreement ID or resource ID cannot be determined. Environment ID: %s, Agreement ID: %s, Resource ID: %s", environmentID, agreementID, resourceID)
					}

					_, err = apiClient.AgreementLanguagesResourcesApi.DeleteAgreementLanguage(ctx, environmentID, agreementID, resourceID).Execute()
					if err != nil {
						t.Fatalf("Failed to delete agreement localisation: %v", err)
					}
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
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckAgreementLocalizationDestroy,
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
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckAgreementLocalizationDestroy,
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
				ExpectError:  regexp.MustCompile(`Invalid import ID specified \(".*"\).  The ID should be in the format "environment_id/agreement_id/agreement_localization_id".`),
			},
			{
				ResourceName:  resourceFullName,
				ImportStateId: "/",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`Invalid import ID specified \(".*"\).  The ID should be in the format "environment_id/agreement_id/agreement_localization_id".`),
			},
			{
				ResourceName:  resourceFullName,
				ImportStateId: "badformat/badformat/badformat",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`Invalid import ID specified \(".*"\).  The ID should be in the format "environment_id/agreement_id/agreement_localization_id".`),
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
`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
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
`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}
