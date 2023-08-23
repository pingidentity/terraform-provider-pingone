package sso_test

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func testAccCheckSignOnPolicyActionDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_sign_on_policy_action" {
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

		body, r, err := apiClient.SignOnPolicyActionsApi.ReadOneSignOnPolicyAction(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["sign_on_policy_id"], rs.Primary.ID).Execute()

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

		return fmt.Errorf("PingOne Sign on Policy Action %s still exists", rs.Primary.ID)
	}

	return nil
}

func testAccGetSignOnPolicyActionIDs(resourceName string, environmentID, signOnPolicyID, resourceID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource Not found: %s", resourceName)
		}

		*resourceID = rs.Primary.ID
		*signOnPolicyID = rs.Primary.Attributes["sign_on_policy_id"]
		*environmentID = rs.Primary.Attributes["environment_id"]

		return nil
	}
}

func TestAccSignOnPolicyAction_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_sign_on_policy_action.%s-3", resourceName)

	name := resourceName

	var resourceID, signOnPolicyID, environmentID string

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSignOnPolicyActionDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccSignOnPolicyActionConfig_Multiple1(resourceName, name),
				Check:  testAccGetSignOnPolicyActionIDs(resourceFullName, &environmentID, &signOnPolicyID, &resourceID),
			},
			// Replan after removal preconfig
			{
				PreConfig: func() {
					var ctx = context.Background()
					p1Client, err := acctest.TestClient(ctx)

					if err != nil {
						t.Fatalf("Failed to get API client: %v", err)
					}

					apiClient := p1Client.API.ManagementAPIClient

					if environmentID == "" || signOnPolicyID == "" || resourceID == "" {
						t.Fatalf("One of environment ID, sign-on policy ID or resource ID cannot be determined. Environment ID: %s, Sign-on policy ID: %s, Resource ID: %s", environmentID, signOnPolicyID, resourceID)
					}

					_, err = apiClient.SignOnPolicyActionsApi.DeleteSignOnPolicyAction(ctx, environmentID, signOnPolicyID, resourceID).Execute()
					if err != nil {
						t.Fatalf("Failed to delete sign-on policy action: %v", err)
					}
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccSignOnPolicyAction_LoginAction(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_sign_on_policy_action.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSignOnPolicyActionDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSignOnPolicyActionConfig_LoginFullWithExt(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "sign_on_policy_id", verify.P1ResourceIDRegexpFullString),
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
				Config: testAccSignOnPolicyActionConfig_LoginMinimal(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "sign_on_policy_id", verify.P1ResourceIDRegexpFullString),
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
				Config: testAccSignOnPolicyActionConfig_LoginFullNoExt(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "sign_on_policy_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "priority", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "registration_external_href", ""),
					resource.TestMatchResourceAttr(resourceFullName, "registration_local_population_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "registration_confirm_user_attributes", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "social_provider_ids.#", "2"),
					resource.TestCheckResourceAttr(resourceFullName, "enforce_lockout_for_identity_providers", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "login.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "login.0.recovery_enabled", "false"),
				),
			},
			{
				Config: testAccSignOnPolicyActionConfig_LoginFullWithExt(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "sign_on_policy_id", verify.P1ResourceIDRegexpFullString),
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
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["sign_on_policy_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSignOnPolicyAction_LoginAction_Gateway(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_sign_on_policy_action.%s", resourceName)

	name := resourceName

	withGateway := resource.TestStep{
		Config: testAccSignOnPolicyActionConfig_LoginFullWithNewUserProvisioning(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "login.0.new_user_provisioning.#", "1"),
			resource.TestCheckResourceAttr(resourceFullName, "login.0.new_user_provisioning.0.gateway.#", "3"),
			resource.TestMatchResourceAttr(resourceFullName, "login.0.new_user_provisioning.0.gateway.0.id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "login.0.new_user_provisioning.0.gateway.0.type", "LDAP"),
			resource.TestMatchResourceAttr(resourceFullName, "login.0.new_user_provisioning.0.gateway.0.user_type_id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "login.0.new_user_provisioning.0.gateway.1.id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "login.0.new_user_provisioning.0.gateway.1.type", "LDAP"),
			resource.TestMatchResourceAttr(resourceFullName, "login.0.new_user_provisioning.0.gateway.1.user_type_id", verify.P1ResourceIDRegexpFullString),
			resource.TestMatchResourceAttr(resourceFullName, "login.0.new_user_provisioning.0.gateway.2.id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "login.0.new_user_provisioning.0.gateway.2.type", "LDAP"),
			resource.TestMatchResourceAttr(resourceFullName, "login.0.new_user_provisioning.0.gateway.2.user_type_id", verify.P1ResourceIDRegexpFullString),
		),
	}

	withoutGateway := resource.TestStep{
		Config: testAccSignOnPolicyActionConfig_LoginMinimal(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "login.0.new_user_provisioning.#", "0"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSignOnPolicyActionDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			withGateway,
			withoutGateway,
			withGateway,
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["sign_on_policy_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
			//Errors
			{
				Config:      testAccSignOnPolicyActionConfig_LoginFullWithNewUserProvisioningWrongGateway(resourceName, name),
				ExpectError: regexp.MustCompile(`Only 'LDAP' type gateways are supported for new user provisioning.`),
			},
		},
	})
}

func TestAccSignOnPolicyAction_IDFirstAction(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_sign_on_policy_action.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSignOnPolicyActionDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSignOnPolicyActionConfig_IDFirstFullWithExt(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "sign_on_policy_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "registration_external_href", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "registration_local_population_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "registration_confirm_user_attributes", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "social_provider_ids.#", "2"),
					resource.TestCheckResourceAttr(resourceFullName, "enforce_lockout_for_identity_providers", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "identifier_first.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "identifier_first.0.recovery_enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "identifier_first.0.discovery_rule.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "identifier_first.0.discovery_rule.0.attribute_contains_text", "domain.com"),
					resource.TestMatchResourceAttr(resourceFullName, "identifier_first.0.discovery_rule.0.identity_provider_id", verify.P1ResourceIDRegexpFullString),
				),
			},
			{
				Config: testAccSignOnPolicyActionConfig_IDFirstMinimal(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "sign_on_policy_id", verify.P1ResourceIDRegexpFullString),
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
				Config: testAccSignOnPolicyActionConfig_IDFirstFullNoExt(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "sign_on_policy_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "registration_external_href", ""),
					resource.TestMatchResourceAttr(resourceFullName, "registration_local_population_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "registration_confirm_user_attributes", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "social_provider_ids.#", "2"),
					resource.TestCheckResourceAttr(resourceFullName, "enforce_lockout_for_identity_providers", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "identifier_first.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "identifier_first.0.recovery_enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "identifier_first.0.discovery_rule.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "identifier_first.0.discovery_rule.0.attribute_contains_text", "pingidentity.com"),
					resource.TestMatchResourceAttr(resourceFullName, "identifier_first.0.discovery_rule.0.identity_provider_id", verify.P1ResourceIDRegexpFullString),
				),
			},
			{
				Config: testAccSignOnPolicyActionConfig_IDFirstFullWithExt(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "sign_on_policy_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "registration_external_href", "https://www.pingidentity.com"),
					resource.TestCheckResourceAttr(resourceFullName, "registration_local_population_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "registration_confirm_user_attributes", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "social_provider_ids.#", "2"),
					resource.TestCheckResourceAttr(resourceFullName, "enforce_lockout_for_identity_providers", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "identifier_first.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "identifier_first.0.recovery_enabled", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "identifier_first.0.discovery_rule.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "identifier_first.0.discovery_rule.0.attribute_contains_text", "domain.com"),
					resource.TestMatchResourceAttr(resourceFullName, "identifier_first.0.discovery_rule.0.identity_provider_id", verify.P1ResourceIDRegexpFullString),
				),
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

						return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["sign_on_policy_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSignOnPolicyAction_MFAAction(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { t.Skipf("Test to be re-defined") }, // test to be re-defined
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSignOnPolicyActionDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccSignOnPolicyActionConfig_IDPFull(resourceName, name),
				ExpectError: regexp.MustCompile(`Not defined`),
			},
		},
	})
}

func TestAccSignOnPolicyAction_IDPAction(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_sign_on_policy_action.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSignOnPolicyActionDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSignOnPolicyActionConfig_IDPFull(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "sign_on_policy_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "registration_external_href", ""),
					resource.TestMatchResourceAttr(resourceFullName, "registration_local_population_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "registration_confirm_user_attributes", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "social_provider_ids.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "identity_provider.#", "1"),
					resource.TestMatchResourceAttr(resourceFullName, "identity_provider.0.identity_provider_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "identity_provider.0.acr_values", "MFA"),
					resource.TestCheckResourceAttr(resourceFullName, "identity_provider.0.pass_user_context", "true"),
				),
			},
			{
				Config: testAccSignOnPolicyActionConfig_IDPMinimal(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "sign_on_policy_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "registration_external_href", ""),
					resource.TestCheckResourceAttr(resourceFullName, "registration_local_population_id", ""),
					resource.TestCheckResourceAttr(resourceFullName, "registration_confirm_user_attributes", "false"),
					resource.TestCheckResourceAttr(resourceFullName, "social_provider_ids.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "identity_provider.#", "1"),
					resource.TestMatchResourceAttr(resourceFullName, "identity_provider.0.identity_provider_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "identity_provider.0.acr_values", ""),
					resource.TestCheckResourceAttr(resourceFullName, "identity_provider.0.pass_user_context", "false"),
				),
			},
			{
				Config: testAccSignOnPolicyActionConfig_IDPFull(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "sign_on_policy_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "registration_external_href", ""),
					resource.TestMatchResourceAttr(resourceFullName, "registration_local_population_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "registration_confirm_user_attributes", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "social_provider_ids.#", "0"),
					resource.TestCheckResourceAttr(resourceFullName, "identity_provider.#", "1"),
					resource.TestMatchResourceAttr(resourceFullName, "identity_provider.0.identity_provider_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "identity_provider.0.acr_values", "MFA"),
					resource.TestCheckResourceAttr(resourceFullName, "identity_provider.0.pass_user_context", "true"),
				),
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

						return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["sign_on_policy_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSignOnPolicyAction_AgreementAction(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_sign_on_policy_action.%s", resourceName)

	name := resourceName

	fullStep := resource.TestStep{
		Config: testAccSignOnPolicyActionConfig_AgreementFull(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "agreement.#", "1"),
			resource.TestMatchResourceAttr(resourceFullName, "agreement.0.agreement_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "agreement.0.show_decline_option", "false"),
		),
	}

	minimalStep := resource.TestStep{
		Config: testAccSignOnPolicyActionConfig_AgreementMinimal(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceFullName, "agreement.#", "1"),
			resource.TestMatchResourceAttr(resourceFullName, "agreement.0.agreement_id", verify.P1ResourceIDRegexpFullString),
			resource.TestCheckResourceAttr(resourceFullName, "agreement.0.show_decline_option", "true"),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSignOnPolicyActionDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			fullStep,
			{
				Config:  testAccSignOnPolicyActionConfig_AgreementFull(resourceName, name),
				Destroy: true,
			},
			// Minimal
			minimalStep,
			{
				Config:  testAccSignOnPolicyActionConfig_AgreementMinimal(resourceName, name),
				Destroy: true,
			},
			// Update
			fullStep,
			minimalStep,
			fullStep,
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["sign_on_policy_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSignOnPolicyAction_ProgressiveProfilingAction(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_sign_on_policy_action.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSignOnPolicyActionDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSignOnPolicyActionConfig_ProgressiveProfilingFull(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "sign_on_policy_id", verify.P1ResourceIDRegexpFullString),
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
				Config: testAccSignOnPolicyActionConfig_ProgressiveProfilingMinimal(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "sign_on_policy_id", verify.P1ResourceIDRegexpFullString),
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
				Config: testAccSignOnPolicyActionConfig_ProgressiveProfilingFull(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "sign_on_policy_id", verify.P1ResourceIDRegexpFullString),
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
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["sign_on_policy_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSignOnPolicyAction_PingIDAction(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_sign_on_policy_action.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckWorkforceEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSignOnPolicyActionDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSignOnPolicyActionConfig_PingID(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "sign_on_policy_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "pingid.#", "1"),
				),
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

						return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["sign_on_policy_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSignOnPolicyAction_PingIDWinLoginPasswordlessAction(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_sign_on_policy_action.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckWorkforceEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSignOnPolicyActionDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSignOnPolicyActionConfig_PingIDWinLoginPasswordless(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "sign_on_policy_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(resourceFullName, "pingid_windows_login_passwordless.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "pingid_windows_login_passwordless.0.unique_user_attribute_name", "username"),
					resource.TestCheckResourceAttr(resourceFullName, "pingid_windows_login_passwordless.0.offline_mode_enabled", "true"),
				),
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

						return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["sign_on_policy_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSignOnPolicyAction_MultipleActionChange(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_sign_on_policy_action.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSignOnPolicyActionDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSignOnPolicyActionConfig_Multiple1(resourceName, name),
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
				Config: testAccSignOnPolicyActionConfig_Multiple2(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-1", resourceFullName), "priority", "1"),
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-1", resourceFullName), "login.#", "1"),
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-2", resourceFullName), "priority", "2"),
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-2", resourceFullName), "progressive_profiling.#", "1"),
				),
			},
			{
				Config: testAccSignOnPolicyActionConfig_Multiple1(resourceName, name),
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

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSignOnPolicyActionDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSignOnPolicyActionConfig_ConditionsSignOnOlderThanSingle(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.last_sign_on_older_than_seconds", "3600"),
				),
			},
			{
				Config: testAccSignOnPolicyActionConfig_LoginMinimal(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "0"),
				),
			},
			{
				Config: testAccSignOnPolicyActionConfig_ConditionsSignOnOlderThanSingle(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.last_sign_on_older_than_seconds", "3600"),
				),
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

						return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["sign_on_policy_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSignOnPolicyAction_ConditionsMemberOfPopulation(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_sign_on_policy_action.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSignOnPolicyActionDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSignOnPolicyActionConfig_ConditionsMemberOfPopulation(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.#", "1"),
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.0.last_sign_on_older_than_seconds", "3600"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.user_is_member_of_any_population_id.#", "1"),
					resource.TestMatchResourceAttr(resourceFullName, "conditions.0.user_is_member_of_any_population_id.0", verify.P1ResourceIDRegexpFullString),
				),
			},
			{
				Config: testAccSignOnPolicyActionConfig_LoginMinimal(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "0"),
				),
			},
			{
				Config: testAccSignOnPolicyActionConfig_ConditionsMemberOfPopulation(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.#", "1"),
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.0.last_sign_on_older_than_seconds", "3600"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.user_is_member_of_any_population_id.#", "1"),
					resource.TestMatchResourceAttr(resourceFullName, "conditions.0.user_is_member_of_any_population_id.0", verify.P1ResourceIDRegexpFullString),
				),
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

						return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["sign_on_policy_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSignOnPolicyAction_ConditionsMemberOfPopulations(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_sign_on_policy_action.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSignOnPolicyActionDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSignOnPolicyActionConfig_ConditionsMemberOfPopulations(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.#", "1"),
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.0.last_sign_on_older_than_seconds", "3600"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.user_is_member_of_any_population_id.#", "3"),
					resource.TestMatchResourceAttr(resourceFullName, "conditions.0.user_is_member_of_any_population_id.0", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "conditions.0.user_is_member_of_any_population_id.1", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "conditions.0.user_is_member_of_any_population_id.2", verify.P1ResourceIDRegexpFullString),
				),
			},
			{
				Config: testAccSignOnPolicyActionConfig_LoginMinimal(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "0"),
				),
			},
			{
				Config: testAccSignOnPolicyActionConfig_ConditionsMemberOfPopulations(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.#", "1"),
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.0.last_sign_on_older_than_seconds", "3600"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.user_is_member_of_any_population_id.#", "3"),
					resource.TestMatchResourceAttr(resourceFullName, "conditions.0.user_is_member_of_any_population_id.0", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "conditions.0.user_is_member_of_any_population_id.1", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(resourceFullName, "conditions.0.user_is_member_of_any_population_id.2", verify.P1ResourceIDRegexpFullString),
				),
			},
		},
	})
}

func TestAccSignOnPolicyAction_ConditionsUserAttributeEqualsSingleString(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_sign_on_policy_action.%s", resourceName)

	name := resourceName

	attributeReference := "${user.lifecycle.status}"
	attributeValue := "ACCOUNT_OK"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSignOnPolicyActionDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSignOnPolicyActionConfig_ConditionsUserAttributeEqualsSingleString(resourceName, name, attributeReference, attributeValue),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.#", "1"),
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.0.last_sign_on_older_than_seconds", "3600"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.user_attribute_equals.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "conditions.0.user_attribute_equals.*", map[string]string{
						"attribute_reference": attributeReference,
						"value":               attributeValue,
					}),
				),
			},
			{
				Config: testAccSignOnPolicyActionConfig_LoginMinimal(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "0"),
				),
			},
			{
				Config: testAccSignOnPolicyActionConfig_ConditionsUserAttributeEqualsSingleString(resourceName, name, attributeReference, attributeValue),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.#", "1"),
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.0.last_sign_on_older_than_seconds", "3600"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.user_attribute_equals.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "conditions.0.user_attribute_equals.*", map[string]string{
						"attribute_reference": attributeReference,
						"value":               attributeValue,
					}),
				),
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

						return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["sign_on_policy_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSignOnPolicyAction_ConditionsUserAttributeEqualsSingleBool(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_sign_on_policy_action.%s", resourceName)

	name := resourceName

	attributeReference := "${user.mfaEnabled}"
	attributeValue := true

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSignOnPolicyActionDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSignOnPolicyActionConfig_ConditionsUserAttributeEqualsSingleBool(resourceName, name, attributeReference, attributeValue),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.#", "1"),
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.0.last_sign_on_older_than_seconds", "3600"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.user_attribute_equals.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "conditions.0.user_attribute_equals.*", map[string]string{
						"attribute_reference": "${user.mfaEnabled}",
						"value_boolean":       strconv.FormatBool(attributeValue),
					}),
				),
			},
			{
				Config: testAccSignOnPolicyActionConfig_LoginMinimal(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "0"),
				),
			},
			{
				Config: testAccSignOnPolicyActionConfig_ConditionsUserAttributeEqualsSingleBool(resourceName, name, attributeReference, attributeValue),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.#", "1"),
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.0.last_sign_on_older_than_seconds", "3600"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.user_attribute_equals.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "conditions.0.user_attribute_equals.*", map[string]string{
						"attribute_reference": "${user.mfaEnabled}",
						"value_boolean":       strconv.FormatBool(attributeValue),
					}),
				),
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

						return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["sign_on_policy_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSignOnPolicyAction_ConditionsUserAttributeEqualsMultiple(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_sign_on_policy_action.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSignOnPolicyActionDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSignOnPolicyActionConfig_ConditionsUserAttributeEqualsMultiple(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.#", "1"),
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.0.last_sign_on_older_than_seconds", "3600"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.user_attribute_equals.#", "4"),
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
					resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "conditions.0.user_attribute_equals.*", map[string]string{
						"attribute_reference": "${user.mfaEnabled}",
						"value_boolean":       "true",
					}),
				),
			},
			{
				Config: testAccSignOnPolicyActionConfig_LoginMinimal(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "0"),
				),
			},
			{
				Config: testAccSignOnPolicyActionConfig_ConditionsUserAttributeEqualsMultiple(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.#", "1"),
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.0.last_sign_on_older_than_seconds", "3600"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.user_attribute_equals.#", "4"),
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
					resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "conditions.0.user_attribute_equals.*", map[string]string{
						"attribute_reference": "${user.mfaEnabled}",
						"value_boolean":       "true",
					}),
				),
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

						return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["sign_on_policy_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSignOnPolicyAction_ConditionsInvalidPriority1(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSignOnPolicyActionDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccSignOnPolicyActionConfig_ConditionsUserAttributeEqualsPriority1(resourceName, name),
				ExpectError: regexp.MustCompile("Condition `user_attribute_equals` is defined cannot be set when the policy action priority is 1"),
			},
			{
				Config:      testAccSignOnPolicyActionConfig_ConditionsMemberOfPopulationsPriority1(resourceName, name),
				ExpectError: regexp.MustCompile("Condition `user_is_member_of_any_population_id` is defined cannot be set when the policy action priority is 1"),
			},
		},
	})
}

