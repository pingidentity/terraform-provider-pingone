// Copyright © 2026 Ping Identity Corporation

package risk_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	acctestlegacysdk "github.com/pingidentity/terraform-provider-pingone/internal/acctest/legacysdk"
	baselegacysdk "github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/base/legacysdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/risk"
	"github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccRiskPolicy_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_risk_policy.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	var riskPolicyID, environmentID string

	var p1Client *client.Client
	var ctx = context.Background()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoBeta(t)

			p1Client = acctestlegacysdk.PreCheckTestClient(ctx, t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             risk.RiskPolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccRiskPolicyConfig_Minimal(resourceName, name),
				Check:  risk.RiskPolicy_GetIDs(resourceFullName, &environmentID, &riskPolicyID),
			},
			{
				PreConfig: func() {
					risk.RiskPolicy_RemovalDrift_PreConfig(ctx, p1Client.API.RiskAPIClient, t, environmentID, riskPolicyID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the environment
			{
				Config: testAccRiskPolicyConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check:  risk.RiskPolicy_GetIDs(resourceFullName, &environmentID, &riskPolicyID),
			},
			{
				PreConfig: func() {
					baselegacysdk.Environment_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccRiskPolicy_NewEnv(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_risk_policy.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             risk.RiskPolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRiskPolicyConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
				),
			},
		},
	})
}

func TestAccRiskPolicy_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_risk_policy.%s", resourceName)

	name := resourceName

	fullCheck := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "name", name),
		resource.TestCheckResourceAttr(resourceFullName, "default_result.type", "VALUE"),
		resource.TestCheckResourceAttr(resourceFullName, "default_result.level", "LOW"),
		resource.TestCheckResourceAttr(resourceFullName, "default", "false"),
		resource.TestCheckResourceAttr(resourceFullName, "overrides.#", "0"),
	)

	minimalCheck := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "name", name),
		resource.TestCheckResourceAttr(resourceFullName, "default_result.type", "VALUE"),
		resource.TestCheckResourceAttr(resourceFullName, "default_result.level", "LOW"),
		resource.TestCheckResourceAttr(resourceFullName, "default", "false"),
		resource.TestCheckResourceAttr(resourceFullName, "overrides.#", "0"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckTestAccFlaky(t) // PND-5900: policy PUT reverts to previous config on next GET, leaving a non-empty refresh plan
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             risk.RiskPolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			{
				Config: testAccRiskPolicyConfig_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config:  testAccRiskPolicyConfig_Full(resourceName, name),
				Destroy: true,
			},
			// Minimal
			{
				Config: testAccRiskPolicyConfig_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config:  testAccRiskPolicyConfig_Minimal(resourceName, name),
				Destroy: true,
			},
			// Change
			{
				Config: testAccRiskPolicyConfig_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config: testAccRiskPolicyConfig_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config: testAccRiskPolicyConfig_Full(resourceName, name),
				Check:  fullCheck,
			},
		},
	})
}

func TestAccRiskPolicy_Scores(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_risk_policy.%s", resourceName)

	name := resourceName

	fullCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "policy_scores.policy_threshold_medium.min_score", "45"),
		resource.TestCheckResourceAttr(resourceFullName, "policy_scores.policy_threshold_medium.max_score", "80"),
		resource.TestCheckResourceAttr(resourceFullName, "policy_scores.policy_threshold_high.min_score", "80"),
		resource.TestCheckResourceAttr(resourceFullName, "policy_scores.policy_threshold_high.max_score", "1000"),
		resource.TestCheckResourceAttr(resourceFullName, "policy_scores.predictors.#", "2"),
		resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "policy_scores.predictors.*", map[string]string{
			"compact_name":              fmt.Sprintf("%s1", name),
			"predictor_reference_value": fmt.Sprintf("${details.%s1.level}", name),
			"score":                     "55",
		}),
		resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "policy_scores.predictors.*", map[string]string{
			"compact_name":              fmt.Sprintf("%s3", name),
			"predictor_reference_value": fmt.Sprintf("${details.%s3.level}", name),
			"score":                     "45",
		}),
		resource.TestCheckResourceAttr(resourceFullName, "evaluated_predictors.#", "3"),
		resource.TestMatchResourceAttr(resourceFullName, "evaluated_predictors.0", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "evaluated_predictors.1", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "evaluated_predictors.2", verify.P1ResourceIDRegexpFullString),
	)

	minimalCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "policy_scores.policy_threshold_medium.min_score", "35"),
		resource.TestCheckResourceAttr(resourceFullName, "policy_scores.policy_threshold_medium.max_score", "70"),
		resource.TestCheckResourceAttr(resourceFullName, "policy_scores.policy_threshold_high.min_score", "70"),
		resource.TestCheckResourceAttr(resourceFullName, "policy_scores.policy_threshold_high.max_score", "1000"),
		resource.TestCheckResourceAttr(resourceFullName, "policy_scores.predictors.#", "1"),
		resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "policy_scores.predictors.*", map[string]string{
			"compact_name":              fmt.Sprintf("%s2", name),
			"predictor_reference_value": fmt.Sprintf("${details.%s2.level}", name),
			"score":                     "45",
		}),
		resource.TestCheckResourceAttr(resourceFullName, "evaluated_predictors.#", "1"),
		resource.TestMatchResourceAttr(resourceFullName, "evaluated_predictors.0", verify.P1ResourceIDRegexpFullString),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckTestAccFlaky(t) // PND-5900: policy PUT reverts to previous config on next GET, leaving a non-empty refresh plan
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             risk.RiskPolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			{
				Config: testAccRiskPolicyConfig_Scores_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config:  testAccRiskPolicyConfig_Scores_Full(resourceName, name),
				Destroy: true,
			},
			// Minimal
			{
				Config: testAccRiskPolicyConfig_Scores_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config:  testAccRiskPolicyConfig_Scores_Minimal(resourceName, name),
				Destroy: true,
			},
			// Change
			{
				Config: testAccRiskPolicyConfig_Scores_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config: testAccRiskPolicyConfig_Scores_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config: testAccRiskPolicyConfig_Scores_Full(resourceName, name),
				Check:  fullCheck,
			},
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("resource not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:  testAccRiskPolicyConfig_Scores_Full(resourceName, name),
				Destroy: true,
			},
			// Errors
			{
				Config:      testAccRiskPolicyConfig_Scores_MediumScoreAboveMaxScore(resourceName, name),
				ExpectError: regexp.MustCompile(`Provided value is not valid`),
			},
			{
				Config:      testAccRiskPolicyConfig_Scores_DefinedPolicyPredictorNotInEvaluated(resourceName, name),
				ExpectError: regexp.MustCompile(`A predictor in the policy set is not listed in "evaluated_predictors".`),
			},
		},
	})
}

