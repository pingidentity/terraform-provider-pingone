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
	acctestlegacysdk "github.com/pingidentity/terraform-provider-pingone/internal/acctest/legacysdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/base"
	baselegacysdk "github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/base/legacysdk"
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

	var alertChannelID, environmentID string

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
		CheckDestroy:             base.AlertChannel_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccAlertChannelConfig_Minimal(resourceName, name),
				Check:  base.AlertChannel_GetIDs(resourceFullName, &environmentID, &alertChannelID),
			},
			{
				PreConfig: func() {
					base.AlertChannel_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID, alertChannelID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the environment
			{
				Config: testAccAlertChannel_NewEnv(environmentName, licenseID, resourceName, name),
				Check:  base.AlertChannel_GetIDs(resourceFullName, &environmentID, &alertChannelID),
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

func TestAccAlertChannel_NewEnv(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_alert_channel.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoBeta(t)
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
			resource.TestCheckResourceAttr(resourceFullName, "alert_name", name),
			resource.TestCheckResourceAttr(resourceFullName, "addresses.#", "3"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "addresses.*", "noreply@pingidentity.com"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "addresses.*", "noreply2@pingidentity.com"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "addresses.*", "noreply3@pingidentity.com"),
			resource.TestCheckResourceAttr(resourceFullName, "channel_type", "EMAIL"),
			resource.TestCheckResourceAttr(resourceFullName, "exclude_alert_types.#", "6"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "exclude_alert_types.*", "LICENSE_EXPIRED"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "exclude_alert_types.*", "LICENSE_EXPIRING"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "exclude_alert_types.*", "LICENSE_USER_SOFT_LIMIT_EXCEEDED"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "exclude_alert_types.*", "LICENSE_90_PERCENT_USER_SOFT_LIMIT"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "exclude_alert_types.*", "LICENSE_ROTATED"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "exclude_alert_types.*", "LICENSE_USER_HARD_LIMIT_EXCEEDED"),
			resource.TestCheckResourceAttr(resourceFullName, "include_alert_types.#", "7"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "include_alert_types.*", "CERTIFICATE_EXPIRING"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "include_alert_types.*", "GATEWAY_VERSION_DEPRECATING"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "include_alert_types.*", "KEY_PAIR_EXPIRING"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "include_alert_types.*", "RISK_CONFIGURATION"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "include_alert_types.*", "CERTIFICATE_EXPIRED"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "include_alert_types.*", "GATEWAY_VERSION_DEPRECATED"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "include_alert_types.*", "KEY_PAIR_EXPIRED"),
			resource.TestCheckResourceAttr(resourceFullName, "include_severities.#", "3"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "include_severities.*", "ERROR"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "include_severities.*", "INFO"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "include_severities.*", "WARNING"),
		),
	}

	fullStep2 := resource.TestStep{
		Config: testAccAlertChannelConfig_PartialFull(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "alert_name", name),
			resource.TestCheckResourceAttr(resourceFullName, "addresses.#", "2"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "addresses.*", "noreply2@pingidentity.com"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "addresses.*", "noreply3@pingidentity.com"),
			resource.TestCheckResourceAttr(resourceFullName, "channel_type", "EMAIL"),
			resource.TestCheckResourceAttr(resourceFullName, "exclude_alert_types.#", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "include_alert_types.#", "3"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "include_alert_types.*", "CERTIFICATE_EXPIRING"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "include_alert_types.*", "GATEWAY_VERSION_DEPRECATING"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "include_alert_types.*", "KEY_PAIR_EXPIRING"),
			resource.TestCheckResourceAttr(resourceFullName, "include_severities.#", "1"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "include_severities.*", "INFO"),
		),
	}

	minimalStep := resource.TestStep{
		Config: testAccAlertChannelConfig_Minimal(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckNoResourceAttr(resourceFullName, "alert_name"),
			resource.TestCheckResourceAttr(resourceFullName, "addresses.#", "1"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "addresses.*", "noreply@pingidentity.com"),
			resource.TestCheckResourceAttr(resourceFullName, "channel_type", "EMAIL"),
			resource.TestCheckResourceAttr(resourceFullName, "exclude_alert_types.#", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "include_alert_types.#", "3"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "include_alert_types.*", "LICENSE_EXPIRED"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "include_alert_types.*", "LICENSE_EXPIRING"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "include_alert_types.*", "LICENSE_ROTATED"),
			resource.TestCheckResourceAttr(resourceFullName, "include_severities.#", "1"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "include_severities.*", "ERROR"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
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
							return "", fmt.Errorf("resource not found: %s", resourceFullName)
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
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
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
}`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
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
    "LICENSE_EXPIRED",
    "LICENSE_EXPIRING",
    "LICENSE_USER_SOFT_LIMIT_EXCEEDED",
    "LICENSE_90_PERCENT_USER_SOFT_LIMIT",
    "LICENSE_ROTATED",
    "LICENSE_USER_HARD_LIMIT_EXCEEDED",
  ]

  include_alert_types = [
    "CERTIFICATE_EXPIRING",
    "GATEWAY_VERSION_DEPRECATING",
    "KEY_PAIR_EXPIRING",
    "RISK_CONFIGURATION",
    "CERTIFICATE_EXPIRED",
    "GATEWAY_VERSION_DEPRECATED",
    "KEY_PAIR_EXPIRED",
  ]

  include_severities = [
    "ERROR",
    "INFO",
    "WARNING",
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
    "CERTIFICATE_EXPIRING",
    "GATEWAY_VERSION_DEPRECATING",
    "KEY_PAIR_EXPIRING",
  ]

  include_severities = [
    "INFO",
  ]
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccAlertChannelConfig_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_alert_channel" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

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
    "ERROR",
  ]
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
