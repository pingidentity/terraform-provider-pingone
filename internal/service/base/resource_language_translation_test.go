// Copyright © 2025 Ping Identity Corporation
// Code generated by ping-terraform-plugin-framework-generator

package base_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/base"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
)

func TestAccLanguageTranslation_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_language_translation.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")
	var environmentId string
	var locale string

	var p1Client *client.Client
	var ctx = context.Background()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)

			p1Client = acctest.PreCheckTestClient(ctx, t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Test removal of the environment
			{
				Config: languageTranslation_NewEnvHCL(environmentName, licenseID, resourceName),
				Check:  languageTranslation_GetIDs(resourceFullName, &environmentId, &locale),
			},
			{
				PreConfig: func() {
					base.Environment_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentId)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccLanguageTranslation_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_language_translation.%s", resourceName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ErrorCheck:               acctest.ErrorCheck(t),
		CheckDestroy:             languageTranslation_CheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: languageTranslation_UpdatedHCL(resourceName),
				Check:  languageTranslation_CheckUpdatedComputedValues(resourceName),
			},
			{
				Config:  languageTranslation_UpdatedHCL(resourceName),
				Destroy: true,
			},
			{
				Config: languageTranslation_InitialHCL(resourceName),
				Check:  languageTranslation_CheckInitialComputedValues(resourceName),
			},
			{
				Config:  languageTranslation_InitialHCL(resourceName),
				Destroy: true,
			},
			{
				Config: languageTranslation_UpdatedHCL(resourceName),
				Check:  languageTranslation_CheckUpdatedComputedValues(resourceName),
			},
			{
				Config: languageTranslation_InitialHCL(resourceName),
				Check:  languageTranslation_CheckInitialComputedValues(resourceName),
			},
			{
				Config: languageTranslation_UpdatedHCL(resourceName),
				Check:  languageTranslation_CheckUpdatedComputedValues(resourceName),
			},
			{
				// Test importing the resource
				ResourceName: fmt.Sprintf("pingone_language_translation.%s", resourceName),
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["locale"]), nil
					}
				}(),
				ImportStateVerifyIdentifierAttribute: "locale",
				ImportState:                          true,
			},
		},
	})
}

func TestAccLanguageTranslation_NewEnvExpectedTranslationCount(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_language_translation.%s", resourceName)
	environmentName := acctest.ResourceNameGenEnvironment()
	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	var environmentId string
	var locale string
	var totalTranslations int

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: languageTranslation_NewEnvHCL(environmentName, licenseID, resourceName),
				Check: resource.ComposeTestCheckFunc(
					languageTranslation_GetIDs(resourceFullName, &environmentId, &locale),
				),
			},
			{
				PreConfig: func() {
					// Fetch all translations for the `en` locale
					p1Client, err := acctest.TestClient(context.Background())
					if err != nil {
						t.Fatalf("Failed to create PingOne client: %v", err)
					}

					apiClient := p1Client.API.ManagementAPIClient
					pagedIterator := apiClient.TranslationsApi.ReadTranslations(context.Background(), environmentId, "en").Execute()

					totalTranslations = 0
					for pageCursor, err := range pagedIterator {
						if err != nil {
							t.Fatalf("Failed to fetch translations: %v", err)
						}

						if translations, ok := pageCursor.EntityArray.Embedded.GetTranslationsOk(); ok {
							totalTranslations += len(translations)
						}
					}
				},
				RefreshState: true,
				Check: func(s *terraform.State) error {
					// Validate the total number of translations
					expectedTranslations := 875
					if totalTranslations != expectedTranslations {
						return fmt.Errorf("Expected %d translations, but got %d", expectedTranslations, totalTranslations)
					}
					return nil
				},
			},
		},
	})
}
func TestAccLanguageTranslation_NewEnv(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: languageTranslation_NewEnvHCL(environmentName, licenseID, resourceName),
				Check:  languageTranslation_CheckInitialComputedValues(resourceName),
			},
		},
	})
}

func TestAccLanguageTranslation_RemoveMiddleEntry(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	environmentName := acctest.ResourceNameGenEnvironment()
	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	initialKeys := []string{
		"flow-ui.button.cancel",
		"flow-ui.button.continue",
		"flow-ui.button.confirm",
	}

	expectedKeysAfterRemoval := []string{
		"flow-ui.button.cancel",
		"flow-ui.button.confirm",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Step 1: Create translations with all initial keys
			{
				Config: languageTranslation_ReorderHCL(environmentName, licenseID, resourceName, initialKeys),
				Check: resource.ComposeTestCheckFunc(
					languageTranslation_ValidateKeys(resourceName, initialKeys),
				),
			},
			// Step 2: Remove the middle entry ("flow-ui.button.continue")
			{
				Config: languageTranslation_ReorderHCL(environmentName, licenseID, resourceName, expectedKeysAfterRemoval),
				Check: resource.ComposeTestCheckFunc(
					languageTranslation_ValidateKeys(resourceName, expectedKeysAfterRemoval),
					languageTranslation_ValidateKeyDoesNotExist(resourceName, "flow-ui.button.continue"),
				),
			},
		},
	})
}