func TestAccRiskPolicy_Weights(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_risk_policy.%s", resourceName)

	name := resourceName

	fullCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "policy_weights.policy_threshold_medium.min_score", "30"),
		resource.TestCheckResourceAttr(resourceFullName, "policy_weights.policy_threshold_medium.max_score", "80"),
		resource.TestCheckResourceAttr(resourceFullName, "policy_weights.policy_threshold_high.min_score", "80"),
		resource.TestCheckResourceAttr(resourceFullName, "policy_weights.policy_threshold_high.max_score", "100"),
		resource.TestCheckResourceAttr(resourceFullName, "policy_weights.predictors.#", "2"),
		resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "policy_weights.predictors.*", map[string]string{
			"compact_name":              fmt.Sprintf("%s1", name),
			"predictor_reference_value": fmt.Sprintf("${details.aggregatedWeights.%s1}", name),
			"weight":                    "4",
		}),
		resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "policy_weights.predictors.*", map[string]string{
			"compact_name":              fmt.Sprintf("%s3", name),
			"predictor_reference_value": fmt.Sprintf("${details.aggregatedWeights.%s3}", name),
			"weight":                    "6",
		}),
		resource.TestCheckResourceAttr(resourceFullName, "evaluated_predictors.#", "3"),
		resource.TestMatchResourceAttr(resourceFullName, "evaluated_predictors.0", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "evaluated_predictors.1", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "evaluated_predictors.2", verify.P1ResourceIDRegexpFullString),
	)

	minimalCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "policy_weights.policy_threshold_medium.min_score", "40"),
		resource.TestCheckResourceAttr(resourceFullName, "policy_weights.policy_threshold_medium.max_score", "70"),
		resource.TestCheckResourceAttr(resourceFullName, "policy_weights.policy_threshold_high.min_score", "70"),
		resource.TestCheckResourceAttr(resourceFullName, "policy_weights.policy_threshold_high.max_score", "100"),
		resource.TestCheckResourceAttr(resourceFullName, "policy_weights.predictors.#", "1"),
		resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "policy_weights.predictors.*", map[string]string{
			"compact_name":              fmt.Sprintf("%s2", name),
			"predictor_reference_value": fmt.Sprintf("${details.aggregatedWeights.%s2}", name),
			"weight":                    "3",
		}),
		resource.TestCheckResourceAttr(resourceFullName, "evaluated_predictors.#", "1"),
		resource.TestMatchResourceAttr(resourceFullName, "evaluated_predictors.0", verify.P1ResourceIDRegexpFullString),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckTestAccFlaky(t) // PND-5900: policy PUT reverts to previous config on next GET, leaving a non-empty refresh plan
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             risk.RiskPolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			{
				Config: testAccRiskPolicyConfig_Weights_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config:  testAccRiskPolicyConfig_Weights_Full(resourceName, name),
				Destroy: true,
			},
			// Minimal
			{
				Config: testAccRiskPolicyConfig_Weights_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config:  testAccRiskPolicyConfig_Weights_Minimal(resourceName, name),
				Destroy: true,
			},
			// Change
			{
				Config: testAccRiskPolicyConfig_Weights_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config: testAccRiskPolicyConfig_Weights_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config: testAccRiskPolicyConfig_Weights_Full(resourceName, name),
				Check:  fullCheck,
			},
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("resource not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:  testAccRiskPolicyConfig_Weights_Full(resourceName, name),
				Destroy: true,
			},
			// Errors
			{
				Config:      testAccRiskPolicyConfig_Weights_MediumScoreAboveMaxScore(resourceName, name),
				ExpectError: regexp.MustCompile(`Provided value is not valid`),
			},
			{
				Config:      testAccRiskPolicyConfig_Weights_DefinedPolicyPredictorNotInEvaluated(resourceName, name),
				ExpectError: regexp.MustCompile(`A predictor in the policy set is not listed in "evaluated_predictors".`),
			},
		},
	})
}

