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

		body, r, err := apiClient.SignOnPoliciesSignOnPolicyActionsApi.ReadOneSignOnPolicyAction(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["sign_on_policy_id"], rs.Primary.ID).Execute()

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
					//resource.TestCheckResourceAttr(resourceFullName, "social_provider_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "confirm_identity_provider_attributes", "false"),
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
					resource.TestCheckResourceAttr(resourceFullName, "social_provider_ids.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "confirm_identity_provider_attributes", "false"),
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
					//resource.TestCheckResourceAttr(resourceFullName, "social_provider_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "confirm_identity_provider_attributes", "false"),
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
					//resource.TestCheckResourceAttr(resourceFullName, "social_provider_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "confirm_identity_provider_attributes", "false"),
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
					//resource.TestCheckResourceAttr(resourceFullName, "social_provider_ids", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "confirm_identity_provider_attributes", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "enforce_lockout_for_identity_providers", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "identifier_first.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "identifier_first.0.recovery_enabled", "false"),
					//resource.TestCheckResourceAttr(resourceFullName, "identifier_first.0.discovery_rule.#", "1"),
				),
			},
			{
				Config: testAccSignOnPolicyActionConfig_IDFirstMinimal(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(resourceFullName, "sign_on_policy_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(resourceFullName, "social_provider_ids.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "confirm_identity_provider_attributes", "false"),
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
					//resource.TestCheckResourceAttr(resourceFullName, "social_provider_ids", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "confirm_identity_provider_attributes", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "enforce_lockout_for_identity_providers", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "identifier_first.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "identifier_first.0.recovery_enabled", "false"),
					//resource.TestCheckResourceAttr(resourceFullName, "identifier_first.0.discovery_rule.#", "1"),
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
					//resource.TestCheckResourceAttr(resourceFullName, "social_provider_ids", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "confirm_identity_provider_attributes", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "enforce_lockout_for_identity_providers", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "identifier_first.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "identifier_first.0.recovery_enabled", "false"),
					//resource.TestCheckResourceAttr(resourceFullName, "identifier_first.0.discovery_rule.#", "1"),
				),
			},
		},
	})
}

// func TestAccSignOnPolicyAction_MFAAction(t *testing.T) {
// }

// func TestAccSignOnPolicyAction_IDPAction(t *testing.T) {
// }

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

// TODO: idp
func testAccSignOnPolicyActionConfig_LoginFullNoExt(environmentName, licenseID, resourceName, name string) string {

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

			registration_local_population_id = "${pingone_environment.%[2]s.default_population_id}"

			login {
				recovery_enabled = false // we set this to false because the calculated default from the api is true
			}
		}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccSignOnPolicyActionConfig_LoginFullWithExt(environmentName, licenseID, resourceName, name string) string {

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

			registration_external_href = "https://www.pingidentity.com"

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

		resource "pingone_sign_on_policy" "%[3]s" {
			environment_id = "${pingone_environment.%[2]s.id}"

			name = "%[4]s"
		}

		resource "pingone_sign_on_policy_action" "%[3]s" {
			environment_id 			 = "${pingone_environment.%[2]s.id}"
			sign_on_policy_id = "${pingone_sign_on_policy.%[3]s.id}"

			priority = 1

			registration_external_href = "https://www.pingidentity.com"

			identifier_first {
				recovery_enabled = false // we set this to false because the calculated default from the api is true
				// discovery_rule {
				// 	condition {
				// 		contains = "domain.com"
				// 		value = "value"
				// 	}
				// 	identity_provider_id =
				// }
			}
		}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccSignOnPolicyActionConfig_IDFirstFullNoExt(environmentName, licenseID, resourceName, name string) string {

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

			registration_local_population_id = "${pingone_environment.%[2]s.default_population_id}"

			identifier_first {
				recovery_enabled = false // we set this to false because the calculated default from the api is true
				// discovery_rule {
				// 	condition {
				// 		contains = "domain.com"
				// 		value = "value"
				// 	}
				// 	identity_provider_id =
				// }
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

// TODO: idp
// func testAccSignOnPolicyActionConfig_IDPFull(environmentName, licenseID, resourceName, name string) string {
// }

// func testAccSignOnPolicyActionConfig_IDPMinimal(environmentName, licenseID, resourceName, name string) string {
// }

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