func TestAccSignOnPolicyAction_ConditionsIPOutOfRangeSingle(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_sign_on_policy_action.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { t.Skipf("Test to be re-defined") }, // test to be re-defined
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSignOnPolicyActionDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSignOnPolicyActionConfig_ConditionsIPOutOfRangeSingle(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.#", "1"),
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.0.last_sign_on_older_than_seconds", "3600"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.ip_out_of_range_cidr.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.ip_out_of_range_cidr.0", "192.168.129.23/17"),
				),
			},
			{
				Config: testAccSignOnPolicyActionConfig_LoginMinimal(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "0"),
				),
			},
			{
				Config: testAccSignOnPolicyActionConfig_ConditionsIPOutOfRangeSingle(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.#", "1"),
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.0.last_sign_on_older_than_seconds", "3600"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.ip_out_of_range_cidr.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.ip_out_of_range_cidr.0", "192.168.129.23/17"),
				),
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

						return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["sign_on_policy_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSignOnPolicyAction_ConditionsIPOutOfRangeMultiple(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_sign_on_policy_action.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { t.Skipf("Test to be re-defined") }, // test to be re-defined
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSignOnPolicyActionDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSignOnPolicyActionConfig_ConditionsIPOutOfRangeMultiple(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.#", "1"),
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.0.last_sign_on_older_than_seconds", "3600"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.ip_out_of_range_cidr.#", "2"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.ip_out_of_range_cidr.0", "192.168.129.23/17"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.ip_out_of_range_cidr.1", "192.168.0.15/24"),
				),
			},
			{
				Config: testAccSignOnPolicyActionConfig_LoginMinimal(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "0"),
				),
			},
			{
				Config: testAccSignOnPolicyActionConfig_ConditionsIPOutOfRangeMultiple(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.#", "1"),
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.0.last_sign_on_older_than_seconds", "3600"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.ip_out_of_range_cidr.#", "2"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.ip_out_of_range_cidr.0", "192.168.129.23/17"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.ip_out_of_range_cidr.1", "192.168.0.15/24"),
				),
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

						return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["sign_on_policy_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSignOnPolicyAction_ConditionsIPHighRisk(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_sign_on_policy_action.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { t.Skipf("Test to be re-defined") }, // test to be re-defined
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSignOnPolicyActionDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSignOnPolicyActionConfig_ConditionsIPHighRisk(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.#", "1"),
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.0.last_sign_on_older_than_seconds", "3600"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.ip_reputation_high_risk", "true"),
				),
			},
			{
				Config: testAccSignOnPolicyActionConfig_LoginMinimal(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "0"),
				),
			},
			{
				Config: testAccSignOnPolicyActionConfig_ConditionsIPHighRisk(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.#", "1"),
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.0.last_sign_on_older_than_seconds", "3600"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.ip_reputation_high_risk", "true"),
				),
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

						return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["sign_on_policy_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSignOnPolicyAction_ConditionsGeovelocity(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_sign_on_policy_action.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { t.Skipf("Test to be re-defined") }, // test to be re-defined
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSignOnPolicyActionDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSignOnPolicyActionConfig_ConditionsGeovelocity(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.#", "1"),
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.0.last_sign_on_older_than_seconds", "3600"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.geovelocity_anomaly_detected", "true"),
				),
			},
			{
				Config: testAccSignOnPolicyActionConfig_LoginMinimal(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "0"),
				),
			},
			{
				Config: testAccSignOnPolicyActionConfig_ConditionsGeovelocity(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.#", "1"),
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.0.last_sign_on_older_than_seconds", "3600"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.geovelocity_anomaly_detected", "true"),
				),
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

						return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["sign_on_policy_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSignOnPolicyAction_ConditionsAnonymousNetwork(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_sign_on_policy_action.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { t.Skipf("Test to be re-defined") }, // test to be re-defined
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSignOnPolicyActionDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSignOnPolicyActionConfig_ConditionsAnonymousNetwork(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.#", "1"),
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.0.last_sign_on_older_than_seconds", "3600"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.anonymous_network_detected", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.anonymous_network_detected_allowed_cidr.#", "0"),
				),
			},
			{
				Config: testAccSignOnPolicyActionConfig_LoginMinimal(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "0"),
				),
			},
			{
				Config: testAccSignOnPolicyActionConfig_ConditionsAnonymousNetwork(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.#", "1"),
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.0.last_sign_on_older_than_seconds", "3600"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.anonymous_network_detected", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.anonymous_network_detected_allowed_cidr.#", "0"),
				),
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

						return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["sign_on_policy_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSignOnPolicyAction_ConditionsAnonymousNetworkWithAllowed(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_sign_on_policy_action.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { t.Skipf("Test to be re-defined") }, // test to be re-defined
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSignOnPolicyActionDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSignOnPolicyActionConfig_ConditionsAnonymousNetworkWithAllowed(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.#", "1"),
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.0.last_sign_on_older_than_seconds", "3600"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.anonymous_network_detected", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.anonymous_network_detected_allowed_cidr.#", "2"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.anonymous_network_detected_allowed_cidr.0", "192.168.129.23/17"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.anonymous_network_detected_allowed_cidr.1", "192.168.0.15/24"),
				),
			},
			{
				Config: testAccSignOnPolicyActionConfig_LoginMinimal(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "0"),
				),
			},
			{
				Config: testAccSignOnPolicyActionConfig_ConditionsAnonymousNetworkWithAllowed(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.#", "1"),
					resource.TestCheckResourceAttr(fmt.Sprintf("%s-id", resourceFullName), "conditions.0.last_sign_on_older_than_seconds", "3600"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.anonymous_network_detected", "true"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.anonymous_network_detected_allowed_cidr.#", "2"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.anonymous_network_detected_allowed_cidr.0", "192.168.129.23/17"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.anonymous_network_detected_allowed_cidr.1", "192.168.0.15/24"),
				),
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

						return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["sign_on_policy_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSignOnPolicyAction_ConditionsCompound(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_sign_on_policy_action.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSignOnPolicyActionDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSignOnPolicyActionConfig_ConditionsCompoundSubset(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.user_is_member_of_any_population_id.#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, "conditions.0.user_attribute_equals.#", "2"),
				),
			},
			// {
			// 	Config: testAccSignOnPolicyActionConfig_ConditionsCompoundFull(resourceName, name),
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
			// 	Config: testAccSignOnPolicyActionConfig_ConditionsCompoundSubset(resourceName, name),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		resource.TestCheckResourceAttr(resourceFullName, "conditions.#", "1"),
			// 		resource.TestCheckResourceAttr(resourceFullName, "conditions.0.user_is_member_of_any_population_id.#", "1"),
			// 		resource.TestCheckResourceAttr(resourceFullName, "conditions.0.user_attribute_equals.#", "2"),
			// 	),
			// },
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["sign_on_policy_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSignOnPolicyAction_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_sign_on_policy_action.%s-3", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSignOnPolicyActionDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccSignOnPolicyActionConfig_Multiple1(resourceName, name),
			},
			// Errors
			{
				ResourceName: resourceFullName,
				ImportState:  true,
				ExpectError:  regexp.MustCompile(`Invalid import ID specified \(".*"\).  The ID should be in the format "environment_id/sign_on_policy_id/sign_on_policy_action_id" and must match regex: .*`),
			},
			{
				ResourceName:  resourceFullName,
				ImportStateId: "/",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`Invalid import ID specified \(".*"\).  The ID should be in the format "environment_id/sign_on_policy_id/sign_on_policy_action_id" and must match regex: .*`),
			},
			{
				ResourceName:  resourceFullName,
				ImportStateId: "badformat/badformat/badformat",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`Invalid import ID specified \(".*"\).  The ID should be in the format "environment_id/sign_on_policy_id/sign_on_policy_action_id" and must match regex: .*`),
			},
		},
	})
}

