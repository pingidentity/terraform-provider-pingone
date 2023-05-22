package credentials_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccCredentialIssuanceRuleDataSource_ByIDFull(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_credential_issuance_rule.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.pingone_credential_issuance_rule.%s", resourceName)

	name := acctest.ResourceNameGen()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCredentiaIssuanceRuleDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCredentialIssuanceRuleDataSource_ByIDFull(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(dataSourceFullName, "credential_type_id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(dataSourceFullName, "credential_issuance_rule_id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(dataSourceFullName, "digital_wallet_application_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "automation.%", resourceFullName, "automation.%"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "filter.%", resourceFullName, "filter.%"),
					resource.TestCheckResourceAttr(dataSourceFullName, "filter.scim", "accountId eq \"12345\" or accountId eq \"98765\" or (address.countryCode eq \"US\")"),
					resource.TestCheckResourceAttr(dataSourceFullName, "automation.issue", "ON_DEMAND"),
					resource.TestCheckResourceAttr(dataSourceFullName, "automation.revoke", "PERIODIC"),
					resource.TestCheckResourceAttr(dataSourceFullName, "automation.update", "ON_DEMAND"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "notification.%", resourceFullName, "notification.%"),
					resource.TestCheckResourceAttr(dataSourceFullName, "notification.methods.#", "2"),
					resource.TestCheckResourceAttr(dataSourceFullName, "notification.template.locale", "en"),
					resource.TestCheckResourceAttr(dataSourceFullName, "notification.template.variant", "template_B"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "status", resourceFullName, "status"),
				),
			},
		},
	})
}

func TestAccCredentialIssuanceRuleDataSource_NotFound(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCredentiaIssuanceRuleDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccCredentialIssuanceRuleDataSource_NotFoundByID(resourceName),
				ExpectError: regexp.MustCompile("Error: Error when calling `ReadOneCredentialIssuanceRule`: The request could not be completed. The requested resource was not found."),
			},
		},
	})
}

func TestAccCredentialIssuanceRuleDataSource_InvalidConfig(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCredentiaIssuanceRuleDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccCredentialIssuanceRuleDataSource_NoEnvironmentID(resourceName),
				ExpectError: regexp.MustCompile("Error: Missing required argument"),
			},
			{
				Config:      testAccCredentialIssuanceRuleDataSource_NoCredentialIssuanceRuleID(resourceName),
				ExpectError: regexp.MustCompile("Error: Missing required argument"),
			},
		},
	})
}

func testAccCredentialIssuanceRuleDataSource_ByIDFull(resourceName, name string) string {
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
        type       = "Directory Attribute"
        title      = "givenName"
        attribute  = "name.given"
        is_visible = false
      },
      {
        type       = "Directory Attribute"
        title      = "surname"
        attribute  = "name.family"
        is_visible = false
      },
      {
        type       = "Directory Attribute"
        title      = "jobTitle"
        attribute  = "title"
        is_visible = false
      },
      {
        type       = "Directory Attribute"
        title      = "displayName"
        attribute  = "name.formatted"
        is_visible = false
      },
      {
        type       = "Directory Attribute"
        title      = "mail"
        attribute  = "email"
        is_visible = false
      },
      {
        type       = "Directory Attribute"
        title      = "preferredLanguage"
        attribute  = "preferredLanguage"
        is_visible = false
      },
      {
        type       = "Directory Attribute"
        title      = "id"
        attribute  = "id"
        is_visible = false
      }
    ]
  }
}

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[2]s"
  enabled        = true

  oidc_options {
    type                        = "NATIVE_APP"
    grant_types                 = ["CLIENT_CREDENTIALS"]
    token_endpoint_authn_method = "CLIENT_SECRET_BASIC"

    mobile_app {
      bundle_id                = "com.pingidentity.ios_wallet_byname"
      package_name             = "com.pingidentity.android_wallet_byname"
      passcode_refresh_seconds = 30
    }
  }
}

resource "pingone_digital_wallet_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  application_id = resource.pingone_application.%[2]s.id
  name           = "%[2]s"
  app_open_url   = "https://www.example.com"

  depends_on = [resource.pingone_application.%[2]s]
}

