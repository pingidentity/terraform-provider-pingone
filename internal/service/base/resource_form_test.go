package base_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/base"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccForm_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_form.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	var formID, environmentID string

	var p1Client *client.Client
	var ctx = context.Background()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)

			p1Client = acctest.PreCheckTestClient(ctx, t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Form_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccFormConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check:  base.Form_GetIDs(resourceFullName, &environmentID, &formID),
			},
			{
				PreConfig: func() {
					base.Form_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID, formID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the environment
			{
				Config: testAccFormConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check:  base.Form_GetIDs(resourceFullName, &environmentID, &formID),
			},
			{
				PreConfig: func() {
					base.Environment_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccForm_NewEnv(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_decision_endpoint.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Form_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFormConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
				),
			},
		},
	})
}

func TestAccForm_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_form.%s", resourceName)

	name := resourceName

	fullStep := resource.TestStep{
		Config: testAccFormConfig_Full(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestCheckResourceAttr(resourceFullName, "description", "This is my awesome form"),
			resource.TestCheckResourceAttr(resourceFullName, "category", "CUSTOM"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_required", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_optional", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "cols", "4"),
			resource.TestCheckResourceAttr(resourceFullName, "language_bundle", "test"),
			resource.TestCheckResourceAttr(resourceFullName, "translation_method", "DEFAULT_VALUE"),
		),
	}

	minimalStep := resource.TestStep{
		Config: testAccFormConfig_Minimal(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestCheckNoResourceAttr(resourceFullName, "description"),
			resource.TestCheckResourceAttr(resourceFullName, "category", "CUSTOM"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_required", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_optional", "false"),
			resource.TestCheckNoResourceAttr(resourceFullName, "cols"),
			resource.TestCheckNoResourceAttr(resourceFullName, "language_bundle"),
			resource.TestCheckNoResourceAttr(resourceFullName, "translation_method"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Form_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full step
			fullStep,
			{
				Config:  testAccFormConfig_Full(resourceName, name),
				Destroy: true,
			},
			// Minimal step
			minimalStep,
			{
				Config:  testAccFormConfig_Minimal(resourceName, name),
				Destroy: true,
			},
			// Change
			fullStep,
			minimalStep,
			fullStep,
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccForm_Multiple(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_form.%s", resourceName)

	name := resourceName

	step1 := resource.TestStep{
		Config: testAccFormConfig_MultipleStep1(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.#", "6"),
			resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "components.fields.*", map[string]string{}),
		),
	}

	step2 := resource.TestStep{
		Config: testAccFormConfig_MultipleStep2(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.#", "4"),
			resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "components.fields.*", map[string]string{}),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Form_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			step1,
			step2,
			step1,
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccForm_FieldText(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_form.%s", resourceName)

	name := resourceName

	fullStep := resource.TestStep{
		Config: testAccFormConfig_FieldTextFull(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "components.fields.*", map[string]string{
				"position.row":                        "0",
				"position.col":                        "0",
				"field_text.label":                    "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"fields.user.username.label\",\"defaultTranslation\":\"Username\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]",
				"field_text.key":                      "user.username",
				"field_text.mode":                     "ATTRIBUTE_MODE",
				"field_text.required":                 "true",
				"field_text.validation.type":          "CUSTOM",
				"field_text.validation.regex":         "[a-zA-Z0-9]+",
				"field_text.validation.error_message": "Regex validation error test message",
			}),
		),
	}

	minimalStep := resource.TestStep{
		Config: testAccFormConfig_FieldTextMinimal(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "components.fields.*", map[string]string{
				"position.row":        "0",
				"position.col":        "0",
				"field_text.key":      "1234567",
				"field_text.required": "false",
			}),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Form_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full step
			fullStep,
			{
				Config:  testAccFormConfig_FieldTextFull(resourceName, name),
				Destroy: true,
			},
			// Minimal step
			minimalStep,
			{
				Config:  testAccFormConfig_FieldTextMinimal(resourceName, name),
				Destroy: true,
			},
			// Change
			fullStep,
			minimalStep,
			fullStep,
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccForm_FieldPassword(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_form.%s", resourceName)

	name := resourceName

	fullStep := resource.TestStep{
		Config: testAccFormConfig_FieldPasswordFull(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestCheckResourceAttr(resourceFullName, "description", "This is my awesome form"),
			resource.TestCheckResourceAttr(resourceFullName, "category", "CUSTOM"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_required", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_optional", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "cols", "4"),
			resource.TestCheckResourceAttr(resourceFullName, "language_bundle", "test"),
			resource.TestCheckResourceAttr(resourceFullName, "translation_method", "DEFAULT_VALUE"),
		),
	}

	minimalStep := resource.TestStep{
		Config: testAccFormConfig_FieldPasswordMinimal(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestCheckNoResourceAttr(resourceFullName, "description"),
			resource.TestCheckResourceAttr(resourceFullName, "category", "CUSTOM"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_required", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_optional", "false"),
			resource.TestCheckNoResourceAttr(resourceFullName, "cols"),
			resource.TestCheckNoResourceAttr(resourceFullName, "language_bundle"),
			resource.TestCheckNoResourceAttr(resourceFullName, "translation_method"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Form_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full step
			fullStep,
			{
				Config:  testAccFormConfig_FieldPasswordFull(resourceName, name),
				Destroy: true,
			},
			// Minimal step
			minimalStep,
			{
				Config:  testAccFormConfig_FieldPasswordMinimal(resourceName, name),
				Destroy: true,
			},
			// Change
			fullStep,
			minimalStep,
			fullStep,
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccForm_FieldPasswordVerify(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_form.%s", resourceName)

	name := resourceName

	fullStep := resource.TestStep{
		Config: testAccFormConfig_FieldPasswordVerifyFull(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestCheckResourceAttr(resourceFullName, "description", "This is my awesome form"),
			resource.TestCheckResourceAttr(resourceFullName, "category", "CUSTOM"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_required", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_optional", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "cols", "4"),
			resource.TestCheckResourceAttr(resourceFullName, "language_bundle", "test"),
			resource.TestCheckResourceAttr(resourceFullName, "translation_method", "DEFAULT_VALUE"),
		),
	}

	minimalStep := resource.TestStep{
		Config: testAccFormConfig_FieldPasswordVerifyMinimal(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestCheckNoResourceAttr(resourceFullName, "description"),
			resource.TestCheckResourceAttr(resourceFullName, "category", "CUSTOM"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_required", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_optional", "false"),
			resource.TestCheckNoResourceAttr(resourceFullName, "cols"),
			resource.TestCheckNoResourceAttr(resourceFullName, "language_bundle"),
			resource.TestCheckNoResourceAttr(resourceFullName, "translation_method"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Form_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full step
			fullStep,
			{
				Config:  testAccFormConfig_FieldPasswordVerifyFull(resourceName, name),
				Destroy: true,
			},
			// Minimal step
			minimalStep,
			{
				Config:  testAccFormConfig_FieldPasswordVerifyMinimal(resourceName, name),
				Destroy: true,
			},
			// Change
			fullStep,
			minimalStep,
			fullStep,
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccForm_FieldRadio(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_form.%s", resourceName)

	name := resourceName

	fullStep := resource.TestStep{
		Config: testAccFormConfig_FieldRadioFull(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestCheckResourceAttr(resourceFullName, "description", "This is my awesome form"),
			resource.TestCheckResourceAttr(resourceFullName, "category", "CUSTOM"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_required", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_optional", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "cols", "4"),
			resource.TestCheckResourceAttr(resourceFullName, "language_bundle", "test"),
			resource.TestCheckResourceAttr(resourceFullName, "translation_method", "DEFAULT_VALUE"),
		),
	}

	minimalStep := resource.TestStep{
		Config: testAccFormConfig_FieldRadioMinimal(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestCheckNoResourceAttr(resourceFullName, "description"),
			resource.TestCheckResourceAttr(resourceFullName, "category", "CUSTOM"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_required", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_optional", "false"),
			resource.TestCheckNoResourceAttr(resourceFullName, "cols"),
			resource.TestCheckNoResourceAttr(resourceFullName, "language_bundle"),
			resource.TestCheckNoResourceAttr(resourceFullName, "translation_method"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Form_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full step
			fullStep,
			{
				Config:  testAccFormConfig_FieldRadioFull(resourceName, name),
				Destroy: true,
			},
			// Minimal step
			minimalStep,
			{
				Config:  testAccFormConfig_FieldRadioMinimal(resourceName, name),
				Destroy: true,
			},
			// Change
			fullStep,
			minimalStep,
			fullStep,
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccForm_FieldCheckbox(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_form.%s", resourceName)

	name := resourceName

	fullStep := resource.TestStep{
		Config: testAccFormConfig_FieldCheckboxFull(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestCheckResourceAttr(resourceFullName, "description", "This is my awesome form"),
			resource.TestCheckResourceAttr(resourceFullName, "category", "CUSTOM"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_required", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_optional", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "cols", "4"),
			resource.TestCheckResourceAttr(resourceFullName, "language_bundle", "test"),
			resource.TestCheckResourceAttr(resourceFullName, "translation_method", "DEFAULT_VALUE"),
		),
	}

	minimalStep := resource.TestStep{
		Config: testAccFormConfig_FieldCheckboxMinimal(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestCheckNoResourceAttr(resourceFullName, "description"),
			resource.TestCheckResourceAttr(resourceFullName, "category", "CUSTOM"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_required", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_optional", "false"),
			resource.TestCheckNoResourceAttr(resourceFullName, "cols"),
			resource.TestCheckNoResourceAttr(resourceFullName, "language_bundle"),
			resource.TestCheckNoResourceAttr(resourceFullName, "translation_method"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Form_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full step
			fullStep,
			{
				Config:  testAccFormConfig_FieldCheckboxFull(resourceName, name),
				Destroy: true,
			},
			// Minimal step
			minimalStep,
			{
				Config:  testAccFormConfig_FieldCheckboxMinimal(resourceName, name),
				Destroy: true,
			},
			// Change
			fullStep,
			minimalStep,
			fullStep,
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccForm_FieldDropdown(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_form.%s", resourceName)

	name := resourceName

	fullStep := resource.TestStep{
		Config: testAccFormConfig_FieldDropdownFull(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestCheckResourceAttr(resourceFullName, "description", "This is my awesome form"),
			resource.TestCheckResourceAttr(resourceFullName, "category", "CUSTOM"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_required", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_optional", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "cols", "4"),
			resource.TestCheckResourceAttr(resourceFullName, "language_bundle", "test"),
			resource.TestCheckResourceAttr(resourceFullName, "translation_method", "DEFAULT_VALUE"),
		),
	}

	minimalStep := resource.TestStep{
		Config: testAccFormConfig_FieldDropdownMinimal(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestCheckNoResourceAttr(resourceFullName, "description"),
			resource.TestCheckResourceAttr(resourceFullName, "category", "CUSTOM"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_required", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_optional", "false"),
			resource.TestCheckNoResourceAttr(resourceFullName, "cols"),
			resource.TestCheckNoResourceAttr(resourceFullName, "language_bundle"),
			resource.TestCheckNoResourceAttr(resourceFullName, "translation_method"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Form_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full step
			fullStep,
			{
				Config:  testAccFormConfig_FieldDropdownFull(resourceName, name),
				Destroy: true,
			},
			// Minimal step
			minimalStep,
			{
				Config:  testAccFormConfig_FieldDropdownMinimal(resourceName, name),
				Destroy: true,
			},
			// Change
			fullStep,
			minimalStep,
			fullStep,
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccForm_FieldCombobox(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_form.%s", resourceName)

	name := resourceName

	fullStep := resource.TestStep{
		Config: testAccFormConfig_FieldComboboxFull(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestCheckResourceAttr(resourceFullName, "description", "This is my awesome form"),
			resource.TestCheckResourceAttr(resourceFullName, "category", "CUSTOM"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_required", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_optional", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "cols", "4"),
			resource.TestCheckResourceAttr(resourceFullName, "language_bundle", "test"),
			resource.TestCheckResourceAttr(resourceFullName, "translation_method", "DEFAULT_VALUE"),
		),
	}

	minimalStep := resource.TestStep{
		Config: testAccFormConfig_FieldComboboxMinimal(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestCheckNoResourceAttr(resourceFullName, "description"),
			resource.TestCheckResourceAttr(resourceFullName, "category", "CUSTOM"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_required", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_optional", "false"),
			resource.TestCheckNoResourceAttr(resourceFullName, "cols"),
			resource.TestCheckNoResourceAttr(resourceFullName, "language_bundle"),
			resource.TestCheckNoResourceAttr(resourceFullName, "translation_method"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Form_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full step
			fullStep,
			{
				Config:  testAccFormConfig_FieldComboboxFull(resourceName, name),
				Destroy: true,
			},
			// Minimal step
			minimalStep,
			{
				Config:  testAccFormConfig_FieldComboboxMinimal(resourceName, name),
				Destroy: true,
			},
			// Change
			fullStep,
			minimalStep,
			fullStep,
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccForm_FieldDivider(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_form.%s", resourceName)

	name := resourceName

	fullStep := resource.TestStep{
		Config: testAccFormConfig_FieldDividerFull(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestCheckResourceAttr(resourceFullName, "description", "This is my awesome form"),
			resource.TestCheckResourceAttr(resourceFullName, "category", "CUSTOM"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_required", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_optional", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "cols", "4"),
			resource.TestCheckResourceAttr(resourceFullName, "language_bundle", "test"),
			resource.TestCheckResourceAttr(resourceFullName, "translation_method", "DEFAULT_VALUE"),
		),
	}

	minimalStep := resource.TestStep{
		Config: testAccFormConfig_FieldDividerMinimal(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestCheckNoResourceAttr(resourceFullName, "description"),
			resource.TestCheckResourceAttr(resourceFullName, "category", "CUSTOM"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_required", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_optional", "false"),
			resource.TestCheckNoResourceAttr(resourceFullName, "cols"),
			resource.TestCheckNoResourceAttr(resourceFullName, "language_bundle"),
			resource.TestCheckNoResourceAttr(resourceFullName, "translation_method"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Form_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full step
			fullStep,
			{
				Config:  testAccFormConfig_FieldDividerFull(resourceName, name),
				Destroy: true,
			},
			// Minimal step
			minimalStep,
			{
				Config:  testAccFormConfig_FieldDividerMinimal(resourceName, name),
				Destroy: true,
			},
			// Change
			fullStep,
			minimalStep,
			fullStep,
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccForm_FieldEmptyField(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_form.%s", resourceName)

	name := resourceName

	fullStep := resource.TestStep{
		Config: testAccFormConfig_FieldEmptyFieldFull(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestCheckResourceAttr(resourceFullName, "description", "This is my awesome form"),
			resource.TestCheckResourceAttr(resourceFullName, "category", "CUSTOM"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_required", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_optional", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "cols", "4"),
			resource.TestCheckResourceAttr(resourceFullName, "language_bundle", "test"),
			resource.TestCheckResourceAttr(resourceFullName, "translation_method", "DEFAULT_VALUE"),
		),
	}

	minimalStep := resource.TestStep{
		Config: testAccFormConfig_FieldEmptyFieldMinimal(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestCheckNoResourceAttr(resourceFullName, "description"),
			resource.TestCheckResourceAttr(resourceFullName, "category", "CUSTOM"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_required", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_optional", "false"),
			resource.TestCheckNoResourceAttr(resourceFullName, "cols"),
			resource.TestCheckNoResourceAttr(resourceFullName, "language_bundle"),
			resource.TestCheckNoResourceAttr(resourceFullName, "translation_method"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Form_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full step
			fullStep,
			{
				Config:  testAccFormConfig_FieldEmptyFieldFull(resourceName, name),
				Destroy: true,
			},
			// Minimal step
			minimalStep,
			{
				Config:  testAccFormConfig_FieldEmptyFieldMinimal(resourceName, name),
				Destroy: true,
			},
			// Change
			fullStep,
			minimalStep,
			fullStep,
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccForm_FieldTextblob(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_form.%s", resourceName)

	name := resourceName

	fullStep := resource.TestStep{
		Config: testAccFormConfig_FieldTextBlobFull(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestCheckResourceAttr(resourceFullName, "description", "This is my awesome form"),
			resource.TestCheckResourceAttr(resourceFullName, "category", "CUSTOM"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_required", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_optional", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "cols", "4"),
			resource.TestCheckResourceAttr(resourceFullName, "language_bundle", "test"),
			resource.TestCheckResourceAttr(resourceFullName, "translation_method", "DEFAULT_VALUE"),
		),
	}

	minimalStep := resource.TestStep{
		Config: testAccFormConfig_FieldTextBlobMinimal(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestCheckNoResourceAttr(resourceFullName, "description"),
			resource.TestCheckResourceAttr(resourceFullName, "category", "CUSTOM"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_required", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_optional", "false"),
			resource.TestCheckNoResourceAttr(resourceFullName, "cols"),
			resource.TestCheckNoResourceAttr(resourceFullName, "language_bundle"),
			resource.TestCheckNoResourceAttr(resourceFullName, "translation_method"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Form_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full step
			fullStep,
			{
				Config:  testAccFormConfig_FieldTextBlobFull(resourceName, name),
				Destroy: true,
			},
			// Minimal step
			minimalStep,
			{
				Config:  testAccFormConfig_FieldTextBlobMinimal(resourceName, name),
				Destroy: true,
			},
			// Change
			fullStep,
			minimalStep,
			fullStep,
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccForm_FieldSlateTextblob(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_form.%s", resourceName)

	name := resourceName

	fullStep := resource.TestStep{
		Config: testAccFormConfig_FieldSlateTextBlobFull(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestCheckResourceAttr(resourceFullName, "description", "This is my awesome form"),
			resource.TestCheckResourceAttr(resourceFullName, "category", "CUSTOM"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_required", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_optional", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "cols", "4"),
			resource.TestCheckResourceAttr(resourceFullName, "language_bundle", "test"),
			resource.TestCheckResourceAttr(resourceFullName, "translation_method", "DEFAULT_VALUE"),
		),
	}

	minimalStep := resource.TestStep{
		Config: testAccFormConfig_FieldSlateTextBlobMinimal(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestCheckNoResourceAttr(resourceFullName, "description"),
			resource.TestCheckResourceAttr(resourceFullName, "category", "CUSTOM"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_required", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_optional", "false"),
			resource.TestCheckNoResourceAttr(resourceFullName, "cols"),
			resource.TestCheckNoResourceAttr(resourceFullName, "language_bundle"),
			resource.TestCheckNoResourceAttr(resourceFullName, "translation_method"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Form_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full step
			fullStep,
			{
				Config:  testAccFormConfig_FieldSlateTextBlobFull(resourceName, name),
				Destroy: true,
			},
			// Minimal step
			minimalStep,
			{
				Config:  testAccFormConfig_FieldSlateTextBlobMinimal(resourceName, name),
				Destroy: true,
			},
			// Change
			fullStep,
			minimalStep,
			fullStep,
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccForm_FieldSubmitButton(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_form.%s", resourceName)

	name := resourceName

	fullStep := resource.TestStep{
		Config: testAccFormConfig_FieldSubmitButtonFull(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestCheckResourceAttr(resourceFullName, "description", "This is my awesome form"),
			resource.TestCheckResourceAttr(resourceFullName, "category", "CUSTOM"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_required", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_optional", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "cols", "4"),
			resource.TestCheckResourceAttr(resourceFullName, "language_bundle", "test"),
			resource.TestCheckResourceAttr(resourceFullName, "translation_method", "DEFAULT_VALUE"),
		),
	}

	minimalStep := resource.TestStep{
		Config: testAccFormConfig_FieldSubmitButtonMinimal(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestCheckNoResourceAttr(resourceFullName, "description"),
			resource.TestCheckResourceAttr(resourceFullName, "category", "CUSTOM"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_required", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_optional", "false"),
			resource.TestCheckNoResourceAttr(resourceFullName, "cols"),
			resource.TestCheckNoResourceAttr(resourceFullName, "language_bundle"),
			resource.TestCheckNoResourceAttr(resourceFullName, "translation_method"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Form_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full step
			fullStep,
			{
				Config:  testAccFormConfig_FieldSubmitButtonFull(resourceName, name),
				Destroy: true,
			},
			// Minimal step
			minimalStep,
			{
				Config:  testAccFormConfig_FieldSubmitButtonMinimal(resourceName, name),
				Destroy: true,
			},
			// Change
			fullStep,
			minimalStep,
			fullStep,
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccForm_FieldErrorDisplay(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_form.%s", resourceName)

	name := resourceName

	fullStep := resource.TestStep{
		Config: testAccFormConfig_FieldErrorDisplayFull(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestCheckResourceAttr(resourceFullName, "description", "This is my awesome form"),
			resource.TestCheckResourceAttr(resourceFullName, "category", "CUSTOM"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_required", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_optional", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "cols", "4"),
			resource.TestCheckResourceAttr(resourceFullName, "language_bundle", "test"),
			resource.TestCheckResourceAttr(resourceFullName, "translation_method", "DEFAULT_VALUE"),
		),
	}

	minimalStep := resource.TestStep{
		Config: testAccFormConfig_FieldErrorDisplayMinimal(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestCheckNoResourceAttr(resourceFullName, "description"),
			resource.TestCheckResourceAttr(resourceFullName, "category", "CUSTOM"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_required", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_optional", "false"),
			resource.TestCheckNoResourceAttr(resourceFullName, "cols"),
			resource.TestCheckNoResourceAttr(resourceFullName, "language_bundle"),
			resource.TestCheckNoResourceAttr(resourceFullName, "translation_method"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Form_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full step
			fullStep,
			{
				Config:  testAccFormConfig_FieldErrorDisplayFull(resourceName, name),
				Destroy: true,
			},
			// Minimal step
			minimalStep,
			{
				Config:  testAccFormConfig_FieldErrorDisplayMinimal(resourceName, name),
				Destroy: true,
			},
			// Change
			fullStep,
			minimalStep,
			fullStep,
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccForm_FieldFlowLink(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_form.%s", resourceName)

	name := resourceName

	fullStep := resource.TestStep{
		Config: testAccFormConfig_FieldFlowLinkFull(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestCheckResourceAttr(resourceFullName, "description", "This is my awesome form"),
			resource.TestCheckResourceAttr(resourceFullName, "category", "CUSTOM"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_required", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_optional", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "cols", "4"),
			resource.TestCheckResourceAttr(resourceFullName, "language_bundle", "test"),
			resource.TestCheckResourceAttr(resourceFullName, "translation_method", "DEFAULT_VALUE"),
		),
	}

	minimalStep := resource.TestStep{
		Config: testAccFormConfig_FieldFlowLinkMinimal(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestCheckNoResourceAttr(resourceFullName, "description"),
			resource.TestCheckResourceAttr(resourceFullName, "category", "CUSTOM"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_required", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_optional", "false"),
			resource.TestCheckNoResourceAttr(resourceFullName, "cols"),
			resource.TestCheckNoResourceAttr(resourceFullName, "language_bundle"),
			resource.TestCheckNoResourceAttr(resourceFullName, "translation_method"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Form_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full step
			fullStep,
			{
				Config:  testAccFormConfig_FieldFlowLinkFull(resourceName, name),
				Destroy: true,
			},
			// Minimal step
			minimalStep,
			{
				Config:  testAccFormConfig_FieldFlowLinkMinimal(resourceName, name),
				Destroy: true,
			},
			// Change
			fullStep,
			minimalStep,
			fullStep,
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccForm_FieldFlowButton(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_form.%s", resourceName)

	name := resourceName

	fullStep := resource.TestStep{
		Config: testAccFormConfig_FieldFlowButtonFull(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestCheckResourceAttr(resourceFullName, "description", "This is my awesome form"),
			resource.TestCheckResourceAttr(resourceFullName, "category", "CUSTOM"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_required", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_optional", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "cols", "4"),
			resource.TestCheckResourceAttr(resourceFullName, "language_bundle", "test"),
			resource.TestCheckResourceAttr(resourceFullName, "translation_method", "DEFAULT_VALUE"),
		),
	}

	minimalStep := resource.TestStep{
		Config: testAccFormConfig_FieldFlowButtonMinimal(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestCheckNoResourceAttr(resourceFullName, "description"),
			resource.TestCheckResourceAttr(resourceFullName, "category", "CUSTOM"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_required", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_optional", "false"),
			resource.TestCheckNoResourceAttr(resourceFullName, "cols"),
			resource.TestCheckNoResourceAttr(resourceFullName, "language_bundle"),
			resource.TestCheckNoResourceAttr(resourceFullName, "translation_method"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Form_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full step
			fullStep,
			{
				Config:  testAccFormConfig_FieldFlowButtonFull(resourceName, name),
				Destroy: true,
			},
			// Minimal step
			minimalStep,
			{
				Config:  testAccFormConfig_FieldFlowButtonMinimal(resourceName, name),
				Destroy: true,
			},
			// Change
			fullStep,
			minimalStep,
			fullStep,
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccForm_FieldRecaptchaV2(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_form.%s", resourceName)

	name := resourceName

	fullStep := resource.TestStep{
		Config: testAccFormConfig_FieldRecaptchaV2Full(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestCheckResourceAttr(resourceFullName, "description", "This is my awesome form"),
			resource.TestCheckResourceAttr(resourceFullName, "category", "CUSTOM"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_required", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_optional", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "cols", "4"),
			resource.TestCheckResourceAttr(resourceFullName, "language_bundle", "test"),
			resource.TestCheckResourceAttr(resourceFullName, "translation_method", "DEFAULT_VALUE"),
		),
	}

	minimalStep := resource.TestStep{
		Config: testAccFormConfig_FieldRecaptchaV2Minimal(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestCheckNoResourceAttr(resourceFullName, "description"),
			resource.TestCheckResourceAttr(resourceFullName, "category", "CUSTOM"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_required", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_optional", "false"),
			resource.TestCheckNoResourceAttr(resourceFullName, "cols"),
			resource.TestCheckNoResourceAttr(resourceFullName, "language_bundle"),
			resource.TestCheckNoResourceAttr(resourceFullName, "translation_method"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Form_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full step
			fullStep,
			{
				Config:  testAccFormConfig_FieldRecaptchaV2Full(resourceName, name),
				Destroy: true,
			},
			// Minimal step
			minimalStep,
			{
				Config:  testAccFormConfig_FieldRecaptchaV2Minimal(resourceName, name),
				Destroy: true,
			},
			// Change
			fullStep,
			minimalStep,
			fullStep,
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccForm_FieldQrCode(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_form.%s", resourceName)

	name := resourceName

	fullStep := resource.TestStep{
		Config: testAccFormConfig_FieldQrCodeFull(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestCheckResourceAttr(resourceFullName, "description", "This is my awesome form"),
			resource.TestCheckResourceAttr(resourceFullName, "category", "CUSTOM"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_required", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_optional", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "cols", "4"),
			resource.TestCheckResourceAttr(resourceFullName, "language_bundle", "test"),
			resource.TestCheckResourceAttr(resourceFullName, "translation_method", "DEFAULT_VALUE"),
		),
	}

	minimalStep := resource.TestStep{
		Config: testAccFormConfig_FieldQrCodeMinimal(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestCheckNoResourceAttr(resourceFullName, "description"),
			resource.TestCheckResourceAttr(resourceFullName, "category", "CUSTOM"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_required", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_optional", "false"),
			resource.TestCheckNoResourceAttr(resourceFullName, "cols"),
			resource.TestCheckNoResourceAttr(resourceFullName, "language_bundle"),
			resource.TestCheckNoResourceAttr(resourceFullName, "translation_method"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Form_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full step
			fullStep,
			{
				Config:  testAccFormConfig_FieldQrCodeFull(resourceName, name),
				Destroy: true,
			},
			// Minimal step
			minimalStep,
			{
				Config:  testAccFormConfig_FieldQrCodeMinimal(resourceName, name),
				Destroy: true,
			},
			// Change
			fullStep,
			minimalStep,
			fullStep,
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccForm_FieldSocialLoginButton(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_form.%s", resourceName)

	name := resourceName

	fullStep := resource.TestStep{
		Config: testAccFormConfig_FieldSocialLoginButtonFull(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestCheckResourceAttr(resourceFullName, "description", "This is my awesome form"),
			resource.TestCheckResourceAttr(resourceFullName, "category", "CUSTOM"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_required", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_optional", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "cols", "4"),
			resource.TestCheckResourceAttr(resourceFullName, "language_bundle", "test"),
			resource.TestCheckResourceAttr(resourceFullName, "translation_method", "DEFAULT_VALUE"),
		),
	}

	minimalStep := resource.TestStep{
		Config: testAccFormConfig_FieldSocialLoginButtonMinimal(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestCheckNoResourceAttr(resourceFullName, "description"),
			resource.TestCheckResourceAttr(resourceFullName, "category", "CUSTOM"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_required", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_optional", "false"),
			resource.TestCheckNoResourceAttr(resourceFullName, "cols"),
			resource.TestCheckNoResourceAttr(resourceFullName, "language_bundle"),
			resource.TestCheckNoResourceAttr(resourceFullName, "translation_method"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Form_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full step
			fullStep,
			{
				Config:  testAccFormConfig_FieldSocialLoginButtonFull(resourceName, name),
				Destroy: true,
			},
			// Minimal step
			minimalStep,
			{
				Config:  testAccFormConfig_FieldSocialLoginButtonMinimal(resourceName, name),
				Destroy: true,
			},
			// Change
			fullStep,
			minimalStep,
			fullStep,
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccForm_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_form.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Form_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccFormConfig_Minimal(resourceName, name),
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
				ImportStateId: "badformat/badformat",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`Unexpected Import Identifier`),
			},
		},
	})
}

func testAccFormConfig_NewEnv(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_form" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  name = "My Awesome Minimal Form"

  mark_required = true
  mark_optional = false

  components = {
    fields = [
      {
        position = {
          row = 0
          col = 0
        }

        field_text = {
          key = "user.username"
          validation = {
            type = "NONE"
          }
          required = true
        }
      }
    ]
  }
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccFormConfig_Full(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_form" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name        = "%[3]s"
  description = "This is my awesome form"

  category = "CUSTOM"

  mark_required = false
  mark_optional = true

  cols = 4

  language_bundle = {
    "button.text"                              = "Submit",
    "fields.user.email.label"                  = "Email Address",
    "fields.user.password.label"               = "Password"
    "fields.user.password.labelPasswordVerify" = "Verify Password",
    "fields.user.username.label"               = "Username",
  }

  translation_method = "DEFAULT_VALUE"

  components = {
    fields = [
      {
        position = {
          row = 0
          col = 0
        }

        field_text = {
          key = "user.username"
          validation = {
            type = "NONE"
          }
          required = true
        }
      }
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccFormConfig_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_form" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  mark_required = true
  mark_optional = false

  components = {
    fields = [
      {
        position = {
          row = 0
          col = 0
        }

        field_text = {
          key = "user.username"
          validation = {
            type = "NONE"
          }
          required = true
        }
      }
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccFormConfig_MultipleStep1(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_form" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  mark_required = true
  mark_optional = false

  components = {
    fields = [
      {
        position = {
          row = 0
          col = 0
        }

        field_slate_textblob = {
          key     = "c8301910-4539-4980-a113-9120ebdd6bd5"
          content = "[{\"children\":[{\"text\":\"Create Your Profile\"}],\"type\":\"heading-1\"},{\"children\":[{\"text\":\"Enter the required information below\"}]},{\"type\":\"divider\",\"children\":[{\"text\":\"\"}]},{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"}]}]"
        }
      },
      {
        position = {
          row = 1
          col = 0
        }

        field_error_display = {
          key = "1e16e184-a8ad-40d6-b87d-866ca41df39e"
        }
      },
      {
        position = {
          row = 2
          col = 0
        }

        field_text = {
          label    = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"fields.user.username.label\",\"defaultTranslation\":\"Username\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]"
          key      = "user.username"
          mode     = "ATTRIBUTE_MODE"
          required = true
          validation = {
            type = "NONE"
          }
        }
      },
      {
        position = {
          row = 3
          col = 0
        }

        field_text = {
          label    = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"fields.user.email.label\",\"defaultTranslation\":\"Email Address\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]"
          key      = "user.email"
          mode     = "ATTRIBUTE_MODE"
          required = true
          validation = {
            type = "NONE"
          }
        }
      },
      {
        "position" : {
          row = 4
          col = 0
        }

        field_password_verify = {
          label               = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"fields.user.password.label\",\"defaultTranslation\":\"Password\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]"
          labelPasswordVerify = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"fields.user.password.labelPasswordVerify\",\"defaultTranslation\":\"Verify Password\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]"
          key                 = "user.password"
          mode                = "ATTRIBUTE_MODE"
          required            = true
          validation = {
            type = "NONE"
          },
          show_password_requirements : true
        }
      },
      {
        position = {
          row = 5
          col = 0
        }

        field_submit_button = {
          key   = "697dc4e9-acf5-4c04-9f7e-3974ea999c37"
          label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"button.text\",\"defaultTranslation\":\"Submit\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]"
        }
      }
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccFormConfig_MultipleStep2(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_form" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  mark_required = true
  mark_optional = false

  components = {
    fields = [
      {
        position = {
          row = 0
          col = 0
        }

        field_slate_textblob = {
          key     = "c8301910-4539-4980-a113-9120ebdd6bd5"
          content = "[{\"children\":[{\"text\":\"Create Your Profile\"}],\"type\":\"heading-1\"},{\"children\":[{\"text\":\"Enter the required information below\"}]},{\"type\":\"divider\",\"children\":[{\"text\":\"\"}]},{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"}]}]"
        }
      },
      {
        position = {
          row = 2
          col = 0
        }

        field_text = {
          label    = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"fields.user.email.label\",\"defaultTranslation\":\"Email Address\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]"
          key      = "user.email"
          mode     = "ATTRIBUTE_MODE"
          required = true
          validation = {
            type = "NONE"
          }
        }
      },
      {
        position = {
          row = 1
          col = 0
        }

        field_text = {
          label    = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"fields.user.username.label\",\"defaultTranslation\":\"Username\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]"
          key      = "user.username"
          mode     = "ATTRIBUTE_MODE"
          required = true
          validation = {
            type = "NONE"
          }
        }
      },
      {
        position = {
          row = 3
          col = 0
        }

        field_submit_button = {
          key   = "697dc4e9-acf5-4c04-9f7e-3974ea999c37"
          label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"button.text\",\"defaultTranslation\":\"Submit\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]"
        }
      }
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccFormConfig_FieldTextFull(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_form" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  mark_required = true
  mark_optional = false

  components = {
    fields = [
      {
        position = {
          row = 0
          col = 0
        }

        field_text = {
          label              = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"fields.user.username.label\",\"defaultTranslation\":\"Username\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]"
          label_mode         = "FLOAT"
          layout             = "VERTICAL"
          options            = ["test", "test1", "test2", "test3"]
          key                = "user.username"
          mode               = "ATTRIBUTE_MODE"
          required           = true
          attribute_disabled = true
          validation = {
            type          = "CUSTOM"
            regex         = "[a-zA-Z0-9]+"
            error_message = "Regex validation error test message"
          }

          other_option_enabled            = true
          other_option_key                = "key.123"
          other_option_label              = "Test label 432"
          other_option_input_label        = "Test label 123"
          other_option_attribute_disabled = true
        }
      },
      {
        position = {
          row = 2
          col = 0
        }

        field_submit_button = {
          key   = "697dc4e9-acf5-4c04-9f7e-3974ea999c37"
          label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"button.text\",\"defaultTranslation\":\"Submit\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]"
        }
      }
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccFormConfig_FieldTextMinimal(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_form" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  mark_required = true
  mark_optional = false

  components = {
    fields = [
      {
        position = {
          row = 0
          col = 0
        }

        field_text = {
          key      = "1234567"
          required = false
        }
      },
      {
        position = {
          row = 2
          col = 0
        }

        field_submit_button = {
          key   = "697dc4e9-acf5-4c04-9f7e-3974ea999c37"
          label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"button.text\",\"defaultTranslation\":\"Submit\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]"
        }
      }
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccFormConfig_FieldPasswordFull(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_form" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  mark_required = true
  mark_optional = false

  components = {
    fields = [
      {
        position = {
          row = 0
          col = 0
        }

        field_password = {
          label    = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"fields.user.username.label\",\"defaultTranslation\":\"Username\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]"
          key      = "user.username"
          mode     = "ATTRIBUTE_MODE"
          required = true
          validation = {
            type          = "CUSTOM"
            regex         = "[a-zA-Z0-9]+"
            error_message = "Regex validation error test message"
          }
        }
      },
      {
        position = {
          row = 2
          col = 0
        }

        field_submit_button = {
          key   = "697dc4e9-acf5-4c04-9f7e-3974ea999c37"
          label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"button.text\",\"defaultTranslation\":\"Submit\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]"
        }
      }
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccFormConfig_FieldPasswordMinimal(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_form" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  mark_required = true
  mark_optional = false

  components = {
    fields = [
      {
        position = {
          row = 0
          col = 0
        }

        field_password = {
          key      = "1234567"
          required = false
        }
      },
      {
        position = {
          row = 2
          col = 0
        }

        field_submit_button = {
          key   = "697dc4e9-acf5-4c04-9f7e-3974ea999c37"
          label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"button.text\",\"defaultTranslation\":\"Submit\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]"
        }
      }
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
