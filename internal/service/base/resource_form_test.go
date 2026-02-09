// Copyright © 2026 Ping Identity Corporation

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
	acctestlegacysdk "github.com/pingidentity/terraform-provider-pingone/internal/acctest/legacysdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/base"
	baselegacysdk "github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/base/legacysdk"
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
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoBeta(t)
			p1Client = acctestlegacysdk.PreCheckTestClient(ctx, t)
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
					baselegacysdk.Environment_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID)
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
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoBeta(t)
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
			resource.TestCheckNoResourceAttr(resourceFullName, "translation_method"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
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
							return "", fmt.Errorf("resource not found: %s", resourceFullName)
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
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.#", "10"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.type", "PASSWORD"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.row", "3"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.col", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.key", "user.password"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.label", "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"fields.user.password.label\",\"defaultTranslation\":\"Password\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.required", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.attribute_disabled", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.show_password_requirements", "false"),
			// "validation.type":            "NONE",
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.1.type", "TEXT"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.1.position.row", "2"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.1.position.col", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.1.key", "user.username"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.1.label", "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"fields.user.username.label\",\"defaultTranslation\":\"Enter your email address\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.1.required", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.1.attribute_disabled", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.1.validation.type", "NONE"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.2.type", "SLATE_TEXTBLOB"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.2.position.row", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.2.position.col", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.2.content", "[{\"children\":[{\"text\":\"Sign On\"}],\"type\":\"heading-1\"},{\"type\":\"divider\",\"children\":[{\"text\":\"\"}]},{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"}]}]"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.3.type", "FLOW_LINK"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.3.position.row", "7"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.3.position.col", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.3.key", "issues"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.3.label", "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Issues signing on?\"}]}]"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.4.type", "FLOW_BUTTON"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.4.position.row", "5"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.4.position.col", "1"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.4.key", "passkey"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.4.label", "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Use your Passkey\"}]}]"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.4.styles.alignment", "CENTER"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.4.styles.background_color", "#FFF"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.4.styles.border_color", "#4462ED"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.4.styles.enabled", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.4.styles.height", "36"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.4.styles.padding.bottom", "5"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.4.styles.padding.left", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.4.styles.padding.right", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.4.styles.padding.top", "5"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.4.styles.text_color", "#4462ED"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.4.styles.width", "100"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.4.styles.width_unit", "PERCENT"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.5.type", "FLOW_LINK"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.5.position.row", "7"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.5.position.col", "1"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.5.key", "register"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.5.label", "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Create your account\"}]}]"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.6.type", "SUBMIT_BUTTON"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.6.position.row", "5"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.6.position.col", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.6.label", "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"button.text.signOn\",\"defaultTranslation\":\"Sign On\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.7.type", "ERROR_DISPLAY"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.7.position.row", "1"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.7.position.col", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.8.type", "DIVIDER"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.8.position.row", "4"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.8.position.col", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.9.type", "DIVIDER"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.9.position.row", "6"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.9.position.col", "0"),
		),
	}

	step2 := resource.TestStep{
		Config: testAccFormConfig_MultipleStep2(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.#", "4"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.type", "ERROR_DISPLAY"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.row", "1"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.col", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.1.type", "TEXT"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.1.position.row", "2"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.1.position.col", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.1.key", "user.username"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.1.label", "[{\"children\":[{\"text\":\"\"},{\"children\":[{\"text\":\"\"}],\"defaultTranslation\":\"Username\",\"inline\":true,\"key\":\"fields.user.username.label\",\"type\":\"i18n\"},{\"text\":\"\"}],\"type\":\"paragraph\"}]"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.1.required", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.1.validation.type", "NONE"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.2.type", "SLATE_TEXTBLOB"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.2.position.row", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.2.position.col", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.2.content", "[{\"children\":[{\"text\":\"Sign On\"}],\"type\":\"heading-1\"},{\"children\":[{\"text\":\"\"}],\"type\":\"divider\"},{\"children\":[{\"text\":\"\"}],\"type\":\"paragraph\"}]"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.3.type", "SUBMIT_BUTTON"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.3.position.row", "3"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.3.position.col", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.3.label", "[{\"children\":[{\"text\":\"\"},{\"children\":[{\"text\":\"\"}],\"defaultTranslation\":\"Sign On\",\"inline\":true,\"key\":\"button.text.signOn\",\"type\":\"i18n\"},{\"text\":\"\"}],\"type\":\"paragraph\"}]"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
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
							return "", fmt.Errorf("resource not found: %s", resourceFullName)
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
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.#", "2"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.row", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.col", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.width", "50"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.type", "CHECKBOX"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.label", "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"fields.user.locale.label\",\"defaultTranslation\":\"Locale\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.label_mode", "FLOAT"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.layout", "VERTICAL"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.key", fmt.Sprintf("user.%s", name)),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.required", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.attribute_disabled", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.other_option_enabled", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.other_option_attribute_disabled", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.options.0.value", "Option1"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.options.0.label", "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Option 1\"}]}]"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.options.1.value", "Option2"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.options.1.label", "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Option 2\"}]}]"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.options.2.value", "Option3"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.options.2.label", "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Option 3\"}]}]"),
		),
	}

	minimalStep := resource.TestStep{
		Config: testAccFormConfig_FieldCheckboxMinimal(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.#", "2"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.row", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.col", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.type", "CHECKBOX"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.label", "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Placeholder\"}]}]"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.key", "checkbox-field"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.layout", "HORIZONTAL"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.required", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.attribute_disabled", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.other_option_enabled", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.other_option_attribute_disabled", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.options.0.value", "Option1"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.options.0.label", "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Option 1\"}]}]"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.options.1.value", "Option3"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.options.1.label", "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Option 3\"}]}]"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Form_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Validate
			{
				Config:      testAccFormConfig_FieldCheckboxMissingRequiredParams(resourceName, name),
				ExpectError: regexp.MustCompile(`Invalid DaVinci form configuration`),
			},
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
							return "", fmt.Errorf("resource not found: %s", resourceFullName)
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
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.#", "2"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.row", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.col", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.width", "50"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.type", "COMBOBOX"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.label", "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"fields.user.locale.label\",\"defaultTranslation\":\"Locale\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.label_mode", "FLOAT"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.layout", "VERTICAL"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.key", "user.locale"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.required", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.attribute_disabled", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.other_option_enabled", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.other_option_attribute_disabled", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.options.0.value", "Option1"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.options.0.label", "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Option 1\"}]}]"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.options.1.value", "Option2"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.options.1.label", "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Option 2\"}]}]"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.options.2.value", "Option3"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.options.2.label", "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Option 3\"}]}]"),
		),
	}

	minimalStep := resource.TestStep{
		Config: testAccFormConfig_FieldComboboxMinimal(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.#", "2"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.row", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.col", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.type", "COMBOBOX"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.label", "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Placeholder\"}]}]"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.key", "combobox-field"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.required", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.attribute_disabled", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.other_option_enabled", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.other_option_attribute_disabled", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.options.0.value", "Option1"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.options.0.label", "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Option 1\"}]}]"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.options.1.value", "Option3"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.options.1.label", "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Option 3\"}]}]"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Form_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Validate
			{
				Config:      testAccFormConfig_FieldComboboxMissingRequiredParams(resourceName, name),
				ExpectError: regexp.MustCompile(`Invalid DaVinci form configuration`),
			},
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
							return "", fmt.Errorf("resource not found: %s", resourceFullName)
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
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.#", "2"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.row", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.col", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.width", "50"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.type", "DROPDOWN"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.label", "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"fields.user.locale.label\",\"defaultTranslation\":\"Locale\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.label_mode", "FLOAT"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.layout", "VERTICAL"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.key", "user.locale"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.required", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.attribute_disabled", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.other_option_enabled", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.other_option_attribute_disabled", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.options.0.value", "Option1"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.options.0.label", "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Option 1\"}]}]"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.options.1.value", "Option2"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.options.1.label", "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Option 2\"}]}]"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.options.2.value", "Option3"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.options.2.label", "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Option 3\"}]}]"),
		),
	}

	minimalStep := resource.TestStep{
		Config: testAccFormConfig_FieldDropdownMinimal(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.#", "2"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.row", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.col", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.type", "DROPDOWN"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.label", "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Placeholder\"}]}]"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.key", "dropdown-field"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.required", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.attribute_disabled", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.other_option_enabled", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.other_option_attribute_disabled", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.options.0.value", "Option1"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.options.0.label", "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Option 1\"}]}]"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.options.1.value", "Option3"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.options.1.label", "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Option 3\"}]}]"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Form_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Validate
			{
				Config:      testAccFormConfig_FieldDropdownMissingRequiredParams(resourceName, name),
				ExpectError: regexp.MustCompile(`Invalid DaVinci form configuration`),
			},
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
							return "", fmt.Errorf("resource not found: %s", resourceFullName)
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
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.#", "2"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.row", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.col", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.width", "50"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.type", "PASSWORD"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.label", "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"fields.user.password.label\",\"defaultTranslation\":\"Password\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.label_mode", "FLOAT"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.layout", "VERTICAL"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.key", "user.password"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.required", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.attribute_disabled", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.other_option_enabled", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.other_option_attribute_disabled", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.show_password_requirements", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.validation.type", "NONE"),
		),
	}

	minimalStep := resource.TestStep{
		Config: testAccFormConfig_FieldPasswordMinimal(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.#", "2"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.row", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.col", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.type", "PASSWORD"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.label", "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Placeholder\"}]}]"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.key", "password-field"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.required", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.attribute_disabled", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.other_option_enabled", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.other_option_attribute_disabled", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.show_password_requirements", "false"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Form_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Validate
			{
				Config:      testAccFormConfig_FieldPasswordMissingRequiredParams(resourceName, name),
				ExpectError: regexp.MustCompile(`Invalid DaVinci form configuration`),
			},
			{
				Config:      testAccFormConfig_FieldPasswordInvalidValidation(resourceName, name),
				ExpectError: regexp.MustCompile(`Invalid DaVinci form configuration`),
			},
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
							return "", fmt.Errorf("resource not found: %s", resourceFullName)
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
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.#", "2"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.row", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.col", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.width", "50"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.type", "PASSWORD_VERIFY"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.label", "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"fields.user.password.label\",\"defaultTranslation\":\"Password\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.label_mode", "FLOAT"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.label_password_verify", "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"fields.user.password.labelPasswordVerify\",\"defaultTranslation\":\"Verify Password\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.layout", "VERTICAL"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.key", "user.password"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.required", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.attribute_disabled", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.other_option_enabled", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.other_option_attribute_disabled", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.show_password_requirements", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.validation.type", "NONE"),
		),
	}

	minimalStep := resource.TestStep{
		Config: testAccFormConfig_FieldPasswordVerifyMinimal(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.#", "2"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.row", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.col", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.type", "PASSWORD_VERIFY"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.label", "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Placeholder\"}]}]"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.key", "password-field"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.required", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.attribute_disabled", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.other_option_enabled", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.other_option_attribute_disabled", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.show_password_requirements", "false"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Form_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Validate
			{
				Config:      testAccFormConfig_FieldPasswordVerifyMissingRequiredParams(resourceName, name),
				ExpectError: regexp.MustCompile(`Invalid DaVinci form configuration`),
			},
			{
				Config:      testAccFormConfig_FieldPasswordVerifyInvalidValidation(resourceName, name),
				ExpectError: regexp.MustCompile(`Invalid DaVinci form configuration`),
			},
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
							return "", fmt.Errorf("resource not found: %s", resourceFullName)
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
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.#", "2"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.row", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.col", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.width", "50"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.type", "RADIO"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.label", "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"fields.user.locale.label\",\"defaultTranslation\":\"Locale\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.label_mode", "FLOAT"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.layout", "VERTICAL"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.key", "user.locale"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.required", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.attribute_disabled", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.other_option_enabled", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.other_option_attribute_disabled", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.options.0.value", "Option1"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.options.0.label", "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Option 1\"}]}]"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.options.1.value", "Option2"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.options.1.label", "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Option 2\"}]}]"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.options.2.value", "Option3"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.options.2.label", "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Option 3\"}]}]"),
		),
	}

	minimalStep := resource.TestStep{
		Config: testAccFormConfig_FieldRadioMinimal(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.#", "2"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.row", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.col", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.type", "RADIO"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.label", "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Placeholder\"}]}]"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.key", "radio-field"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.layout", "HORIZONTAL"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.required", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.attribute_disabled", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.other_option_enabled", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.other_option_attribute_disabled", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.options.0.value", "Option1"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.options.0.label", "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Option 1\"}]}]"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.options.1.value", "Option3"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.options.1.label", "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Option 3\"}]}]"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Form_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Validate
			{
				Config:      testAccFormConfig_FieldRadioMissingRequiredParams(resourceName, name),
				ExpectError: regexp.MustCompile(`Invalid DaVinci form configuration`),
			},
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
							return "", fmt.Errorf("resource not found: %s", resourceFullName)
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
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.#", "1"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.row", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.col", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.width", "50"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.type", "SUBMIT_BUTTON"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.label", "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"button.text\",\"defaultTranslation\":\"Submit\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.styles.width", "25"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.styles.width_unit", "PERCENT"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.styles.height", "36"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.styles.padding.top", "10"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.styles.padding.right", "12"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.styles.padding.bottom", "14"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.styles.padding.left", "16"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.styles.alignment", "RIGHT"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.styles.background_color", "#FF0000"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.styles.text_color", "#00FF00"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.styles.border_color", "#0000FF"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.styles.enabled", "true"),
		),
	}

	minimalStep := resource.TestStep{
		Config: testAccFormConfig_FieldSubmitButtonMinimal(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.#", "1"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.row", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.col", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.type", "SUBMIT_BUTTON"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.label", "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Placeholder\"}]}]"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Form_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Validate
			{
				Config:      testAccFormConfig_FieldSubmitButtonMissingRequiredParams(resourceName, name),
				ExpectError: regexp.MustCompile(`Invalid DaVinci form configuration`),
			},
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
							return "", fmt.Errorf("resource not found: %s", resourceFullName)
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
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.#", "2"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.row", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.col", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.width", "50"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.type", "TEXT"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.label", "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"fields.user.username.label\",\"defaultTranslation\":\"Username\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.label_mode", "FLOAT"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.layout", "VERTICAL"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.key", "user.username"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.required", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.attribute_disabled", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.validation.type", "CUSTOM"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.validation.regex", "[a-zA-Z0-9]+"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.validation.error_message", "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Must be alphanumeric\"}]}]"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.other_option_enabled", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.other_option_attribute_disabled", "false"),
		),
	}

	minimalStep := resource.TestStep{
		Config: testAccFormConfig_FieldTextMinimal(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.#", "2"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.row", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.col", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.type", "TEXT"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.label", "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Placeholder\"}]}]"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.key", "text-field"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.required", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.attribute_disabled", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.validation.type", "NONE"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.other_option_enabled", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.other_option_attribute_disabled", "false"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Form_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Validate
			{
				Config:      testAccFormConfig_FieldTextMissingRequiredParams(resourceName, name),
				ExpectError: regexp.MustCompile(`Invalid DaVinci form configuration`),
			},
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
							return "", fmt.Errorf("resource not found: %s", resourceFullName)
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
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.#", "2"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.row", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.col", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.width", "50"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.type", "DIVIDER"),
		),
	}

	minimalStep := resource.TestStep{
		Config: testAccFormConfig_ItemDividerMinimal(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.#", "2"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.row", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.col", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.type", "DIVIDER"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Form_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Validate - not required
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
							return "", fmt.Errorf("resource not found: %s", resourceFullName)
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
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.#", "3"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.1.position.row", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.1.position.col", "1"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.1.position.width", "50"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.1.type", "EMPTY_FIELD"),
		),
	}

	minimalStep := resource.TestStep{
		Config: testAccFormConfig_ItemEmptyFieldMinimal(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.#", "3"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.1.position.row", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.1.position.col", "1"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.1.type", "EMPTY_FIELD"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Form_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Validate - Not required
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
							return "", fmt.Errorf("resource not found: %s", resourceFullName)
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
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.#", "2"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.row", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.col", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.width", "50"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.type", "ERROR_DISPLAY"),
		),
	}

	minimalStep := resource.TestStep{
		Config: testAccFormConfig_ItemErrorDisplayMinimal(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.#", "2"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.row", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.col", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.type", "ERROR_DISPLAY"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Form_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Validate - Not required
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
							return "", fmt.Errorf("resource not found: %s", resourceFullName)
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

func TestAccForm_ItemFlowButton(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_form.%s", resourceName)

	name := resourceName

	fullStep := resource.TestStep{
		Config: testAccFormConfig_ItemFlowButtonFull(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.#", "2"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.row", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.col", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.width", "50"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.type", "FLOW_BUTTON"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.key", "button-field-full"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.label", "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"button.text\",\"defaultTranslation\":\"Submit\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.styles.width", "25"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.styles.width_unit", "PERCENT"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.styles.height", "36"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.styles.padding.top", "10"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.styles.padding.right", "12"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.styles.padding.bottom", "14"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.styles.padding.left", "16"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.styles.alignment", "RIGHT"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.styles.background_color", "#FF0000"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.styles.text_color", "#00FF00"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.styles.border_color", "#0000FF"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.styles.enabled", "true"),
		),
	}

	minimalStep := resource.TestStep{
		Config: testAccFormConfig_ItemFlowButtonMinimal(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.#", "2"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.row", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.col", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.type", "FLOW_BUTTON"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.key", "button-field"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.label", "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Placeholder\"}]}]"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Form_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Validate
			{
				Config:      testAccFormConfig_ItemFlowButtonMissingRequiredParams(resourceName, name),
				ExpectError: regexp.MustCompile(`Invalid DaVinci form configuration`),
			},
			// Full step
			fullStep,
			{
				Config:  testAccFormConfig_ItemFlowButtonFull(resourceName, name),
				Destroy: true,
			},
			// Minimal step
			minimalStep,
			{
				Config:  testAccFormConfig_ItemFlowButtonMinimal(resourceName, name),
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
							return "", fmt.Errorf("resource not found: %s", resourceFullName)
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

func TestAccForm_ItemFlowLink(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_form.%s", resourceName)

	name := resourceName

	fullStep := resource.TestStep{
		Config: testAccFormConfig_ItemFlowLinkFull(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.#", "2"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.row", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.col", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.width", "50"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.type", "FLOW_LINK"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.key", "link-field-full"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.label", "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"button.text\",\"defaultTranslation\":\"Submit\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.styles.padding.top", "10"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.styles.padding.right", "12"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.styles.padding.bottom", "14"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.styles.padding.left", "16"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.styles.alignment", "RIGHT"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.styles.text_color", "#00FF00"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.styles.enabled", "true"),
		),
	}

	minimalStep := resource.TestStep{
		Config: testAccFormConfig_ItemFlowLinkMinimal(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.#", "2"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.row", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.col", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.type", "FLOW_LINK"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.key", "link-field"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.label", "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Placeholder\"}]}]"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Form_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Validate
			{
				Config:      testAccFormConfig_ItemFlowLinkMissingRequiredParams(resourceName, name),
				ExpectError: regexp.MustCompile(`Invalid DaVinci form configuration`),
			},
			// Full step
			fullStep,
			{
				Config:  testAccFormConfig_ItemFlowLinkFull(resourceName, name),
				Destroy: true,
			},
			// Minimal step
			minimalStep,
			{
				Config:  testAccFormConfig_ItemFlowLinkMinimal(resourceName, name),
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
							return "", fmt.Errorf("resource not found: %s", resourceFullName)
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

func TestAccForm_ItemQRCode(t *testing.T) {
	t.Skipf("Skipping test due to QR code functionality not being fully released in PingOne. Test should be re-enabled when QR code functionality is fully released.")
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_form.%s", resourceName)

	name := resourceName

	fullStep := resource.TestStep{
		Config: testAccFormConfig_ItemQRCodeFull(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.#", "2"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.row", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.col", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.width", "50"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.type", "QR_CODE"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.key", "qr-code-field-full"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.qr_code_type", "MFA_AUTH"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.show_border", "true"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.alignment", "RIGHT"),
		),
	}

	minimalStep := resource.TestStep{
		Config: testAccFormConfig_ItemQRCodeMinimal(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.#", "2"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.row", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.col", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.type", "QR_CODE"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.key", "qr-code-field"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.qr_code_type", "MFA_AUTH"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.show_border", "false"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.alignment", "LEFT"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Form_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Validate
			{
				Config:      testAccFormConfig_ItemQRCodeMissingRequiredParams(resourceName, name),
				ExpectError: regexp.MustCompile(`Invalid DaVinci form configuration`),
			},
			// Full step
			fullStep,
			{
				Config:  testAccFormConfig_ItemQRCodeFull(resourceName, name),
				Destroy: true,
			},
			// Minimal step
			minimalStep,
			{
				Config:  testAccFormConfig_ItemQRCodeMinimal(resourceName, name),
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
							return "", fmt.Errorf("resource not found: %s", resourceFullName)
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

func TestAccForm_ItemRecaptchaV2(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_form.%s", resourceName)

	name := resourceName

	fullStep := resource.TestStep{
		Config: testAccFormConfig_ItemRecaptchaV2Full(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.#", "2"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.row", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.col", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.width", "50"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.type", "RECAPTCHA_V2"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.theme", "LIGHT"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.size", "NORMAL"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.alignment", "RIGHT"),
		),
	}

	minimalStep := resource.TestStep{
		Config: testAccFormConfig_ItemRecaptchaV2Minimal(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.#", "2"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.row", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.col", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.type", "RECAPTCHA_V2"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.theme", "DARK"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.size", "COMPACT"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.alignment", "LEFT"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Form_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Validate
			{
				Config:      testAccFormConfig_ItemRecaptchaV2MissingRequiredParams(resourceName, name),
				ExpectError: regexp.MustCompile(`Invalid DaVinci form configuration`),
			},
			// Full step
			fullStep,
			{
				Config:  testAccFormConfig_ItemRecaptchaV2Full(resourceName, name),
				Destroy: true,
			},
			// Minimal step
			minimalStep,
			{
				Config:  testAccFormConfig_ItemRecaptchaV2Minimal(resourceName, name),
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
							return "", fmt.Errorf("resource not found: %s", resourceFullName)
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
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.#", "2"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.row", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.col", "0"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.position.width", "50"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.type", "SLATE_TEXTBLOB"),
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.0.content", "[{\"children\":[{\"text\":\"Two baguettes in a zoo cage, the sign says 'Bread in captivity'.\"}]}]"),
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
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Form_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Validate
			// {
			// 	Config:      testAccFormConfig_ItemSlateTextblobMissingRequiredParams(resourceName, name),
			// 	ExpectError: regexp.MustCompile(`Invalid DaVinci form configuration`),
			// },
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
							return "", fmt.Errorf("resource not found: %s", resourceFullName)
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

// Deprecated start
func TestAccForm_ItemTextblob(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_form.%s", resourceName)

	name := resourceName

	fullStep := resource.TestStep{
		Config: testAccFormConfig_ItemTextblobFull(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.#", "2"),
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
			resource.TestCheckResourceAttr(resourceFullName, "components.fields.#", "2"),
			resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "components.fields.*", map[string]string{
				"position.row":   "0",
				"position.col":   "0",
				"position.width": "",
				"type":           "TEXTBLOB",
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
		CheckDestroy:             base.Form_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Validate - Not required
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
							return "", fmt.Errorf("resource not found: %s", resourceFullName)
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

// Deprecated end

func TestAccForm_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_form.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.Form_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Test validation
			{
				Config:      testAccFormConfig_MultipleStepDuplicatePosition(resourceName, name),
				ExpectError: regexp.MustCompile(`Invalid DaVinci form configuration`),
			},
			{
				Config:      testAccFormConfig_NoSubmitButton(resourceName, name),
				ExpectError: regexp.MustCompile(`Invalid DaVinci form configuration`),
			},
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
}`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
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
      },
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

  cols = 4

  components = {
    fields = [
      {
        attribute_disabled = false
        key                = "user.password"
        label              = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"fields.user.password.label\",\"defaultTranslation\":\"Password\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]"
        position = {
          col = 0
          row = 3
        }
        required                   = true
        show_password_requirements = false
        type                       = "PASSWORD"
      },
      {
        attribute_disabled = false
        key                = "user.username"
        label              = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"fields.user.username.label\",\"defaultTranslation\":\"Enter your email address\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]"
        position = {
          col = 0
          row = 2
        }
        required = true
        type     = "TEXT"
        validation = {
          type = "NONE"
        }
      },
      {
        content = "[{\"children\":[{\"text\":\"Sign On\"}],\"type\":\"heading-1\"},{\"type\":\"divider\",\"children\":[{\"text\":\"\"}]},{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"}]}]"
        position = {
          col = 0
          row = 0
        }
        type = "SLATE_TEXTBLOB"
      },
      {
        key   = "issues"
        label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Issues signing on?\"}]}]"
        position = {
          col = 0
          row = 7
        }
        type = "FLOW_LINK"
      },
      {
        key   = "passkey"
        label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Use your Passkey\"}]}]"
        position = {
          col = 1
          row = 5
        }
        styles = {
          alignment        = "CENTER"
          background_color = "#FFF"
          border_color     = "#4462ED"
          enabled          = true
          height           = 36
          padding = {
            bottom = 5
            left   = 0
            right  = 0
            top    = 5
          }
          text_color = "#4462ED"
          width      = 100
          width_unit = "PERCENT"
        }
        type = "FLOW_BUTTON"
      },
      {
        key   = "register"
        label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Create your account\"}]}]"
        position = {
          col = 1
          row = 7
        }
        type = "FLOW_LINK"
      },
      {
        label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"button.text.signOn\",\"defaultTranslation\":\"Sign On\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]"
        position = {
          col = 0
          row = 5
        }
        type = "SUBMIT_BUTTON"
      },
      {
        position = {
          col = 0
          row = 1
        }
        type = "ERROR_DISPLAY"
      },
      {
        position = {
          col = 0
          row = 4
        }
        type = "DIVIDER"
      },
      {
        position = {
          col = 0
          row = 6
        }
        type = "DIVIDER"
      },
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

  cols = 4

  components = {
    fields = [
      {
        type = "ERROR_DISPLAY"

        position = {
          row = 1
          col = 0
        }
      },
      {
        type = "TEXT"

        position = {
          row = 2
          col = 0
        }

        key = "user.username"
        label = jsonencode(
          [
            {
              "type" = "paragraph",
              "children" = [
                {
                  "text" = ""
                },
                {
                  "type"               = "i18n",
                  "key"                = "fields.user.username.label",
                  "defaultTranslation" = "Username",
                  "inline"             = true,
                  "children" = [
                    {
                      "text" = ""
                    }
                  ]
                },
                {
                  "text" = ""
                }
              ]
            }
          ]
        )

        required = true

        validation = {
          type = "NONE"
        }
      },
      {
        type = "SLATE_TEXTBLOB"

        position = {
          row = 0
          col = 0
        }

        content = jsonencode([{ "children" : [{ "text" : "Sign On" }], "type" : "heading-1" }, { "type" : "divider", "children" : [{ "text" : "" }] }, { "type" : "paragraph", "children" : [{ "text" : "" }] }])
      },
      {
        type = "SUBMIT_BUTTON"

        position = {
          row = 3
          col = 0
        }

        label = jsonencode(
          [
            {
              "type" = "paragraph",
              "children" = [
                {
                  "text" = ""
                },
                {
                  "type"               = "i18n",
                  "key"                = "button.text.signOn",
                  "defaultTranslation" = "Sign On",
                  "inline"             = true,
                  "children" = [
                    {
                      "text" = ""
                    }
                  ]
                },
                {
                  "text" = ""
                }
              ]
            }
          ]
        )
      }
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccFormConfig_MultipleStepDuplicatePosition(resourceName, name string) string {
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
          row = 1
          col = 0
        }
      },
      {
        type = "TEXT"

        position = {
          row = 2
          col = 0
        }

        key = "user.username"
        label = jsonencode(
          [
            {
              "type" = "paragraph",
              "children" = [
                {
                  "text" = ""
                },
                {
                  "type"               = "i18n",
                  "key"                = "fields.user.username.label",
                  "defaultTranslation" = "Username",
                  "inline"             = true,
                  "children" = [
                    {
                      "text" = ""
                    }
                  ]
                },
                {
                  "text" = ""
                }
              ]
            }
          ]
        )

        required = true

        validation = {
          type = "NONE"
        }
      },
      {
        type = "SLATE_TEXTBLOB"

        position = {
          row = 0
          col = 0
        }

        content = jsonencode([{ "children" : [{ "text" : "Sign On" }], "type" : "heading-1" }, { "type" : "divider", "children" : [{ "text" : "" }] }, { "type" : "paragraph", "children" : [{ "text" : "" }] }])
      },
      {
        type = "SUBMIT_BUTTON"

        position = {
          row = 2
          col = 0
        }

        label = jsonencode(
          [
            {
              "type" = "paragraph",
              "children" = [
                {
                  "text" = ""
                },
                {
                  "type"               = "i18n",
                  "key"                = "button.text.signOn",
                  "defaultTranslation" = "Sign On",
                  "inline"             = true,
                  "children" = [
                    {
                      "text" = ""
                    }
                  ]
                },
                {
                  "text" = ""
                }
              ]
            }
          ]
        )
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

func testAccFormConfig_FieldCheckboxMissingRequiredParams(resourceName, name string) string {
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

func testAccFormConfig_FieldComboboxFull(resourceName, name string) string {
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
        type = "COMBOBOX"

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

func testAccFormConfig_FieldComboboxMinimal(resourceName, name string) string {
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
        type = "COMBOBOX"

        position = {
          row = 0
          col = 0
        }

        label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Placeholder\"}]}]"

        key = "combobox-field"

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

func testAccFormConfig_FieldComboboxMissingRequiredParams(resourceName, name string) string {
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
        type = "COMBOBOX"

        position = {
          row = 0
          col = 0
        }

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

func testAccFormConfig_FieldDropdownMissingRequiredParams(resourceName, name string) string {
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

func testAccFormConfig_FieldPasswordMissingRequiredParams(resourceName, name string) string {
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

func testAccFormConfig_FieldPasswordInvalidValidation(resourceName, name string) string {
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

func testAccFormConfig_FieldPasswordVerifyMissingRequiredParams(resourceName, name string) string {
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

func testAccFormConfig_FieldPasswordVerifyInvalidValidation(resourceName, name string) string {
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

func testAccFormConfig_FieldRadioMissingRequiredParams(resourceName, name string) string {
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

func testAccFormConfig_FieldSubmitButtonFull(resourceName, name string) string {
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
        type = "SUBMIT_BUTTON"

        position = {
          row   = 0
          col   = 0
          width = 50
        }

        label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"button.text\",\"defaultTranslation\":\"Submit\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]"

        styles = {
          width      = 25
          width_unit = "PERCENT"
          height     = 36

          padding = {
            top    = 10
            right  = 12
            bottom = 14
            left   = 16
          }

          alignment        = "RIGHT"
          background_color = "#FF0000"
          text_color       = "#00FF00"
          border_color     = "#0000FF"
          enabled          = true
        }
      }
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccFormConfig_FieldSubmitButtonMinimal(resourceName, name string) string {
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
        type = "SUBMIT_BUTTON"

        position = {
          row = 0
          col = 0
        }

        label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Placeholder\"}]}]"
      }
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccFormConfig_FieldSubmitButtonMissingRequiredParams(resourceName, name string) string {
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
        type = "SUBMIT_BUTTON"

        position = {
          row = 0
          col = 0
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

func testAccFormConfig_FieldTextMissingRequiredParams(resourceName, name string) string {
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

func testAccFormConfig_ItemFlowButtonFull(resourceName, name string) string {
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
        type = "FLOW_BUTTON"

        key = "button-field-full"

        position = {
          row   = 0
          col   = 0
          width = 50
        }

        label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"button.text\",\"defaultTranslation\":\"Submit\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]"

        styles = {
          width      = 25
          width_unit = "PERCENT"
          height     = 36

          padding = {
            top    = 10
            right  = 12
            bottom = 14
            left   = 16
          }

          alignment        = "RIGHT"
          background_color = "#FF0000"
          text_color       = "#00FF00"
          border_color     = "#0000FF"
          enabled          = true
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

func testAccFormConfig_ItemFlowButtonMinimal(resourceName, name string) string {
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
        type = "FLOW_BUTTON"

        key = "button-field"

        position = {
          row = 0
          col = 0
        }

        label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Placeholder\"}]}]"
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

func testAccFormConfig_ItemFlowButtonMissingRequiredParams(resourceName, name string) string {
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
        type = "FLOW_BUTTON"

        position = {
          row = 0
          col = 0
        }

        label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Placeholder\"}]}]"
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

func testAccFormConfig_ItemFlowLinkFull(resourceName, name string) string {
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
        type = "FLOW_LINK"

        key = "link-field-full"

        position = {
          row   = 0
          col   = 0
          width = 50
        }

        label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"button.text\",\"defaultTranslation\":\"Submit\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]"

        styles = {
          padding = {
            top    = 10
            right  = 12
            bottom = 14
            left   = 16
          }

          alignment  = "RIGHT"
          text_color = "#00FF00"
          enabled    = true
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

func testAccFormConfig_ItemFlowLinkMinimal(resourceName, name string) string {
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
        type = "FLOW_LINK"

        key = "link-field"

        position = {
          row = 0
          col = 0
        }

        label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Placeholder\"}]}]"
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

func testAccFormConfig_ItemFlowLinkMissingRequiredParams(resourceName, name string) string {
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
        type = "FLOW_LINK"

        position = {
          row = 0
          col = 0
        }

        label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"Placeholder\"}]}]"
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

func testAccFormConfig_ItemQRCodeFull(resourceName, name string) string {
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
        type = "QR_CODE"

        position = {
          row   = 0
          col   = 0
          width = 50
        }

        key          = "qr-code-field-full"
        qr_code_type = "MFA_AUTH"
        alignment    = "RIGHT"
        show_border  = true
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

func testAccFormConfig_ItemQRCodeMinimal(resourceName, name string) string {
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
        type = "QR_CODE"

        position = {
          row = 0
          col = 0
        }

        key          = "qr-code-field"
        qr_code_type = "MFA_AUTH"
        alignment    = "LEFT"
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

func testAccFormConfig_ItemQRCodeMissingRequiredParams(resourceName, name string) string {
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
        type = "QR_CODE"

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

func testAccFormConfig_ItemRecaptchaV2Full(resourceName, name string) string {
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
        type = "RECAPTCHA_V2"

        position = {
          row   = 0
          col   = 0
          width = 50
        }

        theme     = "LIGHT"
        size      = "NORMAL"
        alignment = "RIGHT"
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

func testAccFormConfig_ItemRecaptchaV2Minimal(resourceName, name string) string {
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
        type = "RECAPTCHA_V2"

        position = {
          row = 0
          col = 0
        }

        theme     = "DARK"
        size      = "COMPACT"
        alignment = "LEFT"
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

func testAccFormConfig_ItemRecaptchaV2MissingRequiredParams(resourceName, name string) string {
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
        type = "RECAPTCHA_V2"

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

// func testAccFormConfig_ItemSlateTextblobMinimal(resourceName, name string) string {
// 	return fmt.Sprintf(`
// 	%[1]s

// resource "pingone_form" "%[2]s" {
//   environment_id = data.pingone_environment.general_test.id

//   name = "%[3]s"

//   mark_required = true
//   mark_optional = false

//   cols = 4

//   components = {
//     fields = [
//       {
//         type = "SLATE_TEXTBLOB"

//         position = {
//           row = 0
//           col = 0
//         }
//       },
//       {
//         type = "SUBMIT_BUTTON"

//         position = {
//           row = 1
//           col = 0
//         }

//         label = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"button.text\",\"defaultTranslation\":\"Submit\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]"
//       }
//     ]
//   }
// }`, acctest.GenericSandboxEnvironment(), resourceName, name)
// }

// Deprecated start
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

// Deprecated end

func testAccFormConfig_NoSubmitButton(resourceName, name string) string {
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
        attribute_disabled = false
        key                = "user.username"
        label              = "[{\"type\":\"paragraph\",\"children\":[{\"text\":\"\"},{\"type\":\"i18n\",\"key\":\"fields.user.username.label\",\"defaultTranslation\":\"Enter your email address\",\"inline\":true,\"children\":[{\"text\":\"\"}]},{\"text\":\"\"}]}]"
        position = {
          col = 0
          row = 0
        }
        required = true
        type     = "TEXT"
        validation = {
          type = "NONE"
        }
      },
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
