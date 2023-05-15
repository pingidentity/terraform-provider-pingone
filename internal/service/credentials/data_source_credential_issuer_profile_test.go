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

	name := acctest.ResourceNameGen()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             nil, // Note: Issuer Profiles aren't deleted once created. Placeholder if this changes.  testAccCheckCredentialIssuerProfileDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCredentialIssuerProfileDataSource_ByEnvironmentIDFull(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(dataSourceFullName, "application_instance_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "name"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "created_at"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "updated_at"),
				),
			},
			{
				Config:  testAccCredentialIssuerProfileDataSource_ByEnvironmentIDFull(resourceName, name),
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
	name := acctest.ResourceNameGen()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             nil, // Note: Issuer Profiles aren't deleted once created. Placeholder if this changes.  testAccCheckCredentialIssuerProfileDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccCredentialIssuerProfileDataSource_NotFound(environmentName, licenseID, resourceName, name),
				ExpectError: regexp.MustCompile("Error when calling `ReadCredentialIssuerProfile`: Issuer not found for environment"),
			},
		},
	})
}

func testAccCredentialIssuerProfileDataSource_ByEnvironmentIDFull(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_credential_issuer_profile" "%[2]s" {
  environment_id = data.pingone_environment.credentials_test.id
}`, acctest.CredentialsSandboxEnvironment(), resourceName, name)
}

func testAccCredentialIssuerProfileDataSource_NotFound(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_credential_issuer_profile" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
}
`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}
