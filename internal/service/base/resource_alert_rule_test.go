package base_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	pingone "github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
)

func testAccCheckAlertRuleDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, pingone.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_alert_rule" {
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

		body, r, err := apiClient.AlertingApi.ReadOneAlertChannel(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.ID).Execute()

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

		return fmt.Errorf("PingOne Alert Rule Instance %s still exists", rs.Primary.ID)
	}

	return nil
}

func TestAccAlertRule_NewEnv(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_alert_rule.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckAlertRuleDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlertRuleConfig_NewEnv(environmentName, licenseID, resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "type", "EMAIL"),
					resource.TestCheckResourceAttr(resourceFullName, "addresses.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "addresses.*", "noreply@pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "addresses.*", "noreply1@pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "include_severities.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "include_alert_types.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "exclude_alert_types.#", "0"),
				),
			},
		},
	})
}

func TestAccAlertRule_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_alert_rule.%s", resourceName)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckAlertRuleDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlertRuleConfig_Full(resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "type", "EMAIL"),
					resource.TestCheckResourceAttr(resourceFullName, "addresses.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "addresses.*", "noreply@pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "addresses.*", "noreply1@pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "include_severities.#", "3"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "include_severities.*", "INFO"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "include_severities.*", "WARNING"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "include_severities.*", "ERROR"),
					resource.TestCheckResourceAttr(resourceFullName, "include_alert_types.#", "3"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "include_alert_types.*", "KEY_PAIR_EXPIRED"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "include_alert_types.*", "CERTIFICATE_EXPIRED"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "include_alert_types.*", "CERTIFICATE_EXPIRING"),
					resource.TestCheckResourceAttr(resourceFullName, "exclude_alert_types.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "exclude_alert_types.*", "GATEWAY_VERSION_DEPRECATED"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "exclude_alert_types.*", "GATEWAY_VERSION_DEPRECATING"),
				),
			},
		},
	})
}

func TestAccAlertRule_Minimal(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_alert_rule.%s", resourceName)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckAlertRuleDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlertRuleConfig_Minimal(resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "type", "EMAIL"),
					resource.TestCheckResourceAttr(resourceFullName, "addresses.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "addresses.*", "noreply@pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "addresses.*", "noreply1@pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "include_severities.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "include_alert_types.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "exclude_alert_types.#", "0"),
				),
			},
		},
	})
}

func TestAccAlertRule_Change(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_alert_rule.%s", resourceName)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckAlertRuleDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlertRuleConfig_Full(resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "type", "EMAIL"),
					resource.TestCheckResourceAttr(resourceFullName, "addresses.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "addresses.*", "noreply@pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "addresses.*", "noreply1@pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "include_severities.#", "3"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "include_severities.*", "INFO"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "include_severities.*", "WARNING"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "include_severities.*", "ERROR"),
					resource.TestCheckResourceAttr(resourceFullName, "include_alert_types.#", "3"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "include_alert_types.*", "KEY_PAIR_EXPIRED"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "include_alert_types.*", "CERTIFICATE_EXPIRED"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "include_alert_types.*", "CERTIFICATE_EXPIRING"),
					resource.TestCheckResourceAttr(resourceFullName, "exclude_alert_types.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "exclude_alert_types.*", "GATEWAY_VERSION_DEPRECATED"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "exclude_alert_types.*", "GATEWAY_VERSION_DEPRECATING"),
				),
			},
			{
				Config: testAccAlertRuleConfig_Minimal(resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "type", "EMAIL"),
					resource.TestCheckResourceAttr(resourceFullName, "addresses.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "addresses.*", "noreply@pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "addresses.*", "noreply1@pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "include_severities.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "include_alert_types.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "exclude_alert_types.#", "0"),
				),
			},
			{
				Config: testAccAlertRuleConfig_FullChange(resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "type", "EMAIL"),
					resource.TestCheckResourceAttr(resourceFullName, "addresses.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "addresses.*", "noreply@pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "include_severities.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "include_severities.*", "INFO"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "include_severities.*", "WARNING"),
					resource.TestCheckResourceAttr(resourceFullName, "include_alert_types.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "include_alert_types.*", "GATEWAY_VERSION_DEPRECATED"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "include_alert_types.*", "GATEWAY_VERSION_DEPRECATING"),
					resource.TestCheckResourceAttr(resourceFullName, "exclude_alert_types.#", "3"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "exclude_alert_types.*", "KEY_PAIR_EXPIRED"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "exclude_alert_types.*", "CERTIFICATE_EXPIRED"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "exclude_alert_types.*", "CERTIFICATE_EXPIRING"),
				),
			},
			{
				Config: testAccAlertRuleConfig_Full(resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "type", "EMAIL"),
					resource.TestCheckResourceAttr(resourceFullName, "addresses.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "addresses.*", "noreply@pingidentity.com"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "addresses.*", "noreply1@pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "include_severities.#", "3"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "include_severities.*", "INFO"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "include_severities.*", "WARNING"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "include_severities.*", "ERROR"),
					resource.TestCheckResourceAttr(resourceFullName, "include_alert_types.#", "3"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "include_alert_types.*", "KEY_PAIR_EXPIRED"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "include_alert_types.*", "CERTIFICATE_EXPIRED"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "include_alert_types.*", "CERTIFICATE_EXPIRING"),
					resource.TestCheckResourceAttr(resourceFullName, "exclude_alert_types.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "exclude_alert_types.*", "GATEWAY_VERSION_DEPRECATED"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "exclude_alert_types.*", "GATEWAY_VERSION_DEPRECATING"),
				),
			},
		},
	})
}

func testAccAlertRuleConfig_NewEnv(environmentName, licenseID, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_alert_rule" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  
  addresses = ["noreply1@pingidentity.com", "noreply@pingidentity.com"]
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, name)
}

func testAccAlertRuleConfig_Full(resourceName string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_alert_rule" "%[2]s" {
	environment_id = data.pingone_environment.general_test.id

  channel_type                = "EMAIL"
  addresses = ["noreply1@pingidentity.com", "noreply@pingidentity.com"]
  include_severities = ["INFO", "WARNING", "ERROR"]
  include_alert_types = ["KEY_PAIR_EXPIRED", "CERTIFICATE_EXPIRED", "CERTIFICATE_EXPIRING"]
  exclude_alert_types = ["GATEWAY_VERSION_DEPRECATED", "GATEWAY_VERSION_DEPRECATING"]
}`, acctest.GenericSandboxEnvironment(), resourceName)
}

func testAccAlertRuleConfig_Minimal(resourceName string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_alert_rule" "%[2]s" {
	environment_id = data.pingone_environment.general_test.id

  addresses = ["noreply1@pingidentity.com", "noreply@pingidentity.com"]
}`, acctest.GenericSandboxEnvironment(), resourceName)
}

func testAccAlertRuleConfig_FullChange(resourceName string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_alert_rule" "%[2]s" {
	environment_id = data.pingone_environment.general_test.id

  channel_type                = "EMAIL"
  addresses = ["noreply@pingidentity.com"]
  include_severities = ["INFO", "WARNING"]
  include_alert_types = ["GATEWAY_VERSION_DEPRECATED", "GATEWAY_VERSION_DEPRECATING"]
  exclude_alert_types = ["KEY_PAIR_EXPIRED", "CERTIFICATE_EXPIRED", "CERTIFICATE_EXPIRING"]
}`, acctest.GenericSandboxEnvironment(), resourceName)
}
