// Copyright © 2025 Ping Identity Corporation

package verify_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	acctestlegacysdk "github.com/pingidentity/terraform-provider-pingone/internal/acctest/legacysdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/verify"
	validation "github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccVerifyPoliciesDataSource_NoFilter(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullName := fmt.Sprintf("data.pingone_verify_policies.%s", resourceName)

	name := acctest.ResourceNameGen()

	environmentName := acctest.ResourceNameGenEnvironment()
	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	findVerifyPolicies := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(dataSourceFullName, "id", validation.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", validation.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(dataSourceFullName, "ids.#", "5"), // includes environment default policy
		resource.TestMatchResourceAttr(dataSourceFullName, "ids.0", validation.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(dataSourceFullName, "ids.1", validation.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(dataSourceFullName, "ids.2", validation.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(dataSourceFullName, "ids.3", validation.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(dataSourceFullName, "ids.4", validation.P1ResourceIDRegexpFullString),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             verify.VerifyPolicy_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVerifyPolicies_NoFilter(environmentName, licenseID, resourceName, name),
				Check:  findVerifyPolicies,
			},
			{
				Config:  testAccVerifyPolicies_NoFilter(environmentName, licenseID, resourceName, name),
				Destroy: true,
			},
		},
	})
}

func testAccVerifyPolicies_NoFilter(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s
resource "pingone_verify_policy" "%[3]s-1" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s-1"
  description    = "%[4]s-1"

  government_id = {
    verify = "REQUIRED"
  }
  depends_on = [pingone_environment.%[2]s]
}

resource "pingone_verify_policy" "%[3]s-2" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s-2"
  description    = "%[4]s-2"

  facial_comparison = {
    verify    = "REQUIRED"
    threshold = "HIGH"
  }
  depends_on = [pingone_environment.%[2]s]
}

resource "pingone_verify_policy" "%[3]s-3" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s-3"
  description    = "%[4]s-3"

  liveness = {
    verify    = "REQUIRED"
    threshold = "HIGH"
  }
  depends_on = [pingone_environment.%[2]s]
}

resource "pingone_verify_policy" "%[3]s-4" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s-4"
  description    = "%[4]s-4"

  liveness = {
    verify    = "REQUIRED"
    threshold = "LOW"
  }
  depends_on = [pingone_environment.%[2]s]
}

data "pingone_verify_policies" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  depends_on = [pingone_verify_policy.%[3]s-1, pingone_verify_policy.%[3]s-2, pingone_verify_policy.%[3]s-3, pingone_verify_policy.%[3]s-4]
}`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}
