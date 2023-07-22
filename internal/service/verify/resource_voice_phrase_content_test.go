package verify_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	validation "github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func testAccCheckVoicePhraseContentsDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.VerifyAPIClient

	mgmtApiClient := p1Client.API.ManagementAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_voice_phrase_content" {
			continue
		}

		_, rEnv, err := mgmtApiClient.EnvironmentsApi.ReadOneEnvironment(ctx, rs.Primary.Attributes["environment_id"]).Execute()

		if err != nil {

			if rEnv == nil {
				return fmt.Errorf("Response object does not exist and no error detected")
			}

			if rEnv.StatusCode == 404 {
				continue
			}

			return err
		}

		body, r, err := apiClient.VoicePhraseContentsApi.ReadOneVoicePhraseContent(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["voice_phrase_id"], rs.Primary.Attributes["id"]).Execute()

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

		return fmt.Errorf("PingOne Voice Phrase Content %s still exists", rs.Primary.ID)
	}

	return nil
}

func TestAccVoicePhraseContent_NewEnv(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_voice_phrase_content.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckVoicePhraseContentsDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVoicePhraseContentConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", validation.P1ResourceIDRegexp),
				),
			},
		},
	})
}

func TestAccVoicePhraseContent_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_voice_phrase_content.%s", resourceName)

	name := acctest.ResourceNameGen()
	updatedName := acctest.ResourceNameGen()

	environmentName := acctest.ResourceNameGenEnvironment()
	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	initialVoicePhraseContent := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", validation.P1ResourceIDRegexp),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", validation.P1ResourceIDRegexp),
		resource.TestMatchResourceAttr(resourceFullName, "voice_phrase_id", validation.P1ResourceIDRegexp),
		resource.TestCheckResourceAttr(resourceFullName, "locale", "en"),
		resource.TestCheckResourceAttr(resourceFullName, "content", "Watch your thoughts; they become words. Watch your words; they become actions. "+
			"Watch your actions; they become habits. Watch your habits; they become character. Watch your character; it becomes your destiny."),
		resource.TestMatchResourceAttr(resourceFullName, "created_at", validation.RFC3339Regexp),
		resource.TestMatchResourceAttr(resourceFullName, "updated_at", validation.RFC3339Regexp),
	)

	updatedVoicePhraseContent := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", validation.P1ResourceIDRegexp),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", validation.P1ResourceIDRegexp),
		resource.TestCheckResourceAttr(resourceFullName, "locale", "en"),
		resource.TestCheckResourceAttr(resourceFullName, "content", "Don't underestimate the importance you can have because history has shown us that "+
			"courage can be contagious and hope can take on a life of its own."),
		resource.TestMatchResourceAttr(resourceFullName, "created_at", validation.RFC3339Regexp),
		resource.TestMatchResourceAttr(resourceFullName, "updated_at", validation.RFC3339Regexp),
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckVoicePhraseDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVoicePhraseContent_Initial(environmentName, licenseID, resourceName, name),
				Check:  initialVoicePhraseContent,
			},
			{
				Config:  testAccVoicePhraseContent_Initial(environmentName, licenseID, resourceName, name),
				Destroy: true,
			},
			{
				Config: testAccVoicePhraseContent_Update(environmentName, licenseID, resourceName, updatedName),
				Check:  updatedVoicePhraseContent,
			},
			{
				Config:  testAccVoicePhraseContent_Update(environmentName, licenseID, resourceName, updatedName),
				Destroy: true,
			},
			// changes
			{
				Config: testAccVoicePhraseContent_Initial(environmentName, licenseID, resourceName, name),
				Check:  initialVoicePhraseContent,
			},
			{
				Config: testAccVoicePhraseContent_Update(environmentName, licenseID, resourceName, updatedName),
				Check:  updatedVoicePhraseContent,
			},
			{
				Config: testAccVoicePhraseContent_Initial(environmentName, licenseID, resourceName, name),
				Check:  initialVoicePhraseContent,
			},
		},
	})
}

func testAccVoicePhraseContentConfig_NewEnv(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_voice_phrase" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s"
}

resource "pingone_voice_phrase_content" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  voice_phrase_id = pingone_voice_phrase.%[3]s.id
  locale = "en"
  content = "Progress is the attraction that moves humanity."

}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccVoicePhraseContent_Initial(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_voice_phrase" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s"
}

resource "pingone_voice_phrase_content" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  voice_phrase_id = pingone_voice_phrase.%[3]s.id
  locale = "en"
  content = "Watch your thoughts; they become words. Watch your words; they become actions. Watch your actions; they become habits. Watch your habits; they become character. Watch your character; it becomes your destiny."

}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccVoicePhraseContent_Update(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_voice_phrase" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s"
}

resource "pingone_voice_phrase_content" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  voice_phrase_id = pingone_voice_phrase.%[3]s.id
  locale = "en"
  content = "Don't underestimate the importance you can have because history has shown us that courage can be contagious and hope can take on a life of its own."

}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}
