package sso_test

// TODO test conditions

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	pingone "github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
)

func testAccCheckSignOnPolicyActionDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, pingone.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_sign_on_policy_action" {
			continue
		}

		_, rEnv, err := apiClient.EnvironmentsApi.ReadOneEnvironment(ctx, rs.Primary.Attributes["environment_id"]).Execute()

		if rEnv.StatusCode == 404 {
			continue
		}

		if err != nil {
			return err
		}

		body, r, err := apiClient.SignOnPolicyActionsApi.ReadOneSignOnPolicyAction(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["sign_on_policy_id"], rs.Primary.ID).Execute()

		if r.StatusCode == 404 {
			continue
		}

		if err != nil {
			tflog.Error(ctx, fmt.Sprintf("Error: %v", body))
			return err
		}

		return fmt.Errorf("PingOne Sign on Policy Action %s still exists", rs.Primary.ID)
	}

	return nil
}

func TestAccSignOnPolicyAction_LoginAction(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_sign_on_policy_action.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckSignOnPolicyActionDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSignOnPolicyActionConfig_LoginFullWithExt(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "sign_on_policy_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "priority", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "registration_external_href", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "registration_local_population_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "registration_confirm_user_attributes", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "social_provider_ids.#", "2"),
					resource.TestCheckResourceAttr(resourceFullName, "enforce_lockout_for_identity_providers", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "login.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "login.0.recovery_enabled", "false"),
				),
			},
			{
				Config: testAccSignOnPolicyActionConfig_LoginMinimal(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "sign_on_policy_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "priority", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "registration_external_href", ""),
					resource.TestCheckResourceAttr(resourceFullName, "registration_local_population_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "registration_confirm_user_attributes", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "social_provider_ids.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "enforce_lockout_for_identity_providers", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "login.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "login.0.recovery_enabled", "true"),
				),
			},
			{
				Config: testAccSignOnPolicyActionConfig_LoginFullNoExt(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "sign_on_policy_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "priority", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "registration_external_href", ""),
					resource.TestMatchResourceAttr(resourceFullName, "registration_local_population_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "registration_confirm_user_attributes", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "social_provider_ids.#", "2"),
					resource.TestCheckResourceAttr(resourceFullName, "enforce_lockout_for_identity_providers", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "login.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "login.0.recovery_enabled", "false"),
				),
			},
			{
				Config: testAccSignOnPolicyActionConfig_LoginFullWithExt(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "sign_on_policy_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "priority", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "registration_external_href", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "registration_local_population_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "registration_confirm_user_attributes", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "social_provider_ids.#", "2"),
					resource.TestCheckResourceAttr(resourceFullName, "enforce_lockout_for_identity_providers", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "login.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "login.0.recovery_enabled", "false"),
				),
			},
		},
	})
}

func TestAccSignOnPolicyAction_IDFirstAction(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_sign_on_policy_action.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckSignOnPolicyActionDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSignOnPolicyActionConfig_IDFirstFullWithExt(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "sign_on_policy_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "registration_external_href", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "registration_local_population_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "registration_confirm_user_attributes", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "social_provider_ids.#", "2"),
					resource.TestCheckResourceAttr(resourceFullName, "enforce_lockout_for_identity_providers", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "identifier_first.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "identifier_first.0.recovery_enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "identifier_first.0.discovery_rule.#", "1"),
				),
			},
			{
				Config: testAccSignOnPolicyActionConfig_IDFirstMinimal(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "sign_on_policy_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "registration_external_href", ""),
					resource.TestCheckResourceAttr(resourceFullName, "registration_local_population_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "registration_confirm_user_attributes", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "social_provider_ids.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "enforce_lockout_for_identity_providers", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "identifier_first.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "identifier_first.0.recovery_enabled", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "identifier_first.0.discovery_rule.#", "0"),
				),
			},
			{
				Config: testAccSignOnPolicyActionConfig_IDFirstFullNoExt(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "sign_on_policy_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "registration_external_href", ""),
					resource.TestMatchResourceAttr(resourceFullName, "registration_local_population_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "registration_confirm_user_attributes", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "social_provider_ids.#", "2"),
					resource.TestCheckResourceAttr(resourceFullName, "enforce_lockout_for_identity_providers", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "identifier_first.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "identifier_first.0.recovery_enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "identifier_first.0.discovery_rule.#", "1"),
				),
			},
			{
				Config: testAccSignOnPolicyActionConfig_IDFirstFullWithExt(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "sign_on_policy_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "registration_external_href", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "registration_local_population_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "registration_confirm_user_attributes", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "social_provider_ids.#", "2"),
					resource.TestCheckResourceAttr(resourceFullName, "enforce_lockout_for_identity_providers", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "identifier_first.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "identifier_first.0.recovery_enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "identifier_first.0.discovery_rule.#", "1"),
				),
			},
		},
	})
}

// func TestAccSignOnPolicyAction_MFAAction(t *testing.T) {
// }

func TestAccSignOnPolicyAction_IDPAction(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_sign_on_policy_action.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckSignOnPolicyActionDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSignOnPolicyActionConfig_IDPFull(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "sign_on_policy_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "registration_external_href", ""),
					resource.TestMatchResourceAttr(resourceFullName, "registration_local_population_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "registration_confirm_user_attributes", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "social_provider_ids.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "enforce_lockout_for_identity_providers", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "identity_provider.#", "1"),
					resource.TestMatchResourceAttr(resourceFullName, "identity_provider.0.identity_provider_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "identity_provider.0.acr_values", "MFA"),
					resource.TestCheckResourceAttr(resourceFullName, "identity_provider.0.pass_user_context", "true"),
				),
			},
			{
				Config: testAccSignOnPolicyActionConfig_IDPMinimal(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "sign_on_policy_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "registration_external_href", ""),
					resource.TestCheckResourceAttr(resourceFullName, "registration_local_population_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "registration_confirm_user_attributes", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "social_provider_ids.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "enforce_lockout_for_identity_providers", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "identity_provider.#", "1"),
					resource.TestMatchResourceAttr(resourceFullName, "identity_provider.0.identity_provider_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "identity_provider.0.acr_values", ""),
					resource.TestCheckResourceAttr(resourceFullName, "identity_provider.0.pass_user_context", "false"),
				),
			},
			{
				Config: testAccSignOnPolicyActionConfig_IDPFull(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "sign_on_policy_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "registration_external_href", ""),
					resource.TestMatchResourceAttr(resourceFullName, "registration_local_population_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "registration_confirm_user_attributes", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "social_provider_ids.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "enforce_lockout_for_identity_providers", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "identity_provider.#", "1"),
					resource.TestMatchResourceAttr(resourceFullName, "identity_provider.0.identity_provider_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "identity_provider.0.acr_values", "MFA"),
					resource.TestCheckResourceAttr(resourceFullName, "identity_provider.0.pass_user_context", "true"),
				),
			},
		},
	})
}

// func TestAccSignOnPolicyAction_AgreementAction(t *testing.T) {
// }

