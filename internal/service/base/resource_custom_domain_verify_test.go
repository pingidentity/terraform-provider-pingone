// Copyright © 2025 Ping Identity Corporation

package base_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/base"
)

func TestAccCustomDomainVerify_CannotVerifyNXDOMAIN(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	domainPrefix := acctest.ResourceNameGen()

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.CustomDomainVerify_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccCustomDomainVerifyConfig_CannotVerifyNXDOMAIN(environmentName, licenseID, resourceName, domainPrefix),
				ExpectError: regexp.MustCompile(`Cannot verify the domain`),
			},
		},
	})
}

func testAccCustomDomainVerifyConfig_CannotVerifyNXDOMAIN(environmentName, licenseID, resourceName, domainPrefix string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_custom_domain" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  domain_name = "%[4]s.terraformdev-verify.ping-eng.com"
}

resource "pingone_custom_domain_verify" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  custom_domain_id = pingone_custom_domain.%[3]s.id

  timeouts = {
    create = "5s"
  }
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, domainPrefix)
}
