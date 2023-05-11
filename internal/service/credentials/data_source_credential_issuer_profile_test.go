package credentials_test

import (
	"fmt"
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

	name := resourceName

	// test i'd prefer - but needs work
	//organizationName := "internal_mikesimontf_444489364" // i need to get the org name dynamically, not a fixed value - how within current framework

	// preference is to do some date comparisons, but limited options other than comparing yyyy-mm-dd components of dates, which has limited utility
	//date := time.Now().Format(time.RFC3339)
	//createdAt := date
	//updatedAt := date

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             nil,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCredentialIssuerProfileDataSourceConfigDataSource_ByEnvironmentIDFull(resourceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceFullName, "id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(dataSourceFullName, "environment_id", verify.P1ResourceIDRegexp),
					resource.TestMatchResourceAttr(dataSourceFullName, "application_instance_id", verify.P1ResourceIDRegexp),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "name"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "created_at"),
					resource.TestCheckResourceAttrSet(dataSourceFullName, "updated_at"),
					// future improvement?
					//resource.TestCheckResourceAttr(dataSourceFullName, "name", organizationName),
					//resource.TestCheckResourceAttr(dataSourceFullName, "created_at", createdAt),
					//resource.TestCheckResourceAttr(dataSourceFullName, "updated_at", updatedAt),
				),
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
		CheckDestroy:             nil,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccCredentialIssuerProfileDataSourceConfigDataSource_NotFound(resourceName),
				ExpectError: regexp.MustCompile("Error when calling `ReadCredentialIssuerProfile`: Issuer not found for environment"),
			},
		},
	})
}

func testAccCredentialIssuerProfileDataSourceConfigDataSource_ByEnvironmentIDFull(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_credential_issuer_profile" "%[2]s" {
	environment_id = data.pingone_environment.credentials_test.id
}`, acctest.CredentialsSandboxEnvironment(), resourceName, name)
}

func testAccCredentialIssuerProfileDataSourceConfigDataSource_NotFound(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_credential_issuer_profile" "%[2]s" {
	environment_id = data.pingone_environment.general_test.id // generic environmet doesn't have P1Creds configured
}
`, acctest.GenericSandboxEnvironment(), resourceName)
}
