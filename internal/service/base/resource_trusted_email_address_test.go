package base_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func testAccCheckTrustedEmailAddressDestroy(s *terraform.State) error {
	var ctx = context.Background()

	p1Client, err := acctest.TestClient(ctx)

	if err != nil {
		return err
	}

	apiClient := p1Client.API.ManagementAPIClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingone_trusted_email_address" {
			continue
		}

		_, rEnv, err := apiClient.EnvironmentsApi.ReadOneEnvironment(ctx, rs.Primary.Attributes["environment_id"]).Execute()

		if err != nil {

			if rEnv == nil {
				return fmt.Errorf("Response object does not exist and no error detected")
			}

			if rEnv.StatusCode == 404 {
				continue
			}

			return err
		}

		body, r, err := apiClient.TrustedEmailAddressesApi.ReadOneTrustedEmailAddress(ctx, rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["email_domain_id"], rs.Primary.ID).Execute()

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

		return fmt.Errorf("PingOne trusted email address %s still exists", rs.Primary.ID)
	}

	return nil
}

func TestAccTrustedEmailAddress_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_trusted_email_address.%s", resourceName)

	verifiedDomain := os.Getenv("PINGONE_VERIFIED_EMAIL_DOMAIN")
	emailAddress := fmt.Sprintf("%s@%s", resourceName, verifiedDomain)

	check := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexp),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexp),
		resource.TestMatchResourceAttr(resourceFullName, "email_domain_id", verify.P1ResourceIDRegexp),
		resource.TestCheckResourceAttr(resourceFullName, "email_address", emailAddress),
		resource.TestCheckResourceAttr(resourceFullName, "status", "VERIFICATION_REQUIRED"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironmentDomainVerified(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckTrustedEmailAddressDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: testAccTrustedEmailAddressConfig_New_DomainVerified(resourceName, verifiedDomain, emailAddress),
				Check:  check,
			},
			{
				Config:  testAccTrustedEmailAddressConfig_New_DomainVerified(resourceName, verifiedDomain, emailAddress),
				Destroy: true,
			},
			{
				Config: testAccTrustedEmailAddressConfig_New_DomainVerified(resourceName, verifiedDomain, emailAddress),
				Check:  check,
			},
		},
	})
}

func TestAccTrustedEmailAddress_NotVerified(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()

	environmentName := acctest.ResourceNameGenEnvironment()
	licenseID := os.Getenv("PINGONE_LICENSE_ID")

	unverifiedDomain := "terraformdev.ping-eng.com"
	unverifiedEmailAddress := fmt.Sprintf("noreply@%s", unverifiedDomain)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheckEnvironment(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckTrustedEmailAddressDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccTrustedEmailAddressConfig_New_Full(environmentName, licenseID, resourceName, unverifiedDomain, unverifiedEmailAddress),
				ExpectError: regexp.MustCompile(`The domain of the given email address is not verified`),
			},
		},
	})
}

func testAccTrustedEmailAddressConfig_New_Full(environmentName, licenseID, resourceName, domain, emailAddress string) string {
	return fmt.Sprintf(`
	%[1]s

resource "pingone_trusted_email_domain" "%[3]s" {
  environment_id = pingone_environment.%[2]s.id

  domain_name = "%[4]s"
}

resource "pingone_trusted_email_address" "%[3]s" {
  environment_id  = pingone_environment.%[2]s.id
  email_domain_id = pingone_trusted_email_domain.%[3]s.id

  email_address = "%[5]s"
}`, acctest.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, domain, emailAddress)
}

func testAccTrustedEmailAddressConfig_New_DomainVerified(resourceName, verifiedDomain, emailAddress string) string {
	return fmt.Sprintf(`
	%[1]s

data "pingone_trusted_email_domain" "%[2]s" {
  environment_id = data.pingone_environment.domainverified_test.id

  domain_name = "%[3]s"
}

resource "pingone_trusted_email_address" "%[2]s" {
  environment_id  = data.pingone_environment.domainverified_test.id
  email_domain_id = data.pingone_trusted_email_domain.%[2]s.id

  email_address = "%[4]s"
}`, acctest.DomainVerifiedSandboxEnvironment(), resourceName, verifiedDomain, emailAddress)
}
