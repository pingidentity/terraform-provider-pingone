package sso_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	pingone "github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
)

func testAccCheckAuthenticationPolicyDestroy(s *terraform.State) error {
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
		if rs.Type != "pingone_authentication_policy" {
			continue
		}

		_, rEnv, err := apiClient.EnvironmentsApi.ReadOneEnvironment(ctx, rs.Primary.Attributes["environment_id"]).Execute()

		if rEnv.StatusCode == 404 {
			continue
		}

		if err != nil {
			return err
		}

		body, r, err := apiClient.SignOnPoliciesSignOnPoliciesApi.ReadOneSignOnPolicy(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.ID).Execute()

		if r.StatusCode == 404 {
			continue
		}

		if err != nil {
			tflog.Error(ctx, fmt.Sprintf("Error: %v", body))
			return err
		}

		return fmt.Errorf("PingOne Authentication Policy %s still exists", rs.Primary.ID)
	}

	return nil
}

func TestAccAuthenticationPolicy_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authentication_policy.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName
	description := "Test description"

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckAuthenticationPolicyDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAuthenticationPolicyConfig_Full(environmentName, licenseID, resourceName, name, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, "id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "environment_id"),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", description),
					resource.TestCheckResourceAttr(resourceFullName, "policy_action.#", "1"),
				),
			},
		},
	})
}

func TestAccAuthenticationPolicy_Minimal(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authentication_policy.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckAuthenticationPolicyDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAuthenticationPolicyConfig_Minimal(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, "id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "environment_id"),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", ""),
					resource.TestCheckResourceAttr(resourceFullName, "policy_action.#", "1"),
				),
			},
		},
	})
}

func TestAccAuthenticationPolicy_LoginAction(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authentication_policy.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckAuthenticationPolicyDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAuthenticationPolicyConfig_LoginFull(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, "id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "environment_id"),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", ""),
					resource.TestCheckResourceAttr(resourceFullName, "policy_action.#", "1"),
				),
			},
			{
				Config: testAccAuthenticationPolicyConfig_LoginMinimal1(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, "id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "environment_id"),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", ""),
					resource.TestCheckResourceAttr(resourceFullName, "policy_action.#", "1"),
				),
			},
			{
				Config: testAccAuthenticationPolicyConfig_LoginMinimal2(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, "id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "environment_id"),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", ""),
					resource.TestCheckResourceAttr(resourceFullName, "policy_action.#", "1"),
				),
			},
			{
				Config: testAccAuthenticationPolicyConfig_LoginFull(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, "id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "environment_id"),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", ""),
					resource.TestCheckResourceAttr(resourceFullName, "policy_action.#", "1"),
				),
			},
		},
	})
}

func TestAccAuthenticationPolicy_IDFirstAction(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authentication_policy.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckAuthenticationPolicyDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAuthenticationPolicyConfig_IDFirstFull(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, "id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "environment_id"),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", ""),
					resource.TestCheckResourceAttr(resourceFullName, "policy_action.#", "1"),
				),
			},
			{
				Config: testAccAuthenticationPolicyConfig_IDFirstMinimal1(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, "id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "environment_id"),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", ""),
					resource.TestCheckResourceAttr(resourceFullName, "policy_action.#", "1"),
				),
			},
			{
				Config: testAccAuthenticationPolicyConfig_IDFirstMinimal2(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, "id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "environment_id"),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", ""),
					resource.TestCheckResourceAttr(resourceFullName, "policy_action.#", "1"),
				),
			},
			{
				Config: testAccAuthenticationPolicyConfig_IDFirstFull(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, "id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "environment_id"),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", ""),
					resource.TestCheckResourceAttr(resourceFullName, "policy_action.#", "1"),
				),
			},
		},
	})
}

func TestAccAuthenticationPolicy_MFAAction(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authentication_policy.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckAuthenticationPolicyDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAuthenticationPolicyConfig_MFAFull(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, "id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "environment_id"),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", ""),
					resource.TestCheckResourceAttr(resourceFullName, "policy_action.#", "1"),
				),
			},
			{
				Config: testAccAuthenticationPolicyConfig_MFAMinimal(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, "id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "environment_id"),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", ""),
					resource.TestCheckResourceAttr(resourceFullName, "policy_action.#", "1"),
				),
			},
			{
				Config: testAccAuthenticationPolicyConfig_MFAFull(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, "id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "environment_id"),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", ""),
					resource.TestCheckResourceAttr(resourceFullName, "policy_action.#", "1"),
				),
			},
		},
	})
}

