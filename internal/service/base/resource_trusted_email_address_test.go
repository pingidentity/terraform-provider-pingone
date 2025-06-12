// Copyright © 2025 Ping Identity Corporation

package base_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest"
	acctestlegacysdk "github.com/pingidentity/terraform-provider-pingone/internal/acctest/legacysdk"
	"github.com/pingidentity/terraform-provider-pingone/internal/acctest/service/base"
	client "github.com/pingidentity/terraform-provider-pingone/internal/client"
	"github.com/pingidentity/terraform-provider-pingone/internal/verify"
)

func TestAccTrustedEmailAddress_RemovalDrift(t *testing.T) {
	t.Parallel()

	resourceName := acctest.ResourceNameGen()
	resourceFullName := fmt.Sprintf("pingone_trusted_email_address.%s", resourceName)

	verifiedDomain := os.Getenv("PINGONE_VERIFIED_EMAIL_DOMAIN")
	emailAddress := fmt.Sprintf("%s@%s", resourceName, verifiedDomain)

	var trustedEmailAddressID, emailDomainID, environmentID string

	var p1Client *client.Client
	var ctx = context.Background()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
			acctest.PreCheckDomainVerification(t)

			p1Client = acctestlegacysdk.PreCheckTestClient(ctx, t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.TrustedEmailAddress_CheckDestroy,
		ErrorCheck:               acctest.ErrorCheck(t),
		Steps: []resource.TestStep{
			// Test removal of the resource
			{
				Config: testAccTrustedEmailAddressConfig_New_DomainVerified(resourceName, verifiedDomain, emailAddress),
				Check:  base.TrustedEmailAddress_GetIDs(resourceFullName, &environmentID, &emailDomainID, &trustedEmailAddressID),
			},
			{
				PreConfig: func() {
					base.TrustedEmailAddress_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID, emailDomainID, trustedEmailAddressID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			// // Test removal of the email domain
			// {
			// 	Config: testAccTrustedEmailAddressConfig_New_DomainVerified(resourceName, verifiedDomain, emailAddress),
			// 	Check:  base.TestAccGetTrustedEmailAddressIDs(resourceFullName, &environmentID, &emailDomainID, &trustedEmailAddressID),
			// },
			// // Replan after removal preconfig
			// {
			// 	PreConfig: func() {
			// 		base.TrustedEmailAddress_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID, emailDomainID, trustedEmailAddressID)
			// 	},
			// 	RefreshState:       true,
			// 	ExpectNonEmptyPlan: true,
			// },
			// // Test removal of the environment
			// {
			// 	Config: testAccApplicationConfig_NewEnv(environmentName, licenseID, resourceName, name),
			// 	Check:  sso.TestAccGetApplicationIDs(resourceFullName, &environmentID, &applicationID),
			// },
			// {
			// 	PreConfig: func() {
			// 		base.Environment_RemovalDrift_PreConfig(ctx, p1Client.API.ManagementAPIClient, t, environmentID)
			// 	},
			// 	RefreshState:       true,
			// 	ExpectNonEmptyPlan: true,
			// },
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
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
			acctest.PreCheckDomainVerification(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.TrustedEmailAddress_CheckDestroy,
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
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNewEnvironment(t)
			acctest.PreCheckNoFeatureFlag(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.TrustedEmailAddress_CheckDestroy,
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
			acctest.PreCheckNoTestAccFlaky(t)
			acctest.PreCheckClient(t)
			acctest.PreCheckNoFeatureFlag(t)
			acctest.PreCheckDomainVerification(t)
		},
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             base.TrustedEmailAddress_CheckDestroy,
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
}`, acctestlegacysdk.MinimalSandboxEnvironment(environmentName, licenseID), environmentName, resourceName, domain, emailAddress)
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
