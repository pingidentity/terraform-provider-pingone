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

func testAccCheckSchemaAttributeDestroy(s *terraform.State) error {
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
		if rs.Type != "pingone_schema_attribute" {
			continue
		}

		_, rEnv, err := apiClient.EnvironmentsApi.ReadOneEnvironment(ctx, rs.Primary.Attributes["environment_id"]).Execute()

		if rEnv.StatusCode == 404 {
			continue
		}

		if err != nil {
			return err
		}

		body, r, err := apiClient.SchemasApi.ReadOneAttribute(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["schema_id"], rs.Primary.ID).Execute()

		if r.StatusCode == 404 {
			continue
		}

		if err != nil {
			tflog.Error(ctx, fmt.Sprintf("Error: %v", body))
			return err
		}

		return fmt.Errorf("PingOne Schema Attribute Instance %s still exists", rs.Primary.ID)
	}

	return nil
}

func TestAccSchemaAttribute_NewEnv(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_schema_attribute.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckSchemaAttributeDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSchemaAttributeConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
				),
			},
		},
	})
}

func TestAccSchemaAttribute_FullString(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_schema_attribute.%s", resourceName)

	name := resourceName
	description := "Test description"

	displayName := fmt.Sprintf("Attribute %s", resourceName)
	attrType := "STRING"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckSchemaAttributeDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSchemaAttributeConfig_Full(resourceName, name, attrType, true, true, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttrSet(resourceFullName, "schema_id"),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "display_name", displayName),
					resource.TestCheckResourceAttr(resourceFullName, "description", description),
					resource.TestCheckResourceAttr(resourceFullName, "type", attrType),
					resource.TestCheckResourceAttr(resourceFullName, "unique", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "required", "false"), // Checking the behaviour of the API
					resource.TestCheckResourceAttr(resourceFullName, "multivalued", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "schema_type", "CUSTOM"),
				),
			},
			{
				Config: testAccSchemaAttributeConfig_Full(resourceName, name, attrType, false, true, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttrSet(resourceFullName, "schema_id"),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "display_name", displayName),
					resource.TestCheckResourceAttr(resourceFullName, "description", description),
					resource.TestCheckResourceAttr(resourceFullName, "type", attrType),
					resource.TestCheckResourceAttr(resourceFullName, "unique", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "required", "false"), // Checking the behaviour of the API
					resource.TestCheckResourceAttr(resourceFullName, "multivalued", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "schema_type", "CUSTOM"),
				),
			},
			{
				Config: testAccSchemaAttributeConfig_Full(resourceName, name, attrType, false, false, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttrSet(resourceFullName, "schema_id"),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "display_name", displayName),
					resource.TestCheckResourceAttr(resourceFullName, "description", description),
					resource.TestCheckResourceAttr(resourceFullName, "type", attrType),
					resource.TestCheckResourceAttr(resourceFullName, "unique", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "required", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "multivalued", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "schema_type", "CUSTOM"),
				),
			},
			{
				Config: testAccSchemaAttributeConfig_Full(resourceName, name, attrType, false, false, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttrSet(resourceFullName, "schema_id"),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "display_name", displayName),
					resource.TestCheckResourceAttr(resourceFullName, "description", description),
					resource.TestCheckResourceAttr(resourceFullName, "type", attrType),
					resource.TestCheckResourceAttr(resourceFullName, "unique", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "required", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "multivalued", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "schema_type", "CUSTOM"),
				),
			},
			{
				Config: testAccSchemaAttributeConfig_Full(resourceName, name, attrType, true, false, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttrSet(resourceFullName, "schema_id"),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "display_name", displayName),
					resource.TestCheckResourceAttr(resourceFullName, "description", description),
					resource.TestCheckResourceAttr(resourceFullName, "type", attrType),
					resource.TestCheckResourceAttr(resourceFullName, "unique", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "required", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "multivalued", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "schema_type", "CUSTOM"),
				),
			},
			{
				Config: testAccSchemaAttributeConfig_Full(resourceName, name, attrType, true, true, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttrSet(resourceFullName, "schema_id"),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "display_name", displayName),
					resource.TestCheckResourceAttr(resourceFullName, "description", description),
					resource.TestCheckResourceAttr(resourceFullName, "type", attrType),
					resource.TestCheckResourceAttr(resourceFullName, "unique", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "required", "false"), // Checking the behaviour of the API
					resource.TestCheckResourceAttr(resourceFullName, "multivalued", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "schema_type", "CUSTOM"),
				),
			},
			{
				Config: testAccSchemaAttributeConfig_Full(resourceName, name, attrType, true, true, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttrSet(resourceFullName, "schema_id"),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "display_name", displayName),
					resource.TestCheckResourceAttr(resourceFullName, "description", description),
					resource.TestCheckResourceAttr(resourceFullName, "type", attrType),
					resource.TestCheckResourceAttr(resourceFullName, "unique", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "required", "false"), // Checking the behaviour of the API
					resource.TestCheckResourceAttr(resourceFullName, "multivalued", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "schema_type", "CUSTOM"),
				),
			},
		},
	})
}