func TestAccAuthenticationPolicy_IDPAction(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authentication_policy.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckAuthenticationPolicyDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAuthenticationPolicyConfig_IDPFull(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, "id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "environment_id"),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", ""),
					resource.TestCheckResourceAttr(resourceFullName, "policy_action.#", "1"),
				),
			},
			{
				Config: testAccAuthenticationPolicyConfig_IDPMinimal(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, "id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "environment_id"),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", ""),
					resource.TestCheckResourceAttr(resourceFullName, "policy_action.#", "1"),
				),
			},
			{
				Config: testAccAuthenticationPolicyConfig_IDPFull(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, "id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "environment_id"),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", ""),
					resource.TestCheckResourceAttr(resourceFullName, "policy_action.#", "1"),
				),
			},
		},
	})
}

func TestAccAuthenticationPolicy_AgreementAction(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authentication_policy.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckAuthenticationPolicyDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAuthenticationPolicyConfig_AgreementFull(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, "id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "environment_id"),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", ""),
					resource.TestCheckResourceAttr(resourceFullName, "policy_action.#", "1"),
				),
			},
			{
				Config: testAccAuthenticationPolicyConfig_AgreementMinimal(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, "id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "environment_id"),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", ""),
					resource.TestCheckResourceAttr(resourceFullName, "policy_action.#", "1"),
				),
			},
			{
				Config: testAccAuthenticationPolicyConfig_AgreementFull(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, "id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "environment_id"),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", ""),
					resource.TestCheckResourceAttr(resourceFullName, "policy_action.#", "1"),
				),
			},
		},
	})
}

func TestAccAuthenticationPolicy_ProgressiveProfilingAction(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authentication_policy.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckAuthenticationPolicyDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAuthenticationPolicyConfig_ProgressiveProfilingFull(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, "id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "environment_id"),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", ""),
					resource.TestCheckResourceAttr(resourceFullName, "policy_action.#", "1"),
				),
			},
			{
				Config: testAccAuthenticationPolicyConfig_ProgressiveProfilingMinimal(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, "id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "environment_id"),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", ""),
					resource.TestCheckResourceAttr(resourceFullName, "policy_action.#", "1"),
				),
			},
			{
				Config: testAccAuthenticationPolicyConfig_ProgressiveProfilingFull(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, "id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "environment_id"),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", ""),
					resource.TestCheckResourceAttr(resourceFullName, "policy_action.#", "1"),
				),
			},
		},
	})
}

func TestAccAuthenticationPolicy_MultipleAction(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authentication_policy.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckAuthenticationPolicyDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAuthenticationPolicyConfig_Multiple1(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, "id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "environment_id"),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", ""),
					resource.TestCheckResourceAttr(resourceFullName, "policy_action.#", "1"),
				),
			},
			{
				Config: testAccAuthenticationPolicyConfig_Multiple2(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, "id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "environment_id"),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", ""),
					resource.TestCheckResourceAttr(resourceFullName, "policy_action.#", "1"),
				),
			},
			{
				Config: testAccAuthenticationPolicyConfig_Multiple1(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, "id"),
					resource.TestCheckResourceAttrSet(resourceFullName, "environment_id"),
					resource.TestCheckResourceAttr(resourceFullName, "name", name),
					resource.TestCheckResourceAttr(resourceFullName, "description", ""),
					resource.TestCheckResourceAttr(resourceFullName, "policy_action.#", "1"),
				),
			},
		},
	})
}

func testAccAuthenticationPolicyConfig_Full(environmentName, licenseID, resourceName, name, description string) string {
	return fmt.Sprintf(`
		resource "pingone_environment" "%[1]s" {
			name = "%[1]s"
			type = "SANDBOX"
			license_id = "%[2]s"
			default_population {}
			service {}
		}

		resource "pingone_authentication_policy" "%[3]s" {
			environment_id = "${pingone_environment.%[1]s.id}"

			name = "%[4]s"
			description = "%[5]s"

			policy_action {
				action_type = "LOGIN"
			}
		}`, environmentName, licenseID, resourceName, name, description)
}