func testAccSignOnPolicyActionConfig_LoginFullNoExt(resourceName, name string) string {

	return fmt.Sprintf(`
		%[1]s

resource "pingone_population" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_identity_provider" "%[2]s-1" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-1"

  google {
    client_id     = "testclientid"
    client_secret = "testclientsecret"
  }
}

resource "pingone_identity_provider" "%[2]s-2" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-2"

  facebook {
    app_id     = "testclientid"
    app_secret = "testclientsecret"
  }
}

resource "pingone_sign_on_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_sign_on_policy_action" "%[2]s" {
  environment_id    = data.pingone_environment.general_test.id
  sign_on_policy_id = pingone_sign_on_policy.%[2]s.id

  priority = 1

  enforce_lockout_for_identity_providers = true
  registration_confirm_user_attributes   = true

  registration_local_population_id = pingone_population.%[2]s.id

  social_provider_ids = [
    pingone_identity_provider.%[2]s-2.id,
    pingone_identity_provider.%[2]s-1.id
  ]

  login {
    recovery_enabled = false // we set this to false because the calculated default from the api is true
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccSignOnPolicyActionConfig_LoginFullWithExt(resourceName, name string) string {

	return fmt.Sprintf(`
		%[1]s

resource "pingone_identity_provider" "%[2]s-1" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-1"

  google {
    client_id     = "testclientid"
    client_secret = "testclientsecret"
  }
}

resource "pingone_identity_provider" "%[2]s-2" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-2"

  facebook {
    app_id     = "testclientid"
    app_secret = "testclientsecret"
  }
}