func TestAccSchemaAttribute_FullJSON(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_schema_attribute.%s", resourceName)

	name := resourceName
	description := "Test description"

	displayName := fmt.Sprintf("Attribute %s", resourceName)
	attrType := "JSON"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckSchemaAttributeDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccSchemaAttributeConfig_Full(resourceName, name, attrType, true, true, true),
				ExpectError: regexp.MustCompile(`Cannot set attribute unique parameter when the attribute type is not STRING.  Attribute type found: JSON`),
			},
			{
				Config: testAccSchemaAttributeConfig_Full(resourceName, name, attrType, false, true, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttrSet(resourceFullName, "schema_id"),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "display_name", displayName),
					resource.TestCheckResourceAttr(resourceFullName, "description", description),
					resource.TestCheckResourceAttr(resourceFullName, "type", attrType),
					resource.TestCheckResourceAttr(resourceFullName, "unique", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "required", "false"), // Checking the behaviour of the API
					resource.TestCheckResourceAttr(resourceFullName, "multivalued", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "schema_type", "CUSTOM"),
				),
			},
			{
				Config: testAccSchemaAttributeConfig_Full(resourceName, name, attrType, false, false, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttrSet(resourceFullName, "schema_id"),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "display_name", displayName),
					resource.TestCheckResourceAttr(resourceFullName, "description", description),
					resource.TestCheckResourceAttr(resourceFullName, "type", attrType),
					resource.TestCheckResourceAttr(resourceFullName, "unique", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "required", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "multivalued", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "schema_type", "CUSTOM"),
				),
			},
			{
				Config: testAccSchemaAttributeConfig_Full(resourceName, name, attrType, false, false, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttrSet(resourceFullName, "schema_id"),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "display_name", displayName),
					resource.TestCheckResourceAttr(resourceFullName, "description", description),
					resource.TestCheckResourceAttr(resourceFullName, "type", attrType),
					resource.TestCheckResourceAttr(resourceFullName, "unique", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "required", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "multivalued", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "schema_type", "CUSTOM"),
				),
			},
			{
				Config:      testAccSchemaAttributeConfig_Full(resourceName, name, attrType, true, false, false),
				ExpectError: regexp.MustCompile(`Cannot set attribute unique parameter when the attribute type is not STRING.  Attribute type found: JSON`),
			},
			{
				Config: testAccSchemaAttributeConfig_Full(resourceName, name, attrType, false, true, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttrSet(resourceFullName, "schema_id"),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "display_name", displayName),
					resource.TestCheckResourceAttr(resourceFullName, "description", description),
					resource.TestCheckResourceAttr(resourceFullName, "type", attrType),
					resource.TestCheckResourceAttr(resourceFullName, "unique", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "required", "false"), // Checking the behaviour of the API
					resource.TestCheckResourceAttr(resourceFullName, "multivalued", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "schema_type", "CUSTOM"),
				),
			},
			{
				Config: testAccSchemaAttributeConfig_Full(resourceName, name, attrType, false, true, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttrSet(resourceFullName, "schema_id"),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "display_name", displayName),
					resource.TestCheckResourceAttr(resourceFullName, "description", description),
					resource.TestCheckResourceAttr(resourceFullName, "type", attrType),
					resource.TestCheckResourceAttr(resourceFullName, "unique", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "required", "false"), // Checking the behaviour of the API
					resource.TestCheckResourceAttr(resourceFullName, "multivalued", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "schema_type", "CUSTOM"),
				),
			},
		},
	})
}

func TestAccSchemaAttribute_Minimal(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_schema_attribute.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckSchemaAttributeDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSchemaAttributeConfig_Minimal(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttrSet(resourceFullName, "schema_id"),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "display_name", ""),
					resource.TestCheckResourceAttr(resourceFullName, "description", ""),
					resource.TestCheckResourceAttr(resourceFullName, "type", "STRING"),
					resource.TestCheckResourceAttr(resourceFullName, "unique", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "required", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "multivalued", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "schema_type", "CUSTOM"),
				),
			},
		},
	})
}

func testAccSchemaAttributeConfig_NewEnv(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

		data "pingone_schema" "%[3]s" {
			environment_id = "${pingone_environment.%[2]s.id}"

			name = "User"
		}

		resource "pingone_schema_attribute" "%[3]s" {
			environment_id = "${pingone_environment.%[2]s.id}"
			schema_id = "${data.pingone_schema.%[4]s.id}"

			name = "%[4]s"
		}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccSchemaAttributeConfig_Full(resourceName, name, attrType string, unique, required, multivalued bool) string {
	return fmt.Sprintf(`
		%[1]s

		data "pingone_schema" "%[2]s" {
			environment_id = "${data.pingone_environment.general_test.id}"
	
			name = "User"
		}

		resource "pingone_schema_attribute" "%[2]s" {
			environment_id = "${data.pingone_environment.general_test.id}"
			schema_id = "${data.pingone_schema.%[2]s.id}"

			name = "%[3]s"
			display_name = "Attribute %[3]s"
			description = "Test description"

			type = "%[4]s"
			unique = %[5]t
			# required = %[6]t
			multivalued = %[7]t
		}`, acctest.GenericSandboxEnvironment(), resourceName, name, attrType, unique, required, multivalued)
}

func testAccSchemaAttributeConfig_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

		data "pingone_schema" "%[2]s" {
			environment_id = "${data.pingone_environment.general_test.id}"

			name = "User"
		}

		resource "pingone_schema_attribute" "%[2]s" {
			environment_id = "${data.pingone_environment.general_test.id}"
			schema_id = "${data.pingone_schema.%[2]s.id}"

			name = "%[3]s"
		}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
