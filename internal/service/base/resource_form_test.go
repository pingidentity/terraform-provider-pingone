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
	resourceFullName := fmt.Sprintf("pingone_form.%s", resourceName)

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
			resource.TestCheckResourceAttr(resourceFullName, "mark_required", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_optional", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "cols", "4"),
			resource.TestCheckResourceAttr(resourceFullName, "language_bundle.%", "1"),
			resource.TestCheckResourceAttr(resourceFullName, "language_bundle.button.text", "Submit"),
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
			resource.TestCheckResourceAttr(resourceFullName, "mark_required", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "mark_optional", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "cols", "4"),
			resource.TestCheckResourceAttr(resourceFullName, "language_bundle.%", "1"),
			resource.TestCheckResourceAttr(resourceFullName, "language_bundle.button.text", "Submit"),
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

// TODO
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

func TestAccForm_FieldCheckbox(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_form.%s", resourceName)

	name := resourceName

	fullStep := resource.TestStep{
		Config: testAccFormConfig_FieldCheckboxFull(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "components.fields.*", map[string]string{
				"position.row":                    "0",
				"position.col":                    "0",
				"position.width":                  "50",
				"type":                            "CHECKBOX",
				"label":                           "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"fields.user.locale.label\",\"defaultTranslation\":\"Locale\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]",
				"label_mode":                      "FLOAT",
				"layout":                          "VERTICAL",
				"key":                             fmt.Sprintf("user.%s", name),
				"required":                        "true",
				"attribute_disabled":              "false",
				"other_option_enabled":            "false",
				"other_option_attribute_disabled": "false",
				"options.0.value":                 "Option1",
				"options.0.label":                 "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Option 1\"}]}]",
				"options.1.value":                 "Option2",
				"options.1.label":                 "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Option 2\"}]}]",
				"options.2.value":                 "Option3",
				"options.2.label":                 "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Option 3\"}]}]",
			}),
		),
	}

	minimalStep := resource.TestStep{
		Config: testAccFormConfig_FieldCheckboxMinimal(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "components.fields.*", map[string]string{
				"position.row":                    "0",
				"position.col":                    "0",
				"type":                            "CHECKBOX",
				"label":                           "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Placeholder\"}]}]",
				"key":                             "checkbox-field",
				"layout":                          "HORIZONTAL",
				"required":                        "false",
				"attribute_disabled":              "false",
				"other_option_enabled":            "false",
				"other_option_attribute_disabled": "false",
				"options.0.value":                 "Option1",
				"options.0.label":                 "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Option 1\"}]}]",
				"options.1.value":                 "Option3",
				"options.1.label":                 "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Option 3\"}]}]",
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
			resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "components.fields.*", map[string]string{
				"position.row":                    "0",
				"position.col":                    "0",
				"position.width":                  "50",
				"type":                            "DROPDOWN",
				"label":                           "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"fields.user.locale.label\",\"defaultTranslation\":\"Locale\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]",
				"label_mode":                      "FLOAT",
				"layout":                          "VERTICAL",
				"key":                             "user.locale",
				"required":                        "true",
				"attribute_disabled":              "false",
				"other_option_enabled":            "false",
				"other_option_attribute_disabled": "false",
				"options.0.value":                 "Option1",
				"options.0.label":                 "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Option 1\"}]}]",
				"options.1.value":                 "Option2",
				"options.1.label":                 "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Option 2\"}]}]",
				"options.2.value":                 "Option3",
				"options.2.label":                 "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Option 3\"}]}]",
			}),
		),
	}

	minimalStep := resource.TestStep{
		Config: testAccFormConfig_FieldDropdownMinimal(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "components.fields.*", map[string]string{
				"position.row":                    "0",
				"position.col":                    "0",
				"type":                            "DROPDOWN",
				"label":                           "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Placeholder\"}]}]",
				"key":                             "dropdown-field",
				"required":                        "false",
				"attribute_disabled":              "false",
				"other_option_enabled":            "false",
				"other_option_attribute_disabled": "false",
				"options.0.value":                 "Option1",
				"options.0.label":                 "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Option 1\"}]}]",
				"options.1.value":                 "Option3",
				"options.1.label":                 "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Option 3\"}]}]",
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

func TestAccForm_FieldPassword(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_form.%s", resourceName)

	name := resourceName

	fullStep := resource.TestStep{
		Config: testAccFormConfig_FieldPasswordFull(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "components.fields.*", map[string]string{
				"position.row":                    "0",
				"position.col":                    "0",
				"position.width":                  "50",
				"type":                            "PASSWORD",
				"label":                           "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"fields.user.password.label\",\"defaultTranslation\":\"Password\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]",
				"label_mode":                      "FLOAT",
				"layout":                          "VERTICAL",
				"key":                             "user.password",
				"required":                        "true",
				"attribute_disabled":              "false",
				"other_option_enabled":            "false",
				"other_option_attribute_disabled": "false",
				"show_password_requirements":      "true",
			}),
		),
	}

	minimalStep := resource.TestStep{
		Config: testAccFormConfig_FieldPasswordMinimal(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "components.fields.*", map[string]string{
				"position.row":                    "0",
				"position.col":                    "0",
				"type":                            "PASSWORD",
				"label":                           "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Placeholder\"}]}]",
				"key":                             "password-field",
				"required":                        "false",
				"attribute_disabled":              "false",
				"other_option_enabled":            "false",
				"other_option_attribute_disabled": "false",
				"show_password_requirements":      "false",
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
			resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "components.fields.*", map[string]string{
				"position.row":                    "0",
				"position.col":                    "0",
				"position.width":                  "50",
				"type":                            "PASSWORD_VERIFY",
				"label":                           "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"fields.user.password.label\",\"defaultTranslation\":\"Password\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]",
				"label_mode":                      "FLOAT",
				"label_password_verify":           "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"fields.user.password.labelPasswordVerify\",\"defaultTranslation\":\"Verify Password\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]",
				"layout":                          "VERTICAL",
				"key":                             "user.password",
				"required":                        "true",
				"attribute_disabled":              "false",
				"other_option_enabled":            "false",
				"other_option_attribute_disabled": "false",
				"show_password_requirements":      "true",
			}),
		),
	}

	minimalStep := resource.TestStep{
		Config: testAccFormConfig_FieldPasswordVerifyMinimal(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "components.fields.*", map[string]string{
				"position.row":                    "0",
				"position.col":                    "0",
				"type":                            "PASSWORD_VERIFY",
				"label":                           "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Placeholder\"}]}]",
				"key":                             "password-field",
				"required":                        "false",
				"attribute_disabled":              "false",
				"other_option_enabled":            "false",
				"other_option_attribute_disabled": "false",
				"show_password_requirements":      "false",
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
			resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "components.fields.*", map[string]string{
				"position.row":                    "0",
				"position.col":                    "0",
				"position.width":                  "50",
				"type":                            "RADIO",
				"label":                           "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"fields.user.locale.label\",\"defaultTranslation\":\"Locale\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]",
				"label_mode":                      "FLOAT",
				"layout":                          "VERTICAL",
				"key":                             "user.locale",
				"required":                        "true",
				"attribute_disabled":              "false",
				"other_option_enabled":            "false",
				"other_option_attribute_disabled": "false",
				"options.0.value":                 "Option1",
				"options.0.label":                 "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Option 1\"}]}]",
				"options.1.value":                 "Option2",
				"options.1.label":                 "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Option 2\"}]}]",
				"options.2.value":                 "Option3",
				"options.2.label":                 "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Option 3\"}]}]",
			}),
		),
	}

	minimalStep := resource.TestStep{
		Config: testAccFormConfig_FieldRadioMinimal(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "components.fields.*", map[string]string{
				"position.row":                    "0",
				"position.col":                    "0",
				"type":                            "RADIO",
				"label":                           "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Placeholder\"}]}]",
				"key":                             "radio-field",
				"layout":                          "HORIZONTAL",
				"required":                        "false",
				"attribute_disabled":              "false",
				"other_option_enabled":            "false",
				"other_option_attribute_disabled": "false",
				"options.0.value":                 "Option1",
				"options.0.label":                 "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Option 1\"}]}]",
				"options.1.value":                 "Option3",
				"options.1.label":                 "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Option 3\"}]}]",
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

func TestAccForm_FieldText(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_form.%s", resourceName)

	name := resourceName

	fullStep := resource.TestStep{
		Config: testAccFormConfig_FieldTextFull(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "components.fields.*", map[string]string{
				"position.row":                    "0",
				"position.col":                    "0",
				"position.width":                  "50",
				"type":                            "TEXT",
				"label":                           "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"fields.user.username.label\",\"defaultTranslation\":\"Username\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]",
				"label_mode":                      "FLOAT",
				"layout":                          "VERTICAL",
				"key":                             "user.username",
				"required":                        "true",
				"attribute_disabled":              "false",
				"validation.type":                 "CUSTOM",
				"validation.regex":                "[a-zA-Z0-9]+",
				"validation.error_message":        "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Must be alphanumeric\"}]}]",
				"other_option_enabled":            "false",
				"other_option_attribute_disabled": "false",
			}),
		),
	}

	minimalStep := resource.TestStep{
		Config: testAccFormConfig_FieldTextMinimal(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "components.fields.*", map[string]string{
				"position.row":                    "0",
				"position.col":                    "0",
				"type":                            "TEXT",
				"label":                           "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Placeholder\"}]}]",
				"key":                             "text-field",
				"required":                        "false",
				"attribute_disabled":              "false",
				"validation.type":                 "NONE",
				"other_option_enabled":            "false",
				"other_option_attribute_disabled": "false",
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

func TestAccForm_ItemDivider(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_form.%s", resourceName)

	name := resourceName

	fullStep := resource.TestStep{
		Config: testAccFormConfig_ItemDividerFull(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "components.fields.*", map[string]string{
				"position.row":   "0",
				"position.col":   "0",
				"position.width": "50",
				"type":           "DIVIDER",
			}),
		),
	}

	minimalStep := resource.TestStep{
		Config: testAccFormConfig_ItemDividerMinimal(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "components.fields.*", map[string]string{
				"position.row": "0",
				"position.col": "0",
				"type":         "DIVIDER",
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
				Config:  testAccFormConfig_ItemDividerFull(resourceName, name),
				Destroy: true,
			},
			// Minimal step
			minimalStep,
			{
				Config:  testAccFormConfig_ItemDividerMinimal(resourceName, name),
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

func TestAccForm_ItemEmptyField(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_form.%s", resourceName)

	name := resourceName

	fullStep := resource.TestStep{
		Config: testAccFormConfig_ItemEmptyFieldFull(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "components.fields.*", map[string]string{
				"position.row":   "0",
				"position.col":   "1",
				"position.width": "50",
				"type":           "EMPTY_FIELD",
			}),
		),
	}

	minimalStep := resource.TestStep{
		Config: testAccFormConfig_ItemEmptyFieldMinimal(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "components.fields.*", map[string]string{
				"position.row": "0",
				"position.col": "1",
				"type":         "EMPTY_FIELD",
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
				Config:  testAccFormConfig_ItemEmptyFieldFull(resourceName, name),
				Destroy: true,
			},
			// Minimal step
			minimalStep,
			{
				Config:  testAccFormConfig_ItemEmptyFieldMinimal(resourceName, name),
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

func TestAccForm_ItemErrorDisplay(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_form.%s", resourceName)

	name := resourceName

	fullStep := resource.TestStep{
		Config: testAccFormConfig_ItemErrorDisplayFull(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "components.fields.*", map[string]string{
				"position.row":   "0",
				"position.col":   "0",
				"position.width": "50",
				"type":           "ERROR_DISPLAY",
			}),
		),
	}

	minimalStep := resource.TestStep{
		Config: testAccFormConfig_ItemErrorDisplayMinimal(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "components.fields.*", map[string]string{
				"position.row": "0",
				"position.col": "0",
				"type":         "ERROR_DISPLAY",
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
				Config:  testAccFormConfig_ItemErrorDisplayFull(resourceName, name),
				Destroy: true,
			},
			// Minimal step
			minimalStep,
			{
				Config:  testAccFormConfig_ItemErrorDisplayMinimal(resourceName, name),
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

func TestAccForm_ItemTextblob(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_form.%s", resourceName)

	name := resourceName

	fullStep := resource.TestStep{
		Config: testAccFormConfig_ItemTextblobFull(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "components.fields.*", map[string]string{
				"position.row":   "0",
				"position.col":   "0",
				"position.width": "50",
				"type":           "TEXTBLOB",
				"content":        "<p>Two baguettes in a zoo cage, the sign says 'Bread in captivity'.</p>",
			}),
		),
	}

	minimalStep := resource.TestStep{
		Config: testAccFormConfig_ItemTextblobMinimal(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "components.fields.*", map[string]string{
				"position.row": "0",
				"position.col": "0",
				"type":         "TEXTBLOB",
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
				Config:  testAccFormConfig_ItemTextblobFull(resourceName, name),
				Destroy: true,
			},
			// Minimal step
			minimalStep,
			{
				Config:  testAccFormConfig_ItemTextblobMinimal(resourceName, name),
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

func TestAccForm_ItemSlateTextblob(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_form.%s", resourceName)

	name := resourceName

	fullStep := resource.TestStep{
		Config: testAccFormConfig_ItemSlateTextblobFull(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "components.fields.*", map[string]string{
				"position.row":   "0",
				"position.col":   "0",
				"position.width": "50",
				"type":           "SLATE_TEXTBLOB",
				"content":        "[{\"children\":[{\"text\":\"Two baguettes in a zoo cage, the sign says 'Bread in captivity'.\"}]}]",
			}),
		),
	}

	// minimalStep := resource.TestStep{
	// 	Config: testAccFormConfig_ItemSlateTextblobMinimal(resourceName, name),
	// 	Check: resource.ComposeTestCheckFunc(
	// 		resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "components.fields.*", map[string]string{
	// 			"position.row": "0",
	// 			"position.col": "0",
	// 			"type":         "SLATE_TEXTBLOB",
	// 		}),
	// 	),
	// }

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
				Config:  testAccFormConfig_ItemSlateTextblobFull(resourceName, name),
				Destroy: true,
			},
			// Minimal step
			// minimalStep,
			// {
			// 	Config:  testAccFormConfig_ItemSlateTextblobMinimal(resourceName, name),
			// 	Destroy: true,
			// },
			// // Change
			// fullStep,
			// minimalStep,
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

  name = "%[4]s"

  mark_required = true
  mark_optional = false

  cols = 4

  components = {
    fields = [
      {
        type = "TEXT"

        position = {
          row = 0
          col = 0
        }

        label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Placeholder\"}]}]"

        key = "text-field"

        validation = {
          type = "NONE"
        }
      },
      {
        type = "SUBMIT_BUTTON"

        position = {
          row = 1
          col = 0
        }

        label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"button.text\",\"defaultTranslation\":\"Submit\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]"
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

  mark_required = true
  mark_optional = true

  cols = 4

  translation_method = "DEFAULT_VALUE"

  components = {
    fields = [
      {
        type = "TEXT"

        position = {
          row = 0
          col = 0
        }

        label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Placeholder\"}]}]"

        key = "text-field"

        validation = {
          type = "NONE"
        }
      },
      {
        type = "SUBMIT_BUTTON"

        position = {
          row = 1
          col = 0
        }

        label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"button.text\",\"defaultTranslation\":\"Submit\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]"
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

  cols = 4

  components = {
    fields = [
      {
        type = "TEXT"

        position = {
          row = 0
          col = 0
        }

        label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Placeholder\"}]}]"

        key = "text-field"

        validation = {
          type = "NONE"
        }
      },
      {
        type = "SUBMIT_BUTTON"

        position = {
          row = 1
          col = 0
        }

        label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"button.text\",\"defaultTranslation\":\"Submit\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]"
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

func testAccFormConfig_FieldCheckboxFull(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_schema_attribute" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  display_name = "%[3]s"

  type        = "STRING"
  unique      = false
  multivalued = true
}

resource "pingone_form" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  mark_required = true
  mark_optional = false

  cols = 4

  components = {
    fields = [
      {
        type = "CHECKBOX"

        position = {
          row   = 0
          col   = 0
          width = 50
        }

        label              = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"fields.user.locale.label\",\"defaultTranslation\":\"Locale\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]"
        label_mode         = "FLOAT"
        layout             = "VERTICAL"
        key                = format("user.%%s", pingone_schema_attribute.%[2]s.name)
        required           = true
        attribute_disabled = false

        options = [
          {
            value = "Option2",
            label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Option 2\"}]}]"
          },
          {
            value = "Option1",
            label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Option 1\"}]}]"
          },
          {
            value = "Option3",
            label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Option 3\"}]}]"
          }
        ]
      },
      {
        type = "SUBMIT_BUTTON"

        position = {
          row = 1
          col = 0
        }

        label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"button.text\",\"defaultTranslation\":\"Submit\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]"
      }
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccFormConfig_FieldCheckboxMinimal(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_schema_attribute" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  display_name = "%[3]s"

  type        = "STRING"
  unique      = false
  multivalued = true
}

resource "pingone_form" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  mark_required = true
  mark_optional = false

  cols = 4

  components = {
    fields = [
      {
        type = "CHECKBOX"

        position = {
          row = 0
          col = 0
        }

        label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Placeholder\"}]}]"

        key = "checkbox-field"

        layout = "HORIZONTAL"

        options = [
          {
            value = "Option1",
            label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Option 1\"}]}]"
          },
          {
            value = "Option3",
            label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Option 3\"}]}]"
          }
        ]
      },
      {
        type = "SUBMIT_BUTTON"

        position = {
          row = 1
          col = 0
        }

        label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"button.text\",\"defaultTranslation\":\"Submit\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]"
      }
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccFormConfig_FieldDropdownFull(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_form" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  mark_required = true
  mark_optional = false

  cols = 4

  components = {
    fields = [
      {
        type = "DROPDOWN"

        position = {
          row   = 0
          col   = 0
          width = 50
        }

        label              = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"fields.user.locale.label\",\"defaultTranslation\":\"Locale\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]"
        label_mode         = "FLOAT"
        layout             = "VERTICAL"
        key                = "user.locale"
        required           = true
        attribute_disabled = false

        options = [
          {
            value = "Option2",
            label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Option 2\"}]}]"
          },
          {
            value = "Option1",
            label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Option 1\"}]}]"
          },
          {
            value = "Option3",
            label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Option 3\"}]}]"
          }
        ]
      },
      {
        type = "SUBMIT_BUTTON"

        position = {
          row = 1
          col = 0
        }

        label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"button.text\",\"defaultTranslation\":\"Submit\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]"
      }
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccFormConfig_FieldDropdownMinimal(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_form" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  mark_required = true
  mark_optional = false

  cols = 4

  components = {
    fields = [
      {
        type = "DROPDOWN"

        position = {
          row = 0
          col = 0
        }

        label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Placeholder\"}]}]"

        key = "dropdown-field"

        options = [
          {
            value = "Option1",
            label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Option 1\"}]}]"
          },
          {
            value = "Option3",
            label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Option 3\"}]}]"
          }
        ]
      },
      {
        type = "SUBMIT_BUTTON"

        position = {
          row = 1
          col = 0
        }

        label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"button.text\",\"defaultTranslation\":\"Submit\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]"
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

  cols = 4

  components = {
    fields = [
      {
        type = "PASSWORD"

        position = {
          row   = 0
          col   = 0
          width = 50
        }

        label                      = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"fields.user.password.label\",\"defaultTranslation\":\"Password\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]"
        label_mode                 = "FLOAT"
        layout                     = "VERTICAL"
        key                        = "user.password"
        required                   = true
        attribute_disabled         = false
        show_password_requirements = true
      },
      {
        type = "SUBMIT_BUTTON"

        position = {
          row = 1
          col = 0
        }

        label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"button.text\",\"defaultTranslation\":\"Submit\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]"
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

  cols = 4

  components = {
    fields = [
      {
        type = "PASSWORD"

        position = {
          row = 0
          col = 0
        }

        label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Placeholder\"}]}]"

        key = "password-field"
      },
      {
        type = "SUBMIT_BUTTON"

        position = {
          row = 1
          col = 0
        }

        label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"button.text\",\"defaultTranslation\":\"Submit\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]"
      }
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccFormConfig_FieldPasswordVerifyFull(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_form" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  mark_required = true
  mark_optional = false

  cols = 4

  components = {
    fields = [
      {
        type = "PASSWORD_VERIFY"

        position = {
          row   = 0
          col   = 0
          width = 50
        }

        label                      = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"fields.user.password.label\",\"defaultTranslation\":\"Password\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]"
        label_mode                 = "FLOAT"
        label_password_verify      = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"fields.user.password.labelPasswordVerify\",\"defaultTranslation\":\"Verify Password\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]"
        layout                     = "VERTICAL"
        key                        = "user.password"
        required                   = true
        attribute_disabled         = false
        show_password_requirements = true
      },
      {
        type = "SUBMIT_BUTTON"

        position = {
          row = 1
          col = 0
        }

        label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"button.text\",\"defaultTranslation\":\"Submit\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]"
      }
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccFormConfig_FieldPasswordVerifyMinimal(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_form" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  mark_required = true
  mark_optional = false

  cols = 4

  components = {
    fields = [
      {
        type = "PASSWORD_VERIFY"

        position = {
          row = 0
          col = 0
        }

        label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Placeholder\"}]}]"

        key = "password-field"
      },
      {
        type = "SUBMIT_BUTTON"

        position = {
          row = 1
          col = 0
        }

        label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"button.text\",\"defaultTranslation\":\"Submit\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]"
      }
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccFormConfig_FieldRadioFull(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_form" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  mark_required = true
  mark_optional = false

  cols = 4

  components = {
    fields = [
      {
        type = "RADIO"

        position = {
          row   = 0
          col   = 0
          width = 50
        }

        label              = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"fields.user.locale.label\",\"defaultTranslation\":\"Locale\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]"
        label_mode         = "FLOAT"
        layout             = "VERTICAL"
        key                = "user.locale"
        required           = true
        attribute_disabled = false

        options = [
          {
            value = "Option2",
            label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Option 2\"}]}]"
          },
          {
            value = "Option1",
            label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Option 1\"}]}]"
          },
          {
            value = "Option3",
            label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Option 3\"}]}]"
          }
        ]
      },
      {
        type = "SUBMIT_BUTTON"

        position = {
          row = 1
          col = 0
        }

        label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"button.text\",\"defaultTranslation\":\"Submit\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]"
      }
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccFormConfig_FieldRadioMinimal(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_form" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  mark_required = true
  mark_optional = false

  cols = 4

  components = {
    fields = [
      {
        type = "RADIO"

        position = {
          row = 0
          col = 0
        }

        label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Placeholder\"}]}]"

        key = "radio-field"

        layout = "HORIZONTAL"

        options = [
          {
            value = "Option1",
            label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Option 1\"}]}]"
          },
          {
            value = "Option3",
            label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Option 3\"}]}]"
          }
        ]
      },
      {
        type = "SUBMIT_BUTTON"

        position = {
          row = 1
          col = 0
        }

        label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"button.text\",\"defaultTranslation\":\"Submit\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]"
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

  mark_required = false
  mark_optional = true

  cols = 4

  components = {
    fields = [
      {
        type = "TEXT"

        position = {
          row   = 0
          col   = 0
          width = 50
        }

        label              = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"fields.user.username.label\",\"defaultTranslation\":\"Username\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]"
        label_mode         = "FLOAT"
        layout             = "VERTICAL"
        key                = "user.username"
        required           = true
        attribute_disabled = false
        validation = {
          type          = "CUSTOM"
          regex         = "[a-zA-Z0-9]+"
          error_message = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Must be alphanumeric\"}]}]"
        }
      },
      {
        type = "SUBMIT_BUTTON"

        position = {
          row = 1
          col = 0
        }

        label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"button.text\",\"defaultTranslation\":\"Submit\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]"
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

  cols = 4

  components = {
    fields = [
      {
        type = "TEXT"

        position = {
          row = 0
          col = 0
        }

        label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Placeholder\"}]}]"

        key = "text-field"

        validation = {
          type = "NONE"
        }
      },
      {
        type = "SUBMIT_BUTTON"

        position = {
          row = 1
          col = 0
        }

        label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"button.text\",\"defaultTranslation\":\"Submit\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]"
      }
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccFormConfig_ItemDividerFull(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_form" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  mark_required = true
  mark_optional = false

  cols = 4

  components = {
    fields = [
      {
        type = "DIVIDER"

        position = {
          row   = 0
          col   = 0
          width = 50
        }
      },
      {
        type = "SUBMIT_BUTTON"

        position = {
          row = 1
          col = 0
        }

        label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"button.text\",\"defaultTranslation\":\"Submit\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]"
      }
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccFormConfig_ItemDividerMinimal(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_form" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  mark_required = true
  mark_optional = false

  cols = 4

  components = {
    fields = [
      {
        type = "DIVIDER"

        position = {
          row = 0
          col = 0
        }
      },
      {
        type = "SUBMIT_BUTTON"

        position = {
          row = 1
          col = 0
        }

        label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"button.text\",\"defaultTranslation\":\"Submit\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]"
      }
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccFormConfig_ItemEmptyFieldFull(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_form" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  mark_required = true
  mark_optional = false

  cols = 4

  components = {
    fields = [
      {
        type = "TEXT"

        position = {
          row   = 0
          col   = 0
          width = 50
        }

        label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Placeholder\"}]}]"

        key = "text-field"

        validation = {
          type = "NONE"
        }
      },
      {
        type = "EMPTY_FIELD"

        position = {
          row   = 0
          col   = 1
          width = 50
        }
      },
      {
        type = "SUBMIT_BUTTON"

        position = {
          row = 1
          col = 0
        }

        label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"button.text\",\"defaultTranslation\":\"Submit\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]"
      }
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccFormConfig_ItemEmptyFieldMinimal(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_form" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  mark_required = true
  mark_optional = false

  cols = 4

  components = {
    fields = [
      {
        type = "TEXT"

        position = {
          row   = 0
          col   = 0
          width = 50
        }

        label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Placeholder\"}]}]"

        key = "text-field"

        validation = {
          type = "NONE"
        }
      },
      {
        type = "EMPTY_FIELD"

        position = {
          row = 0
          col = 1
        }
      },
      {
        type = "SUBMIT_BUTTON"

        position = {
          row = 1
          col = 0
        }

        label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"button.text\",\"defaultTranslation\":\"Submit\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]"
      }
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccFormConfig_ItemErrorDisplayFull(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_form" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  mark_required = true
  mark_optional = false

  cols = 4

  components = {
    fields = [
      {
        type = "ERROR_DISPLAY"

        position = {
          row   = 0
          col   = 0
          width = 50
        }
      },
      {
        type = "SUBMIT_BUTTON"

        position = {
          row = 1
          col = 0
        }

        label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"button.text\",\"defaultTranslation\":\"Submit\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]"
      }
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccFormConfig_ItemErrorDisplayMinimal(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_form" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  mark_required = true
  mark_optional = false

  cols = 4

  components = {
    fields = [
      {
        type = "ERROR_DISPLAY"

        position = {
          row = 0
          col = 0
        }
      },
      {
        type = "SUBMIT_BUTTON"

        position = {
          row = 1
          col = 0
        }

        label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"button.text\",\"defaultTranslation\":\"Submit\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]"
      }
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccFormConfig_ItemTextblobFull(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_form" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  mark_required = true
  mark_optional = false

  cols = 4

  components = {
    fields = [
      {
        type = "TEXTBLOB"

        position = {
          row   = 0
          col   = 0
          width = 50
        }

        content = "<p>Two baguettes in a zoo cage, the sign says 'Bread in captivity'.</p>"
      },
      {
        type = "SUBMIT_BUTTON"

        position = {
          row = 1
          col = 0
        }

        label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"button.text\",\"defaultTranslation\":\"Submit\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]"
      }
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccFormConfig_ItemTextblobMinimal(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_form" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  mark_required = true
  mark_optional = false

  cols = 4

  components = {
    fields = [
      {
        type = "TEXTBLOB"

        position = {
          row = 0
          col = 0
        }
      },
      {
        type = "SUBMIT_BUTTON"

        position = {
          row = 1
          col = 0
        }

        label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"button.text\",\"defaultTranslation\":\"Submit\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]"
      }
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccFormConfig_ItemSlateTextblobFull(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_form" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  mark_required = true
  mark_optional = false

  cols = 4

  components = {
    fields = [
      {
        type = "SLATE_TEXTBLOB"

        position = {
          row   = 0
          col   = 0
          width = 50
        }

        content = "[{\"children\":[{\"text\":\"Two baguettes in a zoo cage, the sign says 'Bread in captivity'.\"}]}]"
      },
      {
        type = "SUBMIT_BUTTON"

        position = {
          row = 1
          col = 0
        }

        label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"button.text\",\"defaultTranslation\":\"Submit\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]"
      }
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccFormConfig_ItemSlateTextblobMinimal(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_form" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  mark_required = true
  mark_optional = false

  cols = 4

  components = {
    fields = [
      {
        type = "SLATE_TEXTBLOB"

        position = {
          row = 0
          col = 0
        }
      },
      {
        type = "SUBMIT_BUTTON"

        position = {
          row = 1
          col = 0
        }

        label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"button.text\",\"defaultTranslation\":\"Submit\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]"
      }
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
