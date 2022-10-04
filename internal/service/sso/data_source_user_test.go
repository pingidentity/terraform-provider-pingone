package sso_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
)

func TestAccUserDataSource_ByNameFull(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_user.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccUserDataSourceConfig_ByNameFull(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "user_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(dataSourceFullName, "username", name),
					resource.TestCheckResourceAttr(dataSourceFullName, "email", fmt.Sprintf("%s@pingidentity.com", name)),
					resource.TestCheckResourceAttr(dataSourceFullName, "status", "OK"),
					resource.TestMatchResourceAttr(dataSourceFullName, "population_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
				),
			},
		},
	})
}

func TestAccUserDataSource_ByEmailFull(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_user.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccUserDataSourceConfig_ByEmailFull(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "user_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(dataSourceFullName, "username", name),
					resource.TestCheckResourceAttr(dataSourceFullName, "email", fmt.Sprintf("%s@pingidentity.com", name)),
					resource.TestCheckResourceAttr(dataSourceFullName, "status", "OK"),
					resource.TestMatchResourceAttr(dataSourceFullName, "population_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
				),
			},
		},
	})
}

func TestAccUserDataSource_ByIDFull(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_user.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccUserDataSourceConfig_ByIDFull(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "user_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(dataSourceFullName, "username", name),
					resource.TestCheckResourceAttr(dataSourceFullName, "email", fmt.Sprintf("%s@pingidentity.com", name)),
					resource.TestCheckResourceAttr(dataSourceFullName, "status", "OK"),
					resource.TestMatchResourceAttr(dataSourceFullName, "population_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
				),
			},
		},
	})
}

func testAccUserDataSourceConfig_ByNameFull(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

		resource "pingone_population" "%[2]s" {
			environment_id = data.pingone_environment.general_test.id
		  
			name = "%[3]s"
		  }
	
		resource "pingone_user" "%[2]s" {
			environment_id = data.pingone_environment.general_test.id
		  
			username      = "%[3]s"
			email         = "%[3]s@pingidentity.com"
			population_id = pingone_population.%[2]s.id
		  }

data "pingone_user" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  username = "%[3]s"

  depends_on = [
	pingone_user.%[2]s,
]
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccUserDataSourceConfig_ByEmailFull(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

		resource "pingone_population" "%[2]s" {
			environment_id = data.pingone_environment.general_test.id
		  
			name = "%[3]s"
		  }
	
		resource "pingone_user" "%[2]s" {
			environment_id = data.pingone_environment.general_test.id
		  
			username      = "%[3]s"
			email         = "%[3]s@pingidentity.com"
			population_id = pingone_population.%[2]s.id
		  }

data "pingone_user" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  email = "%[3]s@pingidentity.com"

  depends_on = [
	pingone_user.%[2]s,
]
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccUserDataSourceConfig_ByIDFull(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

	resource "pingone_population" "%[2]s" {
		environment_id = data.pingone_environment.general_test.id
	  
		name = "%[3]s"
	  }

	resource "pingone_user" "%[2]s" {
		environment_id = data.pingone_environment.general_test.id
	  
		username      = "%[3]s"
		email         = "%[3]s@pingidentity.com"
		population_id = pingone_population.%[2]s.id
	  }

data "pingone_user" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  user_id = pingone_user.%[2]s.id
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
