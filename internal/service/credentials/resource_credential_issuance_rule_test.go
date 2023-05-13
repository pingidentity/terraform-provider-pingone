package credentials_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/patrickcping/pingone-go-sdk-v2/credentials"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func testAccCheckCredentiaIssuanceRuleDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.CredentialsAPIClient
	ctx = context.WithValue(ctx, credentials.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})

	mgmtApiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_credential_issuance_rule" {
			continue
		}

		_, rEnv, err := mgmtApiClient.EnvironmentsApi.ReadOneEnvironment(ctx, rs.Primary.Attributes["environment_id"]).Execute()

		if err != nil {

			if rEnv == nil {
				return fmt.Errorf("Response object does not exist and no error detected")
			}

			if rEnv.StatusCode == 404 {
				continue
			}

			return err
		}

		body, r, err := apiClient.CredentialIssuanceRulesApi.ReadOneCredentialIssuanceRule(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["credential_type_id"], rs.Primary.Attributes["id"]).Execute()

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

		return fmt.Errorf("PingOne Credential Issuance Rule ID %s still exists", rs.Primary.ID)
	}

	return nil
}

func TestAccCredentialIssuanceRule_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_credential_issuance_rule.%s", resourceName)

	name := acctest.ResourceNameGen()

	fullStep := resource.TestStep{
		Config: testAccCredentialIssuanceRule_Full(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
			resource.TestMatchResourceAttr(resourceFullName, "credential_type_id", verify.P1ResourceIDRegexp),
			resource.TestMatchResourceAttr(resourceFullName, "digital_wallet_application_id", verify.P1ResourceIDRegexp),
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
		Config: testAccCredentialIssuanceRule_Minimal(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
			resource.TestMatchResourceAttr(resourceFullName, "credential_type_id", verify.P1ResourceIDRegexp),
			resource.TestMatchResourceAttr(resourceFullName, "digital_wallet_application_id", verify.P1ResourceIDRegexp),
			resource.TestCheckResourceAttr(resourceFullName, "automation.issue", "PERIODIC"),
			resource.TestCheckResourceAttr(resourceFullName, "automation.revoke", "PERIODIC"),
			resource.TestCheckResourceAttr(resourceFullName, "automation.update", "PERIODIC"),
			resource.TestCheckResourceAttr(resourceFullName, "filter.scim", "address.countryCode eq \"CA\""),
			resource.TestCheckNoResourceAttr(resourceFullName, "filter.population_ids.#"),
			resource.TestCheckResourceAttr(resourceFullName, "notification.methods.#", "1"),
			resource.TestCheckResourceAttr(resourceFullName, "notification.methods.0", "EMAIL"),
			resource.TestCheckNoResourceAttr(resourceFullName, "notification.template.%"),
			resource.TestCheckResourceAttr(resourceFullName, "status", "ACTIVE"),
		),
	}

	disabledStep := resource.TestStep{
		Config: testAccCredentialIssuanceRule_Disabled(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
			resource.TestMatchResourceAttr(resourceFullName, "credential_type_id", verify.P1ResourceIDRegexp),
			resource.TestMatchResourceAttr(resourceFullName, "digital_wallet_application_id", verify.P1ResourceIDRegexp),
			resource.TestCheckResourceAttrSet(resourceFullName, "automation.%"),
			resource.TestCheckResourceAttrSet(resourceFullName, "filter.%"),
			resource.TestCheckNoResourceAttr(resourceFullName, "notification.%"),
			resource.TestCheckResourceAttr(resourceFullName, "status", "DISABLED"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             nil,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// full
			fullStep,
			{
				Config:  testAccCredentialIssuanceRule_Full(resourceName, name),
				Destroy: true,
			},
			minimalStep,
			{
				Config:  testAccCredentialIssuanceRule_Minimal(resourceName, name),
				Destroy: true,
			},
			disabledStep,
			{
				Config:  testAccCredentialIssuanceRule_Disabled(resourceName, name),
				Destroy: true,
			},
			fullStep,
			minimalStep,
			fullStep,
			disabledStep,
			fullStep,
		},
	})
}

