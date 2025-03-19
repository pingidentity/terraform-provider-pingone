// Copyright © 2025 Ping Identity Corporation

package sso_test

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/base"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/sso"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
)

func TestAccApplicationV2_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_applicationv2.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	var applicationID, environmentID string

	var p1Client *client.Client
	var ctx = context.Background()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)

			p1Client = acctest.PreCheckTestClient(ctx, t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Application_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccApplicationV2Config_OIDC_MinimalWeb(resourceName, name),
				Check:  sso.Application_GetIDs(resourceFullName, &environmentID, &applicationID),
			},
			{
				PreConfig: func() {
					sso.Application_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID, applicationID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the environment
			{
				Config: testAccApplicationV2Config_NewEnv(environmentName, licenseID, resourceName, name),
				Check:  sso.Application_GetIDs(resourceFullName, &environmentID, &applicationID),
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

func TestAccApplicationV2_OIDCWebUpdate(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	name := resourceName

	data, _ := os.ReadFile("../../acctest/test_assets/image/image-logo.gif")
	image := base64.StdEncoding.EncodeToString(data)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Application_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationV2Config_OIDC_MinimalWeb(resourceName, name),
			},
			{
				Config: testAccApplicationV2Config_OIDC_FullWeb(resourceName, name, image),
			},
			{
				Config: testAccApplicationV2Config_OIDC_MinimalWeb(resourceName, name),
			},
		},
	})
}

func TestAccApplicationV2_SAMLFull(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_applicationv2.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	data, _ := os.ReadFile("../../acctest/test_assets/image/image-logo.gif")
	image := base64.StdEncoding.EncodeToString(data)

	pem_cert := os.Getenv("PINGONE_KEY_PEM_CERT")
	pkcs7_cert := os.Getenv("PINGONE_KEY_PKCS7_CERT")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
			acctest.PreCheckPKCS7Cert(t)
			acctest.PreCheckPEMCert(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Application_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationV2Config_SAML_Full(environmentName, licenseID, resourceName, name, image, pkcs7_cert, pem_cert),
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

func TestAccApplicationV2_SAMLMinimal(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	data, _ := os.ReadFile("../../acctest/test_assets/image/image-logo.gif")
	image := base64.StdEncoding.EncodeToString(data)

	pem_cert := os.Getenv("PINGONE_KEY_PEM_CERT")
	pkcs7_cert := os.Getenv("PINGONE_KEY_PKCS7_CERT")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
			acctest.PreCheckPKCS7Cert(t)
			acctest.PreCheckPEMCert(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Application_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationV2Config_SAML_Minimal(environmentName, licenseID, resourceName, name, image, pkcs7_cert, pem_cert),
			},
		},
	})
}

func TestAccApplicationV2_ExternalLinkFull(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_applicationv2.%s", resourceName)

	name := resourceName

	data, _ := os.ReadFile("../../acctest/test_assets/image/image-logo.gif")
	image := base64.StdEncoding.EncodeToString(data)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Application_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationV2Config_ExternalLinkFull(resourceName, name, image),
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

func TestAccApplicationV2_ExternalLinkMinimal(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.Application_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationV2Config_ExternalLinkMinimal(resourceName, name),
			},
		},
	})
}

func testAccApplicationV2Config_OIDC_MinimalWeb(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_applicationv2" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true

  protocol = "OPENID_CONNECT"
    type                       = "WEB_APP"
    grant_types                = ["REFRESH_TOKEN", "AUTHORIZATION_CODE"]
    response_types             = ["CODE"]
    token_endpoint_auth_method = "CLIENT_SECRET_BASIC"
    redirect_uris              = ["https://www.pingidentity.com"]
  
}
`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccApplicationV2Config_NewEnv(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_applicationv2" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s"
  enabled        = true

  protocol = "OPENID_CONNECT"
    type                       = "WEB_APP"
    grant_types                = ["AUTHORIZATION_CODE", "REFRESH_TOKEN"]
    response_types             = ["CODE"]
    token_endpoint_auth_method = "CLIENT_SECRET_BASIC"
    redirect_uris              = ["https://www.pingidentity.com"]
}
`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccApplicationV2Config_OIDC_FullWeb(resourceName, name, image string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_group" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
}

resource "pingone_image" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  image_file_base64 = "%[4]s"
}

