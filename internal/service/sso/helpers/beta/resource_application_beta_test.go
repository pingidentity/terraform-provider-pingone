// Copyright © 2025 Ping Identity Corporation

// This file relates to a beta feature described in CDI-492

//go:build beta

package beta_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/sso"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccApplication_OIDC_GeneratedClientID(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckFeatureFlag(t, acctest.ENUMFEATUREFLAG_DAVINCI) // its a hack, because I'm not sure this is the best way to do this
			acctest.PreCheckBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Application_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationConfig_OIDC_MinimalCustom(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.client_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.initial_client_secret"),
				),
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

func TestAccApplication_OIDC_ImportedClientIDClientSecret(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_application.%s", resourceName)

	name := resourceName

	clientIDValue := "imported-client-id"
	clientSecretValue1 := "imported-client-secret"
	clientSecretValue2 := "changed-client-secret"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckFeatureFlag(t, acctest.ENUMFEATUREFLAG_DAVINCI) // its a hack, because I'm not sure this is the best way to do this
			acctest.PreCheckBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Application_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationConfig_OIDC_Import(resourceName, name, &clientIDValue, &clientSecretValue1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.client_id", clientIDValue),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.initial_client_secret", clientSecretValue1),
				),
			},
			{
				Config: testAccApplicationConfig_OIDC_Import(resourceName, name, &clientIDValue, &clientSecretValue2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.client_id", clientIDValue),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.initial_client_secret", clientSecretValue2),
				),
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
			{
				Config:  testAccApplicationConfig_OIDC_Import(resourceName, name, &clientIDValue, &clientSecretValue1),
				Destroy: true,
			},
			{
				Config: testAccApplicationConfig_OIDC_Import(resourceName, name, &clientIDValue, nil),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.client_id", clientIDValue),
					resource.TestCheckNoResourceAttr(resourceFullName, "oidc_options.initial_client_secret"),
				),
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
			{
				Config:  testAccApplicationConfig_OIDC_Import(resourceName, name, &clientIDValue, nil),
				Destroy: true,
			},
			{
				Config: testAccApplicationConfig_OIDC_Import(resourceName, name, nil, &clientSecretValue1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "oidc_options.client_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "oidc_options.initial_client_secret", clientSecretValue1),
				),
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

func testAccApplicationConfig_OIDC_MinimalCustom(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.app_import_ff_test.id
  name           = "%[3]s"
  enabled        = true

  oidc_options = {
    type                       = "CUSTOM_APP"
    grant_types                = ["CLIENT_CREDENTIALS"]
    token_endpoint_auth_method = "CLIENT_SECRET_BASIC"
  }
}
`, acctest.AppImportFFSandboxEnvironment(), resourceName, name)
}

func testAccApplicationConfig_OIDC_Import(resourceName, name string, clientID, clientSecret *string) string {

	interpolatedClientID := ""
	if clientID != nil {
		interpolatedClientID = fmt.Sprintf("client_id = \"%s\"", *clientID)
	}

	interpolatedClientSecret := ""
	if clientSecret != nil {
		interpolatedClientSecret = fmt.Sprintf("initial_client_secret = \"%s\"", *clientSecret)
	}

	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.app_import_ff_test.id
  name           = "%[3]s"
  enabled        = true

  oidc_options = {
    type                       = "WEB_APP"
    grant_types                = ["REFRESH_TOKEN", "AUTHORIZATION_CODE"]
    response_types             = ["CODE"]
    token_endpoint_auth_method = "CLIENT_SECRET_BASIC"
    redirect_uris              = ["https://www.pingidentity.com"]
    %[4]s
    %[5]s
  }
}
`, acctest.AppImportFFSandboxEnvironment(), resourceName, name, interpolatedClientID, interpolatedClientSecret)
}