func TestAccRiskPolicy_ChangeType(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_risk_policy.%s", resourceName)

	name := resourceName

	scoresCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "policy_scores.policy_threshold_medium.min_score", "45"),
		resource.TestCheckResourceAttr(resourceFullName, "policy_scores.policy_threshold_medium.max_score", "80"),
		resource.TestCheckResourceAttr(resourceFullName, "policy_scores.policy_threshold_high.min_score", "80"),
		resource.TestCheckResourceAttr(resourceFullName, "policy_scores.policy_threshold_high.max_score", "1000"),
		resource.TestCheckResourceAttr(resourceFullName, "policy_scores.predictors.#", "2"),
		resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "policy_scores.predictors.*", map[string]string{
			"compact_name":              fmt.Sprintf("%s1", name),
			"predictor_reference_value": fmt.Sprintf("${details.%s1.level}", name),
			"score":                     "55",
		}),
		resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "policy_scores.predictors.*", map[string]string{
			"compact_name":              fmt.Sprintf("%s3", name),
			"predictor_reference_value": fmt.Sprintf("${details.%s3.level}", name),
			"score":                     "45",
		}),
		resource.TestCheckResourceAttr(resourceFullName, "evaluated_predictors.#", "3"),
		resource.TestMatchResourceAttr(resourceFullName, "evaluated_predictors.0", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "evaluated_predictors.1", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "evaluated_predictors.2", verify.P1ResourceIDRegexpFullString),
	)

	weightsCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "policy_weights.policy_threshold_medium.min_score", "40"),
		resource.TestCheckResourceAttr(resourceFullName, "policy_weights.policy_threshold_medium.max_score", "70"),
		resource.TestCheckResourceAttr(resourceFullName, "policy_weights.policy_threshold_high.min_score", "70"),
		resource.TestCheckResourceAttr(resourceFullName, "policy_weights.policy_threshold_high.max_score", "100"),
		resource.TestCheckResourceAttr(resourceFullName, "policy_weights.predictors.#", "1"),
		resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "policy_weights.predictors.*", map[string]string{
			"compact_name":              fmt.Sprintf("%s2", name),
			"predictor_reference_value": fmt.Sprintf("${details.aggregatedWeights.%s2}", name),
			"weight":                    "3",
		}),
		resource.TestCheckResourceAttr(resourceFullName, "evaluated_predictors.#", "1"),
		resource.TestMatchResourceAttr(resourceFullName, "evaluated_predictors.0", verify.P1ResourceIDRegexpFullString),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckTestAccFlaky(t) // PND-5900: policy PUT reverts to previous config on next GET, leaving a non-empty refresh plan
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             risk.RiskPolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRiskPolicyConfig_Scores_Full(resourceName, name),
				Check:  scoresCheck,
			},
			{
				Config: testAccRiskPolicyConfig_Weights_Minimal(resourceName, name),
				Check:  weightsCheck,
			},
			{
				Config: testAccRiskPolicyConfig_Scores_Full(resourceName, name),
				Check:  scoresCheck,
			},
		},
	})
}

func TestAccRiskPolicy_PolicyOverrides(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_risk_policy.%s", resourceName)

	name := resourceName

	fullCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "overrides.#", "3"),
		// Ordering: overrides is an ordered list (not a set). The order in HCL
		// determines override priority during evaluation, so the indexed
		// elements must round-trip in the configured order.
		resource.TestCheckResourceAttr(resourceFullName, "overrides.0.name", "my_anon_check"),
		resource.TestCheckResourceAttr(resourceFullName, "overrides.0.priority", "1"),
		resource.TestCheckResourceAttr(resourceFullName, "overrides.1.name", "my_ip_vel_check"),
		resource.TestCheckResourceAttr(resourceFullName, "overrides.1.priority", "2"),
		resource.TestCheckResourceAttr(resourceFullName, "overrides.2.name", "allowed_list"),
		resource.TestCheckResourceAttr(resourceFullName, "overrides.2.priority", "3"),
		resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "overrides.*", map[string]string{
			"name":                                "my_anon_check",
			"priority":                            "1",
			"result.level":                        "HIGH",
			"result.value":                        "starling",
			"result.type":                         "VALUE",
			"condition.type":                      "VALUE_COMPARISON",
			"condition.equals":                    "HIGH",
			"condition.compact_name":              "anonymousNetwork",
			"condition.predictor_reference_value": "${details.anonymousNetwork.level}",
		}),
		resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "overrides.*", map[string]string{
			"name":                                "my_ip_vel_check",
			"priority":                            "2",
			"result.level":                        "MEDIUM",
			"result.value":                        "crow",
			"result.type":                         "VALUE",
			"condition.type":                      "VALUE_COMPARISON",
			"condition.equals":                    "MEDIUM",
			"condition.compact_name":              fmt.Sprintf("%s1", name),
			"condition.predictor_reference_value": fmt.Sprintf("${details.%s1.level}", name),
		}),
		resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "overrides.*", map[string]string{
			"name":                                   "allowed_list",
			"priority":                               "3",
			"result.level":                           "LOW",
			"result.value":                           "sparrow",
			"result.type":                            "VALUE",
			"condition.type":                         "IP_RANGE",
			"condition.ip_range.#":                   "3",
			"condition.predictor_reference_contains": "${transaction.ip}",
		}),
	)

	minimalCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "overrides.#", "1"),
		resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "overrides.*", map[string]string{
			"name":                                "geoVelocity",
			"priority":                            "1",
			"result.level":                        "MEDIUM",
			"result.type":                         "VALUE",
			"condition.type":                      "VALUE_COMPARISON",
			"condition.equals":                    "HIGH",
			"condition.compact_name":              "geoVelocity",
			"condition.predictor_reference_value": "${details.geoVelocity.level}",
		}),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckTestAccFlaky(t) // PND-5900: policy PUT reverts to previous config on next GET, leaving a non-empty refresh plan
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             risk.RiskPolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			{
				Config: testAccRiskPolicyConfig_Overrides_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config:  testAccRiskPolicyConfig_Overrides_Full(resourceName, name),
				Destroy: true,
			},
			// Minimal
			{
				Config: testAccRiskPolicyConfig_Overrides_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config:  testAccRiskPolicyConfig_Overrides_Minimal(resourceName, name),
				Destroy: true,
			},
			// Change
			{
				Config: testAccRiskPolicyConfig_Overrides_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config: testAccRiskPolicyConfig_Overrides_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config: testAccRiskPolicyConfig_Overrides_Full(resourceName, name),
				Check:  fullCheck,
			},
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("resource not found: %s", resourceFullName)
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