func TestAccSignOnPolicyAction_ProgressiveProfilingAction(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_sign_on_policy_action.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckSignOnPolicyActionDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSignOnPolicyActionConfig_ProgressiveProfilingFull(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "sign_on_policy_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "progressive_profiling.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "progressive_profiling.0.prevent_multiple_prompts_per_flow", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "progressive_profiling.0.prompt_interval_seconds", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "progressive_profiling.0.prompt_text", "Sed ut perspiciatis unde omnis iste natus error sit voluptatem accusantium doloremque laudantium, totam rem aperiam, eaque ipsa quae ab illo inventore veritatis et quasi architecto beatae vitae dicta sunt explicabo."),
					resource.TestCheckResourceAttr(resourceFullName, "progressive_profiling.0.attribute.#", "3"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "progressive_profiling.0.attribute.*", map[string]string{
						"name":     "address.postalCode",
						"required": "false",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "progressive_profiling.0.attribute.*", map[string]string{
						"name":     "name.given",
						"required": "true",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "progressive_profiling.0.attribute.*", map[string]string{
						"name":     "name.family",
						"required": "true",
					}),
				),
			},
			{
				Config: testAccSignOnPolicyActionConfig_ProgressiveProfilingMinimal(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "sign_on_policy_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "progressive_profiling.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "progressive_profiling.0.prevent_multiple_prompts_per_flow", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "progressive_profiling.0.prompt_interval_seconds", "7776000"),
					resource.TestCheckResourceAttr(resourceFullName, "progressive_profiling.0.prompt_text", "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."),
					resource.TestCheckResourceAttr(resourceFullName, "progressive_profiling.0.attribute.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "progressive_profiling.0.attribute.*", map[string]string{
						"name":     "address.postalCode",
						"required": "false",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "progressive_profiling.0.attribute.*", map[string]string{
						"name":     "email",
						"required": "true",
					}),
				),
			},
			{
				Config: testAccSignOnPolicyActionConfig_ProgressiveProfilingFull(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "sign_on_policy_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "progressive_profiling.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "progressive_profiling.0.prevent_multiple_prompts_per_flow", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "progressive_profiling.0.prompt_interval_seconds", "5"),
					resource.TestCheckResourceAttr(resourceFullName, "progressive_profiling.0.prompt_text", "Sed ut perspiciatis unde omnis iste natus error sit voluptatem accusantium doloremque laudantium, totam rem aperiam, eaque ipsa quae ab illo inventore veritatis et quasi architecto beatae vitae dicta sunt explicabo."),
					resource.TestCheckResourceAttr(resourceFullName, "progressive_profiling.0.attribute.#", "3"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "progressive_profiling.0.attribute.*", map[string]string{
						"name":     "address.postalCode",
						"required": "false",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "progressive_profiling.0.attribute.*", map[string]string{
						"name":     "name.given",
						"required": "true",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "progressive_profiling.0.attribute.*", map[string]string{
						"name":     "name.family",
						"required": "true",
					}),
				),
			},
		},
	})
}

func TestAccSignOnPolicyAction_MultipleActionChange(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_sign_on_policy_action.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckSignOnPolicyActionDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSignOnPolicyActionConfig_Multiple1(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-1", resourceFullName), "priority", "1"),
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-1", resourceFullName), "identifier_first.#", "1"),
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-2", resourceFullName), "priority", "2"),
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-2", resourceFullName), "login.#", "1"),
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-3", resourceFullName), "priority", "3"),
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-3", resourceFullName), "progressive_profiling.#", "1"),
				),
			},
			{
				Config: testAccSignOnPolicyActionConfig_Multiple2(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-1", resourceFullName), "priority", "1"),
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-1", resourceFullName), "login.#", "1"),
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-2", resourceFullName), "priority", "2"),
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-2", resourceFullName), "progressive_profiling.#", "1"),
				),
			},
			{
				Config: testAccSignOnPolicyActionConfig_Multiple1(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-1", resourceFullName), "priority", "1"),
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-1", resourceFullName), "identifier_first.#", "1"),
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-2", resourceFullName), "priority", "2"),
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-2", resourceFullName), "login.#", "1"),
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-3", resourceFullName), "priority", "3"),
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-3", resourceFullName), "progressive_profiling.#", "1"),
				),
			},
		},
	})
}

func TestAccSignOnPolicyAction_ConditionsSignOnOlderThanSingle(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_sign_on_policy_action.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckSignOnPolicyActionDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSignOnPolicyActionConfig_ConditionsSignOnOlderThanSingle(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.last_sign_on_older_than_seconds", "3600"),
				),
			},
			{
				Config: testAccSignOnPolicyActionConfig_LoginMinimal(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "0"),
				),
			},
			{
				Config: testAccSignOnPolicyActionConfig_ConditionsSignOnOlderThanSingle(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.last_sign_on_older_than_seconds", "3600"),
				),
			},
		},
	})
}

func TestAccSignOnPolicyAction_ConditionsMemberOfPopulation(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_sign_on_policy_action.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckSignOnPolicyActionDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSignOnPolicyActionConfig_ConditionsMemberOfPopulation(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.#", "1"),
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.0.last_sign_on_older_than_seconds", "3600"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.user_is_member_of_any_population_id.#", "1"),
					resource.TestMatchResourceAttr(resourceFullName, "conditions.0.user_is_member_of_any_population_id.0", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
				),
			},
			{
				Config: testAccSignOnPolicyActionConfig_LoginMinimal(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "0"),
				),
			},
			{
				Config: testAccSignOnPolicyActionConfig_ConditionsMemberOfPopulation(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.#", "1"),
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.0.last_sign_on_older_than_seconds", "3600"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.user_is_member_of_any_population_id.#", "1"),
					resource.TestMatchResourceAttr(resourceFullName, "conditions.0.user_is_member_of_any_population_id.0", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
				),
			},
		},
	})
}

func TestAccSignOnPolicyAction_ConditionsMemberOfPopulations(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_sign_on_policy_action.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckSignOnPolicyActionDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSignOnPolicyActionConfig_ConditionsMemberOfPopulations(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.#", "1"),
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.0.last_sign_on_older_than_seconds", "3600"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.user_is_member_of_any_population_id.#", "3"),
					resource.TestMatchResourceAttr(resourceFullName, "conditions.0.user_is_member_of_any_population_id.0", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "conditions.0.user_is_member_of_any_population_id.1", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "conditions.0.user_is_member_of_any_population_id.2", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
				),
			},
			{
				Config: testAccSignOnPolicyActionConfig_LoginMinimal(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "0"),
				),
			},
			{
				Config: testAccSignOnPolicyActionConfig_ConditionsMemberOfPopulations(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.#", "1"),
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.0.last_sign_on_older_than_seconds", "3600"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.user_is_member_of_any_population_id.#", "3"),
					resource.TestMatchResourceAttr(resourceFullName, "conditions.0.user_is_member_of_any_population_id.0", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "conditions.0.user_is_member_of_any_population_id.1", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "conditions.0.user_is_member_of_any_population_id.2", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
				),
			},
		},
	})
}

func TestAccSignOnPolicyAction_ConditionsUserAttributeEqualsSingle(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_sign_on_policy_action.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckSignOnPolicyActionDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSignOnPolicyActionConfig_ConditionsUserAttributeEqualsSingle(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.#", "1"),
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.0.last_sign_on_older_than_seconds", "3600"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.user_attribute_equals.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "conditions.0.user_attribute_equals.*", map[string]string{
						"attribute_reference": "${user.lifecycle.status}",
						"value":               "ACCOUNT_OK",
					}),
				),
			},
			{
				Config: testAccSignOnPolicyActionConfig_LoginMinimal(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "0"),
				),
			},
			{
				Config: testAccSignOnPolicyActionConfig_ConditionsUserAttributeEqualsSingle(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.#", "1"),
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.0.last_sign_on_older_than_seconds", "3600"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.user_attribute_equals.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "conditions.0.user_attribute_equals.*", map[string]string{
						"attribute_reference": "${user.lifecycle.status}",
						"value":               "ACCOUNT_OK",
					}),
				),
			},
		},
	})
}

