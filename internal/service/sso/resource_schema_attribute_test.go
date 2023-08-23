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

func testAccCheckSchemaAttributeDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_schema_attribute" {
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

		body, r, err := apiClient.SchemasApi.ReadOneAttribute(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["schema_id"], rs.Primary.ID).Execute()

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

		return fmt.Errorf("PingOne Schema Attribute Instance %s still exists", rs.Primary.ID)
	}

	return nil
}

func testAccGetSchemaAttributeIDs(resourceName string, environmentID, schemaID, resourceID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource Not found: %s", resourceName)
		}

		*resourceID = rs.Primary.ID
		*schemaID = rs.Primary.Attributes["schema_id"]
		*environmentID = rs.Primary.Attributes["environment_id"]

		return nil
	}
}

func TestAccSchemaAttribute_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_schema_attribute.%s", resourceName)

	name := resourceName

	var resourceID, schemaID, environmentID string

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSchemaAttributeDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccSchemaAttributeConfig_StringMinimal(resourceName, name),
				Check:  testAccGetSchemaAttributeIDs(resourceFullName, &environmentID, &schemaID, &resourceID),
			},
			// Replan after removal preconfig
			{
				PreConfig: func() {
					var ctx = context.Background()
					p1Client, err := acctest.TestClient(ctx)

					if err != nil {
						t.Fatalf("Failed to get API client: %v", err)
					}

					apiClient := p1Client.API.ManagementAPIClient

					if environmentID == "" || schemaID == "" || resourceID == "" {
						t.Fatalf("One of environment ID, schema ID or resource ID cannot be determined. Environment ID: %s, Schema ID: %s, Resource ID: %s", environmentID, schemaID, resourceID)
					}

					_, err = apiClient.SchemasApi.DeleteAttribute(ctx, environmentID, schemaID, resourceID).Execute()
					if err != nil {
						t.Fatalf("Failed to delete schema attribute: %v", err)
					}
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccSchemaAttribute_NewEnv(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_schema_attribute.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSchemaAttributeDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
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

func TestAccSchemaAttribute_String(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_schema_attribute.%s", resourceName)

	name := resourceName
	description := "Test description"

	displayName := fmt.Sprintf("Attribute %s", resourceName)

	fullCheck := resource.TestStep{
		Config: testAccSchemaAttributeConfig_StringFull(resourceName, name, true, true),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "schema_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "schema_name", "User"),
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestCheckResourceAttr(resourceFullName, "display_name", displayName),
			resource.TestCheckResourceAttr(resourceFullName, "description", description),
			resource.TestCheckResourceAttr(resourceFullName, "type", "STRING"),
			resource.TestCheckResourceAttr(resourceFullName, "unique", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "required", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "multivalued", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "schema_type", "CUSTOM"),
		),
	}

	minimalCheck := resource.TestStep{
		Config: testAccSchemaAttributeConfig_StringMinimal(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "schema_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "schema_name", "User"),
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestCheckNoResourceAttr(resourceFullName, "display_name"),
			resource.TestCheckNoResourceAttr(resourceFullName, "description"),
			resource.TestCheckResourceAttr(resourceFullName, "type", "STRING"),
			resource.TestCheckResourceAttr(resourceFullName, "unique", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "required", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "multivalued", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "schema_type", "CUSTOM"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSchemaAttributeDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			fullCheck,
			{
				Config:  testAccSchemaAttributeConfig_StringFull(resourceName, name, true, true),
				Destroy: true,
			},
			// Minimal
			minimalCheck,
			{
				Config:  testAccSchemaAttributeConfig_StringMinimal(resourceName, name),
				Destroy: true,
			},
			// Change
			fullCheck,
			minimalCheck,
			fullCheck,
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["schema_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSchemaAttribute_StringEnumeratedValues(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_schema_attribute.%s", resourceName)

	name := resourceName

	fullCheck := resource.TestStep{
		Config: testAccSchemaAttributeConfig_EnumeratedValues(resourceName, name, "STRING"),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "enumerated_values.#", "6"),
			resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "enumerated_values.*", map[string]string{
				"value":       "value1",
				"archived":    "false",
				"description": "Test description",
			}),
			resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "enumerated_values.*", map[string]string{
				"value":       "value2",
				"description": "Test description",
			}),
			resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "enumerated_values.*", map[string]string{
				"value":       "value3",
				"archived":    "true",
				"description": "Test description",
			}),
			resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "enumerated_values.*", map[string]string{
				"value": "value4",
			}),
			resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "enumerated_values.*", map[string]string{
				"value":    "value5",
				"archived": "true",
			}),
			resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "enumerated_values.*", map[string]string{
				"value":    "value6",
				"archived": "false",
			}),
		),
	}

	minimalCheck := resource.TestStep{
		Config: testAccSchemaAttributeConfig_StringMinimal(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckNoResourceAttr(resourceFullName, "enumerated_values"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSchemaAttributeDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			fullCheck,
			{
				Config:  testAccSchemaAttributeConfig_EnumeratedValues(resourceName, name, "STRING"),
				Destroy: true,
			},
			// Minimal
			minimalCheck,
			{
				Config:  testAccSchemaAttributeConfig_StringMinimal(resourceName, name),
				Destroy: true,
			},
			// Change
			fullCheck,
			minimalCheck,
			fullCheck,
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["schema_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSchemaAttribute_StringRegexValidation(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_schema_attribute.%s", resourceName)

	name := resourceName

	fullCheck := resource.TestStep{
		Config: testAccSchemaAttributeConfig_RegexValidation(resourceName, name, "STRING"),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "regex_validation.pattern", "^[a-zA-Z0-9]*$"),
			resource.TestCheckResourceAttr(resourceFullName, "regex_validation.requirements", "Did you hear about the cow that aced all her tests?  She was outstanding in her field."),
			resource.TestCheckResourceAttr(resourceFullName, "regex_validation.values_pattern_should_match.#", "2"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "regex_validation.values_pattern_should_match.*", "test123"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "regex_validation.values_pattern_should_match.*", "test456"),
			resource.TestCheckResourceAttr(resourceFullName, "regex_validation.values_pattern_should_not_match.#", "2"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "regex_validation.values_pattern_should_not_match.*", "test123!"),
			resource.TestCheckTypeSetElemAttr(resourceFullName, "regex_validation.values_pattern_should_not_match.*", "test456!"),
		),
	}

	minimalCheck := resource.TestStep{
		Config: testAccSchemaAttributeConfig_StringMinimal(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckNoResourceAttr(resourceFullName, "regex_validation"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSchemaAttributeDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			fullCheck,
			{
				Config:  testAccSchemaAttributeConfig_RegexValidation(resourceName, name, "STRING"),
				Destroy: true,
			},
			// Minimal
			minimalCheck,
			{
				Config:  testAccSchemaAttributeConfig_StringMinimal(resourceName, name),
				Destroy: true,
			},
			// Change
			fullCheck,
			minimalCheck,
			fullCheck,
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["schema_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSchemaAttribute_StringParameterCombinations(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_schema_attribute.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSchemaAttributeDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Combos
			{
				Config: testAccSchemaAttributeConfig_StringFull(resourceName, name, true, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "unique", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "required", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "multivalued", "true"),
				),
			},
			{
				Config: testAccSchemaAttributeConfig_StringFull(resourceName, name, false, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "unique", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "required", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "multivalued", "true"),
				),
			},
			{
				Config: testAccSchemaAttributeConfig_StringFull(resourceName, name, false, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "unique", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "required", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "multivalued", "false"),
				),
			},
			{
				Config: testAccSchemaAttributeConfig_StringFull(resourceName, name, true, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "unique", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "required", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "multivalued", "false"),
				),
			},
		},
	})
}

