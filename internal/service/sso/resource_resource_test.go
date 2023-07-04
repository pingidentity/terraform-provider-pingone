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

func testAccCheckResourceDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_resource" {
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

		return fmt.Errorf("PingOne Resource Instance %s still exists", rs.Primary.ID)
	}

	return nil
}

func TestAccResource_NewEnv(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_resource.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckResourceDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
				),
			},
		},
	})
}

func TestAccResource_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_resource.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckResourceDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceConfig_Full(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", "Test Resource"),
					resource.TestCheckResourceAttr(resourceFullName, "type", "CUSTOM"),
					resource.TestCheckResourceAttr(resourceFullName, "audience", fmt.Sprintf("%s-1", name)),
					resource.TestCheckResourceAttr(resourceFullName, "access_token_validity_seconds", "7200"),
					resource.TestCheckResourceAttr(resourceFullName, "introspect_endpoint_auth_method", "CLIENT_SECRET_POST"),
					resource.TestMatchResourceAttr(resourceFullName, "client_secret", regexp.MustCompile(`[a-zA-Z0-9-~_]{10,}`)),
				),
			},
		},
	})
}

func TestAccResource_Minimal(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_resource.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckResourceDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceConfig_Minimal(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", ""),
					resource.TestCheckResourceAttr(resourceFullName, "type", "CUSTOM"),
					resource.TestCheckResourceAttr(resourceFullName, "audience", name),
					resource.TestCheckResourceAttr(resourceFullName, "access_token_validity_seconds", "3600"),
					resource.TestCheckResourceAttr(resourceFullName, "introspect_endpoint_auth_method", "CLIENT_SECRET_BASIC"),
					resource.TestMatchResourceAttr(resourceFullName, "client_secret", regexp.MustCompile(`[a-zA-Z0-9-~_]{10,}`)),
				),
			},
		},
	})
}

func TestAccResource_Change(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_resource.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckResourceDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceConfig_Minimal(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", ""),
					resource.TestCheckResourceAttr(resourceFullName, "type", "CUSTOM"),
					resource.TestCheckResourceAttr(resourceFullName, "audience", name),
					resource.TestCheckResourceAttr(resourceFullName, "access_token_validity_seconds", "3600"),
					resource.TestCheckResourceAttr(resourceFullName, "introspect_endpoint_auth_method", "CLIENT_SECRET_BASIC"),
					resource.TestMatchResourceAttr(resourceFullName, "client_secret", regexp.MustCompile(`[a-zA-Z0-9-~_]{10,}`)),
				),
			},
			{
				Config: testAccResourceConfig_Full(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", "Test Resource"),
					resource.TestCheckResourceAttr(resourceFullName, "type", "CUSTOM"),
					resource.TestCheckResourceAttr(resourceFullName, "audience", fmt.Sprintf("%s-1", name)),
					resource.TestCheckResourceAttr(resourceFullName, "access_token_validity_seconds", "7200"),
					resource.TestCheckResourceAttr(resourceFullName, "introspect_endpoint_auth_method", "CLIENT_SECRET_POST"),
					resource.TestMatchResourceAttr(resourceFullName, "client_secret", regexp.MustCompile(`[a-zA-Z0-9-~_]{10,}`)),
				),
			},
			{
				Config: testAccResourceConfig_Minimal(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", ""),
					resource.TestCheckResourceAttr(resourceFullName, "type", "CUSTOM"),
					resource.TestCheckResourceAttr(resourceFullName, "audience", name),
					resource.TestCheckResourceAttr(resourceFullName, "access_token_validity_seconds", "3600"),
					resource.TestCheckResourceAttr(resourceFullName, "introspect_endpoint_auth_method", "CLIENT_SECRET_BASIC"),
					resource.TestMatchResourceAttr(resourceFullName, "client_secret", regexp.MustCompile(`[a-zA-Z0-9-~_]{10,}`)),
				),
			},
		},
	})
}

func testAccResourceConfig_NewEnv(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_resource" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  name = "%[4]s"
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccResourceConfig_Full(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_resource" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name        = "%[3]s"
  description = "Test Resource"

  audience                      = "%[3]s-1"
  access_token_validity_seconds = 7200

  introspect_endpoint_auth_method = "CLIENT_SECRET_POST"
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccResourceConfig_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_resource" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