func TestAccCredentialIssuanceRule_InvalidConfigs(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	name := acctest.ResourceNameGen()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCredentialTypeDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccCredentialIssuanceRule_InvalidCredentialTypeID(resourceName, name),
				ExpectError: regexp.MustCompile("Error: Error when calling `CreateCredentialIssuanceRule`: The request could not be completed. One or more validation errors were in the request."),
			},
			{
				Config:      testAccCredentialIssuanceRule_InvalidDigitalWalletID(resourceName, name),
				ExpectError: regexp.MustCompile("Error: Error when calling `CreateCredentialIssuanceRule`: The request could not be completed. One or more validation errors were in the request."),
			},
			{
				Config:      testAccCredentialIssuanceRule_InvalidScimFilterAttribute(resourceName, name),
				ExpectError: regexp.MustCompile("Error: Error when calling `CreateCredentialIssuanceRule`: The request could not be completed. One or more validation errors were in the request."),
			},
		},
	})
}

func testAccCredentialIssuanceRule_Full(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_population" "%[2]s" {
	environment_id = data.pingone_environment.credentials_test.id
    name = "%[3]s"
}

resource "pingone_credential_type" "%[2]s" {
  environment_id = data.pingone_environment.credentials_test.id
  title = "%[3]s"
  description = "%[3]s Example Description"
  card_type = "%[3]s"
  card_design_template = "<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 740 480\"><rect fill=\"none\" width=\"736\" height=\"476\" stroke=\"#CACED3\" stroke-width=\"3\" rx=\"10\" ry=\"10\" x=\"2\" y=\"2\"></rect><rect fill=\"$${cardColor}\" height=\"476\" rx=\"10\" ry=\"10\" width=\"736\" x=\"2\" y=\"2\" opacity=\"$${bgOpacityPercent}\"></rect><line y2=\"160\" x2=\"695\" y1=\"160\" x1=\"42.5\" stroke=\"$${textColor}\"></line><text fill=\"$${textColor}\" font-weight=\"450\" font-size=\"30\" x=\"160\" y=\"90\">$${cardTitle}</text><text fill=\"$${textColor}\" font-size=\"25\" font-weight=\"300\" x=\"160\" y=\"130\">$${cardSubtitle}</text></svg>"

  metadata = {
    name = "%[3]s"
    description = "%[3]s Example Description"
    version = 5
    bg_opacity_percent = 100
    card_color = "#000000"
    text_color = "#eff0f1"

    fields = [
      {
        type = "Alphanumeric Text"
        title = "Example Field"
        value = "Demo"
        is_visible = false
      },
    ]
  }
}

resource "pingone_application" "%[2]s" {
	environment_id = data.pingone_environment.credentials_test.id
	name = "%[3]s"
	enabled = true

	oidc_options {
	  type                        = "NATIVE_APP"
	  grant_types                 = ["CLIENT_CREDENTIALS"]
	  token_endpoint_authn_method = "CLIENT_SECRET_BASIC"
	  bundle_id        = "com.pingidentity.ios_%[3]s"
	  package_name     = "com.pingidentity.android_%[3]s"
  
	  mobile_app {
		 bundle_id        = "com.pingidentity.ios_%[3]s"
		 package_name     = "com.pingidentity.android_%[3]s"
		 passcode_refresh_seconds = 30
	  }
	}
}

resource "pingone_digital_wallet_application" "%[2]s" {
	environment_id = data.pingone_environment.credentials_test.id
	application_id = resource.pingone_application.%[2]s.id
	name = "%[3]s"
	app_open_url = "https://www.example.com"	

	depends_on = [resource.pingone_application.%[2]s]
}

