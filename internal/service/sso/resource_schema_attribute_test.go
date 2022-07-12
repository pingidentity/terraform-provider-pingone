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
	pingone "github.com/patrickcping/pingone-go/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
)

func testAccCheckSchemaAttributeDestroy(s *terraform.State) error {
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

func TestAccSchemaAttribute_FullString(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_schema_attribute.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName
	description := "Test description"

	licenseID := os.Getenv("PINGONE_LICENSE_ID")
	region := os.Getenv("PINGONE_REGION")

	displayName := fmt.Sprintf("Attribute %s", resourceName)
	attrType := "STRING"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckSchemaAttributeDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSchemaAttributeConfig_Full(environmentName, licenseID, region, resourceName, name, displayName, description, attrType, true, true, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, "id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "environment_id"),
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
				Config: testAccSchemaAttributeConfig_Full(environmentName, licenseID, region, resourceName, name, displayName, description, attrType, false, true, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, "id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "environment_id"),
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
				Config: testAccSchemaAttributeConfig_Full(environmentName, licenseID, region, resourceName, name, displayName, description, attrType, false, false, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, "id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "environment_id"),
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
				Config: testAccSchemaAttributeConfig_Full(environmentName, licenseID, region, resourceName, name, displayName, description, attrType, false, false, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, "id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "environment_id"),
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
				Config: testAccSchemaAttributeConfig_Full(environmentName, licenseID, region, resourceName, name, displayName, description, attrType, true, false, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, "id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "environment_id"),
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
				Config: testAccSchemaAttributeConfig_Full(environmentName, licenseID, region, resourceName, name, displayName, description, attrType, true, true, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, "id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "environment_id"),
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
				Config: testAccSchemaAttributeConfig_Full(environmentName, licenseID, region, resourceName, name, displayName, description, attrType, true, true, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, "id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "environment_id"),
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

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName
	description := "Test description"

	licenseID := os.Getenv("PINGONE_LICENSE_ID")
	region := os.Getenv("PINGONE_REGION")

	displayName := fmt.Sprintf("Attribute %s", resourceName)
	attrType := "JSON"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckSchemaAttributeDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccSchemaAttributeConfig_Full(environmentName, licenseID, region, resourceName, name, displayName, description, attrType, true, true, true),
				ExpectError: regexp.MustCompile(`Cannot set attribute unique parameter when the attribute type is not STRING.  Attribute type found: JSON`),
			},
			{
				Config: testAccSchemaAttributeConfig_Full(environmentName, licenseID, region, resourceName, name, displayName, description, attrType, false, true, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, "id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "environment_id"),
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
				Config: testAccSchemaAttributeConfig_Full(environmentName, licenseID, region, resourceName, name, displayName, description, attrType, false, false, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, "id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "environment_id"),
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
				Config: testAccSchemaAttributeConfig_Full(environmentName, licenseID, region, resourceName, name, displayName, description, attrType, false, false, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, "id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "environment_id"),
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
				Config:      testAccSchemaAttributeConfig_Full(environmentName, licenseID, region, resourceName, name, displayName, description, attrType, true, false, false),
				ExpectError: regexp.MustCompile(`Cannot set attribute unique parameter when the attribute type is not STRING.  Attribute type found: JSON`),
			},
			{
				Config: testAccSchemaAttributeConfig_Full(environmentName, licenseID, region, resourceName, name, displayName, description, attrType, false, true, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, "id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "environment_id"),
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
				Config: testAccSchemaAttributeConfig_Full(environmentName, licenseID, region, resourceName, name, displayName, description, attrType, false, true, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, "id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "environment_id"),
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

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")
	region := os.Getenv("PINGONE_REGION")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckSchemaAttributeDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSchemaAttributeConfig_Minimal(environmentName, licenseID, region, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, "id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "environment_id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "schema_id"),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckNoResourceAttr(resourceFullName, "display_name"),
					resource.TestCheckNoResourceAttr(resourceFullName, "description"),
					resource.TestCheckResourceAttr(resourceFullName, "type", "STRING"),
					resource.TestCheckResourceAttr(resourceFullName, "unique", "false"),
					//resource.TestCheckResourceAttr(resourceFullName, "required", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "multivalued", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "schema_type", "CUSTOM"),
				),
			},
		},
	})
}

func testAccSchemaAttributeConfig_Full(environmentName, licenseID, region, resourceName, name, displayName, description, attrType string, unique, required, multivalued bool) string {
	return fmt.Sprintf(`
		resource "pingone_environment" "%[1]s" {
			name = "%[1]s"
			type = "SANDBOX"
			license_id = "%[2]s"
			region = "%[3]s"
			default_population {}
			service {}
		}

		data "pingone_schema" "%[4]s" {
			environment_id = "${pingone_environment.%[1]s.id}"
	
			name = "User"
		}

		resource "pingone_schema_attribute" "%[4]s" {
			environment_id = "${pingone_environment.%[1]s.id}"
			schema_id = "${data.pingone_schema.%[4]s.id}"

			name = "%[5]s"
			display_name = "%[6]s"
			description = "%[7]s"

			type = "%[8]s"
			unique = %[9]t
			# required = %[10]t
			multivalued = %[11]t
		}`, environmentName, licenseID, region, resourceName, name, displayName, description, attrType, unique, required, multivalued)
}

func testAccSchemaAttributeConfig_Minimal(environmentName, licenseID, region, resourceName, name string) string {
	return fmt.Sprintf(`
		resource "pingone_environment" "%[1]s" {
			name = "%[1]s"
			type = "SANDBOX"
			license_id = "%[2]s"
			region = "%[3]s"
			default_population {}
			service {}
		}

		data "pingone_schema" "%[4]s" {
			environment_id = "${pingone_environment.%[1]s.id}"

			name = "User"
		}

		resource "pingone_schema_attribute" "%[4]s" {
			environment_id = "${pingone_environment.%[1]s.id}"
			schema_id = "${data.pingone_schema.%[4]s.id}"

			name = "%[5]s"
		}`, environmentName, licenseID, region, resourceName, name)
}
