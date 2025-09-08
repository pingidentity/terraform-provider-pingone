// Copyright © 2025 Ping Identity Corporation

// This file relates to a beta feature described in CDI-492

//go:build beta

package beta_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/sso"
)

func TestAccApplicationSecret_ImportedApp_Rotation(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application_secret.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.pingone_application_secret.%s", resourceName)

	name := resourceName

	clientIDValue := fmt.Sprintf("imported-client-id-%s", resourceName)
	clientSecretValue := "imported-client-secret"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckFeatureFlag(t, acctest.ENUMFEATUREFLAG_APP_AND_RESOURCE_IMPORT)
			acctest.PreCheckBeta(t)
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
			// Single from new
			{
				Config: testAccApplicationSecretConfig_ImportedApp_Rotation1(resourceName, name, clientIDValue, clientSecretValue),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr(dataSourceFullName, "previous.secret"),
					resource.TestCheckNoResourceAttr(dataSourceFullName, "previous.expires_at"),
					resource.TestCheckNoResourceAttr(dataSourceFullName, "previous.last_used"),
					resource.TestCheckResourceAttr(dataSourceFullName, "secret", clientSecretValue),
				),
			},
			{
				Config: testAccApplicationSecretConfig_ImportedApp_Rotation2(resourceName, name, clientIDValue, clientSecretValue),
				Check: resource.ComposeTestCheckFunc(
					// resource.TestCheckResourceAttr(resourceFullName, "previous.secret", clientSecretValue),
					// resource.TestMatchResourceAttr(resourceFullName, "previous.expires_at", verify.RFC3339Regexp),
					resource.TestCheckNoResourceAttr(resourceFullName, "previous.secret"),
					resource.TestCheckNoResourceAttr(resourceFullName, "previous.expires_at"),
					resource.TestCheckNoResourceAttr(resourceFullName, "previous.last_used"),
					resource.TestMatchResourceAttr(resourceFullName, "secret", regexp.MustCompile(`[a-zA-Z0-9-~_]{10,}`)),
				),
			},
			{
				Config:  testAccApplicationSecretConfig_ImportedApp_Rotation2(resourceName, name, clientIDValue, clientSecretValue),
				Destroy: true,
			},
			{
				Config: testAccApplicationSecretConfig_ImportedApp_Rotation2(resourceName, name, clientIDValue, clientSecretValue),
				Check: resource.ComposeTestCheckFunc(
					// resource.TestCheckResourceAttr(resourceFullName, "previous.secret", clientSecretValue),
					// resource.TestMatchResourceAttr(resourceFullName, "previous.expires_at", verify.RFC3339Regexp),
					resource.TestCheckNoResourceAttr(resourceFullName, "previous.secret"),
					resource.TestCheckNoResourceAttr(resourceFullName, "previous.expires_at"),
					resource.TestCheckNoResourceAttr(resourceFullName, "previous.last_used"),
					resource.TestMatchResourceAttr(resourceFullName, "secret", regexp.MustCompile(`[a-zA-Z0-9-~_]{10,}`)),
				),
			},
		},
	})
}

func testAccApplicationSecretConfig_ImportedApp_Rotation1(resourceName, name, clientID, clientSecret string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.app_import_ff_test.id
  name           = "%[3]s"
  enabled        = true

  oidc_options = {
    type                       = "WEB_APP"
    grant_types                = ["AUTHORIZATION_CODE"]
    response_types             = ["CODE"]
    token_endpoint_auth_method = "CLIENT_SECRET_BASIC"
    redirect_uris              = ["https://www.pingidentity.com"]

	client_id             = "%[4]s"
	initial_client_secret = "%[5]s"
  }
}

data "pingone_application_secret" "%[2]s" {
  environment_id = data.pingone_environment.app_import_ff_test.id
  application_id = pingone_application.%[2]s.id
}`, acctest.AppImportFFSandboxEnvironment(), resourceName, name, clientID, clientSecret)
}

func testAccApplicationSecretConfig_ImportedApp_Rotation2(resourceName, name, clientID, clientSecret string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.app_import_ff_test.id
  name           = "%[3]s"
  enabled        = true

  oidc_options = {
    type                       = "WEB_APP"
    grant_types                = ["AUTHORIZATION_CODE"]
    response_types             = ["CODE"]
    token_endpoint_auth_method = "CLIENT_SECRET_BASIC"
    redirect_uris              = ["https://www.pingidentity.com"]

	client_id             = "%[4]s"
	initial_client_secret = "%[5]s"
  }
}

resource "time_offset" "%[2]s" {
  offset_minutes = 10
}

resource "pingone_application_secret" "%[2]s" {
  environment_id = data.pingone_environment.app_import_ff_test.id
  application_id = pingone_application.%[2]s.id
}`, acctest.AppImportFFSandboxEnvironment(), resourceName, name, clientID, clientSecret)
}
