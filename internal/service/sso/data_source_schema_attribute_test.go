// Copyright © 2025 Ping Identity Corporation

package sso_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccSchemaAttributeDataSource_ByNameFull(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_schema_attribute.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSchemaAttributeDataSourceConfig_ByNameFull(resourceName, resourceName, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "schema_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "attribute_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "name", resourceName),
					resource.TestCheckResourceAttr(dataSourceFullName, "display_name", resourceName),
					resource.TestCheckResourceAttr(dataSourceFullName, "description", resourceName),
					resource.TestCheckResourceAttr(dataSourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "type", "STRING"),
					resource.TestCheckResourceAttr(dataSourceFullName, "unique", "false"),
					resource.TestCheckResourceAttr(dataSourceFullName, "multivalued", "false"),
					resource.TestCheckResourceAttr(dataSourceFullName, "required", "false"),
					resource.TestCheckResourceAttr(dataSourceFullName, "schema_type", "CUSTOM"),
				),
			},
			{
				Config: testAccSchemaAttributeDataSourceConfig_ByNameFull(resourceName, resourceName, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "name", resourceName),
				),
			},
		},
	})
}

func TestAccSchemaAttributeDataSource_ByIDFull(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_schema_attribute.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSchemaAttributeDataSourceConfig_ByIDFull(resourceName, resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "schema_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "attribute_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "name", resourceName),
					resource.TestCheckResourceAttr(dataSourceFullName, "display_name", resourceName),
					resource.TestCheckResourceAttr(dataSourceFullName, "description", resourceName),
					resource.TestCheckResourceAttr(dataSourceFullName, "enabled", "true"),
					resource.TestCheckResourceAttr(dataSourceFullName, "type", "STRING"),
					resource.TestCheckResourceAttr(dataSourceFullName, "unique", "false"),
					resource.TestCheckResourceAttr(dataSourceFullName, "multivalued", "false"),
					resource.TestCheckResourceAttr(dataSourceFullName, "required", "false"),
					resource.TestCheckResourceAttr(dataSourceFullName, "schema_type", "CUSTOM"),
				),
			},
		},
	})
}

func TestAccSchemaAttributeDataSource_BuiltIn(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullNameStandardString := fmt.Sprintf("data.pingone_schema_attribute.%s_standard", resourceName)
	dataSourceFullNameStandardComplex := fmt.Sprintf("data.pingone_schema_attribute.%s_standard_complex", resourceName)
	dataSourceFullNameCore := fmt.Sprintf("data.pingone_schema_attribute.%s_core", resourceName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSchemaAttributeDataSourceConfig_StandardAndCore(resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullNameStandardString, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullNameStandardString, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullNameStandardString, "schema_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullNameStandardString, "attribute_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullNameStandardString, "name", "email"),
					resource.TestCheckResourceAttr(dataSourceFullNameStandardString, "type", "STRING"),
					resource.TestCheckResourceAttr(dataSourceFullNameStandardString, "multivalued", "false"),
					resource.TestCheckResourceAttr(dataSourceFullNameStandardString, "schema_type", "STANDARD"),

					resource.TestMatchResourceAttr(dataSourceFullNameStandardComplex, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullNameStandardComplex, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullNameStandardComplex, "schema_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullNameStandardComplex, "attribute_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullNameStandardComplex, "name", "address"),
					resource.TestCheckResourceAttr(dataSourceFullNameStandardComplex, "type", "COMPLEX"),
					resource.TestCheckResourceAttr(dataSourceFullNameStandardComplex, "schema_type", "STANDARD"),
					testCheckSubAttributesAtLeastOne(dataSourceFullNameStandardComplex),

					resource.TestMatchResourceAttr(dataSourceFullNameCore, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullNameCore, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullNameCore, "schema_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullNameCore, "attribute_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullNameCore, "name", "account"),
					resource.TestCheckResourceAttr(dataSourceFullNameCore, "type", "COMPLEX"),
					resource.TestCheckResourceAttr(dataSourceFullNameCore, "multivalued", "false"),
					resource.TestCheckResourceAttr(dataSourceFullNameCore, "schema_type", "CORE"),
				),
			},
		},
	})
}

