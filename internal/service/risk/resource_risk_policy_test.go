package risk_test

import (
	"context"
	"fmt"
	"os"
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
		resource.TestCheckResourceAttr(resourceFullName, "type", "COMPOSITE"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "true"),
		resource.TestCheckResourceAttr(resourceFullName, "policy_composite.composition.condition_json", "{\"not\":{\"or\":[{\"equals\":0,\"type\":\"VALUE_COMPARISON\",\"value\":\"${details.counters.policyLevels.medium}\"},{\"equals\":\"High\",\"type\":\"VALUE_COMPARISON\",\"value\":\"${details.geoVelocity.level}\"},{\"and\":[{\"equals\":\"High\",\"type\":\"VALUE_COMPARISON\",\"value\":\"${details.anonymousNetwork.level}\"}],\"type\":\"AND\"}],\"type\":\"OR\"},\"type\":\"NOT\"}"),
		resource.TestCheckResourceAttr(resourceFullName, "policy_composite.composition.condition", "{\"not\":{\"or\":[{\"equals\":0,\"type\":\"VALUE_COMPARISON\",\"value\":\"${details.counters.policyLevels.medium}\"},{\"equals\":\"High\",\"type\":\"VALUE_COMPARISON\",\"value\":\"${details.geoVelocity.level}\"},{\"and\":[{\"equals\":\"High\",\"type\":\"VALUE_COMPARISON\",\"value\":\"${details.anonymousNetwork.level}\"}],\"type\":\"AND\"}],\"type\":\"OR\"},\"type\":\"NOT\"}"),
		resource.TestCheckResourceAttr(resourceFullName, "policy_composite.composition.level", "HIGH"),
	)

	minimalCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "type", "COMPOSITE"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "true"),
		resource.TestCheckResourceAttr(resourceFullName, "policy_composite.composition.condition_json", "{\"and\":[{\"equals\":5,\"type\":\"VALUE_COMPARISON\",\"value\":\"${details.counters.policyLevels.medium}\"},{\"equals\":\"low\",\"type\":\"VALUE_COMPARISON\",\"value\":\"${details.anonymousNetwork.level}\"},{\"and\":[{\"equals\":\"high\",\"type\":\"VALUE_COMPARISON\",\"value\":\"${details.anonymousNetwork.level}\"},{\"or\":[{\"notEquals\":\"high\",\"type\":\"VALUE_COMPARISON\",\"value\":\"${details.anonymousNetwork.level}\"}]}]}]}"),
		resource.TestCheckResourceAttr(resourceFullName, "policy_composite.composition.condition", "{\"and\":[{\"equals\":5,\"type\":\"VALUE_COMPARISON\",\"value\":\"${details.counters.policyLevels.medium}\"},{\"equals\":\"Low\",\"type\":\"VALUE_COMPARISON\",\"value\":\"${details.anonymousNetwork.level}\"},{\"and\":[{\"equals\":\"High\",\"type\":\"VALUE_COMPARISON\",\"value\":\"${details.anonymousNetwork.level}\"},{\"or\":[{\"notEquals\":\"high\",\"type\":\"VALUE_COMPARISON\",\"value\":\"${details.anonymousNetwork.level}\"}],\"type\":\"OR\"}],\"type\":\"AND\"}],\"type\":\"AND\"}"),
		resource.TestCheckResourceAttr(resourceFullName, "policy_composite.composition.level", "LOW"),
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
		},
	})
}