func TestAccSchemaAttribute_JSON(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_schema_attribute.%s", resourceName)

	name := resourceName
	description := "Test description"

	displayName := fmt.Sprintf("Attribute %s", resourceName)

	fullCheck := resource.TestStep{
		Config: testAccSchemaAttributeConfig_JSONFull(resourceName, name, false, true),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "schema_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "schema_name", "User"),
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestCheckResourceAttr(resourceFullName, "display_name", displayName),
			resource.TestCheckResourceAttr(resourceFullName, "description", description),
			resource.TestCheckResourceAttr(resourceFullName, "type", "JSON"),
			resource.TestCheckResourceAttr(resourceFullName, "unique", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "required", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "multivalued", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "schema_type", "CUSTOM"),
		),
	}

	minimalCheck := resource.TestStep{
		Config: testAccSchemaAttributeConfig_JSONMinimal(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "schema_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "schema_name", "User"),
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestCheckNoResourceAttr(resourceFullName, "display_name"),
			resource.TestCheckNoResourceAttr(resourceFullName, "description"),
			resource.TestCheckResourceAttr(resourceFullName, "type", "JSON"),
			resource.TestCheckResourceAttr(resourceFullName, "unique", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "required", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "multivalued", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "schema_type", "CUSTOM"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSchemaAttributeDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			fullCheck,
			{
				Config:  testAccSchemaAttributeConfig_JSONFull(resourceName, name, false, true),
				Destroy: true,
			},
			// Minimal
			minimalCheck,
			{
				Config:  testAccSchemaAttributeConfig_JSONMinimal(resourceName, name),
				Destroy: true,
			},
			// Change
			fullCheck,
			minimalCheck,
			fullCheck,
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["schema_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSchemaAttribute_JSONInvalidAttrs(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSchemaAttributeDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccSchemaAttributeConfig_RegexValidation(resourceName, name, "JSON"),
				ExpectError: regexp.MustCompile(`Invalid argument combination`),
			},
			{
				Config:      testAccSchemaAttributeConfig_EnumeratedValues(resourceName, name, "JSON"),
				ExpectError: regexp.MustCompile(`Invalid argument combination`),
			},
		},
	})
}