// TODO: idp
func testAccAuthenticationPolicyConfig_LoginFull(environmentName, licenseID, resourceName, name string) string {

	return fmt.Sprintf(`
		resource "pingone_environment" "%[1]s" {
			name = "%[1]s"
			type = "SANDBOX"
			license_id = "%[2]s"
			default_population {}
			service {}
		}

		resource "pingone_authentication_policy" "%[3]s" {
			environment_id = "${pingone_environment.%[1]s.id}"

			name = "%[4]s"

			policy_action {
				action_type = "LOGIN"

				conditions {
					ip_range = "192.168.0.0/32"
					action_session_length_mins = "30"
				}

				login_options {
					confirm_identity_provider_attributes = true
					enforce_lockout_for_identity_providers = true
					recovery {
						enabled = false
					}
					registration {
						enabled = true
						external_href = "https://pingidentity.com"
						population_id = "${pingone_environment.%[1]s.default_population_id}"
					}
					// social_providers
				}

			}
		}`, environmentName, licenseID, resourceName, name)
}

func testAccAuthenticationPolicyConfig_LoginMinimal1(environmentName, licenseID, resourceName, name string) string {

	return fmt.Sprintf(`
		resource "pingone_environment" "%[1]s" {
			name = "%[1]s"
			type = "SANDBOX"
			license_id = "%[2]s"
			default_population {}
			service {}
		}

		resource "pingone_authentication_policy" "%[3]s" {
			environment_id = "${pingone_environment.%[1]s.id}"

			name = "%[4]s"

			policy_action {
				action_type = "LOGIN"
			}
		}`, environmentName, licenseID, resourceName, name)
}

func testAccAuthenticationPolicyConfig_LoginMinimal2(environmentName, licenseID, resourceName, name string) string {

	return fmt.Sprintf(`
		resource "pingone_environment" "%[1]s" {
			name = "%[1]s"
			type = "SANDBOX"
			license_id = "%[2]s"
			default_population {}
			service {}
		}

		resource "pingone_authentication_policy" "%[3]s" {
			environment_id = "${pingone_environment.%[1]s.id}"

			name = "%[4]s"

			policy_action {
				action_type = "LOGIN"

				mfa_options {}
			}
		}`, environmentName, licenseID, resourceName, name)
}

// TODO: idp
func testAccAuthenticationPolicyConfig_IDFirstFull(environmentName, licenseID, resourceName, name string) string {

	return fmt.Sprintf(`
		resource "pingone_environment" "%[1]s" {
			name = "%[1]s"
			type = "SANDBOX"
			license_id = "%[2]s"
			default_population {}
			service {}
		}

		resource "pingone_authentication_policy" "%[3]s" {
			environment_id = "${pingone_environment.%[1]s.id}"

			name = "%[4]s"

			policy_action {
				action_type = "IDENTIFIER_FIRST"

				conditions {
					ip_range = "192.168.0.0/32"
					action_session_length_mins = "30"
				}

				identifier_first_options {
					confirm_identity_provider_attributes = true
					enforce_lockout_for_identity_providers = true
					recovery {
						enabled = false
					}
					registration {
						enabled = true
						external_href = "https://www.pingidentity.com"
						population_id = "${pingone_environment.%[1]s.default_population_id}"
					}
					// social_providers
					// discovery_rule {
					// 	condition {
					// 		contains = "domain.com"
					// 		value = "value"
					// 	}
					// 	identity_provider_id = 
					// }
				}

			}
		}`, environmentName, licenseID, resourceName, name)
}

func testAccAuthenticationPolicyConfig_IDFirstMinimal1(environmentName, licenseID, resourceName, name string) string {

	return fmt.Sprintf(`
		resource "pingone_environment" "%[1]s" {
			name = "%[1]s"
			type = "SANDBOX"
			license_id = "%[2]s"
			default_population {}
			service {}
		}

		resource "pingone_authentication_policy" "%[3]s" {
			environment_id = "${pingone_environment.%[1]s.id}"

			name = "%[4]s"

			policy_action {
				action_type = "IDENTIFIER_FIRST"
			}
		}`, environmentName, licenseID, resourceName, name)
}

func testAccAuthenticationPolicyConfig_IDFirstMinimal2(environmentName, licenseID, resourceName, name string) string {

	return fmt.Sprintf(`
		resource "pingone_environment" "%[1]s" {
			name = "%[1]s"
			type = "SANDBOX"
			license_id = "%[2]s"
			default_population {}
			service {}
		}

		resource "pingone_authentication_policy" "%[3]s" {
			environment_id = "${pingone_environment.%[1]s.id}"

			name = "%[4]s"

			policy_action {
				action_type = "IDENTIFIER_FIRST"

				identifier_first_options {}

			}
		}`, environmentName, licenseID, resourceName, name)
}

