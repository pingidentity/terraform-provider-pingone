package mfa_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccMFAPoliciesDataSource_ByAll(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullName := fmt.Sprintf("data.pingone_mfa_policies.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()

	name := resourceName

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckMFAPolicyDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMFAPoliciesDataSourceConfig_ByAll(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "ids.#", "3"),
					resource.TestMatchResourceAttr(dataSourceFullName, "ids.0", verify.P1ResourceIDRegexpFullString), // the environment default
					resource.TestMatchResourceAttr(dataSourceFullName, "ids.1", verify.P1ResourceIDRegexpFullString), // created by config
					resource.TestMatchResourceAttr(dataSourceFullName, "ids.2", verify.P1ResourceIDRegexpFullString), // created by config
				),
			},
		},
	})
}

func testAccMFAPoliciesDataSourceConfig_ByAll(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_mfa_policy" "%[3]s-1" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s-1"

  sms {
    enabled = true
  }

  voice {
    enabled = true
  }

  email {
    enabled = true
  }

  mobile {
    enabled = true
  }

  totp {
    enabled = true
  }

  fido2 {
    enabled = true
  }

}

resource "pingone_mfa_policy" "%[3]s-2" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s-2"

  sms {
    enabled = true
  }

  voice {
    enabled = true
  }

  email {
    enabled = true
  }

  mobile {
    enabled = true
  }

  totp {
    enabled = true
  }

  fido2 {
    enabled = true
  }

}

data "pingone_mfa_policies" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  depends_on = [
    pingone_mfa_policy.%[3]s-1,
    pingone_mfa_policy.%[3]s-2,
  ]
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}
