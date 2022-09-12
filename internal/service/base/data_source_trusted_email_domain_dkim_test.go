package base_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
)

func TestAccTrustedEmailDomainDKIMDataSource_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_trusted_email_domain_dkim.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheckEnvironment(t) },
		ProviderFactories: acctest.ProviderFactories,
		ErrorCheck:        acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccTrustedEmailDomainDKIMDataSourceConfig_Full(environmentName, licenseID, resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceFullName, "type", "CNAME"),
					resource.TestCheckResourceAttr(dataSourceFullName, "region.#", "2"),
					resource.TestMatchResourceAttr(dataSourceFullName, "region.0.name", regexp.MustCompile(`^[a-z0-9-]*$`)),
					resource.TestCheckResourceAttr(dataSourceFullName, "region.0.status", "VERIFICATION_REQUIRED"),
					resource.TestCheckResourceAttr(dataSourceFullName, "region.0.token.#", "3"),
					resource.TestMatchResourceAttr(dataSourceFullName, "region.0.token.0.key", regexp.MustCompile(`^[a-zA-Z0-9_\.]*$`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "region.0.token.0.value", regexp.MustCompile(`^[a-z0-9_\.]*$`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "region.0.token.1.key", regexp.MustCompile(`^[a-zA-Z0-9_\.]*$`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "region.0.token.1.value", regexp.MustCompile(`^[a-z0-9_\.]*$`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "region.0.token.2.key", regexp.MustCompile(`^[a-zA-Z0-9_\.]*$`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "region.0.token.2.value", regexp.MustCompile(`^[a-z0-9_\.]*$`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "region.1.name", regexp.MustCompile(`^[a-z0-9-]*$`)),
					resource.TestCheckResourceAttr(dataSourceFullName, "region.1.status", "VERIFICATION_REQUIRED"),
					resource.TestCheckResourceAttr(dataSourceFullName, "region.1.token.#", "3"),
					resource.TestMatchResourceAttr(dataSourceFullName, "region.1.token.0.key", regexp.MustCompile(`^[a-zA-Z0-9_\.]*$`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "region.1.token.0.value", regexp.MustCompile(`^[a-z0-9_\.]*$`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "region.1.token.1.key", regexp.MustCompile(`^[a-zA-Z0-9_\.]*$`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "region.1.token.1.value", regexp.MustCompile(`^[a-z0-9_\.]*$`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "region.1.token.2.key", regexp.MustCompile(`^[a-zA-Z0-9_\.]*$`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "region.1.token.2.value", regexp.MustCompile(`^[a-z0-9_\.]*$`)),
				),
			},
		},
	})
}

func testAccTrustedEmailDomainDKIMDataSourceConfig_Full(environmentName, licenseID, resourceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_trusted_email_domain" "%[3]s" {
	environment_id = pingone_environment.%[2]s.id
  
	domain_name = "terraformdev.ping-eng.com"
  }

data "pingone_trusted_email_domain_dkim" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  trusted_email_domain_id = pingone_trusted_email_domain.%[3]s.id
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName)
}
