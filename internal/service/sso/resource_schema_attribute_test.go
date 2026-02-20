// Copyright © 2026 Ping Identity Corporation

package sso_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/legacysdk"
	acctestlegacysdk "github.com/pingidentity/terraform-provider-pingone/internal/acctest/legacysdk"
	baselegacysdk "github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/base/legacysdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/sso"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccSchemaAttribute_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_schema_attribute.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	var schemaAttributeID, schemaID, environmentID string

	var p1Client *client.Client
	var ctx = context.Background()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoBeta(t)
			p1Client = acctestlegacysdk.PreCheckTestClient(ctx, t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.SchemaAttribute_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccSchemaAttributeConfig_StringMinimal(resourceName, name),
				Check:  sso.SchemaAttribute_GetIDs(resourceFullName, &environmentID, &schemaID, &schemaAttributeID),
			},
			{
				PreConfig: func() {
					sso.SchemaAttribute_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID, schemaID, schemaAttributeID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the environment
			{
				Config: testAccSchemaAttributeConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check:  sso.SchemaAttribute_GetIDs(resourceFullName, &environmentID, &schemaID, &schemaAttributeID),
			},
			{
				PreConfig: func() {
					baselegacysdk.Environment_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID)
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
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.SchemaAttribute_CheckDestroy,
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
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestCheckNoResourceAttr(resourceFullName, "display_name"),
			resource.TestCheckNoResourceAttr(resourceFullName, "description"),
			resource.TestCheckResourceAttr(resourceFullName, "type", "STRING"),
			resource.TestCheckResourceAttr(resourceFullName, "unique", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "required", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "multivalued", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "schema_type", "CUSTOM"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.SchemaAttribute_CheckDestroy,
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
							return "", fmt.Errorf("resource not found: %s", resourceFullName)
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
		Config: testAccSchemaAttributeConfig_EnumeratedValues1(resourceName, name, "STRING"),
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
		Config: testAccSchemaAttributeConfig_EnumeratedValues2(resourceName, name, "STRING"),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "enumerated_values.#", "7"),
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
				"archived": "true",
			}),
			resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "enumerated_values.*", map[string]string{
				"value":    "value7",
				"archived": "false",
			}),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.SchemaAttribute_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Change
			fullCheck,
			minimalCheck,
			{
				Config:      testAccSchemaAttributeConfig_EnumeratedValues1(resourceName, name, "STRING"),
				ExpectError: regexp.MustCompile(`Data Loss Protection`),
			},
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("resource not found: %s", resourceFullName)
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
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.SchemaAttribute_CheckDestroy,
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
							return "", fmt.Errorf("resource not found: %s", resourceFullName)
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
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.SchemaAttribute_CheckDestroy,
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
				Config:  testAccSchemaAttributeConfig_StringFull(resourceName, name, true, true),
				Destroy: true,
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
				Config:  testAccSchemaAttributeConfig_StringFull(resourceName, name, false, true),
				Destroy: true,
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
				Config:  testAccSchemaAttributeConfig_StringFull(resourceName, name, false, false),
				Destroy: true,
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
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestCheckNoResourceAttr(resourceFullName, "display_name"),
			resource.TestCheckNoResourceAttr(resourceFullName, "description"),
			resource.TestCheckResourceAttr(resourceFullName, "type", "JSON"),
			resource.TestCheckResourceAttr(resourceFullName, "unique", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "required", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "multivalued", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "schema_type", "CUSTOM"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.SchemaAttribute_CheckDestroy,
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
							return "", fmt.Errorf("resource not found: %s", resourceFullName)
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
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.SchemaAttribute_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccSchemaAttributeConfig_RegexValidation(resourceName, name, "JSON"),
				ExpectError: regexp.MustCompile(`Invalid argument combination`),
			},
			{
				Config:      testAccSchemaAttributeConfig_EnumeratedValues1(resourceName, name, "JSON"),
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
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.SchemaAttribute_CheckDestroy,
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
				Config:  testAccSchemaAttributeConfig_JSONFull(resourceName, name, false, true),
				Destroy: true,
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
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.SchemaAttribute_CheckDestroy,
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

func TestAccSchemaAttribute_DLP(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.SchemaAttribute_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Attribute type
			{
				Config: testAccSchemaAttributeConfig_StringFull(resourceName, name, false, false),
			},
			{
				Config:      testAccSchemaAttributeConfig_JSONFull(resourceName, name, false, false),
				ExpectError: regexp.MustCompile(`Data Loss Protection`),
			},
			{
				Config:  testAccSchemaAttributeConfig_StringFull(resourceName, name, false, false),
				Destroy: true,
			},
			// Multivalued
			{
				Config: testAccSchemaAttributeConfig_StringFull(resourceName, name, true, false),
			},
			{
				Config:      testAccSchemaAttributeConfig_StringFull(resourceName, name, true, true),
				ExpectError: regexp.MustCompile(`Data Loss Protection`),
			},
			{
				Config:  testAccSchemaAttributeConfig_StringFull(resourceName, name, true, false),
				Destroy: true,
			},
			// Enumerated values - enable on existing attribute
			{
				Config: testAccSchemaAttributeConfig_StringFull(resourceName, name, false, false),
			},
			{
				Config:      testAccSchemaAttributeConfig_EnumeratedValues1(resourceName, name, "STRING"),
				ExpectError: regexp.MustCompile(`Data Loss Protection`),
			},
			{
				Config:  testAccSchemaAttributeConfig_StringFull(resourceName, name, false, false),
				Destroy: true,
			},
			// Enumerated values - delete existing value
			{
				Config: testAccSchemaAttributeConfig_EnumeratedValues2(resourceName, name, "STRING"),
			},
			{
				Config:      testAccSchemaAttributeConfig_EnumeratedValues1(resourceName, name, "STRING"),
				ExpectError: regexp.MustCompile(`Data Loss Protection`),
			},
			{
				Config:  testAccSchemaAttributeConfig_EnumeratedValues2(resourceName, name, "STRING"),
				Destroy: true,
			},
		},
	})
}