resource "pingone_sign_on_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_sign_on_policy_action" "%[2]s" {
  environment_id    = data.pingone_environment.general_test.id
  sign_on_policy_id = pingone_sign_on_policy.%[2]s.id

  priority = 1

  registration_external_href = "https://www.pingidentity.com"

  social_provider_ids = [
    pingone_identity_provider.%[2]s-2.id,
    pingone_identity_provider.%[2]s-1.id
  ]

  login {
    recovery_enabled = false // we set this to false because the calculated default from the api is true
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccSignOnPolicyActionConfig_LoginFullWithNewUserProvisioning(resourceName, name string) string {

	return fmt.Sprintf(`
		%[1]s


resource "pingone_population" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
}

resource "pingone_gateway" "%[2]s-1" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-1"
  enabled        = true

  type = "LDAP"

  bind_dn       = "ou=test,dc=example,dc=com"
  bind_password = "dummyPasswordValue"

  vendor = "PingDirectory"

  servers = [
    "ds1.dummyldapservice.com:389",
    "ds3.dummyldapservice.com:389",
    "ds2.dummyldapservice.com:389",
  ]

  user_type {
    name               = "User Set 1"
    password_authority = "LDAP"
    search_base_dn     = "ou=users1,dc=example,dc=com"

    user_link_attributes = ["objectGUID", "objectSid"]

    user_migration {
      lookup_filter_pattern = "(|(uid=$${identifier})(mail=$${identifier}))"

      population_id = pingone_population.%[2]s.id

      attribute_mapping {
        name  = "username"
        value = "$${ldapAttributes.uid}"
      }

      attribute_mapping {
        name  = "email"
        value = "$${ldapAttributes.mail}"
      }
    }

    push_password_changes_to_ldap = true
  }

  user_type {
    name               = "User Set 2"
    password_authority = "PING_ONE"
    search_base_dn     = "ou=users,dc=example,dc=com"

    user_link_attributes = ["objectGUID", "dn", "objectSid"]

    user_migration {
      lookup_filter_pattern = "(|(uid=$${identifier})(mail=$${identifier}))"

      population_id = pingone_population.%[2]s.id

      attribute_mapping {
        name  = "username"
        value = "$${ldapAttributes.uid}"
      }

      attribute_mapping {
        name  = "email"
        value = "$${ldapAttributes.mail}"
      }

      attribute_mapping {
        name  = "name.family"
        value = "$${ldapAttributes.sn}"
      }
    }

    push_password_changes_to_ldap = true
  }

}

resource "pingone_gateway_credential" "%[2]s-1" {
  environment_id = data.pingone_environment.general_test.id
  gateway_id     = pingone_gateway.%[2]s-1.id
}

resource "pingone_gateway" "%[2]s-2" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-2"
  enabled        = true

  type = "LDAP"

  bind_dn       = "ou=test,dc=example,dc=com"
  bind_password = "dummyPasswordValue"

  vendor = "PingDirectory"

  servers = [
    "ds1.dummyldapservice.com:389",
    "ds3.dummyldapservice.com:389",
    "ds2.dummyldapservice.com:389",
  ]

  user_type {
    name               = "User Set 1"
    password_authority = "LDAP"
    search_base_dn     = "ou=users1,dc=example,dc=com"

    user_link_attributes = ["objectGUID", "objectSid"]

    user_migration {
      lookup_filter_pattern = "(|(uid=$${identifier})(mail=$${identifier}))"

      population_id = pingone_population.%[2]s.id

      attribute_mapping {
        name  = "username"
        value = "$${ldapAttributes.uid}"
      }

      attribute_mapping {
        name  = "email"
        value = "$${ldapAttributes.mail}"
      }
    }

    push_password_changes_to_ldap = true
  }
}

resource "pingone_gateway_credential" "%[2]s-2" {
  environment_id = data.pingone_environment.general_test.id
  gateway_id     = pingone_gateway.%[2]s-2.id
}

resource "pingone_sign_on_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_sign_on_policy_action" "%[2]s" {
  environment_id    = data.pingone_environment.general_test.id
  sign_on_policy_id = pingone_sign_on_policy.%[2]s.id

  priority = 1

  login {
    recovery_enabled = false

    new_user_provisioning {
      gateway {
        id           = pingone_gateway.%[2]s-1.id
        user_type_id = pingone_gateway.%[2]s-1.user_type.* [index(pingone_gateway.%[2]s-1.user_type[*].name, "User Set 2")].id
      }

      gateway {
        id           = pingone_gateway.%[2]s-2.id
        user_type_id = pingone_gateway.%[2]s-2.user_type.* [index(pingone_gateway.%[2]s-2.user_type[*].name, "User Set 1")].id
      }

      gateway {
        id           = pingone_gateway.%[2]s-1.id
        user_type_id = pingone_gateway.%[2]s-1.user_type.* [index(pingone_gateway.%[2]s-1.user_type[*].name, "User Set 1")].id
      }
    }
  }

  depends_on = [
    pingone_gateway_credential.%[2]s-1,
    pingone_gateway_credential.%[2]s-2,
  ]
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccSignOnPolicyActionConfig_LoginFullWithNewUserProvisioningWrongGateway(resourceName, name string) string {

	return fmt.Sprintf(`
		%[1]s

resource "pingone_gateway" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  enabled        = true

  type = "PING_FEDERATE"
}

resource "pingone_sign_on_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_sign_on_policy_action" "%[2]s" {
  environment_id    = data.pingone_environment.general_test.id
  sign_on_policy_id = pingone_sign_on_policy.%[2]s.id

  priority = 1

  login {
    recovery_enabled = false

    new_user_provisioning {
      gateway {
        id           = pingone_gateway.%[2]s.id
        user_type_id = pingone_gateway.%[2]s.id
      }
    }
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccSignOnPolicyActionConfig_LoginMinimal(resourceName, name string) string {

	return fmt.Sprintf(`
		%[1]s

resource "pingone_sign_on_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_sign_on_policy_action" "%[2]s" {
  environment_id    = data.pingone_environment.general_test.id
  sign_on_policy_id = pingone_sign_on_policy.%[2]s.id

  priority = 1

  login {}
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

// TODO: idp
func testAccSignOnPolicyActionConfig_IDFirstFullWithExt(resourceName, name string) string {

	return fmt.Sprintf(`
		%[1]s

resource "pingone_identity_provider" "%[2]s-1" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-1"

  google {
    client_id     = "testclientid"
    client_secret = "testclientsecret"
  }
}

resource "pingone_identity_provider" "%[2]s-2" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-2"

  facebook {
    app_id     = "testclientid"
    app_secret = "testclientsecret"
  }
}

resource "pingone_sign_on_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_sign_on_policy_action" "%[2]s" {
  environment_id    = data.pingone_environment.general_test.id
  sign_on_policy_id = pingone_sign_on_policy.%[2]s.id

  priority = 1

  registration_external_href = "https://www.pingidentity.com"

  social_provider_ids = [
    pingone_identity_provider.%[2]s-2.id,
    pingone_identity_provider.%[2]s-1.id
  ]

  identifier_first {
    recovery_enabled = false // we set this to false because the calculated default from the api is true
    discovery_rule {
      attribute_contains_text = "domain.com"
      identity_provider_id    = pingone_identity_provider.%[2]s-1.id
    }
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccSignOnPolicyActionConfig_IDFirstFullNoExt(resourceName, name string) string {

	return fmt.Sprintf(`
		%[1]s

resource "pingone_population" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_identity_provider" "%[2]s-1" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-1"

  google {
    client_id     = "testclientid"
    client_secret = "testclientsecret"
  }
}

resource "pingone_identity_provider" "%[2]s-2" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-2"

  facebook {
    app_id     = "testclientid"
    app_secret = "testclientsecret"
  }
}

resource "pingone_sign_on_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_sign_on_policy_action" "%[2]s" {
  environment_id    = data.pingone_environment.general_test.id
  sign_on_policy_id = pingone_sign_on_policy.%[2]s.id

  priority = 1

  registration_local_population_id = pingone_population.%[2]s.id

  enforce_lockout_for_identity_providers = true
  registration_confirm_user_attributes   = true

  social_provider_ids = [
    pingone_identity_provider.%[2]s-2.id,
    pingone_identity_provider.%[2]s-1.id
  ]

  identifier_first {
    recovery_enabled = false // we set this to false because the calculated default from the api is true
    discovery_rule {
      attribute_contains_text = "pingidentity.com"
      identity_provider_id    = pingone_identity_provider.%[2]s-2.id
    }
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccSignOnPolicyActionConfig_IDFirstMinimal(resourceName, name string) string {

	return fmt.Sprintf(`
		%[1]s

resource "pingone_sign_on_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_sign_on_policy_action" "%[2]s" {
  environment_id    = data.pingone_environment.general_test.id
  sign_on_policy_id = pingone_sign_on_policy.%[2]s.id

  priority = 1

  identifier_first {}

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

// TODO: MFA device policy data source
// func testAccSignOnPolicyActionConfig_MFAFull(resourceName, name string) string {
// }

// func testAccSignOnPolicyActionConfig_MFAMinimal(resourceName, name string) string {
// }

func testAccSignOnPolicyActionConfig_IDPFull(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_population" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_identity_provider" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  openid_connect {
    client_id     = "testclientid"
    client_secret = "testclientsecret"

    authorization_endpoint = "https://pingidentity.com/authz"
    issuer                 = "https://pingidentity.com/issuer"
    jwks_endpoint          = "https://pingidentity.com/jwks"
    scopes                 = ["openid", "profile"]
    token_endpoint         = "https://pingidentity.com/token"
  }
}

resource "pingone_sign_on_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_sign_on_policy_action" "%[2]s" {
  environment_id    = data.pingone_environment.general_test.id
  sign_on_policy_id = pingone_sign_on_policy.%[2]s.id

  priority = 1

  registration_local_population_id = pingone_population.%[2]s.id

  registration_confirm_user_attributes = true

  identity_provider {
    identity_provider_id = pingone_identity_provider.%[2]s.id

    acr_values        = "MFA"
    pass_user_context = true
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccSignOnPolicyActionConfig_IDPMinimal(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_identity_provider" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  openid_connect {
    client_id     = "testclientid"
    client_secret = "testclientsecret"

    authorization_endpoint = "https://pingidentity.com/authz"
    issuer                 = "https://pingidentity.com/issuer"
    jwks_endpoint          = "https://pingidentity.com/jwks"
    scopes                 = ["openid", "profile"]
    token_endpoint         = "https://pingidentity.com/token"
  }
}

resource "pingone_sign_on_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_sign_on_policy_action" "%[2]s" {
  environment_id    = data.pingone_environment.general_test.id
  sign_on_policy_id = pingone_sign_on_policy.%[2]s.id

  priority = 1

  identity_provider {
    identity_provider_id = pingone_identity_provider.%[2]s.id
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccSignOnPolicyActionConfig_AgreementFull(resourceName, name string) string {

	return fmt.Sprintf(`
		%[1]s

resource "pingone_agreement" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name                  = "%[3]s"
  description           = "Before the crowbar was invented, Crows would just drink at home."
  reconsent_period_days = 31

}

data "pingone_language" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  locale = "en"
}

resource "pingone_agreement_localization" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  agreement_id   = pingone_agreement.%[2]s.id
  language_id    = data.pingone_language.%[2]s.id

  display_name = "%[3]s"

  text_checkbox_accept = "Yeah"
  text_button_continue = "Move on"
  text_button_decline  = "Nah"
}

resource "pingone_agreement_localization_revision" "%[2]s" {
  environment_id            = data.pingone_environment.general_test.id
  agreement_id              = pingone_agreement.%[2]s.id
  agreement_localization_id = pingone_agreement_localization.%[2]s.id

  content_type      = "text/html"
  require_reconsent = true
  text              = <<EOT
	<h1>Test</h1>
  EOT

}

resource "pingone_agreement_localization_enable" "%[2]s" {
  environment_id            = data.pingone_environment.general_test.id
  agreement_id              = pingone_agreement.%[2]s.id
  agreement_localization_id = pingone_agreement_localization.%[2]s.id

  enabled = true

  depends_on = [
    pingone_agreement_localization_revision.%[2]s
  ]
}

resource "pingone_agreement_enable" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  agreement_id   = pingone_agreement.%[2]s.id

  enabled = true

  depends_on = [
    pingone_agreement_localization_enable.%[2]s
  ]
}

resource "pingone_sign_on_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_sign_on_policy_action" "%[2]s" {
  environment_id    = data.pingone_environment.general_test.id
  sign_on_policy_id = pingone_sign_on_policy.%[2]s.id

  priority = 1

  agreement {
    agreement_id        = pingone_agreement_enable.%[2]s.agreement_id
    show_decline_option = false
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccSignOnPolicyActionConfig_AgreementMinimal(resourceName, name string) string {

	return fmt.Sprintf(`
		%[1]s


resource "pingone_agreement" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name                  = "%[3]s"
  description           = "Before the crowbar was invented, Crows would just drink at home."
  reconsent_period_days = 31

}

data "pingone_language" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  locale = "en"
}

resource "pingone_agreement_localization" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  agreement_id   = pingone_agreement.%[2]s.id
  language_id    = data.pingone_language.%[2]s.id

  display_name = "%[3]s"

  text_checkbox_accept = "Yeah"
  text_button_continue = "Move on"
  text_button_decline  = "Nah"
}

resource "pingone_agreement_localization_revision" "%[2]s" {
  environment_id            = data.pingone_environment.general_test.id
  agreement_id              = pingone_agreement.%[2]s.id
  agreement_localization_id = pingone_agreement_localization.%[2]s.id

  content_type      = "text/html"
  require_reconsent = true
  text              = <<EOT
	<h1>Test</h1>
  EOT

}

resource "pingone_agreement_localization_enable" "%[2]s" {
  environment_id            = data.pingone_environment.general_test.id
  agreement_id              = pingone_agreement.%[2]s.id
  agreement_localization_id = pingone_agreement_localization.%[2]s.id

  enabled = true

  depends_on = [
    pingone_agreement_localization_revision.%[2]s
  ]
}

resource "pingone_agreement_enable" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  agreement_id   = pingone_agreement.%[2]s.id

  enabled = true

  depends_on = [
    pingone_agreement_localization_enable.%[2]s
  ]
}

resource "pingone_sign_on_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_sign_on_policy_action" "%[2]s" {
  environment_id    = data.pingone_environment.general_test.id
  sign_on_policy_id = pingone_sign_on_policy.%[2]s.id

  priority = 1

  agreement {
    agreement_id = pingone_agreement_enable.%[2]s.agreement_id
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccSignOnPolicyActionConfig_ProgressiveProfilingFull(resourceName, name string) string {

	return fmt.Sprintf(`
		%[1]s

resource "pingone_sign_on_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_sign_on_policy_action" "%[2]s" {
  environment_id    = data.pingone_environment.general_test.id
  sign_on_policy_id = pingone_sign_on_policy.%[2]s.id

  priority = 1

  progressive_profiling {
    prevent_multiple_prompts_per_flow = false // default is true
    prompt_interval_seconds           = 5     // default is 7776000
    prompt_text                       = "Sed ut perspiciatis unde omnis iste natus error sit voluptatem accusantium doloremque laudantium, totam rem aperiam, eaque ipsa quae ab illo inventore veritatis et quasi architecto beatae vitae dicta sunt explicabo."

    attribute {
      name     = "name.given"
      required = true
    }

    attribute {
      name     = "name.family"
      required = true
    }

    attribute {
      name     = "address.postalCode"
      required = false
    }
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccSignOnPolicyActionConfig_ProgressiveProfilingMinimal(resourceName, name string) string {

	return fmt.Sprintf(`
		%[1]s

resource "pingone_sign_on_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_sign_on_policy_action" "%[2]s" {
  environment_id    = data.pingone_environment.general_test.id
  sign_on_policy_id = pingone_sign_on_policy.%[2]s.id

  priority = 1

  progressive_profiling {
    prompt_text = "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."

    attribute {
      name     = "email"
      required = true
    }

    attribute {
      name     = "address.postalCode"
      required = false
    }
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccSignOnPolicyActionConfig_PingID(resourceName, name string) string {

	return fmt.Sprintf(`
		%[1]s

resource "pingone_sign_on_policy" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id

  name = "%[3]s"
}

resource "pingone_sign_on_policy_action" "%[2]s" {
  environment_id    = data.pingone_environment.workforce_test.id
  sign_on_policy_id = pingone_sign_on_policy.%[2]s.id

  priority = 1

  pingid {}
}`, acctest.WorkforceSandboxEnvironment(), resourceName, name)
}

func testAccSignOnPolicyActionConfig_PingIDWinLoginPasswordless(resourceName, name string) string {

	return fmt.Sprintf(`
		%[1]s

resource "pingone_sign_on_policy" "%[2]s" {
  environment_id = data.pingone_environment.workforce_test.id

  name = "%[3]s"
}

resource "pingone_sign_on_policy_action" "%[2]s" {
  environment_id    = data.pingone_environment.workforce_test.id
  sign_on_policy_id = pingone_sign_on_policy.%[2]s.id

  priority = 1

  pingid_windows_login_passwordless {
    unique_user_attribute_name = "username"
    offline_mode_enabled       = true
  }
}`, acctest.WorkforceSandboxEnvironment(), resourceName, name)
}

func testAccSignOnPolicyActionConfig_Multiple1(resourceName, name string) string {

	return fmt.Sprintf(`
		%[1]s

resource "pingone_sign_on_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_sign_on_policy_action" "%[2]s-1" {
  environment_id    = data.pingone_environment.general_test.id
  sign_on_policy_id = pingone_sign_on_policy.%[2]s.id

  priority = 1

  identifier_first {}

}

resource "pingone_sign_on_policy_action" "%[2]s-2" {
  environment_id    = data.pingone_environment.general_test.id
  sign_on_policy_id = pingone_sign_on_policy.%[2]s.id

  priority = 2

  login {}

}

resource "pingone_sign_on_policy_action" "%[2]s-3" {
  environment_id    = data.pingone_environment.general_test.id
  sign_on_policy_id = pingone_sign_on_policy.%[2]s.id

  priority = 3

  progressive_profiling {
    prompt_text = "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."

    attribute {
      name     = "email"
      required = true
    }

    attribute {
      name     = "address.postalCode"
      required = false
    }
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccSignOnPolicyActionConfig_Multiple2(resourceName, name string) string {

	return fmt.Sprintf(`
		%[1]s

resource "pingone_sign_on_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_sign_on_policy_action" "%[2]s-1" {
  environment_id    = data.pingone_environment.general_test.id
  sign_on_policy_id = pingone_sign_on_policy.%[2]s.id

  priority = 1

  login {}

}

resource "pingone_sign_on_policy_action" "%[2]s-2" {
  environment_id    = data.pingone_environment.general_test.id
  sign_on_policy_id = pingone_sign_on_policy.%[2]s.id

  priority = 2

  progressive_profiling {
    prompt_text = "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."

    attribute {
      name     = "email"
      required = true
    }

    attribute {
      name     = "address.postalCode"
      required = false
    }
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

// func testAccSignOnPolicyActionConfig_ConditionsDeclaredNotDefined(resourceName, name string) string {

// 	return fmt.Sprintf(`
// 		%[1]s

// 		resource "pingone_sign_on_policy" "%[2]s" {
// 			environment_id = data.pingone_environment.general_test.id

// 			name = "%[3]s"
// 		}

// 		resource "pingone_sign_on_policy_action" "%[2]s" {
// 			environment_id 			 = data.pingone_environment.general_test.id
// 			sign_on_policy_id = pingone_sign_on_policy.%[2]s.id

// 			priority = 1

// 			conditions {}

// 			login {}

// 		}`, acctest.GenericSandboxEnvironment(), resourceName, name)
// }

func testAccSignOnPolicyActionConfig_ConditionsSignOnOlderThanSingle(resourceName, name string) string {

	return fmt.Sprintf(`
		%[1]s

resource "pingone_sign_on_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_sign_on_policy_action" "%[2]s" {
  environment_id    = data.pingone_environment.general_test.id
  sign_on_policy_id = pingone_sign_on_policy.%[2]s.id

  priority = 1

  conditions {
    last_sign_on_older_than_seconds = 3600
  }

  login {}

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccSignOnPolicyActionConfig_ConditionsMemberOfPopulation(resourceName, name string) string {

	return fmt.Sprintf(`
		%[1]s

resource "pingone_population" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_sign_on_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_sign_on_policy_action" "%[2]s-id" {
  environment_id    = data.pingone_environment.general_test.id
  sign_on_policy_id = pingone_sign_on_policy.%[2]s.id

  priority = 1

  conditions {
    last_sign_on_older_than_seconds = 3600
  }

  identifier_first {}

}

resource "pingone_sign_on_policy_action" "%[2]s" {
  environment_id    = data.pingone_environment.general_test.id
  sign_on_policy_id = pingone_sign_on_policy.%[2]s.id

  priority = 2

  conditions {
    user_is_member_of_any_population_id = [
      pingone_population.%[2]s.id
    ]
  }

  login {}

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccSignOnPolicyActionConfig_ConditionsMemberOfPopulations(resourceName, name string) string {

	return fmt.Sprintf(`
		%[1]s

resource "pingone_population" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_population" "%[2]s-1" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s-1"
}

resource "pingone_population" "%[2]s-2" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s-2"
}

resource "pingone_sign_on_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_sign_on_policy_action" "%[2]s-id" {
  environment_id    = data.pingone_environment.general_test.id
  sign_on_policy_id = pingone_sign_on_policy.%[2]s.id

  priority = 1

  conditions {
    last_sign_on_older_than_seconds = 3600
  }

  identifier_first {}

}

resource "pingone_sign_on_policy_action" "%[2]s" {
  environment_id    = data.pingone_environment.general_test.id
  sign_on_policy_id = pingone_sign_on_policy.%[2]s.id

  priority = 2

  conditions {
    user_is_member_of_any_population_id = [
      pingone_population.%[2]s.id,
      pingone_population.%[2]s-1.id,
      pingone_population.%[2]s-2.id
    ]
  }

  login {}

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccSignOnPolicyActionConfig_ConditionsUserAttributeEqualsSingleString(resourceName, name, attributeReference, attributeValue string) string {

	return fmt.Sprintf(`
		%[1]s

resource "pingone_sign_on_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_sign_on_policy_action" "%[2]s-id" {
  environment_id    = data.pingone_environment.general_test.id
  sign_on_policy_id = pingone_sign_on_policy.%[2]s.id

  priority = 1

  conditions {
    last_sign_on_older_than_seconds = 3600
  }

  identifier_first {}

}

resource "pingone_sign_on_policy_action" "%[2]s" {
  environment_id    = data.pingone_environment.general_test.id
  sign_on_policy_id = pingone_sign_on_policy.%[2]s.id

  priority = 2

  conditions {
    user_attribute_equals {
      attribute_reference = "$%[4]s"
      value               = "%[5]s"
    }
  }

  login {}

}`, acctest.GenericSandboxEnvironment(), resourceName, name, attributeReference, attributeValue)
}

func testAccSignOnPolicyActionConfig_ConditionsUserAttributeEqualsSingleBool(resourceName, name, attributeReference string, attributeValue bool) string {

	return fmt.Sprintf(`
		%[1]s

resource "pingone_sign_on_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_sign_on_policy_action" "%[2]s-id" {
  environment_id    = data.pingone_environment.general_test.id
  sign_on_policy_id = pingone_sign_on_policy.%[2]s.id

  priority = 1

  conditions {
    last_sign_on_older_than_seconds = 3600
  }

  identifier_first {}

}

resource "pingone_sign_on_policy_action" "%[2]s" {
  environment_id    = data.pingone_environment.general_test.id
  sign_on_policy_id = pingone_sign_on_policy.%[2]s.id

  priority = 2

  conditions {
    user_attribute_equals {
      attribute_reference = "$%[4]s"
      value_boolean       = %[5]t
    }
  }

  login {}

}`, acctest.GenericSandboxEnvironment(), resourceName, name, attributeReference, attributeValue)
}

func testAccSignOnPolicyActionConfig_ConditionsUserAttributeEqualsMultiple(resourceName, name string) string {

	return fmt.Sprintf(`
		%[1]s

resource "pingone_sign_on_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_sign_on_policy_action" "%[2]s-id" {
  environment_id    = data.pingone_environment.general_test.id
  sign_on_policy_id = pingone_sign_on_policy.%[2]s.id

  priority = 1

  conditions {
    last_sign_on_older_than_seconds = 3600
  }

  identifier_first {}

}

resource "pingone_sign_on_policy_action" "%[2]s" {
  environment_id    = data.pingone_environment.general_test.id
  sign_on_policy_id = pingone_sign_on_policy.%[2]s.id

  priority = 2

  conditions {
    user_attribute_equals {
      attribute_reference = "$${user.name.given}"
      value               = "Bruce"
    }

    user_attribute_equals {
      attribute_reference = "$${user.name.family}"
      value               = "Wayne"
    }

    user_attribute_equals {
      attribute_reference = "$${user.lifecycle.status}"
      value               = "ACCOUNT_OK"
    }

    user_attribute_equals {
      attribute_reference = "$${user.mfaEnabled}"
      value_boolean       = true
    }
  }

  login {}

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccSignOnPolicyActionConfig_ConditionsUserAttributeEqualsPriority1(resourceName, name string) string {

	return fmt.Sprintf(`
		%[1]s

resource "pingone_sign_on_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_sign_on_policy_action" "%[2]s" {
  environment_id    = data.pingone_environment.general_test.id
  sign_on_policy_id = pingone_sign_on_policy.%[2]s.id

  priority = 1

  conditions {
    user_attribute_equals {
      attribute_reference = "$${user.mfaEnabled}"
      value_boolean       = true
    }
  }

  login {}

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccSignOnPolicyActionConfig_ConditionsMemberOfPopulationsPriority1(resourceName, name string) string {

	return fmt.Sprintf(`
		%[1]s

resource "pingone_population" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_population" "%[2]s-1" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s-1"
}

resource "pingone_population" "%[2]s-2" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s-2"
}

resource "pingone_sign_on_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_sign_on_policy_action" "%[2]s" {
  environment_id    = data.pingone_environment.general_test.id
  sign_on_policy_id = pingone_sign_on_policy.%[2]s.id

  priority = 1

  conditions {
    user_is_member_of_any_population_id = [
      pingone_population.%[2]s.id,
      pingone_population.%[2]s-1.id,
      pingone_population.%[2]s-2.id
    ]
  }

  login {}

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccSignOnPolicyActionConfig_ConditionsIPOutOfRangeSingle(resourceName, name string) string {

	return fmt.Sprintf(`
		%[1]s

resource "pingone_sign_on_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_sign_on_policy_action" "%[2]s-id" {
  environment_id    = data.pingone_environment.general_test.id
  sign_on_policy_id = pingone_sign_on_policy.%[2]s.id

  priority = 1

  conditions {
    last_sign_on_older_than_seconds = 3600
  }

  identifier_first {}

}

resource "pingone_sign_on_policy_action" "%[2]s" {
  environment_id    = data.pingone_environment.general_test.id
  sign_on_policy_id = pingone_sign_on_policy.%[2]s.id

  priority = 2

  conditions {
    ip_out_of_range_cidr = [
      "192.168.129.23/17"
    ]
  }

  mfa {}

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccSignOnPolicyActionConfig_ConditionsIPOutOfRangeMultiple(resourceName, name string) string {

	return fmt.Sprintf(`
		%[1]s

resource "pingone_sign_on_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_sign_on_policy_action" "%[2]s-id" {
  environment_id    = data.pingone_environment.general_test.id
  sign_on_policy_id = pingone_sign_on_policy.%[2]s.id

  priority = 1

  conditions {
    last_sign_on_older_than_seconds = 3600
  }

  identifier_first {}

}

resource "pingone_sign_on_policy_action" "%[2]s" {
  environment_id    = data.pingone_environment.general_test.id
  sign_on_policy_id = pingone_sign_on_policy.%[2]s.id

  priority = 2

  conditions {
    ip_out_of_range_cidr = [
      "192.168.129.23/17",
      "192.168.0.15/24"
    ]
  }

  mfa {}

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccSignOnPolicyActionConfig_ConditionsIPHighRisk(resourceName, name string) string {

	return fmt.Sprintf(`
		%[1]s

resource "pingone_sign_on_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_sign_on_policy_action" "%[2]s-id" {
  environment_id    = data.pingone_environment.general_test.id
  sign_on_policy_id = pingone_sign_on_policy.%[2]s.id

  priority = 1

  conditions {
    last_sign_on_older_than_seconds = 3600
  }

  identifier_first {}

}

resource "pingone_sign_on_policy_action" "%[2]s" {
  environment_id    = data.pingone_environment.general_test.id
  sign_on_policy_id = pingone_sign_on_policy.%[2]s.id

  priority = 2

  conditions {
    ip_reputation_high_risk = true
  }

  mfa {}

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccSignOnPolicyActionConfig_ConditionsGeovelocity(resourceName, name string) string {

	return fmt.Sprintf(`
		%[1]s

resource "pingone_sign_on_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_sign_on_policy_action" "%[2]s-id" {
  environment_id    = data.pingone_environment.general_test.id
  sign_on_policy_id = pingone_sign_on_policy.%[2]s.id

  priority = 1

  conditions {
    last_sign_on_older_than_seconds = 3600
  }

  identifier_first {}

}

resource "pingone_sign_on_policy_action" "%[2]s" {
  environment_id    = data.pingone_environment.general_test.id
  sign_on_policy_id = pingone_sign_on_policy.%[2]s.id

  priority = 2

  conditions {
    geovelocity_anomaly_detected = true
  }

  mfa {}

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccSignOnPolicyActionConfig_ConditionsAnonymousNetwork(resourceName, name string) string {

	return fmt.Sprintf(`
		%[1]s

resource "pingone_sign_on_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_sign_on_policy_action" "%[2]s-id" {
  environment_id    = data.pingone_environment.general_test.id
  sign_on_policy_id = pingone_sign_on_policy.%[2]s.id

  priority = 1

  conditions {
    last_sign_on_older_than_seconds = 3600
  }

  identifier_first {}

}

resource "pingone_sign_on_policy_action" "%[2]s" {
  environment_id    = data.pingone_environment.general_test.id
  sign_on_policy_id = pingone_sign_on_policy.%[2]s.id

  priority = 2

  conditions {
    anonymous_network_detected = true
  }

  mfa {}

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccSignOnPolicyActionConfig_ConditionsAnonymousNetworkWithAllowed(resourceName, name string) string {

	return fmt.Sprintf(`
		%[1]s

resource "pingone_sign_on_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_sign_on_policy_action" "%[2]s-id" {
  environment_id    = data.pingone_environment.general_test.id
  sign_on_policy_id = pingone_sign_on_policy.%[2]s.id

  priority = 1

  conditions {
    last_sign_on_older_than_seconds = 3600
  }

  identifier_first {}

}

resource "pingone_sign_on_policy_action" "%[2]s" {
  environment_id    = data.pingone_environment.general_test.id
  sign_on_policy_id = pingone_sign_on_policy.%[2]s.id

  priority = 2

  conditions {
    anonymous_network_detected = true

    anonymous_network_detected_allowed_cidr = [
      "192.168.129.23/17",
      "192.168.0.15/24"
    ]
  }

  mfa {}

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccSignOnPolicyActionConfig_ConditionsCompoundSubset(resourceName, name string) string {

	return fmt.Sprintf(`
		%[1]s

resource "pingone_population" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_sign_on_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"
}

resource "pingone_sign_on_policy_action" "%[2]s-id" {
  environment_id    = data.pingone_environment.general_test.id
  sign_on_policy_id = pingone_sign_on_policy.%[2]s.id

  priority = 1

  conditions {
    last_sign_on_older_than_seconds = 3600
  }

  identifier_first {}

}

resource "pingone_sign_on_policy_action" "%[2]s" {
  environment_id    = data.pingone_environment.general_test.id
  sign_on_policy_id = pingone_sign_on_policy.%[2]s.id

  priority = 2

  conditions {

    user_is_member_of_any_population_id = [
      pingone_population.%[2]s.id
    ]

    user_attribute_equals {
      attribute_reference = "$${user.name.given}"
      value               = "Bruce"
    }

    user_attribute_equals {
      attribute_reference = "$${user.name.family}"
      value               = "Wayne"
    }

  }

  login {}

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

// func testAccSignOnPolicyActionConfig_ConditionsCompoundFull(resourceName, name string) string {

// 	return fmt.Sprintf(`
// 		%[1]s

// 		resource "pingone_sign_on_policy" "%[2]s" {
// 			environment_id = data.pingone_environment.general_test.id

// 			name = "%[3]s"
// 		}

// 		resource "pingone_sign_on_policy_action" "%[2]s-id" {
// 			environment_id 			 = data.pingone_environment.general_test.id
// 			sign_on_policy_id = pingone_sign_on_policy.%[2]s.id

// 			priority = 1

// 			conditions {
// 				last_sign_on_older_than_seconds = 3600
// 			}

// 			identifier_first {}

// 		}

// 		resource "pingone_sign_on_policy_action" "%[2]s-login" {
// 			environment_id 			 = data.pingone_environment.general_test.id
// 			sign_on_policy_id = pingone_sign_on_policy.%[2]s.id

// 			priority = 2

// 			conditions {

// 				user_is_member_of_any_population_id = [
// 					pingone_population.%[2]s.id
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

// 		resource "pingone_sign_on_policy_action" "%[2]s-login" {
// 			environment_id 			 = data.pingone_environment.general_test.id
// 			sign_on_policy_id = pingone_sign_on_policy.%[2]s.id

// 			priority = 3

// 			conditions {

// 				user_is_member_of_any_population_id = [
// 					pingone_population.%[2]s.id
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

// 		}`, acctest.GenericSandboxEnvironment(), resourceName, name)
// }