resource "pingone_credential_issuance_rule" "%[2]s" {
	environment_id = data.pingone_environment.credentials_test.id
	credential_type_id = resource.pingone_credential_type.%[2]s.id
	digital_wallet_application_id = resource.pingone_digital_wallet_application.%[2]s.id
	status = "ACTIVE"
	
	filter = {
	  population_ids = [resource.pingone_population.%[2]s.id]
	}
  
	automation = {
	  issue = "ON_DEMAND"
	  revoke = "PERIODIC"
	  update = "ON_DEMAND"
	}
  
	notification = {
	  methods = ["EMAIL", "SMS"]
	  template = {
		locale = "en"
		variant = "template_B"
	  }
	}
}`, acctest.CredentialsSandboxEnvironment(), resourceName, name)
}

func testAccCredentialIssuanceRule_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_credential_type" "%[2]s" {
  environment_id = data.pingone_environment.credentials_test.id
  title = "%[3]s"
  description = "%[3]s Example Description"
  card_type = "%[3]s"
  card_design_template = "<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 740 480\"><rect fill=\"none\" width=\"736\" height=\"476\" stroke=\"#CACED3\" stroke-width=\"3\" rx=\"10\" ry=\"10\" x=\"2\" y=\"2\"></rect><rect fill=\"$${cardColor}\" height=\"476\" rx=\"10\" ry=\"10\" width=\"736\" x=\"2\" y=\"2\" opacity=\"$${bgOpacityPercent}\"></rect><line y2=\"160\" x2=\"695\" y1=\"160\" x1=\"42.5\" stroke=\"$${textColor}\"></line><text fill=\"$${textColor}\" font-weight=\"450\" font-size=\"30\" x=\"160\" y=\"90\">$${cardTitle}</text><text fill=\"$${textColor}\" font-size=\"25\" font-weight=\"300\" x=\"160\" y=\"130\">$${cardSubtitle}</text></svg>"

  metadata = {
    name = "%[3]s"
    description = "%[3]s Example Description"
    version = 5
    bg_opacity_percent = 100
    card_color = "#000000"
    text_color = "#eff0f1"

    fields = [
      {
        type = "Alphanumeric Text"
        title = "Example Field"
        value = "Demo"
        is_visible = false
      },
    ]
  }
}

resource "pingone_application" "%[2]s" {
	environment_id = data.pingone_environment.credentials_test.id
	name = "%[3]s"
	enabled = true

	oidc_options {
	  type                        = "NATIVE_APP"
	  grant_types                 = ["CLIENT_CREDENTIALS"]
	  token_endpoint_authn_method = "CLIENT_SECRET_BASIC"
	  bundle_id        = "com.pingidentity.ios_%[3]s"
	  package_name     = "com.pingidentity.android_%[3]s"
  
	  mobile_app {
		 bundle_id        = "com.pingidentity.ios_%[3]s"
		 package_name     = "com.pingidentity.android_%[3]s"
		 passcode_refresh_seconds = 30
	  }
	}
}

resource "pingone_digital_wallet_application" "%[2]s" {
	environment_id = data.pingone_environment.credentials_test.id
	application_id = resource.pingone_application.%[2]s.id
	name = "%[3]s"
	app_open_url = "https://www.example.com"	

	depends_on = [resource.pingone_application.%[2]s]
}

resource "pingone_credential_issuance_rule" "%[2]s" {
	environment_id = data.pingone_environment.credentials_test.id
	credential_type_id = resource.pingone_credential_type.%[2]s.id
	digital_wallet_application_id = resource.pingone_digital_wallet_application.%[2]s.id
	status = "ACTIVE"
	
	filter = {
	  scim = "address.countryCode eq \"CA\""
	}
  
	automation = {
	  issue = "PERIODIC"
	  revoke = "PERIODIC"
	  update = "PERIODIC"
	}
  
	notification = {
	  methods = ["EMAIL"]
	}
}`, acctest.CredentialsSandboxEnvironment(), resourceName, name)
}

func testAccCredentialIssuanceRule_Disabled(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_credential_type" "%[2]s" {
  environment_id = data.pingone_environment.credentials_test.id
  title = "%[3]s"
  description = "%[3]s Example Description"
  card_type = "%[3]s"
  card_design_template = "<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 740 480\"><rect fill=\"none\" width=\"736\" height=\"476\" stroke=\"#CACED3\" stroke-width=\"3\" rx=\"10\" ry=\"10\" x=\"2\" y=\"2\"></rect><rect fill=\"$${cardColor}\" height=\"476\" rx=\"10\" ry=\"10\" width=\"736\" x=\"2\" y=\"2\" opacity=\"$${bgOpacityPercent}\"></rect><line y2=\"160\" x2=\"695\" y1=\"160\" x1=\"42.5\" stroke=\"$${textColor}\"></line><text fill=\"$${textColor}\" font-weight=\"450\" font-size=\"30\" x=\"160\" y=\"90\">$${cardTitle}</text><text fill=\"$${textColor}\" font-size=\"25\" font-weight=\"300\" x=\"160\" y=\"130\">$${cardSubtitle}</text></svg>"

  metadata = {
    name = "%[3]s"
    description = "%[3]s Example Description"
    version = 5
    bg_opacity_percent = 100
    card_color = "#000000"
    text_color = "#eff0f1"

    fields = [
      {
        type = "Alphanumeric Text"
        title = "Example Field"
        value = "Demo"
        is_visible = false
      },
    ]
  }
}

