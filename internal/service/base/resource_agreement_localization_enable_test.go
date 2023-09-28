package base_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func testAccCheckAgreementLocalizationEnableDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_agreement_localization_enable" {
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

		if !body.GetEnabled() {
			continue
		}

		return fmt.Errorf("PingOne agreement localization %s still exists", rs.Primary.ID)
	}

	return nil
}

func TestAccAgreementLocalizationEnable_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_agreement_localization_enable.%s", resourceName)

	name := resourceName

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
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckAgreementLocalizationEnableDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Enabled
			{
				Config: testAccAgreementLocalizationEnableConfig_Enable(resourceName, name),
				Check:  enabledCheck,
			},
			{
				Config:  testAccAgreementLocalizationEnableConfig_Enable(resourceName, name),
				Destroy: true,
			},
			// Disabled
			{
				Config: testAccAgreementLocalizationEnableConfig_Disable(resourceName, name),
				Check:  disabledCheck,
			},
			{
				Config:  testAccAgreementLocalizationEnableConfig_Disable(resourceName, name),
				Destroy: true,
			},
			// Change
			{
				Config: testAccAgreementLocalizationEnableConfig_Enable(resourceName, name),
				Check:  enabledCheck,
			},
			{
				Config: testAccAgreementLocalizationEnableConfig_Disable(resourceName, name),
				Check:  disabledCheck,
			},
			{
				Config: testAccAgreementLocalizationEnableConfig_Enable(resourceName, name),
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

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckAgreementLocalizationEnableDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccAgreementLocalizationEnableConfig_Enable(resourceName, name),
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

func testAccAgreementLocalizationEnableConfig_Enable(resourceName, name string) string {
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

`, acctest.AgreementSandboxEnvironment(), resourceName, name)
}

func testAccAgreementLocalizationEnableConfig_Disable(resourceName, name string) string {
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
`, acctest.AgreementSandboxEnvironment(), resourceName, name)
}
