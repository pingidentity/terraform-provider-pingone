package base_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func testAccCheckWebhookDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_webhook" {
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

		body, r, err := apiClient.SubscriptionsWebhooksApi.ReadOneSubscription(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.ID).Execute()

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

		return fmt.Errorf("PingOne Webhook Instance %s still exists", rs.Primary.ID)
	}

	return nil
}

func testAccGetWebhookIDs(resourceName string, environmentID, resourceID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource Not found: %s", resourceName)
		}

		*resourceID = rs.Primary.ID
		*environmentID = rs.Primary.Attributes["environment_id"]

		return nil
	}
}

func TestAccWebhook_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_webhook.%s", resourceName)

	name := resourceName

	var resourceID, environmentID string

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckWebhookDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccWebhookConfig_Minimal(resourceName, name),
				Check:  testAccGetWebhookIDs(resourceFullName, &environmentID, &resourceID),
			},
			// Replan after removal preconfig
			{
				PreConfig: func() {
					var ctx = context.Background()
					p1Client, err := acctest.TestClient(ctx)

					if err != nil {
						t.Fatalf("Failed to get API client: %v", err)
					}

					apiClient := p1Client.API.ManagementAPIClient

					if environmentID == "" || resourceID == "" {
						t.Fatalf("One of environment ID or resource ID cannot be determined. Environment ID: %s, Resource ID: %s", environmentID, resourceID)
					}

					_, err = apiClient.SubscriptionsWebhooksApi.DeleteSubscription(ctx, environmentID, resourceID).Execute()
					if err != nil {
						t.Fatalf("Failed to delete webhook subsription: %v", err)
					}
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
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckWebhookDestroy,
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
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckWebhookDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccWebhookConfig_Full(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "http_endpoint_url", "https://localhost/"),
					resource.TestCheckResourceAttr(resourceFullName, "http_endpoint_headers.%", "2"),
					resource.TestCheckResourceAttr(resourceFullName, "http_endpoint_headers.Authorization", "Basic usernamepassword"),
					resource.TestCheckResourceAttr(resourceFullName, "http_endpoint_headers.Content-Type", "application/json"),
					resource.TestCheckResourceAttr(resourceFullName, "verify_tls_certificates", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "format", "ACTIVITY"),
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.0.included_action_types.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "filter_options.0.included_action_types.*", "ACCOUNT.LINKED"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "filter_options.0.included_action_types.*", "ACCOUNT.UNLINKED"),
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.0.included_application_ids.#", "3"),
					resource.TestMatchResourceAttr(resourceFullName, "filter_options.0.included_application_ids.0", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "filter_options.0.included_application_ids.1", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "filter_options.0.included_application_ids.2", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.0.included_population_ids.#", "3"),
					resource.TestMatchResourceAttr(resourceFullName, "filter_options.0.included_population_ids.0", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "filter_options.0.included_population_ids.1", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "filter_options.0.included_population_ids.2", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.0.included_tags.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "filter_options.0.included_tags.*", "adminIdentityEvent"),
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.0.ip_address_exposed", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.0.useragent_exposed", "true"),
				),
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
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckWebhookDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccWebhookConfig_Minimal(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "http_endpoint_url", "https://localhost/"),
					resource.TestCheckResourceAttr(resourceFullName, "http_endpoint_headers.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "verify_tls_certificates", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "format", "SPLUNK"),
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.0.included_action_types.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "filter_options.0.included_action_types.*", "ACCOUNT.LINKED"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "filter_options.0.included_action_types.*", "ACCOUNT.UNLINKED"),
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.0.included_application_ids.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.0.included_population_ids.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.0.included_tags.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.0.ip_address_exposed", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.0.useragent_exposed", "false"),
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
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckWebhookDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccWebhookConfig_Full(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "http_endpoint_url", "https://localhost/"),
					resource.TestCheckResourceAttr(resourceFullName, "http_endpoint_headers.%", "2"),
					resource.TestCheckResourceAttr(resourceFullName, "http_endpoint_headers.Authorization", "Basic usernamepassword"),
					resource.TestCheckResourceAttr(resourceFullName, "http_endpoint_headers.Content-Type", "application/json"),
					resource.TestCheckResourceAttr(resourceFullName, "verify_tls_certificates", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "format", "ACTIVITY"),
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.0.included_action_types.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "filter_options.0.included_action_types.*", "ACCOUNT.LINKED"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "filter_options.0.included_action_types.*", "ACCOUNT.UNLINKED"),
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.0.included_application_ids.#", "3"),
					resource.TestMatchResourceAttr(resourceFullName, "filter_options.0.included_application_ids.0", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "filter_options.0.included_application_ids.1", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "filter_options.0.included_application_ids.2", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.0.included_population_ids.#", "3"),
					resource.TestMatchResourceAttr(resourceFullName, "filter_options.0.included_population_ids.0", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "filter_options.0.included_population_ids.1", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "filter_options.0.included_population_ids.2", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.0.included_tags.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "filter_options.0.included_tags.*", "adminIdentityEvent"),
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.0.ip_address_exposed", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.0.useragent_exposed", "true"),
				),
			},
			{
				Config: testAccWebhookConfig_Minimal(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "http_endpoint_url", "https://localhost/"),
					resource.TestCheckResourceAttr(resourceFullName, "http_endpoint_headers.%", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "verify_tls_certificates", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "format", "SPLUNK"),
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.0.included_action_types.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "filter_options.0.included_action_types.*", "ACCOUNT.LINKED"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "filter_options.0.included_action_types.*", "ACCOUNT.UNLINKED"),
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.0.included_application_ids.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.0.included_population_ids.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.0.included_tags.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.0.ip_address_exposed", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.0.useragent_exposed", "false"),
				),
			},
			{
				Config: testAccWebhookConfig_Full(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "http_endpoint_url", "https://localhost/"),
					resource.TestCheckResourceAttr(resourceFullName, "http_endpoint_headers.%", "2"),
					resource.TestCheckResourceAttr(resourceFullName, "http_endpoint_headers.Authorization", "Basic usernamepassword"),
					resource.TestCheckResourceAttr(resourceFullName, "http_endpoint_headers.Content-Type", "application/json"),
					resource.TestCheckResourceAttr(resourceFullName, "verify_tls_certificates", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "format", "ACTIVITY"),
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.0.included_action_types.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "filter_options.0.included_action_types.*", "ACCOUNT.LINKED"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "filter_options.0.included_action_types.*", "ACCOUNT.UNLINKED"),
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.0.included_application_ids.#", "3"),
					resource.TestMatchResourceAttr(resourceFullName, "filter_options.0.included_application_ids.0", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "filter_options.0.included_application_ids.1", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "filter_options.0.included_application_ids.2", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.0.included_population_ids.#", "3"),
					resource.TestMatchResourceAttr(resourceFullName, "filter_options.0.included_population_ids.0", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "filter_options.0.included_population_ids.1", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "filter_options.0.included_population_ids.2", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.0.included_tags.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "filter_options.0.included_tags.*", "adminIdentityEvent"),
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.0.ip_address_exposed", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.0.useragent_exposed", "true"),
				),
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
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckWebhookDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccWebhookConfig_Profile1(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.0.included_application_ids.#", "3"),
					resource.TestMatchResourceAttr(resourceFullName, "filter_options.0.included_application_ids.0", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "filter_options.0.included_application_ids.1", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "filter_options.0.included_application_ids.2", verify.P1ResourceIDRegexp),
				),
			},
			{
				Config: testAccWebhookConfig_Profile2(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.0.included_application_ids.#", "2"),
					resource.TestMatchResourceAttr(resourceFullName, "filter_options.0.included_application_ids.0", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "filter_options.0.included_application_ids.1", verify.P1ResourceIDRegexp),
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
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckWebhookDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccWebhookConfig_Profile1(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.0.included_population_ids.#", "3"),
					resource.TestMatchResourceAttr(resourceFullName, "filter_options.0.included_population_ids.0", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "filter_options.0.included_population_ids.1", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "filter_options.0.included_population_ids.2", verify.P1ResourceIDRegexp),
				),
			},
			{
				Config: testAccWebhookConfig_Profile2(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "filter_options.0.included_population_ids.#", "2"),
					resource.TestMatchResourceAttr(resourceFullName, "filter_options.0.included_population_ids.0", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "filter_options.0.included_population_ids.1", verify.P1ResourceIDRegexp),
				),
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

  filter_options {
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

  oidc_options {
    type                        = "SINGLE_PAGE_APP"
    grant_types                 = ["AUTHORIZATION_CODE"]
    response_types              = ["CODE"]
    pkce_enforcement            = "S256_REQUIRED"
    token_endpoint_authn_method = "NONE"
    redirect_uris               = ["https://www.pingidentity.com"]
  }
}

resource "pingone_application" "%[2]s-2" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-2"
  enabled        = true

  oidc_options {
    type                        = "SINGLE_PAGE_APP"
    grant_types                 = ["AUTHORIZATION_CODE"]
    response_types              = ["CODE"]
    pkce_enforcement            = "S256_REQUIRED"
    token_endpoint_authn_method = "NONE"
    redirect_uris               = ["https://www.pingidentity.com"]
  }
}

