// Copyright © 2025 Ping Identity Corporation

package base_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ErrorCheck:               acctest.ErrorCheck(t),
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

func TestAccTrustedEmailDomainSPFDataSource_NotFound(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// {
			// 	Config:      testAccTrustedEmailDomainSPFDataSourceConfig_NotFoundByName(resourceName),
			// 	ExpectError: regexp.MustCompile("Cannot find domain doesnotexist"),
			// },
			{
				Config:      testAccTrustedEmailDomainSPFDataSourceConfig_NotFoundByID(resourceName),
				ExpectError: regexp.MustCompile("Error when calling `ReadTrustedEmailDomainSPFStatus`: The request could not be completed. The requested resource was not found."),
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

func testAccTrustedEmailDomainSPFDataSourceConfig_NotFoundByID(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

data "pingone_trusted_email_domain_spf" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  trusted_email_domain_id = "9c052a8a-14be-44e4-8f07-2662569994ce" // dummy ID that conforms to UUID v4
}`, acctest.GenericSandboxEnvironment(), resourceName)
}