// TODO: MFA device policy data source
func testAccAuthenticationPolicyConfig_MFAFull(environmentName, licenseID, resourceName, name string) string {

	return fmt.Sprintf(`
		resource "pingone_environment" "%[1]s" {
			name = "%[1]s"
			type = "SANDBOX"
			license_id = "%[2]s"
			default_population {}
			service {}
		}

		resource "pingone_authentication_policy" "%[3]s" {
			environment_id = "${pingone_environment.%[1]s.id}"

			name = "%[4]s"

			policy_action {
				action_type = "MULTI_FACTOR_AUTHENTICATION"

				conditions {
					ip_range = "192.168.0.0/32"
					action_session_length_mins = "30"
				}

				mfa_options {
					device_authentication_policy_id = data.device_policy.id
					no_device_mode = "BLOCK"
				}

			}
		}`, environmentName, licenseID, resourceName, name)
}

func testAccAuthenticationPolicyConfig_MFAMinimal(environmentName, licenseID, resourceName, name string) string {

	return fmt.Sprintf(`
		resource "pingone_environment" "%[1]s" {
			name = "%[1]s"
			type = "SANDBOX"
			license_id = "%[2]s"
			default_population {}
			service {}
		}

		resource "pingone_authentication_policy" "%[3]s" {
			environment_id = "${pingone_environment.%[1]s.id}"

			name = "%[4]s"

			policy_action {
				action_type = "MULTI_FACTOR_AUTHENTICATION"

				mfa_options {
					device_authentication_policy_id = data.device_policy.id
				}

			}
		}`, environmentName, licenseID, resourceName, name)
}

// TODO: idp
func testAccAuthenticationPolicyConfig_IDPFull(environmentName, licenseID, resourceName, name string) string {

	return fmt.Sprintf(`
		resource "pingone_environment" "%[1]s" {
			name = "%[1]s"
			type = "SANDBOX"
			license_id = "%[2]s"
			default_population {}
			service {}
		}

		resource "pingone_authentication_policy" "%[3]s" {
			environment_id = "${pingone_environment.%[1]s.id}"

			name = "%[4]s"

			policy_action {
				action_type = "IDENTITY_PROVIDER"

				conditions {
					ip_range = "192.168.0.0/32"
					action_session_length_mins = "30"
				}

				identity_provider_options {
					acr_values = "Level_3 Level_2 Level_1"
					identity_provider_id = 
					pass_user_context = true
					registration {
						enabled = true
						external_href = "https://www.pingidentity.com"
						population_id = "${pingone_environment.%[1]s.default_population_id}"
					}
				}

			}
		}`, environmentName, licenseID, resourceName, name)
}

func testAccAuthenticationPolicyConfig_IDPMinimal(environmentName, licenseID, resourceName, name string) string {

	return fmt.Sprintf(`
		resource "pingone_environment" "%[1]s" {
			name = "%[1]s"
			type = "SANDBOX"
			license_id = "%[2]s"
			default_population {}
			service {}
		}

		resource "pingone_authentication_policy" "%[3]s" {
			environment_id = "${pingone_environment.%[1]s.id}"

			name = "%[4]s"

			policy_action {
				action_type = "IDENTITY_PROVIDER"

				identity_provider_options {
					identity_provider_id = 
				}

			}
		}`, environmentName, licenseID, resourceName, name)
}

// TODO: agreements
func testAccAuthenticationPolicyConfig_AgreementFull(environmentName, licenseID, resourceName, name string) string {

	return fmt.Sprintf(`
		resource "pingone_environment" "%[1]s" {
			name = "%[1]s"
			type = "SANDBOX"
			license_id = "%[2]s"
			default_population {}
			service {}
		}

		resource "pingone_authentication_policy" "%[3]s" {
			environment_id = "${pingone_environment.%[1]s.id}"

			name = "%[4]s"

			policy_action {
				action_type = "AGREEMENT"

				conditions {
					ip_range = "192.168.0.0/32"
					action_session_length_mins = "30"
				}

				agreement_options {
					agreement_id = 
				}

			}
		}`, environmentName, licenseID, resourceName, name)
}