func TestAccSignOnPolicyAction_ConditionsUserAttributeEqualsMultiple(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_sign_on_policy_action.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckSignOnPolicyActionDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSignOnPolicyActionConfig_ConditionsUserAttributeEqualsMultiple(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.#", "1"),
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.0.last_sign_on_older_than_seconds", "3600"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.user_attribute_equals.#", "3"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "conditions.0.user_attribute_equals.*", map[string]string{
						"attribute_reference": "${user.lifecycle.status}",
						"value":               "ACCOUNT_OK",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "conditions.0.user_attribute_equals.*", map[string]string{
						"attribute_reference": "${user.name.given}",
						"value":               "Bruce",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "conditions.0.user_attribute_equals.*", map[string]string{
						"attribute_reference": "${user.name.family}",
						"value":               "Wayne",
					}),
				),
			},
			{
				Config: testAccSignOnPolicyActionConfig_LoginMinimal(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "0"),
				),
			},
			{
				Config: testAccSignOnPolicyActionConfig_ConditionsUserAttributeEqualsMultiple(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.#", "1"),
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.0.last_sign_on_older_than_seconds", "3600"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.user_attribute_equals.#", "3"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "conditions.0.user_attribute_equals.*", map[string]string{
						"attribute_reference": "${user.lifecycle.status}",
						"value":               "ACCOUNT_OK",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "conditions.0.user_attribute_equals.*", map[string]string{
						"attribute_reference": "${user.name.given}",
						"value":               "Bruce",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "conditions.0.user_attribute_equals.*", map[string]string{
						"attribute_reference": "${user.name.family}",
						"value":               "Wayne",
					}),
				),
			},
		},
	})
}

// func TestAccSignOnPolicyAction_ConditionsIPOutOfRangeSingle(t *testing.T) {
// 	t.Parallel()

// 	resourceName := acctest.ResourceNameGen()
// 	resourceFullName := fmt.Sprintf("pingone_sign_on_policy_action.%s", resourceName)

// 	environmentName := acctest.ResourceNameGenEnvironment()

// 	name := resourceName

// 	licenseID := os.Getenv("PINGONE_LICENSE_ID")

// 	resource.Test(t, resource.TestCase{
// 		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
// 		ProviderFactories: acctest.ProviderFactories,
// 		CheckDestroy:      testAccCheckSignOnPolicyActionDestroy,
// 		ErrorCheck:        acctest.ErrorCheck(t),
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccSignOnPolicyActionConfig_ConditionsIPOutOfRangeSingle(environmentName, licenseID, resourceName, name),
// 				Check: resource.ComposeTestCheckFunc(
// 					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.#", "1"),
// 					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.0.last_sign_on_older_than_seconds", "3600"),
// 					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "1"),
// 					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.ip_out_of_range_cidr.#", "1"),
// 					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.ip_out_of_range_cidr.0", "192.168.129.23/17"),
// 				),
// 			},
// 			{
// 				Config: testAccSignOnPolicyActionConfig_LoginMinimal(environmentName, licenseID, resourceName, name),
// 				Check: resource.ComposeTestCheckFunc(
// 					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "0"),
// 				),
// 			},
// 			{
// 				Config: testAccSignOnPolicyActionConfig_ConditionsIPOutOfRangeSingle(environmentName, licenseID, resourceName, name),
// 				Check: resource.ComposeTestCheckFunc(
// 					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.#", "1"),
// 					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.0.last_sign_on_older_than_seconds", "3600"),
// 					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "1"),
// 					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.ip_out_of_range_cidr.#", "1"),
// 					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.ip_out_of_range_cidr.0", "192.168.129.23/17"),
// 				),
// 			},
// 		},
// 	})
// }

// func TestAccSignOnPolicyAction_ConditionsIPOutOfRangeMultiple(t *testing.T) {
// 	t.Parallel()

// 	resourceName := acctest.ResourceNameGen()
// 	resourceFullName := fmt.Sprintf("pingone_sign_on_policy_action.%s", resourceName)

// 	environmentName := acctest.ResourceNameGenEnvironment()

// 	name := resourceName

// 	licenseID := os.Getenv("PINGONE_LICENSE_ID")

// 	resource.Test(t, resource.TestCase{
// 		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
// 		ProviderFactories: acctest.ProviderFactories,
// 		CheckDestroy:      testAccCheckSignOnPolicyActionDestroy,
// 		ErrorCheck:        acctest.ErrorCheck(t),
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccSignOnPolicyActionConfig_ConditionsIPOutOfRangeMultiple(environmentName, licenseID, resourceName, name),
// 				Check: resource.ComposeTestCheckFunc(
// 					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.#", "1"),
// 					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.0.last_sign_on_older_than_seconds", "3600"),
// 					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "1"),
// 					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.ip_out_of_range_cidr.#", "2"),
// 					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.ip_out_of_range_cidr.0", "192.168.129.23/17"),
// 					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.ip_out_of_range_cidr.1", "192.168.0.15/24"),
// 				),
// 			},
// 			{
// 				Config: testAccSignOnPolicyActionConfig_LoginMinimal(environmentName, licenseID, resourceName, name),
// 				Check: resource.ComposeTestCheckFunc(
// 					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "0"),
// 				),
// 			},
// 			{
// 				Config: testAccSignOnPolicyActionConfig_ConditionsIPOutOfRangeMultiple(environmentName, licenseID, resourceName, name),
// 				Check: resource.ComposeTestCheckFunc(
// 					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.#", "1"),
// 					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.0.last_sign_on_older_than_seconds", "3600"),
// 					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "1"),
// 					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.ip_out_of_range_cidr.#", "2"),
// 					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.ip_out_of_range_cidr.0", "192.168.129.23/17"),
// 					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.ip_out_of_range_cidr.1", "192.168.0.15/24"),
// 				),
// 			},
// 		},
// 	})
// }

// func TestAccSignOnPolicyAction_ConditionsIPHighRisk(t *testing.T) {
// 	t.Parallel()

// 	resourceName := acctest.ResourceNameGen()
// 	resourceFullName := fmt.Sprintf("pingone_sign_on_policy_action.%s", resourceName)

// 	environmentName := acctest.ResourceNameGenEnvironment()

// 	name := resourceName

// 	licenseID := os.Getenv("PINGONE_LICENSE_ID")

// 	resource.Test(t, resource.TestCase{
// 		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
// 		ProviderFactories: acctest.ProviderFactories,
// 		CheckDestroy:      testAccCheckSignOnPolicyActionDestroy,
// 		ErrorCheck:        acctest.ErrorCheck(t),
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccSignOnPolicyActionConfig_ConditionsIPHighRisk(environmentName, licenseID, resourceName, name),
// 				Check: resource.ComposeTestCheckFunc(
// 					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.#", "1"),
// 					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.0.last_sign_on_older_than_seconds", "3600"),
// 					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "1"),
// 					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.ip_reputation_high_risk", "true"),
// 				),
// 			},
// 			{
// 				Config: testAccSignOnPolicyActionConfig_LoginMinimal(environmentName, licenseID, resourceName, name),
// 				Check: resource.ComposeTestCheckFunc(
// 					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "0"),
// 				),
// 			},
// 			{
// 				Config: testAccSignOnPolicyActionConfig_ConditionsIPHighRisk(environmentName, licenseID, resourceName, name),
// 				Check: resource.ComposeTestCheckFunc(
// 					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.#", "1"),
// 					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.0.last_sign_on_older_than_seconds", "3600"),
// 					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "1"),
// 					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.ip_reputation_high_risk", "true"),
// 				),
// 			},
// 		},
// 	})
// }

// func TestAccSignOnPolicyAction_ConditionsGeovelocity(t *testing.T) {
// 	t.Parallel()

// 	resourceName := acctest.ResourceNameGen()
// 	resourceFullName := fmt.Sprintf("pingone_sign_on_policy_action.%s", resourceName)

// 	environmentName := acctest.ResourceNameGenEnvironment()

// 	name := resourceName

// 	licenseID := os.Getenv("PINGONE_LICENSE_ID")

