// Copyright © 2025 Ping Identity Corporation

//go:build beta

package davinci_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/pingidentity/pingone-go-client/pingone"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	acctestlegacysdk "github.com/pingidentity/terraform-provider-pingone/internal/acctest/legacysdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/base"
)

var (
	currentApiKey string
)

func TestAccDavinciApplicationKey_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_davinci_application_key.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")
	var environmentId string
	var id string

	var p1Client *pingone.APIClient
	var ctx = context.Background()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)

			p1Client = acctest.PreCheckTestClient(ctx, t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             davinciApplication_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: davinciApplicationKey_FirstRotateHCL(resourceName),
				Check:  davinciApplicationKey_GetIDs(resourceFullName, &environmentId, &id),
			},
			{
				PreConfig: func() {
					davinciApplication_Delete(ctx, p1Client, t, environmentId, id)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the environment
			{
				Config: davinciApplicationKey_NewEnvHCL(environmentName, licenseID, resourceName),
				Check:  davinciApplicationKey_GetIDs(resourceFullName, &environmentId, &id),
			},
			{
				PreConfig: func() {
					base.Environment_RemovalDrift_PreConfig(ctx, p1Client, t, environmentId)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccDavinciApplicationKey_Rotate(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_davinci_application_key.%s", resourceName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             davinciApplication_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				// Create the application
				Config: davinciApplicationKey_ApplicationOnlyHCL(resourceName),
				Check:  davinciApplicationKey_GetApplicationApiKey(fmt.Sprintf("pingone_davinci_application.%s", resourceName), &currentApiKey),
			},
			{
				// Initial rotation on create
				Config: davinciApplicationKey_FirstRotateHCL(resourceName),
				Check: resource.ComposeTestCheckFunc(
					davinciApplicationKey_checkExpectedKey(t, resourceName, true),
					davinciApplicationKey_CheckComputedValues(resourceName),
					davinciApplicationKey_GetApplicationApiKey(resourceFullName, &currentApiKey),
				),
			},
			{
				// Expect no additional rotation
				Config: davinciApplicationKey_FirstNoRotateHCL(resourceName),
				Check: resource.ComposeTestCheckFunc(
					davinciApplicationKey_checkExpectedKey(t, resourceName, false),
					davinciApplicationKey_CheckComputedValues(resourceName),
					davinciApplicationKey_GetApplicationApiKey(resourceFullName, &currentApiKey),
				),
			},
			{
				// Expect rotation
				Config: davinciApplicationKey_SecondRotateHCL(resourceName),
				Check: resource.ComposeTestCheckFunc(
					davinciApplicationKey_checkExpectedKey(t, resourceName, true),
					davinciApplicationKey_CheckComputedValues(resourceName),
					davinciApplicationKey_GetApplicationApiKey(resourceFullName, &currentApiKey),
				),
			},
			{
				// Expect no additional rotation
				Config: davinciApplicationKey_SecondNoRotateHCL(resourceName),
				Check: resource.ComposeTestCheckFunc(
					davinciApplicationKey_checkExpectedKey(t, resourceName, false),
					davinciApplicationKey_CheckComputedValues(resourceName),
					davinciApplicationKey_GetApplicationApiKey(resourceFullName, &currentApiKey),
				),
			},
			{
				// Test importing the resource
				Config:       davinciApplicationKey_SecondRotateHCL(resourceName),
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["id"]), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
				// The rotation trigger values are terraform-only, so they can't be imported
				ImportStateVerifyIgnore: []string{"rotation_trigger_values"},
			},
		},
	})
}

func TestAccDavinciApplicationKey_NewEnv(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             davinciApplication_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: davinciApplicationKey_NewEnvHCL(environmentName, licenseID, resourceName),
				Check:  davinciApplicationKey_CheckComputedValues(resourceName),
			},
		},
	})
}

func TestAccDavinciApplicationKey_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_davinci_application_key.%s", resourceName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             davinciApplication_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: davinciApplicationKey_FirstRotateHCL(resourceName),
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

func davinciApplicationKey_GetIDs(resourceName string, environmentId, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource Not found: %s", resourceName)
		}
		if environmentId != nil {
			*environmentId = rs.Primary.Attributes["environment_id"]
		}
		if id != nil {
			*id = rs.Primary.Attributes["id"]
		}

		return nil
	}
}

func davinciApplicationKey_GetApplicationApiKey(resourceFullName string, apiKey *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceFullName]
		if !ok {
			return fmt.Errorf("Resource Not found: %s", resourceFullName)
		}
		if apiKey != nil {
			*apiKey = rs.Primary.Attributes["api_key.value"]
		}

		return nil
	}
}

func davinciApplicationKey_CheckComputedValues(resourceName string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(fmt.Sprintf("pingone_davinci_application_key.%s", resourceName), "api_key.enabled", "true"),
		resource.TestCheckResourceAttrSet(fmt.Sprintf("pingone_davinci_application_key.%s", resourceName), "api_key.value"),
		resource.TestCheckResourceAttrSet(fmt.Sprintf("pingone_davinci_application_key.%s", resourceName), "id"),
	)
}

func davinciApplicationKey_checkExpectedKey(_ *testing.T, resourceName string, expectRotation bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		err := resource.TestCheckResourceAttr(fmt.Sprintf("pingone_davinci_application_key.%s", resourceName), "api_key.value", currentApiKey)(s)
		if err != nil && !expectRotation {
			return err
		} else if err == nil && expectRotation {
			return errors.New("Expected the current API key to have rotated, but the key has not changed")
		}
		return nil
	}
}

func davinciApplicationKey_ApplicationOnlyHCL(resourceName string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_davinci_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[2]s"
}
`, acctest.GenericSandboxEnvironment(), resourceName)
}

func davinciApplicationKey_FirstRotateHCL(resourceName string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_davinci_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[2]s"
}

resource "pingone_davinci_application_key" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  application_id = pingone_davinci_application.%[2]s.id
}
`, acctest.GenericSandboxEnvironment(), resourceName)
}

// Ensure that adding triggers doesn't cause a rotation
func davinciApplicationKey_FirstNoRotateHCL(resourceName string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_davinci_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[2]s"
}

resource "pingone_davinci_application_key" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  application_id = pingone_davinci_application.%[2]s.id
  rotation_trigger_values = {
    "trigger" = "initial"
  }
}
`, acctest.GenericSandboxEnvironment(), resourceName)
}

func davinciApplicationKey_SecondRotateHCL(resourceName string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_davinci_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[2]s"
}

resource "pingone_davinci_application_key" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  application_id = pingone_davinci_application.%[2]s.id
  rotation_trigger_values = {
    "trigger"    = "updated"
    "newtrigger" = "new"
  }
}
`, acctest.GenericSandboxEnvironment(), resourceName)
}

// Ensure that removing triggers doesn't cause a rotation
func davinciApplicationKey_SecondNoRotateHCL(resourceName string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_davinci_application" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[2]s"
}

resource "pingone_davinci_application_key" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  application_id = pingone_davinci_application.%[2]s.id
  rotation_trigger_values = {
    "trigger" = "updated"
  }
}
`, acctest.GenericSandboxEnvironment(), resourceName)
}

func davinciApplicationKey_NewEnvHCL(environmentName, licenseID, resourceName string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_davinci_application" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[3]s"
}

resource "pingone_davinci_application_key" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  application_id = pingone_davinci_application.%[3]s.id
  rotation_trigger_values = {
    "trigger" = "initial"
  }
}
`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName)
}
