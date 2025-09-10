// Copyright © 2025 Ping Identity Corporation

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
	"github.com/pingidentity/terraform-provider-pingone/internal/client"
	validation "github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccVerifyVoicePhrase_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_verify_voice_phrase.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	var voicePhraseID, environmentID string

	var p1Client *client.Client
	var ctx = context.Background()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoBeta(t)
			p1Client = acctest.PreCheckTestClient(ctx, t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             verify.VerifyVoicePhrase_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccVerifyVoicePhrase_Full(resourceName, name),
				Check:  verify.VerifyVoicePhrase_GetIDs(resourceFullName, &environmentID, &voicePhraseID),
			},
			{
				PreConfig: func() {
					verify.VerifyVoicePhrase_RemovalDrift_PreConfig(ctx, p1Client.API.VerifyAPIClient, t, environmentID, voicePhraseID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the environment
			{
				Config: testAccVerifyVoicePhraseConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check:  verify.VerifyVoicePhrase_GetIDs(resourceFullName, &environmentID, &voicePhraseID),
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

func TestAccVerifyVoicePhrase_NewEnv(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_verify_voice_phrase.%s", resourceName)

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
		CheckDestroy:             verify.VerifyVoicePhrase_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVerifyVoicePhraseConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", validation.P1ResourceIDRegexp),
				),
			},
		},
	})
}

func TestAccVerifyVoicePhrase_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_verify_voice_phrase.%s", resourceName)

	name := acctest.ResourceNameGen()
	updatedName := acctest.ResourceNameGen()

	initialVoicePhrase := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", validation.P1ResourceIDRegexp),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", validation.P1ResourceIDRegexp),
		resource.TestCheckResourceAttr(resourceFullName, "display_name", name),
		resource.TestMatchResourceAttr(resourceFullName, "created_at", validation.RFC3339Regexp),
		resource.TestMatchResourceAttr(resourceFullName, "updated_at", validation.RFC3339Regexp),
	)

	updatedVoicePhrase := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", validation.P1ResourceIDRegexp),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", validation.P1ResourceIDRegexp),
		resource.TestCheckResourceAttr(resourceFullName, "display_name", updatedName),
		resource.TestMatchResourceAttr(resourceFullName, "created_at", validation.RFC3339Regexp),
		resource.TestMatchResourceAttr(resourceFullName, "updated_at", validation.RFC3339Regexp),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             verify.VerifyVoicePhrase_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVerifyVoicePhrase_Full(resourceName, name),
				Check:  initialVoicePhrase,
			},
			{
				Config:  testAccVerifyVoicePhrase_Full(resourceName, name),
				Destroy: true,
			},
			{
				Config: testAccVerifyVoicePhrase_Full(resourceName, updatedName),
				Check:  updatedVoicePhrase,
			},
			{
				Config:  testAccVerifyVoicePhrase_Full(resourceName, updatedName),
				Destroy: true,
			},
			// changes
			{
				Config: testAccVerifyVoicePhrase_Full(resourceName, name),
				Check:  initialVoicePhrase,
			},
			{
				Config: testAccVerifyVoicePhrase_Full(resourceName, updatedName),
				Check:  updatedVoicePhrase,
			},
			{
				Config: testAccVerifyVoicePhrase_Full(resourceName, name),
				Check:  initialVoicePhrase,
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

						return fmt.Sprintf("%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccVerifyVoicePhrase_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_verify_voice_phrase.%s", resourceName)

	name := acctest.ResourceNameGen()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             verify.VerifyVoicePhrase_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccVerifyVoicePhrase_Full(resourceName, name),
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

func testAccVerifyVoicePhraseConfig_NewEnv(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_verify_voice_phrase" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  display_name   = "%[4]s"

  depends_on = [pingone_environment.%[2]s]
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccVerifyVoicePhrase_Full(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_verify_voice_phrase" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  display_name   = "%[3]s"

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
