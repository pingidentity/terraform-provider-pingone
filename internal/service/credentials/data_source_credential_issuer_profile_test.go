package credentials_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
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
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCredentialIssuerProfilePassthrough,
		//CheckDestroy:           testAccCheckCredentialIssuerProfileDestroy  // Note: Issuer Profiles aren't deleted once created. Uncomment and replace Passthrough if this changes.
		ErrorCheck: acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCredentialIssuerProfileDataSource_ByEnvironmentIDFull(environmentName, licenseID, resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(dataSourceFullName, "application_instance_id", verify.P1ResourceIDRegexp),
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

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCredentialIssuerProfilePassthrough,
		//CheckDestroy:           testAccCheckCredentialIssuerProfileDestroy  // Note: Issuer Profiles aren't deleted once created. Uncomment and replace Passthrough if this changes.
		ErrorCheck: acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccCredentialIssuerProfileDataSource_NotFound(resourceName),
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

func testAccCredentialIssuerProfileDataSource_NotFound(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_credential_issuer_profile" "%[2]s" {
  environment_id = data.pingone_environment.general_test.id
}`, acctest.GenericSandboxEnvironment(), resourceName)
}