func TestAccRiskPolicy_Weights(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_risk_policy.%s", resourceName)

	name := resourceName

	fullCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "type", "COMPOSITE"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "true"),
		resource.TestCheckResourceAttr(resourceFullName, "policy_composite.composition.condition_json", "{\"not\":{\"or\":[{\"equals\":0,\"type\":\"VALUE_COMPARISON\",\"value\":\"${details.counters.policyLevels.medium}\"},{\"equals\":\"High\",\"type\":\"VALUE_COMPARISON\",\"value\":\"${details.geoVelocity.level}\"},{\"and\":[{\"equals\":\"High\",\"type\":\"VALUE_COMPARISON\",\"value\":\"${details.anonymousNetwork.level}\"}],\"type\":\"AND\"}],\"type\":\"OR\"},\"type\":\"NOT\"}"),
		resource.TestCheckResourceAttr(resourceFullName, "policy_composite.composition.condition", "{\"not\":{\"or\":[{\"equals\":0,\"type\":\"VALUE_COMPARISON\",\"value\":\"${details.counters.policyLevels.medium}\"},{\"equals\":\"High\",\"type\":\"VALUE_COMPARISON\",\"value\":\"${details.geoVelocity.level}\"},{\"and\":[{\"equals\":\"High\",\"type\":\"VALUE_COMPARISON\",\"value\":\"${details.anonymousNetwork.level}\"}],\"type\":\"AND\"}],\"type\":\"OR\"},\"type\":\"NOT\"}"),
		resource.TestCheckResourceAttr(resourceFullName, "policy_composite.composition.level", "HIGH"),
	)

	minimalCheck := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "type", "COMPOSITE"),
		resource.TestCheckResourceAttr(resourceFullName, "deletable", "true"),
		resource.TestCheckResourceAttr(resourceFullName, "policy_composite.composition.condition_json", "{\"and\":[{\"equals\":5,\"type\":\"VALUE_COMPARISON\",\"value\":\"${details.counters.policyLevels.medium}\"},{\"equals\":\"low\",\"type\":\"VALUE_COMPARISON\",\"value\":\"${details.anonymousNetwork.level}\"},{\"and\":[{\"equals\":\"high\",\"type\":\"VALUE_COMPARISON\",\"value\":\"${details.anonymousNetwork.level}\"},{\"or\":[{\"notEquals\":\"high\",\"type\":\"VALUE_COMPARISON\",\"value\":\"${details.anonymousNetwork.level}\"}]}]}]}"),
		resource.TestCheckResourceAttr(resourceFullName, "policy_composite.composition.condition", "{\"and\":[{\"equals\":5,\"type\":\"VALUE_COMPARISON\",\"value\":\"${details.counters.policyLevels.medium}\"},{\"equals\":\"Low\",\"type\":\"VALUE_COMPARISON\",\"value\":\"${details.anonymousNetwork.level}\"},{\"and\":[{\"equals\":\"High\",\"type\":\"VALUE_COMPARISON\",\"value\":\"${details.anonymousNetwork.level}\"},{\"or\":[{\"notEquals\":\"high\",\"type\":\"VALUE_COMPARISON\",\"value\":\"${details.anonymousNetwork.level}\"}],\"type\":\"OR\"}],\"type\":\"AND\"}],\"type\":\"AND\"}"),
		resource.TestCheckResourceAttr(resourceFullName, "policy_composite.composition.level", "LOW"),
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
		},
	})
}

func testAccRiskPolicyConfig_NewEnv(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
		%[1]s

resource "pingone_risk_policy" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  name         = "%[4]s"

  policy_scores = {
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

resource "pingone_risk_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[3]s1"
  description  = "When my wife is upset, I let her colour in my black and white tattoos.  She just needs a shoulder to crayon.."

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  policy_anonymous_network = {
    allowed_cidr_list = [
      "10.0.0.0/8",
      "172.16.0.0/12",
      "192.168.0.0/24"
    ]
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPolicyConfig_Scores_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[3]s1"

  policy_anonymous_network = {}
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPolicyConfig_Weights_Full(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[3]s1"
  description  = "When my wife is upset, I let her colour in my black and white tattoos.  She just needs a shoulder to crayon.."

  default = {
    result = {
      level = "MEDIUM"
    }
  }

  policy_anonymous_network = {
    allowed_cidr_list = [
      "10.0.0.0/8",
      "172.16.0.0/12",
      "192.168.0.0/24"
    ]
  }

}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}

func testAccRiskPolicyConfig_Weights_Minimal(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_risk_policy" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  name         = "%[3]s"
  compact_name = "%[3]s1"

  policy_anonymous_network = {}
}`, acctest.GenericSandboxEnvironment(), resourceName, name)
}
