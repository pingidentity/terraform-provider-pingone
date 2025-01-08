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

func TestAccTrustFrameworkCondition_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_trust_framework_condition.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	var conditionID, environmentID string

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
		CheckDestroy:             authorize.TrustFrameworkCondition_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccTrustFrameworkConditionConfig_Minimal(resourceName, name),
				Check:  authorize.TrustFrameworkCondition_GetIDs(resourceFullName, &environmentID, &conditionID),
			},
			{
				PreConfig: func() {
					authorize.TrustFrameworkCondition_RemovalDrift_PreConfig(ctx, p1Client.API.AuthorizeAPIClient, t, environmentID, conditionID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the environment
			{
				Config: testAccTrustFrameworkConditionConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check:  authorize.TrustFrameworkCondition_GetIDs(resourceFullName, &environmentID, &conditionID),
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

func TestAccTrustFrameworkCondition_NewEnv(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_trust_framework_condition.%s", resourceName)

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
		CheckDestroy:             authorize.TrustFrameworkCondition_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccTrustFrameworkConditionConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
				),
			},
		},
	})
}

func TestAccTrustFrameworkCondition_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_trust_framework_condition.%s", resourceName)

	name := resourceName

	fullCheck := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "name", name),
		resource.TestCheckResourceAttr(resourceFullName, "description", "Test application role full"),
		resource.TestCheckResourceAttr(resourceFullName, "full_name", fmt.Sprintf("%[1]s-parent.%[1]s", name)),
		resource.TestMatchResourceAttr(resourceFullName, "parent.id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "type", "CONDITION"),
		resource.TestMatchResourceAttr(resourceFullName, "version", verify.P1ResourceIDRegexpFullString),
	)

	minimalCheck := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "name", name),
		resource.TestCheckResourceAttr(resourceFullName, "description", "Test application role"),
		resource.TestCheckResourceAttr(resourceFullName, "full_name", name),
		resource.TestCheckNoResourceAttr(resourceFullName, "parent"),
		resource.TestCheckResourceAttr(resourceFullName, "type", "CONDITION"),
		resource.TestMatchResourceAttr(resourceFullName, "version", verify.P1ResourceIDRegexpFullString),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckFeatureFlag(t, acctest.ENUMFEATUREFLAG_AUTHORIZEPMTF)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             authorize.TrustFrameworkCondition_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			{
				Config: testAccTrustFrameworkConditionConfig_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config:  testAccTrustFrameworkConditionConfig_Full(resourceName, name),
				Destroy: true,
			},
			// Minimal
			{
				Config: testAccTrustFrameworkConditionConfig_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config:  testAccTrustFrameworkConditionConfig_Minimal(resourceName, name),
				Destroy: true,
			},
			// Change
			{
				Config: testAccTrustFrameworkConditionConfig_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config: testAccTrustFrameworkConditionConfig_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config: testAccTrustFrameworkConditionConfig_Full(resourceName, name),
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

func TestAccTrustFrameworkCondition_ConditionType_And(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_trust_framework_condition.%s", resourceName)

	name := resourceName

	typeCheck1 := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "condition.type", "AND"),
		resource.TestCheckNoResourceAttr(resourceFullName, "condition.comparator"),
		resource.TestCheckNoResourceAttr(resourceFullName, "condition.condition"),
		resource.TestCheckResourceAttr(resourceFullName, "condition.conditions.#", "3"),
		resource.TestMatchTypeSetElemNestedAttrs(resourceFullName, "condition.conditions.*", map[string]*regexp.Regexp{
			"type":        regexp.MustCompile("^COMPARISON$"),
			"comparator":  regexp.MustCompile("^EQUALS$"),
			"left.type":   regexp.MustCompile("^ATTRIBUTE$"),
			"left.id":     verify.P1ResourceIDRegexpFullString,
			"right.type":  regexp.MustCompile("^CONSTANT$"),
			"right.value": regexp.MustCompile("^test2$"),
		}),
		resource.TestMatchTypeSetElemNestedAttrs(resourceFullName, "condition.conditions.*", map[string]*regexp.Regexp{
			"type":        regexp.MustCompile("^COMPARISON$"),
			"comparator":  regexp.MustCompile("^EQUALS$"),
			"left.type":   regexp.MustCompile("^ATTRIBUTE$"),
			"left.id":     verify.P1ResourceIDRegexpFullString,
			"right.type":  regexp.MustCompile("^CONSTANT$"),
			"right.value": regexp.MustCompile("^test1$"),
		}),
		resource.TestMatchTypeSetElemNestedAttrs(resourceFullName, "condition.conditions.*", map[string]*regexp.Regexp{
			"type":                  regexp.MustCompile("^NOT$"),
			"condition.type":        regexp.MustCompile("^COMPARISON$"),
			"condition.comparator":  regexp.MustCompile("^EQUALS$"),
			"condition.left.type":   regexp.MustCompile("^ATTRIBUTE$"),
			"condition.left.id":     verify.P1ResourceIDRegexpFullString,
			"condition.right.type":  regexp.MustCompile("^CONSTANT$"),
			"condition.right.value": regexp.MustCompile("^test3$"),
		}),
		resource.TestCheckNoResourceAttr(resourceFullName, "condition.left"),
		resource.TestCheckNoResourceAttr(resourceFullName, "condition.reference"),
		resource.TestCheckNoResourceAttr(resourceFullName, "condition.right"),
	)

	typeCheck2 := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "condition.type", "AND"),
		resource.TestCheckNoResourceAttr(resourceFullName, "condition.comparator"),
		resource.TestCheckNoResourceAttr(resourceFullName, "condition.condition"),
		resource.TestCheckResourceAttr(resourceFullName, "condition.conditions.#", "2"),
		resource.TestMatchTypeSetElemNestedAttrs(resourceFullName, "condition.conditions.*", map[string]*regexp.Regexp{
			"type":        regexp.MustCompile("^COMPARISON$"),
			"comparator":  regexp.MustCompile("^EQUALS$"),
			"left.type":   regexp.MustCompile("^ATTRIBUTE$"),
			"left.id":     verify.P1ResourceIDRegexpFullString,
			"right.type":  regexp.MustCompile("^CONSTANT$"),
			"right.value": regexp.MustCompile("^test2$"),
		}),
		resource.TestMatchTypeSetElemNestedAttrs(resourceFullName, "condition.conditions.*", map[string]*regexp.Regexp{
			"type":                  regexp.MustCompile("^NOT$"),
			"condition.type":        regexp.MustCompile("^COMPARISON$"),
			"condition.comparator":  regexp.MustCompile("^EQUALS$"),
			"condition.left.type":   regexp.MustCompile("^ATTRIBUTE$"),
			"condition.left.id":     verify.P1ResourceIDRegexpFullString,
			"condition.right.type":  regexp.MustCompile("^CONSTANT$"),
			"condition.right.value": regexp.MustCompile("^test1$"),
		}),
		resource.TestCheckNoResourceAttr(resourceFullName, "condition.left"),
		resource.TestCheckNoResourceAttr(resourceFullName, "condition.reference"),
		resource.TestCheckNoResourceAttr(resourceFullName, "condition.right"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckFeatureFlag(t, acctest.ENUMFEATUREFLAG_AUTHORIZEPMTF)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             authorize.TrustFrameworkCondition_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// From scratch
			{
				Config: testAccTrustFrameworkConditionConfig_Condition_And1(resourceName, name),
				Check:  typeCheck1,
			},
			// Change
			{
				Config: testAccTrustFrameworkConditionConfig_Condition_And2(resourceName, name),
				Check:  typeCheck2,
			},
			{
				Config: testAccTrustFrameworkConditionConfig_Condition_And1(resourceName, name),
				Check:  typeCheck1,
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

func TestAccTrustFrameworkCondition_ConditionType_Comparison(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_trust_framework_condition.%s", resourceName)

	name := resourceName

	typeCheck1 := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "condition.type", "COMPARISON"),
		resource.TestCheckResourceAttr(resourceFullName, "condition.comparator", "CONTAINS"),
		resource.TestCheckNoResourceAttr(resourceFullName, "condition.condition"),
		resource.TestCheckNoResourceAttr(resourceFullName, "condition.conditions"),
		resource.TestCheckResourceAttr(resourceFullName, "condition.left.type", "ATTRIBUTE"),
		resource.TestMatchResourceAttr(resourceFullName, "condition.left.id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckNoResourceAttr(resourceFullName, "condition.reference"),
		resource.TestCheckResourceAttr(resourceFullName, "condition.right.type", "CONSTANT"),
		resource.TestCheckResourceAttr(resourceFullName, "condition.right.value", "test3"),
	)

	typeCheck2 := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "condition.type", "COMPARISON"),
		resource.TestCheckResourceAttr(resourceFullName, "condition.comparator", "EQUALS"),
		resource.TestCheckNoResourceAttr(resourceFullName, "condition.condition"),
		resource.TestCheckNoResourceAttr(resourceFullName, "condition.conditions"),
		resource.TestCheckResourceAttr(resourceFullName, "condition.left.type", "ATTRIBUTE"),
		resource.TestMatchResourceAttr(resourceFullName, "condition.left.id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckNoResourceAttr(resourceFullName, "condition.reference"),
		resource.TestCheckResourceAttr(resourceFullName, "condition.right.type", "ATTRIBUTE"),
		resource.TestMatchResourceAttr(resourceFullName, "condition.right.id", verify.P1ResourceIDRegexpFullString),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckFeatureFlag(t, acctest.ENUMFEATUREFLAG_AUTHORIZEPMTF)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             authorize.TrustFrameworkCondition_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// From scratch
			{
				Config: testAccTrustFrameworkConditionConfig_Condition_Comparison1(resourceName, name),
				Check:  typeCheck1,
			},
			// Change
			{
				Config: testAccTrustFrameworkConditionConfig_Condition_Comparison2(resourceName, name),
				Check:  typeCheck2,
			},
			{
				Config: testAccTrustFrameworkConditionConfig_Condition_Comparison1(resourceName, name),
				Check:  typeCheck1,
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

func TestAccTrustFrameworkCondition_ConditionType_Empty(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_trust_framework_condition.%s", resourceName)

	name := resourceName

	typeCheck1 := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "condition.type", "EMPTY"),
		resource.TestCheckNoResourceAttr(resourceFullName, "condition.comparator"),
		resource.TestCheckNoResourceAttr(resourceFullName, "condition.condition"),
		resource.TestCheckNoResourceAttr(resourceFullName, "condition.conditions"),
		resource.TestCheckNoResourceAttr(resourceFullName, "condition.left"),
		resource.TestCheckNoResourceAttr(resourceFullName, "condition.reference"),
		resource.TestCheckNoResourceAttr(resourceFullName, "condition.right"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckFeatureFlag(t, acctest.ENUMFEATUREFLAG_AUTHORIZEPMTF)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             authorize.TrustFrameworkCondition_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// From scratch
			{
				Config: testAccTrustFrameworkConditionConfig_Condition_Empty(resourceName, name),
				Check:  typeCheck1,
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

func TestAccTrustFrameworkCondition_ConditionType_Not(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_trust_framework_condition.%s", resourceName)

	name := resourceName

	typeCheck1 := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "condition.type", "NOT"),
		resource.TestCheckNoResourceAttr(resourceFullName, "condition.comparator"),
		resource.TestCheckResourceAttr(resourceFullName, "condition.condition.type", "EMPTY"),
		resource.TestCheckNoResourceAttr(resourceFullName, "condition.conditions"),
		resource.TestCheckNoResourceAttr(resourceFullName, "condition.left"),
		resource.TestCheckNoResourceAttr(resourceFullName, "condition.reference"),
		resource.TestCheckNoResourceAttr(resourceFullName, "condition.right"),
	)

	typeCheck2 := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "condition.type", "NOT"),
		resource.TestCheckNoResourceAttr(resourceFullName, "condition.comparator"),
		resource.TestCheckResourceAttr(resourceFullName, "condition.condition.type", "COMPARISON"),
		resource.TestCheckResourceAttr(resourceFullName, "condition.condition.comparator", "EQUALS"),
		resource.TestCheckResourceAttr(resourceFullName, "condition.condition.left.type", "ATTRIBUTE"),
		resource.TestMatchResourceAttr(resourceFullName, "condition.condition.left.id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "condition.condition.right.type", "CONSTANT"),
		resource.TestCheckResourceAttr(resourceFullName, "condition.condition.right.value", "test4"),
		resource.TestCheckNoResourceAttr(resourceFullName, "condition.conditions"),
		resource.TestCheckNoResourceAttr(resourceFullName, "condition.left"),
		resource.TestCheckNoResourceAttr(resourceFullName, "condition.reference"),
		resource.TestCheckNoResourceAttr(resourceFullName, "condition.right"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckFeatureFlag(t, acctest.ENUMFEATUREFLAG_AUTHORIZEPMTF)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             authorize.TrustFrameworkCondition_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// From scratch
			{
				Config: testAccTrustFrameworkConditionConfig_Condition_Not1(resourceName, name),
				Check:  typeCheck1,
			},
			{
				Config: testAccTrustFrameworkConditionConfig_Condition_Not2(resourceName, name),
				Check:  typeCheck2,
			},
			{
				Config: testAccTrustFrameworkConditionConfig_Condition_Not1(resourceName, name),
				Check:  typeCheck1,
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

func TestAccTrustFrameworkCondition_ConditionType_Or(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_trust_framework_condition.%s", resourceName)

	name := resourceName

	typeCheck1 := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "condition.type", "OR"),
		resource.TestCheckNoResourceAttr(resourceFullName, "condition.comparator"),
		resource.TestCheckNoResourceAttr(resourceFullName, "condition.condition"),
		resource.TestCheckResourceAttr(resourceFullName, "condition.conditions.#", "3"),
		resource.TestMatchTypeSetElemNestedAttrs(resourceFullName, "condition.conditions.*", map[string]*regexp.Regexp{
			"type":        regexp.MustCompile("^COMPARISON$"),
			"comparator":  regexp.MustCompile("^EQUALS$"),
			"left.type":   regexp.MustCompile("^ATTRIBUTE$"),
			"left.id":     verify.P1ResourceIDRegexpFullString,
			"right.type":  regexp.MustCompile("^CONSTANT$"),
			"right.value": regexp.MustCompile("^test2$"),
		}),
		resource.TestMatchTypeSetElemNestedAttrs(resourceFullName, "condition.conditions.*", map[string]*regexp.Regexp{
			"type":        regexp.MustCompile("^COMPARISON$"),
			"comparator":  regexp.MustCompile("^EQUALS$"),
			"left.type":   regexp.MustCompile("^ATTRIBUTE$"),
			"left.id":     verify.P1ResourceIDRegexpFullString,
			"right.type":  regexp.MustCompile("^CONSTANT$"),
			"right.value": regexp.MustCompile("^test1$"),
		}),
		resource.TestMatchTypeSetElemNestedAttrs(resourceFullName, "condition.conditions.*", map[string]*regexp.Regexp{
			"type":                  regexp.MustCompile("^NOT$"),
			"condition.type":        regexp.MustCompile("^COMPARISON$"),
			"condition.comparator":  regexp.MustCompile("^EQUALS$"),
			"condition.left.type":   regexp.MustCompile("^ATTRIBUTE$"),
			"condition.left.id":     verify.P1ResourceIDRegexpFullString,
			"condition.right.type":  regexp.MustCompile("^CONSTANT$"),
			"condition.right.value": regexp.MustCompile("^test3$"),
		}),
		resource.TestCheckNoResourceAttr(resourceFullName, "condition.left"),
		resource.TestCheckNoResourceAttr(resourceFullName, "condition.reference"),
		resource.TestCheckNoResourceAttr(resourceFullName, "condition.right"),
	)

	typeCheck2 := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "condition.type", "OR"),
		resource.TestCheckNoResourceAttr(resourceFullName, "condition.comparator"),
		resource.TestCheckNoResourceAttr(resourceFullName, "condition.condition"),
		resource.TestCheckResourceAttr(resourceFullName, "condition.conditions.#", "2"),
		resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "condition.conditions.*", map[string]string{
			"type": "EMPTY",
		}),
		resource.TestMatchTypeSetElemNestedAttrs(resourceFullName, "condition.conditions.*", map[string]*regexp.Regexp{
			"type":                  regexp.MustCompile("^NOT$"),
			"condition.type":        regexp.MustCompile("^COMPARISON$"),
			"condition.comparator":  regexp.MustCompile("^EQUALS$"),
			"condition.left.type":   regexp.MustCompile("^ATTRIBUTE$"),
			"condition.left.id":     verify.P1ResourceIDRegexpFullString,
			"condition.right.type":  regexp.MustCompile("^CONSTANT$"),
			"condition.right.value": regexp.MustCompile("^test1$"),
		}),
		resource.TestCheckNoResourceAttr(resourceFullName, "condition.left"),
		resource.TestCheckNoResourceAttr(resourceFullName, "condition.reference"),
		resource.TestCheckNoResourceAttr(resourceFullName, "condition.right"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckFeatureFlag(t, acctest.ENUMFEATUREFLAG_AUTHORIZEPMTF)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             authorize.TrustFrameworkCondition_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// From scratch
			{
				Config: testAccTrustFrameworkConditionConfig_Condition_Or1(resourceName, name),
				Check:  typeCheck1,
			},
			// Change
			{
				Config: testAccTrustFrameworkConditionConfig_Condition_Or2(resourceName, name),
				Check:  typeCheck2,
			},
			{
				Config: testAccTrustFrameworkConditionConfig_Condition_Or1(resourceName, name),
				Check:  typeCheck1,
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

func TestAccTrustFrameworkCondition_ConditionType_Reference(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_trust_framework_condition.%s", resourceName)

	name := resourceName

	typeCheck1 := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "condition.type", "REFERENCE"),
		resource.TestCheckNoResourceAttr(resourceFullName, "condition.comparator"),
		resource.TestCheckNoResourceAttr(resourceFullName, "condition.condition"),
		resource.TestCheckNoResourceAttr(resourceFullName, "condition.conditions"),
		resource.TestCheckNoResourceAttr(resourceFullName, "condition.left"),
		resource.TestMatchResourceAttr(resourceFullName, "condition.reference.id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckNoResourceAttr(resourceFullName, "condition.right"),
	)

	typeCheck2 := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "condition.type", "REFERENCE"),
		resource.TestCheckNoResourceAttr(resourceFullName, "condition.comparator"),
		resource.TestCheckNoResourceAttr(resourceFullName, "condition.condition"),
		resource.TestCheckNoResourceAttr(resourceFullName, "condition.conditions"),
		resource.TestCheckNoResourceAttr(resourceFullName, "condition.left"),
		resource.TestMatchResourceAttr(resourceFullName, "condition.reference.id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckNoResourceAttr(resourceFullName, "condition.right"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckFeatureFlag(t, acctest.ENUMFEATUREFLAG_AUTHORIZEPMTF)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             authorize.TrustFrameworkCondition_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// From scratch
			{
				Config: testAccTrustFrameworkConditionConfig_Condition_Reference1(resourceName, name),
				Check:  typeCheck1,
			},
			// Change
			{
				Config: testAccTrustFrameworkConditionConfig_Condition_Reference2(resourceName, name),
				Check:  typeCheck2,
			},
			{
				Config: testAccTrustFrameworkConditionConfig_Condition_Reference1(resourceName, name),
				Check:  typeCheck1,
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

func TestAccTrustFrameworkCondition_ConditionType_Change(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_trust_framework_condition.%s", resourceName)

	name := resourceName

	typeCheck1 := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "condition.type", "EMPTY"),
		resource.TestCheckNoResourceAttr(resourceFullName, "condition.comparator"),
		resource.TestCheckNoResourceAttr(resourceFullName, "condition.condition"),
		resource.TestCheckNoResourceAttr(resourceFullName, "condition.conditions"),
		resource.TestCheckNoResourceAttr(resourceFullName, "condition.left"),
		resource.TestCheckNoResourceAttr(resourceFullName, "condition.reference.id"),
		resource.TestCheckNoResourceAttr(resourceFullName, "condition.right"),
	)

	typeCheck2 := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "condition.type", "NOT"),
		resource.TestCheckNoResourceAttr(resourceFullName, "condition.comparator"),
		resource.TestCheckResourceAttr(resourceFullName, "condition.condition.type", "EMPTY"),
		resource.TestCheckNoResourceAttr(resourceFullName, "condition.conditions"),
		resource.TestCheckNoResourceAttr(resourceFullName, "condition.left"),
		resource.TestCheckNoResourceAttr(resourceFullName, "condition.reference"),
		resource.TestCheckNoResourceAttr(resourceFullName, "condition.right"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckFeatureFlag(t, acctest.ENUMFEATUREFLAG_AUTHORIZEPMTF)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             authorize.TrustFrameworkCondition_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// From scratch
			{
				Config: testAccTrustFrameworkConditionConfig_Condition_Empty(resourceName, name),
				Check:  typeCheck1,
			},
			// Change
			{
				Config: testAccTrustFrameworkConditionConfig_Condition_Not1(resourceName, name),
				Check:  typeCheck2,
			},
			{
				Config: testAccTrustFrameworkConditionConfig_Condition_Empty(resourceName, name),
				Check:  typeCheck1,
			},
		},
	})
}

func TestAccTrustFrameworkCondition_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_authorize_trust_framework_condition.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckFeatureFlag(t, acctest.ENUMFEATUREFLAG_AUTHORIZEPMTF)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             authorize.TrustFrameworkCondition_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccTrustFrameworkConditionConfig_Minimal(resourceName, name),
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

func testAccTrustFrameworkConditionConfig_NewEnv(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_authorize_trust_framework_condition" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[3]s"
  description    = "Test application role"

  condition = {
    type = "EMPTY"
  }
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccTrustFrameworkConditionConfig_Full(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_authorize_trust_framework_condition" "%[2]s-parent" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-parent"
  description    = "Test application role"

  condition = {
    type = "EMPTY"
  }
}

resource "pingone_authorize_trust_framework_condition" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Test application role full"

  parent = {
    id = pingone_authorize_trust_framework_condition.%[2]s-parent.id
  }

  condition = {
    type = "EMPTY"
  }
}`, acctest.AuthorizePMTFSandboxEnvironment(), resourceName, name)
}

func testAccTrustFrameworkConditionConfig_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_authorize_trust_framework_condition" "%[2]s-parent" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-parent"
  description    = "Test application role"

  condition = {
    type = "EMPTY"
  }
}`, testAccTrustFrameworkConditionConfig_Condition_Empty(resourceName, name), resourceName, name)
}

