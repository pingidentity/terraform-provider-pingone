package risk_test

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
	"github.com/patrickcping/pingone-go-sdk-v2/risk"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func testAccCheckRiskPolicyDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.RiskAPIClient
	ctx = context.WithValue(ctx, risk.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})

	apiClientManagement := p1Client.API.ManagementAPIClient
	ctxManagement := context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_risk_policy" {
			continue
		}

		_, rEnv, err := apiClientManagement.EnvironmentsApi.ReadOneEnvironment(ctxManagement, rs.Primary.Attributes["environment_id"]).Execute()

		if err != nil {

			if rEnv == nil {
				return fmt.Errorf("Response object does not exist and no error detected")
			}

			if rEnv.StatusCode == 404 {
				continue
			}

			return err
		}

		body, r, err := apiClient.RiskPoliciesApi.ReadOneRiskPolicySet(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.ID).Execute()

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

		return fmt.Errorf("PingOne risk policy %s still exists", rs.Primary.ID)
	}

	return nil
}

func TestAccRiskPolicy_NewEnv(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_risk_policy.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckRiskPolicyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRiskPolicyConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
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
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
		resource.TestCheckResourceAttr(resourceFullName, "name", name),
		resource.TestCheckResourceAttr(resourceFullName, "default_result.type", "VALUE"),
		resource.TestCheckResourceAttr(resourceFullName, "default_result.level", "LOW"),
		resource.TestCheckResourceAttr(resourceFullName, "default", "false"),
		resource.TestCheckResourceAttr(resourceFullName, "evaluated_predictors.#", "3"),
		resource.TestMatchResourceAttr(resourceFullName, "evaluated_predictors.0", verify.P1ResourceIDRegexp),
		resource.TestMatchResourceAttr(resourceFullName, "evaluated_predictors.1", verify.P1ResourceIDRegexp),
		resource.TestMatchResourceAttr(resourceFullName, "evaluated_predictors.2", verify.P1ResourceIDRegexp),
		resource.TestCheckResourceAttr(resourceFullName, "overrides.#", "0"),
	)

	minimalCheck := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
		resource.TestCheckResourceAttr(resourceFullName, "name", name),
		resource.TestCheckResourceAttr(resourceFullName, "default_result.type", "VALUE"),
		resource.TestCheckResourceAttr(resourceFullName, "default_result.level", "LOW"),
		resource.TestCheckResourceAttr(resourceFullName, "default", "false"),
		resource.TestMatchResourceAttr(resourceFullName, "evaluated_predictors.#", regexp.MustCompile(`^(?:[2-9]|[12]\d)\d*$`)),
		resource.TestCheckResourceAttr(resourceFullName, "overrides.#", "0"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckRiskPolicyDestroy,
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
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckRiskPolicyDestroy,
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
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckRiskPolicyDestroy,
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
	)

	weightsCheck := resource.ComposeTestCheckFunc(
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
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckRiskPolicyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRiskPolicyConfig_Scores_Full(resourceName, name),
				Check:  scoresCheck,
			},
			{
				Config: testAccRiskPolicyConfig_Weights_Full(resourceName, name),
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
			"condition.compact_name":              "ipVelocityByUser",
			"condition.predictor_reference_value": "${details.ipVelocityByUser.level}",
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
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckRiskPolicyDestroy,
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
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
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
        compact_name = "ipVelocityByUser"
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
