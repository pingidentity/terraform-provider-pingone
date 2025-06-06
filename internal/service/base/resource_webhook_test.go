// Copyright © 2025 Ping Identity Corporation

package base_test

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
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccWebhook_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_webhook.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	var subscriptionID, environmentID string

	var p1Client *client.Client
	var ctx = context.Background()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)

			p1Client = acctest.PreCheckTestClient(ctx, t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Webhook_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccWebhookConfig_Minimal(resourceName, name),
				Check:  base.Webhook_GetIDs(resourceFullName, &environmentID, &subscriptionID),
			},
			{
				PreConfig: func() {
					base.Webhook_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID, subscriptionID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the environment
			{
				Config: testAccWebhookConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check:  base.Webhook_GetIDs(resourceFullName, &environmentID, &subscriptionID),
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

func TestAccWebhook_NewEnv(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_webhook.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Webhook_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccWebhookConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
				),
			},
		},
	})
}

func TestAccWebhook_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_webhook.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Webhook_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccWebhookConfig_Full(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "http_endpoint_url", "https://localhost/"),
					resource.TestCheckResourceAttr(resourceFullName, "http_endpoint_headers.%", "2"),
					resource.TestCheckResourceAttr(resourceFullName, "http_endpoint_headers.Authorization", "Basic usernamepassword"),
					resource.TestCheckResourceAttr(resourceFullName, "http_endpoint_headers.Content-Type", "application/json"),
					resource.TestCheckResourceAttr(resourceFullName, "verify_tls_certificates", "false"),
					resource.TestCheckNoResourceAttr(resourceFullName, "tls_client_auth_key_pair_id"),
					resource.TestCheckResourceAttr(resourceFullName, "format", "ACTIVITY"),
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.included_action_types.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "filter_options.included_action_types.*", "ACCOUNT.LINKED"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "filter_options.included_action_types.*", "ACCOUNT.UNLINKED"),
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.included_application_ids.#", "3"),
					resource.TestMatchResourceAttr(resourceFullName, "filter_options.included_application_ids.0", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "filter_options.included_application_ids.1", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "filter_options.included_application_ids.2", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.included_population_ids.#", "3"),
					resource.TestMatchResourceAttr(resourceFullName, "filter_options.included_population_ids.0", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "filter_options.included_population_ids.1", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "filter_options.included_population_ids.2", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.included_tags.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "filter_options.included_tags.*", "adminIdentityEvent"),
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.ip_address_exposed", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.useragent_exposed", "true"),
				),
			},
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

func TestAccWebhook_Minimal(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_webhook.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Webhook_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccWebhookConfig_Minimal(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "http_endpoint_url", "https://localhost/"),
					resource.TestCheckResourceAttr(resourceFullName, "http_endpoint_headers.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "verify_tls_certificates", "true"),
					resource.TestCheckNoResourceAttr(resourceFullName, "tls_client_auth_key_pair_id"),
					resource.TestCheckResourceAttr(resourceFullName, "format", "SPLUNK"),
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.included_action_types.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "filter_options.included_action_types.*", "ACCOUNT.LINKED"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "filter_options.included_action_types.*", "ACCOUNT.UNLINKED"),
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.included_application_ids.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.included_population_ids.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.included_tags.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.ip_address_exposed", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.useragent_exposed", "false"),
				),
			},
		},
	})
}

