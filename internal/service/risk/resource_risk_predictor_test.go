// Copyright © 2026 Ping Identity Corporation

package risk_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	acctestlegacysdk "github.com/pingidentity/terraform-provider-pingone/internal/acctest/legacysdk"
	baselegacysdk "github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/base/legacysdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/risk"
	"github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccRiskPredictor_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_risk_predictor.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	var riskPredictorID, environmentID string

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
		CheckDestroy:             risk.RiskPredictor_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccRiskPredictorConfig_Minimal(resourceName, name),
				Check:  risk.RiskPredictor_GetIDs(resourceFullName, &environmentID, &riskPredictorID),
			},
			{
				PreConfig: func() {
					risk.RiskPredictor_RemovalDrift_PreConfig(ctx, p1Client.API.RiskAPIClient, t, environmentID, riskPredictorID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the environment
			{
				Config: testAccRiskPredictorConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check:  risk.RiskPredictor_GetIDs(resourceFullName, &environmentID, &riskPredictorID),
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

func TestAccRiskPredictor_NewEnv(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_risk_predictor.%s", resourceName)

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
		CheckDestroy:             risk.RiskPredictor_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRiskPredictorConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
				),
			},
		},
	})
}

func TestAccRiskPredictor_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_risk_predictor.%s", resourceName)

	name := resourceName

	fullCheck := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "name", name),
		resource.TestCheckResourceAttr(resourceFullName, "compact_name", fmt.Sprintf("%s1", name)),
		resource.TestCheckResourceAttr(resourceFullName, "description", "When my wife is upset, I let her colour in my black and white tattoos.  She just needs a shoulder to crayon.."),
		resource.TestCheckResourceAttr(resourceFullName, "type", "ANONYMOUS_NETWORK"),
		resource.TestCheckResourceAttr(resourceFullName, "licensed", "true"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "true"),
		resource.TestCheckResourceAttr(resourceFullName, "default.result.level", "MEDIUM"),
	)

	minimalCheck := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "name", name),
		resource.TestCheckResourceAttr(resourceFullName, "compact_name", fmt.Sprintf("%s1", name)),
		resource.TestCheckNoResourceAttr(resourceFullName, "description"),
		resource.TestCheckResourceAttr(resourceFullName, "type", "ANONYMOUS_NETWORK"),
		resource.TestCheckResourceAttr(resourceFullName, "licensed", "true"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "true"),
		resource.TestCheckNoResourceAttr(resourceFullName, "default.result.level"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             risk.RiskPredictor_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			{
				Config: testAccRiskPredictorConfig_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config:  testAccRiskPredictorConfig_Full(resourceName, name),
				Destroy: true,
			},
			// Minimal
			{
				Config: testAccRiskPredictorConfig_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config:  testAccRiskPredictorConfig_Minimal(resourceName, name),
				Destroy: true,
			},
			// Change
			{
				Config: testAccRiskPredictorConfig_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config: testAccRiskPredictorConfig_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config: testAccRiskPredictorConfig_Full(resourceName, name),
				Check:  fullCheck,
			},
		},
	})
}

