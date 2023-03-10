package base_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
)

func testAccCheckAgreementLocalizationEnableDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})

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
		resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
		resource.TestMatchResourceAttr(resourceFullName, "agreement_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
		resource.TestMatchResourceAttr(resourceFullName, "agreement_localization_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
		resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
	)

	disabledCheck := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
		resource.TestMatchResourceAttr(resourceFullName, "agreement_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
		resource.TestMatchResourceAttr(resourceFullName, "agreement_localization_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
		resource.TestCheckResourceAttr(resourceFullName, "enabled", "false"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
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
