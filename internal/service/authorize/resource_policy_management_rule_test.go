package authorize_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/authorize"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/base"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccPolicyManagementRule_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_policy_management_rule.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	var ruleID, environmentID string

	var p1Client *client.Client
	var ctx = context.Background()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckFeatureFlag(t, acctest.ENUMFEATUREFLAG_AUTHORIZEPMTF)

			p1Client = acctest.PreCheckTestClient(ctx, t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             authorize.PolicyManagementRule_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccPolicyManagementRuleConfig_Minimal(resourceName, name),
				Check:  authorize.PolicyManagementRule_GetIDs(resourceFullName, &environmentID, &ruleID),
			},
			{
				PreConfig: func() {
					authorize.PolicyManagementRule_RemovalDrift_PreConfig(ctx, p1Client.API.AuthorizeAPIClient, t, environmentID, ruleID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the environment
			{
				Config: testAccPolicyManagementRuleConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check:  authorize.PolicyManagementRule_GetIDs(resourceFullName, &environmentID, &ruleID),
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

func TestAccPolicyManagementRule_NewEnv(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_policy_management_rule.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckFeatureFlag(t, acctest.ENUMFEATUREFLAG_AUTHORIZEPMTF)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             authorize.PolicyManagementRule_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPolicyManagementRuleConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
				),
			},
		},
	})
}

func TestAccPolicyManagementRule_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_policy_management_rule.%s", resourceName)

	name := resourceName

	fullCheck := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "name", name),
		resource.TestCheckResourceAttr(resourceFullName, "description", "Test rule full"),
		resource.TestCheckResourceAttr(resourceFullName, "enabled", "false"),
		// resource.TestCheckResourceAttr(resourceFullName, "statements.#", "1"),
		// resource.TestMatchResourceAttr(resourceFullName, "statements.0.id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "condition.type", "OR"),
		resource.TestCheckResourceAttr(resourceFullName, "condition.conditions.#", "2"),
		resource.TestCheckResourceAttr(resourceFullName, "effect_settings.type", "CONDITIONAL_PERMIT_ELSE_DENY"),
		resource.TestCheckResourceAttr(resourceFullName, "effect_settings.condition.type", "OR"),
		resource.TestCheckResourceAttr(resourceFullName, "effect_settings.condition.conditions.#", "2"),
		resource.TestMatchResourceAttr(resourceFullName, "version", verify.P1ResourceIDRegexpFullString),
	)

	minimalCheck := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "name", name),
		resource.TestCheckResourceAttr(resourceFullName, "description", "Test rule"),
		resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
		resource.TestCheckNoResourceAttr(resourceFullName, "statements"),
		resource.TestCheckResourceAttr(resourceFullName, "condition.type", "EMPTY"),
		resource.TestCheckResourceAttr(resourceFullName, "effect_settings.type", "UNCONDITIONAL_DENY"),
		resource.TestCheckNoResourceAttr(resourceFullName, "effect_settings.condition"),
		resource.TestMatchResourceAttr(resourceFullName, "version", verify.P1ResourceIDRegexpFullString),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckFeatureFlag(t, acctest.ENUMFEATUREFLAG_AUTHORIZEPMTF)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             authorize.PolicyManagementRule_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			{
				Config: testAccPolicyManagementRuleConfig_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config:  testAccPolicyManagementRuleConfig_Full(resourceName, name),
				Destroy: true,
			},
			// Minimal
			{
				Config: testAccPolicyManagementRuleConfig_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config:  testAccPolicyManagementRuleConfig_Minimal(resourceName, name),
				Destroy: true,
			},
			// Change
			{
				Config: testAccPolicyManagementRuleConfig_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config: testAccPolicyManagementRuleConfig_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config: testAccPolicyManagementRuleConfig_Full(resourceName, name),
				Check:  fullCheck,
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

func TestAccPolicyManagementRule_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_policy_management_rule.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckFeatureFlag(t, acctest.ENUMFEATUREFLAG_AUTHORIZEPMTF)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             authorize.PolicyManagementRule_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccPolicyManagementRuleConfig_Minimal(resourceName, name),
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

func testAccPolicyManagementRuleConfig_NewEnv(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_authorize_policy_management_rule" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[3]s"
  description    = "Test rule"

  effect_settings = {
    type = "UNCONDITIONAL_DENY"
  }
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccPolicyManagementRuleConfig_Full(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_authorize_trust_framework_attribute" "%[2]s-current-user-id" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-current-user-id"
  description    = "Test attribute"

  resolvers = [
    {
      type = "CURRENT_USER_ID"
    }
  ]

  value_type = {
    type = "STRING"
  }
}

resource "pingone_authorize_policy_management_rule" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Test rule full"

  enabled = false

  //   statements = [
  //     {
  //       id = pingone_authorize_policy_management_statement.%[2]s.id
  //     }
  //   ]

  condition = {
    type = "OR"

    conditions = [
      {
        type = "EMPTY"
      },
      {
        type = "NOT"

        condition = {
          type       = "COMPARISON"
          comparator = "EQUALS"

          left = {
            type = "ATTRIBUTE"
            id   = pingone_authorize_trust_framework_attribute.%[2]s-current-user-id.id
          }

          right = {
            type  = "CONSTANT"
            value = "test1"
          }
        }
      }
    ]
  }

  effect_settings = {
    type = "CONDITIONAL_PERMIT_ELSE_DENY"

    condition = {
      type = "OR"

      conditions = [
        {
          type = "EMPTY"
        },
        {
          type = "NOT"

          condition = {
            type       = "COMPARISON"
            comparator = "EQUALS"

            left = {
              type = "ATTRIBUTE"
              id   = pingone_authorize_trust_framework_attribute.%[2]s-current-user-id.id
            }

            right = {
              type  = "CONSTANT"
              value = "test1"
            }
          }
        }
      ]
    }
  }
}`, acctest.AuthorizePMTFSandboxEnvironment(), resourceName, name)
}

func testAccPolicyManagementRuleConfig_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_authorize_trust_framework_attribute" "%[2]s-current-user-id" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-current-user-id"
  description    = "Test attribute"

  resolvers = [
    {
      type = "CURRENT_USER_ID"
    }
  ]

  value_type = {
    type = "STRING"
  }
}

resource "pingone_authorize_policy_management_rule" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Test rule"

  effect_settings = {
    type = "UNCONDITIONAL_DENY"
  }
}`, acctest.AuthorizePMTFSandboxEnvironment(), resourceName, name)
}
