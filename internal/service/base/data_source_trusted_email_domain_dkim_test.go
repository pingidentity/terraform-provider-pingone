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
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccTrustedEmailDomainDKIMDataSourceConfig_Full(environmentName, licenseID, resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceFullName, "type", "CNAME"),
					resource.TestCheckResourceAttrWith(dataSourceFullName, "region.#", validateRegionCardinality),
					resource.TestMatchResourceAttr(dataSourceFullName, "region.0.name", regexp.MustCompile(`^[a-z0-9-]*$`)),
					resource.TestCheckResourceAttr(dataSourceFullName, "region.0.status", "VERIFICATION_REQUIRED"),
					resource.TestCheckResourceAttrWith(dataSourceFullName, "region.0.token.#", validateRegionCardinality),
					resource.TestMatchResourceAttr(dataSourceFullName, "region.0.token.0.key", regexp.MustCompile(`^[a-zA-Z0-9_\.]*$`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "region.0.token.0.value", regexp.MustCompile(`^[a-z0-9_\.]*$`)),
				),
			},
		},
	})
}

func TestAccTrustedEmailDomainDKIMDataSource_NotFound(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// {
			// 	Config:      testAccTrustedEmailDomainDKIMDataSourceConfig_NotFoundByName(resourceName),
			// 	ExpectError: regexp.MustCompile("Cannot find domain doesnotexist"),
			// },
			{
				Config:      testAccTrustedEmailDomainDKIMDataSourceConfig_NotFoundByID(resourceName),
				ExpectError: regexp.MustCompile("Error when calling `ReadTrustedEmailDomainDKIMStatus`: The request could not be completed. The requested resource was not found."),
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

func testAccTrustedEmailDomainDKIMDataSourceConfig_NotFoundByID(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

data "pingone_trusted_email_domain_dkim" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  trusted_email_domain_id = "9c052a8a-14be-44e4-8f07-2662569994ce" // dummy ID that conforms to UUID v4
}`, acctest.GenericSandboxEnvironment(), resourceName)
}
