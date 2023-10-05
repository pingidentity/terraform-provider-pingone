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

func TestAccResourceScopePingOneAPI_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_resource_scope_pingone_api.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	var resourceScopeID, openidResourceID, environmentID string

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
		CheckDestroy:             sso.ResourceScopePingOneAPI_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccResourceScopePingOneAPIConfig_Minimal(resourceName, fmt.Sprintf("p1:read:user:%s", name)),
				Check:  sso.ResourceScopePingOneAPI_GetIDs(resourceFullName, &environmentID, &openidResourceID, &resourceScopeID),
			},
			{
				PreConfig: func() {
					sso.ResourceScopePingOneAPI_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID, openidResourceID, resourceScopeID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the environment
			{
				SkipFunc: func() (bool, error) { return true, fmt.Errorf("TBC") },
				Config:   testAccResourceScopePingOneAPIConfig_NewEnv(environmentName, licenseID, resourceName, fmt.Sprintf("p1:read:user:%s", name)),
				Check:    sso.ResourceScopePingOneAPI_GetIDs(resourceFullName, &environmentID, &openidResourceID, &resourceScopeID),
			},
			{
				SkipFunc: func() (bool, error) { return true, fmt.Errorf("TBC") },
				PreConfig: func() {
					base.Environment_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccResourceScopePingOneAPI_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_resource_scope_pingone_api.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.ResourceScopePingOneAPI_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceScopePingOneAPIConfig_Full(resourceName, fmt.Sprintf("p1:read:user:%s", name)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "resource_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", fmt.Sprintf("p1:read:user:%s", name)),
					resource.TestCheckResourceAttr(resourceFullName, "description", "My resource scope"),
					resource.TestCheckResourceAttr(resourceFullName, "schema_attributes.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "schema_attributes.*", "name.given"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "schema_attributes.*", "name.family"),
				),
			},
			{
				Config: testAccResourceScopePingOneAPIConfig_Full(resourceName, fmt.Sprintf("p1:update:user:%s", name)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "resource_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", fmt.Sprintf("p1:update:user:%s", name)),
					resource.TestCheckResourceAttr(resourceFullName, "description", "My resource scope"),
					resource.TestCheckResourceAttr(resourceFullName, "schema_attributes.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "schema_attributes.*", "name.given"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "schema_attributes.*", "name.family"),
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

						return fmt.Sprintf("%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceScopePingOneAPI_Minimal(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_resource_scope_pingone_api.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.ResourceScopePingOneAPI_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceScopePingOneAPIConfig_Minimal(resourceName, fmt.Sprintf("p1:read:user:%s", name)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "resource_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", fmt.Sprintf("p1:read:user:%s", name)),
					resource.TestCheckNoResourceAttr(resourceFullName, "description"),
					resource.TestCheckResourceAttr(resourceFullName, "schema_attributes.#", "0"),
				),
			},
			{
				Config: testAccResourceScopePingOneAPIConfig_Minimal(resourceName, fmt.Sprintf("p1:update:user:%s", name)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "resource_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", fmt.Sprintf("p1:update:user:%s", name)),
					resource.TestCheckNoResourceAttr(resourceFullName, "description"),
					resource.TestCheckResourceAttr(resourceFullName, "schema_attributes.#", "0"),
				),
			},
		},
	})
}

func TestAccResourceScopePingOneAPI_Change(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_resource_scope_pingone_api.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.ResourceScopePingOneAPI_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceScopePingOneAPIConfig_Full(resourceName, fmt.Sprintf("p1:read:user:%s", name)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "resource_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", fmt.Sprintf("p1:read:user:%s", name)),
					resource.TestCheckResourceAttr(resourceFullName, "description", "My resource scope"),
					resource.TestCheckResourceAttr(resourceFullName, "schema_attributes.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "schema_attributes.*", "name.given"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "schema_attributes.*", "name.family"),
				),
			},
			{
				Config: testAccResourceScopePingOneAPIConfig_Minimal(resourceName, fmt.Sprintf("p1:read:user:%s", name)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "resource_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", fmt.Sprintf("p1:read:user:%s", name)),
					//resource.TestCheckResourceAttr(resourceFullName, "description", ""),
					resource.TestCheckResourceAttr(resourceFullName, "schema_attributes.#", "0"),
				),
			},
			{
				Config: testAccResourceScopePingOneAPIConfig_Full(resourceName, fmt.Sprintf("p1:read:user:%s", name)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "resource_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", fmt.Sprintf("p1:read:user:%s", name)),
					resource.TestCheckResourceAttr(resourceFullName, "description", "My resource scope"),
					resource.TestCheckResourceAttr(resourceFullName, "schema_attributes.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "schema_attributes.*", "name.given"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "schema_attributes.*", "name.family"),
				),
			},
		},
	})
}

func TestAccResourceScopePingOneAPI_OverridePredefined(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_resource_scope_pingone_api.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.ResourceScopePingOneAPI_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceScopePingOneAPIConfig_OverridePredefined(environmentName, licenseID, resourceName, "p1:read:user"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "resource_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "name", "p1:read:user"),
					resource.TestCheckResourceAttr(resourceFullName, "schema_attributes.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "schema_attributes.*", "name.given"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "schema_attributes.*", "name.family"),
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

						return fmt.Sprintf("%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceScopePingOneAPI_InvalidParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_resource_scope_pingone_api.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.ResourceScopePingOneAPI_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccResourceScopePingOneAPIConfig_Minimal(resourceName, "testscope"),
				ExpectError: regexp.MustCompile("Invalid Attribute Value Match"),
			},
			{
				Config:      testAccResourceScopePingOneAPIConfig_Full(resourceName, "p1:read:user"),
				ExpectError: regexp.MustCompile("Invalid attribute value"),
			},
			// Configure
			{
				Config: testAccResourceScopePingOneAPIConfig_Minimal(resourceName, fmt.Sprintf("p1:read:user:%s", name)),
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

func testAccResourceScopePingOneAPIConfig_NewEnv(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_resource_scope_pingone_api" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  name        = "%[4]s"
  description = "My resource scope"

  schema_attributes = [
    "name.given",
    "name.family",
  ]
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccResourceScopePingOneAPIConfig_Full(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_resource_scope_pingone_api" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name        = "%[3]s"
  description = "My resource scope"

  schema_attributes = [
    "name.given",
    "name.family",
  ]
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccResourceScopePingOneAPIConfig_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_resource_scope_pingone_api" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccResourceScopePingOneAPIConfig_OverridePredefined(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_resource_scope_pingone_api" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  name = "%[4]s"

  schema_attributes = [
    "name.given",
    "name.family",
  ]
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}
