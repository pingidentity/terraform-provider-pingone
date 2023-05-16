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

func testAccCheckRiskPredictorDestroy(s *terraform.State) error {
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
		if rs.Type != "pingone_risk_predictor" {
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

		body, r, err := apiClient.RiskAdvancedPredictorsApi.ReadOneRiskPredictor(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.ID).Execute()

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

		return fmt.Errorf("PingOne risk predictor %s still exists", rs.Primary.ID)
	}

	return nil
}

func testAccCheckRiskPredictorDestroyUndeletable(s *terraform.State) error {
	return nil
}

func TestAccRiskPredictor_NewEnv(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_risk_predictor.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckRiskPredictorDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRiskPredictorConfig_NewEnv(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
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
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
		resource.TestCheckResourceAttr(resourceFullName, "name", name),
		resource.TestCheckResourceAttr(resourceFullName, "compact_name", fmt.Sprintf("%s1", name)),
		resource.TestCheckResourceAttr(resourceFullName, "description", "When my wife is upset, I let her colour in my black and white tattoos.  She just needs a shoulder to crayon.."),
		resource.TestCheckResourceAttr(resourceFullName, "type", "ANONYMOUS_NETWORK"),
		resource.TestCheckResourceAttr(resourceFullName, "licensed", "true"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "true"),
		resource.TestCheckResourceAttr(resourceFullName, "default.result.level", "MEDIUM"),
	)

	minimalCheck := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
		resource.TestCheckResourceAttr(resourceFullName, "name", name),
		resource.TestCheckResourceAttr(resourceFullName, "compact_name", fmt.Sprintf("%s1", name)),
		resource.TestCheckNoResourceAttr(resourceFullName, "description"),
		resource.TestCheckResourceAttr(resourceFullName, "type", "ANONYMOUS_NETWORK"),
		resource.TestCheckResourceAttr(resourceFullName, "licensed", "true"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "true"),
		resource.TestCheckNoResourceAttr(resourceFullName, "default.result.level"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckRiskPredictorDestroy,
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
		resource.TestCheckResourceAttr(resourceFullName, "composition.condition_json_import", "{\"not\":{\"or\":[{\"equals\":0,\"type\":\"VALUE_COMPARISON\",\"value\":\"${details.counters.predictorLevels.medium}\"},{\"equals\":\"High\",\"type\":\"VALUE_COMPARISON\",\"value\":\"${details.geoVelocity.level}\"},{\"and\":[{\"equals\":\"High\",\"type\":\"VALUE_COMPARISON\",\"value\":\"${details.anonymousNetwork.level}\"}],\"type\":\"AND\"}],\"type\":\"OR\"},\"type\":\"NOT\"}"),
		resource.TestCheckResourceAttr(resourceFullName, "composition.condition", "{\"not\":{\"or\":[{\"equals\":0,\"type\":\"VALUE_COMPARISON\",\"value\":\"${details.counters.predictorLevels.medium}\"},{\"equals\":\"High\",\"type\":\"VALUE_COMPARISON\",\"value\":\"${details.geoVelocity.level}\"},{\"and\":[{\"equals\":\"High\",\"type\":\"VALUE_COMPARISON\",\"value\":\"${details.anonymousNetwork.level}\"}],\"type\":\"AND\"}],\"type\":\"OR\"},\"type\":\"NOT\"}"),
		resource.TestCheckResourceAttr(resourceFullName, "composition.level", "HIGH"),
	)

	fullCheck2 := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "type", "COMPOSITE"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "true"),
		resource.TestCheckResourceAttr(resourceFullName, "composition.condition_json_import", "{\"and\":[{\"equals\":5,\"type\":\"VALUE_COMPARISON\",\"value\":\"${details.counters.predictorLevels.medium}\"},{\"equals\":\"low\",\"type\":\"VALUE_COMPARISON\",\"value\":\"${details.anonymousNetwork.level}\"},{\"and\":[{\"equals\":\"high\",\"type\":\"VALUE_COMPARISON\",\"value\":\"${details.anonymousNetwork.level}\"},{\"or\":[{\"notEquals\":\"high\",\"type\":\"VALUE_COMPARISON\",\"value\":\"${details.anonymousNetwork.level}\"}]}]}]}"),
		resource.TestCheckResourceAttr(resourceFullName, "composition.condition", "{\"and\":[{\"equals\":5,\"type\":\"VALUE_COMPARISON\",\"value\":\"${details.counters.predictorLevels.medium}\"},{\"equals\":\"Low\",\"type\":\"VALUE_COMPARISON\",\"value\":\"${details.anonymousNetwork.level}\"},{\"and\":[{\"equals\":\"High\",\"type\":\"VALUE_COMPARISON\",\"value\":\"${details.anonymousNetwork.level}\"},{\"or\":[{\"notEquals\":\"high\",\"type\":\"VALUE_COMPARISON\",\"value\":\"${details.anonymousNetwork.level}\"}],\"type\":\"OR\"}],\"type\":\"AND\"}],\"type\":\"AND\"}"),
		resource.TestCheckResourceAttr(resourceFullName, "composition.level", "LOW"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckRiskPredictorDestroy,
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

func TestAccRiskPredictor_Anonymous_Network(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_risk_predictor.%s", resourceName)

	name := resourceName

	fullCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "type", "ANONYMOUS_NETWORK"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "true"),
		resource.TestCheckResourceAttr(resourceFullName, "default.result.level", "MEDIUM"),
		resource.TestCheckResourceAttr(resourceFullName, "allowed_cidr_list.#", "3"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "allowed_cidr_list.*", "192.168.0.0/24"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "allowed_cidr_list.*", "10.0.0.0/8"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "allowed_cidr_list.*", "172.16.0.0/12"),
	)

	minimalCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "type", "ANONYMOUS_NETWORK"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "true"),
		resource.TestCheckNoResourceAttr(resourceFullName, "default.result.level"),
		resource.TestCheckResourceAttr(resourceFullName, "allowed_cidr_list.#", "0"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckRiskPredictorDestroy,
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
		resource.TestCheckResourceAttr(resourceFullName, "allowed_cidr_list.#", "3"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "allowed_cidr_list.*", "192.168.0.0/24"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "allowed_cidr_list.*", "10.0.0.0/8"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "allowed_cidr_list.*", "172.16.0.0/12"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckRiskPredictorDestroyUndeletable,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			{
				Config: testAccRiskPredictorConfig_Anonymous_Network_OverwriteUndeletable(resourceName, name, compactName),
				Check:  fullCheck,
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
		resource.TestCheckResourceAttr(resourceFullName, "allowed_cidr_list.#", "3"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "allowed_cidr_list.*", "192.168.0.0/24"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "allowed_cidr_list.*", "10.0.0.0/8"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "allowed_cidr_list.*", "172.16.0.0/12"),
	)

	minimalCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "type", "GEO_VELOCITY"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "true"),
		resource.TestCheckNoResourceAttr(resourceFullName, "default.result.level"),
		resource.TestCheckResourceAttr(resourceFullName, "allowed_cidr_list.#", "0"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckRiskPredictorDestroy,
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
		resource.TestCheckResourceAttr(resourceFullName, "allowed_cidr_list.#", "3"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "allowed_cidr_list.*", "192.168.0.0/24"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "allowed_cidr_list.*", "10.0.0.0/8"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "allowed_cidr_list.*", "172.16.0.0/12"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckRiskPredictorDestroyUndeletable,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			{
				Config: testAccRiskPredictorConfig_Geovelocity_OverwriteUndeletable(resourceName, name, compactName),
				Check:  fullCheck,
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
		resource.TestCheckResourceAttr(resourceFullName, "allowed_cidr_list.#", "3"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "allowed_cidr_list.*", "192.168.0.0/24"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "allowed_cidr_list.*", "10.0.0.0/8"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "allowed_cidr_list.*", "172.16.0.0/12"),
	)

	minimalCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "type", "IP_REPUTATION"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "true"),
		resource.TestCheckNoResourceAttr(resourceFullName, "default.result.level"),
		resource.TestCheckResourceAttr(resourceFullName, "allowed_cidr_list.#", "0"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckRiskPredictorDestroy,
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
		resource.TestCheckResourceAttr(resourceFullName, "allowed_cidr_list.#", "3"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "allowed_cidr_list.*", "192.168.0.0/24"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "allowed_cidr_list.*", "10.0.0.0/8"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "allowed_cidr_list.*", "172.16.0.0/12"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckRiskPredictorDestroyUndeletable,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			{
				Config: testAccRiskPredictorConfig_IP_Reputation_OverwriteUndeletable(resourceName, name, compactName),
				Check:  fullCheck,
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
		resource.TestCheckResourceAttr(resourceFullName, "custom_map.contains", "${event.myshop}"),
		resource.TestCheckResourceAttr(resourceFullName, "custom_map.type", "RANGE"),
		resource.TestCheckResourceAttr(resourceFullName, "custom_map.between_ranges.high.max_score", "6"),
		resource.TestCheckResourceAttr(resourceFullName, "custom_map.between_ranges.high.min_score", "5"),
		resource.TestCheckResourceAttr(resourceFullName, "custom_map.between_ranges.medium.max_score", "4"),
		resource.TestCheckResourceAttr(resourceFullName, "custom_map.between_ranges.medium.min_score", "3"),
		resource.TestCheckResourceAttr(resourceFullName, "custom_map.between_ranges.low.max_score", "2"),
		resource.TestCheckResourceAttr(resourceFullName, "custom_map.between_ranges.low.min_score", "1"),
	)

	minimalCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "type", "MAP"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "true"),
		resource.TestCheckNoResourceAttr(resourceFullName, "default.result.level"),
		resource.TestCheckResourceAttr(resourceFullName, "custom_map.contains", "${event.myshop}"),
		resource.TestCheckResourceAttr(resourceFullName, "custom_map.type", "RANGE"),
		resource.TestCheckNoResourceAttr(resourceFullName, "custom_map.between_ranges.high.max_score"),
		resource.TestCheckNoResourceAttr(resourceFullName, "custom_map.between_ranges.high.min_score"),
		resource.TestCheckResourceAttr(resourceFullName, "custom_map.between_ranges.medium.max_score", "4"),
		resource.TestCheckResourceAttr(resourceFullName, "custom_map.between_ranges.medium.min_score", "3"),
		resource.TestCheckNoResourceAttr(resourceFullName, "custom_map.between_ranges.low.max_score"),
		resource.TestCheckNoResourceAttr(resourceFullName, "custom_map.between_ranges.low.min_score"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckRiskPredictorDestroy,
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
		resource.TestCheckResourceAttr(resourceFullName, "custom_map.contains", "${event.myshop}"),
		resource.TestCheckResourceAttr(resourceFullName, "custom_map.type", "IP_RANGE"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "custom_map.ip_ranges.high.values.*", "192.168.0.0/24"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "custom_map.ip_ranges.high.values.*", "10.0.0.0/8"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "custom_map.ip_ranges.high.values.*", "172.16.0.0/12"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "custom_map.ip_ranges.medium.values.*", "192.0.2.0/24"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "custom_map.ip_ranges.medium.values.*", "192.168.1.0/26"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "custom_map.ip_ranges.medium.values.*", "10.10.0.0/16"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "custom_map.ip_ranges.low.values.*", "172.16.0.0/16"),
	)

	minimalCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "type", "MAP"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "true"),
		resource.TestCheckNoResourceAttr(resourceFullName, "default.result.level"),
		resource.TestCheckResourceAttr(resourceFullName, "custom_map.contains", "${event.myshop}"),
		resource.TestCheckResourceAttr(resourceFullName, "custom_map.type", "IP_RANGE"),
		resource.TestCheckNoResourceAttr(resourceFullName, "custom_map.ip_ranges.high.values"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "custom_map.ip_ranges.medium.values.*", "192.0.2.0/24"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "custom_map.ip_ranges.medium.values.*", "10.0.0.0/8"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "custom_map.ip_ranges.medium.values.*", "172.16.0.0/12"),
		resource.TestCheckNoResourceAttr(resourceFullName, "custom_map.ip_ranges.low.values"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckRiskPredictorDestroy,
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
		resource.TestCheckResourceAttr(resourceFullName, "custom_map.contains", "${event.myshop}"),
		resource.TestCheckResourceAttr(resourceFullName, "custom_map.type", "STRING_LIST"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "custom_map.string_list.high.values.*", "HIGH"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "custom_map.string_list.high.values.*", "HIGH321"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "custom_map.string_list.high.values.*", "HIGH123"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "custom_map.string_list.medium.values.*", "MEDIUM"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "custom_map.string_list.medium.values.*", "MED321"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "custom_map.string_list.medium.values.*", "MED123"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "custom_map.string_list.low.values.*", "LOW"),
	)

	minimalCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "type", "MAP"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "true"),
		resource.TestCheckNoResourceAttr(resourceFullName, "default.result.level"),
		resource.TestCheckResourceAttr(resourceFullName, "custom_map.contains", "${event.myshop}"),
		resource.TestCheckResourceAttr(resourceFullName, "custom_map.type", "STRING_LIST"),
		resource.TestCheckNoResourceAttr(resourceFullName, "custom_map.string_list.high.values"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "custom_map.string_list.medium.values.*", "MEDIUM"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "custom_map.string_list.medium.values.*", "MED321"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "custom_map.string_list.medium.values.*", "MED123"),
		resource.TestCheckNoResourceAttr(resourceFullName, "custom_map.string_list.low.values"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckRiskPredictorDestroy,
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
		},
	})
}