resource "pingone_credential_issuance_rule" "%[2]s" {
  environment_id                = data.pingone_environment.general_test.id
  credential_type_id            = resource.pingone_credential_type.%[2]s.id
  digital_wallet_application_id = resource.pingone_digital_wallet_application.%[2]s.id
  status                        = "ACTIVE"

  filter = {
    scim = "accountId eq \"12345\" or accountId eq \"98765\" or (address.countryCode eq \"US\")"
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
}

data "pingone_credential_issuance_rule" "%[2]s" {
  environment_id              = data.pingone_environment.general_test.id
  credential_type_id          = resource.pingone_credential_type.%[2]s.id
  credential_issuance_rule_id = resource.pingone_credential_issuance_rule.%[2]s.id

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccCredentialIssuanceRuleDataSource_NotFoundByID(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s
resource "pingone_credential_type" "%[2]s" {
  environment_id       = data.pingone_environment.general_test.id
  title                = "%[2]s"
  description          = "%[2]s Example Description"
  card_type            = "%[2]s"
  card_design_template = "<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 740 480\"><rect fill=\"none\" width=\"736\" height=\"476\" stroke=\"#CACED3\" stroke-width=\"3\" rx=\"10\" ry=\"10\" x=\"2\" y=\"2\"></rect><rect fill=\"$${cardColor}\" height=\"476\" rx=\"10\" ry=\"10\" width=\"736\" x=\"2\" y=\"2\" opacity=\"$${bgOpacityPercent}\"></rect><line y2=\"160\" x2=\"695\" y1=\"160\" x1=\"42.5\" stroke=\"$${textColor}\"></line><text fill=\"$${textColor}\" font-weight=\"450\" font-size=\"30\" x=\"160\" y=\"90\">$${cardTitle}</text><text fill=\"$${textColor}\" font-size=\"25\" font-weight=\"300\" x=\"160\" y=\"130\">$${cardSubtitle}</text></svg>"

  metadata = {
    name               = "%[2]s"
    description        = "%[2]s Example Description"
    bg_opacity_percent = 100
    card_color         = "#000000"
    text_color         = "#eff0f1"

    fields = [
      {
        type       = "Directory Attribute"
        title      = "givenName"
        attribute  = "name.given"
        is_visible = false
      },
    ]
  }
}
data "pingone_credential_issuance_rule" "%[2]s" {
  environment_id              = data.pingone_environment.general_test.id
  credential_type_id          = resource.pingone_credential_type.%[2]s.id
  credential_issuance_rule_id = "9c052a8a-14be-44e4-8f07-2662569994ce" // dummy ID that conforms to UUID v4

}`, acctest.GenericSandboxEnvironment(), resourceName)
}

func testAccCredentialIssuanceRuleDataSource_NoEnvironmentID(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s
resource "pingone_credential_type" "%[2]s" {
  environment_id       = data.pingone_environment.general_test.id
  title                = "%[2]s"
  description          = "%[2]s Example Description"
  card_type            = "%[2]s"
  card_design_template = "<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 740 480\"><rect fill=\"none\" width=\"736\" height=\"476\" stroke=\"#CACED3\" stroke-width=\"3\" rx=\"10\" ry=\"10\" x=\"2\" y=\"2\"></rect><rect fill=\"$${cardColor}\" height=\"476\" rx=\"10\" ry=\"10\" width=\"736\" x=\"2\" y=\"2\" opacity=\"$${bgOpacityPercent}\"></rect><line y2=\"160\" x2=\"695\" y1=\"160\" x1=\"42.5\" stroke=\"$${textColor}\"></line><text fill=\"$${textColor}\" font-weight=\"450\" font-size=\"30\" x=\"160\" y=\"90\">$${cardTitle}</text><text fill=\"$${textColor}\" font-size=\"25\" font-weight=\"300\" x=\"160\" y=\"130\">$${cardSubtitle}</text></svg>"

  metadata = {
    name               = "%[2]s"
    description        = "%[2]s Example Description"
    bg_opacity_percent = 20
    card_color         = "#000000"
    text_color         = "#eff0f1"

    fields = [
      {
        type       = "Directory Attribute"
        title      = "givenName"
        attribute  = "name.given"
        is_visible = false
      },
    ]
  }
}
data "pingone_credential_issuance_rule" "%[2]s" {
  credential_type_id          = resource.pingone_credential_type.%[2]s.id
  credential_issuance_rule_id = "9c052a8a-14be-44e4-8f07-2662569994ce" // dummy ID that conforms to UUID v4

}`, acctest.GenericSandboxEnvironment(), resourceName)
}

func testAccCredentialIssuanceRuleDataSource_NoCredentialIssuanceRuleID(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s
resource "pingone_credential_type" "%[2]s" {
  environment_id       = data.pingone_environment.general_test.id
  title                = "%[2]s"
  description          = "%[2]s Example Description"
  card_type            = "%[2]s"
  card_design_template = "<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 740 480\"><rect fill=\"none\" width=\"736\" height=\"476\" stroke=\"#CACED3\" stroke-width=\"3\" rx=\"10\" ry=\"10\" x=\"2\" y=\"2\"></rect><rect fill=\"$${cardColor}\" height=\"476\" rx=\"10\" ry=\"10\" width=\"736\" x=\"2\" y=\"2\" opacity=\"$${bgOpacityPercent}\"></rect><line y2=\"160\" x2=\"695\" y1=\"160\" x1=\"42.5\" stroke=\"$${textColor}\"></line><text fill=\"$${textColor}\" font-weight=\"450\" font-size=\"30\" x=\"160\" y=\"90\">$${cardTitle}</text><text fill=\"$${textColor}\" font-size=\"25\" font-weight=\"300\" x=\"160\" y=\"130\">$${cardSubtitle}</text></svg>"

  metadata = {
    name               = "%[2]s"
    description        = "%[2]s Example Description"
    bg_opacity_percent = 100
    card_color         = "#000000"
    text_color         = "#eff0f1"

    fields = [
      {
        type       = "Directory Attribute"
        title      = "givenName"
        attribute  = "name.given"
        is_visible = false
      },
    ]
  }
}
data "pingone_credential_issuance_rule" "%[2]s" {
  environment_id     = data.pingone_environment.general_test.id
  credential_type_id = resource.pingone_credential_type.%[2]s.id

}`, acctest.GenericSandboxEnvironment(), resourceName)
}