// Helper function to validate that a specific key does not exist
func languageTranslation_ValidateKeyDoesNotExist(resourceName, removedKey string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[fmt.Sprintf("pingone_language_translation.%s", resourceName)]
		if !ok {
			return fmt.Errorf("Resource not found: %s", resourceName)
		}

		// Check that the removed key is not present in the state
		for i := 0; ; i++ {
			attrKey := fmt.Sprintf("translations.%d.key", i)
			if key, ok := rs.Primary.Attributes[attrKey]; ok {
				if key == removedKey {
					return fmt.Errorf("Key %s still exists in translations", removedKey)
				}
			} else {
				break
			}
		}

		return nil
	}
}

func TestAccLanguageTranslation_ValidateKeys(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	environmentName := acctest.ResourceNameGenEnvironment()
	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	expectedKeys := []string{
		"flow-ui.button.cancel",
		"flow-ui.button.continue",
		"flow-ui.button.confirm",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Step 1: Create translations
			{
				Config: languageTranslation_ReorderHCL(environmentName, licenseID, resourceName, expectedKeys),
				Check: resource.ComposeTestCheckFunc(
					languageTranslation_ValidateKeys(resourceName, expectedKeys),
				),
			},
			// Step 2: Reorder the translations
			{
				Config: languageTranslation_ReorderHCL(environmentName, licenseID, resourceName, []string{
					"flow-ui.button.continue",
					"flow-ui.button.confirm",
					"flow-ui.button.cancel",
				}),
				Check: resource.ComposeTestCheckFunc(
					languageTranslation_ValidateKeys(resourceName, expectedKeys),
				),
			},
		},
	})
}

func languageTranslation_ReorderHCL(environmentName, licenseID, resourceName string, keys []string) string {
	translations := ""
	for _, key := range keys {
		translations += fmt.Sprintf(`
	{
		key             = "%s"
		translated_text = "Translated text for %s"
	},`, key, key)
	}

	return fmt.Sprintf(`
%[1]s

resource "pingone_language_translation" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  locale         = "en"
  translations = [
		%[4]s
  ]
}
`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, translations)
}

// Helper function to validate the presence of keys
func languageTranslation_ValidateKeys(resourceName string, expectedKeys []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[fmt.Sprintf("pingone_language_translation.%s", resourceName)]
		if !ok {
			return fmt.Errorf("Resource not found: %s", resourceName)
		}

		// Collect all keys from the state
		actualKeys := map[string]bool{}
		for i := 0; ; i++ {
			attrKey := fmt.Sprintf("translations.%d.key", i)
			if key, ok := rs.Primary.Attributes[attrKey]; ok {
				actualKeys[key] = true
			} else {
				break
			}
		}

		// Ensure all expected keys are present
		for _, key := range expectedKeys {
			if !actualKeys[key] {
				return fmt.Errorf("Expected key %s not found in translations", key)
			}
		}

		return nil
	}
}

func TestAccLanguageTranslation_RemovalDrift_CustomLocale(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_language_translation.%s", resourceName)
	environmentName := acctest.ResourceNameGenEnvironment()
	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Create a custom locale (sv) and translations
			{
				Config: languageTranslation_CustomLocaleHCLFull(environmentName, licenseID, resourceName, "sv"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "translations.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "translations.0.key", "flow-ui.button.createNewAccount"),
					resource.TestCheckResourceAttr(resourceFullName, "translations.0.translated_text", "Skapa nytt konto"),
				),
			},
			// Destroy the custom locale
			{
				Config:  languageTranslation_CustomLocaleHCL(environmentName, licenseID, resourceName, "sv"),
				Destroy: true,
			},
			{
				Config: languageTranslation_CustomLocaleTranslationsHCL(environmentName, licenseID, resourceName, "sv"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "translations.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "translations.0.key", "flow-ui.button.createNewAccount"),
					resource.TestCheckResourceAttr(resourceFullName, "translations.0.translated_text", "Skapa nytt konto"),
				),
			},
			// Verify the translation still exists
			{
				RefreshState: true,
				Check: resource.ComposeTestCheckFunc(
					languageTranslation_ValidateTranslationStillExists(resourceFullName, "sv"),
				),
			},
		},
	})
}

// Helper function to generate HCL for a custom locale
func languageTranslation_CustomLocaleHCL(environmentName, licenseID, resourceName, locale string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_language" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  locale         = "%[4]s"
}
`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, locale)
}

func languageTranslation_CustomLocaleTranslationsHCL(environmentName, licenseID, resourceName, locale string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_language_translation" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  locale         = "%[4]s"
  translations = [
    {
      key             = "flow-ui.button.createNewAccount"
      translated_text = "Skapa nytt konto"
    }
  ]
}
`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, locale)
}

