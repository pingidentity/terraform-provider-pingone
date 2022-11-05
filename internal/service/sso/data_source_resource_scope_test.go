package sso_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
)

func TestAccResourceScopeDataSource_ByNameFull(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_resource_scope.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckResourceScopeDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceScopeDataSourceConfig_ByNameFull(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "resource_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "name", resourceFullName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "description", resourceFullName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "schema_attributes", resourceFullName, "schema_attributes"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "mapped_claims", resourceFullName, "mapped_claims"),
				),
			},
		},
	})
}

func TestAccResourceScopeDataSource_ByNameSystem(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_resource_scope.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckResourceScopeDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceScopeDataSourceConfig_ByNameSystem(resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "resource_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "resource_scope_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(dataSourceFullName, "name", "email"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "description"),
					resource.TestCheckResourceAttr(dataSourceFullName, "schema_attributes.#", "0"),
					resource.TestCheckResourceAttr(dataSourceFullName, "mapped_claims.#", "0"),
				),
			},
		},
	})
}

func TestAccResourceScopeDataSource_ByIDFull(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_resource_scope.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckResourceScopeDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceScopeDataSourceConfig_ByIDFull(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "resource_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "name", resourceFullName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "description", resourceFullName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "schema_attributes", resourceFullName, "schema_attributes"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "mapped_claims", resourceFullName, "mapped_claims"),
				),
			},
		},
	})
}

func TestAccResourceScopeDataSource_ByIDSchemaAttributes(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_resource_scope.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckResourceScopeDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceScopeDataSourceConfig_ByIDSchemaAttributes(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceFullName, "schema_attributes", resourceFullName, "schema_attributes"),
					resource.TestCheckResourceAttr(dataSourceFullName, "mapped_claims.#", "0"),
				),
			},
		},
	})
}

func TestAccResourceScopeDataSource_ByIDMappedClaims(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_resource_scope.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckResourceScopeDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceScopeDataSourceConfig_ByIDMappedClaims(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceFullName, "schema_attributes.#", "0"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "mapped_claims", resourceFullName, "mapped_claims"),
				),
			},
		},
	})
}

func testAccResourceScopeDataSourceConfig_ByNameFull(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_resource" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_resource_scope" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  resource_id    = pingone_resource.%[2]s.id

  name = "%[3]s"
}

data "pingone_resource_scope" "%[3]s" {
  environment_id = data.pingone_environment.general_test.id
  resource_id    = pingone_resource.%[2]s.id

  name = "%[3]s"

  depends_on = [
    pingone_resource_scope.%[2]s
  ]
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccResourceScopeDataSourceConfig_ByIDFull(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_resource" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_resource_scope" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  resource_id    = pingone_resource.%[2]s.id

  name = "%[3]s"
}

data "pingone_resource_scope" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  resource_id    = pingone_resource.%[2]s.id

  resource_scope_id = pingone_resource_scope.%[2]s.id
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccResourceScopeDataSourceConfig_ByNameSystem(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_resource" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "openid"
}

data "pingone_resource_scope" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  resource_id    = data.pingone_resource.%[2]s.id

  name = "email"
}`, acctest.GenericSandboxEnvironment(), resourceName)
}

func testAccResourceScopeDataSourceConfig_ByIDSchemaAttributes(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

	resource "pingone_resource" "%[2]s" {
		environment_id = data.pingone_environment.general_test.id
	  
		name = "%[3]s"
	  }
	  
	  resource "pingone_resource_scope" "%[2]s" {
		environment_id = data.pingone_environment.general_test.id
		resource_id    = pingone_resource.%[2]s.id
	  
		name        = "%[3]s"
		description = "My resource scope"
	  
		schema_attributes = [
		  "name.given",
		  "customAttr",
		  "name.family",
		]
	  }
	  
	data "pingone_resource_scope" "%[2]s" {
		environment_id = data.pingone_environment.general_test.id
		resource_id    = pingone_resource.%[2]s.id
		resource_scope_id = pingone_resource_scope.%[2]s.id
		
	}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccResourceScopeDataSourceConfig_ByIDMappedClaims(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

	resource "pingone_resource" "%[2]s" {
		environment_id = data.pingone_environment.general_test.id
	  
		name = "%[3]s"
	  }
	  
	  resource "pingone_resource_attribute" "%[2]s-1" {
		  environment_id = data.pingone_environment.general_test.id
		  resource_id    = pingone_resource.%[2]s.id
		
		  name        = "%[3]s-1"
		  value = "$${user.name.given}"
		}
	  
		resource "pingone_resource_attribute" "%[2]s-2" {
		  environment_id = data.pingone_environment.general_test.id
		  resource_id    = pingone_resource.%[2]s.id
		
		  name        = "%[3]s-2"
		  value = "$${user.name.family}"
		}
	  
		resource "pingone_resource_attribute" "%[2]s-3" {
		  environment_id = data.pingone_environment.general_test.id
		  resource_id    = pingone_resource.%[2]s.id
		
		  name        = "%[3]s-3"
		  value = "$${user.email}"
		}
	  
	  resource "pingone_resource_scope" "%[2]s" {
		environment_id = data.pingone_environment.general_test.id
		resource_id    = pingone_resource.%[2]s.id
	  
		name        = "%[3]s"
		description = "My resource scope"
	  
		mapped_claims = [
		  pingone_resource_attribute.%[2]s-2.id
		  pingone_resource_attribute.%[2]s-3.id
		  pingone_resource_attribute.%[2]s-1.id
		]
	  }
	  
	 data  "pingone_resource_scope" "%[2]s" {
		environment_id = data.pingone_environment.general_test.id
		resource_id    = pingone_resource.%[2]s.id
		resource_scope_id = pingone_resource_scope.%[2]s.id
		
	}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
