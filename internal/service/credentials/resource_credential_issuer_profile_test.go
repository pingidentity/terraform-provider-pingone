package credentials_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

// Note: Issuer Profiles aren't deleted once created [No API]. Deleted only via deletion of the environment.
// Destroy is a placeholder if this changes. Defined a passthrough for linter purposes.
func testAccCheckCredentialIssuerProfilePassthrough(s *terraform.State) error {

	return nil
}

/*func testAccCheckCredentialIssuerProfileDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.CredentialsAPIClient


	mgmtApiClient := p1Client.API.ManagementAPIClient


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

func TestAccCredentialIssuerProfile_Import(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_credential_issuer_profile.%s", resourceName)

	name := resourceName

	environmentName := acctest.ResourceNameGenEnvironment()
	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCredentialIssuerProfilePassthrough,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccCredentialIssuerProfile_Full(environmentName, licenseID, resourceName, name),
			},
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Errors
			{
				ResourceName: resourceFullName,
				ImportState:  true,
				ExpectError:  regexp.MustCompile(`Invalid import ID specified \(".*"\).  The ID should be in the format "environment_id/credential_issuer_profile_id".`),
			},
			{
				ResourceName:  resourceFullName,
				ImportStateId: "/",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`Invalid import ID specified \(".*"\).  The ID should be in the format "environment_id/credential_issuer_profile_id".`),
			},
			{
				ResourceName:  resourceFullName,
				ImportStateId: "badformat/badformat",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`Invalid import ID specified \(".*"\).  The ID should be in the format "environment_id/credential_issuer_profile_id".`),
			},
		},
	})
}

func TestAccCredentialIssuerProfile_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_credential_issuer_profile.%s", resourceName)

	name := acctest.ResourceNameGen()
	updatedName := acctest.ResourceNameGen()

	environmentName := acctest.ResourceNameGenEnvironment()
	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	initialProfile := resource.TestStep{
		Config: testAccCredentialIssuerProfile_Full(environmentName, licenseID, resourceName, name),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
			resource.TestMatchResourceAttr(resourceFullName, "application_instance_id", verify.P1ResourceIDRegexp),
			resource.TestCheckResourceAttr(resourceFullName, "name", name),
			resource.TestMatchResourceAttr(resourceFullName, "created_at", verify.RFC3339Regexp),
			resource.TestMatchResourceAttr(resourceFullName, "updated_at", verify.RFC3339Regexp),
		),
	}

	updatedProfile := resource.TestStep{
		Config: testAccCredentialIssuerProfile_Full(environmentName, licenseID, resourceName, updatedName),
		Check: resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
			resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
			resource.TestMatchResourceAttr(resourceFullName, "application_instance_id", verify.P1ResourceIDRegexp),
			resource.TestCheckResourceAttr(resourceFullName, "name", updatedName),
			resource.TestMatchResourceAttr(resourceFullName, "created_at", verify.RFC3339Regexp),
			resource.TestMatchResourceAttr(resourceFullName, "updated_at", verify.RFC3339Regexp),
		),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCredentialIssuerProfilePassthrough,
		//CheckDestroy:           testAccCheckCredentialIssuerProfileDestroy  // Note: Issuer Profiles aren't deleted once created. Uncomment and replace Passthrough if this changes.
		ErrorCheck: acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// initial profile
			initialProfile,
			{
				Config:  testAccCredentialIssuerProfile_Full(environmentName, licenseID, resourceName, name),
				Destroy: true,
			},
			// update profile
			updatedProfile,
			{
				Config:  testAccCredentialIssuerProfile_Full(environmentName, licenseID, resourceName, updatedName),
				Destroy: true,
			},
			// changes
			initialProfile,
			updatedProfile,
			initialProfile,
			{
				Config:  testAccCredentialIssuerProfile_Full(environmentName, licenseID, resourceName, name),
				Destroy: true,
			},
			{
				Config:  testAccCredentialIssuerProfile_Full(environmentName, licenseID, resourceName, updatedName),
				Destroy: true,
			},
		},
	})
}

func TestAccCredentialIssuerProfile_InvalidConfig(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	//name := acctest.ResourceNameGen()

	environmentName := acctest.ResourceNameGenEnvironment()
	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             nil, //testAccCheckCredentialIssuerProfilePassthrough,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccCredentialIssuerProfileInvalidConfig_InvalidName(environmentName, licenseID, resourceName, ""),
				ExpectError: regexp.MustCompile("Error: Invalid Attribute Value Length"),
			},
		},
	})
}

func testAccCredentialIssuerProfile_Full(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s
resource "pingone_credential_issuer_profile" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s"

  depends_on = [pingone_environment.%[2]s]

}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}

func testAccCredentialIssuerProfileInvalidConfig_InvalidName(environmentName, licenseID, resourceName, name string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_credential_issuer_profile" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id
  name           = "%[4]s"

  depends_on = [pingone_environment.%[2]s]

}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, name)
}
