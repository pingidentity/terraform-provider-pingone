package credentials_test

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
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/credentials"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccCredentialIssuanceRule_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_credential_issuance_rule.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	var credentialIssuanceRuleID, digitalWalletApplicationID, credentialTypeID, environmentID string

	var ctx = context.Background()
	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		t.Fatalf("Failed to get API client: %v", err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             credentials.TestAccCheckCredentialIssuanceRuleDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Test removal of the resource
			{
				Config: testAccCredentialIssuanceRuleConfig_Full(resourceName, name),
				Check:  credentials.TestAccGetCredentialIssuanceRuleIDs(resourceFullName, &environmentID, &credentialTypeID, &digitalWalletApplicationID, &credentialIssuanceRuleID),
			},
			// Replan after removal preconfig
			{
				PreConfig: func() {
					credentials.CredentialIssuanceRule_RemovalDrift_PreConfig(ctx, p1Client.API.CredentialsAPIClient, t, environmentID, credentialTypeID, credentialIssuanceRuleID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the credential type
			{
				Config: testAccCredentialIssuanceRuleConfig_Full(resourceName, name),
				Check:  credentials.TestAccGetCredentialIssuanceRuleIDs(resourceFullName, &environmentID, &credentialTypeID, &digitalWalletApplicationID, &credentialIssuanceRuleID),
			},
			// Replan after removal preconfig
			{
				PreConfig: func() {
					credentials.CredentialType_RemovalDrift_PreConfig(ctx, p1Client.API.CredentialsAPIClient, t, environmentID, credentialTypeID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the digital wallet application
			{
				Config: testAccCredentialIssuanceRuleConfig_Full(resourceName, name),
				Check:  credentials.TestAccGetCredentialIssuanceRuleIDs(resourceFullName, &environmentID, &credentialTypeID, &digitalWalletApplicationID, &credentialIssuanceRuleID),
			},
			// Replan after removal preconfig
			{
				PreConfig: func() {
					credentials.DigitalWalletApplication_RemovalDrift_PreConfig(ctx, p1Client.API.CredentialsAPIClient, t, environmentID, credentialTypeID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the environment
			{
				Config: testAccCredentialIssuanceRuleConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check:  credentials.TestAccGetCredentialIssuanceRuleIDs(resourceFullName, &environmentID, &credentialTypeID, &digitalWalletApplicationID, &credentialIssuanceRuleID),
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

func TestAccCredentialIssuanceRule_NewEnv(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_credential_issuance_rule.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             credentials.TestAccCheckCredentialIssuanceRuleDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCredentialIssuanceRuleConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
				),
			},
		},
	})
}

func TestAccCredentialIssuanceRule_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_credential_issuance_rule.%s", resourceName)

	name := acctest.ResourceNameGen()

	fullStep := resource.TestStep{
		Config: testAccCredentialIssuanceRuleConfig_Full(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "credential_type_id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "digital_wallet_application_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "automation.issue", "ON_DEMAND"),
			resource.TestCheckResourceAttr(resourceFullName, "automation.revoke", "PERIODIC"),
			resource.TestCheckResourceAttr(resourceFullName, "automation.update", "ON_DEMAND"),
			resource.TestCheckResourceAttrSet(resourceFullName, "filter.population_ids.#"),
			resource.TestCheckResourceAttr(resourceFullName, "notification.methods.#", "2"),
			resource.TestCheckResourceAttr(resourceFullName, "notification.methods.0", "EMAIL"),
			resource.TestCheckResourceAttr(resourceFullName, "notification.methods.1", "SMS"),
			resource.TestCheckResourceAttr(resourceFullName, "notification.template.locale", "en"),
			resource.TestCheckResourceAttr(resourceFullName, "notification.template.variant", "template_B"),
			resource.TestCheckResourceAttr(resourceFullName, "status", "ACTIVE"),
		),
	}

	minimalStep := resource.TestStep{
		Config: testAccCredentialIssuanceRuleConfig_Minimal(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "credential_type_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckNoResourceAttr(resourceFullName, "digital_wallet_application_id"),
			resource.TestCheckResourceAttr(resourceFullName, "filter.scim", "address.countryCode eq \"NG\""),
			resource.TestCheckResourceAttr(resourceFullName, "automation.issue", "PERIODIC"),
			resource.TestCheckResourceAttr(resourceFullName, "automation.revoke", "PERIODIC"),
			resource.TestCheckResourceAttr(resourceFullName, "automation.update", "PERIODIC"),
			resource.TestCheckResourceAttr(resourceFullName, "status", "ACTIVE"),
		),
	}

	disabledStep := resource.TestStep{
		Config: testAccCredentialIssuanceRuleConfig_Disabled(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "credential_type_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckNoResourceAttr(resourceFullName, "digital_wallet_application_id"),
			resource.TestCheckResourceAttr(resourceFullName, "automation.%", "3"),
			resource.TestCheckResourceAttr(resourceFullName, "automation.issue", "PERIODIC"),
			resource.TestCheckResourceAttr(resourceFullName, "automation.revoke", "PERIODIC"),
			resource.TestCheckResourceAttr(resourceFullName, "automation.update", "PERIODIC"),
			resource.TestCheckResourceAttr(resourceFullName, "filter.%", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "notification.#", "0"),
			resource.TestCheckNoResourceAttr(resourceFullName, "notification.methods"),
			resource.TestCheckNoResourceAttr(resourceFullName, "notification.template"),
			resource.TestCheckResourceAttr(resourceFullName, "status", "DISABLED"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             credentials.TestAccCheckCredentialIssuanceRuleDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// full
			fullStep,
			{
				Config:  testAccCredentialIssuanceRuleConfig_Full(resourceName, name),
				Destroy: true,
			},
			//minimalStep,
			{
				Config:  testAccCredentialIssuanceRuleConfig_Minimal(resourceName, name),
				Destroy: true,
			},
			disabledStep,
			{
				Config:  testAccCredentialIssuanceRuleConfig_Disabled(resourceName, name),
				Destroy: true,
			},
			fullStep,
			minimalStep,
			fullStep,
			disabledStep,
			fullStep,
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["credential_type_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCredentialIssuanceRule_InvalidConfigs(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	name := acctest.ResourceNameGen()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             credentials.TestAccCheckCredentialIssuanceRuleDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccCredentialIssuanceRuleConfig_InvalidCredentialTypeID(resourceName, name),
				ExpectError: regexp.MustCompile("Error: Missing required argument"),
				Destroy:     true,
			},
			{
				Config:      testAccCredentialIssuanceRuleConfig_InvalidGroupIdFilterAttribute(resourceName, name),
				ExpectError: regexp.MustCompile("Error: Invalid Attribute Value"),
				Destroy:     true,
			},
			{
				Config:      testAccCredentialIssuanceRuleConfig_InvalidPopulationIdFilterAttribute(resourceName, name),
				ExpectError: regexp.MustCompile("Error: Invalid Attribute Value"),
				Destroy:     true,
			},
			{
				Config:      testAccCredentialIssuanceRuleConfig_InvalidFilterAttribute(resourceName, name),
				ExpectError: regexp.MustCompile("Error: Invalid Attribute Combination"),
				Destroy:     true,
			},
			{
				Config:      testAccCredentialIssuanceRuleConfig_InvalidAutomationAttribute(resourceName, name),
				ExpectError: regexp.MustCompile("Error: Incorrect attribute value type"),
				Destroy:     true,
			},
			{
				Config:      testAccCredentialIssuanceRuleConfig_InvalidNotificationAttribute(resourceName, name),
				ExpectError: regexp.MustCompile("Error: Incorrect attribute value type"),
				Destroy:     true,
			},
			{
				Config:      testAccCredentialIssuanceRuleConfig_InvalidNotificationMethodsAttribute(resourceName, name),
				ExpectError: regexp.MustCompile("Error: Invalid Attribute Value"),
				Destroy:     true,
			},
			{
				Config:      testAccCredentialIssuanceRuleConfig_InvalidNotificationTemplateAttribute(resourceName, name),
				ExpectError: regexp.MustCompile("Error: Incorrect attribute value type"),
				Destroy:     true,
			},
		},
	})
}

func TestAccCredentialIssuanceRule_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_credential_issuance_rule.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             credentials.TestAccCheckCredentialIssuanceRuleDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccCredentialIssuanceRuleConfig_Minimal(resourceName, name),
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

func testAccCredentialIssuanceRuleConfig_NewEnv(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_population" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s"
}

resource "pingone_credential_type" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  title                = "%[4]s"
  description          = "%[4]s Example Description"
  card_type            = "%[4]s"
  card_design_template = "<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 740 480\"><rect fill=\"none\" width=\"736\" height=\"476\" stroke=\"#CACED3\" stroke-width=\"3\" rx=\"10\" ry=\"10\" x=\"2\" y=\"2\"></rect><rect fill=\"$${cardColor}\" height=\"476\" rx=\"10\" ry=\"10\" width=\"736\" x=\"2\" y=\"2\" opacity=\"$${bgOpacityPercent}\"></rect><line y2=\"160\" x2=\"695\" y1=\"160\" x1=\"42.5\" stroke=\"$${textColor}\"></line><text fill=\"$${textColor}\" font-weight=\"450\" font-size=\"30\" x=\"160\" y=\"90\">$${cardTitle}</text><text fill=\"$${textColor}\" font-size=\"25\" font-weight=\"300\" x=\"160\" y=\"130\">$${cardSubtitle}</text></svg>"

  metadata = {
    name               = "%[4]s"
    description        = "%[4]s Example Description"
    bg_opacity_percent = 100
    card_color         = "#000000"
    text_color         = "#eff0f1"

    fields = [
      {
        type       = "Alphanumeric Text"
        title      = "Example Field"
        value      = "Demo"
        is_visible = false
      },
    ]
  }
}

resource "pingone_application" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s"
  enabled        = true

  oidc_options {
    type                        = "NATIVE_APP"
    grant_types                 = ["CLIENT_CREDENTIALS"]
    token_endpoint_authn_method = "CLIENT_SECRET_BASIC"

    mobile_app {
      bundle_id                = "com.pingidentity.ios_%[4]s"
      package_name             = "com.pingidentity.android_%[4]s"
      passcode_refresh_seconds = 30
    }
  }
}

resource "pingone_digital_wallet_application" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  application_id = resource.pingone_application.%[3]s.id
  name           = "%[4]s"
  app_open_url   = "https://www.example.com"

  depends_on = [resource.pingone_application.%[3]s]
}

resource "pingone_credential_issuance_rule" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  credential_type_id            = resource.pingone_credential_type.%[3]s.id
  digital_wallet_application_id = resource.pingone_digital_wallet_application.%[3]s.id
  status                        = "ACTIVE"

  filter = {
    population_ids = [resource.pingone_population.%[3]s.id]
  }

  automation = {
    issue  = "ON_DEMAND"
    revoke = "PERIODIC"
    update = "ON_DEMAND"
  }

  notification = {
    methods = ["EMAIL", "SMS"]
    template = {
      locale  = "en"
      variant = "template_B"
    }
  }
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccCredentialIssuanceRuleConfig_Full(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_population" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
}

resource "pingone_credential_type" "%[2]s" {
  environment_id       = data.pingone_environment.general_test.id
  title                = "%[3]s"
  description          = "%[3]s Example Description"
  card_type            = "%[3]s"
  card_design_template = "<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 740 480\"><rect fill=\"none\" width=\"736\" height=\"476\" stroke=\"#CACED3\" stroke-width=\"3\" rx=\"10\" ry=\"10\" x=\"2\" y=\"2\"></rect><rect fill=\"$${cardColor}\" height=\"476\" rx=\"10\" ry=\"10\" width=\"736\" x=\"2\" y=\"2\" opacity=\"$${bgOpacityPercent}\"></rect><line y2=\"160\" x2=\"695\" y1=\"160\" x1=\"42.5\" stroke=\"$${textColor}\"></line><text fill=\"$${textColor}\" font-weight=\"450\" font-size=\"30\" x=\"160\" y=\"90\">$${cardTitle}</text><text fill=\"$${textColor}\" font-size=\"25\" font-weight=\"300\" x=\"160\" y=\"130\">$${cardSubtitle}</text></svg>"

  metadata = {
    name               = "%[3]s"
    description        = "%[3]s Example Description"
    bg_opacity_percent = 100
    card_color         = "#000000"
    text_color         = "#eff0f1"

    fields = [
      {
        type       = "Alphanumeric Text"
        title      = "Example Field"
        value      = "Demo"
        is_visible = false
      },
    ]
  }
}

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true

  oidc_options {
    type                        = "NATIVE_APP"
    grant_types                 = ["CLIENT_CREDENTIALS"]
    token_endpoint_authn_method = "CLIENT_SECRET_BASIC"

    mobile_app {
      bundle_id                = "com.pingidentity.ios_%[3]s"
      package_name             = "com.pingidentity.android_%[3]s"
      passcode_refresh_seconds = 30
    }
  }
}