resource "pingone_application" "%[2]s-3" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-3"
  enabled        = true

  oidc_options {
    type                        = "SINGLE_PAGE_APP"
    grant_types                 = ["AUTHORIZATION_CODE"]
    response_types              = ["CODE"]
    pkce_enforcement            = "S256_REQUIRED"
    token_endpoint_authn_method = "NONE"
    redirect_uris               = ["https://www.pingidentity.com"]
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

  filter_options {
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

  filter_options {
    included_action_types = ["ACCOUNT.LINKED", "ACCOUNT.UNLINKED"]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
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

  oidc_options {
    type                        = "SINGLE_PAGE_APP"
    grant_types                 = ["AUTHORIZATION_CODE"]
    response_types              = ["CODE"]
    pkce_enforcement            = "S256_REQUIRED"
    token_endpoint_authn_method = "NONE"
    redirect_uris               = ["https://www.pingidentity.com"]
  }
}

resource "pingone_application" "%[2]s-2" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-2"
  enabled        = true

  oidc_options {
    type                        = "SINGLE_PAGE_APP"
    grant_types                 = ["AUTHORIZATION_CODE"]
    response_types              = ["CODE"]
    pkce_enforcement            = "S256_REQUIRED"
    token_endpoint_authn_method = "NONE"
    redirect_uris               = ["https://www.pingidentity.com"]
  }
}

resource "pingone_application" "%[2]s-3" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-3"
  enabled        = true

  oidc_options {
    type                        = "SINGLE_PAGE_APP"
    grant_types                 = ["AUTHORIZATION_CODE"]
    response_types              = ["CODE"]
    pkce_enforcement            = "S256_REQUIRED"
    token_endpoint_authn_method = "NONE"
    redirect_uris               = ["https://www.pingidentity.com"]
  }
}

