// Copyright © 2025 Ping Identity Corporation

package mfa_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	acctestlegacysdk "github.com/pingidentity/terraform-provider-pingone/internal/acctest/legacysdk"
	baselegacysdk "github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/base/legacysdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/mfa"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/sso"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccApplicationPushCredential_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_application_push_credential.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	firebaseCredentials := os.Getenv("PINGONE_GOOGLE_FIREBASE_CREDENTIALS")

	var applicationPushCredentialID, applicationID, environmentID string

	var p1Client *client.Client
	var ctx = context.Background()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoBeta(t)
			p1Client = acctestlegacysdk.PreCheckTestClient(ctx, t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.ApplicationPushCredential_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Test removal of the resource
			{
				Config: testAccApplicationPushCredentialConfig_FCM(resourceName, name, firebaseCredentials),
				Check:  mfa.ApplicationPushCredential_GetIDs(resourceFullName, &environmentID, &applicationID, &applicationPushCredentialID),
			},
			{
				PreConfig: func() {
					mfa.ApplicationPushCredential_RemovalDrift_PreConfig(ctx, p1Client.API.MFAAPIClient, t, environmentID, applicationID, applicationPushCredentialID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the application
			{
				Config: testAccApplicationPushCredentialConfig_FCM(resourceName, name, firebaseCredentials),
				Check:  mfa.ApplicationPushCredential_GetIDs(resourceFullName, &environmentID, &applicationID, &applicationPushCredentialID),
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
				Config: testAccApplicationPushCredentialConfig_NewEnv(environmentName, licenseID, resourceName, name, firebaseCredentials),
				Check:  mfa.ApplicationPushCredential_GetIDs(resourceFullName, &environmentID, &applicationID, &applicationPushCredentialID),
			},
			{
				PreConfig: func() {
					baselegacysdk.Environment_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccApplicationPushCredential_FCM(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_application_push_credential.%s", resourceName)

	name := resourceName

	firebaseCredentials := os.Getenv("PINGONE_GOOGLE_FIREBASE_CREDENTIALS")

	fullFCMCheck := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "application_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttrSet(resourceFullName, "fcm.google_service_account_credentials"),
		resource.TestCheckResourceAttr(resourceFullName, "apns.%", "0"),
		resource.TestCheckResourceAttr(resourceFullName, "hms.%", "0"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
			acctest.PreCheckGoogleFirebaseCredentials(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.ApplicationPushCredential_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// FCM new
			{
				Config: testAccApplicationPushCredentialConfig_FCM(resourceName, name, firebaseCredentials),
				Check:  fullFCMCheck,
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

						return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["application_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"fcm",
					"fcm.google_service_account_credentials",
				},
			},
		},
	})
}

func TestAccApplicationPushCredential_APNS(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_application_push_credential.%s", resourceName)

	name := resourceName

	pkcs8Key := os.Getenv("PINGONE_KEY_PKCS8")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
			acctest.PreCheckAPNSPKCS8Key(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.ApplicationPushCredential_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationPushCredentialConfig_APNS(resourceName, name, pkcs8Key),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "application_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "fcm.%", "0"),
					resource.TestCheckResourceAttrSet(resourceFullName, "apns.key"),
					resource.TestCheckResourceAttrSet(resourceFullName, "apns.team_id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "apns.token_signing_key"),
					resource.TestCheckResourceAttr(resourceFullName, "hms.%", "0"),
				),
			},
			{
				Config: testAccApplicationPushCredentialConfig_APNS(resourceName, name, pkcs8Key),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "application_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "fcm.%", "0"),
					resource.TestCheckResourceAttrSet(resourceFullName, "apns.key"),
					resource.TestCheckResourceAttrSet(resourceFullName, "apns.team_id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "apns.token_signing_key"),
					resource.TestCheckResourceAttr(resourceFullName, "hms.%", "0"),
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

						return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["application_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"apns",
					"apns.key",
					"apns.team_id",
					"apns.token_signing_key",
				},
			},
		},
	})
}

func TestAccApplicationPushCredential_HMS(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_application_push_credential.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.ApplicationPushCredential_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationPushCredentialConfig_HMS(resourceName, name, fmt.Sprintf("%s1", name)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "application_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "fcm.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "apns.%", "0"),
					resource.TestCheckResourceAttrSet(resourceFullName, "hms.client_id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "hms.client_secret"),
				),
			},
			{
				Config: testAccApplicationPushCredentialConfig_HMS(resourceName, name, fmt.Sprintf("%s2", name)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "application_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "fcm.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "apns.%", "0"),
					resource.TestCheckResourceAttrSet(resourceFullName, "hms.client_id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "hms.client_secret"),
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

						return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["application_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"hms",
					"hms.client_id",
					"hms.client_secret",
				},
			},
		},
	})
}

func TestAccApplicationPushCredential_Change(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_application_push_credential.%s", resourceName)

	name := resourceName

	firebaseCredentials := os.Getenv("PINGONE_GOOGLE_FIREBASE_CREDENTIALS")
	pkcs8Key := os.Getenv("PINGONE_KEY_PKCS8")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
			acctest.PreCheckGoogleFirebaseCredentials(t)
			acctest.PreCheckAPNSPKCS8Key(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.ApplicationPushCredential_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationPushCredentialConfig_FCM(resourceName, name, firebaseCredentials),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "application_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttrSet(resourceFullName, "fcm.google_service_account_credentials"),
					resource.TestCheckResourceAttr(resourceFullName, "apns.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "hms.%", "0"),
				),
			},
			{
				Config: testAccApplicationPushCredentialConfig_APNS(resourceName, name, pkcs8Key),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "application_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "fcm.%", "0"),
					resource.TestCheckResourceAttrSet(resourceFullName, "apns.key"),
					resource.TestCheckResourceAttrSet(resourceFullName, "apns.team_id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "apns.token_signing_key"),
					resource.TestCheckResourceAttr(resourceFullName, "hms.%", "0"),
				),
			},
			{
				Config: testAccApplicationPushCredentialConfig_HMS(resourceName, name, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "application_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "fcm.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "apns.%", "0"),
					resource.TestCheckResourceAttrSet(resourceFullName, "hms.client_id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "hms.client_secret"),
				),
			},
		},
	})
}

