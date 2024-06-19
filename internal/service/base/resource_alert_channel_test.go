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

func TestAccAlertChannel_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_alert_channel.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	var AlertChannelID, environmentID string

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
		CheckDestroy:             base.AlertChannel_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccAlertChannelConfig_Minimal(resourceName, name),
				Check:  base.AlertChannel_GetIDs(resourceFullName, &environmentID, &AlertChannelID),
			},
			{
				PreConfig: func() {
					base.AlertChannel_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID, AlertChannelID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the environment
			{
				Config: testAccAlertChannel_NewEnv(environmentName, licenseID, resourceName, name),
				Check:  base.AlertChannel_GetIDs(resourceFullName, &environmentID, &AlertChannelID),
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

func TestAccAlertChannel_NewEnv(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_alert_channel.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.AlertChannel_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlertChannel_NewEnv(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
				),
			},
		},
	})
}

func TestAccAlertChannel_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_alert_channel.%s", resourceName)

	name := resourceName

	fullStep1 := resource.TestStep{
		Config: testAccAlertChannelConfig_Full(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestCheckResourceAttr(resourceFullName, "default", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "quota.#", "2"),
			resource.TestCheckResourceAttr(resourceFullName, "country_limit.type", "DENIED"),
		),
	}

	fullStep2 := resource.TestStep{
		Config: testAccAlertChannelConfig_PartialFull(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestCheckResourceAttr(resourceFullName, "default", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "quota.#", "2"),
			resource.TestCheckResourceAttr(resourceFullName, "country_limit.type", "DENIED"),
		),
	}

	minimalStep := resource.TestStep{
		Config: testAccAlertChannelConfig_Minimal(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestCheckResourceAttr(resourceFullName, "default", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "quota.#", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "country_limit.type", "NONE"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.AlertChannel_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			fullStep1,
			{
				Config:  testAccAlertChannelConfig_Full(resourceName, name),
				Destroy: true,
			},
			// PartialFull
			fullStep2,
			{
				Config:  testAccAlertChannelConfig_PartialFull(resourceName, name),
				Destroy: true,
			},
			// Minimal
			minimalStep,
			{
				Config:  testAccAlertChannelConfig_Minimal(resourceName, name),
				Destroy: true,
			},
			// Change
			fullStep1,
			fullStep2,
			minimalStep,
			fullStep1,
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

func TestAccAlertChannel_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_alert_channel.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.AlertChannel_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccAlertChannelConfig_Minimal(resourceName, name),
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

func testAccAlertChannel_NewEnv(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_alert_channel" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  alert_name = "%[4]s"

  addresses = [
    "noreply@pingidentity.com",
  ]

  channel_type = "EMAIL"

  include_alert_types = [
    "LICENSE_EXPIRED",
    "LICENSE_EXPIRING",
    "LICENSE_ROTATED",
  ]

  include_severities = [
    "INFO",
    "WARNING",
    "ERROR",
  ]
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccAlertChannelConfig_Full(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_alert_channel" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  alert_name = "%[3]s"

  addresses = [
    "noreply3@pingidentity.com",
    "noreply@pingidentity.com",
    "noreply2@pingidentity.com",
  ]

  channel_type = "EMAIL"

  exclude_alert_types = [
    "CERTIFICATE_EXPIRED",
    "CERTIFICATE_EXPIRING",
    "KEY_PAIR_EXPIRED",
    "KEY_PAIR_EXPIRING",
  ]

  include_alert_types = [
    "GATEWAY_VERSION_DEPRECATED",
    "GATEWAY_VERSION_DEPRECATING",
    "LICENSE_EXPIRED",
    "LICENSE_EXPIRING",
    "LICENSE_ROTATED",
    "APPROACHING_USER_LICENSE_LIMIT",
    "USER_LICENSE_LIMIT_REACHED",
    "USER_LICENSE_LIMIT_EXCEEDED",
    "DATA_QUALITY_ISSUE",
  ]

  include_severities = [
    "INFO",
    "WARNING",
    "ERROR",
  ]
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccAlertChannelConfig_PartialFull(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_alert_channel" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  alert_name = "%[3]s"

  addresses = [
    "noreply3@pingidentity.com",
    "noreply2@pingidentity.com",
  ]

  channel_type = "EMAIL"

  include_alert_types = [
    "LICENSE_EXPIRED",
    "LICENSE_EXPIRING",
    "LICENSE_ROTATED",
  ]

  include_severities = [
    "ERROR",
  ]
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccAlertChannelConfig_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_alert_channel" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  alert_name = "%[3]s"

  addresses = [
    "noreply@pingidentity.com",
  ]

  channel_type = "EMAIL"

  include_alert_types = [
    "LICENSE_EXPIRED",
    "LICENSE_EXPIRING",
    "LICENSE_ROTATED",
  ]
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