func TestAccWebhook_Change(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_webhook.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Webhook_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccWebhookConfig_Full(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "http_endpoint_url", "https://localhost/"),
					resource.TestCheckResourceAttr(resourceFullName, "http_endpoint_headers.%", "2"),
					resource.TestCheckResourceAttr(resourceFullName, "http_endpoint_headers.Authorization", "Basic usernamepassword"),
					resource.TestCheckResourceAttr(resourceFullName, "http_endpoint_headers.Content-Type", "application/json"),
					resource.TestCheckResourceAttr(resourceFullName, "verify_tls_certificates", "false"),
					resource.TestCheckNoResourceAttr(resourceFullName, "tls_client_auth_key_pair_id"),
					resource.TestCheckResourceAttr(resourceFullName, "format", "ACTIVITY"),
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.included_action_types.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "filter_options.included_action_types.*", "ACCOUNT.LINKED"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "filter_options.included_action_types.*", "ACCOUNT.UNLINKED"),
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.included_application_ids.#", "3"),
					resource.TestMatchResourceAttr(resourceFullName, "filter_options.included_application_ids.0", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "filter_options.included_application_ids.1", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "filter_options.included_application_ids.2", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.included_population_ids.#", "3"),
					resource.TestMatchResourceAttr(resourceFullName, "filter_options.included_population_ids.0", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "filter_options.included_population_ids.1", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "filter_options.included_population_ids.2", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.included_tags.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "filter_options.included_tags.*", "adminIdentityEvent"),
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.ip_address_exposed", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.useragent_exposed", "true"),
				),
			},
			{
				Config: testAccWebhookConfig_Minimal(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "http_endpoint_url", "https://localhost/"),
					resource.TestCheckResourceAttr(resourceFullName, "http_endpoint_headers.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "verify_tls_certificates", "true"),
					resource.TestCheckNoResourceAttr(resourceFullName, "tls_client_auth_key_pair_id"),
					resource.TestCheckResourceAttr(resourceFullName, "format", "SPLUNK"),
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.included_action_types.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "filter_options.included_action_types.*", "ACCOUNT.LINKED"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "filter_options.included_action_types.*", "ACCOUNT.UNLINKED"),
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.included_application_ids.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.included_population_ids.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.included_tags.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.ip_address_exposed", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.useragent_exposed", "false"),
				),
			},
			{
				Config: testAccWebhookConfig_Full(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "http_endpoint_url", "https://localhost/"),
					resource.TestCheckResourceAttr(resourceFullName, "http_endpoint_headers.%", "2"),
					resource.TestCheckResourceAttr(resourceFullName, "http_endpoint_headers.Authorization", "Basic usernamepassword"),
					resource.TestCheckResourceAttr(resourceFullName, "http_endpoint_headers.Content-Type", "application/json"),
					resource.TestCheckResourceAttr(resourceFullName, "verify_tls_certificates", "false"),
					resource.TestCheckNoResourceAttr(resourceFullName, "tls_client_auth_key_pair_id"),
					resource.TestCheckResourceAttr(resourceFullName, "format", "ACTIVITY"),
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.included_action_types.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "filter_options.included_action_types.*", "ACCOUNT.LINKED"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "filter_options.included_action_types.*", "ACCOUNT.UNLINKED"),
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.included_application_ids.#", "3"),
					resource.TestMatchResourceAttr(resourceFullName, "filter_options.included_application_ids.0", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "filter_options.included_application_ids.1", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "filter_options.included_application_ids.2", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.included_population_ids.#", "3"),
					resource.TestMatchResourceAttr(resourceFullName, "filter_options.included_population_ids.0", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "filter_options.included_population_ids.1", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "filter_options.included_population_ids.2", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.included_tags.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "filter_options.included_tags.*", "adminIdentityEvent"),
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.ip_address_exposed", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.useragent_exposed", "true"),
				),
			},
		},
	})
}

func TestAccWebhook_MTLS(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_webhook.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	pkcs12 := os.Getenv("PINGONE_KEY_PKCS12")
	keystorePassword := os.Getenv("PINGONE_KEY_PKCS12_PASSWORD")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
			acctest.PreCheckPKCS12Key(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Webhook_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccWebhookConfig_MTLS(environmentName, licenseID, resourceName, name, pkcs12, keystorePassword),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "tls_client_auth_key_pair_id", verify.P1ResourceIDRegexpFullString),
				),
			},
			{
				Config: testAccWebhookConfig_NoMTLS(environmentName, licenseID, resourceName, name, pkcs12, keystorePassword),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr(resourceFullName, "tls_client_auth_key_pair_id"),
				),
			},
			{
				Config: testAccWebhookConfig_MTLS(environmentName, licenseID, resourceName, name, pkcs12, keystorePassword),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "tls_client_auth_key_pair_id", verify.P1ResourceIDRegexpFullString),
				),
			},
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

func TestAccWebhook_Applications(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_webhook.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Webhook_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccWebhookConfig_Profile1(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.included_application_ids.#", "3"),
					resource.TestMatchResourceAttr(resourceFullName, "filter_options.included_application_ids.0", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "filter_options.included_application_ids.1", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "filter_options.included_application_ids.2", verify.P1ResourceIDRegexpFullString),
				),
			},
			{
				Config: testAccWebhookConfig_Profile2(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.included_application_ids.#", "2"),
					resource.TestMatchResourceAttr(resourceFullName, "filter_options.included_application_ids.0", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "filter_options.included_application_ids.1", verify.P1ResourceIDRegexpFullString),
				),
			},
		},
	})
}