func TestAccSchemaAttribute_Standard(t *testing.T) {
	t.Parallel()

	environmentName := acctest.ResourceNameGenEnvironment()
	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	fullResourceName := "pingone_schema_attribute.email"
	dataSourceName := "data.pingone_schema_attribute.email"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.SchemaAttribute_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				// Step 1: Resolve IDs for the built-in STANDARD attribute
				Config: testAccSchemaAttributeConfig_StandardDataSource(environmentName, licenseID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "schema_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "environment_id"),
				),
			},
			{
				// Step 2: Import the STANDARD attribute into Terraform state
				Config:       testAccSchemaAttributeConfig_StandardResource(environmentName, licenseID, true),
				ResourceName: fullResourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					ds, ok := s.RootModule().Resources[dataSourceName]
					if !ok {
						return "", fmt.Errorf("Not found: %s", dataSourceName)
					}
					return fmt.Sprintf("%s/%s/%s", ds.Primary.Attributes["environment_id"], ds.Primary.Attributes["schema_id"], ds.Primary.Attributes["id"]), nil
				},
				ImportStateVerify:  false,
				ImportStatePersist: true,
			},
			{
				// Step 3: Manage imported STANDARD attribute
				Config: testAccSchemaAttributeConfig_StandardResource(environmentName, licenseID, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fullResourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(fullResourceName, "schema_type", "STANDARD"),
				),
			},
			{
				// Step 4: Re-import and verify once resource is already managed
				Config:       testAccSchemaAttributeConfig_StandardResource(environmentName, licenseID, false),
				ResourceName: fullResourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[fullResourceName]
					if !ok {
						return "", fmt.Errorf("Not found: %s", fullResourceName)
					}
					return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["schema_id"], rs.Primary.ID), nil
				},
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSchemaAttribute_DeleteBehavior(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	environmentName := acctest.ResourceNameGenEnvironment()
	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	customResourceName := "pingone_schema_attribute.custom"
	standardDataSourceName := "data.pingone_schema_attribute.email"
	standardResourceName := "pingone_schema_attribute.email"

	customAttributeName := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             sso.SchemaAttribute_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// 1) CUSTOM attribute gets deleted completely
			{
				Config: testAccSchemaAttributeScenarioConfig_CustomResource(environmentName, licenseID, customAttributeName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(customResourceName, "schema_type", "CUSTOM"),
				),
			},
			{
				Config: testAccSchemaAttributeScenarioConfig_BuiltInDataSources(environmentName, licenseID),
			},
			{
				Config:      testAccSchemaAttributeScenarioConfig_CustomNotFound(environmentName, licenseID, customAttributeName),
				ExpectError: regexp.MustCompile("Cannot find schema attribute from name"),
			},

			// 2) Imported STANDARD attribute is reset to defaults in API when deleted
			{
				Config:       testAccSchemaAttributeScenarioConfig_StandardImport(environmentName, licenseID, true, false),
				ResourceName: standardResourceName,
				ImportState:  true,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						ds, ok := s.RootModule().Resources[standardDataSourceName]
						if !ok {
							return "", fmt.Errorf("resource not found: %s", standardDataSourceName)
						}

						return fmt.Sprintf("%s/%s/%s", ds.Primary.Attributes["environment_id"], ds.Primary.Attributes["schema_id"], ds.Primary.Attributes["id"]), nil
					}
				}(),
				ImportStateVerify:  false,
				ImportStatePersist: true,
			},
			{
				Config: testAccSchemaAttributeScenarioConfig_StandardResource(environmentName, licenseID, false, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(standardResourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(standardResourceName, "unique", "true"),
					resource.TestCheckResourceAttr(standardResourceName, "regex_validation.pattern", "^[^@]+@[^@]+$"),
					resource.TestCheckResourceAttr(standardResourceName, "regex_validation.requirements", "Must look like an email"),
				),
			},
			{
				Config: testAccSchemaAttributeScenarioConfig_BuiltInDataSources(environmentName, licenseID),
			},
			{
				Config: testAccSchemaAttributeScenarioConfig_BuiltInDataSources(environmentName, licenseID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(standardDataSourceName, "schema_type", "STANDARD"),
					resource.TestCheckResourceAttr(standardDataSourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(standardDataSourceName, "unique", "false"),
					resource.TestCheckNoResourceAttr(standardDataSourceName, "regex_validation"),
				),
			},
		},
	})
}

