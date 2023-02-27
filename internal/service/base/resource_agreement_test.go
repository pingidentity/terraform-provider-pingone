package base_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
)

func testAccCheckAgreementDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_agreement" {
			continue
		}

		_, rEnv, err := apiClient.EnvironmentsApi.ReadOneEnvironment(ctx, rs.Primary.Attributes["environment_id"]).Execute()

		if err != nil {

			if rEnv == nil {
				return fmt.Errorf("Response object does not exist and no error detected")
			}

			if rEnv.StatusCode == 404 {
				continue
			}

			return err
		}

		body, r, err := apiClient.AgreementsResourcesApi.ReadOneAgreement(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.ID).Execute()

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

		return fmt.Errorf("PingOne agreement %s still exists", rs.Primary.ID)
	}

	return nil
}

func TestAccAgreement_NewEnv(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_agreement.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckAgreementDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAgreementConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
				),
			},
		},
	})
}

func TestAccAgreement_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_agreement.%s", resourceName)

	name := resourceName

	fullCheck := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
		resource.TestCheckResourceAttr(resourceFullName, "name", name),
		resource.TestCheckResourceAttr(resourceFullName, "description", "Before the crowbar was invented, Crows would just drink at home."),
		resource.TestCheckResourceAttr(resourceFullName, "enabled", "false"),
		resource.TestCheckResourceAttr(resourceFullName, "reconsent_period_days", "30"),
		resource.TestCheckResourceAttr(resourceFullName, "localized_text.#", "3"),
		resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "localized_text.*", map[string]string{
			"display_name":                        "British English",
			"locale":                              "en-GB",
			"enabled":                             "true",
			"text_checkbox_accept":                "Yeah",
			"text_button_continue":                "Go on",
			"text_button_decline":                 "nah",
			"latest_revision.0.#":                 "1",
			"latest_revision.0.content_type":      "text/html",
			"latest_revision.0.effective_at":      "2100-01-01T01:01:00.000Z",
			"latest_revision.0.require_reconsent": "false",
			"latest_revision.0.text":              "I started a band called 999MB.  We haven't yet got a gig.",
		}),
		resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "localized_text.*", map[string]string{
			"display_name":                        "English",
			"locale":                              "en",
			"enabled":                             "false",
			"text_checkbox_accept":                "Accept",
			"text_button_continue":                "Continue",
			"text_button_decline":                 "Decline",
			"latest_revision.0.#":                 "1",
			"latest_revision.0.content_type":      "text/plain",
			"latest_revision.0.effective_at":      "2100-01-01T01:01:00.000Z",
			"latest_revision.0.require_reconsent": "false",
			"latest_revision.0.text":              "I started a band called 999MB.  We haven't yet got a gig.",
		}),
		resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "localized_text.*", map[string]string{
			"display_name":                        "French",
			"locale":                              "fr",
			"enabled":                             "true",
			"text_checkbox_accept":                "Accepter",
			"text_button_continue":                "Continuer",
			"text_button_decline":                 "Déclin",
			"latest_revision.0.#":                 "1",
			"latest_revision.0.content_type":      "text/html",
			"latest_revision.0.effective_at":      "2100-01-01T01:01:00.000Z",
			"latest_revision.0.require_reconsent": "true",
			"latest_revision.0.text":              "J'ai monté un groupe appelé 999MB. Nous n'avons pas encore de concert.",
		}),
	)

	minimalCheck := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
		resource.TestCheckResourceAttr(resourceFullName, "name", name),
		resource.TestCheckResourceAttr(resourceFullName, "description", ""),
		resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
		resource.TestCheckResourceAttr(resourceFullName, "reconsent_period_days", ""),
		resource.TestCheckResourceAttr(resourceFullName, "localized_text.#", "1"),
		resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "localized_text.*", map[string]string{
			"display_name":                        "British English",
			"locale":                              "en-GB",
			"enabled":                             "true",
			"text_checkbox_accept":                "",
			"text_button_continue":                "",
			"text_button_decline":                 "",
			"latest_revision.0.#":                 "1",
			"latest_revision.0.content_type":      "text/html",
			"latest_revision.0.effective_at":      "2100-01-01T01:01:00.000Z",
			"latest_revision.0.require_reconsent": "false",
			"latest_revision.0.text":              "I started a band called 999MB.  We haven't yet got a gig.",
		}),
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckAgreementDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			{
				Config: testAccAgreementConfig_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config:  testAccAgreementConfig_Full(resourceName, name),
				Destroy: true,
			},
			// Minimal
			{
				Config: testAccAgreementConfig_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config:  testAccAgreementConfig_Minimal(resourceName, name),
				Destroy: true,
			},
			// Change
			{
				Config: testAccAgreementConfig_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config: testAccAgreementConfig_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config: testAccAgreementConfig_Full(resourceName, name),
				Check:  fullCheck,
			},
		},
	})
}

func testAccAgreementConfig_NewEnv(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

		resource "pingone_language" "%[2]s" {
			environment_id = pingone_environment.%[2]s.id
		  
			locale = "en-GB"
		  }
		  
		  resource "pingone_language_update" "%[2]s" {
			environment_id = pingone_environment.%[2]s.id
		  
			language_id = pingone_language.%[2]s.id
			enabled     = true
			default     = false
		  }

resource "pingone_agreement" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  name = "%[4]s"

  localized_text {
	display_name = "British English"

	locale = pingone_language.%[2]s.locale

	latest_revision {
		content_type = "text/html"
		effective_at = "2100-01-01T01:01:00.000Z"
		require_reconsent = false
		text = "I started a band called 999MB.  We haven't yet got a gig."
	}
  }

  depends_on = [
	pingone_language_update.%[2]s
  ]
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccAgreementConfig_Full(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_agreement" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name                  = "%[3]s"
  description           = "Before the crowbar was invented, Crows would just drink at home."
  enabled               = false
  reconsent_period_days = 30

  localized_text {
	display_name = "British English"

	locale = "en-GB"
	enabled = true

	text_checkbox_accept = "Yeah"
	text_button_continue = "Go on"
	text_button_decline = "Nah"

	latest_revision {
		content_type = "text/html"
		effective_at = "2100-01-01T01:01:00.000Z"
		require_reconsent = false
		text = "I started a band called 999MB.  We haven't yet got a gig."
	}
  }

  localized_text {
	display_name = "English"

	locale = "en"
	enabled = false

	text_checkbox_accept = "Accept"
	text_button_continue = "Continue"
	text_button_decline = "Decline"

	latest_revision {
		content_type = "text/plain"
		effective_at = "2100-01-01T01:01:00.000Z"
		require_reconsent = false
		text = "I started a band called 999MB.  We haven't yet got a gig."
	}
  }

  localized_text {
	display_name = "French"

	locale = "fr"
	enabled = true

	text_checkbox_accept = "Accepter"
	text_button_continue = "Continuer"
	text_button_decline = "Déclin"

	latest_revision {
		content_type = "text/html"
		effective_at = "2100-01-01T01:01:00.000Z"
		require_reconsent = true
		text = "J'ai monté un groupe appelé 999MB. Nous n'avons pas encore de concert."
	}
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccAgreementConfig_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_agreement" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  localized_text {
	display_name = "British English"

	locale = "en-GB"

	latest_revision {
		content_type = "text/html"
		effective_at = "2100-01-01T01:01:00.000Z"
		require_reconsent = false
		text = "I started a band called 999MB.  We haven't yet got a gig."
	}
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