func testAccTrustFrameworkConditionConfig_Condition_And1(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_authorize_trust_framework_attribute" "%[2]s-current-user-id" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
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

resource "pingone_authorize_trust_framework_condition" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Test application role"

  condition = {
    type = "AND"

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
            value = "test3"
          }
        }
      }
    ]
  }
}`, acctest.AuthorizePMTFSandboxEnvironment(), resourceName, name)
}

func testAccTrustFrameworkConditionConfig_Condition_And2(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_authorize_trust_framework_attribute" "%[2]s-current-user-id" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
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

resource "pingone_authorize_trust_framework_condition" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Test application role"

  condition = {
    type = "AND"

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
}`, acctest.AuthorizePMTFSandboxEnvironment(), resourceName, name)
}

func testAccTrustFrameworkConditionConfig_Condition_Comparison1(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_authorize_trust_framework_attribute" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Test application role"

  resolvers = [
    {
      type = "CONSTANT"
      value_type = {
        type = "STRING"
      }
      value = "test"
    }
  ]

  value_type = {
    type = "JSON"
  }
}

resource "pingone_authorize_trust_framework_condition" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Test application role"

  condition = {
    type       = "COMPARISON"
    comparator = "CONTAINS"

    left = {
      type = "ATTRIBUTE"
      id   = pingone_authorize_trust_framework_attribute.%[2]s.id
    }

    right = {
      type  = "CONSTANT"
      value = "test3"
    }
  }
}`, acctest.AuthorizePMTFSandboxEnvironment(), resourceName, name)
}

func testAccTrustFrameworkConditionConfig_Condition_Comparison2(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_authorize_trust_framework_attribute" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Test application role"

  resolvers = [
    {
      type = "CONSTANT"
      value_type = {
        type = "STRING"
      }
      value = "test"
    }
  ]

  value_type = {
    type = "JSON"
  }
}

resource "pingone_authorize_trust_framework_condition" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Test application role"

  condition = {
    type       = "COMPARISON"
    comparator = "EQUALS"

    left = {
      type = "ATTRIBUTE"
      id   = pingone_authorize_trust_framework_attribute.%[2]s.id
    }

    right = {
      type = "ATTRIBUTE"
      id   = pingone_authorize_trust_framework_attribute.%[2]s.id
    }
  }
}`, acctest.AuthorizePMTFSandboxEnvironment(), resourceName, name)
}

