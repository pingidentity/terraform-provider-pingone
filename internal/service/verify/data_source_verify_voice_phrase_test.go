package verify_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/verify"
	validation "github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccVerifyVoicePhraseDataSource_All(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("data.pingone_verify_voice_phrase.%s", resourceName)

	name := acctest.ResourceNameGen()
	updatedName := acctest.ResourceNameGen()

	findByID := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", validation.P1ResourceIDRegexp),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", validation.P1ResourceIDRegexp),
		resource.TestCheckResourceAttr(resourceFullName, "display_name", name),
		resource.TestMatchResourceAttr(resourceFullName, "created_at", validation.RFC3339Regexp),
		resource.TestMatchResourceAttr(resourceFullName, "updated_at", validation.RFC3339Regexp),
	)

	findByName := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", validation.P1ResourceIDRegexp),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", validation.P1ResourceIDRegexp),
		resource.TestCheckResourceAttr(resourceFullName, "display_name", updatedName),
		resource.TestMatchResourceAttr(resourceFullName, "created_at", validation.RFC3339Regexp),
		resource.TestMatchResourceAttr(resourceFullName, "updated_at", validation.RFC3339Regexp),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             verify.TestAccCheckVerifyVoicePhraseDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVerifyVoicePhrase_FindByID(resourceName, name),
				Check:  findByID,
			},
			{
				Config:  testAccVerifyVoicePhrase_FindByID(resourceName, name),
				Destroy: true,
			},
			{
				Config: testAccVerifyVoicePhrase_FindByName(resourceName, updatedName),
				Check:  findByName,
			},
			{
				Config:  testAccVerifyVoicePhrase_FindByName(resourceName, updatedName),
				Destroy: true,
			},
		},
	})
}

func TestAccVerifyVoicePhraseDataSource_FailureChecks(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	name := acctest.ResourceNameGen()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             verify.TestAccCheckVerifyVoicePhraseDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccVerifyVoicePhrase_FindByIDFail(resourceName, name),
				ExpectError: regexp.MustCompile("Error: Error when calling `ReadOneVoicePhrase`: voicePhrase could not be found"),
			},
			{
				Config:      testAccVerifyVoicePhrase_FindByNameFail(resourceName, name),
				ExpectError: regexp.MustCompile("Error: Cannot find voice phrase from display name"),
			},
		},
	})
}

func testAccVerifyVoicePhrase_FindByID(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_verify_voice_phrase" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  display_name   = "%[3]s"
}

data "pingone_verify_voice_phrase" "%[2]s" {
  environment_id  = data.pingone_environment.general_test.id
  voice_phrase_id = pingone_verify_voice_phrase.%[2]s.id

  depends_on = [pingone_verify_voice_phrase.%[2]s]
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccVerifyVoicePhrase_FindByName(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_verify_voice_phrase" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  display_name   = "%[3]s"
}

data "pingone_verify_voice_phrase" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  display_name   = "%[3]s"

  depends_on = [pingone_verify_voice_phrase.%[2]s]
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccVerifyVoicePhrase_FindByIDFail(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_verify_voice_phrase" "%[3]s" {
  environment_id  = data.pingone_environment.general_test.id
  voice_phrase_id = "9c052a8a-14be-44e4-8f07-2662569994ce" // dummy ID that conforms to UUID v4


}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccVerifyVoicePhrase_FindByNameFail(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_verify_voice_phrase" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  display_name   = "%[3]s"

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
