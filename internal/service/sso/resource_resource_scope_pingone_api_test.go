package sso_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func testAccCheckResourceScopePingOneAPIDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient

	re, err := regexp.Compile(`^p1:(read|update):user$`)
	if err != nil {
		return fmt.Errorf("Cannot compile regex check for predefined scopes.")
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_resource_scope_pingone_api" {
			continue
		}

		if m := re.MatchString(rs.Primary.Attributes["name"]); m {
			return nil
		} else {

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

			body, r, err := apiClient.ResourcesApi.ReadOneResource(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.ID).Execute()

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

			return fmt.Errorf("PingOne Resource scope Instance %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccGetResourceScopePingOneAPIIDs(resourceName string, environmentID, openidResourceID, resourceID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource Not found: %s", resourceName)
		}

		*resourceID = rs.Primary.ID
		*openidResourceID = rs.Primary.Attributes["resource_id"]
		*environmentID = rs.Primary.Attributes["environment_id"]

		return nil
	}
}

func TestAccResourceScopePingOneAPI_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_resource_scope_pingone_api.%s", resourceName)

	name := resourceName

	var resourceID, openidResourceID, environmentID string

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckResourceScopePingOneAPIDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccResourceScopePingOneAPIConfig_Minimal(resourceName, fmt.Sprintf("p1:read:user:%s", name)),
				Check:  testAccGetResourceScopePingOneAPIIDs(resourceFullName, &environmentID, &openidResourceID, &resourceID),
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

					if environmentID == "" || openidResourceID == "" || resourceID == "" {
						t.Fatalf("One of environment ID, OpenID resource ID or resource ID cannot be determined. Environment ID: %s, OpenID resource ID: %s, Resource ID: %s", environmentID, openidResourceID, resourceID)
					}

					_, err = apiClient.ResourceScopesApi.DeleteResourceScope(ctx, environmentID, openidResourceID, resourceID).Execute()
					if err != nil {
						t.Fatalf("Failed to delete PingOne resource scope: %v", err)
					}
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
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckResourceScopePingOneAPIDestroy,
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
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckResourceScopePingOneAPIDestroy,
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
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckResourceScopePingOneAPIDestroy,
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
		CheckDestroy:             testAccCheckResourceScopePingOneAPIDestroy,
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
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckResourceScopePingOneAPIDestroy,
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