// 	resource.Test(t, resource.TestCase{
// 		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
// 		ProviderFactories: acctest.ProviderFactories,
// 		CheckDestroy:      testAccCheckSignOnPolicyActionDestroy,
// 		ErrorCheck:        acctest.ErrorCheck(t),
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccSignOnPolicyActionConfig_ConditionsGeovelocity(environmentName, licenseID, resourceName, name),
// 				Check: resource.ComposeTestCheckFunc(
// 					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.#", "1"),
// 					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.0.last_sign_on_older_than_seconds", "3600"),
// 					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "1"),
// 					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.geovelocity_anomaly_detected", "true"),
// 				),
// 			},
// 			{
// 				Config: testAccSignOnPolicyActionConfig_LoginMinimal(environmentName, licenseID, resourceName, name),
// 				Check: resource.ComposeTestCheckFunc(
// 					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "0"),
// 				),
// 			},
// 			{
// 				Config: testAccSignOnPolicyActionConfig_ConditionsGeovelocity(environmentName, licenseID, resourceName, name),
// 				Check: resource.ComposeTestCheckFunc(
// 					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.#", "1"),
// 					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.0.last_sign_on_older_than_seconds", "3600"),
// 					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "1"),
// 					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.geovelocity_anomaly_detected", "true"),
// 				),
// 			},
// 		},
// 	})
// }

// func TestAccSignOnPolicyAction_ConditionsAnonymousNetwork(t *testing.T) {
// 	t.Parallel()

// 	resourceName := acctest.ResourceNameGen()
// 	resourceFullName := fmt.Sprintf("pingone_sign_on_policy_action.%s", resourceName)

// 	environmentName := acctest.ResourceNameGenEnvironment()

// 	name := resourceName

// 	licenseID := os.Getenv("PINGONE_LICENSE_ID")

// 	resource.Test(t, resource.TestCase{
// 		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
// 		ProviderFactories: acctest.ProviderFactories,
// 		CheckDestroy:      testAccCheckSignOnPolicyActionDestroy,
// 		ErrorCheck:        acctest.ErrorCheck(t),
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccSignOnPolicyActionConfig_ConditionsAnonymousNetwork(environmentName, licenseID, resourceName, name),
// 				Check: resource.ComposeTestCheckFunc(
// 					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.#", "1"),
// 					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.0.last_sign_on_older_than_seconds", "3600"),
// 					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "1"),
// 					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.anonymous_network_detected", "true"),
// 					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.anonymous_network_detected_allowed_cidr.#", "0"),
// 				),
// 			},
// 			{
// 				Config: testAccSignOnPolicyActionConfig_LoginMinimal(environmentName, licenseID, resourceName, name),
// 				Check: resource.ComposeTestCheckFunc(
// 					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "0"),
// 				),
// 			},
// 			{
// 				Config: testAccSignOnPolicyActionConfig_ConditionsAnonymousNetwork(environmentName, licenseID, resourceName, name),
// 				Check: resource.ComposeTestCheckFunc(
// 					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.#", "1"),
// 					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.0.last_sign_on_older_than_seconds", "3600"),
// 					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "1"),
// 					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.anonymous_network_detected", "true"),
// 					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.anonymous_network_detected_allowed_cidr.#", "0"),
// 				),
// 			},
// 		},
// 	})
// }

// func TestAccSignOnPolicyAction_ConditionsAnonymousNetworkWithAllowed(t *testing.T) {
// 	t.Parallel()

// 	resourceName := acctest.ResourceNameGen()
// 	resourceFullName := fmt.Sprintf("pingone_sign_on_policy_action.%s", resourceName)

// 	environmentName := acctest.ResourceNameGenEnvironment()

// 	name := resourceName

// 	licenseID := os.Getenv("PINGONE_LICENSE_ID")

// 	resource.Test(t, resource.TestCase{
// 		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
// 		ProviderFactories: acctest.ProviderFactories,
// 		CheckDestroy:      testAccCheckSignOnPolicyActionDestroy,
// 		ErrorCheck:        acctest.ErrorCheck(t),
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccSignOnPolicyActionConfig_ConditionsAnonymousNetworkWithAllowed(environmentName, licenseID, resourceName, name),
// 				Check: resource.ComposeTestCheckFunc(
// 					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.#", "1"),
// 					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.0.last_sign_on_older_than_seconds", "3600"),
// 					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "1"),
// 					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.anonymous_network_detected", "true"),
// 					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.anonymous_network_detected_allowed_cidr.#", "2"),
// 					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.anonymous_network_detected_allowed_cidr.0", "192.168.129.23/17"),
// 					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.anonymous_network_detected_allowed_cidr.1", "192.168.0.15/24"),
// 				),
// 			},
// 			{
// 				Config: testAccSignOnPolicyActionConfig_LoginMinimal(environmentName, licenseID, resourceName, name),
// 				Check: resource.ComposeTestCheckFunc(
// 					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "0"),
// 				),
// 			},
// 			{
// 				Config: testAccSignOnPolicyActionConfig_ConditionsAnonymousNetworkWithAllowed(environmentName, licenseID, resourceName, name),
// 				Check: resource.ComposeTestCheckFunc(
// 					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.#", "1"),
// 					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.0.last_sign_on_older_than_seconds", "3600"),
// 					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "1"),
// 					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.anonymous_network_detected", "true"),
// 					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.anonymous_network_detected_allowed_cidr.#", "2"),
// 					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.anonymous_network_detected_allowed_cidr.0", "192.168.129.23/17"),
// 					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.anonymous_network_detected_allowed_cidr.1", "192.168.0.15/24"),
// 				),
// 			},
// 		},
// 	})
// }

func TestAccSignOnPolicyAction_ConditionsCompound(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_sign_on_policy_action.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckSignOnPolicyActionDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSignOnPolicyActionConfig_ConditionsCompoundSubset(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.user_is_member_of_any_population_id.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.user_attribute_equals.#", "2"),
				),
			},
			// {
			// 	Config: testAccSignOnPolicyActionConfig_ConditionsCompoundFull(environmentName, licenseID, resourceName, name),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "1"),
			// 		resource.TestCheckResourceAttr(resourceFullName, "conditions.0.user_is_member_of_any_population_id.#", "1"),
			// 		resource.TestCheckResourceAttr(resourceFullName, "conditions.0.user_attribute_equals.#", "2"),
			// 		resource.TestCheckResourceAttr(resourceFullName, "conditions.0.ip_out_of_range_cidr.#", "1"),
			// 		resource.TestCheckResourceAttr(resourceFullName, "conditions.0.ip_reputation_high_risk", "true"),
			// 		resource.TestCheckResourceAttr(resourceFullName, "conditions.0.geovelocity_anomaly_detected", "true"),
			// 		resource.TestCheckResourceAttr(resourceFullName, "conditions.0.anonymous_network_detected", "true"),
			// 		resource.TestCheckResourceAttr(resourceFullName, "conditions.0.anonymous_network_detected_allowed_cidr.#", "2"),
			// 	),
			// },
			// {
			// 	Config: testAccSignOnPolicyActionConfig_ConditionsCompoundSubset(environmentName, licenseID, resourceName, name),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "1"),
			// 		resource.TestCheckResourceAttr(resourceFullName, "conditions.0.user_is_member_of_any_population_id.#", "1"),
			// 		resource.TestCheckResourceAttr(resourceFullName, "conditions.0.user_attribute_equals.#", "2"),
			// 	),
			// },
		},
	})
}

