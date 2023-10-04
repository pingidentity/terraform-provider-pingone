package verify_test

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
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/verify"
	validation "github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccVerifyVoicePhraseContent_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_verify_voice_phrase_content.%s", resourceName)

	locale := "en"
	phrase := "Experience a better experience."

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	var voicePhraseContentID, voicePhraseID, environmentID string

	var ctx = context.Background()
	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		t.Fatalf("Failed to get API client: %v", err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             verify.VerifyVoicePhraseContents_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Test removal of the resource
			{
				Config: testAccVerifyVoicePhraseContent_Full(resourceName, name, locale, phrase),
				Check:  verify.VerifyVoicePhraseContent_GetIDs(resourceFullName, &environmentID, &voicePhraseID, &voicePhraseContentID),
			},
			// Replan after removal preconfig
			{
				PreConfig: func() {
					verify.VerifyVoicePhraseContent_RemovalDrift_PreConfig(ctx, p1Client.API.VerifyAPIClient, t, environmentID, voicePhraseID, voicePhraseContentID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the voice phrase ID
			{
				Config: testAccVerifyVoicePhraseContent_Full(resourceName, name, locale, phrase),
				Check:  verify.VerifyVoicePhraseContent_GetIDs(resourceFullName, &environmentID, &voicePhraseID, &voicePhraseContentID),
			},
			// Replan after removal preconfig
			{
				PreConfig: func() {
					verify.VerifyVoicePhrase_RemovalDrift_PreConfig(ctx, p1Client.API.VerifyAPIClient, t, environmentID, voicePhraseID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the environment
			{
				Config: testAccVerifyVoicePhraseContentConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check:  verify.VerifyVoicePhraseContent_GetIDs(resourceFullName, &environmentID, &voicePhraseID, &voicePhraseContentID),
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

func TestAccVerifyVoicePhraseContent_NewEnv(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_verify_voice_phrase_content.%s", resourceName)

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
		CheckDestroy:             verify.VerifyVoicePhraseContents_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVerifyVoicePhraseContentConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", validation.P1ResourceIDRegexp),
				),
			},
		},
	})
}

func TestAccVerifyVoicePhraseContent_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_verify_voice_phrase_content.%s", resourceName)

	name := acctest.ResourceNameGen()
	locale := "en"
	phrase := "Watch your thoughts; they become words. Watch your words; they become actions. Watch your actions; " +
		"they become habits. Watch your habits; they become character. Watch your character; it becomes your destiny."

	updatedName := acctest.ResourceNameGen()
	updatedPhrase := "Don't underestimate the importance you can have because history has shown us that courage can " +
		"be contagious and hope can take on a life of its own."

	initialVoicePhraseContent := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", validation.P1ResourceIDRegexp),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", validation.P1ResourceIDRegexp),
		resource.TestMatchResourceAttr(resourceFullName, "voice_phrase_id", validation.P1ResourceIDRegexp),
		resource.TestCheckResourceAttr(resourceFullName, "locale", locale),
		resource.TestCheckResourceAttr(resourceFullName, "content", phrase),
		resource.TestMatchResourceAttr(resourceFullName, "created_at", validation.RFC3339Regexp),
		resource.TestMatchResourceAttr(resourceFullName, "updated_at", validation.RFC3339Regexp),
	)

	updatedVoicePhraseContent := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", validation.P1ResourceIDRegexp),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", validation.P1ResourceIDRegexp),
		resource.TestCheckResourceAttr(resourceFullName, "locale", locale),
		resource.TestCheckResourceAttr(resourceFullName, "content", updatedPhrase),
		resource.TestMatchResourceAttr(resourceFullName, "created_at", validation.RFC3339Regexp),
		resource.TestMatchResourceAttr(resourceFullName, "updated_at", validation.RFC3339Regexp),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             verify.VerifyVoicePhraseContents_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVerifyVoicePhraseContent_Full(resourceName, name, locale, phrase),
				Check:  initialVoicePhraseContent,
			},
			{
				Config:  testAccVerifyVoicePhraseContent_Full(resourceName, name, locale, phrase),
				Destroy: true,
			},
			{
				Config: testAccVerifyVoicePhraseContent_Full(resourceName, updatedName, locale, updatedPhrase),
				Check:  updatedVoicePhraseContent,
			},
			{
				Config:  testAccVerifyVoicePhraseContent_Full(resourceName, updatedName, locale, updatedPhrase),
				Destroy: true,
			},
			// changes
			{
				Config: testAccVerifyVoicePhraseContent_Full(resourceName, name, locale, phrase),
				Check:  initialVoicePhraseContent,
			},
			{
				Config: testAccVerifyVoicePhraseContent_Full(resourceName, updatedName, locale, updatedPhrase),
				Check:  updatedVoicePhraseContent,
			},
			{
				Config: testAccVerifyVoicePhraseContent_UpdateVoicePhraseTestReplace(resourceName, updatedName, locale, updatedPhrase),
				Check:  updatedVoicePhraseContent,
			},
			{
				Config: testAccVerifyVoicePhraseContent_Full(resourceName, name, locale, phrase),
				Check:  initialVoicePhraseContent,
			},
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["voice_phrase_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccVerifyVoicePhraseContent_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_verify_voice_phrase_content.%s", resourceName)

	name := acctest.ResourceNameGen()

	locale := "en"
	phrase := "Watch your thoughts; they become words. Watch your words; they become actions. Watch your actions; " +
		"they become habits. Watch your habits; they become character. Watch your character; it becomes your destiny."

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             verify.VerifyVoicePhraseContents_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccVerifyVoicePhraseContent_Full(resourceName, name, locale, phrase),
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

func testAccVerifyVoicePhraseContentConfig_NewEnv(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_verify_voice_phrase" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  display_name   = "%[4]s"
}

resource "pingone_verify_voice_phrase_content" "%[3]s" {
  environment_id  = pingone_environment.%[2]s.id
  voice_phrase_id = pingone_verify_voice_phrase.%[3]s.id
  locale          = "en"
  content         = "Progress is the attraction that moves humanity."

}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccVerifyVoicePhraseContent_Full(resourceName, name, locale, phrase string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_verify_voice_phrase" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  display_name   = "%[3]s"
}

resource "pingone_verify_voice_phrase_content" "%[2]s" {
  environment_id  = data.pingone_environment.general_test.id
  voice_phrase_id = pingone_verify_voice_phrase.%[2]s.id
  locale          = "%[4]s"
  content         = "%[5]s"

}`, acctest.GenericSandboxEnvironment(), resourceName, name, locale, phrase)
}

func testAccVerifyVoicePhraseContent_UpdateVoicePhraseTestReplace(resourceName, name, locale, phrase string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_verify_voice_phrase" "%[2]s-replace" {
  environment_id = data.pingone_environment.general_test.id
  display_name   = "%[3]s"
}

resource "pingone_verify_voice_phrase_content" "%[2]s" {
  environment_id  = data.pingone_environment.general_test.id
  voice_phrase_id = pingone_verify_voice_phrase.%[2]s-replace.id
  locale          = "%[4]s"
  content         = "%[5]s"

}`, acctest.GenericSandboxEnvironment(), resourceName, name, locale, phrase)
}
