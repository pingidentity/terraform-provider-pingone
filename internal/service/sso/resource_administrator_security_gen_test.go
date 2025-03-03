// Copyright © 2025 Ping Identity Corporation
// Code generated by ping-terraform-plugin-framework-generator

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
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
)

func TestAccAdministratorSecurity_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_administrator_security.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")
	var environmentId string

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
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Test removal of the environment
			{
				Config: administratorSecurity_NewEnvHCL(environmentName, licenseID, resourceName),
				Check:  administratorSecurity_GetIDs(resourceFullName, &environmentId),
			},
			{
				PreConfig: func() {
					base.Environment_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentId)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAdministratorSecurity_MinimalMaximal(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_administrator_security.%s", resourceName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				// Create the resource with a minimal model
				Config: administratorSecurity_MinimalHCL(resourceName),
				Check:  administratorSecurity_CheckComputedValuesMinimal(resourceName),
			},
			{
				// Update to a complete model
				Config: administratorSecurity_CompleteHCL(resourceName),
				Check:  administratorSecurity_CheckComputedValuesComplete(resourceName),
			},
			{
				// Test importing the resource
				Config:       administratorSecurity_CompleteHCL(resourceName),
				ResourceName: fmt.Sprintf("pingone_administrator_security.%s", resourceName),
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s", rs.Primary.Attributes["environment_id"]), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Back to minimal model
				Config: administratorSecurity_MinimalHCL(resourceName),
				Check:  administratorSecurity_CheckComputedValuesMinimal(resourceName),
			},
		},
	})
}

func TestAccAdministratorSecurity_NewEnv(t *testing.T) {
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
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: administratorSecurity_NewEnvHCL(environmentName, licenseID, resourceName),
				Check:  administratorSecurity_CheckComputedValuesMinimal(resourceName),
			},
		},
	})
}
func TestAccAdministratorSecurity_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_administrator_security.%s", resourceName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: administratorSecurity_MinimalHCL(resourceName),
			},
			// Errors
			{
				ResourceName:  resourceFullName,
				ImportStateId: "badformat",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`Unexpected Import Identifier`),
			},
		},
	})
}

// Minimal HCL with only required values set
func administratorSecurity_MinimalHCL(resourceName string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_administrator_security" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  mfa_status = "ENFORCE"
  recovery = true
}
`, acctest.GenericSandboxEnvironment(), resourceName)
}

// Maximal HCL with all values set where possible
func administratorSecurity_CompleteHCL(resourceName string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_identity_provider" "%[2]s-idp" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[2]s-idp"
  enabled                    = true

  microsoft = {
    client_id     = "dummyclientid1"
    client_secret = "dummyclientsecret1"
    tenant_id     = "dummytenantid1"
  }
}

resource "pingone_administrator_security" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  allowed_methods = {
    email = jsonencode(
                {
                 enabled = true
                }
            )
    fido2 = jsonencode(
                {
                 enabled = false
                }
            )
    totp = jsonencode(
                {
                 enabled = false
                }
            )
  }
  authentication_method = "HYBRID"
  mfa_status = "ENFORCE"
  identity_provider = {
    id = pingone_identity_provider.%[2]s-idp.id
  }
  recovery = false
}
`, acctest.GenericSandboxEnvironment(), resourceName)
}

func administratorSecurity_NewEnvHCL(environmentName, licenseID, resourceName string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_administrator_security" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  mfa_status = "ENFORCE"
  recovery = true
}
`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName)
}

// Validate any computed values when applying minimal HCL
func administratorSecurity_CheckComputedValuesMinimal(resourceName string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(fmt.Sprintf("pingone_administrator_security.%s", resourceName), "allowed_methods.email", `{"enabled":true}`),
		resource.TestCheckResourceAttr(fmt.Sprintf("pingone_administrator_security.%s", resourceName), "allowed_methods.fido2", `{"enabled":true}`),
		resource.TestCheckResourceAttr(fmt.Sprintf("pingone_administrator_security.%s", resourceName), "allowed_methods.totp", `{"enabled":true}`),
		resource.TestCheckResourceAttr(fmt.Sprintf("pingone_administrator_security.%s", resourceName), "authentication_method", "PINGONE"),
		resource.TestCheckResourceAttr(fmt.Sprintf("pingone_administrator_security.%s", resourceName), "has_fido2_capabilities", "true"),
		resource.TestCheckResourceAttr(fmt.Sprintf("pingone_administrator_security.%s", resourceName), "is_pingid_in_bom", "false"),
	)
}

// Validate any computed values when applying complete HCL
func administratorSecurity_CheckComputedValuesComplete(resourceName string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(fmt.Sprintf("pingone_administrator_security.%s", resourceName), "has_fido2_capabilities", "true"),
		resource.TestCheckResourceAttr(fmt.Sprintf("pingone_administrator_security.%s", resourceName), "is_pingid_in_bom", "false"),
	)
}

func administratorSecurity_GetIDs(resourceName string, environmentId *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource Not found: %s", resourceName)
		}
		if environmentId != nil {
			*environmentId = rs.Primary.Attributes["environment_id"]
		}

		return nil
	}
}