func TestAccRiskPredictor_Composite(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_risk_predictor.%s", resourceName)

	name := resourceName

	fullCheck1 := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "type", "COMPOSITE"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "true"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_composite.compositions.0.condition_json", "{\"not\":{\"or\":[{\"equals\":0,\"type\":\"VALUE_COMPARISON\",\"value\":\"${details.counters.predictorLevels.medium}\"},{\"equals\":\"High\",\"type\":\"VALUE_COMPARISON\",\"value\":\"${details.geoVelocity.level}\"},{\"startsWith\":\"admin\",\"type\":\"VALUE_COMPARISON\",\"value\":\"${event.user.name}\"},{\"endsWith\":\"@example.com\",\"type\":\"VALUE_COMPARISON\",\"value\":\"${event.user.name}\"},{\"and\":[{\"equals\":\"High\",\"type\":\"VALUE_COMPARISON\",\"value\":\"${details.anonymousNetwork.level}\"},{\"contains\":\"${event.user.groups}\",\"list\":[\"Group Name\"],\"type\":\"GROUPS_INTERSECTION\"}],\"type\":\"AND\"}],\"type\":\"OR\"},\"type\":\"NOT\"}"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_composite.compositions.0.condition", "{\"not\":{\"or\":[{\"equals\":0,\"type\":\"VALUE_COMPARISON\",\"value\":\"${details.counters.predictorLevels.medium}\"},{\"equals\":\"High\",\"type\":\"VALUE_COMPARISON\",\"value\":\"${details.geoVelocity.level}\"},{\"startsWith\":\"admin\",\"type\":\"VALUE_COMPARISON\",\"value\":\"${event.user.name}\"},{\"endsWith\":\"@example.com\",\"type\":\"VALUE_COMPARISON\",\"value\":\"${event.user.name}\"},{\"and\":[{\"equals\":\"High\",\"type\":\"VALUE_COMPARISON\",\"value\":\"${details.anonymousNetwork.level}\"},{\"contains\":\"${event.user.groups}\",\"list\":[\"Group Name\"],\"type\":\"GROUPS_INTERSECTION\"}],\"type\":\"AND\"}],\"type\":\"OR\"},\"type\":\"NOT\"}"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_composite.compositions.0.level", "HIGH"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_composite.compositions.1.condition_json", "{\"and\":[{\"equals\":5,\"type\":\"VALUE_COMPARISON\",\"value\":\"${details.counters.predictorLevels.medium}\"},{\"equals\":\"low\",\"type\":\"VALUE_COMPARISON\",\"value\":\"${details.anonymousNetwork.level}\"},{\"startsWith\":\"test\",\"type\":\"VALUE_COMPARISON\",\"value\":\"${event.user.name}\"},{\"and\":[{\"equals\":\"high\",\"type\":\"VALUE_COMPARISON\",\"value\":\"${details.anonymousNetwork.level}\"},{\"or\":[{\"notEquals\":\"high\",\"type\":\"VALUE_COMPARISON\",\"value\":\"${details.anonymousNetwork.level}\"}]}]}]}"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_composite.compositions.1.condition", "{\"and\":[{\"equals\":5,\"type\":\"VALUE_COMPARISON\",\"value\":\"${details.counters.predictorLevels.medium}\"},{\"equals\":\"low\",\"type\":\"VALUE_COMPARISON\",\"value\":\"${details.anonymousNetwork.level}\"},{\"startsWith\":\"test\",\"type\":\"VALUE_COMPARISON\",\"value\":\"${event.user.name}\"},{\"and\":[{\"equals\":\"high\",\"type\":\"VALUE_COMPARISON\",\"value\":\"${details.anonymousNetwork.level}\"},{\"or\":[{\"notEquals\":\"high\",\"type\":\"VALUE_COMPARISON\",\"value\":\"${details.anonymousNetwork.level}\"}],\"type\":\"OR\"}],\"type\":\"AND\"}],\"type\":\"AND\"}"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_composite.compositions.1.level", "LOW"),
	)

	fullCheck2 := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "type", "COMPOSITE"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "true"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_composite.compositions.0.condition_json", "{\"and\":[{\"equals\":5,\"type\":\"VALUE_COMPARISON\",\"value\":\"${details.counters.predictorLevels.medium}\"},{\"equals\":\"low\",\"type\":\"VALUE_COMPARISON\",\"value\":\"${details.anonymousNetwork.level}\"},{\"endsWith\":\"@test.org\",\"type\":\"VALUE_COMPARISON\",\"value\":\"${event.user.name}\"},{\"and\":[{\"equals\":\"high\",\"type\":\"VALUE_COMPARISON\",\"value\":\"${details.anonymousNetwork.level}\"},{\"or\":[{\"notEquals\":\"high\",\"type\":\"VALUE_COMPARISON\",\"value\":\"${details.anonymousNetwork.level}\"}]}]}]}"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_composite.compositions.0.condition", "{\"and\":[{\"equals\":5,\"type\":\"VALUE_COMPARISON\",\"value\":\"${details.counters.predictorLevels.medium}\"},{\"equals\":\"low\",\"type\":\"VALUE_COMPARISON\",\"value\":\"${details.anonymousNetwork.level}\"},{\"endsWith\":\"@test.org\",\"type\":\"VALUE_COMPARISON\",\"value\":\"${event.user.name}\"},{\"and\":[{\"equals\":\"high\",\"type\":\"VALUE_COMPARISON\",\"value\":\"${details.anonymousNetwork.level}\"},{\"or\":[{\"notEquals\":\"high\",\"type\":\"VALUE_COMPARISON\",\"value\":\"${details.anonymousNetwork.level}\"}],\"type\":\"OR\"}],\"type\":\"AND\"}],\"type\":\"AND\"}"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_composite.compositions.0.level", "LOW"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_composite.compositions.1.condition_json", "{\"not\":{\"or\":[{\"equals\":0,\"type\":\"VALUE_COMPARISON\",\"value\":\"${details.counters.predictorLevels.medium}\"},{\"equals\":\"High\",\"type\":\"VALUE_COMPARISON\",\"value\":\"${details.geoVelocity.level}\"},{\"startsWith\":\"user\",\"type\":\"VALUE_COMPARISON\",\"value\":\"${event.user.name}\"},{\"and\":[{\"equals\":\"High\",\"type\":\"VALUE_COMPARISON\",\"value\":\"${details.anonymousNetwork.level}\"}],\"type\":\"AND\"}],\"type\":\"OR\"},\"type\":\"NOT\"}"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_composite.compositions.1.condition", "{\"not\":{\"or\":[{\"equals\":0,\"type\":\"VALUE_COMPARISON\",\"value\":\"${details.counters.predictorLevels.medium}\"},{\"equals\":\"High\",\"type\":\"VALUE_COMPARISON\",\"value\":\"${details.geoVelocity.level}\"},{\"startsWith\":\"user\",\"type\":\"VALUE_COMPARISON\",\"value\":\"${event.user.name}\"},{\"and\":[{\"equals\":\"High\",\"type\":\"VALUE_COMPARISON\",\"value\":\"${details.anonymousNetwork.level}\"}],\"type\":\"AND\"}],\"type\":\"OR\"},\"type\":\"NOT\"}"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_composite.compositions.1.level", "HIGH"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             risk.RiskPredictor_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			{
				Config: testAccRiskPredictorConfig_Composite_Full_1(resourceName, name),
				Check:  fullCheck1,
			},
			{
				Config:  testAccRiskPredictorConfig_Composite_Full_1(resourceName, name),
				Destroy: true,
			},
			// Minimal
			{
				Config: testAccRiskPredictorConfig_Composite_Full_2(resourceName, name),
				Check:  fullCheck2,
			},
			{
				Config:  testAccRiskPredictorConfig_Composite_Full_2(resourceName, name),
				Destroy: true,
			},
			// Change
			{
				Config: testAccRiskPredictorConfig_Composite_Full_1(resourceName, name),
				Check:  fullCheck1,
			},
			{
				Config: testAccRiskPredictorConfig_Composite_Full_2(resourceName, name),
				Check:  fullCheck2,
			},
			{
				Config: testAccRiskPredictorConfig_Composite_Full_1(resourceName, name),
				Check:  fullCheck1,
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
				ImportStateVerifyIgnore: []string{
					"predictor_composite.compositions.1.condition_json",
				},
			},
			{
				Config:  testAccRiskPredictorConfig_Composite_Full_1(resourceName, name),
				Destroy: true,
			},
			// Error
			{
				Config:      testAccRiskPredictorConfig_Composite_InvalidJSON(resourceName, name),
				ExpectError: regexp.MustCompile(`Cannot parse the condition input JSON`),
			},
		},
	})
}

func TestAccRiskPredictor_Adversary_In_The_Middle(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_risk_predictor.%s", resourceName)

	name := resourceName

	fullCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "type", "ADVERSARY_IN_THE_MIDDLE"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "true"),
		resource.TestCheckResourceAttr(resourceFullName, "default.result.level", "MEDIUM"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_adversary_in_the_middle.allowed_domain_list.#", "3"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "predictor_adversary_in_the_middle.allowed_domain_list.*", "domain1.com"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "predictor_adversary_in_the_middle.allowed_domain_list.*", "domain2.com"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "predictor_adversary_in_the_middle.allowed_domain_list.*", "domain3.com"),
	)

	minimalCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "type", "ADVERSARY_IN_THE_MIDDLE"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "true"),
		resource.TestCheckNoResourceAttr(resourceFullName, "default.result.level"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_adversary_in_the_middle.allowed_domain_list.#", "0"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             risk.RiskPredictor_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			{
				Config: testAccRiskPredictorConfig_Adversary_In_The_Middle_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config:  testAccRiskPredictorConfig_Adversary_In_The_Middle_Full(resourceName, name),
				Destroy: true,
			},
			// Minimal
			{
				Config: testAccRiskPredictorConfig_Adversary_In_The_Middle_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config:  testAccRiskPredictorConfig_Adversary_In_The_Middle_Minimal(resourceName, name),
				Destroy: true,
			},
			// Change
			{
				Config: testAccRiskPredictorConfig_Adversary_In_The_Middle_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config: testAccRiskPredictorConfig_Adversary_In_The_Middle_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config: testAccRiskPredictorConfig_Adversary_In_The_Middle_Full(resourceName, name),
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

func TestAccRiskPredictor_Adversary_In_The_Middle_OverwriteUndeletable(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_risk_predictor.%s", resourceName)

	name := resourceName
	compactName := "adversaryInTheMiddle"

	fullCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "name", name),
		resource.TestCheckResourceAttr(resourceFullName, "compact_name", "adversaryInTheMiddle"),
		resource.TestCheckResourceAttr(resourceFullName, "type", "ADVERSARY_IN_THE_MIDDLE"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "false"),
		resource.TestCheckResourceAttr(resourceFullName, "default.result.level", "MEDIUM"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_adversary_in_the_middle.allowed_domain_list.#", "3"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "predictor_adversary_in_the_middle.allowed_domain_list.*", "domain1.com"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "predictor_adversary_in_the_middle.allowed_domain_list.*", "domain2.com"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "predictor_adversary_in_the_middle.allowed_domain_list.*", "domain3.com"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             risk.RiskPredictor_CheckDestroyUndeletable,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			{
				Config: testAccRiskPredictorConfig_Adversary_In_The_Middle_OverwriteUndeletable(resourceName, name, compactName),
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

func TestAccRiskPredictor_Anonymous_Network(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_risk_predictor.%s", resourceName)

	name := resourceName

	fullCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "type", "ANONYMOUS_NETWORK"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "true"),
		resource.TestCheckResourceAttr(resourceFullName, "default.result.level", "MEDIUM"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_anonymous_network.allowed_cidr_list.#", "3"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "predictor_anonymous_network.allowed_cidr_list.*", "192.168.0.0/24"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "predictor_anonymous_network.allowed_cidr_list.*", "10.0.0.0/8"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "predictor_anonymous_network.allowed_cidr_list.*", "172.16.0.0/12"),
	)

	minimalCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "type", "ANONYMOUS_NETWORK"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "true"),
		resource.TestCheckNoResourceAttr(resourceFullName, "default.result.level"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_anonymous_network.allowed_cidr_list.#", "0"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             risk.RiskPredictor_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			{
				Config: testAccRiskPredictorConfig_Anonymous_Network_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config:  testAccRiskPredictorConfig_Anonymous_Network_Full(resourceName, name),
				Destroy: true,
			},
			// Minimal
			{
				Config: testAccRiskPredictorConfig_Anonymous_Network_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config:  testAccRiskPredictorConfig_Anonymous_Network_Minimal(resourceName, name),
				Destroy: true,
			},
			// Change
			{
				Config: testAccRiskPredictorConfig_Anonymous_Network_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config: testAccRiskPredictorConfig_Anonymous_Network_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config: testAccRiskPredictorConfig_Anonymous_Network_Full(resourceName, name),
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

func TestAccRiskPredictor_Anonymous_Network_OverwriteUndeletable(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_risk_predictor.%s", resourceName)

	name := resourceName
	compactName := "anonymousNetwork"

	fullCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "name", name),
		resource.TestCheckResourceAttr(resourceFullName, "compact_name", "anonymousNetwork"),
		resource.TestCheckResourceAttr(resourceFullName, "type", "ANONYMOUS_NETWORK"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "false"),
		resource.TestCheckResourceAttr(resourceFullName, "default.result.level", "MEDIUM"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_anonymous_network.allowed_cidr_list.#", "3"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "predictor_anonymous_network.allowed_cidr_list.*", "192.168.0.0/24"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "predictor_anonymous_network.allowed_cidr_list.*", "10.0.0.0/8"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "predictor_anonymous_network.allowed_cidr_list.*", "172.16.0.0/12"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             risk.RiskPredictor_CheckDestroyUndeletable,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			{
				Config: testAccRiskPredictorConfig_Anonymous_Network_OverwriteUndeletable(resourceName, name, compactName),
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

func TestAccRiskPredictor_Bot_Detection(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_risk_predictor.%s", resourceName)

	name := resourceName

	fullCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "type", "BOT"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "true"),
		resource.TestCheckResourceAttr(resourceFullName, "default.result.level", "MEDIUM"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_bot_detection.include_repeated_events_without_sdk", "true"),
	)

	minimalCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "type", "BOT"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "true"),
		resource.TestCheckNoResourceAttr(resourceFullName, "default.result.level"),
		resource.TestCheckNoResourceAttr(resourceFullName, "predictor_bot_detection.include_repeated_events_without_sdk"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             risk.RiskPredictor_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			{
				Config: testAccRiskPredictorConfig_Bot_Detection_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config:  testAccRiskPredictorConfig_Bot_Detection_Full(resourceName, name),
				Destroy: true,
			},
			// Minimal
			{
				Config: testAccRiskPredictorConfig_Bot_Detection_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config:  testAccRiskPredictorConfig_Bot_Detection_Minimal(resourceName, name),
				Destroy: true,
			},
			// Change
			{
				Config: testAccRiskPredictorConfig_Bot_Detection_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config: testAccRiskPredictorConfig_Bot_Detection_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config: testAccRiskPredictorConfig_Bot_Detection_Full(resourceName, name),
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

func TestAccRiskPredictor_Bot_Detection_OverwriteUndeletable(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_risk_predictor.%s", resourceName)

	name := resourceName
	compactName := "botDetection"

	fullCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "name", name),
		resource.TestCheckResourceAttr(resourceFullName, "compact_name", "botDetection"),
		resource.TestCheckResourceAttr(resourceFullName, "type", "BOT"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "false"),
		resource.TestCheckResourceAttr(resourceFullName, "default.result.level", "MEDIUM"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_bot_detection.include_repeated_events_without_sdk", "true"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             risk.RiskPredictor_CheckDestroyUndeletable,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			{
				Config: testAccRiskPredictorConfig_Bot_Detection_OverwriteUndeletable(resourceName, name, compactName),
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

func TestAccRiskPredictor_Geovelocity(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_risk_predictor.%s", resourceName)

	name := resourceName

	fullCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "type", "GEO_VELOCITY"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "true"),
		resource.TestCheckResourceAttr(resourceFullName, "default.result.level", "MEDIUM"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_geovelocity.allowed_cidr_list.#", "3"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "predictor_geovelocity.allowed_cidr_list.*", "192.168.0.0/24"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "predictor_geovelocity.allowed_cidr_list.*", "10.0.0.0/8"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "predictor_geovelocity.allowed_cidr_list.*", "172.16.0.0/12"),
	)

	minimalCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "type", "GEO_VELOCITY"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "true"),
		resource.TestCheckNoResourceAttr(resourceFullName, "default.result.level"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_geovelocity.allowed_cidr_list.#", "0"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             risk.RiskPredictor_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			{
				Config: testAccRiskPredictorConfig_Geovelocity_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config:  testAccRiskPredictorConfig_Geovelocity_Full(resourceName, name),
				Destroy: true,
			},
			// Minimal
			{
				Config: testAccRiskPredictorConfig_Geovelocity_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config:  testAccRiskPredictorConfig_Geovelocity_Minimal(resourceName, name),
				Destroy: true,
			},
			// Change
			{
				Config: testAccRiskPredictorConfig_Geovelocity_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config: testAccRiskPredictorConfig_Geovelocity_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config: testAccRiskPredictorConfig_Geovelocity_Full(resourceName, name),
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

func TestAccRiskPredictor_Geovelocity_OverwriteUndeletable(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_risk_predictor.%s", resourceName)

	name := resourceName
	compactName := "geoVelocity"

	fullCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "name", name),
		resource.TestCheckResourceAttr(resourceFullName, "compact_name", "geoVelocity"),
		resource.TestCheckResourceAttr(resourceFullName, "type", "GEO_VELOCITY"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "false"),
		resource.TestCheckResourceAttr(resourceFullName, "default.result.level", "MEDIUM"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_geovelocity.allowed_cidr_list.#", "3"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "predictor_geovelocity.allowed_cidr_list.*", "192.168.0.0/24"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "predictor_geovelocity.allowed_cidr_list.*", "10.0.0.0/8"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "predictor_geovelocity.allowed_cidr_list.*", "172.16.0.0/12"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             risk.RiskPredictor_CheckDestroyUndeletable,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			{
				Config: testAccRiskPredictorConfig_Geovelocity_OverwriteUndeletable(resourceName, name, compactName),
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

func TestAccRiskPredictor_IP_Reputation(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_risk_predictor.%s", resourceName)

	name := resourceName

	fullCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "type", "IP_REPUTATION"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "true"),
		resource.TestCheckResourceAttr(resourceFullName, "default.result.level", "MEDIUM"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_ip_reputation.allowed_cidr_list.#", "3"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "predictor_ip_reputation.allowed_cidr_list.*", "192.168.0.0/24"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "predictor_ip_reputation.allowed_cidr_list.*", "10.0.0.0/8"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "predictor_ip_reputation.allowed_cidr_list.*", "172.16.0.0/12"),
	)

	minimalCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "type", "IP_REPUTATION"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "true"),
		resource.TestCheckNoResourceAttr(resourceFullName, "default.result.level"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_ip_reputation.allowed_cidr_list.#", "0"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             risk.RiskPredictor_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			{
				Config: testAccRiskPredictorConfig_IPReputation_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config:  testAccRiskPredictorConfig_IPReputation_Full(resourceName, name),
				Destroy: true,
			},
			// Minimal
			{
				Config: testAccRiskPredictorConfig_IPReputation_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config:  testAccRiskPredictorConfig_IPReputation_Minimal(resourceName, name),
				Destroy: true,
			},
			// Change
			{
				Config: testAccRiskPredictorConfig_IPReputation_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config: testAccRiskPredictorConfig_IPReputation_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config: testAccRiskPredictorConfig_IPReputation_Full(resourceName, name),
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

func TestAccRiskPredictor_IP_Reputation_OverwriteUndeletable(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_risk_predictor.%s", resourceName)

	name := resourceName
	compactName := "ipRisk"

	fullCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "name", name),
		resource.TestCheckResourceAttr(resourceFullName, "compact_name", "ipRisk"),
		resource.TestCheckResourceAttr(resourceFullName, "type", "IP_REPUTATION"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "false"),
		resource.TestCheckResourceAttr(resourceFullName, "default.result.level", "MEDIUM"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_ip_reputation.allowed_cidr_list.#", "3"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "predictor_ip_reputation.allowed_cidr_list.*", "192.168.0.0/24"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "predictor_ip_reputation.allowed_cidr_list.*", "10.0.0.0/8"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "predictor_ip_reputation.allowed_cidr_list.*", "172.16.0.0/12"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             risk.RiskPredictor_CheckDestroyUndeletable,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			{
				Config: testAccRiskPredictorConfig_IP_Reputation_OverwriteUndeletable(resourceName, name, compactName),
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

func TestAccRiskPredictor_CustomMap_BetweenRanges(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_risk_predictor.%s", resourceName)

	name := resourceName

	fullCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "type", "MAP"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "true"),
		resource.TestCheckResourceAttr(resourceFullName, "default.result.level", "MEDIUM"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_custom_map.contains", "${event.myshop}"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_custom_map.type", "RANGE"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_custom_map.between_ranges.high.max_value", "6"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_custom_map.between_ranges.high.min_value", "5"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_custom_map.between_ranges.medium.max_value", "4"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_custom_map.between_ranges.medium.min_value", "3"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_custom_map.between_ranges.low.max_value", "2"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_custom_map.between_ranges.low.min_value", "1"),
	)

	minimalCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "type", "MAP"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "true"),
		resource.TestCheckNoResourceAttr(resourceFullName, "default.result.level"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_custom_map.contains", "${event.myshop}"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_custom_map.type", "RANGE"),
		resource.TestCheckNoResourceAttr(resourceFullName, "predictor_custom_map.between_ranges.high.max_value"),
		resource.TestCheckNoResourceAttr(resourceFullName, "predictor_custom_map.between_ranges.high.min_value"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_custom_map.between_ranges.medium.max_value", "4"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_custom_map.between_ranges.medium.min_value", "3"),
		resource.TestCheckNoResourceAttr(resourceFullName, "predictor_custom_map.between_ranges.low.max_value"),
		resource.TestCheckNoResourceAttr(resourceFullName, "predictor_custom_map.between_ranges.low.min_value"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             risk.RiskPredictor_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			{
				Config: testAccRiskPredictorConfig_CustomMap_BetweenRanges_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config:  testAccRiskPredictorConfig_CustomMap_BetweenRanges_Full(resourceName, name),
				Destroy: true,
			},
			// Minimal
			{
				Config: testAccRiskPredictorConfig_CustomMap_BetweenRanges_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config:  testAccRiskPredictorConfig_CustomMap_BetweenRanges_Minimal(resourceName, name),
				Destroy: true,
			},
			// Change
			{
				Config: testAccRiskPredictorConfig_CustomMap_BetweenRanges_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config: testAccRiskPredictorConfig_CustomMap_BetweenRanges_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config: testAccRiskPredictorConfig_CustomMap_BetweenRanges_Full(resourceName, name),
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

func TestAccRiskPredictor_CustomMap_IPRanges(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_risk_predictor.%s", resourceName)

	name := resourceName

	fullCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "type", "MAP"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "true"),
		resource.TestCheckResourceAttr(resourceFullName, "default.result.level", "MEDIUM"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_custom_map.contains", "${event.myshop}"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_custom_map.type", "IP_RANGE"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "predictor_custom_map.ip_ranges.high.values.*", "192.168.0.0/24"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "predictor_custom_map.ip_ranges.high.values.*", "10.0.0.0/8"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "predictor_custom_map.ip_ranges.high.values.*", "172.16.0.0/12"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "predictor_custom_map.ip_ranges.medium.values.*", "192.0.2.0/24"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "predictor_custom_map.ip_ranges.medium.values.*", "192.168.1.0/26"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "predictor_custom_map.ip_ranges.medium.values.*", "10.10.0.0/16"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "predictor_custom_map.ip_ranges.low.values.*", "172.16.0.0/16"),
	)

	minimalCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "type", "MAP"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "true"),
		resource.TestCheckNoResourceAttr(resourceFullName, "default.result.level"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_custom_map.contains", "${event.myshop}"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_custom_map.type", "IP_RANGE"),
		resource.TestCheckNoResourceAttr(resourceFullName, "predictor_custom_map.ip_ranges.high.values"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "predictor_custom_map.ip_ranges.medium.values.*", "192.0.2.0/24"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "predictor_custom_map.ip_ranges.medium.values.*", "10.0.0.0/8"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "predictor_custom_map.ip_ranges.medium.values.*", "172.16.0.0/12"),
		resource.TestCheckNoResourceAttr(resourceFullName, "predictor_custom_map.ip_ranges.low.values"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             risk.RiskPredictor_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			{
				Config: testAccRiskPredictorConfig_CustomMap_IPRanges_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config:  testAccRiskPredictorConfig_CustomMap_IPRanges_Full(resourceName, name),
				Destroy: true,
			},
			// Minimal
			{
				Config: testAccRiskPredictorConfig_CustomMap_IPRanges_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config:  testAccRiskPredictorConfig_CustomMap_IPRanges_Minimal(resourceName, name),
				Destroy: true,
			},
			// Change
			{
				Config: testAccRiskPredictorConfig_CustomMap_IPRanges_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config: testAccRiskPredictorConfig_CustomMap_IPRanges_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config: testAccRiskPredictorConfig_CustomMap_IPRanges_Full(resourceName, name),
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

func TestAccRiskPredictor_CustomMap_StringList(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_risk_predictor.%s", resourceName)

	name := resourceName

	fullCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "type", "MAP"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "true"),
		resource.TestCheckResourceAttr(resourceFullName, "default.result.level", "MEDIUM"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_custom_map.contains", "${event.myshop}"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_custom_map.type", "STRING_LIST"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "predictor_custom_map.string_list.high.values.*", "HIGH"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "predictor_custom_map.string_list.high.values.*", "HIGH321"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "predictor_custom_map.string_list.high.values.*", "HIGH123"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "predictor_custom_map.string_list.medium.values.*", "MEDIUM"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "predictor_custom_map.string_list.medium.values.*", "MED321"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "predictor_custom_map.string_list.medium.values.*", "MED123"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "predictor_custom_map.string_list.low.values.*", "LOW"),
	)

	minimalCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "type", "MAP"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "true"),
		resource.TestCheckNoResourceAttr(resourceFullName, "default.result.level"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_custom_map.contains", "${event.myshop}"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_custom_map.type", "STRING_LIST"),
		resource.TestCheckNoResourceAttr(resourceFullName, "predictor_custom_map.string_list.high.values"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "predictor_custom_map.string_list.medium.values.*", "MEDIUM"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "predictor_custom_map.string_list.medium.values.*", "MED321"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "predictor_custom_map.string_list.medium.values.*", "MED123"),
		resource.TestCheckNoResourceAttr(resourceFullName, "predictor_custom_map.string_list.low.values"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             risk.RiskPredictor_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			{
				Config: testAccRiskPredictorConfig_CustomMap_StringList_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config:  testAccRiskPredictorConfig_CustomMap_StringList_Full(resourceName, name),
				Destroy: true,
			},
			// Minimal
			{
				Config: testAccRiskPredictorConfig_CustomMap_StringList_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config:  testAccRiskPredictorConfig_CustomMap_StringList_Minimal(resourceName, name),
				Destroy: true,
			},
			// Change
			{
				Config: testAccRiskPredictorConfig_CustomMap_StringList_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config: testAccRiskPredictorConfig_CustomMap_StringList_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config: testAccRiskPredictorConfig_CustomMap_StringList_Full(resourceName, name),
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

func TestAccRiskPredictor_NewDevice(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_risk_predictor.%s", resourceName)

	name := resourceName

	year, month, day := time.Now().Local().Date()
	activationAt := time.Date(year, month, day, 0, 0, 0, 0, time.UTC).AddDate(0, 0, -1).Format(time.RFC3339)

	fullCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "type", "DEVICE"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "true"),
		resource.TestCheckResourceAttr(resourceFullName, "default.result.level", "MEDIUM"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_device.detect", "NEW_DEVICE"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_device.activation_at", activationAt),
		resource.TestCheckNoResourceAttr(resourceFullName, "predictor_device.should_validate_payload_signature"),
	)

	minimalCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "type", "DEVICE"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "true"),
		resource.TestCheckNoResourceAttr(resourceFullName, "default.result.level"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_device.detect", "NEW_DEVICE"),
		resource.TestCheckNoResourceAttr(resourceFullName, "predictor_device.activation_at"),
		resource.TestCheckNoResourceAttr(resourceFullName, "predictor_device.should_validate_payload_signature"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             risk.RiskPredictor_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			{
				Config: testAccRiskPredictorConfig_NewDevice_Full(resourceName, name, activationAt),
				Check:  fullCheck,
			},
			{
				Config:  testAccRiskPredictorConfig_NewDevice_Full(resourceName, name, activationAt),
				Destroy: true,
			},
			// Minimal
			{
				Config: testAccRiskPredictorConfig_NewDevice_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config:  testAccRiskPredictorConfig_NewDevice_Minimal(resourceName, name),
				Destroy: true,
			},
			// Change
			{
				Config: testAccRiskPredictorConfig_NewDevice_Full(resourceName, name, activationAt),
				Check:  fullCheck,
			},
			{
				Config: testAccRiskPredictorConfig_NewDevice_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config: testAccRiskPredictorConfig_NewDevice_Full(resourceName, name, activationAt),
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

func TestAccRiskPredictor_NewDevice_OverwriteUndeletable(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_risk_predictor.%s", resourceName)

	name := resourceName
	compactName := "newDevice"

	year, month, day := time.Now().Local().Date()
	activationAt := time.Date(year, month, day, 0, 0, 0, 0, time.UTC).AddDate(0, 0, -1).Format(time.RFC3339)

	fullCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "name", name),
		resource.TestCheckResourceAttr(resourceFullName, "compact_name", "newDevice"),
		resource.TestCheckResourceAttr(resourceFullName, "type", "DEVICE"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "false"),
		resource.TestCheckResourceAttr(resourceFullName, "default.result.level", "MEDIUM"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_device.detect", "NEW_DEVICE"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_device.activation_at", activationAt),
		resource.TestCheckNoResourceAttr(resourceFullName, "predictor_device.should_validate_payload_signature"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             risk.RiskPredictor_CheckDestroyUndeletable,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			{
				Config: testAccRiskPredictorConfig_NewDevice_OverwriteUndeletable(resourceName, name, compactName, activationAt),
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

func TestAccRiskPredictor_Email_Reputation(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_risk_predictor.%s", resourceName)

	name := resourceName

	fullCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "type", "EMAIL_REPUTATION"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "true"),
		resource.TestCheckResourceAttr(resourceFullName, "default.result.level", "MEDIUM"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_email_reputation.#", "0"),
	)

	minimalCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "type", "EMAIL_REPUTATION"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "true"),
		resource.TestCheckNoResourceAttr(resourceFullName, "default.result.level"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_email_reputation.#", "0"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             risk.RiskPredictor_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			{
				Config: testAccRiskPredictorConfig_Email_Reputation_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config:  testAccRiskPredictorConfig_Email_Reputation_Full(resourceName, name),
				Destroy: true,
			},
			// Minimal
			{
				Config: testAccRiskPredictorConfig_Email_Reputation_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config:  testAccRiskPredictorConfig_Email_Reputation_Minimal(resourceName, name),
				Destroy: true,
			},
			// Change
			{
				Config: testAccRiskPredictorConfig_Email_Reputation_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config: testAccRiskPredictorConfig_Email_Reputation_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config: testAccRiskPredictorConfig_Email_Reputation_Full(resourceName, name),
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

func TestAccRiskPredictor_Email_Reputation_OverwriteUndeletable(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_risk_predictor.%s", resourceName)

	name := resourceName
	compactName := "emailReputation"

	fullCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "name", name),
		resource.TestCheckResourceAttr(resourceFullName, "compact_name", "emailReputation"),
		resource.TestCheckResourceAttr(resourceFullName, "type", "EMAIL_REPUTATION"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "false"),
		resource.TestCheckResourceAttr(resourceFullName, "default.result.level", "MEDIUM"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_email_reputation.#", "0"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             risk.RiskPredictor_CheckDestroyUndeletable,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			{
				Config: testAccRiskPredictorConfig_Email_Reputation_OverwriteUndeletable(resourceName, name, compactName),
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

func TestAccRiskPredictor_SuspiciousDevice(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_risk_predictor.%s", resourceName)

	name := resourceName

	fullCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "type", "DEVICE"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "true"),
		resource.TestCheckResourceAttr(resourceFullName, "default.result.level", "MEDIUM"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_device.detect", "SUSPICIOUS_DEVICE"),
		resource.TestCheckNoResourceAttr(resourceFullName, "predictor_device.activation_at"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_device.should_validate_payload_signature", "true"),
	)

	minimalCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "type", "DEVICE"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "true"),
		resource.TestCheckNoResourceAttr(resourceFullName, "default.result.level"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_device.detect", "SUSPICIOUS_DEVICE"),
		resource.TestCheckNoResourceAttr(resourceFullName, "predictor_device.activation_at"),
		resource.TestCheckNoResourceAttr(resourceFullName, "predictor_device.should_validate_payload_signature"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             risk.RiskPredictor_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			{
				Config: testAccRiskPredictorConfig_SuspiciousDevice_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config:  testAccRiskPredictorConfig_SuspiciousDevice_Full(resourceName, name),
				Destroy: true,
			},
			// Minimal
			{
				Config: testAccRiskPredictorConfig_SuspiciousDevice_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config:  testAccRiskPredictorConfig_SuspiciousDevice_Minimal(resourceName, name),
				Destroy: true,
			},
			// Change
			{
				Config: testAccRiskPredictorConfig_SuspiciousDevice_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config: testAccRiskPredictorConfig_SuspiciousDevice_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config: testAccRiskPredictorConfig_SuspiciousDevice_Full(resourceName, name),
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

func TestAccRiskPredictor_SuspiciousDevice_OverwriteUndeletable(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_risk_predictor.%s", resourceName)

	name := resourceName
	compactName := "suspiciousDevice"

	fullCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "name", name),
		resource.TestCheckResourceAttr(resourceFullName, "compact_name", "suspiciousDevice"),
		resource.TestCheckResourceAttr(resourceFullName, "type", "DEVICE"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "false"),
		resource.TestCheckResourceAttr(resourceFullName, "default.result.level", "MEDIUM"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_device.detect", "SUSPICIOUS_DEVICE"),
		resource.TestCheckNoResourceAttr(resourceFullName, "predictor_device.activation_at"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             risk.RiskPredictor_CheckDestroyUndeletable,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			{
				Config: testAccRiskPredictorConfig_SuspiciousDevice_OverwriteUndeletable(resourceName, name, compactName),
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

func TestAccRiskPredictor_TrafficAnomaly(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_risk_predictor.%s", resourceName)

	name := resourceName
	compactName := "trafficAnomalyCompactName"

	initialCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "type", "TRAFFIC_ANOMALY"),
		resource.TestCheckResourceAttr(resourceFullName, "compact_name", compactName),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "true"),
		resource.TestCheckResourceAttr(resourceFullName, "default.result.level", "MEDIUM"),
		resource.TestCheckResourceAttr(resourceFullName, "name", name),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_traffic_anomaly.rules.0.enabled", "true"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_traffic_anomaly.rules.0.type", "UNIQUE_USERS_PER_DEVICE"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_traffic_anomaly.rules.0.interval.quantity", "1"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_traffic_anomaly.rules.0.interval.unit", "DAY"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_traffic_anomaly.rules.0.threshold.high", "4"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_traffic_anomaly.rules.0.threshold.medium", "3"),
	)
	changeCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "type", "TRAFFIC_ANOMALY"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "true"),
		resource.TestCheckResourceAttr(resourceFullName, "default.result.level", "HIGH"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_traffic_anomaly.rules.0.enabled", "false"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_traffic_anomaly.rules.0.type", "UNIQUE_USERS_PER_DEVICE"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_traffic_anomaly.rules.0.interval.quantity", "2"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_traffic_anomaly.rules.0.interval.unit", "HOUR"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_traffic_anomaly.rules.0.threshold.high", "5"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_traffic_anomaly.rules.0.threshold.medium", "4"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             risk.RiskPredictor_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Initial (all properties in type are required)
			{
				Config: testAccRiskPredictorConfig_TrafficAnomaly_Initial(resourceName, name, compactName),
				Check:  initialCheck,
			},
			{
				Config:  testAccRiskPredictorConfig_TrafficAnomaly_Initial(resourceName, name, compactName),
				Destroy: true,
			},
			// Change
			{
				Config: testAccRiskPredictorConfig_TrafficAnomaly_Initial(resourceName, name, compactName),
				Check:  initialCheck,
			},
			{
				Config: testAccRiskPredictorConfig_TrafficAnomaly_Change(resourceName, name, compactName),
				Check:  changeCheck,
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

func TestAccRiskPredictor_TrafficAnomaly_OverwriteUndeletable(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_risk_predictor.%s", resourceName)

	name := resourceName
	compactName := "trafficAnomaly"

	fullCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "type", "TRAFFIC_ANOMALY"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "false"),
		resource.TestCheckResourceAttr(resourceFullName, "default.result.level", "HIGH"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_traffic_anomaly.rules.0.enabled", "false"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_traffic_anomaly.rules.0.type", "UNIQUE_USERS_PER_DEVICE"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_traffic_anomaly.rules.0.interval.quantity", "2"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_traffic_anomaly.rules.0.interval.unit", "HOUR"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_traffic_anomaly.rules.0.threshold.high", "5"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_traffic_anomaly.rules.0.threshold.medium", "4"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             risk.RiskPredictor_CheckDestroyUndeletable,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			{
				Config: testAccRiskPredictorConfig_TrafficAnomaly_OverwriteUndeletable(resourceName, name, compactName),
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

func TestAccRiskPredictor_UserLocationAnomaly(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_risk_predictor.%s", resourceName)

	name := resourceName

	fullCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "type", "USER_LOCATION_ANOMALY"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "true"),
		resource.TestCheckResourceAttr(resourceFullName, "default.result.level", "MEDIUM"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_user_location_anomaly.radius.distance", "100"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_user_location_anomaly.radius.unit", "miles"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_user_location_anomaly.days", "90"),
	)

	minimalCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "type", "USER_LOCATION_ANOMALY"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "true"),
		resource.TestCheckNoResourceAttr(resourceFullName, "default.result.level"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_user_location_anomaly.radius.distance", "51"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_user_location_anomaly.radius.unit", "kilometers"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_user_location_anomaly.days", "90"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             risk.RiskPredictor_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			{
				Config: testAccRiskPredictorConfig_UserLocationAnomaly_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config:  testAccRiskPredictorConfig_UserLocationAnomaly_Full(resourceName, name),
				Destroy: true,
			},
			// Minimal
			{
				Config: testAccRiskPredictorConfig_UserLocationAnomaly_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config:  testAccRiskPredictorConfig_UserLocationAnomaly_Minimal(resourceName, name),
				Destroy: true,
			},
			// Change
			{
				Config: testAccRiskPredictorConfig_UserLocationAnomaly_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config: testAccRiskPredictorConfig_UserLocationAnomaly_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config: testAccRiskPredictorConfig_UserLocationAnomaly_Full(resourceName, name),
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

func TestAccRiskPredictor_UserLocationAnomaly_OverwriteUndeletable(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_risk_predictor.%s", resourceName)

	name := resourceName
	compactName := "userLocationAnomaly"

	fullCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "name", name),
		resource.TestCheckResourceAttr(resourceFullName, "compact_name", "userLocationAnomaly"),
		resource.TestCheckResourceAttr(resourceFullName, "type", "USER_LOCATION_ANOMALY"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "false"),
		resource.TestCheckResourceAttr(resourceFullName, "default.result.level", "MEDIUM"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_user_location_anomaly.radius.distance", "100"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_user_location_anomaly.radius.unit", "miles"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_user_location_anomaly.days", "90"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             risk.RiskPredictor_CheckDestroyUndeletable,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			{
				Config: testAccRiskPredictorConfig_UserLocationAnomaly_OverwriteUndeletable(resourceName, name, compactName),
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

func TestAccRiskPredictor_Velocity(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_risk_predictor.%s", resourceName)

	name := resourceName

	byUserCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "type", "VELOCITY"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "true"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_velocity.of", "${event.ip}"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_velocity.by.#", "1"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "predictor_velocity.by.*", "${event.user.id}"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_velocity.measure", "DISTINCT_COUNT"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_velocity.use.type", "POISSON_WITH_MAX"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_velocity.use.medium", "2"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_velocity.use.high", "4"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_velocity.fallback.strategy", "ENVIRONMENT_MAX"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_velocity.fallback.high", "30"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_velocity.fallback.medium", "20"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_velocity.every.unit", "HOUR"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_velocity.every.quantity", "1"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_velocity.every.min_sample", "5"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_velocity.sliding_window.unit", "DAY"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_velocity.sliding_window.quantity", "7"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_velocity.sliding_window.min_sample", "3"),
	)

	byIPCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "type", "VELOCITY"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "true"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_velocity.of", "${event.user.id}"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_velocity.by.#", "1"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "predictor_velocity.by.*", "${event.ip}"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_velocity.measure", "DISTINCT_COUNT"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_velocity.use.type", "POISSON_WITH_MAX"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_velocity.use.medium", "2"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_velocity.use.high", "4"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_velocity.fallback.strategy", "ENVIRONMENT_MAX"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_velocity.fallback.high", "3500"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_velocity.fallback.medium", "2500"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_velocity.every.unit", "HOUR"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_velocity.every.quantity", "1"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_velocity.every.min_sample", "5"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_velocity.sliding_window.unit", "DAY"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_velocity.sliding_window.quantity", "7"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_velocity.sliding_window.min_sample", "3"),
	)

	resource.Test(t, resource.TestCase{
		// PreCheck: func() {
		//	acctest.PreCheckNoTestAccFlaky(t)
		// 	acctest.PreCheckClient(t)
		// 	acctest.PreCheckNoBeta(t)
		//	acctest.PreCheckNoBeta(t)
		// },
		PreCheck:                 func() { t.Skipf("STAGING-21856") },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             risk.RiskPredictor_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// By User
			{
				Config: testAccRiskPredictorConfig_Velocity_ByUser_Full(resourceName, name),
				Check:  byUserCheck,
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
				Config:  testAccRiskPredictorConfig_Velocity_ByUser_Full(resourceName, name),
				Destroy: true,
			},
			// By IP
			{
				Config: testAccRiskPredictorConfig_Velocity_ByIP_Full(resourceName, name),
				Check:  byIPCheck,
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
				Config:  testAccRiskPredictorConfig_Velocity_ByIP_Full(resourceName, name),
				Destroy: true,
			},
			// Change
			{
				Config: testAccRiskPredictorConfig_Velocity_ByUser_Full(resourceName, name),
				Check:  byUserCheck,
			},
			{
				Config: testAccRiskPredictorConfig_Velocity_ByIP_Full(resourceName, name),
				Check:  byIPCheck,
			},
			{
				Config: testAccRiskPredictorConfig_Velocity_ByUser_Full(resourceName, name),
				Check:  byUserCheck,
			},
		},
	})
}

func TestAccRiskPredictor_Velocity_OverwriteUndeletable(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_risk_predictor.%s", resourceName)

	name := resourceName
	compactNameByUser := "ipVelocityByUser"
	compactNameByIP := "userVelocityByIp"

	byUserCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "name", name),
		resource.TestCheckResourceAttr(resourceFullName, "compact_name", compactNameByUser),
		resource.TestCheckResourceAttr(resourceFullName, "type", "VELOCITY"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "false"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_velocity.of", "${event.ip}"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_velocity.by.#", "1"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "predictor_velocity.by.*", "${event.user.id}"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_velocity.measure", "DISTINCT_COUNT"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_velocity.use.type", "POISSON_WITH_MAX"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_velocity.use.medium", "2"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_velocity.use.high", "4"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_velocity.fallback.strategy", "ENVIRONMENT_MAX"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_velocity.fallback.high", "30"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_velocity.fallback.medium", "20"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_velocity.every.unit", "HOUR"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_velocity.every.quantity", "1"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_velocity.every.min_sample", "5"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_velocity.sliding_window.unit", "DAY"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_velocity.sliding_window.quantity", "7"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_velocity.sliding_window.min_sample", "3"),
	)

	byIPCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "name", name),
		resource.TestCheckResourceAttr(resourceFullName, "compact_name", compactNameByIP),
		resource.TestCheckResourceAttr(resourceFullName, "type", "VELOCITY"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_velocity.of", "${event.user.id}"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_velocity.by.#", "1"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "predictor_velocity.by.*", "${event.ip}"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_velocity.measure", "DISTINCT_COUNT"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_velocity.use.type", "POISSON_WITH_MAX"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_velocity.use.medium", "2"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_velocity.use.high", "4"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_velocity.fallback.strategy", "ENVIRONMENT_MAX"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_velocity.fallback.high", "3500"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_velocity.fallback.medium", "2500"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_velocity.every.unit", "HOUR"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_velocity.every.quantity", "1"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_velocity.every.min_sample", "5"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_velocity.sliding_window.unit", "DAY"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_velocity.sliding_window.quantity", "7"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_velocity.sliding_window.min_sample", "3"),
	)

	resource.Test(t, resource.TestCase{
		// PreCheck: func() {
		//	acctest.PreCheckNoTestAccFlaky(t)
		// 	acctest.PreCheckClient(t)
		// 	acctest.PreCheckNoBeta(t)
		//	acctest.PreCheckNoBeta(t)
		// },
		PreCheck:                 func() { t.Skipf("STAGING-21856") },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             risk.RiskPredictor_CheckDestroyUndeletable,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// By User
			{
				Config: testAccRiskPredictorConfig_Velocity_ByUser_Full_OverwriteUndeletable(resourceName, name, compactNameByUser),
				Check:  byUserCheck,
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
				Config:  testAccRiskPredictorConfig_Velocity_ByUser_Full_OverwriteUndeletable(resourceName, name, compactNameByUser),
				Destroy: true,
			},
			// By IP
			{
				Config: testAccRiskPredictorConfig_Velocity_ByIP_Full_OverwriteUndeletable(resourceName, name, compactNameByIP),
				Check:  byIPCheck,
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
				Config:  testAccRiskPredictorConfig_Velocity_ByIP_Full_OverwriteUndeletable(resourceName, name, compactNameByIP),
				Destroy: true,
			},
			// Change
			{
				Config: testAccRiskPredictorConfig_Velocity_ByUser_Full_OverwriteUndeletable(resourceName, name, compactNameByUser),
				Check:  byUserCheck,
			},
			{
				Config: testAccRiskPredictorConfig_Velocity_ByIP_Full_OverwriteUndeletable(resourceName, name, compactNameByIP),
				Check:  byIPCheck,
			},
			{
				Config: testAccRiskPredictorConfig_Velocity_ByUser_Full_OverwriteUndeletable(resourceName, name, compactNameByUser),
				Check:  byUserCheck,
			},
		},
	})
}

func TestAccRiskPredictor_UserRiskBehavior(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_risk_predictor.%s", resourceName)

	name := resourceName

	byUserCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "type", "USER_RISK_BEHAVIOR"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "true"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_user_risk_behavior.prediction_model.name", "points"),
	)

	byOrgCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "type", "USER_RISK_BEHAVIOR"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "true"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_user_risk_behavior.prediction_model.name", "login_anomaly_statistic"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckTestAccFlaky(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             risk.RiskPredictor_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// By User
			{
				Config: testAccRiskPredictorConfig_UserRiskBehavior_ByUser_Full(resourceName, name),
				Check:  byUserCheck,
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
				Config:  testAccRiskPredictorConfig_UserRiskBehavior_ByUser_Full(resourceName, name),
				Destroy: true,
			},
			// By Org
			{
				Config: testAccRiskPredictorConfig_UserRiskBehavior_ByOrg_Full(resourceName, name),
				Check:  byOrgCheck,
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
				Config:  testAccRiskPredictorConfig_UserRiskBehavior_ByOrg_Full(resourceName, name),
				Destroy: true,
			},
			// Change
			{
				Config: testAccRiskPredictorConfig_UserRiskBehavior_ByUser_Full(resourceName, name),
				Check:  byUserCheck,
			},
			{
				Config: testAccRiskPredictorConfig_UserRiskBehavior_ByOrg_Full(resourceName, name),
				Check:  byOrgCheck,
			},
			{
				Config: testAccRiskPredictorConfig_UserRiskBehavior_ByUser_Full(resourceName, name),
				Check:  byUserCheck,
			},
		},
	})
}

func TestAccRiskPredictor_UserRiskBehavior_OverwriteUndeletable(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_risk_predictor.%s", resourceName)

	name := resourceName
	compactNameByUser := "userBasedRiskBehavior"
	compactNameByOrg := "userRiskBehavior"

	byUserCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "name", name),
		resource.TestCheckResourceAttr(resourceFullName, "compact_name", compactNameByUser),
		resource.TestCheckResourceAttr(resourceFullName, "type", "USER_RISK_BEHAVIOR"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "false"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_user_risk_behavior.prediction_model.name", "points"),
	)

	byOrgCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "name", name),
		resource.TestCheckResourceAttr(resourceFullName, "compact_name", compactNameByOrg),
		resource.TestCheckResourceAttr(resourceFullName, "type", "USER_RISK_BEHAVIOR"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "false"),
		resource.TestCheckResourceAttr(resourceFullName, "predictor_user_risk_behavior.prediction_model.name", "login_anomaly_statistic"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckTestAccFlaky(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             risk.RiskPredictor_CheckDestroyUndeletable,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// By User
			{
				Config: testAccRiskPredictorConfig_UserRiskBehavior_ByUser_Full_OverwriteUndeletable(resourceName, name, compactNameByUser),
				Check:  byUserCheck,
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
				Config:  testAccRiskPredictorConfig_UserRiskBehavior_ByUser_Full_OverwriteUndeletable(resourceName, name, compactNameByUser),
				Destroy: true,
			},
			// By Org
			{
				Config: testAccRiskPredictorConfig_UserRiskBehavior_ByOrg_Full_OverwriteUndeletable(resourceName, name, compactNameByOrg),
				Check:  byOrgCheck,
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
				Config:  testAccRiskPredictorConfig_UserRiskBehavior_ByOrg_Full_OverwriteUndeletable(resourceName, name, compactNameByOrg),
				Destroy: true,
			},
			// Change
			{
				Config: testAccRiskPredictorConfig_UserRiskBehavior_ByUser_Full_OverwriteUndeletable(resourceName, name, compactNameByUser),
				Check:  byUserCheck,
			},
			{
				Config: testAccRiskPredictorConfig_UserRiskBehavior_ByOrg_Full_OverwriteUndeletable(resourceName, name, compactNameByOrg),
				Check:  byOrgCheck,
			},
			{
				Config: testAccRiskPredictorConfig_UserRiskBehavior_ByUser_Full_OverwriteUndeletable(resourceName, name, compactNameByUser),
				Check:  byUserCheck,
			},
		},
	})
}

func TestAccRiskPredictor_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_risk_predictor.%s", resourceName)

	name := resourceName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoBeta(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             risk.RiskPredictor_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccRiskPredictorConfig_Minimal(resourceName, name),
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

func testAccRiskPredictorConfig_NewEnv(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_risk_predictor" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  name         = "%[4]s"
  compact_name = "%[4]s1"

  predictor_anonymous_network = {}

}`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccRiskPredictorConfig_Full(resourceName, name string) string {
	return testAccRiskPredictorConfig_Anonymous_Network_Full(resourceName, name)
}

func testAccRiskPredictorConfig_Minimal(resourceName, name string) string {
	return testAccRiskPredictorConfig_Anonymous_Network_Minimal(resourceName, name)
}

func testAccRiskPredictorConfig_Adversary_In_The_Middle_Full(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[3]s1"
  description  = "When my wife is upset, I let her colour in my black and white tattoos.  She just needs a shoulder to crayon.."

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  predictor_adversary_in_the_middle = {
    allowed_domain_list = [
      "domain2.com",
      "domain1.com",
      "domain3.com",
    ]
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPredictorConfig_Adversary_In_The_Middle_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[3]s1"

  predictor_adversary_in_the_middle = {}
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPredictorConfig_Adversary_In_The_Middle_OverwriteUndeletable(resourceName, name, compactName string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[4]s"

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  predictor_adversary_in_the_middle = {
    allowed_domain_list = [
      "domain2.com",
      "domain1.com",
      "domain3.com",
    ]
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name, compactName)
}

func testAccRiskPredictorConfig_Anonymous_Network_Full(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[3]s1"
  description  = "When my wife is upset, I let her colour in my black and white tattoos.  She just needs a shoulder to crayon.."

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  predictor_anonymous_network = {
    allowed_cidr_list = [
      "10.0.0.0/8",
      "172.16.0.0/12",
      "192.168.0.0/24"
    ]
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPredictorConfig_Anonymous_Network_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[3]s1"

  predictor_anonymous_network = {}
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPredictorConfig_Anonymous_Network_OverwriteUndeletable(resourceName, name, compactName string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[4]s"

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  predictor_anonymous_network = {
    allowed_cidr_list = [
      "10.0.0.0/8",
      "172.16.0.0/12",
      "192.168.0.0/24"
    ]
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name, compactName)
}

func testAccRiskPredictorConfig_Bot_Detection_Full(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[3]s1"
  description  = "The neighbours said their dog will retrieve sticks from 10 miles away.  Sounds far fetched to me."

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  predictor_bot_detection = {
    include_repeated_events_without_sdk = true
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPredictorConfig_Bot_Detection_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[3]s1"

  predictor_bot_detection = {}
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPredictorConfig_Bot_Detection_OverwriteUndeletable(resourceName, name, compactName string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[4]s"

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  predictor_bot_detection = {
    include_repeated_events_without_sdk = true
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name, compactName)
}

func testAccRiskPredictorConfig_Composite_Full_1(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[3]s1"
  description  = "When my wife is upset, I let her colour in my black and white tattoos.  She just needs a shoulder to crayon.."

  predictor_composite = {
    compositions = [
      {
        level = "HIGH"

        condition_json = jsonencode({
          "not" : {
            "or" : [{
              "equals" : 0,
              "value" : "$${details.counters.predictorLevels.medium}",
              "type" : "VALUE_COMPARISON"
              }, {
              "equals" : "High",
              "value" : "$${details.geoVelocity.level}",
              "type" : "VALUE_COMPARISON"
              }, {
              "startsWith" : "admin",
              "value" : "$${event.user.name}",
              "type" : "VALUE_COMPARISON"
              }, {
              "endsWith" : "@example.com",
              "value" : "$${event.user.name}",
              "type" : "VALUE_COMPARISON"
              }, {
              "and" : [{
                "equals" : "High",
                "value" : "$${details.anonymousNetwork.level}",
                "type" : "VALUE_COMPARISON"
                }, {
                "list" : ["Group Name"],
                "contains" : "$${event.user.groups}",
                "type" : "GROUPS_INTERSECTION"
              }],
              "type" : "AND"
            }],
            "type" : "OR"
          },
          "type" : "NOT"
        })
      },
      {
        level = "LOW"

        condition_json = jsonencode({
          "and" : [
            {
              "value" : "$${details.counters.predictorLevels.medium}",
              "equals" : 5,
              "type" : "VALUE_COMPARISON"
            },
            {
              "value" : "$${details.anonymousNetwork.level}",
              "equals" : "low",
              "type" : "VALUE_COMPARISON"
            },
            {
              "startsWith" : "test",
              "value" : "$${event.user.name}",
              "type" : "VALUE_COMPARISON"
            },
            {
              "and" : [
                {
                  "value" : "$${details.anonymousNetwork.level}",
                  "equals" : "high",
                  "type" : "VALUE_COMPARISON"
                },
                {
                  "or" : [
                    {
                      "value" : "$${details.anonymousNetwork.level}",
                      "notEquals" : "high",
                      "type" : "VALUE_COMPARISON"
                    }
                  ]
                }
              ]
            }
          ]
        })
      }
    ]
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPredictorConfig_Composite_Full_2(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[3]s1"

  predictor_composite = {
    compositions = [
      {
        level = "LOW"

        condition_json = jsonencode({
          "and" : [
            {
              "value" : "$${details.counters.predictorLevels.medium}",
              "equals" : 5,
              "type" : "VALUE_COMPARISON"
            },
            {
              "value" : "$${details.anonymousNetwork.level}",
              "equals" : "low",
              "type" : "VALUE_COMPARISON"
            },
            {
              "endsWith" : "@test.org",
              "value" : "$${event.user.name}",
              "type" : "VALUE_COMPARISON"
            },
            {
              "and" : [
                {
                  "value" : "$${details.anonymousNetwork.level}",
                  "equals" : "high",
                  "type" : "VALUE_COMPARISON"
                },
                {
                  "or" : [
                    {
                      "value" : "$${details.anonymousNetwork.level}",
                      "notEquals" : "high",
                      "type" : "VALUE_COMPARISON"
                    }
                  ]
                }
              ]
            }
          ]
        })
      },
      {
        level = "HIGH"

        condition_json = jsonencode({
          "not" : {
            "or" : [{
              "equals" : 0,
              "value" : "$${details.counters.predictorLevels.medium}",
              "type" : "VALUE_COMPARISON"
              }, {
              "equals" : "High",
              "value" : "$${details.geoVelocity.level}",
              "type" : "VALUE_COMPARISON"
              }, {
              "startsWith" : "user",
              "value" : "$${event.user.name}",
              "type" : "VALUE_COMPARISON"
              }, {
              "and" : [{
                "equals" : "High",
                "value" : "$${details.anonymousNetwork.level}",
                "type" : "VALUE_COMPARISON"
              }],
              "type" : "AND"
            }],
            "type" : "OR"
          },
          "type" : "NOT"
        })
      },
    ]
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPredictorConfig_Composite_InvalidJSON(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[3]s1"

  predictor_composite = {
    compositions = [
      {
        level = "LOW"

        condition_json = jsonencode({})
      }
    ]
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPredictorConfig_Geovelocity_Full(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[3]s1"
  description  = "When my wife is upset, I let her colour in my black and white tattoos.  She just needs a shoulder to crayon.."

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  predictor_geovelocity = {

    allowed_cidr_list = [
      "10.0.0.0/8",
      "172.16.0.0/12",
      "192.168.0.0/24"
    ]
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPredictorConfig_Geovelocity_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[3]s1"

  predictor_geovelocity = {}

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPredictorConfig_Geovelocity_OverwriteUndeletable(resourceName, name, compactName string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[4]s"

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  predictor_geovelocity = {
    allowed_cidr_list = [
      "10.0.0.0/8",
      "172.16.0.0/12",
      "192.168.0.0/24"
    ]
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name, compactName)
}

func testAccRiskPredictorConfig_IPReputation_Full(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[3]s1"
  description  = "When my wife is upset, I let her colour in my black and white tattoos.  She just needs a shoulder to crayon.."

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  predictor_ip_reputation = {
    allowed_cidr_list = [
      "10.0.0.0/8",
      "172.16.0.0/12",
      "192.168.0.0/24"
    ]
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPredictorConfig_IPReputation_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[3]s1"

  predictor_ip_reputation = {}

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPredictorConfig_IP_Reputation_OverwriteUndeletable(resourceName, name, compactName string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[4]s"

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  predictor_ip_reputation = {
    allowed_cidr_list = [
      "10.0.0.0/8",
      "172.16.0.0/12",
      "192.168.0.0/24"
    ]
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name, compactName)
}

func testAccRiskPredictorConfig_CustomMap_BetweenRanges_Full(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[3]s1"
  description  = "When my wife is upset, I let her colour in my black and white tattoos.  She just needs a shoulder to crayon.."

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  predictor_custom_map = {
    contains = "$${event.myshop}"

    between_ranges = {
      high = {
        max_value = 6
        min_value = 5
      }

      medium = {
        max_value = 4
        min_value = 3
      }

      low = {
        max_value = 2
        min_value = 1
      }
    }
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPredictorConfig_CustomMap_BetweenRanges_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[3]s1"

  predictor_custom_map = {
    contains = "$${event.myshop}"

    between_ranges = {
      medium = {
        max_value = 4
        min_value = 3
      }
    }
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPredictorConfig_CustomMap_IPRanges_Full(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[3]s1"
  description  = "When my wife is upset, I let her colour in my black and white tattoos.  She just needs a shoulder to crayon.."

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  predictor_custom_map = {
    contains = "$${event.myshop}"

    ip_ranges = {
      high = {
        values = [
          "192.168.0.0/24",
          "10.0.0.0/8",
          "172.16.0.0/12"
        ]
      }

      medium = {
        values = [
          "192.0.2.0/24",
          "192.168.1.0/26",
          "10.10.0.0/16"
        ]
      }

      low = {
        values = [
          "172.16.0.0/16"
        ]
      }
    }
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPredictorConfig_CustomMap_IPRanges_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[3]s1"

  predictor_custom_map = {
    contains = "$${event.myshop}"

    ip_ranges = {
      medium = {
        values = [
          "192.0.2.0/24",
          "10.0.0.0/8",
          "172.16.0.0/12"
        ]
      }
    }
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPredictorConfig_CustomMap_StringList_Full(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[3]s1"
  description  = "When my wife is upset, I let her colour in my black and white tattoos.  She just needs a shoulder to crayon.."

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  predictor_custom_map = {
    contains = "$${event.myshop}"

    string_list = {
      high = {
        values = [
          "HIGH",
          "HIGH321",
          "HIGH123"
        ]
      }

      medium = {
        values = [
          "MEDIUM",
          "MED321",
          "MED123"
        ]
      }

      low = {
        values = [
          "LOW"
        ]
      }
    }
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPredictorConfig_CustomMap_StringList_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[3]s1"

  predictor_custom_map = {
    contains = "$${event.myshop}"

    string_list = {
      medium = {
        values = [
          "MEDIUM",
          "MED321",
          "MED123"
        ]
      }
    }
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPredictorConfig_NewDevice_Full(resourceName, name, activationAt string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[3]s1"
  description  = "When my wife is upset, I let her colour in my black and white tattoos.  She just needs a shoulder to crayon.."

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  predictor_device = {
    detect        = "NEW_DEVICE"
    activation_at = "%[4]s"
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name, activationAt)
}

func testAccRiskPredictorConfig_NewDevice_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[3]s1"

  predictor_device = {}
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPredictorConfig_NewDevice_OverwriteUndeletable(resourceName, name, compactName, activationAt string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[4]s"

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  predictor_device = {
    detect        = "NEW_DEVICE"
    activation_at = "%[5]s"
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name, compactName, activationAt)
}

func testAccRiskPredictorConfig_Email_Reputation_Full(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[3]s1"
  description  = "The neighbours said their dog will retrieve sticks from 10 miles away.  Sounds far fetched to me."

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  predictor_email_reputation = {}

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPredictorConfig_Email_Reputation_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[3]s1"

  predictor_email_reputation = {}
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPredictorConfig_Email_Reputation_OverwriteUndeletable(resourceName, name, compactName string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[4]s"

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  predictor_email_reputation = {}

}`, acctest.GenericSandboxEnvironment(), resourceName, name, compactName)
}

func testAccRiskPredictorConfig_SuspiciousDevice_Full(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[3]s1"
  description  = "The neighbours said their dog will retrieve sticks from 10 miles away.  Sounds far fetched to me."

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  predictor_device = {
    detect                            = "SUSPICIOUS_DEVICE"
    should_validate_payload_signature = true
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPredictorConfig_SuspiciousDevice_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[3]s1"

  predictor_device = {
    detect = "SUSPICIOUS_DEVICE"
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPredictorConfig_SuspiciousDevice_OverwriteUndeletable(resourceName, name, compactName string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[4]s"

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  predictor_device = {
    detect = "SUSPICIOUS_DEVICE"
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name, compactName)
}

func testAccRiskPredictorConfig_TrafficAnomaly_Initial(resourceName, name, compactName string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[4]s"

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  predictor_traffic_anomaly = {
    rules = [
      {
        type = "UNIQUE_USERS_PER_DEVICE"
        threshold = {
          medium = 3
          high   = 4
        }
        interval = {
          unit     = "DAY"
          quantity = 1
        }
        enabled = true
      }
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name, compactName)
}

func testAccRiskPredictorConfig_TrafficAnomaly_Change(resourceName, name, compactName string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[4]s"

  default = {
    result = {
      level = "HIGH"
    }
  }

  predictor_traffic_anomaly = {
    rules = [
      {
        type = "UNIQUE_USERS_PER_DEVICE"
        threshold = {
          medium = 4
          high   = 5
        }
        interval = {
          unit     = "HOUR"
          quantity = 2
        }
        enabled = false
      }
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name, compactName)
}

func testAccRiskPredictorConfig_TrafficAnomaly_OverwriteUndeletable(resourceName, name, compactName string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[4]s"

  default = {
    result = {
      level = "HIGH"
    }
  }

  predictor_traffic_anomaly = {
    rules = [
      {
        type = "UNIQUE_USERS_PER_DEVICE"
        threshold = {
          medium = 4
          high   = 5
        }
        interval = {
          unit     = "HOUR"
          quantity = 2
        }
        enabled = false
      }
    ]
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name, compactName)
}

func testAccRiskPredictorConfig_UserLocationAnomaly_Full(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[3]s1"
  description  = "When my wife is upset, I let her colour in my black and white tattoos.  She just needs a shoulder to crayon.."

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  predictor_user_location_anomaly = {
    radius = {
      distance = 100
      unit     = "miles"
    }
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPredictorConfig_UserLocationAnomaly_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[3]s1"

  predictor_user_location_anomaly = {
    radius = {
      distance = 51
    }
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPredictorConfig_UserLocationAnomaly_OverwriteUndeletable(resourceName, name, compactName string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[4]s"

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  predictor_user_location_anomaly = {
    radius = {
      distance = 100
      unit     = "miles"
    }
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name, compactName)
}

func testAccRiskPredictorConfig_Velocity_ByUser_Full(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[3]s1"
  description  = "When my wife is upset, I let her colour in my black and white tattoos.  She just needs a shoulder to crayon.."

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  predictor_velocity = {
    of = "$${event.ip}"
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPredictorConfig_Velocity_ByIP_Full(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[3]s1"

  predictor_velocity = {
    of = "$${event.user.id}"
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPredictorConfig_Velocity_ByUser_Full_OverwriteUndeletable(resourceName, name, compactName string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[4]s"

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  predictor_velocity = {
    of = "$${event.ip}"
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name, compactName)
}

func testAccRiskPredictorConfig_Velocity_ByIP_Full_OverwriteUndeletable(resourceName, name, compactName string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[4]s"

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  predictor_velocity = {
    of = "$${event.user.id}"
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name, compactName)
}

func testAccRiskPredictorConfig_UserRiskBehavior_ByUser_Full(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[3]s1"
  description  = "When my wife is upset, I let her colour in my black and white tattoos.  She just needs a shoulder to crayon.."

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  predictor_user_risk_behavior = {
    prediction_model = {
      name = "points"
    }
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPredictorConfig_UserRiskBehavior_ByOrg_Full(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[3]s1"

  predictor_user_risk_behavior = {
    prediction_model = {
      name = "login_anomaly_statistic"
    }
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPredictorConfig_UserRiskBehavior_ByUser_Full_OverwriteUndeletable(resourceName, name, compactName string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[4]s"

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  predictor_user_risk_behavior = {
    prediction_model = {
      name = "points"
    }
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name, compactName)
}

func testAccRiskPredictorConfig_UserRiskBehavior_ByOrg_Full_OverwriteUndeletable(resourceName, name, compactName string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[4]s"

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  predictor_user_risk_behavior = {
    prediction_model = {
      name = "login_anomaly_statistic"
    }
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name, compactName)
}