func TestAccRiskPolicy_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_risk_policy.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckTestAccFlaky(t) // PND-5900: policy PUT reverts to previous config on next GET, leaving a non-empty refresh plan
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             risk.RiskPolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccRiskPolicyConfig_Minimal(resourceName, name),
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

func testAccRiskPolicyConfig_NewEnv(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_risk_policy" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  name = "%[4]s"

  policy_scores = {
    policy_threshold_medium = {
      min_score = 40
    }

    policy_threshold_high = {
      min_score = 75
    }

    predictors = [
      {
        compact_name = "ipRisk"
        score        = 50
      },
      {
        compact_name = "geoVelocity"
        score        = 50
      }
    ]
  }
}`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccRiskPolicyConfig_Full(resourceName, name string) string {
	return testAccRiskPolicyConfig_Scores_Full(resourceName, name)
}

func testAccRiskPolicyConfig_Minimal(resourceName, name string) string {
	return testAccRiskPolicyConfig_Scores_Minimal(resourceName, name)
}

func testAccRiskPolicyConfig_Scores_Full(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s-1" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s-1"
  compact_name = "%[3]s1"

  predictor_geovelocity = {}
}

resource "pingone_risk_predictor" "%[2]s-2" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s-2"
  compact_name = "%[3]s2"

  predictor_ip_reputation = {}
}

resource "pingone_risk_predictor" "%[2]s-3" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s-3"
  compact_name = "%[3]s3"

  predictor_device = {}
}

resource "pingone_risk_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  default_result = {
    level = "LOW"
  }

  evaluated_predictors = [
    pingone_risk_predictor.%[2]s-1.id,
    pingone_risk_predictor.%[2]s-2.id,
    pingone_risk_predictor.%[2]s-3.id,
  ]

  policy_scores = {
    policy_threshold_medium = {
      min_score = 45
    }

    policy_threshold_high = {
      min_score = 80
    }

    predictors = [
      {
        compact_name = pingone_risk_predictor.%[2]s-1.compact_name
        score        = 55
      },
      {
        compact_name = pingone_risk_predictor.%[2]s-3.compact_name
        score        = 45
      }
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPolicyConfig_Scores_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s-1" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s-1"
  compact_name = "%[3]s1"

  predictor_geovelocity = {}
}

resource "pingone_risk_predictor" "%[2]s-2" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s-2"
  compact_name = "%[3]s2"

  predictor_ip_reputation = {}
}

resource "pingone_risk_predictor" "%[2]s-3" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s-3"
  compact_name = "%[3]s3"

  predictor_device = {}
}

resource "pingone_risk_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  policy_scores = {
    policy_threshold_medium = {
      min_score = 35
    }

    policy_threshold_high = {
      min_score = 70
    }

    predictors = [
      {
        compact_name = pingone_risk_predictor.%[2]s-2.compact_name
        score        = 45
      }
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPolicyConfig_Scores_DefinedPolicyPredictorNotInEvaluated(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s-1" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s-1"
  compact_name = "%[3]s1"

  predictor_geovelocity = {}
}

resource "pingone_risk_predictor" "%[2]s-2" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s-2"
  compact_name = "%[3]s2"

  predictor_ip_reputation = {}
}

resource "pingone_risk_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  evaluated_predictors = [
    pingone_risk_predictor.%[2]s-1.id,
  ]

  policy_scores = {
    policy_threshold_medium = {
      min_score = 35
    }

    policy_threshold_high = {
      min_score = 70
    }

    predictors = [
      {
        compact_name = pingone_risk_predictor.%[2]s-2.compact_name
        score        = 45
      }
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPolicyConfig_Scores_MediumScoreAboveMaxScore(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  policy_scores = {
    policy_threshold_medium = {
      min_score = 80
    }

    policy_threshold_high = {
      min_score = 50
    }

    predictors = [
      {
        compact_name = "ipRisk"
        score        = 45
      }
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPolicyConfig_Weights_Full(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s-1" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s-1"
  compact_name = "%[3]s1"

  predictor_geovelocity = {}

}

resource "pingone_risk_predictor" "%[2]s-2" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s-2"
  compact_name = "%[3]s2"

  predictor_ip_reputation = {}

}

resource "pingone_risk_predictor" "%[2]s-3" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s-3"
  compact_name = "%[3]s3"

  predictor_device = {}
}

resource "pingone_risk_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  default_result = {
    level = "LOW"
  }

  evaluated_predictors = [
    pingone_risk_predictor.%[2]s-1.id,
    pingone_risk_predictor.%[2]s-2.id,
    pingone_risk_predictor.%[2]s-3.id,
  ]

  policy_weights = {
    policy_threshold_medium = {
      min_score = 30
    }

    policy_threshold_high = {
      min_score = 80
    }

    predictors = [
      {
        compact_name = pingone_risk_predictor.%[2]s-1.compact_name
        weight       = 4
      },
      {
        compact_name = pingone_risk_predictor.%[2]s-3.compact_name
        weight       = 6
      }
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPolicyConfig_Weights_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s-1" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s-1"
  compact_name = "%[3]s1"

  predictor_geovelocity = {}

}

resource "pingone_risk_predictor" "%[2]s-2" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s-2"
  compact_name = "%[3]s2"

  predictor_ip_reputation = {}

}

resource "pingone_risk_predictor" "%[2]s-3" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s-3"
  compact_name = "%[3]s3"

  predictor_device = {}
}

resource "pingone_risk_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  policy_weights = {
    policy_threshold_medium = {
      min_score = 40
    }

    policy_threshold_high = {
      min_score = 70
    }

    predictors = [
      {
        compact_name = pingone_risk_predictor.%[2]s-2.compact_name
        weight       = 3
      }
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPolicyConfig_Weights_DefinedPolicyPredictorNotInEvaluated(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s-1" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s-1"
  compact_name = "%[3]s1"

  predictor_geovelocity = {}

}

resource "pingone_risk_predictor" "%[2]s-2" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s-2"
  compact_name = "%[3]s2"

  predictor_ip_reputation = {}

}

resource "pingone_risk_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  evaluated_predictors = [
    pingone_risk_predictor.%[2]s-1.id,
  ]

  policy_weights = {
    policy_threshold_medium = {
      min_score = 40
    }

    policy_threshold_high = {
      min_score = 70
    }

    predictors = [
      {
        compact_name = pingone_risk_predictor.%[2]s-2.compact_name
        weight       = 3
      }
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPolicyConfig_Weights_MediumScoreAboveMaxScore(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  policy_weights = {
    policy_threshold_medium = {
      min_score = 80
    }

    policy_threshold_high = {
      min_score = 50
    }

    predictors = [
      {
        compact_name = "ipRisk"
        weight       = 4
      }
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPolicyConfig_Overrides_Full(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[3]s1"

  predictor_velocity = {
    of = "$${event.user.id}"
  }
}

resource "pingone_risk_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  policy_scores = {
    policy_threshold_medium = {
      min_score = 35
    }

    policy_threshold_high = {
      min_score = 70
    }

    predictors = [
      {
        compact_name = "ipRisk"
        score        = 45
      },
      {
        compact_name = "geoVelocity"
        score        = 45
      }
    ]
  }

  overrides = [
    {
      name = "my_anon_check"

      result = {
        level = "HIGH"
        value = "starling"
      }

      condition = {
        type         = "VALUE_COMPARISON"
        compact_name = "anonymousNetwork"
        equals       = "HIGH"
      }
    },

    {
      name = "my_ip_vel_check"

      result = {
        level = "MEDIUM"
        value = "crow"
      }

      condition = {
        type         = "VALUE_COMPARISON"
        compact_name = pingone_risk_predictor.%[2]s.compact_name
        equals       = "MEDIUM"
      }
    },

    {
      name = "allowed_list"

      result = {
        level = "LOW"
        value = "sparrow"
      }

      condition = {
        type = "IP_RANGE"
        ip_range = [
          "10.0.0.0/8",
          "172.16.0.0/12",
          "192.168.0.0/24"
        ]
      }
    }
  ]
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPolicyConfig_Overrides_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[3]s1"

  predictor_velocity = {
    of = "$${event.user.id}"
  }
}

resource "pingone_risk_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  policy_scores = {
    policy_threshold_medium = {
      min_score = 35
    }

    policy_threshold_high = {
      min_score = 70
    }

    predictors = [
      {
        compact_name = "ipRisk"
        score        = 45
      },
      {
        compact_name = "geoVelocity"
        score        = 45
      }
    ]
  }

  overrides = [
    {
      result = {
        level = "MEDIUM"
      }

      condition = {
        type         = "VALUE_COMPARISON"
        compact_name = "geoVelocity"
        equals       = "HIGH"
      }
    }
  ]
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func TestAccRiskPolicy_Mitigations(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_risk_policy.%s", resourceName)

	name := resourceName

	fullCheck := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "name", name),
		resource.TestCheckResourceAttr(resourceFullName, "mitigations.#", "4"),
		// Ordering: mitigations is an ordered list (not a set). The order in HCL
		// determines evaluation priority, so the indexed elements must round-trip
		// in the configured order.
		resource.TestCheckResourceAttr(resourceFullName, "mitigations.0.name", "anonymousNetwork"),
		resource.TestCheckResourceAttr(resourceFullName, "mitigations.0.priority", "1"),
		resource.TestCheckResourceAttr(resourceFullName, "mitigations.1.name", "geoVelocity"),
		resource.TestCheckResourceAttr(resourceFullName, "mitigations.1.priority", "2"),
		resource.TestCheckResourceAttr(resourceFullName, "mitigations.1.action", "MFA"),
		resource.TestCheckResourceAttr(resourceFullName, "mitigations.2.name", "geoVelocity"),
		resource.TestCheckResourceAttr(resourceFullName, "mitigations.2.priority", "3"),
		resource.TestCheckResourceAttr(resourceFullName, "mitigations.2.action", "VERIFY"),
		resource.TestCheckResourceAttr(resourceFullName, "mitigations.3.name", "ipRisk"),
		resource.TestCheckResourceAttr(resourceFullName, "mitigations.3.priority", "4"),
		resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "mitigations.*", map[string]string{
			"name":                                "anonymousNetwork",
			"priority":                            "1",
			"action":                              "CUSTOM",
			"custom_action":                       "customActionValue",
			"condition.type":                      "VALUE_COMPARISON",
			"condition.equals":                    "HIGH",
			"condition.compact_name":              "anonymousNetwork",
			"condition.predictor_reference_value": "${details.anonymousNetwork.level}",
		}),
		resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "mitigations.*", map[string]string{
			"name":                                "geoVelocity",
			"priority":                            "2",
			"action":                              "MFA",
			"condition.type":                      "VALUE_COMPARISON",
			"condition.equals":                    "HIGH",
			"condition.compact_name":              "geoVelocity",
			"condition.predictor_reference_value": "${details.geoVelocity.level}",
		}),
		resource.TestCheckResourceAttrPair(resourceFullName, "mitigations.1.mfa_authentication_policy_id", fmt.Sprintf("pingone_mfa_device_policy.%s", resourceName), "id"),
		resource.TestCheckResourceAttrPair(resourceFullName, "mitigations.1.mfa_registration_policy_id", fmt.Sprintf("pingone_mfa_device_policy.%s_reg", resourceName), "id"),
		resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "mitigations.*", map[string]string{
			"name":                                "geoVelocity",
			"priority":                            "3",
			"action":                              "VERIFY",
			"condition.type":                      "VALUE_COMPARISON",
			"condition.equals":                    "MEDIUM",
			"condition.compact_name":              "geoVelocity",
			"condition.predictor_reference_value": "${details.geoVelocity.level}",
		}),
		resource.TestCheckResourceAttrPair(resourceFullName, "mitigations.2.verify_policy_id", fmt.Sprintf("pingone_verify_policy.%s", resourceName), "id"),
		resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "mitigations.*", map[string]string{
			"name":                                "ipRisk",
			"priority":                            "4",
			"action":                              "DENY_AND_SUSPEND",
			"condition.type":                      "VALUE_COMPARISON",
			"condition.equals":                    "HIGH",
			"condition.compact_name":              "ipRisk",
			"condition.predictor_reference_value": "${details.ipRisk.level}",
		}),
		resource.TestCheckResourceAttr(resourceFullName, "fallback.action", "APPROVE"),
		resource.TestCheckResourceAttr(resourceFullName, "targets.condition.and.#", "2"),
		resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "targets.condition.and.*", map[string]string{
			"contains": "${event.flow.type}",
			"type":     "STRING_LIST",
		}),
		resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "targets.condition.and.*", map[string]string{
			"contains": "${event.user.groups}",
			"type":     "GROUPS_INTERSECTION",
		}),
	)

	minimalCheck := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "name", name),
		resource.TestCheckResourceAttr(resourceFullName, "mitigations.#", "1"),
		resource.TestCheckTypeSetElemNestedAttrs(resourceFullName, "mitigations.*", map[string]string{
			"name":                                "anonymousNetwork",
			"priority":                            "1",
			"action":                              "DENY",
			"condition.type":                      "VALUE_COMPARISON",
			"condition.equals":                    "HIGH",
			"condition.compact_name":              "anonymousNetwork",
			"condition.predictor_reference_value": "${details.anonymousNetwork.level}",
		}),
		resource.TestCheckResourceAttr(resourceFullName, "fallback.action", "DENY"),
	)

	fallbackMFACheck := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "name", name),
		resource.TestCheckResourceAttr(resourceFullName, "fallback.action", "MFA"),
		resource.TestCheckResourceAttrPair(resourceFullName, "fallback.mfa_authentication_policy_id", fmt.Sprintf("pingone_mfa_device_policy.%s", resourceName), "id"),
	)

	fallbackOnlyCheck := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "name", name),
		resource.TestCheckResourceAttr(resourceFullName, "fallback.action", "DENY"),
		resource.TestCheckNoResourceAttr(resourceFullName, "mitigations.#"),
		resource.TestCheckNoResourceAttr(resourceFullName, "targets.%"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckTestAccFlaky(t) // PND-5900: policy PUT reverts to previous config on next GET, leaving a non-empty refresh plan
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             risk.RiskPolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			{
				Config: testAccRiskPolicyConfig_Mitigations_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config:  testAccRiskPolicyConfig_Mitigations_Full(resourceName, name),
				Destroy: true,
			},
			// Minimal
			{
				Config: testAccRiskPolicyConfig_Mitigations_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config:  testAccRiskPolicyConfig_Mitigations_Minimal(resourceName, name),
				Destroy: true,
			},
			// Change
			{
				Config: testAccRiskPolicyConfig_Mitigations_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config: testAccRiskPolicyConfig_Mitigations_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config: testAccRiskPolicyConfig_Mitigations_Full(resourceName, name),
				Check:  fullCheck,
			},
			// Fallback MFA
			{
				Config: testAccRiskPolicyConfig_Mitigations_FallbackMFA(resourceName, name),
				Check:  fallbackMFACheck,
			},
			// Remove mitigations -> fallback only
			{
				Config: testAccRiskPolicyConfig_Mitigations_FallbackOnly(resourceName, name),
				Check:  fallbackOnlyCheck,
			},
			{
				Config: testAccRiskPolicyConfig_Mitigations_Full(resourceName, name),
				Check:  fullCheck,
			},
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("resource not found: %s", resourceFullName)
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

// TestAccRiskPolicy_MitigationsValidation covers schema-level config validation
// for the mitigations / fallback / targets blocks. These steps never reach a
// successful apply, so they live in their own test to avoid a config-validation
// ExpectError step being the final state the framework tries to destroy.
func TestAccRiskPolicy_MitigationsValidation(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             risk.RiskPolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Mutual exclusion: mitigations + overrides
			{
				Config:      testAccRiskPolicyConfig_Mitigations_ConflictsWithOverrides(resourceName, name),
				ExpectError: regexp.MustCompile(`Attribute "mitigations" cannot be specified when "overrides" is specified`),
			},
			// Mutual exclusion: targets + overrides
			{
				Config:      testAccRiskPolicyConfig_Mitigations_TargetsConflictsWithOverrides(resourceName, name),
				ExpectError: regexp.MustCompile(`Attribute "targets" cannot be specified when "overrides" is specified`),
			},
			// Mutual exclusion: fallback + overrides
			{
				Config:      testAccRiskPolicyConfig_Mitigations_FallbackConflictsWithOverrides(resourceName, name),
				ExpectError: regexp.MustCompile(`Attribute "fallback" cannot be specified when "overrides" is specified`),
			},
			// mitigations requires fallback
			{
				Config:      testAccRiskPolicyConfig_Mitigations_MissingFallback(resourceName, name),
				ExpectError: regexp.MustCompile(`Attribute "fallback" must be specified when "mitigations" is specified`),
			},
			// targets requires fallback
			{
				Config:      testAccRiskPolicyConfig_Mitigations_TargetsMissingFallback(resourceName, name),
				ExpectError: regexp.MustCompile(`Attribute "fallback" must be specified when "targets" is specified`),
			},
			// targets.condition.and must have at least one member
			{
				Config:      testAccRiskPolicyConfig_Mitigations_TargetsEmptyAnd(resourceName, name),
				ExpectError: regexp.MustCompile(`Attribute targets.condition.and (list|set) must contain at least 1 elements`),
			},
			// mitigation condition.type cannot be IP_RANGE
			{
				Config:      testAccRiskPolicyConfig_Mitigations_ConditionIPRange(resourceName, name),
				ExpectError: regexp.MustCompile(`mitigations.*condition\.type value must be one of`),
			},
		},
	})
}

func testAccRiskPolicyConfig_Mitigations_withFallback(resourceName, name, fallback string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_mfa_device_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"

  sms = {
    enabled = true
  }

  voice = {
    enabled = false
  }

  email = {
    enabled = false
  }

  mobile = {
    enabled = false
  }

  totp = {
    enabled = false
  }

  fido2 = {
    enabled = false
  }
}

resource "pingone_mfa_device_policy" "%[2]s_reg" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s-reg"

  sms = {
    enabled = true
  }

  voice = {
    enabled = false
  }

  email = {
    enabled = false
  }

  mobile = {
    enabled = false
  }

  totp = {
    enabled = false
  }

  fido2 = {
    enabled = false
  }
}

resource "pingone_verify_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = "%[3]s"
  description    = "%[3]s"

  facial_comparison = {
    verify    = "REQUIRED"
    threshold = "LOW"
  }
}

resource "pingone_risk_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  policy_scores = {
    policy_threshold_medium = {
      min_score = 35
    }

    policy_threshold_high = {
      min_score = 70
    }

    predictors = [
      {
        compact_name = "ipRisk"
        score        = 45
      },
      {
        compact_name = "geoVelocity"
        score        = 45
      }
    ]
  }

  mitigations = [
    {
      action        = "CUSTOM"
      custom_action = "customActionValue"

      condition = {
        type         = "VALUE_COMPARISON"
        compact_name = "anonymousNetwork"
        equals       = "HIGH"
      }
    },

    {
      action                       = "MFA"
      mfa_authentication_policy_id = pingone_mfa_device_policy.%[2]s.id
      mfa_registration_policy_id   = pingone_mfa_device_policy.%[2]s_reg.id

      condition = {
        type         = "VALUE_COMPARISON"
        compact_name = "geoVelocity"
        equals       = "HIGH"
      }
    },

    {
      action           = "VERIFY"
      verify_policy_id = pingone_verify_policy.%[2]s.id

      condition = {
        type         = "VALUE_COMPARISON"
        compact_name = "geoVelocity"
        equals       = "MEDIUM"
      }
    },

    {
      action = "DENY_AND_SUSPEND"

      condition = {
        type         = "VALUE_COMPARISON"
        compact_name = "ipRisk"
        equals       = "HIGH"
      }
    }
  ]

  %[4]s

  targets = {
    condition = {
      and = [
        {
          list     = ["AUTHENTICATION", "AUTHORIZATION"]
          contains = "$${event.flow.type}"
        },
        {
          list     = ["Sales"]
          contains = "$${event.user.groups}"
        },
      ]
    }
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name, fallback)
}

func testAccRiskPolicyConfig_Mitigations_Full(resourceName, name string) string {
	return testAccRiskPolicyConfig_Mitigations_withFallback(resourceName, name, `  fallback = {
    action = "APPROVE"
  }`)
}

func testAccRiskPolicyConfig_Mitigations_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  policy_scores = {
    policy_threshold_medium = {
      min_score = 35
    }

    policy_threshold_high = {
      min_score = 70
    }

    predictors = [
      {
        compact_name = "ipRisk"
        score        = 45
      },
      {
        compact_name = "geoVelocity"
        score        = 45
      }
    ]
  }

  mitigations = [
    {
      action = "DENY"

      condition = {
        type         = "VALUE_COMPARISON"
        compact_name = "anonymousNetwork"
        equals       = "HIGH"
      }
    }
  ]

  fallback = {
    action = "DENY"
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPolicyConfig_Mitigations_ConflictsWithOverrides(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  policy_scores = {
    policy_threshold_medium = {
      min_score = 35
    }

    policy_threshold_high = {
      min_score = 70
    }

    predictors = [
      {
        compact_name = "ipRisk"
        score        = 45
      }
    ]
  }

  overrides = [
    {
      result = {
        level = "HIGH"
      }

      condition = {
        type         = "VALUE_COMPARISON"
        compact_name = "anonymousNetwork"
        equals       = "HIGH"
      }
    }
  ]

  mitigations = [
    {
      action = "DENY"

      condition = {
        type         = "VALUE_COMPARISON"
        compact_name = "geoVelocity"
        equals       = "HIGH"
      }
    }
  ]

  fallback = {
    action = "DENY"
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPolicyConfig_Mitigations_TargetsConflictsWithOverrides(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  policy_scores = {
    policy_threshold_medium = {
      min_score = 35
    }

    policy_threshold_high = {
      min_score = 70
    }

    predictors = [
      {
        compact_name = "ipRisk"
        score        = 45
      }
    ]
  }

  overrides = [
    {
      result = {
        level = "HIGH"
      }

      condition = {
        type         = "VALUE_COMPARISON"
        compact_name = "anonymousNetwork"
        equals       = "HIGH"
      }
    }
  ]

  targets = {
    condition = {
      and = [
        {
          list     = ["AUTHENTICATION"]
          contains = "$${event.flow.type}"
        },
      ]
    }
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPolicyConfig_Mitigations_FallbackConflictsWithOverrides(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  policy_scores = {
    policy_threshold_medium = {
      min_score = 35
    }

    policy_threshold_high = {
      min_score = 70
    }

    predictors = [
      {
        compact_name = "ipRisk"
        score        = 45
      }
    ]
  }

  overrides = [
    {
      result = {
        level = "HIGH"
      }

      condition = {
        type         = "VALUE_COMPARISON"
        compact_name = "anonymousNetwork"
        equals       = "HIGH"
      }
    }
  ]

  fallback = {
    action = "DENY"
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPolicyConfig_Mitigations_MissingFallback(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  policy_scores = {
    policy_threshold_medium = {
      min_score = 35
    }

    policy_threshold_high = {
      min_score = 70
    }

    predictors = [
      {
        compact_name = "ipRisk"
        score        = 45
      }
    ]
  }

  mitigations = [
    {
      action = "DENY"

      condition = {
        type         = "VALUE_COMPARISON"
        compact_name = "anonymousNetwork"
        equals       = "HIGH"
      }
    }
  ]
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPolicyConfig_Mitigations_TargetsMissingFallback(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  policy_scores = {
    policy_threshold_medium = {
      min_score = 35
    }

    policy_threshold_high = {
      min_score = 70
    }

    predictors = [
      {
        compact_name = "ipRisk"
        score        = 45
      }
    ]
  }

  targets = {
    condition = {
      and = [
        {
          list     = ["AUTHENTICATION"]
          contains = "$${event.flow.type}"
        },
      ]
    }
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPolicyConfig_Mitigations_TargetsEmptyAnd(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  policy_scores = {
    policy_threshold_medium = {
      min_score = 35
    }

    policy_threshold_high = {
      min_score = 70
    }

    predictors = [
      {
        compact_name = "ipRisk"
        score        = 45
      }
    ]
  }

  fallback = {
    action = "DENY"
  }

  targets = {
    condition = {
      and = []
    }
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPolicyConfig_Mitigations_ConditionIPRange(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  policy_scores = {
    policy_threshold_medium = {
      min_score = 35
    }

    policy_threshold_high = {
      min_score = 70
    }

    predictors = [
      {
        compact_name = "ipRisk"
        score        = 45
      }
    ]
  }

  mitigations = [
    {
      action = "DENY"

      condition = {
        type = "IP_RANGE"
      }
    }
  ]

  fallback = {
    action = "DENY"
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPolicyConfig_Mitigations_FallbackMFA(resourceName, name string) string {
	return testAccRiskPolicyConfig_Mitigations_withFallback(resourceName, name, fmt.Sprintf(`fallback = {
    action                       = "MFA"
    mfa_authentication_policy_id = pingone_mfa_device_policy.%s.id
  }`, resourceName))
}

func testAccRiskPolicyConfig_Mitigations_FallbackOnly(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  policy_scores = {
    policy_threshold_medium = {
      min_score = 35
    }

    policy_threshold_high = {
      min_score = 70
    }

    predictors = [
      {
        compact_name = "ipRisk"
        score        = 45
      }
    ]
  }

  fallback = {
    action = "DENY"
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func TestAccRiskPolicy_OverridesMitigationsChangeType(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_risk_policy.%s", resourceName)

	name := resourceName

	overridesCheck := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "overrides.#", "1"),
		resource.TestCheckNoResourceAttr(resourceFullName, "mitigations.#"),
		resource.TestCheckNoResourceAttr(resourceFullName, "fallback.action"),
	)

	migrationsCheck := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "mitigations.#", "1"),
		resource.TestCheckResourceAttr(resourceFullName, "fallback.action", "DENY"),
		resource.TestCheckNoResourceAttr(resourceFullName, "overrides.#"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckTestAccFlaky(t) // PND-5900: policy PUT reverts to previous config on next GET, leaving a non-empty refresh plan
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             risk.RiskPolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Start with overrides
			{
				Config: testAccRiskPolicyConfig_OverridesForTypeChange(resourceName, name),
				Check:  overridesCheck,
			},
			// Switch to mitigations
			{
				Config: testAccRiskPolicyConfig_MitigationsForTypeChange(resourceName, name),
				Check:  migrationsCheck,
			},
			// Switch back to overrides
			{
				Config: testAccRiskPolicyConfig_OverridesForTypeChange(resourceName, name),
				Check:  overridesCheck,
			},
		},
	})
}

func testAccRiskPolicyConfig_MitigationsForTypeChange(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  policy_scores = {
    policy_threshold_medium = {
      min_score = 35
    }

    policy_threshold_high = {
      min_score = 70
    }

    predictors = [
      {
        compact_name = "ipRisk"
        score        = 45
      },
      {
        compact_name = "geoVelocity"
        score        = 45
      }
    ]
  }

  mitigations = [
    {
      action = "DENY"

      condition = {
        type         = "VALUE_COMPARISON"
        compact_name = "geoVelocity"
        equals       = "HIGH"
      }
    }
  ]

  fallback = {
    action = "DENY"
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPolicyConfig_OverridesForTypeChange(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name = "%[3]s"

  policy_scores = {
    policy_threshold_medium = {
      min_score = 35
    }

    policy_threshold_high = {
      min_score = 70
    }

    predictors = [
      {
        compact_name = "ipRisk"
        score        = 45
      },
      {
        compact_name = "geoVelocity"
        score        = 45
      }
    ]
  }

  overrides = [
    {
      result = {
        level = "MEDIUM"
      }

      condition = {
        type         = "VALUE_COMPARISON"
        compact_name = "geoVelocity"
        equals       = "HIGH"
      }
    }
  ]
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