func testAccTrustFrameworkConditionConfig_Condition_Empty(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_authorize_trust_framework_condition" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Test application role"

  condition = {
    type = "EMPTY"
  }
}`, acctest.AuthorizePMTFSandboxEnvironment(), resourceName, name)
}

func testAccTrustFrameworkConditionConfig_Condition_Not1(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_authorize_trust_framework_attribute" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Test application role"

  resolvers = [
    {
      type = "CONSTANT"
      value_type = {
        type = "STRING"
      }
      value = "test"
    }
  ]

  value_type = {
    type = "JSON"
  }
}

resource "pingone_authorize_trust_framework_condition" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Test application role"

  condition = {
    type = "NOT"

    condition = {
      type = "EMPTY"
    }
  }
}`, acctest.AuthorizePMTFSandboxEnvironment(), resourceName, name)
}

func testAccTrustFrameworkConditionConfig_Condition_Not2(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_authorize_trust_framework_attribute" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Test application role"

  resolvers = [
    {
      type = "CONSTANT"
      value_type = {
        type = "STRING"
      }
      value = "test"
    }
  ]

  value_type = {
    type = "JSON"
  }
}

resource "pingone_authorize_trust_framework_condition" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Test application role"

  condition = {
    type = "NOT"

    condition = {
      type       = "COMPARISON"
      comparator = "EQUALS"

      left = {
        type = "ATTRIBUTE"
        id   = pingone_authorize_trust_framework_attribute.%[2]s.id
      }

      right = {
        type  = "CONSTANT"
        value = "test4"
      }
    }
  }
}`, acctest.AuthorizePMTFSandboxEnvironment(), resourceName, name)
}

func testAccTrustFrameworkConditionConfig_Condition_Or1(resourceName, name string) string {
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

resource "pingone_authorize_trust_framework_condition" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Test application role"

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
            value = "test3"
          }
        }
      }
    ]
  }
}`, acctest.AuthorizePMTFSandboxEnvironment(), resourceName, name)
}

