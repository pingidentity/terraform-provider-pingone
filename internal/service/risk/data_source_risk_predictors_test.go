// Copyright © 2026 Ping Identity Corporation

package risk_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/risk"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccRiskPredictorsDataSource_List(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullName := fmt.Sprintf("data.pingone_risk_predictors.%s", resourceName)

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
				Config: testAccRiskPredictorsDataSourceConfig_List(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "ids.0"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "ids.1"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "ids.2"),
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources[dataSourceFullName]
						if !ok {
							return fmt.Errorf("Not found: %s", dataSourceFullName)
						}

						// Check for all 3 created resources
						for i := 1; i <= 3; i++ {
							resName := fmt.Sprintf("pingone_risk_predictor.%s-%d", resourceName, i)
							predictorRs, ok := s.RootModule().Resources[resName]
							if !ok {
								return fmt.Errorf("Not found: %s", resName)
							}
							predictorID := predictorRs.Primary.ID

							// Iterate through attributes starting with "ids."
							found := false
							for k, v := range rs.Primary.Attributes {
								if strings.HasPrefix(k, "ids.") && v == predictorID {
									found = true
									break
								}
							}

							if !found {
								return fmt.Errorf("Predictor ID %s (%s) not found in data source ids list", predictorID, resName)
							}
						}
						return nil
					},
				),
			},
		},
	})
}

func testAccRiskPredictorsDataSourceConfig_List(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_predictor" "%[2]s-1" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s-1"
  compact_name = "%[3]s1"
  description  = "Test 1"

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  predictor_anonymous_network = {
    allowed_cidr_list = [
      "10.0.0.0/8"
    ]
  }

}

resource "pingone_risk_predictor" "%[2]s-2" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s-2"
  compact_name = "%[3]s2"
  description  = "Test 2"

  default = {
    result = {
      level = "HIGH"
    }
  }

  predictor_velocity = {
    of = "$${event.user.id}"
  }

}

resource "pingone_risk_predictor" "%[2]s-3" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s-3"
  compact_name = "%[3]s3"
  description  = "Test 3"

  default = {
    result = {
      level = "LOW"
    }
  }

  predictor_bot_detection = {
    include_repeated_events_without_sdk = true
  }

}

data "pingone_risk_predictors" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  depends_on = [
    pingone_risk_predictor.%[2]s-1,
    pingone_risk_predictor.%[2]s-2,
    pingone_risk_predictor.%[2]s-3,
  ]
}
`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