func TestAccSchemaAttribute_ImportCoreUnsupported(t *testing.T) {
	t.Parallel()

	environmentName := acctest.ResourceNameGenEnvironment()
	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	coreDataSourceName := "data.pingone_schema_attribute.account"
	coreResourceName := "pingone_schema_attribute.account"

	var coreEnvironmentID, coreSchemaID, coreAttributeID string

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSchemaAttributeScenarioConfig_BuiltInDataSources(environmentName, licenseID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(coreDataSourceName, "schema_type", "CORE"),
					func(s *terraform.State) error {
						ds, ok := s.RootModule().Resources[coreDataSourceName]
						if !ok {
							return fmt.Errorf("resource not found: %s", coreDataSourceName)
						}

						coreEnvironmentID = ds.Primary.Attributes["environment_id"]
						coreSchemaID = ds.Primary.Attributes["schema_id"]
						coreAttributeID = ds.Primary.Attributes["id"]

						return nil
					},
				),
			},
			{
				Config:       testAccSchemaAttributeScenarioConfig_CoreImport(environmentName, licenseID),
				ResourceName: coreResourceName,
				ImportState:  true,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						if coreEnvironmentID == "" || coreSchemaID == "" || coreAttributeID == "" {
							return "", fmt.Errorf("core schema attribute IDs were not captured")
						}

						return fmt.Sprintf("%s/%s/%s", coreEnvironmentID, coreSchemaID, coreAttributeID), nil
					}
				}(),
				ExpectError: regexp.MustCompile("Invalid import for CORE schema attribute"),
			},
		},
	})
}

func testAccSchemaAttributeScenarioConfig_CustomResource(environmentName, licenseID, attributeName string) string {
	return fmt.Sprintf(`
		%s

		resource "pingone_schema_attribute" "custom" {
			environment_id = pingone_environment.%s.id

			name = "%s"
			type = "STRING"
		}
	`, legacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, attributeName)
}

func testAccSchemaAttributeScenarioConfig_BuiltInDataSources(environmentName, licenseID string) string {
	return fmt.Sprintf(`
		%s

		data "pingone_schema" "user" {
			environment_id = pingone_environment.%s.id
			name           = "User"
		}

		data "pingone_schema_attribute" "email" {
			environment_id = pingone_environment.%s.id
			schema_id      = data.pingone_schema.user.id
			name           = "email"
		}

		data "pingone_schema_attribute" "account" {
			environment_id = pingone_environment.%s.id
			schema_id      = data.pingone_schema.user.id
			name           = "account"
		}
	`, legacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, environmentName, environmentName)
}

