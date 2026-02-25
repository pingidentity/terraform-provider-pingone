// Copyright © 2026 Ping Identity Corporation

package verify

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/patrickcping/pingone-go-sdk-v2/verify"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/legacysdk"
)

func VerifyVoicePhraseContents_CheckDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := legacysdk.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.VerifyAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_verify_voice_phrase_content" {
			continue
		}

		shouldContinue, err := legacysdk.CheckParentEnvironmentDestroy(ctx, p1Client.API.ManagementAPIClient, rs.Primary.Attributes["environment_id"])
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		_, r, err := apiClient.VoicePhraseContentsApi.ReadOneVoicePhraseContent(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["voice_phrase_id"], rs.Primary.Attributes["id"]).Execute()

		shouldContinue, err = acctest.CheckForResourceDestroy(r, err)
		if err != nil {
			return err
		}

		if shouldContinue {
			continue
		}

		return fmt.Errorf("PingOne Voice Phrase Content %s still exists", rs.Primary.ID)
	}

	return nil
}

func VerifyVoicePhraseContent_GetIDs(resourceName string, environmentID, voicePhraseID, resourceID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		if resourceID != nil {
			*resourceID = rs.Primary.ID
		}

		if voicePhraseID != nil {
			*voicePhraseID = rs.Primary.Attributes["voice_phrase_id"]
		}

		if environmentID != nil {
			*environmentID = rs.Primary.Attributes["environment_id"]
		}

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