func TestAccApplicationPushCredential_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_mfa_application_push_credential.%s", resourceName)

	name := resourceName

	firebaseCredentials := os.Getenv("PINGONE_GOOGLE_FIREBASE_CREDENTIALS")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
			acctest.PreCheckGoogleFirebaseCredentials(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             mfa.ApplicationPushCredential_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccApplicationPushCredentialConfig_FCM(resourceName, name, firebaseCredentials),
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
				ImportStateId: "badformat/badformat",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`Unexpected Import Identifier`),
			},
		},
	})
}

func testAccApplicationPushCredentialConfig_NewEnv(environmentName, licenseID, resourceName, name, key string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s"
  description    = "My test OIDC app for MFA Policy"

  login_page_url = "https://www.pingidentity.com"

  enabled = true

  oidc_options = {
    type                       = "NATIVE_APP"
    grant_types                = ["CLIENT_CREDENTIALS"]
    token_endpoint_auth_method = "CLIENT_SECRET_BASIC"

    mobile_app = {
      package_name = "com.%[4]s.package"

      passcode_refresh_seconds = 45

      integrity_detection = {
        enabled = true
        cache_duration = {
          amount = 30
          units  = "HOURS"
        }
        google_play = {
          verification_type = "INTERNAL"
          decryption_key    = "dummykeydoesnotexist"
          verification_key  = "dummykeydoesnotexist"
        }
      }
    }

    package_name = "com.%[3]s.package"
  }
}

resource "pingone_mfa_application_push_credential" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  application_id = pingone_application.%[3]s.id

  fcm = {
    google_service_account_credentials = jsonencode(%[5]s)
  }
}`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name, key)
}

func testAccApplicationPushCredentialConfig_FCM(resourceName, name, key string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "My test OIDC app for MFA Policy"

  login_page_url = "https://www.pingidentity.com"

  enabled = true

  oidc_options = {
    type                       = "NATIVE_APP"
    grant_types                = ["CLIENT_CREDENTIALS"]
    token_endpoint_auth_method = "CLIENT_SECRET_BASIC"

    mobile_app = {
      package_name = "com.%[2]s.package"

      passcode_refresh_seconds = 45

      integrity_detection = {
        enabled = true
        cache_duration = {
          amount = 30
          units  = "HOURS"
        }
        google_play = {
          verification_type = "INTERNAL"
          decryption_key    = "dummykeydoesnotexist"
          verification_key  = "dummykeydoesnotexist"
        }
      }
    }

    package_name = "com.%[2]s.package"
  }
}

resource "pingone_mfa_application_push_credential" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  application_id = pingone_application.%[2]s.id

  fcm = {
    google_service_account_credentials = jsonencode(%[4]s)
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name, key)
}

func testAccApplicationPushCredentialConfig_APNS(resourceName, name, key string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "My test OIDC app for MFA Policy"

  login_page_url = "https://www.pingidentity.com"

  enabled = true

  oidc_options = {
    type                       = "NATIVE_APP"
    grant_types                = ["CLIENT_CREDENTIALS"]
    token_endpoint_auth_method = "CLIENT_SECRET_BASIC"

    mobile_app = {
      bundle_id = "com.%[2]s.bundle"

      passcode_refresh_seconds = 45

      integrity_detection = {
        enabled = true
        cache_duration = {
          amount = 30
          units  = "HOURS"
        }
        google_play = {
          verification_type = "INTERNAL"
          decryption_key    = "dummykeydoesnotexist"
          verification_key  = "dummykeydoesnotexist"
        }
      }
    }
  }
}

resource "pingone_mfa_application_push_credential" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  application_id = pingone_application.%[2]s.id

  apns = {
    key               = "%[3]s"
    team_id           = "team.id.updated"
    token_signing_key = <<EOT
%[4]s
EOT

  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name, key)
}

func testAccApplicationPushCredentialConfig_HMS(resourceName, name, key string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "My test OIDC app for MFA Policy"

  login_page_url = "https://www.pingidentity.com"

  enabled = true

  oidc_options = {
    type                       = "NATIVE_APP"
    grant_types                = ["CLIENT_CREDENTIALS"]
    token_endpoint_auth_method = "CLIENT_SECRET_BASIC"

    mobile_app = {
      huawei_app_id       = "%[2]s"
      huawei_package_name = "com.%[2]s.huaweipackage"

      passcode_refresh_seconds = 45

      integrity_detection = {
        enabled = true
        cache_duration = {
          amount = 30
          units  = "HOURS"
        }
        google_play = {
          verification_type = "INTERNAL"
          decryption_key    = "dummykeydoesnotexist"
          verification_key  = "dummykeydoesnotexist"
        }
      }
    }
  }
}

resource "pingone_mfa_application_push_credential" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  application_id = pingone_application.%[2]s.id

  hms = {
    client_id     = "%[3]s"
    client_secret = "%[4]s"
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name, key)
}