func TestAccSchemaAttribute_JSONParameterCombinations(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_schema_attribute.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSchemaAttributeDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccSchemaAttributeConfig_JSONFull(resourceName, name, true, true),
				ExpectError: regexp.MustCompile(`Invalid Attribute Combination`),
			},
			{
				Config:      testAccSchemaAttributeConfig_JSONFull(resourceName, name, true, false),
				ExpectError: regexp.MustCompile(`Invalid Attribute Combination`),
			},
			{
				Config: testAccSchemaAttributeConfig_JSONFull(resourceName, name, false, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "unique", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "required", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "multivalued", "true"),
				),
			},
			{
				Config: testAccSchemaAttributeConfig_JSONFull(resourceName, name, false, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "unique", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "required", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "multivalued", "false"),
				),
			},
		},
	})
}

func TestAccSchemaAttribute_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_schema_attribute.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSchemaAttributeDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccSchemaAttributeConfig_StringMinimal(resourceName, name),
			},
			// Errors
			{
				ResourceName: resourceFullName,
				ImportState:  true,
				ExpectError:  regexp.MustCompile(`Unexpected Import Identifier`),
			},
			{
				ResourceName:  resourceFullName,
				ImportStateId: "/",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`Unexpected Import Identifier`),
			},
			{
				ResourceName:  resourceFullName,
				ImportStateId: "badformat/badformat/badformat",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`Unexpected Import Identifier`),
			},
		},
	})
}

func testAccSchemaAttributeConfig_NewEnv(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_schema_attribute" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  name = "%[4]s"
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccSchemaAttributeConfig_StringFull(resourceName, name string, unique, multivalued bool) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_schema_attribute" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  schema_name    = "User"

  name         = "%[3]s"
  display_name = "Attribute %[3]s"
  description  = "Test description"

  type        = "STRING"
  unique      = %[4]t
  multivalued = %[5]t
}`, acctest.GenericSandboxEnvironment(), resourceName, name, unique, multivalued)
}

func testAccSchemaAttributeConfig_StringMinimal(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_schema_attribute" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccSchemaAttributeConfig_JSONFull(resourceName, name string, unique, multivalued bool) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_schema_attribute" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  schema_name    = "User"

  name         = "%[3]s"
  display_name = "Attribute %[3]s"
  description  = "Test description"

  type        = "JSON"
  unique      = %[4]t
  multivalued = %[5]t
}`, acctest.GenericSandboxEnvironment(), resourceName, name, unique, multivalued)
}

func testAccSchemaAttributeConfig_JSONMinimal(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_schema_attribute" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
  type = "JSON"
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccSchemaAttributeConfig_EnumeratedValues(resourceName, name, attrType string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_schema_attribute" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
  type = "%[4]s"

  enumerated_values = [
    {
      value       = "value1"
      archived    = "false"
      description = "Test description"
    },
    {
      value       = "value2"
      description = "Test description"
    },
    {
      value       = "value3"
      archived    = "true"
      description = "Test description"
    },
    {
      value = "value4"
    },
    {
      value    = "value5"
      archived = "true"
    },
    {
      value    = "value6"
      archived = "false"
    }
  ]
}`, acctest.GenericSandboxEnvironment(), resourceName, name, attrType)
}

func testAccSchemaAttributeConfig_RegexValidation(resourceName, name, attrType string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_schema_attribute" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
  type = "%[4]s"

  regex_validation = {
    pattern      = "^[a-zA-Z0-9]*$",
    requirements = "Did you hear about the cow that aced all her tests?  She was outstanding in her field."

    values_pattern_should_match = [
      "test123",
      "test456"
    ]

    values_pattern_should_not_match = [
      "test123!",
      "test456!"
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name, attrType)
}