resource "pingone_webhook" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name              = "%[3]s"
  enabled           = false
  http_endpoint_url = "https://localhost/"

  format = "ACTIVITY"

  filter_options {
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

  oidc_options {
    type                        = "SINGLE_PAGE_APP"
    grant_types                 = ["AUTHORIZATION_CODE"]
    response_types              = ["CODE"]
    pkce_enforcement            = "S256_REQUIRED"
    token_endpoint_authn_method = "NONE"
    redirect_uris               = ["https://www.pingidentity.com"]
  }
}

resource "pingone_application" "%[2]s-new" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-new"
  enabled        = true

  oidc_options {
    type                        = "SINGLE_PAGE_APP"
    grant_types                 = ["AUTHORIZATION_CODE"]
    response_types              = ["CODE"]
    pkce_enforcement            = "S256_REQUIRED"
    token_endpoint_authn_method = "NONE"
    redirect_uris               = ["https://www.pingidentity.com"]
  }
}

resource "pingone_webhook" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name              = "%[3]s"
  enabled           = false
  http_endpoint_url = "https://localhost/"

  format = "ACTIVITY"

  filter_options {
    included_action_types    = ["ACCOUNT.LINKED"]
    included_application_ids = [pingone_application.%[3]s-new.id, pingone_application.%[3]s-1.id]
    included_population_ids  = [pingone_population.%[3]s-new.id, pingone_population.%[3]s-1.id]
    ip_address_exposed       = false
    useragent_exposed        = true
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