func testAccSignOnPolicyActionConfig_LoginFullNoExt(environmentName, licenseID, resourceName, name string) string {

	return fmt.Sprintf(`
		%[1]s

		resource "pingone_identity_provider" "%[3]s-1" {
			environment_id = "${pingone_environment.%[2]s.id}"
			name = "%[4]s-1"
			
			google {
				client_id = "testclientid"
				client_secret = "testclientsecret"
			}
		}

		resource "pingone_identity_provider" "%[3]s-2" {
			environment_id = "${pingone_environment.%[2]s.id}"
			name = "%[4]s-2"
			
			facebook {
				app_id = "testclientid"
				app_secret = "testclientsecret"
			}
		}

		resource "pingone_sign_on_policy" "%[3]s" {
			environment_id = "${pingone_environment.%[2]s.id}"

			name = "%[4]s"
		}

		resource "pingone_sign_on_policy_action" "%[3]s" {
			environment_id 			 = "${pingone_environment.%[2]s.id}"
			sign_on_policy_id = "${pingone_sign_on_policy.%[3]s.id}"

			priority = 1

			registration_local_population_id = "${pingone_environment.%[2]s.default_population_id}"

			social_provider_ids = [,
				"${pingone_identity_provider.%[3]s-2.id}",
				"${pingone_identity_provider.%[3]s-1.id}"
			]

			login {
				recovery_enabled = false // we set this to false because the calculated default from the api is true
			}
		}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccSignOnPolicyActionConfig_LoginFullWithExt(environmentName, licenseID, resourceName, name string) string {

	return fmt.Sprintf(`
		%[1]s

		resource "pingone_identity_provider" "%[3]s-1" {
			environment_id = "${pingone_environment.%[2]s.id}"
			name = "%[4]s-1"
			
			google {
				client_id = "testclientid"
				client_secret = "testclientsecret"
			}
		}

		resource "pingone_identity_provider" "%[3]s-2" {
			environment_id = "${pingone_environment.%[2]s.id}"
			name = "%[4]s-2"
			
			facebook {
				app_id = "testclientid"
				app_secret = "testclientsecret"
			}
		}

		resource "pingone_sign_on_policy" "%[3]s" {
			environment_id = "${pingone_environment.%[2]s.id}"

			name = "%[4]s"
		}

		resource "pingone_sign_on_policy_action" "%[3]s" {
			environment_id 			 = "${pingone_environment.%[2]s.id}"
			sign_on_policy_id = "${pingone_sign_on_policy.%[3]s.id}"

			priority = 1

			registration_external_href = "https://www.pingidentity.com"

			social_provider_ids = [,
				"${pingone_identity_provider.%[3]s-2.id}",
				"${pingone_identity_provider.%[3]s-1.id}"
			]

			login {
				recovery_enabled = false // we set this to false because the calculated default from the api is true
			}
		}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccSignOnPolicyActionConfig_LoginMinimal(environmentName, licenseID, resourceName, name string) string {

	return fmt.Sprintf(`
		%[1]s

		resource "pingone_sign_on_policy" "%[3]s" {
			environment_id = "${pingone_environment.%[2]s.id}"

			name = "%[4]s"
		}

		resource "pingone_sign_on_policy_action" "%[3]s" {
			environment_id 			 = "${pingone_environment.%[2]s.id}"
			sign_on_policy_id = "${pingone_sign_on_policy.%[3]s.id}"

			priority = 1

			login {}
		}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

// TODO: idp
func testAccSignOnPolicyActionConfig_IDFirstFullWithExt(environmentName, licenseID, resourceName, name string) string {

	return fmt.Sprintf(`
		%[1]s

		resource "pingone_identity_provider" "%[3]s-1" {
			environment_id = "${pingone_environment.%[2]s.id}"
			name = "%[4]s-1"
			
			google {
				client_id = "testclientid"
				client_secret = "testclientsecret"
			}
		}

		resource "pingone_identity_provider" "%[3]s-2" {
			environment_id = "${pingone_environment.%[2]s.id}"
			name = "%[4]s-2"
			
			facebook {
				app_id = "testclientid"
				app_secret = "testclientsecret"
			}
		}

		resource "pingone_sign_on_policy" "%[3]s" {
			environment_id = "${pingone_environment.%[2]s.id}"

			name = "%[4]s"
		}

		resource "pingone_sign_on_policy_action" "%[3]s" {
			environment_id 			 = "${pingone_environment.%[2]s.id}"
			sign_on_policy_id = "${pingone_sign_on_policy.%[3]s.id}"

			priority = 1

			registration_external_href = "https://www.pingidentity.com"

			social_provider_ids = [,
				"${pingone_identity_provider.%[3]s-2.id}",
				"${pingone_identity_provider.%[3]s-1.id}"
			]

			identifier_first {
				recovery_enabled = false // we set this to false because the calculated default from the api is true
				discovery_rule {
				 	condition {
				 		contains = "domain.com"
				 		value = "value"
					}
					identity_provider_id ="${pingone_identity_provider.%[3]s-1.id}"
				}
			}
		}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccSignOnPolicyActionConfig_IDFirstFullNoExt(environmentName, licenseID, resourceName, name string) string {

	return fmt.Sprintf(`
		%[1]s

		resource "pingone_identity_provider" "%[3]s-1" {
			environment_id = "${pingone_environment.%[2]s.id}"
			name = "%[4]s-1"
			
			google {
				client_id = "testclientid"
				client_secret = "testclientsecret"
			}
		}

		resource "pingone_identity_provider" "%[3]s-2" {
			environment_id = "${pingone_environment.%[2]s.id}"
			name = "%[4]s-2"
			
			facebook {
				app_id = "testclientid"
				app_secret = "testclientsecret"
			}
		}

		resource "pingone_sign_on_policy" "%[3]s" {
			environment_id = "${pingone_environment.%[2]s.id}"

			name = "%[4]s"
		}

		resource "pingone_sign_on_policy_action" "%[3]s" {
			environment_id 			 = "${pingone_environment.%[2]s.id}"
			sign_on_policy_id = "${pingone_sign_on_policy.%[3]s.id}"

			priority = 1

			registration_local_population_id = "${pingone_environment.%[2]s.default_population_id}"

			social_provider_ids = [,
				"${pingone_identity_provider.%[3]s-2.id}",
				"${pingone_identity_provider.%[3]s-1.id}"
			]

			identifier_first {
				recovery_enabled = false // we set this to false because the calculated default from the api is true
				discovery_rule {
					condition {
						contains = "domain.com"
						value = "value"
				   }
				   identity_provider_id ="${pingone_identity_provider.%[3]s-1.id}"
			   }
			}
		}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccSignOnPolicyActionConfig_IDFirstMinimal(environmentName, licenseID, resourceName, name string) string {

	return fmt.Sprintf(`
		%[1]s

		resource "pingone_sign_on_policy" "%[3]s" {
			environment_id = "${pingone_environment.%[2]s.id}"

			name = "%[4]s"
		}

		resource "pingone_sign_on_policy_action" "%[3]s" {
			environment_id 			 = "${pingone_environment.%[2]s.id}"
			sign_on_policy_id = "${pingone_sign_on_policy.%[3]s.id}"

			priority = 1

			identifier_first {}

		}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

// TODO: MFA device policy data source
// func testAccSignOnPolicyActionConfig_MFAFull(environmentName, licenseID, resourceName, name string) string {
// }

// func testAccSignOnPolicyActionConfig_MFAMinimal(environmentName, licenseID, resourceName, name string) string {
// }

func testAccSignOnPolicyActionConfig_IDPFull(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

		resource "pingone_identity_provider" "%[3]s" {
			environment_id = "${pingone_environment.%[2]s.id}"
			name = "%[4]s"
			
			google {
				client_id = "testclientid"
				client_secret = "testclientsecret"
			}
		}

		resource "pingone_sign_on_policy" "%[3]s" {
			environment_id = "${pingone_environment.%[2]s.id}"

			name = "%[4]s"
		}

		resource "pingone_sign_on_policy_action" "%[3]s" {
			environment_id 			 = "${pingone_environment.%[2]s.id}"
			sign_on_policy_id = "${pingone_sign_on_policy.%[3]s.id}"

			priority = 1

			registration_local_population_id = "${pingone_environment.%[2]s.default_population_id}"

			registration_confirm_user_attributes = true

			identity_provider {
				identity_provider_id = "${pingone_identity_provider.%[3]s.id}"

				acr_values        = "MFA"
				pass_user_context = true
			}
			
		}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccSignOnPolicyActionConfig_IDPMinimal(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

		resource "pingone_identity_provider" "%[3]s" {
			environment_id = "${pingone_environment.%[2]s.id}"
			name = "%[4]s"
			
			google {
				client_id = "testclientid"
				client_secret = "testclientsecret"
			}
		}

		resource "pingone_sign_on_policy" "%[3]s" {
			environment_id = "${pingone_environment.%[2]s.id}"

			name = "%[4]s"
		}

		resource "pingone_sign_on_policy_action" "%[3]s" {
			environment_id 			 = "${pingone_environment.%[2]s.id}"
			sign_on_policy_id = "${pingone_sign_on_policy.%[3]s.id}"

			priority = 1

			identity_provider {
				identity_provider_id = "${pingone_identity_provider.%[3]s.id}"
			}
			
		}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

// TODO: agreements
// func testAccSignOnPolicyActionConfig_AgreementFull(environmentName, licenseID, resourceName, name string) string {
// }

// func testAccSignOnPolicyActionConfig_AgreementMinimal(environmentName, licenseID, resourceName, name string) string {
// }

func testAccSignOnPolicyActionConfig_ProgressiveProfilingFull(environmentName, licenseID, resourceName, name string) string {

	return fmt.Sprintf(`
		%[1]s

		resource "pingone_sign_on_policy" "%[3]s" {
			environment_id = "${pingone_environment.%[2]s.id}"

			name = "%[4]s"
		}

		resource "pingone_sign_on_policy_action" "%[3]s" {
			environment_id 			 = "${pingone_environment.%[2]s.id}"
			sign_on_policy_id = "${pingone_sign_on_policy.%[3]s.id}"

			priority = 1

			progressive_profiling {
				prevent_multiple_prompts_per_flow = false // default is true
				prompt_interval_seconds = 5 // default is 7776000
				prompt_text = "Sed ut perspiciatis unde omnis iste natus error sit voluptatem accusantium doloremque laudantium, totam rem aperiam, eaque ipsa quae ab illo inventore veritatis et quasi architecto beatae vitae dicta sunt explicabo."

				attribute {
					name = "name.given"
					required = true
				}

				attribute {
					name = "name.family"
					required = true
				}

				attribute {
					name = "address.postalCode"
					required = false
				}
			}
			
		}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccSignOnPolicyActionConfig_ProgressiveProfilingMinimal(environmentName, licenseID, resourceName, name string) string {

	return fmt.Sprintf(`
		%[1]s

		resource "pingone_sign_on_policy" "%[3]s" {
			environment_id = "${pingone_environment.%[2]s.id}"

			name = "%[4]s"
		}

		resource "pingone_sign_on_policy_action" "%[3]s" {
			environment_id 			 = "${pingone_environment.%[2]s.id}"
			sign_on_policy_id = "${pingone_sign_on_policy.%[3]s.id}"

			priority = 1

			progressive_profiling {
				prompt_text = "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."

				attribute {
					name = "email"
					required = true
				}

				attribute {
					name = "address.postalCode"
					required = false
				}
			}
		}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccSignOnPolicyActionConfig_Multiple1(environmentName, licenseID, resourceName, name string) string {

	return fmt.Sprintf(`
		%[1]s

		resource "pingone_sign_on_policy" "%[3]s" {
			environment_id = "${pingone_environment.%[2]s.id}"

			name = "%[4]s"
		}

		resource "pingone_sign_on_policy_action" "%[3]s-1" {
			environment_id 			 = "${pingone_environment.%[2]s.id}"
			sign_on_policy_id = "${pingone_sign_on_policy.%[3]s.id}"

			priority = 1

			identifier_first {}

		}
		
		resource "pingone_sign_on_policy_action" "%[3]s-2" {
			environment_id 			 = "${pingone_environment.%[2]s.id}"
			sign_on_policy_id = "${pingone_sign_on_policy.%[3]s.id}"

			priority = 2

			login {}

		}
		
		resource "pingone_sign_on_policy_action" "%[3]s-3" {
			environment_id 			 = "${pingone_environment.%[2]s.id}"
			sign_on_policy_id = "${pingone_sign_on_policy.%[3]s.id}"

			priority = 3

			progressive_profiling {
				prompt_text = "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."

				attribute {
					name = "email"
					required = true
				}

				attribute {
					name = "address.postalCode"
					required = false
				}
			}

		}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccSignOnPolicyActionConfig_Multiple2(environmentName, licenseID, resourceName, name string) string {

	return fmt.Sprintf(`
		%[1]s

		resource "pingone_sign_on_policy" "%[3]s" {
			environment_id = "${pingone_environment.%[2]s.id}"

			name = "%[4]s"
		}

		resource "pingone_sign_on_policy_action" "%[3]s-1" {
			environment_id 			 = "${pingone_environment.%[2]s.id}"
			sign_on_policy_id = "${pingone_sign_on_policy.%[3]s.id}"

			priority = 1

			login {}

		}
		
		resource "pingone_sign_on_policy_action" "%[3]s-2" {
			environment_id 			 = "${pingone_environment.%[2]s.id}"
			sign_on_policy_id = "${pingone_sign_on_policy.%[3]s.id}"

			priority = 2

			progressive_profiling {
				prompt_text = "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."

				attribute {
					name = "email"
					required = true
				}

				attribute {
					name = "address.postalCode"
					required = false
				}
			}

		}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

// func testAccSignOnPolicyActionConfig_ConditionsDeclaredNotDefined(environmentName, licenseID, resourceName, name string) string {

// 	return fmt.Sprintf(`
// 		%[1]s

// 		resource "pingone_sign_on_policy" "%[3]s" {
// 			environment_id = "${pingone_environment.%[2]s.id}"

// 			name = "%[4]s"
// 		}

// 		resource "pingone_sign_on_policy_action" "%[3]s" {
// 			environment_id 			 = "${pingone_environment.%[2]s.id}"
// 			sign_on_policy_id = "${pingone_sign_on_policy.%[3]s.id}"

// 			priority = 1

// 			conditions {}

// 			login {}

// 		}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
// }

func testAccSignOnPolicyActionConfig_ConditionsSignOnOlderThanSingle(environmentName, licenseID, resourceName, name string) string {

	return fmt.Sprintf(`
		%[1]s

		resource "pingone_sign_on_policy" "%[3]s" {
			environment_id = "${pingone_environment.%[2]s.id}"

			name = "%[4]s"
		}

		resource "pingone_sign_on_policy_action" "%[3]s" {
			environment_id 			 = "${pingone_environment.%[2]s.id}"
			sign_on_policy_id = "${pingone_sign_on_policy.%[3]s.id}"

			priority = 1

			conditions {
				last_sign_on_older_than_seconds = 3600
			}

			login {}

		}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccSignOnPolicyActionConfig_ConditionsMemberOfPopulation(environmentName, licenseID, resourceName, name string) string {

	return fmt.Sprintf(`
		%[1]s

		resource "pingone_sign_on_policy" "%[3]s" {
			environment_id = "${pingone_environment.%[2]s.id}"

			name = "%[4]s"
		}

		resource "pingone_sign_on_policy_action" "%[3]s-id" {
			environment_id 			 = "${pingone_environment.%[2]s.id}"
			sign_on_policy_id = "${pingone_sign_on_policy.%[3]s.id}"

			priority = 1

			conditions {
				last_sign_on_older_than_seconds = 3600
			}

			identifier_first {}

		}

		resource "pingone_sign_on_policy_action" "%[3]s" {
			environment_id 			 = "${pingone_environment.%[2]s.id}"
			sign_on_policy_id = "${pingone_sign_on_policy.%[3]s.id}"

			priority = 2

			conditions {
				user_is_member_of_any_population_id = [
					"${pingone_environment.%[2]s.default_population_id}"
				]
			}

			login {}

		}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccSignOnPolicyActionConfig_ConditionsMemberOfPopulations(environmentName, licenseID, resourceName, name string) string {

	return fmt.Sprintf(`
		%[1]s

		resource "pingone_population" "%[3]s-1" {
			environment_id = "${pingone_environment.%[2]s.id}"

			name = "Population 1"
		}

		resource "pingone_population" "%[3]s-2" {
			environment_id = "${pingone_environment.%[2]s.id}"

			name = "Population 2"
		}

		resource "pingone_sign_on_policy" "%[3]s" {
			environment_id = "${pingone_environment.%[2]s.id}"

			name = "%[4]s"
		}

		resource "pingone_sign_on_policy_action" "%[3]s-id" {
			environment_id 			 = "${pingone_environment.%[2]s.id}"
			sign_on_policy_id = "${pingone_sign_on_policy.%[3]s.id}"

			priority = 1

			conditions {
				last_sign_on_older_than_seconds = 3600
			}

			identifier_first {}

		}

		resource "pingone_sign_on_policy_action" "%[3]s" {
			environment_id 			 = "${pingone_environment.%[2]s.id}"
			sign_on_policy_id = "${pingone_sign_on_policy.%[3]s.id}"

			priority = 2

			conditions {
				user_is_member_of_any_population_id = [
					"${pingone_environment.%[2]s.default_population_id}",
					"${pingone_population.%[3]s-1.id}",
					"${pingone_population.%[3]s-2.id}"
				]
			}

			login {}

		}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccSignOnPolicyActionConfig_ConditionsUserAttributeEqualsSingle(environmentName, licenseID, resourceName, name string) string {

	return fmt.Sprintf(`
		%[1]s

		resource "pingone_sign_on_policy" "%[3]s" {
			environment_id = "${pingone_environment.%[2]s.id}"

			name = "%[4]s"
		}

		resource "pingone_sign_on_policy_action" "%[3]s-id" {
			environment_id 			 = "${pingone_environment.%[2]s.id}"
			sign_on_policy_id = "${pingone_sign_on_policy.%[3]s.id}"

			priority = 1

			conditions {
				last_sign_on_older_than_seconds = 3600
			}

			identifier_first {}

		}

		resource "pingone_sign_on_policy_action" "%[3]s" {
			environment_id 			 = "${pingone_environment.%[2]s.id}"
			sign_on_policy_id = "${pingone_sign_on_policy.%[3]s.id}"

			priority = 2

			conditions {
				user_attribute_equals {
					attribute_reference = "$${user.lifecycle.status}"
					value 				= "ACCOUNT_OK"
				}
			}

			login {}

		}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccSignOnPolicyActionConfig_ConditionsUserAttributeEqualsMultiple(environmentName, licenseID, resourceName, name string) string {

	return fmt.Sprintf(`
		%[1]s

		resource "pingone_sign_on_policy" "%[3]s" {
			environment_id = "${pingone_environment.%[2]s.id}"

			name = "%[4]s"
		}

		resource "pingone_sign_on_policy_action" "%[3]s-id" {
			environment_id 			 = "${pingone_environment.%[2]s.id}"
			sign_on_policy_id = "${pingone_sign_on_policy.%[3]s.id}"

			priority = 1

			conditions {
				last_sign_on_older_than_seconds = 3600
			}

			identifier_first {}

		}

		resource "pingone_sign_on_policy_action" "%[3]s" {
			environment_id 			 = "${pingone_environment.%[2]s.id}"
			sign_on_policy_id = "${pingone_sign_on_policy.%[3]s.id}"

			priority = 2

			conditions {
				user_attribute_equals {
					attribute_reference = "$${user.name.given}"
					value 				= "Bruce"
				}

				user_attribute_equals {
					attribute_reference = "$${user.name.family}"
					value 				= "Wayne"
				}

				user_attribute_equals {
					attribute_reference = "$${user.lifecycle.status}"
					value 				= "ACCOUNT_OK"
				}
			}

			login {}

		}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

// func testAccSignOnPolicyActionConfig_ConditionsIPOutOfRangeSingle(environmentName, licenseID, resourceName, name string) string {

// 	return fmt.Sprintf(`
// 		%[1]s

// 		resource "pingone_sign_on_policy" "%[3]s" {
// 			environment_id = "${pingone_environment.%[2]s.id}"

// 			name = "%[4]s"
// 		}

// 		resource "pingone_sign_on_policy_action" "%[3]s-id" {
// 			environment_id 			 = "${pingone_environment.%[2]s.id}"
// 			sign_on_policy_id = "${pingone_sign_on_policy.%[3]s.id}"

// 			priority = 1

// 			conditions {
// 				last_sign_on_older_than_seconds = 3600
// 			}

// 			identifier_first {}

// 		}

// 		resource "pingone_sign_on_policy_action" "%[3]s" {
// 			environment_id 			 = "${pingone_environment.%[2]s.id}"
// 			sign_on_policy_id = "${pingone_sign_on_policy.%[3]s.id}"

// 			priority = 2

// 			conditions {
// 				ip_out_of_range_cidr = [
// 					"192.168.129.23/17"
// 				]
// 			}

// 			mfa {}

// 		}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
// }

// func testAccSignOnPolicyActionConfig_ConditionsIPOutOfRangeMultiple(environmentName, licenseID, resourceName, name string) string {

// 	return fmt.Sprintf(`
// 		%[1]s

// 		resource "pingone_sign_on_policy" "%[3]s" {
// 			environment_id = "${pingone_environment.%[2]s.id}"

// 			name = "%[4]s"
// 		}

// 		resource "pingone_sign_on_policy_action" "%[3]s-id" {
// 			environment_id 			 = "${pingone_environment.%[2]s.id}"
// 			sign_on_policy_id = "${pingone_sign_on_policy.%[3]s.id}"

// 			priority = 1

// 			conditions {
// 				last_sign_on_older_than_seconds = 3600
// 			}

// 			identifier_first {}

// 		}

// 		resource "pingone_sign_on_policy_action" "%[3]s" {
// 			environment_id 			 = "${pingone_environment.%[2]s.id}"
// 			sign_on_policy_id = "${pingone_sign_on_policy.%[3]s.id}"

// 			priority = 2

// 			conditions {
// 				ip_out_of_range_cidr = [
// 					"192.168.129.23/17",
// 					"192.168.0.15/24"
// 				]
// 			}

// 			mfa {}

// 		}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
// }

// func testAccSignOnPolicyActionConfig_ConditionsIPHighRisk(environmentName, licenseID, resourceName, name string) string {

// 	return fmt.Sprintf(`
// 		%[1]s

// 		resource "pingone_sign_on_policy" "%[3]s" {
// 			environment_id = "${pingone_environment.%[2]s.id}"

// 			name = "%[4]s"
// 		}

// 		resource "pingone_sign_on_policy_action" "%[3]s-id" {
// 			environment_id 			 = "${pingone_environment.%[2]s.id}"
// 			sign_on_policy_id = "${pingone_sign_on_policy.%[3]s.id}"

// 			priority = 1

// 			conditions {
// 				last_sign_on_older_than_seconds = 3600
// 			}

// 			identifier_first {}

// 		}

// 		resource "pingone_sign_on_policy_action" "%[3]s" {
// 			environment_id 			 = "${pingone_environment.%[2]s.id}"
// 			sign_on_policy_id = "${pingone_sign_on_policy.%[3]s.id}"

// 			priority = 2

// 			conditions {
// 				ip_reputation_high_risk = true
// 			}

// 			mfa {}

// 		}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
// }

// func testAccSignOnPolicyActionConfig_ConditionsGeovelocity(environmentName, licenseID, resourceName, name string) string {

// 	return fmt.Sprintf(`
// 		%[1]s

// 		resource "pingone_sign_on_policy" "%[3]s" {
// 			environment_id = "${pingone_environment.%[2]s.id}"

// 			name = "%[4]s"
// 		}

// 		resource "pingone_sign_on_policy_action" "%[3]s-id" {
// 			environment_id 			 = "${pingone_environment.%[2]s.id}"
// 			sign_on_policy_id = "${pingone_sign_on_policy.%[3]s.id}"

// 			priority = 1

// 			conditions {
// 				last_sign_on_older_than_seconds = 3600
// 			}

// 			identifier_first {}

// 		}

// 		resource "pingone_sign_on_policy_action" "%[3]s" {
// 			environment_id 			 = "${pingone_environment.%[2]s.id}"
// 			sign_on_policy_id = "${pingone_sign_on_policy.%[3]s.id}"

// 			priority = 2

// 			conditions {
// 				geovelocity_anomaly_detected = true
// 			}

// 			mfa {}

// 		}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
// }

// func testAccSignOnPolicyActionConfig_ConditionsAnonymousNetwork(environmentName, licenseID, resourceName, name string) string {

// 	return fmt.Sprintf(`
// 		%[1]s

// 		resource "pingone_sign_on_policy" "%[3]s" {
// 			environment_id = "${pingone_environment.%[2]s.id}"

// 			name = "%[4]s"
// 		}

// 		resource "pingone_sign_on_policy_action" "%[3]s-id" {
// 			environment_id 			 = "${pingone_environment.%[2]s.id}"
// 			sign_on_policy_id = "${pingone_sign_on_policy.%[3]s.id}"

// 			priority = 1

// 			conditions {
// 				last_sign_on_older_than_seconds = 3600
// 			}

// 			identifier_first {}

// 		}

// 		resource "pingone_sign_on_policy_action" "%[3]s" {
// 			environment_id 			 = "${pingone_environment.%[2]s.id}"
// 			sign_on_policy_id = "${pingone_sign_on_policy.%[3]s.id}"

// 			priority = 2

// 			conditions {
// 				anonymous_network_detected = true
// 			}

// 			mfa {}

// 		}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
// }

// func testAccSignOnPolicyActionConfig_ConditionsAnonymousNetworkWithAllowed(environmentName, licenseID, resourceName, name string) string {

// 	return fmt.Sprintf(`
// 		%[1]s

// 		resource "pingone_sign_on_policy" "%[3]s" {
// 			environment_id = "${pingone_environment.%[2]s.id}"

// 			name = "%[4]s"
// 		}

// 		resource "pingone_sign_on_policy_action" "%[3]s-id" {
// 			environment_id 			 = "${pingone_environment.%[2]s.id}"
// 			sign_on_policy_id = "${pingone_sign_on_policy.%[3]s.id}"

// 			priority = 1

// 			conditions {
// 				last_sign_on_older_than_seconds = 3600
// 			}

// 			identifier_first {}

// 		}

// 		resource "pingone_sign_on_policy_action" "%[3]s" {
// 			environment_id 			 = "${pingone_environment.%[2]s.id}"
// 			sign_on_policy_id = "${pingone_sign_on_policy.%[3]s.id}"

// 			priority = 2

// 			conditions {
// 				anonymous_network_detected = true

// 				anonymous_network_detected_allowed_cidr = [
// 					"192.168.129.23/17",
// 					"192.168.0.15/24"
// 				]
// 			}

// 			mfa {}

// 		}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
// }

func testAccSignOnPolicyActionConfig_ConditionsCompoundSubset(environmentName, licenseID, resourceName, name string) string {

	return fmt.Sprintf(`
		%[1]s

		resource "pingone_sign_on_policy" "%[3]s" {
			environment_id = "${pingone_environment.%[2]s.id}"

			name = "%[4]s"
		}

		resource "pingone_sign_on_policy_action" "%[3]s-id" {
			environment_id 			 = "${pingone_environment.%[2]s.id}"
			sign_on_policy_id = "${pingone_sign_on_policy.%[3]s.id}"

			priority = 1

			conditions {
				last_sign_on_older_than_seconds = 3600
			}

			identifier_first {}

		}

		resource "pingone_sign_on_policy_action" "%[3]s" {
			environment_id 			 = "${pingone_environment.%[2]s.id}"
			sign_on_policy_id = "${pingone_sign_on_policy.%[3]s.id}"

			priority = 2

			conditions {

				user_is_member_of_any_population_id = [
					"${pingone_environment.%[2]s.default_population_id}"
				]

				user_attribute_equals {
					attribute_reference = "$${user.name.given}"
					value 				= "Bruce"
				}

				user_attribute_equals {
					attribute_reference = "$${user.name.family}"
					value 				= "Wayne"
				}

			}

			login {}

		}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

// func testAccSignOnPolicyActionConfig_ConditionsCompoundFull(environmentName, licenseID, resourceName, name string) string {

// 	return fmt.Sprintf(`
// 		%[1]s

// 		resource "pingone_sign_on_policy" "%[3]s" {
// 			environment_id = "${pingone_environment.%[2]s.id}"

// 			name = "%[4]s"
// 		}

// 		resource "pingone_sign_on_policy_action" "%[3]s-id" {
// 			environment_id 			 = "${pingone_environment.%[2]s.id}"
// 			sign_on_policy_id = "${pingone_sign_on_policy.%[3]s.id}"

// 			priority = 1

// 			conditions {
// 				last_sign_on_older_than_seconds = 3600
// 			}

// 			identifier_first {}

// 		}

// 		resource "pingone_sign_on_policy_action" "%[3]s-login" {
// 			environment_id 			 = "${pingone_environment.%[2]s.id}"
// 			sign_on_policy_id = "${pingone_sign_on_policy.%[3]s.id}"

// 			priority = 2

// 			conditions {

// 				user_is_member_of_any_population_id = [
// 					"${pingone_environment.%[2]s.default_population_id}"
// 				]

// 				user_attribute_equals {
// 					attribute_reference = "$${user.name.given}"
// 					value 				= "Bruce"
// 				}

// 				user_attribute_equals {
// 					attribute_reference = "$${user.name.family}"
// 					value 				= "Wayne"
// 				}

// 			}

// 			login {}

// 		}

// 		resource "pingone_sign_on_policy_action" "%[3]s-login" {
// 			environment_id 			 = "${pingone_environment.%[2]s.id}"
// 			sign_on_policy_id = "${pingone_sign_on_policy.%[3]s.id}"

// 			priority = 3

// 			conditions {

// 				user_is_member_of_any_population_id = [
// 					"${pingone_environment.%[2]s.default_population_id}"
// 				]

// 				user_attribute_equals {
// 					attribute_reference = "$${user.name.given}"
// 					value 				= "Bruce"
// 				}

// 				user_attribute_equals {
// 					attribute_reference = "$${user.name.family}"
// 					value 				= "Wayne"
// 				}

// 				ip_out_of_range_cidr = [
// 					"192.168.129.23/17"
// 				]

// 				ip_reputation_high_risk = true

// 				geovelocity_anomaly_detected = true

// 				anonymous_network_detected = true

// 				anonymous_network_detected_allowed_cidr = [
// 					"192.168.129.23/17",
// 					"192.168.0.15/24"
// 				]
// 			}

// 			mfa {}

// 		}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
// }
