package sso_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func testAccCheckUserDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_user" {
			continue
		}

		body, r, err := apiClient.UsersApi.ReadUser(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.ID).Execute()

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

		return fmt.Errorf("PingOne User %s still exists", rs.Primary.ID)
	}

	return nil
}

func TestAccUser_NewEnv(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_user.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckUserDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccUserConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "username", name),
				),
			},
		},
	})
}

func TestAccUser_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_user.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckUserDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccUserConfig_Full(resourceName, name, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "username", name),
					resource.TestCheckResourceAttr(resourceFullName, "email", "noreply@pingidentity.com"),
					resource.TestMatchResourceAttr(resourceFullName, "population_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
				),
			},
			{
				Config: testAccUserConfig_Full(resourceName, name, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "username", name),
					resource.TestCheckResourceAttr(resourceFullName, "email", "noreply@pingidentity.com"),
					resource.TestMatchResourceAttr(resourceFullName, "population_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "false"),
				),
			},
			{
				Config: testAccUserConfig_Full(resourceName, name, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "username", name),
					resource.TestCheckResourceAttr(resourceFullName, "email", "noreply@pingidentity.com"),
					resource.TestMatchResourceAttr(resourceFullName, "population_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
				),
			},
		},
	})
}

func TestAccUser_Minimal(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_user.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckUserDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccUserConfig_Minimal(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "username", name),
					resource.TestCheckResourceAttr(resourceFullName, "email", "noreply@pingidentity.com"),
					resource.TestMatchResourceAttr(resourceFullName, "population_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
				),
			},
			{
				Config: testAccUserConfig_Full(resourceName, name, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "username", name),
					resource.TestCheckResourceAttr(resourceFullName, "email", "noreply@pingidentity.com"),
					resource.TestMatchResourceAttr(resourceFullName, "population_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "enabled", "false"),
				),
			},
		},
	})
}

func TestAccUser_ChangePopulation(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_user.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckUserDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccUserConfig_Minimal(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "username", name),
					resource.TestCheckResourceAttr(resourceFullName, "email", "noreply@pingidentity.com"),
					resource.TestMatchResourceAttr(resourceFullName, "population_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "status", "ENABLED"),
				),
			},
			{
				Config: testAccUserConfig_CustomPopulation(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "username", name),
					resource.TestCheckResourceAttr(resourceFullName, "email", "noreply@pingidentity.com"),
					resource.TestMatchResourceAttr(resourceFullName, "population_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttr(resourceFullName, "status", "ENABLED"),
				),
			},
		},
	})
}

func testAccUserConfig_NewEnv(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_population" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  name = "%[4]s"
}

resource "pingone_user" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  username      = "%[4]s"
  email         = "noreply@pingidentity.com"
  population_id = pingone_population.%[3]s.id
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccUserConfig_Full(resourceName, name string, userEnabled bool) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_population" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_user" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  username      = "%[3]s"
  email         = "noreply@pingidentity.com"
  population_id = pingone_population.%[2]s.id
  enabled        = %[4]t
}`, acctest.GenericSandboxEnvironment(), resourceName, name, userEnabled)
}

func testAccUserConfig_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_population" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_user" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  username      = "%[3]s"
  email         = "noreply@pingidentity.com"
  population_id = pingone_population.%[2]s.id
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccUserConfig_CustomPopulation(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_population" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_user" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  username      = "%[3]s"
  email         = "noreply@pingidentity.com"
  population_id = pingone_population.%[2]s.id
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
