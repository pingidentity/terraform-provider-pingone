package verify_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	validation "github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccVoicePhraseContentDataSource_All(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("data.pingone_voice_phrase_content.%s", resourceName)

	name := acctest.ResourceNameGen()

	environmentName := acctest.ResourceNameGenEnvironment()
	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	findByID := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", validation.P1ResourceIDRegexp),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", validation.P1ResourceIDRegexp),
		resource.TestMatchResourceAttr(resourceFullName, "voice_phrase_id", validation.P1ResourceIDRegexp),
		resource.TestCheckResourceAttr(resourceFullName, "locale", "en"),
		resource.TestCheckResourceAttr(resourceFullName, "content", "Knowing is not enough; we must apply. Being willing is not enough; we must do."),
		resource.TestMatchResourceAttr(resourceFullName, "created_at", validation.RFC3339Regexp),
		resource.TestMatchResourceAttr(resourceFullName, "updated_at", validation.RFC3339Regexp),
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckVoicePhraseContentsDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVoicePhraseContent_FindByID(environmentName, licenseID, resourceName, name),
				Check:  findByID,
			},
			{
				Config:  testAccVoicePhraseContent_FindByID(environmentName, licenseID, resourceName, name),
				Destroy: true,
			},
		},
	})
}

func TestAccVoicePhraseContentDataSource_FailureChecks(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	name := acctest.ResourceNameGen()

	environmentName := acctest.ResourceNameGenEnvironment()
	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckVoicePhraseContentsDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccVoicePhraseContent_FindByIDFail(environmentName, licenseID, resourceName, name),
				ExpectError: regexp.MustCompile("Error: Error when calling `ReadOneVoicePhraseContent`: content could not be found"),
			},
		},
	})
}

func testAccVoicePhraseContent_FindByID(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_voice_phrase" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s"
}

resource "pingone_voice_phrase_content" "%[3]s" {
  environment_id  = pingone_environment.%[2]s.id
  voice_phrase_id = pingone_voice_phrase.%[3]s.id
  locale          = "en"
  content         = "Knowing is not enough; we must apply. Being willing is not enough; we must do."
}

data "pingone_voice_phrase_content" "%[3]s" {
  environment_id          = pingone_environment.%[2]s.id
  voice_phrase_id         = pingone_voice_phrase.%[3]s.id
  voice_phrase_content_id = pingone_voice_phrase_content.%[3]s.id
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccVoicePhraseContent_FindByIDFail(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_voice_phrase" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s"
}

data "pingone_voice_phrase_content" "%[3]s" {
  environment_id          = pingone_environment.%[2]s.id
  voice_phrase_id         = pingone_voice_phrase.%[3]s.id
  voice_phrase_content_id = "9c052a8a-14be-44e4-8f07-2662569994ce" // dummy ID that conforms to UUID v4
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}
