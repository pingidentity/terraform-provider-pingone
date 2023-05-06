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
	"github.com/patrickcping/pingone-go-sdk-v2/management"
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
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})

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

func TestAccResourceScopePingOneAPI_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_resource_scope_pingone_api.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckResourceScopePingOneAPIDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceScopePingOneAPIConfig_Full(resourceName, fmt.Sprintf("p1:read:user:%s", name)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "resource_id", verify.P1ResourceIDRegexp),
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
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "resource_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "name", fmt.Sprintf("p1:update:user:%s", name)),
					resource.TestCheckResourceAttr(resourceFullName, "description", "My resource scope"),
					resource.TestCheckResourceAttr(resourceFullName, "schema_attributes.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "schema_attributes.*", "name.given"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "schema_attributes.*", "name.family"),
				),
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
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckResourceScopePingOneAPIDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceScopePingOneAPIConfig_Minimal(resourceName, fmt.Sprintf("p1:read:user:%s", name)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "resource_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "name", fmt.Sprintf("p1:read:user:%s", name)),
					resource.TestCheckResourceAttr(resourceFullName, "description", ""),
					resource.TestCheckResourceAttr(resourceFullName, "schema_attributes.#", "0"),
				),
			},
			{
				Config: testAccResourceScopePingOneAPIConfig_Minimal(resourceName, fmt.Sprintf("p1:update:user:%s", name)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "resource_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "name", fmt.Sprintf("p1:update:user:%s", name)),
					resource.TestCheckResourceAttr(resourceFullName, "description", ""),
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
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckResourceScopePingOneAPIDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceScopePingOneAPIConfig_Full(resourceName, fmt.Sprintf("p1:read:user:%s", name)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "resource_id", verify.P1ResourceIDRegexp),
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
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "resource_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "name", fmt.Sprintf("p1:read:user:%s", name)),
					//resource.TestCheckResourceAttr(resourceFullName, "description", ""),
					resource.TestCheckResourceAttr(resourceFullName, "schema_attributes.#", "0"),
				),
			},
			{
				Config: testAccResourceScopePingOneAPIConfig_Full(resourceName, fmt.Sprintf("p1:read:user:%s", name)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "resource_id", verify.P1ResourceIDRegexp),
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
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckResourceScopePingOneAPIDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceScopePingOneAPIConfig_OverridePredefined(environmentName, licenseID, resourceName, "p1:read:user"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "resource_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "name", "p1:read:user"),
					resource.TestCheckResourceAttr(resourceFullName, "schema_attributes.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "schema_attributes.*", "name.given"),
					resource.TestCheckTypeSetElemAttr(resourceFullName, "schema_attributes.*", "name.family"),
				),
			},
		},
	})
}

func TestAccResourceScopePingOneAPI_InvalidParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckResourceScopePingOneAPIDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccResourceScopePingOneAPIConfig_Minimal(resourceName, "testscope"),
				ExpectError: regexp.MustCompile("Resource scope name must be either `p1:read:user`, `p1:update:user`, `p1:read:user:{suffix}` or `p1:update:user:{suffix}`"),
			},
			{
				Config:      testAccResourceScopePingOneAPIConfig_Full(resourceName, "p1:read:user"),
				ExpectError: regexp.MustCompile("The scope `p1:read:user` is an existing platform scope.  The description cannot be changed."),
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