resource "pingone_applicationv2" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "My test OIDC app"
  login_page_url = "https://www.pingidentity.com"

  icon = {
    id   = pingone_image.%[2]s.id
    href = pingone_image.%[2]s.uploaded_image.href
  }

  access_control = {
    group = {
      type = "ANY_GROUP"

      groups = [
        {
          id = pingone_group.%[2]s.id
        }
      ]
    }
    role = {
      type = "ADMIN_USERS_ONLY"
    }
  }

  hidden_from_app_portal = true


  enabled = true
  protocol = "OPENID_CONNECT"
    type                            = "WEB_APP"
    grant_types                     = ["REFRESH_TOKEN", "AUTHORIZATION_CODE"]
    response_types                  = ["CODE"]
    token_endpoint_auth_method      = "CLIENT_SECRET_BASIC"
    redirect_uris                   = ["https://www.pingidentity.com", "https://pingidentity.com"]
    allow_wildcard_in_redirect_uris = true
    post_logout_redirect_uris       = ["https://www.pingidentity.com/logout", "https://pingidentity.com/logout"]

    refresh_token_duration                             = 3000000
	# Current application resource sets a default for this which is necessary, because it can't be nullified
    # refresh_token_rolling_duration                     = 30000000
    refresh_token_rolling_grace_period_duration        = 80000
    additional_refresh_token_replay_protection_enabled = false

    home_page_url      = "https://www.pingidentity.com"
    initiate_login_uri = "https://www.pingidentity.com/initiate"
    target_link_uri    = "https://www.pingidentity.com/target"
    pkce_enforcement   = "OPTIONAL"

    par_requirement = "OPTIONAL"
    par_timeout     = 60

    support_unsigned_request_object = true

    cors_settings = {
      behavior = "ALLOW_SPECIFIC_ORIGINS"
      origins = [
        "http://localhost",
        "https://localhost",
        "http://auth.pingidentity.com",
        "https://auth.pingidentity.com",
        "http://*.pingidentity.com",
        "https://*.pingidentity.com",
        "http://192.168.1.1",
        "https://192.168.1.1",
      ]
    }
}
`, acctest.GenericSandboxEnvironment(), resourceName, name, image)
}

func testAccApplicationV2Config_SAML_Full(environmentName, licenseID, resourceName, name, image, pkcs7_cert, pem_cert string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_group" "%[3]s-1" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s-1"
}

resource "pingone_group" "%[3]s-2" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s-2"
}

resource "pingone_key" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  name                = "%[4]s"
  algorithm           = "EC"
  key_length          = 256
  signature_algorithm = "SHA384withECDSA"
  subject_dn          = "CN=%[4]s, OU=Ping Identity, O=Ping Identity, L=, ST=, C=US"
  usage_type          = "SIGNING"
  validity_period     = 365
}

resource "pingone_image" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  image_file_base64 = "%[5]s"
}

resource "pingone_certificate" "%[3]s-1" {
  environment_id = pingone_environment.%[2]s.id

  pkcs7_file_base64 = <<EOT
%[6]s
EOT

  usage_type = "SIGNING"
}

resource "pingone_certificate" "%[3]s-2" {
  environment_id = pingone_environment.%[2]s.id

  pem_file = <<EOT
%[7]s
EOT

  usage_type = "SIGNING"
}

resource "pingone_certificate" "%[3]s-enc" {
  environment_id = pingone_environment.%[2]s.id
  pem_file       = <<EOT
%[7]s
EOT
  usage_type     = "ENCRYPTION"
}

resource "pingone_applicationv2" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s"
  description    = "My test SAML app"
  login_page_url = "https://www.pingidentity.com"

  icon = {
    id   = pingone_image.%[3]s.id
    href = pingone_image.%[3]s.uploaded_image.href
  }

  access_control = {
    group = {
      type = "ANY_GROUP"

      groups = [
        {
          id = pingone_group.%[3]s-2.id
        },
		{
		  id = pingone_group.%[3]s-1.id
		}
      ]
    }
    role = {
      type = "ADMIN_USERS_ONLY"
    }
  }

  hidden_from_app_portal = true

  enabled = true

  protocol = "SAML"
    type               = "WEB_APP"
    home_page_url      = "https://www.pingidentity.com"
    acs_urls           = ["https://www.pingidentity.com", "https://pingidentity.com"]
    assertion_duration = 3600
    sp_entity_id       = "sp:entity:%[3]s"

    sp_encryption = {
      algorithm = "AES_256"
      certificate = {
        id = pingone_certificate.%[3]s-enc.id
      }
    }

    assertion_signed = false
	idp_signing = {
      key = {
        id   = pingone_key.%[3]s.id
      }
      algorithm = pingone_key.%[3]s.signature_algorithm
    }
    enable_requested_authn_context   = true
    name_id_format                    = "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress"
    response_signed               = true
    session_not_on_or_after_duration = 64
    slo_binding                      = "HTTP_REDIRECT"
    slo_endpoint                     = "https://www.pingidentity.com/sloendpoint"
    slo_response_endpoint            = "https://www.pingidentity.com/sloresponseendpoint"
    slo_window                       = 3

    default_target_url = "https://www.pingidentity.com/relaystate"

	sp_verification = {
      authn_request_signed = true
      certificates = [
        {
          id = pingone_certificate.%[3]s-2.id
        },
        {
          id = pingone_certificate.%[3]s-1.id
        }
      ]
    }

    cors_settings = {
      behavior = "ALLOW_SPECIFIC_ORIGINS"
      origins = [
        "http://localhost",
        "https://localhost",
        "http://auth.pingidentity.com",
        "https://auth.pingidentity.com",
        "http://*.pingidentity.com",
        "https://*.pingidentity.com",
        "http://192.168.1.1",
        "https://192.168.1.1",
      ]
    }
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name, image, pkcs7_cert, pem_cert)
}

func testAccApplicationV2Config_SAML_Minimal(environmentName, licenseID, resourceName, name, image, pkcs7_cert, pem_cert string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_group" "%[3]s-1" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s-1"
}

resource "pingone_group" "%[3]s-2" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s-2"
}

resource "pingone_key" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  name                = "%[4]s"
  algorithm           = "EC"
  key_length          = 256
  signature_algorithm = "SHA384withECDSA"
  subject_dn          = "CN=%[4]s, OU=Ping Identity, O=Ping Identity, L=, ST=, C=US"
  usage_type          = "SIGNING"
  validity_period     = 365
}

resource "pingone_image" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  image_file_base64 = "%[5]s"
}

resource "pingone_certificate" "%[3]s-1" {
  environment_id = pingone_environment.%[2]s.id

  pkcs7_file_base64 = <<EOT
%[6]s
EOT

  usage_type = "SIGNING"
}

resource "pingone_certificate" "%[3]s-2" {
  environment_id = pingone_environment.%[2]s.id

  pem_file = <<EOT
%[7]s
EOT

  usage_type = "SIGNING"
}

resource "pingone_certificate" "%[3]s-enc" {
  environment_id = pingone_environment.%[2]s.id
  pem_file       = <<EOT
%[7]s
EOT
  usage_type     = "ENCRYPTION"
}

resource "pingone_applicationv2" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s"
  enabled        = true

  protocol = "SAML"
    type               = "WEB_APP"
    acs_urls           = ["https://pingidentity.com"]
    assertion_duration = 3600
    sp_entity_id       = "sp:entity:%[3]s"

	idp_signing = {
		key = {
		  id   = pingone_key.%[3]s.id
		}
		algorithm = pingone_key.%[3]s.signature_algorithm
	  }
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name, image, pkcs7_cert, pem_cert)
}

func testAccApplicationV2Config_ExternalLinkFull(resourceName, name, image string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_group" "%[2]s-1" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-1"
}

resource "pingone_group" "%[2]s-2" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-2"
}

resource "pingone_image" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  image_file_base64 = "%[4]s"
}

resource "pingone_applicationv2" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "My test external link app"

  icon = {
    id   = pingone_image.%[2]s.id
    href = pingone_image.%[2]s.uploaded_image.href
  }

  access_control = {
     group = {
       type = "ANY_GROUP"

       groups = [
         {
           id = pingone_group.%[2]s-2.id
         },
         {
           id = pingone_group.%[2]s-1.id
         }
       ]
     }
     role = {
       type = "ADMIN_USERS_ONLY"
     }
   }

  hidden_from_app_portal = true

  enabled = true

  protocol = "EXTERNAL_LINK"
 type = "PORTAL_LINK_APP"
    home_page_url = "https://www.pingidentity.com"
}`, acctest.GenericSandboxEnvironment(), resourceName, name, image)
}

func testAccApplicationV2Config_ExternalLinkMinimal(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_applicationv2" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true

  protocol = "EXTERNAL_LINK"
 type = "PORTAL_LINK_APP"
    home_page_url = "https://www.pingidentity.com"
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
