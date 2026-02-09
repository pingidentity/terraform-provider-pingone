// Copyright © 2026 Ping Identity Corporation

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

func TestAccRateLimitConfiguration_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_rate_limit_configuration.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	var rateLimitConfigurationID, environmentID string

	var p1Client *client.Client
	var ctx = context.Background()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoBeta(t)
			p1Client = acctestlegacysdk.PreCheckTestClient(ctx, t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.RateLimitConfiguration_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Test removal of the resource
			{
				Config: testAccRateLimitConfigurationConfig_RemovalDrift(resourceName),
				Check:  base.RateLimitConfiguration_GetIDs(resourceFullName, &environmentID, &rateLimitConfigurationID),
			},
			{
				PreConfig: func() {
					base.RateLimitConfiguration_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID, rateLimitConfigurationID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the environment
			{
				Config: testAccRateLimitConfigurationConfig_IPv4_NewEnv(environmentName, licenseID, resourceName),
				Check:  base.RateLimitConfiguration_GetIDs(resourceFullName, &environmentID, &rateLimitConfigurationID),
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

func TestAccRateLimitConfiguration_NewEnv(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_rate_limit_configuration.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.RateLimitConfiguration_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRateLimitConfigurationConfig_IPv4_NewEnv(environmentName, licenseID, resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
				),
			},
		},
	})
}

func TestAccRateLimitConfiguration_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_rate_limit_configuration.%s", resourceName)

	ipv4Step := resource.TestStep{
		Config: testAccRateLimitConfigurationConfig_IPv4(resourceName),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "type", "WHITELIST"),
			resource.TestCheckResourceAttr(resourceFullName, "value", "192.0.2.1"),
		),
	}

	ipv4CIDRStep := resource.TestStep{
		Config: testAccRateLimitConfigurationConfig_IPv4CIDR(resourceName),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "type", "WHITELIST"),
			resource.TestCheckResourceAttr(resourceFullName, "value", "198.51.100.0/28"),
		),
	}

	ipv6Step := resource.TestStep{
		Config: testAccRateLimitConfigurationConfig_IPv6(resourceName),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "type", "WHITELIST"),
			resource.TestCheckResourceAttr(resourceFullName, "value", "2001:0DB8:0000:0001:0000:0000:0000:0001"),
		),
	}

	ipv6CIDRStep := resource.TestStep{
		Config: testAccRateLimitConfigurationConfig_IPv6CIDR(resourceName),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "type", "WHITELIST"),
			resource.TestCheckResourceAttr(resourceFullName, "value", "2001:0DB8:0001:0000:0000:0000:0000:0000/48"),
		),
	}

	defaultTypeStep := resource.TestStep{
		Config: testAccRateLimitConfigurationConfig_DefaultType(resourceName),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "type", "WHITELIST"),
			resource.TestCheckResourceAttr(resourceFullName, "value", "203.0.113.0/28"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.RateLimitConfiguration_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			ipv4Step,
			{
				Config:  testAccRateLimitConfigurationConfig_IPv4(resourceName),
				Destroy: true,
			},
			ipv4CIDRStep,
			{
				Config:  testAccRateLimitConfigurationConfig_IPv4CIDR(resourceName),
				Destroy: true,
			},
			ipv6Step,
			{
				Config:  testAccRateLimitConfigurationConfig_IPv6(resourceName),
				Destroy: true,
			},
			ipv6CIDRStep,
			{
				Config:  testAccRateLimitConfigurationConfig_IPv6CIDR(resourceName),
				Destroy: true,
			},
			defaultTypeStep,
			{
				Config:  testAccRateLimitConfigurationConfig_DefaultType(resourceName),
				Destroy: true,
			},
			ipv4Step,
			ipv4CIDRStep,
			ipv6Step,
			ipv6CIDRStep,
			defaultTypeStep,
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

func TestAccRateLimitConfiguration_ValidationChecks(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.RateLimitConfiguration_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccRateLimitConfigurationConfig_InvalidIPv4(resourceName),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Match`),
			},
			{
				Config:      testAccRateLimitConfigurationConfig_InvalidIPv6(resourceName),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Match`),
			},
			{
				Config:      testAccRateLimitConfigurationConfig_InvalidCIDR(resourceName),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Match`),
			},
		},
	})
}

func TestAccRateLimitConfiguration_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_rate_limit_configuration.%s", resourceName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.RateLimitConfiguration_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccRateLimitConfigurationConfig_BadParameters(resourceName),
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

func testAccRateLimitConfigurationConfig_IPv4_NewEnv(environmentName, licenseID, resourceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_rate_limit_configuration" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  type  = "WHITELIST"
  value = "192.168.1.1"
}`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName)
}

func testAccRateLimitConfigurationConfig_IPv4(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_rate_limit_configuration" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  type  = "WHITELIST"
  value = "192.0.2.1"
}`, acctest.GenericSandboxEnvironment(), resourceName)
}

func testAccRateLimitConfigurationConfig_IPv4CIDR(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_rate_limit_configuration" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  type  = "WHITELIST"
  value = "198.51.100.0/28"
}`, acctest.GenericSandboxEnvironment(), resourceName)
}

func testAccRateLimitConfigurationConfig_IPv6(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_rate_limit_configuration" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  type  = "WHITELIST"
  value = "2001:0DB8:0000:0001:0000:0000:0000:0001"
}`, acctest.GenericSandboxEnvironment(), resourceName)
}

func testAccRateLimitConfigurationConfig_IPv6CIDR(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_rate_limit_configuration" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  type  = "WHITELIST"
  value = "2001:0DB8:0001:0000:0000:0000:0000:0000/48"
}`, acctest.GenericSandboxEnvironment(), resourceName)
}

func testAccRateLimitConfigurationConfig_DefaultType(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_rate_limit_configuration" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  value = "203.0.113.0/28"
}`, acctest.GenericSandboxEnvironment(), resourceName)
}

func testAccRateLimitConfigurationConfig_InvalidIPv4(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_rate_limit_configuration" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  type  = "WHITELIST"
  value = "not-an-ip-address"
}`, acctest.GenericSandboxEnvironment(), resourceName)
}

func testAccRateLimitConfigurationConfig_InvalidIPv6(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_rate_limit_configuration" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  type  = "WHITELIST"
  value = "ZZZZ:ZZZZ:ZZZZ:ZZZZ:ZZZZ:ZZZZ:ZZZZ:ZZZZ"
}`, acctest.GenericSandboxEnvironment(), resourceName)
}

func testAccRateLimitConfigurationConfig_InvalidCIDR(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_rate_limit_configuration" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  type  = "WHITELIST"
  value = "10.0.0.0/999"
}`, acctest.GenericSandboxEnvironment(), resourceName)
}

func testAccRateLimitConfigurationConfig_RemovalDrift(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_rate_limit_configuration" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  type  = "WHITELIST"
  value = "192.0.2.10"
}`, acctest.GenericSandboxEnvironment(), resourceName)
}

func testAccRateLimitConfigurationConfig_BadParameters(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_rate_limit_configuration" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  type  = "WHITELIST"
  value = "192.0.2.20"
}`, acctest.GenericSandboxEnvironment(), resourceName)
}