resource "pingone_application" "%[2]s" {
	environment_id = data.pingone_environment.credentials_test.id
	name = "%[3]s"
	enabled = true

	oidc_options {
	  type                        = "NATIVE_APP"
	  grant_types                 = ["CLIENT_CREDENTIALS"]
	  token_endpoint_authn_method = "CLIENT_SECRET_BASIC"
	  bundle_id        = "com.pingidentity.ios_%[3]s"
	  package_name     = "com.pingidentity.android_%[3]s"
  
	  mobile_app {
		 bundle_id        = "com.pingidentity.ios_%[3]s"
		 package_name     = "com.pingidentity.android_%[3]s"
		 passcode_refresh_seconds = 30
	  }
	}
}

resource "pingone_digital_wallet_application" "%[2]s" {
	environment_id = data.pingone_environment.credentials_test.id
	application_id = resource.pingone_application.%[2]s.id
	name = "%[3]s"
	app_open_url = "https://www.example.com"	

	depends_on = [resource.pingone_application.%[2]s]
}

resource "pingone_credential_issuance_rule" "%[2]s" {
	environment_id = data.pingone_environment.credentials_test.id
	credential_type_id = resource.pingone_credential_type.%[2]s.id
	digital_wallet_application_id = resource.pingone_digital_wallet_application.%[2]s.id
	status = "DISABLED"
	
	filter = {
	  scim = "address.countryCode eq \"CA\""
	}
  
	automation = {
	  issue = "PERIODIC"
	  revoke = "PERIODIC"
	  update = "PERIODIC"	  	  
	}

}`, acctest.CredentialsSandboxEnvironment(), resourceName, name)
}

func testAccCredentialIssuanceRule_InvalidCredentialTypeID(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s


resource "pingone_application" "%[2]s" {
	environment_id = data.pingone_environment.credentials_test.id
	name = "%[3]s"
	enabled = true

	oidc_options {
	  type                        = "NATIVE_APP"
	  grant_types                 = ["CLIENT_CREDENTIALS"]
	  token_endpoint_authn_method = "CLIENT_SECRET_BASIC"
	  bundle_id        = "com.pingidentity.ios_%[3]s"
	  package_name     = "com.pingidentity.android_%[3]s"
  
	  mobile_app {
		 bundle_id        = "com.pingidentity.ios_%[3]s"
		 package_name     = "com.pingidentity.android_%[3]s"
		 passcode_refresh_seconds = 30
	  }
	}
}

resource "pingone_digital_wallet_application" "%[2]s" {
	environment_id = data.pingone_environment.credentials_test.id
	application_id = resource.pingone_application.%[2]s.id
	name = "%[3]s"
	app_open_url = "https://www.example.com"	

	depends_on = [resource.pingone_application.%[2]s]
}

resource "pingone_credential_issuance_rule" "%[2]s" {
	environment_id = data.pingone_environment.credentials_test.id
	credential_type_id = "9c052a8a-14be-44e4-8f07-2662569994ce" // dummy ID that conforms to UUID v4
	digital_wallet_application_id = resource.pingone_digital_wallet_application.%[2]s.id
	status = "DISABLED"
	
	filter = {
	  scim = "address.countryCode eq \"CA\""
	}
  
	automation = {
	  issue = "PERIODIC"
	  revoke = "PERIODIC"
	  update = "PERIODIC"	  	  
	}

}`, acctest.CredentialsSandboxEnvironment(), resourceName, name)
}

func testAccCredentialIssuanceRule_InvalidDigitalWalletID(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_credential_type" "%[2]s" {
  environment_id = data.pingone_environment.credentials_test.id
  title = "%[3]s"
  description = "%[3]s Example Description"
  card_type = "%[3]s"
  card_design_template = "<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 740 480\"><rect fill=\"none\" width=\"736\" height=\"476\" stroke=\"#CACED3\" stroke-width=\"3\" rx=\"10\" ry=\"10\" x=\"2\" y=\"2\"></rect><rect fill=\"$${cardColor}\" height=\"476\" rx=\"10\" ry=\"10\" width=\"736\" x=\"2\" y=\"2\" opacity=\"$${bgOpacityPercent}\"></rect><line y2=\"160\" x2=\"695\" y1=\"160\" x1=\"42.5\" stroke=\"$${textColor}\"></line><text fill=\"$${textColor}\" font-weight=\"450\" font-size=\"30\" x=\"160\" y=\"90\">$${cardTitle}</text><text fill=\"$${textColor}\" font-size=\"25\" font-weight=\"300\" x=\"160\" y=\"130\">$${cardSubtitle}</text></svg>"

  metadata = {
    name = "%[3]s"
    description = "%[3]s Example Description"
    version = 5
    bg_opacity_percent = 100
    card_color = "#000000"
    text_color = "#eff0f1"

    fields = [
      {
        type = "Alphanumeric Text"
        title = "Example Field"
        value = "Demo"
        is_visible = false
      },
    ]
  }
}

