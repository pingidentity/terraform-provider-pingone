// Copyright © 2026 Ping Identity Corporation

package risk_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/risk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccRiskPolicyDataSource_RiskPolicyID(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_risk_policy.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.pingone_risk_policy.%s", resourceName)

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
			{
				Config: testAccRiskPolicyDataSourceConfig_ByID(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "risk_policy_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "name", resourceFullName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "default", resourceFullName, "default"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "default_result.level", resourceFullName, "default_result.level"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "default_result.type", resourceFullName, "default_result.type"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "evaluated_predictors.#", resourceFullName, "evaluated_predictors.#"),
					resource.TestCheckResourceAttr(dataSourceFullName, "policy_scores.policy_threshold_medium.min_score", "35"),
					resource.TestCheckResourceAttr(dataSourceFullName, "policy_scores.policy_threshold_medium.max_score", "70"),
					resource.TestCheckResourceAttr(dataSourceFullName, "policy_scores.policy_threshold_high.min_score", "70"),
					resource.TestCheckResourceAttr(dataSourceFullName, "policy_scores.policy_threshold_high.max_score", "1000"),
					resource.TestCheckResourceAttr(dataSourceFullName, "policy_scores.predictors.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourceFullName, "policy_scores.predictors.*", map[string]string{
						"compact_name":              fmt.Sprintf("%s2", name),
						"predictor_reference_value": fmt.Sprintf("${details.%s2.level}", name),
						"score":                     "45",
					}),
				),
			},
		},
	})
}

func TestAccRiskPolicyDataSource_Name(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_risk_policy.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.pingone_risk_policy.%s", resourceName)

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
			{
				Config: testAccRiskPolicyDataSourceConfig_ByName(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "risk_policy_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "name", resourceFullName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "default_result.level", resourceFullName, "default_result.level"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "policy_scores.policy_threshold_medium.min_score", resourceFullName, "policy_scores.policy_threshold_medium.min_score"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "policy_scores.predictors.#", resourceFullName, "policy_scores.predictors.#"),
				),
			},
		},
	})
}

func TestAccRiskPolicyDataSource_Weights(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_risk_policy.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.pingone_risk_policy.%s", resourceName)

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
			{
				Config: testAccRiskPolicyDataSourceConfig_Weights(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "name", resourceFullName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "policy_weights.policy_threshold_medium.min_score", resourceFullName, "policy_weights.policy_threshold_medium.min_score"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "policy_weights.policy_threshold_high.min_score", resourceFullName, "policy_weights.policy_threshold_high.min_score"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "policy_weights.predictors.#", resourceFullName, "policy_weights.predictors.#"),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourceFullName, "policy_weights.predictors.*", map[string]string{
						"compact_name": fmt.Sprintf("%s2", name),
						"weight":       "3",
					}),
				),
			},
		},
	})
}

func TestAccRiskPolicyDataSource_Overrides(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_risk_policy.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.pingone_risk_policy.%s", resourceName)

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
			{
				Config: testAccRiskPolicyDataSourceConfig_Overrides(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "name", resourceFullName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "overrides.#", resourceFullName, "overrides.#"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "overrides.0.result.level", resourceFullName, "overrides.0.result.level"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "overrides.0.condition.type", resourceFullName, "overrides.0.condition.type"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "overrides.0.condition.compact_name", resourceFullName, "overrides.0.condition.compact_name"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "overrides.0.condition.equals", resourceFullName, "overrides.0.condition.equals"),
				),
			},
		},
	})
}

func TestAccRiskPolicyDataSource_Mitigations(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_risk_policy.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.pingone_risk_policy.%s", resourceName)

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
			{
				Config: testAccRiskPolicyDataSourceConfig_Mitigations(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "name", resourceFullName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "mitigations.#", resourceFullName, "mitigations.#"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "mitigations.0.action", resourceFullName, "mitigations.0.action"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "mitigations.0.custom_action", resourceFullName, "mitigations.0.custom_action"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "mitigations.0.mfa_authentication_policy_id", resourceFullName, "mitigations.0.mfa_authentication_policy_id"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "fallback.action", resourceFullName, "fallback.action"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "targets.condition.and.#", resourceFullName, "targets.condition.and.#"),
				),
			},
		},
	})
}

