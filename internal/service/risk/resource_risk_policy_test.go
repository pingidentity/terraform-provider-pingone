// Copyright © 2025 Ping Identity Corporation

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
		// PreCheck: func() {
		//  acctest.PreCheckNoTestAccFlaky(t)
		// 	acctest.PreCheckClient(t)
		// 	acctest.PreCheckNoBeta(t)
		//	acctest.PreCheckNoBeta(t)
		// },
		PreCheck:                 func() { t.Skipf("PND-5900") },
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
		// PreCheck: func() {
		//  acctest.PreCheckNoTestAccFlaky(t)
		// 	acctest.PreCheckClient(t)
		// 	acctest.PreCheckNoBeta(t)
		//	acctest.PreCheckNoBeta(t)
		// },
		PreCheck:                 func() { t.Skipf("PND-5900") },
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
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
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
		// PreCheck: func() {
		//	acctest.PreCheckNoTestAccFlaky(t)
		// 	acctest.PreCheckClient(t)
		// 	acctest.PreCheckNoBeta(t)
		//	acctest.PreCheckNoBeta(t)
		// },
		PreCheck:                 func() { t.Skipf("PND-5900") },
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
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
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
		// PreCheck: func() {
		//	acctest.PreCheckNoTestAccFlaky(t)
		// 	acctest.PreCheckClient(t)
		// 	acctest.PreCheckNoBeta(t)
		//	acctest.PreCheckNoBeta(t)
		// },
		PreCheck:                 func() { t.Skipf("PND-5900") },
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
		// PreCheck: func() {
		//	acctest.PreCheckNoTestAccFlaky(t)
		// 	acctest.PreCheckClient(t)
		// 	acctest.PreCheckNoBeta(t)
		//	acctest.PreCheckNoBeta(t)
		// },
		PreCheck:                 func() { t.Skipf("PND-5900") },
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

func TestAccRiskPolicy_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_risk_policy.%s", resourceName)

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
