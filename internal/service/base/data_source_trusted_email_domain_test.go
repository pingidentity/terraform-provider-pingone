package base_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
)

func TestAccTrustedEmailDomainDataSource_ByNameFull(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_trusted_email_domain.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	verifiedDomain := os.Getenv("PINGONE_VERIFIED_EMAIL_DOMAIN")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironmentDomainVerified(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             nil, // The test environment is static and no resources are created, nothing to check on destroy
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccTrustedEmailDomainDataSourceConfig_ByNameFull(resourceName, verifiedDomain),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "trusted_email_domain_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(dataSourceFullName, "domain_name", verifiedDomain),
				),
			},
		},
	})
}

func TestAccTrustedEmailDomainDataSource_ByIDFull(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_trusted_email_domain.%s", resourceName)
	dataSourceFullName := fmt.Sprintf("data.%s", resourceFullName)

	verifiedDomain := os.Getenv("PINGONE_VERIFIED_EMAIL_DOMAIN")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             nil, // The test environment is static and no resources are created, nothing to check on destroy
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccTrustedEmailDomainDataSourceConfig_ByIDFull(resourceName, verifiedDomain),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestMatchResourceAttr(dataSourceFullName, "trusted_email_domain_id", regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)),
					resource.TestCheckResourceAttr(dataSourceFullName, "domain_name", verifiedDomain),
				),
			},
		},
	})
}

func TestAccTrustedEmailDomainDataSource_NotFound(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTrustedEmailDomainDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccTrustedEmailDomainDataSourceConfig_NotFoundByName(resourceName),
				ExpectError: regexp.MustCompile(`Cannot find trusted email domain from domain_name`),
			},
			{
				Config:      testAccTrustedEmailDomainDataSourceConfig_NotFoundByID(resourceName),
				ExpectError: regexp.MustCompile("Error when calling `ReadOneTrustedEmailDomain`: The request could not be completed. The requested resource was not found."),
			},
		},
	})
}

func testAccTrustedEmailDomainDataSourceConfig_ByNameFull(resourceName, verifiedDomain string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_trusted_email_domain" "%[2]s" {
  environment_id = data.pingone_environment.domainverified_test.id

  domain_name = "%[3]s"
}
`, acctest.DomainVerifiedSandboxEnvironment(), resourceName, verifiedDomain)
}

func testAccTrustedEmailDomainDataSourceConfig_ByIDFull(resourceName, verifiedDomain string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_trusted_email_domain" "%[2]s-name" {
  environment_id = data.pingone_environment.domainverified_test.id

  domain_name = "%[3]s"
}

data "pingone_trusted_email_domain" "%[2]s" {
  environment_id = data.pingone_environment.domainverified_test.id

  trusted_email_domain_id = data.pingone_trusted_email_domain.%[2]s-name.id
}`, acctest.DomainVerifiedSandboxEnvironment(), resourceName, verifiedDomain)
}

func testAccTrustedEmailDomainDataSourceConfig_NotFoundByName(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_trusted_email_domain" "%[2]s" {
  environment_id = data.pingone_environment.domainverified_test.id

  domain_name = "doesnotexist.com"
}
`, acctest.DomainVerifiedSandboxEnvironment(), resourceName)
}

func testAccTrustedEmailDomainDataSourceConfig_NotFoundByID(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_trusted_email_domain" "%[2]s" {
  environment_id = data.pingone_environment.domainverified_test.id

  trusted_email_domain_id = "9c052a8a-14be-44e4-8f07-2662569994ce" // dummy ID that conforms to UUID v4
}
`, acctest.DomainVerifiedSandboxEnvironment(), resourceName)
}
