package verify

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/patrickcping/pingone-go-sdk-v2/verify"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
)

func TestAccCheckVerifyVoicePhraseContentsDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.VerifyAPIClient

	mgmtApiClient := p1Client.API.ManagementAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_verify_voice_phrase_content" {
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

func TestAccGetVerifyVoicePhraseContentIDs(resourceName string, environmentID, voicePhraseID, resourceID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource Not found: %s", resourceName)
		}

		*resourceID = rs.Primary.ID
		*voicePhraseID = rs.Primary.Attributes["voice_phrase_id"]
		*environmentID = rs.Primary.Attributes["environment_id"]

		return nil
	}
}

func VerifyVoicePhraseContent_RemovalDrift_PreConfig(ctx context.Context, apiClient *verify.APIClient, t *testing.T, environmentID, voicePhraseID, voicePhraseContentID string) {
	if environmentID == "" || voicePhraseID == "" || voicePhraseContentID == "" {
		t.Fatalf("One of environment ID, voice phrase ID or voice phrase content ID cannot be determined. Environment ID: %s, Voice Phrase ID: %s, Voice Phrase Content ID: %s", environmentID, voicePhraseID, voicePhraseContentID)
	}

	_, err := apiClient.VoicePhraseContentsApi.DeleteVoicePhraseContent(ctx, environmentID, voicePhraseID, voicePhraseContentID).Execute()
	if err != nil {
		t.Fatalf("Failed to delete voice phrase content: %v", err)
	}
}
