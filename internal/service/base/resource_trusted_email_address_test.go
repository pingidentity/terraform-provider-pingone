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

func testAccGetTrustedEmailAddressIDs(resourceName string, environmentID, emailDomainID, resourceID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource Not found: %s", resourceName)
		}

		*resourceID = rs.Primary.ID
		*emailDomainID = rs.Primary.Attributes["email_domain_id"]
		*environmentID = rs.Primary.Attributes["environment_id"]

		return nil
	}
}

func TestAccTrustedEmailAddress_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_trusted_email_address.%s", resourceName)

	verifiedDomain := os.Getenv("PINGONE_VERIFIED_EMAIL_DOMAIN")
	emailAddress := fmt.Sprintf("%s@%s", resourceName, verifiedDomain)

	var resourceID, emailDomainID, environmentID string

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckTrustedEmailAddressDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Test removal of the resource
			{
				Config: testAccTrustedEmailAddressConfig_New_DomainVerified(resourceName, verifiedDomain, emailAddress),
				Check:  testAccGetTrustedEmailAddressIDs(resourceFullName, &environmentID, &emailDomainID, &resourceID),
			},
			// Replan after removal preconfig
			{
				PreConfig: func() {
					var ctx = context.Background()
					p1Client, err := acctest.TestClient(ctx)

					if err != nil {
						t.Fatalf("Failed to get API client: %v", err)
					}

					apiClient := p1Client.API.ManagementAPIClient

					if environmentID == "" || emailDomainID == "" || resourceID == "" {
						t.Fatalf("One of environment ID, email domain ID or resource ID cannot be determined. Environment ID: %s, Email Domain ID: %s, Resource ID: %s", environmentID, emailDomainID, resourceID)
					}

					_, err = apiClient.TrustedEmailAddressesApi.DeleteTrustedEmailAddress(ctx, environmentID, emailDomainID, resourceID).Execute()
					if err != nil {
						t.Fatalf("Failed to delete trusted email address: %v", err)
					}
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// Test removal of the email domain
			{
				Config: testAccTrustedEmailAddressConfig_New_DomainVerified(resourceName, verifiedDomain, emailAddress),
				Check:  testAccGetTrustedEmailAddressIDs(resourceFullName, &environmentID, &emailDomainID, &resourceID),
			},
			// Replan after removal preconfig
			{
				PreConfig: func() {
					var ctx = context.Background()
					p1Client, err := acctest.TestClient(ctx)

					if err != nil {
						t.Fatalf("Failed to get API client: %v", err)
					}

					apiClient := p1Client.API.ManagementAPIClient

					if environmentID == "" || emailDomainID == "" || resourceID == "" {
						t.Fatalf("One of environment ID, email domain ID or resource ID cannot be determined. Environment ID: %s, Email Domain ID: %s, Resource ID: %s", environmentID, emailDomainID, resourceID)
					}

					_, err = apiClient.TrustedEmailDomainsApi.DeleteTrustedEmailDomain(ctx, environmentID, emailDomainID).Execute()
					if err != nil {
						t.Fatalf("Failed to delete trusted email domain: %v", err)
					}
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccTrustedEmailAddress_Full(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_trusted_email_address.%s", resourceName)

	verifiedDomain := os.Getenv("PINGONE_VERIFIED_EMAIL_DOMAIN")
	emailAddress := fmt.Sprintf("%s@%s", resourceName, verifiedDomain)

	check := resource.ComposeTestCheckFunc(
		resource.TestMatchResourceAttr(resourceFullName, "id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "environment_id", verify.P1ResourceIDRegexpFullString),
		resource.TestMatchResourceAttr(resourceFullName, "email_domain_id", verify.P1ResourceIDRegexpFullString),
		resource.TestCheckResourceAttr(resourceFullName, "email_address", emailAddress),
		resource.TestCheckResourceAttr(resourceFullName, "status", "VERIFICATION_REQUIRED"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckDomainVerification(t)
		},
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
			// Test importing the resource
			{
				ResourceName: resourceFullName,
				ImportStateIdFunc: func() resource.ImportStateIdFunc {
					return func(s *terraform.State) (string, error) {
						rs, ok := s.RootModule().Resources[resourceFullName]
						if !ok {
							return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
						}

						return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["environment_id"], rs.Primary.Attributes["email_domain_id"], rs.Primary.ID), nil
					}
				}(),
				ImportState:       true,
				ImportStateVerify: true,
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
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
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

func TestAccTrustedEmailAddress_BadParameters(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_trusted_email_address.%s", resourceName)

	verifiedDomain := os.Getenv("PINGONE_VERIFIED_EMAIL_DOMAIN")
	emailAddress := fmt.Sprintf("%s@%s", resourceName, verifiedDomain)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckTrustedEmailAddressDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Configure
			{
				Config: testAccTrustedEmailAddressConfig_New_DomainVerified(resourceName, verifiedDomain, emailAddress),
			},
			// Errors
			{
				ResourceName: resourceFullName,
				ImportState:  true,
				ExpectError:  regexp.MustCompile(`Invalid import ID specified \(".*"\).  The ID should be in the format "environment_id/email_domain_id/trusted_email_address_id" and must match regex: .*`),
			},
			{
				ResourceName:  resourceFullName,
				ImportStateId: "/",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`Invalid import ID specified \(".*"\).  The ID should be in the format "environment_id/email_domain_id/trusted_email_address_id" and must match regex: .*`),
			},
			{
				ResourceName:  resourceFullName,
				ImportStateId: "badformat/badformat/badformat",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`Invalid import ID specified \(".*"\).  The ID should be in the format "environment_id/email_domain_id/trusted_email_address_id" and must match regex: .*`),
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
