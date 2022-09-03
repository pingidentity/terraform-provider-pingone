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
	pingone "github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
)

func testAccCheckUserDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, pingone.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_user" {
			continue
		}

		body, r, err := apiClient.UsersApi.ReadUser(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.ID).Execute()

		if r.StatusCode == 404 {
			continue
		}

		if err != nil {
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
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckUserDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
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
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckUserDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccUserConfig_Full(resourceName, name, "ENABLED"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "username", name),
					resource.TestCheckResourceAttr(resourceFullName, "email", "noreply@pingidentity.com"),
					resource.TestMatchResourceAttr(resourceFullName, "population_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "status", "ENABLED"),
				),
			},
			{
				Config: testAccUserConfig_Full(resourceName, name, "DISABLED"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "username", name),
					resource.TestCheckResourceAttr(resourceFullName, "email", "noreply@pingidentity.com"),
					resource.TestMatchResourceAttr(resourceFullName, "population_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "status", "DISABLED"),
				),
			},
			{
				Config: testAccUserConfig_Full(resourceName, name, "ENABLED"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "username", name),
					resource.TestCheckResourceAttr(resourceFullName, "email", "noreply@pingidentity.com"),
					resource.TestMatchResourceAttr(resourceFullName, "population_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "status", "ENABLED"),
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
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckUserDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccUserConfig_Minimal(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "username", name),
					resource.TestCheckResourceAttr(resourceFullName, "email", "noreply@pingidentity.com"),
					resource.TestMatchResourceAttr(resourceFullName, "population_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "status", "ENABLED"),
				),
			},
			{
				Config: testAccUserConfig_Full(resourceName, name, "DISABLED"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "username", name),
					resource.TestCheckResourceAttr(resourceFullName, "email", "noreply@pingidentity.com"),
					resource.TestMatchResourceAttr(resourceFullName, "population_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "status", "DISABLED"),
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
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckUserDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccUserConfig_Minimal(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "username", name),
					resource.TestCheckResourceAttr(resourceFullName, "email", "noreply@pingidentity.com"),
					resource.TestMatchResourceAttr(resourceFullName, "population_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "status", "ENABLED"),
				),
			},
			{
				Config: testAccUserConfig_CustomPopulation(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "username", name),
					resource.TestCheckResourceAttr(resourceFullName, "email", "noreply@pingidentity.com"),
					resource.TestMatchResourceAttr(resourceFullName, "population_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
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
			environment_id = "${pingone_environment.%[2]s.id}"

			name = "%[4]s"
		}

		resource "pingone_user" "%[3]s" {
			environment_id = "${pingone_environment.%[2]s.id}"

			username = "%[4]s"
			email = "noreply@pingidentity.com"
			population_id = "${pingone_population.%[3]s.id}"
		}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccUserConfig_Full(resourceName, name, status string) string {
	return fmt.Sprintf(`
		%[1]s

		resource "pingone_population" "%[2]s" {
			environment_id = "${data.pingone_environment.general_test.id}"

			name = "%[3]s"
		}

		resource "pingone_user" "%[2]s" {
			environment_id = "${data.pingone_environment.general_test.id}"

			username = "%[3]s"
			email = "noreply@pingidentity.com"
			population_id = "${pingone_population.%[2]s.id}"
			status = "%[4]s"
		}`, acctest.GenericSandboxEnvironment(), resourceName, name, status)
}

func testAccUserConfig_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

		resource "pingone_population" "%[2]s" {
			environment_id = "${data.pingone_environment.general_test.id}"

			name = "%[3]s"
		}

		resource "pingone_user" "%[2]s" {
			environment_id = "${data.pingone_environment.general_test.id}"

			username = "%[3]s"
			email = "noreply@pingidentity.com"
			population_id = "${pingone_population.%[2]s.id}"
		}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccUserConfig_CustomPopulation(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

		resource "pingone_population" "%[2]s" {
			environment_id = "${data.pingone_environment.general_test.id}"
			
			name = "%[3]s"
		}

		resource "pingone_user" "%[2]s" {
			environment_id = "${data.pingone_environment.general_test.id}"

			username = "%[3]s"
			email = "noreply@pingidentity.com"
			population_id = "${pingone_population.%[2]s.id}"
		}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
