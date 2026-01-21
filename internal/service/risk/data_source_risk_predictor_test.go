// Copyright © 2025 Ping Identity Corporation

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

func TestAccRiskPredictorDataSource_RiskPredictorID(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_risk_predictor.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.pingone_risk_predictor.%s", resourceName)

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
			{
				Config: testAccRiskPredictorDataSourceConfig_ByID(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "name", resourceFullName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "compact_name", resourceFullName, "compact_name"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "description", resourceFullName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "type", resourceFullName, "type"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "licensed", resourceFullName, "licensed"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "deletable", resourceFullName, "deletable"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "default.result.level", resourceFullName, "default.result.level"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "predictor_anonymous_network.allowed_cidr_list.#", resourceFullName, "predictor_anonymous_network.allowed_cidr_list.#"),
				),
			},
		},
	})
}

func TestAccRiskPredictorDataSource_Name(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_risk_predictor.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.pingone_risk_predictor.%s", resourceName)

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
			{
				Config: testAccRiskPredictorDataSourceConfig_ByName(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "name", resourceFullName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "compact_name", resourceFullName, "compact_name"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "description", resourceFullName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "type", resourceFullName, "type"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "licensed", resourceFullName, "licensed"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "deletable", resourceFullName, "deletable"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "default.result.level", resourceFullName, "default.result.level"),
					resource.TestCheckResourceAttrPair(dataSourceFullName, "predictor_anonymous_network.allowed_cidr_list.#", resourceFullName, "predictor_anonymous_network.allowed_cidr_list.#"),
				),
			},
		},
	})
}

func TestAccRiskPredictorDataSource_NotFound(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

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
			{
				Config:      testAccRiskPredictorDataSourceConfig_NotFoundByName(resourceName),
				ExpectError: regexp.MustCompile("Risk Predictor with name .* not found"),
			},
			{
				Config:      testAccRiskPredictorDataSourceConfig_NotFoundByID(resourceName),
				ExpectError: regexp.MustCompile("Risk Predictor with ID .* not found"),
			},
		},
	})
}

func testAccRiskPredictorDataSourceConfig_ByID(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[3]s1"
  description  = "My risk predictor description goes here."

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

}

data "pingone_risk_predictor" "%[2]s" {
  environment_id    = data.pingone_environment.general_test.id
  risk_predictor_id = pingone_risk_predictor.%[2]s.id
}
`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPredictorDataSourceConfig_ByName(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[3]s1"
  description  = "My risk predictor description goes here."

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

}

data "pingone_risk_predictor" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
  name           = pingone_risk_predictor.%[2]s.name
}
`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPredictorDataSourceConfig_NotFoundByName(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_risk_predictor" "%[2]s-notfound" {
  environment_id = data.pingone_environment.general_test.id
  name           = "test_not_found_name"
}
`, acctest.GenericSandboxEnvironment(), resourceName)
}

func testAccRiskPredictorDataSourceConfig_NotFoundByID(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_risk_predictor" "%[2]s-notfound" {
  environment_id    = data.pingone_environment.general_test.id
  risk_predictor_id = "9c052a8a-14be-44e4-8f07-2662569994ce"
}
`, acctest.GenericSandboxEnvironment(), resourceName)
}