func testAccSchemaAttributeScenarioConfig_StandardResource(environmentName, licenseID string, enabled, unique bool) string {
	return fmt.Sprintf(`
		%s

		data "pingone_schema" "user" {
			environment_id = pingone_environment.%s.id
			name           = "User"
		}

		resource "pingone_schema_attribute" "email" {
			environment_id = pingone_environment.%s.id

			name    = "email"
			enabled = %t
			unique  = %t

			regex_validation = {
				pattern      = "^[^@]+@[^@]+$"
				requirements = "Must look like an email"
			}
		}
	`, legacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, environmentName, enabled, unique)
}

func testAccSchemaAttributeScenarioConfig_StandardImport(environmentName, licenseID string, enabled, unique bool) string {
	return fmt.Sprintf(`
		%s

		data "pingone_schema" "user" {
			environment_id = pingone_environment.%s.id
			name           = "User"
		}

		data "pingone_schema_attribute" "email" {
			environment_id = pingone_environment.%s.id
			schema_id      = data.pingone_schema.user.id
			name           = "email"
		}

		resource "pingone_schema_attribute" "email" {
			environment_id = pingone_environment.%s.id

			name    = "email"
			enabled = %t
			unique  = %t

			regex_validation = {
				pattern      = "^[^@]+@[^@]+$"
				requirements = "Must look like an email"
			}
		}
	`, legacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, environmentName, environmentName, enabled, unique)
}

func testAccSchemaAttributeScenarioConfig_CoreImport(environmentName, licenseID string) string {
	return fmt.Sprintf(`
		%s

		data "pingone_schema" "user" {
			environment_id = pingone_environment.%s.id
			name           = "User"
		}

		data "pingone_schema_attribute" "account" {
			environment_id = pingone_environment.%s.id
			schema_id      = data.pingone_schema.user.id
			name           = "account"
		}

		resource "pingone_schema_attribute" "account" {
			environment_id = pingone_environment.%s.id

			name = "account"
			type = "COMPLEX"
		}
	`, legacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, environmentName, environmentName)
}

func testAccSchemaAttributeScenarioConfig_CustomNotFound(environmentName, licenseID, attributeName string) string {
	return fmt.Sprintf(`
		%s

		data "pingone_schema" "user" {
			environment_id = pingone_environment.%s.id
			name           = "User"
		}

		data "pingone_schema_attribute" "custom" {
			environment_id = pingone_environment.%s.id
			schema_id      = data.pingone_schema.user.id
			name           = "%s"
		}
	`, legacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, environmentName, attributeName)
}

func testAccSchemaAttributeConfig_NewEnv(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_schema_attribute" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  name = "%[4]s"
}`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccSchemaAttributeConfig_StringFull(resourceName, name string, unique, multivalued bool) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_schema_attribute" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

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

  name        = "%[3]s"
  type        = "STRING"
  unique      = true
  multivalued = true
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccSchemaAttributeConfig_JSONFull(resourceName, name string, unique, multivalued bool) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_schema_attribute" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

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

  multivalued = true
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccSchemaAttributeConfig_EnumeratedValues1(resourceName, name, attrType string) string {
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

func testAccSchemaAttributeConfig_EnumeratedValues2(resourceName, name, attrType string) string {
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
      archived = "true"
    },
    {
      value    = "value7"
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

  unique      = true
  multivalued = true

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

func testAccSchemaAttributeConfig_StandardDataSource(environmentName, licenseID string) string {
	return fmt.Sprintf(`
		%s

		data "pingone_schema" "user" {
		environment_id = pingone_environment.%s.id
		name           = "User"
		}

		data "pingone_schema_attribute" "email" {
		environment_id = pingone_environment.%s.id
		schema_id      = data.pingone_schema.user.id
		name           = "email"
		}
	`, legacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, environmentName)
}

func testAccSchemaAttributeConfig_StandardResource(environmentName, licenseID string, enabled bool) string {
	return fmt.Sprintf(`
		%s

		data "pingone_schema" "user" {
		environment_id = pingone_environment.%s.id
		name           = "User"
		}

		resource "pingone_schema_attribute" "email" {
		environment_id = pingone_environment.%s.id

		name         = "email"
		enabled      = %t
		}
	`, legacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, environmentName, enabled)
}
