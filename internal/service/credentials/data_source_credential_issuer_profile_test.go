package credentials_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/credentials"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccCredentialIssuerProfileDataSource_ByEnvironmentIDFull(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	dataSourceFullName := fmt.Sprintf("data.pingone_credential_issuer_profile.%s", resourceName)

	environmentName := acctest.ResourceNameGenEnvironment()
	licenseID := os.Getenv("PINGONE_LICENSE_ID")
	name := acctest.ResourceNameGen()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             credentials.CredentialIssuerProfile_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCredentialIssuerProfileDataSource_ByEnvironmentIDFull(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
					resource.TestMatchResourceAttr(dataSourceFullName, "application_instance_id", verify.P1ResourceIDRegexpFullString),
					resource.TestCheckResourceAttr(dataSourceFullName, "name", name),
					resource.TestMatchResourceAttr(dataSourceFullName, "created_at", verify.RFC3339Regexp),
					resource.TestMatchResourceAttr(dataSourceFullName, "updated_at", verify.RFC3339Regexp),
				),
			},
			{
				Config:  testAccCredentialIssuerProfileDataSource_ByEnvironmentIDFull(environmentName, licenseID, resourceName, name),
				Destroy: true,
			},
		},
	})
}

func TestAccCredentialIssuerProfileDataSource_NotFound(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	environmentName := acctest.ResourceNameGenEnvironment()

	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             credentials.CredentialIssuerProfile_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccCredentialIssuerProfileDataSource_NotFound(environmentName, licenseID, resourceName),
				ExpectError: regexp.MustCompile("Error when calling `ReadCredentialIssuerProfile`: Issuer not found for environment"),
			},
		},
	})
}

func testAccCredentialIssuerProfileDataSource_ByEnvironmentIDFull(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_credential_issuer_profile" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s"

  depends_on = [pingone_environment.%[2]s]
}

data "pingone_credential_issuer_profile" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  depends_on = [pingone_environment.%[2]s, resource.pingone_credential_issuer_profile.%[3]s]

}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccCredentialIssuerProfileDataSource_NotFound(environmentName, licenseID, resourceName string) string {
	return fmt.Sprintf(`
resource "pingone_environment" "%[3]s" {
  name       = "%[1]s"
  type       = "SANDBOX"
  license_id = "%[2]s"

  services = [
    {
      type = "SSO"
    },
    {
      type = "MFA"
    },
    {
      type = "Risk"
    }
  ]
}

data "pingone_credential_issuer_profile" "%[3]s" {
  environment_id = pingone_environment.%[3]s.id
}`, environmentName, licenseID, resourceName)
}
