package sso_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	pingone "github.com/patrickcping/pingone-go/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
)

func testAccCheckUserDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API
	ctx = context.WithValue(ctx, pingone.ContextServerVariables, map[string]string{
		"suffix": p1Client.RegionSuffix,
	})

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_user" {
			continue
		}

		body, r, err := apiClient.UsersUsersApi.ReadUser(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.ID).Execute()

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

func TestAccUser_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_user.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")
	region := os.Getenv("PINGONE_REGION")

	email := "noreply@pingidentity.com"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckUserDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccUserConfig_Full(environmentName, licenseID, region, resourceName, email, "ENABLED"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, "id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "environment_id"),
					resource.TestCheckResourceAttr(resourceFullName, "username", name),
					resource.TestCheckResourceAttr(resourceFullName, "email", email),
					resource.TestCheckResourceAttrSet(resourceFullName, "population_id"),
					resource.TestCheckResourceAttr(resourceFullName, "status", "ENABLED"),
				),
			},
			{
				Config: testAccUserConfig_Full(environmentName, licenseID, region, resourceName, email, "DISABLED"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, "id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "environment_id"),
					resource.TestCheckResourceAttr(resourceFullName, "username", name),
					resource.TestCheckResourceAttr(resourceFullName, "email", email),
					resource.TestCheckResourceAttrSet(resourceFullName, "population_id"),
					resource.TestCheckResourceAttr(resourceFullName, "status", "DISABLED"),
				),
			},
			{
				Config: testAccUserConfig_Full(environmentName, licenseID, region, resourceName, email, "ENABLED"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, "id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "environment_id"),
					resource.TestCheckResourceAttr(resourceFullName, "username", name),
					resource.TestCheckResourceAttr(resourceFullName, "email", email),
					resource.TestCheckResourceAttrSet(resourceFullName, "population_id"),
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

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")
	region := os.Getenv("PINGONE_REGION")

	email := "noreply@pingidentity.com"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckUserDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccUserConfig_Minimal(environmentName, licenseID, region, resourceName, email),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, "id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "environment_id"),
					resource.TestCheckResourceAttr(resourceFullName, "username", name),
					resource.TestCheckResourceAttr(resourceFullName, "email", email),
					resource.TestCheckResourceAttrSet(resourceFullName, "population_id"),
					resource.TestCheckResourceAttr(resourceFullName, "status", "ENABLED"),
				),
			},
			{
				Config: testAccUserConfig_Full(environmentName, licenseID, region, resourceName, email, "DISABLED"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, "id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "environment_id"),
					resource.TestCheckResourceAttr(resourceFullName, "username", name),
					resource.TestCheckResourceAttr(resourceFullName, "email", email),
					resource.TestCheckResourceAttrSet(resourceFullName, "population_id"),
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

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")
	region := os.Getenv("PINGONE_REGION")

	email := "noreply@pingidentity.com"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckUserDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccUserConfig_Minimal(environmentName, licenseID, region, resourceName, email),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, "id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "environment_id"),
					resource.TestCheckResourceAttr(resourceFullName, "username", name),
					resource.TestCheckResourceAttr(resourceFullName, "email", email),
					resource.TestCheckResourceAttrSet(resourceFullName, "population_id"),
					resource.TestCheckResourceAttr(resourceFullName, "status", "ENABLED"),
				),
			},
			{
				Config: testAccUserConfig_CustomPopulation(environmentName, licenseID, region, resourceName, email),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, "id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "environment_id"),
					resource.TestCheckResourceAttr(resourceFullName, "username", name),
					resource.TestCheckResourceAttr(resourceFullName, "email", email),
					resource.TestCheckResourceAttrSet(resourceFullName, "population_id"),
					resource.TestCheckResourceAttr(resourceFullName, "status", "ENABLED"),
				),
			},
		},
	})
}

func testAccUserConfig_Full(environmentName, licenseID, region, resourceName, email, status string) string {
	return fmt.Sprintf(`
		resource "pingone_environment" "%[1]s" {
			name = "%[1]s"
			type = "SANDBOX"
			license_id = "%[2]s"
			region = "%[3]s"
			default_population {}
			service {}
		}

		resource "pingone_user" "%[4]s" {
			environment_id = "${pingone_environment.%[1]s.id}"

			username = "%[4]s"
			email = "%[5]s"
			population_id = "${pingone_environment.%[1]s.default_population_id}"
			status = "%[6]s"
		}`, environmentName, licenseID, region, resourceName, email, status)
}

func testAccUserConfig_Minimal(environmentName, licenseID, region, resourceName, email string) string {
	return fmt.Sprintf(`
		resource "pingone_environment" "%[1]s" {
			name = "%[1]s"
			type = "SANDBOX"
			license_id = "%[2]s"
			region = "%[3]s"
			default_population {}
			service {}
		}

		resource "pingone_user" "%[4]s" {
			environment_id = "${pingone_environment.%[1]s.id}"

			username = "%[4]s"
			email = "%[5]s"
			population_id = "${pingone_environment.%[1]s.default_population_id}"
		}`, environmentName, licenseID, region, resourceName, email)
}

func testAccUserConfig_CustomPopulation(environmentName, licenseID, region, resourceName, email string) string {
	return fmt.Sprintf(`
		resource "pingone_environment" "%[1]s" {
			name = "%[1]s"
			type = "SANDBOX"
			license_id = "%[2]s"
			region = "%[3]s"
			default_population {}
			service {}
		}

		resource "pingone_population" "%[4]s" {
			environment_id = "${pingone_environment.%[1]s.id}"
			
			name = "%[4]s"
		}

		resource "pingone_user" "%[4]s" {
			environment_id = "${pingone_environment.%[1]s.id}"

			username = "%[4]s"
			email = "%[5]s"
			population_id = "${pingone_population.%[4]s.id}"
		}`, environmentName, licenseID, region, resourceName, email)
}
