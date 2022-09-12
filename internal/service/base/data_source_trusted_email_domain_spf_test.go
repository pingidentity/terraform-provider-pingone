package base_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
)

func TestAccTrustedEmailDomainSPFDataSource_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_trusted_email_domain_spf.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccTrustedEmailDomainSPFDataSourceConfig_Full(environmentName, licenseID, resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceFullName, "type", "TXT"),
					resource.TestCheckResourceAttr(dataSourceFullName, "status", "VERIFICATION_REQUIRED"),
					resource.TestCheckResourceAttr(dataSourceFullName, "key", "@"),
					resource.TestCheckResourceAttr(dataSourceFullName, "value", "v=spf1 include:amazonses.com ~all"),
				),
			},
		},
	})
}

func testAccTrustedEmailDomainSPFDataSourceConfig_Full(environmentName, licenseID, resourceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_trusted_email_domain" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  domain_name = "terraformdev.ping-eng.com"
}

data "pingone_trusted_email_domain_spf" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  trusted_email_domain_id = pingone_trusted_email_domain.%[3]s.id
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName)
}