func TestAccWebhook_Populations(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_webhook.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Webhook_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccWebhookConfig_Profile1(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.included_population_ids.#", "3"),
					resource.TestMatchResourceAttr(resourceFullName, "filter_options.included_population_ids.0", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "filter_options.included_population_ids.1", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "filter_options.included_population_ids.2", verify.P1ResourceIDRegexpFullString),
				),
			},
			{
				Config: testAccWebhookConfig_Profile2(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.included_population_ids.#", "2"),
					resource.TestMatchResourceAttr(resourceFullName, "filter_options.included_population_ids.0", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "filter_options.included_population_ids.1", verify.P1ResourceIDRegexpFullString),
				),
			},
		},
	})
}

func TestAccWebhook_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_webhook.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Webhook_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccWebhookConfig_Minimal(resourceName, name),
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

func testAccWebhookConfig_NewEnv(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_webhook" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  name              = "%[4]s"
  enabled           = "true"
  http_endpoint_url = "https://localhost/"

  format = "ACTIVITY"

  filter_options = {
    included_action_types = ["ACCOUNT.LINKED", "ACCOUNT.UNLINKED"]
  }
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccWebhookConfig_Full(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_population" "%[2]s-1" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s-1"
}

resource "pingone_population" "%[2]s-2" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s-2"
}

resource "pingone_population" "%[2]s-3" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s-3"
}

resource "pingone_application" "%[2]s-1" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-1"
  enabled        = true

  oidc_options = {
    type                       = "SINGLE_PAGE_APP"
    grant_types                = ["AUTHORIZATION_CODE"]
    response_types             = ["CODE"]
    pkce_enforcement           = "S256_REQUIRED"
    token_endpoint_auth_method = "NONE"
    redirect_uris              = ["https://www.pingidentity.com"]
  }
}

resource "pingone_application" "%[2]s-2" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-2"
  enabled        = true

  oidc_options = {
    type                       = "SINGLE_PAGE_APP"
    grant_types                = ["AUTHORIZATION_CODE"]
    response_types             = ["CODE"]
    pkce_enforcement           = "S256_REQUIRED"
    token_endpoint_auth_method = "NONE"
    redirect_uris              = ["https://www.pingidentity.com"]
  }
}

resource "pingone_application" "%[2]s-3" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-3"
  enabled        = true

  oidc_options = {
    type                       = "SINGLE_PAGE_APP"
    grant_types                = ["AUTHORIZATION_CODE"]
    response_types             = ["CODE"]
    pkce_enforcement           = "S256_REQUIRED"
    token_endpoint_auth_method = "NONE"
    redirect_uris              = ["https://www.pingidentity.com"]
  }
}

resource "pingone_webhook" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name              = "%[3]s"
  enabled           = false
  http_endpoint_url = "https://localhost/"

  http_endpoint_headers = {
    Authorization = "Basic usernamepassword"
    Content-Type  = "application/json"
  }

  verify_tls_certificates = false

  format = "ACTIVITY"

  filter_options = {
    included_action_types    = ["ACCOUNT.LINKED", "ACCOUNT.UNLINKED"]
    included_application_ids = [pingone_application.%[3]s-2.id, pingone_application.%[3]s-3.id, pingone_application.%[3]s-1.id]
    included_population_ids  = [pingone_population.%[3]s-2.id, pingone_population.%[3]s-3.id, pingone_population.%[3]s-1.id]
    included_tags            = ["adminIdentityEvent"]
    ip_address_exposed       = true
    useragent_exposed        = true
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccWebhookConfig_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_webhook" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name              = "%[3]s"
  http_endpoint_url = "https://localhost/"

  format = "SPLUNK"

  filter_options = {
    included_action_types = ["ACCOUNT.LINKED", "ACCOUNT.UNLINKED"]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccWebhookConfig_MTLS(environmentName, licenseID, resourceName, name, pkcs12, keystorePassword string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_key" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  pkcs12_file_base64 = <<EOT
%[5]s
EOT

  pkcs12_file_password = "%[6]s"

  usage_type = "OUTBOUND_MTLS"
}

resource "pingone_webhook" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  name              = "%[4]s"
  enabled           = "true"
  http_endpoint_url = "https://localhost/"

  tls_client_auth_key_pair_id = pingone_key.%[3]s.id

  format = "ACTIVITY"

  filter_options = {
    included_action_types = ["ACCOUNT.LINKED", "ACCOUNT.UNLINKED"]
  }
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name, pkcs12, keystorePassword)
}

