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

func TestAccCheckVerifyVoicePhraseDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.VerifyAPIClient

	mgmtApiClient := p1Client.API.ManagementAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_verify_voice_phrase" {
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

		body, r, err := apiClient.VoicePhrasesApi.ReadOneVoicePhrase(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["id"]).Execute()

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

		return fmt.Errorf("PingOne Voice Phrase %s still exists", rs.Primary.ID)
	}

	return nil
}

func TestAccGetVerifyVoicePhraseIDs(resourceName string, environmentID, resourceID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource Not found: %s", resourceName)
		}

		*resourceID = rs.Primary.ID
		*environmentID = rs.Primary.Attributes["environment_id"]

		return nil
	}
}

func VerifyVoicePhrase_RemovalDrift_PreConfig(ctx context.Context, apiClient *verify.APIClient, t *testing.T, environmentID, voicePhraseID string) {
	if environmentID == "" || voicePhraseID == "" {
		t.Fatalf("One of environment ID or voice phrase ID cannot be determined. Environment ID: %s, Voice Phrase ID: %s", environmentID, voicePhraseID)
	}

	_, err := apiClient.VoicePhrasesApi.DeleteVoicePhrase(ctx, environmentID, voicePhraseID).Execute()
	if err != nil {
		t.Fatalf("Failed to delete voice phrase: %v", err)
	}
}
