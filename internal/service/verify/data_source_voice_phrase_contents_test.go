package verify_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	validation "github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccVoicePhraseContentsDataSource_NoFilter(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullName := fmt.Sprintf("data.pingone_voice_phrase_contents.%s", resourceName)

	name := acctest.ResourceNameGen()

	environmentName := acctest.ResourceNameGenEnvironment()
	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	findByID := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(dataSourceFullName, "id", validation.P1ResourceIDRegexp),
		resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", validation.P1ResourceIDRegexp),
		resource.TestCheckResourceAttr(dataSourceFullName, "ids.#", "4"),
		resource.TestMatchResourceAttr(dataSourceFullName, "ids.0", validation.P1ResourceIDRegexp),
		resource.TestMatchResourceAttr(dataSourceFullName, "ids.1", validation.P1ResourceIDRegexp),
		resource.TestMatchResourceAttr(dataSourceFullName, "ids.2", validation.P1ResourceIDRegexp),
		resource.TestMatchResourceAttr(dataSourceFullName, "ids.3", validation.P1ResourceIDRegexp),
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckVoicePhraseContentsDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVoicePhraseContents_NoFilter(environmentName, licenseID, resourceName, name),
				Check:  findByID,
			},
			{
				Config:  testAccVoicePhraseContents_NoFilter(environmentName, licenseID, resourceName, name),
				Destroy: true,
			},
		},
	})
}

func testAccVoicePhraseContents_NoFilter(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_voice_phrase" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s"
}

resource "pingone_voice_phrase_content" "%[3]s-1" {
  environment_id  = pingone_environment.%[2]s.id
  voice_phrase_id = pingone_voice_phrase.%[3]s.id
  locale          = "es-MX"
  content         = "Anda a ver si ya parió la marrana."
}

resource "pingone_voice_phrase_content" "%[3]s-2" {
  environment_id  = pingone_environment.%[2]s.id
  voice_phrase_id = pingone_voice_phrase.%[3]s.id
  locale          = "sw-KE"
  content         = "Usijenge uadui na adui."
}

resource "pingone_voice_phrase_content" "%[3]s-3" {
  environment_id  = pingone_environment.%[2]s.id
  voice_phrase_id = pingone_voice_phrase.%[3]s.id
  locale          = "sw"
  content         = "Kila jambo na wakati wake."
}

resource "pingone_voice_phrase_content" "%[3]s-4" {
  environment_id  = pingone_environment.%[2]s.id
  voice_phrase_id = pingone_voice_phrase.%[3]s.id
  locale          = "en"
  content         = "I don't have friends, and it's hard for me to make new friends. Right now, the people that are in my life are the people that I work with."
}

data "pingone_voice_phrase_contents" "%[3]s" {
  environment_id          = pingone_environment.%[2]s.id
  voice_phrase_id         = pingone_voice_phrase.%[3]s.id

  depends_on = [pingone_voice_phrase_content.%[3]s-1, pingone_voice_phrase_content.%[3]s-2, pingone_voice_phrase_content.%[3]s-3, pingone_voice_phrase_content.%[3]s-4 ]
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}