func TestAccRiskPredictor_NewDevice(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_risk_predictor.%s", resourceName)

	name := resourceName

	fullCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "type", "DEVICE"),
		resource.TestCheckResourceAttr(resourceFullName, "detect", "NEW_DEVICE"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "true"),
		resource.TestCheckResourceAttr(resourceFullName, "default.result.level", "MEDIUM"),
		resource.TestCheckResourceAttr(resourceFullName, "activation_at", "2023-05-02T00:00:00Z"),
	)

	minimalCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "type", "DEVICE"),
		resource.TestCheckResourceAttr(resourceFullName, "detect", "NEW_DEVICE"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "true"),
		resource.TestCheckNoResourceAttr(resourceFullName, "default.result.level"),
		resource.TestCheckNoResourceAttr(resourceFullName, "activation_at"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckRiskPredictorDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			{
				Config: testAccRiskPredictorConfig_NewDevice_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config:  testAccRiskPredictorConfig_NewDevice_Full(resourceName, name),
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
				Config: testAccRiskPredictorConfig_NewDevice_Full(resourceName, name),
				Check:  fullCheck,
			},
			{
				Config: testAccRiskPredictorConfig_NewDevice_Minimal(resourceName, name),
				Check:  minimalCheck,
			},
			{
				Config: testAccRiskPredictorConfig_NewDevice_Full(resourceName, name),
				Check:  fullCheck,
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

	fullCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "name", name),
		resource.TestCheckResourceAttr(resourceFullName, "compact_name", "newDevice"),
		resource.TestCheckResourceAttr(resourceFullName, "type", "DEVICE"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "false"),
		resource.TestCheckResourceAttr(resourceFullName, "default.result.level", "MEDIUM"),
		resource.TestCheckResourceAttr(resourceFullName, "detect", "NEW_DEVICE"),
		resource.TestCheckResourceAttr(resourceFullName, "activation_at", "2023-05-02T00:00:00Z"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckRiskPredictorDestroyUndeletable,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			{
				Config: testAccRiskPredictorConfig_NewDevice_OverwriteUndeletable(resourceName, name, compactName),
				Check:  fullCheck,
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
		resource.TestCheckResourceAttr(resourceFullName, "radius.distance", "100"),
		resource.TestCheckResourceAttr(resourceFullName, "radius.unit", "miles"),
		resource.TestCheckResourceAttr(resourceFullName, "days", "50"),
	)

	minimalCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "type", "USER_LOCATION_ANOMALY"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "true"),
		resource.TestCheckNoResourceAttr(resourceFullName, "default.result.level"),
		resource.TestCheckResourceAttr(resourceFullName, "radius.distance", "51"),
		resource.TestCheckResourceAttr(resourceFullName, "radius.unit", "kilometers"),
		resource.TestCheckResourceAttr(resourceFullName, "days", "50"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckRiskPredictorDestroy,
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
		resource.TestCheckResourceAttr(resourceFullName, "radius.distance", "100"),
		resource.TestCheckResourceAttr(resourceFullName, "radius.unit", "miles"),
		resource.TestCheckResourceAttr(resourceFullName, "days", "50"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckRiskPredictorDestroyUndeletable,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Full
			{
				Config: testAccRiskPredictorConfig_UserLocationAnomaly_OverwriteUndeletable(resourceName, name, compactName),
				Check:  fullCheck,
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
		resource.TestCheckResourceAttr(resourceFullName, "of", "${event.ip}"),
		resource.TestCheckResourceAttr(resourceFullName, "by.#", "1"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "by.*", "${event.user.id}"),
		resource.TestCheckResourceAttr(resourceFullName, "measure", "DISTINCT_COUNT"),
		resource.TestCheckResourceAttr(resourceFullName, "use.type", "POISSON_WITH_MAX"),
		resource.TestCheckResourceAttr(resourceFullName, "use.medium", "2"),
		resource.TestCheckResourceAttr(resourceFullName, "use.high", "4"),
		resource.TestCheckResourceAttr(resourceFullName, "fallback.strategy", "ENVIRONMENT_MAX"),
		resource.TestCheckResourceAttr(resourceFullName, "fallback.high", "30"),
		resource.TestCheckResourceAttr(resourceFullName, "fallback.medium", "20"),
		resource.TestCheckResourceAttr(resourceFullName, "every.unit", "HOUR"),
		resource.TestCheckResourceAttr(resourceFullName, "every.quantity", "1"),
		resource.TestCheckResourceAttr(resourceFullName, "every.min_sample", "5"),
		resource.TestCheckResourceAttr(resourceFullName, "sliding_window.unit", "DAY"),
		resource.TestCheckResourceAttr(resourceFullName, "sliding_window.quantity", "7"),
		resource.TestCheckResourceAttr(resourceFullName, "sliding_window.min_sample", "3"),
	)

	byIPCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "type", "VELOCITY"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "true"),
		resource.TestCheckResourceAttr(resourceFullName, "of", "${event.user.id}"),
		resource.TestCheckResourceAttr(resourceFullName, "by.#", "1"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "by.*", "${event.ip}"),
		resource.TestCheckResourceAttr(resourceFullName, "measure", "DISTINCT_COUNT"),
		resource.TestCheckResourceAttr(resourceFullName, "use.type", "POISSON_WITH_MAX"),
		resource.TestCheckResourceAttr(resourceFullName, "use.medium", "2"),
		resource.TestCheckResourceAttr(resourceFullName, "use.high", "4"),
		resource.TestCheckResourceAttr(resourceFullName, "fallback.strategy", "ENVIRONMENT_MAX"),
		resource.TestCheckResourceAttr(resourceFullName, "fallback.high", "3500"),
		resource.TestCheckResourceAttr(resourceFullName, "fallback.medium", "2500"),
		resource.TestCheckResourceAttr(resourceFullName, "every.unit", "HOUR"),
		resource.TestCheckResourceAttr(resourceFullName, "every.quantity", "1"),
		resource.TestCheckResourceAttr(resourceFullName, "every.min_sample", "5"),
		resource.TestCheckResourceAttr(resourceFullName, "sliding_window.unit", "DAY"),
		resource.TestCheckResourceAttr(resourceFullName, "sliding_window.quantity", "7"),
		resource.TestCheckResourceAttr(resourceFullName, "sliding_window.min_sample", "3"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckRiskPredictorDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// By User
			{
				Config: testAccRiskPredictorConfig_Velocity_ByUser_Full(resourceName, name),
				Check:  byUserCheck,
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
	compactName := "userLocationAnomaly"

	byUserCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "name", name),
		resource.TestCheckResourceAttr(resourceFullName, "compact_name", compactName),
		resource.TestCheckResourceAttr(resourceFullName, "type", "VELOCITY"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "false"),
		resource.TestCheckResourceAttr(resourceFullName, "of", "${event.ip}"),
		resource.TestCheckResourceAttr(resourceFullName, "by.#", "1"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "by.*", "${event.user.id}"),
		resource.TestCheckResourceAttr(resourceFullName, "measure", "DISTINCT_COUNT"),
		resource.TestCheckResourceAttr(resourceFullName, "use.type", "POISSON_WITH_MAX"),
		resource.TestCheckResourceAttr(resourceFullName, "use.medium", "2"),
		resource.TestCheckResourceAttr(resourceFullName, "use.high", "4"),
		resource.TestCheckResourceAttr(resourceFullName, "fallback.strategy", "ENVIRONMENT_MAX"),
		resource.TestCheckResourceAttr(resourceFullName, "fallback.high", "30"),
		resource.TestCheckResourceAttr(resourceFullName, "fallback.medium", "20"),
		resource.TestCheckResourceAttr(resourceFullName, "every.unit", "HOUR"),
		resource.TestCheckResourceAttr(resourceFullName, "every.quantity", "1"),
		resource.TestCheckResourceAttr(resourceFullName, "every.min_sample", "5"),
		resource.TestCheckResourceAttr(resourceFullName, "sliding_window.unit", "DAY"),
		resource.TestCheckResourceAttr(resourceFullName, "sliding_window.quantity", "7"),
		resource.TestCheckResourceAttr(resourceFullName, "sliding_window.min_sample", "3"),
	)

	byIPCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "name", name),
		resource.TestCheckResourceAttr(resourceFullName, "compact_name", compactName),
		resource.TestCheckResourceAttr(resourceFullName, "type", "VELOCITY"),
		resource.TestCheckResourceAttr(resourceFullName, "of", "${event.user.id}"),
		resource.TestCheckResourceAttr(resourceFullName, "by.#", "1"),
		resource.TestCheckTypeSetElemAttr(resourceFullName, "by.*", "${event.ip}"),
		resource.TestCheckResourceAttr(resourceFullName, "measure", "DISTINCT_COUNT"),
		resource.TestCheckResourceAttr(resourceFullName, "use.type", "POISSON_WITH_MAX"),
		resource.TestCheckResourceAttr(resourceFullName, "use.medium", "2"),
		resource.TestCheckResourceAttr(resourceFullName, "use.high", "4"),
		resource.TestCheckResourceAttr(resourceFullName, "fallback.strategy", "ENVIRONMENT_MAX"),
		resource.TestCheckResourceAttr(resourceFullName, "fallback.high", "3500"),
		resource.TestCheckResourceAttr(resourceFullName, "fallback.medium", "2500"),
		resource.TestCheckResourceAttr(resourceFullName, "every.unit", "HOUR"),
		resource.TestCheckResourceAttr(resourceFullName, "every.quantity", "1"),
		resource.TestCheckResourceAttr(resourceFullName, "every.min_sample", "5"),
		resource.TestCheckResourceAttr(resourceFullName, "sliding_window.unit", "DAY"),
		resource.TestCheckResourceAttr(resourceFullName, "sliding_window.quantity", "7"),
		resource.TestCheckResourceAttr(resourceFullName, "sliding_window.min_sample", "3"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckRiskPredictorDestroyUndeletable,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// By User
			{
				Config: testAccRiskPredictorConfig_Velocity_ByUser_Full_OverwriteUndeletable(resourceName, name, compactName),
				Check:  byUserCheck,
			},
			{
				Config:  testAccRiskPredictorConfig_Velocity_ByUser_Full_OverwriteUndeletable(resourceName, name, compactName),
				Destroy: true,
			},
			// By IP
			{
				Config: testAccRiskPredictorConfig_Velocity_ByIP_Full_OverwriteUndeletable(resourceName, name, compactName),
				Check:  byIPCheck,
			},
			{
				Config:  testAccRiskPredictorConfig_Velocity_ByIP_Full_OverwriteUndeletable(resourceName, name, compactName),
				Destroy: true,
			},
			// Change
			{
				Config: testAccRiskPredictorConfig_Velocity_ByUser_Full_OverwriteUndeletable(resourceName, name, compactName),
				Check:  byUserCheck,
			},
			{
				Config: testAccRiskPredictorConfig_Velocity_ByIP_Full_OverwriteUndeletable(resourceName, name, compactName),
				Check:  byIPCheck,
			},
			{
				Config: testAccRiskPredictorConfig_Velocity_ByUser_Full_OverwriteUndeletable(resourceName, name, compactName),
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
		resource.TestCheckResourceAttr(resourceFullName, "prediction_model.name", "points"),
	)

	byOrgCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "type", "USER_RISK_BEHAVIOR"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "true"),
		resource.TestCheckResourceAttr(resourceFullName, "prediction_model.name", "login_anomaly_statistic"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckRiskPredictorDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// By User
			{
				Config: testAccRiskPredictorConfig_UserRiskBehavior_ByUser_Full(resourceName, name),
				Check:  byUserCheck,
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
	compactName := "userLocationAnomaly"

	byUserCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "name", name),
		resource.TestCheckResourceAttr(resourceFullName, "compact_name", compactName),
		resource.TestCheckResourceAttr(resourceFullName, "type", "USER_RISK_BEHAVIOR"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "false"),
		resource.TestCheckResourceAttr(resourceFullName, "prediction_model.name", "points"),
	)

	byIPCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "name", name),
		resource.TestCheckResourceAttr(resourceFullName, "compact_name", compactName),
		resource.TestCheckResourceAttr(resourceFullName, "type", "USER_RISK_BEHAVIOR"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "false"),
		resource.TestCheckResourceAttr(resourceFullName, "prediction_model.name", "login_anomaly_statistic"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckRiskPredictorDestroyUndeletable,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// By User
			{
				Config: testAccRiskPredictorConfig_UserRiskBehavior_ByUser_Full_OverwriteUndeletable(resourceName, name, compactName),
				Check:  byUserCheck,
			},
			{
				Config:  testAccRiskPredictorConfig_UserRiskBehavior_ByUser_Full_OverwriteUndeletable(resourceName, name, compactName),
				Destroy: true,
			},
			// By Org
			{
				Config: testAccRiskPredictorConfig_UserRiskBehavior_ByOrg_Full_OverwriteUndeletable(resourceName, name, compactName),
				Check:  byIPCheck,
			},
			{
				Config:  testAccRiskPredictorConfig_UserRiskBehavior_ByOrg_Full_OverwriteUndeletable(resourceName, name, compactName),
				Destroy: true,
			},
			// Change
			{
				Config: testAccRiskPredictorConfig_UserRiskBehavior_ByUser_Full_OverwriteUndeletable(resourceName, name, compactName),
				Check:  byUserCheck,
			},
			{
				Config: testAccRiskPredictorConfig_UserRiskBehavior_ByOrg_Full_OverwriteUndeletable(resourceName, name, compactName),
				Check:  byIPCheck,
			},
			{
				Config: testAccRiskPredictorConfig_UserRiskBehavior_ByUser_Full_OverwriteUndeletable(resourceName, name, compactName),
				Check:  byUserCheck,
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

  type = "ANONYMOUS_NETWORK"

}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccRiskPredictorConfig_Full(resourceName, name string) string {
	return testAccRiskPredictorConfig_Anonymous_Network_Full(resourceName, name)
}

func testAccRiskPredictorConfig_Minimal(resourceName, name string) string {
	return testAccRiskPredictorConfig_Anonymous_Network_Minimal(resourceName, name)
}

func testAccRiskPredictorConfig_Anonymous_Network_Full(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[3]s1"
  description  = "When my wife is upset, I let her colour in my black and white tattoos.  She just needs a shoulder to crayon.."

  type = "ANONYMOUS_NETWORK"

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  allowed_cidr_list = [
    "10.0.0.0/8",
    "172.16.0.0/12",
    "192.168.0.0/24"
  ]

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPredictorConfig_Composite_Full_1(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[3]s1"
  description  = "When my wife is upset, I let her colour in my black and white tattoos.  She just needs a shoulder to crayon.."

  type = "COMPOSITE"

  composition = {
    level = "HIGH"

    condition_json_import = jsonencode({
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

  type = "COMPOSITE"

  composition = {
    level = "LOW"

    condition_json_import = jsonencode({
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

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPredictorConfig_Composite_InvalidJSON(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[3]s1"

  type = "COMPOSITE"

  composition = {
    level = "LOW"

    condition_json_import = jsonencode({})
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

  type = "ANONYMOUS_NETWORK"
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPredictorConfig_Anonymous_Network_OverwriteUndeletable(resourceName, name, compactName string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[4]s"

  type = "ANONYMOUS_NETWORK"

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  allowed_cidr_list = [
    "10.0.0.0/8",
    "172.16.0.0/12",
    "192.168.0.0/24"
  ]

}`, acctest.GenericSandboxEnvironment(), resourceName, name, compactName)
}

func testAccRiskPredictorConfig_Geovelocity_Full(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[3]s1"
  description  = "When my wife is upset, I let her colour in my black and white tattoos.  She just needs a shoulder to crayon.."

  type = "GEO_VELOCITY"

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  allowed_cidr_list = [
    "10.0.0.0/8",
    "172.16.0.0/12",
    "192.168.0.0/24"
  ]

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPredictorConfig_Geovelocity_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[3]s1"

  type = "GEO_VELOCITY"

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPredictorConfig_Geovelocity_OverwriteUndeletable(resourceName, name, compactName string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[4]s"

  type = "GEO_VELOCITY"

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  allowed_cidr_list = [
    "10.0.0.0/8",
    "172.16.0.0/12",
    "192.168.0.0/24"
  ]

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

  type = "IP_REPUTATION"

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  allowed_cidr_list = [
    "10.0.0.0/8",
    "172.16.0.0/12",
    "192.168.0.0/24"
  ]

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPredictorConfig_IPReputation_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[3]s1"

  type = "IP_REPUTATION"

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPredictorConfig_IP_Reputation_OverwriteUndeletable(resourceName, name, compactName string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[4]s"

  type = "IP_REPUTATION"

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  allowed_cidr_list = [
    "10.0.0.0/8",
    "172.16.0.0/12",
    "192.168.0.0/24"
  ]

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

  type = "MAP"

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  custom_map = {
    contains = "$${event.myshop}"

    between_ranges = {
      high = {
        max_score = 6
        min_score = 5
      }

      medium = {
        max_score = 4
        min_score = 3
      }

      low = {
        max_score = 2
        min_score = 1
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

  type = "MAP"

  custom_map = {
    contains = "$${event.myshop}"

    between_ranges = {
      medium = {
        max_score = 4
        min_score = 3
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

  type = "MAP"

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  custom_map = {
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

  type = "MAP"

  custom_map = {
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

  type = "MAP"

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  custom_map = {
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

  type = "MAP"

  custom_map = {
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

func testAccRiskPredictorConfig_NewDevice_Full(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[3]s1"
  description  = "When my wife is upset, I let her colour in my black and white tattoos.  She just needs a shoulder to crayon.."

  type   = "DEVICE"
  detect = "NEW_DEVICE"

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  activation_at = "2023-05-02T00:00:00Z"
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPredictorConfig_NewDevice_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[3]s1"

  type = "DEVICE"
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPredictorConfig_NewDevice_OverwriteUndeletable(resourceName, name, compactName string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[4]s"

  type   = "DEVICE"
  detect = "NEW_DEVICE"

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  activation_at = "2023-05-02T00:00:00Z"
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

  type = "USER_LOCATION_ANOMALY"

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  radius = {
    distance = 100
    unit     = "miles"
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

  type = "USER_LOCATION_ANOMALY"

  radius = {
    distance = 51

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

  type = "USER_LOCATION_ANOMALY"

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  radius = {
    distance = 100
    unit     = "miles"
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

  type = "VELOCITY"

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  of = "$${event.ip}"

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPredictorConfig_Velocity_ByIP_Full(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[3]s1"

  type = "VELOCITY"

  of = "$${event.user.id}"

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPredictorConfig_Velocity_ByUser_Full_OverwriteUndeletable(resourceName, name, compactName string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[4]s"

  type = "VELOCITY"

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  of = "$${event.ip}"
}`, acctest.GenericSandboxEnvironment(), resourceName, name, compactName)
}

func testAccRiskPredictorConfig_Velocity_ByIP_Full_OverwriteUndeletable(resourceName, name, compactName string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[4]s"

  type = "VELOCITY"

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  of = "$${event.user.id}"
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

  type = "USER_RISK_BEHAVIOR"

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  prediction_model = {
    name = "points"
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

  type = "USER_RISK_BEHAVIOR"

  prediction_model = {
    name = "login_anomaly_statistic"
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

  type = "USER_RISK_BEHAVIOR"

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  prediction_model = {
    name = "points"
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

  type = "USER_RISK_BEHAVIOR"

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  prediction_model = {
    name = "login_anomaly_statistic"
  }
}`, acctest.GenericSandboxEnvironment(), resourceName, name, compactName)
}