func languageTranslation_CustomLocaleHCLFull(environmentName, licenseID, resourceName, locale string) string {
	return fmt.Sprintf(`
%[1]s
resource "pingone_language" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  locale         = "%[4]s"
}

resource "pingone_language_translation" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  locale         = pingone_language.%[3]s.locale
  translations = [
    {
      key             = "flow-ui.button.createNewAccount"
      translated_text = "Skapa nytt konto"
    }
  ]
}
`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, locale)
}

// Helper function to validate that the translation still exists after locale removal
func languageTranslation_ValidateTranslationStillExists(resourceName, locale string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource not found: %s", resourceName)
		}

		if rs.Primary.Attributes["locale"] != locale {
			return fmt.Errorf("Expected locale %s to still exist, but it does not", locale)
		}

		if rs.Primary.Attributes["translations.0.translated_text"] != "Skapa nytt konto" {
			return fmt.Errorf("Expected translation to still exist, but it does not")
		}

		return nil
	}
}

// Initial HCL with original required values set
func languageTranslation_InitialHCL(resourceName string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_language" "sv" {
  environment_id = data.pingone_environment.general_test.id
  locale         = "sv"
}

resource "pingone_language_translation" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  locale         = pingone_language.sv.locale
  translations = [
    {
      key = "flow-ui.button.createNewAccount"
    }
  ]
}
`, acctest.GenericSandboxEnvironment(), resourceName)
}

// Updated HCL with values updated to reflect a change
func languageTranslation_UpdatedHCL(resourceName string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_language_translation" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  locale         = "sv"
  translations = [
    {
      key             = "flow-ui.button.createNewAccount"
      translated_text = "Skapa nytt konto"
    },
    {
      key             = "flow-ui.label.email"
      translated_text = "E-post"
    }
  ]
}
`, acctest.GenericSandboxEnvironment(), resourceName)
}

func languageTranslation_NewEnvHCL(environmentName, licenseID, resourceName string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_language_translation" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  locale         = "en"
  translations = [
    {
      key = "flow-ui.button.createNewAccount"
    }
  ]
}
`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName)
}

func languageTranslation_CheckInitialComputedValues(resourceName string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(fmt.Sprintf("pingone_language_translation.%s", resourceName), "translations.#", "1"),
		resource.TestCheckResourceAttr(fmt.Sprintf("pingone_language_translation.%s", resourceName), "translations.0.short_key", "button.createNewAccount"),
		resource.TestCheckResourceAttr(fmt.Sprintf("pingone_language_translation.%s", resourceName), "translations.0.reference_text", "Create new Account"),
	)
}

func languageTranslation_CheckUpdatedComputedValues(resourceName string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(fmt.Sprintf("pingone_language_translation.%s", resourceName), "translations.#", "2"),
		resource.TestCheckResourceAttr(fmt.Sprintf("pingone_language_translation.%s", resourceName), "translations.0.short_key", "button.createNewAccount"),
		resource.TestCheckResourceAttr(fmt.Sprintf("pingone_language_translation.%s", resourceName), "translations.0.reference_text", "Create new Account"),
		resource.TestCheckResourceAttr(fmt.Sprintf("pingone_language_translation.%s", resourceName), "translations.0.translated_text", "Skapa nytt konto"),
		resource.TestCheckResourceAttr(fmt.Sprintf("pingone_language_translation.%s", resourceName), "translations.1.short_key", "label.email"),
		resource.TestCheckResourceAttr(fmt.Sprintf("pingone_language_translation.%s", resourceName), "translations.1.reference_text", "Email"),
		resource.TestCheckResourceAttr(fmt.Sprintf("pingone_language_translation.%s", resourceName), "translations.1.translated_text", "E-post"),
	)
}

func languageTranslation_GetIDs(resourceName string, environmentId, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource Not found: %s", resourceName)
		}
		if environmentId != nil {
			*environmentId = rs.Primary.Attributes["environment_id"]
		}
		if id != nil {
			*id = rs.Primary.Attributes["id"]
		}

		return nil
	}
}

func languageTranslation_CheckDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_language_translation" {
			continue
		}

		shouldContinue, err := acctest.CheckParentEnvironmentDestroy(ctx, apiClient, rs.Primary.Attributes["environment_id"])
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		found := false
		pagedIterator := apiClient.TranslationsApi.ReadTranslations(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["locale"]).Execute()

		for pageCursor, err := range pagedIterator {
			if err != nil {
				_, _, err := framework.CheckEnvironmentExistsOnPermissionsError(ctx, apiClient, rs.Primary.Attributes["environment_id"], nil, pageCursor.HTTPResponse, err)
				return err
			}

			if translations, ok := pageCursor.EntityArray.Embedded.GetTranslationsOk(); ok {
				for _, translation := range translations {
					// verifying the translation text is back to the original value for the `sv` locale
					if v, ok := translation.GetIdOk(); ok && *v == rs.Primary.Attributes["id"] && translation.TranslatedText == "" {
						found = true
						break
					}
				}
			}
		}

		if !found {
			continue
		}

		return fmt.Errorf("PingOne Language Translation %s still exists", rs.Primary.ID)
	}

	return nil
}