resource "pingone_application" "%[2]s" {
	environment_id = data.pingone_environment.credentials_test.id
	name = "%[3]s"
	enabled = true

	oidc_options {
	  type                        = "NATIVE_APP"
	  grant_types                 = ["CLIENT_CREDENTIALS"]
	  token_endpoint_authn_method = "CLIENT_SECRET_BASIC"
	  bundle_id        = "com.pingidentity.ios_%[3]s"
	  package_name     = "com.pingidentity.android_%[3]s"
  
	  mobile_app {
		 bundle_id        = "com.pingidentity.ios_%[3]s"
		 package_name     = "com.pingidentity.android_%[3]s"
		 passcode_refresh_seconds = 30
	  }
	}
}

resource "pingone_digital_wallet_application" "%[2]s" {
	environment_id = data.pingone_environment.credentials_test.id
	application_id = resource.pingone_application.%[2]s.id
	name = "%[3]s"
	app_open_url = "https://www.example.com"	

	depends_on = [resource.pingone_application.%[2]s]
}

resource "pingone_credential_issuance_rule" "%[2]s" {
	environment_id = data.pingone_environment.credentials_test.id
	credential_type_id = resource.pingone_credential_type.%[2]s.id
	digital_wallet_application_id = "9c052a8a-14be-44e4-8f07-2662569994ce" // dummy ID that conforms to UUID v4
	status = "DISABLED"
	
	filter = {
	  scim = "address.countryCode eq \"CA\""
	}
  
	automation = {
	  issue = "PERIODIC"
	  revoke = "PERIODIC"
	  update = "PERIODIC"	  	  
	}

}`, acctest.CredentialsSandboxEnvironment(), resourceName, name)
}

func testAccCredentialIssuanceRule_InvalidScimFilterAttribute(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_population" "%[2]s" {
	environment_id = data.pingone_environment.credentials_test.id
    name = "%[3]s"
}

resource "pingone_credential_type" "%[2]s" {
  environment_id = data.pingone_environment.credentials_test.id
  title = "%[3]s"
  description = "%[3]s Example Description"
  card_type = "%[3]s"
  card_design_template = "<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 740 480\"><rect fill=\"none\" width=\"736\" height=\"476\" stroke=\"#CACED3\" stroke-width=\"3\" rx=\"10\" ry=\"10\" x=\"2\" y=\"2\"></rect><rect fill=\"$${cardColor}\" height=\"476\" rx=\"10\" ry=\"10\" width=\"736\" x=\"2\" y=\"2\" opacity=\"$${bgOpacityPercent}\"></rect><line y2=\"160\" x2=\"695\" y1=\"160\" x1=\"42.5\" stroke=\"$${textColor}\"></line><text fill=\"$${textColor}\" font-weight=\"450\" font-size=\"30\" x=\"160\" y=\"90\">$${cardTitle}</text><text fill=\"$${textColor}\" font-size=\"25\" font-weight=\"300\" x=\"160\" y=\"130\">$${cardSubtitle}</text></svg>"

  metadata = {
    name = "%[3]s"
    description = "%[3]s Example Description"
    version = 5
    bg_opacity_percent = 100
    card_color = "#000000"
    text_color = "#eff0f1"

    fields = [
      {
        type = "Alphanumeric Text"
        title = "Example Field"
        value = "Demo"
        is_visible = false
      },
    ]
  }
}

resource "pingone_application" "%[2]s" {
	environment_id = data.pingone_environment.credentials_test.id
	name = "%[3]s"
	enabled = true

	oidc_options {
	  type                        = "NATIVE_APP"
	  grant_types                 = ["CLIENT_CREDENTIALS"]
	  token_endpoint_authn_method = "CLIENT_SECRET_BASIC"
	  bundle_id        = "com.pingidentity.ios_%[3]s"
	  package_name     = "com.pingidentity.android_%[3]s"
  
	  mobile_app {
		 bundle_id        = "com.pingidentity.ios_%[3]s"
		 package_name     = "com.pingidentity.android_%[3]s"
		 passcode_refresh_seconds = 30
	  }
	}
}

resource "pingone_digital_wallet_application" "%[2]s" {
	environment_id = data.pingone_environment.credentials_test.id
	application_id = resource.pingone_application.%[2]s.id
	name = "%[3]s"
	app_open_url = "https://www.example.com"	

	depends_on = [resource.pingone_application.%[2]s]
}

resource "pingone_credential_issuance_rule" "%[2]s" {
	environment_id = data.pingone_environment.credentials_test.id
	credential_type_id = resource.pingone_credential_type.%[2]s.id
	digital_wallet_application_id = resource.pingone_digital_wallet_application.%[2]s.id
	status = "ACTIVE"
	
	filter = {
		scim = "invalidAttribute eq \"Users\""
	}
  
	automation = {
	  issue = "ON_DEMAND"
	  revoke = "PERIODIC"
	  update = "ON_DEMAND"
	}

}`, acctest.CredentialsSandboxEnvironment(), resourceName, name)
}