func testAccAuthenticationPolicyConfig_AgreementMinimal(environmentName, licenseID, resourceName, name string) string {

	return fmt.Sprintf(`
		resource "pingone_environment" "%[1]s" {
			name = "%[1]s"
			type = "SANDBOX"
			license_id = "%[2]s"
			default_population {}
			service {}
		}

		resource "pingone_authentication_policy" "%[3]s" {
			environment_id = "${pingone_environment.%[1]s.id}"

			name = "%[4]s"

			policy_action {
				action_type = "AGREEMENT"

				agreement_options {
					agreement_id = 
				}

			}
		}`, environmentName, licenseID, resourceName, name)
}

func testAccAuthenticationPolicyConfig_ProgressiveProfilingFull(environmentName, licenseID, resourceName, name string) string {

	return fmt.Sprintf(`
		resource "pingone_environment" "%[1]s" {
			name = "%[1]s"
			type = "SANDBOX"
			license_id = "%[2]s"
			default_population {}
			service {}
		}

		resource "pingone_authentication_policy" "%[3]s" {
			environment_id = "${pingone_environment.%[1]s.id}"

			name = "%[4]s"

			policy_action {
				action_type = "PROGRESSIVE_PROFILING"

				conditions {
					ip_range = "192.168.0.0/32"
					action_session_length_mins = "30"
				}

				progressive_profiling_options {
					prevent_multiple_prompts_per_flow = true
					prompt_interval_seconds = 5
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

			}
		}`, environmentName, licenseID, resourceName, name)
}

func testAccAuthenticationPolicyConfig_ProgressiveProfilingMinimal(environmentName, licenseID, resourceName, name string) string {

	return fmt.Sprintf(`
		resource "pingone_environment" "%[1]s" {
			name = "%[1]s"
			type = "SANDBOX"
			license_id = "%[2]s"
			default_population {}
			service {}
		}

		resource "pingone_authentication_policy" "%[3]s" {
			environment_id = "${pingone_environment.%[1]s.id}"

			name = "%[4]s"

			policy_action {
				action_type = "PROGRESSIVE_PROFILING"

				progressive_profiling_options {
					prevent_multiple_prompts_per_flow = true
					prompt_interval_seconds = 5
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

			}
		}`, environmentName, licenseID, resourceName, name)
}

func testAccAuthenticationPolicyConfig_Minimal(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		resource "pingone_environment" "%[1]s" {
			name = "%[1]s"
			type = "SANDBOX"
			license_id = "%[2]s"
			default_population {}
			service {}
		}

		resource "pingone_authentication_policy" "%[3]s" {
			environment_id = "${pingone_environment.%[1]s.id}"

			name = "%[4]s"

			policy_action {
				action_type = "LOGIN"
			}
		}`, environmentName, licenseID, resourceName, name)
}

func testAccAuthenticationPolicyConfig_Multiple1(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		resource "pingone_environment" "%[1]s" {
			name = "%[1]s"
			type = "SANDBOX"
			license_id = "%[2]s"
			default_population {}
			service {}
		}

		resource "pingone_authentication_policy" "%[3]s" {
			environment_id = "${pingone_environment.%[1]s.id}"

			name = "%[4]s"

			policy_action {
				action_type = "IDENTIFIER_FIRST"
			}

			policy_action {
				action_type = "LOGIN"
			}

			policy_action {
				action_type = "MFA"

				mfa_options {
					device_authentication_policy_id = data.device_policy.id
				}
			}
		}`, environmentName, licenseID, resourceName, name)
}

func testAccAuthenticationPolicyConfig_Multiple2(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		resource "pingone_environment" "%[1]s" {
			name = "%[1]s"
			type = "SANDBOX"
			license_id = "%[2]s"
			default_population {}
			service {}
		}

		resource "pingone_authentication_policy" "%[3]s" {
			environment_id = "${pingone_environment.%[1]s.id}"

			name = "%[4]s"

			policy_action {
				action_type = "LOGIN"
			}

			policy_action {
				action_type = "MFA"

				mfa_options {
					device_authentication_policy_id = data.device_policy.id
				}
			}

			policy_action {
				action_type = "PROGRESSIVE_PROFILING"

				progressive_profiling_options {
					prevent_multiple_prompts_per_flow = true
					prompt_interval_seconds = 5
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

			}
		}`, environmentName, licenseID, resourceName, name)
}