resource "pingone_digital_wallet_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  application_id = resource.pingone_application.%[2]s.id
  name           = "%[3]s"
  app_open_url   = "https://www.example.com"

  depends_on = [resource.pingone_application.%[2]s]
}

resource "pingone_credential_issuance_rule" "%[2]s" {
  environment_id                = data.pingone_environment.general_test.id
  credential_type_id            = resource.pingone_credential_type.%[2]s.id
  digital_wallet_application_id = resource.pingone_digital_wallet_application.%[2]s.id
  status                        = "ACTIVE"

  filter = {
    population_ids = [resource.pingone_population.%[2]s.id]
  }

  automation = {
    issue  = "ON_DEMAND"
    revoke = "PERIODIC"
    update = "ON_DEMAND"
  }

  notification = {
    methods = ["EMAIL", "SMS"]
    template = {
      locale  = "en"
      variant = "template_B"
    }
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccCredentialIssuanceRuleConfig_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_credential_type" "%[2]s" {
  environment_id       = data.pingone_environment.general_test.id
  title                = "%[3]s"
  description          = "%[3]s Example Description"
  card_type            = "%[3]s"
  card_design_template = "<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 740 480\"><rect fill=\"none\" width=\"736\" height=\"476\" stroke=\"#CACED3\" stroke-width=\"3\" rx=\"10\" ry=\"10\" x=\"2\" y=\"2\"></rect><rect fill=\"$${cardColor}\" height=\"476\" rx=\"10\" ry=\"10\" width=\"736\" x=\"2\" y=\"2\" opacity=\"$${bgOpacityPercent}\"></rect><line y2=\"160\" x2=\"695\" y1=\"160\" x1=\"42.5\" stroke=\"$${textColor}\"></line><text fill=\"$${textColor}\" font-weight=\"450\" font-size=\"30\" x=\"160\" y=\"90\">$${cardTitle}</text><text fill=\"$${textColor}\" font-size=\"25\" font-weight=\"300\" x=\"160\" y=\"130\">$${cardSubtitle}</text></svg>"

  metadata = {
    name               = "%[3]s"
    description        = "%[3]s Example Description"
    bg_opacity_percent = 100
    card_color         = "#000000"
    text_color         = "#eff0f1"

    fields = [
      {
        type       = "Alphanumeric Text"
        title      = "Example Field"
        value      = "Demo"
        is_visible = false
      },
    ]
  }
}

resource "pingone_credential_issuance_rule" "%[2]s" {
  environment_id     = data.pingone_environment.general_test.id
  credential_type_id = resource.pingone_credential_type.%[2]s.id
  status             = "ACTIVE"

  filter = {
    scim = "address.countryCode eq \"NG\""
  }

  automation = {
    issue  = "PERIODIC"
    revoke = "PERIODIC"
    update = "PERIODIC"
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccCredentialIssuanceRuleConfig_Disabled(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_credential_type" "%[2]s" {
  environment_id       = data.pingone_environment.general_test.id
  title                = "%[3]s"
  description          = "%[3]s Example Description"
  card_type            = "%[3]s"
  card_design_template = "<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 740 480\"><rect fill=\"none\" width=\"736\" height=\"476\" stroke=\"#CACED3\" stroke-width=\"3\" rx=\"10\" ry=\"10\" x=\"2\" y=\"2\"></rect><rect fill=\"$${cardColor}\" height=\"476\" rx=\"10\" ry=\"10\" width=\"736\" x=\"2\" y=\"2\" opacity=\"$${bgOpacityPercent}\"></rect><line y2=\"160\" x2=\"695\" y1=\"160\" x1=\"42.5\" stroke=\"$${textColor}\"></line><text fill=\"$${textColor}\" font-weight=\"450\" font-size=\"30\" x=\"160\" y=\"90\">$${cardTitle}</text><text fill=\"$${textColor}\" font-size=\"25\" font-weight=\"300\" x=\"160\" y=\"130\">$${cardSubtitle}</text></svg>"

  metadata = {
    name               = "%[3]s"
    description        = "%[3]s Example Description"
    bg_opacity_percent = 100
    card_color         = "#000000"
    text_color         = "#eff0f1"

    fields = [
      {
        type       = "Alphanumeric Text"
        title      = "Example Field"
        value      = "Demo"
        is_visible = false
      },
    ]
  }
}

resource "pingone_credential_issuance_rule" "%[2]s" {
  environment_id     = data.pingone_environment.general_test.id
  credential_type_id = resource.pingone_credential_type.%[2]s.id
  status             = "DISABLED"

  automation = {
    issue  = "PERIODIC"
    revoke = "PERIODIC"
    update = "PERIODIC"
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccCredentialIssuanceRuleConfig_InvalidCredentialTypeID(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s


resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true

  oidc_options {
    type                        = "NATIVE_APP"
    grant_types                 = ["CLIENT_CREDENTIALS"]
    token_endpoint_authn_method = "CLIENT_SECRET_BASIC"

    mobile_app {
      bundle_id                = "com.pingidentity.ios_%[3]s"
      package_name             = "com.pingidentity.android_%[3]s"
      passcode_refresh_seconds = 30
    }
  }
}

resource "pingone_digital_wallet_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  application_id = resource.pingone_application.%[2]s.id
  name           = "%[3]s"
  app_open_url   = "https://www.example.com"

  depends_on = [resource.pingone_application.%[2]s]
}

resource "pingone_credential_issuance_rule" "%[2]s" {
  environment_id                = data.pingone_environment.general_test.id
  digital_wallet_application_id = resource.pingone_digital_wallet_application.%[2]s.id
  status                        = "DISABLED"

  filter = {
    scim = "address.countryCode eq \"CA\""
  }

  automation = {
    issue  = "PERIODIC"
    revoke = "PERIODIC"
    update = "PERIODIC"
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccCredentialIssuanceRuleConfig_InvalidGroupIdFilterAttribute(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_population" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
}

resource "pingone_credential_type" "%[2]s" {
  environment_id       = data.pingone_environment.general_test.id
  title                = "%[3]s"
  description          = "%[3]s Example Description"
  card_type            = "%[3]s"
  card_design_template = "<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 740 480\"><rect fill=\"none\" width=\"736\" height=\"476\" stroke=\"#CACED3\" stroke-width=\"3\" rx=\"10\" ry=\"10\" x=\"2\" y=\"2\"></rect><rect fill=\"$${cardColor}\" height=\"476\" rx=\"10\" ry=\"10\" width=\"736\" x=\"2\" y=\"2\" opacity=\"$${bgOpacityPercent}\"></rect><line y2=\"160\" x2=\"695\" y1=\"160\" x1=\"42.5\" stroke=\"$${textColor}\"></line><text fill=\"$${textColor}\" font-weight=\"450\" font-size=\"30\" x=\"160\" y=\"90\">$${cardTitle}</text><text fill=\"$${textColor}\" font-size=\"25\" font-weight=\"300\" x=\"160\" y=\"130\">$${cardSubtitle}</text></svg>"

  metadata = {
    name               = "%[3]s"
    description        = "%[3]s Example Description"
    bg_opacity_percent = 100
    card_color         = "#000000"
    text_color         = "#eff0f1"

    fields = [
      {
        type       = "Alphanumeric Text"
        title      = "Example Field"
        value      = "Demo"
        is_visible = false
      },
    ]
  }
}

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true

  oidc_options {
    type                        = "NATIVE_APP"
    grant_types                 = ["CLIENT_CREDENTIALS"]
    token_endpoint_authn_method = "CLIENT_SECRET_BASIC"

    mobile_app {
      bundle_id                = "com.pingidentity.ios_%[3]s"
      package_name             = "com.pingidentity.android_%[3]s"
      passcode_refresh_seconds = 30
    }
  }
}

resource "pingone_digital_wallet_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  application_id = resource.pingone_application.%[2]s.id
  name           = "%[3]s"
  app_open_url   = "https://www.example.com"

  depends_on = [resource.pingone_application.%[2]s]
}

resource "pingone_credential_issuance_rule" "%[2]s" {
  environment_id                = data.pingone_environment.general_test.id
  credential_type_id            = resource.pingone_credential_type.%[2]s.id
  digital_wallet_application_id = resource.pingone_digital_wallet_application.%[2]s.id
  status                        = "ACTIVE"

  filter = {
    group_ids = []
  }

  automation = {
    issue  = "ON_DEMAND"
    revoke = "PERIODIC"
    update = "ON_DEMAND"
  }

  notification = {
    methods = ["EMAIL", "SMS"]
    template = {
      locale  = "en"
      variant = "template_B"
    }
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccCredentialIssuanceRuleConfig_InvalidPopulationIdFilterAttribute(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_population" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
}

resource "pingone_credential_type" "%[2]s" {
  environment_id       = data.pingone_environment.general_test.id
  title                = "%[3]s"
  description          = "%[3]s Example Description"
  card_type            = "%[3]s"
  card_design_template = "<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 740 480\"><rect fill=\"none\" width=\"736\" height=\"476\" stroke=\"#CACED3\" stroke-width=\"3\" rx=\"10\" ry=\"10\" x=\"2\" y=\"2\"></rect><rect fill=\"$${cardColor}\" height=\"476\" rx=\"10\" ry=\"10\" width=\"736\" x=\"2\" y=\"2\" opacity=\"$${bgOpacityPercent}\"></rect><line y2=\"160\" x2=\"695\" y1=\"160\" x1=\"42.5\" stroke=\"$${textColor}\"></line><text fill=\"$${textColor}\" font-weight=\"450\" font-size=\"30\" x=\"160\" y=\"90\">$${cardTitle}</text><text fill=\"$${textColor}\" font-size=\"25\" font-weight=\"300\" x=\"160\" y=\"130\">$${cardSubtitle}</text></svg>"

  metadata = {
    name               = "%[3]s"
    description        = "%[3]s Example Description"
    bg_opacity_percent = 100
    card_color         = "#000000"
    text_color         = "#eff0f1"

    fields = [
      {
        type       = "Alphanumeric Text"
        title      = "Example Field"
        value      = "Demo"
        is_visible = false
      },
    ]
  }
}

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true

  oidc_options {
    type                        = "NATIVE_APP"
    grant_types                 = ["CLIENT_CREDENTIALS"]
    token_endpoint_authn_method = "CLIENT_SECRET_BASIC"

    mobile_app {
      bundle_id                = "com.pingidentity.ios_%[3]s"
      package_name             = "com.pingidentity.android_%[3]s"
      passcode_refresh_seconds = 30
    }
  }
}

resource "pingone_digital_wallet_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  application_id = resource.pingone_application.%[2]s.id
  name           = "%[3]s"
  app_open_url   = "https://www.example.com"

  depends_on = [resource.pingone_application.%[2]s]
}

resource "pingone_credential_issuance_rule" "%[2]s" {
  environment_id                = data.pingone_environment.general_test.id
  credential_type_id            = resource.pingone_credential_type.%[2]s.id
  digital_wallet_application_id = resource.pingone_digital_wallet_application.%[2]s.id
  status                        = "ACTIVE"

  filter = {
    population_ids = []
  }

  automation = {
    issue  = "ON_DEMAND"
    revoke = "PERIODIC"
    update = "ON_DEMAND"
  }

  notification = {
    methods = ["EMAIL", "SMS"]
    template = {
      locale  = "en"
      variant = "template_B"
    }
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccCredentialIssuanceRuleConfig_InvalidFilterAttribute(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_population" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
}

resource "pingone_credential_type" "%[2]s" {
  environment_id       = data.pingone_environment.general_test.id
  title                = "%[3]s"
  description          = "%[3]s Example Description"
  card_type            = "%[3]s"
  card_design_template = "<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 740 480\"><rect fill=\"none\" width=\"736\" height=\"476\" stroke=\"#CACED3\" stroke-width=\"3\" rx=\"10\" ry=\"10\" x=\"2\" y=\"2\"></rect><rect fill=\"$${cardColor}\" height=\"476\" rx=\"10\" ry=\"10\" width=\"736\" x=\"2\" y=\"2\" opacity=\"$${bgOpacityPercent}\"></rect><line y2=\"160\" x2=\"695\" y1=\"160\" x1=\"42.5\" stroke=\"$${textColor}\"></line><text fill=\"$${textColor}\" font-weight=\"450\" font-size=\"30\" x=\"160\" y=\"90\">$${cardTitle}</text><text fill=\"$${textColor}\" font-size=\"25\" font-weight=\"300\" x=\"160\" y=\"130\">$${cardSubtitle}</text></svg>"

  metadata = {
    name               = "%[3]s"
    description        = "%[3]s Example Description"
    bg_opacity_percent = 100
    card_color         = "#000000"
    text_color         = "#eff0f1"

    fields = [
      {
        type       = "Alphanumeric Text"
        title      = "Example Field"
        value      = "Demo"
        is_visible = false
      },
    ]
  }
}

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true

  oidc_options {
    type                        = "NATIVE_APP"
    grant_types                 = ["CLIENT_CREDENTIALS"]
    token_endpoint_authn_method = "CLIENT_SECRET_BASIC"

    mobile_app {
      bundle_id                = "com.pingidentity.ios_%[3]s"
      package_name             = "com.pingidentity.android_%[3]s"
      passcode_refresh_seconds = 30
    }
  }
}

resource "pingone_digital_wallet_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  application_id = resource.pingone_application.%[2]s.id
  name           = "%[3]s"
  app_open_url   = "https://www.example.com"

  depends_on = [resource.pingone_application.%[2]s]
}

resource "pingone_credential_issuance_rule" "%[2]s" {
  environment_id                = data.pingone_environment.general_test.id
  credential_type_id            = resource.pingone_credential_type.%[2]s.id
  digital_wallet_application_id = resource.pingone_digital_wallet_application.%[2]s.id
  status                        = "ACTIVE"

  filter = {}

  automation = {
    issue  = "ON_DEMAND"
    revoke = "PERIODIC"
    update = "ON_DEMAND"
  }

  notification = {
    methods = ["EMAIL", "SMS"]
    template = {
      locale  = "en"
      variant = "template_B"
    }
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccCredentialIssuanceRuleConfig_InvalidAutomationAttribute(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_population" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
}

resource "pingone_credential_type" "%[2]s" {
  environment_id       = data.pingone_environment.general_test.id
  title                = "%[3]s"
  description          = "%[3]s Example Description"
  card_type            = "%[3]s"
  card_design_template = "<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 740 480\"><rect fill=\"none\" width=\"736\" height=\"476\" stroke=\"#CACED3\" stroke-width=\"3\" rx=\"10\" ry=\"10\" x=\"2\" y=\"2\"></rect><rect fill=\"$${cardColor}\" height=\"476\" rx=\"10\" ry=\"10\" width=\"736\" x=\"2\" y=\"2\" opacity=\"$${bgOpacityPercent}\"></rect><line y2=\"160\" x2=\"695\" y1=\"160\" x1=\"42.5\" stroke=\"$${textColor}\"></line><text fill=\"$${textColor}\" font-weight=\"450\" font-size=\"30\" x=\"160\" y=\"90\">$${cardTitle}</text><text fill=\"$${textColor}\" font-size=\"25\" font-weight=\"300\" x=\"160\" y=\"130\">$${cardSubtitle}</text></svg>"

  metadata = {
    name               = "%[3]s"
    description        = "%[3]s Example Description"
    bg_opacity_percent = 100
    card_color         = "#000000"
    text_color         = "#eff0f1"

    fields = [
      {
        type       = "Alphanumeric Text"
        title      = "Example Field"
        value      = "Demo"
        is_visible = false
      },
    ]
  }
}

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true

  oidc_options {
    type                        = "NATIVE_APP"
    grant_types                 = ["CLIENT_CREDENTIALS"]
    token_endpoint_authn_method = "CLIENT_SECRET_BASIC"

    mobile_app {
      bundle_id                = "com.pingidentity.ios_%[3]s"
      package_name             = "com.pingidentity.android_%[3]s"
      passcode_refresh_seconds = 30
    }
  }
}

resource "pingone_digital_wallet_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  application_id = resource.pingone_application.%[2]s.id
  name           = "%[3]s"
  app_open_url   = "https://www.example.com"

  depends_on = [resource.pingone_application.%[2]s]
}

resource "pingone_credential_issuance_rule" "%[2]s" {
  environment_id                = data.pingone_environment.general_test.id
  credential_type_id            = resource.pingone_credential_type.%[2]s.id
  digital_wallet_application_id = resource.pingone_digital_wallet_application.%[2]s.id
  status                        = "ACTIVE"

  filter = {
    scim = "address.countryCode eq \"NG\""
  }

  automation = {}

  notification = {
    methods = ["EMAIL", "SMS"]
    template = {
      locale  = "en"
      variant = "template_B"
    }
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccCredentialIssuanceRuleConfig_InvalidNotificationAttribute(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_population" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
}

resource "pingone_credential_type" "%[2]s" {
  environment_id       = data.pingone_environment.general_test.id
  title                = "%[3]s"
  description          = "%[3]s Example Description"
  card_type            = "%[3]s"
  card_design_template = "<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 740 480\"><rect fill=\"none\" width=\"736\" height=\"476\" stroke=\"#CACED3\" stroke-width=\"3\" rx=\"10\" ry=\"10\" x=\"2\" y=\"2\"></rect><rect fill=\"$${cardColor}\" height=\"476\" rx=\"10\" ry=\"10\" width=\"736\" x=\"2\" y=\"2\" opacity=\"$${bgOpacityPercent}\"></rect><line y2=\"160\" x2=\"695\" y1=\"160\" x1=\"42.5\" stroke=\"$${textColor}\"></line><text fill=\"$${textColor}\" font-weight=\"450\" font-size=\"30\" x=\"160\" y=\"90\">$${cardTitle}</text><text fill=\"$${textColor}\" font-size=\"25\" font-weight=\"300\" x=\"160\" y=\"130\">$${cardSubtitle}</text></svg>"

  metadata = {
    name               = "%[3]s"
    description        = "%[3]s Example Description"
    bg_opacity_percent = 100
    card_color         = "#000000"
    text_color         = "#eff0f1"

    fields = [
      {
        type       = "Alphanumeric Text"
        title      = "Example Field"
        value      = "Demo"
        is_visible = false
      },
    ]
  }
}

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true

  oidc_options {
    type                        = "NATIVE_APP"
    grant_types                 = ["CLIENT_CREDENTIALS"]
    token_endpoint_authn_method = "CLIENT_SECRET_BASIC"

    mobile_app {
      bundle_id                = "com.pingidentity.ios_%[3]s"
      package_name             = "com.pingidentity.android_%[3]s"
      passcode_refresh_seconds = 30
    }
  }
}

resource "pingone_digital_wallet_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  application_id = resource.pingone_application.%[2]s.id
  name           = "%[3]s"
  app_open_url   = "https://www.example.com"

  depends_on = [resource.pingone_application.%[2]s]
}

resource "pingone_credential_issuance_rule" "%[2]s" {
  environment_id                = data.pingone_environment.general_test.id
  credential_type_id            = resource.pingone_credential_type.%[2]s.id
  digital_wallet_application_id = resource.pingone_digital_wallet_application.%[2]s.id
  status                        = "ACTIVE"

  filter = {
    scim = "address.countryCode eq \"NG\""
  }

  automation = {
    issue  = "ON_DEMAND"
    revoke = "PERIODIC"
    update = "ON_DEMAND"
  }

  notification = {}
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccCredentialIssuanceRuleConfig_InvalidNotificationMethodsAttribute(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_population" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
}

resource "pingone_credential_type" "%[2]s" {
  environment_id       = data.pingone_environment.general_test.id
  title                = "%[3]s"
  description          = "%[3]s Example Description"
  card_type            = "%[3]s"
  card_design_template = "<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 740 480\"><rect fill=\"none\" width=\"736\" height=\"476\" stroke=\"#CACED3\" stroke-width=\"3\" rx=\"10\" ry=\"10\" x=\"2\" y=\"2\"></rect><rect fill=\"$${cardColor}\" height=\"476\" rx=\"10\" ry=\"10\" width=\"736\" x=\"2\" y=\"2\" opacity=\"$${bgOpacityPercent}\"></rect><line y2=\"160\" x2=\"695\" y1=\"160\" x1=\"42.5\" stroke=\"$${textColor}\"></line><text fill=\"$${textColor}\" font-weight=\"450\" font-size=\"30\" x=\"160\" y=\"90\">$${cardTitle}</text><text fill=\"$${textColor}\" font-size=\"25\" font-weight=\"300\" x=\"160\" y=\"130\">$${cardSubtitle}</text></svg>"

  metadata = {
    name               = "%[3]s"
    description        = "%[3]s Example Description"
    bg_opacity_percent = 100
    card_color         = "#000000"
    text_color         = "#eff0f1"

    fields = [
      {
        type       = "Alphanumeric Text"
        title      = "Example Field"
        value      = "Demo"
        is_visible = false
      },
    ]
  }
}

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true

  oidc_options {
    type                        = "NATIVE_APP"
    grant_types                 = ["CLIENT_CREDENTIALS"]
    token_endpoint_authn_method = "CLIENT_SECRET_BASIC"

    mobile_app {
      bundle_id                = "com.pingidentity.ios_%[3]s"
      package_name             = "com.pingidentity.android_%[3]s"
      passcode_refresh_seconds = 30
    }
  }
}

resource "pingone_digital_wallet_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  application_id = resource.pingone_application.%[2]s.id
  name           = "%[3]s"
  app_open_url   = "https://www.example.com"

  depends_on = [resource.pingone_application.%[2]s]
}

resource "pingone_credential_issuance_rule" "%[2]s" {
  environment_id                = data.pingone_environment.general_test.id
  credential_type_id            = resource.pingone_credential_type.%[2]s.id
  digital_wallet_application_id = resource.pingone_digital_wallet_application.%[2]s.id
  status                        = "ACTIVE"

  filter = {
    scim = "address.countryCode eq \"NG\""
  }

  automation = {
    issue  = "ON_DEMAND"
    revoke = "PERIODIC"
    update = "ON_DEMAND"
  }

  notification = {
    methods = []
    template = {
      locale  = "en"
      variant = "template_B"
    }
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccCredentialIssuanceRuleConfig_InvalidNotificationTemplateAttribute(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_population" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
}

resource "pingone_credential_type" "%[2]s" {
  environment_id       = data.pingone_environment.general_test.id
  title                = "%[3]s"
  description          = "%[3]s Example Description"
  card_type            = "%[3]s"
  card_design_template = "<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 740 480\"><rect fill=\"none\" width=\"736\" height=\"476\" stroke=\"#CACED3\" stroke-width=\"3\" rx=\"10\" ry=\"10\" x=\"2\" y=\"2\"></rect><rect fill=\"$${cardColor}\" height=\"476\" rx=\"10\" ry=\"10\" width=\"736\" x=\"2\" y=\"2\" opacity=\"$${bgOpacityPercent}\"></rect><line y2=\"160\" x2=\"695\" y1=\"160\" x1=\"42.5\" stroke=\"$${textColor}\"></line><text fill=\"$${textColor}\" font-weight=\"450\" font-size=\"30\" x=\"160\" y=\"90\">$${cardTitle}</text><text fill=\"$${textColor}\" font-size=\"25\" font-weight=\"300\" x=\"160\" y=\"130\">$${cardSubtitle}</text></svg>"

  metadata = {
    name               = "%[3]s"
    description        = "%[3]s Example Description"
    bg_opacity_percent = 100
    card_color         = "#000000"
    text_color         = "#eff0f1"

    fields = [
      {
        type       = "Alphanumeric Text"
        title      = "Example Field"
        value      = "Demo"
        is_visible = false
      },
    ]
  }
}

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true

  oidc_options {
    type                        = "NATIVE_APP"
    grant_types                 = ["CLIENT_CREDENTIALS"]
    token_endpoint_authn_method = "CLIENT_SECRET_BASIC"

    mobile_app {
      bundle_id                = "com.pingidentity.ios_%[3]s"
      package_name             = "com.pingidentity.android_%[3]s"
      passcode_refresh_seconds = 30
    }
  }
}

resource "pingone_digital_wallet_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  application_id = resource.pingone_application.%[2]s.id
  name           = "%[3]s"
  app_open_url   = "https://www.example.com"

  depends_on = [resource.pingone_application.%[2]s]
}

resource "pingone_credential_issuance_rule" "%[2]s" {
  environment_id                = data.pingone_environment.general_test.id
  credential_type_id            = resource.pingone_credential_type.%[2]s.id
  digital_wallet_application_id = resource.pingone_digital_wallet_application.%[2]s.id
  status                        = "ACTIVE"

  filter = {
    scim = "address.countryCode eq \"NG\""
  }

  automation = {
    issue  = "ON_DEMAND"
    revoke = "PERIODIC"
    update = "ON_DEMAND"
  }

  notification = {
    methods  = ["EMAIL", "SMS"]
    template = {}
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
