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

func TestAccVerifyVoicePhraseDataSource_All(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("data.pingone_verify_voice_phrase.%s", resourceName)

	name := acctest.ResourceNameGen()
	updatedName := acctest.ResourceNameGen()

	environmentName := acctest.ResourceNameGenEnvironment()
	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	findByID := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", validation.P1ResourceIDRegexp),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", validation.P1ResourceIDRegexp),
		resource.TestCheckResourceAttr(resourceFullName, "name", name),
		resource.TestMatchResourceAttr(resourceFullName, "created_at", validation.RFC3339Regexp),
		resource.TestMatchResourceAttr(resourceFullName, "updated_at", validation.RFC3339Regexp),
	)

	findByName := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", validation.P1ResourceIDRegexp),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", validation.P1ResourceIDRegexp),
		resource.TestCheckResourceAttr(resourceFullName, "name", updatedName),
		resource.TestMatchResourceAttr(resourceFullName, "created_at", validation.RFC3339Regexp),
		resource.TestMatchResourceAttr(resourceFullName, "updated_at", validation.RFC3339Regexp),
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckVerifyVoicePhraseDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVerifyVoicePhrase_FindByID(environmentName, licenseID, resourceName, name),
				Check:  findByID,
			},
			{
				Config:  testAccVerifyVoicePhrase_FindByID(environmentName, licenseID, resourceName, name),
				Destroy: true,
			},
			{
				Config: testAccVerifyVoicePhrase_FindByName(environmentName, licenseID, resourceName, updatedName),
				Check:  findByName,
			},
			{
				Config:  testAccVerifyVoicePhrase_FindByName(environmentName, licenseID, resourceName, updatedName),
				Destroy: true,
			},
		},
	})
}

func TestAccVerifyVoicePhraseDataSource_FailureChecks(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	name := acctest.ResourceNameGen()

	environmentName := acctest.ResourceNameGenEnvironment()
	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckVerifyVoicePhraseDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccVerifyVoicePhrase_FindByIDFail(environmentName, licenseID, resourceName, name),
				ExpectError: regexp.MustCompile("Error: Error when calling `ReadOneVoicePhrase`: voicePhrase could not be found"),
			},
			{
				Config:      testAccVerifyVoicePhrase_FindByNameFail(environmentName, licenseID, resourceName, name),
				ExpectError: regexp.MustCompile("Error: Cannot find voice phrase from name"),
			},
		},
	})
}

func testAccVerifyVoicePhrase_FindByID(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_verify_voice_phrase" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s"
}

data "pingone_verify_voice_phrase" "%[3]s" {
  environment_id  = pingone_environment.%[2]s.id
  voice_phrase_id = pingone_verify_voice_phrase.%[3]s.id

  depends_on = [pingone_verify_voice_phrase.%[3]s]
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccVerifyVoicePhrase_FindByName(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_verify_voice_phrase" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s"
}

data "pingone_verify_voice_phrase" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s"

  depends_on = [pingone_verify_voice_phrase.%[3]s]
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccVerifyVoicePhrase_FindByIDFail(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_verify_voice_phrase" "%[3]s" {
  environment_id  = pingone_environment.%[2]s.id
  voice_phrase_id = "9c052a8a-14be-44e4-8f07-2662569994ce" // dummy ID that conforms to UUID v4


}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccVerifyVoicePhrase_FindByNameFail(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_verify_voice_phrase" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s"

}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}