func testAccWebhookConfig_NoMTLS(environmentName, licenseID, resourceName, name, pkcs12, keystorePassword string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_key" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  pkcs12_file_base64 = <<EOT
%[5]s
EOT

  pkcs12_file_password = "%[6]s"

  usage_type = "OUTBOUND_MTLS"
}

resource "pingone_webhook" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  name              = "%[4]s"
  enabled           = "true"
  http_endpoint_url = "https://localhost/"

  format = "ACTIVITY"

  filter_options = {
    included_action_types = ["ACCOUNT.LINKED", "ACCOUNT.UNLINKED"]
  }
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name, pkcs12, keystorePassword)
}

func testAccWebhookConfig_Profile1(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_population" "%[2]s-1" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s-1"
}

resource "pingone_population" "%[2]s-2" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s-2"
}

resource "pingone_population" "%[2]s-3" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s-3"
}

resource "pingone_application" "%[2]s-1" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-1"
  enabled        = true

  oidc_options = {
    type                       = "SINGLE_PAGE_APP"
    grant_types                = ["AUTHORIZATION_CODE"]
    response_types             = ["CODE"]
    pkce_enforcement           = "S256_REQUIRED"
    token_endpoint_auth_method = "NONE"
    redirect_uris              = ["https://www.pingidentity.com"]
  }
}

resource "pingone_application" "%[2]s-2" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-2"
  enabled        = true

  oidc_options = {
    type                       = "SINGLE_PAGE_APP"
    grant_types                = ["AUTHORIZATION_CODE"]
    response_types             = ["CODE"]
    pkce_enforcement           = "S256_REQUIRED"
    token_endpoint_auth_method = "NONE"
    redirect_uris              = ["https://www.pingidentity.com"]
  }
}

resource "pingone_application" "%[2]s-3" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-3"
  enabled        = true

  oidc_options = {
    type                       = "SINGLE_PAGE_APP"
    grant_types                = ["AUTHORIZATION_CODE"]
    response_types             = ["CODE"]
    pkce_enforcement           = "S256_REQUIRED"
    token_endpoint_auth_method = "NONE"
    redirect_uris              = ["https://www.pingidentity.com"]
  }
}

resource "pingone_webhook" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name              = "%[3]s"
  enabled           = false
  http_endpoint_url = "https://localhost/"

  format = "ACTIVITY"

  filter_options = {
    included_action_types    = ["ACCOUNT.LINKED", "ACCOUNT.UNLINKED"]
    included_application_ids = [pingone_application.%[3]s-2.id, pingone_application.%[3]s-3.id, pingone_application.%[3]s-1.id]
    included_population_ids  = [pingone_population.%[3]s-2.id, pingone_population.%[3]s-3.id, pingone_population.%[3]s-1.id]
    ip_address_exposed       = true
    useragent_exposed        = false
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccWebhookConfig_Profile2(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_population" "%[2]s-1" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s-1"
}

resource "pingone_population" "%[2]s-new" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s-new"
}

resource "pingone_application" "%[2]s-1" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-1"
  enabled        = true

  oidc_options = {
    type                       = "SINGLE_PAGE_APP"
    grant_types                = ["AUTHORIZATION_CODE"]
    response_types             = ["CODE"]
    pkce_enforcement           = "S256_REQUIRED"
    token_endpoint_auth_method = "NONE"
    redirect_uris              = ["https://www.pingidentity.com"]
  }
}

resource "pingone_application" "%[2]s-new" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-new"
  enabled        = true

  oidc_options = {
    type                       = "SINGLE_PAGE_APP"
    grant_types                = ["AUTHORIZATION_CODE"]
    response_types             = ["CODE"]
    pkce_enforcement           = "S256_REQUIRED"
    token_endpoint_auth_method = "NONE"
    redirect_uris              = ["https://www.pingidentity.com"]
  }
}

resource "pingone_webhook" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name              = "%[3]s"
  enabled           = false
  http_endpoint_url = "https://localhost/"

  format = "ACTIVITY"

  filter_options = {
    included_action_types    = ["ACCOUNT.LINKED"]
    included_application_ids = [pingone_application.%[3]s-new.id, pingone_application.%[3]s-1.id]
    included_population_ids  = [pingone_population.%[3]s-new.id, pingone_population.%[3]s-1.id]
    ip_address_exposed       = false
    useragent_exposed        = true
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
