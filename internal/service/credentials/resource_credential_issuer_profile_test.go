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

// Note: Issuer Profiles aren't deleted once created [No API]. Deleted only via deletion of the environment.  Placeholder if this changes.
/*func testAccCheckCredentialIssuerProfileDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.CredentialsAPIClient
	ctx = context.WithValue(ctx, credentials.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})

	mgmtApiClient := p1Client.API.ManagementAPIClient
	ctx = context.WithValue(ctx, management.ContextServerVariables, map[string]string{
		"suffix": p1Client.API.Region.URLSuffix,
	})

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_credential_issuer_profile" {
			continue
		}

		_, rEnv, err := mgmtApiClient.EnvironmentsApi.ReadOneEnvironment(ctx, rs.Primary.Attributes["environment_id"]).Execute()

		if err != nil {

			if rEnv == nil {
				return fmt.Errorf("Response object does not exist and no error detected")
			}

			if rEnv.StatusCode == 404 {
				continue
			}

			return err
		}

		body, r, err := apiClient.CredentialIssuersApi.ReadCredentialIssuerProfile(ctx, rs.Primary.Attributes["environment_id"]).Execute()

		if err != nil {

			if r == nil {
				return fmt.Errorf("Response object does not exist and no error detected")
			}

			if r.StatusCode == 404 {
				continue
			}

			tflog.Error(ctx, fmt.Sprintf("Error: %v", body))
			return err
		}

		return fmt.Errorf("PingOne Credential Issuer Profile %s still exists", rs.Primary.ID)
	}

	return nil
}*/

func TestAccCredentialIssuerProfile_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_credential_issuer_profile.%s", resourceName)

	name := acctest.ResourceNameGen()
	updatedName := acctest.ResourceNameGen()

	initialProfile := resource.TestStep{
		Config: testAccCredentialIssuerProfile_Full(resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
			resource.TestMatchResourceAttr(resourceFullName, "application_instance_id", verify.P1ResourceIDRegexp),
			resource.TestCheckResourceAttrSet(resourceFullName, "created_at"),
			resource.TestCheckResourceAttrSet(resourceFullName, "updated_at"),
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
		),
	}

	updatedProfile := resource.TestStep{
		Config: testAccCredentialIssuerProfile_Full(resourceName, updatedName),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
			resource.TestMatchResourceAttr(resourceFullName, "application_instance_id", verify.P1ResourceIDRegexp),
			resource.TestCheckResourceAttrSet(resourceFullName, "created_at"),
			resource.TestCheckResourceAttrSet(resourceFullName, "updated_at"),
			resource.TestCheckResourceAttr(resourceFullName, "name", updatedName),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             nil, // Note: Issuer Profiles aren't deleted once created. Placeholder if this changes.  testAccCheckCredentialIssuerProfileDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// initial profile
			initialProfile,
			{
				Config:  testAccCredentialIssuerProfile_Full(resourceName, name),
				Destroy: true,
			},
			// update profile
			updatedProfile,
			{
				Config:  testAccCredentialIssuerProfile_Full(resourceName, updatedName),
				Destroy: true,
			},
			// changes
			initialProfile,
			updatedProfile,
			initialProfile,
		},
	})
}

func TestAccCredentialIssuerProfile_InvalidConfig(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	name := acctest.ResourceNameGen()

	environmentName := acctest.ResourceNameGenEnvironment()
	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             nil,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccCredentialIssuerProfileInvalidConfig_InvalidName(resourceName, name),
				ExpectError: regexp.MustCompile("Error: Error when calling `CreateCredentialIssuerProfile`: Validation Error : \\[name must not be empty or blank\\]"),
			},
			{
				Config:      testAccCredentialIssuerProfileInvalidConfig_CredentialServuceNotEnabled(environmentName, licenseID, resourceName, name),
				ExpectError: regexp.MustCompile("Error: Error when calling `ReadCredentialIssuerProfile`: Issuer not found for environment"),
			},
		},
	})
}

func testAccCredentialIssuerProfile_Full(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s
resource "pingone_credential_issuer_profile" "%[2]s" {
	environment_id = data.pingone_environment.credentials_test.id
	name = "%[3]s"
}`, acctest.CredentialsSandboxEnvironment(), resourceName, name)
}

func testAccCredentialIssuerProfileInvalidConfig_InvalidName(resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_credential_issuer_profile" "%[3]s" {
	environment_id = data.pingone_environment.credentials_test.id
	name = " "

}`, acctest.CredentialsSandboxEnvironment(), resourceName, name)
}

func testAccCredentialIssuerProfileInvalidConfig_CredentialServuceNotEnabled(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_credential_issuer_profile" "%[3]s" {
	environment_id = pingone_environment.%[2]s.id
	name = "%[3]s"

}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}
