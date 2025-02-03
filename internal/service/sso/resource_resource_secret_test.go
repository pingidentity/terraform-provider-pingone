// Copyright © 2025 Ping Identity Corporation

package sso_test

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
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/sso"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccResourceSecret_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_resource_secret.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	var resourceID, environmentID string

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
		CheckDestroy:             sso.ResourceSecret_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Test removal of the resource
			{
				Config: testAccResourceSecretConfig_Full(resourceName, name),
				Check:  sso.ResourceSecret_GetIDs(resourceFullName, &environmentID, &resourceID),
			},
			{
				PreConfig: func() {
					sso.Resource_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID, resourceID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the environment
			{
				Config: testAccResourceSecretConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check:  sso.ResourceSecret_GetIDs(resourceFullName, &environmentID, &resourceID),
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

func TestAccResourceSecret_Basic(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_resource_secret.%s", resourceName)

	name := resourceName

	fullStep := resource.TestStep{
		Config: testAccResourceSecretConfig_Full(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "resource_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckNoResourceAttr(resourceFullName, "previous.secret"),
			resource.TestCheckNoResourceAttr(resourceFullName, "previous.expires_at"),
			resource.TestCheckNoResourceAttr(resourceFullName, "previous.last_used"),
			resource.TestMatchResourceAttr(resourceFullName, "secret", regexp.MustCompile(`[a-zA-Z0-9-~_]{10,}`)),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.ResourceSecret_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Single from new
			fullStep,
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["resource_id"]), nil
					}
				}(),
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "resource_id",
			},
		},
	})
}

func TestAccResourceSecret_Rotation(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_resource_secret.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {
				VersionConstraint: "0.11.1",
				Source:            "hashicorp/time",
			},
		},
		CheckDestroy: sso.ResourceSecret_CheckDestroy,
		ErrorCheck:   acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Single from new
			{
				Config: testAccResourceSecretConfig_Rotation1(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr(resourceFullName, "previous.secret"),
					resource.TestCheckNoResourceAttr(resourceFullName, "previous.expires_at"),
					resource.TestCheckNoResourceAttr(resourceFullName, "previous.last_used"),
					resource.TestMatchResourceAttr(resourceFullName, "secret", regexp.MustCompile(`[a-zA-Z0-9-~_]{10,}`)),
				),
			},
			{
				Config: testAccResourceSecretConfig_Rotation2(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "previous.secret", regexp.MustCompile(`[a-zA-Z0-9-~_]{10,}`)),
					resource.TestMatchResourceAttr(resourceFullName, "previous.expires_at", verify.RFC3339Regexp),
					resource.TestCheckNoResourceAttr(resourceFullName, "previous.last_used"),
					resource.TestMatchResourceAttr(resourceFullName, "secret", regexp.MustCompile(`[a-zA-Z0-9-~_]{10,}`)),
				),
			},
			{
				Config:  testAccResourceSecretConfig_Rotation2(resourceName, name),
				Destroy: true,
			},
			{
				Config: testAccResourceSecretConfig_Rotation2(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "previous.secret", regexp.MustCompile(`[a-zA-Z0-9-~_]{10,}`)),
					resource.TestMatchResourceAttr(resourceFullName, "previous.expires_at", verify.RFC3339Regexp),
					resource.TestCheckNoResourceAttr(resourceFullName, "previous.last_used"),
					resource.TestMatchResourceAttr(resourceFullName, "secret", regexp.MustCompile(`[a-zA-Z0-9-~_]{10,}`)),
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

						return fmt.Sprintf("%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["resource_id"]), nil
					}
				}(),
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "resource_id",
			},
		},
	})
}

func TestAccResourceSecret_ReplaceTriggers(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_resource_secret.%s", resourceName)

	name := resourceName

	triggerAStep1 := resource.TestStep{
		Config: testAccResourceSecretConfig_ReplaceTriggerA1(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "regenerate_trigger_values.triggerA", "triggerAValue1"),
			resource.TestCheckNoResourceAttr(resourceFullName, "regenerate_trigger_values.triggerB"),
		),
	}

	triggerAStep2 := resource.TestStep{
		Config: testAccResourceSecretConfig_ReplaceTriggerA2(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "regenerate_trigger_values.triggerA", "triggerAValue2"),
			resource.TestCheckNoResourceAttr(resourceFullName, "regenerate_trigger_values.triggerB"),
		),
	}

	addTriggerBStep := resource.TestStep{
		Config: testAccResourceSecretConfig_ReplaceTriggerB1(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "regenerate_trigger_values.triggerA", "triggerAValue2"),
			resource.TestCheckResourceAttr(resourceFullName, "regenerate_trigger_values.triggerB", "triggerBValue1"),
		),
	}

	removeTriggerBStep := resource.TestStep{
		Config: testAccResourceSecretConfig_ReplaceTriggerB2(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "regenerate_trigger_values.triggerA", "triggerAValue2"),
			resource.TestCheckNoResourceAttr(resourceFullName, "regenerate_trigger_values.triggerB"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.ResourceSecret_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			triggerAStep1,
			triggerAStep2,
			addTriggerBStep,
			removeTriggerBStep,
		},
	})
}

func TestAccResourceSecret_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_resource_secret.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.ResourceSecret_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccResourceSecretConfig_Full(resourceName, name),
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
				ImportStateId: "badformat/badformat/badformat",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`Unexpected Import Identifier`),
			},
		},
	})
}

func testAccResourceSecretConfig_NewEnv(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_resource" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  name = "%[4]s"
}

resource "pingone_resource_secret" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  resource_id    = pingone_resource.%[3]s.id
}
`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccResourceSecretConfig_Full(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_resource" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_resource_secret" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  resource_id    = pingone_resource.%[2]s.id
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccResourceSecretConfig_Rotation1(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_resource" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_resource_secret" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  resource_id    = pingone_resource.%[2]s.id
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccResourceSecretConfig_Rotation2(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_resource" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "time_offset" "%[2]s" {
  offset_minutes = 10
}

resource "pingone_resource_secret" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  resource_id    = pingone_resource.%[2]s.id

  previous = {
    expires_at = time_offset.%[2]s.rfc3339
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccResourceSecretConfig_ReplaceTriggerA1(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_resource" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_resource_secret" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  resource_id    = pingone_resource.%[2]s.id

  regenerate_trigger_values = {
    "triggerA" : "triggerAValue1",
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccResourceSecretConfig_ReplaceTriggerA2(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_resource" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_resource_secret" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  resource_id    = pingone_resource.%[2]s.id

  regenerate_trigger_values = {
    "triggerA" : "triggerAValue2",
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccResourceSecretConfig_ReplaceTriggerB1(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_resource" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_resource_secret" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  resource_id    = pingone_resource.%[2]s.id

  regenerate_trigger_values = {
    "triggerA" : "triggerAValue2",
    "triggerB" : "triggerBValue1",
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccResourceSecretConfig_ReplaceTriggerB2(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_resource" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
}

resource "pingone_resource_secret" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  resource_id = pingone_resource.%[2]s.id

  regenerate_trigger_values = {
    "triggerA" : "triggerAValue2",
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
