package sso_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/sso"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccApplicationSecretDataSource_Basic(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application_secret.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {
				VersionConstraint: "0.11.1",
				Source:            "hashicorp/time",
			},
		},
		CheckDestroy: sso.ApplicationSecret_CheckDestroy,
		ErrorCheck:   acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationSecretDataSourceConfig_Full(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "application_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "previous.secret", regexp.MustCompile(`[a-zA-Z0-9-~_]{10,}`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "previous.expires_at", verify.RFC3339Regexp),
					resource.TestCheckNoResourceAttr(dataSourceFullName, "previous.last_used"),
					resource.TestMatchResourceAttr(dataSourceFullName, "secret", regexp.MustCompile(`[a-zA-Z0-9-~_]{10,}`)),
				),
			},
		},
	})
}

func TestAccApplicationSecretDataSource_Rotation(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application_secret.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {
				VersionConstraint: "0.11.1",
				Source:            "hashicorp/time",
			},
		},
		CheckDestroy: sso.ApplicationSecret_CheckDestroy,
		ErrorCheck:   acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationSecretDataSourceConfig_Rotation1(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr(dataSourceFullName, "previous.secret"),
					resource.TestCheckNoResourceAttr(dataSourceFullName, "previous.expires_at"),
					resource.TestCheckNoResourceAttr(dataSourceFullName, "previous.last_used"),
					resource.TestMatchResourceAttr(dataSourceFullName, "secret", regexp.MustCompile(`[a-zA-Z0-9-~_]{10,}`)),
				),
			},
			{
				Config: testAccApplicationSecretDataSourceConfig_Rotation2(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "previous.secret", regexp.MustCompile(`[a-zA-Z0-9-~_]{10,}`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "previous.expires_at", verify.RFC3339Regexp),
					resource.TestCheckNoResourceAttr(dataSourceFullName, "previous.last_used"),
					resource.TestMatchResourceAttr(dataSourceFullName, "secret", regexp.MustCompile(`[a-zA-Z0-9-~_]{10,}`)),
				),
			},
		},
	})
}

func TestAccApplicationSecretDataSource_IncorrectApplicationType(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.ApplicationSecret_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccApplicationSecretDataSourceConfig_IncorrectApplicationType(resourceName, name),
				ExpectError: regexp.MustCompile("Error when calling `ReadApplicationSecret`: The request could not be completed. The requested resource was not found."),
			},
		},
	})
}

func TestAccApplicationSecretDataSource_NotFound(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.ApplicationSecret_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccApplicationSecretDataSourceConfig_NotFoundByID(resourceName),
				ExpectError: regexp.MustCompile("Error when calling `ReadApplicationSecret`: The request could not be completed. The requested resource was not found."),
			},
		},
	})
}

func testAccApplicationSecretDataSourceConfig_Full(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true

  oidc_options = {
    type                        = "SINGLE_PAGE_APP"
    grant_types                 = ["AUTHORIZATION_CODE"]
    response_types              = ["CODE"]
    token_endpoint_authn_method = "CLIENT_SECRET_BASIC"
    redirect_uris               = ["https://www.pingidentity.com"]
  }
}

resource "time_offset" "%[2]s" {
	offset_minutes = 10
  }

resource "pingone_application_secret" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  application_id = pingone_application.%[2]s.id

  previous = {
	expires_at = time_offset.%[2]s.rfc3339
  }
}

data "pingone_application_secret" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  application_id = pingone_application.%[2]s.id

  depends_on = [
	pingone_application_secret.%[2]s,
  ]
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccApplicationSecretDataSourceConfig_Rotation1(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true

  oidc_options = {
    type                        = "SINGLE_PAGE_APP"
    grant_types                 = ["AUTHORIZATION_CODE"]
    response_types              = ["CODE"]
    token_endpoint_authn_method = "CLIENT_SECRET_BASIC"
    redirect_uris               = ["https://www.pingidentity.com"]
  }
}

resource "pingone_application_secret" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  application_id = pingone_application.%[2]s.id
}

data "pingone_application_secret" "%[2]s" {
	environment_id = data.pingone_environment.general_test.id
	application_id = pingone_application.%[2]s.id

	depends_on = [
		pingone_application_secret.%[2]s,
	]
  }`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccApplicationSecretDataSourceConfig_Rotation2(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true

  oidc_options = {
    type                        = "SINGLE_PAGE_APP"
    grant_types                 = ["AUTHORIZATION_CODE"]
    response_types              = ["CODE"]
    token_endpoint_authn_method = "CLIENT_SECRET_BASIC"
    redirect_uris               = ["https://www.pingidentity.com"]
  }
}

resource "time_offset" "%[2]s" {
	offset_minutes = 10
  }

resource "pingone_application_secret" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  application_id = pingone_application.%[2]s.id

  previous = {
	expires_at = time_offset.%[2]s.rfc3339
  }
}

data "pingone_application_secret" "%[2]s" {
	environment_id = data.pingone_environment.general_test.id
	application_id = pingone_application.%[2]s.id

	depends_on = [
		pingone_application_secret.%[2]s,
	]
  }`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccApplicationSecretDataSourceConfig_IncorrectApplicationType(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_key" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name                = "%[3]s"
  algorithm           = "EC"
  key_length          = 256
  signature_algorithm = "SHA384withECDSA"
  subject_dn          = "CN=%[3]s, OU=Ping Identity, O=Ping Identity, L=, ST=, C=US"
  usage_type          = "SIGNING"
  validity_period     = 365
}

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true

  saml_options = {
    acs_urls           = ["https://pingidentity.com"]
    assertion_duration = 3600
    sp_entity_id       = "sp:entity:%[2]s"

    idp_signing_key = {
      key_id    = pingone_key.%[2]s.id
      algorithm = pingone_key.%[2]s.signature_algorithm
    }
  }
}

data "pingone_application_secret" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  application_id = pingone_application.%[2]s.id
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccApplicationSecretDataSourceConfig_NotFoundByID(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_application_secret" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  application_id = "9c052a8a-14be-44e4-8f07-2662569994ce" // dummy ID that conforms to UUID v4
}
`, acctest.GenericSandboxEnvironment(), resourceName)
}