func TestAccSchemaAttributeDataSource_NotFound(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccSchemaAttributeDataSourceConfig_NotFoundByName(resourceName),
				ExpectError: regexp.MustCompile("Cannot find schema attribute from name"),
			},
			{
				Config:      testAccSchemaAttributeDataSourceConfig_NotFoundByID(resourceName),
				ExpectError: regexp.MustCompile("Error when calling `ReadOneAttribute`: The request could not be completed. The requested resource was not found."),
			},
		},
	})
}

func testAccSchemaAttributeDataSourceConfig_ByNameFull(resourceName, name string, insensitivityCheck bool) string {

	// If insensitivityCheck is true, alter the case of the name
	nameComparator := name
	if insensitivityCheck {
		nameComparator = acctest.AlterStringCasing(nameComparator)
	}

	return fmt.Sprintf(`
		%[1]s

data "pingone_schema" "%[2]s_schema" {
  environment_id = data.pingone_environment.general_test.id
  name           = "User"
}

resource "pingone_schema_attribute" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  description  = "%[3]s"
  display_name = "%[3]s"
  enabled      = true
}

data "pingone_schema_attribute" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  schema_id      = data.pingone_schema.%[2]s_schema.id

  name = "%[4]s"

  depends_on = [
    pingone_schema_attribute.%[2]s
  ]
}`, acctest.GenericSandboxEnvironment(), resourceName, name, nameComparator)
}

func testAccSchemaAttributeDataSourceConfig_ByIDFull(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_schema" "%[2]s_schema" {
  environment_id = data.pingone_environment.general_test.id
  name           = "User"
}

resource "pingone_schema_attribute" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  description  = "%[3]s"
  display_name = "%[3]s"
  enabled      = true
}

data "pingone_schema_attribute" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  schema_id      = data.pingone_schema.%[2]s_schema.id

  attribute_id = pingone_schema_attribute.%[2]s.id
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccSchemaAttributeDataSourceConfig_NotFoundByName(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_schema" "%[2]s_schema" {
  environment_id = data.pingone_environment.general_test.id
  name           = "User"
}

data "pingone_schema_attribute" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  schema_id      = data.pingone_schema.%[2]s_schema.id

  name = "doesnotexist"
}`, acctest.GenericSandboxEnvironment(), resourceName)
}

func testAccSchemaAttributeDataSourceConfig_NotFoundByID(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_schema" "%[2]s_schema" {
  environment_id = data.pingone_environment.general_test.id
  name           = "User"
}

data "pingone_schema_attribute" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  schema_id      = data.pingone_schema.%[2]s_schema.id

  attribute_id = "9c052a8a-14be-44e4-8f07-2662569994ce" // dummy ID
}`, acctest.GenericSandboxEnvironment(), resourceName)
}

func testAccSchemaAttributeDataSourceConfig_StandardAndCore(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_schema" "%[2]s_schema" {
  environment_id = data.pingone_environment.general_test.id
  name           = "User"
}

data "pingone_schema_attribute" "%[2]s_standard" {
  environment_id = data.pingone_environment.general_test.id
  schema_id      = data.pingone_schema.%[2]s_schema.id

  name = "email"
}

data "pingone_schema_attribute" "%[2]s_standard_complex" {
  environment_id = data.pingone_environment.general_test.id
  schema_id      = data.pingone_schema.%[2]s_schema.id

  name = "address"
}

data "pingone_schema_attribute" "%[2]s_core" {
  environment_id = data.pingone_environment.general_test.id
  schema_id      = data.pingone_schema.%[2]s_schema.id

  name = "account"
}`, acctest.GenericSandboxEnvironment(), resourceName)
}
