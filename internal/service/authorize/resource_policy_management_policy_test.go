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

func TestAccPolicyManagementPolicy_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_policy_management_policy.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	var policyID, environmentID string

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
		CheckDestroy:             authorize.PolicyManagementPolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccPolicyManagementPolicyConfig_Minimal(resourceName, name),
				Check:  authorize.PolicyManagementPolicy_GetIDs(resourceFullName, &environmentID, &policyID),
			},
			{
				PreConfig: func() {
					authorize.PolicyManagementPolicy_RemovalDrift_PreConfig(ctx, p1Client.API.AuthorizeAPIClient, t, environmentID, policyID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the environment
			{
				Config: testAccPolicyManagementPolicyConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check:  authorize.PolicyManagementPolicy_GetIDs(resourceFullName, &environmentID, &policyID),
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

func TestAccPolicyManagementPolicy_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_policy_management_policy.%s", resourceName)

	name := resourceName

	fullCheck1 := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "name", name),
		resource.TestCheckResourceAttr(resourceFullName, "description", "Test policy full"),
		resource.TestCheckResourceAttr(resourceFullName, "enabled", "false"),
		// resource.TestCheckResourceAttr(resourceFullName, "statements.#", "2"),
		resource.TestCheckResourceAttr(resourceFullName, "condition.type", "OR"),
		resource.TestCheckResourceAttr(resourceFullName, "condition.conditions.#", "2"),
		resource.TestCheckResourceAttr(resourceFullName, "combining_algorithm.algorithm", "FIRST_APPLICABLE"),
		resource.TestCheckResourceAttr(resourceFullName, "children.1.name", "Child 1"),
		resource.TestCheckNoResourceAttr(resourceFullName, "children.1.description"),
		resource.TestCheckResourceAttr(resourceFullName, "children.1.enabled", "true"),
		// resource.TestCheckResourceAttr(resourceFullName, "children.1.statements.#", "2"),
		resource.TestCheckNoResourceAttr(resourceFullName, "children.1.condition"),
		resource.TestCheckResourceAttr(resourceFullName, "children.1.combining_algorithm.algorithm", "DENY_UNLESS_PERMIT"),
		resource.TestCheckNoResourceAttr(resourceFullName, "children.1.children"),
		resource.TestCheckNoResourceAttr(resourceFullName, "children.1.repetition_settings"),
		resource.TestCheckResourceAttr(resourceFullName, "children.2.name", "Child 2"),
		resource.TestCheckResourceAttr(resourceFullName, "children.2.description", "Child 2 description"),
		resource.TestCheckResourceAttr(resourceFullName, "children.2.enabled", "false"),
		// resource.TestCheckResourceAttr(resourceFullName, "children.2.statements.#", "2"),
		resource.TestCheckResourceAttr(resourceFullName, "children.2.condition.type", "OR"),
		resource.TestCheckResourceAttr(resourceFullName, "children.2.condition.conditions.#", "2"),
		resource.TestCheckResourceAttr(resourceFullName, "children.2.combining_algorithm.algorithm", "PERMIT_OVERRIDES"),
		resource.TestCheckNoResourceAttr(resourceFullName, "children.2.children"),
		resource.TestMatchResourceAttr(resourceFullName, "children.2.repetition_settings.source.id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "children.2.repetition_settings.decision", "PERMIT"),
		resource.TestCheckResourceAttr(resourceFullName, "children.3.name", "Child 3"),
		resource.TestCheckNoResourceAttr(resourceFullName, "children.3.description"),
		resource.TestCheckResourceAttr(resourceFullName, "children.3.enabled", "true"),
		// resource.TestCheckResourceAttr(resourceFullName, "children.3.statements.#", "2"),
		resource.TestCheckNoResourceAttr(resourceFullName, "children.3.condition"),
		resource.TestCheckResourceAttr(resourceFullName, "children.3.combining_algorithm.algorithm", "DENY_UNLESS_PERMIT"),
		resource.TestCheckResourceAttr(resourceFullName, "children.3.children.1.name", "Child-Child 1"),
		resource.TestCheckNoResourceAttr(resourceFullName, "children.3.children.1.description"),
		resource.TestCheckResourceAttr(resourceFullName, "children.3.children.1.enabled", "true"),
		// resource.TestCheckResourceAttr(resourceFullName, "children.3.children.1.statements.#", "2"),
		resource.TestCheckNoResourceAttr(resourceFullName, "children.3.children.1.condition"),
		resource.TestCheckResourceAttr(resourceFullName, "children.3.children.1.combining_algorithm.algorithm", "DENY_UNLESS_PERMIT"),
		resource.TestCheckNoResourceAttr(resourceFullName, "children.3.children.1.children"),
		resource.TestCheckNoResourceAttr(resourceFullName, "children.3.children.1.repetition_settings"),
		resource.TestCheckResourceAttr(resourceFullName, "children.3.children.2.name", "Child-Child 2"),
		resource.TestCheckResourceAttr(resourceFullName, "children.3.children.2.description", "Child 2 description"),
		resource.TestCheckResourceAttr(resourceFullName, "children.3.children.2.enabled", "false"),
		// resource.TestCheckResourceAttr(resourceFullName, "children.3.children.2.statements.#", "2"),
		resource.TestCheckResourceAttr(resourceFullName, "children.3.children.2.condition.type", "OR"),
		resource.TestCheckResourceAttr(resourceFullName, "children.3.children.2.condition.conditions.#", "2"),
		resource.TestCheckResourceAttr(resourceFullName, "children.3.children.2.combining_algorithm.algorithm", "PERMIT_OVERRIDES"),
		resource.TestCheckNoResourceAttr(resourceFullName, "children.3.children.2.children"),
		resource.TestMatchResourceAttr(resourceFullName, "children.3.children.2.repetition_settings.source.id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "children.3.children.2.repetition_settings.decision", "PERMIT"),
		resource.TestMatchResourceAttr(resourceFullName, "version", verify.P1ResourceIDRegexpFullString),
	)

	fullCheck2 := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "children.1.name", "Child 2"),
		resource.TestCheckResourceAttr(resourceFullName, "children.1.description", "Child 2 description"),
		resource.TestCheckResourceAttr(resourceFullName, "children.1.enabled", "false"),
		resource.TestCheckResourceAttr(resourceFullName, "children.1.condition.type", "OR"),
		resource.TestCheckResourceAttr(resourceFullName, "children.1.condition.conditions.#", "2"),
		resource.TestCheckResourceAttr(resourceFullName, "children.1.combining_algorithm.algorithm", "PERMIT_OVERRIDES"),
		resource.TestCheckNoResourceAttr(resourceFullName, "children.1.children"),
		resource.TestMatchResourceAttr(resourceFullName, "children.1.repetition_settings.source.id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "children.1.repetition_settings.decision", "PERMIT"),
		resource.TestCheckResourceAttr(resourceFullName, "children.2.name", "Child 1"),
		resource.TestCheckNoResourceAttr(resourceFullName, "children.2.description"),
		resource.TestCheckResourceAttr(resourceFullName, "children.2.enabled", "true"),
		resource.TestCheckNoResourceAttr(resourceFullName, "children.2.condition"),
		resource.TestCheckResourceAttr(resourceFullName, "children.2.combining_algorithm.algorithm", "DENY_UNLESS_PERMIT"),
		resource.TestCheckNoResourceAttr(resourceFullName, "children.2.children"),
		resource.TestCheckNoResourceAttr(resourceFullName, "children.2.repetition_settings"),
	)

	minimalCheck := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "name", "Policies"),
		resource.TestCheckNoResourceAttr(resourceFullName, "description"),
		resource.TestCheckResourceAttr(resourceFullName, "enabled", "true"),
		resource.TestCheckNoResourceAttr(resourceFullName, "statements"),
		resource.TestCheckNoResourceAttr(resourceFullName, "condition"),
		resource.TestCheckResourceAttr(resourceFullName, "combining_algorithm.algorithm", "PERMIT_OVERRIDES"),
		resource.TestCheckNoResourceAttr(resourceFullName, "children"),
		resource.TestMatchResourceAttr(resourceFullName, "version", verify.P1ResourceIDRegexpFullString),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckFeatureFlag(t, acctest.ENUMFEATUREFLAG_AUTHORIZEPMTF)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             authorize.PolicyManagementPolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			{
				Config: testAccPolicyManagementPolicyConfig_Full1(resourceName, name),
				Check:  fullCheck1,
			},
			{
				Config:  testAccPolicyManagementPolicyConfig_Full1(resourceName, name),
				Destroy: true,
			},
			// Minimal
			{
				Config: testAccPolicyManagementPolicyConfig_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config:  testAccPolicyManagementPolicyConfig_Minimal(resourceName, name),
				Destroy: true,
			},
			// Change
			{
				Config: testAccPolicyManagementPolicyConfig_Full1(resourceName, name),
				Check:  fullCheck1,
			},
			{
				Config: testAccPolicyManagementPolicyConfig_Full2(resourceName, name),
				Check:  fullCheck2,
			},
			{
				Config: testAccPolicyManagementPolicyConfig_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config: testAccPolicyManagementPolicyConfig_Full1(resourceName, name),
				Check:  fullCheck1,
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

func TestAccPolicyManagementPolicy_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_policy_management_policy.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckFeatureFlag(t, acctest.ENUMFEATUREFLAG_AUTHORIZEPMTF)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             authorize.PolicyManagementPolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccPolicyManagementPolicyConfig_Minimal(resourceName, name),
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

func testAccPolicyManagementPolicyConfig_NewEnv(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_authorize_policy_management_policy" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[3]s"
  description    = "Test policy"

  combining_algorithm = {
    algorithm = "PERMIT_OVERRIDES"
  }
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccPolicyManagementPolicyConfig_Full1(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_authorize_trust_framework_attribute" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Test attribute"

  value_type = {
    type = "COLLECTION"
  }
}

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

resource "pingone_authorize_policy_management_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Test policy full"

  enabled = false

  //   statements = []

  condition = {
    type = "OR"

    conditions = [
      {
        type       = "COMPARISON"
        comparator = "EQUALS"

        left = {
          type = "ATTRIBUTE"
          id   = pingone_authorize_trust_framework_attribute.%[2]s-current-user-id.id
        }

        right = {
          type  = "CONSTANT"
          value = "test2"
        }
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

  combining_algorithm = {
    algorithm = "FIRST_APPLICABLE"
  }

  children = [
    {
      name = "Child 1"

      combining_algorithm = {
        algorithm = "DENY_UNLESS_PERMIT"
      }
    },
    {
      name        = "Child 2"
      description = "Child 2 description"
      enabled     = false

      condition = {
        type = "OR"

        conditions = [
          {
            type       = "COMPARISON"
            comparator = "EQUALS"

            left = {
              type = "ATTRIBUTE"
              id   = pingone_authorize_trust_framework_attribute.%[2]s-current-user-id.id
            }

            right = {
              type  = "CONSTANT"
              value = "test2"
            }
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

      combining_algorithm = {
        algorithm = "PERMIT_OVERRIDES"
      }

      repetition_settings = {
        source = {
          id = pingone_authorize_trust_framework_attribute.%[2]s.id
        }
        decision = "PERMIT"
      }
    },
    {
      name = "Child 3"

      combining_algorithm = {
        algorithm = "DENY_UNLESS_PERMIT"
      }

      children = [
        {
          name = "Child-Child 1"

          combining_algorithm = {
            algorithm = "DENY_UNLESS_PERMIT"
          }
        },
        {
          name        = "Child-Child 2"
          description = "Child 2 description"
          enabled     = false

          condition = {
            type = "OR"

            conditions = [
              {
                type       = "COMPARISON"
                comparator = "EQUALS"

                left = {
                  type = "ATTRIBUTE"
                  id   = pingone_authorize_trust_framework_attribute.%[2]s-current-user-id.id
                }

                right = {
                  type  = "CONSTANT"
                  value = "test2"
                }
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

          combining_algorithm = {
            algorithm = "PERMIT_OVERRIDES"
          }

          repetition_settings = {
            source = {
              id = pingone_authorize_trust_framework_attribute.%[2]s.id
            }
            decision = "PERMIT"
          }
        },
      ]
    },
  ]


}`, acctest.AuthorizePMTFSandboxEnvironment(), resourceName, name)
}

func testAccPolicyManagementPolicyConfig_Full2(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_authorize_trust_framework_attribute" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Test attribute"

  value_type = {
    type = "COLLECTION"
  }
}

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

resource "pingone_authorize_policy_management_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Test policy full"

  enabled = false

  //   statements = []

  condition = {
    type = "OR"

    conditions = [
      {
        type       = "COMPARISON"
        comparator = "EQUALS"

        left = {
          type = "ATTRIBUTE"
          id   = pingone_authorize_trust_framework_attribute.%[2]s-current-user-id.id
        }

        right = {
          type  = "CONSTANT"
          value = "test2"
        }
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

  combining_algorithm = {
    algorithm = "FIRST_APPLICABLE"
  }

  children = [
    {
      name        = "Child 2"
      description = "Child 2 description"
      enabled     = false

      condition = {
        type = "OR"

        conditions = [
          {
            type       = "COMPARISON"
            comparator = "EQUALS"

            left = {
              type = "ATTRIBUTE"
              id   = pingone_authorize_trust_framework_attribute.%[2]s-current-user-id.id
            }

            right = {
              type  = "CONSTANT"
              value = "test2"
            }
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

      combining_algorithm = {
        algorithm = "PERMIT_OVERRIDES"
      }

      repetition_settings = {
        source = {
          id = pingone_authorize_trust_framework_attribute.%[2]s.id
        }
        decision = "PERMIT"
      }
    },
    {
      name = "Child 1"

      combining_algorithm = {
        algorithm = "DENY_UNLESS_PERMIT"
      }
    },
  ]


}`, acctest.AuthorizePMTFSandboxEnvironment(), resourceName, name)
}

func testAccPolicyManagementPolicyConfig_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_authorize_trust_framework_attribute" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Test attribute"

  value_type = {
    type = "COLLECTION"
  }
}

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

resource "pingone_authorize_policy_management_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  combining_algorithm = {
    algorithm = "PERMIT_OVERRIDES"
  }
}`, acctest.AuthorizePMTFSandboxEnvironment(), resourceName, name)
}