func testAccTrustFrameworkConditionConfig_Condition_Or2(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_authorize_trust_framework_attribute" "%[2]s-current-user-id" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
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

resource "pingone_authorize_trust_framework_condition" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Test application role"

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
}`, acctest.AuthorizePMTFSandboxEnvironment(), resourceName, name)
}

func testAccTrustFrameworkConditionConfig_Condition_Reference1(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_authorize_trust_framework_condition" "%[2]s-ref1" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-ref1"
  description    = "Test application role"

  condition = {
    type = "EMPTY"
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

resource "pingone_authorize_trust_framework_condition" "%[2]s-ref2" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-ref2"
  description    = "Test application role"

  condition = {
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
        value = "test2"
      }
    }
  }
}

resource "pingone_authorize_trust_framework_condition" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Test application role"

  condition = {
    type = "REFERENCE"

    reference = {
      id = pingone_authorize_trust_framework_condition.%[2]s-ref1.id
    }
  }
}`, acctest.AuthorizePMTFSandboxEnvironment(), resourceName, name)
}

func testAccTrustFrameworkConditionConfig_Condition_Reference2(resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_authorize_trust_framework_condition" "%[2]s-ref1" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-ref1"
  description    = "Test application role"

  condition = {
    type = "EMPTY"
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

resource "pingone_authorize_trust_framework_condition" "%[2]s-ref2" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-ref2"
  description    = "Test application role"

  condition = {
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
        value = "test2"
      }
    }
  }
}

resource "pingone_authorize_trust_framework_condition" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "Test application role"

  condition = {
    type = "REFERENCE"

    reference = {
      id = pingone_authorize_trust_framework_condition.%[2]s-ref2.id
    }
  }
}`, acctest.AuthorizePMTFSandboxEnvironment(), resourceName, name)
}
