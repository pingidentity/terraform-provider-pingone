// Copyright © 2025 Ping Identity Corporation

package base_test

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
)

func TestAccTrustedEmailDomainOwnershipDataSource_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	domainPrefix := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_trusted_email_domain_ownership.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

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
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccTrustedEmailDomainOwnershipDataSourceConfig_Full(environmentName, licenseID, resourceName, domainPrefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceFullName, "type", "TXT"),
					resource.TestCheckResourceAttrWith(dataSourceFullName, "regions.#", validateRegionCardinality),
					resource.TestMatchResourceAttr(dataSourceFullName, "regions.0.name", regexp.MustCompile(`^[a-z0-9-]*$`)),
					resource.TestCheckResourceAttr(dataSourceFullName, "regions.0.status", "VERIFICATION_REQUIRED"),
					resource.TestCheckResourceAttr(dataSourceFullName, "regions.0.key", "_amazonses"),
					resource.TestMatchResourceAttr(dataSourceFullName, "regions.0.value", regexp.MustCompile(`^[a-zA-Z0-9-~_//+=]*$`)),
					resource.TestCheckResourceAttr(dataSourceFullName, "environment_dns_record.status", "VERIFICATION_REQUIRED"),
					resource.TestCheckResourceAttr(dataSourceFullName, "environment_dns_record.key", "_pingoneemail"),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_dns_record.value", regexp.MustCompile(`^[a-zA-Z0-9-~_//+=]*$`)),
				),
			},
		},
	})
}

func TestAccTrustedEmailDomainOwnershipDataSource_NotFound(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// {
			// 	Config:      testAccTrustedEmailDomainOwnershipDataSourceConfig_NotFoundByName(resourceName),
			// 	ExpectError: regexp.MustCompile("Cannot find domain doesnotexist"),
			// },
			{
				Config:      testAccTrustedEmailDomainOwnershipDataSourceConfig_NotFoundByID(resourceName),
				ExpectError: regexp.MustCompile("Error when calling `ReadTrustedEmailDomainOwnershipStatus`: The request could not be completed. The requested resource was not found."),
			},
		},
	})
}

func testAccTrustedEmailDomainOwnershipDataSourceConfig_Full(environmentName, licenseID, resourceName, domainPrefix string) string {
	return fmt.Sprintf(`
%[1]s

resource "pingone_trusted_email_domain" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  domain_name = "%[4]s.terraformdev.ping-eng.com"
}

data "pingone_trusted_email_domain_ownership" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  trusted_email_domain_id = pingone_trusted_email_domain.%[3]s.id
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, domainPrefix)
}

func testAccTrustedEmailDomainOwnershipDataSourceConfig_NotFoundByID(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

data "pingone_trusted_email_domain_ownership" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id

  trusted_email_domain_id = "9c052a8a-14be-44e4-8f07-2662569994ce" // dummy ID that conforms to UUID v4
}`, acctest.GenericSandboxEnvironment(), resourceName)
}

func validateRegionCardinality(value string) error {

	valueInt, err := strconv.Atoi(value)
	if err != nil {
		return err
	}

	if valueInt < 1 {
		return fmt.Errorf("region block should have at least one set of values")
	}
	return nil

}