func testAccCredentialIssuanceRule_InvalidAutomation(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_population" "%[2]s" {
	environment_id = data.pingone_environment.credentials_test.id
    name = "%[3]s"
}

resource "pingone_credential_type" "%[2]s" {
  environment_id = data.pingone_environment.credentials_test.id
  title = "%[3]s"
  description = "%[3]s Example Description"
  card_type = "%[3]s"
  card_design_template = "<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 740 480\"><rect fill=\"none\" width=\"736\" height=\"476\" stroke=\"#CACED3\" stroke-width=\"3\" rx=\"10\" ry=\"10\" x=\"2\" y=\"2\"></rect><rect fill=\"$${cardColor}\" height=\"476\" rx=\"10\" ry=\"10\" width=\"736\" x=\"2\" y=\"2\" opacity=\"$${bgOpacityPercent}\"></rect><line y2=\"160\" x2=\"695\" y1=\"160\" x1=\"42.5\" stroke=\"$${textColor}\"></line><text fill=\"$${textColor}\" font-weight=\"450\" font-size=\"30\" x=\"160\" y=\"90\">$${cardTitle}</text><text fill=\"$${textColor}\" font-size=\"25\" font-weight=\"300\" x=\"160\" y=\"130\">$${cardSubtitle}</text></svg>"

  metadata = {
    name = "%[3]s"
    description = "%[3]s Example Description"
    version = 5
    bg_opacity_percent = 100
    card_color = "#000000"
    text_color = "#eff0f1"

    fields = [
      {
        type = "Alphanumeric Text"
        title = "Example Field"
        value = "Demo"
        is_visible = false
      },
    ]
  }
}

resource "pingone_application" "%[2]s" {
	environment_id = data.pingone_environment.credentials_test.id
	name = "%[3]s"
	enabled = true

	oidc_options {
	  type                        = "NATIVE_APP"
	  grant_types                 = ["CLIENT_CREDENTIALS"]
	  token_endpoint_authn_method = "CLIENT_SECRET_BASIC"
	  bundle_id        = "com.pingidentity.ios_%[3]s"
	  package_name     = "com.pingidentity.android_%[3]s"
  
	  mobile_app {
		 bundle_id        = "com.pingidentity.ios_%[3]s"
		 package_name     = "com.pingidentity.android_%[3]s"
		 passcode_refresh_seconds = 30
	  }
	}
}

resource "pingone_digital_wallet_application" "%[2]s" {
	environment_id = data.pingone_environment.credentials_test.id
	application_id = resource.pingone_application.%[2]s.id
	name = "%[3]s"
	app_open_url = "https://www.example.com"	

	depends_on = [resource.pingone_application.%[2]s]
}

resource "pingone_credential_issuance_rule" "%[2]s" {
	environment_id = data.pingone_environment.credentials_test.id
	credential_type_id = resource.pingone_credential_type.%[2]s.id
	digital_wallet_application_id = resource.pingone_digital_wallet_application.%[2]s.id
	status = "ACTIVE"
	
	filter = {
		scim = "invalidAttribute eq \"Users\""
	}
  
	automation = {
		issue = "PERIODIC"
		revoke = "PERIODIC"
		update = "PERIODIC"	  	  
	}

}`, acctest.CredentialsSandboxEnvironment(), resourceName, name)
}
