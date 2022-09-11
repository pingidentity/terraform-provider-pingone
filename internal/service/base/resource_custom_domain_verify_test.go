package base_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
)

func TestAccCustomDomainVerify_CannotVerifyNXDOMAIN(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckCustomDomainDestroy,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccCustomDomainVerifyConfig_CannotVerifyNXDOMAIN(environmentName, licenseID, resourceName),
				ExpectError: regexp.MustCompile(`Cannot verify the domain`),
			},
		},
	})
}

func testAccCustomDomainVerifyConfig_CannotVerifyNXDOMAIN(environmentName, licenseID, resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_custom_domain" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  domain_name = "terraformdev-verify.ping-eng.com"
}

resource "pingone_custom_domain_verify" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  custom_domain_id = pingone_custom_domain.%[3]s.id

  timeouts {
    create = "5s"
  }
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName)
}