func TestAccRiskPolicyDataSource_NotFound(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

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
			{
				Config:      testAccRiskPolicyDataSourceConfig_NotFoundByName(resourceName),
				ExpectError: regexp.MustCompile("Risk Policy with name .* not found"),
			},
			{
				Config:      testAccRiskPolicyDataSourceConfig_NotFoundByID(resourceName),
				ExpectError: regexp.MustCompile("Error when calling `ReadOneRiskPolicySet`"),
			},
		},
	})
}

func TestAccRiskPolicyDataSource_Validation(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

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
			{
				Config:      testAccRiskPolicyDataSourceConfig_IDAndName(resourceName),
				ExpectError: regexp.MustCompile("Invalid Attribute Combination"),
			},
			{
				Config:      testAccRiskPolicyDataSourceConfig_EmptyName(resourceName),
				ExpectError: regexp.MustCompile(`Attribute name string length must be at least 1, got: 0`),
			},
		},
	})
}

func testAccRiskPolicyDataSourceConfig_ByID(resourceName, name string) string {
	return fmt.Sprintf(`%[1]s

data "pingone_risk_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  risk_policy_id = pingone_risk_policy.%[2]s.id
}`, testAccRiskPolicyConfig_Scores_Minimal(resourceName, name), resourceName)
}

func testAccRiskPolicyDataSourceConfig_ByName(resourceName, name string) string {
	return fmt.Sprintf(`%[1]s

data "pingone_risk_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = pingone_risk_policy.%[2]s.name
}`, testAccRiskPolicyConfig_Scores_Minimal(resourceName, name), resourceName)
}

func testAccRiskPolicyDataSourceConfig_Weights(resourceName, name string) string {
	return fmt.Sprintf(`%[1]s

data "pingone_risk_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = pingone_risk_policy.%[2]s.name
}`, testAccRiskPolicyConfig_Weights_Minimal(resourceName, name), resourceName)
}

func testAccRiskPolicyDataSourceConfig_Overrides(resourceName, name string) string {
	return fmt.Sprintf(`%[1]s

data "pingone_risk_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = pingone_risk_policy.%[2]s.name
}`, testAccRiskPolicyConfig_Overrides_Minimal(resourceName, name), resourceName)
}

func testAccRiskPolicyDataSourceConfig_Mitigations(resourceName, name string) string {
	return fmt.Sprintf(`%[1]s

data "pingone_risk_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = pingone_risk_policy.%[2]s.name
}`, testAccRiskPolicyConfig_Mitigations_Full(resourceName, name), resourceName)
}

func testAccRiskPolicyDataSourceConfig_NotFoundByName(resourceName string) string {
	return fmt.Sprintf(`%[1]s

data "pingone_risk_policy" "%[2]s-notfound" {
  environment_id = data.pingone_environment.general_test.id
  name           = "test_not_found_name"
}`, acctest.GenericSandboxEnvironment(), resourceName)
}

func testAccRiskPolicyDataSourceConfig_NotFoundByID(resourceName string) string {
	return fmt.Sprintf(`%[1]s

data "pingone_risk_policy" "%[2]s-notfound" {
  environment_id = data.pingone_environment.general_test.id
  risk_policy_id = "9c052a8a-14be-44e4-8f07-2662569994ce"
}`, acctest.GenericSandboxEnvironment(), resourceName)
}

func testAccRiskPolicyDataSourceConfig_IDAndName(resourceName string) string {
	return fmt.Sprintf(`%[1]s

data "pingone_risk_policy" "%[2]s-validation" {
  environment_id = data.pingone_environment.general_test.id
  risk_policy_id = "9c052a8a-14be-44e4-8f07-2662569994ce"
  name           = "test_validation_name"
}`, acctest.GenericSandboxEnvironment(), resourceName)
}

func testAccRiskPolicyDataSourceConfig_EmptyName(resourceName string) string {
	return fmt.Sprintf(`%[1]s

data "pingone_risk_policy" "%[2]s-validation" {
  environment_id = data.pingone_environment.general_test.id
  name           = ""
}`, acctest.GenericSandboxEnvironment(), resourceName)
}
